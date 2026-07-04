---
agent: hermes-1
idea: parley-learn-playbooks
round: 1
date: 2026-07-04
---

## Summary

The smallest §13 extension that works is a single new advisory output type: after `parley retro` produces findings, it MAY also distill a playbook file under `parley-deck/playbooks/<topic>.md`. No new phases, no new gates, no new quorum class — playbooks sit beside consults in the advisory tier. The protocol change is a paragraph in §13 plus the directory convention; everything else is facilitator practice.

## Proposed approach

Minimal protocol surface — three additions, all in §13:

1. One sentence in §13's output definition: retro findings MAY include a playbook distillation committed to `parley-deck/playbooks/<topic>.md`.
2. One sentence on status: playbooks are advisory — non-canonical for quorum, non-binding on protocol text, exactly like consults.
3. One sentence on lifecycle: substantive playbook revision goes through a normal idea, not in-place edits.

That is the entire protocol delta. The `parley learn <slug>` command and the playbook file shape (phases run, roster/track, step checklist, gotchas+fixes, evidence pattern) are convention and tooling, not protocol — they live in the skill fallback and CLI, not in COOPERATION.md. Keeping the protocol change to "retro may emit a playbook, it's advisory, revisions are normal ideas" makes it additive and trivial to review.

Phase-0 referencing is facilitator choice, not protocol: the facilitator who writes a brief can mention a playbook the same way they'd mention a prior idea or consult. No protocol text forces or even suggests it. This keeps the mechanism opt-in at every layer.

## Concerns / open questions

- Is one closed idea enough signal to generalize a playbook? v1 says yes and defers cross-idea merging — but a single-idea playbook can encode a shape that was idiosyncratic to that idea, not proven. Worth stating explicitly in the playbook frontmatter that it's distilled from N=1.
- Where do playbooks live that aren't repo-specific? The prompt says `parley-deck/playbooks/` — fine for this repo, but antigravity-1's lens asks about consumer decks. If a playbook is useful across decks, it's really a skill, and the non-goals already exclude that. Keep playbooks repo-local.
- Naming/collision: `<topic>.md` is a free-form slug. No registry, no namespace. Probably fine for v1 given low volume, but note it.
- Garbage-in: the prompt relies on normal commit review. That's adequate, but only if playbooks are actually reviewed — if they're seen as auto-generated boilerplate they'll get rubber-stamped. The "human-visible result" constraint helps; the facilitator should be the one who decides a playbook is worth committing, not the retro command by default.

## Risks

- Overreach: auto-apply. The moment a playbook is injected into a new idea's prompt without facilitator choice, it has de facto become protocol — advisory in name, binding in practice. The non-goals exclude this; the risk is drift toward it once playbooks exist and feel useful. The protocol text should say "referencing a playbook is optional" and mean it.
- Overreach: playbook-as-law. If a facilitator treats a playbook as the required shape for a topic ("we have a playbook, follow it"), it overrides the idea-brief-first flow and ossifies a single-idea pattern into a template. Advisory must stay advisory in practice, not just in text.
- Scope creep into retro rewrite. The constraint says §13 findings flow stays and playbooks are an additional output type. The risk is that "retro may also distill a playbook" quietly becomes "retro should always distill a playbook," turning every retro into a heavier artifact and shifting §13's center of gravity. Keep it MAY, not SHOULD.
- Maintenance debt. Playbooks that aren't revised through ideas rot and mislead. The lifecycle sentence (revisions go through normal ideas) is cheap to write and easy to ignore in practice; there's no enforcement. This is acceptable for advisory assets but worth naming as a known drift risk.
