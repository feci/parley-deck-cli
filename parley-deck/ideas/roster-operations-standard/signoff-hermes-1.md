### hermes-1

Status: 🟡 ACCEPT-WITH-RESERVATIONS
Date: 2026-08-06

**Scope.** I read COOPERATION.md in full (§1–§15), 00-prompt.md, all four round-01
and all four round-02 files, the inbox measurement note (with ADDENDUM), and the
consensus draft. This session I verified source directly (all PRIMARY):
`internal/runmanifest/manifest.go:28-56`, `internal/app/roster.go:100-160,380-389`,
`internal/config/runtime.go:538-617`, `internal/runner/runner.go:1094-1122`,
`internal/agents/discover.go:210-239`. I did not run any live `parley` commands, did
not enumerate the 40 decks myself (SECONDARY: claude-1's inbox measurement), did not
read home-directory config files (`~/.parley/agents.toml`, `~/.hermes/config.yaml`),
and did not verify the EFFORT addendum table firsthand (SECONDARY: claude-1's inbox
ADDENDUM). I did not check whether SPEED has the same declared/effective defect —
the consensus notes this is unmeasured, and I confirm it remains unmeasured by me
(RECALL).

---

**VC-2 verdict (non-owner — §15.1).** claude-1 owns the manifest measurement and is
in the rebase camp, so I am the appropriate non-owner to verdict it. I read
`internal/runmanifest/manifest.go:28-56` directly (PRIMARY).

CONFIRMED: `Manifest` carries `Participants []string` (manifest.go:43) and `Step`
carries `AgentID` + `ArtifactPath` (manifest.go:51-52). There is no model, effort,
adapter, speed, or launch-plan field anywhere in the struct (manifest.go:28-46). The
`New()` constructor (manifest.go:79-120) copies `opts.Participants` as a string slice
and nothing more. A completed run records who participated and which artifacts they
wrote, but not what any agent actually launched. The drafter's measurement is
accurate.

**Does the coupling argument hold?** Yes. codex-1's rebase position rests on the run
snapshot guaranteeing reproducibility — "a synced deck stays current" only works if
you can reconstruct what a past run used. Today you cannot: the manifest has
participant IDs only (PRIMARY, above). Without decision 6's snapshot, rebase would
make a run's roster unreconstructable from the repository alone, because the deck
file no longer pins the values and the manifest never recorded them. The coupling is
real: rebase is only admissible if the immutable snapshot ships with it.

**Is rebase safe given decision 6 is unanimous but unshipped?** The user chose
rebase; decision 6 is unanimous and the consensus correctly reads the practical
effect as rebase + snapshot. But "unshipped" means a design agreement, not
implemented code (PRIMARY: the manifest today has no snapshot fields —
manifest.go:28-46).

> **R1 — rebase/snapshot delivery coupling.** Rebase and the immutable run snapshot
> must ship as one atomic delivery unit. If the snapshot implementation slips past
> the rebase implementation, a window exists where rebase is live but reproducibility
> is not guaranteed — exactly the failure mode the snapshot was designed to close.
> FINAL.md must state this as a hard delivery constraint, not a hope. Reproducibility
> must not depend on an unshipped feature shipping later; it must ship with rebase or
> rebase waits.

I accept rebase because decision 6 is unanimous and will ship in the same change. R1
is the guardrail that makes that acceptance safe.

---

**§7 deviation.** §7 requires protocol changes to run as a separate
`meta-protocol-change-*` idea. The user directed the §2 authority change happen
here. I accept the venue. The consensus logs the deviation in `## User direction`
with the user's verbatim direction and explicitly notes "a signer may still block the
protocol text on its merits." This is the §6 rule 3 direct-user-instruction
exception applied to §7's process requirement, and the logging is sufficient: the
deviation is visible, the user's authority is cited, and the protocol edit still
requires full participant ratification at signoff — which is what this is. The user
authorized the venue, not the text; I assess the text separately below.

**Protocol wording — §2 authority (on its merits).** The consensus (§10 + user
direction) makes `parley-deck/agents.toml` the deck authority with §2 a generated
view. I accept the wording. §2 is the store that drifted nine ways across 40 decks
(SECONDARY: claude-1's measurement; the mechanism — hand-edited prose at fleet scale
— is RECALL from my own round-2 analysis). codex-1 and kimi-1 both reversed
round-1 positions to reach this convergence, which is evidence the change was earned,
not defaulted.

> **R2 — §2 generation idempotency.** The generated §2 view must be idempotent
> (running the generator twice produces byte-identical output) and must preserve the
> human-readable prose format (Agent ID, Workspace dir, Role). The consensus does not
> specify the generation mechanism, and a non-idempotent generator would re-create
> drift under a new name. This is an implementation constraint for FINAL.md, not a
> block.

---

**Mass migration — 40 decks.** The drafter's four constraints (CLI-executed, backed
up, dry-run-all-first, skip-and-report on unclean) are necessary but not sufficient.
Gaps:

1. **Inactive-set wiring is a hard prerequisite.** I verified (PRIMARY:
   `internal/app/roster.go:110`) that `resolveRoster` reads
   `active, _, ok := protocol.ReadRosterIDs(root)` — the inactive map is assigned to
   `_` and discarded. The protocol parser does populate it (SECONDARY: kimi-1,
   `internal/protocol/roster.go:62-64`). The migration plan marks 17 retired
   `antigravity-1` rows (and 3 `gemini-1`, 1 `agy-1`) as `inactive`. But with the
   current code, marking a row inactive is cosmetic — `resolveRoster` throws the
   inactive set away and the row still renders as active. Decision 3's `STATE`
   column + wiring up the inactive set must ship in the same change as the migration,
   or the retired-agent cleanup is a no-op. The four constraints do not mention this
   coupling.

2. **Per-deck attended confirmation, not bulk `--yes`.** Constraint 3 says "full diff
   reported to the user before anything applies." With 40 decks, a single bulk diff
   is enormous and a single bulk `--yes` is exactly the mass mutation where one bad
   deck gets swept through. The `roster_change_policy = "confirm-breaking"` setting
   (SECONDARY: my round-2 concern #4, citing `~/.parley/agents.toml:18`) should gate
   each breaking change. The migration should require per-deck or small-batch
   confirmation, not one global `--yes` across 40.

3. **§14.2 explicit compliance.** The migration is human-attended (the user said
   "Prejdem všetkých 40 deckov" — "I'll go through all 40 decks"), but constraint 1
   says "executed by `parley roster` itself" without stating who triggers it. The
   constraint should state explicitly: this is a human-attended operation, not a
   cron/CI/loop hook (§14.2). An automated loop must not modify the active roster.

4. **Backup strategy for dirty and non-git decks.** Constraint 2 says "every deck is
   backed up" but does not say how. Constraint 4 notes several decks are "other
   projects, years old." Some may not be git repos; some may have uncommitted working
   trees. A `git stash` + commit is insufficient for non-git decks. The backup must
   be a file-level copy (e.g., `cp -a` to a timestamped backup) that does not depend
   on git, and must handle dirty working trees without losing uncommitted work.

5. **Idempotency and recovery.** If the migration crashes on deck 23 of 40, what
   state are decks 1-22 in? `roster sync` is idempotent (consensus decision 5), but
   the migration constraints do not mention resumability. The migration must be
   resumable: a re-run after a crash picks up where it left off, and a deck already
   migrated is a no-op.

> **R3 — migration guardrails.** The four constraints need the five additions above.
   These are implementation guardrails for FINAL.md/IMPLEMENTATION.md, not blocks to
   the consensus design.

---

**VC-1 — SOURCE column.** I proposed SOURCE in round 1 and withdrew it in round 2
(CHANGE 3) in favor of codex-1's `--explain AGENT` + JSON `sources` object. I confirm
I would still exclude it. The argument that defeated it is the one the consensus
cites: a single SOURCE column can only name the winning layer for one field (MODEL),
which silently privileges MODEL's provenance over EFFORT, SPEED, and AUTO — whose
winning layers may differ. Per-field provenance belongs in `--explain`/JSON, not in a
12th column that is honest about one field and silent about three. kimi-1's 12-column
set spends the slot on SOURCE by folding ROUTE into MODEL-COMPANY; the consensus's
11-column set drops both SOURCE and ROUTE, which is the cleaner contract.

**VC-3 — scope labels.** My position is `deck|machine` (round-2 CHANGE 4). `local` is
ambiguous between machine-local and project-local; `deck` unambiguously names the
`parley-deck/` directory. On the write target: `--scope deck` must write the
committed `parley-deck/agents.toml`, not the gitignored `agents.local.toml`. I
verified (PRIMARY: `internal/app/roster.go:383-389`) that `rosterTargetPath` maps the
non-machine scope to `filepath.Join(root, protocol.DeckDir, "agents.toml")` — the
committed file. An invisible gitignored change is how a deck silently diverges from
its own repository; claude-1 raised this and I agree. The consensus records this
explicitly, which is correct.

---

**Reservations summary.**

- R1: rebase + snapshot must ship as one atomic delivery unit; if the snapshot slips,
  rebase is held back.
- R2: the §2 generation mechanism must be idempotent and preserve human-readable
  prose.
- R3: migration constraints need five additions (inactive-set wiring as hard
  prerequisite, per-deck confirmation, §14.2 explicit, backup strategy for
  dirty/non-git decks, idempotency/recovery).

None of these are blocks. The core decisions (1–9) are sound, evidence-backed, and
the rounds showed genuine convergence with position reversals on evidence — I
reversed six positions myself in round 2, each forced by another participant's
evidence. The verdict conflicts are either resolved by user direction (VC-2) or by
substantive argument (VC-1, VC-3). The §15.6 correlated-agreement caveat is
addressed: round 1 had substantive disagreement (four different column sets,
different sync semantics, different scope labels, different §2-authority positions),
so §15.6(a) does not bind; §15.6(b) is satisfied by the consensus's
`## Comparison & blind spots` section, which records the shared-prior caveat and
states what would have to be true for the agreed position to be wrong.

I accept with the three reservations above, which should be addressed in FINAL.md as
implementation constraints.
