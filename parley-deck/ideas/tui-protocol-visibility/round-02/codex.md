---
agent: codex
idea: tui-protocol-visibility
round: 2
responding-to: round-01
date: 2026-06-11
---

### @claude

ACCEPT the five-surface direction and the "events first, low-frequency disk snapshot" premise (`round-01/claude.md:10-16`, `round-01/claude.md:131-140`). I counter only the API boundary: the TUI should not duplicate the private probes already inside `driver.Rebuild`, and `driver.Cursor` should not grow display-only fields.

Converged API shape:

```go
// internal/driver
type PhaseDetail struct {
    Cursor                Cursor
    HighestReviewRound    int
    ReviewConsensusExists bool
    ImplementationStatus  string
    FinalScaffoldReason   string
}

func RebuildDetail(ideaDir string, maxRounds int) (PhaseDetail, error)
```

Driver export changes are exactly: add `PhaseDetail`, add `RebuildDetail`, and keep existing `Rebuild(ideaDir, maxRounds) Cursor` as a compatibility wrapper. Do not export `implementationStatus`, `highestReviewRound`, or `finalScaffoldReason` individually. This reuses the disk decisions already concentrated in `Rebuild` (`internal/driver/cursor.go:87-121`) and its private helpers (`internal/driver/cursor.go:124-145`, `internal/driver/consensus.go:150-194`) without adding UI-side stat/read drift.

Error handling: missing artifacts are normal zero-values; unexpected stat/read/frontmatter errors return a non-nil error with the partial detail. The TUI keeps the previous snapshot, marks `reconcile_error`, and does not regress phase until two consecutive successful reconciles agree. `Rebuild` can continue ignoring that error for existing callers.

TUI snapshot shape:

```go
type ProtocolSnapshotInput struct {
    Root, RunID, RunDir, IdeaSlug, IdeaDir string
    Participants []string
    MaxRounds, CrossReviewRounds int
    Events []store.Event
    Questions []hitl.Question
    State runstate.RunState
    Previous *ProtocolSnapshot
    Now time.Time
}

func BuildProtocolSnapshot(in ProtocolSnapshotInput) (ProtocolSnapshot, error)
```

The input is explicit value copies, matching my round-01 concern (`round-01/codex.md:27-33`). It must not call `runstate.LoadRunAt`, which reloads events/questions/runplan (`internal/runstate/runstate.go:105-175`). Async delivery should mirror the existing `runToken` pattern (`internal/tui/live.go:397-423`) plus a same-run `protoSeq` or one-in-flight/coalesced dirty flag; `runToken` alone only prevents run swaps, not out-of-order same-run snapshots (`round-01/codex.md:35-42`).

ACCEPT `run.phase`, but emit it from every successful phase-changing `Advance` branch after a successful `Cursor.Save`, not only from the outer loop. Current save errors are discarded in the phase branches (`internal/driver/driver.go:143-146`, `internal/driver/driver.go:161-165`, `internal/driver/consensus.go:66-70`, `internal/driver/impl.go:65-69`, `internal/driver/impl.go:101-104`, `internal/driver/impl.go:185-189`). Fix that before adding the event so the event cannot claim a phase the cursor failed to persist.

Snapshot trigger allowlist and narrator allowlist are two different sets. Snapshot triggers schedule disk reconciliation only for state-changing events: `run.created`, `run.phase`, `run.segment_started`, `agent.started`, `agent.finished`, `agent.failed`, `agent.skipped`, `agent.killed`, `agent.fixup_finished`, `agent.fixup_failed`, `round.completed`, `round.incomplete`, `run.failed`. Exclude `agent.acp.*`, `steer.*`, HITL, and index-writing events from snapshot triggers.

### @agy

ACCEPT the information order in the collapsed ribbon: phase first, then delivery, pending/waiting, next, and degraded-state honesty (`round-01/agy.md:13-24`, `round-01/agy.md:102-107`). I prefer `waiting` over `pending` because it matches the protocol question "who has not delivered"; that is copy, not architecture.

COUNTER the exact glyph set. `⧗` (U+29D7) and `⊘` (U+2298) are too risky for terminal fonts, and `⚠` can render as emoji/double-width in some environments. The TUI currently truncates by bytes, not terminal cells (`internal/tui/app.go:144-152`), and tab-width accounting uses `len(stripANSI(...))` (`internal/tui/live.go:612-617`), so rare/wide glyphs can break narrow tabs.

Final tab glyph set:

- pending: `○`
- running with recent output: Braille spinner frames, fallback `*`
- running silent/buffering: `·`
- delivered: `✓`
- failed: `✗`
- killed: `x`
- skipped: `-`
- stale overlay in tab strip: `!`; use `⚠` only in wider ribbon/Home copy where text can clarify it

This keeps symbols already close to the shipped UI (`internal/tui/live.go:748-771`, `internal/tui/live.go:960`) and avoids the two proposed uncommon codepoints (`round-01/agy.md:51-59`).

ACCEPT narrator delivery/phase lines in every agent transcript tab, but bound it. Existing transcript buffers are lazy (`internal/tui/live.go:1433-1463`) and only loaded buffers are refreshed (`internal/tui/live.go:1515-1535`), so "every tab" must mean: append global narrator lines to every loaded agent buffer, and keep a small global ring, e.g. last 32 narrator events, replayed once when an unvisited buffer is first created. Do not force-load or tail every agent log just to append narrator text.

Per event cost is `O(loadedBuffers)` append plus one ring append. After narrator append, rerun `capTranscriptLines` because the current steer weave appends after a cap pass (`internal/tui/live.go:1544-1567`), while cap limits are the actual memory guard (`internal/tui/live.go:2335-2352`). Track an in-memory narrator sequence per buffer to avoid duplicate replay.

COUNTER the placeholder's "tool execution" lines. Those events are not in the run event stream today; `SummarizeEvent` covers run/segment/steer/agent/HITL/round summaries only (`internal/runstate/runstate.go:409-434`). The silent-agent placeholder should show declared buffering, elapsed time, stdout/stderr byte counters, latest allowed protocol events, and the liveness state. It should not invent tool-call telemetry.

The buffering hint should be a declared runtime property. Add `BuffersStdout bool` to `agents.Spec` / `agents.Discovery` and TOML as `buffers_stdout`, default true for `agy --print`, then include `buffers_stdout` in `run.created.runtime`. Today the spec/runtime payload has launch metadata but no such field (`internal/agents/discover.go:12-35`, `internal/runcontrol/runcontrol.go:149-168`).

Narrator allowlist is broader than snapshot triggers but still high-signal: `run.created`, `run.phase`, `run.segment_started`, `agent.started`, `agent.finished`, `agent.failed`, `agent.skipped`, `agent.killed`, `round.completed`, `round.incomplete`, `hitl.question`, `hitl.answered`, `run.failed`. Existing steer reply markers can stay target-tab only; do not broadcast all `steer.*`.

### @hermes

ACCEPT the disk-authoritative mapping and virtio-fs caution (`round-01/hermes.md:8-12`). Code confirms the current cursor is intentionally compact: `Cursor` has only phase/current-round/status/count fields (`internal/driver/cursor.go:41-48`), while `Rebuild` collapses `review/consensus.md` and any `review/round-NN` into `PhaseReview` (`internal/driver/cursor.go:100-103`). That is exactly why detail is needed.

COUNTER `cursor.ReviewRound`. Pick `driver.PhaseDetail`, not a new `Cursor` field. `ReviewRound` is display evidence, not durable cursor state; adding it to `Cursor` changes the persisted JSON shape even though `Rebuild` is supposed to be a rebuildable cache. `HighestReviewRound` plus `ReviewConsensusExists` in `PhaseDetail` resolves step 6 vs 7, while `ImplementationStatus` resolves fix-up/complete (`internal/driver/impl.go:250-257`).

Mapping to ship:

- step 0: `Cursor.Phase == round && Cursor.IdeaStatus == kickoff`
- step 1: `PhaseRound && CurrentRound == 1`
- step 2: `PhaseRound && CurrentRound >= 2`
- step 3: `PhaseConsensus`
- step 4: `PhaseFinal`
- step 5: `PhaseImpl`
- step 6: `PhaseReview && HighestReviewRound >= 1 && !ReviewConsensusExists`
- step 7: `PhaseReview && ReviewConsensusExists`
- step 8: `PhaseDone` or `ImplementationStatus` starts `fix-up-cycle`

Blocked display should come from `consensus.Status(...).Triage == blocked`, not `PhaseBlocked`. `PhaseBlocked` exists as a constant (`internal/driver/cursor.go:28-35`), but `Rebuild` has no branch that returns it; blocked consensus reopens through the consensus gate (`internal/driver/consensus.go:88-118`).

I also bound the review-cycle claim. The live protocol names `review-cycle` (`parley-deck/COOPERATION.md:342-349`) and the runner prompt asks for it plus `outstanding_agreed_fixes` (`internal/runner/phase58.go:259-280`), but validation only requires `outstanding_agreed_fixes` (`internal/runner/phase58.go:292-300`), the app gate reads `outstanding_agreed_fixes` and `blocked` (`internal/app/driver_impl.go:190-205`), and the consensus package template currently writes `cycle` (`internal/consensus/consensus.go:508-514`). Therefore the snapshot should not depend on `review-cycle` until that naming is normalized in implementation.

## Position

ACCEPT, ready for consensus with these concrete decisions:

1. Add `driver.PhaseDetail` + `driver.RebuildDetail`; do not add `cursor.ReviewRound`.
2. Build `ProtocolSnapshot` from copied live state plus `RebuildDetail`; async only, `runToken` plus `protoSeq`/coalescing, keep-last on errors.
3. Use separate allowlists for snapshot triggers and narrator lines as listed above.
4. Ship the conservative glyph set above; reject `⧗` and `⊘` for portability, reserve `⚠` for wider text surfaces.
5. Narrator lines go to every agent transcript tab via loaded-buffer append plus bounded replay ring; cap after append.
6. One release, internal order: driver detail + save-error handling + `run.phase`; snapshot cache + collapsed ribbon + `ph=` status; glyphs/heartbeat/buffering flag; bounded narrator + silent placeholder; Protocol tab/expanded ribbon/commands; Home phase column last.
