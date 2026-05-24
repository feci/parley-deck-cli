---
from: codex
to: [claude, gemini, hermes]
idea: continuous-run-tui
round: 1
date: 2026-05-23
status: blocked
---

# Round 1 Agent Blockers

I corrected the `continuous-run-tui` participant roster to use the repository-configured agents from `parley-deck/agents.toml`: `codex`, `claude`, `gemini`, and `hermes`.

The real peer-agent attempts did not produce canonical round artifacts:

- `claude`: installed, but headless run failed with `Not logged in - Please run /login`.
- `gemini`: installed, but isolated `GEMINI_CLI_HOME` run failed because no auth method is configured in the isolated home and no `GEMINI_API_KEY`, `GOOGLE_GENAI_USE_VERTEXAI`, or `GOOGLE_GENAI_USE_GCA` environment variable is available.
- `hermes`: installed, but isolated `HERMES_HOME` run failed before agent work with `Invalid length for parameter modelId, value: 0`.

No files were written by these peer agents under `parley-deck/ideas/continuous-run-tui/round-01/`.

Current canonical state:

- `parley-deck/ideas/continuous-run-tui/00-prompt.md`
- `parley-deck/ideas/continuous-run-tui/round-01/codex.md`

Consensus is blocked until at least one configured non-facilitator agent can run and write its own artifact.
