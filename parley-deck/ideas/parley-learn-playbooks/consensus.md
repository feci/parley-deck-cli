---
idea: parley-learn-playbooks
drafted-by: claude-1
date: 2026-07-04
track: deliberation
participants: [claude-1, codex-1, hermes-1, antigravity-1]
---

## Agreed decisions

Strong convergence: `parley learn` is a small advisory sibling to `parley retro`, not a
new phase or gate. Playbooks sit in the advisory tier beside consults — never quorum
evidence, never overriding protocol text.

1. **`parley learn <closed-idea-slug>`** — a top-level command wired in
   `internal/app/app.go` → `runLearn` (`internal/app/learn.go` owns CLI flags + the
   narrow write boundary), mirroring `parley retro propose`: strict slug validation,
   reject an existing target, `Lstat` symlink guard, and a test that ONLY the allowed
   file is written. `internal/learn` (or a reuse of `internal/retro`) owns read-only
   mining of the closed idea's artifacts.

2. **Precondition:** the idea's `IMPLEMENTATION.md` is `status: complete` (only closed
   ideas are distilled).

3. **Output = one advisory playbook** `parley-deck/playbooks/<topic>.md` with a fixed
   shape (frontmatter `playbook`, `distilled-from`, `distilled`, `status: advisory`;
   sections: When to use / Proven shape / Step checklist / Gotchas & fixes /
   Verification pattern). Idea-specific specifics stripped; transferable process kept.
   `--topic` sets the filename (else inferred).

4. **Advisory status (protocol):** a single §13 paragraph + the `playbooks/` directory
   convention. Playbooks are advisory like consults: referencing one in Phase 0 is
   optional context, never quorum, never overrides protocol; substantive revision goes
   through a normal idea. No new phase, gate, or quorum class.

5. **Distillation is a tooling command, not a Parley round** — it does not need quorum
   (the playbook is advisory and committed, so normal commit review is the quality gate).
   Explicitly NOT a solo-Parley: `parley learn` is a CLI utility, distinct from an idea.

6. **v1 distills from ONE closed idea.** Cross-idea merging + `--refresh` re-distill are
   deferred.

## Deferred follow-ups

- `parley learn --refresh` to re-distill a stale playbook after protocol changes.
- Cross-idea playbook synthesis (merge insights from several ideas into one topic).
- Facilitator auto-suggestion of a matching playbook at Phase 0.

## Dismissed / non-goals

- No auto-application or auto-injection of playbooks into agent prompts.
- No external-agent skill generation (that is the skill repo's job).
- No change to the §13 retro findings flow — playbooks are an additional output type.

## Signoffs

<!-- each participant appends its own block -->

### Signoff: claude-1 — 2026-07-04
Status: ✅ ACCEPT
Drafter. `parley learn` = advisory sibling to `parley retro`; playbooks in the advisory
tier beside consults; one §13 paragraph + playbooks/ directory. v1 = one closed idea.
