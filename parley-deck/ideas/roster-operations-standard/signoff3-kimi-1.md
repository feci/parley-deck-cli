### kimi-1 — revision 3

**Verdict: accept with reservations**

codex-1's four counter-proposal items are met in the operative text, with one deliverable from
item 1 outstanding (the verbatim §2 replacement prose — my own K4 always scoped that to FINAL.md).
My K1 and K2 are adopted into the consensus; K3 and a reduced K4 remain binding FINAL.md
requirements, joined by two new carried items (K5 staging plan, K6 skill-snapshot locator). The
load-bearing drafter-owned claim — that workspace dir, role and host handle are render-only —
survives a non-owner source check; details in §3.

**Scope declared.** Read in full this session (all `PRIMARY`): `parley-deck/COOPERATION.md` §15
(`:1176-1316`) and the §2 roster/host-handle source (`:95-134`); the full revision-3
`consensus.md` (992 lines), including all four embedded revision-1 signoffs (`:491-989`); my
`signoff2-kimi-1.md`; codex-1's `signoff2-codex-1.md` in full. Read partially:
`signoff2-hermes-1.md` — section index plus `:425-486` (verification notes, reservation tail).
Fresh checks executed this session (all `PRIMARY`): full read of `internal/protocol/roster.go`
(70 lines); reads of `internal/app/roster.go:95-134`, `internal/app/app.go:2400-2429`,
`internal/app/preset.go:40-69`, `internal/config/roster.go:60-109`; grep of every `*.go` for
`workspace_dir|Workspace dir|host_handle|Host handle`; grep of every `*.go` for `ReadRosterIDs`
and for `COOPERATION.md`; grep of `consensus.md` for staging language (zero hits); a
repository-wide `find` for COOPERATION.md copies (two); `ls` of the user-scope skill directories
(`~/.kimi-code/skills/parley-deck/`, `~/.claude/skills/parley-deck/`). I did **not**: re-enumerate
the 40 decks; run any `parley` binary or test; read `~/.parley/*` or any foreign deck; re-verify
revision-2-unchanged sections against source (release conditions, migration contract — my
revision-2 `PRIMARY` reads stand, and the text matches at shifted locators); re-run the
byte-for-byte signoff-embedding diff from revision 2; read hermes-1's revision-2 signoff body
(`:50-424`); or inspect any git state. I ran no git command and wrote only this file. Per §15.1
(`COOPERATION.md:1197-1205`) I issue no verdict on claims I own: the `printUsage`/docs omission
(`consensus.md:40-44`), the `DISPLAY-NAME` contradiction (`:76-79`), the discarded inactive set at
`internal/app/roster.go:110`, and the parser-populates-inactive claim from round 1
(`internal/protocol/roster.go:62-64`) — the last two are cited below only as `SECONDARY` via
hermes-1's `PRIMARY` reads (`consensus.md:642-652`; `signoff2-hermes-1.md:427-434`).

#### 1. codex-1's four counter-proposal items (`signoff2-codex-1.md:56-64`)

**Item 1 — normative field table, retention rule, ordering rule, §2 replacement text: MET IN
SUBSTANCE; one deliverable outstanding.** The field table exists and is per-field complete: nine
rows giving committed TOML key, legacy §2 source, absence/conflict behaviour, and runtime-semantic
vs render-only (`consensus.md:353-363`). Inactive-history retention is defined twice,
consistently: the `state` row's "Mark inactive; never delete — history is retained permanently"
(`:357`) and the migration gate's "retained as `active = false`, never removed" (`:387-389`). One
deterministic ordering rule: active before inactive, then agent ID byte-ascending (`:365-367`).
**What is still missing: the verbatim §2 replacement prose.** `:374-384` are normative authority
statements — what the replacement text must establish — not the replacement markdown that will be
committed to `COOPERATION.md` §2. My K4 scoped the verbatim text to FINAL.md, so on my own
requirements this is not a block; codex-1's item 1 demanded it "now, because the consensus itself
says this is required before ratification" (`signoff2-codex-1.md:60`), and on that stricter
reading item 1 is short by exactly that one artifact. codex-1 judges its own demand; the record
should be precise that the definitional deficit revision 2 was blocked for is cured, while the
wording itself is not yet in evidence.

**Item 2 — §7-format `meta/protocol-changelog.md` entry: MET.** Required at `consensus.md:383-384`,
naming this idea and the user-authorized one-off. This also answers hermes-1's revision-2 FINAL.md
note (`signoff2-hermes-1.md:478-479`).

**Item 3 — foreign-deck compatibility gate + explicit retired-row retention: MET.** "A deck whose
protocol/schema version predates this change is skipped and reported, not silently upgraded.
Retired-agent rows are retained as `active = false`, never removed" (`consensus.md:386-389`). The
gate is a named condition distinct from the "unsupported legacy layout" unclean class (`:453-455`)
— exactly the separation my K2 required.

**Item 4 — kimi-1's R4: MET, as both halves.** `--keep <agent>.<field>` ships, **and** the dry-run
and final report must enumerate every deliberate pin the rebase removes, per deck
(`consensus.md:391-394`). See §2.

**The staging plan is NOT recorded.** I was asked to confirm it is recorded faithfully; I cannot,
because it is absent (`PRIMARY`: grep of `consensus.md` for `stage|staged|atomic group` — zero
hits). Neither codex-1's four "Delivery shape" stages (`signoff2-codex-1.md:47-54`) nor my four
stages (`signoff2-kimi-1.md:186-219`) appear in revision 3. What IS recorded are the three atomic
couplings that make any staging safe: snapshot persist+consume with rebase — "Rebase must not
ship first" (`:319-327`); STATE wiring in the same change as the migration (`:396-402`); protocol
change + generated §2 + embedded copy + skill snapshot in one release (`:374-384`). The plan
itself — stage boundaries, what may overlap, what must not start early — remains unwritten.
Carried as K5.

#### 2. hermes-1 / kimi-1 positions

**hermes-1 — nothing weakened.** hermes-1 accepted revision 2 with no carried reservations
(`signoff2-hermes-1.md:5,471-476`). Everything that acceptance rested on is present in revision 3
at shifted locators: release conditions 1-3 (`:317-327`), R2 idempotency (`:404-406`), R3.1 STATE
prerequisite (`:396-402`), R3.2 per-deck confirmation (`:460-461`), R3.3 attended-only (`:464`),
R3.4 file-level backups (`:456-459`), R3.5 resumability (`:462-463`). Revision 3's changes are
additive; I found no revision-2 operative text removed or weakened (comparison against my
revision-2 signoff's quoted citations, which match at the shifted locators; a byte-diff was
impossible — the revision-2 file is overwritten). The R3.1 refinement (`:348-351`) strengthens
hermes-1's point rather than weakening it — see §3, claim C.

**kimi-1 — nothing weakened; R4 is satisfied by adopting both halves.** K1 is adopted in full and
K2 is adopted (§1, items 3-4). K3 (destructive sync against pre-snapshot decks warns and requires
the second confirmation) is unchanged — still a FINAL.md requirement, contradicted nowhere. K4's
changelog half is adopted (`:383-384`); the verbatim-text half carries. On R4 specifically: I
demanded either/or (`signoff2-kimi-1.md:171-173`). Both halves together are strictly stronger and
close each half's hole: `--keep` alone requires knowing every deliberate pin in advance across 40
decks; enumeration alone is retrospective. Shipping both means pins can be exempted prospectively
where known and recovered mechanically where not. **R4: satisfied.**

#### 3. The field table against source — non-owner verdicts

The drafter owns the revision-3 measurement (`consensus.md:341-346`); §15.1 makes a non-owner
verdict the admissibility gate. I am a non-owner for these claims — my round-1 claims are disjoint
(the discard at `roster.go:110`; the inactive map being populated).

**Claim A — `ReadRosterIDs` extracts only the agent ID and the literal `inactive`: CONFIRMED
(`PRIMARY`).** `rosterRowRe` captures one group, the first-cell ID (`internal/protocol/roster.go:17`),
used at `:56-60`; the only other extraction is the case-insensitive substring test for `inactive`
at `:62-64`; the header regex at `:19` anchors the table. No other cell is read anywhere in the
function (`:26-70`).

**Claim B — zero non-test consumers of workspace dir / host handle: CONFIRMED (`PRIMARY`).** My
grep of every `*.go` for `workspace_dir|Workspace dir|host_handle|Host handle` returns six hits,
all in test fixtures: `internal/app/roster_test.go:48,153`, `internal/protocol/drift_test.go:28-29`,
`internal/protocol/roster_test.go:16,22`. No non-test Go code names these columns in any casing
convention.

**Claim C — every row enters `active`, including inactive rows; `inactive` is a separate map:
CONFIRMED (`PRIMARY`).** `active[id] = true` executes unconditionally on every matched row
(`roster.go:61`); `inactive[id] = true` is set additionally at `:62-64`. An inactive agent is in
both maps. The drafter's refinement of hermes-1's R3.1 is accurate.

**Claim D — therefore workspace dir, role and host handle are render-only today: CONFIRMED
(`PRIMARY`).** For workspace dir and host handle this follows from A+B. For role: the only §2
roster parser in the codebase is `ReadRosterIDs` (callers: `internal/app/roster.go:110`,
`internal/app/preset.go:50`, `internal/app/app.go:1793,2412`; consumer signature
`internal/config/roster.go:84`), which reads no role cell, and no other Go code parses §2 roster
prose; `COOPERATION.md:95` independently marks role metadata advisory. Nothing a launch or render
decision depends on reads these three fields. The table's "legacy §2 source" column is also
accurate: col 1 Agent ID / col 2 Workspace dir / col 3 Role carrying the
``(cli `claude`, model `…`)`` prose (`COOPERATION.md:105-110`); the host-handle table at
`:119-126`; and the retention rule's provenance checks out (`:134` — "mark its row as inactive
(do not delete it)").

**Precision correction — strengthens, does not weaken, R3.1.** The sentence "Marking a row
inactive is cosmetic today" (`consensus.md:399`) is over-broad as written. It holds for the roster
display path — `resolveRoster` discards the inactive map (`internal/app/roster.go:110`; owned by
me, cited `SECONDARY` via hermes-1's `PRIMARY`) — but the inactive set is already consumed at
three other sites: `defaultRosterParticipants` skips inactive IDs (`internal/app/app.go:2412-2424`,
the `continue` at `:2418-2420`), preset validation receives it fail-closed
(`internal/app/preset.go:50-53`, `internal/config/roster.go:82-84`), and `app.go:1793` receives
it. Migrating retired rows to `active = false` is therefore not a no-op for default participant
selection even today; it IS a no-op for the `roster show` rendering, which is the defect R3.1
targets. The requirement stands unchanged; FINAL.md should scope the sentence to the display path.

**Locator WRONG — the "three copies" enumeration.** The authority-flip requirement cites the
skill's bundled snapshot as `skills/parley-deck/references/COOPERATION.md`
(`consensus.md:380-381`). That path does not exist in this repository (`PRIMARY`: repository-wide
`find` returns exactly two COOPERATION.md files — `parley-deck/COOPERATION.md` and
`internal/protocol/defaults/COOPERATION.md`; there is no `skills/` tree). The bundled snapshots
exist at user scope: `~/.kimi-code/skills/parley-deck/references/COOPERATION.md` and
`~/.claude/skills/parley-deck/references/COOPERATION.md` (`PRIMARY`: `ls`). As written the locator
is **WRONG** (non-owner verdict; the drafter owns the claim). The requirement itself is sound and,
corrected, must enumerate the actual copies — note there is more than one installed skill
snapshot, so the drift guard is wider than "three copies". Carried as K6.

#### 4. Record accuracy (non-binding)

- The revision-3 preamble's change list (`:12-17`) and the drafter position changelog — fifteen
  entries, four added by this block, entries 12-15 (`:267-274`) — match the document: CONFIRMED
  (`PRIMARY`, read of the table against the cited sections).
- VC-2's header still reads "OPEN" (`:188`) while the body records closure by user direction
  (`:306`), and the preserved framing still camps me as "additive, source-aware pin" (`:192-193`)
  although `round-02/kimi-1.md:22-29` adopted rebase with the `--keep` amendment. Carried from my
  revision-2 signoff, still outstanding; FINAL.md's §15.3 record should quote positions as they
  stood.
- `### Revision 2 — signoffs pending` (`:991`) is a dangling header: the three revision-2 signoffs
  exist as standalone files but are not embedded, unlike the revision-1 blocks (`:487-489`).
  Embed them or strike the header before FINAL.md.
- The fleet figures (40 decks; 17 rosterless; 17 retired-`antigravity-1`; 3 `gemini-1`; 1 `agy-1`)
  remain `SECONDARY` (claude-1's `PRIMARY` measurement, unreproduced by me). Whether `SPEED`
  shares the declared/effective defect remains unmeasured by anyone, me included (`RECALL`;
  recorded as unmeasured at `:294-295`).

#### 5. Carried reservations (checkable FINAL.md/implementation requirements)

- **K3 (unchanged from revision 2).** A `roster sync`/`set` that removes a masking override or
  deactivates a row in a deck with runs created before snapshot support requires the
  breaking-change second confirmation and prints the reconstructability warning; FINAL.md states
  this explicitly, not merely via decision 5's safety list.
- **K4 (reduced).** FINAL.md quotes the verbatim §2 replacement text, scope-limited to the
  authority move plus the generated non-authoritative view, with no code path parsing §2
  afterward. The `meta/protocol-changelog.md` entry requirement is now in the consensus
  (`:383-384`) and drops out of K4.
- **K5 (new).** FINAL.md states the staged delivery plan — stage boundaries, permitted overlaps,
  fleet migration last — preserving the three atomic couplings already in the consensus
  (`:319-327`, `:374-384`, `:396-402`). codex-1's `signoff2-codex-1.md:47-54` and my
  `signoff2-kimi-1.md:186-219` agree on the couplings and can be reconciled there.
- **K6 (new).** The release's protocol-copy enumeration is corrected to the real artifacts: the
  live deck copy, `internal/protocol/defaults/COOPERATION.md`, and every installed skill snapshot
  (at minimum the two user-scope copies found today), replacing the nonexistent
  `skills/parley-deck/references/COOPERATION.md` locator.

With K3-K6 recorded, I sign revision 3.
