package driver

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/store"
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

	// Persist the policy before the first transition. Saved limits remain
	// authoritative on a new Run, including when its flags are omitted.
	binding, err := budget.EnsureStepBinding(ctx, d.cfg.Root, d.cfg.IdeaSlug, d.cfg.MaxDriverSteps, d.cfg.MaxWallClock)
	if err != nil {
		return fmt.Errorf("persistent driver budget: %w", err)
	}
	if binding != nil {
		private := *d
		private.cfg.MaxDriverSteps = binding.Policy.MaxSteps
		private.cfg.MaxWallClock = time.Duration(binding.Policy.WallClockNS)
		d = &private
	}

	deadline := time.Now().Add(roundDeadline)
	start := time.Now()
	steps := 0
	var last Cursor
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if binding != nil {
			state, err := binding.Store.Inspect(ctx)
			if err != nil {
				return fmt.Errorf("read lifetime driver budget: %w", err)
			}
			steps, start = budget.StepCount(state), state.StartedAt
		}

		// LE-5: enforce the loop ceilings BEFORE advancing. A breach escalates (durable
		// inbox note) and halts — it never marks the idea complete.
		if reason := d.loopBudgetBreach(steps, start); reason != "" {
			return d.escalateLoopBudget(last, reason)
		}

		action, c, err := d.Advance(ctx)
		last = c
		// Report the durable charge even if the action failed after reserving it.
		if binding != nil {
			state, inspectErr := binding.Store.Inspect(ctx)
			if inspectErr != nil {
				return fmt.Errorf("read charged driver step: %w", inspectErr)
			}
			steps, start = budget.StepCount(state), state.StartedAt
			if emitErr := d.emitLoopBudget(steps, start); emitErr != nil {
				return emitErr
			}
		}
		if err != nil {
			// A runner failure or a malformed event log halts the driver; capture
			// it in a durable blocking inbox note (consensus D4/AF3), not just
			// stderr, so an unattended --auto run leaves a recovery artifact.
			d.escalate(c, "driver-error",
				fmt.Sprintf("The auto-driver halted with an error while advancing %s:\n\n    %v\n\nInspect the run (events.jsonl / agent logs), fix the cause, then re-run 'parley run --auto'.",
					roundLabel(c.CurrentRound), err))
			return err
		}
		// LE-5: count progress Advances and record budget burn for the TUI/state.
		if binding == nil && isProgressAction(action) {
			steps++
			if err := d.emitLoopBudget(steps, start); err != nil {
				return err
			}
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
		case ActionFinalized:
			fmt.Fprintf(d.cfg.Out, "driver: authored FINAL.md for %s\n", d.cfg.IdeaSlug)
			deadline = time.Now().Add(roundDeadline)
		case ActionImplemented:
			fmt.Fprintf(d.cfg.Out, "driver: ran the implementer (idea status=%s)\n", c.IdeaStatus)
			deadline = time.Now().Add(roundDeadline)
		case ActionReviewOpened:
			fmt.Fprintf(d.cfg.Out, "driver: opened a review round\n")
			deadline = time.Now().Add(roundDeadline)
		case ActionReviewDrafted:
			fmt.Fprintf(d.cfg.Out, "driver: drafted review/consensus.md\n")
			deadline = time.Now().Add(roundDeadline)
		case ActionFixup:
			fmt.Fprintf(d.cfg.Out, "driver: ran a fix-up cycle and opened the next review round\n")
			deadline = time.Now().Add(roundDeadline)
		case ActionAwait:
			if time.Now().After(deadline) {
				return d.escalateDeadline(c)
			}
			time.Sleep(2 * time.Second)
		case ActionComplete:
			fmt.Fprintf(d.cfg.Out, "driver: idea %s is complete (review consensus clean); ready to merge — the driver does not merge/push/release\n", d.cfg.IdeaSlug)
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

// isProgressAction reports whether an action advanced the protocol (so it counts toward
// the driver-step budget and resets the per-tick deadline), as opposed to awaiting or
// terminating.
func isProgressAction(a Action) bool {
	switch a {
	case ActionAwait, ActionComplete, ActionConsensus, ActionSurfaceOnly, ActionEscalated:
		return false
	}
	return true
}

// loopBudgetBreach returns a non-empty reason when a non-zero loop ceiling is exceeded
// (LE-5). 0 ceilings are unlimited. A monetary ceiling requires complete
// per-invocation prices from agent.usage events and stops on unknown accounting.
func (d *Driver) loopBudgetBreach(steps int, start time.Time) string {
	if d.cfg.MaxDriverSteps > 0 && steps >= d.cfg.MaxDriverSteps {
		return fmt.Sprintf("driver-step budget exhausted (%d/%d steps)", steps, d.cfg.MaxDriverSteps)
	}
	if d.cfg.MaxWallClock > 0 {
		if elapsed := time.Since(start); elapsed >= d.cfg.MaxWallClock {
			return fmt.Sprintf("wall-clock budget exhausted (%s/%s)", elapsed.Round(time.Second), d.cfg.MaxWallClock)
		}
	}
	if d.cfg.MaxCostUSD > 0 {
		cost := d.loopCost()
		if !cost.complete {
			return "cost budget cannot be enforced: " + cost.reason
		}
		if cost.usd >= d.cfg.MaxCostUSD {
			return fmt.Sprintf("cost budget exhausted ($%.2f/$%.2f)", cost.usd, d.cfg.MaxCostUSD)
		}
	}
	return ""
}

// emitLoopBudget records budget burn after a progress step so the TUI/state can show it.
// Cost is always reported for observability (F-T2-2); only enforcement is gated by
// MaxCostUSD > 0 (in loopBudgetBreach).
func (d *Driver) emitLoopBudget(steps int, start time.Time) error {
	cost := d.loopCost()
	var total any
	coverage := "incomplete"
	if cost.complete {
		total, coverage = cost.usd, "complete"
	}
	return d.cfg.Events.Append(store.Event{
		Time: time.Now().UTC(),
		Type: "loop.budget",
		Data: map[string]any{
			"idea":              d.cfg.IdeaSlug,
			"steps":             steps,
			"max_driver_steps":  d.cfg.MaxDriverSteps,
			"elapsed_ms":        time.Since(start).Milliseconds(),
			"max_wall_clock_ms": d.cfg.MaxWallClock.Milliseconds(),
			"cost_usd":          total,
			"cost_coverage":     coverage,
			"max_cost_usd":      d.cfg.MaxCostUSD,
		},
	})
}

type loopCostSummary struct {
	usd      float64
	complete bool
	reason   string
}

// A normalized usage event identifies one invocation, including a retry. Only
// identical replays may be deduplicated; a missing price or contradictory replay
// cannot establish the total against a monetary ceiling.
func (d *Driver) loopCost() loopCostSummary {
	if !d.cfg.Events.Enabled() {
		return loopCostSummary{reason: "run event store is unavailable"}
	}
	evs, err := d.cfg.Events.Load()
	if err != nil {
		return loopCostSummary{reason: "run event accounting is missing or unreadable"}
	}
	result := loopCostSummary{complete: true}
	seen := map[string]map[string]any{}
	for _, e := range evs {
		if e.Type != "agent.usage" {
			continue
		}
		id, _ := e.Data["invocation_id"].(string)
		if strings.TrimSpace(id) == "" {
			return loopCostSummary{reason: "usage event has no invocation identity"}
		}
		if prior, ok := seen[id]; ok {
			if !reflect.DeepEqual(prior, e.Data) {
				return loopCostSummary{reason: "conflicting usage for one invocation"}
			}
			continue
		}
		seen[id] = e.Data
		cost, ok := e.Data["cost_usd"].(float64)
		if !ok || math.IsNaN(cost) || math.IsInf(cost, 0) || cost < 0 || cost > 1_000_000 {
			return loopCostSummary{reason: "invocation cost is unknown or invalid"}
		}
		result.usd += cost
		if math.IsInf(result.usd, 0) {
			return loopCostSummary{reason: "cost total is out of range"}
		}
	}
	return result
}

// escalateLoopBudget writes a durable blocking inbox note when a loop ceiling is hit and
// halts cleanly (LE-5: budget hit = escalate, never complete).
func (d *Driver) escalateLoopBudget(c Cursor, reason string) error {
	d.escalate(c, "loop-budget", fmt.Sprintf(
		"The auto-driver hit a loop budget and halted rather than continue:\n\n    %s\n\nThis is a safety ceiling (loop engineering: a budget hit escalates, it does not mark the idea complete). Preserve the saved policy and charges. A new run or changed flags cannot extend a frozen lifetime ceiling; an explicit operator extension or accounting migration is required before continuing.",
		reason))
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
		// Empty or unparseable content means a racing acquirer just created the file
		// with O_EXCL but has not yet written its PID token (the create precedes the
		// token write). Treat it as HELD and refuse — never reclaim it as "stale", or
		// two racers could both remove-and-recreate and both win (review-01 lock race).
		pid, perr := strconv.Atoi(strings.TrimSpace(string(data)))
		if perr != nil {
			return nil, fmt.Errorf("driver.lock present but not yet owned (racing acquirer or corrupt lock); refusing")
		}
		// Held by us (this PID already owns it) or by another live process → refuse.
		// Only a DIFFERENT, dead PID is a stale lock we may reclaim.
		if pid == os.Getpid() || processAlive(pid) {
			return nil, fmt.Errorf("driver.lock held by pid %d", pid)
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
