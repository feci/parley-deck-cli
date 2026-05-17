---
agent: gemini
idea: roadmap-implementation-plan
review_round: 1
date: 2026-05-17
implementation-pr: https://github.com/feci/parley-deck-cli/pull/20
verdict: COMMENT
---

## Summary

The PR successfully implements the `round-index-artifact-pruning` slice. It introduces deterministic `_index.md` generation at the end of a round and a context-only sanitizer for hidden-reasoning fences (`<think>`, `<thought>`, `<thinking>`). The implementation preserves the integrity of participant artifacts while providing a useful, runner-owned summary for downstream context packing.

## Findings

- **MINOR**: The `removeTaggedBlocks` logic in `internal/runner/round_index.go` is "greedy" when encountering an unclosed tag—it drops everything from the open tag to the end of the file. While this is safe from a "hidden reasoning" perspective, it results in a complete loss of summary context for the index if an agent provides a malformed artifact.
- **MINOR**: `extractH2Sections` loads and splits the entire artifact in memory. Given the typical size of agent artifacts, this is currently acceptable, but should be monitored if artifacts grow significantly.
- **NIT**: `approxTokens` uses the simple `(bytes + 3) / 4` heuristic. This is deterministic and follows the design, but remains a very rough approximation for non-ASCII content.
- **NIT**: `escapeTable` replaces newlines and escapes pipes, which is sufficient for current Agent IDs and Statuses, but may need more robust escaping if participants use complex names.

## Tests / verification reviewed

- **Sanitization Tests**: Verified that all three supported tags are removed and that malformed tags are handled safely (by dropping trailing content).
- **Determinism**: Confirmed through tests that `BuildRoundIndex` produces identical output for identical inputs across multiple calls.
- **Round Integration**: Verified via `internal/runner/runner_test.go` that the indexer is correctly wired into `RunRoundOne` and that its failure produces a `Result.Warning` rather than an `ExitError`.
- **Integrity**: Confirmed that source artifacts remain unmodified after indexing.

## Risks / open questions

- **Secret Redaction**: As documented, this sanitizer is for hidden reasoning only. There is a risk that users might mistake it for a general-purpose PII/secret redaction tool. The explicit documentation in `_index.md` is a good mitigation.
- **Resource Usage**: Large artifacts will be processed entirely in memory. If future roadmap slices involve processing many large artifacts simultaneously, a streaming approach might be needed.
