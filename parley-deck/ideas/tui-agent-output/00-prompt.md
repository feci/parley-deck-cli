---
idea: tui-agent-output
author: user
created: 2026-06-07
participants: [claude, codex, agy, hermes]
roles:
  claude: facilitator + TUI transcript rendering (stdout+stderr merge, artifact view, status)
  codex: runner/capture — two-stream tailing, streaming/verbose agent flags, Go correctness, tests
  agy: UX — what the combined transcript looks like, stderr styling, artifact preview, "working…" empty-state
  hermes: keymap/interaction + merged-tail correctness (ordering, scroll, no corruption)
transport: local-dir
cross_review_rounds: 1
status: final
---

## Problem / idea (owner's words)

When `parley tui` runs the agents headless, **nothing shows in the agent tabs** — the
owner sees no output while an agent works. They want: if you run agents headless, **show
what they write**; ideally keep running them headless but **read and display their
output** in the TUI.

**Refinement (owner):** ideally it should feel like the **Codex CLI agent** — a
**scrollable discussion**, where you ask a question and **see a one-line answer that then
rewrites itself in place** as it streams. So the target is a LIVE, streaming, scrollable
transcript whose working line updates in place (not a blank tab that maybe fills at the
end). The steer round-trip (1.18, type-to-agent → reply) is part of this conversation.

## Current state (VERIFIED against the code — design against these facts)

- **The TUI agent transcript tails ONLY stdout.** `renderTranscript`
  (`internal/tui/live.go`) shows `m.buffers[agentID]`, whose `path = a.StdoutPath`
  (`ensureBuffer`). When empty it literally prints
  `"no output yet from <agent> (waiting for the agent to write stdout)"` — exactly what
  the owner sees. **stderr is NEVER shown**, although it IS captured: the runner writes
  `<agentDir>/stderr.log` and `agent.started` records both `stdout` + `stderr` paths
  (`runstate.AgentState.StderrPath` exists).
- **The agents are configured for FINAL-text one-shot output**, so stdout stays empty for
  the whole run and only fills (if at all) at the end (`internal/agents/discover.go`):
  - claude: `-p --output-format text …` (final text only; `--output-format stream-json`
    is available per the Telemetry note).
  - agy: `--print {prompt}` (one-shot; prints the final answer).
  - hermes: `--oneshot {prompt} …` (one-shot).
  - codex: `exec … -` over stdin — codex exec DOES stream text to stdout, so codex is the
    closest to "live", but the others are silent until done.
  So during a (multi-minute) round the agent tab is blank; any progress an agent emits
  goes to **stderr**, which the TUI doesn't show.
- **The artifact the agent produces** (its `round-NN/<agent>.md`, or IMPLEMENTATION.md,
  etc.) is written to `opts.Idea.Path/...` and is NOT shown in the agent tab — yet that
  is literally "co pisu" (what they wrote). There IS a stdout-capture fallback for
  print-only agents (agy) that recovers the artifact from stdout, so agy's stdout.log
  does end up holding its printed artifact after it finishes.
- The TUI's `agentBuffer` tails ONE file via offset-incremental reads
  (`loadFocusTail`/`readAppendedLines`, bounded scrollback). Events (`agent.started`,
  ACP `agent.acp.*` updates, `agent.finished`) are in `events.jsonl` and projected into
  `runstate` (state, duration, latest event) but are not woven into the transcript.
- **In-place updates are currently DESTROYED.** `splitLogLines` does
  `stripANSI(strings.ReplaceAll(data, "\r", ""))` — it strips carriage returns AND ANSI.
  So a CLI that uses `\r` to rewrite a progress line (the Codex-CLI "one line that updates
  itself") or ANSI cursor control has that collapsed/flattened into noise or lost. To get
  the owner's "potom sa to prepíše" behaviour, the renderer must HONOR `\r` (treat it as
  "rewrite the current line") so a streaming line updates in place. The transcript IS
  already scrollable (shift+↑/↓, PgUp/PgDn, ctrl+u/d, Home/End, follow/`/follow`).

## Proposed direction (a STARTING proposal — challenge it in round-01)

Make the agent tab a LIVE, scrollable, Codex-CLI-style transcript:
1. **Stream live output with in-place updates.** HONOR `\r` so a streaming progress line
   rewrites itself in place ("potom sa to prepíše") instead of being stripped. Decide how
   much ANSI to keep (cursor moves / colours) vs strip; at minimum CR-rewrite the current
   line. This is the headline change.
2. **Show stderr too.** Merge the agent's stdout AND stderr (stderr dimmed/`[err]`-tagged)
   so CLI progress/narration is visible — many agents narrate on stderr. Tail BOTH files
   into the agent's scrollback without ordering corruption or unbounded growth.
3. **Weave the steer conversation in.** The discussion = the agent's streamed working
   output + the steers the owner sent + the replies, scrollable as one history (reuse the
   1.18 steer round-trip / events).
4. **Always-on status from events.** A header line so the tab is never blank:
   `agy working… 0:42` / `wrote round-01/agy.md` / `finished in 2m12s` / `failed: …`,
   replacing the bare "no output yet".
5. **Show the produced artifact** when written (preview / toggle) — the concrete thing
   they wrote.
Secondary but important for "live streaming": **enable streaming output** where the CLI
supports it (claude `--output-format stream-json` parsed to readable text; codex exec
already streams) so stdout fills DURING the run — per-agent config, must not break the
existing invocations.

## Round-01 focus questions (answer independently)

1. **In-place streaming (the headline).** How to honor `\r` so a streaming line rewrites
   itself in the scrollback (Codex-CLI feel)? How much ANSI to keep vs strip (colour vs
   cursor control)? Where in the tail pipeline (`splitLogLines`/`readAppendedLines`/render)
   does CR-handling go, keeping bounded scrollback and not corrupting completed lines?
2. **stderr visibility + two-stream tailing.** Merge stdout+stderr (dimmed stderr) vs a
   toggle/split? How to tail TWO files with the offset-incremental `agentBuffer` without
   ordering corruption or unbounded growth?
3. **Conversation/discussion view.** Weave steers + replies + streamed output into one
   scrollable history (reuse steerReplies/events)? What does a "you asked → one-line
   streaming answer that rewrites → final" exchange look like in the tab?
4. **Streaming/verbose agent flags.** Worth switching claude to `stream-json` (parsed to
   text) + confirming codex streams, vs only showing stderr+artifact? Which agents support
   a live-progress stdout mode, and how risky is changing the headless flags (could break
   a working invocation, the artifact-capture fallback, or validation)? Keep it per-agent.
5. **Always-on status + artifact view.** What to show from the projection so the tab is
   never blank; show the produced artifact (preview vs full / toggle) and where its path
   comes from.
6. **Performance, bounds, testability, safety.** Two-file + artifact tailing each tick for
   many agents — keep it cheap (reuse bounded `loadFocusTail`/cap), short terminals. What
   is testable headlessly (model-driven render over fake stdout/stderr/artifact + a CR
   stream). Do NOT break the round/steer path, `--no-tui`, the artifact-capture fallback,
   or durable-kill/liveness.

## Constraints (non-negotiable)
- Reuse the existing `agentBuffer` tailing machinery and the already-captured
  `stderr.log`; bounded scrollback (no unbounded reads). No new third-party deps.
- `internal/tui` reaches data through the projection/seams it already has (StdoutPath,
  StderrPath, events) — no import of runner/app.
- Do NOT break the normal round/steer run path, `--no-tui`, the artifact-capture fallback,
  or the durable-kill/liveness work. Changing an agent's headless flags must keep that
  agent's invocation working (verify against the discover config + validation).
- Testable headlessly (fake stdout/stderr/artifact files driven through the model).

## Non-goals
- No live interactive agent chat (that's steering). No re-architecting how agents run. No
  parsing every vendor's bespoke streaming format beyond what's cheap (claude stream-json
  is the one worth it). No mouse.

## Deliverables (in order)
1. FINAL.md: the agent-tab content model (stdout+stderr merge + styling, artifact view,
   always-on status), the two-stream tailing approach, any per-agent streaming-flag
   changes, performance bounds.
2. Implement in `internal/{tui,runner,agents}` (+ tests), build/vet/test green; a
   manual-smoke note. Then Phases 6-8 review.
