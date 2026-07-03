---
idea: track-aware-driver
author: claude-1
created: 2026-07-03
track: deliberation
participants: [claude-1, codex-1, hermes-1, antigravity-1]
status: complete
---

# Idea: track-aware-driver — deterministic CLI/driver enforcement of §4.0 tracks

## Problem / idea

Follow-up ratified by `meta-protocol-change-devx-speed` (shipped v1.32.0). That idea added the
`track: fast | standard | deliberation` **protocol text** (§4.0) but deferred the **deterministic
CLI/driver enforcement**. Today §4.0 is self-enforcing only by agents reading it. This idea makes
the CLI/driver actually route and gate by track.

## Scope to design (and slice into an MVP-first plan)

1. **`ReadTrack(ideaDir)`** — read `track:` from `00-prompt.md`, default `standard`, exactly
   mirroring `ReadStrictGate`/`ReadAutoImplement` in `internal/driver/transport.go`. Normalize
   fast|standard|deliberation; unknown/empty → standard.
2. **Track → `driver.Config` derivation** — map the §4.0 per-track table onto the existing
   `driver.Config` fields (`internal/driver/driver.go`): `CrossReviewRounds` (fast 0 / standard ≤2 /
   deliberation default), `MaxFixupCycles` (fast 1 / standard 2 / deliberation current), the
   reviewer-count requirement (fast 1 / standard 2 / deliberation all — locate where reviewer
   count is currently enforced: `Impl`/`ImplOps`, `internal/driver/impl.go`), and per-track
   timeouts (`internal/config/runtime.go` + `~/.parley` `[defaults.timeouts]`, ~5/15/30 min).
3. **Fail-safe classifier** — the §4.0 "deliberation-first, then fast, else standard; on doubt →
   stricter" rule, deterministic and script-checkable. **Key design question:** how does the CLI
   obtain the objective inputs (files touched, LOC, reversibility, security/data surface)? Options
   to weigh: (a) a `parley classify` command taking explicit flags (`--files N --loc N --security
   --irreversible --protocol-change …`) that prints the track — fully script-checkable; (b) infer
   from a git diff / idea metadata; (c) author declares `track:` and the CLI validates + warns on
   an obvious under-tier. Pick an MVP.
4. **Invariant enforcement gate** — the driver must HARD-REJECT (escalate, never silently proceed)
   any derived/declared config that would drop a §4.0 all-track invariant: 0 independent reviewers
   (non-solo), or a review path with no refutation. This is the safety-critical part.
5. **init/run templating + timeout seeding** — `parley init` seeds per-track `[defaults.timeouts]`
   into `~/.parley`/deck config; `parley run` accepts/reads `track:` and derives the Config; the
   `00-prompt` template already carries the `track:` line (shipped in v1.32.0).

## Codebase orientation (verified 2026-07-03)

- Frontmatter reads: `internal/driver/transport.go` (`ReadAutoImplement`/`ReadStrictGate`/
  `ReadRequireModelDiversity`; helper `readFrontmatterField` in `cursor.go`). Add `ReadTrack` here.
- Driver config + phase advance: `internal/driver/driver.go` (`Config`, `New` defaults, `Advance`,
  `roundComplete`), `internal/driver/impl.go` (impl/review/strict-gate loop, reviewer logic),
  `internal/driver/consensus.go`.
- Command wiring: `internal/app/app.go`, `internal/app/driver_impl.go` (builds `driver.Config`
  from 00-prompt fields — this is where track→Config derivation plugs in).
- Config/timeouts: `internal/config/runtime.go`.

## Constraints
- Enforce, never weaken, §4.0 invariants and the §14 human brake. `deliberation` behavior must
  remain EXACTLY today's full lifecycle (regression-free: existing driver tests stay green).
- Additive + backward-compatible: an idea with no `track:` behaves as `standard`, and `standard`
  must reproduce today's effective defaults (no silent behavior change for existing ideas).

## Non-goals
- Physical protocol-doc reorganization (that is `protocol-restructure-appendices`).
- New track semantics (owned by the shipped `meta-protocol-change-devx-speed`).

## What each participant delivers in `round-01/<agent-id>.md`
1. **Classifier input model** — recommend (a)/(b)/(c) above (or a hybrid) for the MVP, with the
   exact `parley classify` interface if you pick flags. Must be deterministic & fail-safe.
2. **Track → Config mapping** — the precise field values per track, and where each is threaded
   (cite files/functions). Flag any §4.0 behavior with NO current Config field (e.g. reviewer
   count, collapsed consensus/FINAL for fast) and propose the minimal code to add it.
3. **Invariant enforcement** — exactly how/where the driver hard-rejects a config that would drop
   non-solo/refutation.
4. **MVP slicing** — the smallest safe, shippable first slice vs later slices. Order by
   (safety + value) ÷ risk.
5. **Backward-compat & test plan** — how you guarantee `standard`/absent ≡ today, and which tests prove it.
6. **Risks / what could go wrong.**

Independent analysis; do not read others' round-01 first. Concrete, cite files. English only.
