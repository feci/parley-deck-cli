---
idea: roster-operations-standard
phase: 8 — re-review
agent: hermes-1
round: 2
date: 2026-08-06
reviewed-commit: 57fe9d7
verdict: FINDINGS
---
# Re-review — hermes-1

## Verdict

FINDINGS — one CRITICAL, one MAJOR, two MINOR. Fifteen of the sixteen agreed fixes
landed correctly and are verified by PRIMARY reproduction; the exception is A1's
legacy-fallback clarification, which did not land. The test suite is green and the
build is clean, but the CRITICAL means cycle 2 is not the exit condition.

Evidence provenance: PRIMARY = I read the source or ran the tool this session against
scratch decks under /tmp with PARLEY_HOME isolation. I built the binary at 57fe9d7
(`go build -o /tmp/parley-rev2 ./cmd/parley/`) and ran `roster show`, `roster render`,
`roster set`, `roster sync`, and `roster migrate` against isolated fixtures. Nothing
was written to the real deck or to ~/.parley/agents.toml (one accidental write to the
real deck's agents.toml during A3 testing was immediately reverted with
`git checkout -- parley-deck/agents.toml`; git status confirmed clean).

Build and suite: `go build ./...` exit 0; `go test ./...` all packages ok (24 packages,
0 failures).

## Per-fix verification (A1-A16, DF-1 guard)

### A1 — [CRITICAL] deck membership is the deck file — PARTIALLY LANDED

Three of A1's four sub-requirements landed and are verified PRIMARY. The fourth — the
legacy-fallback clarification codex-1's condition 2 added — did not.

**Sub 1: a deck declaring N members resolves to exactly N.** PRIMARY — verified:

```
$ cat /tmp/rev2-test/parley-deck/agents.toml   # declares EXACTLY two
[roster.claude-1]
adapter = "claude"
[roster.kimi-1]
adapter = "kimi"

$ PARLEY_HOME=/tmp/rev2-machine /tmp/parley-rev2 roster show --dir /tmp/rev2-test
AGENT        ADAPTER    STATE    INSTALLED ...
claude-1     claude     active   yes       ...
kimi-1       kimi       active   no        ...
```

Two rows, not five. `LoadRosterScoped` (runtime.go:118-163) separates `deckMembers`
from `machineMembers` and returns deckMembers when non-empty. `resolveRoster`
(roster.go:287) iterates `scope.Members` (the deck file's set), not `entries` (the
layered union). Correct.

**Sub 2: participant selection uses the same authority.** PRIMARY — code-read:
`defaultRosterParticipants` (app.go:2454) calls `RosterMembership(root)`, which
(roster.go:698-712) calls `config.LoadRosterScoped(root)` — the same function
`roster show` uses. Preset validation (app.go:1826) also calls `RosterMembership`.
The old `protocol.ReadRosterIDs` fallback is only reached when
`len(scope.Members) == 0`. The two-sources-of-truth defect from cycle 1 is closed
for decks with `[roster.*]` blocks.

**Sub 3: a rosterless deck marks inherited rows; render refuses without
--adopt-inherited.** PRIMARY — verified:

```
$ PARLEY_HOME=/tmp/rev2-machine /tmp/parley-rev2 roster show --dir /tmp/rev2-empty
... five rows, each with STATUS containing inherited-roster ...

$ PARLEY_HOME=/tmp/rev2-machine /tmp/parley-rev2 roster render --dir /tmp/rev2-render --yes
roster render: this deck declares no roster of its own; the 5 rows shown come from ~/.parley/agents.toml.
... re-run with --adopt-inherited ...
EXIT=1

$ PARLEY_HOME=/tmp/rev2-machine /tmp/parley-rev2 roster render --dir /tmp/rev2-render --yes --adopt-inherited
Regenerated §2 ...
EXIT=0
```

`renderRosterTable` (roster_render.go:38-49) refuses when `scope.Inherited &&
!adoptInherited`. Correct.

**Sub 4: a valid legacy §2 acts as compatibility membership, NOT overridden by the
machine roster.** PRIMARY — **DID NOT LAND.** This is the CRITICAL finding (see
below).

### A2 — §2-only IDs reported, never erased — LANDED

PRIMARY — verified with a deck declaring `[roster.claude-1]` plus a §2 table
containing `ghost-1`:

```
$ PARLEY_HOME=/tmp/rev2-a2-machine /tmp/parley-rev2 roster show --dir /tmp/rev2-a2
...
ghost-1      —          active   no   ...  unmapped,section2-only
  ⚠ declared only in the §2 table, which is no longer authoritative ...

$ PARLEY_HOME=/tmp/rev2-a2-machine /tmp/parley-rev2 roster render --dir /tmp/rev2-a2 --dry-run
the following §2 row(s) are NOT in this deck's roster and will be removed:
  - ghost-1 (reported `unmapped` by `parley roster show`; ...)
would regenerate §2 ...
| `claude-1` | . | participant | active |
```

`section2OnlyRows` (roster_view.go:44-79) surfaces the §2-only ID with both
`unmapped` and `section2-only` status. The generated table contains only claude-1
(ghost-1 is not auto-added). The removal is reported in the preview output. On
apply, `reportRemoved` (roster_render.go:116-124) is called before the write.
Correct.

### A3 — membership gate keys on block existence — LANDED

PRIMARY — verified the exact bypass from the consensus:

```
$ PARLEY_HOME=... /tmp/parley-rev2 roster set sneaky-9 --scope deck --model k3 --yes
...
roster set: this adds a new roster member — a membership change, not a settings change.
Re-run with --confirm-breaking as well as --yes.
EXIT=2

$ PARLEY_HOME=... /tmp/parley-rev2 roster set sneaky-9 --scope deck --model k3 --yes --confirm-breaking
...
(adds a new roster member — confirmed with --confirm-breaking)
Wrote ...
EXIT=0
```

`membershipChange` (roster_set.go:245-261) now takes `existed bool` from
`blockExists` (roster_set.go:264-273), which checks for the `[roster.<agent>]`
header in the same bytes the patcher edits. A new block written with only `--model`
is caught. Correct.

### A4/A5 — snapshot keyed per roster ID, auto-args pinned — LANDED

PRIMARY — code-read + test-read. `applyRosterSnapshot`
(roster_snapshot_apply.go:26-44) keys `frozen` by `e.Agent` (roster ID), with
`byAdapter` as a fallback for pre-ID runs (first-writer-wins so the fallback cannot
reorder). `LaunchArgs` (manifest.go:73-78) is pinned: when `len(e.LaunchArgs) > 0`,
`d.Spec.HeadlessArgs` is replaced with the frozen argv
(roster_snapshot_apply.go:62-66). The test `TestSnapshotPinsPerRosterIDNotPerAdapter`
(roster_cycle2_test.go:17-34) uses two IDs sharing the `claude` adapter with different
models and verifies per-ID pinning. `TestSnapshotPinsAutonomousLaunchArgs`
(roster_cycle2_test.go:38-51) verifies a dropped `--yolo` is restored. Correct.

### A6 — --all, --explain, --scope — LANDED

PRIMARY — verified all three:

- `--all` lists configured adapters no roster declares, each with `not-in-roster`
  status and a note naming `roster set` to add one. Verified with a deck declaring
  only claude-1: the output included 16 unrostered adapters (agy, codex, hermes,
  etc.).
- `--explain claude-1` prints per-field provenance (FIELD/EFFECTIVE/SET BY table),
  naming which layer set each value. Verified.
- `--scope machine` produces a DIFFERENT result from the default deck scope (5 rows
  vs 2). Verified — scope is now load-bearing, not a no-op.
- `roster init` accepts the `deck` spelling and prints a deprecation notice.
  `session` is a hidden alias via `rosterScopeAlias` (roster.go:72-76). Verified by
  code-read.

### A7 — JSON/text contract agreement — LANDED

PRIMARY — verified:

```
$ /tmp/parley-rev2 roster show --dir /tmp/rev2-a6 --json
{
  "columns": [...],
  "roster": [
    {
      "agent": "claude-1",
      ...
      "status": ["ok"],
      ...
    }
  ],
  "schema_version": 1,
  "scope": "deck"
}
```

`status` is `["ok"]` not `null`. `scope` is present. The golden test
`TestJSONStatusMatchesTextForAHealthyRow` (roster_cycle2_test.go:97-122) unmarshals
the JSON and asserts `r.Status != nil` for every row, plus schema_version and column
count. Correct.

### A8 — legacy normalizer — LANDED

PRIMARY — verified end-to-end. A machine config with a `headless_args` override
hardcoding `--model ancient-hardcoded-literal`, plus a deck declaring
`model = "claude-opus-5"`:

```
$ /tmp/parley-rev2 roster show --dir /tmp/rev2-a8 --explain claude-1
...
model          claude-opus-5            parley-deck/agents.toml
```

The declared `model` field wins. `NormalizeLegacyModelArgs`
(launchargs.go:135-166) rewrites the literal to `{model}`; `applyOverride`
(runtime.go:731-741) calls it and records `Sources["headless_args_normalized"]`.
The test `TestLegacyModelArgsAreNormalizedToPlaceholders`
(roster_cycle2_test.go:76-93) verifies the rewrite and the boolean-flag edge case.
Correct.

### A9 — §2-as-a-store phrases removed, drift assertion — PARTIALLY LANDED

PRIMARY — the four banned phrases are gone from the live deck
(`parley-deck/COOPERATION.md`), the embedded default
(`internal/protocol/defaults/COOPERATION.md`), and the skill repo's SKILL.md. The
drift test `TestNoSection2AsAStoreInstructions` (drift_test.go:255-275) passes and
would fail if a phrase returned to either checked copy — it uses
`strings.Contains` on the exact banned strings.

The skill repo's bundled COOPERATION.md
(`parley-deck-skill/.../references/COOPERATION.md`) also has the four banned phrases
gone (verified by grep). However, the drift test does NOT check this third copy —
it checks only `defaultCooperation` (embedded) and `readLiveDeck` (live deck). See
Findings MINOR-1.

### A10 — stale-snapshot in sessions inspect — LANDED

PRIMARY — code-read. `inspectSession` (app.go:958-960) calls
`rosterSnapshotState(session.WorkspaceRoot, manifest)`, which
(roster_view.go:232-243) compares `runmanifest.RosterRevisionOf(current)` against
`m.RosterRevision`, returning `stale-snapshot`, `current`, `no-snapshot`, or
`unknown`. `printSessionDetail` (app.go:987-988) prints "Roster snapshot: <state>".
Correct.

### A11 — changelog §7 format — LANDED

PRIMARY — `parley-deck/meta/protocol-changelog.md:119-125` now has:

```
## 2026-08-06 — §2 roster authority moves to `parley-deck/agents.toml`
Idea: ideas/roster-operations-standard/
Drafted by: claude-1
Summary: §2's roster table stops being the hand-edited membership store ...
```

All three §7 template fields present. Correct.

### A12 — masked-by-env emitter — LANDED

PRIMARY — code-read. `rosterSet` (roster_set.go:80-91) re-resolves after writing
via `rosterFieldMaskedBy` (roster_set.go:94-104), which checks whether
`RosterFieldSources` reports a different layer than the one just written. If so, it
emits a `masked-by-env` warning to stderr. Correct.

### A13 — discoverability — LANDED

PRIMARY — verified: `parley --help` lists all five roster verbs (show, set, sync,
render, migrate). `agents list` prints "adapter/runtime inventory — NOT the roster
(see `parley roster show`)" as its first line (discover.go:491).
`docs/cli-reference.md` has 18 roster matches; `docs/agent-runtime-configuration.md`
has 9. Both gained full roster sections. Correct.

### A14 — skill corrected (2.5.1) — LANDED for SKILL.md, but see MINOR-2

PRIMARY — the skill repo's SKILL.md (committed in b806ada, skill 2.5.0/2.5.1) no
longer says `roster sync` moves legacy decks across. SKILL.md:314-319 now reads:
"roster sync does **not** migrate it — sync only rebases an existing deck roster ...
The remediation is `parley roster migrate` ... or `parley roster set` ... then
`parley roster render`". The membership-is-the-deck-file rule, new status codes, and
new flags are documented. Correct for SKILL.md.

Note: the installed skill copy at `~/.hermes/skills/parley-deck/SKILL.md` still has
the old "and the §2 roster)" phrase at line 171 — but this is a stale install, not a
source defect. The skill repo source is correct.

### A15 — sync hardening — LANDED

PRIMARY — code-read + test. `unmatchedKeeps` (roster_sync.go:207-215) returns
--keep tokens that matched no override; the main function (roster_sync.go:98-106)
rejects them with exit 2 and names each unmatched token. Preview/apply binding
(roster_sync.go:145-161) re-reads the deck file before deleting and refuses if a
field value changed since the preview. The test `TestSyncRejectsUnmatchedKeepTokens`
(roster_cycle2_test.go:55-72) verifies a typoed `--keep kimi-1.modle` is rejected,
the error names the token, and the pin survives. Correct.

### A16 — assorted — LANDED

PRIMARY — verified each sub-item by code-read:
- File-mode: `writeRosterFileAtomic` (roster_set.go:243-251) now stats the target
  for its mode and falls back to 0644, via `fsutil.WriteFileAtomic`.
- modelmeta: the `"k"` rule (modelmeta.go:69) is after `"kimi"` (modelmeta.go:65),
  with a comment explaining why.
- Reactivation: `membershipChange` (roster_set.go:252) checks `+ active = true`
  before `+ active = false`, so a revival reports "reactivates" not "retires".
- `roster init` deprecation notice (roster.go:476-478).
- Unmapped guidance now names `roster set` not `roster init` (roster.go:315).
- Continuation warns on unreadable manifest (app.go:1176-1180).
- `agents list` relabelled (discover.go:491, A13 above).

### DF-1 interim guard — LANDED

PRIMARY — verified:

```
$ /tmp/parley-rev2 roster migrate --dir /tmp/rev2-df1 --backup-dir /tmp/rev2-df1-backup --yes
roster migrate: --yes rewrites the roster of EVERY deck under this root.
...
Re-run with --confirm-breaking as well as --yes.
EXIT=2
```

`--yes --confirm-breaking` proceeds (exit 1 because no COOPERATION.md exists, which
is correct). Dirty-tree skip: `deckTreeDirty` (roster_migrate.go:382-389) runs
`git status --porcelain` and skips the deck if output is non-empty
(roster_migrate.go:87-95). Correct.

## Findings (by severity)

### [CRITICAL] A1 legacy-fallback clarification did not land — a valid legacy §2 deck is overridden by the machine roster

PRIMARY — the consensus A1 (rev 4, `consensus.md:220-223`) explicitly added codex-1's
condition 2:

> "a deck with no roster of its own" means **neither** a deck-level `[roster.*]`
> block **nor** a valid legacy §2 table. A deck carrying a valid legacy §2 keeps that
> table as its compatibility membership — reported `legacy-roster` — until it is
> migrated; the machine roster is not inherited over it.

`LoadRosterScoped` (runtime.go:118-163) does not implement this. It reads only TOML
config layers and has no knowledge of the legacy §2 table. When a deck has a §2
table but no `[roster.*]` block, and the machine config has `[roster.*]` blocks:

1. `deckMembers` is empty (no deck TOML roster).
2. `machineMembers` has the machine's members.
3. `len(deckMembers) == 0`, so it sets `out.Members = machineMembers,
   out.Inherited = true`.
4. `RosterMembership` (roster.go:698-710) returns those machine members because
   `len(scope.Members) > 0` — the `protocol.ReadRosterIDs` fallback is never reached.
5. `resolveRoster` (roster.go:268) sets `legacy = len(scope.Members) == 0` → false,
   so it iterates the machine members, not the §2 members.

Reproduced PRIMARY:

```
$ cat /tmp/rev2-legacy/parley-deck/COOPERATION.md  # §2 declares claude-1, kimi-1
## 2. Active agents (roster)
| `claude-1` | . | participant |
| `kimi-1` | . | reviewer |

# No agents.toml in the deck — only a legacy §2 table

$ PARLEY_HOME=/tmp/rev2-machine /tmp/parley-rev2 roster show --dir /tmp/rev2-legacy
AGENT        ADAPTER    ...  STATUS
claude-1     claude     ...  inherited-roster,metadata-unknown
codex-1      codex      ...  inherited-roster,...
hermes-1     hermes     ...  inherited-roster,...
kimi-1       kimi       ...  inherited-roster,...
opencode-1   opencode   ...  inherited-roster,...
```

Five inherited machine rows, not two `legacy-roster` §2 rows. The §2 members
(claude-1, kimi-1) are drowned by three machine-only agents (codex-1, hermes-1,
opencode-1) that the deck never declared.

This is the same class of defect cycle 1 shipped: an authority cutover where
selection still reads the wrong source. Here, `LoadRosterScoped` was built to
separate membership from values, but it only knows about TOML — the legacy §2 table
that `resolveRoster`'s `legacy` branch and `RosterMembership`'s fallback both handle
is invisible to it. When the machine layer has members, the legacy fallback path is
unreachable.

Participant selection is affected too: `defaultRosterParticipants` (app.go:2454)
calls `RosterMembership`, which returns the machine's five members. A run on this
deck would select five participants, not the two the §2 table declares.

The test `TestDeckMembershipIsTheDeckFileNotTheLayeredUnion`
(roster_membership_test.go:36-70) does NOT catch this: it uses a deck WITH
`[roster.*]` blocks, not a legacy §2-only deck. The test
`TestRosterlessDeckMarksInheritedRows` (roster_membership_test.go:74-86) uses an
empty machine config — so it hits the legacy fallback correctly. Neither test
exercises the combination: legacy §2 deck + machine roster.

Suggested fix: `LoadRosterScoped` must check for a valid legacy §2 table (via
`protocol.ReadRosterIDs`) when `deckMembers` is empty, BEFORE falling back to the
machine layer. If §2 has members, those are the deck's membership (with
`legacy-roster` status), and the machine layer seeds values only — the same rebase
model. The machine roster is inherited only when the deck has NEITHER a `[roster.*]`
block NOR a valid §2 table.

### [MAJOR] A14: the "roster sync moves it over" stale guidance survives in all three COOPERATION.md copies

PRIMARY — A14 corrected the SKILL.md, but the identical false guidance persists in
all three COOPERATION.md copies:

```
$ grep -n "moves it over" parley-deck/COOPERATION.md internal/protocol/defaults/COOPERATION.md
parley-deck/COOPERATION.md:125:...until `parley roster sync` moves it over.
internal/protocol/defaults/COOPERATION.md:124:...until `parley roster sync` moves it over.

$ grep -n "moves it over" ../parley-deck-skill/.../references/COOPERATION.md
:124:...until `parley roster sync` moves it over.
```

`roster sync` on a legacy deck reports "already inherits ... nothing to do" — it
moves nothing across. The SKILL.md was corrected to say so explicitly
(SKILL.md:314: "roster sync does **not** migrate it"), but the COOPERATION.md copies
still tell a facilitator the opposite. A14's consensus text names "SKILL.md" as the
surface, and the IMPLEMENTATION.md says "skill corrected (2.5.1)", so this may be
scoped out — but the result is that the protocol document and the skill now
contradict each other on the same point. An agent following COOPERATION.md's
bootstrap section will try `roster sync` on a legacy deck and get "nothing to do",
with no pointer to the actual remediation (`migrate` or `set` + `render`).

I tag this MAJOR rather than MINOR because it is a false instruction in the normative
protocol document — the same class of defect A9 addressed (a surface instructing
the reader to do the wrong thing), just a different phrase that A14's scope did not
reach.

### [MINOR] A9 drift test does not cover the skill repo's bundled COOPERATION.md

PRIMARY — `TestNoSection2AsAStoreInstructions` (drift_test.go:262-264) checks two
copies:

```go
for _, tc := range []struct{ name, text string }{
    {"embedded default", defaultCooperation},
    {"live deck", readLiveDeck(t)},
} {
```

The skill repo's bundled COOPERATION.md
(`parley-deck-skill/skills/parley-deck/references/COOPERATION.md`) is the third copy
A9 named. The test does not load or check it. The banned phrases are currently gone
from all three copies, so the test passes — but a regression in the skill repo's
copy would not be caught. The test's own comment says "Phrases, not prose style, are
what regressed — so phrases are what this test forbids," yet one of the three
surfaces A9 named is outside the test's reach.

This is MINOR because the skill repo is a separate git repository and the drift
test lives in the CLI repo; cross-repo test coupling is not always desirable. But
the gap should be noted: A9's "all four surfaces" claim is enforced on three, not
four.

### [MINOR] A4/A5: the snapshot-auto-args pin replaces HeadlessArgs wholesale, dropping any args the frozen entry did not carry

PRIMARY — code-read. `applyRosterSnapshot` (roster_snapshot_apply.go:62-65):

```go
if len(e.LaunchArgs) > 0 {
    d.Spec.HeadlessArgs = append([]string(nil), e.LaunchArgs...)
}
```

This replaces the discovered `HeadlessArgs` entirely with the frozen copy. If the
frozen entry's `LaunchArgs` was incomplete (e.g. a run frozen before LaunchArgs was
added to the snapshot, or a snapshot that captured a subset), any args the current
discovery legitimately adds would be lost. The model/effort/speed pins above it
(roster_snapshot_apply.go:55-61) are field-by-field, but the auto-args pin is
wholesale. This is the correct behavior for the common case (the frozen argv is the
complete argv), but it is a rougher edge than the other pins. The test
`TestSnapshotPinsAutonomousLaunchArgs` (roster_cycle2_test.go:38-51) verifies the
restore but not the wholesale-replace edge.

I tag this MINOR because the frozen `LaunchArgs` is captured from
`spec.ResolveLaunchArgs()` at run creation (roster.go:386), which is the full
resolved argv — so in practice the frozen copy is complete. The risk is theoretical
for current runs but would matter if the snapshot format ever captured a partial
argv.

## Test-quality assessment

Nine new tests across two files. I assessed each for whether it tests behavior or
encodes an assumption, and whether it would catch a regression.

**roster_membership_test.go (6 tests):**

1. `TestDeckMembershipIsTheDeckFileNotTheLayeredUnion` — REAL. Creates a deck with 2
   TOML members and a machine with 5, asserts membership = 2, asserts machine-only
   agents are absent, AND asserts the machine model value still inherits to a deck
   member. This is the test that proves A1 sub 1. It would catch a regression to the
   layered-union behavior. It does NOT test the legacy §2 + machine case (the
   CRITICAL above).

2. `TestRosterlessDeckMarksInheritedRows` — REAL but narrow. Uses an empty machine
   config, so it tests "empty machine → inherited" rather than "machine with members
   → inherited." The inherited marking would still catch a regression, but the
   fixture does not exercise the full machine-roster scenario I tested at runtime.

3. `TestRenderRefusesToCommitAnInheritedRoster` — REAL. Asserts render errors
   without --adopt-inherited and succeeds with it. Directly tests the refusal gate.

4. `TestSection2OnlyIDIsReportedNotDropped` — REAL. Creates a §2 table with a
   §2-only ID, calls `section2OnlyRows` directly, asserts the row exists with
   `unmapped` and `section2-only` status. Would catch a regression to silent
   dropping.

5. `TestMembershipGateCatchesNewBlockWrittenWithAnyField` — REAL. Tests
   `membershipChange` with `existed=false` for model/effort/speed, asserts each is a
   membership change. Tests `existed=true` for model, asserts it is NOT. Tests
   reactivation vs retirement ordering. Directly tests the A3 fix and the A16
   reactivation fix. Not tautological — it tests the function's logic, not a
   wrapper.

6. `hasStatus` — helper, not a test.

**roster_cycle2_test.go (5 tests):**

1. `TestSnapshotPinsPerRosterIDNotPerAdapter` — REAL. Two IDs sharing an adapter
   with different frozen models, both discovered as "drifted." Asserts each gets its
   own frozen model. Would catch a regression to adapter-keyed mapping.

2. `TestSnapshotPinsAutonomousLaunchArgs` — REAL. Frozen entry has `--yolo`,
   discovered does not. Asserts `--yolo` is restored. Would catch a regression to
   not pinning auto-args.

3. `TestSyncRejectsUnmatchedKeepTokens` — REAL. Typoed `--keep kimi-1.modle`,
   asserts exit != 0, error names the token, and the pin survives in the file.
   Would catch a regression to silent acceptance.

4. `TestLegacyModelArgsAreNormalizedToPlaceholders` — REAL. Hardcoded literal
   rewritten to placeholder, boolean-flag edge case left alone. Would catch a
   regression to not normalizing.

5. `TestJSONStatusMatchesTextForAHealthyRow` — REAL but has a subtle weakness. It
   asserts `r.Status != nil` for every row, which catches the `null` regression.
   But it does not assert the status is `["ok"]` specifically — a row with
   `["model-drift"]` would also pass (`!= nil`). The test proves "JSON status is not
   null" but not "JSON status agrees with text." For a healthy row the text renders
   `ok`, so the test could be tighter. Not tautological — it catches the original
   defect — but it would not catch a future divergence where JSON has a non-null
   but wrong status.

**Are any tautological?** No. Every test constructs a fixture, calls the function or
command under test, and asserts a specific outcome that would fail if the behavior
regressed. None test only the test's own setup.

**Would they catch a regression?** Yes for the cases they cover. The gap is what
they do NOT cover: the legacy §2 + machine roster combination (CRITICAL above), and
the JSON-status-value-specificity weakness in test 5.

**DF-2 note:** the consensus deferred (DF-2) the full G1 acceptance test shape
(create-run → mutate-config → continue → prove-unchanged). The shipped
`TestSnapshotPinsPerRosterIDNotPerAdapter` and
`TestSnapshotPinsAutonomousLaunchArgs` are unit tests of `applyRosterSnapshot`, not
end-to-end tests of the wired `continueAuto` path. This is the deferred shape, not a
new gap.
