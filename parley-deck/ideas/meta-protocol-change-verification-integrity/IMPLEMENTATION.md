---
idea: meta-protocol-change-verification-integrity
implementer: claude-1
status: ready-for-re-review
cycle: 1
date: 2026-08-04
---

## What was implemented

`FINAL.md` ratified six rules as a new `§15 Verification integrity` plus two text fixes to existing
sections. This is a **protocol-text change only** — no Go code, no CLI behaviour, no new tooling,
per the ratified constraint that every rule be checkable by a participant reading canonical files.

### Files changed

| File | Role |
|---|---|
| `parley-deck/COOPERATION.md` | live project deck |
| `internal/protocol/defaults/COOPERATION.md` | embedded `parley init` bootstrap template |

The two must stay byte-identical outside five allowlisted project-specific zones
(`internal/protocol/drift_test.go`). Every edit was applied to **both**.

### The edits

1. **New `## 15. Verification integrity`** after §14, subsections 15.1-15.7. §14 was the last
   section, so nothing renumbered.
2. **§4.0 invariant list** — the honest qualifier reconciling §4.0 with §11.A.
3. **§6 rule 4** — the scoping sentence.
4. **Phase pointers** under Phases 1, 2, 3 and 6. Specified by `consensus.md:17-18`.

---

## Cycle 1 — fix-up after review round 1

Review round 1 returned **no `CRITICAL` and nine agreed fixes** (`review/consensus.md`, AF-1 to
AF-9). The finding was that cycle 0 was **not a faithful transcription of the ratified text**: five
semantic changes, two wrong cross-references, one dropped binding sentence and a catalogue of
undisclosed micro-deviations, against exactly one declared deviation.

That was the fourth time in this idea that the same agent moved ratified text silently. hermes-1
called it a pattern rather than an accident; the implementer accepts that.

### The method changed, because the method was the defect

Cycle 0 was hand-transcribed. **Cycle 1 does not transcribe.** §15's rule bodies are now
**extracted mechanically from the ratified blockquotes in `consensus.md`** and assembled with
headings — there is no retyping step for the implementer to corrupt. The extraction and the
verbatim comparison are both scripted, and the comparison is the check reviewers should re-run.

### Fixes applied

| AF | Fix |
|---|---|
| AF-1 | All five semantic deltas removed by regeneration: the *"and requests a verdict"* narrowing, the stable-identifier obligation, the reinstated *"exactly one provenance tag"*, the *"or steps"* widening, and the *"and the artifact"* strengthening. Verified absent by grep in both copies |
| AF-2 | §15.3 now reads *"the existing user-escalation path"* — no section number, as ratified |
| AF-3 | The *"How this rule is actually enforced"* paragraph is **removed** from both copies. Follow-ups 8 and 9 stay open. The implementer argued to keep it and was overruled by all three reviewers |
| AF-4 | One boundary sentence added to §15.1, **verbatim as codex-1 specified and hermes-1 adopted** — see below |
| AF-5 | `— see §15.1` removed from the §4.0 qualifier; the qualifier itself stays |
| AF-6 | *"Binds on every track."* is present (it is inside the regenerated ratified §15.5 body) |
| AF-7 | Phase 3 pointer widened to name §15.5's and §15.6's duties; one line added to Phase 3's gate deferring to the close-conditions **already binding** under §15.3 and §15.6, with an explicit statement that it adds no new condition — per codex-1's constraint that the fix may only point, never create |
| AF-8 | Every micro-deviation is gone by regeneration, not by hand-patching |
| AF-9 | Corrected below |

### AF-4 — the one deliberate addition

Added to §15.1, exactly as specified by codex-1 and adopted by hermes-1 and kimi-1:

> An assignment of `CONFIRMED`, `WRONG`, or `UNVERIFIED`, or equivalent language that classifies a
> claim as true, false, or not established, is a verification verdict; raw source text or command
> output reported without a truth-status classification is evidence, not a verdict.

It adds no mandatory syntax, no identifier obligation, no artifact locator, and no rule that all
factual reporting is a verdict — the four things the Phase 7 signoffs forbade it from adding.

### Every remaining non-ratified line in §15, disclosed

The target agreed in Phase 7 was: **verbatim transcription of the ratified text plus exactly the
AF-4 sentence.** Here is everything in §15 that is not a ratified rule body, so that no reviewer
has to discover it:

| Location | Text | Status |
|---|---|---|
| §15 intro | Three-sentence section framing: what the section governs, its ratifying idea, and that it composes with §4.0 and the Phase 6 no-suppression rule | **Not ratified rule text** — section framing. Every claim in it restates `consensus.md`; a section needs an opening |
| §15.1 first paragraph | *"A factual assertion enters the verification regime only when…"* | **Ratified** — it is `consensus.md:22-24`, the prose lead-in rather than the blockquote |
| §15.1 second paragraph | the AF-4 sentence | **Ratified in Phase 7** as the sole new normative substance |
| §15.7 | the per-track binding table | **Ratified** — `consensus.md`'s `## Per-track binding`, re-rendered with `yes`/`no` cells instead of tick marks |
| §§15.2-15.6 | — | **Nothing.** Byte-identical to the ratified blockquotes |

**One thing was dropped rather than adapted.** `consensus.md` closes §15.3 with *"This composes
with the ratified P6: review-phase findings still close only by reviewer withdrawal, review
consensus, or a quoted operator ruling."* `P6` is a label from a different idea and is not defined
anywhere in `COOPERATION.md`, so shipping it verbatim would create a dangling reference and
adapting it would be a §15 wording delta — which codex-1's signoff expressly forbade. The sentence
is ratified commentary, not rule text, and its content is already in the §15 intro. It is dropped
and disclosed here rather than paraphrased silently.

### AF-9 — corrections to this file

1. **Wrong authority (cycle 0).** This file said *"`FINAL.md` specified pointers"*. `FINAL.md`
   contains no phase-pointer instruction; `consensus.md:17-18` does. Corrected above.
2. **Environment-dependent measurement stated as absolute (cycle 0, raised by the implementer).**
   Cycle 0 recorded *"25 packages ok, 0 failures"* without qualification. Under §15.2 that claim
   must be scoped to the environment that produced it — see the verification table below.

## Verification

| # | Check | Result | Kind |
|---|---|---|---|
| 1 | `go build ./...` | OK | survival |
| 2 | `go test -count=1 ./...` **in the implementer's environment** | rc=0, 25 packages ok, including `internal/runner` | survival, **environment-scoped** |
| 3 | `go test ./internal/protocol/...` — `TestEmbeddedDefaultMatchesLiveDeck` | ok | fix-proving |
| 4 | `diff` of the two copies | 13 differing lines, all inside the five allowlisted zones | fix-proving |
| 5 | **Verbatim check:** each of §§15.1-15.6 compared against the ratified blockquote extracted from `consensus.md` | 15.2-15.6 **byte-identical, 0 extra lines**; 15.1 has 6 extra lines = the ratified prose lead (3) + the AF-4 sentence (3) | fix-proving, **this cycle's central check** |
| 6 | grep for each of the five AF-1 deltas in both copies | 0 occurrences each | fix-proving |
| 7 | grep for `§8 user-escalation` and `How this rule is actually enforced` | 0 each | fix-proving |
| 8 | `TestDefaultCooperationForInit` — no roster/vendor names in bootstrap output | passes | survival |

**On check 2 — the environment scoping matters here.** In codex-1's sandbox
`internal/runner.TestDurableKillEndToEndRealProcess` fails with *"no recorded boot id; not killed"*.
codex-1 established that it **fails identically at the parent commit** and correctly declined to
attribute it to this change, filing no finding. That is the §15.2 discipline applied by a reviewer
to its own measurement, and it is the model the cycle-0 claim failed to follow. The test is
sandbox-sensitive; neither result generalises to the other's environment, and the honest statement
is that this change touches no Go code at all.

**Checks 5 and 6 are the ones this cycle exists for.** Both would have failed at cycle 0 — that is
what makes them fix-proving rather than decorative. Check 5 is scripted and reviewers should re-run
it rather than read the section by eye; reading by eye is what produced the defect.

## Not in scope

`FINAL.md`'s ten follow-ups, including 8 and 9, which AF-3 kept out of this change.

## Note for review round 2

Both codex-1 and hermes-1 required review round 2 over the complete fix-up diff and both protocol
copies, with a fresh source comparison — explicitly refusing a re-run of the cycle-0 mechanical
checks alone. hermes-1's reason: *"the implementer has now been caught moving ratified text four
times in this idea … and a pattern is not checked by a re-run."*

kimi-1's Phase 7 acceptance carried one condition: that the drafter correct deferred-follow-up
item 2 in `review/consensus.md`, which misstated both the scope and the direction of kimi-1's open
question. **That correction is applied**, and it is recorded there as a fifth instance of the same
defect — one that only surfaced because kimi-1 re-read its own filing.
