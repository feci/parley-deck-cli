---
idea: 2026-06-02T21-54-52-continue-the-4-r
status: implemented
implementer: claude
started: 2026-06-02
completed: 2026-06-02
branch: parley-deck-cli#feat/pipeline-followups-1234
head-commit: pending
design-pr: n/a
implementation-pr: pending
---

## Summary of work

Items 1, 2, 4 implemented (additive, full suite green); item 3 escalated as blocked.

- **Item 4 (agy launch robustness)**: generic "write-first" preamble added to all round/review/impl prompts; **stdout-capture fallback** in `runAgent` — if the artifact is absent but captured stdout starts with `---`, persist it as the agent-authored artifact and record `agent.stdout_fallback`. Strict `---` validation rejects narration. Tests: recovers valid stdout; rejects narration.
- **Item 2 (unattended Phase 8 loop)**: `pipeline.ReviewAgreedFixes` (reads `outstanding_agreed_fixes`/`blocked` frontmatter) + pure `Phase8Decision`; `runner.RunReviewConsensus` (drafter writes the machine contract; `Options.Overwrite`), `runner.RunFixup` (Phase 8 implementer re-invoke); `autoDriveImplementationBlock` now loops review→consensus→decide→fix-up bounded by maxCycles=3, fails closed if the contract field is absent; `blockCompleteFunc` recognizes impl completion (IMPLEMENTATION.md status:complete or review consensus 0 fixes); auto advances past a completed impl block.
- **Item 1 (parallel multi-active DAG)**: additive cursor `ready_blocks[]`/`active_blocks[]` (back-compat); `Driver.ComputeDAGStep` (per-wave gate-approved ready set + awaiting gates); `parley pipeline auto --max-active K` (default 4) drives DAG pipelines in concurrent waves of deliberation/watcher blocks (stops at ready action/impl blocks with guidance); `status` shows ready/active.
- **Item 3 (WinGet upstream PR)**: BLOCKED — no GitHub releases/`.exe` assets for v1.6–1.9. Escalated via inbox/claude-to-user_..._winget-blocked.md (create releases? choice). No invented hashes, no placeholder PR.

## Implementation plan / checklist
- [x] Files: internal/runner/{runner,phase58}.go (+round_test, phase58_test), internal/pipeline/{run,executor,review,dag*}.go (+review_test, dag_step_test), internal/app/pipeline_cmd.go (+test).
- [x] Checks: go build ./..., go vet ./..., go test ./... — all green.
- [x] Risk notes: additive cursor fields (old runs load); production gates non-bypassable in DAG waves + decider (low-risk only); stdout fallback strict `---`-gated; Phase 8 loop bounded + fail-closed.

## Deviations from FINAL.md
- Parallel DAG drives deliberation/watcher waves concurrently; ready action/implementation blocks stop the DAG drive with guidance (action=gated execute; impl=single-active auto), rather than parallelizing those kinds. Safer.
- Item 3 not completed (infra-blocked); escalated rather than faked.

## Notes for reviewers
Focus: production never auto-approved (AutoApproveWithDecider low-risk-only decider; DAG waves); stdout fallback strictness; Phase8Decision/ReviewAgreedFixes correctness + fail-closed; cursor back-compat (ready/active additive); ComputeDAGStep topological + gate correctness. Tests: internal/pipeline (review/dag_step), internal/runner (Stdout/Narration/Impl/Review), internal/app (PipelineAuto*).
