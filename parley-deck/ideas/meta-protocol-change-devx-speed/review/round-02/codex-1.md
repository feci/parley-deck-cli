---
agent: codex-1
idea: meta-protocol-change-devx-speed
review-round: 2
date: 2026-07-03
reviewed-commit: dc06103
---

## Summary

Round-01's blocking contradiction is resolved. The new §4.0 text now gives the
per-track table a single authoritative override over the older full-lifecycle defaults in
§4, §5, §9.0, and §11, and it adds the missing Phase-7 review-consensus row.

I accept the narrowed scope for this implementation: protocol text, Quickstart, fail-safe
classifier rules, LE glossary, changelog, and drift-guarded protocol copies land here;
deterministic CLI/driver enforcement and the physical appendix move are carried by named
follow-up ideas.

Verification run:

- `go test ./internal/protocol/...` from `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/parley-deck-cli` -> `ok parley-deck-cli/internal/protocol (cached)`

## Verification of round-01 findings

- **CRITICAL reconcile: RESOLVED.** §4.0 now states that the per-track table is the single
  authoritative per-track gate and explicitly overrides the full-lifecycle defaults in
  §4, §5, §9.0, and §11. The added Phase-7 row also removes the previous ambiguity about
  review-consensus signoff counts. The old phase prose still exists, but it is now
  subordinate for `fast` and `standard`; `deliberation` keeps the full lifecycle.
- **MAJOR CLI/script enforcement: PARTIAL, accepted as deferred.** The CLI classifier,
  init/run templating, timeout seeding, auto-advance, and invariant validation are not
  implemented here. The new `ideas/track-aware-driver/00-prompt.md` stub is concrete and
  covers the exact deferred tooling surface from criteria 1-3. I accept this as a
  Phase-7 scope ratification item rather than a blocker for the protocol-text change.
- **MAJOR classifier fail-closed ordering: RESOLVED.** §4.0 now requires
  deliberation-first evaluation, then fast eligibility, else standard, with doubt and
  boundary cases failing closed to the stricter track.
- **MAJOR appendix/core restructuring: PARTIAL, accepted as deferred.** The physical
  appendix move and <=~200-line core are still not implemented. The Quickstart and
  core-vs-reference map deliver the functional reading-guide improvement, and
  `ideas/protocol-restructure-appendices/00-prompt.md` covers the exact physical move,
  renumber, and cross-reference audit. I accept this deferral for this cycle.
- **MAJOR changelog/protocol metadata: RESOLVED for this review scope.** The
  `meta/protocol-changelog.md` entry is present and names the idea, summary, invariants,
  and both deferred follow-ups. `meta/version.json` remains on the source-deck metadata
  values, but §9.0 treats source freshness as advisory and the implementation records the
  runtime/package refresh as a release-step concern rather than a review blocker here.
- **MINOR heading mismatch: RESOLVED.** The heading is now `### 4.0 — Track selection
  (conditional rigor)`, matching the references.
- **MINOR Quickstart fast wording: RESOLVED.** The Quickstart now describes `fast` as
  round-1 plus collapsed `FINAL.md` signoff plus one refutation-default reviewer and up
  to one fix-up cycle, rather than "one review, then done."

## Scope decision (ACCEPT/BLOCK)

**ACCEPT.** The narrowed scope is acceptable because the protocol is usable and safe now,
while the two deferred surfaces are explicitly captured in concrete deliberation-track
follow-up ideas.

## New findings

No blocking new findings.

Non-blocking note: both follow-up stubs are `status: proposed` but their body text says the
follow-up was "ratified" by this idea. I read that as contingent on this Phase-7 review
consensus accepting the scope decision; the review consensus should make that authority
explicit.

## Signoff (Status: ✅ ACCEPT)

I accept the fix-up cycle and the proposed scope narrowing.
