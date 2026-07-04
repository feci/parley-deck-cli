---
agent: claude-1
idea: parley-learn-playbooks
round: 1
date: 2026-07-04
---

## Summary

Add playbooks as a second, advisory output type of §13 — distilled from ONE closed
idea via `parley learn <slug>`, written to `parley-deck/playbooks/<topic>.md`,
explicitly non-canonical (like consults). The distillation is a real agent pass over
the idea's artifacts; the result is a committed, reviewable markdown asset, not a
hidden model memory. Keep the protocol delta to a single additive §13 sub-section.

## Proposed approach

### 1. What a playbook is
A generalized, idea-agnostic "how we do X here" doc with a fixed shape:
```markdown
---
playbook: protocol-change
distilled-from: ideas/meta-protocol-change-devx-speed
distilled: 2026-07-04
status: advisory        # never quorum evidence
---
## When to use
## Proven shape (phases, roster, track actually used)
## Step checklist
## Gotchas hit & fixes
## Verification pattern (done_when / evidence if used)
```
Idea-specific specifics (exact file names, one-off decisions) are stripped; the
transferable process is kept.

### 2. `parley learn <closed-idea-slug>`
- Precondition: the idea's `IMPLEMENTATION.md` frontmatter is `status: complete`
  (only closed ideas are distilled — matches non-goal).
- Mechanics: the command gathers the idea's 00-prompt, consensus, FINAL,
  IMPLEMENTATION, review/consensus, and asks the facilitator agent to distill the
  playbook to the fixed shape. Output path derived from a `--topic` flag or inferred
  from the idea's meta.
- It is a facilitator-run distillation (one agent), NOT a Parley round — the playbook
  is advisory, so it does not need quorum. It IS committed, so normal review of the
  commit is the quality gate. State this explicitly to avoid the "solo Parley" trap:
  `parley learn` is a tooling command, not a Parley idea.

### 3. §13 protocol text (both COOPERATION.md copies + skill fallback)
Add a sub-section "Playbooks (distilled retro output)":
- defines the artifact + advisory status (never quorum, never overrides protocol,
  referencing one is optional Phase-0 context like a consult);
- substantive revision of a playbook that changes recommended process goes through a
  normal idea (playbooks describe, they do not legislate);
- one cross-ref from §8 consults note ("consults and playbooks are both advisory,
  non-canonical").

### 4. Usage in Phase 0
Facilitator MAY cite a matching playbook in a new idea's 00-prompt (`see
playbooks/protocol-change.md`). No auto-injection; the human/facilitator chooses.

## Concerns / open questions

1. **Overlap with §13 retro findings** — retro already produces advisory findings.
   Playbooks are the reusable-asset form; keep them distinct: findings critique a
   past idea, a playbook templates the next one. State the distinction in the text.
2. **Staleness** — a playbook distilled from an old idea can rot when the protocol
   changes. Mitigation: `distilled-from` + `distilled` date in frontmatter; a playbook
   is advisory so staleness is low-risk, but `parley learn --refresh` could re-distill.
3. **Where does `parley learn` get the agent?** Reuse the facilitator/driver agent
   resolution already used by `parley retro`. If retro is advisory-only today, learn
   piggybacks on the same plumbing.

## Risks

- Scope creep into "auto-apply playbook to new idea" — explicit non-goal; keep it a
  referenced doc only.
- Playbook sprawl (one per idea) — mitigate by topic-keying (`protocol-change`,
  `release-burst`, `tui-feature`) so related ideas refine ONE playbook, not spawn N.
- Distillation leaking specifics/secrets — normal commit review + English/no-secrets
  invariant catch it; the distiller prompt must say "generalize, strip specifics".
