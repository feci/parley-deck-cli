---
idea: roster-operations-standard
phase: 8 — re-review
agent: hermes-1
round: 4
date: 2026-08-06
reviewed-commit: a848e67
verdict: FINDINGS
---
# Re-review round 4 — hermes-1

## Verdict

FINDINGS — one MINOR. Every round-03 finding is verified FIXED at the real call
site, proven behaviorally through the production functions, not inferred from the
diff. All three new cycle-4 tests are effective — each FAILS when its fix is
reverted, confirmed by running the tests against overlaid reverted source. No
regressions introduced by cycle 4. Build clean, full suite green, vet clean.

The one finding is the same skill-release hygiene issue kimi-1 carried open as
N-2 through rounds 02 and 03: the IMPLEMENTATION.md cycle-4 section claims
"Skill 2.5.1 is released with this cycle rather than left uncommitted," but the
skill repo's HEAD is still b806ada (v2.5.0), package.json still reads 2.5.0,
SKILL.md and references/COOPERATION.md are still uncommitted working-tree edits
(` M`), and no v2.5.1 tag exists. The content is correct on all four surfaces;
the release step did not happen.

Evidence provenance: PRIMARY = I read the source, ran the tool, or executed the
test this session against scratch decks under /tmp with PARLEY_HOME isolation.
I built the binary at a848e67 (`go build -o /tmp/parley-r4 ./cmd/parley/`) and
ran `roster show`, `roster show --explain`, and `roster show --all` against
isolated fixtures. For paths the CLI surface cannot print without launching
(participant selection, the continuation boundary), I injected probe tests via
`go test -overlay` from /tmp — the reviewed tree stays byte-identical; the
probes drive the same production functions (`RosterMembership`,
`defaultRosterParticipants`, `applyRosterSnapshotToParticipants` +
`agents.ResolveParticipant`). Both repositories are READ-ONLY per the hard
constraint; the only file written is this review.

Build and suite (exact output):

```
$ go build ./... && echo BUILD_OK
BUILD_OK

$ go test -count=1 ./...
ok  parley-deck-cli/internal/acp         0.543s
ok  parley-deck-cli/internal/agents       0.308s
ok  parley-deck-cli/internal/app         23.036s
ok  parley-deck-cli/internal/config       0.759s
ok  parley-deck-cli/internal/consensus    0.997s
ok  parley-deck-cli/internal/driver       3.984s
ok  parley-deck-cli/internal/fsutil       1.442s
ok  parley-deck-cli/internal/hitl         1.905s
ok  parley-deck-cli/internal/loop         1.673s
ok  parley-deck-cli/internal/pipeline     2.113s
ok  parley-deck-cli/internal/procctl      2.589s
ok  parley-deck-cli/internal/protocol     2.490s
ok  parley-deck-cli/internal/repomap      2.600s
ok  parley-deck-cli/internal/retro        2.609s
ok  parley-deck-cli/internal/runaction    2.663s
ok  parley-deck-cli/internal/runcontrol   2.687s
ok  parley-deck-cli/internal/runmanifest  2.691s
ok  parley-deck-cli/internal/runner       9.808s
ok  parley-deck-cli/internal/runplan      2.714s
ok  parley-deck-cli/internal/runstate     2.717s
ok  parley-deck-cli/internal/sessionstore 2.574s
ok  parley-deck-cli/internal/steer        2.717s
ok  parley-deck-cli/internal/store        2.800s
ok  parley-deck-cli/internal/track        2.789s
ok  parley-deck-cli/internal/tui          2.650s
TEST_EXIT=0

$ go vet ./... && echo VET_OK
VET_OK
```

All 26 packages pass, 0 failures. `internal/runner`'s
`TestDurableKillEndToEndRealProcess` — which failed in codex-1's round-03
sandbox — passes on this machine. No cycle-4 change touches `internal/runner`.

## Round-03 findings: fixed or not

### [CRITICAL codex-1 R3-C1] value-only layers could retire or revive a committed deck member — FIXED (PRIMARY)

The round-03 CRITICAL: `LoadRosterScoped` gated which layers may add membership
IDs but left `active` merging from every layer via `mergeRosterEntry`
(runtime.go:251-253), so `[roster.claude-1] active = false` in the gitignored
`agents.local.toml`, `$PARLEY_HEADLESS_AGENT_CONFIG`, or the machine file could
silently drop a committed member from the quorum, or revive one the deck retired.

PRIMARY — The fix is `RosterScope.applyAuthorityState` (runtime.go:209-221),
called after all value layering completes. It iterates `r.Members` and forces
each member's `Active` to the value the authority layer declared — deck, legacy
§2, or machine when inherited — discarding any `active` a value-only layer
merged in. The three authority branches each call it before returning
(runtime.go:184, 199, 204).

PRIMARY — I verified this behaviorally with conflicting `active` values across
ALL FOUR layers, which the round-03 finding explicitly requested:

  - Deck: `[roster.claude-1] active = true` + `[roster.kimi-1] active = true`
  - Machine: `[roster.claude-1] active = false` + `[roster.codex-1] active = true`
  - agents.local.toml: `[roster.kimi-1] active = false`
  - $PARLEY_HEADLESS_AGENT_CONFIG: `[roster.claude-1] active = false`

CLI surface (`roster show --json` with all four layers active): claude-1 and
kimi-1 both `state=active`. No retirement leaked from any value-only layer.

PRIMARY — Participant selection, which the CLI cannot print without launching,
was driven through an overlay probe calling the real `RosterMembership(root)`:
both claude-1 and kimi-1 are in the `active` set, not `inactive`. A second probe
calling `defaultRosterParticipants(root, discovered)` returns
`ids=[claude-1 kimi-1]` — the run quorum is the deck's, not the machine's or
the local layer's.

PRIMARY — I also verified the REVERSE direction: deck retires claude-1
(`active = false`), local layer tries to revive it (`active = true`). The
overlay probe confirms claude-1 stays `inactive` — the local layer cannot
revive a member the deck retired.

PRIMARY — The shipped test `TestValueLayersCannotChangeMembershipState`
(roster_cycle2_test.go:241-258) covers the retire direction across three of
four layers (deck, machine, local) but NOT the env config layer. My probe
covers all four including env. The env layer follows the same code path as
local (both are `!item.membership` layers skipped at runtime.go:150), so the
omission is a test-coverage gap, not a behavioral risk.

### [MAJOR codex-1 R3-M1 / kimi-1 N-1] cycle 3's regression: bare-family participant lost its freeze — FIXED (PRIMARY)

The round-03 MAJOR: `applyRosterSnapshotToParticipants` appended the frozen
discovery instead of replacing it. The runner resolves by FIRST exact-ID match
(resolve.go:47-53), so for a participant spelled exactly like its adapter family
(`claude`, not `claude-1`), the live unfrozen record won — the freeze existed in
memory but never reached the launch.

PRIMARY — The fix replaces in place (roster_snapshot_apply.go:98-113): after
freezing, it scans `out` for a discovery with the same `Spec.ID` and overwrites
it. If no match is found (a new participant), it appends. This is the correct
fix — the frozen record now occupies the position the runner's first-match rule
will find.

PRIMARY — I verified this through the real resolution path. An overlay probe
fed `applyRosterSnapshotToParticipants` a bare-family discovery
(`{ID: "claude", Model: "drifted"}`) with a snapshot pinning `claude` to
`frozen-model`, then resolved through `agents.ResolveParticipant("claude", out,
map{})` — the exact function the runner calls at runner.go:366. Result:
`picked.Spec.Model = "frozen-model"`. The runner would launch the frozen record.

PRIMARY — I also verified the composite-ID path still works (no regression from
the replace-in-place fix): two participants `claude-1` and `claude-2` sharing
the `claude` adapter, each with a different frozen model. Both resolve to their
own frozen model — `claude-1 → frozen-1`, `claude-2 → frozen-2`.

### [MAJOR codex-1 R3-M2] --scope machine leaked through --all and misattributed provenance — FIXED (PRIMARY)

The round-03 MAJOR had two parts:

(a) `--scope machine --all` used unscoped `LoadAgentSpecs` and
`LoadRosterAdapters` in `unrosteredRows`, so deck/local/env values leaked into
machine-scope unrostered rows.

PRIMARY — `unrosteredRows` now takes a `machineOnly bool` parameter
(roster_view.go:98) and calls `LoadAgentSpecsScoped` and
`LoadRosterAdaptersScoped` (roster_view.go:100, 106). The single call site at
roster.go:437 passes `opts.scope == "machine"`. I verified behaviorally: a deck
setting `model = "deck-only-model"` for claude-1 and a local layer setting
`model = "local-leaked-model"` for the claude agent block — `--scope machine
--all --json` shows neither leaked value in the unrostered rows. Only machine-
scoped values appear.

(b) `--explain` reported "built-in default" for values that reach the launch
through `[agents.<family>]` blocks.

PRIMARY — `rosterExplain` now calls `config.AgentFieldSources(root,
row.Adapter, opts.scope == "machine")` (roster_view.go:175) and falls through
to it when `RosterFieldSourcesScoped` returns no source for a field
(roster_view.go:178-180). I verified behaviorally: a machine file with
`[agents.claude] model = "machine-claude-model" reasoning = "high"` —
`roster show --scope machine --explain claude-1` reports `model
machine-claude-model SET BY ~/.parley/agents.toml` and `effort high SET BY
~/.parley/agents.toml`, not "built-in default". Speed (not set in the agent
block) correctly falls through to "built-in default".

### [MAJOR codex-1 R3-M3] the four-surface guard overclaimed — FIXED (PRIMARY)

The round-03 MAJOR: the drift guard's `t.Logf` said "checked N of 4 surfaces"
without stating that the unread surfaces were NOT enforced, implying
four-surface coverage when only two were readable.

PRIMARY — The log message now reads: "NOTE: checked N of 4 surfaces — the
bundled skill protocol and SKILL.md live in the sibling parley-deck-skill
checkout and were not readable here; they are NOT enforced by this run"
(drift_test.go:289-292). The message explicitly states the surfaces are not
enforced. I ran the test verbosely: with the sibling checkout present (as in
this workspace), all 4 surfaces are read and no NOTE appears. Without the
sibling checkout, the NOTE would fire — it no longer overclaims.

### [MINOR codex-1 R3-m1] the JSON contract test could not catch a reverted fix — FIXED (PRIMARY)

The round-03 MINOR: `TestJSONStatusMatchesTextForAHealthyRow` only asserted
`status` was non-null and checked the column count, so re-adding `display_name`
or `note` to the JSON output would not fail it.

PRIMARY — The new test `TestJSONRowHasExactlyTheFrozenColumns`
(roster_cycle2_test.go:262-296) asserts the exact eleven-key set in BOTH
directions: every key in the row must be in the expected set (no extras), and
every expected key must be present in the row (no missing). I verified
behaviorally: `roster show --json` returns exactly 11 keys per row — `agent,
adapter, state, installed, model, model_family, model_company, effort, speed,
autonomous, status` — with no `display_name` or `note`.

### [MINOR both] skill 2.5.1 unshipped — NOT FIXED (PRIMARY)

PRIMARY — The IMPLEMENTATION.md cycle-4 section (line 165) states: "Skill 2.5.1
is released with this cycle rather than left uncommitted." This is false. I
verified in the skill repo:

  - HEAD: b806ada (`[claude-1] roster-operations-standard: skill 2.5.0`)
  - `git status --short`: ` M skills/parley-deck/SKILL.md`, ` M
    skills/parley-deck/references/COOPERATION.md` — still uncommitted
  - `package.json`: `"version": "2.5.0"` — not bumped
  - `git tag -l 'v2.5*'`: only `v2.5.0` — no v2.5.1 tag

The content is correct on all four surfaces (grep for all banned phrases
returns zero hits in the working tree), but the release step — commit, version
bump, tag — did not happen. This is the same state kimi-1's round-03 N-2
described. The hard constraint makes both repos READ-ONLY, so this cannot be
fixed in this review. See New findings below.

## New findings (by severity, or "none")

### MINOR — R4-m1: IMPLEMENTATION.md claims skill 2.5.1 is released; it is not

PRIMARY — The cycle-4 IMPLEMENTATION.md section (line 165) says "Skill 2.5.1 is
released with this cycle rather than left uncommitted." The skill repo
contradicts this: HEAD is b806ada (v2.5.0), package.json is 2.5.0, the corrected
SKILL.md and references/COOPERATION.md are uncommitted working-tree edits, and
no v2.5.1 tag exists. The content is verified correct on all four surfaces; the
release did not happen. This is a documentation-accuracy finding: the
IMPLEMENTATION.md claims a release that was not performed. The fix is to either
commit and tag the skill repo as 2.5.1, or correct the IMPLEMENTATION.md claim
to say the skill edits remain uncommitted and the release is pending.

This is the same issue kimi-1 carried open as N-2 through rounds 02 and 03. It
is not chargeable to commit a848e67's code; it is chargeable to the exit
checklist. If Phase 8 closes here, the uncommitted skill edits must be committed
and both versions cut as part of that close, or the content this cycle verified
never reaches an installed skill.

## Test-quality assessment

PRIMARY — I ran all three new cycle-4 tests individually (all PASS), then
verified each would FAIL if its fix were reverted by running the test against
overlaid reverted source. No tautologies found.

1. `TestSnapshotFreezeReachesBareFamilyParticipants` (R3-M1) — REAL. Calls
   `applyRosterSnapshotToParticipants` with a bare-family discovery
   (`ID="claude"`, `Model="drifted"`) and a snapshot pinning `claude` to
   `frozen-model`, then resolves through `agents.ResolveParticipant("claude",
   out, map{})` — the exact function the runner calls at runner.go:366. Asserts
   `picked.Spec.Model == "frozen-model"`. REVERT VERIFIED: replacing the
   replace-in-place block with append-only (the cycle-3 bug) makes the test
   FAIL with `runner would launch "drifted", not the frozen "frozen-model"`.
   Not tautological — it observes the record the runner would actually launch,
   not just the helper's output slice.

2. `TestValueLayersCannotChangeMembershipState` (R3-C1) — REAL. Creates a deck
   with two active members, a machine file retiring claude-1, and a local layer
   retiring kimi-1. Calls `RosterMembership(root)` and asserts both stay
   active. REVERT VERIFIED: removing `applyAuthorityState` calls from
   `LoadRosterScoped` makes the test FAIL with both members retired by the
   value-only layers. Not tautological. Coverage gap: does not test the env
   config layer (`$PARLEY_HEADLESS_AGENT_CONFIG`); my overlay probe covers all
   four layers including env, and the env layer follows the same code path as
   local (both are `!item.membership` at runtime.go:150).

3. `TestJSONRowHasExactlyTheFrozenColumns` (R3-m1) — REAL. Calls `rosterShow`
   with `--json`, unmarshals the roster array, and asserts each row carries
   exactly the eleven frozen keys — no extras, no missing. REVERT VERIFIED:
   changing `json:"-"` back to `json:"display_name"` and `json:"note"` on the
   `rosterRow` struct makes the test FAIL with `row carries "display_name"` and
   `row carries "note"`. Not tautological — the bidirectional key-set check
   catches both additions and removals.

4. Drift test guard note (R3-M3) — not a new test, but a message clarity fix.
   The `t.Logf` now explicitly states unread surfaces are "NOT enforced by this
   run." This is a documentation improvement to the test's output, not a new
   assertion. The guard's functional behavior (check readable surfaces, fail on
   banned phrases) is unchanged. No tautology risk.

Regression sweep: `go build ./...` clean, `go test -count=1 ./...` all 26
packages ok, `go vet ./...` clean. The `unrosteredRows` signature change (adding
`machineOnly bool`) has exactly one call site (roster.go:437), which was updated
correctly. The `AgentFieldSources` function (runtime.go:1135-1168) correctly
maps `ov.Reasoning` → `"effort"` and handles missing `[agents.<family>]` blocks
by returning an empty map (no panic, no misattribution). The
`applyAuthorityState` function (runtime.go:209-221) correctly handles the
`!declared` case (`e.Active = !declared || state` — if a member ID is not in
the authority active map, it defaults to active). No cycle-4 change touches
`internal/runner`. The CLI repo is clean (`git status --short` empty).
