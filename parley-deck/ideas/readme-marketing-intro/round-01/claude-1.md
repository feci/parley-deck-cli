---
agent: claude-1
idea: readme-marketing-intro
round: 1
date: 2026-06-20
---

## Summary

Open the README with a 3-part intro: (1) a one-sentence hook, (2) a "What you get"
feature cluster, (3) a "Standing on the shoulders of" inspirations block. Lead with
the *problem* it solves (ad-hoc multi-agent is a solo model talking to itself), then
the *shape* of the answer (a durable, auditable, vendor-neutral protocol), then the
*pedigree* (the prior art we studied and adapted). Keep it skimmable: bold leads +
short lines, not paragraphs.

## Proposed approach

Structure (for `parley-deck-cli/README.md`, before `## Install`):

```
# parley-deck-cli

> Real multi-agent deliberation with a durable audit trail — not one model
> pretending to be a committee.

**Parley Deck** is a vendor-neutral protocol (+ the `parley` CLI) for getting
several AI agents to actually *cooperate* on a hard change: each writes its own
analysis, they cross-review, reach recorded consensus, implement, and review the
implementation — every step a file you can read, diff, and resume.

### What's in the box
- **8-phase idea lifecycle** — kickoff → independent analysis → cross-review →
  consensus → FINAL → implementation → code review → fix-up. Append-only trail.
- **Non-solo by design** — per-agent canonical artifacts; one file per agent per
  round; no agent overwrites another.
- **Transport-agnostic** — local files, GitHub PRs, or GitLab MRs.
- **Vendor/model-agnostic roster** — Claude, Codex, Gemini, GLM, … by stable ID.
- **Readiness preflight** (`parley preflight`) — protocol freshness + a live ping
  of every agent before an idea starts.
- **Retrospective optimization** (`parley retro`) — advisory, quorum-gated mining
  of the deck's own history.
- **Pipelines & action stages** — supervised-first automation with gates,
  idempotency, and a durable effects ledger.
- **Living docs** — self-contained `FINAL.md` + an `IMPLEMENTATION.md` you can
  resume from the document alone.
- Plus a live TUI, auto-drive, agent supervision, and central `~/.parley` defaults.

### Inspired by (adopted & adapted)
- **OpenRouter Fusion** → consensus "compare-not-merge" blind-spots lens.
- **OpenAI ExecPlans (PLANS.md)** → resume-from-the-doc FINAL/IMPLEMENTATION.
- **RHO** → retro pass, but advisory & quorum-gated, not self-graded.
- **kindly** → review dispositions, no-suppression, supervision knobs.
```

## Concerns / open questions

- Length: I lean toward ONE tight "What's in the box" list (8–10 bullets) — more
  than that and nobody reads it. Codex should sanity-check each bullet maps to a
  real section/command.
- Skill README vs CLI README: the skill README should get a shorter variant framed
  around "install the protocol skill into your agents", reusing the inspirations block.
- Tone: a single pull-quote hook is enough; avoid hype adjectives ("revolutionary").
- Attribution wording: "inspired by / adopted & adapted" — must not imply endorsement
  by those projects.

## Risks

- Over-claiming: if a bullet oversells (e.g. "fully autonomous"), it's false — code
  mutation is gated. Keep claims hedged to what ships.
- Drift: the README intro will age as the protocol grows; keep it feature-cluster
  level, not exhaustive, so it stays true longer.
