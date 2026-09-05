package app

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/telemetry"
)

func measuredFixture(t *testing.T, script string) (string, string) {
	t.Helper()
	root := t.TempDir()
	t.Setenv("PARLEY_HOME", t.TempDir())
	t.Setenv("PARLEY_HEADLESS_AGENT_CONFIG", "")
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "test-agent")
	body := "#!/bin/sh\nif [ \"$1\" = \"--version\" ]; then echo fixture; exit 0; fi\n" + script + "\n"
	if err := os.WriteFile(path, []byte(body), 0o700); err != nil {
		t.Fatal(err)
	}
	writeAgentsLocalConfig(t, root, fakeAgentConfig{ID: "fixture", Path: path})
	prompt := filepath.Join(root, "input.txt")
	if err := os.WriteFile(prompt, []byte("Private prompt, never telemetry"), 0o600); err != nil {
		t.Fatal(err)
	}
	return root, prompt
}

func TestAgentsExecRecordsManualLaunch(t *testing.T) {
	root, prompt := measuredFixture(t, "cat >/dev/null; printf 'agent-owned' > answer.md")
	var out, errOut bytes.Buffer
	code := Run([]string{"agents", "exec", "--dir", root, "--agent", "fixture", "--prompt-file", prompt,
		"--artifact", "answer.md", "--json", "--yes"}, &out, &errOut)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, errOut.String())
	}
	var record telemetry.Record
	if err := json.Unmarshal(out.Bytes(), &record); err != nil {
		t.Fatal(err)
	}
	if record.InvocationID == "" || record.StartedAt == nil || record.Outcome.ArtifactSHA256 == nil || record.Metadata.Phase != "manual" {
		t.Fatalf("record: %+v", record)
	}
	if record.Outcome.Usage.CostUSD != nil || record.Outcome.Usage.ReportedModel != nil {
		t.Fatal("unknown usage/identity invented")
	}
	if strings.Contains(out.String(), "Private prompt") || strings.Contains(out.String(), "agent-owned") {
		t.Fatal("content leaked into public output")
	}
	code = Run([]string{"agents", "exec", "--dir", root, "--agent", "fixture", "--prompt-file", prompt,
		"--artifact", "answer.md", "--yes"}, &out, &errOut)
	if code != 2 {
		t.Fatal("existing artifact overwritten")
	}
}

func TestAgentsExecRequiresExplicitLaunchAndSafePaths(t *testing.T) {
	for _, extra := range [][]string{nil, {"--yes", "--run-id", "run/../../escape"}, {"--yes", "--artifact", "../escape"}} {
		root, prompt := measuredFixture(t, "touch launched")
		args := append([]string{"agents", "exec", "--dir", root, "--agent", "fixture", "--prompt-file", prompt}, extra...)
		var out, errOut bytes.Buffer
		if code := Run(args, &out, &errOut); code != 2 {
			t.Fatalf("code=%d %s", code, errOut.String())
		}
		if _, err := os.Stat(filepath.Join(root, "launched")); !os.IsNotExist(err) {
			t.Fatal("rejected request launched")
		}
	}
}

func TestAgentsExecRetainsFailedStart(t *testing.T) {
	root, prompt := measuredFixture(t, "exit 0")
	if err := os.Remove(filepath.Join(root, "test-agent")); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	code := Run([]string{"agents", "exec", "--dir", root, "--agent", "fixture", "--prompt-file", prompt, "--json", "--yes"}, &out, &errOut)
	if code != 1 {
		t.Fatalf("code=%d %s", code, errOut.String())
	}
	var record telemetry.Record
	if err := json.Unmarshal(out.Bytes(), &record); err != nil {
		t.Fatal(err)
	}
	if record.InvocationID == "" || record.StartedAt != nil || record.Outcome.Status != "failed" {
		t.Fatalf("record: %+v", record)
	}
}

func TestResolveMeasuredAgentReusesRosterMapping(t *testing.T) {
	d := []agents.Discovery{{Spec: agents.Spec{ID: "family", Commands: []string{"absent-command"}}}}
	got, err := resolveMeasuredAgent("named-1", d, map[string]string{"named-1": "family"})
	if err != nil || got.ID != "named-1" || got.Adapter() != "family" || got.Path != "absent-command" {
		t.Fatalf("%+v %v", got, err)
	}
	if d[0].Found {
		t.Fatal("mutated discovery readiness")
	}
	if _, err := resolveMeasuredAgent("../bad", d, nil); err == nil {
		t.Fatal("unsafe identity accepted")
	}
}
