---
agent: codex
idea: interactive-agent-mode
round: 1
date: 2026-05-14
---

## Summary

The CLI should support explicit launch modes rather than treating every participant as headless automation. I propose adding a per-agent `launch_mode` with values `headless`, `interactive`, and `manual`, plus command-level overrides, where interactive/manual modes produce a handoff prompt and rely on the user-driven agent session to create the canonical artifact before `parley` validates it.

## Proposed approach

Add configuration:

```toml
[agents.claude]
launch_mode = "interactive" # headless | interactive | manual
interactive_command = "claude"
```

Command overrides:

```text
parley run --agent-mode claude=interactive ...
parley consensus request-signoffs --agent-mode claude=interactive ...
parley agents list
```

Behavior:

- `headless`: existing behavior. `parley` invokes configured `headless_args`, streams/captures output, then validates the expected artifact.
- `interactive`: `parley` writes a handoff prompt file under the run directory, prints the exact command and target artifact path, optionally opens the configured interactive command, then waits for the artifact to appear or for the user to confirm completion. After that it runs the same validation path as headless mode.
- `manual`: `parley` writes the handoff prompt file and exits with next steps. A later `parley resume` or validation command verifies the artifact.

Interactive handoff should be transparent:

- show the provider/backend and launch mode in `agents list`, `run`, `resume`, and `request-signoffs --dry-run`;
- print that interactive use is user-driven and may be accounted differently by the provider;
- keep `--yes` gating for hosted headless launches, but do not treat an interactive handoff as automated provider launch unless `parley` actually starts a process.

Implementation shape:

- Extend `agents.Spec` / config override with `LaunchMode` and `InteractiveCommand`.
- Extend runtime matrix output.
- Add a small handoff prompt builder in `runner` for round artifacts and in the consensus request-signoffs path for signoffs.
- For round execution, a selected interactive agent should not run in the parallel headless goroutine as if it were automated. It should return a "waiting/manual" result or be handled by an orchestration path that prints/polls deterministically.
- For consensus request-signoffs, interactive mode can be simpler: print the prompt path and target consensus path, wait for the signoff block, then validate with the existing post-invocation validator.
- Store handoff prompt files under `parley-deck/runs/<run-id>/agents/<agent>/handoff-prompt.md` when a run exists. For command paths without a run ID, create a timestamped directory under `parley-deck/runs/`.

Tests:

- config parsing for `launch_mode` and `interactive_command`;
- `agents list` includes launch mode;
- dry-run shows which participants are headless vs interactive;
- interactive consensus signoff path validates an artifact created by a fake/manual writer without invoking headless args;
- manual mode writes prompt and exits with clear next steps;
- headless behavior remains unchanged.

## Concerns / open questions

- Should `interactive` actually spawn `claude` in the foreground, or should v1 only print the command and prompt path? I prefer v1 print-only/manual handoff because it avoids PTY complexity and avoids pretending automation is interactive.
- Should the default remain `headless` for all agents? I think yes, with project/local config allowing users to set Claude to `interactive`.
- How should `parley run` coordinate multiple participants when one is interactive? Sequential orchestration may be necessary for mixed modes; parallel headless plus interactive waiting is possible but more complex.
- Do we need a dedicated `parley validate` command for artifacts produced manually? It would make manual mode cleaner, but `resume` may be enough.

## Risks

- If the CLI silently turns headless work into interactive handoff, automation/CI will break. Avoid by making defaults explicit and preserving headless as default.
- If interactive mode spawns a process and tries to drive it programmatically, the distinction becomes ambiguous. Keep v1 as user-driven handoff.
- Mixed participant modes can make progress state harder to understand. The TUI/status views must label waiting/manual states clearly.
- Provider policy language can change. The CLI should describe modes technically, not promise billing outcomes.
