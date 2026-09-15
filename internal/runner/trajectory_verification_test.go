package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/trajectory"
)

func verifierLaunchFixture(t *testing.T) (trajectory.VerificationTicket, string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("POSIX captured verification")
	}
	root, idea := trajectoryRuntimeFixture(t)
	builder := telemetryShell("printf 'broken\\n' > source; exit 7", false)
	if record, err := RunMeasured(context.Background(), ExecOptions{Root: root, Agent: builder, Prompt: "synthetic patch", Timeout: 30 * time.Second, Info: LaunchInfo{RunID: "patch-fixture", Idea: idea.Slug, Phase: "fixup"}}); err == nil || record.StartedAt == nil {
		t.Fatalf("charged patch did not execute: %+v %v", record, err)
	}
	ticket, err := trajectory.PrepareCapturedVerification(context.Background(), root, idea.Slug, "reviewer", "verification-fixture")
	if err != nil {
		t.Fatal(err)
	}
	b, err := budget.LoadCycleBinding(context.Background(), root, idea.Slug, budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	return ticket, filepath.Join(filepath.Dir(b.Store.Dir), "trajectory-verifications", ticket.Request.Charge.EntryKey)
}

func TestCapturedVerifierLaunchReservesBeforeSpawnAndRefusesReplay(t *testing.T) {
	ticket, journal := verifierLaunchFixture(t)
	ctx, err := WithCapturedVerification(context.Background(), ticket)
	if err != nil {
		t.Fatal(err)
	}
	// Mutating caller-owned criteria after binding cannot rewrite the ticket.
	ticket.Request.Criteria[0].CommandSHA256 = strings.Repeat("0", 64)
	quote := "'" + strings.ReplaceAll(filepath.Join(journal, "launch.json"), "'", "'\\''") + "'"
	agent := telemetryShell("test -f "+quote+" || exit 9; printf '%s\\n' \"$PARLEY_PROC_MARKER\" >> .parley-runtime/verifier-starts", false)
	agent.ID = "reviewer"
	options := ExecOptions{Root: ticket.Root, Agent: agent, Prompt: "synthetic verifier", Timeout: 30 * time.Second,
		Info: LaunchInfo{RunID: ticket.RunID, Idea: ticket.Request.Idea, Phase: CapturedVerificationPhase}}
	record, err := RunMeasured(ctx, options)
	if err != nil || record.StartedAt == nil {
		t.Fatalf("bound verifier did not execute: %+v %v", record, err)
	}
	launch, err := os.ReadFile(filepath.Join(journal, "launch.json"))
	if err != nil {
		t.Fatal(err)
	}
	var reserved struct {
		InvocationID string `json:"invocation_id"`
	}
	if err = json.Unmarshal(launch, &reserved); err != nil || reserved.InvocationID != record.InvocationID {
		t.Fatal("launch reservation not bound to actual telemetry", err)
	}
	if denied, err := RunMeasured(ctx, options); err == nil || denied.StartedAt != nil || denied.Outcome == nil {
		t.Fatalf("ticket silently replayed: %+v %v", denied, err)
	}
	after, err := os.ReadFile(filepath.Join(journal, "launch.json"))
	if err != nil || !bytes.Equal(launch, after) {
		t.Fatal("replay changed the original invocation reservation")
	}
	starts, err := os.ReadFile(filepath.Join(ticket.Root, ".parley-runtime/verifier-starts"))
	if err != nil || string(starts) != record.InvocationID+"\n" {
		t.Fatalf("expected exactly one observed verifier: %s %v", starts, err)
	}
}

func TestCapturedVerifierLaunchRejectsChangedBinding(t *testing.T) {
	for _, kind := range []string{"run", "idea", "phase", "agent", "root", "missing-ticket", "publication", "handoff"} {
		t.Run(kind, func(t *testing.T) {
			ticket, journal := verifierLaunchFixture(t)
			ctx, err := WithCapturedVerification(context.Background(), ticket)
			if err != nil {
				t.Fatal(err)
			}
			agent := telemetryShell("printf started > .parley-runtime/verifier-starts", false)
			agent.ID = "reviewer"
			info := LaunchInfo{RunID: ticket.RunID, Idea: ticket.Request.Idea, Phase: CapturedVerificationPhase}
			root := ticket.Root
			switch kind {
			case "run":
				info.RunID = "another-run"
			case "idea":
				info.Idea = "another-idea"
			case "phase":
				info.Phase = "review"
			case "agent":
				agent.ID = "test-1"
			case "root":
				root = t.TempDir()
				if err = protocol.InitWorkspace(root); err != nil {
					t.Fatal(err)
				}
				declareTestLaunchSource(t, root)
			case "missing-ticket":
				ctx = context.Background()
			case "publication":
				if err = os.Mkdir(filepath.Join(journal, "launch.json"), 0700); err != nil {
					t.Fatal(err)
				}
			}
			if kind == "handoff" {
				if _, err = beginLaunch(WithLaunchInfo(ctx, info), root, ticket.RunID, agent, launchHandoff); err == nil {
					t.Fatal("unobserved handoff reserved a verifier execution")
				}
			} else if record, err := RunMeasured(ctx, ExecOptions{Root: root, Agent: agent, Prompt: "synthetic verifier", Timeout: 30 * time.Second, Info: info}); err == nil || record.StartedAt != nil || record.Outcome == nil {
				t.Fatalf("changed binding launched: %+v %v", record, err)
			}
			if _, err := os.Stat(filepath.Join(root, ".parley-runtime/verifier-starts")); !os.IsNotExist(err) {
				t.Fatal("refused verifier actually started")
			}
			if kind != "publication" {
				if _, err := os.Stat(filepath.Join(journal, "launch.json")); !os.IsNotExist(err) {
					t.Fatal("changed binding consumed the prepared ticket")
				}
			}
		})
	}
}
