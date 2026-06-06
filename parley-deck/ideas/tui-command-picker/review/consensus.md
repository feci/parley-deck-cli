---
idea: tui-command-picker
phase: review-consensus
drafter: claude
date: 2026-06-06
participants: [claude, codex, agy, hermes]
---

## Review consensus

Phase 6 review (round-01) raised findings; the implementer (claude) applied fix-up
cycle 1; Phase 8 re-review (round-02) confirms every blocking item is resolved with no
new findings. **All three reviewers (codex, agy, hermes) ACCEPT at round-02.** Zero
agreed fixes remain → the implementation is ready to mark complete.

### Agreed fixes (all applied in fix-up cycle 1, verified FIXED in round-02)
- **MAJOR — FINAL §8 async refresh** (codex, agy #3): `refreshPickerItems()` rebuilds an
  open picker's `Items` from cached data, preserving `Filter`/`Index`/`Offset` and
  re-clamping; wired into the `questionsMsg` handler (pickerAnswer) and `refreshHomeRuns`
  (pickerOpen). Test `TestPickerAnswerRefreshesOnBackgroundUpdate`.
- **MAJOR — answer lifecycle** (agy #1): `submitInput` no longer pre-clears
  `composing`/`answerQID`; `answerQuestion` clears them only on a successful write, so a
  failed answer keeps the user in compose. Test `TestPickerAnswerFailureKeepsCompose`.
- **MAJOR — lose-on-error** (agy #2): `launchIdea` keeps `composing` + the typed task on
  launch failure.
- **MAJOR — picker hint** (agy #4): `renderInputRow` switches the hint to
  `↑/↓ select · type filter · Enter choose · esc cancel` while the picker is active.
- **MINOR — window math** (codex, agy #5): `renderPicker` reclamps a local copy
  (`reclamp` is the single source of truth); `WindowSizeMsg` re-clamps an open picker.
- **NIT — label order** (codex): run rows `run  <run-id>  <idea>  [status]`; idea rows
  `idea  <slug>  [status]` (FINAL §3).

### Dismissed (agreed by all)
- **NIT — mixed value/pointer receivers** (agy #6): dismissed. The value-receiver
  Update-path methods (which return the mutated model copy) plus pointer-receiver
  mutating helpers match the existing Bubble Tea idiom in `live.go`; not a defect.
  codex and hermes concurred; agy agreed in round-02.

### Deferred follow-ups (non-blocking, not for this slice)
- Async **refresh-on-open** (fetch fresh sources the moment a picker opens) remains a
  deferred follow-up per FINAL §8; the current slice uses cached data + reactive rebuild
  on ticks, which the reviewers accepted.

### Verification
`go build ./... && go vet ./... && go test ./...` green (repo-local GOCACHE/GOMODCACHE).

## Signoffs

(each participant appends its own block below — append-only, do not edit others')

### claude — ✅ ACCEPT
2026-06-06. As implementer/drafter I confirm all six agreed fixes are applied and the
one NIT is dismissed by agreement. Build/vet/full test suite green; the two new tests
lock the async-refresh and answer-failure contracts. Ready to mark complete and ship.

### codex — ✅ ACCEPT
The agreed fixes from fix-up cycle 1 are correctly reflected in the implementation and
the round-02 re-review already verified the prior blockers as resolved. I accept marking
the idea complete.

### agy — ✅ ACCEPT
All four of my major findings—covering the answer lifecycle, launch error recovery, asynchronous picker refreshing, and contextual picker hints—have been successfully resolved and verified. I accept marking this idea complete.


### hermes — ✅ ACCEPT
Key-routing and interaction fidelity invariants remain intact after the fix-up cycle. Round-02 acceptance is reaffirmed with no blockers. Idea marked complete.
