---
agent: claude
idea: tui-agent-output
round: 1
date: 2026-06-07
---

## Summary

The tab is blank for two compounding reasons: (a) the agents are configured for
final-text one-shot output AND any progress goes to **stderr**, which the TUI never shows;
(b) even when a CLI streams a self-rewriting progress line via `\r`, `splitLogLines`
**strips `\r`**, so the Codex-CLI "one line that rewrites itself" is destroyed. Fix in
three moves, headline first: honor `\r` for in-place line updates + surface the live
partial line; merge stderr into the scrollback (dimmed); and turn the steer round-trip
into a woven, scrollable conversation. Streaming agent flags (claude `stream-json`) are a
worthwhile but secondary multiplier.

## Proposed approach

### 1. In-place streaming (the headline) — honor `\r`, show the live partial line
Today `readAppendedLines` advances the offset only past the last `\n` and DROPS the
trailing partial line; `splitLogLines` does `ReplaceAll("\r","")`. Two changes:
- **CR-collapse each line**: when turning bytes into a display line, keep only the text
  after the LAST `\r` in that line (the final rewrite). A progress line
  `"… 10%\r… 50%\r… 90%"` renders as `"… 90%"`. Applied to both committed lines and the
  live partial.
- **Surface the live partial line**: keep a per-buffer `partial string` = the CR-collapsed
  text after the last `\n` (the not-yet-newlined current line). Render it as the last line
  of the transcript so a streaming line UPDATES in place each tick. When a `\n` finally
  arrives it commits and the partial clears. This is exactly "vidim jeden riadok a potom
  sa to prepíše."
- **ANSI**: keep stripping ANSI control for v1 (cursor moves would corrupt the line model);
  optionally preserve SGR colour later. CR-rewrite is the ask; colour is a nice-to-have.

### 2. Show stderr too (merged, dimmed)
Tail BOTH `stdout.log` and `stderr.log` into the agent's scrollback. Keep two underlying
offset-cursors but ONE rendered history: each tick, append new stdout lines (plain) then
new stderr lines (dimmed, `[err]`-tagged). Per-tick ordering is good enough (true
cross-file interleave needs timestamps we don't have). Reuse the bounded
`loadFocusTail`/`capFocusLines` for both so growth stays bounded. This alone makes silent
agents (agy/hermes, which narrate on stderr) visible immediately.

### 3. Weave the steer conversation into the scrollback (not a replacing panel)
1.18's `steerReplies` REPLACES the transcript while a reply streams. The owner wants a
scrollable DISCUSSION. Change it so a steer renders as conversation entries appended into
the agent's history: `› you: <question>` then the reply streamed below (CR-in-place), part
of the same scrollback you can scroll. The agent's own working output and your steers/
replies live in one chronological, scrollable view (Codex-CLI feel).

### 4. Always-on status header (never blank)
A one-line header from the projection above the transcript: `agy · working 0:42` /
`wrote round-01/agy.md` / `finished 2m12s` / `failed: …`. Replaces the bare "no output
yet" so the tab is alive even before the first byte.

### 5. Show the produced artifact
When the agent's artifact exists (path from the projection/events), offer it: a key (e.g.
`a`) toggles "live output ⇄ artifact", or show a short preview at the bottom. The artifact
is the concrete "co napísal."

### 6. Streaming agent flags (secondary multiplier)
- claude: `--output-format stream-json` streams token/JSON deltas → needs a tiny parser to
  extract text deltas into stdout-display; gives true live streaming. Worth it but adds a
  parser; confirm it doesn't affect the file-artifact path (claude writes the file via its
  tools, independent of stdout format).
- codex exec already streams text to stdout → works once `\r`/stderr handling lands.
- agy/hermes: one-shot, no streaming flag → rely on stderr + status + artifact.
Keep per-agent; do NOT change a flag that breaks the invocation or the artifact-capture
fallback.

## Concerns / open questions
- The live-partial-line surfacing interacts with `follow`/scroll: when following, the
  partial is the bottom line; when scrolled up, it should not yank the view. Reuse the
  existing follow logic.
- Merging stderr could be noisy for chatty agents; a toggle to hide stderr may be wanted.
- claude stream-json parsing scope — keep it minimal (text deltas only) or defer to a
  follow-up and ship CR+stderr+conversation first? I lean: ship 1-4 solidly; stream-json
  as a fast-follow if it risks the invocation.
- Two-file tail bound: cap each stream; the merged history honors the existing line/byte cap.

## Risks
- Mis-handling `\r` could corrupt completed multi-line output (e.g. a real `\r\n`). Treat
  `\r\n` as a normal newline; only a lone `\r` (not followed by `\n`) triggers in-place
  rewrite. Test with `\r`-progress, `\r\n`, and plain `\n`.
- Changing claude/codex flags could break a working headless run or the artifact-capture
  fallback — gate flag changes behind verification; keep agy/hermes untouched.
- Performance: tailing two files + artifact each tick × many agents — reuse bounded reads;
  only refresh the active + visited buffers (already the case).
- Steer-conversation rewrite must not regress the 1.18 reply behaviour or the esc-dismiss.
