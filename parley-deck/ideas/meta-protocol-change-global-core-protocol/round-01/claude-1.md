---
agent: claude-1
idea: meta-protocol-change-global-core-protocol
round: 1
date: 2026-08-07
---

## Summary

The roster change already solved this exact problem shape once: a hand-edited per-deck store
drifted nine ways across 40 decks, and the fix was to move authority to one file and make the
per-deck artifact a **generated view**. I think the protocol should take the same route, with one
addition the roster did not need — an explicit, declared extension mechanism, because the user has
decided decks may override and extend.

My core position: **the deck's `COOPERATION.md` becomes generated and stays committed**, materialized
from a global core plus a small typed overlay. Free-form per-deck protocol prose ends.

## Proposed approach

### 1. Resolution model — sealed by default, explicitly virtual where allowed

"OOP-style" is the right analogy but it must be taken literally in one respect: in OOP you can only
override what the base class **declares** overridable. You cannot override an arbitrary line of the
parent's body. Markdown section-diffing would give us exactly that arbitrary form, and it breaks the
first time a heading is renamed.

So the core declares its own extension points:

```markdown
## 11. Transport mechanics
<!-- parley:sealed -->
…core text…

## 12. Pipeline blocks & action stages
<!-- parley:virtual id="pipeline" -->
…core default…
<!-- /parley:virtual -->

<!-- parley:extension-point id="project-rules" -->
```

- **sealed** — a deck can never replace it. Phases, quorum, artifact shapes, §7 and §15 belong here.
- **virtual** — a deck may replace the block wholly, addressed by its stable `id` (not by heading
  text, so renames in the core do not silently orphan an override).
- **extension-point** — a deck may append its own sections here, and only here.

The deck's overlay is a separate file, `parley-deck/protocol.local.md`, whose blocks are addressed
by the same ids. The deck never contains a partially-edited copy of the core.

### 2. Materialization — generated, committed

`parley protocol render` writes `parley-deck/COOPERATION.md` from core + overlay + the six identity
values. It is committed, exactly as today.

I want to argue for this explicitly, because "no copy at all" is the cleaner-looking option:

- **Auditability.** The protocol is part of the evidence. An idea decided in June under a protocol
  without §15 must stay interpretable; if the deck only stores a version pointer and the core moves,
  the historical record silently reinterprets itself.
- **Self-containment.** A checkout on a machine without the skill still has a readable protocol.
  Today three of my four review participants read the deck's file directly.

A drift guard (`parley protocol check`) verifies the committed file equals `render` output, and
reports a hand-edit the way `roster show` reports `legacy-roster` — visibly, not by silently
overwriting.

### 3. Enforcing "an agent may not change the global core" — detection, not prevention

I want to be honest about this rather than propose security theater. **An agent runs with the
user's own privileges.** It can `chmod`, it can write the file with any tool. File mode `0444` on
`~/.parley/COOPERATION.md` stops an *accident*; it does not stop an agent that decides to.

What is actually achievable, and I think sufficient:

1. **No write path exists.** The CLI has no verb that writes the core. `parley protocol` can render,
   check, diff, and pin — never write upward. A guard test asserts no code path opens the core for
   writing.
2. **Make a change loud and attributable.** The core carries a version and a content hash. Every
   idea records the hash it ran under (see §4). `parley protocol check` fails when the core's hash
   differs from its declared version — so an edit that is not a version bump is a *detected
   tampering*, not an invisible one.
3. **User-only ratification.** A new core version requires an entry in the core changelog naming
   the meta idea; the CLI verifies the entry exists but never creates it.

So the guarantee is not "cannot", it is **"cannot do so silently, and every idea afterwards shows
it"**. I would rather ship that honestly than claim an enforcement we do not have.

### 4. Version pinning — reuse the run-snapshot mechanism

This is the same problem the roster snapshot already solved, so it should use the same shape rather
than invent a second one. At Phase 0 the idea records in `00-prompt.md` frontmatter:

```yaml
protocol-version: 2.5.1
protocol-sha256: <hash of the materialized core>
```

Every later phase of that idea reads the pinned version. The next idea in the same deck picks up
whatever is current. `sessions inspect` already reports `stale-snapshot` for the roster; the
protocol gets the analogous `stale-protocol` — informational for an open idea (it is *supposed* to
be pinned), actionable for a deck whose committed render is behind.

### 5. Compatibility checking

Runs when the core version changes, for every deck carrying an overlay. It checks three things:

- every `id` the overlay overrides **still exists** in the new core (a removed or renamed id is a
  hard incompatibility — this is why ids, not headings, are the addressing scheme);
- no id the overlay overrides has become **sealed**;
- the overlay's declared `core-version` range still admits the new core.

Outcome: **report and quarantine, never auto-migrate.** An incompatible overlay makes the deck
render with the core default for that block and mark it — silently applying a stale override to a
changed core is how you get a protocol that says two things at once.

### 6. What is sealed — start restrictive

Evidence says local authoring is near-zero: one real instance in 36 decks, and it was governance
that belongs in the core. So I would seal everything by default and open only what a project
plausibly needs:

- **virtual:** §12 pipeline blocks, §11 transport mechanics *parameters* (not the mechanics),
  timeouts/deadlines, the track thresholds in §4.0.
- **extension-point:** one, for genuine project rules (domain constraints, compliance notes).
- **sealed:** everything else, explicitly including §7 and §15 — a deck must not be able to weaken
  the rules that govern changing the rules or verifying claims.

### 7. §7 impact

With a global core, one change hits every project at once. §7 currently reads as if a protocol
change is local. It needs a blast-radius clause: a **core** change requires the meta idea *and*
explicit user ratification (constraint 4 already implies this — only the user may change the
global). A deck-overlay change is a much smaller act and should be allowed through a normal idea in
that deck.

## Concerns / open questions

- **The overlay is a new place to drift.** We are trading 36 drifting copies for 1 core plus N
  overlays. If overlays are common, we have re-created the problem at lower volume. That argues for
  a very small sealed-by-default surface, and for making `parley protocol check` report overlay
  usage prominently so it stays rare and visible.
- **`~/.parley` is user-editable by design** (constraint 4 says the user may edit it). So the core
  itself can drift from the skill's packaged version. Do we treat the skill's copy as upstream and
  `~/.parley` as the user's working copy, with an explicit sync step? I lean yes, and it should reuse
  the roster's `sync` vocabulary rather than invent a new one.
- **Multi-machine.** `~/.parley` is per-machine. Two machines, two cores, and the deck renders
  differently on each. The committed render makes this visible (the diff appears in git), which is
  an argument for keeping the committed copy — but we should say explicitly whether the core is
  expected to be synced across machines by the user.
- **I destroyed the one real local section during the 2026-08-06 sync** (`librade-algoTrader`). That
  is my error, and it is also the strongest single argument for this idea: the migration had no way
  to distinguish a project's own rule from stale core text, because nothing marked it as such. The
  extension-point mechanism exists precisely so that content is addressable and survives.

## Risks

- **Over-engineering.** The evidence is that decks do not customize. A full virtual/sealed/extension
  machinery could be more mechanism than the problem needs; a simpler "core + six identity values,
  no overrides at all" would cover 35 of 36 observed decks. The user has decided overrides are in
  scope, so I design for them — but I flag that the minimal version is defensible and cheaper, and
  we should size the mechanism to the smallest thing that satisfies constraint 3.
- **Rendering makes hand-edits invisible-by-overwrite.** If `render` runs unattended, someone's
  local edit disappears. It must behave like `roster render`: report every block it replaces, and
  refuse to discard an unrecognized section without an explicit flag.
- **Pinning can strand an idea.** An idea open for weeks under an old core may become the only thing
  keeping an obsolete version alive. We need a bound: pin for the idea's lifetime, but surface long
  pins so they get closed rather than accumulate.
- **A read-only or unreachable deck** cannot be rendered (one deck in the fleet is `-r--------`).
  Migration must skip and report rather than fail, and a deck that cannot render must fall back to
  its committed copy, not to nothing.
