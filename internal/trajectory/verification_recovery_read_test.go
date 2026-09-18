package trajectory

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/telemetry"
)

// Read-only recovery reader and the attended preview/apply authority. Every
// fixture is built through the public APIs (real refused launches, real
// recovery admission); the reader must never write, reserve or publish, and
// missing or drifted authority must refuse rather than recreate a launch.

// backdateRetainedRecoveryIntoRefusalWindow rewrites the admitted recovery so
// its At falls strictly between the reservation and the refusal completion.
func backdateRetainedRecoveryIntoRefusalWindow(t *testing.T, path, root, refused string) {
	t.Helper()
	var launch verificationLaunch
	if err := json.Unmarshal(snapshotRead(t, filepath.Join(path, "launch.json")), &launch); err != nil {
		t.Fatal(err)
	}
	var refusedTerminal telemetry.Record
	if err := json.Unmarshal(snapshotRead(t, filepath.Join(root, ".parley-runtime", "invocations", refused, "terminal.json")), &refusedTerminal); err != nil {
		t.Fatal(err)
	}
	if refusedTerminal.CompletedAt == nil || !refusedTerminal.CompletedAt.After(launch.At) {
		t.Fatalf("fixture refusal did not complete after its reservation: %+v", refusedTerminal)
	}
	var admitted verificationRecovery
	if err := json.Unmarshal(snapshotRead(t, filepath.Join(path, "recovery.json")), &admitted); err != nil {
		t.Fatal(err)
	}
	admitted.At = launch.At.Add(refusedTerminal.CompletedAt.Sub(launch.At) / 2)
	data, err := canonical(admitted)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(path, "recovery.json"), data, 0600); err != nil {
		t.Fatal(err)
	}
}

// TestReadCapturedVerificationRecoveryReturnsValidatedIdentity: after a genuine
// recovery the read-only reader returns the exact bound identity — the intended
// invocation, the refused invocation, the bound evidence digests, the prior
// launch digest and the recovery artifact's own digest — with the journal and
// accounting byte-identical and no bound invocation directory allocated.
func TestReadCapturedVerificationRecoveryReturnsValidatedIdentity(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX captured verification")
	}
	ctx := context.Background()
	ticket, path := journalFixture(t, capturedCriteria(), dirtyCapturedChild(t))
	b, err := budget.LoadCycleBinding(ctx, ticket.Root, "fixture", budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	ledgerBefore := snapshotRead(t, filepath.Join(b.Store.Dir, "ledger.json"))
	refused := refusedVerifierLaunchFixture(t, ticket, nil, false)
	launchBefore := snapshotRead(t, filepath.Join(path, "launch.json"))
	if err = RecoverCapturedVerificationLaunch(ctx, ticket, refused, "recovered-verifier-invocation"); err != nil {
		t.Fatal(err)
	}
	recoveryBytes := snapshotRead(t, filepath.Join(path, "recovery.json"))
	var recovery verificationRecovery
	if err = json.Unmarshal(recoveryBytes, &recovery); err != nil {
		t.Fatal(err)
	}

	identity, err := ReadCapturedVerificationRecovery(ctx, ticket)
	if err != nil {
		t.Fatalf("read-only recovery resolution refused a valid retained recovery: %v", err)
	}
	if identity.InvocationID != "recovered-verifier-invocation" || identity.RefusedInvocationID != refused ||
		identity.RefusedRequestedSHA256 != recovery.RefusedRequestedSHA256 ||
		identity.RefusedTerminalSHA256 != recovery.RefusedTerminalSHA256 ||
		identity.PriorLaunchSHA256 != recovery.PriorLaunchSHA256 ||
		identity.RecoverySHA256 != digest(recoveryBytes) || !identity.At.Equal(recovery.At) {
		t.Fatalf("reader returned an identity diverging from the retained recovery: %+v vs %+v", identity, recovery)
	}
	// The read wrote nothing: journal inventory and bytes, accounting and the
	// (not yet allocated) bound invocation directory are unchanged.
	if names := recoveryJournalNames(t, path); len(names) != 3 || !names["request.json"] || !names["launch.json"] || !names["recovery.json"] {
		t.Fatalf("read changed the retained journal inventory: %v", names)
	}
	if !bytes.Equal(launchBefore, snapshotRead(t, filepath.Join(path, "launch.json"))) {
		t.Fatal("read rewrote the original launch reservation")
	}
	if !bytes.Equal(recoveryBytes, snapshotRead(t, filepath.Join(path, "recovery.json"))) {
		t.Fatal("read rewrote the immutable recovery")
	}
	if !bytes.Equal(ledgerBefore, snapshotRead(t, filepath.Join(b.Store.Dir, "ledger.json"))) {
		t.Fatal("read changed the original fixup accounting")
	}
	if _, statErr := os.Stat(filepath.Join(ticket.Root, ".parley-runtime", "invocations", "recovered-verifier-invocation")); !os.IsNotExist(statErr) {
		t.Fatal("read allocated the bound invocation directory")
	}
}

// TestReadCapturedVerificationRecoveryRefusesMissingAuthority: a never-reserved
// ticket, a refused-but-never-recovered launch, and a recovery whose authority
// files were BOTH deleted afterwards each refuse — and nothing is written. The
// deleted-both case is the deterministic form of the former probe's race: the
// read-only reader can never recreate a launch reservation there.
func TestReadCapturedVerificationRecoveryRefusesMissingAuthority(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX captured verification")
	}
	ctx := context.Background()
	const message = "requires a retained launch reservation and immutable recovery record"

	t.Run("never-reserved", func(t *testing.T) {
		ticket, path := journalFixture(t, capturedCriteria(), dirtyCapturedChild(t))
		if _, err := ReadCapturedVerificationRecovery(ctx, ticket); err == nil || !strings.Contains(err.Error(), message) {
			t.Fatalf("never-reserved ticket was admitted: %v", err)
		}
		if names := recoveryJournalNames(t, path); len(names) != 1 || !names["request.json"] {
			t.Fatalf("refused read wrote artifacts: %v", names)
		}
	})

	t.Run("never-recovered", func(t *testing.T) {
		ticket, path := journalFixture(t, capturedCriteria(), dirtyCapturedChild(t))
		refusedVerifierLaunchFixture(t, ticket, nil, false)
		launchBefore := snapshotRead(t, filepath.Join(path, "launch.json"))
		if _, err := ReadCapturedVerificationRecovery(ctx, ticket); err == nil || !strings.Contains(err.Error(), message) {
			t.Fatalf("unrecovered refused launch was admitted: %v", err)
		}
		if !bytes.Equal(launchBefore, snapshotRead(t, filepath.Join(path, "launch.json"))) {
			t.Fatal("refused read rewrote the launch reservation")
		}
		if _, statErr := os.Stat(filepath.Join(path, "recovery.json")); !os.IsNotExist(statErr) {
			t.Fatal("refused read wrote a recovery artifact")
		}
	})

	t.Run("authority-deleted-after-recovery", func(t *testing.T) {
		ticket, path := journalFixture(t, capturedCriteria(), dirtyCapturedChild(t))
		refused := refusedVerifierLaunchFixture(t, ticket, nil, false)
		if err := RecoverCapturedVerificationLaunch(ctx, ticket, refused, "recovered-verifier-invocation"); err != nil {
			t.Fatal(err)
		}
		for _, name := range []string{"launch.json", "recovery.json"} {
			if err := os.Remove(filepath.Join(path, name)); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := ReadCapturedVerificationRecovery(ctx, ticket); err == nil || !strings.Contains(err.Error(), message) {
			t.Fatalf("deleted authority was admitted: %v", err)
		}
		// The reader must not recreate either authority file: presence is never
		// authority, and a read is never a reservation or publication.
		for _, name := range []string{"launch.json", "recovery.json"} {
			if _, statErr := os.Stat(filepath.Join(path, name)); !os.IsNotExist(statErr) {
				t.Fatalf("refused read recreated %s", name)
			}
		}
		if names := recoveryJournalNames(t, path); len(names) != 1 || !names["request.json"] {
			t.Fatalf("refused read changed the retained journal inventory: %v", names)
		}
	})
}

// TestReadCapturedVerificationRecoveryRefusesDrift: canonical post-recovery
// drift of the retained lineage fails the read with the lineage error and
// writes nothing.
func TestReadCapturedVerificationRecoveryRefusesDrift(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX captured verification")
	}
	ctx := context.Background()
	for _, tc := range []struct {
		name    string
		mutate  func(t *testing.T, ticket VerificationTicket, path, refused string)
		message string
	}{
		{name: "refused-evidence-rewritten", mutate: func(t *testing.T, ticket VerificationTicket, _, refused string) {
			rewriteRefusedEvidenceCanonical(t, ticket.Root, refused)
		}, message: "no longer matches its retained refused-launch evidence"},
		{name: "refused-terminal-deleted", mutate: func(t *testing.T, ticket VerificationTicket, _, refused string) {
			if err := os.Remove(filepath.Join(ticket.Root, ".parley-runtime", "invocations", refused, "terminal.json")); err != nil {
				t.Fatal(err)
			}
		}, message: "refused verifier launch lacks its exact pre-start budget-refusal terminal"},
		{name: "forged-prior-launch", mutate: func(t *testing.T, _ VerificationTicket, path, _ string) {
			var admitted verificationRecovery
			if err := json.Unmarshal(snapshotRead(t, filepath.Join(path, "recovery.json")), &admitted); err != nil {
				t.Fatal(err)
			}
			admitted.PriorLaunchSHA256 = digest([]byte("forged"))
			data, err := canonical(admitted)
			if err != nil {
				t.Fatal(err)
			}
			if err = os.WriteFile(filepath.Join(path, "recovery.json"), data, 0600); err != nil {
				t.Fatal(err)
			}
		}, message: "contradicts its original launch reservation"},
		{name: "recovery-predates-refusal-completion", mutate: func(t *testing.T, ticket VerificationTicket, path, refused string) {
			backdateRetainedRecoveryIntoRefusalWindow(t, path, ticket.Root, refused)
		}, message: "predates the retained refusal completion"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ticket, path := journalFixture(t, capturedCriteria(), dirtyCapturedChild(t))
			refused := refusedVerifierLaunchFixture(t, ticket, nil, false)
			if err := RecoverCapturedVerificationLaunch(ctx, ticket, refused, "recovered-verifier-invocation"); err != nil {
				t.Fatal(err)
			}
			tc.mutate(t, ticket, path, refused)
			if _, err := ReadCapturedVerificationRecovery(ctx, ticket); err == nil || !strings.Contains(err.Error(), tc.message) {
				t.Fatalf("%s was admitted: %v", tc.name, err)
			}
			if _, statErr := os.Stat(filepath.Join(ticket.Root, ".parley-runtime", "invocations", "recovered-verifier-invocation")); !os.IsNotExist(statErr) {
				t.Fatalf("%s allocated the bound invocation directory", tc.name)
			}
		})
	}
}

// TestCapturedRecoveryPreviewApplyExactRecheck: the attended preview binds the
// exact refused launch and the selected replacement; the apply rechecks the
// inspected digest before the single immutable write; replay is exact and
// write-free; conflicts refuse against the retained artifact.
func TestCapturedRecoveryPreviewApplyExactRecheck(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX captured verification")
	}
	ctx := context.Background()
	ticket, path := journalFixture(t, capturedCriteria(), dirtyCapturedChild(t))
	b, err := budget.LoadCycleBinding(ctx, ticket.Root, "fixture", budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	ledgerBefore := snapshotRead(t, filepath.Join(b.Store.Dir, "ledger.json"))
	refused := refusedVerifierLaunchFixture(t, ticket, nil, false)
	launchBefore := snapshotRead(t, filepath.Join(path, "launch.json"))
	ticketSHA, err := ticket.SHA256()
	if err != nil {
		t.Fatal(err)
	}

	p, err := PreviewCapturedVerificationRecovery(ctx, ticket, "recovered-verifier-invocation")
	if err != nil {
		t.Fatal(err)
	}
	if p.Version != 1 || p.Root != ticket.Root || p.Idea != "fixture" || p.RunID != ticket.RunID ||
		p.TicketSHA256 != ticketSHA || p.RefusedInvocationID != refused ||
		!validHash(p.RefusedRequestedSHA256) || !validHash(p.RefusedTerminalSHA256) ||
		p.RefusedCompletedAt.IsZero() || p.InvocationID != "recovered-verifier-invocation" {
		t.Fatalf("preview does not bind the exact refused launch and replacement: %+v", p)
	}
	again, err := PreviewCapturedVerificationRecovery(ctx, ticket, "recovered-verifier-invocation")
	if err != nil || again.SHA256() != p.SHA256() {
		t.Fatal("preview is not exact for the same replacement")
	}
	other, err := PreviewCapturedVerificationRecovery(ctx, ticket, "other-verifier-invocation")
	if err != nil || other.SHA256() == p.SHA256() {
		t.Fatal("preview did not bind the selected replacement identity")
	}
	if names := recoveryJournalNames(t, path); len(names) != 2 {
		t.Fatalf("preview wrote artifacts: %v", names)
	}

	if _, err = ApplyCapturedVerificationRecovery(ctx, ticket, "recovered-verifier-invocation", digest([]byte("stale"))); err == nil ||
		!strings.Contains(err.Error(), "changed since preview") {
		t.Fatalf("stale preview was applied: %v", err)
	}
	if _, err = ApplyCapturedVerificationRecovery(ctx, ticket, "other-verifier-invocation", p.SHA256()); err == nil {
		t.Fatal("a different replacement applied with another preview's digest")
	}
	if _, statErr := os.Stat(filepath.Join(path, "recovery.json")); !os.IsNotExist(statErr) {
		t.Fatal("refused apply wrote a recovery artifact")
	}

	applied, err := ApplyCapturedVerificationRecovery(ctx, ticket, "recovered-verifier-invocation", p.SHA256())
	if err != nil || applied.SHA256() != p.SHA256() {
		t.Fatalf("exact apply did not reproduce the inspected preview: %+v %v", applied, err)
	}
	recoveryBytes := snapshotRead(t, filepath.Join(path, "recovery.json"))
	var recovery verificationRecovery
	if err = json.Unmarshal(recoveryBytes, &recovery); err != nil {
		t.Fatal(err)
	}
	if recovery.InvocationID != "recovered-verifier-invocation" || recovery.RefusedInvocationID != refused ||
		recovery.RefusedRequestedSHA256 != p.RefusedRequestedSHA256 || recovery.RefusedTerminalSHA256 != p.RefusedTerminalSHA256 {
		t.Fatalf("applied recovery does not bind the inspected facts: %+v", recovery)
	}
	if !bytes.Equal(launchBefore, snapshotRead(t, filepath.Join(path, "launch.json"))) {
		t.Fatal("apply rewrote the original launch reservation")
	}
	if !bytes.Equal(ledgerBefore, snapshotRead(t, filepath.Join(b.Store.Dir, "ledger.json"))) {
		t.Fatal("apply changed the original fixup accounting")
	}

	replayed, err := ApplyCapturedVerificationRecovery(ctx, ticket, "recovered-verifier-invocation", p.SHA256())
	if err != nil || replayed.SHA256() != p.SHA256() {
		t.Fatalf("exact replay changed the retained recovery: %+v %v", replayed, err)
	}
	if !bytes.Equal(recoveryBytes, snapshotRead(t, filepath.Join(path, "recovery.json"))) {
		t.Fatal("exact replay rewrote the retained recovery")
	}
	if _, err = PreviewCapturedVerificationRecovery(ctx, ticket, "recovered-verifier-invocation"); err != nil {
		t.Fatalf("bound preview replay failed after recovery: %v", err)
	}
	if _, err = PreviewCapturedVerificationRecovery(ctx, ticket, "conflicting-verifier-invocation"); err == nil ||
		!strings.Contains(err.Error(), "contradicts the retained immutable recovery") {
		t.Fatalf("conflicting preview was admitted after recovery: %v", err)
	}
	if _, err = ApplyCapturedVerificationRecovery(ctx, ticket, "conflicting-verifier-invocation", other.SHA256()); err == nil ||
		!strings.Contains(err.Error(), "contradicts the retained immutable recovery") {
		t.Fatalf("conflicting apply was admitted after recovery: %v", err)
	}
}

// TestCapturedRecoveryPreviewRefusesUnavailableAuthority: preview and apply
// refuse a never-reserved ticket and an ordinary (never-refused) reservation
// without writing anything.
func TestCapturedRecoveryPreviewRefusesUnavailableAuthority(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX captured verification")
	}
	ctx := context.Background()

	t.Run("never-reserved", func(t *testing.T) {
		ticket, path := journalFixture(t, capturedCriteria(), dirtyCapturedChild(t))
		if _, err := PreviewCapturedVerificationRecovery(ctx, ticket, "recovered-verifier-invocation"); err == nil ||
			!strings.Contains(err.Error(), "requires the original launch reservation") {
			t.Fatalf("never-reserved ticket previewed: %v", err)
		}
		if names := recoveryJournalNames(t, path); len(names) != 1 {
			t.Fatalf("refused preview wrote artifacts: %v", names)
		}
	})

	t.Run("ordinary-reservation-not-a-refusal", func(t *testing.T) {
		ticket, path := journalFixture(t, capturedCriteria(), dirtyCapturedChild(t))
		reserveJournal(t, ticket)
		if _, err := PreviewCapturedVerificationRecovery(ctx, ticket, "recovered-verifier-invocation"); err == nil ||
			!strings.Contains(err.Error(), "refused verifier launch request is unavailable or changed") {
			t.Fatalf("ordinary reservation previewed as a recovery: %v", err)
		}
		if names := recoveryJournalNames(t, path); len(names) != 2 {
			t.Fatalf("refused preview wrote artifacts: %v", names)
		}
	})

	t.Run("replacement-identities-rejected", func(t *testing.T) {
		ticket, _ := journalFixture(t, capturedCriteria(), dirtyCapturedChild(t))
		refusedVerifierLaunchFixture(t, ticket, nil, false)
		for _, replacement := range []string{"", ticket.Request.InvocationID} {
			if _, err := PreviewCapturedVerificationRecovery(ctx, ticket, replacement); err == nil {
				t.Fatalf("invalid replacement %q previewed", replacement)
			}
		}
	})
}
