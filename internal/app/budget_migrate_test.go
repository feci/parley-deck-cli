package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"parley-deck-cli/internal/budget"
)

func migrationCLIArgs(t *testing.T, root string) []string {
	t.Helper()
	var out, errout bytes.Buffer
	if code := runBudgetPlatformControl(context.Background(), []string{"migrate", "inspect", "--dir", root, "--idea", "idea"}, &out, &errout, false, false); code != 0 {
		t.Fatalf("read-only inspection failed: %d %s", code, errout.String())
	}
	var i budget.LaunchMigrationInventory
	if err := json.Unmarshal(out.Bytes(), &i); err != nil {
		t.Fatal(err)
	}
	return []string{"migrate", "apply", "--dir", root, "--idea", "idea", "--kind", "launch",
		"--expected-history-sha256", i.HistorySHA256, "--decision-id", "fixture", "--reason", "Explicit fixture accounting",
		"--started-at", time.Now().UTC().Add(-time.Hour).Format(time.RFC3339Nano), "--additional-launches", "0",
		"--max-launches", "1", "--max-cost-micros", "0", "--wall-clock", "2h", "--writers-stopped", "--yes"}
}

func TestBudgetMigrationCLIRequiresAttendanceAndAllExplicitFields(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	args := migrationCLIArgs(t, root)
	for _, state := range []struct{ supported, attended bool }{{true, false}, {false, true}, {false, false}} {
		if code := runBudgetPlatformControl(ctx, args, io.Discard, io.Discard, state.supported, state.attended); code != 2 {
			t.Fatalf("unsupported/unattended control returned %d", code)
		}
	}
	for _, field := range []string{"--expected-history-sha256", "--decision-id", "--reason", "--started-at", "--additional-launches", "--max-launches", "--max-cost-micros", "--wall-clock", "--writers-stopped", "--yes"} {
		t.Run(field, func(t *testing.T) {
			var incomplete []string
			for n := 0; n < len(args); n++ {
				if args[n] == field {
					if field != "--writers-stopped" && field != "--yes" {
						n++
					}
					continue
				}
				incomplete = append(incomplete, args[n])
			}
			if code := runBudgetPlatformControl(ctx, incomplete, io.Discard, io.Discard, true, true); code == 0 {
				t.Fatal("missing explicit field accepted")
			}
		})
	}
	for _, tail := range [][]string{{"--reserve-micros", "-1"}, {"--kind", "step"}, {"--kind", "cycle"}, {"--wall-clock", "1ns"}, {"--started-at", "invalid"}, {"--additional-launches", "-1"}, {"unexpected"}, {"--attended"}} {
		invalid := append(append([]string{}, args...), tail...)
		if code := runBudgetPlatformControl(ctx, invalid, io.Discard, io.Discard, true, true); code == 0 {
			t.Fatalf("invalid call accepted: %v", tail)
		}
	}
	for _, tail := range [][]string{{"--yes=false"}, {"--max-launches", "0"}, {"--writers-stopped=false"}, {"--kind", "unknown"}, {"unexpected"}} {
		inspect := append([]string{"migrate", "inspect", "--dir", root, "--idea", "idea"}, tail...)
		if code := runBudgetPlatformControl(ctx, inspect, io.Discard, io.Discard, false, false); code == 0 {
			t.Fatal("invalid inspect accepted", tail)
		}
	}
	if _, err := os.Stat(filepath.Join(root, ".parley-runtime")); !os.IsNotExist(err) {
		t.Fatal("invalid/readonly controls created state", err)
	}
}

type migrationOutputFailure struct{}

func (migrationOutputFailure) Write([]byte) (int, error) {
	return 0, errors.New("fixture output unavailable")
}

func TestBudgetMigrationCLIOutputFailureAndReplayPreserveLaterCharge(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	args := migrationCLIArgs(t, root)
	var errout bytes.Buffer
	if code := runBudgetPlatformControl(ctx, args, migrationOutputFailure{}, &errout, true, true); code != 1 {
		t.Fatal("output failure not reported", code)
	}
	b, err := budget.LoadLaunchBinding(ctx, root, "idea")
	if err != nil || b == nil {
		t.Fatal("output failure lost published import", err)
	}
	initial, err := b.Inspect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := b.Reserve(ctx, "later-action"); err != nil {
		t.Fatal(err)
	}
	for n := 0; n < 2; n++ {
		var out bytes.Buffer
		if code := runBudgetPlatformControl(ctx, args, &out, &errout, true, true); code != 0 {
			t.Fatalf("replay: %d %s", code, errout.String())
		}
		var status budget.PolicyStatus
		if err := json.Unmarshal(out.Bytes(), &status); err != nil {
			t.Fatal(err)
		}
		if status.Spent != 1 || !status.StartedAt.Equal(initial.StartedAt) {
			t.Fatalf("replay reset accounting: %+v", status)
		}
	}
	if _, err := b.Reserve(ctx, "over-cap"); !errors.Is(err, budget.ErrLimit) {
		t.Fatal("replay reopened capacity", err)
	}
}
