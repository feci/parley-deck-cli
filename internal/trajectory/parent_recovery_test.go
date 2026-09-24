package trajectory

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/telemetry"
)

func TestParentRecoveryRejectsContraryPartialPublication(t *testing.T) {
	want := ParentResult{Version: 1, RunID: "original-run", InvocationID: "original-invocation", TrajectoryPending: true, Assessment: &Assessment{Outcome: Regression}}
	for _, raw := range []string{
		`{"run_id":"other",`, `{"failure_stage":"acceptance",`,
		`{"trajectory_pending":false,`, `{"trajectory_pending":null,`,
		`{"run_id":null,`, `{"run_id":"other`,
		`{"run_id":"original-run","run_id":`, `{"unknown":`,
		`{"assessment":{"outcome":"no-regression"},`,
		`{"assessment":{"outcome":"no-regression",`,
		`{"version":"1",`, `{"version":0,`, `{"version":1, garbage`,
		`{"version":1} {}`, `[]`, `garbage`,
	} {
		t.Run(raw, func(t *testing.T) {
			if _, _, err := decodeUnpublishedParent([]byte(raw), want); err == nil {
				t.Fatal("conflicting or ambiguous partial publication was accepted")
			}
		})
	}
	for _, raw := range []string{"", "{", `{"ver`, `{"version":1,"run_`, `{"version":1,"run_id":`, `{"run_id":"orig`, `{"run_id":"original-run",`, `{"failure_stage":"parent-publication",`} {
		if _, complete, err := decodeUnpublishedParent([]byte(raw), want); err != nil || complete {
			t.Fatalf("compatible partial publication %q: complete=%v err=%v", raw, complete, err)
		}
	}
}

type recoveryProcessInput struct {
	Ticket             VerificationTicket
	Criteria           []Criterion
	Invocation, Parent string
}

func TestParentRecoveryProcessHelper(t *testing.T) {
	path := os.Getenv("PARLEY_PARENT_RECOVERY_PROCESS_INPUT")
	if path == "" {
		return
	}
	var input recoveryProcessInput
	raw, err := os.ReadFile(path)
	if err == nil {
		err = json.Unmarshal(raw, &input)
	}
	if err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		_, err = ExecuteCapturedVerification(ctx, input.Ticket, input.Invocation, input.Criteria, input.Parent)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(19)
	}
	os.Exit(0)
}

func parentRecoveryFixture(t *testing.T) (string, ParentRecoveryPreview) {
	t.Helper()
	criteria := capturedCriteria()
	ticket, _ := journalFixture(t, criteria, dirtyCapturedChild(t))
	root := ticket.Root
	scope := "---\nparticipants: [builder, reviewer]\nchecks:\n  - name: material\n    command: >\n      " + criteria[0].Command + "\n---\n"
	snapshotWrite(t, root, "parley-deck/ideas/fixture/00-prompt.md", []byte(scope), 0600)
	base := filepath.Join(".parley-runtime", "trajectory-verification", ticket.RunID)
	req, _ := canonical(HelperRequest{1, ticket, []string{"builder", "reviewer"}, criteria})
	snapshotWrite(t, root, filepath.Join(base, "request.json"), req, 0600)
	inv, err := telemetry.Begin(filepath.Join(root, ".parley-runtime", "invocations"), telemetry.Metadata{RunID: ticket.RunID, Idea: "fixture", Phase: "trajectory-verification", Agent: "reviewer", LaunchMode: "headless"})
	if err != nil {
		t.Fatal(err)
	}
	if err = ReserveCapturedVerificationLaunch(context.Background(), ticket, inv.ID); err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal(recoveryProcessInput{ticket, criteria, inv.ID, t.TempDir()})
	input := filepath.Join(t.TempDir(), "request.json")
	if err = os.WriteFile(input, raw, 0600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestParentRecoveryProcessHelper$")
	cmd.Env = append(os.Environ(), "PARLEY_PARENT_RECOVERY_PROCESS_INPUT="+input)
	var out bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &out
	if err = cmd.Start(); err != nil {
		t.Fatal(err)
	}
	startErr := inv.Started(cmd.Process.Pid)
	if err = cmd.Wait(); err != nil {
		dumpAllRetainedVerifications(t, root)
		t.Fatalf("helper: %v %s", err, &out)
	}
	if startErr != nil {
		t.Fatal(startErr)
	}
	zero := 0
	if err = inv.Finish(telemetry.Outcome{Status: "process-exited", ExitCode: &zero}); err != nil {
		t.Fatal(err)
	}
	p, err := PreviewParentRecovery(context.Background(), root, "fixture", ticket.RunID)
	if err != nil {
		dumpAllRetainedVerifications(t, root)
		t.Fatal(err)
	}
	return root, p
}

func TestParentRecoveryPublicationInterruptionAndExactReplay(t *testing.T) {
	for _, phase := range []string{"before-publication", "rename", "persistence-barrier", "after-publication"} {
		t.Run(phase, func(t *testing.T) {
			root, p := parentRecoveryFixture(t)
			base := filepath.Join(".parley-runtime", "trajectory-verification", p.RunID)
			path := filepath.Join(root, base, parentRecoveryName)
			calls := 0
			fault := errors.New("injected parent publication interruption")
			publish := func(dir *os.Root, base string, r ParentRecovery) error {
				calls++
				switch phase {
				case "before-publication":
					return fault
				case "rename":
					if err := dir.Mkdir(filepath.Join(base, parentRecoveryName), 0700); err != nil {
						return err
					}
					defer dir.Remove(filepath.Join(base, parentRecoveryName))
					return publishParentRecovery(dir, base, r)
				case "persistence-barrier":
					return publishParentRecoveryWithSync(dir, base, r, func(*os.Root, string) error { return fault })
				default:
					if err := publishParentRecovery(dir, base, r); err != nil {
						return err
					}
					return fault
				}
			}
			if _, err := parentRecovery(context.Background(), root, "fixture", p.RunID, p.SHA256(), true, publish); err == nil || calls != 1 {
				t.Fatal("publication interruption was ignored", err)
			}
			var retained []byte
			published := phase == "after-publication" || phase == "persistence-barrier"
			if published {
				retained = snapshotRead(t, path)
			} else if _, err := os.Lstat(path); !os.IsNotExist(err) {
				t.Fatal("failed unpublished observation became a record", err)
			}
			entries, err := os.ReadDir(filepath.Join(root, base))
			if err != nil {
				t.Fatal(err)
			}
			for _, entry := range entries {
				if strings.HasPrefix(entry.Name(), ".parent-recovery-") {
					t.Fatal("failed publication leaked its staging file")
				}
			}
			publishCalls := 0
			finish := func(dir *os.Root, base string, r ParentRecovery) error {
				publishCalls++
				return publishParentRecovery(dir, base, r)
			}
			replay, err := parentRecovery(context.Background(), root, "fixture", p.RunID, p.SHA256(), true, finish)
			if err != nil || !sameJSON(replay, p) {
				t.Fatal("exact recovery after interruption failed", err)
			}
			if published && (publishCalls != 0 || !bytes.Equal(retained, snapshotRead(t, path))) {
				t.Fatal("lost publication output rewrote the original record")
			}
			if _, err = PreviewReconciliation(context.Background(), root, "fixture", p.RunID); err != nil {
				t.Fatal(err)
			}
		})
	}
}
