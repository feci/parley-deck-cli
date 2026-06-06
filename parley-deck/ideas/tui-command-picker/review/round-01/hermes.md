---
agent: hermes
idea: tui-command-picker
phase: review
round: 1
date: 2026-06-06
---

## Summary
looks-good — implementation is faithful to FINAL.md on key routing, state reset, reclamp, dedup, empty states, and test coverage; no CRITICAL/MAJOR issues.

## Findings
MINOR: internal/tui/live.go:1246 (updatePicker) — "up"/"down" cases mutate Index before the final reclamp at 1268, but reclamp already runs on every mutation path; no off-by-one because reclamp clamps both before Enter selection (1252) and before every render (1366-1381). Correct.

MINOR: internal/tui/live.go:1294 (openItems) — dedup via seen map on Slug/RunID is present and matches FINAL §3; no duplicates possible.

MINOR: internal/tui/live.go:1204 (pickerRows) + 421 (renderTabbed) — transcript height shrink by pickerBlock line count plus floor-to-1 and min(8,transcriptHeight) guarantees no overflow on tiny terminals; empty-state path also returns a single-line box.

NIT: internal/tui/live.go:145 (pickerState.reclamp receiver) — pointer receiver is appropriate and consistent with mutation sites; value receiver on filtered() is correct since it returns a copy.

NIT: internal/tui/live_test.go:779 (TestPickerReclampClampsIndex) — covers the exact FINAL §9.5 contract; all 11 §9 items have direct assertions.

## Verdict
ACCEPT — no blocking items. All key-routing, two-step answerQID clearing (activateRun 391, clearComposition 1222, esc 737/1243, select 1277/1280, submit 1038), and re-entrancy paths are clean.