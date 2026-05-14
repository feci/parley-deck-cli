---
agent: claude
idea: interactive-agent-mode
review-round: 2
date: 2026-05-14
reviewed-commit: 30d5fe6
responding-to: [review/round-01/claude, review/round-01/gemini, review/round-01/hermes, review/consensus]
---

## Summary
Every agreed fix in `review/consensus.md` is correctly applied at 30d5fe6, and the CRITICAL append-only bypass I raised in round-01 is fully closed. No new blockers introduced. Ready to mark fix-up cycle 1 complete.

## Findings

### [CRITICAL] None
All round-01 CRITICAL findings resolved. Spot checks:
- `before_len` + `before_sha256` are now persisted on `agent.handoff.started` and `agent.handoff.pending` events (`internal/app/consensus_request_signoffs.go:472–475`, `:524–526`).
- `validateResumedConsensusHandoff` (`internal/app/app.go:582–629`) reads the artifact, verifies `sha256Hex(prefix) == before_sha256`, runs `validateAppendOnlyContent`, and rejects `BLOCK` via `consensus.CanonicalStatus`. Resume only emits `agent.handoff.completed` after this passes; failures surface as errors rather than silent "complete" events.
- `eventDataInt` (`internal/app/app.go:631–640`) tolerates the JSON float64 round-trip from the events store, so the integer prefix length survives reload.

### [MAJOR] None
Both round-01 MAJORs resolved:
- Placeholder expansion is unified through `runner.ExpandInteractiveArgs` (`internal/runner/handoff.go:98–107`), consumed by `handoffInstructions` (`:54–56`), `commandPreview` (`internal/app/consensus_request_signoffs.go:813–845`), and `runInteractiveTTY` (`:558`). Printed handoff, dry-run preview, and `spawn-tty` exec now substitute the same set of placeholders.
- Manual handoff output and `writeSignoffHandoff` point at `parley resume <runID>` (`internal/app/consensus_request_signoffs.go:518`, `:546`). The runID is now plumbed through `HandoffOptions.ResumeCommand` and rendered under `Resume:` in `handoff.md`.

### [MINOR] None
- `ValidateArtifact` → `ValidateRoundOneArtifact` rename completed (`internal/runner/validation.go:11`); the only caller `runner.RunRoundOne` (`internal/runner/runner.go:265`) is updated. Contract is honest for this slice.
- Pending manual handoffs return exit code `3` via `errRequestPending` (`internal/app/consensus_request_signoffs.go:55–57`, `:205–208`), documented in `docs/agent-runtime-configuration.md:61`.
- `validateLaunchModes` carries a code comment explaining why manual mode still validates the shared `interactive_*` shape (`internal/app/consensus_request_signoffs.go:352–353`).

### [NIT] None
- The dead `agents.InteractiveCommandOrDefault(agent.Spec)` fallback is gone from `runInteractiveTTY` (`internal/app/consensus_request_signoffs.go:550–571`); `agent.Path` is the sole fallback for `command`, consistent with the guarantee from `requestSignoffAgents`.
- `IMPLEMENTATION.md` frontmatter now uses `implementer: codex` per COOPERATION.md §4 / §11.B.

## Open questions
- No blocking questions. The unit-test surface around `resumePendingConsensusSignoffs` / `validateResumedConsensusHandoff` is currently empty (no matches for `before_sha256`, `validateResumedConsensusHandoff`, or `handoff.pending` in `internal/app/app_test.go`). The consensus did not require tests in this cycle, but I'd flag adding coverage as a candidate for the next slice's IMPLEMENTATION notes — not a blocker for this fix-up cycle.
