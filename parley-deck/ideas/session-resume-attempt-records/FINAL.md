---
idea: session-resume-attempt-records
status: final
author: codex
consensus-date: 2026-05-25
participants: [codex, claude, gemini, hermes]
---

## Final plan / specification

Implement slice 2 of the session resume roadmap: versioned per-agent attempt
records and prompt hashes.

Attempt records are written by the parent runner/facilitator process, not by
participant agents. They describe an invocation and its expected artifact; they
do not replace participant-owned protocol artifacts.

### Workspace attempt records

Add a small internal package for schema version 1 attempt records. The exact
filename can be settled during implementation, but records should live under
the workspace run directory, for example:

```text
parley-deck/runs/<run-id>/attempts/<step-id>-<agent-id>-attempt-<n>.json
```

Each workspace record should include:

- schema version;
- run ID, idea slug, step ID, agent ID, and attempt number;
- expected artifact path;
- prompt hash algorithm, prompt SHA-256, and prompt byte length;
- selected model/profile/effort when known;
- redacted launch command metadata;
- start/end timestamps, duration, exit code, and timeout/failure status;
- bounded redacted diagnostics suitable for committed artifacts;
- validation status: missing, exists, invalid, or valid.

### Local-only state

Full rendered prompts, raw stdout/stderr, and native session handles remain
local-only under `~/.parley-deck` or another explicit local Parley home. They
must not be committed by default.

### Write behavior

Attempt records should be written atomically:

1. create a pending/start record before invoking the agent when possible;
2. update the record with process outcome and artifact validation afterward;
3. preserve a minimal failed record if the agent crashes, times out, or writes
   no artifact.

### Inspect behavior

Extend `parley sessions inspect` to show attempt records when present. Keep the
display concise and safe: status, agent, step, expected artifact, validation
result, duration, and redacted diagnostic summary.

### Out of scope

- Native resume handles.
- Heartbeats and stale-process recovery.
- Cache rebuild/rebind.
- TUI retry/resume actions.
- Full lifecycle/terminal status semantics beyond the minimal status needed for
  coherent inspect output.

## References

- Consensus: ./consensus.md
- Rounds: ./round-01/
- Parent roadmap: ../session-resume-cache-plan/FINAL.md
