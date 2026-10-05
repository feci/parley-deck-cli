---
agent: claude-1
idea: meta-protocol-change-quota-auto-exclude
review-round: 4
date: 2026-10-05
reviewed-commit: d8b729af187a18ebe3891a36dfc379c4b04c8b95
skill-commit: dcb7d593130247709cfad19715aa1f53e5cc8be0
product-commit: 906857b9af33306158d3b0e91b94e6ef4b6f7458
current-head: a6f0b74bebc0b328ae2a8dd3e3fffc81ce34d70e
responding-to: [review/round-03/claude-1.md, review/round-03/consensus.md]
---

## Summary

This is the full-scope Phase 6 re-review after fix-up cycle 1. It is the single owner-authorized relaunch
of round 04. It covers CLI product `906857b` (review snapshot `d8b729a`) and skill `dcb7d59`.

Fix-up cycle 1 resolves both round-03 CRITICALs and most of the cycle-1 list. I checked these on this
deck's AppleVirtIOFS volume and on `/tmp`:

- New ideas are created with the knob on and off.
- The idea lease is exclusive across processes.
- Foreign or unknown host/boot owners are never reclaimed.
- Committed owner authority survives permitted inbox archiving and deletion.
- Every wrong-positive zcode family I demonstrated in round 03 is now rejected.

The implementation is **not ready for an attended close**. Three MAJOR issues remain:

1. **R4-MAJOR-1.** The zcode retry-wrapper grammar requires every retry's whole human message to be
   identical. The provider message contains a countdown that changes with time. So the SDK's standard
   "retries exhausted" report of a real exhaustion is always rejected. The only accepting retry fixture
   repeats one error object three times, which cannot happen natively.
2. **R4-MAJOR-2 (R2).** On knob-off ideas, re-inclusion now needs a new `included:` marker, and a catch-up
   join needs a committed owner answer with an exact directive. Neither the protocol nor the skill defines
   these. That is a third knob-off behavior change that the owner never listed.
3. **R4-MAJOR-3.** Native AC2 acceptance evidence is still missing. This needs an owner decision, and
   R4-MAJOR-1 makes it worse.

Counts: CRITICAL 0, MAJOR 3, MINOR 2, NIT 1. This review signs nothing and grants no close.

## Protocol context and coverage

- **Protocol packet (PRIMARY, executed 2026-10-05T19:31Z).**
  - Command: `parley protocol packet --dir . --phase 6 --track deliberation --idea
    meta-protocol-change-quota-auto-exclude --flag protocol_change --json` (parley 1.50.0).
  - Attestation: `context_mode=full`,
    `source_sha256=packet_sha256=73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e`,
    `fallback_reason` absent. Shadow packet 83,009 B, not used.
  - Body file: `.parley-runtime/protocol-packets/full-phase6-deliberation-73613f95….md`. Its SHA-256 equals
    the deck `COOPERATION.md` (1,501 lines, 125,862 B).
  - **Actual full read.** I read the ENTIRE body in four bounded chunks (lines 1–380, 381–760, 761–1140,
    1141–1502). Three lines exceed the reader's line limit (59: 1,976 chars; 452: 2,928; 1002: 2,205). I
    printed those three in full separately, so no truncation applies. This completes the coverage that
    round 03 lacked, where I read only selected sections.
- **Authorities read in full:**
  - `FINAL.md` (902 lines);
  - the owner scope-reset answer;
  - the new owner review-quota answer (quoted below);
  - my `review/round-03/claude-1.md` (594 lines);
  - the archived signed `review/round-03/consensus.md`, including R1–R4. Its SHA-256 `0feba153…cccc`
    matches the dispatch.
- **Living and producer artifacts:**
  - `IMPLEMENTATION.md` is a living file. I read it at working-tree SHA-256 `11f27fc5…727e`: the
    deviations, fix-up and resumption sections, plus the header.
  - Both producer maps, `codex-1-fixup-1-evidence.md` and `-host-evidence.md`.
  - The retained native tail, `source-context/provenance-review/codex-1-retained-native-tail.md`.
- **Snapshot identity (PRIMARY).**
  - `git diff --name-only 906857b HEAD -- internal cmd docs go.mod go.sum parley-deck/COOPERATION.md
    parley-deck/meta` is empty.
  - `906857b..a6f0b74` touches only orchestration files: IMPLEMENTATION, organizer notes and usage, the
    round-03 consensus move, host evidence, the quota-stop note and the ledger. CLI product is
    `906857b`, review snapshot `d8b729a`, current HEAD `a6f0b74`, skill `dcb7d59`.
  - All 24 CLI-product, 2 skill-product and 11 fix-up chunk SHA-256 values match `diff-manifest.json`.
- **Complete product diffs read: all 24 CLI chunks (`cli-product.diff`, 475,024 B) and both skill chunks
  (39,018 B), in full.** The incremental fix-up diff was not used as a substitute.

**Coverage table (every changed product file).** Every file listed was read in full in the product diff.
"Exec" names my own executed checks: probes P1–P5 (below), or the focused test run T1 (below).

| Area (chunks) | Files | Exec |
|---|---|---|
| docs (1) | `docs/cli-reference.md`, `docs/quota-membership.md` | P3 (grammar vs docs) |
| app (1–6) | `agents_exec.go`, `app.go`, `consensus_request_signoffs.go`, `driver_consensus.go`, `driver_impl.go`, `organizer.go`, `pipeline_cmd.go`, `preflight.go`, `preflight_liveness.go`, `quota.go`, `quota_consumers.go`, `quota_pipeline.go`, `quota_revision.go`, `quota_signoff.go`, `wait.go`; tests `quota_fixup_test.go`, `quota_signoff_test.go`, `quota_surfaces_test.go`, `quota_test.go` | T1 |
| config (6) | `runtime.go`, `runtime_test.go` | T2 |
| consensus (6–7) | `consensus.go`, `quota_test.go` | T1 |
| driver (7) | `driver.go`, `impl.go`, `loop.go`, `phasedigest.go`, `quota.go`, `quota_test.go` | T2 |
| fsutil (7–8) | `sync_darwin.go`, `sync_darwin_test.go` | T1, P1 |
| membership (8–10) | `fixup_test.go`, `gates.go`, `lock.go`, `membership.go`, `membership_test.go`, `revision.go`, `revision_test.go` | P1, P3, T1 |
| pidlease (10–11) | `lease.go`, `lease_test.go`, `live_unix.go`, `live_windows.go` (Windows: read only) | P1, P4, T1 |
| protocol (11–14) | `defaults/COOPERATION.md`, `participantartifact.go`, `quota.go`, `quota_manual.go`, `quota_retained.go`, `workspace.go` | P1, P3, T1 |
| quota (14–17) | `authority.go`, `authority_test.go`, `history.go`, `history_test.go`, `quota.go`, `quota_test.go`, `record.go`, `revision.go`, `revision_snapshot.go` | P3, P5, T1 |
| quotatest (17) | `authority.go` | via T1 |
| runcontrol, runmanifest (17) | `runcontrol.go`, `quota_test.go`, `manifest.go` | P1 (create); T2 |
| runner (17–19) | `acp.go`, `failclass.go`, `handoff.go`, `phase58.go`, `quota.go`, `quota_target.go`, `quota_test.go`, `runner.go`, `telemetry.go`, `telemetry_test.go`, `validation.go` | T2 |
| runplan, runstate, store (19) | `runplan.go`, `runstate.go`, `events.go` | read only |
| telemetry (19–23) | `provider.go`, `quota.go`, `quota_framing.go`, `quota_framing_test.go`, `quota_test.go`, `quota_zcode.go`, `quota_zcode_test.go`, `record.go`, `usage.go`; `testdata/quota/` README, harness and all 12 fixtures | P1 (zadv), P2, T1 |
| deck (23–24) | `parley-deck/COOPERATION.md`, `parley-deck/meta/protocol-changelog.md` | PRIMARY: deck and defaults hunks identical (`diff` exit 0) |
| skill (1–2) | `SKILL.md`, `parley-addon.json`, `references/COOPERATION.md`, `references/ROSTER_AND_PROTOCOL.md` | `cmp` deck vs skill copy exit 0 |

**Executed checks (PRIMARY; raw logs under the git-ignored `.parley-runtime/claude1-r4/`).** No provider was
invoked. Discovery was shadowed by `fixup-1/host/no-provider-bin` on `PATH`. No product, test, roster,
credential, worktree or other-agent file was edited.

- **P1, `probes/rerun.sh`.** My UNCHANGED round-03 programs `create`, `knoboff`, `reinclude`, `lease` and
  `zadv`, run on the shared volume and on `/tmp`. Logs: `probes/shared.log`, `probes/local.log`,
  `probes/zadv.log`.
- **P2, `zretry/main.go`.** Changing-countdown zcode records and RetryError wrappers. Log:
  `zretry/zretry.log`.
- **P3, `r4extra/main.go`.** R2 knob-off grammar and return. Logs: `r4extra/r4extra.log`,
  `r4extra/r2-return.log`.
- **P4, also `r4extra`.** R3 foreign, unknown and same-host/boot lease owners, on both volumes.
- **P5, also `r4extra`.** R1 committed authority in a scratch git repository, on both volumes.
- **T1.** `TMPDIR=<shared volume>/.parley-runtime/claude1-r4/tests/shared-tmp go test ./internal/app
  ./internal/membership ./internal/pidlease ./internal/fsutil ./internal/consensus ./internal/quota
  ./internal/telemetry ./internal/protocol -run
  'TestQuota|TestLease|TestSyncFile|TestEmbeddedDefaultMatchesLiveDeck' -count=1`.
  - Result: all 8 packages `ok`, exit 0 (`tests/shared-focused.log`).
  - T1 includes the unweakened native `TestLeaseProcess` on AppleVirtIOFS. The host boot identity is
    available: `{ sec = 1791044926 … }`.
- **T2.** The same shared-volume TMPDIR pattern: `go test ./internal/driver ./internal/runner
  ./internal/config ./internal/runcontrol ./internal/runstate -run 'TestQuota' -count=1`.
  - Result: `ok` for driver, runner, config and runcontrol. runstate had no matching tests. Exit 0
    (`tests/shared-focused-2.log`).

**Not executed by me (SECONDARY, producer evidence only):**

- the full `go test ./...` (host, 529 s);
- the race run;
- `go vet`;
- the Windows cross-compile;
- the full skill suite.

Windows runtime is unexecuted by anyone. Packets for phases 0, 5 and 8 were not re-rendered this round.
The packet is full-mode and its SHA equals the deck, so their content follows from that.

## User direction

- **Scope/reset answer.** Unchanged. My exact verbatim quote of
  `parley-deck/inbox/user-to-codex-1_meta-protocol-change-quota-auto-exclude_scope-reset-answer.md` is in
  `review/round-03/claude-1.md:62-109`. I re-read the answer in full this round, and it matches that
  quote. The bounded zcode stderr rule and the display-clock rule are owner policy. I check the code
  against their exact gates and do not relitigate them.
- **New owner answer.** Quoted verbatim from
  `parley-deck/inbox/user-to-codex-1_meta-protocol-change-quota-auto-exclude_review-quota-answer.md`
  (body below its frontmatter):

> ## Owner answer to `codex-1-to-user_meta-protocol-change-quota-auto-exclude_review-quota-stop-20261004.md`
>
> The answer was relayed by the owner's Claude Code session. The question and the answer, verbatim (Slovak),
> with a translation below:
>
> > Question: "Implementácia je hotová: obe etapy aj prvé kolo opráv, všetky testy sú zelené. Review claude-1
> > (kolo 4) však o 06:54 spadlo na gateway s chybou „429 All claude accounts have exhausted their quota
> > (reset after 5m)“. Claude teraz znova odpovedá, o 06:57 vrátil PONG. Povoliť jedno opakovanie review?"
> > Selected: **"Áno, raz teraz (Recommended)"**. The option read: "Codex spustí rovnaké review kolo 4 nad
> > rovnakým snapshotom ešte raz. Ak znova spadne na kvótu, zastaví sa a ozve sa."
>
> Translation: "Yes, once, now. Codex relaunches the same round-4 review on the same snapshot once more. If it
> fails on quota again, it stops and reports back."
>
> The answer arrived on 2026-10-05 at about 21:28 CEST, much later than the question. The relay probed again
> at 21:28 CEST, and both `claude/claude-opus-5-5[1m]` and `gpt-6-astra` returned `PONG` (liveness only).
>
> This authorizes **exactly one** relaunch of the claude-1 round-04 review on the same product snapshot
> (CLI `906857b`, review snapshot `d8b729a`, skill `dcb7d59`). Keep the review brief focused, to save Claude
> tokens. If it fails on a quota, credit or auth error again, stop and write a new blocking note with the
> verbatim error. Then continue as `IMPLEMENTATION.md` says: review consensus, any fix-up, the
> attended-close request, and the release after the owner confirms the close.

(Original language of the quoted question/option: Slovak; translation as relayed.)

This answer authorizes the relaunch only. It approves no code, waives no gate and grants no close. No
quota, credit or authentication error occurred in this run.

## Refutation attempts

Tags: PRIMARY = my executed check or my direct reading at the locator on `906857b`; SECONDARY = producer
evidence I did not rerun. Every result below is scoped.

### Agreed fixes G1–G9 and reservations R1–R4

| Item | Attempt | Result |
|---|---|---|
| G1 durability | P1 `create` on AppleVirtIOFS and `/tmp`. Knob false, knob true and the nil (legacy) API, against round 03's `inappropriate ioctl for device`. T1 `TestSyncFile*`, `TestQuotaFixupCreateRunOnFilesystem` (actual `parley run --quota-auto-exclude=true/false` on the shared TMPDIR). | **Resolved.** PRIMARY: all six create calls return `err=<nil>`, kickoff written. Raw syncs now go through `fsutil.SyncFile` (`quota/record.go`, `protocol/quota.go`, `app/quota.go`, `runcontrol.go`, `pidlease/lease.go`, `quota/history.go`). |
| G2 authorized changes | P1 `knoboff`/`reinclude`; P3; T1 revision, withdrawal and off-mode tests. | **Partially resolved.** The policy-on `parley quota revise` path works and is bound to committed authority (reading `membership/revision.go:78-181`; T1 PRIMARY pass). Knob-off exclusion in the documented grammar works without a new command. Knob-off return and join now need new markers: R4-MAJOR-2. |
| G3 serialization | P1 `lease` (run-A holds, separate run-B process); P4; T1 `TestLease*`, `TestQuotaCrossProcessLease`, `TestQuotaFixupProjectionWaitAndProcessExclusivity`. | **Resolved.** PRIMARY on both volumes: run-B is refused with `idea driving lock held … lease held: demo/run-A pid=… host=… boot=…`, and acquisition succeeds after release. Lease files live under `.parley-runtime/membership/<hash>/`; the idea directory has none (PRIMARY `ls`). |
| G4 framing | P1 `zadv` (54 cases); P2. | **Wrong-positives resolved; false negative found.** PRIMARY: 54/54 round-03 rows are ineligible, including raw quotes and every mixed JS-error form. P2 shows the retry-wrapper false negative: R4-MAJOR-1. |
| G5 handoff | Reading `consensus_request_signoffs.go:613-626` (passes `ctx`) and `handoff.go:50` (adds `ArtifactPath`); T1 `TestQuotaFixupManualAndInteractiveHandoffKeepLease`. | **Resolved** (PRIMARY: reading plus T1 pass). |
| G6 retained authority | Reading `protocol/quota_retained.go:129-229`; P5; T1 `TestQuotaRetainedVetoOnlyOwnerRulingCanDispose` and `TestQuotaFixupReincludeWithdrawThenAppendSignoff`. | **Resolved structurally.** Rulings bind committed path, commit, blob, digest and quote. Withdrawal needs the original author, `quota-revision`, a matching `reinclusion` batch with owner authority, the author present in the current set, and a later round. Human identity and truth remain out of scope, as documented. |
| G7 receipts | Reading `protocol/quota.go:104-115` (a bad receipt marks the idea pending) and `membership/membership.go:79-182` (no receipt shortcut); T1 `TestQuotaFixupCorruptReceiptRecoveryChecksAllProjections`, `TestQuotaFixupCLIRevisionRecovery`. | **Resolved** (PRIMARY: reading plus T1). Repeated recovery gives one notice and one terminal evaluation. |
| G8 lifetime and targets | Reading the TUI path (`app.go` `defer cancelRun(); driveWait.Wait()` runs before `quotaRelease`), `quota_pipeline.go`, `runner/telemetry.go:87-118` and `quota_target.go`; T1 `TestQuotaFixupTUIExitDrainsDriverBeforeRelease`, `…PipelineBlockDriverOwnsAllStages`, `…AgentsExecCanonicalTargetWithoutIdea`. | **Resolved** (PRIMARY: reading plus T1). The runner-level identity tests passed in T2. |
| G9 diagnostics | Reading `telemetry/provider.go` and the `record.go` Notice gate list; T1 `TestQuotaFixupPreflightBare503AndNoise`, `TestQuotaFixupNoticeNamesRemainingGates`. | **Resolved.** The runner-side 503 test passed in T2. |
| R1 | P5, both volumes: bind a committed answer, then validate it live, after a move to `inbox/archived/` and after deletion. Also a modified archived copy, an unrelated path, a fabricated commit, a wrong sha256, an absent quote and a wrong idea. | **Met.** PRIMARY: the first three validate `<nil>`. The others fail: `user-answer working copy changed`, `committed authority unavailable` (×2), `user-answer digest changed`, `wrong user-answer attribution, idea or verbatim quote` (×2). |
| R2 | P3 (below), plus SKILL.md and ROSTER text. | **Not met for return and join:** R4-MAJOR-2. Met for exclusion in the documented grammar. |
| R3 | P4, both volumes: a dead PID (2147480000) with foreign host, foreign boot, unknown host or unknown boot; then same host and boot. | **Met.** PRIMARY: all four are refused with `only a proven dead local owner may be reclaimed`, owner bytes unchanged. Same host and boot is reclaimed. Windows `definitelyDead` returns false (fail-closed, read only). |
| R4 | Reading `quota_framing.go` (allowlisted header, stack, fields and values; `[Object]` rejected); P1; P2. | **Allowlist met**, and no reconstruction is called native. The AC2 gap stays visible: R4-MAJOR-3. Framing is too strict for real retries: R4-MAJOR-1. |

**P3 raw results** (identical on the shared volume and `/tmp`, `r4extra.log` and `r2-return.log`). The idea
is a knob-off idea created through `CreateIdeaWithQuota`, then `participants: [a, b, c]` plus the line shown:

- this idea's own unbracketed style, `excluded: d — HTTP 429 … (preflight class
  provider-failure:rate-limit) — confirmed 2026-10-03` → `current=[a b c] known=[a b c d] err=<nil>`;
  `RecordManual` revision 1;
- `excluded: [d — unavailable — confirmed 2026-10-04]` → accepted;
- `[d — unavailable — quota exhausted — confirmed 2026-10-04]`, `… confirmed 2026-10-04 via relay`,
  `[d]   # §9.0 — see inbox/…` and `[d - unavailable - confirmed …]` → `manual exclusion d lacks recorded
  owner confirmation`;
- known return of `d` by editing `participants:` only → `manual inclusion d lacks recorded owner
  confirmation`;
- with `included: [d — returned after owner confirmation — confirmed 2026-10-05]` → accepted, revision 2,
  current `[a b c d]`.

My round-03 free-form line (`excluded: d — unavailable; user confirmed 2026-10-04`) is not the
protocol's documented §9.0 grammar. Its refusal is therefore not by itself an R2 breach.

**P2 raw results** (`zretry.log`). "Fixture" is the producer's
`zcode-sdk-retry-source-derived.stderr`. All P2 inputs use complete allowlisted SDK framing, and the
controls (A, B3, E-control) built by the same code are accepted:

```
A source-derived retry fixture (identical retries)         eligible=true  rule="quota.named-reset-ge-60m.v1" reset=2026-10-04T22:14:57.003Z
B aggregated retry, changing countdown, obs=final          eligible=false reason="diagnostic-only: native terminal provider-error provenance unsupported"
B aggregated retry, changing countdown, obs=first          eligible=false (same)
B2 aggregated retry of the two retained countdowns         eligible=false (same)
B3 control: probe-built identical retries                  eligible=true  reset=2026-10-04T22:14:57.003Z
C two top-level records, one clock (obs=first|obs=final)   eligible=false (both)
C1 first record alone (obs=first) / C2 second alone        eligible=true / eligible=true
D two records, real collector clock 2 s apart              eligible=true  rule="quota.named-reset-ge-60m.v1"
E aggregated, backoff-like countdowns 6s/4s/0s, one write  eligible=false
E-control identical nested, same clock                     eligible=true
```

- Inputs B and B2 use the two retained native values: 176930 with "49h 8m 50s", and 176890 with
  "49h 8m 10s". The absolute `reset_at` 2026-10-04T22:14:57.003Z is identical, and `lastError` equals the
  final nested error.
- `lastError` equality itself is checked by fingerprint (`quota_framing.go:187-193`). That part is sound.

### Acceptance criteria AC1–AC21

- **AC1. PASS (scoped).** PRIMARY:
  - the deck and embedded hunks are identical;
  - `cmp` of deck against skill copy gives exit 0;
  - `meta/protocol-changelog.md:1-4` has the entry;
  - T1 ran `TestEmbeddedDefaultMatchesLiveDeck` and it passed;
  - the full packet SHA equals the deck.
  Packets for phases 0, 5 and 8 were not re-rendered. "Not yet in force" is absent, which is consistent
  only because both stages ship together.
- **AC2. NOT MET with native evidence; owner action required (R4-MAJOR-3).**
  - Source-derived positives are eligible (P2 A, C1, C2, D).
  - Both retained native excerpts are rejected alone. PRIMARY: P1 zadv row P1 (the 4-line excerpt); and
    `tail/tail.log`, where the second retained tail, extracted verbatim from
    `codex-1-retained-native-tail.md`, gives `eligible=false`.
  - The SDK retry-wrapper with a realistic changing countdown is rejected (R4-MAJOR-1).
- **AC3. PASS.** PRIMARY, P1 zadv: 54/54 negatives are ineligible. They include:
  - raw and enveloped quotes;
  - `Error:`, `TypeError`, `RangeError`, `SyntaxError`, `AbortError` and `FetchError` lines;
  - an unhandled rejection;
  - a RetryError whose last error is a different one;
  - 503 and 500 records;
  - 30-minute and 59m59s resets;
  - display-only and offset cases;
  - contradictory, past and zone-less resets;
  - hourly and 5-hour allowances with no reset;
  - a generic `quota exceeded` and a bare `credit balance`;
  - exit 0, a valid artifact, a later success, a watchdog and truncation.

  The codex 503s and 401 are diagnostic-only by adapter. The zcode path requires `statusCode: 429`
  (`quota_framing.go`).
- **AC4. PASS.** PRIMARY: `ClassifyQuota` short-circuits every non-zcode adapter (`telemetry/quota.go`),
  and zadv shows codex, claude and kimi ineligible.
- **AC5. PASS.** PRIMARY reading of `quota/quota.go` `Evaluate`; T1 `TestQuotaBatchPermutations`
  (24 orders) and `TestQuotaFloorRolesAndSuccess`.
- **AC6. PASS on both volumes.** PRIMARY: P1 create; T1 actual `parley run` creation on the shared volume.
  Filtering happens before the prompt, `run.created` and the manifest (reading `runcontrol.go:56-58`).
  T2 passed `TestQuotaC1KickoffReplayAndFrozenScope` (runcontrol). The `agent.started` absence is
  SECONDARY.
- **AC7. PASS.** PRIMARY: the shared `ProviderUnavailablePattern` and the T1 preflight bare-503 and noise
  tests. The runner half passed in T2 (`TestQuotaFixupRunnerBare503AndNoise`, `TestQuotaBare503RunnerGate`).
- **AC8. PASS.** PRIMARY: `Evaluate` checks protected roles, and `membership.Roles` reads them. T1
  `TestQuotaStubSurvivorAndProtectedRolesBlockWholeBatch` passed (one blocking note across two attempts).
- **AC9. PASS** for policy-on ideas and documented knob-off markers.
  - PRIMARY: `consensus.Status` takes known signers from history and makes retained obligations
    `TriageBlocked`. `AppendSignoff` accepts only current members.
  - T1 passed the excluded-veto, kickoff-never-known, re-include → withdraw → append and new-identical-veto
    tests.
  - Knob-off return is subject to R4-MAJOR-2.
- **AC10. PASS.** PRIMARY reading of `membership/gates.go` `CheckGates`, used in `Settle`,
  `FinalizeContext` and `driverImplOps.Complete`. T1 passed the review, diversity and strict-gate test.
- **AC11. PASS.**
  - PRIMARY: reading of `DurableWrite` and `CommitBatch` replay-by-id. T1 projection-fault (5 stages),
    corrupt-receipt and truncation tests passed.
  - Driver consumer rebinding (`driver/quota.go`) is by reading, and its test passed in T2.
- **AC12. PASS on both volumes.**
  - PRIMARY: P1 lease and P4.
  - Same-PID competing run: refused by reading `lease.go:84` (`o.PID == os.Getpid()`), and T1
    `TestLeaseSamePIDPartialForeignAndReleaseIdentity` passed.
  - Policy-off, legacy and kickoff-only ideas have no singleton: T1 `TestQuotaLifetimeLockDifferentRunAndOffScope`
    and `TestQuotaFixupLegacyHasNoLifetimeOrProjectionSingleton`.
  - The stale-reaper race (T1 `TestLeaseStaleReaperGenerationClaim`) passed.
- **AC13. PASS.** PRIMARY: reading `ValidateRound` and `runAgent`'s preserved incomplete artifact.
  `OpenReviewRound` now preserves an invalid artifact instead of deleting it (`driver_impl.go`). T1 replay
  and incomplete-preservation tests passed.
- **AC14. PASS.** PRIMARY: `ReadHistory` is strict and never repairs. P1 `reinclude` rows 3–5 show
  contradictory edits blocking.
- **AC15. PASS.**
  - PRIMARY: markers are "automatic quota" and never "confirmed" (`record.go` `Markers`).
  - The notice names the gates.
  - `wait` re-reads workspace status each poll, and `boundaryReached` refuses while pending.
  - T1 passed the read-only status, wait and brief agreement test.
- **AC16. PARTIAL.** Recorded policy reuse, no widening on resume, explicit widening via `revise`, and
  legacy ideas staying on confirmation all pass (PRIMARY: T1 and P1 `reinclude` row 5). Config layering
  passed in T2. The knob-off behavior change in R4-MAJOR-2 is unlisted.
- **AC17. PARTIAL.** PRIMARY by reading: no `agents.toml` write, no timer, and `RelaunchHint` is the reset
  plus 5 minutes. Return on knob-off ideas: R4-MAJOR-2.
- **AC18. PASS** (PRIMARY: zadv quote rows). Obligation capture runs after the decision (`Settle` order).
  The owner-accepted subagent residual is unchanged.
- **AC19. CONDITIONAL, honored.** `status: fix-up-cycle-1`.
- **AC20. PARTIAL.** My T1 run (8 packages) and T2 run (5 packages), both on the shared volume, passed.
  The full suite, race, vet, skill and Windows checks are SECONDARY or unexecuted.
- **AC21. NOT MET.** It needs review consensus on this round, both signoffs, and a NEW owner close answer.

### Membership paths requested by the dispatch

- **Off-mode ordinary confirmations without a new CLI.** Exclusion works (P3). Return and join: R4-MAJOR-2.
- **Manual exclusion and catch-up.** Policy-on: `revise` only (R4-MINOR-2). Knob-off: P3.
- **Known-agent re-inclusion and scope change.** Through `revise` (reading; T1).
- **Veto withdrawal by its re-included author, then a new signoff append.** Passes in T1 for policy-on and
  for knob-off with documented markers.
- **Committed ruling surviving archive or delete.** P5.

## Findings

### [MAJOR] R4-MAJOR-1: The zcode retry wrapper demands identical retry messages and one clock, so real SDK retry exhaustion is never recognized

**Locators.** `internal/telemetry/quota_framing.go:187-198` and `internal/telemetry/quota_zcode.go:40-104`.

**What is wrong.**

- After validating count, reason and `lastError`, the RetryError branch rejects any nested message that
  differs from the last one (`if msg != last { "mixed retry errors" }`).
- The provider message embeds a time-relative countdown ("reset after 49h 8m 50s"). The two retained
  native records of this incident show the countdown changing while the absolute reset stays the same:
  176930 / 49h 8m 50s against 176890 / 49h 8m 10s, `reset_at` 2026-10-04T22:14:57.003Z
  (`codex-1-retained-native-tail.md`).
- Nested records in one dump also share a single receipt clock. So each earlier attempt's
  `retry_after` disagrees with its own absolute `reset_at` by the retry delay. This is PRIMARY by
  reading `quota_zcode.go:40-104`. Executed case C shows the same effect for top-level records on one
  clock.
- The ratified rule requires records to "agree on the exhaustion class and reset within one second"
  (`COOPERATION.md:913`), not identical human text.
- The only accepting retry fixture is generated with `errors:[e,e,e]` (`internal/telemetry/testdata/quota/sdk-framing-harness.cjs:22`):
  one object repeated. A native retry sequence cannot produce that.

**Executed (PRIMARY, P2).** B, B2 and E are ineligible. The identical-message controls are eligible.

**Impact.**

- This fails closed and authorizes nothing wrongly.
- But the only supported recognizer cannot classify the AI SDK's standard `maxRetriesExceeded` report of
  this provider's exhaustion.
- AC2's retry positive is demonstrated only on an impossible input.

**Suggested fix.**

- Compare nested and top-level records by exhaustion class and absolute machine reset (`reset_at`).
- Check each record's internal consistency without assuming one observation instant: `retry_after`,
  header and countdown text should agree with each other.
- Apply the 60-minute threshold to the absolute reset against observation.
- Keep `lastError` fingerprint equality, and keep mixed class/status gating.
- Add fixtures with the two retained values and with backoff-like deltas.
- Any loosening of the 1-second agreement must stay within the owner's text: class plus reset agreement.

### [MAJOR] R4-MAJOR-2: Knob-off return and join need new, undocumented markers — an unlisted third knob-off behavior change (R2)

**Locators.**

- `internal/protocol/quota_manual.go:31-63`, `internal/quota/revision_snapshot.go:38-60` and
  `internal/protocol/quota.go:80-101`.
- `skills/parley-deck/SKILL.md:169` ("Knob-off ideas keep their ordinary recorded confirmations").
- `references/ROSTER_AND_PROTOCOL.md:286`.
- The only definition is `docs/quota-membership.md:61-71`.

**What is wrong.** Every idea created after delivery carries immutable history, even with the policy off.

- A known agent's return must carry `included: [id — reason — confirmed YYYY-MM-DD]`.
  `COOPERATION.md` §9.0 (`:984-985`) requires owner confirmation but defines no such marker, and the skill
  never mentions it.
- A new-identity catch-up join must carry `owner-answer/commit/blob/sha256/quote` frontmatter in the late
  round-1. It also needs a committed owner answer containing exactly `Catch-up join: <id> from round-02`.
  The §5 catch-up rule (`:761`) requires neither.
- Otherwise `QuotaMembers` fails. That blocks consensus status, signoff, finalize and reopen, plus the
  launch gate.

**Executed (PRIMARY, P3).**

- Return without the marker: `manual inclusion d lacks recorded owner confirmation`.
- With the marker: accepted.
- The join requirement is PRIMARY by reading. The producer's
  `TestQuotaFixupOffModeCatchupNeedsCommittedOwnerAnswer` is in T1 and passed.

**Why MAJOR.**

- FINAL §10–§11 allow exactly two knob-off behavior changes, C1 and the bare 503. A third must be the
  owner's decision and must be listed with its own regression expectations.
- R2 reserves exactly this: implement within the existing behavior, or record the deviation for the owner
  instead of absorbing it.
- `IMPLEMENTATION.md` "Deviations from FINAL.md" calls these only "documented confirmation markers" and
  does not list them.

**Fix (either).**

- (a) Record this as a deviation for the owner decision, together with protocol and skill text that
  defines `included:` and the catch-up evidence, a release note, and tests.
- (b) Accept the existing recorded confirmation forms for knob-off return and join.

### [MAJOR] R4-MAJOR-3: Native AC2 acceptance evidence is missing (owner decision required)

**What is missing.**

- The 24,833-byte original was not recovered.
- Both retained native excerpts are rejected by the strict grammar.
- The positives are source-derived, labeled as such, and never presented as native.
- R4-MAJOR-1 shows the realistic retry shape fails.
- Native inspect depth is UNVERIFIED (open question 1): default `console.error` depth can render
  `requestBodyValues.messages` items as `[Object]`, which the grammar rejects.

**Why it matters.** AC2 cannot be marked met from current-tree evidence. I do not waive it.

**What can close it.** Only a native capture or an explicit owner ruling. Obtaining a capture would need an
owner-authorized zcode invocation; zcode is excluded from this idea. The implementer cannot close it alone.

### [MINOR] R4-MINOR-1: The knob-off exclusion grammar is an exact em-dash split with a non-actionable error

**Locator.** `internal/quota/revision_snapshot.go:54`.

**What is wrong.** `ConfirmedInPrompt` accepts only exactly three ` — `-separated fields with a bare
`confirmed YYYY-MM-DD`. A reason that itself contains ` — ` (the protocol's own house style), a note after
the date, or hyphen separators all fail with the same `lacks recorded owner confirmation` (P3).

**Fix.** Name the expected grammar in the error and document it next to §9.0. Optionally, anchor the
parse on the trailing `— confirmed <date>`.

### [MINOR] R4-MINOR-2: Policy-on ideas reject the §9.0 prompt recording, or silently revert it

**Locators.** `internal/protocol/quota.go:96-101` and `:177-181`.

**What is wrong.** For default-on ideas, a human-path §9.0 change recorded in `00-prompt.md` as §9.0
instructs (`COOPERATION.md:889-892`) behaves in one of two ways:

- It is `contradictory quota membership projection`, a hard block that does not name `parley quota revise`.
- Or, when it equals an earlier revision, it is "pending", and the next driver `Before` rewrites
  `participants:` back with no notice. P1 `reinclude` row 2 shows the rewrite `[a, b, c, d] -> [a, b, c]`.

`revise` is the agreed G2 path, but §9.0 does not mention it.

**Fix.** Name `parley quota revise` in the error, and issue an escalation instead of silently rewriting an
owner edit.

### [NIT] R4-NIT-1: Misleading reason when zcode framing is rejected

**Locator.** `internal/telemetry/quota.go:61`.

When `parseSDKCapture` rejects a supported zcode capture (P2 B, C, E), the evidence reason reads
`diagnostic-only: native terminal provider-error provenance unsupported`. Owner notices and logs then
cannot tell a framing or retry rejection apart from an unsupported adapter. Return the parser's reason.

## Prior finding dispositions

All are my own evaluation.

- **Round-03 findings, resolved:**
  - CRITICAL-1 (G1): resolved.
  - MAJOR-2 / MINOR-1 / MINOR-4 (G3): resolved.
  - MAJOR-4 (G5): resolved.
  - MINOR-2 (G6): resolved, structural limits disclosed.
  - MINOR-3 (G8): resolved.
  - MINOR-5 (G7): resolved.
  - NIT-2 and NIT-3 (G9): resolved.
- **Round-03 findings, partly resolved:**
  - CRITICAL-2 (G2): resolved for exclusion. Its knob-off return/join residue is R4-MAJOR-2.
  - MAJOR-1 (G2/G6): resolved for policy-on via `revise` and bound withdrawal. The knob-off residue is
    R4-MAJOR-2, and discoverability is R4-MINOR-2.
  - MAJOR-3 (G4): every demonstrated wrong-positive is resolved. A new false negative is R4-MAJOR-1.
- **Round-03 dispositions I keep:** NIT-1 (completed-drafter protection kept, consensus dismissal) and NIT-4
  (literal AC1 byte equality, with the bootstrap instruction residual). I concur.
- **Reservations:**
  - R1: met.
  - R2: not met for return and join (R4-MAJOR-2).
  - R3: met.
  - R4: allowlist met; the AC2 gap stays visible (R4-MAJOR-3).
- **Round-03 open questions:**
  - Q2 (option A or B): answered by G2.
  - Q3 (TUI and pipeline leads): fixed (G8, T1).
  - Q1 (zcode raw stderr): the raw-quote half is moot now that the allowlist rejects quotes.
- **Round-02 owner rulings** (bounded stderr channel, display clock): still closed by the quoted owner
  ruling. Not reopened.

## Open questions

1. **Native inspect depth (UNVERIFIED, RECALL-level reasoning only).** Which call writes zcode's SDK errors
   to stderr, and with what `util.inspect` depth? In `vendor/zcode.cjs` (sha256 `3e3433d9…`), the only
   `depth:null` uses are unrelated: YAML `LOG_STREAM` and the metrics exporter. If the sink uses the
   default depth 2, real dumps contain `[Object]`, which `quota_framing.go` rejects. This needs evidence
   before any AC2 owner ruling.
2. **Crashed writers.** `membership.RequireStopped` treats any `started.json` without `terminal.json` for
   the idea as a live writer. I found no path that settles a crashed parley process's invocation.
   Whether such an orphan permanently blocks `Before`, `Settle`, `RecordManual` and `quota recover` is
   UNVERIFIED. Please confirm a recovery route.
3. **Manifest before commit.** P3's first attempt had no kickoff run manifest. `RecordManual` committed
   revision 1, then failed in reconcile (`quota manifest run-1: … no such file`), leaving permanent
   pending. The scenario is artificial, because `parley run` writes the manifest. Should manual imports
   validate manifests before committing?
4. **Next steps.** The owner rules on R4-MAJOR-2(a) and R4-MAJOR-3. Fix R4-MAJOR-1, R4-MINOR-1/2 and
   R4-NIT-1, then re-review. Completion still requires both signoffs and a new attended owner close.
