---
agent: zcode-1
idea: meta-protocol-change-participant-dropout-review-gate-timing
kind: goal-check
date: 2026-10-09
checked-commit: 8a939761175fa70226cc36ed09c4cd5be33f8838
product-commit: b89e2abaa056813a4b238c2fa4278195b0f7096b
skill-commit: 74cc831b18ce33e48e705fc8992c89be85cadf6f
context_mode: full
source_sha256: acbd4dbc0c0702bc191176bb80bcee32c5093c8b4ebbee42e036df6a9b7d1137
packet_sha256: acbd4dbc0c0702bc191176bb80bcee32c5093c8b4ebbee42e036df6a9b7d1137
fallback_reason: none
---

# Fresh independent goal check — zcode-1

## Launch attestation

context_mode=full, no fallback. I re-hashed the packet file myself
(`shasum -a 256` → `acbd4dbc0c0702bc191176bb80bcee32c5093c8b4ebbee42e036df6a9b7d1137`,
matches `packet_sha256`) and read all 1,552 lines of
`.parley-runtime/protocol-packets/full-phase8-deliberation-acbd4dbc….md`. The packet hash
equals the live `parley-deck/COOPERATION.md` hash, verified below. This native process
(native start `goal.native-start.json`, PID 8391, started 2026-10-09T09:10:07Z, timeout
1800 s, model zai/glm-5.3) is the attended goal assessment authorized by the controlling
brief — it is NOT a claim that the product's 120 s auto-goal parser executed it. The
product goal bound stays 120 s. Elapsed at artifact authorship ≈ 8 min (>120 s); final
duration is recorded by the process's own exit receipt for the disclosed timeout-tuning
follow-up. I authored only this file and my own signoff; no product edits, no commits, no
GitHub messages, no roster/credential/model changes.

## Tree state (verified this launch)

- HEAD `8a93976` = stated `9311bfc` + one commit adding only
  `source-context/goal-brief.md` (14 lines). `git diff --stat b89e2ab..HEAD -- . ':(exclude)parley-deck'`
  is empty and `git status --porcelain` over product paths is clean: **no product change
  after b89e2ab**. Only untracked state: two `parley-deck/runs/<ts>/` driver bookkeeping
  dirs (`events.jsonl`/`run.json`) — run state, not product.
- Skill worktree `review-gate-timing-skill`: HEAD exactly `74cc831`, clean.
  `SKILL.md` 18,990 / 20,000 B, sha256 `0b9769e8…` — matches
  `implementation-compaction-proof.json`.
- Kimi exclusion is recorded in `00-prompt.md` as owner-confirmed manual
  (`excluded: [kimi-1 — … confirmed 2026-10-09]`); the idea's history has no qualifying
  automatic ≥2→1 transition, and nothing in these artifacts claims the product
  auto-close exception for this run.

## Independent commands and results (all this launch, current tree = product b89e2ab)

Test env per `test-storage.json`: `TMPDIR`/`GOTMPDIR=/private/var/tmp/pd-review-gate-tests-bazott2q`,
shared `GOCACHE`, `GOFLAGS=-p=2`.

1. `go test ./internal/{app,membership,driver,runner,agents} -run 'TestReviewGate|TestDropoutSingleReviewerChangesOnlyCountGate' -count=1 -timeout 15m`
   → **exit 0** (2026-10-09T09:13:18Z–09:13:38Z): app ok 16.772s, membership ok 4.266s,
   driver ok 0.499s, runner ok 8.425s, agents ok 0.199s.
2. `go test ./internal/app -run 'TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail' -count=1 -v`
   → **PASS**; the test's own live R-2 measurement prints facilitator body **68,846 B**
   against the 70,000 B guardrail (matches the recorded preview figure).
3. `go test ./internal/protocol/... ./internal/protocolpacket/... -count=1 -timeout 10m`
   → **exit 0** (protocol 0.237s incl. `TestEmbeddedDefaultMatchesLiveDeck` and
   `TestDriftAnchorsAcceptTheGeneratedRosterTable`; protocolpacket 1.595s).
4. `go build ./...` → **exit 0**; `go vet ./...` → **exit 0**.
5. Full-suite log re-hash and recount (below) and my own AC11 byte-reconstruction (below).

## AC coverage (current tree)

**AC1 — shared derivation, kickoff/mid/prospective: MET.** One function,
`membership.SingleReviewerAfterDropout` (`internal/membership/review_gate.go:26-140`),
feeds every consumer: `Settle` passes the in-memory settled decision to `checkGates` pre-commit
(`membership.go:314`), whose prospective branch only requires extending current history
(`review_gate.go:61-65`, no first-transition deadlock); the driver wrapper
(`driver_impl.go:773-778`) is the single derivation call and is consumed at :627, :692,
:740, :807 (goal refusal), :962 and :1047. Kickoff transitions qualify on the recorded
proposed batch with readiness-step binding (`review_gate.go:83-88, 168`). Focused tests
`TestReviewGateProspectiveCommitAndPolicyRevision` and `TestReviewGateKickoffAndLegacyRules`
ran green in my run (item 1).

**AC2 — negatives: MET.** Source re-read: pending/manual history errors (`:31-33`), no
history → no exception (`:34-37`), membership must equal immutable history (`:67-69`),
latest membership-changing batch with owner-edit invalidation (`:70-82`), exact
Before/After and removal-set matching with protected/facilitator refusal (`:90-115`),
candidates must match removed IDs exactly (`:113-115`), arbitrary RuleIDs fail in
`reviewFailureEvidence` (`:166-181`). `TestReviewGateProspectiveAdversarialEvidence`
(11 negative variants) and `TestReviewGateSnapshotModelsAndNoEvidence` (marker-only,
manual, two-person-by-design, unrecorded single reviewer still blocked) ran green.

**AC3 — typed evidence allowlist: MET.** `reviewFailureEvidence` (`:166-181`):
`participant-failure.v1` only via the paired-attempt validator with kickoff readiness
binding; else exactly the three frozen legacy RuleIDs requiring `zcode` adapter, frozen
provenance `zcode.stderr.owner-deviation.2026-10-04`, `Failure == nil`, and the per-rule
reset constraints (named-reset ≥60 min with raw reset; both no-reset variants require no
reset evidence). No prose/elapsed-time path exists. `TestReviewGateLegacyEvidenceNegatives`
(all three rules × valid-boundary/wrong-adapter/wrong-provenance/reset-contradiction,
incl. +60 min pass vs +59:59 fail) and the arbitrary-`test` negative ran green.

**AC4 — known distinct snapshot models, diversity=false included: MET.**
`knownModel` (`:158-164`) maps empty/`unknown`/`cli-default` to unknown; trimmed
case-insensitive comparison with distinct-model requirement (`:128-138`); run/idea
binding (`:125-127`); duplicate identity fails closed (`:130-132`, error). The model
fixture sets `require_model_diversity: false` and still rejects the unknown/colliding
variants — ran green. NoModelBinding stays a disclosed configured-authority caveat in
docs/skill; nothing claims observed native models.

**AC5 — count-only relaxation: MET.** `driver/impl.go:291-296`: `minimum` stays
`MinReviewers` (2) unless `SingleReviewerAfterDropout && ReviewerCount == 1`; signer
duties, reservations, retained BLOCK, strict-gate and protected/floor checks are outside
that branch. `TestDropoutSingleReviewerChangesOnlyCountGate` ran green (driver package),
including the unqualified-still-blocked and reservation/goal-fail/negative scans.

**AC6 — same-reviewer goal check, 120 s: MET.** `GoalCheck` keeps
`Timeout: 2 * time.Minute` (`driver_impl.go:838`), refuses a missing/self checker
(`:818-820`) and now also refuses when the shared derivation errors (`:806-808`); fresh
`CommandFor` launch semantics and fail-closed escalations unchanged.

**AC7 — scoped supervision defaults: MET.** `supervision.go:31-36`: 120 s first output,
1,800 s default stall, **300 s participant-step stall**, 60 s heartbeat;
`supervisionForStep` tightens only `ParticipantStepActive` contexts when no explicit
override exists, honoring disables and hard-ceiling clamping; heartbeats never count as
activity; buffered declarations disable soft guards. `TestReviewGateScopedSupervisionDefaults`
ran green (runner package), including the 1800 s non-participant pin, overrides,
disables, buffered and clamp cases.

**AC8 — terminal classification, two attempts, valid output wins: MET.**
`TestReviewGateSupervisedCommandTerminalAndReplay` ran green with real subprocess
fixtures: no_first_output/stalled/timeout classes with exactly 2 attempts, replay mints
no third attempt, buffered silent-then-success survives, disabled guards fall to hard
timeout. `RunSupervised` finalizes telemetry after classification and cleanup; valid
own output (BLOCK included) wins over failure candidacy.

**AC9 — batch settlement, rebind, retained dissent: MET.**
`TestReviewGateSignoffWatchdogAndSameTickRebind` ran green end-to-end (stalled signer
settles after others accepted; same-tick rebind to `[beta]`; standalone run inherits the
snapshot and still derives the exception), and
`TestReviewGateSignoffSnapshotAndManifestIntegrity` ran green covering
absent-snapshot-success-with-warning, missing (`os.IsNotExist`), corrupt and foreign
(`contradictory quota manifest`) manifests with zero-child assertions. Retained-BLOCK and
below-floor/protected durable Block paths unchanged (`Settle`).

**AC10 — buffering audit: MET.** `discover.go`: `BuffersStdout: true` exactly on default
Claude text (`:249`), zcode (`:433`) and pre-existing agy (`:274`); codex/kimi remain
streaming. `TestReviewGateBufferedDefaultContracts` ran green (agents package). Docs
disclose buffered/manual/interactive hard-only limits and the override responsibility;
no Claude task process appears in any validation record (kimi-readiness/native-launch
records only — checked the runtime evidence listing).

**AC11 — 14 hunks, three copies, staged core, caps/map: MET (independently
reconstructed).** My own script this launch, from
`source-context/protocol-hunks.json` (sha256 re-verified
`54d0dd07d747cf54ce9d0e234b9d65c34a8282b525bd8667090aaf19cbffcc39`, `base_source_sha256`
equals the `85c5a8f` live base hash `091e6fb8…`): applying the 14 substitutions — each
`old` occurring exactly once pre-replacement — reproduces byte-identically the live
`parley-deck/COOPERATION.md` (sha256 `acbd4dbc…` = the packet), the embedded
`internal/protocol/defaults/COOPERATION.md` (`d8cc6dbd…`), and the staged core preview
(`~/.parley/staging/COOPERATION-2.16.0.md` `8c6b95b6…` + hunks ==
`core-2.17.0-preview.md` `6073c311…`, 124,976 B). Every `new` occurs exactly once
post-replacement in both protocol copies. The skill `references/COOPERATION.md` is
byte-identical to the live copy. Packet guard test PASS with live measurement 68,846 /
70,000 B; SKILL.md 18,990 / 20,000 B with frontmatter/Core rule/required sections intact;
drift tests (`internal/protocol`, `internal/protocolpacket`) exit 0; no applicability-map
change in either product diff. Staging/publication itself is AC14 post-close.

**AC12 — validation, re-hashed and recounted: MET.**
- `full-host-tests-cycle1.jsonl` re-hashed:
  `04d2f0529c7e38046ac64c7be04603fbdd05dd0cc461b796d5d8c65f27724db5` (3,416,558 B) —
  matches the durable receipt and `fixup1-validation.json`. My own recount: **14,858
  events; 34/34 packages terminal pass; 0 `fail` events; 3,330 distinct passing
  test/subtest IDs; 4 test-level skips** (budget×2, runner×2) plus 3 no-test packages;
  37 packages seen, none without a terminal result. Identical to the recorded counts.
- Receipts `full-host-tests-cycle1/build-cycle1/vet-cycle1-result.json` (all `exit_code: 0`,
  head `b89e2ab`) are **byte-identical (cmp)** between the shared volume and the native
  task temp `/private/var/tmp/pd-review-gate-tests-bazott2q/validation-receipts/`; build/vet
  logs are empty (sha `e3b0c442…`, silent success). The committed
  `source-context/fixup1-validation.json` matches all three receipts field-for-field.
- Skill unchanged at `74cc831`; `skill-tests-v3.log` re-hashed
  `3509b5ce4c7db449fb6f72ab9055215ce6ad7d0c29bbe367ff83a0aa0cc35a4b` (399 Node / 54
  Python / six manifests, recorded exit 0).
- My own targeted runs (items 1–4 above) all exit 0 on the current tree. Failed-history
  logs (v1/v2, fixup v1/v2) remain retained as failures; nothing failed was counted.
- No pending or failed run is counted as a pass anywhere above.

**AC13 — attended close gates: MET by the required combination.** Final independent
review `review/round-02/zcode-1.md` (separate process PID 40945, exit 0) has no open
CRITICAL/MAJOR/MINOR; this artifact is the fresh zcode-1 goal process's current-tree
AC1–AC13 evidence and verdict; both current participants' final review-consensus ACCEPTs
exist once my own signoff is appended (codex-1's is already present); fix-up cycles used
1 of 5. This run claims no product single-reviewer exception for itself (manual Kimi
exclusion, kickoff-only history — re-confirmed in `00-prompt.md` and the packet's §9.0
wording).

**AC14 — post-close delivery: NOT CLAIMED, NOT COUNTED.** Explicitly post-close per
FINAL; nothing here verifies delivery, and none is needed for this verdict.

## Findings and disposition assessments (nothing suppressed)

1. **NIT-1 (duplicate-snapshot error branch has no failing fixture) — deferral
   concurring.** The branch (`review_gate.go:130-132`) is fail-closed in source (returns
   an error; the exception cannot be earned), re-read this launch. The adjacent
   identity/model table covers empty/unknown/cli-default/collision variants, and the
   branch is two lines above the empty-map failure that IS fixture-pinned
   (`TestReviewGateAbsentSnapshotCannotQualify`). A one-row fixture is a coverage
   improvement, not a correctness gap; deferring to the next membership-test touch
   matches the round-02 reviewer's own non-blocking suggestion. Not blocking.
2. **NIT-2 (cycle-1 receipts not yet in the committed record) — resolved, verified.**
   `source-context/fixup1-validation.json` (added at `ce7ec30`) exists now and matches
   the actual receipts exactly (see AC12); receipts are byte-identical across both
   storages. The timing-overlap bookkeeping concern is discharged. Not blocking.
3. **Observation (new, non-blocking, evidence-description only):** round-02's note
   "`go test ./internal/app -run 'TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail|TestProtocolDrift'` → ok"
   cites a test name that matches no test anywhere (`grep -rn "TestProtocolDrift" internal/`
   is empty; the drift tests live under other names in `internal/protocol` /
   `internal/protocolpacket`). The package still exits 0 because Go does not fail on
   zero matched tests, so the recorded exit is honest but that alternation's second half
   was vacuous. No impact on any AC outcome: the recounted full suite at b89e2ab covers
   those packages, and I independently ran both packages this launch — exit 0 (item 3).
4. Prior MINOR withdrawal (round-01, prior-manifest dependency) was re-challenged in
   round-02 with a failing foreign-manifest fixture proving the `validateManifests`
   backstop catches what the new earlier load lets past; the withdrawal survives my
   re-read of `quota_signoff.go:45-75` and `membership.go` `validateManifests`. The
   receipt-ENOENT watch-condition did not recur (dual receipts byte-identical this
   cycle) and its deferral condition remains unmet. Receipt ENOENT follow-up stays
   disclosed/TBD if recurring — agreed.

## Goal verdict

GOAL-CHECK: PASS

Current-tree evidence independently verifies AC1–AC13 against product `b89e2ab` /
skill `74cc831`; AC14 is post-close and not claimed. This textual PASS is defense in
depth on top of the current-tree criterion evidence above; it substitutes for nothing
and certifies no exemption. This run's attended close authority remains the
controlling brief's combination: no open CRITICAL/MAJOR, both final ACCEPTs, this
fresh goal PASS, and recorded current-tree evidence.
