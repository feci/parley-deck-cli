---
idea: tui-workspace-sessions
drafted-by: codex
date: 2026-05-17
---

## Agreed decisions

1. The TUI will become a workspace session console for runs, not a second protocol store.
   - Repository-local `parley-deck/runs/<run-id>/events.jsonl`, HITL question files, agent logs, and idea files remain canonical.
   - User-local metadata under `~/.parley-deck` is only an advisory discovery and resume index.

2. Add a small `internal/sessionstore` package for `~/.parley-deck/sessions.json`.
   - The file is schema-versioned JSON, written atomically.
   - Each entry stores only safe metadata: workspace root, run ID, idea slug, task summary if already present in `run.created`, participants, created/updated timestamps, last event timestamp, and terminal flag.
   - It must tolerate stale/missing workspaces and dedupe by workspace root plus run ID.
   - Tests use an override such as `PARLEY_HOME` so they do not write to the real home directory.

3. Factor run startup into a reusable orchestration path.
   - Existing `parley run` behavior must be preserved.
   - The workspace TUI should start new ideas through the same run creation path as the CLI: create idea, create run ID, append `run.created`, register the session, and call `runner.RunRoundOneAsync` for TUI-managed runs.
   - This avoids duplicating startup logic inside `internal/tui`.

4. The first TUI slice uses a workspace-level multi-run model.
   - The left pane lists current-workspace sessions/runs, sorted by attention state and recency.
   - The main pane shows a time-ordered event stream for the selected run.
   - The detail pane shows selected run state, selected agent state, artifact path, stdout/stderr tail, and agent-filtered events.
   - The action pane shows open HITL questions and answer entry.
   - The existing agent launch-mode override controls stay available and apply to future TUI-started runs.

5. The TUI can start multiple headless-capable ideas in parallel while it is open.
   - Runs started by the TUI are represented by in-memory handles and cancel functions.
   - Previously created runs are loaded from workspace run artifacts and session metadata for observation/resume.
   - True detached execution after the parent process exits is out of scope for this MVP; the UI must not imply it can reattach to killed child processes.

6. User-action status is derived, not manually tracked.
   - `ACTION`: open HITL questions, pending/manual handoff, or blocked/pending consensus work.
   - `RUNNING`: non-terminal run with at least one running participant and no open question.
   - `FAILED`: failed run, incomplete terminal round, or failed required participant.
   - `DONE`: completed terminal run with no pending action.
   - `STALE`: non-terminal run with no recent event beyond a conservative threshold and no open question.
   - `IDLE`: known run without active process evidence.
   - The derivation should live in `internal/runstate` or another shared package so CLI status and TUI do not diverge.

7. Event and log rendering should reuse existing live-TUI mechanics.
   - Keep offset-based event reads.
   - Maintain a larger in-memory event tail for the selected run, while preserving the existing compact recent window for summaries.
   - Reuse or move existing log-tailing and HITL answer helpers instead of creating a separate behavior.

8. Resume means observation and protocol continuation in this slice.
   - `parley resume` remains supported.
   - TUI resume opens a selected run from persisted event/question artifacts and validates pending consensus signoffs where existing CLI behavior already does so.
   - Durable background supervision, daemon behavior, heartbeats, and PID-based liveness are deferred.

9. Tests are required for the state and persistence model.
   - `internal/sessionstore`: read/write/update/dedupe, atomic writes, stale paths, and home override.
   - `internal/runstate`: attention/status derivation across action, running, failed, completed, stale, and idle cases.
   - `internal/tui`: rendering and key transitions for multiple runs, selected run details, event stream, per-agent view, question/action banner, and start-run input.
   - run startup path: `parley run` still writes `run.created`, registers sessions, and preserves existing no-TUI/live-TUI behavior.

## Agreed trade-offs

- Prefer a small JSON index under `~/.parley-deck` over a database or daemon.
- Prefer current event logs and question files as the source of truth over a separate TUI state file.
- Prefer a conservative MVP without true detached execution. This reduces lifecycle risk and keeps the first implementation shippable.
- Prefer polling with bounded tails in the first slice. Heartbeats, file locks, fsnotify, PID tracking, and background polling backoff can be added after the basic model is proven.
- Preserve current CLI commands and flags. The TUI adds another entrypoint; it must not break `parley run`, `parley resume`, `parley status`, or `parley answer`.

## Open items deferred to implementation

- Exact keymap names can change during implementation if conflicts appear, especially around `n`, `tab`, `q`, and `ctrl+c`.
- The first UI can use compact panes instead of a full separate per-agent screen if the selected-agent panel gives enough visibility.
- Cross-workspace recent sessions may be listed later than current-workspace runs if that keeps the first cut smaller.
- Heartbeat events, PID liveness, fsnotify, true detached supervisor/daemon behavior, and global question queue are follow-up candidates, not MVP requirements.

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: codex - 2026-05-17
Status: ✅ ACCEPT
Notes: Accept. MVP keeps canonical state in run artifacts, adds only advisory session metadata, and defers daemon-style detached execution.

### Signoff: claude - 2026-05-17
Status: 🟡 ACCEPT-WITH-RESERVATIONS
Notes: Accept. Canonical-state rule, advisory `~/.parley-deck/sessions.json`, shared run startup path, derived action status in `internal/runstate`, per-run keying, and resume-as-observation all match the correctness lens from round-01. Two reservations to handle in implementation, not blockers: (a) the `n` keybinding already means "next question" in `live.go` - pick a non-colliding key (e.g. `N`) for "new idea" before shipping; (b) the consensus does not specify what quitting the TUI does to TUI-started runs - since true detachment is out of scope, the implementation must either cancel cleanly with confirmation or document that quit cancels, so users do not assume background continuation.

### Signoff: gemini - 2026-05-17
Status: ✅ ACCEPT
Notes: The state model correctly prioritizes repository-local artifacts as the source of truth, with `sessions.json` serving as an advisory index. Centralizing derived status logic in `internal/runstate` is critical for consistency. I echo Claude's reservation regarding clear lifecycle management for TUI-started runs to ensure that the "in-memory handle" model doesn't lead to user confusion about background persistence.

### Signoff: hermes - 2026-05-17
Status: ✅ ACCEPT
Notes: Accept. The MVP clarifies handling of local `~/.parley-deck` metadata, including observation and failure visibility.
