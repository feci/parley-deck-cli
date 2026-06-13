---
agent: codex
idea: embedded-default-protocol-resync
review-round: 1
date: 2026-06-13
reviewed-commit: efe76d0
---

## Summary

Reviewed `efe76d0` against `FINAL.md`, `IMPLEMENTATION.md`, and `consensus.md`. The required §12 diff is empty, both files end with a newline, the embedded header and §2 tables match the current agreed bootstrap shape, and `defaultCooperationForInit()` remains the transport-swap-only helper. `go build ./...` and `go test ./internal/protocol/` pass when run with a writable temporary Go cache.

## Findings

### [MINOR] Guard normalizes embedded-only defaults without asserting them

What is wrong: `internal/protocol/drift_test.go:87` and `internal/protocol/drift_test.go:91` normalize any `Workspace` / `Created` line value and drop all §2 table body rows for both files. The live deck values and table rows are intentionally project-specific, but the embedded default's placeholder header values and empty table bodies are explicit D2/D3 bootstrap requirements. A future edit that adds an illustrative/default row to the embedded roster, or alters the created line while keeping the key prefix, can pass the drift guard; the init-output test only blocks the current parley-deck roster names.

Why it matters: This widens the blind spot called out in the implementation notes. An in-allowlist edit to embedded bootstrap behavior can slip through even though new projects must start with exact placeholders and empty quorum tables.

Concrete fix: Keep the five-zone live-vs-embedded comparison allowlist, but add pre-normalization assertions for the embedded default: exact D2 header lines (`Workspace`, `Created`, `Transport`, no `Protocol synced`) and zero body rows after each §2 separator. Use the same anchored table parser so missing or duplicated anchors still fail closed.

### [NIT] IMPLEMENTATION.md records a stale head commit

What is wrong: `parley-deck/ideas/embedded-default-protocol-resync/IMPLEMENTATION.md:8` says `head-commit: bc0af15`, but the reviewed branch head and requested reviewed commit are `efe76d0`.

Why it matters: The protocol artifact's audit trail points reviewers at an earlier/amended commit rather than the reviewed one.

Concrete fix: Update `head-commit` to `efe76d0` or its full hash after the final review-ready commit is settled.

## Open questions

None.

## Verdict

ACCEPT-WITH-FIXES
