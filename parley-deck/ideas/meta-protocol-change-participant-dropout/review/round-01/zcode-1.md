---
agent: zcode-1
idea: meta-protocol-change-participant-dropout
review-round: 1
date: 2026-10-09
reviewed-commit: e4681cf5b9144ed786b269d8ab14b8d65b23d581
skill-commit: efe296c7acf13a147ab820ce6cbf8e6705b68691
---

# Review round 1 — zcode-1 (independent non-implementer review)

Protocol context attestation: `context_mode=full`,
`source_sha256=packet_sha256=091e6fb841685c85fa153f7e2f05e2329c3f28dbd9b0b4c72831bc88dc458bbf`,
`fallback_reason` absent. Verified by `shasum -a 256` on
`.parley-runtime/protocol-packets/full-phase6-deliberation-091e6fb841685c85fa153f7e2f05e2329c3f28dbd9b0b4c72831bc88dc458bbf.md`;
the body is identical to the live protocol supplied at launch (PRIMARY, executed 2026-10-09).
Active transport: `github-pr`. `meta/protocol-changelog.md` top entry is this idea's staged
2.16.0 hunk record (read, PRIMARY).

Scope. I reviewed both complete diffs — CLI `e4681cf5b9144ed786b269d8ab14b8d65b23d581` vs
design merge `431d6b0748ef7c16722c09899d1cc0e609e5eccb` (46 files, +2547/−341, mirrored at
`.parley-runtime/participant-dropout-launches/cli-review01.diff`) and skill
`efe296c7acf13a147ab820ce6cbf8e6705b68691` vs `352a475` (8 files, mirrored at
`skill-review01.diff`) — the frozen FINAL.md, IMPLEMENTATION.md, the three COOPERATION.md
copies, and the validator logs. Independent probes ran in a scratch copy (`git archive
e4681cf | tar -x -C /tmp/zcode-review01`); no product, source, test or canonical artifact of
either worktree was modified, and I wrote no one else's artifact or signoff.

## Summary

The implementation is a faithful, unusually well-tested realization of FINAL D1–D7. The
one-reducer extension is real: `quota.Policy` gains only an optional saved `trigger`
(legacy two-field JSON round-trips byte-identically — asserted by
`TestDropoutPolicyLegacyBytesAndStrictExtension` and re-read by me); the runner gains a small
durable per-step attempt ledger (`internal/runner/dropout.go`) keyed by idea/agent/logical
step over immutable invocation telemetry, with a PID lease, private 0600 preserved partials,
and strict two-attempt linkage (`appendParticipantAttempt` ordinal/retry-of checks;
re-validated at commit by `ValidateFailureEvidence`); control-plane classes (cancellation,
budget/protocol/telemetry refusal, tampering) are excluded from candidacy by
`ParticipantFailureClass` and block the batch instead; `membership.Settle` re-runs
`CheckGates` BEFORE `CommitBatch`, so `auto_implement` 3→2 still refuses with one
independent reviewer (tested in both gate variants); every rejoin path I could find or
construct (owner revision, manual edit, catch-up, recovery, prompt reconcile) refuses a
history-derived permanent drop even after owner downgrade to policy-off; the three
COOPERATION.md copies carry byte-identical normative hunks P1–P6 differing only in project
header zones, and the packet caps hold (phase-1 facilitator body 69,963 B ≤ 70,000 measured
in-run; SKILL.md 19,995 B ≤ 20,000).

I attempted to break the headline claims and broke two of them, both around restart
durability at the edges of the new ledger:

- **Z1 (MAJOR)** — the goal-check step's validator reads in-memory state, so a restart in
  the crash window between terminal telemetry and the validation receipt misrecords a
  valid first attempt as failed; the false reason is then bound into the immutable dropout
  evidence, and one failed relaunch drops a checker whose first attempt actually succeeded.
  This contradicts FINAL D2/AC1 ("a structurally valid artifact on either attempt wins").
- **Z2 (MAJOR)** — the kickoff readiness step key embeds a per-process random `ProbeID`
  that is never persisted, so a restarted `parley run` grants a fresh two-attempt probe
  budget, contradicting FINAL D2, the Idempotence & recovery section, AC2's restart clause,
  and the very §9.0 text this change ships.

Neither is a safety inversion (both err toward more retries/one illegitimate drop of a
non-protected participant, never toward keeping a failed agent or waiving a gate), but both
are acceptance-criterion violations with protocol-text contradictions and must be fixed or
disputed before close. Separately, AC13 is not yet claimable: the recorded final full-suite
runs each contain one non-reproducing, load-sensitive failure (details under "Pending
validation"), and vet/build evidence is still being written.

Verdict: **no CRITICAL findings; two MAJOR findings (Z1, Z2) and two MINOR findings
(Z3, Z4)**; implementation not ready for review consensus until Z1/Z2 are resolved and the
AC13 suite evidence is clean. Nothing here suppresses or narrows any future finding; the
five-cycle cap does not bound what I report.

## Refutation attempts

All commands executed by me on 2026-10-09 unless noted; provenance PRIMARY for executed
checks, SECONDARY where I rely on codex-1's recorded validator logs (marked). Scratch tree =
`git archive e4681cf` (byte-identical product sources to the reviewed commit).

**Against AC1–AC3 (two attempts, eligibility, control failures):**

1. `cd /tmp/zcode-review01 && go build ./...` — builds clean.
2. `go test ./internal/quota/ ./internal/telemetry/ ./internal/membership/ ./internal/protocol/ ./internal/consensus/ -count=1` — all `ok` (2.1s/2.6s/18.0s/0.6s/6.1s).
3. `go test ./internal/runner/ ./internal/app/ -count=1 -timeout 25m` — both `ok` (131s / 793s) on the scratch archive.
4. Read `internal/runner/dropout_test.go` end-to-end: real child processes (`/bin/sh` via telemetry shells) exercise exec exit codes, watchdog `no_first_output`, timeout, ACP started-exit, restart between failures (simulated driver exit before retry), cancellation/missing-protocol/shared-file/usage-event refusals, concurrent admission via the step lease, and 0600 private partials. The retry loop I feared (watchdog minting a third attempt) does not exist on the participant path: `runAgent`'s participant branch calls `runExecAttempt`/`runACPAgent` once per `RunParticipantStep` slot; the legacy first-output retry loop is only on the non-participant path (runner.go:586-617).
5. Probe A (Z1), new scratch test `TestZcodeProbeGoalCheckCrashWindowMisrecordsValidAttempt`: attempt 1 exits 0 with a preserved `GOAL-CHECK: PASS` verdict; deleting only `participant-result.json` (the crash window) and restarting with the amnesiac validator shape of driver_impl.go:838-841 plus a failing relaunch yields two-attempt dropout evidence whose attempt-1 `ValidatorReason` is `"missing or invalid GOAL-CHECK verdict"`. CONFIRMED (test PASSes by demonstrating the violation).
6. Probe B (Z2), new scratch test `TestZcodeProbeKickoffProbeIDResetsBudgetAcrossRestarts`: two consecutive `checkRoster` kickoff invocations against the same roster produced two independent settled ledgers — `gamma attempts=2 idea=kickoff-20261008T224903.726756000Z`, then `gamma attempts=2 idea=kickoff-20261008T224908.899237000Z` — four real failing children for one agent/step. CONFIRMED.

**Against AC4–AC5 (floor, protected roles, precommit gates):**

7. Read `quota.Evaluate` (quota.go:159-232) and `membership.Settle` (membership.go:282-349): protected-role failure blocks before any candidate check; designee/pin positive-usability is required only when a reduction is proposed; floor counts usable non-facilitators; `CheckGates` runs before `CommitBatch` and `CommitBatch` re-validates via `Batch.Validate` (history.go:314), which binds each candidate's `Failure.Idea` to the batch idea. `TestDropoutRoundUsesRealChildrenAndSettlesOnce` (gate variant) proves the auto_implement 3→2 refusal end-to-end with a real reduction refused and `QuotaBlocked` set.
8. Cross-idea evidence injection: `Evaluate` validates with idea="" (loose), but `Batch.Validate`→`ValidateFailureEvidence(e, b.Idea, …)` rejects at commit and on every later `ReadHistory` — fail-closed, no finding.
9. Escalation surface: `writeQuotaBlock` now publishes through `quota.PublishNotice` (safe-inbox semantics) and embeds `quota.OwnerOptions`; blocked kickoff publication failure prints the decision + options to stderr instead of failing the preflight path (preflight.go diff); `TestDropoutKickoffBlockCreatesSafeInbox` covers content, idempotency and the symlinked-inbox refusal.

**Against AC6–AC8 (immutable proof, permanent return, legacy compatibility):**

10. `TestDropoutEveryReturnPathAfterDowngrade` + my read of the guards: owner `Revise` rejoin, hand-edited `participants:` prompt, `Before` recovery, and `ManualCatchupTarget` all refuse the dropped id after an owner-authorized policy downgrade; `Dropped()` is derived from immutable kickoff/batch candidates, not `excluded:` markers. `TestDropoutKickoffNoticeReceiptReplay` covers crash-before-receipt replay and preservation of an owner-edited archived copy.
11. Legacy bytes: `TestDropoutPolicyLegacyBytesAndStrictExtension` asserts `{"enabled":true,"scope":"kickoff-and-mid-idea"}` unmarshal/remarshal byte-identity, rejection of `null`/`""`/`false`/`"future"`/duplicate/unknown trigger fields, and that the per-idea `false` opt-out still disables the new trigger. `ReconcileQuotaPrompt` filters any existing trigger line before re-projecting, keeping omission omission for legacy snapshots (protocol/quota.go diff); `InspectQuota` treats a trigger field without immutable history as an error and requires exact trigger/policy match.
12. Crash-derived attempts: `membership.ParticipantCrash` requires `started.json` + authenticated crash-settlement, synthesizes only an in-memory `invocation.crash-settled` (status failed, class `crash`), never writes terminal or invents success; a request with no start stays "unresolved participant request" and blocks (`TestDropoutCrashDerivedProofNeverInventsTerminal`).

**Against AC9–AC12 (dispatch surfaces, status, notices, protocol hunks):**

13. Wiring read: `runAgent` participant branch (runner.go:520-585) covers exec and ACP rounds with a stable `RoundLabel[/ArtifactName]` step id; `runSignoffAgent` wraps headless/interactive signoffs with prefix-ownership integrity checks and an own-suffix `AfterFailure` restore (AC does not apply — ACP signoffs route to the pre-existing manual flow, no child); `runDrafter` marks the draft started (and thus protected) BEFORE the bounded-launch guard can reject it; `GoalCheck` wraps the consult in the step ledger; `beginLaunch` refuses unbounded canonical launches under the new policy (`participant-failure policy requires a bounded logical step`), exempting read-only evidence-verification/consult; standalone `checkRoster` with no `QuotaPolicy` produces no evidence (`TestDropoutReadinessRetriesBeforeMembership` standalone half). `recoverParticipantSignoffs` settles an interrupted invalid own suffix recover-only, without a new child.
14. `diff` of the three COOPERATION.md copies at the reviewed commits: `parley-deck/COOPERATION.md` == skill `references/COOPERATION.md` byte-for-byte; defaults copy differs only in the generic workspace/date/roster header zone, as FINAL D5 requires. P1–P6 hunks match the packet I was launched with; the legacy quota predicate is retained verbatim-in-substance under its own label; §5/§0/Phase-0/Phase-5 wording matches the reviewed FINAL hunk table. Protocol changelog entry present; VERSION=1.52.0, version.json deckVersion=2.16.0 with `protocolSha256=091e6fb…` equal to the live packet hash; skill package 2.16.0.
15. Packet caps re-measured: built `cmd/parley` from the scratch archive and rendered `protocol packet --phase 1 --track deliberation --audience facilitator` → body file 70,037 B including the 74-byte wrapper comment; in-code `len(c.Body)`=69,963 B ≤ 70,000 guardrail (`TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail` PASSes on the committed copy; phase-6 facilitator body is 73,240 B but the ratified guardrail scope is phase 1, recorded per open item 2). SKILL.md = 19,995 B ≤ 20,000. `packet check`/drift tests green (also SECONDARY: packet-check-3.log `internal/protocolpacket ok`).

**Against AC13 (current-tree checks):**

16. Reproduction attempts on the one failing final-suite test — `go test ./internal/app/ -run TestEvidenceVerifierProductionClosure -count=1` (whole test, PASS 43s) plus `-run …/review-persistence-failure-recovery` three times (PASS 4.9-7.5s) on the scratch archive of e4681cf, and twice on a scratch archive of the design base 431d6b0 (PASS) — 6/6 isolated passes on both trees; the failure appears only under full-suite parallel load. Not a demonstrated regression, but the recorded suite run is not a pass.
17. Skill side: the failing `install is fleet-wide too: an immovable destination writes nothing anywhere` (`test/bidding-addon.test.js:2472`) passes in isolation on the skill worktree at efe296c (`node --test --test-name-pattern`, PASS, 24.5s); the skill diff touches no test, script, or bidding-addon file.

Remaining uncertainty: I did not execute the full 37-package `go test ./...` myself (I ran
all seven packages the CLI diff touches, green, plus the probes); I did not exercise a real
provider 429/503 against the new trigger (all evidence is local supervisor facts by design,
and the tests fabricate the full class matrix); cross-host/PID-namespace crash recovery and
Windows behavior remain disclosed untested surfaces, as FINAL states.

## Findings

### [MAJOR] Z1 — goal-check replay in the crash window misrecords a valid attempt as failed, enabling an illegitimate drop

- Where: `internal/app/driver_impl.go:838-841` (Validate closure over in-memory `res.Answer`),
  `internal/runner/dropout.go:108-123` (replay path re-validates and durably writes the
  receipt from the caller's `Validate`), `internal/telemetry/dropout.go:15-54` (the false
  `ValidatorReason` is frozen into `FailedAttempt`/`Excerpt`).
- What is wrong: `RunParticipantStep`'s restart contract requires `Validate` to be
  re-derivable from durable state — the round path re-reads the artifact file and the
  signoff path re-reads `consensus.md`, but the goal-check path validates
  `parseGoalVerdict(res.Answer)` over a zero-valued `ConsultResult` after a restart. If the
  driver dies between `invocation.Finish` (terminal.json) and `writeParticipantValidation`
  — a narrow but real window that includes the validator's own execution — the replay writes
  an immutable receipt marking the attempt failed with `"missing or invalid GOAL-CHECK
  verdict"` even though the child exited 0 and its preserved `goal-check.stdout.log` (already
  durably copied by the `Files` preservation) contains a valid verdict. One failed relaunch
  then produces two-attempt dropout evidence and the goal-checker is permanently dropped
  despite a valid first attempt.
- Why it blocks: violates FINAL D2 ("A structurally valid artifact on either attempt … wins
  over an ordinary execution failure and prevents dropout for that step") and AC1's first
  sentence; the committed evidence carries a factually false validator reason, which is
  exactly the "record the verbatim failure evidence" integrity duty (D4) this feature exists
  to guarantee. Probe A reproduces it end-to-end on the reviewed commit.
- Suggested fix: make the goal-check validator re-parse the preserved
  `goal-check.stdout.log` (or a verdict digest persisted next to the terminal record) in the
  replay path, or have `RunParticipantStep` refuse to mint a receipt from a validator whose
  inputs are not durable (leaving the attempt unsettled for explicit recovery) — any of
  these keeps the ledger honest without granting a third attempt.

### [MAJOR] Z2 — kickoff readiness two-attempt bound is replenished by a preflight restart (ephemeral ProbeID)

- Where: `internal/app/preflight.go:923-924` (`opts.ProbeID = "kickoff-" + store.NewRunID(time.Now())`
  per process, never persisted), `:958` (`Idea: opts.ProbeID` feeds
  `participantKey(idea, agent, "readiness")`), `internal/runner/dropout.go:87-88,226-228`.
- What is wrong: the durable ledger key for readiness probes embeds a random per-process
  batch id. A driver crash/timeout between a first probe failure and the batch settle — the
  precise "restart duplication" refutation target named in FINAL's Known risks — leads the
  relaunched `parley run` to probe under a new key with an empty ledger: two fresh attempts.
  Probe B measured two consecutive invocations each settling two attempts for the same
  agent (four failing children total). Mid-idea paths are unaffected (they key on the real
  idea slug).
- Why it blocks: contradicts FINAL D2 ("Run IDs, restarts, changed prompt inputs or a new
  driver process do not replenish it"; "Kickoff readiness uses the same two-attempt rule
  within the proposed batch"), the Idempotence & recovery section ("never grant a third
  attempt"), AC2's restart clause, the §9.0 P5 text this change itself ships ("Bind
  idea/agent/logical step: restart, run ID and prompt edits grant no third attempt"), and
  the skill SKILL.md guidance ("Durable idea/agent/logical-step identity prevents a third
  attempt across restarts/input/run changes"). No AC13 evidence statement can honestly
  claim AC2 while this stands. Blast radius is bounded (extra probe invocations of an
  already-failing agent; committed evidence remains internally valid), but the shipped rule
  and the behavior disagree.
- Suggested fix: persist the readiness batch id for the proposed kickoff (e.g. derived
  from the target idea slug/participants set, or written to a pre-idea scratch record under
  `.parley-runtime/` and reused by a relaunched identical proposal), so the same proposed
  batch shares one ledger across restarts.

### [MINOR] Z3 — replay `continue` for control-plane records skips valid-output recognition

- Where: `internal/runner/dropout.go:127-136`; `internal/telemetry/dropout.go:24-26`
  (class check precedes the `valid` check in `ParticipantAttempt`).
- A cancelled (or other control-class) invocation that nevertheless produced a structurally
  valid own artifact is `continue`d in replay and never recognized as valid; a new child is
  launched. At the round and signoff surfaces the pre-dispatch artifact-exists / prefix
  checks usually absorb this, so the exposure is narrow (e.g. goal-check, overwrite
  dispatches). It errs toward an extra bounded attempt, never toward an unsafe drop or a
  waived gate, hence MINOR. Worth aligning: check validity before control-class skipping,
  or document that control-class terminals are never success carriers.

### [MINOR] Z4 — phase-1 facilitator packet body is at 99.95% of its ceiling; SKILL.md at 19,995/20,000 B

- Where: `internal/protocol/defaults/COOPERATION.md` (P5 hunk),
  `internal/app/facilitator_packet_live_test.go:21` (70,000 B guardrail), skill
  `skills/parley-deck/SKILL.md` (19,995 B vs the 20,000 B core cap).
- Both caps hold on the reviewed tree (I re-measured; see refutation 15), and AC12 requires
  them unchanged — but five bytes of SKILL.md headroom and 37 B of facilitator headroom
  mean essentially any future additive protocol/skill sentence breaches a cap and will force
  a compaction round. Not a defect today; recording it so the next protocol edit plans for
  it rather than discovering it in CI.

### Pending validation (not code findings — close must not be claimed on current evidence)

- AC13: the FINAL host run (`full-host-tests-final.jsonl`, complete: 37/37 packages
  reported, 33 pass / 1 fail / 3 skip) has `internal/app` failing on
  `TestEvidenceVerifierProductionClosure/review-persistence-failure-recovery` under
  full-suite parallel load. I could not reproduce it in six isolated runs across both the
  reviewed commit and the design base, so I classify it load-sensitive rather than a
  regression — but the recorded suite evidence is a failure and a clean full-suite rerun is
  required before close. The skill run (`skill-tests-final.log`) is 398/399 with the
  bidding-addon install test failing (non-reproducing in isolation; no test/script change in
  this diff), and because `node --test` failed, the chained `run-python-tests.js` and
  `build-addon-manifest.js --check` stages did not run — their evidence is missing, not
  pass. `vet-candidate.log` was still empty (running) at review time; the recorded vet/build
  evidence in IMPLEMENTATION.md is from a preceding tree and is explicitly superseded.
- AC14: delivery has not started (releases, Homebrew, WinGet PR, installs, core 2.16.0
  staging/publication) — expected at this phase; nothing here simulates or claims it.

## AC assessment (current tree, my own evidence)

| AC | Assessment |
| --- | --- |
| AC1 | Substantiated by tests + reads for the class matrix, valid-artifact/BLOCK precedence; **violated in the Z1 goal-check crash-window corner**. |
| AC2 | Substantiated mid-idea (durable bound, watchdog slot sharing, exec/ACP, restart between failures, input/run-id stability); **violated for kickoff restarts (Z2)**. |
| AC3 | Holds — control classes excluded and batch-blocking; tampering preserved; crash settlement never invents success (tests + code reads). |
| AC4 | Holds — whole-batch evaluation, positive-usability floor incl. designee/pin, protected-role blocks, no role-only seat (code + tests). |
| AC5 | Holds — precommit `CheckGates` before commit; auto_implement 3→2 refused end-to-end; one actionable escalation with owner options; safe inbox + stderr fallback tested. |
| AC6 | Holds — paired immutable attempts with linkage re-validated at commit and on every history read; private 0600 partials; obligations/vetoes retained (tests). |
| AC7 | Holds — every return path I attacked after downgrade refuses; kickoff candidates included; next idea probes afresh (tests + reads). |
| AC8 | Holds — legacy bytes/hash identity asserted; strict trigger parsing fail-closed; per-idea opt-out precedence kept; new default only for new ideas. |
| AC9 | Holds for kickoff/mid-idea/round/review/signoff/drafter dispatch; standalone preflight report-only verified; protected drafters escalate; goal-check carries the Z1 corner. |
| AC10 | Spot-verified (quotaSurface trigger + permanent label, unchanged wait exit contract, reads never repair); not exhaustively re-probed — residual, not a finding. |
| AC11 | Holds — receipt-idempotent notices, kickoff replay incl. owner-edited archived copies, block dedup across run IDs, delivery failure diagnostic-only (tests). |
| AC12 | Holds on inspection and re-measurement — identical hunks on all three copies, project zones preserved, changelog/versions staged, caps unchanged and respected (Z4 headroom note). |
| AC13 | **OPEN** — final suite runs each carry one non-reproducing failure; chained skill stages and vet/build evidence outstanding; see Pending validation. |
| AC14 | **OPEN** — delivery not started; owner-only npm/TTY publication steps remain owner actions. |

## Open questions

1. For Z2: is there an intended reading of "within its proposed readiness batch" under
   which a restarted kickoff is a new batch? FINAL D2's restart sentence and the shipped
   §9.0 text give no such carve-out; if the owner wants one, it is a protocol-text change,
   not an implementation detail.
2. For Z1: should the shared helper enforce durable-state validators structurally (refuse
   to mint replay receipts when the validator cannot attest its inputs), rather than fixing
   only the goal-check closure? The round/signoff callers already satisfy it; a structural
   guard would prevent the next caller from repeating the mistake.
3. Note for the fix-up planner, not a defect: ACP-mode consensus signoffs are not wired
   (pre-existing manual routing), so the new ledger simply does not see them; if signoff
   ACP wiring lands later it must join the participant wrapper explicitly.

## Verdict

The architecture, the protected-role/floor/gate semantics, the permanent-return enforcement
and the protocol/skill edits are implemented as designed and survive my refutation attempts;
the two MAJOR findings are both restart-durability defects at the ledger's edges (goal-check
crash window; kickoff ProbeID), each with a narrow fix that does not require redesign. I
file this review with **no CRITICAL, two MAJOR (Z1, Z2), two MINOR (Z3, Z4)** and AC13/AC14
open as described. Review consensus should wait for Z1/Z2 fixes (or an owner ruling if
either is disputed) and clean current-tree suite evidence; the brief's attended-close
preconditions are not yet met. This verdict is mine alone; codex-1 implements, and nothing
in this file waives, narrows or pre-authorizes anything for anyone.
