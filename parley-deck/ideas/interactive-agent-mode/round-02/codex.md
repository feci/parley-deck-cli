---
agent: codex
idea: interactive-agent-mode
round: 2
date: 2026-05-14
responding-to: [claude/round-01, gemini/round-01, hermes/round-01]
---

## Position changes since round 1

The core proposal remains `headless | interactive | manual`, but the user has clarified a stronger requirement: the CLI must be configurable not only by mode, but also by how each agent CLI is invoked in interactive versus non-interactive use. That means the design should model both launch mode and per-mode command shape.

I also now think v1 should support two interactive behaviors:

- `interactive` with `open = false`: write handoff files, print the command and target artifact, then wait/poll for the user-driven artifact.
- `interactive` with `open = true`: start the configured interactive command attached to the user's terminal, without piping the task prompt through stdin or a non-interactive prompt flag. The user still drives the session.

This keeps the distinction honest: headless mode is programmatic execution; interactive mode is a user-driven session with explicit handoff material.

## Responses to others

### @claude — round-01

I agree with `launch_mode` instead of a generic `mode`; it is clearer in config and less likely to collide with model/provider "mode" naming. I also agree that validation parity is the load-bearing invariant and that stricter artifact validation should be shared by headless, interactive, and manual flows.

I partially disagree with "never invoke the agent CLI in interactive or manual mode" as an absolute rule. For `manual`, yes: it should only write instructions and exit. For `interactive`, the user specifically needs configuration for how the CLI should be called in interactive mode. The safe compromise is:

- allow an `interactive_command` / `interactive_args` command shape;
- only run it when a TTY is attached and the user or config explicitly enables opening it;
- attach it to the user's terminal;
- never feed the full prompt into stdin or a provider's non-interactive prompt flag;
- make billing/usage caveats visible without promising provider accounting outcomes.

I agree with separate `interactive_timeout_ms` and `interactive_poll_ms`, stale artifact guards, and a `resume` path. I would avoid provider-specific billing text like "counts against Claude.ai subscription" in the core CLI because provider policy can change. The CLI should say what it is doing technically and allow per-agent notes in local config.

### @gemini — round-01

I agree with the handoff workflow and with surfacing launch mode in dry-run and event logs. The validation loop is important: interactive mode should not accept an incomplete markdown file just because it exists.

I would not rely only on existing "skip if artifact exists" behavior for manual/resume. We need an explicit pending state or validation command so stale artifacts and half-written files do not look like successful agent output. Existing skip behavior can remain for backwards compatibility, but interactive/manual should record enough state to verify intent.

I also agree that interactive handoffs should be sequential. Headless agents can still run in parallel before or after the human handoff, but a single terminal should not be asked to supervise multiple interactive agents at once.

### @hermes — round-01

I agree on the default remaining headless and on applying the same mode logic to round artifacts and consensus/review signoffs. The canonical target artifact is the common denominator.

I disagree on using a bare `mode` key. The final schema should use `launch_mode` and reserve `mode` for future provider/model meanings. The config example should show both configured command paths:

```toml
[agents.claude]
launch_mode = "interactive"
headless_args = ["--print"]
interactive_args = []
interactive_open = true
```

Exact field names should still be fitted to the current `agents.Spec` shape during implementation, but the concept should survive: headless and interactive command invocations are independently configurable.

## New concerns / questions

- We need to inspect the current `agents.Spec` before freezing field names. If the existing schema already models command templates in a particular way, the implementation should extend that style rather than introduce a nested schema from scratch.
- `interactive_open = true` must fail fast outside a TTY unless an explicit override is present. CI should never hang.
- `interactive_open = true` must not accidentally become `claude --print <prompt>` with a different name. If arguments include known headless flags, validation should warn or fail when possible, while still staying provider-agnostic.
- The handoff prompt should include the target artifact path, required frontmatter, required sections, and the exact validation/resume command. It should not contain secrets beyond the context already sent to headless agents.
- Provider billing text should be cautious: "This is a user-driven interactive handoff. Provider billing/accounting depends on the provider and your account. Headless mode is programmatic execution."

## Current proposal

Adopt a staged design:

1. Add launch-mode resolution to agent config and CLI overrides.
   - Values: `headless`, `interactive`, `manual`.
   - Default: `headless`.
   - Config supports separate interactive command/args/open flags in addition to existing headless command/args.
   - `agents list`, `run --dry-run`, request-signoff dry-runs, and event logs show the resolved mode.

2. Add a shared handoff package writer and shared artifact validation.
   - Handoff location: under the run directory when available, otherwise a timestamped directory under `parley-deck/runs/`.
   - Validation applies equally to headless, interactive, and manual/resume output.
   - Validation checks frontmatter identity and required sections for round artifacts, plus signoff block rules for consensus/review signoffs.

3. Implement interactive/manual first for consensus request-signoffs.
   - This is the lowest-risk path because signoffs are naturally sequential and append-only.
   - It immediately solves the Claude interactive use case for approvals without disturbing round execution concurrency.

4. Implement mixed-mode round execution.
   - Headless agents can run in parallel.
   - Interactive agents are handled sequentially, either by printing handoff instructions and polling or by opening the configured TTY command when explicitly enabled.
   - Manual agents produce pending state and require `resume` or `validate`.

5. Document recommended local config.
   - For Claude after June 15, 2026, recommend setting `launch_mode = "interactive"` locally when the user wants an interactive session, and leaving `headless` only when they intentionally want programmatic execution.
   - Do not claim this avoids or changes billing; document the technical difference and the provider caveat.

If the other agents agree, consensus should lock the field name `launch_mode`, the three launch modes, the no-secret-PTY rule, the per-mode command configuration requirement, and the staged implementation order above.
