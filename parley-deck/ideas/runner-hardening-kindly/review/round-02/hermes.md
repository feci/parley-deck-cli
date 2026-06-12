---
agent: hermes
idea: runner-hardening-kindly
round: 2
phase: review
date: 2026-06-12
---

## Fix verification

fix 1 (ACP activity + attempt_id threading + retry-once loop): VERIFIED — acp.go now marks initialized/session_opened/prompt_completed as activity; attempt_id threaded through procctl marker and started payload; runAgent ACP branch has the attempt loop.
fix 2 (conditional snapshot cleanup + Abandon keeps dir + publishArtifact always live path): VERIFIED — reviewsnapshot.go now retains snapshot on MoveArtifactBack failure (marker removed), publishArtifact always returns live path, terminal events report live path even on move failure.
fix 3 (moveAsideInvalidArtifact unique destination + remove-on-rename-failure + TestMoveAsideInvalidArtifact): VERIFIED — runner.go implements unique destination and remove-on-failure path; hardening_test.go has the test.
fix 4 (RunFixup via execAgentProcess + supervision + watchdog + no retry): VERIFIED — phase58.go switched RunFixup to execAgentProcess path with groups/procctl/watchdog; no retry loop present.
fix 5 (failEarly classified payload): VERIFIED — failEarly now routes through classifyFailure producing failure_class/recovery_hint.
fix 6 (consult frontmatter + ledger session_id + EFFECTIVE timeout_ms): VERIFIED — consult.go writes session_id (even when empty) and effective timeout_ms into frontmatter.
fix 7 (TestClassifyFailure exact class/hint contract lock): VERIFIED — hardening_test.go asserts exact pairs from the implemented table.

## New findings

None.

## Dispositions

TestDurableKillEndToEndRealProcess fails under the codex seatbelt sandbox (sysctl kern.boottime restricted). This is an environment artifact, not a code defect. The test passes in a normal shell. Disposition: dismissed.

## Verdict

ACCEPT (zero agreed fixes remaining)