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
	"github.com/charmbracelet/lipgloss"

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
	mode       liveMode
	answerText string
	answerErr  string
	logPreview string
	errText    string
	done       bool

	// Per-agent focus view (Slice 2): a scrollable, follow-capable viewport over
	// the focused agent's full stdout log, fed by offset-incremental reads.
	// Active when mode == modeAgentDetail.
	follow      bool
	focusAgent  string
	focusLines  []string
	focusOffset int64
	focusScroll int
	focusTrunc  bool

	// Steering composer (Slice 4): active when mode == modeCompose.
	composeTarget steer.Target
	composeAgent  string
	composeText   string
	composeErr    string
	statusMsg     string
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
	model.state = ProjectEvents(opts.Participants, nil, now)
	model.refreshLogPreview()
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
		switch m.mode {
		case modeAnswerQuestion:
			return m.updateAnswerMode(msg)
		case modeAgentDetail:
			return m.updateFocusMode(msg)
		case modeCompose:
			return m.updateComposeMode(msg)
		case modeHelp:
			return m.updateHelpMode(msg)
		}
		switch msg.String() {
		case "q", "esc":
			return m, tea.Quit
		case "ctrl+c":
			if m.opts.Cancel != nil {
				m.opts.Cancel()
			}
			return m, tea.Quit
		case "?":
			m.mode = modeHelp
			return m, nil
		case "i":
			m.openCompose(steer.TargetAgent)
			return m, nil
		case "I":
			m.openCompose(steer.TargetDeck)
			return m, nil
		case "enter", "o":
			m.enterFocus()
			return m, nil
		case "j", "down", "tab":
			m.selectNext()
			m.refreshLogPreview()
			return m, nil
		case "k", "up", "shift+tab":
			m.selectPrev()
			m.refreshLogPreview()
			return m, nil
		case "n", "]":
			m.selectNextQuestion()
			return m, nil
		case "p", "[":
			m.selectPrevQuestion()
			return m, nil
		case "a":
			if m.canAnswerSelectedQuestion() {
				m.mode = modeAnswerQuestion
				m.answerText = ""
				m.answerErr = ""
			}
			return m, nil
		}
	case eventsMsg:
		m.offset = msg.offset
		if msg.err != nil {
			m.errText = msg.err.Error()
		} else {
			m.errText = ""
			m.events = append(m.events, msg.events...)
			m.state = ProjectEvents(m.opts.Participants, m.events, m.now)
			m.clampSelection()
			m.refreshLogPreview()
			m.refreshFocus()
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
			m.answerErr = msg.err.Error()
		} else {
			m.questions = msg.questions
			m.clampQuestionSelection()
		}
		return m, nil
	case elapsedTickMsg:
		m.now = time.Time(msg)
		m.state = ProjectEvents(m.opts.Participants, m.events, m.now)
		m.refreshLogPreview()
		m.refreshFocus()
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
	if m.mode == modeCompose {
		return m.renderCompose(width, height)
	}
	if m.mode == modeAgentDetail {
		return m.renderAgentDetail(width, height)
	}
	bodyWidth := width - 4
	if bodyWidth < 76 {
		bodyWidth = 76
	}

	header := m.liveHeader("normal")
	if height < compactLiveHeight {
		return m.renderCompactLive(width, height, m.liveHeader("compact"))
	}

	leftWidth := clampInt(bodyWidth/2, 42, 54)
	rightWidth := bodyWidth - leftWidth - 4
	if rightWidth < 40 {
		rightWidth = 40
		leftWidth = bodyWidth - rightWidth - 4
	}

	leftBody := strings.Join([]string{
		m.renderAgentTable(),
		"",
		m.renderQuestionsPane(),
	}, "\n")
	rightBody := strings.Join([]string{
		m.renderEventPane(),
		"",
		m.renderLogPane(),
	}, "\n")
	usableRows := height - 6
	left := boxStyle.Width(leftWidth).Render(clipLines(leftBody, usableRows))
	right := boxStyle.Width(rightWidth).Render(clipLines(rightBody, usableRows))
	footer := m.renderLiveFooter()

	return strings.Join([]string{
		header,
		"",
		lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", right),
		"",
		footer,
		"",
	}, "\n")
}

func (m liveModel) liveHeader(layout string) string {
	header := fmt.Sprintf(
		"Parley Deck  idea=%s  round=%s  run=%s  status=%s",
		m.opts.Idea.Slug,
		displayRoundLabel(m.opts.Idea.Status),
		m.opts.RunID,
		displayRoundStatus(m.state.RoundStatus, m.done, m.opts.Resume),
	)
	if layout == "compact" {
		header += "  layout=compact"
	}
	return headerStyle.Render(header)
}

func (m liveModel) renderLiveFooter() string {
	footerText := "Keys: j/k/tab agent  enter focus  i steer  n/p question  a answer  ? help  q/esc detach TUI  ctrl+c cancel run"
	if m.opts.Resume {
		footerText = "Keys: j/k/tab agent  enter focus  i steer  n/p question  a answer  ? help  q/esc/ctrl+c close resume view"
	}
	footer := mutedStyle.Render(footerText)
	if m.errText != "" {
		footer = warnStyle.Render(m.errText) + "\n" + footer
	}
	if m.statusMsg != "" {
		footer = okStyle.Render(m.statusMsg) + "\n" + footer
	}
	if m.mode == modeAnswerQuestion {
		footer = warnStyle.Render("Answer mode: type answer, enter submit, esc cancel") + "\n" + footer
	}
	return footer
}

func (m liveModel) renderCompactLive(width, height int, header string) string {
	var b strings.Builder
	b.WriteString(header)
	b.WriteString("\n")
	if m.errText != "" {
		b.WriteString(warnStyle.Render(m.errText))
		b.WriteString("\n")
	}
	if m.mode == modeAnswerQuestion {
		b.WriteString(warnStyle.Render("Answer mode: type answer, enter submit, esc cancel"))
		b.WriteString("\n")
	}
	b.WriteString(sectionTitle("Agents"))
	b.WriteString("\n")
	b.WriteString(m.renderCompactLiveAgents(4))
	b.WriteString("\n")
	b.WriteString(sectionTitle("Latest events"))
	b.WriteString("\n")
	b.WriteString(m.renderCompactLatestEvent())
	b.WriteString("\n")
	b.WriteString(sectionTitle("Questions"))
	b.WriteString("\n")
	b.WriteString(m.renderCompactLiveQuestion())
	b.WriteString("\n")
	b.WriteString(sectionTitle("Log preview"))
	b.WriteString("\n")
	b.WriteString(m.renderCompactLogLine(width - 12))
	b.WriteString("\n")
	b.WriteString(m.renderLiveFooter())
	return clipLines(b.String(), height) + "\n"
}

func (m liveModel) renderCompactLiveAgents(limit int) string {
	if len(m.state.Agents) == 0 {
		return mutedStyle.Render("no agents")
	}
	var b strings.Builder
	for i, agent := range m.state.Agents {
		if i >= limit {
			b.WriteString(fmt.Sprintf("  ... %d more agent(s)\n", len(m.state.Agents)-i))
			break
		}
		marker := " "
		if i == m.selected {
			marker = ">"
		}
		b.WriteString(fmt.Sprintf("%s %-8s %-14s %-8s %s\n",
			marker,
			agent.ID,
			stateBadge(agent.State),
			formatAgentDuration(agent, m.now),
			truncateText(agent.LatestEvent, 52),
		))
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m liveModel) renderCompactLatestEvent() string {
	if len(m.state.Recent) == 0 {
		return mutedStyle.Render("waiting for events")
	}
	event := m.state.Recent[len(m.state.Recent)-1]
	return fmt.Sprintf("%s  %-16s %s", event.Time.Local().Format("15:04:05"), event.Type, truncateText(event.Text, 72))
}

func (m liveModel) renderCompactLiveQuestion() string {
	if len(m.questions) == 0 {
		return mutedStyle.Render("no questions")
	}
	question := m.questions[clampIndex(m.selectedQ, len(m.questions))]
	prefix := " "
	if question.Status == hitl.StatusOpen {
		prefix = warnStyle.Render("!")
	}
	text := fmt.Sprintf("%s %-24s %-13s %-10s %s", prefix, question.ID, question.Status, riskBadge(question.Risk), truncateText(question.Prompt, 72))
	if m.mode == modeAnswerQuestion {
		text += "\n" + warnStyle.Render("answer> ") + m.answerText
	}
	if m.answerErr != "" {
		text += "\n" + warnStyle.Render(m.answerErr)
	}
	return text
}

func (m liveModel) renderCompactLogLine(limit int) string {
	agent := m.selectedAgent()
	if agent == nil {
		return mutedStyle.Render("no agent selected")
	}
	line := firstVisibleLine(m.logPreview)
	if line == "" {
		return fmt.Sprintf("%s: %s", agent.ID, mutedStyle.Render("no log output yet"))
	}
	return fmt.Sprintf("%s: %s", agent.ID, truncateText(line, limit))
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

func (m liveModel) renderLogPane() string {
	agent := m.selectedAgent()
	if agent == nil {
		return "Log preview\n" + mutedStyle.Render("no agent selected")
	}
	if strings.TrimSpace(m.logPreview) == "" {
		return fmt.Sprintf("Log preview: %s\n%s", agent.ID, mutedStyle.Render("no log output yet"))
	}
	return fmt.Sprintf("Log preview: %s\n%s", agent.ID, m.logPreview)
}

func (m liveModel) updateAnswerMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		if m.opts.Cancel != nil {
			m.opts.Cancel()
		}
		return m, tea.Quit
	case "esc":
		m.mode = modeOverview
		m.answerText = ""
		m.answerErr = ""
		return m, nil
	case "enter":
		question := m.selectedQuestion()
		if question == nil {
			m.mode = modeOverview
			return m, nil
		}
		if _, err := hitl.New(m.opts.RunDir).Answer(question.ID, m.answerText, false); err != nil {
			m.answerErr = err.Error()
			return m, nil
		}
		m.mode = modeOverview
		m.answerText = ""
		m.answerErr = ""
		return m, readQuestionsCmd(m.opts.RunDir)
	case "backspace", "ctrl+h":
		if len(m.answerText) > 0 {
			runes := []rune(m.answerText)
			m.answerText = string(runes[:len(runes)-1])
		}
		return m, nil
	}
	if len(msg.Runes) > 0 {
		m.answerText += string(msg.Runes)
	}
	return m, nil
}

func (m *liveModel) selectNext() {
	if len(m.state.Agents) == 0 {
		m.selected = 0
		return
	}
	m.selected = (m.selected + 1) % len(m.state.Agents)
}

func (m *liveModel) selectPrev() {
	if len(m.state.Agents) == 0 {
		m.selected = 0
		return
	}
	m.selected--
	if m.selected < 0 {
		m.selected = len(m.state.Agents) - 1
	}
}

func (m *liveModel) clampSelection() {
	if len(m.state.Agents) == 0 {
		m.selected = 0
		return
	}
	if m.selected >= len(m.state.Agents) {
		m.selected = len(m.state.Agents) - 1
	}
}

func (m *liveModel) selectNextQuestion() {
	if len(m.questions) == 0 {
		m.selectedQ = 0
		return
	}
	m.selectedQ = (m.selectedQ + 1) % len(m.questions)
}

func (m *liveModel) selectPrevQuestion() {
	if len(m.questions) == 0 {
		m.selectedQ = 0
		return
	}
	m.selectedQ--
	if m.selectedQ < 0 {
		m.selectedQ = len(m.questions) - 1
	}
}

func (m *liveModel) clampQuestionSelection() {
	if len(m.questions) == 0 {
		m.selectedQ = 0
		return
	}
	if m.selectedQ >= len(m.questions) {
		m.selectedQ = len(m.questions) - 1
	}
}

func (m liveModel) selectedAgent() *AgentState {
	if len(m.state.Agents) == 0 || m.selected < 0 || m.selected >= len(m.state.Agents) {
		return nil
	}
	return &m.state.Agents[m.selected]
}

func (m liveModel) selectedQuestion() *hitl.Question {
	if len(m.questions) == 0 || m.selectedQ < 0 || m.selectedQ >= len(m.questions) {
		return nil
	}
	return &m.questions[m.selectedQ]
}

func (m liveModel) canAnswerSelectedQuestion() bool {
	question := m.selectedQuestion()
	return question != nil && question.Status == hitl.StatusOpen
}

func (m *liveModel) refreshLogPreview() {
	agent := m.selectedAgent()
	if agent == nil {
		m.logPreview = ""
		return
	}
	budget := m.previewLineBudget()
	var parts []string
	if agent.StdoutPath != "" {
		if tail := tailLogFile(agent.StdoutPath, maxFocusBytes, budget); strings.TrimSpace(tail) != "" {
			parts = append(parts, "stdout:\n"+tail)
		}
	}
	if agent.StderrPath != "" {
		if tail := tailLogFile(agent.StderrPath, maxFocusBytes, budget); strings.TrimSpace(tail) != "" {
			parts = append(parts, "stderr:\n"+tail)
		}
	}
	m.logPreview = strings.Join(parts, "\n\n")
}

// previewLineBudget is how many tail lines the overview log pane shows per
// stream, derived from the layout height (the right pane is clipped to fit
// anyway), so the preview shows as many lines as fit rather than a hard six.
func (m liveModel) previewLineBudget() int {
	return clampInt((tuiHeight(m.height, defaultLiveHeight)-12)/2, 6, 30)
}

// updateFocusMode handles keys while the per-agent focus view is open. j/k and
// page keys scroll (and drop follow); f toggles follow; g/G jump to top/bottom;
// tab cycles the focused agent; esc returns to the overview.
func (m liveModel) updateFocusMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		if m.opts.Cancel != nil {
			m.opts.Cancel()
		}
		return m, tea.Quit
	case "q":
		return m, tea.Quit
	case "esc":
		m.exitFocus()
		return m, nil
	case "j", "down":
		m.scrollFocus(1)
		return m, nil
	case "k", "up":
		m.scrollFocus(-1)
		return m, nil
	case "pgdown", "ctrl+f", " ":
		m.scrollFocus(m.focusBodyHeight())
		return m, nil
	case "pgup", "ctrl+b":
		m.scrollFocus(-m.focusBodyHeight())
		return m, nil
	case "g":
		m.follow = false
		m.focusScroll = 0
		return m, nil
	case "G":
		m.follow = true
		m.focusScroll = m.focusBottom()
		return m, nil
	case "f":
		m.follow = !m.follow
		if m.follow {
			m.focusScroll = m.focusBottom()
		}
		return m, nil
	case "tab":
		m.selectNext()
		m.enterFocus()
		return m, nil
	case "shift+tab":
		m.selectPrev()
		m.enterFocus()
		return m, nil
	}
	return m, nil
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

// openCompose opens the steering composer for an agent (the selected one) or
// the whole deck.
func (m *liveModel) openCompose(target steer.Target) {
	if target == steer.TargetAgent {
		agent := m.selectedAgent()
		if agent == nil {
			return
		}
		m.composeAgent = agent.ID
	} else {
		m.composeAgent = ""
	}
	m.mode = modeCompose
	m.composeTarget = target
	m.composeText = ""
	m.composeErr = ""
}

// updateComposeMode handles keys while the steering composer is open. enter
// queues the steer (recorded as steer.* events); esc cancels.
func (m liveModel) updateComposeMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		if m.opts.Cancel != nil {
			m.opts.Cancel()
		}
		return m, tea.Quit
	case "esc":
		m.mode = modeOverview
		m.composeText = ""
		m.composeErr = ""
		return m, nil
	case "enter":
		text := strings.TrimSpace(m.composeText)
		if text == "" {
			m.composeErr = "type an instruction first"
			return m, nil
		}
		result, err := steer.Submit(m.opts.RunDir, steer.Request{
			Target:    m.composeTarget,
			Agent:     m.composeAgent,
			Text:      text,
			CreatedBy: "tui",
			SegmentID: m.composeSegment(),
		}, time.Now().UTC())
		if err != nil {
			m.composeErr = err.Error()
			return m, nil
		}
		if m.composeTarget == steer.TargetAgent {
			m.statusMsg = fmt.Sprintf("recorded %s for %s (queued; auto-exec not wired yet)", result.ID, m.composeAgent)
		} else {
			m.statusMsg = fmt.Sprintf("recorded %s for the deck (queued; auto-exec not wired yet)", result.ID)
		}
		m.mode = modeOverview
		m.composeText = ""
		m.composeErr = ""
		return m, readEventsCmd(filepath.Join(m.opts.RunDir, "events.jsonl"), m.offset)
	case "backspace", "ctrl+h":
		if len(m.composeText) > 0 {
			runes := []rune(m.composeText)
			m.composeText = string(runes[:len(runes)-1])
		}
		return m, nil
	}
	if len(msg.Runes) > 0 {
		m.composeText += string(msg.Runes)
	}
	return m, nil
}

func (m liveModel) composeSegment() string {
	if m.composeTarget != steer.TargetAgent {
		return ""
	}
	for i := range m.state.Agents {
		if m.state.Agents[i].ID == m.composeAgent {
			return m.state.Agents[i].Segment
		}
	}
	return ""
}

func (m liveModel) renderCompose(width, height int) string {
	targetLabel := "the deck"
	if m.composeTarget == steer.TargetAgent {
		targetLabel = "agent " + m.composeAgent
	}
	header := headerStyle.Render("Steer " + targetLabel)
	lines := []string{
		mutedStyle.Render("Type the next instruction. It is RECORDED as a queued follow-up for the"),
		mutedStyle.Render("agent (agents are one-shot — nothing is injected into a running process)."),
		mutedStyle.Render("Note: auto-running queued steers is not wired up yet (a later slice);"),
		mutedStyle.Render("for now this captures the instruction durably. enter record · esc cancel"),
		"",
		warnStyle.Render("steer> ") + m.composeText,
	}
	if m.composeErr != "" {
		lines = append(lines, "", warnStyle.Render(m.composeErr))
	}
	body := boxStyle.Width(clampInt(width-4, 40, 84)).Render(strings.Join(lines, "\n"))
	return clipLines(strings.Join([]string{header, "", body, ""}, "\n"), height)
}

// enterFocus opens the focus view on the currently selected agent and loads the
// bounded tail of its stdout log, positioned at the bottom with follow on.
func (m *liveModel) enterFocus() {
	agent := m.selectedAgent()
	if agent == nil {
		return
	}
	m.mode = modeAgentDetail
	m.follow = true
	m.focusAgent = agent.ID
	lines, offset, truncated := loadFocusTail(agent.StdoutPath)
	m.focusLines = lines
	m.focusOffset = offset
	m.focusTrunc = truncated
	m.focusScroll = m.focusBottom()
}

func (m *liveModel) exitFocus() {
	m.mode = modeOverview
	m.focusAgent = ""
	m.focusLines = nil
	m.focusOffset = 0
	m.focusScroll = 0
	m.focusTrunc = false
}

func (m *liveModel) scrollFocus(delta int) {
	m.follow = false
	m.focusScroll += delta
	if m.focusScroll < 0 {
		m.focusScroll = 0
	}
	if bottom := m.focusBottom(); m.focusScroll > bottom {
		m.focusScroll = bottom
	}
}

// focusBottom is the top-line index that shows the last page of output.
func (m liveModel) focusBottom() int {
	bottom := len(m.focusLines) - m.focusBodyHeight()
	if bottom < 0 {
		return 0
	}
	return bottom
}

// focusBodyHeight is how many log lines the focus viewport can show, mirroring
// the chrome (header + info + footer) rendered by renderAgentDetail.
func (m liveModel) focusBodyHeight() int {
	rows := tuiHeight(m.height, defaultLiveHeight) - 8
	if rows < 3 {
		rows = 3
	}
	return rows
}

func (m liveModel) clampedFocusScroll() int {
	scroll := m.focusScroll
	if scroll < 0 {
		scroll = 0
	}
	if bottom := m.focusBottom(); scroll > bottom {
		scroll = bottom
	}
	return scroll
}

func (m liveModel) focusedAgentState() *AgentState {
	for i := range m.state.Agents {
		if m.state.Agents[i].ID == m.focusAgent {
			return &m.state.Agents[i]
		}
	}
	return nil
}

// refreshFocus appends newly written stdout bytes to the focus viewport using an
// incremental offset read. If the log was truncated/rotated (a new segment
// re-creates the file), it reloads the tail. Honors follow mode.
func (m *liveModel) refreshFocus() {
	if m.mode != modeAgentDetail {
		return
	}
	agent := m.focusedAgentState()
	if agent == nil || agent.StdoutPath == "" {
		return
	}
	path := agent.StdoutPath
	if info, err := os.Stat(path); err == nil && info.Size() < m.focusOffset {
		lines, offset, truncated := loadFocusTail(path)
		m.focusLines = lines
		m.focusOffset = offset
		m.focusTrunc = truncated
		m.focusScroll = m.clampedFocusScroll()
		if m.follow {
			m.focusScroll = m.focusBottom()
		}
		return
	}
	newLines, newOffset := readAppendedLines(path, m.focusOffset)
	m.focusOffset = newOffset
	if len(newLines) > 0 {
		m.focusLines = append(m.focusLines, newLines...)
		capped := false
		m.focusLines, capped = capFocusLines(m.focusLines)
		if capped {
			m.focusTrunc = true
		}
	}
	if m.follow {
		m.focusScroll = m.focusBottom()
	} else {
		m.focusScroll = m.clampedFocusScroll()
	}
}

func (m liveModel) renderAgentDetail(width, height int) string {
	agent := m.focusedAgentState()
	stateLabel, segment, elapsed, artifact := "-", "-", "-", "-"
	if agent != nil {
		stateLabel = stateBadge(agent.State)
		segment = valueOr(agent.Segment, "-")
		elapsed = formatAgentDuration(*agent, m.now)
		artifact = valueOr(agent.ArtifactPath, "-")
	}
	header := headerStyle.Render(fmt.Sprintf(
		"Parley Deck  agent=%s  %s  segment=%s  elapsed=%s  run=%s",
		valueOr(m.focusAgent, "?"), stateLabel, segment, elapsed, m.opts.RunID,
	))
	info := mutedStyle.Render("artifact=" + artifact)
	body := m.renderFocusBody(width-4, m.focusBodyHeight())
	footer := m.renderFocusFooter()
	return strings.Join([]string{header, info, "", body, "", footer, ""}, "\n")
}

func (m liveModel) renderFocusBody(width, rows int) string {
	if width < 8 {
		width = 8
	}
	if len(m.focusLines) == 0 {
		return mutedStyle.Render("no log output yet (the agent has not written stdout)")
	}
	start := m.clampedFocusScroll()
	end := start + rows
	if end > len(m.focusLines) {
		end = len(m.focusLines)
	}
	var b strings.Builder
	if start == 0 && m.focusTrunc {
		b.WriteString(mutedStyle.Render("… earlier output truncated"))
		b.WriteString("\n")
	}
	for i := start; i < end; i++ {
		b.WriteString(truncateText(m.focusLines[i], width))
		b.WriteString("\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m liveModel) renderFocusFooter() string {
	follow := "off"
	if m.follow {
		follow = "on"
	}
	position := "0/0"
	if len(m.focusLines) > 0 {
		position = fmt.Sprintf("%d/%d", m.clampedFocusScroll()+1, len(m.focusLines))
	}
	text := fmt.Sprintf(
		"Keys: j/k scroll  pgup/pgdn page  g/G top/bottom  f follow(%s)  tab agent  esc back  q detach  ctrl+c cancel   [%s]",
		follow, position,
	)
	footer := mutedStyle.Render(text)
	if m.errText != "" {
		footer = warnStyle.Render(m.errText) + "\n" + footer
	}
	return footer
}

func (m liveModel) renderHelp(width, height int) string {
	header := headerStyle.Render("Parley Deck — Help")
	lines := []string{
		"Overview",
		"  j / k / tab        select agent",
		"  enter / o          open agent focus view",
		"  n / p              select question",
		"  a                  answer the selected question",
		"  i / I              steer: record a follow-up prompt (queued; not auto-run yet)",
		"  ?                  toggle this help",
		"  q / esc            detach TUI (the run keeps going)",
		"  ctrl+c             cancel the run",
		"",
		"Agent focus view",
		"  j / k              scroll one line",
		"  pgup / pgdn        scroll one page",
		"  g / G              jump top / bottom",
		"  f                  toggle follow (tail) mode",
		"  tab / shift+tab    cycle focused agent",
		"  esc                back to overview",
		"",
		mutedStyle.Render("press esc, q, ? or enter to close"),
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

func displayRoundLabel(status string) string {
	if strings.TrimSpace(status) == "" {
		return "round-01"
	}
	return status
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
