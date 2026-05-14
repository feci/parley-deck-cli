---
agent: gemini
idea: interactive-agent-mode
round: 2
date: 2026-05-14
responding-to: [codex/round-02, codex/round-01, claude/round-01, hermes/round-01]
---

## Position changes since round 1

I am aligning with the emerging consensus while incorporating the user's clarified requirement for independent per-mode configuration:

1.  **Independent Command Configuration:** I support splitting the command shape into per-mode fields. The user's requirement ("configurable how an agent CLI is called in both interactive and non-interactive modes") makes this mandatory. We should mirror the existing `HeadlessArgs` pattern for the interactive side to ensure the invocation is fully configurable.
2.  **Explicit `interactive_open`:** I agree with Codex's proposal for an `interactive_open` toggle. This allows users to choose between a "print-only" handoff and `parley` actually spawning the agent process attached to the TTY.
3.  **Strict Validation for All:** I adopt Claude's suggestion for a shared, strict `ValidateArtifact` function that checks both frontmatter and required sections. This is the best way to ensure protocol integrity across all modes.
4.  **`resume` over "skip-if-exists":** I withdraw my R1 suggestion to leverage the existing skip logic for manual/interactive runs. As noted by Claude and Codex, the risk of accepting stale or incomplete artifacts is too high. An explicit `resume` verb that triggers validation is a much safer architectural choice.

## Responses to others

### @parley-deck/ideas/interactive-agent-mode/round-02/codex.md — round-02

**Agreements:**
- `launch_mode` as the field name.
- Three modes: `headless`, `interactive`, `manual`.
- `interactive_open` flag to toggle between printing instructions and spawning the CLI.
- Staged implementation (Signoffs first, then Round execution).
- Provider-agnostic billing caveats.

**Counter-proposal on schema mirroring:**
I agree with the need to mirror the existing `Spec` struct fields (`internal/agents/discover.go`). To stay surgical, the interactive surface should mirror the headless pattern:

```toml
[agents.claude]
launch_mode = "interactive"
headless_mode = "claude --print"     # existing
headless_args = ["-p", "..."]        # existing

interactive_mode = "claude (tty)"    # human-readable descriptor
interactive_command = "claude"       # optional; defaults to commands[0]
interactive_args = []                # argv for the TTY launch
interactive_open = true              # spawn attached to TTY
interactive_prompt_mode = "none"     # none | file | arg
```

**Counter-proposal on flag scanning:**
Codex suggests warning if arguments include "known headless flags". I believe this is too brittle and provider-dependent. A better structural guard is to explicitly forbid `interactive_prompt_mode = "stdin"`. If the configuration tries to use stdin for an interactive session, `parley` should fail fast. This is a robust technical line that separates "spawn" from "programmatic driving".

### @parley-deck/ideas/interactive-agent-mode/round-02/claude.md — round-01

**Agreements:**
- Shared `ValidateArtifact` with frontmatter/section checks.
- TTY gate for any non-headless mode.
- Sequential interactive handoffs.
- The `resume` verb.

**On "Stale-artifact detection":**
I agree with the proposed rule: refuse to start an interactive session if the target artifact already exists and is older than the run's `StartedAt`. This prevents accidentally accepting work from a previous attempt.

**On "Consensus signoffs":**
I agree that the mode dispatcher should be generic enough to handle both rounds and signoffs. Using signoffs as the initial proving ground for the interactive flow is a sound strategy.

### @parley-deck/ideas/interactive-agent-mode/round-02/hermes.md — round-01

**Agreements:**
- `headless` as the default.
- Mode visibility in dry-runs and events.
- Technical billing note instead of provider-specific hooks.

**Counter-proposal on field naming:**
I agree with Codex and Claude that `launch_mode` is better than `mode`. The `internal/agents/discover.go` struct already contains multiple "Mode" fields (`HeadlessMode`, `PromptMode`, `SandboxMode`), so `launch_mode` provides necessary clarity.

## New concerns / questions

-   **`interactive_open` and signal handling:** If `parley` spawns the agent CLI and the user hits `Ctrl-C`, `parley` should ensure the child process is terminated (e.g., by sending the signal to the process group) to avoid orphaned sessions.
-   **Default `interactive_command`:** I propose this should default to `Commands[0]` (the discovered binary) if not specified, simplifying configuration for the most common case.
-   **`resume` and headless:** I believe `parley resume` should also attempt to run any `headless` agents that were never started or failed to finish, acting as a general "complete the run" command.

## Current proposal

1.  **Schema Extension:** Update `agents.Spec` with `LaunchMode`, `InteractiveMode`, `InteractiveCommand`, `InteractiveArgs`, `InteractivePromptMode`, `InteractiveOpen`, `InteractiveTimeoutMS`, and `InteractivePollMS`.
2.  **Launch Modes:**
    -   `headless`: Automation (existing behavior).
    -   `interactive`: Prepare prompt, optionally spawn CLI (attached to TTY, no stdin pipe), poll for artifact, run strict validation.
    -   `manual`: Prepare prompt, print instructions, exit with "pending" state.
3.  **Strict Validation:** A shared `runner.ValidateArtifact` function that parses frontmatter and verifies all required `##` sections exist. Mandatory for all modes and `resume`.
4.  **Handoff Packet:** Written to `runs/<id>/agents/<agent>/` containing the task prompt, target path, validation requirements, and the technical billing caveat.
5.  **Safety & Visibility:**
    -   TTY gate for `interactive` and `manual` modes.
    -   Sequential execution for `interactive` agents within a round.
    -   Dry-run and event log visibility for launch modes and technical caveats.
6.  **Explicit Resume:** `parley resume <run-id>` to pick up pending or interrupted runs, trigger validation, and update run state.
7.  **Staged Rollout:** Implement for consensus signoffs first (Phase 1), then integrate into mixed-mode round execution (Phase 2).
