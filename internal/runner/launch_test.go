package runner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/telemetry"
)

func interactiveFixture(t *testing.T, root, script string) (agents.Discovery, *os.File) {
	t.Helper()
	file, err := os.CreateTemp(root, "terminal-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { file.Close() })
	agent := telemetryShell(script, false)
	agent.LaunchMode = agents.LaunchInteractive
	agent.InteractiveCommand = "/bin/sh"
	agent.InteractivePromptMode = agents.InteractivePromptFile
	agent.InteractiveArgs = []string{"-c", script, "fixture", "{prompt_path}"}
	return agent, file
}

func TestInteractiveProcessHasDistinctFreshEvidence(t *testing.T) {
	root := t.TempDir()
	writeLaunchProtocol(t, root)
	agent, terminal := interactiveFixture(t, root, `cat "$1"; printf '\nchild output\n'`)
	packet, err := WriteHandoffPacket(HandoffOptions{Root: root, RunID: "interactive", Agent: agent, Prompt: "old task"})
	if err != nil {
		t.Fatal(err)
	}
	// A printed handoff is not the prompt authority of a later process launch.
	if err := os.WriteFile(packet.PromptPath, []byte("tampered handoff"), 0o600); err != nil {
		t.Fatal(err)
	}
	ctx := WithLaunchInfo(context.Background(), LaunchInfo{RunID: "interactive", Phase: "consensus"})
	if err := RunInteractive(ctx, root, agent, "current task", "", terminal, terminal, terminal); err != nil {
		t.Fatal(err)
	}
	output, err := os.ReadFile(terminal.Name())
	if err != nil || !strings.Contains(string(output), "Mandatory source obligation.") || !strings.Contains(string(output), "current task") || strings.Contains(string(output), "tampered handoff") {
		t.Fatalf("process prompt: %s; %v", output, err)
	}
	records := terminalRecords(t, root)
	if len(records) != 2 {
		t.Fatalf("expected handoff plus process, got %d", len(records))
	}
	for _, r := range records {
		if r.InvocationID == packet.InvocationID {
			if r.StartedAt != nil || r.Outcome.Status != "unobserved-handoff" {
				t.Fatalf("handoff claimed execution: %+v", r)
			}
			continue
		}
		if r.StartedAt == nil || r.PID == nil || r.Outcome.ExitCode == nil || *r.Outcome.ExitCode != 0 || r.Outcome.Status != "process-exited" || r.Metadata.Context.Mode != "full" {
			t.Fatalf("missing process evidence: %+v", r)
		}
		if r.Outcome.Observation.StdoutBytes != nil || r.Outcome.Observation.StderrBytes != nil || r.Outcome.Observation.StreamCoverage != "not-observed-terminal" || r.Outcome.Observation.FirstActivityMS != nil || r.Outcome.Usage.CostUSD != nil || r.Outcome.Usage.ReportedModel != nil {
			t.Fatalf("invented terminal observations: %+v", r.Outcome)
		}
	}
}

func TestInteractiveProcessFailureEvidence(t *testing.T) {
	for _, scenario := range []string{"missing-delivery", "missing-authority", "failed-start", "failed-exit", "timeout", "start-write", "terminal-write"} {
		t.Run(scenario, func(t *testing.T) {
			root := t.TempDir()
			if scenario != "missing-authority" {
				writeLaunchProtocol(t, root)
			}
			script := "exit 7"
			if scenario == "timeout" || scenario == "start-write" {
				script = "exec sleep 20"
			}
			agent, terminal := interactiveFixture(t, root, script)
			if scenario == "missing-delivery" {
				agent.InteractivePromptMode = agents.InteractivePromptNone
			}
			if scenario == "failed-start" {
				agent.InteractiveCommand = filepath.Join(root, "missing")
			}
			ctx := WithLaunchInfo(context.Background(), LaunchInfo{Observe: func(r telemetry.Record) {
				if r.Type == "invocation.requested" && (scenario == "start-write" || scenario == "terminal-write") {
					name := "started.json"
					if scenario == "terminal-write" {
						name = "terminal.json"
					}
					if err := os.Mkdir(filepath.Join(root, ".parley-runtime", "invocations", r.InvocationID, name), 0o700); err != nil {
						t.Fatal(err)
					}
				}
			}})
			if scenario == "timeout" {
				agent.InteractiveTimeoutMS = 100
			}
			if err := RunInteractive(ctx, root, agent, "task", "", terminal, terminal, terminal); err == nil {
				t.Fatal("invalid process accepted")
			}
			if scenario == "terminal-write" {
				return
			}
			records := terminalRecords(t, root)
			if len(records) != 1 || records[0].Outcome.Status != "failed" {
				t.Fatalf("missing failure: %+v", records)
			}
			want := map[string]string{"missing-delivery": "protocol_context_refused", "missing-authority": "protocol_context_refused", "failed-start": "start_failure", "failed-exit": "process_failure", "timeout": "timeout", "start-write": "telemetry_failure"}[scenario]
			if got := records[0].Outcome.FailureClass; got == nil || *got != want {
				t.Fatalf("failure class: %v, want %s", got, want)
			}
		})
	}
}

func TestInteractiveTimeoutKillsDescendants(t *testing.T) {
	root := t.TempDir()
	writeLaunchProtocol(t, root)
	agent, terminal := interactiveFixture(t, root, `(sleep 2; touch survived) & wait`)
	agent.InteractiveTimeoutMS = 150
	if err := RunInteractive(context.Background(), root, agent, "task", "", terminal, terminal, terminal); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("timeout did not fail: %v", err)
	}
	time.Sleep(2100 * time.Millisecond)
	if _, err := os.Stat(filepath.Join(root, "survived")); !os.IsNotExist(err) {
		t.Fatal("interactive descendant survived cancellation")
	}
}

// Run with PARLEY_TTY_TEST=1 under a real controlling PTY. The external harness
// supplies two lines and a hard deadline: both the child and the restored parent
// must be able to read the terminal, and all child descriptors must remain TTYs.
func TestInteractiveRealTerminal(t *testing.T) {
	if os.Getenv("PARLEY_TTY_TEST") != "1" {
		t.Skip("requires a controlling PTY harness")
	}
	root := t.TempDir()
	writeLaunchProtocol(t, root)
	agent, _ := interactiveFixture(t, root, `test -t 0 && test -t 1 && test -t 2 && IFS= read -r answer && test "$answer" = child`)
	agent.InteractiveTimeoutMS = 5000
	input, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer input.Close()
	fmt.Println("TTY_CHILD_READY")
	if err := RunInteractive(context.Background(), root, agent, "task", "", input, os.Stdout, os.Stderr); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", `IFS= read -r answer && test "$answer" = parent`)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	fmt.Println("TTY_PARENT_READY")
	if err := cmd.Run(); err != nil {
		t.Fatalf("parent terminal not restored: %v", err)
	}
	r := terminalRecords(t, root)[0]
	if r.StartedAt == nil || r.Outcome.Observation.StreamCoverage != "not-observed-terminal" {
		t.Fatalf("missing terminal evidence: %+v", r)
	}
}

func TestInteractiveRealTerminalFailures(t *testing.T) {
	if os.Getenv("PARLEY_TTY_TEST") != "1" {
		t.Skip("requires a controlling PTY harness")
	}
	for _, scenario := range []string{"timeout", "start-write", "failed-start", "terminal-write"} {
		t.Run(scenario, func(t *testing.T) {
			root := t.TempDir()
			writeLaunchProtocol(t, root)
			agent, _ := interactiveFixture(t, root, `(sleep 2; touch survived) & wait`)
			agent.InteractiveTimeoutMS = 150
			if scenario == "failed-start" {
				agent.InteractiveCommand = filepath.Join(root, "missing")
			}
			ctx := WithLaunchInfo(context.Background(), LaunchInfo{Observe: func(r telemetry.Record) {
				if r.Type != "invocation.requested" {
					return
				}
				name := ""
				if scenario == "start-write" {
					name = "started.json"
				}
				if scenario == "terminal-write" {
					name = "terminal.json"
				}
				if name != "" {
					if err := os.Mkdir(filepath.Join(root, ".parley-runtime", "invocations", r.InvocationID, name), 0o700); err != nil {
						t.Fatal(err)
					}
				}
			}})
			if err := RunInteractive(ctx, root, agent, "task", "", os.Stdin, os.Stdout, os.Stderr); err == nil {
				t.Fatal("failure accepted")
			}
			parentCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			cmd := exec.CommandContext(parentCtx, "/bin/sh", "-c", `IFS= read -r answer && test "$answer" = parent`)
			cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
			fmt.Println("TTY_PARENT_READY")
			if err := cmd.Run(); err != nil {
				t.Fatalf("parent foreground not restored: %v", err)
			}
			time.Sleep(2100 * time.Millisecond)
			if _, err := os.Stat(filepath.Join(root, "survived")); !os.IsNotExist(err) {
				t.Fatal("terminal descendant survived failure")
			}
		})
	}
}

func TestInteractiveArgumentDeliveryAndBound(t *testing.T) {
	for _, large := range []bool{false, true} {
		root := t.TempDir()
		writeLaunchProtocol(t, root)
		agent, out := interactiveFixture(t, root, `printf '%s' "$1"`)
		agent.InteractivePromptMode = agents.InteractivePromptArg
		agent.InteractiveArgs = []string{"-c", `printf '%s' "$1"`, "fixture", "{prompt}"}
		prompt := "argument task"
		if large {
			prompt = strings.Repeat("x", 121<<10)
		}
		err := RunInteractive(context.Background(), root, agent, prompt, "", out, out, out)
		if large {
			if err == nil || !strings.Contains(err.Error(), "too-large-use-file") {
				t.Fatalf("oversized argument: %v", err)
			}
			if terminalRecords(t, root)[0].StartedAt != nil {
				t.Fatal("oversized argument spawned")
			}
		} else {
			if err != nil {
				t.Fatal(err)
			}
			body, err := os.ReadFile(out.Name())
			if err != nil || !strings.Contains(string(body), "Mandatory source obligation.") || !strings.HasSuffix(string(body), prompt) {
				t.Fatalf("arg prompt not delivered: %v", err)
			}
		}
	}
}

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
		if len(records) != 1 || (records[0].Outcome.Observation.StdoutBytes == nil || *records[0].Outcome.Observation.StdoutBytes != 5) || (records[0].Outcome.Observation.StderrBytes == nil || *records[0].Outcome.Observation.StderrBytes != 5) {
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
