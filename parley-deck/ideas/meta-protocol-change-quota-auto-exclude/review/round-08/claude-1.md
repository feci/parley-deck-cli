---
agent: claude-1
idea: meta-protocol-change-quota-auto-exclude
review-round: 8
date: 2026-10-08
reviewed-commit: 2705a1e850132f74aed2dbb87149df5491bfe94b
reviewed-skill-commit: 99b3f3f9fee161e8e61e585ad6c8e5b1bbbd4d2b
---

## Summary

This is the full-scope round-08 review after the last authorized fix-up cycle (cycle 5). It covers the
complete CLI product diff since FINAL `27e42b8` (124 files) and the complete skill diff since `a5664d8`,
frozen at CLI `2705a1e` and skill `99b3f3f`.

- **Cycle 5 delivers G21–G25 as signed, including all three of my plan reservations.** I could not break
  any of them, in my own executions on both volumes:
  - G21: the Windows limits are disclosed.
  - G22: aliased decks are refused without writes, and the out-of-root lease from R7-MINOR-2 is gone.
  - G23: the bind-time wording and the status-rollback guidance are in place.
  - G24: no unneeded clarification receipt is written.
  - G25: a kickoff notice no longer aborts a created run.
- **Findings: 0 CRITICAL, 0 MAJOR, 1 MINOR, 2 NIT.**
  - R8-MINOR-1: the kickoff floor/role **blocking** escalation still uses a plain exclusive open. With a
    missing or unwritable `inbox/`, `parley run` stops with only `preflight failed: open …` and writes no
    escalation; candidates and arithmetic appear nowhere. This is stage-1 code that I missed in rounds 3–7.
  - R8-NIT-1: the kickoff notice has no replay. A crash between the manifest write and publication loses
    it permanently. Also pre-existing.
  - R8-NIT-2: on an aliased deck with the policy off, every plain `participants:` edit freezes all
    signoffs. The refusal tells the user to disable a policy that is already off, and the disclosure does
    not name plain edits as manual imports.
  - None of the three is a regression in the cycle-5 membership gates. None needs a scope or FINAL change.
- **AC2 stays NOT MET / owner-waived** (round05-answer Q2). It is not a PASS. AC21 is not met: the final
  review-consensus signoffs do not exist.
- The trajectory is 15 → 6 → 3 → 3 → 4 → 3, with no CRITICAL or MAJOR since round 06. Under finish-now
  point 4, findings remain after cycle 5. What follows from that is the owner's call. This review signs
  nothing, accepts no code and grants no close.

## Protocol context and coverage

**Packet (PRIMARY, executed 2026-10-07T23:28:12Z).** Command:
`parley protocol packet --dir . --phase 6 --track deliberation --idea meta-protocol-change-quota-auto-exclude --flag protocol_change --json`.
I ran it with my own build of `2705a1e` (`.parley-runtime/claude1-r8/retry2/bin/parley`) and with the
installed parley 1.50.0. Both exit 0 with empty stderr and produce byte-identical 26,904 B JSON.

- The installed binary's first output file read back as 0 bytes on the shared volume. A re-run, and a
  later re-read of the first file, are byte-identical. I treat it as virtiofs read-after-write lag and
  disclose it.
- `context_mode=full`.
- `source_sha256=packet_sha256=73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e`, the full
  expected source, unchanged.
- `fallback_reason` is absent (the key is not in the JSON). The shadow packet (`6d56607c…`, 83,009 B,
  32 omitted blocks) was not used. No optimized or shadow context was used.
- Body: the emitted live file `.parley-runtime/protocol-packets/full-phase6-deliberation-73613f95….md`,
  1,501 lines, 125,862 B. Its hash reproduces.
  - `cmp` against deck `COOPERATION.md` exits 0, and `cmp` against skill
    `skills/parley-deck/references/COOPERATION.md` exits 0.
  - I read **the entire body** in bounded chunks 1–300, 301–600, 601–900, 901–1200 and 1201–1501.
  - The over-long lines 59, 452 and 1002 (1,957 / 2,894 / 2,188 characters) displayed in full, and I
    printed their last 400 characters separately to confirm the tails.
- Phases 0, 5 and 8 also render `full`, source = packet = `73613f95…`, with no fallback. Each contains the
  §9.0 rule, the §5 cross-reference, the Phase 0 `quota_auto_exclude: false` line and the Phase 5
  designee cross-reference.

**Authorities read in full:**

- `…_finish-now.md` (`caafe192…`, newest), `…_scope-reset-answer.md`, `…_round05-answer.md` and
  `…_round06-answer.md`;
- FINAL.md (`f90577f1…`, unchanged), all 903 lines including AC1–AC21;
- `review/consensus.md` (`e5b56a7f…`), including my full signed cycle-5 block;
- my raw `review/round-07/claude-1.md` (`4705ad7c…`).

**Read in part:** IMPLEMENTATION.md (`862ac5f4…`): the frontmatter, the summary/checklist, and the
cycle-5 sections (lines 1057–1146).

**Producer material read** (SECONDARY, not verdicts): `codex-1-fixup-5-evidence.md`,
`codex-1-fixup-5-host-evidence.md`, `codex-1-cycle5-platform-ci.md`, `codex-1-release-notes-draft.md`
(`437f04d8…`), the `windows-portability` handoff note, `fixup-5/{focused-initial,focused-corrected,regression-local}.log`
and `fixup-cycle-5-20261008/host-summary.json`.

**Diffs (PRIMARY, regenerated).**

- `git diff 27e42b8 2705a1e -- . ':(exclude)parley-deck/{ideas,inbox,runs}/**' ':(exclude)graphify-out/**' ':(exclude).parley-runtime/**'`
  is byte-identical to the manifest's `cli-product.diff`: sha256 `3b625ac4…`, 613,727 B, 15,110 lines,
  124 files (85 added, 39 modified).
- `git diff a5664d8 99b3f3f` (skill) is byte-identical to `skill-product.diff`: sha256 `fe1b7a4c…`,
  47,402 B, 5 files.
- All 37 manifest chunk hashes verify.
- **I read every chunk in full myself.** That is 31 CLI product chunks, the 2 CLI cycle-5 chunks
  (`e04852e..2705a1e`, 11 files), the 3 skill product chunks and the skill cycle-5 chunk. No helper agent
  read anything for me.

CLI coverage by chunk (paths under `internal/` unless noted):

- 001–002: CHANGELOG, docs/cli-reference, docs/quota-membership; app agents_exec, app.go,
  consensus_request_signoffs.
- 003–007: app driver_consensus, driver_impl, organizer, pipeline_cmd, preflight, preflight_liveness,
  quota.go, quota_consumers, quota_cycle3/4/fixup/signoff/surfaces tests, quota_test, quota_pipeline,
  quota_revision, quota_signoff.
- 008–009: app/wait; config runtime(+test); consensus catchup(+test), consensus.go, quota_test; driver.go.
- 010: driver impl, loop, phasedigest, quota(+test); fsutil sync_darwin(+test); membership crash,
  cycle2_test.
- 011–015: membership cycle3/4/5 tests, fixup_test, gates, lock, membership(+test), notice, revision(+test).
- 015–016: pidlease lease(+test), live_unix, live_windows; protocol/defaults/COOPERATION.md.
- 017–019: protocol participantartifact, quota, quota_manual, quota_retained, quota_scope, workspace.
- 019–022: quota authority(+cycle4/tests), history(+test), notice, quota(+test), receipt, record,
  revision, revision_cycle4_test, revision_snapshot; quotatest/authority.
- 022–025: runcontrol quota_cycle5_test, quota_test, runcontrol; runmanifest; runner acp, failclass,
  handoff, phase58, quota, quota_target, quota_test, runner, telemetry(+test), validation; runplan;
  runstate; store/events.
- 025–030: telemetry provider, quota, quota_cycle2_test, quota_framing(+test), quota_test,
  quota_zcode(+test), record, usage; every telemetry/testdata/quota fixture plus the README and harness.
- 030–031: parley-deck/COOPERATION.md and parley-deck/meta/protocol-changelog.md.

Skill coverage: CHANGELOG, SKILL.md, parley-addon.json, references/COOPERATION.md and
references/ROSTER_AND_PROTOCOL.md.

**Execution environment.**

- Every command ran on the macOS host. PATH started with my guard stubs
  (`.parley-runtime/claude1-r8/no-provider-bin`), then the producer guard
  (`fixup-5/host/no-provider-bin`).
- `claude1-r8/guard-denied.log` is absent, so no provider CLI was reached. Probes launched only local stub
  scripts.
- No provider, participant, capture, worktree, roster, commit, merge or release operation was performed.
  OpenViking was not consulted; this review rests on local canonical files.
- Tracked CLI and skill files are unchanged across every check (`summary-*.json`: no source changes).
- Both `agents.toml` hashes are unchanged before and after my work (`f8cc2ab5…` machine, `f52a0a77…` deck).
- No `.parley-runtime` exists at `/Volumes/My Shared Files/AI_WORKSPACE/` or
  `…/AI_WORKSPACE/parley-deck/`.

**Earlier attempts of this same step.** They died on provider 502s before producing any artifact, and
left check logs under `claude1-r8/checks` and `claude1-r8/retry1/checks`.

- Attempt 0's shared-TMPDIR runs had three failures:
  - `TestQuotaCycle4CLIRoundOnePlainEdits/return-kept` (a stub-written artifact read back with an empty
    `agent:` field);
  - `pidlease` `signal: killed`;
  - a driver test binary `exec format error`.
- Retry 1 and this attempt (retry 2) passed the same commands.
- I attribute these to shared-volume (virtiofs) coherence, like the empty packet read above. That is
  reasoning, not proof, and I disclose it rather than drop it.

## User direction

The newest owner note, verbatim (Slovak), from `…_finish-now.md` (2026-10-07):

> "sakra tak to fixni a dokonci a deployni cez vsetky kanaly, taha sa to dlho"
>
> Translation: "Damn, then fix it, finish it and deploy it through all channels, this is dragging on."

Points 4 and 5 bind this review, verbatim:

> 4. **Fix-up.** Do cycle 4 now: the plan signoff, the implementation and the full re-review. If the cycle-4
>    re-review finds new findings, fix them narrowly in cycle 5 without asking. That is the last cycle the
>    protocol's deliberation cap allows. Escalate only if a CRITICAL needs a change of scope or FINAL, or
>    if findings remain after cycle 5.
> 5. **Close is pre-confirmed.** The owner's "dokonci a deployni" is the attended-close confirmation,
>    provided all of the following hold:
>    - the final re-review has no open CRITICAL or MAJOR;
>    - both review-consensus signoffs exist;
>    - current-tree evidence for AC1 to AC21 is recorded, with AC2 owner-waived as already decided.

The AC2 authority, `…_round05-answer.md` Q2, verbatim: Selected **"Vydať s obmedzením + follow-up
(Recommended)"**. Translation: "Release with the limitation plus a follow-up."

I infer no further owner gate and no obligation to pass. §15.5 role concentration: codex-1 organizes and
implements, and supplies no independent verdict here.

## Refutation attempts

Tags:

- **E** = PRIMARY execution by me;
- **S** = PRIMARY source reading by me;
- **P** = SECONDARY producer evidence, named, not reproduced.

**My checks** on the frozen inputs (`.parley-runtime/claude1-r8/retry2/checks/`, exits in `results-*.json`):

| Check | Exit | Result |
| --- | --- | --- |
| `go build ./...` / binary build | 0 / 0 | |
| `go vet ./...` | 0 | |
| `GOOS=windows GOARCH=amd64 go build ./...`; `GOOS=windows go vet` on pidlease, membership, fsutil | 0 / 0 | compile only, not runtime proof |
| `go test -race` on 9 packages, `-run TestLease\|TestQuotaCycle5\|…\|TestQuota` | 0 | 411 PASS, 0 FAIL, 39 cycle-5 PASS lines |
| focused regex `TestQuotaCycle5\|Cycle4\|Cycle3\|Cycle2\|Fixup\|Lifetime\|CrossProcess\|TestLease\|TestSyncFile`, 8 packages, shared TMPDIR | 0 | 244 PASS, 0 FAIL, 0 SKIP |
| the same, local `/tmp` | 0 | 244 PASS, 0 FAIL, 0 SKIP |
| `go test ./internal/runner ./internal/driver -run TestQuota`, shared TMPDIR | 0 | 10 PASS |
| `go test ./internal/protocol -run TestEmbeddedDefaultMatchesLiveDeck -v` | 0 | PASS |
| skill `npm test` at `99b3f3f` | 0 | 399 Node, 54 Python, 6 add-on manifests ok, 126.7 s |
| **full `go test ./... -count=1 -timeout 45m`** | **0** | 616.8 s, 34 packages `ok`, 0 `FAIL` |
| `gofmt -l` on all 100 changed Go files | | clean |

The read-only kickoff-inbox case really executes on both volumes (0 SKIP).

**My probes**, injected with `go test -overlay` from `retry2/probe/`, never written into the tree.
Logs are in `retry2/logs/`. Each ran on both the shared volume and local `/tmp`.

- P-A, `TestClaude1R8KickoffCLINotice`: a real `parley run` (app `Run`) with stub agents a–d, where `d`'s
  readiness probe carries eligible synthetic quota evidence. Six inbox modes.
- P-B, `TestClaude1R8KickoffFloorBlockInbox`: a 3→1 kickoff floor block with the inbox present, missing
  and read-only.
- P-C, `TestClaude1R8AliasCLI`:
  - a symlinked deck with the default policy, and with the policy off plus a round-1 plain edit;
  - a workspace ancestor alias;
  - a symlinked idea directory on an existing scoped idea.
- P-D: my round-07 probes, rebuilt on `2705a1e`: notice6, notice7, fifo-notice, symlinked-deck and
  plain-edits. The baseline comparison reuses the round-07 `baseline-plain-edits` binary, built then from
  an oracle verified file-by-file against `27e42b8` (491 files). I did not re-verify the oracle this round.

**Cycle-5 fixes, each judged independently.**

- **G21, Windows disclosure: verified.**
  - S: the CLI CHANGELOG, the skill CHANGELOG and the release draft say plainly that Windows 1.51.0
    cannot create ideas, including policy-off. They name scoped driving/signing and manual-import,
    owner-revision and transition failures, keep Windows experimental and hold CLI winget.
  - S: the handoff `inbox/codex-1-to-all_windows-portability_quota-durable-sites.md` lists the five
    sites. It distinguishes the local `windows-portability` branch at `3526b82` (2026-09-28,
    `IMPLEMENTATION` in progress) from `origin/windows-portability` at `9a60390`, both confirmed by
    `git rev-parse` on local refs. It promises no fix.
  - E on the stored log `review-cycle-5/windows-ci.log`: sha256 `5933115e…`, 126,208 B, 327 `--- FAIL`
    lines, 142 "Access is denied".
    - All six new cycle-5 test functions fail there only with the kickoff `sync … Access is denied`.
    - `TestQuotaCycle5CreationScope/deck/mid` and `/ideas/mid` (refusal before writes) and all `legacy`
      subtests are absent from the failures. So no new Windows-specific alias misclassification is
      visible.
    - The unchanged `syscall.Mkfifo` test-compile failure is present.
  - P: the GitHub run URLs. I did not access GitHub.
- **G22, alias refusal: verified.** There is a disclosure-precision residue (R8-NIT-2).
  - S: `QuotaLeaseRoot` derives the root lexically, resolves only the workspace root, and then requires
    `parley-deck`, `ideas` and `<idea>` to resolve to themselves. `lockPath` keys on
    `ideas/<base>` under that root and never walks resolved ancestors.
  - S: callers are `acquireScoped`, `requireLease`/`CheckLease`, `projectionLock`, `RecordManual`,
    `Revise`, and `IntegrityBlock` when scoped. `Acquire` still no-ops for legacy, off and kickoff-only.
  - E (P-C): a symlinked deck with the default policy gives `run create failed: quota leases do not
    support symlinked deck or idea scopes…`.
    - The workspace tree is byte-unchanged.
    - On the shared volume the misleading outer `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck`
      exists, yet `…/AI_WORKSPACE/.parley-runtime/membership` was not created. That was R7-MINOR-2's
      out-of-root write.
  - E (P-C): with the policy off on the same aliased deck, the run succeeds and creates no
    `.parley-runtime/membership`.
  - E (P-C): with an ancestor alias, `parley run --dir <alias>` exits 0. The lease lives under the
    physical root, nothing appears beside the alias, and a rival through the physical path is refused
    while the alias path holds the lease.
  - E (P-C): with a symlinked idea directory, `consensus signoff` and `quota recover` refuse, status exits
    0, and nothing is written in the deck or the idea.
  - E (P-D): the round-07 symlinked-deck probe now refuses at `CreateIdeaWithQuota` on both volumes.
  - Missed entry paths, examined:
    - pipeline blocks use `parley-deck/ideas/<slug>__<block>`, so the layout check holds (S);
    - the TUI launch path takes its lease inside the runner's `quotaBefore` (S);
    - `quota.IdeaRoot` walks lexically to the idea's own `parley-deck`, so it is not the R7-MINOR-2
      pattern (S).
  - P: the producer's disclosed first failure (`focused-initial.log`) came from a "legacy" fixture that
    deleted the kickoff record but kept scoped prompt fields. Refusing that state is the correct
    fail-closed result, and the corrected true-legacy fixture is right. I concur.
- **G23, wording: verified (S).**
  - `docs/quota-membership.md` uses exactly the split I asked for: "At binding, a contradictory live or
    archived copy is rejected; a changed digest, wrong idea/author, missing object or fabricated path is
    rejected at binding and on every later read."
  - "Do not roll status back to round-01 to bypass catch-up" appears in the docs and in the skill
    ROSTER_AND_PROTOCOL.
  - COOPERATION.md is unchanged (`73613f95…`).
- **G24, unneeded clarification receipt: verified.**
  - S: in `membership/notice.go`, `ReadApplied` errors still return first. A failed inspection without
    historical evidence returns before correction inspection, publication or receipt.
  - E: `TestQuotaCycle5ManualClarificationUnsafePath` (race, shared, local) asserts:
    - 2/1 diagnostics per call (unapplied/applied);
    - no receipt, and no write through the aliased inbox;
    - a symlinked real receipt still errors;
    - the historical label is corrected exactly once after inbox repair;
    - a malformed receipt is repaired by checked replay.
  - S, by design: a historical label whose correction destination is a non-regular object still records
    an attempt receipt with a diagnostic. That is consistent with the signed attempt-not-delivery
    semantics.
- **G25, kickoff notices: verified.**
  - S: `runcontrol.Create` now calls `quota.PublishNotice` after `run.created` and the manifest and only
    reports a diagnostic. This is the same helper as mid-idea publication.
  - E (P-A, all six modes on both volumes; Created idea printed):
    - `parley run` exits 0 and the transition candidate is `d`;
    - prompt, `run.created` and the manifest are all `[a b c]`, `agent.started` is `[a b c]`, and `d`
      never starts;
    - exactly 1 notice and 0 diagnostics for present/missing, with the missing inbox created;
    - 0 notices and exactly 1 diagnostic for read-only, inbox-symlink, inbox-file and archive-symlink;
    - nothing is written through an alias, and the inbox-file owner bytes are preserved;
    - `status` exits 0 and shows `automatic exclusion: d`;
    - `Acquire`+`Before` under the kickoff run return nil, and the notice count is unchanged after replay.
  - Observation (not a finding): `durableMkdir` recreates a missing inbox with mode 0700, where
    `parley init` uses 0755. Git does not track directory modes.

**AC1–AC21.**

- **AC1, protocol text: PASS.**
  - E: phase 0/5/6/8 packets are full, source = packet = `73613f95…`; the drift test passes; the deck,
    the skill reference and the packet body are byte-identical.
  - S: chunk 031 has the changelog entry, and chunks 016–017 and 030–031 carry the hunks.
  - Both stages ship, so "not yet in force" is correctly absent.
- **AC2: NOT MET / owner-waived.** Not a PASS.
  - S: `QuotaSupport` and the testdata README say native-positive evidence is NOT MET and owner-waived.
  - S: the cycle-5 diff touches no telemetry file, so there is no capture and no parser relaxation.
  - Source-derived positives are not native verification.
- **AC3, negatives: PASS.**
  - E (race, telemetry ok): `TestQuotaStrictSemantics`, `TestQuotaSuccessAndWatchdogsWin`,
    `TestQuotaZcodeAdversarialFixtures`, the framing allowlist tests and the 503 noise tests.
  - S: `refineNativeQuota`, `quotaReset` and `classifyZcodeQuota`.
- **AC4, unsupported adapters: PASS.** E: `TestQuotaUnsupportedAndAdversarialProvenance`. S:
  `ClassifyQuota` accepts zcode only.
- **AC5, floor: PASS, except R8-MINOR-1.**
  - E: `TestQuotaBatchPermutations` (24 permutations), `TestQuotaFloorRolesAndSuccess` and
    `TestQuotaPreflightWholeBatchAndOneBlock`.
  - E (P-B): a 3→1 kickoff with the inbox present applies nothing, creates no idea, exits 1 and writes
    one escalation listing `[b c]`.
  - With the inbox missing or read-only, no escalation is written (R8-MINOR-1).
- **AC6, kickoff (C1): PASS.**
  - E (P-A): the excluded id is absent from the first prompt, `run.created`, the manifest and
    `agent.started`.
  - E: the C1 and `--yes` tests pass: `TestQuotaC1KickoffReplayAndFrozenScope` in the race and full runs,
    and `TestQuotaCycle4CLIRoundOnePlainEdits` (a real `run --yes` kickoff) in the race, shared and local
    runs.
  - E: standalone preflight reports only (`TestQuotaPreflightReportsOnlyAndBare503`).
- **AC7, bare 503: PASS.** E: preflight and runner bare-503 tests. S: the shared
  `ProviderUnavailablePattern`.
- **AC8, protected roles: PASS.** E: `TestQuotaStubSurvivorAndProtectedRolesBlockWholeBatch` and
  `TestQuotaFloorRolesAndSuccess` (designee, pin, drafter, global-default fall-through).
- **AC9, filed artifacts: PASS.** E (race, consensus ok): retained-veto, owner-ruling, withdrawal and
  kickoff-never-known tests. S: `quotaDisposition` still calls `ValidateAuthority`.
- **AC10, gates: PASS.** E: `TestQuotaReviewDiversityAndStrictGatesOnProspectiveMembers`. S: `gates.go`.
- **AC11, durability and replay: PASS.**
  - E (P-D, both volumes):
    - notice6 archive/delete and notice7 annotate-live/archived: `Before` ×2 nil, status without a
      gate, survivor signoff 0, applied receipt, one `round.completed`;
    - fifo-notice (local; the shared volume cannot create a FIFO): one diagnostic, no hang, FIFO
      preserved.
  - E: the cycle-4 matrices and the receipt/history gates pass.
- **AC12, serialization: PASS.**
  - E: `TestQuotaLifetimeLockDifferentRunAndOffScope`, `TestQuotaCrossProcessLease`, `TestLease*` and
    `TestQuotaCycle5AncestorAliasSharesLease`.
  - E (P-C): cross-path serialization through an alias.
- **AC13, partial artifacts: PASS.** E: `TestQuotaTransitionReplayHistoryAndIncompletePreservation` and
  `TestQuotaRunnerProcessFailureTransitionsAfterWritersStop`.
- **AC14, history: PASS.** E: driver `TestQuotaDriverReconcilesPendingBeforeRebindingEveryConsumer` (shared
  TMPDIR) and `TestQuotaLaterRunReadsHistoryAndContradictionsBlock`.
- **AC15, records and surfaces: PASS under the signed G18 meaning** (at most one notice; exactly one when
  the destination is safe), with the R8-NIT-1 crash-window caveat.
  - E (P-A): kickoff notices.
  - E: `TestQuotaStatusWaitBriefAgreePendingAndAppliedReadOnly`.
- **AC16, configuration: PASS.**
  - E (P-D): the plain-edit differential, current vs the `27e42b8` baseline binary, on both volumes,
    matches round 07 exactly:
    - round-01 marker-kept, marker-removed and new-id joiner signoffs exit 0, as on baseline;
    - round-02 joiners exit 1 with pending catch-up and exit 0 after catch-up. These are the disclosed
      G20 differences.
  - E: the presence-aware and deck-override config tests pass in the full suite.
- **AC17, return: PASS.** E: both `agents.toml` files are unchanged. S: there is no timer or rejoin, and
  `RelaunchHint` is the reset plus 5 minutes, or none when the reset is unknown.
- **AC18, integrity: PASS.** E: the quoted assistant/tool and content negatives. S: evidence comes only
  from supervisor-captured terminal streams.
- **AC19, completion: honored.** S: IMPLEMENTATION is `status: fix-up-cycle-5`, not complete.
- **AC20, checks: PASS on the macOS host** (the table above).
  - Windows: cross-build and cross-vet only, plus the stored CI log (P/E-log). Windows runtime is known
    broken, as disclosed.
- **AC21, close: NOT MET.** `review/consensus.md` holds the cycle-5 plan. No final review-consensus
  signoffs exist. This review evaluates no close condition.

## Findings

### [MINOR] R8-MINOR-1 — Kickoff floor/role escalation is lost when `inbox/` is missing or unwritable

**Cause (S).** `writeQuotaBlock` (`internal/app/quota.go:26–27`) writes the kickoff `blocking: yes`
escalation with a plain `os.OpenFile(…, O_CREATE|O_EXCL)` into `parley-deck/inbox/`.

- On error, `preflight` returns `report, 1, err`.
- `runTaskPreflight` then prints only `preflight failed: <err>`. The gate detail, with candidates and
  arithmetic, is never printed on this path.
- By contrast, mid-idea `membership.Block` uses `quota.DurableWrite`, which creates a missing inbox.
  Cycle-5 G25 made kickoff *notices* do the same. The kickoff *escalation* was not converted.
- The code dates from stage 1 (`3aa05cf`, 2026-10-04). I missed it in rounds 3–7.

**Evidence (E, P-B, both volumes).** Participants a, b, c; `b` and `c` carry eligible synthetic quota
evidence (3→1).

| `inbox/` | Exit | Ideas | Escalations | User-visible output |
| --- | --- | --- | --- | --- |
| present | 1 | 0 | 1, listing `[b c]` | `preflight: hard stop — fewer than 2 participants after exclusion (§1 non-solo).` |
| missing | 1 | 0 | 0 | `preflight failed: open …/inbox/parley-to-user_quota-kickoff-….md: no such file or directory` |
| read-only | 1 | 0 | 0 | `preflight failed: open …: permission denied` |

**Impact.**

- What still holds: nothing is applied and no idea is created, so the gate fails closed.
- What fails: FINAL §6 / AC5 require "one blocking escalation lists both candidates and the arithmetic".
  With a missing inbox, as in a fresh clone of a deck whose `inbox/` was empty and so untracked, the
  escalation is not issued. The candidates and arithmetic are recorded nowhere, and the message does not
  say this was a quota floor block.
- Reachability is narrow: it needs eligible kickoff evidence (today zcode only, and AC2 says it may be
  inert), a floor break, and a missing or unwritable inbox.

**Severity: MINOR, not MAJOR.** The gate fails closed, no state is corrupted, and the precondition set is
rare. It is still a refutation of an explicit AC5 obligation in a real environment.

**Suggested fix (any one):**

- (a) Mirror `membership.Block`: return if the escalation already exists, otherwise
  `quota.DurableWrite(path, text, true)`. The body embeds the date, so deduplicate by path first.
- (b) On any write failure, also print the decision (candidates and arithmetic) to stderr before exiting
  non-zero.
- (c) Disclose it and route the fix to a follow-up.

A missing/read-only inbox test belongs with (a) or (b).

### [NIT] R8-NIT-1 — The kickoff notice has no replay; a crash after the manifest loses it permanently

**Cause (S).**

- `runcontrol.Create` (`runcontrol.go:114–118`) attempts the kickoff notice once, after `run.created` and
  the manifest. It is the only kickoff `Notice()` publication site in the tree (`grep`).
- `membership.Before`/`reconcileLocked` replays only `h.Batches`, and kickoff transitions have no
  `quota-applied` receipt.
- So a crash or kill between the manifest write and publication leaves zero notices, even on a safe
  destination.
- The kickoff record, the marker, status, wait and the brief still show the exclusion (E, P-A:
  `automatic exclusion: d`).
- Pre-existing since stage 1. Not executed: I did not inject a crash.

**Severity: NIT.** The window is milliseconds and other surfaces remain. It still departs from the signed
AC15 meaning ("exactly one when the destination is safe"), and from the receipt-replayed mid-idea path.

**Suggested fix (either):**

- publish right after `CreateIdeaWithQuota` returns (the helper preserves existing copies), or give
  kickoff transitions a one-attempt receipt replayed by `Before`;
- or disclose it.

### [NIT] R8-NIT-2 — Aliased deck with the policy off: plain edits freeze every signoff, and the remedy text misleads

**Evidence (E, P-C, both volumes).** Symlinked deck, `--quota-auto-exclude=false`; the run succeeds.

| Step | Exit | Output |
| --- | --- | --- |
| signoff by `a` before any edit | 0 | |
| round-1 plain edit `[a, b, c]` → `[a, b, c, d]`, then signoff by existing member `b` | 1 | `quota leases do not support symlinked deck or idea scopes: …/parley-deck; use a physical deck or disable quota_auto_exclude for ordinary driving` |
| signoff by joiner `d` | 1 | the same |
| `status` | 0 | |

**Cause (S).**

- `RecordManual` takes the revision lease, so the alias refusal fires for any captured manual revision.
- That includes the round-1 plain joins the CHANGELOG calls unchanged ("Round-1 policy-off joins/returns
  use plain edits as before"), and an ordinary §9.0 `excluded: … — confirmed …` edit.
- Every signer, driver tick (`Before`) and `agents exec` catch-up import then fails until the deck is
  physical.
- `agents exec` launches the agent (`agents_exec.go:130`) before its import fails (`:152–162`; S, not
  executed on an aliased deck).
- In 1.50.0 the signoff path had no alias-dependent code (S, `git grep EvalSymlinks 27e42b8` shows none in
  consensus or workspace).

**Why NIT.** The behavior matches my own signed G22(a), which named manual-import leases as refusal points.
The docs, CHANGELOG and release draft do say "manual/owner revisions still require a physical scope".
The defects are precision only:

- the refusal's second remedy ("disable quota_auto_exclude") cannot apply to a policy-off idea;
- the disclosure does not say that plain `participants:` edits and confirmed exclusions are such manual
  imports, or that they block all signers;
- secondary: with the default policy, the refusal happens only after the readiness pings (provider
  calls) have run, because `runTask` calls preflight before `runcontrol.Create` (S).

**Suggested fix.** Tailor the refusal text on the manual and revision paths: "restore a physical deck
path; the participants edit stays pending". Name plain edits and confirmed exclusions in the aliased-deck
bullet. Optionally check the deck alias before preflight when the resolved policy is enabled mid-idea.

## Prior dispositions

- **R7-MINOR-1 → G21 (disclosure-only): concur, resolved.** The wording matches the executed CI evidence.
  The branch routing is explicit and promises nothing. Windows runtime stays deferred to the unmerged
  `windows-portability` track.
- **R7-MINOR-2 → G22: concur, resolved** (E: no out-of-root lease, refusal without writes, ancestor
  aliases and serialization kept). The residual precision issue is R8-NIT-2.
- **R7-NIT-1 → G23: concur, resolved.**
- **R7-NIT-2 → G24: concur, resolved.**
- **Reservation 3 → G25: concur, resolved** (E through real `parley run`, six modes, both volumes).
- **Round-07 open question 1 (status rollback): answered** by the G23 guidance, with no new gate. I concur.
- **R5-MAJOR-2 / AC2: concur.** Owner-accepted and deferred to
  `quota-zcode-native-exhaustion-capture` (present at `2705a1e`). NOT MET, never PASS.
- **Other deferrals: concur.** D6 accounting, unsupported-adapter provenance, container/PID namespaces,
  and Windows crash settlement.
- **Producer's disclosed initial failure: concur.** It was a fixture error, not a product defect (see G22).

## Open questions

1. R8-MINOR-1 and both NITs are pre-existing or disclosure-level, and none needs a scope or FINAL change.
   No cycle 6 is authorized. Does the owner want them fixed in a follow-up, disclosed in the release
   notes, or both? This is the owner's call under finish-now point 4. This review does not presume it.
2. Should the release notes add one line on R8-MINOR-1, namely that a fresh clone needs an `inbox/`
   directory before an automatic kickoff floor block can record its escalation, until it is fixed?

## Result, path and limits

- **File:** `parley-deck/ideas/meta-protocol-change-quota-auto-exclude/review/round-08/claude-1.md`.
- **Counts:** 0 CRITICAL, 0 MAJOR, 1 MINOR, 2 NIT. AC2 NOT MET / owner-waived. AC21 not met. Every other
  AC passes on the macOS host, with AC5 (R8-MINOR-1) and AC15 (R8-NIT-1) caveats.
- **Limits:**
  - Windows runtime is not verified (cross-build only, plus the stored CI log).
  - The GitHub CI runs are producer evidence.
  - The baseline binary is reused from round 07.
  - R8-NIT-1 rests on source reading only.
  - This review grants no close, signoff or release.
