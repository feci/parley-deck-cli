package runner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"parley-deck-cli/internal/agents"
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
	selected := selectedAgents(opts.Idea.Participants, opts.Agents)
	if len(selected) == 0 {
		return Result{AgentID: "implementer", ExitError: "no implementer available in participants"}
	}
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
// agreed fixes from review/consensus.md and update IMPLEMENTATION.md. Unlike
// RunImplementation it edits existing files, so success is exit-code based (no
// single new-artifact validation). opts.Idea.Participants must be [implementer].
func RunFixup(ctx context.Context, opts Options) Result {
	selected := selectedAgents(opts.Idea.Participants, opts.Agents)
	if len(selected) == 0 {
		return Result{AgentID: "implementer", ExitError: "no implementer available in participants"}
	}
	agent := selected[0]
	now := time.Now().UTC()
	agentDir := filepath.Join(opts.Root, protocol.DeckDir, "runs", opts.RunID, "agents", agent.ID)
	if err := os.MkdirAll(agentDir, 0o755); err != nil {
		return Result{AgentID: agent.ID, ExitError: err.Error()}
	}
	stdoutPath := filepath.Join(agentDir, "stdout.log")
	stderrPath := filepath.Join(agentDir, "stderr.log")
	result := Result{AgentID: agent.ID, StartedAt: now, StdoutPath: stdoutPath, StderrPath: stderrPath}

	prompt := BuildFixupPrompt(agent, opts.Idea)
	cctx, cancel := context.WithTimeout(ctx, timeoutForAgent(opts.Timeout, agent))
	defer cancel()
	cmd, cleanup, err := CommandFor(cctx, opts.Root, agent, prompt)
	if cleanup != nil {
		defer cleanup()
	}
	if err != nil {
		result.ExitError = err.Error()
		result.CompletedAt = time.Now().UTC()
		return result
	}
	cmd.Dir = opts.Root
	so, err := os.Create(stdoutPath)
	if err != nil {
		result.ExitError = err.Error()
		return result
	}
	defer so.Close()
	se, err := os.Create(stderrPath)
	if err != nil {
		result.ExitError = err.Error()
		return result
	}
	defer se.Close()
	cmd.Stdout = so
	cmd.Stderr = se
	if agent.PromptMode == agents.PromptStdin {
		cmd.Stdin = strings.NewReader(prompt)
	}
	runErr := cmd.Run()
	result.CompletedAt = time.Now().UTC()
	result.Duration = result.CompletedAt.Sub(now)
	if runErr != nil {
		result.ExitError = runErr.Error()
		if cctx.Err() != nil {
			result.ExitError = cctx.Err().Error()
		}
	} else {
		result.ArtifactOK = true
	}
	eventType := "agent.fixup_finished"
	if !result.ArtifactOK {
		eventType = "agent.fixup_failed"
	}
	_ = opts.Store.Append(store.Event{Time: result.CompletedAt, Type: eventType, Data: map[string]any{"agent": agent.ID, "error": result.ExitError}})
	return result
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

Rules:
- Review the implementation against FINAL.md and IMPLEMENTATION.md (below).
- Create exactly this review file: %s
- Immediately use your file-writing tool to create that file with its required frontmatter; do not narrate instead of writing.
- Do not edit implementation files or other agents' files.
- Use ONLY these severity tags: CRITICAL, MAJOR, MINOR, NIT.
- Each finding states what is wrong, why it matters, and the concrete fix.
- The first line of the file must be exactly "---".
- Return only a short confirmation with the path written.

Required file shape:
---
agent: %s
idea: %s
review-round: %d
date: %s
---

## Summary
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
	selected := selectedAgents(opts.Idea.Participants, opts.Agents)
	if len(selected) == 0 {
		return Result{AgentID: "drafter", ExitError: "no drafter available in participants"}
	}
	return runAgent(ctx, opts, selected[0])
}

// BuildReviewConsensusPrompt is the Phase 7 drafter prompt; it MUST set the
// machine-readable outstanding_agreed_fixes so the unattended loop can decide.
func BuildReviewConsensusPrompt(agent agents.Discovery, idea protocol.IdeaStatus, outputPath, context string) string {
	return fmt.Sprintf(`You are %s, drafting the Phase 7 review consensus for Parley Deck idea %s.

Rules:
- Create exactly this file: %s (overwrite any prior draft).
- Synthesize the review-round findings below into agreed fixes, deferred follow-ups, and dismissed findings.
- The frontmatter MUST include outstanding_agreed_fixes: <integer count of fixes still to apply> and blocked: <true|false>.
- Set outstanding_agreed_fixes to 0 only when nothing remains to fix.
- The first line of the file must be exactly "---".
- Return only a short confirmation with the path written.

Required file shape:
---
idea: %s
review-cycle: 1
drafted-by: %s
date: %s
outstanding_agreed_fixes: 0
blocked: false
---

## Agreed fixes
## Deferred follow-ups
## Dismissed findings
## Signoffs

Review context (IMPLEMENTATION.md + review rounds):
%s
`, agent.ID, idea.Slug, outputPath, idea.Slug, agent.ID, time.Now().Format("2006-01-02"), context)
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
func ValidateReviewArtifact(path, agentID, ideaSlug string, round int) error {
	meta, err := protocol.ReadFrontmatter(path)
	if err != nil {
		return err
	}
	for key, want := range map[string]string{
		"agent":        agentID,
		"idea":         ideaSlug,
		"review-round": strconv.Itoa(round),
	} {
		if got := strings.Trim(strings.TrimSpace(meta[key]), `"'`); got != want {
			return fmt.Errorf("%s frontmatter %s=%q, want %q", path, key, got, want)
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if !strings.Contains(string(data), "## Findings") {
		return fmt.Errorf("%s missing '## Findings'", path)
	}
	return nil
}
