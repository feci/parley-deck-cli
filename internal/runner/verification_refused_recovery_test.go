package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
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

// recoveredVerifierRuntimeFixture mirrors trajectoryRuntimeFixture but freezes
// a criterion whose command emits the structured evidence envelope, so the
// captured helper can actually complete after recovery. All agents are local
// /bin/sh fixtures; no model is invoked.
func recoveredVerifierRuntimeFixture(t *testing.T) (string, protocol.IdeaStatus, []trajectory.Criterion) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("POSIX captured verification")
	}
	root := t.TempDir()
	git := func(args ...string) {
		t.Helper()
		base := []string{"-C", root, "-c", "core.hooksPath=/dev/null", "-c", "commit.gpgSign=false", "-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid"}
		cmd := exec.Command("git", append(base, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("Git: %v %s", err, out)
		}
	}
	git("init", "-q")
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	declareTestLaunchSource(t, root)
	idea, err := protocol.CreateIdea(root, "Verifier refusal recovery fixture", []string{"test-1", "reviewer"})
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(idea.Path, "00-prompt.md"), "---\nidea: "+idea.Slug+"\nparticipants: [test-1, reviewer]\ntrack: deliberation\nstatus: implemented\n---\nFixture\n")
	mustWrite(t, filepath.Join(idea.Path, "IMPLEMENTATION.md"), "---\nidea: "+idea.Slug+"\nstatus: implemented\n---\n\n## Summary of work\nFixture\n")
	mustWrite(t, filepath.Join(root, ".gitignore"), ".parley-runtime/\nparley-deck/runs/\n")
	mustWrite(t, filepath.Join(root, "source"), "original\n")
	git("add", ".")
	git("commit", "-qm", "Recovery fixture")
	criteria := []trajectory.Criterion{{Name: "material", Command: `printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":0}\n'`}}
	if _, err = budget.EnsureCycleBinding(context.Background(), root, idea.Slug, budget.Fixup, 5, 0, "", idea.Path); err != nil {
		t.Fatal(err)
	}
	p, expected, err := trajectory.NewPolicy(context.Background(), root, idea.Slug, "test-1", criteria)
	if err != nil {
		t.Fatal(err)
	}
	if err = trajectory.Activate(context.Background(), root, expected, p); err != nil {
		t.Fatal(err)
	}
	return root, idea, criteria
}

// TestVerifierBudgetRefusalRecoversExplicitly drives the real public launch
// path into its pre-start budget refusal, repairs it with the explicit public
// recovery operation, and proves the subsequent launch is separately checked,
// stale handles cannot control the replacement, and the restored verification
// completes while the original accounting is retained byte-for-byte.
func TestVerifierBudgetRefusalRecoversExplicitly(t *testing.T) {
	ctx := context.Background()
	root, idea, criteria := recoveredVerifierRuntimeFixture(t)
	builder := telemetryShell("printf 'broken\\n' > source; exit 7", false)
	if attempt, err := RunMeasured(ctx, ExecOptions{Root: root, Agent: builder, Prompt: "synthetic patch", Timeout: 30 * time.Second, Info: LaunchInfo{RunID: "patch-attempt", Idea: idea.Slug, Phase: "fixup"}}); err == nil || attempt.StartedAt == nil {
		t.Fatalf("changed-source attempt did not actually execute: %+v %v", attempt, err)
	}
	ticket, err := trajectory.PrepareCapturedVerification(ctx, root, idea.Slug, "reviewer", "verifier-run")
	if err != nil {
		t.Fatal(err)
	}
	wantTicket, err := ticket.SHA256()
	if err != nil {
		t.Fatal(err)
	}
	b, err := budget.LoadCycleBinding(ctx, root, idea.Slug, budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	cycleLedgerPath := filepath.Join(b.Store.Dir, "ledger.json")
	cycleBefore, err := os.ReadFile(cycleLedgerPath)
	if err != nil {
		t.Fatal(err)
	}
	journal := filepath.Join(filepath.Dir(b.Store.Dir), "trajectory-verifications", ticket.Request.Charge.EntryKey)

	// A real pre-start budget refusal through the public launch path.
	bound, err := WithCapturedVerification(ctx, ticket)
	if err != nil {
		t.Fatal(err)
	}
	refusedStore := budget.Store{Dir: filepath.Join(t.TempDir(), "launch"), Scope: "verifier-launch-refusal"}
	denied := WithLaunchBudget(bound, LaunchBudget{Store: refusedStore, Limits: budget.Limits{Denied: map[budget.Kind]bool{budget.Launch: true}}})
	verifier := telemetryShell("touch .parley-runtime/verifier-spawned; exit 0", false)
	verifier.ID = "reviewer"
	options := ExecOptions{Root: root, Agent: verifier, Prompt: "synthetic verifier", Timeout: 30 * time.Second,
		Info: LaunchInfo{RunID: ticket.RunID, Idea: idea.Slug, Phase: CapturedVerificationPhase}}
	refused, err := RunMeasured(denied, options)
	if err == nil || refused.StartedAt != nil || refused.PID != nil || refused.Outcome == nil || refused.Outcome.ExitCode != nil ||
		refused.Outcome.FailureClass == nil || *refused.Outcome.FailureClass != "budget_refused" || refused.InvocationID == "" {
		t.Fatalf("verifier launch was not refused before start: %+v %v", refused, err)
	}
	if _, err = os.Stat(filepath.Join(root, ".parley-runtime", "verifier-spawned")); !os.IsNotExist(err) {
		t.Fatal("refused verifier spawned a child process")
	}
	launchRaw, err := os.ReadFile(filepath.Join(journal, "launch.json"))
	if err != nil {
		t.Fatal(err)
	}
	var reserved struct {
		TicketSHA256 string `json:"ticket_sha256"`
		InvocationID string `json:"invocation_id"`
	}
	if err = json.Unmarshal(launchRaw, &reserved); err != nil || reserved.InvocationID != refused.InvocationID || reserved.TicketSHA256 != wantTicket {
		t.Fatalf("refused launch reservation not bound to the refused invocation and ticket: %+v %v", reserved, err)
	}
	invDir := filepath.Join(root, ".parley-runtime", "invocations", refused.InvocationID)
	for _, name := range []string{"requested.json", "terminal.json"} {
		if _, err = os.Stat(filepath.Join(invDir, name)); err != nil {
			t.Fatal("refused launch lifecycle not retained", err)
		}
	}
	if _, err = os.Stat(filepath.Join(invDir, "started.json")); !os.IsNotExist(err) {
		t.Fatal("pre-start refusal retained a started lifecycle")
	}
	terminalRaw, err := os.ReadFile(filepath.Join(invDir, "terminal.json"))
	if err != nil {
		t.Fatal(err)
	}
	var terminal telemetry.Record
	if err = json.Unmarshal(terminalRaw, &terminal); err != nil {
		t.Fatal(err)
	}
	if terminal.InvocationID != refused.InvocationID || terminal.Metadata.Phase != CapturedVerificationPhase ||
		terminal.StartedAt != nil || terminal.PID != nil || terminal.Outcome == nil ||
		terminal.Outcome.FailureClass == nil || *terminal.Outcome.FailureClass != "budget_refused" || terminal.Outcome.ExitCode != nil {
		t.Fatalf("retained terminal is not the pre-start budget refusal: %+v", terminal)
	}

	// Recovery is one explicit public operation; it is never a launch-path retry.
	if err = trajectory.RecoverCapturedVerificationLaunch(ctx, ticket, refused.InvocationID, "recovered-verifier-invocation"); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(journal)
	if err != nil {
		t.Fatal(err)
	}
	names := map[string]bool{}
	for _, e := range entries {
		names[e.Name()] = true
	}
	if len(names) != 3 || !names["request.json"] || !names["launch.json"] || !names["recovery.json"] {
		t.Fatalf("recovery changed the retained journal inventory: %v", names)
	}
	if after, err := os.ReadFile(filepath.Join(journal, "launch.json")); err != nil || !bytes.Equal(launchRaw, after) {
		t.Fatal("recovery rewrote the refused launch reservation", err)
	}
	if after, err := os.ReadFile(cycleLedgerPath); err != nil || !bytes.Equal(cycleBefore, after) {
		t.Fatal("recovery changed the original fixup accounting", err)
	}
	ledger, err := b.Store.Inspect(ctx)
	if err != nil || b.Count(ledger) != 1 {
		t.Fatalf("original charge count changed: %v %+v", err, ledger)
	}
	if _, err = os.Stat(filepath.Join(refusedStore.Dir, "ledger.json")); !os.IsNotExist(err) {
		t.Fatal("denied launch budget persisted or refunded a reservation")
	}
	recoveryRaw, err := os.ReadFile(filepath.Join(journal, "recovery.json"))
	if err != nil {
		t.Fatal(err)
	}

	// The subsequent launch is separately checked: a real allowed verifier launch
	// whose own generated invocation differs from the recovery binding is
	// refused at the reservation boundary, before its budget and before spawn.
	retryStore := budget.Store{Dir: filepath.Join(t.TempDir(), "launch-retry"), Scope: "verifier-launch-retry"}
	allowed := WithLaunchBudget(bound, LaunchBudget{Store: retryStore})
	retry, err := RunMeasured(allowed, options)
	if err == nil || retry.StartedAt != nil || retry.InvocationID == "" || retry.InvocationID == refused.InvocationID {
		t.Fatalf("separately checked verifier launch rode the recovery: %+v %v", retry, err)
	}
	if _, err = os.Stat(filepath.Join(root, ".parley-runtime", "verifier-spawned")); !os.IsNotExist(err) {
		t.Fatal("separately checked launch spawned a child process")
	}
	if _, err = os.Stat(filepath.Join(retryStore.Dir, "ledger.json")); !os.IsNotExist(err) {
		t.Fatal("refused subsequent launch was charged")
	}
	if after, err := os.ReadFile(filepath.Join(journal, "recovery.json")); err != nil || !bytes.Equal(recoveryRaw, after) {
		t.Fatal("refused subsequent launch rewrote the recovery artifact", err)
	}
	if after, err := os.ReadFile(filepath.Join(journal, "launch.json")); err != nil || !bytes.Equal(launchRaw, after) {
		t.Fatal("refused subsequent launch rewrote the original reservation", err)
	}

	// Stale handles cannot control the replacement.
	if err = trajectory.StopCapturedVerification(ctx, ticket, refused.InvocationID); err == nil || !strings.Contains(err.Error(), "superseded by an explicit recovery") {
		t.Fatalf("stale stop controlled the replacement: %v", err)
	}
	if _, err = os.Stat(filepath.Join(journal, "stop.json")); !os.IsNotExist(err) {
		t.Fatal("stale stop wrote a durable stop record")
	}
	if _, _, err = trajectory.ReadCapturedVerification(ctx, ticket, refused.InvocationID); err == nil {
		t.Fatal("stale read observed the replacement")
	}
	if _, err = trajectory.ExecuteCapturedVerification(ctx, ticket, refused.InvocationID, criteria, ""); err == nil {
		t.Fatal("stale execution consumed the replacement")
	}
	if _, err = os.Stat(filepath.Join(journal, "claim.json")); !os.IsNotExist(err) {
		t.Fatal("stale execution wrote a helper claim")
	}

	// Restored verification completes with original accounting: the helper half
	// executes the frozen criteria under the recovered invocation and its
	// receipt validates against the complete recovered journal lineage.
	receipt, err := trajectory.ExecuteCapturedVerification(ctx, ticket, "recovered-verifier-invocation", criteria, t.TempDir())
	if err != nil || receipt.Steps != 4 || receipt.FailureStage != "" || receipt.InvocationID != "recovered-verifier-invocation" {
		t.Fatalf("recovered verification did not complete: %+v %v", receipt, err)
	}
	read, observation, err := trajectory.ReadCapturedVerification(ctx, ticket, "recovered-verifier-invocation")
	if err != nil || read.TicketSHA256 != wantTicket || read.LastSHA256 != receipt.LastSHA256 {
		t.Fatalf("recovered receipt did not validate: %+v %v", read, err)
	}
	assessment, err := trajectory.AssessCaptured(ticket.Request, observation)
	if err != nil || assessment.Outcome != trajectory.NoRegression {
		t.Fatalf("recovered observation lost its material assessment: %+v %v", assessment, err)
	}
	if after, err := os.ReadFile(cycleLedgerPath); err != nil || !bytes.Equal(cycleBefore, after) {
		t.Fatal("recovered execution changed the original fixup accounting", err)
	}
	// The refusal itself stays spent history: its telemetry is never removed.
	for _, name := range []string{"requested.json", "terminal.json"} {
		if _, err = os.Stat(filepath.Join(invDir, name)); err != nil {
			t.Fatal("refused launch evidence not retained after recovery", err)
		}
	}
}

// TestVerifierRecoveryRefusedEvidenceMutationFailsHandles: the refused launch's
// retained telemetry — produced here by the real public launch path, not a
// fixture — is part of the recovery binding. Rewriting both records after
// recovery, canonically and self-consistently (still a semantically valid
// pre-start budget refusal), is detectable only by digest equality against the
// retained bytes: every stale and current handle must fail closed before any
// write or charge, and the evidence itself stays retained.
func TestVerifierRecoveryRefusedEvidenceMutationFailsHandles(t *testing.T) {
	ctx := context.Background()
	root, idea, criteria := recoveredVerifierRuntimeFixture(t)
	builder := telemetryShell("printf 'broken\\n' > source; exit 7", false)
	if attempt, err := RunMeasured(ctx, ExecOptions{Root: root, Agent: builder, Prompt: "synthetic patch", Timeout: 30 * time.Second, Info: LaunchInfo{RunID: "patch-attempt", Idea: idea.Slug, Phase: "fixup"}}); err == nil || attempt.StartedAt == nil {
		t.Fatalf("changed-source attempt did not actually execute: %+v %v", attempt, err)
	}
	ticket, err := trajectory.PrepareCapturedVerification(ctx, root, idea.Slug, "reviewer", "verifier-run")
	if err != nil {
		t.Fatal(err)
	}
	bound, err := WithCapturedVerification(ctx, ticket)
	if err != nil {
		t.Fatal(err)
	}
	refusedStore := budget.Store{Dir: filepath.Join(t.TempDir(), "launch"), Scope: "verifier-launch-refusal"}
	denied := WithLaunchBudget(bound, LaunchBudget{Store: refusedStore, Limits: budget.Limits{Denied: map[budget.Kind]bool{budget.Launch: true}}})
	verifier := telemetryShell("true", false)
	verifier.ID = "reviewer"
	refused, err := RunMeasured(denied, ExecOptions{Root: root, Agent: verifier, Prompt: "synthetic verifier", Timeout: 30 * time.Second,
		Info: LaunchInfo{RunID: ticket.RunID, Idea: idea.Slug, Phase: CapturedVerificationPhase}})
	if err == nil || refused.StartedAt != nil || refused.Outcome == nil || refused.Outcome.FailureClass == nil ||
		*refused.Outcome.FailureClass != "budget_refused" || refused.InvocationID == "" {
		t.Fatalf("verifier launch was not refused before start: %+v %v", refused, err)
	}
	if err = trajectory.RecoverCapturedVerificationLaunch(ctx, ticket, refused.InvocationID, "recovered-verifier-invocation"); err != nil {
		t.Fatal(err)
	}
	b, err := budget.LoadCycleBinding(ctx, root, idea.Slug, budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	cycleLedgerPath := filepath.Join(b.Store.Dir, "ledger.json")
	cycleBefore, err := os.ReadFile(cycleLedgerPath)
	if err != nil {
		t.Fatal(err)
	}
	journal := filepath.Join(filepath.Dir(b.Store.Dir), "trajectory-verifications", ticket.Request.Charge.EntryKey)

	// Canonical, self-consistent rewrite of both retained records: only the
	// recovery's bound digests can catch it.
	invDir := filepath.Join(root, ".parley-runtime", "invocations", refused.InvocationID)
	for _, name := range []string{"requested.json", "terminal.json"} {
		data, err := os.ReadFile(filepath.Join(invDir, name))
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
		if err = os.WriteFile(filepath.Join(invDir, name), append(rewritten, '\n'), 0600); err != nil {
			t.Fatal(err)
		}
	}
	const mismatch = "no longer matches its retained refused-launch evidence"
	if err = trajectory.StopCapturedVerification(ctx, ticket, refused.InvocationID); err == nil || !strings.Contains(err.Error(), mismatch) {
		t.Fatalf("stale stop did not fail on changed retained evidence: %v", err)
	}
	if err = trajectory.StopCapturedVerification(ctx, ticket, "recovered-verifier-invocation"); err == nil || !strings.Contains(err.Error(), mismatch) {
		t.Fatalf("current stop did not fail on changed retained evidence: %v", err)
	}
	if _, _, err = trajectory.ReadCapturedVerification(ctx, ticket, "recovered-verifier-invocation"); err == nil || !strings.Contains(err.Error(), mismatch) {
		t.Fatalf("current read did not fail on changed retained evidence: %v", err)
	}
	if _, err = trajectory.ExecuteCapturedVerification(ctx, ticket, "recovered-verifier-invocation", criteria, t.TempDir()); err == nil || !strings.Contains(err.Error(), mismatch) {
		t.Fatalf("current execution did not fail on changed retained evidence: %v", err)
	}
	if _, err = os.Stat(filepath.Join(journal, "claim.json")); !os.IsNotExist(err) {
		t.Fatal("refused operation wrote a helper claim")
	}
	if _, err = os.Stat(filepath.Join(journal, "stop.json")); !os.IsNotExist(err) {
		t.Fatal("refused operation wrote a stop record")
	}
	if after, err := os.ReadFile(cycleLedgerPath); err != nil || !bytes.Equal(cycleBefore, after) {
		t.Fatal("refused operations changed the original fixup accounting", err)
	}
	for _, name := range []string{"requested.json", "terminal.json"} {
		if _, err = os.Stat(filepath.Join(invDir, name)); err != nil {
			t.Fatal("refused launch evidence was removed", err)
		}
	}
}
