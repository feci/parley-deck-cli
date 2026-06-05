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
			// A runner failure or a malformed event log halts the driver; capture
			// it in a durable blocking inbox note (consensus D4/AF3), not just
			// stderr, so an unattended --auto run leaves a recovery artifact.
			d.escalate(c, "driver-error",
				fmt.Sprintf("The auto-driver halted with an error while advancing %s:\n\n    %v\n\nInspect the run (events.jsonl / agent logs), fix the cause, then re-run 'parley run --auto'.",
					roundLabel(c.CurrentRound), err))
			return err
		}
		switch action {
		case ActionPromoted:
			fmt.Fprintf(d.cfg.Out, "driver: opened %s (idea status=%s)\n", roundLabel(c.CurrentRound), c.IdeaStatus)
			deadline = time.Now().Add(roundDeadline)
		case ActionConsensusDrafted:
			fmt.Fprintf(d.cfg.Out, "driver: drafted consensus.md for %s\n", d.cfg.IdeaSlug)
			deadline = time.Now().Add(roundDeadline)
		case ActionSignoffsRequested:
			fmt.Fprintf(d.cfg.Out, "driver: requested missing consensus signoffs\n")
			deadline = time.Now().Add(roundDeadline)
		case ActionReopened:
			fmt.Fprintf(d.cfg.Out, "driver: consensus blocked — reopened %s (idea status=%s)\n", roundLabel(c.CurrentRound), c.IdeaStatus)
			deadline = time.Now().Add(roundDeadline)
		case ActionAwait:
			if time.Now().After(deadline) {
				return d.escalateDeadline(c)
			}
			time.Sleep(2 * time.Second)
		case ActionFinalized:
			fmt.Fprintf(d.cfg.Out, "driver: authored FINAL.md; idea %s is final (implementation auto-drive is a later slice)\n", d.cfg.IdeaSlug)
			return nil
		case ActionConsensus:
			fmt.Fprintf(d.cfg.Out, "driver: cross-review complete at %s; next step is `parley consensus draft %s` (consensus auto-drive not wired)\n", roundLabel(c.CurrentRound), d.cfg.IdeaSlug)
			return nil
		case ActionSurfaceOnly:
			fmt.Fprintf(d.cfg.Out, "driver: auto-advance not enabled here (needs --auto and local-dir transport); idea left at %s\n", c.IdeaStatus)
			return nil
		case ActionEscalated:
			return nil
		}
	}
}

// escalateDeadline writes a blocking inbox escalation when a round does not
// complete within the deadline, then halts (consensus D12).
func (d *Driver) escalateDeadline(c Cursor) error {
	d.escalate(c, "round-deadline", fmt.Sprintf(
		"Round %s did not complete within %s. The driver halted rather than spin.\n\nSome participant artifact under %s/%s is missing or did not validate, and no round.completed event was reconstructable.\n\nInspect the round, re-run the missing participant, or adjust the roster, then re-run 'parley run --auto'.",
		roundLabel(c.CurrentRound), roundDeadline, d.cfg.IdeaSlug, roundLabel(c.CurrentRound)))
	return nil
}

// escalate writes a durable blocking inbox note and logs it. topic is a short
// kebab slug; body is the human-facing explanation + recovery steps. Centralizes
// the escalation path so every halt (deadline, malformed log, runner failure)
// leaves a recovery artifact (consensus D11/D12/AF3).
func (d *Driver) escalate(c Cursor, topic, body string) {
	note := fmt.Sprintf(`---
from: claude
to: user
idea: %s
phase: %s
blocking: yes
date: %s
---

## Question / blocker
%s
`, d.cfg.IdeaSlug, c.Phase, time.Now().UTC().Format("2006-01-02"), body)

	inboxDir := filepath.Join(d.cfg.RunDir, "..", "..", "inbox")
	if abs, err := filepath.Abs(inboxDir); err == nil {
		inboxDir = abs
	}
	_ = os.MkdirAll(inboxDir, 0o755)
	name := fmt.Sprintf("claude-to-user_%s_%s.md", d.cfg.IdeaSlug, topic)
	_ = os.WriteFile(filepath.Join(inboxDir, name), []byte(note), 0o644)
	fmt.Fprintf(d.cfg.Out, "driver: halted (%s); wrote blocking escalation %s\n", topic, name)
}

// acquireLock atomically creates a PID lock file (O_EXCL), refusing if a live
// process already holds it. A stale lock (dead PID) is reclaimed and the atomic
// create retried, so two racing starters cannot both win (consensus D10/AF1).
// Release removes the file only if it still carries this process' token.
func acquireLock(path string) (func(), error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	token := strconv.Itoa(os.Getpid())
	for attempt := 0; attempt < 2; attempt++ {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
		if err == nil {
			_, _ = f.WriteString(token)
			_ = f.Close()
			return func() { releaseLock(path, token) }, nil
		}
		if !os.IsExist(err) {
			return nil, err
		}
		data, rerr := os.ReadFile(path)
		if rerr != nil {
			return nil, fmt.Errorf("driver.lock present but unreadable: %w", rerr)
		}
		// Held by us (this PID already owns it) or by another live process → refuse.
		// Only a DIFFERENT, dead PID is a stale lock we may reclaim.
		if pid, perr := strconv.Atoi(strings.TrimSpace(string(data))); perr == nil {
			if pid == os.Getpid() || processAlive(pid) {
				return nil, fmt.Errorf("driver.lock held by pid %d", pid)
			}
		}
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}
	return nil, fmt.Errorf("could not acquire driver.lock at %s", path)
}

// releaseLock removes the lock only if it still holds our ownership token, so we
// never delete a lock another process re-acquired after a stale reclaim.
func releaseLock(path, token string) {
	if data, err := os.ReadFile(path); err == nil && strings.TrimSpace(string(data)) == token {
		_ = os.Remove(path)
	}
}

func processAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return proc.Signal(syscall.Signal(0)) == nil
}
