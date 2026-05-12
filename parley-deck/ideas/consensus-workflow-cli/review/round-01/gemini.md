---
agent: gemini
idea: consensus-workflow-cli
review-round: 1
date: 2026-05-12
implementation-pr: https://github.com/feci/parley-deck-cli/pull/10
---

## Summary

The implementation provides the core deterministic consensus workflow as specified in `FINAL.md`. However, there are significant issues in the triage logic that could lead to safety violations (finalizing blocked ideas) and several protocol deviations regarding file naming and template structure.

## Findings

### [CRITICAL] Triage logic bypass via non-canonical status strings
In `internal/consensus/consensus.go`, `validateDocument` parses signoff statuses but fails to use the normalized canonical value for triage state transitions.

```go
// internal/consensus/consensus.go:235
for _, signoff := range doc.Signoffs {
    // ...
    if _, err := CanonicalStatus(signoff.Status); err != nil {
        summary.Errors = append(summary.Errors, fmt.Sprintf("line %d: %v", signoff.Line, err))
        continue
    }
    // BUG: using raw signoff.Status instead of normalized value from CanonicalStatus
    switch signoff.Status {
    case StatusReservations:
        hasReservations = true
    case StatusBlock:
        hasBlock = true
    }
}
```

If a participant manually writes `Status: block`, `CanonicalStatus` will validate it, but `hasBlock` will remain `false` because `"block" != "❌ BLOCK"`. If all other signoffs are `accept`, the triage will result in `ready`, allowing `parley consensus finalize` to succeed despite an active block.

**Fix:** Assign the result of `CanonicalStatus` to a variable and use it in the `switch` and for triage.

### [MAJOR] Incorrect naming of aborted consensus files
In `internal/consensus/consensus.go`, `nextAbortedPath` uses a simple counter `i` to generate the round number in the filename, which is misleading.

```go
func nextAbortedPath(baseDir string) string {
	for i := 1; ; i++ {
		path := filepath.Join(baseDir, fmt.Sprintf("round-%02d-consensus-aborted.md", i))
```

If consensus is blocked in `round-03`, it will be saved as `round-01-consensus-aborted.md` (if it's the first abortion). The filename should reflect the actual round being aborted and use a separate attempt counter if needed (e.g., `round-03-consensus-aborted-1.md`).

### [MAJOR] Strict signoff line parsing
`parseDocument` uses `strings.CutPrefix(line, "Status:")` etc. This is very strict and will include trailing text if the file contains comments like the one suggested in `COOPERATION.md` example: `Status: ✅ ACCEPT (or 🟡 ACCEPT-WITH-RESERVATIONS, or ❌ BLOCK)`. While `CanonicalStatus` tries to handle this, the aforementioned triage bug (using the raw string) makes this particularly fragile.

### [MINOR] Review consensus template deviations
The `Draft` template for review consensus (in `draftTemplate`) has two deviations from `COOPERATION.md` and `FINAL.md`:
1. It uses `cycle: review` instead of a numeric cycle (e.g., `cycle: 1`).
2. It is missing the `reviewed-commit` field in the frontmatter, which is a protocol requirement for review consensus.

### [MINOR] Missing non-goals section in FINAL.md scaffold
The `finalTemplate` in `internal/consensus/consensus.go` does not include the `### Non-goals` section, which was explicitly requested in the `FINAL.md` specification for the `FINAL.md` scaffold.

### [NIT] Reopen does not update implementation status
When `parley consensus reopen --review` is called, it renames the review consensus file and updates the idea status in `00-prompt.md` to `discussion` (if not review). However, it should also ideally update the `IMPLEMENTATION.md` status to reflect that the consensus attempt failed and it's back to a fix-up or discussion state.

## Tests reviewed

- `internal/consensus/consensus_test.go`:
  - `TestDraftAndAppendSignoffTriage`: Covers basic triage flow. Note: it uses `Status: accept` and `Status: reserve` which happen to work because `CanonicalStatus` is called in `AppendSignoff` before writing, but doesn't catch the manual edit bug in `validateDocument`.
  - `TestMalformedSignoffs`: Verifies error reporting for unknown participants and missing counter-proposals.
  - `TestFinalizeCreatesFinalAndUpdatesStatus`: Verifies the finalization transition.
  - `TestReservedFinalizeRequiresOpenItems`: Verifies the safety check for reserved consensus.
  - `TestReopenBlockedConsensus`: Verifies the file renaming and status rollback.
  - `TestReviewDraftUsesReviewPath`: Verifies path selection for `--review`.

## Verdict

REQUEST_CHANGES
