---
from: user
to: all
idea: roster-membership-overlay
topic: uniform-inheritance-path-c
date: 2026-08-19
binding: yes
---

# Owner ruling — the model is class inheritance, and membership is just another property

The owner rejected the framing of both (a) and (c) and supplied a third, **Path C**. Quoted
verbatim, then restated; the quote governs where the restatement is unclear.

> "neviem preco to komplikujes a rozdelujes vyznam feature a membership, proste sa na to pozri ako
> pri objektovom programovani je pretazenie. Global config je Rodicovska trieda, od ktorej moze
> dedit nejaka session svoj configuracnu triedu a tam pretazit hociktoru hodnotu z rodicovskej. Ale
> parley-deck si samozrejme najprv nacita rodica, potom z potomka pretazi rodicovske hodnoty a
> dostane vysledny config. To co pisem je as cesta C, ale je to trocha podobne ceste A. Akurat by
> som mal mat moznost pretazit nie len hociktoru roster premennu ale aj membership a vobec hocijaku
> configuraciu je v global."

## The model

```
class GlobalConfig:            # ~/.parley/agents.toml
    members = [claude-1, codex-1, hermes-1, kimi-1, opencode-1, zcode-1]
    kimi_1_speed = "deep"
    kimi_1_model = "kimi-code/k3"
    ...

class DeckConfig(GlobalConfig): # parley-deck/agents.toml
    kimi_1_speed = "fast"       # override ONE property
```

Resolution is the ordinary one: **load the parent, apply the child's declared overrides property by
property, and the result is the effective config.** A property the child does not mention is
inherited. Nothing else happens.

`DeckConfig` above has six members, because it never overrode `members`.

**Membership is not an authority. It is a property**, and it is overridable like any other:

```
class DeckConfig(GlobalConfig):
    members = [claude-1, codex-1]     # explicit override — now the deck has two
```

And the owner's scope is wider than the roster: *"nie len hociktoru roster premennu ale aj
membership a vobec hocijaku configuraciu je v global"* — **any** global setting must be
overridable this way, not only roster fields.

## Why this dissolves the idea's central conflict

Today's resolver does something no object model would accept: **defining any property in the child
replaces the entire parent class.** `LoadRosterScoped` asks whether the deck file contains any
`[roster.*]` block and, if so, treats that file as the complete roster — without reading the block's
contents. That is D-A's root cause, and it is why

```toml
[roster.kimi-1]
speed = "fast"
```

leaves a six-member deck with one member.

Under Path C that stops being a special case needing a fix: the block overrides `speed`, membership
was never mentioned, so membership is inherited. **The (a)/(c) split was an argument about which
patch to apply to a resolution rule that the owner has now replaced.** Both sides were arguing
inside a premise that no longer holds.

## What the participants must now answer

Path C is binding as a **direction**. Its engineering is not settled, and these are open:

1. **The compatibility hazard is unchanged and is the hard part.** 37 fleet decks carry
   `[roster.*]` blocks holding values. Under C those blocks stop declaring membership, so **37
   quorums silently change** — the exact hazard @claude-1 raised against the overlay in round 1 and
   which @codex-1's design avoided by being opt-in. C cannot avoid it by being opt-in without
   ceasing to be the default model. **Say how you would migrate**, and whether a versioned marker
   (`schema = 2`) on the deck file is the answer, or something better.
2. **`members` as a list property.** In an object model, overriding a property REPLACES its value.
   So `members = [...]` in the child replaces the parent's list. Is that sufficient, or does the
   deck also need the `super.members + [x]` / `super.members - [x]` forms? Note this is the same
   question as @codex-1's `add`/`remove`, arrived at from the owner's model rather than from a new
   syntax — say whether that changes your answer.
3. **Scope beyond the roster.** *"hocijaku configuraciu"* — the sandbox, approval, timeout and
   isolated-home settings in `[agents.*]` blocks are also global-vs-deck. Does C apply to them
   uniformly, and do they already behave this way? **Verify rather than assume.**
4. **What C does NOT change.** D-A, D-B and D-C are still real defects and §1 of `consensus.md`
   remains unanimous. C changes what the correct end state is; it does not retire the three verbs
   that currently misreport their own effects.
5. **Does C weaken or strengthen your signoff?** If your position was (a) "keep the authority
   model", C rejects that model outright — say whether you accept the ruling and what you would
   build under it. If you believe C has a defect the owner has not seen, **say so plainly**; a
   binding direction is not a prohibition on stating a technical objection to it.

## Status of the signoffs already in flight

The Phase 3 signoffs were dispatched before this ruling and answer the (a)/(c) question. They are
**not void** — their §1 confirmations, D-C positions and attribution corrections all stand. Their
(a)/(c) answers become input to Path C rather than a decision.
