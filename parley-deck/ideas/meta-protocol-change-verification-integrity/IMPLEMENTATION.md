---
idea: meta-protocol-change-verification-integrity
implementer: claude-1
status: ready-for-review
cycle: 0
date: 2026-08-04
---

## What was implemented

`FINAL.md` ratified six rules as a new `§15 Verification integrity` plus two text fixes to
existing sections. This is a **protocol-text change only** — no Go code, no CLI behaviour, no new
tooling. That is deliberate: the ratified constraint was that every rule must be checkable by a
participant reading canonical files, without new tooling.

### Files changed

Both copies of the protocol, per the embedded-default drift guard:

| File | Change |
|---|---|
| `parley-deck/COOPERATION.md` | live project deck |
| `internal/protocol/defaults/COOPERATION.md` | embedded `parley init` bootstrap template |

The two files must stay byte-identical outside five allowlisted project-specific zones
(`internal/protocol/drift_test.go`). Every edit below was applied to **both**.

### The edits

1. **New `## 15. Verification integrity`**, appended after §14, with subsections 15.1-15.7:
   scope/ownership/location, provenance, conflicting verdicts, exemption-claim admissibility, role
   concentration, correlated agreement, and a per-track binding table. §14 was the last section, so
   nothing renumbered.

2. **§4.0 invariant list** — the honest qualifier. The protocol previously listed round-1
   independence among invariants *"never dropped for speed"* while §11.A said *"There is no
   enforcement beyond agent discipline."* Now reads:

   > round-1 independence discipline (Phase 1; an unenforced cooperative convention unless kickoff
   > selects §11.B sub-branches or per-agent isolated staging — see §15.1)

3. **§6 rule 4** — the scoping sentence, requiring material gathered while scoping to be copied
   into `00-prompt.md` before participants are invoked, with the disclosure obligation when it
   cannot be shared.

4. **One-line pointers** added under Phases 1, 2, 3 and 6:
   *"Verification verdicts, their provenance, and verdict conflicts follow §15."*
   `FINAL.md` specified pointers rather than duplicated text, so the normative content lives in one
   place and cannot drift between phase sections.

## Deviations from `FINAL.md`

**One, and it is an addition rather than a departure.** §15.5 carries a paragraph *"How this rule
is actually enforced"* that `FINAL.md` records as **follow-up 8** (codex-1's compliance-model
finding) rather than as ratified rule text. I implemented it now because:

- `FINAL.md` states the finding is *"the most useful output"* of the idea and that the section
  "should say what actually enforces it, so that a future drafter does not read §15.5 as a box it
  can tick alone";
- deferring it would ship the exact misreading the follow-up exists to prevent, in the same commit
  that creates the rule;
- it is descriptive and `SHOULD`-level — it adds one obligation (a signer *should* state the scope
  it read, which is follow-up 9, kimi-1's) and no gate.

**This is the deviation to review.** If reviewers judge it out of scope for Phase 5, it comes out
and both follow-ups stay open. It is flagged here rather than folded in silently, which is the
§15.5 discipline applied to implementation.

## Verification

| # | Check | Result |
|---|---|---|
| 1 | `go build ./...` | OK |
| 2 | `go test ./...` | 25 packages ok, 0 failures |
| 3 | `go test ./internal/protocol/...` — the drift guard `TestEmbeddedDefaultMatchesLiveDeck` | ok |
| 4 | `diff parley-deck/COOPERATION.md internal/protocol/defaults/COOPERATION.md` | differs **only** in the five allowlisted zones: `**Workspace:**`, `**Created:**`, `**Protocol synced:**`, and the two §2 table bodies |
| 5 | `grep -n "^## " parley-deck/COOPERATION.md \| tail` | `## 15. Verification integrity` present; §14 still last before it; no renumbering |
| 6 | `grep -c "follow \*\*§15\*\*"` in each copy | 4 and 4 |
| 7 | `TestDefaultCooperationForInit` `mustNotContain` list (no roster CLI names in the bootstrap output) | passes — §15 text contains no agent-family names |

**Fix-proving vs survival:** checks 1-2 and 7 are survival guards — they passed before this change
and must still pass. Checks 3-6 are the ones this change is about; 4 and 6 would have failed had an
edit landed in only one copy, which is the failure mode the drift guard exists for.

**Not verified, deliberately:** the rules' *effect*. Whether §15 improves verification quality is
not testable by this repository's test suite; the ratifying idea's own four-round record is the
only evidence, and `FINAL.md` states plainly that the untested question — whether compliance
happens when nobody runs the check — is still untested.

## Not in scope

The nine tooling follow-ups in `FINAL.md` (T1-T6 fixes, compliance tooling, and the two
protocol-text follow-ups, one of which is implemented above and one of which — a signer stating
its scope — is included as `SHOULD` in the same paragraph). None is required by this change and
none is started here.
