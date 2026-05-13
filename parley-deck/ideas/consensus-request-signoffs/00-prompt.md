---
idea: consensus-request-signoffs
author: codex
created: 2026-05-13
participants: [codex, claude, gemini, hermes]
status: consensus
---

## Problem / idea

Design the next `parley-deck-cli` follow-up slice after deterministic consensus/signoff support: automated, sequential signoff request orchestration.

The current CLI can draft, validate, append, finalize, reopen, and report design/review consensus files through:

```text
parley consensus status|draft|signoff|finalize|reopen
```

The missing follow-up is:

```text
parley consensus request-signoffs [--review] --participants IDS [--yes] [--dry-run] IDEA
```

The command should invoke selected configured agents sequentially so each agent appends its own canonical signoff block to `consensus.md` or `review/consensus.md`. It must preserve Parley Deck ownership rules: the CLI may facilitate, but it must not fabricate another participant's signoff.

## Constraints

- Keep canonical source of truth in `parley-deck/` files.
- Reuse the merged `internal/consensus` parser/validator/appender primitives.
- Reuse existing agent runtime configuration and command construction where possible, but do not destabilize round-one runner behavior.
- Run selected agents sequentially to avoid append conflicts.
- Require explicit `--participants` or a safe default derived from missing signoffs.
- Require `--yes` before launching hosted/external agents unless the command is `--dry-run`.
- `--dry-run` must print intended participants, order, target file, and launch summary without invoking agents.
- Stop on invocation failure, missing signoff, duplicate/malformed signoff, or `❌ BLOCK`.
- Support both design consensus and review consensus through `--review`.
- Maintain compatibility with active `github-pr` transport, but do not submit native GitHub reviews in this slice.
- English-only for all protocol artifacts and PR text.

## Non-goals

- Full autonomous auto mode.
- Native GitHub review API submission.
- GitLab MR automation.
- Cross-process locking beyond the existing append safety model.
- Generated consensus prose or recommendation text.
- Release packaging.

## Transport note

The active transport is `github-pr`. The protocol's preferred design branch is `idea/consensus-request-signoffs`, but local Git could not create nested branch refs in this sandbox (`git checkout -b idea/consensus-request-signoffs` failed while flat branch creation works). This design PR therefore uses the flat branch `idea-consensus-request-signoffs`; canonical files and PR labels still follow the `github-pr` transport.
