---
agent: gemini
idea: agent-runtime-config
review-round: 1
date: 2026-05-11
reviewed-commit: d38e7f738471caa61042ca087c445b4df5bd2606
---

## Summary
The implementation successfully delivers the `agent-runtime-config` slice with proper configuration layering, source tracking, and expanded agent spec fields. It introduces the `parley agents` command group and correctly threads resolved runtime values into the runner and agent prompts.

## Findings
### [CRITICAL]
None.

### [MAJOR] Codex Git smoke runs directly in CLI instead of via agent
In `internal/app/app.go`, the `runCodexGitSmoke` function (lines 280–303) executes Git commands directly using `exec.CommandContext(ctx, "git", ...)`. According to `FINAL.md`, these checks are specifically intended to verify Codex's `workspace-write` sandbox and `on-failure` approval policy. By running them in the CLI's own process environment, the implementation fails to verify the agent's actual runtime behavior and sandbox restrictions, providing false confidence in the agent's capability to participate in a round.

### [MINOR] Agents list hardcodes "not-probed" for headless status
The `PrintRuntimeMatrix` function in `internal/agents/discover.go` (lines 183–223) hardcodes the `headless` status variable to `"not-probed"`. `FINAL.md` requested "clear columns for installed, version, and headless status" in `parley agents list`. Always showing `"not-probed"`—even after a successful `verify --full` has been performed—makes this column unhelpful. At minimum, it should indicate if a headless command is configured or if the agent has passed a recent probe.

### [NIT] PrintDiscovery vs PrintRuntimeMatrix inconsistency
The legacy-style `PrintDiscovery` function still displays the actual `HeadlessMode` command string, whereas the new `PrintRuntimeMatrix` (used by `agents list`) replaces it with the "not-probed" status. This results in a loss of visibility for users who want to inspect the exact command being used for headless execution.

### [NIT] Default values for non-proven fields
While `internal/agents/discover.go` correctly uses `CLIDefault` for many fields, `ExternalBackend` defaults to `ExternalUnknown` ("unknown") and `Speed` defaults to `DefaultSpeed` ("balanced"). While reasonable, this slightly diverges from the goal of using `cli-default` for all fields that cannot be explicitly proven or discovered from the environment.

## Open questions
- Should the results of `parley agents verify --full` be persisted (e.g., in a local state file) so that `parley agents list` can display a meaningful headless status?
- Is there a plan to implement Step 1 of the resolution order ("Explicit CLI flag where the command already has one") for runtime fields like model or timeout in a follow-up slice?
