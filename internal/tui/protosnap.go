package tui

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/driver"
	"parley-deck-cli/internal/hitl"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runplan"
	"parley-deck-cli/internal/runstate"
	"parley-deck-cli/internal/store"
)

// ProtocolSnapshotInput carries explicit value copies — the builder runs in an
// async tea.Cmd and must never reach back into the live model or reload run
// state (consensus tui-protocol-visibility D5; never runstate.LoadRunAt).
type ProtocolSnapshotInput struct {
	Root              string
	RunID             string
	RunDir            string
	IdeaSlug          string
	IdeaDir           string
	Participants      []string
	MaxRounds         int
	CrossReviewRounds int // 0 → read cross_review_rounds from 00-prompt.md (default 1)
	Events            []store.Event
	Questions         []hitl.Question
	State             runstate.RunState
	Previous          *ProtocolSnapshot
	Now               time.Time
}

// AgentDelivery is one participant's row in the current step's delivery matrix.
type AgentDelivery struct {
	ID          string
	State       string // pending | running | delivered | failed | killed | skipped
	Unvalidated bool   // events and disk disagree → render a trailing "?"
	At          time.Time
	Note        string
}

// ProtocolSnapshot is the cached, display-ready protocol state for the ribbon,
// the status line, and the Protocol tab. It is produced asynchronously and the
// View only ever renders the cached value.
type ProtocolSnapshot struct {
	Step         int    // 0..8
	StepName     string // kickoff, round-01, cross-review, consensus, final, implement, review, review-consensus, fix-up, complete
	Phase        driver.Phase
	Blocked      bool
	CurrentRound int
	TotalRounds  int // 1 + cross_review_rounds (deliberation budget)
	RoundLabel   string
	IdeaStatus   string
	Implementer  string
	Delivery     []AgentDelivery
	Waiting      []string
	Signoffs     *consensus.Summary
	Next         *runplan.NextAction
	DiskFallback bool // some delivery rows came from disk stats, not events
	ReconciledAt time.Time
	Err          string // last reconcile error (previous snapshot retained)
	// regressSeen remembers a lower step observed once; the regression is only
	// accepted when two consecutive reconciles agree (virtio-fs discipline).
	regressSeen int
}

// narratorTypes is the display allowlist (consensus D7) — what gets woven into
// transcripts. protoTriggerTypes (D6) is the stricter snapshot-trigger set.
var narratorTypes = map[string]bool{
	"run.created": true, "run.phase": true, "run.segment_started": true,
	"agent.started": true, "agent.finished": true, "agent.failed": true,
	"agent.skipped": true, "agent.killed": true,
	"agent.no_first_output": true, "agent.stalled": true,
	"round.completed": true, "round.incomplete": true,
	"hitl.question": true, "hitl.answered": true,
	"run.failed": true, "run.manifest_deferred": true,
	// agent.heartbeat is deliberately ABSENT (consensus D4): heartbeats are
	// status-surface data, never transcript lines.
}

var protoTriggerTypes = map[string]bool{
	"run.created": true, "run.phase": true, "run.segment_started": true,
	"agent.started": true, "agent.finished": true, "agent.failed": true,
	"agent.skipped": true, "agent.killed": true,
	"agent.no_first_output": true, "agent.stalled": true,
	"agent.fixup_finished": true, "agent.fixup_failed": true,
	"round.completed": true, "round.incomplete": true, "run.failed": true,
	// agent.heartbeat never triggers snapshot reconciliation (consensus D4).
}

var stepNames = [9]string{"kickoff", "round-01", "cross-review", "consensus", "final", "implement", "review", "review-consensus", "complete"}

// displayStep maps the 7-state cursor detail to the 9-step pipeline
// (consensus D3). Pure; golden-tested.
func displayStep(d driver.PhaseDetail) (int, string) {
	switch {
	case d.Cursor.Phase == driver.PhaseDone:
		return 8, "complete"
	case strings.HasPrefix(d.ImplementationStatus, "fix-up-cycle"):
		return 8, "fix-up"
	case d.Cursor.Phase == driver.PhaseReview && d.ReviewConsensusExists:
		return 7, stepNames[7]
	case d.Cursor.Phase == driver.PhaseReview:
		return 6, stepNames[6]
	case d.Cursor.Phase == driver.PhaseImpl:
		return 5, stepNames[5]
	case d.Cursor.Phase == driver.PhaseFinal:
		return 4, stepNames[4]
	case d.Cursor.Phase == driver.PhaseConsensus:
		return 3, stepNames[3]
	case d.Cursor.IdeaStatus == "kickoff":
		return 0, stepNames[0]
	case d.Cursor.CurrentRound >= 2:
		return 2, stepNames[2]
	default:
		return 1, stepNames[1]
	}
}

// BuildProtocolSnapshot derives the protocol state from the idea dir and the
// already-loaded run events. Disk reads are bounded: one RebuildDetail pass,
// the 00-prompt frontmatter, at most one consensus.Status parse, and a per-
// participant stat fallback for the active round dir.
func BuildProtocolSnapshot(in ProtocolSnapshotInput) (ProtocolSnapshot, error) {
	detail, derr := driver.RebuildDetail(in.IdeaDir, in.MaxRounds)
	if derr != nil {
		return keepLast(in, fmt.Sprintf("reconcile: %v", derr)), derr
	}
	step, stepName := displayStep(detail)
	out := ProtocolSnapshot{
		Step:         step,
		StepName:     stepName,
		Phase:        detail.Cursor.Phase,
		CurrentRound: detail.Cursor.CurrentRound,
		RoundLabel:   roundLabelFor(detail.Cursor.CurrentRound),
		IdeaStatus:   detail.Cursor.IdeaStatus,
		ReconciledAt: in.Now,
	}
	out.TotalRounds = 1 + in.CrossReviewRounds
	if in.CrossReviewRounds <= 0 {
		out.TotalRounds = 1 + readCrossReviewRoundsFM(in.IdeaDir)
	}
	out.Implementer = implementerOf(in.IdeaDir)

	// Participants precedence (D5): live = opts → run.created payload →
	// 00-prompt frontmatter; display = ordered union. Delivery rows render the
	// display set; the waiting list counts only the live set (a dropped agent's
	// delivered artifact still shows, but nobody waits on it).
	live := append([]string(nil), in.Participants...)
	fromEvents := participantsFromEvents(in.Events)
	fromFM := participantsFromFrontmatter(in.IdeaDir)
	if len(live) == 0 {
		live = fromEvents
	}
	if len(live) == 0 {
		live = fromFM
	}
	display := unionOrdered(live, fromEvents, fromFM)

	// Signoff matrices: design schema at steps 3-4, review schema at 7+.
	switch {
	case step == 3 || step == 4:
		if sum, err := consensus.Status(in.Root, in.IdeaSlug, false); err == nil {
			out.Signoffs = &sum
			out.Blocked = sum.Triage == consensus.TriageBlocked
		}
	case step >= 7:
		if sum, err := consensus.Status(in.Root, in.IdeaSlug, true); err == nil {
			out.Signoffs = &sum
		}
	}

	out.Delivery, out.Waiting, out.DiskFallback = deliveryMatrix(in, detail, step, out.Implementer, out.Signoffs, display, live)

	terminal, outcome := runstate.Outcome(in.Events)
	if actions := runplan.Plan(in.Root, runplan.Input{
		RunID:        in.RunID,
		IdeaSlug:     in.IdeaSlug,
		Participants: append([]string(nil), in.Participants...),
		Terminal:     terminal,
		Outcome:      outcome,
		Questions:    append([]hitl.Question(nil), in.Questions...),
		Agents:       plannerAgents(in.State.Agents),
		RoundStatus:  in.State.RoundStatus,
		CurrentRound: out.RoundLabel,
	}); len(actions) > 0 {
		out.Next = &actions[0]
	}

	// Two-consecutive-agreement before regressing the step (virtio-fs: a stale
	// directory read must not bounce the pipeline backwards).
	if prev := in.Previous; prev != nil && out.Step < prev.Step {
		if prev.regressSeen != out.Step {
			kept := *prev
			kept.regressSeen = out.Step
			kept.ReconciledAt = in.Now
			kept.Err = ""
			return kept, nil
		}
	}
	return out, nil
}

// keepLast returns the previous snapshot annotated with the reconcile error, or
// a minimal error snapshot when there is no previous one.
func keepLast(in ProtocolSnapshotInput, msg string) ProtocolSnapshot {
	if in.Previous != nil {
		kept := *in.Previous
		kept.Err = msg
		kept.ReconciledAt = in.Now
		return kept
	}
	return ProtocolSnapshot{Step: -1, StepName: "unknown", Err: msg, ReconciledAt: in.Now}
}

// deliveryMatrix merges event-projected agent state (primary) with a bounded
// disk fallback (secondary) for the artifacts awaited at the current step.
// display is the wider union roster (rows rendered); live is the narrower
// active roster (waiting math).
func deliveryMatrix(in ProtocolSnapshotInput, detail driver.PhaseDetail, step int, implementer string, signoffs *consensus.Summary, display, live []string) ([]AgentDelivery, []string, bool) {
	switch step {
	case 3, 7:
		// Consensus steps: "delivery" is the signoff matrix; waiting = missing signers.
		if signoffs != nil {
			return nil, append([]string(nil), signoffs.Missing...), false
		}
		return nil, nil, false
	case 4, 8:
		return nil, nil, false
	}

	awaited := append([]string(nil), display...)
	roundDir := filepath.Join(in.IdeaDir, roundLabelFor(detail.Cursor.CurrentRound))
	if step == 5 {
		if implementer != "" {
			awaited = []string{implementer}
		}
		roundDir = "" // IMPLEMENTATION.md presence is already encoded in the step
	}
	if step == 6 {
		awaited = without(awaited, implementer)
		roundDir = filepath.Join(in.IdeaDir, "review", reviewRoundLabel(detail.HighestReviewRound))
	}

	byAgent := map[string]runstate.AgentState{}
	for _, a := range in.State.Agents {
		byAgent[a.ID] = a
	}
	var delivery []AgentDelivery
	var waiting []string
	diskFallback := false
	for _, id := range awaited {
		row := AgentDelivery{ID: id, State: "pending"}
		a, hasEvents := byAgent[id]
		if hasEvents {
			switch a.State {
			case runstate.StateRunning:
				row.State = "running"
			case runstate.StateFinished:
				row.State = "delivered"
				row.At = a.StartedAt.Add(a.Duration)
			case runstate.StateFailed:
				row.State = "failed"
				row.Note = a.Error
			case runstate.StateKilled:
				row.State = "killed"
			case runstate.StateSkipped:
				row.State = "skipped"
				row.Note = a.Reason
			}
		}
		if roundDir != "" {
			onDisk := artifactOnDisk(filepath.Join(roundDir, id+".md"))
			switch {
			case row.State == "delivered" && !onDisk:
				row.Unvalidated = true
				row.Note = "event says delivered, file not visible"
			case row.State != "delivered" && row.State != "running" && onDisk:
				// Disk shows the artifact but events do not (reattach, other run,
				// stale projection) → trust the file, marked unvalidated.
				row.State = "delivered"
				row.Unvalidated = true
				diskFallback = true
			}
		}
		if row.State != "delivered" && contains(live, id) {
			waiting = append(waiting, id)
		}
		delivery = append(delivery, row)
	}
	return delivery, waiting, diskFallback
}

// participantsFromEvents reads the roster from the run.created payload.
func participantsFromEvents(events []store.Event) []string {
	for _, e := range events {
		if e.Type != "run.created" {
			continue
		}
		raw, _ := e.Data["participants"].([]any)
		out := make([]string, 0, len(raw))
		for _, v := range raw {
			if s, _ := v.(string); s != "" {
				out = append(out, s)
			}
		}
		if ps, _ := e.Data["participants"].([]string); len(ps) > 0 {
			out = append(out, ps...)
		}
		return out
	}
	return nil
}

// participantsFromFrontmatter reads the roster from 00-prompt.md ("[a, b, c]").
func participantsFromFrontmatter(ideaDir string) []string {
	meta, err := protocol.ReadFrontmatter(filepath.Join(ideaDir, "00-prompt.md"))
	if err != nil {
		return nil
	}
	raw := strings.TrimSpace(meta["participants"])
	raw = strings.TrimPrefix(raw, "[")
	raw = strings.TrimSuffix(raw, "]")
	var out []string
	for _, part := range strings.Split(raw, ",") {
		if p := strings.TrimSpace(part); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// unionOrdered merges rosters preserving first-seen order.
func unionOrdered(lists ...[]string) []string {
	seen := map[string]bool{}
	var out []string
	for _, list := range lists {
		for _, item := range list {
			if item == "" || seen[item] {
				continue
			}
			seen[item] = true
			out = append(out, item)
		}
	}
	return out
}

func contains(items []string, want string) bool {
	for _, item := range items {
		if item == want {
			return true
		}
	}
	return false
}

// artifactOnDisk is the bounded disk-fallback probe: one Stat with a single
// immediate retry on a transient (non-NotExist) error.
func artifactOnDisk(path string) bool {
	for attempt := 0; attempt < 2; attempt++ {
		info, err := os.Stat(path)
		if err == nil {
			return info.Mode().IsRegular()
		}
		if errors.Is(err, fs.ErrNotExist) {
			return false
		}
	}
	return false
}

// implementerOf reads the implementer's agent id from IMPLEMENTATION.md
// frontmatter; "" when absent.
func implementerOf(ideaDir string) string {
	meta, err := protocol.ReadFrontmatter(filepath.Join(ideaDir, "IMPLEMENTATION.md"))
	if err != nil {
		return ""
	}
	for _, key := range []string{"agent", "implementer", "by"} {
		if v := strings.TrimSpace(meta[key]); v != "" {
			return v
		}
	}
	return ""
}

// readCrossReviewRoundsFM reads cross_review_rounds from 00-prompt.md (default 1).
func readCrossReviewRoundsFM(ideaDir string) int {
	meta, err := protocol.ReadFrontmatter(filepath.Join(ideaDir, "00-prompt.md"))
	if err != nil {
		return 1
	}
	if raw, ok := meta["cross_review_rounds"]; ok {
		if n, err := strconv.Atoi(strings.TrimSpace(raw)); err == nil && n >= 0 {
			return n
		}
	}
	return 1
}

func roundLabelFor(n int) string { return fmt.Sprintf("round-%02d", maxInt(n, 1)) }

func reviewRoundLabel(n int) string { return fmt.Sprintf("round-%02d", maxInt(n, 1)) }

func plannerAgents(agents []runstate.AgentState) []runplan.AgentState {
	out := make([]runplan.AgentState, 0, len(agents))
	for _, a := range agents {
		out = append(out, runplan.AgentState{
			ID: a.ID, State: a.State, ArtifactPath: a.ArtifactPath, Error: a.Error, Reason: a.Reason,
		})
	}
	return out
}

func without(items []string, drop string) []string {
	if drop == "" {
		return items
	}
	out := items[:0]
	for _, it := range items {
		if it != drop {
			out = append(out, it)
		}
	}
	return out
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// hasProtoTrigger reports whether an event batch contains a snapshot-trigger
// type (in-memory check, zero I/O).
func hasProtoTrigger(events []store.Event) bool {
	for _, e := range events {
		if protoTriggerTypes[e.Type] {
			return true
		}
	}
	return false
}

// shortPhaseName compresses a step name for the status line (D10).
func shortPhaseName(step int, stepName string, roundLabel string) string {
	switch step {
	case 1:
		return "r01"
	case 2:
		return "xrev-r" + strings.TrimPrefix(roundLabel, "round-")
	case 7:
		return "rcon"
	default:
		return stepName
	}
}
