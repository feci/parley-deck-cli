---
agent: agy
idea: runner-hardening-kindly
round: 1
phase: review
date: 2026-06-12
---

## Summary

Reviewed the implementation of `runner-hardening-kindly` (D1-D12) under the designated lenses: failure taxonomy & hints (D5/D6), heartbeat view contract (D4), finished-with-exit rendering (D7), and consult command UX (D10). The core supervision architecture, watchdog kills, and artifact-wins rules are cleanly implemented and behave calmly. A consensus divergence was identified where the `session_id` field is missing from the consult frontmatter.

## Findings

### [MAJOR] Missing session_id in consult frontmatter
**File:** [internal/app/consult.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/app/consult.go#L110-L124)
**Why:** Consensus D10 specifies that the consult artifact frontmatter must include a `session_id` field. The current implementation omits this field entirely when writing the markdown frontmatter template.
**Suggested fix:** Include the `session_id` in the frontmatter template using the format string. Because the runner executes the consult using a session ID composed of `"consult:"+agentID`, the frontmatter can write `session_id: consult:<agentID>`:
```diff
--- internal/app/consult.go
+++ internal/app/consult.go
@@ -118,6 +118,7 @@
 timeout_ms: %d
 exit_code: %d
 agent_exit: %d
+session_id: consult:%s
 stdout_log: %s
 stderr_log: %s
 quorum: false
@@ -126,3 +127,3 @@
 	`, agentID, valueOr(agent.Model, "cli-default"), created.Format(time.RFC3339), slug, question,
-		root, timeout.Milliseconds(), exitCode, res.AgentExit,
+		root, timeout.Milliseconds(), exitCode, res.AgentExit, agentID,
 		relOrSelf(root, stdoutLog), relOrSelf(root, stderrLog))
```

### [NIT] Auth recovery hint phrasing alignment
**File:** [internal/runner/failclass.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runner/failclass.go#L27)
**Why:** The implemented recovery hint for the `auth` failure class (`"Run the agent CLI's auth command (e.g. 'claude login', 'hermes auth') to refresh credentials."`) varies slightly from agy's round-01 table (`"Run client auth command (e.g., 'claude login') to refresh credentials."`). This change is an improvement that incorporates `hermes` commands and updates wording, so no code change is requested.
**Suggested fix:** None; accept as an improvement.

### [NIT] Phrasing updates in recovery hints
**File:** [internal/runner/failclass.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runner/failclass.go#L41-L52)
**Why:** Phrasing for several hints (e.g., `invalid-request`, `stalled`, `unknown`) has been adjusted relative to agy's round-01 table (adding "the" where appropriate, generalizing `stalled` to refer to "the stall window" rather than a hardcoded "30m", and updating `unknown` to refer to stdout/stderr log files without hardcoding `.local/logs`). These adjustments are logical improvements.
**Suggested fix:** None; accept as improvements.

## Dispositions

- **Finding/disposition: failure-class names use hyphens (rate-limit, model-not-found) except watchdog classes which use underscores (no_first_output, stalled).**  
  **Prior disposition:** Accepted trade-off (consensus D5).  
  **Rationale:** Watchdog class names must equal their event type suffixes (agent.no_first_output).  
  **Concurrence:** I concur. Matching the watchdog class names to the underlying event type suffixes simplifies log parsing and avoids introducing mapping noise, and the hyphenated-vs-underscored inconsistency is a minor stylistic cost.

## Verdict

ACCEPT-WITH-FIXES
