---
agent: claude
idea: parley-deck-cli-plan
round: 1
date: 2026-05-10
---

## Summary

Build `parley-deck-cli` as a thin, durable orchestrator on top of the existing
filesystem protocol in `parley-deck/`. It supervises N local agent CLIs in
parallel, persists every observable transition as an append-only event log
under `parley-deck/runs/<run-id>/`, and renders that log in a TUI. The TUI is
read-only on the protocol — the source of truth is files, not memory. Default
mode is human-in-the-loop; "auto" mode is the same supervisor with the HITL
prompt auto-answered by a configured policy (e.g. "advance on consensus,
otherwise wait").

Recommendation: **TypeScript/Node.js + Ink**, distributed both as `npx`
package and as standalone single binaries via `bun build --compile` for
macOS/Linux/Windows. Rationale below.

## Proposed approach

### Language choice — Go vs Python vs TS/Node

| Axis | Go | Python | TS/Node |
|---|---|---|---|
| Single binary | Best (native cross-compile via GoReleaser) | Worst (PyInstaller is large, brittle, sign/notarize pain) | Good (`bun --compile`, `pkg`/`@yao-pkg/pkg`) |
| `npx` runnable | Needs an npm wrapper that downloads a platform binary | Needs Node anyway → awkward | Native — just publish to npm |
| TUI ecosystem | Bubble Tea, tview (excellent) | Textual, Rich (excellent) | Ink, blessed (good) |
| Spawning agent CLIs | Fine | Fine | Native — `codex`/`gemini` are themselves Node, same runtime |
| Concurrency | Best (goroutines) | Decent (asyncio) | Decent (event loop, child_process) |
| Dev velocity for protocol files (JSON, YAML, markdown) | Verbose | Fastest | Fast |
| Match to existing parley-deck assets | Neutral | Neutral | Aligned (codex CLI, gemini CLI are Node) |

**Recommend TS/Node** because the user explicitly required `npx` runnability,
two of three discovered agent CLIs are Node-based (so Node is already on the
target machine in practice), and Bun's `--compile` gives a credible
single-binary story for the three OSes. **Go is the strong second choice** if
single-binary fidelity and startup time matter more than `npx` ergonomics; we
would then ship an npm wrapper that shells out to a downloaded binary.
Python is not recommended — it loses on distribution without giving us
anything we cannot get elsewhere.

Trade-off being made: we accept a slightly heavier compiled binary
(~40–80 MB with Bun) in exchange for one source tree powering both `npx` and
the standalone executables.

### Packaging / distribution

- Primary: publish `@parley-deck/cli` to npm → `npx parley-deck` works
  anywhere Node ≥ 20 is available.
- Standalone binaries: `bun build --compile --target=bun-{darwin,linux,windows}-{x64,arm64}`
  in CI, attached to GitHub Releases.
- macOS: signed + notarized via GitHub Actions; Homebrew tap formula
  pinning the release.
- Linux: tarball + checksum; optional `.deb` later.
- Windows: zip + Scoop manifest; signing deferred.
- Install one-liner: `curl -fsSL .../install.sh | sh` mirroring the
  GoReleaser pattern.
- Version is a single source of truth from `package.json`; binaries embed
  it at build time.

### Architecture (small first cut)

```
src/
  cli.ts              # commander entrypoint
  run.ts              # creates run-id, materializes state, starts supervisor
  supervisor.ts       # spawns agent adapters, watches files, emits events
  adapters/
    codex.ts          # codex exec
    claude.ts         # claude --print
    gemini.ts         # gemini --prompt
    base.ts           # capability flags: { canStreamTokens, canStreamProgress }
  protocol/
    phases.ts         # ideation -> implementation -> review -> consensus
    consensus.ts      # auto-mode decision policy
    events.ts         # append-only JSONL writer
  tui/
    app.tsx           # Ink root
    components/...
  state/
    store.ts          # reads JSONL, materializes per-agent view
```

Durable state lives in `parley-deck/runs/<run-id>/`:

- `events.jsonl` — append-only event log (only writer is the supervisor).
- `state.json` — periodic snapshot, derivable from events.
- `agents/<agent>/stdout.log`, `stderr.log`, `meta.json`.
- `questions/<id>.md` — HITL questions; user replies by editing or via TUI.

This means: crash the TUI, re-launch, state is intact. The TUI is a viewer
over files, not a controller of in-memory state. Matches the existing
`parley-deck/ideas/...` filesystem-as-protocol pattern.

### Agent adapter contract

Each adapter exposes:

```ts
interface AgentAdapter {
  name: string;
  detect(): Promise<{ ok: boolean; version?: string; path?: string }>;
  run(input: RunInput, hooks: Hooks): Promise<RunResult>;
  capabilities: {
    headlessFlag: string;       // 'exec' | '--print' | '--prompt'
    tokenReporting: 'none' | 'final' | 'stream';
    progressReporting: 'none' | 'lines' | 'json';
  };
}
```

Token accounting is **best-effort and per-adapter**. We never block a run
on missing token data — tokens are an optional column in the TUI. Adapters
parse final-usage lines if the CLI emits them; otherwise the token cell
shows `—`.

### Phases, rounds, parallelism

The CLI does not invent a new protocol — it executes
`parley-deck/COOPERATION.md`. A phase = a directory of expected artifacts;
a round = a numbered subdirectory. The supervisor:

1. Reads `participants` from the idea frontmatter.
2. Spawns all participants in parallel for the current round.
3. Watches their target file paths; an agent is "done" when its file
   exists with valid frontmatter.
4. Emits `agent-finished` events as files appear.
5. Advances phase/round when completion criteria are met (all done, or
   consensus reached in auto mode).

### TUI library + initial layout

Candidates: **Ink** (React-based, strong fit for our component model),
blessed/neo-blessed (lower-level), tui (Rust bindings — overkill).
Recommend **Ink** with `ink-spinner`, `ink-table`, `ink-text-input`.

Initial layout (80×24 minimum):

```
┌─ parley-deck-cli  idea: parley-deck-cli-plan  phase: round-01  mode: HITL ─┐
│                                                                            │
│  Agent     Status      Elapsed    In tok   Out tok   Activity              │
│  codex     ● running   00:42      1,204    312       writing summary…      │
│  claude    ● running   00:39      1,180    —         writing summary…      │
│  gemini    ◌ waiting   —          —        —         queued                │
│                                                                            │
├─ Questions (1) ────────────────────────────────────────────────────────────┤
│  [codex] "Confirm: should we keep token tracking optional?"  (a) answer    │
├─ Consensus ────────────────────────────────────────────────────────────────┤
│  auto-mode off                                                             │
├─ Logs (codex) ─────────────────────────────────────────────────────────────┤
│  > drafting proposed approach…                                             │
└─ q quit · a answer · tab focus · l logs · m mode ──────────────────────────┘
```

Three regions: top status bar, agent grid, bottom stack of HITL queue +
consensus + focused log + keybind footer. A single layout in v1 — splits
and tabs deferred.

### HITL vs auto mode

- HITL (default): when an agent writes a file matching
  `questions/*.md`, supervisor pauses that agent, surfaces it in the
  Questions pane; user answers in TUI or by editing the file; reply is
  written and agent resumes.
- Auto: a `consensus.ts` policy reads each agent's round artifact, applies
  a configured rule (default: unanimous on a `decision:` field; fallback
  to majority); if reached, the round is marked decided and next phase
  begins. If not reached after a timeout, fall back to HITL.

### Minimal first milestone

1. `parley-deck run <idea-path>` spawns the three adapters in HITL mode,
   produces `events.jsonl`, renders the TUI above. → verify: re-running
   after a kill resumes from the log.
2. Token parsing for codex (best case adapter) only. → verify: tokens
   appear for codex, `—` for the others.
3. Auto mode with unanimous consensus rule. → verify: a fixture idea with
   matching `decision:` fields advances without prompting.

## Concerns / open questions

These are the trade-offs I want the user to weigh before we start:

1. **`npx` vs single-binary priority.** If `npx` is non-negotiable, TS/Node
   wins. If you want zero Node dependency on user machines, Go wins and we
   pay an `npx` wrapper tax.
2. **TUI richness.** Ink gives us a React-style component model and is
   pleasant to build, but Bubble Tea (Go) and Textual (Python) are more
   polished out of the box. How much "wow" do you want in v1?
3. **Token tracking strictness.** Are you OK with token columns showing
   `—` for adapters that don't expose usage, or should we delay the
   feature until we have a uniform answer?
4. **State storage.** Append-only JSONL + snapshot is simple and
   diff-friendly; SQLite would give us better querying for long runs but
   adds a native dep. JSONL recommended.
5. **HITL surface.** Answer questions in the TUI only, or also via
   editing `questions/<id>.md` and committing? Editing is more durable
   and Git-reviewable; TUI is faster.
6. **Consensus rule defaults.** Unanimous on `decision:`, majority, or
   weighted by agent? Default unanimous, with override per idea.
7. **Auto-detect agents vs explicit config.** Should the CLI probe
   `$PATH` and offer what it finds, or require a `parley-deck.toml`
   with explicit adapter entries? Probing is friendlier; config is more
   reproducible.
8. **Resume semantics.** On re-launch, do we replay the event log
   passively (read-only view) or also restart non-finished agents?
   Restarting is convenient but risks duplicate spend.

## Risks

- **Headless output drift.** `codex exec`, `claude --print`,
  `gemini --prompt` may change output formats between versions; adapter
  parsing breaks silently. Mitigation: capability probe on startup +
  pinned tested versions in CI.
- **Windows TTY quirks.** Ink + Windows ConPTY is mostly fine but
  Bun-compiled binaries on Windows are the least-tested target. Plan
  for a manual smoke matrix.
- **Bun compile binary size & cold start.** ~40–80 MB and ~80 ms cold
  start. Acceptable for a CLI launched once per session, but worth
  measuring early.
- **Consensus illusions.** Auto-mode may declare consensus on agents
  that are confidently wrong in the same way. Mitigation: require
  explicit `decision:` fields and a non-empty rationale; default to
  HITL on close calls.
- **Token under-reporting hides cost.** Without a uniform usage API,
  users may run expensive auto-mode rounds blind. Mitigation: surface
  "tokens unknown" prominently and add a wall-clock-based soft cap.
- **Filesystem-as-protocol race conditions.** Two processes writing
  the same artifact path. Mitigation: only the supervisor writes
  protocol files; agents write to their own pre-assigned paths.
- **Scope creep into transports.** GitHub/GitLab transports are
  explicitly later. Risk that "later" leaks into v1. Mitigation: keep
  the adapter and transport interfaces cleanly separated from day one
  but ship only `local` in v1.
