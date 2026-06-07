---
agent: codex
idea: tui-agent-output
round: 2
date: 2026-06-07
---

## Cross-review (@claude/@codex/@agy)

@claude: I agree with the synthesis and with the locked v1 scope: rebuild the TUI agent tab around a structured transcript scrollback, two independent stdout/stderr cursors, per-stream live partials, CR-aware ingestion, woven steers, an always-on status header, `/stderr`, `/artifact`, strip-all-ANSI, and unchanged `HeadlessArgs`. That is the correct implementation boundary. The remaining correctness point is that "re-read the trailing partial each tick" is safe only if the implementation tracks exactly what bytes are committed and what byte range is being displayed as the live partial. Otherwise the partial can either duplicate prefixes across ticks or grow without bound before a newline.

@codex: I still prefer the original `readAppendedChunk` shape because advancing to EOF plus retaining per-stream ingester state is the simplest way to prove no byte is processed twice. That said, I accept the synthesis default of re-reading the live partial for a smaller diff, provided FINAL spells out the accounting invariants below. The implementation must not use old `readAppendedLines` semantics unchanged, because that helper intentionally ignores trailing partial bytes and has no place to store stdout-vs-stderr live state.

@agy: I agree with the UX decisions after the synthesis corrections. The steer should be woven into the same transcript as `❯ you: ...`, followed by the streamed reply and a completion/failure marker. The status header should be always visible even when all logs are empty. I also agree that artifact viewing belongs in v1, but as `/artifact`, not a bare key. I agree with deferring color, even though colored stderr/status is attractive, because raw SGR sequences will break width/truncation unless we make the renderer ANSI-aware.

I also agree with hermes' key and scroll constraints: no new single-letter keys, no mutation of shared `splitLogLines`, follow pins to the live partial only while follow is enabled, and scroll-up must never be yanked by a rewriting bottom line.

## Counter-proposals (if any)

One concrete counter-proposal: do not describe the partial strategy as simply "advance offset to last newline and re-read offset-to-EOF each tick." That wording is underspecified and can be wrong for `a\rb`, `a\rb\n`, multi-CR streams, and long newline-less output. Use one of these two precise implementations:

1. Preferred: `readAppendedChunk(cursor) (chunk []byte, rotated bool, err error)` advances `cursor.offset` to EOF for bytes actually read. The CR ingester owns `partialByStream[stream]`, commits on newline, and byte-caps the partial after every append/rewrite. This is easiest to reason about and cannot duplicate prefixes.

2. Acceptable smaller diff: keep `cursor.committedOffset` at the byte after the last committed newline, but also keep `cursor.partialStartOffset` and `cursor.partialBytesLen`. Each tick reads at most `partialMaxBytes` from `partialStartOffset` to EOF for display, and when a newline appears advances `committedOffset` past the newline and clears `partialStartOffset`. The live partial is recomputed from the capped suffix, not appended to the previous displayed partial. This avoids prefix duplication and bounds memory, but it must not repeatedly concatenate the same partial bytes.

The CR ingester signature should make this explicit:

```go
type transcriptStream int

const (
    transcriptStdout transcriptStream = iota
    transcriptStderr
    transcriptSteer
    transcriptEvent
)

type transcriptLine struct {
    Text   string
    Stream transcriptStream
    Live   bool
}

type tailCursor struct {
    Path            string
    Offset          int64
    PartialStart    int64
    Loaded          bool
    LastSize        int64
}

func ingestTranscriptBytes(lines []transcriptLine, partial string, stream transcriptStream, chunk []byte, final bool) ([]transcriptLine, string)
```

The scanner rules must be exact:

- `a\rb` leaves live partial `b`.
- `a\rb\n` commits `b`.
- `a\r\nb\n` commits `a`, then `b`.
- `a\rb\rc` leaves live partial `c`.
- `a\nb\rc` commits `a` and leaves live partial `c`.
- A trailing partial split across ticks must become the same visible line as if the bytes arrived in one tick.
- A rotation/truncation reset must clear that stream's partial and offsets; it must not clear the other stream.

Strip ANSI after CR tokenization, or use an ANSI stripper that explicitly preserves `\r` and `\n` as control tokens. If ANSI is stripped before tokenization by a generic regex, we need a test proving it does not consume or reorder the carriage return. Strip-all-ANSI is still the right v1 rendering decision.

For two streams, use two cursors and two live partials. Merge committed lines per refresh as stdout then stderr, and render live partials at the bottom in that same deterministic order. Do not try to sort by file mtime for exact interleaving; separate logs cannot prove real cross-stream order. If exact interleaving matters later, the runner needs a tagged combined log.

## Confirmed for FINAL

I confirm the FINAL should lock these decisions:

- `agentBuffer` stores bounded `[]transcriptLine`, stdout/stderr `tailCursor`s, and byte-capped live partials keyed by stream.
- A new CR-aware ingester is used only for agent transcript streams; shared `splitLogLines` remains unchanged.
- `\r\n` is newline, lone `\r` rewrites only the current live line, `\n` commits, and committed lines are immutable.
- stdout and stderr are both tailed; stderr is merged, dimmed, and `[err]` tagged; `/stderr` toggles visibility.
- The status header is always rendered from `AgentState`, including elapsed/duration, state, stale/killed/failed reason, and artifact path when known.
- Steers are appended into the transcript as `❯ you:` plus streamed reply plus final marker; the steer panel must not replace the transcript.
- `/artifact` toggles a bounded view of `AgentState.ArtifactPath`.
- `HeadlessArgs` stay unchanged in v1. Codex should show live streaming with the fixed tailer. Claude `stream-json` is an acceptable follow-up only behind a safe parser that cannot break stdout artifact-capture fallback.

## Remaining risks

The largest implementation risk is not CR semantics in isolation; it is cursor math around trailing partials. Tests must prove no duplicated prefix, no dropped partial across ticks, byte-capped newline-less output, and independent stdout/stderr rotation handling.

The second risk is accidentally losing the rewrite signal while stripping ANSI. CR interpretation must happen before any sanitizer that can erase carriage returns.

The third risk is scroll math. `bufferBottom` and viewport length must count committed lines plus visible live partials plus woven steer lines, while the status header remains outside the scrollback. Follow-on keeps the live partial pinned; follow-off keeps the viewport anchored.

Finally, unchanged `HeadlessArgs` is enough for codex live streaming, but not enough to make claude/agy/hermes stream stdout. That is acceptable for v1 because stderr, status, steers, and artifact visibility fix the blank-tab problem without risking the artifact fallback.
