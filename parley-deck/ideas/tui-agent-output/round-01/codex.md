---
agent: codex
idea: tui-agent-output
round: 1
date: 2026-06-07
---

## Summary

The blank agent tab is primarily a tail-pipeline problem, not a runner-capture problem. The runner already captures `stdout.log` and `stderr.log`, and `agent.started` already projects both paths through `runstate.AgentState`. The TUI currently throws most of that away: `ensureBuffer` binds one `agentBuffer.path` to `StdoutPath`, `refreshBuffers` advances only that one file, `readAppendedLines` returns only newline-terminated content, and `splitLogLines` strips both ANSI and `\r`. That combination misses stderr progress, hides trailing partial output, and destroys the exact carriage-return rewrite behavior the owner wants.

My recommendation is to first make the TUI transcript model correct for two live log streams, carriage-return rewrites, and partial trailing lines. Then add per-agent streaming flags only where they can be parsed without breaking artifact capture. In concrete scope: merge stdout and stderr in the agent tab, dim/tag stderr, surface a live partial line, honor lone `\r` as "replace current line", keep `\r\n` as a normal newline, and keep stdout fallback untouched. Claude `stream-json` is worth adding only behind a per-agent adapter that writes readable text deltas to the log or a separate display stream; it should not be a blind replacement for `--output-format text` in the existing artifact-producing invocation.

## Proposed approach (concrete)

1. In-place streaming and CR handling

Do not keep `strings.ReplaceAll(data, "\r", "")`. Replace the current `splitLogLines(data []byte) []string` shape with a small stateful log ingester owned by `agentBuffer`.

The buffer needs to track completed scrollback plus one live line:

```go
type agentBuffer struct {
    streams map[logStream]*tailCursor
    lines   []transcriptLine
    partial *transcriptLine
    scroll  int
    follow  bool
    trunc   bool
}

type tailCursor struct {
    path    string
    info    os.FileInfo
    offset  int64
    loaded  bool
}

type transcriptLine struct {
    text   string
    stream logStream // stdout, stderr, event, steer
    time   time.Time // best-effort ingest/event time
    live   bool
}
```

If the implementation wants a smaller diff, keep `lines []string` initially and add `partial string`, but the two-stream styling pushes toward `transcriptLine` quickly.

The CR logic belongs after bounded file reading and before appending to scrollback. `completeLinesFrom` should stop being the only gate for display because the live trailing line must be visible before a newline arrives. A better helper is:

```go
func readAppendedChunk(path string, offset int64) (chunk []byte, newOffset int64, jumped bool)
```

It advances the offset to the bytes actually read, not merely to the last newline. This prevents rereading the same partial bytes on every tick and makes the live partial stable. The ingester, not the offset code, decides which text is completed.

For the ingester:

- Treat `\r\n` as `\n`.
- Treat lone `\r` as "replace the current live line with what follows".
- When `\n` is seen, commit the current live line into bounded scrollback and clear the partial.
- When bytes end without `\n`, update `partial` and render it as the last line with `live: true`.
- Strip ANSI after CR interpretation, or strip recognized ANSI while preserving `\r`/`\n`. Stripping first is acceptable only if the ANSI stripper cannot consume control bytes around `\r`.
- For a segment like `one\rtwo\rthree`, the visible live line should be `three`, not `onetwothree`.
- For `one\rtwo\n`, commit `two`.
- For `one\r\ntwo`, commit `one`, then live/commit `two` according to following bytes.

This should not mutate already completed lines. A lone `\r` rewrites only `partial` or the not-yet-newline-committed current line. Once a newline commits a line to `lines`, later CRs apply to a new live line. That is the main corruption boundary.

`renderTranscript` should render `lines` plus `partial` as the logical scrollback length. If `follow` is on, the partial line is visible at the bottom and rewrites in place each refresh. If the user scrolls up, `follow` is off and the partial may continue changing below the viewport without yanking scroll position. `bufferBottom` must count the partial line when present.

ANSI should stay conservative. Keep color/styling controlled by the TUI, not arbitrary CLI escape sequences. Strip ANSI color and cursor controls for now, but honor CR semantics. Raw cursor movement sequences are too easy to misapply in a scrollback model and can corrupt old lines.

2. stderr visibility and two-stream tailing

Use one `agentBuffer` per agent, but give it two cursors:

```go
type agentBuffer struct {
    stdout tailCursor
    stderr tailCursor
    lines  []transcriptLine
    partial map[logStream]string
    ...
}
```

`ensureBuffer` should populate `stdout.path = a.StdoutPath` and `stderr.path = a.StderrPath`. `refreshBuffers` should advance both cursors for each loaded agent. Reuse the current rotation/truncation check per cursor.

Ordering cannot be perfect from two independent files because stdout and stderr writes are captured separately and the files do not carry per-line timestamps. The clean, honest rule is ingest-order ordering by TUI refresh: on each tick, read appended stdout and stderr chunks, assign each produced line an ingest sequence/time, and append in a deterministic order for same-tick data. I recommend stdout first, then stderr, with stderr dimmed and prefixed as `[err]` or styled in a dedicated column. This avoids pretending we can recover kernel-level interleaving after the streams have been split. If precise ordering becomes required, the runner must tee both streams into a single combined log at capture time; that is a larger runner change and not necessary for the owner problem.

Bounded scrollback remains shared per agent. After appending from either stream, run `capFocusLines` or its `[]transcriptLine` equivalent on completed lines. Partial lines should have their own byte cap so a newline-less verbose process cannot grow memory; if a single partial exceeds the byte budget, keep the tail with a truncation marker.

For display, merge by default. A toggle/split can be a later UX refinement, but the immediate fix should show stderr because that is where many CLIs put progress. Dim or tag stderr so errors/progress are distinguishable without hiding them.

3. Conversation/discussion view

The current `renderTranscript` completely replaces the round transcript with `renderSteerReply` whenever `m.steerReplies[agentID]` exists. That is the opposite of a scrollable discussion. Steers should be appended as transcript entries:

- `you: <question>` as an event line when the steer is submitted.
- `agent: <live partial reply>` from the steer attempt stdout, using the same CR/partial ingester.
- a final status line such as `reply complete` or `reply failed`.

This can reuse `steerReplies`, but the rendering should not swap out the base transcript. Treat each steer reply as another stream segment in the agent tab. A minimal version can render base `b.lines`, then a divider/query, then the steer buffer lines/partial, all inside the same scroll math. The better version is to add steer lines into the same `agentBuffer` transcript with `stream: steer`.

The target exchange should look like:

```text
[run] reading prompt
[err] resolving model...
you: What is the blocker?
agent: checking state...
agent: validating state transitions
agent: The blocker is the missing stderr tail.
```

During streaming, the `agent:` line rewrites in place. After newline/final completion, it becomes an ordinary completed line in the scrollback.

4. Streaming and verbose agent flags

Do not globally change all headless flags as part of the first TUI fix. The safe baseline is:

- `codex exec ... -`: likely already streams useful stdout; keep it unchanged and let the improved TUI show its live partial and CR rewrites.
- `claude -p --output-format text`: consider `--output-format stream-json`, but only behind a Claude-specific parser that extracts text deltas into readable display text. Blindly writing JSON events into `stdout.log` would make the TUI noisy and would break or confuse stdout artifact fallback because stdout would no longer be a plain Markdown artifact candidate.
- `agy --print`: keep unchanged unless there is a documented streaming mode. Its stdout fallback is important and already tested; changing the shape of stdout risks losing recovered artifacts.
- `hermes --oneshot`: keep unchanged unless a known streaming mode is discovered and validated.

For Claude, the cleanest scope is either:

- keep the artifact-producing invocation as text and rely on stderr plus artifact/status for now; or
- add a per-agent output adapter in the runner: raw stdout remains captured if needed, parsed text deltas are written to a display log, and artifact validation reads the final artifact file or a known final text field, not raw stream JSON.

Given the owner wants live output, Claude streaming is worth doing, but after the TUI tail model is fixed and with tests around parser output, stdout fallback, and validation. It should be opt-in per `Spec` capability, not a blanket `HeadlessArgs` edit that changes the contract for every phase.

5. Always-on status and artifact view

The tab should never render "no output yet ... waiting for stdout". Use projected state from `AgentState`:

- pending: `claude pending - no process started yet`
- running: `claude working... 0:42` plus liveness `STALE` if applicable
- finished: `claude finished in 2m12s - wrote round-01/claude.md`
- failed/killed/skipped: show the state and reason/error

The source is already present: `AgentState.State`, `StartedAt`, `Duration`, `LatestEvent`, `Error`, `Reason`, `ArtifactPath`, `StdoutPath`, and `StderrPath`.

Artifact preview should come from `AgentState.ArtifactPath`, not from stdout. The path in events appears to be absolute in the runner, so the renderer can tail or preview it directly if present. I would not mix full artifact preview into the first implementation of the transcript tailer; instead show a short status line like `artifact: round-01/codex.md` and add a key/toggle later for preview. If preview is added now, it should reuse bounded tail helpers and be visually separate from live process output so the user can tell "what the agent is saying now" from "what file it produced".

6. Performance, bounds, testability, and safety

The existing bounded tail strategy is the right base: max bytes per read, max retained lines, and lazy buffers for active/visited agents. Extending from one file to two files doubles cheap `stat`/bounded reads for loaded buffers only. That is acceptable for normal Parley rosters. For many agents, keep the current lazy behavior: load active and visited tabs, not every possible agent artifact. If a status indicator needs freshness for non-active agents, use events/liveness rather than tailing all logs eagerly.

Headless tests should cover the model without a real TUI:

- `split/ingest` tests for `a\rb`, `a\rb\n`, `a\r\nb\n`, multiple CRs, ANSI around CR, and newline-less partial updates.
- Offset tests proving partial bytes are not duplicated and oversized partials stay bounded.
- Two-stream tests with fake stdout/stderr files appended across ticks, verifying both cursors advance independently and stderr styling metadata survives.
- Rotation/truncation tests per stream.
- Render tests where follow is on and the partial line appears at the bottom, then follow is off and refresh does not pull the user back down.
- Steer tests proving a steer reply is appended/woven into the transcript rather than replacing it.
- Runner regression tests preserving stdout fallback for print-only agents and validation rejection for malformed stdout candidates.

## Concerns / open questions

The biggest open ordering question is whether "without ordering corruption" means "do not scramble within each stream" or "recover exact stdout/stderr interleaving". The former is straightforward with two cursors. The latter is impossible once the streams are stored as separate files without timestamps or a combined capture log. If exact interleaving is required, the runner should write a third `combined.log` with stream tags at capture time, while retaining the existing stdout/stderr files for diagnostics and artifact fallback.

Claude `stream-json` needs a real sample before committing the parser. The adapter should extract only assistant text deltas and ignore metadata/progress JSON for the transcript, unless there is a clear UX reason to show selected status events. A malformed or changed stream format must fail soft: keep raw stderr/status visible and do not break the agent run.

It is also unclear how much artifact preview belongs in the initial fix. I think status plus path is necessary; a full preview/toggle is useful but secondary to making the live discussion nonblank and correct.

## Risks

The main correctness traps are:

- `\r` versus `\r\n`: treating all `\r` as rewrite will corrupt Windows newlines. Normalize `\r\n` to `\n` first or handle it as one newline token in the scanner.
- Partial-line offset accounting: if offsets advance only to completed newlines, the same partial bytes are reread every tick; if offsets advance without preserving partial state, trailing output is lost. The fix needs both advancing byte offsets and per-stream partial state.
- Completed-line corruption: CR should rewrite only the current live line. It must not edit already committed scrollback lines.
- Two-stream ordering: separate stdout/stderr files cannot prove exact cross-stream order. Do not invent false precision; use deterministic ingest order or add a runner combined log.
- Follow/scroll interaction: when `follow` is on, new completed lines and live partial updates keep the viewport pinned to bottom. When the user scrolls, `follow` must turn off and refreshes must not snap back merely because a partial changed.
- Artifact fallback: changing stdout to JSON or progress chatter can break the existing recovery path. Preserve current stdout semantics for print-only agents or isolate display streaming from artifact capture.
- ANSI handling: arbitrary cursor movement does not map safely into a scrollback transcript. Strip ANSI for now, but preserve CR semantics before stripping destroys the rewrite signal.
