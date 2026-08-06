---
idea: roster-operations-standard
phase: 8 — re-review
agent: kimi-1
round: 4
date: 2026-08-06
reviewed-commit: a848e67
verdict: FINDINGS
---
# Re-review round 4 — kimi-1

## Verdict

**FINDINGS** — one MINOR. Every cycle-4 *code* fix is verified landed at the real call
site, proven behaviorally at a848e67 and proven reversion-sensitive by overlay: the
CRITICAL quorum-authority fix, the bare-family freeze regression, the machine-scope leak
and provenance, the guard wording, and the JSON contract test. The one false claim is the
release step: IMPLEMENTATION.md says "Skill 2.5.1 is released with this cycle rather than
left uncommitted," and the skill repo is byte-for-byte in its round-3 state — same HEAD,
same tag, same version, both files still uncommitted. Four cycles, and this is the cycle-4
instance of "documented as landed, not landed." All claims PRIMARY (verified this session)
unless tagged otherwise.

**§15.1 ownership disclosure.** R3-M1 is co-filed as my round-3 N-1; that position is
mine. What I verdict here are fresh implementation facts about code that first exists at
a848e67 (the replace-in-place helper, `applyAuthorityState`, `AgentFieldSources`, the new
tests), gathered this session. I issue no verdict on the round-3 position itself. codex-1
and hermes-1 verdict independently.

**Method.** Binary built from a848e67 to `/tmp/parley-rev4`; scratch decks under `/tmp/kr4`
with `PARLEY_HOME` isolation. Neither repository was written: `git status --short` in the
CLI repo shows only the untracked `review/round-04/` directory this file lands in, and the
skill repo's two ` M` entries predate this session (they are the round-3 state). Where the
CLI surface cannot observe an internal path I injected probe tests via `go test -overlay`
from `/tmp` — the reviewed tree stays byte-identical. Reversion sensitivity was proven the
same way: overlay a reverted copy of the fixed file, run the new test, watch it fail.

**Mandated command, run by me at a848e67.** `go build ./...` exited 0 silently. I ran the
suite uncached (`go test -count=1 ./...`, stricter than the mandated plain `go test
./...`): exit 0, all 26 packages ok, 0 failures. Exact output:

```text
BUILD_OK
?   	parley-deck-cli/cmd/parley	[no test files]
ok  	parley-deck-cli/internal/acp	0.243s
ok  	parley-deck-cli/internal/agents	0.480s
ok  	parley-deck-cli/internal/app	24.123s
ok  	parley-deck-cli/internal/config	1.082s
ok  	parley-deck-cli/internal/consensus	0.853s
ok  	parley-deck-cli/internal/driver	4.004s
ok  	parley-deck-cli/internal/fsutil	1.503s
ok  	parley-deck-cli/internal/hitl	1.732s
ok  	parley-deck-cli/internal/loop	2.138s
ok  	parley-deck-cli/internal/pipeline	1.925s
ok  	parley-deck-cli/internal/procctl	2.595s
ok  	parley-deck-cli/internal/protocol	2.559s
ok  	parley-deck-cli/internal/repomap	2.741s
ok  	parley-deck-cli/internal/retro	2.739s
ok  	parley-deck-cli/internal/runaction	2.660s
ok  	parley-deck-cli/internal/runcontrol	2.663s
ok  	parley-deck-cli/internal/runmanifest	2.653s
ok  	parley-deck-cli/internal/runner	9.822s
ok  	parley-deck-cli/internal/runplan	2.668s
ok  	parley-deck-cli/internal/runstate	2.686s
ok  	parley-deck-cli/internal/sessionstore	2.441s
ok  	parley-deck-cli/internal/steer	2.660s
ok  	parley-deck-cli/internal/store	2.678s
ok  	parley-deck-cli/internal/track	2.674s
ok  	parley-deck-cli/internal/tui	2.645s
TEST_EXIT=0
```

`internal/runner` — including `TestDurableKillEndToEndRealProcess`, which failed in
codex-1's sandbox in rounds 2 and 3 — passes on this machine, as it did in my round-3 run.
`go vet ./...` is clean.

## Round-03 findings: fixed or not

**[CRITICAL, codex-1] R3-C1 value-only layers could retire/revive a committed deck member
— FIXED, proven on BOTH required surfaces with all four layers in conflict.** Fixture:
deck commits `claude-1` + `kimi-1` (active) and `codex-1` (`active = false`, retired by
the deck itself); the machine file tries to retire claude-1 and revive codex-1;
`agents.local.toml` tries to retire kimi-1 and revive codex-1;
`$PARLEY_HEADLESS_AGENT_CONFIG` tries to retire claude-1 and kimi-1. At a848e67:

- `roster show --json` (deck scope): claude-1 `state=active`, kimi-1 `state=active`,
  codex-1 `state=inactive` — every value-layer flip discarded, the deck's own retirement
  intact.
- Participant selection (overlay probe driving the production functions):
  `RosterMembership` → `active={claude-1,kimi-1}`, `inactive={codex-1}`;
  `defaultRosterParticipants` → `[claude-1 kimi-1]`. The run quorum is the committed
  deck's, not the layer stack's.

The mechanism is where the fix claims it: only `parley-deck/agents.toml` carries
`membership: true` (internal/config/runtime.go:389), only membership/machine layers
record `active` into `deckActive`/`machineActive` (runtime.go:158,168), and
`applyAuthorityState` then forces each member's `Active` to the authority's value
(runtime.go:211-222), called on all three authority branches (deck runtime.go:184, legacy
runtime.go:199, inherited-machine runtime.go:204). `mergeRosterEntry` still merges every
layer's `active` into `Entries` (runtime.go:251-253) — the discard happens at the
authority boundary, which is the correct shape: values layer, state does not.

Reversion proof: overlaying a no-op'd `applyAuthorityState` makes the shipped
`TestValueLayersCannotChangeMembershipState` FAIL, and my probe reproduces the exact
round-3 hole — the run quorum becomes `[codex-1]` (both committed members retired, the
deck-retired member revived). The fix, not the fixture, is what holds the quorum.

Same class on the legacy branch (cycle 4 re-routed it through `applyAuthorityState`, so I
probed it): legacy §2 deck (claude-1, kimi-1 active; hermes-1 `(inactive)`) plus machine
file retiring claude-1 and reviving hermes-1, local retiring kimi-1, env retiring
claude-1 → `RosterMembership` returns exactly the §2 states. Under the same revert
overlay all three rows flip. The legacy authority holds state the same way.

**[MAJOR, codex-1 R3-M1 / my N-1] bare-family freeze regression — FIXED, proven through
the runner's own resolver.** The helper now replaces in place when the frozen clone's ID
already exists in the discovery list and appends only otherwise
(internal/app/roster_snapshot_apply.go:98-113). My probe feeds the production input shape
(adapter-keyed discovery `{ID: claude, Model: drifted}`) through
`applyRosterSnapshotToParticipants` and then resolves each participant with
`agents.ResolveParticipant` exactly as the runner's `selectedAgents` does
(internal/runner/runner.go:364-372; rule 1 = first exact ID match,
internal/agents/resolve.go:47-53):

- bare family `claude` → launches `frozen-model` with `--frozen` argv (the round-3 defect
  shape: drifted model AND live args both dead);
- composite `claude-1` → `frozen-a`;
- two composites on one adapter (`claude-1`, `claude-2`) → `opus-a` / `opus-b`;
- bare + composite side by side (`claude`, `claude-1`) → `frozen-bare` / `frozen-comp`,
  no shadowing in either direction.

Reversion proof: overlaying the cycle-3 append shape makes the shipped
`TestSnapshotFreezeReachesBareFamilyParticipants` FAIL with "runner would launch
\"drifted\", not the frozen \"frozen-model\"", and under the same overlay my probe's two
bare-family subtests fail while both composite subtests pass — matching round 3's
diagnosis that composite IDs never collided. The still-unguarded one-line call site at
app.go:1184 (reverting the wiring, not the helper, passes the suite) is unchanged from
round 3 and remains deferred by consensus (DF-2); my probe closed it for this review.

**[MAJOR, codex-1] R3-M2 machine-scope leak + misattribution — FIXED, both halves.**
`unrosteredRows` is scope-threaded (roster_view.go:98-105, called with `opts.scope ==
"machine"` at roster.go:437). Fixture: machine file sets `[agents.claude] model =
"machine-central-model"`; env config sets `[agents.claude] model = "env-leak-model"` and
`[agents.kimi] model = "env-kimi-model"`. With the env layer live, `roster show --scope
machine --all --json` reports the unrostered claude row as `machine-central-model` and the
kimi-1 member row as `unknown` — no env value anywhere. Control: the same invocation at
deck scope shows `env-leak-model` / `env-kimi-model`, proving the env layer was actually
loaded and the machine-scope result is scoping, not a dead config. Baseline without the
env layer is identical to the with-env machine-scope result.

Provenance half: machine file carries `[agents.claude] model = "machine-central-model"`
plus `[roster.claude-1] adapter = "claude"` — codex-1's exact round-3 scenario. `roster
show --scope machine --explain claude-1` now reports `model machine-central-model
~/.parley/agents.toml` where round 3 printed `built-in default`. `AgentFieldSources`
(internal/config/runtime.go:1138-1168) attributes model/effort/speed to the last layer
whose `[agents.<family>]` block sets them, in the same low-to-high order the spec loader
applies, and `rosterExplain` consults it only when no `[roster.<id>]` source claims the
field (roster_view.go:175-181) — roster-block provenance still wins when both exist,
which matches which value actually reaches the launch.

**[MAJOR, codex-1] R3-M3 four-surface guard overclaim — FIXED, both branches observed.**
With the sibling checkout present (this workspace), `TestNoSection2AsAStoreInstructions`
passes with no skip logs — all four surfaces read. With an overlay pointing the sibling
paths at nonexistent `/tmp` locations, the test passes AND logs, verbatim: "NOTE: checked
2 of 4 surfaces — the bundled skill protocol and SKILL.md live in the sibling
parley-deck-skill checkout and were not readable here; they are NOT enforced by this
run" (drift_test.go:290-292). The guard no longer implies four-surface coverage it did
not perform. The remaining substance of R3-M3 — the skill surfaces' committed content —
is the carried-open release item below.

**[MINOR, codex-1] R3-m1 JSON contract test — FIXED, reversion-sensitive as claimed.**
`TestJSONRowHasExactlyTheFrozenColumns` (roster_cycle2_test.go:264-296) unmarshals real
`rosterShow` JSON output and asserts the exact eleven-key set in both directions: any
extra key fails, any missing frozen column fails. Reversion proof: overlaying `rosterRow`
with `json:"display_name"` / `json:"note"` restored makes it FAIL, naming both forbidden
keys. The cycle-3 shape-only test's hole is closed.

**[MINOR, both] skill 2.5.1 — NOT FIXED, and the cycle-4 handoff says otherwise.** See
New findings, R4-m1.

## New findings (by severity, or "none")

### R4-m1 — [MINOR] Cycle 4 documents skill 2.5.1 as released; the skill repo is unchanged from round 3

PRIMARY. IMPLEMENTATION.md's cycle-4 section states: "**[MINOR, both] Skill 2.5.1** is
released with this cycle rather than left uncommitted." Checked directly in
`.../parley-deck-skill`: HEAD is still `b806ada` (the round-3 HEAD), the newest tag is
still `v2.5.0`, `package.json` still reads `"version": "2.5.0"`, and `git status --short`
still shows ` M skills/parley-deck/SKILL.md` and ` M
skills/parley-deck/references/COOPERATION.md` — the same uncommitted working-tree edits
codex-1 and I recorded in rounds 2 and 3. Committed content still carries both banned
phrases: `SKILL.md:295` ("`parley roster sync` moves it across.") and
`references/COOPERATION.md:124` ("until `parley roster sync` moves it over"). A clean
checkout of the released skill therefore still ships the false migration instruction;
only the dirty working tree — which is what the CLI's drift guard reads — is correct.

Severity stays MINOR: the content is verified correct on all four surfaces in this
workspace and is drift-guarded here, and the CLI's own version is likewise still 1.40.1
(VERSION file, CHANGELOG's latest entry) — the cut was always the post-Phase-8 step. But
it must be recorded, not closed, for two reasons: the round-2/3 item is still open, and
cycle 4 now asserts in writing that it is done. "Documented as landed and not landed" is
the exact failure class this round exists to catch; the remedy is the one line in the
exit checklist my round-3 N-2 named: commit the two skill files, bump to 2.5.1, tag it,
and cut the CLI release as part of the Phase-8 close — or strike the sentence from
IMPLEMENTATION.md until that happens.

### Anything cycle 4 broke — nothing found

PRIMARY. Full uncached suite green (above), `go vet` clean, and the cycle-3 tests all
still pass natively (`TestSnapshotPinsSurviveParticipantResolution`,
`TestLegacySection2BeatsTheMachineRoster`,
`TestDeckMembershipIsTheDeckFileNotTheLayeredUnion` re-run individually). Behavioral
spot-checks beyond the probes: deck-scope `--all` still layers values correctly (the
control half of the R3-M2 fixture), machine-scope `--explain` still names
`built-in default` for genuinely unset fields, and the legacy-authority probe doubles as
a render-path regression check on `LoadRosterScoped`'s reordered branches. The cycle-4
diff touches five source files; each hunk is exercised above or by the suite.

## Test-quality assessment

PRIMARY. I read all three new tests and the modified guard, ran each natively at a848e67
(all PASS: `TestSnapshotFreezeReachesBareFamilyParticipants`,
`TestValueLayersCannotChangeMembershipState`, `TestJSONRowHasExactlyTheFrozenColumns`),
and — the decisive check this round — ran each against an overlaid revert of its own fix.
Results, by test:

- `TestValueLayersCannotChangeMembershipState` — REAL. Reverting `applyAuthorityState`
  to a no-op fails it (both committed members reported retired). Coverage note, below
  finding threshold: the shipped test exercises the machine and local layers; the env
  layer and the revival direction are covered only by my probe. The mechanism is uniform
  (non-membership layers never record state), so the untested permutations share the
  tested code path.
- `TestSnapshotFreezeReachesBareFamilyParticipants` — REAL, and it resolves through
  `agents.ResolveParticipant` rather than re-implementing the runner's rule. Reverting
  the helper to append fails it with the round-3 signature ("runner would launch
  \"drifted\""). Composite-ID coverage lives in the cycle-3 test; my probe confirmed all
  four ID shapes end-to-end. The app.go:1184 wiring remains suite-unguarded (carried
  from round 3, consensus-deferred).
- `TestJSONRowHasExactlyTheFrozenColumns` — REAL. Both directions asserted against real
  command output; restoring the two JSON tags fails it naming `display_name` and `note`.
  A twelfth column under any name, or a renamed frozen column, also fails it.
- Extended drift guard — the wording fix is observable, not just cosmetic: I forced the
  degraded branch via overlay and the NOTE logs exactly what is and is not enforced. In
  this workspace all four surfaces are read and the guard passes on content.

**No tautologies.** Each new test constructs the defect fixture, drives the production
path, and asserts the runner-visible or operator-visible outcome; each fails when its
fix is reverted. The nearest miss is the shipped R3-C1 test's env-layer/revival
permutations living only in my probe — a coverage gap, not a vacuous assertion.

**Reversion-sensitivity summary (all proven, none assumed):** no-op
`applyAuthorityState` → shipped C1 test FAILS, quorum collapses to `[codex-1]`; append
instead of replace → shipped M1 test FAILS with "drifted"; restore JSON tags → shipped
contract test FAILS naming both keys. That is the property the last four cycles kept
missing, and cycle 4 has it.
