---
agent: agy
idea: tui-protocol-visibility
round: 1
date: 2026-06-11
---

## 1. Collapsed Ribbon Line

### Proposal
**ACCEPT** the design seed concept, but re-order and polish the layout for immediate scannability. At a glance, the user needs to know **where we are**, **progress**, **who is blocking**, and **what is next**.

Exact 1-line layout:
`◆ Ph 2: Cross-Review (R02) · Delivered 1/3 · Pending: agy, hermes · Next: consensus`

If the state is degraded (stale or disk fallback):
`◆ Ph 2: Cross-Review (R02) · Delivered 1/3? · STALE (reconciled 12m ago) · Pending: agy`

*Rationale*: 
- `◆` (U+25C6) provides a clean starting delimiter.
- Phase index and name (`Ph 2: Cross-Review`) must come first to establish context.
- `Delivered X/Y` shows progression immediately.
- `Pending` (renamed from `waiting` for cleaner tone) explicitly lists blockers.
- `Next` sets expectations for the upcoming state transition.

---

## 2. Expanded (Ctrl+P) 3-Line View

To stay within 100 columns and remain highly scannable, we must condense the pipeline names and display clear, honest metadata.

### Proposal
```
Pipeline: Kick ── R01 ── XRev [▶] ── Cons ── Final ── Impl ── Revw ── RCon ── Fixp
Delivery: claude (✓ 12:04) · codex (✓ 12:05) · agy (working 2m) · hermes (pending)
System:   Next: consensus · Reconciled 14s ago (Disk fallback)
```

*Rationale*:
- **Line 1 (Pipeline)**: Condensed phase names fit easily within 80 cols. The `[▶]` marker represents the active phase.
- **Line 2 (Delivery detail)**: Shows exact status and times.
- **Line 3 (System metadata)**: Explicitly calls out virtio-fs cache lag (`Reconciled 14s ago`) and whether the data is from events or a disk scan.

---

## 3. Tab Activity Glyphs

Using text descriptors like `RUN` or `FIN` clutters the tab strip. Replacing them with distinct, single-character glyphs enhances layout breathing room.

### Proposal
The exact glyph set to ship, optimized for terminal font portability:
- **Pending**: `○` (U+25CB, White Circle) — e.g., `  hermes ○`
- **Running (actively writing output)**: Braille spinner `⠋` (rotating frames)
- **Running (silent / buffering stdout)**: `⧗` (U+29D7, White Hourglass) — e.g., `  agy ⧗`
- **Finished (delivered artifact)**: `✓` (U+2713, Check Mark) — e.g., `  codex ✓`
- **Failed**: `✗` (U+2717, Ballot X)
- **Killed**: `⊘` (U+2298, Circled Division Sign)
- **Skipped**: `–` (U+2013, En-dash)
- **STALE** (Process unresponsive or lost): `⚠` (U+26A0, Warning Sign) — e.g., `  hermes ⚠`

---

## 4. Narrator Lines

### Proposal
- **Trim the Allowlist**: Narrator lines must be high-signal. Restrict them to:
  - Phase start/end transitions.
  - Agent status changes (started, finished, failed, killed).
  - HITL questions posed and user answers submitted.
  Exclude all command execution details, file reads, and acp chunk writes.
- **Tab Weaving**: Delivery lines and phase changes **must** appear in **every** tab's transcript. If a user is focused on a silent tab (like `agy`), seeing `── codex delivered round-02/codex.md ✓ ──` in their current view instantly confirms the run is alive and moving.

---

## 5. The Buffered-Agent Placeholder

To eliminate the "dead tab" look for agents that buffer stdout (like `agy`), the empty transcript space should be populated with a dynamic placeholder.

### Proposal
When stdout is empty, render a structured activity block using internal event tracking:

```
┌────────────────────────────────────────────────────────────────────────┐
│ agy buffers all stdout until exit.                                     │
│ Live Status: RUNNING (elapsed: 2m 14s)                                 │
│ Stdout: 0 bytes (buffered)   Stderr: 412 bytes (live)                  │
│                                                                        │
│ Activity Log:                                                          │
│ 12:04:10 · Agent started                                               │
│ 12:04:12 · Tool execution: view_file (internal/tui/live.go)            │
│ 12:04:45 · Tool execution: grep_search (func .*render)                 │
│ 12:05:12 · Thinking (deep reasoning mode)...                           │
└────────────────────────────────────────────────────────────────────────┘
```

*Rationale*: Showing the stderr byte counter and a log of recent tool calls provides immediate assurance that the agent is actively executing.

---

## 6. Degraded-State Honesty

### Proposal
We must never hide staleness or disk discrepancies behind Ctrl+P.
- **STALE Runs**: If a run's attention state becomes `STALE` (no process or >10m inactivity), the collapsed ribbon must prefix itself with `[STALE]` and color-flip the phase tag to yellow/red.
- **Disk Fallback / Absent run.json**: Append a `?` to the delivery count (e.g., `Delivered 1/3?`) to signal that we are working off file scans rather than verified memory events.
- **Reconciliation Age**: If reconciliation lag exceeds 30s, display the age directly in the collapsed ribbon: `· Reconciled 45s ago`.

---

## 7. Status Line Grammar

### Proposal
**COUNTER-PROPOSE** the `ph=N:name wait:...` grammar to make it slightly cleaner to parse and read while retaining the key fields.

Adjusted format:
`ph=2:xrev-r02 wait=agy,hermes`

- `ph=2:xrev-r02` is more compact than `cross-review-r02`, saving precious status line space.
- Use `=` instead of `:` for list values (`wait=agy,hermes`) to align with typical key-value query parameters.

---

## 8. Home Phase Column

### Proposal
Replace the raw `idea.Status` string with aligned columns displaying the structured Phase and the Attention badge.

Layout:
```
  Ideas
    tui-protocol-visibility       Ph 2: Cross-Review    ● RUNNING
    another-completed-idea        Ph 8: Complete        ✓ DONE
    broken-idea-slug              Ph 5: Implement       ✗ FAILED
    waiting-for-user-idea         Ph 3: Consensus       ⚠ ACTION
```

*Rationale*:
This gives the user a clear picture of their active directory landscape without opening any individual ideas.
