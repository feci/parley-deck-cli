---
agent: claude-1
idea: tui-editor-composer
round: 1
date: 2026-07-04
---

## Summary

Implement `/editor` (+ ctrl+e) via bubbletea's `tea.ExecProcess`, which is the
framework's sanctioned suspend-run-resume mechanism and already repaints cleanly.
Content flows back into the existing composer buffer — we do NOT send directly, so
the user always sees exactly what will be sent and can still edit/cancel.

## Proposed approach

1. **Command surface**: add `/editor` to the existing slash-command autocomplete table
   in `internal/tui/live.go`; bind ctrl+e in the composer keymap to the same action.
2. **Editor resolution**: `$VISUAL` → `$EDITOR` → `vi`. Split respecting simple args
   (e.g. `code --wait`) via `strings.Fields` — document that exotic quoting is out of
   scope v1.
3. **Flow**:
   - write current composer content to `os.CreateTemp(os.TempDir(), "parley-steer-*.md")`,
     `chmod 0600`;
   - `tea.ExecProcess(exec.Command(editor, args..., tmpPath), func(err) tea.Msg { return editorDoneMsg{path, err} })`;
   - on `editorDoneMsg`: err != nil or editor exit != 0 ⇒ delete temp, keep prior
     composer content, status-line notice "editor cancelled";
   - else read file, `strings.TrimRight(content, "\n")`, set composer value, delete
     temp file (defer-safe: delete in the handler regardless of branch).
4. **Multi-line reality check**: the composer is single-line today. v1 rule: the
   returned text is sent as one steer message preserving newlines in the payload;
   the composer input shows a `[3 lines from editor]` placeholder + first line, with
   the full text held in a pending buffer — enter sends the full buffer, esc discards
   it. This avoids rewriting the input widget (non-goal) while making long steers
   real.
5. **Tests**: unit-test editor resolution order + cancel/success handlers with a fake
   `editorDoneMsg`; temp-file perms and cleanup asserted. TUI smoke via existing
   harness patterns in live_test.go.

Estimated footprint: live.go + a small new editor.go + tests ⇒ 3 files, well under
300 LOC. Fast-track constraints hold.

## Concerns / open questions

1. Does the composer context (steer vs /answer free-text) share one input model? If
   /answer uses the picker sub-mode with a free-text branch, /editor must work there
   too — needs a look at answerItems/picker flow during implementation.
2. tmux/screen users: ExecProcess handles tty handoff, but a status repaint glitch on
   resume is possible — acceptable; bubbletea repaints on resume by design.

## Risks

- Leaving temp files with sensitive steer text — mitigated: 0600 + unconditional
  delete in the done handler.
- `code` without `--wait` returns immediately (empty edit) — document; detect common
  no-wait editors is out of scope v1.
