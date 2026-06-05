---
agent: agy
idea: deliberation-driver
round: 1
date: 2026-06-05
---

## Summary

The proposed `deliberation-driver` seeks to automate the multi-agent deliberation loop past `round-01` using a linear dependency-gated graph. From the lens of protocol correctness and compliance with [COOPERATION.md](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/parley-deck/COOPERATION.md), I support the core architecture of a re-entrant, disk-derived cursor to compute ready phases. However, strict guards must be put in place to ensure that automation does not bypass mandatory deliberation, fabricate consensus, or auto-finalize empty/placeholder specifications.

## Proposed approach

To ensure absolute adherence to the rules of [COOPERATION.md](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/parley-deck/COOPERATION.md):
1. **Disk as Single Source of Truth**: The `Rebuild` method of the `Cursor` must derive state purely from filesystem artifacts. Crucially, the current transport must be parsed at the **individual idea level** (by checking the `transport:` field in `00-prompt.md`) rather than just the global workspace default in [COOPERATION.md](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/parley-deck/COOPERATION.md). This accommodates ideas like this one which explicitly select `local-dir` to enable auto-advance inside an otherwise `github-pr` governed repository.
2. **Default Deliberation (Cross-Review)**: Since Round 1 files are written independently without reading others' contributions, a cross-review step is conceptually mandatory for deliberation. The driver should default to running at least one cross-review round (`round-02`) before consensus drafting. A straight-to-consensus flow from `round-01` should be illegal unless explicitly bypassed.
3. **Escalation and Circuit Breaking**: The `MaxRounds` limit (default 4) must act as a strict circuit breaker. If a consensus draft is blocked and reopened repeatedly, exceeding `MaxRounds` must write a `blocking: yes` escalation message to `inbox/` and halt execution. The driver must not progress the idea until a user direction is recorded and subsequently quoted in the next round's artifacts.

## Concerns / open questions

1. **Auto-Finalization Content Gap**: When a consensus triage status reaches `TriageReady`, the driver is proposed to promote it to `final` autonomously. However, [consensus.Finalize](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/consensus/consensus.go#L196) only generates an empty/scaffold `FINAL.md` template. If the driver auto-finalizes and immediately starts Phase 5 (Implementation), the implementer agent will build against a blank specification. We must define a guard: the driver must either invoke the initiator agent to compile the consensus into a drafted `FINAL.md`, or halt/escalate for human intervention to populate the specification before proceeding.
2. **Authenticity of Signoffs**: Under `--auto`, the driver must invoke actual participant agents to author their own signoffs via [consensus.AppendSignoff](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/consensus/consensus.go#L140). The driver must never mock or fabricate ACCEPT signoffs, and it must verify that all expected participant artifacts exist before executing finalization.
3. **Transport Swapping mid-idea**: If the transport configuration in `00-prompt.md` is changed during an active run, the driver must detect this immediately. The re-entrant `Advance` loop must evaluate the transport gate on every tick.

## Risks

1. **Building Empty Specifications**: The highest risk is that the driver automates the transition from `consensus` to `final` and then to `implementation` without checking if the generated `FINAL.md` actually contains the specification details, wasting agent tokens and producing corrupted or blank code changes.
2. **Deliberation Bypass**: Bypassing cross-review rounds and going straight to consensus prevents agents from formally debating and responding to other participants' independent round-1 ideas.
3. **Cursor Desynchronization**: If the cursor state diverges from the physical files or fails to handle malformed consensus files gracefully, it could cause the driver to enter an infinite loop or execute out-of-order phase actions.
