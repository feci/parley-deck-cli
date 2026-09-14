package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/protocolpacket"
	"parley-deck-cli/internal/trajectory"
)

func TestKnownProtocolRefusalDoesNotSpendNewAttempt(t *testing.T) {
	root, idea := trajectoryRuntimeFixture(t)
	ctx := context.Background()
	step, err := budget.EnsureStepBinding(ctx, root, idea.Slug, 5, 0)
	if err != nil {
		t.Fatal(err)
	}
	launch := budget.Store{Dir: filepath.Join(t.TempDir(), "launch-ledger"), Scope: "refusal-fixture"}
	ctx = WithLaunchBudget(ctx, LaunchBudget{Store: launch, Limits: budget.Limits{Actions: map[budget.Kind]int{budget.Launch: 5}}})
	opts := ExecOptions{Root: root, Agent: telemetryShell("touch .parley-runtime/spawned; exit 7", false), Prompt: "fixture", Timeout: 5 * time.Second, Info: LaunchInfo{Idea: idea.Slug, RunID: "known-refusal", Phase: "fixup"}}
	restore := refusePreparedProtocolCache(t, root, opts.Info)
	r, err := RunMeasured(ctx, opts)
	if err == nil || r.InvocationID == "" || r.StartedAt != nil || r.PID != nil || r.Outcome == nil || r.Outcome.FailureClass == nil || *r.Outcome.FailureClass != "protocol_context_refused" {
		t.Fatalf("known refusal lost exact unstarted terminal: %+v %v", r, err)
	}
	s, err := trajectory.Inspect(ctx, root, idea.Slug)
	if err != nil || s == nil || len(s.Attempts) != 0 {
		t.Fatalf("known refusal charged a fresh trajectory attempt: %+v %v", s, err)
	}
	for _, store := range []budget.Store{step.Store, launch} {
		ledger, err := store.Inspect(ctx)
		if store.Dir == launch.Dir && os.IsNotExist(err) {
			continue
		}
		if err != nil || len(ledger.Entries) != 0 {
			t.Fatalf("known refusal spent another budget: %+v %v", ledger, err)
		}
	}
	if _, err = os.Stat(filepath.Join(root, ".parley-runtime", "spawned")); !os.IsNotExist(err) {
		t.Fatal("known refusal spawned")
	}
	restore()
	next, err := RunMeasured(ctx, opts)
	if err == nil || next.StartedAt == nil || next.InvocationID == r.InvocationID {
		t.Fatalf("restored protocol lost ordinary launch: %+v %v", next, err)
	}
	s, err = trajectory.Inspect(ctx, root, idea.Slug)
	if err != nil || len(s.Attempts) != 1 || s.Attempts[0].Launch == nil || s.Attempts[0].Launch.InvocationID != next.InvocationID {
		t.Fatalf("ordinary launch did not spend one fresh attempt: %+v %v", s, err)
	}
	for _, store := range []budget.Store{step.Store, launch} {
		ledger, err := store.Inspect(ctx)
		if err != nil || len(ledger.Entries) != 1 {
			t.Fatalf("ordinary launch lost required reservation: %+v %v", ledger, err)
		}
	}
}
func TestKnownProtocolRefusalPreservesAlreadySpentAttempt(t *testing.T) {
	root, idea := trajectoryRuntimeFixture(t)
	b, err := budget.LoadCycleBinding(context.Background(), root, idea.Slug, budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	ctx, finish, err := budget.OpenCycleSession(budget.WithCycleObserver(context.Background(), &trajectory.Observer{Root: root}), b)
	if err != nil {
		t.Fatal(err)
	}
	defer finish()
	if _, err = budget.ChargeCycle(ctx, budget.Fixup); err != nil {
		t.Fatal(err)
	}
	before, err := trajectory.Inspect(ctx, root, idea.Slug)
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal(before)
	refusePreparedProtocolCache(t, root, LaunchInfo{Idea: idea.Slug, Phase: "fixup"})
	r, err := RunMeasured(ctx, ExecOptions{Root: root, Agent: telemetryShell("exit 0", false), Prompt: "fixture", Info: LaunchInfo{Idea: idea.Slug, RunID: "precharged-refusal", Phase: "fixup"}})
	if err == nil || r.StartedAt != nil {
		t.Fatalf("precharged known refusal executed: %+v %v", r, err)
	}
	after, err := trajectory.Inspect(ctx, root, idea.Slug)
	if err != nil {
		t.Fatal(err)
	}
	actual, _ := json.Marshal(after)
	if !bytes.Equal(raw, actual) || len(after.Attempts) != 1 || after.Attempts[0].Launch != nil || after.Attempts[0].Terminal != nil {
		t.Fatalf("known refusal rewrote already spent source/launch evidence: %+v", after)
	}
	ledger, err := b.Store.Inspect(ctx)
	if err != nil || b.Count(ledger) != 1 {
		t.Fatal("known refusal refunded a precharge", err)
	}
	if err = trajectory.RequireResolved(ctx, root, idea.Slug); err == nil {
		t.Fatal("known refusal established completion")
	}
}
func TestKnownProtocolRefusalLeavesPreparedVerifierUnconsumed(t *testing.T) {
	ticket, journal := verifierLaunchFixture(t)
	ctx, err := WithCapturedVerification(context.Background(), ticket)
	if err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(ticket.Root, "parley-deck", "COOPERATION.md")
	original, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	originalInfo, err := os.Stat(file)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.Remove(file); err != nil {
		t.Fatal(err)
	}
	agent := telemetryShell("exit 0", false)
	agent.ID = "reviewer"
	opts := ExecOptions{Root: ticket.Root, Agent: agent, Prompt: "fixture", Timeout: 5 * time.Second, Info: LaunchInfo{RunID: ticket.RunID, Idea: ticket.Request.Idea, Phase: CapturedVerificationPhase}}
	r, err := RunMeasured(ctx, opts)
	if err == nil || r.StartedAt != nil || r.Outcome == nil || r.Outcome.FailureClass == nil || *r.Outcome.FailureClass != "protocol_context_refused" {
		t.Fatalf("known refusal lost verifier failure: %+v %v", r, err)
	}
	if _, err = os.Stat(filepath.Join(journal, "launch.json")); !os.IsNotExist(err) {
		t.Fatal("known refusal consumed prepared verifier ticket")
	}
	if err = os.WriteFile(file, original, originalInfo.Mode().Perm()); err != nil {
		t.Fatal(err)
	}
	next, err := RunMeasured(ctx, opts)
	if err != nil || next.StartedAt == nil {
		t.Fatalf("restored protocol could not use previously unconsumed ticket: %+v %v", next, err)
	}
	if _, err = os.Stat(filepath.Join(journal, "launch.json")); err != nil {
		t.Fatal("ordinary verifier lacked launch reservation", err)
	}
}

// Corrupt only the ignored cached publication. Material source still equals the
// frozen activation baseline, so a later source mismatch cannot mask charging.
func refusePreparedProtocolCache(t *testing.T, root string, info LaunchInfo) func() {
	t.Helper()
	if _, _, err := prepareProtocolPrompt(root, "fixture", info); err != nil {
		t.Fatal(err)
	}
	paths, err := filepath.Glob(filepath.Join(protocolpacket.RuntimeDir(root), "*.md"))
	if err != nil || len(paths) != 1 {
		t.Fatalf("expected one prepared full context: %v %v", paths, err)
	}
	original, err := os.ReadFile(paths[0])
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(paths[0], []byte("tampered cached context"), 0600); err != nil {
		t.Fatal(err)
	}
	return func() {
		if err := os.WriteFile(paths[0], original, 0600); err != nil {
			t.Fatal(err)
		}
	}
}
