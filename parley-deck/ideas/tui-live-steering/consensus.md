---
idea: tui-live-steering
phase: consensus
drafter: claude
date: 2026-06-06
participants: [claude, codex, agy, hermes]
---

## Consensus

After round-01 (4 independent positions) and round-02 (cross-review), all four
participants agree on the design below. No open disagreements remain. codex's `Handle`
design is the runner spine; agy owns the UX; hermes the keymap/concurrency. The
autocomplete-as-dedicated-`suggest`-mode override (vs reusing the modal picker) is
accepted by all; agy's error-surfacing/short-terminal/styling refinements and
codex+hermes's implementation cautions are folded in.

### Goal
Make `parley tui` an active console: (1) slash-command **autocomplete**, (2) **kill** one
running agent, (3) a **steer round-trip** so typing to an agent sends it and shows the
reply. Today steers are write-only, agents are one-shot, and there's no per-agent kill.

### A. Slash-command autocomplete (pure TUI)
- A single `commandSpecs` table `{Name, Usage, RequiresRun, OpensPicker}` is the ONE
  source for `runCommand`, `/help`, the suggestion list, and hints.
- A dedicated **non-modal `suggest` sub-mode** (`suggest bool`, filtered command slice,
  `suggestIndex int`) — NOT the modal `picker` (which clears/owns `inputText`). It shows a
  slim menu above the input row while `inputText` starts with `/` and contains no space;
  printable keys still edit `inputText` and re-filter live; the menu is cleared on any
  non-slash input change or Esc. Suppressed on tiny terminals (`transcriptHeight() < 3`),
  replaced by an inline hint.
- **Tab** completes the longest common prefix of the matches; if exactly one match,
  completes the full command + a trailing space when it takes args. No Tab-cycling. ↑/↓
  move the highlighted suggestion; Enter accepts it (no-arg command executes; an
  arg/`OpensPicker` command + trailing space hands off to the existing picker).
- **Conditional Tab**: `tab`/`shift+tab` switch tabs ONLY when `inputText` is not
  slash-prefixed; `←/→` always switch tabs (muscle memory preserved). Hint line documents
  this.

### B. Kill one agent (runner + app + TUI)
- TUI: `ctrl+k` while on an agent tab whose state is running/steering opens a **modal**
  `confirmKillAgentID` sub-mode — the highest-priority interceptor in `updateMain` — that
  renders `kill <agent>? (y/N)` in warn style and blocks every other key; `y`/`enter`
  confirms → `opts.KillAgent(id)`, `n`/`esc` cancels. Killed agent shows a `KILLED` badge;
  the run continues.
- Runner `Handle` holds a mutex-guarded registry of per-attempt `context.CancelFunc` keyed
  by agentID (created in the agent-run path before the process starts, removed on exit).
  `KillAgent(agentID) KillResult` cancels ONLY that attempt's child ctx; if normal
  completion already deregistered, returns `Killed:false` and synthesizes no event
  ("first writer wins"). Emits `agent.killed` only when it wins; the attempt lands in a
  terminal failed/killed state (never promoted to success). **Never** calls run-wide
  cancel.

### C. Steer round-trip (the crux)
- A steer to an agent EXECUTES as a fresh single-agent attempt; its stdout is the reply,
  shown inline in that agent's tab.
- Flow / deterministic event ordering:
  1. TUI `submitInput` (steer on an agent tab) → `opts.SubmitSteer(req)`.
  2. App records `steer.requested` + `steer.delivered` (audit trail kept).
  3. Record-only paths (observational `/open` run, CLI `parley steer`, no live handle)
     return a record-only status and stop here.
  4. Else `Handle.RunSteerAttempt` accepts or rejects under the handle mutex (depth-1 queue
     per agent; a second steer while one is pending/running is **rejected** —
     "<agent> is already replying").
  5. When the attempt starts: under the handle **segment lock**, atomically compute the
     next segment id + append `run.segment_started` (`reason:"steer"`, `steer_id`) + create
     the per-steer dir + publish the segment id; then `steer.reply_started`; then
     register/spawn the process.
  6. On completion: exactly one terminal event — `steer.replied` (reply_path, stdout,
     duration_ms) or `steer.reply_failed` (error).
- Isolation (mandatory): steer attempts write ONLY under
  `runs/<runID>/agents/<agentID>/steers/<steerID>/{stdout.log,stderr.log,reply.md}` —
  never the round `stdout.log` (it's `os.Create`d). Steer attempt mode overrides ALL
  round-path defaults: no skip-if-`round-01/<id>.md`-exists, no protocol-artifact fallback,
  no artifact validation (validate only that a reply/stdout exists).
- `BuildSteerPrompt`: steer text + the idea `00-prompt.md` + the agent's OWN latest
  artifact (if any) + the tail of that agent's stdout (~50 lines / 4KB). Tells the agent:
  answer only this follow-up, do not edit protocol artifacts/other agents' files, print the
  reply to stdout and write `reply.md`.
- Reply surfacing: the TUI tails the returned `StdoutPath` into the target agent's
  transcript buffer behind a divider `── steer <steerID>: "<truncated query>" ──`, reply
  text in a dimmed/faint style; while running, a "<agent> is replying…" spinner in the
  status row / tab badge. When queued, a hint: `[steer queued — ctrl+k the active attempt
  to run it sooner]`. Reject a separate tab/pane.
- Errors: a synchronous `SubmitSteer`/`KillAgent` error → transient red `inputErr`; an
  async steer failure (`steer.reply_failed`) → a warning block appended to the agent tab.

### D. Seams (decoupling)
- `internal/tui` imports neither `internal/runner` nor `internal/app`. New capabilities
  are plain func fields on BOTH `LiveOptions` and `LaunchResult`:
  `SubmitSteer SteerFunc` and `KillAgent KillAgentFunc` (TUI-shaped request/result types
  defined in `internal/tui`). `activateRun` MUST copy both (else Home-launched runs lose
  them — a required regression test).
- Runner-native: extend the existing `Handle` (returned by `RunRoundOneAsync`) to own
  attempt state (store `opts`, root ctx, `mu`, `segmentMu`, `active map[string]…`,
  per-agent steer-busy state) and expose `KillAgent` + `RunSteerAttempt`. App binds these
  to the TUI seams. `RunRoundOne` (sync) can stay; the live async path returns the
  authoritative `*Handle`.
- Scope: kill/steer target the async round `Handle` path used by `parley tui`. Fix-up /
  implementation phases (driver-driven, headless; `RunFixup` doesn't use the shared
  agent-run helper) and deck-level steer fan-out are OUT of scope for v1 (noted as
  follow-ups). CLI `parley steer` stays record-only.

### E. Tests (no live terminal, fake agent command)
- Runner: `RunSteerAttempt` (printer fake → reply captured, per-steer paths used, events
  ordered) and `KillAgent` (sleeper fake → killed; a SECOND sleeper survives one kill;
  kill-vs-normal-completion race; duplicate-steer rejection) — include `-race`.
- TUI model tests (drive `Update`): Tab longest-common-prefix completion; suggest menu
  show/filter/clear; conditional-Tab (tab switches tabs when not slash-prefixed);
  confirm-kill modal dispatch (y/n, blocks other keys); `submitInput` steer via a fake
  `SubmitSteer`; `activateRun` copies `SubmitSteer`+`KillAgent` (Home-launched run).

### Non-goals
No long-lived interactive ACP chat; no multi-turn memory beyond the passed context; no
protocol change; no mouse; no fuzzy-finder lib; the normal round path is unchanged.

## Signoffs

(each participant appends its own block below — append-only, do not edit others')

### claude — ✅ ACCEPT
2026-06-06. As drafter I accept. This locks the round-02 resolution across all three
features: dedicated `suggest` autocomplete + conditional-Tab; per-attempt-registry kill
that never touches run-wide cancel; steer = isolated single-agent attempt (per-steer dir,
handle-locked segment, deterministic `steer.reply_*` events, depth-1 queue) surfaced
inline; injected `SubmitSteer`/`KillAgent` seams copied in `activateRun`. No blockers —
ready for FINAL + implementation.

### codex — ✅ ACCEPT
2026-06-06. I accept the consensus as drafted and confirm it reflects my round-02
position: dedicated `suggest` mode with conditional Tab, per-agent kill through the
`Handle` registry without run-wide cancellation, and steer replies as isolated
single-agent attempts with deterministic events and inline surfacing. The injected
`SubmitSteer`/`KillAgent` seams and `activateRun` copy requirement close the main
integration risk, so this is ready for FINAL + implementation.

### agy — ✅ ACCEPT
2026-06-06. I accept the consensus as drafted, which successfully incorporates the dedicated `suggest` mode, conditional-Tab, and inline steer feedback. The clean decoupling of `SubmitSteer` and `KillAgent` seams, alongside runner handle isolation, ensures safety and performance. This design is fully ready for FINAL review and implementation.

### hermes — ✅ ACCEPT
2026-06-06. I accept the final consensus. All three features are locked exactly as refined in round-02: non-modal suggest autocomplete with conditional-Tab, per-agent kill via Handle registry (never touching run-wide cancel), and isolated steer round-trips (per-steer dir, handle-locked segments, reply events, depth-1 queue, inline surfacing) plus the SubmitSteer/KillAgent seams copied in activateRun. Ready for FINAL + implementation.
