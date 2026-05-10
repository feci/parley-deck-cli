---
idea: meta-protocol-change-github-pr-transport
author: user
created: 2026-05-10
participants: [codex]
status: final
---

## Problem / idea

The repository now has a GitHub remote at `https://github.com/feci/parley-deck-cli`. Switch future Parley Deck coordination from local-directory-only mechanics to GitHub Pull Request transport while preserving the canonical `parley-deck/` files in the repository history.

## Constraints

- Preserve existing local Parley Deck history in Git.
- Keep files canonical; GitHub PR conversations are ergonomic mirrors.
- Use merge commits rather than squash merges for future Parley PRs.
- Do not rewrite existing local artifacts.

## Non-goals

- Do not migrate existing closed ideas into historical PRs.
- Do not require branch protection configuration in this initial push.
