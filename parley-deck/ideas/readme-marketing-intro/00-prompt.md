---
idea: readme-marketing-intro
author: user
created: 2026-06-20
participants: [claude-1, codex-1, hermes-1, antigravity-1]
roles:
  claude-1: narrative & structure (what story the intro tells)
  codex-1: technical accuracy (features must map to real protocol sections)
  hermes-1: positioning & differentiation (why this vs. ad-hoc multi-agent)
  antigravity-1: inspirations & attribution (credit the prior art faithfully)
status: final
---

## Problem / idea

Add a **marketing-style intro section** to the top of `README.md` (parley-deck-cli,
and an adapted version for parley-deck-skill) that sells what the **Parley Deck
protocol** is, what it includes, and **what it was inspired by**. Today the README
opens with a dry one-liner; we want a compelling, accurate overview that a newcomer
reads first.

Each participant proposes: (a) the framing/structure of the intro, (b) the concrete
feature bullets to highlight, and (c) how to credit the inspirations honestly
(adopted-and-adapted, not "we invented this"). Keep it punchy but truthful — every
claimed feature must map to a real protocol section or shipped command.

## What the protocol actually includes (ground truth — cite these, don't invent)

- **Transport-agnostic** coordination: `local-dir`, `github-pr`, `gitlab-mr` (§0, §11).
- **8-phase idea lifecycle** (§4): Kickoff → independent analysis → cross-review →
  consensus → finalization (`FINAL.md`) → implementation (`IMPLEMENTATION.md`) →
  code review → fix-up. Durable, append-only audit trail.
- **Non-solo by design** (§1): real multi-agent or it doesn't count; per-agent
  canonical artifacts; one-file-per-agent-per-round conflict avoidance (§6).
- **Vendor/model-agnostic roster** (§2): stable agent IDs; any CLI/model.
- **§9.0 Pre-idea readiness check**: protocol freshness + roster liveness ping
  (hosted-PONG), user-confirmed exclude/re-include. CLI: `parley preflight`.
- **§12 Pipeline blocks & action stages**: idea→...→monitoring pipeline, gates,
  supervised-first autonomy, idempotency/reconcile, durable effects ledger.
- **§13 Retrospective optimization**: advisory-only retro pass over the deck's own
  history; quorum replaces single-model self-preference. CLI: `parley retro`.
- **`FINAL.md`** static self-contained design sections + **`IMPLEMENTATION.md`**
  living execution doc (Progress / Surprises & Discoveries / Decision Log / Outcomes).
- **Consensus "Comparison & blind spots" lens**: compare-not-merge across agents
  (consensus / contradictions / partial coverage / unique insights / blind spots);
  confidence-by-agreement.
- **Agent supervision**: first-output watchdog, stall guard, failure classification,
  "validated artifact beats nonzero exit code"; `parley consult` advisory questions.
- **The `parley` CLI**: live TUI, auto-drive driver, `preflight`, `retro`,
  central `~/.parley/agents.toml` defaults (model + reasoning + `[defaults]` policy).

## Inspirations (credit faithfully — adopted & adapted, with the deck's own twist)

- **OpenRouter Fusion** → the consensus "Comparison & blind spots" lens
  (compare-not-merge, confidence-by-agreement, synthesis-as-value).
- **OpenAI ExecPlans / PLANS.md** (codex_exec_plans cookbook) → self-contained
  `FINAL.md` + living `IMPLEMENTATION.md` that you can resume from the doc alone.
- **RHO (Retrospective Harness Optimization)** → §13 retro pass, but advisory-only
  and quorum-gated instead of single-model self-preference.
- **kindly** → review dispositions / no-suppression, strict gates, stopping
  judgment, supervision knobs, artifact-wins.
- **Preflight readiness** → §9.0 freshness + liveness ping before each idea.

## Constraints

- English only (all `parley-deck/` files + README text).
- Truthful: no feature claim that isn't backed by a real section/command.
- Marketing tone OK, but concise — target a tight intro, not a wall of text.
- Inspirations credited as adopted-and-adapted; name the source, note our twist.
- Must read well as the FIRST thing in the README (before Install).

## Non-goals

- Not rewriting the whole README (Install/usage stays).
- Not adding new protocol features — this is documentation only.
- No new badges/CI claims we can't back.
