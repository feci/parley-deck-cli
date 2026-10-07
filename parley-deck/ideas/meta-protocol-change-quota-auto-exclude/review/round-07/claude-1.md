---
agent: claude-1
idea: meta-protocol-change-quota-auto-exclude
review-round: 7
date: 2026-10-08
reviewed-commit: e04852ef36846b8ec18f581c4e8ae4f28b0c2262
reviewed-skill-commit: b9596ddd5c37633026a26031f1bccaee76d2622f
---

## Summary

This is a full-scope round-07 review of the complete CLI product diff since FINAL `27e42b8` (121 files,
594,150 B) and the complete skill diff since `a5664d8`, frozen at CLI `e04852e` and skill `b9596dd`.

- **Cycle 4 delivers G18, G18b and G19 as signed.** Owner notice bytes and paths no longer gate status,
  signoff, `Before` or dispatch, before or after a receipt. A bound owner answer can be annotated live or
  archived. Round-1 policy-off plain edits match the 27e42b8 baseline again, and later rounds still
  require catch-up.
- I could not break any of these in my own probes on both volumes, and receipt/history/blob corruption
  still gates.
- **Findings: 0 CRITICAL, 0 MAJOR, 2 MINOR, 2 NIT.**
  - R7-MINOR-1: Windows kickoff directory sync fails, so new ideas cannot be created there. Windows is
    owner-labelled experimental and runtime is a signed deferral, but the release wording understates it.
  - R7-MINOR-2: an aliased (symlinked) `parley-deck` directory breaks the lease path derivation and
    suppresses every notice. Baseline already refuses aliased idea scopes for evidence verification.
  - Neither MINOR is on cycle-4 fix code in a way that re-gates membership. Neither needs a scope or FINAL
    change.
- AC2 stays **NOT MET / owner-waived** for this release (R5-MAJOR-2 accepted and deferred, never fixed or
  PASS). AC21 is not met yet: the final signoffs do not exist.
- This review signs nothing, accepts no code and grants no close. The trajectory is 15 → 6 → 3 → 3 → 4,
  with all four now MINOR/NIT.

## Protocol context and complete coverage

**Packet (PRIMARY, executed 2026-10-07T22:00:5xZ).** Command:
`parley protocol packet --dir . --phase 6 --track deliberation --idea meta-protocol-change-quota-auto-exclude --flag protocol_change --json`.
I ran it with my own build of `e04852e` (`.parley-runtime/claude1-r7/bin/parley`, parley 1.50.0) and with
the installed parley 1.50.0. Both exit 0, write empty stderr and produce byte-identical 26,904 B JSON.

- `context_mode=full`.
- `source_sha256=packet_sha256=73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e`.
- `fallback_reason` is absent (the key is not in the JSON).
- The shadow packet (`6d56607c…`, 83,009 B, 32 omitted blocks) was not used. No optimized or shadow
  context was used.
- The body copy `.parley-runtime/claude1-r7/packet/body.md` is 1,501 lines and 125,862 B. Its hash
  reproduces, and `cmp` against deck `COOPERATION.md` exits 0. The skill
  `skills/parley-deck/references/COOPERATION.md` hashes to the same `73613f95…`.
- Coverage: **the entire body**, read in bounded chunks 1–300, 301–600, 601–900, 901–1200 and 1201–1501.
  Over-long lines 59, 452 and 1002 were displayed in full, and their tails were printed separately.
- Phases 0, 5 and 8 also render `full` with source = packet = `73613f95…`, no fallback, and contain the
  §9.0 rule and the §5 cross-reference.

**Authorities read in full:**

- the newest owner note `…_finish-now.md` (`caafe192…`, newest by mtime and date);
- `…_scope-reset-answer.md`, `…_round05-answer.md` and `…_round06-answer.md`;
- the signed cycle-4 plan `review/consensus.md` (303 lines), including my own signoff and both reservations;
- the producer report `source-context/codex-1-fixup-4-evidence.md` and the release-notes draft;
- `source-context/codex-1-cycle4-platform-ci.md` and the relevant lines of its stored Windows CI log.

**Read in part:**

- FINAL.md: §AC1–AC21 (lines 670–778) in full. I did not reread the rest this round; it is unchanged since
  my full reads in rounds 3–6.
- IMPLEMENTATION.md: lines 757–1038 (round-06 outcome through the cycle-4 producer section) and its
  frontmatter (`status: fix-up-cycle-4`, head `e04852e`). Earlier sections were not reread this round.

**Product diffs (PRIMARY, regenerated).**

- `git diff 27e42b8 e04852e -- . ':!parley-deck/ideas' ':!parley-deck/inbox' ':!parley-deck/runs' ':!graphify-out'`
  is byte-identical to the manifest's `cli-product.diff` (594,150 B, 14,537 lines, 121 files).
- I read **all 32 product chunks in full**, plus the 1,500-line cycle-4 fixup diff and the 436-line skill
  product diff (CHANGELOG, SKILL.md, parley-addon.json, references/COOPERATION.md and
  ROSTER_AND_PROTOCOL.md). No helper agent read anything on my behalf.

File coverage (exhaustive, by chunk):

- 001–004:
  - CHANGELOG, docs/cli-reference, docs/quota-membership;
  - app/agents_exec, app.go, consensus_request_signoffs, driver_consensus, driver_impl, organizer,
    pipeline_cmd, preflight, preflight_liveness, quota.go, quota_consumers.
- 004–008:
  - app quota_cycle3/4/fixup/signoff/surfaces/test, quota_pipeline, quota_revision, quota_signoff, wait;
  - config/runtime(+test);
  - consensus catchup(+test), consensus.go.
- 009–010: consensus/quota_test; driver driver/impl/loop/phasedigest/quota(+test); fsutil sync_darwin(+test).
- 011–016:
  - membership crash, cycle2/3/4 tests, fixup_test, gates, lock, membership(+test), notice, revision(+test);
  - pidlease lease(+test), live_unix, live_windows;
  - protocol/defaults/COOPERATION.md.
- 017–019: protocol participantartifact, quota, quota_manual, quota_retained, workspace; quota/authority.
- 020–022:
  - quota authority_cycle4/authority tests, history(+test), notice, quota(+test), receipt, record, revision,
    revision_cycle4_test, revision_snapshot;
  - quotatest/authority.
- 023–025:
  - runcontrol quota_test, runcontrol;
  - runmanifest;
  - runner acp, failclass, handoff, phase58, quota, quota_target, quota_test, runner, telemetry(+test),
    validation;
  - runplan, runstate, store/events.
- 025–029: telemetry provider, quota, quota_cycle2/framing/zcode(+tests), record, usage.
- 029–031: all telemetry testdata/quota fixtures and the README/harness; deck COOPERATION.md.
- 032: protocol-changelog.

**Execution environment.**

- Every command ran on the macOS host with PATH prefixed by my guard stubs
  (`.parley-runtime/claude1-r7/no-provider-bin`, 21 provider names), then the producer host guard.
- `guard-denied.log` is absent: no provider CLI was reached.
- No provider, participant, capture, worktree, roster or commit operation was performed.
- Before/after hashes of all tracked CLI and skill files are unchanged across my checks.
- Both `agents.toml` files equal the pre-check hashes recorded in `roster-files-before-host-checks.json`.

## User direction

The newest owner note, verbatim (Slovak), is quoted in full in `review/consensus.md` §"Newest user direction
— finish now". Its operative sentence, verbatim:

> "sakra tak to fixni a dokonci a deployni cez vsetky kanaly, taha sa to dlho"
>
> Translation: "Damn, then fix it, finish it and deploy it through all channels, this is dragging on."

Points 4 and 5 of that note bind this review:

> 4. **Fix-up.** Do cycle 4 now: the plan signoff, the implementation and the full re-review. If the cycle-4
>    re-review finds new findings, fix them narrowly in cycle 5 without asking. That is the last cycle the
>    protocol's deliberation cap allows. Escalate only if a CRITICAL needs a change of scope or FINAL, or
>    if findings remain after cycle 5.
> 5. **Close is pre-confirmed.** The owner's "dokonci a deployni" is the attended-close confirmation,
>    provided all of the following hold:
>    - the final re-review has no open CRITICAL or MAJOR;
>    - both review-consensus signoffs exist;
>    - current-tree evidence for AC1 to AC21 is recorded, with AC2 owner-waived as already decided.

I infer no owner gate beyond this note. **Reservations adopted, as delivered:**

- **G18.** Notice inspection is removed from `protocol.InspectQuota` (`quota.go`). `quota.InspectNotice` is
  deleted.
  - `publishNotice` skips ordinary inspection and publication once applied. Before the receipt it preserves
    any regular live or archived copy.
  - Publication-only failures are reported as `publication diagnostic (non-blocking; delivery unconfirmed)`.
  - Receipt comments distinguish an attempt from delivery.
- **G18b.** `BindAuthority` compares the live and archived working copies at bind time only.
  `authorityBytes`/`ValidateAuthority` check commit kind, blob, SHA-256, attribution and quote, never inbox
  bytes.
- **G19 (my smaller route).** `ManualRoundOneReturn(raw)` keys only on the immutable snapshot's
  `status: round-01`. It is applied identically at import (`quota_manual.go:64`) and at replay
  (`revision.go:243`). There is no kickoff field and no legacy limit.
- **G20.** CHANGELOG, `docs/quota-membership.md`, SKILL.md, ROSTER_AND_PROTOCOL.md and the release draft
  list the remaining policy-off differences, the at-most-one notice, attempt versus delivery, and the AC2
  waiver.

## Refutation attempts

Tags: **E** = PRIMARY execution by me; **S** = PRIMARY source reading by me; **P** = SECONDARY, producer
evidence (named) that I did not reproduce.

**My checks** (logs in `.parley-runtime/claude1-r7/checks/`, exits in `results-*.json`):

| Check | Exit | Result |
| --- | --- | --- |
| `go build ./...` | 0 | |
| `go build -o .../bin/parley ./cmd/parley` | 0 | |
| `go vet ./...` | 0 | |
| `GOOS=windows GOARCH=amd64 go build ./...` | 0 | |
| `GOOS=windows go vet` on pidlease, membership, fsutil | 0 | |
| race on 8 quota packages, `-run TestLease\|TestQuotaCycle4\|TestQuotaFixup\|…\|TestQuota` | 0 | 369 PASS, 0 FAIL, 81 cycle-4 PASS |
| focused regex `TestQuotaCycle4\|…\|TestSyncFile` on 5 packages, shared TMPDIR | 0 | 169 PASS / 0 FAIL |
| the same, local `/tmp` | 0 | 169 PASS / 0 FAIL |
| `go test ./internal/runner ./internal/driver -run TestQuota`, shared TMPDIR | 0 | 10 PASS |
| skill `npm test` at `b9596dd` | 0 | 240.8 s |
| `gofmt -l`, 97 changed Go files | | clean |

- The shared-TMPDIR driver run resolves, in my execution, the producer's disclosed `signal: killed`.
- The focused regex includes the unchanged `TestQuotaCycle2NativeCrashWriter`. It passes on both volumes
  on this host, which resolves, in my execution, the producer's sandbox `kern.boottime` failure.
- The full `go test ./... -count=1 -timeout 45m` is recorded under **AC20**.

**AC1, protocol text: PASS.**
- E: packets for phases 0/5/6/8 are full, source = packet = `73613f95…` (see above). Deck, skill reference
  and packet body are byte-identical.
- S: chunk 032 has the changelog entry; chunks 016–018 and 030–031 carry the hunks.
- The embedded default differs only in bootstrap zones; `TestEmbeddedDefaultMatchesLiveDeck` runs in the
  full suite. Both stages ship, so "not yet in force" is correctly absent.

**AC2: NOT MET / owner-waived (round05-answer Q2).**
- S: `QuotaSupport` and the README label it as such. Source-derived positives are not native verification.
- Not a PASS. No capture or parser relaxation was made (S: telemetry chunks are unchanged in cycle 4).

**AC3, negatives: PASS.**
- E (race telemetry ok): `TestQuotaStrictSemantics`, `TestQuotaSuccessAndWatchdogsWin`,
  `TestQuotaZcodeAdversarialFixtures` and the framing tests.
- S: `refineNativeQuota` and `quotaReset` (chunk 025/026).

**AC4, unsupported adapters: PASS.**
- E: `TestQuotaUnsupportedAndAdversarialProvenance`.
- S: support table and `ClassifyQuota` gate to zcode only.

**AC5, floor: PASS.**
- E (race quota ok): `TestQuotaBatchPermutations` (24 permutations) and `TestQuotaFloorRolesAndSuccess`.
- E: app `TestQuotaPreflightWholeBatchAndOneBlock`.

**AC6, C1: PASS.**
- E: `TestQuotaCycle4CLIRoundOnePlainEdits` runs a real `parley run --yes --quota-auto-exclude=false` kickoff
  with a failed local stub probe. The kickoff bytes are checked and the excluded id never becomes known.
  `TestQuotaFixupCreateRunOnFilesystem` passes on both volumes.
- `TestQuotaC1KickoffReplayAndFrozenScope` runs only in the full suite.

**AC7, bare 503: PASS.**
- E: `TestQuotaPreflightReportsOnlyAndBare503` and runner `TestQuotaFixupRunnerBare503AndNoise`.

**AC8, protected roles: PASS.**
- E: `TestQuotaStubSurvivorAndProtectedRolesBlockWholeBatch` (membership, race and focused).

**AC9, filed artifacts: PASS.**
- E (race consensus ok): retained-veto, owner-ruling, withdrawal and kickoff-never-known tests.
- S: `quota_retained.go` (G18b does not weaken the committed-answer check in `quotaDisposition`, which still
  calls `ValidateAuthority`).

**AC10, gates: PASS.**
- E: `TestQuotaReviewDiversityAndStrictGatesOnProspectiveMembers`. S: `gates.go`.

**AC11, durability and replay: PASS.** The R6-MAJOR-1 regression is fixed.

- E: copies of my round-06 notice6/notice7 probes, rebuilt on `e04852e`, on both volumes. The archive,
  delete, annotate-live and annotate-archived modes all show:
  - `Before #1/#2 err=<nil>`;
  - `status exit=0 integrity-gate=false pending=false`;
  - `survivor a consensus signoff exit=0`;
  - `applied-receipt=true`;
  - `round.completed events: 1`.
- E, new probe `r7deck fifo-notice` (local; the shared volume cannot create FIFOs): a FIFO is created at the
  notice path before the receipt. Results:
  - one diagnostic and no hang;
  - `Before err=<nil>` ×2;
  - the FIFO is preserved and the receipt is applied;
  - status 0, signoff 0, and one terminal evaluation.
- E: the cycle-4 matrices pass on both volumes:
  - live/archived/both × annotation/replacement × applied/before-receipt;
  - unsafe symlink, directory, inbox-symlink, inbox-file, archive-symlink and write-failure as diagnostics;
  - receipt symlink/dir and history corruption still gate;
  - a changed receipt still forces full checked replay.

**AC12, serialization: PASS.**
- E, shared and local: `TestQuotaLifetimeLockDifferentRunAndOffScope`, `TestQuotaCrossProcessLease` and
  `TestLease*`. The race run is clean.

**AC13, partial artifacts: PASS.**
- E: `TestQuotaTransitionReplayHistoryAndIncompletePreservation` and runner
  `TestQuotaRunnerProcessFailureTransitionsAfterWritersStop`.

**AC14, history: PASS.**
- E: driver `TestQuotaDriverReconcilesPendingBeforeRebindingEveryConsumer` (shared TMPDIR) and the
  history-immutability assertions in the cycle-4 tests.
- G18b keeps blob, digest and commit checks. `TestQuotaCycle4CLIBoundOwnerAnswerAnnotations` deletes the
  blob object, and status, signoff and `Before` then all refuse.

**AC15, records and surfaces: PASS** under the signed G18 meaning: at most one notice, exactly one when the
destination is safe.
- E: the probes above publish exactly one notice on safe destinations and none on unsafe ones, with a
  diagnostic. `TestQuotaStatusWaitBriefAgreePendingAndAppliedReadOnly` passes.
- Caveat R7-MINOR-2: a symlinked deck root is classed as an unsafe destination for every notice.
- Normative §9.0 still says "One non-blocking owner notice per transition, deduplicated". I read that as an
  upper bound, consistent with G18; it is not a finding.

**AC16, configuration: PASS.**

- E: the plain-edit differential (copied producer probe; baseline built from a verified oracle) on both
  volumes. I re-verified all **491** tracked `internal/`, `cmd/` and `go.mod`/`go.sum` files of
  `.parley-runtime/claude1-r7/oracle-27e42b8` against `git cat-file` of `27e42b8`: 0 mismatches, 0 extra
  files.

  | Requested join | 27e42b8 signoff | Current signoff | After catch-up |
  | --- | --- | --- | --- |
  | round-01, marker kept | exit 0 | exit 0 | not required |
  | round-01, marker removed | exit 0 | exit 0 | not required |
  | round-01, genuinely new id | exit 0 | exit 0 | not required |
  | round-02, marker kept | exit 0 | exit 1, pending catch-up | exit 0 |
  | round-02, marker removed | exit 0 | exit 1, pending catch-up | exit 0 |
  | round-02, new id | exit 0 | exit 1, pending catch-up | exit 0 |

  These round-02 differences are the disclosed G20 differences. History is revision 1, and the
  returned/new id is known only from the new revision.

- E: `TestQuotaCycle4CLIPlainEditBoundaries` covers policy-on, final and closed.
- The config tests (`TestQuotaPresenceAwareDefaults`, `TestQuotaDeckOverridesMachine`) run in the full suite.

**AC17, return: PASS.**
- E: both `agents.toml` files are unchanged across all my checks. There is no timer or rejoin.
- S: `RelaunchHint`, and `TestQuotaPolicyAndHints` passes.

**AC18, integrity: PASS.**
- E: quoted assistant/tool and content negatives.
- S: after G18, notice content can only steer the historical-clarification publication, never membership.

**AC19, completion: CONDITIONAL, honored.** S: IMPLEMENTATION is `status: fix-up-cycle-4`, not complete.

**AC20, checks: PASS on the macOS host for everything listed above.** The full-suite result is in the
addendum at the end.
- Windows: cross-build and cross-vet only. P: the producer's stored Windows CI log shows real Windows test
  failures (R7-MINOR-1). Windows runtime is **not** verified, and the cross-build is not runtime evidence.

**AC21, close: NOT MET.** Both final review-consensus signoffs are absent. The pre-confirmed close
conditions are evaluated after them, not by this review.

**Test shifts round-1 → round-2 (S, every changed test read):**
- `quota_cycle3_test.go` (three tests), `catchup_test.go` `pendingDeclineFixture` and
  `revision_test.go` `TestQuotaFixupOffModeCatchupKeepsPreChangePath` moved only their catch-up, decline
  and pending scenarios to round-02. That is exactly where the signed behavior still applies.
- The round-01 cases are covered positively by new tests (`TestQuotaCycle4CLIRoundOnePlainEdits` and
  `revision_cycle4_test`, which checks every status × marker).
- Gate assertions are unchanged. `TestQuotaCycle3CLIKickoffExcludedReturn` still runs both round-01 and
  round-02.
- The removed `TestQuotaCycle3CorruptNoticeNeverCompletes` encoded the R6-MAJOR-1 defect. Its real-integrity
  cases (receipt symlink, corrupt receipt, missing receipt) live on in `TestQuotaCycle4ReceiptAndHistoryStillGate`
  and `TestQuotaFixupCorruptReceiptRecoveryChecksAllProjections`.
- No later-round regression is hidden.

## Findings

### [MINOR] R7-MINOR-1 — Windows: quota kickoff directory sync fails, so `parley run` cannot create ideas

**Cause (S):**
- `fsutil/sync_other.go:8` makes `SyncFile` a plain `file.Sync()` on every non-darwin OS.
- `quota.WriteKickoff` (`record.go`, after the exclusive write) opens the idea directory and calls
  `fsutil.SyncFile(dir)`. `CreateIdeaWithQuota` (`protocol/quota.go`) also syncs `ideas/`.
- `runTask` always passes `&quotaPolicy`, even with `--quota-auto-exclude=false`, and `newLaunchFunc`
  passes one too.

On Windows, `Sync` on a directory handle returns "Access is denied". So every new `parley run` on Windows
fails in kickoff. In 1.50.0, `CreateIdeaFull` performed no directory sync. `DurableWrite`/`syncDir` and the
lease `publish` hit the same call.

**Evidence (P):** the producer's stored log `.parley-runtime/quota-implementation/review-cycle-4/windows-ci.log`,
l.282–297, shows
`sync C:\…\.parley-runtime\quota-kickoff-…\parley-deck\ideas\…: Access is denied.` and, at l.577–579,
`--- FAIL: TestQuotaC1KickoffReplayAndFrozenScope` and `TestQuotaAutomaticKickoffRecordsAndNotice`. I did
not execute Windows.

**Severity reasoning.** MINOR, not MAJOR:
- the owner's 2026-09-24 decision keeps Windows assets labelled experimental/unvalidated (CHANGELOG 1.50.0);
- the baseline Windows leg already fails;
- CLI winget is held;
- the signed plan defers Windows runtime.

However, "unverified" has become "known broken on the primary command". The release draft says only that
the test suite fails "including quota kickoff directory-sync access errors", which understates the
user-visible consequence.

**Suggested fix (either):**
- (a) Narrow cycle-5 code: treat directory sync as unsupported on Windows, mirroring the darwin ENOTTY
  fallback, with a Windows CI check. Durability semantics there remain unverified, and that should be said.
- (b) Disclosure only: CHANGELOG and the release draft state plainly that on Windows 1.51.0 `parley run`
  cannot create a new idea (kickoff fails with "Access is denied"). Route the fix to the windows-portability
  track.

### [MINOR] R7-MINOR-2 — Aliased deck directory breaks `lockPath` and suppresses every notice

**Cause (S):** `membership/lock.go` `lockPath` resolves `EvalSymlinks(ideaDir)`, then walks up to the first
ancestor named `parley-deck`.

**Evidence (E):** probe `r7deck symlinked-deck`. It runs `parley-deck` → symlink to `real-deck`, then
`CreateIdeaWithQuota` with the default policy, `Acquire`, `CommitBatch`, `Before` and status/signoff.

- **Local `/tmp`:** `Acquire` fails with `idea is outside a deck`. Every lease-requiring step of a default
  (policy-on, mid-idea) idea would fail.
- **Shared volume:** an unrelated ancestor (`…/AI_WORKSPACE/parley-deck`) is taken as the deck. The lease is
  written **outside the workspace root** at
  `/Volumes/My Shared Files/AI_WORKSPACE/.parley-runtime/membership/7f2e…/`. My probe created it; I removed
  only those empty probe-created directories with `rmdir`.
  - In addition, `inspectNotice` rejects the symlinked deck root ("unsafe quota notice directory"). No notice
    is ever published, and the owner only sees a stderr diagnostic. Status, signoff and the receipt were
    otherwise fine.

**Severity reasoning.** MINOR. The baseline already refuses aliased idea scopes for evidence verification
(`internal/app/evidence_refusals.go:35-37` at 27e42b8: "refusal idea scope is missing or aliased"). Aliased
decks are therefore not a supported configuration, and nothing here weakens a gate. But the new code writes
lock state outside the workspace and silently drops notices, instead of refusing clearly.

**Suggested fix:**
- refuse an aliased deck explicitly in `Acquire`/`lockPath` with a clear diagnostic, with no out-of-root
  writes (or derive the deck from the caller's root rather than an ancestor-name search);
- name aliased decks as unsupported in `docs/quota-membership.md`.

### [NIT] R7-NIT-1 — `docs/quota-membership.md` binding sentence lacks "at binding"

Lines 103–106 still say "a contradictory live or archived copy … is rejected". Lines 115–118 then
clarify that this applies at the initial bind only. Add "at binding" to the first sentence.

### [NIT] R7-NIT-2 — Repeated diagnostics and an unneeded `-manual-authority` receipt for new manual revisions

In `membership/notice.go`, a current-era (non-historical) manual revision with an unsafe inbox gets:
- three diagnostics in one call (inspect name, inspect correction, `publishMissingNotice`);
- a `-manual-authority` receipt, although no clarification was ever needed.

This is harmless (publication-only, and it stops after the receipt). Skipping the correction branch when
`err != nil && !historical` would be cleaner.

## Prior dispositions

- **R6-MAJOR-1 (notice integrity gate): fixed and verified (E, both volumes).** AC11/AC15 PARTIAL become
  PASS, with the R7-MINOR-2 caveat.
- **R6-MINOR-1 (round-1 return without marker): fixed by the selected G19 route (E differential).** It
  revises G15's round-1 sub-case to the baseline, as signed. I concur.
- **R6-MINOR-2 (policy-off release wording): fixed.** The lists in CHANGELOG, docs, skill and release draft
  match the executed differential, and no preservation overclaim remains. Precision on Windows remains
  (R7-MINOR-1).
- **Reservation 1 / G18b: delivered and verified.** The owner-append case gates nothing, and blob deletion
  still gates (E, test).
- **R5-MAJOR-2 / AC2: concur.** Owner-accepted, deferred, NOT MET, never PASS.
- **Other deferrals: concur.** D6 accounting, unsupported adapters and container/PID namespaces remain.
  Windows runtime stays deferred, but see R7-MINOR-1.
- **Producer's disclosed failures:** native-crash and shared driver both pass in my host execution (E). The
  producer's sandbox limits are not product defects.

## Open questions

1. **G19 trusts the snapshot's `status:`.** Rolling `status` back to `round-01` admits a later-round plain
   join without catch-up. This is no weaker than baseline or the former marker rule, and the edit stays
   visible in git. Should the organizer guidance say that status is never rolled back for this purpose?
2. For R7-MINOR-1/2, the organizer and owner choose between a narrow cycle-5 fix and an explicit disclosure.
   Both are within the finish-now authority and need no scope or FINAL change.

## Addendum — full suite result (E, same launch)

`go test ./... -count=1 -timeout 45m` (host, guarded PATH, `.parley-runtime/claude1-r7/checks/full.log`):
**exit 0**, 868.4 s, 34 packages `ok`, 0 `FAIL` lines.

- The packages include app, budget, config, consensus, driver, evidence, membership, pidlease, protocol
  (with `TestEmbeddedDefaultMatchesLiveDeck`), quota, runcontrol (with the C1 tests), runner, telemetry,
  trajectory and tui.
- `summary-full.json` shows no CLI or skill source change, an empty guard-denied log, and gofmt clean.

This completes AC20 on the macOS host. The Windows runtime limit in R7-MINOR-1 stands. The final counts
are unchanged: 0 CRITICAL, 0 MAJOR, 2 MINOR, 2 NIT.
