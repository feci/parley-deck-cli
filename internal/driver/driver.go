package driver

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/store"
)

// RoundRunner launches cross-review round N, writing round-NN/<agent>.md for each
// participant. The production adapter (NewRunnerAdapter) wraps runner.RunRound;
// unit tests inject a fake. This is the driver core's only agent-launch
// dependency, which keeps Advance testable without live agents (consensus D14).
type RoundRunner interface {
	RunRound(ctx context.Context, round int) error
}

// Action is the outcome of one Advance tick.
type Action string

const (
	ActionPromoted          Action = "promoted"           // opened the next round
	ActionAwait             Action = "await"              // round present but incomplete; wait
	ActionConsensus         Action = "consensus-ready"    // cross-review budget spent, gate unwired (slice-1 stop)
	ActionConsensusDrafted  Action = "consensus-drafted"  // authored consensus.md → PhaseConsensus
	ActionSignoffsRequested Action = "signoffs-requested" // invoked missing signers
	ActionFinalized         Action = "finalized"          // authored FINAL.md → PhaseFinal
	ActionReopened          Action = "reopened"           // BLOCK → reopened a cross-review round
	ActionSurfaceOnly       Action = "surface-only"       // not auto-drivable here (gate/phase)
	ActionEscalated         Action = "escalated"          // halted; caller writes escalation
)

// Config carries everything the driver needs; it never imports internal/app.
type Config struct {
	IdeaDir           string      // …/ideas/<slug>
	IdeaSlug          string      // idea slug (for artifact validation + events)
	Participants      []string    // expected participant set (00-prompt participants)
	RunDir            string      // …/runs/<runID> (cursor + lock live here)
	Root              string      // workspace root (for per-tick transport read, D8)
	Events            store.Store // run event store (Append + Load events.jsonl)
	CrossReviewRounds int         // default 1 (one independent + N cross-review rounds)
	MaxRounds         int         // circuit breaker, default 4
	Auto              bool        // --auto flag; transport is read from disk per tick
	// Consensus, when set, enables the consensus gate (slice 2). nil keeps the
	// slice-1 behavior: stop at the consensus boundary (ActionConsensus).
	Consensus ConsensusOps
	// Impl, when set, enables the implementation/review gate (driver-impl-phase).
	// nil keeps the prior behavior: stop at FINAL.md (surface-only).
	Impl ImplOps
	// AutoImplement gates the code-writing phases (Implement/Fixup): true only when
	// the idea opted in (00-prompt auto_implement) AND --no-implement was not set.
	AutoImplement bool
	// MaxFixupCycles bounds the review→fix-up loop (default 3).
	MaxFixupCycles int
	Out            io.Writer // progress output (nil → discard)
}

// Driver advances one idea through the deliberation phases via Advance ticks.
type Driver struct {
	cfg    Config
	runner RoundRunner
}

// New constructs a Driver, applying defaults.
func New(cfg Config, r RoundRunner) *Driver {
	if cfg.MaxRounds <= 0 {
		cfg.MaxRounds = 4
	}
	if cfg.CrossReviewRounds < 0 {
		cfg.CrossReviewRounds = 1
	}
	if cfg.MaxFixupCycles <= 0 {
		cfg.MaxFixupCycles = 3
	}
	if cfg.Out == nil {
		cfg.Out = io.Discard
	}
	return &Driver{cfg: cfg, runner: r}
}

func (d *Driver) cursorPath() string { return filepath.Join(d.cfg.RunDir, "driver.json") }

// saveCursor is a seam so tests can force a cursor-save failure.
var saveCursor = func(c Cursor, path string) error { return c.Save(path) }

// commitCursor persists a phase-changing cursor and, only after a successful
// save, appends a best-effort run.phase event so the TUI sees phase flips
// event-driven (consensus tui-protocol-visibility D4). A failed save is
// returned so the branch escalates instead of letting the event log claim a
// phase the cursor never persisted; a failed event append is healed by the
// snapshot's disk reconcile.
func (d *Driver) commitCursor(c Cursor, action Action, previous Phase) error {
	if err := saveCursor(c, d.cursorPath()); err != nil {
		return fmt.Errorf("save driver cursor: %w", err)
	}
	_ = d.cfg.Events.Append(store.Event{
		Time: time.Now().UTC(),
		Type: "run.phase",
		Data: map[string]any{
			"idea":           d.cfg.IdeaSlug,
			"run_id":         filepath.Base(d.cfg.RunDir),
			"action":         string(action),
			"phase":          string(c.Phase),
			"previous_phase": string(previous),
			"current_round":  c.CurrentRound,
			"round_label":    roundLabel(c.CurrentRound),
			"idea_status":    c.IdeaStatus,
			"rounds_run":     c.RoundsRun,
			"max_rounds":     c.MaxRounds,
			"source":         "driver",
		},
	})
	return nil
}

// autoDriveEnabled reports whether the driver should auto-advance. As of 1.27.0
// this is transport-independent (relaxing the original local-dir-only gate from
// consensus D8): the canonical artifacts (rounds, consensus, FINAL, …) are the
// source of truth under every transport, so auto-drive advances them on
// github-pr / gitlab-mr too. The driver does NOT create PR/MR branches — that
// mirroring stays a manual, ergonomic step. Only the --auto flag (default on;
// --no-auto to disable) gates auto-drive now.
func (d *Driver) autoDriveEnabled() bool {
	return d.cfg.Auto
}

// Advance performs ONE re-entrant, idempotent tick. Disk is authoritative: it
// rebuilds the cursor from disk, then runs at most one gated action. Every branch
// is a no-op when its output already exists, so a duplicated tick or crash-restart
// cannot double-produce.
func (d *Driver) Advance(ctx context.Context) (Action, Cursor, error) {
	c := Rebuild(d.cfg.IdeaDir, d.cfg.MaxRounds)

	// Auto gate (1.27.0): auto-drive is transport-independent; only --auto/--no-auto
	// decides. See autoDriveEnabled.
	if !d.autoDriveEnabled() {
		return ActionSurfaceOnly, c, nil
	}
	switch c.Phase {
	case PhaseRound:
		return d.advanceRound(ctx, c)
	case PhaseConsensus:
		return d.advanceConsensus(ctx, c)
	case PhaseFinal:
		return d.advanceFinal(ctx, c)
	case PhaseImpl:
		return d.advanceImpl(ctx, c)
	case PhaseReview:
		return d.advanceReview(ctx, c)
	default:
		// PhaseDone / PhaseBlocked: nothing to auto-drive (ready to merge / halted).
		return ActionSurfaceOnly, c, nil
	}
}

// advanceRound drives the round phase: promote a completed round to the next
// cross-review round, or (budget spent) draft consensus / stop.
func (d *Driver) advanceRound(ctx context.Context, c Cursor) (Action, Cursor, error) {
	done, err := d.roundComplete(c.CurrentRound)
	if err != nil {
		return ActionEscalated, c, err
	}
	if !done {
		return ActionAwait, c, nil
	}
	// Round complete. Cross-review policy: rounds 1..(1+CrossReviewRounds).
	if c.CurrentRound >= 1+d.cfg.CrossReviewRounds {
		if d.cfg.Consensus == nil {
			return ActionConsensus, c, nil // gate not wired (slice-1 stop)
		}
		if err := d.cfg.Consensus.Draft(ctx); err != nil {
			return ActionEscalated, c, fmt.Errorf("draft consensus: %w", err)
		}
		c.Phase = PhaseConsensus
		c.IdeaStatus = "consensus"
		c.UpdatedAt = nowRFC3339()
		if err := d.commitCursor(c, ActionConsensusDrafted, PhaseRound); err != nil {
			return ActionEscalated, c, err
		}
		return ActionConsensusDrafted, c, nil
	}

	next := c.CurrentRound + 1
	// Idempotent: if the next round already completed (re-entry / crash after the
	// run but before the cursor save), do not re-dispatch — just advance.
	if nextDone, _ := d.roundComplete(next); !nextDone {
		if err := d.runner.RunRound(ctx, next); err != nil {
			return ActionEscalated, c, fmt.Errorf("run %s: %w", roundLabel(next), err)
		}
	}
	if err := setIdeaStatus(d.cfg.IdeaDir, roundLabel(next)); err != nil {
		return ActionEscalated, c, fmt.Errorf("set idea status: %w", err)
	}
	c.CurrentRound = next
	c.RoundsRun = next
	c.UpdatedAt = nowRFC3339()
	if err := d.commitCursor(c, ActionPromoted, PhaseRound); err != nil {
		return ActionEscalated, c, err
	}
	return ActionPromoted, c, nil
}

// roundComplete is the two-signal gate with reconciliation (consensus D4): all
// participant artifacts present + valid AND a terminal round.completed event,
// reconstructing the event from validated artifacts when events.jsonl is
// truncated. Returns (false, nil) for an incomplete round and (false, err) only
// for a malformed event log (→ escalate). Re-emission never fires on file
// presence alone — every artifact must validate first.
func (d *Driver) roundComplete(round int) (bool, error) {
	label := roundLabel(round)
	roundDir := filepath.Join(d.cfg.IdeaDir, label)
	if !dirExists(roundDir) {
		return false, nil
	}
	if len(d.cfg.Participants) == 0 {
		return false, nil
	}
	for _, participant := range d.cfg.Participants {
		artifact := filepath.Join(roundDir, participant+".md")
		info, err := os.Stat(artifact)
		if err != nil || !info.Mode().IsRegular() {
			return false, nil
		}
		if err := runner.ValidateRoundArtifact(artifact, participant, d.cfg.IdeaSlug, round); err != nil {
			return false, nil // present but not yet valid → incomplete, not an error
		}
		if round >= 2 {
			// D4 cross-review evidence: responding-to frontmatter AND a per-agent
			// `### @<other>` heading for every other participant (the runner's
			// BuildRoundPrompt now emits these). Absence → not a real cross-review
			// response yet → incomplete.
			if !hasRespondingTo(artifact) {
				return false, nil
			}
			if err := validateCrossReviewBody(artifact, participant, d.cfg.Participants); err != nil {
				return false, nil
			}
		}
	}
	terminal, err := d.terminalRoundEvent(label)
	if err != nil {
		return false, err // malformed event log → escalate
	}
	switch terminal {
	case "round.completed":
		return true, nil
	case "round.incomplete":
		return false, nil // authoritative block
	default:
		// Artifacts all valid but no terminal event (events.jsonl missing/
		// truncated): reconstruct round.completed, scoped to the current run.
		_ = d.cfg.Events.Append(store.Event{
			Time: time.Now().UTC(),
			Type: "round.completed",
			Data: map[string]any{
				"idea":          d.cfg.IdeaSlug,
				"round":         label,
				"completed":     len(d.cfg.Participants),
				"total":         len(d.cfg.Participants),
				"reconstructed": true,
			},
		})
		return true, nil
	}
}

// terminalRoundEvent returns the last "round.completed"/"round.incomplete" event
// for the given round label in the current run, or "" if none. A malformed event
// log returns an error (caller escalates rather than guessing).
func (d *Driver) terminalRoundEvent(label string) (string, error) {
	events, err := d.cfg.Events.Load()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", fmt.Errorf("read event log: %w", err)
	}
	terminal := ""
	for _, e := range events {
		if e.Type != "round.completed" && e.Type != "round.incomplete" {
			continue
		}
		if r, _ := e.Data["round"].(string); r != label {
			continue
		}
		// D4: the terminal event is scoped to idea+round. Tolerate a missing idea
		// field (older events) but reject a mismatching one.
		if idea, _ := e.Data["idea"].(string); idea != "" && idea != d.cfg.IdeaSlug {
			continue
		}
		terminal = e.Type
	}
	return terminal, nil
}

func hasRespondingTo(path string) bool {
	_, ok := readFrontmatterField(path, "responding-to")
	return ok
}

// validateCrossReviewBody enforces the D4 cross-review evidence: the artifact must
// contain a `### @<other>` heading for every other active participant.
func validateCrossReviewBody(path, agent string, participants []string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	body := string(data)
	for _, other := range participants {
		if other == agent {
			continue
		}
		if !strings.Contains(body, "### @"+other) {
			return fmt.Errorf("%s missing cross-review heading ### @%s", path, other)
		}
	}
	return nil
}

// setIdeaStatus rewrites the status: field in the idea 00-prompt.md frontmatter
// (the code path the brief found missing — nothing else writes round-02). Atomic.
func setIdeaStatus(ideaDir, status string) error {
	path := filepath.Join(ideaDir, "00-prompt.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	inFrontmatter := false
	replaced := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			}
			if !replaced {
				lines = append(lines[:i], append([]string{"status: " + status}, lines[i:]...)...)
				replaced = true
			}
			break
		}
		if inFrontmatter && strings.HasPrefix(trimmed, "status:") {
			lines[i] = "status: " + status
			replaced = true
		}
	}
	if !replaced {
		return fmt.Errorf("%s has no frontmatter status field", path)
	}
	return writeFileAtomic(path, []byte(strings.Join(lines, "\n")), 0o644)
}

func writeFileAtomic(path string, data []byte, perm os.FileMode) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, perm); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// NewRunnerAdapter wraps runner.RunRound as a RoundRunner. opts is the same
// runner.Options used for round-01; the adapter sets Round/RoundLabel and forces
// Overwrite=false so already-written artifacts are skipped (idempotent re-entry).
func NewRunnerAdapter(opts runner.Options) RoundRunner {
	return roundRunnerAdapter{base: opts}
}

type roundRunnerAdapter struct {
	base runner.Options
}

func (a roundRunnerAdapter) RunRound(ctx context.Context, round int) error {
	opts := a.base
	opts.Round = round
	opts.RoundLabel = roundLabel(round)
	opts.Overwrite = false
	for _, result := range runner.RunRound(ctx, opts) {
		if result.ExitError != "" {
			return fmt.Errorf("agent %s failed: %s", result.AgentID, result.ExitError)
		}
	}
	return nil
}
