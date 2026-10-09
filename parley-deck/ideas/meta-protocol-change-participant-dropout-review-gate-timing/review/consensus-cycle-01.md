---
idea: meta-protocol-change-participant-dropout-review-gate-timing
review-cycle: 1
outstanding_agreed_fixes: 3
blocked: false
drafted-by: zcode-1
date: 2026-10-09
reviewed-commit: a26f588
---

## Agreed fixes

Drafter's note (revision 2 of this own draft): this cycle-1 consensus codifies the fix plan
for the round-1 findings, amended after the implementer's MINOR-1 rebuttal
(`source-context/implementer-response-minor1.md`) was independently verified and accepted
(see `## Verdict conflicts` and `## Drafter position changes`; the full withdrawal rationale
with PRIMARY evidence is appended to `review/round-01/zcode-1.md`). The first revision of
this draft and its zcode-1 signoff are preserved byte-identically at
`review/consensus-cycle-01-proposed.md` (sha256
5e9baccc8fca2b2566f224d58f3d10d2dbf174045baf8ec80e7656d3f46bc6b3, verified by
`cmp`/`shasum` this launch) and are superseded by this text. The two NIT fixes and the
deferred receipt follow-up carry over unchanged.

- From `codex-1/implementer-response-minor1` (rebuttal accepted in full by zcode-1
  adjudication; supersedes the fix text this draft first recorded under the now-withdrawn
  round-1 MINOR) — **standalone-signoff seam hardening, `internal/app/quota_signoff.go`**:
  1. **Preserve the shipped missing-manifest integrity gate.** Keep the hard failure at the
     new prior-run `runmanifest.Load` (quota_signoff.go:60-64) for missing, corrupt and
     foreign-identity prior manifests; do NOT add an `os.IsNotExist` catch there. Verified
     basis: baseline 85c5a8f reaches the identical gate through `membership.Before` →
     `Reconcile` → `reconcileLocked` → `validateManifests` (membership.go:68/74/82/93 →
     :157-197; required runs = current + kickoff + every batch at :160-163; any load error
     aborts at :177-179 — "Missing original run identity is an integrity gate, not an
     opportunity to repair legacy gap 11"), those lines are unchanged in a26f588, and a
     tolerant catch could not enable continuation because the same missing kickoff/batch
     manifest still fails `validateManifests` one gate later. Relaxing `validateManifests`
     instead would change the shared history gate behind `Settle` (membership.go:335) and
     both policy-revision paths (revision.go:65, :173) — outside FINAL D2/D4 scope.
  2. **Clarify the comment** at quota_signoff.go:52-55 to distinguish a missing *snapshot*
     (manifest loads, `RosterSnapshot` absent → inherits empty; empty never invents models
     and never qualifies the single-reviewer exception, per FINAL D2 "Missing snapshot does
     not qualify") from a missing *manifest* (integrity gate; blocks the standalone signoff,
     exactly as already on 1.52.0).
  3. **Explicit standalone tests at this seam:** (a) valid prior manifest with an absent
     snapshot → empty inheritance + clear diagnostic, and the exception never derives;
     (b) missing prior manifest → the signoff still blocks; (c) corrupt prior manifest →
     blocks; (d) foreign-identity prior manifest → blocks; (e) rerun/retain the existing
     same-tick rebind test (`TestReviewGateSignoffWatchdogAndSameTickRebind`) as the
     covering test (resolves round-01 Open question 3).
- From `zcode-1/review/round-01` **[NIT] "Missing blank line before new heading in
  docs/agent-cli-mechanics.md"** — insert the Markdown blank line before the
  `## Participant-step timing and buffered output (1.53.0)` heading.
- From `zcode-1/review/round-01` **[NIT] "Legacy-rule negatives not exercised at the
  review-gate seam"** — add adversarial table rows at the `reviewFailureEvidence` seam
  for: a short reset on `quota.named-reset-ge-60m.v1`, a contradictory reset on a
  no-reset rule, a non-zcode adapter, and wrong provenance. Preserve the positive
  boundary cases and all three enumerated legacy rules. (This adds wrong-provenance to
  the three rows the finding suggested — accepted.)

## Verdict conflicts

- Claim: "On 1.52.0 this dependency did not exist, so an idea whose prior driving-run
  manifest is missing or unreadable now cannot collect survivor signoffs at all" — author
  zcode-1, `review/round-01/zcode-1.md`, finding [MINOR] "Standalone mid-idea signoff now
  hard-depends on the prior run's manifest". Tag: RECALL — the baseline behavior was
  asserted without a baseline citation; the round-1 PRIMARY label covered the a26f588-tree
  reads, not this premise.
- Counter-claim: "absent prior manifests already block standalone signoff in
  baseline85c5a8f. Its quotaSignoffStart calls membership.Before, then
  Reconcile/reconcileLocked, then validateManifests … errors on any missing run.json …
  These lines are unchanged by this idea." — author codex-1,
  `source-context/implementer-response-minor1.md`. Tag: PRIMARY (named functions and
  files, independently re-derived below).
- Resolution (zcode-1, PRIMARY, this launch): codex-1 CONFIRMED, the round-1 premise
  WRONG. `git show 85c5a8f:internal/app/quota_signoff.go` and
  `85c5a8f:internal/membership/membership.go` show the identical chain and an identical
  `validateManifests` (required runs {current, kickoff, every batch}; any load error →
  abort); `git diff 85c5a8f..a26f588 -- internal/membership/membership.go` touches
  exactly one unrelated line (Settle's `CheckGates`→`checkGates`); `runmanifest.Load`
  (internal/runmanifest/manifest.go:248-252) returns the raw `os.ReadFile` error, so the
  same missing prior manifest fails in both trees — baseline later at `validateManifests`
  (wrapped by `Before`'s IntegrityBlock), current earlier at the new load. The conflict is
  resolved: no supported pre-existing path regressed, and the finding is withdrawn by its
  author (see `## Dismissed findings`).

## Deferred follow-ups

- From `zcode-1/review/round-01` **[NIT] "full-host-tests-v2 receipt needed recovery;
  on-disk *-result.json unreadable"** — defer the investigation of the shared-volume
  ENOENT behavior to a follow-up idea (slug `TBD`; open it if the volume repeats this).
  Disposition rationale, which I concur with openly: my own PRIMARY recount re-hashed
  the complete actual exit-0 jsonl (sha256 match, 3,315 distinct test/subtest IDs
  recounted: 3,311 pass + 4 skip, 34/34 packages, 0 fail events), so the pass is
  independently established from the full raw log; the failure mode is disclosed
  evidence-storage trouble on the shared volume, not a product defect in this diff and
  not a test exemption. All original and failed logs, the exact recovered emitted
  receipt, and the hash/recount proof stay retained
  (`.parley-runtime/review-gate-timing/full-host-tests-v2*`, committed
  `source-context/implementation-validation.json`).
- From `zcode-1/review/round-01` **Open question 2** — the deferred timeout-tuning
  follow-up idea should record the fresh goal-check durations this run produces (per
  FINAL.md "Known risks"), alongside the already-disclosed fact that a genuinely hung
  buffered zcode signoff is detected only at its hard ceiling.

## Dismissed findings

- `zcode-1/review/round-01` **[MINOR] "Standalone mid-idea signoff now hard-depends on the
  prior run's manifest"** — withdrawn by its author on 2026-10-09 after verifying the
  implementer rebuttal (rationale and PRIMARY evidence in `## Verdict conflicts` above and
  in the adjudication note appended to `review/round-01/zcode-1.md`). The actionable
  residue — the snapshot-vs-manifest comment clarity and the missing regression tests at
  this seam — is carried as the first Agreed fix above, so nothing actionable is lost by
  the withdrawal. No other finding was withdrawn and none was judged not-an-issue.

## Drafter position changes

- Round-1 MINOR "Standalone mid-idea signoff now hard-depends on the prior run's
  manifest" (`review/round-01/zcode-1.md`): prior position — open MINOR with suggested
  fix "treat `os.IsNotExist` on the prior manifest as empty inheritance (with a warning),
  keeping hard failure for corrupt/foreign manifests". New position — withdrawn; the
  suggested catch would be ineffective for its stated purpose (`validateManifests` still
  blocks the same missing manifest) and misleading (it implies continuation without the
  prior manifest is possible). The amended plan instead preserves the integrity gate and
  adds the comment clarification plus the snapshot-vs-manifest regression tests.
- All other positions are unchanged from `review/round-01/zcode-1.md` and from revision 1
  of this draft: both NIT fixes, the deferred receipt follow-up, and the deferred
  goal-duration recording.

## Coverage & blind spots

- **Reviewer coverage this round:** one independent reviewer (zcode-1). kimi-1 was
  owner-excluded before Phase 0 (recorded in `00-prompt.md`), so this floor run had
  exactly one non-implementer participant; the run correctly claims no product
  single-reviewer exception for itself (no mid-idea automatic transition exists to
  derive one from — verified in the round-1 review, AC13). Cross-reviewer overlap is
  therefore undefined this round; coverage rests on the round-1 refutation sweep across
  all 14 ACs plus independent reruns of the focused Go suites, the packet/drift guard
  tests, `go build ./...` and `go vet ./...`, and now the PRIMARY baseline-vs-current
  re-verification behind the MINOR-1 adjudication.
- **Blind spots disclosed in the review:** no attempt was made to hand-forge a history
  file (tampering is outside this rule's threat model; stops-for-repair is the shipped
  stance); the skill-suite pass at 74cc831 is the dependency-refreshed v3 run after
  `npm ci` (v1/v2 failures retained as history); protocol/skill text equality was
  established by independent byte-exact reconstruction, not by eyeball. New this
  revision: the round-1 MINOR's baseline premise was itself unverified (RECALL) until
  this adjudication — a reminder that "unchanged from baseline" claims need baseline
  citations.
- **Process state after this consensus:** this is a fix cycle, not a close. There is no
  codex-1 signoff on this revised draft yet; codex-1 signs only after weighing this
  adjudication, then implements the agreed plan and requests a fresh independent
  full-scope re-review. Still required before completion: that fresh re-review of the
  fix-up diff, a later Phase 7 consensus with zero Agreed fixes, both current
  participants' final review-consensus ACCEPT blocks, the fresh zcode-1 goal-check
  process, and the current-tree AC evidence (owner-attended close authority per
  `00-prompt.md` Constraints).

## Process provenance (revision 2)

- The driver refused `consensus reopen` for this cycle (triage=partial). This is
  organizer-reported in the phase-7 dispatch brief; I located no durable run-log record
  of the refusal and claim none. No CLI reopen occurred and none is claimed: this
  revision is the organizer-authorized fallback in which the original drafter and only
  signer (zcode-1) revised its own unratified draft — which carried no other agent's
  signoff — and replaced its own prior signoff with a fresh verdict on the revised text.
  No other agent's artifact or signoff was edited.
- A preceding native review-consensus-01b process was deliberately cancelled by the
  organizer within seconds to correct an inaccurate prompt premise about driver reopen;
  it supplied no accepted output and was not a timeout/provider failure. This
  adjudication launch supersedes it (organizer-reported).
- The organizer has not signed and has not changed product; product remains CLI a26f588 /
  skill 74cc831. In this step zcode-1 wrote no product files, made no commits, and sent
  no GitHub messages. Launch attestation: context_mode=full,
  source_sha256=packet_sha256=acbd4dbc0c0702bc191176bb80bcee32c5093c8b4ebbee42e036df6a9b7d1137,
  fallback absent (packet file re-hashed by me this launch, match).

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: zcode-1 — 2026-10-09
Status: ✅ ACCEPT
Notes: This signoff is my fresh verdict on revision 2 of this own draft; it supersedes my
signoff on revision 1, which is preserved byte-identically (with that draft) in
review/consensus-cycle-01-proposed.md. ACCEPT approves the cycle-1 fix plan only, not
product close. I verified the implementer's MINOR-1 rebuttal against baseline 85c5a8f and
current a26f588 (PRIMARY; commands and line cites in ## Verdict conflicts and in the
adjudication note appended to my review/round-01 file): the missing/corrupt/foreign prior
manifest already blocked standalone signoff on 1.52.0 through the unchanged
validateManifests integrity gate, my original MINOR rested on an unverified baseline
premise, and I withdraw it; the amended first Agreed fix preserves that gate, clarifies
the snapshot-vs-manifest comment, and pins the seam with explicit tests. I concur with
both NIT fixes and with deferring the shared-volume ENOENT receipt investigation (my
PRIMARY recount re-hashed the complete exit-0 jsonl, so the headline suite pass is
independently established and the failure is disclosed storage trouble, not a test
exemption). Awaiting codex-1's weighing of this adjudication and signoff, the fix-up
implementation, a fresh full-scope independent re-review, a zero-fix consensus, both
final ACCEPT blocks, the fresh goal-check process and current-tree AC evidence.

### Signoff: codex-1 — 2026-10-09
Status: ✅ ACCEPT
Notes: I accept the revised cycle-1 plan after weighing the independently verified baseline rebuttal. Preserve the historical manifest integrity gate; clarify absent snapshot versus absent manifest, add the requested diagnostics and negative/seam coverage, and fix the Markdown blank line. I concur with the disclosed receipt-storage follow-up and retained raw evidence. This approves implementation of the fix plan only, not close or a self-review verdict.
