---
agent: agy
idea: runner-hardening-kindly
round: 3
phase: review
date: 2026-06-12
---

## Fix verification

- **Cycle-2 Fix 1 (`finishACP` live-path mirror):** VERIFIED. [finishACP](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runner/acp.go#L218) in [internal/runner/acp.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runner/acp.go) has been updated to mirror [finalizeExecResult](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runner/runner.go#L596). It assigns `result.OutputPath` to the non-empty `livePath` returned by `publishArtifact` even when the move-back operation fails, ensuring that the ACP terminal event reports the live canonical path.
- **Cycle-2 Fix 2 (`TestMoveAsideInvalidArtifact` rename-failure):** VERIFIED. [TestMoveAsideInvalidArtifact](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runner/hardening_test.go#L325) in [internal/runner/hardening_test.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runner/hardening_test.go) has been extended to cover the rename-failure case. It asserts that when the rename fails (simulated using a destination path exceeding `NAME_MAX`), the invalid artifact is deleted from the canonical path and no recovery file remains.

## New findings

None.

## Dispositions

We concur that the `TestDurableKillEndToEndRealProcess` test failure under the codex seatbelt sandbox (due to restricted access to `sysctl kern.boottime` / boot-id verification) is an environment-specific artifact. It remains dismissed.

## Verdict

ACCEPT (zero agreed fixes remaining)
