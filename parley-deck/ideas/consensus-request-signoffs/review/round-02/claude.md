---
agent: claude
idea: consensus-request-signoffs
review-round: 2
date: 2026-05-13
reviewed-commit: 7cb518507fc8f8389eee0d014162334e5ed20e71
---

## Findings

- `NIT` internal/app/consensus_request_signoffs.go:121-131 — `printPartialProgress` is now correctly called on the four post-invocation failure paths (statusErr, afterRaw read, validateErr, runErr-after-valid-append), but the three pre-invocation error paths inside the loop (`consensus.Status` failure, `before.Errors` non-empty, `os.ReadFile(before.Path)` failure) still return without printing accumulated successes. These only matter for batches of two or more participants where iteration N+1 fails before invoking its agent, but the FINAL.md rule ("print which signoffs succeeded before failure") is symmetric across mid-loop failures. Cheap to fix by adding `printPartialProgress(stdout, successes)` to those three returns.
- `NIT` internal/app/consensus_request_signoffs.go:345-375 — `signoffHeaderAgents` and `newSignoffAgents` use a local `var agents []string` that shadows the `parley-deck-cli/internal/agents` package import. Compiles fine and isn't referenced as a package inside either function, but reads oddly next to other helpers in this file that use `agents.*`. Rename to `ids` or similar.

## Open questions

- The canonical block in the prompt template now always shows `Counter-proposal (required if ❌): <counter-proposal if status is ❌ BLOCK>` regardless of the status the agent will return. The prose qualifier "required if ❌" is clear to a human, but a literal-minded agent might emit the placeholder line on an ACCEPT signoff. Worth confirming the `consensus` parser tolerates an extra `Counter-proposal:` line on non-BLOCK statuses, or alternatively hoist the line into a status-conditional hint above the canonical block. The happy-path test in `TestConsensusRequestSignoffsHappyPath` does not exercise this.

## Summary

All seven agreed fixes from `review/consensus.md` are present and behave as specified at the inspected commit. The "exactly one new signoff" enforcement is now defended from two angles — the raw-text prefix + header-count check in `validateAppendOnlyContent` (lines 327-343) and the parsed `newSignoffAgents` check in `validateRequestedSignoff` (line 298) — and the two new regression tests (`TestConsensusRequestSignoffsRejectsForgedExtraSignoff`, `TestConsensusRequestSignoffsRejectsExistingContentEdit`) pin both the forged-foreign-signoff and the existing-content-edit scenarios I flagged in round 1. The preflight exit-code mapping, prompt prose quoting, BLOCK template line, and partial-progress reporting on the most likely failure paths are all wired correctly. Residual risks: pre-invocation mid-loop failures still skip the partial-progress summary (NIT above), the parser's treatment of a stray `Counter-proposal:` line on ACCEPT is untested, and crash/timeout mid-append leaving the file truncated is still uncovered by tests. No regressions spotted in the rest of the diff.
