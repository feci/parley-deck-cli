package evidence

import (
	"context"
	"runtime"
	"testing"
)

func TestDetailedExecutionSeparatesObservedFailureFromIncompleteCapture(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("criterion execution currently requires POSIX")
	}
	for _, tc := range []struct {
		name, command string
		complete      bool
	}{
		{"pass", `printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":0}\n'`, true},
		{"test-failure", `printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":1}\n'; exit 1`, true},
		{"masked-test-failure", `printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":1}\n'`, true},
		{"malformed", `printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":null}\n'`, false},
		{"overflow", `printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":1}\n'; head -c 5000000 /dev/zero; exit 1`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := RunCriterionDetailed(context.Background(), t.TempDir(), "fixture", tc.command, "fixture-verifier")
			if r.Complete != tc.complete {
				t.Fatalf("wrong typed completeness: %+v", r)
			}
			if tc.name != "pass" && r.Record.Status != StatusFail {
				t.Fatalf("failure status changed: %+v", r)
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r := RunCriterionDetailed(ctx, t.TempDir(), "cancelled", "true", "fixture-verifier")
	if r.Complete || r.Record.Status != StatusFail {
		t.Fatalf("cancelled execution admitted: %+v", r)
	}
}
