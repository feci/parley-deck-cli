package driver

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/protocol"
)

// ConsensusOps is the agent-launch seam for the consensus phase (consensus D6/D9).
// The driver core depends only on this interface, so the consensus gate is
// unit-testable with a fake and the driver never imports internal/app. The
// production adapter (provided by the app layer) reuses the existing
// request-signoffs path and a FINAL/consensus content drafter.
//
// All methods are real agent launches or pure consensus-package calls; they must
// be idempotent enough for re-entry (e.g. Draft is a no-op if consensus.md exists).
type ConsensusOps interface {
	// Status reads the current consensus triage from disk.
	Status() (consensus.Summary, error)
	// Draft authors consensus.md content from the completed rounds (real agent),
	// then sets idea status to "consensus". No-op if consensus.md already exists.
	Draft(ctx context.Context) error
	// RequestSignoffs invokes each missing participant to append its own signoff.
	RequestSignoffs(ctx context.Context, missing []string) error
	// DraftFinal authors FINAL.md content from the accepted consensus (real agent).
	DraftFinal(ctx context.Context) error
	// Reopen runs consensus.Reopen with the BLOCK reason (the driver then opens the
	// next cross-review round and invalidates the stale consensus.md/FINAL.md).
	Reopen(ctx context.Context, reason string) error
}

// advanceConsensus runs one tick of the consensus gate (consensus D6/D7). The
// caller has already confirmed the transport/auto gate and that the cursor is in
// PhaseConsensus (consensus.md exists). Returns the (possibly updated) cursor.
func (d *Driver) advanceConsensus(ctx context.Context, c Cursor) (Action, Cursor, error) {
	if d.cfg.Consensus == nil {
		return ActionSurfaceOnly, c, nil // gate not wired (slice-1 behavior)
	}
	summary, err := d.cfg.Consensus.Status()
	if err != nil {
		return ActionEscalated, c, fmt.Errorf("read consensus status: %w", err)
	}
	switch summary.Triage {
	case consensus.TriageReady, consensus.TriageReserved:
		finalPath := filepath.Join(d.cfg.IdeaDir, "FINAL.md")
		// Validate even if FINAL.md already exists: a stale scaffold from a prior
		// failed/partial draft must be re-drafted, never accepted (D7, AF1). The
		// idea status is committed to "final" ONLY after the content validates, so a
		// failed drafter can never strand the idea at status=final with scaffold
		// content.
		if finalScaffoldReason(finalPath) != "" {
			if err := d.cfg.Consensus.DraftFinal(ctx); err != nil {
				return ActionEscalated, c, fmt.Errorf("draft FINAL.md: %w", err)
			}
			if reason := finalScaffoldReason(finalPath); reason != "" {
				return ActionEscalated, c, fmt.Errorf("FINAL.md is not acceptable after drafting: %s", reason)
			}
		}
		if err := setIdeaStatus(d.cfg.IdeaDir, "final"); err != nil {
			return ActionEscalated, c, fmt.Errorf("commit idea status final: %w", err)
		}
		c.Phase = PhaseFinal
		c.IdeaStatus = "final"
		c.UpdatedAt = nowRFC3339()
		if err := d.commitCursor(c, ActionFinalized, PhaseConsensus); err != nil {
			return ActionEscalated, c, err
		}
		return ActionFinalized, c, nil

	case consensus.TriagePartial:
		if err := d.cfg.Consensus.RequestSignoffs(ctx, summary.Missing); err != nil {
			return ActionEscalated, c, fmt.Errorf("request signoffs: %w", err)
		}
		// RequestSignoffs is synchronous (it invokes the missing signers and waits).
		// If the same participants are still missing, the agents ran but did not
		// sign — re-requesting would loop, so escalate for human help instead.
		after, err := d.cfg.Consensus.Status()
		if err != nil {
			return ActionEscalated, c, fmt.Errorf("read consensus status after signoffs: %w", err)
		}
		if after.Triage == consensus.TriagePartial {
			return ActionEscalated, c, fmt.Errorf("signoffs still missing after request: %s", strings.Join(after.Missing, ", "))
		}
		return ActionSignoffsRequested, c, nil

	case consensus.TriageBlocked:
		next := highestRound(d.cfg.IdeaDir) + 1
		// The §4.0 per-track cross-review ceiling binds HERE too, not only on the
		// initially scheduled budget. Without this, a BLOCKed consensus reopens rounds
		// under MaxRounds alone and walks straight past the printed cap.
		if d.cfg.HardCrossReviewCap > 0 && next > 1+d.cfg.HardCrossReviewCap {
			return ActionEscalated, c, fmt.Errorf(
				"consensus still blocked after %d cross-review round(s) after round 1; round %d would exceed the §4.0 cap of %d; escalating for human review",
				next-2, next, d.cfg.HardCrossReviewCap)
		}
		if next > 1+d.cfg.MaxRounds {
			return ActionEscalated, c, fmt.Errorf("consensus still blocked at round %d (MaxRounds=%d); escalating", next, d.cfg.MaxRounds)
		}
		// Open the re-deliberation round FIRST (beyond the cross_review budget,
		// bounded by MaxRounds). If RunRound fails, the BLOCK state is preserved —
		// consensus.md is still present, so the next tick re-enters this branch and
		// retries, rather than forgetting the BLOCK and re-drafting the old round
		// (AF5).
		if err := d.runner.RunRound(ctx, next); err != nil {
			return ActionEscalated, c, fmt.Errorf("run reopened %s: %w", roundLabel(next), err)
		}
		// Commit the reopen: consensus.Reopen renames consensus.md to an aborted
		// artifact and sets idea status back to the latest (now `next`) round.
		if err := d.cfg.Consensus.Reopen(ctx, "consensus blocked; opening another cross-review round seeded by the BLOCK counter-proposal"); err != nil {
			return ActionEscalated, c, fmt.Errorf("reopen consensus: %w", err)
		}
		// Invalidate any remaining stale consensus.md/FINAL.md so Rebuild classifies
		// the idea as back in the round phase (consensus D6, hermes/agy point). A
		// failure here must halt rather than leave a stale FINAL.md that Rebuild
		// would misread as PhaseFinal (AF4).
		if err := d.invalidateStale(); err != nil {
			return ActionEscalated, c, fmt.Errorf("invalidate stale consensus/FINAL after reopen: %w", err)
		}
		c.CurrentRound = next
		c.RoundsRun = next
		c.Phase = PhaseRound
		c.UpdatedAt = nowRFC3339()
		if err := d.commitCursor(c, ActionReopened, PhaseConsensus); err != nil {
			return ActionEscalated, c, err
		}
		return ActionReopened, c, nil

	default: // TriageMalformed or unknown
		return ActionEscalated, c, fmt.Errorf("consensus triage=%q; cannot auto-advance, escalating", summary.Triage)
	}
}

// invalidateStale renames a stale consensus.md / FINAL.md to *.bak after a BLOCK
// reopen so Rebuild cannot misclassify the reopened round as already at consensus
// or final (consensus D6). It removes a pre-existing *.bak first so the rename
// works across BLOCK cycles and on Windows (where os.Rename fails if the
// destination exists), and returns an error so the caller can halt rather than
// silently leave a stale FINAL.md behind (AF4).
func (d *Driver) invalidateStale() error {
	for _, name := range []string{"consensus.md", "FINAL.md"} {
		path := filepath.Join(d.cfg.IdeaDir, name)
		if !fileExists(path) {
			continue
		}
		bak := path + ".bak"
		if fileExists(bak) {
			if err := os.Remove(bak); err != nil {
				return err
			}
		}
		if err := os.Rename(path, bak); err != nil {
			return err
		}
	}
	return nil
}

// finalScaffoldReason returns a non-empty reason when FINAL.md is still a scaffold
// or otherwise unacceptable (consensus D7), or "" when it is acceptable content.
// finalScaffoldReason reports why a FINAL.md is not a closable final, or "" when it is.
// Everything it checks now lives in protocol.ValidateFinal, shared with manual finalization.
func finalScaffoldReason(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return "FINAL.md is missing"
	}
	if len(data) <= 250 {
		return "FINAL.md is shorter than 250 bytes (scaffold)"
	}
	// One gate, shared with manual finalization (review round 1, @codex-1 MAJOR): the two used to
	// check different things, so an artifact rejected by one path closed an idea through the other.
	if reason := protocol.ValidateFinal(string(data), filepath.Base(filepath.Dir(path))); reason != "" {
		return "FINAL.md " + reason
	}
	return ""
}
