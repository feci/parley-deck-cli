package trajectory

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/procctl"
)

func TestCapturedVerificationStopBeforeHelperIsDurable(t *testing.T) {
	criteria := capturedCriteria()
	ticket, path := journalFixture(t, criteria, dirtyCapturedChild(t))
	reserveJournal(t, ticket)
	if err := StopCapturedVerification(context.Background(), ticket, "wrong-invocation"); err == nil {
		t.Fatal("unrelated invocation stopped the verification")
	}
	if _, err := os.Stat(filepath.Join(path, "stop.json")); !os.IsNotExist(err) {
		t.Fatal("unrelated invocation wrote a stop")
	}
	if err := StopCapturedVerification(context.Background(), ticket, "verifier-invocation"); err != nil {
		t.Fatal(err)
	}
	before := snapshotRead(t, filepath.Join(path, "stop.json"))
	if err := StopCapturedVerification(context.Background(), ticket, "verifier-invocation"); err != nil || !bytes.Equal(before, snapshotRead(t, filepath.Join(path, "stop.json"))) {
		t.Fatal("stop replay rewrote its durable identity", err)
	}
	cmd, output := journalCommand(t, ticket, criteria)
	if err := cmd.Run(); err == nil {
		t.Fatalf("stopped ticket ran in a fresh helper: %s", output)
	}
	for _, name := range []string{"claim.json", "process-001.json", "step-001.json"} {
		if _, err := os.Stat(filepath.Join(path, name)); !os.IsNotExist(err) {
			t.Fatal("stopped ticket advanced", name)
		}
	}
}

func TestCapturedVerificationProcessPublicationFailurePreventsMaterialStart(t *testing.T) {
	trace := filepath.Join(t.TempDir(), "material-started")
	criteria := capturedCriteria()
	criteria[0].Command = "printf started > '" + strings.ReplaceAll(trace, "'", "'\\''") + "'; " + criteria[0].Command
	ticket, path := journalFixture(t, criteria, dirtyCapturedChild(t))
	reserveJournal(t, ticket)
	if err := os.Mkdir(filepath.Join(path, "process-001.json"), 0700); err != nil {
		t.Fatal(err)
	}
	receipt, err := ExecuteCapturedVerification(context.Background(), ticket, "verifier-invocation", criteria, t.TempDir())
	if err == nil || receipt.FailureStage != "execution" {
		t.Fatal("failed process publication accepted", err)
	}
	if _, err := os.Stat(trace); !os.IsNotExist(err) {
		t.Fatal("material command ran before its process identity was durable")
	}
	if _, _, err := ReadCapturedVerification(context.Background(), ticket, "verifier-invocation"); err == nil {
		t.Fatal("unregistered execution accepted")
	}
}

func TestCapturedVerificationLegacyReceiptRemainsReadable(t *testing.T) {
	criteria := capturedCriteria()
	ticket, path := journalFixture(t, criteria, dirtyCapturedChild(t))
	reserveJournal(t, ticket)
	receipt, err := ExecuteCapturedVerification(context.Background(), ticket, "verifier-invocation", criteria, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	// Reconstruct the previous on-disk schema in this test-owned journal.
	// Production only reads it; it never rewrites an existing claim or receipt.
	rewrite := func(name string, value any) string {
		t.Helper()
		data, err := canonical(value)
		if err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(filepath.Join(path, name), data, 0600); err != nil {
			t.Fatal(err)
		}
		return digest(data)
	}
	var claim verificationClaim
	if err = json.Unmarshal(snapshotRead(t, filepath.Join(path, "claim.json")), &claim); err != nil {
		t.Fatal(err)
	}
	claim.Version = 1
	receipt.Version, receipt.ClaimSHA256 = 1, rewrite("claim.json", claim)
	var prepared verificationPrepared
	if err = json.Unmarshal(snapshotRead(t, filepath.Join(path, "prepared.json")), &prepared); err != nil {
		t.Fatal(err)
	}
	prepared.ClaimSHA256 = receipt.ClaimSHA256
	receipt.PreparedSHA256 = rewrite("prepared.json", prepared)
	previous := ""
	for ordinal := 1; ordinal <= receipt.Steps; ordinal++ {
		name := fmt.Sprintf("step-%03d.json", ordinal)
		var step verificationStep
		if err = json.Unmarshal(snapshotRead(t, filepath.Join(path, name)), &step); err != nil {
			t.Fatal(err)
		}
		step.Version, step.ClaimSHA256, step.PreviousSHA256, step.ProcessSHA256 = 1, receipt.ClaimSHA256, previous, ""
		previous = rewrite(name, step)
		if err = os.Remove(filepath.Join(path, fmt.Sprintf("process-%03d.json", ordinal))); err != nil {
			t.Fatal(err)
		}
	}
	receipt.LastSHA256 = previous
	rewrite("receipt.json", receipt)
	if actual, _, err := ReadCapturedVerification(context.Background(), ticket, "verifier-invocation"); err != nil || actual.Version != 1 {
		t.Fatal("legacy completed journal became unreadable", err)
	}
}

func TestCapturedVerificationStopReapsRegisteredExecAndRefusesIdentityChange(t *testing.T) {
	for _, changeIdentity := range []bool{false, true} {
		t.Run(map[bool]string{false: "owned-exec", true: "identity-refusal"}[changeIdentity], func(t *testing.T) {
			signal := filepath.Join(t.TempDir(), "started")
			criteria := capturedCriteria()
			criteria[0].Command = "trap '' TERM; printf started > '" + strings.ReplaceAll(signal, "'", "'\\''") + "'; exec sleep 40"
			ticket, path := journalFixture(t, criteria, dirtyCapturedChild(t))
			reserveJournal(t, ticket)
			cmd, output := journalCommand(t, ticket, criteria)
			if err := cmd.Start(); err != nil {
				t.Fatal(err)
			}
			var identity procctl.Spawned
			waited := false
			defer func() {
				if !waited {
					_, _, _ = procctl.KillTreeAttributed(identity)
					_ = cmd.Process.Kill()
					_ = cmd.Wait()
				}
			}()
			deadline := time.After(20 * time.Second)
			tick := time.NewTicker(10 * time.Millisecond)
			defer tick.Stop()
		waitStarted:
			for {
				select {
				case <-deadline:
					t.Fatal("registered criterion never started")
				case <-tick.C:
					if _, err := os.Stat(signal); err == nil {
						break waitStarted
					}
				}
			}
			processPath := filepath.Join(path, "process-001.json")
			original := snapshotRead(t, processPath)
			var process verificationProcess
			if err := json.Unmarshal(original, &process); err != nil {
				t.Fatal(err)
			}
			identity = process.Identity
			if ok, reason := procctl.Attributed(identity); !ok {
				t.Fatal("material exec changed stable group identity", reason)
			}
			if changeIdentity {
				process.Identity.ProcStart = "wrong-start"
				changed, err := canonical(process)
				if err != nil {
					t.Fatal(err)
				}
				if err = os.WriteFile(processPath, changed, 0600); err != nil {
					t.Fatal(err)
				}
				if err = StopCapturedVerification(context.Background(), ticket, "verifier-invocation"); err == nil || !procctl.Alive(identity) {
					t.Fatal("changed process identity was signalled or accepted", err)
				}
				// Test-owned restoration permits cleanup; production never repairs
				// an attribution mismatch or turns that failure into acceptance.
				if err = os.WriteFile(processPath, original, 0600); err != nil {
					t.Fatal(err)
				}
			}
			if err := StopCapturedVerification(context.Background(), ticket, "verifier-invocation"); err != nil {
				t.Fatal(err)
			}
			err := cmd.Wait()
			waited = true
			if err == nil || procctl.Alive(identity) {
				t.Fatalf("stopped helper accepted or retained supervisor: %v %s", err, output)
			}
			if _, _, err := ReadCapturedVerification(context.Background(), ticket, "verifier-invocation"); err == nil {
				t.Fatal("stopped journal accepted")
			}
			if _, err := os.Stat(filepath.Join(path, "process-002.json")); !os.IsNotExist(err) {
				t.Fatal("another criterion was started after stop")
			}
		})
	}
}
