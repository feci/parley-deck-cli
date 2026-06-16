---
idea: rho-retro-tooling
status: complete
implementer: claude
started: 2026-06-16
completed: 2026-06-16
branch: parley-deck-cli#feature/rho-protocol-proposal
head-commit: 816dde2
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

## Fix-up cycle 1 (review/round-01 → review/consensus.md)

All six agreed fixes applied:

1. **[MAJOR] propose write-boundary hardening.** `retroPropose` now validates the
   slug as strict kebab-case (`reSlug`), `Lstat`s `ideas/<slug>` and fails closed
   if anything already exists there (covers a pre-existing dir without
   `00-prompt.md` and a symlinked entry — Lstat does not follow links), creates
   exactly the new dir with `os.Mkdir` (not `MkdirAll`), and writes the prompt
   with `os.OpenFile(O_CREATE|O_EXCL|O_WRONLY)`. New tests: existing dir without
   prompt, symlinked slug, non-kebab/space/uppercase/double-hyphen slugs.
2. **[MAJOR] design-churn classification.** `classify` now uses `s.Rounds > 1`
   (was `> 2`), matching the `score` friction threshold, so a 2-design-round idea
   is bucketed `design-churn` (score > 0) and kept in the coreset instead of
   being dropped as low-friction.
3. **[MAJOR] blocker detection.** `reBlocker` now matches both `Status: ❌`
   (consensus signoffs) and `Verdict: BLOCK|❌` (reviewer files).
4. **[MAJOR] D4 signals.** Added `Abandoned` (from `status:` frontmatter in
   IMPLEMENTATION.md/00-prompt.md → bucketed with `blocked-or-abandoned`) and
   `RunFailures` (a new `scanRuns` reads structured `parley-deck/runs/*/events.jsonl`
   — NOT raw transcripts — counting `agent.failed`/`agent.no_first_output`/
   `agent.stalled`/`driver.error`, attributed to the idea via the run's
   `run.created` slug; new `runtime-failure` bucket). Scored + tested.
5. **[MINOR] neutral generated author.** The scaffolded `00-prompt.md` now writes
   `author: <fill: author>` (was hard-coded `claude`), so Phase-4 ownership is the
   facilitator's to confirm.
6. **[NIT] test helper.** Round-dir naming uses `fmt.Sprintf("round-%02d", i)`
   (removed the single-digit `itoa`).

Deviation note on fix 4: drift-guard failures are not persisted to any deck
artifact (the drift guard is a Go test), so there is nothing to mine for them in
v1; abandoned-status and run-event-log failures (which ARE persisted) are covered.

`gofmt` clean; `go build ./...` and full `go test ./...` green.

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
