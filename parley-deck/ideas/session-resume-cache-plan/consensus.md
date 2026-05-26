---
idea: session-resume-cache-plan
drafted-by: codex
date: 2026-05-25
---

## Agreed decisions

- The recovered design trail is valid: `codex`, `claude`, `gemini`, and
  `hermes` now each have canonical round-01, round-02, and round-03 artifacts.
- Slice 1 is treated as complete and should not be reopened:
  repository-local `parley-deck/runs/<run-id>/run.json`, read-only
  `parley sessions list`, and read-only `parley sessions inspect`.
- Parley promises workflow-level recovery and per-step retry. Native model
  conversation continuity is best-effort, reported per agent, and never
  required for correctness.
- Full rendered prompts and native session handles are sensitive local state.
  Store them only under local cache such as `~/.parley-deck` or another
  explicitly local Parley home. Committed workspace artifacts should contain
  hashes, redacted launch metadata, bounded diagnostics, expected artifact
  paths, and validation results.
- `headless-agents.local.json` is local, uncommitted machine configuration. It
  may carry explicit capability metadata, including resume capability tiers, but
  it is not canonical project state.
- Do not rely on parsing private external CLI cache formats. Prefer explicit
  launch configuration, captured handles from invocations, version probes, and
  smoke tests.
- New durable schemas must include a `schema_version` or equivalent version
  field from their first slice.

## Agreed trade-offs

- JSON files remain acceptable for the next slices if the code keeps a clean
  access boundary for a future SQLite migration.
- Minimal run status updates may be folded into slice 2 when they make attempt
  records and inspect output coherent. A separate status-only slice is not
  required unless implementation proves it is the smaller, safer path.
- Capability discovery should be explicit and conservative. Help-text parsing
  may inform diagnostics, but it should not become the only contract for resume
  behavior.
- Workspace records should stay portable and reviewable. Local cache can be
  richer, but cache loss must degrade to prompt replay/re-invocation rather
  than corrupting canonical artifacts.

## Open items deferred to implementation

- Exact attempt record path and retry naming, e.g. whether retries use
  `<round>-<agent>-attempt-N.json` or timestamped IDs.
- Exact bounded-log limits and redaction policy for stdout/stderr diagnostics.
- File locking and stale-heartbeat behavior for concurrent TUI/CLI access.
- Cache rebuild/rebind command names and handling of moved workspaces.
- Whether any workspace-local ignored cache is useful in addition to
  `~/.parley-deck`; default should remain outside committed artifacts.

## Roadmap

1. Slice 1 - complete: `run.json`, `parley sessions list`, and
   `parley sessions inspect`.
2. Slice 2 - per-agent attempts and prompt hashes: record launch configuration,
   prompt hash, bounded diagnostics, expected artifact path, validation result,
   and minimal status updates needed for coherent inspect output.
3. Slice 3 - native resume handles and capability config: capture native
   session IDs only when the launched CLI exposes or accepts a stable handle;
   add configured resume capability tiers for Claude, Codex, Gemini, and
   Hermes.
4. Slice 4 - cache rebuild/rebind: rebuild the local index from workspace run
   manifests and repair moved workspace paths.
5. Slice 5 - TUI resume/retry/manual handoff: expose user-facing recovery
   actions after the lower-level state is reliable.

## Signoffs

<!-- Each agent appends its own signoff block. Do not edit others' blocks. -->

### Signoff: codex - 2026-05-25
Status: ACCEPT
Notes: Accept. The consensus reflects the recovered multi-agent trail and
keeps sensitive prompt/native-session data out of committed artifacts by
default.

### Signoff: claude - 2026-05-25
Status: ACCEPT
Notes: Consensus is sound. The local-cache-only policy for sensitive state, workflow-level recovery over native continuity, and incremental slice roadmap all align with the design trail. No objections.

### Signoff: gemini - 2026-05-25
Status: ACCEPT
Notes: The plan is sound. I agree with the proposed approach.

### Signoff: hermes - 2026-05-25
Status: ACCEPT
Notes: Looks solid.
