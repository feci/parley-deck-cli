---
idea: meta-protocol-change-github-pr-transport
drafted-by: codex
date: 2026-05-10
---

## Agreed decisions

- Switch active transport to `github-pr`.
- Preserve the existing local audit trail by committing and pushing all current source and `parley-deck/` artifacts.
- Treat future GitHub PR comments and reviews as ergonomic mirrors; files remain canonical.

## Agreed trade-offs

- Existing pre-remote work will enter GitHub as an initial commit rather than reconstructed PR history.
- Future Parley workflow should use merge commits, not squash merges, to preserve per-agent commits.

## Open items deferred to implementation

- Configure GitHub branch protection, labels, and merge strategy in repository settings.

## Signoffs

### Signoff: codex — 2026-05-10
Status: ACCEPT
Notes: User explicitly requested GitHub usage and remote push; this transport change is appropriate.
