---
agent: claude
idea: antigravity-agent-migration
review-round: 2
date: 2026-05-27
responding-to: [review/consensus]
---

## Summary

All five agreed fixes from review/consensus.md have been correctly applied.
No new findings. The implementation is ready to proceed.

## Verification of agreed fixes

### 1. `installedAgentIDs` scoped to exclude only legacy `gemini` — VERIFIED

`internal/app/app.go:1934-1942` — The function now checks only
`result.Found && result.ID != "gemini"`. The extra `LaunchMode == LaunchHeadless`
and `len(HeadlessArgs) > 0` conditions from my round-01 MAJOR #1 are removed.
Non-headless and ACP-capable agents are no longer silently dropped from CLI
default selection.

### 2. Hermes `--yolo` removed from built-in and shared defaults — VERIFIED

`internal/agents/discover.go:192` — HeadlessArgs are
`["--oneshot", "{prompt}", "--model", "xai/grok-4.3", "--accept-hooks"]` with no
`--yolo`. `ApprovalPolicy` at line 199 is `"accept-hooks"`, not
`"yolo,accept-hooks"`. `parley-deck/agents.toml:54` matches:
`approval_policy = "accept-hooks"`.

### 3. `agy` isolated-home limitation documented — VERIFIED

`docs/agent-runtime-configuration.md:31-34` states that `agy` exposes no stable
home-isolation environment variable and that users must provide their own
`isolated_home_env` mapping before enabling `isolate_home`. The built-in spec in
`discover.go:140-162` correctly omits `IsolateHome` and `IsolatedHomeEnv`.

### 4. `agy --print` argument-order constraint documented in code — VERIFIED

`internal/agents/discover.go:146` — Comment reads: `// Keep {prompt} immediately
after --print and last; agy treats --print as a value-taking flag.`

### 5. Legacy annotation on retained `gemini` host mapping — VERIFIED

Both `parley-deck/COOPERATION.md:94` and
`internal/protocol/defaults/COOPERATION.md:94` show
`| gemini | not mapped (legacy) |` in the host-handle table, consistent with
the `inactive legacy` role in the roster table at line 77.

## Findings

No findings.
