package app

// lean-organizer C.5 brief contract tests: <= 8,192 B, byte-identical across two
// runs over an unchanged tree, writes no file (asserted against a read-only deck),
// and computed — never stored.

import (
	"bytes"
	"fmt"
	"os"
	"parley-deck-cli/internal/fsacl"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/protocol"
)

func briefFor(t *testing.T, root string) (int, string) {
	t.Helper()
	var out, errOut bytes.Buffer
	code := runOrganizer([]string{"brief", "--dir", root, "--idea", "wait-idea"}, &out, &errOut)
	if code != 0 {
		t.Fatalf("brief exit %d: %s", code, errOut.String())
	}
	return code, out.String()
}

func TestOrganizerBriefContract(t *testing.T) {
	root, ideaDir := seedWaitIdea(t, []string{"claude-1", "kimi-1"})
	os.WriteFile(filepath.Join(ideaDir, "round-01", "claude-1.md"), []byte(validRoundOne("claude-1")), 0o644)

	_, a := briefFor(t, root)
	_, b := briefFor(t, root)
	if a != b {
		t.Fatalf("brief must be byte-identical across two runs over an unchanged tree")
	}
	if len(a) > organizerBriefByteCap {
		t.Fatalf("brief is %d B, cap %d B", len(a), organizerBriefByteCap)
	}
	for _, want := range []string{"# Organizer brief", "## Protocol context", "## Idea state", "## PhaseDigest", "## Next action"} {
		if !strings.Contains(a, want) {
			t.Errorf("brief must contain %q", want)
		}
	}
	if !strings.Contains(a, "wait-idea") || !strings.Contains(a, "claude-1") {
		t.Errorf("brief must name the idea and participants; got:\n%s", a)
	}
}

func TestOrganizerBriefWritesNoFileReadOnlyDeck(t *testing.T) {
	root, ideaDir := seedWaitIdea(t, []string{"claude-1"})
	os.WriteFile(filepath.Join(ideaDir, "round-01", "claude-1.md"), []byte(validRoundOne("claude-1")), 0o644)

	// Write-detection that separates PRODUCT WRITES from mtime churn:
	// structure (the set of paths) must be identical, and a sentinel file's
	// BYTES must be untouched. ModTime alone is not evidence of a product
	// write (the hosted failure showed mtime noise under the working deny).
	snapshot := func() string {
		var sb strings.Builder
		filepath.Walk(filepath.Join(root, protocol.DeckDir), func(p string, info os.FileInfo, err error) error {
			if err == nil {
				sb.WriteString(p + "|" + info.Name() + "|" + fmt.Sprint(info.IsDir()) + "\n")
			}
			return nil
		})
		return sb.String()
	}
	sentinel := filepath.Join(root, protocol.DeckDir, "ideas", "wait-idea", "sentinel.md")
	if err := os.WriteFile(sentinel, []byte("untouched\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	before := snapshot()
	// Read-only deck: the brief must still compute and must write nothing into it.
	deck := filepath.Join(root, protocol.DeckDir)
	// GENUINE native deny (§D.9): chmod is a no-op for write access on
	// Windows — the historical pass proved nothing there. DenyWrite installs
	// a deny-ACE (GENERIC_WRITE on a directory = add-file/add-subdirectory,
	// deliberately NOT delete-child so cleanup removals still work), and
	// AllowWrite restores by rebuilding the allow set WITHOUT the protected
	// flag (the invocation-13 regression root cause, now fixed in fsacl).
	if err := fsacl.DenyWrite(deck); err != nil {
		t.Fatalf("cannot make deck read-only: %v", err)
	}
	t.Cleanup(func() { _ = fsacl.AllowWrite(deck) })
	code, out := briefFor(t, root)
	if code != 0 {
		t.Fatalf("brief must work against a read-only deck, exit %d", code)
	}
	after := snapshot()
	if before != after {
		// Hosted delta print: the EXACT created/removed paths, so the
		// product-write investigation starts from evidence.
		beforeSet, afterSet := map[string]bool{}, map[string]bool{}
		for _, l := range strings.Split(before, "\n") {
			if parts := strings.SplitN(l, "|", 2); len(parts) == 2 {
				beforeSet[parts[0]] = true
			}
		}
		for _, l := range strings.Split(after, "\n") {
			if parts := strings.SplitN(l, "|", 2); len(parts) == 2 {
				afterSet[parts[0]] = true
				if !beforeSet[parts[0]] {
					t.Errorf("brief CREATED: %s", parts[0])
				}
			}
		}
		for p := range beforeSet {
			if !afterSet[p] {
				t.Errorf("brief REMOVED: %s", p)
			}
		}
		t.Fatalf("brief wrote into the deck tree (structure changed)")
	}
	if data, err := os.ReadFile(sentinel); err != nil || string(data) != "untouched\n" {
		t.Fatalf("brief disturbed deck content: %q %v", data, err)
	}
	if len(out) == 0 {
		t.Fatalf("brief produced no output")
	}
}

// TestOrganizerBriefPhaseNeverZeroForLivePhaseFivePlus (fix-up F5, claude-1 MAJ-5):
// a live Phase 5–8 idea must never resolve to the phase-0 facilitator packet (whose
// body carries no §15 at all). The deck's real vocabulary — IMPLEMENTATION.md status
// `implemented` / `fix-up-cycle-1`, prompt status left at `final` — must map to
// phases 5 and 8 respectively.
func TestOrganizerBriefPhaseNeverZeroForLivePhaseFivePlus(t *testing.T) {
	cases := []struct {
		name         string
		promptStatus string
		implStatus   string
		wantPhase    int
	}{
		{"implementation status maps to phase 5", "implementation", "", 5},
		{"final plus implemented implementation maps to phase 5", "final", "implemented", 5},
		{"fix-up-cycle-1 maps to phase 8", "final", "fix-up-cycle-1", 8},
		{"complete maps to phase 8", "complete", "complete", 8},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root, ideaDir := seedWaitIdea(t, []string{"claude-1", "kimi-1"})
			if err := protocol.InitWorkspace(root); err != nil {
				t.Fatal(err)
			}
			fm := "---\nidea: wait-idea\nauthor: user\ntrack: standard\nparticipants: [claude-1, kimi-1]\nstatus: " + tc.promptStatus + "\n---\n\n## Problem\n"
			if err := os.WriteFile(filepath.Join(ideaDir, "00-prompt.md"), []byte(fm), 0o644); err != nil {
				t.Fatal(err)
			}
			if tc.implStatus != "" {
				impl := "---\nidea: wait-idea\nstatus: " + tc.implStatus + "\nimplementer: claude-1\nhead-commit: deadbeef\n---\n\n## Summary\nDone.\n"
				if err := os.WriteFile(filepath.Join(ideaDir, "IMPLEMENTATION.md"), []byte(impl), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			_, brief := briefFor(t, root)
			marker := fmt.Sprintf("-phase%d-", tc.wantPhase)
			if !strings.Contains(brief, marker) {
				t.Fatalf("brief must render the phase-%d facilitator packet for prompt status %q / implementation status %q; brief:\n%s", tc.wantPhase, tc.promptStatus, tc.implStatus, brief)
			}
			if strings.Contains(brief, "-phase0-") {
				t.Fatalf("a live Phase 5–8 idea must never resolve to the phase-0 packet (no §15); brief:\n%s", brief)
			}
		})
	}
}
