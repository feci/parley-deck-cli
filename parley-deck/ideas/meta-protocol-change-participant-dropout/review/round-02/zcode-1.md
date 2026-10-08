---
agent: zcode-1
idea: meta-protocol-change-participant-dropout
review-round: 2
date: 2026-10-09
reviewed-commit: f8f4f1f17bc99832aeb46c0a8ee44c6331325d1b
skill-commit: efe296c7acf13a147ab820ce6cbf8e6705b68691
responding-to: [review/round-01/zcode-1.md, review/consensus.md]
---

# Review round 2 — zcode-1 (independent non-implementer review of fix-up cycle 1)

Protocol context attestation: `context_mode=full`,
`source_sha256=packet_sha256=091e6fb841685c85fa153f7e2f05e2329c3f28dbd9b0b4c72831bc88dc458bbf`,
`fallback_reason` absent. Verified by `shasum -a 256` on
`.parley-runtime/protocol-packets/full-phase6-deliberation-091e6fb841685c85fa153f7e2f05e2329c3f28dbd9b0b4c72831bc88dc458bbf.md`
at session start; the body is identical to the live protocol supplied at launch (PRIMARY, executed
2026-10-09). Active transport: `github-pr`.

Scope and trees. I reviewed the complete product diff `431d6b0..f8f4f1f` (verified byte-identical,
modulo `index` lines, to `.parley-runtime/participant-dropout-launches/cli-review02.diff` =
`git diff 431d6b0..f8f4f1f` excluding only `parley-deck/ideas/`), the cycle-1 diff
(`cycle1-review02.diff`, product content = `85ea372..cb78e9f`: CHANGELOG.md, docs/quota-membership.md,
app/{driver_impl,preflight,preflight_liveness,dropout_test}.go, runner/{dropout,dropout_test}.go),
and the unchanged skill diff (`skill-review01.diff` == `git diff 352a475..efe296c` in the skill
worktree, verified byte-identical). `f8f4f1f` is deck-only over `cb78e9f` (IMPLEMENTATION.md,
organizer-notes.md), so the reviewed product tree is `cb78e9f`'s. While my review was in flight the
producer added deck-only `9b6edba` (validation record + an inbox note to me); no product source has
changed since `cb78e9f` and my checks below ran on a `git archive f8f4f1f` scratch tree, whose
product content equals the current tree. I read ORGANIZER-BRIEF.md, frozen FINAL.md, IMPLEMENTATION.md,
my round-01 review and both cycle-1 consensus signoffs in full. I modified no product source, test,
or canonical artifact of either worktree; my probes live only in the scratch tree; both worktrees
verified clean (`git status --porcelain`) after all my runs.

## Summary

Fix-up cycle 1 lands both agreed code fixes (Z1, Z3) and the ACP fixture repair exactly as the
signed consensus specified, and my attempts to re-break them failed — which is the desired outcome.
`ValidateRecord` validators are now bound to the immutable invocation identity and reconstruct their
inputs from the invocation's own retained output (`ParticipantOutput`: private per-invocation copy
first, the ORIGINAL run's public log as fallback — never the new driver's run directory —
non-regular files fail closed); receipts bind `SourceSHA256`, so changed or missing previously-valid
evidence blocks instead of authorizing a replacement; a nonzero or non-exited child can no longer
establish goal completion merely by printing PASS; control-class terminals with valid output are
blocking (never dropout, never completion). My four new adversarial probes (cancelled valid PASS
through the missing-receipt crash window; retained-copy symlink; public-log tamper after a valid
receipt; nonzero-exit readiness failures under a stable ProbeID across restarts) all PASS — i.e. the
implementation survives them. The Z2 withdrawal and Z4 deferral both still hold on my re-evaluation.

On AC13 I can now distinguish finished passes from the environment failure: the producer's first
full-host run (`full-host-tests-cycle1.jsonl`, external-volume TMPDIR) failed 7 tests in 4 packages
and is retained as non-evidence; my independent timeline analysis confirms those failures are
external-mount/git/lock/path-space breakage of fixtures untouched by this idea, not product
regressions. The replacement run (`full-host-tests-cycle1-local.jsonl`, local TMPDIR/GOTMPDIR,
external GOCACHE) completed **exit 0 at f8f4f1f** — 37 packages, 3262 test passes, zero failures —
which I verified from its result JSON and raw log myself. With my own focused runs, vet, build, and
a full independent skill suite rerun (399 Node + 54 Python + 6 manifests, zero failures) all green,
AC1–AC13 are evidenced on the current tree. AC14 (delivery) correctly remains post-close.

Verdict: **no CRITICAL and no MAJOR findings remain open; no new CRITICAL/MAJOR found this round;
zero agreed fixes are required from this review.** The zero-fix final review consensus and the
brief's attended-close prerequisites are met on the merits once that consensus is drafted and signed
by both participants — details and the one process caveat (evidence-file naming) under "Close
prerequisites".

## Evaluation of signed dispositions

### Z1 (was MAJOR) — durable goal-check/readiness replay validation — CONCUR: fixed

- Fix verified in source (PRIMARY): `driver_impl.go` GoalCheck now passes
  `ValidateRecord(telemetry.Record)` whose fallback path is
  `runs/<r.Metadata.RunID>/agents/<checker>/goal-check.stdout.log` — the attempt's ORIGINAL run —
  with kimi-adapter unwrapping and `SourceSHA256` over the raw retained bytes;
  `preflight.go`/`preflight_liveness.go` reconstruct readiness from per-invocation
  `probe.stdout.log`/`probe.stderr.log` (64 KiB cap, regular-file check, digest, outcome-derived
  exit/timeout), and `runner.ParticipantOutput` reads retained-first, durably re-preserves a
  fallback read, and fails closed on non-regular evidence.
- Replay flow (PRIMARY, read `runner/dropout.go` current tree): a missing last receipt is
  re-validated from durable state and written; a previously-valid receipt is re-validated every
  replay and blocks on `!Valid`, `Integrity`, or `SourceSHA256` mismatch ("previously valid
  participant output changed; no new attempt authorized"); completion requires
  `outcome.Status == "process-exited"` and exit 0 — a failed child printing PASS returns
  `goal-check checker failed; completion is unverified`; FAIL is preserved as valid dissent
  (no dropout, verdict surfaced).
- Committed tests pass in my independent run and both full-host runs:
  `TestDropoutGoalCheckReplaysActualOutputAcrossRuns` (PASS/FAIL/malformed × terminal-before-receipt
  × before-preservation × changed run ID × a planted contradictory log in the NEW run's directory),
  `TestDropoutGoalCheckValidEvidenceCannotChangeOrExcuseFailure` (changed-valid/missing/nonzero),
  `TestDropoutReadinessSameBatchReplay` (PONG/malformed/changed-valid/missing, incl. durable PONG
  receipt repair and retained paired-attempt proof).
- My new probes (scratch only, all PASS): cancelled-child-with-valid-PASS through the missing-receipt
  crash window blocks on first call and replay with no dropout and no replacement child;
  a symlinked retained copy fails closed; tampering only the original public log cannot flip a
  retained valid PASS.

### Z3 (was MINOR) — validity before control-class skip — CONCUR: fixed

- Replay now validates retained output before any control-class skip: a control-plane terminal with
  a valid receipt re-validates and, on the `ParticipantAttempt` control error, returns
  "valid participant output retained after a control-plane refusal; no replacement authorized";
  integrity receipts remain hard blocks ("preserved participant integrity failure"); an invalid
  control-class record still `continue`s so a repaired pre-dispatch refusal can get its first real
  child (`TestDropoutRepairedRefusalWithoutValidOutputMayDispatch` asserts the refusal consumed no
  child attempt).
- `TestDropoutCancelledChildWithValidArtifactNeverReplaced` (child- and parent-cancel variants)
  proves a cancelled child's valid artifact is never replaced across two replays; the replay
  assertions added to `TestDropoutCancellationAndIntegrityNeverReduce` hold. Cancellation and
  integrity never produce dropout evidence (`h.Dropped` asserted false in my probes) and never
  completion success (GoalCheck false asserted).

### Z2 (withdrawn in my cycle-1 signoff) — re-evaluated: disposition STANDS

I re-read frozen FINAL D2 and the signed design consensus again this round (PRIMARY): the
non-replenishment sentences are keyed to the per-idea ledger, kickoff readiness is expressly scoped
"within the proposed batch", and the restart bound is staged only "after an idea exists". My round-01
probe exercised two complete proposals (next-proposal probing), not a mid-batch resume. The
disposition's own commitments all landed: the stable-ProbeID same-batch replay test exists and
passes (PONG/malformed/changed-valid/missing), my new probe extends it to the nonzero-exit
process-failure class across two `checkRoster` invocations (exactly two children total, paired
attempts retained), and the boundary is visibly retained in IMPLEMENTATION.md, CHANGELOG.md and
docs/quota-membership.md ("resuming an uncreated proposal across commands is not yet a supported
lifecycle"). The persistent pre-idea proposal/resume/abandon lifecycle remains a recorded TBD
follow-up. No reopening.

### ACP started-exit fixture — CONCUR: fixed as agreed

`TestDropoutExecWatchdogTimeoutAndACPShareTwoSlots`: `acp-started-exit` launches via
`LaunchACP`/`ACPArgs ["-c","exit 7"]` under a 5 s ceiling; the deliberate timeout case keeps 150 ms
(`exec sleep 5`, no exit) and watchdog keeps 2 s + 30 ms first-event; both attempts share the
ceiling and the existing two-attempt cap (exactly two terminal records asserted); the class
assertion now prints the actual class. Stability: 3 consecutive focused repetitions in my scratch
run (10.7 s total) plus passes in both full-host runs. Runtime failure classification is unchanged.

### Z4 (deferred maintenance risk) — re-evaluated: STANDS as accepted limitation

Re-measured by me on the current tree: SKILL.md 19,995 B ≤ 20,000; deck COOPERATION.md ==
skill `references/COOPERATION.md` byte-for-byte; defaults copy differs only in the generic header
zone; phase-1 facilitator packet guard tests (`protocolpacket`, `protocol`) pass in both full-host
runs and my package run; version metadata consistent (VERSION 1.52.0, deckVersion 2.16.0,
protocolSha256 == live packet hash). No limit raised, no criterion waived. The headroom risk is
correctly recorded as a TBD compaction budget for the next additive change.

## Refutation attempts

Refutation scope: AC1–AC14, all on the current tree; every attempt was executed by me unless
marked SECONDARY.

Environment for all my commands: scratch tree `git archive f8f4f1f` at
`.pd-dropout-th904gpe/zc2-scratch`; storage per updated `test-storage.json`
(TMPDIR=GOTMPDIR=/private/var/tmp/pd-dropout-tests-201b7y78, GOCACHE external shared);
`go vet ./...` PASS; `go build ./cmd/parley` PASS (my `/tmp/zc2-parley` reports `parley 1.52.0`).

| # | Attempts and result |
| --- | --- |
| AC1 | Class matrix, valid-artifact/BLOCK precedence and the Z1 corner re-verified: runner `Dropout` suite ok (6.0 s) and app `Dropout\|GoalCheck\|Readiness` suite ok (45.7 s) in my run; 130 dropout-scope tests pass in BOTH full-host runs; cancelled-valid-PASS probe (missing receipt) blocks without dropout; changed/missing previously-valid goal-check evidence blocks (committed tests + code read). |
| AC2 | Restart/input/run-ID stability re-proven by the cycle-1 replay tests across changed run IDs; watchdog slot sharing and exec/ACP parity in `TestDropoutExecWatchdogTimeoutAndACPShareTwoSlots` (3× count run ok); my same-ProbeID nonzero-exit probe holds exactly two children across two `checkRoster` invocations (pre-idea batch per the Z2 disposition; per-idea bound unchanged since round 1). |
| AC3 | Cancellation/control/integrity exclusion re-verified: cancelled-child tests (both variants) + replay assertions in `TestDropoutCancellationAndIntegrityNeverReduce` + my cancelled goal-check and symlink probes — always blocking, never dropout (`h.Dropped` false), never completion (GoalCheck false). |
| AC4/AC5 | Floor/protected/gate code (`quota.Evaluate`, `membership.Settle`, `CheckGates` before `CommitBatch`) is unchanged by cycle 1 — round-1 verification carries; membership (17.1 s), quota (2.4 s) package tests green in my run and both full-host runs. |
| AC6/AC7 | Immutable paired-attempt evidence, private partials, permanent-return guards unchanged by cycle 1; quota/telemetry/membership tests green in my run and both full-host runs; round-1 probes (every return path after downgrade) unrefuted. |
| AC8 | Legacy policy byte-identity and fail-closed trigger parsing unchanged (protocol package green, 0.7 s, plus both full-host runs; round-1 test-level verification carries). |
| AC9 | Goal-check and readiness dispatch now use the record-aware validators (the two surfaces cycle 1 changed) — verified by the committed tests plus my probes; standalone-preflight report-only path unchanged since round 1; protected drafter/implementer gating code unchanged. |
| AC10 | **Beyond spot-check this round, with my own binary from the scratch tree**: `status --json`, `organizer brief` and `wait --for review` all agree on membership `[codex-1, zcode-1]`, implementation `fix-up-cycle-1` at head `cb78e9f`, and the pending `review/round-02/zcode-1.md`; the stale advisory run surfaces as `STALE / No recoverable action` (known D6 phase-pointer gap) without any repair; before/after `find -newer` + `git status` snapshots prove zero writes into `parley-deck/` from all three read-only verbs; `wait -timeout 5s` exits 3 with a loud outstanding list (contract unchanged). The idea has no automatic-exclusion history (kimi-1's exclusion is the owner-authorized frontmatter path), so the empty quota surface is correct. |
| AC11 | Notice/receipt/replay code unchanged by cycle 1; membership package green; round-1 verification (receipt-idempotent notices, kickoff replay, owner-edited copies, block dedup across run IDs) carries. |
| AC12 | Re-verified by me today: the three COOPERATION.md copies agree (deck == skill references byte-for-byte; defaults differs only in the generic header zone); protocol changelog entry present; VERSION 1.52.0 / deckVersion 2.16.0 / skill 2.16.0; `protocolSha256` equals the live packet hash; packet-guard tests pass in both full-host runs. Caps unchanged (19,995/20,000 skill; facilitator packet guard green). |
| AC13 | **Met on the current tree.** My independent runs: vet PASS, build PASS, quota 2.4 s, telemetry 2.7 s, membership 17.1 s, protocol 0.7 s, consensus 11.0 s, runner `Dropout` 6.0 s, app `Dropout\|GoalCheck\|Readiness` 45.7 s, ACP fixture ×3, four adversarial probes — all PASS. Full host (adjudicated below): `full-host-tests-cycle1-local-result.json` **exit 0 at f8f4f1f**, 731.1 s, 37 packages (34 pass, 3 no-test skips), 3262 test passes, zero failures — recounted by me from the raw jsonl. Producer's vet/build exit 0 (`vet-build-cycle1-result.json`, cb78e9f) consistent with mine. Skill (my own full `npm test` on unchanged efe296c): 399 Node + 54 Python + 6 manifests, zero failures — independently reproduces `skill-tests-rerun.log`. Windows limits remain disclosed, not passed. |
| AC14 | Not started (no release, formulae, WinGet PR, installs, core staging) — correct at this phase; remains post-close owner+implementer work. Nothing simulates or claims it. |

**Full-host evidence adjudication (honest record).** The first cycle-1 full-host run
(`full-host-tests-cycle1.jsonl`, external-volume TMPDIR per the then-current `test-storage.json`)
recorded 7 test failures across driver/budget/runner/app: git `worktree add` "invalid reference"
failures, a git-shim fixture whose `seen-env` file never appeared, a flock-target non-refusal, a
concurrency-cap count change, an autodrive fixture child exit, and an organizer-brief write into a
read-only deck. My independent analysis: (i) every failing test lives in code UNTOUCHED by this
idea (no `internal/budget` or `internal/driver` file in the product diff; the three app failures are
in pre-existing fixtures); (ii) all 7 passed in BOTH earlier full runs on internal-disk temp
storage and in the replacement run; (iii) all 130 dropout-scope tests passed in the failing run
itself; (iv) the failure timestamps (01:15:52–01:19:03 local) precede the second run's start
(01:19:11), so the two overlapping producer runs did not cause them — the external shared-folder
mount (spaces in path, no flock semantics, git ref/worktree quirks, and Go 1.27 `testing.TempDir`
honoring GOTMPDIR) did. This matches the producer's own diagnosis, which I did not adopt
uncritically — I verified the timeline, the test selection, and the pass/fail splits myself. Run A
(`exit 1`, result JSON present) is correctly retained as non-evidence; run B (local TMPDIR/GOTMPDIR,
external GOCACHE) is the qualifying pass. Process caveat, not a code defect: my launch contract
named `full-host-tests-cycle1-result.json` as the only completion evidence; that file records the
FAILED run, and the qualifying evidence is the `-local` pair. The producer disclosed and
re-contracted this in the 9b6edba record and the inbox note, retaining both runs; the final
consensus/IMPLEMENTATION should keep pointing at the `-local` result (or re-run the suite under the
final storage path) so the evidence chain is unambiguous at close.

## Findings

No new CRITICAL or MAJOR findings. No new MINOR findings. Historical findings: **Z1 and Z3 are
resolved by cb78e9f and verified above; Z2 remains withdrawn under the FINAL D2 pre-idea/idea
boundary with its disposition commitments landed; Z4 remains a deferred maintenance limitation.**
Zero agreed fixes are requested by this review.

Non-finding observations for the record (no action required for this idea):

1. Residual corner, inherent to the design: if a valid attempt's receipt, retained private copy AND
   original public log are all destroyed before replay, no durable witness of validity exists; the
   attempt re-validates invalid and may consume the second slot (bounded by the two-attempt cap; no
   third child is possible). Consistent with FINAL D2's evidence duties; noting it so a future
   design that wants verdict digests inside terminal telemetry does so deliberately.
2. Pre-existing test-robustness gap exposed by the storage incident (outside this idea's diff): an
   inherited gitprobe fixture embeds an unquoted TMPDIR path (breaks on spaces), and several
   budget/git fixtures assume internal-disk filesystem semantics. A candidate small hygiene
   follow-up idea; explicitly NOT a finding against this implementation.
3. The two overlapping producer full-host runs (01:15:40 and 01:19:09 local) were untidy process
   hygiene, but the timeline shows the second run did not cause the first run's failures, and the
   qualifying run is clean; no evidence was relabeled.

## Open questions

0. Self-correction of my round-01 artifact (invited by the consensus's validation note, my file
   stays unedited): I there wrote "review-persistence-failure-recovery"; the actual failing subtest
   was `TestEvidenceVerifierProductionClosure/report-persistence-failure-recovery`. The correction
   changes nothing else in that finding — and the subtest passes in the qualifying exit-0 full-host
   run, closing that pending-validation item.
1. For the final consensus drafter: record explicitly that the qualifying full-host evidence for
   cycle 1 is `full-host-tests-cycle1-local-result.json` (exit 0, f8f4f1f) and that
   `full-host-tests-cycle1-result.json` (exit 1) is retained environment-failure history — so the
   close record cannot be read as claiming the failed file as its pass.
2. Nothing else is open from my side. The TBD items (pre-idea proposal lifecycle; protocol/skill
   headroom compaction budget; channel delivery under AC14) are already recorded as follow-ups and
   do not gate this idea.

## Verdict and remaining uncertainty

Fix-up cycle 1 is complete and correct against the signed fix plan and the frozen FINAL: both code
fixes and the fixture repair landed exactly as agreed, all agreed-fix tests plus my four new
adversarial probes pass, every FINAL AC1–AC12 criterion is independently evidenced on the current
tree, and AC13's focused, full-host, vet/build and skill legs are all green — the full-host leg via
the qualifying exit-0 run I verified myself, with the failed external-storage run honestly retained
as non-evidence. AC14 remains post-close delivery.

**Zero-fix final consensus and attended close prerequisites**: with this review filed, no CRITICAL
or MAJOR remains open anywhere in the review lifecycle (round 1's Z1/Z3 fixed and re-verified; Z2
withdrawn with commitments landed; Z4 a recorded deferral), and this round requests zero agreed
fixes — so a Phase-7 final consensus listing **zero Agreed fixes** is supportable, and after both
signoffs (codex-1, zcode-1) the brief's attended-close preconditions (both final review-consensus
signoffs, no open CRITICAL/MAJOR, current-tree AC evidence recorded, standing owner authority) are
met, subject to open question 1's evidence-naming record. `strict_gate` is not set on this idea; I
have nonetheless reviewed the complete diff at the current commit, not only fix-up deltas, and my
probes were not limited by the fixture hints. I do not sign here — that is the final consensus's
step, after codex-1 reads this round.

Remaining uncertainty: I did not exercise a real provider 429/503/400 against the new trigger (all
evidence is local supervisor facts by design; the test matrix fabricates the classes); cross-host /
PID-namespace crash recovery, real hosted error timing, and native Windows remain disclosed
untested surfaces, as FINAL states; my probes ran on the archived f8f4f1f product tree (product-
identical to the current tree, re-verified by git) rather than the live worktree, to keep the
worktree untouched.

<!-- Correction 2026-10-09 (zcode-1, own artifact): the section heading above was normalized from
     "## Refutation attempts — AC1–AC14 (all current-tree; executed by me unless marked SECONDARY)"
     to the exact "## Refutation attempts" required by the CLI review-artifact validator; the scope
     annotation moved into the paragraph under the heading. No finding, evidence, or verdict changed. -->
