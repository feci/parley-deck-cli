---
agent: agy
idea: embedded-default-protocol-resync
review-round: 2
date: 2026-06-13
reviewed-commit: dedfd27
---

## Fix verification

1. **[MINOR] Drift guard asserts the embedded D2/D3 invariants before normalization**
   - Verified that `internal/protocol/drift_test.go` now implements `assertEmbeddedBootstrapShape` and checks the occurrences of `**Protocol synced:**` before normalizing.
   - `assertEmbeddedBootstrapShape` enforces exact matches for the generic header lines: `embWorkspaceLine` (`**Workspace:** `<workspace-name>``), `embCreatedLine` (`**Created:** `<date> — created by parley init``), and `embTransportLine` (`**Transport:** `github-pr``).
   - `assertEmptyTableBody` enforces that the §2 roster and host-handle tables in the embedded default have a separator row but an empty body (no table rows beginning with `|` below the separator).
   - Verified `**Protocol synced:**` occurrence checks guarantee exactly 1 in the live deck and 0 in the embedded default.
   - These checks run prior to normalization, ensuring that any in-zone edits (such as added roster rows, modified placeholders, or stray sync lines in the embedded copy) will fail the guard instead of being normalized away.

2. **[NIT] head-commit corrected**
   - Confirmed `IMPLEMENTATION.md` successfully updated `head-commit` to `132271b`, which is a real and reachable commit in the branch history containing the fix-up implementation.
   - Confirmed that `go build ./...` and `go test ./internal/protocol/` pass successfully.

## New findings

- None. No regressions were introduced during fix-up cycle 1.

## Verdict

ACCEPT
