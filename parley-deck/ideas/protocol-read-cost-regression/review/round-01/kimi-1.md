---
agent: kimi-1
idea: protocol-read-cost-regression
round: 1
date: 2026-08-10
---

## Summary

The design-side wiring is sound: round 2 is byte-identical to the old path, the fallback renderer is
the old walker verbatim, the boundary is honoured in every validator and in
`buildConsensusDraftPrompt` (verified untouched in the diff — G4 holds design-side), and `go build`,
`go vet`, `go test ./internal/runner/` are all green on my own run. But the review-side wiring has
two CRITICAL defects the declared deviations do not name, and the fallback net catches almost none
of the uncertainty states my signoff contract enumerated. The two flagged deviations are real; they
are not the worst of what I found. This cannot ship as-is.

## Refutation attempts

- Re-ran `go build ./...`, `go vet ./internal/runner/`, `go test ./internal/runner/` — green (9.3s).
- Diffed `internal/runner/runner.go` and `internal/runner/phase58.go` against HEAD — the change is
  exactly the frontier call, the instruction edit, and the `gatherReviewContext` split; nothing else.
- Confirmed `internal/app/driver_consensus.go` absent from the diff and
  `buildConsensusDraftPrompt` (`driver_consensus.go:113`) still orders the full-history read.
- Traced `buildPromptForRound` case `"review-consensus"` (`runner.go:919-924`) →
  `gatherReviewContext(path, round+1)` → `frontierContext` — the Phase 7 drafter compacts (Finding 1).
- Traced `gatherReviewContext(path, 1)` and `(path, 2)` by hand — head is prepended to a `full()`
  result that already contains it: FINAL.md/IMPLEMENTATION.md embedded twice (Finding 3).
- Constructed the head-drop case: a review round whose earlier rounds contain the literal phrase
  "carry-forward fallback" (this review file is such an artifact) makes
  `strings.Contains(rounds, "carry-forward fallback")` true on the non-fallback path → head skipped
  (Finding 4).
- Enumerated the round-2+ template's own sections ("## Remaining disagreements", "## Refined
  position") and the round-1/review templates ("## Concerns / open questions", "## Open questions")
  against `ledgerMarkers` — their natural wording carries no marker (Finding 2).
- Checked each G3 trigger against `buildCarryForward`/`frontierContext` line by line (Finding 5).

## The six questions I was asked to hit hardest

**1. Derived ledger, no lifecycle — does it block?** It blocks as built; the missing lifecycle does
not. Concretely: a ledger with no close operation *cannot fabricate disposal* — no item can be
RESOLVED or SUPERSEDED by anyone, owner or not, so the owner-disposes rule cannot be violated
inside the ledger. The cost of no-lifecycle is stale items accumulating forever (the window is
always rounds 1..N−2 from r=1 — nothing ever "falls out of scope", contrary to IMPLEMENTATION.md's
claim), which is my signoff's false-*non*-convergence direction: wasted rounds, not false consent.
That direction does not block v1. What blocks is that derivation-by-markers recreates orphaned
dissent through the extraction gap (Finding 2) and misattributes ownership (Finding 6), and the
fallback net does not catch either (Finding 5). The contract pieces that may legitimately wait for
v2: immutable IDs, SHA-256, append-only transition history, DISPUTED joining (tracked separately by
G6) — they serve audit and re-litigation detection, not consent integrity. Verdict: DERIVED is
acceptable for v1 only after Findings 2, 5 and 6 land; before that, no.

**2. G6 knowingly unmet — severity?** MAJOR, and it keeps the idea open. The degradation is
partially fail-safe: both verdict lines arrive verbatim and the extractor never invents a
resolution — which the contract forbids — so this is not a silence-into-consent conversion. But a
conflict that reaches a round as two unrelated lines is not "joined as DISPUTED" and never triggers
fallback, so G6 and the last trigger of G3 both stay red. Shipping with a named gate knowingly
unmet is the owner's call to record, not the implementer's to declare done; IMPLEMENTATION.md
handled the honesty correctly and the review must not convert that honesty into acceptance. Minimum
for v1: the ledger header states conflicts are NOT joined and opposing verdict lines must be read
as DISPUTED; the G6 fixture exists as a failing-then-passing test. Finding 7.

**3. The boundary.** Honoured everywhere except one place, and that place is CRITICAL. No
validator, no close condition, no artifact-acceptance rule references the ledger; selection is
prompt-assembly only; the fallback announces itself in the prompt, keeping the optimization
falsifiable. The one breach is in effect, not in mechanism: the Phase 7 review-consensus drafter —
the function that decides `outstanding_agreed_fixes` and writes "## Dismissed findings" — now
receives compacted context (Finding 1). A consensus-close decision running on a marker-filtered
history is exactly the artifact-validity-adjacent use the boundary forbids.

**4. Fallback completeness.** Enumerated in Finding 5: of the six required states, only "globally
zero extracted items" and "empty previous round" fire. Missing (partial), invalid, ambiguous,
challenged, unresolved hash/locator, and verdict-conflict-not-DISPUTED do not.

**5. Silence = implicit agreement.** Two paths convert omission into recorded consent: the marker
gap (Finding 2) and the Phase 7 drafter compaction (Finding 1). Both are CRITICAL per the rule, not
MAJOR.

**6. What the marker list misses.** Enumerated in Finding 2 — including the round-2+ template's own
"## Remaining disagreements" prose, "Objection:", "I object", lowercase severities, "verdict:
RESERVED", and every NIT finding.

**Does it break `parley run`?** Design round 2: no — `frontierContext(round≤2)` returns `full()`,
byte-identical to the old `gatherPriorRounds`; the instruction edit does not touch validation.
Review round 1: runs and validates, but FINAL/IMPLEMENTATION are embedded twice (Finding 3) — a
regression, not a breakage. Design consensus drafter: untouched, verified. Review-consensus
drafter: compacted — Finding 1, CRITICAL.

## Findings

### [CRITICAL] The Phase 7 review-consensus drafter now receives compacted history

`buildPromptForRound` case `"review-consensus"` (`internal/runner/runner.go:919-924`) calls
`gatherReviewContext(opts.Idea.Path, roundNumber(opts)+1)`, which since this change routes through
`frontierContext`. After review round 2 (context round 3), review round 1 reaches the drafter only
as marker-matched lines; with the measured 5.1-round average the drafter almost always compacts.
The drafter's job is to synthesize **every** finding into agreed fixes and dismissed findings and
to set `outstanding_agreed_fixes` — the count that closes the auto fix-up loop. A finding from an
earlier review round whose wording matches no marker (see Finding 2) leaves the drafter's context
without its owner disposing of it; the loop can then close with the fix never applied. That is
silence converted into recorded consensus at the exact point the protocol closes review. FINAL.md's
table says "Consensus drafter — unchanged, full history"; G4 names `buildConsensusDraftPrompt`, and
that function is indeed untouched — but the review-side drafter is a consensus drafter too, and it
was made *worse than today* without being named in IMPLEMENTATION.md. This is an undeclared third
deviation and a boundary breach in effect. **Fix:** in the `"review-consensus"` case call
`gatherReviewContextFull(opts.Idea.Path, roundNumber(opts)+1)`; keep frontier selection for the
`"review"` case only. Add a test asserting the review-consensus context contains every prior review
round body.

### [CRITICAL] Marker extraction is fail-open: objections without a marker silently drop out (G5 broken in the general case)

`extractItems` (`internal/runner/frontier.go:120-141`) carries a line only if it contains one of 16
fixed case-sensitive substrings. Anything else from rounds ≤ N−2 never enters the ledger, no
fallback fires, and Phase 2 rule 1 records the omission as agreement. The G5 fixture tests only the
easy case (a literal "CRITICAL:"). What the list misses in real deliberation wording:

- **The round-2+ template's own output.** "## Remaining disagreements" and "## Refined position"
  are free prose: "- The retry cap is still unsafe under timeout" contains no marker. The template
  the runtime itself mandates produces the exact lines the extractor drops.
- "Objection:", "I object", "we disagree", "disagreement", "I don't agree" (misses "I disagree").
- Lowercase or prose severities: "critical: the lock is unchecked", "this is a major problem" —
  matching is case-sensitive.
- "verdict: RESERVED" — a reserved signoff *is* a conditional objection; BLOCK is carried, RESERVED
  is not.
- "## Concerns / open questions" (round-1 template) and "## Open questions" (review template)
  prose — where unresolved challenges naturally live.
- "This breaks/regresses X", "still unaddressed", "do not ship", "I remain unconvinced".
- Every NIT finding — the never-cut list says "every live objection or finding" and no severity
  floor was adopted anywhere in FINAL.md.

**Fix:** stop keying on severity substrings. Carry the designated disagreement-bearing sections
verbatim per file ("## Remaining disagreements", "## Concerns / open questions", "## Risks",
"## Findings", "## Open questions", "## Signoffs"), and make extraction fail-closed per file: any
non-trivial participant artifact from which zero items were extracted either forces the fallback or
is carried in full. Add the adversarial G5 fixture: a round-1 objection phrased with no marker must
reach round 3 or trip the fallback.

### [MAJOR] Review rounds 1–2 embed FINAL.md and IMPLEMENTATION.md twice

`gatherReviewContext` (`internal/runner/phase58.go:486-507`) always prepends the head unless the
frontier result contains "carry-forward fallback". For round ≤ 2 `frontierContext` returns
`full()` — i.e. `gatherReviewContextFull`, which already emits FINAL/IMPLEMENTATION — so every
review round 1 and 2 prompt now carries the two largest head artifacts twice. Review round 1 is
therefore **not** unchanged from today, and the change's first effect on the most common review
call is to double the part this idea exists to shrink. No information is lost, which is why this is
MAJOR not CRITICAL. **Fix:** prepend the head only when `round >= 3` and the frontier did not fall
back (structured flag, per Finding 4), or have the round≤2 path return the rounds-only renderer.

### [MAJOR] Head-drop via substring sniffing: quoting the banner removes FINAL/IMPLEMENTATION from a reviewer prompt

The dedup test `strings.Contains(rounds, "carry-forward fallback")` (`phase58.go:503`) runs over
content that includes verbatim lines extracted from earlier rounds. Any earlier-round line
containing that phrase — a reviewer quoting the banner, discussing the mechanism, or reviewing
*this very idea* (this file contains the phrase repeatedly) — makes the check true on the
non-fallback path, so the head is skipped and the reviewer is told to "review against FINAL.md and
IMPLEMENTATION.md (below)" with neither present. Silent, content-triggered, and near-certain on
this idea's own review round 3. It violates the non-negotiable head (my reservation 1) and the
never-cut list (FINAL.md, acceptance criteria, check results). Borders CRITICAL for review quality;
MAJOR because the artifacts remain on disk and a diligent reviewer can recover. **Fix:** return a
struct from `frontierContext` (`text string; fellBack bool`) and branch on the flag, never on
content; or at minimum `strings.HasPrefix(rounds, "\n===== FULL HISTORY (carry-forward fallback)")`.

### [MAJOR] G3 is mostly unmet: four of six required fallback triggers do not exist

Required triggers and what the code actually does:

- **missing** — only the global version: zero items across all older rounds, or an empty previous
  round. A *partially* missing round is silent: `os.IsNotExist(err) → continue`
  (`frontier.go:87-89`) skips a missing round dir, so if round-01 is absent but round-02 yields one
  marker line, the ledger ships with round 1 gone and no fallback.
- **invalid** — nothing is validated, so nothing is rejected; a truncated or garbage artifact is
  carried or skipped without a trigger. Per-file zero-extraction is not a trigger either.
- **ambiguous** — no check (e.g. two files for one owner in one round).
- **challenged** — the contract's "citation/provenance challenge expands the source verbatim" is
  not implemented; no fallback.
- **unresolved hash/locator** — the ledger records no hashes, so there is nothing to fail to
  resolve; the trigger is vacuous rather than implemented.
- **verdict conflict not marked DISPUTED** — not detected, per Finding 7.

Separately: non-IsNotExist I/O errors abort the prompt build instead of falling back — visible and
fail-closed, so acceptable, but it is not what G3's wording promises. **Fix:** per-file
zero-extraction on a non-trivial file triggers fallback; a missing round dir inside the 1..N−2
sequence triggers fallback; the ledger header prints a coverage line (rounds and files read) so
gaps are visible to the reader; verdict-conflict detection per Finding 7. Extend the G3 test to one
case per trigger.

### [MAJOR] Owner attribution by filename misattributes quoted objections

`extractItems` attributes every matching line to the artifact's filename owner
(`frontier.go:44,102`). Deliberation files quote each other constantly — "> @b-1 wrote: CRITICAL:
…", "@codex-1's MAJOR is …". A quoted objection is rendered under the *quoter*, and the ledger
header then tells the reader an objection "is disposed of only by its own owner" — pointing at the
wrong agent. The true owner's disposition can never be matched to the item, and a non-owner appears
to hold it. That corrupts the one rule the ledger exists to protect. **Fix:** skip blockquote
("> ") lines or, on lines containing "@<agent>", attribute to the mention and mark the line
"quoted-by <file owner>"; carry the raw line prefix so quote context survives rendering.

### [MAJOR] G6 knowingly unmet: not DISPUTED-joined, never falls back — keep the gate red

As declared by the implementer. My judgment: the failure is fail-safe in direction (both verdict
lines arrive verbatim; no resolution is invented, which the contract forbids), so MAJOR, not
CRITICAL — but it also leaves G3's last trigger dead. **Fix:** v1 minimum — the ledger header
states that verdict conflicts are NOT joined and that opposing verdict lines must be read as
DISPUTED by the reader; plus the G6 fixture (same claim reworded under a new ID, opposing PRIMARY
verdicts across the compaction boundary) implemented as a test that asserts join-or-fallback. The
idea must not be recorded as fully implemented while this test is absent.

### [MINOR] NIT findings and RESERVED verdicts are never carried

Never-cut says "every live objection or finding"; no severity floor was adopted. "NIT" is absent
from `ledgerMarkers` and "RESERVED" is absent while "BLOCK" is present. **Fix:** add both markers,
or amend the never-cut claim in FINAL.md to name the floor explicitly — silently having a floor is
the one thing this idea's own data argued against.

### [MINOR] The "byte-for-byte" claim fails for indented lines

`renderLedger` emits `strings.TrimSpace(it.Line)` (`frontier.go:162`), so an indented finding (a
nested list item, a quoted block) loses its leading whitespace. The verbatim test only covers an
unindented line. **Fix:** store and emit the raw line; trim only a copy for the emptiness check;
add an indented-line case to `TestLedgerCarriesVerbatimLinesNotParaphrase`.

### [MINOR] No test covers `gatherReviewContext` — the two MAJORs above escaped through it

`frontier_test.go` tests `frontierContext` and the design side only; the review-side wiring (head
prepend, dedup, fallback interplay) has zero coverage. **Fix:** tests for review round 1
head-exactly-once; round ≥ 3 = head + ledger + previous round; fallback delivers head exactly once;
and the banner-phrase-in-carried-content case.

### [MINOR] `TestFrontierDoesNotTouchTheConsensusDrafterPath` is vacuous

It iterates two strings and asserts nothing (`frontier_test.go:153-158`); it cannot go red, so G4
has no automated enforcement and G7 could never be satisfied for it. **Fix:** a real assertion
(e.g. fail if any file under `internal/app` references `frontierContext`), or delete the test and
say G4 is review-enforced.

### [NIT] Instruction claims a ledger even when the fallback fired

`runner.go:997` tells every round ≥ 2 participant "older rounds appear as a verbatim carry-forward
ledger" — false when the fallback delivered full history (the banner contradicts the instruction),
and vacuous for round 2, which has no older rounds. **Fix:** make the sentence conditional, or
generic: "older rounds appear either in full or as a carry-forward ledger; a banner says which."

### [NIT] IMPLEMENTATION.md's "until its round falls out of scope" is wrong

The window is rounds 1..N−2 starting at r=1 for every N — no round ever falls out of scope, so a
withdrawn objection's line is carried forever (fail-safe direction, but cumulative growth, which is
codex-1's packet-growth caution and my reservation 3's unmeasured write-side cost). Fix the
sentence; consider measuring ledger bytes per round on one real idea.

### [NIT] G7 is incomplete and its harness is unsafe

Two of three reverts were verified; "fallback announces itself" was not, and the killed harness
left the revert applied in the tree. The restore was manual. Re-run the third revert with a harness
that restores on kill (defer + post-run tree check) before signoff; the incident demonstrates the
harness can silently ship a reverted banner.

## Open questions

- Does the owner accept G6 shipping red with the header-notice mitigation, or does the fixture gate
  the idea's completion? (Review can flag; it cannot amend gates.)
- Should the Phase 7 drafter get the frontier treatment eventually (with an authored ledger), or is
  full history permanently correct there given its close-the-loop duty?
