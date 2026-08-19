---
agent: kimi-1
idea: roster-membership-overlay
round: 1
date: 2026-08-19
---

# roster-membership-overlay — kimi-1, round 1

## Summary

**NO CHANGE to the membership authority model.** The owner's 99% case — inherit the machine
roster, declare nothing locally — already works today in this deck (verified below, PRIMARY), and
nobody, the facilitator included, has named a concrete deck that needs machine-plus-one or
machine-minus-one. What my verification actually turned up is that the live risk sits elsewhere:
the documented tool gesture for "change one thing locally" (`parley roster set --scope deck` on an
inheriting deck) silently collapses a six-member roster to one (PRIMARY, reproduced), and
`parley roster render` writes a §2 table shape this repo's own drift guard rejects (PRIMARY,
reproduced). Fix those; defer the overlay behind a named trigger, exactly as
`protocol-overlay-local-extension` deferred `replace`. If the deck votes to build anyway, the only
defensible v1 is extend-only (`+`) in TOML under new syntax, full `[roster.*]` lists unchanged,
rule 2 untouched.

## What I actually ran (verification log)

All commands run on 2026-08-19 against this repo (parley 1.45.0 at /opt/homebrew/bin/parley).
Temp experiments used copies under /tmp; the shared working tree was never modified.

- `parley roster show` — all six agents report `inherited-roster`; this deck declares no roster
  of its own. **PRIMARY.**
- `parley roster show --explain zcode-1` and `--explain kimi-1` — per-field layering confirmed:
  each field carries its own SET BY line (machine TOML, the agent's own config read at launch,
  built-in default). Matches prompt claim 1. **PRIMARY.**
- `parley roster show --scope machine` — machine roster of six, STATUS `ok`. **PRIMARY.**
- `sed -n '160,210p' internal/config/runtime.go` — the AUTHORITY ORDER comment is verbatim as the
  prompt quotes: (1) committed deck blocks, (2) else a VALID legacy §2 table, (3) else machine
  roster marked Inherited. Also present: "Step 2 was ratified and then omitted: LoadRosterScoped
  knew only about TOML, so any machine roster silently outranked a legacy deck's four declared
  members." **PRIMARY.** `internal/app/roster.go` (~line 293) adds: "MEMBERSHIP IS THE DECK FILE,
  not the layered union. Iterating `entries` here … is what made a deck declaring two members
  resolve to five." **PRIMARY.**
- `cat parley-deck/agents.toml` — no `[roster.*]` blocks; a dated comment records the owner
  instruction of 2026-08-19 ("do not override anything locally, use the global roster") and
  explains this is deliberately not `roster sync`. **PRIMARY.**
- Claim-3 reproduction in `/tmp/kimi-roster-exp` (copy of this deck): with the §2 roster table
  empty and no TOML blocks → `inherited-roster`. After inserting four rows into the §2 table of
  the temp copy → every row flips to `legacy-roster` and the roster shrinks to exactly those four
  (hermes-1 and opencode-1 drop out). After emptying the rows again → `inherited-roster` returns.
  The facilitator's incident reproduces exactly as described. **PRIMARY.**
- `go test -count=1 ./internal/protocol/...` — `ok … 0.255s` (uncached). **PRIMARY.**
- Drift-guard edits in `/tmp/kimi-drift-exp` (temp module copy: go.mod, go.sum,
  internal/protocol, internal/fsutil, parley-deck/COOPERATION.md):
  - prose line inserted before the §2 roster table header → `TestEmbeddedDefaultMatchesLiveDeck`
    FAILS ("differ after normalizing the allowlisted zones … first divergence at normalized
    line 132"). **PRIMARY.**
  - duplicated roster header → FAILS ("appears 2 times, want exactly 1 (drift guard fails
    closed)"). **PRIMARY.** (Missing header is the same `assertExactLineOnce` count≠1 path —
    same code, not separately executed.)
  - control: an extra row *inside* the table body → test PASSES (body is normalized away).
    **PRIMARY.**
- Closed STATUS vocabulary — read in the parley-deck skill's SKILL.md (~line 272): `ok`,
  `unmapped`, `not-installed`, `model-drift`, `model-unbound`, `effort-unknown`,
  `metadata-unknown`, `masked-by-env`, `legacy-roster`, `inactive`, `stale-snapshot`,
  `section2-only`, `inherited-roster`, `not-in-roster`. **PRIMARY** (for what the document says).
- `protocol-overlay-local-extension` precedent — FINAL.md B1: "Extend-only. One operation kind,
  `extend`, at most one instance. No `replace`. (User ruling.)" `protocol-overlay-replace` is
  explicitly deferred. **PRIMARY.**
- 1.41.0 history — CHANGELOG.md: "The headline defect: a deck's declared membership was not what
  ran… a deck listing two agents resolved to five whenever `~/.parley/agents.toml` listed five."
  COOPERATION.md §2 records the nine-rosters-across-40-decks measurement. **PRIMARY** that these
  documents say this; the underlying 40-deck measurement itself is **SECONDARY** (I did not
  re-measure 40 decks).
- Fleet size — `parley-deck/meta/protocol-changelog.md` (2026-08-07) records "measured across 36
  decks before a one-off sync — eight different `deckVersion` values, §15 present in 5 of 36, the
  §2 roster-authority change in 1 of 36." **PRIMARY** for the text. The claim "36 decks are synced
  to this protocol (fleet-protocol-sync, 2026-08-06)" with "several carrying uncommitted local
  changes" is **UNVERIFIED**: I found no machine-readable fleet registry in this repo
  (`cache/projects.json` lists exactly one project — this one), so the current fleet state cannot
  be confirmed from here.

### Two findings the prompt did not contain

- **F1 — `roster set` collapses an inherited roster.** In the temp copy of this (roster-less,
  inheriting) deck: `parley roster set ghost-1 --scope deck --adapter kimi --yes
  --confirm-breaking` wrote a single `[roster.ghost-1]` block; `parley roster show` then listed
  exactly one member — the five inherited members vanished from the quorum. The tool's own
  confirmation text framed it as "adds a new roster member." The preview-by-default behavior is
  good; the semantic is a trap, because under winner-takes-all the first local write replaces the
  whole roster. **PRIMARY** (command output observed; consistent with the authority-order code
  quoted above).
- **F2 — `roster render` disagrees with the drift guard.** `internal/app/roster_render.go:73`
  renders `| Agent ID | Workspace dir | Role | State |` (compact, 4 columns); the drift guard
  anchors on `| Agent ID       | Workspace dir                       | Role          |` (padded, 3
  columns, drift_test.go:28). Reproduced end-to-end: `parley roster render --yes` in the temp
  deck rewrote §2 to the compact shape; copying that file into the temp module and running the
  guard failed closed ("anchor … appears 0 times, want exactly 1"). Consequence: the first
  post-migration (or `--adopt-inherited`) `roster render` on this repo's own deck breaks its own
  test suite. **PRIMARY.**

Independence disclosure: a repo-wide grep for `fleet-protocol-sync` returned two single lines from
peers' round-01 files (hermes-1, zcode-1) inside its result snippets. I did not open or read any
peer artifact; those two lines played no role in forming the above, which rests on the commands
listed. Reporting this under the round-1 independence discipline and §15 honesty norms.

## Proposed approach

Do not build the membership overlay now. Instead, in order of urgency:

1. **Fix F1 (the set-collapse).** When a deck's roster is currently inherited, `roster set
   --scope deck` must not silently materialize a one-member override. Options, cheapest first:
   refuse with "deck currently inherits; use `--adopt-inherited`-style materialization first", or
   materialize the full inherited roster and show the complete before/after membership diff for
   confirmation. This is the operation the owner's sentence actually maps to — "if something is
   changed locally, apply it" — and today it is the most dangerous roster gesture in the CLI.
2. **Fix F2 (render vs. guard).** Pick one §2 header shape and make both the renderer and the
   drift-guard anchor use it. Any work touching §2 (including any future overlay) is blocked on
   this anyway.
3. **Close the incident's actual trap.** The facilitator's failure was not the absence of an
   overlay; it was that a *tool-generated* §2 view is indistinguishable from a *hand-written*
   legacy table, so deleting TOML blocks silently promoted the view to rule-2 authority (my
   temp-deck reproduction confirms the mechanism). A targeted fix: `render` emits a marker comment
   above the generated table, and the rule-2 reader declines marker-tagged tables (or at minimum
   `roster show`/preflight warns loudly when membership resolves to `legacy-roster` while a
   machine roster exists). Small, local, testable — unlike an authority-model change.
4. **Defer the overlay with a named trigger**, mirroring the deferred `protocol-overlay-replace`:
   reopen when a concrete deck files a need for machine-plus-one or machine-minus-one that full
   enumeration plus `roster sync` cannot service. That converts an evidence-free build into an
   evidence-gated one.

If the deck nonetheless ratifies the overlay, the only shape I could sign: extend-only `+` deltas
in TOML under **new syntax** (an unmarked full `[roster.*]` list keeps today's full-override
meaning forever); rule 2 untouched; additive STATUS terms only; `--explain` gains a
membership-provenance line. Details in the answers below.

## Answers to the six questions

**1. Operations — extend-only, and the sibling precedent transfers with extra force.**
`protocol-overlay-local-extension` ratified extend-only for v1 and deferred `replace` (PRIMARY,
FINAL.md B1). For membership the asymmetry is stronger than for protocol text: a protocol
`replace` rewrites prose a human reads; a membership `-` rewrites *who deliberates*, and the
consequence is computed, not read. A removal delta is sticky and invisible — `-zcode-1` survives
every machine-side change and silently narrows every future quorum — while an extend delta only
ever widens. Removal is already expressible today as a full enumerated deck roster, which has the
one property a removal must have: the resulting quorum is fully visible in one committed file.
Note also that the base floats under any inheritance scheme — a new machine agent auto-joins the
deck quorum with no deck-level decision. That is already true of plain inheritance (ratified de
facto by the owner instruction recorded in agents.toml); a delta does not create that property,
but `-` would weaponize it into "auto-join everyone except the excluded." So: `+` only, if
anything; `-` stays expressible solely through full enumeration.

**2. The 36 existing decks — unmodified files must never change meaning.**
A full `[roster.*]` list stays a full override, permanently, not within a migration window. The
prompt is right that reinterpreting an unmodified file is the dangerous case, and the changelog
evidence makes it worse: at measurement time the §2 roster-authority change was present in only 1
of 36 decks (PRIMARY for the changelog text), so a semantic re-interpretation would land hardest
exactly on decks whose protocol copies predate the cutover — decks least positioned to notice.
Therefore deltas need new syntax or an explicit opt-in marker in the deck file; absence of the
marker means current semantics. That preserves back-compat by construction and gives
`roster show` a cheap discriminator. (Caveat: the current fleet state is UNVERIFIED from this
repo — see the verification log — so any migration plan must measure first.)

**3. Rule 2, the legacy §2 table — keep it an authority; never make it a delta; do not stop
reading it.**
My temp-deck run confirms a §2 row is a full membership declaration (PRIMARY). Existing legacy
tables are full rosters by ratified semantics; reinterpreting them as deltas against the machine
roster would be the same unmodified-file sin as in answer 2, committed against prose whose authors
may not know it is load-bearing. Stopping rule 2 would silently re-quorum every pre-cutover deck —
at measurement time that was potentially 35 of 36. The drift-guard claims in the prompt are
accurate (PRIMARY: prose before the header fails via line-count divergence; a duplicated header
fails the exactly-once anchor assertion; an in-table row is normalized away). But F2 shows the
guard and the renderer already disagree on the header shape today — that latent conflict must be
repaired before §2 is touched by anything, overlay or not.

**4. Visibility / STATUS — the mechanism exists; extend it additively or not at all.**
`roster show` already answers "who is in this deck's quorum" in one frozen table, and the STATUS
vocabulary is closed and documented in the skill's SKILL.md (PRIMARY — the 14 terms listed in my
verification log). The deck scope already distinguishes `inherited-roster` from `legacy-roster`,
so provenance display for layered membership is a solved pattern. If a delta ever ships: add new
STATUS terms (allowed — adding to the contract is permitted), never change what
`inherited-roster` means; the table prints the computed set only, never delta syntax; and
`--explain` gains one membership-provenance line ("membership: machine base + deck delta
(+hermes-2)"). Anything less re-creates the three-tables-one-question failure the SKILL.md
explicitly warns about.

**5. The anti-goal — the disease was unvalidated hand maintenance, not local difference.**
Nine rosters across 40 decks happened because prose tables were hand-edited and no tool computed
an effective set; drift was invisible by construction (SECONDARY — the measurement as recorded in
§2 and the CHANGELOG). A TOML delta with a tool-computed effective set does not recreate that
*mechanism* — validation stays centralized. What it does recreate in miniature is unbounded
per-deck divergence, so the honest answer includes the metric: (a) count of decks with any local
membership declaration (target ≈ 0, the owner's 99%); (b) count of decks whose effective roster
differs from the machine roster; (c) dangling deltas (`+x` where the machine has since retired
`x`). If (a) grows toward "most decks carry deltas," the overlay became the new prose table.
Candor requires noting that this measurement capability does not exist today — I found no fleet
registry in this repo (PRIMARY: `cache/projects.json` holds one project) — so the anti-drift
instrumentation is itself unbuilt prerequisite work, which belongs on the cost side of the
overlay's ledger.

**6. Do nothing — yes, and it is the recommendation.**
The 99% case is not a proposal; it is the observed state of this deck: six rows, all
`inherited-roster`, zero local membership declarations, owner instruction dated and recorded
(PRIMARY). The concrete situations an overlay would service — machine-plus-one for a
project-specific reviewer, machine-minus-one for a constrained deck — are both expressible today
via full enumeration, and `roster sync` exists precisely to pay the staleness tax on those copies
(rebase semantics, machine → deck, deliberate pins enumerated, per the skill's SKILL.md — PRIMARY
that the document specifies this). Neither situation has been named as an actual need by any deck;
the prompt asks us to produce one and I cannot produce one from evidence. Under §15.4, a
recommendation claiming to avoid a named obstacle needs a witness; the overlay's obstacle — "full
enumeration is unserviceable" — currently has none. What the owner's sentence does collide with
today is F1: the obvious tool gesture for "change one thing locally" replaces the whole roster.
Fixing F1 serves the owner's stated intent more directly than building the overlay the sentence is
being read as requesting.

## Concerns / open questions

- Was F1 (`roster set` on an inherited deck → one-member override) ever a deliberate design
  decision, or is it simply the unexamined composition of "set writes one block" with
  winner-takes-all? The changelog and code comments I read do not say. If deliberate, the
  reasoning needs to be on record before this idea closes, because it bears directly on the
  owner's request.
- F2 means this repo's deck cannot survive a future `roster render` without a test failure. Is
  the compact 4-column header the intended future shape (guard anchor stale), or is the padded
  3-column shape canonical (renderer stale)? Someone who ratified `embedded-default-protocol-resync`
  should answer; I only demonstrate the conflict.
- The "36 decks, several with uncommitted local changes" exposure claim is decision-relevant and
  currently UNVERIFIED from anything I can reach. Before any authority-model change ships, the
  fleet measurement the changelog describes needs to be re-run and committed as an artifact this
  idea can cite as PRIMARY.
- If consensus converges on NO CHANGE, the deferral trigger must be written down in FINAL.md as
  precisely as `protocol-overlay-replace`'s was, or "deferred" will decay into "forgotten."
- Per §15.6 and the facilitator's own declared bias: if most participants reach NO CHANGE, that
  is *disagreement with the facilitator's framing* and carries independent evidentiary weight;
  if most reach YES, consensus.md must record the shared prior explicitly. The drafter should not
  treat either unanimity as convergent evidence.

## Risks

- **Of building the overlay:** a second membership semantics (delta alongside full override) on
  top of a base whose existing edges (F1, F2) are still sharp; sticky invisible removals if `-`
  is ever admitted; quorum drift in decks nobody is watching, with no fleet instrumentation to
  detect it; and a semantic re-interpretation risk for 36 decks whose protocol state is
  historically heterogeneous and currently unmeasured from here.
- **Of NO CHANGE (stated honestly):** a deck that genuinely needs machine-plus-one pays the
  full-enrollment copy tax and its staleness risk, mitigated but not eliminated by `roster sync`;
  and if F1/F2 are *not* fixed, the trap that triggered this idea remains armed for the next
  facilitator — NO CHANGE to the authority model must not become NO ACTION on the defects the
  investigation exposed.
- **Of the small fixes I propose:** F1's fix changes a shipped CLI behavior (needs its own idea
  and track classification; likely `standard`); the §2 marker for generated tables only helps
  forward — existing hand-written legacy tables are indistinguishable from old generated ones,
  so rule 2 must keep reading unmarked tables, which is precisely why the marker can only narrow,
  never abolish, the legacy path.
