---
agent: codex-1
idea: addon-manifest-coverage
review-round: 4
date: 2026-08-02
reviewed-commit: 065985e
---

## Ruling on the managed / round-13 conflict

Do not block the inversion. The earlier assertion was ratified, but its `managed !== true`
half coupled payload health to ownership and was broader than the round-13 guarantee it served.
The guarantee was that a symlinked manifest cannot act as payload authority: verification must
fail, `doctor` must not report `valid`, and health must fail. All three still hold.

I reproduced the conflict from archived commits. With an intact marker and a symlinked manifest,
`f61e66b` returned `verifyPayload.ok: false`, `status: malformed`, `doctor.ok: false`, and
`managed: false`, while an unforced uninstall returned `ok: true` and removed the tree. At
`065985e` the three health answers are unchanged, while `managed: true` now agrees with that
same successful mutation. The old assertion therefore preserved the exact reporting
contradiction from my round-3 finding; it was not necessary to preserve the manifest trust
boundary.

I also attacked the marker predicate at `065985e`. A wrong package name, wrong skill identity,
invalid/unreadable marker, and symlinked marker each produced `managed: false`; unforced
uninstall refused each tree and left it in place. A valid marker plus a damaged payload produced
`managed: true`, `status: malformed`, and a successful unforced uninstall. I found no path that
reports `managed: true` for a tree the mutation predicate does not own.

## Other findings

### [MINOR] The new exit cleanup still leaves two bidding fixture roots per full run

The cycle-3 claim of zero relevant temp directories after a full run is not reproducible. I ran
`npm test` from a clean archive of `065985e` with an initially empty isolated `TMPDIR`; all
385 Node tests, 54 Python tests, and six manifest checks passed, but two
`parley-bidding-test-*` roots remained.

Both reproduce independently with focused green tests:

- `a frozen owned destination completes the install and names the debris` leaves a moved
  `.parley-deck.*.bak` tree whose directories remain mode `0555`.
- `one unreadable subdirectory deep in a destination no longer blocks anything` leaves a moved
  `.parley-bidding.*.bak/references` directory at mode `000`.

In each test the installer renames the hardened old destination to a backup. The test's
`finally` block restores permissions through the original destination path, which now names the
new tree, not the moved backup. The new process-exit handler then fails to remove the backup and
silently catches that failure, so the suite stays green while the tracked root leaks.

Restore permissions on the moved backup in each test, or make tracked cleanup normalize
permissions without following symlinks before removal, and add an isolated-TMPDIR assertion
that no `parley-bidding-test-*` roots remain. The separate `node -e` child cleanup is correct: a
focused runtime-probe run left zero entries in its isolated `TMPDIR`.

No other findings. The full suite and manifest checks passed, and the ownership/health behavior
changed only as cycle 3 states.

## Verdict

FINDINGS — listed above
