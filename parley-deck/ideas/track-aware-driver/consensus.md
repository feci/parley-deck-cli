---
idea: track-aware-driver
drafted-by: claude-1
date: 2026-07-03
---

## Consensus — design round 1 (unanimous on substance)

Four independent round-01 analyses (claude-1, codex-1, hermes-1, antigravity-1) converged
strongly. No cross-review round needed on substance; the one genuine design decision (how
`standard` reconciles with backward-compat) is resolved below.

### D1 — Classifier: hybrid, declared-track is the driver's source of truth
The driver reads the **declared** `track:` (deterministic); a separate **pure** classifier
computes the recommended track. MVP: `ReadTrack` + a fail-safe **contradiction reject**
(machine-checkable from frontmatter). `parley classify` (explicit flags, §4.0 verbatim,
`--json`, fail-closed "on doubt → stricter") lands as Slice 2. Git-diff inference is advisory-
only, deferred (all four agreed; codex/hermes explicitly rejected diff-as-truth for the MVP).

### D2 — Config shape (additive, backward-compatible)
Add to `driver.Config`: `Track string`, `MaxReviewers int` (reviewer-set cap; 0 = unlimited),
`MinReviewers int` (LE-11 auto-complete minimum; 0 = none). `driver.New` defaults
`MinReviewers = 2` when `AutoImplement && MinReviewers <= 0` (preserves LE-11 for Configs built
directly in existing tests). A single `applyTrack(cfg *driver.Config, ideaDir, runOpts)` in a
new `internal/app/track_config.go` is called from **all three** `driver.Config` construction
sites (`app.go` ~1154 `continueAuto`, ~1827 no-TUI, ~1881 TUI auto-drive) — one point of truth,
no three-site drift.

### D3 — THE key decision: absent `track:` ≡ today; explicit `track:` applies §4.0
`ReadTrack` returns `(track, present)` (like `readFrontmatterField`). This reconciles the
00-prompt backward-compat constraint with §4.0's "the table overrides the defaults":
- **absent / empty / unknown → today's behavior byte-for-byte** (`Track="standard"`,
  CrossReviewRounds default, MaxFixupCycles 3, MaxReviewers 0 = all reviewers, MinReviewers 2,
  timeout unchanged). Existing pre-v1.32.0 ideas are untouched.
- **explicit `track: standard`** → §4.0 standard: CrossReviewRounds `min(2, frontmatter)`,
  MaxFixupCycles 2, MaxReviewers 2, MinReviewers 2 (→1 if ≤2 participants, §4.0 degradation),
  timeout 15m.
- **`track: fast`** → CrossReviewRounds 0, MaxFixupCycles 1, MaxReviewers 1, MinReviewers 1,
  timeout 5m, `RequireModelDiversity=true` (fail-safe: fast's single reviewer must be diverse).
- **`track: deliberation`** → today's defaults (constraint: deliberation = today's full
  lifecycle) + explicit `Track="deliberation"` + timeout 30m.

### D4 — Reviewer count is the core enforcement point
- `newDriverImplOps` (`driver_impl.go`) truncates the built `reviewers` slice to `MaxReviewers`
  when `> 0`, preserving the model-diversity ordering (never drop the only diverse reviewer).
- The LE-11 guard (`impl.go:240`) generalizes from hard-coded `< 2` to `< d.cfg.MinReviewers`.
- Non-solo is fail-closed: `applyTrack` errors if a track needs ≥1 reviewer but < 1 distinct
  non-implementer exists; the existing `len(o.reviewers)==0` guard stays as the runtime backstop.

### D5 — Invariants the driver hard-rejects (escalate, never silently proceed)
- **Non-solo:** any derived config with 0 available independent reviewers → error.
- **Refutation stays structural & non-optional:** no flag disables review or the
  `## Refutation attempts` validator (`runner.ValidateReviewArtifact` / `phase58`); the driver
  ensures the review path runs (Phase 6) — it cannot be skipped on any track.
- **Contradiction reject:** `track: fast` + `auto_implement: true` → error; `track: fast` +
  `strict_gate: true` → error (both are §4.0 deliberation triggers).

### D6 — MVP scope (this idea) vs deferred
- **MVP = Slice 1 + Slice 2:** ReadTrack; Config fields + `New` default; `applyTrack` at all
  three sites; reviewer truncation; LE-11 generalization; non-solo + contradiction gates;
  `parley classify` pure command (`internal/track`); write `track: standard` into newly
  created 00-prompt (`protocol.CreateIdeaWithExclusions`).
- **Deferred (documented, own follow-up slices):** §9.0 ping-skip for fast (needs preflight
  reordering); collapsed consensus/FINAL for fast (consensus-package change); per-phase human
  gates; per-track `[defaults.timeouts]` seeding + `roundDeadline` per-track; mid-idea upgrade
  via diff scan. These are ceremony/nicety, not safety — the invariants above hold without them.

### D7 — Backward-compat is the acceptance bar
`deliberation` and absent-`track:` MUST produce today's cursor sequence (golden). All existing
`internal/driver` + `internal/app` tests stay green. New tests: `ReadTrack` table; `applyTrack`
per-track values incl. absent≠explicit-standard; reviewer truncation; MinReviewers guard
(fast=1 completes, fast=0 escalates, standard=1 escalates); classifier table (every §4.0 trigger
+ boundary fail-safe).

## Agreed trade-offs
- Explicit `track: standard` intentionally reduces reviewers/fix-up/timeout vs today — but only
  for ideas that opt in by declaring it; absent stays today. Escalation at a cap is correct, not
  a regression (author can declare `deliberation`).
- `auto_implement`/`strict_gate` force `deliberation` (fast is rejected) — §4.0 normative, not
  our invention; not a regression (those ideas already run the full lifecycle today).

## Open items deferred to implementation
- Exact placement of `MaxReviewers` truncation vs the model-diversity selection order.
- Whether `standard`'s `CrossReviewRounds` reads `min(2, frontmatter)` or a flat 2 (lean: min).

## Comparison & blind spots
- Unanimous: hybrid classifier, Config field additions, applyTrack single-point, LE-11
  generalization, MVP=fast-first, backward-compat bar.
- hermes uniquely nailed the **absent-vs-explicit-standard** reconciliation (D3) — adopted.
- codex uniquely enumerated the exact construction sites + `runcontrol.CreateOptions.Track` +
  that `CreateIdeaWithExclusions` doesn't yet write `track:` — folded into MVP.
- antigravity uniquely flagged the `auto_implement`-forces-deliberation tension (D5 contradiction reject).
- Blind spot: timeout precedence (`agents.TimeoutMS` overrides `Options.Timeout`) means the
  per-track timeout only bites if `runOpts.Timeout` is actually threaded — must be tested (codex R, hermes §2.6).

## Signoffs
<!-- Authored by each participant in signoffs/<agent>.md; assembled here verbatim. -->

### Signoff: hermes-1 — 2026-07-03
Status: ✅ ACCEPT
Notes: Consensus faithfully captures my round-01 (D1–D7 verified point-by-point, incl. the
absent-vs-explicit-standard reconciliation and the timeout-precedence must-test). Safe to
implement as MVP Slice 1+2.

### Signoff: antigravity-1 — 2026-07-03
Status: ✅ ACCEPT

### Signoff: codex-1 — 2026-07-03
Status: 🟡 ACCEPT-WITH-RESERVATIONS
Notes: Safe MVP; captures the safety-critical positions. Reservation: consensus applies §4.0
to explicit `track: standard` (not byte-for-byte with today). Accept as a deliberate opt-in
**provided tests lock absent-track compatibility and explicit-standard behavior separately.**
(Implementer note: satisfied — TestNewAbsentTrackIsLegacy, TestNewDeliberationIsLegacy, and
TestNewExplicitStandardAppliesCaps assert these independently; the default 00-prompt template
was left absent, so no new idea silently changes behavior.)

### Signoff: claude-1 (facilitator/implementer) — 2026-07-03
Status: ✅ ACCEPT — unanimous on substance (2×✅, 1×🟡 whose reservation the test plan resolves).
Proceeding to Phase 5 (MVP Slice 1+2) → Phase 6 review.
