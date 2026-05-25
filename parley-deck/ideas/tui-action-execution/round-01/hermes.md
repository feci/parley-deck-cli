---
agent: hermes
idea: tui-action-execution
round: 1
date: 2026-05-25
---

## Summary
The tui-action-execution slice extends the existing planner (runplan + runaction) so the TUI can safely trigger actions such as retry-agent, draft-consensus, request-signoffs, and finalize. The lens is long-running execution (multi-minute agent runs), recovery from partial failures, and operational safety (no silent state mutation, RequiresYes gating, stable run identity).

## Proposed approach
- TUI receives an injected action runner (never imports internal/app or shells directly); runner handles ExecuteAction semantics while TUI only triggers and observes.
- Before execution, re-validate the action via runplan.Plan using fresh runstate + canonical files; reject if the artifact now exists or consensus state has advanced.
- For RequiresYes=true actions, surface an explicit confirmation modal in the TUI that quotes the action Summary and Risk.
- Launch via the existing runner but under TUI-local supervision: track the child PID, surface live stdout tail in a dedicated pane, and on exit re-invoke refreshRuns + planner to update the selected run.
- Recovery: on TUI restart or focus switch, any in-flight TUI-owned run is treated as "inspect" until its events.jsonl shows terminal state; never auto-retry without a fresh planner action.
- Safety: all mutations still go through canonical files only; TUI never edits participant .md files directly.

## Concerns / open questions
- How to bound concurrent TUI-launched runs without a daemon (current continuous-run-tui non-goal)?
- Should long-running output be persisted to a per-action log under runs/<run-id>/ or only shown live?
- Does the TUI need to surface "cancel" for a running agent, and if so how does that interact with the protocol's non-silent advancement rule?

## Risks
- Race between TUI action dispatch and another agent (or CLI) writing the target artifact between Plan and execute.
- User closes TUI while an agent is still writing its round artifact, leaving the run in an ambiguous "pending" state until manual inspection.
- Confirmation fatigue if too many normal-risk actions require explicit yes; may encourage blind acceptance.