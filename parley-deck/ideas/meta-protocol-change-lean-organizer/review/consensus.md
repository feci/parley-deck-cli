---
idea: meta-protocol-change-lean-organizer
review-cycle: 3
outstanding_agreed_fixes: 7
blocked: false
drafted-by: zcode-1
date: 2026-09-24
reviewed-commit: 913f8ba
---

<!-- blocked: false and outstanding_agreed_fixes: 7 are the DRAFTER'S PROPOSAL, not a
     resolution. The two reviewers DISAGREE about readiness for closure (VC-7) and
     about the materiality of claude-1's R3-MIN-1 (VC-8). Both conflicts stay open
     until the signoffs decide them — never by count, drafter, or organizer. Signoffs
     are empty in this draft. -->

Draft of the Phase-7 review consensus for review cycle 3 (after fix-up cycle 2; reviewed
commit `913f8ba`, fix-up-2 source `998346c`). The previous signed consensus is archived
**VERBATIM** at `review/consensus-cycle-02.md` (55,430 B, `cmp`-silent against the file it
replaces, sha256 `b4976343…f6d1d4`, including all three signoff blocks and reservations);
this fresh draft reuses no stale signoff. Inputs read in full: the live phase-7 packet
attested below; frozen `FINAL.md` (frozen at `120a9bf`; `git diff 120a9bf..913f8ba --
…/FINAL.md` empty — re-verified this session); `IMPLEMENTATION.md` at `913f8ba` including
all of fix-up cycle 2; both complete round-03 review files
(`review/round-03/claude-1.md`, `review/round-03/kimi-1.md`); the entire signed cycle-2
consensus (now `review/consensus-cycle-02.md`); `review/consensus-cycle-01.md`;
`00-prompt.md`; `organizer-notes.md`. No code, FINAL, release-plan, reviewer-file,
organizer-note, or inbox edits were made in this invocation; this file (plus the verbatim
archive of its predecessor) is the only output. Nothing is committed, tagged, installed,
or published, and `parley protocol publish` is not run.

**Protocol context attestation (Phase-7 drafting):**

```json
{"context_mode": "full", "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "fallback_reason": null, "body_path": "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer/.parley-runtime/protocol-packets/full-phase7-deliberation-8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7.md"}
```

Drafter cross-check (PRIMARY): `shasum -a 256` over the packet body returns
`8ce83cde…a9db7` — the attested packet IS the live authority; both round-03 reviews attest
the same hash against their own trees.

**Drafter disclosure (§15.1).** The drafter is the implementer (zcode-1), drafting per
Phase 7's default. The implementer issues **no verification verdict on its own
implementation**: every disposition below is a **PROPOSAL** awaiting the reviewers'
independent evaluation and the signoff gate. Where the reviewers materially disagree —
this round: readiness for closure (VC-7) and the materiality of R3-MIN-1 (VC-8) — the
conflict is recorded under `## Verdict conflicts` and stays open; nothing is resolved by
the drafter's self-judgment, by participant count, or by the organizer (verdicts read
only). No brief suppressed or narrowed any finding: all 7 claude-1 findings and kimi-1's
zero-findings position (with its recorded G4(a) observation) carry an explicit disposition
below, and the raw review files remain canonical. The drafter's positions are unchanged
since its cycle-2 signoff block (archived verbatim); no drafter position changes to
record. Drafter-verified record facts this session (tagged DRAFTER-PRIMARY inline):
packet/authority hash equality above; HEAD `913f8ba`; FINAL frozen; the archive's byte
identity; `phasedigest.go`'s strict `After` comparison; the `runs/` record shapes; the
`auto_implement: true` frontmatter line; the driver close-gate lines cited in H3.

## Agreed fixes

*(Proposed — fix-up cycle 3 (cap 5 on `deliberation`; this is cycle 3). Six of the seven
items are record corrections or documentation on this idea's own artifacts; the one code
item (H1) touches the same file and surface the signed G4 already changed, inside frozen
FINAL B. None changes the packet map, any byte cap, the guardrail's phase-1 scope, or any
frozen protocol text. No protocol-text edit is proposed, so no core restage is expected;
if signoffs add one, the restage, its independent recheck, and the core-publish note
update happen together in the same step. H1 is contingent on VC-8: it is proposed because
one reviewer filed it at fix-worthy severity with PRIMARY provenance and the other
reviewer's live-surface verification is not incompatible with it; a signer who sustains
the lesser position amends or strikes it at signoff — that is the conflict closing, not
this draft.)*

### Finding → disposition map (complete)

All 7 findings filed this round (all claude-1; kimi-1 files none), severity and
provenance kept. Prior-round outcomes for traceability: 11 of 11 claude-1 round-02
findings withdrawn-as-fixed by their filer except R2-MIN-1, whose (a) half is narrowed
into R3-MIN-1 (R2-NIT-3 remains retired-by-owner); all 4 kimi-1 round-02 findings
verified fixed (→G9/G10/G2/G4).

| Finding (severity, provenance) | Disposition |
|---|---|
| claude-1 R3-MIN-1 (G4(a) inert when mtimes equal — every fresh clone/checkout reads `await implementation`; silent wrong-direction failure on the organizer's branch surface) [MINOR, PRIMARY — fresh-checkout reproduction at `913f8ba` with `stat` output quoted; mechanism isolated at `phasedigest.go:204` strict `After`] | Fix H1 (contingent on VC-8); narrows R2-MIN-1(a) |
| claude-1 R3-MIN-2 (`## Fix-up cycle 2` lacks the Phase-8-required `### Deviations from agreed fixes`; the G1 figure divergence belongs there) [MINOR, PRIMARY — `grep -i deviat` over the cycle-2 span; Phase-8 shape quoted] | Fix H2 |
| claude-1 R3-MIN-3 (`auto_implement: true` close conditions (LE-7/LE-11) unaccounted for in any round-02 artifact or the signed consensus's Cycle-mechanics paragraph) [MINOR, PRIMARY — code read `internal/driver/impl.go:283-297`; artifact scan finding no commissioned goal-done check] | Fix H3 (documentation + commissioned independent check; no code change, no flag change) |
| claude-1 R3-NIT-1 (23,230 ↔ 23,223 reconciliation arithmetically unexplained: per-block newline would add 28, not 7) [NIT, PRIMARY — block count per section; includes filer's §15.1 SELF-CORRECTION superseding its own 23,223 B] | Fix H4 |
| claude-1 R3-NIT-2 ("19,736 B mislabelled `floor`" attaches the label to the wrong number — the `floor` variable held 15,759 B; 19,736 B was the `+§2 reference`) [NIT, PRIMARY — both quantities re-measured under the old code's own convention] | Fix H5 |
| claude-1 R3-NIT-3 (run-record residual stale and mischaracterized: a post-F7 run exists; recent runs write no `run.json` at all, not "zero-width") [NIT, PRIMARY — `runs/` records re-dated and re-listed] | Fix H6 (facts DRAFTER-verified this session) |
| claude-1 R3-NIT-4 (record overstates the G4(a) test's second assertion — the relaxation path asserts enumeration validity only) [NIT, PRIMARY — test body quoted] | Fix H7 |
| kimi-1 — no findings; recorded observation (not a finding): G4(a) fresh-checkout fall-through "remains defensible… the live-deck operational surface behaves correctly" | No fix item; the observation enters VC-8 as the minority materiality position on R3-MIN-1 |

**Organizer-owned corrections: none required this round.** No finding names an
organizer-owned artifact; R2-NIT-3 stays resolved-by-owner from cycle 2.

### Proposed fix plan

- **H1 — from claude-1 R3-MIN-1 (MINOR; contingent on VC-8).** In
  `internal/driver/phasedigest.go`, stop keying the G4(a) fix-up-published branch on
  mtime alone. **Primary shape (claude-1's suggested fix and its open-question 3 lean):
  the content signal.** `d.Implementation.Status` is `fix-up-cycle-N` and `d.Review`'s
  latest complete round is `round-0M`; when a fix-up cycle is published and `M ≤ N`, a
  further review round is awaited regardless of timestamps (an incomplete next round is
  already caught earlier by `Completed < Total`; a complete `round-0M` with `M > N`
  reviewed that fix-up and falls through exactly as today). **Signoff-amendable
  alternative:** keep mtime and relax `implInfo.ModTime().After(newest)`
  (`phasedigest.go:204`) to `!…Before(newest)` so equality engages. Either way, add the
  missing test case: the fresh-checkout shape (no next-round directory; `IMPLEMENTATION.md`
  and every latest-round artifact sharing ONE mtime — what `git clone`/`git checkout`
  produce) must read `await review artifact`; today it reads `await implementation`
  (claude-1 PRIMARY; the shipped test pins strict ordering via `os.Chtimes`, so the
  equal-mtime case is never exercised). **Record correction (mandatory under either
  shape):** `IMPLEMENTATION.md:503-508`'s unconditional "Fix-up G4(a) retired the
  residual" claim is corrected transparently — the cycle-2 behavior was conditional on
  mtime ordering (equal mtimes left the enumeration at `await implementation`), the
  imprecision is owned, and the correction states which shape restored the claim (content
  signal: unconditionally; `>=` relaxation: including equality). No false claim that the
  retired behavior predates this fix. Grounds: the failure direction is silent and wrong
  on both organizer surfaces (`--json` digest, `parley organizer brief`) — the opposite
  of G8's loud-failure mtime approximation; the live deck is unaffected only via the
  older directory-based path (empty `round-03/`), which is state-dependent, not a
  guarantee.
- **H2 — from claude-1 R3-MIN-3's sibling R3-MIN-2 (MINOR).** Add the Phase-8-required
  `### Deviations from agreed fixes` subsection to `## Fix-up cycle 2` in
  `IMPLEMENTATION.md` (cycle 1 carries one at `:494`; cycle 2 has none). Content: one
  bullet recording that G1's delivered floor figure (**52,295 B**) diverges +10.2 KB /
  +24 % from the "≈ 42 KB with §2" named by the signed G1 plan text and FINAL's
  acceptance row, why (the round-1 ≈ 42.1 KB was a derivation subtracting a
  whole-section quantity from an already-within-section-optimized body — 65,516 − 27,420
  + §2 at `ffa4587`; both reviewers independently hold the recomputed whole-block
  retained total is the conceptually correct realization and the criterion met as
  frozen — claude-1 withdrew R2-MAJ-1; kimi-1 "criterion met as frozen; no wording
  amendment needed"), and that FINAL stays frozen; `None` for the other nine fixes. The
  divergence is currently disclosed in two narrative places (the G1 "Recorded-number
  note for round 3"; the closure-conditions residual) but not in the one place Phase 8's
  shape puts decision-relevant divergence. Per claude-1's open question 1, the
  agreed-fixes deviation bullet is sufficient; no `## Deviations from FINAL.md` entry is
  required (that section stays for mechanism changes like G7's).
- **H3 — from claude-1 R3-MIN-3 (MINOR; documentation + commissioned verification; no
  code change, no `auto_implement` change, no manufactured pass).** Two halves:
  **(a) Document the real close conditions** — the subsection immediately after this
  plan records them, and fix-up cycle 3 adds a short pointer line in `IMPLEMENTATION.md`
  so the closing consensus is planned against the complete rule set. **(b) Arrange the
  independent LE-7 goal-done check before close** — requested of the organizer here, not
  performed by the implementer: a **fresh non-implementer participant session**
  (claude-1 or kimi-1; never zcode-1 the implementer, never codex-1 — FINAL A's
  role-ineligibility excludes the declared facilitator from the goal-done role, enforced
  by the driver's predicate) verifies FINAL.md's observable acceptance criteria with its
  own commands and §15.2 PRIMARY provenance **at the closing HEAD** (the tree the zero-fix
  consensus will certify — after fix-up cycle 3 lands; a check at `913f8ba` now would be
  staled by any cycle-3 edit), filing its verdict in its own participant-owned canonical
  artifact (proposed: `review/goal-done/<agent-id>.md`, review-file shape with
  `goal-check: true`; signoff-amendable), editing nothing. A confident fail or
  inconclusive verdict blocks close and escalates to the user; a pass is cited as
  evidence by the closing consensus. Until it runs, nothing in this idea may claim the
  close conditions satisfied.
- **H4 — from claude-1 R3-NIT-1 (NIT).** Correct the G1 clause "claude-1's independently
  measured 23,223 B **plus this convention's per-block newline**": the seven sections
  contain 28 blocks, so a per-block newline would add 28, not the observed 7; the
  measured cause is one byte per *section* at the round-02 extraction boundary.
  Replacement clause (suggested): "claude-1's round-02 23,223 B was one byte short per
  section at the extraction boundary; the exact source total is 23,230 B." The filer's
  §15.1 SELF-CORRECTION (23,223 → 23,230, a weakening, effective immediately) is already
  on record in `review/round-03/claude-1.md`; this item corrects the implementer's record
  that quoted the old figure.
- **H5 — from claude-1 R3-NIT-2 (NIT).** Correct the G1 clause "the old undercounted sum
  was 19,736 B mislabelled `floor`": at `118b245` the variable named `floor` held
  **15,759 B** (logged `named-omission-set bytes`); **19,736 B** was the `+§2(…)
  reference` figure it was compared against (the right comparison basis for FINAL's
  ≈ 42.1 KB — the record's substantive point stands). Replacement clause (suggested):
  "the `floor` variable held 15,759 B, logged as `named-omission-set bytes`; the `+§2
  reference` it was compared against was 19,736 B."
- **H6 — from claude-1 R3-NIT-3 (NIT).** Replace the cycle-2 residual clause "no driver
  transition occurred during this cycle either — the deck's `runs/` records are still
  the historical zero-width set" (and the cycle-2 blind-spot sentence dating the newest
  record to `20260924T003301…Z`) with the corrected statement (facts DRAFTER-verified
  this session: the four recent `runs/` directories list exactly as described): "the
  newest run record is `20260924T025208.421874000Z` (02:52:08Z, `mode:
  consensus-signoff`), which postdates F7 (`64a622c`, committed 2026-09-24T01:48:45Z);
  it and the two before it (`20260923T211817Z`, `20260924T003301Z`) write no `run.json`
  at all (events.jsonl only), and the only recent manifest (`20260923T202501Z`) is
  zero-width — so no live attribution window has yet been produced." Conclusion
  unchanged and reinforced; whether a `consensus-signoff` runner launch counts as a
  "driver transition" stays arguable, but the dated facts are now correct.
- **H7 — from claude-1 R3-NIT-4 (NIT).** Correct the G4 record at
  `IMPLEMENTATION.md:112-113` ("both the fix-up-published state **and the
  newer-review-artifact relaxation**"): the test's second half asserts enumeration
  validity only — the expected value is `t.Logf`-logged, not asserted (the in-test
  comment is honest; the IMPLEMENTATION summary was not). Primary: reword the record to
  "the fix-up-published state asserted; the newer-review-artifact relaxation checked for
  enumeration validity only". Signoff-amendable stronger alternative: assert the expected
  `d2.Next` for that fixture's state, replacing the log-only check, and keep the original
  wording.

**Close conditions for this idea (R3-MIN-3(a) — existing conditions documented, nothing
new created).** `00-prompt.md:8` sets `auto_implement: true` (`strict_gate` absent;
`checks:` absent, so the list-form completion contract does not engage). Under
`auto_implement`, zero Agreed fixes is necessary but **not** sufficient:

1. **Phase-7 consensus listing zero Agreed fixes** (default Phase-8 close rule; the only
   rule kimi-1's round-02 close statement captured).
2. **LE-11 — no reservations triage:** any ACCEPT-WITH-RESERVATIONS among the signoffs
   escalates instead of completing (`internal/driver/impl.go:283-286` — "reservations
   need human review before completion"). Both cycle-2 signoffs were 🟡; the closing
   consensus therefore needs all-✅ signoffs **or** a recorded operator ruling accepting
   a 🟡 close.
3. **LE-11 — reviewer floor:** fewer than `MinReviewers` (2) independent reviewers
   escalates (`:287-289`) — satisfied by claude-1 + kimi-1.
4. **LE-7 — goal-done check:** a fresh non-implementer verifies FINAL's observable
   acceptance criteria and the check must pass (`:293-297`); a checker that is missing,
   is the implementer, or cannot be resolved and launched; a failed or non-zero
   execution; and an inconclusive or reservations-only verdict each leave completion
   **unverified** and escalate. The check can only withhold a close, never establish
   one; a textual pass never substitutes for current-tree criterion evidence. No such
   check has been commissioned or run on this idea to date (H3(b) arranges it).
5. Driver-side riders when the driver closes: completion-contract re-read/observation
   and `trajectory.RequireResolved` (`:299-309`) — moot for a manual close but recorded
   for completeness.

**This run is manually driven** under the owner-authorized manual fallback
(organizer-notes, Phase-2 fallback and the FINAL-boundary re-entry stop), so the
driver's auto-complete path will not close this idea — but the conditions the code
enforces are the protocol's own (§4 Phase 8 close-decision integrity, LE-7/LE-11), and
the manual close must satisfy and record their equivalents: zero Agreed fixes, all-✅
signoffs or a recorded operator ruling for any 🟡, two independent reviewers, and the
commissioned goal-done verdict filed before `status: complete` is set.

**Concrete dispatch steps for the organizer (H3(b)), avoiding the documented driver
re-entry gaps** (`parley run` mints timestamp slugs and cannot target the live idea;
`parley continue` resolved a newer signoff-request run at the FINAL boundary and began
drafting FINAL before being stopped; organizer-notes: "Further Phase-5/6 dispatch uses
the explicit authorized manual fallback"; "successful leaf CLI signoff requests remain
usable"):

1. **Do not re-enter `parley run`/`parley continue` for this dispatch** — the collision
   class is recorded in organizer-notes; a re-entrant driver attempting to draft or
   complete could collide with the live review-consensus state.
2. **Commission by inbox note** (organizer-owned, audit-trailed):
   `parley-deck/inbox/codex-1-to-<agent>_meta-protocol-change-lean-organizer_goal-done.md`
   naming the idea, the exact closing HEAD sha, the duty (verify FINAL.md observable
   acceptance criteria A–D + cross-cutting at that tree, own commands, §15.2 PRIMARY
   tags), the output path (`review/goal-done/<agent-id>.md`), and the prohibitions
   (fresh session; no source, reviewer, consensus, or organizer edits).
3. **Launch the checker through its configured CLI adapter directly** — the same
   hand-launch mechanics the organizer already used for round-01 and the heading
   repairs (organizer-notes Phase-0/1) — in a fresh session of claude-1 or kimi-1.
4. **Gate arrival with the leaf, read-only `parley wait`** (proven safe this cycle),
   e.g. `--for review` for the round-04 files; await the goal-done artifact in the
   organizer's own dispatch loop — never by restarting the driver.
5. **Record the outcome** in the closing consensus and organizer-notes; an
   inconclusive/failed/unavailable checker escalates to the user rather than closing.

**Closure conditions for fix-up cycle 3 (not new fixes):** both suites green at the
cycle-3 HEAD with the G3-named command (`go test ./... -count=1 -timeout 2400s`; skill
`npm test`); H1's equal-mtime regression case added (and H7's assertion if the
alternative is picked); `IMPLEMENTATION.md` gets the Phase-8 fix-up-cycle-3 section —
including its own `### Deviations from agreed fixes` (`None` or items, per the restored
discipline) — with per-fix commit references and bumped frontmatter; the core-publish
note's facts still match the staged file (no protocol-text edit expected, so no restage;
the G6 rule continues); no release, merge, tag, global install, or publish in the cycle.

**Cycle mechanics (readiness honestly recorded — see VC-7).** This consensus lists 7
proposed Agreed fixes, so it is not a closing consensus under any position in VC-7. If
the signoffs sustain the plan: cycle 3 lands H1–H7 (H1 is the only code change; H2–H7
are `IMPLEMENTATION.md` record corrections plus H3's pointer line), review round 04
verifies (fix-verification scope plus whatever fresh scope reviewers choose — it should
include the fresh-checkout equal-mtime shape), the goal-done check is commissioned at
the cycle-3 HEAD per H3(b), and consensus cycle 4 aims to be the zero-fix closing record
citing the goal-done verdict. Budget: cycle 3 of 5; trajectory per §4 stopping judgment
is the converging shape (round 1: 25 findings incl. 1 CRITICAL → round 2: 15, 3 MAJOR →
round 3: 7, 0 CRITICAL, 0 MAJOR, no unmet acceptance criterion). If the signoffs strike
items (kimi-1's READY position), the struck items' dispositions move to
`## Deferred follow-ups`/`## Dismissed findings` by the signoffs' own written decision.
One mechanical note the signoffs should weigh (drafter's reading of the packet's Phase
8, not an adjudication): claude-1's route 2 (zero-fix close with dispositions only)
still requires the `:503-508` wording softening it itself conditions its signature on,
and after `status: complete` the closed `IMPLEMENTATION.md` may not be edited — the
implementer's only edit channel for record corrections is a fix-up cycle authorized by
Agreed fixes, so a close that also corrects the records is not mechanically available;
striking all seven would close over records this idea's own discipline treats as
inaccurate. This is recorded as a reason the short third fix-up is needed, not as a
premature closure either way.

## Deferred follow-ups

- **DF-1 `meta-protocol-change-consensus-duty-gates`** — unchanged; inactive candidate,
  no quorum staffed, not launched. G7's deviation bullet points here.
- **DF-2 `meta-protocol-change-facilitator-integrity-phase-coverage`** — unchanged;
  inactive candidate (phase-5/8 §15.5/§15.6 pin question).
- **DF-3 `facilitator-packet-per-phase-bounds`** — unchanged; inactive candidate
  (phase 7 sits 86 B over the 70,000 B figure at the test path; per-phase policy is
  this slug's question, not this cycle's).
- **DF-4 `release-binary-reproducibility`** — unchanged; inactive candidate
  (embedded-build-path mechanism; kimi-1's two-party byte-identical rebuild this round
  is further evidence the mechanism, not the source, explains hash variance).
- **Advisory, no slug opened:** `internal/trajectory` / `internal/app` package
  durations — this round's measurements (claude-1: 644.7 s; kimi-1: 609.4 s; both over
  Go's 600 s per-package default on their machines) further evidence the G3 flag is
  load-bearing. Nothing launched; opening a slug stays the organizer's/owner's call.
- **No new DFs this round:** all seven round-03 findings are dispositioned in-cycle
  (H1–H7). The three carried residuals — GitHub-hosted runner wall-clock, live
  attribution windows (with H6's corrected record facts), Windows `wait`/`usage`
  portability — remain disclosed residuals inside the signed scope, not findings and
  not DFs.

## Dismissed findings

None. No finding is dismissed, withdrawn, or downgraded by this draft. This round's
withdrawals are the filer's own: claude-1 withdrew R2-MAJ-1/-2/-3 and its lesser round-02
findings as fixed (verified independently by both reviewers), narrowed R2-MIN-1(a) into
R3-MIN-1, and issued a §15.1 SELF-CORRECTION superseding its own 23,223 B figure —
filer-owned actions recorded in its review file, not draft dismissals. kimi-1 filed
nothing to dismiss; its G4(a) fresh-checkout observation is recorded (in the map and
VC-8), not dismissed. Dismissal or withdrawal can only happen through a reviewer's own
self-correction or a signoff resolution, never through this drafter.

## Coverage & blind spots

*(Advisory; signoffs remain the gate.)*

**Both reviewers independently (PRIMARY each, separate isolated worktrees and
binaries):** all ten G1–G10 verified at code and real-entrypoint level; both acceptance
elements claude-1 named unmet at round 2 now met — C.3/R-2's floor half and
cross-cutting both-suites-green at FINAL's named command (each ran the full CLI suite to
exit 0, 31/31 packages, and the skill suite to 399 pass, on their own machines); the
52,295 / 30,262 / 23,230 B figures reproduced through disjoint measurement paths
(claude-1: rendered-body surgery + index-sum under the renderer's layout convention;
kimi-1: awk per-section range sums + body surgery); staged core and core-publish note
facts re-verified (`fc907e59…`, 109,772 B, mtime 2026-09-24 03:35:27); `~/.parley/
protocol/core/` holds only `2.10.0` — no publish by anyone; kimi-1 additionally
reproduced the implementer's recorded binary byte-identically from the same clean clone
(two-party reproduction).

**Only claude-1 saw:** R3-MIN-1 — its fresh checkout carries equalized mtimes, while
kimi-1's reused review worktrees carry real mtimes (the coverage gap that hid the edge
from round 2's filer and from kimi-1 this round); R3-MIN-2 (the Phase-8 section shape);
R3-MIN-3 (the `auto_implement` close conditions — kimi-1's round-02 close rule, though
right about `strict_gate`, omitted the LE-7/LE-11 rider); the four record-archaeology
NITs.

**Only kimi-1 saw:** nothing new — zero findings; its fresh-checkout mtime observation
converges with R3-MIN-1's facts while diverging on materiality (VC-8), and its live-deck
G4(a) verification is the round's only first-party check of the live surface.

**Blind spots — what no reviewer covered:**

1. **The LE-7 goal-done check has never been run on this idea** — no artifact records one
   commissioned or executed; until H3(b)'s check lands, "close conditions satisfied" is
   untested testimony, not evidence.
2. **No shipped test exercises equal-mtime/fresh-clone surfaces** (R3-MIN-1; H1 adds the
   case) — round-03 verification used worktrees with real mtimes except claude-1's
   checkout reproduction.
3. **GitHub-hosted runner wall-clock remains unmeasured** (G3 removes the known cliff,
   not the unknown; both reviewers' local over-600 s measurements sharpen it).
4. **Live attribution windows remain test-proven only** — with H6's corrected facts: the
   newest run record postdates F7 but writes no manifest; the first live window still
   awaits a real driver transition.
5. **Windows `wait`/`usage` portability remains macOS-verified only** (carried from
   cycles 1–2; unchanged).

**Correlated agreement (§15.6(b)).** Two model families, but the cross-confirmations
rest on shared priors (PRIMARY execution over prose; deterministic tooling trustworthy)
and partially identical command sequences. Their agreement is strongest where least
independent; this round's value is again the divergences — VC-7 (readiness) and VC-8
(materiality) — which the signoff requests ask each reviewer to re-examine against the
other's evidence rather than their own.

## Verdict conflicts

*(§15.3: quoted verbatim with author, tag, and evidence. Both are **DISPUTED** and stay
open; this draft closes over neither. Resolution is by reviewer self-correction (§15.1)
or signoff — never by participant count (the reviewers split 1-vs-1), never by the
drafter (the implementer issues no self-verdict), and never by the organizer (verdicts
read only, no code verdict). The idea's closure depends on VC-7; H1 depends on VC-8.)*

**VC-7 — "Is the implementation ready for a zero-fix closing consensus?"**

- claude-1, `review/round-03/claude-1.md`, closing position — **NOT YET**: "**Readiness
  for a zero-fix closing consensus: not yet, but the gap is now dispositional rather
  than substantive.** No CRITICAL, no MAJOR, no unmet acceptance criterion, no broken
  gating property, and the trajectory is the converging shape §4's stopping judgment
  describes … Concretely, I would sign a zero-fix closing consensus that does **either**
  of the following, and I do not insist on the first: 1. lands R3-MIN-1 and R3-MIN-2 as
  a short cycle 3 and dispositions the four NITs; **or** 2. dispositions all seven by
  written signoff decision — provided R3-MIN-1's "retired" wording at
  `IMPLEMENTATION.md:503-508` is softened to match the mtime condition … Either route
  additionally requires R3-MIN-3 answered in the record." (Its round-02 NOT-READY
  rationale — two unmet frozen-FINAL acceptance elements — is withdrawn as spent: "both
  are now met and I say so plainly".)
- kimi-1, `review/round-03/kimi-1.md`, verdict — **READY**: "**Readiness for zero-fix
  closing consensus: READY.** All ten signed fixes verified at the reviewed commits,
  both suites green under my own hands, no new findings, no scope creep, FINAL frozen,
  no publish/merge/tag/install action taken by anyone per the record — and none by me.
  In my judgment the next Phase-7 consensus can be the zero-fix closing record."
- **Status: DISPUTED — and it cannot be resolved inside this draft.** The positions are
  not factually incompatible (both hold: no CRITICAL, no MAJOR, no unmet acceptance
  criterion, all ten fixes verified); they differ on whether the seven findings must be
  fixed/dispositioned before close. This draft proposes the short third fix-up (H1–H7),
  which makes the question concrete: sustaining the plan at signoff is a vote for
  claude-1's route 1; striking items is a vote for kimi-1's branch — with the mechanical
  note in **Cycle mechanics** that claude-1's own route-2 precondition requires a record
  edit whose only implementer channel is a fix-up cycle. VC-7 must be resolved by the
  signoffs, a reviewer self-correction, or an operator ruling before any closing
  consensus.

**VC-8 — materiality of G4(a)'s equal-mtime fall-through: R3-MIN-1 [MINOR, fix in-cycle]
vs kimi-1's recorded observation [not a finding, defensible].**

- claude-1, R3-MIN-1 — **MINOR, fix in-cycle**, PRIMARY (fresh checkout at `913f8ba`;
  `stat` shows `IMPLEMENTATION.md` and `review/round-02/claude-1.md` sharing
  `2026-09-24 06:04:18`; digest reads `await implementation`; `touch IMPLEMENTATION.md`
  flips it to `await review artifact`; mechanism `phasedigest.go:204` strict `After`):
  "A `git clone`, `git checkout`, worktree creation, archive extraction or any `rsync`
  without `-t` gives `IMPLEMENTATION.md` and the review artifacts **the same** mtime, so
  the predicate is false and `nextAction` falls through to `NextAwaitImplementation` —
  the exact wrong output R2-MIN-1(a) named. … G4(a) introduces a second mtime dependence
  whose failure direction is the opposite: **silent**, and wrong in the organizer's
  branch surface."
- kimi-1, `review/round-03/kimi-1.md`, observation (not a finding), PRIMARY for the
  live surface: "G4(a)'s fix-up-published signal is mtime-based by signed design ('the
  same signal wait uses'); on a fresh checkout with equalized mtimes the digest falls
  through to `await implementation`, which **remains defensible there (the implementer
  owes the next artifact)** and the digest's status line (`status=fix-up-cycle-N`) stays
  visible. The live-deck operational surface — where mtimes are real — behaves
  correctly, as verified above."
- **Status: DISPUTED (materiality only — the facts carry both reviewers' agreement:
  the equal-mtime fall-through exists on fresh checkouts; the live deck reads
  `await review artifact`).** H1 is proposed contingent on this conflict, primary shape
  the content signal (claude-1's lean: deterministic, survives any filesystem operation,
  removes a mtime dependence rather than widening one); kimi-1's defensible-as-is
  position stays open for it to sustain with a counter-proposal or concur at signoff.
  claude-1's own offered falsifier stands: a fresh-checkout reproduction that reads
  `await review artifact` withdraws R3-MIN-1.

**§15.3 dependency check.** H1 is explicitly contingent on VC-8; no other item and no
acceptance statement in this draft is derived from a disputed claim — the facts under
both conflicts (the fall-through exists on fresh checkouts; the live surface is correct;
no CRITICAL/MAJOR exists and every acceptance row passes) carry both reviewers'
independent PRIMARY provenance, and several H-item facts are DRAFTER-verified this
session (the strict `After` comparison; the `runs/` shapes; `auto_implement: true`;
FINAL frozen at the reviewed commit). Closure of the idea additionally requires VC-7
resolved by signoff, reviewer self-correction, or operator ruling, and the H3 goal-done
check commissioned and passed at the closing HEAD. claude-1's §15.1 SELF-CORRECTION
(round-02 23,223 B superseded by 23,230 B) is a weakening of its own claim, effective
immediately, and creates no conflict.

## Signoffs

*(Deliberation track: all three participants sign — reviewers and implementer. Signoffs
are intentionally empty in this draft; the cycle-2 signoffs live verbatim in
`review/consensus-cycle-02.md` and are NOT reused. Append-only; do not edit others'
blocks. Each signoff request asks the signer to: (1) state a binding position on **VC-7**
(short cycle 3 vs close-now) and **VC-8** (content signal vs `>=` relaxation vs
defensible-as-is); (2) confirm or amend the H1–H7 plan and the map — every filed finding
must be traceable to a disposition; (3) kimi-1: reproduce claude-1's fresh-checkout
equal-mtime reproduction in its own fresh clone/checkout (its reused worktrees carry
real mtimes — the round's one coverage gap; claude-1's offered falsifier: an
`await review artifact` reading there withdraws R3-MIN-1) and self-correct if warranted
(§15.1); (4) claude-1: evaluate H3(b)'s dispatch plan against its open question 2 and
H1's shape against its open question 3; (5) both: note that the signoffs on THIS
consensus do not trigger LE-11 (the driver is not closing this manually-driven run),
but the closing consensus's signoffs interact with it per H3 — a 🟡 there requires a
recorded operator ruling to close. The implementer signs per §15.1 with no
self-verification verdict on its own implementation.)*

### Signoff: claude-1 — (pending)

### Signoff: kimi-1 — (pending)

### Signoff: zcode-1 — (pending)

### Signoff: claude-1 — 2026-09-24
Status: ✅ ACCEPT
Notes:

**What this ✅ accepts.** The H1–H7 plan, the complete finding→disposition map, the
dispositions, and the decision to leave VC-7 and VC-8 open. It is **not** a close: this
consensus lists seven Agreed fixes, so nothing here certifies completion, and my VC-7
position is unchanged by signing it. Two amendments to the draft's own record are below;
neither conditions my acceptance.

**VC-7 — binding position: NOT YET, and this draft is route 1, the route I named first.**
I sustain the round-03 position verbatim: no CRITICAL, no MAJOR, no unmet acceptance
criterion, no broken gating property; the gap is dispositional, not substantive. I
sustain the plan as drafted — cycle 3 lands H1–H7, review round 04 verifies, and cycle 4
aims to be the zero-fix closing record. **I concede the drafter's mechanical point about
my route 2** (Cycle mechanics): route 2 required softening `IMPLEMENTATION.md:503-508`,
that softening is itself a record edit, the implementer's only edit channel for record
edits is a fix-up cycle authorized by Agreed fixes, and a closed `IMPLEMENTATION.md` may
not be edited after `status: complete`. So pure route 2 — disposition all seven and change
nothing — is not available to me, and I withdraw it as a live option rather than leave it
on the record as one. @kimi-1: this is the concrete reason I do not think READY and my
NOT-YET are as far apart as the labels suggest; striking all seven closes over records
this idea has corrected in every prior cycle, and the correction channel closes with the
idea. **Forward commitment, per §4 stopping judgment:** the trajectory is converging (25
findings incl. 1 CRITICAL → 15 incl. 3 MAJOR → 7 incl. 0 CRITICAL/0 MAJOR) and this is
cycle 3 of 5. If review round 04 produced another comparable crop of record findings on
unchanged ground, I would read that as churn, not convergence, and would support closing
over dispositions rather than opening a cycle 5.

**VC-8 — binding position: R3-MIN-1 sustained at MINOR, fix in cycle 3 (H1).**
Re-reproduced at signoff time by a **second, disjoint path** — `git archive HEAD | tar -x`
into `/tmp/c1-signoff/repro` (archive extraction is one of the operations the finding
names), binary `/tmp/c1-signoff/parley` built by me at `913f8ba`, sha256 `f87bab4f…f56`
(PRIMARY):

    $ stat -f '%Sm  %N' IMPLEMENTATION.md review/round-02/claude-1.md review/round-02/kimi-1.md
    Sep 24 06:01:35 2026  IMPLEMENTATION.md
    Sep 24 06:01:35 2026  review/round-02/claude-1.md      # equal
    Sep 24 06:01:35 2026  review/round-02/kimi-1.md        # equal
    # digest: implementation: present=true status=fix-up-cycle-2 implementer=zcode-1
    #         review round-02 filed 2/2 valid; no round-03 dir
    ## Next action (fixed enumeration)
    await implementation                 # WRONG
    $ touch IMPLEMENTATION.md
    await review artifact                # correct

@kimi-1 — one distinction is, I think, the whole of our disagreement, and I may have
stated it too thinly in round 03. Your "defensible" was **established on your control
fixture**, where the review artifact is *newer* than the implementation: there a review
round arrived after the implementation and the implementer genuinely does owe the next
artifact, so `await implementation` is right. The equal-mtime case is not that control.
In it the digest's own status line says `fix-up-cycle-2` — the implementer has published
and signalled ready-for-re-review, and the artifact owed is round 03, by us. That is the
exact state G4(a) was signed to repair, and it is the state every fresh clone, checkout,
worktree creation, archive extraction or `rsync` without `-t` produces. Your live-deck
verification is correct and I do not contest it; the live deck reads correctly through the
older directory-based path (empty `round-03/`), which is state-dependent, not a guarantee.
**My falsifier stands unchanged:** a fresh clone/checkout/archive in that state reading
`await review artifact` withdraws R3-MIN-1 — I have now run it twice and not obtained it.

**H1 shape (my open question 3) — content signal, as drafted, and it needs no new
machinery.** Verified this session (PRIMARY): `PhaseImplSection.Status` carries
`fix-up-cycle-N` verbatim (closed vocabulary, `internal/protocol/reviewartifact.go:103`
`^fix-up-cycle-\d+$`), `PhaseRoundSection` carries `Label`/`Completed`/`Total`, and
`roundNumberFromLabel` (`phasedigest.go:333`) already parses `round-0M` → M. So `N` and
`M` are both in hand and the predicate is a comparison, not a new parser. I therefore
endorse the primary shape over the `>=` relaxation: it removes the mtime dependence rather
than widening it, and it is deterministic under any filesystem operation. If the quorum
picks the relaxation instead, the equal-mtime regression case and the `:503-508` softening
are still mandatory from my side — those are the part I will not trade.

**Amendment 1 (locator; PRIMARY at `913f8ba`).** The draft cites `phasedigest.go:204` for
the strict `After` twice (finding map and H1). At the reviewed commit, `:204` is a comment
line inside `nextAction`'s doc block. The correct locators: the function
`implementationNewerThanLatestReview` spans **`:182-200`**, and
`return implInfo.ModTime().After(newest)` is at **`:199`**. The mechanism is unambiguous
because the draft quotes the expression, so no claim changes — but cycle 3's record should
carry `:199`.

**Amendment 2 (H6 — the phrasing rots, and it already has).** H6's replacement text names
`20260924T025208.421874000Z` (02:52:08Z) as "the newest run record". It was newest when
the draft was written (consensus.md mtime 04:40:16Z) and stopped being so ~72 seconds
later: `parley-deck/runs/20260924T044128.336773000Z/` exists — created
`2026-09-24T04:41:28.336793Z`, `mode: consensus-signoff`, a single `run.created` event
(202 B), **no `run.json`** — and is, by its timing and participant list, almost certainly
the run carrying this very signoff round (PRIMARY: `ls`, `head events.jsonl`, `wc -c`).
Neither the drafter nor I erred; the *formulation* does, because each signoff-request
launch mints another record, including the launch that files the correction. Suggested
durable replacement for H6, which I ask cycle 3 to use instead: *"Every recent `runs/`
record for this idea writes `events.jsonl` only, with no `run.json` at all, except
`20260923T202501Z`, whose manifest has `created_at == updated_at` (verified: both
`2026-09-23T20:25:01.377412Z`) — a zero-width window. Records postdating F7 (`64a622c`,
2026-09-24T01:48:45Z) do exist, including `mode: consensus-signoff` runs at 02:52:08Z and
04:41:28Z, so the cycle-2 claim that no post-F7 record exists is wrong; none has produced
a non-zero-width attribution window, so the conclusion is unchanged and reinforced: live
attribution windows remain test-proven only."* This is falsified by a real driver
transition, not by the next signoff launch.

**H3(b) (my open question 2) — accepted as drafted, with one thing checked so it cannot
bite.** The dispatch plan answers the question I asked: the check is commissioned, by
inbox note, at the **closing HEAD**, from a fresh non-implementer non-facilitator session,
filing its own artifact and editing nothing, with an inconclusive/failed/unavailable
checker escalating rather than closing. I verified the proposed path is inert to the
surfaces this idea is about: `latestRoundSection` filters candidate directories on
`strings.HasPrefix(e.Name(), "round-")` (`phasedigest.go:305`), so a `review/goal-done/`
directory is not enumerated as a review round and cannot perturb the digest, the brief, or
`wait --for review` (PRIMARY, code read at `913f8ba`). I accept the path as drafted. I am
eligible for the role and willing to take it if the organizer dispatches it to me; I hold
no preference between me and kimi-1, and whoever takes it must run it at the post-cycle-3
tree, not at `913f8ba`.

**Open question 1 — answered, and H2 matches my lean.** The agreed-fixes deviation bullet
is the right and sufficient home for the ≈ 42 KB → 52,295 B divergence; `## Deviations
from FINAL.md` stays for mechanism changes like G7's. I do not re-litigate branch 1: the
recomputed whole-block total is the conceptually correct realization and the criterion is
met as frozen. H2's premise re-verified this session: `## Fix-up cycle 2` spans
`IMPLEMENTATION.md:14-289` and contains no `### Deviations from agreed fixes`; the only
one is at `:494`, inside cycle 1 (PRIMARY).

**§15.1 SELF-CORRECTION (my own locator).** My round-03 wrote
"`implementationNewerThanLatestReview` (`internal/driver/phasedigest.go:187-205`)". The
function spans `:182-200`; my range both started late and ran past the function's end.
Superseded by `:182-200` / comparison at `:199`. A locator correction with no effect on
R3-MIN-1 or on any other finding; effective immediately.

**Map completeness and no suppression.** All seven of my findings are traceable to a
disposition (R3-MIN-1→H1, R3-MIN-2→H2, R3-MIN-3→H3, R3-NIT-1→H4, R3-NIT-2→H5,
R3-NIT-3→H6, R3-NIT-4→H7), severity and provenance preserved; kimi-1's zero-findings
position and its G4(a) observation are recorded as its own, not converted into a fix.
Nothing was narrowed, downgraded or dismissed by the draft, and no brief told me what not
to report. H3's close conditions match the code I read independently
(`internal/driver/impl.go`: `:283-286` reservations triage, `:287-289` reviewer floor,
goal-done gate at `:294-298` — the draft's `:293-297` overlaps it and my own round-03 cite
was the same, so I raise no item, only the exact span for cycle 3) and `00-prompt.md:8` is
`auto_implement: true` (PRIMARY).

**On the closing consensus and LE-11.** Noted: this signoff does not route through the
driver's close gate. For planning, not as a pre-commitment — my cycle-2 🟡 rested on two
MAJOR-based unmet acceptance elements, both discharged and withdrawn; nothing currently
outstanding on my side would force a 🟡 from me at close if H1–H7 land as planned and the
goal-done check passes at the closing HEAD. I will judge the closing consensus on its own
artifact.

**Provenance for this signoff.** All PRIMARY checks above are mine, run this session in
the live deck worktree at `913f8ba` and in the disposable `/tmp/c1-signoff/repro` archive
extraction; binary `/tmp/c1-signoff/parley` (sha256 `f87bab4f0339fdfa75b07d434f54800e845d1789f0898c9a428720fdfe257d56`),
built by me with `go build -o /tmp/c1-signoff/parley ./cmd/parley`. I edited exactly one
file — this one — appending exactly this block; no existing line was altered. No source,
FINAL, IMPLEMENTATION, peer-review, organizer or inbox artifact was modified; nothing was
committed, merged, tagged, installed, released or published, and `parley protocol publish`
was not run.

### Signoff: kimi-1 — 2026-09-24
Status: ✅ ACCEPT
Notes:

**What this ✅ accepts.** The H1–H7 plan, the complete finding→disposition map, the close
conditions record (R3-MIN-3(a)), and the decision to leave VC-7/VC-8 open for the signoffs.
This is not a close — seven Agreed fixes are listed — and my cycle-2 🟡 reservations (CI
wall-clock, live attribution windows, Windows portability) remain carried disclosures in
this draft's blind spots, not waived; my closing-consensus signoff will be judged on its
own artifact. Two of my positions change below; both changes are stated plainly with their
evidence.

**§15.1 SELF-CORRECTION (my round-03 G4(a) observation) — effective immediately.** The
signoff request asked me to reproduce claude-1's fresh-checkout reproduction in my own
fresh extraction; I did, and the result supersedes my round-03 observation. PRIMARY, my own
hands this session: `git archive 913f8ba | tar -x` into `/tmp/kimi1-signoff/fresh` (an
operation R3-MIN-1 names), run with MY round-03 binary `/tmp/kimi1-r3/parley` (sha256
`79dbf7bf…`, built by me at `998346c`, code-identical to `913f8ba`):

    $ stat -f '%Sm  %N' IMPLEMENTATION.md review/round-02/claude-1.md review/round-02/kimi-1.md
    Sep 24 06:01:35 2026  IMPLEMENTATION.md
    Sep 24 06:01:35 2026  review/round-02/claude-1.md      # equal
    Sep 24 06:01:35 2026  review/round-02/kimi-1.md        # equal
    # digest: implementation: present=true status=fix-up-cycle-2 implementer=zcode-1
    #         review round-02 filed 2/2 valid; no round-03 dir in the committed tree
    ## Next action (fixed enumeration)
    await implementation                 # WRONG
    $ touch IMPLEMENTATION.md
    await review artifact                # correct — strict After is the sole trigger

claude-1's falsifier did not fire; R3-MIN-1 stands, and claude-1's signoff distinction is
the correct reading of my error: my round-03 "defensible" was established on my control
fixture where the review artifact is *newer* than the implementation — there `await
implementation` is right because the implementer genuinely owes the next artifact. The
equal-mtime fresh-checkout state is not that control: the digest's own status line says
`fix-up-cycle-2` — the implementer has published and signalled ready-for-re-review, and the
artifacts owed are the round-03 reviews, by the reviewers. That is the exact state G4(a)
was signed to repair, and every fresh clone/checkout/archive/`rsync` without `-t` produces
it. My parenthetical "the implementer owes the next artifact" is therefore wrong for that
state and is superseded; the rest of the observation (the fall-through exists; the live
deck read correctly via the directory path at round-03 time) stands and converges with
claude-1's facts, as the draft records.

**VC-8 — binding position: R3-MIN-1 sustained at MINOR; fix in cycle 3 via H1, content
signal as the primary shape.** Given the reproduction above, defensible-as-is is no longer
a position I hold. Between the two fix shapes I endorse the content signal over the `>=`
relaxation, and I verified this session that it needs no new machinery (PRIMARY, code read
at `913f8ba`): `roundNumberFromLabel` already parses `round-0M` → M
(`phasedigest.go:333-339`), `implementationNewerThanLatestReview` spans `:182-200` with the
strict `After` at **`:199`** — claude-1's Amendment 1 is correct, the draft's `:204` sits
past the function's closing brace — and the digest demonstrably carries
`status=fix-up-cycle-N`. The content signal removes a filesystem-fragile dependence instead
of widening it (coarse timestamp granularity and `rsync -t` ties keep the relaxation
breakable); the relaxation is an acceptable fallback only with both mandatory companions —
the equal-mtime regression case and the `IMPLEMENTATION.md:503-508` correction. That
correction is mandatory under either shape: I re-read `:503-508` this session (PRIMARY) and
the unconditional "Fix-up G4(a) retired the residual" claim is falsified by the
reproduction above — with no next-round directory and equal mtimes the enumeration reads
`await implementation` today. For H7 I mildly prefer the signoff-amendable stronger
alternative (assert the expected `d2.Next`), which pairs naturally with the content-signal
shape; the choice is the implementer's.

**VC-7 — binding position: NOT YET; I sustain the short cycle 3 (claude-1's route 1), and
I record this as a position change, not a reversal under pressure.** My round-03 READY was
filed independently, before reading claude-1's round-03 file; re-examined against its
evidence as the signoff request asks, three points move me. (1) The findings are real: I
reproduced the load-bearing one myself (above), and I verified R3-MIN-2's premise this
session (PRIMARY: `## Fix-up cycle 2` spans `IMPLEMENTATION.md:14-289`; the only
`### Deviations from agreed fixes` is at `:494`, inside cycle 1) and R3-MIN-3's premise
(`00-prompt.md:8` is `auto_implement: true`; the LE-7/LE-11 close conditions the draft
records are the ones I omitted from my round-02 close statement — its rule captured only
the default). (2) claude-1's mechanical point closes route 2: the record corrections have
no channel after `status: complete`, so a close-now would freeze factual errors — the
"retired" claim, the mislabelled `floor` figure, the stale `runs/` facts — into a sealed
record, in the one idea whose own discipline has corrected exactly these errors every
cycle. (3) The trajectory supports exactly one short cycle: 25 findings (1 CRITICAL) → 15
(3 MAJOR) → 7 (0 CRITICAL/0 MAJOR, no unmet acceptance criterion). I align with claude-1's
forward commitment: a comparable crop of fresh findings on unchanged ground at round 04
would read as churn, not convergence, and I would then support closing over dispositions
rather than opening cycle 5. Six of seven items are record corrections; H1 is one code
change on the already-signed G4 surface inside frozen FINAL B. This stays well inside
budget (cycle 3 of 5).

**H3(b) dispatch plan — accepted; path verified inert.** PRIMARY: `latestRoundSection`
filters candidate directories on `strings.HasPrefix(e.Name(), "round-")`
(`phasedigest.go:305`), so `review/goal-done/` is not enumerated as a review round and
cannot perturb the digest, the brief, or `wait --for review`. The plan — inbox-commissioned,
fresh non-implementer non-facilitator session, at the post-cycle-3 closing HEAD, filing its
own artifact and editing nothing, inconclusive/failed/unavailable escalating — is the right
shape. I am eligible and willing to take the check if the organizer dispatches it to me;
whoever takes it must run it at the cycle-3 tree, not at `913f8ba`.

**Amendment 2 (H6 durable rewording) — endorsed.** The named-newest-record phrasing rots at
each signoff-request launch; the proposed durable replacement (facts about record *shapes*:
`events.jsonl`-only runs, the one zero-width manifest, post-F7 records existing, no non-zero
attribution window yet) is falsified only by a real driver transition. Use it in cycle 3.

**Map completeness, no suppression, LE-11 noted.** All seven claude-1 findings trace to
H1–H7 with severity and provenance preserved; my zero-findings position and my G4(a)
observation are recorded as mine, not converted or dismissed; no brief narrowed anything I
reported. Noted per the signoff request: this consensus's signoffs do not route through the
driver's LE-11 close gate (manually driven run), and the closing consensus's signoffs do —
a 🟡 there requires a recorded operator ruling. One honest live-state observation, not a
finding and no gate: at this moment the live deck reads `await implementation` because
round-03 is now 2/2 with real, newer mtimes — the fall-through, but directionally tolerable
here (the pending artifacts are signoffs, a class the digest's consensus line surfaces
separately), and H1's content signal falls through identically in this state (M=3 > N=2),
so the plan neither causes nor worsens it.

**Provenance for this signoff.** All checks above are mine, run this session: fresh
`git archive` extraction of `913f8ba` at `/tmp/kimi1-signoff/fresh` with equal mtimes
confirmed by `stat`; digest via my own round-03 binary `/tmp/kimi1-r3/parley` (sha256
`79dbf7bf153fb69d6f37b10d4a4808c16c3d901625e2ddeffa8a2d11daa2c0d0`, built by me from my
clean clone at `998346c`, code-identical to the reviewed commit); `touch` control; live-deck
brief before/after; `IMPLEMENTATION.md` structure grep and `:500-510` read; `phasedigest.go`
`:182-200`, `:303-307`, `:330-340` reads at `913f8ba`. I edited exactly one file — this one —
appending exactly this block; no existing line was altered, and the `(pending)` placeholders
above are the drafter's, left untouched. No source, FINAL, IMPLEMENTATION, peer-review,
organizer or inbox artifact was modified; nothing was committed, merged, tagged, installed,
released or published, and `parley protocol publish` was not run.

### Signoff: zcode-1 — 2026-09-24
Status: ✅ ACCEPT
Notes:

**What this ✅ accepts.** As drafter of this record and the implementer, I accept the
H1–H7 plan, the complete finding→disposition map, the close-conditions record (R3-MIN-3(a)),
and the decision to leave VC-7/VC-8 to the signoffs. The signoffs have now resolved both:
kimi-1's §15.1 SELF-CORRECTION removes the minority position in VC-8 (its fresh-`git
archive` reproduction joins claude-1's two, so R3-MIN-1 stands with three independent
PRIMARY reproductions and the falsifier unfired), and its recorded position change on VC-7
leaves both reviewers at NOT-YET sustaining route 1 — resolved by signoff and
self-correction, exactly the channels this draft reserved; nothing was closed by count,
drafter, or organizer. **This is not a close**: the consensus lists seven Agreed fixes, so
it authorizes fix-up cycle 3 and nothing more; completion remains the zero-fix closing
consensus's call, gated per H3.

**§15.1 — implementer disclosure.** I issue no verification verdict on my own
implementation; every code-level verdict above belongs to claude-1 and kimi-1 in their own
files and signoff blocks. The facts I checked this session are record facts I will act on
in cycle 3, not adjudications of my implementation's correctness: at HEAD `913f8ba`,
`implementationNewerThanLatestReview` spans `phasedigest.go:182-200` with `return
implInfo.ModTime().After(newest)` at `:199`, and `:204` is a comment line inside
`nextAction`'s doc block — claude-1's Amendment 1 is correct and the draft's `:204`
citation was imprecise (the quoted expression kept the mechanism unambiguous); and
`parley-deck/runs/20260924T044128.336773000Z/` exists with `events.jsonl` only (no
`run.json`), first event `run.created` for this idea, `mode: consensus-signoff`,
participants claude-1/kimi-1/zcode-1 — Amendment 2's triggering fact, now also carrying two
`agent.usage` events appended after claude-1's check, which is the phrasing-rot argument
demonstrating itself.

**Amendments accepted into the cycle-3 plan.** Amendment 1: cycle 3's record cites
`:182-200` / comparison at `:199`. Amendment 2: H6 uses the durable shape-based rewording
both reviewers endorsed (facts about record shapes and the zero-width manifest, not a
named-newest record). Both are commitments of the fix-up cycle, not conditions of this
acceptance.

**Binding positions on the plan items.** VC-7: the reviewers own readiness of my
implementation; on the route question I record only a process position — sustain the short
cycle 3, with the mechanical point (record corrections have no channel after `status:
complete`) now conceded by claude-1's route-2 withdrawal and kimi-1's concurrence. H1 goes
forward on its primary shape, the content signal (both endorsers; claude-1's falsifier
checks and the mandatory equal-mtime fresh-checkout regression test and the
`IMPLEMENTATION.md:503-508` correction included). H7 — taking the choice kimi-1 left me:
the stronger alternative, asserting the expected `d2.Next` for the relaxation fixture
instead of the log-only check, keeping the original record wording; claude-1 offered this
as the signoff-amendable alternative and kimi-1 mildly prefers it, so the quorum
accommodates it and it pairs with H1's test additions. H3(b): the goal-done check is not
mine to run — it is commissioned from a fresh non-implementer session at the closing HEAD;
nothing in cycle 3 may claim the close conditions satisfied before that verdict lands, and
I note both reviewers independently verified `review/goal-done/` inert to the digest via
the `HasPrefix("round-")` filter (`phasedigest.go:305`).

**Drafter position changes: none.** Verified against my archived cycle-2 signoff block
this session: same stance — accept the record and plan, issue no self-verdict, execute the
authorized fix-up cycle, claim no completion.

**Provenance for this signoff.** One file edited (this one), one block appended; no
existing line altered and the `(pending)` placeholders left untouched, matching both
co-signers. No source, FINAL, IMPLEMENTATION, peer-review, organizer, or inbox artifact
modified; nothing committed, merged, tagged, installed, released, or published;
`parley protocol publish` not run.
