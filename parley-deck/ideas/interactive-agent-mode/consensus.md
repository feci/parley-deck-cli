---
idea: interactive-agent-mode
drafted-by: codex
date: 2026-05-14
---

## Agreed decisions

- Add explicit per-agent launch modes with `launch_mode = "headless" | "interactive" | "manual"`. Default is `headless` to preserve existing automation and CI behavior.
- Keep `launch_mode` as the config field name. Do not use a bare `mode` key.
- Model interactive and headless invocation as separate command shapes. The existing headless fields remain; the interactive side mirrors the same style rather than introducing a nested schema from scratch.
- Extend agent specs and TOML overrides with interactive fields conceptually equivalent to:

```toml
[agents.claude]
launch_mode = "interactive"
headless_mode = "claude --print"
headless_args = ["-p", "--output-format", "text"]

interactive_mode = "claude (tty)"
interactive_command = "claude"
interactive_args = []
interactive_prompt_mode = "none" # none | file | arg
interactive_invoke = "print-only" # print-only | spawn-tty
interactive_timeout_ms = 1800000
interactive_poll_ms = 2000
interactive_notes = ""
```

- `interactive_command` defaults to the first resolved command for the agent when omitted. `interactive_args` is independently configurable.
- `interactive_invoke = "print-only"` writes handoff material, prints the command and target artifact path, then waits/polls according to the command being run. This is the default interactive behavior.
- `interactive_invoke = "spawn-tty"` starts the configured interactive command attached to the user's terminal. It must not pipe the task prompt through stdin, use a provider's non-interactive prompt flag, allocate a pseudo-terminal, scrape terminal output, or drive the session with an `expect`-style loop.
- `interactive_prompt_mode` describes how the user receives the prepared prompt: `none`, `file`, or `arg`. `stdin` is explicitly forbidden for interactive mode.
- `manual` writes the same handoff packet and exits with clear next steps. It does not poll or invoke an agent.
- Handoff packets live under the run directory when a run exists, otherwise under a timestamped directory in `parley-deck/runs/`. They include the task prompt, target artifact path, validation contract, resume/validate command, and technical usage caveat.
- Core CLI usage language must be provider-agnostic: describe that headless mode is programmatic execution and interactive mode is user-driven. Do not claim or promise a provider-specific billing outcome. Operator-owned provider notes may live in local config.
- Add one shared artifact validation path. It must validate round artifacts and consensus/review signoffs produced by headless, interactive, and manual/resume flows.
- For round artifacts, validation checks frontmatter identity and required sections. For consensus/review signoffs, validation reuses the append-only signoff validator and the existing canonical status rules.
- Add an explicit resume/validation flow for pending manual or interrupted interactive runs. Do not rely on "artifact exists" alone for interactive/manual completion.
- `parley resume <run-id>` validates pending artifacts and records completion events. For headless agents that were selected but never started, resume may rerun them headlessly.
- Headless agents may continue to run in parallel. Interactive handoffs are sequential within a run to avoid terminal/user confusion.
- `agents list`, runtime matrix output, dry-runs, launch summaries, and event logs show resolved `launch_mode` and the relevant command shape.
- Staged implementation order:
  1. Schema/config plumbing, handoff packet writer, shared validation, resume, mode visibility, and consensus signoff support.
  2. Mixed-mode round execution using the same primitives.
  3. Optional polish for `spawn-tty` signal handling and TUI states if not finished in step 1 or 2.

## Agreed trade-offs

- We prefer an explicit `interactive_invoke` enum over a boolean such as `interactive_open`, because it makes the launch behavior auditable and leaves room for future invocation styles without changing the field type.
- We do not use flag-name heuristics to detect headless behavior in interactive config. The structural guard is `interactive_prompt_mode`: interactive mode cannot deliver the task prompt through stdin.
- We keep default behavior conservative. Existing users get headless execution until they opt into local or command-line interactive/manual modes.
- `spawn-tty` is useful, but it is not required for the first useful delivery if print-only handoff plus polling/resume lands first.
- Strict validation may reject some incomplete human drafts, but this is necessary because interactive/manual modes broaden the ways an artifact can be created.
- Provider billing/accounting language stays technical and cautious. Exact account-plan guidance belongs in user docs or local notes that the operator controls.

## Open items deferred to implementation

- Final exact TOML field names must be fitted to the existing `agents.Spec` and `agentOverride` structs during implementation. The conceptual shape above is binding; spelling may be adjusted only for consistency with the current code.
- Define the exact pending/manual exit code.
- Define exact signal behavior for `spawn-tty`, with a preference for terminating the child process group on `Ctrl-C` unless implementation proves that unsafe.
- Decide how much of `resume` ships in the first implementation slice. At minimum it must validate pending signoff artifacts.
- Decide whether non-TTY print-only interactive handoff requires an explicit `--allow-interactive-noninteractive` flag. `spawn-tty` always requires a TTY.
- TUI state labels such as `waiting (interactive)` and `pending (manual)` may be implemented with the mixed-mode round execution slice if not needed for signoff support.

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: codex — 2026-05-14
Status: ✅ ACCEPT
Notes: Accept. The consensus captures the user requirement for independently configurable interactive and non-interactive CLI invocation while keeping the provider-accounting boundary explicit and compliant.

### Signoff: claude — 2026-05-14
Status: ✅ ACCEPT
Notes: All my round-02 locking targets are in: `launch_mode` field name with `headless` default; three modes; per-mode command/args split; `interactive_prompt_mode` enum with `stdin` forbidden as a structural guard (no flag-string heuristic); shared validation for round artifacts and signoffs; explicit `resume` verb; TTY gate; sequential interactive handoffs; provider-agnostic core language with operator-owned `interactive_notes`; staged delivery with signoffs first. The `interactive_invoke = "print-only" | "spawn-tty"` enum is a better shape than my original `interactive_open` boolean — auditable and extensible, as the trade-offs section notes. My open questions (Ctrl-C process group, resume re-running never-started headless agents, `interactive_command` defaulting to the first resolved command, exact TOML spelling deferred to implementation) are all carried with sensible leans in the deferred-items section.

### Signoff: gemini — 2026-05-14
Status: ✅ ACCEPT
Notes: ACCEPT. The consensus successfully incorporates the independent per-mode configuration and structural guards (no stdin for interactive sessions) that I advocated for in Round 02. The transition from a boolean `interactive_open` to the `interactive_invoke` enum is a superior architectural choice for auditability and future extensibility.

### Signoff: hermes — 2026-05-14
Status: ✅ ACCEPT
Notes: Accept. The design cleanly separates launch modes while preserving headless defaults, shared validation, and protocol invariants for multi-agent signoffs.
