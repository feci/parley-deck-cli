---
agent: claude-1
idea: meta-protocol-change-quota-auto-exclude
review-round: 6
date: 2026-10-06
reviewed-commit: 25ea1da73ab0eab5e6099f759e2720aff54ada88
product-commit: fac40aaa1bba5350351bab0310ea8ede30ed3e6e
skill-commit: e976f7c9515f3250761380c1e09295e0f6c59985
responding-to: [review/round-05/claude-1.md, review/round-05/consensus.md]
---

## Summary

This is the full-scope Phase-6 re-review after the narrow fix-up cycle 3. It covers CLI product `fac40aa`
(review snapshot `25ea1da`; the product paths of the two commits are identical) and skill `e976f7c`, against
FINAL `27e42b8` and skill `a5664d8`.

Cycle 3 repairs most of what round 05 found. On both volumes the actual pre-change CLI and the current CLI
now give the same exit outcomes for D1–D4 (CLI-dispatched catch-up, edit-first join, kickoff-excluded return
during round 1, and the §5 decline). An archived or deleted applied notice is no longer re-published. I ran
the broad checks myself this round: the full Go suite, build, vet, race, the Windows cross-build, focused
shared/local suites, the skill suite, packets for phases 0/5/6/8, and a real `parley` supervisor crash on both
volumes. All of them pass.

The implementation is **not ready for an attended close**:

1. **R6-MAJOR-1 (new, cycle-3 G16 code).** If the owner appends a note to an applied transition notice
   (live or archived), the whole idea is integrity-gated. `parley status` reports `consensus=error`, a
   surviving member's `consensus signoff` exits 1, and the driver's `Before` blocks. The same probe on the
   cycle-2 product (`0ee1888`) leaves status and signoff working. Executed on both volumes.
2. **R6-MINOR-1 (cycle-3 G15 code).** A kickoff-excluded agent's round-1 return is recognized only while
   the prompt still carries its `excluded:` display marker. If the re-inclusion edit also removes that
   stale marker, the CLI treats the return as a new identity's pending catch-up. The returned agent's
   signoff then exits 1. The baseline accepts the same edit. Executed on both volumes.
3. **R6-MINOR-2 (W1 disclosure).** The CHANGELOG says existing policy-off membership behavior is preserved,
   and the release-notes draft lists only C1 and the bare 503 as policy-off changes. The remaining
   policy-off differences from `27e42b8` are not listed as owner-visible changes.

Counts: CRITICAL 0, MAJOR 1, MINOR 2, NIT 0. Trajectory: 15 (round 03) → 6 (round 04) → 3 (round 05) → 3
(round 06). R6-MAJOR-1 is a fresh MAJOR on cycle-3 fix code. Under the signed stopping condition and the
owner's Q1, this is **another trajectory stop for the owner, not an automatic cycle 4**. Native-positive AC2
remains **NOT MET**. The owner waived it for this release (Q2); this review does not record it as PASS. This
review signs nothing, accepts no code and grants no close.

## Protocol context and complete file coverage

**Packet (PRIMARY, executed 2026-10-06T16:13:52Z).** Command:
`parley protocol packet --dir . --phase 6 --track deliberation --idea meta-protocol-change-quota-auto-exclude --flag protocol_change --json`.
I ran it with my own build of `25ea1da` (`.parley-runtime/claude1-r6/bin/parley`, parley 1.50.0) and with the
installed parley 1.50.0. Both exit 0 and produce byte-identical JSON.

- `context_mode=full`.
- `source_sha256=packet_sha256=73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e`.
- `fallback_reason` is absent.
- The shadow packet (`6d56607c…`, 83,009 B, 32 omitted blocks) was not used. No optimized or shadow
  context was used.

The body `.parley-runtime/protocol-packets/full-phase6-deliberation-73613f95….md` is copied to
`.parley-runtime/claude1-r6/packet/body.md`. Its hash reproduces. `cmp` against the deck `COOPERATION.md`
exits 0, and the deck and skill `references/COOPERATION.md` are byte-identical (1,501 lines, 125,862 B). The
embedded `internal/protocol/defaults/COOPERATION.md` is `a92d8123…`; it differs only by the intended bootstrap
substitutions, and the unmodified `TestEmbeddedDefaultMatchesLiveDeck` passed in my full run.

**Actual coverage: the entire body**, read in four bounded chunks (1–380, 381–760, 761–1140, 1141–1501). The
over-long lines 59, 452 and 1002 were displayed in full, and their tails were also printed separately.

**Authorities read in full this round:**

- FINAL.md (902 lines);
- the scope-reset answer;
- the round05-answer;
- both standing-permission answers;
- the archived scope-reset question;
- IMPLEMENTATION.md (living, 747 lines);
- my `review/round-05/claude-1.md`;
- the signed cycle-3 plan `review/round-05/consensus.md` (G15–G17, W1–W5). Its hash is `ed4d8891…`, equal
  to `review/consensus.md`.
- the producer maps `codex-1-fixup-3-evidence.md` (subprocess, including its historical W2 stop),
  `codex-1-fixup-3-parent-evidence.md` and `codex-1-fixup-3-host-evidence.md`;
- the host `summary.json` and `host-checks.py`;
- the release-notes draft;
- the follow-up `quota-zcode-native-exhaustion-capture/00-prompt.md`.

**Diff identity (PRIMARY).** I regenerated every diff with the manifest's pathspec, and each matches the
manifest:

| Diff | Command | SHA-256 | Bytes |
|---|---|---|---|
| CLI product | `git diff 27e42b8 25ea1da` | `31440914…1ca` | 565,265 |
| skill product | `a5664d8..e976f7c` | `75c73ecd…58d` | 44,443 |
| CLI fix-up | `0ee1888..25ea1da` | `bfb1631d…b72` | 61,623 |
| skill fix-up | `e2f3649..e976f7c` | `8675238b…338` | 10,693 |

All 29 CLI product chunks, 3 skill product chunks, 4 CLI fix-up chunks and 1 skill fix-up chunk match their
manifest hashes and byte counts. Each set concatenates to its full diff
(`.parley-runtime/claude1-r6/diff/chunk-hash-check.txt`). `git diff fac40aa 25ea1da` over the product pathspec
is empty. Both worktrees are clean on product paths, at `25ea1da` and `e976f7c`.

**File coverage (117 CLI files, 5 skill files). Actual coverage this round: the complete product diffs.** I read
every chunk in the manifest in full, in bounded tool calls:

- all 29 CLI product chunks (565,265 B);
- all 3 skill product chunks (44,443 B);
- all 4 CLI fix-up chunks and the skill fix-up chunk, as aids.

Each file listed below was therefore read in its complete product hunk. The closing note records the order
and timing.

- **Cycle-3 code: 19 CLI files and 4 skill files, all read in full this round** through the complete
  incremental diffs (all fix-up chunks):
  - `CHANGELOG.md`, `docs/quota-membership.md`;
  - `internal/app/agents_exec.go`, `internal/app/quota_cycle3_test.go` (new);
  - `internal/consensus/catchup.go` (new), `catchup_test.go` (new), `consensus.go`, `quota_test.go`;
  - `internal/membership/cycle3_test.go` (new), `membership.go`, `notice.go` (new);
  - `internal/protocol/quota.go`, `quota_manual.go`;
  - `internal/quota/notice.go` (new), `receipt.go` (new), `revision.go`, `revision_snapshot.go`;
  - `internal/runner/telemetry.go`, `internal/telemetry/quota.go`;
  - skill: `CHANGELOG.md`, `SKILL.md`, `parley-addon.json`, `references/ROSTER_AND_PROTOCOL.md`.

  The eight files that are new to the product diff since round 05 (`CHANGELOG.md` and the seven new Go
  files) appear in full in these hunks.
- **The other 98 CLI files** are untouched since `0ee1888`. A mechanical per-file check finds that their
  hunks in `git diff 27e42b8 25ea1da` are byte-identical to the round-05 product diff
  (`diff/unchanged-since-r5.txt`: 98 identical, 0 differ). I re-read them in full this round in the product
  chunks, and also read them in round 05. The areas are:
  - docs `cli-reference.md`;
  - app (14 non-cycle-3 files);
  - config;
  - driver;
  - fsutil;
  - membership `crash`/`gates`/`lock`/`revision` and its older tests;
  - pidlease (Windows read-only);
  - protocol (`participantartifact`, `quota_retained`, `workspace`, embedded defaults);
  - quota (`authority`, `history`, `quota`, `record`, `revision_snapshot` base);
  - quotatest, runcontrol, runmanifest;
  - runner (all except `telemetry.go`);
  - runplan, runstate, store;
  - telemetry (all except `quota.go`), including the testdata/fixtures;
  - deck `COOPERATION.md` and `protocol-changelog.md`.
- **Product-chunk reading this round:** complete (see the closing coverage note).

**My executed checks (PRIMARY; raw logs under the git-ignored `.parley-runtime/claude1-r6/`).**

PATH setup: `PATH` was prefixed with my own guard directory (`no-provider-bin`, 21 provider names; version
queries are stubbed, any other call is denied and logged), then the producer guard. `guard-denied.log` does
not exist. No provider, participant, zcode or Claude entry point was invoked. No product, test, roster,
credential, worktree, review or consensus file was edited. Product and skill source hashes were unchanged
across my checks. Both `agents.toml` hashes equal the producer's recorded values (`f8cc2ab5…`, `f52a0a77…`).

- **G (broad Go, `checks.py go`, host macOS arm64, go1.27.1).** All of these exit 0:
  - `go test ./... -count=1 -timeout 45m` (528.0 s; 34 packages `ok`, no FAIL);
  - `go build ./...`;
  - `go vet ./...`;
  - `GOOS=windows GOARCH=amd64 go build ./...`, and `go vet` of pidlease/membership/fsutil for Windows;
  - `go test -race` over pidlease/membership/app/consensus/telemetry/runner with the producer's `-run`
    pattern (262 PASS lines, 0 FAIL);
  - the focused cycle-3/cycle-2/fixup/lease/sync suite once with `TMPDIR` on this shared volume and once on
    `/tmp` (118 PASS, 0 FAIL each; 52 `TestQuotaCycle3*` pass lines; `TestQuotaCycle2NativeCrashWriter` passes).

  gofmt over all 93 changed Go files finds nothing unformatted. The skill suite and my probes ran
  concurrently with the full suite, which still passed.
- **S (skill).** `npm test` in the skill worktree exits 0 in 82.9 s: 399/399 Node, 54 Python OK, all six addon
  manifests ok.
- **K (packets).** Phases 0, 5 and 8 render `full` with source = packet = `73613f95…`, and no fallback.
- **D (actual `app.Run` differential).** The unchanged extended D1–D4 main (`7a87344b…`, identical on both
  sides; only `setup.go` differs) ran against the `27e42b8` archive and against the current tree. Each side
  ran on `/tmp` and on this shared volume. All four runs exit 0 (`d14/*.log`).
- **V (my variants, same helpers, identical main `b85dabcd…` on both sides, both volumes).** The cases are:
  - V1: kickoff-excluded return with the marker removed;
  - V2: edit-first, all existing members sign while the joiner is pending, then the joiner's own CLI
    catch-up and signoff;
  - V3: own incomplete stub that exits 1, then a retry in place;
  - V4: decline, then catch-up anyway;
  - V5: decline, then the decliner is removed from `participants:`.

  Logs: `vprobe/{baseline,current}-{local,shared}.log`.
- **N (notices, current, both volumes).** `notice6` covers owner archive and owner delete after the applied
  receipt, then two `Before` calls and `status`. `notice7` covers an owner-appended note on the live or the
  archived copy, then `Before`, `status` and a survivor's `consensus signoff`. The same `notice7` (minus the
  cycle-3-only `ReadApplied` print) also ran against a `git archive 0ee1888` cycle-2 tree on `/tmp`.
- **C (real `parley` crash, both volumes).** I inspected `real-crash/run.py` and `init.go`. My copy changes
  only the guard path and the binary (my build), and writes to new evidence directories:
  - It launches a real `parley run` supervisor (`start_new_session`) with two local stub agents that exec
    a private Python sleeper.
  - It binds both `started.json` records to that supervisor PID and verifies pid == pgid and the exact
    worker command.
  - It SIGKILLs only that supervisor and those two verified groups.
  - It then runs `parley quota recover` twice.

  Result on both volumes: two immutable `crash-settlement.json` files, unchanged on replay, and no
  `terminal.json`. The settlements record the same host and boot, the supervisor PID, the writer PID
  (= process group), and the proof `same-host-and-boot; supervisor and writer PIDs absent; supervised process
  group absent`. **My variant (`run_live.py`):** with the supervisor dead and both writers still alive,
  `quota recover` exits 1 with `unsettled writer … only proven-dead local supervisor, writer and process group
  may be settled`. No settlement or terminal is written. After the writers are killed, recovery settles as
  above. Writer PIDs are confirmed dead afterward.

**Not executed by me:** a Windows runtime (no Windows pass is inferred from cross-compilation); a Linux
container or PID-namespace run; a rerun of the G14 `source-audit.py` extraction; any native zcode capture
(none authorized). The producer's 567 s full run, race, Windows build and real-crash results remain
SECONDARY. My own runs of the same commands are PRIMARY.

## User direction

**New owner answer since round 05: round05-answer, quoted verbatim.** Source:
`inbox/user-to-codex-1_meta-protocol-change-quota-auto-exclude_round05-answer.md`, re-read in full. The
originals are in Slovak, and the relay's translations follow them.

> ## Owner answers to `codex-1-to-user_meta-protocol-change-quota-auto-exclude_round05-trajectory-ac2.md`
>
> Relayed by the owner's Claude Code session on 2026-10-06 at about 15:50 CEST. The questions and answers
> below are verbatim (Slovak), each followed by a translation.
>
> **Q1.** "Po druhom kole opráv zostali 3 nálezy (predtým 15, potom 6). Prvý je nový MAJOR: pri vypnutom
> nastavení oprava pokazila neskoré pripojenie agenta k idei a návrat agenta, ktorý bol vyradený pri štarte.
> Povoliť tretí, úzko ohraničený cyklus opráv?"
>
> Selected: **"Áno, úzky cyklus 3 (Recommended)"**. The option read: "Vrátiť pôvodné správanie pri vypnutom
> nastavení a opraviť drobnosť s archivovanou notifikáciou. Plán podpíše claude-1 a potom urobí review."
>
> Translation of the answer: "Yes, narrow cycle 3. Restore the original behavior with the setting off and fix
> the small issue with the archived notice. claude-1 signs the plan and then reviews."
>
> **Q2.** "Druhý MAJOR: rozpoznávač chyby zcode môže byť v praxi neúčinný. Reálny stderr zcode vnorené objekty
> skracuje na „[Object]“ a rozpoznávač taký záznam odmietne. Codex navrhuje zachytiť skutočný výstup zcode.
> Overil som však, že zcode má teraz kvótu: o 14:57 vrátil PONG. Takto by sa zachytila len úspešná odpoveď,
> nie chyba kvóty. Čo s tým?"
>
> Selected: **"Vydať s obmedzením + follow-up (Recommended)"**. The option read: "Vydať so známym obmedzením:
> rozpoznávanie zcode môže byť neúčinné a namiesto automatiky sa agent opýta teba, čo je bezpečné. Výstup
> najbližšieho skutočného vyčerpania zcode sa zachytí a gramatika sa podľa neho upraví v nadväzujúcej idei."
>
> Translation of the answer: "Release with the limitation plus a follow-up. Release with a known limitation:
> zcode recognition may be inert, so instead of the automatic path you are asked, which is safe. The output of
> the next real zcode exhaustion is captured, and the grammar is adjusted to it in a follow-up idea."
>
> **Relay fact (PRIMARY).** At 14:57 CEST the relay ran one zcode invocation with the prompt "Reply exactly
> PONG. Do not call tools, read or write files, or execute commands." It exited 0 and printed `PONG`, with
> empty stderr. zcode is not exhausted now, so a capture today could not produce the qualifying native
> error.
>
> ## What this authorizes
>
> - **Q1.** One narrow fix-up cycle 3 within the bounds in your note (R5-MAJOR-1 and R5-MINOR-1), with the
>   normal Phase-7 plan signed by claude-1 and a separate full re-review. A fresh CRITICAL or MAJOR on the
>   cycle-3 fix code escalates again, as you proposed.
> - **Q2.** An explicit owner waiver of AC2's native positive evidence for this release, with no capture and no
>   grammar relaxation now. Record it in `IMPLEMENTATION.md` as an owner-accepted known limitation, and state
>   it in the protocol or skill wording and in the release notes. The wording should say that zcode
>   auto-exclusion may not fire on real native output, and that unrecognized failures fall back to the
>   owner-confirmed path. Open a linked follow-up idea whose precondition is a captured native zcode
>   exhaustion. If cheap and in scope, the CLI should retain the raw scrubbed stderr of a failed zcode
>   invocation privately, so the next real exhaustion becomes that capture. Otherwise, name this in the
>   follow-up. R5-MAJOR-2 is dispositioned as owner-accepted and deferred, not fixed.
> - Everything else is unchanged: the claude-1 signoffs, a NEW attended-close request with current-tree
>   evidence, then the release. Both standing retry permissions still apply to failed invocations.

This quotes the whole answer body after its frontmatter.

**Previous answers.** They are re-read and unchanged; my exact earlier quotations stand:

- scope-reset: `review/round-03/claude-1.md:62-109`;
- review-quota: `review/round-04/claude-1.md:147-168`, with its erratum at `:516-531`;
- quota-standing and timeout-standing: `review/round-04/consensus.md:158-235`.

The bounded zcode stderr rule and the display-clock rule are owner policy; I check the code against them and do
not relitigate them. The standing permissions govern only organizer relaunches. This process made no retry or
provider probe. The Q2 waiver changes only AC2's release disposition. I treat no other gate as waived.

## Refutation attempts

PRIMARY = my execution (section above) or my direct reading of `25ea1da` at the cited locator. SECONDARY =
producer evidence that I did not rerun.

### G15–G17 (signed cycle-3 plan)

| Item | Attempt | Result |
|---|---|---|
| G15 catch-up / kickoff return | D1–D3 baseline vs current on both volumes; V1–V3; reading `protocol/quota_manual.go` (`manualQuotaRevision`, `ManualCatchupTarget`, `manualCatchupPathSafe`), `runner/telemetry.go` `beginLaunch`, `app/agents_exec.go` import, `membership.go` `Before`, `quota/revision_snapshot.go` `ManualRoundOneReturn`; the `TestQuotaCycle3CLI*` boundaries (policy-on, foreign path, later round, known-excluded, kickoff-excluded without an edit, closed, wrong idea, symlink, foreign stub, complete artifact) passed in my focused runs | **Mostly resolved.** D1: exit 0 and the artifact is written, as on the baseline. D2: existing signers append (V2: all four), status shows `quota transition pending: pending catch-up for e; … parley agents exec --agent e --artifact "…/round-01/e.md" …` and no integrity gate; the joiner's own CLI catch-up then imports it (revision 1, known), and its signoff makes consensus `ready`, as on the baseline. D3: both signoffs exit 0, revision 1, manual, and the kickoff is unchanged. V3: a retry over its own incomplete stub now works (the baseline refused it). Residual: R6-MINOR-1 (return with the marker removed). |
| G15/W2 decline | D4 and V4/V5 on both trees and volumes; reading `consensus/catchup.go`, `AppendSignoff`, `validateDocumentWithDecliners`; `TestQuotaCycle3PendingDecline*` and `TestQuotaExcludedVeto…` (known excluded id cannot use the decline) passed | **Resolved, with the disclosed difference.** The exact BLOCK form exits 0 (`Consensus: blocked`) on both trees; the literal status exits 1 on both; an existing signer's later append exits 0 on both. Current keeps `e` in `Missing signoffs` (`a,b,c,d,e`; the baseline shows `a,b,c,d`), as G15 requires. After the decliner catches up (V4), both trees read `blocked`, `Missing a,b,c,d`. Removing the decliner (V5) makes consensus `malformed` on **both** trees, which predates this idea (Open question 3). Exactness: only the trimmed note `❌ NON-PARTICIPANT` with a counter-proposal, design consensus only, policy-off, pending ids only; it adds no history (revision 0). |
| G16 notice | N on both volumes; reading `membership/notice.go`, `quota/notice.go`, `quota/receipt.go`, the cycle-3 `InspectQuota` hunk; `TestQuotaCycle3NoticeDeliveryRecovery` (applied/before-receipt archive and delete, before-notice, corrupt-receipt), `…CorruptNoticeNeverCompletes`, `…HistoricalManualClarificationSurvivesInboxActions` passed | **Archive and delete resolved; adjacent new defect R6-MAJOR-1.** After the receipt, archive gives 0 in the inbox and 1 archived; delete gives 0 and 0. Two later `Before` calls re-publish nothing, `status` exits 0 with no integrity/pending text, the receipt is valid, and `round.completed` = 1, on both volumes. An owner-appended note gates the idea. |
| G17 disclosure | Reading the support table (`telemetry/quota.go`), `docs/quota-membership.md:211-226`, `CHANGELOG.md`, the skill `SKILL.md:173-177` and `ROSTER_AND_PROTOCOL.md:329-333`, the release-notes draft, IMPLEMENTATION AC2 rows (`:171`, `:557`), and the follow-up prompt | **Met** (see W4/W5). |

### AC1–AC21 (my own scoped evidence)

- **AC1 PASS (scoped).**
  - The phase 0/5/6/8 packets are full, with packet = source = deck hash (K).
  - The deck and skill copies are byte-identical (`cmp`).
  - `TestEmbeddedDefaultMatchesLiveDeck` passed in G.
  - No protocol hunk changed in cycle 3 (the fix-up diff touches no COOPERATION or changelog file).
  - The `meta/protocol-changelog.md` entry is present (product chunk 029, re-read this round; unchanged in
    cycle 3).
- **AC2 NOT MET; owner-waived for this release (Q2).**
  - R5-MAJOR-2 is owner-accepted and deferred, not fixed and not PASS.
  - Source-derived fixtures keep their labels: no fixture file changed in cycle 3, and the support-table
    limitation now says "native-positive evidence NOT MET and owner-waived".
- **AC3 PASS.** The full and focused runs include every recognizer negative. No telemetry grammar changed in
  cycle 3; only the support-table limitation string changed.
- **AC4 PASS.** By reading, `QuotaSupport` lists zcode as the only supported adapter (owner deviation) and the
  rest as diagnostic-only.
- **AC5 PASS** (G: whole-batch, permutation, floor, duplicate and facilitator tests).
- **AC6 PASS.**
  - Kickoff filtering and C1 pass (G).
  - D1 shows the joiner's own dispatch is restored.
  - The exception never launches a known excluded member (`known-excluded` boundary test) or an id outside
    its own `round-01/<id>.md`.
- **AC7 PASS** (G).
- **AC8 PASS** (G).
- **AC9 PASS.**
  - Retained vetoes hold (G).
  - A pending joiner gains no known status (V2 history before import: revision 0).
  - An excluded known id cannot use the decline route.
  - A kickoff-excluded id is known only from revision 1 (D3 `verify`).
- **AC10 PASS** (G).
- **AC11 PARTIAL.**
  - Fault, pending, replay and crash pass (G, C).
  - Archive and delete pass (N), with one evaluation.
  - An owner-annotated notice gates the idea (R6-MAJOR-1).
- **AC12 PASS** (G on both volumes, race).
- **AC13 PASS.**
  - G.
  - V3: the incomplete own stub is preserved and is not imported (revision 0) until a valid retry.
- **AC14 PASS** (G; D3 and V2 history readouts).
- **AC15 PARTIAL** (R6-MAJOR-1).
- **AC16 PARTIAL.** There are remaining knob-off differences (R6-MINOR-1), and they are not disclosed in
  release surfaces (R6-MINOR-2).
- **AC17 PASS (scoped).**
  - Roster hashes are unchanged.
  - There is no timer or automatic rejoin.
  - A same-idea return works by the plain edit, with the R6-MINOR-1 caveat.
- **AC18 PASS** (G: quoted, tool and content negatives).
- **AC19 CONDITIONAL, honored.** `status: fix-up-cycle-3`; the cycle-3 subsection says it is not overall
  completion.
- **AC20 PASS on this macOS host (PRIMARY: G, S, K, gofmt).** The Windows runtime is unexecuted.
- **AC21 NOT MET.** No final signoffs and no NEW attended close.

## Findings

### [MAJOR] R6-MAJOR-1: An owner-annotated applied notice integrity-gates the whole idea

**Locators (25ea1da).**

- `internal/protocol/quota.go`: the cycle-3 hunk in `InspectQuota` calls `quota.InspectNotice` for every batch
  notice, and for the `-manual-authority` clarification, before any read result is returned.
- `internal/quota/notice.go` `InspectNotice` requires the exact bytes of `Notice()` or `LegacyManualNotice()`
  for every extant copy in `inbox/` **and** `inbox/archived/`, and otherwise returns
  `corrupt or mismatched quota notice`.
- `internal/membership/notice.go` `checkedNotice`: "Validate every existing copy, even with an applied
  receipt."

**Executed (PRIMARY, on both volumes: `notice7/local-*.log` on `/tmp` and `notice7/shared-*.log` on this shared volume).** After a settled policy-on
batch with a valid applied receipt, the owner appends `## Owner answer / Acknowledged; no re-inclusion this
idea.` to the notice. Two cases were run: the notice kept in `inbox/`, and the notice moved to
`inbox/archived/` and then annotated. Results on current `25ea1da`:

- `Before` → `quota batch blocked: integrity/recovery gate: corrupt or mismatched quota notice: …`.
- `parley status` → `consensus=error`, `quota transition pending/integrity gate: corrupt or mismatched quota notice`.
- A survivor's `consensus signoff --agent a` → exit 1, with the same error.

On the cycle-2 product `0ee1888`, the same program gives these results:

- Annotated live notice: `Before` is blocked (`immutable quota record differs`), but `status` is clean and the
  survivor's signoff exits 0.
- Annotated archived notice: no gate, and the signoff exits 0. The notice is re-published, which was
  R5-MINOR-1.

Cycle 3 therefore introduced the idea-wide gate on every read and write path. On this policy-on idea, each
gated `Before` also files a `blocking: yes` integrity escalation for the owner
(`inbox/parley-to-user_quota-block-*.md`, `phase: quota-recovery`) that calls the notice corrupt. I observed
it in both runs.

**Why it matters.**

- §4 lets the owner answer a `to: user` inbox file "by any medium: … an appended answer block in the
  escalation file itself", and it says the inbox file "is not the authoritative record".
- After a valid applied receipt, publication is already proven. The bytes of a non-authoritative,
  owner-owned file then decide whether any participant can sign, dispatch or wait.
- The diagnostic calls an owner note "corrupt". Recovery means deleting the file or restoring its exact
  bytes; deletion is accepted, but nothing tells the owner that.
- W3 required that normal owner inbox handling never become the gating "ambiguous state". The Q1 bound was
  "fix the small issue with the archived notice". This fix widens the gating surface instead.
- The exposure is not limited to hand-written notes. Any tool or editor that rewrites the file's bytes
  (line endings, a final newline, a markdown formatter or pre-commit hook) has the same effect. That broader
  exposure is my RECALL-level reasoning, not executed.
- The documentation does not cover this. The skill (`ROSTER_AND_PROTOCOL.md`: "Contradictory extant notices
  still gate") and the CLI docs ("still require investigation") disclose that a gate exists. Neither says
  that an owner note, or any byte change, counts as "contradictory", and neither says to answer in a
  separate file.

**Suggested fix.**

- Once `ReadApplied` is true, never gate on extant notice bytes. A regular file at either location is
  delivered, and its content is the owner's.
- Before the receipt, treat a mismatched copy as "not delivery proof" and publish once. Do not raise an
  idea-wide integrity gate.
- Keep the symlink and non-regular-file refusals if wanted, but scope them to the publication step, not to
  `status`, `consensus` or dispatch.
- Alternatively, the owner may accept this as a documented limitation ("never edit transition notices;
  answer in a separate file").

### [MINOR] R6-MINOR-1: A kickoff-excluded agent's round-1 return depends on keeping its `excluded:` display marker

**Locators.**

- `internal/quota/revision_snapshot.go` `ManualRoundOneReturn`: the return is recognized only when
  `status: round-01` **and** `ConfirmedInPrompt(raw, "excluded", id)`.
- `internal/protocol/quota_manual.go` `manualQuotaRevision`: otherwise, an id that is not known needs a valid
  late `round-01/<id>.md`, or it becomes `manualCatchupPending`.
- `internal/membership/membership.go` `Before`: `len(v.Catchup) > 0` returns the pending error and blocks
  ordinary driving.

**Executed (PRIMARY, V1 on both volumes).** At kickoff, `d — unavailable — confirmed 2026-10-04` is excluded.
During round 1, the edit adds `d` to `participants:` and deletes the now-stale `excluded:` line.

- Baseline `27e42b8`: signoffs by `a` and `d` both exit 0.
- Current:
  - `a` exits 0.
  - `d` exits 1, with `pending catch-up for d; complete each own late round-1 artifact … parley agents exec --agent d …`.
  - History stays at revision 0, and `d` is not known.

With the marker kept (D3), both exit 0.

**Why it matters.**

- The docs (`docs/quota-membership.md:77-78`) and the skill (`ROSTER_AND_PROTOCOL.md:298`) promise that "a
  kickoff-excluded id may return by a plain edit during round 1". The skill also says never to derive quorum
  from `excluded:` markers.
- Removing a stale exclusion line is a natural part of re-inclusion, yet the code infers the kind of import
  from that mutable display marker. The immutable kickoff record already holds the exclusion.
- The outcome fails closed and is actionable, and during round 1 the agent owes its round-1 file anyway. So
  this is a narrow knob-off residual of the owner's "restore the original behavior", not a hard block.
- Driven runs halt in `Before` until someone runs the manual `agents exec` (by reading; covered by
  `TestQuotaCycle3CLICatchupBothOrdersAndIncompleteRetry`'s `Before` assertion).

**Fix.** Decide "kickoff-excluded" from immutable kickoff history, not from the current prompt marker.
Otherwise, document that the marker must stay, and list this difference (see R6-MINOR-2).

### [MINOR] R6-MINOR-2: Release surfaces overclaim policy-off preservation (W1 listing)

**Locators.**

- `CHANGELOG.md` (Unreleased): "Preserve historical vetoes/findings and existing policy-off membership behavior".
- `source-context/codex-1-release-notes-draft.md`: "Two fixes apply even when the policy is off: …" (only C1
  and the bare 503).

**The differences from `27e42b8` that remain with the policy off.** These include changes G15 deliberately
introduced:

1. A joiner cannot sign or complete the quorum before its late round-1 imports. The baseline signs right after
   the edit. Mechanism executed in V1; asserted by the passing `TestQuotaCycle3*` tests.
2. The exact decline leaves the joiner in `Missing signoffs` (D4).
3. A kickoff-excluded agent returning after round 1 needs catch-up (the `round-02` case of
   `TestQuotaCycle3CLIKickoffExcludedReturn`).
4. A pending catch-up stops ordinary driving in `Before` until the manual exec.
5. Membership edits on a `final` or `closed` idea are rejected (`closed idea membership is frozen`; by
   reading only).
6. R6-MINOR-1.

**Why it matters.** W1 said: "List any remaining difference as an owner-visible knob-off change; FINAL §11
permits only C1 and the bare 503." The CLI docs and skill describe items 1–3, but the owner-facing release
surfaces claim preservation.

**Fix.** List items 1–6, or whichever remain after fixes, in the release notes and CHANGELOG as policy-off
behavior changes for the owner's attended close.

## Prior finding and W1–W5 dispositions

All of these are my own evaluation.

- **R5-MAJOR-1: resolved for D1–D3 as specified** (executed on both trees and both volumes). The residuals are
  R6-MINOR-1 and the R6-MINOR-2 listing.
- **R5-MAJOR-2: owner-accepted and deferred (Q2), not fixed.** AC2 is NOT MET and waived for this release
  only. I do not withdraw the finding. The disclosure is consistent (W4).
- **R5-MINOR-1: resolved for archive and delete** (N on both volumes, with one evaluation). The adjacent
  defect is R6-MAJOR-1.
- **W1: partially met.**
  - Met: D1, D2 and D3 outcomes equal the baseline. Existing-member appends work while catch-up is pending
    (V2). `status` shows pending with the exact command, not an integrity gate. CLI completion works, and so
    does the own incomplete-stub retry (V3). Excluded known ids cannot dispatch (boundary test). There is no
    false quorum: with all four existing members accepting and `e` pending, the result is `partial`,
    `Missing e`, on both trees.
  - Not met: the listing (R6-MINOR-2) and the marker variant (R6-MINOR-1).
  - In my run, the D2 status line also reports the probe-setup-only missing `events.jsonl`, as the producer
    disclosed. That comes from the shared probe's setup, not the product.
- **W2: met.**
  - The baseline outcome is preserved for exit codes and triage.
  - The Missing-list difference is the disclosed G15 consequence.
  - The parent's shared-consensus API addition is narrow: exact trimmed note, counter-proposal required,
    BLOCK only, design only, policy-off pending ids only.
  - It is consistent on the append and read paths.
  - A hand-filed duplicate or out-of-scope decline reads as `malformed`, as the passing tests show.
- **W3: met for archive and delete.**
  - Neither re-publishes or gates.
  - Interrupted publication is stated (`docs/quota-membership.md`: one benign re-publication after deletion
    before the receipt) and tested.
  - Owner annotation: R6-MAJOR-1.
- **W4: met.**
  - (a) The skill, CLI docs, CHANGELOG and release notes say capture is not automatic. They tell the owner to
    keep the run's private, unscrubbed `parley-deck/runs/<run-id>/agents/<agent-id>/stderr.log` in place
    (path verified at `internal/runner/runner.go:404-406`), never to copy or commit it, and they say its
    completeness and cleanup are unverified.
  - (b) IMPLEMENTATION's AC2 rows, `docs/quota-membership.md`, the support table and the skill reference say
    NOT MET and owner-waived, never PASS.
  - No fixture is relabeled native. Nothing invoked a provider or zcode, captured output or changed the
    grammar (no telemetry grammar or fixture bytes changed in cycle 3).
  - The core `SKILL.md` says "owner-waived" and points to the reference that says NOT MET.
- **W5: met.** The follow-up is `status: candidate` with no `participants:` line, and the intended
  organizer/quorum appears in prose only (`00-prompt.md:54`).
- **Retained earlier gates** (G10–G14, V1–V9, R1–R4, round-03 dispositions) are unchanged from round 05. Their
  code is untouched by cycle 3, apart from the `beginLaunch` known-and-current condition.
  - The broad G/C runs re-exercise them.
  - V8/R3: re-verified with a real process (C), including refusal while writers are alive.

## Open questions

1. **Trajectory stop (signed condition, owner Q1).** R6-MAJOR-1 is a fresh MAJOR on cycle-3 fix code. The
   trajectory is 15 → 6 → 3 → 3, with no CRITICAL since round 03. The owner's options are:
   - (a) authorize a further narrow fix for R6-MAJOR-1, optionally with R6-MINOR-1 and R6-MINOR-2;
   - (b) accept R6-MAJOR-1 as a documented limitation ("do not edit transition notices"), plus listing;
   - (c) another disposition.

   I recommend (a): the fix is small and local to `notice.go` and `InspectQuota`. There is no cycle 4 by
   default.
2. **Windows (limitation; RECALL reasoning, UNVERIFIED).** Several checks compare exact bytes: notices,
   receipts and hash-bound history. A Windows checkout with `core.autocrlf=true` could therefore gate ideas.
   Only cross-compilation was executed, and Windows remains experimental.
3. **Pre-existing, not this idea.** Resolving a §5 decline by removing the decliner from `participants:`
   makes consensus `malformed` on the baseline and on the current tree alike (V5). Should a follow-up define
   the decline's resolution?
4. **Container/PID namespaces (UNVERIFIED reasoning, carried from round 05).** Same-host/boot proof uses the
   hostname and boot time. A writer in another PID namespace that shares the hostname could look dead to
   `kill(pid, 0)`. My crash runs cover one macOS host only.

### Closing coverage note

The draft of this artifact was written at about 16:26Z, before the CLI product chunks had been re-read. I then
read every remaining chunk in order and updated this note after each batch. Read in full in round 06:

- all 4 CLI fix-up chunks (61,623 B);
- the skill fix-up chunk (10,693 B);
- all 3 skill product chunks (44,443 B);
- all 29 CLI product chunks, 001–029 (565,265 B), finished at about 16:29:40Z.

That reading changed no verdict. It added two facts:

- On policy-on ideas, the R6-MAJOR-1 gate also files a blocking owner escalation (`membership.Before`
  `IntegrityBlock`). This is now recorded in R6-MAJOR-1.
- `consensus.Draft`, `AppendSignoff` and `Reopen` call `membership.RecordManual` first. A pending catch-up
  returns no error there: `InspectQuota` returns before it sets `Manual`. This matches V2.

No chunk was skipped or truncated. Six lines are longer than 2,000 characters:

- CLI 014:329, 015:135, 028:24 and 029:32;
- skill 001:157 and 002:121.

All six are unchanged context lines: the §4 Phase 5 paragraph and §9 checklist item 1. They were displayed in
full, and their tails were checked separately; they end "…re-pin." and "…full COOPERATION.md.". They are
byte-identical to packet-body lines 452 and 1002.
