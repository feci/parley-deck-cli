# Agent Runtime Configuration

Parley Deck CLI should treat agent runtime setup as an explicit part of cooperation, not as a hidden user environment detail. A run is only effective when each selected agent has a known command, headless invocation mode, writable protocol-artifact path, timeout policy, and sandbox/permission behavior.

## Configuration Files

Runtime defaults are resolved from highest to lowest precedence:

1. Explicit CLI flags where a command already exposes one.
2. `PARLEY_HEADLESS_AGENT_CONFIG`.
3. `parley-deck/agents.local.toml`.
4. `parley-deck/agents.toml`.
5. Built-in defaults.

`parley-deck/agents.toml` is checked in and documents shared project defaults. `parley-deck/agents.local.toml` is gitignored and is the place for machine-specific paths or temporary overrides.

Example local override:

```toml
[agents.codex]
command = "/opt/homebrew/bin/codex"
sandbox_mode = "workspace-write"
approval_policy = "on-failure"
timeout_ms = 1800000

[agents.gemini]
isolate_home = true

[agents.gemini.isolated_home_env]
GEMINI_CLI_HOME = "{tempdir}"
```

Use `cli-default` when a model, reasoning, effort, or profile value cannot be proven from the target CLI. Do not invent model names in shared config.

## Codex Recommended Settings

When Codex participates as an agent, recommend this sandbox profile:

```toml
sandbox_mode = "workspace-write"
approval_policy = "on-failure"
```

Expected behavior:

- Git writes inside the selected repository should work normally.
- If a Git write fails because of sandbox or permission restrictions, retry the exact failed command through the runtime's approval/escalation path.
- Do not silently work around Git write failures by switching to a separate work directory unless the user explicitly approves that workaround.
- Prefer repository-local writes for protocol artifacts and implementation changes.

Useful startup smoke checks for Codex-style runs:

```sh
git status
git branch tmp-parley-git-test
git branch -D tmp-parley-git-test
printf test | git hash-object -w --stdin
```

`parley agents verify --full --agent codex --yes` includes the Codex Git smoke. The CLI does not implement a generic stderr-based retry engine; Codex approval retry is expected to come from Codex's own `on-failure` runtime behavior.

## Adding Agents

Users should add an agent only after they can describe these fields:

- Stable Parley agent ID, such as `codex`, `claude`, `gemini`, or `hermes`.
- CLI command or absolute executable path.
- Headless prompt mode, preferably stdin when supported.
- Arguments required to allow narrow workspace writes.
- Model, reasoning, effort, or profile settings when the CLI exposes them.
- Timeout for real round/review work; default to 30 minutes unless the agent is known to be fast.
- Whether the agent needs an isolated home directory for logs, sessions, auth metadata, or settings.
- Whether the agent may send task and repository context to an external backend.

Inspect the effective matrix before selecting an agent for quorum:

```sh
parley agents list
```

Before selecting an agent for quorum, Parley Deck CLI can verify:

- The executable is installed.
- `--version` or equivalent command returns successfully when available.
- A minimal headless prompt can create exactly one assigned artifact under `parley-deck/`.
- The agent does not need to write outside the workspace except through documented isolated-home setup.

Cheap verification:

```sh
parley agents verify
parley agents verify --agent codex
```

Without `--agent`, `parley agents verify` checks every configured agent and exits non-zero if any configured agent is missing or has a failing version probe. Use `--agent ID` when you intentionally want to verify only one agent.

Full behavioral verification can spend hosted-backend tokens and write probe files under `parley-deck/meta/runtime-probes/`, so it requires `--yes`:

```sh
parley agents verify --full --yes
parley agents verify --full --agent codex --yes
```

## Effective Cooperation Defaults

Recommended defaults for multi-agent cooperation:

- Use all installed, verified, headless-capable agents only when the user has not selected participants explicitly.
- Show the capability matrix before starting a new idea, review cycle, or implementation cycle.
- Keep model and reasoning fields as `cli-default` when they cannot be discovered; do not invent model names.
- Prefer isolated temporary homes for CLIs that otherwise write logs or sessions into user-global directories.
- Treat credentials, customer data, and unrelated private documents as explicit-consent data, not normal repository context.
- Keep every protocol artifact under `parley-deck/` in English.

## Isolated Homes

Some CLIs write sessions, logs, or auth metadata into user-global directories. Parley Deck can prepare temporary isolated homes for those CLIs while copying only the minimum known auth/config files needed for headless execution.

Configured isolated-home environment values may use placeholders:

- `{root}`: selected repository root.
- `{deck}`: selected repository's `parley-deck/` directory.
- `{tempdir}` inside `isolated_home_env`: temporary isolated home prepared for that specific agent launch.
- `{tempdir}` in command/path fields: the operating system temp directory resolved when configuration is loaded.

Probe outputs are not source artifacts and are ignored by Git.
