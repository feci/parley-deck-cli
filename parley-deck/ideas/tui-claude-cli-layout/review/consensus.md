---
idea: tui-claude-cli-layout
review-cycle: 1
drafted-by: claude
date: 2026-06-04
reviewed-commit: e763ab8
---

Synthesis of Phase 6 review-round-01 from codex and hermes (review/round-01/).
Agreed fixes applied in fix-up cycle 1 on branch feat/tui-claude-cli-layout.

## Agreed fixes

### AF1 — Tab strip must keep the active tab visible (codex MAJOR + hermes MAJOR)
`renderTabStrip` joins all tabs then plain-truncates, so a far-right active tab
(many agents / narrow terminal) can be clipped out entirely. Fix: overflow-aware
window keyed by stable ids — always include the active tab, add neighbors while
width permits, and show `… +N` markers for omitted sides. Tests: active tab near
the right edge at narrow width; active `Status`.

### AF2 — Reload a buffer when its log file is replaced (codex MAJOR)
`refreshBuffers` only reloads when `size < offset`; a replaced stdout file that
already grew past the old offset by the next tick is read from a stale offset,
dropping its earlier lines. Fix: record the file identity (`os.FileInfo`) at load
and reload when `!os.SameFile(old, cur)` (portable). Test: replace the log with a
different file larger than the prior offset → buffer reloads from the new tail.

### AF3 — Remove the retired modal state (codex MINOR/NIT + hermes NIT)
Delete the now-dead modes/fields from the old layout: `modeAgentDetail`,
`modeCompose`, `modeAnswerQuestion`; `answerText`/`answerErr`/`logPreview`/
`focus*`/`compose*` fields; the unreachable `modeAnswerQuestion` branch in
`renderQuestionsPane`; orphaned `previewLineBudget`. Keep `modeOverview` (default)
and `modeHelp`.

### AF4 — Add non-conflicting line scroll (hermes MAJOR, resolved)
The owner-mandated `↑/↓ = switch tabs` (consensus D3) deliberately drops arrow
line-scroll. To recover fine scroll without breaking that, add `shift+↑`/`shift+↓`
= scroll one line (and keep `PgUp/PgDn`/`ctrl+u`/`ctrl+d`/`Home`/`End`). Plain
`↑/↓` stay tab switches. Documented in `/help`.

### AF5 — Strengthen the routing tests (hermes MINOR)
Add tests for: `Enter` answer-before-steer when a question is open vs steer when
not; `/deck` records a deck steer; `/status` switches to the Status tab; a scroll
key drops follow on an agent buffer. (`TestKeyRoutingPrintableAppendsNotHotkey`
stays for the printable-append contract.)

### AF6 — Comment drift (hermes NIT)
Fix the bounded-scrollback comment to match the constants (20000 lines / 4 MiB).

## Deferred follow-ups

- Mouse-wheel scroll; `bubbles/textinput` cursor/paste editing; `/agent`,`/tab`
  commands; Events tab — all carried (FINAL deferred list).
- (Carried) executing queued steers, live ACP delivery, opt-in thoughts.

## Dismissed findings

- hermes "no arrow line-scroll" as a blocker: the plain-`↑/↓`-for-tabs mapping is
  an explicit owner request ratified in consensus D3; AF4 adds `shift+↑/↓` line
  scroll as the compromise rather than reverting the owner's mapping.

## Signoffs

### claude — ACCEPT (2026-06-04)
Accept the synthesis. AF1 (active-tab-visible) and AF2 (replace-file reload) are
correct must-fixes; AF3 finishes the dead-state cleanup; AF4 recovers line scroll
without breaking the owner's ↑/↓ mapping; AF5/AF6 tighten tests/docs. No blockers.

<!-- codex appends its signoff after re-review -->

<!-- hermes appends its signoff after re-review -->
