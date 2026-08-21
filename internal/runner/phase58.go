package runner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/fsutil"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/store"
)

// RunImplementation runs Phase 5: it launches a single implementer to produce
// IMPLEMENTATION.md (and code on a branch) per FINAL.md. opts.Idea.Participants
// must contain exactly the implementer. Reuses the shared launch machinery.
func RunImplementation(ctx context.Context, opts Options) Result {
	opts.Phase = "implementation"
	opts.ArtifactName = "IMPLEMENTATION.md"
	if opts.RoundLabel == "" {
		opts.RoundLabel = "implementation"
	}
	selected, _ := selectedAgents(opts.Idea.Participants, opts.Agents, resolveMapping(opts))
	if len(selected) == 0 {
		return Result{AgentID: "implementer", ExitError: "no implementer available in participants"}
	}
	opts.SegmentID = appendSegmentStarted(opts, "continue", agentIDs(selected))
	return runAgent(ctx, opts, selected[0])
}

// RunReviewRound runs Phase 6 review round N: each reviewer (opts.Idea.Participants
// should be the non-implementer reviewers) writes review/round-NN/<agent>.md.
func RunReviewRound(ctx context.Context, opts Options) []Result {
	if opts.Round < 1 {
		opts.Round = 1
	}
	opts.Phase = "review"
	opts.RoundLabel = filepath.Join("review", roundLabel(opts.Round))
	return RunRoundOne(ctx, opts)
}

// RunFixup runs a Phase 8 fix-up: it re-invokes the implementer to apply the
// agreed fixes from review/consensus.md and update IMPLEMENTATION.md. Success
// requires the updated IMPLEMENTATION.md to validate (ValidateFixupArtifact);
// an ordinary nonzero exit with a valid artifact succeeds with agent_exit
// (consensus D7). opts.Idea.Participants must be [implementer].
func RunFixup(ctx context.Context, opts Options) Result {
	selected, _ := selectedAgents(opts.Idea.Participants, opts.Agents, resolveMapping(opts))
	if len(selected) == 0 {
		return Result{AgentID: "implementer", ExitError: "no implementer available in participants"}
	}
	agent := selected[0]
	opts.SegmentID = appendSegmentStarted(opts, "retry", []string{agent.ID})
	now := time.Now().UTC()
	agentDir := filepath.Join(opts.Root, protocol.DeckDir, "runs", opts.RunID, "agents", agent.ID)
	if err := fsutil.MkdirAllResilient(agentDir, 0o755); err != nil {
		return Result{AgentID: agent.ID, ExitError: err.Error()}
	}
	stdoutPath := filepath.Join(agentDir, "stdout.log")
	stderrPath := filepath.Join(agentDir, "stderr.log")
	result := Result{AgentID: agent.ID, StartedAt: now, StdoutPath: stdoutPath, StderrPath: stderrPath}

	prompt := BuildFixupPrompt(agent, opts.Idea)
	hardTimeout := timeoutForAgent(opts.Timeout, agent)
	cctx, cancel := context.WithTimeout(ctx, hardTimeout)
	defer cancel()

	// The fix-up runs through the same hardened exec path as every other agent
	// launch (review fix 4): process group + procctl marker + participant env
	// shedding + counting writers + supervised wait. No retry — a code-mutating
	// phase is not safely re-runnable after a watchdog kill.
	act := &activityTracker{}
	cfg := supervisionForAgent(agent, hardTimeout)
	hooks := supervisionHooks{
		onHeartbeat: func(snap activitySnapshot, elapsed time.Duration) {
			_ = opts.Store.Append(store.Event{
				Time: time.Now().UTC(),
				Type: "agent.heartbeat",
				Data: map[string]any{
					"agent": agent.ID, "segment_id": opts.SegmentID, "attempt_id": 1,
					"phase": "fixup", "launch": agents.LaunchHeadless,
					"elapsed_ms": elapsed.Milliseconds(), "timeout_ms": hardTimeout.Milliseconds(),
					"stdout_bytes": snap.StdoutBytes, "stderr_bytes": snap.StderrBytes,
					"last_activity_ms_ago": activityAgeMS(snap),
				},
			})
		},
		onWatchdog: func(kind string, snap activitySnapshot, elapsed time.Duration) {
			_ = opts.Store.Append(store.Event{
				Time: time.Now().UTC(),
				Type: "agent." + kind,
				Data: map[string]any{
					"agent": agent.ID, "segment_id": opts.SegmentID, "attempt_id": 1,
					"elapsed_ms":   elapsed.Milliseconds(),
					"stdout_bytes": snap.StdoutBytes, "stderr_bytes": snap.StderrBytes,
					"action": "failed",
				},
			})
		},
	}
	_, runErr := execAgentProcess(cctx, opts.Root, opts.RunID, agent.ID, opts.RunID+":"+agent.ID+":fixup", agent, prompt, stdoutPath, stderrPath, nil, act, cfg, hooks)
	result.CompletedAt = time.Now().UTC()
	result.Duration = result.CompletedAt.Sub(now)
	watchdog := ""
	switch {
	case errors.Is(runErr, errNoFirstOutput):
		watchdog = "no_first_output"
	case errors.Is(runErr, errStalled):
		watchdog = "stalled"
	}
	cancelled := cctx.Err() != nil || watchdog != ""
	if runErr != nil {
		result.ExitError = runErr.Error()
		if cctx.Err() != nil {
			result.ExitError = cctx.Err().Error()
		}
	}

	// Fix-up success is no longer exit-code based (consensus D7): the updated
	// IMPLEMENTATION.md must validate, and only then can an ordinary nonzero
	// exit be overridden (artifact-wins). Timeouts/cancellations always fail.
	validateErr := ValidateFixupArtifact(opts.Idea.Path, opts.Idea.Slug)
	if validateErr == nil {
		result.ArtifactOK = true
		var exitErr *exec.ExitError
		if !cancelled && runErr != nil && errors.As(runErr, &exitErr) {
			result.AgentExit = exitErr.ExitCode()
			result.AgentExitKind = "exec"
			result.ExitError = ""
		}
	} else if !cancelled && runErr == nil {
		result.ExitError = combineError(result.ExitError, validateErr)
	}

	failed := result.ExitError != "" || !result.ArtifactOK
	data := map[string]any{"agent": agent.ID, "error": result.ExitError, "segment_id": opts.SegmentID}
	eventType := "agent.fixup_finished"
	if failed {
		eventType = "agent.fixup_failed"
		class, hint := terminalFailureClass(watchdog, cctx.Err(), false, stderrPath, stdoutPath, result.ExitError)
		result.FailureClass = class
		result.RecoveryHint = hint
		data["failure_class"] = class
		data["recovery_hint"] = hint
	} else if result.AgentExitKind != "" {
		data["agent_exit"] = result.AgentExit
		data["agent_exit_kind"] = result.AgentExitKind
	}
	_ = opts.Store.Append(store.Event{Time: result.CompletedAt, Type: eventType, Data: data})
	return result
}

// ValidateFixupArtifact checks that a fix-up cycle left IMPLEMENTATION.md in a
// reviewable state: the file exists, its frontmatter idea matches, its status
// is review-ready, and a fix-up section is present (consensus D7).
func ValidateFixupArtifact(ideaPath, ideaSlug string) error {
	implPath := filepath.Join(ideaPath, "IMPLEMENTATION.md")
	data, err := os.ReadFile(implPath)
	if err != nil {
		return fmt.Errorf("fix-up validation: %w", err)
	}
	meta, err := protocol.ReadFrontmatter(implPath)
	if err != nil {
		return fmt.Errorf("fix-up validation: read frontmatter: %w", err)
	}
	if got := strings.TrimSpace(meta["idea"]); got != "" && got != ideaSlug {
		return fmt.Errorf("fix-up validation: IMPLEMENTATION.md frontmatter idea=%q, want %q", got, ideaSlug)
	}
	status := strings.TrimSpace(meta["status"])
	reviewReady := status == "implemented" || status == "ready-for-review" || strings.HasPrefix(status, "fix-up-cycle")
	if !reviewReady {
		return fmt.Errorf("fix-up validation: IMPLEMENTATION.md status=%q is not review-ready", status)
	}
	if !strings.Contains(string(data), "## Fix-up cycle") {
		return fmt.Errorf("fix-up validation: IMPLEMENTATION.md has no \"## Fix-up cycle\" section")
	}
	return nil
}

// BuildFixupPrompt is the Phase 8 fix-up prompt for the implementer.
func BuildFixupPrompt(agent agents.Discovery, idea protocol.IdeaStatus) string {
	rc := filepath.Join(idea.Path, "review", "consensus.md")
	impl := filepath.Join(idea.Path, "IMPLEMENTATION.md")
	return fmt.Sprintf(`You are %s, the implementer for Parley Deck idea %s (Phase 8 fix-up).

Rules:
- Read %s and apply EVERY agreed fix listed there to the code.
- Then append a "## Fix-up cycle" section to %s describing the fixes applied, and keep its frontmatter status accurate.
- Do not edit other agents' review or signoff files.
- Run the project's build/tests if available.
- Return a short confirmation with the fixes applied and checks run.
`, agent.ID, idea.Slug, rc, impl)
}

// BuildImplementationPrompt is the Phase 5 implementer prompt.
func BuildImplementationPrompt(agent agents.Discovery, idea protocol.IdeaStatus, outputPath string) string {
	finalPath := filepath.Join(idea.Path, "FINAL.md")
	return fmt.Sprintf(`You are %s, the implementer for Parley Deck idea %s (Phase 5).

Rules:
- Implement strictly according to %s. Read it first.
- Before multi-file changes or changes outside parley-deck/, write your plan/checklist into the file below.
- Create exactly this file: %s
- Immediately use your file-writing tool to create that file; do not explore the workspace first beyond what FINAL.md requires, and report a blocker instead of narrating.
- Record unavoidable deviations from FINAL.md in that file.
- Do not edit other agents' protocol files.
- The first line of the file must be exactly "---".
- Return a short confirmation with branch, files changed, checks run.

Required file shape:
---
idea: %s
status: implemented
implementer: %s
started: %s
completed: %s
branch: <repo-path>#<branch-name>
head-commit: <sha-or-short-sha>
---

## Summary of work
## Implementation plan / checklist
## Deviations from FINAL.md
## Notes for reviewers
`, agent.ID, idea.Slug, finalPath, outputPath, idea.Slug, agent.ID, time.Now().Format("2006-01-02"), time.Now().Format("2006-01-02"))
}

// BuildReviewPrompt is the Phase 6 reviewer prompt for review round N.
func BuildReviewPrompt(agent agents.Discovery, idea protocol.IdeaStatus, round int, outputPath, context string) string {
	return fmt.Sprintf(`You are %s, a reviewer for Parley Deck idea %s (Phase 6, review round %d).

Refutation-default posture: assume the implementation is WRONG until you fail to break it.
For each observable acceptance criterion in FINAL.md, actively try to construct a failing
case or run the relevant check. Report "no findings" only after stating, under
"## Refutation attempts", what you tried that failed to break it.

Rules:
- Review the implementation against FINAL.md and IMPLEMENTATION.md (below).
- Create exactly this review file: %s
- Immediately use your file-writing tool to create that file with its required frontmatter; do not narrate instead of writing.
- Do not edit implementation files or other agents' files.
- Use ONLY these severity tags: CRITICAL, MAJOR, MINOR, NIT.
- Each finding states what is wrong, why it matters, and the concrete fix.
- Under "## Refutation attempts", record the criteria/cases you tried to break and the result; an empty "## Findings" is only credible when refutation attempts are recorded.
- The first line of the file must be exactly "---".
- Record the exact tree you reviewed in "reviewed-commit:" — run 'git rev-parse HEAD' and paste it.
  A review that does not name its tree cannot be told apart from a stale one later.
- Return only a short confirmation with the path written.

Required file shape:
---
agent: %s
idea: %s
review-round: %d
reviewed-commit: <output of 'git rev-parse HEAD'>
date: %s
---

## Summary
## Refutation attempts
## Findings
### [CRITICAL] <title>
### [MAJOR] <title>
### [MINOR] <title>
### [NIT] <title>
## Open questions

Context (FINAL.md, IMPLEMENTATION.md, prior review rounds):
%s
`, agent.ID, idea.Slug, round, outputPath, agent.ID, idea.Slug, round, time.Now().Format("2006-01-02"), context)
}

// gatherReviewContext concatenates FINAL.md, IMPLEMENTATION.md, and any prior
// review-round artifacts to seed a code-review round.
func gatherReviewContext(ideaPath string, round int) (string, error) {
	var b strings.Builder
	for _, name := range []string{"FINAL.md", "IMPLEMENTATION.md"} {
		data, err := os.ReadFile(filepath.Join(ideaPath, name))
		if err == nil {
			fmt.Fprintf(&b, "\n===== %s =====\n%s\n", name, string(data))
		}
	}
	for r := 1; r < round; r++ {
		dir := filepath.Join(ideaPath, "review", roundLabel(r))
		entries, err := os.ReadDir(dir)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return "", err
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") || e.Name() == "_index.md" {
				continue
			}
			data, err := os.ReadFile(filepath.Join(dir, e.Name()))
			if err != nil {
				return "", err
			}
			fmt.Fprintf(&b, "\n===== review/%s/%s =====\n%s\n", roundLabel(r), e.Name(), string(data))
		}
	}
	return b.String(), nil
}

// validateArtifactForPhase dispatches artifact validation by phase.
func validateArtifactForPhase(opts Options, outputPath, agentID string) error {
	switch opts.Phase {
	case "implementation":
		return ValidateImplementationArtifact(outputPath, opts.Idea.Slug)
	case "review":
		return ValidateReviewArtifact(outputPath, agentID, opts.Idea.Slug, roundNumber(opts))
	case "review-consensus":
		return ValidateReviewConsensusArtifact(outputPath)
	}
	return ValidateRoundArtifact(outputPath, agentID, opts.Idea.Slug, roundNumber(opts))
}

// RunReviewConsensus launches the drafter to (re)write review/consensus.md with
// the machine-readable Phase-7 contract (outstanding_agreed_fixes). Overwrites
// any prior draft so each fix-up cycle records the current count.
func RunReviewConsensus(ctx context.Context, opts Options) Result {
	opts.Phase = "review-consensus"
	opts.ArtifactName = filepath.Join("review", "consensus.md")
	opts.Overwrite = true
	selected, _ := selectedAgents(opts.Idea.Participants, opts.Agents, resolveMapping(opts))
	if len(selected) == 0 {
		return Result{AgentID: "drafter", ExitError: "no drafter available in participants"}
	}
	opts.SegmentID = appendSegmentStarted(opts, "continue", agentIDs(selected))
	return runAgent(ctx, opts, selected[0])
}

// BuildReviewConsensusPrompt is the Phase 7 drafter prompt; it MUST set the
// machine-readable outstanding_agreed_fixes so the unattended loop can decide.
func BuildReviewConsensusPrompt(agent agents.Discovery, idea protocol.IdeaStatus, outputPath, context string, strictGate bool) string {
	strictRule := ""
	strictFields := ""
	if strictGate {
		// LE-2: under strict_gate the close decision needs a certified-clean closing
		// round, expressed via two machine-readable fields the driver checks.
		strictRule = "- strict_gate is ON: also set closing_review_round: <the review round number this consensus certifies> and strict_gate_clean: <true ONLY if EVERY reviewer in that round reported zero findings of any severity>.\n"
		strictFields = "closing_review_round: 0\nstrict_gate_clean: false\n"
	}
	return fmt.Sprintf(`You are %s, drafting the Phase 7 review consensus for Parley Deck idea %s.

Rules:
- Create exactly this file: %s (overwrite any prior draft).
- Synthesize the review-round findings below into agreed fixes, deferred follow-ups, and dismissed findings.
- The frontmatter MUST include outstanding_agreed_fixes: <integer count of fixes still to apply> and blocked: <true|false>.
- Set outstanding_agreed_fixes to 0 only when nothing remains to fix.
%s- The first line of the file must be exactly "---".
- Return only a short confirmation with the path written.

Required file shape:
---
idea: %s
review-cycle: 1
drafted-by: %s
date: %s
outstanding_agreed_fixes: 0
blocked: false
%s---

## Agreed fixes
## Deferred follow-ups
## Dismissed findings
## Signoffs

Review context (IMPLEMENTATION.md + review rounds):
%s
`, agent.ID, idea.Slug, outputPath, strictRule, idea.Slug, agent.ID, time.Now().Format("2006-01-02"), strictFields, context)
}

// ValidateReviewConsensusArtifact requires the machine-readable contract field.
func ValidateReviewConsensusArtifact(path string) error {
	meta, err := protocol.ReadFrontmatter(path)
	if err != nil {
		return err
	}
	if strings.TrimSpace(meta["outstanding_agreed_fixes"]) == "" {
		return fmt.Errorf("%s frontmatter missing outstanding_agreed_fixes (the auto fix-up loop fails closed without it)", path)
	}
	return nil
}

// ValidateImplementationArtifact checks an IMPLEMENTATION.md.
func ValidateImplementationArtifact(path, ideaSlug string) error {
	meta, err := protocol.ReadFrontmatter(path)
	if err != nil {
		return err
	}
	if got := strings.Trim(strings.TrimSpace(meta["idea"]), `"'`); got != ideaSlug {
		return fmt.Errorf("%s frontmatter idea=%q, want %q", path, got, ideaSlug)
	}
	if strings.TrimSpace(meta["status"]) == "" {
		return fmt.Errorf("%s frontmatter missing status", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if !strings.Contains(string(data), "## Summary of work") {
		return fmt.Errorf("%s missing '## Summary of work'", path)
	}
	return nil
}

// ValidateReviewArtifact checks a review/round-NN/<agent>.md.
//
// Deprecated shim: the rule lives in internal/protocol so the manual review-consensus path
// enforces exactly the same contract as the driver.
func ValidateReviewArtifact(path, agentID, ideaSlug string, round int) error {
	return protocol.ValidateReviewArtifact(path, agentID, ideaSlug, round)
}
