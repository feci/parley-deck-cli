---
agent: hermes
idea: tui-agent-output
round: 1
date: 2026-06-07
---

## Summary
The agent tab must deliver a Codex-CLI-like live scrollable transcript: a one-line streaming answer that rewrites itself in place via \r, merged stdout+stderr tail, steers woven in, never blank, bounded scrollback. Current splitLogLines destroys \r semantics by stripping them; agentBuffer already has scroll/follow but only tails stdout. Interaction must preserve all existing scroll/picker/suggest/confirm-kill/steer keys while adding at most one non-colliding toggle.

## Proposed approach
Honor \r as "rewrite current line" only when it is a lone \r (not followed by \n). Rule: split on \r\n and \n normally; on a bare \r replace the last line in the accumulating buffer instead of appending. This lives in a new splitLogLinesCR variant called from readAppendedLines (after the existing offset read, before bounded cap). Keep ANSI colours, strip only cursor/position escapes (or none initially). Merge: tail both StdoutPath and StderrPath each tick into the same agentBuffer.lines (stderr lines tagged [err] and dimmed on render); interleave by file modtime or simple append order per tick (acceptable for progress narration). Follow pins viewport to live partial line; scroll-up pauses follow so a rewriting bottom line never yanks the view. No new keys needed: reuse /follow and existing scroll keys; stderr toggle can ride on a future 'e' only if it does not collide (current keys have no 'e').

Key-table delta (none required):
- No additions; /follow already exists and covers live pinning.

Conversation weave reuses steerReplies render path exactly as-is; streamed lines simply continue below the steer divider.

## Concerns / open questions
- Partial-line offset accounting on bare \r: must track "current line start offset" so re-read on next tick does not duplicate the prefix.
- Ordering when stdout and stderr interleave rapidly: modtime sort per chunk or accept occasional stderr-before-stdout.
- Artifact preview: where does its path come from in the projection (already in round-NN dir) without new seams?

## Risks
- Changing splitLogLines affects every tail (including non-agent); must be a parallel helper.
- Two-file tail per tick for 8 agents stays cheap only because bounded loadFocusTail is reused.
- Must not touch 1.18 steer/esc-dismiss paths.