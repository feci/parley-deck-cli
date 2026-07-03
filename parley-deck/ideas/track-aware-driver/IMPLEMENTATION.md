---
idea: track-aware-driver
status: ready-for-review
implementer: claude-1
track: deliberation
started: 2026-07-03
head-commit: (pending)
---

## Progress

MVP (Slice 1 + Slice 2) implemented on branch `track-aware-driver-impl`. Full suite green
(`go test ./...`), `go vet` clean.

### What landed
- **`internal/track/` (new, pure, no deps):** `Normalize`, `Classify(Inputs)` (§4.0
  deliberation-first, fail-safe), `PolicyFor` (per-track knobs + contradiction/non-solo errors).
  Fully table-tested (`track_test.go`).
- **`internal/driver/transport.go`:** `ReadTrack(ideaDir) (track.Track, bool)` — returns
  `present=false` for absent/empty/unknown (legacy path), mirroring `ReadStrictGate`.
- **`internal/driver/driver.go`:** `Config.{Track, MaxReviewers, MinReviewers}`; `New` derives
  the track policy from `cfg.IdeaDir` (applies overrides only for an explicit `fast`/`standard`;
  records a `trackErr` for a §4.0 contradiction / non-solo violation); `New` defaults
  `MinReviewers=2` (preserves today's LE-11 `< 2`). `Advance` hard-gates on `trackErr` (escalates).
- **`internal/driver/impl.go`:** LE-11 auto-complete guard generalized `< 2` → `< d.cfg.MinReviewers`.
- **`internal/app/driver_impl.go`:** `newDriverImplOps` truncates the reviewer set to the
  track's `MaxReviewers` (fast=1, standard=2; absent/deliberation keep all).
- **`internal/app/classify.go` + `app.go`:** `parley classify` command (pure; `--json`;
  `--declared` under-tier check → exit 4).

### Deviations from consensus (recorded per §4 Phase 5)
- **D2 placement — track derivation lives in `driver.New` (driver-internal), NOT an `applyTrack`
  helper at the three `app.go` `driver.Config` sites.** Rationale: `New` is the single
  construction chokepoint every site already calls, so this is a stronger single-point-of-truth
  than editing three call sites (eliminates the three-site-drift risk hermes/codex both flagged)
  and needs zero `app.go` changes. Same observable result. Reviewer truncation reads the track in
  `newDriverImplOps` (same deterministic source).
- **Per-track timeout DEFERRED.** Consensus listed timeout in the mapping; hermes flagged the
  `agents.TimeoutMS`-overrides-`Options.Timeout` precedence as a must-test and a threading hazard.
  Deferred to a follow-up (`track-timeouts`) rather than shipped half-wired. Not a safety item.
- **Default `00-prompt` template left ABSENT (no `track: standard`).** codex's 🟡 reservation
  asked that absent stay byte-for-byte today; templating `standard` into every new CLI-created
  idea would silently change new-idea behaviour. Kept absent = today; auto-defaulting the template
  is a deferred policy decision. (`parley classify` + explicit `track:` are the opt-in path.)
- **Reviewer truncation uses participant order** (not model-diversity-preserving order) — a
  documented follow-up nicety; the non-solo floor and the existing model-diversity gate still apply.

### Deferred (documented; own follow-ups) — matches consensus D6
§9.0 ping-skip for fast; collapsed consensus/FINAL for fast; per-phase human gates; per-track
timeouts + `roundDeadline`; mid-idea upgrade via diff scan.

## Verification
- `go test ./...` → all green (existing driver/app suites unchanged — backward-compat bar).
- New: `internal/track` table tests; `internal/driver/track_test.go` (fast/standard/deliberation/
  absent + contradiction + non-solo escalation).
- CLI smoke: `parley classify --files 2 --loc 20 --reversible --mechanically-verifiable` → `fast`;
  `+ --security` → `deliberation`; `--declared fast --files 30` → exit 4.
- `go vet ./internal/{track,driver,app}/...` clean.

## Observable acceptance criteria status
1. classify §4.0 verbatim + fail-safe — **met** (track tests + smoke).
2. per-track Config values (fast/standard/deliberation/absent) — **met** (driver track tests).
3. hard-reject contradiction + non-solo — **met** (driver track tests).
4. refutation structural/non-optional every track — **met** (unchanged validators; no bypass added).
5. backward-compat, full suite green — **met**.
