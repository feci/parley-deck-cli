---
agent: zcode-1
idea: meta-protocol-change-participant-dropout-review-gate-timing
review-round: 2
date: 2026-10-09
reviewed-commit: b89e2abaa056813a4b238c2fa4278195b0f7096b
scope: full
responding-to: [review/round-01/zcode-1.md]
---

## Summary

Full-scope re-review after fix-up cycle 1. Reviewed the entire CLI diff
85c5a8f..b89e2ab (implementation a26f588 plus the cycle-1 fix commit b89e2ab; the later
handoff commit ef652e6 adds only idea artifacts — `git diff --stat b89e2ab..ef652e6 --
internal/ docs/ README.md CHANGELOG.md VERSION go.mod cmd/` is empty, verified this
launch) and the entire skill diff dbdb919..74cc831, re-read the complete live protocol
packet, FINAL.md, IMPLEMENTATION.md, the signed review/consensus.md (revision 2, both
ACCEPTs), source-context/implementer-response-minor1.md, my round-1 review and phase-7
adjudication, and every validation artifact named in the brief.

**Launch attestation:** context_mode=full,
source_sha256=packet_sha256=acbd4dbc0c0702bc191176bb80bcee32c5093c8b4ebbee42e036df6a9b7d1137,
no `fallback_reason` (`review02-attestation.json`, PRIMARY; I also re-hashed the packet
file with `shasum -a 256` — match — and read it in full, 1,552 lines / 125,132 B). This
launch is the recorded configured-native review fallback (native start
`review02.native-start.json`, PID 40945, 2026-10-09T08:37:50Z, prompt sha 08474bc0…)
because the driver has no recoverable review-dispatch action and the original driving run
has the D6 refusal. I authored only this file; no product files, no commits, no GitHub
messages, no consensus drafting or signing.

**Verdict: no CRITICAL, no MAJOR, no MINOR. Two NIT.** Every agreed fix from the signed
cycle-1 plan is implemented exactly as signed and pinned by real tests. The final
current-source validation at b89e2ab is now complete and I verified it independently:
full Go suite exit 0 (my own recount: 34/34 packages, 3,330 distinct passing test/subtest
IDs, 4 skips, 0 fail events; jsonl sha256 04d2f052… matches the durable receipt, which is
byte-identical across native task temp and shared storage), build/vet exit 0, and the
skill suite at the unchanged 74cc831 stands on its verified v3 pass. My round-1 MINOR
withdrawal survives a deliberate re-challenge (below); the deferred shared-volume receipt
follow-up did not recur this cycle and its deferral condition ("open it if the volume
repeats this") remains unmet.

## Refutation attempts

All attempts are **PRIMARY** — I read the cited code at b89e2ab and ran the cited
commands myself this launch. Paths are relative to the CLI worktree
`/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/review-gate-timing`; the
skill worktree is the sibling `review-gate-timing-skill` (HEAD 74cc831, clean).

### Cycle-1 agreed fixes, verified against the signed plan

- **Fix 1 (preserve the missing-manifest integrity gate).** quota_signoff.go:62-66
  loads the prior run manifest and returns the error with no `os.IsNotExist` catch —
  exactly the signed "do NOT add" instruction. `validateManifests`
  (membership.go:157-197, unchanged from baseline 85c5a8f) still requires the current
  run, the kickoff run and every batch run, and still aborts on any load error or
  foreign identity ("contradictory quota manifest", :182). Re-challenged my own
  withdrawal of round-1 MINOR-1 rather than resting on it: the new earlier load can only
  fail where `validateManifests` would also fail, because `priorRun` (kickoff, or last
  batch, quota_signoff.go:58-61) is always a member of the required set
  (membership.go:160-163); a *foreign* prior manifest passes the new load (valid JSON)
  and is then caught by `validateManifests`' identity check — now proven by a failing
  fixture, not just by my reading (see Fix 3). The withdrawal stands: baseline 85c5a8f
  blocked the same missing prior manifest at `validateManifests` (verified in my
  adjudication, `git show 85c5a8f:internal/app/quota_signoff.go` +
  membership.go:68/74/82/93→157-197), so no supported pre-existing continuation path
  regressed; only the failure point moved earlier with a clearer message.
- **Fix 2 (absent snapshot vs absent manifest comment).** quota_signoff.go:52-57 now
  distinguishes exactly what the plan demanded: absent snapshot in a valid manifest
  stays absent and cannot earn the exception; missing/corrupt/foreign manifest still
  blocks, naming `validateManifests` as the also-requiring backstop.
- **Fix 3 (real standalone tests at the seam).**
  `TestReviewGateSignoffSnapshotAndManifestIntegrity` (app/review_gate_timing_test.go:87-172)
  runs the real `requestConsensusSignoffs` entry four ways: absent-snapshot → signoffs
  succeed (TriageReady), the "no roster snapshot" warning is emitted
  (consensus_request_signoffs.go:136-138), exactly two run manifests exist and neither
  invented a snapshot; missing-manifest → `os.IsNotExist`; corrupt-manifest → "unexpected
  end of JSON input"; foreign-manifest → "contradictory quota manifest" (i.e. the
  `validateManifests` backstop, not the new load — the strongest of the four, it proves
  Fix 1's equivalence claim behaviorally). All three negative variants additionally
  assert zero dropout telemetry records, so no child launched before the integrity gate.
- **Fix 4 (empty snapshot cannot qualify even after valid dropout).**
  `TestReviewGateAbsentSnapshotCannotQualify` (membership/review_gate_test.go:176-193)
  settles a real 4→2 participant-failure batch, then nils the run manifest's
  RosterSnapshot and requires `SingleReviewerAfterDropout` to fail with an error
  (review_gate.go:121-137: load succeeds, models map is empty → "requires known distinct
  snapshot models"). This is the FINAL D2 "Missing snapshot does not qualify" rule
  pinned end-to-end through a committed cause.
- **Fix 5 (legacy-rule negatives at the shared predicate).**
  `TestReviewGateLegacyEvidenceNegatives` (membership/review_gate_test.go:195-228) runs
  all three enumerated legacy rules × {valid-boundary, wrong-adapter,
  wrong-provenance, reset-contradiction}: named-reset at exactly +60 min passes and +59:59
  fails (`quota.MinimumReset`, review_gate.go:176); the no-reset rules fail when any
  reset is present (:177-178); non-`zcode` adapter and non-frozen provenance fail
  (:170-171). The wrong-provenance row exceeds the signed plan (four rows were agreed) —
  accepted, as recorded in the consensus. `TestReviewGateKickoffAndLegacyRules` retains
  the positives plus the arbitrary-`test` negative.
- **Fix 6 (Markdown blank line).** docs/agent-cli-mechanics.md:94 — the blank line
  before `## Participant-step timing and buffered output (1.53.0)` is present.
- **Covering test retained.** `TestReviewGateSignoffWatchdogAndSameTickRebind` is
  unchanged and still passes (my rerun below), resolving round-1 Open question 3 as
  planned.

### AC1 — shared derivation, kickoff/mid/prospective

Unchanged since round-1 review (cycle 1 touched no membership/driver/runner predicate
code; `git diff --stat a26f588..b89e2ab -- internal/membership internal/driver
internal/runner internal/agents` shows test files only). Re-verified by reading:
`Settle` feeds the in-memory settled decision to `checkGates(..., &d)`
(membership.go:314) whose prospective branch requires only that the decision extend
current history (review_gate.go:61-65) — no first-transition deadlock; later gates
(driver_impl.go:627, 692, 740-741, 804-806, 962-964, 1047-1049) all call the same
`SingleReviewerAfterDropout`; the driver consumes only the derived
`ReviewStatus.SingleReviewerAfterDropout` (driver/impl.go:59-66, 288-296), never
frontmatter. Reran `TestReviewGateProspectiveCommitAndPolicyRevision` and
`TestReviewGateKickoffAndLegacyRules` myself (below). Could not refute.

### AC2 — negatives

Re-ran the full adversarial table `TestReviewGateProspectiveAdversarialEvidence`
(arbitrary-rule, wrong-agent, missing-attempt, duplicate-candidate, added-id,
wrong-before, not-applied, protected, unknown-snapshot, pending, corrupt) plus
`TestReviewGateSnapshotModelsAndNoEvidence` (marker-only/manual exclusion, two-person-
by-design, unrecorded single reviewer still blocked by `CheckGates`,
review_gate_test.go:82-94). Code paths re-read: manual/owner transitions excluded
(review_gate.go:31-33, 77-79), latest-wins walk with owner-edit invalidation (:70-82),
Before/After/removal exact-match (:99-115), pending/manual history errors (:31-33).
Blind spot unchanged and disclosed: I did not hand-forge a history file; tampering is
outside this rule's threat model (stops-for-repair is the shipped stance). Could not
refute.

### AC3 — typed evidence allowlist

Re-read `reviewFailureEvidence` (review_gate.go:166-181): `participant-failure.v1` via
the paired-attempt validator with kickoff readiness-step binding (:168), else the three
frozen legacy RuleIDs with `zcode` adapter + frozen provenance + `Failure == nil` + the
per-rule reset constraints. No prose/elapsed-time path exists. The new negatives
(Fix 5) close round-1 NIT-4's coverage gap at exactly this seam. Could not refute.

### AC4 — known distinct snapshot models, diversity=false included

Fixture sets `require_model_diversity: false` and still rejects `""`, `unknown`,
`cli-default`, trim/case-colliding models while accepting `" model-B "`
(review_gate_test.go:70-81; review_gate.go:128-138, `knownModel` :158-164; constants
re-verified `agents.Unknown`/`agents.CLIDefault`). Snapshot identity binding (:125-127)
and duplicate-identity error (:130-132) re-read — see NIT-1 on the missing duplicate
fixture. NoModelBinding stays a disclosed configured-authority caveat in docs and skill
text; nothing claims an observed model. Could not refute.

### AC5 — count-only relaxation

Re-ran `TestDropoutSingleReviewerChangesOnlyCountGate`: `qualified` completes only after
a goal-check call; no-proof, no-reviewer, reservations, goal-fail, blocked and the
strict-gate NIT scan never complete with the exception flag true. Unqualified
single-reviewer close keeps `minimum` 2 (driver/impl.go:290-296). Signer duties,
protected/floor checks and `Settle`'s durable Block path unchanged. Could not refute.

### AC6 — same-reviewer goal check, 120 s

GoalCheck keeps `Timeout: 2 * time.Minute` (driver_impl.go:838, re-verified this
launch), keeps the local implementer refusal, and now also refuses when the derivation
errors (driver_impl.go:804-806). Fresh-process semantics (new `CommandFor` launch) and
the fail-closed escalations unchanged. Could not refute.

### AC7 — scoped supervision defaults

Re-ran `TestReviewGateScopedSupervisionDefaults`: participant step 120/300/60, ordinary
120/1800/60, explicit overrides, all-negative disables, buffered (soft guards off,
heartbeat on), short-hard clamp — all as FINAL D3 (supervision.go:41-46 +
`supervisionForAgent` precedence). Non-participant paths (runner.go:691, acp.go:161,
consult.go:126) use `supervisionForStep` with non-step contexts — behaviorally
identical to baseline there (the 1800 s case is pinned by the test). Could not refute.

### AC8 — terminal classification, two attempts, valid output wins

Re-ran `TestReviewGateSupervisedCommandTerminalAndReplay` with real subprocess fixtures:
no_first_output/stalled/timeout classes recorded with exactly 2 attempts, replay mints
no third attempt, buffered silent-then-success survives (1 call, no failure class),
disabled guards fall to hard timeout. `RunSupervised` (launch.go:157-207) finalizes
telemetry after classification and cleanup; valid-own-output precedence and
cancellation semantics come from the unchanged `RunParticipantStep`/TriageBlocked paths.
Could not refute.

### AC9 — batch settlement, rebind, retained dissent

Re-ran `TestReviewGateSignoffWatchdogAndSameTickRebind` end-to-end: stalled signer
settles with a 2-attempt ledger after the others accepted, remaining signers reach
TriageReady, the already-constructed adapter rebinds reviewers to `[beta]` in the same
tick, and the new standalone run inherits the snapshot and still derives the exception.
The new integrity test (Fix 3) extends this seam to the absent-snapshot and
invalid-manifest variants with no-child assertions. Retained-BLOCK refusal
(TriageBlocked) and below-floor/protected durable Block via `Settle` unchanged. Could
not refute.

### AC10 — buffering audit

`BuffersStdout: true` on zcode and default claude text only (discover.go:249, 433),
pinned by `TestReviewGateBufferedDefaultContracts` (codex/kimi false; zcode/claude/agy
true — rerun). Audit evidence (healthy 775.745 s zcode invocation, first activity at
775.662 s) recorded in both docs; no Claude task process appears in any validation log
(kimi-readiness and native-launch records only). Docs disclose buffered/manual/
interactive hard-only limits and the override responsibility. Could not refute.

### AC11 — 14 hunks, three copies, staged core, caps/map

Independent byte-exact reconstruction, my own script, this launch: base live
COOPERATION.md at 85c5a8f (sha 091e6fb8…, 125,879 B, equals `base_source_sha256`;
protocol-hunks.json sha 54d0dd07… equals FINAL D5) + the 14 hunks == committed live copy
(acbd4dbc…, 125,132 B == the packet); same for the embedded copy (d8cc6dbd…, 124,886 B)
and the staged core: ~/.parley/staging/COOPERATION-2.16.0.md (8c6b95b6…) + the same
hunks == the retained core-2.17.0-preview.md (6073c311…, 124,976 B). Every `old` occurs
exactly once pre-replacement; every `new` occurs exactly once post-replacement in both
protocol copies. The skill's references/COOPERATION.md is byte-identical to the live
copy (verified). Packet guardrail + drift tests re-run: exit 0. Applicability map
absent from both diffs. Skill SKILL.md 18,990/20,000 B with frontmatter, Core Rule and
required sections intact (verified directly, see AC12 for the suite). Could not refute.

### AC12 — validation, independently recounted

- **Full current-source Go suite at b89e2ab** (`full-host-tests-cycle1.jsonl` +
  `full-host-tests-cycle1-result.json`, head recorded b89e2ab, cmd
  `go test ./... -json -count=1 -timeout 45m`, started 08:37:49Z, finished while this
  review waited on it): receipt exit_code 0, and I re-hashed the complete 3.4 MB jsonl
  myself — sha256 04d2f0529c7e38046ac64c7be04603fbdd05dd0cc461b796d5d8c65f27724db5,
  matching the receipt's `log_sha256`. My own event recount: 14,858 events, 34/34
  packages with terminal `pass`, 0 `fail` events, 0 packages without a terminal result,
  3 no-test packages, 3,330 distinct passing test/subtest IDs and the same 4 built-in
  skips as v2 (budget×2, runner×2). The +19 IDs over v2's 3,311 are exactly the new
  cycle-1 tests (1 absent-snapshot + 13 legacy-negative parent/subtests + 5 integrity
  parent/subtests) — arithmetically consistent with the diff. The receipts are
  byte-identical across the native task temp
  (`/private/var/tmp/pd-review-gate-tests-bazott2q/validation-receipts/`) and the shared
  volume (sha256-compared by me) — no ENOENT anomaly this time. Pending was not counted
  as pass: I watched the run live (PID 40918, 10:37–11:01 local) and recounted only
  after exit.
- **build-cycle1 / vet-cycle1**: exit 0 at b89e2ab (receipts read directly; empty logs,
  sha e3b0c442… as expected for silent success).
- **Skill**: worktree still exactly 74cc831, clean; cycle 1 touched no skill file, so
  the v3 pass stands — `skill-tests-v3.log` sha256 3509b5ce… matches my `shasum` and
  the committed implementation-validation.json (399 Node, 54 Python, six manifests, all
  `ok`, aggregate d380b77a… which I additionally verified equals the recomputed
  parley-addon.json aggregate and that all 10 manifest file hashes match the actual
  bytes on disk).
- **Focused fix-up runs**: fixup1-focused-v1/v2 exit 1 retained as failure history
  (fixture setup defects in the new test — verified in the logs); v3 exit 0, log sha
  8548172a… matches the receipt and IMPLEMENTATION.md.
- **My own independent runs** (native scratch TMPDIR/GOTMPDIR under
  `.parley-runtime/review-gate-timing/zcode-review-02`, shared GOCACHE, `GOFLAGS=-p=2`,
  concurrent-safe with the then-running full suite):
  `go test ./internal/{app,membership,driver,runner,agents} -run
  'TestReviewGate|TestDropoutSingleReviewerChangesOnlyCountGate' -count=1 -timeout 10m`
  → all ok; `go test ./internal/app -run
  'TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail|TestProtocolDrift'` → ok.
- Deferred receipt follow-up re-evaluated: this cycle's receipts were readable on both
  storages at first attempt with matching hashes, so the recorded watch-condition ("open
  the follow-up if the volume repeats this") is still unmet and the deferral remains
  sound. My v2 recount (which independently established the earlier pass from the raw
  log) is now corroborated by a clean recurrence-free cycle.
- See NIT-2: the cycle-1 receipt identities are not yet appended to the *committed*
  validation record.

Could not refute the pass claims; the two NITs below are bookkeeping, not failures.

### AC13 — attended close gates (process)

Still a process gate by design. This review is the required final independent full-scope
review and it carries no open CRITICAL/MAJOR/MINOR. Still outstanding and correctly
unclaimed: the Phase 7 consensus with zero Agreed fixes, both current participants' final
review-consensus ACCEPT blocks, the fresh zcode-1 goal-check process, and the
current-tree AC evidence record. The run correctly claims no product single-reviewer
exception for itself (owner-marker kimi exclusion recorded at Phase 0; no mid-idea
automatic transition exists to derive from — re-confirmed: the idea's history has the
kickoff only). Nothing to refute; no completion is asserted by codex-1
(IMPLEMENTATION.md cycle-1 section explicitly leaves this review, the zero-fix
consensus, ACCEPTs and goal process open). Fix-up cycles used: 1 of 5.

### AC14 — post-close delivery

Expressly post-close; no delivery claimed, none verified beyond PR existence (CLI #78 /
skill #10 open at round-1, gh-verified then; not re-verified this launch since no
delivery claim depends on it now). Nothing to refute.

## Findings

### [NIT] Duplicate roster-snapshot identity error branch has no failing fixture
`internal/membership/review_gate.go:130-132` returns "single-reviewer roster snapshot has
duplicate identity" when two snapshot entries share an agent ID — verified by reading
only. The model table (TestReviewGateSnapshotModelsAndNoEvidence) and the new cycle-1
tests cover empty/unknown/cli-default/collision/absent variants but not a duplicated
`Agent` entry. One table row would pin the error branch the way the legacy negatives now
pin theirs. Same class as round-1 NIT-4, which cycle 1 resolved for the evidence seam.

### [NIT] Cycle-1 validation receipts not yet in the committed validation record
`source-context/implementation-validation.json` (committed at 5a3e653) still records the
a26f588 v2 run as the latest full-suite identity; the b89e2ab cycle-1 receipts
(full-host-tests-cycle1/build-cycle1/vet-cycle1, exit 0, log sha 04d2f052…/e3b0c442…/
e3b0c442…) currently exist only under `.parley-runtime/review-gate-timing/` (private but
durable and independently verified by this review, including byte-identity across both
storages). IMPLEMENTATION.md's committed cycle-1 section, written before the run
finished, announces the run without claiming its result — honest, and the private
retention satisfies the letter of AC12's "retain all failed and passing evidence". For
the attended close, the exit-0 identities should be appended to the committed record (or
a cycle-1 validation note) so the close's current-tree evidence does not rest on
reviewer testimony alone. This self-resolves with the implementer's next artifact
commit; nothing about it is a product defect.

## Open questions

1. For the fresh goal-check process (AC13): both NITs above are non-blocking for this
   reviewer, but the goal process should confirm the cycle-1 receipt identities are
   appended to the committed record before the attended close is claimed (NIT-2), and
   may optionally fold the one-row duplicate-snapshot fixture in with any future touch
   of the membership tests (NIT-1).
2. Carried from round-1 Open question 2, unchanged: the deferred timeout-tuning
   follow-up idea should record the fresh goal-check duration this run produces (per
   FINAL "Known risks"), alongside the disclosed fact that a genuinely hung buffered
   zcode signoff is detected only at its hard ceiling. The 775.7 s healthy buffered
   invocation remains the standing evidence for that disclosure.
3. None blocking. My round-1 MINOR-1 stands withdrawn after this launch's re-challenge;
   its actionable residue is fully discharged by the signed and implemented cycle-1
   fixes, and no new issue at the seam survived probing (Fix 1/3 above).
