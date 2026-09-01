---
idea: protocol-mutation-diversity
status: final
drafted-by: claude-1
date: 2026-08-31
rounds: 2
participants: [claude-1, codex-1, hermes-1, kimi-1]
signoffs: accept=4 reservations=0 block=0
corpus-freeze-required: true
---

# Verdict

**The diagnosis is right. Mutation-as-randomness does not transfer. One stochastic mechanism
survives and is worth one bounded experiment.**

- **Premature convergence is real, measured, and worst where it should be least.**
- **Blind randomness is rejected** — on the literature and, independently, on this deck's budget.
- **One genuine stochastic operator survives:** drawing *semantic material* at random from a
  disjoint closed idea. It is closer to GA crossover than any curated scheme proposed here.
- **Nothing changes in production.** No core version, no runner change, no new artifact class.

# The finding the verdict rests on

Of **80** idea directories, **28** closed after a single round — no cross-review at all. By track:

| track | single-round ideas |
| --- | --- |
| `<none>` (predates the track field) | 19 |
| `standard` | 4 |
| **`deliberation`** | **4** |
| `fast` | **1** (`tui-editor-composer`) |

The gate on this entire idea was a falsification hypothesis the drafter proposed against his own
position in round-02: *"these are just small `fast` ideas, and closing them in one round is
correct."* **It failed.** Exactly one is `fast`. Four are `deliberation` — the highest-rigour track
— closed with **zero** cross-review, and two of those four are protocol changes:
`meta-protocol-change-devx-speed`, `protocol-restructure-appendices`, plus `track-aware-driver` and
`parley-learn-playbooks`.

**Independently recounted by codex-1, hermes-1 and kimi-1 at signoff; all three reproduced
28 of 80 and the 19/4/4/1 split.** The drafter had this number wrong twice before, so it was
verified three times.

**Where the defect is not:** of **141** round-02+ artifacts carrying the mandated
`## Position changes since prior round`, **141 (100%)** have substantive content and **23**
explicitly report no change. Cross-review moves positions when it runs. *Instrument limits:* the
"no change" figure is a keyword match, so **23 is a lower bound** and ~84% movement is an **upper
bound**; it screens for movement, not diversity, and cannot detect four participants moving
together toward a wrong answer — which is §15.6(b)'s concern.

**So the defect is the round-02 trigger, not intra-round divergence.**

# The protocol gap (real, deliberately left open)

Verified verbatim. Current §15.6(a) covers "the mechanisms the proposal builds by hand, and for each
what the toolchain **already ships**" — an implementation-reuse duty about *tooling*. The clause
removed on 2026-08-29 covered "the strongest **rejected or unconsidered alternative** ...
steelmanned, with its best supporting evidence and an observation that would change the
recommendation" — a deliberation duty about *competing positions*. **Different objects.** §15.6(b)
carries only a disconfirmation fragment; it develops no alternative and assigns nobody.

The gap is real. **It does not follow that filling it improves outcomes**, and it is not filled
here.

# Binding decisions

**D1. No production change.** No core version, no runner change, no new artifact class, no
`adversarial.md`, no change to `consensus.md` or `FINAL.md` semantics. The removed clause is **not**
restored as protocol text.

**D2. Blind randomness rejected**, for two independent reasons: for LLMs mutation is "a learned
proposal operator with **semantic priors**" and "random mutations [are] **inefficient** in large
solution spaces"; and the GA case for randomness rests on many cheap trials while we have ~4
expensive ones per round. **This is an argument from cost, not principle** — at 400 cheap
participants the owner's original framing would likely beat every structured scheme proposed here.

**D3. Carrier is the existing `roles:` field**, never the protocol: per-idea, advisory by
construction, cannot change quorum or signoff weight, reversible, no core version.

**D4. One sealed benchmark, three arms.**
1. **Control** — ordinary advisory role.
2. **Structured reframe** — one scheduled transform (`boundary`, `mechanism`, `representation`,
   `objective`); emits an alternative only if endorsable, else null.
3. **Stochastic semantic donor** — a **recorded seed** selects, without replacement, one bounded
   mechanism/constraint/test tuple from a **disjoint frozen idea**; the participant asks whether
   that donor, its inverse, or its abstraction transfers. **Randomness selects semantic material —
   not tokens, not temperature, and never a position the agent is ordered to defend.**

Budget: 12 targets → 36 generation calls + 2 blind evaluators × 12 → 24 = **60 deep calls**;
repeated once on 12 held-out targets only if the first batch passes; **hard ceiling 120 calls**.

**D5. Endorsability is mandatory.** An arm emits an alternative **only if the participant can
endorse it**, else null, and null is a finding. No participant ever argues a position it does not
hold. This is the agreed guard against manufactured dissent.

**D6. Cost measurement is blocked, and says so.** Headless runners emit no provider input-token
telemetry (`internal/driver/loop.go:174-175`). The benchmark reports **call counts and wall time**
and must not report currency or token savings.

**D7. No success criterion may be phrased per idea.** With n=4 there are 6 pairwise distances per
round; per-idea diversity is statistically unreadable. Every criterion is corpus-level.

**D8. The corpus must be frozen at a stated commit before the experiment starts**, and the running
idea must be excluded from its own target set — see the methodological finding. *Carried from
hermes-1's signoff reservation: the 28-of-80 figure is valid only against that stated freeze.*

**D9. Mixed-outcome tie-break** *(carried from kimi-1's signoff)*: both "gap exists, tool does not
work" and "tool works, but the trigger is the real problem" close **without** a core change. Only a
result that is positive on the tool *and* points at production changes anything, and even then via
a separate protocol-change idea.

**D10. §15.6(b) applied to ourselves.** All four answered a request for *randomness* with
*structure*, and three of four opened with "no new mechanism". That is a shared prior, not
evidence. What would have made it wrong is stated above and nearly happened: had the 28 skewed
`fast`, our collective "no mechanism needed" would have been the very premature convergence the
owner suspected.

# Methodological finding

**The corpus we measured contained the idea doing the measuring.** Three participants reported
three different single-round counts — 29, 52 cumulative, and 28 — and **all three were true when
taken**: this idea sat in the single-round bucket until it opened round-02, then moved. No glob was
wrong. This is why D8 is binding rather than advisory.

# Unresolved — the FINAL drafter is required to leave this open

**A sealed benchmark may not be able to answer this question.** A reframe's value is whether it
changes a *live* deliberation's outcome; replaying against closed ideas measures whether it produces
*different text*, not a *better decision*. codex-1's blind two-evaluator adjudication is the
proposed mitigation, and it may be insufficient. The drafter raised this against the position he
had just conceded to and has no clean counter-proposal. All three signoffs explicitly required that
this not be presented as settled. **It is not settled.**

# Rejected

| Rejected | Reason |
| --- | --- |
| Temperature / sampling variance as the operator | D2; and no configured variance exists today to baseline against |
| Forced advocacy of an unheld position | manufactured dissent; superseded by D5 |
| Crossover between participants' artifacts | no defined splice exists; D4 arm 3 is corpus-to-participant, not participant-to-participant |
| A population, generations, or a fitness function | requires cheap automatic fitness we do not have |
| Restoring the removed steelman clause as protocol text | re-creates an unvalidatable judgment duty; §15.6's own preamble warns of exactly this decay |
| A new `adversarial.md` artifact class | adds a file that never reaches consensus and duplicates existing carriers |
| Any per-idea diversity metric | D7 |

# Deferred follow-ups (named, unowned)

- **`round-two-trigger`** — what the measurements actually point at: what opens round-02, given the
  facilitator is also a participant and bears the cost of the extra round. **This is the highest-value
  follow-up from this idea and it is not designed here.**
- **`parley-context-telemetry`** — carried over; prerequisite for any cost claim (D6).
- **Drawn-axis operator** — fallback only, if D4 shows the donor arm works and the structured arm
  does not.

# Corrections logged

| # | Correction | Found by |
| --- | --- | --- |
| M1 | kickoff said 2 ideas reached a fourth round; it is 1 — a `ls -d round-*` glob matched a stray file | kimi-1 → claude-1 conceded |
| M2 | hermes-1's claim that §15.6(a) carries the removed steelman's function — wrong on the object | claude-1 → hermes-1 withdrew |
| M3 | claude-1's round-01 mechanism withdrawn — the instrument he built to satisfy kimi-1 showed the defect is the trigger, not intra-round divergence | claude-1 (own) |
| M4 | claude-1's falsification hypothesis tested and **failed** | claude-1 (own) |
| M5 | 29 vs 52 vs 28 — not an error; the corpus changed because this idea was in it | claude-1 |

**M3 and M4 are the transferable ones.** The drafter's proposal was destroyed by an instrument he
built to satisfy a critic, and his replacement premise survived only because the data went the
other way. Neither changed by argument. That is the mechanism this idea was asking for, and it
already works — it just is not called mutation.
