---
idea: roster-operations-standard
phase: 8 — re-review
agent: kimi-1
round: 6
date: 2026-08-06
reviewed-commit: d312c53
verdict: CLEAN
---
# Re-review round 6 — kimi-1

## Verdict

CLEAN. Round 05 left exactly one code finding, filed twice (codex-1 MAJOR R5-M1, kimi-1
MINOR R5-m1): cycle 5's `active` provenance override ignored scope, so `--scope machine
--explain` attributed state to `parley-deck/agents.toml` while its own membership header
named the machine file, and the post-write masking check made the same unscoped
comparison. Cycle 6 fixes both halves and nothing else. I verified on one fixture with
OPPOSITE deck/machine `active` values that all three surfaces — `roster show` state,
`--explain` provenance, and the `--explain` membership header — agree in BOTH scopes and
on every authority branch (deck TOML, legacy §2, inherited machine); that the masked
warning no longer fires on machine-scope state writes and still fires where it should;
that `TestActiveProvenanceIsScopeAware` FAILS when the fix is reverted; and that the
mandated command exits 0. No new findings.

**§15.1 ownership disclosure.** R5-m1 was my round-5 finding; I own it and issue no
self-verdict on it. FIXED labels below adjudicate the implementer's cycle-6 claims (as
stated in IMPLEMENTATION.md "Fix-up cycle 6"), not my own prior position. The skill
2.5.1 release is handled at the Phase 8 close per the brief — not a code finding, not
re-filed. All claims PRIMARY (read, run, or executed by me this session) unless tagged
otherwise.

**Method.** Binary built from d312c53 to `/tmp/parley-r6`; the cycle-5 binary for
before/after was built from 8c8a8f1 extracted read-only via `git archive` to
`/tmp/r6prev` (`/tmp/parley-r5prev`). Scratch decks under `/tmp/kr6*` with `PARLEY_HOME`
isolation. Reversion sensitivity proven with `go test -overlay` from `/tmp`. The
reviewed tree stayed byte-identical throughout: `git status --short` at the end shows
only the untracked `review/round-06/` directory this file lands in. No git write
commands used.

**Mandated command, run by me at d312c53.** `go build ./... && go test ./...` exited 0.
Exact output:

```text
?   	parley-deck-cli/cmd/parley	[no test files]
ok  	parley-deck-cli/internal/acp	(cached)
ok  	parley-deck-cli/internal/agents	(cached)
ok  	parley-deck-cli/internal/app	(cached)
ok  	parley-deck-cli/internal/config	(cached)
ok  	parley-deck-cli/internal/consensus	(cached)
ok  	parley-deck-cli/internal/driver	(cached)
ok  	parley-deck-cli/internal/fsutil	(cached)
ok  	parley-deck-cli/internal/hitl	(cached)
ok  	parley-deck-cli/internal/loop	(cached)
ok  	parley-deck-cli/internal/pipeline	(cached)
ok  	parley-deck-cli/internal/procctl	(cached)
ok  	parley-deck-cli/internal/protocol	(cached)
ok  	parley-deck-cli/internal/repomap	(cached)
ok  	parley-deck-cli/internal/retro	(cached)
ok  	parley-deck-cli/internal/runaction	(cached)
ok  	parley-deck-cli/internal/runcontrol	(cached)
ok  	parley-deck-cli/internal/runmanifest	(cached)
ok  	parley-deck-cli/internal/runner	(cached)
ok  	parley-deck-cli/internal/runplan	(cached)
ok  	parley-deck-cli/internal/runstate	(cached)
ok  	parley-deck-cli/internal/sessionstore	(cached)
ok  	parley-deck-cli/internal/steer	(cached)
ok  	parley-deck-cli/internal/store	(cached)
ok  	parley-deck-cli/internal/track	(cached)
ok  	parley-deck-cli/internal/tui	0.184s
MANDATED_EXIT=0
```

Mostly cached, so I also ran the suite uncached (`go test -count=1 ./...`, stricter than
mandated): exit 0, zero failures — including `internal/runner`, which fails only in
codex-1's sandbox and passes on this machine as in my rounds 3–5. `go vet ./...` clean.

## Round-05 finding: fixed or not

### [MAJOR codex-1 R5-M1 / MINOR kimi-1 R5-m1] machine-scope provenance and masking ignored scope — FIXED

PRIMARY — The fixture the finding named: committed deck roster `claude-1 active = false`,
machine roster `claude-1 active = true` — opposite values for the same agent. At d312c53:

- Deck scope, all three surfaces agree: `roster show --json` → `{'claude-1': 'inactive'}`;
  `--explain` prints header `claude-1 — membership from parley-deck/agents.toml` and
  `active  inactive  parley-deck/agents.toml`.
- Machine scope, all three agree: `roster show --scope machine --json` →
  `{'claude-1': 'active'}`; `--scope machine --explain` prints header
  `claude-1 — membership from /tmp/kr6/home/agents.toml` and
  `active  active  /tmp/kr6/home/agents.toml`.

PRIMARY — Before/after on the IDENTICAL fixture proves the delta is what fixed it. The
8c8a8f1 binary prints `active  active  parley-deck/agents.toml` under a header naming
the machine file (the round-5 defect, reproduced); the d312c53 binary prints
`active  active  /tmp/kr6cmp/home/agents.toml`, agreeing with its own header. The fix is
exactly what the finding suggested: `rosterExplain` overrides `layers["active"]` with
`scope.Source` — the object it already resolved for THIS scope (roster_view.go:181-183)
— instead of re-deriving the deck authority.

PRIMARY — Both other authority branches are self-consistent in both scopes too. Legacy
deck (§2 table only, machine asserting the opposite): deck-scope explain shows
`active  active  COOPERATION.md §2` under a `COOPERATION.md §2` header; machine-scope
shows `active  inactive  /tmp/kr6/home/agents.toml` under the machine header. Inherited
deck (no deck roster): deck-scope shows `active  inactive  ~/.parley/agents.toml` under
the `~/.parley/agents.toml (INHERITED — …)` header; machine scope agrees with its own
header. (Cosmetic, below threshold: machine-scope output labels the file by absolute
path where other rows use the `~/.parley/agents.toml` display label — consistent within
each output, which is what the test asserts.)

PRIMARY — Write surface, codex-1's exact R5-M1 repro: deck says `active = false`,
machine file flipped to `active = true` via
`roster set claude-1 --scope machine --state active --yes --confirm-breaking`. At
8c8a8f1 this wrote AND printed the spurious
`warning: active = "true" is MASKED — parley-deck/agents.toml sets it at a higher layer…`
while `roster show --scope machine` reported `active` (both directions warned at cycle
5, reproduced). At d312c53 the same write prints NO warning, machine-scope show flips to
`active`, deck-scope show keeps `inactive` — the write is described relative to the
scope it targeted, via `config.RosterStateSourceForTarget(root, target)`
(roster_set.go:107; runtime.go:1185-1197), which maps the central agents.toml to the
machine authority and everything else to `LoadRosterScoped`'s deck chain.

PRIMARY — Warning correctness in both directions, no spurious fires, no lost true fires:

- Machine-scope retire and revive (deck asserting the opposite both times): no warning
  either way; each machine view reflected the write. Correct — the machine roster
  governs machine scope.
- Deck-scope revive while `$PARLEY_HEADLESS_AGENT_CONFIG` asserts `active = false`: no
  warning (the R4-M1 surface stays fixed), deck show/explain flipped to
  `active … parley-deck/agents.toml`.
- True-positive control for a LAYERED field: deck `--model k3` write while the env file
  pins `model = "env-model"` still warns
  `model = "k3" is MASKED — PARLEY_HEADLESS_AGENT_CONFIG:/tmp/kr6/env.toml …` — the
  general layer-stack branch (roster_set.go:119-133) is untouched and live.
- Legacy edges: a deck-scope state write in a §2-only deck creates the block, flips the
  authority to the deck file, and warns nothing (the new `ap = authority` fallback at
  roster_set.go:111-114 is not reachable as a false positive here, because post-write
  `LoadRosterScoped` no longer returns the §2 label); a machine-scope state write in
  that deck warns nothing and flips the machine view only.

PRIMARY — One honest accounting note, not a finding: my round-5 review called the
machine-write MASKED warning a "TRUE-positive" because its substance was true of the
DECK view. codex-1's R5-M1 measured the same warning against the MACHINE view the
command was scoped to, where it was false. Cycle 6 adopts the scoped reading — the
warning is evaluated against the authority of the scope being written — and under that
reading the old warning was spurious and its absence now is correct. I verified the new
behavior is self-consistent per scope; the semantics question was settled by the
round-5 consensus, not re-litigated here.

### [MINOR, all three] skill 2.5.1 — carried; not re-filed

PRIMARY — The brief states the release happens at the Phase 8 close; nothing in
d312c53's diff touches it. RECALL — its state is recorded in the round-4/5 reviews; the
close must still perform it.

### Test quality — the new test is reversion-sensitive

PRIMARY — `TestActiveProvenanceIsScopeAware` (roster_cycle2_test.go:341-372) passes
natively at d312c53, then FAILS when its fix is reverted via `go test -overlay` from
`/tmp` (reviewed tree untouched): overlaying roster_view.go with the cycle-5 form
(`layers["active"]` from an unconditional `config.RosterStateSource(root)`) makes it
fail, `go test` exit 1, with the exact R5-m1 signature —
`machine scope: provenance "active … parley-deck/agents.toml" contradicts its own
membership header "/var/folders/…/agents.toml"`. The test drives the real `rosterExplain`
at both scopes and asserts the two halves of one output agree — precisely the guard
whose absence let R5-m1 ship. Not tautological.

PRIMARY — Coverage observation, below finding threshold: the masking half of the fix is
not reversion-guarded by any committed test. Overlaying roster_set.go to call the
cycle-5 unscoped `config.RosterStateSource(root)` instead of `RosterStateSourceForTarget`
leaves the whole `internal/app` suite GREEN (exit 0, verified) — the cycle-5 test
`TestActiveProvenanceAndMaskingFollowTheAuthority` targets only the deck authority, so
the machine-write warning path is guarded solely by the behavioral runs above. codex-1's
R5-M1 text suggested a test covering that surface too. Not a code defect: the behavior
is verified correct this session; recording so the coverage map is accurate.

## New findings (by severity, or "none")

none

PRIMARY — The cycle-6 diff touches four files; each hunk is exercised above or by the
suite. Full uncached suite green, vet clean. `RosterStateSource` (the cycle-5 helper) is
now referenced only by its cycle-5 test, production code having moved to
`RosterStateSourceForTarget` — an exported config API kept green by that test, not dead
code in any problematic sense. Deck-scope value layering, machine-scope scoping
(machine-sourced `model` row visible in both scopes' explains above), and the
legacy/inherited branches all behaved as at 8c8a8f1 except for the fixed attribution.
