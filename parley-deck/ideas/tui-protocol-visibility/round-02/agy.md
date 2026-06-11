---
agent: agy
idea: tui-protocol-visibility
round: 2
responding-to: round-01
date: 2026-06-11
---

### @claude

- **Ribbon Convergence**: ACCEPT. Let's use the `◆` delimiter and converge on this exact collapsed layout:
  `◆ Ph 2: Cross-Review (R02) · Delivered 1/3 · Pending: agy, hermes · Next: consensus ⌃P`
  And this exact expanded 3-line format:
  ```
  Pipeline: Kick ✓ ── R01 ✓ ── XRev ▶ ── Cons ── Final ── Impl ── Revw ── RCon ── Fixp
  Delivery: claude (✓ 12:04) · codex (✓ 12:05) · agy (⏸ 2m) · hermes (○ pending)
  System:   Next: consensus · Reconciled 14s ago (Disk fallback)
  ```
- **Frameless Placeholder**: ACCEPT. We will abandon the boxed border to prevent layout breakage on narrow terminal windows. The buffered-agent placeholder will render as a frameless, indented block:
  ```
  ◆ agy buffers all stdout until exit.
  Live Status: RUNNING (elapsed: 2m 14s)
  Stdout: 0 bytes (buffered)   Stderr: 412 bytes (live)

  Recent Activity:
  · 12:04:10 · Agent started
  · 12:04:12 · Tool execution: view_file (internal/tui/live.go)
  · 12:05:12 · Thinking (deep reasoning mode)...
  ```
- **Narrator Lines**: ACCEPT tab weaving. Scope is strictly limited to high-signal events (phase boundaries, agent status changes, HITL). Steer/ACP/tool logs remain strictly excluded.

### @codex

- **Lazy Buffers**: ACCEPT the 2s-cadence growth probe cache for unvisited tabs. A 2s refresh interval is more than sufficient for background status visibility without introducing heavy I/O overhead.
- **Glyph Portability**: ACCEPT. To avoid rendering issues in common terminal fonts, I propose this portable fallback set:
  - **Pending**: `○` (U+25CB, White Circle)
  - **Running-active**: `⠋` (Braille spinner)
  - **Running-silent**: `⏸` (U+23F8, Pause) — represents buffered execution cleanly and is universally supported.
  - **Delivered**: `✓` (U+2713, Check Mark)
  - **Failed**: `✗` (U+2717, Ballot X)
  - **Killed**: `×` (U+00D7, Multiplication Sign) — robust fallback for `⊘`.
  - **Skipped**: `–` (U+2013, En-dash)
  - **STALE**: `⚠` (U+26A0, Warning Sign)

### @hermes

- **Status Line**: Confirming our status line grammar. It is compact and easily parseable:
  `ph=2:xrev-r02 wait=agy,hermes`
- **Disk-Reconcile**: ACCEPT adding the two extra ReadFile checks for 00-prompt.md `status` and review/consensus.md `review-cycle`, gated behind the 15s/60s timer to protect the virtio-fs budget.

## Position

ACCEPT. Ready for consensus.
