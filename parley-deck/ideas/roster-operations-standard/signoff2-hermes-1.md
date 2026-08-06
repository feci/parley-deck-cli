### hermes-1 — revision 2

Date: 2026-08-06

**Verdict: accept**

Revision 1 filed three reservations (R1, R2, R3 with five sub-points) and an
ACCEPT-WITH-RESERVATIONS. Revision 2 converts every reservation into binding
consensus text, upholds codex-1's four-part block in full, and closes VC-1 and
VC-3 the way I argued they should close. I accept without reservation.

---

#### Scope declared (§15.1, §15.2)

- `PRIMARY` — I read `parley-deck/COOPERATION.md` in full (§1–§15, including
  §4.0 at `:172-228`, §6 rule 3 at `:708`, §7 at `:717-731`, §14.2 at
  `:1153-1161`, §15.1–§15.7 at `:1183-1316`). I read the full revision-2
  `consensus.md` (938 lines) including all four embedded revision-1 signoffs
  and my own revision-1 signoff at `:498-677`.
- `PRIMARY` — I read my own `round-01/hermes-1.md` (168 lines) and
  `round-02/hermes-1.md` (127 lines) in full this session.
- `PRIMARY` — Fresh source checks this session: `internal/app/roster.go:100-139`
  (`resolveRoster` discards inactive into `_` at `:110`; `rosterTargetPath`
  maps non-machine scope to `agents.toml` at `:383-389`),
  `internal/protocol/roster.go:22-69` (`ReadRosterIDs` populates the inactive
  map at `:62-63` but returns it to a caller that drops it),
  `internal/runmanifest/manifest.go:28-46` (`Manifest` carries `Participants
  []string` at `:43` and no model/effort/adapter/snapshot field),
  `internal/app/app.go:1145-1169` (`continueAuto` re-discovers config at
  `:1148` and passes it as `Agents: discovered` at `:1159`).
- `SECONDARY` — I rely on claude-1's `PRIMARY` 40-deck fleet measurement
  (`consensus.md:135-137`, sourced from the inbox measurement note) for the
  fleet figures. I did not re-enumerate the 40 decks.
- `SECONDARY` — I rely on codex-1's `PRIMARY` reads of
  `internal/runcontrol/runcontrol.go:152-177` (the `run.created` event records
  declared `model`/`reasoning`/`sources` but no materialized invocation plan)
  and `docs/agent-runtime-configuration.md:7-15` (precedence of
  `agents.local.toml` vs `agents.toml`).
- I did not run any live `parley` command, run tests, inspect foreign decks,
  or read `~/.parley/agents.toml` or `~/.hermes/config.yaml` this session.
- Per §15.1 I issue no verdict on any claim I own (my round-1/round-2 findings:
  the `resolveRoster` inactive-discard, the `rosterTargetPath` mapping, the
  EFFORT declared/effective split for hermes/codex/kimi, the
  `meta/headless-agents.local.json` non-reader). I verdict codex-1's and
  claude-1's claims as a non-owner below.

---

#### codex-1's four required changes — are they now met?

codex-1's block (`consensus.md:437-496`) listed four required changes. I assess
each as a non-owner (§15.1 — codex-1 owns these claims; I verdict them).

**1. Snapshot-consumption acceptance gate + legacy-run skip (consensus.md:493
item 1, and the three binding release conditions at `:307-315`).**
MET. `PRIMARY` (I read `:307-315` directly): condition 1 requires the change
exposing rebase to also persist AND consume the immutable effective snapshot,
with an acceptance test — "create a run, change machine/deck config, continue
the run, and prove adapter, model, effort and autonomous-write args are
unchanged." Condition 2 requires fleet migration to skip and report every
nonterminal legacy run lacking that snapshot, and forbids manufacturing
snapshots by rewriting `participants:` or run artifacts. Condition 3 requires
that if the gate is relaxed, the residual risk is stated as "unsafe for
pre-snapshot resumable runs." This matches codex-1's three items at
`:457-459` item-for-item. The coupling is now a binding release gate, not the
"practical effect" hand-wave revision 1 offered (`:254` entry 7 documents the
prior wording). My own R1 (revision-1 `:544-553`) demanded exactly this:
"rebase and the immutable run snapshot must ship as one atomic delivery unit;
if the snapshot implementation slips past the rebase implementation, a window
exists where rebase is live but reproducibility is not guaranteed." Revision 2
makes R1 binding.

**2. §7 deviation wording corrected + track upgraded to `deliberation` + §2
authority spec completed (`consensus.md:494` item 2).**
MET on all three sub-parts.

  - **§7 wording.** `PRIMARY` (I read `:354-368` directly): revision 2 calls
    the deviation "an explicit user-authorized ONE-OFF, not a protocol
    exception" and states it "rests on the user's explicit one-off
    authorization alone and creates no general precedent." It explicitly
    corrects the mis-citation: "Revision 1 called this 'the protocol's
    direct-user-instruction exception'. That was a mis-citation by the
    drafter. The only such exception in the protocol is `COOPERATION.md` §6
    rule 3, and its text is scoped to 'Never edit another agent's file'".
    I verified the source independently: `PRIMARY` —
    `COOPERATION.md:706-708` shows §6 rule 3 is "Never edit another agent's
    file" with the exception scoped to direct user instruction for that
    specific rule. `COOPERATION.md:717-731` (§7) contains no
    direct-user-instruction exception — only the version-sync carve-out at
    `:726-730`. codex-1's block at `:463-465` was correct, and revision 2
    corrects it verbatim. This also resolves my own revision-1 signoff's
    error: I wrote that "this is the §6 rule 3 direct-user-instruction
    exception applied to §7's process requirement" (`:561-562`). That was
    wrong for the same reason codex-1 identified — §6 rule 3's exception is
    scoped to file-editing, not to §7. Revision 2's wording is the correct
    one, and it is stricter than what I accepted in revision 1.

  - **Track upgrade.** `PRIMARY` (consensus.md frontmatter `:5`):
    `track: deliberation`. `:370-374` documents the upgrade: the §4.0
    classifier (`COOPERATION.md:181`) forces `deliberation` if any trigger
    fires; this idea fires three — protocol change (§7), data migration, and
    an irreversible/destructive fleet operation — and the fail-closed rule
    independently requires the stricter track. I verified the classifier:
    `PRIMARY` — `COOPERATION.md:181` lists the `deliberation` triggers
    including "protocol change (§7); ... data migration / irreversible /
    destructive op." All three fire here. MET.

  - **§2 authority spec.** `PRIMARY` (I read `:326-348` directly): revision 2
    requires, for each of workspace dir, role, host handle, active/inactive
    history, and row ordering — (a) which file is the canonical source,
    (b) the migration path for existing values, (c) that the generated §2 is
    non-authoritative and is a rendering not a store, (d) that runtime code
    MUST NOT parse the generated view as roster authority (citing
    `resolveRoster` at `internal/app/roster.go:110`), and (e) that every
    other protocol reference calling §2 authoritative, plus the embedded
    protocol copy and the skill's bundled snapshot, must change in the same
    release. This matches codex-1's requirement at `:469` item-for-item. I
    confirmed the §2 contents codex-1 names: `PRIMARY` —
    `COOPERATION.md:101-117` stores Agent ID, Workspace dir, and Role;
    `:119-126` stores the host-handle table. The commands at `consensus.md:92`
    manage adapter/state/model/effort/speed only, so without this spec a
    generated view would drop workspace dir, role, and host handle. MET.

**3. Fleet migration contract expanded (`consensus.md:495` item 3).**
MET. `PRIMARY` (I read `:390-414` directly): revision 2's migration contract
adds all five of codex-1's requirements — (1) machine-readable inventory
(`:393-395`), (2) compare-and-swap between dry-run and apply with batch
approval and second confirmation for membership changes (`:396-398`), (3)
precise "unclean" definition including dirty worktree, parse/validation error,
unsupported legacy layout, path/symlink ambiguity, concurrent modification,
and nonterminal pre-snapshot run (`:399-401`), (4) backups with recorded
location and hashes, verified restore procedure, atomic per-deck writes,
post-write validation, automatic rollback, and file-level (not git) backups
(`:402-405`), (5) final machine-readable report marking every deck
applied/skipped/failed-and-restored/unchanged with before-and-after hashes and
backup reference, and no automatic commit/push/edit to locked participants or
historical run artifacts (`:411-414`). This matches codex-1's five items at
`:477-481`. MET.

**4. VC-1 closed by engagement + VC-3 ratified as `deck|machine` with
`--scope deck` targeting committed `agents.toml` (`consensus.md:496` item 4).**
MET. VC-1 is closed at `:155-165`: kimi-1 withdrew `SOURCE` at signoff and the
resolution is by argument (one column can only name the winning layer for one
field; a `SOURCE` header is read as row provenance regardless of docs), not by
the 3-to-1 count — which is the §15.3 requirement. VC-3 is closed at
`:222-230`: `deck|machine` unanimous, `--scope deck` writes the committed
`parley-deck/agents.toml` never the gitignored `agents.local.toml`, grounded in
codex-1's source citation (`docs/agent-runtime-configuration.md:7-15`,
`internal/config/runtime.go:134-151`). MET.

All four of codex-1's required changes are met. Nothing is missing.

---

#### hermes-1 — are R1, R2, and R3 (all five sub-points) now addressed?

**R1 (revision-1 `:544-553`): rebase + snapshot must ship as one atomic
delivery unit.**
ADDRESSED, and upgraded from a reservation to a binding release gate.
`consensus.md:307-315` condition 1 requires the rebase-exposing change to also
persist and consume the immutable snapshot, with an acceptance test. My R1
said "reproducibility must not depend on an unshipped feature shipping later;
it must ship with rebase or rebase waits." Revision 2 condition 1 says exactly
that, and condition 3 adds the "unsafe for pre-snapshot resumable runs"
residual-risk statement if the gate is ever relaxed. R1 is satisfied and
exceeded.

**R2 (revision-1 `:575-580`): generated §2 must be idempotent and preserve
human-readable prose.**
ADDRESSED. `consensus.md:350-352`: "Generated §2 must be idempotent
(hermes-1's R2). Running the generator twice produces byte-identical output,
and the human-readable prose shape is preserved. A non-idempotent generator
recreates the drift under a new name." This is my R2 verbatim in intent. The
consensus also requires at `:337` that "runtime code MUST NOT parse the
generated view as roster authority" — which is the deeper form of R2: not just
idempotent generation, but no dual-store parsing at all. R2 is satisfied.

**R3 (revision-1 `:627-629`): migration guardrails, five sub-points.**

- **R3.1 — inactive-set wiring as hard prerequisite.** ADDRESSED.
  `consensus.md:342-348`: "`STATE` wiring is a hard prerequisite for the
  migration (hermes-1's R3.1, confirmed by the drafter as a non-owner).
  `resolveRoster` reads `active, _, ok := protocol.ReadRosterIDs(root)` — the
  inactive map is discarded into `_` (`internal/app/roster.go:110`) ...
  Marking a row inactive is cosmetic today ... Decision 3's `STATE` column and
  the inactive-set wiring MUST ship in the same change as the migration." I
  re-verified this session: `PRIMARY` — `internal/app/roster.go:110` reads
  `active, _, ok := protocol.ReadRosterIDs(root)`; `internal/protocol/roster.go:62-63`
  populates `inactive[id] = true` but the caller drops it. The consensus
  correctly identifies the coupling and makes it binding. R3.1 satisfied.

- **R3.2 — per-deck confirmation, not bulk `--yes`.** ADDRESSED.
  `consensus.md:406-407`: "Per-deck or small-batch confirmation, not one bulk
  `--yes` across 40 (hermes-1 R3.2), honoring `roster_change_policy =
  \"confirm-breaking\"`." This is my R3.2 directly. Satisfied.

- **R3.3 — human-attended only, never from a loop/cron/CI hook (§14.2).**
  ADDRESSED. `consensus.md:410`: "Human-attended only, never from a loop, cron
  or CI hook (hermes-1 R3.3, §14.2)." I verified §14.2 is the right citation:
  `PRIMARY` — `COOPERATION.md:1153-1161` ("What an automated loop MUST NOT do
  without a recorded human or full-quorum gate") includes "Modify the active
  roster (§2)" at `:1159`. R3.3 satisfied and correctly grounded.

- **R3.4 — file-level backups for dirty and non-git decks.** ADDRESSED.
  `consensus.md:404-405`: "Backups are file-level copies, not git operations —
  several decks may not be git repositories or may have uncommitted work
  (hermes-1 R3.4)." This is my R3.4 directly. Satisfied.

- **R3.5 — resumability.** ADDRESSED. `consensus.md:408-409`: "Resumability
  (hermes-1 R3.5): a crash on deck 23 leaves 1-22 in a known state, a re-run
  resumes, and an already-migrated deck is a no-op." This is my R3.5 directly.
  Satisfied.

All five R3 sub-points are addressed. R1, R2, and R3 are fully resolved.

---

#### kimi-1 — are its reservations addressed, and is VC-1 closure recorded as it argued?

kimi-1's revision-1 signoff (`consensus.md:679-856`) filed four reservations
(R1–R4) and several accuracy corrections. I assess these as a non-owner
(§15.1 — kimi-1 owns these positions).

**R1 (kimi-1, `:837-841`): rebase coupling as release gate, with resume
consuming the snapshot not merely comparing it.**
ADDRESSED. `consensus.md:307-315` condition 1 requires the snapshot to be
persisted AND consumed (not merely compared), with a continuation acceptance
test. kimi-1's R1 specifically required "resume consuming the snapshot, not
merely comparing it" — the consensus says "persist and consume the immutable
effective snapshot" (`:307-308`). The destructive-sync second confirmation and
reconstructability warning at `:840-841` map to revision 2's condition 3
(`:313-315`). Satisfied.

**R2 (kimi-1, `:842-846`): §7 follow-through — verbatim §2 replacement text in
FINAL.md, `meta/protocol-changelog.md` entry, and stated one-time-vs-codified
disposition.**
This is a FINAL.md requirement, not a consensus-text requirement. The
consensus at `:326-368` sets up the authority for it: it names the user
direction as authority, requires the generated §2 to be non-authoritative with
no code path parsing it, and explicitly records the one-off nature at
`:354-368`. The changelog entry and the one-time-vs-codified statement are
FINAL.md artifacts. The consensus does not explicitly name
`meta/protocol-changelog.md`, but it does require at `:339-340` that "every
other protocol reference calling §2 authoritative, plus the embedded protocol
copy and the skill's bundled snapshot, must change in the same release" — which
is the substantive content. I note this as a FINAL.md carry-item, not a gap in
the consensus. The consensus text is sufficient to require it.

**R3 (kimi-1, `:847-851`): migration constraints extended — locked
participants/live runs untouched, retired agents marked inactive never deleted,
foreign-deck protocol-version gate, written skip-class definitions, post-apply
re-resolution, documented restore path, stated git disposition, fleet form of
membership second confirmation.**
ADDRESSED across `:390-414` and `:342-348`. The inactive-never-deleted rule is
grounded: `PRIMARY` — `COOPERATION.md:134` says "mark its row as inactive (do
not delete it)." The consensus at `:342-348` wires the inactive set and the
migration marks rows inactive. Compare-and-swap (`:396-398`) covers the
locked-participants and live-runs protection — condition 2 says migration
"skips and report[s] every nonterminal legacy run lacking that snapshot" and
"Existing `participants:` lists and run artifacts are never rewritten"
(`:312-313`). The skip-class definition at `:399-401` covers kimi-1's named
classes. The final report at `:411-414` covers the per-deck enumeration. The
foreign-deck protocol-version gate is implicitly covered by the "unclean"
definition (a deck whose protocol copy is at an incompatible version would be a
parse/validation skip), though the consensus does not name "protocol-version
gate" explicitly. This is a minor gap in specificity, not in coverage — the
unclean definition is broad enough to capture it. I flag this as a note for
FINAL.md, not a block.

**R4 (kimi-1, `:852-853`): `--keep` ships, or per-deck enumeration of removed
pins is mandatory.**
This is a command-surface and report detail. The consensus at `:92-96` shows
the `roster set` command without `--keep`, and the migration report at
`:411-414` requires before-and-after hashes per deck. The per-deck enumeration
of removed pins is not explicitly named in the report contract, but the
compare-and-swap dry-run diff (`:396-398`) would surface removed pins in the
diff that goes to the user before apply. This is a FINAL.md/implementation
detail, not a consensus gap. I note it for the implementer.

**VC-1 closure — is it recorded as kimi-1 argued it?**
YES. `PRIMARY` (I read `:155-165` directly): "VC-1 — Does the canonical table
carry `SOURCE`? CLOSED in revision 2 — excluded, by argument." The resolution
states kimi-1 "withdrew its own proposal at signoff" and records the argument:
"the row-wide version is incoherent because `MODEL`, `EFFORT`, `SPEED` and
`AUTO` can each win at a different layer; the narrowed `MODEL`-scoped version
fails on its own merits because a frozen additive-only API column must carry
permanent width, `STATUS` already flags the surprising cases ... and a header
named `SOURCE` will be read as row provenance regardless of documentation."
It explicitly states "Three of the four reached that position by reversing
their own earlier one (claude-1 and hermes-1 in round 2, kimi-1 at signoff),
which is what makes this resolution-by-argument rather than 3-to-1 attrition."
This is exactly the resolution kimi-1 argued for at `:794-806`: "VC-1 is
resolved by argument, not by the 3-to-1 count, which is the resolution §15.3
requires." The closure is recorded correctly and as kimi-1 argued it.

**kimi-1's accuracy corrections (`:821-833`).** The VC-2 mislabeling (kimi-1
was in the rebase camp after round-2, not additive-pin) is corrected: the
"Original positions, preserved" at `:167-179` preserves the round-1 positions
for the record, while the user direction at `:290-294` closes VC-2 by
selecting rebase. The `--keep` disappearance is noted as R4 above — it is a
carry-item, not a consensus error.

kimi-1's reservations are addressed. VC-1 is closed as it argued. The two
minor FINAL.md carry-items (explicit protocol-version gate naming, `--keep` or
removed-pin enumeration) are implementation details within the consensus's
already-broad contract, not gaps that block signoff.

---

#### All — should this ship as ONE change or be staged?

This should be STAGED. The implementation is large — placeholders + resolver,
11-column contract + JSON schema, `roster set`/`sync`, `modelmeta` registry,
`STATE` wiring, run snapshot, generated §2 + protocol change, skill update,
migration command, docs — and several internal dependencies make a single
atomic delivery risky. But the rebase gate (`consensus.md:307-315`) already
forces snapshot+rebase to be atomic, and the `STATE`-wiring prerequisite
(`:342-348`) forces STATE+migration to be atomic. So the staging is not
arbitrary; it follows the dependency graph the consensus already identifies.

Proposed stages:

**Stage 1 — the data contract and display layer (must land as one change).**
This is the foundation everything else depends on. It includes:
- The 11-column contract + JSON schema with `schema_version`/`columns` and
  golden tests (`consensus.md:59-76`).
- The `modelmeta` resolver (CLI-owned, versioned, tested) with gateway-prefix
  peeling (`:78-88`).
- `STATE` column + wiring up the inactive set in `resolveRoster`
  (`internal/app/roster.go:110` — stop discarding the inactive map). This is
  the hard prerequisite for the migration (`:342-348`).
- `{model}`/`{effort}` placeholder substitution in `HeadlessArgs` + the legacy
  normalizer for deck overrides that hardcode model literals (`:116-120`).
  This is the model-argv fix that makes `MODEL` and `EFFORT` effective rather
  than declared.
- `parley roster show` with `--scope deck|machine`, `--all`, `--json`,
  `--explain AGENT`, and its appearance in `parley --help` and docs
  (`:33-34`, `:92-97`).
- `parley roster set` and `parley roster sync` (command surface at `:92-97`).
  `roster init` becomes a deprecated alias.
- `--scope deck` writes the committed `parley-deck/agents.toml` (`:222-230`).
- `--yes` refused for membership changes; second confirmation required
  (`:100-103`).

This stage ships the effective-model fix, the STATE wiring, the column
contract, and the command surface. It is self-contained: a user can run
`roster show` and get an honest table, and `roster set`/`sync` work. It does
NOT ship rebase, the run snapshot, the §2 protocol change, or the migration.

**Stage 2 — the run snapshot + rebase (must land as one change, after Stage 1,
and the two are atomic per the release gate).**
This stage is gated by `consensus.md:307-315`:
- The immutable run snapshot: at run creation, write a secret-free roster
  snapshot plus `roster_revision` into run state; every later phase uses it
  (`:108-113`).
- `sessions inspect` reports `stale-snapshot` when the deck roster has moved
  since (`:112`).
- `roster sync` rebase semantics: sync removes deck overrides that mask
  machine values, so the deck keeps inheriting (`:182-184`, user direction at
  `:290-294`).
- The acceptance test: create a run, change machine/deck config, continue the
  run, prove adapter/model/effort/autonomous-write args are unchanged
  (`:308-310`).

Stage 1 and Stage 2 could be one PR if the implementer is confident, but the
release gate makes them logically atomic anyway — rebase cannot ship before
the snapshot. Splitting them lets the data contract stabilize (and get
reviewed) before the behavior change lands. If the snapshot implementation
slips, rebase waits — which is exactly R1 and the release gate.

**Stage 3 — the §2 protocol change + generated §2 + skill update (after Stage
1; can be concurrent with Stage 2).**
This is the protocol change the user authorized as a one-off in this idea
(`:320-368`):
- `parley-deck/agents.toml` becomes the deck authority; §2 becomes a generated,
  non-authoritative rendering (`:326-348`).
- The generated §2 is idempotent and preserves human-readable prose (R2,
  `:350-352`).
- Runtime code MUST NOT parse the generated view (`:337`, `resolveRoster`
  changes to read `agents.toml` not §2).
- Every protocol reference calling §2 authoritative, the embedded protocol
  copy, and the skill's bundled snapshot change in the same release
  (`:339-340`).
- The skill invokes `parley roster show` and reproduces its output; it never
  parses §2, TOML, or `agents list` (`:127-131`).
- `meta/protocol-changelog.md` gets the §7-format entry (kimi-1 R2).

This stage does not depend on the run snapshot, so it can proceed in parallel
with Stage 2. It does depend on Stage 1 (the command surface and `agents.toml`
as authority must exist before §2 is generated from it).

**Stage 4 — the fleet migration (after Stages 1, 2, and 3).**
This is the 40-deck migration (`:376-414`):
- `parley roster sync --dry-run` across all 40, full diff to the user.
- Compare-and-swap, file-level backups, verified restore, per-deck
  confirmation, resumability, attended-only, final report.
- Nonterminal legacy runs skipped and reported (release condition 2).
- Foreign decks skipped and reported by name.

This must be last because it requires the STATE wiring (Stage 1), the run
snapshot (Stage 2, for the legacy-run skip), and `agents.toml` as authority
(Stage 3, since the migration writes to it). It is the irreversible fleet op
that triggered the `deliberation` track in the first place.

**What MUST land together:** Stage 1 is internally atomic (STATE wiring + column
contract + command surface + model-argv fix are mutually dependent — the
column contract promises effective values, so the placeholder fix must ship
with it, and STATE must be wired for the migration prerequisite). Stage 2 is
internally atomic by the release gate (snapshot + rebase). Stage 3 is
internally atomic (protocol text + generator + skill + code authority switch
in one release, per `:339-340`). Stage 4 is internally atomic per the migration
contract. The cross-stage dependencies are: Stage 2 requires Stage 1; Stage 3
requires Stage 1; Stage 4 requires Stages 1+2+3.

If the implementer prefers fewer stages, Stages 1+2+3 can collapse into one
large change (the release gates still enforce the internal atomicity), but I
recommend the four-stage split for reviewability — the `deliberation` track
means all non-implementers review, and a single monolithic change spanning
placeholders, snapshot, protocol rewrite, and migration is hard to review
well.

---

#### Verification notes

- `PRIMARY` — `internal/app/roster.go:110`: `active, _, ok :=
  protocol.ReadRosterIDs(root)`. The inactive map is discarded. Confirmed
  revision 2's claim at `:344-345` is accurate. (I own this finding from
  round-2; I do not verdict it — I confirm the consensus cites the correct
  locator.)
- `PRIMARY` — `internal/protocol/roster.go:62-63`: `if strings.Contains(...) {
  inactive[id] = true }`. The parser populates the inactive map. Confirmed
  revision 2's claim at `:345` is accurate.
- `PRIMARY` — `internal/runmanifest/manifest.go:28-46`: `Manifest` carries
  `Participants []string` (`:43`) and no model/effort/adapter/snapshot field.
  Confirmed revision 2's VC-2 measurement at `:201-209` is accurate. (codex-1
  and kimi-1 both verdicted this `CONFIRMED` in revision 1; I concur as a
  non-owner — `SECONDARY` would be the tag if I relied on their verdicts, but
  I read the source directly this session, so this is `PRIMARY`.)
- `PRIMARY` — `internal/app/app.go:1148-1160`: `continueAuto` calls
  `discoverConfigured(ctx, root)` and passes `Agents: discovered` at `:1159`.
  Confirmed revision 2's claim at `:299-302` is accurate — a continuation
  re-discovers config and can silently switch models today.
- `PRIMARY` — `COOPERATION.md:706-708`: §6 rule 3 exception scoped to
  "Never edit another agent's file." Confirmed revision 2's correction at
  `:358-364` is accurate — the exception is not general to §7.
- `PRIMARY` — `COOPERATION.md:717-731`: §7 contains no
  direct-user-instruction exception, only the version-sync carve-out at
  `:726-730`. Confirmed revision 2's correction at `:355-359` is accurate.
- `PRIMARY` — `COOPERATION.md:181`: the `deliberation` classifier triggers
  include "protocol change (§7); ... data migration / irreversible /
  destructive op." Confirmed revision 2's track-upgrade justification at
  `:370-374` is accurate — three triggers fire.
- `PRIMARY` — `COOPERATION.md:134`: "mark its row as inactive (do not delete
  it)." Confirmed the migration's inactive-never-deleted rule is protocol-
  grounded.

No `DISPUTED` claims. No `EXEMPTION-CLAIM UNVERIFIED`. The consensus's
`## Comparison & blind spots` section (`:264-283`) satisfies §15.6(b): it
records the correlated-agreement caveat (four related models converged fast),
states what would have to be true for the agreed position to be wrong (an
adapter resolves its model server-side from a source parley cannot inspect),
and names unmeasured items (SPEED, VC-2's deciding question). §15.6(a) does
not bind: round 1 had substantive disagreement (four different column sets,
different sync semantics, different scope labels, different §2-authority
positions), so the unanimity is not a shared-prior-only convergence.

---

#### Reservation from revision 1: none carried forward

R1, R2, and R3 (all five sub-points) are addressed in binding consensus text.
The reservations I filed in revision 1 were implementation constraints for
FINAL.md; revision 2 elevates them to release gates and prerequisites in the
consensus itself. I have no remaining reservations.

One note for FINAL.md (not a reservation): the consensus does not explicitly
name `meta/protocol-changelog.md` as the changelog target (kimi-1 R2 names it).
The consensus at `:339-340` requires "every other protocol reference calling §2
authoritative ... must change in the same release," which is the substantive
content, but the changelog entry is a §7 format requirement
(`COOPERATION.md:719-724`) that FINAL.md should make explicit. This is a
drafting detail, not a gap in the consensus design.

I accept.
