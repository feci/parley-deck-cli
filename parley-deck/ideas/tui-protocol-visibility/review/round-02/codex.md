---
agent: codex
idea: tui-protocol-visibility
round: 2
phase: review
date: 2026-06-12
---

## Fix verification

1. VERIFIED - RebuildDetail error surfacing: `RebuildDetail` now uses `statRegular` and `readFrontmatterFieldErr`, keeps the first unexpected error, and treats an existing-but-unreadable `FINAL.md` as an error (`internal/driver/cursor.go:121`, `internal/driver/cursor.go:134`, `internal/driver/cursor.go:140`, `internal/driver/cursor.go:147`, `internal/driver/cursor.go:269`, `internal/driver/cursor.go:288`); `TestRebuildDetailSurfacesReadErrors` covers unreadable `00-prompt.md`, `IMPLEMENTATION.md`, and `FINAL.md` (`internal/driver/phase_event_test.go:147`).
2. VERIFIED - Participants precedence plus display/live split: `BuildProtocolSnapshot` resolves live participants as opts -> `run.created` -> `00-prompt.md`, renders the ordered union, and waits only on the live set (`internal/tui/protosnap.go:148`, `internal/tui/protosnap.go:161`, `internal/tui/protosnap.go:176`, `internal/tui/protosnap.go:289`); `TestSnapshotParticipantsFallback` covers event and frontmatter fallback (`internal/tui/protosnap_test.go:165`).
3. VERIFIED - Scoped render-path pure-cache test plus `/artifact` exception: `TestRenderPathPureCache` renders ribbon, protocol panes, status segment, and tab glyphs with nonexistent paths (`internal/tui/protosnap_test.go:210`); `IMPLEMENTATION.md` documents `/artifact` as the bounded user-triggered View-time tail-read exception and defers the async rework (`parley-deck/ideas/tui-protocol-visibility/IMPLEMENTATION.md:102`).
4. VERIFIED - `buffers_stdout` tri-state: `noteRuntimeFlags` records explicit true and false, and `agentBuffersStdout` uses the declaration before the heuristic (`internal/tui/protocolui.go:176`, `internal/tui/protocolui.go:370`); `TestBuffersStdoutTriState` and `TestNoteRuntimeFlagsTriState` cover explicit false (`internal/tui/protosnap_test.go:241`, `internal/tui/protosnap_test.go:261`).
5. VERIFIED - Full 9-site `run.phase` emission matrix: the shared `wantLastRunPhase` helper is applied to finalized and reopened consensus branches plus implemented, review-opened, complete, and both fix-up paths (`internal/driver/phase_event_test.go:36`, `internal/driver/consensus_test.go:135`, `internal/driver/consensus_test.go:262`, `internal/driver/impl_test.go:140`, `internal/driver/impl_test.go:193`, `internal/driver/impl_test.go:248`, `internal/driver/impl_test.go:281`, `internal/driver/impl_test.go:335`), in addition to the existing promoted and consensus-drafted checks (`internal/driver/phase_event_test.go:53`, `internal/driver/phase_event_test.go:80`).
6. VERIFIED - Home Attention badge: `renderHome` looks up the matching `homeRuns` entry and appends `attentionBadge(r.Attention)` without render-path disk I/O (`internal/tui/live.go:1058`, `internal/tui/live.go:1063`).
7. VERIFIED - Placeholder indentation: `renderSilentPlaceholder` funnels every placeholder line through `writeLine`, which prepends two spaces (`internal/tui/protocolui.go:673`).
8. VERIFIED - Step-8 pipeline title: the current row uses `p.StepName`, so step 8 can render `Fix-Up` during fix-up cycles and `Complete` when done (`internal/tui/protocolui.go:606`).
9. VERIFIED - Collapsed ribbon round format: `renderRibbon` now emits `(R%02d)` without the `/total` denominator, and the pure-cache test guards the D9 string (`internal/tui/protocolui.go:427`, `internal/tui/protosnap_test.go:224`).
10. VERIFIED - Narrator/activity noise: `friendlyEventText` maps event types to human-readable verbs or self-describing summaries, and both narrator lines and recent activity use it (`internal/tui/protocolui.go:216`, `internal/tui/protocolui.go:733`).
11. VERIFIED - Home double status: the home chip column renders a neutral placeholder when no phase chip exists instead of repeating `idea.Status` (`internal/tui/live.go:1058`).

## New findings

None.

## Verdict

ACCEPT

Verification commands run:

- `go build ./...`
- `go test ./internal/driver/ ./internal/tui/`
