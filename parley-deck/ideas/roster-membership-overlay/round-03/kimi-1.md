---
agent: kimi-1
idea: roster-membership-overlay
round: 3
date: 2026-08-20
responding-to: [user/path-c-ruling, claude-1/path-c-measurement, codex-1/signoff-block]
---

# roster-membership-overlay — kimi-1, round 3: the experiment, run

Everything below was executed in an isolated copy at `/tmp/kimi-r3-e1.9A8zCK`
(`internal/`, `cmd/`, `go.mod`, `go.sum`, `parley-deck/agents.toml` + `COOPERATION.md`,
plus `VERSION` and a copy of `.gomodcache`; `GOCACHE`/`GOMODCACHE` pointed inside the
copy). The shared working tree was never written; `git status --porcelain` on it is empty
(checked after all runs). Two binaries were built from the same copy, differing only by the
patch: `parley-shipped` (= 1.45.0, unmodified source) and `parley-patched`. Full logs:
`/tmp/kimi-r3-e1.9A8zCK/census2/`. Unless tagged otherwise, every claim is PRIMARY (ran or
read this session).

**The patch** (`internal/config/runtime.go`, inside `LoadRosterScoped`'s membership-layer
branch, shipped :166-173). One rule, chosen before any result was known: **a committed deck
block declares membership iff it binds the ID to a family — non-empty `adapter`.** A block
carrying only values no longer constitutes membership; its values still merge, and
membership falls through to legacy §2 or the machine roster.

```go
// EXPERIMENT E1 (Path C probe, NOT a shipped rule): a committed deck block
// declares membership only when it binds the ID to a family — non-empty
// `adapter`. ...
for _, id := range ids {
    if strings.TrimSpace(entries[id].Adapter) == "" {
        continue
    }
    deckMembers[id] = true
    deckActive[id] = entries[id].Active == nil || *entries[id].Active
}
if deckSource == "" && len(deckMembers) > 0 {
    deckSource = item.source
}
```

Why `adapter` as the key, stated before the data: it is the one field a values-only
override of an inherited member never needs to restate (it is inherited by layering), and
the one field a genuinely new member cannot lack (else the row is `unmapped` and cannot
launch). And it is the field D-A's written block lacked. My signoff set the bar for this
cheap shape: "the chosen key provably never reclassifies an existing block" — measured
below (E1c).

## E1 result

### a) Six members with kimi-1 at speed=fast — YES (PRIMARY)

Fixture: machine home with six `[roster.*]` members (adapter each, kimi-1 `speed =
"deep"`); deck `parley-deck/agents.toml` containing exactly:

```toml
[roster.kimi-1]
speed = "fast"
```

```
$ PARLEY_HOME=fix-home ./parley-shipped roster show --dir fix-deck
AGENT        ADAPTER    STATE    ...  SPEED    AUTO STATUS
kimi-1       kimi       active   ...  fast     yes  model-unbound,effort-from-config,metadata-unknown
```
One member. The defect, reproduced as control.

```
$ PARLEY_HOME=fix-home ./parley-patched roster show --dir fix-deck
claude-1     claude     active   ...  deep     yes  inherited-roster
codex-1      codex      active   ...  deep     yes  inherited-roster,...
hermes-1     hermes     active   ...  deep     yes  inherited-roster,...
kimi-1       kimi       active   ...  fast     yes  inherited-roster,model-unbound,effort-from-config,...
opencode-1   opencode   active   ...  deep     yes  inherited-roster,...
zcode-1      zcode      active   ...  deep     yes  inherited-roster,model-from-config,...
```

**Six members, kimi-1 at speed=fast, every row marked `inherited-roster`.** Provenance is
already truthful without any new display work:

```
$ ./parley-patched roster show --dir fix-deck --explain kimi-1
kimi-1 — membership from ~/.parley/agents.toml (INHERITED — this deck declares no roster of its own)
FIELD          EFFECTIVE                SET BY
adapter        kimi                     ~/.parley/agents.toml
speed          fast                     parley-deck/agents.toml
active         active                   ~/.parley/agents.toml
```

That is the owner's requirement working, verbatim: the child overrode one property;
membership was never mentioned, so it was inherited.

The same result through the actual CLI gesture, not a hand edit: on a block-less
inheriting deck, `./parley-patched roster set zcode-1 --scope deck --speed fast --yes
--confirm-breaking` wrote `[roster.zcode-1]\nspeed = "fast"` and the subsequent show
listed six inherited members with zcode-1 at speed=fast. (`--confirm-breaking` was
demanded by the stale gate — see b.) Shipped control on that same committed file:
`zcode-1` alone.

### b) `go test ./...` — NOTHING breaks; and that is partly a coverage gap (PRIMARY)

First run failed only `TestVersionFileMatchesBinaryVersion` (`open ../../VERSION`) — my
copy omitted `VERSION` (the brief's copy list doesn't include it); artifact of the copy,
not the patch. With `VERSION` added: **all 28 packages pass, `internal/app` 45.2s,
`internal/config` 0.4s included.**

Classification: zero genuine regressions AND zero old-rule assertions broke — because
**no test ever writes a values-only block into a committed deck file**. I grepped every
`[roster.<id>]` fixture in `internal/app` and `internal/config` tests: all committed-deck
blocks carry `adapter`; the only adapter-less blocks live in `agents.local.toml` and the
env layer (`roster_cycle2_test.go:248,306`), which never carry membership under either
rule. `LoadRosterScoped` has no direct unit test in `internal/config` at all (grep:
`loadrosterscoped|rosterscope|members` over `internal/config/*_test.go` → nothing). So the
green suite proves the patch doesn't disturb any PINNED behavior; it cannot prove the new
branch is safe, because the branch is entered by no test. The two real caller-level breaks
I found, I found by running, not by the suite:

**(b1) The D-A gate inverts its lie.** Under the patched resolver, on a pure-inheritance
deck:

```
$ ./parley-patched roster set zcode-1 --scope deck --speed fast --yes
roster set: this adds a new roster member — a membership change, not a settings change.
Re-run with --confirm-breaking as well as --yes.
```

The resolver says membership is UNCHANGED (six inherited before and after); the gate
(`membershipChange`, `roster_set.go:287-290`, keys on block existence in the file) still
claims a member is added and demands `--confirm-breaking`. Under Path C the gate is no
longer misleading in the old direction — it is now wrong in BOTH directions, and wrong on
the exact gesture the owner wants to be easy. Re-keying the gate to the resolver's
before/after member sets (the already-agreed §1.2 D-A shape) stops being cosmetic and
becomes **release-blocking**.

**(b2) `roster sync` becomes a silent membership verb.** Five-block adapter-only deck
(the fleet's dominant shape), machine has the same adapters plus zcode-1:

```
$ ./parley-shipped roster sync --dir fix-deck2 --yes   # then show
Wrote .../agents.toml — 5 redundant override(s) and 0 deliberate pin(s) removed; the deck now inherits.
→ members after: claude-1,codex-1,hermes-1,kimi-1,opencode-1        # STILL 5 (D-C lie, reproduced again)
$ ./parley-patched roster sync --dir fix-deck2 --yes   # identical file and output text
→ members after: claude-1,codex-1,hermes-1,kimi-1,opencode-1,zcode-1  # 5 → 6
```

Same command, same bytes written (adapters stripped, empty headers left), same preview
text. Under the adapter key, the member set changed 5→6 with **no membership intent
expressed anywhere in the preview** — sync's ratified contract is "value rebase,
membership survives" (`roster_sync.go:13-24`; `roster_sync_test.go:35-40` requires the
headers to survive, and that test still passes while the membership meaning of those
headers silently changes). The D-C message accidentally becomes true while the actual
change becomes invisible. This is the single strongest measurement of this round: **a
content-keyed membership rule makes the value-rebase verb a quorum-changing verb, unless
`adapter` is exempted from sync's rebase.**

### c) Fleet census, patched vs shipped — ZERO decks change (PRIMARY)

Method: `find "/Volumes/My Shared Files/AI_WORKSPACE" -type d -name parley-deck
-not -path '*/.git/*'` → **49 deck dirs** (scope note: no depth cap, includes 4 worktree
copies of this repo, nested duplicate deck dirs, and `.r16-kimi-scratch` fixtures; my
round-2 census used a different root and maxdepth 4, hence 41 — same population,
different frame). For each workspace: `roster show --dir <ws> --json` with BOTH binaries
against the REAL machine config (no `PARLEY_HOME` override), then compared the sorted
member set and the sorted active set. Script: `/tmp/kimi-r3-e1.9A8zCK/census2.sh`; data:
`census2/compare.tsv`.

```
=== DIFFS (member set or active set differs shipped vs patched) ===
=== END DIFFS ===            ← empty: 0 of 49
=== decks where EITHER binary errored ===
=== END ERRORS ===           ← empty: 0 of 49
```

**How many decks change their active member set: ZERO. Names: there are none.** Spot rows
(identical under both binaries): BYTE `claude-1,codex-1,hermes-1,kimi-1,opencode-1`;
Finance (carries two `active=false` blocks) members
`claude,claude-1,codex,codex-1,hermes-1,kimi-1,opencode-1` / active
`claude-1,codex-1,hermes-1,kimi-1,opencode-1`; this repo (`parley-deck/parley-deck-cli`)
all six, inherited, both binaries.

Corroborating block-CONTENTS census (every prior census, mine included, counted blocks,
not fields; my signoff made this census part of the experiment — script
`contents_census.py`): **38 decks carry `[roster.*]` blocks; every block — 100% — carries
`adapter`. Zero adapter-less committed blocks exist on the volume.** Shapes:
`adapter`-only (dominant; 24 decks purely so), `adapter+active` (13 decks, partial),
`adapter+effort` (igm-app), `adapter+effort+model` (mcp_anywhere). Decks with ≥1
adapter-less block: **0**.

So the ruling's "37 quorums silently change" hazard is real for a content-blind reading
of C ("blocks stop declaring membership" — under that rule all 38 flip) and **measured
zero for the adapter key**, because every fleet block carries the family binding. The
hazard was the argument against guessing the key from any-value presence; it is not an
argument against this key. My signoff's dead-on-arrival condition ("if the census finds
even one existing deck whose membership depends on adapter-less blocks, the cheap shape
is dead") did not trigger.

### d) The rule I would ship

**Membership iff the committed deck block carries a non-empty `adapter`** — with three
ride-alongs the experiment proved are release-blocking, not polish:

1. **D-A gate re-keyed to the resolver's before/after member sets** (b1). Under the new
   rule a values-only write needs no `--confirm-breaking`; an adapter-adding write does.
   This is the §1.2 shape, now mandatory.
2. **`roster sync` must exempt `adapter` from rebase removal** (b2). `adapter` is the
   seat-binding, not a value; stripping it must not be a field operation. With the
   exemption, sync's contract ("membership survives") becomes true by construction, and
   the D-C message fix reduces to honest value-reporting. (codex-1's stronger D-C shape —
   state value-vs-membership outcomes and post-resolve — I adopt as the acceptance test.)
3. **Define committed-deck `active` on an inherited member.** Today `applyAuthorityState`
   (`runtime.go:209-221`) discards it. Under the owner's model, `active` is a property and
   the committed child may override it: `[roster.zcode-1] active = false` on an
   inheriting deck should retire zcode-1 from the quorum — visibly (STATE=inactive,
   `--explain` attribution), reversibly, in the committed record. That is codex-1's
   `remove`, arrived at from the owner's model exactly as the ruling predicted, with zero
   new syntax. (The ratified non-layering of `active` targeted the gitignored/env layers —
   `runtime.go:127-134` — and must stand: those stay unable to touch quorum state.) If the
   idea does not take this now, the alternative is to reject `active`-only committed
   blocks at write time; what must not ship is today's silent discard.

Why this beats the alternatives:

- **`schema = 2` marker:** its only function is insurance against reclassifying existing
  files. Measured reclassification under the adapter key is zero (c), so the marker buys
  nothing on this fleet while permanently splitting membership into two regimes — and an
  opt-in marker makes C non-default, which the ruling explicitly excludes ("C cannot
  avoid it by being opt-in without ceasing to be the default model"). codex-1's
  compatibility boundary (§1.5) was written for an overlay ADDING a second semantics;
  here the resolver change is one semantics, measured byte-safe.
- **Explicit `members = [...]` list:** the owner's literal model, but it is new syntax, a
  new coherence problem (dangling value blocks for non-members), and it still needs the
  attended migration for 38 existing files. It expresses nothing the adapter key doesn't
  already express with today's bytes. If the fleet ever shows demand for list arithmetic
  (`super.members ± x`), the object model already answers: `add` = a block with `adapter`
  for an ID the parent lacks; `remove` = `active = false` (ride-along 3). The demand gate
  my (a) required is preserved in reduced form: those two gestures are now expressible,
  so any future syntax is sugar, and sugar needs demonstrated demand.

## E2 result

**Does Path C extend beyond the roster? Measured: yes, almost everywhere already —
`[roster.*]` membership is the only wholesale exception; three narrower non-uniformities
exist that claude-1's four-field check could not reach.**

Empirical (PRIMARY, shipped binary, fixture: machine sets `model`, `sandbox_mode`,
`isolate_home`, `timeout_ms`, `notes` for `[agents.kimi]`; deck sets only
`approval_policy`, `speed`, `timeout_ms`):

```
$ PARLEY_HOME=fix-home-e2 ./parley-shipped agents list --dir fix-e2
kimi  yes  0.36.1  headless  configured  read-only  never  kimi-code/k3-machine  222000  yes  yes  hosted
  sources: sandbox=~/.parley/agents.toml approval=parley-deck/agents.toml model=~/.parley/agents.toml timeout=parley-deck/agents.toml
```

Six fields, three sources, resolved independently — including `isolate_home` (HOME=yes,
inherited from machine; not one of claude-1's four and not printed in the sources line)
and `speed` (deck). Code-verified for the remaining ~24 fields (PRIMARY,
`runtime.go:759-922`): every `[agents.*]` field is presence-gated per field; the
pointer-typed supervision knobs additionally distinguish explicit `0` from absent.

Properties that do NOT layer per field, named:

1. **`[roster.*]` membership** — the exception this idea exists to remove (E1).
2. **`[roster.*] active`** — deliberately authority-bound, ratified (`runtime.go:127-134`,
   `209-221`): non-authority layers' `active` is discarded. See E1d ride-along 3 for the
   committed-deck question C reopens.
3. **Integer fields gated `> 0` cannot be overridden TO zero**: `timeout_ms`,
   `interactive_timeout_ms`, `interactive_poll_ms` (`runtime.go:829-838`). Absence and
   zero conflate, so a deck cannot express "no timeout" over a machine that sets one. The
   pointer-typed supervision knobs (`first_event_timeout_ms`, `stall_timeout_ms`,
   `heartbeat_ms`) are the in-codebase repair pattern for exactly this. This is the one
   genuine gap in "any global property is overridable" that I found beyond membership.
4. **Composite properties replace wholesale per layer** — `commands`, `headless_args`,
   `acp_args`, `interactive_args`, `isolated_home_env` (`runtime.go:765-828, 901-907`): a
   higher layer replaces the whole list/map, no per-element merge. That IS the object
   model's rule for a list property (override replaces the value), so it is consistent
   with C — but it must be documented as such, because operators expect env maps to merge.

Checked and layering per the model (so nobody has to re-verify): `[defaults]` scalars and
`[defaults.timeouts]` per field (`mergeDefaults`, `runtime.go:534-574`); `[defaults.loop]`
presence-aware pointers; `[defaults.track_rosters]` per track key; `[rosters.*]` presets
per preset-name, whole-list replace (`internal/config/roster.go:31-64`). `[agents.*]`
block EXISTENCE unions across layers (a deck can add a family the machine lacks,
`runtime.go:709-729`) — under C that is a child declaring a new property, which the
owner's model permits; not an exception.

## Position under Path C

I signed (a). The ruling rejects (a)'s premise — I was defending an authority model in
which membership is categorically different from values, and the owner has ruled that
membership is a property. I accept the ruling, and I note what does NOT survive contact
with the measurements: my (a) rested on "zero measured demand for membership deltas" and
on the authority invariant being what made D-A diagnosable. The demand question is
settled by the ruling itself (the owner is the demand, and the demand is the VALUE case,
exactly as claude-1's reframing said). The diagnosability argument survives but in
weakened form — see the fourth objection below.

What I now build, replacing both my (a) and codex-1's (c): **the adapter-keyed resolver
rule, the three ride-alongs of E1d, D-B first, then the attended `roster inherit` verb**
(Migration below). No `[membership]` stanza, no schema marker, no add/remove syntax.
codex-1's block-condition is met by measurement: the experiment ran, the six-member
fixture stays six after a committed speed override (E1a), unmarked full-declared fixtures
retain their exact member sets (E1c, 49/49), and the D-A/D-C honesty requirements are
ride-alongs 1-2.

Defects in C the owner has not seen, stated plainly per the ruling's invitation:

1. **C as literally worded ("a deck block which declares no membership intent") has no
   visible marker of intent in the file.** My content rule infers intent from the
   `adapter` key; that is an inference, measured safe on THIS fleet today, but it is a
   rule a reader must know in order to read a deck file. I judge the one-sentence rule
   ("a block that binds a family declares a seat") plus truthful `--explain` (already
   shipped, shown in E1a) sufficient; the alternative readings of C either reclassify 38
   decks (content-blind) or make C opt-in (marker). This is a judgment call the FINAL
   must state, not hide.
2. **Measured, not hypothetical: C makes `roster sync` a membership verb** (b2). Shipping
   C without ride-along 2 converts D-C from a false message into a silent quorum change.
3. **Measured: C without the gate fix is worse than the status quo on the target
   gesture** (b1) — the owner would be asked to `--confirm-breaking` to change one speed.
4. **The invariant my (a) defended — one committed file answers "who deliberates" —
   genuinely weakens.** Under C, membership is answered by the merge of two files, and a
   deck's negative space (what it doesn't override) carries quorum semantics. A `roster
   show` reader must understand `inherited-roster` to know where the quorum came from.
   `--explain` already answers provenance per field, and the committed record still
   answers "what does this deck decide differently" — I judge the diagnosability cost
   paid back by the removal of the 37 stale full copies' false authority. But it is a
   real cost, and it belongs in the FINAL's tradeoff section, not in a footnote.

## Migration

Measured first: **the 38 block-carrying decks need no file change and no quorum change —
E1c shows their member and active sets identical under the patch, byte-for-byte files
untouched.** The ruling's compatibility hazard, for the adapter key, is zero today. So
migration is no longer "37 quorums must move"; it is "the model changes under decks that
choose to move." Concretely:

1. **D-B lands first** (unchanged from my round-2 batching; renderer, drift anchor and
   embedded default atomically, codex-1's render-then-guard acceptance). Every step after
   this one writes or displays roster state; §2 writes are unsafe until it lands.
2. **The resolver change ships with all three ride-alongs in one release.** The gate
   tells the resolver's truth (1), sync cannot change membership (2), `active` on an
   inherited member is defined and visible (3). This release also fixes D-C's message per
   codex-1's strengthened shape.
3. **Decks that want to move from declared to inherited get one explicit, attended verb:**
   zcode-1's round-1 `roster inherit`, now with a well-defined mechanism under the
   adapter key — remove the deck's adapter bindings (the blocks themselves, if
   values-free), preview the resolver's before/after member sets, bind apply to the
   preview, verify `Inherited` after the write. Per-deck, git-first: hermes-1/zcode-1's
   dirty-tree and non-git counts (18/26, 15/41 — SECONDARY, their censuses) still gate
   bulk operation. Never `roster sync` (its contract is values; ride-along 2 makes that
   enforceable).
4. **No deck is moved by default, and nothing is inferred from omission** (§1.3 stands).
   A deck that keeps its five adapter blocks keeps its five members — that is now a
   stable, even self-recommending state (its blocks visibly bind families), not a stale
   accident.
5. **The instrumented demand count survives, reduced.** Each attended `inherit` records
   the deck's reason. If ≥2 decks ask for machine±1-while-tracking, the add/remove
   question reopens — but under C those are already expressible (E1d), so the count now
   measures demand for syntax sugar, and the bar for spending a contract change on sugar
   stays high.

## What I would sign

- The adapter-keyed membership rule as the Path C implementation, with E1d's three
  ride-alongs release-blocking, D-B first, `roster inherit` attended. (PRIMARY basis:
  E1a/E1b/E1c above.)
- E2's boundary list recorded in the FINAL: membership (removed by this change), roster
  `active` (ride-along 3), the `> 0`-gated int fields that cannot be overridden to zero,
  wholesale composite replacement. No other `[agents.*]` or `[defaults]` property was
  found non-layering.
- consensus.md §1 in full; my signoff's two attribution corrections (§1.4 round-1
  migration credit; value-case trigger is zcode-1's) still stand against the record.
- I withdraw (a) as the end-state — it remains the correct record of the pre-C decision
  — and sign C-with-ride-alongs. codex-1's block was the correct call: this experiment
  had to run before FINAL, and its result (zero fleet flips, the sync hazard, the gate
  inversion) is not what either side of (a)/(c) predicted in full.
- I do not sign: any C implementation keyed on block CONTENT other than the adapter rule
  without re-running E1c against that key; any release that lets `sync` strip the
  membership-keying field; any silent discard of a committed-deck `active` override; and
  any FINAL that presents the zero-flip result as guaranteed rather than measured on this
  volume (49 dirs, one SMB share, scopes stated — the caveat travels with the number).

**Verification basis.** All command output quoted above: PRIMARY, ran this session in
`/tmp/kimi-r3-e1.9A8zCK` (shipped and patched binaries built from the same copy differing
only by the patch; fixtures and scripts under that dir; census data
`census2/compare.tsv`). Code locators — `runtime.go:122-221, 534-574, 759-922`;
`roster_set.go:287-290`; `roster_sync.go:13-24, 76-81, 169-170`; `roster.go:280-302`;
`roster_view.go:29-50`; `internal/config/roster.go:31-64`: PRIMARY, read this session.
Fleet/dirty-tree figures from hermes-1/zcode-1: SECONDARY, credited. No RECALL claims
relied on. Shared tree: `git status --porcelain` empty after all runs.
