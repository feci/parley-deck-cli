---
agent: hermes-1
idea: meta-protocol-change-devx-speed
round: 1
date: 2026-07-03
---

## Summary

The protocol is 1046 lines across 15 sections with a 9-phase lifecycle, and it has accreted four major amendments (§12 pipelines, §13 retro, §14 outer loop, and the LE-1 through LE-11 loop-engineering patches woven through §4/§8). It reads as a document authored *for* protocol authors — dense, undifferentiated, and with no graded path for low-risk work. The biggest wall-clock sinks are the 30-minute per-agent deep-reasoning timeouts (confirmed in `agents.toml`: `timeout_ms = 1800000` for every agent), the strictly sequential round structure where each round is a full sync barrier, and the mandatory full-quorum consensus gate at every phase transition. My proposal is a three-track tiering model with objective triggers, a restructured document that uses progressive disclosure (a ~150-line core + appendices), and a set of speed levers that collapse phases for low-risk work without touching the safety core.

## 1. Time-sink audit

The top 6 places the protocol spends wall-clock or effort out of proportion to the risk it buys:

**1. Deep-reasoning timeouts — 30 min × N agents × N rounds.**
`agents.toml` sets `timeout_ms = 1800000` (30 minutes) for every agent — claude (opus, reasoning=max), codex, hermes, agy. For a 4-agent roster doing 2 design rounds, that's potentially 4 hours of wall-clock just in agent execution time, before any human review. The timeout is uniform regardless of whether the task is "rename a variable" or "redesign the consensus mechanism." There is no per-track timeout. The protocol itself doesn't set these timeouts (§0 says they live in `agents.toml` `[defaults]`), but the protocol also never suggests tuning them down for low-risk work. *Cost-to-risk ratio: terrible for trivial work, acceptable for risky work.*

**2. Sequential round structure — every round is a full sync barrier.**
Phase 1 (§4) requires all participants to write round-01 files before anyone opens round 2 (§4 Phase 2: "Once all participants submitted round N, any agent may open round N+1"). Round 1 is parallel, but every subsequent round is a serial dependency: wait for all → read all → write response → wait for all again. With 3+ rounds (common for contested ideas), this is 3 sequential agent-execution cycles. A round that one agent finishes in 5 minutes still blocks on the slowest agent's 25-minute run. *Cited: §4 Phase 1, §4 Phase 2.*

**3. Mandatory full-quorum consensus at every phase transition (§5).**
Quorum = all agents in `participants:` (§5). Every active participant must ✅ sign off (§4 Phase 3: "✅ from every active participant = consensus reached"). One slow or absent agent blocks the entire idea. The async rules (§5: "Agent silent past deadline is treated as ✅ — but only if they were pinged via `inbox/` first") require an *additional* inbox ping round-trip before silence counts as assent — adding another sync point to work around the first one. For a 2-line typo fix, requiring all 4 agents to sign off is ceremony out of proportion to risk. *Cited: §4 Phase 3, §5.*

**4. §9.0 readiness check — liveness ping ceremony before every idea.**
Before *every* idea, the facilitator runs a readiness check (§9.0): protocol-freshness comparison, roster liveness ping of every participant, available/unavailable table, and user confirmation for any exclusion. This is valuable for a 4-agent deck where someone might be offline, but it's a fixed tax on every idea regardless of size. The protocol-freshness sync logic alone (§9.0: `protocolRole` branching, `source` vs `consumer`, SHA comparison, auto-sync with zone preservation) is ~30 lines of conditional rules that fire before any actual work begins. *Cited: §9.0.*

**5. The design→implement→review→fix-up loop (Phases 4→5→6→7→8).**
Each transition is a sync point with its own consensus: Phase 3 consensus → Phase 4 FINAL → Phase 5 IMPLEMENTATION → Phase 6 review (all non-implementers) → Phase 7 review consensus → Phase 8 fix-up → repeat 6→7→8 until zero agreed fixes. For a change that touches one file, this is 5+ sync barriers, each requiring agent execution cycles and signoff rounds. The review loop (§4 Phase 6-8) mandates that *every* active participant except the implementer writes a review file with refutation attempts (LE-1), and review consensus requires all signoffs again (§4 Phase 7). *Cited: §4 Phases 3-8, LE-1.*

**6. The 1046-line protocol document itself.**
A developer joining the project must read COOPERATION.md to participate. At 1046 lines it's longer than most RFCs. The structure front-loads transport selection (§0, ~40 lines), roster mechanics (§2, ~40 lines), and directory layout (§3, ~40 lines) before reaching the actual workflow (§4). §11 (transport mechanics, ~200 lines) is in the main body despite being reference material. §12 (pipeline, ~120 lines), §13 (retro, ~50 lines), and §14 (outer loop, ~50 lines) are additive amendments that most ideas never touch. The §10 TL;DR (§10) is 12 lines — too terse to orient a newcomer, and buried at line 684. *This is the core DevX problem: the document has no graded entry point.*

## 2. Tiering model

I propose three named tracks. The track is determined by objective triggers at Phase 0 (kickoff), recorded in `00-prompt.md` frontmatter as `track: fast | standard | deliberation`.

### Track A — `fast` (lightweight)

**Objective triggers (ANY one suffices):**
- Files touched ≤ 3 AND no security/auth/secret-handling surface
- Change is fully reversible (git revert with no data migration, no external side effects)
- Diff is mechanically verifiable (lint, type-check, unit test) AND the idea is `auto_implement`
- Documentation-only or protocol-text-only change with no behavioral semantics change

**Phases kept / collapsed / made optional:**
- Phase 0: kept (lightweight — no §9.0 liveness ping; freshness check only if `protocolRole: consumer`)
- Phase 1: kept (single round, parallel — this is the valuable part)
- Phase 2: **collapsed** — no cross-review rounds; go straight to consensus after round 1
- Phase 3 + Phase 4: **collapsed into one step** — drafter writes a combined `consensus+FINAL.md` (or a single `FINAL.md` with embedded signoffs); signoff from **one non-author non-facilitator participant** suffices (not full quorum)
- Phase 5: kept (implement to FINAL.md)
- Phase 6: **single-reviewer fast-path** — one non-implementer review with refutation attempts; no multi-reviewer rounds
- Phase 7: **collapsed** — single reviewer's ✅ = review consensus
- Phase 8: kept (fix-up if the single reviewer finds issues)
- Timeout: **5 min per agent** (override `agents.toml` default for this idea)

### Track B — `standard` (normal)

**Objective triggers (DEFAULT track if none of the fast or deliberation triggers fire):**
- Files touched 4-20 OR moderate complexity
- Change is reversible but may require coordination (multiple files, config changes)
- Not a production deployment, not a protocol change, not a security-sensitive change

**Phases kept / collapsed / made optional:**
- Phase 0: kept (full §9.0 readiness check)
- Phase 1: kept (parallel round 1)
- Phase 2: **capped at 2 cross-review rounds** — if no ❌ BLOCK after round 2, proceed to consensus (currently the protocol allows unbounded rounds: §4 Phase 2 "Continue until nobody has new substantive objections")
- Phase 3 + Phase 4: **kept separate** but with a **fast-consensus** option: if all participants ✅ within a configurable window (default 1 round), skip additional rounds
- Phase 5: kept
- Phase 6: **2 reviewers minimum** (not full quorum) — if >2 non-implementer participants exist, only 2 need to review; others may opt out via inbox note
- Phase 7: kept (signoff from reviewers who actually reviewed)
- Phase 8: kept, with `MaxFixupCycles` defaulting to 2 (currently unbounded for standard ideas)
- Timeout: **15 min per agent**

### Track C — `deliberation` (full protocol, unchanged)

**Objective triggers (ANY one forces this track):**
- `type: meta-protocol-change` (changes to COOPERATION.md itself — §7)
- `risk: production` or `risk: high` in pipeline manifest (§12.3)
- Security-sensitive: auth, secrets, permissions, data migration, external side effects
- Irreversible changes (data loss potential, production deployment without rollback)
- `strict_gate: true` in `00-prompt.md`
- Files touched > 20 OR the change modifies the non-solo execution requirement (§1)

**Phases:** All phases 0-8, full quorum, full refutation-default (LE-1), unbounded rounds (subject to MaxRounds), full §9.0 readiness, 30-min timeouts. **This is the current protocol, unchanged.**

### Track-selection mechanics

The track is set at Phase 0 by the idea author, recorded as `track: <name>` in `00-prompt.md`. The track is **binding** — a participant may challenge the track assignment (open a round-1 objection), and the challenge is resolved by the quorum. **Any participant may force escalation to `deliberation`** by posting an inbox note before round-1 closes; this is the safety valve against a facilitator under-tiering a risky change. The `strict_gate: true` flag (§4 Phase 8) automatically forces `deliberation` regardless of other triggers.

## 3. DevX improvements

### 3.1 Document restructure — progressive disclosure

The single highest-impact DevX change. Today the protocol is a flat 1046-line document. Restructure into:

**Core (target ~150 lines, the front of the document):**
- §0 Choose transport (compressed to the table + 2 sentences)
- §1 Scope (compressed: the 4 goals + non-solo requirement, ~15 lines)
- §4 Protocol phases (compressed: one paragraph per phase with the track variations, ~60 lines)
- §5 Quorum (compressed: 5 lines)
- §6 Conflict-avoidance (unchanged, 6 rules, ~15 lines)
- §10 TL;DR (expanded slightly and **moved to line ~20**, right after the header — it's the first thing a newcomer should read)
- New: **§0.5 Quickstart** (see below)

**Appendices (moved to the back, clearly marked as reference):**
- Appendix A: Adopting in a new project (already an appendix)
- Appendix B (was §11): Transport mechanics — local, GitHub PR, GitLab MR. ~200 lines of reference material that most agents read once.
- Appendix C (was §12): Pipeline blocks & action stages. ~120 lines, opt-in, only for pipeline users.
- Appendix D (was §13): Retrospective optimization. ~50 lines.
- Appendix E (was §14): Automated outer loop. ~50 lines.
- Appendix F (was §9): Session-start checklist. ~40 lines, reference.

**Roster and meta stay in the core** (§2, §3) because every participant needs them, but §3 (directory layout) can be compressed to the tree diagram + 5 lines.

This gives a developer a ~150-line core that covers the entire workflow, with "see Appendix B for GitHub PR mechanics" pointers. The 1046 lines still exist — they're just not in the path of a first-time reader.

### 3.2 Quickstart (new §0.5, ~30 lines, in the core)

```
## Quickstart — start an idea in 5 minutes

1. Create ideas/<slug>/00-prompt.md (copy the template in Appendix A).
   Set track: fast | standard | deliberation (see §4 triggers).
2. Each participant writes round-01/<agent-id>.md independently.
   Do NOT read others' files first.
3. For `fast` track: the drafter writes FINAL.md with embedded signoffs.
   One non-author participant reviews and ✅. Done.
4. For `standard` track: cross-review (max 2 rounds), consensus, FINAL.md.
   2 reviewers review the implementation. Done.
5. For `deliberation` track: follow the full phase lifecycle (§4).

Defaults if you're unsure: track=standard, transport=local-dir,
quorum=all participants, timeout=15min/agent.
```

### 3.3 Defaults that "just work"

- **Default track: `standard`** — if `track:` is absent in `00-prompt.md`, the idea runs standard. This is safe-by-default without being ceremonial.
- **Default transport: `local-dir`** — for new projects without a git host. §0 already lists this as "simplest setup" but doesn't call it the default.
- **Default timeout: tiered by track** — `agents.toml` currently hardcodes 30 min for everyone. The protocol should recommend track-specific timeouts (5/15/30 min) and the skill should seed them.
- **Default reviewers: 2 for standard, 1 for fast** — not full quorum. Full quorum is the `deliberation` default.

### 3.4 Role-based framing

Add a short "Who are you?" table near the top of the core:

```
| If you are...        | Read...                                    | You do...                    |
|----------------------|--------------------------------------------|------------------------------|
| Developer (first time| §0.5 Quickstart, §4 phases, §6 rules       | Write round files, sign off  |
| Facilitator          | §0.5, §4, §9 checklist, §11 (your transport)| Open ideas, drive consensus  |
| Reviewer             | §4 Phase 6-8, LE-1 refutation-default       | Review code, find issues     |
| Pipeline builder     | Core + Appendix C (pipelines)              | Design pipeline manifests    |
```

This is 8 lines and tells a developer exactly where to start instead of "read 1046 lines."

### 3.5 What to move out of the main doc

- **§11 (transport mechanics, ~200 lines)** → Appendix B. It's reference material read once at setup.
- **§12 (pipeline, ~120 lines)** → Appendix C. It's opt-in and most ideas never use it.
- **§13 (retro, ~50 lines)** → Appendix D. It's a periodic process, not per-idea.
- **§14 (outer loop, ~50 lines)** → Appendix E. It's automation policy, not per-idea.
- **§9 (session-start checklist, ~40 lines)** → Appendix F. It's a reference checklist.
- **The LE-1 through LE-11 patches** — these are woven into §4 and §8 as inline paragraphs. They should be **extracted into a consolidated "Loop engineering rules" subsection** (or Appendix G) with cross-references from the phases they affect. Today a reader hits "LE-1" in §4 Phase 6 with no context — it's insider jargon. Either inline the rule as plain English ("Reviewers must attempt to break each acceptance criterion...") or collect them in one place.

This moves ~460 lines to appendices, leaving ~586 in the core — and with compression, the core target is ~150 lines.

## 4. Speed levers

### 4.1 Parallel rounds (cross-review)

Today, cross-review rounds (Phase 2) are sequential: all of round N must complete before round N+1 opens. But within a round, agents already work in parallel (each writes their own file). The serialization is at the *transition* between rounds.

**Lever:** For `standard` track, allow **streaming starts** — an agent may open round N+1 as soon as they have read the round-N files that are *already present*, without waiting for the slowest agent. Their round N+1 file records `responding-to:` only the files they actually read. Late-arriving round-N files are addressed in a subsequent round. This breaks the "wait for all" barrier without losing independence (the agent still writes their own file before reading responses to it).

**Risk:** An agent might miss a round-N file that arrives after they start round N+1. Mitigation: the facilitator posts an inbox note when all round-N files are in, and any agent who started early must address the late arrival in their next round. For `deliberation` track, keep the strict "wait for all" rule.

### 4.2 Single-reviewer fast-path (Phase 6)

For `fast` track, one non-implementer review with refutation attempts (LE-1) suffices. For `standard` track, 2 reviewers suffice (not full quorum). Full quorum review remains only for `deliberation`.

This is the single biggest wall-clock win for implementation review: today, every non-implementer must write a review file (§4 Phase 6), and with 4 agents that's 3 review files × 15-30 min each = 45-90 min of review, all sequential because review consensus (Phase 7) needs all signoffs. Reducing to 1-2 reviewers cuts this to 15-30 min.

**Safety:** The non-solo requirement (§1) is preserved — one non-implementer reviewer is still independent verification. The refutation-default (LE-1) is preserved — the single reviewer still attempts to break each criterion. What's dropped is *redundant* multi-reviewer coverage, which is appropriate for low-risk reversible work.

### 4.3 Collapse consensus + finalize (Phases 3+4)

For `fast` track, merge Phase 3 (consensus) and Phase 4 (FINAL) into a single step: the drafter writes `FINAL.md` with embedded signoff blocks, and the one required non-author participant appends their ✅. No separate `consensus.md` file.

For `standard` track, keep them separate but allow **simultaneous drafting** — the drafter writes `FINAL.md` as soon as the last round closes, and participants sign off on `consensus.md` while reading `FINAL.md` (not sequentially: consensus → then finalize). If a participant ❌ BLOCKs, the FINAL is withdrawn and a new round opens.

### 4.4 Auto-advance / driver

The protocol already has a driver concept (§4 LE-2, LE-5, LE-7, §12). Extend it: for `fast` and `standard` tracks, the driver may **auto-advance** between phases when the track's signoff predicate is met — no human gate at every transition. Block-boundary gates (§12.8) remain `human` only for `deliberation` track and for `production` risk.

Specifically:
- `fast` track: driver auto-advances from round-1 → FINAL → implementation → review → complete, pausing only for the required signoffs (1 non-author for design, 1 non-implementer for review).
- `standard` track: driver auto-advances through phase transitions once the signoff predicate passes, but pauses for a human gate at FINAL → implementation (the "are we sure this is ready to build?" moment).
- `deliberation` track: every transition is a human gate (current behavior).

### 4.5 Timeout tuning

`agents.toml` sets `timeout_ms = 1800000` (30 min) for every agent. This is the single largest wall-clock component. Proposal:

- `fast` track: 5 min/agent (override at the idea level via `00-prompt.md` `timeout_ms: 300000`)
- `standard` track: 15 min/agent
- `deliberation` track: 30 min/agent (current)

The protocol should **recommend** tiered timeouts in §0 or §4, and the skill should seed them per-track in `agents.toml`. The timeout is an agent-config concern (§0), but the protocol currently gives no guidance on tuning it — a 30-min timeout for a typo fix is absurd.

### 4.6 When one agent is genuinely enough

**Never for the artifact.** The non-solo requirement (§1) is a hard invariant: at least one non-facilitator must write a canonical artifact. One agent is never enough to satisfy Parley Deck.

**One agent IS enough for:** the facilitator's own analysis (round-1 is independent by design), the implementation (Phase 5: the default implementer is the FINAL drafter, solo), the fix-up (Phase 8: the implementer applies fixes solo). The protocol already allows solo work in these phases. The speed win is in *review*: for `fast` track, one non-implementer reviewer is enough (today it's all non-implementers).

### 4.7 Round cap for standard track

§4 Phase 2 says "Continue until nobody has new substantive objections" — unbounded. For `standard` track, cap at **2 cross-review rounds**. If there's still a ❌ BLOCK after round 2, escalate to the user (§4 escalation) or force-upgrade to `deliberation`. This prevents a 2-participant disagreement from spinning for 5 rounds.

## 5. Modern agentic concepts mapped to concrete protocol edits

### 5.1 Conditional rigor / right-altitude → tiering model (§4)

**Concept:** Not every task needs the same rigor. Match the process altitude to the task's risk and reversibility. (See: Anthropic's "right altitude" framing, OpenAI's conditional agent autonomy levels, 2025-2026.)

**Protocol edit:** Add `track: fast | standard | deliberation` to `00-prompt.md` frontmatter (§4 Phase 0). Add a §4.1 "Track selection" subsection with the objective triggers from my §2 above. Modify each phase description to show the track-specific behavior (collapsed/kept/optional). This is the structural change that enables most other speed levers.

### 5.2 Lead-agent + subagent orchestration → explicit in §1

**Concept:** A lead agent orchestrates subagents for parallel subtasks. (See: Claude's subagent dispatch, OpenAI Codex's parallel agents, 2025-2026.)

**Protocol edit:** §1 "Internal helpers" already permits subagents ("An agent MAY use internal helper mechanisms such as subagents..."). Strengthen this: explicitly allow a participant to use subagents for **parallel refutation** in Phase 6 — e.g., one subagent per acceptance criterion, each attempting to break it independently, results consolidated by the participant into their review file. Add a sentence to §4 Phase 6: "A reviewer MAY dispatch internal subagents to parallelize refutation attempts across acceptance criteria; the participant remains accountable for the consolidated review file."

### 5.3 Deterministic workflows vs model-driven control flow → §12 + driver

**Concept:** Deterministic, declarative workflows for predictable paths; model-driven control flow only where judgment is needed. (See: the "bitter lesson" applied to agent orchestration — less scaffolding, more model judgment where it matters, deterministic rails where it doesn't.)

**Protocol edit:** §12 (pipelines) already embodies this (deterministic block sequencing, driver executes approved side effects). Extend the principle to the base protocol: the **track selection** (§5.1 above) is a *deterministic* routing decision based on objective triggers — no model judgment needed. The *phases within a track* are deterministic. Model judgment is reserved for the round content (analysis, review, refutation) — which is where it adds value. Add a §4 preamble: "The protocol separates deterministic routing (track selection, phase sequencing) from model-driven work (analysis, review, refutation). Track selection is objective and mechanical; round content is agent-driven."

### 5.4 Plan mode → Phase 0 + track selection as a "plan"

**Concept:** A planning step that gets approved before execution begins. (See: Claude's plan mode, Cursor's plan mode, 2025-2026.)

**Protocol edit:** Phase 0 (kickoff) + track selection already IS a plan — `00-prompt.md` defines the problem, participants, and now the track. Make the **track assignment** an approvable artifact: any participant may object to the track in round-1 (forcing upgrade to `deliberation`), which is the "approve the plan" moment. No new phase needed — just make the track challengeable in round 1, which I already propose in §2.

### 5.5 Spec-driven development → FINAL.md is already the spec

**Concept:** Write the spec first, implement to the spec, verify against the spec. (See: GitHub Spec Kit, OpenAI's spec-driven dev, 2025-2026.)

**Protocol edit:** FINAL.md (§4 Phase 4) already IS the spec, and §4 Phase 5 says "Implements strictly to FINAL.md." The protocol is already spec-driven. The gap is that FINAL.md's "Observable acceptance criteria" (§4 Phase 4) are not always leveraged in review. Strengthen the link: in §4 Phase 6, make the refutation-default (LE-1) explicitly map each acceptance criterion to a refutation attempt — "for each observable acceptance criterion in FINAL.md, the reviewer MUST record at least one failing-case attempt." This is already the spirit of LE-1; make it mechanically explicit.

### 5.6 Context engineering / progressive disclosure → document restructure (§3.1)

**Concept:** Give the model (or reader) only the context they need at this altitude, with progressive disclosure for deeper layers. (See: Anthropic's context engineering, 2025-2026.)

**Protocol edit:** The §3.1 document restructure IS this concept applied to the protocol itself. Core ~150 lines = the context every participant needs. Appendices = progressive disclosure for transport mechanics, pipelines, retro, outer loop. A developer reads the core; a pipeline builder reads core + Appendix C; a facilitator reads core + Appendix F + Appendix B. This is context engineering applied to the meta-document.

### 5.7 Parallel worktrees → already in §11.B, extend to all tracks

**Concept:** Parallel isolated worktrees for independent work streams. (See: git worktree-based agent isolation, 2025-2026.)

**Protocol edit:** §11.B already has the sub-branch protocol for round-1 isolation ("each agent works on `idea/<slug>/round-01-<agent-id>`"). Extend this to `standard` track cross-review rounds: agents may work on sub-branches and merge when ready, rather than all pushing to `idea/<slug>` and creating implicit serialization. This is a transport-mechanics change (Appendix B) and is opt-in.

### 5.8 Verification/refutation gates → LE-1, already present, strengthen

**Concept:** A gate that requires the agent to *try to break* the implementation, not just read it. (See: refutation-based verification, 2025-2026.)

**Protocol edit:** LE-1 (§4 Phase 6) already mandates refutation attempts. The concept is adopted. The edit is to make it **track-aware**: `fast` track requires refutation from 1 reviewer; `standard` from 2; `deliberation` from all. The gate itself (refutation-default) is non-negotiable across all tracks — it's in the "MUST stay" list (§6 below).

### 5.9 Right-sized autonomy → §14 human brake + track-aware auto-advance

**Concept:** Give automation as much autonomy as the task's risk allows, no more. (See: graduated autonomy levels, 2025-2026.)

**Protocol edit:** §14 (outer loop) already implements this for automated loops (discovery only, human gate for promotion). Extend to the driver (§4.4 above): `fast` track gets more driver autonomy (auto-advance through phases), `deliberation` gets less (human gate at every transition). This is the same principle applied at a different layer.

### 5.10 Skills → the protocol is already skill-based

**Concept:** Reusable, composable procedural knowledge packs. (See: Claude skills, agent skill ecosystems, 2025-2026.)

**Protocol edit:** The protocol is already distributed as a skill (`parley-deck-skill`). The edit is to make the **track selection** a skill-level concern: the skill's `parley init` / `parley run` commands should ask for or infer the track, seed the appropriate timeout in `agents.toml`, and template `00-prompt.md` with the `track:` field. This is tooling, not protocol text — but the protocol should *recommend* it in §0 or Appendix A.

### 5.11 "Closing the loop" → LE-7/LE-11 goal-done check, already present

**Concept:** A final verification that the goal was actually achieved, not just that the process completed. (See: loop-closing verification, 2025-2026.)

**Protocol edit:** LE-7/LE-11 (§4 Phase 8) already implements a goal-done check ("a fresh non-implementer agent verifies the FINAL.md observable acceptance criteria"). The concept is adopted. The edit is to make it **track-aware**: `fast` track skips the goal-done check (the single review is sufficient); `standard` track runs it only if a reviewer flagged ACCEPT-WITH-RESERVATIONS; `deliberation` track always runs it (current behavior).

### 5.12 The bitter lesson (less scaffolding) → collapse ceremony for low-risk work

**Concept:** Over-engineered process scaffolding often adds less value than simply letting a strong model do the work with minimal ceremony. (See: Rich Sutton's bitter lesson applied to agent workflows, 2025-2026.)

**Protocol edit:** The `fast` track IS this principle: for low-risk reversible work, collapse the 9-phase lifecycle to "write → one review → done" and trust the model + one independent check. The scaffolding (multi-round, full quorum, fix-up cycles) is reserved for `deliberation` where it earns its cost. This is not "delete the protocol" — it's "right-size the protocol to the task."

## 6. What MUST stay

These are the non-negotiables I would defend against a speed-at-all-costs push:

1. **Non-solo execution (§1).** At least one non-facilitator participant MUST write a canonical artifact. No track may collapse to a solo facilitator process. This is the defining property of Parley Deck — without it, the protocol is just a checklist.

2. **Durable audit trail — files are canonical (§0, §10).** Every decision, argument, and signoff lives in a file that survives context compaction. PR/MR conversations are ergonomic, not authoritative. No track may move canonical state to a non-durable surface.

3. **Refutation-default review (LE-1, §4 Phase 6).** Reviewers assume the implementation is wrong until they fail to break it. A "no findings" review is credible only with refutation attempts recorded. This applies to `fast` track too — the single reviewer must still attempt refutation. The number of reviewers may shrink; the refutation discipline may not.

4. **Human brake on automation (§14).** Automated loops may discover and draft only; promotion, implementation, merge, and consensus override require a human or full-quorum gate. No speed lever may bypass this.

5. **Signoff mechanism — consensus gates (§4 Phase 3, Phase 7).** Every phase transition that produces a canonical artifact (FINAL, IMPLEMENTATION complete) requires explicit signoff. The *number* of signoffs may be tiered (1 for fast, 2 for standard, all for deliberation), but the *mechanism* — append-only signoff blocks with ✅/🟡/❌ — is invariant.

6. **English-only (§6.6).** Cross-agent interoperability and reviewability depend on a single working language. No track may override this.

7. **No-secret rules.** The protocol never reads, prints, or commits secrets. No speed optimization may relax this.

8. **Irreversibility escalation.** Any change touching production, security, auth, secrets, or data migration MUST land in `deliberation` track. The objective triggers (§2) enforce this — no participant may silently down-tier a risky change.

9. **Track escalation safety valve.** Any participant may force-upgrade to `deliberation` by posting an inbox note before round-1 closes. This is the defense against a facilitator under-tiering.

10. **Independence of round 1 (§4 Phase 1).** Round 1 is written without reading other agents' files. No speed lever (including streaming starts in §4.1) may violate this for any track.

## 7. Prioritized shortlist

Ordered by (impact on G1+G2) ÷ (cost + risk). Each item is one line.

### MUST (highest impact-to-cost ratio, do first)

1. **MUST — Add `track: fast|standard|deliberation` to `00-prompt.md` with objective triggers (§2).** This is the structural enabler for every other speed lever. Without tiering, all changes pay the full ceremony cost. (High impact, low cost, low risk — additive frontmatter field.)

2. **MUST — Restructure the document: ~150-line core + appendices for §11/§12/§13/§14/§9 (§3.1).** This is the single biggest DevX win. A developer reads 150 lines instead of 1046. (High impact, medium cost — reorganization, not rewriting — low risk — content preserved.)

3. **MUST — Move §10 TL;DR to line ~20 and expand it into a quickstart (§3.2, §3.4).** The first thing a newcomer reads should orient them in 30 lines, not be buried at line 684. (High impact, very low cost, zero risk.)

4. **MUST — Single-reviewer fast-path for `fast` track; 2-reviewer for `standard` (§4.2).** Cuts review wall-clock from N×30min to 1-2×timeout. Non-solo preserved (1 reviewer is still independent). (High impact, low cost, medium risk — mitigated by refutation-default being mandatory.)

5. **MUST — Collapse Phases 3+4 for `fast` track; simultaneous draft for `standard` (§4.3).** Eliminates one sync barrier for low-risk work. (Medium-high impact, low cost, low risk — deliberation track unchanged.)

### SHOULD (high impact, slightly higher cost or risk)

6. **SHOULD — Tiered timeouts: 5/15/30 min per track, recommended in protocol and seeded by skill (§4.5).** Directly attacks the biggest wall-clock sink. (High impact, low cost, low risk — agents.toml already supports per-agent overrides; just add protocol guidance.)

7. **SHOULD — Cap cross-review at 2 rounds for `standard` track; escalate or upgrade if still blocked (§4.7).** Prevents unbounded rounds. (Medium impact, very low cost, low risk.)

8. **SHOULD — Auto-advance for `fast` and `standard` tracks; human gates only at `deliberation` transitions and FINAL→impl for standard (§4.4).** Reduces human coordination overhead. (Medium-high impact, medium cost — driver changes — medium risk — mitigated by signoff predicates still required.)

9. **SHOULD — Extract LE-1 through LE-11 from inline §4/§8 patches into a consolidated subsection or Appendix G (§3.5).** Reduces cognitive load — "LE-1" is insider jargon that a developer cannot parse. (Medium impact, low cost, low risk — text reorganization.)

10. **SHOULD — Streaming starts for `standard` track cross-review: agents may open round N+1 after reading available round-N files (§4.1).** Breaks the "wait for all" barrier. (Medium impact, low cost, medium risk — mitigated by facilitator inbox ping and late-arrival handling.)

### COULD (worthwhile but lower priority or higher risk)

11. **COULD — Parallel refutation via subagents in Phase 6 (§5.2).** Speeds review for ideas with many acceptance criteria. (Medium impact, medium cost — protocol edit + agent capability — low risk.)

12. **COULD — Extend sub-branch isolation (§11.B) from round-1 to all `standard` track rounds (§5.7).** Reduces merge serialization. (Low-medium impact, medium cost — transport mechanics change — low risk.)

13. **COULD — Track-aware goal-done check: skip for `fast`, conditional for `standard`, always for `deliberation` (§5.11).** Reduces close-decision ceremony for low-risk work. (Low-medium impact, low cost, medium risk — LE-7/LE-11 is a safety net; skipping it for fast track is acceptable because the single reviewer + refutation already covers it.)

14. **COULD — Default `track: standard` when absent, default `transport: local-dir` for new projects (§3.3).** Safe-by-default without ceremony. (Low impact, very low cost, zero risk — but mostly a skill/tooling concern, not a protocol-text change.)

15. **COULD — Role-based "Who are you?" table in the core (§3.4).** 8 lines, orients newcomers. (Low-medium DevX impact, very low cost, zero risk.)

## Concerns / open questions

1. **Track-gaming risk.** A facilitator motivated by speed could systematically under-tier changes. The escalation safety valve (§6 item 9) mitigates this, but only if participants are willing to challenge. Is one inbox note a sufficient barrier, or should track assignment require explicit quorum confirmation at Phase 0?

2. **Track-drift during implementation.** An idea starts as `standard` but implementation reveals it's riskier than expected (touches auth, say). The protocol needs an explicit **mid-idea upgrade** path — any participant may force-upgrade the track via an inbox note, and the idea re-runs from the current phase under the new track's rules. This is not in my proposal yet and needs to be specified.

3. **Appendix boundary.** Moving §12 (pipelines) to an appendix is safe — it's opt-in. But moving §9 (session-start checklist) to an appendix could cause agents to skip it. The checklist is a per-session obligation, not a per-idea one — is an appendix the right home, or should it stay in the core as a short "every session, do this" block?

4. **Timeout enforcement.** The protocol can *recommend* tiered timeouts, but `agents.toml` is the actual enforcement point. Should the protocol mandate that the skill seeds track-specific timeouts, or just recommend? Mandating crosses into tooling-spec territory (§14: "specified separately").

5. **Two-participant ideas under tiering.** §5 notes that two-participant ideas use the same rules. Under `fast` track, "one non-author non-facilitator reviewer" = the only other participant. Under `standard`, "2 reviewers" is impossible with 2 participants (one is the implementer, one is left). The tiering model needs an explicit floor: `standard` track with 2 participants → 1 reviewer (same as fast). The trigger should account for roster size, not just risk.

## Risks

1. **Restructuring risk.** Moving 460 lines to appendices changes the document structure that existing agents and skills depend on. Section numbers shift, cross-references break. Mitigation: do the restructure as a single atomic edit with updated cross-references, and run a `parley retro` pass afterward to catch broken references. The `protocolSha256` in `meta/version.json` will change — consumers of this protocol (if any) need a §9.0 freshness sync.

2. **Tiering complexity risk.** Adding tracks adds a routing decision at Phase 0. If the triggers are not objective and mechanical, the track selection becomes a source of disagreement — adding ceremony instead of removing it. Mitigation: the triggers in §2 are designed to be checkable by a script (file count, risk class, strict_gate flag). The skill should automate track suggestion.

3. **Single-reviewer degradation.** Reducing to 1 reviewer for `fast` track means a single agent's blind spot becomes the implementation's blind spot. The refutation-default (LE-1) mitigates this but doesn't eliminate it — a single reviewer may not think to refute a criterion they didn't notice. Mitigation: the `fast` track is only for low-risk reversible work; irreversible work is forced to `deliberation` with full quorum.

4. **Speed-vs-safety culture drift.** Once a `fast` track exists, there's social pressure to use it for everything ("it's faster, just use fast"). Mitigation: the objective triggers are binding — a change that touches auth cannot be `fast` regardless of facilitator preference. The escalation valve lets any participant enforce this.

5. **Amendment accretion continues.** This proposal is itself a meta-protocol-change that will add length. If the restructure moves 460 lines to appendices but the new tiering rules add 100 lines to the core, the net core reduction is only ~60 lines (from ~150 target + 100 tiering = 250). The tiering rules must be *concise* — a trigger table and a per-track phase summary, not a second 1046-line document. If the tiering section can't fit in ~40 lines, it's too complex.
