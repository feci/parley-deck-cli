package runner

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/budget"
)

// MINOR-1 known cross-resource exhaustion: a launch whose driver-step budget
// is ALREADY exhausted at the known state must refuse before spending the
// protocol cycle. The read-only preflight charges nothing; the real charge
// path stays the authority for anything it cannot know.
func TestLaunchPreflightExhaustedStepCannotSpendCycle(t *testing.T) {
	root := t.TempDir()
	writeLaunchProtocol(t, root)
	ideaDir := filepath.Join(root, "parley-deck", "ideas", "idea")
	mustWrite(t, filepath.Join(ideaDir, "00-prompt.md"), "---\nparticipants: [a, b]\ntrack: deliberation\n---\nFixture\n")
	ctx := context.Background()
	steps, err := budget.EnsureStepBinding(ctx, root, "idea", 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	sctx, finish, err := budget.OpenStepSession(ctx, steps)
	if err != nil {
		t.Fatal(err)
	}
	if err := budget.ChargeStep(sctx); err != nil {
		t.Fatal(err)
	}
	finish()
	cycles, err := budget.EnsureCycleBinding(ctx, root, "idea", budget.Fixup, 5, 0, "", ideaDir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := RunMeasured(ctx, ExecOptions{Root: root, Agent: telemetryShell("", false), Prompt: "fixture", Timeout: time.Second, Info: LaunchInfo{Idea: "idea", Phase: "fixup"}}); err == nil || !strings.Contains(err.Error(), "budget") {
		t.Fatalf("known exhausted step launched: %v", err)
	}
	state, err := cycles.Store.Inspect(ctx)
	if err != nil || cycles.Count(state) != 0 {
		t.Fatalf("known exhausted step spent a cycle: %+v %v", state, err)
	}
	stepState, err := steps.Store.Inspect(ctx)
	if err != nil || budget.StepCount(stepState) != 1 {
		t.Fatalf("charged step was refunded or duplicated: %+v %v", stepState, err)
	}
	records := terminalRecords(t, root)
	if len(records) != 1 || records[0].StartedAt != nil || records[0].Outcome.FailureClass == nil || *records[0].Outcome.FailureClass != "budget_refused" {
		t.Fatalf("known refusal lifecycle: %+v", records)
	}
}

// MINOR-1 reuse invariant: a grouped session that is ALREADY charged for this
// logical action stays reusable at the cap on both resources. A naive
// current-count-only preflight would refuse this; the helpers must not, and
// the grouped charges themselves keep their exact reuse semantics.
func TestLaunchPreflightChargedGroupedSessionReusableAtCap(t *testing.T) {
	root := t.TempDir()
	writeLaunchProtocol(t, root)
	ideaDir := filepath.Join(root, "parley-deck", "ideas", "idea")
	mustWrite(t, filepath.Join(ideaDir, "00-prompt.md"), "---\nparticipants: [a, b]\ntrack: deliberation\n---\nFixture\n")
	ctx := context.Background()
	steps, err := budget.EnsureStepBinding(ctx, root, "idea", 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	cycles, err := budget.EnsureCycleBinding(ctx, root, "idea", budget.Fixup, 1, 0, "", ideaDir)
	if err != nil {
		t.Fatal(err)
	}
	sctx, sfinish, err := budget.OpenStepSession(ctx, steps)
	if err != nil {
		t.Fatal(err)
	}
	defer sfinish()
	if err := budget.ChargeStep(sctx); err != nil {
		t.Fatal(err)
	}
	cctx, cfinish, err := budget.OpenCycleSession(sctx, cycles)
	if err != nil {
		t.Fatal(err)
	}
	defer cfinish()
	if n, err := budget.ChargeCycle(cctx, budget.Fixup); err != nil || n != 1 {
		t.Fatalf("grouped cycle charge: ordinal=%d err=%v", n, err)
	}
	// Both budgets are now AT the cap for this grouped logical action.
	if err := budget.PreflightStepCharge(cctx, steps); err != nil {
		t.Fatalf("preflight refused a charged grouped step session at the cap: %v", err)
	}
	if err := budget.PreflightCycleCharge(cctx, cycles); err != nil {
		t.Fatalf("preflight refused a charged grouped cycle session at the cap: %v", err)
	}
	if err := budget.ChargeStep(cctx); err != nil {
		t.Fatalf("charged grouped step session not reusable at the cap: %v", err)
	}
	if n, err := budget.ChargeCycle(cctx, budget.Fixup); err != nil || n != 1 {
		t.Fatalf("charged grouped cycle session not reusable at the cap: ordinal=%d err=%v", n, err)
	}
	state, err := cycles.Store.Inspect(ctx)
	if err != nil || cycles.Count(state) != 1 {
		t.Fatalf("grouped reuse duplicated the cycle charge: %+v %v", state, err)
	}
	stepState, err := steps.Store.Inspect(ctx)
	if err != nil || budget.StepCount(stepState) != 1 {
		t.Fatalf("grouped reuse duplicated the step charge: %+v %v", stepState, err)
	}
}

// MINOR-1 race honesty: the preflight is a known-state read, not a lock. Two
// launches can both pass it while one step remains; the loser still refuses
// at the real step authority AFTER charging its cycle, and that cycle stays
// spent — no refund, no dedup, no atomicity claim. The interleave below is
// the exact reserveBudget order (cycle charge before step join), made
// deterministic without goroutines.
func TestLaunchPreflightLateRaceRefusalStaysSpent(t *testing.T) {
	root := t.TempDir()
	writeLaunchProtocol(t, root)
	ideaDir := filepath.Join(root, "parley-deck", "ideas", "idea")
	mustWrite(t, filepath.Join(ideaDir, "00-prompt.md"), "---\nparticipants: [a, b]\ntrack: deliberation\n---\nFixture\n")
	ctx := context.Background()
	steps, err := budget.EnsureStepBinding(ctx, root, "idea", 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	cycles, err := budget.EnsureCycleBinding(ctx, root, "idea", budget.Fixup, 5, 0, "", ideaDir)
	if err != nil {
		t.Fatal(err)
	}
	// Both launches preflight the same known state: one step still remains.
	if err := budget.PreflightStepCharge(ctx, steps); err != nil {
		t.Fatal(err)
	}
	if err := budget.PreflightStepCharge(ctx, steps); err != nil {
		t.Fatal(err)
	}
	// Launch A charges cycle then step (reserveBudget order) and wins the step.
	actx, afinish, err := budget.OpenStepSession(ctx, steps)
	if err != nil {
		t.Fatal(err)
	}
	defer afinish()
	acctx, acfinish, err := budget.OpenCycleSession(actx, cycles)
	if err != nil {
		t.Fatal(err)
	}
	defer acfinish()
	if _, err := budget.ChargeCycle(acctx, budget.Fixup); err != nil {
		t.Fatal(err)
	}
	if err := budget.ChargeStep(acctx); err != nil {
		t.Fatal(err)
	}
	// Launch B passed the same preflight, charges its cycle, and only then
	// reaches the step authority: the late refusal must not refund its cycle.
	bcctx, bcfinish, err := budget.OpenCycleSession(ctx, cycles)
	if err != nil {
		t.Fatal(err)
	}
	defer bcfinish()
	if _, err := budget.ChargeCycle(bcctx, budget.Fixup); err != nil {
		t.Fatal(err)
	}
	if _, _, err := budget.OpenStepSession(bcctx, steps); err == nil {
		t.Fatal("late race escaped the step authority")
	}
	state, err := cycles.Store.Inspect(ctx)
	if err != nil || cycles.Count(state) != 2 {
		t.Fatalf("late race cycle was refunded or over-charged: %+v %v", state, err)
	}
	stepState, err := steps.Store.Inspect(ctx)
	if err != nil || budget.StepCount(stepState) != 1 {
		t.Fatalf("late race duplicated the step charge: %+v %v", stepState, err)
	}
}
