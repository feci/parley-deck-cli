---
agent: antigravity-1
idea: meta-protocol-change-devx-speed
review-round: 1
date: 2026-07-03
reviewed-commit: a224621
---

## Summary
The implementation succeeds in delivering the functional core of the `meta-protocol-change-devx-speed` idea—the three-track system, the objective classifier, and the Quickstart guide. The invariants (§1 non-solo, LE-1 refutation-default) are explicitly protected. However, the implementation significantly under-delivers on the "progressive disclosure" and "core compression" goals of `FINAL.md`. While the deferral of the physical appendix move is a defensible risk-mitigation strategy, the failure to consolidate LE jargon and the introduction of cross-reference mismatches (§4.0 vs Phase 0.0) create avoidable friction. In `REFUTATION-DEFAULT` mode, the change is accepted as safe but rejected as "noisy" and incomplete.

## Refutation attempts
- **Classifier Escape (Risky change to Fast/Standard):** Attempted to find a loophole where a security-relevant or irreversible change could be routed to `fast`. Result: **FAIL**. The `deliberation` triggers are "Force if ANY" and are checked first. Even a single-line change to `auth` or `secrets` is caught by the deliberation-forcing row.
- **Classifier Ambiguity (Overlapping triggers):** Checked for cases where a change could match both `fast` and `deliberation` without clear priority. Result: **FAIL**. The row-order priority ("deliberation triggers are checked first") and the "ALL" requirement for `fast` ensure unambiguous routing.
- **Invariant Weakening (Speed vs Safety):** Checked if the `fast` track's skipping of Phase 2 (cross-review) or §9.0 (ping) removes the non-solo or refutation-default requirements. Result: **FAIL**. The "Invariants on every track" section (line 72 in `COOPERATION.md`) explicitly forbids dropping these. Phase 1 (independent analysis) and Phase 6 (reviewers) still exist, ensuring non-solo execution and refutation discipline.
- **Byte-Identity Check:** Ran `go test ./internal/protocol/...` to ensure the two protocol copies are in sync. Result: **SUCCESS (Green)**.
- **Broken References:** Searched for existing cross-references (e.g., to §11, §12) to see if the additive change broke them. Result: **SUCCESS**. Since the implementer deferred physical reordering, existing links remain valid.

## Findings

### [MAJOR] Cross-reference mismatch: `§4.0` vs `Phase 0.0`
The Quickstart guide and the `00-prompt.md` template both refer to the track selection section as `§4.0`. However, the actual header in `COOPERATION.md` is `### Phase 0.0 — Track selection (conditional rigor)`. A user following the Quickstart instructions ("Read §4.0") will find no section with that label, creating avoidable friction in the new "5-minute start" flow.
**Fix:** Rename the header to `### 4.0 Track selection (conditional rigor)` to match the references.

### [MAJOR] Failure to consolidate LE jargon
`FINAL.md` §4 explicitly required: "Consolidate LE-1…LE-11 inline jargon into one 'Loop-engineering rules' block in plain English, cross-referenced from the phases that use them." The implementation instead kept the jargon scattered throughout the document and only added a brief note in the invariants section. This leaves the "core" cluttered with internal implementation IDs (LE-3, LE-4, etc.) that the DevX speed change specifically aimed to abstract away.
**Fix:** Create a "Loop-engineering rules" appendix or consolidation block and move the detailed LE-N rule text there, replacing inline occurrences with plain English.

### [MINOR] Confusion between Phase 0.0 and Phase 0
The protocol now has `Phase 0.0` (Track Selection) followed by `Phase 0` (Kickoff). In common parlance, "Phase 0" is used to refer to the very beginning. Having a "Phase 0.0" before "Phase 0" is counter-intuitive and violates the sequential expectation of phase numbering.
**Fix:** If renumbering all phases is too risky, rename `Phase 0.0` to `4.0 Track Selection` to separate it from the sequential phase labels.

### [MINOR] "Core" length exceeds target (Missing compression)
AC 4 called for a core ≤ ~200 lines. The current core (§0-§8) remains ~660 lines because §1, §5, §6, and §7 were not compressed as planned in `FINAL.md`. While the Quickstart helps navigate this, the underlying complexity of the core text remains high, partially defeating the DevX goal.
**Fix:** Perform the compression of §1, §5, §6, and §7 as specified in `FINAL.md`, moving non-essential detail to appendices.

### [NIT] Table formatting in `COOPERATION.md`
The per-track behavior table uses inconsistent spacing in some cells compared to existing tables (e.g., `skipped` vs `skipped `).
**Fix:** Normalize whitespace in the markdown table.

## Open questions
- The `IMPLEMENTATION.md` records the physical appendix move as deferred to `protocol-restructure-appendices`. Can we at least commit to a "Core" that is physically at the front, even if we don't renumber (e.g., §15 follows §8, but the Quickstart calls it an appendix)?
