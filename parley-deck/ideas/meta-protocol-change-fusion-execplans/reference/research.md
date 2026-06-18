---
idea: meta-protocol-change-fusion-execplans
kind: reference
drafted-by: claude (facilitator)
date: 2026-06-18
sources:
  - https://openrouter.ai/docs/guides/routing/routers/fusion-router
  - https://openrouter.ai/docs/guides/features/plugins/fusion
  - https://openrouter.ai/blog/announcements/fusion-beats-frontier/
  - https://developers.openai.com/cookbook/articles/codex_exec_plans
---

# Research: OpenRouter Fusion + OpenAI ExecPlans (PLANS.md)

This is the shared evidence corpus for the deliberation. It records the two
external systems **faithfully** (Part A, Part B), then lists **neutral questions**
for the panel (Part C). The mapping to parley-deck is deliberately left thin so
each participant forms an independent view in round-01.

---

## Part A — OpenRouter Fusion

**One-line:** Fusion turns a single API call into a *multi-model deliberation*:
a panel of diverse models answers in parallel, a judge model **compares (does not
merge)** their answers into a structured analysis, and a final model writes the
answer using that analysis.

### A1. Pipeline (5 stages)
1. **Request entry** — caller sends `model: "openrouter/fusion"`; the router
   resolves the alias to a concrete model and attaches the `openrouter:fusion`
   server tool.
2. **Model decision** — the outer model decides whether deliberation is warranted;
   it may answer directly unless `tool_choice: "required"`. (Tactical prompts skip
   fusion entirely.)
3. **Panel execution** — up to **8 models** answer the prompt **in parallel**, each
   with `openrouter:web_search` + `openrouter:web_fetch` enabled. Independent
   responses, generated simultaneously.
4. **Judge analysis** — a judge model receives **all** panel responses and
   **compares them — it does not merge them.** Returns structured JSON analysis.
5. **Final answer** — the outer model writes the final answer *using the judge's
   structured analysis* (not the raw panel text).

### A2. The judge's structured analysis — the core idea
The judge emits JSON with these dimensions:
- **Consensus / agreement** — points *all or most* models agree on, explicitly
  **"treated as higher-confidence consensus."** (Confidence scales with breadth of
  agreement.)
- **Contradictions / disagreements** — where panel answers conflict and need
  reconciliation.
- **Partial coverage** — what only *some* models covered (incomplete treatment).
- **Unique insights** — distinctive contributions from an individual model.
- **Blind spots** — *issues no model addressed at all.* (A structured "what did
  everyone miss" lens.)

### A3. Configuration
- `analysis_models` — panel composition; **1–8 models**. Default = "Quality preset"
  (`~anthropic/claude-opus-latest`, `~openai/gpt-latest`, `~google/gemini-pro-latest`).
- `preset` — curated slugs, e.g. `general-high`, `general-budget`.
- `model` — the judge model (defaults to the first preset model / outer model).
- `max_tool_calls` — web search/fetch iteration cap per model, default `8`, range 1–16.
- `enabled` — toggle.
- **Recursion protection** — depth tracked via `x-openrouter-fusion-depth`; panel
  and judge models **cannot** recursively invoke fusion (single deliberation level).

### A4. Cost
~**4–5× a single completion** with the default 3-model panel (N panel calls + 1
judge call). Scales linearly with panel size.

### A5. Design rationale (from the announcement)
- **Diversity wins:** "Bringing multiple different perspectives to complex problems
  yields superior results" — mirrors human teams. Panels consistently outperform
  individual models (e.g. Fable 5 + GPT-5.5 fusion 69.0% vs Fable 5 alone 65.3%).
- **Synthesis is itself valuable:** even fusing a model *with itself* (Opus 4.8 ×
  Opus 4.8) gained **+6.7 points** — the *synthesizer/judge step* adds value beyond
  model diversity.
- **Cost-efficiency via diversity:** budget-model panels approximated frontier
  quality at ~half the cost.
- **Benchmark hygiene:** evaluated on DRACO (100 deep-research tasks, 10 domains)
  with a rubric that has **negative-weight criteria** — "confidently stating wrong
  things gets punished" — so synthesis quality is tested, not verbosity.

---

## Part B — OpenAI ExecPlans (PLANS.md)

**One-line:** An ExecPlan is a **single, self-contained, living Markdown design
document** that lets a *stateless* agent (or a human novice) carry a complex,
multi-hour task from design through implementation and **resume from the document
alone**, with no external context.

### B1. Problem & workflow
- Targets **complex multi-hour tasks** needing research + design + implementation.
- A thorough self-contained doc is written *before* execution; the agent then
  **continuously updates it** while working; work can resume from the plan alone.
- Wired in via **AGENTS.md** shorthand, e.g.:
  `When writing complex features or significant refactors, use an ExecPlan (as
  described in .agent/PLANS.md) from design to implementation.`
- Autonomy instruction: *"do not prompt the user for 'next steps'; simply proceed to
  the next milestone. Keep all sections up to date."*

### B2. Structure (one fenced `md` block)
**Living sections (non-negotiable, updated continuously):**
- **Progress** — checklist with ISO timestamps.
- **Surprises & Discoveries** — unexpected findings *with evidence*.
- **Decision Log** — `Decision / Rationale / Date·Author`.
- **Outcomes & Retrospective** — summary at major milestones.

**Content sections:**
- **Purpose / Big Picture** — user-visible behavior enabled.
- **Context and Orientation** — repo state, file paths, term definitions.
- **Plan of Work** — prose sequence of edits (files, locations, insertions).
- **Concrete Steps** — exact commands, working dirs, expected output.
- **Validation and Acceptance** — how to *observe* success (behavior-focused).
- **Idempotence and Recovery** — safe retry/rollback paths.
- **Artifacts and Notes** — concise transcripts, diffs, snippets.
- **Interfaces and Dependencies** — specific libraries, modules, signatures.

**Progress format:**
```
- [x] (2025-10-01 13:00Z) Example completed step.
- [ ] Example incomplete step.
- [ ] Example partial step (completed: X; remaining: Y).
```
**Decision Log format:**
```
- Decision: …
  Rationale: …
  Date/Author: …
```

### B3. Lifecycle
- **Draft:** skeleton → research → embed findings in prose → proof-of-concept
  milestones for high-risk areas → fully self-contained.
- **Execute:** read the whole plan → run concrete steps, observe outputs → update
  Progress with timestamps at every stopping point → record discoveries → commit
  frequently.
- **Resume:** plan has all context to restart; next agent/human reads top-to-bottom,
  continues from the last Progress checkpoint; updates Decision Log if course changes.
- **Complete:** write Outcomes & Retrospective (achievements, gaps, lessons); the
  plan becomes reference documentation.

### B4. Principles
- **Self-containment is paramount** — "all knowledge and instructions needed for a
  novice to succeed"; define every term; embed library knowledge; don't point to
  external blogs or prior plans; repeat assumptions.
- **Narrative prose over lists** — "Prefer sentences over lists. Avoid checklists,
  tables, and long enumerations unless brevity would obscure meaning." Each milestone
  reads "as a story: goal, work, result, proof."
- **Behavior-focused acceptance** — e.g. "after starting the server, navigating to
  http://localhost:8080/health returns HTTP 200 with body OK"; show test commands +
  expected pass counts; demonstrate impact beyond compilation.
- **Idempotency & recovery** — "run multiple times without causing damage or drift";
  safe retry/rollback; backups for destructive changes; prefer additive, testable,
  incrementally-validated changes.
- **De-risking** — explicit prototyping milestones for high-risk areas; feasibility
  spikes evaluated independently; parallel implementations during large migrations.
- **Living-document discipline** — update all sections at every stopping point;
  record decisions + rationale contemporaneously; note plan changes at the bottom
  with reasoning; the plan must always be sufficient to restart alone.
- **Anti-patterns:** undefined jargon; "compiles but does nothing meaningful";
  leaving ambiguity unresolved (resolve it *in the plan* and explain why).

---

## Part C — Neutral questions for the panel

The deliberation question is: **which of these concepts, if any, could improve the
parley-deck protocol — and which are already covered, redundant, or a net negative?**
parley-deck is *itself* a multi-agent deliberation + artifact protocol, so several
ideas may already exist here in a stronger form. Be honest about overlap.

Prompts to consider (not a checklist to satisfy — pick what matters, add your own):

1. **Judge "compare-don't-merge" + structured analysis.** parley-deck's
   `consensus.md` is a *drafted synthesis*. Is there value in a structured analysis
   artifact with explicit **consensus / contradictions / partial-coverage /
   unique-insights / blind-spots** dimensions? Is the **"blind spots none addressed"**
   lens a genuine gap in our review/consensus phases?
2. **Confidence-by-breadth-of-agreement.** Our signoffs are binary
   (✅ ACCEPT / ❌ BLOCKER). Should consensus carry a graded confidence based on how
   many participants agree, and should that change anything downstream?
3. **Synthesis-as-distinct-value.** Fusion shows the judge step helps *even without
   diversity*. Does parley-deck under-use a dedicated synthesis/judge role distinct
   from the consensus *drafter*? Tradeoffs vs. our append-only signoff model?
4. **Negative-weight / "confidently wrong" failure mode.** DRACO punishes confident
   wrongness. Do our review severities (CRITICAL/MAJOR/MINOR/NIT) or §13 retro signals
   capture "an agent was confidently wrong"? Worth a dedicated lens?
5. **ExecPlan self-containment for FINAL.md / IMPLEMENTATION.md.** Could making these
   *self-contained living ExecPlans* let a fresh headless agent (codex/agy/hermes) —
   or the auto-drive **driver** — resume implementation **from the artifact alone**?
   This is exactly the cross-invocation resumption the driver needs.
6. **ExecPlan living sections → §13 retro evidence.** Decision Log / Surprises &
   Discoveries / Outcomes & Retrospective are structured signals. Would adopting them
   in IMPLEMENTATION.md feed the `parley retro` evidence corpus far better than what
   we mine today?
7. **Idempotence & Recovery section → auto_implement safety.** Would an explicit
   idempotency/recovery section in FINAL.md harden the gated auto-drive
   implementation phase (clean-tree, no-land)?
8. **Behavior-focused acceptance criteria.** Should FINAL.md require observable
   acceptance criteria (à la ExecPlans) that Phase 6 review and the driver check?
9. **Cost/overhead realism.** Fusion costs 4–5×; ExecPlans add doc-maintenance
   overhead. parley-deck already pays a multi-agent tax. Which (if any) of these are
   worth the added weight, and which would just add ceremony?
10. **Anything we should explicitly NOT adopt**, and why (e.g. things that conflict
    with our append-only / one-file-per-agent / human-gate invariants).

> Note for the panel: this is a **design/brainstorm round for human review** — no
> protocol text is being changed yet. Per §7, any actual protocol change needs its
> own ratified meta-protocol-change idea + human approval. Output hypotheses and a
> clear recommendation, not edits.
