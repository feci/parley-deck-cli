package runner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/telemetry"
)

func noParticipantWait(t *testing.T) {
	t.Helper()
	old := participantRetryWait
	participantRetryWait = func(ctx context.Context, d time.Duration) error {
		if d <= 0 || d > 5*time.Second {
			t.Errorf("unexpected retry delay: %s", d)
		}
		return ctx.Err()
	}
	t.Cleanup(func() { participantRetryWait = old })
}

func TestDropoutRoundUsesRealChildrenAndSettlesOnce(t *testing.T) {
	noParticipantWait(t)
	for _, gate := range []bool{false, true} {
		t.Run(fmt.Sprint(gate), func(t *testing.T) {
			opts := quotaRunnerFixture(t, quota.NewParticipantPolicy(nil, nil))
			writeLaunchProtocol(t, opts.Root)
			if gate {
				p := filepath.Join(opts.Idea.Path, "00-prompt.md")
				raw, _ := os.ReadFile(p)
				raw = []byte(strings.Replace(string(raw), "status:", "auto_implement: true\nimplementer: a\nstatus:", 1))
				os.WriteFile(p, raw, 0600)
			}
			own := filepath.Join(opts.Idea.Path, "round-01", "c.md")
			script := "cat >/dev/null; printf 'invalid own output' > '" + strings.ReplaceAll(own, "'", "'\\''") + "'; exit 7"
			for _, id := range opts.Idea.Participants {
				opts.Agents = append(opts.Agents, agents.Discovery{Spec: agents.Spec{ID: id, HeadlessArgs: []string{"-c", script}, PromptMode: agents.PromptStdin}, Found: true, Path: "/bin/sh"})
			}
			results := RunRoundOne(context.Background(), opts)
			h, err := quota.ReadHistory(opts.Idea.Path)
			if err != nil {
				t.Fatal(err)
			}
			if gate {
				if h.Revision != 0 || !results[len(results)-1].QuotaBlocked {
					t.Fatalf("review gate bypassed: %+v %+v", h, results)
				}
			} else if h.Revision != 1 || !h.Dropped("c") {
				t.Fatalf("no reduction: %+v %+v", h, results)
			}
			if len(terminalRecords(t, opts.Root)) != 2 {
				t.Fatal("attempt cap")
			}
			again := RunRoundOne(context.Background(), opts)
			if len(terminalRecords(t, opts.Root)) != 2 {
				t.Fatalf("replay launched again: %+v", again)
			}
		})
	}
}

func TestDropoutExecWatchdogTimeoutAndACPShareTwoSlots(t *testing.T) {
	noParticipantWait(t)
	for _, kind := range []string{"timeout", "watchdog", "acp-started-exit"} {
		t.Run(kind, func(t *testing.T) {
			opts := quotaRunnerFixture(t, quota.NewParticipantPolicy(nil, nil))
			writeLaunchProtocol(t, opts.Root)
			opts.Timeout = 150 * time.Millisecond
			if kind == "watchdog" {
				opts.Timeout = 2 * time.Second
			}
			if kind == "acp-started-exit" {
				opts.Timeout = 5 * time.Second
			}
			for _, id := range opts.Idea.Participants {
				a := agents.Discovery{Spec: agents.Spec{ID: id, HeadlessArgs: []string{"-c", "exec sleep 5"}, PromptMode: agents.PromptStdin}, Found: true, Path: "/bin/sh"}
				if kind == "watchdog" {
					a.FirstEventTimeoutMS = 30
				}
				if kind == "acp-started-exit" {
					a.LaunchMode = agents.LaunchACP
					a.ACPArgs = []string{"-c", "exit 7"}
				}
				opts.Agents = append(opts.Agents, a)
			}
			results := RunRoundOne(context.Background(), opts)
			h, err := quota.ReadHistory(opts.Idea.Path)
			if err != nil || !h.Dropped("c") {
				t.Fatalf("%s: %+v %v %+v", kind, h, err, results)
			}
			records := terminalRecords(t, opts.Root)
			if len(records) != 2 {
				t.Fatalf("watchdog minted a third attempt: %d", len(records))
			}
			want := map[string]string{"timeout": "timeout", "watchdog": "no_first_output", "acp-started-exit": "process_failure"}[kind]
			for _, r := range records {
				if r.Outcome.FailureClass == nil || *r.Outcome.FailureClass != want {
					got := "<nil>"
					if r.Outcome.FailureClass != nil {
						got = *r.Outcome.FailureClass
					}
					t.Fatalf("%s: failure class=%q want=%q outcome=%+v", kind, got, want, r.Outcome)
				}
			}
		})
	}
}

func TestDropoutDurableStepBound(t *testing.T) {
	noParticipantWait(t)
	for _, successAt := range []int{0, 1, 2} {
		t.Run(fmt.Sprint(successAt), func(t *testing.T) {
			root := t.TempDir()
			writeLaunchProtocol(t, root)
			agent := telemetryShell("exit 7", false)
			valid, calls := false, 0
			opts := ParticipantStepOptions{Root: root, Idea: "idea", Agent: agent, Step: "round-01/own.md", Validate: func() StepValidation {
				return StepValidation{Valid: valid, Reason: "missing own artifact"}
			}}
			run := func(ctx context.Context, ordinal int, retry string) error {
				calls++
				valid = ordinal == successAt
				cmd, clean, err := CommandFor(ctx, root, agent, "input changes must not reset the budget")
				if clean != nil {
					defer clean()
				}
				if err != nil {
					return err
				}
				return cmd.Run()
			}
			r, err := RunParticipantStep(context.Background(), opts, run)
			want := 2
			if successAt == 1 {
				want = 1
			}
			if calls != want || r.Valid != (successAt > 0) || (err == nil) != (successAt > 0) {
				t.Fatalf("calls=%d result=%+v err=%v", calls, r, err)
			}
			if successAt == 0 && (r.Evidence == nil || len(r.Evidence.Failure.Attempts) != 2) {
				t.Fatalf("missing paired proof: %+v", r)
			}
			first := calls
			ctx := WithLaunchInfo(context.Background(), LaunchInfo{RunID: "different-run", Phase: "round-01"})
			again, _ := RunParticipantStep(ctx, opts, run)
			if calls != first || !again.Replayed || again.Valid != r.Valid {
				t.Fatalf("replay relaunched: calls=%d result=%+v", calls, again)
			}
			records := terminalRecords(t, root)
			if len(records) != want {
				t.Fatalf("records=%d", len(records))
			}
		})
	}
}

func TestDropoutRestartBetweenFailures(t *testing.T) {
	noParticipantWait(t)
	root := t.TempDir()
	writeLaunchProtocol(t, root)
	agent := telemetryShell("exit 9", false)
	opts := ParticipantStepOptions{Root: root, Idea: "restart", Agent: agent, Step: "round-01/own.md", Validate: func() StepValidation { return StepValidation{Reason: "missing"} }}
	run := func(ctx context.Context, _ int, _ string) error {
		cmd, clean, err := CommandFor(ctx, root, agent, "attempt")
		if clean != nil {
			defer clean()
		}
		if err != nil {
			return err
		}
		return cmd.Run()
	}
	participantRetryWait = func(context.Context, time.Duration) error { return errors.New("simulated driver exit before retry") }
	first, err := RunParticipantStep(context.Background(), opts, run)
	if err == nil || first.Evidence != nil || len(terminalRecords(t, root)) != 1 {
		t.Fatalf("%+v %v", first, err)
	}
	participantRetryWait = func(context.Context, time.Duration) error { return nil }
	second, err := RunParticipantStep(context.Background(), opts, run)
	if err == nil || second.Evidence == nil || len(terminalRecords(t, root)) != 2 {
		t.Fatalf("%+v %v", second, err)
	}
}

func TestDropoutCancellationAndIntegrityNeverReduce(t *testing.T) {
	noParticipantWait(t)
	for _, kind := range []string{"cancelled", "missing-protocol", "shared-file", "usage-event"} {
		t.Run(kind, func(t *testing.T) {
			root := t.TempDir()
			if kind != "missing-protocol" {
				writeLaunchProtocol(t, root)
			}
			agent := telemetryShell("exit 0", false)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			opts := ParticipantStepOptions{Root: root, Idea: "controls", Agent: agent, Step: "own", Validate: func() StepValidation {
				if kind == "shared-file" {
					return StepValidation{Integrity: errors.New("another author's block changed")}
				}
				return StepValidation{Valid: true}
			}}
			calls := 0
			run := func(c context.Context, _ int, _ string) error {
				calls++
				if kind == "cancelled" {
					cancel()
				}
				if kind == "usage-event" {
					// The terminal observer cannot make a later telemetry failure disappear.
					c = WithLaunchInfo(c, LaunchInfo{Observe: func(r telemetry.Record) {
						if r.Type == "invocation.terminal" {
							c.Value(participantStepKey{}).(*participantStepLaunch).integrity = errors.New("event publication failed")
						}
					}})
				}
				cmd, clean, err := CommandFor(c, root, agent, "attempt")
				if clean != nil {
					defer clean()
				}
				if err != nil {
					return err
				}
				return cmd.Run()
			}
			got, err := RunParticipantStep(ctx, opts, run)
			if err == nil || got.Evidence != nil || got.Valid || calls != 1 {
				t.Fatalf("%s: %+v %v calls=%d", kind, got, err, calls)
			}
			// A fresh parent must not replace valid output just because the old
			// terminal class was cancellation/refusal; integrity stays blocking.
			if kind != "missing-protocol" { // This refusal precedes the step ledger.
				again, replayErr := RunParticipantStep(context.Background(), opts, run)
				if replayErr == nil || again.Evidence != nil || again.Valid || calls != 1 {
					t.Fatalf("%s replay: %+v %v calls=%d", kind, again, replayErr, calls)
				}
			}
		})
	}
}

func TestDropoutCancelledChildWithValidArtifactNeverReplaced(t *testing.T) {
	for _, cancelParent := range []bool{false, true} {
		t.Run(fmt.Sprint(cancelParent), func(t *testing.T) {
			root := t.TempDir()
			writeLaunchProtocol(t, root)
			path := filepath.Join(root, "own.md")
			agent := telemetryShell("printf valid > '"+strings.ReplaceAll(path, "'", "'\\''")+"'; exec sleep 40", false)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			opts := ParticipantStepOptions{Root: root, Idea: "cancelled-output", Agent: agent, Step: "own", Files: []string{path}, Validate: func() StepValidation {
				raw, _ := os.ReadFile(path)
				return StepValidation{Valid: string(raw) == "valid"}
			}}
			calls := 0
			run := func(c context.Context, _ int, _ string) error {
				calls++
				stop := cancel
				if !cancelParent {
					c, stop = context.WithCancel(c)
					defer stop()
				}
				cmd, clean, err := CommandFor(c, root, agent, "attempt")
				if clean != nil {
					defer clean()
				}
				if err != nil {
					return err
				}
				if err := cmd.Start(); err != nil {
					return err
				}
				deadline := time.Now().Add(5 * time.Second)
				for !opts.Validate().Valid && time.Now().Before(deadline) {
					time.Sleep(10 * time.Millisecond)
				}
				stop()
				return cmd.Wait()
			}
			first, err := RunParticipantStep(ctx, opts, run)
			if err == nil || first.Valid || first.Evidence != nil || !opts.Validate().Valid {
				t.Fatalf("cancelled=%+v %v", first, err)
			}
			for i := 0; i < 2; i++ {
				again, err := RunParticipantStep(context.Background(), opts, run)
				if err == nil || again.Valid || again.Evidence != nil || calls != 1 {
					t.Fatalf("cancelled replay=%+v %v calls=%d", again, err, calls)
				}
			}
		})
	}
}

func TestDropoutRepairedRefusalWithoutValidOutputMayDispatch(t *testing.T) {
	root := t.TempDir()
	agent := telemetryShell("exit 0", false)
	valid, calls := false, 0
	opts := ParticipantStepOptions{Root: root, Idea: "repair", Agent: agent, Step: "own", Validate: func() StepValidation {
		return StepValidation{Valid: valid, Reason: "missing own artifact"}
	}}
	run := func(ctx context.Context, ordinal int, retry string) error {
		calls++
		if ordinal != 1 || retry != "" {
			t.Error("undispatched refusal consumed child attempt", ordinal, retry)
		}
		cmd, clean, err := CommandFor(ctx, root, agent, "attempt")
		if clean != nil {
			defer clean()
		}
		if err != nil {
			return err
		}
		err = cmd.Run()
		valid = err == nil
		return err
	}
	if got, err := RunParticipantStep(context.Background(), opts, run); err == nil || got.Valid || got.Evidence != nil {
		t.Fatalf("refusal=%+v %v", got, err)
	}
	writeLaunchProtocol(t, root)
	if got, err := RunParticipantStep(context.Background(), opts, run); err != nil || !got.Valid || got.Evidence != nil || calls != 2 {
		t.Fatalf("repair=%+v %v calls=%d", got, err, calls)
	}
}

func TestDropoutConcurrentAndPrivatePartial(t *testing.T) {
	noParticipantWait(t)
	root := t.TempDir()
	writeLaunchProtocol(t, root)
	agent := telemetryShell("exit 1", false)
	file := filepath.Join(root, "partial.md")
	os.WriteFile(file, []byte("original dissent"), 0600)
	opts := ParticipantStepOptions{Root: root, Idea: "same", Agent: agent, Step: "own", Files: []string{file}, Validate: func() StepValidation { return StepValidation{Reason: "invalid own artifact"} }}
	entered, release, done := make(chan struct{}), make(chan struct{}), make(chan struct{})
	var calls atomic.Int32
	run := func(ctx context.Context, ordinal int, _ string) error {
		calls.Add(1)
		if ordinal == 1 {
			close(entered)
			<-release
		}
		os.WriteFile(file, []byte(fmt.Sprintf("invalid %d", ordinal)), 0600)
		cmd, clean, err := CommandFor(ctx, root, agent, "attempt")
		if clean != nil {
			defer clean()
		}
		if err != nil {
			return err
		}
		return cmd.Run()
	}
	go func() { defer close(done); RunParticipantStep(context.Background(), opts, run) }()
	<-entered
	if _, err := RunParticipantStep(context.Background(), opts, run); err == nil {
		t.Fatal("concurrent step admitted")
	}
	close(release)
	<-done
	if calls.Load() != 2 {
		t.Fatal(calls.Load())
	}
	saved, _ := filepath.Glob(filepath.Join(root, ".parley-runtime", "invocations", "*", "participant-files", "00-partial.md"))
	if len(saved) != 2 {
		t.Fatalf("partial evidence=%v", saved)
	}
	for _, p := range saved {
		st, _ := os.Stat(p)
		if st.Mode().Perm() != 0600 {
			t.Fatalf("non-private evidence %s", p)
		}
	}
}
