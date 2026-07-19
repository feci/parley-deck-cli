---
agent: codex-1
idea: composite-agent-naming-and-roster-reinit
review-round: 4
date: 2026-07-19
reviewed-commit: dd11f5b4b13914641ea4108a3be6bfddddd8a801
responding-to: [codex-1/review/round-03]
---

## Summary

All three MAJOR correctness findings left open in round 03 are **RESOLVED** at HEAD. The
default path now hard-stops when a parsed non-empty §2 roster has no active resolvable
member, excludes inactive rows, and retains the raw-family fallback only when
`ReadRosterIDs` reports no readable roster. Machine initialization now admits only built-in
and central families for both proposals and target validation.

I found no new cycle-5 correctness regression. Two non-blocking review-hygiene items remain:
the two default-run fixes lack the regression tests explicitly requested in round 03, and
`IMPLEMENTATION.md` still lacks the Phase-8 cycle metadata. Neither is a CRITICAL/MAJOR
runtime defect.

## Refutation attempts

- Read `IMPLEMENTATION.md` cycles 1–5 and its deferred list, the supplied `cdef991..HEAD`
  fix-up diff, the focused `5760a476..dd11f5b` cycle-5 diff, round 03, and the HEAD call sites.
- Traced no-flag selection with: no readable roster; one active mapped/installed member; a
  readable roster whose active members are all unmapped/unavailable; one active plus one
  inactive mapped member; and an all-inactive roster.
- Traced machine initialization with built-in, central custom, deck-only, local-only, and
  environment-only family catalogs, including validation of adapters already present in the
  central target. Traced session initialization separately with its intentionally layered
  catalog.
- `go build ./...`, `go vet ./...`, formatting of all cycle-5 Go files, and
  `git diff --check 5760a476..HEAD` pass.
- `go test ./internal/protocol ./internal/config ./internal/app ./internal/agents` passes.
  `go test ./...` passes every package except the previously documented environment-dependent
  `internal/runner.TestDurableKillEndToEndRealProcess`: this host has no recorded boot ID, so
  the safety check refuses to kill the attributed process. Cycle 5 does not touch that path.

## Round-03 MAJOR dispositions

### [MAJOR] RESOLVED — Default-run raw-family escape hatch

`internal/app/app.go:1818-1836` enters the roster-default branch whenever
`protocol.ReadRosterIDs` succeeds, collects only IDs that resolve through the effective
mapping, and returns exit 1 when that set is empty. It therefore cannot reach
`selectedParticipantIDs` with an empty flag in the zero-resolved case. The legacy fallback
at `internal/app/app.go:2419-2422` is reached only when the roster reader did not succeed;
`internal/protocol/roster.go:22-25,69` defines that success boundary.

### [MAJOR] RESOLVED — Default selection launched inactive members

`internal/app/app.go:1819-1828` now retains the inactive subset returned by
`ReadRosterIDs` and skips each inactive ID before resolution or insertion into the default
participant flag. An all-inactive roster consequently reaches the same hard stop at
`app.go:1831-1833`; it does not fall back to installed families.

### [MAJOR] RESOLVED — Machine scope consumed project-only families

`internal/config/runtime.go:248-276` builds `MachineFamilyCatalog` from
`agents.DefaultSpecs()` plus only the central `agents.toml` `[agents.*]` keys. Machine init
loads that catalog and passes it into proposal resolution at
`internal/app/roster.go:216-228`; `resolveRoster` rejects any family outside it at
`roster.go:135-140`. The same catalog validates every adapter already present in the target
at `roster.go:241-258`. Session scope deliberately leaves `allowed` nil and retains layered
behavior. Thus deck/local/environment-only families are neither proposed nor accepted for
the machine target.

## New cycle-5 regressions

None found.

## Remaining non-blocking findings

### [MINOR] The two default-run fixes have no direct regression tests

Round 03 explicitly required a no-flag zero-resolved test and an active-plus-inactive
default-selection test. Cycle 5 changes `internal/app/app.go` but adds tests only for the
machine family filter/catalog (`internal/app/roster_test.go:41-76`); there is no test that
drives or extracts the new `runTask` default-selection branch. The code is correct on review,
but these exact quorum-boundary failures can regress silently.

**Suggested fix:** extract the roster-default selection into an error-returning helper or
exercise `runTask` with controlled discovery, then lock the zero-resolved, inactive-filter,
all-inactive, and no-readable-roster fallback cases.

### [MINOR] Phase-8 metadata remains stale

`IMPLEMENTATION.md:1-7` still says `status: implemented` and has no top-level `completed` or
`head-commit`, while cycle 5 is recorded only as prose at lines 146-160. This leaves resume
and audit consumers unable to identify the reviewed fix-up HEAD from the required metadata.

**Suggested fix:** record `status: fix-up-cycle-5`, `completed: 2026-07-19`, and
`head-commit: dd11f5b4b13914641ea4108a3be6bfddddd8a801` in the frontmatter plus the cycle-5
status block. Change to `status: complete` only after review consensus records zero agreed
fixes.

## Documented deferrals

The `roster diff` surface, TOML `autonomous_write` override, OS-enforced sandbox/live
sentinel, S6 display plumbing, and fully interactive model/effort picker are accepted scope
deferrals and are not findings in this review.

## Verdict

**ACCEPT-WITH-FIXES**

The three release-blocking MAJORs are closed and no new blocker was found. The remaining
work is limited to regression coverage and protocol metadata.
