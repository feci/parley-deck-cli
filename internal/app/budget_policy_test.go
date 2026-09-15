package app

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"parley-deck-cli/internal/budget"
)

func TestBudgetRuntimePolicyControlsRequireAttendanceAndReplayExactly(t *testing.T) {
	ctx := context.Background()
	for _, surface := range []string{"launch", "step"} {
		t.Run(surface, func(t *testing.T) {
			root := t.TempDir()
			kind := budget.Launch
			if surface == "launch" {
				if _, err := budget.ConfigureLaunchBudget(ctx, root, "idea", budget.LaunchPolicy{MaxLaunches: 1}); err != nil {
					t.Fatal(err)
				}
			} else {
				kind = budget.DriverStep
				if _, err := budget.EnsureStepBinding(ctx, root, "idea", 1, 0); err != nil {
					t.Fatal(err)
				}
			}
			var out, errout bytes.Buffer
			inspect := []string{surface, "inspect", "--dir", root, "--idea", "idea"}
			if code := runBudgetPlatformControl(ctx, inspect, &out, &errout, false, false); code != 0 {
				t.Fatalf("read-only inspect: %d %s", code, errout.String())
			}
			var preview budget.PolicyStatus
			if err := json.Unmarshal(out.Bytes(), &preview); err != nil {
				t.Fatal(err)
			}
			args := []string{surface, "extend", "--dir", root, "--idea", "idea", "--wall-clock", "0", "--expected-policy-sha256", preview.PolicySHA256, "--decision-id", "fixture", "--reason", "Explicit fixture decision", "--yes"}
			if kind == budget.Launch {
				args = append(args, "--max-launches", "2", "--max-cost-micros", "0")
			} else {
				args = append(args, "--max-steps", "2")
			}
			for _, state := range []struct{ supported, attended bool }{{true, false}, {false, true}} {
				out.Reset()
				errout.Reset()
				if code := runBudgetPlatformControl(ctx, args, &out, &errout, state.supported, state.attended); code != 2 {
					t.Fatalf("unattended/unsupported grant: %d %s", code, errout.String())
				}
			}
			for i := 0; i < 2; i++ {
				out.Reset()
				errout.Reset()
				if code := runBudgetPlatformControl(ctx, args, &out, &errout, true, true); code != 0 {
					t.Fatalf("grant/replay: %d %s", code, errout.String())
				}
			}
			current, err := budget.InspectRuntimeBudget(ctx, root, "idea", kind)
			if err != nil || current.Spent != 0 || current.PolicySHA256 == preview.PolicySHA256 {
				t.Fatalf("grant not recorded: %+v %v", current, err)
			}
			var p struct {
				Extensions []budget.PolicyExtension `json:"extensions"`
			}
			if err := json.Unmarshal(current.Policy, &p); err != nil {
				t.Fatal(err)
			}
			if len(p.Extensions) != 1 {
				t.Fatal("exact replay duplicated grant")
			}
		})
	}
}

func TestBudgetRuntimePolicyInvalidCallsCreateNoState(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	for _, surface := range []string{"launch", "step"} {
		for _, tail := range [][]string{{"inspect"}, {"inspect", "--wall-clock", "0"}, {"extend", "--yes"}, {"extend", "--max-steps", "2", "--yes"}} {
			args := append([]string{surface}, tail...)
			args = append(args, "--dir", root, "--idea", "idea")
			var out, errout bytes.Buffer
			if code := runBudgetPlatformControl(ctx, args, &out, &errout, true, true); code == 0 {
				t.Fatalf("invalid command passed: %v", args)
			}
		}
	}
	if _, err := os.Stat(filepath.Join(root, ".parley-runtime")); !os.IsNotExist(err) {
		t.Fatal("inspection/refusal initialized state")
	}
}

func TestBudgetRuntimePolicyWallClockIsAbsolute(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	if _, err := budget.EnsureStepBinding(ctx, root, "idea", 1, time.Hour); err != nil {
		t.Fatal(err)
	}
	s, err := budget.InspectRuntimeBudget(ctx, root, "idea", budget.DriverStep)
	if err != nil {
		t.Fatal(err)
	}
	var out, errout bytes.Buffer
	args := []string{"step", "extend", "--dir", root, "--idea", "idea", "--max-steps", strconv.Itoa(1), "--wall-clock", "2h", "--expected-policy-sha256", s.PolicySHA256, "--decision-id", "time", "--reason", "Explicit absolute clock grant", "--yes"}
	if code := runBudgetControl(ctx, args, &out, &errout, true); code != 0 {
		t.Fatalf("clock grant: %d %s", code, errout.String())
	}
	b, err := budget.LoadStepBinding(ctx, root, "idea")
	if err != nil || b.Policy.WallClockNS != int64(2*time.Hour) {
		t.Fatalf("clock grant: %+v %v", b, err)
	}
}
