### kimi-1 — revision 2

**Verdict: accept with reservations**

codex-1's four required changes are all met. hermes-1's R1, R2 and R3 (all five sub-points) are
addressed. Of my own revision-1 reservations, R1 and R2 are answered at the level a consensus can
answer them; parts of R3 and all of R4 are not carried into revision 2's operative text and are
restated below as binding FINAL.md/implementation requirements. None of the residuals contradicts
the consensus text, so they are reservations, not a block.

**Scope declared.** Read in full this session: `parley-deck/COOPERATION.md` §2, §4.0, §5, §6, §7,
§13.2, §14, §15 (all `PRIMARY`, locators quoted below); the full revision-2 `consensus.md` (938
lines) including all four embedded revision-1 signoffs; my own `signoff-kimi-1.md`. Fresh checks
executed this session (all `PRIMARY`): read of `internal/app/app.go:1143-1160` (`continueAuto`);
read of `internal/app/roster.go:100-130`; a programmatic byte-for-byte diff of each embedded
revision-1 signoff block in `consensus.md` against its standalone `signoff-<agent>.md` file; a grep
of `consensus.md` for `keep|deliberate pin|foreign|protocol sync|protocol-version`. I did **not**:
re-enumerate the 40 decks, run any `parley` binary, read `~/.parley/*` or any foreign deck, re-read
`manifest.go`/`runcontrol.go` this session (my revision-1 `PRIMARY` reads of
`internal/runmanifest/manifest.go:28-56` and `internal/runcontrol/runcontrol.go:140-189` stand),
or read `signoff2-codex-1.md` — it exists on disk but I wrote this before reading it, to keep this
signoff independent. I ran no git command and wrote only this file. Per §15.1 I issue no verdict on
claims I own: the `printUsage`/docs omission (`consensus.md:33-37`), the `DISPLAY-NAME`
contradiction (`consensus.md:69-71`), and the discarded inactive set at `internal/app/roster.go:110`
— the last is cited below only as `SECONDARY` via hermes-1's `PRIMARY` verification
(`consensus.md:506-507,588-596`).

**Standing non-owner verdicts, unchanged from revision 1.** The drafter's VC-2 manifest measurement:
`CONFIRMED` (my revision-1 `PRIMARY`, full read of `manifest.go`). New this revision: codex-1's
continuation-mechanism claim, which I left partially `UNVERIFIED` in revision 1, is now `CONFIRMED`
(`PRIMARY`): `continueAuto` calls `discoverConfigured(ctx, root)` and passes the freshly discovered
values as `Agents: discovered` (`internal/app/app.go:1148-1159`), with no snapshot lookup anywhere
in the path. Combined with the manifest carrying no invocation plan, a continuation after a machine
config change re-resolves from current config. That is the load-bearing fact for release
condition 1, and it is now verified by a non-owner.

#### 1. codex-1's four required changes

**Required change 1 — snapshot acceptance gate + legacy-run skip, or a plain "unsafe" record: MET.**
All three binding release conditions are adopted (`consensus.md:305-315`): the change exposing
rebase must persist *and consume* the immutable effective snapshot, with a named acceptance test
(create run, change machine/deck config, continue, prove adapter/model/effort/autonomous-write args
unchanged); fleet migration must skip and report every nonterminal legacy run lacking the snapshot,
with `participants:` lists and run artifacts never rewritten to manufacture one; and any relaxation
must be stated plainly as **unsafe for pre-snapshot resumable runs**, with "decision 6 is unanimous"
never offered as present protection. The revision also records *why* this is binding rather than
advisory — three participants converged on it independently (`consensus.md:317-318`).

**Required change 2 — §7 wording corrected, track upgraded, §2 authority spec completed: MET at
consensus level.** The frontmatter now reads `track: deliberation` (`consensus.md:5`), and the
upgrade is justified against the classifier: three triggers fire (protocol change §7, data
migration, irreversible fleet op) plus the fail-closed rule (`consensus.md:370-374`; classifier
verified `PRIMARY` at `COOPERATION.md:179-190`). The §7 mis-citation is corrected in full
(`consensus.md:354-368`): the only textual exception is §6 rule 3, scoped to editing another
agent's file (verified `PRIMARY`, `COOPERATION.md:708`; §7 itself, `:717-731`, contains no
user-instruction exception), and the deviation now rests on the user's explicit one-off
authorization alone, creating **no precedent** (`:362-364`). The §2 authority spec is completed
(`consensus.md:326-340`): workspace dir, role, host handle, active/inactive history and ordering
must each get a canonical source and migration path; the generated §2 is non-authoritative, a
rendering not a store; runtime code must not parse it; and every protocol reference, the embedded
protocol copy and the skill's bundled snapshot must change in the same release (`:339-340`).
One honest scoping note: the consensus now fully enumerates *what the protocol text must define*;
the concrete per-item canonical mapping (which file/key holds workspace dir, role, host handle) is
still owed in FINAL.md's verbatim protocol text. That is consistent with the revision's own rule
that a signer may still block the protocol text on its merits (`:366-368`) — the venue question is
closed, the wording ratification is not.

**Required change 3 — migration contract additions: MET.** The revision-2 migration contract
(`consensus.md:390-414`) contains every element codex-1 listed: machine-readable inventory of
roots, frozen source roster revision, pre-migration hashes, protocol/schema version, worktree
state, dry-run disposition (`:393-395`); compare-and-swap between dry-run and apply with explicit
batch approval and the membership second confirmation (`:396-398`); a precise six-class "unclean"
definition — dirty worktree, parse/validation error, unsupported legacy layout, path/symlink
ambiguity, concurrent modification, nonterminal pre-snapshot run — skipped, never guessed
(`:399-401`); backups with recorded location and hashes, verified restore, atomic per-deck writes,
post-write `roster show`/schema validation, automatic rollback on validation failure (`:402-405`);
and a final machine-readable report marking every deck `applied` / `skipped` /
`failed-and-restored` / `unchanged` with before/after hashes and the backup reference, with no
commit, push, or edit to locked `participants:` lists or historical artifacts (`:411-414`).

**Required change 4 — VC-1 closed by engagement, VC-3 ratified as `deck|machine` writing the
committed file: MET.** VC-1 is closed by argument, not count (`consensus.md:155-165`), with the
three position reversals documented and the original positions preserved (`:167-179`). VC-3 is
closed unanimous (`:222-230`): `deck|machine`, `--scope deck` writes the committed
`parley-deck/agents.toml` and never the gitignored `agents.local.toml`, grounded in source
(`docs/agent-runtime-configuration.md:7-15`, `internal/config/runtime.go:134-151`), and the
revision records that the conflict it closed was already stale — a drafter error, corrected
(`:228-230`).

#### 2. hermes-1's R1, R2, R3

- **R1 (rebase + snapshot as one atomic delivery unit): addressed.** Release condition 1
  (`consensus.md:306-310`) makes it a hard gate — "Rebase must not ship first" — and the text
  records that hermes-1 filed it as R1 alongside two other independent convergences (`:317-318`).
- **R2 (generated §2 idempotent, prose shape preserved): addressed.** Byte-identical output on a
  second run and preserved human-readable prose, with the correct rationale that a non-idempotent
  generator recreates drift under a new name (`consensus.md:350-352`).
- **R3.1 (inactive-set wiring a hard prerequisite): addressed, and independently confirmed.**
  `resolveRoster` discards the inactive map into `_` (`SECONDARY` via hermes-1's `PRIMARY` read of
  `internal/app/roster.go:110`; the claim is mine from round 1, so no verdict from me per §15.1).
  The consensus states marking rows inactive is cosmetic today and that the `STATE` column plus the
  inactive-set wiring must ship in the same change as the migration (`consensus.md:342-348`).
- **R3.2 (per-deck attended confirmation, not bulk `--yes`): addressed** (`consensus.md:406-407`),
  honoring `roster_change_policy = "confirm-breaking"`.
- **R3.3 (§14.2 explicit): addressed** — human-attended only, never from a loop, cron or CI hook
  (`consensus.md:410`; §14.2 verified `PRIMARY`, `COOPERATION.md:1153-1160`, "Modify the active
  roster" at `:1159`).
- **R3.4 (file-level backups for dirty/non-git decks): addressed** — file-level copies, not git
  operations (`consensus.md:404-405`). The facilitator note (`:416-419`) is properly candid that
  the pre-taken backup satisfies *existence*, not the *verified restore* requirement.
- **R3.5 (resumability): addressed** — crash on deck 23 leaves 1-22 known, re-run resumes,
  already-migrated deck is a no-op (`consensus.md:408-409`).

#### 3. My revision-1 reservations, and VC-1

**VC-1 closure is recorded as I argued it.** The closure (`consensus.md:155-165`) matches my
revision-1 §4: the row-wide version is incoherent because `MODEL`, `EFFORT`, `SPEED` and `AUTO` can
each win at a different layer; the narrowed `MODEL`-scoped version fails on its own merits —
permanent width in a frozen additive-only API, `STATUS` already flagging the surprising cases, and
a header named `SOURCE` being read as row provenance regardless of documentation. Three of four
reached exclusion by reversing their own earlier position, making it resolution-by-argument under
§15.3 rather than 3-to-1 attrition. The eleven-column contract stands. No residual.

**R1 (rebase coupling as release gate): substantively answered.** Conditions 1-3
(`consensus.md:305-315`) are exactly the conversion of the coupling from prose into a release gate
I asked for, including the consume-don't-compare requirement and the acceptance test. One
operational detail from my R1 is not restated: a destructive single-deck `roster sync`/`set`
against a deck with pre-snapshot runs should print the reconstructability warning and require the
breaking-change second confirmation. Decision 5's safety list (`consensus.md:99-103`) implies the
second confirmation for membership changes but does not name the warning. Carried as K3 below.

**R2 (§7 follow-through): the mischaracterization is corrected; the follow-through remains owed.**
Revision 2 answers the one-off-vs-amendment question the right way — one-off, no precedent, a
future protocol change still requires its own meta idea (`consensus.md:362-364`). Still owed in
FINAL.md: the verbatim §2 replacement text and the `meta/protocol-changelog.md` entry in the §7
format naming this idea and the 2026-08-06 user direction as authority. Carried as K4, with the
same scope limit as revision 1: anything beyond the authority move — in particular anything
touching quorum or signoff rules — voids this signoff.

**R3 (migration constraints): mostly answered, one sub-point not adopted.** Locked
`participants:`/live-run artifacts untouched (`:413-414`), retired rows migrated to `inactive`
rather than deleted (the migration target state at `:345-346`, consistent with
`COOPERATION.md:134`), written skip classes (`:399-401`), post-apply re-resolution and verified
restore (`:402-405`), stated git disposition (`:413-414`), and the fleet form of the membership
second confirmation (`:396-398,406-407`) are all present. **Not adopted: the foreign-deck
protocol-version gate** (my revision-1 signoff, `consensus.md:775-779`). A deck whose own
`COOPERATION.md` copy still instructs hand-editing §2 is not a named skip class; the contract
names foreign projects in the dry-run diff (`:387-388`) and skips "unsupported legacy layout"
(`:399-401`), but a stale-*protocol* deck can present a perfectly clean roster layout while its
protocol text contradicts the authority model being written into it. `PRIMARY`: a grep of
`consensus.md` shows `protocol sync`/`protocol-version` gating appears only inside my embedded
revision-1 block, nowhere in the operative revision-2 text. Carried as K2.

**R4 (deliberate pins survive discoverably): not addressed.** Neither `--keep` nor a mandatory
per-deck enumeration of removed pins appears anywhere in revision 2's operative text (`PRIMARY`:
grep for `keep|deliberate pin`; both terms occur only inside my embedded revision-1 block and in
drafter position change 4, `consensus.md:250`, which is history, not a requirement). Under rebase,
sync *removes* deck overrides that mask machine values; across 40 decks, some of those overrides
are deliberate pins, and codex-1's round-2 answer (preview labels each; user re-applies with
`roster set`) is impractical at fleet scale without an enumeration. Carried as K1.

**Record accuracy (non-binding, for FINAL.md).** Two items from my revision-1 §6 were not taken:
the VC-2 header still says "OPEN" (`consensus.md:181`) while the body closes it by user direction
(`:294`), and the preserved framing still camps me as "additive, source-aware pin" (`:186`)
although `round-02/kimi-1.md:22-29` adopted codex-1's rebase with the `--keep` amendment, making
the post-round-2 split 3-to-1. The user's direction closes VC-2 either way, but §15.3's conflict
record should quote positions as they actually stood when consensus opened.

**Carried reservations (each a checkable FINAL.md/implementation requirement):**

- **K1 — deliberate pins survive discoverably.** `--keep <agent>.<field>` ships, or the migration
  report enumerates every deliberate pin the rebase removes, per deck, so re-application is a
  checklist rather than archaeology.
- **K2 — foreign-deck protocol-version gate.** Migrating a deck whose `COOPERATION.md` copy still
  instructs hand-editing §2 is gated on that deck's §9.0 protocol sync, or the deck is a named skip
  class in the machine-readable report — not absorbed into "unsupported legacy layout".
- **K3 — destructive sync against pre-snapshot decks warns.** A `roster sync`/`set` that removes a
  masking override or deactivates a row in a deck with runs created before snapshot support
  requires the breaking-change second confirmation and prints the reconstructability warning;
  FINAL.md states this, not merely implies it via decision 5.
- **K4 — §7 follow-through.** FINAL.md quotes the verbatim §2 replacement text (scope-limited to
  the authority move plus the generated non-authoritative view; no code path parses §2 afterward)
  and includes the `meta/protocol-changelog.md` entry in the §7 format (`COOPERATION.md:719-724`)
  naming this idea and the user direction as authority.

#### 4. One change or staged

**Staged.** The work list is ten workstreams across CLI, protocol text, skill and 40 foreign
repositories; a single change would satisfy every coupling trivially but produce an unreviewable
diff and bury the release gates the block just fought to establish. Staging keeps each gate
enforceable — provided the couplings already named in the consensus are treated as atomic units,
not ordering hopes. Four stages, within this one idea's implementation:

- **Stage 1 — data truth and the read path.** `{model}`/`{effort}` placeholders + resolver + the
  legacy `headless_args` normalizer; the versioned `modelmeta` registry with golden tests; `STATE`
  wiring so `resolveRoster` consumes the inactive set; the 11-column contract + JSON schema +
  golden tests; `roster show`/`--json`/`--explain`/`--all`; help and docs listing. These land
  together because the contract must not freeze "effective" before an effective value exists
  (drafter position change 3, `consensus.md:249`), and the drift symptoms are one defect at three
  altitudes, to be fixed as one change (`:276-280`).
- **Stage 2 — the write path and the snapshot, atomically.** The immutable run snapshot (persist
  *and consume*, with the acceptance test) and `roster set` + `roster sync` exposing rebase. The
  rebase gate already forces this atomicity: "the change that exposes rebase MUST also persist and
  consume the immutable effective snapshot… Rebase must not ship first" (`consensus.md:306-310`).
  If staging pressure ever splits them, the snapshot ships first and sync waits — never the reverse.
- **Stage 3 — the authority flip, one release.** `agents.toml` becomes the deck authority; the
  generated, idempotent §2 view; runtime stops parsing §2; every protocol reference that calls §2
  authoritative, the embedded protocol copy and the skill's bundled snapshot change in the same
  release (`consensus.md:339-340`); the skill's roster section becomes the pointer plus the three
  verbs. Splitting this stage re-creates the dual-store failure with a generator.
- **Stage 4 — the fleet migration.** Attended execution of the revision-2 migration contract,
  last. It depends on Stage 1 (`STATE` wiring — R3.1, `consensus.md:342-348`), Stage 2 (the
  snapshot, so nonterminal pre-snapshot runs are a defined skip class) and Stage 3 (the authority
  model being migrated to), and on K1-K2 above being implemented in the migration command.

Must land together, as hard constraints: **(a)** snapshot persist+consume with `roster sync`
(release condition 1); **(b)** `STATE` wiring with the migration (R3.1); **(c)** protocol change +
generated §2 + skill + embedded copies in one release (`consensus.md:339-340`). Stages 1-2 may
overlap or swap order only if (a) is preserved.

#### 5. Verification-integrity notes

- All four revision-1 signoffs are embedded byte-for-byte verbatim: `CONFIRMED` (`PRIMARY` —
  programmatic diff of each embedded block against `signoff-claude-1.md`, `signoff-codex-1.md`,
  `signoff-hermes-1.md`, `signoff-kimi-1.md`; all four identical).
- The revision-2 changelog's count of eleven drafter position changes, seven forced at signoff
  (`consensus.md:241-262`), matches the table as written: `CONFIRMED` (`PRIMARY`, read of the
  table against the cited sources; the seven signoff-forced entries are items 5-11).
- The fleet figures (40 decks, 17 rosterless, 17 retired-`antigravity-1`) remain `SECONDARY`
  (claude-1's `PRIMARY` inbox measurement, unreproduced by me). Whether `SPEED` shares the
  declared/effective defect remains unmeasured by anyone, me included (`RECALL` — the consensus
  correctly records it as unmeasured, `consensus.md:282-283`).
- §15.6: round 1 had substantive disagreement and every participant reversed positions on evidence,
  so the correlated-agreement clause is satisfied as recorded in `## Comparison & blind spots`
  (`consensus.md:264-283`).

With K1-K4 recorded, I sign revision 2.
