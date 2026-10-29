package runner

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/membership"
	"parley-deck-cli/internal/quota"
)

func exhaustionChild(t *testing.T, root string) agents.Discovery {
	t.Helper()
	raw, err := os.ReadFile("../telemetry/testdata/quota/zcode-sdk-source-derived.stderr")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "provider-fixture.stderr")
	if err := os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	a := telemetryShell("cat >/dev/null; cat '"+strings.ReplaceAll(path, "'", "'\\''")+"' >&2; exit 7", false)
	a.AdapterID = "zcode"
	return a
}

func TestExhaustionOneChildAndDurableAdmission(t *testing.T) {
	root := t.TempDir()
	writeLaunchProtocol(t, root)
	agent := exhaustionChild(t, root)
	opts := ParticipantStepOptions{Root: root, Idea: "idea", Agent: agent, Step: "round-01/own.md", Exhaustion: true,
		Validate: func() StepValidation { return StepValidation{Reason: "missing own output"} }}
	calls := 0
	run := func(ctx context.Context, _ int, _ string) error {
		calls++
		cmd, cleanup, err := CommandFor(ctx, root, agent, "fixture")
		if cleanup != nil {
			defer cleanup()
		}
		if err != nil {
			return err
		}
		return cmd.Run()
	}
	r, err := RunParticipantStep(context.Background(), opts, run)
	if err == nil || calls != 1 || r.Evidence == nil || r.Evidence.RuleID != quota.ProviderExhaustionRule || quota.ValidateExhaustionEvidence(*r.Evidence, "idea", agent.ID) != nil {
		t.Fatalf("calls=%d result=%+v err=%v", calls, r, err)
	}
	if r.Evidence.ResetAt != nil {
		t.Fatal("past fixture reset did not become unknown")
	}
	// Replay and a downgraded/new logical step must not reconnect.
	opts.Exhaustion = false
	again, err := RunParticipantStep(context.Background(), opts, run)
	if err == nil || !again.Replayed || calls != 1 {
		t.Fatalf("%+v %v calls=%d", again, err, calls)
	}
	ctx := WithLaunchInfo(context.Background(), LaunchInfo{RunID: "changed-run", Idea: "idea", Phase: "consult"})
	if _, err := beginLaunch(ctx, root, "changed-run", agent); err == nil || !strings.Contains(err.Error(), "permanently suppresses") {
		t.Fatal(err)
	}
	if len(terminalRecords(t, root)) != 1 {
		t.Fatal("admission launched another child")
	}
	if _, err := membership.SupplementUsability(filepath.Join(root, "parley-deck", "ideas", "idea"), []string{agent.ID}, []quota.Member{{ID: agent.ID, ValidArtifact: true}}); err == nil {
		t.Fatal("old valid artifact restored a suppressed identity's floor seat")
	}
}

func TestExhaustionStaleImplementationDoesNotWin(t *testing.T) {
	root := t.TempDir()
	writeLaunchProtocol(t, root)
	agent := exhaustionChild(t, root)
	path := filepath.Join(root, "IMPLEMENTATION.md")
	os.WriteFile(path, []byte("---\nidea: idea\nstatus: implemented\n---\n\n## Summary of work\nOld completed work.\n"), 0600)
	opts := ParticipantStepOptions{Root: root, Idea: "idea", Agent: agent, Step: "implementation", Exhaustion: true, SingleAttempt: true, Files: []string{path},
		Validate: func() StepValidation { return freshExecutionArtifact(root, "idea", agent.ID, "implementation", path) }}
	r, err := RunParticipantStep(context.Background(), opts, func(ctx context.Context, _ int, _ string) error {
		cmd, cleanup, err := CommandFor(ctx, root, agent, "fixture")
		if cleanup != nil {
			defer cleanup()
		}
		if err != nil {
			return err
		}
		return cmd.Run()
	})
	if err == nil || r.Valid || r.Evidence == nil {
		t.Fatalf("stale artifact masked failure: %+v %v", r, err)
	}
}

func TestExhaustionTerminalValidationRecoveryNoChild(t *testing.T) {
	root := t.TempDir()
	writeLaunchProtocol(t, root)
	agent := exhaustionChild(t, root)
	opts := ParticipantStepOptions{Root: root, Idea: "idea", Agent: agent, Step: "readiness", Exhaustion: true,
		Validate: func() StepValidation { return StepValidation{Reason: "not PONG"} }}
	r, _ := RunParticipantStep(context.Background(), opts, func(ctx context.Context, _ int, _ string) error {
		cmd, cleanup, err := CommandFor(ctx, root, agent, "fixture")
		if cleanup != nil {
			defer cleanup()
		}
		if err != nil {
			return err
		}
		return cmd.Run()
	})
	if r.Evidence == nil {
		t.Fatal("fixture did not exhaust")
	}
	path := filepath.Join(root, ".parley-runtime", "invocations", r.Record.InvocationID, "participant-result.json")
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := exhaustionAdmission(root, "idea", agent.ID); err == nil || !strings.Contains(err.Error(), "validation recovery") {
		t.Fatal(err)
	}
	again, err := RunParticipantStep(context.Background(), opts, func(context.Context, int, string) error { t.Fatal("recovery launched a child"); return nil })
	if err == nil || again.Evidence == nil || !again.Replayed {
		t.Fatalf("%+v %v", again, err)
	}
	if err := os.WriteFile(path, []byte("{}"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := exhaustionAdmission(root, "idea", agent.ID); err == nil {
		t.Fatal("corrupt receipt admitted a child")
	}
}

func TestExhaustionValidOutputAndSingleAttemptTransient(t *testing.T) {
	for _, valid := range []bool{false, true} {
		root := t.TempDir()
		writeLaunchProtocol(t, root)
		agent := exhaustionChild(t, root)
		if !valid {
			agent = telemetryShell("cat >/dev/null; exit 9", false)
		}
		opts := ParticipantStepOptions{Root: root, Idea: "idea", Agent: agent, Step: "implementation", Exhaustion: true, SingleAttempt: true,
			Validate: func() StepValidation { return StepValidation{Valid: valid, Reason: "missing"} }}
		calls := 0
		r, err := RunParticipantStep(context.Background(), opts, func(ctx context.Context, _ int, _ string) error {
			calls++
			cmd, cleanup, err := CommandFor(ctx, root, agent, "fixture")
			if cleanup != nil {
				defer cleanup()
			}
			if err != nil {
				return err
			}
			return cmd.Run()
		})
		if calls != 1 || r.Evidence != nil || r.Valid != valid || (err == nil) != valid {
			t.Fatalf("%+v %v calls=%d", r, err, calls)
		}
		if err := exhaustionAdmission(root, "idea", agent.ID); err != nil {
			t.Fatal(err)
		}
	}
}
