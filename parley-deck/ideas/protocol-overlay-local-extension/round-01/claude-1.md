---
agent: claude-1
idea: protocol-overlay-local-extension
round: 1
date: 2026-08-07
---

## Summary

The overlay should ship **with** the block-ID registry, not after it, because there is no third
addressing option in the ratified record and the migration cost of adding a registry is exactly zero
today and never again. The §2 roster table must **not** become overlay content — `agents.toml` is
already the roster authority and a second roster surface would recreate the drift this whole line of
work exists to end; the overlay should own the *annotations around* the table, which is precisely the
content the 2026-08-06 sync destroyed. One hazard is disqualifying as a substrate: table zones are
located by **prose match**, so a core column rename silently empties every deck's roster, and that
must be fixed in this slice or the overlay's central promise is false on arrival.

## Proposed approach

### D-a — Overlay file syntax: YAML manifest + delimited Markdown payloads

`parley-deck/protocol-overlay.md`, one file, structured as a YAML frontmatter **manifest** followed
by **payload sections** delimited by HTML comments:

```markdown
---
core-version-range: ">=1.0.0 <2.0.0"
operations:
  - op: replace
    target: s6.6
    base-sha256: <hash of the block this was written against>
    rationale: "German-language deck; working language overridden by user directive 2026-08-07."
  - op: extend
    at: ext-1
    id: librade.packaged-reference-drift
    depends-on: [all-sealed]
    depends-sha256: <hash>
    rationale: "Records that the packaged reference must not overwrite this deck's copy."
---

<!-- parley:payload s6.6 -->
6. **Working language.** ...
<!-- /parley:payload s6.6 -->

<!-- parley:payload librade.packaged-reference-drift -->
## Project-specific packaged-reference drift
...
<!-- /parley:payload librade.packaged-reference-drift -->
```

Defence. The payload **is** protocol prose in Markdown; authoring it inside a YAML or TOML string
means escaping Markdown into a string literal, which is how multi-line protocol text becomes
unreadable and unreviewable in a diff. The metadata is structured and must be machine-validated, so
it belongs in frontmatter. The delimiters are HTML comments because they are invisible when the file
is read as Markdown, greppable, and — importantly — D2's ban on inline markup binds the **core
body**, not the overlay, so nothing here contradicts it.

One rule the delimiters need: a payload's delimiter id MUST match a manifest entry and vice versa,
one-to-one, fail closed on either orphan. Otherwise the file has two sources of truth about what it
declares.

### D-b — Ship the registry in this slice

**Yes, and the sequencing is the whole argument.** D2 forbids addressing by heading text and forbids
inline markup in the core. Those two prohibitions leave exactly one mechanism — a registry — and the
record contains no third option. Shipping the overlay without one means inventing an unratified
addressing scheme under delivery pressure, which is how `droppedContent` got nine review rounds.

The cost argument is time-sensitive in a way that is easy to miss: a release directory currently
holds exactly one file, releases are **write-once**, and **zero releases exist**. Adding a registry
file to the release layout costs nothing today. The moment a core version is published, the layout
is frozen for that version forever and every later addition is a migration.

**Concrete consequence for the user, and it is urgent:** I handed over the
`parley protocol publish --version 1.0.0` command earlier today. That publish should be **held**
until this idea settles the release layout. Publishing now does not break anything — it just spends
the free window and makes 1.0.0 a release that can never carry a registry.

Block extent, given that `s6.6` is a **list item** and not a heading: the registry stores the block's
identity and its extent explicitly rather than inferring it from Markdown structure. My preference is
a registry entry naming the enclosing heading plus the ordinal of the list item, validated at render
time against a stored hash — so a drifting core fails closed rather than replacing the wrong item.
Byte offsets are the obvious alternative and I think they are worse: they turn any whitespace edit in
an unrelated part of the core into a registry migration.

### D-k — The §2 roster table is NOT overlay content

Position (c) from the brief, sharpened. The overlay owns **nothing inside** the table.

- `agents.toml` is the roster authority. That shipped in 1.41.0 and the §2 table is explicitly a
  generated, non-authoritative view. Letting an overlay write roster rows creates a second roster
  surface and re-opens the exact schism that work closed.
- The bespoke columns 11 of 29 decks invented (CLI/runtime, Model, Effort, State, Backend,
  Invocation) are **schema requests**, not local protocol rules. They belong in
  `parley roster render` growing columns sourced from `agents.toml` — one renderer, one authority.
- What the overlay SHOULD own is the material `agents.toml` genuinely cannot express and that the
  sync actually destroyed: dated user directives, invocation gotchas, MANUAL-Bash caveats,
  roster-swap history. In the four decks I repaired today that is *all* of the real loss.

This matters beyond tidiness. If the overlay may write roster rows, then a deck can put a
non-rostered agent into the rendered table, and the roster table stops being evidence of anything.

### D-c / D-g — Placement, and what the loss report means

**Placement:** a named block ID, not "end of file". "Append to the end" is not a position, it is the
absence of one, and it is unstable the moment the core grows an appendix.

**Loss report — option (iii), with one addition.** Accept a noisy first report plus an explicit
migration note, and grow `RenderResult` a **provenance** field so each output block is labelled
`core` / `identity` / `overlay`. That lets the report say *"carried by overlay"* factually instead of
**exempting** anything.

I want to be blunt about why option (i) is dangerous. Exempting overlay-carried content from the loss
report reintroduces the silent-erasure class that cost nine review cycles — a deck whose overlay
*claims* to carry a block but whose payload is empty or misaddressed would report nothing lost. The
constraint that must survive verbatim is the one already in the source: an empty report means **no
line disappeared**, not that no meaning was lost.

## Concerns / open questions

- **Do we actually need `replace` in v1?** The measured evidence is that the fleet's real need is
  *additive*: dated directives, caveats, one genuine new section. `s6.6` (working language) is the
  only replaceable block, and I found no deck that overrides it. A v1 with **extend only** would be
  materially smaller and would still discharge the owner's constraint 3. I am not confident enough to
  propose dropping `replace` outright, but it deserves an explicit decision rather than inheritance.
- **D-u — is the owner's constraint 3 currently met?** kimi-1's signoff made the near-empty open
  surface acceptable *on condition that rank 3 shipped that cycle*, and said constraint 3 would
  otherwise be unmet. It did not ship. Someone other than me should judge whether this idea
  discharges that condition, since I closed the cycle that broke it.
- **kimi-1's `parley protocol audit` was dropped with no follow-up number.** Its own signoff proposed
  it become DF-5 and no DF-5 exists. Either number it or record the refusal.

## Risks

- **H9 is disqualifying as a substrate (PRIMARY).** `isTableHeader` matches on prose:
  `internal/protocolcore/render.go:129-133` requires the line to start with `| Agent ID` and contain
  `"Workspace"` or `"Host handle"`. A core that renames a column to `Workdir` therefore matches
  nothing, and the render emits header and separator with no data rows — the deck's entire roster
  silently gone. Shipping an overlay whose selling point is "your local content survives a render"
  onto a substrate that erases the single highest-volume local content class would be a promise the
  code contradicts. **D-t must be "now".**
- **Four writers, one file (SECONDARY, from the brief).** `protocol render`, `roster render` and
  `preflight.syncConsumerProtocol` each write `parley-deck/COOPERATION.md` with their own location
  logic. The overlay is a fourth content source into a file with no single owner. I regard
  consolidating to one writer as in scope; leaving it is how a render is silently undone by the next
  command.
- **The two promissory notes must retire atomically (PRIMARY).** `internal/app/protocol.go:211`
  prints "the overlay that will carry it is ratified but not shipped", and
  `parley-deck/COOPERATION.md:767` says "the deck's own overlay, once that ships". If they do not
  change in the same commit, the CLI and the protocol text contradict each other in opposite
  directions.
- **First render is noisy fleet-wide (SECONDARY).** Most decks lag the packaged default purely from
  sync timing, so nearly every deck's first real render produces a non-empty removal report before
  the overlay contributes anything. If that lands without a migration note, operators learn to
  ignore the G1 report — which is worse than not having it.
- **Process risk on this idea specifically.** I wrote the kickoff, produced the scoping brief, and
  hold strong positions on D-b, D-k and H9. Reviewers should treat my framing as anchoring and push
  back on it directly; §15.1 forbids me from ruling on my own claims.
