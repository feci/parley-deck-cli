---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-11
inspected-commit: 9ea4e0f
artifact-kind: implementation call-site audit
not-a-signoff: true
---

# Remaining shared budget integration boundary

This is a source-inspection map for the next implementation step, not evidence
that these callers are already budgeted. The budget package has real operator
inspection/reconciliation callers but no model-launch or driver-action callers.

## Actual process paths

| Path | Current process boundary | Common request boundary |
| --- | --- | --- |
| Manual `agents exec` and headless signoffs | `internal/runner/launch.go` `AgentCommand.Start` | `trackedCommandFor` → `beginProtocolLaunch` → `beginLaunch` |
| Normal rounds, steer, consult/goal-check and phase-5/8 commands | `internal/runner/runner.go` `execAgentProcess` owns a raw `exec.Cmd.Start`, called also by `consult.go` and `phase58.go` | `beginProtocolLaunch` → `beginLaunch` |
| ACP participant | `internal/runner/acp.go` `runACPAgent` → `internal/acp/spawn.go` `Spawn` | `beginProtocolLaunch` → `beginLaunch` |
| Measured interactive session | `internal/runner/launch.go` `RunInteractive` owns a raw `exec.Cmd.Start` | `beginProtocolLaunch` → `beginLaunch`, or an explicit recorded configuration refusal |
| Real readiness/runtime probe | `internal/app/preflight.go` and `internal/app/app.go` → `ProbeCommandFor` → `AgentCommand.Start` | `beginLaunch` with explicit `probe-only` context |
| Printed interactive/manual handoff | `internal/runner/handoff.go` `WriteHandoffPacket`; no child process | `beginLaunch`, followed by `unobserved-handoff` |

The first four rows are distinct spawn implementations. Budgeting only
`AgentCommand.Start` would miss normal rounds/steer, ACP and direct terminals.
Budgeting every `beginLaunch` without distinguishing the last row would charge
an unobserved handoff as a process. A failed build/start must still retain its
attempt identity; failed evidence writes must not grant execution.

Existing `WithLaunchInfo` construction occurs separately in normal runner,
ACP, phase58, consult, steer and preflight. Passing a policy only in one newly
constructed LaunchInfo would be lost when another surface replaces it.

## Implementation order and invariants

1. Resolve one durable budget scope for the idea across run IDs, resume and
   worktrees. A separate ledger in each worktree's runtime directory cannot
   enforce a shared cap. Keep non-idea probes/consult scope explicit. A scope
   mismatch or unknown legacy accounting must stop; do not silently create a
   fresh zero-count grant after moving a worktree or resuming a run.
2. Attach attempt reservation to the common request/lifecycle boundary, with
   explicit handoff-versus-process intent. Keep unique invocation identity,
   before-spawn publication and failed attempts. Do not infer that no process
   ran merely because `started.json` is absent: a start-evidence write can fail
   after the actual child has already started.
3. Reconcile observed cost at terminal handling without refunding the action.
   An expired process context must not skip durable terminal accounting. A
   conservative operator ceiling is distinct from provider-observed cost;
   unknown telemetry remains null. A monetary ceiling with no known maximum
   for a new attempt must visibly refuse rather than guess zero.
4. Connect charged fixup and cross-review actions separately from model launch
   count. Preserve the existing maximum of cursor/driver marker safety state.
   Migrate legacy charges explicitly before accepting new work. All failed
   attempts stay spent, BLOCK cannot buy another round, and verification of
   the last allowed attempt remains legal at equality with the cap.
5. Persist driver progress and elapsed origin across re-entry. The current
   `Run` locals (`start := time.Now()`, `steps := 0`) still reset those limits;
   the tested budget-store wall clock does not fix that until it is wired.
6. Make actual operator extensions explicit controls with durable decisions,
   preserving existing track gates and charged history. The new cost-recovery
   control addresses unknown prices only; it is not yet a count/time extension.
7. Add the opt-in independently confirmed, patch-linked regression trigger.
   Repeated criticism of unchanged code is not a new patch regression.

## Required negative verification

Before claiming AC-B1/T1: actual process fixtures through every row above;
failed build/start and start/terminal evidence-write failures; ordinary and
cross-worktree concurrent callers; fresh run/resume with old cursor charges;
manual/driver/BLOCK entrypoints; inclusive 5/6 fixup and 3/4 cross-review limits;
unknown/conflicting cost; corrupt state; no fake grant from participant
frontmatter; exact operator replay; no reset after failure, retry or inspection.

The >=20-real-attempt gate remains separate from those fixtures. Twelve real
terminal attempts are currently reconciled; test processes and printed
handoffs are not live model attempts. This map does not certify all surfaces
from the current manual-facilitation inventory.

## Retrieval provenance

Graphify vocabulary `[budget, reservation, ledger, launch, cursor]` returned
the ledger and cursor source locations with a 1000-token traversal budget.
The graph is navigation evidence only. The exact process map above was then
verified by directly reading the named current source files and searching
their non-test spawn/constructor call sites; it is not inferred from a graph
edge or an earlier facilitator summary.
