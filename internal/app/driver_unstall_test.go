package app

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/telemetry"
)

func TestGoalTimeoutTrackAndConfiguredBounds(t *testing.T) {
	for _, tc := range []struct {
		track string
		ms    int
		want  time.Duration
		bad   bool
	}{
		{"", 0, 15 * time.Minute, false}, {"fast", 0, 5 * time.Minute, false}, {"standard", 0, 15 * time.Minute, false}, {"deliberation", 0, 30 * time.Minute, false},
		{"deliberation", 240000, 4 * time.Minute, false}, {"fast", 1800000, 5 * time.Minute, false}, {"standard", -1, 15 * time.Minute, false}, {"standard", int(^uint(0) >> 1), 15 * time.Minute, false},
		{"standrd", 0, 0, true}, {"empty", 0, 0, true}, {"[deliberation]", 0, 0, true},
	} {
		dir := t.TempDir()
		fm := ""
		if tc.track != "" {
			v := tc.track
			if v == "empty" {
				v = ""
			}
			fm = "track: " + v + "\n"
		}
		writePrompt(t, dir, fm)
		got, err := goalCheckTimeout(dir, tc.ms)
		if (err != nil) != tc.bad || got != tc.want {
			t.Fatalf("%+v: got %s %v", tc, got, err)
		}
	}
}

func goalUnstallFixture(t *testing.T, script string, timeoutMS int) (driverImplOps, string) {
	t.Helper()
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	declareAppTestSource(t, root)
	idea := filepath.Join(root, "parley-deck", "ideas", "demo")
	writePrompt(t, idea, "participants: [author, reviewer]\ntrack: deliberation\n")
	counter := filepath.Join(root, "children")
	a := agents.Discovery{Spec: agents.Spec{ID: "reviewer", Commands: []string{"sh"}, HeadlessArgs: []string{"-c", "cat >/dev/null; printf 'child\\n' >> \"$1\"; " + script, "fixture", counter}, PromptMode: agents.PromptStdin, TimeoutMS: timeoutMS, BuffersStdout: true}, Found: true, Path: "/bin/sh"}
	o := newOpsFor(root, idea, []agents.Discovery{a}, "author", []string{"reviewer"})
	o.drafter = "reviewer"
	o.base.RunID = "first"
	o.out = io.Discard
	return o, counter
}

func TestGoalUnstallRetriesAndReplayRemainBounded(t *testing.T) {
	for _, tc := range []struct {
		name, script string
		want         bool
		attempts     int
	}{
		{"valid-fail", "printf 'GOAL-CHECK: FAIL — criterion unmet\\n'", false, 1},
		{"failed-pass", "printf 'GOAL-CHECK: PASS\\n'; exit 7", false, 1},
		{"malformed", "printf 'no verdict\\n'", false, 2},
		{"child-failure", "exit 9", false, 2},
		{"retry-success", "if [ \"$(wc -l < \"$1\" | tr -d ' ')\" = 1 ]; then exit 9; fi; printf 'GOAL-CHECK: PASS\\n'", true, 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			o, counter := goalUnstallFixture(t, tc.script, 240000)
			// The checker is a protected FINAL drafter: retry is execution policy,
			// never authority to remove that participant from membership.
			if err := os.WriteFile(filepath.Join(o.ideaDir, "FINAL.md"), []byte("---\nauthor: reviewer\n---\n"), 0600); err != nil {
				t.Fatal(err)
			}
			ok, why := o.GoalCheck(context.Background())
			if ok != tc.want {
				t.Fatalf("got %v %s", ok, why)
			}
			o.base.RunID = "resumed"
			o.base.Agents[0].TimeoutMS = 1000
			again, why := o.GoalCheck(context.Background())
			if again != tc.want {
				t.Fatalf("replay got %v %s", again, why)
			}
			raw, err := os.ReadFile(counter)
			if err != nil || strings.Count(string(raw), "child\n") != tc.attempts {
				t.Fatalf("attempt amplification: %q %v", raw, err)
			}
			files, _ := filepath.Glob(filepath.Join(o.root, ".parley-runtime", "invocations", "*", "terminal.json"))
			if len(files) != tc.attempts {
				t.Fatalf("terminal count %d", len(files))
			}
			for _, file := range files {
				raw, err := os.ReadFile(file)
				if err != nil {
					t.Fatal(err)
				}
				var r telemetry.Record
				if err = json.Unmarshal(raw, &r); err != nil {
					t.Fatal(err)
				}
				if r.Metadata.ParticipantTimeoutNS != int64(4*time.Minute) {
					t.Fatalf("wrong frozen ceiling: %+v", r.Metadata)
				}
			}
			prompt, _ := os.ReadFile(filepath.Join(o.ideaDir, "00-prompt.md"))
			if !bytes.Contains(prompt, []byte("participants: [author, reviewer]")) {
				t.Fatal("retry changed protected membership")
			}
		})
	}
}

func TestGoalUnstallBufferedHardDeadlineAndCancellation(t *testing.T) {
	const ceiling = time.Second
	o, counter := goalUnstallFixture(t, "exec sleep 30", int(ceiling/time.Millisecond))
	start := time.Now()
	if ok, why := o.GoalCheck(context.Background()); ok {
		t.Fatalf("timed-out check passed: %s", why)
	}
	if time.Since(start) > 15*time.Second {
		t.Fatal("buffered hard deadline not enforced")
	}
	// A launched shell can hit its hard deadline before consuming stdin and
	// appending the counter. Supervisor records are the start/attempt oracle.
	readAttempts := func() map[string]int {
		t.Helper()
		files, err := filepath.Glob(filepath.Join(o.root, ".parley-runtime", "invocations", "*", "terminal.json"))
		if err != nil || len(files) != 2 {
			t.Fatalf("hard timeout terminal records: %d %v", len(files), err)
		}
		attempts := make(map[string]int)
		ordinals := make(map[int]bool)
		for _, file := range files {
			raw, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			var record telemetry.Record
			if err := json.Unmarshal(raw, &record); err != nil {
				t.Fatal(err)
			}
			if record.StartedAt == nil || record.PID == nil || *record.PID <= 0 || record.Outcome == nil ||
				record.Outcome.FailureClass == nil || *record.Outcome.FailureClass != "timeout" ||
				record.Metadata.ParticipantTimeoutNS != int64(ceiling) {
				t.Fatalf("missing started timeout evidence: %+v", record)
			}
			attempts[record.InvocationID] = record.Metadata.AttemptOrdinal
			ordinals[record.Metadata.AttemptOrdinal] = true
		}
		if len(attempts) != 2 || !ordinals[1] || !ordinals[2] {
			t.Fatalf("hard timeout attempt identities: %v", attempts)
		}
		return attempts
	}
	attempts := readAttempts()
	raw, _ := os.ReadFile(counter)
	if strings.Count(string(raw), "child\n") > len(attempts) {
		t.Fatalf("more shell executions than observed starts: %q", raw)
	}
	o.base.RunID = "resumed"
	if ok, why := o.GoalCheck(context.Background()); ok {
		t.Fatalf("replayed timed-out check passed: %s", why)
	}
	for id, ordinal := range readAttempts() {
		if attempts[id] != ordinal {
			t.Fatalf("timeout replay changed invocation identity: %s ordinal %d", id, ordinal)
		}
	}
	replayedCounter, _ := os.ReadFile(counter)
	if !bytes.Equal(raw, replayedCounter) {
		t.Fatal("timeout replay launched an unrecorded shell")
	}
	c, counter2 := goalUnstallFixture(t, "exec sleep 30", 5000)
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()
	if ok, _ := c.GoalCheck(ctx); ok {
		t.Fatal("cancelled checker passed")
	}
	raw, _ = os.ReadFile(counter2)
	if strings.Count(string(raw), "child\n") > 1 {
		t.Fatal("cancellation retried")
	}
}

// Opt-in real process witness crosses the former 120-second product ceiling.
// The normal suite covers the exact policy/telemetry and fast failure cases.
func TestGoalUnstallBeyondFormerCeiling(t *testing.T) {
	if os.Getenv("PARLEY_TEST_LONG_GOAL_CHECK") != "1" {
		t.Skip("set PARLEY_TEST_LONG_GOAL_CHECK=1 for the 121-second process witness")
	}
	o, _ := goalUnstallFixture(t, "sleep 121; printf 'GOAL-CHECK: PASS\\n'", 180000)
	if ok, why := o.GoalCheck(context.Background()); !ok {
		t.Fatalf("check beyond 120s failed: %s", why)
	}
}

func TestBudgetLegacyAttendanceAndReadOnlyPreview(t *testing.T) {
	root := t.TempDir()
	path := "parley-deck/runs/legacy"
	dir := filepath.Join(root, path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "events.jsonl"), []byte(`{"time":"2026-05-10T19:40:03Z","type":"run.created","data":{"task":"smoke"}}`+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	var out, diagnostic bytes.Buffer
	code := runBudgetControl(context.Background(), []string{"legacy", "inspect", "--dir", root, "--run", path}, &out, &diagnostic, false)
	if code != 0 {
		t.Fatalf("inspect: %d %s", code, &diagnostic)
	}
	var p budget.LegacyPreview
	if err := json.Unmarshal(out.Bytes(), &p); err != nil {
		t.Fatal(err)
	}
	args := []string{"legacy", "apply", "--dir", root, "--run", path, "--expected-preview-sha256", p.PreviewSHA256, "--decision-id", "fixture", "--reason", "unknown fixture", "--writers-stopped", "--acknowledge-unknown-history", "--yes"}
	out.Reset()
	diagnostic.Reset()
	if code := runBudgetControl(context.Background(), args, &out, &diagnostic, false); code != 2 || !strings.Contains(diagnostic.String(), "attended") {
		t.Fatalf("unattended apply: %d %s", code, &diagnostic)
	}
	if code := runBudgetPlatformControl(context.Background(), args, io.Discard, io.Discard, false, true); code != 2 {
		t.Fatal("unsupported platform admitted apply")
	}
	if _, err := os.Stat(filepath.Join(root, ".parley-runtime")); !os.IsNotExist(err) {
		t.Fatal("read-only or refused command wrote authority")
	}
	if code := runBudgetControl(context.Background(), args, io.Discard, io.Discard, true); code != 0 {
		t.Fatalf("attended fixture apply: %d", code)
	}
	if code := runBudgetControl(context.Background(), args, io.Discard, io.Discard, true); code != 0 {
		t.Fatalf("attended fixture replay: %d", code)
	}
}
