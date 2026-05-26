---
idea: session-resume-attempt-records
drafted-by: codex
date: 2026-05-25
---

## Agreed decisions

- Add versioned per-agent attempt records written by the parent
  runner/facilitator, not by participant agents.
- Workspace attempt records should be committed-safe: prompt hash and byte
  length, redacted launch metadata, expected artifact path, validation result,
  process timing/outcome, and bounded redacted diagnostics.
- Full rendered prompts, raw logs, and native session handles remain local-only
  state under `~/.parley-deck` or another explicit local Parley home.
- Attempt records must be written atomically, with a pending/start record before
  invocation and an outcome update after invocation when possible.
- The first schema must include a version field and a hash algorithm field.
- Validation can be shallow in v1: expected artifact missing, exists, invalid,
  or valid.

## Agreed trade-offs

- Bounded diagnostics are useful for `sessions inspect`, but they are the
  highest privacy risk. The implementation should use conservative redaction
  and small defaults, and it may keep richer diagnostics local-only.
- Retry naming can use deterministic attempt numbers. Exact path names may be
  settled during implementation as long as they are stable and documented.
- Heartbeats and native resume handles remain out of scope for this slice.

## Open items deferred to implementation

- Exact workspace path and filename scheme for attempt records.
- Exact diagnostic byte limits.
- Exact redaction patterns and tests.
- Whether CLI version capture belongs in v1 or the native-resume capability
  slice.

## Signoffs

<!-- Each agent appends its own signoff block. Do not edit others' blocks. -->

### Signoff: codex - 2026-05-25
Status: ACCEPT
Notes: Accept. This is small enough for a focused implementation slice.

### Signoff: gemini - 2026-05-25
Status: ACCEPT
Notes: Agreed, a focused approach to versioned per-agent attempt records is a good first step.

### Signoff: hermes - 2026-05-25
Status: ACCEPT
Notes: Accept. Small focused slice with clear boundaries.

### Signoff: claude - 2026-05-25
Status: ACCEPT
Notes: The agreed decisions are well-scoped and the privacy-first framing (conservative redaction, local-only for raw logs/handles) is the right default. Atomic pending→outcome writes and a versioned schema with hash algorithm field give a solid foundation. Open items are appropriately deferred without blocking the slice.
