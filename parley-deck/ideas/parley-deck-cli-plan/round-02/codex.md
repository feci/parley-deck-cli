---
agent: codex
idea: parley-deck-cli-plan
round: 2
date: 2026-05-10
responding-to: [claude/round-01]
---

## Position changes since prior round

I agree with Claude that `npx` is a first-class requirement, not a nice-to-have. I also agree that the first version should use a simple append-only runtime event log and treat the TUI as a projection over durable state.

I do not agree that TypeScript/Node should be the default solely because `npx` matters. The user's requirements include both "completely self contained" and "installable as executable" on Windows, Linux, and macOS. A native binary is the cleaner default for that, while an npm wrapper can still satisfy the Node path.

Updated position: make the plan present two viable choices, with a recommendation:

- Recommended default if self-contained executable and operational reliability matter most: Go + Bubble Tea + npm wrapper.
- Recommended default if npm-native development and React-style TUI speed matter most: TypeScript + Ink.

The final plan should explicitly ask the user to choose which priority wins before implementation.

## Responses to others

### @claude

Your strongest point is that `npx` should not feel bolted on. I agree the npm path needs to be designed from day one, including a real `@parley-deck/cli` package and a tested Windows path.

Where I disagree is using Bun-compiled TypeScript as the default native-executable answer. That may be workable, but it creates a second runtime/toolchain dependency and asks us to validate Node/Bun executable behavior across terminal edge cases. For a supervisor that launches long-running child processes, handles terminal control, and must be trusted on all three major OSes, Go's native binary story is lower risk.

I also think "two of the discovered CLIs are Node-based" is not a decisive reason to build the orchestrator in Node. The orchestrator's hard parts are process supervision, durable state, terminal rendering, release packaging, and crash recovery; those are not materially easier because the child CLIs happen to be Node programs.

Counter-proposal:

1. Implement core in Go.
2. Build TUI with Bubble Tea, Bubbles, and Lip Gloss.
3. Publish native release artifacts with GoReleaser.
4. Publish npm package as a launcher/downloader for the native binary, with `npx parley-deck` treated as a required release path, not a demo.
5. Keep a TypeScript/Ink spike as the first alternative if the user says npm-native source is more important than native binary simplicity.

## New concerns / questions

- We should not overpromise exact token counts. The final plan should make "unknown" an explicit, visible state.
- We should avoid implementing SQLite in v1. JSONL events plus snapshots are enough until the run history needs complex queries.
- The TUI should not be the only interface. `run`, `status`, `answer`, `resume`, and `agents discover` should all work in plain terminal mode.

## Current proposal

Recommend Go + Bubble Tea as the default architecture, but frame it as a priority-based decision:

- choose Go if the main product promise is "download one executable and it works";
- choose TypeScript/Ink if the main product promise is "npm-native, hackable, React-style CLI";
- choose Python/Textual only if rich TUI speed matters more than distribution simplicity.

The final plan should ask the user to pick the priority order before implementation starts.
