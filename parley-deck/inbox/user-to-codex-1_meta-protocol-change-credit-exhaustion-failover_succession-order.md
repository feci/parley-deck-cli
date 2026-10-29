---
from: user
to: codex-1
idea: meta-protocol-change-credit-exhaustion-failover
phase: round-01
blocking: no
date: 2026-10-29
---

## Owner ruling: the succession order (binding)

The owner gave this ruling at about 19:30 CET on 2026-10-29, before round 1. The relay recorded it verbatim:

> "no ked vypadne organizer, tak implementer ho preberie, ak vypadne implementer tak organizer ho preberie,
> ak vypadne aj organizer aj implementer, tak ho preberie prvy participant a participanti ked vypadnu, tak sa
> nedej nic, len proste vypadnu"

Translation: "Well, when the organizer drops out, the implementer takes it over. If the implementer drops out,
the organizer takes it over. If both the organizer and the implementer drop out, the first participant takes
it over. When participants drop out, nothing happens; they just drop out."

## What this settles

1. The organizer drops out: the implementer takes over the organizer role.
2. The implementer drops out: the organizer takes over the implementer role.
3. Both drop out: the first participant takes over both roles. "First" is the first remaining agent in
   `participants:` order that is neither the dropped organizer nor the dropped implementer.
4. A participant (not the organizer, not the implementer) drops out: no succession. It only drops out.

This replaces the "recommended default order" in brief section 3. FINAL must implement this ruling. The
participants decide only the details that the ruling does not settle. Examples:

- one agent that is both organizer and implementer (as in this idea) counts as "both drop out";
- the rule "a declared facilitator does not implement" yields to rule 2;
- the drafter roles and the reviewer count after a succession;
- the floor and the single-reviewer exception after a succession.
