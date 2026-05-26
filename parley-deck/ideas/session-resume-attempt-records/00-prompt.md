---
idea: session-resume-attempt-records
author: codex
created: 2026-05-25
participants: [codex, claude, gemini, hermes]
roles:
  codex: Go implementation and runner integration
  claude: schema and privacy review
  gemini: cross-agent replay and config interoperability
  hermes: long-running process diagnostics and recovery
status: final
---

## Problem / idea

Design the next session resume roadmap slice from
`session-resume-cache-plan/FINAL.md`: per-agent attempt records and prompt
hashes.

Slice 2 should make each agent invocation auditable and replayable without
committing sensitive prompt bodies. It should record launch metadata, prompt
hashes, bounded diagnostics, expected artifact paths, validation results, and
minimal run status updates needed for coherent `parley sessions inspect`
output.

## Constraints

- Follow `session-resume-cache-plan/FINAL.md`.
- Full rendered prompts and native session handles must stay in local cache,
  not committed workspace artifacts by default.
- Workspace artifacts may include hashes, redacted launch configuration,
  bounded diagnostics, expected artifact paths, validation results, and schema
  versions.
- Keep the first implementation slice small enough to review.
- Preserve canonical Parley ownership rules: attempt records describe
  invocations; they do not replace participant-owned protocol artifacts.

## Non-goals

- Do not implement native resume handles in this slice.
- Do not implement cache rebuild/rebind in this slice.
- Do not add TUI resume/retry actions in this slice.
