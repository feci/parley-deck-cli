---
agent: claude-1
idea: meta-protocol-change-designated-implementer
review-round: 1
date: 2026-09-25
reviewed-commit: 0893989911f61077ae2535c0842ce08e83f196e7
phase: review-round-01
role: independent non-implementer reviewer
implementer: kimi-1
verdict: CHANGES REQUESTED
branch: designated-implementer
reviewed-commit-cli: 0893989911f61077ae2535c0842ce08e83f196e7
reviewed-commit-skill: bf7e0495848db05327790a6540a7313db19c37c8
baseline-commit-cli: e4640bf2840249db0c1a1ecab7493813f4dacfdb
baseline-commit-skill: 8161e5e
---

## Protocol context attestation

`parley protocol packet --dir . --phase 6 --track deliberation --idea
meta-protocol-change-designated-implementer --flag auto_implement --flag protocol_change
--audience participant --json`, run in the original CLI worktree against the live source
authority (no refusal, no fallback):

```
context_mode     : "full"
source.path      : <worktree>/parley-deck/COOPERATION.md
source.role      : "source"   (live source file; global core drift never substituted)
source.bytes     : 114771
source_sha256    : c749218255c96c4efeecc8d598abc6192f195a294291eff6c505096f0091568f
packet_sha256    : c749218255c96c4efeecc8d598abc6192f195a294291eff6c505096f0091568f
index blocks     : 69
shadow           : {packet_sha256 3cab52e7…, 72857 bytes, 37 included / 32 omitted}  — NOT used
```

- `fallback_reason` is **ABSENT**, verified by enumerating the top-level keys:
  `body_path, context_mode, index, packet_sha256, request, shadow, source, source_sha256`.
- Rendered body verified byte-identical to the authority: `shasum -a 256` of
  `.parley-runtime/protocol-packets/full-phase6-deliberation-c7492182….md` and of
  `parley-deck/COOPERATION.md` both `c749218255c96c4efeecc8d598abc6192f195a294291eff6c505096f0091568f`,
  1400 lines / 114 771 bytes each. The Phase-6 body was read from the live file.
- The authority is the **amended** COOPERATION.md (post-hunk). The Phase-5 attestation in
  `IMPLEMENTATION.md` records the same hash for its post-edit run, so implementer and reviewer
  read the same text.

## Review provenance and method

- **Review worktrees.** Two private read-only clones on local disk, pinned to the exact commits
  under review, to avoid contending with zcode-1's fixtures/caches and to avoid the slow shared
  mount: `/private/tmp/claude1-rev-di/cli` @ `0893989` and `/private/tmp/claude1-rev-di/skill` @
  `bf7e049`. No file in either canonical worktree was modified by this review except this artifact.
- **Break tests.** Findings B1–B8 below were produced by Go test files I wrote myself
  (`internal/app/zz_claude1_review*_test.go`) in the throwaway clone only. They are **not** part of
  the delta, were never committed, and were moved out before I re-ran the implementer's own suites.
  Every "OBSERVED" line quoted below is verbatim `go test -v` output from `0893989`.
- **Refutation default.** No implementer pass claim was accepted on its word. Every AC below was
  re-derived from the tree; where I could not derive it, I say so under
  `## Limitations — what I did NOT verify`.
- **Scope discipline.** Windows architecture stays deferred; nothing below asks for it. I propose no
  edits to peers' artifacts, FINAL, IMPLEMENTATION, signatures or global config.
- Environment: go1.27.1 darwin/arm64.

## Verdict

**CHANGES REQUESTED.** The mechanism is well built, genuinely dormant when unset (independently
confirmed byte-for-byte), and the three-copy protocol fidelity is exact. Three MAJOR findings block
a clean pass: one escalation the FINAL designed can be bypassed by *deleting* the line it guards;
one sentence of the **newly shipped protocol text** asserts a gate the code does not implement on
design-only runs; and the tier-2 availability gate can be cleared by a waiver naming a *different*
agent. None is a crash, a data-loss or an unset-path regression.

---

## Findings

### MAJOR-1 — R29's anti-silent-reassignment escalation is bypassed by removing the designation

**Where.** `internal/app/driver_impl.go` — `Implement` guards the whole re-entry comparison behind
`if o.implDesignated`, and `checkImplementerReentry` is therefore unreachable once the designation
is gone.

**What FINAL says.** R29: *"a live tier-2/tier-3 change against a recorded dispatch **escalates**
rather than reassigning silently. That comparison **fires only when the recorded source was
`designation` or `global-default`**."* The firing condition FINAL states is a property of the
**recorded** event. The implementation additionally requires a designation to be **currently
present**, and the two conditions come apart in exactly the case that matters.

**Repro (B1, verbatim).**

```
=== RUN   TestB1_DeletingTheDesignationAfterARecordedDispatch
    after deletion: implementer="aa-first" source="" designated=false
    SILENT REASSIGNMENT: recorded dispatch zz-impl via designation; after deleting the line
    the run dispatched "aa-first" with no escalation (err=<nil>)
        stdout:
        driver: implementing via aa-first ...
```

Steps: idea with `participants: [aa-first, zz-impl]` and `implementer: zz-impl`; append
`agent.implementer_resolved {idea: demo, implementer: zz-impl, source: designation}` (i.e. a dispatch
happened but Phase 5 did not get as far as writing `IMPLEMENTATION.md`, so there is no pin); delete
the `implementer:` line; construct ops and call `Implement`. Result: `aa-first` is dispatched, no
error, no notice.

**Tier 3 is the same (B2, verbatim).**

```
=== RUN   TestB2_ClearingTheGlobalDefaultAfterARecordedDispatch
    SILENT REASSIGNMENT (tier 3): recorded zz-impl via global-default; after clearing it
    the run dispatched "aa-first" with no escalation
```

**Why it matters.** R28 rejects ALT-15 because *"self-appointment with a paper trail is not
authorization"*. Deleting one line is self-appointment with **no** paper trail, and it is strictly
easier than the route R28 blocks: changing the designation to another id **does** escalate
(the implementer's own `TestReentryComparisonEscalatesOnDesignationChange` proves it), so the
cheapest way past the gate is to remove the field rather than edit it. The same hole covers the
combination "write an `IMPLEMENTATION.md` pinning yourself **and** delete the designation" — the R28
conflict gate needs the designation to still be there to see a conflict.

**Honest counter-argument.** R29 carries the `[D]` tag, and `[D]` is defined as *"fires only when a
designation is present (tier 2 or tier 3)"*. Under a literal reading of the tag alone, the
implementation is compliant and the rule simply has this hole by design. I do not think that reading
survives the rule body ("escalates rather than reassigning silently") or R28's stated rationale, but
the tension is in FINAL, not invented by me, so this is a rule-interpretation call the review
consensus should make explicitly rather than absorb silently.

**Suggested fix (one of).** (a) Hoist the comparison out of the `implDesignated` guard: run
`checkImplementerReentry()` in `Implement` unconditionally — it already returns `nil` when the store
is empty, when no matching event exists, and when `implSource == pin`, so the unset-path cost is one
`Store.Load()` and **no** behaviour change (AC-4's four clauses are all about the event, the stdout
line, `ExpectedRoundParticipants` and the three legacy tests — none of them is touched by a read).
Or (b) if the review consensus prefers the literal `[D]` reading, record the hole explicitly in
`IMPLEMENTATION.md` and in the deferred register, rather than leaving R29 reading as a closed
protection.

### MAJOR-2 — the newly shipped protocol text claims a launch-blocking validity gate that does not exist on design-only runs

**Where.** The Phase-5 paragraph added to **all three** `COOPERATION.md` copies:

> "Validity gates are hard and fire on **any** run: an empty value, an id outside this idea's
> eligible participants, or a malformed value **blocks the launch**…"

**What the code does.** The gate rides `roleErr`, which is read by exactly four role actions
(`Implement`, `OpenReviewRound`, `GoalCheck`, `Fixup`). A design-only run reaches none of them:
`internal/driver/impl.go:81` — `if !d.cfg.AutoImplement { return ActionSurfaceOnly, c, nil // idea
did not opt in / --no-implement → stop at FINAL }` — returns before `Impl.Implement` is ever called,
and `OpenReviewRound` is only reached from `advanceImpl`, i.e. after `PhaseImpl`. Nothing else
consults `roleErr`, and F4 (the optional idea-scoped preflight copy of R26) was deliberately not
shipped. So on `auto_implement: false` or `parley run --no-implement`, a defective `implementer:`
line does not block anything, and the launch is never blocked on any run — the four role actions are
mid-run, not launch.

**Evidence (B6, verbatim escalation surface for a present-empty designation).**

```
validity-gate escalation surface: map[Complete:true Fixup:true GoalCheck(blocked):true
  Implement:true ImplementationStatus(err):true OpenReviewRound:true PrecheckFixup:false
  ReviewStatus(err):true]
```

(`Complete:true` / `ImplementationStatus`/`ReviewStatus` here are *other* errors — no
`IMPLEMENTATION.md`, no review dir — not the role gate; `Complete` provably has no `roleErr` check,
see the note under "confirmations" below. The four genuine `roleErr` sites are the ones FINAL names.)

**What FINAL says.** R16 [A]: *"Validity gates are hard and fire on **any** run, **including a
design-only run**: … A defective frontmatter line is worth stopping on even where Phase 5 is
unreachable."* This is not met. The implementer's F4 rationale — *"the `roleErr` authority is
complete without it"* (`IMPLEMENTATION.md`, Decisions) — is the specific claim that fails: R25's
`roleErr` siting is the right **authority**, but on a design-only run it has no **trigger**.

**Why this is MAJOR rather than MINOR.** Declining the optional preflight copy (R26/F4) is the
implementer's licensed choice. Shipping a sentence into three protocol copies that asserts a gate
the tooling does not apply is not: this text is the thing about to be published, and the protocol is
the artifact agents are told to trust over their own memory. The mismatch is self-inflicted by the
delta and is cheap to remove.

**Suggested fix (either, not both).** (a) Ship the R26 idea-scoped preflight copy so the claim
becomes true (R26 permits it explicitly; it must not iterate all ideas — `facilitatorConflictGates`
shape is the trap R25 names); **or** (b) correct the protocol sentence in all three copies to state
what the code does, e.g. *"Validity gates are hard and fire on any run that reaches an implementer,
reviewer, goal-check or fix-up action; a design-only run surfaces them when it first dispatches."*
Option (b) preserves AC-1 by keeping the hunks identical and is the smaller change; option (a) is
what R16 actually asked for. This is a fix-up call for kimi-1 with the consensus, not a reviewer's
to make.

### MAJOR-3 — the tier-2 availability gate is cleared by a waiver naming a different agent

**Where.** `internal/protocol/implementer.go:183` —
`if !strings.Contains(strings.TrimSpace(parts[0]), id) { return false }`. Substring, not identity.

**Repro (B5, verbatim).**

```
=== RUN   TestB5_WaiverPrefixCollision
    prefix-collision waiver: dispatchErr="" source="fall-through-unavailable"
    OBSERVED: a waiver naming zz-impl-2 cleared zz-impl's unavailability gate (substring match)
```

Fixture: `implementer: zz-impl`, `implementer_waived: zz-impl-2 — on leave — confirmed 2026-09-25`,
with `zz-impl` undiscovered. Expected: the R18 blocking gate. Actual: no gate, fall-through recorded
as `fall-through-unavailable`.

**Why it matters.** R18's gate exists precisely to require an **explicit, per-agent, user-confirmed**
record before a designation is set aside — it borrows §9.0's confirmation shape on purpose. A reader
that matches any id **contained in** the waiver line fails open: one confirmed waiver for
`claude-1-fast` silently waives `claude-1`. Note the asymmetry: the sibling reader
`ParseImplementerReassignment` gets this right (exact `==` on both sides of the separator), so the
laxity is isolated to one function and one line.

**Suggested fix.** Compare identity, not containment — split the subject on whitespace and require
an exact element match (or `strings.TrimSpace(parts[0]) == id` if the shape is meant to be
`{agent-id} — {reason} — confirmed {date}`, which is what R18 and the protocol text both write).
Add a test with two ids where one is a prefix of the other.

### MINOR-1 — `none` is treated as "designation present" for the `[D]` kickoff surface, so the documented suppressor does not suppress

**Repro (B4 / B4b / B7b, verbatim).**

```
none: designated=true out="driver: WARNING designated run leaves a single non-implementer reviewer
  (zz-impl); aa-first implements and never reviews itself — consider a third participant…"
global none: designated=true out="driver: WARNING designated run leaves a single non-implementer
  reviewer (zz-impl); …"
agent.model_diversity events — designated: 2, unset: 1, `implementer: none`: 2
```

`resolveDispatchDesignation` returns `present: true` for both `implementer: none` and
`default_implementer = "none"`, and `newDriverImplOps` calls `kickoffDesignationChecks()` on
`present`. Consequences on a deck that explicitly declined designation: (a) the operator is told
this is a "designated run" when it is the opposite; (b) `checkModelDiversity` runs an extra time at
construction, so the always-on `agent.model_diversity` event is recorded **twice** instead of once.

R45 does legitimately count an opt-out as "present" **for emission** — that part is correct and
required. The defect is that one boolean is doing two jobs: R45's emission presence and R36–R39's
`[D]` behaviour presence. `none` is an explicit non-designation, so the `[D]` behaviours should not
fire. The blast radius is widest for exactly the configuration R9 documents as the deck-wide off
switch: `default_implementer = "none"` makes **every** idea on the deck print designated-run output.

**Suggested fix.** Split the flag: keep `present` for the event/line, add e.g. `live bool` (true only
for `DesignationSet` and an applicable tier-3 id) and gate `kickoffDesignationChecks()` on that.
`designatedImplementerPreference` already draws this line correctly for R36 — reuse its predicate.

### MINOR-2 — a malformed `default_implementer` is pin-shadowed: it gates some ideas and is silently ignored on others

**Repro (B3, verbatim).**

```
with pin: roleErr="" implementer="aa-first" source="pin" designated=true
OBSERVED: a malformed default_implementer does NOT gate an idea that already has a pin
```

With `default_implementer = "zz-impl # copied from a comment"`, an idea with no `IMPLEMENTATION.md`
hard-gates (correct, AC-10 / R19 concession 3), while an idea that already carries a pin resolves
cleanly and records `source: pin`, `designated: true`. The tier-3 value is read (to set `present`)
but never validated on that branch — `driver_impl.go`, the `DesignationAbsent`+`pinOK` case.

The rank-1-stops-the-chain reading of R13 defends this, and it never mis-dispatches (the pin is what
the owner wanted). But R19 concession 3 states the hard failure without a pin exception, and the
operator-visible result of one typo in one config file is that half a deck stops and half does not,
with no message on the half that does not. I record it as MINOR because it is fail-safe, not
fail-open.

**Suggested fix (if the consensus wants uniformity).** Validate the tier-3 value on the
`DesignationAbsent`+`pinOK` branch too, or — cheaper and arguably better — emit the malformed-value
notice there instead of gating, so the typo is always visible somewhere.

### MINOR-3 — a standing `default_implementer` is silently dropped when the layered config errors

`globalDefaultImplementer` (`driver_impl.go`) does `if err != nil { return "" }`. `LoadDefaults`
returns an error for (a) malformed TOML in **any** layer and (b) a missing
`$PARLEY_HEADLESS_AGENT_CONFIG` file, which `configLayers` marks `optional: false` — R11 itself
calls that layer *"highest, non-optional when set"*. In both cases the owner's standing preference
disappears with no message and the run falls through to positional dispatch. The comment says these
paths "already surface" the error elsewhere; that is true of the commands that call `LoadDefaults`
directly, and not of a driver dispatch that never calls them.

**Suggested fix.** Carry the read failure as a one-line notice on the designation path (it cannot be
a gate without violating R16's scoping, and should not be), or at minimum record it in the resolved
event.

### MINOR-4 — `IMPLEMENTATION.md` frontmatter points at the baseline, not the implementation

`head-commit: e4640bf` is the frozen-FINAL commit. The implementation is CLI `0893989` and skill
`bf7e049`; `e4640bf` contains none of it. `branch:` likewise records only the CLI worktree path and
does not name the skill worktree/commit, although the delta spans both repositories and AC-1/AC-3
both depend on the skill side.

I did **not** treat the stale frontmatter as permission to review the baseline; everything above is
measured at `0893989` / `bf7e049`. Recording it because §15 traceability is the point of the field,
and a fresh agent resuming "from FINAL + IMPLEMENTATION alone" (R40's own requirement, which this
idea just strengthened) would check out the wrong tree.

**Suggested fix.** `head-commit: 0893989` plus a `skill-commit: bf7e049` (or a two-line
`branch:`/`commits:` block naming both worktrees).

### MINOR-5 — duplicate `## Validation evidence` section in `IMPLEMENTATION.md`

Two headings with the same name: line 182 (populated) and line 247 (an empty "(Living — filled as
checks run…)" placeholder), both present in the committed file at `0893989`. Any reader or tool that
takes the first or last match gets a different answer. **Suggested fix:** delete the empty one at
line 247.

### NIT-1 — the R57 skill line drops the undesignated-claim case

`skills/parley-deck/references/ROSTER_AND_PROTOCOL.md:64` now reads "…then the FINAL drafter as
fallback; **a claim does not override a live designation**; implementation must follow `FINAL.md`…".
The baseline line said "…unless another participant claims it". The new line states the restriction
and drops the permission, so a reader of this file alone learns that claims are subordinate but not
that the volunteer route still exists with no designation — which is R15's explicit second half
("with no designation, a claim is today's normal volunteer route"). One clause would restore it.

### NIT-2 — `agent.model_diversity` is now recorded twice per designated run

Measured: 2 events vs 1 on an unset run (B7b above). This is a direct and *licensed* consequence of
R37 ("`checkModelDiversity` runs **unchanged** … Under a designation it **also** runs at kickoff …
its always-on `agent.model_diversity` event … untouched"), so it is not a rule violation. Recording
it because it is a real change to the event stream on the designation path and any consumer that
counts these events (TUI/state) will see the doubling. Worth one sentence in `IMPLEMENTATION.md`.

### NIT-3 — pre-existing drafter precheck does not mirror the new preference

`precheckDrafterLaunch` (`internal/app/driver_precheck.go:73`) selects with `firstHeadlessAgent`,
not `firstEligibleHeadlessAgentPreferring`, so the protocol precheck can render for a different
agent than `runDrafter` will actually launch. **This divergence pre-dates the delta** (the precheck
already ignored facilitator ineligibility), and R36 widens it in kind rather than creating it. Out of
scope for this idea — flagged, not requested, per the "mention, don't fix adjacent code" rule.

---

## Refutation attempts

Per acceptance criterion. Method: for each AC I tried to construct a case that breaks it, and report what I actually ran.

**AC-1 — Three-copy fidelity. PASS (independently re-derived, not taken from the tail-hash claim).**
Refutation attempted two ways. (i) Tail hashes at `0893989`/`bf7e049`: `tail -n +160` of the deck
copy and `tail -n +153` of `internal/protocol/defaults/COOPERATION.md` and of the skill's
`references/COOPERATION.md` all hash to
`da9704d57888efe549a2c781ba3bcdf34012c3caa81de9977fff3a1bc7e8eaf7`. (ii) Because the §0 hunk sits at
line ~59 — **above** the tail window, so the tail hash cannot cover it — I diffed the files whole:
deck-vs-embedded yields exactly `3c3`, `6,7c6`, `154,159d152`; embedded-vs-skill yields exactly
`5,6c5,6`. I then re-derived the same two difference sets at the baselines (`e4640bf` /
`8161e5e`, tail hash `f67af415…`): **identical hunk sets before and after**, so the delta introduced
no new divergence and the §0 sentence is byte-identical in all three. Sizes 1400/1393/1393 lines.

**AC-2 — Changelog. PASS.** `parley-deck/meta/protocol-changelog.md` gains a dated entry carrying
`Idea:`, `Drafted by:`, `Summary:` and an explicit `**Status: UNRELEASED.**`, prepended above the
2026-09-23 lean-organizer entry in the same shape.

**AC-3 — Skill companion consistency. PASS (with NIT-1).** Ran the prescribed check in the skill
worktree: `grep -rn "default implementer" skills/ lib/ bin/` → **0 hits**. Widened it myself:
`grep -rn "FINAL drafter"` and `grep -rn "default_implementer"` outside `references/COOPERATION.md`
each return exactly one hit, `ROSTER_AND_PROTOCOL.md:64`, which states the amended chain. `SKILL.md`
restates no Phase-5 implementer rule (its only role text is the facilitator exclusion at :119–121,
which the delta does not contradict). No contradicting restatement survives anywhere in the tree.

**On R57 specifically (the disclosed-unreviewed item).** I find the skill-companion edit is a
**correct and required consequence of the signed consistency requirement, not an expansion of it**:
AC-3 — a signed criterion — names `ROSTER_AND_PROTOCOL.md`'s Phase-5 line explicitly and requires it
to state the amended rule. Leaving it unedited would have failed a signed AC. The content is
accurate against the shipped protocol text (same four ranks, same order, same claim subordination),
and the file is documentation with no runtime effect. My only reservation is NIT-1 above. So: the
unreviewed status was disclosed honestly, and on review the item holds.

**AC-4 — Unset-path invariance. PASS (independently reproduced).** My own B8 asserts the stdout
byte-for-byte on a real dispatch through the fixture implementer: output is exactly
`"driver: implementing via aa-first ...\n"` and nothing else. Static confirmation of the other three
clauses: the event and the extra line are both inside `if o.implDesignated`, which is false on an
unset deck (B1's log line `designated=false` shows it directly); `expectedRoundParticipants` contains
no designation read at all (grep for `ImplementerDesignation|DefaultImplementer|default_implementer`
in `internal/consensus/*.go` non-test → one **comment**, zero code); the three R33 legacy tests are
unmodified in the diff (`git show --stat` lists no `app_test.go` / `roundgate_test.go`).
Refutation I could not achieve: I tried to find any construction-time write on the unset path —
`kickoffDesignationChecks` is guarded by `present`, and the only added unset-path work is extra
*reads* (`LoadDefaults`, one extra `ReadFrontmatter` for the pin). No output, no event.

**AC-5 — Four-state parse. PASS.** `internal/protocol` suite green (`ok … 0.383s`); the malformed
inline-comment value is kept literal by `ImplementerDesignationFromMeta` (I read the function: only
`TrimSpace` + quote-trim, no comment stripping) and fails the membership check downstream.

**AC-6 — `none` does not collapse. PASS.** Distinguishable from absent in the parser
(`DesignationNone` vs `DesignationAbsent`) and at dispatch (`source: none`, `present: true` vs zero
state). My B4 exercised it and confirms the state is separate — the finding there is about the
*consequence* of `present`, not about the state collapsing.

**AC-7 — Config precedence. PASS.** `internal/config` suite green (`ok … 0.487s`) including the
two-real-layer case. Code re-read: the merge is the non-empty-string form (R8), so `""` cannot
clear — confirmed by reading `mergeDefaults`, not only by the test.

**AC-8 — Tier-2 gate and its exits. PARTIAL — the gate and exits 1 and 3 hold; exit 2 fails open
(MAJOR-3).** The gate itself is correct and correctly scoped to the dispatch actions.

**AC-9 — Dormant designation paths. PASS for all five states.** Re-read each branch and confirmed
distinct sources: absent → zero state; non-participant tier 3 → `fall-through-inapplicable`; `none`
→ `none`; own declared non-participating facilitator at tier 3 → `fall-through-inapplicable` (and
correctly a hard gate at tier 2, R20's split); ping-failed tier 3 → `fall-through-unavailable`, with
the notice naming the designee and all three exits. Refutation attempted: I checked whether any of
the five can reach `roleErr` or `dispatchErr` — none does, so "no gate" holds.

**AC-10 — Invalid tier-3 still fails. PASS on the stated case; MINOR-2 on an unstated one.** A
whitespace-bearing value hard-gates. The implementer's line-drawing (whitespace after trim =
malformed; whitespace-free-but-ineligible = inapplicable) is a reasonable resolution of R19
concession 3 against R20/R22 and is recorded in `IMPLEMENTATION.md`. The gate is pin-shadowed
(MINOR-2).

**AC-11 — Pin conflict escalates. PASS.** Escalation fires on disagreement; `implementer_reassigned:`
re-pins only with a confirmed record, and its parser uses **exact** id matching on both sides (unlike
the waiver reader — see MAJOR-3). Refutation attempted: a reassignment record naming the right pair
but missing "confirmed" → correctly rejected (read of `ParseImplementerReassignment`); a record with
no separator → rejected (`len(parts) < 2`).

**AC-12 — Degenerate drafter fallback. PASS.** Two-pass `firstEligibleHeadlessAgentPreferring`; pass 1
re-admits the designee, and `preferNot == ""` short-circuits to a single pass so undesignated
selection is unchanged.

**AC-13 — Gate scoping. PASS.** Every read in `resolveDispatchDesignation` is keyed to `ideaDir`;
nothing iterates `status.Ideas`. Refutation attempted: grepped the delta for the
`facilitatorConflictGates` shape (`for _, idea := range status.Ideas`) — not present. R25's trap is
avoided.

**AC-14 — Two-participant kickoff warns. PASS** (it prints and returns; no error path). See MINOR-1
for the case where the same warning fires when it should not.

**AC-15 — Review exclusion reads the record. PASS.** `internal/consensus` suite green
(`ok … 0.933s`). Structurally re-derived rather than trusted: the only change to
`internal/consensus/consensus.go` is `resolveImplementer` delegating to
`ResolveImplementerChain(LegacyImplementerCandidates(), participants)` — same two artifacts, same
key order, same raw (unfiltered) eligibility list, same `""` tail. `ExpectedRoundParticipants` and its
three consumers (`reviewConsensusVoters`, the completeness check, `phasedigest.go:364` for
`parley wait`) are untouched. R34's divergence is now an explicit parameter at both call sites, as
required.

**AC-16 — Inertness. PASS.** `grep -rn "impl-claim" --include="*.go" .` → **0 hits**. Writers to
`00-prompt.md`: only `internal/protocol/workspace.go:225` (`CreateIdeaFull`) and
`internal/app/pipeline_cmd.go:1024` (the pre-existing pipeline seeder) appear as real file writes;
the diff adds only `ReadFrontmatter` reads and markdown text — verified by scanning every `+` line
of the Go delta for `WriteFile|Create|OpenFile|Rename`.

**AC-17 — Whole-tree health. PASS for build/vet/gofmt and for every test I was able to run;
the full `go test ./...` had not finished inside my window — see Limitations.** `go build ./...` OK,
`go vet ./...` OK (both exit 0, go1.27.1). `gofmt -l` over every `.go` file in the commit: clean.
`go test -count=1`: `internal/protocol` ok 0.383s, `internal/config` ok 0.487s, `internal/consensus`
ok 0.933s. In `internal/app`, all fourteen new designation tests plus the R33 legacy pin test ran
green on a clean `-count=1` selection (`ok … 0.654s`): `TestResolveImplementerFromRoleMetadata`,
`TestUnsetPathIsByteIdentical`, `TestTier2DesignationDispatchesDesignee`,
`TestTwoLayerGlobalDefaultDispatchesHigherLayer`, `TestNoneSuppressesTier3`,
`TestTier2UnavailabilityGateAndExits`, `TestPinDesignationConflictEscalates`,
`TestDormantTier3StatesFallThrough`, `TestMalformedTier3HardFails`, `TestMalformedTier2Gates`,
`TestDesignationGateIsScopedToTheIdea`, `TestTwoParticipantDesignatedKickoffWarnsNotBlocks`,
`TestKickoffModelDiversityUnderDesignation`, `TestReentryComparisonEscalatesOnDesignationChange`,
`TestDrafterSeparationPreference` — all PASS. The **whole** `internal/app` and `internal/driver`
packages were still running when I finished.

**AC-18 — End-to-end designation demonstration. PASS (re-run independently).** Both
`TestTier2DesignationDispatchesDesignee` and `TestTwoLayerGlobalDefaultDispatchesHigherLayer` drive a
real fixture implementer process through `RunImplementation` — a genuine end-to-end dispatch, not a
resolver assertion — and both PASS on a clean `-count=1` run at `0893989` (see AC-17). I also
reproduced the surrounding harness in my own B-tests (same `designationOps` path, real `Implement`,
real stdout): B1 and B8 both observe the real dispatch line.

**AC-19 — Ships UNSET. PASS.** `centralDefaultTemplate` emits `# default_implementer = "agent-id"`
— commented out, with a sentence recording the deliberate deviation from the template's active-key
shape. `grep -rn "default_implementer" --include="*.toml" .` over the whole tree → **0 hits**, so no
deck config carries the key. The only `default_implementer` strings in the repo are protocol prose,
Go source and tests.

**AC-20 — Attended boundaries intact. PASS.** Scanned every added line of the commit for
`pty|tty|IsTerminal|publish` — no match in any non-prose hunk. The TTY refusal
(`internal/app/protocol.go:363–377`, "requires a controlling terminal… Only the user may change the
global protocol") is untouched by the delta. No channel/publication/global-config action appears
anywhere in the diff.

**AC-21 — Independent review actually happened.** This file is my half of it (non-implementer,
populated `## Refutation attempts`). zcode-1 owns the other half; the goal-done check belongs to a
fresh non-implementer and is not claimed here. Not self-satisfiable, as FINAL requires.

## Other claims I checked rather than accepted

- **FINAL's own locator correction, recorded by the implementer as a "Locator note, not a
  deviation": CONFIRMED CORRECT.** At the baseline `e4640bf`, `internal/app/driver_impl.go:504` is
  `Fixup`'s `roleErr` check and `:518` is `func … Complete`, which has **no** `roleErr` check
  (`Complete` delegates straight to `completeWithWriter`). FINAL R25's citation of `Complete (:504)`
  was wrong; kimi-1 re-measured instead of transcribing, exactly as VC-4's standing rule demands.
- **`ResolveImplementerChain` is behaviour-preserving for the legacy path: CONFIRMED** by reading old
  and new side by side — same source order, same per-key order, same trim, same eligibility
  predicate, same `""`/`participants[0]` tails left at the call sites.
- **R27 (rank 1 wins when tiers 1 and 2 agree): CONFIRMED**, and note the conflict check is
  predicated on `pinOK`, i.e. on a pin that is itself eligible — an ineligible pin simply cannot
  conflict. That is sensible and worth nobody's time, recorded only so it is not mistaken for a gap.

## Limitations — what I did NOT verify

Stated plainly rather than papered over:

1. **I did not see `go test ./...` finish.** The whole `internal/app` and `internal/driver`
   packages were still running at the end of my window (even on local disk). I verified
   `internal/protocol`, `internal/config` and `internal/consensus` green with `-count=1`, the
   fifteen named `internal/app` tests green with `-count=1`, and `go build`/`go vet`/`gofmt` clean.
   I make **no** PASS claim for AC-17's full `go test ./...` clause and no claim about the remaining
   `internal/app`/`internal/driver`/`internal/trajectory` tests; the implementer's 505 s/605 s
   report is plausible but unverified by me. In particular I did **not** re-run the two R33 legacy
   `internal/consensus/roundgate_test.go` tests by name — they are inside the green
   `internal/consensus` package run and are unmodified in the diff, which is how I read that clause
   of AC-4 as satisfied.
2. **I did not run the driver end-to-end** (`parley run` against a real designated deck). All
   dispatch evidence is at the `driverImplOps` boundary, which is where the mechanism lives, but a
   real `parley run --no-implement` on a deck with a defective `implementer:` would make MAJOR-2
   demonstrated rather than traced. The trace is from source (`internal/driver/impl.go:81`) and I am
   confident in it, but it is a trace.
3. **I did not exercise the skill's installer/manifest/packaging.** The skill delta is two markdown
   files and I checked their content and AC-1/AC-3 only; I did not run the skill test suite,
   `npm test`, the manifest `--check` gate, or `npm pack`. The commit claims no installer/manifest
   change and the file list is consistent with that, but I did not prove the packaging is unaffected.
4. **No Windows evaluation.** Deferred by the owner; I looked for nothing and claim nothing. The one
   adjacent observation I can make factually: this delta adds no platform-conditional code and no new
   path handling beyond `filepath.Join`, so it introduces no new Windows surface of its own. It also
   does not reduce the existing deferred risk.
5. **Release/staging rules R50–R57** bind the organizer's release step, not the code. I verified only
   R57 (the skill companion) and confirmed no publication/channel action is present in either commit.
   I did not assess release readiness, versions, or channel state.
6. **Prior-run inherited limitations** (`source-context/prior-release-handoff.md`) were read for
   context but not re-verified; nothing in my findings depends on them.

## Nothing requiring human-only escalation

MAJOR-1's rule-interpretation tension (the `[D]` tag versus R29's body) is a decision for the review
consensus, and MAJOR-2 offers two legitimate dispositions — but both are ordinary protocol-review
judgement with the evidence in hand, not owner-only calls, so I am raising neither to the user. No
finding touches quorum, roster, attended boundaries, or the global default's UNSET shipping posture.

## Signoff

Not appended in this invocation, per the review instructions. This file is claude-1's own artifact;
no peer artifact, signature, FINAL, IMPLEMENTATION or global config was modified.
