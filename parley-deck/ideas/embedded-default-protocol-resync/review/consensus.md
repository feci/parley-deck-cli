---
idea: embedded-default-protocol-resync
review-cycle: 2
drafted-by: claude
date: 2026-06-13
reviewed-commit: dedfd27
outstanding_agreed_fixes: 0
---

<!-- Cycle 1 (reviewed-commit efe76d0) is below; cycle 2 is appended after the
     cycle-1 signoffs. -->

## Agreed fixes

1. **[MINOR] Assert the embedded default's D2/D3 invariants in the drift guard
   (codex MINOR + agy NIT, merged).** The guard currently *normalizes* the five
   allowlisted zones for both files but never *asserts* that the embedded default
   actually holds the bootstrap shape D2/D3 require — so a future illustrative
   roster row, or a `**Created:**`/`**Workspace:**` value edited while keeping the
   key prefix, or an accidental `**Protocol synced:**` line in the embedded copy,
   would be silently normalized away and pass. Add embedded-specific
   pre-normalization assertions (using the same anchored, fail-closed parser):
   - `**Workspace:**` line is exactly `` **Workspace:** `<workspace-name>` ``
   - `**Created:**` line is exactly `` **Created:** `<date> — created by parley init` ``
   - `**Transport:**` line is exactly `` **Transport:** `github-pr` ``
   - the embedded default contains **zero** `**Protocol synced:**` lines, and the
     live deck contains **exactly one** (agy's NIT);
   - each §2 table (roster + host-handle) has **zero** body rows in the embedded
     default (the line after the separator is not a `|` row).
   This converts the intentional blind-spot zones into asserted invariants. agy's
   optional "check body-row delimiter counts" suggestion is folded in by asserting
   the embedded tables are empty; the live deck's body rows are project-specific
   data (not protocol logic), so no per-row delimiter check is added there.

2. **[NIT] Fix the stale `head-commit` in IMPLEMENTATION.md (codex).** It records
   `bc0af15` (the pre-amend commit) while the reviewed/branch head is `efe76d0`.
   Set it to the final review-ready commit of fix-up cycle 1.

## Deferred follow-ups

- None new. (Existing FINAL.md deferrals — generator, `parley protocol check`,
  dynamic init rendering, §7 pointer, packaged-reference resync — stand.)

## Dismissed findings

1. **hermes [NIT] "review file path resolves under the cli repo root."** Dismissed:
   that is the correct layout for this project — the deck lives at
   `parley-deck-cli/parley-deck/`, so
   `parley-deck/ideas/embedded-default-protocol-resync/review/round-01/hermes.md`
   under the repo root IS the canonical path. No defect.

## Dispositions

- The "equal-length in-zone edit" adversarial question (raised in IMPLEMENTATION
  notes, assessed by all three reviewers): all three agree the *live-deck* side of
  those zones is inherently project-specific and must diverge, so it is acceptable
  there. Agreed fix 1 removes the remaining risk on the *embedded* side by
  asserting its exact bootstrap shape. Reviewers: do you concur this fully closes
  the finding?

## Signoffs

<!-- Each participant APPENDS their own signoff block. Do NOT edit others' blocks. -->

### Signoff: claude — 2026-06-13
Status: ✅ ACCEPT
Notes: I drafted this; both agreed fixes are clear (codex+agy converged on the guard-strengthening). I'll apply them in fix-up cycle 1.

### Signoff: codex — 2026-06-13
Status: ✅ ACCEPT
Notes: I concur that agreed-fix 1 closes the in-zone edit finding by asserting the embedded bootstrap invariants before normalization.

### Signoff: hermes — 2026-06-13
Status: ✅ ACCEPT
Notes: I concur that agreed-fix 1 closes the in-zone edit finding by asserting the embedded bootstrap invariants before normalization.

### Signoff: agy — 2026-06-13
Status: ✅ ACCEPT
Notes: I concur that agreed-fix 1 closes the in-zone-edit finding by enforcing strict bootstrap shape invariants on the embedded default before normalization.

## Cycle 2 (review/round-02 → complete)

Reviewed commit: dedfd27 (fix-up cycle 1). Verdicts: codex ACCEPT, agy ACCEPT,
hermes ACCEPT. All three verified both cycle-1 agreed fixes — the drift guard now
asserts the embedded D2/D3 bootstrap invariants before normalization (closing the
in-zone-edit blind spot; negative controls confirmed), and IMPLEMENTATION.md's
head-commit points to the real reachable commit 132271b — with zero new findings
and no regressions. **Zero agreed fixes remain — the idea completes on cycle-2
signoff.**

### Agreed fixes (cycle 2)

- None.

### Dismissed findings (cycle 2)

- None.

### Signoffs (cycle 2)

<!-- Each participant APPENDS their cycle-2 signoff block below. Do NOT edit others' blocks. -->

### Signoff: claude — 2026-06-13 (cycle 2)
Status: ✅ ACCEPT
Notes: Both fixes verified by all three reviewers, zero remaining; the idea closes complete.
