---
agent: hermes-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-11
scope: owned files only — internal/app/preflight.go (no source edits by me), preflight_liveness.go, preflight_liveness_test.go; plus this new OWN artifact
status: edits completed; tests NOT run by me (facilitator runs independently); unverified
---

## Actual changes completed (final bounded correction attempt)

- Created this note at the required absolute path (FIRST, before any source edit) and updated after edits.
- preflight_liveness.go: replaced incorrect scanObject token walk with correct `dec.More()` approach — read key token (string), read value token directly (NO `:` or `,` delimiter tokens — Token() never returns them), recurse on nested `{`/`[`, enforce duplicate/case-aliased semantic keys (`isSemanticKey`) and present-null (`valTok == nil`) correctly.
- preflight_liveness.go: added empty-Write-at-cap guard (`remaining == 0`) so boundedWriter does not truncate when at cap; kept observedBytes independent.
- preflight_liveness.go: timeout-after-output uses actual byte presence (`stdout == "" && stderr == ""`) not TrimSpace; comment only (code already used direct comparison).
- preflight_liveness.go: scrubSecrets augmented with JSON credential-field redaction regex (`"token":...`, etc.) in addition to existing bearer/labeled/token patterns — secrets including JSON credential fields redacted.
- preflight_liveness_test.go: appended `false, ""` to the four 4-argument `classifyReadiness` calls at lines 378,382,388,394 (now 6-argument, matching signature `stdout, stderr, exitCode, timedOut, truncated, truncationReason`).
- preflight_liveness_test.go: kept positive fixtures for concrete supported schemas (`{"message":"PONG"}`, `{"text":"PONG"}`, `{"result":"PONG"}`, `{"role":"assistant","content":"PONG"}`, wrapper schemas); added explicit `{"type":"result","subtype":"success","is_error":false,"result":"PONG"}` positive child fixture; kept `{"content":"PONG"}` as malformed (no provenance); kept `{"type":"error",...}` and provider-error expectations unchanged (not fabricated to malformed); kept exact plain PONG (`PONG\n`) supported.

## Unchanged / out of scope (honest record)

- `preflight.go`: NO edits by me (only pre-existing `undefined: rosterEntry` at line 193/205; not my defect).
- No `go test` executed by me; no fabrication of test results; no other files edited.
- The 4-argument fix is source-only; facilitator will independently compile/run the corrected source.
