---
agent: kimi-1
idea: addon-manifest-coverage
review-round: 4
date: 2026-08-02
reviewed-commit: 065985e
---

## Ruling on the managed / round-13 conflict

**The resolution stands. No ratified guarantee was overturned. Not a block.**

What round 13 actually ratified (codex-1 MAJOR, `integrate-parley-bidding-addon/review/round-13/codex-1.md:76-94`):
the manifest file must be inside the trust boundary — `lstat`, regular non-symlink file, never read or hashed
through a link — and a tree whose manifest is external must not read `valid`. The words "and `managed`" appeared
in that finding's *symptom description* ("`doctor` reported the skill `valid` and `managed`, with no problems"),
not in its normative demand, which was confined to the manifest trust boundary and health.

Reproduced at `065985e` (fresh install, bidding manifest renamed out and symlinked back):

- `status` is `malformed` — the round-13 guarantee holds.
- `doctor.ok` is `false` — health still fails, now asserted a second time in the test.
- `managed` is `true` — the inverted half.
- **An unforced `uninstall` removes that same tree.** This is the decisive fact: the round-13-era assertion
  `managed !== true` was written while `managed` was a health-flavored aggregate. The review process itself has
  since redefined it — cycle 2 (hermes-1 R2-2(b), co-signed by me, codex-1 agreeing) made ownership marker
  evidence, and codex-1's round-3 MINOR completed it: `managed` must answer what the mutation paths ask, or
  `doctor` contradicts `uninstall` about the same directory.

Reproduced the contradiction at `f61e66b`: the same symlinked-manifest tree reported `managed: false` while an
unforced `uninstall` removed it without protest. Both changed assertions fail at `f61e66b` and pass at `065985e`,
as the notes claim.

The only way to keep both halves would be to change `installerOwnsDestination` to disown defective trees — which
would strand every damaged install (install and uninstall both refusing without `--force`) and which round 13
never asked for. The inversion is flagged in the test comment and put to this round explicitly, not settled
quietly. Ownership and health are separate questions; round 13 settled health, and health is unchanged.

## Other findings

**Attack on the marker predicate — no path to `managed: true` for a tree this installer does not own.**
Four attacks on a freshly installed core, each checking `doctor` *and* the mutation path for agreement:

- wrong `name` in marker → `managed: false`, unforced uninstall refuses
- wrong `skill` in marker → `managed: false`, unforced uninstall refuses
- malformed marker JSON → `managed: false`, unforced uninstall refuses
- marker replaced by a symlink to the byte-identical marker → `managed: false` (`lstat` → not a regular file →
  present-but-unreadable), unforced uninstall refuses

Reporting and behavior agree in every case, in both directions. codex-1's round-3 scenario also reproduced at
`065985e`: delete installed `plugin.json`, valid marker remains → `status: malformed`, `missing: ["plugin.json"]`,
`managed: true`, unforced uninstall removes the tree.

**Temp-directory cleanup — the exit handler is sound, but the "0 before, 0 after" measurement does not reproduce.
[MINOR]** Measured with an isolated `TMPDIR`, one full `node --test` run at `065985e` (385/385) leaves **18
directories** behind:

1. **2 from `test/bidding-addon.test.js` — one of the three files this cycle fixed.** Reproduced individually:
   - `a frozen owned destination completes the install and names the debris` (test/bidding-addon.test.js:1466):
     `freeze()` makes the old tree 0555/0444; the install renames it to `.parley-deck.<pid>.<ts>.bak`; `thaw()`
     chmods the *original* path, which by then is the new tree, so the debris stays frozen; the exit handler's
     `rmSync` gets EACCES and the `catch` swallows it.
   - `one unreadable subdirectory deep in a destination no longer blocks anything`
     (test/bidding-addon.test.js:1490): same shape with a 0000 `references` inside the `.bak`.
   These leftovers **resist a plain `rm -rf`** (verified — removal needs a chmod pass first), which is exactly
   the failure mode that filled the disk twice. Two directories per run instead of hundreds, but they accumulate
   forever and the standard cleanup command cannot take them.
2. **16 from `skills/parley-tracker/bin/claim.test.js` and `validate.test.js`** — both create temp dirs via
   `mkdtempSync` with no cleanup at all, so "all three test files that create them" undercounts: at least five
   test files create them.

The exit handler itself is clean: it runs after the runner fixes the exit code, cannot mask a failure, and no
test observes temp dirs after exit, so no test's meaning changes. The `node -e` child (`b6-` prefix) cleans its
own directory — verified, zero `b6-` leftovers — and on a red child the leak lands on an already-failing test,
which is acceptable. Suggested direction, not prescription: best-effort `chmod u+rwx` before the exit-handler
`rmSync` (or thaw the renamed debris by tracking it), extend `trackTemp` to the two tracker files, and re-run
the measurement with an isolated `TMPDIR` before claiming 0/0.

**Everything else in the cycle-3 notes reproduces.** `node --test` 385/385 at `065985e`, and 385/385 again
under a PATH whose only python3 is `/usr/bin/python3` 3.9.6; python leg 54/54; `--check` green on all six
manifests; `skills/parley-bidding` aggregate still `sha256:7854adf1…`. (The shared `/var/folders` count was not
usable as a before/after measurement — hermes-1's round-4 suite is running concurrently on this machine and its
in-flight fixtures share the namespace; the isolated-`TMPDIR` measurement above is the controlled one.)

## Verdict

FINDINGS — listed above (one MINOR: temp-cleanup incomplete in `bidding-addon.test.js`, two tracker files
uncovered, "0 before, 0 after" not reproducible). The `managed` resolution is correct and the round-13 health
guarantee is intact; nothing here blocks on the question this round was called to rule on.
