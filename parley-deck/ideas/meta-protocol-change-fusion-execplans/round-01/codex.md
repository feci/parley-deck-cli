---
agent: codex
idea: meta-protocol-change-fusion-execplans
round: 1
date: 2026-06-18
---

## Summary

My main hypothesis: ExecPlans are more useful to parley-deck than Fusion, but only if adapted into the existing artifact chain instead of replacing it. Parley-deck already owns the hard part of Fusion in a stronger, auditable form: named participants, independent round-01 files, cross-review, `consensus.md`, and append-only signoffs. The gaps worth addressing are narrower: better structured comparison in `consensus.md` / `review/consensus.md`, and more self-contained, resumable `FINAL.md` and `IMPLEMENTATION.md` artifacts for long-running work and the auto-drive driver.

Priority should go to low-ceremony additions that make future agents and `parley retro` better at reconstructing what happened. Do not import hidden judges, majority confidence gates, automatic web-enabled panels, or a single living plan document that absorbs the current multi-file protocol.

## Prioritized recommendation

### 1. ADAPT ExecPlan self-containment into `FINAL.md` and `IMPLEMENTATION.md`

Verdict: adapt. Benefit high. Cost medium, with the cost mostly in disciplined writing and keeping the implementation artifact current.

`FINAL.md` is already the single source of truth after Phase 4, but today its required shape is thin: final plan/specification plus references. That is sufficient for human review of a small idea, but weak for a fresh stateless implementer or the auto-drive driver resuming later. The ExecPlan concept that matters is not "one big plan file"; it is "the artifact contains enough orientation, observable success criteria, and recovery context to continue without reconstructing the round history."

I would adapt this as a content expectation for non-trivial `FINAL.md`, especially implementation or pipeline ideas:

- Purpose and user-visible outcome: what should be true after implementation.
- Context and orientation: relevant paths, existing modules, important terms, and constraints discovered during deliberation.
- Observable acceptance criteria: concrete commands, UI checks, API behavior, or artifact checks that Phase 6 reviewers and the driver can verify.
- Idempotence and recovery notes: what is safe to rerun, what must not be repeated blindly, and what to inspect after a partial failure.
- Known risks and de-risking steps: any prototype/spike or high-risk dependency the implementer should validate before broad edits.

`FINAL.md` should not become "living" after finalization. The current immutability rule is stronger: if the plan is invalidated, open a new idea or v2. The ExecPlan living-document discipline belongs primarily in `IMPLEMENTATION.md`, where Phase 5-8 already allow updates, deviations, fix-up cycles, and status changes.

For `IMPLEMENTATION.md`, I would adapt more of ExecPlans directly:

- Progress checkpoints with timestamps for multi-file, risky, or driver-managed work.
- Surprises and discoveries with evidence, especially when they alter implementation choices.
- Decision log entries for choices made after `FINAL.md`, with rationale and author/date.
- Validation and acceptance evidence, including commands run and what they proved.
- Outcomes and retrospective at completion, focused on what should feed section 13 / `parley retro`.

This would materially improve the auto-drive driver. Section 12 already gives the driver a durable effects ledger, idempotency keys, gates, and provider capability checks. What it still needs from human/agent-authored artifacts is task-level orientation: how to tell whether the implementation is done, what recovery means for this plan, and what deviations were intentional. A richer `IMPLEMENTATION.md` gives the driver a resumable narrative without letting it write participant artifacts or bypass gates.

### 2. ADAPT Fusion's compare-don't-merge rubric into consensus artifacts

Verdict: adapt. Benefit medium-high. Cost low-medium.

The Fusion judge's useful contribution is the structured comparison frame:

- consensus / agreement
- contradictions / disagreements
- partial coverage
- unique insights
- blind spots

Parley-deck already has independent panel work and visible raw answers. It does not need a hidden judge. But `consensus.md` today can become a drafted synthesis too quickly: "Agreed decisions", "Agreed trade-offs", and "Open items deferred to implementation" do not explicitly force the drafter to show what was compared or what everyone may have missed.

I would adapt this as a drafting checklist or optional section in `consensus.md`, not as a new authority. The value is especially in two fields:

- Contradictions: prevents disagreement from being smoothed into a vague trade-off.
- Blind spots: forces the drafter to ask "what did no participant cover?" before signoff.

`review/consensus.md` would benefit from the same frame. It already has "Agreed fixes", "Deferred follow-ups", and "Dismissed findings"; adding a comparison lens would help distinguish:

- findings everyone independently saw
- findings only one reviewer saw but the group accepts
- areas no reviewer inspected deeply
- dismissed findings that are factual non-issues versus accepted risks

Append-only signoffs remain the real gate. Any participant can block if the comparison is inaccurate. The file should remain canonical, and PR/MR comments should remain ergonomic mirrors.

### 3. ADAPT behavior-focused acceptance criteria as a first-class bridge from design to review

Verdict: adapt. Benefit high. Cost low-medium.

ExecPlans' behavior-focused validation maps cleanly onto parley-deck's existing phases. `FINAL.md` should give reviewers and implementers more than intent; it should say how success will be observed. This is not just test-command detail. It is a contract between Phase 4, Phase 5, Phase 6, and the auto-drive driver.

Concrete effect:

- `FINAL.md` names observable acceptance criteria.
- `IMPLEMENTATION.md` records which criteria were satisfied and with what evidence.
- Phase 6 reviewers classify findings against those criteria where possible.
- Phase 7 review consensus can separate true agreed fixes from deferred follow-ups or dismissed findings using the criteria as shared ground.

This should not alter the existing severities. `CRITICAL`, `MAJOR`, `MINOR`, and `NIT` are already a good implementation review vocabulary. Acceptance criteria would make severity assignment less subjective, not replace it.

### 4. ADAPT "confidently wrong" as a retro signal, not a review severity

Verdict: adapt lightly. Benefit medium. Cost low, but cultural risk if framed as blame.

Fusion's negative-weight benchmark idea points at a real protocol failure mode: an agent or consensus artifact can be confidently wrong, and ordinary review severities do not capture the process lesson. In parley-deck, this should not become a fifth severity. `CRITICAL` and `MAJOR` should stay code/artifact findings about what must be fixed.

The better home is section 13 evidence. `parley retro` needs structured inputs that distinguish:

- false factual claim accepted by the group
- unsupported assumption that shaped `FINAL.md`
- missed risk that later caused fix-up churn
- reviewer confidence that was not backed by inspected evidence

This can be captured in `IMPLEMENTATION.md` outcomes, `review/consensus.md`, or a retro-signal note attached to an idea. The point is to improve future harness proposals, not to score participants or weaken equal signoff.

### 5. REJECT confidence-by-breadth as a gate or weighting scheme

Verdict: reject for decision semantics; maybe record as non-binding context. Benefit low. Cost/risk medium.

Fusion treats agreement across model outputs as higher confidence. That makes sense for a one-shot router trying to decide how much to trust panel answers. Parley-deck is different: it has named accountable participants, explicit counter-proposals, and all-participant signoff. A 3-1 majority should not be allowed to override a real blocker, and a unanimous group can still be wrong.

Recording "all participants agreed on X" is useful context. Using breadth of agreement to change downstream gates is net-negative. It would erode the append-only signoff model and make the protocol less protective of minority but correct objections.

## Concept-by-concept verdicts

### Fusion panel execution

Verdict: already covered here in stronger form; reject importing it as protocol mechanics.

Parley-deck's Phase 1 already gives independent analyses from named participants before anchoring. It is stronger than Fusion's opaque panel because every participant's file is durable, attributable, and reviewable. The roster, quorum, and one-file-per-agent mechanics are better suited to protocol work than a dynamic panel of up to eight hidden model calls.

Optional internal helper calls are already allowed by the protocol, and they should remain non-canonical. If an agent wants to use a Fusion-like helper internally, that is fine, but the named participant remains accountable for its own `round-01/<agent-id>.md`.

### Fusion judge compare-don't-merge

Verdict: adapt as structure, reject as authority.

The drafter of `consensus.md` is already a synthesizer. The missing piece is an explicit comparison pass before decisions. A lightweight "comparison notes" discipline would improve transparency without creating a new judge role. A separate judge with special authority would conflict with equal participant signoff.

### Fusion final answer from judge analysis only

Verdict: reject.

Parley-deck should not hide raw round files from the final drafter or make the final artifact depend only on a judge summary. `FINAL.md` should continue to reference `consensus.md` and the rounds. The audit trail is a feature, not noise.

### Fusion unique insights and partial coverage

Verdict: adapt.

These are good fields for `consensus.md` and `review/consensus.md`. They help preserve minority insights that are not blockers and identify areas where only one participant did real coverage. This pairs well with append-only signoffs because participants can object if their contribution is misrepresented.

### Fusion blind spots

Verdict: adopt the lens, adapt the mechanics.

"What did nobody address?" is the highest-value Fusion concept for parley-deck consensus. It is not naturally produced by the current templates. The section 13 retro can discover blind spots after failure, but consensus and review need a forward-looking blind-spot check before finalization or completion.

### Fusion recursion protection

Verdict: mostly redundant; adapt only as a reminder for automation.

The protocol already says internal helpers are not participants and cannot own canonical files unless explicitly listed. Section 12 also keeps the driver as the sole side-effect actor. That is the important recursion boundary. We do not need Fusion-style depth headers in the human-readable protocol unless a future CLI implementation starts nesting automated parley runs.

### ExecPlan single self-contained living document

Verdict: reject as a replacement, adapt as artifact discipline.

Parley-deck should not collapse `round-01/`, later rounds, `consensus.md`, `FINAL.md`, `IMPLEMENTATION.md`, and review artifacts into one living file. The multi-file shape is what gives the deck independent thought, append-only signoffs, and collision avoidance.

But the self-containment principle is valuable inside the artifacts that are already intended to carry state forward: `FINAL.md` for design intent and `IMPLEMENTATION.md` for execution state.

### ExecPlan progress, surprises, decision log, outcomes

Verdict: adapt mainly into `IMPLEMENTATION.md`; selectively summarize in `FINAL.md`.

`consensus.md` already captures agreed decisions and trade-offs at design time. Duplicating a full Decision Log there would add ceremony. During implementation, however, decisions and surprises happen after consensus; today they can disappear into session memory, PR comments, or vague "Notes for reviewers." That is exactly the gap ExecPlan living sections solve.

These sections would also give section 13 a better evidence corpus. Retrospectives are only as good as the durable data they mine.

### ExecPlan concrete steps and validation

Verdict: adapt, with scope control.

For high-risk work, driver-managed work, and multi-file implementation, concrete steps and expected validation outputs are worth the overhead. For small design-only ideas, they should be optional or abbreviated. The protocol already warns against unnecessary multi-agent overhead; the same principle should apply to plan verbosity.

### ExecPlan idempotence and recovery

Verdict: adapt strongly for implementation and action/pipeline work.

Section 12 already has robust low-level idempotency for effects: stable keys, reconcile before retry, dry-run limitations, and human gates. That is stronger than generic ExecPlan advice for external side effects. The missing layer is task-specific recovery in `FINAL.md` / `IMPLEMENTATION.md`: what state matters, what partial progress looks like, what can be retried, and what requires a human gate.

### ExecPlan autonomy instruction

Verdict: reject where it conflicts with parley-deck gates.

"Do not ask for next steps; proceed to the next milestone" is useful for a single-agent coding plan. It is dangerous as a parley-deck rule. The protocol intentionally stops for quorum, blockers, human gates, provider capability failures, production risk, and protocol changes. The auto-drive driver may advance only inside approved boundaries; it must not turn plan prose into permission to skip consensus or gate files.

### ExecPlan narrative prose preference

Verdict: reject as a protocol preference.

Parley-deck artifacts need to be skimmed by multiple agents, humans, PR UIs, and possible tooling. Structured headings, bullets, and checklists are a strength here. Narrative is useful for orientation, but forbidding tables/checklists would make `IMPLEMENTATION.md`, review consensus, and driver parsing worse.

## What parley-deck already covers better

- Multi-agent diversity: named roster participants with durable files are stronger than an opaque automatic panel.
- Independence: Phase 1 explicitly forbids reading other round-01 files before writing one's own artifact.
- Accountability: every canonical artifact is authored by an agent ID; helper models do not get hidden ownership.
- Consensus: all-participant signoff is stronger than confidence by breadth of agreement.
- Audit trail: raw rounds, cross-review, consensus, final, implementation, review, and fix-up files survive context loss.
- Conflict avoidance: one file per agent per round and append-only signoffs are stronger than a merged hidden transcript.
- Side-effect safety: section 12's driver boundary, idempotency keys, reconcile rules, capability checks, and human gates are stronger than generic plan autonomy.
- Retrospective governance: section 13 already keeps harness improvement advisory, multi-agent, and human-approved.

## Explicitly not to adopt

- Do not replace `consensus.md` with a judge JSON artifact or give a judge role authority over participants.
- Do not let majority confidence override a `BLOCK` signoff or change signoff weight.
- Do not hide raw round files from `FINAL.md` behind a judge summary.
- Do not make Fusion-style automatic web search/fetch a default for every participant; cost, privacy, and prompt-injection risk are too high.
- Do not expand the active roster dynamically per call without the protocol's roster/quorum mechanics.
- Do not let internal helper agents create canonical artifacts unless they are listed participants.
- Do not make `FINAL.md` a living file after Phase 4; use a v2 idea when final design is invalidated.
- Do not collapse the deck into one ExecPlan file.
- Do not import "proceed without prompting" autonomy across human gates, production gates, protocol-change gates, or unresolved blockers.
- Do not create a "confidently wrong" review severity. Capture it as retro evidence, not as a blame label or merge blocker category.

## Suggested adoption path

First, trial the ExecPlan-inspired `IMPLEMENTATION.md` discipline on one medium or large implementation idea, especially one involving the auto-drive driver or pipeline blocks. Measure whether a fresh agent can resume from `FINAL.md` plus `IMPLEMENTATION.md` without reading session transcripts.

Second, trial a Fusion-inspired comparison block in one `consensus.md` and one `review/consensus.md`. Keep it advisory and let signoffs test whether participants find it accurate.

Third, feed the resulting artifacts into a future section 13 retro pass. If `parley retro` can extract better recurring failure modes from progress, surprises, decision logs, acceptance evidence, and blind-spot notes, then the additions are earning their keep. If the sections become boilerplate, keep them opt-in for risky or long-running work only.
