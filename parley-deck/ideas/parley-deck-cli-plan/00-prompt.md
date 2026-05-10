---
idea: parley-deck-cli-plan
author: codex
created: 2026-05-10
participants: [codex, claude]
status: final
---

## Problem / idea

Create an implementation plan for `parley-deck-cli`: a self-contained cross-platform CLI/TUI for orchestrating Parley Deck multi-agent workflows.

The CLI must let a user start a task and see:

- which cooperation phase is active;
- which round is active;
- what each participating agent/model is doing in parallel;
- basic per-agent statistics such as elapsed time, input tokens, output tokens, and total tokens when available;
- questions from agents that need human-in-the-loop answers;
- automatic consensus progress when running in fully automatic mode.

The default operating mode should be human-in-the-loop. Fully automatic mode must also exist, where agents proceed through the cooperation protocol and decide by consensus when possible.

The plan must compare Go, Python, and TypeScript/Node.js for this project and recommend realistic packaging/distribution options. The CLI should be installable as a simple executable on Windows, Linux, and macOS, while also being runnable from Node.js when the user has Node installed, for example through `npx`.

The plan must suggest candidate TUI libraries and describe the initial TUI layout.

The plan should end by asking the user which trade-offs matter most before implementation starts.

## Constraints

- Keep the initial implementation small and goal-driven.
- Prefer durable protocol state over transient UI state.
- Do not assume all agent CLIs expose token counts; design token accounting as best-effort per adapter.
- Support local CLI agents first; keep room for GitHub/GitLab transport later.
- External backend disclosure is allowed for this design brief and necessary repository context, with secrets excluded.
- Current discovered candidate CLIs:
  - `codex` at `/Users/tomasfecko/.nvm/versions/node/v24.15.0/bin/codex`, version `codex-cli 0.130.0`, headless mode `codex exec`.
  - `claude` at `/opt/homebrew/bin/claude`, version `2.1.126`, headless mode `claude --print`.
  - `gemini` at `/Users/tomasfecko/.nvm/versions/node/v24.15.0/bin/gemini`, version `0.41.2`, headless mode `gemini --prompt`, excluded from quorum because the local invocation did not produce an artifact during recovery.

## Non-goals

- Do not implement the CLI in this idea.
- Do not design cloud-hosted orchestration.
- Do not require users to install all supported agent CLIs.
- Do not make token/cost tracking mandatory when a backend cannot expose it.
