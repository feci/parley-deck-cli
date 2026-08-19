---
agent: zcode-1
idea: roster-membership-overlay
round: 3
date: 2026-08-19
responding-to: [user/path-c-ruling, claude-1/path-c-measurement, codex-1/signoff-block]
---

# Round 3 — zcode-1: the experiment, run

Everything below marked PRIMARY was executed by me today in `/tmp/zcode-r3/` copies of
`internal/`, `cmd/`, `go.mod`, `go.sum`, `parley-deck/` (plus `VERSION`, which the app test
suite reads). Two module copies: `ctrl` (unpatched) and `pathc` (patched). Both built as
v1.45.0; `ctrl` was spot-checked against the shipped `/opt/homebrew/bin/parley` on three
decks (cms, ldx, millenniumProblems) — identical member/state output (PRIMARY), so
"patched vs shipped" below is a fair name for the diff. The shared working tree was never
written; `git status --porcelain` is empty. Read-only `roster show --dir` runs against the
fleet used the real `~/.parley` as the machine layer — that is the environment these decks
actually resolve in.

## E1 result

### The patch (PRIMARY — `/tmp/zcode-r3/pathc/internal/config/runtime.go`)

Ungated, content-free Path C, faithful to the ruling's object model:

- A `[roster.<id>]` block in `parley-deck/agents.toml` no longer adds the ID to
  `deckMembers` (the deleted lines are the old `runtime.go:166-169` loop body). Blocks
  still merge as values, exactly as before (`runtime.go:144-146`, untouched).
- Membership is a property: a new top-level `members = [...]` key in the deck file is the
  explicit override; its list REPLACES the machine's list.
- Authority order becomes: explicit `members` key → else, if the deck has any
  `[roster.*]` blocks, the machine roster (Inherited=true; legacy §2 does NOT intercept —
  a deck with TOML roster blocks has left prose behind) → else valid legacy §2 (unchanged)
  → else machine (unchanged).
- `active` becomes a layered property for the layers that may speak: machine declaration
  as parent, the deck block's **explicit** `active` as override. `applyAuthorityState`
  receives `overlayActive(machineActive, deckActive)`. The gitignored `agents.local.toml`
  and env layers remain barred from `active` — the ratified quorum guard
  (`runtime.go:127-133`) is preserved; only its scope shrinks from "authority layer" to
  "machine ⊕ deck".

Build: `go build ./cmd/parley` — clean.

### (a) The owner's requirement — WORKS (PRIMARY)

Fixture: `PARLEY_HOME=/tmp/zcode-r3/mh` (copy of the real six-member machine config),
deck at `/tmp/zcode-r3/fix1` whose entire `parley-deck/agents.toml` is:

```toml
[roster.kimi-1]
speed = "fast"
```

```
$ PARLEY_HOME=.../mh ./parley-ctrl  roster show --dir .../fix1   # shipped rule
AGENT ADAPTER STATE ... SPEED STATUS
kimi-1 kimi active ... fast effort-from-config          # ONE member — the D-A collapse

$ PARLEY_HOME=.../mh ./parley-pathc roster show --dir .../fix1   # patched
AGENT ADAPTER STATE ... SPEED STATUS
claude-1 ... deep inherited-roster
codex-1 ... deep inherited-roster
hermes-1 ... deep inherited-roster
kimi-1 ... fast inherited-roster,effort-from-config     # SIX members, value overridden
opencode-1 ... deep inherited-roster
zcode-1 ... deep inherited-roster,model-from-config,...
```

Six members, kimi-1 at `fast`, everyone else at the machine default `deep`, local values
flowing, membership inherited. The second half of the owner's model also works (PRIMARY):
a deck with `members = ["claude-1", "codex-1"]` plus the same kimi block resolves to
exactly claude-1 + codex-1 under the patch (ctrl: kimi-1 only). A control deck with no
roster config at all resolves identically on both binaries — the pure-inheritance path is
untouched.

### (b) `go test ./...` — 25 packages ok, `internal/app` FAIL, 13 tests (PRIMARY)

`internal/config` itself **passes**. The 13 failures in `internal/app` sort into three
bins:

**Bin 1 — tests asserting the old rule (3).** `TestDeckMembershipIsTheDeckFileNotTheLayeredUnion`
(the old rule's eponymous test: fixture deck declares 2 blocks, machine has 4; patch
resolves 4 inherited, test demands 2); `TestActiveProvenanceAndMaskingFollowTheAuthority`
(`--explain` names `~/.parley/agents.toml` as state authority; under C that IS the
membership layer for a values-only deck — the truthfulness requirement stands, the
assertion re-points); `TestValueLayersCannotChangeMembershipState` — the interesting one:
its fixture's **machine** file declares `claude-1 active = false`, and under the patch
that retirement now propagates into the deck, because the machine is the parent of `active`.
That is Path C working as specified, not a regression — but it is a behavior change the
guard's ratifiers did not vote on (see Position). The gitignored-layer half of the guard
still holds in the patch.

**Bin 2 — genuine gap: no parent roster (6).** Fixtures with deck blocks and NO machine
layer resolve to an empty scope and silently fall back to legacy §2
(`TestDefaultRosterParticipants` case 1, `TestZcodeRowReports…`,
`TestDeckDeclaredModelNeverOverrides…`, `TestUnreadableAgentConfigFallsBackToUnknown`,
`TestJSONStatusMatchesTextForAHealthyRow`, `TestJSONRowHasExactlyTheFrozenColumns`).
Values-only blocks over a nonexistent parent roster must **fail closed** ("deck declares
value overrides but has no machine roster to inherit"), never read §2. This is kimi-1's
predicted "rule-2 fall-through" break, now demonstrated.

**Bin 3 — genuine gap: the verbs cannot see the new state (4).** All four
`TestRosterRender*` fail because `renderRosterTable` treats `scope.Inherited` as pure
inheritance and refuses without `--adopt-inherited`. A values-only deck inherits
membership but commits local values; refusing to render it as "adopting" anything is
wrong — render needs the inherited-vs-values-only distinction the patch doesn't yet
expose. Same class, measured separately (PRIMARY): `roster set` on the fix1 deck still
refuses with "this adds a new roster member — a membership change" — under the patch that
sentence is false in the opposite direction from D-A (the write would NOT change
membership). The resolver patch alone is not shippable; set/sync/render/explain all need
the state.

### (c) Fleet census — 36 of 38 decks change their active member set (PRIMARY)

Every deck under `/Volumes/My Shared Files/AI_WORKSPACE` (37 carrying `[roster.*]` blocks,
plus this repo's own block-free inheriting deck), each resolved read-only by both binaries
against the real `~/.parley`:

```
while read -r d; do
  ./parley-ctrl  roster show --dir "$d" --json > "c-$key.json"
  ./parley-pathc roster show --dir "$d" --json > "p-$key.json"
done < decks.txt
```

**36 decks change their active member set. Every single one of them gains `zcode-1`.**
No deck loses an active member; no inactive member is reactivated; every dropped member
was already inactive (verified row-by-row). The hazard is not 37 — it is **36**, and the
miss is informative: `millenniumProblems` carries blocks for all six members including
`zcode-1`, so inheritance changes its STATUS (`ok` → `inherited-roster`) but not its
quorum. Unchanged decks: `millenniumProblems` and `parley-deck-cli` itself (pure
inheritor).

Additionally, 16 decks lose retired (inactive) members from membership — their
`active = false` blocks become dangling value overrides: 7 decks drop `antigravity-1`
(ai_prezz, ecb-api, librade-algoTrader, the four wt-* decks), the rest drop pre-composite
bare IDs (claude/codex/hermes/kimi/claude-code/gemini/opencode in various combinations:
Finance, auftra, cms, ecb-ai-prezz, ecb-meeting-2026.05, ldx, ldx-wt-mail-fixups,
rev-kimi-scratch, servers). Not quorum changes — but those IDs' `workspace_dir`/`role`
render-only values stop being fleet-visible, and the stale §2 rows resurface as
`section2-only` display noise under inheritance.

Value overrides keep flowing: on decks with pinned values (e.g. BYTE), ctrl and pathc
produce byte-identical MODEL/EFFORT/SPEED cells for all pre-existing members; the only
row-level diff is the added `zcode-1`.

For the record against my own signoff: it predicted "ungated, all 37 flip." Measured: 36
of 37 block-carrying decks. The 37th already declares the machine's full set.

### (d) The rule I would ship

**A versioned marker on the deck file.** `schema = 2` in `parley-deck/agents.toml`
(codex-1's `[membership] mode` stanza is the same mechanism with different spelling; I
don't care which spelling, I care about the three properties it buys). Under it: blocks
are values-only; membership is the `members` property (explicit list replaces, per the
ruling's own `class DeckConfig` example); absence of `members` inherits the machine
roster. Unmarked files: all three current branches — deck-declared, legacy §2, pure
inheritance — byte-for-byte.

Why it beats the alternatives:

- **Beats content-keying** (any rule reading block contents to infer intent): the census
  above IS the content-keyed fleet outcome — 36 silent quorum changes, which §1.3 and
  §1.5 forbid and which codex-1's signoff ruled out ("a content rule that flips any
  unmodified file is the §1.3 sin"). Content-keying also breaks under `sync`'s
  field-stripping: a file's membership meaning would change because a VALUE was removed
  (D-C's preserved empty headers are the live example).
- **Beats ungated C** (measured, not argued: 36 decks gain a member silently).
- **Beats keeping the authority rule**: the ruling is binding, and claude-1's measurement
  stands — membership is the one property whose resolution rule no other property uses;
  E1a shows the unified model needs no new syntax to express the owner's case.

`members` replaces; `super.members + [x]` arithmetic stays OUT of v1 — the measured fleet
demand is exactly zero delta instances across five censuses, and plain replacement covers
the ruling's own example. That is codex-1's falsifier intact, arrived at from the owner's
model.

## E2 result

**All 32 `[agents.*]` properties layer per field, machine → deck. PRIMARY, measured** with
a scratch-layer test (`TestE2EveryAgentsFieldLayers`, ctrl copy, shipped code): machine
sets every property of a family; (i) a deck silent about that family inherits every field
with per-field source `~/.parley/agents.toml`; (ii) a deck overriding every field takes
every deck value with source `parley-deck/agents.toml`. That extends claude-1's four
(sandbox, approval, model, timeout) to the full set: command(s), version_args, launch/
headless/interactive modes and args, prompt_mode, model_label, reasoning, profile, speed,
the supervision knobs, isolate_home, buffers_stdout, isolated_home_env,
external_backend, telemetry, notes. `[defaults]` (incl. `[defaults.loop]`, presence-aware
per F-T2-1) layers per field the same way.

**What does NOT layer, measured or read:**

1. **Roster membership** — the subject of this idea; block presence, not a property.
2. **Roster `active`** — deliberately authority-bound (`runtime.go:127-133`); under Path C
   it must become layered (my patch layers machine ⊕ deck), which re-opens the question
   the 2026 ratification closed: a machine-wide `active = false` now deactivates the
   member in every values-only deck. Intended by the model, unratified as a behavior
   change (Bin 1 above).
3. **Numeric `>0`-guarded fields cannot be overridden to zero** (PRIMARY,
   `TestE2ZeroOverrides`): `timeout_ms`, `interactive_timeout_ms`, `interactive_poll_ms`,
   and `[defaults.timeouts].*` treat 0 as "unset", so a deck cannot lower them to 0;
   pointer-typed siblings (`first_event_timeout_ms`, `stall_timeout_ms`, `heartbeat_ms`,
   loop ceilings, `isolate_home`, `buffers_stdout`) CAN be explicitly zeroed/falsed.
4. **Empty-slice inconsistency** (PRIMARY, `TestE2EmptySliceEdges`): `acp_args = []`
   clears the parent's list (non-nil replaces); `headless_args = []` is ignored (len
   guard). Two slice properties, two rules at the empty edge.
5. **Blanket-default shadowing** (PRIMARY, `TestE2DeckDefaultSpeedBeatsMachinePerAgentSpeed`):
   a higher layer's `[defaults] speed` overwrites a lower layer's per-agent `[agents.X]
   speed`. Defensible under C (the child spoke, later), but it means a parent's specific
   value does not survive a child's generic one — worth a line in the schema-2 docs.

So claude-1's claim generalizes: Path C is the implemented model for every `[agents.*]`
and `[defaults]` field, with quibbles at the zero/empty edges, and roster membership is
the one wholesale exception — `active` is a deliberate, documented second exception that
Path C converts into an ordinary layered property at the cost of re-opening a ratified
guard question.

## Position under Path C

I argued (a). The ruling rejects (a)'s premise — that membership-declares-values is a
ratified design worth preserving — and claude-1's measurement plus my E1a show the
replacement model is already 95% shipped. **I accept the direction and I am not
re-litigating it.** My round-2 file contained the seeds of this: I conceded the coupling
("no way to change one local setting without owning the whole membership list"), my
signoff said a positive §2.1 result "would change something I am signing" (the D-A fix
shape), and my amended trigger's owner-instruction disjunct has now fired in the
strongest possible form — the owner didn't just name a deck that needs value-overrides
under inheritance, he redrew the model. The demand question my (a) rested on is answered
by instruction of record; what remains of my position is sequencing, and it survives
intact: **D-A, D-B, D-C land first (unchanged, unanimous §1), then the schema-2 resolver
+ verbs as one contract change, then attended migration.** What I now build is C's end
state wearing codex-1's compatibility boundary — which codex-1's own signoff already
converged on: its `inherit-values-v1` mode stanza is my `schema = 2` with a different
key name.

Defects in C the owner has not seen — stated plainly, per the ruling's own invitation:

1. **The parent class is machine-local.** Inheritance binds a committed file's effective
   quorum to the `~/.parley/agents.toml` of whoever runs the command. With an explicit
   `members` list this is safe; the moment a deck drops the key to track the machine, its
   quorum differs per machine (codex-1's machine-dependence point, now structural).
   Mitigation already exists in spirit: rosters are snapshotted at run kickoff — that
   must stay mandatory, and `roster show` should print the machine source path it
   inherited from (it does, via `--explain`).
2. **The no-parent edge is a hole until fixed**: values-only blocks over a machine with
   no roster silently read legacy §2 today (Bin 2, measured). Must fail closed.
3. **`active` layering re-opens a ratified guard.** Machine-wide retirements will
   propagate into values-only decks. I believe that is correct inheritance and I signed
   the patch that way, but the 2026 ratification ("retiring a member belongs to the
   record that grants membership") should be re-ratified for the machine layer, not
   silently overridden by this idea.
4. **The verbs are state-blind.** render/set/sync/explain cannot distinguish
   pure-inheritance from values-only-inheritance (Bins 2-3, measured). D-C's truthfulness
   bar applies to all of them under C.

None of these is a reason to reject C; all four are reasons C ships as a versioned
contract change with the three verb fixes, not as a resolver patch.

## Migration

The 36 changing decks must not move silently, and under the marker rule they don't move
at all until each is attended:

1. **Fix D-C first** (codex-1's condition; my signoff's precondition). Today's `sync`
   still lies about inheritance; no migration instrument may be built on it.
2. **One new verb** (the `roster inherit`/`adopt-schema-2` shape from my round-2
   proposal, now writing the marker): per deck, git-first (18/26 tracked dirty, 15/41 in
   no git work tree — round-2 PRIMARY; the `--backup-dir` machinery from `roster migrate`
   covers the untracked), it prints the resolver's before/after member sets, then writes
   **`schema = 2` AND `members = [<the deck's current effective set>]` in the same edit**.
   The quorum is preserved byte-for-byte at conversion: a five-member deck converts to an
   explicit five-member list, `zcode-1` does NOT silently join, and "track the machine
   roster" becomes a separate, deliberate, previewed act (delete the `members` key). The
   adopt-machine vs preserve-set choice is made per deck by the operator, never inferred
   (§1.3). Post-write, the verb re-resolves and verifies the member set is unchanged.
3. **Retired-member blocks** (the 16 decks): conversion leaves them as value overrides or
   offers to delete the block in the same attended step — an explicit choice, since
   dropping an `antigravity-1 active=false` block changes rendered §2 history.
4. `millenniumProblems` and this repo's deck need nothing; they are the proof the marker
   rule is byte-preserving.

## What I would sign

- **Path C as the target model, gated by a versioned marker** (`schema = 2` or the
  `[membership]` stanza — spelling undecided, mechanism not): unmarked files
  byte-for-byte on all three current branches; `members` replaces; no add/remove
  arithmetic in v1; no content-keyed discrimination (the census above is what
  content-keying costs); no-parent fails closed; `active` layers machine ⊕ deck with the
  gitignored/env guard kept and the machine-propagation behavior explicitly re-ratified.
- **D-A, D-B, D-C as fast-track `standard` ideas, unchanged from §1**, now as hard
  prerequisites of the schema-2 contract change rather than of an overlay — D-A's gate
  must state resolver before/after sets against the NEW rule (my fix1b run shows the
  current gate text is false under C in the opposite direction).
- **The attended migration verb** above as its own idea, sequenced after D-C.
- **Runs record their effective member set at kickoff** (codex-1's point, structural
  under C — machine-local parent).
- **Will not sign:** ungated Path C (measured: 36 silent quorum changes); any rule that
  infers membership from block contents; FINAL while `opencode-1` owes an artifact unless
  the human gate has decided its §5 status (unchanged from my signoff).

My signoff's request — "run §2.1 before FINAL" — is discharged by this file; codex-1's
seven-point matrix is satisfied in substance by E1a (1), the census (3: unmarked decks
retained exactly — the marker rule, and the ctrl runs are that control), fix1b and Bin 2/3
(5, 6, 7 partially: D-A print and `--explain` truthfulness are specified, not
implemented). Point 2 of its matrix (machine churn reflected without deck rewrites) is
demonstrated by fix3, and point 4 (legacy §2) is preserved by the marker rule and fails
closed at the no-parent edge. I consider its block conditions addressable by this
round's evidence; whether they are ADDRESSED is codex-1's call, not mine.
