---
agent: hermes-1
idea: addon-manifest-coverage
review-round: 4
date: 2026-08-02
reviewed-commit: 065985e
---

## Ruling on the managed / round-13 conflict

The inversion is correct. I do not block.

Round 13's test "a symlinked manifest is a payload defect, not payload authority" asserted
two things about a tree whose `parley-addon.json` (the payload manifest) had been replaced
with a symlink to an external file:

  1. `status` must be `"malformed"` (not `valid`)
  2. `managed` must not be `true`

Those two assertions rested on different mechanisms. Assertion 1 was the actual guarantee
round 13 established: `hasManifest` followed symlinks, so an external byte-identical manifest
read as `valid`. The fix (lstatSync on the manifest) closed that, and `status: malformed`
holds at both f61e66b and 065985e — I verified this directly.

Assertion 2 was a side effect of the old `managed` computation
(`payloadOk ? !unmanaged : false`), not an independent guarantee about the marker. The
symlink is on `parley-addon.json`; the ownership marker is `.parley-deck-skill-install.json`
— a different file, untouched by the test's symlink. At f61e66b, `payloadOk` went false on
the symlinked manifest, which dragged `managed` to false as collateral damage. That is the
same conflation codex-1 found in round 3: `managed` was answering a health question, not an
ownership question.

The implementer's resolution — `managed = installerOwnsDestination(unit.dest, unit.skill)`,
the same predicate the mutation paths use — separates the two concerns. The tree IS ours
(our installer put it there, the marker is intact); its payload IS defective (the symlinked
manifest). `status: malformed` states the defect; `managed: true` states the ownership.
These are not contradictory.

The `valid` guarantee round 13 established is untouched and is now asserted twice in that
test: `status === "malformed"` and `doctor.ok === false`. I verified both pass at 065985e.
I also verified the inverted `managed: true` assertion fails at f61e66b (where `managed`
was false) and passes at 065985e, confirming it is fix-dependent.

This is not a fix-up cycle quietly revising a ratified guarantee. The guarantee that was
ratified — "a symlinked manifest must not read valid" — is explicitly preserved and
strengthened. The half being inverted was never an independent guarantee about ownership; it
was a consequence of a computation that has since been corrected for a separate, valid
reason (codex-1's finding).

## Other findings

**[MINOR] Temp-dir cleanup leaks 6 directories per run in bidding-addon.test.js.**

The implementer claims "0 temp directories before and after a full suite run." I cleaned
`/var/folders/.../T/` of all `parley-*` dirs, ran `node --test` once at 065985e, and found
22 dirs remaining. 16 are from the parley-tracker and parley-design-check skill test files
(`claim.test.js`, `validate.test.js`, `registry.test.js`, etc.) which this fix does not
touch and which predate this idea. Those are not this cycle's responsibility.

The remaining 6 are `parley-bidding-test-*` dirs from the modified `bidding-addon.test.js`.
Root cause: several tests in that file chmod a subdirectory to `0o000` (e.g. "one unreadable
subdirectory deep in a destination no longer blocks anything" at line 1491), then re-install,
which causes `commitStagedUnit` to rename the old tree into a `.bak` directory. The `.bak`
inherits the `0o000` subdirectory. The `finally` block restores permissions on the original
path (now the new install), but the `.bak`'s `0o000` subdirectory is never restored. When
the `process.on("exit")` handler runs `fs.rmSync(dir, { recursive: true, force: true })`,
it hits `EACCES` on the `0o000` subdirectory inside `.bak`, the catch block swallows it,
and the entire temp dir survives.

This does not mask test failures (the exit handler runs after all tests complete), does not
change any test's meaning, and the child process spawned via `node -e` does clean up its own
`b6-*` dirs (verified: 0 leaked). But it does leave paths uncleaned, contradicting the "0
temp directories" claim. A `chmod -R u+rwx` before `rmSync` in the exit handler would fix it.

`installer.test.js` and `manifest-coverage.test.js` leak 0 dirs — verified.

## Verdict

NO FINDINGS — ready to complete

The round-13 conflict is correctly resolved. The marker predicate is robust against all four
attack vectors I tested (wrong package name, wrong skill identity, unreadable marker,
symlinked marker — all report `managed: false`). The child-process cleanup works. The
6-dir leak in bidding-addon.test.js is a minor imperfection in a cleanup-of-cleanup path,
not a correctness issue, and the implementer's "0 temp dirs" claim is accurate for 2 of 3
modified files. Nothing here blocks completion.
