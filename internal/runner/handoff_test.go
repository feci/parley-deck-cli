package runner

import (
	"path/filepath"
	"testing"

	"parley-deck-cli/internal/agents"
)

func TestHandoffTelemetryDoesNotInventLaunch(t *testing.T) {
	root := t.TempDir()
	packet, err := WriteHandoffPacket(HandoffOptions{Root: root, RunID: "run-handoff",
		Agent:  agents.Discovery{Spec: agents.Spec{ID: "human-1", LaunchMode: agents.LaunchManual}},
		Prompt: "private prompt", TargetPath: filepath.Join(root, "target.md")})
	if err != nil {
		t.Fatal(err)
	}
	r := terminalRecords(t, root)[0]
	if r.InvocationID != packet.InvocationID || r.StartedAt != nil || r.PID != nil || r.Outcome.ExitCode != nil {
		t.Fatalf("handoff invented process evidence: %+v", r)
	}
	if r.Outcome.Status != "unobserved-handoff" || r.Outcome.Usage.CostUSD != nil {
		t.Fatalf("handoff: %+v", r.Outcome)
	}
}
