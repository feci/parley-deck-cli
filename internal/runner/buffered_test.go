package runner

import (
	"context"
	"io"
	"os"
	"regexp"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
)

// TestBufferedQuietLateSuccessHelper is re-executed as a real child process: a
// declared-buffered agent that is silent for a while and then produces a valid
// artifact (quiet late success). It must not be killed by a soft watchdog.
func TestBufferedQuietLateSuccessHelper(t *testing.T) {
	if !hasArg("parley-fake-buffered-late") {
		return
	}
	input, _ := io.ReadAll(os.Stdin)
	p := regexp.MustCompile(`(?m)^- Create exactly this file and no other protocol artifact: (.+)$`).FindStringSubmatch(string(input))
	idea := regexp.MustCompile(`(?m)^idea: (.+)$`).FindStringSubmatch(string(input))
	if len(p) != 2 || len(idea) != 2 {
		os.Exit(7)
	}
	time.Sleep(900 * time.Millisecond) // buffered: nothing on stdout until done
	body := "---\nagent: fake\nidea: " + idea[1] + "\nround: 1\ndate: 2026-06-12\n---\n\n## Summary\nx\n\n## Proposed approach\nx\n\n## Existing alternatives\nnone\n\n## Concerns / open questions\nx\n\n## Risks\nx\n"
	_ = os.WriteFile(p[1], []byte(body), 0o644)
	os.Exit(0)
}

// TestBufferedNeverEndingHelper is re-executed as a real child: a declared-
// buffered agent that emits nothing and never exits. The soft first-output guard
// must be suppressed; only the hard per-agent timeout bounds it.
func TestBufferedNeverEndingHelper(t *testing.T) {
	if !hasArg("parley-fake-buffered-hang") {
		return
	}
	_, _ = io.ReadAll(os.Stdin)
	time.Sleep(30 * time.Second)
	os.Exit(0)
}

func bufferedFakeAgent(helper, arg string, spec agents.Spec) agents.Discovery {
	spec.HeadlessArgs = []string{"-test.run=" + helper, "--", arg}
	spec.PromptMode = agents.PromptStdin
	return agents.Discovery{Spec: spec, Path: os.Args[0], Found: true}
}

// TestSupervisionForAgentBuffersStdoutDisablesSoftGuards locks the soft-guard
// suppression: a declared-buffered transport disables the first-output and stall
// guards while retaining the heartbeat (which never counts as activity).
func TestSupervisionForAgentBuffersStdoutDisablesSoftGuards(t *testing.T) {
	cfg := supervisionForAgent(agents.Discovery{Spec: agents.Spec{
		ID: "agy", FirstEventTimeoutMS: 5_000, StallTimeoutMS: 5_000, HeartbeatMS: 60_000, BuffersStdout: true,
	}}, time.Minute)
	if cfg.FirstEventTimeout != 0 {
		t.Fatalf("first-event=%v want 0 for a buffered transport", cfg.FirstEventTimeout)
	}
	if cfg.StallTimeout != 0 {
		t.Fatalf("stall=%v want 0 for a buffered transport", cfg.StallTimeout)
	}
	if cfg.HeartbeatInterval == 0 {
		t.Fatal("heartbeat must stay on (never counts as activity)")
	}

	// Contrast: a non-buffered agent keeps both soft guards.
	plain := supervisionForAgent(agents.Discovery{Spec: agents.Spec{
		ID: "claude", FirstEventTimeoutMS: 5_000, StallTimeoutMS: 5_000,
	}}, time.Minute)
	if plain.FirstEventTimeout == 0 || plain.StallTimeout == 0 {
		t.Fatalf("non-buffered soft guards must remain: %+v", plain)
	}
}

// TestBufferedQuietLateSuccessNotKilledBySoftGuard: a quiet, late-success child
// must survive a short first-event window when it is a declared-buffered
// transport — the soft guard is disabled, and the child's own exit-0 wins.
func TestBufferedQuietLateSuccessNotKilledBySoftGuard(t *testing.T) {
	root, idea, st := setupRunnerIdea(t)
	results := RunRoundOne(context.Background(), Options{
		Root: root, RunID: "test-run", Idea: idea, Task: "t",
		Agents: []agents.Discovery{bufferedFakeAgent("TestBufferedQuietLateSuccessHelper",
			"parley-fake-buffered-late",
			agents.Spec{ID: "fake", FirstEventTimeoutMS: 150, StallTimeoutMS: 150, HeartbeatMS: -1, BuffersStdout: true})},
		Timeout: 10 * time.Second, Store: st,
	})
	if len(results) != 1 {
		t.Fatalf("results=%d", len(results))
	}
	r := results[0]
	if !r.Success() || !r.ArtifactOK {
		t.Fatalf("buffered quiet late success must not be killed: %+v", r)
	}
	if r.FailureClass != "" {
		t.Fatalf("failure_class=%q want empty", r.FailureClass)
	}
}

// TestBufferedNeverEndingChildBoundedByHardTimeout: a buffered child that never
// outputs is NOT classified no_first_output (soft guard suppressed) and is
// instead bounded by the hard timeout, then cleaned up.
func TestBufferedNeverEndingChildBoundedByHardTimeout(t *testing.T) {
	root, idea, st := setupRunnerIdea(t)
	results := RunRoundOne(context.Background(), Options{
		Root: root, RunID: "test-run", Idea: idea, Task: "t",
		Agents: []agents.Discovery{bufferedFakeAgent("TestBufferedNeverEndingHelper",
			"parley-fake-buffered-hang",
			agents.Spec{ID: "fake", FirstEventTimeoutMS: 150, StallTimeoutMS: 150, HeartbeatMS: -1, BuffersStdout: true})},
		Timeout: 4 * time.Second, Store: st,
	})
	r := results[0]
	if r.Success() {
		t.Fatalf("a never-ending buffered child must eventually fail: %+v", r)
	}
	if r.FailureClass == "no_first_output" {
		t.Fatalf("buffered transport must NOT be killed by the soft first-output guard: %+v", r)
	}
	if r.FailureClass != "timeout" {
		t.Fatalf("failure_class=%q want timeout (hard timeout bounds the run)", r.FailureClass)
	}
}

// TestNonBufferedSilentChildStillKilledBySoftGuard is the contrast: without the
// declared buffer, the same silence is a soft-watchdog kill (no universal hang
// diagnosis, but the guard still fires where it is appropriate).
func TestNonBufferedSilentChildStillKilledBySoftGuard(t *testing.T) {
	root, idea, st := setupRunnerIdea(t)
	results := RunRoundOne(context.Background(), Options{
		Root: root, RunID: "test-run", Idea: idea, Task: "t",
		Agents: []agents.Discovery{fakeAgent("TestFakeSilentAgentHelper", agents.Spec{
			ID: "fake", FirstEventTimeoutMS: 150, StallTimeoutMS: -1, HeartbeatMS: -1,
		})},
		Timeout: 20 * time.Second, Store: st,
	})
	if results[0].FailureClass != "no_first_output" {
		t.Fatalf("non-buffered silence must be a soft-guard kill, got %+v", results[0])
	}
}
