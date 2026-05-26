---
agent: claude
idea: session-resume-attempt-records
round: 1
date: 2026-05-25
---

## Summary

Slice 2 introduces per-agent attempt records that make each invocation auditable
and replayable without leaking sensitive prompt content into committed artifacts.
The core tension is between recording enough metadata for meaningful replay and
diagnostics while keeping full prompts, session handles, and potentially
identifying content strictly in local cache. This review focuses on schema
design, privacy boundaries, and the redaction contract.

## Proposed approach

### Attempt record schema

Each attempt should live at a predictable path under the run directory:

```
parley-deck/runs/<run-id>/attempts/<agent-id>-<attempt-number>.json
```

Minimal schema (v1):

```jsonc
{
  "schema_version": 1,
  "agent_id": "claude",
  "attempt": 1,
  "timestamp_start": "2026-05-25T10:00:00Z",
  "timestamp_end": "2026-05-25T10:02:31Z",
  "duration_ms": 151000,

  "launch": {
    "model": "claude-opus-4-7",
    "profile": "default",
    "effort": null,
    "cli_args_redacted": ["--print", "--model", "***", "--prompt-file", "***"]
  },

  "prompt": {
    "sha256": "ab12cd34...",
    "byte_length": 4820
  },

  "expected_artifact": "parley-deck/ideas/session-resume-attempt-records/round-01/claude.md",
  "artifact_validated": true,
  "validation_error": null,

  "exit_status": 0,

  "diagnostics": {
    "stderr_head": "",
    "stderr_tail": "",
    "truncated": false
  }
}
```

### Privacy boundary

Two strict rules:

1. **Committed fields are non-invertible or public.** Prompt SHA-256 hashes
   cannot recover the prompt. CLI args use positional redaction (the flag name is
   visible; the value is `***`). Diagnostic logs are bounded and redacted.
2. **Local cache holds everything sensitive.** Full rendered prompts, raw
   stderr/stdout, native session handles, and unredacted launch args go to
   `~/.parley-deck/cache/runs/<run-id>/attempts/<agent-id>-<attempt>.json`.

### Redaction policy for diagnostics

- `stderr_head` and `stderr_tail`: cap at 512 bytes each.
- Apply a simple pattern-based scrub before writing: strip tokens matching
  `sk-...`, `ghp_...`, `Bearer ...`, absolute home-directory paths, and
  environment variable expansions.
- If any scrub fires, set a `redacted: true` flag so inspectors know the output
  was altered.

### Schema versioning

The `schema_version` field is mandatory from v1 per the FINAL invariants. The
reader must reject or warn on unknown versions rather than silently parsing.

## Concerns / open questions

1. **Attempt numbering vs. retry naming.** The FINAL defers "exact retry naming."
   A simple monotonic integer (`attempt: 1, 2, 3`) is sufficient for slice 2, but
   should we distinguish "scheduled retry" from "manual re-run"? A `trigger`
   field (`auto` | `manual` | `tui`) would help without complicating the schema.

2. **Hash algorithm agility.** SHA-256 is fine today. Should the schema carry an
   `algorithm` field (`"sha256"`) so a future migration doesn't require a schema
   version bump? Low cost to add now.

3. **Diagnostic redaction completeness.** Pattern-based scrubbing is best-effort.
   A missed pattern leaks secrets into committed history. An alternative: never
   commit raw diagnostics at all; keep them local-only and commit only a boolean
   `has_diagnostics: true` pointer. This is safer but less useful for
   `sessions inspect` without local cache.

4. **Artifact validation timing.** Should validation run immediately after the
   agent exits, or allow a grace window for agents that write asynchronously
   (e.g., a background file sync)? Immediate validation is simpler and matches
   current behavior; document it as the contract.

5. **Multi-artifact agents.** The schema assumes one `expected_artifact` per
   attempt. If a future slice needs multiple artifacts per invocation, this
   should be an array. Starting as a single string is fine if the migration path
   is documented.

## Risks

- **Secret leakage via diagnostics.** Even with redaction, committed stderr
  excerpts are the highest-risk field. If any participant's CLI emits secrets on
  stderr (API keys in error messages, tokens in debug output), the scrub list
  must cover it or the data leaks into git history permanently. Mitigation:
  default diagnostics to local-only in v1; let users opt in to committed
  diagnostics via a config flag.

- **Hash collision is not a practical risk** for SHA-256 at this scale, but
  prompt hashes still act as fingerprints. Two runs with identical hashes reveal
  identical prompts. If prompt content is confidential between participants, hash
  equality leaks that two runs used the same input. Acceptable for internal
  tooling; worth noting for public repos.

- **Schema migration burden.** Every field added now becomes a compatibility
  contract. Keep v1 minimal; use the deferred-details list for fields that don't
  have a clear consumer yet (locking, heartbeats, capability tiers belong in
  slice 3+).

- **Path length on Windows.** Nested paths like
  `parley-deck/runs/<uuid>/attempts/<agent>-<n>.json` can approach the 260-char
  limit. Use short run IDs (8-char prefix) and keep agent IDs concise.
