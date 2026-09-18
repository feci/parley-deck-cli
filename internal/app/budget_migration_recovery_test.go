package app

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/telemetry"
)

func recoveryCLIFixture(t *testing.T, root string, kind budget.Kind) ([]string, string) {
	t.Helper()
	ctx := context.Background()
	var status budget.PolicyStatus
	var err error
	rawKind := string(kind)
	if kind == budget.DriverStep {
		rawKind = "step"
	}
	if kind == budget.Launch {
		var out, errout bytes.Buffer
		if code := runBudgetPlatformControl(ctx, migrationCLIArgs(t, root), &out, &errout, true, true); code != 0 {
			t.Fatalf("initial fixture: %d %s", code, errout.String())
		}
		if err := json.Unmarshal(out.Bytes(), &status); err != nil {
			t.Fatal(err)
		}
	} else {
		i, err := budget.InspectProtocolMigration(ctx, root, "idea", kind, "")
		if err != nil {
			t.Fatal(err)
		}
		status, err = budget.MigrateProtocolBudget(ctx, root, "idea", kind, budget.ProtocolMigrationRequest{ExpectedHistorySHA256: i.HistorySHA256, DecisionID: "initial-fixture", Reason: "Explicit fixture-only accounting", StartedAt: time.Now().UTC().Add(-time.Hour), TotalActions: 1, Maximum: 3, WritersStopped: true})
		if err != nil {
			t.Fatal(err)
		}
	}
	// Model a lost final publication marker only in this isolated fixture.
	dir := filepath.Dir(status.LedgerDir)
	if err := os.Remove(filepath.Join(dir, "migration-active")); err != nil {
		t.Fatal(err)
	}
	i, err := telemetry.Begin(filepath.Join(root, ".parley-runtime", "invocations"), telemetry.Metadata{Idea: "idea", Phase: "fixup", Agent: "fixture", RunID: "fixture", SegmentID: "fixture"})
	if err != nil {
		t.Fatal(err)
	}
	if err := i.Finish(telemetry.Outcome{Status: "failed", FailureClass: telemetry.String("budget_refused")}); err != nil {
		t.Fatal(err)
	}
	var out, errout bytes.Buffer
	inspect := []string{"migrate", "recover", "inspect", "--dir", root, "--idea", "idea", "--kind", rawKind}
	if code := runBudgetPlatformControl(ctx, inspect, &out, &errout, false, false); code != 0 {
		t.Fatalf("read-only preview: %d %s", code, errout.String())
	}
	var preview budget.MigrationRecoveryPreview
	if err := json.Unmarshal(out.Bytes(), &preview); err != nil {
		t.Fatal(err)
	}
	countFlag, count := "--additional-launches", "0"
	if kind != budget.Launch {
		countFlag, count = "--total-actions", "1"
	}
	return []string{"migrate", "recover", "apply", "--dir", root, "--idea", "idea", "--kind", rawKind,
		"--expected-import-sha256", preview.ImportSHA256, "--expected-history-sha256", preview.HistorySHA256,
		"--decision-id", "recovery-fixture", "--reason", "Explicit fixture-only recovery decision",
		"--started-at", preview.Retained.StartedAt.Format(time.RFC3339Nano), countFlag, count, "--writers-stopped", "--yes"}, dir
}

func recoveryCLIState(t *testing.T, root string) map[string]string {
	t.Helper()
	state := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		h := sha256.Sum256(b)
		state[path] = hex.EncodeToString(h[:])
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return state
}

func TestBudgetRecoveryCLIRejectsUnattendedAndIncompleteDecisions(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	args, _ := recoveryCLIFixture(t, root, budget.Launch)
	before := recoveryCLIState(t, root)
	for _, s := range []struct{ supported, attended bool }{{true, false}, {false, true}, {false, false}} {
		if code := runBudgetPlatformControl(ctx, args, io.Discard, io.Discard, s.supported, s.attended); code != 2 {
			t.Fatalf("unattended/unsupported decision returned %d", code)
		}
	}
	for _, field := range []string{"--expected-import-sha256", "--expected-history-sha256", "--decision-id", "--reason", "--started-at", "--additional-launches", "--writers-stopped", "--yes"} {
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
			t.Fatal("missing explicit decision accepted", field)
		}
	}
	for _, extra := range [][]string{{"--total-actions", "0"}, {"--max-launches", "99"}, {"--max-cost-micros", "0"}, {"--attended"}, {"--writers-stopped=false"}, {"--yes=false"}, {"--additional-launches", "-1"}, {"unexpected"}} {
		invalid := append(append([]string{}, args...), extra...)
		if code := runBudgetPlatformControl(ctx, invalid, io.Discard, io.Discard, true, true); code == 0 {
			t.Fatal("invalid recovery accepted", extra)
		}
	}
	for _, extra := range [][]string{{"--yes=false"}, {"--additional-launches", "0"}, {"--kind", "unknown"}, {"--expected-import-sha256", "bad"}} {
		inspect := append([]string{"migrate", "recover", "inspect", "--dir", root, "--idea", "idea"}, extra...)
		if code := runBudgetPlatformControl(ctx, inspect, io.Discard, io.Discard, false, false); code == 0 {
			t.Fatal("invalid inspect accepted", extra)
		}
	}
	if after := recoveryCLIState(t, root); !reflect.DeepEqual(before, after) {
		t.Fatal("refused controls changed import files")
	}
}

func TestBudgetRecoveryCLIAllKindsOutputFailureAndReplay(t *testing.T) {
	for _, kind := range []budget.Kind{budget.Launch, budget.DriverStep, budget.Fixup, budget.CrossReview} {
		t.Run(string(kind), func(t *testing.T) {
			ctx := context.Background()
			root := t.TempDir()
			args, dir := recoveryCLIFixture(t, root, kind)
			var errout bytes.Buffer
			if code := runBudgetPlatformControl(ctx, args, migrationOutputFailure{}, &errout, true, true); code != 1 {
				t.Fatalf("lost output failure: %d %s", code, errout.String())
			}
			before := recoveryCLIState(t, dir)
			var out bytes.Buffer
			if code := runBudgetPlatformControl(ctx, args, &out, &errout, true, true); code != 0 {
				t.Fatalf("exact replay: %d %s", code, errout.String())
			}
			var status budget.PolicyStatus
			if err := json.Unmarshal(out.Bytes(), &status); err != nil {
				t.Fatal(err)
			}
			want := 1
			if kind == budget.Launch {
				want = 0
			}
			if status.Spent != want {
				t.Fatalf("replay changed accounting: %+v", status)
			}
			if after := recoveryCLIState(t, dir); !reflect.DeepEqual(before, after) {
				t.Fatal("exact active replay rewrote files")
			}
		})
	}
}
