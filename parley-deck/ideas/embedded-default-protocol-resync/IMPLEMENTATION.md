---
idea: embedded-default-protocol-resync
status: fix-up-cycle-1
implementer: claude
started: 2026-06-13
completed: 2026-06-13
branch: parley-deck-cli#feature/embedded-default-resync
head-commit: 9fb1d42
reviewed-commit-cycle-1: efe76d0
design-pr: n/a
implementation-pr: n/a
---

## Summary of work

Resynced the embedded default protocol (`internal/protocol/defaults/COOPERATION.md`,
the `parley init` bootstrap template) with the live deck and added a fail-closed
Go drift guard so the two copies cannot silently diverge again. Per FINAL.md /
consensus D1–D7.

- **§12 carried verbatim (D1).** Appended `## 12. Pipeline blocks & action stages`
  to the embedded default after §11 by copying the live deck's `sed -n '805,$p'`
  tail, so the section (incl. its "ratified by idea
  `meta-protocol-change-end-to-end-pipeline` (2026-06-02)" provenance line and the
  exact final newline) is byte-identical. Verified by `diff` of the two §12
  sections (empty).
- **Header genericized (D2).** Embedded header now reads
  `**Workspace:** \`<workspace-name>\`` and
  `**Created:** \`<date> — created by parley init\``; `**Transport:** \`github-pr\``
  retained; no `**Protocol synced:**` line added.
- **§2 tables emptied (D3).** Both §2 tables (roster + host-handle) keep their
  header and separator rows; the parley-deck agent body rows were removed.
- **InitWorkspace unchanged (D4).** `defaultCooperationForInit()` still performs
  only the `github-pr` → `local-dir` swap; no renderer, no clock seam.
- **Two Go tests (D5, D6)** in `internal/protocol/drift_test.go`:
  - `TestEmbeddedDefaultMatchesLiveDeck` — the drift guard. Normalizes exactly the
    five allowlisted zones (deck-only `**Protocol synced:**` line, `**Workspace:**`
    value, `**Created:**` value, §2 roster table body, §2 host-handle table body)
    and compares the rest line-for-line. Fails closed: a missing deck file, or a
    missing/duplicated structural anchor, is a hard failure rather than a silent
    skip or a widened normalization region.
  - `TestDefaultCooperationForInit` — asserts bootstrap output has
    `Transport: local-dir`, contains §12 incl. provenance, contains the static
    placeholders, contains no `**Protocol synced:**` line, and contains none of
    the parley-deck roster names.

## Implementation plan / checklist

- [x] Append §12 verbatim to the embedded default (byte-identity verified).
- [x] Genericize the two embedded header values; keep transport; omit Protocol-synced.
- [x] Empty both §2 table bodies (keep header + separator).
- [x] Add `internal/protocol/drift_test.go` (guard + init-output tests).
- [x] `gofmt` clean; `go build ./...` green.
- [x] Negative control: injected drift makes the guard FAIL with a clear message;
      reverting makes it PASS.
- [x] Full `go test ./...` green.

## The dual-copy invariant (D7 — where it lives instead of §7)

Per consensus D7 (codex's scope objection, concurred 4/4), the live deck's §7 was
**not** edited in this idea. The "any protocol edit must touch BOTH
`internal/protocol/defaults/COOPERATION.md` and `parley-deck/COOPERATION.md`"
invariant is enforced by `TestEmbeddedDefaultMatchesLiveDeck` and documented in:
1. the package-level comment block at the top of `drift_test.go`, and
2. that test's failure message, which names both files and the allowlist.
A §7 protocol-text pointer remains available as a future meta-protocol-change.

## Fix-up cycle 1 (review/round-01 → review/consensus.md)

Both agreed fixes applied:

1. **[MINOR] Drift guard now asserts the embedded D2/D3 invariants** (closes the
   in-zone-edit blind spot codex + agy converged on). `drift_test.go` gained
   `assertEmbeddedBootstrapShape` (the three genericized header lines verbatim,
   via `embWorkspaceLine`/`embCreatedLine`/`embTransportLine`) and
   `assertEmptyTableBody` (each §2 table has its separator and an empty body),
   plus explicit `**Protocol synced:**` occurrence checks (0 in the embedded
   default, exactly 1 in the live deck). These run before normalization, so an
   illustrative roster row, an altered placeholder value, or a stray sync line in
   the embedded copy now fails the guard instead of being normalized away.
   Negative controls confirmed: a stray Protocol-synced line fails the D2 check;
   an `agent-1` roster row fails the D3 empty-body check; the clean tree passes.
   (agy's optional body-row delimiter check is subsumed — the embedded bodies are
   asserted empty; the deck's rows are project-specific data, not protocol logic.)
2. **[NIT] head-commit corrected** from the stale `bc0af15` to the fix-up cycle 1
   commit.

`gofmt` clean; `go build ./...` and full `go test ./...` green.

## Deviations from FINAL.md

- None.

## Notes for reviewers

- The drift guard reads `../../parley-deck/COOPERATION.md` relative to the package
  dir (where `go test` runs). It is in-repo-only by design; the out-of-repo
  packaged skill reference is a separate, deferred concern.
- Allowlist breadth is the key review axis: confirm the five zones are exactly
  right and that the anchored normalization cannot mask a real protocol change
  (the negative control exercised the length-mismatch path; consider whether an
  equal-length in-zone edit could slip through).
