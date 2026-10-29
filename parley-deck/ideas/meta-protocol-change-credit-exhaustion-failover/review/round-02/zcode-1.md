---
agent: zcode-1
idea: meta-protocol-change-credit-exhaustion-failover
review-round: 2
date: 2026-10-29
reviewed-commit: 31c039c41a455a03385401085d132234a17603f6
responding-to: [kimi-1/review/round-01]
---

## Summary

I re-reviewed the full CLI diff `c2d2431..31c039c` and the skill diff `f557c5e..70d29c1` against FINAL D1–D12 and AC-01 through AC-15. All six agreed fixes F1–F6 are applied and verified at the fixed source. I ran my own focused suites, my own hunk reconstruction, and my own hash checks. I found no new CRITICAL, MAJOR, or MINOR defect. One NIT follows below.

The F1 disposition satisfies the signed requirement. I concur. The eager fallback/default calculation stays read-only. It grants no protection at either protection site. Only `provider-exhaustion.v1` losses request succession. Selection uses membership after every settled loss. Immutable history re-derives the same result.

The Windows disposition is scoped and honest. I concur with it. The retained log shows only platform artifacts on top of 149 pre-existing baseline failures. CLI Windows stays experimental. CLI WinGet stays held under D12.

Release gates stay open, and this review does not close them. The required full Go command `go test ./... -count=1 -timeout 45m` was still running when I wrote this file. I saw the live process and the three-line log. I do not treat a pending result as a pass. No hosted run of the fixed source exists yet. AC-15 release and installation evidence does not exist yet. These are gates, not product findings.

Provenance: I tag each verdict below. All executed checks and source citations are `PRIMARY`. I made no provider call and induced no exhaustion.

Protocol attestation: context_mode full; source_sha256 and packet_sha256 `b9284d7d7062ef49d59439a23117dd7434e494340b5b564a64651a12eb324709`; fallback_reason none.

## Position changes since prior review round

1. My round-01 AC-05 and AC-08 results carried no provenance tags. The consensus read them as `RECALL` and `UNVERIFIED`. I accept that reading. Kimi's counterexamples were real at `931550c`. I verified the sites myself this round: the old `Roles.Protected` included `Assignment.Implementer`, and `RebuildDetail` overwrote the stat result. My "not broken" claims for the complete scope were wrong. The narrower test results I reported still passed then and pass now.
2. My round-01 MINOR (umask fixture), NIT (wait wording), and NIT (empty `Unreleased` heading) are all fixed. I verified each fix at the current source.
3. My round-01 hunk reconstruction, packet-size result, capture product-code assessment, and release-gate reading all reproduced exactly this round. No position change there.

## Responses to other reviewers

### @kimi-1 — round-01

Your three MAJOR findings were correct at `931550c`. I confirmed each site in the old and new source:

- V2 transient protection: fixed. `internal/quota/quota.go:143-145` no longer protects `Assignment.Implementer`. The inline map at `internal/quota/quota.go:240` drops it too. `internal/app/preflight.go:971` sets `SingleAttempt: opts.QuotaRoles.Protected(a.ID)`, which is now false for a fallback/default identity. Your requested gating of succession is present at `internal/quota/quota.go:275` and `internal/quota/roles.go:210`: both pass `exhaustedCandidateIDs(d.Candidates)`, never all candidate IDs. I grepped for other protection sites (`grep -rn 'Assignment.Implementer' internal/ --include='*.go'`); only the three read-only assignments remain (`internal/app/app.go:1967,1970`, `internal/membership/succession.go:90`). There is no third protection site.
- Version metadata: fixed. `cat VERSION` printed `1.56.0`; `internal/app/version.go:3` says `const version = "1.56.0"`.
- Phase-read regression: fixed. `internal/driver/cursor.go:270-276` applies `ExecutionStarted`'s boolean only when its error is nil and keeps the stat evidence. I ran `go test ./internal/driver -run 'TestRebuildDetailSurfacesReadErrors' -count=1 -v`; all three subtests passed, including `unreadable_IMPLEMENTATION.md`.
- Umask fixture: fixed as you specified; see AC-12 below.

The whole-batch tests you required exist and pass: `TestExhaustionSuccessionUsesNativeRoleLossAndWholeBatch` covers mixed ordinary-plus-exhaustion losses, both-role loss with an ordinary first survivor, an unusable ordinary counterpart, an ordinary pinned implementer, and an ordinary designated implementer, for both `legacy` and `global-default` sources. I ran it verbosely; all 14 subtests passed. `Kickoff.Validate` runs in each case, so the immutable history validation is exercised, not only `Roles.Protected`. I concur with the AC-05 and AC-08 conflict resolutions as written in `review/consensus.md`.

### Both round-01 files — fix plan dispositions

- F1: I concur. See Summary and Refutation attempts. The eager calculation is a pure derivation over current prompt, config, and history. It stores nothing. An exhaustion-backed transition binds the identity and source in immutable history (`internal/quota/roles.go:150-158`, `internal/quota/history.go` batch validation). After a transition, `internal/membership/succession.go:24-27` reads `previous.After` and ignores the config chain, so later config changes cannot rewrite a bound successor. Validated history outranks retired pins and designations, as the protocol requires.
- F2, F3, F4, F5, F6: applied and verified; commands and outputs below.
- New wait test: verified. `TestWaitImplementationAbsentOrEarlyHandover` runs both `handover=false` and `handover=true` through the real wait command. The `handover=true` case commits a validated kickoff with a role transition, writes the successor handover, and checks `protocol.ResolveActiveExecution` before waiting. Both subtests passed and expect exit 3 with the new message.
- Readiness tests: verified as claimed. `TestExhaustionReadinessUnpinnedDefault` uses real synthetic child processes (re-exec of the test binary). The transient case asserts two attempts with at least five seconds between them. The native case asserts one attempt and `RoleTransition.After.Implementer == "beta"`. The global-default case sets `default_implementer = "alpha"`, drops alpha through an ordinary failure, and checks `designationOps` resolves `"beta"` through the actual resolver.
- Windows disposition: I concur. See AC-14 and the Windows assessment below.
- Driver refusal: the record is consistent. FINAL D12 says native fallback receipts supply no driver accounting or membership authority. The IMPLEMENTATION and consensus record it the same way.

## Refutation attempts

I assumed each criterion was broken and tried to break it. Commands ran at HEAD `7a1d498` with a clean tree; `git diff --stat 31c039c..HEAD -- . ':!parley-deck'` is empty, so no product change follows the target commit.

### AC-01 — v2 freeze, v1 preservation, owner-only widening

Attempt: break legacy byte stability; widen a saved policy without owner authority.
Result: `CONFIRMED (PRIMARY)` for executed scope. Under `umask 077` I ran `go test ./internal/quota ./internal/membership ./internal/protocol ./internal/driver ./internal/config -count=1 -timeout 12m`. Output: `ok parley-deck-cli/internal/quota 4.007s`, `ok parley-deck-cli/internal/membership 36.684s`, `ok parley-deck-cli/internal/protocol 1.040s`, `ok parley-deck-cli/internal/driver 25.897s`, `ok parley-deck-cli/internal/config 1.095s`. That includes `TestRoleFieldsDoNotChangeLegacySerialization` and `TestOwnerWideningFreezesRolesWithoutKickoffRewrite` (I also ran the latter verbosely: `--- PASS`). New runs select v2 at `internal/app/app.go:1952` and `:2509` (diff read).

### AC-02 — native framing only, separate reset, valid output wins

Attempt: break the reset separation; feed near-miss framing.
Result: `CONFIRMED (PRIMARY)` for synthetic fixtures. `go test ./internal/telemetry -count=1 -timeout 10m` under `umask 077`: `ok parley-deck-cli/internal/telemetry 2.527s`. Source: `internal/telemetry/exhaustion.go:200-205` takes a reset hint only from the unchanged legacy classifier. A failed hint leaves `e.ResetAt` nil and the reason "explicit provider exhaustion; reset unknown". `internal/quota/exhaustion.go:41-43` rejects a non-nil short reset as fabricated evidence; real recognizers never emit one. Valid output wins in `RunParticipantStep`; the in-flight window in `ExhaustionAdmission` requires `StartedAt` before the failure's completion and completion after it. No native positive is claimed.

### AC-03 — one child, durable suppression

Attempt: replay under a changed run, prompt, phase, and blocked batch.
Result: `CONFIRMED (PRIMARY)`. `go test ./internal/runner -run 'TestExhaustion|TestSelectedHandover|TestActiveExecution|TestDraft' -count=1`: `ok`. The readiness native case asserts exactly one child. Suppression derives from the durable invocation ledger in `internal/telemetry/exhaustion_admission.go:78-143`, not memory. `beginLaunch` calls it before any child (`internal/runner/telemetry.go`, exhaustion admission block). `SupplementUsability` calls it for observed and unobserved identities, so a blocked batch still suppresses.

### AC-04 — recovery admits no new child; stale artifacts cannot mask

Attempt: remove or corrupt a validation receipt; present an old valid execution record.
Result: `CONFIRMED (PRIMARY)`. The runner suite above passed, including the recovery and stale-artifact tests. `internal/runner/dropout.go:163-179` recovers a missing receipt only for the last invocation, without a child. `internal/runner/exhaustion.go:69-82` rejects an unchanged prior record as "not output from this code step".

### AC-05 — transient parity, protected and code-mutating single attempt

Attempt: break the corrected pre-pin path; find a protected role that lost its gate.
Result: `CONFIRMED (PRIMARY)`. The transient readiness case asserts two attempts and a five-second minimum spacing (`proof[1].ObservedAt.Sub(proof[0].ObservedAt) >= 5s`). Designee, pin, organizer, facilitator, and started-drafter protection remain in `Roles.Protected` and the inline map. `internal/runner/runner.go:570` keeps the implementation phase single-attempt for any implementer. `runExhaustionFixup` uses `SingleAttempt: true`. No new retry exists for code-mutating steps.

### AC-06 — owner order over the whole batch

Attempt: make an ordinary failure synthesize a transition; make the selector skip or substitute.
Result: not broken, `CONFIRMED (PRIMARY)`. `TestOwnerSuccessionOrder` (8 subtests) and `TestExhaustionSuccessionUsesNativeRoleLossAndWholeBatch` (14 subtests) passed under my run. `SelectSuccession` (`internal/quota/roles.go:113-160`) returns an error for an unusable selection and never searches past it. The both-lost case selects from `s.Order` over `d.After`, so an ordinary-failed member cannot win a seat; the "ordinary first survivor" subtest proves it (successor `z`, not `peer`).

### AC-07 — frozen facilitator exception and floor

Attempt: add another membership identity; count the facilitator toward the floor.
Result: `CONFIRMED (PRIMARY)`. The full `./internal/quota` run passed, including `TestExhaustionRoleBatchFloorAndExternalFacilitator`. `internal/quota/quota.go:282-284` appends only the facilitator who becomes implementer. `internal/quota/quota.go:199` excludes the facilitator from `UsableSurvivors`. `validateRoleDecision` re-checks the floor and the single addition.

### AC-08 — one handover, byte-exact chain, checked resolver

Attempt: break the F3 error path; find a fixed `IMPLEMENTATION.md` reader that bypasses the resolver; break a base hash.
Result: `CONFIRMED (PRIMARY)`. `TestRebuildDetailSurfacesReadErrors` passed (three subtests). The resolver (`internal/protocol/execution.go:29-79`) selects only from validated typed history, rejects non-regular files and symlinked parents, and returns `PendingHandover` for a missing tip. Every current-state consumer I inspected routes through it: `driver/quota.go` (`quotaBefore`), `phase58.go` (implementation, fix-up, prompts, review context), `evidence_table.go`, `driver_evidence.go`, `phasedigest.go` (`StatExecution`, `implSection`), `gates.go` (`StatExecution`), `cursor.go` (`ExecutionPath`), `wait.go` digest, `learn.go`, `pipeline_cmd.go`, `retro.go`, `protosnap.go`, and `membership.Roles` (`ReadExecutionFrontmatter`). A resolver error in `implSection` surfaces as `Present: true, Status: "unparsed"` instead of hiding the defect.

### AC-09 — draft integrity

Attempt: move a signed, malformed, or complete draft.
Result: `CONFIRMED (PRIMARY)`. `./internal/protocol` and `./internal/membership` passed, including `TestDraftRecoveryIntegrityAndExactBinding`, `TestDraftArchiveReplayPreservesNewCanonical`, and `TestSuccessionSettlesWholeBatchAndRetainsFindings`. `internal/protocol/draft_recovery.go:17-81` considers only drafts of exhausted authors, blocks signed or frozen drafts (`draftSignoff`, emoji, `status: final|complete`), blocks truncated fences and contradictory identity, and skips complete drafts. `reconcileDrafts` moves bytes exactly once to a hash-bound archive.

### AC-10 — retained obligations and review independence

Attempt: break the before/after reviewer proof; let a successor review its own code.
Result: `CONFIRMED (PRIMARY)`. The membership suite passed. `internal/membership/review_gate.go` computes `beforeReviewers` with the before-implementer excluded, keeps the frozen facilitator excluded from both sets, permits the facilitator-implementer only through `D5`, and requires exhaustion evidence for any removed protected role (`exhaustionCandidate` revalidates the proof). `internal/membership/gates.go:41-48` forces the successor implementer into the reviewer-exclusion check and suppresses stale designations after a transition. Retained obligations flow through `CaptureQuotaObligations` at `membership.go:351-355`.

### AC-11 — exact Codex contract

Attempt: enable recognition without explicit JSON or the exact version; accept a near-miss message.
Result: `CONFIRMED (PRIMARY)` for the source-backed scope. `internal/telemetry/exhaustion.go:21-53` reads effective argv only; `:57-76` requires explicit JSON and `CodexContractVersion == "0.165.0"`; `:78-137` requires a single-key `turn.failed` envelope with one of the four exact messages, rejects unknown events, duplicate keys, mixed terminals, and divergent prior errors. `internal/quota/exhaustion.go:44-48` re-checks version, JSON flag, message, and the absence of reset data. The smoke receipt exists at `source-context/codex-json-smoke-receipt.json`. No binary equivalence and no captured native exhaustion are claimed, as FINAL limits require.

### AC-12 — capture defaults off, private, bounded

Attempt: enable capture through PARLEY_HOME, deck, or environment; break scrubbing, permissions, or caps; find an authority path that reads the store.
Result: `CONFIRMED (PRIMARY)`. Under `umask 077`, `./internal/config` and `./internal/telemetry` passed. `internal/config/capture.go` reads only `$HOME/.parley/agents.toml`; `TestCaptureOnlyMachineConfiguration` proves `PARLEY_HOME` and `PARLEY_CAPTURE_ZCODE_FAILURES` cannot enable it and that a string value errors. `internal/telemetry/capture.go:30-95` writes nothing when framing is uncertain and marks `[REDACTED]`/`[TRUNCATED]`. Modes, symlinks, caps, retention, and exclusive creation are enforced in `captureDirectory`, `captureFiles`, and `CaptureZcodeFailure`. The only readers are `runner/capture.go` (write path) and the owner status/purge command. `ClassifyExhaustion` reads live stderr from the collector, never the store.

### AC-13 — exact hunks, staged core, size limits

Attempt: reconstruct each copy from its own before-bytes plus the 29 signed hunks and find a mismatch.
Result: `CONFIRMED (PRIMARY)`. The signed JSON hashes to `f140e81ada0a8fb08644109b889979e54945f829183c3d82b47c760eae67147c` (`shasum -a 256`). My own script applied each hunk exactly once per base (uniqueness enforced) and compared bytes:
- staged 2.19.0 `16323e04…` + hunks = 2.20.0 `6c437ed5…` — exact.
- live copy before `aae6e830…` (at `c2d2431`) + hunks = `b9284d7d…` — exact.
- bundled CLI default before `ea8f26f1…` + hunks = `43f2be0e…` — exact.
- skill reference before `aae6e830…` (at `f557c5e`) + hunks = `b9284d7d…` — exact.

Packet test executed: `go test ./internal/app -run 'TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail' -count=1 -v`. Output: `facilitator body=68448 B`; `phase-1 protocol_change facilitator body=69880 B`; both under 70,000. `wc -c skills/parley-deck/SKILL.md` = `19989`, under 20,000. Header zones are unchanged because reconstruction from the baselines is exact.

### AC-14 — versions, guidance, full suites

Attempt: find a version mismatch or missing guidance; close the criterion without the full suite.
Result: metadata and guidance `CONFIRMED (PRIMARY)`; full-suite portion open, not passed. `VERSION` and `internal/app/version.go` both say 1.56.0; the version tests passed in my app run. `docs/credit-exhaustion.md`, README, CHANGELOG (F6 verified: no empty `Unreleased` heading), the skill `EXHAUSTION_AND_HANDOVER.md` reference, and `parley quota capture` help are consistent with FINAL. The required full command was still running (live `go test ./... -count=1 -timeout 45m` process; `fixup01-full-go-v2.log` holds three lines). I do not record AC-14 as satisfied until it passes. The skill commit is unchanged at `70d29c1`; its prior 402 Node and 54 Python results are implementer-reported, not rerun by me.

Windows assessment: the retained log `review01-cli-windows.log` and `windows-baseline-comparison.json` show baseline commit `f3b8a4f` with 149 top-level failures and the review01 source with 154; 149 shared, 5 added. The five added failures are: four `sync ... Access is denied` durable directory-sync errors in the new v2 tests (`TestActiveExecutionHistoryChainAndCorruption`, `TestSuccessionSettlesWholeBatchAndRetainsFindings`, `TestDraftArchiveReplayPreservesNewCanonical`, `TestOwnerWideningFreezesRolesWithoutKickoffRewrite`) and the capture private-mode test. The capture test now expects refusal and no file; `captureDirectory` requires `Perm() == 0700`, which Windows never reports, so the product refuses before any write. `docs/credit-exhaustion.md` states this limit. The comparison record itself says it "does not establish Windows correctness". I concur with the disposition: no clean baseline is claimed, CLI Windows stays experimental, and CLI WinGet stays held under D12.

### AC-15 — release gates

Attempt: find release, channel, or installation evidence.
Result: `UNVERIFIED`, correctly open. No release or installation exists. A fresh independent goal check, channel hashes, runtime installs, and a separate channel review remain. This is a release gate, not a product finding.

### Decision spot checks (D1–D12)

D1: v2 constructors on new runs; legacy tests passed. D2: reset separation verified in source. D3: ledger-derived admission verified. D4: whole-batch selection verified. D5: single membership addition and floor verified. D6: resolver, base hashes, and `BindHandover` verified. D7: draft recovery verified. D8: causal reviewer proof verified. D9: exact-version recognizer verified. D10: machine-only capture verified. D11: all four copies reconstructed exactly. D12: versions and notes staged; publication gates untouched.

## Updated findings

I found no new CRITICAL, MAJOR, or MINOR defect. All round-01 findings are resolved by F1–F6, which I verified above.

### [NIT] `parley quota capture status` deletes expired evidence

`CaptureStatus` calls `captureCleanup(dir, now, 0, false)` (`internal/telemetry/capture.go:225-233`). With `purge=false` it still removes every file older than seven days and rewrites `cleanup.json`. So the inspection command mutates the store it inspects. An owner who inspects retained evidence before a completeness review silently loses files already past retention.

FINAL D10 says "On each use, delete expired files first". The words cover this behavior, so this is not a rule violation. It is an operational surprise in a diagnostic facility.

Concrete fix: keep `status` read-only. Report the expired count without deletion, and let only `purge` and a new capture enforce retention. Alternatively, state in `docs/credit-exhaustion.md` that `status` also enforces retention. Location: `internal/telemetry/capture.go:219-234`, `internal/app/quota_capture.go:24-31`.

## Open questions

1. Will `go test ./... -count=1 -timeout 45m` pass at `31c039c`? Its log `fixup01-full-go-v2.log` is incomplete. The consensus makes a pass a close condition.
2. No hosted run of the fixed source exists yet, on any OS. Which run will supply hosted confirmation before release?
3. AC-15 needs the fresh independent goal PASS, channel hashes, runtime installs, and a separate channel review. Which process performs the channel verification that D12 names?
4. The driver refused phase 6 with "has no runs yet". The release inbox note must keep listing the pending owner commands, and the native fallback receipts must stay clearly separated from driver accounting, as D12 requires.
5. Does the owner want the NIT above fixed in this cycle or deferred to the open Zcode follow-up? It changes no acceptance criterion.

Independence: `review/round-02/` was empty when I started. I did not read the peer's round-02 output. I wrote this file in a separate process. I read the full live protocol at the attested hash before the code review.

---

## Supplemental assessment — 2026-10-29 (round 02, after new CI evidence)

This addendum answers the review02-supplement brief. New hosted and local evidence arrived after my main file above. I assessed it in a fresh process. Every verdict below is `PRIMARY` unless tagged otherwise. I executed each check myself at the unchanged product source `31c039c`; `git diff 31c039c..HEAD --stat -- ':!parley-deck'` is empty and `git status --porcelain -- ':!parley-deck'` is empty.

Protocol attestation for this process: context_mode full; source_sha256 and packet_sha256 `b9284d7d7062ef49d59439a23117dd7434e494340b5b564a64651a12eb324709` (verified: `shasum -a 256 parley-deck/COOPERATION.md` prints that hash); fallback_reason none. The phase-6 driver refusal and the native-fallback boundary are unchanged from my main file.

### Self-correction (§15.1)

My main file, Open questions, item 2 says: "No hosted run of the fixed source exists yet, on any OS." That statement is now stale and I replace it. Hosted runs of the fixed source exist. Run `40963775887` (job `123762942534`, PR 87, merge `af6dfda` of `7a1d498` into `24c6184`) failed on Ubuntu in `TestStepBudgetConcurrentIndependentGroupsShareCap` at `step_budget_test.go:153: started=2 refused=5` (`PRIMARY`: I read `review02-cli-ubuntu.log`; the checkout SHA and the failing line are in the log). The corresponding push run `40963768431` passed Ubuntu at the same source (`SECONDARY`: the organizer's named claim in the inbox note; I have no local log of it). The correct open question is no longer "which run will supply hosted confirmation" but "which hosted matrix run passes after the concurrency defect below is fixed".

My main file, Open questions, item 1 asked whether the full command would pass. It exits 0 sometimes and fails sometimes. A pass cannot dispose the reproduced defect. See the finding below.

### New evidence assessed

1. Local full suite (`PRIMARY`). `fixup01-full-go-v2.log` shows `ok` for 34 packages, no failure, and the receipt `fixup01-full-go-v2.receipt.json` records exit 0 at product commit `31c039c` with log sha256 `3ea94698…`. This closes the "still running" state my main file recorded. It does not close AC-14: the receipt's own `limit` field says a separate 100-run check reproduced an unresolved admission race, and I reproduced that race myself (below). One passing run of a flaky suite proves nothing about absence of the defect.
2. Hosted Ubuntu failure (`PRIMARY` for the failure, `SECONDARY` for the passing push run). See the self-correction. The failed job checked out merge `af6dfda`; the branch tip `7a1d498` changes `parley-deck/` records only, so the product files equal `31c039c`. All other packages in that log are `ok`.
3. Local repetition (`PRIMARY`). I executed `go test ./internal/runner -run '^TestStepBudgetConcurrentIndependentGroupsShareCap$' -count=100 -timeout 5m`. Output: exit 1, 7 failures of 100 (`--- FAIL: TestStepBudgetConcurrentIndependentGroupsShareCap`, six `started=2 refused=5`, one `started=2 refused=4`). Log: `.parley-runtime/credit-failover/zcode1-review02-repro.log`. The diagnostic overlay log `review02-concurrent-budget-diagnostic.log` shows the missing launch in each failing run returns `unreadable invocation request during exhaustion admission`. I verified that error string exists at exactly one place in the product: `internal/telemetry/exhaustion_admission.go:94` (`grep -rn "during exhaustion admission" internal/ --include='*.go'`). No other code can produce it. This ties the failure to the new common admission scanner, not to the budget ledger.
4. Organizer hypothesis (`CONFIRMED (PRIMARY)` by source reading). `Invocation.write` opens `requested.json` with `os.O_CREATE|os.O_EXCL|os.O_WRONLY` at `internal/telemetry/record.go:409` and writes the bytes at `:413`. The file is observable while empty or partial. `ExhaustionAdmission` globs every `invocations/*/requested.json` and parses each before any identity filter: `internal/telemetry/exhaustion_admission.go:82-95`. The parse at `:93` fails on a concurrent launch's incomplete file, and `:94` refuses the wrong launch. `beginLaunch` calls `exhaustionAdmission` at `internal/runner/telemetry.go:119` before `recordLaunchRequest` at `:179`, so the refused launch writes no invocation record at all. That is why `refused` drops below 6: the eighth launch leaves no terminal record. The writer pattern pre-dates this idea (`git show c2d2431:internal/telemetry/record.go` shows the same `O_EXCL` write). The regression is the new reader: `git diff c2d2431..31c039c -- internal/runner/telemetry.go` adds the `exhaustionAdmission` call, and `internal/telemetry/exhaustion_admission.go` is a new file from commit `931550c`. The same parse-before-filter pattern pre-exists in `internal/runner/dropout.go:416-418`, but this idea moved such a scan into the hot path of every launch. The failing test itself pre-dates the idea (`git show c2d2431:internal/runner/step_budget_test.go` contains it).
5. Private proposed correction (`PRIMARY`; tested by overlay, product unchanged). Assessment below.
6. Capture-status NIT disposition. I concur; see the disposition response.

### New finding

#### [CRITICAL] Common exhaustion admission races with concurrent record creation and refuses unrelated launches

What is wrong: `ExhaustionAdmission` reads and parses every `requested.json` under `.parley-runtime/invocations/` on every launch (`internal/telemetry/exhaustion_admission.go:82-95`, called from `internal/runner/telemetry.go:119`). `Invocation.write` exposes `requested.json` empty or partial between exclusive creation (`internal/telemetry/record.go:409`) and the completing write (`:413`). A concurrent ordinary launch in the same root therefore fails admission with `unreadable invocation request during exhaustion admission` instead of its normal budget-refused outcome.

Why it blocks: the defect is introduced by this idea's new launch-path scanner and breaks a test that existed and passed before the idea. I reproduced it 7 times in 100 runs locally, and hosted run `40963775887` failed on it. AC-14's full-suite evidence cannot be satisfied while the suite is flaky-red, so AC-14 and AC-15 cannot close. The spurious refusal also leaves no invocation record, because admission runs before `recordLaunchRequest` (`internal/runner/telemetry.go:119` precedes `:179`); the refusal is silent in the audit ledger this CLI exists to keep. Impact scope: any concurrent launches in one root, regardless of idea, agent, or exhaustion policy, because the parse at `exhaustion_admission.go:93` precedes the identity filter at `:96`.

Concrete suggested fix: publish each record atomically from the writer, exactly as the private proposal does — stage the complete synced bytes in a temporary file inside the invocation directory, then publish with an exclusive `os.Link` so the published name never exists with incomplete bytes and never replaces an existing record. I verified this fix shape myself (next section) and it covers every `requested.json` reader at once (`exhaustion_admission.go:82`, `runner/dropout.go:405`, `trajectory/verification_recovery.go:322`, `trajectory/unchanged.go:160`, `trajectory/reconcile.go:317`, `budget/binding.go:401`, `budget/step_history.go:59`, `budget/cycle_history.go:189`). Do not fix the reader by tolerating parse errors: after atomic publication, any parse failure is genuine corruption and must keep blocking. Apply this as fix-up cycle 2, rerun the concurrent test at high `-count`, and only then re-evaluate AC-14.

### Executed checks on the proposed correction

I ran the receipt's exact command myself with the private overlay (no product edit):

```sh
go test -overlay .parley-runtime/credit-failover/record-atomic-proposal-overlay.json ./internal/runner ./internal/telemetry -run 'TestStepBudgetConcurrentIndependentGroupsShareCap|TestBeginBound|TestWriteFailuresAreReturned|TestInvocation' -count=100 -timeout 5m
```

Output: exit 0; `ok parley-deck-cli/internal/runner 26.078s`; `ok parley-deck-cli/internal/telemetry 197.113s`. Log: `.parley-runtime/credit-failover/zcode1-review02-atomic-check.log`. I also ran `go test -overlay .parley-runtime/credit-failover/record-atomic-proposal-overlay.json ./internal/runner -run '^TestStepBudgetConcurrentIndependentGroupsShareCap$' -race -count=20 -timeout 10m`: `ok parley-deck-cli/internal/runner 6.767s`, no race report.

Scope: `write` serves `requested.json` and `terminal.json` for every invocation, so the change covers all readers listed in the finding, not only admission. Non-overwrite semantics: `os.Link` fails with `EEXIST` when the target name exists in any form, including a symlink or a directory, so it preserves the exact `O_EXCL` refusal. `TestWriteFailuresAreReturned` (`internal/telemetry/record_test.go:237-259`) proves the failure paths still refuse: it plants a directory named `terminal.json` and requires `Finish` to fail twice. `TestBeginBoundRefusesEveryReuseAndCollision` covers collision refusal. Both passed 100 times under the overlay in my run. Unsafe path handling: the temporary file is created inside the caller-exclusive invocation directory with mode 0600, matching the old modes; the published link shares that inode. No symlink is followed or created. Crash residue: a crash between staging and cleanup can leave a `.record-*` file inside one invocation directory. I checked every enumerator: `budget/binding.go:390`, `budget/step_history.go:93`, and `budget/cycle_history.go:77` enumerate the invocations root and expect directories; glob readers match `requested.json` or `terminal.json` names only. No reader enumerates the interior of an invocation directory, so the residue is inert. Persistence: data bytes are synced before publication, same as before. Neither the old nor the proposed code fsyncs the directory entry after publication; that is parity, not a regression. An optional follow-up may fsync `i.Dir` after `Link`. Windows note: `os.Link` needs NTFS; CLI Windows stays experimental, and a failure there surfaces as a launch error, not silent corruption.

My verdict on the proposal: `CONFIRMED (PRIMARY)` for the tested scope. I concur with the approach and found no defect in it. I propose no different correction; writer-side atomic publication is the correct shape, and reader-side tolerance would weaken corruption handling against D2/D3.

### Effect on prior verdicts

- AC-14: my main file said "full-suite portion open, not passed". The evidence is now negative, not merely missing: hosted run `40963775887` failed, and my local repetition failed 7 of 100. The exit-0 full run in evidence 1 does not change this. AC-14 stays blocking until the CRITICAL finding is fixed and the suite passes, including a repeated run of the concurrent test.
- AC-15: unchanged (open), and now additionally blocked by the CRITICAL finding.
- AC-05: my verified scope stands; the retry-parity, spacing, and protection tests still pass. The race adds a separate concurrency defect. Risk I could not verify by execution: a spurious admission refusal inside a participant step surfaces as an invocation-less failure, and how that interacts with attempt accounting deserves a check in fix-up cycle 2.
- AC-03: suppression semantics and durability results stand. The common-admission concurrency behavior is defective; the CRITICAL finding covers it.
- My main file's claim "I found no new CRITICAL, MAJOR, or MINOR defect" is superseded by this addendum's CRITICAL finding.
- All other round-02 verdicts (AC-01, AC-02, AC-04, AC-06 through AC-13, D1–D12, F1–F6) are unaffected by the new evidence.

### Disposition response: capture-status NIT

The implementer proposes my named documentation alternative: keep D10's on-use cleanup, and state in `docs/credit-exhaustion.md` that `parley quota capture status` removes expired and over-capacity files and updates the cleanup receipt. I concur. FINAL D10 says "On each use, delete expired files first", so the behavior follows the ratified text; a read-only `status` would contradict it. This is an agreed clarification within FINAL; no owner question is needed. Apply the doc line in fix-up cycle 2.

### Independence

I did not read Kimi's round-02 output or any `review02-kimi-1-*` process log. I resolved the evidence myself: the hosted log, both repetition runs, the diagnostic overlay source, the product source, the git history, and the private overlay checks above. The product tree stayed unchanged; my checks wrote only under `.parley-runtime/credit-failover/`.

---

## Findings

This index supplies the heading that the Phase 6 template requires. The earlier organizer brief requested only `## Updated findings`. That brief omitted this required heading; this structural repair, appended 2026-10-29, closes the gap.

Each entry below points to the full statement elsewhere in this file. Each preserves its severity, claim, evidence, and suggested fix. Both findings remain open. This index adds no new finding, replaces none, and weakens none.

1. [CRITICAL] Common exhaustion admission races with concurrent record creation and refuses unrelated launches — status: open. Full statement: "New finding" in the supplemental assessment above. Evidence cited there: local repetition (7 failures in 100 runs), hosted run `40963775887`, and the named source sites. Suggested fix cited there: writer-side atomic publication of invocation records, to apply in fix-up cycle 2.
2. [NIT] `parley quota capture status` deletes expired evidence — status: open. Full statement: `## Updated findings` above, with its evidence and suggested fix. Agreed disposition: the documentation alternative in `docs/credit-exhaustion.md`, to apply in fix-up cycle 2.

No other CRITICAL, MAJOR, MINOR, or NIT finding is open in this file. All round-01 findings are resolved by F1–F6, verified above. "Effect on prior verdicts" records that the CRITICAL finding supersedes the main file's "no new CRITICAL, MAJOR, or MINOR defect" claim.

This repair changed structure only. No new test run, product change, or verdict change accompanies it.
