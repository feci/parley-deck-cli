---
agent: zcode-1
idea: meta-protocol-change-participant-dropout-review-gate-timing
review-round: 1
date: 2026-10-09
reviewed-commit: a26f588
scope: full
---

## Summary

Full-scope refutation review of the CLI product diff 85c5a8f..a26f588 (HEAD of the reviewed
implementation; PR feci/parley-deck-cli#78 head 5a3e653 adds only non-product artifacts after
a26f588 — verified: `git diff --stat a26f588..5a3e653 -- internal/ docs/ README.md CHANGELOG.md
VERSION` is empty) and the skill diff dbdb919..74cc831 (PR feci/parley-deck-skill#10 head
74cc831, gh-verified). I read the controlling brief, the complete live protocol packet, frozen
FINAL.md/IMPLEMENTATION.md, every product/test/protocol file in both diffs, the three
source-context transformation artifacts, and the validation logs, and I ran my own probes.

**Launch attestation (recorded per brief):** this phase-6 launch used
`.parley-runtime/protocol-packets/full-phase6-deliberation-acbd4dbc….md` with
`context_mode=full`, `source_sha256=packet_sha256=acbd4dbc0c0702bc191176bb80bcee32c5093c8b4ebbee42e036df6a9b7d1137`,
`fallback_reason` absent (`review-attestation.json`, PRIMARY). I independently re-hashed the
packet file with `shasum -a 256` and got the same digest. The packet was read in full (1,553
lines) before reviewing. Review dispatch: the driver's `continue` selected consensus-signoff
with no recoverable review-dispatch action (documented D6 pre-dispatch refusal in the original
driving run), so this review runs as the configured-native zcode-1 headless fallback
(PID 95819, launched 09:53 local, `.parley-runtime/review-gate-timing/review01.native-start.json`).

**Verdict: no CRITICAL, no MAJOR findings. One MINOR, three NIT.** The causal predicate in
`internal/membership/review_gate.go` is tight against every attack I attempted (listed per AC
below); the gate is re-derived at every consumer from immutable history plus the roster
snapshot, never from writable metadata; only the numeric minimum becomes one; the supervision
changes are scoped, honor overrides/disables/buffering, and reuse the shipped two-attempt
ledger; all 14 protocol substitutions reproduce all three COOPERATION.md copies and the staged
core 2.17.0 preview byte-exactly (I reconstructed them independently); the full Go suite at
a26f588 exited 0 (34/34 packages, 0 fail events — receipt recovered and hash-verified by me,
see NIT-1) and the full skill suite at 74cc831 exited 0.

## Refutation attempts

Provenance: all attempts below are **PRIMARY** — I read the cited code in the a26f588 tree and
ran the cited commands myself unless stated otherwise. Paths are relative to the CLI worktree
`/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/review-gate-timing`; the skill
worktree is the sibling `review-gate-timing-skill`.

- **AC1 (shared derivation, kickoff/mid/prospective).** Attempted to break the one-predicate
  claim: (a) precommit deadlock — `Settle` passes the in-memory decision to
  `checkGates(..., &d)` (membership.go:314) whose prospective branch requires only
  `prospective.Before == h.Current` (review_gate.go:61-65), so the first eligible loss is not
  deadlocked by a not-yet-committed history; (b) later gates — `driverImplOps.singleReviewerAfterDropout`
  (app/driver_impl.go) feeds OpenReviewRound, DraftReviewConsensus, ReviewStatus, GoalCheck,
  RequestReviewSignoffs and Complete, all calling the same `membership.SingleReviewerAfterDropout`;
  the driver consumes only the derived `ReviewStatus.SingleReviewerAfterDropout` (driver/impl.go:59-66,
  288-296), never frontmatter. Ran `TestReviewGateProspectiveCommitAndPolicyRevision` and
  `TestReviewGateKickoffAndLegacyRules` (my run, below). Could not refute: one derivation, both
  entry paths, all consumers.
- **AC2 (negatives).** Attempted: marker-only/manual (`v.Manual != nil` → error; owner batch →
  `b.Owner != nil` → no exception, review_gate.go:31-33, 77-79); unrelated/stale cause (latest
  membership-changing transition only, :72-82; `decision.After` must equal current ids, :91);
  later manual edit invalidates cause even when it restores an older set (latest-wins walk +
  test's `new` catch-up case); Before/After/removal mismatches (`removed` must exactly equal
  candidate IDs, :99-115); duplicates (`sameMemberSet` enforces uniqueness via `quota.Unique`);
  pending/corrupt history (error, fail closed); two-person-by-design (no transition → no
  decision → false; `TestReviewGateSnapshotModelsAndNoEvidence` also pins `CheckGates` still
  blocking an unrecorded auto single reviewer). Could not refute. Blind spot: I did not attempt
  to hand-forge a history file; tampering is outside this rule's threat model (stops-for-repair
  is the shipped stance).
- **AC3 (typed evidence allowlist).** Attempted: arbitrary RuleID (`reviewFailureEvidence`
  accepts only `participant-failure.v1` via the paired-attempt validator
  (`quota.ValidateFailureEvidence`, quota/dropout.go:41-62: provenance `supervisor-terminal`,
  exactly 2 linked attempts, retry linkage, no reset fields) plus the three frozen legacy rules
  with `zcode.stderr.owner-deviation.2026-10-04` provenance, adapter `zcode`, `Failure == nil`,
  reset ≥ 60 min for `named-reset` (`quota.MinimumReset = time.Hour`), and no reset for the two
  no-reset rules (review_gate.go:166-181); kickoff evidence must carry `Step == "readiness"`
  (:168). Prose/elapsed-time inference has no code path. Could not refute; see NIT-4 on
  negative-fixture coverage at this seam.
- **AC4 (distinct known snapshot models even with diversity=false).** The fixture prompt sets
  `require_model_diversity: false` and the test table still rejects `""`, `unknown`,
  `cli-default`, `" model-A "` (trim+casefold collides with implementer) and accepts
  `" model-B "` (review_gate.go:128-138, `knownModel` :158-164; constants verified
  `agents.Unknown="unknown"`, `agents.CLIDefault="cli-default"`). Snapshot binding:
  `m.RunID != runID || m.IdeaSlug != base(ideaDir)` → error; duplicate snapshot identity →
  error; missing snapshot → `runmanifest.Load` error → no exception. NoModelBinding stays
  configured authority — nothing in the diff claims observed models (docs and skill text
  disclose this). Could not refute.
- **AC5 (count-only relaxation).** Ran `TestDropoutSingleReviewerChangesOnlyCountGate`
  (driver): `qualified` completes only with a `goal-check` call; `no-proof`, `no-reviewer`,
  `reservations`, `goal-fail`, `blocked`, `strict-nit` never reach complete — reservations
  (LE-11) and the strict-gate deterministic NIT scan still veto with the exception flag set
  true. Signer duties: review-consensus signers follow current participants (consensus.Status,
  unchanged); protected/floor checks live in `Settle`/`quota.Evaluate` (unchanged) and the
  predicate re-checks `UsableSurvivors >= Floor` and protected/facilitator removal
  (review_gate.go:90, 105-112). Unqualified single-reviewer close stays blocked
  (`minimum` stays 2, driver/impl.go:290-296). Could not refute.
- **AC6 (goal check by the same reviewer, 120 s).** GoalCheck still launches a fresh process
  with `Timeout: 2 * time.Minute` (app/driver_impl.go:838, unchanged), keeps the local
  `checker == o.implementer` refusal, and now refuses when the derivation errors
  (app/driver_impl.go:804-806). Fresh-process semantics unchanged (new `CommandFor` launch);
  the failure/inconclusive escalations are pre-existing and untouched. Could not refute.
- **AC7 (scoped defaults).** `supervisionForStep` applies the 300 s stall default only when
  `ParticipantStepActive(ctx) && agent.StallTimeoutMS == 0` (supervision.go:41-46); explicit
  overrides, negative disables, `BuffersStdout` soft-guard disable, and the hard-ceiling stall
  clamp keep their shipped precedence (`supervisionForAgent`, :53-86). Ran
  `TestReviewGateScopedSupervisionDefaults`: step 120/300/60, ordinary 120/1800/60, override,
  disabled, buffered (soft guards off, heartbeat on), short-hard clamp all as expected.
  Non-participant paths (`runner.go`, `acp.go`, `consult.go`) switched from
  `supervisionForAgent` to `supervisionForStep` with non-step contexts — behaviorally
  identical there (verified by reading; the 1800 s default case is pinned by the test). Could
  not refute.
- **AC8 (terminal classification, two attempts, valid output wins).** `RunSupervised`
  (runner/launch.go:157-207) wires counting writers, process-group capture/kill, persists
  heartbeat/watchdog events, and finalizes telemetry after classification and cleanup. Ran
  `TestReviewGateSupervisedCommandTerminalAndReplay` with real subprocess fixtures: no-first-
  output/stalled/timeout terminal classes recorded, exactly 2 attempts, replay mints no third
  attempt, buffered silent-then-success survives (1 call, no failure class), disabled guards
  fall to hard timeout. Valid BLOCK/own-output precedence comes from the unchanged
  `RunParticipantStep` validation plus `requestConsensusSignoffs`'s existing
  TriageBlocked refusal (consensus_request_signoffs.go:106); cancellation returns
  `agentCtx.Err()` (unchanged). Could not refute.
- **AC9 (batch settlement, rebind, retained dissent).** Ran
  `TestReviewGateSignoffWatchdogAndSameTickRebind` end-to-end through the real standalone
  signoff entry: stalled signer settles as `stalled` with a 2-attempt ledger after the other
  two accepted, history revision 1, remaining signers reach TriageReady, the already-constructed
  adapter rebinds in the same tick (`o.quotaCurrent()` → reviewers `[beta]`), the new standalone
  run inherits the snapshot and still derives the exception. Standalone inheritance:
  quota_signoff.go:52-65 copies `RosterSnapshot/RosterRevision` from the history's prior run;
  `requestConsensusSignoffs` binds launches with `applyRosterSnapshot` under `MidIdea()`
  (consensus_request_signoffs.go:131-137). Retained earlier BLOCK still binds via the unchanged
  blocked-consensus refusal. Below-floor/protected loss produces a durable Block through
  `Settle` (unchanged). Could not refute; see MINOR-1 on the new prior-manifest dependency.
- **AC10 (buffering audit).** `BuffersStdout: true` added for zcode and default claude text
  only (agents/discover.go:249, 433), pinned by `TestReviewGateBufferedDefaultContracts`
  (codex/kimi false, zcode/claude/agy true — ran it). The audit evidence (a healthy 775.745 s
  zcode invocation with first activity at 775.662 s) is recorded in docs/agent-cli-mechanics.md
  and docs/agent-runtime-configuration.md; no Claude task process appears anywhere in the
  validation logs (kimi-readiness and native-launch records only). Docs disclose
  buffered/manual/interactive hard-only limits. Could not refute.
- **AC11 (14 hunks, 3 copies, staged core, caps/map).** PRIMARY reconstruction, my own script:
  base live COOPERATION.md at 85c5a8f (`git archive`, sha256 091e6fb8…, 125,879 B, matches
  `base_source_sha256`) + the 14 hunks from protocol-hunks.json (file sha256 54d0dd07…,
  matches FINAL D5) == live copy byte-exact (acbd4dbc…, 125,132 B); same for the embedded
  copy (d8cc6dbd…, 124,886 B) and, with the identical hunk set, the staged core:
  ~/.parley/staging/COOPERATION-2.16.0.md (8c6b95b6…) → 6073c311…, 124,976 B ==
  `.parley-runtime/review-gate-timing/core-2.17.0-preview.md`. Skill references copy is
  byte-identical to live; embedded differs only in the generic workspace/roster zones (the
  drift-guard allowlist). Each `old` occurs exactly once pre-replacement and each `new` exactly
  once post-replacement. SKILL.md reconstructs byte-exactly from skill-compaction.json
  (6 hunks) at 18,990/20,000 B with frontmatter, Core Rule and all 11 `##` sections intact.
  Applicability map not in either diff; packet/drift guard tests passed (my run, below).
  Phase-1 guard 68,846/70,000 B recorded in the proof and enforced by the passing
  `TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail`. Could not refute.
- **AC12 (validation).** Evidence I recounted myself: focused-v1 exit 0 (5 pkgs); static-
  checks-v1 exit 0 (packet-drift/build/vet); validation-fix-v1 exit 0 (the two v1-failed
  tests, matching v1's 5 fail events confined to `internal/app` `TestGoalCheckFailedOrUnverifiable
  ExecutionCannotPass` and `internal/quota` `TestQuotaFixupNoticeNamesRemainingGates`); full-
  host-tests-v1 exit 1 preserved as failure history; **full-host-tests-v2 at a26f588: exit 0,
  34/34 packages, 0 `"Action":"fail"` events, jsonl sha256 04e3fa05… — the on-disk
  *-result.json was unreadable (shared-volume ENOENT) and was recovered into
  `full-host-tests-v2-recovered-result.json` + committed `source-context/implementation-
  validation.json` with a documented recovery note; I re-hashed the jsonl and recounted 3,315
  distinct test/subtest IDs (receipt: 3,311 pass + 4 skip) and the hash matches**; build-v2 and
  vet-v2 exit 0 at a26f588 (result jsons read directly); skill-tests-v1 exit 1 (stale manifest)
  and v2 exit 1 (missing `commonmark` dev dependency — "Error: Cannot find module 'commonmark'",
  not counted as pass) preserved; **skill-tests-v3 after `npm ci` (exit 0): exit 0 at 74cc831**
  (399 node tests incl. every published `node --test` command, 54 python tests, 6 addon
  manifests ok, aggregate sha256 d380b77… matches parley-addon.json). My own independent runs
  (env per test-storage.json: native TMPDIR/GOTMPDIR, shared GOCACHE, GOFLAGS=-p=2):
  `go test ./internal/{membership,driver,app,runner,agents} -run 'TestReviewGate|TestDropout
  SingleReviewerChangesOnlyCountGate' -count=1 -timeout 10m` → all ok, exit 0;
  `go test ./internal/app -run 'TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail|TestProtocolDrift'`
  → ok, exit 0; `go build ./...` → OK; `go vet ./...` → OK. Could not refute; see NIT-1.
- **AC13 (attended close gates).** Process gates still pending by design at review time: this
  review is the required independent final review; the review-consensus ACCEPT blocks, the
  fresh zcode-1 goal process and the recorded current-tree AC evidence follow it. The run
  correctly claims no product exception for itself (owner marker only, no dropout
  history/snapshot — verified: the idea's history contains the kickoff with kimi owner-excluded
  before Phase 0, so no mid-idea automatic transition exists to derive from). Nothing to
  refute yet; no completion is asserted by codex-1 (IMPLEMENTATION.md:47).
- **AC14 (post-close delivery).** Expressly post-close; no delivery is claimed and none is
  verified here. PRs #78/#10 exist and are open (gh-verified); merging main, releases,
  Homebrew/winget/19 targets and the separate channel verification all remain post-close work.
  Nothing to refute yet.

## Findings

### [MINOR] Standalone mid-idea signoff now hard-depends on the prior run's manifest
`internal/app/quota_signoff.go:56-64`: a new standalone signoff run loads the prior
(kickoff/last-batch) run's manifest and aborts the whole signoff on any load error. On 1.52.0
this dependency did not exist, so an idea whose prior driving-run manifest is missing or
unreadable now cannot collect survivor signoffs at all, even though the code comment's own
philosophy is "missing old snapshots stay missing and cannot earn the single-reviewer
exception" — i.e. an absent snapshot could safely degrade to empty inheritance rather than
aborting. Fail-closed direction, clear error, no data loss, and modern runs always have
manifests, so this is an edge regression only. Suggested fix: treat `os.IsNotExist` on the
prior manifest as empty inheritance (with a warning), keeping hard failure for corrupt/foreign
manifests.

### [NIT] full-host-tests-v2 receipt needed recovery; on-disk *-result.json unreadable
`full-host-tests-v2-result.json` on the shared volume returns ENOENT despite the directory
listing it; the pass is evidenced by the recovered stdout receipt plus the hash-verified
complete jsonl (I matched sha256 04e3fa05… and recounted packages/events myself), documented
in `full-host-tests-v2-receipt-recovery.md` and committed `implementation-validation.json`.
The accounting is honest and verifiable, but the canonical on-disk receipt artifact for the
headline suite run is a recovery rather than the promised at-exit file; worth a hard look if
the shared volume repeats this.

### [NIT] Missing blank line before new heading in docs/agent-cli-mechanics.md
The new `## Participant-step timing and buffered output (1.53.0)` heading follows the
pre-existing last line ("…the `buffers_stdout` flag).") with no blank line
(docs/agent-cli-mechanics.md:93-94; the file previously ended without a trailing newline).
Markdown renders, but it is malformed relative to the file's own style.

### [NIT] Legacy-rule negatives not exercised at the review-gate seam
`reviewFailureEvidence`'s legacy branch (review_gate.go:170-180) enforces reset ≥ 60 min for
`quota.named-reset-ge-60m.v1` and no-reset for the other two rules, adapter `zcode`, and the
frozen provenance string, but the adversarial table in `review_gate_test.go` only mutates the
participant-failure path (plus `arbitrary-rule` and a positive legacy row per rule). A short
reset, a contradictory reset on a no-reset rule, or a non-zcode adapter would currently be
caught only by code reading, not by a failing fixture. One table row each would pin it.

## Open questions

1. For the organizer, not blocking: when MINOR-1's edge triggers on a real deck, is the
   intended owner recovery "restore the prior run manifest" or should a future idea relax
   `IsNotExist` to empty inheritance? The comment in quota_signoff.go suggests the latter was
   considered.
2. The 775.7 s buffered zcode invocation documented as the buffering witness also means a
   genuinely hung zcode signoff is only detected at its hard ceiling (as disclosed). FINAL
   defers timeout tuning to a later idea — should that follow-up also record the fresh
   goal-check durations this run produces, as FINAL's Known risks suggests?
3. AC13's fresh goal process will re-read this review; if the fix-up cycle for MINOR-1 touches
   `quota_signoff.go`, the same-tick rebind test should be rerun (it is the covering test for
   that seam).

---

## Adjudication note (2026-10-09, phase-7 deliberation launch)

Appended by zcode-1 after cycle-1 consensus drafting; nothing above is altered (the
`## Refutation attempts` heading and all findings stand as filed). Launch attestation:
context_mode=full,
source_sha256=packet_sha256=acbd4dbc0c0702bc191176bb80bcee32c5093c8b4ebbee42e036df6a9b7d1137,
fallback absent — re-hashed by me this launch (`shasum -a 256`, match).

### [MINOR] "Standalone mid-idea signoff now hard-depends on the prior run's manifest" — WITHDRAWN

The implementer rebuttal (`source-context/implementer-response-minor1.md`) is correct and my
finding's baseline premise was wrong. Verified PRIMARY this launch against both trees:

- Baseline 85c5a8f has the same standalone-creation path and then calls
  `membership.Before` → `Reconcile` → `reconcileLocked` → `validateManifests`
  (`git show 85c5a8f:internal/app/quota_signoff.go`; membership.go:68/74/82/93 in that
  tree). `validateManifests` builds the required run set from {current run,
  `h.Kickoff.RunID`, every `h.Batches[*].RunID`} (membership.go:160-163) and fails on any
  load error (`quota manifest %s: %w`, :177-179), under the shipped comment "Missing
  original run identity is an integrity gate, not an opportunity to repair legacy gap 11"
  (:158-159). So on 1.52.0 an idea with a missing prior (kickoff/last-batch) manifest ALSO
  could not collect survivor signoffs — it aborted at `validateManifests` (wrapped by
  `Before`'s IntegrityBlock) instead of at the new load. No pre-existing continuation path
  regressed. My sentence "On 1.52.0 this dependency did not exist" asserted baseline
  behavior I had not actually checked: it was RECALL wearing my review's PRIMARY label —
  my error, and precisely why this independent re-check was warranted.
- The integrity lines are unchanged by this idea: `git diff 85c5a8f..a26f588 --
  internal/membership/membership.go` changes exactly one unrelated line (Settle's
  `CheckGates` → `checkGates`, :314); the `quota_signoff.go` diff is exactly the +14-line
  snapshot-inheritance block. Catching `os.IsNotExist` only at the new prior-manifest load
  could not enable continuation anyway: `priorRun` is the kickoff or last-batch run, both
  always present in `validateManifests`' required set, so the same missing manifest would
  still abort the signoff one gate later.
- Relaxing `validateManifests` instead would change shared, unrelated history rules: it also
  gates `Settle` (membership.go:335) and both policy-revision paths (revision.go:65, :173) —
  identical call sites in baseline — while FINAL D2 ("preserve normal snapshot binding …
  Missing snapshot does not qualify") and D4 ("Regress these seams rather than redesigning
  shipped membership/retry accounting") scope this idea to regression coverage of the
  seams, not to redesigning that gate.
- Adjacent consumer checked for completeness: `requestConsensusSignoffs`
  (consensus_request_signoffs.go:131-137) loads the *current* run's manifest — which
  `quotaSignoffStart` has already created — so no separate prior-manifest dependency
  exists there.

Residual issue, folded into the amended cycle-1 plan rather than left as an open finding:
the new comment at quota_signoff.go:52-55 conflates a missing *snapshot* (valid manifest,
absent `RosterSnapshot` → inherits empty, never qualifies the exception) with a missing
*manifest* (integrity gate, blocks). The amended plan keeps the hard failure at the new
load, clarifies that comment, and adds explicit standalone tests: absent-snapshot →
inherits empty + diagnostic, never earns the exception; missing/corrupt/foreign prior
manifest → still blocks; plus the existing same-tick rebind test
(`TestReviewGateSignoffWatchdogAndSameTickRebind`) as the covering test (this also resolves
my Open question 3 above). The two NIT fixes and the disclosed receipt follow-up are
unchanged.

Adjudication is recorded in my revised own draft `review/consensus.md` (revision 2,
including a §15.3 `## Verdict conflicts` record). The prior draft and my prior signoff stay
byte-identical in `review/consensus-cycle-01-proposed.md` (sha256
5e9baccc8fca2b2566f224d58f3d10d2dbf174045baf8ec80e7656d3f46bc6b3, verified by
`cmp`/`shasum` this launch). The driver refused `consensus reopen` (triage=partial,
organizer-reported); the CLI reopened nothing — this is the drafter's own revision of its
own unratified draft under organizer authorization, replacing its own prior signoff with a
fresh verdict on the revised text.
