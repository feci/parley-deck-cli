### kimi-1

**Verdict: accept with reservations**

**Scope declared.** I read in full: `parley-deck/COOPERATION.md` §§2, 6, 7, 14, 15
(plus the section index); `00-prompt.md`; all four round-01 and all four round-02 files;
`parley-deck/inbox/claude-1-to-all_roster-operations-standard_measured-drift.md` including the
ADDENDUM; and `consensus.md`. Fresh checks executed this session (all `PRIMARY`): full read of
`internal/runmanifest/manifest.go`; read of `internal/runcontrol/runcontrol.go:140-189`; grep of
`.gitignore`; grep of `COOPERATION.md` for user-exception patterns. I did **not**: re-enumerate
the 40 decks, run any `parley` binary, read `~/.parley/*` or any foreign deck, re-read
`discover.go`/`runtime.go`/`runner.go` this session (I rely on round-1/2 citations for those),
inspect resume code, or run any git command. Per §15.1 I issue no verdict on any claim I own
(e.g. my round-1 "`roster` undocumented" and "`resolveRoster` discards the inactive set"
findings).

#### 1. The VC-2 measurement — verified, and rebase is safe only as gated below

**Non-owner verdict on the drafter's measurement: CONFIRMED (`PRIMARY`).** I read
`internal/runmanifest/manifest.go` this session. The `Manifest` struct (lines 28–46; the drafter
wrote "28-45" — the struct body ends at line 46, a locator nit that changes nothing) contains
exactly: `schema_version`, `run_id`, `workspace_root`, `idea_slug`, `task`, `mode`, `transport`,
`status`, `phase`, `idea_status`, `current_round`, `active_steps`, `last_action_at`,
`next_actions`, `participants []string`, `created_at`, `updated_at`. `Step` (lines 48–56) carries
`agent_id` and `artifact_path`. **No model, no effort, no adapter, no invocation, no roster
revision.** The drafter's enumeration is accurate.

**One refinement, stated as evidence, not as refutation (`PRIMARY`).** The manifest is not the
whole run record: the `run.created` event (`internal/runcontrol/runcontrol.go:60`) records
per-discovery `model`, `reasoning`, `speed`, `profile`, `sources` and runtime fields via
`RuntimeEventData` (`runcontrol.go:152-177`). So "a run records participant IDs and nothing else"
is too broad for run state as a whole. But it does not rescue reproducibility: the event stores
`result.Model` / `result.Reasoning` — the declaration-level spec fields, not the materialized
argv — and no adapter-keyed invocation, no effort-as-launched, and no `roster_revision` anywhere
(the identification of `result.Model` as the declared value follows the round-1 code chain
`applyOverride` → `spec.Model`, `runtime.go:594-605` — `SECONDARY` via codex-1's P4/P6, which
terminate in his `PRIMARY` reads). Whether resume re-resolves from current config I did not
check; that part of codex-1's concern 4 stands `UNVERIFIED` by me.

**Position on rebase without a hard snapshot prerequisite.** The user chose plain rebase; the
consensus itself notes the practical effect is rebase *with* the snapshot because decision 6 is
unanimous and in this same document. On the verified facts: under rebase, the deck's committed
files no longer contain the values a past run used, and — per the measurement above — no run
artifact written today can reconstruct them. Reproducibility of every pre-snapshot run then rests
on deck git history plus an unversioned `~/.parley/agents.toml`, which is not a designed property
but an accident. I therefore accept rebase **only under reservation R1**: it is unsafe if the
coupling is allowed to slip, and the FINAL must convert the coupling from prose into a release
gate. The user did not overrule the coupling argument; they declined to make it a separate
prerequisite, which is compatible with gating the release rather than the decision.

#### 2. The §7 deviation — venue accepted, logging necessary but mischaracterized

`PRIMARY`: §7 (`COOPERATION.md:717-731`) requires a `meta-protocol-change-*` idea and names
exactly one carve-out — the version-sync exception at `:726-730`. It contains **no**
direct-user-instruction exception. The protocol does contain user-authorized exception
*precedents elsewhere* — the solo exception (§1, `:84`) and the §6 rule-3 override for editing
another agent's file under direct user instruction (`:708`) — and §13.2 (`:1113`) repeats
"changeable only via a meta-protocol-change idea (§7) with human approval." So when consensus.md
calls this "the protocol's direct-user-instruction exception", it names an exception §7 does not
textually contain; the accurate description is **a user-directed deviation from §7, recorded,
with §7's ratification and audit requirements preserved in full**.

I accept the venue. §7's purpose is auditability and no-silent-edits, not friction for its own
sake; an explicit, verbatim-quoted user decision (`consensus.md:262-273`) with all-participant
ratification at signoff satisfies that purpose. The logging is sufficient only with the
follow-through in **reservation R2** (verbatim §2 replacement text in FINAL.md; a
`meta/protocol-changelog.md` entry in the §7 format naming this idea and the user direction as
authority; an explicit statement whether this is a one-time deviation or §7 is amended to codify
the exception — otherwise the next protocol change hits the same ambiguity).

On the **merits of the protocol wording**: the wording does not exist yet — consensus states only
the direction ("`parley-deck/agents.toml` becomes the deck authority, §2 is generated"). My
acceptance is conditional on the text itself: the generated §2 must be explicitly
non-authoritative, CLI-regenerated, never hand-edited, and **no code path may parse it**
(codex-1's round-2 concern 2 — a generated table that code still reads is a second stale view,
which re-creates the dual-store failure with a generator). Anything beyond the stated scope — in
particular anything touching quorum or signoff rules, an explicit non-goal of this idea
(`00-prompt.md:84`) — voids this signoff.

#### 3. Mass migration — the four constraints are necessary, not sufficient

The fleet figures are `SECONDARY` (facilitator's `PRIMARY` measurement, inbox note incl. the
9-roster table; I did not re-enumerate): 40 decks, 17 with no §2 roster, 17 naming retired
`antigravity-1`, 3 `gemini-1`, 1 `agy-1`. The drafter's four constraints (CLI-executed after the
CLI ships, backed up, dry-run-all-first with full diff, skip-and-report on anything unclean) are
the right spine. What is missing, each item checkable:

- **Locked quorum and live runs.** `PRIMARY` (consensus.md:83-99): decision 5's safety list does
  not restate codex-1's property that sync never edits an open idea's `participants:` or a live
  run snapshot. For a fleet run this must be explicit. Seeding the 17 rosterless decks from the
  machine roster — which contains `opencode-1` (`SECONDARY`, facilitator measurement) — changes
  future default quorum in those decks; the report must enumerate that per deck, not aggregate
  it.
- **Retention rule for retired agents.** `PRIMARY` (`COOPERATION.md:134`): "mark its row as
  inactive (do not delete it)". The 21 retired-agent occurrences must become `inactive`, never
  deleted; the four constraints do not say so.
- **Foreign decks carry foreign protocol copies.** The measurement implies 40 separate
  `COOPERATION.md` files at various protocol versions (`SECONDARY`, facilitator). Writing a
  generated-§2 authority model into a deck whose protocol copy still instructs hand-editing §2
  creates a cross-deck contradiction. Migration of a non-parley-deck-cli deck must be gated on
  that deck's §9.0 protocol sync, or skipped and reported.
- **"Clean mapping" needs a written definition before the dry run**, else skip-and-report is
  discretionary: unparseable §2, deliberate deck pin masking a machine value, rostered adapter
  not installed, and pre-existing `masked-by-env` conditions should each be a named skip class
  with a per-deck, machine-readable report.
- **Post-apply verification and restore.** Each migrated deck is re-resolved and compared to the
  approved diff (the same post-write re-resolve decision 5 gives `set`); "backed up" needs a
  documented restore path, not just a file copy.
- **Git disposition in 39 repositories this project does not own.** The constraints say parley
  executes the writes; they do not say whether changes are committed. Require: working-tree
  changes left uncommitted, or committed only with per-deck user approval — stated in FINAL.md.
- **Fleet form of the membership second confirmation.** `PRIMARY` (consensus.md:93-95): `--yes`
  alone is refused when membership changes. Migration inherently changes membership in ~21 decks;
  define whether the second confirmation is per-deck, per-class, or one aggregate act.

#### 4. VC-1 — the "one column, one field" argument defeats my proposal; I withdraw `SOURCE`

This was my column (round-01, narrowed in round-02 to "winning layer for `MODEL`"). The argument
proves that a row-wide `SOURCE` is incoherent: `MODEL`, `EFFORT`, `SPEED`, and `AUTO` can each
win at a different layer, so one cell silently privileges `MODEL` and misleads about the rest.
That does not logically touch the narrowed MODEL-scoped version — but the narrowed version fails
on its own merits: a frozen, additive-only API column must carry permanent width, `STATUS`
already flags the cases where `MODEL`'s layer is surprising (`model-drift`, `masked-by-env`), the
exact layer is one `--explain AGENT` away, and a header named `SOURCE` cannot carry the scoping —
it will be read as row provenance no matter what the docs say. The two documented reversals
(claude-1, hermes-1) were position changes on the merits, not votes, and mine now joins them the
same way. **VC-1 is resolved by argument, not by the 3-to-1 count, which is the resolution §15.3
requires.** The 11-column contract in decision 3 stands as written.

#### 5. VC-3 — `deck|machine`, and `--scope deck` writes the committed file

My position: `deck|machine`, as I proposed in round-01. Note for the record that VC-3's framing
is round-1-stale: codex-1's round-02 SELF-CORRECTION C1 already reads "`roster set AGENT --scope
deck|machine`" (`PRIMARY`, `round-02/codex-1.md:15-17`), so the label question was effectively
converged before consensus opened. On the file-target question: **yes, unconditionally** —
`--scope deck` must write the committed `parley-deck/agents.toml`, never the gitignored
`agents.local.toml` (`PRIMARY`: `.gitignore:6` ignores `parley-deck/agents.local.toml`; `:7`
ignores `meta/headless-agents.local.json`). A verb named "deck" writing an invisible machine-local
file re-creates the exact silent-divergence disease this idea exists to cure; `agents.local.toml`
survives purely as a manual escape-hatch layer, and when it masks a deck value the row should
report `masked-by-env`, not silently win.

#### 6. Accuracy corrections to consensus.md (positions, not verdicts)

- **VC-2 mislabels my camp.** Consensus lists "hermes-1, kimi-1 — additive, source-aware pin"
  (`consensus.md:167`). My round-02 position change 1 adopted codex-1's rebase, with the `--keep`
  amendment (`PRIMARY` as a document quote, `round-02/kimi-1.md:22-29`). Post-round-2 the split
  was 3-to-1 for rebase, not 2-to-2. The user's direction closes VC-2 either way, but the
  conflict record §15.3 requires should quote positions as they actually stood.
- **`--keep` disappeared silently.** Decision 5's command line has no `--keep`
  (`consensus.md:89`). Codex-1's round-2 answer — preview labels each overwritten pin; the user
  aborts or re-applies with `roster set --scope deck` — is acceptable for single-deck use but
  impractical across 40 decks. **Reservation R4**: either `--keep <agent>.<field>` ships, or the
  fleet report must enumerate every deliberate pin the rebase removes, per deck, so re-application
  is a checklist rather than an archaeological dig.

#### Reservations (consolidated; each is a checkable FINAL.md/implementation requirement)

- **R1 — rebase coupling is a release gate.** Decision 6 (immutable run snapshot, with resume
  consuming the snapshot, not merely comparing it) ships in the same release as `roster sync`;
  FINAL.md states this as a gate, not an ordering hope. A destructive sync (removes a masking
  override or deactivates a row) against a deck with runs created before snapshot support
  requires the breaking-change second confirmation and prints the reconstructability warning.
- **R2 — §7 follow-through.** FINAL.md quotes the verbatim §2 replacement text (scope-limited to
  the authority move + generated non-authoritative view; no code path parses §2 afterward);
  `meta/protocol-changelog.md` gets the §7-format entry naming this idea and the 2026-08-06 user
  direction as authority; FINAL.md states whether §7 is amended to codify a user-directed venue
  exception or this stands as a recorded one-time deviation.
- **R3 — migration constraints extended** per §3 above: locked `participants:`/live runs
  untouched; retired agents marked `inactive`, never deleted; foreign-deck protocol-version gate;
  written skip-class definitions with per-deck machine-readable report; post-apply re-resolution
  against the approved diff; documented restore path; stated git disposition; defined fleet form
  of the membership second confirmation.
- **R4 — deliberate pins survive discoverably.** `--keep` ships, or the per-deck enumeration of
  removed pins is mandatory in the migration report.

None of the reservations contradicts the consensus text; all of them bind what FINAL.md and the
implementation must make explicit. With them recorded, I sign.
