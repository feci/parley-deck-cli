//go:build windows

package trajectory

import (
	"bytes"
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
// windowsChargedStateFixture constructs a VALID charged-cycle state entirely
// test-side (kimi-1 F3a recipe, FINAL §D.4 cross-OS premise): the ledger
// charge via the store's own Reserve (writeSynced — no dir barrier), the
// pre-charge trajectory state via the product's own read (the fixture's
// Activate already wrote it — no rewrite needed), and the intent built from
// the SAME fields the product's PrepareCycleReservation assembles, written
// with plain Mkdir/WriteFile under the accounting EntryKey filename. No
// product publication, no §C.1 refusal. The corruption switch then mutates
// one guard's input. A CLEAN CONTROL (no mutation) must pass Preview's
// identity checks before any mutated case runs — otherwise the fixture
// itself is malformed and every "refusal" would be a false pass.
func windowsChargedStateFixture(t *testing.T, change string) (string, *budget.CycleBinding, string) {
	t.Helper()
	root, b, _ := accountingFixture(t)
	ctx := context.Background()
	// Ledger charge (the store path — durable via ReplaceSyncedFile).
	zero := int64(0)
	if _, err := b.Store.Reserve(ctx, budget.Request{ID: "fixture-charge", Kind: budget.Fixup, ReserveMicros: &zero}, budget.Limits{}); err != nil {
		t.Fatal(err)
	}
	snap, err := b.Store.Inspect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var entry string
	for k := range snap.Entries {
		entry = k
	}
	if entry == "" {
		t.Fatal("no ledger charge entry")
	}
	// The pre-charge state IS the fixture's trajectory.json (Activate wrote
	// it; the charge appends only at the product's own later step).
	s, raw, err := readState(statePath(*b))
	if err != nil {
		t.Fatal(err)
	}
	accounting := budget.CycleReservationIntent{
		Version: 1, Policy: b.Policy, StartedAt: snap.StartedAt,
		EntryKey: entry, ReserveMicros: &zero,
	}
	// canonicalRoot: the product compares i.Root against the canonicalized
	// (Abs + EvalSymlinks) worktree — the raw t.TempDir string differs
	// hosted (58f926a's F4.1 finding: "recovery worktree differs").
	croot, err := canonicalRoot(root)
	if err != nil {
		t.Fatal(err)
	}
	i := reservationIntent{Version: 1, Root: croot, PreparedAt: snap.StartedAt.UTC(), Accounting: accounting, Before: s, BeforeSHA256: digest(raw)}
	ibody, err := canonical(i)
	if err != nil {
		t.Fatal(err)
	}
	intentDir := filepath.Join(filepath.Dir(b.Store.Dir), "reservation-intents")
	if change != "missing-intent" {
		if err := os.MkdirAll(intentDir, 0o700); err != nil {
			t.Fatal(err)
		}
		body := ibody
		switch change {
		case "partial-intent":
			body = []byte("{\n")
		case "changed-root":
			i.Root = croot + "-other"
			if body, err = canonical(i); err != nil {
				t.Fatal(err)
			}
		case "changed-before":
			i.Before.Policy.Baseline.Tree.SHA256 = digest([]byte("other"))
			if body, err = canonical(i); err != nil {
				t.Fatal(err)
			}
		case "changed-limits":
			i.Accounting.Policy.Maximum++
			if body, err = canonical(i); err != nil {
				t.Fatal(err)
			}
		case "changed-action":
			i.Accounting.Action = &budget.ActionIdentity{Version: 1, EntrySHA256: digest([]byte("other"))}
			if body, err = canonical(i); err != nil {
				t.Fatal(err)
			}
		case "changed-trajectory":
			i.Before.Policy.Implementer = "other"
			i.BeforeSHA256 = digest([]byte("changed"))
			if body, err = canonical(i); err != nil {
				t.Fatal(err)
			}
		case "symlink-intent":
			target := filepath.Join(t.TempDir(), "retained.json")
			if err := os.WriteFile(target, ibody, 0o600); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(target, filepath.Join(intentDir, entry+".json")); err != nil {
				t.Fatal(err)
			}
			return root, b, entry
		}
		if err := os.WriteFile(filepath.Join(intentDir, entry+".json"), body, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	// extra-charge and missing-archive mutate the STORE/STATE sides.
	switch change {
	case "extra-charge":
		if _, err := b.Store.Reserve(ctx, budget.Request{ID: "unexpected", Kind: budget.Fixup, ReserveMicros: &zero}, budget.Limits{}); err != nil {
			t.Fatal(err)
		}
	case "missing-archive":
		if err := os.RemoveAll(snapshotDirectory(*b)); err != nil {
			t.Fatal(err)
		}
	}
	return root, b, entry
}

// The clean control FIRST (the discipline the correction demands): on the
// unmutated construction, Preview must NOT refuse with an evidence error —
// it may report a preview status, but any "incomplete or changed" text here
// means the FIXTURE is malformed and every mutated case below would be a
// false pass. Then each of the ten original adversarial cases runs with a
// NON-NIL refusal (guard-specific evidence printed via t.Logf for hosted
// pinning; missing-intent additionally pinned to its exact designed text).
func TestPreviewAdversarialRefusalsOnConstructedState(t *testing.T) {
	t.Run("clean-control", func(t *testing.T) {
		root, _, entry := windowsChargedStateFixture(t, "")
		_, err := PreviewReservationRecovery(context.Background(), root, "fixture", entry)
		if err != nil && strings.Contains(err.Error(), "incomplete or changed") {
			t.Fatalf("fixture is malformed (clean control refused): %v", err)
		}
		t.Logf("clean control outcome: err=%v", err)
	})
	for _, change := range []string{"missing-intent", "partial-intent", "changed-root", "changed-before", "changed-limits", "changed-action", "missing-archive", "extra-charge", "changed-trajectory", "symlink-intent"} {
		t.Run(change, func(t *testing.T) {
			root, b, entry := windowsChargedStateFixture(t, change)
			// V4 (kimi gap 4): the read-path refusal must leave BOTH the
			// trajectory state AND the ledger byte-identical.
			stateBefore := snapshotRead(t, statePath(*b))
			ledgerBefore := snapshotRead(t, filepath.Join(b.Store.Dir, "ledger.json"))
			_, err := PreviewReservationRecovery(context.Background(), root, "fixture", entry)
			if err == nil {
				t.Fatalf("%s: corrupt evidence was accepted by Preview", change)
			}
			if change == "missing-intent" && !strings.Contains(err.Error(), "original reservation intent directory is missing") {
				t.Fatalf("missing-intent refusal text: %v", err)
			}
			if got := snapshotRead(t, statePath(*b)); !bytes.Equal(got, stateBefore) {
				t.Fatalf("%s: read-path Preview rewrote trajectory state", change)
			}
			if got := snapshotRead(t, filepath.Join(b.Store.Dir, "ledger.json")); !bytes.Equal(got, ledgerBefore) {
				t.Fatalf("%s: read-path Preview rewrote the ledger", change)
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
	// RecoverReservation requires expected == the exact Preview digest —
	// take the Preview first (the product's own contract), then apply with it.
	preview, perr := PreviewReservationRecovery(context.Background(), root, "fixture", entry)
	if perr != nil {
		t.Fatalf("preview failed on the valid construction: %v", perr)
	}
	stateBefore, beforeErr := os.ReadFile(statePath(*b))
	_, err := RecoverReservation(context.Background(), root, "fixture", entry, preview.SHA256())
	if err == nil {
		t.Fatal("recovery apply succeeded on Windows — the §B barrier was bypassed")
	}
	// The barrier CLASS: either the named durability refusal (the dir part)
	// or the raw Access-denied from the intent file's read-only-handle sync
	// (the D4a file part — hosted 444d930 evidence; both are pre-persist
	// fail-closed; the unnamed raw text is recorded as a §B polish item).
	if !strings.Contains(err.Error(), "directory-entry durability is not available on Windows") &&
		!strings.Contains(err.Error(), "Access is denied") {
		t.Fatalf("apply refusal is not the barrier class: %v", err)
	}
	stateAfter, afterErr := os.ReadFile(statePath(*b))
	if beforeErr != afterErr || string(stateBefore) != string(stateAfter) {
		t.Fatal("refused apply rewrote the original state")
	}
}
