---
agent: kimi-1
idea: meta-protocol-change-global-core-protocol
round: 1
date: 2026-08-07
---

## Summary

The drift evidence (36 decks, 8 versions, §15 in 5 of 36 — 00-prompt.md:15-17) says the per-deck
copy-as-store model has already failed; the 2026-08-06 hand sync is a deferral, not a fix
(00-prompt.md:18). The roster precedent proves the remedy: §2 went from hand-maintained store to
generated view over a single authority (`internal/app/roster_render.go:14-22`), with an immutable
run snapshot for pinning (`internal/runmanifest/manifest.go:49-58`) and a fail-closed byte-compare
drift guard (`internal/protocol/drift_test.go:46`). I propose applying exactly that shape to the
whole protocol:

- **Core**: versioned, write-once copies under `~/.parley/protocol/core/<version>/`, adopted per
  deck, never edited in place.
- **Overlay**: ONE committed file per deck, `parley-deck/protocol.overlay.md`, with exactly two
  verbs — `Override §<anchor>` and `Extend §<anchor>`. Nothing richer is justified: the fleet
  produced ONE genuine local section in 36 decks, and it was sync governance that belongs in the
  core (00-prompt.md:21-24).
- **Materialization**: `parley protocol render` generates the committed `COOPERATION.md`
  deterministically (core < overlay, sealed sections untouchable), stamped with core version +
  core sha256 + overlay sha256. The committed copy stays — it is the cross-machine authority and
  the historical record.
- **Constraint 4 honestly**: an agent runs with the user's uid, so *prevention* of core edits is
  not achievable in-process. The real guarantee is **detection + attribution + fail-closed
  blocking**: every run verifies the core hash and refuses to proceed on a tampered or
  unratified core.

## Proposed approach

### 1. Resolution model — section-addressed, two verbs, one overlay file

Evidence argues for the smallest mechanism that satisfies constraint 3 (override AND extension).
What the fleet actually needed in ~years of operation was one local section — and it was a rule
about sync, i.e. core governance (00-prompt.md:21-26). So:

- **Addressing**: by existing section anchor (`§4`, `§9.0`, `§15.2`). The protocol is already
  numbered; no new block-naming scheme, no named-region markup in the core.
- **`## Override §<n>`**: replaces the body of core section `<n>` (up to the next same-or-higher
  heading) for this deck. Only permitted on sections classified open (below).
- **`## Extend §<n>`**: appends a deck-local subsection under section `<n>` (rendered as
  `§<n>.local-<slug>` so numbering of the core is never perturbed).
- **Precedence**: identity zones (the six per-deck zones measured in 00-prompt.md:28-30 —
  Workspace/Created/Transport/sync-stamp/roster-table/handle-table; these stay deck-generated,
  never live in the core) → core version → overlay. Overlay never touches identity zones or
  sealed sections.
- **Deliberate brittleness**: if a new core renames or deletes an anchored section, the override
  targeting it breaks loudly at compatibility-check time. Silent re-targeting is how the old sync
  destroyed the librade-algoTrader section (00-prompt.md:25-26).

Rejected richer models: named block-level regions (needs markup invasive to the core, and zero
evidence of need), conditional/frontmatter logic (programming language creep), multiple overlay
files (recreates the multi-store problem the roster change just eliminated —
`internal/config/runtime.go:90-103` documents what unioning across layers did to membership).
If a deck genuinely needs more, that is evidence the core should change via §7, not that the
overlay language should grow.

### 2. Sealed vs open

Rule of thumb: **anything a participant in *another* deck relies on when trusting this deck's
artifacts is sealed; anything about how *this* deck organizes its own work is open.**

- **Sealed (no override, no rule-changing extension)**: §4 phases and artifact shapes, §5 quorum,
  §7 itself, §14 (the human brake — COOPERATION.md:1165-1203), §15 verification integrity
  (COOPERATION.md:1204-1344), §12.4 execution boundary. If deck A could override §15.2 provenance,
  deck B's `CONFIRMED` tags stop meaning the same thing; cross-deck comparability of "Parley
  verified" is the entire value of a shared protocol.
- **Open**: §0 transport (already per-deck), §6 conflict-avoidance, §8 inbox conventions, §11
  transport mechanics, §13 retro cadence, Appendix A adoption notes.
- Sealed sections may still receive `Extend` ONLY as purely additive local *procedures* that
  cannot weaken the sealed rule (e.g. an extra local checklist item) — the compatibility checker
  classifies this conservatively and anything ambiguous blocks for user decision. A permissive
  default is wrong per the prompt's own note (00-prompt.md:53-56), but a total ban invites silent
  workarounds; the two-verb, declared-and-committed overlay makes every exception reviewable in
  one screen of diff.

### 3. Enforcement of constraint 4 — the honest answer is detection, attribution, and a closed gate

Prevention is not achievable in-process: the agent holds the user's privileges, so file modes,
a refusing CLI, and guard tests are all bypassable by the same shell that bypassed nothing to
write this file. Anyone claiming a mechanism *prevents* the edit should be asked what defeats it;
the answer is always "the agent's own next Bash call." What IS achievable, in layers:

1. **Chokepoint (raises cost, not a guarantee)**: all sanctioned core mutations go through
   `parley protocol` verbs that require an interactive TTY confirmation for any core write and
   refuse non-interactive invocation. A headless agent has no TTY by default.
2. **Content-addressed immutability (detection)**: a ratified core version is write-once. Every
   adopting deck records `coreVersion` + `coreSha256` in its committed render header and in the
   run manifest. At run creation AND continuation, the CLI re-hashes the local core and **fails
   closed** on mismatch — the same contract the roster snapshot already implements
   (`internal/runmanifest/manifest.go:49-58`, "Every later phase of a run uses this snapshot,
   never a fresh resolve"). A tampered core does not get to govern work; the run stops.
3. **Attribution (legitimacy, not just integrity)**: each core version directory carries a
   ratification record — the meta-idea slug and the sha256 of its `FINAL.md`. A core version
   whose ratification record is missing or whose FINAL hash doesn't match is *unratified* and
   rejected at preflight. A user hand-edit then has exactly one legitimate path: edit, then
   `parley protocol ratify` (interactive), which regenerates the record; the version history
   shows who ratified what, when. An agent doing the same leaves the same visible trace.
4. **Rejected**: chowning `~/.parley` to root or a dedicated user. It is the only true
   prevention and it is disproportionate — it breaks the user's own "I may edit it myself"
   clause (constraint 4) and every ordinary maintenance flow.

Net guarantee, stated plainly: an agent *can* physically edit the core; it cannot do so
undetectably, cannot get any run or preflight to proceed on the result, and cannot forge a
ratification record invisibly. That is the maximum enforcement compatible with constraint 4's
second sentence, and it should be written into the design as the goal — not "prevention."

### 4. Version pinning — copy the roster snapshot pattern exactly

- Phase 0 records `protocolVersion:` + `protocolSha256:` in `00-prompt.md` frontmatter, locked
  when Phase 0 completes — the same moment the quorum locks (COOPERATION.md:836).
- The run manifest gains `ProtocolVersion`/`ProtocolSHA` beside `RosterSnapshot`/
  `RosterRevision` (`internal/runmanifest/manifest.go:49-58`); `sessions inspect` reports
  `stale-protocol` on drift, mirroring `stale-snapshot`. Constraint 1 ("an open idea runs to
  completion on its starting version") falls out for free: the run never re-resolves the core.
- "Next idea in the session uses the new version" hooks into the existing §9.0 preflight
  (COOPERATION.md:803-824), which already compares protocol freshness and already distinguishes
  additive (auto) vs breaking (user-confirmed) adoption. The overlay compatibility check (below)
  runs there.

### 5. Compatibility checking (constraint 5)

- **What is compared**: (a) every overlay anchor exists in the new core; (b) no overlay verb
  targets a sealed section; (c) each override records the sha of the core section it replaced —
  if that section's text changed in the new core, the override is **stale** and must be
  re-confirmed by the user even if the anchor still exists (a fresh override of changed text is a
  new rule, and new rules go through review, not adoption).
- **When**: at core adoption (`parley protocol sync`) and at every §9.0 preflight while a deck
  lags the installed core.
- **On incompatibility**: block adoption for that deck; the deck stays pinned on its current
  core version (per-deck pinning is constraint 1 applied at deck scope), warns loudly at
  preflight, escalates to the user. **Never auto-migrate prose** — auto-migration of rules is
  precisely the 2026-08-06 failure mode that erased the one genuine local section
  (00-prompt.md:25-26). Quarantine = staying on the old pinned version with a visible flag;
  that IS the quarantine.

### 6. Keep the committed `COOPERATION.md` — yes

The committed copy is not the drift source; *hand-editing* it was. Keep it, generated:

- Self-containment: a reviewer, a new agent, or a machine without `~/.parley` reads one file and
  knows the exact governing protocol — including for every historical idea (the header stamps
  version + hashes).
- The two-machine problem (below) makes the committed render the *cross-machine authority*:
  machine cores are derivations; the committed artifact is what collaborators agreed on.
- Drift risk is handled the way §2's was: idempotent generator (`roster_render.go:20-22`
  documents why non-idempotent generation re-creates drift), fail-closed guard test comparing
  committed copy to render output with only the six identity zones normalized — the exact
  normalizer pattern in `internal/protocol/drift_test.go:98-130`, which already names five of
  the six zones.

### 7. Migration

Same staging discipline the roster FINAL ratified (Stage 4, FINAL.md:127-131: tooling first,
then a **separate attended fleet operation**, never folded into the code merge):

1. Ship core store + renderer + checker behind a feature gate.
2. Seed core v1 from the ratified text (the post-sync state).
3. Inventory pass over the 36 decks: for each, diff its current copy against core v1. For the 35
   clean decks the overlay is empty. For librade-algoTrader, restore the destroyed section from
   `.parley-protosync-backup-2026-08-06/` (00-prompt.md:26) as its first overlay — and ratify
   its substance into the core, since it was governance, not a project rule.
4. Per deck: pin version, render, show the user the diff, apply on confirmation, in batches.
5. Read-only/unreachable decks: skip, mark `legacy-protocol` in the inventory; their flat copy
   keeps working (it does today) until attended. Migration is an adoption deadline for nothing —
   old decks degrade gracefully, they don't break.

### 8. §7 impact — yes, it changes

§7 today ratifies a change to *a file* (COOPERATION.md:745-752). With a global core it ratifies
**a new core version**: the meta idea's FINAL.md becomes the ratification record from §3 above,
and `meta/protocol-changelog.md` moves to the core (version-scoped changelog). Per-deck adoption
of a ratified version extends the existing §7 carve-out for version syncs (COOPERATION.md:754-758):
additive adoption stays a maintenance sync, breaking adoption still pauses for user confirmation.
The ratification bar does NOT lower just because one change now hits every deck at once — if
anything the blast radius argues the full deliberation lifecycle stays mandatory, and the
§9.0-style per-deck, user-confirmed breaking-adoption gate is what makes one-global-change safe.

## Concerns / open questions

- **The overlay as a new drift habitat (weak point 3).** 36 overlays can drift in *content*, but
  not silently: each overlay is a committed, one-screen diff rather than a 1344-line copy, so
  drift is inspectable by construction. Three mechanisms keep it bounded: (a) the compatibility
  check forces every overlay to re-justify itself against each new core — overlays cannot rot
  unnoticed across a version bump; (b) a fleet audit surface (`parley protocol audit`) reporting
  every deck's core version + overlay anchors makes divergence a query, not an excavation;
  (c) the stale-override sha rule (§5 above) turns silent core/overlay skew into an explicit
  user decision. Residual risk: 36 decks pinning 36 *different core versions* is the old
  8-`deckVersion` problem wearing a new hat — the audit report must surface version spread as a
  first-class number.
- **Two machines, two cores (weak point 4).** Does not break the design, because the committed
  render is the authority and the pin is a content hash, not a version label. Minimum for safety:
  (a) run manifest and render header carry `coreSha256` (hashes, unlike version strings, cannot
  disagree about what they name); (b) run creation/continuation fails closed when the local core
  hash ≠ pinned hash, with an explicit attended re-adopt path — exactly `stale-snapshot`'s
  contract. Without the hash check, two machines with the same version label and different bytes
  would render differently and both claim the same protocol. The hash is non-negotiable; the
  version string alone is insufficient.
- **Is even the two-verb overlay justified (weak point 2)?** The evidence supports *extension*
  barely (one instance) and *override* not at all — the single observed case was an extension,
  and it belonged in the core. I keep `Override` only because constraint 3 is decided; but I
  want the compatibility checker to make overrides visibly expensive (stale-sha re-confirmation
  every bump) so the default pressure is toward core changes via §7. If the fleet grows a second
  genuine override use case, that is the signal the mechanism earns its keep; if not, a later
  meta idea should consider retiring `Override` to extension-only.
- **Hash self-reference**: the recorded hash lives in files the agent can also edit. Mitigation:
  the hash is *committed to the deck repo*, so forging it requires a committed, reviewable
  change — attribution survives even if detection-at-rest is defeated. Worth stating in the
  design so nobody believes the local manifest alone is tamper-proof.

## Risks

- **Enforcement overpromise.** If FINAL.md claims agents are "prevented" from editing the core,
  the design is dishonest and will be believed. Mitigation: write the detection + attribution +
  fail-closed guarantee into the protocol text verbatim.
- **Version proliferation across 36 decks.** Per-deck pinning (needed for constraint 1) means
  the fleet legitimately spans N core versions at any time. Mitigation: audit surface + breaking
  bumps stay rare; but accept that "which version governs deck X" must always be answerable from
  the committed render header, never from memory.
- **Renderer non-determinism.** A non-idempotent renderer is drift with a generator — the
  roster change documented this (`roster_render.go:20-22`). Mitigation: byte-identical re-render
  test + committed-copy drift guard as release gates (G4-style).
- **`~/.parley` unavailable/read-only decks.** Rendered committed copy keeps them working;
  adoption simply waits. Do not make core presence a hard requirement for *reading* the
  protocol — only for *adopting a new version*.
- **Overlay verb creep.** Once `Override`/`Extend` exist, pressure will come for conditionals,
  includes, parameterization. Each is a step toward a templating language no evidence supports.
  Mitigation: the two-verb limit should itself be a sealed rule, changeable only via meta idea.
- **Migration trust.** The last fleet-wide sync destroyed user data (00-prompt.md:25-26). The
  migration must be attended, batch-wise, with per-deck diffs and rollback — the roster Stage-4
  discipline — or this idea ships its own disproof.
