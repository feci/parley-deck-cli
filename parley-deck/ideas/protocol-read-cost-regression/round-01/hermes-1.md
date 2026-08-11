---
agent: hermes-1
idea: protocol-read-cost-regression
round: 1
date: 2026-08-10
---

## Summary

The diagnosis is directionally right — the protocol grew 2.11× in ten weeks and every agent reads it in full — but it is incomplete in one load-bearing way: the CLI does not put the protocol into the prompt. The protocol enters context because the **parley-deck skill instructs agents to read it** ("Always read `parley-deck/COOPERATION.md` first", `SKILL.md:12`) and §9 reinforces that ("Read `parley-deck/COOPERATION.md`", `COOPERATION.md:857`). The CLI's `BuildRoundOnePrompt` constructs the prompt from `00-prompt.md` only (`runner.go:822`); COOPERATION.md is never embedded by the runtime. This distinction matters because it means the read-cost lever is not in the CLI prompt builder — it is in the skill's loading instruction and in the §9 checklist, neither of which the diagnosis names. The competing explanation is not that the diagnosis is wrong, but that it has located the cause one layer off: the felt slowdown is a **skill-loading protocol** problem, and the protocol's growth is the multiplier, not the mechanism.

My independent measurements confirm the growth and the per-round re-read. I did not replicate the exact absolute byte counts at every checkpoint (I measured 60,825 bytes at `efe76d0` 2026-06-13 vs the stated 54,124, and 68,244 at `8477ea6` 2026-06-19 vs 71,714 — likely different commit SHAs for the same dates), but the trend and the 2.1× ratio hold: `49,380 → 104,480 bytes`, `721 → 1360 lines`, monotonic (PRIMARY, `git show <sha>:internal/protocol/defaults/COOPERATION.md | wc -c`).

## Proposed approach

**Recommend a phase-scoped generated view, built on the core-protocol machinery already ratified in `meta-protocol-change-global-core-protocol`.**

That idea (status `final`, `FINAL.md` D1–D12) already moves COOPERATION.md to a generated view from an immutable, versioned, content-addressed core under `~/.parley/protocol/core/<version>/`, with a registry of permanent block IDs. The render verb (`parley protocol render`) is the natural insertion point for a `--scope <phase>` flag that emits only the blocks the phase requires. The cost is low because the block-ID registry (D2) already exists as a ratified design — scoping is a read-filter over it, not a new architecture.

What this gives:
- Phase 1 agent receives §1 (scope), §3 (layout — the artifact paths), §4.0–Phase 1 (the round-1 rules + file shape), §6 rules 1–3 and 6 (conflict avoidance, English-only), §15.2 (provenance), and the Quickstart. That is ~8–10K bytes, not 104K.
- Phase 6 reviewer receives §4 Phase 6 + §15 in full + the review brief rules. The full §15 (8,432 bytes, the third-largest section) is only load-bearing during review — it is dead weight in Phase 1.
- §11 (transport mechanics, 16,311 bytes — the second-largest section) is reference material: an agent in transport A (local-dir) does not need §11.B (GitHub PRs) or §11.C (GitLab MRs). These are on-demand.

**What it costs:**
1. The block-ID registry must exist and be maintained (already ratified as D2, not yet built — `FINAL.md` says rank 1 ships `render`, `check`, and the core store; the registry is part of that).
2. A phase→blocks mapping table, owned by the protocol, that says which blocks each phase needs. This is new protocol text — a §7 change. But it is additive and small.
3. A risk that scoping hides a rule an agent needed but did not load. Mitigated by keeping the full protocol one command away (`parley protocol render` with no `--scope`) and by the §9 checklist telling agents they may open more.

This composes with the ratified `protocol-restructure-appendices` (physical reorg, already done) rather than competing with it — the appendices reorg made core-first ordering, and scoping is the runtime version of the same principle.

## Concerns and open questions

**1. The diagnosis locates the cause one layer off.** The prompt says "every participant reads it in full, every round" and treats this as a protocol property. It is a **skill-loading instruction** property. `SKILL.md:12` says "Always read `parley-deck/COOPERATION.md` first" and `SKILL.md:24` says "Load the full cooperation protocol before acting." §9.1 (`COOPERATION.md:857`) says "Read `parley-deck/COOPERATION.md`." None of these say "read the parts you need for this phase" — they say read the whole thing. The protocol grew; the loading instruction never adapted. The fix is therefore not just in the protocol text but in the skill's loading instruction and the §9 checklist. If we scope the rendered view but the skill still says "always read the full protocol first," agents will read the full protocol anyway and the scoping saves nothing. **The skill instruction and §9.1 must change in lockstep with any rendering change, or the optimization is a no-op.**

**2. The 71% figure needs scrutiny.** The prompt says "protocol alone = 71% of the round-1 read." I measured the prior deliberation (`protocol-overlay-local-extension`) and got different numbers: protocol = 104,480 bytes, round-01 sum of all four agent files = 98,843 bytes, so protocol = 51.3% of (protocol + round-01 artifacts) (PRIMARY, `wc -c` on the files). The 71% may be measuring something different — perhaps the prompt's measurement includes the `00-prompt.md` and `00-scoping-brief.md` (37,291 bytes for that brief alone) in the non-protocol side, or counts only a single agent's read rather than the sum. Either way, the protocol is the single largest item, but "71%" and "51%" lead to different urgency calibrations. The diagnosis should state exactly what the denominator is.

**3. The re-read compounding may be overstated for this idea's track.** The prompt cites `protocol-overlay-local-extension` (track `deliberation`, 3 rounds). This idea is track `standard` (§4.0 caps cross-review at 2 rounds). The 2.24× round-3 cost is a `deliberation`-track phenomenon. On `standard`, the maximum is round-02 (round-1 + round-2 = ~2× round-1), and round-03 does not exist unless escalated. The compounding is real but its magnitude is track-dependent, and most ideas run `standard` or `fast` (the `fast` track skips cross-review entirely). The felt slowdown across "the last few versions" would be dominated by the protocol-growth term (which hits every track, every round), not the re-read term (which only hits `deliberation` and `standard`).

**4. Did skipping the full COOPERATION.md read cost me?** Yes, marginally — but less than expected. I opened the §4 phase rules, §9 checklist, §7, §15 provenance, §6, and the §2/§3 structural sections on demand. I did not read §11 (transport mechanics, 16KB), §12 (pipeline blocks, 7.5KB), §13 (retrospective, 5KB), or §14 (loop engineering, 2.2KB) at all — and none of them were needed to answer the five questions. That is ~31KB (30% of the document) I never loaded, and I did not miss it. This is direct evidence that a phase-scoped view would work: I just ran a manual version of it. The one thing I had to go back for was the exact §9.1 line to confirm the loading instruction — which proves the point that the instruction is the lever, not the protocol body.

## Risks

**A. Scoping can silently drop a rule that prevents a real failure.** The prompt names §15 verification integrity, §7 protocol change, §6 rule 3, and §14 human brake as rules bought with real failures. If a phase-scoped view omits §15 from a Phase 1 prompt and an agent makes an unprovenanced claim, that is exactly the failure §15 was written to stop. **Mitigation:** the phase→blocks mapping must be conservative — when in doubt, include. §15.2 (provenance, ~21 lines) is cheap and should be in every phase's scope. §15.3–15.7 (conflicting verdicts, role concentration, correlated agreement) are review-phase-specific and can be scoped out of Phase 1 safely.

**B. The skill and §9 must change in lockstep (see Concern 1).** If they do not, the generated view is built but never used because the loading instruction still says "read the full protocol." This is the highest-probability failure mode — the protocol gets a new render flag, the skill doesn't update, and the read cost is unchanged.

**C. The core-protocol machinery is not yet built.** `meta-protocol-change-global-core-protocol/FINAL.md` says the implementation slice ships `parley protocol render`, `check`, and the core store — but the status of that idea's `00-prompt.md` is `round-01` (its own status field says so), and the FINAL is dated 2026-08-07, two days before this idea. If the core store and render verb are not yet implemented, the scoping mechanism depends on unfinished work. The fallback is prompt-side scoping by the facilitator (the 00-prompt.md names which sections to read), which costs nothing in infrastructure but puts the burden on the facilitator and is not enforceable.

**D. The re-read of prior rounds is not clearly separable from peer accountability.** A digest of a peer's round-1 analysis loses the nuance of their argumentation — the specific trade-offs they weighed, the evidence they cited. In a refutation-default review culture (LE-1, `COOPERATION.md:254`), a reviewer who reads a digest of an implementation critique rather than the full critique may miss the refutation attempt the original author made. **What breaks:** the reviewer's ability to judge whether the original author's refutation was valid, not just what they concluded. A digest that preserves the author's conclusions but drops their evidence chain degrades review quality. The safer approach for the re-read problem is a **structured digest** (frontmatter + summary + key claims + cited evidence locators), not a freeform summary — so a reviewer can follow the locators to the full text when they need the chain.

**E. The bundled skill copy is already stale.** The skill's `references/COOPERATION.md` is 104,570 bytes vs the embedded default's 104,480 — they differ (md5 mismatch, PRIMARY). This means agents loading the fallback instead of the live file get a slightly different protocol. This is a drift problem (§9.0 addresses it) but it also means any scoping solution must scope the **same** content across all three copies, or the phase mapping desynchronizes.
