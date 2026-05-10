---
idea: parley-deck-cli-plan
status: complete
implementer: codex
started: 2026-05-10
completed: 2026-05-10
branch: /Users/tomasfecko/AI_WORKSPACE/parley-deck/parley-deck-cli#master
head-commit: uncommitted
design-pr: n/a
implementation-pr: n/a
---

## Summary of work

Implemented the first Go-based CLI slice for `parley-deck-cli`:

- Added a Go module and `cmd/parley` entrypoint.
- Added standard-library command dispatch for:
  - `parley init`
  - `parley agents discover`
  - `parley status`
  - `parley run --no-tui [--auto] TASK`
  - `parley tui`
- Added protocol workspace initialization and frontmatter/status parsing.
- Added append-only JSONL run event storage under `parley-deck/runs/<run-id>/events.jsonl`.
- Added agent discovery/probing for Codex, Claude, Gemini, and Hermes.
- Added a first Bubble Tea/Lip Gloss dashboard shell that renders protocol status and discovered agents.
- Added focused tests for protocol frontmatter/workspace parsing and event-store append/load behavior.
- Added `.gitignore` entries for generated build artifacts.

This covers the bootstrap portion of M0, the basic event-store portion of M1, and the discovery/probe portion of M2 from `FINAL.md`.

## Deviations from FINAL.md

- The implementation does not yet launch agents for real protocol rounds.
- The TUI is a dashboard shell, not yet a polished live orchestration surface.
- GitHub PR and GitLab MR transports are not implemented in this first slice.
- `npx` packaging, Homebrew, and cross-platform release automation are not implemented yet.
- Token accounting is not implemented yet; the current code only records that telemetry support is adapter-specific.
- Gemini isolated-home execution is documented in discovery notes but not yet automated by the runner.
- `head-commit` is `uncommitted` because this repository has no commits yet.

## Notes for reviewers

- Build/test commands need writable Go caches in this sandbox:
  - `GOCACHE=/tmp/parley-go-build GOPATH=/tmp/parley-go go test ./...`
  - `GOCACHE=/tmp/parley-go-build GOPATH=/tmp/parley-go go run ./cmd/parley agents discover`
- `parley agents discover` found all four requested agents in this environment:
  - Codex
  - Claude
  - Gemini
  - Hermes
- A smoke run was created with:
  - `GOCACHE=/tmp/parley-go-build GOPATH=/tmp/parley-go go run ./cmd/parley run --no-tui --auto 'smoke implementation run'`
- That command produced `parley-deck/runs/20260510T194003Z/events.jsonl`.
- Cross-compile smoke checks passed for:
  - `GOOS=darwin GOARCH=arm64`
  - `GOOS=linux GOARCH=amd64`
  - `GOOS=windows GOARCH=amd64`

## Fix-up cycle 1
status: complete
completed: 2026-05-10
head-commit: uncommitted

### Fixes applied

- Changed `parley run` from a misleading run-only event into a partial but real workflow bootstrap:
  - requires an initialized workspace;
  - creates a new idea directory;
  - writes `00-prompt.md`;
  - creates `round-01/`;
  - records selected installed participants in the run event;
  - clearly states that agent execution is not wired yet.
- Removed the dispatcher-level 15-second timeout and replaced it with signal-based cancellation.
- Stopped `parley run` from silently initializing workspaces in arbitrary directories.
- Added clearer missing-workspace handling for `parley status`.
- Added `agents probe` as an alias for `agents discover`.
- Added placeholder dispatcher entries for `resume` and `answer` so their current state is explicit.
- Added telemetry details to `agents discover` output.
- Made run IDs nanosecond-resolution to avoid same-second collisions.
- Switched JSONL loading to `bufio.Scanner`.
- Updated the TUI to use Bubble Tea alt screen and dynamic panel widths from `tea.WindowSizeMsg`.
- Updated usage text to include `resume`, `answer`, and `version`.
- Re-ran cross-compile smoke checks for macOS, Linux, and Windows targets.

### Deviations from agreed fixes

- `parley status` is still frontmatter-oriented and does not yet compute full runtime phase, round, participant state, HITL inbox, or artifact completion from event logs.
- `resume` and `answer` are explicit stubs only; full durable resume and HITL answering remain M5 work.
- Agent execution, log streaming, artifact watching, and token parsing remain deferred to the next runner slice.

## Fix-up cycle 2
status: ready-for-re-review
completed: 2026-05-10
head-commit: uncommitted

### Fixes applied

- Added a first real runner layer:
  - selects installed participants from the idea roster;
  - launches all selected agents in parallel for `round-01`;
  - writes per-agent stdout/stderr logs under `parley-deck/runs/<run-id>/agents/<agent-id>/`;
  - checks whether each agent created its required `round-01/<agent-id>.md` artifact;
  - records `agent.started`, `agent.finished`, `agent.failed`, `agent.skipped`, `round.completed`, and `round.incomplete` events.
- Wired `parley run` to execute the runner after creating the idea and run event.
- Added headless launch specifications for Codex, Claude, Gemini, and Hermes.
- Automated Gemini isolated-home execution by creating a temporary `GEMINI_CLI_HOME` with copied OAuth files and minimal `oauth-personal` settings.
- Added a runner prompt builder that tells each agent exactly which protocol artifact it owns and forbids editing other agents' files.
- Added a runner unit test using a fake headless agent process, so the execution contract is verified without spending model tokens or invoking real CLIs.
- Kept `--no-tui` behavior as "run without opening the TUI after the round"; the default still opens the dashboard after execution.

### Deviations from agreed fixes

- Live TUI streaming is not implemented yet; the CLI runs agents and then opens the current dashboard.
- Token accounting is still adapter metadata only. The runner does not yet parse Codex/Claude/Gemini/Hermes token usage.
- HITL question capture and `parley answer` remain explicit stubs.
- Durable resume remains an explicit stub.
- I did not run `parley run` against the real installed agents in this verification pass, because the new command now launches all detected models. The runner behavior was verified with a fake agent test instead.

### Verification

- `GOCACHE=/tmp/parley-go-build GOPATH=/tmp/parley-go go test ./...`
- `GOCACHE=/tmp/parley-go-build GOPATH=/tmp/parley-go go run ./cmd/parley agents probe`
- `GOCACHE=/tmp/parley-go-build GOPATH=/tmp/parley-go GOOS=darwin GOARCH=arm64 go build ./cmd/parley`
- `GOCACHE=/tmp/parley-go-build GOPATH=/tmp/parley-go GOOS=linux GOARCH=amd64 go build ./cmd/parley`
- `GOCACHE=/tmp/parley-go-build GOPATH=/tmp/parley-go GOOS=windows GOARCH=amd64 go build ./cmd/parley`

## Fix-up cycle 3
status: complete
completed: 2026-05-10
head-commit: uncommitted

### Review inputs addressed

- Formal review: `review/round-02/claude.md`
- Advisory review: `../../inbox/gemini-to-all_parley-deck-cli-plan_runner-fixup-advisory-review.md`

### Fixes applied

- Added participant selection and a HITL confirmation gate:
  - `parley run --participants codex,claude ...` limits the selected agents;
  - unknown or unavailable agent IDs fail before any idea/run is created;
  - default HITL mode asks for confirmation before launching agents;
  - `--yes` or `--auto` explicitly skips the confirmation prompt.
- Reused the already discovered agent list when opening the post-run TUI, avoiding an immediate second discovery/probe pass.
- Changed runner event handling so store append errors are surfaced through `Result.ExitError` instead of being discarded.
- Added `agent.failed` events for setup-time failures before the child process is launched.
- Added `CompletedAt` and `Duration` population for early-return runner failures.
- Serialized JSONL `Store.Append` calls with a process-wide mutex to avoid concurrent append interleaving across platforms.
- Removed the unused `agents.Spec.IsolatedHomeID` field.
- Made Gemini isolated-home setup use `os.UserHomeDir()` and fail fast when neither OAuth files nor Gemini API key environment variables are available.
- Reduced round prompt duplication by relying on the stored `00-prompt.md` as the idea prompt.
- Updated CLI usage to include `agents discover|probe`, `help`, and the new `run` flags.
- Added runner tests for:
  - successful fake headless agent execution;
  - skip behavior when the artifact already exists;
  - failed child process event recording.

### Remaining deviations from FINAL.md

- Live TUI streaming is still not implemented; the dashboard opens after synchronous round execution.
- Token accounting is still adapter metadata only.
- HITL question capture and `parley answer` remain explicit stubs.
- Durable resume remains an explicit stub.
- Hermes remains supported because the user explicitly requested Hermes after the original finalized plan. It is still noted as less standardized than Codex/Claude/Gemini.
- I still did not run `parley run --yes` against the real installed agents in this verification pass. The runner contract is covered by fake-agent tests, and the default smoke check verified that HITL mode does not launch paid/model-backed agents without confirmation.

### Verification

- `GOCACHE=/tmp/parley-go-build GOPATH=/tmp/parley-go go test ./...`
- `GOCACHE=/tmp/parley-go-build GOPATH=/tmp/parley-go go run ./cmd/parley agents probe`
- `GOCACHE=/tmp/parley-go-build GOPATH=/tmp/parley-go go run ./cmd/parley run --no-tui 'non interactive hitl confirmation smoke'`
- `GOCACHE=/tmp/parley-go-build GOPATH=/tmp/parley-go GOOS=darwin GOARCH=arm64 go build ./cmd/parley`
- `GOCACHE=/tmp/parley-go-build GOPATH=/tmp/parley-go GOOS=linux GOARCH=amd64 go build ./cmd/parley`
- `GOCACHE=/tmp/parley-go-build GOPATH=/tmp/parley-go GOOS=windows GOARCH=amd64 go build ./cmd/parley`

### Review consensus

- Formal reviewer round: `review/round-03/claude.md`
- Review consensus: `review/consensus.md`
- Outcome: accepted with no remaining CRITICAL or MAJOR findings.

## Hermes integration follow-up
status: complete
completed: 2026-05-10
head-commit: uncommitted

### Fixes applied

- Enabled isolated-home execution for Hermes:
  - sets `HERMES_HOME` to a temporary writable directory;
  - creates writable `logs`, `sessions`, and `home` subdirectories;
  - copies existing `config.yaml`, `.env`, `auth.json`, and `SOUL.md` from the user's `~/.hermes` when present;
  - sets `HERMES_ACCEPT_HOOKS=1` and `HERMES_SESSION_SOURCE=parley` for the child process.
- Updated `parley agents probe` notes to show Hermes uses isolated `HERMES_HOME` for writable logs.

### Verification

- `hermes --version`
- `GOCACHE=/tmp/parley-go-build GOPATH=/tmp/parley-go go test ./...`
- `GOCACHE=/tmp/parley-go-build GOPATH=/tmp/parley-go go run ./cmd/parley agents probe`
- `GOCACHE=/tmp/parley-go-build GOPATH=/tmp/parley-go GOOS=windows GOARCH=amd64 go build ./cmd/parley`
- Initial Hermes smoke run failed because Hermes tried to write `/Users/tomasfecko/.hermes/logs/agent.log` from inside the sandbox.
- Follow-up Hermes-only smoke run succeeded:
  - run: `parley-deck/runs/20260510T201528.433687000Z/events.jsonl`
  - artifact: `parley-deck/ideas/hermes-integration-smoke-test-verify-that-hermes-can-participate-in-a-parley-deck-round-by-writing-its-own-round-01-artifact-only-2/round-01/hermes.md`
  - result: `round.completed`, `completed=1`, `total=1`.
