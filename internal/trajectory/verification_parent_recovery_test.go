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
	"time"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/telemetry"
)

// The recovered parent-lineage fixtures reconstruct, through public APIs where
// they exist, exactly what a verifier-run resolution retains after an explicit
// verifier-launch recovery: the real refused launch and admission, the real
// helper execution under the bound replacement invocation, and the replacement
// invocation's own observed requested/started/terminal lifecycle retained as
// telemetry bytes. telemetry.Begin cannot predetermine an invocation ID, so
// the replacement lifecycle is written here with telemetry's exact canonical
// encoding around the real helper execution — isolating the derivation reader.
// The runner-side bound relaunch (telemetry.BeginBound under
// WithCapturedVerificationRecovery) now produces these records natively and is
// tested at the runner level; only the attended app CLI wiring remains
// separate.

type recoveredParentLineage struct {
	ticket        VerificationTicket
	path          string
	refused       string
	replacement   string
	b             *budget.CycleBinding
	ledgerPath    string
	ledgerBefore  []byte
	launchBefore  []byte
	recoveryBytes []byte
	preview       ReconciliationPreview
}

// retainedReplacementLifecycle writes the replacement invocation's observed
// lifecycle exactly as telemetry retains it — requested and started before the
// helper runs so its receipt lands inside the observed window, the successful
// terminal last — and returns the successful terminal's retained digest.
// lifecycleTimes models the recovered chronology: nil uses real now-timestamps
// (the fresh launch begins after the recovery); a non-nil override backdates
// the replacement to contradict the recovery.
func retainedReplacementLifecycle(t *testing.T, ticket VerificationTicket, replacement string, lifecycleTimes func(recoveryAt time.Time) (time.Time, time.Time), recoveryAt time.Time, runHelper func()) string {
	t.Helper()
	requestedAt, startedAt := time.Now().UTC(), time.Now().UTC()
	if lifecycleTimes != nil {
		requestedAt, startedAt = lifecycleTimes(recoveryAt)
	}
	metadata := telemetry.Metadata{RunID: ticket.RunID, Idea: ticket.Request.Idea, Phase: "trajectory-verification",
		Agent: ticket.Request.Verifier, Adapter: "claude", LaunchMode: "headless", AttemptOrdinal: 1}
	invBase := filepath.Join(ticket.Root, ".parley-runtime", "invocations", replacement)
	if err := os.MkdirAll(invBase, 0700); err != nil {
		t.Fatal(err)
	}
	pid := os.Getpid()
	requested := telemetry.Record{SchemaVersion: telemetry.SchemaVersion, Type: "invocation.requested",
		InvocationID: replacement, Metadata: metadata, RequestedAt: requestedAt}
	started := requested
	started.Type = "invocation.started"
	started.StartedAt, started.PID = &startedAt, &pid
	write := func(name string, record telemetry.Record) []byte {
		data, err := canonical(record)
		if err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(filepath.Join(invBase, name), data, 0600); err != nil {
			t.Fatal(err)
		}
		return data
	}
	write("requested.json", requested)
	write("started.json", started)
	runHelper()
	terminal := requested
	terminal.Type = "invocation.terminal"
	completedAt := time.Now().UTC()
	terminal.StartedAt, terminal.PID, terminal.CompletedAt = &startedAt, &pid, &completedAt
	exitCode := 0
	terminal.Outcome = &telemetry.Outcome{Status: "process-exited", ExitCode: &exitCode}
	return digest(write("terminal.json", terminal))
}

// recoveredParentLineageFixture builds the complete recovered lineage and, for
// a real-now chronology (lifecycleTimes == nil), proves the clean preview: the
// parent publishes exactly the facts the derivation must reproduce, and fresh
// PreviewReconciliation resolves through the recovered lineage. A backdating
// override skips that control — its preview must fail, which the caller
// asserts.
func recoveredParentLineageFixture(t *testing.T, lifecycleTimes func(recoveryAt time.Time) (time.Time, time.Time)) recoveredParentLineage {
	t.Helper()
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
	refused := refusedVerifierLaunchFixture(t, ticket, nil, false)
	launchBefore := snapshotRead(t, filepath.Join(path, "launch.json"))
	const replacement = "recovered-verifier-invocation"
	if err = RecoverCapturedVerificationLaunch(ctx, ticket, refused, replacement); err != nil {
		t.Fatal(err)
	}
	var recovery verificationRecovery
	recoveryBytes := snapshotRead(t, filepath.Join(path, "recovery.json"))
	if err = json.Unmarshal(recoveryBytes, &recovery); err != nil {
		t.Fatal(err)
	}
	// Reconstruct the verifier-run resolution inputs exactly as the app
	// retains them, with the live scope still matching the frozen criteria.
	scope := "---\nparticipants: [builder, reviewer]\nchecks:\n  - name: material\n    command: >\n      " + criteria[0].Command + "\n---\n"
	snapshotWrite(t, ticket.Root, "parley-deck/ideas/fixture/00-prompt.md", []byte(scope), 0600)
	reqData, err := canonical(HelperRequest{Version: 1, Ticket: ticket, Participants: []string{"builder", "reviewer"}, Criteria: criteria})
	if err != nil {
		t.Fatal(err)
	}
	base := filepath.Join(ticket.Root, ".parley-runtime", "trajectory-verification", ticket.RunID)
	if err = os.MkdirAll(base, 0700); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(base, "request.json"), reqData, 0600); err != nil {
		t.Fatal(err)
	}
	var receipt VerificationReceipt
	terminalSHA := retainedReplacementLifecycle(t, ticket, replacement, lifecycleTimes, recovery.At, func() {
		var execErr error
		receipt, execErr = ExecuteCapturedVerification(ctx, ticket, replacement, criteria, t.TempDir())
		if execErr != nil || receipt.Steps != 4 || receipt.FailureStage != "" || receipt.InvocationID != replacement {
			t.Fatalf("recovered helper execution did not complete: %+v %v", receipt, execErr)
		}
	})
	// The parent publishes exactly the facts the derivation must reproduce;
	// any divergence is refused by the ordinary parent-result binding.
	read, observation, err := ReadCapturedVerification(ctx, ticket, replacement)
	if err != nil {
		t.Fatal(err)
	}
	wantTicket, err := ticket.SHA256()
	if err != nil {
		t.Fatal(err)
	}
	if read.TicketSHA256 != wantTicket || read.LastSHA256 != receipt.LastSHA256 {
		t.Fatalf("recovered receipt did not validate against its journal: %+v", read)
	}
	assessment, err := AssessCaptured(ticket.Request, observation)
	if err != nil || assessment.Outcome != Regression {
		t.Fatalf("recovered observation lost its material assessment: %+v %v", assessment, err)
	}
	encoded, err := canonical(read)
	if err != nil {
		t.Fatal(err)
	}
	result := ParentResult{Version: 1, RunID: ticket.RunID, RequestPath: filepath.Join(base, "request.json"),
		RequestSHA256: digest(reqData), InvocationID: replacement, TerminalSHA256: terminalSHA,
		ReceiptSHA256: digest(encoded), Assessment: &assessment, TrajectoryPending: true}
	resultData, err := canonical(result)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(base, "parent-result.json"), resultData, 0600); err != nil {
		t.Fatal(err)
	}
	f := recoveredParentLineage{ticket: ticket, path: path, refused: refused, replacement: replacement,
		b: b, ledgerPath: ledgerPath, ledgerBefore: ledgerBefore, launchBefore: launchBefore, recoveryBytes: recoveryBytes}
	if lifecycleTimes != nil {
		return f
	}
	preview, err := PreviewReconciliation(ctx, ticket.Root, "fixture", ticket.RunID)
	if err != nil {
		t.Fatalf("clean recovered lineage did not resolve: %v", err)
	}
	if preview.Version != 1 || preview.Sequence != 1 || preview.Root != ticket.Root || preview.RunID != ticket.RunID ||
		preview.Assessment.Outcome != Regression || preview.ParentSHA256 != digest(resultData) || preview.RecoverySHA256 != "" {
		t.Fatalf("clean recovered preview is not the derived parent binding: %+v", preview)
	}
	f.preview = preview
	return f
}

// TestRecoveredReplacementLifecycleResolvesAndReconciles: a fully observed
// replacement lifecycle plus the real local helper criteria under the bound
// invocation produce a valid fresh preview and an explicit resolution with the
// original accounting, refusal records and reservation retained byte-exact.
// Later state operations re-derive the recorded resolution through the same
// recovered lineage (checkResolutions runs on every state guard), so the
// resolution stays bound to the evidence it was accepted on.
func TestRecoveredReplacementLifecycleResolvesAndReconciles(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX captured verification")
	}
	ctx := context.Background()
	f := recoveredParentLineageFixture(t, nil)
	if err := RequireResolved(ctx, f.ticket.Root, "fixture"); err == nil {
		t.Fatal("recovered lineage resolved the trajectory without an explicit reconciliation")
	}
	expected := f.preview.SHA256()
	if err := Reconcile(ctx, f.ticket.Root, "fixture", f.ticket.RunID, expected); err != nil {
		t.Fatalf("explicit resolution refused the recovered lineage: %v", err)
	}
	if err := Reconcile(ctx, f.ticket.Root, "fixture", f.ticket.RunID, expected); err != nil {
		t.Fatalf("exact resolution replay changed its original decision: %v", err)
	}
	replayed, err := PreviewReconciliation(ctx, f.ticket.Root, "fixture", f.ticket.RunID)
	if err != nil || replayed.SHA256() != expected {
		t.Fatalf("retained resolution did not replay exactly: %+v %v", replayed, err)
	}
	state, err := Inspect(ctx, f.ticket.Root, "fixture")
	if err != nil || len(state.Resolutions) != 1 || state.Resolutions[0].SHA256 != expected ||
		state.Resolutions[0].Preview.ChargeKey != state.Attempts[0].Charge.EntryKey {
		t.Fatalf("recovered resolution was not retained verbatim: %v %+v", err, state)
	}
	if !bytes.Equal(f.ledgerBefore, snapshotRead(t, f.ledgerPath)) {
		t.Fatal("recovered resolution spent, refunded or reset the original fixup accounting")
	}
	ledger, err := f.b.Store.Inspect(ctx)
	if err != nil || f.b.Count(ledger) != 1 {
		t.Fatalf("original charge count changed: %v %+v", err, ledger)
	}
	if !bytes.Equal(f.launchBefore, snapshotRead(t, filepath.Join(f.path, "launch.json"))) ||
		!bytes.Equal(f.recoveryBytes, snapshotRead(t, filepath.Join(f.path, "recovery.json"))) {
		t.Fatal("resolution rewrote the retained launch reservation or recovery artifact")
	}
	invBase := filepath.Join(f.ticket.Root, ".parley-runtime", "invocations", f.refused)
	for _, name := range []string{"requested.json", "terminal.json"} {
		if _, statErr := os.Stat(filepath.Join(invBase, name)); statErr != nil {
			t.Fatal("refused invocation evidence not retained", statErr)
		}
	}
	if _, statErr := os.Stat(filepath.Join(invBase, "started.json")); !os.IsNotExist(statErr) {
		t.Fatal("refused invocation fabricated a started lifecycle")
	}
	// The material regression outcome legitimately keeps the trajectory open
	// for the fix-up loop; the resolution itself is complete and retained.
	if err = RequireResolved(ctx, f.ticket.Root, "fixture"); err == nil {
		t.Fatal("regression resolution closed the trajectory without a clean outcome")
	}
}

// TestRecoveredParentResolutionRequiresExactLineage: every mutation of the
// retained lineage after the clean preview fails the fresh derivation, keeps
// the explicit resolution unavailable, and never auto-resolves.
func TestRecoveredParentResolutionRequiresExactLineage(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX captured verification")
	}
	ctx := context.Background()
	for _, tc := range []struct {
		name    string
		times   func(time.Time) (time.Time, time.Time)
		mutate  func(t *testing.T, f recoveredParentLineage)
		message string
		handle  func(t *testing.T, f recoveredParentLineage)
	}{
		{name: "old-requested-deleted", mutate: func(t *testing.T, f recoveredParentLineage) {
			removePath(t, filepath.Join(f.ticket.Root, ".parley-runtime", "invocations", f.refused, "requested.json"))
		}, message: "refused verifier launch request is unavailable or changed"},
		{name: "old-terminal-deleted", mutate: func(t *testing.T, f recoveredParentLineage) {
			removePath(t, filepath.Join(f.ticket.Root, ".parley-runtime", "invocations", f.refused, "terminal.json"))
		}, message: "refused verifier launch lacks its exact pre-start budget-refusal terminal"},
		{name: "old-evidence-rewritten", mutate: func(t *testing.T, f recoveredParentLineage) {
			rewriteRefusedEvidenceCanonical(t, f.ticket.Root, f.refused)
		}, message: "no longer matches its retained refused-launch evidence"},
		{name: "old-launch-deleted", mutate: func(t *testing.T, f recoveredParentLineage) {
			removePath(t, filepath.Join(f.path, "launch.json"))
		}, message: "original verifier launch is unavailable"},
		{name: "recovery-deleted", mutate: func(t *testing.T, f recoveredParentLineage) {
			removePath(t, filepath.Join(f.path, "recovery.json"))
		}, message: "successful observed parent terminal binding is unavailable"},
		{name: "recovery-rebinds-decoy", mutate: func(t *testing.T, f recoveredParentLineage) {
			rewriteRecoveryInvocation(t, f, "decoy-verifier-invocation")
		}, message: "successful observed parent terminal binding is unavailable",
			handle: func(t *testing.T, f recoveredParentLineage) {
				if _, _, err := ReadCapturedVerification(ctx, f.ticket, f.replacement); err == nil ||
					!strings.Contains(err.Error(), "differs from its durable reservation") {
					t.Fatalf("rebind stale handle did not fail closed: %v", err)
				}
			}},
		{name: "replacement-lifecycle-moved", mutate: func(t *testing.T, f recoveredParentLineage) {
			if err := os.Rename(filepath.Join(f.ticket.Root, ".parley-runtime", "invocations", f.replacement),
				filepath.Join(f.ticket.Root, ".parley-runtime", "invocations", f.replacement+"-moved")); err != nil {
				t.Fatal(err)
			}
		}, message: "successful observed parent terminal binding is unavailable"},
		{name: "recovery-predates-refusal-completion", mutate: func(t *testing.T, f recoveredParentLineage) {
			backdateRecoveryIntoRefusalWindow(t, f)
		}, message: "verification recovery predates the retained refusal completion",
			handle: func(t *testing.T, f recoveredParentLineage) {
				if _, _, err := ReadCapturedVerification(ctx, f.ticket, f.replacement); err == nil ||
					!strings.Contains(err.Error(), "predates the retained refusal completion") {
					t.Fatalf("predated recovery handle did not fail closed: %v", err)
				}
			}},
		{name: "replacement-lifecycle-predates-recovery",
			times: func(recoveryAt time.Time) (time.Time, time.Time) {
				return recoveryAt.Add(-time.Hour), recoveryAt.Add(-30 * time.Minute)
			}, message: "original started lifecycle differs from the verifier terminal"},
		{name: "missing-replacement-terminal", mutate: func(t *testing.T, f recoveredParentLineage) {
			removePath(t, filepath.Join(f.ticket.Root, ".parley-runtime", "invocations", f.replacement, "terminal.json"))
		}, message: "successful observed parent terminal binding is unavailable"},
		{name: "unknown-journal-record", mutate: func(t *testing.T, f recoveredParentLineage) {
			if err := os.WriteFile(filepath.Join(f.path, "quarantine.json"), []byte("{}\n"), 0600); err != nil {
				t.Fatal(err)
			}
		}, message: "unexpected or out-of-order"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := recoveredParentLineageFixture(t, tc.times)
			if tc.mutate != nil {
				tc.mutate(t, f)
			}
			_, err := PreviewReconciliation(ctx, f.ticket.Root, "fixture", f.ticket.RunID)
			if err == nil || !strings.Contains(err.Error(), tc.message) {
				t.Fatalf("%s did not fail closed: %v", tc.name, err)
			}
			if tc.handle != nil {
				tc.handle(t, f)
			}
			if err = RequireResolved(ctx, f.ticket.Root, "fixture"); err == nil {
				t.Fatalf("%s auto-resolved the trajectory", tc.name)
			}
			state, inspectErr := Inspect(ctx, f.ticket.Root, "fixture")
			if inspectErr != nil || len(state.Resolutions) != 0 {
				t.Fatalf("%s recorded a resolution: %v %+v", tc.name, inspectErr, state)
			}
		})
	}
}

// TestRecoveredParentStalePreviewFailsExactRecheck: a preview accepted on the
// clean recovered lineage does not survive mutation of that lineage — the
// explicit acceptance re-derives every fact, so the stale expected SHA is
// refused and no resolution is written.
func TestRecoveredParentStalePreviewFailsExactRecheck(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX captured verification")
	}
	ctx := context.Background()
	for _, tc := range []struct {
		name    string
		mutate  func(t *testing.T, f recoveredParentLineage)
		message string
	}{
		{name: "refused-evidence-rewritten", mutate: func(t *testing.T, f recoveredParentLineage) {
			rewriteRefusedEvidenceCanonical(t, f.ticket.Root, f.refused)
		}, message: "no longer matches its retained refused-launch evidence"},
		{name: "recovery-deleted", mutate: func(t *testing.T, f recoveredParentLineage) {
			removePath(t, filepath.Join(f.path, "recovery.json"))
		}, message: "successful observed parent terminal binding is unavailable"},
		{name: "parent-result-rewritten", mutate: func(t *testing.T, f recoveredParentLineage) {
			path := filepath.Join(f.ticket.Root, ".parley-runtime", "trajectory-verification", f.ticket.RunID, "parent-result.json")
			var result ParentResult
			if err := json.Unmarshal(snapshotRead(t, path), &result); err != nil {
				t.Fatal(err)
			}
			result.ReceiptSHA256 = digest([]byte("forged"))
			data, err := canonical(result)
			if err != nil {
				t.Fatal(err)
			}
			if err = os.WriteFile(path, data, 0600); err != nil {
				t.Fatal(err)
			}
		}, message: "parent did not retain the independently derived successful comparison"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := recoveredParentLineageFixture(t, nil)
			expected := f.preview.SHA256()
			tc.mutate(t, f)
			if err := Reconcile(ctx, f.ticket.Root, "fixture", f.ticket.RunID, expected); err == nil ||
				!strings.Contains(err.Error(), tc.message) {
				t.Fatalf("stale preview was accepted: %v", err)
			}
			if err := RequireResolved(ctx, f.ticket.Root, "fixture"); err == nil {
				t.Fatalf("%s auto-resolved the trajectory", tc.name)
			}
			state, err := Inspect(ctx, f.ticket.Root, "fixture")
			if err != nil || len(state.Resolutions) != 0 {
				t.Fatalf("%s recorded a resolution: %v %+v", tc.name, err, state)
			}
		})
	}
}

func removePath(t *testing.T, path string) {
	t.Helper()
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
}

// rewriteRecoveryInvocation rewrites the admitted recovery artifact, canonically
// and with every other binding intact, to name a different replacement.
func rewriteRecoveryInvocation(t *testing.T, f recoveredParentLineage, invocation string) {
	t.Helper()
	var admitted verificationRecovery
	if err := json.Unmarshal(f.recoveryBytes, &admitted); err != nil {
		t.Fatal(err)
	}
	admitted.InvocationID = invocation
	data, err := canonical(admitted)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(f.path, "recovery.json"), data, 0600); err != nil {
		t.Fatal(err)
	}
}

// backdateRecoveryIntoRefusalWindow rewrites the admitted recovery artifact so
// its At falls strictly between the original reservation and the refusal
// terminal's completion — a recovery that predates the refusal it repairs.
func backdateRecoveryIntoRefusalWindow(t *testing.T, f recoveredParentLineage) {
	t.Helper()
	var launch verificationLaunch
	if err := json.Unmarshal(snapshotRead(t, filepath.Join(f.path, "launch.json")), &launch); err != nil {
		t.Fatal(err)
	}
	var refusedTerminal telemetry.Record
	terminalErr := json.Unmarshal(snapshotRead(t, filepath.Join(f.ticket.Root, ".parley-runtime", "invocations", f.refused, "terminal.json")), &refusedTerminal)
	if terminalErr != nil || refusedTerminal.CompletedAt == nil || !refusedTerminal.CompletedAt.After(launch.At) {
		t.Fatalf("fixture refusal did not complete after its reservation: %v %+v", terminalErr, refusedTerminal)
	}
	var admitted verificationRecovery
	if err := json.Unmarshal(f.recoveryBytes, &admitted); err != nil {
		t.Fatal(err)
	}
	admitted.At = launch.At.Add(refusedTerminal.CompletedAt.Sub(launch.At) / 2)
	data, err := canonical(admitted)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(f.path, "recovery.json"), data, 0600); err != nil {
		t.Fatal(err)
	}
}

// refusedParentRecoveryLineage builds the complete app-lifecycle evidence the
// recovered parent observation binds: the real refused launch, the real
// recovery, the retained run request, the retained ORIGINAL failed parent
// result exactly as the production verify path writes it on a pre-start budget
// refusal, and the fully observed replacement/helper lineage.
type refusedParentRecoveryLineage struct {
	ticket      VerificationTicket
	path        string
	base        string
	refused     string
	replacement string
	original    ParentResult
	originalRaw []byte
}

func refusedParentRecoveryFixture(t *testing.T) refusedParentRecoveryLineage {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("POSIX captured verification")
	}
	ctx := context.Background()
	criteria := capturedCriteria()
	ticket, path := journalFixture(t, criteria, dirtyCapturedChild(t))
	refused := refusedVerifierLaunchFixture(t, ticket, nil, false)
	const replacement = "recovered-verifier-invocation"
	if err := RecoverCapturedVerificationLaunch(ctx, ticket, refused, replacement); err != nil {
		t.Fatal(err)
	}
	var recovery verificationRecovery
	if err := json.Unmarshal(snapshotRead(t, filepath.Join(path, "recovery.json")), &recovery); err != nil {
		t.Fatal(err)
	}
	scope := "---\nparticipants: [builder, reviewer]\nchecks:\n  - name: material\n    command: >\n      " + criteria[0].Command + "\n---\n"
	snapshotWrite(t, ticket.Root, "parley-deck/ideas/fixture/00-prompt.md", []byte(scope), 0600)
	reqData, err := canonical(HelperRequest{Version: 1, Ticket: ticket, Participants: []string{"builder", "reviewer"}, Criteria: criteria})
	if err != nil {
		t.Fatal(err)
	}
	base := filepath.Join(ticket.Root, ".parley-runtime", "trajectory-verification", ticket.RunID)
	if err = os.MkdirAll(base, 0700); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(base, "request.json"), reqData, 0600); err != nil {
		t.Fatal(err)
	}
	terminalRaw := snapshotRead(t, filepath.Join(ticket.Root, ".parley-runtime", "invocations", refused, "terminal.json"))
	original := ParentResult{Version: 1, RunID: ticket.RunID, RequestPath: filepath.Join(base, "request.json"),
		RequestSHA256: digest(reqData), InvocationID: refused, TerminalSHA256: digest(terminalRaw),
		FailureStage: "launch", TrajectoryPending: true}
	originalRaw, err := canonical(original)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(base, "parent-result.json"), originalRaw, 0600); err != nil {
		t.Fatal(err)
	}
	retainedReplacementLifecycle(t, ticket, replacement, nil, recovery.At, func() {
		receipt, execErr := ExecuteCapturedVerification(ctx, ticket, replacement, criteria, t.TempDir())
		if execErr != nil || receipt.Steps != 4 || receipt.FailureStage != "" || receipt.InvocationID != replacement {
			t.Fatalf("recovered helper execution did not complete: %+v %v", receipt, execErr)
		}
	})
	return refusedParentRecoveryLineage{ticket: ticket, path: path, base: base, refused: refused, replacement: replacement, original: original, originalRaw: originalRaw}
}

// TestRecoveredParentObservationResolvesAndReconciles: with the original failed
// parent retained and the replacement lineage fully observed, the immutable
// recovered observation publishes exactly once, fresh reconciliation resolves
// with its provenance pinned, the explicit resolution replays exactly, and the
// old failure plus all prior records stay byte-identical throughout.
func TestRecoveredParentObservationResolvesAndReconciles(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX captured verification")
	}
	ctx := context.Background()
	f := refusedParentRecoveryFixture(t)

	// Before publication the retained failure blocks resolution, and the
	// missing-parent recovery stays a separate mechanism that still refuses
	// the complete failed parent.
	if _, err := PreviewReconciliation(ctx, f.ticket.Root, "fixture", f.ticket.RunID); err == nil ||
		!strings.Contains(err.Error(), "parent did not retain the independently derived successful comparison") {
		t.Fatalf("retained failure resolved without the recovered observation: %v", err)
	}
	if _, err := PreviewParentRecovery(ctx, f.ticket.Root, "fixture", f.ticket.RunID); err == nil ||
		!strings.Contains(err.Error(), "conflicting identity, failure or assessment") {
		t.Fatalf("missing-parent recovery conflated the retained failure: %v", err)
	}

	p, err := PreviewRecoveredParent(ctx, f.ticket.Root, "fixture", f.ticket.RunID)
	if err != nil {
		t.Fatalf("recovered preview refused the complete lineage: %v", err)
	}
	if p.Version != 1 || p.Root != f.ticket.Root || p.Idea != "fixture" || p.RunID != f.ticket.RunID ||
		p.Original != (ParentObservation{Kind: "failed-budget-refusal", SHA256: digest(f.originalRaw)}) ||
		!sameJSON(p.OriginalResult, f.original) || p.Recovery.InvocationID != f.replacement ||
		p.Recovery.RefusedInvocationID != f.refused || !validHash(p.Recovery.RecoverySHA256) ||
		p.Derived.Result.InvocationID != f.replacement || p.Derived.Result.Assessment == nil ||
		p.Derived.Result.Assessment.Outcome != Regression || !validHash(p.ParentSHA256) {
		t.Fatalf("recovered preview is not the complete lineage binding: %+v", p)
	}
	if _, statErr := os.Stat(filepath.Join(f.base, recoveredParentName)); !os.IsNotExist(statErr) {
		t.Fatal("preview wrote the observation")
	}
	if _, err = PublishRecoveredParent(ctx, f.ticket.Root, "fixture", f.ticket.RunID, digest([]byte("stale"))); err == nil ||
		!strings.Contains(err.Error(), "changed since preview") {
		t.Fatalf("stale preview was published: %v", err)
	}
	published, err := PublishRecoveredParent(ctx, f.ticket.Root, "fixture", f.ticket.RunID, p.SHA256())
	if err != nil || published.SHA256() != p.SHA256() {
		t.Fatalf("exact publish did not reproduce the inspected preview: %+v %v", published, err)
	}
	if !bytes.Equal(f.originalRaw, snapshotRead(t, filepath.Join(f.base, "parent-result.json"))) {
		t.Fatal("publication overwrote the original failed parent result")
	}
	replayed, err := PublishRecoveredParent(ctx, f.ticket.Root, "fixture", f.ticket.RunID, p.SHA256())
	if err != nil || replayed.SHA256() != p.SHA256() {
		t.Fatalf("exact replay changed the observation: %+v %v", replayed, err)
	}
	if _, err = PublishRecoveredParent(ctx, f.ticket.Root, "fixture", f.ticket.RunID, digest([]byte("forged"))); err == nil ||
		!strings.Contains(err.Error(), "replay changed its original preview") {
		t.Fatalf("conflicting replay was admitted: %v", err)
	}

	preview, err := PreviewReconciliation(ctx, f.ticket.Root, "fixture", f.ticket.RunID)
	if err != nil {
		t.Fatalf("fresh reconciliation refused the recovered lineage: %v", err)
	}
	if preview.RecoverySHA256 == "" || preview.Assessment.Outcome != Regression || preview.ParentSHA256 != p.ParentSHA256 {
		t.Fatalf("recovered preview lost its provenance, outcome or parent binding: %+v", preview)
	}
	expected := preview.SHA256()
	if err = Reconcile(ctx, f.ticket.Root, "fixture", f.ticket.RunID, expected); err != nil {
		t.Fatalf("explicit resolution refused the recovered lineage: %v", err)
	}
	if err = Reconcile(ctx, f.ticket.Root, "fixture", f.ticket.RunID, expected); err != nil {
		t.Fatalf("exact resolution replay changed its decision: %v", err)
	}
	state, err := Inspect(ctx, f.ticket.Root, "fixture")
	if err != nil || len(state.Resolutions) != 1 || state.Resolutions[0].SHA256 != expected ||
		state.Resolutions[0].Preview.RecoverySHA256 != preview.RecoverySHA256 {
		t.Fatalf("recovered resolution was not retained with its provenance: %v %+v", err, state)
	}
	if !bytes.Equal(f.originalRaw, snapshotRead(t, filepath.Join(f.base, "parent-result.json"))) {
		t.Fatal("reconciliation overwrote the original failed parent result")
	}
	if err = RequireResolved(ctx, f.ticket.Root, "fixture"); err == nil {
		t.Fatal("regression resolution closed the trajectory without a clean outcome")
	}
}

// TestRecoveredParentObservationRequiresExactLineage: every post-publication
// mutation or deletion of the original failure, the replacement lineage, the
// journal recovery or the record itself fails the fresh validation and the
// later state guards, and never auto-resolves.
func TestRecoveredParentObservationRequiresExactLineage(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX captured verification")
	}
	ctx := context.Background()
	for _, tc := range []struct {
		name    string
		mutate  func(t *testing.T, f refusedParentRecoveryLineage)
		message string
	}{
		{name: "original-result-deleted", mutate: func(t *testing.T, f refusedParentRecoveryLineage) {
			removePath(t, filepath.Join(f.base, "parent-result.json"))
		}, message: "lost its original refused parent result"},
		{name: "original-result-rewritten", mutate: func(t *testing.T, f refusedParentRecoveryLineage) {
			forged := f.original
			forged.TerminalSHA256 = digest([]byte("forged"))
			data, err := canonical(forged)
			if err != nil {
				t.Fatal(err)
			}
			if err = os.WriteFile(filepath.Join(f.base, "parent-result.json"), data, 0600); err != nil {
				t.Fatal(err)
			}
		}, message: "original refused parent result changed after recovery"},
		{name: "original-result-success-stage", mutate: func(t *testing.T, f refusedParentRecoveryLineage) {
			forged := f.original
			forged.FailureStage = ""
			data, err := canonical(forged)
			if err != nil {
				t.Fatal(err)
			}
			if err = os.WriteFile(filepath.Join(f.base, "parent-result.json"), data, 0600); err != nil {
				t.Fatal(err)
			}
		}, message: "original refused parent result changed after recovery"},
		{name: "replacement-terminal-deleted", mutate: func(t *testing.T, f refusedParentRecoveryLineage) {
			removePath(t, filepath.Join(f.ticket.Root, ".parley-runtime", "invocations", f.replacement, "terminal.json"))
		}, message: "successful observed parent terminal binding is unavailable"},
		{name: "replacement-terminal-rewritten", mutate: func(t *testing.T, f refusedParentRecoveryLineage) {
			path := filepath.Join(f.ticket.Root, ".parley-runtime", "invocations", f.replacement, "terminal.json")
			var record telemetry.Record
			if err := json.Unmarshal(snapshotRead(t, path), &record); err != nil {
				t.Fatal(err)
			}
			record.Metadata.AttemptOrdinal = 99
			data, err := canonical(record)
			if err != nil {
				t.Fatal(err)
			}
			if err = os.WriteFile(path, data, 0600); err != nil {
				t.Fatal(err)
			}
		}, message: "original requested lifecycle differs from the verifier terminal"},
		{name: "observation-record-rewritten", mutate: func(t *testing.T, f refusedParentRecoveryLineage) {
			var record RecoveredParentRecord
			if err := json.Unmarshal(snapshotRead(t, filepath.Join(f.base, recoveredParentName)), &record); err != nil {
				t.Fatal(err)
			}
			record.Preview.OriginalResult.InvocationID = "decoy-verifier-invocation"
			data, err := canonical(record)
			if err != nil {
				t.Fatal(err)
			}
			if err = os.WriteFile(filepath.Join(f.base, recoveredParentName), data, 0600); err != nil {
				t.Fatal(err)
			}
		}, message: "invalid recovered parent observation"},
		{name: "recovery-deleted", mutate: func(t *testing.T, f refusedParentRecoveryLineage) {
			removePath(t, filepath.Join(f.path, "recovery.json"))
		}, message: "successful observed parent terminal binding is unavailable"},
		{name: "refused-evidence-rewritten", mutate: func(t *testing.T, f refusedParentRecoveryLineage) {
			rewriteRefusedEvidenceCanonical(t, f.ticket.Root, f.refused)
		}, message: "no longer matches its retained refused-launch evidence"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := refusedParentRecoveryFixture(t)
			p, err := PreviewRecoveredParent(ctx, f.ticket.Root, "fixture", f.ticket.RunID)
			if err != nil {
				t.Fatal(err)
			}
			if _, err = PublishRecoveredParent(ctx, f.ticket.Root, "fixture", f.ticket.RunID, p.SHA256()); err != nil {
				t.Fatal(err)
			}
			tc.mutate(t, f)
			if _, err = PreviewReconciliation(ctx, f.ticket.Root, "fixture", f.ticket.RunID); err == nil ||
				!strings.Contains(err.Error(), tc.message) {
				t.Fatalf("%s did not fail the fresh validation: %v", tc.name, err)
			}
			if err = RequireResolved(ctx, f.ticket.Root, "fixture"); err == nil {
				t.Fatalf("%s auto-resolved the trajectory", tc.name)
			}
			state, inspectErr := Inspect(ctx, f.ticket.Root, "fixture")
			if inspectErr != nil || len(state.Resolutions) != 0 {
				t.Fatalf("%s recorded a resolution: %v %+v", tc.name, inspectErr, state)
			}
		})
	}
}

// TestRecoveredParentResolutionGuardsRevalidate: a resolution accepted over the
// recovered observation is re-derived by every later state guard; mutating the
// original failure, the replacement lineage or the observation afterwards
// breaks the guard instead of passing on the recorded facts.
func TestRecoveredParentResolutionGuardsRevalidate(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX captured verification")
	}
	ctx := context.Background()
	for _, tc := range []struct {
		name    string
		mutate  func(t *testing.T, f refusedParentRecoveryLineage)
		message string
	}{
		{name: "original-result-deleted", mutate: func(t *testing.T, f refusedParentRecoveryLineage) {
			removePath(t, filepath.Join(f.base, "parent-result.json"))
		}, message: "lost its original refused parent result"},
		{name: "replacement-terminal-deleted", mutate: func(t *testing.T, f refusedParentRecoveryLineage) {
			removePath(t, filepath.Join(f.ticket.Root, ".parley-runtime", "invocations", f.replacement, "terminal.json"))
		}, message: "successful observed parent terminal binding is unavailable"},
		{name: "observation-deleted", mutate: func(t *testing.T, f refusedParentRecoveryLineage) {
			removePath(t, filepath.Join(f.base, recoveredParentName))
		}, message: "parent did not retain the independently derived successful comparison"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := refusedParentRecoveryFixture(t)
			p, err := PreviewRecoveredParent(ctx, f.ticket.Root, "fixture", f.ticket.RunID)
			if err != nil {
				t.Fatal(err)
			}
			if _, err = PublishRecoveredParent(ctx, f.ticket.Root, "fixture", f.ticket.RunID, p.SHA256()); err != nil {
				t.Fatal(err)
			}
			preview, err := PreviewReconciliation(ctx, f.ticket.Root, "fixture", f.ticket.RunID)
			if err != nil {
				t.Fatal(err)
			}
			if err = Reconcile(ctx, f.ticket.Root, "fixture", f.ticket.RunID, preview.SHA256()); err != nil {
				t.Fatal(err)
			}
			tc.mutate(t, f)
			if _, err = Inspect(ctx, f.ticket.Root, "fixture"); err == nil || !strings.Contains(err.Error(), tc.message) {
				t.Fatalf("%s did not fail the later state guard: %v", tc.name, err)
			}
			if _, err = PreviewReconciliation(ctx, f.ticket.Root, "fixture", f.ticket.RunID); err == nil ||
				!strings.Contains(err.Error(), tc.message) {
				t.Fatalf("%s did not fail the fresh reconciliation: %v", tc.name, err)
			}
		})
	}
}

// TestRecoveredParentPublicationRefusesNonRefusalOriginal: a successful
// (non-refusal) retained parent result is never conflated with this mechanism —
// the observation preview refuses it, while the ordinary same-bytes parent
// resolution path stays available and byte-identical.
func TestRecoveredParentPublicationRefusesNonRefusalOriginal(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX captured verification")
	}
	ctx := context.Background()
	f := recoveredParentLineageFixture(t, nil)
	if _, err := PreviewRecoveredParent(ctx, f.ticket.Root, "fixture", f.ticket.RunID); err == nil ||
		!strings.Contains(err.Error(), "original retained result is not the exact pre-start budget refusal") {
		t.Fatalf("non-refusal original previewed as a recovered observation: %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(f.ticket.Root, ".parley-runtime", "trajectory-verification", f.ticket.RunID, recoveredParentName)); !os.IsNotExist(statErr) {
		t.Fatal("refused preview wrote an observation")
	}
	// The ordinary resolution path is unaffected and still resolves.
	preview, err := PreviewReconciliation(ctx, f.ticket.Root, "fixture", f.ticket.RunID)
	if err != nil || preview.RecoverySHA256 != "" {
		t.Fatalf("ordinary recovered lineage lost its plain resolution path: %+v %v", preview, err)
	}
}
