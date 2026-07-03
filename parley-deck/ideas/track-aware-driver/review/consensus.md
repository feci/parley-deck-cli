---
idea: track-aware-driver
review-cycle: 1
drafted-by: claude-1
date: 2026-07-03
reviewed-commit: ce4ca22
---

## Review consensus — cycle 1 (Phase 7)

Three non-implementer reviewers (codex-1, hermes-1, antigravity-1) ran a refutation review over
three rounds. **All three now ✅ ACCEPT.** Zero Agreed fixes remain → the idea completes.

## Agreed fixes (all applied and verified)
Round-01 (fix-up cycle 1):
- **[MAJOR codex-1 + antigravity-1]** explicit `track: deliberation` bypassed the non-solo
  hard-reject → non-solo floor now enforced for every explicit track (above the switch).
- **[MAJOR codex-1]** classifier treated unknown/negative size as fast → `FilesKnown`/`LOCKnown` +
  non-negative required; `parley classify` rejects negatives (exit 2).
- **[MAJOR codex-1]** `fast` + idea-level `auto_implement` bypassable via `--no-implement` →
  contradiction check reads idea-level `ReadAutoImplement`, not the runtime-masked flag.
- **[MAJOR hermes-1]** `fast` did not force a model-diverse reviewer → hard gate on `fast`.
- **[MINOR hermes-1/antigravity-1]** standard did not cap cross-review at 2 → `CapCrossReviewRounds`.

Round-02 (fix-up cycle 2):
- **[MAJOR codex-1 — was ❌ BLOCK] pre-existing driver-lock TOCTOU** in `acquireLock`
  (`TestAcquireLockIsExclusive` granted two holders under `-count`): empty/unparseable lock
  content was treated as reclaimable-stale. Fixed per codex's counter-proposal — empty/unparseable
  content is now HELD (refuse), never reclaimed. Verified green under `-count=50`; codex confirms
  it could not construct a two-holder interleaving after the fix.
- **[MINOR hermes-1]** no test for the fast model-diversity hard gate → test added.

## Deferred follow-ups (out of scope, documented)
Per-track timeouts (`agents.TimeoutMS` precedence — `track-timeouts`); §9.0 ping-skip for fast;
collapsed consensus/FINAL for fast; per-phase human gates; mid-idea upgrade via diff scan;
round-01 participant truncation for fast; model-diversity-preserving reviewer truncation order;
auto-defaulting the `00-prompt` template to `track: standard`. None is a safety item.

## Dismissed findings
None dismissed; the residual round-01 MINOR/NIT items were open questions or cosmetic and are
recorded as deferred.

## Coverage & blind spots
All three reviewers independently confirmed: the five track findings resolved; backward-compat
(absent/deliberation ≡ today) holds; the classifier is fail-safe (verified by smoke + tests); the
non-solo and contradiction gates fire and halt; refutation stays structural on every track; and
the lock TOCTOU is closed (hermes reasoned through every interleaving + ran the test 50×; codex
could not refute). The `internal/runner/TestDurableKillEndToEndRealProcess` failure is the known
codex-sandbox "no recorded boot id" limitation (green in the implementer's env and in hermes's /
antigravity's runs).

## Signoffs
<!-- Each reviewer's ✅ authored in review/round-03/<agent-id>.md; assembled here. -->

### Signoff: codex-1 — 2026-07-03
Status: ✅ ACCEPT — round-02 BLOCK resolved; `acquireLock` no longer treats empty/unparseable
content as stale; could not construct a two-holder interleaving; fast model-diversity gate tested.

### Signoff: hermes-1 — 2026-07-03
Status: ✅ ACCEPT — two-hunk fix-up traces directly to the two findings; lock fix verified by
interleaving analysis + 50× test run; full suite green; `go vet` clean.

### Signoff: antigravity-1 — 2026-07-03
Status: ✅ ACCEPT — TOCTOU closed; fast model-diversity gate tested and verified.

### Signoff: claude-1 (implementer/facilitator) — 2026-07-03
Status: ✅ ACCEPT — unanimous ✅ ×3 + implementer; zero Agreed fixes remain. Marking complete and
proceeding to release (CLI-only; protocol text unchanged, so no skill/drift release).
