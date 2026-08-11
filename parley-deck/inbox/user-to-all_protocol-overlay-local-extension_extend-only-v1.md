---
from: user
to: all
idea: protocol-overlay-local-extension
phase: round-02
blocking: yes
date: 2026-08-07
---

## BINDING: v1 is extend-only. `replace` is dropped from this slice.

Escalated by claude-1 during round 2 and answered by the user out of band. Recorded here verbatim
and quoted into round 3 as required.

**Question put to the user.** Their constraint 3 reads "lokálne pretaženie, ALEBO rozšírenie
protokolu" — *a local override, OR an extension, of the protocol*. Does v1 actually need the
override (a deck replacing a core rule), or is the extension (a deck adding local content) enough?

**The user's answer: "Len rozšírenie"** — *extension only*.

The option they selected read, in full:

> "Deck môže PRIDAŤ lokálny obsah, nemôže prepísať pravidlo jadra. Register potrebuje len jeden
> vkladací bod. Zmiznú tým všetky najťažšie problémy návrhu naraz: kontrolný hash cieľového bloku,
> rozsahy blokov, aj neriešiteľná segmentácia položky zoznamu. Prepísanie sa otvorí ako pomenovaný
> follow-up, keď ho reálne budete potrebovať."
>
> *(A deck may ADD local content; it may not replace a core rule. The registry then needs only one
> insertion point. This removes the design's hardest problems at a stroke: the target-block
> verification hash, block extents, and the intractable list-item segmentation. Override opens as a
> named follow-up when it is actually needed.)*

*(Slovak original quoted verbatim per §6 rule 6, with an English translation.)*

## Why it was escalated

claude-1's round-1 file raised it and no participant engaged:

> "Do we actually need `replace` in v1? … `s6.6` (working language) is the only replaceable block,
> and I found no deck that overrides it."

Three participants then spent most of their round-1 budget designing `replace` machinery. The
measured taxonomy says the fleet's real need is additive: 0 of 29 decks override the working-language
rule, and the one true local section (librade-algoTrader's) is an addition.

## What this binds, and what it opens

**Binding.** v1 has ONE operation: extend at `ext-1`. No `replace`. Any round-3 or consensus text
that specifies replace machinery for v1 is out of scope. Override becomes a named follow-up, opened
when a real deck needs it — not built speculatively.

**Not decided by this, and it is the central question for round 3.** All four participants agreed in
round 1 that this slice must ship the block-ID registry. That agreement rested on D2's prohibitions —
no heading-text addressing, no inline markup — which leave a registry as the only way to *address a
block*. With `replace` gone, the only thing v1 must address is a single insertion point, and
codex-1's round-2 formulation defines that point as the end of the normalized core body rather than
as a lookup.

So round 3 must answer plainly, and it may overturn the one thing round 1 was unanimous about:

1. With extend-only, does v1 still need a per-block registry at all?
2. If `depends` collapses from a per-sealed-block hash set to a single whole-core-body hash, what
   compatibility guarantee is actually lost? D10's default dependency set was already "all sealed".
3. If the registry is NOT needed for v1, does the free-migration window argument still require adding
   it now anyway — the release layout is frozen per version once anything is published — or does that
   argument evaporate along with the need?

Do not treat round 1's unanimity as settling any of these. It was reached under a scope the user has
now removed, and unanimity reached from a shared brief is one analysis with four signatures.
