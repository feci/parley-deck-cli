---
agent: codex-1
idea: protocol-restructure-appendices
review-round: 1
date: 2026-07-03
reviewed-commit: bfaea9b
---

## Summary

The pure-move claim survived the sorted-line diff checks for both protocol copies, and the top-level order matches FINAL.md: core sections first, then §10, §9, §11, Appendix A, and §12-§14. I did not find positional prose broken by the move.

I did refute two acceptance/evidence claims: the full suite is not green in this checkout, and a header-based cross-reference audit is not zero-dangling because `§6.6` has no matching markdown header. The `§6.6` issue appears pre-existing rather than introduced by the move, but the current implementation evidence still overclaims the audit result.

## Refutation attempts (what you tried to break, and the result)

- Read `FINAL.md`, `IMPLEMENTATION.md`, `review/round-01/DIFF.txt`, `parley-deck/COOPERATION.md`, `internal/protocol/defaults/COOPERATION.md`, and the idea prompt.
- Pure move: ran the requested sorted-line diffs against `main` for both `parley-deck/COOPERATION.md` and `internal/protocol/defaults/COOPERATION.md`. Both commands exited 0 with empty output, so I found no added, removed, or altered lines.
- Diff inspection: `DIFF.txt` shows only the `§10` block moving before `§9` in both copies. That is equivalent to the required final state where `§9` follows `§10`; I did not find a banner or rule-text edit.
- New order: `grep -nE '^## '` shows `§0`-`§8`, then `§10`, then `§9`, then `§11`, Appendix A, and `§12`-`§14` in both live files.
- Drift guard / byte identity: `go test ./internal/protocol/...` passed.
- Cross-refs: scanned markdown headings and all `§N` / `Appendix` references in both live files. All section/header references resolve except `§6.6`, which appears in the English-only invariant but has no `### 6.6` or equivalent heading. Because the sorted diff is empty, this was preserved, not created by the move.
- Positional prose: searched for `above`, `below`, `earlier`, `later`, `table above`, `table below`, `see above`, and `see below`. The hits are temporal wording or local pointers inside §4/§5/§11 templates and transport subsections; I did not find one falsified by moving §9.
- Full suite: `go test ./...` failed in `internal/runner`. A focused rerun of `go test ./internal/runner -run TestDurableKillEndToEndRealProcess -count=1 -v` failed the same way.

## Findings

### [MAJOR] Full suite is not green

`go test ./...` exits 1 at reviewed commit `bfaea9b`. The failing test is `TestDurableKillEndToEndRealProcess` in `internal/runner`, with:

```text
process verification failed (no recorded boot id); not killed
```

A focused rerun of that test also fails. The implementation diff does not touch runner code, so this is likely unrelated to the protocol reorder itself, but it still refutes FINAL.md acceptance criterion 5 and the IMPLEMENTATION.md claim that the full suite is green.

Fix: make `internal/runner` pass in this environment or record an explicit accepted exception for this pre-existing/environment-specific failure, then rerun `go test ./...` and update `IMPLEMENTATION.md` with accurate validation evidence.

### [MINOR] Header-based cross-reference audit is not zero-dangling

Both protocol copies contain `English-only rule (§6.6)`, but there is no markdown header for `§6.6`; the existing header is `## 6. Conflict-avoidance mechanics`, and English-only is numbered list item 6 under that section. Under the prompt's "resolve to an existing header" rule, this is dangling in both files.

This is not move-created: the sorted-line diff against `main` is empty, so the reference was preserved. It does, however, contradict the implementation's "0 dangling" cross-reference evidence if the audit is header-based.

Fix: either change the reference to `§6 rule 6` / `§6` through a separate content-changing protocol fix, or explicitly define the cross-reference audit as accepting numbered list-item anchors and update the evidence accordingly.

## Signoff

Status: ❌ BLOCK
