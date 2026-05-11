---
agent: hermes
idea: agent-runtime-config
review-round: 1
date: 2026-05-11
reviewed-commit: d38e7f738471caa61042ca087c445b4df5bd2606
---

## Summary
Reviewed the implementation at commit d38e7f7 (post 1e1c757 implementation). Changes add cmd/parley/main.go, internal/config/runtime.go + tests, extend agents and app/runner layers for resolved runtime config with TOML layering and source tracking, update Codex defaults, add agents list/verify commands, and thread runtime values into runs. Matches FINAL.md scope closely with noted deviations in IMPLEMENTATION.md.

## Findings
### [CRITICAL] None
### [MAJOR] None
### [MINOR] agents.toml missing some fields present in runtime model
In parley-deck/agents.toml (lines 1-44), several fields from agentOverride struct in internal/config/runtime.go (e.g. speed, telemetry, version_args) are absent from the default TOML. This may lead to incomplete examples for users configuring new agents. No tests assert full field coverage from TOML.

### [NIT] Hardcoded timeout in verify smoke test
internal/app/app_test.go:152 uses a fixed 30s timeout for Codex Git smoke without using the resolved timeout_ms from config or a constant. Minor maintainability issue if defaults change.

## Open questions
- Does the placeholder expansion in runtime.go handle nested or cross-agent references safely, or should it be limited to ${HOME} style only?
- Should `parley agents verify --full` output include the exact resolved command line for reproducibility, beyond the current matrix?