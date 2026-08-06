---
idea: roster-operations-standard
phase: 8 — re-review
agent: codex-1
round: 4
date: 2026-08-06
reviewed-commit: a848e67
verdict: FINDINGS
---
# Re-review round 4 — codex-1

## Verdict

PRIMARY — FINDINGS. Cycle 4 fixes the quorum behavior, the bare-family continuation regression,
machine-scoped `--all`, `[agents.<family>]` provenance, the guard's missing-sibling disclosure,
and the exact JSON row contract. It also leaves state diagnostics on the old layered semantics:
`roster show --explain` names an ignored value-only layer as the source of the effective state,
and `roster set --state` falsely says a successful authority change was masked. That is a new
MAJOR finding below.

PRIMARY — I own the new finding and recommendation below, so I issue no §15.1 truth verdict on
them. Labels such as FIXED adjudicate cycle-4 implementation claims first made by the implementer;
they are not self-verdicts on my round-03 findings.

PRIMARY — The exact required command, `go build ./... && go test ./...`, exited 1. `go build
./...` succeeded silently. `go test ./...` failed in the pre-existing durable-kill case;
`git diff --name-only 7220715 a848e67 -- internal/runner` printed no path. The exact combined
output was:

```text
?   	parley-deck-cli/cmd/parley	[no test files]
ok  	parley-deck-cli/internal/acp	(cached)
ok  	parley-deck-cli/internal/agents	(cached)
ok  	parley-deck-cli/internal/app	30.487s
ok  	parley-deck-cli/internal/config	(cached)
ok  	parley-deck-cli/internal/consensus	(cached)
ok  	parley-deck-cli/internal/driver	0.939s
ok  	parley-deck-cli/internal/fsutil	(cached)
ok  	parley-deck-cli/internal/hitl	(cached)
ok  	parley-deck-cli/internal/loop	(cached)
ok  	parley-deck-cli/internal/pipeline	(cached)
ok  	parley-deck-cli/internal/procctl	0.251s
ok  	parley-deck-cli/internal/protocol	(cached)
ok  	parley-deck-cli/internal/repomap	(cached)
ok  	parley-deck-cli/internal/retro	(cached)
ok  	parley-deck-cli/internal/runaction	(cached)
ok  	parley-deck-cli/internal/runcontrol	(cached)
ok  	parley-deck-cli/internal/runmanifest	(cached)
--- FAIL: TestDurableKillEndToEndRealProcess (0.02s)
    durablekill_test.go:116: a live attributed process should be killed, got {AgentID:sleeper Killed:false Cleared:false Failed:true SegmentID:segment-0001 Message:process verification failed (no recorded boot id); not killed}
FAIL
FAIL	parley-deck-cli/internal/runner	7.251s
ok  	parley-deck-cli/internal/runplan	(cached)
ok  	parley-deck-cli/internal/runstate	(cached)
ok  	parley-deck-cli/internal/sessionstore	(cached)
ok  	parley-deck-cli/internal/steer	(cached)
ok  	parley-deck-cli/internal/store	(cached)
ok  	parley-deck-cli/internal/track	(cached)
ok  	parley-deck-cli/internal/tui	0.311s
FAIL
```

## Round-03 findings: fixed or not

### [CRITICAL codex-1 R3-C1] authority state — FIXED behaviorally; diagnostics diverge

PRIMARY — I created a scratch deck through the real CLI with committed `claude-1 active=true`
and `kimi-1 active=false`. The machine file, `agents.local.toml`, and
`$PARLEY_HEADLESS_AGENT_CONFIG` all contained both IDs with the opposite values. The production
`roster show --json` result was:

```json
{"claude-1":"active","kimi-1":"inactive"}
```

PRIMARY — I then put scratch-only `claude` and `kimi` executables (symlinks to `/usr/bin/true`)
on `PATH` and ran the production no-flag selection path with `parley run --no-tui --no-auto
--no-preflight --no-ping --yes`. It printed `Starting round-01 with participants: claude-1`,
and the created scratch `00-prompt.md` recorded `participants: [claude-1]`. Thus both requested
consumers honor the committed deck state in both directions: value-only layers neither retire
`claude-1` nor revive `kimi-1`.

PRIMARY — The implementation basis is `RosterScope.applyAuthorityState` at
`internal/config/runtime.go:184,199,204,209-220`. The shipped test does not cover the full
behavioral matrix; that test-quality gap is assessed below. The new diagnostic mismatch caused
by the same change is R4-M1.

### [MAJOR codex-1 R3-M1 / kimi-1 N-1] first-match snapshot regression — FIXED

PRIMARY — I drove `applyRosterSnapshotToParticipants` and then the runner's actual public
consumer, `agents.ResolveParticipant`, for both a bare ID and a composite ID sharing the same
adapter. The executed probe logged:

```text
claude -> model=frozen-bare effort=max argv=[--bare-frozen]
claude-1 -> model=frozen-composite effort=low argv=[--composite-frozen]
```

PRIMARY — `internal/app/roster_snapshot_apply.go:98-113` replaces an exact-ID collision in place
and appends only when no collision exists. The bare record therefore precedes no live duplicate,
while the composite record remains independently addressable.

### [MAJOR codex-1 R3-M2] machine scope and agent-block provenance — FIXED as scoped

PRIMARY — With an environment fixture setting `[agents.agy].model = "Gemini 3.5 Flash (High)"`,
`PARLEY_HEADLESS_AGENT_CONFIG=<fixture> parley roster show --scope machine --all --json` returned
the unrostered `agy` row with the machine/built-in value `Gemini 3.6 Flash (High)`, not the env
value. `unrosteredRows` now threads machine scope through both loaders at
`internal/app/roster_view.go:98-105`.

PRIMARY — On the same production surface, `roster show --scope machine --explain claude-1`
printed `model claude-opus-5[1m] ~/.parley/agents.toml`; the value comes from central
`[agents.claude]`. `AgentFieldSources` at `internal/config/runtime.go:1135-1167` fixes the exact
`[agents.<family>]` attribution named in R3-M2. R4-M1 records other effective-source cases that
the new source reconstruction still gets wrong.

### [MAJOR codex-1 R3-M3] four-surface guard claim — FIXED

PRIMARY — I ran the committed test from a `git archive a848e67` extraction with no sibling skill
checkout. It logged both missing skill surfaces and then this exact line before passing:

```text
NOTE: checked 2 of 4 surfaces — the bundled skill protocol and SKILL.md live in the sibling parley-deck-skill checkout and were not readable here; they are NOT enforced by this run
```

PRIMARY — The behavior is at `internal/protocol/drift_test.go:275-303`. A run with the sibling
present checks four surfaces; an isolated CLI run now truthfully reports that it checks only two.

### [MINOR codex-1 R3-m1] exact JSON row contract — FIXED

PRIMARY — `TestJSONRowHasExactlyTheFrozenColumns` calls production `rosterShow`, decodes every row
as a key map, rejects any key outside the frozen eleven, and rejects every missing frozen key
(`internal/app/roster_cycle2_test.go:264-296`). The test passed in the targeted run.

### [MINOR both] skill 2.5.1 co-release — pending release step, not re-filed against this commit

PRIMARY — The skill repository currently remains at `b806ada` / tag `v2.5.0`, with
`skills/parley-deck/SKILL.md` and `skills/parley-deck/references/COOPERATION.md` modified and
`package.json` still at `2.5.0`. The corrected content is present in those working files.

PRIMARY — The current operator brief explicitly says the skill is being released alongside this
CLI version. I therefore treat the version/commit/tag work as the stated coordinated release gate,
not as a new defect in CLI commit `a848e67`; this review does not independently establish that the
future release has occurred.

### [NIT codex-1 R3-N1] implementation handoff metadata — retained without a self-verdict

PRIMARY — `IMPLEMENTATION.md` still begins with blank lines followed by `## Fix-up cycle 2`; it
still has no document frontmatter, top-level `status`, or top-level `head-commit`. Cycle 4 appends
its section at line 133 but does not name this prior NIT among its fixes. Because codex-1 owns the
original claim, I retain the evidence and position without issuing a §15.1 verdict on it.

## New findings (by severity, or "none")

### [MAJOR] R4-M1 — effective-state provenance and masking still use the ignored layer stack

PRIMARY — I own this finding and issue no self-verdict. In the same four-layer fixture used for
R3-C1, `roster show` correctly reported `claude-1 active` and `kimi-1 inactive`, but `roster show
--explain` attributed both effective states to the env file that asserted the opposite:

```text
claude-1 — membership from parley-deck/agents.toml
active         active                   PARLEY_HEADLESS_AGENT_CONFIG:<scratch-env>/parley-deck/agents.toml

kimi-1 — membership from parley-deck/agents.toml
active         inactive                 PARLEY_HEADLESS_AGENT_CONFIG:<scratch-env>/parley-deck/agents.toml
```

PRIMARY — The write surface contradicts the effective result too. I changed the committed deck
state through `roster set` while the env file asserted the opposite. Both writes succeeded and
`roster show` reflected them, but each emitted this warning:

```text
warning: active = "true" is MASKED — PARLEY_HEADLESS_AGENT_CONFIG:<scratch-env>/parley-deck/agents.toml sets it at a higher layer, so the effective value did not change.
```

PRIMARY — The effective value did change, so this is false operator guidance on the command used
to retire or reactivate quorum members. `applyAuthorityState` overwrites layered `Active`, but
`RosterFieldSourcesScoped` still assigns `active` to the last raw block
(`internal/config/runtime.go:1101-1132`); `rosterExplain` trusts that source
(`internal/app/roster_view.go:165-190`), and `rosterFieldMaskedBy` uses the same raw-source model
for its warning (`internal/app/roster_set.go:80-115`).

PRIMARY — The new `AgentFieldSources` reconstruction is also incomplete for effective spec
sources: the current machine config has `[defaults] speed = "deep"`, and production machine
`--explain` printed `speed deep built-in default`. The helper scans only per-family `model`,
`reasoning`, and `speed` fields and never the `[defaults]` source that `LoadAgentSpecsScoped`
actually applies (`internal/config/runtime.go:1138-1167`).

PRIMARY — Suggested fix: derive provenance from the same resolved objects that produce the
effective row. `active` must name the membership authority and must never be considered masked by
a value-only layer; model/effort/speed should use the resolved spec's source metadata (with a
roster-ID override taking precedence), including `[defaults]`. Add one end-to-end test that uses
all four conflicting layers and asserts `roster show`, no-flag selection, `--explain active`, and
the absence of a false `masked-by-env` warning after `roster set --state`.

## Test-quality assessment

PRIMARY — I ran the three new assertion tests plus the prior composite and selection tests:

```text
TestSnapshotPinsSurviveParticipantResolution                 PASS
TestSnapshotFreezeReachesBareFamilyParticipants              PASS
TestValueLayersCannotChangeMembershipState                   PASS
TestJSONRowHasExactlyTheFrozenColumns                        PASS
TestDefaultRosterParticipants                                PASS
TestNoSection2AsAStoreInstructions                           PASS
```

PRIMARY — `TestSnapshotFreezeReachesBareFamilyParticipants` is reversion-sensitive: restoring
cycle 3's append makes `agents.ResolveParticipant("claude", ...)` select `drifted`, so its model
assertion fails. It tests the actual first-match consumer for the bare ID.

PRIMARY — The shipped composite test is still weaker: `TestSnapshotPinsSurviveParticipantResolution`
folds the helper output into a last-write-wins map and never calls `ResolveParticipant`, so it
would still pass with cycle 3's append bug. My ephemeral overlay probe exercised both bare and
composite records through `ResolveParticipant`, but that composite consumer assertion is not in
the committed suite.

PRIMARY — `TestValueLayersCannotChangeMembershipState` is reversion-sensitive: removing the
authority-state overwrite lets the machine retire `claude-1` and the local layer retire `kimi-1`,
failing its assertions. It does not satisfy the full claimed matrix: the deck declares both
members active, there is no env config, it tests retirement but not reactivation, and it calls
only `RosterMembership` rather than `rosterShow` and `defaultRosterParticipants`. My scratch CLI
run supplies the missing behavioral evidence, but those regressions remain unguarded in-tree.

PRIMARY — `TestJSONRowHasExactlyTheFrozenColumns` is reversion-sensitive in both directions:
restoring `display_name`/`note` or any other twelfth serialized field fails the unexpected-key
loop, while removing any frozen field fails the missing-key loop.

PRIMARY — Cycle 4 adds no committed regression test for the two R3-M2 code paths
(`--scope machine --all` and `[agents.<family>]` provenance). I verified both behaviorally, but
reverting either change leaves the new test set green.

PRIMARY — The claimed fourth test, the “corrected guard note,” is not a new assertion. It changes
only `t.Logf`; reverting that wording leaves `TestNoSection2AsAStoreInstructions` green. The
isolated behavioral run proves the current message is honest, but the honesty text is diagnostic,
not enforced.

PRIMARY — No cycle-4 assertion is a literal constant tautology. The closest functional tautology
remains the pre-existing composite snapshot test: it observes the helper through a map whose
last-write semantics erase the runner's first-match behavior. The uncovered diagnostic mismatch
in R4-M1 is exactly the kind of real-call-site divergence the new state test cannot detect.
