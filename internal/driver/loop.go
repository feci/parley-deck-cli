package driver

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// roundDeadline bounds how long a single tick waits for an incomplete round
// before escalating (consensus D12). A round whose agents never finish must
// terminate, not spin.
const roundDeadline = 30 * time.Minute

// Run drives the idea forward under a mandatory advisory lock until the round
// phase is exhausted (cross-review complete → consensus is a later slice), the
// gate disables auto-advance, the deadline trips, or ctx is cancelled. The
// contract is single-driver + idempotent re-entry, NOT multi-writer (consensus
// D10): if the lock is held, Run stops cleanly without driving.
func (d *Driver) Run(ctx context.Context) error {
	release, err := acquireLock(filepath.Join(d.cfg.RunDir, "driver.lock"))
	if err != nil {
		fmt.Fprintf(d.cfg.Out, "driver: not auto-advancing — %v\n", err)
		return nil
	}
	defer release()

	deadline := time.Now().Add(roundDeadline)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		action, c, err := d.Advance(ctx)
		if err != nil {
			fmt.Fprintf(d.cfg.Out, "driver: halting — %v\n", err)
			return err
		}
		switch action {
		case ActionPromoted:
			fmt.Fprintf(d.cfg.Out, "driver: opened %s (idea status=%s)\n", roundLabel(c.CurrentRound), c.IdeaStatus)
			deadline = time.Now().Add(roundDeadline)
		case ActionAwait:
			if time.Now().After(deadline) {
				return d.escalateDeadline(c)
			}
			time.Sleep(2 * time.Second)
		case ActionConsensus:
			fmt.Fprintf(d.cfg.Out, "driver: cross-review complete at %s; next step is `parley consensus draft %s` (consensus auto-drive is a later slice)\n", roundLabel(c.CurrentRound), d.cfg.IdeaSlug)
			return nil
		case ActionSurfaceOnly:
			fmt.Fprintf(d.cfg.Out, "driver: auto-advance not enabled here (needs --auto and local-dir transport, round phase); idea left at %s\n", c.IdeaStatus)
			return nil
		case ActionEscalated:
			return nil
		}
	}
}

// escalateDeadline writes a blocking inbox escalation when a round does not
// complete within the deadline, then halts (consensus D12).
func (d *Driver) escalateDeadline(c Cursor) error {
	note := fmt.Sprintf(`---
from: claude
to: user
idea: %s
phase: %s
blocking: yes
date: %s
---

## Question
Round %s did not complete within %s. The driver halted rather than spin.

## Context
Some participant artifact under %s/%s is missing or did not validate, and no
round.completed event was reconstructable. The auto-driver stops here.

## What I need from you
Inspect the round, re-run the missing participant, or adjust the roster, then
re-run 'parley run --auto'.
`, d.cfg.IdeaSlug, c.Phase, time.Now().UTC().Format("2006-01-02"), roundLabel(c.CurrentRound), roundDeadline, d.cfg.IdeaSlug, roundLabel(c.CurrentRound))

	inboxDir := filepath.Join(d.cfg.RunDir, "..", "..", "inbox")
	if abs, err := filepath.Abs(inboxDir); err == nil {
		inboxDir = abs
	}
	_ = os.MkdirAll(inboxDir, 0o755)
	name := fmt.Sprintf("claude-to-user_%s_round-deadline.md", d.cfg.IdeaSlug)
	_ = os.WriteFile(filepath.Join(inboxDir, name), []byte(note), 0o644)
	fmt.Fprintf(d.cfg.Out, "driver: %s did not complete within %s; wrote blocking escalation %s\n", roundLabel(c.CurrentRound), roundDeadline, name)
	return nil
}

// acquireLock writes a PID lock file, refusing if a live process already holds
// it. A stale lock (dead PID) is reclaimed. Returns a release func on success.
func acquireLock(path string) (func(), error) {
	if data, err := os.ReadFile(path); err == nil {
		if pid, perr := strconv.Atoi(strings.TrimSpace(string(data))); perr == nil && pid != os.Getpid() && processAlive(pid) {
			return nil, fmt.Errorf("driver.lock held by live pid %d", pid)
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(path, []byte(strconv.Itoa(os.Getpid())), 0o644); err != nil {
		return nil, err
	}
	return func() { _ = os.Remove(path) }, nil
}

func processAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return proc.Signal(syscall.Signal(0)) == nil
}
