---
agent: agy
idea: rho-retro-tooling
review-round: 2
date: 2026-06-16
reviewed-commit: 816dde2
---
## Fix verification

I have verified all six agreed fixes from consensus.md on the diff between 984c757 and 816dde2:

1. **Hardened `propose` boundary**: Slug is validated as strict kebab-case via `reSlug`. `Lstat` is called on `ideas/<slug>` and fails closed if anything already exists at that path (protecting against pre-existing directories without a prompt file and symlink attacks). It creates the target directory with `os.Mkdir` (not `os.MkdirAll`) and opens the destination file using `O_CREATE|O_EXCL|O_WRONLY`. Tests have been added to cover non-kebab slugs, existing directories without prompts, and symlinked slug entries.
2. **Design-churn classification**: The `classify` function now uses `s.Rounds > 1` (instead of `> 2`), properly matching the scoring friction threshold, ensuring that 2-round ideas are classified as `design-churn` and retained in the coreset.
3. **Blocker detection**: `reBlocker` has been updated to match both `Status: ❌` and `Verdict: BLOCK|❌`.
4. **D4 signals**:
   - `Abandoned` work is extracted from the `status` frontmatter in `IMPLEMENTATION.md` or `00-prompt.md` and correctly bucketed under `blocked-or-abandoned`.
   - `RunFailures` are scanned from the structured logs at `parley-deck/runs/*/events.jsonl` (counting `agent.failed`, `agent.no_first_output`, `agent.stalled`, and `driver.error` events) and mapped to the idea via the `run.created` event slug, bucketed under `runtime-failure`. Scored, mapped, and tested.
5. **Neutral generated author**: The scaffolded `00-prompt.md` template now uses `author: <fill: author>` instead of hard-coding `claude`.
6. **Test helper round naming**: The `itoa` helper was removed, and round directories are named using `fmt.Sprintf("round-%02d", i)`.

All tests pass (`go test ./internal/retro/ ./internal/app/`) and the project compiles successfully.

## New findings

None. No regressions were introduced during fix-up cycle 1.

## Verdict

ACCEPT
