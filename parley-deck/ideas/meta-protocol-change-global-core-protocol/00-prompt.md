---
idea: meta-protocol-change-global-core-protocol
author: user
created: 2026-08-07
participants: [claude-1, codex-1, hermes-1, kimi-1, opencode-1]
status: final
track: deliberation
---

## Problem / idea

`COOPERATION.md` is copied into every deck and hand-editable there. Measured across the fleet
on 2026-08-06, before a one-off sync:

- 36 decks, **8 different `deckVersion` values**, 6 with no version metadata at all.
- §15 (Verification integrity) present in **5 of 36**. §13 in 16 of 36. The §2 roster-authority
  change in **1 of 36**.
- The sync that fixed this was a hand-written script, **not a mechanism**. It will drift again.

Crucially, the drift is NOT project customization. Of 36 decks, only **two** carried headings the
ratified protocol lacks, and only **one** was a real local section — `librade-algoTrader`'s
"Project-specific packaged-reference drift", which said the packaged reference must not overwrite
the deck copy. So the single instance of local protocol authoring in the whole fleet was a rule
about **how the protocol is synced**, i.e. governance that belongs in the core, not a project rule.
(That section was destroyed by the 2026-08-06 sync; it is preserved in
`.parley-protosync-backup-2026-08-06/librade-algoTrader/`.)

Empirically, the only genuinely per-deck content is six identity zones: `**Workspace:**`,
`**Created:**`, `**Transport:**`, the `**Protocol synced:**` stamp, the §2 roster table body, and
the host-handle table body. The protocol sync replaced everything else on 35 of 36 decks with no
loss.

The proposal is to do to `COOPERATION.md` what was already done to §2: stop treating the per-deck
copy as the store, and derive it from a single core.

## Constraints (DECIDED by the user — design within these, do not relitigate)

1. **Version pinning per idea.** An idea that is already open runs to completion under the protocol
   version it started with. The **next** idea in that session uses the new version.
2. **The core lives in `~/.parley`.** Not in the deck, not only inside the installed skill.
3. **Local override AND extension are allowed**, in the sense of object-oriented inheritance: a
   deck may override a defined part or extend with its own.
4. **Neither an idea nor a local session may change the GLOBAL protocol.** The global core changes
   only through a new version. **The user may edit it himself; an agent may not.**
5. **A deck carrying a local override/extension MUST check its compatibility against each new core
   version.**

## What this idea must decide

- **Resolution model.** Concretely: what does "override" and "extend" mean for a Markdown protocol?
  Section-level? Named blocks? What is the precedence order, and how is the effective protocol
  materialized for an agent to read?
- **What is sealed vs open.** Which parts may a deck never override (phases, quorum, artifact
  shapes, §7 itself, §15?) and which may it extend? Evidence above says local authoring is
  near-zero, so a permissive default may be wrong — but a total ban invites silent workarounds
  instead of declared, reviewable exceptions.
- **Enforcement of constraint 4.** How is "an agent cannot change the global core" made REAL rather
  than merely stated? File mode? A CLI that refuses? A guard/test? Note that an agent runs with the
  user's own privileges, so this needs a mechanism, not an honor system.
- **Version pinning mechanics.** How does an open idea record and keep its protocol version, and
  what reads it? Compare with the existing run-roster snapshot (`RosterSnapshot` /
  `RosterRevision`), which solved the same shape of problem for the roster.
- **Compatibility checking.** What exactly is compared, when does it run, and what happens when a
  deck's override is incompatible with a new core — block, warn, quarantine, auto-migrate?
- **Does the deck still keep a committed `COOPERATION.md`?** Trade-off: a generated committed copy
  keeps the deck self-contained and preserves what protocol governed each historical idea; no copy
  removes drift entirely but loses both. Recommend one and say why.
- **Migration** from today's 36 flat copies to the new model, and what happens to decks whose
  protocol is read-only or unreachable.
- **§7 impact.** Today a protocol change needs its own meta idea. With a global core, one change
  hits every project at once. Does §7 need to change too?

## Non-goals

- Redesigning the phases, quorum, or artifact shapes.
- Changing the roster mechanism (settled in `roster-operations-standard`).
- Choosing a specific file format before the resolution model is agreed.
- Implementation. This idea produces a design; code follows in Phases 5-8.
