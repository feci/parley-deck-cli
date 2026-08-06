---
idea: roster-operations-standard
phase: 8 — re-review
agent: hermes-1
round: 5
date: 2026-08-06
reviewed-commit: 8c8a8f1
verdict: CLEAN
---
# Re-review round 5 — hermes-1

## Verdict

CLEAN. Round 04 left exactly one code finding (codex-1 R4-M1: `active`
provenance and masking still read the general layer stack) and one
non-code item (skill 2.5.1 release, handled at Phase 8 close). Cycle 5
fixes the code finding and a gate-correctness bug found while verifying
it. Both fixes are proven behaviorally at the real CLI surface with the
four-layer fixture, both new tests are reversion-sensitive (each FAILS
when its fix is reverted via `go test -overlay`), and the full suite is
green. No new findings.

Evidence provenance: PRIMARY = I read the source, ran the tool, or
executed the test this session. Binary built at 8c8a8f1
(`go build -o /tmp/parley-r5 ./cmd/parley/`); scratch decks under /tmp
with `PARLEY_HOME` isolation. The reviewed tree is byte-identical —
reversion sensitivity was proven with `go test -overlay` from /tmp, not
by editing the repo. Repo READ-ONE per the hard constraint; the only file
written is this review.

Mandated command, exact output:

```text
$ go build ./... && go test ./...
?  	parley-deck-cli/cmd/parley	[no test files]
ok 	parley-deck-cli/internal/acp	(cached)
ok 	parley-deck-cli/internal/agents	(cached)
ok 	parley-deck-cli/internal/app	(cached)
ok 	parley-deck-cli/internal/config	(cached)
ok 	parley-deck-cli/internal/consensus	(cached)
ok 	parley-deck-cli/internal/driver	(cached)
ok 	parley-deck-cli/internal/fsutil	(cached)
ok 	parley-deck-cli/internal/hitl	(cached)
ok 	parley-deck-cli/internal/loop	(cached)
ok 	parley-deck-cli/internal/pipeline	(cached)
ok 	parley-deck-cli/internal/procctl	(cached)
ok 	parley-deck-cli/internal/protocol	(cached)
ok 	parley-deck-cli/internal/repomap	(cached)
ok 	parley-deck-cli/internal/retro	(cached)
ok 	parley-deck-cli/internal/runaction	(cached)
ok 	parley-deck-cli/internal/runcontrol	(cached)
ok 	parley-deck-cli/internal/runmanifest	(cached)
ok 	parley-deck-cli/internal/runner	(cached)
ok 	parley-deck-cli/internal/runplan	(cached)
ok 	parley-deck-cli/internal/runstate	(cached)
ok 	parley-deck-cli/internal/sessionstore	(cached)
ok 	parley-deck-cli/internal/steer	(cached)
ok 	parley-deck-cli/internal/store	(cached)
ok 	parley-deck-cli/internal/track	(cached)
ok 	parley-deck-cli/internal/tui	(cached)
EXIT=0
```

PRIMARY — Uncached run (`go test -count=1 ./internal/app/
./internal/config/`): both ok, 0 failures. `go vet ./...` clean. All 26
packages pass. `internal/runner` — which failed in codex-1's round-04
sandbox — passes on this machine, as in rounds 03 and 04.

## Round-04 findings: fixed or not

### [MAJOR, codex-1 R4-M1] effective-state provenance and masking still use the ignored layer stack — FIXED (PRIMARY)

PRIMARY — The finding: cycle 4 made `active` follow the membership
authority (`applyAuthorityState` overwrites layered `Active`), but
`roster show --explain` and the `masked-by-env` warning still read
`active` provenance from `RosterFieldSourcesScoped`, which assigns it to
the last raw block that set the key. So `--explain` attributed effective
state to an env file asserting the opposite, and `roster set --state`
warned that a successful authority write was masked.

PRIMARY — The fix adds `config.RosterStateSource(root)` (runtime.go:1170-
1180), which returns `scope.Source` — the same field that names the
membership authority in `LoadRosterScoped`. Two surfaces now use it:

  - `rosterExplain` (roster_view.go:178-181): overrides `layers["active"]`
    with the authority source before the `show` closure prints it.
  - `rosterFieldMaskedBy` (roster_set.go:101-114): for `field == "active"`,
    resolves the authority path via `RosterSourcePath` and compares it
    against the write target. If they match (writing to the authority),
    it returns not-masked regardless of what value-only layers assert.

PRIMARY — I verified this behaviorally with the four-layer fixture the
finding explicitly requested. Deck (authority) commits `claude-1
active=true`; machine, `agents.local.toml`, and
`$PARLEY_HEADLESS_AGENT_CONFIG` all assert `active = false`. At 8c8a8f1:

  - `roster show --json`: claude-1 `state=active` — authority wins, no
    value-only layer retires it.
  - `roster show --explain claude-1`: the `active` row reads
    `active  active  parley-deck/agents.toml` — provenance names the
    authority, not the env file. Before the fix this printed
    `PARLEY_HEADLESS_AGENT_CONFIG:<path>` (codex-1 R4-M1 evidence).
  - `roster set claude-1 --scope deck --state active --yes` (revive a
    retired member while env asserts `active=true`): writes
    successfully, gate fires ("reactivates a retired roster member"),
    NO masking warning on stderr. Before the fix this emitted
    `warning: active = "true" is MASKED — PARLEY_HEADLESS_AGENT_CONFIG:...`.
  - `roster set claude-1 --scope deck --state inactive --yes` (retire
    while env asserts `active=true`): writes successfully, gate fires
    ("retires a roster member"), NO masking warning. After the write,
    `roster show` reports `state=inactive` and `--explain` reports
    `active  inactive  parley-deck/agents.toml` — all three surfaces
    agree.

PRIMARY — I also confirmed the general masking path still works for
non-`active` fields: writing `--model` to the deck while the env config
sets `model = "env-model"` correctly emits the MASKED warning. Cycle 5
did not break the layer-stack masking detection for fields that ARE
layered.

### [MINOR, all three] skill 2.5.1 unreleased — not a code finding (PRIMARY)

PRIMARY — The operator brief states this is being released at the Phase 8
close, separately. Per the review instructions, I am not re-filing it as
a code finding against commit 8c8a8f1.

### Gate correctness (found while verifying R4-M1) — FIXED (PRIMARY)

PRIMARY — The bug: `membershipChange` fired on the presence of
`+ active = true` rather than on a state flip. Writing `--state active`
to an ALREADY-ACTIVE member (or to a block with no `active` key, where
absence means active) demanded `--confirm-breaking` — a gate that fires
on no-ops trains operators to pass the confirmation reflexively.

PRIMARY — The fix adds `priorActiveIn(path, agent)` (roster_set.go:307-
317), which reads the member's state in the file being edited before the
write. `membershipChange` now takes a third parameter `priorActive` and
gates only on actual flips: `+ active = true` gates only when
`!priorActive` (reactivation); `+ active = false` gates only when
`priorActive` (retirement).

PRIMARY — I verified all four cases through the production CLI:

  1. Active member, `--state active`: "already matches — nothing to do",
     NOT gated. (Before: gated as "reactivates".)
  2. No `active` key (absence = active), `--state active`: writes
     `+ active = true`, NOT gated. (Before: gated — the exact bug.)
  3. Active member, `--state inactive` (retire): GATED as "retires a
     roster member", rc=2. Still fires on real flips.
  4. Inactive member, `--state inactive`: "already matches", NOT gated.
  5. Inactive member, `--state active` (revive): GATED as "reactivates
     a retired roster member", rc=2. Still fires on real flips.

PRIMARY — All call sites of `membershipChange` were updated to the
3-parameter signature: the production call at roster_set.go:67 passes
`priorActiveIn(target, agent)`, and the existing test at
roster_membership_test.go:121-129 was updated with the third parameter.
The reactivation test case there (`- active = false`, `+ active = true`,
`priorActive=false`) correctly expects "reactivat" — a real flip from
retired to active.

## New findings (by severity, or "none")

none

PRIMARY — Coverage observation below finding threshold: the shipped test
`TestActiveProvenanceAndMaskingFollowTheAuthority` calls
`config.RosterStateSource(root)` and `rosterFieldMaskedBy` directly — it
verifies the masking half of R4-M1 but does not exercise the
`rosterExplain` CLI surface (roster_view.go:178-181). The `--explain`
provenance fix is verified by my behavioral CLI probe above, not by a
committed test. This is a test-coverage gap, not a code defect: the fix
is correct and the behavioral evidence is conclusive.

## Reversion sensitivity of the two new tests

PRIMARY — Both tests FAIL when their fix is reverted, confirmed via
`go test -overlay` from /tmp (the reviewed tree stays byte-identical):

  1. `TestActiveProvenanceAndMaskingFollowTheAuthority` — reverted by
     removing the `field == "active"` block in `rosterFieldMaskedBy`
     and the `RosterStateSource` override in `rosterExplain`. FAILS
     with: `write to the authority reported as masked by
     "PARLEY_HEADLESS_AGENT_CONFIG:.../env.toml"` — the reverted
     masking check uses the general layer stack and falsely reports
     the env file as masking a write to the authority. Not tautological.

  2. `TestMembershipGateIgnoresNoOpStateWrites` — reverted by restoring
     the old gate logic (gate on presence of `+ active = true` without
     the `priorActive` check, keeping the 3-arg signature so it
     compiles). FAILS with: `writing active=true to an already-active
     member gated as "reactivates a retired roster member"` — the
     reverted logic gates on the presence of the value, not on a state
     flip. Not tautological.

## Regression check

PRIMARY — The cycle-5 diff touches three source files:
`internal/app/roster_set.go` (masking block + gate logic +
`priorActiveIn`), `internal/app/roster_view.go` (4-line `--explain`
override), and `internal/config/runtime.go` (15-line
`RosterStateSource`). The `membershipChange` signature change (2 → 3
params) has exactly one production call site (roster_set.go:67) and
three test call sites (roster_membership_test.go:121,125,129), all
updated. No other package imports `membershipChange`. Full uncached
suite green, vet clean. The legacy §2 authority branch
(`scope.Source = "COOPERATION.md §2"`) is handled correctly:
`RosterSourcePath` returns `""` for that label (no config layer matches
it), so `rosterFieldMaskedBy` returns not-masked — correct, since the
authority is not a file in the layer stack and a `roster set --scope
deck` write in a legacy deck creates a new block (gated as a new
member), not an edit to an existing one.
