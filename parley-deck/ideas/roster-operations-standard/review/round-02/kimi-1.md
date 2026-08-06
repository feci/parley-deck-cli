---
idea: roster-operations-standard
phase: 8 — re-review
agent: kimi-1
round: 2
date: 2026-08-06
reviewed-commit: 57fe9d7
verdict: FINDINGS
---

# Re-review — kimi-1

## Verdict

**FINDINGS** — two MAJOR, three MINOR. `PRIMARY` unless tagged otherwise.

Cycle 2 is materially better than cycle 1: fourteen of sixteen agreed fixes landed and I
verified each behaviorally against scratch decks under `/tmp` with `PARLEY_HOME` isolation
(binary built from 57fe9d7 to `/tmp`; both repositories untouched — `git status --short`
clean in the CLI repo; the skill repo's pre-existing uncommitted edits are the
implementer's, see F-4). `go build ./...`, `go vet ./...` and the full `go test ./...`
are green at 57fe9d7 (ran them myself). But the adversarial premise of this re-review was
vindicated twice more:

- **F-1 (MAJOR):** A1's adopted legacy-fallback clarification is recorded in
  IMPLEMENTATION.md as landed and is **not implemented** — a deck whose only roster is a
  valid legacy §2 table still gets the machine roster inherited over it, quorum included.
  Same defect class as cycle 1's "documented as done, not done".
- **F-2 (MAJOR):** the new `masked-by-env` emitter (A12) false-positives on **every**
  machine-scope write — it warns that the file just written is masking itself. A new
  defect introduced by this cycle, in the one emitter that shipped without a test.

**§15.1 ownership disclosure.** The A1 resolution adopted my round-01 M4 position; A12's
post-write re-resolve was my round-2 adopted position; R4 (`--keep`) was mine. I issue no
verdict on those positions. What I verdict below are fresh implementation facts about code
that first exists at 57fe9d7 — gathered by PRIMARY evidence (command output quoted, file:line
cited). codex-1 and hermes-1 verdict independently.

## Per-fix verification (A1–A16, DF-1 guard)

**A1 (CRITICAL) — core FIXED; adopted legacy clarification NOT honored (F-1).** `PRIMARY`.
Deck declaring 2 members with a 5-member machine roster now resolves to exactly 2:

```
$ parley roster show --dir /tmp/kr/deckA        # deck: claude-1, kimi-1; machine: 5
claude-1     claude     active   yes       claude-opus-5  ... ok
kimi-1       kimi       active   no        kimi-code/k3   ... not-installed,effort-unknown
```

Participant selection reads the same authority: both call sites use `RosterMembership`
(app.go:1826, app.go:2454), which consumes `config.LoadRosterScoped` (roster.go:697-708)
— membership = deck file, values still layer through (verified: claude-1's undeclared
model inherits the machine value while machine-only members vanish). A rosterless deck on
a rostered machine marks every row `inherited-roster` (verified), and `roster render`
refuses: exit 1 with "this deck declares no roster of its own; the 5 rows shown come from
~/.parley/agents.toml …" (roster_render.go:40-47), while `--adopt-inherited` permits it
(verified, preview shown). **However** — see F-1: the codex-1 condition-2 clarification
(a valid legacy §2 is compatibility membership, machine roster NOT inherited over it) is
not implemented anywhere.

**A2 — FIXED.** `PRIMARY`, behavioral. Deck with TOML roster + §2-only `ghost-1`:
`roster show` reports `ghost-1 … unmapped,section2-only` with an explanatory note
(roster_view.go:56-84); the generated §2 excludes it (never auto-added); `roster render`
names the removal in the preview ("the following §2 row(s) are NOT in this deck's roster
and will be removed: - ghost-1 …") **and** again on apply; the applied file loses the row;
the second apply is a byte-stable no-op ("§2 already matches the roster — nothing to do").

**A3 — FIXED, adversarially probed.** `PRIMARY`. All four bypass shapes now refuse with
exit 2 and write nothing: `--model` only, `--effort` only, `--state active` only, and
`--adapter` only, each on a non-existent block. The gate keys on block existence
(`blockExists`, roster_set.go:112-124, checked against the same bytes the patcher edits).
Controls: an existing member's `--model` change applies with `--yes` alone (no
over-gating); a new member with `--adapter --yes --confirm-breaking` applies; a set into
a deck with no `agents.toml` yet is gated. Nothing half-written in any refusal.

**A4/A5 — FIXED.** `PRIMARY` (code-read + shipped unit tests). Snapshot entries are keyed
by `Agent` with a first-writer-wins adapter fallback for pre-ID runs
(roster_snapshot_apply.go:33-41); `LaunchArgs` carries the resolved argv
(manifest.go:73-77, captured via `spec.ResolveLaunchArgs()` at roster.go:386, pinned back
at roster_snapshot_apply.go:62-67), so a machine-config change dropping an auto-approve
flag can no longer change a continuation's autonomy posture. `TestSnapshotPinsPerRosterIDNotPerAdapter`
and `TestSnapshotPinsAutonomousLaunchArgs` exercise the production function and would
catch a regression. One hole noted as NIT-1 below (revision hash excludes the args).

**A6 — FIXED.** `PRIMARY`, behavioral. `--scope machine` vs `--scope deck` now return
different, correct views (5 machine rows vs 2 deck rows on the same deck); a bogus scope
is rejected with exit 2. `--all` appends `not-in-roster` rows for configured-but-unclaimed
adapters (the opencode-invisibility answer). `--explain claude-1` prints per-field
provenance ("adapter … parley-deck/agents.toml / model … ~/.parley/agents.toml / effort …
built-in default") and exits 1 for a non-member. `roster init` accepts `--scope deck`,
prints the deprecation notice, and `session` is no longer the visible default in flag
help (roster.go:96).

**A7 — core FIXED; one clause unaddressed (F-3).** `PRIMARY`. A healthy row now renders
`ok` in text and `"status": ["ok"]` in JSON (verified side by side); the payload carries
`schema_version`, ordered `columns`, and the new `scope` key. The out-of-contract fields
were neither removed nor relocated: `display_name` still serializes on every resolved row
(and `note` when present). Details in F-3.

**A8 — FIXED.** `PRIMARY`, behavioral. Machine config with `headless_args = ["--print",
"--model", "ancient-literal", "--effort", "low", …]` plus deck-declared `model`/`effort`
now resolves the row to the declared values with status `ok` — pre-fix this was
`model-drift` with the literal winning. `NormalizeLegacyModelArgs` (launchargs.go:131-168)
rewrites only the value, leaves flags/positions alone, and the rewrite is recorded in
`Sources["headless_args_normalized"]` (runtime.go:727-741).

**A9 — text removed from all four surfaces; assertion covers two of them (F-4).**
`PRIMARY`. `rg` for all four banned phrases across both repositories' working trees: zero
hits. The drift assertion `TestNoSection2AsAStoreInstructions` (drift_test.go:257-282)
pins the phrases in the embedded default and the live deck — and I ran the negative
control: injecting "Fill in §2 roster" into a `/tmp` copy of the embedded default makes
the test FAIL as designed. The skill repo's bundled copy and SKILL.md have no equivalent
assertion, and the skill-side edits are uncommitted — F-4.

**A10 — FIXED, verified end-to-end.** `PRIMARY`. With a fabricated session index and
manifest under scratch dirs: a manifest whose frozen revision matches the deck reports
`Roster snapshot: current`; mutating the deck's effective model flips it to
`stale-snapshot`; a snapshot-less manifest reports `no-snapshot`. Text and JSON
(`roster_snapshot_state`) agree. The comparison lives at roster_view.go:232-244, wired
into `inspectSession` at app.go:956-960. (I replicated `RosterRevisionOf`'s hash in
Python to construct the matching case — the "current" result confirms the hash inputs are
the snapshot fields, no hidden salt.)

**A11 — FIXED.** `PRIMARY`. protocol-changelog.md:119-124 now carries the §7 template:
`Idea: ideas/roster-operations-standard/`, `Drafted by: claude-1`, `Summary: …`, with the
substantive prose preserved beneath.

**A12 — landed for deck scope, broken for machine scope (F-2).** `PRIMARY`, behavioral.
Deck-scope write masked by `agents.local.toml` correctly warns ("model = \"brand-new-model\"
is MASKED — parley-deck/agents.local.toml sets it at a higher layer …"), and an unmasked
deck write stays silent. But **every** machine-scope write false-positives — see F-2.

**A13 — FIXED.** `PRIMARY`. `parley --help` lists all five verbs and the `agents list`
relabel ("adapter/runtime inventory — NOT the roster", also printed by `agents list`
itself, discover.go:490). Docs went from 0 roster mentions to 18 (cli-reference.md) and 9
(agent-runtime-configuration.md); I read the new cli-reference section — it documents the
membership rule, statuses, and flags accurately (its `legacy-roster` row, :77, actually
describes the AGREED semantics the code doesn't implement — see F-1).

**A14 — FIXED in the skill working tree, uncommitted (F-4).** `PRIMARY`. The SKILL.md
diff deletes "keeps working … until `parley roster sync` moves it across" and replaces it
with "roster sync does **not** migrate it … The remediation is `parley roster migrate
--backup-dir DIR --dry-run` … or `parley roster set …`, then `parley roster render`" —
factually correct about released behavior. The membership-is-the-deck rule, the three new
status codes, and the new flags are documented. But `git status` in the skill repo shows
`M skills/parley-deck/SKILL.md`, `M skills/parley-deck/references/COOPERATION.md`,
package.json still reads 2.5.0, and CHANGELOG.md's latest entry is 2.5.0 — the agreed
"skill 2.5.1" does not exist yet.

**A15 — FIXED.** `PRIMARY`, behavioral. `--keep kimi-1.modle --yes` exits 2, names the
unmatched token, writes nothing, and the pin survives on disk. A correct `--keep
kimi-1.model` exempts the pin while the rebase removes the rest. The preview/apply binding
re-reads the deck file and refuses if any to-be-dropped field's value changed since the
preview (roster_sync.go:145-160) — verified by code-read; the residual theoretical gap
(an edit to a *non-dropped* field landing between the two same-invocation reads) is
inside what the consensus asked for. `--drop-pins` correctly absent.

**A16 — mostly FIXED; attribution item NOT done (F-5).** `PRIMARY`. File mode: a
machine-scope write to a `0644` file now preserves `0644` (stat-verified;
writeRosterFileAtomic → fsutil.WriteFileAtomic with the target's mode,
roster_set.go:243-254). Continuation on an unreadable manifest now warns on stderr
(app.go:1176-1182, code-read). Stale guidance: the `unmapped` note now points at `roster
set` (roster.go:315), and `init` carries the deprecation notice (verified). modelmeta:
`kimi-k2-thinking` now classifies as family "Kimi" (the `kimi` rule is reachable; the
broad `k` rule moved last, modelmeta.go:64-69). Reactivation messages "reactivates a
retired roster member" and is gated (verified). **Not done:** CHANGELOG.md:7 still reads
"defect in 1.40.0, found by codex-1 and hermes-1" — the agreed credit correction did not
land (CHANGELOG.md is untouched by 57fe9d7).

**DF-1 interim guard — FIXED.** `PRIMARY`, behavioral. `roster migrate --yes` without
`--confirm-breaking` exits 2 with the attended-only refusal (roster_migrate.go:62-69);
dry-run is unaffected; a dirty-tree deck is skipped and reported ("working tree has
uncommitted changes; commit or stash first …") while a clean deck in the same run migrates
— verified with two git-initialized scratch decks. The guard's deliberate choice to treat
a non-repository as clean is stated in deckTreeDirty's comment (roster_migrate.go:379-387)
and is reasonable for an interim guard.

## Findings (by severity)

### F-1 — [MAJOR] A1's adopted legacy-fallback clarification is not implemented: a valid legacy §2 is still overridden by the machine roster

`PRIMARY`. Consensus A1 (codex-1's condition 2, **adopted**): "'a deck with no roster of
its own' means **neither** a deck-level `[roster.*]` block **nor** a valid legacy §2
table. A deck carrying a valid legacy §2 keeps that table as its compatibility membership
— reported `legacy-roster` — until it is migrated; the machine roster is not inherited
over it." IMPLEMENTATION.md:15-16 records this as done ("Legacy fallback per codex-1's
condition …"). It is not done. `LoadRosterScoped` (runtime.go:119-160) consults only
config layers — it never looks at §2 — and `resolveRoster` triggers the legacy path
solely on `legacy := len(scope.Members) == 0` (roster.go:268). Behavioral proof, deck
with a valid §2 table (`claude-1`, `kimi-1`) and NO config roster of its own, on a
5-member machine:

```
$ parley roster show --dir /tmp/kr/deckC
claude-1 … inherited-roster
codex-1  … inherited-roster,effort-unknown
hermes-1 … inherited-roster,effort-unknown
kimi-1   … inherited-roster,not-installed,effort-unknown
opencode-1 … inherited-roster,model-unbound,effort-unknown,metadata-unknown
```

The §2 table never appears; no row is marked `legacy-roster`. Because `RosterMembership`
(roster.go:697-708) reads the same scoped view and falls back to
`protocol.ReadRosterIDs` only when `scope.Members` is empty, a **run on this deck selects
the five machine members as its quorum** — the exact outcome the clarification was adopted
to prevent, for exactly the population it named (not-yet-migrated legacy decks on rostered
machines). `roster render`'s refusal message ("this deck declares no roster of its own")
is also wrong for such a deck — it does declare one, in §2. The new documentation even
states the agreed semantics the code contradicts: docs/cli-reference.md:77 defines
`legacy-roster` as "this deck has no `[roster.*]` block; §2 is the fallback". The fix is
small and localized: before inheriting machine members, consult
`protocol.ReadRosterIDs(root)`; if a valid §2 exists, treat membership as legacy (show:
`legacy-roster`; membership: the §2 sets; render: keep refusing, with a §2-aware message).
Severity MAJOR not CRITICAL because the CRITICAL's core case (a deck TOML roster being
unioned away) is genuinely fixed; this is the same defect class confined to the
legacy-deck path — but it is the adopted text of A1, claimed as landed, and it is not.

### F-2 — [MAJOR] NEW in cycle 2: `masked-by-env` false-positives on every machine-scope `roster set`

`PRIMARY`, behavioral:

```
$ parley roster set claude-1 --scope machine --dir /tmp/kr/deckF --model even-newer --yes
Wrote /tmp/kr/home3/agents.toml

warning: model = "even-newer" is MASKED — ~/.parley/agents.toml sets it at a higher layer, so the effective value did not change.

$ parley roster show --dir /tmp/kr/deckF --scope machine | grep claude-1
claude-1  claude  active  yes  even-newer  …        ← the write DID take effect
```

The warning names the very file just written as its own masker. Root cause (code-read):
`RosterFieldSources` returns display labels, and the machine layer's label is the literal
string `"~/.parley/agents.toml"` (runtime.go:325), while `rosterFieldMaskedBy`
(roster_set.go:96-105) compares that label against the resolved absolute `target` path
with `HasSuffix`/`Contains` — never equal, so `masked=true` unconditionally for machine
scope. Deck scope works only because the label `"parley-deck/agents.toml"` happens to be a
path suffix. A12's agreed emitter therefore fires at the wrong time on a primary verb:
every routine machine-scope change (the normal way to change your default model) cries
wolf, teaching operators to ignore the genuine signal the fix exists to provide. Filed
MAJOR because the message asserts the write was ineffective — an operational falsehood on
every machine-scope use — not because anything is corrupted. The obvious fix is to compare
resolved paths (or skip when `src` labels the same file as `target`). No test covers this
emitter; see the test-quality section.

### F-3 — [MINOR] A7 residual: `display_name`/`note` still serialized beyond the frozen 11 columns

`PRIMARY` (JSON output quoted in A7 above): the consensus fix read "either remove the
out-of-contract fields or place them formally in the `--explain` provenance object";
neither happened — `display_name` (and `note`, when set) still ride the JSON payload
(roster.go:205-211, `omitempty`). The divergence core (null vs `ok`) and the golden are
fixed, so this is the contract-purity residue only.

### F-4 — [MINOR] A9/A14 skill side is working-tree-only, and the drift assertion covers 2 of 4 surfaces

`PRIMARY`. (a) The bundled `references/COOPERATION.md` and `SKILL.md` edits are
uncommitted in the skill repo (`git status --short` → ` M` on both), `package.json` is
still `2.5.0`, and there is no 2.5.1 changelog entry — the consensus's "ship as … skill
2.5.1" has not shipped. The text itself is correct (verified, A14 above); it simply isn't
a release artifact yet. (b) `TestNoSection2AsAStoreInstructions` pins the banned phrases
in the embedded default and the live deck only; nothing pins them in the skill repo's
bundled copy or in SKILL.md, so on those two surfaces the contradiction *can* silently
return — the exact thing A9's "add drift assertions" clause exists to prevent. (The
bundled copy currently matches the embedded one outside the header placeholders —
diff-verified — so this is a guard gap, not a content gap.) If the skill release step
commits these files and tags 2.5.1, part (a) evaporates.

### F-5 — [MINOR] A16 attribution item not landed

`PRIMARY`. CHANGELOG.md:7 still credits the 1.40.1 defect findings to "codex-1 and
hermes-1" alone; A16's "Correct the credit line" did not happen (the file is untouched by
57fe9d7). My round-01 file independently corroborated both CRITICALs
(review/round-01/kimi-1.md:59-106), which is what the correction was to record.

### NITs

- **NIT-1** `PRIMARY` (code-read): `RosterRevisionOf` (manifest.go:82-94) hashes
  Agent/Adapter/Model/Effort/Speed/Auto/Installed — not `LaunchArgs`. Now that A4 makes
  the resolved argv part of the frozen identity, an args-only drift that leaves AUTO
  unchanged (e.g. `--output-format text` → `json`) still reports `current` from A10's
  staleness audit. Auto-affecting changes are caught (AUTO flips the hash), so this is a
  narrow residue.
- **NIT-2** `PRIMARY`: the A12 warning's parenthetical "(status `masked-by-env`; …)"
  implies the STATUS column carries the code, but no `roster show` row ever emits it —
  the vocabulary code remains table-unemittable, mentioned only in set's stderr warning.
- **NIT-3** `PRIMARY`: IMPLEMENTATION.md:75 says "Nine new tests"; eleven shipped (five in
  roster_membership_test.go, five in roster_cycle2_test.go, one in drift_test.go). Record
  accuracy only.

## Test-quality assessment

`PRIMARY`. I read all eleven new tests and ran the suite. **None are tautological** —
each exercises production code against real fixtures and would catch its target
regression: the membership test drives `RosterMembership`/`resolveRoster` against a
layered fixture and asserts both the count AND value inheritance (a union-semantics
regression fails it); the render-refusal test asserts both the refusal and the
`--adopt-inherited` escape; the sync test goes end-to-end through `rosterSync` and asserts
the exit code, the named token, AND the pin's survival on disk — the strongest of the set;
the snapshot tests assert per-ID pins survive and frozen args are restored. The A9
assertion's negative control I ran myself (injection → FAIL, quoted in A9 above).

Gaps, stated plainly:

1. **The two defects I found live exactly in untested seams.** F-2's emitter
   (`rosterFieldMaskedBy`) has no test at all; F-1's clarification has no test either —
   the membership test suite covers deck-roster and rosterless-deck cases but never a
   legacy-§2 deck on a rostered machine. Cycle 2's test investment was real, but the
   uncovered corners are where the defects are, same as cycle 1.
2. `TestMembershipGateCatchesNewBlockWrittenWithAnyField` tests `membershipChange` with a
   hand-passed `existed` flag; nothing tests `blockExists` or that `rosterSet` wires the
   two together (a regression making `blockExists` always-true would silently un-gate new
   members). I verified the wiring behaviorally (A3 above), but the test alone wouldn't
   catch its unwiring.
3. `TestJSONStatusMatchesTextForAHealthyRow` never renders the text table — it asserts
   the JSON side is non-null and trusts a shared convention for the "matches text" half.
   The two `ok` renderings are produced by different code paths (the status-normalizing
   loop in `rosterShow` vs `statusOrOK`), so cross-representation drift would not fail
   this test. The name overclaims what it pins.
4. A2's removal *report* (the preview/apply messaging), the DF-1 migrate guard, and A10's
   `rosterSnapshotState` have no automated tests; all three are behaviorally verified in
   this review, but a regression would sail through the suite.

**Regression sweep (anything NEW broken):** full suite, build, and vet green at 57fe9d7;
edge cases probed — empty-everywhere deck fails closed with a clear message, set into a
non-existent file is gated, machine scope excludes §2-only rows, `--explain` on a
non-member exits 1. The one new defect found is F-2.
