---
agent: codex
idea: parley-deck-cli-plan
round: 1
date: 2026-05-10
---

## Summary

The strongest initial direction is a native-first CLI with a clear event model, durable run state, and a terminal dashboard that can observe several agent processes at once. I would choose Go with Bubble Tea for the default implementation, then publish an npm wrapper so `npx parley-deck` works by downloading or dispatching the correct platform binary. TypeScript/Ink is the best alternate if `npx` and JS hackability matter more than a single native binary.

## Proposed approach

Use a layered architecture:

- `protocol`: a state machine for Parley Deck phases, rounds, participant files, consensus, signoffs, HITL questions, and auto-mode policy.
- `runner`: per-agent adapters for launching local CLIs, streaming stdout/stderr, detecting questions, collecting exit status, and parsing token telemetry when available.
- `store`: durable append-only run events plus derived state. Start with JSONL event logs and a small state snapshot; move to SQLite only when querying/history becomes painful.
- `tui`: live dashboard over the event stream.
- `dist`: release binaries for macOS, Linux, and Windows, plus npm, Homebrew, and direct archive installs.

Recommended stack:

- Go core and TUI: best fit for single executable distribution, subprocess management, low startup overhead, and Homebrew/native package release.
- Bubble Tea, Bubbles, Lip Gloss: mature Go terminal UI stack with tables, viewports, progress bars, input areas, and full-screen UI patterns.
- GoReleaser or equivalent release automation: publish GitHub release artifacts, Homebrew tap/cask, checksums, and packages.
- npm wrapper: `@parley-deck/cli` with a tiny JS launcher that downloads or selects the platform binary, enabling `npx parley-deck`.

Alternative stack:

- TypeScript/Node.js with Ink: fastest if the CLI should be primarily npm-native and contributors know React. Packaging into true self-contained executables is weaker than Go, though Node single executable applications are improving.
- Python with Textual: excellent for rich TUI iteration, but packaging self-contained executables is heavier and `npx` is unnatural unless it is only a wrapper around a packaged binary.
- Rust/Ratatui: excellent binary story and performance, but slower implementation unless the team already wants Rust.

Initial feature set:

1. `parley init`: create `parley-deck/COOPERATION.md`, `ideas/`, `inbox/`, `meta/`, and optional local agent config.
2. `parley agents discover`: inspect configured/local CLIs and show capability matrix.
3. `parley run "<task>"`: create an idea, launch selected agents for round 1, and stream events to the TUI.
4. `parley status`: non-interactive summary of active ideas/runs.
5. `parley resume <run-or-idea>`: reopen a running or paused task.
6. `parley answer`: answer a pending HITL question from the terminal without opening the full TUI.
7. `parley auto`: explicit automatic mode for consensus-oriented runs.

TUI layout:

```text
┌ Parley Deck: parley-deck-cli-plan ─ phase: round-01 ─ mode: HITL ─ elapsed 08:14 ┐
│ Timeline: kickoff > round-01 running > consensus pending                         │
├ Agents ────────────────────────┬ Live Activity ──────────────────────────────────┤
│ codex   running  2m14s  9.1k   │ codex: writing proposed architecture...          │
│ claude  blocked  question 4.2k │ claude asks: should npm be primary install path? │
│ gemini  running  2m01s  5.8k   │ gemini: comparing TS vs Go packaging...          │
├ Protocol ──────────────────────┼ Stats ──────────────────────────────────────────┤
│ 00-prompt.md                   │ input 18.2k  output 4.1k  total 22.3k            │
│ round-01/codex.md ✓            │ cost n/a     wall 8m14s    failures 0            │
│ round-01/claude.md pending     │ per-agent sparkline / table                      │
│ round-01/gemini.md pending     │                                                 │
├ HITL Inbox ──────────────────────────────────────────────────────────────────────┤
│ [1] claude asks: Pick priority: easiest install, richest TUI, or fastest MVP?     │
└ Enter=respond  Tab=panel  a=auto  p=pause  r=retry  q=quit-but-keep-running ─────┘
```

Use a simple split layout first: left agent list/protocol tree, right live log/details, bottom HITL input. Add tabs later for `Overview`, `Agents`, `Artifacts`, `Logs`, and `Settings`.

HITL/auto semantics:

- HITL default: agents can ask questions; the run pauses only the blocked phase/agent; the TUI exposes a response box; answers are written into the next protocol artifact.
- Auto mode: agents proceed through rounds until consensus or a policy threshold is hit. User intervention is required only for configured escalation categories such as destructive operations, credentials, external private data, or unresolved blockers after N rounds.
- Both modes use the same durable event stream, so quitting the TUI must not kill running processes unless the user explicitly cancels them.

Token/statistics strategy:

- First-class fields: elapsed time, process status, retries, output bytes/lines, artifacts written.
- Best-effort fields: input tokens, output tokens, total tokens, cost, model, reasoning profile.
- Adapter-specific parsers: use JSON/stream events when a CLI exposes them; otherwise record `unknown` rather than guessing.
- Store raw telemetry events so parsers can improve without losing history.

## Concerns / open questions

- The npm path can mean two different things: a TypeScript app that runs on Node, or a Node wrapper that installs/runs a native binary. These lead to different implementation choices.
- Token accounting may be inconsistent because local CLIs differ in whether they expose usage metadata.
- Rich TUI interactions are useful, but the first release should still work well in plain non-interactive commands for CI and remote shells.
- The project needs a strict boundary between protocol files and orchestration runtime state so the CLI never corrupts canonical Parley Deck artifacts.

## Risks

- Building too much TUI before the protocol runner is reliable would create a nice dashboard over weak state.
- Automatic mode can silently make bad decisions unless policy is explicit and auditable.
- A Node-only executable strategy may disappoint users who expect one downloaded binary with no runtime.
- A native-only strategy may disappoint users who expect `npx parley-deck` to work everywhere without install friction.
