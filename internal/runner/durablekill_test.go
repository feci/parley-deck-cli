package runner

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/procctl"
	"parley-deck-cli/internal/store"
)

func startedEvent(agentID string, pid int) store.Event {
	return store.Event{
		Time: time.Now().UTC(),
		Type: "agent.started",
		Data: map[string]any{
			"agent": agentID, "segment_id": "segment-0001",
			"pid": pid, "pgid": pid, "boot_id": procctl.CurrentBootID(),
			"proc_start": "x", "command": "y",
		},
	}
}

// A recorded process that is no longer alive is reconciled (stale badge cleared
// via a synthetic agent.failed), never signaled.
func TestDurableKillClearsStaleDeadPid(t *testing.T) {
	runDir := t.TempDir()
	st := store.New(runDir)
	_ = st.Append(startedEvent("codex", 2147480000)) // a pid that is not alive

	res := KillAgentDurable(st, "codex")
	if res.Killed {
		t.Fatalf("a dead recorded pid must not be 'killed', got %+v", res)
	}
	if !res.Cleared {
		t.Fatalf("a dead recorded pid should clear the stale badge, got %+v", res)
	}
	events, _ := st.Load()
	if !hasEventAgent(events, "agent.failed", "codex") {
		t.Fatalf("clear-stale should append a synthetic agent.failed, got %v", eventTypes(events))
	}
}

// Once a terminal event exists, durable kill is a no-op (idempotency / two-parley).
func TestDurableKillNoopAfterTerminal(t *testing.T) {
	runDir := t.TempDir()
	st := store.New(runDir)
	_ = st.Append(startedEvent("codex", 2147480000))
	_ = st.Append(store.Event{Time: time.Now().UTC(), Type: "agent.finished", Data: map[string]any{"agent": "codex"}})

	res := KillAgentDurable(st, "codex")
	if res.Killed || res.Cleared {
		t.Fatalf("after a terminal event, durable kill must be a no-op, got %+v", res)
	}
	events, _ := st.Load()
	n := 0
	for _, e := range events {
		if e.Type == "agent.failed" {
			n++
		}
	}
	if n != 0 {
		t.Fatalf("no synthetic failure should be written after a terminal event, got %d", n)
	}
}

// End-to-end on a real OS process: spawn a real agent in its own group, persist
// its identity in agent.started (as runAgent does), then durably kill it via the
// event log alone — exercising the full restart-path mechanism on this host.
func TestDurableKillEndToEndRealProcess(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("durable kill is unsupported on windows")
	}
	runDir := t.TempDir()
	st := store.New(runDir)
	agentDir := filepath.Join(runDir, "agents", "sleeper")
	if err := os.MkdirAll(agentDir, 0o755); err != nil {
		t.Fatal(err)
	}
	agent := agents.Discovery{Spec: agents.Spec{ID: "sleeper", HeadlessArgs: []string{"30"}}, Path: "/bin/sleep", Found: true}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan struct{})
	go func() {
		defer close(done)
		onStarted := func(sp procctl.Spawned) {
			_ = st.Append(store.Event{
				Time: time.Now().UTC(), Type: "agent.started",
				Data: map[string]any{
					"agent": "sleeper", "segment_id": "segment-0001",
					"pid": sp.PID, "pgid": sp.PGID, "boot_id": sp.BootID,
					"proc_start": sp.ProcStart, "command": sp.Command,
				},
			})
		}
		_, _ = execAgentProcess(ctx, runDir, "run-x", "sleeper", "m", agent,
			"", filepath.Join(agentDir, "stdout.log"), filepath.Join(agentDir, "stderr.log"), onStarted)
	}()

	// Wait until the process is recorded + alive (a separate "parley" would see this).
	var sp procctl.Spawned
	waitFor(t, 5*time.Second, func() bool {
		var found, terminal bool
		sp, _, found, terminal = latestAgentProc(st, "sleeper")
		return found && !terminal && procctl.Alive(sp)
	})

	res := KillAgentDurable(st, "sleeper")
	if !res.Killed {
		t.Fatalf("a live attributed process should be killed, got %+v", res)
	}
	select {
	case <-done: // execAgentProcess returned because the process was reaped
	case <-time.After(5 * time.Second):
		t.Fatal("process was not killed (execAgentProcess still running)")
	}
	if procctl.Alive(sp) {
		t.Fatal("process should be gone after durable kill")
	}
}

// SAFETY: a LIVE but unattributable recorded process (e.g. tampered/incomplete
// identity, simulating PID reuse) is NEVER signaled — durable kill refuses and
// leaves the process running.
func TestDurableKillRefusesUnattributableLiveProcess(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("durable kill unsupported on windows")
	}
	cmd := exec.Command("/bin/sleep", "30")
	procctl.SetNewProcessGroup(cmd)
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = cmd.Process.Kill() }()
	go func() { _ = cmd.Wait() }()
	sp := procctl.Capture(cmd, "m")
	if sp.BootID == "" {
		t.Skip("identity probe unavailable here")
	}

	runDir := t.TempDir()
	st := store.New(runDir)
	// Record the REAL live pid/pgid/boot but a WRONG start time → unattributable.
	_ = st.Append(store.Event{
		Time: time.Now().UTC(), Type: "agent.started",
		Data: map[string]any{
			"agent": "codex", "segment_id": "segment-0001",
			"pid": sp.PID, "pgid": sp.PGID, "boot_id": sp.BootID,
			"proc_start": "BOGUS-START-TIME", "command": sp.Command,
		},
	})

	res := KillAgentDurable(st, "codex")
	if res.Killed {
		t.Fatal("an unattributable live process must NOT be killed")
	}
	if !res.Failed {
		t.Fatalf("expected a verification-failed result, got %+v", res)
	}
	if !procctl.Alive(sp) {
		t.Fatal("the (unattributable) process must still be alive — we must not signal it")
	}
	events, _ := st.Load()
	if hasEventAgent(events, "agent.killed", "codex") || hasEventAgent(events, "agent.failed", "codex") {
		t.Fatalf("a refusal must not write a kill/clear event, got %v", eventTypes(events))
	}
}

func TestAgentLivenessStaleAndAbsent(t *testing.T) {
	runDir := t.TempDir()
	st := store.New(runDir)
	if got := AgentLiveness(st, "codex"); got != "" {
		t.Fatalf("no started event → liveness '', got %q", got)
	}
	_ = st.Append(startedEvent("codex", 2147480000)) // dead pid, no terminal
	if got := AgentLiveness(st, "codex"); got != "stale" {
		t.Fatalf("dead recorded pid → 'stale', got %q", got)
	}
	_ = st.Append(store.Event{Time: time.Now().UTC(), Type: "agent.finished", Data: map[string]any{"agent": "codex"}})
	if got := AgentLiveness(st, "codex"); got != "" {
		t.Fatalf("terminal event → liveness '', got %q", got)
	}
}
