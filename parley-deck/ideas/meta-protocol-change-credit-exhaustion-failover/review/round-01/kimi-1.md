---
agent: kimi-1
idea: meta-protocol-change-credit-exhaustion-failover
review-round: 1
date: 2026-10-29
reviewed-commit: 931550cb4c57f041f9a8a1d1337571cb271a873a
---

## Summary

The implementation is not ready for review consensus or release. Three MAJOR findings and one MINOR finding remain. Native synthetic checks and exact protocol reconstruction passed within their stated limits. The full Go suite failed, and AC-15 remains open.

CLI target: `931550cb4c57f041f9a8a1d1337571cb271a873a`, compared against `c2d2431`.
Skill target: `70d29c1aebd9a632f6f6239123fd2f52e3775a48`, compared against `f557c5e`.

Protocol attestation:

context_mode: full
source_sha256: b9284d7d7062ef49d59439a23117dd7434e494340b5b564a64651a12eb324709
packet_sha256: b9284d7d7062ef49d59439a23117dd7434e494340b5b564a64651a12eb324709
fallback_reason: none

I reviewed D1-D12 and every acceptance criterion. I made no provider call and induced no exhaustion. Verdicts use §15 provenance. Code locators and executed checks supply the `PRIMARY` evidence.

## Refutation attempts

### AC-01 — v2 creation and saved-policy compatibility

Attempt: trace all production policy constructors. Reserialize omitted-trigger and v1 records. Reject null, unknown, and duplicate triggers. Try widening without owner-bound roles.

Result: `CONFIRMED (PRIMARY)` for this scope. `internal/app/app.go:1952` and `internal/app/app.go:2512` select v2 for new ideas. `internal/quota/exhaustion_test.go` covers legacy bytes and malformed triggers. `internal/membership/succession_test.go` covers owner widening without kickoff rewrite.

Executed: `go test ./internal/quota ./internal/membership ./internal/protocol -count=1 -timeout 10m`. All three packages passed. I did not execute a hosted new-idea launch.

### AC-02 — native evidence and reset separation

Attempt: test all four exact Codex messages. Add detail, generic quota text, context errors, rollout budgets, assistant quotes, and tool quotes. Break JSON framing and terminal order. Vary Zcode resets and append mixed framing.

Result: `CONFIRMED (PRIMARY)` for synthetic fixtures. The recognizers rejected the tested near-misses. Short, absent, past, malformed, and contradictory resets did not negate valid Zcode exhaustion.

Executed: `go test ./internal/telemetry -run '^TestExhaustion' -count=1 -timeout 10m`. Result: `ok parley-deck-cli/internal/telemetry`. No captured native Codex exhaustion or Zcode native-positive claim follows.

### AC-03 — one child and durable suppression

Attempt: exhaust a synthetic Zcode child. Replay with a changed run and downgraded policy. Try a different launch phase. Inspect common admission for prior proof and validation receipts.

Result: `CONFIRMED (PRIMARY)` for synthetic lifecycle coverage. `internal/runner/exhaustion_test.go` observed one child and durable suppression. `internal/telemetry/exhaustion_admission.go` enforces the common gate before launch.

Executed: `go test ./internal/runner -run '^TestExhaustion' -count=1 -timeout 10m`. Result: `ok parley-deck-cli/internal/runner`. SDK-internal retries remain outside the verified boundary.

### AC-04 — terminal/validation recovery and stale output

Attempt: remove a validation receipt. Replay recovery with a run function that fails if called. Corrupt the receipt. Present an old valid implementation artifact to a failed new step.

Result: `CONFIRMED (PRIMARY)` for the tested recovery cases. Recovery launched no child. Corrupt proof blocked admission. The stale artifact did not mask failure. Fix-up freshness uses the same helper at `internal/runner/exhaustion.go:69`.

Scope limit: the tests remove or corrupt receipts after completion. They do not inject a fatal crash at every durable-write instruction.

### AC-05 — transient parity and protected-step compatibility

Attempt: trace a transient failure for a legacy fallback implementer and a global-default implementer. Compare the v2 path with the prior two-attempt and global-default fallthrough rules.

Result: `WRONG (PRIMARY)`. `internal/app/app.go:1965-1970` assigns either identity to `Assignment.Implementer`. `internal/quota/quota.go:143-145` treats that assignment as protected. `internal/app/preflight.go:964-972` then selects one attempt.

This contradicts the legacy comment at `internal/quota/quota.go:140`. It can also block global-default fallthrough. See the first finding.

### AC-06 — owner-order succession

Attempt: exercise ordinary, organizer, implementer, both, same-holder, simultaneous-loss, and original-order cases. Make the selected counterpart unusable. Request a skip or substitute.

Result: `CONFIRMED (PRIMARY)` for table-driven and history tests. `internal/quota/roles_test.go` applies the owner order and rejects an unusable selection. `internal/membership/succession_test.go` settles the complete loss set before selection.

The focused quota, membership, and protocol packages passed. No live role takeover was attempted.

### AC-07 — frozen facilitator exception and floor

Attempt: transfer implementation to the external frozen facilitator. Remove another usable non-facilitator. Try to count the facilitator toward the floor. Exhaust the external organizer for a role-only transition.

Result: `CONFIRMED (PRIMARY)` for reducer coverage. `internal/quota/roles_test.go` produced the sole membership addition and kept the fixed floor. The external organizer loss changed roles without membership deletion. Role validation rejects other additions.

### AC-08 — active execution authority

Attempt: build two handovers. Add a misleading filename. Change and remove predecessor bytes. Substitute a symlink. Replay a pending handover. Audit fixed `IMPLEMENTATION.md` consumers and the phase reader.

Result: `WRONG (PRIMARY)` for the complete current-reader scope. Resolver and pending-action tests passed. The fixed-path consumers mostly use `ResolveActiveExecution`. However, the phase reader has a confirmed error-handling regression.

Executed: `go test ./internal/driver -run '^TestRebuildDetailSurfacesReadErrors/unreadable_IMPLEMENTATION.md$' -count=1 -v`.

Output: `partial detail phase=round, want impl (the stat still saw the file)`. See the third finding.

### AC-09 — shared-draft integrity

Attempt: recover incomplete consensus, FINAL, and review-consensus drafts. Add a signoff, malformed identity, duplicate metadata, and an open fence. Keep a complete unsigned consensus draft unchanged.

Result: `CONFIRMED (PRIMARY)` for the structural fixtures. `internal/protocol/execution_test.go` preserved complete bytes and rejected the adversarial cases. `internal/membership/succession_test.go` archived one eligible draft byte-exactly and retained participant findings.

### AC-10 — retained obligations and independent review

Attempt: move a reviewer into implementation. Retain its earlier MAJOR finding. Vary snapshot models and Before/After membership. Inspect retained-veto and owner-ruling dispositions.

Result: `CONFIRMED (PRIMARY)` for synthetic succession and review-gate coverage. The focused membership tests passed. `internal/membership/review_gate.go` derives both reviewer sets with their own implementers and the frozen facilitator excluded.

Scope: an external-facilitator causal handover remains `UNVERIFIED` in a live child process. No fresh goal PASS exists yet.

### AC-11 — exact Codex support

Attempt: use text mode, an unknown version, a future version, and malformed effective argv. Compare source near-misses and exact messages. Inspect the ordinary smoke receipt.

Result: `CONFIRMED (PRIMARY)` for exact-version, explicit-JSON, source-backed support only. Local `codex --version` printed `codex-cli 0.165.0`. The recognizer tests passed. The smoke receipt records ordinary JSON behavior and no exhaustion.

No installed-binary equivalence or captured native positive is established. Default text, other versions, Kimi, and Claude remain diagnostic-only.

### AC-12 — private capture controls and security

Attempt: enable capture through environment and alternate homes. Scrub credentials, prompts, headers, requests, and unrelated payloads. Test symlinks, permissions, truncation, age, capacity, status, and purge.

Result: `CONFIRMED (PRIMARY)` for implementation behavior, with one test-portability finding. The scrubber writes nothing for uncertain grammar. The capture store supplies no recognizer or membership authority.

A real CLI check with a forced `0644` capture file printed `capture storage: unsafe capture file` and exited 1. The Go test fails under umask `0077`, because its fixture becomes `0600`. It passes under umask `0022`. See the fourth finding.

### AC-13 — exact protocol hunks and packet limits

Attempt: hash the signed JSON. Apply each of its 29 unique hunks independently to each base. Compare reconstructed bytes with all three copies and staged core 2.20.0. Render both phase-1 packets.

Result: `CONFIRMED (PRIMARY)`. The signed JSON SHA256 is `f140e81ada0a8fb08644109b889979e54945f829183c3d82b47c760eae67147c`. Every hunk occurred exactly once per base. All four reconstructions reported `exact_bytes True`.

The live and skill copies hash to `b9284d7d7062ef49d59439a23117dd7434e494340b5b564a64651a12eb324709`. The bundled CLI copy hashes to `43f2be0ef5fb03aa1a669697362446616925d647563f0bafa985d89e41c96994`. Staged core hashes to `6c437ed5fa0805bc8929b1aa25371664f9c4df1fc34c731a9df02bf02f2bd357`.

Executed: `go test ./internal/app -run '^TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail$' -count=1 -v`. It passed with 68,448 unflagged bytes and 69,880 meta bytes. `SKILL.md` is 19,989 bytes.

### AC-14 — versions, guidance, and full suites

Attempt: run the required Go and skill commands. Compare all version metadata. Review adapter, saved-policy, handover, capture, and account-alias guidance.

Result: `WRONG (PRIMARY)`. Executed: `go test ./... -count=1 -timeout 45m`. It failed three checks:

- `VERSION="1.55.0" internal version="1.56.0"`.
- `partial detail phase=round, want impl (the stat still saw the file)`.
- `unsafe status/purge accepted` under umask `0077`.

The targeted race run failed the same capture fixture and reported no data race. The other listed race packages passed.

Executed: `npm test` in the skill worktree. Result: 402 Node tests and 54 Python tests passed. Payload manifests passed. Package dry-run and native portable build remain implementer-reported checks; I did not rerun them.

### AC-15 — independent close and release

Attempt: inspect current review, goal, and release gates. Look for current-tree evidence, release hashes, channel evidence, runtime hashes, and separate channel review.

Result: `UNVERIFIED`. These findings remain open. No fresh independent goal PASS exists. No channel is released. Release and install verification belongs after fixes and fresh goal checking. I performed no release or install action.

## Findings

### [MAJOR] V2 changes transient retry and global-default protection

`internal/app/app.go:1965-1970` eagerly assigns the fallback or global-default implementer. `internal/membership/succession.go:85-90` also resolves the fallback before a transition. `internal/quota/quota.go:143-145` then protects that identity.

At readiness, `internal/app/preflight.go:964-972` gives the identity only one attempt. A transient failure can therefore block the batch. Legacy behavior gives two attempts and lets an unavailable global default fall through.

This violates D1, D4, and AC-05. It also contradicts the explicit legacy-protection comment at `internal/quota/quota.go:140`.

Concrete fix:

- Separate succession assignment from legacy protected-role status.
- Defer fallback and global-default assignment until exhaustion requires a transition.
- Keep one-attempt gates limited to roles protected under the prior rules.
- Add v2 readiness tests for transient fallback, global-default fallthrough, and exhaustion succession.

### [MAJOR] CLI release metadata still says 1.55.0

`VERSION:1` contains `1.55.0`, while `internal/app/version.go:3` contains `1.56.0`. The required full suite confirms the mismatch through `TestVersionFileMatchesBinaryVersion`.

This violates D12 and AC-14. Source archives and release metadata cannot consistently identify the intended CLI release.

Concrete fix: set `VERSION` to `1.56.0`. Rerun the version test and the complete Go suite at the fixed HEAD.

### [MAJOR] Execution read errors erase the implementation phase

`internal/driver/cursor.go:268-278` overwrites `implExists` with the result of `ExecutionStarted`. An unreadable implementation file returns `false` and an error. The later switch then selects `PhaseRound`, despite the successful stat.

The existing test returns `partial detail phase=round, want impl (the stat still saw the file)`. This regresses current-state error reporting and leaves AC-08 and AC-14 unsatisfied.

Concrete fix: preserve stat presence when `ExecutionStarted` returns an error. Apply its boolean only when its error is nil. Retain tests for early handover, unreadable files, and corrupt history.

### [MINOR] Capture permission test depends on the process umask

`internal/telemetry/capture_test.go:98` requests mode `0644` through `os.WriteFile`. Umask `0077` turns that request into `0600`. `CaptureStatus` then correctly accepts the safe file, so the test reports `unsafe status/purge accepted`.

The implementation rejected an actual `0644` file in my real CLI check. The defect is test portability, but the required full suite remains red until the fixture is corrected.

Concrete fix: call `os.Chmod(path, 0644)` after creation and check every setup error. Run this test under restrictive and permissive umasks.

## Open questions

- What are the full Go and targeted race results at the fixed HEAD?
- Can a fresh independent goal check pass every criterion after those fixes?
- Which release hashes, hosted assets, runtime installs, and separate channel checks will satisfy AC-15?
