---
idea: tui-agent-output
phase: final
status: final
drafter: claude
implementer: claude
date: 2026-06-07
participants: [claude, codex, agy, hermes]
supersedes: consensus.md
---

# FINAL — live, scrollable, Codex-CLI-style agent transcript

Ratified by claude, codex, agy, hermes (all ✅ ACCEPT in consensus.md). Build spec for
Phase 5; scope is `internal/tui/live.go` (+ `live_test.go`). No agent-flag changes, no new
deps; `internal/tui` stays decoupled (projection + existing seams only).

## 1. Transcript model (replace the stdout-only `agentBuffer`)
```go
type transcriptStream int
const ( transcriptStdout transcriptStream = iota; transcriptStderr; transcriptSteer; transcriptEvent )
type transcriptLine struct { Text string; Stream transcriptStream; Live bool }
type tailCursor struct { Path string; Offset int64; Info os.FileInfo; Loaded bool }
```
`agentBuffer`: `stdout, stderr tailCursor`; `lines []transcriptLine` (bounded by the
existing line/byte caps); `partial map[transcriptStream]string` (each byte-capped);
`scroll int`, `follow bool`, `trunc bool`; plus `hideStderr bool`, `showArtifact bool`.

## 2. CR-aware ingestion (headline) — NEW helpers, leave shared `splitLogLines` untouched
- `readAppendedChunk(c *tailCursor) (chunk []byte, rotated bool)`: stat; on shrink/rotate
  (`size < Offset` or `!SameFile`) signal rotated + reload from `max(0, size-maxFocusBytes)`;
  else read `Offset..EOF` bounded to `maxFocusBytes` (jump-ahead on a huge burst), advance
  `Offset` to the bytes actually read.
- `ingestTranscriptBytes(lines []transcriptLine, partial string, stream transcriptStream,
  chunk []byte) (out []transcriptLine, newPartial string)`: scan rune/byte-wise:
  `\r\n`→newline; `\n`→commit `partial` as `transcriptLine{cleanANSI(partial),stream}` (drop
  if empty after clean), clear; lone `\r`→`partial = ""` (rewrite — keep only text after it);
  else append to `partial`; cap `partial` to `partialMaxBytes` (keep tail). Committed lines
  are never mutated. `cleanANSI` strips ANSI AFTER CR handling.
- Scanner cases (tests): `a\rb`→partial `b`; `a\rb\n`→commit `b`; `a\r\nb\n`→commit `a`,`b`;
  `a\rb\rc`→partial `c`; `a\nb\rc`→commit `a`, partial `c`; split-across-ticks consistent;
  rotation clears ONLY that stream.

## 3. refreshBuffers / ensureBuffer (two streams)
`ensureBuffer`: `b.stdout.Path = a.StdoutPath`, `b.stderr.Path = a.StderrPath`; initial load
via `readAppendedChunk` + ingest for each. `refreshBuffers`: for each loaded buffer, for
each cursor (stdout then stderr): `chunk, rotated := readAppendedChunk(&cur)`; if rotated
reset that stream's `partial`; `b.lines, b.partial[stream] = ingestTranscriptBytes(b.lines,
b.partial[stream], stream, chunk)`; cap `b.lines`; if `follow` set `scroll = bufferBottom`.

## 4. renderTranscript (status header + merged scroll + partial + artifact)
- ALWAYS render a status header (outside the scrollback) from `AgentState`:
  `● <id> working… M:SS` (warn, `+ STALE` if `agentLiveness=="stale"`) / `✓ <id> wrote
  <rel-artifact> (dur)` (ok) / `✗ <id> failed: <err> (dur)` (danger) / `◌ <id> killed`
  (muted) / pending. Replaces "no output yet".
- Body = the visible transcript: committed `lines` filtered (`hideStderr` drops stderr;
  `/artifact` mode shows the artifact instead) + the live `partial`s appended at the bottom
  (stdout then stderr), each styled by `Stream` (stderr/`transcriptEvent`=dimmed `[err]`/
  tag, `transcriptSteer` `❯`/reply style). Scroll/follow math counts committed + partial +
  woven lines; `bufferBottom` includes the partials.
- `/artifact` mode: header `[Viewing Artifact: <rel-path>]`, body = bounded tail of
  `AgentState.ArtifactPath` (reuse `loadFocusTail`); missing/empty → `[Artifact not yet
  written]`. No fd leak.

## 5. Conversation weave (remove the replace-panel)
Drop the `steerReplies`-replaces-transcript behavior. On steer submit, append a
`transcriptSteer` line `❯ you: <text>` into the active agent's `b.lines`. Tail the steer
attempt's stdout (the `SubmitSteer` `StdoutPath`) as a `transcriptSteer` stream into the
SAME buffer (a third cursor or fold into the steer reply state), CR/partial-handled; on the
`steer.replied`/`reply_failed` event append `[reply complete — Ns]` / `[reply failed]`. esc
no longer "dismisses a panel" (it's just scrollback now). Keep the steer input row + cyan
prefix + record-only messaging from 1.18.

## 6. Keys / commands
Add `/stderr` (toggle `hideStderr`) and `/artifact` (toggle `showArtifact`) to
`commandSpecs` (autocomplete + `/help`) and `runCommand`. No new single-letter keys; scroll
keys unchanged; `/follow` still re-pins.

## 7. Tests (model-driven, no terminal)
The §2 scanner cases; no-prefix-duplication + bounded partial across two ingest ticks; two
cursors advance independently; per-stream rotation resets only that stream; stderr `[err]`
tag + `/stderr` hides it; follow pins to a rewriting partial + scroll-up doesn't yank;
steer woven (a `❯ you:` line + reply appended, transcript NOT replaced); status header per
state (working/finished/failed/killed); `/artifact` toggles + missing-file graceful. Keep
the runner stdout-fallback + validation tests green (unchanged).

## Acceptance
`go build ./... && go vet ./... && go test ./...` green; a manual-smoke note. Then Phases
6-8 review to zero agreed fixes. (Out of v1: ANSI colour passthrough, exact cross-stream
interleave, claude `stream-json`.)
