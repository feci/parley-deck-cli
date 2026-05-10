---
agent: claude
idea: live-run-tui
round: 2
date: 2026-05-10
responding-to: [codex/round-01, gemini/round-01, hermes/round-01]
---

## Position changes since round 1

All four round-1 proposals converged on the same core: keep the runner unchanged
behind a goroutine, treat `events.jsonl` as the single source of truth, and let a
Bubble Tea model poll the file. I am not retreating from any of my round-1
positions, but I want to make three small adjustments based on what the others
raised:

1. **Confirm `q` = detach, `ctrl+c` = cancel.** Hermes and gemini both flag that
   default TUI behavior could surprise users in CI. I want the footer to spell
   this out and the help text to repeat it, but the binding I proposed in round 1
   stands - the prompt explicitly says "quit without killing durable state".
2. **Drop the artifact column from the agent table** that codex proposed. The
   event vocabulary already encodes that via `agent.finished` vs `agent.failed`;
   a separate `artifact_ok` column is a duplicate projection of the same fact.
3. **Cap log preview to the selected agent only**, not every running/failed
   agent. Gemini asked about I/O overhead and hermes asked about minimal polling
   interval - bounding the read set to one agent at a time makes both questions
   moot.

## Responses to others

### @codex - round-01

Strong agreement on the overall shape: same goroutine model, same poll-based
tailing, same `RunLive(...)` entry point with `runDir` + done channel, same "do
not invent states the runner does not emit" discipline (your point about
`writing-artifact` is exactly right).

Two specific pushbacks:

- **`artifact` column in the agent table.** I would drop it. `agent.finished`
  already implies an ok artifact under the current runner, and `agent.failed`
  already implies the opposite. Showing both `state` and `artifact` columns
  invites a third projection ("finished but artifact missing") that the events
  do not currently support. Counter-proposal: keep five columns - `agent`,
  `state`, `elapsed`, `last event`, `log path` - and let the log path imply the
  artifact location.
- **`tab` cycles selected agent.** Fine, but I would bind `j`/`k` as well for
  parity with the existing `viewport` muscle memory in `internal/tui/app.go`.
  Document `tab` first since it is more discoverable for new users.

Agreement on the runner/TUI lifecycle race: buffered done channel, TUI returns
first, then `printRunResults` runs on stdout after the alt-screen tears down.
That matches my round-1 wire-up for `--no-tui` and the default path both.

### @gemini - round-01

Good framing of "event sourcing" and the three-pane layout - my round-1 sketch
is the same shape. Three responses to your open questions:

- **"Polling vs. Go channel for log streaming?"** Stay with disk polling for v1.
  Counter-proposal if you want lower latency without coupling: have the runner
  fan out a `<-chan Event` *in addition to* writing `events.jsonl`, and let the
  TUI prefer the channel when present and fall back to disk tail. But that is a
  later slice - it doubles the surface area for no observable user benefit at
  current event volumes.
- **"Strict file read vs. channel?"** File. Two reasons beyond robustness: it
  keeps `parley resume` (deferred) trivially achievable later, and it forces us
  to keep `events.jsonl` honest as the durable record - any state visible in the
  TUI must already be in the file.
- **"events.jsonl might not exist at TUI start."** Handle it by treating "file
  not found" as "no events yet" in the tailer. The TUI should render the
  participant list in `pending` and only transition on real events. The runner
  emits `run.created` synchronously before any agent starts, so the gap is
  bounded to milliseconds in practice.

On flicker/debounce: a 250 ms event tick plus a 1 s elapsed tick is already a
form of debouncing. I would not add an explicit rate-limiter until we see
flicker.

### @hermes - round-01

The high-level summary matches my plan. Concrete answers to your open
questions, since they are the same questions a reviewer will ask:

- **Terminal resize.** Use `tea.WindowSizeMsg` and let `lipgloss` boxes recompute
  widths from the new dimensions. The current static `View()` already needs
  this fix - it hard-codes column widths. Counter-proposal: include a resize
  test in the snapshot suite (re-render at 80x24, 120x40, and 60x20) so we do
  not regress.
- **Polling interval.** 250 ms for the events tail, 1 s for the elapsed-time
  refresh, on-demand for the log preview (only when an agent is selected and
  it is `running` or `failed`). That is the floor before reads become noisy at
  four agents.
- **Raw vs. summarized events.** Both. Render a one-line summary per event in
  the latest-events panel (`HH:MM:SS  agent.finished  codex  (3.2s)`) and let
  the user press `enter` to expand into a raw-JSON modal. If we cannot fit the
  modal in this slice, drop the expand-to-JSON and ship the summary only.
- **"Round complete" without false positives.** Trigger off the
  `round.completed` / `round.incomplete` events the runner already emits. Do
  *not* infer "all agents finished -> round complete" in the TUI - that
  duplicates runner logic and will go wrong the first time we add a phase 2.
  Counter-proposal: if the round-status event is missing when the goroutine
  exits, render `unknown` for the round status, not `completed`.

On your risk that "TUI bugs could mask runner failures": this is the strongest
argument for keeping `--no-tui` exactly as today and routing it through a
sync wrapper, not the async one. Already in my plan; flagging here for shared
visibility.

## New concerns / questions

- **`run.created` timing relative to TUI start.** The runner writes
  `run.created` synchronously in `runTask` before spawning the goroutine; the
  TUI then opens `events.jsonl`. There is a small window where the file exists
  but is empty if the OS has not flushed - acceptable but worth confirming the
  store calls `f.Sync()` or `bufio.Flush()` on each append.
- **Selecting an agent that is `pending`.** Should pressing `j`/`k` onto a
  pending agent show an empty log panel, or skip it? My preference is to show
  "waiting to start" rather than skip, because skipping makes the list feel
  buggy.
- **Re-attaching to a finished run.** If the goroutine has already exited when
  the TUI starts (e.g. a very fast round), the TUI should detect "file done,
  no `agent.started` pending" and render the final state immediately. This is
  the same code path that `parley resume` would use later; worth keeping the
  entry point parameterised by run dir for that reason.
- **Selected-agent indicator visibility.** With four agents and a single-line
  highlight, the selection should be unambiguous in mono terminals (no
  truecolor). Use a `>` prefix in addition to color so screen readers and
  16-color terminals still work.

## Current proposal

Unchanged from my round-1 artifact in all the load-bearing parts; deltas only:

1. Runner exposes `RunRoundOneAsync(ctx, opts) <-chan Result` returning a handle
   with `Done()` and `RunDir()`. `--no-tui` calls the existing sync function
   directly; the default path calls the async wrapper.
2. `tui.RunLive(handle, status, discovered)` opens an alt-screen Bubble Tea
   program. Two `tea.Cmd`s: a 250 ms events tailer (byte-offset, `\n`-bounded)
   and a 1 s render-only tick.
3. State reduction over events:
   - no event - `pending`
   - `agent.started` - `running`, record `started_at`
   - `agent.finished` - `finished`, record `duration_ms`
   - `agent.failed` - `failed`, record `error`
   - `agent.skipped` - `skipped`, record `reason`
   - `round.completed` / `round.incomplete` - round status
   - non-participant - `unknown`
4. Layout: header (idea/round/run/status), left agent table
   (`agent / state / elapsed / last-event / log-path`), right events panel
   (latest N event summaries), bottom log preview for the selected agent only
   (last 4 KiB of stdout + stderr, ANSI-stripped, whole lines).
5. Keys: `j`/`k` and `tab` select agent, `q`/`esc` detach (runner continues),
   `ctrl+c` cancel + exit. Footer states this verbatim.
6. After the TUI returns, `printRunResults` runs on stdout so non-interactive
   captures still see the outcome. `--no-tui` path stays untouched.
7. Tests: state-reduction unit test against an `events.jsonl` fixture covering
   all transitions and the `unknown`/`pending` cases; log-tail unit test for
   mid-line truncation; one async runner integration test asserting event
   ordering and channel close; `View()` snapshot tests at 80x24, 120x40, 60x20.
8. No new payload shapes, no token accounting, no HITL input, no resume - all
   deferred per the prompt's non-goals.
