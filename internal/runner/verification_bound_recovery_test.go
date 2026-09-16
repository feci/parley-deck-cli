package runner

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/telemetry"
	"parley-deck-cli/internal/trajectory"
)

// The bound verifier relaunch. A real budget-refused verifier launch is
// recovered explicitly; a trusted recovery context then lets one real
// RunMeasured relaunch execute under the invocation ID the immutable recovery
// already pinned — telemetry exclusively allocates that ID, the reservation
// boundary revalidates the whole recovered authority, and the verifier child's
// captured helper executes once inside the observed launch window. All agents
// are local /bin/sh fixtures; no provider or model is invoked.

const boundVerifierInvocation = "recovered-verifier-invocation"

type boundVerifierFixture struct {
	root        string
	idea        protocol.IdeaStatus
	criteria    []trajectory.Criterion
	ticket      trajectory.VerificationTicket
	journal     string
	binding     *budget.CycleBinding
	cycleLedger []byte
	launchRaw   []byte
	refusedID   string
	recoveryRaw []byte
}

// boundVerifierRefusalFixture drives the real public launch path into its
// pre-start budget refusal (no verifier process), and by default repairs it
// with the explicit public recovery pinning boundVerifierInvocation.
func boundVerifierRefusalFixture(t *testing.T, recover bool) *boundVerifierFixture {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("POSIX captured verification")
	}
	ctx := context.Background()
	root, idea, criteria := recoveredVerifierRuntimeFixture(t)
	// Match the canonical origin persisted in verification tickets.
	root, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	builder := telemetryShell("printf 'broken\\n' > source; exit 7", false)
	if attempt, err := RunMeasured(ctx, ExecOptions{Root: root, Agent: builder, Prompt: "synthetic patch", Timeout: 30 * time.Second, Info: LaunchInfo{RunID: "patch-attempt", Idea: idea.Slug, Phase: "fixup"}}); err == nil || attempt.StartedAt == nil {
		t.Fatalf("changed-source attempt did not actually execute: %+v %v", attempt, err)
	}
	ticket, err := trajectory.PrepareCapturedVerification(ctx, root, idea.Slug, "reviewer", "verifier-run")
	if err != nil {
		t.Fatal(err)
	}
	b, err := budget.LoadCycleBinding(ctx, root, idea.Slug, budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	f := &boundVerifierFixture{root: root, idea: idea, criteria: criteria, ticket: ticket, binding: b,
		journal: filepath.Join(filepath.Dir(b.Store.Dir), "trajectory-verifications", ticket.Request.Charge.EntryKey)}
	f.cycleLedger = readFixtureBytes(t, filepath.Join(b.Store.Dir, "ledger.json"))
	bound, err := WithCapturedVerification(ctx, ticket)
	if err != nil {
		t.Fatal(err)
	}
	refusedStore := budget.Store{Dir: filepath.Join(t.TempDir(), "launch"), Scope: "verifier-launch-refusal"}
	denied := WithLaunchBudget(bound, LaunchBudget{Store: refusedStore, Limits: budget.Limits{Denied: map[budget.Kind]bool{budget.Launch: true}}})
	verifier := telemetryShell("touch .parley-runtime/verifier-spawned", false)
	verifier.ID = "reviewer"
	refused, err := RunMeasured(denied, ExecOptions{Root: root, Agent: verifier, Prompt: "synthetic verifier", Timeout: 30 * time.Second,
		Info: LaunchInfo{RunID: ticket.RunID, Idea: idea.Slug, Phase: CapturedVerificationPhase}})
	if err == nil || refused.StartedAt != nil || refused.Outcome == nil || refused.Outcome.FailureClass == nil ||
		*refused.Outcome.FailureClass != "budget_refused" || refused.InvocationID == "" {
		t.Fatalf("verifier launch was not refused before start: %+v %v", refused, err)
	}
	if _, err = os.Stat(filepath.Join(root, ".parley-runtime", "verifier-spawned")); !os.IsNotExist(err) {
		t.Fatal("refused verifier spawned a child process")
	}
	f.refusedID = refused.InvocationID
	f.launchRaw = readFixtureBytes(t, filepath.Join(f.journal, "launch.json"))
	if recover {
		if err = trajectory.RecoverCapturedVerificationLaunch(ctx, ticket, f.refusedID, boundVerifierInvocation); err != nil {
			t.Fatal(err)
		}
		f.recoveryRaw = readFixtureBytes(t, filepath.Join(f.journal, "recovery.json"))
	}
	return f
}

func readFixtureBytes(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func fixtureDigest(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// canonicalFixtureJSON matches the exact canonical encoding every private
// runtime artifact reader requires (indented JSON plus one trailing newline).
func canonicalFixtureJSON(t *testing.T, v any) []byte {
	t.Helper()
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	return append(data, '\n')
}

func decodeFixtureJSON(t *testing.T, path string, v any) {
	t.Helper()
	if err := json.Unmarshal(readFixtureBytes(t, path), v); err != nil {
		t.Fatal(err)
	}
}

// rewriteFixtureRecovery rewrites the retained recovery artifact canonically,
// changing only the bound replacement invocation. Field order mirrors the
// retained artifact so the rewrite stays byte-canonical for the authority
// reader; this is the same-UID honesty boundary, used here to prove the
// launch boundary does not trust the admission-time validation alone.
func rewriteFixtureRecovery(t *testing.T, f *boundVerifierFixture, invocation string) {
	t.Helper()
	var admitted struct {
		Version                int       `json:"version"`
		TicketSHA256           string    `json:"ticket_sha256"`
		PriorLaunchSHA256      string    `json:"prior_launch_sha256"`
		RefusedInvocationID    string    `json:"refused_invocation_id"`
		RefusedRequestedSHA256 string    `json:"refused_requested_sha256"`
		RefusedTerminalSHA256  string    `json:"refused_terminal_sha256"`
		InvocationID           string    `json:"invocation_id"`
		At                     time.Time `json:"at"`
	}
	if err := json.Unmarshal(f.recoveryRaw, &admitted); err != nil {
		t.Fatal(err)
	}
	admitted.InvocationID = invocation
	if err := os.WriteFile(filepath.Join(f.journal, "recovery.json"), canonicalFixtureJSON(t, admitted), 0o600); err != nil {
		t.Fatal(err)
	}
}

// relaunchBoundVerifier runs one real RunMeasured verifier relaunch under the
// pinned invocation with a budget-authorized store. The verifier child does
// its work once and then waits for the captured helper's receipt; the helper
// executes concurrently, triggered by the pinned invocation's own started
// record, so its receipt lands inside the observed launch window exactly as
// an in-process helper would produce it.
func relaunchBoundVerifier(t *testing.T, f *boundVerifierFixture, bound context.Context) telemetry.Record {
	t.Helper()
	retryStore := budget.Store{Dir: filepath.Join(t.TempDir(), "launch-relaunch"), Scope: "verifier-relaunch"}
	relaunch := WithLaunchBudget(bound, LaunchBudget{Store: retryStore})
	quote := "'" + strings.ReplaceAll(filepath.Join(f.journal, "receipt.json"), "'", "'\\''") + "'"
	script := "printf 'ran\\n' >> .parley-runtime/verifier-relaunched; i=0; while [ ! -f " + quote + " ] && [ \"$i\" -lt 600 ]; do sleep 0.05; i=$((i+1)); done; test -f " + quote
	verifier := telemetryShell(script, false)
	verifier.ID = "reviewer"
	parent := t.TempDir()
	type helperResult struct {
		receipt trajectory.VerificationReceipt
		err     error
	}
	helper := make(chan helperResult, 1)
	go func() {
		started := filepath.Join(f.root, ".parley-runtime", "invocations", boundVerifierInvocation, "started.json")
		deadline := time.Now().Add(25 * time.Second)
		for {
			_, err := os.Stat(started)
			if err == nil {
				break
			}
			if !os.IsNotExist(err) {
				helper <- helperResult{err: err}
				return
			}
			if time.Now().After(deadline) {
				helper <- helperResult{err: errors.New("bound verifier relaunch never started")}
				return
			}
			time.Sleep(20 * time.Millisecond)
		}
		receipt, err := trajectory.ExecuteCapturedVerification(context.Background(), f.ticket, boundVerifierInvocation, f.criteria, parent)
		helper <- helperResult{receipt: receipt, err: err}
	}()
	record, err := RunMeasured(relaunch, ExecOptions{Root: f.root, Agent: verifier, Prompt: "synthetic verifier relaunch", Timeout: 60 * time.Second,
		Info: LaunchInfo{RunID: f.ticket.RunID, Idea: f.idea.Slug, Phase: CapturedVerificationPhase}})
	done := <-helper
	if err != nil {
		t.Fatalf("bound verifier relaunch failed: %v (helper: %v)", err, done.err)
	}
	if done.err != nil || done.receipt.Steps != 4 || done.receipt.FailureStage != "" || done.receipt.InvocationID != boundVerifierInvocation {
		t.Fatalf("captured helper did not complete inside the bound relaunch: %+v %v", done.receipt, done.err)
	}
	if record.InvocationID != boundVerifierInvocation || record.StartedAt == nil || record.Outcome == nil ||
		record.Outcome.Status != "process-exited" || record.Outcome.ExitCode == nil || *record.Outcome.ExitCode != 0 {
		t.Fatalf("bound relaunch did not execute under the pinned invocation: %+v", record)
	}
	return record
}

// TestBoundVerifierRecoveryRelaunchExecutesPinnedInvocationAndResolves: the
// real refusal, the explicit recovery, and one real budget-authorized
// RunMeasured relaunch under the pinned ID. The pinned invocation's lifecycle
// is native runner output (requested/started/terminal), the helper executes
// once inside the observed window, retained lineage and accounting stay
// byte-identical, and a fresh parent resolution over the native records
// previews, reconciles and replays idempotently.
func TestBoundVerifierRecoveryRelaunchExecutesPinnedInvocationAndResolves(t *testing.T) {
	ctx := context.Background()
	f := boundVerifierRefusalFixture(t, true)
	bound, err := WithCapturedVerificationRecovery(ctx, f.ticket, boundVerifierInvocation)
	if err != nil {
		t.Fatalf("trusted recovery context refused a valid retained recovery: %v", err)
	}
	relaunchBoundVerifier(t, f, bound)

	invBase := filepath.Join(f.root, ".parley-runtime", "invocations", boundVerifierInvocation)
	var requested, started, terminal telemetry.Record
	decodeFixtureJSON(t, filepath.Join(invBase, "requested.json"), &requested)
	decodeFixtureJSON(t, filepath.Join(invBase, "started.json"), &started)
	decodeFixtureJSON(t, filepath.Join(invBase, "terminal.json"), &terminal)
	var recovery struct {
		At time.Time `json:"at"`
	}
	if err = json.Unmarshal(f.recoveryRaw, &recovery); err != nil {
		t.Fatal(err)
	}
	if requested.RequestedAt.Before(recovery.At) {
		t.Fatal("bound relaunch was requested before the recovery that pins it")
	}
	if terminal.Metadata.RunID != f.ticket.RunID || terminal.Metadata.Idea != f.idea.Slug || terminal.Metadata.Phase != CapturedVerificationPhase ||
		terminal.Metadata.Agent != "reviewer" || terminal.Metadata.LaunchMode != "headless" ||
		terminal.StartedAt == nil || terminal.PID == nil || *terminal.PID <= 0 || terminal.CompletedAt == nil ||
		terminal.Outcome == nil || terminal.Outcome.Status != "process-exited" || terminal.Outcome.ExitCode == nil || *terminal.Outcome.ExitCode != 0 {
		t.Fatalf("native bound terminal is not the successful observed verifier: %+v", terminal)
	}
	if started.StartedAt == nil || !started.StartedAt.Equal(*terminal.StartedAt) || started.PID == nil || *started.PID != *terminal.PID {
		t.Fatalf("native bound started record disagrees with the terminal: %+v", started)
	}
	if marker := readFixtureBytes(t, filepath.Join(f.root, ".parley-runtime", "verifier-relaunched")); string(marker) != "ran\n" {
		t.Fatalf("verifier work did not execute exactly once: %q", marker)
	}
	if _, err = os.Stat(filepath.Join(f.root, ".parley-runtime", "verifier-spawned")); !os.IsNotExist(err) {
		t.Fatal("relaunch reused the refused verifier's process marker")
	}
	if after := readFixtureBytes(t, filepath.Join(f.journal, "launch.json")); !bytes.Equal(f.launchRaw, after) {
		t.Fatal("relaunch rewrote the original launch reservation")
	}
	if after := readFixtureBytes(t, filepath.Join(f.journal, "recovery.json")); !bytes.Equal(f.recoveryRaw, after) {
		t.Fatal("relaunch rewrote the immutable recovery")
	}
	if after := readFixtureBytes(t, filepath.Join(f.binding.Store.Dir, "ledger.json")); !bytes.Equal(f.cycleLedger, after) {
		t.Fatal("relaunch changed the original fixup accounting")
	}
	refusedBase := filepath.Join(f.root, ".parley-runtime", "invocations", f.refusedID)
	for _, name := range []string{"requested.json", "terminal.json"} {
		if _, statErr := os.Stat(filepath.Join(refusedBase, name)); statErr != nil {
			t.Fatal("refused launch evidence not retained", statErr)
		}
	}
	if _, statErr := os.Stat(filepath.Join(refusedBase, "started.json")); !os.IsNotExist(statErr) {
		t.Fatal("refused invocation fabricated a started lifecycle")
	}
	read, observation, err := trajectory.ReadCapturedVerification(ctx, f.ticket, boundVerifierInvocation)
	if err != nil {
		t.Fatal(err)
	}
	assessment, err := trajectory.AssessCaptured(f.ticket.Request, observation)
	if err != nil || assessment.Outcome != trajectory.NoRegression {
		t.Fatalf("bound observation lost its material assessment: %+v %v", assessment, err)
	}

	// Fresh parent resolution over the NATIVE lifecycle records: publish the
	// app-side verifier-run request and parent result exactly as the app
	// retains them, then preview, reconcile and replay the resolution.
	scope := "---\nparticipants: [test-1, reviewer]\nchecks:\n  - name: material\n    command: >\n      " + f.criteria[0].Command + "\n---\n"
	mustWrite(t, filepath.Join(f.idea.Path, "00-prompt.md"), scope)
	base := filepath.Join(f.root, ".parley-runtime", "trajectory-verification", f.ticket.RunID)
	if err = os.MkdirAll(base, 0o700); err != nil {
		t.Fatal(err)
	}
	reqData := canonicalFixtureJSON(t, trajectory.HelperRequest{Version: 1, Ticket: f.ticket, Participants: []string{"test-1", "reviewer"}, Criteria: f.criteria})
	if err = os.WriteFile(filepath.Join(base, "request.json"), reqData, 0o600); err != nil {
		t.Fatal(err)
	}
	result := trajectory.ParentResult{Version: 1, RunID: f.ticket.RunID, RequestPath: filepath.Join(base, "request.json"),
		RequestSHA256: fixtureDigest(reqData), InvocationID: boundVerifierInvocation,
		TerminalSHA256: fixtureDigest(readFixtureBytes(t, filepath.Join(invBase, "terminal.json"))),
		ReceiptSHA256:  fixtureDigest(canonicalFixtureJSON(t, read)), Assessment: &assessment, TrajectoryPending: true}
	if err = os.WriteFile(filepath.Join(base, "parent-result.json"), canonicalFixtureJSON(t, result), 0o600); err != nil {
		t.Fatal(err)
	}
	preview, err := trajectory.PreviewReconciliation(ctx, f.root, f.idea.Slug, f.ticket.RunID)
	if err != nil {
		t.Fatalf("native recovered lineage did not resolve: %v", err)
	}
	if preview.Version != 1 || preview.Sequence != 1 || preview.Root != f.root || preview.RunID != f.ticket.RunID || preview.Assessment.Outcome != trajectory.NoRegression {
		t.Fatalf("native recovered preview is not the derived parent binding: %+v", preview)
	}
	expected := preview.SHA256()
	if err = trajectory.Reconcile(ctx, f.root, f.idea.Slug, f.ticket.RunID, expected); err != nil {
		t.Fatalf("explicit resolution refused the native recovered lineage: %v", err)
	}
	if err = trajectory.Reconcile(ctx, f.root, f.idea.Slug, f.ticket.RunID, expected); err != nil {
		t.Fatalf("exact resolution replay changed its decision: %v", err)
	}
	state, err := trajectory.Inspect(ctx, f.root, f.idea.Slug)
	if err != nil || len(state.Resolutions) != 1 || state.Resolutions[0].SHA256 != expected {
		t.Fatalf("recovered resolution was not retained: %v %+v", err, state)
	}
}

// TestBoundVerifierRecoveryContextNeverWrites: admission of the trusted
// recovery context is read-only. A ticket whose launch was never refused must
// be refused WITHOUT the reservation probe creating the missing launch
// reservation or any bound invocation directory.
func TestBoundVerifierRecoveryContextNeverWrites(t *testing.T) {
	ctx := context.Background()
	root, idea, _ := recoveredVerifierRuntimeFixture(t)
	builder := telemetryShell("printf 'broken\\n' > source; exit 7", false)
	if attempt, err := RunMeasured(ctx, ExecOptions{Root: root, Agent: builder, Prompt: "synthetic patch", Timeout: 30 * time.Second, Info: LaunchInfo{RunID: "patch-attempt", Idea: idea.Slug, Phase: "fixup"}}); err == nil || attempt.StartedAt == nil {
		t.Fatalf("changed-source attempt did not actually execute: %+v %v", attempt, err)
	}
	ticket, err := trajectory.PrepareCapturedVerification(ctx, root, idea.Slug, "reviewer", "verifier-run")
	if err != nil {
		t.Fatal(err)
	}
	b, err := budget.LoadCycleBinding(ctx, root, idea.Slug, budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	journal := filepath.Join(filepath.Dir(b.Store.Dir), "trajectory-verifications", ticket.Request.Charge.EntryKey)
	if _, err = WithCapturedVerificationRecovery(ctx, ticket, boundVerifierInvocation); err == nil ||
		!strings.Contains(err.Error(), "requires a retained launch reservation and immutable recovery record") {
		t.Fatalf("recovery context admitted a launch that was never refused: %v", err)
	}
	if _, err = os.Stat(filepath.Join(journal, "launch.json")); !os.IsNotExist(err) {
		t.Fatal("admission probe wrote a launch reservation")
	}
	if _, err = os.Stat(filepath.Join(journal, "recovery.json")); !os.IsNotExist(err) {
		t.Fatal("admission probe wrote a recovery artifact")
	}
	if _, err = os.Stat(filepath.Join(root, ".parley-runtime", "invocations", boundVerifierInvocation)); !os.IsNotExist(err) {
		t.Fatal("admission probe allocated the bound invocation directory")
	}
}

// TestBoundVerifierRecoveryContextRefusesEveryContradiction: the trusted
// context is admitted only for the exact retained authority — a stale refused
// invocation, a decoy invocation, an unsafe identifier, a recovery rebound to
// a decoy, or mutated refused evidence each fails closed with the journal
// byte-unchanged and no bound invocation directory allocated.
func TestBoundVerifierRecoveryContextRefusesEveryContradiction(t *testing.T) {
	ctx := context.Background()
	for _, tc := range []struct {
		name            string
		build           func(t *testing.T) (*boundVerifierFixture, string)
		message         string
		mutatedRecovery bool
	}{
		{name: "stale-refused-invocation", build: func(t *testing.T) (*boundVerifierFixture, string) {
			f := boundVerifierRefusalFixture(t, true)
			return f, f.refusedID
		}, message: "superseded by an explicit recovery"},
		{name: "mismatched-invocation", build: func(t *testing.T) (*boundVerifierFixture, string) {
			f := boundVerifierRefusalFixture(t, true)
			return f, "decoy-verifier-invocation"
		}, message: "differs from its durable reservation"},
		{name: "unsafe-invocation-id", build: func(t *testing.T) (*boundVerifierFixture, string) {
			f := boundVerifierRefusalFixture(t, true)
			return f, "escape/attempt"
		}, message: "safe single-element invocation identifier"},
		{name: "recovery-rebound-to-decoy", build: func(t *testing.T) (*boundVerifierFixture, string) {
			f := boundVerifierRefusalFixture(t, true)
			rewriteFixtureRecovery(t, f, "decoy-verifier-invocation")
			return f, boundVerifierInvocation
		}, message: "differs from its durable reservation", mutatedRecovery: true},
		{name: "refused-evidence-rewritten", build: func(t *testing.T) (*boundVerifierFixture, string) {
			f := boundVerifierRefusalFixture(t, true)
			invDir := filepath.Join(f.root, ".parley-runtime", "invocations", f.refusedID)
			for _, name := range []string{"requested.json", "terminal.json"} {
				var record telemetry.Record
				decodeFixtureJSON(t, filepath.Join(invDir, name), &record)
				record.Metadata.AttemptOrdinal = 99
				if err := os.WriteFile(filepath.Join(invDir, name), canonicalFixtureJSON(t, record), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			return f, boundVerifierInvocation
		}, message: "no longer matches its retained refused-launch evidence"},
		{name: "missing-recovery", build: func(t *testing.T) (*boundVerifierFixture, string) {
			f := boundVerifierRefusalFixture(t, false)
			return f, boundVerifierInvocation
		}, message: "requires a retained launch reservation and immutable recovery record"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, intended := tc.build(t)
			if _, err := WithCapturedVerificationRecovery(ctx, f.ticket, intended); err == nil ||
				!strings.Contains(err.Error(), tc.message) {
				t.Fatalf("%s was admitted: %v", tc.name, err)
			}
			if after := readFixtureBytes(t, filepath.Join(f.journal, "launch.json")); !bytes.Equal(f.launchRaw, after) {
				t.Fatalf("%s rewrote the launch reservation", tc.name)
			}
			switch {
			case tc.mutatedRecovery:
				if _, statErr := os.Stat(filepath.Join(f.journal, "recovery.json")); statErr != nil {
					t.Fatalf("%s lost the mutated recovery artifact: %v", tc.name, statErr)
				}
			case f.recoveryRaw != nil:
				if after := readFixtureBytes(t, filepath.Join(f.journal, "recovery.json")); !bytes.Equal(f.recoveryRaw, after) {
					t.Fatalf("%s rewrote the retained recovery artifact", tc.name)
				}
			default:
				if _, statErr := os.Stat(filepath.Join(f.journal, "recovery.json")); !os.IsNotExist(statErr) {
					t.Fatalf("%s wrote a recovery artifact", tc.name)
				}
			}
			if intended == boundVerifierInvocation {
				if _, err := os.Stat(filepath.Join(f.root, ".parley-runtime", "invocations", boundVerifierInvocation)); !os.IsNotExist(err) {
					t.Fatalf("%s allocated the bound invocation directory", tc.name)
				}
			}
		})
	}
}

// TestBoundVerifierRelaunchFailsClosedWhenAuthorityChanged: a context admitted
// while the authority was valid must still fail closed at the actual launch if
// the recovery changed in between — the reservation boundary revalidates the
// full recovered lineage, so no process spawns, no budget is charged, and the
// pinned invocation directory is consumed by the failed attempt rather than
// silently reused.
func TestBoundVerifierRelaunchFailsClosedWhenAuthorityChanged(t *testing.T) {
	ctx := context.Background()
	f := boundVerifierRefusalFixture(t, true)
	bound, err := WithCapturedVerificationRecovery(ctx, f.ticket, boundVerifierInvocation)
	if err != nil {
		t.Fatal(err)
	}
	rewriteFixtureRecovery(t, f, "decoy-verifier-invocation")
	changedRecovery := readFixtureBytes(t, filepath.Join(f.journal, "recovery.json"))
	retryStore := budget.Store{Dir: filepath.Join(t.TempDir(), "launch-relaunch"), Scope: "verifier-relaunch"}
	relaunch := WithLaunchBudget(bound, LaunchBudget{Store: retryStore})
	verifier := telemetryShell("printf 'ran\\n' >> .parley-runtime/verifier-relaunched", false)
	verifier.ID = "reviewer"
	record, err := RunMeasured(relaunch, ExecOptions{Root: f.root, Agent: verifier, Prompt: "synthetic verifier relaunch", Timeout: 30 * time.Second,
		Info: LaunchInfo{RunID: f.ticket.RunID, Idea: f.idea.Slug, Phase: CapturedVerificationPhase}})
	if err == nil || record.StartedAt != nil || !strings.Contains(err.Error(), "differs from its durable reservation") {
		t.Fatalf("changed recovery authority still executed the bound relaunch: %+v %v", record, err)
	}
	if _, statErr := os.Stat(filepath.Join(f.root, ".parley-runtime", "verifier-relaunched")); !os.IsNotExist(statErr) {
		t.Fatal("relaunch with changed authority executed verifier work")
	}
	invBase := filepath.Join(f.root, ".parley-runtime", "invocations", boundVerifierInvocation)
	if _, statErr := os.Stat(filepath.Join(invBase, "started.json")); !os.IsNotExist(statErr) {
		t.Fatal("failed bound relaunch fabricated a started lifecycle")
	}
	if _, statErr := os.Stat(filepath.Join(invBase, "requested.json")); statErr != nil {
		t.Fatal("failed bound relaunch lost its request evidence", statErr)
	}
	if _, statErr := os.Stat(filepath.Join(retryStore.Dir, "ledger.json")); !os.IsNotExist(statErr) {
		t.Fatal("refused bound relaunch was charged")
	}
	if after := readFixtureBytes(t, filepath.Join(f.journal, "recovery.json")); !bytes.Equal(changedRecovery, after) {
		t.Fatal("refused bound relaunch rewrote the retained recovery")
	}
	if after := readFixtureBytes(t, filepath.Join(f.journal, "launch.json")); !bytes.Equal(f.launchRaw, after) {
		t.Fatal("refused bound relaunch rewrote the launch reservation")
	}
}

// TestBoundVerifierRecoveryContextAdmissionIsReadOnly: a valid admission writes
// nothing — the journal, the refused evidence and the accounting stay
// byte-identical, and the bound invocation directory is allocated only by the
// actual launch, never by the context constructor.
func TestBoundVerifierRecoveryContextAdmissionIsReadOnly(t *testing.T) {
	ctx := context.Background()
	f := boundVerifierRefusalFixture(t, true)
	if _, err := WithCapturedVerificationRecovery(ctx, f.ticket, boundVerifierInvocation); err != nil {
		t.Fatalf("trusted recovery context refused a valid retained recovery: %v", err)
	}
	if after := readFixtureBytes(t, filepath.Join(f.journal, "launch.json")); !bytes.Equal(f.launchRaw, after) {
		t.Fatal("admission rewrote the original launch reservation")
	}
	if after := readFixtureBytes(t, filepath.Join(f.journal, "recovery.json")); !bytes.Equal(f.recoveryRaw, after) {
		t.Fatal("admission rewrote the immutable recovery")
	}
	if after := readFixtureBytes(t, filepath.Join(f.binding.Store.Dir, "ledger.json")); !bytes.Equal(f.cycleLedger, after) {
		t.Fatal("admission changed the original fixup accounting")
	}
	if _, err := os.Stat(filepath.Join(f.root, ".parley-runtime", "invocations", boundVerifierInvocation)); !os.IsNotExist(err) {
		t.Fatal("admission allocated the bound invocation directory")
	}
}

// TestBoundVerifierRecoveryContextRefusesMissingAuthorityWithoutWrite: with
// BOTH the launch reservation and the recovery record deleted after a genuine
// recovery, admission is read-only and refuses — it recreates no launch
// reservation (the failure the former existence-probe-plus-mutator design
// risked in its race window), writes no recovery artifact and allocates no
// bound invocation directory.
func TestBoundVerifierRecoveryContextRefusesMissingAuthorityWithoutWrite(t *testing.T) {
	ctx := context.Background()
	f := boundVerifierRefusalFixture(t, true)
	for _, name := range []string{"launch.json", "recovery.json"} {
		if err := os.Remove(filepath.Join(f.journal, name)); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := WithCapturedVerificationRecovery(ctx, f.ticket, boundVerifierInvocation); err == nil ||
		!strings.Contains(err.Error(), "requires a retained launch reservation and immutable recovery record") {
		t.Fatalf("missing authority was admitted: %v", err)
	}
	for _, name := range []string{"launch.json", "recovery.json"} {
		if _, err := os.Stat(filepath.Join(f.journal, name)); !os.IsNotExist(err) {
			t.Fatalf("admission recreated %s", name)
		}
	}
	entries, err := os.ReadDir(f.journal)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != "request.json" {
		t.Fatalf("admission changed the retained journal inventory: %v", entries)
	}
	if _, err := os.Stat(filepath.Join(f.root, ".parley-runtime", "invocations", boundVerifierInvocation)); !os.IsNotExist(err) {
		t.Fatal("admission allocated the bound invocation directory")
	}
}

// TestBoundVerifierRecoveryContextRefusesEvidenceDriftWithoutWrite: deleting
// the refused launch's retained terminal after recovery fails admission with
// the lineage error; the journal stays byte-identical and nothing is allocated.
func TestBoundVerifierRecoveryContextRefusesEvidenceDriftWithoutWrite(t *testing.T) {
	ctx := context.Background()
	f := boundVerifierRefusalFixture(t, true)
	if err := os.Remove(filepath.Join(f.root, ".parley-runtime", "invocations", f.refusedID, "terminal.json")); err != nil {
		t.Fatal(err)
	}
	if _, err := WithCapturedVerificationRecovery(ctx, f.ticket, boundVerifierInvocation); err == nil ||
		!strings.Contains(err.Error(), "refused verifier launch lacks its exact pre-start budget-refusal terminal") {
		t.Fatalf("evidence drift was admitted: %v", err)
	}
	if after := readFixtureBytes(t, filepath.Join(f.journal, "launch.json")); !bytes.Equal(f.launchRaw, after) {
		t.Fatal("refused admission rewrote the launch reservation")
	}
	if after := readFixtureBytes(t, filepath.Join(f.journal, "recovery.json")); !bytes.Equal(f.recoveryRaw, after) {
		t.Fatal("refused admission rewrote the retained recovery")
	}
	if _, err := os.Stat(filepath.Join(f.root, ".parley-runtime", "invocations", boundVerifierInvocation)); !os.IsNotExist(err) {
		t.Fatal("refused admission allocated the bound invocation directory")
	}
}

// TestBoundVerifierRelaunchExecutesExactlyOnce: the pinned invocation is
// exclusively allocated — two concurrent relaunches sharing one recovery
// context produce exactly one execution, and a replay after the fact fails
// without a second process, second work execution or second charge.
func TestBoundVerifierRelaunchExecutesExactlyOnce(t *testing.T) {
	ctx := context.Background()
	f := boundVerifierRefusalFixture(t, true)
	bound, err := WithCapturedVerificationRecovery(ctx, f.ticket, boundVerifierInvocation)
	if err != nil {
		t.Fatal(err)
	}
	retryStore := budget.Store{Dir: filepath.Join(t.TempDir(), "launch-exactly-once"), Scope: "verifier-relaunch"}
	relaunch := WithLaunchBudget(bound, LaunchBudget{Store: retryStore})
	quote := "'" + strings.ReplaceAll(filepath.Join(f.journal, "receipt.json"), "'", "'\\''") + "'"
	script := "printf 'ran\\n' >> .parley-runtime/verifier-relaunched; i=0; while [ ! -f " + quote + " ] && [ \"$i\" -lt 600 ]; do sleep 0.05; i=$((i+1)); done; test -f " + quote
	launch := func() (telemetry.Record, error) {
		verifier := telemetryShell(script, false)
		verifier.ID = "reviewer"
		return RunMeasured(relaunch, ExecOptions{Root: f.root, Agent: verifier, Prompt: "synthetic verifier relaunch", Timeout: 60 * time.Second,
			Info: LaunchInfo{RunID: f.ticket.RunID, Idea: f.idea.Slug, Phase: CapturedVerificationPhase}})
	}
	parent := t.TempDir()
	type result struct {
		record telemetry.Record
		err    error
	}
	results := make(chan result, 2)
	helperDone := make(chan error, 1)
	go func() {
		started := filepath.Join(f.root, ".parley-runtime", "invocations", boundVerifierInvocation, "started.json")
		deadline := time.Now().Add(25 * time.Second)
		for {
			_, err := os.Stat(started)
			if err == nil {
				break
			}
			if !os.IsNotExist(err) {
				helperDone <- err
				return
			}
			if time.Now().After(deadline) {
				helperDone <- errors.New("bound verifier relaunch never started")
				return
			}
			time.Sleep(20 * time.Millisecond)
		}
		_, err := trajectory.ExecuteCapturedVerification(context.Background(), f.ticket, boundVerifierInvocation, f.criteria, parent)
		helperDone <- err
	}()
	for range 2 {
		go func() {
			record, err := launch()
			results <- result{record: record, err: err}
		}()
	}
	var winner telemetry.Record
	successes := 0
	for range 2 {
		r := <-results
		if r.err == nil {
			successes++
			winner = r.record
		}
	}
	if err = <-helperDone; err != nil {
		t.Fatalf("captured helper failed: %v", err)
	}
	if successes != 1 {
		t.Fatalf("bound relaunch executed %d times concurrently", successes)
	}
	if winner.InvocationID != boundVerifierInvocation || winner.StartedAt == nil {
		t.Fatalf("winning relaunch did not execute under the pinned invocation: %+v", winner)
	}
	if replay, replayErr := launch(); replayErr == nil || replay.StartedAt != nil {
		t.Fatalf("replayed bound relaunch executed again: %+v %v", replay, replayErr)
	}
	if marker := readFixtureBytes(t, filepath.Join(f.root, ".parley-runtime", "verifier-relaunched")); string(marker) != "ran\n" {
		t.Fatalf("verifier work did not execute exactly once: %q", marker)
	}
	ledger, err := retryStore.Inspect(ctx)
	if err != nil || len(ledger.Entries) != 1 {
		t.Fatalf("bound relaunch charged %d launch reservations: %v %+v", len(ledger.Entries), err, ledger)
	}
	if _, statErr := os.Stat(filepath.Join(f.journal, "receipt.json")); statErr != nil {
		t.Fatal("captured helper receipt missing from the verification journal", statErr)
	}
}
