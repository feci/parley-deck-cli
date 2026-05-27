# Changelog

## v1.5.4 - 2026-05-27

- Treat ACP as a selectable launch mode on an existing agent instead of
  exposing duplicate `*-acp` agent IDs for Codex, Claude, and Hermes.
- Add the TUI `a` key for session-only ACP launch overrides and show ACP
  command details in the selected-agent panel.
- Add `acp_args` runtime configuration so local installs can enable ACP for
  CLIs when their concrete ACP launch args are known.
- Apply TUI launch-mode overrides to newly started runs and record effective
  launch metadata in run runtime events.

## v1.5.3 - 2026-05-27

- Add Antigravity CLI `agy` as a first-class headless agent and default
  replacement for Gemini.
- Mark Gemini as legacy compatibility while keeping existing overrides working.
- Prefer verified stronger defaults for Claude (`opus`/`max`) and Hermes
  (`xai/grok-4.3`) while keeping Antigravity model/thinking fields at
  `cli-default` until the CLI exposes flags.
- Update project and embedded protocol rosters, docs, examples, and runtime
  configuration defaults for the Antigravity migration.

## v1.5.2 - 2026-05-26

- Add TUI planner action execution and focus-aware action controls.
- Refresh dashboard and live TUI layouts with height-aware compact modes,
  two-column normal views, semantic badges, and short-terminal tests.
- Embed the full default `COOPERATION.md` protocol for workspace initialization
  while preserving local-dir bootstrap transport.
- Record concrete project roster metadata and ignore machine-local headless
  agent launch config.
- Improve `parley version --all` project status probing timeouts and fallback
  behavior.

## v1.5.1 - 2026-05-25

- Complete Parley review cycles for the continuous-run TUI planner slice and
  version-awareness project sync.
- Unify planner and manifest next-action serialization through a shared
  `internal/runaction` type.
- Make continuation planning round-aware and remove hardcoded `codex` ownership
  from generated continuation commands.
- Add `parley version --dir DIR --all` project targeting, indented JSON output,
  and cleaner missing-installer fallback errors.
