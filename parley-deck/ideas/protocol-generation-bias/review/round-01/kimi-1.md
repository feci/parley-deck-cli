---
agent: kimi-1
idea: protocol-generation-bias
review-round: 1
date: 2026-08-29
reviewed-commit: 59eb663
---

## Summary

The protocol text moved correctly and the byte math is honest — §15.6 went from 1,372 B to
896 B (net −476 B), all three `COOPERATION.md` copies carry the identical new section, and the
prompt template change is real. But at 59eb663 the **enforcement story is mostly absent**: the
new round-1 gate (`protocol.ValidateRoundOneArtifact`) has zero runtime callers, the leg-3
disposition scanner does not exist while its protocol clause promises one, and leg 2 was deferred
by implementer decision without any amendment to the ratified acceptance criteria it voids. This
change ships prose and a test, not behavior — one release of the exact defect class the idea
exists to kill, twice over (unwired validator; unenforced clause (c)).

**Verdict: request changes.** One CRITICAL (wire the gate), two MAJORs that are cheap fixes
(§15.7 row; FINAL.md amendment regularizing D4), one MAJOR that is real work (clause (c) scanner
or strip its enforcement sentence). Rulings on D1 and D4 are in Findings 4 and 5.

Note on review conditions: during my review, uncommitted working-tree edits appeared in
`internal/runner/validation.go` (wiring the §15.6(a) check into the runtime validator) and
`internal/protocol/defaults/COOPERATION.md` (fixing the §15.7 row) — someone is mid-fix-up on
exactly Findings 1 and 2. I did not touch those files; findings below stand against 59eb663 and
the live edits corroborate rather than cover them. (PRIMARY: `git status --porcelain` →
` M internal/protocol/defaults/COOPERATION.md`, ` M internal/runner/validation.go`; diff inspected.)

## Refutation attempts

I assumed the implementation was wrong and tried to break each load-bearing claim. What failed
to break it:

1. **Byte claim.** `git show protocol-generation-bias-baseline:parley-deck/COOPERATION.md | sed -n
   '1346,1368p' | wc -c` → **1372**; current `sed -n '1346,1361p' parley-deck/COOPERATION.md | wc -c`
   → **896**; defaults copy (`1339,1354p`) → 896; skill copy → 896; pairwise `diff` of the three
   sections → empty. Net **−476 B** (implementer's figure), better than the ratified −237 B; both
   net-negative, so criterion 1 survives. Could not break it. (PRIMARY, CONFIRMED.)
2. **Gate sensitivity.** Inverted the condition in `internal/protocol/roundartifact.go:28`
   (`if !HasNonEmptySection` → `if HasNonEmptySection`, my own mutation, distinct from the
   implementer's `len(raw) < 0`): `TestValidateRoundOneArtifactRequiresExistingAlternatives`
   **failed immediately** ("a round-1 artifact without an Existing alternatives section must be
   rejected"). After restore, `-count=1` PASS and `git diff` empty for the file. The gate is a
   gate. (The restore raced with a concurrent process — my second Edit found the file already
   back to the committed state; I verified final content byte-identical to HEAD.) Could not break
   the gate itself; what broke instead is its call graph — see Finding 1. (PRIMARY.)
3. **D4's premise.** `grep -ri "exchange" --include='*.go'` over the whole repo → **zero
   matches**. No runner stage exists; withholding the clause was substantively justified. Premise
   survives. (PRIMARY, CONFIRMED.)
4. **Criterion 5 via the back door.** Maybe the "transfer unverified" label landed only in the
   unguarded skill copy: `grep -c` → **0** there too (and 0 in both guarded copies). Criterion 5
   fails in all three. This one broke the implementation, not the claim — Finding 4. (PRIMARY.)
5. **The implementer's HiddenBench read.** Fetched `arXiv:2505.11556` myself
   (`curl … | pdftotext -layout - -`): "The protocol consists of two stages. In the Exchange stage
   **(2 rounds)**, each agent shares 1–2 decision-relevant facts and gives one reason the current
   front-runner may be incorrect. In the **Decide stage (1 pass)**…"; "18 HIDDENBENCH tasks … 5
   runs each, **4 agents**"; Table 7 GPT-4.1 **0.037 → 0.800**; works "without explicitly
   informing agents of asymmetry". The implementer's PRIMARY read is accurate — and it confirms,
   not refutes, the D1 structural gap. (PRIMARY, CONFIRMED.)
6. **Retroactive rejection probe.** Ran the new section requirement over this idea's own six
   round-01 files: all six lack `## Existing alternatives` — the gate would reject the full
   pre-adoption corpus, including the late `opencode-1.md` that FINAL.md says implementation must
   read. Not a defect (the rule is not retroactive; that corpus is exactly the frozen baseline),
   and it makes R3's ritualisation comparison concrete: baseline compliance is 0/6. Probe fails
   to indicted the change. (PRIMARY.)
7. **"No consensus-side machinery to imitate."** Checked whether clause (c)'s missing scanner at
   least had no existing home: false — `internal/consensus/consensus.go:433ff` already has
   consensus.md schema machinery. The scanner had a place to live and wasn't built. Strengthens
   Finding 3 rather than breaking it. (PRIMARY.)
8. **Skill manifest freshness.** `shasum -a 256` on `SKILL.md` and `references/COOPERATION.md`
   matches `parley-addon.json` exactly — "manifest re-run after the payload edit" is true.
   (PRIMARY, CONFIRMED.) What did not survive: those changes are uncommitted — Finding 6.

## Findings

### [CRITICAL] The new §15.6(a) gate has no caller on the runtime path

`protocol.ValidateRoundOneArtifact` (`internal/protocol/roundartifact.go:23`) is invoked only by
its own test file. (PRIMARY: repo-wide grep for `ValidateRoundOneArtifact` and
`protocol.RoundOneRequiredSection` — every non-test hit is the *other*, pre-existing
`runner.ValidateRoundOneArtifact(path, agentID, ideaSlug)` at `internal/runner/validation.go:63`,
a different function in a different package with a different signature.) The actual runtime path
is `validateArtifactForPhase` (`internal/runner/phase58.go:312,321`) → `ValidateRoundArtifact`
(`validation.go:15`, round ≤ 1 branch) → `runner.ValidateRoundOneArtifact`, which still checks
only the old four headings (`## Summary`, `## Proposed approach`, `## Concerns / open questions`,
`## Risks`); the driver also calls it directly at `internal/driver/driver.go:367`. Neither path
references `## Existing alternatives`. (PRIMARY, CONFIRMED.)

Why it matters: this is the idea's own headline defect class — a rule whose validator exists but
binds nowhere — shipped inside the change that names it. Worse, the new §15.6 text asserts
*"The executing wording lives in the round prompt templates and is **validated there**"* — a
statement that is false at the reviewed commit. Acceptance criterion 2 is met only inside a test
harness; in production a round-1 artifact without the section completes round 1 and advances.
The duplicated function name across the two packages is an additional hazard: a future maintainer
can wire, or believe wired, the wrong one.

Fix: call the check from the runtime validator — e.g. in `runner.ValidateRoundOneArtifact`, after
the frontmatter checks, `if !protocol.HasNonEmptySection(body, protocol.RoundOneRequiredSection)`
— or delegate to `protocol.ValidateRoundOneArtifact` outright, and add a test that exercises the
runtime-path validator (not the orphaned one) on a missing-section artifact. (An uncommitted
working-tree edit doing roughly the first form appeared during this review; finding stands
against 59eb663.)

### [MAJOR] §15.7's per-track row contradicts the new §15.6

The 15.7 table's last row still reads `15.6 correlated agreement | no | yes (section in an
existing round-02 file) | yes (assigned round artifact)` (`parley-deck/COOPERATION.md:1372`,
unchanged from baseline — PRIMARY). It (a) denies binding on `fast`, while the new §15.6 preamble
says "Unconditional on every track" and FINAL.md leg 1 explicitly ratifies **all tracks including
`fast`**; (b) names mechanisms the change deleted ("section in an existing round-02 file",
"assigned round artifact"); (c) leaves clauses (a) and (c) with no row at all. The protocol is
self-contradictory in adjacent sections, and the track binding FINAL.md ratified for leg 1 is
negated two paragraphs later. Fix: replace the row with the unconditional form (`15.6
alternatives & correlated agreement | yes | yes | yes`) in both guarded copies. (A matching
uncommitted edit to the defaults copy appeared mid-review; the deck copy must move with it or
`TestEmbeddedDefaultMatchesLiveDeck` fails closed.)

### [MAJOR] Clause (c) prints an enforcement promise with no scanner — D4's principle applied inconsistently

§15.6(c) as shipped: "*a contradiction blocks signoff and escalates to the owner. The scanner
never auto-halts.*" There is no scanner: zero Go matches for `Alternatives disposition` or
`ALT-` (PRIMARY, CONFIRMED), `internal/consensus/consensus.go` untouched by the diff, and
acceptance criterion 6 ("Both halves tested") has no tests. The implementer withheld clause (b)
precisely because printing an unenforced rule would be "the fourth instance of this deck's named
defect class" — then printed clause (c) in exactly that state. The principle was right; its
selective application shipped the defect one clause over. Fix: either build the consensus-side
check (the schema machinery at `internal/consensus/consensus.go:433ff` is the natural home) with
both halves of criterion 6 tested, or strip the enforcement sentence from (c) until the scanner
exists — under the implementer's own D4 logic, those are the only two consistent options.

### [MAJOR] Ruling on D4 — withholding was substantively right and procedurally an overstep; regularize it

Two questions were put to review. On the merits: **the call was correct.** Printing an exchange
duty with no exchange stage would have been the fourth printed-rule-that-binds-nowhere, committed
by the idea that named the class (premise verified: zero exchange code — Refutation 3). On
authority: **it was an overstep as executed.** FINAL.md ratified leg 2 *and* acceptance criteria
3–5 that reference it; an implementer does not get to un-ratify a clause, only to propose. The
evidence the criteria are now void is mechanical: criterion 5 requires the string "transfer
unverified" in `COOPERATION.md`; it occurs 0 times in all three copies (PRIMARY, Refutation 4) —
a ratified, grep-checkable criterion fails at the reviewed commit, and the deferral is recorded
only in `IMPLEMENTATION.md`, while FINAL.md's own follow-up list (1–4) contains no leg-2
implementation idea. The byte consequence (−476 B instead of −237 B) is real and verified
(Refutation 1).

Ruling: **sustain the withholding; require regularization.** Amend FINAL.md (or file the
successor) to record the leg-2 deferral, re-scope criteria 3–5 to the leg-2 implementation idea,
and name its owner — and apply the same rule to clause (c) per Finding 3. With that amendment,
D4 becomes a disclosed, ratified-scope change instead of a silent one.

### [MAJOR] Ruling on D1 — the implementer was right to build the ratified form; the design record now miscites its own foundation

Was it right to implement one sealed packet rather than upgrade to the measured two rounds +
Decide on post-ratification evidence? **Yes.** Implementing the ratified text and escalating the
discrepancy for review is the correct division of labor; silently re-designing a ratified leg on
an implementer's own primary read would be the larger procedural defect, and the implementer did
the R2 gate honestly. The defect is therefore not in the implementation but in the ratified
record: FINAL.md leg 2 cites +76.3pp from a protocol whose structure it does not adopt — the
measured intervention is *Exchange (2 rounds) + Decide (1 pass) at 4 agents* (verified by my own
extraction — Refutation 5, PRIMARY, CONFIRMED), not one packet with no Decide pass. The
"structurally derived; transfer unverified" label honestly covers unknown transfer; it does not
cover *known-different* structure, which is what the record now holds. Because D4 deferred leg 2
entirely, no mis-shaped code shipped, so this does not block the current change. Ruling: **do not
block on D1; gate it forward.** The leg-2 implementation idea's first decision must be one-packet
vs two-rounds-plus-Decide with the paper's structure quoted, and FINAL.md's citation must carry
the structural-difference qualifier until then. Shipping the one-packet runner stage while citing
the two-round effect size would convert an honest label into the mis-citation class of R6.

### [MINOR] The third copy's move is uncommitted in its own repo

`../parley-deck-skill` carries the §15.6 change, the SKILL.md template update, and the refreshed
manifest as **uncommitted working-tree modifications** (`git -C ../parley-deck-skill status` →
` M skills/parley-deck/SKILL.md`, ` M skills/parley-deck/parley-addon.json`,
` M skills/parley-deck/references/COOPERATION.md`; PRIMARY). Content and hashes are verified
consistent today (Refutations 1, 8), but the copy no drift guard covers is also the copy no
commit pins: one stray `git checkout --` in that repo silently reverts it while the two guarded
copies stay new. Fix: commit the skill-side change before this idea is called complete.

### [MINOR] D2's R4 wording correction lives nowhere binding

D2 prescribes softening R4's "+0.6% at 7 agents against +34.8% at 3" indictment, and the
implementer's Table-4 read (non-monotonic; N=6 beats N=5; 5–7 tasks per cell except N=4) is
correct per the paper. But `git diff --stat 9d4f45c..59eb663` touches no FINAL.md (PRIMARY), so
the ratified text still asserts the monotonic reading; the correction exists only in
IMPLEMENTATION.md. Fold the R4 wording fix into the same FINAL.md amendment as Finding 4.

### [NIT] IMPLEMENTATION.md checklist left unchecked

`- [ ] **Files or areas to change:**` remains open while the frontmatter says
`status: implemented`. Cosmetic, but this deck measures the gap between stated and actual state;
the file should not exhibit it.

## Open questions

1. Who owns the leg-2 implementation idea, and is the one-packet vs two-rounds-plus-Decide ruling
   its gate-zero item, before any runner code? (D1's forward path.)
2. What is the amendment mechanism for a ratified FINAL.md — in-place amendment vs successor
   idea? D4's deferral, the criteria 3–5 re-scoping, and the R4 wording all need whichever form
   the protocol recognizes; I found no precedent cited either way.
3. Disposition vs REFRAME: leg 3 absorbs alternatives into `## Alternatives disposition` in
   consensus.md (signoff block); `round-01/opencode-1.md` absorbs frames into `## Frames
   considered` inside FINAL.md (freeze gate, "silence is not rejection"). One carrier or two with
   a stated boundary? The deferred vocabulary follow-up must answer this explicitly; the two
   destinations will otherwise compete. On item 5 from the review brief: the disposition leg's
   *ratified content* is complete without REFRAME — FINAL.md already ruled "no new finding-class
   vocabulary" for this design — so reading opencode-1 is not what the leg is missing; the
   scanner is (Finding 3).
4. A concurrent fix-up touched `internal/` during this review round. Which commit closes
   round-01 review — 59eb663 as reviewed, or a fix-up commit the reviewers then re-review?
5. Until the clause-(c) scanner exists, is the disposition duty prose-only — i.e., did instance
   #4 of the defect class ship at 59eb663 anyway, one clause to the left of where the implementer
   refused to ship it?
