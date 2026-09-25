package consensus

import (
	"os"
	"path/filepath"
	"testing"
)

// T-11 / AC-15 (R32): review exclusion reads the RECORD of who implemented, never the
// instruction of who was designated. Tiers 2 and 3 never enter the consensus-side
// read, so a designation cannot shrink a review round's expected artifact set.
func TestReviewExclusionReadsThePinNotTheDesignation(t *testing.T) {
	idea := t.TempDir()
	write := func(name, body string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(idea, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	participants := []string{"pin-agent", "des-agent", "rev-agent"}

	// Tier-2 designation names des-agent; the IMPLEMENTATION.md pin names pin-agent.
	// Review expectations must exclude the PIN (who implemented), not the designee.
	write("00-prompt.md", "---\nidea: x\nparticipants: [pin-agent, des-agent, rev-agent]\nimplementer: des-agent\n---\n")
	write("IMPLEMENTATION.md", "---\nidea: x\nstatus: implemented\nimplementer: pin-agent\n---\n")
	got := expectedRoundParticipants(idea, participants, true)
	if len(got) != 2 || got[0] != "des-agent" || got[1] != "rev-agent" {
		t.Fatalf("review must exclude the recorded implementer (pin), got %v", got)
	}

	// With no pin, the FINAL drafter fallback reads FINAL.md only — still never the
	// designation.
	if err := os.Remove(filepath.Join(idea, "IMPLEMENTATION.md")); err != nil {
		t.Fatal(err)
	}
	write("FINAL.md", "---\nidea: x\ndrafted-by: rev-agent\n---\n")
	got = expectedRoundParticipants(idea, participants, true)
	if len(got) != 2 || got[0] != "pin-agent" || got[1] != "des-agent" {
		t.Fatalf("review must exclude the FINAL drafter, got %v", got)
	}

	// With no artifacts at all the designation changes nothing: fail closed to the
	// full list, exactly as before the feature existed.
	if err := os.Remove(filepath.Join(idea, "FINAL.md")); err != nil {
		t.Fatal(err)
	}
	got = expectedRoundParticipants(idea, participants, true)
	if len(got) != 3 {
		t.Fatalf("an unresolvable record must expect everyone regardless of the designation, got %v", got)
	}
	if got := ExpectedRoundParticipants(idea, participants, true); len(got) != 3 {
		t.Fatalf("the exported read must be equally designation-blind, got %v", got)
	}
}

// AC-4 parity: ExpectedRoundParticipants is unaffected by a designation in every
// state — design rounds, review rounds, the unresolvable case, the FINAL-drafter case.
func TestExpectedRoundParticipantsUnaffectedByDesignation(t *testing.T) {
	participants := []string{"a", "b", "c"}
	states := []struct {
		name  string
		files map[string]string
	}{
		{"design round", map[string]string{}},
		{"review round, unresolvable", map[string]string{}},
		{"review round, FINAL drafter", map[string]string{"FINAL.md": "---\ndrafted-by: a\n---\n"}},
		{"review round, pin", map[string]string{"IMPLEMENTATION.md": "---\nimplementer: b\n---\n"}},
	}
	for _, tc := range states {
		t.Run(tc.name, func(t *testing.T) {
			review := tc.name != "design round"
			without := t.TempDir()
			for name, body := range tc.files {
				if err := os.WriteFile(filepath.Join(without, name), []byte(body), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			with := t.TempDir()
			for name, body := range tc.files {
				if err := os.WriteFile(filepath.Join(with, name), []byte(body), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			if err := os.WriteFile(filepath.Join(with, "00-prompt.md"), []byte("---\nidea: x\nimplementer: c\n---\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			gotWithout := ExpectedRoundParticipants(without, participants, review)
			gotWith := ExpectedRoundParticipants(with, participants, review)
			if len(gotWithout) != len(gotWith) {
				t.Fatalf("designation changed the expected set: %v vs %v", gotWithout, gotWith)
			}
			for i := range gotWithout {
				if gotWithout[i] != gotWith[i] {
					t.Fatalf("designation changed the expected set: %v vs %v", gotWithout, gotWith)
				}
			}
		})
	}
}
