package app

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/config"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/telemetry"
	"parley-deck-cli/internal/trajectory"
)

// Attended verifier-launch recovery through the app entrypoints. The original
// refusal is driven through the complete production parent path
// (verifyTrajectoryWithAgent) with a real denied launch budget; the recovery
// preview/apply, the separately budget-authorized real relaunch with its real
// helper, and the immutable recovered parent observation all run through the
// new app commands against local fake-CLI processes only — no provider calls.

type verifierRecoveryCLIOutput struct {
	Preview trajectory.CapturedRecoveryPreview `json:"preview"`
	SHA256  string                             `json:"sha256"`
	Applied bool                               `json:"applied"`
}

type verifierRelaunchCLIOutput struct {
	Plan            verifierRelaunchPlan               `json:"plan"`
	SHA256          string                             `json:"sha256"`
	Recovered       *trajectory.RecoveredParentPreview `json:"recovered,omitempty"`
	RecoveredSHA256 string                             `json:"recovered_sha256,omitempty"`
	Applied         bool                               `json:"applied"`
}

type verifierRecoveryAppFixture struct {
	root        string
	trace       string
	journal     string
	binary      string
	runID       string
	refusedID   string
	originalRaw []byte
	original    trajectoryVerificationResult
	launchRaw   []byte
	ledgerRaw   []byte
	binding     *budget.CycleBinding
}

// refusedVerifierAppFixture drives the production verify path into its
// pre-start budget refusal with a real denied budget: the verifier process
// never starts, and the failed parent result is retained exactly as production
// writes it.
func refusedVerifierAppFixture(t *testing.T) *verifierRecoveryAppFixture {
	t.Helper()
	binary := trajectoryHelperBinary(t)
	root, trace, journal, _ := trajectoryHelperFixture(t, "actual-helper")
	ctx := context.Background()
	binding, err := budget.LoadCycleBinding(ctx, root, "idea-x", budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	discovered, err := discoverConfigured(ctx, root)
	if err != nil {
		t.Fatal(err)
	}
	mapping, err := config.LoadRosterAdapters(root)
	if err != nil {
		t.Fatal(err)
	}
	agent, err := resolveMeasuredAgent("reviewer", discovered, mapping)
	if err != nil {
		t.Fatal(err)
	}
	refusedStore := budget.Store{Dir: filepath.Join(t.TempDir(), "launch-refusal"), Scope: "verifier-launch-refusal"}
	denied := runner.WithLaunchBudget(ctx, runner.LaunchBudget{Store: refusedStore, Limits: budget.Limits{Denied: map[budget.Kind]bool{budget.Launch: true}}})
	var progress bytes.Buffer
	result, err := verifyTrajectoryWithAgent(denied, root, "idea-x", agent, 90*time.Second, binary, &progress)
	if err == nil || result.FailureStage != "launch" || result.InvocationID == "" || result.TerminalSHA256 == "" ||
		result.Assessment != nil || !result.TrajectoryPending || result.RunID == "" {
		t.Fatalf("original verifier launch was not refused before start: %+v %v\n%s", result, err, progress.String())
	}
	if _, statErr := os.Stat(filepath.Join(root, ".parley-runtime", "verifier-starts")); !os.IsNotExist(statErr) {
		t.Fatal("refused verifier launched a process")
	}
	base := filepath.Join(root, ".parley-runtime", "trajectory-verification", result.RunID)
	originalRaw, err := os.ReadFile(filepath.Join(base, "parent-result.json"))
	if err != nil {
		t.Fatal(err)
	}
	var original trajectoryVerificationResult
	if err = json.Unmarshal(originalRaw, &original); err != nil {
		t.Fatal(err)
	}
	if original.Version != 1 || original.RunID != result.RunID || original.InvocationID != result.InvocationID ||
		original.TerminalSHA256 != result.TerminalSHA256 || original.RequestSHA256 != result.RequestSHA256 ||
		original.RequestPath != filepath.Join(base, "request.json") || original.ReceiptSHA256 != "" ||
		original.Assessment != nil || original.FailureStage != "launch" || !original.TrajectoryPending {
		t.Fatalf("retained original is not the exact pre-start budget refusal: %+v", original)
	}
	launchRaw, err := os.ReadFile(filepath.Join(journal, "launch.json"))
	if err != nil {
		t.Fatal(err)
	}
	ledgerRaw, err := os.ReadFile(filepath.Join(binding.Store.Dir, "ledger.json"))
	if err != nil {
		t.Fatal(err)
	}
	if _, statErr := os.Stat(filepath.Join(journal, "recovery.json")); !os.IsNotExist(statErr) {
		t.Fatal("refusal wrote a recovery artifact")
	}
	return &verifierRecoveryAppFixture{root: root, trace: trace, journal: journal, binary: binary,
		runID: result.RunID, refusedID: result.InvocationID, originalRaw: originalRaw, original: original,
		launchRaw: launchRaw, ledgerRaw: ledgerRaw, binding: binding}
}

func (f *verifierRecoveryAppFixture) runBase() string {
	return filepath.Join(f.root, ".parley-runtime", "trajectory-verification", f.runID)
}

func (f *verifierRecoveryAppFixture) verifierStarts(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(f.root, ".parley-runtime", "verifier-starts"))
	if os.IsNotExist(err) {
		return ""
	}
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func (f *verifierRecoveryAppFixture) assertUnchanged(t *testing.T) {
	t.Helper()
	if after, err := os.ReadFile(filepath.Join(f.runBase(), "parent-result.json")); err != nil || !bytes.Equal(f.originalRaw, after) {
		t.Fatal("original failed parent result was modified", err)
	}
	if after, err := os.ReadFile(filepath.Join(f.journal, "launch.json")); err != nil || !bytes.Equal(f.launchRaw, after) {
		t.Fatal("original launch reservation was modified", err)
	}
	if after, err := os.ReadFile(filepath.Join(f.binding.Store.Dir, "ledger.json")); err != nil || !bytes.Equal(f.ledgerRaw, after) {
		t.Fatal("original fixup accounting was modified", err)
	}
	invBase := filepath.Join(f.root, ".parley-runtime", "invocations", f.refusedID)
	for _, name := range []string{"requested.json", "terminal.json"} {
		if _, err := os.Stat(filepath.Join(invBase, name)); err != nil {
			t.Fatal("refused launch evidence not retained", err)
		}
	}
	if _, err := os.Stat(filepath.Join(invBase, "started.json")); !os.IsNotExist(err) {
		t.Fatal("refused invocation fabricated a started lifecycle")
	}
}

func runVerifierRecoveryCLI(t *testing.T, ctx context.Context, args ...string) (int, verifierRecoveryCLIOutput, string) {
	t.Helper()
	var out, errout bytes.Buffer
	code := runTrajectoryVerifierRecovery(ctx, args, &out, &errout)
	var decoded verifierRecoveryCLIOutput
	if code == 0 {
		if err := json.Unmarshal(out.Bytes(), &decoded); err != nil {
			t.Fatalf("cannot decode recovery output: %v %s", err, out.String())
		}
	}
	return code, decoded, errout.String()
}

func runVerifierRelaunchCLI(t *testing.T, ctx context.Context, executable string, args ...string) (int, verifierRelaunchCLIOutput, string) {
	t.Helper()
	var out, errout bytes.Buffer
	code := runTrajectoryVerifierRelaunchWith(ctx, args, &out, &errout, executable)
	var decoded verifierRelaunchCLIOutput
	if code == 0 {
		if err := json.Unmarshal(out.Bytes(), &decoded); err != nil {
			t.Fatalf("cannot decode relaunch output: %v %s", err, out.String())
		}
	}
	return code, decoded, errout.String()
}

// applyVerifierRecovery previews (through the app dispatcher) and applies the
// exact inspected recovery, asserting the stale-digest and wrong-identity
// refusals write nothing.
func applyVerifierRecovery(t *testing.T, f *verifierRecoveryAppFixture) verifierRecoveryCLIOutput {
	t.Helper()
	ctx := context.Background()
	// Preview through the app dispatcher entrypoint.
	var out, errout bytes.Buffer
	if code := runTrajectory(ctx, []string{"recover-verifier", "--dir", f.root, "--idea", "idea-x", "--run", f.runID}, &out, &errout); code != 0 {
		t.Fatalf("recovery preview refused: %s", errout.String())
	}
	var preview verifierRecoveryCLIOutput
	if err := json.Unmarshal(out.Bytes(), &preview); err != nil {
		t.Fatal(err)
	}
	if preview.Applied || preview.Preview.RunID != f.runID || preview.Preview.Idea != "idea-x" ||
		preview.Preview.RefusedInvocationID != f.refusedID ||
		preview.Preview.RefusedTerminalSHA256 != f.original.TerminalSHA256 ||
		preview.Preview.InvocationID == "" || preview.Preview.InvocationID == f.refusedID ||
		preview.Preview.RefusedCompletedAt.IsZero() {
		t.Fatalf("preview does not bind the exact refused launch and replacement: %+v", preview)
	}
	if _, err := os.Stat(filepath.Join(f.journal, "recovery.json")); !os.IsNotExist(err) {
		t.Fatal("preview wrote a recovery artifact")
	}
	if starts := f.verifierStarts(t); starts != "" {
		t.Fatalf("preview launched a process: %q", starts)
	}
	// Stale digest and wrong identity refuse before any write.
	if code, _, errText := runVerifierRecoveryCLI(t, ctx, "--dir", f.root, "--idea", "idea-x", "--run", f.runID,
		"--replacement", preview.Preview.InvocationID, "--sha256", strings.Repeat("0", 64), "--yes"); code == 0 ||
		!strings.Contains(errText, "changed since preview") {
		t.Fatalf("stale preview was applied: %d %s", code, errText)
	}
	if code, _, errText := runVerifierRecoveryCLI(t, ctx, "--dir", f.root, "--idea", "idea-x", "--run", f.runID,
		"--replacement", "decoy-verifier-invocation", "--sha256", preview.SHA256, "--yes"); code == 0 {
		t.Fatalf("wrong replacement identity was applied: %d %s", code, errText)
	}
	if _, err := os.Stat(filepath.Join(f.journal, "recovery.json")); !os.IsNotExist(err) {
		t.Fatal("refused applies wrote a recovery artifact")
	}
	// Exact apply, then an exact write-free replay, then a conflicting refusal.
	code, applied, errText := runVerifierRecoveryCLI(t, ctx, "--dir", f.root, "--idea", "idea-x", "--run", f.runID,
		"--replacement", preview.Preview.InvocationID, "--sha256", preview.SHA256, "--yes")
	if code != 0 || !applied.Applied || applied.SHA256 != preview.SHA256 {
		t.Fatalf("exact recovery apply failed: %d %+v %s", code, applied, errText)
	}
	recoveryRaw, err := os.ReadFile(filepath.Join(f.journal, "recovery.json"))
	if err != nil {
		t.Fatal(err)
	}
	code, replayed, errText := runVerifierRecoveryCLI(t, ctx, "--dir", f.root, "--idea", "idea-x", "--run", f.runID,
		"--replacement", preview.Preview.InvocationID, "--sha256", preview.SHA256, "--yes")
	if code != 0 || replayed.SHA256 != preview.SHA256 {
		t.Fatalf("exact replay changed the retained recovery: %d %+v %s", code, replayed, errText)
	}
	if after, err := os.ReadFile(filepath.Join(f.journal, "recovery.json")); err != nil || !bytes.Equal(recoveryRaw, after) {
		t.Fatal("exact replay rewrote the immutable recovery", err)
	}
	if code, _, errText := runVerifierRecoveryCLI(t, ctx, "--dir", f.root, "--idea", "idea-x", "--run", f.runID,
		"--replacement", "conflicting-verifier-invocation", "--sha256", strings.Repeat("0", 64), "--yes"); code == 0 ||
		!strings.Contains(errText, "contradicts the retained immutable recovery") {
		t.Fatalf("conflicting apply was admitted: %d %s", code, errText)
	}
	if starts := f.verifierStarts(t); starts != "" {
		t.Fatalf("recovery or replay launched a process: %q", starts)
	}
	f.assertUnchanged(t)
	return preview
}

// relaunchRecoveredVerifierApp drives the plan preview, the stale-plan refusal
// and the real budget-authorized relaunch, returning the plan and the relaunch
// output.
func relaunchRecoveredVerifierApp(t *testing.T, f *verifierRecoveryAppFixture, replacement string) (verifierRelaunchCLIOutput, verifierRelaunchCLIOutput) {
	t.Helper()
	ctx := context.Background()
	code, plan, errText := runVerifierRelaunchCLI(t, ctx, f.binary, "--dir", f.root, "--idea", "idea-x", "--run", f.runID)
	if code != 0 || plan.Applied || plan.Plan.Recovery.InvocationID != replacement ||
		plan.Plan.Recovery.RefusedInvocationID != f.refusedID || plan.Plan.Verifier != "reviewer" ||
		plan.Plan.RunID != f.runID || plan.Plan.Idea != "idea-x" {
		t.Fatalf("relaunch plan does not bind the validated recovery identity: %d %+v %s", code, plan, errText)
	}
	if starts := f.verifierStarts(t); starts != "" {
		t.Fatalf("relaunch plan preview launched a process: %q", starts)
	}
	if code, _, errText := runVerifierRelaunchCLI(t, ctx, f.binary, "--dir", f.root, "--idea", "idea-x", "--run", f.runID,
		"--sha256", strings.Repeat("0", 64), "--yes"); code == 0 || !strings.Contains(errText, "plan changed since preview") {
		t.Fatalf("stale relaunch plan was applied: %d %s", code, errText)
	}
	if starts := f.verifierStarts(t); starts != "" {
		t.Fatalf("stale relaunch plan launched a process: %q", starts)
	}
	authorizedStore := budget.Store{Dir: filepath.Join(t.TempDir(), "launch-authorized"), Scope: "verifier-relaunch-authorized"}
	authorized := runner.WithLaunchBudget(ctx, runner.LaunchBudget{Store: authorizedStore})
	code, applied, errText := runVerifierRelaunchCLI(t, authorized, f.binary, "--dir", f.root, "--idea", "idea-x", "--run", f.runID,
		"--sha256", plan.SHA256, "--timeout", "90s", "--yes")
	if code != 0 || !applied.Applied || applied.Recovered == nil {
		t.Fatalf("budget-authorized relaunch failed: %d %+v %s", code, applied, errText)
	}
	ledger, err := authorizedStore.Inspect(ctx)
	if err != nil || len(ledger.Entries) != 1 {
		t.Fatalf("relaunch was not separately charged exactly once: %v %+v", err, ledger)
	}
	return plan, applied
}

// TestAppVerifierRecoveryFullLifecycle: the complete attended lifecycle through
// the app entrypoints — production refusal, exact preview/apply, separately
// budget-authorized real relaunch plus real helper under the pinned invocation,
// immutable recovered observation, explicit fresh reconciliation — with the
// original failure, refusal evidence, reservation and accounting byte-exact.
func TestAppVerifierRecoveryFullLifecycle(t *testing.T) {
	f := refusedVerifierAppFixture(t)
	ctx := context.Background()
	preview := applyVerifierRecovery(t, f)
	replacement := preview.Preview.InvocationID

	_, applied := relaunchRecoveredVerifierApp(t, f, replacement)

	// The real verifier ran exactly once under the pinned invocation, and the
	// real helper executed the original AB/BA criteria.
	if starts := f.verifierStarts(t); starts != replacement+"\n" {
		t.Fatalf("verifier did not run exactly once under the pinned invocation: %q", starts)
	}
	if data, err := os.ReadFile(f.trace); err != nil ||
		string(data) != "reviewer:original\nreviewer:broken\nreviewer:broken\nreviewer:original\n" {
		t.Fatalf("independent helper did not execute original AB/BA criteria: %q %v", data, err)
	}
	invBase := filepath.Join(f.root, ".parley-runtime", "invocations", replacement)
	var terminal telemetry.Record
	if _, err := readVerificationJSON(filepath.Join(invBase, "terminal.json"), &terminal); err != nil ||
		terminal.InvocationID != replacement || terminal.Metadata.RunID != f.runID || terminal.Metadata.Idea != "idea-x" ||
		terminal.Metadata.Agent != "reviewer" || terminal.Metadata.Phase != runner.CapturedVerificationPhase ||
		terminal.StartedAt == nil || terminal.PID == nil || terminal.CompletedAt == nil ||
		terminal.Outcome == nil || terminal.Outcome.Status != "process-exited" ||
		terminal.Outcome.ExitCode == nil || *terminal.Outcome.ExitCode != 0 {
		t.Fatalf("replacement lifecycle is not the successful observed verifier: %+v %v", terminal, err)
	}
	for _, name := range []string{"requested.json", "started.json"} {
		if _, err := os.Stat(filepath.Join(invBase, name)); err != nil {
			t.Fatal("replacement lifecycle evidence missing", err)
		}
	}
	recovered := applied.Recovered
	if recovered.OriginalResult.InvocationID != f.refusedID || recovered.OriginalResult.TerminalSHA256 != f.original.TerminalSHA256 ||
		recovered.OriginalResult.FailureStage != "launch" || recovered.Recovery.InvocationID != replacement ||
		recovered.Derived.Result.InvocationID != replacement || recovered.Derived.Result.Assessment == nil ||
		recovered.Derived.Result.Assessment.Outcome != trajectory.Regression {
		t.Fatalf("recovered observation does not tie the original failure to the replacement lineage: %+v", recovered)
	}
	if _, err := os.Stat(filepath.Join(f.runBase(), "parent-recovered.json")); err != nil {
		t.Fatal("immutable recovered parent observation missing", err)
	}
	f.assertUnchanged(t)

	// Explicit fresh reconciliation over the complete validated lineage, with
	// the observation's provenance pinned, then an exact replay.
	recon, err := trajectory.PreviewReconciliation(ctx, f.root, "idea-x", f.runID)
	if err != nil {
		t.Fatalf("fresh reconciliation refused the recovered lineage: %v", err)
	}
	if recon.RecoverySHA256 == "" || recon.Assessment.Outcome != trajectory.Regression || recon.ParentSHA256 != recovered.ParentSHA256 {
		t.Fatalf("recovered reconciliation lost its provenance, outcome or parent binding: %+v", recon)
	}
	expected := recon.SHA256()
	if err = trajectory.Reconcile(ctx, f.root, "idea-x", f.runID, expected); err != nil {
		t.Fatalf("explicit reconciliation refused the recovered lineage: %v", err)
	}
	if err = trajectory.Reconcile(ctx, f.root, "idea-x", f.runID, expected); err != nil {
		t.Fatalf("exact reconciliation replay changed its decision: %v", err)
	}
	state, err := trajectory.Inspect(ctx, f.root, "idea-x")
	if err != nil || len(state.Resolutions) != 1 || state.Resolutions[0].SHA256 != expected {
		t.Fatalf("recovered resolution was not retained: %v %+v", err, state)
	}
	if err = trajectory.RequireResolved(ctx, f.root, "idea-x"); err == nil {
		t.Fatal("regression resolution closed the trajectory without a clean outcome")
	}
	f.assertUnchanged(t)

	// Duplicate relaunch: the pinned invocation is consumed — no second
	// process, no second charge, nothing published again.
	authorizedStore := budget.Store{Dir: filepath.Join(t.TempDir(), "launch-duplicate"), Scope: "verifier-relaunch-duplicate"}
	authorized := runner.WithLaunchBudget(ctx, runner.LaunchBudget{Store: authorizedStore})
	code, _, errText := runVerifierRelaunchCLI(t, authorized, f.binary, "--dir", f.root, "--idea", "idea-x", "--run", f.runID,
		"--sha256", applied.SHA256, "--timeout", "90s", "--yes")
	if code == 0 {
		t.Fatalf("duplicate relaunch executed again: %d %s", code, errText)
	}
	if starts := f.verifierStarts(t); starts != replacement+"\n" {
		t.Fatalf("duplicate relaunch started another process: %q", starts)
	}
	if data, err := os.ReadFile(f.trace); err != nil ||
		string(data) != "reviewer:original\nreviewer:broken\nreviewer:broken\nreviewer:original\n" {
		t.Fatalf("duplicate relaunch repeated the helper: %q %v", data, err)
	}
	if _, err := os.Stat(filepath.Join(authorizedStore.Dir, "ledger.json")); !os.IsNotExist(err) {
		t.Fatal("duplicate relaunch was charged")
	}
	// Recovered-observation publish replay stays exact and write-free.
	replayed, err := trajectory.PublishRecoveredParent(ctx, f.root, "idea-x", f.runID, applied.RecoveredSHA256)
	if err != nil || replayed.SHA256() != applied.Recovered.SHA256() {
		t.Fatalf("recovered observation replay changed its decision: %+v %v", replayed, err)
	}
	f.assertUnchanged(t)
}

// TestAppVerifierRelaunchRequiresRecovery: without the applied recovery the
// relaunch refuses at the read-only identity resolution and launches nothing.
func TestAppVerifierRelaunchRequiresRecovery(t *testing.T) {
	f := refusedVerifierAppFixture(t)
	ctx := context.Background()
	code, _, errText := runVerifierRelaunchCLI(t, ctx, f.binary, "--dir", f.root, "--idea", "idea-x", "--run", f.runID)
	if code == 0 || !strings.Contains(errText, "requires a retained launch reservation and immutable recovery record") {
		t.Fatalf("relaunch preview was admitted without a recovery: %d %s", code, errText)
	}
	if starts := f.verifierStarts(t); starts != "" {
		t.Fatalf("refused relaunch launched a process: %q", starts)
	}
	f.assertUnchanged(t)
}

// TestAppVerifierRecoveryRefusedEvidenceMutationFailsRelaunch: canonical
// post-recovery mutation of the refused launch's retained evidence fails the
// relaunch plan at the read-only reader — no process, no write.
func TestAppVerifierRecoveryRefusedEvidenceMutationFailsRelaunch(t *testing.T) {
	f := refusedVerifierAppFixture(t)
	ctx := context.Background()
	applyVerifierRecovery(t, f)
	invBase := filepath.Join(f.root, ".parley-runtime", "invocations", f.refusedID)
	for _, name := range []string{"requested.json", "terminal.json"} {
		var record telemetry.Record
		if _, err := readVerificationJSON(filepath.Join(invBase, name), &record); err != nil {
			t.Fatal(err)
		}
		record.Metadata.AttemptOrdinal = 99
		data, err := json.MarshalIndent(record, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(filepath.Join(invBase, name), append(data, '\n'), 0600); err != nil {
			t.Fatal(err)
		}
	}
	code, _, errText := runVerifierRelaunchCLI(t, ctx, f.binary, "--dir", f.root, "--idea", "idea-x", "--run", f.runID)
	if code == 0 || !strings.Contains(errText, "no longer matches its retained refused-launch evidence") {
		t.Fatalf("relaunch was admitted over mutated refused evidence: %d %s", code, errText)
	}
	if starts := f.verifierStarts(t); starts != "" {
		t.Fatalf("refused relaunch launched a process: %q", starts)
	}
	if _, err := os.Stat(filepath.Join(f.runBase(), "parent-recovered.json")); !os.IsNotExist(err) {
		t.Fatal("refused relaunch published an observation")
	}
}

// TestAppVerifierRecoveryFailedReplacementRemainsUnresolved: a budget-refused
// relaunch consumes the pinned invocation with its own failed terminal — no
// started lifecycle, no observation, no resolution, and no second chance
// through a repeated recovery apply.
func TestAppVerifierRecoveryFailedReplacementRemainsUnresolved(t *testing.T) {
	f := refusedVerifierAppFixture(t)
	ctx := context.Background()
	preview := applyVerifierRecovery(t, f)
	replacement := preview.Preview.InvocationID

	code, plan, errText := runVerifierRelaunchCLI(t, ctx, f.binary, "--dir", f.root, "--idea", "idea-x", "--run", f.runID)
	if code != 0 {
		t.Fatalf("relaunch plan preview refused: %d %s", code, errText)
	}
	deniedStore := budget.Store{Dir: filepath.Join(t.TempDir(), "launch-denied-again"), Scope: "verifier-relaunch-denied"}
	denied := runner.WithLaunchBudget(ctx, runner.LaunchBudget{Store: deniedStore, Limits: budget.Limits{Denied: map[budget.Kind]bool{budget.Launch: true}}})
	code, _, errText = runVerifierRelaunchCLI(t, denied, f.binary, "--dir", f.root, "--idea", "idea-x", "--run", f.runID,
		"--sha256", plan.SHA256, "--timeout", "90s", "--yes")
	if code == 0 {
		t.Fatalf("budget-refused relaunch reported success: %d %s", code, errText)
	}
	invBase := filepath.Join(f.root, ".parley-runtime", "invocations", replacement)
	var terminal telemetry.Record
	if _, err := readVerificationJSON(filepath.Join(invBase, "terminal.json"), &terminal); err != nil ||
		terminal.InvocationID != replacement || terminal.Outcome == nil || terminal.Outcome.FailureClass == nil ||
		*terminal.Outcome.FailureClass != "budget_refused" {
		t.Fatalf("failed replacement did not retain its own budget refusal: %+v %v", terminal, err)
	}
	if _, err := os.Stat(filepath.Join(invBase, "started.json")); !os.IsNotExist(err) {
		t.Fatal("failed replacement fabricated a started lifecycle")
	}
	if starts := f.verifierStarts(t); starts != "" {
		t.Fatalf("budget-refused relaunch started a process: %q", starts)
	}
	if _, err := os.Stat(filepath.Join(f.runBase(), "parent-recovered.json")); !os.IsNotExist(err) {
		t.Fatal("failed replacement published an observation")
	}
	if _, err := os.Stat(filepath.Join(deniedStore.Dir, "ledger.json")); !os.IsNotExist(err) {
		t.Fatal("denied relaunch persisted or refunded a reservation")
	}
	if _, err := trajectory.PreviewReconciliation(ctx, f.root, "idea-x", f.runID); err == nil {
		t.Fatal("failed replacement resolved the trajectory")
	}
	if err := trajectory.RequireResolved(ctx, f.root, "idea-x"); err == nil {
		t.Fatal("failed replacement closed the trajectory")
	}
	// The recovery grants no second chance: a repeated apply fails closed, and
	// a repeated relaunch cannot consume the pinned invocation again.
	if code, _, errText := runVerifierRecoveryCLI(t, ctx, "--dir", f.root, "--idea", "idea-x", "--run", f.runID,
		"--replacement", "another-verifier-invocation", "--sha256", strings.Repeat("0", 64), "--yes"); code == 0 ||
		!strings.Contains(errText, "contradicts the retained immutable recovery") {
		t.Fatalf("failed replacement was re-recovered: %d %s", code, errText)
	}
	authorizedStore := budget.Store{Dir: filepath.Join(t.TempDir(), "launch-authorized"), Scope: "verifier-relaunch-authorized"}
	authorized := runner.WithLaunchBudget(ctx, runner.LaunchBudget{Store: authorizedStore})
	if code, _, errText := runVerifierRelaunchCLI(t, authorized, f.binary, "--dir", f.root, "--idea", "idea-x", "--run", f.runID,
		"--sha256", plan.SHA256, "--timeout", "90s", "--yes"); code == 0 {
		t.Fatalf("consumed pinned invocation relaunched: %d %s", code, errText)
	}
	if starts := f.verifierStarts(t); starts != "" {
		t.Fatalf("replayed relaunch started a process: %q", starts)
	}
	if _, err := os.Stat(filepath.Join(authorizedStore.Dir, "ledger.json")); !os.IsNotExist(err) {
		t.Fatal("replayed relaunch was charged")
	}
	f.assertUnchanged(t)
}

// TestAppVerifierRecoveredLineageMutationFailsReconciliation: after the full
// lifecycle, deletion or mutation of the original failure or the replacement
// evidence fails fresh reconciliation and the later state guards.
func TestAppVerifierRecoveredLineageMutationFailsReconciliation(t *testing.T) {
	ctx := context.Background()
	for _, tc := range []struct {
		name    string
		mutate  func(t *testing.T, f *verifierRecoveryAppFixture, replacement string)
		message string
	}{
		{name: "original-parent-deleted", mutate: func(t *testing.T, f *verifierRecoveryAppFixture, _ string) {
			if err := os.Remove(filepath.Join(f.runBase(), "parent-result.json")); err != nil {
				t.Fatal(err)
			}
		}, message: "lost its original refused parent result"},
		{name: "original-parent-rewritten", mutate: func(t *testing.T, f *verifierRecoveryAppFixture, _ string) {
			forged := f.original
			forged.TerminalSHA256 = strings.Repeat("f", 64)
			data, err := json.MarshalIndent(forged, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			if err = os.WriteFile(filepath.Join(f.runBase(), "parent-result.json"), append(data, '\n'), 0600); err != nil {
				t.Fatal(err)
			}
		}, message: "original refused parent result changed after recovery"},
		{name: "replacement-terminal-deleted", mutate: func(t *testing.T, f *verifierRecoveryAppFixture, replacement string) {
			if err := os.Remove(filepath.Join(f.root, ".parley-runtime", "invocations", replacement, "terminal.json")); err != nil {
				t.Fatal(err)
			}
		}, message: "successful observed parent terminal binding is unavailable"},
		{name: "replacement-terminal-rewritten", mutate: func(t *testing.T, f *verifierRecoveryAppFixture, replacement string) {
			path := filepath.Join(f.root, ".parley-runtime", "invocations", replacement, "terminal.json")
			var record telemetry.Record
			if _, err := readVerificationJSON(path, &record); err != nil {
				t.Fatal(err)
			}
			record.Metadata.AttemptOrdinal = 99
			data, err := json.MarshalIndent(record, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			if err = os.WriteFile(path, append(data, '\n'), 0600); err != nil {
				t.Fatal(err)
			}
		}, message: "original requested lifecycle differs from the verifier terminal"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := refusedVerifierAppFixture(t)
			preview := applyVerifierRecovery(t, f)
			relaunchRecoveredVerifierApp(t, f, preview.Preview.InvocationID)
			recon, err := trajectory.PreviewReconciliation(ctx, f.root, "idea-x", f.runID)
			if err != nil {
				t.Fatal(err)
			}
			if err = trajectory.Reconcile(ctx, f.root, "idea-x", f.runID, recon.SHA256()); err != nil {
				t.Fatal(err)
			}
			tc.mutate(t, f, preview.Preview.InvocationID)
			if _, err = trajectory.PreviewReconciliation(ctx, f.root, "idea-x", f.runID); err == nil ||
				!strings.Contains(err.Error(), tc.message) {
				t.Fatalf("%s did not fail fresh reconciliation: %v", tc.name, err)
			}
			if _, err = trajectory.Inspect(ctx, f.root, "idea-x"); err == nil || !strings.Contains(err.Error(), tc.message) {
				t.Fatalf("%s did not fail the later state guard: %v", tc.name, err)
			}
			if err = trajectory.RequireResolved(ctx, f.root, "idea-x"); err == nil {
				t.Fatalf("%s closed the trajectory", tc.name)
			}
		})
	}
}
