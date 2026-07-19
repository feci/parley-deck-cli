---
agent: codex-1
idea: composite-agent-naming-and-roster-reinit
review-round: 3
date: 2026-07-19
reviewed-commit: 5760a476d7f94dde5f2f88c2e7665dce46d06d0b
responding-to: [codex-1/review/round-02]
---

## Summary

Cycle 4 resolves the autonomous-confinement representation, candidate-TOML validation,
agy-only tier rendering, and the stale `selectedAgents` contract comment. The default-run
happy path also now uses roster IDs when at least one active row resolves.

The branch is still not safe to accept. The no-flag path falls back to raw installed
families when a readable roster yields zero resolvable IDs, and its new selection loop also
includes rows explicitly marked inactive. Separately, machine-scope initialization still
uses deck/local/environment agent specs both for proposals and for validating the machine
target. It can therefore write or accept a machine mapping whose adapter exists only in the
current deck and is unresolved everywhere else. These are correctness and autonomous-launch
scope bugs, not convenience-command deferrals.

## Refutation attempts

- Read the supplied 1,116-line `cdef991..HEAD` code diff, the complete repository diff, HEAD
  source, FINAL/consensus, both prior reviews, and all four implementation cycles plus their
  deferrals.
- Traced the no-flag path from `runTask` through `protocol.ReadRosterIDs`,
  `agents.ResolveParticipant`, and `selectedParticipantIDs`, including a fully mapped roster,
  a readable but completely unmapped roster, a roster with only unavailable adapters, and a
  mapped installed row marked `inactive`.
- Traced session and machine init through `resolveRoster`, `config.LoadAgentSpecs`,
  `config.LoadRosterAdapters`, `RosterAdaptersInFile`, candidate validation, and atomic write.
  Challenged both a missing machine mapping proposed from a deck-only custom family and an
  existing machine mapping whose family is defined only by the deck.
- Verified the quoted/empty-key defense: the complete candidate is parsed by
  `ValidateAgentsConfigBytes` before `WriteFileAtomic`, so equivalent quoted and unquoted
  duplicate table declarations are rejected without replacing the target.
- Checked every built-in `AutonomousWrite` declaration and the `Declared`/`Confined` split;
  only codex retains `Scope="workspace"`. Checked non-agy parenthesized labels and empty labels
  through the updated render implementation and tests.
- Ran `go build ./...`, `go vet ./...`, formatting over every cycle-4 Go file, and
  `git diff --check cdef991..HEAD`: all clean. `go test ./internal/agents ./internal/app
  ./internal/config` passes. `go test ./...` has only the previously documented,
  environment-dependent `internal/runner.TestDurableKillEndToEndRealProcess` failure because
  this host has no recorded boot ID; it is unrelated to this change.

## Round-02 OPEN item dispositions

### [CRITICAL] RESOLVED — Autonomous-write confinement honesty

`internal/agents/discover.go:101-108` now separates configured autonomy (`Declared`) from a
demonstrated sandbox (`Confined`). Codex keeps `Scope="workspace"` with
`--sandbox workspace-write`; claude, agy, and hermes have empty scope at
`discover.go:204,229,278`. The data model no longer claims those bypass/yolo modes are
workspace-confined.

### [MAJOR] OPEN — Default run path still has a raw-family escape hatch

The new `internal/app/app.go:1816-1829` loop correctly writes roster IDs into
`participantsFlag` when at least one ID resolves. But when the §2 roster is readable and
non-empty yet none of its IDs resolve, it leaves the flag empty. `selectedParticipantIDs`
then takes its legacy empty-input branch at `app.go:2413-2415` and returns every installed
family. A deck with roster `claude-1` and no mapping therefore still starts a no-flag run as
`claude`; a deck whose roster contains no installed family can launch unrelated installed
agents that are not roster members at all. This contradicts both the fix claim (“falls back
only when there is no §2 roster”) and the locked roster/quorum boundary.

**Required fix:** distinguish “no readable §2 roster” from “readable roster with zero
installed/resolvable members.” Use raw families only in the former case; hard-stop in the
latter. Add a no-flag regression test for the zero-resolved case.

### [MAJOR] RESOLVED — Empty/quoted roster table cannot be atomically duplicated

`internal/app/roster.go:347-354` validates the full candidate through
`config.ValidateAgentsConfigBytes` before replacement. This structurally catches the quoted
empty-table case and leaves the target unchanged. The earlier text-spelling dependency is
gone.

### [MAJOR] OPEN — Machine scope still consumes project-only families

Cycle 4 validates target mappings, but it builds `known` with
`config.LoadAgentSpecs(root)` at `internal/app/roster.go:224-228`. That loader deliberately
includes central, deck, local, and environment layers (`internal/config/runtime.go:134-165`).
The proposals were already built from those same layers through `resolveRoster` at
`roster.go:207`.

Consequently, if the deck defines `[agents.acme]` and maps a roster ID to `acme`,
`roster init --scope machine --yes` can write only `[roster.<id>] adapter="acme"` to the
central file, without copying the `acme` spec. Outside this deck the machine mapping cannot
resolve. The same layered catalog makes an already-existing central `adapter="acme"` pass
cycle-4 validation merely because this deck defines `acme`, masking a broken machine target.
This is the exact “never copies deck values up” correctness boundary in consensus §B, so its
documentation as a MINOR follow-up does not make the deferral acceptable.

**Required fix:** for machine scope, build proposals and validate existing adapters from a
machine-only catalog (built-ins plus the central target), excluding deck, local, and
`PARLEY_HEADLESS_AGENT_CONFIG` layers. Add both copy-up and masked-target regression tests.

### [MINOR] RESOLVED — Parenthesized tier extraction is agy-only

`internal/agents/naming.go:193-206` gates `StripParenTier` and tier substitution on
`family == "agy"`; other families sanitize the complete model label and retain their
configured effort. The new non-agy and empty-model cases cover the intended behavior.

### [NIT] RESOLVED — `selectedAgents` documentation matches fail-closed behavior

`internal/runner/runner.go:358-363` now states that unresolved identities are returned to the
caller and materialized as failed results rather than skipped.

### [MINOR] OPEN — `roster diff` and TOML `autonomous_write` remain deferred

Neither surface is implemented. I accept these as documented additive scope deferrals for
this release; they do not change the verdict.

### [NIT] OPEN — Unresolved pseudo-results still lack participant events

`RunRoundOne` still appends unresolved `Result` values after `run.segment_started`, whose
targets contain only resolved agents, and emits no per-agent failure event. This remains the
non-blocking projection-parity issue from round 2.

### [MINOR] OPEN — Phase-8 metadata remains stale

`IMPLEMENTATION.md` still has `status: implemented` and no top-level `completed` or
`head-commit`; cycle 4 also lacks the protocol-required status/completed/head-commit block.
Its lines 159-164 additionally repeat the superseded claim that claude/agy/hermes assert
`Scope="workspace"`, contradicting lines 126-130 and current code. This is an audit/resume
defect, not a runtime blocker, but the round-02 item is not resolved.

## New cycle-4 regressions

### [MAJOR] Default selection launches inactive roster members

`protocol.ReadRosterIDs` returns the roster set and a separate inactive subset. The new
default loop discards that second value (`active, _, ok`) at `internal/app/app.go:1817` and
attempts to resolve every row. Thus a mapped, installed member explicitly retained as
`inactive` is silently included in a no-flag run. That defeats the protocol's inactive
retention semantics and can autonomously launch a participant the roster deliberately
disabled.

**Required fix:** exclude IDs in the returned inactive set before resolving/sorting, and add
a no-flag test with one active and one inactive installed mapping.

## Documented deferrals

The `roster diff`, TOML `autonomous_write`, OS-enforced sandbox/live-sentinel, S6 display
plumbing, and fully interactive model/effort picker omissions are acceptable documented scope
for this review. The layered machine-scope behavior is not: it writes or blesses semantically
invalid machine state and is therefore a correctness bug.

## Verdict

**BLOCK**

The autonomous, TOML-structure, render, and documentation fixes are valid, but the raw-family
fallback, inactive-member launch, and project-to-machine family leakage are genuine unfixed
MAJOR correctness issues.
