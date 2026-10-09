package driver

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"parley-deck-cli/internal/consensus"
)

func TestDropoutSingleReviewerChangesOnlyCountGate(t *testing.T) {
	for _, kind := range []string{"qualified", "no-proof", "no-reviewer", "reservations", "goal-fail", "blocked", "strict-nit"} {
		t.Run(kind, func(t *testing.T) {
			ideaDir, runDir, parts := setupReviewPhase(t, "auto_implement: true\n")
			os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0600)
			rs := closeReady(consensus.TriageReady, 0, 1)
			rs.SingleReviewerAfterDropout = kind != "no-proof"
			fi := &fakeImpl{roundComplete: true, review: rs, goalFail: kind == "goal-fail"}
			if kind == "no-reviewer" {
				fi.review.ReviewerCount = 0
			}
			if kind == "reservations" {
				fi.review.Summary.Triage = consensus.TriageReserved
			}
			if kind == "blocked" {
				fi.review.Blocked = true
			}
			d := newImplDriver(ideaDir, runDir, parts, true, fi)
			if kind == "strict-nit" {
				fi.review.StrictGateClean, fi.review.ClosingReviewRound = true, 1
				os.WriteFile(filepath.Join(ideaDir, "review", "round-01", "codex.md"), []byte("---\nagent: codex\nidea: demo\nreview-round: 1\n---\n\n## Refutation attempts\nchecked\n\n## Findings\n### [NIT] wrong comment\nObjective factual error.\n"), 0600)
				d = newStrictDriver(ideaDir, runDir, parts, 3, fi)
			}
			action, _, err := d.Advance(context.Background())
			if kind == "qualified" {
				if err != nil || action != ActionComplete || !contains(fi.calls, "goal-check") {
					t.Fatalf("qualified: %s %v %v", action, err, fi.calls)
				}
			} else if action == ActionComplete || contains(fi.calls, "complete") {
				t.Fatalf("%s bypassed another gate: %s %v", kind, action, err)
			}
		})
	}
}
