---
agent: codex-1
idea: tui-editor-composer
review-round: 1
date: 2026-07-04
reviewed-commit: 1f46e21
---

## Summary

The core command wiring is mostly correct: typed `/editor` in `/answer` picker composition is intercepted before `answerQuestion`, `ctrl+e` is ordered after confirm-kill and picker handling but before suggestion/menu editing, and `openEditorCmd` appears to create a 0600 temp file and remove it on all normal callback/error paths. I ran `go test ./internal/tui/...`, `go build ./...`, and `go vet ./internal/tui/...`; all passed.

I found one behavior gap against FINAL.md and two review/verification gaps. The largest issue is that only the input row uses `editorPreview`; the agent transcript `you:` echo still renders raw multi-line text.

## Findings

### [MAJOR] Multi-line `you:` steer echo bypasses the display-only preview

`renderInputRow` correctly calls `editorPreview(m.inputText)`, but the agent steer echo still appends raw text at `internal/tui/live.go:1897`:

```go
b.lines = append(b.lines, transcriptLine{Text: "❯ you: " + text, Stream: transcriptSteer})
```

FINAL.md explicitly requires the display-only preview helper for both the input row and the `you:` steer echo, while preserving the raw multi-line `inputText` for submission. Because `styleTranscriptLine` ultimately calls `truncateText`, and `truncateText` trims but does not flatten embedded newlines, an editor-composed multi-line steer can render as multiple transcript lines with only the first line carrying the `you:` prefix. That breaks the single-line echo affordance and can make the transcript display inconsistent even though the submitted steer text is raw.

Concrete fix: keep `SteerRequest.Text` and `steer.Submit` using the raw `text`, but render the transcript echo with `editorPreview(text)`, for example:

```go
b.lines = append(b.lines, transcriptLine{Text: "❯ you: " + editorPreview(text), Stream: transcriptSteer})
```

Add a seam test that submits `"line one\nline two"` and asserts the request text remains raw while `transcriptText` contains the flattened preview with `[2 lines]`.

### [MINOR] Help and key hints do not expose the editor entry points

`/editor` is present in `commandSpecs` and in the slash-command hint, but the help overlay at `internal/tui/live.go:2392` through `internal/tui/live.go:2405` omits both `/editor` and `ctrl+e`. The main/composing/steer hints around `internal/tui/live.go:1056`, `internal/tui/live.go:1066`, `internal/tui/live.go:1068`, and `internal/tui/live.go:1072` also do not mention `ctrl+e`.

This matters because FINAL.md calls out help text and command hints as part of the feature. Users who do not already know the shortcut will not discover the external editor from the normal help surface, and composing/answer mode gives no hint that `ctrl+e` works there.

Concrete fix: add a help row for `ctrl+e` under Input and a slash-command row for `/editor`; add a compact `ctrl+e editor` mention to the main, answer-composing, new-idea-composing, and steer hints where space allows.

### [MINOR] Required editor-path tests are still missing

`internal/tui/editor_test.go` covers editor resolution, preview formatting, command-spec presence, and direct `editorFinishedMsg` handling, but it does not cover several verification cases listed in FINAL.md: temp-file mode/content/removal, success preserving `composing` and `answerQID`, the `/answer` picker-mode `/editor` exception, and `ctrl+e` precedence over suggestions while remaining blocked by picker/confirm-kill. `TestOpenEditorCmdNonNil` only proves a command object is returned.

This matters because the riskiest parts of the feature are exactly the lifecycle and routing paths: private temp-file handling, no accidental answer submission from picker-composition, and shortcut ordering relative to modal/picker/suggestion states.

Concrete fix: factor the editor command path enough to inject a fake editor/runner or temp directory, then add tests that assert 0600 creation and removal, non-zero exit cancellation, success preserving answer composition state, `/editor` in picker-answer composition opening the editor instead of answering, and `ctrl+e` being ignored by picker/confirm-kill but accepted before the suggestion menu.

## Open questions

None.
