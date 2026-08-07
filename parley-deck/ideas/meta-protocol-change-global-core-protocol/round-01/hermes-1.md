---
agent: hermes-1
idea: meta-protocol-change-global-core-protocol
round: 1
date: 2026-08-07
---

## Summary

The proposal generalizes what roster-operations-standard already did to §2: stop
treating the per-deck `COOPERATION.md` copy as the store, and derive it from a single
core in `~/.parley`. The mechanism already half-exists — `syncConsumerProtocol`
(`internal/app/preflight.go:488`) rewrites the deck copy from the skill's packaged
protocol body, preserving the five identity zones via `mergePreservingZones`
(`preflight.go:520`). This idea makes that the primary model rather than a repair:
the core is the authority, the deck keeps a *generated* committed copy, and a small
overlay carries the near-zero genuinely-local content.

My central thesis on the four weak points: (1) prevention of an agent editing
`~/.parley/COOPERATION.md` is NOT achievable — the agent has the user's privileges —
so the honest design is detection + attribution plus a read-only speed bump; (2) the
override/extension machinery should be the SMALLEST thing that satisfies constraint 3,
because the evidence is 1-in-36, and a rich DSL is unjustified; (3) overlay drift is
contained by making the overlay versioned, additive, and compatibility-checked on every
core bump — the same shape as `RosterRevisionOf` (`internal/runmanifest/manifest.go:83`);
(4) the per-machine core is made safe by deck-level version pinning in `meta/version.json`
plus a render-time refusal to commit a machine-local core into a shared file — the exact
shape of `adoptInherited` in `internal/app/roster_render.go:40`.

## Proposed approach

### Resolution model — section-level, generated view

Override and extend operate at the **Markdown section** level (`## N. Title` headings),
because that is already the unit the codebase parses: `mergePreservingZones` splits at
`## 3.` (`preflight.go:538`), `replaceRosterSection` finds `## 2. Active agents (roster)`
(`roster_render.go:158`), and the drift guard anchors on `rosterSectionLine`
(`drift_test.go:27`). Introducing named-blocks or a key-value DSL would be a new parsing
surface for zero evidence-driven gain.

The effective protocol is materialized by a new `parley protocol render` command,
analogous to `parley roster render`:

  core (`~/.parley/COOPERATION.md`)  +  deck overlay (`parley-deck/protocol-overlay.md`)
  ─────────────────────────────────────────────────────────────────────────────────────
  → generated `parley-deck/COOPERATION.md` (committed, the file agents read)

Precedence: the overlay's override entries replace named core sections by heading match;
the overlay's extension entries append after the last core section. The generated file
carries a `**Protocol synced:** <date> — parley-deck-core <version> (<sha256-prefix>)`
provenance line — the same shape `refreshProtocolSyncedLine` already writes
(`preflight.go:582`). The generated file is what every agent reads at session start (§9
step 1); the core and overlay are build inputs, not session reading material.

The overlay file format — minimal, markdown-native:

    ---
    coreVersion: "2.5.1"
    coreSha256: "<16-char prefix>"
    ---

    ## override §0. Choose the transport
    (replacement body for §0 — the identity zone already permits this)

    ## extend: Project-specific packaged-reference drift
    (new section body, appended after the last core section)

`## override §N` replaces the matching core section. `## extend: <title>` appends a new
titled section. No finer granularity: sub-section override is not offered. If a deck
needs to override a sub-section, it overrides the whole section (rare, and reviewable).

### What is sealed vs open

SEALED (a deck may never override, and the render refuses to place an override entry
targeting them):

- §1 (scope, purpose, non-solo execution requirement) — the definitional invariant.
- §4 including §4.0/§4.0.1 (phases, track classifier, LE rules) — the process engine.
- §5 (quorum and async participation) — the signoff math.
- §6 rules 1-5 (conflict-avoidance mechanics, minus the existing §6.6 working-language
  override note which is already an acknowledged project choice).
- §7 (changing the protocol) — a deck cannot override the change mechanism itself.
- §15 (verification integrity) — the epistemic engine.

OPEN (override permitted) — the five identity zones the drift guard already allowlists
(`drift_test.go:24-30`): `**Workspace:**`, `**Created:**`, `**Transport:**`, the
`**Protocol synced:**` stamp, and the §2 roster table body + host-handle table body.
These are already per-deck by construction; "override" here just means the generated
copy carries the deck's values, which is what `mergePreservingZones` already does.

OPEN (extend permitted) — a deck may append project-specific sections after the last
core section. This is the one capability that is genuinely new. The evidence (1 real
local section in 36 decks, and it was governance-about-sync) says this slot will almost
never be used, but constraint 3 requires it exist, and a declared, reviewable extension
file is better than silent in-place edits to the core copy.

Trade-off stated: a total ban on override+extension (the 1-in-36 evidence would support
it) would violate constraint 3. A permissive default (any section overridable) would
re-create drift under a new name. The sealed/open split above is the middle: the process
engine is immovable, identity zones are per-deck by design, and extension is a single
declared append-only slot.

### Enforcement of constraint 4 — detection + attribution, not prevention

I address weak point 1 directly: an agent runs with the user's own filesystem
privileges. `~/.parley/COOPERATION.md` is a regular file. Any agent that can run
`os.WriteFile` (or `cat >`) can overwrite it. File-mode prevention (chmod 0444) is a
speed bump — the agent can chmod it back, because it has the user's uid. `chattr +i`
(immutable) on Linux is stronger but the agent can `chattr -i` it. On macOS there is no
equivalent of `chattr +i` for a regular user file. So **prevention against a privileged
actor is not achievable; the honest design is detection + attribution.**

The mechanism, layered:

1. **Read-only file mode (speed bump, not wall).** The CLI writes the core with mode
   0444. A naive `os.WriteFile` fails with EPERM. An agent that blindly writes via the
   CLI's own commands is stopped. An agent that explicitly `os.Chmod` first is not — but
   that is a detectable, attributable action.

2. **Content-hash verification (detection).** `~/.parley/version.json` records
   `coreSha256` for the installed core version. `parley core verify` (and preflight's
   freshness path) recomputes sha256 of the on-disk core and compares. A mismatch is
   reported as `drifted-core` — the same shape as `stale-snapshot` for the roster
   (`internal/runmanifest/manifest.go:56-58`). This catches any edit, by agent or
   by accident, regardless of how it was made.

3. **No CLI write path for agents (attack-surface reduction).** The CLI exposes
   `parley core upgrade` (writes a new ratified version, requires `--yes`) and
   `parley core verify` (read-only). There is no `parley core edit`. The only legitimate
   write path is a version upgrade; the user edits the file himself with his editor,
   which is constraint 4's allowance ("The user may edit it himself").

4. **Attribution via version control (optional, recommended).** If `~/.parley` is a git
   repo (the user's choice), `git diff` on the core attributes any change. The CLI does
   not require this, but `parley core init` should recommend it.

What defeats it: an agent that (a) chmods the file, (b) writes arbitrary content, and
(c) updates `~/.parley/version.json`'s `coreSha256` to match its edit, all without the
user noticing, has fully escaped. That is the same trust boundary as "the agent can edit
any file the user can edit" — it is fundamental, not fixable by this protocol. The design
makes it **loud** (hash mismatch on next preflight, `drifted-core` status) and
**attributable** (git diff if versioned), not impossible.

### Version pinning mechanics — same shape as the roster snapshot

Constraint 1 (an open idea runs to completion under the protocol version it started with)
is the exact problem `RosterSnapshot` / `RosterRevision` solved for the roster
(`internal/runmanifest/manifest.go:49-58`). The parallel structure:

| roster (shipped) | protocol (this idea) |
|---|---|
| `RosterSnapshot []RosterSnapshotEntry` | `ProtocolSnapshot` — version string + sha256 of the effective protocol |
| `RosterRevision string` (content hash) | `ProtocolRevision string` (content hash) |
| `stale-snapshot` when deck roster moved | `stale-protocol` when deck's effective protocol moved |
| frozen at run creation, every later phase uses it | frozen at Phase 0, every later phase uses it |

Concretely: `runmanifest.Manifest` gains `ProtocolVersion string`,
`ProtocolSha256 string`, `ProtocolRevision string` fields (all `omitempty` for backward
compatibility, same as `RosterSnapshot` is today). At run creation, the CLI computes the
effective protocol (core + overlay, rendered) and writes its hash. Every later phase of
the run reads from the pinned version, never a fresh resolve — exactly the
"`continueAuto` re-discovers config" hazard the roster snapshot fixed
(`internal/app/app.go:1148-1160`, per roster-operations-standard G1).

The deck's `meta/version.json` already carries `deckVersion` and `protocolSha256`
(`internal/app/preflight.go:116-119`). This idea extends it: `coreVersion` (the
`~/.parley` core the deck is pinned to) and `coreSha256`. The deck's pin is the source of
truth for what version it runs; the machine's installed core may be newer, but the deck
does not auto-advance until it bumps its pin (additive) or the user confirms (breaking).

### Compatibility checking

What is compared: the overlay's `## override §N` targets against the new core's section
list. Specifically: (a) does every overridden section still exist in the new core? (b)
has an overridden section's heading text changed? (c) has a sealed section been moved or
merged?

When it runs: in `classifyAndSyncFreshness` (`preflight.go:348`), after the bump is
classified. Today the additive consumer path auto-syncs and the breaking path pauses
(`preflight.go:425-454`). This idea adds a third outcome: **the overlay is incompatible
with the new core** → the deck stays on its pinned core version and is reported
`overlay-incompatible`, regardless of whether the bump is additive or breaking. The user
must run `parley protocol migrate` (which diffs the overlay against the new core and
either auto-updates the `## override §N` targets if the section was renumbered but
unchanged, or reports the conflict for human resolution).

Outcome on incompatibility: **block + quarantine, never auto-migrate protocol text.**
Auto-migration of protocol prose is unsafe for the same reason §7 requires a meta idea:
protocol text is the rules, and silent rewriting of rules is the failure mode this whole
idea exists to prevent. The deck stays pinned to its old core; the new core is available
but not applied until the overlay is reconciled. This mirrors the `breaking-freshness`
gate (`preflight.go:66`).

### Does the deck keep a committed COOPERATION.md?

YES — a generated, committed copy. Reasons:

1. **Self-contained deck.** An agent or human cloning the repo gets the protocol that
   governs the deck without needing `~/.parley` or the CLI installed. This is the same
   argument that kept a generated §2 in `COOPERATION.md` rather than only in
   `agents.toml`.
2. **Historical preservation.** Git history of the generated file records exactly what
   protocol governed each past idea — the audit trail §1 mandates. A no-copy model
   loses this.
3. **The drift guard needs a file.** `TestEmbeddedDefaultMatchesLiveDeck`
   (`drift_test.go:46`) compares the embedded default with the live deck. Without a
   committed deck copy, there is nothing to compare and drift detection loses its
   anchor.

The generated copy is NOT hand-editable in spirit: `parley protocol render` regenerates
it, and hand-edits to non-identity zones are overwritten on the next render — the same
contract `roster render` already enforces for §2 (`roster_render.go:14-22`). Drift is
not eliminated by removing the copy; it is eliminated by making the copy generated and
checked.

### Migration

The 2026-08-06 sync already proved 35 of 36 decks are pure-core + identity zones
(00-prompt.md:28-31). Migration is therefore near-trivial for most decks:

1. `parley core init` writes `~/.parley/COOPERATION.md` from the installed skill's
   packaged protocol (`skillStatusPackagedProtocol`, `preflight.go:823`) or the embedded
   default (`internal/protocol/workspace.go:22`). Records `coreSha256` in
   `~/.parley/version.json`.
2. For each deck, `parley protocol migrate --dry-run` diffs the deck's current
   `COOPERATION.md` against the core. The five identity zones are extracted into the
   generated identity block (they stay where they are). Any remaining diff — the
   near-zero local content — becomes an `## extend` entry in `protocol-overlay.md`.
3. `parley protocol render` regenerates `COOPERATION.md` from core + overlay. The result
   should be byte-identical to the pre-migration file for 35/36 decks (modulo the
   provenance line), proving no-loss. The one deck with a real local section
   (`librade-algoTrader`, preserved in
   `.parley-protosync-backup-2026-08-06/librade-algoTrader/`) gets that section as an
   `## extend` entry.
4. Decks that are read-only, unreachable, or whose `COOPERATION.md` does not parse
   cleanly against the core are **skipped and reported** — the same `unclean` → skip
   contract from roster-operations-standard's migration (consensus.md:453-456).

The migration is an attended, per-deck-or-small-batch operation, never a bulk `--yes`
across 36 decks — directly inheriting the roster migration constraints
(consensus.md:444-468). It is human-attended only, never from a loop/cron (§14.2).

### §7 impact

§7 itself does not need to change in substance: a meta-protocol-change idea still
produces a new protocol version, and the drafting agent still logs the change in
`meta/protocol-changelog.md`. What changes is the **unit of application**:

- Today: the idea's FINAL.md supersedes `COOPERATION.md`, and the drafting agent updates
  it in-place. With 36 flat copies, "in-place" means 36 files.
- After: the idea's FINAL.md produces a new **core version** that ships via the skill/CLI
  (`parley core upgrade`). Each deck picks it up on its next preflight: additive →
  auto-sync of the generated copy (the §7 carve-out, `COOPERATION.md:754-758`, still
  holds), breaking → pause for user confirmation.

The one textual addition to §7: note that the core version is the unit of protocol
change, and that a deck with an overlay MUST pass compatibility checking before the new
core applies. The carve-out ("a version sync is not a protocol change") is unchanged —
syncing a ratified core version into a deck is maintenance, not a meta idea.

## Concerns / open questions

1. **Overlay format — TOML vs markdown.** I propose markdown-with-frontmatter for the
   overlay because the content is markdown sections and the render already works in
   heading units. A TOML/YAML overlay (`[override."§0"]`, `[[extend]]`) would be more
   machine-parseable but introduces a second parsing surface. The trade-off is
   parseability vs. consistency with the existing heading-anchored code. I lean markdown;
   this is a round-2 question.

2. **Is `parley protocol render` idempotent?** It must be — two renders produce
   byte-identical output — or it recreates drift under a new name. This is
   roster-operations-standard G4 (`consensus.md:104`) applied to the protocol. The
   ordering of extension sections must be deterministic (file order, stable). I state
   this as a requirement, not yet a mechanism.

3. **The embedded default and the skill's packaged copy.** Today there are three
   protocol surfaces: the embedded default (`internal/protocol/defaults/COOPERATION.md`,
   `//go:embed`), the live deck, and the skill's bundled snapshot
   (`skills/parley-deck/references/COOPERATION.md`). The drift guard
   (`drift_test.go:46`) enforces embedded == live (modulo zones). Where does
   `~/.parley/COOPERATION.md` (the core) fit? My reading: the core is a *fourth*
   surface, initialized from the packaged skill protocol on `parley core init`, and
   `parley core upgrade` updates it from the skill's packaged body on version bumps. The
   drift guard should then compare embedded == core (the core is the user's installed
   copy of the ratified protocol) and core == deck-generated (modulo zones + overlay).
   This needs explicit test design — a round-2 question.

4. **What happens to `protocolRole: source`?** Today a `source` deck (this repo) is the
   protocol's upstream and preflight never auto-writes it (`preflight.go:389-395`). With
   a global core, the source repo's `~/.parley/COOPERATION.md` IS the core, and the deck
   copy is generated from it. The `source` role should mean "the core is authored here"
   — the user edits the core directly and `parley protocol render` generates the deck
   copy from it. This inverts the current flow and needs explicit handling.

## Risks

1. **The overlay becomes the new drift surface (weak point 3).** If 36 decks each grow
   an overlay with `## extend` sections, those sections drift independently. Mitigation:
   the overlay is versioned (`coreVersion`/`coreSha256` pin), compatibility-checked on
   every core bump, and the extension slot is the *only* override-capable surface beyond
   the identity zones. The drift guard does not enforce overlay identity (overlays
   legitimately differ) — it enforces core identity and overlay *compatibility*. The
   residual risk is extension-content drift, which is bounded by the near-zero evidence
   (1 in 36) and by the fact that extension sections are reviewable, declared, and
   append-only.

2. **Two machines, two cores, same deck (weak point 4).** If machine A has core v2.5.1
   and machine B has v2.6.0, and both render the same deck, the committed
   `COOPERATION.md` differs. This breaks transport B/C (noise diffs on PR/MR) and the
   audit trail. Mitigation: the deck pins its core version in `meta/version.json`
   (`coreVersion`/`coreSha256`); `parley protocol render` renders from the *pinned*
   version, not the machine's installed version. A machine with a newer core sees the
   deck is pinned to an older one and reports `stale-core` (like `stale-snapshot`),
   but does NOT auto-render a different version into the committed file. The minimum that
   makes this safe: **the deck's pin is authoritative for rendering, the machine's
   installed core is only the source for upgrades.** If `~/.parley` lacks the pinned
   version, the render fails closed with a clear message ("deck pins core v2.5.1; install
   it with `parley core install 2.5.1`") rather than silently substituting the machine's
   version. This is the same fail-closed discipline as `roster render`'s
   `adoptInherited` refusal (`roster_render.go:40-47`).

3. **The generated `COOPERATION.md` is large (1344 lines today).** A render that
   regenerates the whole file on every core bump produces large diffs even for small
   protocol changes, making PR review harder. This is already true of the §2 render but
   at a smaller scale. Mitigation: the render is deterministic and idempotent (only the
   changed sections diff), and the provenance line is the only non-content change. The
   alternative — generating only a diff — is more complex and not warranted at this
   scale.

4. **The detection-only enforcement (weak point 1) is honest but may not satisfy the
   user.** If the user expects "an agent *cannot* change the global core" to mean a
   hard technical barrier, this design delivers "an agent's change is detected,
   attributed, and loud, but a determined agent can still make it." I believe that is the
   correct honest answer — the agent has the user's privileges — but it is a gap between
   the constraint's wording and what is achievable. This should be surfaced explicitly in
   consensus rather than papered over.

5. **Migration of the one real local section.** `librade-algoTrader`'s
   "Project-specific packaged-reference drift" section was destroyed by the 2026-08-06
   sync (00-prompt.md:25-26). The backup exists. Migration must restore it as an
   `## extend` entry, but if that deck has diverged further since the backup, the
   overlay may not capture all local content. A pre-migration diff against the backup is
   the safety check.
