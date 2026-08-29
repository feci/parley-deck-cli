---
agent: zcode-1
idea: protocol-generation-bias
review-round: 1
date: 2026-08-29
reviewed-commit: 59eb663
---

## Summary

I verified every claim I could execute myself, ruled on both flagged deviations, and found one
thing the implementer did not flag: **the new §15.6(a) gate is dead code.**
`protocol.ValidateRoundOneArtifact` (`internal/protocol/roundartifact.go:23`) has zero non-test
callers; the runtime path (`internal/driver/driver.go:367`, `internal/runner/phase58.go:321` →
`runner.ValidateRoundArtifact` → the runner's own same-named `ValidateRoundOneArtifact`,
`internal/runner/validation.go:63`) still checks only the four legacy sections. I fed it a round-1
artifact with no `## Existing alternatives` — it passed with `err=nil`. FINAL.md acceptance
criterion 2 is unmet in the executed system, and this is precisely the idea's own defect class: a
validator that exists and is never invoked. The mutation test the implementer ran proves the
function rejects when called; it proves nothing calls it.

Everything else checks out mechanically: the byte claim re-derives exactly (1,372 B → 896 B, net
−476 B, `PRIMARY`), all three `COOPERATION.md` copies moved and are byte-identical in §15.6
(`PRIMARY`), the gate's own test fails under my own mutation and the full suite is green
(26 packages, exit 0, `PRIMARY`), and the HiddenBench quotes the implementer extracted are
verbatim in the arXiv HTML v4, which I fetched myself (`PRIMARY`) — the D1 factual basis is sound.

**Ruling on D1 (exchange fidelity):** the implementer was right to build the ratified one-packet
form rather than silently upgrade a ratified design on post-ratification evidence — an agent
proposes, it does not amend ratification. But the ruling this review must now record is that the
**follow-up exchange idea implements the measured form (two exchange rounds plus a Decide pass) or
cites no effect size at all.** +76.3pp belongs to a structure we would not be building; a
one-packet form is an unmeasured interpolation, and FINAL.md R1 currently misattributes the number.

**Ruling on D4 (clause (b) withheld):** the withholding is **correct in substance and irregular in
form.** Printing a ratified duty no code performs would be the fourth instance of this deck's named
defect class, committed by the idea that named it — and the new §15.6's own preamble
("the executing wording lives in the round prompt templates and is validated there") would have
been falsified by its own clause (b) in the same section. But an implementer cannot both cite
FINAL.md as authority and delete a ratified leg without a recorded re-scope: acceptance criteria 3,
4 and 5 now silently lapse, and "transfer unverified" appears nowhere in the protocol (I grepped
all three copies — zero hits). Endorsed, conditional on this review's consensus recording the
re-scope and on the follow-up idea existing before 2.11.0 publishes.

Two further majors the implementer disclosed only partially: leg 3 shipped prose-only (no template,
no scanner, no SKILL.md carriage — the §15.6 preamble's "validated there" is false for clause (c)),
and the skill-repo half of the implementation sits **uncommitted** in `../parley-deck-skill`'s
working tree, undisclosed in IMPLEMENTATION.md.

## Refutation attempts

I assumed the implementation was wrong and tried to break each verified claim. Attempts that
failed to break it are recorded here with what I ran.

1. **Tried to break the byte accounting.** Re-derived at the baseline tag:
   `git show protocol-generation-bias-baseline:parley-deck/COOPERATION.md | sed -n '1346,1368p' | wc -c`
   → **1372**, matching FINAL.md's one-command check. New §15.6 at `parley-deck/COOPERATION.md:1346-1361`
   → `wc -c` → **896**. Net −476 B. Could not break it. (Char-count cross-check: defaults copy is
   104,413 via `wc -m` — the implementer's number is chars, not bytes; offsets +246/+90 hold in bytes.)
2. **Tried to find drift among the three COOPERATION.md copies.** Extracted §15.6 from
   `internal/protocol/defaults/COOPERATION.md:1339`, `parley-deck/COOPERATION.md:1346`, and
   `../parley-deck-skill/skills/parley-deck/references/COOPERATION.md:1339` and diffed all three:
   byte-identical. The old carve-out string *"primarily a judgment rather than a mechanically
   decidable"* is present in all three at baseline and absent in all three now. No drift. (But see
   the finding on commit state — the skill copy's movement is not in any commit.)
3. **Tried to defeat the gate's test by mutation, myself.** Changed the condition in
   `internal/protocol/roundartifact.go:28` to `if len(raw) < 0` (never fires, still compiles):
   `TestValidateRoundOneArtifactRequiresExistingAlternatives` **failed**; restored it; full suite
   `go test ./... -count=1` → exit 0, 26 packages ok; drift guard
   `TestDriftAnchorsAcceptTheGeneratedRosterTable` passes. The test is genuinely sensitive to the
   gate. Could not break it.
4. **Tried to falsify the implementer's HiddenBench extraction as fabricated quotes.** Fetched
   `arxiv.org/html/2505.11556v4` myself: *"In the Exchange stage (2 rounds), each agent shares 1–2
   decision-relevant facts and gives one reason the current front-runner may be incorrect"*);
   *"In the Decide stage (1 pass), each agent summarizes the strongest evidence and remaining
   uncertainty before voting"*; evaluation at *"5 runs each, 4 agents"*; protocol works *"without
   explicitly informing agents of asymmetry"*. All verbatim as quoted in IMPLEMENTATION.md. The
   quotes survive; D1's factual basis is real. (`PRIMARY` — my fetch, my locator.)
5. **Tried to find exchange code that would falsify D4's premise.** `grep -rn "Evidence
   exchange|sealed packet|exchange stage" --include="*.go" internal/ cmd/` → no matches. The runner
   genuinely has no exchange stage; D4's premise holds.
6. **Tried to find a fast-track bypass of leg 1.** The round-1 prompt dispatch
   (`internal/runner/runner.go:929`) has a single `BuildRoundOnePrompt` with no track branch, so the
   enumerated instruction reaches every track including `fast`, as ratified. Could not break it.
7. **Tried to pass a §15.6(a)-noncompliant artifact through the RUNTIME gate — and it passed.**
   Temporary test in `internal/runner` (written, run, deleted): a fully-compliant round-1 artifact
   (frontmatter + all four legacy sections non-empty) with no `## Existing alternatives` →
   `runner.ValidateRoundOneArtifact(path, "hermes", "demo")` returned `err=nil`. This refutation
   attempt **succeeded** and is the CRITICAL finding below.
8. **Tried to find any other caller of the new gate** (a CLI subcommand, a consensus hook):
   `grep -rn "ValidateRoundOneArtifact|RoundOneRequiredSection" --include="*.go"` over the whole
   module — the protocol-package function is referenced only by its own test file. Confirms finding 1.

## Findings

### [CRITICAL] The §15.6(a) gate is dead code — the runtime validator never enforces it

**What is wrong.** `protocol.ValidateRoundOneArtifact` (`internal/protocol/roundartifact.go:23`)
is called by nothing except its own test. The runtime path for design round-1 artifacts is
`internal/driver/driver.go:367` and `internal/runner/phase58.go:321` → `runner.ValidateRoundArtifact`
→ (round ≤ 1) the runner's own `ValidateRoundOneArtifact` (`internal/runner/validation.go:63`),
whose required-section list is still `## Summary`, `## Proposed approach`,
`## Concerns / open questions`, `## Risks` — the new section was never added. Demonstrated
empirically: an artifact lacking `## Existing alternatives` passes the runtime gate (`PRIMARY`,
temp test executed and deleted). Worse, the runner and protocol packages now export two
same-named functions with different signatures, so grepping "ValidateRoundOneArtifact" finds a
live-looking caller that is not this gate.

**Why it matters.** FINAL.md acceptance criterion 2 requires rejection *"by the validator, not
merely warned about."* As shipped, the duty is prompt-carried only — and the idea's own headline
finding (kimi-1's refinement, quoted in FINAL.md) is that *"the half a scanner cannot check decays
to prose rates even when the prompt carries it."* The implementation of the idea that named this
defect class shipped a second instance of it (after successfully avoiding the first in D4), and the
mutation test in IMPLEMENTATION.md masks it: watching the function fail when called says nothing
about whether anything calls it.

**Fix.** In `runner.ValidateRoundOneArtifact`, add `protocol.RoundOneRequiredSection` to the
required-sections loop (reusing the existing `sectionHasContent`), or delegate to
`protocol.ValidateRoundOneArtifact(path)` after the frontmatter checks. Add a runner-level test
that a round-1 artifact missing the section fails *through* `ValidateRoundArtifact(..., round=1)`.
Rename one of the two colliding functions (suggest: the runner's becomes
`validateRoundOneContract`). Two lines plus a test; this is a block, not a fast-follow — the idea's
thesis is the stakes.

### [MAJOR] Leg 3 disposition is prose-only, and §15.6's preamble asserts carriage that does not exist

**What is wrong.** `## Alternatives disposition` appears in exactly one place in the three
repositories: the protocol prose (`internal/protocol/defaults/COOPERATION.md:1351`). No Go
template carries it, no validator checks it, and `SKILL.md`'s Phase 3 consensus template was not
updated (the skill diff touches only the round-1 template). The new §15.6 opens with *"The
executing wording lives in the round prompt templates and is validated there; this section carries
the duty only"* — for clause (c) that sentence is false in both halves.

**Why it matters.** By this idea's own carrier thesis, a prose-only duty runs in single digits.
FINAL.md acceptance criterion 6 (contradiction blocks signoff, escalates, scanner never auto-halts,
both halves tested) is unmet — IMPLEMENTATION.md admits the scanner is "remaining code work," but
does not admit the template half is also missing, and the protocol text now makes a false claim
about its own enforcement, which is worse than an honest gap. Clause (c) is the B1 fix — the case
that started this idea (the withdrawn PBS proposal vanishing from FINAL.md) is exactly a
disposition failure, and it got the single-digit-carrier form.

**Fix.** Carry `## Alternatives disposition` in the consensus drafter prompt (Go) and the SKILL.md
Phase 3 template, and shape-check it in `internal/consensus` (every `ALT-` id raised anywhere in
the idea appears with adopt/reject; contradiction detection may stay human, per "the scanner never
auto-halts"). Alternatively amend §15.6(c) to state the scanner lands with the follow-up idea — but
then the preamble sentence must be qualified. Gate the 2.11.0 publish on one of the two.

### [MAJOR] The skill-repo half of the implementation is uncommitted and undisclosed

**What is wrong.** The reviewed commit 59eb663 contains only the CLI-repo surfaces. The skill
changes — `skills/parley-deck/references/COOPERATION.md`, `skills/parley-deck/SKILL.md`,
`skills/parley-deck/parley-addon.json` — are uncommitted working-tree modifications in
`../parley-deck-skill` (HEAD `3a16a43`, "release: parley-deck-skill 2.10.0"). IMPLEMENTATION.md
reports "All three `COOPERATION.md` copies moved" and "SKILL.md round-1 template updated" without
disclosing that the skill half is in no commit.

**Why it matters.** A reviewer or reproducer checking out the reviewed commit gets a skill copy
stale at 2.10.0 against a deck protocol carrying 2.11.0 content — the exact three-copy drift the
implementer itself warned about in its own Notes ("The skill's references/COOPERATION.md has gone
stale before"). The record of this implementation cannot be verified from the record; it is
verifiable only from this machine's working tree right now.

**Fix.** Commit the skill repo on its own cadence and reference that commit hash in
IMPLEMENTATION.md; or, if the skill release is deliberately deferred, disclose the pending state
and the owner action the way D3 discloses the publish handoff.

### [MAJOR] The deferred vocabulary question was deferred again, silently

**What is wrong.** FINAL.md binds the implementation: *"implementation must read
`round-01/opencode-1.md` before deciding it"* (the `REFRAME` class / `## Frames considered`
question). IMPLEMENTATION.md's only discharge is Notes item 5, which tells **reviewers** to read it.
No decision was made, no deviation entry exists, and the shipped §15.6(c) engages with none of it.

**Why it matters.** I read `round-01/opencode-1.md` in full. Its proposal overlaps leg 3's cargo
substantially (absorb-before-freeze, stable identifiers, adopt/reject-with-reason) and diverges on
three load-bearing points: destination (`## Frames considered` inside `FINAL.md` vs
`## Alternatives disposition` inside `consensus.md`), checkability (a one-string presence gate —
"FINAL.md MUST NOT freeze while any REFRAME ... is absent" — vs leg 3's semantic
contradiction check that nothing scans), and the witness field (a locator per entry, the explicit
anti-theater rule; leg 3's ALT- entries require only "the decisive reason"). Two half-formed,
mutually ignorant schemes for the same cargo guarantee a collision in whichever follow-up lands
second. Note also that opencode-1's Phase 5–8 rule — "IMPLEMENTATION.md MAY NOT introduce a frame
absent from `## Frames considered`" — bears directly on how D4 was executed: an implementer
withholding a ratified frame decision is the shape that rule exists to catch. The disposition leg
is **not complete** without a recorded decision here.

**Fix.** This review round's consensus records the decision, minimally: ALT- identifiers subsume
the REFRAME identity question (or do not, and `REFRAME` is added as a parallel class per
opencode-1, keeping it out of the `impl.go:445` severity switch); the witness field is merged into
ALT- entries (recommended — it is mechanically checkable and is the anti-ritual rule §15.6(c)
currently lacks); and the freeze-gate-vs-contradiction-check choice is made once, for both schemes.

### [MAJOR] D1: the design record still attaches +76.3pp to a structure that was never measured

**What is wrong.** Nothing exchange-shaped was committed, so nothing in the code is wrong. What is
wrong is the record FINAL.md leaves as the spec for the follow-up: leg 2 and R1 cite HiddenBench's
+76.3pp next to a one-packet, no-Decide design, and the primary source (verified by my own fetch)
measured two exchange rounds plus a Decide pass at 4 agents.

**Why it matters.** This is a provenance defect in the deck's own §15.2 terms: a headline number
attributed to an intervention other than the one specified. Any future reader building leg 2 from
FINAL.md alone inherits the misattribution, and the ratified "transfer unverified" label becomes an
understatement — the transfer would be unmeasured *and* structurally truncated.

**Fix.** Record the ruling in review consensus: the follow-up implements the measured form (two
exchange rounds plus a Decide pass, 4-agent binding), or implements one packet as an explicitly
unmeasured variant with its own instrumentation and **cites no effect size**. Amend R1's wording
when the follow-up opens. The implementer's restraint in not silently upgrading the ratified design
was correct; the corollary is that the upgrade must now happen on the record, here.

### [MINOR] D4: withholding clause (b) is endorsed, but the re-scope of FINAL.md is unrecorded

**What is wrong.** The withholding itself I endorse (see Summary). What is missing: FINAL.md's
acceptance criteria 3, 4 and 5 (exchange adds zero rounds/files/agents; no asymmetry assertion in
the execution prompt; *"transfer unverified"* in `COOPERATION.md`) are now unmet with no recorded
disposition — I verified the label string appears in none of the three copies and in no Go file.
An implementer cannot produce that state from a ratified package without a destination for the
delta; that is precisely the "contradicted adoption" shape leg 3 exists to catch, applied to this
idea's own FINAL.md.

**Why it matters.** The alternative — printing an unenforced clause — is worse, and the
implementer's reasoning (the carrier thesis applied to its own protocol text) is the correct
application of the idea's thesis. But an endorsed deviation and a silent lapse differ in whether
the next agent can reconstruct why the protocol jumped from three legs to two.

**Fix.** Review consensus explicitly re-scopes criteria 3–5 to the follow-up exchange idea; that
idea re-adopts them plus the D1 ruling; it opens before `parley protocol publish --version 2.11.0`
runs, so the owner never publishes a core whose ratification record and text have diverged with no
pointer.

### [MINOR] D3 overstates readiness: no 2.11.0 artifact exists to byte-verify

**What is wrong.** IMPLEMENTATION.md D3 says *"The 2.11.0 core is prepared and byte-verified here."*
`~/.parley/protocol/core/` contains only `2.10.0`; a filesystem-wide search finds no `2.11.0`
directory anywhere. What exists is the new body in `internal/protocol/defaults/COOPERATION.md`,
which `parley protocol publish --version 2.11.0 --from <file>` (attended, TTY-gated — verified at
`internal/app/protocol.go:361`) would consume.

**Why it matters.** "Prepared and byte-verified" reads as an artifact a verifier can hash; there is
none. The handoff to the owner is right in spirit — publish is correctly attended-only — but the
owner is not told the exact command or which file is the body.

**Fix.** State the exact invocation and body path in IMPLEMENTATION.md (and note whether the
defaults copy is publish-shaped or needs a wrapper), or stage the intended body at a durable path
so "byte-verified" is checkable by anyone.

### [NIT] IMPLEMENTATION.md mixes character and byte counts for the copy sizes

**What is wrong.** "defaults 104,413 chars; deck +246; skill +90" — 104,413 is `wc -m` characters;
the offsets hold in bytes (`wc -c`: 105,415 / 105,661 / 105,505). The skill offset is +89 chars but
+90 bytes.

**Why it matters.** Minor, but this idea's record was already burned once by a number inherited
across measurement methods (the −726 B error FINAL.md corrects); state the unit once and keep it.

**Fix.** One sentence: "sizes in bytes via `wc -c`" with the byte figures.

## Open questions

1. **Is the dead-gate CRITICAL a block for this commit or a fast-follow?** I recommend block: the
   fix is two lines plus a test, and shipping an implementation of *this* idea with a
   never-invoked validator is the one outcome the idea cannot survive.
2. **Who records the D1/D4 dispositions with authority** — this review's consensus, a FINAL.md
   addendum, or the follow-up idea's opening state? FINAL.md's own leg 3 says a contradicted
   adoption escalates to the owner; the cleanest form here is review consensus plus owner
   dismissal-with-reason or ratification of the re-scope.
3. **Witness field into ALT- entries?** opencode-1's four-field payload (current frame / other
   frame / witness / stay-condition) is the only anti-ritual mechanism either scheme proposes for
   the disposition section; §15.6(c) as shipped has none. Recommended for the follow-up, decided
   now (see the vocabulary finding).
4. **Arrival date of record for `opencode-1.md`:** its frontmatter says `date: 2026-08-28` while
   FINAL.md says it was filed at 14:05 on 2026-08-29, after exclusion. Which is authoritative for
   the §9.0 record? (No action for this implementation; a record-accuracy question for the deck.)
5. **The dose-response non-monotonicity (D2)** — I did not re-verify Table 4 cell-by-cell against
   the paper; my fetch confirmed the endpoints and the "4 agents" evaluation condition, which is
   what Gate R2 needed. If R4's softened wording is adopted, someone should re-derive the N=5/6/7
   cells from the PDF before the 88-idea comparison runs.
