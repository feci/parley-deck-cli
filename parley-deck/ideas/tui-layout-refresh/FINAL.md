---
idea: tui-layout-refresh
date: 2026-05-26
status: final
participants: [codex, claude, gemini, hermes]
---

# Final

The TUI layout refresh is implemented and reviewed.

## Outcome

- Dashboard now uses terminal height and switches below 30 rows into a compact layout that avoids boxed vertical stacking.
- Normal dashboard layout is now a two-column view: sessions on the left, selected run workspace on the right, and a compact footer.
- Live run view now uses two columns in normal mode and a clipped compact layout below 30 rows.
- Semantic badges and section titles make state, risk, and attention easier to scan.
- Existing keyboard behavior is preserved.
- Review feedback was addressed, including removal of obsolete render methods, normal-mode footer warning restoration, compact header consistency, and threshold-boundary tests.

## Verification

- `go test ./internal/tui` passed.
- `go test ./...` passed.

## Review

- Claude: conditional approval; actionable findings addressed.
- Gemini: no blocking findings.
- Hermes: no blocking findings.
