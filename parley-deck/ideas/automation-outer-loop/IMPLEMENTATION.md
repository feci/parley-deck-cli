---
idea: automation-outer-loop
status: implemented
implementer: claude-1
started: 2026-06-24
completed: 2026-06-24
branch: parley-deck-cli#loop-engineering-impl
head-commit: (this commit)
---

## Summary of work

Tier 4 (the outer loop), per `FINAL.md`: LE-8 (human-brake protocol §) + LE-9
(`parley loop tick`, the one-shot human-braked discovery command). `gofmt`,
`go build ./...`, `go vet`, `go test -count=1 ./...`, and the drift guard are green.

## Implementation checklist

- [x] **LE-8 — §14 human brake** — added `## 14. Automated outer loop (loop engineering)
  — the human brake` to BOTH `COOPERATION.md` copies (live + embedded default), byte
  identical (drift guard green). It binds any automated/standing/scheduled loop to
  discover-and-draft-only: §14.1 (MAY: draft `status: candidate`, no `participants:`
  claim), §14.2 (MUST NOT without a recorded human/full-quorum gate: promote, run,
  implement, land/merge/push, finalize, edit the roster, override consensus), §14.3
  (fail-safe — when uncertain, do less).
- [x] **LE-9 — `internal/loop` package** — `Candidate`, `Config` (Enabled, disabled by
  default), `TickResult`; `ReadConfig` (absent → disabled/no-error, malformed → fail
  closed), `ReadSignals` (absent → empty/no-error), `SlugFor`/`fingerprintOf` (explicit
  fingerprint kept whole; else 8-char sha256 of source+id), and `Tick` — drafts a
  `status: candidate` 00-prompt.md per not-yet-seen signal (dedupe by slug) and returns.
  It never staffs a quorum, runs, pushes, merges, or finalizes (§14 enforced in code).
- [x] **LE-9 — `parley loop tick` command** — `internal/app/loop_cmd.go`: `runLoop` +
  `runLoopTick` (`--dir`, `--signals`, `--enable`, `--json`). Resolves `<root>/parley-deck`,
  reads config (`--enable` force-enables this one-off run, still candidate-only), reads
  signals, calls `loop.Tick`, prints a human/JSON summary. Exits 0 when disabled or when
  zero candidates (cron/idempotent-safe). Wired `case "loop"` + a usage line in `app.go`.
- [x] **Tests** — `internal/loop/loop_test.go` (disabled-writes-nothing, enabled-drafts-
  candidate, dedupe, fingerprint default + explicit, ReadSignals/ReadConfig edge cases);
  `internal/app/loop_cmd_test.go` (disabled exits 0, `--enable` drafts candidate-only,
  `--json` shape).

## Deviations from FINAL.md

- None of substance. The MVP reads a signals file (as specified); live connectors
  (GitHub/CI APIs) and optional human-confirmed `parley run` from a promoted candidate are
  explicit out-of-scope follow-ups.

## Notes for reviewers

- **The §14 brake is the security boundary.** Try to break it: can any code path in
  `loop.Tick` / `runLoopTick` stand up a quorum, set `participants:`, flip a candidate to
  `round-01`, call `parley run`, push, merge, or finalize? It must not.
- **Disabled-by-default + fail-safe.** Confirm an absent config writes nothing; a malformed
  config fails closed (disabled), never silently enables.
- **Dedupe integrity.** Confirm the slug is stable for the same signal and that an existing
  candidate dir is skipped (no overwrite). Check an unexpected `os.Stat` error fails closed.
- **Prose vs frontmatter.** The candidate prompt mentions `participants:` in its Promotion
  note (prose); the invariant is that there is no frontmatter `participants:` key.
