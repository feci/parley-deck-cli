---
idea: round-two-trigger
author: claude-1
created: 2026-09-01
participants: [claude-1, codex-1, hermes-1, kimi-1]
roles:
  claude-1: protocol fit — this idea IS allowed to propose protocol text; where exactly
  codex-1: skeptic — is a rule the right instrument, or is this a tooling/UX problem
  hermes-1: mechanism design — the concrete detector and what it does when it fires
  kimi-1: measurement — what makes the trigger falsifiable rather than ceremonial
status: final
track: deliberation
corpus-freeze: "2d17478, 2026-09-01T06:55:05Z — 80 idea dirs with >=1 round dir under parley-deck/ideas/"
---

## Problem / idea

**Who decides that a deliberation is finished, and should that person be allowed to?**

Today the facilitator decides whether to open `round-02`. The facilitator is also a participant,
also usually the `FINAL.md` drafter, and **bears the cost of the extra round** — minutes of wall
clock and a real bill. Nothing detects convergence, nothing challenges the decision, and the
decision is invisible in the record: an idea that closed after one round and an idea that *should
have* opened a second look identical afterwards.

Two closed ideas pointed here independently, from different directions, and neither designed it:

- `protocol-mutation-diversity` (2026-08-31) named this **the highest-value follow-up** after its
  measurements relocated the defect: divergence *inside* a round is not the problem.
- `openviking-context-structure` (2026-08-31) established the doctrine this idea must obey: a
  passing measurement does not authorise a default flip, and a CLI-only change that quietly alters
  deliberation semantics is worse than a protocol change, not better.

### The measurement, frozen (PRIMARY)

**Corpus frozen at commit `2d17478`, 2026-09-01T06:55:05Z.** This idea is excluded from its own
target set. (The freeze is not ceremony: in the previous idea three participants reported three
different single-round counts and **all three were true when taken**, because the running idea sat
inside the corpus it measured.)

Of **80** idea directories with at least one round directory, **28** closed after a single round:

| track | single-round ideas |
| --- | --- |
| `<none>` (predates the track field) | 19 |
| `standard` | 4 |
| **`deliberation`** | **4** |
| `fast` | **1** (`tui-editor-composer`) |

The falsification hypothesis — "these are just small `fast` ideas, closing them in one round is
correct" — was tested in the previous idea and **failed**. Exactly one is `fast`. Four are
`deliberation`, the highest-rigour track, closed with **zero** cross-review; two of those four are
protocol changes: `meta-protocol-change-devx-speed`, `protocol-restructure-appendices` (the others
are `track-aware-driver`, `parley-learn-playbooks`).

**Counter-evidence, and it is real:** of **141** round-02+ artifacts carrying the mandated
`## Position changes since prior round`, **141** have substantive content and only **23** explicitly
report no change. Cross-review demonstrably moves positions *when it runs*. That is the argument
that the extra round is worth opening — and also the argument that the current judgment is mostly
right, since 52 of 80 ideas did open one.

### What to decide

- **The detector.** What observable condition says "this closed too early"? §15.6(b)'s existing
  language is *"round 1 closes with no substantive disagreement"* — never machine-detected, and a
  judgment call by the conflicted party. Can it be made checkable? If not, say so plainly rather
  than shipping a rule nothing can check.
- **What fires.** Opening a full round is the expensive answer. Is there a cheaper one — a single
  assigned artifact, a `parley consult`, an inbox escalation to the owner, a recorded decision?
- **Where it lives.** This idea **is** allowed to propose protocol text; that is its purpose, and
  it is the one thing that would earn a core version 2.12.0. But protocol text is not automatically
  the right carrier: `roles:`, the CLI, or a recorded facilitator declaration may be better.
- **The conflict of interest.** Even a perfect detector still reports to the party paying the cost.
  Does anything change that, and is it worth changing?
- **How we would know it worked**, at corpus level (n=4 gives 6 pairwise distances per round; per-idea
  criteria are not readable).

## Constraints

- **A protocol change is a new core version, published attended by the owner.** Do not edit
  `COOPERATION.md` directly. **Note the standing blocker: core 2.11.0 is staged but NOT published;
  the installed core is 2.10.0 and 40 of 41 decks lack §15.6 in prose.** Any proposal must say how
  it behaves on decks that do not have the current version.
- **The carrier thesis applies** (ratified in `protocol-generation-bias`): rules carried by prompt
  templates plus a validator reach near-universal compliance; identical rules in prose run in single
  digits. §15.6's own preamble says only clause (a) is machine-validated **and says so explicitly
  rather than implying a gate that does not exist**. A new unvalidatable duty will decay the same way.
- **Do not raise the floor cost of a `fast` two-round idea.** Track-aware or trigger-gated.
- **Zero new runtime dependencies.**
- **Determinism and auditability.** If the trigger fires, the record must show it fired and what
  happened; if it does not fire, the record must show it was evaluated.
- **§15 applies:** tag claims PRIMARY / SECONDARY / RECALL. Untagged is RECALL. Verify locators —
  a locator is a factual claim (openviking C8).
- English only under `parley-deck/`.

## Non-goals

- Re-litigating mutation, randomness, or the semantic-donor benchmark. That is closed
  (`protocol-mutation-diversity` FINAL, D1–D10). This idea is about the **phase transition**.
- Restoring the removed steelman clause wholesale. It was explicitly not restored; if a fragment
  of it belongs in the trigger, argue that fragment on its own merits.
- Cost measurement in tokens or currency — headless runners emit no provider input-token telemetry
  (`internal/driver/loop.go:174-175`). Call counts and wall time only.
- Automating the owner out of ratification.
