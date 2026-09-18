package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"parley-deck-cli/internal/budget"
)

func TestBudgetCycleOperatorControlAttendancePreviewAndReplay(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	b, err := budget.EnsureCycleBinding(ctx, root, "idea", budget.Fixup, 1, 1, "", "")
	if err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	inspect := []string{"cycle", "inspect", "--dir", root, "--idea", "idea", "--kind", "fixup"}
	if code := runBudgetPlatformControl(ctx, inspect, &stdout, &stderr, false, false); code != 0 {
		t.Fatalf("read-only inspect: %d %s", code, stderr.String())
	}
	var preview budget.CycleStatus
	if err := json.Unmarshal(stdout.Bytes(), &preview); err != nil {
		t.Fatal(err)
	}
	args := []string{"cycle", "extend", "--dir", root, "--idea", "idea", "--kind", "fixup", "--max-cycles", "2", "--expected-policy-sha256", preview.PolicySHA256, "--decision-id", "one", "--reason", "Explicit fixture grant", "--yes"}
	policy := filepath.Join(filepath.Dir(b.Store.Dir), "policy.json")
	before, err := os.ReadFile(policy)
	if err != nil {
		t.Fatal(err)
	}
	for _, platform := range []struct{ supported, attended bool }{{true, false}, {false, true}} {
		stdout.Reset()
		stderr.Reset()
		if code := runBudgetPlatformControl(ctx, args, &stdout, &stderr, platform.supported, platform.attended); code != 2 {
			t.Fatalf("unattended/unsupported grant: %d %s", code, stderr.String())
		}
		after, _ := os.ReadFile(policy)
		if !bytes.Equal(before, after) {
			t.Fatal("refused control changed policy")
		}
	}
	for i := 0; i < 2; i++ {
		stdout.Reset()
		stderr.Reset()
		if code := runBudgetPlatformControl(ctx, args, &stdout, &stderr, true, true); code != 0 {
			t.Fatalf("attended grant/replay: %d %s", code, stderr.String())
		}
	}
	current, err := budget.InspectCycleBudget(ctx, root, "idea", budget.Fixup)
	if err != nil || current.Spent != 1 || current.Policy.Maximum != 2 || len(current.Policy.Extensions) != 1 {
		t.Fatalf("operator history: %+v %v", current, err)
	}
	if n, err := b.Reserve(ctx, "second"); err != nil || n != 2 {
		t.Fatalf("control did not grant actual work: %d %v", n, err)
	}
	if _, err := b.Reserve(ctx, "third"); !errors.Is(err, budget.ErrLimit) {
		t.Fatalf("control granted unlimited work: %v", err)
	}
}

func TestBudgetCycleInspectAndInvalidControlsDoNotInitializeState(t *testing.T) {
	root := t.TempDir()
	ctx := context.Background()
	for _, args := range [][]string{
		{"cycle", "inspect", "--dir", root, "--idea", "idea", "--kind", "fixup"},
		{"cycle", "inspect", "--dir", root, "--idea", "idea", "--kind", "fixup", "--max-cycles", "7"},
		{"cycle", "extend", "--dir", root, "--idea", "idea", "--kind", "fixup", "--max-cycles", "7", "--yes"},
		{"cycle", "extend", "--dir", root, "--idea", "idea", "--kind", "launch", "--yes"},
	} {
		var out, errout bytes.Buffer
		if code := runBudgetPlatformControl(ctx, args, &out, &errout, true, true); code == 0 {
			t.Fatalf("invalid control accepted: %v", args)
		}
	}
	if _, err := os.Stat(filepath.Join(root, ".parley-runtime")); !os.IsNotExist(err) {
		t.Fatal("inspection/refusal created runtime state")
	}
}
