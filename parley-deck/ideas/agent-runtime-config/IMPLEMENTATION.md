---
idea: agent-runtime-config
status: implemented
implementer: codex
started: 2026-05-11
completed: 2026-05-11
branch: parley-deck-cli#feature-agent-runtime-config
head-commit: 1e1c757ddd74a3c81f8ceda3f0d186a4ae8adc7a
design-pr: https://github.com/feci/parley-deck-cli/pull/5
implementation-pr: https://github.com/feci/parley-deck-cli/pull/6
---

## Summary of work

Implemented the `agent-runtime-config` slice from `FINAL.md`:

- Added `cmd/parley/main.go` delegating to the existing `internal/app.Run`.
- Added `internal/config` for project-local TOML runtime configuration, including `parley-deck/agents.toml`, optional gitignored `parley-deck/agents.local.toml`, `PARLEY_HEADLESS_AGENT_CONFIG`, source tracking, and placeholder expansion.
- Extended `agents.Spec` / `agents.Discovery` with resolved runtime fields: sandbox, approval, model, reasoning/profile, speed, timeout, isolated-home env templates, external backend, notes, and field sources.
- Changed Codex built-in defaults from approval `never` to `on-failure` with `sandbox_mode = "workspace-write"`.
- Added `parley agents list` and `parley agents verify [--agent ID] [--full] [--yes]`, while keeping `discover|probe` as compatibility aliases.
- Added consent-gated full verification with probe artifacts under `parley-deck/meta/runtime-probes/<run-id>/` and a Codex Git smoke stage.
- Threaded resolved runtime values into `parley run`, the run event log, runner command construction, per-agent timeout handling, isolated-home setup, and round prompt text.
- Updated `.gitignore` and `docs/agent-runtime-configuration.md`.
- Added focused tests for config layering/source tracking, placeholder expansion, Codex defaults, agent list/verify behavior, Codex Git smoke failure reporting, and runner prompt runtime values.

## Deviations from FINAL.md

- The implementation branch is `feature-agent-runtime-config` instead of the example-style `feature/agent-runtime-config`. Creating the nested Git ref failed locally with `unable to create directory for .git/refs/heads/feature/agent-runtime-config`, while a flat branch ref worked. This does not change the implementation surface.
- Exact table formatting was resolved during implementation: `parley agents list` prints a compact fixed-width matrix plus source detail lines.
- `agents discover|probe` remain compatibility aliases to `list|verify` rather than being removed.

## Notes for reviewers

- Full verification can launch hosted agent CLIs and therefore requires `--yes` for non-local backends.
- On this machine, `go run ./cmd/parley agents list` reports the installed Codex binary as present but with a version probe error. That is expected from the local Codex installation state and is exactly the kind of distinction this slice is intended to surface.
- Verification run: `GOPATH=/private/tmp/parley-go GOCACHE=/private/tmp/parley-go-cache GOMODCACHE=/private/tmp/parley-go/pkg/mod go test ./...`
- CLI smoke: `GOPATH=/private/tmp/parley-go GOCACHE=/private/tmp/parley-go-cache GOMODCACHE=/private/tmp/parley-go/pkg/mod go run ./cmd/parley agents list`
