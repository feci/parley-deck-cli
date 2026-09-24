---
idea: meta-protocol-change-lean-organizer
review-cycle: 4
outstanding_agreed_fixes: 0
blocked: false
drafted-by: zcode-1
date: 2026-09-24
reviewed-commit: 4df0855
---

<!-- outstanding_agreed_fixes: 0 and blocked: false are the DRAFTER'S PROPOSAL, not a
     resolution. Both reviewers filed READY positions this round, but kimi-1 has not
     yet stated a position on the three dispositions below (it filed zero findings and
     did not address claude-1's R4 items in its review file) — its concurrence is
     requested at signoff, never presumed or majority-resolved. Signoffs are empty in
     this draft. -->

Draft of the Phase-7 review consensus for review cycle 4 — proposed as the **zero-fix
closing record** (after fix-up cycle 3; reviewed commit `4df0855`, fix-up-3 source
`c3baf09`). The previous signed consensus is archived **VERBATIM** at
`review/consensus-cycle-03.md` (58,715 B, `cmp`-silent against the file it replaces,
sha256 `419900f1…b9f116b`, including all three ✅ signoff blocks); this fresh draft reuses
no stale signoff. Inputs read in full this session: the live phase-7 packet attested
below; frozen `FINAL.md` (frozen at `120a9bf`; `git diff 120a9bf..HEAD -- …/FINAL.md`
empty — re-verified this session); `IMPLEMENTATION.md` at `4df0855` including all of
fix-up cycle 3; both complete round-04 review files (`review/round-04/claude-1.md`,
`review/round-04/kimi-1.md`); the entire signed cycle-3 consensus (now
`review/consensus-cycle-03.md`); the commissioned LE-7 goal-done verdict
`review/goal-done/kimi-1.md` (**PASS** at `4df0855` / source `c3baf09` / skill `b06a65a`);
`review/consensus-cycle-02.md`; `review/consensus-cycle-01.md`; `00-prompt.md`;
`organizer-notes.md`. No code, FINAL, release-plan, reviewer-file, goal-done-file,
organizer-note, or inbox edits were made in this invocation; this file (plus the verbatim
archive of its predecessor) is the only output. Nothing is committed, tagged, installed,
published, or released; `parley protocol publish` is not run; `IMPLEMENTATION.md` is NOT
touched — `status: complete` is set only after all-✅ signoffs (Phase 8 close), never by
this draft.

**Protocol context attestation (Phase-7 drafting):**

```json
{"context_mode": "full", "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "fallback_reason": null, "body_path": "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer/.parley-runtime/protocol-packets/full-phase7-deliberation-8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7.md"}
```

Drafter cross-check (PRIMARY): `shasum -a 256` over the packet body returns
`8ce83cde…a9db7` — the attested packet IS the live authority; both round-04 reviews and
the goal-done check attest the same hash against their own trees.

**Drafter disclosure (§15.1).** The drafter is the implementer (zcode-1), drafting per
Phase 7's default. The implementer issues **no verification verdict on its own
implementation**: the zero-Agreed-fixes proposal and every disposition below follow the
**filer's own filed positions** (claude-1's review states each suggested disposition and
conditions its signoff only on their being carried into the record); the reviewers'
verification evidence is theirs (PRIMARY in their files), not re-asserted by the drafter.
kimi-1's concurrence on the three dispositions is **not yet on record** and is requested
at signoff — it is not presumed from its READY position, and nothing here is resolved by
participant count, drafter self-judgment, or organizer action. No brief suppressed or
narrowed any finding: all three claude-1 findings and kimi-1's zero-findings position
(with its two checked-and-resolved observations) carry an explicit disposition below, and
the raw review files remain canonical. The drafter's positions are unchanged since its
cycle-3 signoff block (archived verbatim); no drafter position changes to record.

Drafter-verified record facts this session (tagged DRAFTER-PRIMARY inline): packet hash
equality above; HEAD `4df0855` with `git show --stat` touching only `IMPLEMENTATION.md`;
`c3baf09` touching only `internal/driver/phasedigest.go` + `internal/app/wait_test.go`;
FINAL frozen (empty diff `120a9bf..HEAD`); the archive's byte identity; `~/.parley/
protocol/core/` holding only `2.10.0` (no publish by anyone); the archived cycle-3
signoff blocks (three ✅, H1–H7 authorized).

## Agreed fixes

**None — `outstanding_agreed_fixes: 0`.** Under the Phase-8 default close rule
(`strict_gate` absent; `checks:` absent, so the LE-4 list-form completion contract does
not engage), a Phase-7 consensus listing zero Agreed fixes is the close condition. This
section is deliberately empty; the three round-04 findings are dispositioned below as
two deferred follow-ups and one dismissed finding, exactly as their filer proposed
("none of them asks for a fix-up cycle 4"; "I do not ask for a cycle to fix a typo").
If any signer amends a disposition into a fix at signoff, this stops being the closing
record: a fix-up cycle 4 opens, and the goal-done verdict must be **re-commissioned at
the new HEAD** (claude-1 open question 1 — LE-7 rejects a stale code tree).

### Finding → disposition map (complete)

All 3 findings filed this round (all claude-1; kimi-1 files none). Prior-round outcomes
for traceability: all 7 claude-1 round-03 findings **withdrawn-as-fixed by their filer**
(R3-MIN-1 discharged on the filer's own pre-committed falsifier, fired at `4df0855` in
the filer's hands and independently in kimi-1's archive-extraction reproduction);
kimi-1's round-03 count was zero.

| Finding (severity, provenance) | Disposition |
|---|---|
| claude-1 R4-MIN-1 (`wait --for review` boundary is `latest existing round complete`, with no notion of which fix-up cycle that round reviewed — post-H1 the `--json` envelope's `next: await review artifact` and the terminal `boundary reached`/exit 0 contradict each other on a state this deck reaches every cycle; not an H1 regression, violates no frozen criterion — FINAL B.3 specifies the command and exit code without pinning the boundary to the published cycle; `boundaryReached` is referenced only by `wait` itself, so the driver's phase machine is unaffected and FINAL B.5's pure-observer property holds) [MINOR, PRIMARY — fresh extraction at `4df0855`, quoted envelope; control showing correct behavior once `round-04/` exists; repo-wide grep of the reference] | **Deferred follow-up** — TBD slug `wait-boundary-vs-published-fixup` (filer's own suggestion; candidate per LE-10). Filer does not condition its closing signoff on either offered fix shape (content signal in `boundaryReached`, or documenting the directory-creation convention in the skill's `wait` section). |
| claude-1 R4-NIT-1 (a quoted `status:` is ready-for-review but unparseable — `implSection` trims whitespace while `ValidImplementationStatus` also strips `"`/`'`, so `"fix-up-cycle-2"` defeats H1's content signal into R3-MIN-1's narrower failure and a quoted `"complete"` misses the exact-match arm; zero observed instances in this repo's 80-plus idea corpus; quote tolerance predates this idea at `815c93a`) [NIT, PRIMARY — byte-identical-to-shipped-test probe with quotes added, both arms reproduced] | **Deferred follow-up** — same TBD slug as R4-MIN-1; filer's suggested one-line root fix (`strings.Trim(strings.TrimSpace(…), "\"'")` in `implSection`) fixes both arms. Not worth a fix-up cycle on its own (filer's words). |
| claude-1 R4-NIT-2 (typo at `IMPLEMENTATION.md:453` — "cycle 1 carries **its** at the equivalent position" should read "one"/"its own"; underlying fact correct) [NIT, PRIMARY — locator quoted, cycle-1 `:815` verified] | **Dismissed as immaterial** — the filer's own requested disposition, with eyes open: closing freezes the typo into the record, which the filer explicitly consents to ("dismiss as immaterial and close over it"). |
| kimi-1 — no findings (0/0/0/0). Two checked-and-resolved observations recorded in its review (design-consensus `triage=reserved` line is accurate history, not a mis-parse; the brief's advisory run-phase pointer is the raw run cursor, with the attested packet phase correctly reading `phase8`) | No disposition needed — both resolved by kimi-1's own code reads (`phasedigest.go:142`, `organizer.go:200-217`); recorded here for coverage, dismissed by no one. |

**Organizer-owned corrections: none required this round.** No finding names an
organizer-owned artifact.

**Open questions carried, not dispositioned** (observations, not findings — no action
proposed, no scope added): claude-1's open question 2 (`head-commit:` frontmatter carries
an explanatory parenthetical where the phase-5/8 template shows a bare sha; no code
parses the field; the filer records it as a judgment call, not a defect) is left as-is
unless a signer prefers the template shape — amending it would itself require a fix-up
cycle, contradicting every filed position this round.

## Deferred follow-ups

- **DF-5 (new this round) `wait-boundary-vs-published-fixup` — TBD, inactive candidate,
  not launched.** Carries R4-MIN-1 + R4-NIT-1 (both filer-deferred, table above).
  **Suggested inactive candidate path:** `parley-deck/ideas/wait-boundary-vs-published-fixup/00-prompt.md`
  with `status: candidate` per LE-10 — the organizer may create this record before
  signoffs so the deferred items link a real slug instead of `TBD`; creating it staffs
  no quorum and launches nothing (no `participants:` set, no round-01, status stays
  `candidate` until a human flips it). This draft does not create it — organizer-owned
  per the division both reviewers recorded.
- **DF-1 `meta-protocol-change-consensus-duty-gates`** — unchanged; inactive candidate,
  no quorum staffed, not launched. Cycle-2's G7 deviation bullet points here.
- **DF-2 `meta-protocol-change-facilitator-integrity-phase-coverage`** — unchanged;
  inactive candidate (phase-5/8 §15.5/§15.6 pin question).
- **DF-3 `facilitator-packet-per-phase-bounds`** — unchanged; inactive candidate (phase 7
  sits 86 B over the 70,000 B figure at the test path; per-phase policy is this slug's
  question, not this cycle's).
- **DF-4 `release-binary-reproducibility`** — unchanged; inactive candidate
  (embedded-build-path mechanism; kimi-1's second consecutive two-party byte-identical
  rebuild this round is further evidence the mechanism, not the source, explains hash
  variance). Post-deploy channel audit territory — organizer/owner work after Phase 8,
  not this record's.
- **Advisory, no slug opened:** `internal/trajectory` / `internal/app` package durations
  — this round's measurements (claude-1: suite durations over Go's 600 s per-package
  default on its machine; kimi-1: 626.9 s in `internal/trajectory`) further evidence the
  G3 `-timeout 45m` flag is load-bearing. Nothing launched.
- **No other new DFs.** The three carried residuals — GitHub-hosted runner wall-clock,
  live attribution windows (test-proven only; recent `runs/` records events-only or
  zero-width), Windows `wait`/`usage` portability (macOS-verified only) — remain
  disclosed residuals inside the signed scope, restated by claude-1's round-04 closing
  position and the goal-done check's own residuals section; they are not findings and
  not DFs.

## Dismissed findings

- **claude-1 R4-NIT-2 (typo at `IMPLEMENTATION.md:453`) — dismissed as immaterial.**
  Rationale: purely textual, zero behavioral surface, filer-verified correct underlying
  fact; the filer itself filed the dismissal request with explicit consent to freezing
  the typo into the closed record. Recorded knowingly per the filer's own framing ("a
  closing consensus should decide knowingly whether to freeze it rather than discover it
  later"). **This dismissal, like the two deferrals, needs the signoffs' concurrence to
  stand** — a reviewer who disagrees amends at signoff, which reopens a fix-up cycle
  rather than closing over a dispute.

Beyond this item, no other finding is dismissed, withdrawn, or downgraded by this draft.
The round-03 withdrawals are the filer's own (all seven, on its own evidence); kimi-1
filed nothing to dismiss. Dismissal or withdrawal happens only through a reviewer's own
self-correction or a signoff resolution, never through this drafter.

## Coverage & blind spots

*(Advisory; signoffs remain the gate.)*

**Both reviewers independently (PRIMARY each, separate isolated worktrees and own
binaries):** all seven H1–H7 fixes verified at the reviewed commits — H1 held under
four separate refutation attacks (claude-1) and the real equal-mtime fresh-checkout
behavior plus its regression test (kimi-1, including the pre-H1 binary reading `await
implementation` on the identical tree where the current binary reads `await review
artifact`); H2–H7 each held with re-measurement or fact-by-fact re-verification; both
Amendments landed; the frozen FINAL A–D + cross-cutting regression sweep re-attacked at
this HEAD with every row passing on their own measurements; both suites green under
their own hands (CLI 31/31 packages exit 0; skill 399 pass at unchanged `b06a65a`);
FINAL frozen; staged-core facts re-verified (`fc907e59…`, 109,772 B); no
publish/merge/tag/release/install action by anyone — and none by either reviewer;
kimi-1's second consecutive two-party byte-identical binary reproduction
(`19831a24…06bdc`, `vcs.revision c3baf09`, `vcs.modified=false`).

**Only claude-1 saw:** R4-MIN-1 (the `wait` boundary semantics — kimi-1's round-04
verified H1's digest behavior and the goal-done criteria, not `wait`'s boundary
definition); R4-NIT-1 (the quoted-status probe); R4-NIT-2 (the typo). All three are
claude-1-primary with reproduced evidence, and none rests on implementer assertions.

**Only kimi-1 saw:** nothing new — zero findings; its two checked-and-resolved
observations (design-consensus triage line; brief run-phase pointer) are the round's
only first-party resolutions of those surfaces, both by code read.

**Blind spots — what no reviewer covered:**

1. **`wait --for review`'s boundary versus the published cycle** is now a recorded,
   deferred gap (R4-MIN-1) — verified real, explicitly not fixed this cycle.
2. **Quoted-`status:` frontmatter shapes** remain unhandled (R4-NIT-1, deferred) — no
   observed instance, root fix deferred to the same candidate.
3. **GitHub-hosted runner wall-clock remains unmeasured** (G3 removes the known cliff,
   not the unknown; both reviewers' local over-600 s measurements sharpen it).
4. **Live attribution windows remain test-proven only** — no non-zero-width window from
   a real driver transition yet; live state honestly `ambiguous`.
5. **Windows `wait`/`usage` portability remains macOS-verified only** (carried from
   cycles 1–3; unchanged).
6. **The LE-7 goal-done check has now run exactly once** — commissioned by the organizer
   via inbox note, executed by kimi-1 at the pinned closing HEAD, verdict PASS. It is
   defense-in-depth on top of this consensus, never its substitute (LE-7's own words in
   the verdict file).

**Correlated agreement (§15.6(b)).** Two model families, but the cross-confirmations
rest on shared priors (PRIMARY execution over prose; deterministic tooling trustworthy)
and partially identical command sequences. Their agreement is strongest where least
independent; this round's value is that the last two divergences (VC-7, VC-8) closed by
evidence and signoff rather than persisting, and the one remaining open item — kimi-1's
concurrence with claude-1's three dispositions — is put to kimi-1 directly rather than
inferred from its READY position.

## Verdict conflicts

**None open this round.** Both prior conflicts closed through the signoff mechanism
before cycle 3 opened (never by count, drafter, or organizer), as recorded in the signed
cycle-3 consensus and re-stated in both round-04 files:

- **VC-7 (readiness) — RESOLVED at cycle-3 signoff to the short cycle 3; both reviewers
  now READY.** claude-1 moved NOT YET → READY this round ("the gap is now dispositional
  rather than substantive" — confirmed by cycle 3 being one code fix on the signed G4
  surface plus six record corrections, with nothing new under re-attack); kimi-1 moved
  READY → NOT-YET (route 1) at cycle-3 signoff and back to **READY** at round 04 on its
  own verification (zero new findings, trajectory converging). The positions arrive at
  READY from opposite directions, each stated as a position change in its own file.
- **VC-8 (materiality of the equal-mtime fall-through) — RESOLVED at cycle-3 signoff at
  MINOR-fix-in-cycle; verified fixed.** claude-1 withdrew R3-MIN-1 on its own
  pre-committed falsifier, fired in its own hands at `4df0855` and independently in
  kimi-1's `git archive` extraction (equal mtimes, current binary `await review
  artifact`, pre-H1 binary `await implementation` on the identical tree).

**Readiness disagreement honestly stated: none exists.** Both reviewers file READY under
explicitly compatible conditions (claude-1's four-point closure position; kimi-1's
recorded close conditions). The conditions are restated in **Close conditions** below and
are the same list — nothing is majority-resolved because nothing is disputed. The one
item NOT yet on record from any reviewer is kimi-1's position on the three R4
dispositions themselves (defer/defer/dismiss): kimi-1's review neither concurs nor
objects — its signoff request below asks for that position explicitly, because a
deferral under the template means "*findings everyone agrees* are out of scope" and that
agreement must actually be supplied, not assumed from READY.

## Close conditions (LE-7/LE-11 manual equivalents — carried from the signed cycle-3 consensus H3(a))

This run is manually driven under the owner-authorized manual fallback, so the driver's
auto-complete path will not close the idea — but the manual close must satisfy and record
the protocol's own conditions (both round-04 reviewers file READY against exactly this
list):

1. **Zero Agreed fixes** — proposed by this draft (§ Agreed fixes: None).
2. **LE-11 reservations:** all-✅ signoffs, or a recorded operator ruling accepting any
   🟡. Nothing on either reviewer's filed position forces a 🟡 today; if one lands, the
   close escalates rather than completing.
3. **LE-11 reviewer floor:** two independent reviewers — satisfied by claude-1 + kimi-1.
4. **LE-7 goal-done:** commissioned, independent, non-implementer check **passed at the
   closing HEAD** — `review/goal-done/kimi-1.md`, verdict PASS at `4df0855` / source
   `c3baf09` / skill `b06a65a`, commissioned by the organizer via inbox note
   (`codex-1-to-kimi-1_…_goal-done.md`, `blocking: no`). The check can only withhold a
   close, never establish one: this PASS is defense-in-depth evidence, not the close.
   **Pinning dependency (claude-1 open question 1):** the verdict is pinned to `4df0855`
   and remains valid iff this consensus closes with zero fixes; any signoff amendment
   that creates a fix-up cycle 4 voids the pin and requires re-commissioning at the new
   HEAD.
5. **No release actions before the organizer's authorized release step** — nothing
   released, tagged, published, installed, or merged (verified untouched by both
   reviewers as of their reviews, and by the drafter this session — next section).
   The release step, the attended core publish, and any **post-deploy channel audit
   remain the organizer's/owner's post-Phase-8 work** — separate from and after this
   record; nothing here performs or implies them.

With all-✅ signoffs on THIS consensus, the implementer's Phase-8 close (set
`status: complete` in `IMPLEMENTATION.md` frontmatter, publish the completion message)
follows on the same branch. That step is NOT taken by this draft.

## Source-pinned evidence (concrete)

The reviewed source remains pinned and unmodified (DRAFTER-PRIMARY this session unless
attributed):

- `git show --stat 4df0855` → touches **only** `…/IMPLEMENTATION.md` (344 insertions,
  13 deletions). The record commit changes no code.
- `git show --stat c3baf09` → touches **only** `internal/driver/phasedigest.go` +
  `internal/app/wait_test.go`. Fix-up-3 source is exactly the H1+H7 surface the signed
  plan authorized.
- Frozen FINAL: `git diff 120a9bf..HEAD -- …/FINAL.md` → **empty**.
- Skill pinned: `b06a65a` unchanged since cycle 2 (verified PRIMARY by both round-04
  reviewers and the goal-done check against their own trees).
- Goal-done check (kimi-1, PRIMARY): CLI record tip `4df0855`; code `c3baf09` with
  `go version -m` → `vcs.revision=c3baf0977b82c2e9fc386b174f45de9e5a1dc1d9`,
  `vcs.modified=false`; two-party byte-identical rebuild of the recorded task binary
  (`19831a24…06bdc`).
- Staged core unchanged: `fc907e59…`, 109,772 B (cycle-3 record; re-verified against the
  live file by kimi-1 round-04).
- No publish by anyone: `~/.parley/protocol/core/` holds **only `2.10.0`**
  (DRAFTER-PRIMARY this session).

## Signoffs

*(Deliberation track: all three participants sign — reviewers and implementer. Signoffs
are intentionally empty in this draft; the cycle-3 signoffs live verbatim in
`review/consensus-cycle-03.md` and are NOT reused. Append-only; do not edit others'
blocks. This is the proposed **closing** consensus: all ✅ → Phase-8 close; any ❌ → new
review round with the blocker's counter-proposal as the starting point; any 🟡 → LE-11
escalation (operator ruling required to close over it). Each signoff request asks the
signer to: (1) state a binding position on **readiness** — both reviewers filed READY;
sustain or dissent with evidence; (2) **kimi-1 especially: concur or amend each of the
three dispositions** (R4-MIN-1 defer, R4-NIT-1 defer, R4-NIT-2 dismiss) — its review did
not address them, and the template's "everyone agrees" for deferrals requires its
explicit yes; (3) confirm the close-conditions list above is complete and satisfied-in-
pending (goal-done PASS pinned to this HEAD — noting any amendment creating a fix voids
the pin); (4) confirm the DF-5 candidate path suggestion and the preserved DF-1…DF-4 /
residuals list, with no scope additions; (5) confirm no release action is taken or
implied by signing. The implementer signs per §15.1 with no self-verification verdict on
its own implementation.)*

### Signoff: claude-1 — (pending)

### Signoff: kimi-1 — (pending)

### Signoff: zcode-1 — (pending)

### Signoff: claude-1 — 2026-09-24
Status: ✅ ACCEPT
Notes:

**What this ✅ accepts.** The zero-Agreed-fixes proposal as the **closing** record at
`4df0855`, the complete finding→disposition map, the DF list, and the close-conditions
section. This is a close, so I verified the three conditions I do not own **myself this
session** rather than reading them off the draft. (Appended after the drafter's empty
`(pending)` scaffold rather than overwriting it — same convention as the cycle-3 archive.)

**Readiness — binding position: READY, sustained.** The position change I filed this round
(NOT YET → READY) stands, and nothing in this draft moved me off it. Round 4 produced 3
items from me (0 CRITICAL, 0 MAJOR, 1 MINOR, 2 NIT) against 7 / 15 / 25 in the prior
cycles, the only non-trivial one needs a *new slug* rather than an edit to this idea, and I
do not ask for a fix-up cycle 4. That is the converging shape §4's stopping judgment
describes.

**My four round-04 conditions, re-verified (PRIMARY, my own tree, this session):**

1. **Dispositions written, and they match what I filed.** R4-MIN-1 → deferred to
   `wait-boundary-vs-published-fixup`; R4-NIT-1 → deferred to the same slug; R4-NIT-2 →
   dismissed as immaterial. I am the filer of all three and the draft records my own
   proposals without narrowing them. I re-confirm the dismissal **knowingly**: closing
   freezes the `IMPLEMENTATION.md:453` typo into the record and I consent to that.
2. **LE-7 goal-done at the closing HEAD.** I read `review/goal-done/kimi-1.md` itself, not
   the draft's summary of it: `reviewed-commit: 4df0855`, **verdict PASS**, author kimi-1
   — a non-implementer (implementer zcode-1) and not the declared facilitator — commissioned
   by the organizer, criteria verdicts PRIMARY. Per LE-7 I treat it as defense-in-depth that
   can only withhold a close; my ✅ does not substitute for it, and it does not substitute
   for the review consensus.
3. **LE-11.** Two-independent-reviewer floor met by kimi-1 and me. Nothing on my side forces
   a 🟡; I file a clean ✅ so no operator ruling is needed on my account.
4. **No release action.** `HEAD == 4df08551e965bc06262a3cd417220f4a407a0b07`;
   `git tag --contains HEAD` empty; HEAD reachable from no branch but `lean-organizer` and
   one kimi review branch (not `main`); `~/.parley/protocol/core/` holds **only `2.10.0`**;
   working tree carries no code modification (deck record files only); FINAL still frozen
   (`git diff 120a9bf..HEAD -- …/FINAL.md` empty).

**The pin (my open question 1) — I sign with it explicit.** The goal-done verdict is pinned
to `4df0855` and is valid **iff** this consensus closes with zero fixes. If any signer
amends a disposition into a fix, this stops being the closing record: cycle 4 opens and the
verdict must be re-commissioned at the new HEAD. My ✅ is a signoff on *this* record, not a
vote to close over a stale tree.

**Open question 2 — closes with no action.** I do **not** prefer the template shape for
`head-commit:`; the parenthetical is honest, no code parses the field, and amending it would
itself require a cycle. Leaving it as-is is my position, so the question closes.

**Dispositions are kimi-1's to concur in, not mine to supply.** The draft is right to put
the three dispositions to kimi-1 directly. My signoff does not carry its concurrence, and
READY is not concurrence — a deferral's "everyone agrees" has to be supplied, not inferred.

**DF-5 and scope.** I verified the suggested candidate exists at
`parley-deck/ideas/wait-boundary-vs-published-fixup/00-prompt.md` with `status: candidate`,
**no `participants:` staffed**, and "Inactive deferral record only. No participants are
staffed, no round opens" — LE-10/§14.1-conforming, staffing no quorum, adding no scope. It
is currently untracked; committing it is the organizer's, not a condition of my signoff.
DF-1…DF-4 and the residuals I filed (GitHub-hosted runner wall-clock unmeasured; live
attribution windows test-proven only; Windows `wait`/`usage` macOS-verified only) are
preserved as I filed them, and R4-MIN-1/R4-NIT-1 now join that list as I proposed. No scope
additions.

**Signing authorizes nothing downstream.** I take, and my ✅ implies, no release, tag,
publish, install, core publication, or merge. The release step, the attended core publish,
and the post-deploy channel audit remain the organizer's and owner's post-Phase-8 work. The
Phase-8 close (`status: complete` in `IMPLEMENTATION.md`) is the implementer's, after
all-✅.

**§15.1.** I issue no verification verdict on any claim I own; the facts above are PRIMARY
checks I ran in my own tree, and the reviewers' and checker's evidence remains theirs.

### Signoff: kimi-1 — 2026-09-24
Status: ✅ ACCEPT
Notes:

**What this ✅ accepts.** The zero-Agreed-fixes proposal as the **closing** record at
`4df0855`, the complete finding→disposition map, DF-5 joining DF-1…DF-4, and the
close-conditions section. Appended after the drafter's `(pending)` scaffold, per the
cycle-3 convention claude-1 followed. Before signing I re-ran the load-bearing record
facts myself this session (PRIMARY, this tree): `HEAD == 4df08551e965bc06262a3cd417220f4a407a0b07`;
`git show --stat 4df0855` touches only `IMPLEMENTATION.md` (+344/−13); `git show --stat
c3baf09` touches only `internal/driver/phasedigest.go` + `internal/app/wait_test.go`;
`git diff 120a9bf..HEAD -- …/FINAL.md` empty (FINAL frozen); `git tag --contains HEAD`
empty; HEAD reachable only from `lean-organizer` and my own preserved review branch, not
`main`; `~/.parley/protocol/core/` holds only `2.10.0` (no publish by anyone); staged core
`fc907e59…62c9f`, 109,772 B unchanged; the cycle-3 archive is 58,715 B at sha256
`419900f1…b9f116b`, exactly as the draft records; and my commissioned goal-done verdict
stands at `review/goal-done/kimi-1.md` (**PASS**, pinned `4df0855` / `c3baf09` /
`b06a65a`).

**Readiness — binding position: READY, sustained.** My round-04 READY was filed on my own
verification (zero findings, all seven H-fixes held under re-attack in my own trees, both
suites green under my own hands); nothing this session moves me off it. The finding
trajectory — 25 → 15 → 7 → 3 from one filer with none requesting a cycle, and 0 from me —
is the converging shape §4's stopping judgment describes. My cycle-3 churn commitment
binds symmetrically, and I honor it in the direction it binds me: round 04 produced no
fresh findings on fix-up code and no re-litigated ground, so I support closing over
dispositions rather than opening a cycle 5.

**The three dispositions — my explicit concurrence, each on its own merits.** My round-04
review filed no findings and did not address claude-1's R4 items; the draft is right that
a deferral's "everyone agrees" must be supplied, not inferred from READY. I supply it:

1. **R4-MIN-1 → defer: I CONCUR.** claude-1 asked for my read directly: is "create the
   round directory when you dispatch" the intended contract? My round-04 evidence says the
   deck's actual dispatch flow already satisfies it — at the live cycle-3-published state
   `parley wait --for review` exits 3 naming both reviewers *because* `review/round-04/`
   existed, and claude-1's own control shows correct behavior returns the moment the next
   round's directory exists. The residual gap (one `--json` envelope saying
   `next: await review artifact` while the terminal line exits 0) is real but bounded: it
   violates no frozen criterion (FINAL B.3 does not pin the boundary to the published
   cycle), it cannot mis-advance a phase (`boundaryReached` is referenced only by `wait`
   itself, so B.5's pure-observer property holds), and both offered fix shapes — the
   content signal in `boundaryReached`, or documenting the directory-creation convention
   in the skill — are legitimately the new slug's question to decide between. A fix-up
   cycle 4 is disproportionate to a MINOR with a working operational convention.
2. **R4-NIT-1 → defer to the same slug: I CONCUR.** I verified the zero-instance claim
   myself this session (PRIMARY): a repo-wide grep for quoted `status:` frontmatter over
   the 80-plus-idea corpus returns **0 files**. The quote tolerance predates this idea
   (`815c93a`), and the one-line root fix in `implSection` touches the same digest-parsing
   surface as R4-MIN-1 — one candidate, one conversation.
3. **R4-NIT-2 → dismiss as immaterial: I CONCUR, knowingly.** The typo at
   `IMPLEMENTATION.md:453` is real (re-read this session: "cycle 1 carries its at the
   equivalent position"), the underlying fact is correct (cycle 1's deviations subsection
   at `:815`, verified in my round-04), the surface is purely textual with zero behavioral
   effect, and the filer consents to freezing it. Closing over a named, known typo is an
   honest record; a cycle to fix it is not justified.

**Close conditions — the list is complete, and satisfied-in-pending.** The five items are
the complete rule set for this run: zero Agreed fixes; LE-11 all-✅-or-recorded-operator-
ruling (I file a clean ✅, so no ruling is needed on my account); the two-independent-
reviewer floor met by claude-1 and me; the LE-7 goal-done verdict at the closing HEAD; and
no release action before the organizer's authorized step. **The pin:** my goal-done
verdict is pinned to `4df0855` and remains valid iff this consensus closes with zero
fixes — my concurrence above maintains that condition. Any signoff amendment that creates
a fix-up cycle 4 voids the pin and requires re-commissioning the check at the new HEAD
(claude-1 open question 1; LE-7 rejects a stale code tree). Per LE-7 that check can only
withhold a close, never establish one; my ✅ here is a review-consensus signoff, not a
re-issue of my own verdict — §15.1: I own that claim and do not verdict it again here.

**DF-5 and scope — confirmed.** The candidate exists at
`parley-deck/ideas/wait-boundary-vs-published-fixup/00-prompt.md` with `status: candidate`,
**no `participants:` staffed**, and "Inactive deferral record only" (PRIMARY, read this
session) — LE-10/§14.1-conforming, staffing no quorum, launching nothing, adding no scope.
DF-1…DF-4 are preserved unchanged, and the three disclosed residuals — GitHub-hosted
runner wall-clock unmeasured, live attribution windows test-proven only, Windows
`wait`/`usage` macOS-verified only — stand as recorded, non-blocking under the frozen
FINAL. No scope additions.

**Signing authorizes nothing downstream.** I take, and my ✅ implies, no release, tag,
publish, install, core publication, or merge. The release step, the attended core publish,
and the post-deploy channel audit remain the organizer's/owner's post-Phase-8 work. The
Phase-8 close (`status: complete` in `IMPLEMENTATION.md`) is the implementer's, after
all-✅.

### Signoff: zcode-1 — 2026-09-24
Status: ✅ ACCEPT
Notes:

**What this ✅ accepts.** The zero-Agreed-fixes closing record at `4df0855` as drafted: the
complete finding→disposition map (R4-MIN-1 defer, R4-NIT-1 defer to the same slug, R4-NIT-2
dismiss as immaterial), DF-5 joining the preserved DF-1…DF-4 and residuals, and the
close-conditions section. Appended after the `(pending)` scaffold, per the convention both
reviewers followed.

**§15.1 — implementer discipline.** I authored this draft and the implementation it
dispositions; I issue **no verification verdict on my own implementation**. The
verification evidence in this record is the reviewers' (PRIMARY in their round-04 files and
re-verified in their signoff blocks above) and the commissioned checker's (kimi-1, PASS at
`4df0855`); their blocks carry it, and I do not re-assert it. What I attest is
implementer-side record fact: fix-up cycle 3 landed the signed H1–H7 plan with no
deviations (`## Deviations from agreed fixes (cycle 3): None` is my own record),
`IMPLEMENTATION.md` frontmatter stands at `status: fix-up-cycle-3` awaiting exactly this
close, and the draft's dispositions carry the filers' own filed positions without
narrowing — claude-1's round-04 disposition table and kimi-1's explicit item-by-item
concurrence above each supply that directly. My positions are unchanged since my cycle-3
signoff block (archived verbatim in `review/consensus-cycle-03.md`); no drafter position
changes beyond those already on record.

**Readiness — binding position: READY, as implementer.** Both reviewers filed READY on
their own evidence, from opposite directions, and their signoffs above sustain it; the
finding trajectory (25 → 15 → 7 → 3, none asking for a cycle) is the converging shape
§4's stopping judgment describes. Nothing in this session's checks gives me grounds to
dissent.

**Close conditions — complete and satisfied-in-pending.** The five-item list is the
complete rule set for this run. With this ✅ the all-✅ condition is met, the
two-independent-reviewer floor stands (claude-1, kimi-1), and the LE-7 goal-done PASS
stands pinned to `4df0855` — valid iff this consensus closes with zero fixes; any signoff
amendment that creates a fix-up cycle 4 voids the pin and requires re-commissioning the
check at the new HEAD. Record facts re-verified myself this session before signing (git,
not summaries): `git show --stat 4df0855` touches only `…/IMPLEMENTATION.md` (+344/−13);
`c3baf09` touches only `internal/driver/phasedigest.go` + `internal/app/wait_test.go`;
`git diff 120a9bf..HEAD -- …/FINAL.md` empty (FINAL frozen); `review/goal-done/kimi-1.md`
records verdict PASS at the pinned HEAD; the DF-5 candidate exists at
`parley-deck/ideas/wait-boundary-vs-published-fixup/00-prompt.md` with `status: candidate`,
no participants staffed, no round opened — LE-10/§14.1-conforming, no scope additions.

**Signing authorizes nothing downstream.** I take, and this ✅ implies, no release, tag,
publish, install, core publication, or merge; those remain the organizer's/owner's
post-Phase-8 work. The Phase-8 close — `status: complete` in `IMPLEMENTATION.md` plus the
completion message — is mine to perform on this branch after this all-✅, and is not
performed in this invocation.
