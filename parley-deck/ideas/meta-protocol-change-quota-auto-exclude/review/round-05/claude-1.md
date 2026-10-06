---
agent: claude-1
idea: meta-protocol-change-quota-auto-exclude
review-round: 5
date: 2026-10-06
reviewed-commit: e7bf96c2f8914fea151ed1b862ba658b6f3852b9
product-commit: 0ee18889977a29d6baf53d3016b501112760ae12
skill-commit: e2f3649eb938e870367c76943b441382fe7acf65
responding-to: [review/round-04/claude-1.md, review/round-04/consensus.md]
---

## Summary

Full-scope Phase 6 re-review after fix-up cycle 2: CLI product `0ee1888` (review snapshot `e7bf96c`; the
product paths of the two commits are byte-identical) and skill `e2f3649`, against FINAL `27e42b8` / skill
`a5664d8`.

Cycle 2 fixes most of what round 04 found. Realistic changing-countdown retries now qualify and the 60-minute
threshold is anchored on terminal receipt (R4-MAJOR-1 resolved). The exclusion grammar, the policy-on
revision boundary, crash settlement and precommit manifest validation are in place.

The implementation is **not ready for an attended close**:

1. **R5-MAJOR-1 (new, G11 membership path).** On policy-off ideas the ordinary §5 catch-up join and the
   return of a kickoff-excluded agent no longer work through the CLI. Executed against the actual pre-change
   CLI: all three paths succeed on `27e42b8` and fail on the current tree, on both volumes.
2. **R5-MAJOR-2 (carried R4-MAJOR-3, owner decision).** Native AC2 evidence is still absent. By the
   producer's own source evidence, the realistic native shape (`[Object]` at default inspect depth) is
   rejected, so the only supported recognizer would likely never fire on real zcode output.
3. **R5-MINOR-1.** An applied transition's owner notice is re-created in `inbox/` after the owner archives
   it.

Counts: CRITICAL 0, MAJOR 2, MINOR 1, NIT 0. Trajectory: 15 (round 03) → 6 (round 04) → 3 (round 05).
R5-MAJOR-1 is a fresh MAJOR on G11's membership path. Under my signed stopping condition this is
**escalation with the trajectory, not an automatic cycle 3** (see Open questions). This review signs nothing,
accepts no code and grants no close.

## Protocol context and coverage

- **Packet (PRIMARY, executed 2026-10-06T12:32Z).** `parley protocol packet --dir . --phase 6 --track
  deliberation --idea meta-protocol-change-quota-auto-exclude --flag protocol_change --json` (parley 1.50.0),
  exit 0. `context_mode=full`,
  `source_sha256=packet_sha256=73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e`,
  `fallback_reason` absent. Shadow packet (`6d56607c…`, 83,009 B) not used. Body
  `.parley-runtime/protocol-packets/full-phase6-deliberation-73613f95….md`, copied to
  `.parley-runtime/claude1-r5/packet/body.md`; `cmp` against the deck `COOPERATION.md` exits 0 (1,501 lines,
  125,862 B).
- **Actual coverage: the entire body**, read in four bounded chunks (1–380, 381–760, 761–1140, 1141–1501).
  The three over-long lines (59, 452, 1002) were displayed in full and their tails printed separately.
- **Authorities read in full this round:** FINAL.md (902 lines); the scope-reset answer; both standing-permission
  answers; IMPLEMENTATION.md (510 lines, living); my `review/round-04/claude-1.md` including its erratum; the
  signed cycle-2 `review/round-04/consensus.md` with V1–V9; both producer maps
  (`codex-1-fixup-2-evidence.md`, `codex-1-fixup-2-host-evidence.md`) and the host `summary.json`.
- **Diff identity (PRIMARY).** I regenerated `git diff 27e42b8 e7bf96c` with the manifest's pathspec:
  SHA-256 `e7780acf…6257`, 517,979 B, identical to the manifest; skill diff `6f7d4ca3…7c15` identical.
  All 27 CLI and 3 skill product chunks (and the 6 fix-up chunks) match their manifest hashes and concatenate
  to the full diffs (`.parley-runtime/claude1-r5/01-chunk-hash-check.txt`). `git diff 0ee1888 e7bf96c` over the
  product pathspec is empty.
- **Complete product diffs read: all 27 CLI chunks (517,979 B) and all 3 skill chunks (41,259 B).** The cycle-2
  incremental diffs were used only as aids.

**Coverage table (every changed product file; 109 CLI files plus 4 skill files, all read in full in the
product diff).** "Exec" names my own executed checks (below). Numstats are in `02-*-numstat.txt`.

| Area (chunks) | Files | Exec |
|---|---|---|
| docs (1) | `docs/cli-reference.md`, `docs/quota-membership.md` | read vs. D1–D3 |
| app (1–6) | `agents_exec.go`, `app.go`, `consensus_request_signoffs.go`, `driver_consensus.go`, `driver_impl.go`, `organizer.go`, `pipeline_cmd.go`, `preflight.go`, `preflight_liveness.go`, `quota.go`, `quota_consumers.go`, `quota_pipeline.go`, `quota_revision.go`, `quota_signoff.go`, `wait.go`; tests `quota_fixup_test.go`, `quota_signoff_test.go`, `quota_surfaces_test.go`, `quota_test.go` | T, D1–D3 |
| config (6) | `runtime.go`, `runtime_test.go` | T |
| consensus (6–7) | `consensus.go`, `quota_test.go` | T, D2–D3 |
| driver (7–8) | `driver.go`, `impl.go`, `loop.go`, `phasedigest.go`, `quota.go`, `quota_test.go` | T |
| fsutil (8) | `sync_darwin.go`, `sync_darwin_test.go` | T |
| membership (8–11) | `crash.go`, `cycle2_test.go`, `fixup_test.go`, `gates.go`, `lock.go`, `membership.go`, `membership_test.go`, `revision.go`, `revision_test.go` | T, N |
| pidlease (11–12) | `lease.go`, `lease_test.go`, `live_unix.go`, `live_windows.go` (Windows: read only) | T |
| protocol (12–15) | `defaults/COOPERATION.md`, `participantartifact.go`, `quota.go`, `quota_manual.go`, `quota_retained.go`, `workspace.go` | T, D2–D3 |
| quota (15–18) | `authority.go`, `authority_test.go`, `history.go`, `history_test.go`, `quota.go`, `quota_test.go`, `record.go`, `revision.go`, `revision_snapshot.go` | T |
| quotatest, runcontrol, runmanifest (18) | `quotatest/authority.go`, `runcontrol/quota_test.go`, `runcontrol.go`, `runmanifest/manifest.go` | T |
| runner (18–20) | `acp.go`, `failclass.go`, `handoff.go`, `phase58.go`, `quota.go`, `quota_target.go`, `quota_test.go`, `runner.go`, `telemetry.go`, `telemetry_test.go`, `validation.go` | T, D1 |
| runplan, runstate, store (20) | `runplan.go`, `runstate.go`, `events.go` | read; T (runstate had no matching tests) |
| telemetry (20–25) | `provider.go`, `quota.go`, `quota_cycle2_test.go`, `quota_framing.go`, `quota_framing_test.go`, `quota_test.go`, `quota_zcode.go`, `quota_zcode_test.go`, `record.go`, `usage.go`; `testdata/quota/` README, harness and all 16 fixtures | T, Z |
| deck (25–27) | `parley-deck/COOPERATION.md`, `parley-deck/meta/protocol-changelog.md` | packet `cmp`; drift test in T |
| skill (1–3) | `SKILL.md`, `parley-addon.json`, `references/COOPERATION.md`, `references/ROSTER_AND_PROTOCOL.md` | read; addon hash of `references/COOPERATION.md` equals the packet hash |

**My executed checks (PRIMARY; raw logs under the git-ignored `.parley-runtime/claude1-r5/`).** No provider,
participant or zcode entry point was invoked; `PATH` was prefixed with `fixup-2/host/no-provider-bin`. No
product, test, roster, credential, worktree or other-agent file was edited.

- **T (focused suites).** `tests/run-focused.sh`: `go test ./internal/{telemetry,membership,app,pidlease,fsutil,
  consensus,quota,protocol,driver,runner,config,runcontrol,runstate} -run
  'Quota|Cycle2|Lease|SyncFile|EmbeddedDefaultMatchesLiveDeck' -count=1`, once with `TMPDIR` on this deck's
  shared volume and once on `/tmp`. Both exit 0; 12 packages `ok`, `runstate` has no matching tests
  (`tests/shared-focused.log`, `tests/local-focused.log`). The pattern includes the native
  `TestQuotaCycle2NativeCrashWriter`, `TestLeaseProcess` and all `TestQuotaCycle2*` regressions.
- **D1–D3 (differential, actual `app.Run`).** `diff-common/main.go` is byte-identical on both sides; only
  `setup.go` differs (baseline `CreateIdeaFull`; current policy-off `CreateIdeaWithQuota` plus its manifest).
  Baseline is a fresh `git archive 27e42b8` at `/tmp/claude1-r5-prechange-27e42b8` (no worktree operation);
  its `app.go`, `consensus.go` and `COOPERATION.md` hashes equal the producer's `baseline-identity.txt`
  (`03-prechange-archive-identity.txt`). Run on the shared volume and on `/tmp`:
  `diff-{baseline,current}.{shared,local}.log`.
- **N (notice republication).** `notice/main.go` on both volumes (`notice/{shared,local}.log`).
- **P (my round-04 P3/P4/P5 program, unchanged).** `go run ./.parley-runtime/claude1-r4/r4extra <fresh root>` on both
  volumes (`r4extra-rerun/{shared,local}.log`; identical except paths). Its round-04 outputs were not touched.
- **H (G14 harness).** `env -i LC_ALL=C LANG=C TZ=UTC /opt/homebrew/bin/node internal/telemetry/testdata/quota/sdk-framing-harness.cjs <mode>`
  for `default-console`, `retained-retries`, `retained-top-level` and `backoff-retries`: all exit 0, and each stdout is
  byte-identical to its committed `zcode-<mode>-source-derived.stderr`. The harness evaluates only extracted SDK
  class text in an empty `vm` context and reports source SHA-256 `3e3433d9…685f`; no zcode entry point, configuration,
  network or credential. `node -p` under `env -i` reports default inspect depth `2` (`g14harness/results.txt`).
- **S (stack forms and single-record V1).** `g14stack/main.go`, offline `ClassifyQuota` only (`g14stack/g14stack.log`).
- **Z (my round-04 P2 inputs).** `zretry/main.go` is my unchanged round-04 program; the only edit is its output
  path (`zretry/only-change.diff`). Log `zretry/zretry.log`.

**Not executed by me (SECONDARY producer evidence or unexecuted):** the full `go test ./...` (host, 607.6 s),
`-race`, `go vet ./...`, `go build`, the Windows/amd64 compile, the full skill suite, a rerun of the G14
`source-audit.py` slice extraction (I reran the harness, H), a rerun of my round-04 P1 programs, an end-to-end crash of a real
`parley` process, and re-rendering packets for phases 0/5/8. Windows runtime is unexecuted by anyone; no
Windows pass is inferred from compilation.

## User direction

- **Scope/reset answer** (`user-to-codex-1_…_scope-reset-answer.md`). My verbatim quote is in
  `review/round-03/claude-1.md:62-109`. Re-read in full this round; unchanged. The bounded zcode stderr rule
  and the display-clock rule are owner policy; I check the code against them and do not relitigate them.
- **Review-quota answer.** Quoted verbatim in `review/round-04/claude-1.md:147-168`, with my dated erratum at
  `:516-531` (`IMPL-ORGANIZER-BRIEF.md`, not `IMPLEMENTATION.md`).
- **Quota-standing and timeout-standing answers.** Quoted verbatim in my Phase-7 signoff in
  `review/round-04/consensus.md:158-235`. Both govern only organizer relaunches. This process made no retry or
  provider probe.
- **New owner direction since round 04:** none. No owner has waived AC2 or approved a close.

## Refutation attempts

PRIMARY = my execution or direct reading at the locator in `0ee1888`; SECONDARY = producer evidence I did not
rerun.

### G10–G14 (signed cycle-2 fix list)

| Item | Attempt | Result |
|---|---|---|
| G10 retries / receipt | Z (all round-04 P2 inputs); reading `quota_zcode.go:34-156` (record loop; receipt/start/high-water check at `:113`; final `first.Sub(receipt)` check at `:151`), `quota_framing.go` (`lastError` fingerprint vs. last `errors` entry); T (`TestQuotaCycle2*`, framing tests) | **Resolved.** PRIMARY Z: B-final, B2 (the two retained countdowns 176930/176890), C-final (two top-level records, one clock) are now eligible with reset `2026-10-04T22:14:57.003Z`, and D (separate real-clock collector writes) and E (6/4/0-second backoff, one write) are eligible at their real-clock resets; B-first and C-first reject as `contradictory inferred observation order/start/receipt` (receipt earlier than later attempts). Rejections now carry their own reason (`zcode evidence rejected: …`). |
| G11 knob-off path | D1–D3; reading `quota_manual.go`, `revision_snapshot.go`, `runner/telemetry.go:85-116`; T | **Partially resolved.** Exclusion grammar and known return work; retained vetoes cannot be withdrawn by a manual import. Catch-up join and kickoff-excluded return fail: R5-MAJOR-1. |
| G12 revision boundary | Reading `protocol/quota.go:91-105` (latest-receipt discriminator) and `membership.go` reconcile; T (`TestQuotaCycle2AppliedEditAndUnknownPendingNeverRewritten` and the projection-fault tests) | **Resolved.** An edit after an applied revision, even to an older set, blocks with a message naming `parley quota revise`; only a missing/invalid latest receipt with the prompt on that transition's before/after set is replayed. |
| G13 crash / manifest | Reading `membership/crash.go`, `pidlease/live_unix.go`, `revision.go`; T (native crash writer, settlement proof/replay, missing-manifest tests) on both volumes | **Resolved within V8.** Settlement needs same host/boot, dead supervisor and writer PIDs and an absent writer-owned group; it writes an immutable digest-bound record and never a terminal. Runner launches use `Setsid` (`runner.go:1124`, `launch.go:69`), so pgid == pid holds in practice. Manifests are validated before every commit. Residual: open question 3. |
| G14 evidence | Reading the harness, README and producer shape map; H (harness rerun), S (stack forms), Z | **Evidence prepared; AC2 still open:** R5-MAJOR-2. |

### Reservations V1–V9

- **V1. Met.** The threshold is `first.Sub(receipt) >= 60m` at terminal receipt (`quota_zcode.go`, final check).
  `TestQuotaCycle2TerminalThresholdAndObservationOrder` (61-minute first / 59-minute terminal → ineligible)
  passed in T. PRIMARY (S): a single record whose inferred observation is 61 minutes before the reset is eligible
  at a 61-minute receipt, and ineligible at a 59-minute receipt (`known reset below 60 minutes at terminal receipt`).
- **V2. Met.** Inferred observations are checked against record receipt, terminal receipt, recorded start and a
  high-water mark, each with one second of tolerance; the runner now passes the start time
  (`runner/telemetry.go:266-270`). PRIMARY Z (B-first, C-first) and T.
- **V3. Met** for top-level records (Z: C-final, D) and nested retries (Z: B, B2, E). Shape labels: see
  R5-MAJOR-2.
- **V4. Partly met.** Met: manual imports never become owner authority (`quota_retained.go:203` requires
  `r.Owner.Authority != nil`; `TestQuotaFixupOffModeReturnCannotAuthorizeWithdrawal` passed in T), and a known
  return needs only a `participants:` edit. Not met: the §5 catch-up join and the return of a kickoff-excluded
  agent differ from the pre-change CLI path (R5-MAJOR-1).
- **V5. Met.** PRIMARY (P): this idea's own unbracketed form, the bracketed form and a reason containing an em
  dash are accepted (the last was rejected in round 04); a trailing note, an inbox reference and hyphen separators
  are rejected with `manual exclusion d requires excluded: [d — reason — confirmed YYYY-MM-DD] (brackets optional;
  em dashes allowed in reason)…` (`revision_snapshot.go:96-98`). A known return without any `included:` marker is
  accepted (rejected in round 04). Docs list the rejected variants (`docs/quota-membership.md:70-75`).
- **V6. Met** (reading plus T): durable latest-receipt state, not prompt content, discriminates; replay keeps
  one terminal evaluation (N: `round.completed` stays 1).
- **V7. Met.** No COOPERATION wording changed in cycle 2: the packet source hash still equals the deck file and
  the skill addon hash, and command guidance lives in `docs/` and the skill.
- **V8. Met** (see G13).
- **V9. Met.** PRIMARY (H): the four source-derived modes reproduce byte-identically offline under a sterile
  environment; nothing loads zcode's entry point or configuration.

### Acceptance criteria AC1–AC21 (scoped)

- **AC1 PASS (scoped).** Packet hash equals the deck; `TestEmbeddedDefaultMatchesLiveDeck` passed in T; the
  skill addon records the same hash for `references/COOPERATION.md`; the changelog entry is present. Phase
  0/5/8 packets not re-rendered.
- **AC2 NOT MET** (R5-MAJOR-2).
- **AC3 PASS.** Negatives in T and Z; the identical-message control and contradictory forms still reject.
- **AC4 PASS** (reading `ClassifyQuota`: only zcode is supported).
- **AC5 PASS** (T: permutations, floor, duplicates, facilitator).
- **AC6 PASS** for kickoff filtering (T). Note D1: the same canonical-launch gate now also refuses never-excluded
  joiners (R5-MAJOR-1).
- **AC7 PASS** (T: preflight and runner bare 503 and noise).
- **AC8 PASS** (T).
- **AC9 PASS** for retained vetoes; knob-off return of new identities: R5-MAJOR-1.
- **AC10 PASS** (T: gates on prospective members).
- **AC11 PARTIAL.** Fault, pending and replay pass (T), and replay keeps one evaluation (N); the notice is
  re-created after archival (R5-MINOR-1).
- **AC12 PASS** (T, both volumes).
- **AC13 PASS** (T).
- **AC14 PASS** (T; reading `ReadHistory`).
- **AC15 PARTIAL** (R5-MINOR-1).
- **AC16 PARTIAL** (R5-MAJOR-1: an unlisted knob-off behavior change).
- **AC17 PARTIAL** (R5-MAJOR-1: same-idea return of a kickoff-excluded agent).
- **AC18 PASS** (T: quoted/tool/content negatives).
- **AC19 CONDITIONAL, honored** (`status: fix-up-cycle-2`).
- **AC20 PARTIAL.** My focused runs pass on both volumes; the full suite, race, vet, skill and Windows results
  are SECONDARY or unexecuted.
- **AC21 NOT MET.**

## Findings

### [MAJOR] R5-MAJOR-1: Policy-off catch-up join and kickoff-excluded return no longer work through the CLI

**Locators.** `internal/protocol/quota_manual.go:31-52` (at `:44` a never-known id needs a valid
`round-01/<id>.md` before the `participants:` edit is importable); `internal/runner/telemetry.go:85-116` (`beginLaunch`
refuses any canonical-artifact launch for an id outside the current set, with the message "excluded participant" at `:111`);
`internal/quota/record.go` (kickoff-excluded ids are never known).

**Executed (PRIMARY, identical program, shared volume and `/tmp`; `diff-*.log`).**

| Case | `27e42b8` | current |
|---|---|---|
| D1: `parley agents exec --agent e --artifact parley-deck/ideas/<slug>/round-01/e.md --yes` (late round-1 for joiner `e`, local stub) | exit 0; artifact written | exit 1, `excluded participant e cannot dispatch`; no artifact |
| D2: `participants:` edit adding `e`, then `consensus signoff --agent a` | exit 0 | exit 1, `open …/round-01/e.md: no such file or directory`; `status` shows an integrity gate |
| D3: re-include `d` (excluded with confirmation at kickoff) during round 1, then signoffs by `a` and `d` | both exit 0 | both exit 1 (`…/round-01/d.md: no such file`) |

**Why it matters.**

- The joiner cannot be dispatched to write its late round-1 (D1), and the join edit cannot come first (D2).
  Through the CLI, the §5 catch-up is now impossible; only an out-of-band file write works. The producer's
  differential wrote the late round-1 bytes directly to disk, so it never exercised this.
- A kickoff-excluded agent is a "new identity", so its return during round 1 needs a complete round-1 artifact
  first (D3). The pre-change CLI accepted the plain `participants:` edit.
- Each failure blocks the entire idea, with a raw file-not-found error: consensus status and signoff (executed),
  and the driver's `Before` (by reading).
- The new docs (`docs/quota-membership.md:77-82`) prescribe writing the late round-1 before the edit, but give no
  CLI route to write it (D1) and do not mention kickoff-excluded returns (D3).
- This is a third knob-off behavior change. FINAL §11 allows only C1 and the bare 503, and V4/G11 required the
  pre-change path as the oracle.

**Suggested fix (for the owner's decision; see Open questions).** For policy-off ideas: allow a canonical launch
whose target is the launching id's own late `round-01/<id>.md`; accept a join edit whose artifact is still
missing as pending membership (not an integrity block), importing it once the artifact validates; accept a
kickoff-excluded id's return by the plain `participants:` edit (it becomes known only from that revision, so the
kickoff signer rule is unchanged); and make the diagnostic actionable. Alternatively, list it
as a knob-off behavior change for the owner.

### [MAJOR] R5-MAJOR-2: Native AC2 evidence is still absent, and the realistic native shape is rejected (owner decision)

Carried R4-MAJOR-3, now concrete.

- The 24,833-byte original is not recovered; both retained native excerpts reject.
- The producer's source map locates the SDK default `console.error` sink with no inspect options, and the
  adapter's own retry loop with SDK `maxRetries: 0`. At Node's default depth 2, a real request's
  `requestBodyValues.messages` renders `[Object]`; `zcode-default-console-source-derived.stderr` rejects
  (`TestQuotaCycle2DefaultConsoleDepthStaysRejected` passed in T).
- Every accepted positive uses a synthetic one-frame stack (`sdk-framing-harness.cjs:27`,
  `e.stack=…split('\n')[0]+'\n    at offlineFixture…'`). Real multi-frame Node stacks are never exercised.
  PRIMARY (S): substituting real V8 frame forms into the fixture, bundle-path frames, `processTicksAndRejections`
  and `at async k7r (…)` are accepted, but `at async Promise.all (index 0)`, `at <anonymous>` and
  `at Generator.next (<anonymous>)` reject (`framing: unknown SDK stack framing`). Whether zcode's native
  stacks contain such frames is UNVERIFIED.
- Consequence: as shipped, the only supported recognizer most likely cannot classify a real zcode exhaustion.
  This fails closed and authorizes nothing wrongly, but the feature's practical benefit depends on it.

**Accepted and rejected shapes (current code):** accepted are complete allowlisted top-level or `RetryError`
records with a synthetic stack, `messages: []`, consistent 429 exhaustion and reset ≥ 60 minutes at terminal
receipt. Rejected are default-depth `[Object]`, partial native tails, mixed, contradictory or truncated records,
and some real V8 stack-frame forms (S).

**Owner decision needed (exact):** either (a) authorize one native zcode capture with invocation, start and
terminal-receipt facts, then accept or adjust the grammar against it; or (b) explicitly accept AC2 as unmet and
ship zcode support knowing it is effectively inert in practice; or (c) authorize a narrowly bounded grammar
extension (for example, inspect placeholders only inside `requestBodyValues`, plus listed real stack-frame
forms), validated by (a). I do not waive AC2.

### [MINOR] R5-MINOR-1: An applied transition's owner notice is re-created after the owner archives it

**Locator.** `internal/membership/membership.go:79` `reconcileLocked` → `:153` `publishNotice` → `quota.DurableWrite(path,
…, true)` for every batch on every `Before`.

**Executed (PRIMARY, `notice/{shared,local}.log`).** After a settled batch, 1 notice. The owner moves it to
`inbox/archived/`, leaving 0. After the next `Before`, there is 1 again. `round.completed` stays 1.

**Why.** §9.0 and AC11/AC15 call for one deduplicated notice per transition. Archiving is the normal inbox
workflow (§4), so every later driver tick re-surfaces an old exclusion as if it were new. **Fix:** treat an
applied receipt as proof of publication, or also check `inbox/archived/` before re-publishing.

## Prior finding dispositions

All are my own evaluation.

- **R4-MAJOR-1: resolved** (Z, T).
- **R4-MAJOR-2: partially resolved.** Exclusion, known return and veto protection are preserved; the catch-up
  and kickoff-excluded return residue is R5-MAJOR-1.
- **R4-MAJOR-3: open as R5-MAJOR-2** (owner).
- **R4-MINOR-1: resolved** (V5).
- **R4-MINOR-2: resolved** (G12).
- **R4-NIT-1: resolved** (rejections carry their own reason).
- **Round-04 open questions:** Q1 (inspect depth) answered by the producer's source evidence (SECONDARY):
  default depth with no override (I confirmed the default depth 2 and the `[Object]` rendering offline, H). Q2 (crashed writers) and Q3 (manifest before commit):
  answered by G13.
- **Round-03 dispositions:** unchanged from round 04 (CRITICAL-1/2, MAJOR-1–4, MINOR-1–5, NIT-1–4 as recorded
  there). R1 and R3 are still met (P: foreign/unknown host or boot never reclaimed; same host and boot with a dead PID reclaimed;
  committed authority survives archive and deletion, and modified, fabricated or mis-attributed answers fail); R2 has the residue in R5-MAJOR-1; R4 is still met for the allowlist, with
  the gap in R5-MAJOR-2.
- **Owner rulings** (bounded zcode stderr, display clock): not reopened.

## Open questions

1. **Stopping judgment (my signed condition).** R5-MAJOR-1 is a fresh MAJOR on G11's membership path. I
   recommend escalating to the owner with this trajectory (15 → 6 → 3 findings; both round-03 CRITICALs resolved in round 04, none since):
   either authorize a narrowly scoped cycle 3 for R5-MAJOR-1 and R5-MINOR-1, or record R5-MAJOR-1 as an
   accepted knob-off behavior change. Do not open cycle 3 by default.
2. **AC2:** the owner choice (a), (b) or (c) under R5-MAJOR-2.
3. **Proof-of-death residual (V8/R3, no finding; RECALL-level reasoning, UNVERIFIED).** "Same host and boot" uses
   `os.Hostname()` plus the boot id. On Linux, a writer in a container sharing the hostname but with its own PID namespace could look dead to
   `kill(pid, 0)`. Should the documentation state this limit?
4. **Unexecuted by me:** the full suite, race, vet, Windows, the skill suite, the G14 `source-audit.py` slice
   extraction and an end-to-end kill of a real `parley` process. Before an attended close, the owner should decide whether
   SECONDARY host evidence suffices for these.
