---
agent: codex
idea: meta-protocol-change-github-pr-transport
round: 1
date: 2026-05-10
---

## Summary

GitHub PR transport is appropriate now that the user created `https://github.com/feci/parley-deck-cli`. The first push should include the current Parley Deck audit trail as repository history, then future ideas can use branches and PRs according to section 11.B of `COOPERATION.md`.

## Proposed approach

- Update the transport header in `parley-deck/COOPERATION.md` to `github-pr`.
- Add a protocol changelog entry documenting the user-directed switch.
- Push the current repository, including `parley-deck/`, to the new GitHub remote.
- Use `main` as the integration branch unless the remote rejects it.

## Concerns / open questions

- The existing design and implementation work happened before the GitHub remote existed, so it cannot be mirrored as historical PR discussions without fabricating history.
- Branch protection and merge strategy still need to be configured in GitHub settings by a repository admin.

## Risks

- A single initial commit has coarse attribution for all pre-remote work. This is acceptable because the canonical agent artifacts preserve per-agent authorship inside `parley-deck/`.
