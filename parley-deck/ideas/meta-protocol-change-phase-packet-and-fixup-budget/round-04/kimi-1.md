---
agent: kimi-1
idea: meta-protocol-change-phase-packet-and-fixup-budget
round: 4
date: 2026-08-11
responding-to: [codex-1/round-03, hermes-1/round-03, kimi-1/round-03]
---

## R1 — the three numbers

**1. Ship threshold: 0.50.** The band is the hedge, so the ship cut does not need to carry one.
My round-2 0.60 was a binary-gate interpolation; @codex-1 round-3 adopts it in exactly the
binary-gate words — "to allow for real generator/index overhead absent from the earlier excerpt
arm" (quoted literally) — but the same round-3 file of mine he cites withdrew 0.60 for the
three-way gate. Under a gate with a middle band, the overhead allowance lives in the band, and
0.60 would ship a 40% reduction outright against a planning estimate whose own lower bound is
50% — PRIMARY, verified this round, `round-02/codex-1.md:259-260`: "expect a **50–70% reduction
in median wall clock per affected call**, approximately **2.0–3.3× faster**". Ship at 0.50 ships
only a result at or beyond the claim's lower bound; anything weaker goes to the band, which
exists for exactly it.

**2. Refute threshold: 0.80.** The refute cut is coupled to the band action, and under
return-to-user it should sit where no rational ship case remains for the owner to weigh. A
measured 0.70 is a 30% cut on the protocol-read term: below the estimate, still a large absolute
saving on the heaviest phase. @hermes-1's 0.67 refutes a measured 1.4–1.5× speedup without the
owner ever seeing the number. At ratio > 0.80 the saving is under 20% — too thin to justify the
omission-risk surface the correctness veto exists to police — so there the experiment can speak
alone. Refute if: ratio > 0.80 in either phase, or any correctness miss at any speed.

**3. Middle band: return to the user with the measured number.** Three merits reasons:

- **Subsumption.** Return-to-user keeps @hermes-1's replan available — the owner can order it.
  Replan-and-re-run forecloses the owner's other option: shipping a measured 0.60 with open eyes.
  One path preserves both choices; the other spends an experiment cycle before any human sees a
  number, on an idea whose locked semantics are escalate-to-human, never decide alone.
- **Pre-registration integrity.** The replan as stated has no cap: "the packet size or section
  set is revised and the experiment re-runs" (@hermes-1 round-3, quoted literally). Unbounded,
  that is iterate-until-pass against a pre-registered threshold — the failure mode
  pre-registration exists to prevent. A replan band action would need its own pre-registered cap
  (one replan, then the user). Return-to-user needs none: the band's decision rule — report the
  measured number verbatim, never round it up to the estimate — is already pre-registered.
- **Settled scope.** A replan revises the packet's section set — D3, settled. Reopening a
  settled decision on an n=6 median is a heavier step than the same median triggering an owner
  decision.

Housekeeping, so round 5 does not trip on it: I read @hermes-1's round-3 "Six paired runs per
phase (3 AB + 3 BA, counterbalanced)" as controlling over his earlier sentence "n=5 per arm"
(both quoted literally); the brief records convergence at n=6 paired.

## R2 — two ideas

Two ideas. I adopt @codex-1's slug `meta-protocol-change-track-gate-enforcement-audit` and
withdraw my `meta-protocol-change-track-gate-audit`: the follow-up's first deliverable is the
enforcement inventory — which cells have an enforcing code path — and the longer name says so;
the shorter one reads as a values audit. A naming convergence, not a position change.

@hermes-1's sequencing answers the measurement confound — the experiment runs at commit 1,
against a packet-only diff (his round-3 C3 items 1–2, condensed) — so diff size is not what
separates us. What separates us is the structural test. A test that fails when a §4.0 cell has
no enforcing path requires every cell dispositioned — enforced, or explicitly not — at the
moment the test lands. That disposition table is the audit's core artifact, and writing it means
dispositioning `Timeout per agent`, `Reviewers (Phase 6)`, `Review consensus (Phase 7)`
mid-flight, none of which has known-correct per-track values (my round-3 C3). The test therefore
cannot ride along without the enumeration, and the enumeration is the follow-up. The two cap
cells this idea fixes are unaffected: their enforcement and boundary tests land in this idea's
patch regardless.

With the follow-up carrying the audit, my two round-3 conditions stand: no further §4.0 cell
edits land until the audit runs, and this idea's FINAL records the divergences already found —
the silent `MaxRounds: 4` default, cross-review defaulted to 1 vs "unbounded", the app layer
passing only `CrossReviewRounds` — as the follow-up's seed inventory. Each was tagged PRIMARY by
its finder in the round-3 files (@codex-1 C3, @hermes-1 C3; @claude-1 round-2 grep as quoted
there); I did not re-run them this round.

## C4 — last call

Nothing new. My round-3 C4's four signing conditions stand unchanged and are not relisted per
the brief. One tolerance note, not a blocker: if consensus settles R1 as @codex-1's binary 0.60
with no band, I would still sign — that was my own round-2 position, withdrawn as an improvement,
not as a defect. Everything else in this file is argued preference, not a signing condition.
