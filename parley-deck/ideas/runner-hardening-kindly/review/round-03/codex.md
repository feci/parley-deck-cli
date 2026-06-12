---
agent: codex
idea: runner-hardening-kindly
round: 3
phase: review
date: 2026-06-12
---

## Fix verification

- Cycle-2 fix 1: VERIFIED — `finishACP` now assigns `result.OutputPath = livePath` whenever `publishArtifact` returns a non-empty path before handling the move-back error, so the terminal event reports the LIVE canonical path even on publish failure.
- Cycle-2 fix 2: VERIFIED — `TestMoveAsideInvalidArtifact` now covers the deterministic rename-failure case with a near-`NAME_MAX` basename; after `moveAsideInvalidArtifact`, the canonical path is gone and no recovery file remains.

## New findings

None.

## Dispositions

- I concur that `TestDurableKillEndToEndRealProcess` stays dismissed as a codex seatbelt sandbox artifact. Cycle 2 already records all reviewers re-confirmed the `kern.boottime` restriction, and this fix-up touched only the ACP publish-path assignment and the invalid-artifact rename-failure test.
- Verification run: `go test ./internal/runner/ -run 'TestMoveAsideInvalidArtifact'` passed, `go test ./internal/runner/ -run 'TestMoveAsideInvalidArtifact' -count=1` passed, and `go build ./...` passed.

## Verdict

ACCEPT
