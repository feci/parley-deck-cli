---
idea: readme-marketing-intro
drafter: claude-1
date: 2026-06-20
status: final
supersedes: consensus.md
---

## Purpose

Add a marketing-style intro to the top of both READMEs that sells what the Parley
Deck protocol is and credits its inspirations honestly. Documentation only.

## Context

Synthesized from four round-01 lenses (claude-1 structure, hermes-1 hook, codex-1
accuracy/feature-map, antigravity-1 lineage). Every feature claim maps to a real
COOPERATION.md section or shipped `parley` command; inspirations are credited as
adopted-and-adapted with a non-endorsement disclaimer.

## Deliverable A — `parley-deck-cli/README.md` (insert after `# parley-deck-cli`, before `## Install`)

```markdown
> **Real multi-agent deliberation with a durable, reviewable audit trail** — not
> three agents in three terminals, and not one model role-playing a committee.

`parley` runs **Parley Deck**: a transport-agnostic protocol for getting several
AI agents to genuinely cooperate on a hard change — plus the CLI that makes it
usable instead of just specified. Each agent writes its own analysis, they
cross-review, reach a recorded consensus, implement, and review the
implementation — every step a file you can read, diff, and resume.

**Why not just spawn three agents in three terminals?** Ad-hoc multi-agent gives
you no audit trail, no conflict discipline, no consensus step, and no way to
resume. Parley Deck ships those as protocol.

**Why not ask one model to play a committee?** That's solo reasoning in a costume —
there's no second voice to actually disagree, and it reintroduces the single-model
self-preference that quorum-gated review exists to defeat. Parley Deck is non-solo
by design.

### What you get

- **An 8-phase idea lifecycle** (§4) — kickoff → independent analysis → cross-review
  → consensus → `FINAL.md` → `IMPLEMENTATION.md` → code review → fix-up.
  Append-only, and resumable from the documents alone.
- **Non-solo by design** (§1) — stable agent IDs (§2), one file per agent per round
  (§6); no agent overwrites another.
- **Compare, don't merge** — the consensus "Comparison & blind spots" lens rates
  confidence by agreement and surfaces contradictions and blind spots instead of
  averaging them away.
- **Transport-agnostic** (§0, §11) — `local-dir`, `github-pr`, or `gitlab-mr`: the
  same protocol whether agents share a filesystem or review each other through a PR.
- **Vendor/model-agnostic roster** — Claude, Codex, Gemini, GLM, and more by stable
  ID (subject to the CLIs you have installed and authorized).
- **Readiness preflight** (§9.0, `parley preflight`) — protocol-freshness plus a
  live roster ping before an idea starts; any exclusion is user-confirmed.
- **Advisory retrospectives** (§13, `parley retro`) — a quorum-gated pass over the
  deck's own history that *proposes* improvements through the normal workflow,
  never applies them automatically.
- **Supervised automation** — a live TUI and an auto-drive driver advance protocol
  *phases*; agent supervision (watchdog, stall guard, validated-artifact-beats-
  nonzero-exit) catches hung agents; code implementation and side effects stay gated.

### Inspired by — adopted & adapted

Parley Deck didn't invent these ideas; it wired them into one repository-backed,
quorum-gated protocol:

- **OpenRouter Fusion** → the compare-not-merge consensus lens
  (confidence-by-agreement, blind-spots), applied to asynchronous multi-round
  markdown instead of a real-time API ensemble.
- **OpenAI ExecPlans / PLANS.md** → resume-from-the-doc state, split into a static
  `FINAL.md` and a living `IMPLEMENTATION.md` governed by review-consensus.
- **RHO (Retrospective Harness Optimization)** → §13 retro, but advisory-only and
  quorum-gated instead of single-model self-preference.
- **kindly** → strict gates, stopping judgment, no-suppression review dispositions,
  and artifact-wins supervision.
- **Preflight readiness** → §9.0 protocol-freshness and roster liveness before each idea.

*Reference to these projects is for attribution and lineage only; no endorsement,
sponsorship, or affiliation is implied.*
```

## Deliverable B — `parley-deck-skill/README.md` (insert after the H1, before existing body)

```markdown
> **Install the Parley Deck cooperation protocol into your AI agents** — real
> multi-agent deliberation with a durable audit trail, not one model role-playing
> a committee.

`parley-deck-skill` installs the vendor-neutral **Parley Deck** cooperation
instructions (and a fallback protocol snapshot) into supported agent runtimes —
Codex, Claude Code, Antigravity, Gemini, Hermes, or a custom skill directory — then
helps you check and sync project metadata. It teaches agents *how to participate*
in the protocol; the companion `parley` CLI orchestrates the runs.

**Parley Deck** is a transport-agnostic protocol where several agents genuinely
cooperate on a change: each writes its own analysis, they cross-review, reach a
recorded consensus, implement, and review the implementation — every step a file
you can read, diff, and resume.

### What the protocol gives your agents

- **An 8-phase idea lifecycle** — kickoff → independent analysis → cross-review →
  consensus → `FINAL.md` → `IMPLEMENTATION.md` → code review → fix-up; append-only
  and resumable from the documents alone.
- **Non-solo by design** — stable agent IDs, one canonical file per agent per round;
  no agent overwrites another.
- **Compare, don't merge** — a consensus "Comparison & blind spots" lens that rates
  confidence by agreement and surfaces blind spots instead of averaging them away.
- **Discipline that travels** — no-suppression review dispositions, strict gates,
  and pre-idea readiness checks, the same across every runtime.

### Inspired by — adopted & adapted

Parley Deck didn't invent these ideas; it wired them into one quorum-gated protocol:

- **OpenRouter Fusion** → the compare-not-merge consensus lens.
- **OpenAI ExecPlans / PLANS.md** → resume-from-the-doc `FINAL.md` + living `IMPLEMENTATION.md`.
- **RHO** → advisory, quorum-gated retrospective optimization.
- **kindly** → strict gates, stopping judgment, no-suppression dispositions, artifact-wins.
- **Preflight readiness** → protocol-freshness + roster liveness before each idea.

*Reference to these projects is for attribution and lineage only; no endorsement,
sponsorship, or affiliation is implied.*
```

## Validation

- Each "What you get" bullet maps to a COOPERATION.md section or shipped command
  (codex-1 feature-claim map). No ungated-autonomy claims; auto-drive advances
  phases, implementation stays gated.
- Inspirations credited as adopted-and-adapted + non-endorsement disclaimer.
- CLI vs skill command surfaces kept separate (the skill README claims only
  install/doctor/status/sync-project/uninstall/paths/--version via its existing body).
- English only; intro stays ~one screenful.

## Idempotence & recovery

Insertion is text-only and additive (prepend a section; existing Install/usage
unchanged). Re-running the implementation must not duplicate the intro block.

## Outcomes & retrospective

(to be filled by IMPLEMENTATION.md after the READMEs are written and pushed)
