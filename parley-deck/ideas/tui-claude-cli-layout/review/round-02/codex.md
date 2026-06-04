---
agent: codex
idea: tui-claude-cli-layout
review-round: 2
date: 2026-06-04
responding-to: [claude/IMPLEMENTATION fix-up-cycle-1, hermes/review-round-01]
---

## Position changes since prior review round

The fix-up satisfies review/consensus.md AF1-AF6 and resolves my round-01 findings.

Verified `git diff e763ab8..HEAD`, re-read `internal/tui/live.go` and
`internal/tui/live_test.go`, and checked the dead-state symbols with `rg`.
`renderTabStrip` now windows around the active tab and emits clipped-side `…+N`
markers; `agentBuffer` records `os.FileInfo` and `refreshBuffers` reloads when
`!os.SameFile`; retired modal modes/fields/branches are gone; `shift+↑` and
`shift+↓` line scroll without disturbing plain `↑`/`↓` tab switching; AF2/AF4/AF5
tests and the AF6 comment update are present.

Attempted checks with workspace-local caches:
`GOCACHE=$PWD/.gocache GOMODCACHE=$PWD/.gomodcache go build ./...`,
`go vet ./...`, `go test ./internal/tui`, and `go test ./...`. They remain
blocked by restricted network plus uncached modules (`bubbletea`, `lipgloss`,
`go-toml/v2`, `yaml.v3`). In `go test ./...`, cached packages outside the blocked
dependency chain passed.

Carried deferrals are acceptable: mouse-wheel scroll, richer textinput editing,
`/agent`/`/tab`, Events tab, queued-steer execution, and live ACP delivery remain
outside this review cycle.

## Responses to other reviewers

### @hermes

Agreed. Your prior tab-strip, line-scroll, key-routing, and comment concerns are
covered by AF1, AF4, AF5, and AF6 respectively. I also agree with your round-02
position that plain `↑`/`↓` line scroll remains intentionally dismissed by
consensus in favor of the owner-approved tab mapping plus `shift+↑`/`shift+↓`.

## Updated findings

### [MAJOR] Active tab can be clipped out of the tab strip

RESOLVED. `renderTabStrip` no longer joins all tabs and right-truncates the
result. It resolves the active stable tab id, creates a window around that tab,
expands to neighbors within budget, and marks omitted tabs with `…+N`.
`TestTabStripKeepsActiveTabVisible` covers a far-right active agent in an
overflowing strip.

### [MINOR] Log rotation detection can miss replacement files that grow past the old offset

RESOLVED. `agentBuffer` now stores the `os.FileInfo` captured at load, and
`refreshBuffers` reloads when the current file is not `os.SameFile` with the
recorded one, even if the new file has already grown past the old offset.
`TestBufferReloadsOnFileReplace` covers the replacement-past-offset case.

### [NIT] Retired modal state and answer-mode rendering remain in the model

RESOLVED. The retired live modes, answer/compose/focus fields, unreachable
`modeAnswerQuestion` render branch, and orphaned `previewLineBudget` are removed.
The remaining modes are `modeOverview` and `modeHelp`.

No new CRITICAL/MAJOR/MINOR/NIT findings.

ZERO agreed fixes remain; Phase 7 can finalize.

## Open questions

None.
