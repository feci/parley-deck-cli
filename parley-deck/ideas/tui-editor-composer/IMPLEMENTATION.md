---
idea: tui-editor-composer
status: implemented
implementer: claude-1
started: 2026-07-04
completed: 2026-07-04
branch: parley-deck-cli#editor-composer-design
head-commit: 219b949
design-pr: https://github.com/feci/parley-deck-cli/pull/69
implementation-pr: same
---

## Summary of work

Added a `/editor` slash command + `ctrl+e` keybinding to the live TUI composer. Both
open `$VISUAL`/`$EDITOR`/`vi` on a 0600 temp file via `tea.ExecProcess`; on a clean
exit the file content is dropped into `m.inputText` (NOT sent), so the existing Enter
routing keeps steer/answer/launch semantics and their visual cues.

## Implementation plan / checklist

- [x] New `internal/tui/editor.go`: `editorFinishedMsg`, `resolveEditor`,
      `openEditorCmd` (temp-file 0600 + seed + readback + unconditional remove),
      `editorPreview` (multi-line → single-line preview).
- [x] `internal/tui/live.go`:
  - `/editor` added to `commandSpecs`, command hint, and slash help line.
  - `ctrl+e` handled in `updateMain` after picker/confirm interceptors, before the
    suggestion menu and ordinary editing keys.
  - `/editor` case in `runCommand` (clears the token, opens editor empty).
  - composing-mode exception in `submitInput` so typed `/editor` while answering a
    picker-selected question opens the editor instead of submitting the literal text
    (codex-1 finding #1).
  - `editorFinishedMsg` case in `Update`: success fills input + clears suggest;
    cancel keeps prior input + status; error keeps prior input + surfaces error.
  - `renderInputRow` uses `editorPreview` so a multi-line draft never breaks the
    one-line row while the raw value is still submitted (codex-1 finding #2).
- [x] Tests (`internal/tui/editor_test.go`): editor resolution order; single/multi-line
      preview; `editorFinishedMsg` success/cancel/error handling; `/editor` in specs +
      `/edi` autocomplete.
- [x] Checks run: `go build ./...` (green), `go vet ./internal/tui/` (clean),
      `go test ./internal/tui/` (ok), `gofmt -l` (clean).

## Deviations from FINAL.md

None. Both of codex-1's round-01 findings are implemented as agreed. Direct-send
remains out of scope; multi-line content is submitted raw and only previewed flattened.

## Notes for reviewers

- `openEditorCmd` returns an opaque `tea.Cmd`; its readback/cancel branches are covered
  via `editorFinishedMsg` handler tests rather than by launching a real editor in CI.
- Temp-file cleanup is best-effort (`defer os.Remove` in the ExecProcess callback);
  a host SIGKILL while the editor is open cannot be covered — documented in FINAL §6.
- `$EDITOR`/`$VISUAL` are split on fields (`code --wait` works). Exotic shell quoting
  in those vars is out of scope (documented).
