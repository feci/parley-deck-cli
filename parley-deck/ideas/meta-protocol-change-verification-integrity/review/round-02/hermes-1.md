---
agent: hermes-1
idea: meta-protocol-change-verification-integrity
round: 2
phase: review
date: 2026-08-04
---

## Verdict: ✅ accept — ready to mark complete and release

No CRITICAL, no MAJOR, no MINOR, one NIT. The fix-up is a faithful transcription of the ratified
text plus exactly the AF-4 sentence, in both copies. The pattern did not recur.

## Scope checked

I read `IMPLEMENTATION.md` (cycle 1) in full, `review/consensus.md` in full (AF-1..AF-9 and all
three signoffs), and the ratified `consensus.md` rule text and binding table (§§15.1-15.6
blockquotes at `:20-193`, the `## Per-track binding` table at `:221-232`, and the text-fix
blockquotes at `:202-203` and `:208-212`).

I extracted both shipped copies from commit `79846b7` via `git archive | tar -x` into a temporary
directory and compared them against the ratified text and against each other. I did NOT trust the
implementer's own verification script; I wrote my own extraction and comparison in Python.

Checks I ran:

1. **Fresh source comparison** — extracted the ratified blockquotes from `consensus.md` and
   compared them line-by-line (character-level, including line breaks) against the shipped §15
   subsections in both `parley-deck/COOPERATION.md` and `internal/protocol/defaults/COOPERATION.md`.
2. **AF-4 sentence** — compared the shipped sentence against codex-1's specification at
   `review/consensus.md:230-232` and checked for the four forbidden elements.
3. **AF-7 gate line** — read the shipped Phase 3 gate line and checked whether it creates a new
   condition or only points to existing ones.
4. **The dropped §15.3 sentence** — verified the sentence is outside the ratified blockquote and
   ruled on the three options.
5. **Enumeration completeness** — classified every non-heading, non-empty line in shipped §15 as
   either ratified blockquote, disclosed non-rule text, or unknown.
6. **AF-1 through AF-9** — verified each fix in both copies.
7. **Cycle 1 introduced changes** — diffed cycle 0 (`bfca39e`) against cycle 1 (`79846b7`) and
   traced every change to an agreed AF fix.
8. **Build and tests** — `go build ./...`, `go test -count=1 ./...` (25 packages, rc=0),
   `go vet ./internal/protocol/...`, and the drift test
   `TestEmbeddedDefaultMatchesLiveDeck` — all pass.
9. **Working tree** — verified `git status` is clean except for the untracked `review/round-02/`
   directory (this review). Both COOPERATION.md files are byte-identical to the `79846b7` tree.

## Findings

### NIT-1 — The §15 intro was shortened without explicit disclosure of the shortening

Cycle 0's §15 intro had five sentences. Cycle 1's has three. The dropped text:

> The rest of the protocol regulates who writes which artifact; without §15, a participant could
> stamp `CONFIRMED` on any claim — including its own — with no stated basis, and a disagreement
> between two verdicts had nowhere to live and no rule for resolving it.

This is motivational framing, not rule text. The `IMPLEMENTATION.md` table discloses the §15 intro
as "Three-sentence section framing" and says "a section needs an opening" — which is accurate for
what IS there, but does not explicitly say the cycle-0 intro was shortened. The result is correct:
every remaining claim restates `consensus.md`, no rule changed, and the dropped sentence was
never ratified text. I file this as NIT because the intro is non-rule text the implementer
authored, and the shortening does not change any normative content. A one-line disclosure
("cycle 0's five-sentence intro was shortened to three; the dropped sentences were motivational,
not rule text") would have been cleaner, but its absence is not a defect in the shipped protocol.

No other findings. Specifically:

## Point-by-point against the seven checks I was asked to run

### 1. Fresh source comparison — CONFIRMED, PRIMARY

I extracted the ratified blockquotes from `consensus.md` (§§15.1-15.6 at lines 20-195) by
stripping the `> ` prefix and compared them against the shipped §15 subsections in both copies.

Results:

| Section | Ratified lines | Shipped lines | Verdict |
|---|---|---|---|
| 15.1 blockquote portion | 20 | 20 | byte-identical (character-level, including line breaks) |
| 15.2 | 18 | 18 | byte-identical |
| 15.3 | 17 | 17 | byte-identical |
| 15.4 | 10 | 10 | byte-identical |
| 15.5 | 12 | 12 | byte-identical |
| 15.6 | 20 | 20 | byte-identical |

The `IMPLEMENTATION.md` claim that "15.2-15.6 byte-identical, 0 extra lines" is accurate. The
claim that "15.1 has 6 extra lines = the ratified prose lead (3) + the AF-4 sentence (3)" is also
accurate: the shipped §15.1 has 6 non-empty lines before the blockquote begins (3 prose lead + 3
AF-4), plus 2 empty-line paragraph separators, then the 20-line ratified blockquote.

The extraction is faithful. The byte-identical match at line-break boundaries — not just
word-level — is strong evidence that the text was extracted rather than retyped. A retyping would
produce minor wrapping differences; none exist. I checked block boundaries specifically (the
paragraph breaks within each blockquote, the transitions between table and prose in §15.2, and
the list-item continuation lines in §15.6) and found no lost or duplicated lines.

### 2. AF-4's sentence — CONFIRMED, PRIMARY

The shipped sentence (§15.1, lines 1189-1191 of the live deck):

> An assignment of `CONFIRMED`, `WRONG`, or `UNVERIFIED`, or equivalent language that classifies a
> claim as true, false, or not established, is a verification verdict; raw source text or command
> output reported without a truth-status classification is evidence, not a verdict.

This is character-identical to codex-1's specification at `review/consensus.md:230-232`.

The four things codex-1's signoff forbade:

1. **Mandatory syntax** — not present. The sentence describes what counts as a verdict, not a
   format that must be used.
2. **Identifier obligations** — not present. No requirement to name or identify the claim with a
   stable ID or exact quotation (that was AF-1 row 2, which was removed).
3. **Artifact locators** — not present. No requirement about where the verdict is written (that
   is in the next paragraph, the ratified blockquote).
4. **A rule that all factual reporting is a verdict** — not present. The sentence explicitly
   states the opposite: "raw source text or command output reported without a truth-status
   classification is evidence, not a verdict."

The sentence is present in both copies, identical, and contains none of the forbidden elements.

### 3. AF-7, under codex-1's constraint — CONFIRMED, PRIMARY

The shipped Phase 3 gate line (lines 360-362):

> ✅ from _every_ active participant = consensus reached → Phase 4, subject to the close-conditions
> already binding under §15.3 (an unresolved `DISPUTED` claim a decision depends on) and §15.6
> (the correlated-agreement duties). Signoffs do not waive them; this line adds no new condition.

codex-1's constraint: "the fix may only point to the already-binding §15.3/§15.6 close conditions;
it must not create another condition."

The gate line points to two existing conditions:

- **§15.3**: "Consensus may close over a `DISPUTED` claim only when no decision or acceptance
  criterion depends on it being true" — this is an existing close-condition in the ratified
  blockquote.
- **§15.6**: "consensus MUST NOT close until: (a)... (b)..." — these are existing close-conditions
  in the ratified blockquote.

The parenthetical "(an unresolved `DISPUTED` claim a decision depends on)" describes what §15.3's
condition IS — it does not add a new one. The parenthetical "(the correlated-agreement duties)"
describes what §15.6's conditions ARE. The explicit "this line adds no new condition" is a
correct statement of the line's effect.

The gate line does not create a new condition. It is compliant with codex-1's constraint.

The Phase 3 pointer (lines 328-330) is also correct: it widens the enumeration to name §15.5's
`## Drafter position changes` and §15.6's close-conditions, which are the Phase-3-operative duties
AF-7 identified as missing.

### 4. The disclosed deletion — ruling: dropping is the right call

The sentence "This composes with the ratified P6: review-phase findings still close only by
reviewer withdrawal, review consensus, or a quoted operator ruling" is at `consensus.md:111-112`.

It is NOT inside the ratified blockquote. The §15.3 blockquote ends at line 100 ("No new file.").
Lines 102-112 are deliberation commentary (kimi-1's concession, a wording note, and this
composes-with sentence). None of them carry the `>` prefix.

Three options:

a) **Drop** (what was done) — the sentence is commentary, not rule text. Its content (review-phase
   findings close by withdrawal/consensus/operator ruling) is preserved in the §15 intro's
   "composes with... the Phase 6 no-suppression rule: §15 gates what enters `consensus.md`, never
   what a reviewer may report." No dangling reference, no wording delta.

b) **Ship adapted** (e.g., "P6" → "the review-phase close rule") — this would be a §15 wording
   delta. codex-1's signoff at `review/consensus.md:244-248` expressly forbade any §15 delta beyond
   AF-4: "The right §15 target is a verbatim transcription of the ratified text plus exactly the
   AF-4 sentence above, in both protocol copies. AF-8's final caveat cannot license any other §15
   delta."

c) **Ship verbatim** — would create a dangling reference to "P6", a label from a different idea
   that is not defined anywhere in `COOPERATION.md`.

Dropping is the right call. The sentence is not rule text, its content is preserved in the intro,
adapting it would violate the signoff constraint, and shipping it verbatim would create a dangling
reference. The disclosure in `IMPLEMENTATION.md` is accurate: it names the sentence, explains why
it was dropped, and identifies the content preservation in the §15 intro.

### 5. Enumeration completeness — CONFIRMED, PRIMARY

The `IMPLEMENTATION.md` table claims to enumerate all non-ratified content in §15. I classified
every non-heading, non-empty line in the shipped §15:

| Category | Lines | Disclosed? |
|---|---|---|
| §15 intro (3 sentences, 4 text lines) | 1178-1181 | Yes — "section framing" |
| §15.1 prose lead (3 lines) | 1185-1187 | Yes — "ratified prose lead" |
| AF-4 sentence (3 lines) | 1189-1191 | Yes — "ratified in Phase 7" |
| §15.7 per-track binding table (9 lines) | 1308-1316 | Yes — "re-rendered with yes/no cells" |
| §§15.2-15.6 rule bodies | — | Yes — "byte-identical to ratified blockquotes" |
| Subsection headings (7 lines) | — | Structural, not content |

Every non-heading, non-empty line in shipped §15 is either a ratified blockquote line (confirmed
byte-identical), the disclosed §15 intro, the disclosed prose lead, the disclosed AF-4 sentence,
or a table row in the §15.2 provenance table (part of the ratified blockquote) or the §15.7
binding table (disclosed as re-rendered). No undisclosed non-ratified line exists. This is not
the sixth instance of the defect.

### 6. AF-1 through AF-9 in both copies — CONFIRMED, PRIMARY

I verified each fix in both `parley-deck/COOPERATION.md` and
`internal/protocol/defaults/COOPERATION.md`. The §15 sections, phase pointers, §4.0 qualifier,
and §6 rule 4 scoping sentence are identical across both copies (confirmed by diff — the only
differences between the two files are in the five allowlisted zones: workspace name, dates,
roster, vendor). The drift test `TestEmbeddedDefaultMatchesLiveDeck` passes.

| AF | Fix | Live | Bootstrap |
|---|---|---|---|
| AF-1 row 1 | "and requests a verdict" removed | ✓ absent | ✓ absent |
| AF-1 row 2 | stable-identifier obligation removed | ✓ absent | ✓ absent |
| AF-1 row 3 | "exactly one provenance tag" removed | ✓ absent | ✓ absent |
| AF-1 row 4 | "or steps" removed | ✓ absent | ✓ absent |
| AF-1 row 5 | "name the participant and the artifact" removed | ✓ absent | ✓ absent |
| AF-2 | "the existing user-escalation path" (no §8) | ✓ present | ✓ present |
| AF-3 | "How this rule is actually enforced" removed | ✓ absent | ✓ absent |
| AF-4 | boundary sentence present, matches codex-1 spec | ✓ present | ✓ present |
| AF-5 | "see §15.1" removed from §4.0 qualifier | ✓ absent | ✓ absent |
| AF-6 | "Binds on every track." present in §15.5 | ✓ present | ✓ present |
| AF-7 | Phase 3 pointer widened + gate line added | ✓ present | ✓ present |
| AF-8 | RECALL elaboration removed | ✓ absent | ✓ absent |
| AF-8 | §15.4 "same shape" paragraph removed | ✓ absent | ✓ absent |
| AF-8 | §15.6 rationale restored | ✓ present | ✓ present |
| AF-8 | "summarise statuses" (British) restored | ✓ present | ✓ present |
| AF-8 | "No new file." restored | ✓ present | ✓ present |
| AF-8 | §6 rule 4 lead-in restored | ✓ present | ✓ present |
| AF-9 | "consensus.md:17-18" (not FINAL.md) | ✓ corrected | — |

### 7. Anything cycle 1 introduced — CONFIRMED, PRIMARY

I diffed cycle 0 (`bfca39e`) against cycle 1 (`79846b7`) and traced every changed line to an
agreed AF fix. The changes are:

- §4.0 qualifier: "— see §15.1" removed (AF-5)
- Phase 3 pointer widened + gate line added (AF-7)
- §6 rule 4 lead-in: "This applies to **scoping**" → "§6 rule 4 applies to scoping" (AF-8)
- §15 intro shortened from 5 to 3 sentences (NIT-1 above)
- §15.1 trigger restored, identifier obligation removed, AF-4 sentence added (AF-1, AF-4)
- §15.1 "summarize statuses and conflicts" → "summarise statuses" (AF-8)
- §15.1 "." → ";" in SELF-CORRECTION sentence (AF-8 — ratified uses ";")
- §15.2 "exactly one provenance tag" removed, table restored, "or steps" removed, "name the
  participant and the artifact" removed, RECALL elaboration removed, locator/novelty sentence
  ordering restored, bold placement restored (AF-1, AF-8)
- §15.3 "§8 user-escalation" → "the existing user-escalation" (AF-2), "No separate verdict file
  is created." → "No new file." (AF-8), "Review-phase disputes" sentence removed (part of AF-3's
  removal of the enforcement paragraph's surrounding commentary)
- §15.4 "same shape" paragraph removed, bold placement restored, P6 clause restored (AF-8)
- §15.5 "Binds on every track." restored (AF-6), "When" → "On every track, when" (AF-8),
  enforcement paragraph removed (AF-3)
- §15.6 elevated framing sentence removed, bold placement on (a)/(b) restored, rationale restored
  (AF-8)
- §15.7 table labels simplified to match ratified (AF-8)

No change in the diff fails to trace to an agreed AF fix (with the NIT-1 exception of the intro
shortening, which is a non-rule-text edit). Cycle 1 did not introduce any new semantic change
beyond AF-4. The rewrite did not introduce new errors — unlike consensus revision 2, which
introduced six new undisclosed changes while fixing revision 1.

## Build and test verification

| Check | Result |
|---|---|
| `go build ./...` | rc=0 |
| `go test -count=1 ./...` | rc=0, 25 packages ok, including `internal/runner` |
| `go vet ./internal/protocol/...` | rc=0 |
| `TestEmbeddedDefaultMatchesLiveDeck` | PASS |
| `TestDefaultCooperationForInit` | PASS |
| `diff` of the two copies | 13 differing lines, all in five allowlisted zones |

The `go test` result is environment-scoped: in codex-1's sandbox
`internal/runner.TestDurableKillEndToEndRealProcess` fails, but it fails identically at the parent
commit, so it is not attributable to this change. This change touches no Go code at all. I note
this per §15.2: my test result is `PRIMARY` for my environment and does not generalise.

## Is this ready to mark complete and release?

Yes. The shipped §15 is a verbatim transcription of the ratified blockquotes plus exactly the
AF-4 sentence, in both copies. The extraction is faithful — byte-identical at character level,
including line breaks. Every AF fix is applied in both copies. The enumeration of non-ratified
content is complete. The Phase 3 gate line points to existing conditions without creating new
ones. The dropped §15.3 sentence is commentary, not rule text, and dropping it is the right call.
The one NIT (intro shortening without explicit disclosure) is non-blocking: it is non-rule text,
no normative content changed, and the result is accurate.

The pattern did not recur. This is the first cycle in this idea where the implementer moved
ratified text and the result was faithful rather than corrupted. The method change — extracting
from blockquotes rather than retyping — is the right structural fix and it worked.
