---
agent: zcode-1
idea: roster-membership-overlay
round: 1
date: 2026-08-19
---

# Round 1 — zcode-1

## Summary

**NO CHANGE to the membership authority model.** The owner's 99% case — "nothing local,
global roster rules" — is already served by the inherit state this deck runs today
(PRIMARY: `parley roster show` reports all six members `inherited-roster`; this deck's
`agents.toml` contains no `[roster.*]` blocks). The delta case the overlay would serve has
**zero demonstrated demand**: in a census I ran across the 41 decks on this machine's shared
volume, not one deck expresses a membership delta — 37 carry full local declarations that
spot-checking shows are stale *copies* of old machine rosters, and 4 inherit (PRIMARY, method
below). The measured disease today is stale copies, not missing overlay expressiveness, and
the cure for that is migration onto the existing inherit state, not a new authority model.
I propose: migrate the fleet to inheritance behind a safe one-command transition (the manual
path trips rule 2 — I reproduced the trap), and record a measurable trigger that re-opens a
narrow extend-only overlay later if real demand appears.

## What I ran (verification log)

All commands run 2026-08-19 from the repo root unless noted. Temp-dir tests were performed
on copies; the shared working tree was not modified.

| # | What | Result | Tag |
|---|------|--------|-----|
| 1 | `parley roster show` | 6 agents, every row `inherited-roster` | PRIMARY |
| 2 | `parley roster show --scope machine` | same 6 agents, STATUS `ok` | PRIMARY |
| 3 | `parley roster show --explain zcode-1` | matches the prompt's quoted block verbatim: membership from `~/.parley/agents.toml` (INHERITED), per-field SET BY layers | PRIMARY |
| 4 | `sed -n '60,260p' internal/config/runtime.go` | AUTHORITY ORDER comment exactly as the prompt quotes it; plus two facts the prompt doesn't quote (see below) | PRIMARY |
| 5 | `go test ./internal/protocol/...` | `ok` (unmodified tree) | PRIMARY |
| 6 | `cat parley-deck/agents.toml` | no `[roster.*]` blocks; comment records the owner instruction of 2026-08-19 ("lokalne nepretazuj nic, pouzivaj globalny roster") and why inheritance ≠ `roster sync` | PRIMARY |
| 7 | temp-dir repro of rule 2: copied the deck to `/tmp`, inserted one row into the §2 generated view, `parley roster show --dir` | membership = exactly that row, STATUS `legacy-roster,unmapped`; **all six machine agents vanish** | PRIMARY |
| 8 | drift guard, on a `/tmp` copy of the module: duplicated the roster header / inserted prose before it / removed it, `go test ./internal/protocol/...` each time | FAIL / FAIL / FAIL — all three fail closed, matching the prompt's claim | PRIMARY |
| 9 | fleet census: `find "/Volumes/My Shared Files" -maxdepth 4 -type d -name parley-deck` | 41 decks | PRIMARY |
| 10 | census classification (grep `[roster.*` blocks; parsed §2 table bodies) | 37 full-declared (5–10 blocks each), 0 legacy-table-only, 4 inherit/empty | PRIMARY |
| 11 | git state of each deck's parent dir | 26 inside a git work tree: **18 dirty, 8 clean**; **15 outside any git work tree** | PRIMARY |
| 12 | spot-check of the first 12 declaring decks' roster IDs | all carry the same five IDs (`claude-1,codex-1,hermes-1,kimi-1,opencode-1`); 4 also `antigravity-1`; 2 also unprefixed duplicates (`claude-code,codex,hermes,kimi`); **none contain `zcode-1`** | PRIMARY |
| 13 | `git log --grep` | `3abbd16 release: v1.41.0 — roster membership authority` exists; `49afe45` (stop declaring locally) commit message and diff corroborate the prompt's failure story in the facilitator's own words | PRIMARY (that the commit says so; the events described are claude-1's testimony) |
| 14 | read `ideas/protocol-overlay-local-extension/FINAL.md` | B1 extend-only ratified with a user ruling; `protocol-overlay-replace` deferred with named prerequisites (block IDs, extents, target hashes, tombstones, registry) | PRIMARY |

Claims I could NOT verify: the historical "nine rosters across 40 decks / 17 without roster /
17 naming a retired agent" measurement — recorded in COOPERATION.md §2 prose and the prompt;
not reproducible from here (SECONDARY, named source: COOPERATION.md §2). The "36 decks synced
2026-08-06 via fleet-protocol-sync" — no fleet registry exists anywhere I looked (in-repo grep
finds the string only in this idea's prompt; no `parley fleet` command; nothing in `~/.parley/`);
my census found 41 deck directories on this volume, which makes ~36 plausible but the synced
subset is unverifiable (SECONDARY for the event, PRIMARY for the absence of a registry).

Two things the code says that the prompt under-plays (both PRIMARY, #4):

- **The union overlay was already tried and was the disease.** The `RosterScope` comment:
  pre-1.41.0 `LoadRoster` "unions membership across every layer, which meant a deck declaring
  two members resolved to five whenever the machine file listed five — and, because `roster
  render` writes §2 from the same view, that inherited membership got committed into
  COOPERATION.md, where it went stale on the next machine change. That is the exact drift
  this change exists to end." The overlay question is not greenfield; the naive version ran
  in production and produced the nine-roster drift.
- **Membership is deliberately confined to the committed record — in both directions.**
  `configLayers` marks only `parley-deck/agents.toml` as a membership layer; the machine
  file seeds values, and the gitignored `agents.local.toml` / env layers are barred from
  membership entirely. `active` was also re-anchored to the authority layer ("Retiring a
  member is a membership change; it belongs to the same record that grants membership") —
  so even `active` does not layer today. The system's authors already answered "should
  membership layer?" once, deliberately, for both add and remove.

## Proposed approach

Build nothing in the authority model. Do three smaller things:

1. **Migrate the fleet onto the inherit state.** The census shows the owner's complaint is
   real but mislocated: 37 decks resolve quorums from frozen copies that predate the
   machine roster's current shape (spot-checked copies lack `zcode-1`, the newest machine
   member, and four still list `antigravity-1`, which this deck's `agents.toml` records as
   inactive since 2026-07-18). Those decks will not see machine changes until their local
   declarations are removed. That is migration work with existing tools, per-deck, attended.
2. **Make the declaration→inherit transition one safe command** (small, fast-track idea;
   tooling, not authority semantics). Today the transition is manual, two-file, and order-
   sensitive: deleting the `[roster.*]` blocks alone falls to rule 2 because the §2
   generated view still holds rows — reproduced empirically (#7), and it is exactly the
   trap the facilitator hit (`49afe45` had to empty both). A `parley roster inherit`-style
   verb that clears the blocks and the §2 body in one reviewed write, refusing when the
   deck would resolve to no roster at all, removes the footgun for all 37 migrations.
   Sequence the fleet cleanup first: 18 of 26 git-tracked decks are dirty and 15 decks are
   in no git work tree at all (#11) — membership changes landing there are unauditable.
3. **Record the re-open trigger now, and the design shape if it fires.** Trigger: after
   migration, decks that demonstrably need "machine roster ± a delta" — say, two or three
   real cases or an explicit owner instruction — open the extend-only overlay idea. Shape,
   so the deliberation isn't lost: a *new, explicit* grammar (e.g. a `[roster-overlay]`
   block with `add`), never a reinterpretation of `[roster.*]` blocks; add-only in v1;
   fail-closed if both override blocks and an overlay block are present; an additive STATUS
   term (e.g. `overlay-roster`) with `roster show --explain` naming the origin per member;
   and the render/rule-2 interaction addressed (a rendered view of a computed roster must
   not become the authority when the overlay block is later deleted — same trap as #7).

## Answers to the six questions

### 1. Operations (extend-only vs +/−)

Not built now; if ever, **extend-only in v1**. The sibling precedent's *mechanical* reason
does not transfer: `protocol-overlay-replace` was deferred because replacing protocol prose
needs addressing machinery (block IDs, extents, target hashes, tombstones — PRIMARY, FINAL.md
B4/Deferred), whereas membership deltas address discrete agent IDs and need none of it. What
does transfer is the *risk* logic: removal is the one operation that silently shrinks a
quorum, and it already has a safe path — a full declaration, deliberate and reviewable.
A `-` delta's only added value over full declaration is tracking machine churn, which for
removals means exclusions persist silently while new machine members join — audit-hostile.
The code's own ratified principle points the same way: membership changes, including
retirement, belong in the committed record that grants membership (PRIMARY, #4).

### 2. The existing decks (what a full `[roster.*]` list means)

Measured, not estimated: 37 of 41 decks declare full rosters (#10). Therefore a full list
**must remain a full override**. Reinterpreting it as a delta would silently change the
quorum resolution of 90% of the fleet in one release — each deck would suddenly resolve to
machine ∪ local instead of local. That is Q2's "dangerous case" with a measured blast
radius, and it is the single strongest argument against any implicit overlay. Back-compat
is achievable only via new opt-in syntax (approach §3), and even then the file carries two
grammars whose mixing must fail closed.

### 3. Rule 2, the legacy §2 table

It works as documented and I reproduced it: one row in the §2 view with no TOML blocks makes
that row the deck's entire membership, reported `legacy-roster`, machine agents gone (#7).
The drift-guard anchors are exactly as claimed — duplicated header, prose before the header,
and missing header each fail the suite (#8). What the census adds: **zero decks on this
volume rely on rule 2** — every deck is either TOML-declared (37) or inheriting (4). Rule 2
guards a population of ~zero here, though this volume may not be the whole fleet (SMB share;
caveat below). My position: leave rule 2 alone in this idea — neither authority-promoted,
delta-ified, nor deleted. It is back-compat for pre-cutover decks elsewhere, and retiring it
is a breaking change deserving its own idea once the fleet state is known. What should change
is the *ergonomics around it* (approach §2): the generated view doubling as the legacy
authority is what turns a routine cleanup into a silent authority flip.

### 4. Visibility / STATUS

With NO CHANGE, nothing to do — `parley roster show` already answers "who is in this deck's
quorum" in one table, and `--explain` already states the membership source per agent
(PRIMARY, #1–#3). If an overlay ever ships: membership becomes a computed set, so STATUS
must say so — but only *additively* (e.g. `overlay-roster` for added members, base rows
staying `inherited-roster`). The vocabulary is closed and documented in the skill's
`SKILL.md`; a new term is a contract change that must be made there explicitly, and changing
what an existing column means is out of bounds — the frozen-columns constraint stands.

### 5. The anti-goal (nine rosters across 40 decks)

The disease had two entry doors, both visible in the code (PRIMARY, #4): hand-edited,
unvalidated §2 tables, and — decisively — `roster render` committing *resolved inherited
membership* into §2, where it went stale on the next machine change. NO CHANGE keeps both
doors closed; nothing new is built that could reopen them. My census says the live form of
the disease today is stale full-copies (#10, #12), which inheritance cures. What I would
measure, now and after any future overlay: (a) the count of decks whose committed membership
differs from the machine base without a recorded reason — today that count is 37 and should
drop toward 0 under migration; (b) after migration, how many decks re-introduce local
declarations and why (each should carry a dated rationale in its `agents.toml`, as this
deck's did for `opencode-1`); (c) if an overlay ships, that `roster render` output is never
parseable as a rule-2 authority for a deck whose overlay was removed. Note no fleet-audit
command exists today (PRIMARY — no `parley fleet` in `--help`); my census was shell one-liners.

### 6. Do nothing

This deck already lives the owner's 99% case: zero local membership declarations, machine
roster inherited at read time, `zcode-1`'s row marked INHERITED (PRIMARY, #1/#3/#6). The
concrete situation that is *unserviceable* without an overlay is exactly one shape: "this
deck wants the machine roster **plus or minus one member** without freezing a full local
copy." Across 41 decks I found zero instances of that shape being attempted — no deltas, no
half-declarations; only full copies and pure inheritance (#10, #12). So the honest statement
is: the residue is real but currently has no customers, while the problem that does have
customers (stale copies, 37 decks) is fully serviceable with what already shipped. Recommend
NO CHANGE, with approach §3's trigger as the path back if that ever flips.

## Concerns / open questions

- **Census scope.** I counted one volume (`/Volumes/My Shared Files`, maxdepth 4). The share
  is SMB and may be multi-machine; decks elsewhere (other hosts, deeper paths) are not
  counted, and the "36 synced on 2026-08-06" subset remains unverified — there is no fleet
  registry to check against. Building a fleet registry would itself be a useful small idea.
- **Worktree clones inflate the census.** Several of the 41 are `wt-*` worktree copies of
  this very repo (e.g. `wt-roster-presets`), so some "stale decks" are this repo's own past
  states. Independent decks (`aditoLeads`, `IGBCE`, `lustrator`, `rev-kimi-scratch`,
  `design-mail-fe`) show the same stale-copy shape, so the conclusion survives, but the true
  fleet denominator is smaller than 41.
- **Deliberate vs copied.** I cannot distinguish a deliberate per-deck override from a
  `roster sync` copy by block counts alone; the 12-deck spot check (identical five-ID core,
  no `zcode-1`) argues copies, but the migration in approach §1 should confirm per deck
  before deleting a declaration — a deck that genuinely wants a different membership keeps
  its full override, and that remains a supported, if rare, state.
- **Two decks carry unprefixed duplicate IDs** (`claude-code`, `codex`, …) alongside the
  numbered ones (#12). That looks like adapter-mapping residue, is orthogonal to this idea,
  and deserves its own small audit during migration.
- **Is `active`-not-layering part of the owner's ask?** The owner said "apply local on top
  if something is changed locally" — values do layer, `active` deliberately does not (#4).
  I read the current design as correct (a gitignored layer must not be able to drop a
  quorum member), but if the owner expects per-deck *activation* deltas, that is a separate
  question this idea should answer explicitly rather than inherit silently.

## Risks

- **Of NO CHANGE:** the ±1 residue stays unserviceable, and a future deck needing it will
  be tempted to re-freeze a full copy — recreating staleness locally. Mitigated by the
  recorded trigger and design sketch (approach §3); the cost of getting it wrong is one
  stale deck, visible in the next census.
- **Of the migration I propose instead:** it edits committed files in 37 decks, 18 of whose
  repos are dirty and 15 of which have no git at all — bulk-editing those is how orphaned
  unauditable states happen. Mitigation: per-deck, attended, git-first ordering, and the
  one-command transition verb so no deck is hand-edited into the rule-2 trap.
- **Of building the overlay now (why I reject it):** the silent-semantic-change hazard on a
  measured 37-deck population; two grammars in one file; the render/rule-2 interaction; an
  additive-but-real STATUS contract change; and the precedent of v1-minimalism paying off
  (the registry was unanimously wanted and then unanimously dropped once scoped —
  protocol-overlay-local-extension, "round 1's unanimity was worth very little", PRIMARY #14).
  Against that stack of costs stands a benefit whose current demand is measured at zero.
- **Correlated agreement:** the facilitator declared a bias toward the overlay and asked for
  an independent NO CHANGE check; I reached NO CHANGE via the census rather than from the
  prompt's framing, but if the round closes unanimous, §15.6 applies — consensus must record
  the shared prior and steelman the overlay (this file's approach §3 is my attempt at that
  steelman, with the observation that would change my recommendation: post-migration decks
  demonstrably needing ±1).
