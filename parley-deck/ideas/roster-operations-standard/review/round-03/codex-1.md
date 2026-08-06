---
idea: roster-operations-standard
phase: 8 — re-review
agent: codex-1
round: 3
date: 2026-08-06
reviewed-commit: 7220715
verdict: FINDINGS
---
# Re-review round 3 — codex-1

## Verdict

PRIMARY — FINDINGS. Cycle 3 fixes the three-way ID-set precedence, the `roster init` gate, `LaunchArgs` hashing, machine-write masking, the frozen JSON shape, and the A16 residuals. It does not fully fix quorum authority, continuation pinning, machine scope, or the four-surface protocol guard.

PRIMARY — I own the new findings and recommendations below, so I issue no §15.1 truth verdict on them. Labels such as FIXED, PARTIAL, and NOT FIXED adjudicate cycle-3 implementation claims that I do not own.

PRIMARY — The exact required command, `go build ./... && go test ./...`, exited 1. `go build ./...` succeeded silently; the test run failed only in the pre-existing durable-kill test. `git diff --name-only 57fe9d7..7220715 -- internal/runner` produced no path, so cycle 3 did not modify that package. The exact combined output was:

```text
?   	parley-deck-cli/cmd/parley	[no test files]
ok  	parley-deck-cli/internal/acp	(cached)
ok  	parley-deck-cli/internal/agents	(cached)
ok  	parley-deck-cli/internal/app	35.049s
ok  	parley-deck-cli/internal/config	(cached)
ok  	parley-deck-cli/internal/consensus	(cached)
ok  	parley-deck-cli/internal/driver	1.140s
ok  	parley-deck-cli/internal/fsutil	(cached)
ok  	parley-deck-cli/internal/hitl	(cached)
ok  	parley-deck-cli/internal/loop	(cached)
ok  	parley-deck-cli/internal/pipeline	(cached)
ok  	parley-deck-cli/internal/procctl	0.307s
ok  	parley-deck-cli/internal/protocol	(cached)
ok  	parley-deck-cli/internal/repomap	(cached)
ok  	parley-deck-cli/internal/retro	(cached)
ok  	parley-deck-cli/internal/runaction	(cached)
ok  	parley-deck-cli/internal/runcontrol	(cached)
ok  	parley-deck-cli/internal/runmanifest	(cached)
--- FAIL: TestDurableKillEndToEndRealProcess (0.02s)
    durablekill_test.go:116: a live attributed process should be killed, got {AgentID:sleeper Killed:false Cleared:false Failed:true SegmentID:segment-0001 Message:process verification failed (no recorded boot id); not killed}
FAIL
FAIL	parley-deck-cli/internal/runner	7.316s
ok  	parley-deck-cli/internal/runplan	(cached)
ok  	parley-deck-cli/internal/runstate	(cached)
ok  	parley-deck-cli/internal/sessionstore	(cached)
ok  	parley-deck-cli/internal/steer	(cached)
ok  	parley-deck-cli/internal/store	(cached)
ok  	parley-deck-cli/internal/track	(cached)
ok  	parley-deck-cli/internal/tui	0.375s
FAIL
```

## Round-02 findings: fixed or not

### A1 — legacy fallback and membership authority: PARTIAL

PRIMARY — The requested three-case authority order is fixed for roster IDs:

- PRIMARY — Committed deck blocks win. `PARLEY_HEADLESS_AGENT_CONFIG='/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/wt-editor-composer/parley-deck/agents.toml' go run ./cmd/parley roster show --dir . --json` used an environment config containing an extra `antigravity-1`, yet its roster array contained only `claude-1`, `codex-1`, `hermes-1`, and `kimi-1`. `TestDeckMembershipIsTheDeckFileNotTheLayeredUnion` (`internal/app/roster_membership_test.go:36-69`) independently constructs a two-member deck over a five-member machine roster and gets exactly the two deck IDs from `RosterMembership`; `resolveRoster`, which backs `roster show`, gets the same authority.

- PRIMARY — A valid legacy §2 table wins when committed deck blocks are absent. `TestLegacySection2BeatsTheMachineRoster` (`internal/app/roster_cycle2_test.go:159-185`) constructs a two-member legacy table over a five-member machine roster; both `RosterMembership` and the production show resolver return only the two legacy IDs, and the rows are `legacy-roster`, not `inherited-roster`.

- PRIMARY — The machine roster wins only when neither deck authority exists. `go run ./cmd/parley roster show --dir /private/tmp/parley-r3-neither-does-not-exist --json` returned `claude-1`, `codex-1`, `hermes-1`, `kimi-1`, and `opencode-1`, all marked `inherited-roster`. Participant selection consumes the same `RosterMembership` result at `internal/app/app.go:2454`; the no-roster branch is also exercised by `TestDefaultRosterParticipants` case 3 (`internal/app/roster_test.go:91-101`) when all three authorities are empty.

PRIMARY — `roster show` reaches `resolveRoster` directly at `internal/app/roster.go:428`, while default participant selection reaches `RosterMembership` at `internal/app/app.go:2454`. The behavioral tests therefore cover the two production authority consumers rather than a second test-only implementation.

PRIMARY — Non-machine layers can no longer add IDs: `LoadRosterScoped` records only the committed deck layer as deck membership and skips local/env IDs at `internal/config/runtime.go:139-150`. This fixes the round-02 union finding for ID membership.

PRIMARY — Quorum authority is nevertheless still layered through `active`; see new CRITICAL finding R3-C1. Thus the narrower claim “non-machine layers cannot add IDs” is fixed, but the larger claim that only the committed deck controls membership/state is not.

### A3 — `roster init --yes` membership gate: FIXED

PRIMARY — `rosterInit` now receives `confirmBreaking` and refuses an applying membership change without it. `TestRosterInitRequiresConfirmBreaking` (`internal/app/roster_cycle2_test.go:199-217`) exercises the full command handler and would fail if that gate were removed; the targeted test passed.

### A5 — per-ID snapshot pins at continuation: NOT FIXED

PRIMARY — Cycle 3 moves participant resolution before snapshot application at the `continueAuto` call site (`internal/app/app.go:1184`), so two composite IDs sharing one adapter now receive distinct appended frozen discoveries. That repairs the exact composite-ID example from round 2.

PRIMARY — The new helper preserves the original discovery slice before appending frozen discoveries (`internal/app/roster_snapshot_apply.go:88-96`). The real runner resolver scans from the start and returns its first exact-ID match (`internal/agents/resolve.go:47-53`). A legacy or pre-ID participant whose ID is already the adapter family therefore selects the original, current-config discovery and ignores the appended frozen duplicate. See new MAJOR finding R3-M1.

PRIMARY — `TestSnapshotPinsSurviveParticipantResolution` does not invoke `continueAuto` or the runner resolver, so it does not prove the continuation boundary and would still pass if line 1184 were reverted. Its last-write-wins map also hides the production first-match behavior.

### A6 — machine scope: NOT FIXED

PRIMARY — The primary roster rows now use scoped spec, mapping, and roster-source loaders, so ordinary machine-scope rows no longer inherit deck roster fields.

PRIMARY — Two production surfaces remain unscoped. `--scope machine --all` calls `unrosteredRows`, which uses unscoped `LoadAgentSpecs` and `LoadRosterAdapters` (`internal/app/roster_view.go:98-105`); `rosterShow` invokes it without scope at `internal/app/roster.go:436-438`. With an environment config setting the `agy` model to `Gemini 3.5 Flash (High)`, machine-scope `--all --json` reported that environment value; without the environment layer it reported `Gemini 3.6 Pro (High)`.

PRIMARY — Machine `--explain` provenance is also inaccurate for effective values originating from central `[agents.<family>]`. `RosterFieldSourcesScoped` inspects only `[roster.<id>]` entries (`internal/config/runtime.go:1083-1106`). In the current machine config, `claude`'s effective model is set under `[agents.claude]`, but `roster show --scope machine --explain claude-1` printed `built-in default` as its source.

### A7 — frozen eleven-column contract: behavior FIXED; regression guard INCOMPLETE

PRIMARY — `rosterRow.Display` and `rosterRow.Note` now use `json:"-"` (`internal/app/roster.go:213-214`). A production JSON invocation exposed exactly the eleven declared columns and exactly their eleven row keys; neither `display_name` nor `note` was present.

PRIMARY — The agreed golden is still absent. `TestJSONStatusMatchesTextForAHealthyRow` unmarshals only `agent` and `status` and checks only the column count (`internal/app/roster_cycle2_test.go:103-121`), so restoring either forbidden JSON field leaves the test green. See R3-m1.

### A9/A14 and the stale `roster sync` claim: NOT FIXED

PRIMARY — The embedded CLI protocol and live CLI protocol no longer contain either stale `roster sync moves it over/across` phrase.

PRIMARY — The skill repository's working files contain the corrections, but those are uncommitted modifications. Its HEAD is still tagged `v2.5.0`, not the claimed 2.5.1, and the committed `skills/parley-deck/SKILL.md` still says “`parley roster sync` moves it across”; committed `skills/parley-deck/references/COOPERATION.md` still says it “moves it over.” A clean skill checkout therefore still ships the false migration instruction.

PRIMARY — `TestNoSection2AsAStoreInstructions` reads the two skill surfaces only from an optional sibling checkout and logs-and-continues on read failure (`internal/protocol/drift_test.go:275-288`). In this workspace it consumed the dirty corrected files and passed; in an isolated CLI checkout it checks only two of four surfaces and still passes. This does not satisfy the claimed four-surface drift guard.

### A10 — `RosterRevisionOf` omitted `LaunchArgs`: FIXED

PRIMARY — `RosterRevisionOf` now hashes `strings.Join(e.LaunchArgs, "\x1f")` at `internal/runmanifest/manifest.go:94-96`. `TestRosterRevisionCoversLaunchArgs` compares otherwise-identical snapshots and would fail on a reversion; it passed.

### A12 — machine-scope `masked-by-env` false positive: FIXED

PRIMARY — `roster set` now converts the winning source label to its concrete path through `config.RosterSourcePath` before comparing it with the write target. `TestMachineScopeWriteIsNotReportedAsMasked` would fail on the former label-versus-path comparison; it passed.

### A16 residuals: FIXED, with one test gap

PRIMARY — Unmapped guidance now includes `--confirm-breaking` at `internal/app/roster.go:324`; `kimiCodename` accepts only `k` followed by a digit at `internal/agents/modelmeta.go:130-135`; and the 1.40.1 changelog attribution names kimi-1 as well as the other reviewers at `CHANGELOG.md:3-8`.

PRIMARY — The model-meta golden covers a qualified Kimi K3 input, but it has no unrelated `k*` negative case. Reintroducing another broad K-prefix path would not necessarily fail that test.

### IMPLEMENTATION.md accuracy NIT: NOT FIXED

PRIMARY — `IMPLEMENTATION.md` still begins with two blank lines and the heading “Fix-up cycle 2”; it has no document frontmatter or top-level implementation status. Cycle 3 says test-count and full-suite claims were corrected, but the earlier text still says “`go test ./...` green on this machine,” while this required run failed. The nearby caveat names the environment-dependent failure, but it does not make the green claim or missing handoff metadata accurate.

## New findings (by severity, or "none")

### CRITICAL — R3-C1: local, environment, and machine layers can still retire or reactivate committed deck members

PRIMARY — `LoadRosterScoped` correctly limits the ID set, but it first merges every layer into `scope.Entries` (`internal/config/runtime.go:127-138`), and `mergeRosterEntry` applies every layer's `active` value (`internal/config/runtime.go:220-229`). Both `resolveRoster` and `RosterMembership` then decide active/inactive state from those fully layered entries (`internal/app/roster.go:296-300` and `:723-728`).

PRIMARY — Consequently, `[roster.claude-1] active = false` in `agents.local.toml`, the environment config, or the machine file can silently remove a committed deck member from the selected participant set; a higher non-deck layer can likewise reactivate a member the deck retired. This changes run quorum using files collaborators may not see, despite correctly blocking those layers from adding a new ID.

PRIMARY — The ratified field table assigns state to `[roster.<id>].active`, makes absence mean `true`, and assigns `parley-deck/agents.toml` authority for every field (`parley-deck/ideas/roster-operations-standard/consensus.md:353-375`). The remedy is to derive membership state from the same authority record that supplied the ID set, while continuing to layer only the permitted runtime values such as model, effort, and speed. A regression test should set conflicting `active` values in deck, local, environment, and machine layers and assert both `roster show` state and `defaultRosterParticipants` output.

### MAJOR — R3-M1: cycle 3 regresses snapshot pinning for legacy/bare-family participant IDs

PRIMARY — For a participant named `claude`, the helper starts with the live `{ID: "claude", Model: "drifted"}` discovery, resolves and freezes a second `{ID: "claude", Model: <snapshot>}` discovery, then appends it. The runner's first-exact-match rule returns the live first record. The snapshot is present in memory but is not the record launched.

PRIMARY — Cycle 2's `applyRosterSnapshot` adapter fallback updates the single bare-family discovery in place (`internal/app/roster_snapshot_apply.go:40-69`), so this failure is introduced by cycle 3's append-based helper. Composite IDs such as `claude-1` and `claude-2` work because their appended IDs do not collide with the original adapter ID.

PRIMARY — The continuation boundary should replace/deduplicate any discovery with the participant ID, or construct a participant-keyed slice passed to the runner. The regression test must exercise the runner's `ResolveParticipant` selection on the helper output for both a composite ID and a bare-family ID; an end-to-end `continueAuto` test would be stronger.

### MAJOR — R3-M2: `--scope machine` still leaks layered values through `--all` and misattributes central agent values

PRIMARY — The behavioral and source evidence is recorded under A6. Machine scope must be threaded into `unrosteredRows`, and provenance must account for central `[agents.<family>]` fields that produce the effective spec, not just `[roster.<id>]` fields.

### MAJOR — R3-M3: the four-surface protocol fix exists only in a dirty sibling checkout

PRIMARY — The behavioral and repository evidence is recorded under A9/A14. The skill changes need a committed versioned release, and the CLI guard must not claim four-surface enforcement when two surfaces are optional. A self-contained fixture, submodule/package input, or separate mandatory skill-repo test is needed.

### MINOR — R3-m1: the frozen JSON row shape has no reversion-sensitive test

PRIMARY — The behavior is correct at 7220715, but the test suite accepts extra row keys and therefore does not enforce the frozen eleven-column JSON contract. A golden or exact-key assertion over both text and JSON should fail if `display_name`, `note`, or any twelfth field is serialized.

### NIT — R3-N1: the cycle-3 handoff correction is itself incomplete

PRIMARY — The malformed/missing top-level implementation metadata and stale green-suite sentence are recorded under the IMPLEMENTATION.md item above.

## Test-quality assessment

PRIMARY — I ran `go test ./internal/app -run 'TestDeckMembershipIsTheDeckFileNotTheLayeredUnion|TestRosterlessDeckMarksInheritedRows|TestLegacySection2BeatsTheMachineRoster|TestSnapshotPinsSurviveParticipantResolution|TestRosterRevisionCoversLaunchArgs|TestMachineScopeWriteIsNotReportedAsMasked|TestRosterInitRequiresConfirmBreaking' -count=1 -v`; all seven named checks passed and the package reported `ok parley-deck-cli/internal/app 4.714s`:

```text
TestDeckMembershipIsTheDeckFileNotTheLayeredUnion
TestRosterlessDeckMarksInheritedRows
TestLegacySection2BeatsTheMachineRoster
TestSnapshotPinsSurviveParticipantResolution
TestRosterRevisionCoversLaunchArgs
TestMachineScopeWriteIsNotReportedAsMasked
TestRosterInitRequiresConfirmBreaking
```

PRIMARY — I also ran `go test ./internal/protocol -run TestNoSection2AsAStoreInstructions -count=1 -v` and `go test ./internal/agents -run TestDeriveModelMetaGolden -count=1 -v`; each test passed and its package reported `ok` in the present workspace.

PRIMARY — Reversion sensitivity by new cycle-3 test:

- PRIMARY — `TestLegacySection2BeatsTheMachineRoster` is effective: reverting legacy-before-machine authority makes it fail.

- PRIMARY — `TestRosterRevisionCoversLaunchArgs` is effective: removing `LaunchArgs` from the hash makes it fail.

- PRIMARY — `TestMachineScopeWriteIsNotReportedAsMasked` is effective for the A12 path/label bug.

- PRIMARY — `TestRosterInitRequiresConfirmBreaking` is effective for the A3 apply gate.

- PRIMARY — `TestSnapshotPinsSurviveParticipantResolution` is not a continuation-boundary test. Reverting the production call at `continueAuto` does not affect it, and its last-write-wins map masks the runner's first-match semantics.

- PRIMARY — The extended drift test is effective for the two in-repo surfaces, but only conditionally effective for the two skill surfaces. Removing the sibling checkout makes those checks disappear without failure; dirty sibling files can mask stale committed content.

PRIMARY — No new test is a literal constant assertion. The closest functional tautology is `TestSnapshotPinsSurviveParticipantResolution`: it calls the newly introduced helper directly and then proves that the helper appended the two records by folding them into a last-write-wins map. It does not observe which record the continuation's actual consumer selects.

PRIMARY — There is no cycle-3 regression test for A6's scoped `--all` or effective provenance, no exact-key/golden test for A7, no conflicting-layer `active` test for A1, and no unrelated `k*` negative test for A16. Those omissions correspond directly to the remaining failures above.
