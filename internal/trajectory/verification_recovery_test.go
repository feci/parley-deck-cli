package trajectory

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/telemetry"
)

// refusedVerifierLaunchFixture reproduces the production ordering of a
// pre-start budget refusal through public APIs: telemetry.Begin persists the
// request, ReserveCapturedVerificationLaunch writes the durable launch
// reservation for that exact invocation, and only then the terminal failure is
// retained. No lifecycle beyond requested+terminal is ever fabricated.
func refusedVerifierLaunchFixture(t *testing.T, ticket VerificationTicket, edit func(*telemetry.Metadata, *telemetry.Outcome), started bool) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("POSIX captured verification")
	}
	metadata := telemetry.Metadata{RunID: ticket.RunID, Idea: ticket.Request.Idea, Phase: "trajectory-verification",
		Agent: ticket.Request.Verifier, Adapter: "claude", LaunchMode: "headless"}
	outcome := telemetry.Outcome{Status: "failed", FailureClass: telemetry.String("budget_refused")}
	if edit != nil {
		edit(&metadata, &outcome)
	}
	inv, err := telemetry.Begin(filepath.Join(ticket.Root, ".parley-runtime", "invocations"), metadata)
	if err != nil {
		t.Fatal(err)
	}
	if err = ReserveCapturedVerificationLaunch(context.Background(), ticket, inv.ID); err != nil {
		t.Fatal(err)
	}
	if started {
		if err = inv.Started(os.Getpid()); err != nil {
			t.Fatal(err)
		}
	}
	if err = inv.Finish(outcome); err != nil {
		t.Fatal(err)
	}
	return inv.ID
}

func recoveryJournalNames(t *testing.T, path string) map[string]bool {
	t.Helper()
	entries, err := os.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}
	names := map[string]bool{}
	for _, e := range entries {
		names[e.Name()] = true
	}
	return names
}

func TestRecoveredRefusedVerifierLaunchRestoresHelperCompletion(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX captured verification")
	}
	ctx := context.Background()
	criteria := capturedCriteria()
	ticket, path := journalFixture(t, criteria, dirtyCapturedChild(t))
	b, err := budget.LoadCycleBinding(ctx, ticket.Root, "fixture", budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	ledgerPath := filepath.Join(b.Store.Dir, "ledger.json")
	ledgerBefore := snapshotRead(t, ledgerPath)
	requestBefore := snapshotRead(t, filepath.Join(path, "request.json"))
	refused := refusedVerifierLaunchFixture(t, ticket, nil, false)
	launchBefore := snapshotRead(t, filepath.Join(path, "launch.json"))

	if err = RecoverCapturedVerificationLaunch(ctx, ticket, refused, "recovered-verifier-invocation"); err != nil {
		t.Fatal(err)
	}
	if names := recoveryJournalNames(t, path); len(names) != 3 || !names["request.json"] || !names["launch.json"] || !names["recovery.json"] {
		t.Fatalf("recovery changed the retained journal inventory: %v", names)
	}
	if !bytes.Equal(requestBefore, snapshotRead(t, filepath.Join(path, "request.json"))) ||
		!bytes.Equal(launchBefore, snapshotRead(t, filepath.Join(path, "launch.json"))) {
		t.Fatal("recovery rewrote the original reservation bytes")
	}
	var recovery verificationRecovery
	if err = json.Unmarshal(snapshotRead(t, filepath.Join(path, "recovery.json")), &recovery); err != nil {
		t.Fatal(err)
	}
	wantTicket, err := ticket.SHA256()
	if err != nil {
		t.Fatal(err)
	}
	if recovery.Version != 1 || recovery.TicketSHA256 != wantTicket || recovery.RefusedInvocationID != refused ||
		recovery.InvocationID != "recovered-verifier-invocation" || !validHash(recovery.PriorLaunchSHA256) ||
		!validHash(recovery.RefusedRequestedSHA256) || !validHash(recovery.RefusedTerminalSHA256) || recovery.At.IsZero() {
		t.Fatalf("recovery artifact is not bound to the refused launch and intended invocation: %+v", recovery)
	}
	if after := snapshotRead(t, ledgerPath); !bytes.Equal(ledgerBefore, after) {
		t.Fatal("recovery spent, refunded or reset the original fixup accounting")
	}
	ledger, err := b.Store.Inspect(ctx)
	if err != nil || b.Count(ledger) != 1 {
		t.Fatalf("original charge count changed: %v %+v", err, ledger)
	}
	invBase := filepath.Join(ticket.Root, ".parley-runtime", "invocations", refused)
	for _, name := range []string{"requested.json", "terminal.json"} {
		if _, err = os.Stat(filepath.Join(invBase, name)); err != nil {
			t.Fatal("refused invocation evidence not retained", err)
		}
	}
	if _, err = os.Stat(filepath.Join(invBase, "started.json")); !os.IsNotExist(err) {
		t.Fatal("refused invocation fabricated a started lifecycle")
	}

	// The subsequent launch is separately checked at the reservation boundary:
	// the bound invocation is admitted with no write, any other is refused.
	if err = ReserveCapturedVerificationLaunch(ctx, ticket, "recovered-verifier-invocation"); err != nil {
		t.Fatalf("bound subsequent invocation was not admitted: %v", err)
	}
	if names := recoveryJournalNames(t, path); len(names) != 3 {
		t.Fatalf("admitted subsequent reservation wrote artifacts: %v", names)
	}
	if err = ReserveCapturedVerificationLaunch(ctx, ticket, "unbound-invocation"); err == nil ||
		!strings.Contains(err.Error(), "differs from its durable reservation") {
		t.Fatalf("unbound subsequent invocation was admitted: %v", err)
	}

	// Stale handles cannot control the replacement.
	for _, op := range []string{"stop", "read", "execute"} {
		var err error
		switch op {
		case "stop":
			err = StopCapturedVerification(ctx, ticket, refused)
		case "read":
			_, _, err = ReadCapturedVerification(ctx, ticket, refused)
		case "execute":
			_, err = ExecuteCapturedVerification(ctx, ticket, refused, criteria, t.TempDir())
		}
		if !errors.Is(err, errSupersededVerificationInvocation) {
			t.Fatalf("stale %s did not fail closed as superseded: %v", op, err)
		}
	}
	if _, err = os.Stat(filepath.Join(path, "stop.json")); !os.IsNotExist(err) {
		t.Fatal("superseded stop wrote a durable stop record")
	}
	if _, err = os.Stat(filepath.Join(path, "claim.json")); !os.IsNotExist(err) {
		t.Fatal("superseded execution consumed the ticket")
	}

	// Restored verification completes with the original accounting retained.
	receipt, err := ExecuteCapturedVerification(ctx, ticket, "recovered-verifier-invocation", criteria, t.TempDir())
	if err != nil || receipt.Steps != 4 || receipt.FailureStage != "" || receipt.InvocationID != "recovered-verifier-invocation" {
		t.Fatalf("recovered verification did not complete: %+v %v", receipt, err)
	}
	read, observation, err := ReadCapturedVerification(ctx, ticket, "recovered-verifier-invocation")
	if err != nil || read.TicketSHA256 != wantTicket || read.LastSHA256 != receipt.LastSHA256 {
		t.Fatalf("recovered receipt did not validate against its journal: %+v %v", read, err)
	}
	assessment, err := AssessCaptured(ticket.Request, observation)
	if err != nil || assessment.Outcome != Regression {
		t.Fatalf("recovered observation lost its material assessment: %+v %v", assessment, err)
	}
	if after := snapshotRead(t, ledgerPath); !bytes.Equal(ledgerBefore, after) {
		t.Fatal("recovered execution changed the original fixup accounting")
	}
	// The current reader keeps refusing unknown journal records whole.
	if err = os.WriteFile(filepath.Join(path, "quarantine.json"), []byte("{}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, _, err = ReadCapturedVerification(ctx, ticket, "recovered-verifier-invocation"); err == nil ||
		!strings.Contains(err.Error(), "unexpected or out-of-order") {
		t.Fatalf("unknown journal record was admitted: %v", err)
	}
	if err = os.Remove(filepath.Join(path, "quarantine.json")); err != nil {
		t.Fatal(err)
	}
	// The effective invocation retains control; the recovery grants no reissue.
	if err = StopCapturedVerification(ctx, ticket, "recovered-verifier-invocation"); err != nil {
		t.Fatalf("effective invocation lost stop control: %v", err)
	}
	if _, err = PrepareCapturedVerification(ctx, ticket.Root, "fixture", "reviewer", "another-run"); err == nil ||
		!strings.Contains(err.Error(), "verification already reserved") {
		t.Fatalf("recovery reissued the one-time reservation: %v", err)
	}
	if err = RequireResolved(ctx, ticket.Root, "fixture"); err == nil {
		t.Fatal("recovery or helper execution resolved the trajectory by itself")
	}
}

func TestRecoveryAdmissionRefusesEveryContradiction(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX captured verification")
	}
	code := 9
	for _, tc := range []struct {
		name  string
		build func(t *testing.T, ticket VerificationTicket, path string) (string, string)
	}{
		{name: "started-lifecycle", build: func(t *testing.T, ticket VerificationTicket, _ string) (string, string) {
			return refusedVerifierLaunchFixture(t, ticket, nil, true), "recovered-verifier-invocation"
		}},
		{name: "exit-code", build: func(t *testing.T, ticket VerificationTicket, _ string) (string, string) {
			return refusedVerifierLaunchFixture(t, ticket, func(_ *telemetry.Metadata, o *telemetry.Outcome) { o.ExitCode = &code }, false), "recovered-verifier-invocation"
		}},
		{name: "wrong-failure-class", build: func(t *testing.T, ticket VerificationTicket, _ string) (string, string) {
			return refusedVerifierLaunchFixture(t, ticket, func(_ *telemetry.Metadata, o *telemetry.Outcome) { o.FailureClass = telemetry.String("timeout") }, false), "recovered-verifier-invocation"
		}},
		{name: "wrong-run", build: func(t *testing.T, ticket VerificationTicket, _ string) (string, string) {
			return refusedVerifierLaunchFixture(t, ticket, func(m *telemetry.Metadata, _ *telemetry.Outcome) { m.RunID = "another-run" }, false), "recovered-verifier-invocation"
		}},
		{name: "wrong-idea", build: func(t *testing.T, ticket VerificationTicket, _ string) (string, string) {
			return refusedVerifierLaunchFixture(t, ticket, func(m *telemetry.Metadata, _ *telemetry.Outcome) { m.Idea = "another-idea" }, false), "recovered-verifier-invocation"
		}},
		{name: "wrong-agent", build: func(t *testing.T, ticket VerificationTicket, _ string) (string, string) {
			return refusedVerifierLaunchFixture(t, ticket, func(m *telemetry.Metadata, _ *telemetry.Outcome) { m.Agent = "builder" }, false), "recovered-verifier-invocation"
		}},
		{name: "wrong-phase", build: func(t *testing.T, ticket VerificationTicket, _ string) (string, string) {
			return refusedVerifierLaunchFixture(t, ticket, func(m *telemetry.Metadata, _ *telemetry.Outcome) { m.Phase = "fixup" }, false), "recovered-verifier-invocation"
		}},
		{name: "missing-request", build: func(t *testing.T, ticket VerificationTicket, _ string) (string, string) {
			refused := refusedVerifierLaunchFixture(t, ticket, nil, false)
			if err := os.Remove(filepath.Join(ticket.Root, ".parley-runtime", "invocations", refused, "requested.json")); err != nil {
				t.Fatal(err)
			}
			return refused, "recovered-verifier-invocation"
		}},
		{name: "missing-terminal", build: func(t *testing.T, ticket VerificationTicket, _ string) (string, string) {
			refused := refusedVerifierLaunchFixture(t, ticket, nil, false)
			if err := os.Remove(filepath.Join(ticket.Root, ".parley-runtime", "invocations", refused, "terminal.json")); err != nil {
				t.Fatal(err)
			}
			return refused, "recovered-verifier-invocation"
		}},
		{name: "claim-present", build: func(t *testing.T, ticket VerificationTicket, path string) (string, string) {
			if err := os.WriteFile(filepath.Join(path, "claim.json"), []byte("{}\n"), 0600); err != nil {
				t.Fatal(err)
			}
			return refusedVerifierLaunchFixture(t, ticket, nil, false), "recovered-verifier-invocation"
		}},
		{name: "receipt-present", build: func(t *testing.T, ticket VerificationTicket, path string) (string, string) {
			if err := os.WriteFile(filepath.Join(path, "receipt.json"), []byte("{}\n"), 0600); err != nil {
				t.Fatal(err)
			}
			return refusedVerifierLaunchFixture(t, ticket, nil, false), "recovered-verifier-invocation"
		}},
		{name: "stop-present", build: func(t *testing.T, ticket VerificationTicket, path string) (string, string) {
			if err := os.WriteFile(filepath.Join(path, "stop.json"), []byte("{}\n"), 0600); err != nil {
				t.Fatal(err)
			}
			return refusedVerifierLaunchFixture(t, ticket, nil, false), "recovered-verifier-invocation"
		}},
		{name: "unknown-record-present", build: func(t *testing.T, ticket VerificationTicket, path string) (string, string) {
			if err := os.WriteFile(filepath.Join(path, "stray.json"), []byte("{}\n"), 0600); err != nil {
				t.Fatal(err)
			}
			return refusedVerifierLaunchFixture(t, ticket, nil, false), "recovered-verifier-invocation"
		}},
		{name: "wrong-refused-id", build: func(t *testing.T, ticket VerificationTicket, _ string) (string, string) {
			refusedVerifierLaunchFixture(t, ticket, nil, false)
			return "not-the-refused-invocation", "recovered-verifier-invocation"
		}},
		{name: "missing-launch", build: func(t *testing.T, ticket VerificationTicket, path string) (string, string) {
			refused := refusedVerifierLaunchFixture(t, ticket, nil, false)
			if err := os.Remove(filepath.Join(path, "launch.json")); err != nil {
				t.Fatal(err)
			}
			return refused, "recovered-verifier-invocation"
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ticket, path := journalFixture(t, capturedCriteria(), dirtyCapturedChild(t))
			refused, next := tc.build(t, ticket, path)
			err := RecoverCapturedVerificationLaunch(context.Background(), ticket, refused, next)
			if err == nil {
				t.Fatalf("%s was admitted", tc.name)
			}
			if _, statErr := os.Stat(filepath.Join(path, "recovery.json")); !os.IsNotExist(statErr) {
				t.Fatalf("%s wrote a recovery artifact before refusing", tc.name)
			}
		})
	}
	t.Run("same-invocation", func(t *testing.T) {
		ticket, _ := journalFixture(t, capturedCriteria(), dirtyCapturedChild(t))
		refused := refusedVerifierLaunchFixture(t, ticket, nil, false)
		if err := RecoverCapturedVerificationLaunch(context.Background(), ticket, refused, refused); err == nil {
			t.Fatal("recovery admitted the refused invocation as its own replacement")
		}
	})
	t.Run("implementer-invocation", func(t *testing.T) {
		ticket, _ := journalFixture(t, capturedCriteria(), dirtyCapturedChild(t))
		refused := refusedVerifierLaunchFixture(t, ticket, nil, false)
		if err := RecoverCapturedVerificationLaunch(context.Background(), ticket, refused, ticket.Request.InvocationID); err == nil {
			t.Fatal("recovery admitted the implementer's own invocation")
		}
	})
	t.Run("recovery-without-reservation", func(t *testing.T) {
		ticket, _ := journalFixture(t, capturedCriteria(), dirtyCapturedChild(t))
		if err := RecoverCapturedVerificationLaunch(context.Background(), ticket, "never-reserved-invocation", "recovered-verifier-invocation"); err == nil {
			t.Fatal("recovery admitted a launch that was never reserved")
		}
	})
}

func TestRepeatedAndConflictingRecoveryApplyFailsClosed(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX captured verification")
	}
	ctx := context.Background()
	ticket, path := journalFixture(t, capturedCriteria(), dirtyCapturedChild(t))
	refused := refusedVerifierLaunchFixture(t, ticket, nil, false)
	if err := RecoverCapturedVerificationLaunch(ctx, ticket, refused, "recovered-verifier-invocation"); err != nil {
		t.Fatal(err)
	}
	retained := snapshotRead(t, filepath.Join(path, "recovery.json"))
	for _, next := range []string{"recovered-verifier-invocation", "conflicting-invocation"} {
		if applyErr := RecoverCapturedVerificationLaunch(ctx, ticket, refused, next); applyErr == nil ||
			!strings.Contains(applyErr.Error(), "already recovered") {
			t.Fatalf("repeated or conflicting apply was not refused: %v", applyErr)
		}
		if after := snapshotRead(t, filepath.Join(path, "recovery.json")); !bytes.Equal(retained, after) {
			t.Fatal("refused reapply rewrote the retained recovery artifact")
		}
	}
	// Contradictory retained artifacts are refused by every reader, not repaired.
	// Each variant starts from a genuinely admitted recovery so only the varied
	// binding is wrong. Readers revalidate the retained refusal evidence against
	// the recovery's bound digests, so a correct structural binding carrying
	// wrong evidence digests is refused too; only a fully self-consistent
	// same-UID forgery (recomputed digests over fabricated canonical evidence)
	// remains beyond byte evidence — the documented honesty boundary.
	for _, variant := range []struct {
		name    string
		mutate  func(VerificationTicket, *verificationRecovery)
		message string
	}{
		{name: "foreign-prior-launch", mutate: func(_ VerificationTicket, r *verificationRecovery) {
			r.PriorLaunchSHA256 = digest([]byte("forged"))
		}, message: "contradicts its original launch reservation"},
		{name: "wrong-evidence-digests", mutate: func(_ VerificationTicket, r *verificationRecovery) {
			r.RefusedRequestedSHA256 = digest([]byte("forged"))
			r.RefusedTerminalSHA256 = digest([]byte("forged"))
		}, message: "no longer matches its retained refused-launch evidence"},
		{name: "reversed-chronology", mutate: func(_ VerificationTicket, r *verificationRecovery) {
			r.At = r.At.Add(-24 * time.Hour)
		}, message: "contradicts its original launch reservation"},
		{name: "self-bound-invocation", mutate: func(_ VerificationTicket, r *verificationRecovery) {
			r.InvocationID = r.RefusedInvocationID
		}, message: "contradicts its original launch reservation"},
		{name: "implementer-bound-invocation", mutate: func(ticket VerificationTicket, r *verificationRecovery) {
			r.InvocationID = ticket.Request.InvocationID
		}, message: "contradicts its original launch reservation"},
	} {
		t.Run(variant.name, func(t *testing.T) {
			ticket, path := journalFixture(t, capturedCriteria(), dirtyCapturedChild(t))
			refused := refusedVerifierLaunchFixture(t, ticket, nil, false)
			if err := RecoverCapturedVerificationLaunch(ctx, ticket, refused, "recovered-verifier-invocation"); err != nil {
				t.Fatal(err)
			}
			var admitted verificationRecovery
			if err := json.Unmarshal(snapshotRead(t, filepath.Join(path, "recovery.json")), &admitted); err != nil {
				t.Fatal(err)
			}
			variant.mutate(ticket, &admitted)
			forged, err := canonical(admitted)
			if err != nil {
				t.Fatal(err)
			}
			if err = os.WriteFile(filepath.Join(path, "recovery.json"), forged, 0600); err != nil {
				t.Fatal(err)
			}
			if err = RecoverCapturedVerificationLaunch(ctx, ticket, refused, "recovered-verifier-invocation"); err == nil {
				t.Fatal("contradictory recovery state was repaired in place")
			}
			for _, op := range []func() error{
				func() error { return StopCapturedVerification(ctx, ticket, "recovered-verifier-invocation") },
				func() error {
					_, _, err := ReadCapturedVerification(ctx, ticket, "recovered-verifier-invocation")
					return err
				},
				func() error {
					_, err := ExecuteCapturedVerification(ctx, ticket, "recovered-verifier-invocation", capturedCriteria(), t.TempDir())
					return err
				},
			} {
				if err = op(); err == nil || !strings.Contains(err.Error(), variant.message) {
					t.Fatalf("contradictory recovery artifact was not refused by every reader: %v", err)
				}
			}
		})
	}
}

func TestConcurrentRecoveryAppliesExactlyOnce(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX captured verification")
	}
	ctx := context.Background()
	ticket, path := journalFixture(t, capturedCriteria(), dirtyCapturedChild(t))
	refused := refusedVerifierLaunchFixture(t, ticket, nil, false)
	newInvocations := []string{"recovered-a", "recovered-b"}
	errs := make([]error, 2)
	var wg sync.WaitGroup
	for i := range newInvocations {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			errs[i] = RecoverCapturedVerificationLaunch(ctx, ticket, refused, newInvocations[i])
		}(i)
	}
	wg.Wait()
	if (errs[0] == nil) == (errs[1] == nil) {
		t.Fatalf("exactly one concurrent recovery must win: %v / %v", errs[0], errs[1])
	}
	winner := 0
	if errs[1] == nil {
		winner = 1
	}
	var recovery verificationRecovery
	if err := json.Unmarshal(snapshotRead(t, filepath.Join(path, "recovery.json")), &recovery); err != nil {
		t.Fatal(err)
	}
	if recovery.InvocationID != newInvocations[winner] {
		t.Fatalf("retained recovery does not bind the winning invocation: %+v", recovery)
	}
	// The losing invocation cannot execute: it never became authority.
	if _, err := ExecuteCapturedVerification(ctx, ticket, newInvocations[1-winner], capturedCriteria(), t.TempDir()); err == nil {
		t.Fatal("losing concurrent recovery invocation executed")
	}
}

// TestRecoveredJournalFailsClosedForPreRecoveryReaders: a reader built before
// recovery admits no journal record outside its inventory, so the recovered
// journal is refused whole by pre-recovery binaries — never partially
// interpreted behind the new artifact's back.
func TestRecoveredJournalFailsClosedForPreRecoveryReaders(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX captured verification")
	}
	ticket, path := journalFixture(t, capturedCriteria(), dirtyCapturedChild(t))
	refused := refusedVerifierLaunchFixture(t, ticket, nil, false)
	if err := RecoverCapturedVerificationLaunch(context.Background(), ticket, refused, "recovered-verifier-invocation"); err != nil {
		t.Fatal(err)
	}
	oldAllowed := map[string]bool{"request.json": true, "launch.json": true, "claim.json": true,
		"prepared.json": true, "receipt.json": true, "stop.json": true}
	entries, err := os.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}
	var unknown []string
	for _, e := range entries {
		if !oldAllowed[e.Name()] && !strings.HasPrefix(e.Name(), "step-") && !strings.HasPrefix(e.Name(), "process-") {
			unknown = append(unknown, e.Name())
		}
	}
	if len(unknown) != 1 || unknown[0] != "recovery.json" {
		t.Fatalf("pre-recovery readers would silently admit the recovered journal: %v", unknown)
	}
}

// TestFreshResolutionValidatesRecoveredLineageAndFailsClosed: the verifier-run
// resolution follows the immutable recovery and re-derives the COMPLETE launch
// lineage. A helper-only journal — claim, steps and receipt under the bound
// replacement invocation — is not resolution evidence on its own: without the
// replacement invocation's own observed requested/started/terminal lifecycle
// the derivation fails closed at the terminal binding and the trajectory stays
// unresolved. Partial helper evidence never becomes acceptance of
// unimplemented functionality.
func TestFreshResolutionValidatesRecoveredLineageAndFailsClosed(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX captured verification")
	}
	ctx := context.Background()
	criteria := []Criterion{{Name: "material", Command: `printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":0}\n'`}}
	ticket, _ := journalFixture(t, criteria, dirtyCapturedChild(t))
	refused := refusedVerifierLaunchFixture(t, ticket, nil, false)
	if err := RecoverCapturedVerificationLaunch(ctx, ticket, refused, "recovered-verifier-invocation"); err != nil {
		t.Fatal(err)
	}
	if _, err := ExecuteCapturedVerification(ctx, ticket, "recovered-verifier-invocation", criteria, t.TempDir()); err != nil {
		t.Fatal(err)
	}
	// Reconstruct the verifier-run resolution inputs exactly as the app retains
	// them, with the live scope still matching the frozen criteria.
	scope := "---\nparticipants: [builder, reviewer]\nchecks:\n  - name: material\n    command: >\n      " + criteria[0].Command + "\n---\n"
	snapshotWrite(t, ticket.Root, "parley-deck/ideas/fixture/00-prompt.md", []byte(scope), 0600)
	req := HelperRequest{Version: 1, Ticket: ticket, Participants: []string{"builder", "reviewer"}, Criteria: criteria}
	data, err := canonical(req)
	if err != nil {
		t.Fatal(err)
	}
	base := filepath.Join(ticket.Root, ".parley-runtime", "trajectory-verification", ticket.RunID)
	if err = os.MkdirAll(base, 0700); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(base, "request.json"), data, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = PreviewReconciliation(ctx, ticket.Root, "fixture", ticket.RunID); err == nil ||
		!strings.Contains(err.Error(), "successful observed parent terminal binding is unavailable") {
		t.Fatalf("fresh resolution did not validate the recovered launch lineage: %v", err)
	}
	if err = RequireResolved(ctx, ticket.Root, "fixture"); err == nil {
		t.Fatal("recovered helper execution closed the trajectory")
	}
}

// rewriteRefusedEvidenceCanonical rewrites both retained refusal records with
// telemetry's exact canonical encoding, changing only an ordinal no admission
// or reader rule inspects. The result is still a semantically valid pre-start
// budget refusal, so only the recovery's bound-digest equality can catch it.
func rewriteRefusedEvidenceCanonical(t *testing.T, root, refused string) {
	t.Helper()
	base := filepath.Join(root, ".parley-runtime", "invocations", refused)
	for _, name := range []string{"requested.json", "terminal.json"} {
		data, err := os.ReadFile(filepath.Join(base, name))
		if err != nil {
			t.Fatal(err)
		}
		var record telemetry.Record
		if err = json.Unmarshal(data, &record); err != nil {
			t.Fatal(err)
		}
		record.Metadata.AttemptOrdinal = 99
		rewritten, err := json.MarshalIndent(record, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(filepath.Join(base, name), append(rewritten, '\n'), 0600); err != nil {
			t.Fatal(err)
		}
	}
}

// TestRecoveredRefusedEvidenceMutationFailsEveryHandle: the recovery binds the
// refused launch's exact retained evidence, and every recovered operation
// revalidates it against the bytes retained now. Mutating or deleting that
// evidence after recovery fails both stale (refused) and current (recovered)
// handles across reserve, stop, read and execute, with no journal write. The
// rewritten case stays canonical and self-consistent, so its rejection proves
// the digest-equality binding rather than a parse failure.
func TestRecoveredRefusedEvidenceMutationFailsEveryHandle(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX captured verification")
	}
	ctx := context.Background()
	criteria := capturedCriteria()
	for _, tc := range []struct {
		name    string
		mutate  func(t *testing.T, root, refused string)
		message string
	}{
		{name: "evidence-rewritten", mutate: func(t *testing.T, root, refused string) {
			rewriteRefusedEvidenceCanonical(t, root, refused)
		}, message: "no longer matches its retained refused-launch evidence"},
		{name: "requested-deleted", mutate: func(t *testing.T, root, refused string) {
			if err := os.Remove(filepath.Join(root, ".parley-runtime", "invocations", refused, "requested.json")); err != nil {
				t.Fatal(err)
			}
		}, message: "refused verifier launch request is unavailable or changed"},
		{name: "terminal-deleted", mutate: func(t *testing.T, root, refused string) {
			if err := os.Remove(filepath.Join(root, ".parley-runtime", "invocations", refused, "terminal.json")); err != nil {
				t.Fatal(err)
			}
		}, message: "refused verifier launch lacks its exact pre-start budget-refusal terminal"},
		{name: "started-lifecycle-added", mutate: func(t *testing.T, root, refused string) {
			if err := os.WriteFile(filepath.Join(root, ".parley-runtime", "invocations", refused, "started.json"), []byte("{}\n"), 0600); err != nil {
				t.Fatal(err)
			}
		}, message: "must not retain any started lifecycle"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ticket, path := journalFixture(t, criteria, dirtyCapturedChild(t))
			refused := refusedVerifierLaunchFixture(t, ticket, nil, false)
			if err := RecoverCapturedVerificationLaunch(ctx, ticket, refused, "recovered-verifier-invocation"); err != nil {
				t.Fatal(err)
			}
			tc.mutate(t, ticket.Root, refused)
			for _, op := range []struct {
				label string
				run   func() error
			}{
				{label: "reserve-current", run: func() error {
					return ReserveCapturedVerificationLaunch(ctx, ticket, "recovered-verifier-invocation")
				}},
				{label: "stop-current", run: func() error {
					return StopCapturedVerification(ctx, ticket, "recovered-verifier-invocation")
				}},
				{label: "stop-stale", run: func() error {
					return StopCapturedVerification(ctx, ticket, refused)
				}},
				{label: "read-current", run: func() error {
					_, _, err := ReadCapturedVerification(ctx, ticket, "recovered-verifier-invocation")
					return err
				}},
				{label: "execute-current", run: func() error {
					_, err := ExecuteCapturedVerification(ctx, ticket, "recovered-verifier-invocation", criteria, t.TempDir())
					return err
				}},
			} {
				if err := op.run(); err == nil || !strings.Contains(err.Error(), tc.message) {
					t.Fatalf("%s did not fail closed on changed retained lineage: %v", op.label, err)
				}
			}
			if _, err := os.Stat(filepath.Join(path, "stop.json")); !os.IsNotExist(err) {
				t.Fatal("refused operation wrote a stop record")
			}
			if _, err := os.Stat(filepath.Join(path, "claim.json")); !os.IsNotExist(err) {
				t.Fatal("refused operation wrote a helper claim")
			}
			if names := recoveryJournalNames(t, path); len(names) != 3 || !names["request.json"] || !names["launch.json"] || !names["recovery.json"] {
				t.Fatalf("refused operations changed the retained journal inventory: %v", names)
			}
		})
	}
}

// TestRecoveredReceiptRevalidatesRefusedEvidence: a complete, journal-valid
// helper receipt still cannot be read over changed refusal evidence — the
// receipt read resolves the launch authority first and revalidates the whole
// retained recovery lineage, not just the receipt's own bindings.
func TestRecoveredReceiptRevalidatesRefusedEvidence(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX captured verification")
	}
	ctx := context.Background()
	criteria := capturedCriteria()
	ticket, _ := journalFixture(t, criteria, dirtyCapturedChild(t))
	refused := refusedVerifierLaunchFixture(t, ticket, nil, false)
	if err := RecoverCapturedVerificationLaunch(ctx, ticket, refused, "recovered-verifier-invocation"); err != nil {
		t.Fatal(err)
	}
	if _, err := ExecuteCapturedVerification(ctx, ticket, "recovered-verifier-invocation", criteria, t.TempDir()); err != nil {
		t.Fatal(err)
	}
	if _, _, err := ReadCapturedVerification(ctx, ticket, "recovered-verifier-invocation"); err != nil {
		t.Fatal(err)
	}
	rewriteRefusedEvidenceCanonical(t, ticket.Root, refused)
	if _, _, err := ReadCapturedVerification(ctx, ticket, "recovered-verifier-invocation"); err == nil ||
		!strings.Contains(err.Error(), "no longer matches its retained refused-launch evidence") {
		t.Fatalf("valid receipt was read over changed refused-launch evidence: %v", err)
	}
}
