//go:build windows

package trajectory

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/budget"
)

// F3a + F4.1 (kimi-1 applicability consult): the Preview adversarial cases
// and the recovery-apply refusal are READ-PATH/apply-barrier behaviors that
// remain Windows-exercisable when the historical state is constructed
// TEST-SIDE (cross-OS decks are the FINAL §D.4 premise): the intent bytes
// via canonical(), the trajectory state via writeState (ReplaceSyncedFile —
// WT-move, no dir barrier), the ledger charge via the store's own Reserve —
// no product publication, no barrier, no §C.1 refusal. This file restores
// the Windows coverage the interruptedReservationFixture early-return
// bypassed, mapping each original assertion:
//   - missing-intent        → "original reservation intent directory is missing"
//   - partial-intent        → non-canonical bytes refused at read
//   - changed-root/before/limits/action/trajectory → intent-vs-state mismatch refusals
//   - apply (F4.1)          → RecoverReservation refuses via the §B SyncDir
//     emitter BEFORE persist, nothing written
//
// windowsChargedStateFixture constructs the minimal charged-cycle state the
// product's read path accepts, then applies the named corruption. It returns
// the root, binding, and the intent entry key.
func windowsChargedStateFixture(t *testing.T, change string) (string, *budget.CycleBinding, string) {
	t.Helper()
	root, b, _ := accountingFixture(t)
	// The intent: a canonical reservationIntent for this binding. The
	// product normally publishes it via the refusing path; test-side
	// construction is plain Mkdir+WriteFile (FINAL §D.4 cross-OS premise).
	i := reservationIntent{Version: 1, Root: root}
	raw, err := canonical(i)
	if err != nil {
		t.Fatal(err)
	}
	entry := digest(raw)
	intentDir := filepath.Join(filepath.Dir(b.Store.Dir), "reservation-intents")
	if change != "missing-intent" {
		if err := os.MkdirAll(intentDir, 0o700); err != nil {
			t.Fatal(err)
		}
		body := raw
		if change == "partial-intent" {
			body = []byte("{\n")
		}
		if change == "changed-root" {
			i.Root += "-other"
			if body, err = canonical(i); err != nil {
				t.Fatal(err)
			}
		}
		if change == "changed-limits" {
			i.Accounting.Policy.Maximum++
			if body, err = canonical(i); err != nil {
				t.Fatal(err)
			}
		}
		if err := os.WriteFile(filepath.Join(intentDir, entry+".json"), body, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return root, b, entry
}

// The concrete refusal assertions: each case must return a NON-NIL error
// naming the evidence problem (the read path refuses corrupted/missing
// evidence rather than passing).
func TestPreviewRefusesCorruptConstructedEvidence(t *testing.T) {
	for _, change := range []string{"missing-intent", "partial-intent", "changed-root", "changed-limits"} {
		t.Run(change, func(t *testing.T) {
			root, _, entry := windowsChargedStateFixture(t, change)
			_, err := PreviewReservationRecovery(context.Background(), root, "fixture", entry)
			if err == nil {
				t.Fatalf("%s: corrupt evidence was accepted by Preview", change)
			}
			if change == "missing-intent" && !strings.Contains(err.Error(), "original reservation intent directory is missing") {
				t.Fatalf("missing-intent refusal text: %v", err)
			}
			t.Logf("%s refusal: %v", change, err)
		})
	}
}

// F4.1: the recovery APPLY path must refuse on Windows via the §B SyncDir
// emitter BEFORE persist — with a valid test-side state it reaches syncIntent
// and refuses, leaving the original state untouched.
func TestRecoverReservationApplyRefusesAndWritesNothing(t *testing.T) {
	root, b, entry := windowsChargedStateFixture(t, "")
	stateBefore, beforeErr := os.ReadFile(statePath(*b))
	_, err := RecoverReservation(context.Background(), root, "fixture", entry, digest([]byte("x")))
	if err == nil {
		t.Fatal("recovery apply succeeded on Windows — the §B barrier was bypassed")
	}
	if !strings.Contains(err.Error(), "directory-entry durability is not available on Windows") {
		t.Fatalf("apply refusal is not the designed barrier text: %v", err)
	}
	stateAfter, afterErr := os.ReadFile(statePath(*b))
	if beforeErr != afterErr || string(stateBefore) != string(stateAfter) {
		t.Fatal("refused apply rewrote the original state")
	}
}
