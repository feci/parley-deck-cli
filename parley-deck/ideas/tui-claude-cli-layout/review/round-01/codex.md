---
agent: codex
idea: tui-claude-cli-layout
review-round: 1
date: 2026-06-04
---

## Summary

The default TUI surface is largely aligned with FINAL.md D3/D4: printable keys append to the input row, arrows/tab route tabs, scroll keys affect only transcript tabs, Enter answers the active agent's open question before steering, Status Enter records a deck steer, esc clears then detaches, and ctrl+c still cancels. The per-agent transcript buffers are bounded and loaded before normal render paths via `newLiveModel`, tab switches, and tick/event handlers.

I found one spec-level issue around tab-strip overflow, one rotation edge in the buffer refresh path, and some remaining dead modal state from the retired layout.

Checks attempted with workspace-local caches: `GOCACHE=$PWD/.gocache GOMODCACHE=$PWD/.gomodcache go build ./...`, `go vet ./...`, and `go test ./...`. All three were blocked for TUI/config/pipeline packages because modules were uncached and network is unavailable (`proxy.golang.org: no such host`); `go test` still ran and passed the packages whose dependencies were already available.

## Findings

### [MAJOR] Active tab can be clipped out of the tab strip

`renderTabStrip` builds all tabs in fixed order and then calls `truncateText` on the joined string (`internal/tui/live.go:258-265`). With enough agents or a narrow terminal, a selected far-right agent or the Status tab can be completely hidden even though that tab is active. FINAL.md requires overflow to keep the active tab plus neighbors visible and show a clipped-side `... +N` marker; this implementation explicitly deviates from that and can make keyboard navigation feel broken because the top strip stops showing the current location.

Concrete fix: replace plain right truncation with an overflow-aware tab window keyed by stable tab IDs. Always include the active tab, include adjacent tabs while width permits, and render a left and/or right clipped marker such as `... +N` for omitted tabs. Add a narrow-width test with the active tab near the right edge and another for Status.

### [MINOR] Log rotation detection can miss replacement files that grow past the old offset

`refreshBuffers` treats rotation as `os.Stat(path).Size() < b.offset` (`internal/tui/live.go:598-603`). That catches truncation and many rotations, but not a replaced stdout file that has already grown to at least the old offset by the next 250 ms refresh. In that case `readAppendedLines` seeks to the stale offset in the new file, dropping the new file's earlier lines from the transcript. The code comment promises reload on log truncation/rotation, so this is a correctness gap in the per-agent buffer path.

Concrete fix: track file identity in `agentBuffer` when loading, and reload when identity changes. On Unix this can use `syscall.Stat_t` dev/inode; a portable fallback can combine size decrease with modtime/created path heuristics if needed. Add a test that replaces the log file with a different file whose size is greater than or equal to the previous offset and verifies the buffer reloads from the new tail.

### [NIT] Retired modal state and answer-mode rendering remain in the model

The old top-level modes and fields are still present even though the new default surface only uses `modeHelp`: `modeAgentDetail`, `modeCompose`, `modeAnswerQuestion`, `answerText`, focus fields, compose fields, and `previewLineBudget` (`internal/tui/live.go:42-50`, `internal/tui/live.go:93-110`, `internal/tui/live.go:841-842`, `internal/tui/live.go:861-863`). This is not currently breaking behavior, but it makes future routing changes harder to audit because dead branches still look like valid UI states.

Concrete fix: remove unused live modes and fields, delete the unreachable `modeAnswerQuestion` branch in `renderQuestionsPane`, and update the printable-key test to assert against the remaining help/default mode shape rather than `modeOverview`.

## Open questions

- Is the plain tab-strip truncation intended to ship as an accepted v1 product deviation, or should FINAL.md's active-plus-neighbors overflow behavior remain a Phase 6 acceptance requirement?
