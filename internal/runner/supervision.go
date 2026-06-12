package runner

import (
	"errors"
	"sync/atomic"
	"time"

	"parley-deck-cli/internal/agents"
)

// Supervision implements the kindly-derived three-layer agent watchdog
// (runner-hardening-kindly D1-D4): a first-output grace window, a stall guard
// after first output, and periodic persisted heartbeats. Activity is tracked
// in-process by counting writers wrapped around the agent's stdout/stderr
// sinks (exec mode) or by ACP protocol events — never by filesystem probing.

// Watchdog sentinel errors: runAgent keys retry/classification on these.
var (
	errNoFirstOutput = errors.New("no output within the first-output grace window")
	errStalled       = errors.New("no output growth within the stall window")
)

// SupervisionConfig carries the effective per-agent watchdog windows.
// A zero or negative duration disables that guard.
type SupervisionConfig struct {
	FirstEventTimeout time.Duration
	StallTimeout      time.Duration
	HeartbeatInterval time.Duration
}

const (
	defaultFirstEventTimeoutMS = 120_000
	defaultStallTimeoutMS      = 1_800_000
	defaultHeartbeatMS         = 60_000
)

// supervisionForAgent derives the effective windows from the agent spec
// (consensus D2): 0 means "use the default"; a negative spec value disables
// the guard (the TOML layer maps an explicit 0 override to -1). The stall
// window is clamped under the hard timeout so the stall guard always fires
// first and yields the more actionable failure class.
func supervisionForAgent(agent agents.Discovery, hardTimeout time.Duration) SupervisionConfig {
	pick := func(specMS, defMS int) time.Duration {
		switch {
		case specMS < 0:
			return 0 // explicitly disabled
		case specMS == 0:
			return time.Duration(defMS) * time.Millisecond
		default:
			return time.Duration(specMS) * time.Millisecond
		}
	}
	cfg := SupervisionConfig{
		FirstEventTimeout: pick(agent.FirstEventTimeoutMS, defaultFirstEventTimeoutMS),
		StallTimeout:      pick(agent.StallTimeoutMS, defaultStallTimeoutMS),
		HeartbeatInterval: pick(agent.HeartbeatMS, defaultHeartbeatMS),
	}
	if cfg.StallTimeout > 0 && hardTimeout > 0 && cfg.StallTimeout >= hardTimeout {
		cfg.StallTimeout = hardTimeout - time.Second
		if cfg.StallTimeout <= 0 {
			cfg.StallTimeout = 0
		}
	}
	return cfg
}

// activityTracker accumulates output bytes and the last-activity instant.
// Counting writers (and the ACP handler) call Mark; the supervisor reads
// snapshots. All fields are atomics — writers and the supervisor goroutine
// never share a lock.
type activityTracker struct {
	stdoutBytes  atomic.Int64
	stderrBytes  atomic.Int64
	lastUnixNano atomic.Int64
}

func (t *activityTracker) Mark(stream string, n int) {
	if n <= 0 {
		return
	}
	switch stream {
	case "stderr":
		t.stderrBytes.Add(int64(n))
	default:
		t.stdoutBytes.Add(int64(n))
	}
	t.lastUnixNano.Store(time.Now().UnixNano())
}

// MarkEvent records non-byte activity (ACP protocol events).
func (t *activityTracker) MarkEvent() {
	t.lastUnixNano.Store(time.Now().UnixNano())
}

type activitySnapshot struct {
	StdoutBytes  int64
	StderrBytes  int64
	LastActivity time.Time
}

func (t *activityTracker) snapshot() activitySnapshot {
	var last time.Time
	if ns := t.lastUnixNano.Load(); ns > 0 {
		last = time.Unix(0, ns)
	}
	return activitySnapshot{
		StdoutBytes:  t.stdoutBytes.Load(),
		StderrBytes:  t.stderrBytes.Load(),
		LastActivity: last,
	}
}

// countingWriter wraps an output sink, attributing written bytes to the
// tracker. The underlying writes are unchanged — zero extra filesystem I/O.
type countingWriter struct {
	w      interface{ Write([]byte) (int, error) }
	t      *activityTracker
	stream string
}

func (c *countingWriter) Write(p []byte) (int, error) {
	n, err := c.w.Write(p)
	c.t.Mark(c.stream, n)
	return n, err
}

// supervisionHooks lets the wait loop append typed events without knowing
// about the store: the runner injects closures that build the payloads.
// onWatchdog MUST append its event BEFORE the kill fires (consensus D1) —
// waitSupervised guarantees that ordering by calling it before kill().
type supervisionHooks struct {
	onHeartbeat func(snap activitySnapshot, elapsed time.Duration)
	onWatchdog  func(kind string, snap activitySnapshot, elapsed time.Duration)
}

// waitSupervised waits for the agent like the previous wait select, adding the
// three supervision layers. kill must terminate the whole process tree; the
// function always drains waitErr before returning so the child is reaped.
// Returns errNoFirstOutput / errStalled for watchdog kills, ctx.Err() for the
// outer timeout/cancel, or the process's own wait error.
func waitSupervised(done <-chan struct{}, ctxErr func() error, waitErr <-chan error, kill func(), act *activityTracker, cfg SupervisionConfig, hooks supervisionHooks) error {
	started := time.Now()
	tick := time.NewTicker(time.Second)
	defer tick.Stop()
	lastHeartbeat := started
	for {
		select {
		case e := <-waitErr:
			return e
		case <-done:
			kill()
			<-waitErr
			return ctxErr()
		case <-tick.C:
			now := time.Now()
			snap := act.snapshot()
			elapsed := now.Sub(started)
			if snap.LastActivity.IsZero() {
				if cfg.FirstEventTimeout > 0 && elapsed >= cfg.FirstEventTimeout {
					if hooks.onWatchdog != nil {
						hooks.onWatchdog("no_first_output", snap, elapsed)
					}
					kill()
					<-waitErr
					return errNoFirstOutput
				}
			} else if cfg.StallTimeout > 0 && now.Sub(snap.LastActivity) >= cfg.StallTimeout {
				if hooks.onWatchdog != nil {
					hooks.onWatchdog("stalled", snap, elapsed)
				}
				kill()
				<-waitErr
				return errStalled
			}
			if cfg.HeartbeatInterval > 0 && now.Sub(lastHeartbeat) >= cfg.HeartbeatInterval {
				lastHeartbeat = now
				if hooks.onHeartbeat != nil {
					hooks.onHeartbeat(snap, elapsed) // never counts as activity
				}
			}
		}
	}
}
