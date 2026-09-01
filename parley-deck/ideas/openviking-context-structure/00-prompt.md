---
idea: openviking-context-structure
author: user
created: 2026-08-31
participants: [claude-1, codex-1, hermes-1, kimi-1]
roles:
  claude-1: protocol fit — where this would land in COOPERATION.md and the CLI
  codex-1: cost/benefit skeptic — is this a problem we actually have?
  hermes-1: integration reality — it ships an OpenViking client already
  kimi-1: evidence and measurement — what would we measure before adopting
status: final
track: standard
---

## Problem / idea

The owner asks: **look at the OpenViking project, decide whether the way it structures data is
interesting enough to inspire Parley Deck, and if so, what we should take.**

OpenViking is an open-source "context database for AI agents" that unifies memory, resources and
skills behind a **filesystem paradigm** with a `viking://` URI scheme and **hierarchical context
tiers** — a ~100-token `.abstract.md` (L0) and ~2,000-token `.overview.md` (L1) sidecar per
directory, with the full document (L2) loaded only on demand.

All the external material you need is copied into
`reference/openviking-notes.md` in this idea directory, with provenance tags. **Read it first.**
Do not assume you can browse the web; several participants cannot. If you do verify something
independently, say so and tag it PRIMARY.

This is an **evaluation**, not a foregone adoption. "Interesting but we should not adopt it" is a
fully acceptable answer, and so is "adopt exactly one narrow piece". What is NOT acceptable is a
verdict without a named mechanism and a named cost.

### Why this is being asked now — the concrete pressure

Parley Deck is already a filesystem context store. It has the same shape OpenViking formalizes:

```
parley-deck/
├── COOPERATION.md              # generated view of ~/.parley/protocol/core/<version>/
├── agents.toml                 # roster authority
├── inbox/<from>-to-<to>_<slug>_<topic>.md
└── ideas/<slug>/
    ├── 00-prompt.md
    ├── round-NN/<agent-id>.md
    ├── consensus.md
    ├── FINAL.md
    ├── IMPLEMENTATION.md
    └── review/round-NN/<agent-id>.md
```

Three pressures are real and measurable in this repository today — do not take them on trust,
they are checkable:

1. **Round-N prompt bloat.** Every cross-review round must carry the prior round's artifacts into
   every participant's prompt. The `protocol-generation-bias` idea that closed on 2026-08-29
   produced **22 artifacts totalling roughly 491 kB**. Cost grows with rounds x participants.
   L0/L1 tiering is aimed squarely at this.
2. **No cross-idea recall.** There are **41 decks** on this machine. There is no mechanism to ask
   "has any deck already deliberated this?" A prior idea's `reference/research-brief.md` is
   reachable only by knowing its path. `parley retro` and `parley learn` exist — establish what
   they actually cover before claiming a gap.
3. **No stable addressing.** Artifacts are referenced by ad-hoc relative paths that break when an
   idea is renamed or a deck moves. OpenViking assigns every file a deterministic id
   (`md5(account_id:uri)`).

### The question to answer

For each candidate idea you would take from OpenViking, state:

- **the mechanism** — what file, format or command changes, concretely;
- **what it costs** — who generates the summaries, when, and with what token/latency budget
  (an `.abstract.md` is not free: something must write it);
- **how we would know it worked** — the measurement, before adoption;
- **the failure mode** — specifically, what goes wrong when a summary is *stale or wrong* and an
  agent acts on the summary instead of the artifact.

## Constraints

- **Zero new runtime dependencies is strongly preferred.** Parley Deck is a single Go binary plus a
  Node installer. A design that requires running an OpenViking server, a vector DB, or a network
  service is very likely out of scope — if you propose one, justify it against a no-dependency
  alternative that reaches most of the benefit.
- **Do not propose editing `COOPERATION.md` directly.** It is a generated view of a write-once core
  (§ protocol). Protocol changes are a new core version, published attended by the owner.
- **The deck must keep working with plain `cat`, `grep` and `git`.** A structure only a tool can
  read is a regression: the audit trail's value is that a human can read it years later.
- **Backwards compatibility.** 41 existing decks and every closed idea must remain valid and
  readable. A migration that rewrites history is not acceptable.
- **English only** for every file under `parley-deck/`.
- Respect §15: tag claims PRIMARY / SECONDARY / RECALL. An untagged claim is treated as RECALL.
  Do not state a verdict about OpenViking's internals that the notes file does not support.

## Non-goals

- Building or vendoring any part of OpenViking.
- Replacing the round/consensus/review protocol itself. This is about how artifacts are
  **stored, addressed and summarized**, not about how agents deliberate.
- Solving cross-machine or multi-user sync. Single machine, local filesystem.
- Anything that requires the network at round time.
