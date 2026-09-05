package runner

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/store"
	"parley-deck-cli/internal/telemetry"
)

func terminalRecords(t *testing.T, root string) []telemetry.Record {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(root, ".parley-runtime", "invocations", "*", "terminal.json"))
	if err != nil {
		t.Fatal(err)
	}
	var out []telemetry.Record
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		var record telemetry.Record
		if err := json.Unmarshal(data, &record); err != nil {
			t.Fatal(err)
		}
		out = append(out, record)
	}
	return out
}

func telemetryShell(script string, structured bool) agents.Discovery {
	args := []string{"-c", script, "fixture"}
	if structured {
		args = append(args, "--output-format", "json")
	}
	return agents.Discovery{Spec: agents.Spec{ID: "test-1", AdapterID: "claude",
		Model: "requested-model", Reasoning: "max", Speed: "fast",
		HeadlessArgs: args, PromptMode: agents.PromptStdin}, Path: "/bin/sh"}
}

func runTelemetryFixture(ctx context.Context, root string, agent agents.Discovery) error {
	_, err := execAgentProcess(ctx, root, "run-test", agent.ID, "fixture", agent,
		"private prompt", filepath.Join(root, "stdout.log"), filepath.Join(root, "stderr.log"),
		nil, nil, SupervisionConfig{}, supervisionHooks{})
	return err
}

func TestExecTelemetryLifecycleAndUsageEvent(t *testing.T) {
	root := t.TempDir()
	sink := store.New(filepath.Join(root, "events"))
	ctx := WithLaunchInfo(context.Background(), LaunchInfo{RunID: "run-1", SegmentID: "segment-2",
		Idea: "task-C01", Phase: "review", Store: sink})
	agent := telemetryShell(`printf '%s\n' '{"type":"result","usage":{"input_tokens":10,"output_tokens":4},"total_cost_usd":0.125,"modelUsage":{"reported-model":{}}}'`, true)
	if err := runTelemetryFixture(ctx, root, agent); err != nil {
		t.Fatal(err)
	}
	records := terminalRecords(t, root)
	if len(records) != 1 {
		t.Fatalf("records: %d", len(records))
	}
	r := records[0]
	if r.Metadata.RunID != "run-1" || r.Metadata.SegmentID != "segment-2" || r.Metadata.Phase != "review" {
		t.Fatalf("metadata: %+v", r.Metadata)
	}
	if r.StartedAt == nil || r.CompletedAt == nil || r.PID == nil || r.DurationMS == nil {
		t.Fatalf("missing lifecycle: %+v", r)
	}
	if r.StartedAt.Before(r.RequestedAt) || r.CompletedAt.Before(*r.StartedAt) {
		t.Fatal("out of order lifecycle")
	}
	if r.Outcome.Status != "process-exited" || r.Outcome.ExitCode == nil || *r.Outcome.ExitCode != 0 {
		t.Fatalf("outcome: %+v", r.Outcome)
	}
	u := r.Outcome.Usage
	if u.CostUSD == nil || *u.CostUSD != .125 || u.ReportedModel == nil || *u.ReportedModel != "reported-model" {
		t.Fatalf("usage: %+v", u)
	}
	if *r.Metadata.RequestedModel != "requested-model" {
		t.Fatal("requested model replaced by reported")
	}
	if r.Outcome.Observation.FirstActivityMS == nil {
		t.Fatal("missing observed activity")
	}
	events, err := sink.Load()
	if err != nil || len(events) != 1 || events[0].Type != "agent.usage" {
		t.Fatalf("events: %+v %v", events, err)
	}
	if events[0].Data["invocation_id"] != r.InvocationID || events[0].Data["cost_usd"] != .125 {
		t.Fatalf("event: %+v", events[0])
	}
	for _, name := range []string{"requested.json", "started.json", "terminal.json"} {
		path := filepath.Join(root, ".parley-runtime", "invocations", r.InvocationID, name)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(data), "private prompt") || strings.Contains(string(data), "printf") {
			t.Fatal("content leaked into telemetry")
		}
	}
	for _, name := range []string{"stdout.log", "stderr.log"} {
		info, err := os.Stat(filepath.Join(root, name))
		if err != nil || info.Mode().Perm() != 0o600 {
			t.Fatalf("private log: %v %v", info, err)
		}
	}
}

func TestExecTelemetryFailedStartAndUnknownUsage(t *testing.T) {
	root := t.TempDir()
	agent := telemetryShell("", false)
	agent.Path = filepath.Join(root, "missing-executable")
	if err := runTelemetryFixture(context.Background(), root, agent); err == nil {
		t.Fatal("missing executable passed")
	}
	records := terminalRecords(t, root)
	if len(records) != 1 {
		t.Fatalf("records: %d", len(records))
	}
	r := records[0]
	if r.StartedAt != nil || r.PID != nil || r.Outcome.ExitCode != nil {
		t.Fatal("failed start invented process observation")
	}
	if *r.Outcome.FailureClass != "start_failure" || r.Outcome.Usage.CostUSD != nil {
		t.Fatalf("outcome: %+v", r.Outcome)
	}
}

func TestExecTelemetryProviderErrorCannotBeOverriddenByConsultOutput(t *testing.T) {
	root := t.TempDir()
	agent := telemetryShell(`printf '%s\n' '{"type":"result","subtype":"success","is_error":true,"api_error_status":429}'; exit 1`, true)
	res := RunConsult(context.Background(), ConsultOptions{Root: root, Agent: agent, Timeout: time.Second,
		StdoutPath: filepath.Join(root, "out"), StderrPath: filepath.Join(root, "err")})
	if res.ExitError == "" || res.AgentExit != 0 || res.InvocationID == "" {
		t.Fatalf("provider error passed as answer: %+v", res)
	}
	r := terminalRecords(t, root)[0]
	if *r.Outcome.FailureClass != "rate-limit" {
		t.Fatalf("failure: %s", *r.Outcome.FailureClass)
	}
}

func TestExecTelemetryRequestFailurePreventsLaunch(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, ".parley-runtime"), []byte("file"), 0o600); err != nil {
		t.Fatal(err)
	}
	err := runTelemetryFixture(context.Background(), root, telemetryShell("touch launched", false))
	if err == nil {
		t.Fatal("unrecordable launch passed")
	}
	if _, err := os.Stat(filepath.Join(root, "launched")); !os.IsNotExist(err) {
		t.Fatal("process launched despite request failure")
	}
}

func TestExecTelemetryPersistenceFailureCannotPass(t *testing.T) {
	for _, name := range []string{"started.json", "terminal.json"} {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			ctx := WithLaunchInfo(context.Background(), LaunchInfo{Observe: func(r telemetry.Record) {
				if r.Type == "invocation.requested" {
					path := filepath.Join(root, ".parley-runtime", "invocations", r.InvocationID, name)
					if err := os.Mkdir(path, 0o700); err != nil {
						t.Fatal(err)
					}
				}
			}})
			if err := runTelemetryFixture(ctx, root, telemetryShell("exit 0", false)); err == nil {
				t.Fatal("lost evidence passed")
			}
		})
	}
}

func TestExecTelemetryUniqueAttemptsAndRetryLineage(t *testing.T) {
	root := t.TempDir()
	var previous string
	for ordinal := 1; ordinal <= 2; ordinal++ {
		info := LaunchInfo{AttemptOrdinal: ordinal, RetryOf: previous, Observe: func(r telemetry.Record) { previous = r.InvocationID }}
		ctx := WithLaunchInfo(context.Background(), info)
		if err := runTelemetryFixture(ctx, root, telemetryShell("exit 0", false)); err != nil {
			t.Fatal(err)
		}
	}
	records := terminalRecords(t, root)
	if len(records) != 2 || records[0].InvocationID == records[1].InvocationID {
		t.Fatalf("records: %+v", records)
	}
	var first, second telemetry.Record
	for _, r := range records {
		if r.Metadata.AttemptOrdinal == 1 {
			first = r
		} else {
			second = r
		}
	}
	if second.Metadata.RetryOf == nil || *second.Metadata.RetryOf != first.InvocationID {
		t.Fatal("retry lineage lost")
	}
}

func TestExecTelemetryCancelledBeforeStart(t *testing.T) {
	root := t.TempDir()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := runTelemetryFixture(ctx, root, telemetryShell("touch launched", false)); err == nil {
		t.Fatal("cancelled launch passed")
	}
	r := terminalRecords(t, root)[0]
	if r.StartedAt != nil || *r.Outcome.FailureClass != "cancelled" {
		t.Fatalf("record: %+v", r)
	}
}
