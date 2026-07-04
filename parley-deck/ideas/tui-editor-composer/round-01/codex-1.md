---
agent: codex-1
idea: tui-editor-composer
round: 1
date: 2026-07-04
---

## Summary

The feature is feasible in `internal/tui/live.go`, but the safe shape is narrower than "run editor and send directly". The current composer is one `inputText` string on `liveModel` (`internal/tui/live.go:291-304`), and Enter already routes that string through the correct semantics: slash command first, then active question answer, then agent steer, then deck steer (`internal/tui/live.go:1765-1804`). `/answer` picker mode is different: selecting an item closes the picker and sets `composing=true` plus `answerQID` (`internal/tui/live.go:2121-2133`), then `submitInput` sends the answer through `answerQuestion` (`internal/tui/live.go:1778-1782`, `internal/tui/live.go:1898-1910`).

So `/editor` should primarily fill `inputText` and let the existing Enter path send. Direct-send would bypass the visible answer/steer confirmation cues in `renderInputRow` (`internal/tui/live.go:1003-1035`) and make cancel/error handling riskier.

## Proposed approach

Add `/editor` to `commandSpecs` near `/deck` and `/answer` (`internal/tui/live.go:1124-1138`) so autocomplete (`internal/tui/live.go:1157-1179`) and suggestion acceptance (`internal/tui/live.go:1205-1221`) pick it up. Also update the hard-coded command hint (`internal/tui/live.go:1052-1053`) and help text (`internal/tui/live.go:2363-2369`). Add `ctrl+e` in `updateMain` after confirm-kill and picker interception (`internal/tui/live.go:1300-1317`) but before ordinary editing keys (`internal/tui/live.go:1428-1442`). Do not open the editor while `picker.Active`; the picker should keep owning keys until the user selects or cancels it.

Use `tea.ExecProcess` from the key or slash-command handler, with a callback that returns an `editorFinishedMsg`. `Update` already centralizes async messages in one switch (`internal/tui/live.go:413-536`), so adding an `editorFinishedMsg` case fits the model. The callback should read the temp file after a zero exit, remove it, and return `{content, err, canceled}`. On success, set `m.inputText = content`, clear `inputErr`, clear suggestions, and leave `composing` and `answerQID` unchanged. On non-zero exit or exec error, remove the temp file, keep the prior composer content unchanged, and set a status or input error without sending.

Editor resolution should be explicit and testable: `$VISUAL`, then `$EDITOR`, then fallback `vi`. For env-provided values, support common argument-bearing editors such as `EDITOR="code --wait"`; either parse conservatively or run the env command through the user's shell with the temp path passed as a positional parameter. For fallback `vi`, use a direct `exec.Command("vi", path)`.

Temp-file hygiene should be isolated in helper functions. Create the file under `os.TempDir()` with `os.CreateTemp(os.TempDir(), "parley-editor-*.md")`, assert `0600` with `Chmod(0o600)`, write the current composer text, close it before launching the editor, and remove it in the `ExecProcess` callback on every exit path. The initial content for `ctrl+e` should be the current `m.inputText`. For `/editor`, use empty content or optional text after the command; clear the slash command itself before launching so cancel does not leave `/editor` in the composer.

Tests should avoid launching a real editor where possible. Add helper-level tests for editor resolution, temp permissions/content/removal, non-zero cancel preserving prior input, and success replacing `inputText` while preserving `composing` plus `answerQID`. Add a TUI unit test that `/editor` appears in slash suggestions, and a picker-answer test extending the existing `TestPickerAnswerTwoStepSubmit` coverage in `internal/tui/live_test.go:850-880`.

## Concerns / open questions

The biggest flow bug is `/answer` picker composition. Once the picker selects a question, `submitInput` handles `m.composing` before slash dispatch (`internal/tui/live.go:1767-1783`). That means typing `/editor` in that state would currently submit the literal text `/editor` as the answer. To make slash `/editor` work there, `submitInput` needs a specific composing-mode exception before `answerQuestion`. Otherwise only `ctrl+e` works for picker-selected answers.

The active-agent open-question row is different: when not in picker composition, slash commands are dispatched before active-tab question answering (`internal/tui/live.go:1791-1800`), so `/editor` can be intercepted there.

Multi-line content also exposes a rendering issue. `truncateText` trims but does not flatten embedded newlines (`internal/tui/app.go:144-152`), while `renderInputRow` concatenates `label + m.inputText` directly (`internal/tui/live.go:1027-1035`). A saved multi-line draft can therefore break the single-line input row. Use a display-only preview helper for input rendering, while keeping the raw multi-line `inputText` for submission. The steer echo path also appends `"you: " + text` as one transcript line (`internal/tui/live.go:1862-1867`), so it should either preview or split display lines.

## Risks

Temp cleanup cannot be absolute if the whole process is killed while the external editor is open. The implementation can guarantee best-effort cleanup on editor exit, non-zero status, read failure, and exec failure, but not on host crash or SIGKILL.

Using shell semantics for `$VISUAL` or `$EDITOR` improves compatibility with `code --wait`, `vim -n`, and similar values, but it means a malicious editor environment can execute arbitrary commands. Since `$EDITOR` is user-controlled local configuration, this is acceptable if the file path is passed safely and not string-concatenated unsafely.

Returning content to the composer rather than direct-sending adds one extra Enter, but it preserves the current color-flipped answer warning and steer labeling (`internal/tui/live.go:1003-1035`) and makes cancellation unambiguous. Direct-send should remain out of scope for the fast-track version.
