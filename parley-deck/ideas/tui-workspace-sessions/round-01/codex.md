---
agent: codex
idea: tui-workspace-sessions
round: 1
date: 2026-05-17
---

## Summary

The TUI should become a workspace-level run console, not a second protocol store. The existing `parley-deck/runs/<run-id>/events.jsonl`, HITL question files, agent log files, and idea frontmatter are already close to the right source of truth; the missing parts are a durable cross-workspace session index, a reusable run-starting API, and a TUI model that can supervise more than one run at a time.

I would ship this as an MVP without a background daemon. Runs started from the TUI are managed by the current TUI process while it is open, and previous runs are discoverable again from persisted metadata plus workspace run artifacts.

## Proposed approach

### 1. Add a small user-local session registry

Create a narrow package, for example `internal/sessionstore`, that owns a JSON file under `~/.parley-deck/sessions.json`. Use a schema-versioned document, written atomically, with entries like:

```json
{
  "schema_version": 1,
  "sessions": [
    {
      "workspace_root": "/repo",
      "run_id": "20260517T120000.000000000Z",
      "idea_slug": "tui-workspace-sessions",
      "task": "Upgrade the TUI...",
      "participants": ["claude", "gemini"],
      "created_at": "2026-05-17T12:00:00Z",
      "updated_at": "2026-05-17T12:05:00Z",
      "last_event_at": "2026-05-17T12:04:59Z"
    }
  ]
}
```

The registry must not store secrets, logs, prompts from third-party CLIs, or canonical protocol state. It should tolerate missing workspaces and stale runs. Tests can use an env override such as `PARLEY_HOME` to avoid writing into a real home directory while production defaults to `~/.parley-deck`.

Register sessions from `parley run`, from runs started inside the TUI, and opportunistically when the TUI lists workspace runs that are not yet in the registry. This keeps the registry useful even if a run was created before the feature existed.

### 2. Factor run startup out of `internal/app`

Right now `runTask` owns idea creation, run ID creation, event initialization, auto-answer startup, and `runner.RunRoundOne` invocation. The TUI should not duplicate all of that. Add a small orchestration package such as `internal/runcontrol` with a function along these lines:

```go
type StartOptions struct {
    Root         string
    Task         string
    Participants []string
    Discovered   []agents.Discovery
    Auto         bool
    TUIManaged   bool
}

type StartedRun struct {
    Idea   protocol.IdeaStatus
    RunID  string
    RunDir string
    Handle *runner.Handle
    Cancel context.CancelFunc
}
```

`runTask` can call this package for both `--no-tui` and live-TUI paths. The workspace TUI can call the same path asynchronously when the user starts a new task. That keeps the CLI and TUI behavior aligned and gives one place to register the session in `~/.parley-deck`.

### 3. Replace the current dashboard with a run console model

Extend `internal/tui/app.go` from an idea/agent dashboard into a workspace console with these panes:

- Left pane: sessions/runs for the current workspace, sorted by newest event. Include badges such as `ACTION`, `RUNNING`, `FAILED`, `DONE`, `STALE`, and `IDLE`.
- Center pane: event time series for the selected run. This should show more than the current eight-event `RunState.Recent` window, either by loading all run events or by keeping a larger TUI-local tail.
- Right pane: selected run details, with an agent tab/selector. The active agent view should show state, duration, latest event, artifact path, stdout/stderr tail, and agent-filtered events.
- Bottom/action pane: open HITL questions for the selected run, answer entry mode, current key hints, and start-run input mode.

Suggested keys for the first slice:

- `n`: new idea/task input mode.
- `enter`: open or focus selected run.
- `tab` / `shift+tab`: cycle panes or agent tabs depending on focus.
- `j/k` or arrows: move within the focused list.
- `a`: answer selected open question.
- `r`: refresh runs and registry.
- `h/i/m/x`: keep the existing per-agent launch-mode override behavior, but apply it to the runtime config used for newly started TUI runs.
- `q`/`esc`: close or leave the current mode.

### 4. Derive user-action state from existing artifacts

The "needs user" signal should be computed, not manually tracked:

- `ACTION`: at least one open HITL question, pending manual handoff, or blocked consensus signoff.
- `RUNNING`: at least one selected participant has `agent.started` without a terminal agent event and the run is not terminal.
- `FAILED`: terminal outcome failed, round incomplete, or any failed required participant.
- `DONE`: terminal completed.
- `STALE`: non-terminal run with no recent event for a threshold, such as 10 minutes, and no open question.
- `IDLE`: known run without active process evidence.

`internal/runstate` can expose a small `RunBadge` or `Attention` helper so CLI status, TUI, and future tests do not reimplement the same triage logic.

### 5. Keep resume pragmatic in this slice

There are two different meanings of "resume":

- Resume observation: find a prior run and reopen its live or workspace view from the event log. This is achievable now.
- Resume execution after the parent TUI process exits: continue or reattach to still-running child agent processes. This requires durable process supervision or a daemon.

For the MVP, I recommend implementing resume observation plus multiple concurrent in-process TUI runs. If the TUI starts three ideas and stays open, all three can run in parallel and stream into the same console. If the TUI exits, the next launch can rediscover the sessions and show their latest persisted state, but it should not pretend it can reattach to killed child processes.

If true detached execution is required later, open a separate idea for a small supervisor process. That is a larger lifecycle and cleanup problem than this TUI slice needs.

### 6. Test plan

Add unit tests before relying on manual terminal testing:

- `internal/sessionstore`: create/update/list, atomic JSON write, schema version, stale workspace tolerance, env override for home.
- `internal/runstate`: attention badge derivation for open questions, running, failed, completed, stale, and idle.
- `internal/tui`: render multiple runs, selected run details, event stream, per-agent tab, question/action banner, and input mode transitions using synthetic run summaries.
- `internal/app` or `internal/runcontrol`: starting a run writes `run.created`, registers a session, and preserves existing `parley run` behavior.

Manual verification should then use a fake or very short configured agent before exercising real hosted CLIs.

## Concerns / open questions

- The current `tui.Run(status, discovered)` signature does not accept root, run summaries, or callbacks. It needs either a new options struct or a new `RunWorkspace` entrypoint.
- `runstate.ProjectEvents` currently stores only eight recent events. That is fine for status summaries but too small for a useful event time series.
- The existing live TUI code has useful file-tail and question-answer helpers. We should reuse or move helpers rather than create a second implementation with different behavior.
- It is unclear whether the user expects true detached execution after quitting the TUI. I recommend explicitly documenting the MVP behavior and deferring daemon/process supervision unless needed.
- Interactive third-party agent sessions should remain out of scope for automatic TUI orchestration. The TUI can display/choose launch modes, but parallel managed runs should use headless-capable agents.

## Risks

- Starting multiple real hosted agents from a TUI can consume quota quickly. The start flow should show selected participants and require an explicit start action, while preserving existing `--yes` semantics for command-line use.
- A home-level session registry can become stale if workspaces are moved or deleted. The UI must show stale entries as stale and never fail startup because of them.
- If the TUI reads whole event logs forever, very long sessions may become slow. The first implementation can cap the display while keeping the file as source of truth.
- Factoring `runTask` into a reusable run-control package touches the main CLI path. Keep the change small and protect it with tests around existing `parley run`, `status`, and `resume` behavior.
