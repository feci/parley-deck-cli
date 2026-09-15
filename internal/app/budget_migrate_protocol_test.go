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

func protocolMigrationCLIArgs(t *testing.T, root, kind string) []string {
	t.Helper()
	var out, errout bytes.Buffer
	if code := runBudgetPlatformControl(context.Background(), []string{"migrate", "inspect", "--dir", root, "--idea", "idea", "--kind", kind}, &out, &errout, false, false); code != 0 {
		t.Fatalf("read-only protocol inventory: %d %s", code, errout.String())
	}
	var i budget.ProtocolMigrationInventory
	if err := json.Unmarshal(out.Bytes(), &i); err != nil {
		t.Fatal(err)
	}
	args := []string{"migrate", "apply", "--dir", root, "--idea", "idea", "--kind", kind,
		"--expected-history-sha256", i.HistorySHA256, "--decision-id", "fixture", "--reason", "Explicit fixture protocol reconciliation",
		"--started-at", time.Now().UTC().Add(-time.Hour).Format(time.RFC3339Nano), "--total-actions", "1", "--writers-stopped", "--yes"}
	if kind == "step" {
		return append(args, "--max-steps", "2", "--wall-clock", "2h")
	}
	return append(args, "--max-cycles", "2")
}

func TestBudgetProtocolMigrationCLIRequiresExplicitFieldsAndAttendance(t *testing.T) {
	for _, kind := range []string{"step", "fixup", "cross-review"} {
		t.Run(kind, func(t *testing.T) {
			ctx, root := context.Background(), t.TempDir()
			args := protocolMigrationCLIArgs(t, root, kind)
			for _, state := range []struct{ supported, attended bool }{{true, false}, {false, true}, {false, false}} {
				if code := runBudgetPlatformControl(ctx, args, io.Discard, io.Discard, state.supported, state.attended); code != 2 {
					t.Fatalf("unsupported/unattended: %d", code)
				}
			}
			fields := []string{"--idea", "--expected-history-sha256", "--decision-id", "--reason", "--started-at", "--total-actions", "--writers-stopped", "--yes"}
			if kind == "step" {
				fields = append(fields, "--max-steps", "--wall-clock")
			} else {
				fields = append(fields, "--max-cycles")
			}
			for _, field := range fields {
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
					t.Fatal("missing field accepted", field)
				}
			}
			invalid := [][]string{{"--max-launches", "0"}, {"--additional-launches", "0"}, {"--max-cost-micros", "0"}, {"--reserve-micros", "0"}, {"--total-actions", "-1"}, {"--started-at", "invalid"}, {"--writers-stopped=false"}, {"--yes=false"}, {"--idea-path", "parley-deck/../escape"}, {"unexpected"}}
			if kind == "step" {
				invalid = append(invalid, []string{"--max-cycles", "0"}, []string{"--max-steps", "-1"}, []string{"--wall-clock", "-1s"})
			} else {
				invalid = append(invalid, []string{"--max-steps", "0"}, []string{"--wall-clock", "0"}, []string{"--max-cycles", "-1"})
			}
			for _, tail := range invalid {
				if code := runBudgetPlatformControl(ctx, append(append([]string{}, args...), tail...), io.Discard, io.Discard, true, true); code == 0 {
					t.Fatal("mixed/invalid fields accepted", tail)
				}
			}
			for _, tail := range [][]string{{"--yes=false"}, {"--total-actions", "0"}, {"--max-cycles", "0"}, {"--max-steps", "0"}, {"--wall-clock", "0"}, {"--writers-stopped=false"}} {
				args := append([]string{"migrate", "inspect", "--dir", root, "--idea", "idea", "--kind", kind}, tail...)
				if code := runBudgetPlatformControl(ctx, args, io.Discard, io.Discard, false, false); code == 0 {
					t.Fatal("inspect accepted mutation flags", tail)
				}
			}
			if _, err := os.Stat(filepath.Join(root, ".parley-runtime")); !os.IsNotExist(err) {
				t.Fatal("invalid/readonly command created accounting", err)
			}
		})
	}
}

func TestBudgetProtocolMigrationCLIOutputFailureAndReplayPreserveSpend(t *testing.T) {
	for _, tc := range []struct {
		name string
		kind budget.Kind
	}{{"step", budget.DriverStep}, {"fixup", budget.Fixup}, {"cross-review", budget.CrossReview}} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, root := context.Background(), t.TempDir()
			args := protocolMigrationCLIArgs(t, root, tc.name)
			var errout bytes.Buffer
			if code := runBudgetPlatformControl(ctx, args, migrationOutputFailure{}, &errout, true, true); code != 1 {
				t.Fatal("output failure not reported", code)
			}
			var charge func() error
			if tc.kind == budget.DriverStep {
				b, err := budget.LoadStepBinding(ctx, root, "idea")
				if err != nil || b == nil {
					t.Fatal("published step import lost", err)
				}
				charge = func() error {
					scoped, finish, err := budget.OpenStepSession(ctx, b)
					if err != nil {
						return err
					}
					defer finish()
					return budget.ChargeStep(scoped)
				}
			} else {
				b, err := budget.LoadCycleBinding(ctx, root, "idea", tc.kind)
				if err != nil || b == nil {
					t.Fatal("published cycle import lost", err)
				}
				charge = func() error {
					scoped, finish, err := budget.OpenCycleSession(ctx, b)
					if err != nil {
						return err
					}
					defer finish()
					_, err = budget.ChargeCycle(scoped, tc.kind)
					return err
				}
			}
			if err := charge(); err != nil {
				t.Fatal(err)
			}
			var origin time.Time
			for n := 0; n < 2; n++ {
				var out bytes.Buffer
				if code := runBudgetPlatformControl(ctx, args, &out, &errout, true, true); code != 0 {
					t.Fatalf("replay: %d %s", code, errout.String())
				}
				var s budget.PolicyStatus
				if err := json.Unmarshal(out.Bytes(), &s); err != nil {
					t.Fatal(err)
				}
				if n == 0 {
					origin = s.StartedAt
				}
				if s.Spent != 2 || !s.StartedAt.Equal(origin) {
					t.Fatalf("replay reset count/epoch: %+v", s)
				}
			}
			if err := charge(); !errors.Is(err, budget.ErrLimit) {
				t.Fatal("replay reopened capacity", err)
			}
		})
	}
}
