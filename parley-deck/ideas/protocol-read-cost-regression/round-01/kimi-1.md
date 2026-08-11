---
agent: kimi-1
idea: protocol-read-cost-regression
round: 1
date: 2026-08-10
---

## Summary

**The diagnosis is half right, and the half that is wrong is the half that matters for the
owner's question.** The owner asked why Parley Deck *feels slower*. claude-1's answer is "read
cost". My verdict: protocol growth is real and verified, but (a) the compounding term in the
round-by-round numbers is peer round files, not the protocol; (b) there is no wall-clock
evidence that read volume is what the owner feels; and (c) at least three other variables moved
in the same ten-week window and none were measured. The diagnosis is a measurable proxy standing
in for an unmeasured cause.

What I verified myself (PRIMARY, my re-runs, not restatements):

- **Protocol growth confirmed.** `git log --follow` over `internal/protocol/defaults/COOPERATION.md`:
  49,380 B / 721 L at 2026-05-26 (rev `697cb66`) → 104,480 B / 1,360 L at 2026-08-07 (rev
  `8ed3c4b`). 2.12×, matching claude-1's table within ±1 line. What his table does not show:
  growth is **bursty, not steady** — 4 commits on 06-19, 3 on 06-24, 5 on 07-03/04, 8 between
  08-04 and 08-07 — and the ideas directory holds ~15 `meta-protocol-change-*` ideas
  (PRIMARY: `ls parley-deck/ideas/`). The protocol grows because the retro→rule pipeline works
  as designed and nothing in it ever subtracts. Accretion is the policy, not drift.
- **Obligations grew faster than prose.** `grep -c 'MUST'`: 15 → 37 (2.5×); `'MUST NOT'`:
  6 → 15 (2.5×) over the same revisions. The number of *mandated behaviors per idea* outgrew
  the byte count.
- **The round-3 compounding decomposes to peer files.** In `protocol-overlay-local-extension`
  (PRIMARY: `wc -c` over its round files): round-01 files total 98,843 B, round-02 files 85,129 B.
  An agent opening round 3 re-reads ~184 KB of peer text against 104.5 KB of protocol. The
  protocol re-read is a *constant* across rounds; **the entire 2.24× growth term is prior-round
  files**, and those are agent-asymmetric: hermes-1 writes ~34 KB/round, kimi-1 ~29–32 KB,
  claude-1 ~10 KB, codex-1 ~12–22 KB. claude-1's framing — "the per-round total compounds
  because each round also re-reads every prior round file" — names the right mechanism but the
  wrong culprit: the protocol is the fixed cost, verbosity is the compounding cost.
- **No wall-clock signal in the data he cited.** File mtimes for that same deliberation
  (PRIMARY: `stat`): round-1 window 20:53→21:08 (15 min), round-2 21:42→21:55 (13 min),
  round-3 12:44→13:04 (20 min). Read volume grew 2.24× across those rounds; completion windows
  did not. The deliberation spanned ~43 hours (08-07 20:53 → consensus 08-09 15:42), dominated
  by between-round gaps, not in-round work. And runner telemetry can't settle it: the newest
  `parley-deck/runs/` event log is 2026-06-02 (`duration_ms` 17,874 / 33,668 for a round-1 back
  then) — **no duration data exists for the period of the reported slowdown.**
- **Confounders moved in the same window** (PRIMARY: `git log -- parley-deck/agents.toml`):
  hermes pinned to GLM 5.2 (06-20), kimi-1 activated with pinned models (07-28), antigravity-1
  reactivated as a *fifth* participant (07-30). Model swaps change per-call latency; a fifth
  participant adds a writer and a reader to every round and makes each round complete at the
  slowest of five.

**Competing explanation:** the felt slowdown is write-side and process-side, not read-side.
Input tokens are processed in parallel (and with prompt caching the unchanged protocol prefix is
near-free on re-read); output tokens are serial. A 30 KB round file (~7.5k output tokens) takes
longer to *generate* than 104 KB takes to *read*. Meanwhile MUST-density grew 2.5× (more duties
per idea → more artifacts, longer files, more cycles) and the roster grew. Read cost is the
easiest of these to measure post-hoc, which is why it became the diagnosis — a selection effect,
not an independent confirmation.

**Evidence that would separate them:** per-invocation telemetry for one `standard` idea —
input bytes, output bytes, duration per agent call. The runner already emits `duration_ms`
(`events.jsonl`); it just isn't being used for skill-facilitated deliberations. My hypothesis
predicts duration correlates with output bytes and round count, not input bytes. claude-1's
predicts duration tracks input bytes. One instrumented idea decides it.

**Cost of the restraint this prompt imposed** (required disclosure): I read the heading
skeleton, the Quickstart, §9, and the Phase 0–3 text — roughly 15–20 KB of the protocol's
104.5 KB. It cost me exactly one thing: my "never cut" endorsement of §6 rule 3 and §14 is
*inherited from the prompt's characterization*, not verified against their text, and my
provenance tagging follows the prompt's restatement of §15.2 rather than §15.2 itself. For a
diagnostic idea the loss was ~zero — which is itself a datapoint for phase-scoped reading. For a
verdict-producing idea I would have had to open §15.1–15.2 (8.4 KB, 8% of the file). Restraint
scaled fine; it would have scaled even with the sections the prompt calls untouchable.

## Proposed approach

### 1. What is load-bearing per phase (with bytes, PRIMARY: `sed -n` ranges of the 104,480 B file)

- **Phase 1 (independent analysis):** Quickstart, §0–§3 (transport header, scope/non-solo,
  roster, layout), §4.0 track, Phase 0–1 text (lines 262–316), §5, §6, §15 (Phase 1 cites it at
  line 299 — provenance duties bind from round 1). **Not needed:** Phases 4–8 (15,614 B —
  implementation/review/fix-up for work that doesn't exist yet), §9.0 (3,765 B, facilitator-only
  pre-flight), §11 (16,338 B; only the active transport's subsection matters, and for a
  `local-dir` project §11.B + §11.C = 13,245 B of pure dead weight), §12–§14 (~15 KB), Appendix A
  (1,963 B). Phase-1 need ≈ 40–45 KB, not 104.5 KB.
- **Phase 3 (consensus):** Phase 1 set minus Phase 2 templates, plus Phase 3 text and the close
  conditions it cites — §15.3, §15.5, §15.6 (lines 347–349) — plus Escalation (lines 677–709,
  2,166 B). Still not Phases 5–8, §11 inactives, §12–§14.
- **Phase 6 (code review):** the 15.6 KB of Phases 5–8 text *becomes* load-bearing here, plus
  §15.1 (verification ownership) and §7 only if the fix touches protocol. The design-phase
  templates (Phases 0–3) become skippable. Symmetric conclusion: **no phase needs more than
  ~half the file, and every phase needs a different half.**

### 2. What to cut, what never to cut

Never cut from the always-loaded set (the rules bought with failures): §15 (8,350 B), §7
(2,811 B), §6 (1,814 B), §14 (~3 KB), Escalation (2,166 B), the §1 non-solo requirement. Total
≈ 18–19 KB — under 19% of the file. The safety core is *cheap*; that is exactly why "we must
load everything to keep the rules" is a false dilemma.

Defer to on-demand (reference, opened when the task touches them): the two inactive transports
(13,245 B), §12 pipelines (~3.5 KB), §13 retro (~2.6 KB), §9.0 facilitator pre-flight (3,765 B —
facilitator-only), Appendix A (1,963 B), and the not-yet-relevant phase blocks (up to 15.6 KB).
Combined deferrable for a design-phase participant: **~55–60 KB of the 104.5 KB**.

### 3. Mechanism — recommendation

**Physically split core from reference via the already-ratified restructure, and flip the
read mandate to match.** Notably, the machinery and the doctrine both already exist:

- The protocol's own Quickstart already declares progressive disclosure: *"You do not need to
  read all of this… The core every participant needs is §0–§8. The rest are reference
  appendices — skip them until a task needs them"* (PRIMARY:
  `internal/protocol/defaults/COOPERATION.md:13-14, 35-38`).
- The launcher contradicts it: `SKILL.md:12` "Always read `parley-deck/COOPERATION.md` first";
  `SKILL.md:24` "Load the full cooperation protocol before acting"; `SKILL.md:116` and §9
  checklist item 1 (`COOPERATION.md:857`) repeat the unqualified mandate. **The regression is
  not missing machinery — it is a mandate that didn't age with the file.** The proof: this very
  idea's prompt had to carve out an explicit exemption ("Deliberately NOT required: reading
  COOPERATION.md in full"). A default that needs per-idea exemptions is the bug.
- The generation pipeline exists: `COOPERATION.md` is already "a generated view of a global
  core" (PRIMARY: `SKILL.md:304-307`), and §2 already became a generated view (commit
  `3b94f85`). `protocol-restructure-appendices` is already ratified. So: ship the deck view as
  **core (safety sections + active-phase text + active transport)**, move the rest to
  `parley-deck/protocol-ref/*.md` opened on demand, keep § numbers immutable across the split,
  and rewrite the three mandate lines to "load the core view; open reference sections on
  demand." The Quickstart's role/phase table (lines 28–38) becomes the discovery map.

I reject two of the listed options as primary levers: **per-track context budgets** — track
scales process rigor; a `fast` idea still needs §15, so context scoping is orthogonal to track;
and **pure prompt-side scoping** — the Quickstart doctrine proves self-description rots when the
mandate says otherwise. Prompt-side scoping is the *transitional* form of the fix (edit three
lines in SKILL.md + §9 first), not the end state.

**Cost:** one §7 meta-protocol-change idea (mandatory, and correct — this changes protocol
behavior); a discovery risk that an agent needs a deferred section and doesn't know it exists
(mitigated by the always-loaded safety core + the map); and cross-reference churn if anyone
renumbers (don't).

### 4. Prior-round re-read: separable, but bound the source first

Yes, separable — but the cheapest fix is upstream of the digest question entirely: **cap
round-file length in the Phase 1/2 templates.** The compounding term is peer text (§Summary),
the template's only length guidance is "concise restatement" (line 336), and "Address every
other active agent explicitly" (line 341) mandates O(n²) cross-reference text per round. A
bound (e.g. round files ≤ ~12 KB / digest-first structure, "cite line numbers, don't quote
blocks") attacks the growth term at its source.

On digests specifically: a **facilitator-written** digest is dangerous — it injects the
facilitator's framing into every response (the very anchoring round-1 independence exists to
prevent) and concentrates narrative power in one agent (§15.5/§15.6 exist because correlated
agreement was a real failure). A **self-authored, fixed-format digest block** appended to each
round file (position, open objections, DISPUTED flags verbatim, ≤ ~400 words) is safe: round
N+1 reads all digests + the full text only of files it will dispute. What breaks if an agent
never reads a peer's full text: quote-level rebuttal, detection of hedged or self-contradictory
claims, and §15.3 provenance tracing — a digest can launder a contested claim into a settled
one. The consensus section already distrusts summaries for this reason: *"raw round files are
never hidden behind the summary"* (line 385). So: **digest for routing attention, full text for
adjudication.**

## Concerns and open questions

- **The owner could be feeling something none of our numbers capture.** The only evidence that
  read cost dominates *felt* time is absent; my mtime analysis is suggestive but confounded
  (unknown per-agent start times, parallel launches, human gaps). Before any restructure ships,
  one instrumented idea should confirm the mechanism — otherwise we optimize the measurable term,
  not the felt one. This is my strongest disagreement with the prompt: it asks "what would you
  cut" before establishing that reading is what the owner feels.
- **The mandate lives in two places** (skill SKILL.md *and* §9 checklist). Fixing one without
  the other leaves the full-read default in force for whichever surface an agent entered
  through.
- **Verbosity may be agent character, not protocol ambiguity.** hermes-1 and kimi-1 (this file
  included) write 3× claude-1's bytes under identical instructions. A length cap will be
  experienced as a style constraint and needs framing as a *protocol* rule or it will be
  gamed/ignored.
- **Generated views add a sync failure mode.** If `protocol-ref/` and the deck view drift
  (the two existing COOPERATION.md copies already differ by 577 B of project zones — expected,
  but proof that copies drift), agents may cite stale reference text. The version.json sha256
  provenance mechanism (§9.0) needs to cover the split files.
- I did not independently verify §6 rule 3 or §14 text (restraint disclosure, §Summary); if the
  prompt mischaracterized them, my never-cut list inherits the error.

## Risks

- **Optimizing the wrong term (highest risk).** If the write-side/process-side hypothesis is
  right, cutting read cost ~in half yields little felt speedup and the restructure gets blamed.
  Mitigation: instrument one idea first; ship the round-file length bound together with the
  core/reference split so both terms move.
- **Unknown-unknowns from scoped reads.** An agent that never loads §13 won't run a retro when
  it should; one that skips §12 may mishandle pipeline work. Mitigation: safety core always
  loaded; Quickstart map enumerates what exists; facilitator's per-phase prompt names the
  reference sections relevant to that phase.
- **Digest laundering.** Self-authored digests can quietly upgrade a DISPUTED claim to settled.
  Mitigation: digests must carry §15 verdict states verbatim; any dispute forces full-text reads
  of the cited files.
- **Citation rot across the split.** Old ideas cite §-numbers; renumbering would orphan them.
  Mitigation: split physically, never renumber; new sections append.
- **Scope creep into a rewrite.** The correct deliverable is a move + three mandate lines + one
  template bound, executed as a §7 meta-protocol-change idea — not a prose overhaul. Every byte
  "improved" in passing is a new unreviewed rule.
