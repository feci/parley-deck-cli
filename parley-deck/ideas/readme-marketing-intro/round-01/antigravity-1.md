---
agent: antigravity-1
idea: readme-marketing-intro
round: 1
date: 2026-06-20
---

## Summary

As `antigravity-1`, our analysis focuses on inspirations and attribution. We believe that presenting a marketing-oriented introduction to the README must also accurately credit the prior art and intellectual foundations of the Parley Deck protocol. By adopting an "adopted-and-adapted" framing, we can honestly attribute our core mechanics to their original sources, outline our specific twists, and explicitly avoid any false implication of endorsement.

## Proposed approach

We propose including an **Inspirations & Lineage** subsection within the new README introduction. This section uses the "adopted-and-adapted" framework to credit prior work while highlighting Parley Deck's unique multi-agent and protocol twists:

### Inspirations & Lineage (Adopted & Adapted)

Parley Deck builds upon and adapts key concepts from the following prior art and tool architectures to fit our transport-agnostic, multi-agent consensus workflow:

*   **OpenRouter Fusion (Adapted for Multi-Agent Consensus):**
    *   *Lineage:* We adopt the *compare-not-merge* philosophy and the focus on *synthesis-as-value* for consensus.
    *   *Our Twist:* Instead of a real-time LLM API ensemble, Parley Deck applies this to asynchronous, multi-round markdown discussions, creating structured "Comparison & blind spots" analysis sections directly in consensus logs.
*   **OpenAI ExecPlans / PLANS.md (Adapted for Multi-Agent Execution):**
    *   *Lineage:* We adopt the concept of self-contained state encapsulation where execution can be resumed entirely from the documentation.
    *   *Our Twist:* We formalize this into a strict division between static design (`FINAL.md`) and a living progress document (`IMPLEMENTATION.md`), governed by multi-agent review-consensus cycles.
*   **RHO / Retrospective Harness Optimization (Adapted for Quorum Gates):**
    *   *Lineage:* We adopt retrospective optimization to audit past executions and adjust local settings.
    *   *Our Twist:* Rather than single-model self-preference, RHO optimization cycles in Parley Deck are advisory-only and gated by multi-agent quorum consensus (`parley retro`).
*   **kindly (Adapted for Local CLI Pipelines):**
    *   *Lineage:* We adopt strict gates, stopping judgments, and non-suppression of review dispositions ("validated artifact beats nonzero exit code").
    *   *Our Twist:* We integrate these supervision features into a local markdown pipeline using standardized severity tiers (CRITICAL, MAJOR, MINOR, NIT) and append-only signoff blocks.
*   **Preflight Readiness (Adapted for Distributed Roster Liveness):**
    *   *Lineage:* We adopt checking system readiness before initiating actions.
    *   *Our Twist:* We extend this into §9.0 project protocol freshness and active roster liveness pings (`hosted-PONG`) to ensure agent alignment before any round begins (`parley preflight`).

*Disclaimer: Reference to these projects and organizations is for attribution and lineage purposes only; no endorsement, sponsorship, or affiliation is implied.*

## Concerns / open questions

*   **Marketing Pitch Balance:** How prominent should the lineage section be in the final intro? We need to balance marketing punchiness with faithful attribution. Placing this section directly under the key features list provides context without disrupting the main user flow.
*   **Terminology and Links:** Should we provide external links to the original projects (e.g., the `codex_exec_plans` cookbook or RHO papers)? Adding references could enhance credibility, but we must make sure they don't draw users away from the core README.

## Risks

*   **Implied Endorsement:** Using names like OpenAI, OpenRouter, and RHO might lead some users to believe Parley Deck is sponsored or officially recommended by these entities. We mitigate this risk using the explicit disclaimer and emphasizing the "adopted-and-adapted" twist framing.
*   **Lineage Dilution:** If the specific modifications (our twists) are not clearly highlighted, users might dismiss Parley Deck as a redundant re-implementation of existing workflows. We must keep the focus on how we operationalized these concepts into a durable, multi-agent protocol.
