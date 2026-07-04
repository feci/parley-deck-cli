---
idea: tui-editor-composer
status: final
drafter: claude-1
track: fast
date: 2026-07-04
participants: [claude-1, codex-1]
---

## Decision

Add a `/editor` slash command and `ctrl+e` keybinding to the parley live TUI composer
that opens `$VISUAL`/`$EDITOR`/`vi` on a temp file via bubbletea `tea.ExecProcess`, and
on save+clean-exit **fills `m.inputText`** with the content (does NOT send directly).
The existing Enter path then sends it with the correct steer/answer semantics and
visible cues. This folds in codex-1's two round-01 findings.

## Agreed design

1. **Command + key wiring** (`internal/tui/live.go`):
   - Add `/editor` to `commandSpecs` (near `/deck`, `/answer`, ~line 1124) so
     autocomplete, suggestion acceptance, the command hint (~line 1052), and help text
     (~line 2363) pick it up.
   - Bind `ctrl+e` in `updateMain` AFTER confirm-kill and picker interception
     (~line 1300) but BEFORE ordinary editing keys. Do NOT open the editor while
     `picker.Active` — the picker keeps owning keys.

2. **Fill-composer, not direct-send.** The callback sets `m.inputText = content`,
   clears `inputErr` and suggestions, and leaves `composing` + `answerQID` unchanged, so
   the existing Enter routing (slash → active question answer → agent steer → deck
   steer) and the answer/steer visual cues are preserved. Direct-send is out of scope.

3. **`/answer` picker-composition exception** (codex-1 finding #1). When the picker has
   selected a question (`composing=true`, `answerQID` set), `submitInput` handles the
   composing branch before slash dispatch — so typed `/editor` would be sent as the
   literal answer. Fix: add a specific composing-mode exception so `/editor` (and only
   `/editor`) is intercepted before `answerQuestion`. `ctrl+e` works in that state
   regardless.

4. **Multi-line rendering safety** (codex-1 finding #2). `renderInputRow` concatenates
   `label + m.inputText` on one line and `truncateText` does not flatten newlines. Add a
   **display-only preview helper** that flattens/first-lines multi-line content for the
   input row and the `you:` steer echo, while the raw multi-line `inputText` is what gets
   submitted. Show a `[N lines]` affordance when content has newlines.

5. **Editor resolution** — `$VISUAL` → `$EDITOR` → `vi`. Support argument-bearing values
   (`EDITOR="code --wait"`) by splitting on fields; fallback `vi` is a direct
   `exec.Command("vi", path)`. Non-zero exit or exec error ⇒ keep prior composer content,
   status notice, no send.

6. **Temp-file hygiene** — `os.CreateTemp(os.TempDir(), "parley-editor-*.md")`,
   `Chmod(0o600)`, seed with current `m.inputText` (for `ctrl+e`) or empty (for
   `/editor`), close before launch, and `os.Remove` in the callback on EVERY exit path.
   Best-effort only (host SIGKILL mid-edit cannot be covered) — documented.

## Verification (done criteria)

- `go build ./...` and `go test ./internal/tui/...` green.
- New helper-level tests: editor resolution order; temp perms/content/removal;
  non-zero-exit cancel preserves prior input; success replaces `inputText` while
  preserving `composing`+`answerQID`; `/editor` appears in slash suggestions; picker-
  answer path (extend `TestPickerAnswerTwoStepSubmit`).
- Footprint ≤5 files, ≤~300 LOC (fast-track envelope).

## Non-goals

- No markdown preview, no multi-line inline composer widget rewrite, no persistent
  drafts, no direct-send.

## Signoffs

<!-- each participant appends its own block -->

### Signoff: claude-1 — 2026-07-04
Status: ✅ ACCEPT
Drafter. FINAL folds both of codex-1's round-01 findings (picker-mode literal-submit,
multi-line row rendering) into the agreed design. Fast-track envelope respected.

### Signoff: codex-1 — 2026-07-04
Status: ✅ ACCEPT
The final plan covers the picker-mode `/editor` exception and display-only multi-line previews while preserving raw composer submission semantics. I accept the fast-track scope and verification criteria.
