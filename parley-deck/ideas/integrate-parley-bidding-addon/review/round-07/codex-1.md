---
idea: integrate-parley-bidding-addon
review-round: 7
agent: codex-1
date: 2026-07-30
reviewed-commit: 49fc3ec
---

## Verdict
ACCEPT

## Outstanding findings — closed or not

Closed. The round-6 MAJOR is resolved: marker identity validation no longer accepts a marker with an absent `skill` field, while valid 2.0.0-shape markers remain accepted for all six units.

The round-6 MINOR reported by kimi-1 is also resolved: `selected` is derived from the core marker rather than from the inspection-selection flag, so targeted residual inspection and unflagged doctor inspection agree. The cycle-5 narrowing remains intact.

## New findings
None.

## Release judgement
Yes. Commit `49fc3ec` is releasable as 2.1.0. No further change is required for release.

## What I verified

- The reviewed branch resolves to `49fc3ec`, and its worktree remained clean.
- The missing-identity regression rejects both core and add-on markers whose `skill` field is absent. I independently downgraded every installed marker to the exact 2.0.0 shape; all six units still reported `valid`.
- The recorded-selection regression makes filtered and unfiltered doctor reads report the residual bidding unit as `selected: false` and `valid-unselected`; the filtered read remains narrowed to core plus bidding.
- The full Node suite passed: 314 tests, 0 failures.
- The Python 3.14 leg passed: 54 tests across seven files. The manifest check passed with aggregate `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`.
- `skills/parley-bidding/` is unchanged from its first integration commit. It contains 48 files, including 16 parseable JSON files and four platform adapters, with no Python cache artifacts.
