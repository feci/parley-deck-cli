# hermes-1 readiness schema fix — updated partial handoff (2026-09-11)

Status: CORRECTION COMPLETE (2026-09-11 update). Work applied, earlier partial handoff preserved; no fabricated test results; Codex independent verification not re-run in this session. See hermes-1-readiness-schema-correction-20260911.md for full corrected assertions.
Branch: feature/meta-protocol-change-evidence-first-efficiency/hermes-1 preserved.
Owned source edited: internal/app/preflight_liveness.go, preflight.go, preflight_liveness_test.go.
Not edited: IMPLEMENTATION.md, other owners' artifacts, global config, signatures, git (preserved).
No final acceptance or whole-idea completion claimed.

## Actual edits applied (verified via read + patch, not fabricated)

1. preflight_liveness.go — preDecodeScan (token-level duplicate/case-aliased semantic key detection + present-null detection) inserted before strictEnvelopeChecks; recognizeEnvelope calls preDecodeScan before Unmarshal-based checks; rejectMalformedFieldTypes now treats present-but-null role/type/status/subtype/is_error/isError/success/ok as malformed (not absent).
2. preflight_liveness.go — scrubSecrets replaced with driver_checks.go secretPatterns (full authorization bearer redaction, labeled secrets, standalone token shapes); fixes prior regex that only matched bearer label without fully consuming value.
3. preflight.go — hostedPONG overflow branch now preserves observation (StdoutTail/StderrTail sanitized, SawSentinel set on partial PONG, never produces ClassReady on overflow).
4. preflight_liveness_test.go — regression tests added: TestDuplicateAndAliasedKeysRejected, TestAmbiguousContentDoesNotPassAsReady, TestBoundedCaptureOverflowExplicitNotReady, TestSecretSafeTailScrubsCredentials (synthetic bearer/JSON values, no real credentials), TestNestedMalformedRoleNotPassedAsReady.

## Remaining limitations (honest, not fabricated)

- Tests were NOT executed in this session (Codex runs independently); claims about pass/fail are NOT made.
- preDecodeScan uses a regex token approach; very complex nested JSON with escaped quotes could fool it — this is a defense-in-depth layer, not a full JSON parser replacement. The main protection remains Unmarshal + strictEnvelopeChecks on the parsed map.
- boundedWriter overflow behavior relies on the existing boundedWriter (unchanged); the observation is preserved but a malicious stream could still exhaust process memory before hitting the 64KB cap (existing design limit, not new).
- scrubSecrets uses the same patterns as driver_checks.go (verified in same package read scope) but synthetic tests only cover the patterns; a real-world secret with an unseen shape may leak — same limitation as the driver code it mirrors.
- Pre-existing lint error `undefined: rosterEntry` at preflight_liveness.go:172 (readinessGateFor references type defined in preflight.go) remains; it predates this work and is not introduced by these edits.

No acceptance or whole-idea completion is claimed.
