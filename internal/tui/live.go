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
	stateUnknown  = runstate.StateUnknown
)

// Bounded scrollback for the per-agent focus view: keep at most maxFocusLines
// lines and read at most maxFocusBytes from the tail, so a very long log never
// blows up memory.
const (
	maxFocusLines = 20000
	maxFocusBytes = 4 << 20 // 4 MiB
)

// liveMode is the live-run TUI view state. It replaces the previously overloaded
// answerMode/focus booleans with one explicit machine (Slice 3). modeCompose is
// reserved for the steering composer (Slice 4).
type liveMode int

const (
	modeOverview liveMode = iota
	modeAgentDetail
	modeCompose
	modeAnswerQuestion
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
}

// statusTabID is the stable id of the dashboard tab; agent tabs are "agent:<id>".
const statusTabID = "status"

// agentBuffer is one agent's bounded, follow-capable transcript over its stdout
// log, fed by offset-incremental reads (reuses loadFocusTail/readAppendedLines/
// capFocusLines). Kept per agent so switching tabs is instant and retains scroll.
type agentBuffer struct {
	path   string
	lines  []string
	offset int64
	scroll int
	follow bool
	trunc  bool
	loaded bool
}

type liveModel struct {
	opts       LiveOptions
	width      int
	height     int
	now        time.Time
	offset     int64
	events     []store.Event
	state      RunState
	selected   int
	questions  []hitl.Question
	selectedQ  int
	mode       liveMode // only modeHelp is used now (overlay opened by /help)
	answerText string
	answerErr  string
	logPreview string
	errText    string
	done       bool

	// Retained focus fields (reused by the Status dashboard's log preview).
	follow      bool
	focusAgent  string
	focusLines  []string
	focusOffset int64
	focusScroll int
	focusTrunc  bool

	composeTarget steer.Target
	composeAgent  string
	composeText   string
	composeErr    string
	statusMsg     string

	// Claude-CLI tabbed layout (tui-claude-cli-layout): the default surface.
	activeTab string                  // "" = resolve to default; else "agent:<id>" or statusTabID
	inputText string                  // persistent bottom prompt
	inputErr  string                  // transient input/route error
	buffers   map[string]*agentBuffer // per-agent transcript, lazily loaded for active+visited
}

type RunState = runstate.RunState
type AgentState = runstate.AgentState
type EventSummary = runstate.EventSummary

type eventsMsg struct {
	events []store.Event
	offset int64
	err    error
}

type questionsMsg struct {
	questions []hitl.Question
	err       error
}

type eventTickMsg time.Time
type elapsedTickMsg time.Time
type doneMsg struct{}

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
	model.ensureActiveBuffer()
	return model
}

func (m liveModel) Init() tea.Cmd {
	cmds := []tea.Cmd{
		readEventsCmd(filepath.Join(m.opts.RunDir, "events.jsonl"), m.offset),
		readQuestionsCmd(m.opts.RunDir),
		eventTickCmd(),
		elapsedTickCmd(),
	}
	if m.opts.Done != nil {
		cmds = append(cmds, waitDoneCmd(m.opts.Done))
	}
	return tea.Batch(cmds...)
}

func (m liveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		if m.mode == modeHelp {
			return m.updateHelpMode(msg)
		}
		return m.updateMain(msg)
	case eventsMsg:
		m.offset = msg.offset
		if msg.err != nil {
			m.errText = msg.err.Error()
		} else {
			m.errText = ""
			m.events = append(m.events, msg.events...)
			m.state = ProjectEvents(m.opts.Participants, m.events, m.now)
			m.ensureActiveBuffer()
			m.refreshBuffers()
		}
		if m.done {
			return m, tea.Quit
		}
		return m, nil
	case eventTickMsg:
		return m, tea.Batch(
			readEventsCmd(filepath.Join(m.opts.RunDir, "events.jsonl"), m.offset),
			readQuestionsCmd(m.opts.RunDir),
			eventTickCmd(),
		)
	case questionsMsg:
		if msg.err != nil {
			m.errText = msg.err.Error()
		} else {
			m.questions = msg.questions
		}
		return m, nil
	case elapsedTickMsg:
		m.now = time.Time(msg)
		m.state = ProjectEvents(m.opts.Participants, m.events, m.now)
		m.ensureActiveBuffer()
		m.refreshBuffers()
		return m, elapsedTickCmd()
	case doneMsg:
		m.done = true
		return m, readEventsCmd(filepath.Join(m.opts.RunDir, "events.jsonl"), m.offset)
	}
	return m, nil
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
	var main string
	if agentID, ok := agentTab(m.activeTabResolved()); ok {
		main = m.renderTranscript(agentID, width, rows)
	} else {
		main = m.renderStatusTab(width, rows)
	}
	parts := []string{
		m.renderTabStrip(width),
		"",
		clipLines(main, rows),
	}
	if banner := m.renderQuestionBanner(width); banner != "" {
		parts = append(parts, banner)
	}
	parts = append(parts, m.renderStatusLine(width), m.renderInputRow(width))
	return strings.Join(parts, "\n")
}

func (m liveModel) renderTabStrip(width int) string {
	active := m.activeTabResolved()
	parts := make([]string, 0, len(m.state.Agents)+1)
	for _, a := range m.state.Agents {
		parts = append(parts, m.styleTab(a.ID+" "+shortState(a.State), "agent:"+a.ID == active))
	}
	parts = append(parts, m.styleTab("Status", active == statusTabID))
	return truncateText(strings.Join(parts, " "), width)
}

func (m liveModel) styleTab(label string, on bool) string {
	if on {
		return headerStyle.Render("▸ " + label)
	}
	return mutedStyle.Render("  " + label)
}

func shortState(s string) string {
	switch s {
	case stateRunning:
		return "RUN"
	case stateFinished:
		return "FIN"
	case stateFailed:
		return "ERR"
	case stateSkipped:
		return "SKIP"
	case statePending:
		return "·"
	default:
		return "?"
	}
}

func (m liveModel) renderTranscript(agentID string, width, rows int) string {
	b := m.buffers[agentID]
	if b == nil || len(b.lines) == 0 {
		return mutedStyle.Render("no output yet from " + agentID + " (waiting for the agent to write stdout)")
	}
	start := b.scroll
	if start < 0 {
		start = 0
	}
	if bottom := m.bufferBottom(b); start > bottom {
		start = bottom
	}
	end := start + rows
	if end > len(b.lines) {
		end = len(b.lines)
	}
	var sb strings.Builder
	if start == 0 && b.trunc {
		sb.WriteString(mutedStyle.Render("… earlier output truncated"))
		sb.WriteString("\n")
	}
	for i := start; i < end; i++ {
		sb.WriteString(truncateText(b.lines[i], width-1))
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

// renderStatusTab reuses the original dashboard panes as the Status tab content.
func (m liveModel) renderStatusTab(width, rows int) string {
	body := strings.Join([]string{
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
	line := fmt.Sprintf("run=%s idea=%s round=%s  %s %s %s  q:%d",
		m.opts.RunID, m.opts.Idea.Slug,
		displayRoundStatus(m.state.RoundStatus, m.done, m.opts.Resume),
		label, stateStr, follow, openQ)
	out := mutedStyle.Render(truncateText(line, width))
	if m.errText != "" {
		out = warnStyle.Render(truncateText(m.errText, width)) + "\n" + out
	} else if m.statusMsg != "" {
		out = okStyle.Render(truncateText(m.statusMsg, width)) + "\n" + out
	}
	return out
}

func (m liveModel) renderInputRow(width int) string {
	active := m.activeTabResolved()
	answer := false
	label := "deck › "
	if agentID, ok := agentTab(active); ok {
		if q := m.openQuestionFor(agentID); q != nil {
			label = fmt.Sprintf("answer %s/%s › ", agentID, q.ID)
			answer = true
		} else {
			label = "steer " + agentID + " › "
		}
	}
	var row string
	if answer {
		// Colour-flip the whole row in answer mode so an answer is never accidental.
		row = warnStyle.Render(truncateText(label+m.inputText, width))
	} else {
		row = okStyle.Render(label) + truncateText(m.inputText, width-len(label))
	}
	if m.inputErr != "" {
		row += "  " + warnStyle.Render(m.inputErr)
	}
	if strings.HasPrefix(m.inputText, "/") {
		row += "\n" + mutedStyle.Render("commands: /help  /status  /follow  /deck <t>  /answer <qid> <t>  /quit")
	} else {
		row += "\n" + mutedStyle.Render("↑/↓ tabs · PgUp/PgDn scroll · Enter steer/answer · /help · ctrl+c cancel")
	}
	return row
}

// --- Claude-CLI tabbed layout: navigation, buffers, input routing ---

// updateMain handles all keys on the default tabbed surface.
func (m liveModel) updateMain(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		if m.opts.Cancel != nil {
			m.opts.Cancel()
		}
		return m, tea.Quit
	case "esc":
		if m.inputText != "" {
			m.inputText = ""
			m.inputErr = ""
			return m, nil
		}
		return m, tea.Quit
	case "up", "left", "shift+tab":
		m.switchTab(-1)
		return m, nil
	case "down", "right", "tab":
		m.switchTab(1)
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
		return m, nil
	}
	if len(msg.Runes) > 0 {
		m.inputText += string(msg.Runes)
		m.inputErr = ""
	}
	return m, nil
}

func (m liveModel) tabIDs() []string {
	ids := make([]string, 0, len(m.state.Agents)+1)
	for _, a := range m.state.Agents {
		ids = append(ids, "agent:"+a.ID)
	}
	return append(ids, statusTabID)
}

// activeTabResolved resolves the active tab, defaulting to the first running
// agent (else first agent, else Status) so the owner lands on a live transcript.
func (m liveModel) activeTabResolved() string {
	if m.activeTab != "" {
		for _, id := range m.tabIDs() {
			if id == m.activeTab {
				return id
			}
		}
	}
	for _, a := range m.state.Agents {
		if a.State == stateRunning {
			return "agent:" + a.ID
		}
	}
	if len(m.state.Agents) > 0 {
		return "agent:" + m.state.Agents[0].ID
	}
	return statusTabID
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
	rows := tuiHeight(m.height, defaultLiveHeight) - 7
	if rows < 3 {
		rows = 3
	}
	return rows
}

func (m liveModel) bufferBottom(b *agentBuffer) int {
	if bottom := len(b.lines) - m.transcriptHeight(); bottom > 0 {
		return bottom
	}
	return 0
}

// ensureBuffer lazily creates and loads an agent's transcript buffer.
func (m *liveModel) ensureBuffer(agentID string) *agentBuffer {
	if m.buffers == nil {
		m.buffers = map[string]*agentBuffer{}
	}
	b := m.buffers[agentID]
	if b == nil {
		b = &agentBuffer{follow: true}
		m.buffers[agentID] = b
	}
	if b.path == "" {
		if a := m.agentByID(agentID); a != nil {
			b.path = a.StdoutPath
		}
	}
	if !b.loaded && b.path != "" {
		b.lines, b.offset, b.trunc = loadFocusTail(b.path)
		b.loaded = true
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

// refreshBuffers advances every loaded buffer (active + visited) with the new
// stdout bytes, reloading on log truncation/rotation, honoring follow.
func (m *liveModel) refreshBuffers() {
	for id, b := range m.buffers {
		if b.path == "" {
			if a := m.agentByID(id); a != nil {
				b.path = a.StdoutPath
			}
		}
		if b.path == "" || !b.loaded {
			continue
		}
		if info, err := os.Stat(b.path); err == nil && info.Size() < b.offset {
			b.lines, b.offset, b.trunc = loadFocusTail(b.path)
			if b.follow {
				b.scroll = m.bufferBottom(b)
			}
			continue
		}
		newLines, newOff := readAppendedLines(b.path, b.offset)
		b.offset = newOff
		if len(newLines) > 0 {
			b.lines = append(b.lines, newLines...)
			capped := false
			b.lines, capped = capFocusLines(b.lines)
			if capped {
				b.trunc = true
			}
		}
		if b.follow {
			b.scroll = m.bufferBottom(b)
		} else if bottom := m.bufferBottom(b); b.scroll > bottom {
			b.scroll = bottom
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
	if text == "" {
		if agentID, ok := agentTab(m.activeTabResolved()); ok && m.openQuestionFor(agentID) != nil {
			m.inputErr = "type an answer first"
		}
		return m, nil
	}
	if strings.HasPrefix(text, "/") {
		return m.runCommand(text)
	}
	if agentID, ok := agentTab(m.activeTabResolved()); ok {
		if q := m.openQuestionFor(agentID); q != nil {
			return m.answerQuestion(q.ID, text)
		}
		return m.submitSteer(steer.TargetAgent, agentID, text)
	}
	return m.submitSteer(steer.TargetDeck, "", text)
}

func (m liveModel) submitSteer(target steer.Target, agentID, text string) (tea.Model, tea.Cmd) {
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
	if target == steer.TargetAgent {
		m.statusMsg = fmt.Sprintf("recorded %s for %s (queued; auto-exec not wired yet)", res.ID, agentID)
	} else {
		m.statusMsg = fmt.Sprintf("recorded %s for the deck (queued; auto-exec not wired yet)", res.ID)
	}
	m.inputText, m.inputErr = "", ""
	return m, readEventsCmd(filepath.Join(m.opts.RunDir, "events.jsonl"), m.offset)
}

func (m liveModel) answerQuestion(qid, text string) (tea.Model, tea.Cmd) {
	if _, err := hitl.New(m.opts.RunDir).Answer(qid, text, false); err != nil {
		m.inputErr = err.Error()
		return m, nil
	}
	m.statusMsg = "answered " + qid
	m.inputText, m.inputErr = "", ""
	return m, readQuestionsCmd(m.opts.RunDir)
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
	case "/status":
		m.activeTab = statusTabID
		m.inputText, m.inputErr = "", ""
		return m, nil
	case "/follow":
		if agentID, ok := agentTab(m.activeTabResolved()); ok {
			b := m.ensureBuffer(agentID)
			b.follow = true
			b.scroll = m.bufferBottom(b)
		}
		m.inputText, m.inputErr = "", ""
		return m, nil
	case "/deck":
		if rest == "" {
			m.inputErr = "usage: /deck <text>"
			return m, nil
		}
		m.inputText = ""
		return m.submitSteer(steer.TargetDeck, "", rest)
	case "/answer":
		af := strings.Fields(rest)
		if len(af) < 2 {
			m.inputErr = "usage: /answer <qid> <text>"
			return m, nil
		}
		ans := strings.TrimSpace(strings.TrimPrefix(rest, af[0]))
		m.inputText = ""
		return m.answerQuestion(af[0], ans)
	default:
		m.inputErr = "unknown command: " + cmd + " (try /help)"
		return m, nil
	}
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
		} else if m.mode == modeAnswerQuestion {
			b.WriteString(warnStyle.Render("answer> ") + m.answerText + "\n")
		}
	}
	if m.answerErr != "" {
		b.WriteString(warnStyle.Render(m.answerErr))
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m liveModel) selectedQuestion() *hitl.Question {
	if len(m.questions) == 0 || m.selectedQ < 0 || m.selectedQ >= len(m.questions) {
		return nil
	}
	return &m.questions[m.selectedQ]
}

// previewLineBudget is how many tail lines the overview log pane shows per
// stream, derived from the layout height (the right pane is clipped to fit
// anyway), so the preview shows as many lines as fit rather than a hard six.
func (m liveModel) previewLineBudget() int {
	return clampInt((tuiHeight(m.height, defaultLiveHeight)-12)/2, 6, 30)
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
		"  PgUp / PgDn        scroll the transcript a page",
		"  ctrl+u / ctrl+d    scroll a half page · Home/End top/bottom",
		"",
		"Input (always typeable)",
		"  type + Enter       answer the active agent's open question, else",
		"                     record a steer for it (deck steer on the Status tab)",
		"  esc                clear the input, or detach the TUI when empty",
		"  ctrl+c             cancel the run",
		"",
		"Slash commands",
		"  /help              this overlay        /status   jump to Status tab",
		"  /follow            re-pin to the bottom (tail)",
		"  /deck <text>       record a deck-level steer",
		"  /answer <qid> <t>  answer a specific question",
		"  /quit              detach the TUI",
		"",
		mutedStyle.Render("steers are recorded/queued; auto-execution is a later slice. esc/Enter to close"),
	}
	body := boxStyle.Width(clampInt(width-4, 40, 84)).Render(strings.Join(lines, "\n"))
	return clipLines(strings.Join([]string{header, "", body, ""}, "\n"), height)
}

// loadFocusTail reads the last <= maxFocusBytes of a log file as cleaned,
// complete lines, capped to the line + byte scrollback budget. The returned
// offset points just past the last complete (newline-terminated) line, so any
// trailing partial line is re-read once it completes (no fragmentation).
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

func readEventsCmd(path string, offset int64) tea.Cmd {
	return func() tea.Msg {
		events, nextOffset, err := readEventsFromOffset(path, offset)
		return eventsMsg{events: events, offset: nextOffset, err: err}
	}
}

func readQuestionsCmd(runDir string) tea.Cmd {
	return func() tea.Msg {
		questions, err := hitl.New(runDir).List()
		return questionsMsg{questions: questions, err: err}
	}
}

func eventTickCmd() tea.Cmd {
	return tea.Tick(250*time.Millisecond, func(t time.Time) tea.Msg {
		return eventTickMsg(t)
	})
}

func elapsedTickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return elapsedTickMsg(t)
	})
}

func waitDoneCmd(done <-chan struct{}) tea.Cmd {
	return func() tea.Msg {
		if done == nil {
			return nil
		}
		<-done
		return doneMsg{}
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
