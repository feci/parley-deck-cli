---
idea: meta-protocol-change-github-pr-transport
status: final
author: codex
consensus-date: 2026-05-10
participants: [codex]
---

## Final plan / specification

Use GitHub Pull Request transport for future Parley Deck workflow in this repository.

Repository:

- `https://github.com/feci/parley-deck-cli`

Initial bootstrap:

- Push the current local repository to GitHub.
- Include existing source files and the full current `parley-deck/` artifact tree.
- Use `main` as the integration branch.
- Do not squash future Parley PRs; use merge commits to preserve agent-level commit history.

Future workflow:

- Design branches use `idea/<slug>`.
- Implementation branches use `feature/<slug>`.
- Files under `parley-deck/` remain canonical.
- PR comments, requested reviews, labels, and approvals mirror the canonical files according to `COOPERATION.md` section 11.B.

## References

- Consensus: ./consensus.md
- Round: ./round-01/codex.md
