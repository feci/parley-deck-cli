---
from: user
to: claude-1
idea: meta-protocol-change-quota-auto-exclude
phase: round-02
blocking: no
date: 2026-10-03
---

## Owner decision: stale worktree registrations pruned (driver gap 6)

Relayed by the owner's Claude Code session. The owner's answer, verbatim:

> Question: "Should I remove the 13 stale git worktree registrations that stop the parley driver at every
> phase transition?"
> Selected: **"Yes, prune with a backup (Recommended)"**. The option read: "Back up .git/worktrees/<name> for
> all 13, then run `git worktree prune`. Their directories are already gone, and every commit is reachable
> from a branch, so nothing is lost and the backup restores it. The driver then works again and the organizer
> stops launching rounds by hand."

Done at 2026-10-03 about 13:58 CEST, in the shared repository `parley-deck-cli/.git`:

- All 13 admin directories were backed up first to
  `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/runs-handoff/worktree-prune-backup-20261003/`,
  together with `worktree-list-before.txt` and `worktree-list-after.txt`.
- Before pruning, the HEAD of each one was checked: 11 branch refs and 2 detached commits (`09e9ab3`,
  `c4a8b83`). All of them resolve, and both detached commits are contained in branch `designated-implementer`.
- `git worktree prune -v` removed exactly those 13: `checkout`, `checkout1` to `checkout10`,
  `revert-test2` and `f2repo`. Afterwards `git worktree list --porcelain | grep -c '^prunable'` returns 0,
  and all 21 remaining registered worktree paths exist.

What this means for the run: you may use `parley continue --auto --no-implement` again, as RESUME-BRIEF-1
step 4 says. If the driver now refuses for a different reason, record it verbatim as a new driver gap and
use the fallback. Do not prune or declare anything yourself. The driver's earlier note
`claude-to-user_…_driver-error.md` has its cause resolved by this decision. You may archive it under
`inbox/archived/` and record that in `organizer-notes.md`.
