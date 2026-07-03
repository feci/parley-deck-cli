---
agent: claude-1
idea: track-aware-driver
round: 1
date: 2026-07-03
---

## Summary

The existing driver already has almost every field §4.0 needs (`CrossReviewRounds`,
`MaxFixupCycles`, the reviewer set, the LE-11 `ReviewerCount < 2` gate). The MVP is therefore
mostly a **derivation + capping** layer, not new machinery: read `track:`, cap the reviewer
slice and set the loop fields per track, and make the one hard invariant (≥1 independent
reviewer, never 0) impossible to violate. The classifier belongs in a **separate, pure,
flag-driven `parley classify` helper** so the driver stays deterministic (it reads the declared
`track:`, it does not guess). Keep `standard`/absent byte-for-byte equal to today's behavior.

## 1. Classifier input model (MVP recommendation)

**Hybrid (a)+(c), with (a) as the deterministic core.**
- **(a) `parley classify` — a pure, flag-driven advisor.** `parley classify --files N --loc N
  [--security] [--irreversible] [--protocol-change] [--auto-implement] [--pipeline]
  [--api-break] [--data-migration]` prints exactly one of `fast|standard|deliberation` (and
  `--json` with the matched trigger). It implements §4.0 verbatim: evaluate deliberation
  triggers first → if any, `deliberation`; else if ALL fast conditions hold (reversible, files
  ≤5, loc ≤300, no security/data) → `fast`; else `standard`. **Fail-safe:** any unknown/omitted
  fast input ⇒ not-fast; any plausible-but-unset deliberation flag the user passes as
  `--maybe-*` ⇒ deliberation. This is fully script-checkable and unit-testable.
- **(c) declare + validate.** The driver reads the declared `track:` (`ReadTrack`). `parley run`
  may call the same classifier over what it *can* observe (e.g. `git diff --numstat` file/LOC
  counts when a code diff exists) and **warn on an obvious under-tier** (declared `fast` but
  files/LOC exceed the fast ceiling), never silently override.
- **Reject (b) as the MVP core:** inferring "reversibility / security surface" from a task
  string is not deterministic; git-diff inference is only available late (post-implementation)
  and only covers size, not risk class. Use it as the *warn* signal in (c), not the source of truth.

Rationale: the driver must be deterministic and fail-closed; a pure declared-track read + a
separate pure classifier is the cleanest split (mirrors "deterministic routing vs model-driven
content" from §4.0).

## 2. Track → Config mapping (exact, with threading points)

Add `ReadTrack(ideaDir) string` to `internal/driver/transport.go` (mirror `ReadStrictGate`;
normalize, default `standard`). Then derive in `internal/app/driver_impl.go` (where the
participant/reviewer set and `driver.Config` are built) and `driver.New`:

| Field / behavior | fast | standard | deliberation | Where |
| --- | --- | --- | --- | --- |
| Reviewer set (Phase 6) | `reviewers[:1]` | `reviewers[:2]` | all (today) | `newDriverImplOps` in `driver_impl.go:45-62` — cap the built `reviewers` slice by track |
| LE-11 auto-complete guard | `< 1` escalates | `< 2` | `< 2` (today) | `impl.go:240` — make the threshold track-aware (fast min = 1) |
| `CrossReviewRounds` | 0 | min(2) | default (today) | `driver.Config` build; `driver.New` defaults |
| `MaxFixupCycles` | 1 | 2 | current default 3 | `driver.Config` build |
| Per-agent timeout | ~5m | ~15m | ~30m | `internal/config/runtime.go` `[defaults.timeouts]` → round/review deadline |
| Collapsed consensus/FINAL | fast only | separate | separate | **no current field** — see below |

**§4.0 behaviors with NO current Config field:**
- **Reviewer *count* cap** is not a Config field — today all non-implementers review. Minimal
  add: cap the `reviewers` slice in `newDriverImplOps` by track (deterministic prefix; keep
  model-diversity ordering so the kept reviewer isn't the implementer's model when avoidable).
- **Collapsed consensus/FINAL for fast** — the consensus→final path is two phases in
  `consensus.go`/`impl.go`. This is the largest piece; **defer to a later slice** (fast still
  works correctly with separate consensus/FINAL, just not collapsed — a speed nicety, not a
  correctness gap).
- **§9.0 ping skip for fast** — lives in `parley preflight`, not the driver; separate slice.

## 3. Invariant enforcement (hard-reject, never silent)

The one non-negotiable: **≥1 independent reviewer on every track** (non-solo). Enforce at two
points, both fail-closed:
- `impl.go:221` already errors on `len(o.reviewers) == 0` ("no non-implementer reviewers") —
  keep it; the track cap must **never produce 0** (cap is `min(trackMax, len)` with a floor of 1
  when any non-implementer exists).
- Add a guard in the Config derivation: if a track would set reviewer max `< 1`, refuse to build
  the driver (`return error`), never proceed. Refutation is already structurally enforced
  (`ValidateReviewArtifact` requires the refutation section) and is track-independent — no track
  may bypass `ValidateReviewArtifact`.

## 4. MVP slicing (by (safety+value) ÷ risk)

- **Slice 1 (MVP, safety-critical):** `ReadTrack` + reviewer-slice cap + track-aware LE-11
  threshold + `CrossReviewRounds`/`MaxFixupCycles` derivation + the ≥1-reviewer floor. Ships the
  real enforcement. Backward-compat: absent/`standard` reproduces today exactly.
- **Slice 2:** `parley classify` pure flag classifier (+ `--json`) with full §4.0 unit tests;
  `parley run`/`init` read `track:` and wire the derivation; warn-on-under-tier.
- **Slice 3:** per-track timeout seeding (`init` → `[defaults.timeouts]`; driver round/review deadline).
- **Slice 4 (largest, deferrable):** collapsed consensus/FINAL for fast; §9.0 ping skip for fast.

Ship Slice 1 (+2) as this idea; 3/4 can be follow-on if the diff grows too large to review well.

## 5. Backward-compat & test plan

- **`standard` ≡ today:** derive `standard` to the *current* effective values (CrossReviewRounds
  default, MaxFixupCycles 3, reviewers capped at 2 — NOTE: today all non-implementers review, so
  with a 3-agent roster standard drops from 3→2 reviewers; call this out as the one intended
  behavior change and gate it so a 2-participant roster still works). **Safer MVP alt:** make
  `standard` = "≤2 but all if roster ≤3" so existing 3-agent decks are unaffected — decide in review.
- **absent `track:` ⇒ standard**, and standard = today (modulo the reviewer cap decision above).
- Tests: `ReadTrack` table test (fast/standard/deliberation/unknown/empty→standard); a
  `classify` table test covering every §4.0 row + fail-safe cases; a driver test per track
  asserting reviewer cap, CrossReviewRounds, MaxFixupCycles, and the ≥1 floor; a
  regression test that an existing `deliberation`/absent idea produces the identical cursor
  sequence as before (golden). All existing `internal/driver` tests MUST stay green.

## 6. Risks
- **Silent reviewer reduction** for existing decks if `standard` caps 3→2 — mitigate with the
  roster-size carve-out above; make it explicit, not silent.
- **Reviewer-cap ordering** could drop the only model-diverse reviewer — cap must preserve the
  diversity preference (keep a non-implementer-model reviewer first).
- **Two sources of truth** (declared `track:` vs classifier) — the driver MUST use only the
  declared track; classifier is advisory/warn to avoid nondeterminism.
- **Scope creep** into consensus/FINAL collapse — keep it out of the MVP.
