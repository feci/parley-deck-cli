---
agent: codex
idea: tui-claude-cli-layout
round: 1
date: 2026-06-04
---

## Summary

Make the live TUI default to a tabbed transcript surface, not a dashboard. The shape should be:

- top: a one-line tab strip with `Status` plus one tab per agent;
- main: the selected agent's stdout transcript, loaded and refreshed through the existing bounded focus pipeline;
- bottom: an always-present prompt row and one status line.

My main recommendation for the Bubble Tea input model: **do not keep the current modal single-letter command scheme on the default screen**. Once there is an always-typeable prompt, printable keys must be text. Put commands behind slash commands and non-printable keys:

- printable characters append to the prompt;
- `enter` submits the prompt;
- `left/right` switch tabs;
- `up/down`, `pgup/pgdn`, `home/end` scroll the active transcript/status pane;
- `ctrl+c` cancels the run when live and closes resume view;
- `esc` clears the prompt if non-empty, otherwise detaches/closes;
- `?` is no longer a global printable shortcut; use `/help` or `ctrl+?` if Bubble Tea reports it reliably;
- `q`, `a`, `f`, `i`, `h`, `j`, `k` are typed into the prompt, not intercepted.

This gives the Claude/Codex CLI feel without vim modes: the user can always type naturally, and navigation uses keys that terminal users expect to be outside the text stream. It also avoids the biggest usability trap: a user typing "help me..." should not open help, and "quit..." should not detach.

## Proposed approach

Keep `internal/tui/live.go` as the implementation center. `RunLive` already owns the live Bubble Tea program (`internal/tui/live.go:123`) and the live model already receives `WindowSizeMsg`, `KeyMsg`, `eventsMsg`, `questionsMsg`, elapsed ticks, and done messages in one place (`internal/tui/live.go:155`). The redesign is a view/state refactor around that existing loop, not new runner or engine work.

State model:

- Replace `selected` as "selected agent in overview" with `activeTab int` or `activeTabID`, where tab 0 is `Status` and later tabs map to `m.state.Agents`.
- Keep `modeHelp` for an overlay (`internal/tui/live.go:52`, `internal/tui/live.go:1028`), but stop using `modeOverview`, `modeAgentDetail`, and `modeCompose` as the normal top-level workflow. The normal state should be a single tabbed view with an input buffer.
- Replace `composeText` with persistent `inputText`; keep `composeErr`/`statusMsg` equivalents for bottom-line feedback. The current composer proves manual input handling is enough for single-line Unicode append/backspace (`internal/tui/live.go:803`, `internal/tui/live.go:810`), so I would not add `bubbles/textinput` for the first slice.
- Add per-tab scroll/follow state. At minimum: `follow bool` plus `scrollByTab map[string]int`. Prefer also `buffers map[string]focusBuffer` so agent tabs can retain loaded tails and offsets independently.

The current mode-specific handlers are useful references but should not survive unchanged. Today, default overview consumes printable single-letter commands like `q`, `?`, `i`, `I`, `j`, `k`, `n`, `p`, and `a` (`internal/tui/live.go:172`). Agent focus mode does the same for `q`, `j`, `k`, `g`, `G`, and `f` (`internal/tui/live.go:679`). That is correct for a modal dashboard, but it conflicts with an always-typeable prompt. The new key routing should run in this order:

1. `ctrl+c`: cancel live run if `Cancel` exists, then quit; for resume, quit/close. This preserves current behavior (`internal/tui/live.go:175`, `internal/tui/live.go:681`).
2. If help overlay is open, only overlay dismissal keys are routed to it. Current help mode already isolates this (`internal/tui/live.go:729`).
3. Non-printable navigation: `left/right` switch tabs; `shift+tab/tab` can also cycle tabs if desired, but arrows are the required primary keys; `up/down/pgup/pgdn/home/end` scroll the active main pane and turn follow off when moving away from bottom; `end` or `ctrl+e` jumps bottom and re-enables follow.
4. `esc`: if `inputText != ""`, clear input; otherwise detach/close.
5. `enter`: submit input.
6. `backspace`/`ctrl+h`: delete one rune from input.
7. `msg.Runes`: append to input, including `?`, `q`, `a`, `f`, and `/`.

Prompt semantics:

- Empty `enter`: no-op, or if an active open HITL question exists, optionally prefill/submit should still require text. I recommend no-op with a status hint, because accidental answers are worse than an extra keystroke.
- Non-empty `enter` when the active agent has an open question: answer that question using `hitl.New(m.opts.RunDir).Answer(...)`, same API as current answer mode (`internal/tui/live.go:550`). This matches the prompt's requirement.
- Non-empty `enter` otherwise: record a steer with `steer.Submit`, same request fields as current composer (`internal/tui/live.go:783`). Keep `CreatedBy: "tui"` and `SegmentID: m.composeSegment()` behavior (`internal/tui/live.go:787`, `internal/tui/live.go:816`).
- Slash commands are only interpreted when the submitted input starts with `/`. Minimum set: `/help`, `/quit`, `/deck <text>`, `/agent <id> <text>`, `/follow`, `/status`, `/tab <id|status>`, `/answer <question-id> <text>`. Text not starting with `/` is normal steer/answer text.

This slash-command recommendation beats the alternatives:

- Single-letter actions are incompatible with always-typeable input; they steal ordinary prose.
- Empty-input heuristic (`q` quits only when input is empty) is tempting but still surprising when the prompt starts empty and the first typed character is a command. It also makes `?` impossible as literal first input.
- Modifier keys (`ctrl+q`, `ctrl+f`, etc.) are usable for a few actions, but not for the whole command surface, and terminal portability gets uneven.
- `bubbles/textinput` would provide cursor movement and editing polish, but it brings a new dependency and does not solve command routing. Since `go.mod` currently depends on Bubble Tea and lipgloss only, not `bubbles`, adding it should be a later polish slice, not required for the architecture.

Tab rendering:

- Render `Status` first, then agents in `m.state.Agents` order. `ProjectEvents` preserves participant ordering and appends unknown event agents stably (`internal/runstate/runstate.go:321`, `internal/runstate/runstate.go:388`), so the strip can reuse that order.
- Each agent tab should show `id` plus a compact state marker derived from `stateBadge`/badge colors (`internal/tui/layout.go:81`). For example: `[Status] [codex RUN] [claude FIN]`.
- Use active styling, not a box. The current `headerStyle` and badge styles already fit a one-line terminal surface (`internal/tui/app.go:118`, `internal/tui/layout.go:17`).
- Overflow: keep active tab visible, then fill left/right neighbors until width is exhausted. Show `... +N` on the clipped side(s). On very narrow terminals, show `Status`, active tab, and counts. Do not wrap the tab strip; wrapping steals transcript rows and shifts the main content.
- Preserve `tuiWidth`'s 80-column minimum and `tuiHeight` fallback (`internal/tui/layout.go:26`, `internal/tui/layout.go:33`), but remove the old compact dashboard as the default fallback. Compact should still be tab strip + clipped main + bottom rows.

Transcript buffers:

- Reuse `loadFocusTail`, `readAppendedLines`, and `capFocusLines` exactly for stdout (`internal/tui/live.go:1059`, `internal/tui/live.go:1098`, `internal/tui/live.go:1152`). They already enforce the 20k-line/4 MiB budget (`internal/tui/live.go:34`) and handle partial trailing lines without fragmentation (`internal/tui/live.go:1055`).
- Prefer **per-agent buffers** over reload-on-switch. Reload-on-switch is simpler, but it creates visible lag, loses per-agent scroll position, and makes tabbing feel unlike Claude/Codex CLI. Per-agent buffers cost up to N * 4 MiB if every agent has been visited; with the usual 2-4 participant roster this is acceptable, and it can be capped by loading only visited tabs.
- Refresh only loaded buffers on each event/elapsed tick, and always load the active agent tail when switching to an unloaded agent. The current `refreshFocus` is gated by `modeAgentDetail` (`internal/tui/live.go:928`); in the new model that condition becomes "active tab is an agent with a loaded buffer" plus optionally "loaded background buffers".
- Follow mode should be per active transcript. Any scroll-up disables follow; bottom jump re-enables it. The current focus behavior is already correct (`internal/tui/live.go:875`, `internal/tui/live.go:958`).

Status tab, events, and HITL:

- `Status` tab is the old dashboard content, adapted to a single main pane. It should include agent table, latest events, open questions, queued steers, and recent errors. Existing render helpers can be reused initially: agent table (`internal/tui/live.go:439`), events (`internal/tui/live.go:460`), questions (`internal/tui/live.go:478`).
- Do not add an `Events` tab in slice 1 unless the status tab becomes too dense. `runstate.ProjectEvents` already keeps the recent event summary at eight items (`internal/runstate/runstate.go:331`), enough for a status pane. A dedicated events tab can be a later slice if users want full event history.
- HITL questions should appear in two places: a bottom/banner hint when the active agent has an open question, and the full list/detail on `Status`. The HITL store already records questions and emits `hitl.question`/`hitl.answered` events (`internal/hitl/hitl.go:54`, `internal/hitl/hitl.go:215`, `internal/hitl/hitl.go:228`).
- Selection of questions by `n/p` should go away on the default screen. Active-agent question routing should pick the oldest open question for that agent. For runner/deck questions or multiple questions, require `/answer <question-id> <text>` or use the `Status` tab's highlighted oldest open question.

Steer and state plumbing:

- Keep `internal/steer` unchanged. Its contract is exactly right for this layout: record `steer.requested` and `steer.delivered` with `new_attempt` queued delivery, without live injection (`internal/steer/steer.go:1`, `internal/steer/steer.go:84`). The bottom status line must retain honest wording like "queued; auto-exec not wired yet", matching current copy (`internal/tui/live.go:795`).
- Keep segment plumbing unchanged. `runstate.ProjectEvents` resets targeted agents on `run.segment_started` and assigns segment IDs to later agent events (`internal/runstate/runstate.go:358`, `internal/runstate/runstate.go:432`). `composeSegment` can keep deriving the current target segment from projected agent state (`internal/tui/live.go:816`).
- Keep `--no-tui` unchanged. `runTask` bypasses `RunLive` entirely when `--no-tui` is set (`internal/app/app.go:1647`), and `resume` prints run detail instead of opening the TUI (`internal/app/app.go:993`). This proposal should not touch those branches.

Bottom rows:

- Prompt row: one physical line, for example `codex > <input>` or `answer codex/qid > <input>` when an open question is active. Clip from the left if needed so the cursor-side text remains visible.
- Status line: idea, run, round status, active tab, active agent state, segment, follow on/off, scroll position, open question count, queued steer count, and last error/status message. It should replace the verbose key footer (`internal/tui/live.go:325`) with concise state; `/help` carries the command reference.
- Resume mode should make status explicit and never imply the process is live, preserving the existing resume test expectation around unverified status (`internal/tui/live_test.go:226`).

Incremental slice plan:

1. Introduce tab state and render the new top/main/bottom frame while leaving current manual input as simple `inputText`. Default active tab should be first running agent, else first agent, else `Status`.
2. Move the focus transcript into the default agent-tab main area and add per-agent buffer structs using the existing focus read helpers. Update focus tests to assert default transcript rendering instead of enter-to-focus.
3. Replace modal composer/answer key handling with persistent prompt routing and slash commands. Keep `steer.Submit` and `hitl.Answer` call sites; remove `i/I/a` as global printable commands.
4. Port the dashboard into `Status` tab and add tab strip overflow tests at narrow widths.
5. Polish status line, queued steer counts via `steer.List`, and help text. Consider `bubbles/textinput` only here if cursor movement, paste handling, or selection becomes visibly deficient.

## Concerns / open questions

- Should `enter` with non-empty text on an agent tab always answer an open question for that agent, or should answers require an `answer` prefix? The prompt asks for answer-first behavior, and I agree, but the UI must make the active question unmistakable in the prompt label.
- What should `Status` tab prompt submissions target? I recommend defaulting to deck-level steer from `Status`, while agent tabs default to agent-level steer.
- Do we want `tab/shift+tab` as aliases for left/right tab switching? It is convenient, but tab may later be useful for prompt completion if `bubbles/textinput` or completions are added. I would support it as an alias in slice 1 and keep left/right documented as primary.
- How much background refreshing should happen for non-active agent buffers? Refreshing only visited buffers is bounded and simple. Refreshing all agents gives instant switching but can read many files every 250 ms (`internal/tui/live.go:1216`). I recommend active + visited only.

## Risks

- The biggest risk is preserving old printable shortcuts in the name of compatibility. That would undercut the core Claude/Codex CLI layout because normal typing would be unsafe. This is a deliberate behavior change and should be called out in help/tests.
- Per-agent buffers multiply the 4 MiB cap by visited agents. It is acceptable for normal Parley rosters, but the implementation should avoid preloading every possible agent in very large runs.
- Manual input handling is enough for first delivery, but it lacks cursor movement and paste-friendly editing. If users expect full prompt editing parity, `bubbles/textinput` may become necessary later.
- Tests currently assert the old dashboard panels and key footer (`internal/tui/live_test.go:119`). They will need intentional rewrites, not minor string edits, because the default surface changes from overview to transcript.
- There is a subtle accidental-answer risk when active-agent HITL questions take precedence over steering. The prompt label and status line must clearly show `answer <question-id>` before `enter` submits.
