---
idea: continuous-run-tui
drafted-by: codex
date: 2026-05-23
status: accepted-with-user-authorized-exception
---

# Consensus

## Protocol exception

This is not a clean multi-agent consensus. The configured non-facilitator agents from `parley-deck/agents.toml` were invoked, but none produced a canonical round artifact because of local auth/model/sandbox blockers recorded in `parley-deck/inbox/codex-to-all_continuous-run-tui_round-01-agent-blockers.md`.

The user explicitly authorized continuing without waiting in `parley-deck/inbox/user-to-all_continuous-run-tui_solo-exception.md`.

## Agreed decisions

- Treat `resume` as read-only state restoration and live-view reopening.
- Introduce `continue` as the mutating workflow-progress command; the initial implementation may be read-only planning until safe step execution exists.
- Add a durable run action planner that reconciles canonical files first, then events, then manifest/cache state.
- Extend `run.json` as a snapshot/cache, not as the single source of truth.
- Model continuation at the workflow-step level: answer open questions, retry one failed/missing participant artifact, draft consensus, request signoffs, finalize, or inspect.
- Keep TUI parallelism process-local for the first useful slice; do not promise detached process survival until a daemon/supervisor exists.
- Make the TUI consume the same planner as the CLI so the command and workbench agree on recommended actions.

## Agreed trade-offs

- Native model-session reattachment is best-effort and not required for the MVP.
- The first implementation should favor inspection and explicit user-triggered actions over automatic phase mutation.
- Sessionstore remains a cross-workspace index; canonical workflow truth remains under `parley-deck/ideas/` and `parley-deck/runs/`.
- Parallel runs need a conservative concurrency cap before TUI exposes broad launch controls.

## Deferred items

- Native agent resume handles.
- Durable detached supervisor/daemon.
- Full TUI action execution for every Parley phase.
- GitHub PR refresh before phase advancement.

## Signoffs

### Signoff: codex - 2026-05-23
Status: accept
Notes: User-authorized exception after peer-agent auth/model/sandbox blockers.

I accept the plan as the implementer. The peer-agent blockers are recorded, and the user explicitly instructed implementation to proceed without waiting.

### Signoff: claude — 2026-05-24
Status: ✅ ACCEPT
Notes: Retroactive review under the user-authorized solo-exception. The proposal is sound from a TUI UX and continuation semantics perspective. The `resume` (read-only) vs `continue` (mutating) split is the right semantic boundary and avoids the overloaded-resume confusion. The file-first reconciliation precedence (canonical artifacts > events > snapshots) is correct and will prevent the state-split bugs flagged in the risks section. The incremental slicing is pragmatic — deferring detached execution and native resume handles until the core planner and step-retry loop are proven is the right order. No blocking concerns.

### Signoff: gemini — 2026-05-24
Status: ✅ ACCEPT
Notes: The proposed design for `internal/runplan` and the file-first reconciliation strategy provides a robust foundation for workflow continuation. I especially support the `resume` (read-only) vs `continue` (mutating) distinction, which clarifies intent for both users and agents. The decision to treat `run.json` as a cache and prioritize canonical artifacts ensures that the state remains auditable and consistent across different environments. The incremental implementation slices are sensible, and I agree with deferring detached execution until a dedicated supervisor is implemented.

### Signoff: hermes — 2026-05-24
Status: ✅ ACCEPT
Notes: From long-running run recovery perspective, the file-first reconciliation and deferral of native resume handles until supervisor exists is the correct conservative stance. The `resume` vs `continue` split and planner approach prevent stale-process pitfalls. Accept under the recorded user-authorized exception.
