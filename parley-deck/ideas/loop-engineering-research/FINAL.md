---
idea: loop-engineering-research
drafter: claude-1
date: 2026-06-22
status: final
supersedes: consensus.md
participants: [claude-1, codex-1, hermes-1, antigravity-1]
signoffs: claude-1 ACCEPT · codex-1 ACCEPT-WITH-RESERVATIONS · hermes-1 ACCEPT-WITH-RESERVATIONS · antigravity-1 ACCEPT
---

## Purpose

What "loop engineering" (the June-2026 successor to prompt → context → harness
engineering, crystallized by Boris Cherny, named by Addy Osmani) teaches the Parley
Deck protocol and the `parley-deck-cli` engine — and a **prioritized, spin-off-ready
backlog** of concrete changes. **Design-only:** nothing is implemented in this idea;
each accepted item below becomes its own idea (slugs given).

## What loop engineering is (one paragraph)

Stop prompting the agent turn-by-turn; instead **write the loop** that prompts it,
reads the result, decides if it's done, and re-prompts — on a schedule or until a goal
holds. Cherny: *"My job is to write loops."* The outer loop = **Discover → Hand →
Verify → Persist → Decide**, built from six blocks: **Automations/Trigger, Worktrees,
Skills, Connectors (MCP), Maker/Checker sub-agents, Durable State**. Its non-negotiable
discipline is **stopping conditions** (max iterations / cost / duration + a goal-checker)
and **refutation-mode verification by a different model** — because *"the model that
wrote the code is way too nice grading its own homework."* Its standing warning
(Osmani): the lever moved, but **you stay the engineer** — verification, comprehension,
and the merge decision remain human.

## Headline finding

**Parley Deck is already ~70% a loop-engineering system** — and the frame's biggest
payoff was using it to **audit our own engine**, which exposed three gates we *believe*
we have but don't:

1. **`strict_gate` is protocol prose the driver ignores** — the promised
   `ReadStrictGate` + `strict_gate_clean`/`closing_review_round` fields (from idea
   `review-gate-honesty`) were never built; `advanceReview` closes on
   `outstanding_agreed_fixes == 0` regardless.
2. **`RunChecks` is a no-op outside Go** — it runs `go test ./...` only if `go.mod`
   exists; every non-Go or design-only idea passes the pre/post-fix-up guard for free.
3. **The auto-loop closes on one agent's self-report** — no objective re-verification
   that the agreed fixes hold or that a single FINAL.md acceptance criterion was checked.

These are the real "false-green" seams, and they are where the cheapest, highest-value
work lives.

## Six-block mapping

| LE block | Parley | Note |
|---|---|---|
| Automations / Trigger | **MISSING** | `parley run` is human-triggered; nothing *opens* ideas (and `pipeline watch`'s `openRemediationIdea` opens them with `participants: []`, a non-solo violation). |
| Worktrees | **COVERED** | `parley-worktrees` addon. |
| Skills | **COVERED** | Parley *is* a SKILL.md; `parley-tracker`. |
| Connectors (MCP) | **PARTIAL** | `github-pr` writes PRs; no ingestion of issues/CI; §12 keeps side-effects gated (correct). |
| Maker / Checker | **STRUCTURAL, gates inert** | Phase 5 vs 6 by agent id (AF3 moved the consensus drafter off the implementer) — but confirmatory-not-refutational prompt, `strict_gate` unimplemented, model-diversity unenforced, Go-only checks. |
| Durable State | **PER-IDEA yes / CROSS-RUN no** | Strong per-idea artifacts; no cross-run "what's next" ledger (contested — see LE-12). |

## Prioritized backlog → spin-off ideas

> Each row = one future idea. Effort S/M/L; calls reflect the 4/4 consensus.

### Tier 1 — verification honesty (ship first; cheap, unanimous)
1. **LE-1 — Refutation-default review** *(Adopt, S)* → `meta-protocol-change-refutation-default`
   `BuildReviewPrompt` adopts an adversarial posture ("assume it's wrong; try to break
   each FINAL.md criterion; only report no-findings after stating what you tried");
   `ValidateReviewArtifact` requires a `## Refutation attempts` section so an
   empty-findings review must show its work. + one normative line in §Phase 6.
2. **LE-2 — Implement `strict_gate`** *(Adopt, M)* → `strict-gate-enforcement`
   Build the never-shipped `ReadStrictGate` + `strict_gate_clean`/`closing_review_round`;
   under `strict_gate:true`, require a fresh full-scope zero-finding review round before
   `Complete()`. (`internal/driver/impl.go`, `internal/app/driver_impl.go`, §Phase 8.)
3. **LE-3 — Model-diversity guard** *(Adopt, S)* → folds into LE-1
   Warn (configurable: refuse auto-complete) when every reviewer shares the implementer's
   `Model`. No-op in today's 4-model roster; a real guard for 2-agent same-model decks.
4. **LE-4 — Generalize `RunChecks` + fix "artifact-wins"** *(Adopt, S-M)* → `generalize-runchecks`
   Read a `checks:` command from frontmatter / `~/.parley [defaults]`; **fail closed**
   for code-writing auto ideas with no checks; and tie the "artifact-wins" fix-up
   override (`phase58.go:128-139`) to `RunChecks` passing, not shape alone (hermes #8).

### Tier 2 — stopping conditions / budgets (the safety floor)
5. **LE-5 — Unified loop-budget contract** *(Adopt, M)* → `driver-loop-budgets`
   One `LoopBudget` (max driver-steps / rounds / fixup-cycles / wall-clock / cost);
   **budget hit = escalate, never complete**. Consolidates today's scattered caps
   (`loop.go:16` hard-coded 30m, `MaxRounds`, `MaxFixupCycles`, `pipeline maxCycles=3`).
   Protocol says *that* limits exist; CLI flags + `~/.parley [defaults.loop]` give numbers.
6. **LE-6 — Best-effort cost telemetry** *(Adapt, M-L)* → folds into LE-5
   Emit `agent.usage`/`loop.budget` events where the CLI exposes tokens; strict mode
   treats unknown cost as a halt. Don't pass provider spend-caps that abort mid-artifact.
7. **LE-10 — Fix `openRemediationIdea` quorum** *(Adopt, S)* → `remediation-idea-quorum-fix`
   Watcher-created ideas must be non-active candidates until they satisfy Phase-0 quorum
   (today `participants: []` violates non-solo). Do this before any trigger expansion.

### Tier 3 — close-decision integrity
8. **LE-7 — Goal-done gate** *(Adopt, run-once form, M)* → `goal-done-gate`
   Before `Complete()` under `auto_implement`/`strict_gate`, evaluate FINAL.md observable
   acceptance criteria via the existing advisory `consult.go` (the "separate small
   model"). Run **once before close**, not every tick (cost/fatigue).
9. **LE-11 — HITL-fatigue guardrails** *(Adopt, S-M)* → `hitl-fatigue-guardrails`
   Batch/rate-limit driver-opened questions; under `auto_implement` don't auto-close on
   `TriageReserved` (escalate reservations); refuse auto-complete with `< 2` reviewers.

### Tier 4 — the outer loop (only on top of Tiers 1-3)
10. **LE-8 — Human-brake invariant** *(Adopt, S)* → `automation-human-brake`
    New §: an automated loop may **discover + draft (Phase 0/1) only**; never push to
    quorum, implement, land/merge, finalize, modify the roster, or override/bypass
    consensus without a recorded human or full-quorum gate.
11. **LE-9 — `parley loop tick`** *(Adapt, L)* → `standing-loop-watch-mode`
    One-shot, scheduler-friendly (cron/GH-Actions/MCP) — **not a daemon**: discover
    candidates (commits/CI/issues), dedupe, write provenance into a *candidate* prompt,
    optionally call `parley run` within the LE-5 budget. Disabled by default;
    human-confirm mandatory; gated by LE-8.

### Contested — defer to its own idea
12. **LE-12 — Cross-run durable goal state** *(CONTESTED → `durable-backlog-ledger`)*
    claude-1 + codex-1 **Adopt** (a `parley-deck/STATE.md` and/or `parley goals` CLI,
    append/update-disciplined); antigravity-1 **Rejects** the shared-mutable-file form
    (context bloat + write conflicts; `parley-tracker` already maps external work). The
    *need* is real; the *home* (new file vs tracker vs `runs/`-derived view) is unresolved.
    **Build only after LE-5/8** — a backlog with no bounded loop to consume it is premature.

### Rejected (unanimous)
- A fully autonomous **discover → implement → merge/deploy daemon**. §12 gates stay;
  merge/release, roster edits, consensus overrides are never automated.
- **Uniform** refutation rigor on trivial/design-only ideas — keep **conditional rigor**
  (default `outstanding_agreed_fixes == 0` close stays for low-risk ideas; refutation /
  strict-gate scale with `auto_implement`/`strict_gate`).

### Framing (deferred to `readme-marketing-intro`)
Position Parley Deck as a **"loop-engineering substrate — the outer loop with a
human-gated consensus brake"** — but only after LE-1/2/5/8 land, so the claim is true.

## Recommended sequencing
1. **LE-1..4** (verification honesty) — closes the false-green seams the audit found.
2. **LE-5 + LE-10** (budgets + remediation quorum) — the safety floor.
3. **LE-7 + LE-11** (goal-done + HITL guardrails) — close-decision integrity.
4. **LE-8 → LE-9** (human-brake §, then the trigger) — the outer loop, gated.
5. **LE-12** in its own idea once 1-4 exist.

## Provenance
Four independent round-01 lenses (claude-1 synthesis, codex-1 CLI/engine, hermes-1
maker/checker, antigravity-1 risk/guardrails), each tracing recommendations to exact
files/lines, over a June-2026 web research pass (Cherny via The New Stack/Digg; Osmani
`addyosmani.com/blog/loop-engineering`; Oflight six-block deep-dive; Greyling;
explainx/aibuilderclub guides). Consensus 4/4 (no BLOCK).
