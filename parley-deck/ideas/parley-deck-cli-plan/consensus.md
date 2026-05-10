---
idea: parley-deck-cli-plan
drafted-by: codex
date: 2026-05-10
---

## Agreed decisions

- The CLI should be a durable local orchestrator over the existing Parley Deck filesystem protocol, not a separate hidden state machine.
- Runtime state should start as append-only JSONL events plus derived snapshots; SQLite is deferred until history querying or multi-run analytics require it.
- The TUI should be a live projection over durable run state. Killing or closing the TUI must not corrupt the run or imply that agent processes are cancelled.
- Default mode is human-in-the-loop. Automatic mode is explicit and policy-driven; it proceeds by consensus only when required decision fields and thresholds are satisfied, otherwise it falls back to HITL.
- The initial supported transport should be local directory only. GitHub PR and GitLab MR transports remain later extensions.
- Agent support should be adapter-based. Adapters launch local CLIs, stream output, watch artifact paths, parse telemetry when exposed, and render unknown token/cost data as `unknown` or `-`, never guessed.
- Initial commands should include `init`, `agents discover` or `agents probe`, `run`, `status`, `resume`, and `answer`.
- Non-interactive CLI commands must work without the TUI; the TUI is the rich live surface, not the only control path.
- The TUI should show current idea, phase, round, mode, timeline, agent statuses, live activity, artifact completion, HITL inbox, consensus state, elapsed time, and best-effort token/time statistics.
- The TUI should begin as one clear dashboard layout before adding advanced tabs or multi-idea dashboards.
- The language choice should be presented as a priority decision rather than hidden behind a single unqualified recommendation.

## Agreed trade-offs

- Go + Bubble Tea is the strongest option when the top priority is native self-contained executables, low operational packaging risk, and predictable cross-platform process supervision.
- TypeScript/Node.js + Ink is the strongest option when the top priority is npm-native distribution, React-style UI development, and direct `npx` ergonomics without a binary-downloader wrapper.
- Python + Textual is attractive for rich TUI iteration but is the weakest fit for the requested install story because executable packaging and `npx` distribution are less natural.
- If implementation must start immediately without a spike, Codex recommends Go + Bubble Tea as the provisional default. Claude will not block that choice, provided the npm path is treated as a first-class supported install route.
- If the user wants evidence before choosing, run a short decision spike: build the same tiny supervisor/TUI in Go and TypeScript, test Windows/macOS/Linux child-process handling, packaged executable behavior, and `npx` path, then choose from measured results.

## Open items deferred to user direction

- Whether to pick the stack now or spend roughly two days on a Go-vs-TypeScript packaging/TUI spike.
- Whether `npx parley-deck` should mean a real Node/TypeScript app or an npm package that downloads/dispatches a native binary.
- Whether the v1 TUI should optimize for maximum polish or for the smallest reliable dashboard.
- How strict token accounting must be when agent CLIs do not expose usage metadata.
- How much authority automatic mode should have before pausing for the user.

## Signoffs

### Signoff: codex — 2026-05-10
Status: ✅ ACCEPT
Notes: Accepts the consensus. The final plan should recommend a priority-based choice, with Go + Bubble Tea as provisional default if the user wants to start building before running a spike.

### Signoff: claude — 2026-05-10
Status: 🟡 ACCEPT-WITH-RESERVATIONS
Notes: Agree with the architecture (durable JSONL events, TUI-as-projection, HITL default, adapter-based agents, local-only transport for v1). Reservation on stack choice: I do not object to Go + Bubble Tea as provisional default, but the npm/`npx parley-deck` install route must remain first-class — either via a real Node implementation or via a thin npm wrapper that downloads the native binary with verified checksums and a clear offline fallback. If the wrapper path is chosen, that distribution mechanism should itself be validated during the spike (or early in implementation) before stack lock-in, since a broken `npx` story would undermine a stated user-facing requirement.
