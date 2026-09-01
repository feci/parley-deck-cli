---
agent: claude-1
idea: openviking-context-structure
round: 1
date: 2026-08-31
---

## Summary

OpenViking is worth stealing from, but only one of its four ideas survives our constraints.

**Take the tiering. Reject the addressing.** The `.abstract.md` / `.overview.md` / full-document
ladder maps onto a cost we can measure in this repository today. The `viking://` URI scheme,
the deterministic file ids and the vector layer do not survive "must keep working with `cat`,
`grep` and `git`" plus "41 decks stay valid".

The honest headline is uncomfortable for my own proposal: the tiering saves **41.9%** of a
round-03 prompt, but **all of that saving comes from round-01**, because round-01 is the only
round whose file shape mandates a `## Summary`. Measured: **6 of 16** round artifacts in
`protocol-generation-bias` carry one. So the mechanism is real, cheap — and bounded in a way that
the OpenViking analogy hides.

**Correction to `00-prompt.md`:** it says `protocol-generation-bias` is "22 artifacts, roughly
491 kB". Measured today, `find … -name '*.md' | wc -c` gives **29 files, 586 kB**. The direction
of the argument is unchanged (it is larger, not smaller), but the number in the kickoff is low and
should not be re-quoted. Command:
`find protocol-generation-bias -name '*.md' -type f -exec cat {} + | wc -c`.

## Proposed approach

**One mechanism, staged, no new dependency.**

A round-N prompt currently carries every prior round in full. Change it to carry:

- round **N-1** in **full** — cross-review requires responding to specific claims, and a summary
  cannot support `### @<other-agent>` rebuttal;
- rounds **< N-1** as **summaries only**, with the full path printed beside each summary and an
  explicit instruction that the artifact is one `cat` away.

**Where the summary comes from — this is the part that costs nothing.** Round-01's required file
shape already mandates `## Summary`. Extracting it is `awk`, in Go, at prompt-build time. No agent
writes anything new, no model call, no file is generated, nothing to go stale — because nothing is
stored. This is the decisive difference from OpenViking: their `.abstract.md` is a **materialised**
sidecar that something must write and keep fresh; mine is **derived at read time** from a section
the protocol already requires.

Measured on `protocol-generation-bias`:

| | bytes |
| --- | --- |
| round-03 prompt today (r01 full + r02 full) | 220,115 |
| round-03 prompt tiered (r01 summaries + r02 full) | 127,735 |
| saving per participant | 92,380 (41.9%) |
| × 4 participants, one round | 369,520 |

**What I am NOT proposing, and why.** No `.abstract.md` files, no `parley-deck://` scheme, no
content index, no vector store, no cross-deck registry. Each of those is a materialised artifact
with a staleness failure mode, and the constraint list rules out the infrastructure that would keep
them fresh.

## Existing alternatives

The mechanisms my proposal would build by hand, and what already ships:

**1. Extracting a per-artifact summary at prompt-build time.**
- Built by hand: an `awk`-equivalent section extractor in Go, plus the round-N prompt assembly change.
- Already ships: **`parley learn <closed-idea-slug>`** (`parley --help`; run today on
  `protocol-generation-bias`). It writes `parley-deck/playbooks/<slug>.md` — **1,444 B distilled
  from 586 kB, a 415× compression**. This is the closest existing thing to an L1 tier and I nearly
  missed it.
  **But it does not cover this case, and the reason matters:** the emitted playbook is a
  **skeleton with placeholders**, not a distillation. Its mechanical fields are auto-derived
  (track, participants, fix-up cycles, lifecycle), while every judgment field is an instruction to
  a future reader — verbatim: `(Generalize this line: what class of task does this playbook
  cover? …)` and `(Fill from this idea's review consensus + IMPLEMENTATION.md deviations …)`.
  It is also **per closed idea**, not per round, and runs **after** an idea closes, so it cannot
  feed a round-03 prompt of the idea still running. Constraint-forced: yes — the whole point is
  in-flight rounds.

**2. Finding prior deliberation on a topic.**
- Built by hand: a content index across ideas/decks.
- Already ships: **`parley retro scan`**. Verified today — it returns `SCORE FAILURE-TYPE IDEA`
  over 29 closed ideas (`blocked-or-abandoned`, `escalation`, `fix-up-heavy`, `review-churn`).
  It indexes **process health, not content**: it can tell me `integrate-parley-bidding-addon`
  scored 105.0 as blocked-or-abandoned; it cannot tell me whether any idea has discussed context
  tiering. Not a substitute — but it means "no recall mechanism at all" (00-prompt, pressure 2)
  is **overstated** and I withdraw that framing.
- Also ships: `parley context repo-map` — maps the **code repo**, not the deck; wrong corpus.
- Also ships: `git log`/`git grep` across the deck — real, and genuinely covers "find the text";
  it does not summarize. Inherited, not constraint-forced.

**3. Stable addressing for artifacts.**
- Built by hand: an id scheme.
- Already ships: **git blob SHAs and paths**, plus `parley sessions inspect RUN_ID` and
  `parley consensus status IDEA`. OpenViking's `md5(account_id:uri)` buys deterministic identity
  across a rename; git buys content identity, which is strictly stronger for an audit trail.
  I judge this alternative **sufficient** and drop the addressing proposal entirely.

**4. Scope separation (`resources` / `user` / `agent`).**
- Already ships, under different names: `~/.parley/protocol/core/` (account-global),
  `parley-deck/` (project), `parley-deck/inbox/` (peer-to-peer), `agents.toml` (agent capability).
  The separation OpenViking formalizes, we already have. Nothing to take.

**Sources consulted for this section:** `parley --help`, `parley retro scan --k 3`,
`parley learn protocol-generation-bias`, `parley context repo-map --help`, the generated
`parley-deck/playbooks/protocol-generation-bias.md`, and `reference/openviking-notes.md`.

## Concerns / open questions

1. **The saving is smaller than it looks, and I want it challenged.** Round-02+ files use a
   different shape (`## Position changes since prior round`, `## Responses to others`,
   `## Current proposal`) with **no `## Summary`** — measured 0 B of extractable summary for
   round-02. So an idea that reaches round-04 saves on r01 and r02… but r02 has nothing to
   extract. Either later rounds gain a mandated summary section (a protocol change, i.e. a new
   core version) or the mechanism only ever pays off once. I lean toward the narrow version and
   accepting the smaller win.
2. **Is a 4-agent, ~40-idea deck even in OpenViking's problem class?** Their design targets
   large-corpus RAG with a vector layer. Adopting a structure sized for that is the classic
   over-build. I want at least one participant to argue the null: change nothing.
3. `hermes-1` ships a working `OpenVikingMemoryProvider` (see notes §7). That makes an experiment
   cheap — but availability is not a reason to adopt, and I would like hermes to say explicitly
   whether its client would need a running server.

## Risks

- **Summary-driven error.** An agent rebuts a position it only saw in summary and misses a
  load-bearing qualifier. This is the real failure mode and it is not hypothetical: this very
  session shipped a rule whose defect was that *a printed rule binds only where enforcement lives*.
  Mitigation: full text for round N-1 (where rebuttal happens), path printed beside every summary,
  and the prompt stating that summaries are lossy.
- **Derived-vs-stored is the whole safety argument, and it is load-bearing.** If a later
  implementation quietly materialises `.abstract.md` files "for speed", every staleness failure
  mode I excluded comes back. That should be an explicit non-goal in `FINAL.md`.
- **Protocol coupling.** Extraction depends on the round-01 file shape. If a future core version
  renames `## Summary`, prompt assembly silently degrades to carrying nothing. Needs a test that
  fails loudly, not a fallback that hides it.
- **Scope creep toward a context database.** Every step from "derive a summary" to "index the
  fleet" is individually reasonable. `FINAL.md` should draw the line explicitly.
