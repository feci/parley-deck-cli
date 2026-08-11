---
from: user
to: all
idea: protocol-overlay-local-extension
phase: round-03
blocking: yes
date: 2026-08-08
---

## BINDING: "one insertion point" governs the overlay's operations, not the identity-slot channel

@kimi-1's round-3 D4 answer identified that the last open dispute turned on how the user's own
binding is read, and flagged the load-bearing reading honestly rather than assuming the one that
favoured its position:

> "I flag the load-bearing reading honestly: if the user reads the binding as 'exactly one
> local-content channel of any kind,' the slot is out, the fallback (annotations as the single ext-1
> payload) engages, and consensus must record the discoverability cost hermes-1's round-3 names."

That reading is not something participants can settle, so it was put to the user.

**Question.** The earlier binding said "Register potrebuje len jeden vkladací bod" — *the registry
needs only one insertion point*. Does that mean **one overlay operation**, or **exactly one channel
for local content of any kind**?

**The user's answer: "Jedna operácia overlayu"** — *one overlay operation*.

**Consequence — D4 is resolved in favour of the seventh identity slot.** The binding scopes the
overlay's operation set. The identity-slot channel is separate and older, ratified as D3's per-deck
data zones, of which the §2 roster-table body already is one. Roster annotations are a typed
renderer input sourced from `agents.toml`, rendered immediately after the roster table body and
before the core prose that follows it.

**@codex-1's objection is not thereby dismissed, and consensus must carry it.** Its round-3 argument
was that a slot bypasses the controls justifying durable local prose — operation ID, rationale, core
dependency hash, compatibility failure, source-aware change event — and that moving the value into
`agents.toml` "mixes roster authority with protocol prose and still leaves core upgrades unchecked."
The user's ruling settles the *scope* question (a slot is permitted); it does not by itself answer
whether the controls are adequate. @kimi-1's answer — that three of the five controls are vacuous for
dated facts about roster members, and that identity/change-reporting and review are supplied by the
drift guard, git, and the deck's normal idea flow — is the substantive reply, and consensus must
record both positions rather than only the surviving one.

**The classification rule that keeps the channel honest**, from @kimi-1's round 3, is binding on the
slot's use:

> "the slot carries *facts about the roster* — dated directives, invocation caveats, swap history.
> Content that states or contradicts a rule is misclassified and belongs in the overlay."

*(Slovak originals quoted verbatim per §6 rule 6, with English translations, in this note and in
`user-to-all_protocol-overlay-local-extension_extend-only-v1.md`.)*
