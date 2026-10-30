---
agent: zcode-1
idea: meta-protocol-change-credit-exhaustion-failover
review-round: 1
date: 2026-10-29
reviewed-commit: 931550cb4c57f041f9a8a1d1337571cb271a873a
---

## Summary

I reviewed the full CLI diff `c2d2431..931550c` (79 files) and the skill diff `f557c5e..70d29c1` against FINAL D1–D12 and AC-01 through AC-15. Skill commit: `70d29c1aebd9a632f6f6239123fd2f52e3775a48`. The implementation is healthy. The succession order, active execution resolver, draft recovery, retained findings, and review gates match the owner ruling and the frozen decisions.

I reconstructed the exact hunk application myself. All 29 signed hunks (`f140e81ada0a8fb08644109b889979e54945f829183c3d82b47c760eae67147c`) applied to the four before-bytes reproduce the four current files exactly. I ran the five focused Go packages and both hard packet tests myself. One capture subtest fails under umask 077; see the MINOR finding.

Protocol attestation: context_mode full; source_sha256 and packet_sha256 `b9284d7d7062ef49d59439a23117dd7434e494340b5b564a64651a12eb324709`; fallback_reason none. I hashed the live `parley-deck/COOPERATION.md` myself and it matches.

The full Go suite, release, and install checks stay open. AC-14 and AC-15 are not complete yet. I did not treat any pending check as passed.

## Refutation attempts

I assumed each criterion was broken and tried to break it. I list the attempt and the result for each criterion.

### AC-01 — v2 freeze, v1 preservation, owner-only widening

Attempt: I searched for a path that widens a saved policy without owner authority. I read `Policy.UnmarshalJSON` (`internal/quota/quota.go:29`), `validateRevision` (`internal/quota/revision.go:52`), and both `runTask` policy constructors (`internal/app/app.go`). I checked legacy kickoff bytes through `TestRoleFieldsDoNotChangeLegacySerialization`.
Result: not broken. The v1 constructor stays for saved policies. A v2 kickoff without a frozen role order fails validation. A widening revision without an owner-quoted `role_snapshot` fails. `TestOwnerWideningFreezesRolesWithoutKickoffRewrite` proves the old kickoff bytes stay identical.

### AC-02 — native framing only, separate reset, valid output wins

Attempt: I fed the recognizer the source near-misses and mixed terminal sequences in `TestExhaustionCodexExactSourceContract`. I traced the reset paths in `TestExhaustionZcodeResetCannotNegateClass`. I checked valid-output precedence in `RunParticipantStep` (`internal/runner/dropout.go`) and the in-flight success window in `ExhaustionAdmission` (`internal/telemetry/exhaustion_admission.go:130`).
Result: not broken. Only reviewed native framing yields exhaustion. Reset uncertainty never negates the class; the reset hint comes from the unchanged legacy classifier only. Valid output wins. A success counts only when it started before and completed after the failure, so a later fresh launch cannot win. Failed execution is required (`ValidateExhaustionEvidence`, `internal/quota/exhaustion.go:26`).

### AC-03 — one child, durable suppression

Attempt: I ran `TestExhaustionOneChildAndDurableAdmission`. I traced `beginLaunch` (`internal/runner/telemetry.go`) and checked suppression on a changed run and a consult phase. I checked the blocked-batch path through `SupplementUsability`.
Result: not broken. Exactly one child ran. Replay launched nothing. A changed run, prompt, and phase stayed suppressed. A stale valid artifact did not restore a suppressed identity's floor seat. Admission reads the durable invocation ledger, so it survives restarts and policy changes.

### AC-04 — recovery admits no new child, stale artifacts do not mask

Attempt: I ran `TestExhaustionTerminalValidationRecoveryNoChild` and `TestExhaustionStaleImplementationDoesNotWin`. I read `freshExecutionArtifact` (`internal/runner/exhaustion.go:69`) and the receipt checks in `ExhaustionAdmission`.
Result: not broken. A removed receipt forces recovery without a child. A corrupt receipt blocks. An unchanged prior `IMPLEMENTATION.md` fails validation as "not output from this code step".

### AC-05 — transient parity, protected single attempt

Attempt: I checked attempt limits, the five-second delay, and the frozen ceiling in `RunParticipantStep`. I ran `TestExhaustionValidOutputAndSingleAttemptTransient`. I checked `SingleAttempt` on implementation, fix-up, protected steps, and started drafters.
Result: not broken. Ordinary steps keep two attempts, the delay, and the frozen ceiling. Code-mutating and protected steps run once. No new automatic retry exists for them.

### AC-06 — owner order over the whole batch

Attempt: I ran the table in `TestOwnerSuccessionOrder`. I tried the same-holder case, simultaneous third loss, an unusable first survivor, and a missing counterpart.
Result: not broken. `SelectSuccession` (`internal/quota/roles.go:113`) applies the owner order after all losses settle. An unusable selected successor blocks; the selector never skips or substitutes. An unresolved implementer blocks. `validateRoleDecision` re-runs the selection and rejects any deviation.

### AC-07 — frozen facilitator exception only

Attempt: I ran `TestExhaustionRoleBatchFloorAndExternalFacilitator`. I traced the single membership-addition path in `Batch.Validate` (`internal/quota/history.go:121`) and the known-signer update in `ReadHistory`.
Result: not broken. Only the facilitator implementer is appended, and only through a validated transition. The facilitator never counts toward the floor or usable survivors. An external organizer loss can be role-only with no membership change. Below the floor, no transition applies and the decision blocks.

### AC-08 — one handover, byte-exact chain, checked resolver

Attempt: I ran `TestActiveExecutionHistoryChainAndCorruption` and `TestSelectedHandoverChildAndReplay`. I wrote a filename-order adversary check into my reading of the test (a `zz-fake.md` beside the handover). I changed predecessor bytes and swapped in a symlink in the test flow. I then audited every current-state consumer.
Result: not broken. The resolver selects only from validated typed history. Filename order supplies nothing. A changed base hash, a corrupt tip, or a symlink blocks without fallback. A pending handover admits only the selected successor's handover action (`internal/runner/telemetry.go`, `selectedHandover`). All consumers I audited route through the resolver: driver checks, evidence table and report guard, driver evidence, cursor and phase digest, `advanceFinal`, implementation status and completion, fix-up paths, goal-check prompt, signoff recovery, learn, pipeline, TUI, retro, and action identity. `ExecutionStarted` keeps a pre-Phase-5 handover from implying review readiness.

### AC-09 — draft integrity

Attempt: I ran `TestDraftRecoveryIntegrityAndExactBinding`, `TestSuccessionSettlesWholeBatchAndRetainsFindings`, and `TestDraftArchiveReplayPreservesNewCanonical`. I tried signed, truncated, contradictory-identity, and duplicate-author drafts.
Result: not broken. Complete drafts stay. Only structurally valid incomplete unsigned drafts move, byte-exactly once, to a recorded archive. Any signoff, emoji, or frozen status blocks relocation. Participant rounds and reviews never move. Replay does not touch the successor's new canonical draft.

### AC-10 — retained obligations and review independence

Attempt: I read `CaptureQuotaObligations` and `UnresolvedQuotaObligations` (`internal/protocol/quota_retained.go`). I checked reviewer exclusion in `checkGates` and `newDriverImplOps`. I traced the single-reviewer proof and its re-derivation sites.
Result: not broken. Vetoes, CRITICAL/MAJOR findings, disputes, and the successor's own earlier findings are retained as obligations with independent-evidence dispositions. The successor implementer is excluded from reviewers; the frozen facilitator is excluded too. `SingleReviewerAfterDropout` proves the full before/after cause, requires at least two before reviewers and exactly one after reviewer, and requires known distinct snapshot models. It is re-derived at precommit (`Settle` → `checkGates`) and through `ReviewerStatus` at review, review consensus, goal check, and close (`internal/app/driver_impl.go:789`, `internal/driver/impl.go:297`).

### AC-11 — Codex inert outside the exact contract

Attempt: I checked argv parsing in `TestExhaustionExplicitEffectiveArgs`, version binding in `recordLaunch`, and the four exact messages against `source-context/codex-native-source.md`. I read the smoke receipt.
Result: not broken. Recognition needs explicit effective `exec --json`, discovered CLI version exactly `0.165.0`, and a terminal primary `turn.failed` envelope with one reviewed message. The source excerpt confirms the single-key error body and the four verbatim strings. Added detail, mixed terminals, quotes, and near-misses reject. The preflight prints the support notice. Limits stay disclosed: no binary equivalence and no captured native positive. The smoke receipt covers ordinary envelope behavior only.

### AC-12 — capture defaults off, private, bounded

Attempt: I read `internal/config/capture.go`, `internal/telemetry/capture.go`, and `internal/runner/capture.go`. I checked enablement, scrub conservatism, permissions, symlinks, caps, retention, and cleanup. I grepped for any authority path that reads the store.
Result: broken in test robustness only; see the MINOR finding. The product code is correct: machine file only, `PARLEY_HOME` and deck layers cannot enable it, the scrubber writes nothing when uncertain, redactions and truncation are marked, directories are 0700, files are 0600 and exclusive, symlinks refuse, and no recognizer or reducer reads the store.

### AC-13 — exact hunks, staged core, size limits

Attempt: I rebuilt each target from its before-bytes plus the 29 signed hunks and compared SHA256. I checked the header zones. I ran `TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail` myself.
Result: not broken. Live copy: before `aae6e830…`, after `b9284d7d…`. Bundled CLI default: `43f2be0e…`. Skill reference: `b9284d7d…`. Staged core: `2.19.0` `16323e04…` plus the hunks equals `2.20.0` `6c437ed5…`. All four reproduce exactly; each header zone is unchanged. `SKILL.md` is 19,989 bytes. The packet test measured 68,448 bytes unflagged and 69,880 bytes with `protocol_change`, both under 70,000. The unflagged case is retained. The two disclosed compaction changes are present, and I found no rule loss in the compacted zones (§4.0, §7, §14, LE-4); the goal-check ceiling rule survives in the references and protocol copies.

### AC-14 — versions and guidance

Attempt: I checked `internal/app/version.go` (1.56.0), skill `package.json` and `compatibility.json` (2.20.0), `docs/credit-exhaustion.md`, README, CHANGELOG, `parley quota capture` help, and the skill references and templates.
Result: not complete. The content I checked is consistent. The full Go suite was still in progress at review time, so I do not record AC-14 as passed.

### AC-15 — release gates

Attempt: I checked for release evidence.
Result: not passed, correctly. No channel is released. Review consensus, fresh goal check, merge/release hashes, Homebrew, WinGet, runtime hashes, and the separate channel review are all outstanding. This is the recorded release gate.

### Decision spot checks (D1–D12)

D1: v2 trigger frozen at kickoff; v1 and omitted-trigger policies keep behavior and bytes. D2: reset evaluation is separate and never negates exhaustion. D3: suppression is common-launch admission and survives restart and blocked batches. D4: the full owner order is enforced once over the whole loss set. D5: only the frozen facilitator enters membership, without a floor or review seat. D6: one successor-owned handover per transition, base-bound by hash. D7: shared-draft recovery retains every filed position. D8: the causal single-reviewer exception changes only the numeric minimum. D9: exact-version, explicit-JSON Codex only; default text, Kimi, and Claude stay diagnostic. D10: capture is machine-only, scrubbed, bounded, and non-authoritative. D11: the exact signed hunks are on all four copies. D12: versions and release notes are staged; release execution stays behind AC-15.

## Findings

### [MINOR] Capture unsafe-mode test depends on the process umask

`TestZcodeCaptureRejectsUnsafeStorage/file-mode` (`internal/telemetry/capture_test.go:97-98`) writes `inv.log` with mode 0644 through `os.WriteFile`. Under umask 077 the file is created 0600, the fixture then looks safe, and `CaptureStatus` succeeds. The test fails with "unsafe status/purge accepted".

I reproduced this myself. My shell umask is 077. I ran `go test ./internal/telemetry -run TestZcodeCaptureRejectsUnsafeStorage -count=1`; only this subtest failed. The five focused packages otherwise passed (`quota`, `membership`, `protocol`, `runner` all `ok`; `telemetry` failed only here).

Why it matters: the "focused tests pass" claim is environment-dependent. A strict-umask machine or CI reports a failure that is not a product defect. The full-suite gate at close may then block on a test bug.

The product code is not at fault. It creates files with explicit 0600, and umask can only tighten that mode.

Fix: after `os.WriteFile`, force the mode with `os.Chmod(filepath.Join(dir, "inv.log"), 0644)`. The `directory-mode` case already uses `os.Chmod`; make the file case match it. Location: `internal/telemetry/capture_test.go:98`.

### [NIT] wait message names only the initial execution record

`outstandingAgents` (`internal/app/wait.go:566`) prints "IMPLEMENTATION.md not filed" when no started implementation exists. After a succession, the authority is `role-handover/<transition-id>.md`. The digest itself is resolver-routed; only the display string is stale. Fix: print the active path from `protocol.ExecutionPath`.

### [NIT] Empty Unreleased section above the new version header

`CHANGELOG.md` keeps `## Unreleased` with no items directly above `## 1.56.0 — 2026-10-29`. Remove the empty header or keep unreleased items under it.

## Open questions

1. The full Go suite must pass at the reviewed tree before close. It was still running when I wrote this review. Will it run under umask 022, or will the MINOR finding surface there too?
2. AC-15 release evidence is not produced yet. Which process will verify hosted Windows skill hashes after the release workflow replaces portable assets?
3. The Zcode grammar keeps the owner-deviation provenance and the native-positive limit. The open follow-up idea owns the completeness precondition. Is that follow-up scheduled after this release?
4. A structurally invalid handover draft from the selected successor blocks the handover action itself, because a corrupt tip blocks without fallback. Repair then needs a manual file removal or owner action. Is that the intended operator runbook, and should the docs state it?

Independence: `review/round-01/` was empty when I started. I wrote this file in a separate process. I did not read the other reviewer's round-01 file before I finished this review. I read the full live protocol at the attested hash before any code review.
