---
agent: kimi-1
idea: protocol-overlay-local-extension
round: 1
date: 2026-08-07
---

## Summary

Ship the overlay with the block-ID registry in this slice — D1's ratified text already places the
registry inside the release (consensus.md:83-84), zero published releases make the layout change
free (IMPLEMENTATION.md:475-479), and without it no compliant addressing exists at all. Syntax:
YAML frontmatter (operations + provenance) over Markdown payload sections — the convention every
other Parley artifact already uses, parsed with the yaml.v3 dependency the repo already has. The §2
roster table is **not** overlay content: rows render from `agents.toml`, and the measured annotation
class (dated roster directives, decision logs) gets one new typed renderer input rendered after the
roster table. `ext-1` renders last in registry block order; the loss report learns "moved, carried
by overlay operation X" as a distinct category from "removed". My rank-3 signoff condition from the
meta idea is currently **unmet**; this slice is its discharge.

Availability note: the 2026-08-07 preflight found kimi-1 `no-pong` (00-prompt.md:102-104). I am
present and participating; this file is my round-1 artifact, filed asynchronously per §5.

## Proposed approach

### 0. Provenance and evidence base

`PRIMARY` for every claim about this repository's code or ratified record — I read or re-read each
cited file this session, and each carries a `file:line` locator. The 29-deck empirical taxonomy is
the scoping brief's measurement, produced by internal helper readers who are **not participants**
(00-prompt.md:52-54); the decks live outside this repository and my instructions confine me to it,
so where a load-bearing claim rests on fleet numbers I cite the brief by line and mark it
`[BRIEF-TESTIMONY — UNVERIFIED by me]` per §15.1's transcription rule. I own no verdict on those
numbers; my design deliberately does not depend on their exactness — only on their shape (divergence
concentrates in the header and §2), which the in-repo code independently corroborates (the five
identity slots at render.go:9-20 are exactly those two zones).

### 1. D-a — overlay file syntax: YAML frontmatter + Markdown payload sections

The record genuinely never decided this (brief:151-153 confirms all three candidate syntaxes are
dead), so this is a proposal, not a ratification recall. I choose **YAML-frontmatter-plus-Markdown-
payloads**, one file at `parley-deck/protocol-overlay.md` (D4's path, consensus.md:109-110 —
settling R1.1's discrepancy in writing: my own meta-idea round-01 wrote `protocol.overlay.md`
(round-01/kimi-1.md:20 of that idea); D4 governs and I conform).

Grammar:

```
overlay     := "---\n" frontmatter "---\n" body
frontmatter := YAML mapping, unknown keys fail closed:
  overlay: 1                          # format version; required; must be 1
  core-version-range: "<range>"       # R1.4 (consensus.md:205-206, DPC-4); required
  operations: [<op>, ...]             # v1: 1..2 entries (at most one replace — see D-w)
op          := mapping, unknown keys fail closed:
  op: replace | extend                # required
  target: s6.6                        # replace only; must exist, be registry-replaceable
  at: ext-1                           # extend only; v1's only extension point
  id: "<deck-slug>.<slug>"            # extend only; shape per §6 below (D-d)
  base-sha256: "<64 hex>"             # replace only; FORBIDDEN on extend
                                      # (rev-3 D4, consensus.md:111-114 — an extend
                                      #  replaces nothing; codex-1's block stands)
  depends-on: [{block: <id>, sha256: "<64 hex>"}, ...]
                                      # extend only; ABSENT means all sealed blocks
                                      # (D10 default, consensus.md:200-203)
  rationale: "<non-empty string>"     # required on both (D4)
body        := one payload section per frontmatter operation, 1:1:
  "## replace: <target>\n" <markdown> | "## extend: <id>\n" <markdown>
```

Validation, all fail-closed (G4 posture, consensus.md:263-264): exactly one frontmatter block;
every frontmatter operation has exactly one matching body section and vice versa; a payload that is
empty or whitespace is rejected (R1.5's empty-overlay ban, consensus.md:115-117, extended to the
operation level so `operations: []` with frontmatter only is equally rejected); body headings
matching `^## (replace|extend): ` are machine-owned — a payload containing one is rejected, so the
splitter can never be confused by content. Payloads are **free-form Markdown, HTML comments
included**: D2's no-inline-markup rule binds the core body (consensus.md:93), not payloads, and the
four DF-4 decks carry their local content precisely as HTML-comment blocks and header lines
(IMPLEMENTATION.md:440-444). This answers D-p's coupled sub-question in the only way that lets DF-4
content ever become an `ext-1` payload.

Why not the other two:

- **Fenced directive blocks.** CommonMark fence info-strings are unstructured; D10's dependency
  set is a *list of id+hash pairs*, which does not fit an info string, so the design degenerates
  into a hand-rolled attribute grammar validated by regex — the H4 failure mode (render.go:380's
  format-locked, silently-failing stamp regex) reincarnated at the center of the new mechanism.
- **TOML sidecar.** Either it is a second file — violating R1.1's one-file rule
  (consensus.md:109-110) — or TOML embedded in the Markdown file, where Markdown payloads (full of
  `#`, `---`, quotes) fight string escaping for no gain. Markdown payloads belong in Markdown.

Why this one: every canonical Parley artifact an agent writes today — `00-prompt.md`, round files,
`consensus.md` — is YAML frontmatter over Markdown (Phase 0/1 templates, COOPERATION.md:275-294,
312-322), so the syntax is already the fleet's muscle memory; `yaml.v3` is already a direct
dependency (go.mod:10) already used to parse frontmatter (`internal/driver/checks.go:42`); and the
metadata/payload split matches the two audiences exactly — the checker reads only the frontmatter,
the reviewer reads only the diff.

### 2. D-b — the registry ships in this slice

**Yes, now.** Three independent arguments, each sufficient:

1. **It is already ratified as part of the release.** D1's first sentence: the release directory
   "holds the exact core Markdown plus its registry, both hashed" (consensus.md:83-84). Rank 1
   shipped without it — a release directory holds exactly one file today (`CoreFileName =
   "COOPERATION.md"`, core.go:26; `Publish(version, body string)` has nowhere to put one,
   core.go:137,167-175 — H7 confirmed). Publishing a first release without the registry would
   *contradict the ratified layout*, and releases are write-once: the omission could never be
   repaired in that release.
2. **The free window closes at first publish.** The core store is empty on this machine and no
   `protocol-lock.yaml` exists anywhere under the workspace (H14; IMPLEMENTATION.md:475-479).
   Changing the release layout today costs nothing; after the first publish it costs a migration
   per release, forever.
3. **There is no third addressing option.** The kickoff asks what it is if the registry slips.
   Honest answer: none exists. D2 forbids heading text and forbids inline markup
   (consensus.md:88-93). The only remaining conceivable option — byte extents recorded in the
   *deck's* overlay/lock instead of in the release — is just a registry distributed into per-deck
   files: it breaks D1's self-contained, doubly-hashed release, duplicates extent data across 29+
   decks, and puts the integrity-critical mapping in the less-trusted store. It is strictly worse
   and I will block it if proposed. And the blocker is concrete, not aesthetic: v1's single
   replaceable block `s6.6` is an ordered **list item** — `6. **English only.** …` at
   parley-deck/COOPERATION.md:743 under the `## 6.` heading at :732 (H8 confirmed) — so
   heading-delimited segmentation cannot express even the one block v1 needs.

**Sub-decision — how extents are delimited.** Registry file `registry.yaml` inside the release
directory, hashed alongside the body per D1. Each entry:

```
- id: s6.6
  kind: sealed | replaceable | extension-point | tombstone
  extent: {start: <byte>, end: <byte>}   # into the CRLF-normalized release body
  sha256: "<hex of body[start:end)>"     # integrity, computed at publish, verified at load
  structure: "ordered-list item 6 within s6"   # review documentation, never parsed
- id: ext-1
  kind: extension-point                  # no extent; position per §4 below
```

Byte extents — not structural paths, not core restructuring. Because releases are write-once and
content-addressed, extents computed by the publisher cannot drift within a release; the per-block
hash makes any corruption fail closed at load; and `Publish` grows to accept and validate the
registry rather than hand authors computing offsets. Structural descriptors stay as human review
text because a structural path language ("6th list item under s6") is a second addressing scheme to
specify and test — one scheme, machine-checked, is enough. Rejected: restructuring the core so every
addressable block is a heading — that rewrites ratified core *text* (a §7, user-ratified act) to
serve a mechanism that byte extents solve for free.

### 3. D-k — the §2 roster table is not overlay content

**Option (c): the overlay owns nothing in the table; rows always render from `agents.toml`; the
overlay's only §2-adjacent role is the annotation block, carried as a new typed renderer input.**

Rows first. `agents.toml` is the roster authority — the protocol text says so verbatim
(defaults/COOPERATION.md:102-103), and R7.4 forbids a competing surface. Option (a) (overlay owns
rows) creates exactly that competitor and double-books the drift class the brief reports: 23/29
tables already disagree with their own `agents.toml` `[BRIEF-TESTIMONY — UNVERIFIED by me:
brief:72]`. Option (a) is out.

But (c) has a dividend the brief does not fully cash in: **once rows render from config, H9's
roster-wipe vector dies by construction.** Today's wipe path is extraction: `ExtractIdentity` →
`tableRows(body, "Workspace")` (render.go:155-156) scans the *deck's* table by prose match
(`| Agent ID` + "Workspace", render.go:129-133,169-189); a deck whose header reads "Workdir"
extracts zero rows, and the render emits header + separator with **zero data rows**, reported only
as `## 2 — 2 lines not carried forward` (H9's probe; mechanism confirmed in the cited lines). When
the caller populates `IdentitySlots.RosterTable/HandleTable` from `config.LoadRosterScoped`
(roster_render.go:30 already does exactly this for `roster render`) instead of from parsing the
prior deck body, no deck's bespoke schema can empty its roster. Rows whose bespoke *columns* have no
config home yet appear in the removal report with a remediation pointer — loud, explained, reversible
— instead of a silent empty table. This is also what the ratified text already implies: D3's slot
list is *values*, and the roster slot's value's authority is config.

The bespoke columns `[BRIEF-TESTIMONY — UNVERIFIED by me: brief:15 — 11/29 decks; Model ×7,
CLI/runtime ×5, State ×4, …]` mostly duplicate data that already has a config home: model and
reasoning live in `agents.toml` via `roster set` (defaults/COOPERATION.md:56), and `State` is
already a rendered column (roster_render.go:73). The right home is `parley roster render` growing
those columns from config — **DF-8**, named below — not the overlay. Note the live inconsistency
this follow-up must reconcile atomically: the drift guard demands the padded 3-column header
`| Agent ID       | Workspace dir … | Role          |` exactly once (drift_test.go:28,59-61) while
`roster render` emits the 4-column `| Agent ID | Workspace dir | Role | State |`
(roster_render.go:73) — this repo's own deck carries the former (parley-deck/COOPERATION.md:133),
so `roster render --yes` here today *breaks* `TestEmbeddedDefaultMatchesLiveDeck`. H13's second half
is not hypothetical; it is the current state.

Annotations. The genuinely inexpressible residue — dated user directives, roster-decision logs,
MANUAL-Bash caveats; the auftra 10-line HTML-comment directive whose loss IMPLEMENTATION.md:465-469
measured — is deck-local *data about the roster*, not protocol override. It is the same class as the
other identity slots (R2.4: "data substitution, not OOP-style replacement"). So I propose a
**seventh identity slot**, `RosterAnnotations []string`, rendered verbatim immediately after the
roster table body — the position the fleet actually uses `[BRIEF-TESTIMONY — UNVERIFIED by me:
brief:17 — 3/29 decks append after the §2 table]`. A fixed render position, not a zone extraction:
I checked this repo's deck, and core prose follows the handle table ("When a new agent joins…",
parley-deck/COOPERATION.md:156-162), so any "everything after the last table" zone boundary would
capture core text. The renderer emits core §2 prose, then the roster table, then the deck's
annotation block, then the remaining core §2 prose and handle table. One caveat I flag for round 2:
this amends D3's ratified enumeration of six slots (consensus.md:99-100). I judge it an
implementation-level amendment within this idea's mandate — the slot list lives in the renderer and
the consensus record, not in sealed core text — but consensus.md for *this* idea must record the
amendment explicitly rather than smuggle it. If participants reject it, the fallback is annotations
as `ext-1` payload with a cross-reference, accepting that nothing in §2 can point at them (§2 prose
is core-owned) and that auftra's directive lands at end-of-file rather than beside the roster it
governs. I consider that fallback materially worse and say so now.

Consequences, stated not implied: the drift guard must then also normalize this slot
(drift_test.go:102-130 pattern); D-m's five-vs-six mismatch (struct at render.go:14-20 declares
five fields while its own comment at :9 says six) is resolved by rewriting the slot accounting
honestly in the same commit; and `**Parley deck:** ./parley-deck/` (parley-deck/COOPERATION.md:4)
gets its disposition: **not a slot** — the deck directory is a tool constant
(`protocol.DeckDir`, used at protocol.go:101), so a divergent value is not project data, it is a
self-contradicting deck; the core keeps the line and no mechanism preserves overrides of it.

### 4. D-c / D-g — `ext-1` placement, and what the loss report means now

**D-c: `ext-1` renders last in registry block order — end of body in v1 — recorded in the
registry, not in prose.** I am abandoning neither determinism nor my round-2 "after the final core
section"; I am making it precise. codex-1's "after §8, before the TL;DR/reference appendices"
(round-02/codex-1.md:46-48 of the meta idea, quoted in brief:161) rests on a layout premise that
does not hold in the actual core: in this deck, `## 10. TL;DR` sits at line 810, *before* `## 9.`
at 826, and `## Appendix A` at 1084 *precedes* §12–§15 (1110-1369). "Before the appendices" would
render deck-local rules mid-document, before §15 verification integrity. hermes-1's "end of file is
non-deterministic" (brief:163) is half right: EOF is deterministic per release but drifts as the
core grows. The registry dissolves all three positions: block order is a publish-time fact recorded
per release, `ext-1` is an entry in that order (last, v1), and a future core that wants deck content
elsewhere moves the entry in a new release — where D10's change report shows the move instead of
silently re-anchoring it. One declared position, per D3 (consensus.md:98), now machine-checked.

**D-g: option (ii), scoped to declared overlay operations — and I reject (i) and (iii)
explicitly.** The H2 mechanism is confirmed in source: `droppedContent` is an order-sensitive LCS
by deliberate design (render.go:193-202,217-272), so a deck section repositioned to `ext-1` reports
every moved line as lost while every line survives. My rule: when the lines a deck loses are
*exactly* the lines a declared overlay operation carries, the report classifies them as **moved** —
`carried by overlay operation <id> (was: <heading>)` — never as removed. Unscoped (ii) (suppress any
loss whose lines reappear anywhere) is rejected: it blinds the report to real meaning-changes the
LCS was built to catch (render.go:198) and lets duplication mask deletion. (i) is rejected outright:
exempting overlay-carried content from reporting re-creates the silent-erasure class nine review
cycles bought (IMPLEMENTATION.md:409-416) — the report must *praise* the carry, not go quiet. (iii)
is rejected because it trains operators to ignore the report on exactly the decks the overlay exists
to rescue: H17 already guarantees a noisy first render fleet-wide (legacy stamps are genuinely
replaced — a true positive), and stacking an accepted-false-positive class on top is how reports
stop being read. The invariant the kickoff demands survives intact: **an empty report still means no
line disappeared** — moved lines are reported as moved, so empty means neither dropped nor moved.

Two coupled report changes ship in the same slice:

- **`RenderResult` grows `Applied []AppliedOp`** ({op, target/id, lines, hash-check outcome}) — H3
  is real (struct at render.go:35-39 has Body/Removed/Preserved only), and without it the overlay's
  entire payload is invisible in preview and on apply. G1's "reports every block it replaces or
  removes" (consensus.md:255-257) extends honestly: a replace is already a G1 event; an extend is
  reported because silent additions are how the 2026-08-06 damage hid. Both call sites print it
  (protocol.go:202-216, 283-288), which also answers H1 lightly: `check --json` names *which input*
  diverges (core vs lock, overlay vs lock, body vs render) instead of the two-value enum
  (protocol.go:258-274) that conflates four states.
- **Stamp and regex land atomically (D-j).** The stamp gains overlay + effective hashes and is
  derived from the deck lock, not `rel.SHA256` (R3.3; today's stamp lies — render.go:95, H5), and
  `generatedStampRe` (render.go:380) changes in the *same commit*: H4 is confirmed in source — the
  exemption dies silently on any format change and every render then reports a spurious header loss
  (render.go:233-235,384-395). H11's asymmetric guard (core stamp lines dropped at render.go:65-67,
  no overlay equivalent) gets its one-line twin in that commit. H10's fence-blindness (stamp
  insertion via `strings.Replace` at render.go:96-98; unscoped first-prefix `findLine` at
  :160-167) is fixed alongside: scope identity matching to the header block and never substitute
  inside fenced code — mandatory once overlay payloads legitimately *quote* protocol syntax.

### 5. Disqualifying hazards — the overlay may NOT ship onto this substrate as-is

**Disqualifying — must be fixed inside this slice:**

- **H9 (probe-confirmed roster wipe).** The overlay may not ship while identity *extraction*
  prose-matches deck tables. My §3 answer kills the worst of it structurally (rows from config, not
  from parsing decks), and D-t is answered "now" for the residual: core-side table location moves to
  registry block IDs in this slice — cheap precisely because §2's registry ships here. Shipping
  without this means the overlay's first fleet render still empties 11 decks' rosters
  `[BRIEF-TESTIMONY — UNVERIFIED by me: brief:15]` — the exact loss class the overlay exists to end.
- **H2 (repositioning false-loss).** Fixed via §4's moved-category; otherwise DF-2's first act on
  every rescued deck is a false data-loss report.
- **H3 (additions invisible).** Fixed via `Applied`; an overlay whose payload never appears in the
  report is 2026-08-06 wearing a friendlier hat.
- **H4+H5 (stamp lies; regex dies silently).** Fixed via the atomic D-j commit. Either half alone
  is disqualifying.
- **H6 (lock prefix scan silently ignores unknown keys).** Disqualifying the moment the lock gains
  `overlay-sha256:` — today's binary reads such a lock *as if the overlay did not exist*
  (protocol.go:92-98) and then reports the overlay's content as not carried forward. D-i is
  answered now: strict parse, known-key allowlist, unknown key fails closed. yaml.v3 is already in
  go.mod (:10) so this is posture, not a new dependency; and zero production locks exist (H14), so
  strictness breaks nothing that isn't this repo's own tests. §5's rank-2 entanglement is respected:
  this slice writes only `core-version`, `core-sha256`, `overlay-sha256|none` — the minimum D8
  already ratifies (consensus.md:166-167) — and states them.
- **H12 (preflight's shadow writer).** `syncConsumerProtocol` still rewrites consumer decks from the
  packaged *skill* body via `mergePreservingZones` with a `## 3.` byte anchor (preflight.go:488-536,
  anchor at :539) — a coarser preservation model rank 1 demoted to scaffolding
  (consensus.md:224-229). Left live, preflight after `protocol render` reintroduces a body the core
  did not produce and clobbers the overlay's render. In-slice minimum (D-s): that path must refuse
  or no-op on any deck carrying a protocol lock; one pure render function owns COOPERATION.md's
  core-derived bytes, and `roster render` keeps only its §2 UX while sourcing the same table
  generator. Full writer consolidation beyond that is not needed for correctness and is not claimed.
- **H7+H8 (no addressing substrate; s6.6 not heading-delimited).** Subsumed by §2 shipping the
  registry with byte extents.
- **H15 (two promissory notes).** protocol.go:211 and parley-deck/COOPERATION.md:767-768 must
  change in the ship commit — confirmed at those lines; trivial, but non-atomic retirement leaves
  CLI and protocol text lying in opposite directions.

**Not disqualifying — decided or deferred below:** H1 (handled inside the `Applied` call-site work),
H10/H11 (folded into the D-j commit), H13 (D-r decision below), H14 (not a defect — it *is* the
window argument), H16 (test-discipline: G7b end-to-end tests through the real entry point per
protocol_test.go:439-457's pattern, plus real unit coverage inside protocolcore for registry and
resolver — render_test.go is 45 lines of LCS-only today), H17 (true-positive noise; the legacy
stamp genuinely is replaced; the migration note says so), H18 (a DF-2 extraction rule, below).

### 6. Decision ranking — settle in this slice vs defer with a named follow-up

**Settle in this slice (the overlay cannot ship correctly without the answer):**

| Decision | Answer (one line) |
|---|---|
| D-a syntax | YAML frontmatter + Markdown payload sections (§1) |
| D-b registry | Ship now; `registry.yaml` in the release, byte extents + per-block sha (§2) |
| D-c placement | `ext-1` = registry block order, last in v1 (§4) |
| D-d extension IDs | `<deck-slug>.<slug>`; segments `[a-z0-9]([a-z0-9-]*[a-z0-9])?`, ≥1 dot, ≤64 chars; uniqueness within the file; duplicate provider fails closed (D4); deck-slug recorded, not enforced against any authority — none exists to check it |
| D-g loss report | Scoped-(ii) moved-category + `Applied` (§4) |
| D-h missing/unreadable overlay | Fail closed both ways: lock names an overlay hash but file missing/unreadable → **block** (the missing-release analogue, protocol.go:118-128 — rendering without it re-enacts 2026-08-06); lock says `none` but file present → refuse with a named remedy. "Absent + lock says none" is the only free pass (R1.5) |
| D-i lock parsing | Strict known-key parse, fail closed (§5); this slice writes `core-version`, `core-sha256`, `overlay-sha256` only |
| D-j stamp + regex | Core+overlay+effective from the lock; regex in the same commit (§4) |
| D-k §2 roster | Option (c): rows from `agents.toml`; seventh slot for annotations; columns → DF-8 (§3) |
| D-s writers | Minimum: preflight merge-path refuses locked decks; one render function owns core-derived bytes (§5) |
| D-t table addressing | Now: registry IDs for core-side table location; deck-side row extraction deleted (§3/§5) |
| D-w multi-replace | v1 permits at most one replace (exactly one replaceable block exists); ordering spec is due when a second block opens, not before |
| D-x line endings | Deck convention decides output; overlay normalized per-source exactly like core (the render.go:53-57 pattern, R3.5) — one line to decide, a known non-convergence bug if not |
| D-r CLI's own deck | Standing rule: the source repo's deck carries no overlay until the drift guard compares deck-to-*render* (DF-6). Zero overlays exist; nothing dogfoods early |
| D-e reconfirmation artifact | Answered with zero new build: reconfirmation = the attended adoption updating the committed lock (core version + overlay + effective hashes); the lock diff in git *is* the reviewable receipt. A separate receipt file amends D8's ratified field set — do not add one here; rank 2 owns the remaining lock fields (brief:242) |
| D-p payload half | Payloads are free-form Markdown, HTML comments allowed (§1). The promote-librade-to-core half is a separate meta idea if anyone wants it — DF-4 already closed by raw restore (IMPLEMENTATION.md:433-444) |
| D-q provenance | D4's per-operation `rationale:` is the home; the three improvised fleet conventions become rationale text. No separate changelog structure for three decks |
| D-m slot count | Folded into §3's commit: slots become seven, doc comment and D3 language corrected together |

**Deferred, each with a named follow-up:**

| Decision | Follow-up | Why it can wait |
|---|---|---|
| D-e receipt *file* | rank 2 (per-idea pinning) | The lock update already proves reconfirmation; rank 2 owns the remaining lock/manifest fields |
| D-f authoring surface | **DF-7** — `protocol overlay set|remove --dry-run --yes` | v1 is hand-written files in a normal deck idea (R5.1, D11), validated by the existing `render --dry-run`/`check`; no attended gate is implied (D11) and none is added |
| D-l transport: identity or core | **DF-2** sub-decision | The asymmetry is real but benign today: the renderer preserves Transport (render.go:9-20,116) while the drift guard does not exempt it (drift_test.go:24-30) — and the guard runs only on this repo's deck, whose Transport matches the default (parley-deck/COOPERATION.md:5). The servers-deck reconciliation evidence (brief:197) is migration material |
| D-n Workspace derivation | **DF-2** inventory polish | Free-text slot keeps working; derivation-from-directory changes no guarantee |
| D-o handle table under `local-dir` | **DF-8** | Cosmetic noise; one renderer rule + guard update, bundled with the columns work |
| D-v never-synced decks | **DF-2** extraction rule | No overlay extraction from never-synced decks until synced; absent sections are staleness, not local removals |
| H18 placeholders | **DF-2** extraction rule | Extraction ignores template placeholders; R1.6 (no auto-created overlays, round-01/codex-1.md:149-151 of the meta idea per brief:36) already forbids the pollution path |
| D-r implementation | **DF-6** — render-based drift guard | Guard compares deck to rendered (core, overlay, slots); lifts the standing rule above |
| Roster columns + guard reconciliation | **DF-8** — `roster render` v2 | Grows Model/Reasoning/CLI columns from `agents.toml`, suppresses the handle table under `local-dir`, updates the drift guard's `rosterHeaderLine` anchor in the same commit. **Ordering constraint: DF-8 must land before DF-2's fleet renders**, or first renders drop bespoke columns from views fleet-wide (data safe in config, but needlessly alarming reports) |

### 7. D-u — is user constraint 3 currently satisfied?

**No — and this slice is its discharge.** My rev-1 signoff accepted D3's near-empty surface "on
condition that rank 3 actually ships this cycle (… if it slipped to deferred, constraint 3 would be
unmet)" (signoffs/kimi-1.md:95-97 of the meta idea — my own prior artifact, re-read this session).
Rank 3 did not ship: it is "Deferred, ratified, not built" (IMPLEMENTATION.md:418-419) and DF-2 is
blocked on it (IMPLEMENTATION.md:471-473). My rev-3 ACCEPT carried "Conditions: None" and covered
only the four retained textual fixes (signoffs/rev3/kimi-1.md:15-57) — it did not waive the rank-3
condition, which was about the shipped cycle outcome, not the consensus text. So the record should
state plainly: **constraint 3 stands unmet today; it is discharged exactly by this idea shipping the
overlay; if this idea slips again I will re-raise it at the meta level** as a failed owner
constraint, not a design preference.

Unfinished business reclaimed: my `parley protocol audit` fleet surface (overlay count, targeted
IDs, core-version spread) was dropped from the meta idea with no follow-up number, though my own
signoff proposed DF-5 (signoffs/kimi-1.md:58-60). **DF-5 is hereby named**: the audit surface, to be
scheduled after DF-2. With D-d's deck-namespaced IDs and a committed lock per deck, the audit is a
read-only query over committed files — cheap precisely if this slice lands its parts.

## Concerns / open questions

- **The taxonomy is unverified testimony in this round.** My §3 design leans on the *shape* of the
  29-deck survey, which in-repo code corroborates (identity slots = exactly the header+§2 zones,
  render.go:9-20), but the counts (11/29 bespoke schemas, 23/29 toml disagreement, 3/29 annotation
  decks) are `[BRIEF-TESTIMONY — UNVERIFIED by me]`. If any participant can verify them with
  cross-workdir access, a CONFIRMED verdict would strengthen the seventh-slot proposal; if the
  counts collapse, the fallback (annotations at `ext-1`) becomes relatively more attractive. Flagged
  per §15.1 rather than silently relied on.
- **Seventh slot vs D3's ratified enumeration.** I judged this an implementation-level amendment
  (§3), but a participant could reasonably read D3's "six identity slots" as frozen by user
  ratification of the meta idea. If that reading wins, the annotation class falls to `ext-1` and I
  want the consensus to record the discoverability cost it accepted.
- **Interaction with rank 2 is one-directional only if D-i lands.** If the strict lock parse slips,
  rank 2's new fields inherit H6's silent-read vector. D-i is cheap now and expensive later; I will
  resist deferring it.
- **Byte extents make `Publish` the only extent computer.** Tooling must never accept a hand-written
  registry. G2's attended-publisher gate (consensus.md:257-260) already scopes this correctly, but
  the acceptance tests should assert a hand-forged registry/body mismatch fails closed at load.

## Risks

- **Scope accretion.** This file already assigns the overlay slice: registry, strict lock, stamp
  rework, writer-path guard, seventh slot, report rework. Each is individually justified by a
  measured hazard, but the slice is large; if it must shrink, the registry (D-b) and H9/D-k are the
  last two things to cut — cutting either ships an overlay that cannot address its one block or that
  ships onto a roster-wiping substrate. Everything else is more deferrable than those.
- **False confidence from a green build.** The H13 rosterHeaderLine/roster-render mismatch means the
  tree already contains two mutually inconsistent table contracts that no test currently fails on
  together. DF-8's atomic reconciliation must come with the compiling-revert discipline of R6.5
  (IMPLEMENTATION.md:382-388) or we will "fix" it into a third contract.
- **The moved-category weakens operator alarm over time.** Scoped-(ii) is safe today because moved
  lines require a declared operation; if a future version widens suppression, the G1 report's alarm
  value erodes. The scope limit should be written next to the code, not just in this file.
- **My own conflict of interest, disclosed:** D-u's resolution and DF-5's naming both reclaim
  positions from my prior signoffs. I have tried to argue them from the constraint text and the
  record rather than from authorship; round 2 should test them on those grounds.
