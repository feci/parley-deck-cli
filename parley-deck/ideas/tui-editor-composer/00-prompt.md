---
idea: tui-editor-composer
author: user
created: 2026-07-04
track: fast
participants: [claude-1, codex-1]
roles:
  claude-1: implementer
  codex-1: refutation-default reviewer (model-diverse)
status: round-01
---

## Problem / idea

Inspired by Hermes Agent v0.18.0 `/prompt`: compose a long, multi-line prompt in
`$EDITOR` instead of fighting a one-line input box.

Parley TUI analogue: the live-run composer (steering, /answer) is a single-line
input. Long steers and structured answers are painful to type.

Proposal: a `/editor` slash command (and/or ctrl+e keybinding) in the parley TUI
composer that suspends the TUI, opens `$EDITOR` (fallback vi) on a temp file, and on
save+exit puts the file content into the composer (or sends it directly as the steer/
answer being composed), then resumes the TUI.

## Constraints

- Follow the existing bubbletea suspend/resume pattern (tea.ExecProcess) — the TUI
  must repaint cleanly after the editor exits.
- Respect $VISUAL then $EDITOR then fallback; non-zero editor exit ⇒ cancel, keep
  prior composer content, no send.
- Temp file in os.TempDir with 0600 perms, deleted after use; content may be
  sensitive (steer text) — never leave it behind.
- Works in both contexts where the composer accepts free text (steer + answer);
  slash-command autocomplete lists /editor.
- Fast track: ≤5 files, ≤300 LOC, mechanically verifiable (go test + manual TUI
  smoke), fully reversible.

## Non-goals

- No markdown preview, no multi-line inline composer rewrite.
- No persistent drafts across sessions.
