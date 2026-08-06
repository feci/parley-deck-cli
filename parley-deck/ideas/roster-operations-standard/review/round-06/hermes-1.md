---
idea: roster-operations-standard
phase: 8 — re-review
agent: hermes-1
round: 6
date: 2026-08-06
reviewed-commit: d312c53
verdict: CLEAN
---
# Re-review round 6 — hermes-1

## Verdict

CLEAN. Round 05 left exactly one code finding (codex-1 MAJOR R5-M1, kimi-1 MINOR
R5-m1 — the same defect): cycle 5's `active` provenance override in
`rosterExplain` was unscoped, so `roster show --scope machine --explain`
attributed state to `parley-deck/agents.toml` while its own membership header
named the machine file. Cycle 6 fixes it: `rosterExplain` now uses the
already-resolved `scope.Source`, and the masking check uses
`config.RosterStateSourceForTarget`, which recognizes a machine-scope write as
governed by the machine roster. The new test
`TestActiveProvenanceIsScopeAware` asserts the provenance line agrees with the
membership header of the same output, in both scopes. I verified this
behaviorally at the real CLI surface with opposite deck/machine `active`
values, confirmed the masking warning fires correctly and not spuriously,
proved the new test fails when the fix is reverted (via `go test -overlay`,
repo byte-identical), and ran the full mandated command — all 26 packages
pass, exit 0, vet clean. No new findings.

Evidence provenance: PRIMARY = I read the source, ran the tool, or executed
the test this session. Binary built at d312c53
(`go build -o /tmp/parley-r6 ./cmd/parley/`); scratch decks under /tmp with
`PARLEY_HOME` isolation. Reversion sensitivity was proven with
`go test -overlay` from /tmp — the reviewed tree stayed byte-identical
(`git status --short` clean). Repo READ-ONLY per the hard constraint; the
only file written is this review.

Mandated command, exact output:

```text
$ go build ./... && go test ./...
?   parley-deck-cli/cmd/parley   [no test files]
ok  parley-deck-cli/internal/acp (cached)
ok  parley-deck-cli/internal/agents (cached)
ok  parley-deck-cli/internal/app 18.076s
ok  parley-deck-cli/internal/config (cached)
ok  parley-deck-cli/internal/consensus (cached)
ok  parley-deck-cli/internal/driver (cached)
ok  parley-deck-cli/internal/fsutil (cached)
ok  parley-deck-cli/internal/hitl (cached)
ok  parley-deck-cli/internal/loop (cached)
ok  parley-deck-cli/internal/pipeline (cached)
ok  parley-deck-cli/internal/procctl (cached)
ok  parley-deck-cli/internal/protocol (cached)
ok  parley-deck-cli/internal/repomap (cached)
ok  parley-deck-cli/internal/retro (cached)
ok  parley-deck-cli/internal/runaction (cached)
ok  parley-deck-cli/internal/runcontrol (cached)
ok  parley-deck-cli/internal/runmanifest (cached)
ok  parley-deck-cli/internal/runner (cached)
ok  parley-deck-cli/internal/runplan (cached)
ok  parley-deck-cli/internal/runstate (cached)
ok  parley-deck-cli/internal/sessionstore (cached)
ok  parley-deck-cli/internal/steer (cached)
ok  parley-deck-cli/internal/store (cached)
ok  parley-deck-cli/internal/track (cached)
ok  parley-deck-cli/internal/tui (cached)
EXIT=0
```

PRIMARY — Uncached run (`go test -count=1 ./...`): exit 0, all 26 packages
ok, 0 failures. `go vet ./...` clean. `internal/runner` passes on this
machine (10.060s), as in rounds 03–05.

## Round-05 finding: fixed or not

### [MAJOR, codex-1 R5-M1 / MINOR, kimi-1 R5-m1] cycle 5's unscoped provenance override — FIXED (PRIMARY)

PRIMARY — The finding: `rosterExplain` (roster_view.go:176-180 at 8c8a8f1)
called `config.RosterStateSource(root)`, which always resolves the DECK
authority via `LoadRosterScoped(root)`, regardless of the `--scope` being
asked about. So `roster show --scope machine --explain AGENT` printed a
membership header naming the machine file while the `active` SET BY cell named
`parley-deck/agents.toml` — the same output contradicting itself. The matching
problem existed in `rosterFieldMaskedBy` (roster_set.go:101-114 at 8c8a8f1),
which used the same unscoped helper, so a machine-scope `--state` write was
described relative to the deck authority.

PRIMARY — The fix makes two changes:

1. `rosterExplain` (roster_view.go:181-183): replaces
   `config.RosterStateSource(root)` with `scope.Source` — the field already
   resolved by `rosterScopeFor(root, opts.scope)` at line 140. For
   `--scope machine`, `rosterScopeFor` returns `scope.Source = path` (the
   central agents.toml path, roster_view.go:42). For `--scope deck`,
   `LoadRosterScoped` returns `scope.Source = deckSource` (e.g.
   `parley-deck/agents.toml`). The provenance now follows the scope in view.

2. `rosterFieldMaskedBy` (roster_set.go:107): replaces
   `config.RosterStateSource(root)` with
   `config.RosterStateSourceForTarget(root, target)`. The new helper
   (runtime.go:1185-1198) checks whether `target` is the central agents.toml
   — if so, returns `"~/.parley/agents.toml"` (machine authority); otherwise
   falls through to `LoadRosterScoped(root)` (deck authority). A machine-scope
   write targets the central file, so the check recognizes it as governed by
   the machine roster; a deck-scope write targets the deck file, so it
   recognizes the deck authority.

PRIMARY — A secondary change at roster_set.go:112-113: the `ap == ""` fallback
changed from `return "", false` to `ap = authority`. This handles the
`"COOPERATION.md §2"` label, for which `RosterSourcePath` returns `""` (no
config layer matches that label). I traced this edge case: the fallback only
fires for a deck-scope write when the authority is §2, but the masking check
runs AFTER the write, and by then the deck file has roster blocks so
`LoadRosterScoped` returns `"parley-deck/agents.toml"` instead of §2. The §2
fallback is therefore unreachable in the post-write masking path. I confirmed
this behaviorally with a legacy §2 fixture (see below).

PRIMARY — I verified the fix behaviorally at the real CLI surface. Fixture:
deck `agents.toml` commits `claude-1 active = false`; machine `agents.toml`
commits `claude-1 active = true` (OPPOSITE values). At d312c53:

```text
$ parley roster show --scope deck --json
  "agent": "claude-1", "state": "inactive"

$ parley roster show --scope deck --explain claude-1
claude-1 — membership from parley-deck/agents.toml

FIELD          EFFECTIVE                SET BY
adapter        claude                   parley-deck/agents.toml
active         inactive                 parley-deck/agents.toml
status: inactive

$ parley roster show --scope machine --json
  "agent": "claude-1", "state": "active"

$ parley roster show --scope machine --explain claude-1
claude-1 — membership from /tmp/parley-r6-fix/home/agents.toml

FIELD          EFFECTIVE                SET BY
adapter        claude                   ~/.parley/agents.toml
active         active                   /tmp/parley-r6-fix/home/agents.toml
```

PRIMARY — In both scopes, all three surfaces agree: `roster show` state,
`--explain` provenance (SET BY), and the `--explain` membership header all
name the same authority and report the same state. The machine-scope output
names the machine file in both the header and the `active` SET BY cell — the
exact contradiction R5-M1 identified is gone.

PRIMARY — Masking warning: `roster set --scope machine --state inactive
--yes --confirm-breaking` (real retire flip, writing to the machine file)
wrote successfully, changed the machine-scope state to `inactive`, and emitted
NO masking warning — correct, because the machine file is the authority for
machine scope. `roster set --scope deck --state active --yes --confirm-breaking`
(real revive flip, writing to the deck file) wrote successfully, changed the
deck-scope state to `active`, and emitted NO masking warning — correct,
because the deck file is the authority for deck scope. No spurious warnings
in either direction.

PRIMARY — Legacy §2 edge case: a deck with only a §2 table (no `agents.toml`)
correctly shows `active ... COOPERATION.md §2` in `--explain`, with the header
matching. Writing `--scope deck --state inactive` to create a new deck roster
block wrote successfully, emitted NO masking warning (the post-write
`LoadRosterScoped` finds the new block and returns `parley-deck/agents.toml`,
so the authority path resolves and `tp == app` → not masked), and the
subsequent `--explain` reported `active inactive parley-deck/agents.toml`
with the header matching. The `ap = authority` fallback at roster_set.go:113
is safe in practice because the §2 authority is superseded by the new deck
file before the masking check runs.

### [MINOR, all three] skill 2.5.1 release — not a code finding (PRIMARY)

PRIMARY — The operator brief states this is released at Phase 8 close,
separately. Per the review instructions, I am not re-filing it as a code
finding against commit d312c53.

## New findings (by severity, or "none")

none

PRIMARY — Coverage observation below finding threshold: the new test
`TestActiveProvenanceIsScopeAware` drives `rosterExplain` directly (not via
the CLI), so it guards the `roster_view.go` call site that
`TestActiveProvenanceAndMaskingFollowTheAuthority` (cycle 5) did not. However,
neither committed test exercises `rosterFieldMaskedBy` through a real
`roster set --state` invocation — the masking half is guarded by the
cycle-5 helper-level test and by my behavioral CLI probes above, not by a
committed end-to-end test. This is the same test-coverage gap noted in
round 05 by all three reviewers (the explain surface) and by codex-1
specifically (the masking surface); cycle 6 narrowed the explain gap but
the masking e2e gap remains. It is a coverage observation, not a code
defect: the fix is correct and the behavioral evidence is conclusive.

## Reversion sensitivity of the new test

PRIMARY — `TestActiveProvenanceIsScopeAware` FAILS when the fix is reverted,
confirmed via `go test -overlay` from /tmp (the reviewed tree stays
byte-identical, `git status --short` clean):

Reverted roster_view.go:181-183 from `if scope.Source != ""` back to
`if src, serr := config.RosterStateSource(root); serr == nil && src != ""`.
The test failed with:

```text
=== RUN   TestActiveProvenanceIsScopeAware
    roster_cycle2_test.go:369: machine scope: provenance "active         active
    parley-deck/agents.toml" contradicts its own membership header
    "/var/folders/.../TestActiveProvenanceIsScopeAware.../002/agents.toml"
--- FAIL: TestActiveProvenanceIsScopeAware (0.78s)
FAIL
```

PRIMARY — The failure is exactly R5-M1: the machine-scope `active` SET BY
cell names `parley-deck/agents.toml` (the deck authority, re-derived by the
unscoped `RosterStateSource`) while the membership header names the machine
file. The test catches this because it cross-checks the `active` line's
provenance against the membership header of the same output. Not tautological.

## Regression check

PRIMARY — The cycle-6 diff touches four source/test files:
`internal/app/roster_view.go` (4-line override change),
`internal/app/roster_set.go` (masking helper + fallback),
`internal/config/runtime.go` (19-line `RosterStateSourceForTarget`), and
`internal/app/roster_cycle2_test.go` (39-line new test). The pre-existing
`config.RosterStateSource` (runtime.go:1174-1179) is still used by the
cycle-5 test `TestActiveProvenanceAndMaskingFollowTheAuthority`
(roster_cycle2_test.go:311), so it is not dead code. No other package
imports the changed functions. Full uncached suite green, vet clean. No
behavioral regression detected in the deck-scope or machine-scope surfaces,
the gate, or the legacy §2 / inherited-machine authority branches.
