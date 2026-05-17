---
idea: tui-workspace-sessions
status: final
author: codex
consensus-date: 2026-05-17
participants: [codex, claude, gemini, hermes]
---

## Final plan / specification

### Goal

Upgrade `parley tui` into a workspace session console that can start and supervise multiple Parley Deck runs while preserving the repository-local protocol artifacts as the source of truth.

### Scope

- Add user-local session discovery metadata under `~/.parley-deck/sessions.json`.
- Register runs created by `parley run` and runs created from the TUI.
- Show current-workspace runs/sessions in the TUI with attention-first status.
- Show selected-run event time series, HITL questions, selected-agent state, artifact paths, and stdout/stderr tails.
- Allow starting new headless-capable ideas/runs from the TUI.
- Allow multiple TUI-started runs to proceed in parallel while the TUI process remains alive.
- Allow observation/resume of prior runs from event/question artifacts and the session registry.

### Implementation details

1. Add `internal/sessionstore`.
   - Default path: `~/.parley-deck/sessions.json`.
   - Test override: `PARLEY_HOME`.
   - JSON schema version 1.
   - Dedupe key: absolute workspace root plus run ID.
   - Atomic write via temp file and rename.
   - Store safe metadata only: workspace root, run ID, idea slug, task summary already present in `run.created`, participants, created/updated timestamps, last event timestamp, and terminal flag.
   - Missing workspace/run paths must not make TUI startup fail.

2. Add a reusable run-start path.
   - Factor the current `runTask` startup sequence into a small shared package or helper that both CLI and TUI can call.
   - Preserve current `parley run` behavior and flags.
   - Startup sequence remains: discover/select participants, create idea, create run ID, append `run.created`, register session metadata, call the runner.
   - For TUI starts, use `runner.RunRoundOneAsync` and keep the handle/cancel function in the TUI model.

3. Extend shared run state.
   - Add a shared attention/status derivation for `ACTION`, `RUNNING`, `FAILED`, `DONE`, `STALE`, and `IDLE`.
   - Base it on run terminal outcome, participant state, last event age, open HITL questions, and known pending handoff/consensus conditions where available.
   - Reuse it in TUI and, where practical, CLI status output.

4. Replace the static TUI dashboard with a workspace run console.
   - Left pane: current-workspace sessions/runs, sorted by attention state and recency.
   - Main pane: selected run event time series with a bounded in-memory tail.
   - Detail pane: selected run and selected agent details, including artifact path, state, duration, last event, stdout/stderr tail, and agent-filtered events.
   - Action pane: selected run's open HITL questions and answer entry.
   - Keep existing launch-mode override controls for future TUI-started runs.
   - Avoid the existing `n` collision with question navigation; use a non-colliding binding such as `N` for new idea if needed.

5. Define quit and resume semantics clearly.
   - `q`/`esc` should never imply true detached execution.
   - If live TUI-started runs exist, the implementation must either prompt before canceling/closing or document and render the exact behavior in the footer.
   - `ctrl+c` may remain the explicit cancellation path, but multi-run cancellation should be deliberate and visible.
   - Resume in this slice means reopening/observing persisted run artifacts and continuing existing protocol handoff validation where supported by `parley resume`.

6. Keep deferred lifecycle features out of the MVP.
   - No daemon.
   - No true reattach to child processes after parent exit.
   - No PID liveness.
   - No heartbeat events.
   - No fsnotify or cross-process file locking unless implementation discovers a concrete need.

### Tests

- `internal/sessionstore`: read/write/update/dedupe, atomic write, home override, stale workspace tolerance.
- `internal/runstate`: table tests for attention/status derivation across open question, running, failed, completed, stale, and idle cases.
- `internal/tui`: rendering and key-transition tests for multiple runs, selected run details, event stream, per-agent selection, question/action banner, answer mode, and start-run input mode.
- shared run-start path: tests that existing `parley run` behavior still creates `run.created`, registers the session, and preserves no-TUI/live-TUI behavior.
- Existing tests must continue passing with `go test ./...`.

### Non-goals

- No hosted service or background daemon.
- No terminal automation of interactive third-party CLIs.
- No replacement of canonical `parley-deck/` files.
- No complex database.
- No promise that a TUI-started child process can survive parent process exit.
- No large log or secret storage under `~/.parley-deck`.

### Verification

- `go test ./...`
- `go run ./cmd/parley status --dir .`
- `go run ./cmd/parley tui --dir .` for manual UI inspection.
- Start at least one short test run from the TUI with headless-capable agents and verify:
  - the run appears in the session list,
  - event stream updates,
  - per-agent details/log tails render,
  - open questions show as `ACTION`,
  - `~/.parley-deck/sessions.json` contains safe metadata only.

## References

- Consensus: ./consensus.md
- Rounds: ./round-01/
