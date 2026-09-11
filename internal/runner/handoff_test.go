package runner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/agents"
)

func TestHandoffTelemetryDoesNotInventLaunch(t *testing.T) {
	root := t.TempDir()
	writeLaunchProtocol(t, root)
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
	if r.Outcome.Observation.StdoutBytes != nil || r.Outcome.Observation.StderrBytes != nil {
		t.Fatal("handoff invented measured bytes")
	}
	if !strings.Contains(packet.PromptPath, string(filepath.Separator)+".parley-runtime"+string(filepath.Separator)) {
		t.Fatal("private prompt written to a tracked path")
	}
	body, err := os.ReadFile(packet.PromptPath)
	if err != nil || !strings.Contains(string(body), "Mandatory source obligation.") || !strings.Contains(string(body), "private prompt") || r.Metadata.Context.Mode != "full" {
		t.Fatalf("handoff lost attested protocol or task: %+v %v", r.Metadata.Context, err)
	}
}

func TestHandoffWithoutAuthorityPreservesRefusal(t *testing.T) {
	root := t.TempDir()
	_, err := WriteHandoffPacket(HandoffOptions{Root: root, RunID: "refused-handoff",
		Agent: agents.Discovery{Spec: agents.Spec{ID: "human-1", LaunchMode: agents.LaunchManual}}, Prompt: "task"})
	if err == nil {
		t.Fatal("handoff emitted without authority")
	}
	r := terminalRecords(t, root)[0]
	if r.Metadata.Context.Mode != "refused" || r.Outcome.Status != "failed" || r.StartedAt != nil {
		t.Fatalf("invalid handoff refusal: %+v", r)
	}
	if _, err := os.Stat(filepath.Join(root, "parley-deck", "runs", "refused-handoff", "agents", "human-1", "handoff-prompt.md")); !os.IsNotExist(err) {
		t.Fatal("refused handoff emitted a prompt")
	}
}

func TestHandoffAttemptsKeepDistinctPromptBytes(t *testing.T) {
	root := t.TempDir()
	writeLaunchProtocol(t, root)
	opts := HandoffOptions{Root: root, RunID: "same-run", Agent: agents.Discovery{Spec: agents.Spec{ID: "human-1", LaunchMode: agents.LaunchManual}}, Prompt: "first-task"}
	first, err := WriteHandoffPacket(opts)
	if err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(first.PromptPath)
	if err != nil {
		t.Fatal(err)
	}
	opts.Prompt = "second-task"
	second, err := WriteHandoffPacket(opts)
	if err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(first.PromptPath)
	if err != nil || string(before) != string(after) || first.PromptPath == second.PromptPath {
		t.Fatalf("earlier handoff replaced: %v", err)
	}
}
