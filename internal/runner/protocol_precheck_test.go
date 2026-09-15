package runner

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/telemetry"
	"parley-deck-cli/internal/trajectory"
)

func TestProtocolPrecheckRetainsRefusalAndRechecksAtLaunch(t *testing.T) {
	root, idea := trajectoryRuntimeFixture(t)
	ctx := context.Background()
	agent := telemetryShell("touch .parley-runtime/n1-spawned; exit 7", false)
	var observed telemetry.Record
	info := LaunchInfo{Idea: idea.Slug, RunID: "precheck", Phase: "fixup", Observe: func(r telemetry.Record) { observed = r }}
	if err := PrecheckProtocolLaunch(ctx, root, agent, info); err != nil {
		t.Fatal(err)
	}
	if observed.InvocationID != "" {
		t.Fatal("successful precheck created an invocation")
	}
	restore := refusePreparedProtocolCache(t, root, info)
	if err := PrecheckProtocolLaunch(ctx, root, agent, info); err == nil {
		t.Fatal("precheck accepted refused protocol")
	}
	if observed.InvocationID == "" || observed.Type != "invocation.terminal" || observed.StartedAt != nil || observed.PID != nil || observed.Outcome == nil || observed.Outcome.FailureClass == nil || *observed.Outcome.FailureClass != "protocol_context_refused" {
		t.Fatalf("precheck lost actual unstarted refusal: %+v", observed)
	}
	s, err := trajectory.Inspect(ctx, root, idea.Slug)
	if err != nil || len(s.Attempts) != 0 {
		t.Fatal("precheck consumed a charge", s, err)
	}
	refused, err := os.ReadFile(filepath.Join(root, ".parley-runtime", "invocations", observed.InvocationID, "terminal.json"))
	if err != nil {
		t.Fatal(err)
	}
	restore()
	if err = PrecheckProtocolLaunch(ctx, root, agent, info); err != nil {
		t.Fatal(err)
	}
	// A later refusal after a real caller charge must not reuse the precheck or
	// refund the reservation. The actual launch renderer still owns attestation.
	b, err := budget.LoadCycleBinding(ctx, root, idea.Slug, budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	ctx, finish, err := budget.OpenCycleSession(budget.WithCycleObserver(ctx, &trajectory.Observer{Root: root}), b)
	if err != nil {
		t.Fatal(err)
	}
	defer finish()
	if _, err = budget.ChargeCycle(ctx, budget.Fixup); err != nil {
		t.Fatal(err)
	}
	refusePreparedProtocolCache(t, root, info)
	record, err := RunMeasured(ctx, ExecOptions{Root: root, Agent: agent, Prompt: "fixture", Timeout: 5 * time.Second, Info: LaunchInfo{Idea: idea.Slug, RunID: "late-refusal", Phase: "fixup"}})
	if err == nil || record.StartedAt != nil || record.InvocationID == observed.InvocationID || record.Outcome == nil || record.Outcome.FailureClass == nil || *record.Outcome.FailureClass != "protocol_context_refused" {
		t.Fatal("successful precheck bypassed actual launch attestation", record, err)
	}
	s, err = trajectory.Inspect(ctx, root, idea.Slug)
	if err != nil || len(s.Attempts) != 1 || s.Attempts[0].Launch != nil || s.Attempts[0].Terminal != nil {
		t.Fatal("late refusal rewrote precharge evidence", s, err)
	}
	if _, err = os.Stat(filepath.Join(root, ".parley-runtime", "n1-spawned")); !os.IsNotExist(err) {
		t.Fatal("refused protocol spawned", err)
	}
	current, err := os.ReadFile(filepath.Join(root, ".parley-runtime", "invocations", observed.InvocationID, "terminal.json"))
	if err != nil || string(current) != string(refused) {
		t.Fatal("earlier refusal terminal changed", err)
	}
}
