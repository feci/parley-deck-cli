---
idea: rho-retro-tooling
status: implemented
implementer: claude
started: 2026-06-16
completed: 2026-06-16
branch: parley-deck-cli#feature/rho-protocol-proposal
head-commit: 6a15dde
design-pr: https://github.com/feci/parley-deck-cli/pull/52
implementation-pr: https://github.com/feci/parley-deck-cli/pull/52
---

## Summary of work

Implemented the `parley retro` tooling per the parent FINAL (D5) and COOPERATION.md
§13.

- **`internal/retro/retro.go`** — read-only mining of structured artifacts:
  - `Scan(root)` walks `parley-deck/ideas/*` and builds one `IdeaSignals` per idea
    from canonical artifacts only (round/review dir counts, `## Fix-up cycle`
    count in IMPLEMENTATION.md, NOT-FIXED occurrences, dismissed-findings count,
    `❌` blocker signoffs, and `inbox/*to-user*` escalations referencing the slug).
    No raw session JSONL (v1).
  - `score`/`classify` — deterministic failure-density weighting and a
    failure-mode bucket (blocked / escalation / fix-up-heavy / review-churn /
    design-churn / low-friction) used as the diversity key.
  - `Select(signals, k)` — type-diverse coreset: pass 1 takes one representative
    per failure type (excluding low-friction/zero-score), pass 2 fills by score.
  - `Diagnose(coreset)` — grouped, deterministic report; states the output is
    hypotheses, not findings.
- **`internal/app/retro.go`** — `parley retro scan|select|diagnose|propose`.
  scan/select/diagnose are strictly read-only (`--json` for scan/select, `--k`
  for select/diagnose). `propose --slug SLUG` scaffolds **only** a single new
  `ideas/<slug>/00-prompt.md` seeded from the diagnosis, fail-if-exists, rejects
  invalid slugs (path separators / dotfiles); it writes nothing else and never an
  existing/participant file.
- Wired `case "retro"` + a usage line in `internal/app/app.go`.

## Implementation plan / checklist

- [x] `internal/retro` package (Scan/score/classify/Select/Diagnose).
- [x] `parley retro` command (read-only default; propose = one 00-prompt.md).
- [x] Tests: `internal/retro/retro_test.go` (scan/score/classify, type-diverse
      select, diagnose) + `internal/app/retro_test.go` (propose writes only
      00-prompt.md, fail-closed on existing, rejects invalid slugs).
- [x] `gofmt` clean; `go build ./...`; full `go test ./...` green.
- [x] Smoke-run on this repo: top hard case = `runner-hardening-kindly` (the
      2-fix-up-cycle idea), with a type-diverse coreset.
- [ ] Reviewers' round-01 (codex/agy/hermes), fix-up, complete.

## Deviations from FINAL.md

- None. DPP/embeddings, raw-JSONL ingestion, re-rollout, best-of-N, auto-apply,
  and a persistent quarantine registry are deferred per the FINAL (cut from v1).

## Notes for reviewers

- Confirm read-only discipline: scan/select/diagnose write nothing; `propose`
  writes exactly one new `00-prompt.md` and fails closed if it exists (covered by
  `TestRetroProposeWritesOnly00PromptAndFailsClosed`).
- Confirm the failure-density signals and the deterministic, type-diverse coreset
  match §13 / D4, and that no raw session JSONL is read in v1.
- The generated 00-prompt carries `drafted-by: parley retro` and a `created:
  <fill: date>` placeholder for the facilitator — by design (the tool is not a
  participant and does not stamp authority).
