package trajectory

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/procctl"
)

func TestStateValidationDoesNotBlockRegisteredCriterionStop(t *testing.T) {
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
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	tick := time.NewTicker(10 * time.Millisecond)
	defer tick.Stop()
waitStarted:
	for {
		select {
		case <-ctx.Done():
			t.Fatal("registered criterion never started", ctx.Err())
		case <-tick.C:
			if _, err := os.Stat(signal); err == nil {
				break waitStarted
			}
		}
	}
	var process verificationProcess
	if err := json.Unmarshal(snapshotRead(t, filepath.Join(path, "process-001.json")), &process); err != nil {
		t.Fatal(err)
	}
	identity = process.Identity
	if ok, reason := procctl.Attributed(identity); !ok {
		t.Fatal("registered process lacks attribution", reason)
	}
	entered, release := make(chan struct{}), make(chan struct{})
	done := make(chan error, 1)
	go func() {
		done <- withStateResolutionCheck(ctx, ticket.Root, ticket.Request.Idea, func(ctx context.Context, b budget.CycleBinding, s State) error {
			if err := checkResolutions(ctx, b, s); err != nil {
				return err
			}
			close(entered)
			select {
			case <-release:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		}, func(budget.CycleBinding, budget.Snapshot, State) error { return nil })
	}()
	select {
	case <-entered:
	case err := <-done:
		t.Fatal("reader never paused", err)
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
	stopCtx, stop := context.WithTimeout(ctx, 2*time.Second)
	stopErr := StopCapturedVerification(stopCtx, ticket, "verifier-invocation")
	stop()
	close(release)
	readerErr := <-done
	if stopErr != nil {
		t.Fatal("full-history reader blocked registered criterion stop", stopErr)
	}
	err := cmd.Wait()
	waited = true
	if err == nil || procctl.Alive(identity) {
		t.Fatalf("stopped helper accepted or criterion survived: %v %s", err, output)
	}
	if readerErr != nil {
		t.Fatal("unchanged authority reader failed after stop", readerErr)
	}
	if _, err = os.Stat(filepath.Join(path, "stop.json")); err != nil {
		t.Fatal("durable stop absent", err)
	}
	if _, err = os.Stat(filepath.Join(path, "process-002.json")); !os.IsNotExist(err) {
		t.Fatal("another criterion started after stop", err)
	}
	if _, _, err = ReadCapturedVerification(ctx, ticket, "verifier-invocation"); err == nil {
		t.Fatal("stopped helper became acceptance")
	}
}
