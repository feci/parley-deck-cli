---
idea: runtime-status-resume
status: implemented
implementer: codex
started: 2026-05-12
completed: 2026-05-12
branch: parley-deck-cli#feature/runtime-status-resume
head-commit: 653bac85d34e03f3fc8f70da31b44caa81b7bd21
design-pr: https://github.com/feci/parley-deck-cli/pull/7
implementation-pr: https://github.com/feci/parley-deck-cli/pull/8
---

## Summary of work

Implemented the `runtime-status-resume` slice from `FINAL.md`:

- Added `internal/runstate` to project durable run state from `events.jsonl`, `run.created`, HITL question files through the `hitl` package, and protocol frontmatter.
- Moved the event-state reducer out of TUI-only ownership by making `internal/tui` delegate `ProjectEvents` and event summaries to `internal/runstate`.
- Added conservative run state semantics:
  - terminal `outcome=completed|incomplete|failed`;
  - non-terminal `liveness=unverified|idle`;
  - no process reattachment or live-process claim.
- Extended `parley status`:
  - default workspace overview now includes run summaries;
  - `--run RUN_ID`;
  - `--idea SLUG`;
  - unstable `--json`.
- Implemented `parley resume [--dir DIR] [--no-tui] RUN_OR_IDEA`:
  - exact run ID resolution first;
  - newest run for idea slug second;
  - `--no-tui` prints run detail;
  - default opens a resume TUI over durable run files.
- Added an explicit resume TUI option and footer so the resume path is not accidental `Done=nil` behavior.
- Added tests for run projection, status output, JSON output, resume `--no-tui`, and resume TUI exit behavior.

## Deviations from FINAL.md

- `internal/tui` keeps thin compatibility wrappers named `ProjectEvents` and `summarizeEvent` so existing TUI tests and call sites remain stable while the logic lives in `internal/runstate`.
- Default `parley status` currently lists up to 10 newest runs rather than only one latest run per idea. Detail flags still inspect a specific run or idea.
- Artifact reporting comes from event artifact paths in this slice. Deeper protocol artifact validation remains deferred.

## Notes for reviewers

- `resume` is intentionally read-only with respect to agent execution. It can answer HITL questions because that writes durable question/event files, but it does not retry or control subprocesses.
- JSON output is explicitly unstable and developer-oriented for this slice.
- Follow-up ideas from consensus remain deferred: retry, supervised run state, workspace watch, and runstate cache.

## Verification

- `GOPATH=/private/tmp/parley-go GOCACHE=/private/tmp/parley-go-cache GOMODCACHE=/private/tmp/parley-go/pkg/mod go test ./...`
- `GOPATH=/private/tmp/parley-go GOCACHE=/private/tmp/parley-go-cache GOMODCACHE=/private/tmp/parley-go/pkg/mod go run ./cmd/parley status`
- `GOPATH=/private/tmp/parley-go GOCACHE=/private/tmp/parley-go-cache GOMODCACHE=/private/tmp/parley-go/pkg/mod go run ./cmd/parley status --run 20260510T201528.433687000Z`
- `GOPATH=/private/tmp/parley-go GOCACHE=/private/tmp/parley-go-cache GOMODCACHE=/private/tmp/parley-go/pkg/mod go run ./cmd/parley resume --no-tui 20260510T201528.433687000Z`
- `GOPATH=/private/tmp/parley-go GOCACHE=/private/tmp/parley-go-cache GOMODCACHE=/private/tmp/parley-go/pkg/mod go run ./cmd/parley status --run 20260510T201528.433687000Z --json | jq -r '.run_id + " " + .idea_slug + " " + .outcome'`
- Manual resume TUI smoke: opened `parley resume 20260510T201528.433687000Z` in a PTY, verified the resume footer, and exited with `q`.
