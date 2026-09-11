package app

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/driver"
	"parley-deck-cli/internal/evidence"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/store"
	"parley-deck-cli/internal/track"
)

// driverImplOps is the production driver.ImplOps adapter (driver-impl-phase). It
// reuses the runner Phase 5-8 helpers + review-mode consensus; the driver depends
// only on the driver.ImplOps interface, so internal/driver never imports
// internal/app.
type driverImplOps struct {
	base            runner.Options // the round-01 runOpts (Root, RunID, Store, Agents, Idea, Timeout)
	root            string
	ideaSlug        string
	ideaDir         string
	implementer     string   // FINAL drafter / first participant
	reviewers       []string // non-implementer participants
	drafter         string   // review-consensus drafter (facilitator)
	out             io.Writer
	verificationCLI string // internal process-fixture seam; empty uses this running CLI
}

func newDriverImplOps(base runner.Options, root, ideaSlug, ideaDir string, participants []string, out io.Writer) driver.ImplOps {
	implementer := resolveImplementer(ideaDir, participants)
	// Dedupe to distinct non-implementer IDs (review CF1). Duplicate participant
	// IDs (e.g. [impl, rev, rev]) would otherwise (a) inflate ReviewerCount so the
	// LE-11 `< 2` guard passes with a single real reviewer, and (b) make
	// RunReviewRound launch one goroutine per duplicate, all writing the same
	// agents/<id>/stdout.log and review/round-NN/<id>.md concurrently (a write race).
	var reviewers []string
	seen := make(map[string]bool)
	for _, p := range participants {
		if p != implementer && !seen[p] {
			seen[p] = true
			reviewers = append(reviewers, p)
		}
	}
	// Track-aware reviewer cap (idea track-aware-driver): an EXPLICIT fast/standard
	// track reduces the reviewer set per §4.0 (fast=1, standard=2); absent and
	// deliberation keep all non-implementers. The truncation is deterministic
	// (participant order) and never drops below the non-solo floor of 1.
	if t, present := driver.ReadTrack(ideaDir); present {
		if pol, err := track.PolicyFor(t, present, len(reviewers), driver.ReadAutoImplement(ideaDir), driver.ReadStrictGate(ideaDir)); err == nil && pol.MaxReviewers > 0 && len(reviewers) > pol.MaxReviewers {
			reviewers = reviewers[:pol.MaxReviewers]
		}
	}
	// The review-consensus drafter MUST be a non-implementer so the implementer
	// cannot filter reviewer findings out of the consensus (AF3). Fall back to the
	// implementer only if there are somehow no reviewers.
	drafter := implementer
	if len(reviewers) > 0 {
		drafter = reviewers[0]
	}
	return driverImplOps{
		base: base, root: root, ideaSlug: ideaSlug, ideaDir: ideaDir,
		implementer: implementer, reviewers: reviewers, drafter: drafter, out: out,
	}
}

// resolveImplementer picks the implementer from durable role metadata (D10/AF6):
// IMPLEMENTATION.md `implementer` (on re-entry), else FINAL.md `implementer` /
// `drafted-by`, validated against participants; otherwise participants[0].
func resolveImplementer(ideaDir string, participants []string) string {
	isParticipant := func(id string) bool {
		for _, p := range participants {
			if p == id {
				return true
			}
		}
		return false
	}
	for _, src := range []struct {
		file string
		keys []string
	}{
		{"IMPLEMENTATION.md", []string{"implementer"}},
		{"FINAL.md", []string{"implementer", "drafted-by"}},
	} {
		meta, err := protocol.ReadFrontmatter(filepath.Join(ideaDir, src.file))
		if err != nil {
			continue
		}
		for _, k := range src.keys {
			id := strings.Trim(strings.TrimSpace(meta[k]), `"'`)
			if id != "" && isParticipant(id) {
				return id
			}
		}
	}
	if len(participants) > 0 {
		return participants[0]
	}
	return ""
}

// withParticipants clones the base run options with a narrowed participant set.
func (o driverImplOps) withParticipants(ids ...string) runner.Options {
	opts := o.base
	opts.Idea.Participants = append([]string(nil), ids...)
	return opts
}

// modelOf returns the configured model id for an agent id, or "" if unknown. It
// matches the id directly OR its mapped family (claude-1 -> claude) so the LE-3
// model-diversity gate is not silently disabled for a roster-id deck (review MINOR,
// kimi-1). It does NOT require the agent to be installed — the configured model
// counts for the diversity comparison regardless of discovery.
func (o driverImplOps) modelOf(id string) string {
	family := rosterMappingFor(o.root)[id]
	for _, a := range o.base.Agents {
		if a.ID == id || (family != "" && a.ID == family) {
			return a.Model
		}
	}
	return ""
}

// reviewersShareImplementerModel reports whether every reviewer uses the same model as
// the implementer (LE-3). Returns (false, "") when the implementer's model is unknown
// or there are no reviewers — diversity can't be asserted, so it never fires spuriously.
func (o driverImplOps) reviewersShareImplementerModel() (bool, string) {
	implModel := o.modelOf(o.implementer)
	if implModel == "" || len(o.reviewers) == 0 {
		return false, ""
	}
	for _, r := range o.reviewers {
		if !strings.EqualFold(o.modelOf(r), implModel) { // review fix F6: case-insensitive
			return false, ""
		}
	}
	return true, implModel
}

// checkModelDiversity (LE-3) emits the agent.model_diversity event and either warns
// (default) or, with require_model_diversity, returns an escalation error. Extracted so
// OpenReviewRound stays thin and the event/warn/escalate logic is unit-testable without
// launching reviewers (review fixes F4 + F8).
func (o driverImplOps) checkModelDiversity() error {
	same, model := o.reviewersShareImplementerModel()
	if !same {
		return nil
	}
	required := driver.ReadRequireModelDiversity(o.ideaDir)
	// §4.0: the fast track's single reviewer MUST be model-diverse — a same-model
	// solo checker is the highest rubber-stamp risk (review-01 F4). Make it a hard
	// gate on fast regardless of the frontmatter flag.
	if t, present := driver.ReadTrack(o.ideaDir); present && t == track.Fast {
		required = true
	}
	action := "warn"
	if required {
		action = "escalate"
	}
	if o.base.Store != (store.Store{}) { // F4: best-effort durable event for TUI/state consumers
		_ = o.base.Store.Append(store.Event{
			Time: time.Now().UTC(),
			Type: "agent.model_diversity",
			Data: map[string]any{
				"idea": o.ideaSlug, "implementer": o.implementer,
				"reviewers": append([]string(nil), o.reviewers...),
				"model":     model, "required": required, "action": action,
			},
		})
	}
	if required {
		return fmt.Errorf("require_model_diversity: every reviewer shares the implementer's model %q; refusing to open review (LE-3)", model)
	}
	fmt.Fprintf(o.out, "driver: WARNING model-diversity — every reviewer shares the implementer's model %q; a same-model checker is more likely to rubber-stamp (LE-3). Set require_model_diversity: true to make this a hard gate.\n", model)
	return nil
}

func (o driverImplOps) Implement(ctx context.Context) error {
	fmt.Fprintf(o.out, "driver: implementing via %s ...\n", o.implementer)
	r := runner.RunImplementation(ctx, o.withParticipants(o.implementer))
	if !r.Success() {
		return fmt.Errorf("implementer %s: %s", r.AgentID, r.ExitError)
	}
	return nil
}

func (o driverImplOps) ImplementationStatus() (string, error) {
	meta, err := protocol.ReadFrontmatter(filepath.Join(o.ideaDir, "IMPLEMENTATION.md"))
	if err != nil {
		return "", err
	}
	return strings.Trim(strings.TrimSpace(meta["status"]), `"'`), nil
}

// RunChecks runs the verification gate (LE-4). Resolution order:
//  1. an explicit `checks:` command from 00-prompt frontmatter → `sh -c <command>`;
//  2. else `go test ./...` when the workspace is a Go module;
//  3. else, for a code-writing (auto_implement) idea with no checks → FAIL CLOSED;
//  4. else (design-only, non-Go) → nothing to check.
//
// Because advanceImpl (pre-review) and advanceReview (post-fix-up) both escalate when
// RunChecks fails, step 3 transitively ties the "artifact-wins" fix-up override to a
// real check: a fix-up that wrote a valid-shaped artifact but cannot be verified no
// longer auto-passes (hermes #8).
func (o driverImplOps) RunChecks(ctx context.Context) (bool, string) {
	run := func(name string, cmd *exec.Cmd) (bool, string) {
		fmt.Fprintf(o.out, "driver: running checks (%s) ...\n", name)
		cmd.Dir = o.root
		var buf bytes.Buffer
		cmd.Stdout = &buf
		cmd.Stderr = &buf
		err := cmd.Run()
		return err == nil, buf.String()
	}
	// List-form `checks:` = the completion contract (completion-contracts-evidence-ledger):
	// run each criterion, write the evidence table into IMPLEMENTATION.md, fail closed on
	// any non-zero exit. A malformed list is a hard fail (present but invalid).
	if criteria, isList, err := driver.ReadChecksContract(o.ideaDir); err != nil {
		return false, "checks: contract invalid — " + err.Error()
	} else if isList {
		return o.runChecksContract(ctx, criteria)
	}
	checks := ""
	if meta, err := protocol.ReadFrontmatter(filepath.Join(o.ideaDir, "00-prompt.md")); err == nil {
		checks = strings.TrimSpace(strings.Trim(strings.TrimSpace(meta["checks"]), `"'`))
	}
	if checks != "" {
		return run(checks, exec.CommandContext(ctx, "sh", "-c", checks))
	}
	if _, err := os.Stat(filepath.Join(o.root, "go.mod")); err == nil {
		return run("go test ./...", exec.CommandContext(ctx, "go", "test", "./..."))
	}
	if driver.ReadAutoImplement(o.ideaDir) {
		return false, "no go.mod and no `checks:` configured for a code-writing (auto_implement) idea; set checks: in 00-prompt.md so a fix-up cannot pass unverified"
	}
	return true, "no go.mod and no checks configured; nothing to check"
}

func (o driverImplOps) OpenReviewRound(ctx context.Context, round int) error {
	if len(o.reviewers) == 0 {
		return fmt.Errorf("no non-implementer reviewers available")
	}
	// LE-3 model-diversity: a checker that shares the implementer's model is more
	// likely to rubber-stamp. Default = warn (+event); require_model_diversity escalates.
	if err := o.checkModelDiversity(); err != nil {
		return err
	}
	fmt.Fprintf(o.out, "driver: opening review round %d (reviewers: %s) ...\n", round, strings.Join(o.reviewers, ", "))
	// AF5: drop any reviewer artifact that exists but fails validation, so
	// RunReviewRound (Overwrite=false) regenerates it instead of skipping a
	// malformed file forever (which would spin the driver to the deadline).
	dir := filepath.Join(o.ideaDir, "review", roundDirLabel(round))
	for _, reviewer := range o.reviewers {
		path := filepath.Join(dir, reviewer+".md")
		if _, err := os.Stat(path); err == nil {
			if runner.ValidateReviewArtifact(path, reviewer, o.ideaSlug, round) != nil {
				_ = os.Remove(path)
			}
		}
	}
	opts := o.withParticipants(o.reviewers...)
	opts.Round = round
	results := runner.RunReviewRound(ctx, opts)
	failed := 0
	for _, r := range results {
		if r.ExitError != "" {
			failed++
		}
	}
	if failed == len(results) && failed > 0 {
		return fmt.Errorf("all reviewers failed in review round %d", round)
	}
	return nil
}

func (o driverImplOps) ReviewRoundComplete(round int) (bool, error) {
	dir := filepath.Join(o.ideaDir, "review", roundDirLabel(round))
	for _, reviewer := range o.reviewers {
		path := filepath.Join(dir, reviewer+".md")
		if _, err := os.Stat(path); err != nil {
			return false, nil
		}
		if err := runner.ValidateReviewArtifact(path, reviewer, o.ideaSlug, round); err != nil {
			return false, nil
		}
	}
	return len(o.reviewers) > 0, nil
}

func (o driverImplOps) DraftReviewConsensus(ctx context.Context, round int) error {
	fmt.Fprintf(o.out, "driver: drafting review consensus via %s ...\n", o.drafter)
	strict := driver.ReadStrictGate(o.ideaDir)
	opts := o.withParticipants(o.drafter)
	opts.Round = round
	opts.StrictGate = strict // LE-2: emit the close fields under strict_gate
	r := runner.RunReviewConsensus(ctx, opts)
	if !r.Success() {
		return fmt.Errorf("review-consensus drafter %s: %s", r.AgentID, r.ExitError)
	}
	path := filepath.Join(o.ideaDir, "review", "consensus.md")
	if err := runner.ValidateReviewConsensusArtifact(path); err != nil {
		return err
	}
	// Review fix F7: under strict_gate, the drafter MUST emit the close fields. If they
	// are absent, fail fast here rather than spinning fresh rounds to MaxFixupCycles.
	if strict {
		meta, err := protocol.ReadFrontmatter(path)
		if err != nil {
			return err
		}
		if strings.TrimSpace(meta["strict_gate_clean"]) == "" || strings.TrimSpace(meta["closing_review_round"]) == "" {
			return fmt.Errorf("strict_gate: review/consensus.md must set strict_gate_clean and closing_review_round")
		}
	}
	return nil
}

func (o driverImplOps) ReviewStatus() (driver.ReviewStatus, error) {
	summary, err := consensus.Status(o.root, o.ideaSlug, true)
	if err != nil {
		return driver.ReviewStatus{}, err
	}
	meta, err := protocol.ReadFrontmatter(filepath.Join(o.ideaDir, "review", "consensus.md"))
	if err != nil {
		return driver.ReviewStatus{}, err
	}
	rawFixes := strings.Trim(strings.TrimSpace(meta["outstanding_agreed_fixes"]), `"'`)
	fixes, err := strconv.Atoi(rawFixes)
	if err != nil || fixes < 0 {
		return driver.ReviewStatus{}, fmt.Errorf("review/consensus.md outstanding_agreed_fixes=%q is not a non-negative integer", rawFixes)
	}
	blocked := strings.EqualFold(strings.Trim(strings.TrimSpace(meta["blocked"]), `"'`), "true")
	// LE-2 strict_gate close fields (absent/zero on non-strict ideas → harmless defaults).
	strictClean := strings.EqualFold(strings.Trim(strings.TrimSpace(meta["strict_gate_clean"]), `"'`), "true")
	closingRound := 0
	if v := strings.Trim(strings.TrimSpace(meta["closing_review_round"]), `"'`); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			closingRound = n
		}
	}
	return driver.ReviewStatus{
		Summary: summary, OutstandingAgreedFixes: fixes, Blocked: blocked,
		StrictGateClean: strictClean, ClosingReviewRound: closingRound,
		ReviewerCount: len(o.reviewers), // LE-11
	}, nil
}

// discoveryFor returns the discovered agent for an id, matching the id directly OR
// its mapped family (claude-1 -> claude) so the LE-7 goal-check is not silently
// skipped for a roster-id deck (review MINOR, kimi-1).
func (o driverImplOps) discoveryFor(id string) (agents.Discovery, bool) {
	family := rosterMappingFor(o.root)[id]
	for _, a := range o.base.Agents {
		if a.ID == id || (family != "" && a.ID == family) {
			return a, true
		}
	}
	return agents.Discovery{}, false
}

// GoalCheck (LE-7) runs a fresh non-implementer agent (the review drafter) to check the
// FINAL.md acceptance criteria before close, reusing the consult execution path with a
// verdict prompt. Missing, self, failed or ambiguous execution cannot establish
// completion. A textual pass remains defense in depth, not criterion evidence.
func (o driverImplOps) GoalCheck(ctx context.Context) (bool, string) {
	checker := o.drafter
	// CF6: GoalCheck must use a non-implementer checker. The upstream guards
	// (ReviewerCount < 2 under auto; OpenReviewRound under strict) already prevent
	// the drafter==implementer fallback from reaching here, but enforce the contract
	// locally too — never run the implementer as its own goal checker.
	if checker == "" || checker == o.implementer {
		return false, "goal-check has no independent checker"
	}
	agent, err := agents.ResolveParticipant(checker, o.base.Agents, rosterMappingFor(o.root))
	if err != nil {
		return false, "goal-check checker unavailable"
	}
	fmt.Fprintf(o.out, "driver: goal-done check via %s ...\n", checker)
	dir := filepath.Join(o.root, protocol.DeckDir, "runs", o.base.RunID, "agents", checker)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return false, "goal-check cannot create its evidence directory"
	}
	ctx = runner.WithLaunchInfo(ctx, runner.LaunchInfo{RunID: o.base.RunID,
		Idea: o.ideaSlug, Phase: "goal-check", Store: o.base.Store})
	res := runner.RunConsult(ctx, runner.ConsultOptions{
		Root:  o.root,
		Agent: agent,
		// Keep the existing bounded one-shot deadline. A timeout leaves
		// completion unverified and halts instead of silently passing.
		Timeout:    2 * time.Minute,
		Prompt:     runner.BuildGoalCheckPrompt(agent, o.base.Idea),
		StdoutPath: filepath.Join(dir, "goal-check.stdout.log"),
		StderrPath: filepath.Join(dir, "goal-check.stderr.log"),
		Progress:   o.out,
	})
	if res.ExitError != "" || res.AgentExit != 0 {
		return false, "goal-check checker failed; completion is unverified"
	}
	switch parseGoalVerdict(res.Answer) {
	case "FAIL":
		return false, res.Answer
	case "PASS":
		return true, ""
	default:
		return false, "goal-check inconclusive; completion is unverified"
	}
}

// parseGoalVerdict accepts PASS only when every stated verdict is an exact PASS.
// A failure is sticky; any unknown verdict makes a pass ambiguous. A trailing
// format example therefore cannot erase an earlier failure or reservation.
func parseGoalVerdict(answer string) string {
	seenPass, seenFail, ambiguous := false, false, false
	for _, line := range strings.Split(answer, "\n") {
		t := strings.ToUpper(strings.TrimSpace(line))
		// CF2: strip leading markdown / quote wrappers (heading, bold, blockquote,
		// underscore, backtick, single/double quote) so a wrapped marker like
		// `GOAL-CHECK: FAIL` or **GOAL-CHECK:** still matches the prefix.
		t = strings.TrimLeft(t, "#*->_ \t`\"'")
		if !strings.HasPrefix(t, "GOAL-CHECK:") {
			continue
		}
		rest := strings.TrimSpace(strings.TrimPrefix(t, "GOAL-CHECK:"))
		// CF2: a bolded/quoted marker ("**GOAL-CHECK:** FAIL") leaves "** FAIL" in
		// rest — strip the leading wrapper run before the PASS/FAIL prefix check.
		rest = strings.Trim(rest, "*`\"'_ ")
		switch {
		case rest == "PASS":
			seenPass = true
		case strings.HasPrefix(rest, "FAIL"):
			seenFail = true
		default:
			ambiguous = true
		}
	}
	if seenFail {
		return "FAIL"
	}
	if seenPass && !ambiguous {
		return "PASS"
	}
	return ""
}

func (o driverImplOps) RequestReviewSignoffs(ctx context.Context, missing []string) error {
	return requestConsensusSignoffs(ctx, requestSignoffsOptions{
		Root:            o.root,
		IdeaSlug:        o.ideaSlug,
		Review:          true,
		ParticipantsRaw: strings.Join(missing, ","),
		Yes:             true,
	}, o.out, o.out)
}

func (o driverImplOps) Fixup(ctx context.Context, cycle int) error {
	fmt.Fprintf(o.out, "driver: running fix-up cycle %d via %s ...\n", cycle, o.implementer)
	r := runner.RunFixup(ctx, o.withParticipants(o.implementer))
	if !r.Success() {
		return fmt.Errorf("fix-up implementer %s: %s", r.AgentID, r.ExitError)
	}
	return nil
}

// Complete writes IMPLEMENTATION.md status=complete. This is a deterministic file
// write by the orchestrator (NOT an implementer agent), so an implementer cannot
// short-circuit review (consensus D5).
func (o driverImplOps) Complete(ctx context.Context) error {
	return evidence.WithReportWriter(ctx, o.ideaDir, func(_ *evidence.ReportWriter) error { return o.completeWithWriter(ctx) })
}

func (o driverImplOps) completeWithWriter(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	pin := ""
	if o.base.Store.Enabled() {
		cursor, err := driver.LoadCursor(filepath.Join(o.base.Store.Directory(), "driver.json"))
		if err == nil {
			pin = cursor.ChecksContractSHA256
		} else if !os.IsNotExist(err) {
			return err
		}
	}
	contract, err := driver.ObserveChecksContract(o.ideaDir, pin)
	if err != nil {
		return err
	} else if contract != "" {
		if runtime.GOOS == "windows" {
			return fmt.Errorf("independent evidence completion requires a POSIX execution host; Windows runtime is not supported")
		}
		if err := o.requireAcceptedVerification(); err != nil {
			return err
		}
		gate := o.EvidenceCloseGate(o.drafter)
		if !gate.Allowed {
			return fmt.Errorf("independent evidence changed before completion: %s", strings.Join(gate.Reasons, "; "))
		}
	}
	path := filepath.Join(o.ideaDir, "IMPLEMENTATION.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if contract != "" {
		completed, _, err := evidence.TransitionStatusToComplete(data)
		if err != nil {
			return err
		}
		report, err := evidence.Load(o.ideaDir)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(o.root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if err := verifyValidationEvidence(data, rel, report); err != nil {
			return err
		}
		beforeRest, _, err := splitValidationEvidence(data)
		if err != nil {
			return err
		}
		if report.ExtraDigests[rel] != sha256Hex(string(beforeRest)) {
			return fmt.Errorf("implementation changed before completion write")
		}
		afterRest, _, err := splitValidationEvidence(completed)
		if err != nil {
			return err
		}
		if reasons := evidence.VerifyCompletionTransition(report, rel, report.ExtraDigests[rel], afterRest, o.drafter); len(reasons) > 0 {
			return fmt.Errorf("completion status was not independently authorized: %s", strings.Join(reasons, "; "))
		}
		return writeVerificationBytes(path, completed)
	}
	lines := strings.Split(string(data), "\n")
	inFrontmatter := false
	replaced := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			}
			if !replaced {
				lines = append(lines[:i], append([]string{"status: complete"}, lines[i:]...)...)
				replaced = true
			}
			break
		}
		if inFrontmatter && strings.HasPrefix(trimmed, "status:") {
			lines[i] = "status: complete"
			replaced = true
		}
	}
	if !replaced {
		return fmt.Errorf("%s has no frontmatter status field", path)
	}
	return writeVerificationBytes(path, []byte(strings.Join(lines, "\n")))
}

func roundDirLabel(n int) string { return fmt.Sprintf("round-%02d", n) }
