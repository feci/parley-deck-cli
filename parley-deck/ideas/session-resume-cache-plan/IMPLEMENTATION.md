---
idea: session-resume-cache-plan
status: in-progress
implementer: codex
date: 2026-05-18
branch: feature/session-resume-cache-plan
design-pr: https://github.com/feci/parley-deck-cli/pull/25
---

## Current protocol state

The implementation proceeds incrementally under explicit user instruction from
2026-05-18. Round 1 has canonical artifacts from `codex` and `claude`.
`gemini` and `hermes` remain blocked by local CLI/runtime issues recorded in
`parley-deck/inbox/codex-to-all_session-resume-cache-plan_agent-blockers.md`.

Because another participant is available, `claude` is being used for slice
review and follow-up artifacts. The implementation keeps changes small enough
to review and release one slice at a time.

## Slice 1: local run manifest and sessions CLI

Status: implemented by `codex` on 2026-05-18; review round 1 fixes applied;
Claude re-review is blocked by local authentication (`401 Invalid
authentication credentials`).

Scope:

- Add a lightweight repository-local run manifest at
  `parley-deck/runs/<run-id>/run.json`. The manifest includes
  `schema_version`, run ID, workspace root, idea slug, mode, transport,
  status, participants, and timestamps.
- Keep `~/.parley-deck/sessions.json` as the global local index.
- Add `parley sessions list` for known local sessions.
- Add `parley sessions inspect <run-id>` for a detailed run/session view.
- Add tests for manifest creation and sessions CLI behavior.

Implemented files:

- `internal/runmanifest/manifest.go`
- `internal/runcontrol/runcontrol.go`
- `internal/app/app.go`
- `internal/runcontrol/runcontrol_test.go`
- `internal/app/app_test.go`

Verification:

- `GOCACHE=/Users/tomasfecko/AI_WORKSPACE/parley-deck/.gocache GOMODCACHE=/Users/tomasfecko/AI_WORKSPACE/parley-deck/.gomodcache go test ./...`
- `git diff --check`

## Review round 1 fixes

Reviewer: `claude`

Agreed fixes applied:

- Added a `status` field to `run.json`, defaulting to `running`, so schema
  version 1 can express lifecycle state before later terminal updates are
  implemented.
- Kept the `run.json` field names `idea_slug`, `workspace_root`, `created_at`,
  `updated_at`, `task`, and `mode`. These are intentional slice-1 contract
  fields because they match existing Go naming and make `sessions inspect`
  useful without reading every event.
- Updated `sessions list` to show `workspace`, `started`, and manifest-backed
  `status` when available.
- Added current-workspace/`--dir` fallback for `sessions inspect` when a run is
  not yet present in the global session index.
- Improved the legacy-run message for runs without `run.json`.
- Documented the read-only/single-active-process caveat in CLI help.
- Expanded tests for empty list, JSON list, legacy inspect without manifest,
  unindexed `--dir` fallback, missing run failure, and direct manifest
  round-trip/defaults.

Deferred from review round 1:

- Malformed `sessions.json` recovery remains out of scope for slice 1. The
  current store treats malformed JSON as an error, consistent with the existing
  cache reader.
- Terminal status updates are deferred until the per-run lifecycle slice.

Out of scope for this slice:

- Native agent session resume handles.
- Prompt input packs.
- Heartbeats and stale process recovery.
- TUI action menus for retry/resume.
- Cache rebuild and workspace rebind commands.

## Planned follow-up slices

1. Record per-agent attempts, prompt hashes, and capped logs in the run cache.
2. Add native resume capability probes and session handles for Claude, Codex,
   Gemini, and Hermes where supported.
3. Add cache rebuild/rebind commands for moved workspaces and wiped local state.
4. Expose resume/retry actions and per-agent activity in the TUI.
5. Add retention, pruning, and redaction controls for sensitive local state.
