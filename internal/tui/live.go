package tui

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"

	"parley-deck-cli/internal/hitl"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runstate"
	"parley-deck-cli/internal/steer"
	"parley-deck-cli/internal/store"
)

const (
	statePending  = runstate.StatePending
	stateRunning  = runstate.StateRunning
	stateFinished = runstate.StateFinished
	stateFailed   = runstate.StateFailed
	stateSkipped  = runstate.StateSkipped
	stateKilled   = runstate.StateKilled
	stateUnknown  = runstate.StateUnknown
)

// Bounded scrollback for per-agent transcript buffers: keep at most 20000 lines
// and at most 4 MiB from the tail, so a very long log never blows up memory.
const (
	maxFocusLines = 20000
	maxFocusBytes = 4 << 20 // 4 MiB
)

// liveMode is the live-run TUI view state. The default tabbed surface is
// modeOverview; modeHelp is the only overlay (opened by /help).
type liveMode int

const (
	modeOverview liveMode = iota
	modeHelp
)

type LiveOptions struct {
	Status       protocol.WorkspaceStatus
	Idea         protocol.IdeaStatus
	Participants []string
	RunID        string
	RunDir       string
	Done         <-chan struct{}
	Cancel       func()
	Resume       bool

	// Unified TUI (unified-tui-home):
	Home  bool       // open to the Home tab (no active run yet)
	Root  string     // workspace root, for refreshing the Home recent-runs list
	Start LaunchFunc // launch a new run (N / new idea); nil disables launching

	// Live steering / kill (tui-live-steering). Injected func seams so the TUI
	// stays decoupled from internal/runner and internal/app. nil = record-only /
	// kill-unavailable (e.g. observational /open runs).
	SubmitSteer SteerFunc     // execute a steer (agent reply round-trip)
	KillAgent   KillAgentFunc // terminate one agent; the rest of the run continues

	// Durable kill / liveness (tui-durable-kill). ReattachKill/ReattachLiveness
	// build per-run seams for an opened run from only its run id (kills via
	// persisted process identity across a parley restart; RUN/STALE badge).
	// Liveness reports the active run's agent process states for the badge.
	ReattachKill     func(runID string) KillAgentFunc
	ReattachLiveness func(runID string) LivenessFunc
	Liveness         LivenessFunc

	// StartAutoDrive kicks the auto-driver for this run (advance past round-01:
	// cross-review → consensus → finalize → opted-in implementation). It is
	// idempotent — calling it again while already driving is a no-op. Wired by
	// `parley run`; nil on observational/opened runs. Bound to the `/run`
	// command so a --no-auto run can be advanced on demand.
	StartAutoDrive func()
}

// LaunchRequest is a request to launch a new idea/run from the TUI.
type LaunchRequest struct{ Task string }

// LaunchResult is what a launch returns: enough to attach to the real runner.
type LaunchResult struct {
	Idea         protocol.IdeaStatus
	Participants []string
	RunID        string
	RunDir       string
	Done         <-chan struct{}
	Cancel       func()
	SubmitSteer  SteerFunc
	KillAgent    KillAgentFunc
	Liveness     LivenessFunc
}

// LaunchFunc launches a run through the existing runner path and returns a
// handle to attach to — the TUI never runs a parallel engine.
type LaunchFunc func(LaunchRequest) (LaunchResult, error)

// SteerRequest asks the live run to deliver/execute a steer. SteerFunc executes
// an agent-targeted steer as a fresh single-agent attempt whose stdout is the
// reply; deck/record-only steers just persist the audit trail.
type SteerRequest struct {
	Target  steer.Target
	AgentID string
	Text    string
}

// SteerResult reports what happened. StdoutPath (when non-empty) is the steer
// attempt's live stdout, which the TUI tails into the agent's tab as the reply.
type SteerResult struct {
	ID         string
	Status     string // "running" | "queued" | "recorded" | "rejected"
	SegmentID  string
	StdoutPath string
}

// SteerFunc executes/records a steer. KillAgentFunc terminates one agent (or, on
// a restarted/observational run, kills it via its persisted process identity, or
// clears a stale "running" badge). It returns a user-facing outcome message.
type SteerFunc func(SteerRequest) (SteerResult, error)
type KillAgentFunc func(agentID string) (string, error)

// LivenessFunc reports a running agent's true process state for the badge:
// "live" (attributed process alive), "stale" (no attributable live process), or
// "" (unknown / not applicable). Injected so internal/tui stays decoupled.
type LivenessFunc func(agentID string) string

// Tab ids: Home is always first; agent tabs are "agent:<id>"; Status is last.
const (
	homeTabID   = "home"
	statusTabID = "status"
)

// agentBuffer is one agent's bounded, follow-capable transcript over its stdout
// log, fed by offset-incremental reads (reuses loadFocusTail/readAppendedLines/
// capFocusLines). Kept per agent so switching tabs is instant and retains scroll.
// transcriptStream tags where a transcript line came from, so it can be styled.
type transcriptStream int

const (
	transcriptStdout transcriptStream = iota
	transcriptStderr
	transcriptSteer
	transcriptEvent
)

// partialMaxBytes bounds a single not-yet-newlined live line (a verbose,
// newline-less progress stream can never grow memory).
const partialMaxBytes = 8192

// transcriptLine is one committed line of an agent's transcript, tagged by stream.
type transcriptLine struct {
	Text   string
	Stream transcriptStream
}

// tailCursor tracks the read position of one log file (stdout/stderr/steer).
type tailCursor struct {
	path   string
	offset int64
	info   os.FileInfo // identity, to detect truncation/rotation
}

// agentBuffer is one agent's live, scrollable transcript: it tails stdout, stderr
// AND the in-flight steer reply into a single CR-aware history, keeping a live
// "partial" (not-yet-newlined) line per stream that rewrites in place on \r.
type agentBuffer struct {
	stdout tailCursor
	stderr tailCursor
	steer  tailCursor

	lines     []transcriptLine
	partial   map[transcriptStream]string
	crPending map[transcriptStream]bool // a \r at a chunk boundary, awaiting \n
	scroll    int
	follow    bool
	trunc     bool
	loaded    bool

	hideStderr   bool // /stderr toggle
	showArtifact bool // /artifact toggle

	// Protocol visibility (tui-protocol-visibility): lastGrowthAt feeds the
	// spinner-vs-silent tab glyph; narratorSeq dedups the woven replay ring.
	lastGrowthAt time.Time
	narratorSeq  int
}

// pickerKind discriminates what a confirmed picker selection dispatches to.
type pickerKind string

const (
	pickerOpen   pickerKind = "open"   // select an idea/run to open
	pickerAnswer pickerKind = "answer" // select an open question to answer
)

// pickerItem is one selectable row: Label is shown, Value is what the command
// runs with. No per-item kind/callback — dispatch is by the active picker.Kind.
type pickerItem struct{ Label, Value string }

// pickerState is the reusable arrow-key picker sub-mode (mirrors composing). It
// temporarily owns ↑/↓/Enter/esc/printable/backspace while Active; everything
// else is a no-op. Index/Offset index into the FILTERED view.
type pickerState struct {
	Active bool
	Kind   pickerKind
	Title  string
	Items  []pickerItem
	Index  int
	Filter string // case-insensitive substring over Label+Value
	Offset int    // top of the visible scroll window
}

// filtered returns the items matching Filter (case-insensitive substring over
// Label+Value); empty filter returns all items.
func (p pickerState) filtered() []pickerItem {
	needle := strings.ToLower(strings.TrimSpace(p.Filter))
	if needle == "" {
		return p.Items
	}
	var out []pickerItem
	for _, it := range p.Items {
		if strings.Contains(strings.ToLower(it.Label+" "+it.Value), needle) {
			out = append(out, it)
		}
	}
	return out
}

// reclamp keeps Index inside the filtered list and Offset a window of `visible`
// rows that always contains Index.
func (p *pickerState) reclamp(visible int) {
	n := len(p.filtered())
	if n == 0 {
		p.Index, p.Offset = 0, 0
		return
	}
	if p.Index < 0 {
		p.Index = 0
	}
	if p.Index >= n {
		p.Index = n - 1
	}
	if visible < 1 {
		visible = 1
	}
	if p.Index < p.Offset {
		p.Offset = p.Index
	}
	if p.Index >= p.Offset+visible {
		p.Offset = p.Index - visible + 1
	}
	maxOff := n - visible
	if maxOff < 0 {
		maxOff = 0
	}
	if p.Offset > maxOff {
		p.Offset = maxOff
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
}

type liveModel struct {
	opts      LiveOptions
	width     int
	height    int
	now       time.Time
	offset    int64
	events    []store.Event
	state     RunState
	selected  int
	questions []hitl.Question
	selectedQ int
	mode      liveMode // modeOverview (tabbed surface) or modeHelp (overlay)
	errText   string
	statusMsg string
	done      bool

	// Claude-CLI tabbed layout (tui-claude-cli-layout): the default surface.
	activeTab string                  // "" = resolve to default; else home/agent:<id>/status
	inputText string                  // persistent bottom prompt
	inputErr  string                  // transient input/route error
	buffers   map[string]*agentBuffer // per-agent transcript, lazily loaded for active+visited

	// Unified TUI (unified-tui-home):
	composing bool                  // new-idea input mode (opened by N)
	runToken  int                   // identifies the active run; stale ticks/reads are ignored
	homeRuns  []runstate.RunSummary // recent runs for the Home tab

	// Command picker (tui-command-picker): a bare arg-taking slash command
	// (/open, /answer) opens an arrow-key picker instead of erroring.
	picker    pickerState
	answerQID string // non-empty while composing an answer chosen via the /answer picker

	// Live steering / kill / autocomplete (tui-live-steering).
	suggest          bool          // slash-command suggestion menu visible
	suggestItems     []commandSpec // current matches for inputText prefix
	suggestIndex     int           // highlighted suggestion
	confirmKillAgent string        // non-empty = modal "kill <agent>? (y/N)" confirm

	// Protocol visibility (tui-protocol-visibility): the cached snapshot plus
	// the async-build gates, the woven-narrator ring, and the stat growth cache.
	proto         *ProtocolSnapshot
	protoSeq      int
	protoBusy     bool
	protoDirty    bool
	ribbonMode    int // ribbonCollapsed | ribbonExpanded | ribbonHidden (Ctrl+P)
	narrateMode   int // narrateProtocol | narrateVerbose | narrateOff (/narrate)
	narratorRing  []narratorEntry
	narratorTotal int
	growth        map[string]growthInfo // unvisited-tab activity (2s stat cache)
	buffersStdout map[string]bool       // declared buffers_stdout per agent (run.created)
	ideaPhases    map[string]string     // Home phase chips, computed in refreshHomeRuns
}

// hasRun reports whether an active run is attached (vs the no-run Home state).
func (m liveModel) hasRun() bool { return strings.TrimSpace(m.opts.RunDir) != "" }

// attached reports whether this run was launched by the TUI (cancelable) vs
// opened observationally (/open, resume).
func (m liveModel) attached() bool { return m.opts.Cancel != nil }

type RunState = runstate.RunState
type AgentState = runstate.AgentState
type EventSummary = runstate.EventSummary

type eventsMsg struct {
	events []store.Event
	offset int64
	err    error
	token  int
}

type questionsMsg struct {
	questions []hitl.Question
	err       error
	token     int
}

type eventTickMsg struct {
	t     time.Time
	token int
}
type elapsedTickMsg struct {
	t     time.Time
	token int
}
type doneMsg struct{ token int }

var ansiSequence = regexp.MustCompile(`\x1b\[[0-?]*[ -/]*[@-~]`)

func RunLive(opts LiveOptions) error {
	model := newLiveModel(opts)
	program := tea.NewProgram(model, tea.WithAltScreen())
	_, err := program.Run()
	return err
}

func newLiveModel(opts LiveOptions) liveModel {
	now := time.Now()
	model := liveModel{
		opts:  opts,
		width: 100,
		now:   now,
	}
	model.buffers = map[string]*agentBuffer{}
	model.state = ProjectEvents(opts.Participants, nil, now)
	model.refreshHomeRuns()
	model.ensureActiveBuffer()
	return model
}

func (m liveModel) Init() tea.Cmd {
	// The elapsed clock is global (one loop forever); run-specific reads only
	// start when a run is attached.
	cmds := []tea.Cmd{elapsedTickCmd(m.runToken)}
	cmds = append(cmds, m.runCmds()...)
	return tea.Batch(cmds...)
}

// runCmds are the run-specific commands for the currently attached run (empty
// at Home). They carry m.runToken so stale ones are ignored after a run swap.
func (m liveModel) runCmds() []tea.Cmd {
	if !m.hasRun() {
		return nil
	}
	cmds := []tea.Cmd{
		readEventsCmd(filepath.Join(m.opts.RunDir, "events.jsonl"), m.offset, m.runToken),
		readQuestionsCmd(m.opts.RunDir, m.runToken),
		eventTickCmd(m.runToken),
		// First protocol snapshot lands ~1s after attach; the handler re-arms
		// at the 15s/60s reconcile cadence. Growth probes run on their own 2s tick.
		protoTickCmd(m.runToken, time.Second),
		growthTickCmd(m.runToken),
	}
	if m.opts.Done != nil {
		cmds = append(cmds, waitDoneCmd(m.opts.Done, m.runToken))
	}
	return cmds
}

func (m liveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.picker.Active {
			m.picker.reclamp(m.pickerRows()) // visible-row count may have changed
		}
		return m, nil
	case tea.KeyMsg:
		if m.mode == modeHelp {
			return m.updateHelpMode(msg)
		}
		return m.updateMain(msg)
	case eventsMsg:
		if msg.token != m.runToken {
			return m, nil // stale read from a previous run
		}
		m.offset = msg.offset
		if msg.err != nil {
			m.errText = msg.err.Error()
		} else {
			m.errText = ""
			m.events = append(m.events, msg.events...)
			m.state = ProjectEvents(m.opts.Participants, m.events, m.now)
			m.ensureActiveBuffer()
			m.refreshBuffers()              // drain steer reply stdout into the transcript FIRST
			m.appendSteerEvents(msg.events) // then commit the trailing line + the marker
			m.noteRuntimeFlags(msg.events)
			m.appendProtocolEvents(msg.events)
			if hasProtoTrigger(msg.events) {
				return m, m.scheduleProtoRefresh()
			}
		}
		return m, nil // never quit on done — the TUI stays open
	case eventTickMsg:
		if msg.token != m.runToken {
			return m, nil // stale tick loop dies (no re-arm)
		}
		return m, tea.Batch(
			readEventsCmd(filepath.Join(m.opts.RunDir, "events.jsonl"), m.offset, m.runToken),
			readQuestionsCmd(m.opts.RunDir, m.runToken),
			eventTickCmd(m.runToken),
		)
	case questionsMsg:
		if msg.token != m.runToken {
			return m, nil
		}
		if msg.err != nil {
			m.errText = msg.err.Error()
		} else {
			m.questions = msg.questions
			m.refreshPickerItems() // keep an open /answer picker in sync
		}
		return m, nil
	case elapsedTickMsg:
		// Global clock — not gated by runToken; one loop forever.
		m.now = msg.t
		if m.hasRun() {
			m.state = ProjectEvents(m.opts.Participants, m.events, m.now)
			m.ensureActiveBuffer()
			m.refreshBuffers()
		}
		return m, elapsedTickCmd(m.runToken)
	case doneMsg:
		if msg.token != m.runToken {
			return m, nil
		}
		m.done = true
		m.refreshHomeRuns()
		// one final read + a final protocol reconcile; do NOT quit
		cmds := []tea.Cmd{readEventsCmd(filepath.Join(m.opts.RunDir, "events.jsonl"), m.offset, m.runToken)}
		if c := m.scheduleProtoRefresh(); c != nil {
			cmds = append(cmds, c)
		}
		return m, tea.Batch(cmds...)
	case protoMsg:
		if msg.token != m.runToken || msg.seq != m.protoSeq {
			return m, nil // stale snapshot (run swap or superseded build)
		}
		m.protoBusy = false
		snap := msg.snap
		m.proto = &snap
		if m.protoDirty {
			m.protoDirty = false
			return m, m.scheduleProtoRefresh()
		}
		return m, nil
	case protoTickMsg:
		if msg.token != m.runToken {
			return m, nil
		}
		cmds := []tea.Cmd{protoTickCmd(m.runToken, m.protoInterval())}
		if c := m.scheduleProtoRefresh(); c != nil {
			cmds = append(cmds, c)
		}
		return m, tea.Batch(cmds...)
	case growthTickMsg:
		if msg.token != m.runToken {
			return m, nil
		}
		cmds := []tea.Cmd{growthTickCmd(m.runToken)}
		if paths := m.growthPaths(); len(paths) > 0 {
			cmds = append(cmds, statGrowthCmd(paths, m.runToken))
		}
		return m, tea.Batch(cmds...)
	case growthMsg:
		if msg.token != m.runToken {
			return m, nil
		}
		if m.growth == nil {
			m.growth = map[string]growthInfo{}
		}
		for id, sizes := range msg.sizes {
			prev := m.growth[id]
			next := growthInfo{stdout: sizes[0], stderr: sizes[1], lastGrowth: prev.lastGrowth}
			if sizes[0] > prev.stdout || sizes[1] > prev.stderr {
				next.lastGrowth = msg.at
			}
			m.growth[id] = next
		}
		return m, nil
	case editorFinishedMsg:
		switch {
		case msg.err != nil && !msg.canceled:
			m.inputErr = "editor: " + msg.err.Error()
		case msg.canceled:
			m.statusMsg = "editor cancelled — composer unchanged"
		default:
			// Success: drop the edited text into the composer; the existing Enter
			// path still decides steer vs answer vs launch. Nothing is sent here.
			m.inputText = msg.content
			m.inputErr = ""
			m.suggest, m.suggestItems, m.suggestIndex = false, nil, 0
		}
		return m, nil
	}
	return m, nil
}

// refreshHomeRuns reloads the Home recent-runs list (best-effort).
func (m *liveModel) refreshHomeRuns() {
	if strings.TrimSpace(m.opts.Root) == "" {
		return
	}
	if runs, err := runstate.ListRuns(m.opts.Root); err == nil {
		if len(runs) > 8 {
			runs = runs[:8]
		}
		m.homeRuns = runs
		m.refreshPickerItems() // keep an open /open picker in sync
	}
	// Home phase chips: per-idea pipeline position, computed only here (on
	// launch/open/done) — never on a tick (consensus D14).
	chips := make(map[string]string, len(m.opts.Status.Ideas))
	for _, idea := range m.opts.Status.Ideas {
		if chip := ideaPhaseChip(idea.Path); chip != "" {
			chips[idea.Slug] = chip
		}
	}
	m.ideaPhases = chips
}

// activateRun swaps the model onto a newly launched (or opened) run in place:
// it resets run state, bumps runToken (so stale ticks/reads from the previous
// run are dropped), and returns the new run's command batch.
func (m *liveModel) activateRun(s LaunchResult) tea.Cmd {
	m.runToken++
	m.opts.Idea = s.Idea
	m.opts.Participants = s.Participants
	m.opts.RunID = s.RunID
	m.opts.RunDir = s.RunDir
	m.opts.Done = s.Done
	m.opts.Cancel = s.Cancel
	// Copy the live-steering/kill seams onto the active run — without this, runs
	// launched from Home (newLaunchFunc) would silently lose steering and kill.
	m.opts.SubmitSteer = s.SubmitSteer
	m.opts.KillAgent = s.KillAgent
	m.opts.Liveness = s.Liveness
	m.opts.Resume = false
	m.opts.Home = false
	m.offset = 0
	m.events = nil
	m.questions = nil
	m.done = false
	m.errText = ""
	m.composing = false
	m.answerQID = ""
	m.picker = pickerState{}
	m.suggest = false
	m.suggestItems = nil
	m.confirmKillAgent = ""
	m.inputText = ""
	m.inputErr = ""
	m.buffers = map[string]*agentBuffer{}
	m.activeTab = "" // resolve to the first agent
	m.state = ProjectEvents(s.Participants, nil, m.now)
	// Reset the protocol-visibility caches for the new run.
	m.proto = nil
	m.protoSeq = 0
	m.protoBusy = false
	m.protoDirty = false
	m.narratorRing = nil
	m.narratorTotal = 0
	m.growth = map[string]growthInfo{}
	m.buffersStdout = map[string]bool{}
	m.ensureActiveBuffer()
	return tea.Batch(m.runCmds()...)
}

func (m liveModel) View() string {
	width := tuiWidth(m.width)
	height := tuiHeight(m.height, defaultLiveHeight)
	if m.mode == modeHelp {
		return m.renderHelp(width, height)
	}
	return m.renderTabbed(width, height)
}

// renderTabbed is the default Claude-CLI-style surface: a top tab strip, the
// active agent's live transcript (or the Status dashboard) as the main area,
// then a status line and a persistent input row.
func (m liveModel) renderTabbed(width, height int) string {
	rows := m.transcriptHeight()
	// The picker / suggestion menu overlay above the input row; shrink the
	// transcript area by their height so a tiny terminal never overflows.
	var overlay string
	if m.picker.Active {
		overlay = m.renderPicker(width)
	} else if m.suggest {
		overlay = m.renderSuggest(width)
	}
	if overlay != "" {
		rows -= strings.Count(overlay, "\n") + 1
		if rows < 1 {
			rows = 1
		}
	}
	active := m.activeTabResolved()
	var main string
	switch {
	case active == homeTabID:
		main = m.renderHome(width, rows)
	case strings.HasPrefix(active, "agent:"):
		agentID, _ := agentTab(active)
		main = m.renderTranscript(agentID, width, rows)
	default:
		main = m.renderStatusTab(width, rows)
	}
	parts := []string{m.renderTabStrip(width)}
	if m.ribbonHeight() > 0 {
		parts = append(parts, m.renderRibbon(width))
	}
	parts = append(parts, "", clipLines(main, rows))
	if banner := m.renderQuestionBanner(width); banner != "" {
		parts = append(parts, banner)
	}
	if overlay != "" {
		parts = append(parts, overlay)
	}
	parts = append(parts, m.renderStatusLine(width), m.renderInputRow(width))
	return strings.Join(parts, "\n")
}

// renderSuggest draws the slim slash-command autocomplete menu above the input row.
func (m liveModel) renderSuggest(width int) string {
	var b strings.Builder
	b.WriteString(headerStyle.Render("commands"))
	b.WriteString("\n")
	for i, s := range m.suggestItems {
		label := s.Name
		if s.Usage != "" {
			label = s.Usage
		}
		if i == m.suggestIndex {
			b.WriteString(okStyle.Render("> " + truncateText(label, width-6)))
		} else {
			b.WriteString("  " + mutedStyle.Render(truncateText(label, width-6)))
		}
		b.WriteString("\n")
	}
	b.WriteString(mutedStyle.Render("Tab complete · ↑/↓ pick · Enter run"))
	return boxStyle.Width(clampInt(width-4, 40, 84)).Render(strings.TrimRight(b.String(), "\n"))
}

func (m liveModel) renderTabStrip(width int) string {
	type tab struct{ id, label string }
	var tabs []tab
	for _, id := range m.tabIDs() {
		switch {
		case id == homeTabID:
			tabs = append(tabs, tab{id, "Home"})
		case id == statusTabID:
			tabs = append(tabs, tab{id, "Status"})
		default:
			agentID, _ := agentTab(id)
			// Conservative single-cell activity glyph (consensus D8): spinner =
			// output flowing, · = running silent, ! = STALE, ✓/✗/x/-/○ terminal.
			tabs = append(tabs, tab{id, agentID + " " + m.agentGlyph(agentID)})
		}
	}

	active := m.activeTabResolved()
	ai := 0
	for i, t := range tabs {
		if t.id == active {
			ai = i
			break
		}
	}
	rendered := make([]string, len(tabs))
	w := make([]int, len(tabs))
	for i, t := range tabs {
		rendered[i] = m.styleTab(t.label, i == ai)
		w[i] = len(stripANSI(rendered[i])) + 1 // +1 for the join space (conservative)
	}
	// Window that always includes the active tab, expanding to neighbors while
	// it fits, with "…+N" markers for any clipped tabs (reserve ~8 cols).
	budget := width - 8
	if budget < 1 {
		budget = width
	}
	lo, hi, total := ai, ai, w[ai]
	for {
		grew := false
		if hi+1 < len(tabs) && total+w[hi+1] <= budget {
			hi++
			total += w[hi]
			grew = true
		}
		if lo-1 >= 0 && total+w[lo-1] <= budget {
			lo--
			total += w[lo]
			grew = true
		}
		if !grew {
			break
		}
	}
	strip := strings.Join(rendered[lo:hi+1], " ")
	if lo > 0 {
		strip = mutedStyle.Render(fmt.Sprintf("…+%d ", lo)) + strip
	}
	if hi < len(tabs)-1 {
		strip = strip + mutedStyle.Render(fmt.Sprintf(" …+%d", len(tabs)-1-hi))
	}
	return strip
}

func (m liveModel) styleTab(label string, on bool) string {
	// While the picker is active, dim every tab so focus reads as "on the picker"
	// and the disabled tab navigation is visually obvious.
	if on && !m.picker.Active {
		return headerStyle.Render("▸ " + label)
	}
	return mutedStyle.Render("  " + label)
}

func (m liveModel) renderTranscript(agentID string, width, rows int) string {
	header := m.renderAgentStatusHeader(agentID, width)
	bodyRows := rows - 1 // header takes a row
	if bodyRows < 1 {
		bodyRows = 1
	}
	b := m.ensureBuffer(agentID)

	// /artifact mode: show the produced artifact instead of the live logs.
	if b.showArtifact {
		return header + "\n" + m.renderArtifactView(agentID, width, bodyRows)
	}

	lines := b.visibleLines()
	var sb strings.Builder
	sb.WriteString(header)
	sb.WriteString("\n")
	if len(lines) == 0 {
		// A running-but-silent agent gets the structured placeholder (declared
		// buffering, liveness, byte counters, own recent activity) instead of a
		// dead-looking blank pane.
		if a := m.agentByID(agentID); a != nil && a.State == stateRunning {
			sb.WriteString(m.renderSilentPlaceholder(agentID, width))
			return strings.TrimRight(sb.String(), "\n")
		}
		sb.WriteString(mutedStyle.Render("  (no output yet — the agent is starting; stderr/replies appear here live)"))
		return strings.TrimRight(sb.String(), "\n")
	}
	start := b.scroll
	if start < 0 {
		start = 0
	}
	if bottom := m.bufferBottom(b); start > bottom {
		start = bottom
	}
	end := start + bodyRows
	if end > len(lines) {
		end = len(lines)
	}
	if start == 0 && b.trunc {
		sb.WriteString(mutedStyle.Render("… earlier output truncated"))
		sb.WriteString("\n")
	}
	for i := start; i < end; i++ {
		sb.WriteString(styleTranscriptLine(lines[i], width))
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

// styleTranscriptLine renders one transcript line styled by its stream: stdout
// plain, stderr dimmed + [err], steer cyan with a ❯ user prefix, events dimmed.
func styleTranscriptLine(l transcriptLine, width int) string {
	switch l.Stream {
	case transcriptStderr:
		return mutedStyle.Render(truncateText("[err] "+l.Text, width-1))
	case transcriptSteer:
		return steerStyle.Render(truncateText(l.Text, width-1))
	case transcriptEvent:
		return okStyle.Render(truncateText(l.Text, width-1))
	default:
		return truncateText(l.Text, width-1)
	}
}

// renderAgentStatusHeader is the always-on one-line status so the tab is never
// blank, derived from the projection (state, elapsed/duration, artifact, error).
func (m liveModel) renderAgentStatusHeader(agentID string, width int) string {
	a := m.agentByID(agentID)
	if a == nil {
		return mutedStyle.Render(agentID)
	}
	dur := formatAgentDuration(*a, m.now)
	var line string
	switch a.State {
	case stateRunning:
		tag := "● " + agentID + " working… " + dur
		if lv := m.agentLiveness(agentID); lv == "stale" {
			tag += " · STALE"
		} else if lv != "" {
			tag += " · proc:" + lv
		}
		if b := m.buffers[agentID]; b != nil && b.loaded {
			tag += fmt.Sprintf(" · stdout %s · stderr %s", formatByteCount(b.stdout.offset), formatByteCount(b.stderr.offset))
		}
		if art := relArtifact(a.ArtifactPath); art != "" {
			tag += " · → " + art
		}
		line = warnStyle.Render(truncateText(tag, width-1))
	case stateFinished:
		art := relArtifact(a.ArtifactPath)
		msg := "✓ " + agentID + " finished " + dur
		if a.AgentExit != 0 {
			// Artifact-wins completion: succeeded despite a nonzero exit; keep
			// the rendering calm — no error styling (consensus D7/agy UX).
			msg += fmt.Sprintf(" (exit %d, artifact verified)", a.AgentExit)
		}
		if art != "" {
			msg += " · wrote " + art
		}
		line = okStyle.Render(truncateText(msg, width-1))
	case stateFailed:
		msg := "✗ " + agentID + " failed " + dur
		if a.FailureClass != "" {
			msg += " · class: " + a.FailureClass
			if a.RecoveryHint != "" {
				msg += " — " + a.RecoveryHint
			}
		} else if a.Error != "" {
			msg += ": " + a.Error
		}
		line = warnStyle.Render(truncateText(msg, width-1))
	case stateKilled:
		line = mutedStyle.Render(truncateText("◌ "+agentID+" killed", width-1))
	case stateSkipped:
		line = mutedStyle.Render(truncateText("– "+agentID+" skipped: "+a.Reason, width-1))
	default:
		line = mutedStyle.Render(truncateText("○ "+agentID+" pending", width-1))
	}
	return line
}

// renderArtifactView shows a bounded tail of the agent's produced artifact.
func (m liveModel) renderArtifactView(agentID string, width, rows int) string {
	a := m.agentByID(agentID)
	if a == nil || strings.TrimSpace(a.ArtifactPath) == "" {
		return mutedStyle.Render("[Artifact not yet written]")
	}
	lines, _, _ := loadFocusTail(a.ArtifactPath)
	var sb strings.Builder
	sb.WriteString(okStyle.Render(truncateText("[Viewing Artifact: "+relArtifact(a.ArtifactPath)+"]  (/artifact to return)", width-1)))
	sb.WriteString("\n")
	if len(lines) == 0 {
		sb.WriteString(mutedStyle.Render("[Artifact not yet written]"))
		return sb.String()
	}
	if len(lines) > rows {
		lines = lines[len(lines)-rows:] // show the newest rows of the bounded tail
	}
	for _, l := range lines {
		sb.WriteString(truncateText(l, width-1))
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

// relArtifact shows the artifact path from the idea slug down (readable).
func relArtifact(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	if i := strings.LastIndex(path, "/ideas/"); i >= 0 {
		return path[i+len("/ideas/"):]
	}
	if strings.HasPrefix(path, "ideas/") {
		return path[len("ideas/"):]
	}
	return filepath.Base(path)
}

// renderStatusTab is the Protocol tab: the pipeline/delivery/signoff/next panes
// above the original dashboard panes (/status and /protocol both land here).
func (m liveModel) renderStatusTab(width, rows int) string {
	body := strings.Join([]string{
		m.renderProtocolPanes(),
		"",
		m.renderAgentTable(),
		"",
		m.renderEventPane(),
		"",
		m.renderQuestionsPane(),
	}, "\n")
	return body
}

func (m liveModel) renderQuestionBanner(width int) string {
	agentID, ok := agentTab(m.activeTabResolved())
	if !ok {
		return ""
	}
	q := m.openQuestionFor(agentID)
	if q == nil {
		return ""
	}
	return warnStyle.Render(truncateText(fmt.Sprintf("? %s — type the answer below + Enter (or /answer %s …)", q.Prompt, q.ID), width-1))
}

func (m liveModel) renderStatusLine(width int) string {
	active := m.activeTabResolved()
	label, stateStr := "Status", ""
	if agentID, ok := agentTab(active); ok {
		label = agentID
		if a := m.agentByID(agentID); a != nil {
			stateStr = a.State + " " + formatAgentDuration(*a, m.now)
		}
	}
	follow := ""
	if b := m.activeBuffer(); b != nil {
		if b.follow {
			follow = "follow:on"
		} else {
			follow = "follow:off"
		}
	}
	openQ := 0
	for _, q := range m.questions {
		if q.Status == hitl.StatusOpen {
			openQ++
		}
	}
	doneTag := ""
	if m.done {
		doneTag = "  [done]"
	}
	line := fmt.Sprintf("run=%s idea=%s %s  %s %s %s  q:%d%s",
		m.opts.RunID, m.opts.Idea.Slug,
		m.statusPhaseSegment(),
		label, stateStr, follow, openQ, doneTag)
	out := mutedStyle.Render(truncateText(line, width))
	if m.errText != "" {
		out = warnStyle.Render(truncateText(m.errText, width)) + "\n" + out
	} else if m.statusMsg != "" {
		out = okStyle.Render(truncateText(m.statusMsg, width)) + "\n" + out
	}
	return out
}

func (m liveModel) renderInputRow(width int) string {
	// Modal kill confirmation takes over the input row.
	if m.confirmKillAgent != "" {
		prompt := "kill agent " + m.confirmKillAgent + " and its process tree? (y/N)"
		hint := "y confirm · n/esc cancel (the rest of the run keeps going)"
		if m.agentLiveness(m.confirmKillAgent) == "stale" {
			prompt = "clear stale running status for " + m.confirmKillAgent + "? (y/N)"
			hint = "y clears the stale badge (no live process found) · n/esc cancel"
		}
		row := warnStyle.Render(truncateText(prompt, width))
		return row + "\n" + mutedStyle.Render(hint)
	}
	active := m.activeTabResolved()
	answer := false
	steerRow := false
	steerAgent := ""
	var label string
	switch {
	case m.composing && m.answerQID != "":
		// Step 2 of the /answer picker: colour-flip so an answer is never accidental.
		label = "answer " + m.answerQID + " › "
		answer = true
	case m.composing:
		label = "new idea › "
	case active == homeTabID:
		label = "› "
	default:
		if agentID, ok := agentTab(active); ok {
			if q := m.openQuestionFor(agentID); q != nil {
				label = fmt.Sprintf("answer %s/%s › ", agentID, q.ID)
				answer = true
			} else {
				label = "steer " + agentID + " › "
				steerRow = true
				steerAgent = agentID
			}
		} else {
			label = "deck › "
		}
	}
	// Flatten any multi-line draft (from /editor) to a single-line preview so it
	// never breaks the input row; m.inputText keeps the raw value for submission.
	shown := editorPreview(m.inputText)
	var row string
	switch {
	case answer:
		// Colour-flip the whole row in answer mode so an answer is never accidental.
		row = warnStyle.Render(truncateText(label+shown, width))
	case steerRow:
		// Cyan prefix marks "this goes to the agent" so a steer is never accidental.
		row = steerStyle.Render(label) + truncateText(shown, width-len(label))
	default:
		row = okStyle.Render(label) + truncateText(shown, width-len(label))
	}
	if m.inputErr != "" {
		row += "  " + warnStyle.Render(m.inputErr)
	}
	hint := "↑/↓ tabs · N new idea · Enter steer/answer · ctrl+k kill · /help · ctrl+c"
	if m.done {
		hint = "[done] · ↑/↓ tabs · N new idea · /open (pick) · /quit or esc to exit"
	}
	switch {
	case m.picker.Active:
		hint = "↑/↓ select · type filter · Enter choose · esc cancel"
	case m.suggest:
		hint = "Tab complete · ↑/↓ pick · Enter run · esc close"
	case m.composing && m.answerQID != "":
		hint = "type the answer · Enter submit · esc cancel"
	case m.composing:
		hint = "type a task · Enter launch · esc cancel"
	case strings.HasPrefix(m.inputText, "/"):
		hint = "Tab to complete · commands: /help /protocol /refresh /narrate /follow /deck /editor /answer /open /home /quit"
	case steerRow:
		hint = "Enter sends to " + steerAgent + " · ctrl+k kill · ↑/↓ tabs · /help"
	case active == homeTabID && !m.done:
		hint = "N new idea · /open (pick) · ↑/↓ tabs · /help · ctrl+c"
	}
	return row + "\n" + mutedStyle.Render(hint)
}

// renderHome lists open ideas and recent runs; N starts a new idea, /open reopens.
func (m liveModel) renderHome(width, rows int) string {
	var b strings.Builder
	if n := m.staleAgentCount(); n > 0 {
		b.WriteString(warnStyle.Render(truncateText(fmt.Sprintf("⚠ %d stale agent process(es) — open the agent's tab and ctrl+k to clear", n), width-1)))
		b.WriteString("\n\n")
	}
	b.WriteString(sectionTitle("Ideas"))
	b.WriteString("\n")
	if len(m.opts.Status.Ideas) == 0 {
		b.WriteString(mutedStyle.Render("  no ideas yet — press N to start one"))
	} else {
		for _, idea := range m.opts.Status.Ideas {
			chip := m.ideaPhases[idea.Slug]
			if chip == "" {
				chip = "—" // never duplicate the raw status into the chip column
			}
			row := fmt.Sprintf("  %-34s %-22s %-14s", truncateText(idea.Slug, 34), truncateText(chip, 22), idea.Status)
			// Attention badge from the idea's latest run (D14); homeRuns is
			// already in memory — no I/O on the render path.
			for _, r := range m.homeRuns {
				if r.IdeaSlug == idea.Slug {
					if r.Attention != "" {
						row += " " + attentionBadge(r.Attention)
					}
					break
				}
			}
			b.WriteString(strings.TrimRight(row, " "))
			b.WriteString("\n")
		}
	}
	b.WriteString("\n")
	b.WriteString(sectionTitle("Recent runs"))
	b.WriteString("\n")
	if len(m.homeRuns) == 0 {
		b.WriteString(mutedStyle.Render("  no runs yet"))
	} else {
		for _, r := range m.homeRuns {
			st := "active"
			if r.Terminal {
				st = strings.ToUpper(valueOr(r.Outcome, "done"))
			}
			b.WriteString(fmt.Sprintf("  %-26s %-10s %s\n", truncateText(r.IdeaSlug, 26), st, r.RunID))
		}
	}
	b.WriteString("\n")
	b.WriteString(mutedStyle.Render("Press N to start a new idea (all available agents) · /open <slug|run> to reopen"))
	return strings.TrimRight(b.String(), "\n")
}

// --- slash-command autocomplete + per-agent kill (tui-live-steering) ---

// commandSpec describes one slash command — the single source for runCommand
// dispatch hints, the autocomplete menu, and /help.
type commandSpec struct {
	Name        string
	Usage       string
	OpensPicker bool // bare form opens the existing picker (/open, /answer)
	TakesArg    bool // needs free-text after the name (/deck) or a picker target
}

var commandSpecs = []commandSpec{
	{Name: "/help"},
	{Name: "/status"},
	{Name: "/protocol"},
	{Name: "/refresh"},
	{Name: "/narrate"},
	{Name: "/follow"},
	{Name: "/home"},
	{Name: "/stderr"},
	{Name: "/artifact"},
	{Name: "/run", Usage: "/run — advance the protocol now (auto-drive)"},
	{Name: "/deck", Usage: "/deck <text>", TakesArg: true},
	{Name: "/editor", Usage: "/editor — compose in $EDITOR"},
	{Name: "/open", Usage: "/open [slug|run]", OpensPicker: true, TakesArg: true},
	{Name: "/answer", Usage: "/answer [qid text]", OpensPicker: true, TakesArg: true},
	{Name: "/quit"},
}

// commandMatches returns the specs whose name has the typed command as a prefix
// ("/" alone matches all). token is the command word (before any space).
func commandMatches(input string) []commandSpec {
	token := strings.ToLower(strings.TrimSpace(input))
	var out []commandSpec
	for _, s := range commandSpecs {
		if strings.HasPrefix(s.Name, token) {
			out = append(out, s)
		}
	}
	return out
}

// recomputeSuggest refreshes the autocomplete menu after any input edit. It is
// visible only while the input is a bare slash command (starts with "/", no
// space) with at least one match and enough vertical room.
func (m *liveModel) recomputeSuggest() {
	m.suggest = false
	m.suggestItems = nil
	if !strings.HasPrefix(m.inputText, "/") || strings.Contains(m.inputText, " ") {
		m.suggestIndex = 0
		return
	}
	// Suppress on very short terminals. transcriptHeight() is clamped to >=3, so
	// check the raw (pre-clamp) available rows.
	if tuiHeight(m.height, defaultLiveHeight)-7 < 3 {
		return
	}
	matches := commandMatches(m.inputText)
	if len(matches) == 0 {
		m.suggestIndex = 0
		return
	}
	m.suggestItems = matches
	if m.suggestIndex >= len(matches) {
		m.suggestIndex = 0
	}
	m.suggest = true
}

// applyTabComplete completes the input to the longest common prefix of the
// matches; a single match completes the full command (+ space when it takes args).
func (m *liveModel) applyTabComplete() {
	matches := m.suggestItems
	if len(matches) == 0 {
		return
	}
	if len(matches) == 1 {
		m.inputText = matches[0].Name
		if matches[0].TakesArg {
			m.inputText += " "
		}
	} else {
		names := make([]string, len(matches))
		for i, s := range matches {
			names[i] = s.Name
		}
		m.inputText = longestCommonPrefix(names)
	}
	m.recomputeSuggest()
}

// acceptSuggestion runs/extends the highlighted command: a no-arg command runs
// immediately; an arg command is completed (a picker command opens its picker).
func (m liveModel) acceptSuggestion() (tea.Model, tea.Cmd) {
	if m.suggestIndex < 0 || m.suggestIndex >= len(m.suggestItems) {
		return m, nil
	}
	spec := m.suggestItems[m.suggestIndex]
	m.suggest, m.suggestItems = false, nil
	if spec.OpensPicker {
		// Bare /open or /answer opens the existing arrow-key picker.
		m.inputText = ""
		return m.runCommand(spec.Name)
	}
	if spec.TakesArg {
		m.inputText = spec.Name + " "
		return m, nil
	}
	m.inputText = ""
	return m.runCommand(spec.Name)
}

func longestCommonPrefix(items []string) string {
	if len(items) == 0 {
		return ""
	}
	prefix := items[0]
	for _, s := range items[1:] {
		for !strings.HasPrefix(s, prefix) {
			prefix = prefix[:len(prefix)-1]
			if prefix == "" {
				return ""
			}
		}
	}
	return prefix
}

// agentIsRunning reports whether the named agent currently has a running attempt.
func (m liveModel) agentIsRunning(agentID string) bool {
	a := m.agentByID(agentID)
	return a != nil && a.State == stateRunning
}

// agentLiveness returns the true process state ("live"/"stale"/"") of a running
// agent via the injected seam, for the RUN/STALE badge and confirm copy.
func (m liveModel) agentLiveness(agentID string) string {
	if m.opts.Liveness == nil {
		return ""
	}
	return m.opts.Liveness(agentID)
}

// staleAgentCount counts projected-running agents whose process is actually gone.
func (m liveModel) staleAgentCount() int {
	if m.opts.Liveness == nil {
		return 0
	}
	n := 0
	for i := range m.state.Agents {
		a := m.state.Agents[i]
		if a.State == stateRunning && m.agentLiveness(a.ID) == "stale" {
			n++
		}
	}
	return n
}

// updateConfirmKill handles the modal "kill <agent>? (y/N)" confirmation. It
// blocks every other key; y/enter confirms, n/esc cancels.
func (m liveModel) updateConfirmKill(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y", "enter":
		id := m.confirmKillAgent
		m.confirmKillAgent = ""
		if m.opts.KillAgent == nil {
			m.inputErr = "kill is not available for this run"
			return m, nil
		}
		msg, err := m.opts.KillAgent(id)
		if err != nil {
			m.inputErr = "kill failed: " + err.Error()
		} else if msg != "" {
			m.statusMsg = msg
		} else {
			m.statusMsg = "killed " + id
		}
		return m, readEventsCmd(filepath.Join(m.opts.RunDir, "events.jsonl"), m.offset, m.runToken)
	case "n", "N", "esc":
		m.confirmKillAgent = ""
		return m, nil
	}
	return m, nil // modal: ignore everything else
}

// --- Claude-CLI tabbed layout: navigation, buffers, input routing ---

// updateMain handles all keys on the default tabbed surface.
func (m liveModel) updateMain(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// ctrl+c is global; then the picker (when active) owns navigation/filter keys
	// BEFORE any tab/input handling can see them — this ordering is the whole
	// correctness gate for the feature.
	if msg.String() == "ctrl+c" {
		if m.opts.Cancel != nil {
			m.opts.Cancel()
		}
		return m, tea.Quit
	}
	// Modal confirm-kill is the highest-priority interceptor: it blocks every
	// other key so a stray Enter can't send a steer mid-confirmation.
	if m.confirmKillAgent != "" {
		return m.updateConfirmKill(msg)
	}
	if m.picker.Active {
		return m.updatePicker(msg)
	}
	key := msg.String()
	// ctrl+e opens $EDITOR seeded with the current composer text (tui-editor-composer).
	// Placed after the picker/confirm interceptors so those keep owning keys, but
	// before the suggestion menu and ordinary editing so it works in any compose state.
	if key == "ctrl+e" {
		m.suggest = false
		return m, openEditorCmd(m.inputText)
	}
	slash := strings.HasPrefix(m.inputText, "/")
	// Slash-command autocomplete: while the suggestion menu is visible it owns
	// Tab/↑/↓/Enter/Esc; printable keys still edit the input below and re-filter.
	if m.suggest {
		switch key {
		case "tab":
			m.applyTabComplete()
			return m, nil
		case "up":
			if m.suggestIndex > 0 {
				m.suggestIndex--
			}
			return m, nil
		case "down":
			if m.suggestIndex < len(m.suggestItems)-1 {
				m.suggestIndex++
			}
			return m, nil
		case "enter":
			return m.acceptSuggestion()
		case "esc":
			m.suggest, m.suggestItems = false, nil
			return m, nil
		}
	}
	switch key {
	case "esc":
		if m.composing {
			m.clearComposition()
			return m, nil
		}
		if m.inputText != "" {
			m.inputText = ""
			m.inputErr = ""
			m.recomputeSuggest()
			return m, nil
		}
		return m, tea.Quit
	case "ctrl+p":
		// Cycle the protocol ribbon: collapsed → expanded → hidden.
		m.ribbonMode = (m.ribbonMode + 1) % 3
		return m, nil
	case "ctrl+k":
		// Kill the focused agent (modal confirm). Only when on a running agent tab
		// of a live run with the kill seam wired.
		if id, ok := agentTab(m.activeTabResolved()); ok && m.agentIsRunning(id) {
			if m.opts.KillAgent == nil {
				m.inputErr = "kill is not available for this run"
			} else {
				m.confirmKillAgent = id
				m.inputErr = ""
			}
		}
		return m, nil
	case "N":
		// Start a new idea — only as a command (empty input) so it never eats
		// a capital N typed mid-steer.
		if !m.composing && m.inputText == "" && !m.picker.Active && m.opts.Start != nil {
			m.composing = true
			m.answerQID = ""
			m.inputErr = ""
			return m, nil
		}
		m.inputText += "N"
		m.inputErr = ""
		m.recomputeSuggest()
		return m, nil
	case "tab":
		// Conditional Tab: switch tabs ONLY when the input is not slash-prefixed
		// (when it is, Tab drives autocomplete — handled in the suggest block above).
		if !slash {
			m.switchTab(1)
		}
		return m, nil
	case "shift+tab":
		if !slash {
			m.switchTab(-1)
		}
		return m, nil
	case "up", "left":
		m.switchTab(-1)
		return m, nil
	case "down", "right":
		m.switchTab(1)
		return m, nil
	case "shift+up":
		m.scrollActive(-1)
		return m, nil
	case "shift+down":
		m.scrollActive(1)
		return m, nil
	case "pgup":
		m.scrollActive(-m.transcriptHeight())
		return m, nil
	case "pgdown":
		m.scrollActive(m.transcriptHeight())
		return m, nil
	case "ctrl+u":
		m.scrollActive(-m.transcriptHeight() / 2)
		return m, nil
	case "ctrl+d":
		m.scrollActive(m.transcriptHeight() / 2)
		return m, nil
	case "home":
		m.scrollActiveTop()
		return m, nil
	case "end":
		m.scrollActiveBottom()
		return m, nil
	case "enter":
		return m.submitInput()
	case "backspace", "ctrl+h":
		if r := []rune(m.inputText); len(r) > 0 {
			m.inputText = string(r[:len(r)-1])
		}
		m.inputErr = ""
		m.recomputeSuggest()
		return m, nil
	}
	if len(msg.Runes) > 0 {
		m.inputText += string(msg.Runes)
		m.inputErr = ""
		m.recomputeSuggest()
	}
	return m, nil
}

// tabIDs are Home first, then the active run's agent tabs, then Status (only
// when a run is attached). No-run ⇒ [Home] only.
func (m liveModel) tabIDs() []string {
	ids := []string{homeTabID}
	for _, a := range m.state.Agents {
		ids = append(ids, "agent:"+a.ID)
	}
	if m.hasRun() {
		ids = append(ids, statusTabID)
	}
	return ids
}

// activeTabResolved resolves the active tab: an explicit selection if valid,
// else Home (no run, or opened via `parley tui`), else the first running agent.
func (m liveModel) activeTabResolved() string {
	if m.activeTab != "" {
		for _, id := range m.tabIDs() {
			if id == m.activeTab {
				return id
			}
		}
	}
	if !m.hasRun() || m.opts.Home {
		return homeTabID
	}
	for _, a := range m.state.Agents {
		if a.State == stateRunning {
			return "agent:" + a.ID
		}
	}
	if len(m.state.Agents) > 0 {
		return "agent:" + m.state.Agents[0].ID
	}
	return homeTabID
}

func agentTab(tabID string) (string, bool) {
	if strings.HasPrefix(tabID, "agent:") {
		return strings.TrimPrefix(tabID, "agent:"), true
	}
	return "", false
}

func (m *liveModel) switchTab(delta int) {
	ids := m.tabIDs()
	if len(ids) == 0 {
		return
	}
	cur := m.activeTabResolved()
	idx := 0
	for i, id := range ids {
		if id == cur {
			idx = i
			break
		}
	}
	idx = (idx + delta) % len(ids)
	if idx < 0 {
		idx += len(ids)
	}
	m.activeTab = ids[idx]
	if agentID, ok := agentTab(m.activeTab); ok {
		m.ensureBuffer(agentID)
	}
}

func (m liveModel) agentByID(id string) *AgentState {
	for i := range m.state.Agents {
		if m.state.Agents[i].ID == id {
			return &m.state.Agents[i]
		}
	}
	return nil
}

func (m liveModel) transcriptHeight() int {
	rows := tuiHeight(m.height, defaultLiveHeight) - 7 - m.ribbonHeight()
	if rows < 3 {
		rows = 3
	}
	return rows
}

// visibleLines is the transcript to display: committed lines (stderr filtered
// out when hidden) followed by the live partials (the streaming lines that
// rewrite in place) at the bottom.
func (b *agentBuffer) visibleLines() []transcriptLine {
	out := make([]transcriptLine, 0, len(b.lines)+3)
	for _, l := range b.lines {
		if b.hideStderr && l.Stream == transcriptStderr {
			continue
		}
		out = append(out, l)
	}
	for _, s := range []transcriptStream{transcriptStdout, transcriptSteer, transcriptStderr} {
		if s == transcriptStderr && b.hideStderr {
			continue
		}
		if text := cleanLogText(b.partial[s]); strings.TrimSpace(text) != "" {
			out = append(out, transcriptLine{Text: text, Stream: s})
		}
	}
	return out
}

func (m liveModel) bufferBottom(b *agentBuffer) int {
	body := m.transcriptHeight() - 1 // the always-on status header takes one row
	if body < 1 {
		body = 1
	}
	if bottom := len(b.visibleLines()) - body; bottom > 0 {
		return bottom
	}
	return 0
}

// ensureBuffer lazily creates and loads an agent's transcript buffer (stdout +
// stderr; the steer cursor is pointed when a steer is sent).
func (m *liveModel) ensureBuffer(agentID string) *agentBuffer {
	if m.buffers == nil {
		m.buffers = map[string]*agentBuffer{}
	}
	b := m.buffers[agentID]
	if b == nil {
		b = &agentBuffer{follow: true, partial: map[transcriptStream]string{}, crPending: map[transcriptStream]bool{}}
		m.buffers[agentID] = b
	}
	if b.partial == nil {
		b.partial = map[transcriptStream]string{}
	}
	if b.crPending == nil {
		b.crPending = map[transcriptStream]bool{}
	}
	if a := m.agentByID(agentID); a != nil {
		if b.stdout.path == "" {
			b.stdout.path = a.StdoutPath
		}
		if b.stderr.path == "" {
			b.stderr.path = a.StderrPath
		}
	}
	if !b.loaded && (b.stdout.path != "" || b.stderr.path != "") {
		m.advanceBuffer(b)
		b.loaded = true
		m.replayNarrator(b) // backfill woven protocol lines exactly once
		b.scroll = m.bufferBottom(b)
	}
	return b
}

func (m *liveModel) ensureActiveBuffer() {
	if agentID, ok := agentTab(m.activeTabResolved()); ok {
		m.ensureBuffer(agentID)
	}
}

func (m liveModel) activeBuffer() *agentBuffer {
	if agentID, ok := agentTab(m.activeTabResolved()); ok && m.buffers != nil {
		return m.buffers[agentID]
	}
	return nil
}

// advanceBuffer ingests new bytes from all of an agent's streams (stdout, stderr,
// the in-flight steer reply) into the CR-aware transcript, honoring rotation and
// the bounded scrollback.
func (m *liveModel) advanceBuffer(b *agentBuffer) {
	for _, sc := range []struct {
		c *tailCursor
		s transcriptStream
	}{
		{&b.stdout, transcriptStdout},
		{&b.stderr, transcriptStderr},
		{&b.steer, transcriptSteer},
	} {
		if sc.c.path == "" {
			continue
		}
		chunk, rotated, jumped := readAppendedChunk(sc.c)
		if rotated {
			b.partial[sc.s] = ""
			b.crPending[sc.s] = false
		}
		if jumped {
			b.trunc = true
			b.partial[sc.s] = ""
			b.crPending[sc.s] = false
		}
		if len(chunk) > 0 {
			b.lines, b.partial[sc.s], b.crPending[sc.s] = ingestTranscriptBytes(b.lines, b.partial[sc.s], b.crPending[sc.s], sc.s, chunk)
			b.lastGrowthAt = m.now // feeds the spinner-vs-silent tab glyph
		}
	}
	capped := false
	b.lines, capped = capTranscriptLines(b.lines)
	if capped {
		b.trunc = true
	}
}

// refreshBuffers advances every loaded buffer (active + visited), honoring follow.
func (m *liveModel) refreshBuffers() {
	for id, b := range m.buffers {
		if b.partial == nil {
			b.partial = map[transcriptStream]string{}
		}
		if b.crPending == nil {
			b.crPending = map[transcriptStream]bool{}
		}
		if a := m.agentByID(id); a != nil {
			if b.stdout.path == "" {
				b.stdout.path = a.StdoutPath
			}
			if b.stderr.path == "" {
				b.stderr.path = a.StderrPath
			}
		}
		if !b.loaded {
			continue
		}
		m.advanceBuffer(b)
		if b.follow {
			b.scroll = m.bufferBottom(b)
		} else if bottom := m.bufferBottom(b); b.scroll > bottom {
			b.scroll = bottom
		}
	}
}

// appendSteerEvents weaves steer outcomes into the agent's transcript: a final
// marker line when a reply completes or fails (the conversation stays scrollable).
func (m *liveModel) appendSteerEvents(events []store.Event) {
	for _, e := range events {
		if e.Type != "steer.replied" && e.Type != "steer.reply_failed" {
			continue
		}
		agentID, _ := e.Data["agent"].(string)
		b := m.buffers[agentID]
		if b == nil {
			continue
		}
		// Drain any last reply bytes, then commit the trailing (no-newline) reply
		// line before clearing the cursor, so a final answer line isn't lost.
		m.advanceBuffer(b)
		if p := cleanLogText(b.partial[transcriptSteer]); strings.TrimSpace(p) != "" {
			b.lines = append(b.lines, transcriptLine{Text: p, Stream: transcriptSteer})
		}
		marker := "[reply complete]"
		if e.Type == "steer.reply_failed" {
			marker = "[reply failed]"
		}
		b.lines = append(b.lines, transcriptLine{Text: marker, Stream: transcriptEvent})
		// The cap is the actual memory guard — re-run it after the append (the
		// pre-existing path appended after advanceBuffer's cap pass).
		var capped bool
		b.lines, capped = capTranscriptLines(b.lines)
		if capped {
			b.trunc = true
		}
		b.steer = tailCursor{} // stop tailing the finished steer
		b.partial[transcriptSteer] = ""
		if b.crPending != nil {
			b.crPending[transcriptSteer] = false
		}
		if b.follow {
			b.scroll = m.bufferBottom(b)
		}
	}
}

func (m *liveModel) scrollActive(delta int) {
	agentID, ok := agentTab(m.activeTabResolved())
	if !ok {
		return
	}
	b := m.ensureBuffer(agentID)
	b.follow = false
	b.scroll += delta
	if b.scroll < 0 {
		b.scroll = 0
	}
	if bottom := m.bufferBottom(b); b.scroll > bottom {
		b.scroll = bottom
	}
}

func (m *liveModel) scrollActiveTop() {
	if agentID, ok := agentTab(m.activeTabResolved()); ok {
		b := m.ensureBuffer(agentID)
		b.follow = false
		b.scroll = 0
	}
}

func (m *liveModel) scrollActiveBottom() {
	if agentID, ok := agentTab(m.activeTabResolved()); ok {
		b := m.ensureBuffer(agentID)
		b.follow = true
		b.scroll = m.bufferBottom(b)
	}
}

func (m liveModel) openQuestionFor(agentID string) *hitl.Question {
	for i := range m.questions {
		if m.questions[i].Agent == agentID && m.questions[i].Status == hitl.StatusOpen {
			return &m.questions[i]
		}
	}
	return nil
}

func (m liveModel) segmentFor(agentID string) string {
	if a := m.agentByID(agentID); a != nil {
		return a.Segment
	}
	return ""
}

// submitInput routes the input row on Enter: slash command, else answer the
// active agent's open question, else steer the active agent (deck on Status).
func (m liveModel) submitInput() (tea.Model, tea.Cmd) {
	text := strings.TrimSpace(m.inputText)
	if m.composing {
		// The composing branch returns before slash dispatch, so a typed "/editor"
		// would be submitted as the literal answer/task. Intercept it here so the
		// editor works while answering a picker-selected question or composing an idea.
		if text == "/editor" {
			return m, openEditorCmd("")
		}
		if text == "" {
			if m.answerQID != "" {
				m.inputErr = "type an answer first"
			} else {
				m.inputErr = "type a task for the new idea"
			}
			return m, nil
		}
		if m.answerQID != "" {
			// Keep composing/answerQID set; answerQuestion clears them only on a
			// successful write so a failed answer leaves the user in compose to retry.
			return m.answerQuestion(m.answerQID, text)
		}
		return m.launchIdea(text)
	}
	if text == "" {
		if agentID, ok := agentTab(m.activeTabResolved()); ok && m.openQuestionFor(agentID) != nil {
			m.inputErr = "type an answer first"
		}
		return m, nil
	}
	if strings.HasPrefix(text, "/") {
		return m.runCommand(text)
	}
	if active := m.activeTabResolved(); active == homeTabID {
		m.inputErr = "press N to start an idea, or /open <slug|run>"
		return m, nil
	}
	if agentID, ok := agentTab(m.activeTabResolved()); ok {
		if q := m.openQuestionFor(agentID); q != nil {
			return m.answerQuestion(q.ID, text)
		}
		return m.submitSteer(steer.TargetAgent, agentID, text)
	}
	return m.submitSteer(steer.TargetDeck, "", text)
}

// launchIdea starts a new run via the LaunchFunc and attaches to it in place.
func (m liveModel) launchIdea(task string) (tea.Model, tea.Cmd) {
	if m.opts.Start == nil {
		m.inputErr, m.composing = "launching is not available here", false
		return m, nil
	}
	res, err := m.opts.Start(LaunchRequest{Task: task})
	if err != nil {
		// Keep the user in compose with their task text intact so they can retry.
		m.inputErr = "launch failed: " + err.Error()
		return m, nil
	}
	m.statusMsg = "launched " + res.RunID + " (" + strings.Join(res.Participants, ", ") + ")"
	return m, m.activateRun(res)
}

// openRun attaches the TUI to an existing run observationally (no cancel).
func (m liveModel) openRun(target string) (tea.Model, tea.Cmd) {
	run, err := runstate.ResolveRun(m.opts.Root, target)
	if err != nil {
		m.inputErr = "open failed: " + err.Error()
		return m, nil
	}
	m.statusMsg = "opened " + run.RunID
	res := LaunchResult{
		Idea:         protocol.IdeaStatus{Slug: run.IdeaSlug, Status: run.CurrentRound},
		Participants: run.Participants,
		RunID:        run.RunID,
		RunDir:       run.RunDir,
		Done:         nil, // observational: no done-wait
		Cancel:       nil, // observational: not cancelable
	}
	// An opened run is observational, but it can still durably kill its agents and
	// show RUN/STALE via their persisted process identity.
	if m.opts.ReattachKill != nil {
		res.KillAgent = m.opts.ReattachKill(run.RunID)
	}
	if m.opts.ReattachLiveness != nil {
		res.Liveness = m.opts.ReattachLiveness(run.RunID)
	}
	return m, m.activateRun(res)
}

func (m liveModel) submitSteer(target steer.Target, agentID, text string) (tea.Model, tea.Cmd) {
	// When a live seam is injected (TUI-launched run), the steer EXECUTES: an
	// agent-targeted steer spawns a fresh attempt whose stdout reply is tailed
	// inline into the agent's tab. Without a seam (observational /open run) it
	// falls back to recording only.
	if m.opts.SubmitSteer != nil {
		res, err := m.opts.SubmitSteer(SteerRequest{Target: target, AgentID: agentID, Text: text})
		if err != nil {
			m.inputErr = err.Error()
			return m, nil
		}
		m.inputText, m.inputErr = "", ""
		if target == steer.TargetAgent && res.StdoutPath != "" {
			// Weave the steer into the agent's scrollable transcript: a "❯ you:"
			// line, then tail the reply stdout as a steer stream (streams in place).
			b := m.ensureBuffer(agentID)
			// Flatten a multi-line (editor-composed) steer for the echo only; the raw
			// text was already submitted above so the agent still gets every line.
			b.lines = append(b.lines, transcriptLine{Text: "❯ you: " + editorPreview(text), Stream: transcriptSteer})
			b.steer = tailCursor{path: res.StdoutPath}
			b.partial[transcriptSteer] = ""
			b.follow = true
			b.scroll = m.bufferBottom(b)
			if res.Status == "queued" {
				m.statusMsg = agentID + " busy — steer queued (ctrl+k the active attempt to run it sooner)"
			} else {
				m.statusMsg = agentID + " is replying…"
			}
		} else {
			m.statusMsg = "recorded steer " + res.ID
		}
		return m, readEventsCmd(filepath.Join(m.opts.RunDir, "events.jsonl"), m.offset, m.runToken)
	}

	res, err := steer.Submit(m.opts.RunDir, steer.Request{
		Target:    target,
		Agent:     agentID,
		Text:      text,
		CreatedBy: "tui",
		SegmentID: m.segmentFor(agentID),
	}, time.Now().UTC())
	if err != nil {
		m.inputErr = err.Error()
		return m, nil
	}
	m.statusMsg = fmt.Sprintf("recorded steer %s (this run is observational — open the live run to get a reply)", res.ID)
	m.inputText, m.inputErr = "", ""
	return m, readEventsCmd(filepath.Join(m.opts.RunDir, "events.jsonl"), m.offset, m.runToken)
}

func (m liveModel) answerQuestion(qid, text string) (tea.Model, tea.Cmd) {
	if _, err := hitl.New(m.opts.RunDir).Answer(qid, text, false); err != nil {
		// On failure keep any answer-composition state so the user can edit + retry.
		m.inputErr = err.Error()
		return m, nil
	}
	// Success: exit answer composition cleanly (no-op for the non-picker callers).
	m.composing = false
	m.answerQID = ""
	m.statusMsg = "answered " + qid
	m.inputText, m.inputErr = "", ""
	return m, readQuestionsCmd(m.opts.RunDir, m.runToken)
}

func (m liveModel) runCommand(text string) (tea.Model, tea.Cmd) {
	fields := strings.Fields(text)
	cmd := fields[0]
	rest := strings.TrimSpace(strings.TrimPrefix(text, cmd))
	switch cmd {
	case "/help":
		m.mode = modeHelp
		m.inputText = ""
		return m, nil
	case "/quit":
		return m, tea.Quit
	case "/editor":
		// Clear the "/editor" token so a cancel doesn't leave it in the composer;
		// seed the editor empty (ctrl+e is the "edit current text" entry point).
		m.inputText, m.inputErr = "", ""
		return m, openEditorCmd("")
	case "/status", "/protocol":
		m.activeTab = statusTabID
		m.inputText, m.inputErr = "", ""
		return m, m.scheduleProtoRefresh() // the tab shows reconciled state
	case "/refresh":
		m.inputText, m.inputErr = "", ""
		m.statusMsg = "reconciling protocol state…"
		return m, m.scheduleProtoRefresh()
	case "/run":
		m.inputText, m.inputErr = "", ""
		if m.opts.StartAutoDrive == nil {
			m.inputErr = "auto-drive is not available for this run"
			return m, nil
		}
		m.opts.StartAutoDrive()
		m.statusMsg = "auto-drive started — advancing the protocol (watch the Protocol tab)"
		return m, m.scheduleProtoRefresh()
	case "/narrate":
		m.inputText, m.inputErr = "", ""
		switch m.narrateMode {
		case narrateProtocol:
			m.narrateMode = narrateVerbose
			m.statusMsg = "narrator: verbose (protocol + agent tool activity)"
		case narrateVerbose:
			m.narrateMode = narrateOff
			m.statusMsg = "narrator: off"
		default:
			m.narrateMode = narrateProtocol
			m.statusMsg = "narrator: protocol events"
		}
		return m, nil
	case "/follow":
		if agentID, ok := agentTab(m.activeTabResolved()); ok {
			b := m.ensureBuffer(agentID)
			b.follow = true
			b.scroll = m.bufferBottom(b)
		}
		m.inputText, m.inputErr = "", ""
		return m, nil
	case "/stderr":
		if agentID, ok := agentTab(m.activeTabResolved()); ok {
			b := m.ensureBuffer(agentID)
			b.hideStderr = !b.hideStderr
			b.scroll = m.bufferBottom(b)
			if b.hideStderr {
				m.statusMsg = "stderr hidden (/stderr to show)"
			} else {
				m.statusMsg = "stderr shown"
			}
		}
		m.inputText, m.inputErr = "", ""
		return m, nil
	case "/artifact":
		m.inputText = ""
		if agentID, ok := agentTab(m.activeTabResolved()); ok {
			b := m.ensureBuffer(agentID)
			b.showArtifact = !b.showArtifact
			m.inputErr = ""
		} else {
			m.inputErr = "open an agent tab to view its artifact"
		}
		return m, nil
	case "/deck":
		if rest == "" {
			m.inputErr = "usage: /deck <text>"
			return m, nil
		}
		m.inputText = ""
		return m.submitSteer(steer.TargetDeck, "", rest)
	case "/answer":
		if rest == "" {
			items := m.answerItems()
			if len(items) == 0 {
				m.inputText, m.inputErr = "", "no open questions"
				return m, nil
			}
			m.openPicker(pickerAnswer, "Answer a question", items)
			return m, nil
		}
		af := strings.Fields(rest)
		if len(af) < 2 {
			m.inputErr = "usage: /answer <qid> <text> (or just /answer to pick)"
			return m, nil
		}
		ans := strings.TrimSpace(strings.TrimPrefix(rest, af[0]))
		m.inputText = ""
		return m.answerQuestion(af[0], ans)
	case "/open":
		if rest == "" {
			items := m.openItems()
			if len(items) == 0 {
				m.inputText, m.inputErr = "", "nothing to open yet"
				return m, nil
			}
			m.openPicker(pickerOpen, "Open an idea or run", items)
			return m, nil
		}
		m.inputText = ""
		return m.openRun(rest)
	case "/home":
		m.activeTab = homeTabID
		m.inputText, m.inputErr = "", ""
		return m, nil
	default:
		m.inputErr = "unknown command: " + cmd + " (try /help)"
		return m, nil
	}
}

// --- command picker (tui-command-picker) ---

// pickerRows is the number of picker rows visible at once: capped at 8 and at
// the available transcript height so tiny terminals don't overflow.
func (m liveModel) pickerRows() int {
	rows := m.transcriptHeight() // already floored at 3
	if rows > 8 {
		rows = 8
	}
	if rows < 1 {
		rows = 1
	}
	return rows
}

// clearPicker zeroes the picker sub-mode.
func (m *liveModel) clearPicker() { m.picker = pickerState{} }

// clearComposition exits both new-idea and answer composition cleanly so a later
// composition can never submit to a stale qid.
func (m *liveModel) clearComposition() {
	m.composing = false
	m.answerQID = ""
	m.inputText = ""
	m.inputErr = ""
}

// openPicker activates the picker over items, leaving compose/input cleared.
func (m *liveModel) openPicker(kind pickerKind, title string, items []pickerItem) {
	m.composing = false
	m.answerQID = ""
	m.inputText = ""
	m.inputErr = ""
	m.picker = pickerState{Active: true, Kind: kind, Title: title, Items: items}
	m.picker.reclamp(m.pickerRows())
}

// refreshPickerItems rebuilds an open picker's candidate list from the latest
// cached data when a background tick changes the source, preserving the filter
// and selection cursor and re-clamping (FINAL §8). A no-op when no picker is open.
func (m *liveModel) refreshPickerItems() {
	if !m.picker.Active {
		return
	}
	switch m.picker.Kind {
	case pickerOpen:
		m.picker.Items = m.openItems()
	case pickerAnswer:
		m.picker.Items = m.answerItems()
	}
	m.picker.reclamp(m.pickerRows())
}

// updatePicker owns ↑/↓/Enter/esc/printable/backspace while the picker is active;
// every other key is a no-op (scroll/paging keys are intentionally ignored).
func (m liveModel) updatePicker(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	visible := m.pickerRows()
	switch msg.String() {
	case "esc":
		m.clearPicker()
		m.inputErr = ""
		return m, nil
	case "up":
		m.picker.Index--
	case "down":
		m.picker.Index++
	case "enter":
		items := m.picker.filtered()
		m.picker.reclamp(visible)
		if len(items) == 0 {
			return m, nil // empty state: select nothing, keep the picker open
		}
		return m.selectPickerItem(items[m.picker.Index])
	case "backspace", "ctrl+h":
		if r := []rune(m.picker.Filter); len(r) > 0 {
			m.picker.Filter = string(r[:len(r)-1])
			m.picker.Index = 0
		}
	default:
		if len(msg.Runes) > 0 {
			m.picker.Filter += string(msg.Runes)
			m.picker.Index = 0
		}
	}
	m.picker.reclamp(visible)
	m.inputErr = ""
	return m, nil
}

// selectPickerItem dispatches a confirmed selection by the active picker kind.
func (m liveModel) selectPickerItem(item pickerItem) (tea.Model, tea.Cmd) {
	switch m.picker.Kind {
	case pickerOpen:
		m.clearPicker()
		return m.openRun(item.Value)
	case pickerAnswer:
		m.clearPicker()
		m.composing = true
		m.answerQID = item.Value
		m.inputText = ""
		m.inputErr = ""
		return m, nil
	default:
		m.clearPicker()
		return m, nil
	}
}

// openItems are the /open candidates: ideas (from Status) then recent runs,
// deduped by Value, each labelled with its kind.
func (m liveModel) openItems() []pickerItem {
	var items []pickerItem
	seen := map[string]bool{}
	for _, idea := range m.opts.Status.Ideas {
		if idea.Slug == "" || seen[idea.Slug] {
			continue
		}
		seen[idea.Slug] = true
		items = append(items, pickerItem{
			Label: fmt.Sprintf("idea  %-30s [%s]", truncateText(idea.Slug, 30), idea.Status),
			Value: idea.Slug,
		})
	}
	for _, r := range m.homeRuns {
		if r.RunID == "" || seen[r.RunID] {
			continue
		}
		seen[r.RunID] = true
		st := "active"
		if r.Terminal {
			st = strings.ToUpper(valueOr(r.Outcome, "done"))
		}
		items = append(items, pickerItem{
			Label: fmt.Sprintf("run   %-24s %-20s [%s]", truncateText(r.RunID, 24), truncateText(r.IdeaSlug, 20), st),
			Value: r.RunID,
		})
	}
	return items
}

// answerItems are the /answer candidates: the currently open questions.
func (m liveModel) answerItems() []pickerItem {
	var items []pickerItem
	for _, q := range m.questions {
		if q.Status != hitl.StatusOpen {
			continue
		}
		who := q.Agent
		if who == "" {
			who = "deck"
		}
		items = append(items, pickerItem{
			Label: fmt.Sprintf("%-26s %-8s %s", truncateText(q.ID, 26), who, truncateText(strings.TrimSpace(q.Prompt), 48)),
			Value: q.ID,
		})
	}
	return items
}

// renderPicker draws the picker list (a bordered box) shown above the input row.
func (m liveModel) renderPicker(width int) string {
	items := m.picker.filtered()
	var b strings.Builder
	title := m.picker.Title
	if m.picker.Filter != "" {
		title += "  filter: " + m.picker.Filter
	}
	b.WriteString(headerStyle.Render(truncateText(title, width-4)))
	b.WriteString("\n")
	if len(items) == 0 {
		empty := "(no matches)" // a picker only opens with candidates, so empty == filtered out
		if m.picker.Filter == "" {
			switch m.picker.Kind {
			case pickerOpen:
				empty = "(no recent runs or ideas to open)"
			case pickerAnswer:
				empty = "(no open questions to answer)"
			}
		}
		b.WriteString(mutedStyle.Render("  " + empty))
		return boxStyle.Width(clampInt(width-4, 40, 84)).Render(strings.TrimRight(b.String(), "\n"))
	}
	// reclamp a local copy so the view never mutates model state but Index/Offset
	// are always consistent (reclamp is the single source of the window math).
	visible := m.pickerRows()
	p := m.picker
	p.reclamp(visible)
	off := p.Offset
	end := off + visible
	if end > len(items) {
		end = len(items)
	}
	if off > 0 {
		b.WriteString(mutedStyle.Render("  ↑ more"))
		b.WriteString("\n")
	}
	for i := off; i < end; i++ {
		line := truncateText(items[i].Label, width-6)
		if i == p.Index {
			b.WriteString(okStyle.Render("> " + line))
		} else {
			b.WriteString("  " + mutedStyle.Render(line))
		}
		b.WriteString("\n")
	}
	if end < len(items) {
		b.WriteString(mutedStyle.Render("  ↓ more"))
		b.WriteString("\n")
	}
	b.WriteString(mutedStyle.Render("↑/↓ select · type filter · Enter choose · esc cancel"))
	return boxStyle.Width(clampInt(width-4, 40, 84)).Render(strings.TrimRight(b.String(), "\n"))
}

func (m liveModel) renderAgentTable() string {
	var b strings.Builder
	b.WriteString(sectionTitle("Agents"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("%-2s %-8s %-10s %-8s %s\n", "", "ID", "STATE", "ELAPSED", "LAST EVENT"))
	for i, agent := range m.state.Agents {
		marker := " "
		if i == m.selected {
			marker = ">"
		}
		b.WriteString(fmt.Sprintf("%-2s %-8s %-10s %-8s %s\n",
			marker,
			agent.ID,
			stateBadge(agent.State),
			formatAgentDuration(agent, m.now),
			agent.LatestEvent,
		))
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m liveModel) renderEventPane() string {
	var b strings.Builder
	b.WriteString(sectionTitle("Latest events"))
	b.WriteString("\n")
	if len(m.state.Recent) == 0 {
		b.WriteString(mutedStyle.Render("waiting for events"))
		return b.String()
	}
	for _, event := range m.state.Recent {
		b.WriteString(fmt.Sprintf("%s  %-16s %s\n",
			event.Time.Local().Format("15:04:05"),
			event.Type,
			event.Text,
		))
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m liveModel) renderQuestionsPane() string {
	var b strings.Builder
	b.WriteString(sectionTitle("Questions"))
	b.WriteString("\n")
	if len(m.questions) == 0 {
		b.WriteString(mutedStyle.Render("no questions"))
		return b.String()
	}
	for i, question := range m.questions {
		marker := " "
		if i == m.selectedQ {
			marker = ">"
		}
		prompt := question.Prompt
		if len(prompt) > 72 {
			prompt = prompt[:69] + "..."
		}
		b.WriteString(fmt.Sprintf("%-2s %-28s %-13s %-6s %s\n", marker, question.ID, question.Status, question.Risk, prompt))
	}
	if selected := m.selectedQuestion(); selected != nil {
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("Selected: %s from %s\n", selected.ID, selected.Agent))
		if selected.Prompt != "" {
			b.WriteString(fmt.Sprintf("Prompt: %s\n", selected.Prompt))
		}
		if selected.Details != "" {
			b.WriteString(selected.Details + "\n")
		}
		if selected.DefaultAnswer != "" {
			b.WriteString(mutedStyle.Render("default: "+selected.DefaultAnswer) + "\n")
		}
		if selected.Status != hitl.StatusOpen {
			b.WriteString(okStyle.Render("answer: "+selected.Answer) + "\n")
		} else {
			b.WriteString(mutedStyle.Render("open this agent's tab and type the answer + Enter (or /answer "+selected.ID+" …)") + "\n")
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m liveModel) selectedQuestion() *hitl.Question {
	if len(m.questions) == 0 || m.selectedQ < 0 || m.selectedQ >= len(m.questions) {
		return nil
	}
	return &m.questions[m.selectedQ]
}

// updateHelpMode handles keys while the help overlay is open; esc/?/q/enter
// dismiss it, ctrl+c still cancels the run.
func (m liveModel) updateHelpMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		if m.opts.Cancel != nil {
			m.opts.Cancel()
		}
		return m, tea.Quit
	case "esc", "?", "q", "enter":
		m.mode = modeOverview
		return m, nil
	}
	return m, nil
}

func (m liveModel) renderHelp(width, height int) string {
	header := headerStyle.Render("Parley Deck — Help")
	lines := []string{
		"Tabs & navigation",
		"  ↑ / ↓              switch tab (agents, then Status)",
		"  ← / → · tab        switch tab (aliases)",
		"  shift+↑ / shift+↓  scroll one line",
		"  PgUp / PgDn        scroll a page · ctrl+u/ctrl+d half · Home/End ends",
		"",
		"Input (always typeable)",
		"  type + Enter       on an agent tab: send a steer — the agent replies in",
		"                     this tab; answers its open question if one is pending",
		"  ctrl+e             compose the input in $EDITOR (also /editor); multi-line ok",
		"  ctrl+k             kill the focused agent (confirm y/N); the run continues",
		"  esc                close a reply/menu, clear the input, or detach when empty",
		"  ctrl+c             cancel the attached run, else quit",
		"",
		"Slash commands  (type / then Tab to complete; ↑/↓ + Enter to pick)",
		"  /help              this overlay        /status   jump to Status tab",
		"  /follow            re-pin to the bottom (tail)",
		"  /deck <text>       record a deck-level steer",
		"  /editor            compose the input in $EDITOR (same as ctrl+e)",
		"  /open              pick an idea/run to open (↑/↓ + Enter); /open <slug|run> direct",
		"  /answer            pick an open question, then type the answer; /answer <qid> <t> direct",
		"  /quit              detach the TUI",
		"",
		"Picker (after a bare /open or /answer)",
		"  ↑ / ↓              move the selection      type   filter the list",
		"  Enter              choose                  esc    cancel the picker",
		"",
		mutedStyle.Render("a steer to an agent runs a fresh reply attempt; esc/Enter to close this overlay"),
	}
	body := boxStyle.Width(clampInt(width-4, 40, 84)).Render(strings.Join(lines, "\n"))
	return clipLines(strings.Join([]string{header, "", body, ""}, "\n"), height)
}

// loadFocusTail reads the last <= maxFocusBytes of a log file as cleaned,
// complete lines, capped to the line + byte scrollback budget. The returned
// offset points just past the last complete (newline-terminated) line, so any
// trailing partial line is re-read once it completes (no fragmentation).
// readAppendedChunk returns the bytes appended since cursor.offset (bounded to
// the last maxFocusBytes on a huge burst) and advances the offset to EOF.
// rotated is true when the file shrank or was replaced (truncation/rotation),
// so the caller resets that stream's partial. jumped is true when older unread
// bytes were skipped (bounded scrollback) — the caller drops the leading partial.
func readAppendedChunk(c *tailCursor) (chunk []byte, rotated, jumped bool) {
	if c.path == "" {
		return nil, false, false
	}
	f, err := os.Open(c.path)
	if err != nil {
		return nil, false, false
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return nil, false, false
	}
	size := st.Size()
	if size < c.offset || (c.info != nil && !os.SameFile(c.info, st)) {
		rotated = true
		c.offset = 0
	}
	c.info = st
	if size <= c.offset {
		return nil, rotated, false
	}
	start := c.offset
	if size-start > int64(maxFocusBytes) {
		start = size - int64(maxFocusBytes)
		jumped = true
	}
	if _, err := f.Seek(start, io.SeekStart); err != nil {
		return nil, rotated, false
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, rotated, false
	}
	// Advance the offset by the bytes ACTUALLY read (not the stat size), so bytes
	// the process appended between Stat and ReadAll aren't re-read next tick.
	c.offset = start + int64(len(data))
	if jumped {
		// Started mid-file: drop the partial leading line so we don't render a
		// fragment as if it were a whole line.
		if i := bytes.IndexByte(data, '\n'); i >= 0 {
			data = data[i+1:]
		} else {
			data = nil
		}
	}
	return data, rotated, jumped
}

// ingestTranscriptBytes folds a new chunk into the committed lines + the live
// partial for one stream, honoring carriage returns: \r\n and \n COMMIT the
// current partial line; a lone \r REWRITES the live line (Codex-CLI "potom sa to
// prepíše"); committed lines are never mutated. ANSI is stripped at commit and a
// line that is empty after cleaning is dropped. The partial is byte-capped.
func ingestTranscriptBytes(lines []transcriptLine, partial string, crPending bool, stream transcriptStream, chunk []byte) ([]transcriptLine, string, bool) {
	seg := []byte(partial)
	commit := func() {
		text := cleanLogText(string(seg))
		if strings.TrimSpace(text) != "" {
			lines = append(lines, transcriptLine{Text: text, Stream: stream})
		}
		seg = seg[:0]
	}
	capSeg := func() {
		if len(seg) > partialMaxBytes {
			seg = seg[len(seg)-partialMaxBytes:]
			for len(seg) > 0 && !utf8.RuneStart(seg[0]) {
				seg = seg[1:] // keep the cap on a rune boundary (no mojibake)
			}
		}
	}
	i := 0
	for i < len(chunk) {
		// A \r seen at the end of a previous chunk is resolved here: followed by
		// \n it is a CRLF newline; otherwise it was a lone \r (rewrite).
		if crPending {
			crPending = false
			if chunk[i] == '\n' {
				commit()
				i++
				continue
			}
			seg = seg[:0] // lone \r rewrites the live line; reprocess chunk[i]
		}
		switch chunk[i] {
		case '\n':
			commit()
			i++
		case '\r':
			crPending = true // defer: disambiguate against the next byte
			i++
		default:
			j := i
			for j < len(chunk) && chunk[j] != '\n' && chunk[j] != '\r' {
				j++
			}
			seg = append(seg, chunk[i:j]...)
			capSeg()
			i = j
		}
	}
	return lines, string(seg), crPending
}

// cleanLogText strips ANSI escapes and stray carriage returns from raw log text
// (the TUI applies its own styling; raw colour/cursor escapes would corrupt the
// width/scroll math). CR handling already happened in the ingester.
func cleanLogText(s string) string {
	return strings.ReplaceAll(stripANSI(s), "\r", "")
}

// capTranscriptLines bounds the committed transcript to maxFocusLines + the byte
// budget by evicting the oldest lines, reporting whether anything was dropped.
func capTranscriptLines(lines []transcriptLine) ([]transcriptLine, bool) {
	truncated := false
	if len(lines) > maxFocusLines {
		lines = lines[len(lines)-maxFocusLines:]
		truncated = true
	}
	total := 0
	for _, l := range lines {
		total += len(l.Text) + 1
	}
	for total > maxFocusBytes && len(lines) > 1 {
		total -= len(lines[0].Text) + 1
		lines = lines[1:]
		truncated = true
	}
	return lines, truncated
}

func loadFocusTail(path string) (lines []string, offset int64, truncated bool) {
	if path == "" {
		return nil, 0, false
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, 0, false
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return nil, 0, false
	}
	size := stat.Size()
	start := int64(0)
	if size > int64(maxFocusBytes) {
		start = size - int64(maxFocusBytes)
		truncated = true
	}
	if _, err := file.Seek(start, io.SeekStart); err != nil {
		return nil, size, truncated
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, size, truncated
	}
	complete, consumed := completeLinesFrom(data, start > 0)
	capped := false
	lines, capped = capFocusLines(splitLogLines(complete))
	return lines, start + int64(consumed), truncated || capped
}

// readAppendedLines returns the complete lines written after offset and the new
// offset. The read window is bounded to maxFocusBytes: if more than that was
// appended between ticks, it skips ahead to the last maxFocusBytes (older
// unread bytes are dropped — bounded scrollback) so a giant burst or a
// newline-less megabyte line can never make the read or buffer grow unbounded.
// The offset advances only past complete lines, so a partial trailing line is
// read again next tick.
func readAppendedLines(path string, offset int64) ([]string, int64) {
	file, err := os.Open(path)
	if err != nil {
		return nil, offset
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return nil, offset
	}
	size := stat.Size()
	if size <= offset {
		return nil, offset
	}
	start := offset
	jumped := false
	if size-start > int64(maxFocusBytes) {
		start = size - int64(maxFocusBytes)
		jumped = true
	}
	if _, err := file.Seek(start, io.SeekStart); err != nil {
		return nil, offset
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, offset
	}
	complete, consumed := completeLinesFrom(data, jumped)
	return splitLogLines(complete), start + int64(consumed)
}

// completeLinesFrom returns the newline-terminated byte region of data (after
// optionally dropping a leading partial line when the read started mid-file) and
// the number of bytes consumed up to and including the last newline. Any
// trailing partial line is left unconsumed.
func completeLinesFrom(data []byte, dropLeadingPartial bool) (complete []byte, consumed int) {
	begin := 0
	if dropLeadingPartial {
		idx := bytes.IndexByte(data, '\n')
		if idx < 0 {
			return nil, len(data) // one oversized partial line: consume it, show nothing
		}
		begin = idx + 1
	}
	region := data[begin:]
	lastNewline := bytes.LastIndexByte(region, '\n')
	if lastNewline < 0 {
		return nil, begin // no complete line yet; keep the trailing partial
	}
	return region[:lastNewline+1], begin + lastNewline + 1
}

// capFocusLines bounds a focus buffer to BOTH maxFocusLines and maxFocusBytes by
// evicting the oldest lines, reporting whether anything was dropped.
func capFocusLines(lines []string) ([]string, bool) {
	truncated := false
	if len(lines) > maxFocusLines {
		lines = lines[len(lines)-maxFocusLines:]
		truncated = true
	}
	total := 0
	for _, l := range lines {
		total += len(l) + 1
	}
	for total > maxFocusBytes && len(lines) > 1 {
		total -= len(lines[0]) + 1
		lines = lines[1:]
		truncated = true
	}
	// A single retained line can still exceed the byte cap; head-truncate it so
	// the buffer honors the budget while keeping the newest content visible.
	if total > maxFocusBytes && len(lines) == 1 {
		keep := maxFocusBytes - 1
		if keep < 0 {
			keep = 0
		}
		if l := lines[0]; len(l) > keep {
			lines[0] = strings.ToValidUTF8(l[len(l)-keep:], "")
		}
		truncated = true
	}
	return lines, truncated
}

func splitLogLines(data []byte) []string {
	if len(data) == 0 {
		return nil
	}
	cleaned := stripANSI(strings.ReplaceAll(string(data), "\r", ""))
	cleaned = strings.TrimRight(cleaned, "\n")
	if cleaned == "" {
		return nil
	}
	return strings.Split(cleaned, "\n")
}

func ProjectEvents(participants []string, events []store.Event, now time.Time) RunState {
	return runstate.ProjectEvents(participants, events, now)
}

func summarizeEvent(event store.Event) EventSummary {
	return runstate.SummarizeEvent(event)
}

func readEventsCmd(path string, offset int64, token int) tea.Cmd {
	return func() tea.Msg {
		events, nextOffset, err := readEventsFromOffset(path, offset)
		return eventsMsg{events: events, offset: nextOffset, err: err, token: token}
	}
}

func readQuestionsCmd(runDir string, token int) tea.Cmd {
	return func() tea.Msg {
		questions, err := hitl.New(runDir).List()
		return questionsMsg{questions: questions, err: err, token: token}
	}
}

func eventTickCmd(token int) tea.Cmd {
	return tea.Tick(250*time.Millisecond, func(t time.Time) tea.Msg {
		return eventTickMsg{t: t, token: token}
	})
}

func elapsedTickCmd(token int) tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return elapsedTickMsg{t: t, token: token}
	})
}

func waitDoneCmd(done <-chan struct{}, token int) tea.Cmd {
	return func() tea.Msg {
		if done == nil {
			return nil
		}
		<-done
		return doneMsg{token: token}
	}
}

func readEventsFromOffset(path string, offset int64) ([]store.Event, int64, error) {
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, offset, nil
	}
	if err != nil {
		return nil, offset, err
	}
	defer file.Close()

	if _, err := file.Seek(offset, io.SeekStart); err != nil {
		return nil, offset, err
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, offset, err
	}
	lastNewline := bytes.LastIndexByte(data, '\n')
	if lastNewline < 0 {
		return nil, offset, nil
	}
	complete := data[:lastNewline+1]
	nextOffset := offset + int64(len(complete))

	lines := bytes.Split(complete, []byte{'\n'})
	events := make([]store.Event, 0, len(lines))
	for i, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		var event store.Event
		if err := json.Unmarshal(line, &event); err != nil {
			return nil, offset, fmt.Errorf("%s:%d: %w", path, i+1, err)
		}
		events = append(events, event)
	}
	return events, nextOffset, nil
}

func tailLogFile(path string, maxBytes, maxLines int) string {
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return ""
	}
	start := int64(0)
	if stat.Size() > int64(maxBytes) {
		start = stat.Size() - int64(maxBytes)
	}
	if _, err := file.Seek(start, io.SeekStart); err != nil {
		return ""
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return ""
	}
	if start > 0 {
		if idx := bytes.IndexByte(data, '\n'); idx >= 0 {
			data = data[idx+1:]
		}
	}
	if len(data) > 0 && data[len(data)-1] != '\n' {
		if idx := bytes.LastIndexByte(data, '\n'); idx >= 0 {
			data = data[:idx+1]
		} else {
			data = nil
		}
	}
	cleaned := stripANSI(strings.ReplaceAll(string(data), "\r", ""))
	lines := strings.Split(strings.TrimRight(cleaned, "\n"), "\n")
	if len(lines) > maxLines {
		lines = lines[len(lines)-maxLines:]
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func stripANSI(value string) string {
	return ansiSequence.ReplaceAllString(value, "")
}

func displayDuration(agent AgentState, now time.Time) time.Duration {
	if agent.Duration > 0 {
		return agent.Duration
	}
	if agent.State == stateRunning && !agent.StartedAt.IsZero() {
		return now.Sub(agent.StartedAt)
	}
	return 0
}

func formatDuration(value time.Duration) string {
	if value <= 0 {
		return "-"
	}
	if value < time.Second {
		return value.Round(time.Millisecond).String()
	}
	return value.Round(time.Second).String()
}

func formatAgentDuration(agent AgentState, now time.Time) string {
	duration := displayDuration(agent, now)
	if duration <= 0 {
		switch agent.State {
		case stateFinished, stateFailed, stateSkipped:
			return "0s"
		default:
			return "-"
		}
	}
	return formatDuration(duration)
}

func displayRoundStatus(status string, done bool, resume bool) string {
	if status != "" && status != "pending" {
		return status
	}
	if resume {
		return runstate.LivenessUnverified
	}
	if done {
		return "unknown"
	}
	return "running"
}
