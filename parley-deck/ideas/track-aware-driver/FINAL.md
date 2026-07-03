---
idea: track-aware-driver
status: final
author: claude-1
consensus-date: 2026-07-03
participants: [claude-1, codex-1, hermes-1, antigravity-1]
---

## Final plan / specification

Make the CLI/driver deterministically enforce the §4.0 tracks (MVP = Slice 1 + Slice 2 of the
consensus). A new pure `internal/track` package classifies and derives per-track policy; the
driver reads the declared `track:` and applies reduced ceremony for `fast`/explicit-`standard`
while preserving today's behaviour for absent/`deliberation`.

## Purpose / user-visible outcome

- `parley classify [--files N --loc N --security …] [--declared T] [--json]` — a pure,
  script-checkable §4.0 classifier (deliberation-first, fail-safe); `--declared` exits 4 on an
  under-tier so CI can gate.
- An idea with `track: fast` runs with 1 model-diverse-preferred reviewer, no cross-review
  rounds, and a 1-cycle fix-up cap; `track: standard` (explicit) caps reviewers at 2 and fix-up
  at 2; `track: deliberation` and absent-`track:` behave exactly as today.

## Context & orientation

`internal/track/track.go` (new, pure): `Normalize`, `Classify(Inputs)`, `PolicyFor`.
`internal/driver/transport.go`: `ReadTrack`. `internal/driver/driver.go`: `Config.{Track,
MaxReviewers,MinReviewers}`, track derivation + `trackErr` in `New`, hard-gate in `Advance`.
`internal/driver/impl.go`: LE-11 guard uses `MinReviewers`. `internal/app/driver_impl.go`:
reviewer-set truncation by track. `internal/app/classify.go` + `app.go` dispatch: the command.

## Observable acceptance criteria

1. `parley classify` implements §4.0 verbatim (deliberation-first; fast only when
   reversible+mechanically-verifiable+≤5 files+≤300 LOC; else standard) and fails safe; verified
   by `internal/track` table tests + a CLI smoke run.
2. `track: fast` → `Config` MaxReviewers 1 / MinReviewers 1 / CrossReviewRounds 0 / MaxFixupCycles 1;
   `track: standard` (explicit) → 2/2/·/2; `deliberation` & absent → today's values. (`internal/driver` tests.)
3. The driver hard-rejects (escalates) `fast`+`auto_implement`, `fast`+`strict_gate`, and a
   non-solo config (0 available reviewers). (`internal/driver` tests.)
4. Refutation stays structural and non-optional on every track (no config disables review or the
   `## Refutation attempts` validator).
5. Backward-compat: absent-`track:` and `deliberation` reproduce today; the whole existing suite
   stays green (`go test ./...`).

## Idempotence & recovery

Additive, read-only derivation from `00-prompt.md` per `New`; no persisted state changes. Absent
track = today. Re-running `New`/`Advance` is deterministic.

## Known risks / de-risking

- Explicit `track: standard` intentionally reduces ceremony vs today — opt-in only; default
  template left absent. - Per-track timeout deferred (agents.TimeoutMS precedence needs careful
  threading — follow-up). - Reviewer truncation uses participant order (model-diversity-preserving
  order is a documented follow-up).

## References
- Consensus: ./consensus.md
- Rounds: ./round-01/{claude-1,codex-1,hermes-1,antigravity-1}.md
