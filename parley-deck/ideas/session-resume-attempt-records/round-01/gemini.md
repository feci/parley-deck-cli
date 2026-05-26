---
agent: gemini
idea: session-resume-attempt-records
round: 1
date: 2026-05-25
---

## Summary

I support the proposal to introduce per-agent attempt records. This is a crucial step for audibility and reliable replay across different agents. A structured, versioned schema is essential for interoperability, and keeping sensitive prompt data out of the committed workspace is the correct approach. My focus is on ensuring these records are standardized and sufficient for consistent cross-agent configuration and replay.

## Proposed approach

I propose a versioned JSON schema for the attempt records to ensure consistency. These records should live in the workspace under a new directory, for example `parley-deck/runs/<run-id>/attempts/<timestamp>-<agent-id>.json`.

The full, unredacted prompts should be cached locally in the user's home directory (e.g., `~/.parley-deck/prompt-cache/`) indexed by their SHA-256 hash. This allows for exact replay for debugging purposes without committing sensitive data.

A draft schema for each attempt record file:
```json
{
  "schemaVersion": "1.0",
  "runId": "<run-id>",
  "attemptId": "<timestamp-nanos>-<agent-id>",
  "agentId": "<agent-id>",
  "timestampUTC": "YYYY-MM-DDTHH:MM:SS.ssssssZ",
  "execution": {
    "status": "completed" | "failed" | "timed_out",
    "durationMs": 12345,
    "exitCode": 0,
    "processId": 54321
  },
  "config": {
    "model": "gemini-2.5-pro",
    "profile": "balanced",
    "cliArgs": ["--model", "gemini-2.5-pro", "--temperature", "0.7"],
    "headlessConfig": {
      "source": "parley-deck/meta/headless-agents.local.json"
    }
  },
  "prompt": {
    "sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "byteLength": 1234
  },
  "output": {
    "expectedArtifact": {
      "path": "parley-deck/ideas/foo/round-01/gemini.md",
      "validation": {
        "status": "valid" | "invalid" | "missing" | "not_checked"
      }
    },
    "diagnostics": {
      "stdoutTail": "...",
      "stderrTail": "..."
    }
  }
}
```

## Concerns / open questions

1.  **Redaction Policy:** We need a clearly defined and consistently implemented policy for redacting sensitive information from `cliArgs` and the `diagnostics` logs. Should this be a shared library function that all agent runners use?
2.  **Configuration Source:** The `config` block should unambiguously trace back to the source configuration (`headless-agents.local.json`, CLI flags, defaults). This is key for replayability. How do we represent the hierarchy of applied settings?
3.  **Log Sizing:** What is a reasonable maximum size for `stdoutTail` and `stderrTail`? I propose we start with 4KB each.
4.  **Atomic Writes:** How do we handle abrupt process termination? The runner/facilitator should probably write a "pending" record before invoking the agent, and then update it with the outcome. This ensures we have a record even if the agent process crashes without a trace.

## Risks

*   **Schema Adherence:** Without strong, centralized validation of the attempt records, we risk schema drift between agents, which would undermine the goal of interoperable tooling. The `parley` CLI should own validation.
*   **Replay Fidelity:** Replaying an attempt accurately depends on the environment (CLI versions, auth state) being similar. The record should capture more context, like the agent CLI's version (`<cli> --version`), to help diagnose replay failures.
*   **Incomplete Records:** If an agent crashes catastrophically, it cannot write its own attempt record. The parent runner process must be responsible for creating a minimal record indicating the launch attempt and failure.