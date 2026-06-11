---
idea: tui-protocol-visibility
cycle: 1
drafted-by: claude
date: 2026-06-12
reviewed-commit: 65cf1ca
outstanding_agreed_fixes: 11
---

## Agreed fixes

From codex (review/round-01/codex.md):

1. **[MAJOR] RebuildDetail error surfacing (D2).** Probe FINAL.md /
   IMPLEMENTATION.md / review/consensus.md with a stat helper that propagates
   non-NotExist errors; surface 00-prompt.md and IMPLEMENTATION.md frontmatter
   read errors; a FINAL.md that stats as existing but reads as "missing" is an
   error. Aggregate the first unexpected error into the returned detail. Tests
   for the unreadable-file cases.
2. **[MAJOR] Participants precedence + display/live split (D5).** Resolve the
   roster inside BuildProtocolSnapshot: live = opts → run.created payload →
   00-prompt frontmatter (first non-empty); display = ordered union. Delivery
   rows render the display set; the waiting list counts only the live set.
   Tests with empty input participants.
3. **[MAJOR→scoped] View-time I/O in /artifact (D8/D17).** The artifact view
   predates this idea (1.20.0) and is an explicit user action with a bounded
   tail read; making it async is deferred (follow-up below). Agreed fix here:
   add the D17 render-path test proving the NEW surfaces (ribbon, glyphs,
   status segment, protocol panes) render correctly with nonexistent
   RunDir/IdeaDir paths (pure cache), and document the /artifact exception in
   IMPLEMENTATION.md.
4. **[MINOR] buffers_stdout tri-state.** Record declared true AND false from
   the run.created runtime payload; the 30s/0B heuristic applies only when no
   declaration exists. Regression test for explicit false.
5. **[MINOR] Full run.phase emission matrix (D17).** Extend the existing
   driver phase tests so every one of the 9 commit sites asserts its run.phase
   event (finalized, reopened, implemented, review-opened, both fix-up paths,
   complete — in addition to the covered promoted/consensus-drafted).

From agy (review/round-01/agy.md):

6. **[MAJOR] Home Attention badge (D14).** Per-idea attention resolved from
   the idea's latest run in homeRuns, rendered via the existing attentionBadge
   next to the phase chip.
7. **[MAJOR] Placeholder indentation (D12).** Two-space indent on every
   renderSilentPlaceholder line.
8. **[MINOR] Step-8 pipeline title.** Use the snapshot's StepName for the
   current row so a fix-up cycle shows "Fix-Up", not "Complete".
9. **[MINOR] Collapsed ribbon round format.** `(R02)` per the D9 string — drop
   the `/total` denominator.
10. **[MINOR] Narrator/activity noise.** Humanize the line: drop the raw event
    type when the summary text is self-describing; otherwise map the type to a
    concise verb ("hermes started", "round-02 complete (3/3)").
11. **[NIT] Home double status.** When the phase chip is unavailable, render
    "—" in the chip column instead of repeating idea.Status.

## Deferred follow-ups

- ASCII fallback for the braille spinner (agy NIT): no terminal-capability
  probe exists in the codebase; the rest of the UI already assumes UTF-8
  (✓/●/⚠). Follow-up if a real terminal report shows breakage.
- Spinner animates at 1 fps (agy NIT): bounded by the 1s elapsed tick by
  design; a faster animation tick is not worth the render churn. Documented.
- Async artifact view (codex MAJOR remainder): move /artifact's loadFocusTail
  to a tea.Cmd cache — follow-up idea (pre-existing behavior since 1.20.0).
- impl.checks event + review-cycle naming normalization (carried over from the
  design consensus).

## Dismissed findings

- None. (internal/runner TestDurableKillEndToEndRealProcess was pre-excluded
  as the known codex-sandbox artifact and was not reported.)

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: claude — 2026-06-12
Status: ✅ ACCEPT
Notes: All 11 fixes accepted for fix-up cycle 1; three deferrals recorded above.

### Signoff: codex — 2026-06-12
Status: ✅ ACCEPT
Notes: Consensus faithfully triages my five findings; the scoped artifact-mode fix plus async deferral is acceptable.

### Signoff: agy — 2026-06-12
Status: ✅ ACCEPT
Notes: Triage faithfully represents my findings, including the visual alignment fixes and spinner deferrals.
