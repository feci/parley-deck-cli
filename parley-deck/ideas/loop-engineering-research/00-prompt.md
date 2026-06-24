---
idea: loop-engineering-research
author: user
created: 2026-06-22
participants: [claude-1, codex-1, hermes-1, antigravity-1]
roles:
  claude-1: synthesis — map loop-engineering's 6 blocks onto Parley Deck's existing surface; protocol (§) change recommendations
  codex-1: CLI/engine lens — outer-loop driver, automations/triggers, stopping conditions (max iter/cost/duration), durable goal-state across runs
  hermes-1: maker/checker & verification lens — Phase 5/6 mapping, refutation-mode verifier, self-grading risk, conditional rigor
  antigravity-1: risk & guardrail lens — runaway cost, HITL fatigue, cognitive surrender, comprehension debt; what NOT to adopt
status: round-01
---

## Problem / idea

**Research "loop engineering" thoroughly and extract concrete inspiration for the
Parley Deck protocol (`COOPERATION.md`) AND the `parley-deck-cli` engine.**

"Loop engineering" is the June-2026 paradigm (successor to prompt → context →
harness engineering) crystallized by **Boris Cherny** (creator/head of Claude Code
at Anthropic) — *"I don't prompt Claude anymore. I have loops running that prompt
Claude and figuring out what to do. My job is to write loops."* — echoed by OpenAI's
Peter Steinberger and named/codified by Google's Addy Osmani.

Parley Deck is **already** a partial loop-engineering system (durable on-disk
artifacts, maker/checker via Phase 5 implementer vs Phase 6 reviewers, worktrees
addon, skills, a CLI driver that auto-advances phases). The question is **what the
loop-engineering frame teaches us that we are missing or could sharpen** — in the
protocol and in the CLI. Each participant leads with their lens but must touch both
**protocol** and **CLI** and end with **prioritized, concrete** recommendations
(a recommendation must name the file/section/feature it would change).

Output: a `FINAL.md` design doc — a prioritized backlog of protocol-change ideas and
CLI features, each with rationale, the loop-engineering principle it embodies, rough
cost/risk, and an explicit "adopt / adapt / reject" call. This idea is **design-only**
(no implementation in this idea); accepted items spin off their own ideas.

---

## Embedded research corpus (June 2026)

> Copied here so every participant works from the same source (some agents lack web
> access). Add your own knowledge; cite where you go beyond this corpus.

### Definition & lineage
- **Loop engineering** = designing systems that autonomously prompt and iterate with
  AI agents instead of prompting them turn-by-turn yourself. "You design the system
  that does it instead." (Osmani)
- Lineage: **Prompt engineering** (2022–24) → **Context engineering** (2025, Tobi
  Lütke; Anthropic "Effective Context Engineering for AI Agents") → **Harness
  engineering** (early 2026, the env a *single* agent runs in) → **Loop engineering**
  (June 2026, "one floor above the harness" — runs on timers, spawns helpers,
  self-feeds). The "factory model": the system that builds software itself.

### The core loop cycle (outer loop)
Discovery-and-execution, repeated on a schedule or until a goal condition holds:
1. An **automation** discovers work (reads logs, issues, commits).
2. A **triage skill** writes findings to persistent state (markdown / Linear).
3. For viable findings, **isolated worktrees** spawn sub-agents to draft fixes.
4. A **separate verifier** agent reviews against project skills + tests.
5. **Connectors** open PRs, update tickets, send notifications automatically.

Lifecycle table: **Discover → Hand → Verify → Persist → Decide**
| Phase | Action | Owner | Tools |
|---|---|---|---|
| Discover | trigger identifies work (PR/Slack/timer) | Automation/Trigger | GH Actions, cron, MCP events |
| Hand | maker sub-agent gets task + Skills + State + Context | Maker | Workflow, `/goal` |
| Verify | checker inspects output; **defaults to refutation** | Checker | separate model, adversarial prompts |
| Persist | loop state + results to durable storage | State layer | STATE.md, Linear, Jira, DB |
| Decide | termination condition evaluated; retry or exit | Orchestrator | max iters / max cost / max duration / goal |

### The six building blocks
1. **Automations / Trigger — "the heartbeat."** Scheduled (cron) or event-driven
   starts; no manual "run" button. Embeds the loop into infrastructure.
2. **Worktrees — collision-free parallelism.** Isolated git checkouts so parallel
   sub-agents don't collide. (Note: "Your review bandwidth decides how many you can
   actually run, not the tool.")
3. **Skills (SKILL.md / CLAUDE.md) — externalize intent.** Project knowledge in repo
   files, not re-stated in prompts. "The more shared context you don't have to
   restate, the more stable the loop." Reduces "intent debt."
4. **Plugins / Connectors via MCP — side-effect authority.** Permission to open PRs,
   update tickets, deploy. Moves agents from proposals to closed-loop action.
5. **Maker / Checker sub-agents — separate generation from verification.** Ideally
   different models. "The model that wrote the code is way too nice grading its own
   homework." Verifier defaults to refutation, retries until passing.
6. **Durable State — memory on disk, not context.** STATE.md / JSON / Linear / DB,
   read back each cycle. "The agent forgets, the repo doesn't." Survives
   context-window exhaustion across long-running loops.

### Stopping conditions & cost control (mandatory hard limits)
- **Max iterations** (e.g. 100 steps), **Max cost** (e.g. $10/loop), **Max duration**
  (e.g. 30 min). Whichever triggers first halts and escalates to human review.
  (Claude Code's Workflow has a `budget` param.)
- `/goal` directive: "keeps going until a condition you wrote is actually true, and
  after every turn a separate small model checks whether you are done."

### Inner vs outer loop
- **Inner loop:** single LLM + tool-call cycle (ReAct / Reflexion / CodeAct).
- **Outer loop (where loop engineering lives):** timer/event trigger → spawn
  maker/checker → verify in refutation mode → persist state → decide next.

### Risks / failure modes (Osmani + Oflight)
1. **Runaway-loop cost explosion** — loose termination compounds token use.
2. **Verifier mis-grading** — same model generating + verifying = optimistic self-assessment.
3. **HITL approval fatigue** — too many notifications → mechanical clicks → the human brake breaks.
4. **Loop brittleness** — passes benchmarks, fails production edge cases.
5. **Comprehension debt** — code ships faster than understanding grows.
6. **Cognitive surrender** — using loops to avoid thinking instead of to accelerate judgment.
- Osmani's rule: *"Verification is still on you… your job is to ship code you
  confirmed works."* / *"Build the loop. But build it like someone who intends to
  stay the engineer, not just the person who presses go."*

### Boris Cherny's practiced workflow (context)
- Treats Claude Code as infrastructure, not magic: memory files, permission configs,
  verification loops, formatting hooks. Optimizes for **throughput**, not conversation.
- Runs ~10–15 concurrent sessions (terminal tabs + browser + mobile).
- "Giving Claude a way to verify its work increases quality 2–3×."
- Aggressive durable project memory (shared CLAUDE.md, updated when the agent repeats
  a mistake). Reported impact: ~8× more code/day, Claude authoring >80% of merged
  production code by May 2026.

---

## What Parley Deck already has (don't re-propose; sharpen instead)
- **Durable state on disk:** canonical artifacts (`round-NN/`, `consensus.md`,
  `FINAL.md`, `IMPLEMENTATION.md`), append-only signoffs, inbox. ("The repo doesn't forget.")
- **Maker/checker:** Phase 5 implementer (default = FINAL drafter) vs Phase 6
  reviewers (every *non-implementer*) with fixed severities; Phase 7/8 fix-up loop
  until zero agreed fixes.
- **Worktrees:** the `parley-worktrees` addon (sibling `../worktrees/`, lock manifest).
- **Skills:** Parley Deck *is* a SKILL.md; `parley-tracker` externalizes ticket intent.
- **A driver / auto-advance:** the CLI `internal/driver` auto-advances past round-01
  and (opt-in `auto_implement`) drives Phases 5–8; `parley run` auto-drives by default.
- **Readiness ping (§9.0):** per-idea roster liveness (hosted-PONG) before an idea.
- **Retrospective optimization (§13):** advisory post-idea improvement signal.

## Research questions (each participant: touch all, lead with your lens)
1. **Mapping.** Which of the six blocks does Parley Deck already cover, partially
   cover, or miss? Be specific about the seam.
2. **Outer-loop / automations.** Parley today is human-triggered per idea. Should the
   CLI gain a *standing loop* mode — cron/event-triggered discovery (read commits /
   issues / CI logs) that **auto-opens ideas**? What's the protocol's role vs the CLI's?
3. **Stopping conditions.** Loop engineering mandates max-iters / max-cost / max-duration
   + a goal-checker. The driver auto-advances phases — does it have explicit budget /
   iteration / wall-clock ceilings and a "are we actually done" gate? Where do they belong
   (protocol invariant vs CLI flag vs `[defaults]` in `~/.parley`)?
4. **Refutation-mode verification.** Loop engineering says the verifier should *default
   to refutation* and use a *different model*. Parley's Phase 6 already uses non-implementer
   agents — should the protocol make "default-to-refute" explicit, and require a different
   model than the implementer? (We already learned this lesson — the Phase-6 false-green.)
5. **Durable goal-state across runs.** Loop engineering keeps a STATE.md the loop reads
   each cycle. Parley has per-idea state but no cross-idea/cross-run "what's next" goal
   ledger. Worth a `parley-deck/STATE.md` or backlog? Protocol or CLI?
6. **Guardrails.** Which loop-engineering risks (runaway cost, HITL fatigue, comprehension
   debt, cognitive surrender) does Parley's human-in-the-loop, consensus-gated design
   already mitigate — and which does an *automated* outer loop newly expose? What must we
   refuse to automate?
7. **Naming / framing.** Is there value in Parley Deck explicitly positioning itself as a
   *loop-engineering substrate* (README / marketing)? Cherny/Osmani give us vocabulary.

## Deliverable shape (for FINAL.md)
A prioritized table: `# | Recommendation | Protocol/CLI | LE principle | Effort | Risk |
Adopt/Adapt/Reject | Spin-off idea slug`. Plus short rationale per top item. Bias toward
a **small number of high-leverage** changes over a long wishlist. Reject loudly where
Parley's deliberate human-gate conflicts with naive automation.
