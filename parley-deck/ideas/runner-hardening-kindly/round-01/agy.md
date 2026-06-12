---
agent: agy
idea: runner-hardening-kindly
round: 1
date: 2026-06-12
---

# Runner Hardening (Kindly) - Operator UX Analysis

This analysis addresses the UX and operator ergonomics aspect of the kindly integration.

## P2: Failure-Class Taxonomy and Recovery Hints

We define the 12 failure classes, mapping each to a precise recovery hint. Stderr parsing will classify errors and attach `failure_class` + `recovery_hint` to the `agent.failed` event payload.

| Failure Class | Recovery Hint |
| :--- | :--- |
| `rate-limit` | `Wait for reset or switch provider keys/endpoints.` |
| `overloaded` | `Retry in a few minutes or choose a less busy model.` |
| `auth` | `Run client auth command (e.g., 'claude login') to refresh credentials.` |
| `billing` | `Check your API account balance and credit card status.` |
| `invalid-request` | `Verify prompt structure and system constraints in config.` |
| `model-not-found` | `Check model spelling and access permissions in your API settings.` |
| `context-window` | `Reduce prompt size or prune file attachments/logs from scope.` |
| `sandbox` | `Adjust local sandbox configuration or run with lower restriction.` |
| `budget` | `Increase session budget limit (e.g., raise spend caps in settings).` |
| `no_first_output` | `Verify agent executable is not blocking or waiting for stdin.` |
| `stalled` | `Check process tree; agent did not emit output for 30m.` |
| `unknown` | `Check agent stderr/stdout log files in .local/logs for details.` |

### UI Integration Points

1. **Narrator Line Format** (in `friendlyEventText`):
   When `agent.failed` occurs, if `failure_class` and `recovery_hint` are present, format as:
   `── [TIMESTAMP] [agent] failed: [failure_class] — [recovery_hint] ──`
   *Example:* `── 17:08:15 claude failed: auth — Run client auth command (e.g., 'claude login') to refresh credentials. ──`

2. **Agent Status Header** (in TUI status panel):
   Show the class and the hint dynamically on failure:
   `Status: FAILED | Class: [failure_class] | Hint: [recovery_hint]`

3. **Parley Status Output** (via `parley status`):
   When an agent fails, print a detailed block:
   ```
   Agent: [agent]
   Status: Failed
   Failure Class: [failure_class]
   Recovery Hint: [recovery_hint]
   ```

## P1: Watchdog Kills and Retries UX

During watchdog events, the narrator should keep the operator informed of background interventions using clear, active-voice state updates.

- **No First Output (Grace Period Exceeded):**
  `── [TIMESTAMP] [agent] no first output (grace window elapsed) — killing process tree ──`
  `── [TIMESTAMP] [agent] retrying (attempt 2/2) ──`
  *If retry fails:*
  `── [TIMESTAMP] [agent] failed: no_first_output — no output within grace window ──`

- **Stalled Execution (Silence limit reached):**
  `── [TIMESTAMP] [agent] stalled (no new output for 30m) — killing process tree ──`
  `── [TIMESTAMP] [agent] failed: stalled — killed due to inactivity ──`

- **Heartbeat Progress Updates:**
  To maintain a lively TUI without spamming logs, heartbeats are rendered inline/overwritten on the status line:
  `── [TIMESTAMP] [agent] running (elapsed: [elapsed], bytes: [bytes], last event: [event_type]) ──`

## P3: Nonzero Exit with Valid Artifact UX

When the process exits nonzero but the artifact validates successfully:
1. Record state as `SUCCESS` via `agent.finished`.
2. Save the exit code in `agent_exit`.
3. Render in the narrator with a clear indication that it completed successfully despite the exit code:
   `── [TIMESTAMP] [agent] finished (exit code: [code], artifact verified) ──`
4. In the status header:
   `State: Finished (Exit Code: [code], Artifact: OK)`
This avoids false alarm UI indicators (red colors/error banners).

## P8: Consult Command UX

### Invocation Ergonomics
Users can run questions directly or pipe them:
- `parley consult <agent> "<question>"`
- `cat query.txt | parley consult <agent>`

### Waiting UI
While executing, print live heartbeats directly to stderr so the stdout remains clean for redirection:
```
Consulting [agent]: "[question]"
[17:08:15] Starting read-only sandbox session...
[17:08:20] [agent] thinking (elapsed: 5s, events: 4, bytes: 124B)
```

### Artifact Schema
Saved to `parley-deck/consults/YYYYMMDD-HHMMSS-<agent>-<question-slug>.md` with frontmatter:
```yaml
---
consult: true
agent: <agent>
timestamp: 2026-06-12T17:08:15+02:00
question: "..."
model: <model_name>
session_id: <session_uuid>
exit_code: 0
---
```

### Discovery of Past Consults
1. **Durable Index:** Every consult appends to `parley-deck/consults/index.jsonl`.
2. **List Command:** `parley consults list` prints a formatted table:
   ```
   DATE                 AGENT    QUESTION                         PATH
   2026-06-12 17:08:15  claude   How to clone on virtio-fs?       parley-deck/consults/20260612-170815-claude-how-to-clone.md
   ```

## Position

The proposed UX makes runner interventions, failure modes, and advisory consults clear and actionable without TUI clutter. I support proceeding with these specifications.
