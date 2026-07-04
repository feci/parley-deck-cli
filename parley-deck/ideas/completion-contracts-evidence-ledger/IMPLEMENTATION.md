---
idea: completion-contracts-evidence-ledger
status: fix-up-cycle-1
implementer: claude-1
started: 2026-07-04
completed: 2026-07-04
branch: parley-deck-cli#completion-contracts-design
head-commit: b443481
design-pr: https://github.com/feci/parley-deck-cli/pull/67
implementation-pr: same
---

## Summary of work

Extended `checks:` from a scalar to an optional named list = the completion contract.
The driver runs each criterion, writes a per-criterion evidence table into the existing
`## Validation evidence` section of IMPLEMENTATION.md, and fails closed on any non-zero
exit (delivering the Phase-8 veto through the existing RunChecks gate). No new
`done_when:` key, no separate ledger artifact. Scalar/absent `checks:` is unchanged.

## Implementation plan / checklist

- [x] `internal/driver/checks.go` (new): `ReadChecksContract(ideaDir)` (yaml.v3 over the
      frontmatter block) — absent/scalar → legacy; list → validated `[]{name, command}`;
      malformed list (empty/dup name, empty command, empty list) → fail closed.
- [x] `internal/app/driver_checks.go` (new): `runChecksContract` runs each criterion
      (`sh -c`, cwd = repo root), captures exit + duration + secret-scrubbed/truncated
      (~100 lines / 4KB) output; `writeValidationEvidence` overwrites the
      `## Validation evidence` section each cycle (`replaceSection`); returns fail-closed
      with a descriptive detail naming the failing criterion.
- [x] `internal/app/driver_impl.go`: `RunChecks` dispatches to the contract path when
      `checks:` is a list, before the scalar/go-test fallback. The Phase-8 veto is
      delivered here (RunChecks is already the pre-review / post-fixup gate).
- [x] Protocol text (BOTH `parley-deck/COOPERATION.md` and
      `internal/protocol/defaults/COOPERATION.md`, byte-identical): LE-4 (scalar or list;
      list activates the contract + driver writes `## Validation evidence`), Phase-5
      IMPLEMENTATION.md template (driver-populated when list), Phase-8 (list-checks
      completion requires latest all-pass; fail-closed veto independent of strict_gate).
- [x] Skill fallback `parley-deck-skill/references/COOPERATION.md` re-synced
      (body-identical from line 7; only header lines 5-6 neutralized).
- [x] `meta/version.json` protocolSha256 refreshed; `meta/protocol-changelog.md` entry.
- [x] Tests: `internal/driver/checks_test.go` (absent/scalar-legacy/list/malformed-fail-
      closed), `internal/app/driver_checks_test.go` (scrub, replaceSection, evidence write
      + veto). Drift guard `TestEmbeddedDefaultMatchesLiveDeck` green.
- [x] Checks: `go build ./...`, `go vet`, `gofmt -l` clean; targeted suites green.

## Deviations from FINAL.md

- The Phase-8 veto is delivered through the existing `RunChecks` gate (which already runs
  pre-review and post-fixup) rather than a separate new completion-path branch in
  impl.go — same fail-closed outcome (a failing contract blocks the fix-up from
  converging to complete), smaller surface, no new call site. Documented.

## Notes for reviewers

- `ReadChecksContract` distinguishes scalar vs list via `yaml.Node.Kind`; confirm an
  unparseable frontmatter safely falls back to legacy (scalar path) rather than erroring.
- `scrubAndTruncate` redacts token/secret/password/api-key/bearer/authorization patterns
  and caps to the last 100 lines / 4KB; confirm no raw unbounded output can reach the
  ledger.
- `replaceSection` overwrites `## Validation evidence` up to the next `## `; confirm it
  preserves following sections (`TestReplaceSection`).

## Fix-up cycle 1
status: complete
completed: 2026-07-04

### Fixes applied (review round-01)
- [CRITICAL, codex-1] The zero-agreed-fixes completion path was un-gated. Added a
  completion-contract gate in `impl.go` (before GoalCheck/Complete): with list-form
  `checks:`, completion runs a fresh RunChecks and vetoes (ActionEscalated) on failure;
  scalar/absent is unchanged. Test `TestPhaseReviewListChecksVetoCompletion`.
- [MAJOR, hermes-1] Evidence write dirtied the tree and tripped the next cycle's
  `gitTreeClean` guard. The driver now COMMITS the IMPLEMENTATION.md evidence write
  immediately (`commitEvidence`, mirroring driver-artifact commits) — best-effort,
  no-op-safe, non-git-safe.
- [MAJOR, codex-1] Hardened the secret scrubber: full `Authorization: Bearer <token>`,
  standalone `bearer <token>`, labeled key/value, and standalone provider shapes
  (`sk-`, `gh[pousr]_`, `xox[baprs]-`, `AKIA…`, JWT). Vacuous test replaced with a
  multi-shape leak assertion (also fixes hermes-1 MAJOR-2).
- [MAJOR, codex-1] A YAML syntax error in list-form `checks:` now FAILS CLOSED
  (`looksLikeChecksList` detects block form) instead of silently falling back to the
  legacy scalar/go-test path. Test `TestReadChecksContractSyntaxErrorListFailsClosed`.

### Deviations from agreed fixes
None.
