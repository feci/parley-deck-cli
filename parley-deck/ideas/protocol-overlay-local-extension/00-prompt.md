---
idea: protocol-overlay-local-extension
author: claude-1
created: 2026-08-07
participants: [claude-1, codex-1, hermes-1, kimi-1, opencode-1]
status: final
track: deliberation
---

## Problem / idea

Design the **protocol overlay** — rank 3 of the staging ratified in
`ideas/meta-protocol-change-global-core-protocol`. It is the mechanism that lets a deck carry
project-local protocol content across a core render without letting any deck weaken the core.

It is the owner's binding constraint 3, stated in their own words:

> "moze existovat lokalne pretazenie, alebo rozsirenie protokolu, nieco ako v objektovom
> programovani"
> *(there may be a local override, or extension, of the protocol — something like in
> object-oriented programming)*

with constraints 4 and 5 bounding it:

> "ziadna idea nemoze alebo lokalna session nemoze zmenit globalny protokol, ten je zmenitelny len
> novou verziou. Samozrejme user si ho sam zmenit moze. Ale agent nie."
> *(no idea and no local session may change the global protocol; it changes only by a new version.
> The user may change it themselves. An agent may not.)*

> "Ak sa to ale stane, vzdy si potom musi kontrolovat kompatibilitu s novou verziou protokolu."
> *(if that does happen, compatibility with the new protocol version must always be re-checked.)*

*(Slovak original quoted verbatim per §6 rule 6, with an English translation.)*

## Why now

1. **DF-2 is blocked on it.** A fleet migration to the core today would re-inflict the 2026-08-06
   damage on every deck holding local content — announced instead of silent. Measured: rendering
   `auftra` onto a core release reports 13 lines lost in §2, and applying it drops genuine content
   (`ideas/meta-protocol-change-global-core-protocol/IMPLEMENTATION.md`, "Post-completion
   operations").
2. **The core store is still empty**, so the release layout can still change for free. That window
   closes on the first published release, because releases are write-once.
3. **A signoff condition may be unmet.** kimi-1 signed off on the near-empty open surface *on
   condition that rank 3 shipped in that cycle*, and noted constraint 3 would otherwise be unmet.
   It did not ship. Round 1 should say plainly whether that condition is discharged or must be
   re-raised.

## Source material

`00-scoping-brief.md` in this directory is **required reading** and is copied here in full rather
than referenced across workdirs (§6 rule 4). It was produced by four independent internal readers
(ratified design, renderer, CLI surface, and an empirical survey of 29 real decks). **Internal
helpers are not participants** and own no canonical artifact — the brief is input, not a round.

Treat it as evidence, not as decided design. It contains:

- **§1–2** what the ratified record actually says, plus a measured taxonomy of what 29 real decks
  carry (0 of 29 are byte-identical to the packaged default; for 27 of them ALL divergence lives in
  just two zones — the header block and §2).
- **§3** eighteen design hazards (H1–H18) with `file:line` evidence, several probe-confirmed.
- **§4** twenty-four genuinely open decisions (D-a … D-x) phrased as either/or.
- **§5** what is out of scope — guard against creep into rank 2 or DF-1.

## What round 1 must produce

Each participant writes `round-01/<agent-id>.md` independently, before reading anyone else's.

Answer these, and say **explicitly when the record never decided something** rather than inventing
a ratification:

1. **D-a — overlay file syntax.** Never decided; all three candidate syntaxes are dead in the
   record. Propose one concrete grammar and defend it.
2. **D-b — does this slice ship the block-ID registry?** D2 forbids addressing by heading text and
   forbids inline markup, and no registry exists (a release directory holds exactly one file). If
   not now, say what the third addressing option is — because the record does not contain one.
   Weigh it against the closing free-migration window above.
3. **D-k — is the §2 roster table overlay content at all?** This is the largest real override class
   (11 of 29 decks rewrite its schema) and D3's open surface does not cover it. `agents.toml` is
   the roster authority; do not create a competing roster surface.
4. **D-c / D-g — placement and loss reporting.** Where `ext-1` renders, and what the G1 loss report
   means once a second content source exists. The hard constraint that must survive: *an empty
   report means no line disappeared, NOT that no meaning was lost.*
5. **The hazards you judge disqualifying.** H9 is probe-confirmed: table zones are addressed by
   prose match, so a core column rename silently empties every deck's roster. Say whether the
   overlay may ship onto that substrate.

Then rank the remaining open decisions by whether they must be settled **in this slice** or can be
deferred with a named follow-up.

## Constraints

- **Read-only review of the repository**: no edits and no git write commands while round 1 is open.
- The tree does not move while a round is open.
- English only for every file under `parley-deck/`.
- Redact obvious secrets; never copy tokens into artifacts.
- §15 applies: tag claim provenance (`PRIMARY` / `SECONDARY` / `RECALL`; untagged reads as RECALL),
  and a PRIMARY claim needs a stable locator or quoted command output.

## Roster note

`parley preflight` on 2026-08-07 found **kimi-1 unavailable (`no-pong`)**. It is kept in
`participants:` rather than dropped, because excluding it is a §9.0 gate that needs the user's
confirmation, and async participation is permitted. If it produces no artifact, that is recorded in
`inbox/`, not treated as agreement.
