package trajectory

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/procctl"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCapturedCleanupSurvivesUnavailableArchive(t *testing.T) {
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

			binding, bindErr := budget.LoadCycleBinding(context.Background(), ticket.Root, ticket.Request.Idea, budget.Fixup)
			if bindErr != nil {
				t.Fatal(bindErr)
			}
			retained, _, readErr := readState(statePath(*binding))
			if readErr != nil {
				t.Fatal(readErr)
			}
			if removeErr := os.Remove(snapshotPath(snapshotDirectory(*binding), retained.BaselineArchive)); removeErr != nil {
				t.Fatal(removeErr)
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
				if err = StopCapturedVerification(context.Background(), ticket, "verifier-invocation"); err == nil || !strings.Contains(err.Error(), "process start time mismatch") || !procctl.Alive(identity) {
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

func TestCapturedCleanupRetainsStopWhenEarlierResolutionEvidenceDisappears(t *testing.T) {
	ctx := context.Background()
	root, binding := unchangedFixture(t, "true")
	original := unchangedProcess(t, root, binding, "exit 0", false)
	preview, err := PreviewUnchanged(ctx, root, "fixture", 1)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = ReconcileUnchanged(ctx, root, "fixture", 1, preview.SHA256()); err != nil {
		t.Fatal(err)
	}
	promotion, err := PreviewContinuation(ctx, root, "fixture")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = Continue(ctx, root, "fixture", promotion.SHA256(), "cleanup-fixture-decision", "Synthetic fixture acknowledgment", false, true); err != nil {
		t.Fatal(err)
	}
	unchangedProcess(t, root, binding, "printf changed > source", false)
	ticket, err := PrepareCapturedVerification(ctx, root, "fixture", "reviewer", "cleanup-later-ticket")
	if err != nil {
		t.Fatal(err)
	}
	reserveJournal(t, ticket)
	if err = os.Remove(filepath.Join(root, ".parley-runtime", "invocations", original.InvocationID, "terminal.json")); err != nil {
		t.Fatal(err)
	}
	if _, err = Inspect(ctx, root, "fixture"); err == nil {
		t.Fatal("missing original result permitted ordinary history inspection")
	}
	if err = StopCapturedVerification(ctx, ticket, "verifier-invocation"); err != nil {
		t.Fatalf("cleanup depended on missing earlier result: %v", err)
	}
	dir, err := openVerificationDirectory(*binding, ticket.Request.Charge.EntryKey, false)
	if err != nil {
		t.Fatal(err)
	}
	defer dir.Close()
	var stop verificationStop
	before, err := readVerificationArtifact(dir, "stop.json", &stop)
	if err != nil {
		t.Fatal(err)
	}
	if err = StopCapturedVerification(ctx, ticket, "verifier-invocation"); err != nil {
		t.Fatal(err)
	}
	after, err := readVerificationArtifact(dir, "stop.json", &stop)
	if err != nil || after != before {
		t.Fatal("cleanup replay changed original stop", err)
	}
	if _, err = dir.Lstat("claim.json"); !os.IsNotExist(err) {
		t.Fatal("cleanup invented a helper execution")
	}
}

func TestCapturedCleanupBeforeHelperRetainsStopWithoutArchive(t *testing.T) {
	ticket, path := journalFixture(t, capturedCriteria(), dirtyCapturedChild(t))
	reserveJournal(t, ticket)
	binding, err := budget.LoadCycleBinding(context.Background(), ticket.Root, ticket.Request.Idea, budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	retained, _, err := readState(statePath(*binding))
	if err != nil {
		t.Fatal(err)
	}
	if err = os.Remove(snapshotPath(snapshotDirectory(*binding), retained.BaselineArchive)); err != nil {
		t.Fatal(err)
	}
	if err = StopCapturedVerification(context.Background(), ticket, "verifier-invocation"); err != nil {
		t.Fatalf("cleanup depended on missing archive: %v", err)
	}
	before := snapshotRead(t, filepath.Join(path, "stop.json"))
	if err = StopCapturedVerification(context.Background(), ticket, "verifier-invocation"); err != nil || !bytes.Equal(before, snapshotRead(t, filepath.Join(path, "stop.json"))) {
		t.Fatal("cleanup replay rewrote original stop", err)
	}
	if _, err = ExecuteCapturedVerification(context.Background(), ticket, "verifier-invocation", capturedCriteria(), t.TempDir()); err == nil {
		t.Fatal("cleanup permitted new helper execution")
	}
	if _, err = os.Stat(filepath.Join(path, "claim.json")); !os.IsNotExist(err) {
		t.Fatal("cleanup created a helper claim")
	}
}

func TestCapturedCleanupPreservesOriginalLedgerAuthority(t *testing.T) {
	ticket, path := journalFixture(t, capturedCriteria(), dirtyCapturedChild(t))
	reserveJournal(t, ticket)
	binding, err := budget.LoadCycleBinding(context.Background(), ticket.Root, ticket.Request.Idea, budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	state, _, err := readState(statePath(*binding))
	if err != nil {
		t.Fatal(err)
	}
	one := int64(1)
	state.Attempts[ticket.Request.Sequence-1].Charge.ReserveMicros = &one
	// Coordinate the test-owned state/ticket/launch rewrites while the actual
	// published charge stays unchanged. Only the original-ledger predicate may
	// authorize control; self-consistent rewritten request bytes are insufficient.
	ticket.Request, err = capturedRequestAt(state, ticket.Request.Verifier, ticket.Request.Sequence)
	if err != nil {
		t.Fatal(err)
	}
	if err = writeState(statePath(*binding), state); err != nil {
		t.Fatal(err)
	}
	ticketSHA, err := ticket.SHA256()
	if err != nil {
		t.Fatal(err)
	}
	var launch verificationLaunch
	if err = json.Unmarshal(snapshotRead(t, filepath.Join(path, "launch.json")), &launch); err != nil {
		t.Fatal(err)
	}
	launch.TicketSHA256 = ticketSHA
	for name, value := range map[string]any{"request.json": ticket, "launch.json": launch} {
		raw, e := canonical(value)
		if e != nil {
			t.Fatal(e)
		}
		if e = os.WriteFile(filepath.Join(path, name), raw, 0600); e != nil {
			t.Fatal(e)
		}
	}
	err = StopCapturedVerification(context.Background(), ticket, "verifier-invocation")
	if err == nil || !strings.Contains(err.Error(), "active accounting session lost or changed its original reservation") {
		t.Fatalf("cleanup accepted charge/ticket rewritten against original ledger: %v", err)
	}
	if _, err = os.Stat(filepath.Join(path, "stop.json")); !os.IsNotExist(err) {
		t.Fatal("changed ledger authority published a stop")
	}
}
