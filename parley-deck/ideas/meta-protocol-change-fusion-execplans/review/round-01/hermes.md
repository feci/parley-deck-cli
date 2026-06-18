---
agent: hermes
idea: meta-protocol-change-fusion-execplans
review-round: 1
date: 2026-06-19
reviewed-commit: 8477ea6
---

## Summary

The implementation faithfully realizes FINAL.md Edits 1–7 as additive protocol text across both `COOPERATION.md` copies, with the shared zones byte-identical (drift guard green) and the embedded default still genericized. No reject-list concept (confidence gate, judge authority, autonomy-across-gates, new retro mining code) was imported; FINAL.md immutability and the `CRITICAL/MAJOR/MINOR/NIT` severity vocabulary are explicitly preserved. One MINOR fidelity gap (a spec parenthetical dropped from the Decision Log) and two NITs; otherwise clean.

## Findings

[MINOR] Decision Log omits Edit 2's deviations cross-reference.
Edit 2 specifies the `## Decision Log` head as "decisions made *after* `FINAL.md`: `Decision / Rationale / Date·Author` (deviations still also recorded under the existing 'Deviations' head)." The implemented description reads only "(Decisions made *after* FINAL.md — Decision / Rationale / Date·Author.)" and drops the parenthetical. The retained `## Deviations from FINAL.md` head ("Any divergence, with rationale. 'None' is a valid answer.") still sits directly above, so the two heads coexist with overlapping subject matter — a deviation *is* a post-FINAL decision, so a drafter could be unsure whether to file it under Deviations, Decision Log, or both. The spec's parenthetical existed precisely to close that seam. Impact is low (the heads are distinguishable on close reading) but it is a realized spec instruction that went missing. Fix: append the parenthetical, e.g. "(Decisions made *after* FINAL.md — Decision / Rationale / Date·Author. Deviations still go under `## Deviations from FINAL.md` above.)", byte-identically in both copies.

[NIT] Rigor-trigger vocabulary is slightly inconsistent, inherited from FINAL.md.
The guiding principle (FINAL.md preamble) lists four REQUIRED triggers — complex / `auto_implement` / driver-managed / pipeline — while the Phase 4 self-containedness sentence and the Phase 5 Progress/closing prose list three (omitting "pipeline"), and the Phase 4 idempotence clause lists a different three (`auto_implement` / action / pipeline). This is **not an implementer defect**: Edit 1's own sentence names only three triggers and the implementer reproduced it verbatim, while Edit 4's sentence names "pipeline" and the implementer included it there. The tension lives in FINAL.md itself (preamble vs. edit-level sentences). Flagged only so consensus can decide whether to harmonize "pipeline" into the general trigger for defensiveness — in practice pipeline ideas are driver-managed (§12.5/§12.8), so the gap is mostly theoretical, but a pure-deliberation pipeline that is neither complex nor `auto_implement` could slip through the three-trigger phrasing.

[NIT] Phase 4 prose paragraph is dense.
The new Phase 4 paragraph carries Edits 1, 3, and 4 in a single block (static/frozen nature, self-containedness requirement, observable acceptance criteria, idempotence-as-recovery-contract, and the N/A escape hatch). Every clause is accurate and maps cleanly to a spec edit, and the spec itself described these as one combined prose addition, so this is not a correctness issue — only a readability one. Splitting the idempotence clause into its own sentence/paragraph would make the two distinct requirement triggers (general self-containedness vs. idempotence-specific) easier to scan. Optional.

## Verification

- Edits 1–7 all present and faithful:
  - Edit 1 — five design-time heads added to the FINAL.md template (Purpose / user-visible outcome; Context & orientation; Observable acceptance criteria; Idempotence & recovery; Known risks / de-risking) and the self-containedness + N/A prose; "open a v2" immutability rule retained and "static" reaffirmed.
  - Edit 2 — five living heads added to the IMPLEMENTATION.md template (Progress with ISO `(YYYY-MM-DD HH:MMZ)` and `(completed: X; remaining: Y)`; Decision Log; Surprises & Discoveries; Validation evidence; Outcomes & Retrospective) plus the "living companion" prose tying it to §12.
  - Edit 3 — observable-behavior framing in Phase 4 and the Phase 6 sentence ("Where `FINAL.md` states observable acceptance criteria, reviewers should check the implementation against them and may cite a criterion in a finding; this does **not** change the severity vocabulary — it only makes severity assignment less subjective.").
  - Edit 4 — idempotence folded into the Phase-4 head with the driver recovery-contract prose ("required for `auto_implement` / action / pipeline ideas, where the driver treats it as the recovery contract").
  - Edit 5 — Phase 3 `## Comparison & blind spots` (HTML-comment advisory) + advisory-not-a-gate consensus bullet; Phase 7 `## Coverage & blind spots`; "raw round files are never hidden behind the summary" preserved.
  - Edit 6 — §13.2 evidence-corpus confident-error signal with the explicit guard ("**never** a new review severity, a blame label, or a merge gate"); §13 footer amended to credit `meta-protocol-change-fusion-execplans` (2026-06-18).
  - Edit 7 — §3 layout comments updated to "static, self-contained authoritative artifact" and "living execution doc (Progress / Decision Log / Surprises / Outcomes)".
- Byte-identical shared zones: `go test ./internal/protocol/ -run TestEmbeddedDefaultMatchesLiveDeck` PASS. The diff hunks are textually identical across both files (only line offsets differ, due to the allowlisted roster/workspace/header zones). The edits touch only shared (non-allowlisted) zones.
- Embedded default stays genericized: `internal/protocol/defaults/COOPERATION.md` keeps `<workspace-name>`, `<date> — created by parley init`, empty roster/host-handle tables, and no `Protocol synced:` line. The diff does not touch any allowlisted zone.
- Build/test: `go build ./...` succeeds; `go test ./...` fully green. Acceptance criteria 2 and 3 met.
- No reject-list import: no confidence-by-breadth/majority gate, no dedicated judge role with authority, no hiding of raw rounds (explicitly forbidden in the new text), no Fusion panel/recursion-depth/cost/web-search machinery, no deck-into-one-file collapse, no "proceed-without-prompting" autonomy across gates, no anti-list/anti-table maximalism, and no new `parley retro` mining code (the change is text-only).
- Invariants held: FINAL.md immutability reaffirmed; severities unchanged (Phase 6 states the vocabulary is not changed); non-solo/quorum/signoff mechanics untouched.
- Surgical commit: 8477ea6 touches only the two `COOPERATION.md` copies (47 insertions each) plus the idea's own deliberation artifacts; zero `.go` changes — consistent with FINAL.md's "no Go logic required beyond the drift guard".
- Footer convention: only §12 and §13 carry "Changing this section follows §7… ratified/amended by idea…" footers; §3 and §4 never have (confirmed by search — two matches, both in §12/§13). Amending only the §13 footer, and leaving §3/§4 without footers, is consistent with established precedent, so FINAL.md acceptance criterion 5 is satisfied.
- Conditional-rigor "N/A" escape hatch is unambiguous: stated three times (Phase 4 "these added sections may be `N/A`"; Phase 5 Progress "Required for complex / `auto_implement` / driver-managed ideas; 'N/A' for trivial or design-only work"; Phase 5 closing "may be `N/A` for trivial or design-only work") with a consistent trigger and a consistent escape. The idempotence head carries its own tighter trigger (`auto_implement` / action / pipeline), matching Edit 4.
