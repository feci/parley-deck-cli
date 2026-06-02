---
idea: 2026-06-02T21-14-49-finish-the-6-rem
status: implemented
implementer: claude
started: 2026-06-02
completed: 2026-06-02
branch: parley-deck-cli#feat/pipeline-finish-6
head-commit: pending
design-pr: n/a
implementation-pr: pending
---

## Summary of work

All six items implemented per FINAL.md / consensus.md, additive, full suite green.

- **Item 3 — `execute --json`**: `internal/app/pipeline_cmd.go` adds `--json` emitting `ExecuteJSON` (schema_version 1; status dry_run|pending_gate|ready_for_harness; provider_call, effect_digest, idempotency_key, gate). Ledger written before print. Contract doc `references/EXECUTION_CONTRACT.md`.
- **Item 5 — WinGet manifests**: generated `winget-pkgs/.../ParleyDeckCli/{1.6.0,1.7.0,1.8.0}/` and `parley-deck-skill/packaging/winget/.../ParleyDeckSkill/1.3.0/` (version+installer+locale). `InstallerSha256` = `PLACEHOLDER-FILL-FROM-RELEASE-ASSET` (never invented).
- **Item 4 — `pipeline watch`**: `internal/pipeline/watch_eval.go` (`SignalSource`+`FileSignalSource`, `LoadMonitoring`, threshold eval, `BreachRecord` persistence under `breaches/`, lifecycle) + `Signal.Class`; `runPipelineWatch` in app opens low-risk remediation ideas or notify/gate, marks recovered breaches resolved. One pass per invocation (cron-friendly).
- **Item 2 — Phase 5-8 runner**: `internal/runner/phase58.go` (`RunImplementation`, `RunReviewRound`, `BuildImplementationPrompt`, `BuildReviewPrompt`, `ValidateImplementationArtifact`, `ValidateReviewArtifact`); `Options.Phase`+`Options.ArtifactName`; `runAgent` dispatches prompt/validation/path by phase.
- **Item 1 — auto for all kinds**: `runPipelineAuto` dispatch — deliberation/watcher drive+advance (watcher prints watch hint); action drives plan then STOPS `needs_human_gate` (never executes/advances past); implementation drives Phase 5 + review round-01 via the runner then STOPS `needs_artifact`.
- **Item 6 — transport/decider/DAG**: `Block.Transport` + `Manifest.EffectiveTransport`; `Manifest.Decider` + `AutoApproveWithDecider` (low-risk non-prod only; prod never); `Manifest.Execution` (linear|dag) + `validateDAG` (Kahn cycle check) + `dag.go` `ReadyBlocks`/`AllBlocksComplete`.

## Implementation plan / checklist
- [x] Files changed: internal/app/pipeline_cmd.go (+test), internal/pipeline/{manifest,gate,executor,watcher,dag,watch_eval}.go (+dag_test), internal/runner/{runner,phase58}.go (+phase58_test), references/EXECUTION_CONTRACT.md, winget manifests.
- [x] Checks run: `go build ./...`, `go vet ./...`, `go test ./...` — all green.
- [x] Review/risk notes: invariants held — additive, linear back-compat, agents-write-markdown/driver-executes, non-bypassable production (boundary gate when block risk=production; execute gate always; decider/auto never approve production).

## Deviations from FINAL.md

- Action blocks: auto STOPS at the action block (never auto-runs even dry-run execute, since the target is not derivable from the manifest without the plan); the human runs `execute`. Safer than the FINAL note's "auto dry-run then stop"; still never executes production.
- Implementation blocks: auto drives Phase 5 + review round-01 then stops for human-guided review consensus + fix-up (the zero-agreed-fixes loop is not fully unattended).
- DAG: shipped the validation + `ReadyBlocks` engine (+ per-block transport/decider); the linear `Driver.Advance` remains the default executor. Full parallel multi-active DAG auto-drive over `ReadyBlocks` is wired at the engine level but `auto` still walks the linear cursor.

## Notes for reviewers

Focus: (1) production-safety invariants in `AutoApproveWithDecider`, `runPipelineExecute` gate handling, and `runPipelineAuto` action/impl stops; (2) `validateDAG` cycle/endpoint correctness; (3) `thresholdBreached` operator parsing; (4) backward-compat of `Options`/manifest additions (zero-value defaults). Tests: `internal/pipeline` (dag/watch/gate/effects/provider/manifest), `internal/runner` (round/phase58), `internal/app` (PipelineAuto*).
