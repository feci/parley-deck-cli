---
idea: roster-operations-standard
phase: 8 — re-review
agent: codex-1
round: 2
date: 2026-08-06
reviewed-commit: 57fe9d7
verdict: FINDINGS
---
# Re-review — codex-1

## Verdict

- `PRIMARY` — **FINDINGS.** The implementation claims in `IMPLEMENTATION.md:5-76` are not all true at `57fe9d7`: A1 still overrides a valid legacy §2 with machine membership; A3 is bypassable through `roster init --yes`; A5 collapses same-adapter roster IDs at the real continuation call site; A6's machine scope reads deck values; A7 still has a JSON-only field; A9's guard covers only two of four required surfaces; and A10's revision omits frozen launch args. A12, A14, and A16 are also incomplete.
- `PRIMARY` — I own the findings and recommendations in this file and therefore issue no §15.1 verdict on those claims. The `CONFIRMED`, `PARTIAL`, and `WRONG` labels below adjudicate only the fix-up-cycle-2 implementation claims in `IMPLEMENTATION.md`, which I do not own.
- `PRIMARY` — The reviewed boundary is exactly `a39c2a3e33f1b05961e78d960945807621832116..57fe9d7`; `git rev-parse 57fe9d7^` returned `a39c2a3e33f1b05961e78d960945807621832116`.

## Per-fix verification (A1-A16, DF-1 guard)

### A1 — `PRIMARY` — WRONG

- `PRIMARY` — The narrow N-member case is fixed only when the committed deck layer has at least one `[roster.*]` block: `LoadRosterScoped` returns `deckMembers` before `machineMembers` (`internal/config/runtime.go:155-161`), and both display and participant selection call that resolver (`internal/app/roster.go:263-293`, `internal/app/roster.go:698-711`). The cycle-2 test confirms two committed deck blocks resolve to two (`internal/app/roster_membership_test.go:36-70`).
- `PRIMARY` — The implementation does **not** implement the adopted legacy rule. `LoadRosterScoped` knows nothing about §2 and selects machine membership whenever it finds it (`internal/config/runtime.go:119-162`); `resolveRoster` consults §2 only when `len(scope.Members) == 0` (`internal/app/roster.go:263-283`); `RosterMembership` makes the same choice (`internal/app/roster.go:698-711`). Thus show and participant selection now agree, but agree on the wrong authority.
- `PRIMARY` — I built `57fe9d7`, removed the deck TOML roster from a disposable copy, retained its valid four-row §2, and added `machine-only` only to the machine roster. The valid legacy table was overridden:

  ```text
  $ PARLEY_HOME=<tmp>/home-legacy parley roster show --dir <tmp>/legacy --json \
      | jq '{agents:[.roster[].agent],statuses:(.roster|map({key:.agent,value:.status})|from_entries)}'
  {
    "agents": ["claude-1","codex-1","hermes-1","kimi-1","machine-only"],
    "statuses": {
      "claude-1": ["inherited-roster"],
      "machine-only": ["inherited-roster","effort-unknown"]
    }
  }
  ```

- `PRIMARY` — The same unchanged deck with an empty machine home returned exactly the four §2 IDs, each with `legacy-roster`; this proves the fixture's §2 was parser-valid and that adding machine configuration alone caused the authority switch:

  ```text
  $ PARLEY_HOME=<tmp>/home-empty parley roster show --dir <tmp>/legacy --json | jq ...
  "agents": ["claude-1","codex-1","hermes-1","kimi-1"]
  "claude-1": ["legacy-roster","unmapped"]
  ```

- `PRIMARY` — The same resolver also treats every non-machine layer as membership: `agents.local.toml` and `$PARLEY_HEADLESS_AGENT_CONFIG` are both `machine:false` (`internal/config/runtime.go:318-336`) and their IDs are inserted into `deckMembers` (`internal/config/runtime.go:147-153`). That contradicts the adopted “committed deck file owns membership; higher layers override values” rule in `review/consensus.md:215-223`.
- `PRIMARY` — The inherited-row marker and render refusal themselves work for a truly rosterless, §2-less deck (`internal/app/roster.go:301-306`, `internal/app/roster_render.go:37-46`), but they do not cure the valid-legacy failure above.

### A2 — `PRIMARY` — CONFIRMED

- `PRIMARY` — A disposable TOML-authoritative deck with an extra §2 `ghost-1` row produced `{"agent":"ghost-1","status":["unmapped","section2-only"]}`. Preview and apply both printed `the following §2 row(s) are NOT in this deck's roster and will be removed:` followed by `ghost-1`; apply then printed `Regenerated §2 ...`, and the post-apply §2 had no `ghost-1` row.
- `PRIMARY` — The code matches that behavior: `section2OnlyRows` adds both statuses without adding membership (`internal/app/roster_view.go:52-92`), while render excludes the row and reports `removed` before both preview and apply (`internal/app/roster_render.go:53-63`, `internal/app/roster_render.go:114-151`).

### A3 — `PRIMARY` — WRONG

- `PRIMARY` — The specific shape bypass is closed for `roster set`. On a missing block, `parley roster set sneaky-9 --scope deck --model k3 --yes` exited 2 and printed:

  ```text
  changing  [roster.sneaky-9] in <tmp>/parley-deck/agents.toml
    + model = "k3"

  roster set: this adds a new roster member — a membership change, not a settings change.
  Re-run with --confirm-breaking as well as --yes.
  ```

- `PRIMARY` — The gate is nevertheless bypassable through the still-exposed deprecated verb. `runRoster` calls `rosterInit` without passing `confirmBreaking` (`internal/app/roster.go:166-167`), and `rosterInit` writes every missing roster block when `--yes` is set (`internal/app/roster.go:536-556`). On a disposable valid-legacy deck with an empty target, `parley roster init --scope deck --yes` wrote four new membership blocks without `--confirm-breaking` and exited 0:

  ```text
  note: `parley roster init` is deprecated; prefer `parley roster set AGENT --adapter FAMILY` ...
  Wrote 4 mapping(s) to parley-deck/agents.toml. The driver can now run this roster.
  ```

- `PRIMARY` — This violates D5's membership-change confirmation regardless of which public verb creates the block. The fix needs one membership gate shared by `set` and `init`, or `init --yes` must be made non-writing/require the same confirmation.

### A4 — `PRIMARY` — CONFIRMED in isolation; A5 — `PRIMARY` — WRONG at the real call site

- `PRIMARY` — A4's data path exists: the snapshot captures resolved `LaunchArgs` (`internal/app/roster.go:659-684`) and the consumer restores them to `HeadlessArgs` (`internal/app/roster_snapshot_apply.go:62-68`). A single-adapter continuation therefore retains its frozen autonomous argv.
- `PRIMARY` — A5 is still broken in actual continuation execution. `continueAuto` applies the snapshot immediately to fresh discoveries (`internal/app/app.go:1169-1185`). Fresh discoveries are adapter-family IDs because `discoverConfigured` merely loads specs and calls `agents.Discover` (`internal/app/app.go:2089-2095`). Roster IDs are assigned only later when the runner calls `ResolveParticipant` (`internal/runner/runner.go:204-205`, `internal/runner/runner.go:358-372`; `internal/agents/resolve.go:55-61`).
- `PRIMARY` — Consequently `frozen[d.Spec.ID]` cannot hit a normal roster ID at `internal/app/roster_snapshot_apply.go:42`; execution falls back to `byAdapter` at lines 43-45. If `claude-1` and `claude-2` share adapter `claude`, both consume the first adapter fallback, collapsing the per-ID model, effort, speed, and auto-args despite the new map being keyed by `Agent`.
- `PRIMARY` — `TestSnapshotPinsPerRosterIDNotPerAdapter` manufactures discoveries with `Spec.ID` already equal to `claude-1` and `claude-2` (`internal/app/roster_cycle2_test.go:17-25`). That is the post-resolution identity, but `applyRosterSnapshot` is called pre-resolution. The test therefore passes a state the production call site never supplies and would not catch this regression.

### A6 — `PRIMARY` — WRONG

- `PRIMARY` — `--all` and `--explain` parse, and `--scope machine` chooses machine membership/entries in `rosterScopeFor` (`internal/app/roster_view.go:28-49`). The implementation is still not load-bearing: `resolveRoster` independently loads fully layered agent specs and roster mappings from `root` (`internal/app/roster.go:239-257`), and explain independently loads fully layered provenance (`internal/app/roster_view.go:165-177`).
- `PRIMARY` — I changed only the deck's Claude agent-model value to `deck-only-model` while leaving the machine source unchanged. Machine scope returned the deck value, and explain attributed it nonsensically:

  ```text
  $ PARLEY_HOME=<tmp>/home parley roster show --dir <tmp>/deck --scope machine --json | jq ...
  {"scope":"machine","claude":{"model":"deck-only-model","status":["metadata-unknown"]}}

  $ PARLEY_HOME=<tmp>/home parley roster show --dir <tmp>/deck --scope machine --explain claude-1
  claude-1 — membership from <tmp>/home/agents.toml
  FIELD          EFFECTIVE                SET BY
  adapter        claude                   parley-deck/agents.toml
  model          deck-only-model          built-in default
  ```

- `PRIMARY` — `--scope machine` therefore still answers with deck-layer values and deck-layer provenance. Scope needs to select one coherent spec, mapping, membership, and provenance loader, not only a membership map.

### A7 — `PRIMARY` — WRONG as the frozen v1 contract; healthy status itself is fixed

- `PRIMARY` — A healthy row now prints `ok` and emits `"status":["ok"]`; the prior null-vs-text contradiction is fixed by normalization at `internal/app/roster.go:430-445`.
- `PRIMARY` — Text and JSON still do not have the same v1 shape. The declared contract has eleven fields (`internal/app/roster.go:180-184`), but `rosterRow` still exports JSON-only `display_name` and `note` (`internal/app/roster.go:193-212`). Actual healthy-row output showed eleven `columns` but twelve row keys, including `display_name`, while the text table had only the eleven advertised columns:

  ```text
  $ parley roster show --dir <tmp> | head -2
  AGENT ... AUTO STATUS
  claude-1 ... yes  ok

  $ parley roster show --dir <tmp> --json | jq '{columns,healthy:(.roster[]|select(.agent=="claude-1")|{keys:(keys),status})}'
  "keys": ["adapter","agent","autonomous","display_name","effort","installed",
           "model","model_company","model_family","speed","state","status"],
  "status": ["ok"]
  ```

- `PRIMARY` — This leaves the explicit A7 choice unresolved: remove the extra fields or formally relocate them to the explain contract (`review/consensus.md:291-298`). The claimed “golden test asserts text and JSON together” at `IMPLEMENTATION.md:34-35` is false; see the test-quality section.

### A8 — `PRIMARY` — CONFIRMED

- `PRIMARY` — `NormalizeLegacyModelArgs` rewrites literal model/effort operands to placeholders (`internal/agents/launchargs.go:18-63`), and config application invokes it after an override (`internal/config/runtime.go:726-742`).
- `PRIMARY` — In a disposable deck whose Hermes override hardcoded `--model glm-5p2`, setting the roster model to `xai/grok-4.5` made `roster show` report `{"model":"xai/grok-4.5","status":["effort-unknown"]}` with no `model-drift`; the declared model beat the old literal.

### A9 — `PRIMARY` — PARTIAL / WRONG as a regression guard

- `PRIMARY` — An exact `rg` over all three current COOPERATION.md copies and the external `SKILL.md` returned no matches for the four forbidden phrases. The prose is corrected in the current working trees.
- `PRIMARY` — The drift assertion does **not** cover all four surfaces. `TestNoSection2AsAStoreInstructions` iterates only `defaultCooperation` and `readLiveDeck(t)` (`internal/protocol/drift_test.go:255-275`): the CLI embedded copy and the CLI live copy. It never reads the external skill repo's bundled COOPERATION.md or `SKILL.md`, so either external phrase can return while the test remains green. That directly fails the guard required at `review/consensus.md:316-319`.
- `PRIMARY` — The external corrections are also only uncommitted working-tree modifications: `git -C <skill-repo> status --short` returned `M skills/parley-deck/SKILL.md` and `M skills/parley-deck/references/COOPERATION.md` at HEAD `b806adabdae349b80e83f3126312a66fceae1848`. They are not part of reviewed CLI commit `57fe9d7` or a landed skill commit.

### A10 — `PRIMARY` — WRONG as a complete stale-snapshot detector

- `PRIMARY` — The emitter exists: `rosterSnapshotState` compares the current revision (`internal/app/roster_view.go:225-243`), and `sessions inspect` prints it (`internal/app/app.go:975-989`).
- `PRIMARY` — The compared revision is incomplete. `RosterSnapshotEntry` now says `LaunchArgs` are part of the frozen launch identity (`internal/runmanifest/manifest.go:62-77`), but `RosterRevisionOf` hashes only agent, adapter, model, effort, speed, auto, and installed (`internal/runmanifest/manifest.go:80-93`). A change from an auto-approve argv to another argv with the same `Auto` boolean leaves the revision equal, so inspect reports `current` although the frozen auto-args differ from the deck. Include launch args deterministically in the revision and test the emitter.

### A11 — `PRIMARY` — CONFIRMED

- `PRIMARY` — The protocol changelog now has `Idea:` as a path, `Drafted by:`, and `Summary:` in the required order (`parley-deck/meta/protocol-changelog.md:119-124`) while retaining the substantive sections below.

### A12 — `PRIMARY` — WRONG as an accurate masked-write emitter

- `PRIMARY` — A post-write warning exists (`internal/app/roster_set.go:94-108`), but its source/target comparison mixes an absolute target with the display label `~/.parley/agents.toml`. Writing the machine file itself produced a false “higher layer” warning:

  ```text
  $ PARLEY_HOME=<tmp>/home parley roster set machine-only --scope machine --model gpt-test --yes
  Wrote <tmp>/home/agents.toml
  warning: model = "gpt-test" is MASKED — ~/.parley/agents.toml sets it at a higher layer,
  so the effective value did not change.
  ```

- `PRIMARY` — `RosterFieldSources` deliberately returns the label `~/.parley/agents.toml` (`internal/config/runtime.go:318-336`, `internal/config/runtime.go:939-973`), while `rosterFieldMaskedBy` compares that label to the absolute target with suffix/substring tests (`internal/app/roster_set.go:101-108`). Source identity needs a canonical layer/path, not a presentation string.

### A13 — `PRIMARY` — CONFIRMED

- `PRIMARY` — Top-level help lists show/set/sync/render/migrate and labels `agents list` as inventory, not roster (`internal/app/app.go:119-127`); the runtime matrix carries the same label (`internal/agents/discover.go:494`). Both named docs now explain authority, layering, flags, and verbs (`docs/cli-reference.md:53-121`, `docs/agent-runtime-configuration.md:38-82`).

### A14 — `PRIMARY` — WRONG as “skill corrected (2.5.1)” landing

- `PRIMARY` — The external `SKILL.md` working-tree content now correctly says `roster sync` does not migrate legacy decks and names migrate/set/render (`<skill-repo>/skills/parley-deck/SKILL.md:314-319`).
- `PRIMARY` — It has not landed as skill 2.5.1. The skill-repo changes are uncommitted, `<skill-repo>/package.json:3` remains `"version": "2.5.0"`, and the status command reported installer 2.5.0 plus every installed runtime skill at 2.5.0:

  ```text
  $ parley-deck-skill status --target all --project . --json | jq '{installer:.installer.version,cli:.parleyCli.version,runtimes:(.runtimeInstalls|map({target,version}))}'
  "installer": "2.5.0"
  "cli": "parley 1.40.1"
  "runtimes": [{"target":"codex","version":"2.5.0"}, ..., {"target":"opencode","version":"2.5.0"}]
  ```

### A15 — `PRIMARY` — CONFIRMED for the agreed scope

- `PRIMARY` — Every unmatched `--keep` token is rejected non-zero before a write (`internal/app/roster_sync.go:52-105`), and the cycle-2 test proves the typoed pin survives (`internal/app/roster_cycle2_test.go:53-72`).
- `PRIMARY` — Apply re-reads the current entries and compares each field-to-drop with its previewed old value before editing (`internal/app/roster_sync.go:140-165`). That satisfies the adopted field-old-value binding. A full-file CAS remains DF-1/DF-3 territory and is not charged here.

### A16 — `PRIMARY` — PARTIAL

- `PRIMARY` — File-mode preservation landed (`internal/app/roster_set.go:240-251`), unreadable manifests warn (`internal/app/app.go:1176-1184`), `agents list` is relabelled (`internal/agents/discover.go:494`), and reactivation is detected before retirement (`internal/app/roster_set.go:261-275`).
- `PRIMARY` — The unmapped remediation now names `roster set`, but the emitted command omits `--confirm-breaking` (`internal/app/roster.go:311-316`) and is therefore refused by the A3 gate for a missing block. The replacement guidance is still not executable as printed.
- `PRIMARY` — Moving broad prefix `"k"` after `"kimi"` fixes the unreachable Kimi rule, but `strings.HasPrefix` still classifies every unrelated lower-case `k*` model as Moonshot (`internal/agents/modelmeta.go:54-69`, `internal/agents/modelmeta.go:103-108`). The agreed “unrelated k* misclassification” half remains.
- `PRIMARY` — The agreed attribution correction did not land: `CHANGELOG.md:7` still says the defects were “found by codex-1 and hermes-1” and omits kimi-1, contrary to `review/consensus.md:387-388`.

### DF-1 interim guard — `PRIMARY` — CONFIRMED

- `PRIMARY` — `roster migrate --yes` without `--confirm-breaking` exited 2 before traversal:

  ```text
  roster migrate: --yes rewrites the roster of EVERY deck under this root.
  The full migration contract (compare-and-swap, per-deck confirmation, version gate) is not
  implemented yet, so this operation is attended-only. Re-run with --confirm-breaking as well as --yes.
  ```

- `PRIMARY` — The guard is explicit at `internal/app/roster_migrate.go:59-70`. An applying traversal checks `git status --porcelain` and records a dirty repository as skipped before `migrateOneDeck` (`internal/app/roster_migrate.go:82-95`, `internal/app/roster_migrate.go:378-388`). I did not mutate either reviewed repository to manufacture a dirty-tree execution case.

## Findings (by severity, or "none")

### CRITICAL

1. `PRIMARY` — **A1: valid legacy §2 compatibility membership is overridden by the machine roster.** `resolveRoster` and `RosterMembership` treat any machine roster as config authority before consulting §2 (`internal/app/roster.go:263-283`, `internal/app/roster.go:698-711`). A four-member legacy deck demonstrably displayed and would select a fifth machine-only member. Higher local/env roster blocks can also become membership because every non-machine layer is pooled as `deckMembers` (`internal/config/runtime.go:124-160`). Fix authority selection before value layering: committed deck blocks first, else valid legacy §2, else explicitly inherited machine display; participant selection must consume that same result.

### MAJOR

1. `PRIMARY` — **A3: `roster init --yes` bypasses the membership confirmation.** It writes missing roster blocks without receiving or checking `--confirm-breaking` (`internal/app/roster.go:166-167`, `internal/app/roster.go:536-556`).
2. `PRIMARY` — **A5: per-roster-ID snapshot pins still collapse at continuation.** Snapshot application runs against adapter-ID discoveries before participant resolution, so every real lookup misses `Agent` and uses the adapter fallback (`internal/app/app.go:1169-1185`, `internal/app/roster_snapshot_apply.go:40-45`, `internal/runner/runner.go:358-372`). Resolve/clone each participant to its roster ID before applying its snapshot, and test through the continuation/selection boundary with two IDs sharing one adapter.
3. `PRIMARY` — **A6: machine scope is polluted by deck values and provenance.** Only membership/entries are scoped; specs, mappings, and field sources remain layered (`internal/app/roster.go:239-257`, `internal/app/roster_view.go:165-177`). An actual deck-only model appeared in `--scope machine`.
4. `PRIMARY` — **A7: the frozen v1 text/JSON contract still diverges.** `display_name` and `note` remain JSON-only fields outside the eleven declared columns (`internal/app/roster.go:180-212`), and the new test is not a text+JSON golden.
5. `PRIMARY` — **A9: the required four-surface drift guard covers only two surfaces.** The external bundled protocol and `SKILL.md` can regress without failing `TestNoSection2AsAStoreInstructions` (`internal/protocol/drift_test.go:255-275`).
6. `PRIMARY` — **A10: stale-snapshot comparison omits frozen launch args.** `LaunchArgs` are part of the snapshot but absent from `RosterRevisionOf` (`internal/runmanifest/manifest.go:62-93`), so some autonomy drift reports `current`.

### MINOR

1. `PRIMARY` — **A12: machine-scope writes falsely report themselves as masked.** The emitter compares an absolute target path with a display source label (`internal/app/roster_set.go:94-108`).
2. `PRIMARY` — **A14: the skill correction is neither committed nor versioned as 2.5.1.** External content is better, but the skill repo remains at HEAD `b806adab...`, dirty in the two relevant files, with package/runtime version 2.5.0.
3. `PRIMARY` — **A16: three agreed cleanup items remain incomplete.** Printed unmapped guidance omits the newly mandatory confirmation, the broad `k*` classification remains, and kimi-1 attribution remains missing (`internal/app/roster.go:311-316`, `internal/agents/modelmeta.go:54-69`, `CHANGELOG.md:7`).

### NIT

1. `PRIMARY` — **The new implementation handoff is malformed and overclaims verification.** `IMPLEMENTATION.md:1-3` starts with blank lines and a section heading rather than the required top-level Phase-5 status/head-commit metadata, while `IMPLEMENTATION.md:75-76` claims “Full suite green” despite the reproducible runner failure below. The runner failure predates this cycle, so this is handoff accuracy/audit debt, not a new roster regression.

## Test-quality assessment

- `PRIMARY` — The required command was run exactly. Build succeeded, but the combined command exited 1 in a package untouched by `a39c2a3..57fe9d7`:

  ```text
  $ go build ./... && go test ./...
  ...
  ok  parley-deck-cli/internal/runmanifest (cached)
  --- FAIL: TestDurableKillEndToEndRealProcess (0.02s)
      durablekill_test.go:116: a live attributed process should be killed, got
      {AgentID:sleeper Killed:false ... Message:process verification failed
      (no recorded boot id); not killed}
  FAIL parley-deck-cli/internal/runner
  FAIL
  ```

- `SECONDARY` — This exact runner failure was already recorded in my round-01 review (`review/round-01/codex-1.md:42-47`), and `git diff --name-only a39c2a3 57fe9d7 -- internal/runner` is empty. I therefore do not classify it as newly broken by cycle 2.
- `PRIMARY` — The changed/high-relevance packages and vet pass:

  ```text
  $ go test ./internal/app ./internal/agents ./internal/config ./internal/protocol ./internal/runmanifest -count=1 && go vet ./...
  ok  parley-deck-cli/internal/app          27.592s
  ok  parley-deck-cli/internal/agents        0.138s
  ok  parley-deck-cli/internal/config        0.418s
  ok  parley-deck-cli/internal/protocol      0.702s
  ok  parley-deck-cli/internal/runmanifest   0.273s
  ```

- `PRIMARY` — The new A1 test covers committed-deck-versus-machine union but omits a valid legacy §2 (`internal/app/roster_membership_test.go:36-70`). Its rosterless fixture has no COOPERATION table (`internal/app/roster_membership_test.go:72-86`), so it cannot catch the critical legacy override.
- `PRIMARY` — The A2 test calls `section2OnlyRows` directly with a hand-supplied membership map (`internal/app/roster_membership_test.go:99-114`). It never runs `roster show`, render preview, render apply, or post-apply exclusion. The shipped behavior happened to pass my end-to-end checks, but this test would not catch a missing render report.
- `PRIMARY` — The A3 test directly calls `membershipChange(..., false)` (`internal/app/roster_membership_test.go:116-131`). The field-shaped strings are irrelevant once `existed` is false, so the test nearly hard-wires the helper's deciding predicate and does not exercise block parsing, CLI flag plumbing, or the `roster init` bypass.
- `PRIMARY` — The A4/A5 per-ID test uses already-resolved roster IDs even though production snapshot application receives unresolved adapter IDs (`internal/app/roster_cycle2_test.go:17-34`). This is the same wrong-path test pattern called out in the re-review prompt: it passes the defect.
- `PRIMARY` — The A7 test never renders text, never asserts exact `status == ["ok"]`, and unmarshals a subset that silently discards extra JSON fields (`internal/app/roster_cycle2_test.go:95-121`). It is not the claimed joint golden and cannot detect `display_name`/`note` contract drift.
- `PRIMARY` — Searches found no cycle-2 tests for load-bearing `--scope machine`, `--all`, `--explain`, `rosterSnapshotState`, or DF-1's migration gates/dirty-tree skip. A6 and A10 were therefore able to remain wrong while the targeted suites passed.
- `PRIMARY` — The A9 regression test scans only two in-repo strings (`internal/protocol/drift_test.go:255-275`), not the four ratified surfaces. The A15 test covers unmatched keeps (`internal/app/roster_cycle2_test.go:53-72`) but no test injects a changed old value between preview computation and apply.
- `PRIMARY` — None of the new tests is literally a constant assertion, but the A3 test hard-wires the deciding boolean, the A5 test supplies a production-impossible identity state, and the A7 test's “text and JSON” premise is not exercised. Those are assumption-encoding tests that would pass the regressions documented above.
