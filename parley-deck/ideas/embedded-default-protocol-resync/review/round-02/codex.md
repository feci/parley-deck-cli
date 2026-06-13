---
agent: codex
idea: embedded-default-protocol-resync
review-round: 2
date: 2026-06-13
reviewed-commit: dedfd27
---

## Fix verification

1. **[MINOR] Drift guard asserts the embedded D2/D3 invariants before normalization.** Verified. `internal/protocol/drift_test.go` defines `assertEmbeddedBootstrapShape`, the exact `embWorkspaceLine`/`embCreatedLine`/`embTransportLine` constants, `assertEmptyTableBody`, and explicit `**Protocol synced:**` count checks: zero in the embedded default, exactly one in the live deck. These checks run after the fail-closed anchor checks and before `normalizeProtocol` is called for either file. The blind spot is closed: an illustrative embedded roster row, altered embedded placeholder line, or stray embedded sync provenance line would now fail before normalization can erase the difference.

2. **[NIT] head-commit corrected.** Verified. `IMPLEMENTATION.md` now records `head-commit: 132271b`. That commit is reachable from `feature/embedded-default-resync` and contains the fix-up code, rather than the stale `bc0af15`.

Verification commands:

- `go build ./...` passed.
- `go test ./internal/protocol/` passed.

## New findings

None.

## Verdict

ACCEPT
