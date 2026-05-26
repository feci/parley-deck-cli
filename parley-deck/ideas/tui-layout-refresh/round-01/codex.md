---
agent: codex
idea: tui-layout-refresh
round: 1
date: 2026-05-26
---

## Summary

The current dashboard wastes height by rendering a header, a three-column row,
then full-width Actions/Questions and Footer panels. On short terminals the
user sees fragments instead of a coherent dashboard. The refresh should make
the top viewport dense and scannable: title/status line, two-column main area,
and a compact footer; for low height, collapse into a summary mode.

## Proposed approach

- Add small layout helpers for usable terminal dimensions, content clipping,
  section titles, colored status badges, and compact key/value rows.
- Dashboard:
  - Use one header line with product, transport, focus, and selected run
    status.
  - On normal width, split into a left navigation column and a right content
    column. The left column lists sessions or agents; the right column combines
    run summary, agents, latest events, and actions/questions.
  - On low height, show only header, selected run summary, condensed agent
    statuses, top actions/questions, and one-line footer.
  - Avoid rendering a separate footer box; use a single muted key line.
- Live run view:
  - Use top line with idea/run/status.
  - Use two columns on normal screens: agents/questions left, events/log right.
  - On low height, collapse to agents + latest event + selected question/action
    summary.
- Add color hierarchy:
  - green/success for done/valid/found;
  - amber for running/waiting/questions;
  - red for failed/error;
  - cyan/blue for focus and section titles.
- Tests should assert the compact view appears at low height and that key
  sections remain present on normal width.

## Concerns / open questions

- Bubble Tea viewport widgets would be useful, but a focused string-rendering
  change is safer for this slice.
- Exact terminal height thresholds should be conservative. I would use compact
  dashboard at height under 22 and compact live view under 20.
- Existing tests check substrings, so we can preserve labels while improving
  structure.

## Risks

- Too much color can reduce readability on low-contrast themes. Keep color
  semantic and preserve plain labels.
- Aggressive truncation can hide important IDs. Keep full IDs available in the
  selected detail area when space allows.
- A layout rewrite can break navigation focus assumptions if it changes which
  data is rendered. Keep focus and selection state unchanged.
