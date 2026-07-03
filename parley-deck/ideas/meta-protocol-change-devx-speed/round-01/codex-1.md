---
agent: codex-1
idea: meta-protocol-change-devx-speed
round: 1
date: 2026-07-03
---

## 1. Time-sink audit

### 1. Full-quorum all-the-time is the largest wall-clock multiplier

The current protocol makes quorum equal to every agent listed in `00-prompt.md` (§5) and requires every active participant to sign off for design consensus (§4 Phase 3) and review consensus (§4 Phase 7). That is defensible for protocol changes, security-sensitive work, production actions, and irreversible decisions, but it is too expensive as the default for ordinary developer tasks.

The cost is not just "more model calls." It creates a slowest-agent critical path at multiple gates:

- Phase 1 waits for every participant's independent analysis (§4 Phase 1).
- Phase 3 waits for every participant's signoff (§4 Phase 3).
- Phase 6 waits for every non-implementer reviewer (§4 Phase 6).
- Phase 7 waits again for every participant, including the implementer (§4 Phase 7).

The async rule "silent past deadline is treated as accepted only after an inbox ping" (§4 Phase 3, §5) is safe but adds another round trip. The protocol should keep full quorum for high-risk tracks and make it conditional elsewhere.

### 2. The Phase 1 to Phase 4 design lifecycle is too heavy for low-risk implementation work

The conceptual flow in §4 requires kickoff, independent analysis, optional cross-review rounds, consensus, and finalization before implementation. For an architectural decision, that is valuable. For a small bug fix, dependency bump, local refactor, documentation change, test addition, or mechanical CLI improvement, it creates a mini-ADR process before any code is touched.

The heaviest parts are:

- Separate `consensus.md` and `FINAL.md` artifacts (§4 Phase 3 and Phase 4), often repeating decisions.
- Mandatory frozen `FINAL.md` even when the useful artifact is the implementation plan plus acceptance criteria.
- Cross-review rounds (§4 Phase 2) that remain open-ended until "nobody has new substantive objections."

The protocol should distinguish "design is the work" from "implementation is the work." Many developer tasks need a short plan and explicit acceptance criteria, not a full design-finalization ceremony.

### 3. Review and fix-up cycles are rigorous but over-applied

The Phase 6 to Phase 8 loop is strong: review files include refutation attempts (§4 Phase 6), review consensus records agreed fixes (§4 Phase 7), and fix-up repeats until zero agreed fixes (§4 Phase 8). The optional strict gate is even stronger. That is appropriate for risky changes, but it is too slow as a baseline.

Specific cost drivers:

- Every non-implementer reviews, even if the change is low-risk and narrow (§4 Phase 6).
- Review consensus is a separate artifact and signoff gate after review files already contain findings (§4 Phase 7).
- Fix-up requires consensus, implementation update, re-review, and another consensus cycle (§4 Phase 8).
- `strict_gate: true` requires a fresh full-scope review round with zero findings of any severity (§4 Phase 8), which should remain opt-in and rare.

The protocol needs a single-reviewer fast path and a "review disposition in implementation PR" shortcut for low-risk tasks, while preserving refutation-default review for normal and high-risk tracks.

### 4. §9.0 readiness and session-start checks front-load cost before developers see value

The §9.0 pre-idea readiness check does useful things: protocol freshness, roster liveness, user-confirmed exclusions, and quorum lock. But for a developer trying to start a small task, it can feel like setup friction before the work exists. The full §9 session checklist also asks every agent to read inbox, open prompts, and PR/MR actions before proceeding.

The disproportionate costs are:

- Per-idea liveness probes of every rostered participant (§9.0), even for tasks that should not require full roster participation.
- User confirmation for excluding and re-including unavailable agents (§9.0), even when the selected track only needs one reviewer.
- Protocol freshness sync behavior in the same path as ordinary idea startup (§9.0), which mixes maintenance with task initiation.

The readiness check should become track-aware. A fast track should probe only the selected participants and background-check the rest where possible. Full freshness and roster reconciliation should remain mandatory for full deliberation, protocol changes, and automation.

### 5. GitHub/GitLab transport mirrors duplicate canonical-file work

For transport B/C, agents write canonical files and also mirror actions in PR/MR reviews, labels, assignees, descriptions, branch state, and native approvals (§11.B, §11.C). The canonical-file rule is correct, but the transport layer makes many phase transitions operationally noisy.

Examples:

- Phase 3 signoff requires both a `consensus.md` commit and a native PR review (§11.B Phase 3).
- Phase 6 review requires a review markdown file plus a PR review on the implementation PR (§11.B Phase 6).
- Phase labels must be bumped at each round, review round, consensus, fix-up, and completion (§11.B).
- Phase 5 commits `IMPLEMENTATION.md` directly to the parley-deck integration branch, while code changes live in another PR (§11.B Phase 5), forcing cross-repo bookkeeping.

The PR/MR mirror should be generated by tooling or reduced to "link canonical artifact + native review verdict" instead of hand-maintained ceremony.

### 6. The main protocol mixes user guide, normative rules, automation spec, and edge cases

`COOPERATION.md` is currently 1046 lines across the core lifecycle, transport mechanics, pipelines (§12), retrospective optimization (§13), and automated outer loop (§14). This is too much for a developer to read before contributing. The problem is not only length; the entry point is buried.

The current §10 TL;DR helps but appears after the full lifecycle and before detailed transport mechanics. The developer must still infer:

- Which path applies to their task.
- Which sections are mandatory for a quick fix.
- Which artifacts they personally owe right now.
- Which automation sections are opt-in and irrelevant to normal work.

The protocol needs a "start here" decision tree and a smaller normative core, with advanced automation and pipeline material moved into appendices or separate reference documents.

## 2. Tiering model with objective triggers

I propose replacing the single lifecycle default with four named tracks. The protocol should require a `track:` field in `00-prompt.md`, defaulted by a deterministic classifier and overridable only with a recorded rationale.

### Track A: Solo Scratch

Purpose: tasks where Parley Deck should not be invoked as a full multi-agent workflow because the safety benefit is negligible.

Objective triggers:

- Purely local, reversible work.
- No code shipped to users, no production data, no credentials, no security/auth/payment surface.
- At most one file or one generated artifact, or docs/comments/tests only.
- No protocol semantics, no public API contract, no dependency upgrade with runtime behavior impact.
- Expected diff under roughly 50 changed lines, or the task is an exploratory command/result request.

Phases:

- No Parley idea required.
- If the user explicitly invoked Parley anyway, record a short `inbox/` note or lightweight `ideas/<slug>/00-prompt.md` with `track: solo-scratch` and `solo_reason: low-risk-user-authorized`.
- Keep normal code review outside Parley if the repository requires it.

Why this is acceptable: §1 says Parley Deck is non-solo when a request uses Parley. I would preserve that rule for real Parley workflows, but add an explicit off-ramp: not every tiny developer action should become a Parley workflow. This is not a "solo Parley" claim; it is "Parley not needed."

### Track B: Fast Dev

Purpose: normal developer tasks where independent verification is useful but full deliberation is excessive.

Objective triggers:

- Reversible change.
- 1-5 files touched, no broad module boundary change.
- Expected diff under roughly 300 changed lines.
- No security/auth/secrets/payment/production-infra/data-migration surface.
- No protocol change, no public API break, no irreversible external side effect.
- Tests or deterministic checks exist or can be added.

Phases:

- Keep Phase 0, but use a compact `00-prompt.md` with `track: fast-dev`, objective triggers, acceptance criteria, selected driver, and one required reviewer.
- Collapse Phase 1-4 into a single `PLAN.md` or `FINAL.md` written by the driver, with one reviewer signoff before implementation only when the change is non-obvious. For obvious changes, the implementation PR plus acceptance criteria is enough.
- Phase 5 stays, but `IMPLEMENTATION.md` can be compact: summary, files changed, checks run, deviations.
- Phase 6 uses one non-implementer reviewer, selected by availability and relevant expertise.
- Phase 7 review consensus collapses into the reviewer's review file plus implementer disposition in `IMPLEMENTATION.md`.
- Phase 8 fix-up can complete after reviewer confirms agreed fixes, without a full new consensus file, unless new MAJOR/CRITICAL issues appear.

Required safety:

- At least one non-implementer reviewer for any code merge.
- Checks run and recorded.
- Escalate to Standard Deliberation if the reviewer finds architectural uncertainty, hidden coupling, or risk beyond the triggers.

### Track C: Standard Deliberation

Purpose: the current "normal" Parley Deck flow, but with bounded defaults.

Objective triggers:

- Cross-module behavior change.
- New feature or user-visible workflow.
- 6-15 files or roughly 300-1000 changed lines.
- Public API or persisted format additions that are backward compatible.
- Meaningful product/design choice.
- Multiple plausible implementation approaches.
- Any participant requests elevation with a concrete risk.

Phases:

- Keep Phase 0.
- Keep Phase 1 independent analysis, but run participants in parallel and use 2-3 participants by default.
- Make Phase 2 optional: open one cross-review round only if there are substantive disagreements, missing analysis, or a blocker.
- Collapse Phase 3 and Phase 4 when there are no blockers: `consensus.md` may include an embedded "Final artifact" section, or `FINAL.md` may include a signoff appendix. Pick one canonical final source to avoid duplicative prose.
- Keep Phase 5 implementation.
- Keep Phase 6 refutation-default review, but reviewers may be a representative subset of at least two non-implementers unless triggers require full quorum.
- Keep Phase 7/8, but allow a "fix-only verification" review when agreed fixes are narrow and no new risk is introduced.

Required safety:

- Non-solo independent analysis before major decisions.
- Refutation attempts against observable acceptance criteria.
- Durable artifacts and signoffs.

### Track D: Full / High-Risk

Purpose: preserve the current safety core for work where mistakes are expensive.

Objective triggers:

- Protocol changes under §7.
- Security, authentication, authorization, secrets, payments, privacy, compliance, or production infrastructure.
- Data migration, destructive operation, irreversible side effect, deployment/rollback, or provider action.
- Public API breaking change or persistent schema change.
- Broad changes over roughly 15 files or 1000 changed lines.
- `strict_gate: true`, `auto_implement`, pipeline/action block (§12), automated loop behavior (§14), or any production gate.
- User explicitly requests maximum rigor.

Phases:

- Keep Phase 0-8 as the current full lifecycle.
- Keep §9.0 full readiness.
- Keep full quorum unless a user-confirmed exclusion is recorded.
- Keep review consensus and fix-up loops.
- Use `strict_gate` when objective risk warrants it.

### Track escalation and de-escalation rules

Add a required "track check" at each phase boundary:

- Escalate if new files, risk surfaces, unavailable checks, review findings, or unknowns exceed the selected track's triggers.
- De-escalation is allowed only before implementation starts and only with recorded participant agreement.
- A reviewer can force escalation from Fast Dev to Standard by filing a MAJOR or CRITICAL concern that cites a trigger.
- A human can always choose a stricter track.

## 3. DevX improvements

### Add a real "Start here" quickstart at the top

Move a developer-facing quickstart before §0. It should answer "what do I do now?" in under one page:

1. Pick a track using the decision tree.
2. Run `parley start "<task>"` or create `ideas/<slug>/00-prompt.md`.
3. Write or wait for the artifact you owe.
4. Implement only after the track's plan gate is satisfied.
5. Review using the track's review rule.

The quickstart should include three examples:

- "Small bug fix" → Fast Dev.
- "Architecture decision" → Standard Deliberation.
- "Protocol/security/production change" → Full / High-Risk.

### Split the document into user-guide and reference layers

Keep `COOPERATION.md` as the normative core, but reorganize:

- Top: Quickstart, track decision tree, current transport/roster summary.
- Core: invariants, artifacts, phases by track, quorum/signoff rules.
- Reference appendices: detailed transport mechanics (§11), pipelines (§12), retro (§13), automated outer loop (§14), bootstrap instructions.

Alternatively, keep one file for canonicality but make the main path short and push detail behind headings named "Advanced reference." A developer should not need to read pipeline/action-stage rules to start a normal implementation.

### Make defaults explicit and tool-friendly

The default should be:

- `track: fast-dev` for small reversible code tasks.
- `track: standard-deliberation` for new features and ambiguous design.
- `track: full-high-risk` for protocol, security, production, automation, and irreversible changes.
- Participants: driver plus one reviewer for Fast Dev; 2-3 participants for Standard; full active roster for Full.
- Timeouts: 5-10 minutes for Fast Dev reviewer tasks; 15 minutes for Standard round files; 30 minutes only for Full/deep review.
- Transport: sticky current transport from the header (§0).
- Checks: required for code tasks, inferred from repo when possible.

The protocol already has model/reasoning bootstrap and central config (§0, §2). Reframe that as "setup once, then do not think about it."

### Add role-based instructions

Developers, PMs, designers, reviewers, and facilitators need different entry text. Add short role cards:

- Developer/implementer: "start here, choose track, write plan/implementation, run checks, record deviations."
- Reviewer: "look at acceptance criteria, try to break them, file findings with fixed severity."
- PM/designer: "state outcome, constraints, non-goals, risk tolerance; do not manage protocol mechanics."
- Facilitator: "only you handle roster/transport/phase labels unless automated."

This avoids making every participant internalize the whole system.

### Provide copy-paste minimal templates

Current templates are complete but verbose. Add minimal track-specific templates:

- Fast Dev `00-prompt.md`: task, track triggers, acceptance criteria, checks, participants.
- Fast Dev `IMPLEMENTATION.md`: summary, changed files, checks, deviations, reviewer disposition.
- Standard `FINAL.md`: only the sections needed for design decisions.
- Full remains as-is.

### First-time developer path in under five minutes

Concrete target flow:

1. Developer runs `parley start --fast "Fix CLI error when config is missing"`.
2. Tool creates `00-prompt.md` with `track: fast-dev`, selected driver, one reviewer, inferred checks, and `status: implementation-plan`.
3. Driver writes a 10-line plan/acceptance criteria.
4. Reviewer either approves the plan or says "elevate to Standard because <trigger>."
5. Developer implements and runs checks.

If tooling is unavailable, the same flow should be possible manually from the quickstart with one template and one command/check.

## 4. Speed levers

### Parallelize independent work by default

Phase 1 independent analyses are naturally parallel (§4 Phase 1). The protocol should state that facilitators launch all participants concurrently unless the transport or file ownership requires sequencing. The same applies to Phase 6 non-implementer reviews (§4 Phase 6). Sequential invocation should be the exception, not the normal path.

### Replace full quorum with track-specific quorum

Use:

- Fast Dev: driver + one independent reviewer.
- Standard: 2-3 participants, at least one non-facilitator and one non-implementer reviewer.
- Full: all active participants from `00-prompt.md`.

This preserves §1's non-solo property where it matters while avoiding full-roster latency for routine work.

### Collapse consensus and finalization when the content is duplicative

For Standard with no blockers and Fast Dev by default, collapse Phase 3 and Phase 4:

- Option A: `FINAL.md` includes decisions, trade-offs, deferred items, and signoffs.
- Option B: `consensus.md` becomes the final artifact and includes the plan/spec.

Do not maintain both unless the idea is Full / High-Risk or design-heavy. The current split (§4 Phase 3 and Phase 4) is useful for complex deliberation but duplicative for small tasks.

### Use an auto-advance driver for mechanical phase transitions

The protocol already has a human brake for automated loops (§14) and a driver concept in pipelines (§12). Extend that with a deterministic intra-idea driver:

- Detect required files exist.
- Validate frontmatter.
- Bump phase labels.
- Create empty next directories.
- Request reviews.
- Generate PR/MR links to canonical files.
- Refuse to advance if a gate is unmet.

The driver should not author participant content or decide contested issues. It should remove mechanical waiting.

### Tune timeouts by track

Current skill guidance often uses 20-30 minute agent timeouts. Keep that for deep and high-risk tasks. Add protocol defaults:

- Fast Dev: 5 minute plan/review target, 10 minute hard timeout per agent.
- Standard: 10-15 minute target, 20 minute timeout.
- Full: 30 minute timeout, extendable.

Timeouts should be budgets with escalation, not silent acceptance. If a Fast Dev reviewer times out, pick another available reviewer or escalate to Standard.

### Stream partial status, not partial artifacts

For long-running rounds, agents can stream status to PR comments or local orchestration logs, but canonical artifacts remain complete files. This improves human perception of speed without corrupting the audit trail.

### Make one-agent work legitimate outside Parley

The protocol should explicitly say when not to use Parley. If a task is Track A Solo Scratch, do the work normally and do not claim Parley Deck verification. This is cleaner than stretching §1 with solo exceptions for tiny work.

### Use fix-only verification after narrow fix-ups

After review consensus lists agreed fixes (§4 Phase 7), a narrow fix-up should not always require a full fresh Phase 6 review. For Fast Dev and Standard:

- If fixes touch only the lines/files cited in agreed findings, the reviewer may perform fix-only verification.
- If the fix touches broader code or introduces new behavior, run a fresh review.
- Full / High-Risk and `strict_gate: true` keep the current full-scope rule.

## 5. Modern agentic concepts mapped to concrete edits

### Conditional rigor / right-altitude planning

Concept: agent workflows should scale ceremony to risk and uncertainty, not run maximum process everywhere.

Concrete edit: add `track:` to `00-prompt.md`, define objective triggers, and make phase/quorum/review requirements track-specific. This is the main protocol change.

### Lead-agent plus subagent orchestration

Concept: a lead agent owns coordination and may use helper agents/tools internally, while accountability remains clear.

Concrete edit: keep §1 "Internal helpers" but add a "driver/lead responsibilities" section. Helpers may inspect, test, summarize, and propose, but only named participants write canonical artifacts. For Fast Dev, the lead implementer can use subagents for search/test generation without turning them into quorum members.

### Deterministic workflows over model-driven control flow

Concept: LLMs should produce judgment artifacts; deterministic code should enforce gates and state transitions.

Concrete edit: specify a machine-checkable state model for phases/tracks: required files, frontmatter fields, signoff statuses, reviewer counts, checks, and close predicates. The driver advances only when predicates are satisfied. This extends §12's ledger/driver discipline to ordinary ideas without requiring pipelines.

### Plan mode and spec-driven development

Concept: separate "decide what to build" from "build it," but keep the plan as small as the risk permits.

Concrete edit: for Fast Dev, replace full Phase 1-4 with a compact plan containing acceptance criteria, affected files, and checks. For Standard/Full, keep `FINAL.md` with observable acceptance criteria (§4 Phase 4).

### Context engineering / progressive disclosure

Concept: give agents the minimum relevant context first, with deeper references available on demand.

Concrete edit: reorganize `COOPERATION.md` into quickstart, track matrix, and reference appendices. Add participant prompts that cite only the relevant sections for the selected track. Move §12-§14 out of the default reading path for ordinary developer work.

### Parallel worktrees

Concept: isolate concurrent implementation/review attempts to avoid collisions and enable parallel exploration.

Concrete edit: for Standard/Full designs with competing implementation approaches, allow participants to prototype in separate worktrees before consensus. Record only conclusions and evidence in round files. Do not require worktrees for Fast Dev.

### Verification and refutation gates

Concept: reviews should actively try to falsify acceptance criteria, not just summarize.

Concrete edit: preserve §4 Phase 6 refutation attempts for Standard and Full. For Fast Dev, require a smaller "verification evidence" section: checks run, one negative/edge case if applicable, and reviewer confirmation. Keep strict refutation for high risk.

### Right-sized autonomy

Concept: automation can execute routine transitions but must stop at human/risk gates.

Concrete edit: add an auto-advance mode inside an active idea that can create directories, update labels, request reviews, and validate files. Keep §14's human brake for candidate promotion, finalization, production actions, roster changes, and consensus overrides.

### Skills and reusable task packs

Concept: agent skills package domain-specific instructions and templates.

Concrete edit: define "protocol skill packs" for common tracks: `fast-dev`, `standard-design`, `full-risk`, `reviewer`, `facilitator`. Each pack contains only the relevant template and rules. This makes participant prompts shorter and reduces mistakes.

### Closing the loop

Concept: workflows should feed outcomes back into future process tuning.

Concrete edit: keep §13 retrospective optimization but make it lightweight by default: each completed idea records `actual_track`, `escalations`, `rounds`, `review_cycles`, and `wall_clock_estimate`. Retros can then identify over-ceremony without reading every artifact.

### The bitter lesson: less bespoke scaffolding unless it buys measurable value

Concept: avoid process complexity that tries to encode every judgment path; use simple scalable primitives and strong verification.

Concrete edit: remove duplicate manual ceremonies where a simpler invariant suffices. The core primitives should be: track, owner, independent reviewer(s), acceptance criteria, checks, canonical artifact, signoff. Advanced pipelines and loop engineering stay available but should not dominate the common path.

## 6. What MUST stay

- Durable canonical files under `parley-deck/`; PR/MR conversations remain ergonomic mirrors (§0, §11).
- Non-solo independent verification for risky changes (§1), especially Standard and Full tracks.
- Round 1 independence for deliberation tracks (§4 Phase 1); no anchoring on other agents before the first analysis.
- One file per agent per round and no editing another agent's artifact (§6).
- Append-only signoffs for consensus and review consensus where those gates apply (§4 Phase 3, Phase 7).
- Human-authorized solo exceptions only when Parley is still claimed (§1, §5).
- English-only canonical artifacts unless a project deliberately overrides (§6.6).
- No-secret / sensitive-data caution from the skill and transport discipline; never optimize speed by leaking credentials or unrelated private material.
- Refutation-default review for meaningful code changes (§4 Phase 6), at least in Standard and Full tracks.
- Strict gate semantics for tasks that opt in or objectively require high assurance (§4 Phase 8).
- Human brake on automated outer loops (§14); automation may discover and scaffold but must not silently promote, merge, finalize, or override consensus.
- Immutable `FINAL.md` after closure for Standard/Full ideas (§4 Phase 4); later invalidation opens a new idea.
- Escalation-to-user path for human-only judgment and unresolved tie/block conditions (§4 Escalation, §5).

## 7. Prioritized shortlist

### MUST

- Add `track:` to `00-prompt.md` with objective triggers and four tracks: Solo Scratch, Fast Dev, Standard Deliberation, Full / High-Risk.
- Make quorum, required phases, reviewer count, and timeout defaults track-specific instead of full-roster by default.
- Add a one-page "Start here" quickstart and decision tree before §0 so developers can start without reading the full protocol.
- Preserve the existing full lifecycle unchanged for protocol changes, security/production/irreversible work, automation, and `strict_gate`.
- Collapse Phase 3 + Phase 4 for Fast Dev and low-conflict Standard ideas to remove duplicate `consensus.md`/`FINAL.md` prose.
- Add Fast Dev templates for compact `00-prompt.md`, implementation plan, review, and fix disposition.
- State explicitly when not to use Parley Deck: tiny Solo Scratch work should not claim Parley verification.

### SHOULD

- Parallelize Phase 1 and Phase 6 launches by default and make sequential invocation an exception.
- Add deterministic auto-advance tooling for mechanical phase transitions, frontmatter validation, labels, review requests, and missing-file checks.
- Tune default agent timeouts by track: short for Fast Dev, moderate for Standard, long only for Full/deep review.
- Permit one-reviewer Fast Dev review with mandatory escalation triggers when risk is discovered.
- Allow fix-only verification for narrow Fast Dev/Standard fix-ups while keeping full re-review for high-risk and strict-gate work.
- Reorganize §12-§14 into advanced reference material outside the default developer reading path.
- Add role cards for developer, reviewer, PM/designer, and facilitator.

### COULD

- Add lightweight metrics fields (`actual_track`, rounds, review cycles, wall-clock estimate, escalations) to support §13 retrospectives.
- Define reusable skill packs or prompt packs per track and role.
- Add optional parallel worktree guidance for competing Standard/Full implementation approaches.
- Generate PR/MR mirror actions automatically from canonical files to reduce §11 transport bookkeeping.
- Add a machine-readable track schema so tools can validate phase requirements without parsing prose.
