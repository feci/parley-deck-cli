---
idea: tui-action-execution
author: codex
created: 2026-05-25
participants: [codex, claude, gemini, hermes]
roles:
  codex: current CLI/TUI architecture and implementation slice owner
  claude: terminal UX and operator feedback lens
  gemini: planner/action command contract and test strategy lens
  hermes: long-running execution, recovery, and operational safety lens
status: final
---

## Problem / idea

The `continuous-run-tui` slice added continuation planning, persisted `run.json`
next actions, and a read-only TUI action display. The next slice should make the
TUI action area operational without turning the TUI into a broad shell wrapper.

Design and implement a narrow first step where a user can focus planner actions,
select one action, and trigger supported safe actions from the TUI. Unsupported,
interactive, or higher-risk actions must produce a clear in-TUI message that
shows what is missing or what command should be run manually.

## Constraints

- Keep the TUI reusable and testable: `internal/tui` should receive an injected
  action runner instead of importing `internal/app` or shelling out directly.
- Reuse the existing planner `NextAction` data model and the existing CLI command
  semantics where possible.
- Avoid a daemon, background scheduler, or silent phase advancement.
- Execution must be asynchronous from Bubble Tea's perspective so long-running
  actions do not freeze rendering.
- Refresh run/idea state after successful mutating actions.
- Preserve the current read-only planner output for actions that cannot yet be
  performed directly.
- Use hosted/headless participant CLIs for this idea; do not introduce local
  model servers or local agent runtimes.

## Non-goals

- Building full native input flows for answering human questions.
- Implementing true agent conversation continuation or retry scheduling.
- Replacing the existing CLI command surface.
- Adding broad workflow automation beyond the selected action.
