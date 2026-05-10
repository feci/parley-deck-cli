---
agent: claude
idea: live-run-tui
round: 1
date: 2026-05-10
---

## Summary

Turn the current "render-then-exit" TUI into a live observability view by running the
runner on a goroutine and deriving all displayed state from `events.jsonl` plus on-disk
log tails. Re-use the existing event vocabulary; do not invent a parallel in-memory
status model. Keep `--no-tui` exactly as today (sequential runner + text summary) and
make `parley run` default to: confirm -> start runner goroutine -> open TUI that tails
events and logs until the goroutine exits or the user quits.

## Proposed approach

1. **Decouple runner from caller.** Add `runner.RunRoundOneAsync(ctx, opts) <-chan Result`
   (or a `Start` returning a small handle with `Done()`/`Wait()`). Reuse `runAgent` and
   the existing event writes verbatim - no new payload shapes for the runner. The
   handle also exposes `RunDir` so the TUI knows where `events.jsonl` lives.

2. **Tail events as the single source of truth.** A `tea.Cmd` opens
   `runs/<run-id>/events.jsonl`, keeps a file offset, and emits an `eventsTickMsg`
   every ~250 ms with any new lines. Per-agent state is reduced from the event stream:

   | event              | transition                             |
   | ------------------ | -------------------------------------- |
   | (no event yet)     | `pending`                              |
   | `agent.started`    | `running` + record `started_at`        |
   | `agent.finished`   | `finished` + record `duration_ms`      |
   | `agent.failed`     | `failed` + record `error`              |
   | `agent.skipped`    | `skipped` + record `reason`            |
   | `round.completed`  | round status -> `completed`            |
   | `round.incomplete` | round status -> `incomplete`           |

   An agent that is selected but has no event yet stays `pending`; an agent that is
   not selected this round renders as `unknown` per the prompt's vocabulary.

3. **Elapsed time.** A second `tea.Cmd` ticks every 1 s and only re-renders. For
   `running` agents, elapsed = `now - started_at`; for terminal states use
   `duration_ms` from the event. Wall-clock skew is not a concern at second
   granularity.

4. **Log previews.** For each agent currently `running` or `failed`, re-read the last
   ~4 KiB of `stdout.log` and `stderr.log`, drop a partial trailing line, ANSI-strip,
   and show the last N (e.g. 6) lines per stream in a panel. Simple, no fsnotify.

5. **Layout (single screen, no input modes yet).**

   ```
   Parley Deck  idea=live-run-tui  round=round-01  run=<id>  [running|completed|incomplete]
   +-------------------+----------------------------------+
   | Agents            | Latest events (events.jsonl tail) |
   |  codex   running  |  T agent.started codex            |
   |  claude  finished |  T agent.finished claude          |
   |  gemini  failed   |  T agent.failed  gemini           |
   |  hermes  pending  |  ...                              |
   +-------------------+----------------------------------+
   | Log preview: <selected agent> stdout/stderr (last N)  |
   +-------------------------------------------------------+
   Keys: j/k select agent  q quit (runner keeps going)  ctrl+c cancel run
   ```

   Use the existing `lipgloss` styles already in `internal/tui/app.go`; do not
   restyle.

6. **Quit semantics.**
   - `q` / `esc`: stop the TUI only. The runner goroutine keeps writing to
     `events.jsonl` and log files - this satisfies "quit without killing durable
     state".
   - `ctrl+c`: cancel the run context (kills child agents) and then exit the TUI.
     Print the same `printRunResults` summary on stdout after the alt-screen tears
     down, so non-interactive scripts still see the outcome.

7. **`--no-tui` path** stays exactly as today: synchronous `RunRoundOne`, then
   `printRunResults`. The async wrapper just runs the sync function in a goroutine
   and forwards results, so behavior parity is mechanical to verify.

8. **Wire-up in `internal/app/app.go`.** In `runTask`, after confirmation:
   - create idea + run dir + initial `run.created` event (unchanged),
   - call `runner.RunRoundOneAsync(...)`,
   - if `--no-tui`: `<-handle.Done()`, then `printRunResults`,
   - else: `tui.RunLive(handle, status, discovered)` which blocks until the user
     quits or `handle.Done()` fires and the user dismisses the final screen.

9. **Tests (small and verifiable).**
   - Unit: a "reduce events to state" function gets a fixed `events.jsonl` fixture
     and produces the expected per-agent map; covers all transitions including
     `unknown` and `pending`.
   - Unit: log-tail helper on a fixture with mid-line truncation returns whole
     lines only.
   - Integration: existing `runner_test.go` continues to pass via the sync wrapper;
     add one async test that drives a fake agent and asserts `events.jsonl`
     contains `agent.started` before `agent.finished` and the channel closes.
   - TUI rendering is exercised via a snapshot of `View()` against a hand-built
     model - no real terminal.

## Concerns / open questions

- **Confirm prompt vs. alt-screen.** `confirmLaunch` reads from `os.Stdin` before the
  TUI takes over. That keeps the TUI free of input modes this slice. Worth confirming
  we are happy leaving HITL confirmation pre-TUI; moving it inside would force a
  small input layer now.
- **`q` semantics.** "Quit without killing durable state" reads to me as detach, not
  cancel. I have proposed that, but if the team prefers `q` = cancel + exit, the
  binding should swap. Either way the choice should be visible in the footer.
- **Multiple concurrent `parley run` invocations.** Out of scope per the prompt, but
  the design assumes one run per process; the TUI is bound to a single `RunID`.
  Worth a one-line note in the help text.
- **Re-attach to an existing run.** `parley resume` is deferred, but if we model
  state purely from `events.jsonl` here, resume becomes "open the TUI pointed at an
  existing run dir" - cheap to deliver later. Worth keeping the TUI entry point
  parameterised by run dir, not by a live handle alone.
- **ANSI handling in log previews.** I default to stripping; if agents print
  progress bars (e.g. spinners with `\r`), stripping plus last-N-lines may look
  jittery. Acceptable for v1, but flagging.
- **Event payload stability.** The TUI reduces over `data.agent`, `data.duration_ms`,
  `data.error`, `data.reason`. These keys exist in the runner today; this slice
  should treat them as a small implicit schema and avoid changing them.

## Risks

- **Runner detachment on `q`.** If the goroutine outlives the TUI, an unexpected
  panic in a child agent or in the runner itself can go unseen by the user. Mitigate
  with a `defer recover` in the runner goroutine that appends a `run.failed` event;
  the file remains the durable record.
- **Partial-line reads** on `events.jsonl` and `stdout.log` while another process is
  writing. Mitigation: keep a byte offset, only advance past `\n`, retain a small
  carry buffer. The existing `Store.Append` writes whole lines under a mutex, so
  this is safe for events.
- **Render cost on every tick.** Re-reading two log files per agent every second is
  cheap at four agents, but unbounded in principle. Cap to the agents actually
  shown in the preview panel (selected agent only) to keep this O(1).
- **Terminal-size edge cases.** The current `View()` already hard-codes minimum
  widths; the live layout should keep degrading gracefully on narrow terminals
  rather than overflowing lipgloss boxes.
- **Behavior drift for `--no-tui`.** Easy to accidentally route the synchronous path
  through the async wrapper and change exit codes or ordering. Lock this down with
  the existing `runner_test.go` cases as a regression net.
- **Scope creep into token accounting / HITL / resume.** All explicitly deferred.
  The temptation to "just add a token column" or "just add `parley answer` as a
  keybind" should be resisted in this slice.
