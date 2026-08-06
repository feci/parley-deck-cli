---
idea: roster-operations-standard
phase: 8 — re-review
agent: hermes-1
round: 3
date: 2026-08-06
reviewed-commit: 7220715
verdict: CLEAN
---
# Re-review round 3 — hermes-1

## Verdict

CLEAN. Every round-02 finding I or another reviewer raised is now actually fixed
at the real call site, verified by PRIMARY reproduction. The build is clean, the
full test suite is green (24 packages, 0 failures), and `go vet` is clean. No new
defects were introduced by cycle 3.

Evidence provenance: PRIMARY = I read the source, ran the tool, or executed the
test this session against scratch decks under /tmp with PARLEY_HOME isolation. I
built the binary at 7220715 (`go build -o /tmp/parley-rev3 ./cmd/parley/`) and ran
`roster show`, `roster set`, `roster init`, and `roster show --explain` against
isolated fixtures. Nothing was written to the real deck or to
~/.parley/agents.toml. Both repositories are READ-ONLY per the hard constraint;
the only file written is this review.

Build and suite (exact output):

```
$ go build ./... && echo BUILD_OK
BUILD_OK

$ go test ./...
ok  parley-deck-cli/internal/app          19.628s
ok  parley-deck-cli/internal/runner        8.398s
... (all 24 packages ok)
TEST_EXIT=0

$ go vet ./...
VET_EXIT=0
```

Note: `internal/runner`'s `TestDurableKillEndToEndRealProcess` — which failed in
codex-1's round-02 environment — passes on this machine. No cycle-3 change
touches `internal/runner`.

## Round-02 findings: fixed or not

### CRITICAL: A1 legacy-fallback clarification — FIXED (PRIMARY, all three cases)

The round-02 CRITICAL (all three reviewers): a deck whose only roster is a valid
legacy §2 table had the machine roster inherited over it, including the run
quorum. `LoadRosterScoped` knew nothing about §2.

`LoadRosterScoped` (runtime.go:164-192) now decides authority in an explicit
order before any value layering:

  1. committed deck blocks (`deckMembers` non-empty) → return immediately
  2. else a valid legacy §2 table (`protocol.ReadRosterIDs(root)`) → `out.Legacy
     = true`, `out.Members` = §2 IDs, return
  3. else the machine roster (`out.Inherited = true`)

PRIMARY — I proved all three cases behaviorally with scratch fixtures:

CASE 1 (committed deck blocks block the machine roster):
  Deck declares [roster.claude-1] + [roster.kimi-1]; machine has 5 members.
  `roster show --json` returned exactly 2 agents: claude-1, kimi-1. No
  machine-only agents (codex-1, hermes-1, opencode-1) leaked. No legacy-roster
  or inherited-roster status. Correct.

CASE 2 (legacy §2 only, no deck agents.toml, machine has 5) — THE CRITICAL FIX:
  Deck has a valid §2 table (header `| Agent ID | Workspace dir |`) declaring
  claude-1 + kimi-1; no agents.toml in the deck; machine has 5 members.
  `roster show --json` returned exactly 2 agents: claude-1, kimi-1. Both with
  status `legacy-roster`. No machine-only agents leaked. No inherited-roster
  status. Correct — this is the exact scenario that was broken in cycle 2.

  (Note: my first attempt used an invalid §2 header format — the parser requires
  `| Agent ID | Workspace dir |` per `rosterHeaderRe` in protocol/roster.go:19.
  Once I used the correct format, the fix worked. The cycle-2 failure was real;
  the cycle-3 fix is real.)

CASE 3 (neither deck blocks nor §2 — truly rosterless):
  Deck has COOPERATION.md with no §2 table and no agents.toml; machine has 5.
  `roster show --json` returned 5 agents, all with `inherited-roster`. Correct.

Participant selection: `defaultRosterParticipants` (app.go:2453) calls
`RosterMembership` (roster.go:719), which calls `LoadRosterScoped` — the same
function `roster show` uses. For the legacy §2 case, `LoadRosterScoped` returns
the §2 IDs in `scope.Members`, so `RosterMembership` returns exactly those IDs.
The test `TestLegacySection2BeatsTheMachineRoster` (roster_cycle2_test.go:159-186)
verifies both `RosterMembership` AND `resolveRoster` for this case. PRIMARY — I
read the test and traced the code path.

### codex-1: membership pooled every non-machine layer — FIXED (PRIMARY)

codex-1 noted that `agents.local.toml` and `$PARLEY_HEADLESS_AGENT_CONFIG` were
pooled as `deckMembers`, so they could add members.

`configLayers` (runtime.go:355-365) now marks only `parley-deck/agents.toml`
with `membership: true`. In `LoadRosterScoped` (runtime.go:142-156), layers with
`item.machine || !item.membership` are skipped for membership — their IDs go to
`machineMembers` (if machine) or are ignored entirely (if local/env). Only the
`membership: true` layer populates `deckMembers`.

PRIMARY — I proved this behaviorally:
  - Deck with [roster.claude-1] + [roster.kimi-1] + agents.local.toml declaring
    [roster.sneaky-local-1]: `roster show` returned 2 agents. sneaky-local-1 did
    NOT appear. Correct.
  - Same deck + `$PARLEY_HEADLESS_AGENT_CONFIG` pointing to a file declaring
    [roster.env-sneaky-1]: `roster show` returned 2 agents. env-sneaky-1 did NOT
    appear. Correct.

### MAJOR A3: `roster init --yes` bypassed the membership gate — FIXED (PRIMARY)

`rosterInit` (roster.go:474) now receives `confirmBreaking` (roster.go:167) and
gates on it: `--yes` without `--confirm-breaking` exits 2 with "this is a
membership change — re-run with --confirm-breaking as well as --yes." (roster.go:554-557,
601-604). No file is written.

PRIMARY — verified:
```
$ parley roster init --scope deck --dir <case> --yes
roster init (deck) will add 1 roster member(s) to parley-deck/agents.toml.
this is a membership change — re-run with --confirm-breaking as well as --yes.
EXIT=2
# agents.toml NOT written

$ parley roster init --scope deck --dir <case> --yes --confirm-breaking
Wrote 1 mapping(s) to parley-deck/agents.toml. The driver can now run this roster.
EXIT=0
```

Two existing tests were updated to add `--confirm-breaking` to their `init --yes`
calls (roster_test.go:192, 258, 290). The new test
`TestRosterInitRequiresConfirmBreaking` (roster_cycle2_test.go:197-218) verifies
both the refusal (exit 2, no file) and the success (exit 0, file written).

### MAJOR A5: per-ID snapshot pins collapsed at the real call site — FIXED (PRIMARY)

`continueAuto` (app.go:1184) now calls `applyRosterSnapshotToParticipants`
instead of `applyRosterSnapshot`. The new function
(roster_snapshot_apply.go:83-100) resolves each participant to its roster ID via
`agents.ResolveParticipant` (which uses the mapping to find the adapter-level
discovery and rewrites `Spec.ID` to the roster ID), then applies
`applyRosterSnapshot` to each resolved discovery individually. This means
`frozen[e.Agent]` hits the roster ID, not the adapter fallback.

PRIMARY — code-read: the resolution path is correct. `ResolveParticipant`
(resolve.go:56-64) uses `mapping[participant]` to find the family, matches the
discovery with that family ID, and sets `d.Spec.ID = participant` (the roster
ID). So `applyRosterSnapshot` receives a discovery whose `Spec.ID` is the roster
ID, which matches `frozen[e.Agent]`.

The new test `TestSnapshotPinsSurviveParticipantResolution`
(roster_cycle2_test.go:123-143) proves the boundary: it uses adapter-level
discoveries (ID="claude") and two participants ("claude-1", "claude-2") sharing
that adapter with different frozen models, and asserts each gets its own model.
This is the exact scenario that was broken — the test would fail if
`applyRosterSnapshotToParticipants` were replaced with `applyRosterSnapshot`.

### MAJOR A6: machine scope still read deck values/provenance — FIXED (PRIMARY)

`resolveRoster` (roster.go:249) now calls `config.LoadAgentSpecsScoped(root,
opts.scope == "machine")` which restricts the layer stack to the machine layer
only (via `machineOnlyLayers`, runtime.go). Adapter mappings use
`config.LoadRosterAdaptersScoped` (roster.go:265). `rosterExplain`
(roster_view.go:165) uses `config.RosterFieldSourcesScoped`.

PRIMARY — verified:
  Deck declares [roster.claude-1] with `model = "deck-only-model"`; machine has
  `model = "claude-opus-5"`.
```
$ parley roster show --dir <case> --scope machine --json
  claude-1 model: claude-opus-5  (NOT deck-only-model)

$ parley roster show --dir <case> --scope machine --explain claude-1
  model    claude-opus-5    ~/.parley/agents.toml  (NOT "built-in default" with deck-only-model)
```

### MAJOR A7: display_name+note were JSON-only fields outside the frozen 11 columns — FIXED (PRIMARY)

`rosterRow` (roster.go:213-214) now has `Display string \`json:"-"\`` and
`Note string \`json:"-"\``. The `json:"-"` tag suppresses serialization entirely.

PRIMARY — verified:
```
$ parley roster show --dir <case> --json | jq '.roster[0] | keys'
  ["adapter","agent","autonomous","effort","installed","model",
   "model_company","model_family","speed","state","status"]
```
11 keys, no `display_name`, no `note`. The frozen 11-column contract is now
honored in both text and JSON.

### MAJOR A9: drift guard covered 2 of 4 surfaces — FIXED (PRIMARY)

`TestNoSection2AsAStoreInstructions` (drift_test.go:255-299) now reads all four
surfaces: embedded default, live deck, bundled skill protocol, and SKILL.md. The
last two are read from the sibling skill repo checkout
(`../../../parley-deck-skill/...`), with a `t.Logf` skip when the checkout is
absent — so the guard never silently passes on a surface it did not read.

The banned-phrase list gained `"roster sync` moves it over"` and
`"roster sync` moves it across"` (drift_test.go:263-264).

PRIMARY — ran the test: `PASS` (0.298s). The skill repo is present as a sibling
checkout, so all 4 surfaces are checked.

### MAJOR A10: RosterRevisionOf omitted LaunchArgs — FIXED (PRIMARY)

`RosterRevisionOf` (manifest.go:94-96) now hashes
`strings.Join(e.LaunchArgs, "\x1f")` as part of the revision. The `\x1f` (unit
separator) delimiter prevents ambiguity between args that could concatenate to
the same string.

The test `TestRosterRevisionCoversLaunchArgs` (roster_cycle2_test.go:149-155)
verifies two entries with identical Agent/Adapter/Auto but different LaunchArgs
produce different revisions. Would fail if the `LaunchArgs` hash were removed.

### MAJOR (hermes-1): stale "roster sync moves it over" claim in all three COOPERATION.md copies — FIXED (PRIMARY)

PRIMARY — checked all four copies (three COOPERATION.md + SKILL.md):
  - `parley-deck/COOPERATION.md`: CLEAN
  - `internal/protocol/defaults/COOPERATION.md`: CLEAN
  - `parley-deck-skill/.../references/COOPERATION.md`: CLEAN
  - `parley-deck-skill/.../SKILL.md`: CLEAN

All four now contain the correct guidance: "`roster sync` does NOT migrate it"
and name `roster migrate` or `roster set` + `roster render` as the remediation.

### MINOR A12: masked-by-env false-positived on every machine-scope write — FIXED (PRIMARY)

`rosterFieldMaskedBy` (roster_set.go:105-116) now resolves the display label to
a path via `config.RosterSourcePath(root, src)` and compares absolute paths
(`filepath.Abs` on both sides). If the resolved path matches the target, no
warning is emitted.

PRIMARY — verified:
```
$ parley roster set claude-1 --scope machine --model test-model-xyz --yes
  (no MASKED warning on stderr)
  effective model after write: test-model-xyz
```

The test `TestMachineScopeWriteIsNotReportedAsMasked`
(roster_cycle2_test.go:189-195) verifies `rosterFieldMaskedBy` returns
`masked=false` for a machine-scope write.

### MINOR A16: residuals — FIXED (PRIMARY, all three sub-items)

1. Unmapped guidance: `roster.go:324` now reads
   `parley roster set %s --scope deck --adapter <family> --yes --confirm-breaking`.
   PRIMARY — verified: `roster show` on a deck with an unmapped agent that
   matches a family prefix prints the note with `--confirm-breaking`.

2. Broad `k*` classification: the `{"k", ...}` prefix rule is removed from
   `prefixRules` (modelmeta.go). `kimiCodename` (modelmeta.go:133-136) matches
   only `k` followed by a digit (k2, k3, k2-0711). PRIMARY — verified the logic:
   `k2` matches, `klms` does not, `kimi-code` does not, `kafka` does not.

3. kimi-1 attribution: `CHANGELOG.md:7` now reads "found by codex-1, hermes-1
   and kimi-1 — all three corroborated both CRITICALs independently." PRIMARY.

### NIT: IMPLEMENTATION.md accuracy — ADDRESSED (PRIMARY)

The cycle-2 text (IMPLEMENTATION.md:75-79) already names the
`internal/runner` test as environment-dependent. The cycle-3 NIT section
(line 126-127) references this. The cycle-2 test count says "Ten new tests" —
kimi-1's NIT-3 counted 11. This is a minor count discrepancy in the cycle-2
text that was not corrected, but the cycle-3 section's own claim of "Six more"
tests is accurate (5 new test functions + 1 extended drift test).

## New findings (by severity, or "none")

None.

One observation below NIT threshold that I do not raise as a finding: the
`rosterUsage` string (roster.go:63) does not list `--confirm-breaking` for the
`init` subcommand, though the flag is registered and functional. The flag help
output does show it. This is a help-text completeness gap, not a functional
defect — the A3 gate works correctly as verified above.

The skill repo's working-tree changes (SKILL.md, references/COOPERATION.md)
remain uncommitted with package.json still at 2.5.0. This is the same state
codex-1 and kimi-1 noted in round 02. The cycle-3 IMPLEMENTATION.md does not
claim to have committed or versioned the skill; it claims to have fixed the
stale COOPERATION.md text and extended the drift guard, both of which are
verified done. The drift test reads from disk (working tree), so it passes with
the current correct content. If the skill changes were discarded, the drift
test would fail — but that is a skill-repo commit decision, not a cycle-3 code
fix claim. I do not raise this as a finding because it was not a cycle-3 fix
claim, and the hard constraint makes both repositories read-only.

## Test-quality assessment

Six new tests added in cycle 3 (roster_cycle2_test.go:123-218) plus the
extended drift test. I assessed each for whether it tests behavior or encodes
an assumption, and whether it would fail if the fix were reverted.

1. `TestSnapshotPinsSurviveParticipantResolution` (A5) — REAL. Uses
   adapter-level discoveries (ID="claude") and two participants sharing that
   adapter with different frozen models. Asserts each gets its own model. If
   the fix were reverted (replacing `applyRosterSnapshotToParticipants` with
   `applyRosterSnapshot`), the adapter-level discovery would take the adapter
   fallback and both participants would get the same model — the test would
   fail. Not tautological.

2. `TestRosterRevisionCoversLaunchArgs` (A10) — REAL. Two entries with
   identical fields except LaunchArgs. Asserts different revisions. If the
   LaunchArgs hash were removed, the revisions would be equal — the test would
   fail. Not tautological.

3. `TestLegacySection2BeatsTheMachineRoster` (A1) — REAL. Creates a legacy §2
   deck on a 5-member machine, calls both `RosterMembership` and
   `resolveRoster`. Asserts exactly 2 members (not 5), no `inherited-roster`,
   all rows have `legacy-roster`. If the §2 authority check in
   `LoadRosterScoped` were removed, `scope.Members` would come from
   `machineMembers` (5), and the test would fail on the count assertion. Not
   tautological. This is the test that proves the CRITICAL fix.

4. `TestMachineScopeWriteIsNotReportedAsMasked` (A12) — REAL. Calls
   `rosterFieldMaskedBy` with the machine target path. Asserts `masked=false`.
   If the fix were reverted (comparing display label against absolute path),
   the label `"~/.parley/agents.toml"` would not match the absolute target —
   `masked=true` — the test would fail. Not tautological.

5. `TestRosterInitRequiresConfirmBreaking` (A3) — REAL. Calls `runRoster`
   with `init --yes` (no `--confirm-breaking`), asserts exit 2 and no file
   written. Then with `--confirm-breaking`, asserts exit 0. If the gate were
   removed, the first call would exit 0 and write the file — the test would
   fail. Not tautological. Tests the actual CLI flag plumbing, not just the
   helper.

6. Extended `TestNoSection2AsAStoreInstructions` (A9) — REAL. Now reads 4
   surfaces (2 in-repo + 2 sibling skill repo files). Added 2 new banned
   phrases. If a stale phrase returned to any checked surface, the test would
   fail with `t.Errorf`. The sibling-skip is a `t.Logf`, not a silent pass —
   the guard never silently passes on a surface it did not read. Not
   tautological.

No tautologies found. Every new test constructs a fixture, calls the function
or command under test, and asserts a specific outcome that would fail if the
behavior regressed. The tests cover the exact scenarios that were broken in
cycle 2: the legacy §2 + machine roster combination (A1), the continuation
boundary with adapter-keyed discoveries (A5), the machine-scope self-masking
(A12), and the init gate bypass (A3).

Regression check: `go build ./...` clean, `go test ./...` all 24 packages ok,
`go vet ./...` clean. The `internal/runner` test that failed in codex-1's
round-02 environment passes here. No cycle-3 change touches `internal/runner`.
The `resolveRoster` legacy-path refactor (roster.go:277-301) correctly handles
inactive agents in the legacy case via `scope.Entries[id].Active`. The
`applyRosterSnapshotToParticipants` function preserves the original discoveries
and appends frozen participant discoveries, so the runner's resolver matches
the frozen entries by exact ID (rule 1) without ambiguity.
