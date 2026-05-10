---
from: gemini
to: all
idea: parley-deck-cli-plan
phase: runner-fixup-advisory-review
blocking: no
date: 2026-05-10
---

## Summary
The runner implementation in fix-up cycle 2 successfully establishes the parallel execution core and protocol artifact validation. The use of a fake-agent helper for unit testing the execution lifecycle is excellent. Gemini-specific isolation via temporary `GEMINI_CLI_HOME` addresses known auth-hang issues. The primary areas for improvement are redundant discovery overhead, lack of participant selection, and minor prompt injection duplicates.

## Findings
### CRITICAL
None.

### MAJOR
1. **Redundant Agent Discovery**: `parley run` performs agent discovery (including version probes) at the start of the task, and `runTUIView` performs it again immediately after execution finishes. Since each probe has a 4s timeout, this can introduce significant UI lag, especially if multiple agents are installed but non-responsive. Discovered results should be passed into the TUI or cached.
2. **Implicit Participation**: `parley run` automatically selects all discovered agents as participants for the new idea. While acceptable for a bootstrap slice, this lacks the "Selects participants" control promised in `FINAL.md` and may lead to unexpected costs or unwanted model involvement if the user has many agents configured.
3. **Prompt Content Redundancy**: `BuildRoundOnePrompt` injects the `task` string into the `Task:` section and then appends the full `00-prompt.md` (which already contains the task under `## Problem / idea`). This consumes unnecessary context tokens and may cause confusion for some models.

### MINOR
1. **Hardcoded Protocol Logic**: The Round 1 prompt structure and "one agent, one file" rules are hardcoded in `internal/runner/runner.go`. To support future rounds and custom protocols without recompiling, this logic should eventually move to a template or the `protocol` layer.
2. **Gemini Isolation Scope**: `isolatedGeminiHome` only copies `oauth_creds.json` and `google_accounts.json`. Users with custom credential filenames or specific `.gemini/settings.json` requirements (other than the hardcoded `oauth-personal` override) may find the isolation too restrictive.
3. **Unchecked Store Errors**: `runner.RunRoundOne` ignores return values from `opts.Store.Append`. While currently non-blocking, failing to record `round.completed` or `agent.started` events should at least be logged or surfaced.

### NIT
1. `NewRunID` nanosecond resolution produces very long directory names; millisecond or second resolution is usually sufficient for CLI tools.
2. `printRunResults` uses fixed-width padding (`%-8s`) which may truncate longer Agent IDs in the future.

## Verification notes
- Parallel agent execution via goroutines and `sync.WaitGroup` is correctly implemented and waited upon.
- The `TestRunRoundOneCreatesArtifactWithHeadlessAgent` successfully validates the execution contract using a real subprocess helper.
- Artifact existence checks before execution correctly prevent accidental overwrites.

## Recommendation
Proceed with the current slice but prioritize passing discovered agent data into the TUI to eliminate redundant probes. Add a mechanism (flag or interactive prompt) for participant selection before moving to M4. Consider refactoring the prompt builder to avoid double-task injection.
