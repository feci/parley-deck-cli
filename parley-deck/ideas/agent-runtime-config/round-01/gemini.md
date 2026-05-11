---
agent: gemini
idea: agent-runtime-config
round: 1
date: 2026-05-11
---

## Summary
Propose an explicit agent configuration system for `parley-deck-cli` that replaces hardcoded discovery logic with a user-manageable runtime definition. The goal is to ensure agents operate in a predictable, narrow sandbox while providing clear escalation paths for restricted operations (like Git writes).

## Proposed approach

### 1. Unified Agent Configuration Schema
Replace the current `Spec` and `Discovery` structs with a serialized `AgentConfig` (stored in `~/.config/parley/agents.toml` or similar):

```toml
[agents.codex]
command = "codex"
headless_args = ["exec", "--cd", "{root}", "--sandbox", "workspace-write", "--ask-for-approval", "never", "-"]
prompt_mode = "stdin"
sandbox_mode = "workspace-write"
approval_policy = "on-failure"
timeout = "30m"
isolate_home = true
external_backend = true
```

### 2. Enhanced Agent Management CLI
Implement subcommands to manage this lifecycle:
- `parley agents add <id>`: Interactive wizard to define the fields above.
- `parley agents check <id>`: Perform automated verification:
    - Path/executable existence.
    - Version/Help command success.
    - **Smoke Test**: Attempt to write a temporary Git branch and a small artifact in `parley-deck/` using the headless mode.
- `parley agents matrix`: Display a table of all configured agents and their effective capabilities (Sandbox, Approval, Backend).

### 3. "Retry with Approval" Runtime Loop
Modify the runner to intercept failures that look like permission/sandbox blocks:
1. Agent attempts a Git operation (e.g., `git push`, `git commit`).
2. If it fails with a "permission denied" or "sandbox violation" signature:
    - If `approval_policy == "on-failure"`, the runner pauses and presents the exact command to the user.
    - Upon user approval, the runner executes the command in a higher-privilege context (or with sandbox temporarily relaxed).
3. Explicitly forbid "alternate workdir" workarounds unless the user has opted-in via `AgentConfig`.

### 4. Generalizing "Isolated Home"
Move the logic for `GEMINI_CLI_HOME` and `HERMES_HOME` from `runner.go` into a configurable template in `AgentConfig`. This allows users to add new agents that require similar isolation without modifying the Go source.

## Concerns / open questions
- **Error Signatures**: How reliably can we distinguish a sandbox violation from a genuine Git error across different agents? We may need agent-specific regex patterns for error matching.
- **Project-Local Config**: Should some settings (like `sandbox_mode`) be overrideable in a project-local `.parley.toml`? Some repos might require stricter isolation than others.

## Risks
- **User Friction**: Frequent "on-failure" approval prompts could degrade the "autonomous" feel of the agents.
- **Security**: The "Retry with Approval" path must be strictly limited to the exact failed command to prevent agents from escalating to arbitrary command execution.
- **Portability**: Isolated home strategies involving `os.MkdirTemp` and file copying might behave differently across OS/Filesystems (e.g., permissions on Windows).
