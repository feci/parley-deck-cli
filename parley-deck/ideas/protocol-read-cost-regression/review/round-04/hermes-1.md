---
agent: hermes-1
idea: protocol-read-cost-regression
round: 4
date: 2026-08-11
reviewed-commit: 41e6cd6 (v1.43.1)
---

verdict: CLEAN

# Review — fix-up cycle 4, round 4 (SHIPPED state)

## Summary

v1.43.1 adopts @codex-1's counter-proposal in full: `frontier.go` and
`frontier_test.go` are deleted, and `runner.go` / `phase58.go` are byte-identical
to their pre-idea form (v1.42.1, commit d4256a2). The dormant machinery that was
Finding B is gone, not merely disabled. My round-3 CLEAN verdict was contingent on
the constant being a sufficient gate; @codex-1's stricter reading won, and the
removal resolves the objection at its root rather than relocating it. I find no
defect in the shipped artifact that warrants anything above NIT.

## 1. Are runner.go and phase58.go genuinely back to pre-idea behaviour?

YES. I diffed both files against d4256a2 (v1.42.1, the commit before this idea):

    git diff d4256a2 HEAD -- internal/runner/runner.go   →  0 lines
    git diff d4256a2 HEAD -- internal/runner/phase58.go  →  0 lines

Zero diff. Not "functionally equivalent" — byte-identical. The two code paths this
idea touched are:

- `buildPromptForRound` (runner.go:909) dispatches to `gatherPriorRounds` for
  round ≥ 2 and `gatherReviewContext` for review phases — the original walkers,
  with no `frontierContext` intermediary.
- `gatherPriorRounds` (runner.go:938) and `gatherReviewContext` (phase58.go:278)
  are the verbatim pre-idea walkers. They exclude `_index.md` only; there is no
  `ledgerFileName` / `_ledger.md` exclusion, because that concept no longer exists
  in the codebase.

I searched `internal/` for `frontier`, `frontierContext`, `compactionEnabled`,
`ledgerFileName`, and `_ledger` — zero matches. The files `frontier.go` and
`frontier_test.go` do not exist on disk. There is no residue: no import, no
constant, no helper, no test, no comment referencing the removed feature.

## 2. Is my round-3 objection resolved, or merely relocated?

RESOLVED. My round-3 CLEAN was the losing position on Finding B. I held that
carrying compiled-but-unreachable machinery behind a `const false` was acceptable
because the gate was a reviewed source change. @codex-1's objection — that a
constant-false branch is executed by no test, so "compiled" was never "verified",
and its guards had to be asserted by source-text matching rather than behaviour —
was the stricter reading and it was correct.

Deletion resolves this at the root. There is no constant to flip, no dormant path
to rot, no source-level test asserting properties of code that never runs. The
risk I dismissed as "acceptable" is not merely gated — it is eliminated. Nothing
carries the risk forward. The enablement gate now lives entirely in
`IMPLEMENTATION.md` as a documented contract for a *future* idea, not as code in
this release.

I am not @codex-1, so I cannot speak to whether Finding B is "closed by deletion"
in codex-1's own terms. From my perspective, the substance of Finding B —
unreachable safety code behind a constant — is gone, and with it the maintenance
hazard. My round-3 concern about re-reviewing dead content checks when the
constant flips is moot: there is no constant and no content checks to re-review.
That concern transfers to whatever future idea re-introduces compaction, which is
correct.

## 3. Does anything shipped in 1.43.0/1.43.1 still change what an agent receives?

NO — not from this idea. The two active input changes I flagged in round-3 are
both gone:

1. **`_ledger.md` exclusion.** The walkers in runner.go:951 and phase58.go:296
   exclude only `_index.md`. There is no `_ledger.md` exclusion because the
   `ledgerFileName` constant no longer exists. An agent receives exactly the
   files it received before this idea. No deck contains a `_ledger.md` file
   today, so this is not a practical difference — but the code change is also
   gone, so even the theoretical input perturbation is removed.

2. **Instruction wording.** The cross-review prompt (runner.go:989) reads:

     "READ every prior-round artifact below and respond to the other participants
      by name: where you agree, where you disagree, what you refine. Converge
      toward consensus."

   This is the original pre-idea text. The forward-looking sentence about "older
   rounds appear either in full or as a carry-forward ledger; a banner above says
   which" — which I flagged as a NIT in round-2 and which @codex-1 raised as
   Finding A in round-3 — is gone. The prompt no longer describes a mechanism that
   does not exist. It tells the agent the truth: every prior round appears below
   in full.

The review prompt (phase58.go:236) and review-consensus prompt (phase58.go:348)
were never touched by this idea and remain unchanged. They say nothing about
ledgers or compaction.

What DOES ship in 1.43.1 that was not in v1.42.1 is the **protocol overlay
machinery** — `internal/protocolcore/overlay.go`, `lock.go`, the `render.go`
changes, and `internal/app/protocol.go` additions (`parley protocol overlay
show|validate`). This is from the separate `protocol-overlay-local-extension`
idea, not from this idea. It does not change what an agent receives in a prompt;
it changes protocol rendering and validation plumbing. It is out of scope for this
review except to note: the protocol text in all three COOPERATION.md copies now
describes the overlay as "partially implemented," which is accurate (see §4
below).

## 4. Does the protocol text (all three copies) now describe the overlay accurately?

The three copies of COOPERATION.md are:

1. `parley-deck/COOPERATION.md` (deck view)
2. `internal/protocol/defaults/COOPERATION.md` (bundled default)
3. `~/.hermes/skills/parley-deck/references/COOPERATION.md` (skill fallback)

Copies 1 and 2 were updated in this release cycle (the 7-line diff against
d4256a2). They now read:

  "A deck change — the deck's own overlay (`parley-deck/protocol-overlay.md`) —
  is a smaller act and goes through a normal idea in that deck. The overlay is
  partially implemented: the file grammar, the `parley.protocol-lock/v2` lock and
  composition at the terminal boundary exist and are extend-only; the
  roster-annotation identity slot and the removal of prose-matched zone addressing
  do NOT. Do not rely on the parts that are absent."

This is accurate. The overlay parser, lock, and render-composition code exist in
`internal/protocolcore/`. The roster-annotation identity slot and prose-matched
zone addressing removal are not in the shipped code. The text says "partially
implemented" and names exactly what is and is not present. It does not describe an
intended future as present fact.

Copy 3 (the skill fallback) still reads "the deck's own overlay, once that ships"
— the pre-overlay wording. This is the bundled fallback in the Hermes skill
directory, not in this repo. It is stale relative to copies 1 and 2 but that is a
skill-portability issue, not a defect in this release's shipped code or in this
idea's artifacts. The skill's own drift-check section says to prefer the live
project file and warns the fallback may be stale. I note it for completeness; it
is not a finding against v1.43.1.

Neither the word "overlay" nor "carry-forward" appears in the protocol text in the
context of this idea (read-cost-regression). The overlay text is from the separate
`protocol-overlay-local-extension` idea. The protocol text contains no trace of
the frontier/compaction/ledger concepts this idea introduced and then removed.

## 5. Anything that should make the owner yank or supersede 1.43.1?

No. v1.43.1 is a clean removal. The runner code is byte-identical to v1.42.1. The
build is green, vet is green, and `internal/runner` tests pass (7.6s). No path in
the shipped code compacts, excludes, or perturbs agent input. The protocol text
describes the overlay accurately. The CHANGELOG honestly records both the
technical reason (inert code removed) and the process failure (1.43.0 released on
a non-CLEAN verdict).

The one thing worth saying, and it is a NIT not a blocker: v1.43.0 remains
published. The owner ruled it stays, and 1.43.1 supersedes it. That is a
defensible position — yanking a published release has its own cost — but anyone
who pinned to 1.43.0 specifically carries the dormant machinery and the two
prompt perturbations. The CHANGELOG and IMPLEMENTATION.md make this visible
enough that a consumer can upgrade. I would not yank 1.43.0; I would make sure
the release notes point to 1.43.1 as the recommended version. That is already what
the CHANGELOG does.

## Findings

### [NIT] v1.43.0 remains published with the prompt perturbations

1.43.0 changed what agents received (the `_ledger.md` exclusion and the
forward-looking instruction wording) while delivering no speedup. 1.43.1 reverts
both. Anyone on 1.43.0 should move to 1.43.1. The CHANGELOG says this clearly.
No action required beyond what the owner already ruled.

## What ships

The measured diagnosis, the located quadratic code paths, the signed carry-forward
ledger contract, and the enablement gate — all in this idea's artifacts. The code
ships nothing new from this idea: no compaction, no frontier, no ledger exclusion,
no instruction change. The runner is byte-identical to v1.42.1. That is the
correct outcome for an optimization idea whose mechanism could not be shown safe:
the diagnosis survives, the machinery does not.
