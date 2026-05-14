---
idea: interactive-agent-mode
author: codex
created: 2026-05-14
participants: [codex, claude, gemini, hermes]
status: final
---

## Problem / idea

Design configurable agent launch modes for `parley-deck-cli` so each configured CLI agent can be used either through automated headless execution or through an explicit interactive handoff.

The immediate driver is Claude's announced Agent SDK / `claude -p` credit model starting June 15, 2026. `parley` currently invokes Claude as a headless/programmatic participant, which is correct for automation but draws from programmatic usage budgets. Users also need a first-class interactive mode where `parley` prepares the exact task, output path, and validation contract, then pauses while the user drives the agent manually in an interactive terminal session.

The goal is not to hide or bypass provider billing rules. The goal is to make the mode explicit, configurable, auditable, and safe:

- `headless`: `parley` invokes the agent command and validates the artifact.
- `interactive`: `parley` prints or writes a handoff prompt, opens or instructs the user to open the agent interactively, waits/polls for the target artifact, then validates it.
- `manual`: `parley` only prepares the prompt/artifact path and exits with clear next steps.

## Constraints

- Preserve Parley Deck ownership rules: each agent writes its own canonical artifact or signoff.
- Do not fake another participant's output.
- Do not attempt to misclassify programmatic usage as interactive usage.
- Keep billing/usage implications visible in dry-run and launch output.
- Make mode configurable per agent, with a sensible command-level override.
- Existing headless behavior must remain the default for automation and CI.
- Interactive mode must work for round artifacts and consensus/review signoffs.
- Interactive handoff must validate the final file the same way headless mode does.
- Avoid adding a full terminal multiplexer or pseudo-terminal manager unless the design proves it is necessary.
- Maintain compatibility with existing `agents.toml` / `agents.local.toml` runtime configuration.
- English-only for protocol artifacts and PR text.

## Non-goals

- Circumventing Anthropic, Google, OpenAI, or other provider billing policies.
- Automating an interactive session by secretly driving it through a PTY.
- Native provider billing integration.
- Full usage/cost dashboard.
- Replacing headless mode.
- Release packaging.

## Transport note

The active transport is `github-pr`. The protocol's preferred design branch would be `idea/interactive-agent-mode`, but nested refs have failed in this local Git sandbox before; this design PR uses flat branch `idea-interactive-agent-mode` while canonical files follow the protocol layout.
