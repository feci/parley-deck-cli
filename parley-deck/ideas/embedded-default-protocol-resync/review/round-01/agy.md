---
agent: agy
idea: embedded-default-protocol-resync
review-round: 1
date: 2026-06-13
reviewed-commit: efe76d0
---

## Summary

The implementation is very clean, runs successfully, passes all tests, and strictly adheres to the decisions outlined in [FINAL.md](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/parley-deck/ideas/embedded-default-protocol-resync/FINAL.md) and [IMPLEMENTATION.md](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/parley-deck/ideas/embedded-default-protocol-resync/IMPLEMENTATION.md). Specifically:
- §12 has been successfully appended to the embedded default verbatim.
- The embedded default's headers have been genericized and its tables emptied of active agents.
- The drift guard test is fail-closed, verifying exact matches on structural anchors.
- The init-output test successfully validates the behavior of `defaultCooperationForInit()`.

A minor adversarial blind spot and a small lax occurrence check are flagged below as a NIT/MINOR finding, but they do not block acceptance.

## Findings

### [NIT] Lax `**Protocol synced:**` occurrence checks in drift guard
* **What is wrong:** The normalizer strips the `**Protocol synced:**` line from both files if it starts with the prefix `**Protocol synced:**`, but the test does not assert that this line is absent in the embedded default or present in the live deck. If a developer accidentally adds a `**Protocol synced:**` line to the embedded default, the drift guard will silently strip it and pass.
* **Why it matters:** Decision D2 states that the embedded default must not contain a `**Protocol synced:**` line. While `TestDefaultCooperationForInit` catches this indirectly by asserting the init output lacks it, verifying the occurrence counts directly in the drift guard makes it robust and fail-closed.
* **Concrete suggested fix:** In `TestEmbeddedDefaultMatchesLiveDeck` in [drift_test.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/protocol/drift_test.go#L35-L72), assert that the prefix `**Protocol synced:**` appears exactly 0 times in `defaultCooperation` and exactly 1 time in the live `deck` file:
  ```go
  nEmb := 0
  for _, line := range strings.Split(defaultCooperation, "\n") {
  	if strings.HasPrefix(line, protocolSyncPrefix) {
  		nEmb++
  	}
  }
  if nEmb != 0 {
  	t.Fatalf("embedded default: header line %q* appears %d times, want exactly 0", protocolSyncPrefix, nEmb)
  }

  nDeck := 0
  for _, line := range strings.Split(deck, "\n") {
  	if strings.HasPrefix(line, protocolSyncPrefix) {
  		nDeck++
  	}
  }
  if nDeck != 1 {
  	t.Fatalf("live deck: header line %q* appears %d times, want exactly 1 (drift guard fails closed)", protocolSyncPrefix, nDeck)
  }
  ```

### [MINOR] Key adversarial zone assessment
* **What is wrong:** Normalization of allowlisted zones is prefix-based for header values (`**Workspace:**` and `**Created:**`) and discards the entire body of the §2 tables. Under this model, someone could place arbitrary text inside a §2 table body row starting with `|` or on a header line, and it would be silently normalization-collapsed or ignored without flagging drift.
* **Why it matters:** This creates a theoretical blind spot. However, because:
  1. These zones are inherently project-specific (the template must be empty/generic, while the live deck is populated), they *must* diverge.
  2. The roster table body contains only row data mapping Agent IDs to directories and roles, which is not evaluated as protocol logic.
  3. Structural changes to table columns are caught because the header and separator rows are matched exactly (`line == rosterHeaderLine` / `line == handleHeaderLine`).
  Thus, this behavior is acceptable.
* **Concrete suggested fix:** Keep the current design but optionally check that ignored table body rows have the correct number of delimiters (e.g. exactly 4 `|` characters for the roster table rows) to prevent arbitrary formatted text from being nested inside.

## Open questions

None.

## Verdict

ACCEPT-WITH-FIXES
