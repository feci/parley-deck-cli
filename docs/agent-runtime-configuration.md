# Agent Runtime Configuration

Parley Deck CLI should treat agent runtime setup as an explicit part of cooperation, not as a hidden user environment detail. A run is only effective when each selected agent has a known command, headless invocation mode, writable protocol-artifact path, timeout policy, and sandbox/permission behavior.

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

Before selecting an agent for quorum, Parley Deck CLI should verify:

- The executable is installed.
- `--help` or equivalent command returns successfully.
- `--version` or equivalent command returns successfully when available.
- A minimal headless prompt can create exactly one assigned artifact under `parley-deck/`.
- The agent does not need to write outside the workspace except through documented isolated-home setup.

## Effective Cooperation Defaults

Recommended defaults for multi-agent cooperation:

- Use all installed, verified, headless-capable agents only when the user has not selected participants explicitly.
- Show the capability matrix before starting a new idea, review cycle, or implementation cycle.
- Keep model and reasoning fields as `cli-default` when they cannot be discovered; do not invent model names.
- Prefer isolated temporary homes for CLIs that otherwise write logs or sessions into user-global directories.
- Treat credentials, customer data, and unrelated private documents as explicit-consent data, not normal repository context.
- Keep every protocol artifact under `parley-deck/` in English.
