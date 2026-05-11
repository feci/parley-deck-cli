---
idea: agent-runtime-config
status: final
author: codex
consensus-date: 2026-05-11
participants: [codex, claude, gemini, hermes]
---

## Final plan

Implement the `agent-runtime-config` slice as the next `parley-deck-cli` delivery after `hitl-tui-questions`.

The goal is to make agent runtime setup explicit, auditable, and verifiable before a cooperation run starts. The CLI should stop depending on hidden shell state for agent behavior and should show the effective runtime contract for Codex, Claude, Gemini, Hermes, and future agents.

## Implementation scope

### CLI entrypoint

Add a minimal `cmd/parley/main.go`:

```go
package main

import (
	"os"

	"parley-deck-cli/internal/app"
)

func main() {
	os.Exit(app.Run(os.Args[1:], os.Stdout, os.Stderr))
}
```

Do not introduce Cobra in this slice. Keep command routing aligned with the current `internal/app` and stdlib `flag` style.

### Runtime configuration files

Add project-local TOML configuration:

- `parley-deck/agents.toml`: checked-in shared defaults.
- `parley-deck/agents.local.toml`: user-local overrides, gitignored.

The implementation should use one focused TOML dependency. Do not support both TOML and JSON config formats in this slice.

Add these ignore rules:

```gitignore
parley-deck/agents.local.toml
parley-deck/meta/runtime-probes/
```

### Resolution order

Resolve each runtime field in this order, highest precedence first:

1. Explicit CLI flag where the command already has one.
2. `PARLEY_HEADLESS_AGENT_CONFIG`.
3. `parley-deck/agents.local.toml`.
4. `parley-deck/agents.toml`.
5. Built-in defaults.

Track the source for each resolved field so `parley agents list` can explain where values came from.

### Runtime model

Add an `AgentRuntime` / `ResolvedAgent` layer over the existing `agents.Spec` and `agents.Discovery` structs. Do not replace the adapter model wholesale.

The resolved runtime must cover at least:

- stable agent ID;
- CLI path or command candidates;
- headless invocation arguments;
- prompt mode;
- sandbox mode;
- approval policy;
- model;
- reasoning/profile;
- timeout;
- isolated-home environment template;
- external-backend class: `hosted`, `local`, or `unknown`;
- telemetry note;
- human-readable notes;
- field source metadata.

Use `cli-default` for model, reasoning, and profile when the CLI cannot prove supported options. Do not invent model names or reasoning levels.

### Built-in defaults

Change the built-in Codex runtime to:

```toml
sandbox_mode = "workspace-write"
approval_policy = "on-failure"
```

The CLI must preserve the no-silent-alternate-workdir rule: if a configured agent cannot write in the target repository, surface the failure instead of relocating protocol work elsewhere.

For agents without known portable sandbox or approval controls, record the literal configured flags when present or `cli-default` when unknown.

### Agent commands

Add the user-facing commands:

```text
parley agents list
parley agents verify [--agent ID] [--full] [--yes]
```

`parley agents list` must be cheap and non-spending. It should show path lookup, version probe, resolved runtime values, and clear columns for installed, version, and headless status.

`parley agents verify` without `--full` must remain cheap: path lookup and version probing only.

`parley agents verify --full` may run behavioral headless probes. It must require explicit `--yes` for agents whose `external_backend` is not `local`.

The current `agents discover|probe` commands may remain as compatibility aliases if that keeps the implementation smaller. User-facing docs should use `list|verify`.

### Full verification

Full verification probe artifacts must be written under:

```text
parley-deck/meta/runtime-probes/<run-id>/
```

That path must be gitignored. Probe artifacts must include a sentinel that lets the verifier confirm the intended agent wrote the current probe output.

For Codex, `verify --full --agent codex` must include this Git write smoke sequence in the target repository:

```bash
git status
git branch tmp-codex-git-test
git branch -D tmp-codex-git-test
printf test | git hash-object -w --stdin
```

If the Codex Git smoke needs approval, the retry must be handled by Codex's native `on-failure` runtime behavior. Do not build a generic stderr-pattern retry engine in the CLI.

### Run integration

`parley run` must show the resolved capability matrix before launching agents. `--yes` or `--auto` may skip an interactive confirmation prompt, but the effective runtime should still be visible in stdout or in the run event log.

Thread resolved runtime values into the runner instead of hardcoded defaults. The runner prompt and command construction should use the resolved model, reasoning/profile, timeout, sandbox, approval, and isolated-home settings.

### Documentation

Update `docs/agent-runtime-configuration.md` so it is the user-facing reference for:

- adding a cooperating agent;
- configuring Codex with `workspace-write` and `on-failure`;
- verifying Git write behavior in the target repository;
- external-backend consent;
- isolated-home behavior;
- why `agents.local.toml` is not committed;
- what `cli-default` means.

Do not update `parley-deck/meta/protocol-changelog.md`; this is a CLI behavior change, not a cooperation protocol change.

## Tests and verification

Add focused tests for:

- config layering and precedence;
- source tracking for resolved fields;
- TOML parsing of project and local files;
- placeholder expansion for root/temp paths;
- Codex built-in runtime defaults;
- `agents list` output on a controlled fake PATH;
- `agents verify` cheap path;
- `agents verify --full --agent codex` failure behavior when one Git smoke step fails;
- `parley run` using resolved runtime values rather than hardcoded defaults.

Before opening the implementation PR, run:

```bash
go test ./...
```

If the implementation adds or changes the CLI entrypoint, also verify a local binary invocation:

```bash
go run ./cmd/parley agents list
```

## Non-goals

- Interactive `parley agents add` wizard.
- Machine-global config under `~/.config/parley/`.
- Generic model discovery and automatic strongest-model selection.
- Per-run model/profile/timeout override flags beyond existing flags.
- Portable CLI-side approval retry engine.
- `parley resume` and richer runtime status/resume behavior.
- Automatic consensus progression policy.
- GitHub/GitLab transport automation.
- Token telemetry parsing and cost reporting.

## References

- Consensus: ./consensus.md
- Round 1: ./round-01/
- Round 2: ./round-02/
- User-facing note already started: ../../../docs/agent-runtime-configuration.md
