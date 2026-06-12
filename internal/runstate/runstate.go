package runstate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"parley-deck-cli/internal/hitl"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runmanifest"
	"parley-deck-cli/internal/runplan"
	"parley-deck-cli/internal/store"
)

const (
	StatePending  = "pending"
	StateRunning  = "running"
	StateFinished = "finished"
	StateFailed   = "failed"
	StateSkipped  = "skipped"
	StateKilled   = "killed"
	StateUnknown  = "unknown"

	OutcomeCompleted  = "completed"
	OutcomeIncomplete = "incomplete"
	OutcomeFailed     = "failed"

	LivenessUnverified = "unverified"
	LivenessIdle       = "idle"

	AttentionAction  = "ACTION"
	AttentionRunning = "RUNNING"
	AttentionFailed  = "FAILED"
	AttentionDone    = "DONE"
	AttentionStale   = "STALE"
	AttentionIdle    = "IDLE"
)

const AttentionStaleAfter = 10 * time.Minute

type RunState struct {
	Agents      []AgentState   `json:"agents"`
	RoundStatus string         `json:"round_status"`
	Recent      []EventSummary `json:"recent"`
}

type AgentState struct {
	ID           string        `json:"id"`
	State        string        `json:"state"`
	Segment      string        `json:"segment,omitempty"`
	StartedAt    time.Time     `json:"started_at,omitempty"`
	Duration     time.Duration `json:"duration"`
	LatestEvent  string        `json:"latest_event,omitempty"`
	Error        string        `json:"error,omitempty"`
	Reason       string        `json:"reason,omitempty"`
	ArtifactPath string        `json:"artifact_path,omitempty"`
	StdoutPath   string        `json:"stdout_path,omitempty"`
	StderrPath   string        `json:"stderr_path,omitempty"`
	// Killed is sticky: set by an agent.killed event so the subsequent (canceled)
	// agent.failed does not downgrade the badge from killed to a generic failure.
	Killed bool `json:"killed,omitempty"`
	// Failure classification + operator hint (runner-hardening-kindly D5/D6).
	FailureClass string `json:"failure_class,omitempty"`
	RecoveryHint string `json:"recovery_hint,omitempty"`
	// AgentExit preserves an artifact-wins nonzero exit on a FINISHED agent.
	AgentExit int `json:"agent_exit,omitempty"`
}

type EventSummary struct {
	Time  time.Time `json:"time"`
	Type  string    `json:"type"`
	Agent string    `json:"agent,omitempty"`
	Text  string    `json:"text,omitempty"`
}

// RunSummary is the intentionally small, unstable developer JSON surface for
// this slice. Keep fields conservative until real consumers exist.
type RunSummary struct {
	RunID         string               `json:"run_id"`
	RunDir        string               `json:"-"`
	IdeaSlug      string               `json:"idea_slug"`
	Task          string               `json:"task,omitempty"`
	Mode          string               `json:"mode,omitempty"`
	CurrentRound  string               `json:"current_round,omitempty"`
	Participants  []string             `json:"participants,omitempty"`
	Terminal      bool                 `json:"terminal"`
	Outcome       string               `json:"outcome,omitempty"`
	Liveness      string               `json:"liveness,omitempty"`
	LastEventAt   time.Time            `json:"last_event_at,omitempty"`
	LastEventAge  time.Duration        `json:"last_event_age"`
	OpenQuestions int                  `json:"open_questions"`
	Attention     string               `json:"attention,omitempty"`
	Questions     []hitl.Question      `json:"questions,omitempty"`
	NextActions   []runplan.NextAction `json:"next_actions,omitempty"`
	State         RunState             `json:"state"`
	Error         string               `json:"error,omitempty"`
}

func RunDir(root, runID string) string {
	return filepath.Join(root, protocol.DeckDir, "runs", runID)
}

func LoadRun(root, runID string) (RunSummary, error) {
	return LoadRunAt(root, runID, time.Now())
}

func LoadRunAt(root, runID string, now time.Time) (RunSummary, error) {
	runDir := RunDir(root, runID)
	manifest, hasManifest := loadManifestSnapshot(root, runID)
	events, err := store.New(runDir).Load()
	if err != nil {
		return RunSummary{RunID: runID, RunDir: runDir}, err
	}

	summary := RunSummary{
		RunID:    runID,
		RunDir:   runDir,
		IdeaSlug: "unknown",
	}
	if hasManifest {
		applyManifestDefaults(&summary, manifest)
	}
	for _, event := range events {
		if !event.Time.IsZero() && (summary.LastEventAt.IsZero() || event.Time.After(summary.LastEventAt)) {
			summary.LastEventAt = event.Time
		}
		if event.Type != "run.created" {
			continue
		}
		if value := dataString(event.Data, "idea"); value != "" {
			summary.IdeaSlug = value
		}
		summary.Task = dataString(event.Data, "task")
		summary.Mode = dataString(event.Data, "mode")
		summary.Participants = dataStringSlice(event.Data, "participants")
	}
	if !summary.LastEventAt.IsZero() {
		summary.LastEventAge = now.Sub(summary.LastEventAt)
		if summary.LastEventAge < 0 {
			summary.LastEventAge = 0
		}
	}

	if len(summary.Participants) == 0 {
		summary.Participants = inferParticipants(root, summary.IdeaSlug)
	}
	summary.State = ProjectEvents(summary.Participants, events, now)
	if summary.CurrentRound == "" {
		summary.CurrentRound = inferCurrentRound(summary.State)
	}
	summary.Terminal, summary.Outcome = deriveOutcome(events)
	if !summary.Terminal {
		summary.Liveness = deriveLiveness(summary.State)
	}

	questions, err := hitl.New(runDir).List()
	if err != nil {
		return summary, err
	}
	summary.Questions = questions
	for _, question := range questions {
		if question.Status == hitl.StatusOpen {
			summary.OpenQuestions++
		}
	}
	summary.Attention = Attention(summary)
	summary.NextActions = runplan.Plan(root, runplan.Input{
		RunID:        summary.RunID,
		IdeaSlug:     summary.IdeaSlug,
		Participants: append([]string(nil), summary.Participants...),
		Terminal:     summary.Terminal,
		Outcome:      summary.Outcome,
		Questions:    append([]hitl.Question(nil), summary.Questions...),
		Agents:       plannerAgents(summary.State.Agents),
		RoundStatus:  summary.State.RoundStatus,
		CurrentRound: summary.CurrentRound,
	})

	return summary, nil
}

func loadManifestSnapshot(root, runID string) (runmanifest.Manifest, bool) {
	manifest, err := runmanifest.Load(root, runID)
	if err != nil {
		return runmanifest.Manifest{}, false
	}
	return manifest, true
}

func applyManifestDefaults(summary *RunSummary, manifest runmanifest.Manifest) {
	if manifest.IdeaSlug != "" && (summary.IdeaSlug == "" || summary.IdeaSlug == "unknown") {
		summary.IdeaSlug = manifest.IdeaSlug
	}
	if summary.Task == "" {
		summary.Task = manifest.Task
	}
	if summary.Mode == "" {
		summary.Mode = manifest.Mode
	}
	if summary.CurrentRound == "" {
		summary.CurrentRound = manifest.CurrentRound
	}
	if len(summary.Participants) == 0 {
		summary.Participants = append([]string(nil), manifest.Participants...)
	}
	if summary.LastEventAt.IsZero() && !manifest.UpdatedAt.IsZero() {
		summary.LastEventAt = manifest.UpdatedAt
	}
}

func inferCurrentRound(state RunState) string {
	best := ""
	for _, agent := range state.Agents {
		round := roundFromArtifact(agent.ArtifactPath)
		if round > best {
			best = round
		}
	}
	if best == "" {
		return "round-01"
	}
	return best
}

func roundFromArtifact(path string) string {
	path = filepath.ToSlash(path)
	parts := strings.Split(path, "/")
	if len(parts) == 0 || !strings.HasPrefix(parts[0], "round-") {
		return ""
	}
	return parts[0]
}

func Attention(run RunSummary) string {
	if run.Error != "" {
		return AttentionFailed
	}
	if run.OpenQuestions > 0 {
		return AttentionAction
	}
	if run.Terminal {
		switch run.Outcome {
		case OutcomeCompleted:
			return AttentionDone
		case OutcomeIncomplete, OutcomeFailed:
			return AttentionFailed
		default:
			return AttentionFailed
		}
	}
	for _, agent := range run.State.Agents {
		if agent.State == StateFailed {
			return AttentionFailed
		}
	}
	for _, agent := range run.State.Agents {
		if agent.State == StateRunning {
			return AttentionRunning
		}
	}
	if !run.LastEventAt.IsZero() && run.LastEventAge >= AttentionStaleAfter {
		return AttentionStale
	}
	return AttentionIdle
}

func ListRuns(root string) ([]RunSummary, error) {
	runsDir := filepath.Join(root, protocol.DeckDir, "runs")
	entries, err := os.ReadDir(runsDir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	now := time.Now()
	runs := make([]RunSummary, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		run, err := LoadRunAt(root, entry.Name(), now)
		if err != nil {
			run = RunSummary{
				RunID:    entry.Name(),
				RunDir:   RunDir(root, entry.Name()),
				IdeaSlug: "unknown",
				Error:    err.Error(),
			}
		}
		runs = append(runs, run)
	}
	sort.SliceStable(runs, func(i, j int) bool {
		return runs[i].RunID > runs[j].RunID
	})
	return runs, nil
}

func ResolveRun(root, target string) (RunSummary, error) {
	target = strings.TrimSpace(target)
	if target == "" {
		return RunSummary{}, fmt.Errorf("run ID or idea slug is required")
	}
	if info, err := os.Stat(RunDir(root, target)); err == nil && info.IsDir() {
		return LoadRun(root, target)
	}

	runs, err := ListRuns(root)
	if err != nil {
		return RunSummary{}, err
	}
	for _, run := range runs {
		if run.IdeaSlug == target && run.Error == "" {
			return run, nil
		}
	}
	if ideaExists(root, target) {
		return RunSummary{}, fmt.Errorf("idea %q has no runs yet", target)
	}
	if len(runs) == 0 {
		return RunSummary{}, fmt.Errorf("no runs found for %q", target)
	}
	return RunSummary{}, fmt.Errorf("no run or idea %q found; available runs: %s", target, strings.Join(runIDs(runs), ", "))
}

func ProjectEvents(participants []string, events []store.Event, now time.Time) RunState {
	agentsByID := map[string]*AgentState{}
	order := append([]string(nil), participants...)
	for _, id := range participants {
		agentsByID[id] = &AgentState{ID: id, State: StatePending}
	}

	state := RunState{RoundStatus: "pending"}
	currentSegment := ""
	for _, event := range events {
		if event.Type == "agent.heartbeat" {
			continue // status surfaces read the live snapshot, not Recent (D4)
		}
		summary := SummarizeEvent(event)
		if summary.Type != "" {
			state.Recent = append(state.Recent, summary)
			if len(state.Recent) > 8 {
				state.Recent = state.Recent[len(state.Recent)-8:]
			}
		}

		switch event.Type {
		case "agent.started", "agent.finished", "agent.failed", "agent.skipped", "agent.killed":
			id := dataString(event.Data, "agent")
			if id == "" {
				continue
			}
			agent, ok := agentsByID[id]
			if !ok {
				agent = &AgentState{ID: id, State: StateUnknown}
				agentsByID[id] = agent
				order = append(order, id)
				agent.LatestEvent = event.Type
				continue
			}
			if agent.State == StateUnknown {
				agent.LatestEvent = event.Type
				continue
			}
			applyAgentEvent(agent, event, now, currentSegment)
		case "run.segment_started":
			// A segment boundary scopes agent state to the current run segment.
			// Resetting the targeted agents here means a stale terminal badge
			// from a prior segment (e.g. an earlier round that finished) never
			// lingers after continue/resume/retry. The next agent.* event in
			// this segment sets the current state; non-targeted agents keep
			// theirs. Old, unsegmented runs emit no such event and so behave
			// exactly as before.
			seg := dataString(event.Data, "segment_id")
			currentSegment = seg
			for _, id := range dataStringSlice(event.Data, "targets") {
				agent, ok := agentsByID[id]
				if !ok {
					continue
				}
				agent.State = StatePending
				agent.StartedAt = time.Time{}
				agent.Duration = 0
				agent.Error = ""
				agent.Reason = ""
				agent.Killed = false
				agent.Segment = seg
				agent.LatestEvent = event.Type
			}
		case "round.completed":
			state.RoundStatus = "completed"
		case "round.incomplete":
			state.RoundStatus = "incomplete"
		}
	}

	sort.SliceStable(order, func(i, j int) bool {
		leftKnown := contains(participants, order[i])
		rightKnown := contains(participants, order[j])
		if leftKnown != rightKnown {
			return leftKnown
		}
		return order[i] < order[j]
	})
	for _, id := range order {
		if agent := agentsByID[id]; agent != nil {
			state.Agents = append(state.Agents, *agent)
		}
	}
	return state
}

func SummarizeEvent(event store.Event) EventSummary {
	agent := dataString(event.Data, "agent")
	text := agent
	switch event.Type {
	case "run.created":
		text = fmt.Sprintf("idea=%s", dataString(event.Data, "idea"))
	case "run.segment_started":
		text = strings.TrimSpace(fmt.Sprintf("%s %s", dataString(event.Data, "segment_id"), dataString(event.Data, "reason")))
	case "steer.requested":
		text = strings.TrimSpace(fmt.Sprintf("%s %s %s", dataString(event.Data, "id"), dataString(event.Data, "target"), dataString(event.Data, "agent")))
	case "steer.delivered":
		text = strings.TrimSpace(fmt.Sprintf("%s %s/%s", dataString(event.Data, "id"), dataString(event.Data, "mode"), dataString(event.Data, "status")))
	case "run.phase":
		text = strings.TrimSpace(fmt.Sprintf("%s %s %s",
			dataString(event.Data, "phase"), dataString(event.Data, "round_label"), dataString(event.Data, "action")))
	case "agent.finished":
		if artifact := relIdeaPath(dataString(event.Data, "artifact")); artifact != "" {
			text = fmt.Sprintf("%s wrote %s", agent, artifact)
		}
	case "agent.failed":
		if class := dataString(event.Data, "failure_class"); class != "" {
			text = fmt.Sprintf("%s failed: %s", agent, class)
			if hint := dataString(event.Data, "recovery_hint"); hint != "" {
				text += " — " + hint
			}
		} else if errText := dataString(event.Data, "error"); errText != "" {
			text = fmt.Sprintf("%s %s", agent, errText)
		}
	case "agent.no_first_output", "agent.stalled":
		text = strings.TrimSpace(fmt.Sprintf("%s %s after %ds (%s)",
			agent, strings.TrimPrefix(event.Type, "agent."), dataInt(event.Data, "elapsed_ms")/1000, dataString(event.Data, "action")))
	case "agent.heartbeat":
		text = strings.TrimSpace(fmt.Sprintf("%s %ds elapsed, %d+%d bytes",
			agent, dataInt(event.Data, "elapsed_ms")/1000, dataInt(event.Data, "stdout_bytes"), dataInt(event.Data, "stderr_bytes")))
	case "agent.skipped":
		text = fmt.Sprintf("%s %s", agent, dataString(event.Data, "reason"))
	case "hitl.question":
		text = fmt.Sprintf("%s question %s %s", agent, dataString(event.Data, "question_id"), dataString(event.Data, "risk"))
	case "hitl.answered":
		text = fmt.Sprintf("%s answered %s %s", agent, dataString(event.Data, "question_id"), dataString(event.Data, "status"))
	case "round.completed", "round.incomplete":
		text = fmt.Sprintf("%s/%s agents", dataString(event.Data, "completed"), dataString(event.Data, "total"))
	}
	return EventSummary{Time: event.Time, Type: event.Type, Agent: agent, Text: strings.TrimSpace(text)}
}

func applyAgentEvent(agent *AgentState, event store.Event, now time.Time, currentSegment string) {
	agent.LatestEvent = event.Type
	// An explicit segment_id wins; otherwise an untagged event inherits the
	// current segment (FINAL backward-compat rule).
	if seg := dataString(event.Data, "segment_id"); seg != "" {
		agent.Segment = seg
	} else if currentSegment != "" {
		agent.Segment = currentSegment
	}
	if artifact := dataString(event.Data, "artifact"); artifact != "" {
		agent.ArtifactPath = artifact
	}
	switch event.Type {
	case "agent.started":
		agent.State = StateRunning
		agent.StartedAt = event.Time
		agent.StdoutPath = dataString(event.Data, "stdout")
		agent.StderrPath = dataString(event.Data, "stderr")
	case "agent.finished":
		agent.State = StateFinished
		agent.Duration = dataDuration(event.Data, "duration_ms", now.Sub(agent.StartedAt))
		agent.AgentExit = dataInt(event.Data, "agent_exit")
		agent.FailureClass = ""
		agent.RecoveryHint = ""
	case "agent.killed":
		// The user terminated this agent; sticky so a later agent.failed keeps it.
		agent.Killed = true
		agent.State = StateKilled
		if agent.Error == "" {
			agent.Error = "killed by user"
		}
	case "agent.failed":
		if agent.Killed {
			agent.State = StateKilled // killed wins over the trailing cancel failure
		} else {
			agent.State = StateFailed
		}
		agent.Error = dataString(event.Data, "error")
		agent.FailureClass = dataString(event.Data, "failure_class")
		agent.RecoveryHint = dataString(event.Data, "recovery_hint")
		agent.Duration = dataDuration(event.Data, "duration_ms", now.Sub(agent.StartedAt))
	case "agent.skipped":
		agent.State = StateSkipped
		agent.Reason = dataString(event.Data, "reason")
		agent.Duration = dataDuration(event.Data, "duration_ms", 0)
	}
}

// relIdeaPath shortens an artifact path to the idea-relative form for display.
func relIdeaPath(path string) string {
	path = filepath.ToSlash(strings.TrimSpace(path))
	if path == "" {
		return ""
	}
	if i := strings.LastIndex(path, "/ideas/"); i >= 0 {
		return path[i+len("/ideas/"):]
	}
	return filepath.Base(path)
}

// Outcome exposes the terminal/outcome projection for callers that already hold
// the event list (the TUI protocol snapshot) and must not reload run state.
func Outcome(events []store.Event) (bool, string) {
	return deriveOutcome(events)
}

func deriveOutcome(events []store.Event) (bool, string) {
	terminal := false
	outcome := ""
	for _, event := range events {
		switch event.Type {
		case "round.completed":
			terminal = true
			outcome = OutcomeCompleted
		case "round.incomplete":
			terminal = true
			outcome = OutcomeIncomplete
		case "run.failed":
			terminal = true
			outcome = OutcomeFailed
		}
	}
	return terminal, outcome
}

func deriveLiveness(state RunState) string {
	for _, agent := range state.Agents {
		if agent.State == StateRunning {
			return LivenessUnverified
		}
	}
	return LivenessIdle
}

func inferParticipants(root, ideaSlug string) []string {
	if ideaSlug == "" || ideaSlug == "unknown" {
		return nil
	}
	status, err := protocol.ReadWorkspaceStatus(root)
	if err != nil {
		return nil
	}
	for _, idea := range status.Ideas {
		if idea.Slug == ideaSlug {
			return append([]string(nil), idea.Participants...)
		}
	}
	return nil
}

func ideaExists(root, ideaSlug string) bool {
	if ideaSlug == "" {
		return false
	}
	status, err := protocol.ReadWorkspaceStatus(root)
	if err != nil {
		return false
	}
	for _, idea := range status.Ideas {
		if idea.Slug == ideaSlug {
			return true
		}
	}
	return false
}

func plannerAgents(agents []AgentState) []runplan.AgentState {
	out := make([]runplan.AgentState, 0, len(agents))
	for _, agent := range agents {
		out = append(out, runplan.AgentState{
			ID:           agent.ID,
			State:        agent.State,
			ArtifactPath: agent.ArtifactPath,
			Error:        agent.Error,
			Reason:       agent.Reason,
		})
	}
	return out
}

func runIDs(runs []RunSummary) []string {
	ids := make([]string, 0, len(runs))
	for _, run := range runs {
		ids = append(ids, run.RunID)
	}
	return ids
}

func dataString(data map[string]any, key string) string {
	if data == nil {
		return ""
	}
	switch value := data[key].(type) {
	case string:
		return value
	case fmt.Stringer:
		return value.String()
	case int:
		return fmt.Sprintf("%d", value)
	case int64:
		return fmt.Sprintf("%d", value)
	case float64:
		return fmt.Sprintf("%.0f", value)
	default:
		return ""
	}
}

func dataInt(data map[string]any, key string) int {
	switch v := data[key].(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	}
	return 0
}

func dataStringSlice(data map[string]any, key string) []string {
	if data == nil {
		return nil
	}
	switch value := data[key].(type) {
	case []string:
		return append([]string(nil), value...)
	case []any:
		values := make([]string, 0, len(value))
		for _, item := range value {
			if text, ok := item.(string); ok && strings.TrimSpace(text) != "" {
				values = append(values, text)
			}
		}
		return values
	case string:
		if strings.TrimSpace(value) == "" {
			return nil
		}
		parts := strings.Split(value, ",")
		values := make([]string, 0, len(parts))
		for _, part := range parts {
			if part = strings.TrimSpace(part); part != "" {
				values = append(values, part)
			}
		}
		return values
	default:
		return nil
	}
}

func dataDuration(data map[string]any, key string, fallback time.Duration) time.Duration {
	if data == nil {
		return fallback
	}
	switch value := data[key].(type) {
	case int:
		return time.Duration(value) * time.Millisecond
	case int64:
		return time.Duration(value) * time.Millisecond
	case float64:
		return time.Duration(value) * time.Millisecond
	default:
		return fallback
	}
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
