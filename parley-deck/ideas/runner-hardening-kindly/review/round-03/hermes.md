---
agent: hermes
idea: runner-hardening-kindly
round: 3
phase: review
date: 2026-06-12
---

## Fix verification
1. finishACP (internal/runner/acp.go): VERIFIED — livePath from publishArtifact is now assigned to result.OutputPath before the error branch, matching finalizeExecResult.
2. TestMoveAsideInvalidArtifact (internal/runner/hardening_test.go): VERIFIED — rename-failure case added (NAME_MAX basename triggers ENAMETOOLONG; artifact removed, no recovery file).

## New findings
- None (zero regressions from cycle-2 fix-up).

## Dispositions
- TestDurableKillEndToEndRealProcess: concurs with prior dismissal (sandbox sysctl kern.boottime restriction only).

## Verdict
ACCEPT (zero agreed fixes remaining)