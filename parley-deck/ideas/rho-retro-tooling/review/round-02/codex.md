---
agent: codex
idea: rho-retro-tooling
review-round: 2
date: 2026-06-16
reviewed-commit: 816dde2
---
## Fix verification

Verified against `git diff 984c757..816dde2 -- internal/retro internal/app/retro.go internal/app/retro_test.go`.

1. `propose` write-boundary hardening is present. `reSlug` enforces strict lowercase kebab-case, `retroPropose` `Lstat`s `ideas/<slug>` and refuses any existing entry, creates the slug directory with `os.Mkdir`, and writes `00-prompt.md` with `os.OpenFile(O_CREATE|O_EXCL|O_WRONLY)`. Tests cover an existing slug dir without `00-prompt.md`, a symlinked slug, and non-kebab slugs.
2. `classify` now uses `s.Rounds > 1`, so a two-round idea is `design-churn`, has positive score, and survives `Select`.
3. `reBlocker` now matches both `Status: ❌` and `Verdict: BLOCK|❌`.
4. D4 signals are covered: abandoned frontmatter is represented as `Abandoned` and bucketed as `blocked-or-abandoned`; `scanRuns` reads structured `parley-deck/runs/*/events.jsonl` and counts `agent.failed`, `agent.no_first_output`, `agent.stalled`, and `driver.error` against the run-created idea; `runtime-failure` is scored, classified, and tested.
5. Generated `00-prompt.md` author is neutral: `author: <fill: author>`.
6. The round directory test helper now uses `fmt.Sprintf("round-%02d", i)` and the old `itoa` helper is gone.

Checks run:

- `go build ./...` passed.
- `go test ./internal/retro/ ./internal/app/` passed.

The workspace was on later `HEAD` `adc5622`, but `git diff 816dde2..HEAD -- internal/retro internal/app/retro.go internal/app/retro_test.go` was empty, so the checks exercised the reviewed files.

## New findings

None.

## Verdict

ACCEPT
