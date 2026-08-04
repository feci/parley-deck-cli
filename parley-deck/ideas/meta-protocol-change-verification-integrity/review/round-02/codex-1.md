---
agent: codex-1
idea: meta-protocol-change-verification-integrity
round: 2
phase: review
date: 2026-08-04
---

## Summary

One `NIT` remains: the shipped §15.2 heading adds the unratified word `Verdict`, and the
`IMPLEMENTATION.md` table that claims to enumerate everything outside the ratified rule bodies
does not disclose that adaptation. I found no `CRITICAL`, `MAJOR`, or `MINOR` issue.

The rule bodies survived a fresh, independent source comparison. In both protocol copies,
§§15.2-15.6 are byte-identical to the de-quoted ratified blockquotes. Section 15.1 has exactly six
additional nonblank lines: the ratified prose lead and the exact AF-4 sentence; two blank Markdown
separators make the physical prefix eight lines, after which the ratified blockquote is
byte-identical. No boundary line was lost or duplicated.

This is **not ready to mark complete and release** while the NIT is unresolved. The fix is small
and non-normative, but Phase 8 has not yet reached a review consensus with zero agreed fixes.

## Scope checked

- Read `review/consensus.md` in full, including AF-1 through AF-9 and the codex-1, hermes-1, and
  kimi-1 signoffs; read cycle 1 of `IMPLEMENTATION.md` in full.
- Read the live `parley-deck/COOPERATION.md` in full. Inspected the embedded
  `internal/protocol/defaults/COOPERATION.md` through a complete diff against the live copy and
  direct extraction of every changed surface: §4.0 qualifier, four phase pointers, Phase 3 gate,
  §6 rule 4, and §§15.1-15.7.
- Read the ratified `consensus.md` §§15.1-15.6 and `## Per-track binding`, plus the matching
  `FINAL.md` rule and text-fix summaries.
- Extracted `bfca39e` with `git archive bfca39e | tar -x -C <temporary-directory>` and reviewed the
  complete cycle-1 protocol and `IMPLEMENTATION.md` diffs against that archive. No repository
  file was used as a temporary comparison output.
- Ran `go test -count=1 ./internal/protocol/...`; it returned
  `ok parley-deck-cli/internal/protocol`.
- I did not rerun the full Go survival suite. No Go source changed in cycle 1, and the review
  requirement was a fresh source comparison rather than another reliance on cycle 0's mechanical
  checks.

## Ownership and verdict discipline

AF-1, AF-3, AF-4, AF-5, and AF-9's wrong-authority item include claims originating in my
round-01 review. I issue no verification verdict on those owned claims. For them I report the
source text, comparison result, or command output without a truth-status classification.

Formal verdicts below are limited to claims owned by the implementer or by other reviewers. Each
`PRIMARY` verdict is based on the located files and the fresh comparisons described here.

## Fresh source comparison

I did not use or invoke the implementer's comparison script. For each ratified subsection, I used
an independent `awk` extraction bounded by its `### 15.N` heading and the next `###` heading. The
source extraction emitted only lines beginning with `>` and removed exactly the blockquote prefix;
the target extraction emitted the shipped subsection body and trimmed only leading and trailing
blank lines. I then used `cmp -s` on the extracted bytes. For §15.1 I separately compared the
target after its eight-line physical prefix.

**Implementer's §§15.2-15.6 body-fidelity claim: `CONFIRMED` — `PRIMARY`.** Relevant output from
the independent comparison, identical for the live and embedded copies:

| Section | Ratified block | Shipped body | Result |
|---|---:|---:|---|
| 15.2 | 1,316 bytes / 18 lines | 1,316 bytes / 18 lines | byte-identical |
| 15.3 | 1,184 bytes / 17 lines | 1,184 bytes / 17 lines | byte-identical |
| 15.4 | 729 bytes / 10 lines | 729 bytes / 10 lines | byte-identical |
| 15.5 | 945 bytes / 12 lines | 945 bytes / 12 lines | byte-identical |
| 15.6 | 1,340 bytes / 20 lines | 1,340 bytes / 20 lines | byte-identical |

For §15.1, the ratified blockquote is 1,389 bytes / 20 lines and the shipped body is 1,891 bytes /
28 lines. Lines 1-3 are the consensus prose lead at `consensus.md:22-24`, word-for-word after
removing the deliberation label `Unanimous.` and reflowing. Lines 4 and 8 are blank separators.
Lines 5-7 are AF-4. Lines 9-28 are byte-identical to the complete ratified blockquote. Thus the
implementation's count is accurate as six additional content lines, and the block boundary is
faithful.

The AF-4 sentence at live `COOPERATION.md:1189-1191` and embedded
`COOPERATION.md:1180-1182` is a physical-line byte match to codex-1's required sentence at
`review/consensus.md:230-232`. It introduces no mandatory syntax, identifier obligation, artifact
locator, or classification of all factual reporting as a verdict.

The §15.7 table is cell-identical to `consensus.md:223-231` after only the disclosed rendering of
`✔` as `yes` and `—` as `no`; this holds in both copies.

## AF-1 through AF-9

Rows described as evidence-only below deliberately do not assign a truth status to claims I own.

| AF | Result of this review |
|---|---|
| AF-1 | Evidence only: searches of both copies returned zero occurrences of all five removed deltas, and the independently extracted §§15.1-15.2 bodies contain the ratified wording. |
| AF-2 | **`CONFIRMED` — `PRIMARY`.** Both §15.3 bodies say `the existing user-escalation path`; neither contains `§8 user-escalation`. |
| AF-3 | Evidence only: `How this rule is actually enforced` occurs zero times in both copies. `FINAL.md` follow-ups 8 and 9 remain follow-ups. |
| AF-4 | Evidence only: the sentence is byte-exact in both copies, with none of the four forbidden additions. |
| AF-5 | Evidence only: the §4.0 qualifier in both copies ends at `per-agent isolated staging)`; `— see §15.1` occurs zero times. |
| AF-6 | **`CONFIRMED` — `PRIMARY`.** `Binds on every track.` is present in both §15.5 bodies and in the ratified source. |
| AF-7 | The Phase 3 pointer names §15.5's drafter-position section and §15.6's close conditions. The gate line says it is subject to conditions already binding under §15.3 and §15.6, and expressly says it adds no new condition. My ruling is that it only points: it neither creates a third condition nor changes either referenced condition. |
| AF-8 | **`CONFIRMED` — `PRIMARY`.** The byte comparison covers every §15 micro-deviation. Separate searches found none of the removed elaborations, and both §6 copies use the ratified `§6 rule 4 applies to scoping:` lead. |
| AF-9 | Evidence only: `IMPLEMENTATION.md:37` attributes the phase pointers to `consensus.md:17-18`, and check 2 at `IMPLEMENTATION.md:110` scopes its result to the implementer's environment. |

The complete live-versus-embedded diff contains only the five allowlisted project-specific zones:
the header fields and two §2 table bodies. Therefore every AF protocol edit above is present in
both copies, not only in the live deck.

## Ruling on the disclosed §15.3 deletion

Dropping the sentence is the right call. It is outside the ratified §15.3 blockquote: the quoted
protocol body ends with `No new file.` at `consensus.md:109`; the `P6` sentence is commentary at
`:111-112`. Copying every surrounding consensus comment would also import deliberation rationale,
while adapting this one sentence would create the additional §15 wording delta my Phase 7 signoff
forbade.

No binding close rule is lost. The actual rule remains at live `COOPERATION.md:525-528` and
embedded `COOPERATION.md:516-519`: a disputed review finding closes only by reviewer withdrawal,
review consensus, or an explicit operator ruling. The §15 intro points to that Phase 6 rule and
states that §15 does not override it. More precisely than `IMPLEMENTATION.md` phrases it, the intro
contains the non-override pointer, while Phase 6 itself contains the three close routes.

Shipping the sentence verbatim would add another idea-local `P6` reference that has no definition
in `COOPERATION.md`. This choice does not cure the already-ratified `P6` wording inside §15.4;
changing that blockquote would be a separate wording decision, not a reason to duplicate the
§15.3 commentary now.

## Findings

### CRITICAL

None.

### MAJOR

None.

### MINOR

None.

### NIT

#### [NIT] The non-body enumeration omits the adapted §15.2 heading

**Verdict on `IMPLEMENTATION.md:78-87` claiming to list everything in §15 outside the ratified
rule bodies: `WRONG` — `PRIMARY`.** The ratified heading is `### 15.2 — Provenance (CRITICAL-2)`
at `consensus.md:50`, and `FINAL.md:31` also names the rule `Provenance`. Both shipped copies instead
say `### 15.2 Verdict provenance` at live `COOPERATION.md:1214` and embedded
`COOPERATION.md:1205`. The extra word is outside the mechanically extracted body and is absent from
the disclosure table.

This does not change the rule's meaning, so the severity is `NIT`. It does violate the agreed
verbatim-fidelity target and makes the supposedly exhaustive enumeration incomplete. Restore both
headings to `### 15.2 Provenance`, and update the implementation disclosure to account for the
structural heading adaptations (removal of deliberation-only dashes/finding labels and numbering
of `Per-track binding` as §15.7), or narrow its exhaustive claim explicitly to non-structural body
prose.

## Release decision

**No — not yet ready to mark complete and release.** Resolve or explicitly dismiss the one NIT in
review consensus. I found no other obstacle in the complete fix-up diff.
