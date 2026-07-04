---
idea: parley-learn-playbooks
status: final
drafter: claude-1
track: deliberation
date: 2026-07-04
participants: [claude-1, codex-1, hermes-1, antigravity-1]
---

## Decision

Add `parley learn <closed-idea-slug>` — an advisory sibling to `parley retro` — that
distills a COMPLETED idea into a reusable playbook at `parley-deck/playbooks/<topic>.md`.
Playbooks are advisory (beside consults): never quorum, never overriding protocol.
Unanimous consensus (✅ ×4). Protocol change = one new §13.5 subsection.

## Design (ratified in consensus.md)

1. **`parley learn <slug> [--topic NAME] [--dir DIR]`** wired in `internal/app/app.go`
   → `runLearn` (`internal/app/learn.go`). Read-only over the idea; writes exactly ONE
   new playbook file, fail-closed if the target exists (`Lstat` symlink guard, strict
   kebab-case slug/topic — the `parley retro propose` safety pattern).

2. **Precondition:** the idea's `IMPLEMENTATION.md` frontmatter is `status: complete`.

3. **Output** = an advisory playbook with a fixed shape (frontmatter `playbook`,
   `distilled-from`, `distilled`, `status: advisory`; sections When to use / Proven
   shape / Step checklist / Gotchas & fixes / Verification pattern). v1 distillation is
   a deterministic skeleton seeded from the idea's frontmatter (track, participants,
   fix-up count) + the lifecycle; the human refines transferable prose before commit.

4. **Protocol** = §13.5 (a single paragraph + the `playbooks/` directory convention):
   playbooks are advisory/non-canonical like consults; `parley learn` is a tooling
   command, NOT a Parley round (advisory → no quorum needed; commit review is the gate).

5. **v1 scope:** one closed idea → one playbook. `--refresh`, cross-idea synthesis, and
   Phase-0 auto-suggestion are deferred follow-ups.

## Verification (done criteria)

- Tests: `parley learn` writes the playbook; rejects a non-complete idea; fails closed on
  an existing target; `--topic` after the slug works. `go build ./...`, `go vet`,
  `gofmt -l` clean. Drift guard green (both COOPERATION.md copies identical); skill
  fallback re-synced; `meta/version.json` protocolSha256 + changelog updated.

## Non-goals

No auto-application/injection of playbooks; no external-agent skill generation; no change
to the §13 retro findings flow.

## Signoffs

<!-- each participant appends its own block -->
