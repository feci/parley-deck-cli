package runner

import (
	"context"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
)

func TestReviewGateScopedSupervisionDefaults(t *testing.T) {
	step := context.WithValue(context.Background(), participantStepKey{}, &participantStepLaunch{})
	for _, test := range []struct {
		name                    string
		ctx                     context.Context
		spec                    agents.Spec
		hard                    time.Duration
		first, stall, heartbeat time.Duration
	}{
		{"step", step, agents.Spec{}, time.Hour, 120 * time.Second, 300 * time.Second, 60 * time.Second},
		{"ordinary", context.Background(), agents.Spec{}, time.Hour, 120 * time.Second, 1800 * time.Second, 60 * time.Second},
		{"override", step, agents.Spec{FirstEventTimeoutMS: 20, StallTimeoutMS: 30, HeartbeatMS: 40}, time.Hour, 20 * time.Millisecond, 30 * time.Millisecond, 40 * time.Millisecond},
		{"disabled", step, agents.Spec{FirstEventTimeoutMS: -1, StallTimeoutMS: -1, HeartbeatMS: -1}, time.Hour, 0, 0, 0},
		{"buffered", step, agents.Spec{BuffersStdout: true}, time.Hour, 0, 0, 60 * time.Second},
		{"short-hard", step, agents.Spec{}, 90 * time.Second, 120 * time.Second, 89 * time.Second, 60 * time.Second},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := supervisionForStep(test.ctx, agents.Discovery{Spec: test.spec}, test.hard)
			if got.FirstEventTimeout != test.first || got.StallTimeout != test.stall || got.HeartbeatInterval != test.heartbeat {
				t.Fatalf("%+v", got)
			}
		})
	}
}

func TestReviewGateSupervisedCommandTerminalAndReplay(t *testing.T) {
	noParticipantWait(t)
	for _, kind := range []string{"no_first_output", "stalled", "timeout", "buffered-success", "disabled"} {
		t.Run(kind, func(t *testing.T) {
			root := t.TempDir()
			writeLaunchProtocol(t, root)
			script := "cat >/dev/null; exec sleep 30"
			if kind == "stalled" {
				script = "cat >/dev/null; printf 'ready'; exec sleep 30"
			}
			if kind == "buffered-success" {
				script = "cat >/dev/null; sleep 2; printf 'done'"
			}
			agent := telemetryShell(script, false)
			agent.FirstEventTimeoutMS, agent.StallTimeoutMS, agent.HeartbeatMS = 50, 50, 50
			hard := 5 * time.Second
			if kind == "timeout" {
				agent.BuffersStdout = true
				hard = 200 * time.Millisecond
			}
			if kind == "buffered-success" {
				agent.BuffersStdout = true
			}
			if kind == "disabled" {
				agent.FirstEventTimeoutMS, agent.StallTimeoutMS = -1, -1
				hard = 200 * time.Millisecond
			}
			valid, calls := false, 0
			opts := ParticipantStepOptions{Root: root, Idea: "review-gate-command", Agent: agent, Step: "consensus", Validate: func() StepValidation { return StepValidation{Valid: valid, Reason: "no signoff"} }}
			run := func(parent context.Context, _ int, _ string) error {
				calls++
				ctx, cancel := context.WithTimeout(parent, hard)
				defer cancel()
				cmd, clean, err := CommandFor(ctx, root, agent, "signoff fixture")
				if clean != nil {
					defer clean()
				}
				if err != nil {
					return err
				}
				err = cmd.RunSupervised(agent, hard)
				valid = err == nil
				return err
			}
			result, err := RunParticipantStep(context.Background(), opts, run)
			want := 2
			if kind == "buffered-success" {
				want = 1
			}
			if calls != want || result.Valid != (kind == "buffered-success") || (err == nil) != (kind == "buffered-success") {
				t.Fatalf("calls=%d result=%+v err=%v", calls, result, err)
			}
			for _, r := range terminalRecords(t, root) {
				if kind == "buffered-success" {
					if r.Outcome.FailureClass != nil {
						t.Fatalf("buffered success classified failed: %+v", r.Outcome)
					}
					continue
				}
				class := kind
				if kind == "disabled" {
					class = "timeout"
				}
				if r.Outcome.FailureClass == nil || *r.Outcome.FailureClass != class {
					t.Fatalf("want %s: %+v", class, r.Outcome)
				}
			}
			if kind != "buffered-success" && (result.Evidence == nil || len(result.Evidence.Failure.Attempts) != 2) {
				t.Fatal("no paired watchdog evidence")
			}
			_, _ = RunParticipantStep(context.Background(), opts, run)
			if calls != want {
				t.Fatal("watchdog replay minted another attempt")
			}
		})
	}
}
