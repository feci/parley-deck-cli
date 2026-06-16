---
agent: codex
idea: rho-retro-tooling
review-round: 1
date: 2026-06-16
reviewed-commit: 984c757
---

## Summary

The implementation is mostly aligned with the approved v1 shape: `scan`, `select`, and `diagnose` are read-only; raw session JSONL is not ingested; ranking is deterministic from the current scan order; and `propose` writes only a prompt file on the normal path in the covered happy path. I found two merge-blocking issues in the fail-closed/write-boundary behavior and in coverage of ratified structured failure signals.

Verification run:

- `go build ./...` - pass
- `go test ./internal/retro/ ./internal/app/ ./internal/protocol/` - pass

## Findings

### [MAJOR] `propose` does not require the slug directory to be new

`internal/app/retro.go:96` builds `ideas/<slug>`, but `internal/app/retro.go:98` checks only whether `00-prompt.md` exists before `internal/app/retro.go:105` calls `os.MkdirAll` and `internal/app/retro.go:109` writes the file. That means an existing idea directory with other canonical artifacts but no `00-prompt.md` still accepts a retro prompt. It also leaves the write boundary dependent on any pre-existing directory entry at `ideas/<slug>`; for example, a symlinked slug directory can make the final `WriteFile` land outside the intended new idea directory.

This violates D5's explicit "non-existing slug" and fail-closed requirement, and it weakens the adversarial guarantee that `propose` can only create one fresh `ideas/<slug>/00-prompt.md`.

Concrete fix: validate the slug with the same kebab-case rule as the protocol, `Lstat` `ideas/<slug>` and fail if anything already exists there, create exactly that final directory with `os.Mkdir` rather than `MkdirAll`, and create `00-prompt.md` with `os.OpenFile(..., O_CREATE|O_EXCL|O_WRONLY, 0644)`. Add tests for an existing slug directory without `00-prompt.md`, a symlinked slug entry, and non-kebab slugs.

### [MAJOR] The scanner omits structured failure signals required by D4

`internal/retro/retro.go:22` through `internal/retro/retro.go:33` defines signals for rounds, review rounds, fix-up cycles, NOT-FIXED, dismissed findings, escalations, and blocker signoffs; `internal/retro/retro.go:85` through `internal/retro/retro.go:117` scores only those fields. It does not inspect `00-prompt.md` status for abandoned work, does not read structured run event logs under `parley-deck/runs/*/events.jsonl` for watchdog or `agent.failed` events, and does not capture drift-guard failures where they are recorded.

FINAL D4 explicitly includes blocked/abandoned work, drift-guard failures, and watchdog/`agent.failed` events in the v1 failure-density signals. Missing them means real hard cases can be ranked as low-friction or excluded from the coreset, so the retro pass is biased away from the operational failures it was ratified to mine.

Concrete fix: extend `IdeaSignals` and scoring for abandoned status, drift-guard failure evidence, and structured run failures. Parse `00-prompt.md` frontmatter for `status: abandoned`; scan `parley-deck/runs/*/events.jsonl` as structured event logs, not raw session transcripts, and attribute `agent.failed`/watchdog events by artifact path, slug, or run metadata; then cover these cases in `internal/retro/retro_test.go`.

### [MINOR] Generated prompts hard-code Claude as the author

`internal/app/retro.go:121` writes `author: claude` into every generated retro idea, while `internal/app/retro.go:126` separately says the prompt was drafted by `parley retro`. A user or another agent can run the command, and Phase 4 ownership depends on the `author:` field, so hard-coding Claude can assign idea stewardship to the wrong participant.

Concrete fix: either require an explicit `--author` value, derive the current agent from the app context if such a convention exists, or write a neutral placeholder such as `author: user` / `<fill: author>` so the facilitator confirms ownership before starting the round.

## Open questions

None.

## Verdict

ACCEPT-WITH-FIXES
