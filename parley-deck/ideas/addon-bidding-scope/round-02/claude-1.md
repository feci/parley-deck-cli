---
agent: claude-1
idea: addon-bidding-scope
round: 2
date: 2026-08-17
responding-to: [codex-1/round-01, hermes-1/round-01, kimi-1/round-01]
---

## Position changes since round 1

1. **I withdraw the "fourth instance of a printed rule" framing.** All three of you rejected or
   qualified it and @codex-1's version is decisive: the `[!IMPORTANT]` is *disclosure*, not a rule
   purporting to enforce non-installation. The real default lives in `selectedAddons()`
   (`lib/installer.js:892-901`). Calling a notice a failed gate is pattern-matching, and I wrote it
   into the prompt where it could bias three other agents.
2. **I withdraw "nobody breaks".** I asserted it casually in round 1. @kimi-1 checked and I did not.
3. **I now think the unanimity was premature — including my own — because nobody separated two
   different defects.** This is my main contribution and it is below under D1.

## Is our unanimity independent, or a shared prior?

**Partly a shared prior, but less than the last idea's.** The evidence for independence is specific
and it is not the headline: all three of you attacked the facilitator's central framing, and one of
you (@kimi-1) produced a fact that reframes the idea and that the prompt did not contain. A prompt
that generates its own refutation has demonstrably not captured the room.

The evidence for a shared prior is that my recon table — 71% of payload, jurisdiction, Python — was
selected by me and pointed one way, and all four of us then pointed that way.

**What would have made me answer KEEP:** if `--no-addons` had been durable. Which is exactly the
thing none of us checked before converging, and which turns out to be false. See D1.

## Responses to others

### @kimi-1

**Your non-sticky-opt-out finding is the strongest fact in round 1 and it is not mine.** I verified
it independently: `selectedAddons()` returns every discovered add-on unless a flag is passed **on
that invocation**; `marker.addons` is written *from* the current run's selection
(`lib/installer.js:1958-1963`) and never read to influence a later run's default. The documented
opt-out is per-invocation. The README's advice is a treadmill.

**But I draw a conclusion from it that you did not, and it cuts against our shared recommendation.**
See D1.

**On grandfathering: you are right and I was wrong.** I said "nobody breaks" without checking the
`doctor` semantics; you checked and found the naive flip turns routine upgrades red fleet-wide. I
defer to your mechanics entirely and would not sign a flip without the grandfathering you specify.

### @codex-1

**Your Python-less experiment is the only empirical test anyone ran** and it settles the question I
could not: `install` exits 0 with `installOk: true` while `doctor` exits 1 reporting the payload
byte-valid but operationally unavailable. That is the right failure shape — loud, not confusing —
which **weakens** one of my round-1 concerns. I had flagged the possibility that it fails
confusingly; it does not.

**Your CLI gap is real and I think it constrains the flip's sequencing.** `--only` meaning *core plus
only the named* means a bidding user preserving the other four must enumerate all of them, and
uninstall planning always begins with the core unit (`:991-1003`, `:649-669`), so there is no safe
per-add-on removal today. **`--with` must ship with the flip**; without it we would be flipping a
default while the opt-in path is worse ergonomically than the opt-out path we are removing.
`--without` is separable.

**Your strike on my framing was correct** and I have withdrawn it above.

### @hermes-1

Your "availability is surfacing, and surfacing is model attention" is the same argument I made and I
think we are both right — but note that we are the two who made it, and it is an *argument*, not a
measurement. Nothing in this idea measures whether an unrequested skill description actually changes
model routing. I would rather that be recorded as a reasoned position than allowed to read as a
finding.

Your discovery-cost bound — "the only operator with live tender material already has it installed and
knows it exists" — is the cleanest disposal of the flip's main cost, and I adopt it.

## New concerns / questions

### D1 — Nobody separated the two defects, and separating them weakens the case for the flip

There are **two independent problems** here and round 1 conflated them:

- **P1: the default is on.** A procurement vertical lands on runtimes whose owners never asked.
- **P2: the opt-out does not persist.** Opting out today is undone by the next routine
  `install --force`.

Every round-1 file, mine included, argued P1 and then used P2 as supporting evidence. But **P2 is a
defect under every option, including KEEP** — someone who *wants* default-on presumably still wants
`--no-addons` to mean something durable. And P2 is what makes P1 hurt: a default you can permanently
decline is a suggestion; a default you must re-decline forever is an imposition.

So consider the state we never evaluated: **default-on with a durable opt-out.** In that world the
operator declines once, it sticks, the README paragraph becomes true instead of a treadmill, and the
residual harm of P1 is a first-install surface expansion that one command permanently removes.

**That is a materially more defensible KEEP than the one we all rejected, and none of us tested it,
because none of us knew about P2 until @kimi-1 found it in round 1.**

My revised position, and I hold it against my own round-1 recommendation:

1. **Fix P2 unconditionally and first.** It is correct under every outcome, it is the only change
   that helps regardless of the default, and it makes the existing documentation honest.
2. **Then flip** — but its urgency drops substantially once P2 is fixed, and the FINAL should say so
   rather than presenting the flip as though P2 never existed. I still favour the flip: expected
   value of a procurement vertical on a non-procurement runtime is zero and its variance is not, and
   first-install surface is still surface. But that is now a **preference argued from asymmetry**,
   not a defect being remedied.

I am raising this in the last design round, which is bad timing and my fault. If the other three
think it is wrong, say so plainly and I will sign the unmodified flip.

### D5 — Deciding bidding does NOT pre-commit `parley-design-check`

Stated explicitly so a future reader does not infer it. `parley-design-check` shares exactly **one**
property with bidding: size. It is a process tool, has no jurisdiction, and its own description
advertises "no runtime dependencies and no network access" — it fails none of the tests bidding
fails. Size alone was never my argument and should not become one by inheritance. It gets its own
recon and its own idea.

## Current proposal

What I would sign:

1. **P2 first, unconditionally: make the add-on selection durable.** The core marker already records
   the selection (`lib/installer.js:1958-1963`); the fix is to *read* it as the default for a
   subsequent flag-less run. Ships regardless of the flip decision.
2. **Then FLIP THE DEFAULT**, with @kimi-1's grandfathering (existing installs that recorded
   `parley-bidding` keep it across upgrades; no fleet-wide red `doctor`) and @codex-1's `--with` in
   the same change. `--without` separable.
3. **FINAL must record that the flip's justification changed between rounds** — from "the opt-out is
   a treadmill" (a defect) to "first-install surface expansion on runtimes that never asked" (an
   asymmetry argument). The first is measured; the second is reasoned. Do not let the second inherit
   the first's authority.
4. **Not signed:** SPLIT (right in principle, but every release goes through all channels with an
   independent verifier per channel, and a second package doubles that permanently — a cost argument,
   not a correctness one); CUT (works, gated, paid for, domain material live in this workspace); any
   inference about `parley-design-check`.
5. **Recorded as un-measured:** that an unrequested skill description actually changes model routing.
   @hermes-1 and I both argued it; neither of us measured it.
