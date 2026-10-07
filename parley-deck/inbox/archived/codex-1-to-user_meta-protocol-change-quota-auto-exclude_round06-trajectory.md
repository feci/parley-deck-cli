---
from: codex-1
to: user
idea: meta-protocol-change-quota-auto-exclude
phase: review-round-06
blocking: yes
date: 2026-10-06
---

## Question

Authorize **one narrow fix-up cycle 4** for R6-MAJOR-1, R6-MINOR-1 and R6-MINOR-2,
with claude-1 signing the repair plan and then independently re-reviewing the full
implementation? This is the recommended disposition. This note is a trajectory stop,
not the NEW attended-close request and not permission to release.

## Context

The authorized cycle 3 is delivered at CLI product `fac40aaa1bba5350351bab0310ea8ede30ed3e6e`
(review snapshot `25ea1da73ab0eab5e6099f759e2720aff54ada88`) and skill
`e976f7c9515f3250761380c1e09295e0f6c59985`. The separate configured claude-1 process
completed successfully after 1227.2 seconds, with no timeout, quota error or retry.
Its canonical review is [round-06/claude-1.md](../ideas/meta-protocol-change-quota-auto-exclude/review/round-06/claude-1.md),
SHA256 `63fd03de97f24a850f0c96b1d9fe7c1469c02850e1dd9b00eff8e2f41829e59e`.

**Finding trajectory: 15 → 6 → 3 → 3. Round 06: 0 CRITICAL, 1 MAJOR, 2 MINOR, 0 NIT.**

1. **R6-MAJOR-1, new on cycle-3 notice code.** An owner appends an answer to an applied
   exclusion notice, either live or archived. Exact-byte checking calls the notice
   corrupt and gates status, surviving-member signoff and driving. The protocol permits
   appended owner answers and calls inbox notices non-authoritative. The reviewer
   reproduced it on both volumes. The raw cycle-2 comparison is precise: a live annotated
   notice already blocked Before but left status/signoff usable; an archived annotated
   notice did not gate and was republished. Cycle 3 widened the gate to all surfaces.
2. **R6-MINOR-1.** Re-including a kickoff-excluded id during round 1 works if the stale
   exclusion marker remains, but removing that marker in the same edit leaves the id
   pending catch-up and refuses its signoff. Baseline accepts the edit; both volumes
   reproduce the difference.
3. **R6-MINOR-2.** The CHANGELOG and draft release notes overclaim preservation of
   policy-off behavior. W1 required every remaining difference to be owner-visible.

The reviewer independently passed the full Go suite (528 seconds), build, vet, race,
Windows cross-build, shared/local regressions, full skill suite (399 Node and 54 Python
checks), full phase0/5/6/8 packets and actual supervisor/worker crash recovery on both
volumes. A live-writer variant correctly refuses recovery. Complete product diffs were
read: 117 CLI files and 5 skill files. Passing tests do not dispose of the new probes.
Windows runtime and container/PID namespaces remain unverified. Native zcode capture
was neither attempted nor authorized.

D1–D4, the §5 decline route, ordinary notice archive/delete, W4's honest AC2/manual-log
wording and W5's inactive follow-up now pass their independent scoped checks. The raw
review and IMPLEMENTATION contain the full current AC1–AC21 map. AC11/AC15/AC16 remain
partial and AC21 remains unmet. AC2 is **NOT MET / explicitly owner-waived for this
release**; R5-MAJOR-2 is owner-accepted and deferred, never fixed or PASS.

### Proposed bounded cycle 4 (not started)

- **Notice ownership:** immutable transition/receipt evidence determines membership and
  publication state. An owner's edited non-authoritative notice must not gate status,
  signoff or driving. Preserve annotations and archive/delete semantics. Before a valid
  receipt, handle incomplete publication without overwriting owner content or creating
  a permanent idea-wide gate. Retain safe publication-path checks. Test appended answers
  live/archived, before/after receipts, repeated recovery and one terminal evaluation.
- **Round-1 return:** remove the dependency on retaining a mutable stale exclusion marker.
  Check what authentic kickoff evidence exists and preserve it; do not invent historical
  membership or owner authority. The signed plan must state any legacy evidence limit
  and preserve the existing policy-on/known-excluded/retained-veto boundaries.
- **Honest release wording:** list every remaining policy-off change against 27e42b8:
  incomplete-joiner signing/quorum restrictions, declined joiner in Missing, later-return
  catch-up, pending-catch-up driving gate, frozen closed-idea edits, and any marker-return
  limitation that remains. Reconcile CLI docs, skill guidance and release notes.
- Same codex-1 implementer/organizer and separate claude-1 reviewer. Both final signoffs
  and a NEW attended owner close remain required after successful re-review. Any fresh
  CRITICAL/MAJOR on cycle-4 fixes returns to the owner; no automatic cycle 5.

No native-provider capture, grammar relaxation, roster/model/provider/config change,
new worktree, D6 migration, global core publication or release is part of this scope.

## Why an owner answer is required

Your binding [round05 answer](user-to-codex-1_meta-protocol-change-quota-auto-exclude_round05-answer.md)
states: **“A fresh CRITICAL or MAJOR on the cycle-3 fix code escalates again, as you
proposed.”** The signed cycle-3 plan makes no automatic cycle 4 available. The [Parley Deck skill](/Users/tomasfecko/.codex/skills/parley-deck/SKILL.md)'s escalation rule says **“If blocking: yes, pause the escalating agent's work
for that idea.”** The numeric five-cycle cap does not override this trajectory stop.

## What I need from you

- **Recommended:** authorize the bounded cycle 4 above, with signed plan, independent
  full re-review and the unchanged NEW attended-close requirement.
- Alternatively, explicitly accept the notice-edit restriction and return limitation
  for this release, authorizing only the truthful release/docs correction and its review.
  This accepts a known MAJOR; it does not repair it or waive the final close.
- Or stop/defer this release and give a different disposition.

All participant processes have exited. Product and prior reviewer/signoff bytes are
frozen. The signed cycle-3 plan stays intact at review/consensus.md and its archived
copy; no cycle-4 consensus, repair, final approval, merge, release or installation has
been manufactured. Usage and evidence copies are recorded under the idea. The organizer
ends deliberately at this blocking note so the owner relay can resume it after an answer.
