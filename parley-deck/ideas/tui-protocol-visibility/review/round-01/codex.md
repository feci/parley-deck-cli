---
agent: codex
idea: tui-protocol-visibility
round: 1
phase: review
date: 2026-06-12
---

## Findings

### [MAJOR] RebuildDetail still swallows several unexpected disk errors

internal/driver/cursor.go:114, internal/driver/cursor.go:132, internal/driver/cursor.go:135, internal/driver/cursor.go:249, internal/driver/cursor.go:272

D2 requires missing artifacts to become zero values, but unexpected stat/read errors to return a non-nil error with partial detail. `RebuildDetail` only propagates `ReadDir` errors from the round/review directories. It still ignores read errors from `00-prompt.md` and `IMPLEMENTATION.md`, turns any `FINAL.md` read error into the string "FINAL.md is missing", and treats any `os.Stat` error in `fileExists` as absence. That means the async snapshot cannot keep the last good state on these failures, because `BuildProtocolSnapshot` only enters its keep-last path when `RebuildDetail` returns an error.

Suggested fix: replace the boolean/string helpers used by `RebuildDetail` with `(value, exists, err)` style helpers, aggregate the first unexpected error into the returned `PhaseDetail`, and add tests for unreadable `00-prompt.md`, `FINAL.md`, `IMPLEMENTATION.md`, and `review/consensus.md`.

### [MAJOR] Snapshot participants ignore the D5 fallback and live/display split

internal/tui/protocolui.go:143, internal/tui/protosnap.go:218

D5 says participant resolution is `opts -> run.created payload -> 00-prompt frontmatter`, with the wider set used for display and the live set used for waiting math. The current snapshot input copies only `m.opts.Participants`, and `deliveryMatrix` waits on that exact slice. `BuildProtocolSnapshot` never reads participants from the in-memory `run.created` event or from `00-prompt.md`, and it has no separate display-vs-live participant set. If a reattached or older run arrives with empty/stale `LiveOptions.Participants`, the ribbon/protocol tab can show no delivery rows and no waiting agents even though the idea/run artifacts carry the roster.

Suggested fix: resolve participants inside the snapshot builder from the agreed precedence, keep separate live and display participant slices, and add tests where `ProtocolSnapshotInput.Participants` is empty but `run.created` and/or `00-prompt.md` provides the roster.

### [MAJOR] Artifact mode still performs filesystem I/O from View()

internal/tui/live.go:760, internal/tui/live.go:863, internal/tui/live.go:868, internal/tui/live.go:2486

D8/D17 call out View-time I/O explicitly. The new protocol surfaces render from cache, but `View()` still reaches `renderArtifactView`, which calls `loadFocusTail`, which opens and reads the artifact file synchronously. A user can trigger this path with `/artifact`, so a slow or blocked filesystem can still stall rendering, and the promised View-no-disk-read invariant is not tested.

Suggested fix: move artifact loading to an async command/cache path, or explicitly scope and document an exception for artifact mode and add the D17 test that would catch any protocol/ribbon/glyph render-path disk access.

### [MINOR] Explicit buffers_stdout=false is not preserved

internal/tui/protocolui.go:190, internal/tui/protocolui.go:318

D12 says the silent placeholder should use the declared `buffers_stdout` flag first, with the heuristic only as a fallback. `noteRuntimeFlags` only records true values, so an explicit `buffers_stdout=false` from TOML/runtime is indistinguishable from "not declared". After 30 seconds with zero stdout, `agentBuffersStdout` can still show the "buffers all stdout until exit" copy for an agent whose configuration explicitly opted out.

Suggested fix: store `buffers_stdout` as a tri-state, for example `map[string]*bool` or a value plus presence map, and only run the heuristic when no runtime declaration exists. Add a regression test for explicit false.

### [MINOR] D17 run.phase emission matrix is only partially covered

internal/driver/phase_event_test.go:35, internal/driver/phase_event_test.go:62, internal/driver/phase_event_test.go:84, internal/driver/phase_event_test.go:103

D17 requires the driver tests to cover exactly one `run.phase` after save for every phase-changing action, no event for await/surface/error branches, and no event on save failure. The current dedicated tests cover promotion, round-to-consensus drafting, await/surface, and one save-failure path, but they do not assert `run.phase` emission for finalization, reopen, implementation, review-open, either fix-up path, or completion. Those branches currently use `commitCursor`, but the test suite would not catch a future branch bypassing it.

Suggested fix: add a table-driven emission test over all nine commit sites, using the existing fakes from `consensus_test.go` and `impl_test.go`, and assert the D4 payload fields after each successful save.

## Verdict

ACCEPT-WITH-FIXES - the implementation builds and the main architecture matches consensus, but the D2/D5/D8 conformance gaps should be fixed before this release is marked complete.
