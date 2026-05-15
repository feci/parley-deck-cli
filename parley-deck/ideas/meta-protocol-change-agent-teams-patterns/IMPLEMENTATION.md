---
idea: meta-protocol-change-agent-teams-patterns
status: complete
implementer: codex
started: 2026-05-14
completed: 2026-05-14
branch: /Users/tomasfecko/AI_WORKSPACE/parley-deck/parley-deck-cli#main-working-tree
head-commit: uncommitted-working-tree-on-181fd08
design-pr: n/a
implementation-pr: n/a
---

## Summary of work

Implemented the user-approved points 1-5 from `FINAL.md`:

- Updated `parley-deck/COOPERATION.md` with advisory per-idea `roles:` metadata, participant sizing guidance, internal-helper accountability, Phase 5 plan-gate guidance, and inbox mirroring rules.
- Updated `parley-deck/meta/protocol-changelog.md` with the protocol change entry.
- Updated the Parley Deck skill source at `/Users/tomasfecko/AI_WORKSPACE/parley-deck/parley-deck-skill/SKILL.md`.
- Synchronized the skill fallback protocol snapshot at `/Users/tomasfecko/AI_WORKSPACE/parley-deck/parley-deck-skill/references/COOPERATION.md`.
- Updated the installed local skill at `/Users/tomasfecko/.codex/skills/parley-deck/SKILL.md`.
- Synchronized the installed skill fallback protocol snapshot at `/Users/tomasfecko/.codex/skills/parley-deck/references/COOPERATION.md`.

## Implementation plan / checklist

- Keep `participants:` as a list of IDs and add `roles:` as a sibling advisory map.
- State explicitly that roles do not affect quorum, signoff weight, artifact ownership, or drafter eligibility.
- State explicitly that internal helpers do not count as Parley Deck participants or satisfy the non-solo requirement.
- Add participant sizing guidance as process guidance, not a new artifact or phase.
- Add Phase 5 planning as a checklist inside `IMPLEMENTATION.md`, not a new phase.
- Keep inbox handoffs non-authoritative and require phase-affecting decisions to be mirrored into canonical artifacts.
- Avoid introducing nested sub-ideas, task boards, delegate rounds, locks, or hook-driven transitions.

## Deviations from FINAL.md

- Scope extended per direct user instruction to keep `/Users/tomasfecko/AI_WORKSPACE/parley-deck/parley-deck-skill/SKILL.md`, `/Users/tomasfecko/AI_WORKSPACE/parley-deck/parley-deck-skill/references/COOPERATION.md`, and the installed `/Users/tomasfecko/.codex/skills/parley-deck/` copies in sync with the new protocol text. Protocol intent unchanged.
- The active transport is `github-pr`, but this turn applied the approved changes directly in the local working tree rather than creating a feature branch and implementation PR. The working tree is uncommitted on top of `181fd08`; create a follow-up branch/PR if strict GitHub transport mirroring is required.
- The deferred items in `FINAL.md` remain deferred.

## Notes for reviewers

- Changes are documentation/protocol-only; no Go code was changed for this idea.
- The working tree already contained unrelated untracked review artifacts under `parley-deck/ideas/tui-agent-controls/review/round-02/`; they were not modified for this implementation.
- The installed skill under `/Users/tomasfecko/.codex/skills/parley-deck/` was updated directly because the user explicitly requested updating the active skill.

## Fix-up cycle 1

status: complete
completed: 2026-05-14
head-commit: uncommitted-working-tree-on-181fd08

### Fixes applied

- Updated `IMPLEMENTATION.md` metadata and deviations to truthfully describe the local uncommitted working tree and user-directed skill scope expansion.
- Clarified Phase 5 wording in `COOPERATION.md` so `IMPLEMENTATION.md` can be opened with a plan/checklist before code changes and finalized at completion.
- Added a small implementation plan/checklist example to the Phase 5 `IMPLEMENTATION.md` template.
- Added `IMPLEMENTATION.md` to the §8 inbox mirroring rule.
- Added the idea directory reference to `parley-deck/meta/protocol-changelog.md`.

### Deviations from agreed fixes

None.

## Re-review cycle 2

status: complete
completed: 2026-05-14
head-commit: uncommitted-working-tree-on-181fd08

### Review result

- `review/round-02/claude.md`: no remaining findings.
- `review/round-02/gemini.md`: no remaining findings.
- `review/round-02/hermes.md`: no remaining findings.

### Remaining agreed fixes

None.
