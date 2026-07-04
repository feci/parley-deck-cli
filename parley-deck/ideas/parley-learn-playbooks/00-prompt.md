---
idea: parley-learn-playbooks
author: user
created: 2026-07-04
track: deliberation
participants: [claude-1, codex-1, hermes-1, antigravity-1]
roles:
  claude-1: facilitation + §13 coherence
  codex-1: CLI/Go implementation shape
  hermes-1: protocol minimalism — smallest §13 extension that works
  antigravity-1: consumer-deck adoption — are playbooks useful outside this repo?
status: round-01
---

## Problem / idea

Inspired by Hermes Agent v0.18.0 `/learn`: point it at anything (a directory, a URL,
a workflow you just walked through) and it distills a reusable skill automatically.

Parley analogue: §13 (`parley retro`) analyzes closed ideas and proposes process
optimizations — but its output is advisory findings, not reusable assets. When a deck
ships its 5th protocol-change idea or 3rd release-burst, the hard-won shape of "how we
do X here" lives only in old idea directories and the facilitator's memory.

Proposal to deliberate:

- **`parley learn <closed-idea-slug>`** — distill a completed idea into a reusable
  **playbook**: `parley-deck/playbooks/<topic>.md` capturing the proven shape
  (phases actually run, roster+track used, checklist of steps, gotchas hit and their
  fixes, evidence/verification pattern) — generalized, with idea-specific details
  stripped.
- **Playbook usage**: when a new idea's brief matches an existing playbook topic, the
  facilitator (and driver prompts) reference it in Phase 0 so the new idea starts from
  the proven skeleton instead of from scratch.
- **§13 extension**: define playbooks as a first-class retro output — advisory,
  non-canonical for quorum (like consults), maintained via normal ideas when they need
  substantive revision.

## Constraints

- Playbooks are ADVISORY: never quorum evidence, never override protocol text, never
  auto-applied. Referencing one is optional context, like a consult.
- Distillation must be a real multi-step process with a human-visible result — the
  playbook is committed and reviewable; garbage-in protection is normal review of the
  commit.
- Additive protocol change to §13 only (both COOPERATION.md copies + skill fallback);
  no new phases, no new gates.
- English-only, no secrets (playbooks summarize process, not credentials/tokens).
- v1 distills from ONE closed idea; merging insights across ideas can come later.

## Non-goals

- No auto-generated skills for external agent CLIs (that is the skill repo's job).
- No playbook auto-injection into agent prompts without facilitator choice.
- No retro rewrite — §13 findings flow stays; playbooks are an additional output type.
