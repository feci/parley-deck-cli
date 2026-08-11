---
idea: protocol-read-cost-regression
author: claude-1
created: 2026-08-10
participants: [claude-1, codex-1, hermes-1, kimi-1]
status: final
track: standard
---

## The report

The owner reports that Parley Deck has got **slower over the last few versions** and asks why, and
whether it can be optimized.

## What was measured before opening this idea

All figures are PRIMARY — reproduced from this repository. Token counts are bytes/4, a rough but
consistent proxy; treat ratios as reliable and absolutes as approximate.

**1. It is not the CLI.** Best of three, installed `parley 1.42.1`:

```
parley --version            0.00s      parley roster show        0.292s
parley --help               0.00s      parley protocol status    0.127s
parley agents list          0.08s      parley preflight -no-ping 0.074s
parley-deck-skill status --target all --project . --json        0.914s
```

Nothing here is a plausible source of a felt slowdown.

**2. The protocol body doubled in ten weeks.** `git log --follow` over
`internal/protocol/defaults/COOPERATION.md`:

```
2026-05-26    720 lines    49,379 bytes   ~12,300 tokens
2026-06-13    802           54,124         ~13,500
2026-06-19    959           71,714         ~17,900
2026-07-03   1138           88,726         ~22,200
2026-08-04   1306          100,490         ~25,100
2026-08-07   1359          104,479         ~26,100
```

**2.1× growth, monotonic.** Every version added; none removed.

**3. Every participant reads it in full, every round.** Measured on the deliberation that closed
two days ago (`ideas/protocol-overlay-local-extension`, 3 rounds, 4 agents):

```
round        bytes read/agent   ~tokens/agent   x4 agents
round-01              148,046          37,011      148,044
round-02              246,889          61,722      246,888
round-03              332,018          83,004      332,016

protocol alone = 71% of the round-1 read
round-3 read   = 2.24x the round-1 read
whole deliberation, read side only ≈ 727,000 tokens
```

**4. The existing lever does not pull on this.** §4.0's `track: fast | standard | deliberation`
scales rounds, reviewer count, fix-up cycles and timeouts. It does **not** scale what an agent
reads: a `fast` idea still loads the same ~26,000-token protocol. Conditional rigor was applied to
the *process* and never to the *context*.

## Working diagnosis — attack it

> The slowdown is not execution time; it is **read cost**. The protocol is 71% of what every agent
> loads before it can think, it doubled in ten weeks, and the per-round total compounds because each
> round also re-reads every prior round file. The felt effect is every round getting slower, which
> is what the owner reports.

I hold this diagnosis and I framed the measurements that support it, so **treat it as the weakest
claim here.** If it is wrong, say so and say what the real cause is.

## What round 1 must produce

Write `round-01/<agent-id>.md`. Answer:

1. **Is the diagnosis right?** Name a competing explanation and say what evidence would separate
   them. If you can measure something I did not, do.
2. **What is actually load-bearing per phase?** Most of the protocol is not needed to write a
   round-1 analysis. Which sections must an agent have loaded in Phase 1, in Phase 3, in Phase 6 —
   and which are reference an agent should open only on demand?
3. **What would you cut, and what must never be cut?** §15 verification integrity, §7 protocol
   change, §6 rule 3, and the §14 human brake exist because each was bought with a real failure.
   Reducing read cost must not quietly reduce the rules that stop bad work.
4. **Mechanism.** Options include a phase-scoped generated view, a physical restructure into
   appendices (already a ratified follow-up, `protocol-restructure-appendices`), a per-track context
   budget, prompt-side scoping by the facilitator, or something none of us has proposed. Recommend
   one and say what it costs.
5. **Is the compounding re-read of prior rounds separable?** Round 3 costs 2.24× round 1 largely
   because every prior round file is re-read. Is that necessary, or is a digest enough — and what
   breaks if an agent does not read a peer's full text?

## Note on how this idea is being run

This prompt deliberately does **not** ask you to read `COOPERATION.md` in full. You may open any
part you need, and you should say if that choice cost you something. Running the investigation the
expensive way would have been its own answer.

Track is `standard`, not `deliberation`, for the same reason. If the outcome turns out to change
protocol text, that escalates to a meta-protocol-change idea under §7 — this idea diagnoses and
recommends.

## Constraints

- Read-only: no edits, no git write commands.
- English only. Redact obvious secrets.
- §15 provenance: `PRIMARY` needs a stable locator or quoted command output; untagged reads as
  `RECALL`. Do not present my measurements above as your own verification — they are mine, and
  re-running one is more useful than repeating it.
