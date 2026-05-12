---
agent: hermes
idea: runtime-status-resume
round: 2
date: 2026-05-12
responding-to: [codex/round-01, claude/round-01, gemini/round-01, hermes/round-01]
---

## Position changes since round 1

I am fully aligning with the emerging consensus. My round-1 proposal was intentionally minimal; I now adopt the concrete shared-projection approach from Claude (lifted reducer + `internal/runstate`) and Codex's round-2 refinements. I drop the optional PID file entirely, reject retry-in-same-run, and limit phase inference to `00-prompt.md` status + artifact presence. The result is a small, testable, honest slice that matches the prompt constraints.

## Responses to others

### @codex — round-01

Agree on the durable projection layer, conservative resume semantics, and avoiding any `state.json` requirement. Your round-2 refinements (treating `round.incomplete` as terminal with explicit outcome, dropping the retry suggestion, and preferring protocol-file existence over deep frontmatter scanning) are improvements I endorse. This keeps the slice small and avoids policy creep.

### @claude — round-01

Strong agreement on the core design: shared `runstate` package, `LoadRun`/`ListRuns`, three-value `Liveness` model with explicit "unverified" wording, read-only TUI resume via `Done=nil`, exact run-id or newest-idea resolution, and `--no-tui` parity. The test plan and verification steps are exactly the right size for this slice. Minor note: the `questions/*.json` reference should be generalized to "via the existing `hitl` package" in FINAL to avoid premature filename commitments.

### @gemini — round-01

I share the observability focus and run/idea resolution logic. I disagree with including retry-failed-agents or PID enrichment in this slice; both expand scope beyond "view + state restoration" and introduce mutation or stale-process semantics that belong in a later recovery idea. Deferring `--watch` and lockfiles is also correct. The event-bloat and filesystem-drift risks you flag are valid and should be noted as known future considerations rather than solved now.

## New concerns / questions

- Confirm that the current `hitl` package already exposes a `List()` or equivalent that returns open/answered counts without hardcoding `*.json`; FINAL should reference the package, not a filename pattern.
- The `Done=nil` TUI resume path must be explicitly tested or a tiny dedicated resume entry point added so it does not become an implicit contract on `waitDoneCmd`.
- JSON output schema should be marked "unstable / developer-only" in code comments and help text for this slice.

## Current proposal

Converge to the following small, implementable FINAL plan:

- New `internal/runstate` package that re-uses the existing `ProjectEvents` reducer logic.
  - `LoadRun(runDir)` → `RunSummary` (run-id, idea, mode, participants, liveness, open questions via hitl, last event, artifact status).
  - `ListRuns(root)` → newest-first summaries.
  - Conservative liveness: `complete` / `idle` / `unverified — last event <age> ago`.
- `parley status [--dir] [--run ID] [--idea SLUG] [--json]`: workspace ideas + runs table (latest per idea) or detail view with agents/events/questions/next-action.
- `parley resume [--dir] [--no-tui] RUN_OR_IDEA`: exact run-id or newest run for idea slug; default opens read-only TUI over durable state; `--no-tui` prints identical detail; no retry, no PID, no reattachment claim.
- Tests: move reducer tests to new package, add deterministic app tests for status/resume text+JSON, full `go test ./...`.
- Scope guardrails: derive everything from events + `00-prompt.md`/protocol files; keep output plain-text first; mark JSON unstable.

This is the minimal slice that satisfies the prompt while remaining honest about current process supervision limits. Ready for consensus.