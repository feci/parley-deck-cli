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
	answerMode bool
	answerText string
	answerErr  string
	logPreview string
	errText    string
	done       bool
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
		if m.answerMode {
			return m.updateAnswerMode(msg)
		}
		switch msg.String() {
		case "q", "esc":
			return m, tea.Quit
		case "ctrl+c":
			if m.opts.Cancel != nil {
				m.opts.Cancel()
			}
			return m, tea.Quit
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
		case "a", "?":
			if m.canAnswerSelectedQuestion() {
				m.answerMode = true
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
		return m, elapsedTickCmd()
	case doneMsg:
		m.done = true
		return m, readEventsCmd(filepath.Join(m.opts.RunDir, "events.jsonl"), m.offset)
	}
	return m, nil
}

func (m liveModel) View() string {
	width := m.width
	if width < 80 {
		width = 80
	}
	bodyWidth := width - 4
	if bodyWidth < 76 {
		bodyWidth = 76
	}

	header := headerStyle.Render(fmt.Sprintf(
		"Parley Deck  idea=%s  round=%s  run=%s  status=%s",
		m.opts.Idea.Slug,
		displayRoundLabel(m.opts.Idea.Status),
		m.opts.RunID,
		displayRoundStatus(m.state.RoundStatus, m.done),
	))

	leftWidth := 36
	rightWidth := bodyWidth - leftWidth - 4
	if rightWidth < 36 {
		rightWidth = 36
	}

	left := boxStyle.Width(leftWidth).Render(m.renderAgentTable())
	right := boxStyle.Width(rightWidth).Render(m.renderEventPane())
	questions := boxStyle.Width(bodyWidth).Render(m.renderQuestionsPane())
	logs := boxStyle.Width(bodyWidth).Render(m.renderLogPane())
	footerText := "Keys: j/k/tab agent  n/p question  a answer  q/esc detach TUI  ctrl+c cancel run"
	if m.opts.Resume {
		footerText = "Keys: j/k/tab agent  n/p question  a answer  q/esc close resume view"
	}
	footer := mutedStyle.Render(footerText)
	if m.errText != "" {
		footer = warnStyle.Render(m.errText) + "\n" + footer
	}
	if m.answerMode {
		footer = warnStyle.Render("Answer mode: type answer, enter submit, esc cancel") + "\n" + footer
	}

	return strings.Join([]string{
		header,
		"",
		lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", right),
		"",
		questions,
		"",
		logs,
		"",
		footer,
		"",
	}, "\n")
}

func (m liveModel) renderAgentTable() string {
	var b strings.Builder
	b.WriteString("Agents\n")
	b.WriteString(fmt.Sprintf("%-2s %-8s %-10s %-8s %s\n", "", "ID", "STATE", "ELAPSED", "LAST EVENT"))
	for i, agent := range m.state.Agents {
		marker := " "
		if i == m.selected {
			marker = ">"
		}
		b.WriteString(fmt.Sprintf("%-2s %-8s %-10s %-8s %s\n",
			marker,
			agent.ID,
			agent.State,
			formatAgentDuration(agent, m.now),
			agent.LatestEvent,
		))
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m liveModel) renderEventPane() string {
	var b strings.Builder
	b.WriteString("Latest events\n")
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
	b.WriteString("Questions\n")
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
		if selected.Details != "" {
			b.WriteString(selected.Details + "\n")
		}
		if selected.DefaultAnswer != "" {
			b.WriteString(mutedStyle.Render("default: "+selected.DefaultAnswer) + "\n")
		}
		if selected.Status != hitl.StatusOpen {
			b.WriteString(okStyle.Render("answer: "+selected.Answer) + "\n")
		} else if m.answerMode {
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
		m.answerMode = false
		m.answerText = ""
		m.answerErr = ""
		return m, nil
	case "enter":
		question := m.selectedQuestion()
		if question == nil {
			m.answerMode = false
			return m, nil
		}
		if _, err := hitl.New(m.opts.RunDir).Answer(question.ID, m.answerText, false); err != nil {
			m.answerErr = err.Error()
			return m, nil
		}
		m.answerMode = false
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
	var parts []string
	if agent.StdoutPath != "" {
		if tail := tailLogFile(agent.StdoutPath, 4096, 6); strings.TrimSpace(tail) != "" {
			parts = append(parts, "stdout:\n"+tail)
		}
	}
	if agent.StderrPath != "" {
		if tail := tailLogFile(agent.StderrPath, 4096, 6); strings.TrimSpace(tail) != "" {
			parts = append(parts, "stderr:\n"+tail)
		}
	}
	m.logPreview = strings.Join(parts, "\n\n")
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

func displayRoundStatus(status string, done bool) string {
	if status != "" && status != "pending" {
		return status
	}
	if done {
		return "unknown"
	}
	return "running"
}
