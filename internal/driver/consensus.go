package driver

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"parley-deck-cli/internal/consensus"
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
		if !fileExists(finalPath) {
			if err := d.cfg.Consensus.DraftFinal(ctx); err != nil {
				return ActionEscalated, c, fmt.Errorf("draft FINAL.md: %w", err)
			}
			if reason := finalScaffoldReason(finalPath); reason != "" {
				return ActionEscalated, c, fmt.Errorf("FINAL.md is not acceptable: %s", reason)
			}
		}
		c.Phase = PhaseFinal
		c.IdeaStatus = "final"
		c.UpdatedAt = nowRFC3339()
		_ = c.Save(d.cursorPath())
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
		if next > 1+d.cfg.MaxRounds {
			return ActionEscalated, c, fmt.Errorf("consensus still blocked at round %d (MaxRounds=%d); escalating", next, d.cfg.MaxRounds)
		}
		if err := d.cfg.Consensus.Reopen(ctx, "consensus blocked; opening another cross-review round seeded by the BLOCK counter-proposal"); err != nil {
			return ActionEscalated, c, fmt.Errorf("reopen consensus: %w", err)
		}
		// Invalidate the stale consensus.md/FINAL.md so Rebuild classifies the idea
		// as back in the round phase, not at a completed consensus (consensus D6,
		// hermes/agy cross-review point).
		d.invalidateStale()
		// Open the re-deliberation round now (beyond the cross_review budget, bounded
		// by MaxRounds) so the next tick re-drafts consensus from the new round.
		if err := d.runner.RunRound(ctx, next); err != nil {
			return ActionEscalated, c, fmt.Errorf("run reopened %s: %w", roundLabel(next), err)
		}
		if err := setIdeaStatus(d.cfg.IdeaDir, roundLabel(next)); err != nil {
			return ActionEscalated, c, fmt.Errorf("set idea status: %w", err)
		}
		c.CurrentRound = next
		c.RoundsRun = next
		c.Phase = PhaseRound
		c.UpdatedAt = nowRFC3339()
		_ = c.Save(d.cursorPath())
		return ActionReopened, c, nil

	default: // TriageMalformed or unknown
		return ActionEscalated, c, fmt.Errorf("consensus triage=%q; cannot auto-advance, escalating", summary.Triage)
	}
}

// invalidateStale renames a stale consensus.md / FINAL.md to *.bak after a BLOCK
// reopen so Rebuild cannot misclassify the reopened round as already at consensus
// or final (consensus D6).
func (d *Driver) invalidateStale() {
	for _, name := range []string{"consensus.md", "FINAL.md"} {
		path := filepath.Join(d.cfg.IdeaDir, name)
		if fileExists(path) {
			_ = os.Rename(path, path+".bak")
		}
	}
}

// finalScaffoldReason returns a non-empty reason when FINAL.md is still a scaffold
// or otherwise unacceptable (consensus D7), or "" when it is acceptable content.
func finalScaffoldReason(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return "FINAL.md is missing"
	}
	if len(data) <= 250 {
		return "FINAL.md is shorter than 250 bytes (scaffold)"
	}
	if status, _ := readFrontmatterField(path, "status"); status != "final" {
		return fmt.Sprintf("FINAL.md frontmatter status=%q, want final", status)
	}
	body := string(data)
	const section = "## Final plan / specification"
	idx := strings.Index(body, section)
	if idx < 0 {
		return "FINAL.md missing '## Final plan / specification' section"
	}
	rest := body[idx+len(section):]
	if next := strings.Index(rest, "\n## "); next >= 0 {
		rest = rest[:next]
	}
	content := 0
	for _, line := range strings.Split(rest, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "<!--") {
			continue
		}
		if strings.Contains(trimmed, "<") && strings.Contains(trimmed, ">") {
			return "FINAL.md '## Final plan / specification' still contains an unexpanded <…> placeholder"
		}
		content++
	}
	if content < 3 {
		return "FINAL.md '## Final plan / specification' has fewer than 3 content lines (scaffold)"
	}
	return ""
}
