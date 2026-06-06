---
agent: codex
idea: tui-command-picker
round: 1
date: 2026-06-06
---
## Summary

Use one reusable picker sub-mode on `liveModel`, not per-command modes. The model already has one orthogonal overlay mode (`modeHelp`) and one input sub-mode (`composing`); a picker should be another input sub-state that temporarily owns navigation, filtering, confirm, and cancel keys while leaving the existing tabs, transcript buffers, and slash-command execution paths intact.

The most important implementation detail is key routing order. If the picker is active, it must intercept `up`, `down`, `enter`, `esc`, printable runes, and backspace before the existing branches that switch tabs, submit answers/steers, clear input, or start composing with `N`.

## Proposed approach

Add small value types near `liveModel`:

```go
type pickerKind string

const (
	pickerOpen   pickerKind = "open"
	pickerAnswer pickerKind = "answer"
)

type pickerItem struct {
	Label string
	Value string
	Kind  string
}

type pickerState struct {
	Active bool
	Kind   pickerKind
	Title  string
	Items  []pickerItem
	Index  int
	Filter string
}
```

Put `picker pickerState` on `liveModel`. Avoid storing callbacks in the picker. A function field would be awkward for tests, copies of Bubble Tea models, and future equality/debugging. `Kind` plus a single `selectPickerItem` method is simpler Go and keeps command effects in ordinary model methods:

```go
func (m liveModel) selectPickerItem(item pickerItem) (tea.Model, tea.Cmd) {
	switch m.picker.Kind {
	case pickerOpen:
		m.clearPicker()
		return m.openRun(item.Value)
	case pickerAnswer:
		m.answerQID = item.Value
		m.clearPicker()
		m.composing = true
		m.inputText = ""
		m.inputErr = ""
		return m, nil
	default:
		m.clearPicker()
		return m, nil
	}
}
```

For filtering, expose a derived view and clamp selection every time it is used:

```go
func (p pickerState) filtered() []pickerItem {
	if strings.TrimSpace(p.Filter) == "" {
		return p.Items
	}
	needle := strings.ToLower(strings.TrimSpace(p.Filter))
	var out []pickerItem
	for _, item := range p.Items {
		hay := strings.ToLower(item.Label + " " + item.Value + " " + item.Kind)
		if strings.Contains(hay, needle) {
			out = append(out, item)
		}
	}
	return out
}

func (p *pickerState) clamp() {
	n := len(p.filtered())
	if n == 0 {
		p.Index = 0
		return
	}
	if p.Index < 0 {
		p.Index = 0
	}
	if p.Index >= n {
		p.Index = n - 1
	}
}
```

I would add `answerQID string` to `liveModel` and reuse `composing` for the second step of `/answer`. A separate answer-text mode would mostly duplicate `composing` semantics: it owns the input row, Enter submits text, Esc cancels, and the field being edited is still `inputText`. The only real difference is the submit target, so make that explicit:

```go
// Unified TUI:
composing bool
answerQID  string // non-empty means composing an answer selected by /answer picker
picker     pickerState
```

Then `submitInput` starts with:

```go
if m.composing {
	if text == "" {
		if m.answerQID != "" {
			m.inputErr = "type an answer first"
		} else {
			m.inputErr = "type a task for the new idea"
		}
		return m, nil
	}
	if m.answerQID != "" {
		qid := m.answerQID
		m.answerQID = ""
		m.composing = false
		return m.answerQuestion(qid, text)
	}
	return m.launchIdea(text)
}
```

`renderInputRow` can use `answerQID` before the normal active-agent answer detection:

```go
case m.composing && m.answerQID != "":
	label = "answer " + m.answerQID + " › "
	answer = true
case m.composing:
	label = "new idea › "
```

This avoids pre-filling `inputText` with `/answer <qid> `. Pre-fill is easy to implement, but it leaks command syntax back into the flow, makes backspace/cancel behavior less clean, and allows users to accidentally edit the qid. The current model already treats `inputText` as the text payload for the active context, so keeping the qid in state is cleaner.

The exact `updateMain` placement should be immediately after `ctrl+c` and before the current `esc`, `N`, arrow, and `enter` cases:

```go
func (m liveModel) updateMain(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		if m.opts.Cancel != nil {
			m.opts.Cancel()
		}
		return m, tea.Quit
	}
	if m.picker.Active {
		return m.updatePicker(msg)
	}
	switch msg.String() {
	case "esc":
		...
	case "N":
		if !m.composing && m.inputText == "" && !m.picker.Active && m.opts.Start != nil {
			m.composing = true
			m.answerQID = ""
			m.inputErr = ""
			return m, nil
		}
		...
	case "up", "left", "shift+tab":
		m.switchTab(-1)
	...
	case "enter":
		return m.submitInput()
	}
	...
}
```

The explicit `!m.picker.Active` in the `N` guard is redundant if the early picker branch stays in place, but worth keeping because it documents the intended invariant: `N` is a command only when `!composing && inputText == "" && !picker.Active`. While picker filtering is active, a typed `N` should become filter text, not start a new idea.

`updatePicker` should own only picker keys:

```go
func (m liveModel) updatePicker(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.clearPicker()
		return m, nil
	case "up":
		m.picker.Index--
	case "down":
		m.picker.Index++
	case "enter":
		items := m.picker.filtered()
		m.picker.clamp()
		if len(items) == 0 {
			m.inputErr = "no matches"
			return m, nil
		}
		return m.selectPickerItem(items[m.picker.Index])
	case "backspace", "ctrl+h":
		if r := []rune(m.picker.Filter); len(r) > 0 {
			m.picker.Filter = string(r[:len(r)-1])
		}
	default:
		if len(msg.Runes) > 0 {
			m.picker.Filter += string(msg.Runes)
		}
	}
	m.picker.clamp()
	m.inputErr = ""
	return m, nil
}
```

For `/open`, keep `/open <slug|run-id>` as-is. For bare `/open`, build items from cached `opts.Status.Ideas` and `homeRuns`, deduping by `Value` so an idea slug and run id can both appear but duplicate run ids do not. Put ideas first, then recent runs, matching Home. Empty candidates should not activate the picker; set `inputErr = "nothing to open yet"`.

For `/answer`, keep `/answer <qid> <text>` as-is. For bare `/answer`, build items from `questions` with `StatusOpen`. Empty candidates should set `inputErr = "no open questions"`. Selection sets `answerQID`, enters `composing`, clears `inputText`, and renders an answer row. The active agent tab shortcut remains unchanged: if an agent tab already has an open question, typing text + Enter still answers that question.

Rendering should be simple and testable. Add a `renderPicker(width, rows)` list above the input row, probably inserted in `renderTabbed` after the question banner and before `renderStatusLine`. The input hint should switch while active: `↑/↓ select · type filter · Enter choose · esc cancel`. `/help` and slash-command hints should mention bare `/open` and `/answer` opening pickers.

Tests should mirror the existing `live_test.go` model-driving style: no terminal, no Bubble Tea program. I would add table-driven tests around direct `Update(tea.KeyMsg{...})` calls and assert state fields.

Recommended tests:

- Bare `/open` activates `pickerOpen`, clears or leaves command input consistently, and starts at index 0.
- With picker active, `down` changes `picker.Index` and does not change `activeTabResolved()`.
- With picker active, printable runes update `picker.Filter` and do not append to `inputText`.
- Filter narrowing clamps index when the current index is beyond the filtered length.
- `enter` on an `/open` picker calls the same open path as explicit `/open <value>`; for a pure unit test, selecting an invalid value can assert `inputErr` from `openRun`, or a temp workspace run can assert `opts.RunID`.
- `esc` with picker active cancels only the picker and does not quit or clear unrelated model state unexpectedly.
- Bare `/answer` with open questions activates `pickerAnswer`; selecting a question sets `composing=true`, `answerQID=<qid>`, and leaves `inputText=""`.
- Enter while composing an answer calls `answerQuestion`, clears `answerQID`, clears `composing`, and persists the answer via `hitl.New(runDir).List()`.
- `N` is gated: when picker active, rune `N` goes to `picker.Filter`; when no picker, not composing, and empty input, it opens the new-idea composer.

## Concerns / open questions

The only state-shape concern is overloading `composing`. I think it is still the right fit if `answerQID` is treated as the discriminator, but all cancellation paths must clear both fields. `activateRun`, `esc` from composing, successful `launchIdea`, and successful `answerQuestion` should leave `composing=false`, `answerQID=""`, `picker.Active=false`, and `inputText=""`.

The `/open` candidate value needs a deliberate choice. For ideas, `Value` should be the slug. For runs, `Value` should be the run id. Labels should include enough context to distinguish them, for example `idea  tui-command-picker  round-01` and `run   20260606-...  tui-command-picker  active`. Avoid making labels part of command semantics.

Refreshing data on picker open would be nice, but the current architecture has cached `homeRuns` and `questions` already fed by ticks/messages. For this slice I would use cached state only. A later async refresh-on-open can be added without changing picker state.

## Risks

The main regression risk is key collision from placing the picker branch too late. If it comes after `up/down/enter/esc`, the feature will look implemented but still switch tabs, submit steers, or quit. This should be covered by routing tests, especially active-tab invariance under picker arrows.

The second risk is stale answer state. If `answerQID` survives cancellation or run activation, later `composing` input could answer the wrong question. Clear it aggressively anywhere `composing` is cancelled or the run changes.

The third risk is testing only rendering. Rendering tests are useful for hints, but the important contract is model state and effects after key messages. The picker needs table-driven state-machine tests first; rendering can be a small follow-up assertion that the picker and hint appear in `View()`.
