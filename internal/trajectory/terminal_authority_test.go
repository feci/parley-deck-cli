package trajectory

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/telemetry"
)

func TestTerminalPublicationRetainsActualOutcomeWithUnavailableHistory(t *testing.T) {
	for _, mode := range []string{"missing-archive", "missing-earlier-result", "changed-original-ledger", "changed-original-state"} {
		t.Run(mode, func(t *testing.T) {
			ctx := context.Background()
			root, binding := unchangedFixture(t, "true")
			var earlier telemetry.Record
			if mode == "missing-earlier-result" {
				earlier = unchangedProcess(t, root, binding, "exit 0", false)
				p, err := PreviewUnchanged(ctx, root, "fixture", 1)
				if err != nil {
					t.Fatal(err)
				}
				if _, err = ReconcileUnchanged(ctx, root, "fixture", 1, p.SHA256()); err != nil {
					t.Fatal(err)
				}
				c, err := PreviewContinuation(ctx, root, "fixture")
				if err != nil {
					t.Fatal(err)
				}
				if _, err = Continue(ctx, root, "fixture", c.SHA256(), "terminal-authority-decision", "Synthetic fixture decision", false, true); err != nil {
					t.Fatal(err)
				}
			}
			charged := chargeFixture(t, root, binding)
			invocation, err := telemetry.Begin(filepath.Join(root, ".parley-runtime", "invocations"), telemetry.Metadata{RunID: "terminal-authority-fixture", Idea: "fixture", Phase: "fixup", Agent: "builder", LaunchMode: "headless"})
			if err != nil {
				t.Fatal(err)
			}
			run, err := Begin(charged, root, "fixture", "builder", invocation.ID)
			if err != nil {
				t.Fatal(err)
			}
			cmd := exec.Command("sh", "-c", "printf 'changed\\n' > source; exit 7")
			cmd.Dir = root
			if err = cmd.Start(); err != nil {
				t.Fatal(err)
			}
			if err = invocation.Started(cmd.Process.Pid); err != nil {
				cmd.Process.Kill()
				cmd.Wait()
				t.Fatal(err)
			}
			var failure *exec.ExitError
			if err = cmd.Wait(); !errors.As(err, &failure) {
				t.Fatal("expected actual exit 7", err)
			}
			code := failure.ExitCode()
			if code != 7 {
				t.Fatal(code)
			}
			if err = invocation.Finish(telemetry.Outcome{Status: "failed", ExitCode: &code, FailureClass: telemetry.String("process_failure")}); err != nil {
				t.Fatal(err)
			}
			state, _, err := readState(statePath(*binding))
			if err != nil {
				t.Fatal(err)
			}
			ledger, err := binding.Store.Inspect(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if mode == "missing-archive" {
				if err = os.Remove(snapshotPath(snapshotDirectory(*binding), state.BaselineArchive)); err != nil {
					t.Fatal(err)
				}
			} else if mode == "missing-earlier-result" {
				if err = os.Remove(filepath.Join(root, ".parley-runtime", "invocations", earlier.InvocationID, "terminal.json")); err != nil {
					t.Fatal(err)
				}
			} else if mode == "changed-original-state" {
				state.Attempts[len(state.Attempts)-1].ReservationIntentSHA256 = strings.Repeat("a", 64)
				if err = writeState(statePath(*binding), state); err != nil {
					t.Fatal(err)
				}
			} else {
				one := int64(1)
				state.Attempts[len(state.Attempts)-1].Charge.ReserveMicros = &one
				if err = writeState(statePath(*binding), state); err != nil {
					t.Fatal(err)
				}
			}
			before := snapshotRead(t, statePath(*binding))
			err = run.Finish(ctx, "failed", &code)
			if mode == "changed-original-state" {
				if err == nil || !strings.Contains(err.Error(), "trajectory state changed after the original launch") {
					t.Fatal("changed launch state authorized terminal", err)
				}
				if !bytes.Equal(before, snapshotRead(t, statePath(*binding))) {
					t.Fatal("changed launch state wrote terminal")
				}
				return
			}
			if mode == "changed-original-ledger" {
				if err == nil || !strings.Contains(err.Error(), "active accounting session lost or changed its original reservation") {
					t.Fatal("changed original ledger authorized terminal", err)
				}
				if !bytes.Equal(before, snapshotRead(t, statePath(*binding))) {
					t.Fatal("changed ledger wrote terminal")
				}
				return
			}
			if err != nil {
				t.Fatalf("actual terminal depended on unavailable historical evidence: %v", err)
			}
			retained, _, err := readState(statePath(*binding))
			if err != nil {
				t.Fatal(err)
			}
			a := retained.Attempts[len(retained.Attempts)-1]
			if a.Terminal == nil || a.Terminal.Status != "failed" || a.Terminal.ExitCode == nil || *a.Terminal.ExitCode != 7 || a.Terminal.SnapshotError != "" || a.After == nil || a.AfterArchive == nil || a.After.Tree.SHA256 == a.Before.Tree.SHA256 {
				t.Fatal("lost actual failed patch", a)
			}
			afterLedger, err := binding.Store.Inspect(ctx)
			if err != nil || !sameJSON(ledger, afterLedger) {
				t.Fatal("terminal changed ledger", err)
			}
			if _, err = Inspect(ctx, root, "fixture"); err == nil {
				t.Fatal("unavailable history accepted by Inspect")
			}
			if err = RequireResolved(ctx, root, "fixture"); err == nil {
				t.Fatal("unavailable history permitted completion")
			}
			before = snapshotRead(t, statePath(*binding))
			if err = run.Finish(ctx, "failed", &code); err == nil {
				t.Fatal("terminal replay accepted")
			}
			if !bytes.Equal(before, snapshotRead(t, statePath(*binding))) {
				t.Fatal("terminal replay changed original bytes")
			}
		})
	}
}
