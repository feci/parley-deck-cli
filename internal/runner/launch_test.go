package runner

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/telemetry"
)

func TestTrackedCommandForRunAndCombinedOutput(t *testing.T) {
	for _, combined := range []bool{false, true} {
		root := t.TempDir()
		writeLaunchProtocol(t, root)
		cmd, cleanup, err := trackedCommandFor(context.Background(), root, telemetryShell("printf 'hello'; printf 'world' >&2", false), "prompt")
		if err != nil {
			t.Fatal(err)
		}
		var output []byte
		if combined {
			output, err = cmd.CombinedOutput()
		} else {
			output, err = cmd.Output()
		}
		cleanup()
		if err != nil || !strings.Contains(string(output), "hello") {
			t.Fatalf("output: %q %v", output, err)
		}
		if combined && !strings.Contains(string(output), "world") {
			t.Fatalf("stderr missing: %q", output)
		}
		records := terminalRecords(t, root)
		if len(records) != 1 || records[0].Outcome.Observation.StdoutBytes != 5 || records[0].Outcome.Observation.StderrBytes != 5 {
			t.Fatalf("records: %+v", records)
		}
		if err := cmd.Run(); err == nil {
			t.Fatal("command reused")
		}
	}
}

func TestTrackedCommandForSeparateStartWait(t *testing.T) {
	root := t.TempDir()
	writeLaunchProtocol(t, root)
	cmd, cleanup, err := trackedCommandFor(context.Background(), root, telemetryShell("cat", false), "input")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	var output bytes.Buffer
	cmd.Stdout = &output
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	if len(terminalRecords(t, root)) != 0 {
		t.Fatal("terminal event before wait")
	}
	if err := cmd.Wait(); err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(output.String(), "Launch task:\ninput") || !strings.Contains(output.String(), "Mandatory source obligation.") {
		t.Fatalf("stdin changed: %q", output.String())
	}
	if len(terminalRecords(t, root)) != 1 {
		t.Fatal("terminal record missing")
	}
}

func TestTrackedCommandForTimeoutKillsChildGroup(t *testing.T) {
	root := t.TempDir()
	writeLaunchProtocol(t, root)
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()
	cmd, cleanup, err := trackedCommandFor(ctx, root,
		telemetryShell("(sleep 1; touch survived) & wait", false), "prompt")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	if err := cmd.Run(); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("timeout: %v", err)
	}
	time.Sleep(time.Second)
	if _, err := os.Stat(filepath.Join(root, "survived")); !os.IsNotExist(err) {
		t.Fatal("descendant survived timeout")
	}
	r := terminalRecords(t, root)[0]
	if r.StartedAt == nil || *r.Outcome.FailureClass != "timeout" {
		t.Fatalf("record: %+v", r)
	}
}

func TestTrackedCommandForStartEvidenceFailureReaps(t *testing.T) {
	root := t.TempDir()
	writeLaunchProtocol(t, root)
	ctx := WithLaunchInfo(context.Background(), LaunchInfo{Observe: func(r telemetry.Record) {
		if r.Type == "invocation.requested" {
			if err := os.Mkdir(filepath.Join(root, ".parley-runtime", "invocations", r.InvocationID, "started.json"), 0o700); err != nil {
				t.Fatal(err)
			}
		}
	}})
	cmd, cleanup, err := trackedCommandFor(ctx, root, telemetryShell("sleep 20", false), "prompt")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	if err := cmd.Run(); err == nil {
		t.Fatal("start evidence failure passed")
	}
	if cmd.ProcessState == nil {
		t.Fatal("process was not reaped")
	}
	r := terminalRecords(t, root)[0]
	if *r.Outcome.FailureClass != "telemetry_failure" {
		t.Fatalf("record: %+v", r)
	}
}

func TestTrackedCommandForAbandonedAndFailedStart(t *testing.T) {
	for _, start := range []bool{false, true} {
		root := t.TempDir()
		writeLaunchProtocol(t, root)
		agent := telemetryShell("exit 0", false)
		agent.Path = filepath.Join(root, "does-not-exist")
		cmd, cleanup, err := trackedCommandFor(context.Background(), root, agent, "prompt")
		if err != nil {
			t.Fatal(err)
		}
		if start && cmd.Start() == nil {
			t.Fatal("nonexistent executable started")
		}
		cleanup()
		r := terminalRecords(t, root)[0]
		if r.StartedAt != nil || r.Outcome.Status != "failed" {
			t.Fatalf("record: %+v", r)
		}
	}
}
