package driver

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"parley-deck-cli/internal/consensus"
)

// ImplOps is the agent-launch seam for Parley Deck Phases 5-8 (consensus D1). The
// driver core depends only on this interface, so the impl/review gate is testable
// with a fake and the driver never imports internal/app. The production adapter
// (app layer) reuses runner.RunImplementation/RunReviewRound/RunFixup/
// RunReviewConsensus + the review-mode consensus package.
type ImplOps interface {
	Implement(ctx context.Context) error                       // Phase 5: one implementer → IMPLEMENTATION.md + code
	ImplementationStatus() (string, error)                     // IMPLEMENTATION.md status: frontmatter
	RunChecks(ctx context.Context) (bool, string)              // build/test gate; (ok, detail)
	OpenReviewRound(ctx context.Context, round int) error      // Phase 6: non-implementer reviewers
	ReviewRoundComplete(round int) (bool, error)               // all reviewer artifacts present + valid
	DraftReviewConsensus(ctx context.Context, round int) error // Phase 7: draft review/consensus.md
	ReviewStatus() (ReviewStatus, error)                       // review-mode triage + outstanding_agreed_fixes
	RequestReviewSignoffs(ctx context.Context, missing []string) error
	Fixup(ctx context.Context, cycle int) error // Phase 8: re-invoke implementer for agreed fixes
	Complete(ctx context.Context) error         // driver writes IMPLEMENTATION.md status=complete
}

// ReviewStatus wraps the review-mode consensus summary plus the machine-readable
// Phase-7 fields the unattended loop decides on (consensus D5).
type ReviewStatus struct {
	Summary                consensus.Summary
	OutstandingAgreedFixes int
	Blocked                bool
	// StrictGateClean + ClosingReviewRound carry the drafter's strict_gate
	// certification (LE-2): the drafter sets strict_gate_clean true only when the
	// round it names in closing_review_round had zero findings of any severity.
	StrictGateClean    bool
	ClosingReviewRound int
}

const (
	ActionImplemented   Action = "implemented"    // ran the implementer → PhaseImpl
	ActionReviewOpened  Action = "review-opened"  // opened a review round → PhaseReview
	ActionReviewDrafted Action = "review-drafted" // drafted review/consensus.md
	ActionFixup         Action = "fixup"          // ran a fix-up cycle → new review round
	ActionComplete      Action = "complete"       // review clean → status=complete → PhaseDone
)

// advanceFinal decides whether to auto-implement after FINAL.md (consensus D3/D6).
func (d *Driver) advanceFinal(ctx context.Context, c Cursor) (Action, Cursor, error) {
	if d.cfg.Impl == nil {
		return ActionSurfaceOnly, c, nil // impl auto-drive not wired (slice 1-2 behavior)
	}
	if fileExists(filepath.Join(d.cfg.IdeaDir, "IMPLEMENTATION.md")) {
		return ActionSurfaceOnly, c, nil // Rebuild moves to PhaseImpl/Review next tick
	}
	if !d.cfg.AutoImplement {
		return ActionSurfaceOnly, c, nil // idea did not opt in / --no-implement → stop at FINAL
	}
	if !gitTreeClean(d.cfg.Root) {
		return ActionEscalated, c, fmt.Errorf("git working tree is dirty; refusing to run the implementer (commit or stash first)")
	}
	if err := d.cfg.Impl.Implement(ctx); err != nil {
		return ActionEscalated, c, fmt.Errorf("implement: %w", err)
	}
	c.Phase = PhaseImpl
	c.IdeaStatus = "implemented"
	c.UpdatedAt = nowRFC3339()
	if err := d.commitCursor(c, ActionImplemented, PhaseFinal); err != nil {
		return ActionEscalated, c, err
	}
	return ActionImplemented, c, nil
}

// advanceImpl validates the implementation, runs the checks gate, and opens the
// first review round (consensus D6).
func (d *Driver) advanceImpl(ctx context.Context, c Cursor) (Action, Cursor, error) {
	if d.cfg.Impl == nil {
		return ActionSurfaceOnly, c, nil
	}
	status, err := d.cfg.Impl.ImplementationStatus()
	if err != nil {
		return ActionEscalated, c, fmt.Errorf("read implementation status: %w", err)
	}
	if !implReadyForReview(status) {
		// A known in-progress status means the implementation attempt is mid-flight
		// (e.g. a crash/re-entry); await rather than fail. An empty/unknown status is
		// fail-closed (AF7/D6).
		if implInProgress(status) {
			return ActionAwait, c, nil
		}
		return ActionEscalated, c, fmt.Errorf("IMPLEMENTATION.md status=%q is neither review-ready nor a known in-progress state", status)
	}
	// Build/test gate before spending reviewer agents (D4); failure escalates.
	if ok, detail := d.cfg.Impl.RunChecks(ctx); !ok {
		return ActionEscalated, c, fmt.Errorf("checks failed before review:\n%s", strings.TrimSpace(detail))
	}
	// PhaseImpl only opens the FIRST review round; subsequent rounds come from the
	// fix-up loop in advanceReview. Re-entry is safe: once review/round-01 exists,
	// Rebuild routes to PhaseReview instead.
	if err := d.cfg.Impl.OpenReviewRound(ctx, 1); err != nil {
		return ActionEscalated, c, fmt.Errorf("open review round 1: %w", err)
	}
	c.Phase = PhaseReview
	c.UpdatedAt = nowRFC3339()
	if err := d.commitCursor(c, ActionReviewOpened, PhaseImpl); err != nil {
		return ActionEscalated, c, err
	}
	return ActionReviewOpened, c, nil
}

// advanceReview runs the Phase 6-8 loop: complete the round → draft review
// consensus → signoffs → done | fix-up | escalate (consensus D5/D6).
func (d *Driver) advanceReview(ctx context.Context, c Cursor) (Action, Cursor, error) {
	if d.cfg.Impl == nil {
		return ActionSurfaceOnly, c, nil
	}
	round := highestReviewRound(d.cfg.IdeaDir)
	if round < 1 {
		round = 1
	}
	// AF2 crash-idempotency: if this round's fix-up already completed (a crash after
	// Fixup+RunChecks but before the next round opened), finish the transition
	// without re-running Fixup or re-drafting the consensus.
	fixupMarker := filepath.Join(d.cfg.IdeaDir, "review", roundLabel(round), ".fixup-done")
	if fileExists(fixupMarker) {
		if err := d.archiveReviewConsensus(round); err != nil {
			return ActionEscalated, c, fmt.Errorf("archive review consensus: %w", err)
		}
		if err := d.cfg.Impl.OpenReviewRound(ctx, round+1); err != nil {
			return ActionEscalated, c, fmt.Errorf("open review round %d: %w", round+1, err)
		}
		c.Phase = PhaseReview
		c.UpdatedAt = nowRFC3339()
		if err := d.commitCursor(c, ActionFixup, PhaseReview); err != nil {
			return ActionEscalated, c, err
		}
		return ActionFixup, c, nil
	}
	complete, err := d.cfg.Impl.ReviewRoundComplete(round)
	if err != nil {
		return ActionEscalated, c, fmt.Errorf("check review round %d: %w", round, err)
	}
	if !complete {
		// Fill any missing reviewer artifacts (RunReviewRound is Overwrite=false), then
		// await; the loop deadline bounds a reviewer that never finishes.
		if err := d.cfg.Impl.OpenReviewRound(ctx, round); err != nil {
			return ActionEscalated, c, fmt.Errorf("re-open review round %d: %w", round, err)
		}
		if done, _ := d.cfg.Impl.ReviewRoundComplete(round); !done {
			return ActionAwait, c, nil
		}
	}

	reviewConsensus := filepath.Join(d.cfg.IdeaDir, "review", "consensus.md")
	if !fileExists(reviewConsensus) {
		if err := d.cfg.Impl.DraftReviewConsensus(ctx, round); err != nil {
			return ActionEscalated, c, fmt.Errorf("draft review consensus: %w", err)
		}
		return ActionReviewDrafted, c, nil
	}

	rs, err := d.cfg.Impl.ReviewStatus()
	if err != nil {
		return ActionEscalated, c, fmt.Errorf("read review status: %w", err) // malformed → escalate
	}
	if rs.Blocked {
		return ActionEscalated, c, fmt.Errorf("review consensus is blocked; escalating")
	}
	if len(rs.Summary.Missing) > 0 {
		if err := d.cfg.Impl.RequestReviewSignoffs(ctx, rs.Summary.Missing); err != nil {
			return ActionEscalated, c, fmt.Errorf("request review signoffs: %w", err)
		}
		after, err := d.cfg.Impl.ReviewStatus()
		if err != nil {
			return ActionEscalated, c, fmt.Errorf("read review status after signoffs: %w", err)
		}
		if len(after.Summary.Missing) > 0 {
			return ActionEscalated, c, fmt.Errorf("review signoffs still missing: %s", strings.Join(after.Summary.Missing, ", "))
		}
		rs = after
	}
	if rs.Summary.Triage != consensus.TriageReady && rs.Summary.Triage != consensus.TriageReserved {
		return ActionEscalated, c, fmt.Errorf("review consensus triage=%q; cannot complete", rs.Summary.Triage)
	}

	if rs.OutstandingAgreedFixes == 0 {
		if d.cfg.StrictGate {
			// strict_gate (LE-2): completion requires a FRESH full-scope closing review
			// round with zero findings of any severity — certified by the drafter
			// (strict_gate_clean + closing_review_round for THIS round) and not
			// contradicted by a deterministic finding scan. The scan can only veto a
			// clean claim (fail closed); it never auto-passes one.
			certifiedClean := rs.StrictGateClean && rs.ClosingReviewRound == round
			if certifiedClean && reviewRoundHasFindings(d.cfg.IdeaDir, round) {
				return ActionEscalated, c, fmt.Errorf("strict_gate: drafter certified round %d clean but a review file has a real finding; escalating", round)
			}
			if !certifiedClean {
				// Not a certified-clean closing round yet → run one more fresh review
				// round. Bounded by MaxFixupCycles so the strict-close loop terminates.
				if round >= d.cfg.MaxFixupCycles {
					return ActionEscalated, c, fmt.Errorf("strict_gate: no clean closing review round after %d round(s) (MaxFixupCycles=%d); escalating", round, d.cfg.MaxFixupCycles)
				}
				if err := d.archiveReviewConsensus(round); err != nil {
					return ActionEscalated, c, fmt.Errorf("archive review consensus: %w", err)
				}
				if err := d.cfg.Impl.OpenReviewRound(ctx, round+1); err != nil {
					return ActionEscalated, c, fmt.Errorf("open strict closing review round %d: %w", round+1, err)
				}
				c.Phase = PhaseReview
				c.UpdatedAt = nowRFC3339()
				if err := d.commitCursor(c, ActionReviewOpened, PhaseReview); err != nil {
					return ActionEscalated, c, err
				}
				return ActionReviewOpened, c, nil
			}
		}
		// DONE (D5): the driver — not the implementer — writes status=complete.
		if err := d.cfg.Impl.Complete(ctx); err != nil {
			return ActionEscalated, c, fmt.Errorf("mark complete: %w", err)
		}
		c.Phase = PhaseDone
		c.IdeaStatus = "complete"
		c.UpdatedAt = nowRFC3339()
		if err := d.commitCursor(c, ActionComplete, PhaseReview); err != nil {
			return ActionEscalated, c, err
		}
		return ActionComplete, c, nil
	}

	// Outstanding agreed fixes → a fix-up cycle, bounded by MaxFixupCycles.
	cycle := round
	if cycle >= d.cfg.MaxFixupCycles {
		return ActionEscalated, c, fmt.Errorf("review still has %d agreed fixes after %d cycle(s) (MaxFixupCycles=%d); escalating", rs.OutstandingAgreedFixes, cycle, d.cfg.MaxFixupCycles)
	}
	if !d.cfg.AutoImplement {
		return ActionEscalated, c, fmt.Errorf("review has agreed fixes but auto_implement (code-writing) is not enabled; escalating")
	}
	if !gitTreeClean(d.cfg.Root) {
		return ActionEscalated, c, fmt.Errorf("git working tree is dirty; refusing to run a fix-up")
	}
	if err := d.cfg.Impl.Fixup(ctx, cycle); err != nil {
		return ActionEscalated, c, fmt.Errorf("fix-up cycle %d: %w", cycle, err)
	}
	// AF1: re-run the build/test gate after the fix-up, before spending reviewers on
	// the next round; a fix-up that broke the build escalates immediately.
	if ok, detail := d.cfg.Impl.RunChecks(ctx); !ok {
		return ActionEscalated, c, fmt.Errorf("checks failed after fix-up cycle %d before review:\n%s", cycle, strings.TrimSpace(detail))
	}
	// AF2: mark the fix-up done so a crash before the next round opens does not
	// re-run Fixup on re-entry.
	if err := os.WriteFile(fixupMarker, []byte(nowRFC3339()+"\n"), 0o644); err != nil {
		return ActionEscalated, c, fmt.Errorf("write fix-up marker: %w", err)
	}
	// Archive the just-reviewed consensus and open the next review round so the next
	// tick drafts a fresh consensus for round N+1 (consensus D8 + agy's archiving).
	if err := d.archiveReviewConsensus(round); err != nil {
		return ActionEscalated, c, fmt.Errorf("archive review consensus: %w", err)
	}
	if err := d.cfg.Impl.OpenReviewRound(ctx, round+1); err != nil {
		return ActionEscalated, c, fmt.Errorf("open review round %d: %w", round+1, err)
	}
	c.Phase = PhaseReview
	c.UpdatedAt = nowRFC3339()
	if err := d.commitCursor(c, ActionFixup, PhaseReview); err != nil {
		return ActionEscalated, c, err
	}
	return ActionFixup, c, nil
}

// archiveReviewConsensus moves the active review/consensus.md to
// review/round-NN/consensus.md before the next cycle re-drafts the root file,
// preserving a per-cycle audit trail (and clearing the root so the gate re-drafts).
func (d *Driver) archiveReviewConsensus(round int) error {
	src := filepath.Join(d.cfg.IdeaDir, "review", "consensus.md")
	if !fileExists(src) {
		return nil
	}
	dst := filepath.Join(d.cfg.IdeaDir, "review", roundLabel(round), "consensus.md")
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	if fileExists(dst) {
		if err := os.Remove(dst); err != nil {
			return err
		}
	}
	return os.Rename(src, dst)
}

func implReadyForReview(status string) bool {
	switch {
	case status == "implemented", status == "ready-for-review":
		return true
	case strings.HasPrefix(status, "fix-up-cycle"):
		return true
	}
	return false
}

// reviewRoundHasFindings scans a review round's reviewer files for at least one real
// finding — a "### [CRITICAL|MAJOR|MINOR|NIT] <title>" heading whose title is non-empty
// and not the literal "<title>" placeholder from the prompt template. It is a
// deterministic veto on a strict_gate_clean claim (LE-2): it can only fail a clean
// claim closed, never auto-pass one.
func reviewRoundHasFindings(ideaDir string, round int) bool {
	dir := filepath.Join(ideaDir, "review", roundLabel(round))
	entries, err := os.ReadDir(dir)
	if err != nil {
		// Fail closed (review fix F1): an absent dir means nothing to veto
		// (ReviewRoundComplete guards presence). But an unreadable dir that should
		// exist must VETO (return true → escalate), never silently auto-pass — the
		// scan's stated invariant is "veto-only, never auto-pass."
		return !errors.Is(err, fs.ErrNotExist)
	}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".md") || name == "_index.md" || name == "consensus.md" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return true // F1: a file confirmed present but now unreadable → fail closed (veto)
		}
		if scanHasRealFinding(string(data)) {
			return true
		}
	}
	return false
}

// scanHasRealFinding reports whether the content has a severity-tagged finding heading.
// It is whitespace-tolerant and case-insensitive on the severity tag (review fix F2), and
// counts an empty-title heading (finding text on the next line) as a finding (review fix
// F3); only the literal template placeholder "<title>" is ignored.
func scanHasRealFinding(content string) bool {
	for _, line := range strings.Split(content, "\n") {
		t := strings.TrimSpace(line)
		if !strings.HasPrefix(t, "###") {
			continue
		}
		rest := strings.TrimSpace(strings.TrimPrefix(t, "###"))
		if !strings.HasPrefix(rest, "[") {
			continue
		}
		closeIdx := strings.Index(rest, "]")
		if closeIdx < 0 {
			continue
		}
		switch strings.ToUpper(strings.TrimSpace(rest[1:closeIdx])) {
		case "CRITICAL", "MAJOR", "MINOR", "NIT":
		default:
			continue
		}
		if strings.TrimSpace(rest[closeIdx+1:]) != "<title>" {
			return true // empty or concrete title both count; only the placeholder is ignored
		}
	}
	return false
}

func implInProgress(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "in-progress", "in_progress", "wip", "draft", "drafting":
		return true
	}
	return false
}

// gitTreeClean reports whether the workspace git tree has no uncommitted changes.
// A non-git workspace (or missing git) is treated as clean — there is nothing to
// clobber. But INSIDE a real repo, a git error (e.g. a stale index.lock) is treated
// as dirty/unsafe rather than clean, so we never run a code-writing agent on a
// possibly-dirty tree (consensus D3/AF9).
func gitTreeClean(root string) bool {
	// GIT_OPTIONAL_LOCKS=0: read-only probes must never take optional locks
	// (index refresh writes .git) on the weakly-coherent virtio-fs mount
	// (runner-hardening-kindly D8).
	probe := func(args ...string) *exec.Cmd {
		cmd := exec.Command("git", args...)
		cmd.Env = append(os.Environ(), "GIT_OPTIONAL_LOCKS=0")
		return cmd
	}
	if err := probe("-C", root, "rev-parse", "--is-inside-work-tree").Run(); err != nil {
		return true // not inside a git work tree (or git missing) → nothing to clobber
	}
	out, err := probe("-C", root, "status", "--porcelain").Output()
	if err != nil {
		return false // inside a repo but the status probe failed → assume unsafe
	}
	return strings.TrimSpace(string(out)) == ""
}
