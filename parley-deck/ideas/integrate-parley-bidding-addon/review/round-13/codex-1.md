---
idea: integrate-parley-bidding-addon
review-round: 13
agent: codex-1
date: 2026-07-31
reviewed-commit: dd8d756
---

## Verdict

BLOCK

Cycle 16 closes all four regressions it claims to close, including every finding from my
round-12 review. Full-scope review found two new MAJOR boundary failures and one MINOR dry-run
fidelity failure. Stored data can still redirect manifest reads outside the destination tree,
and an install transaction can still corrupt another successful transaction's result.

## Outstanding findings — closed or not

- **Round-12 CRITICAL, unknown marker name deletes an unrelated sibling: CLOSED.** A marker name
  is now accepted only when the package ships that add-on or the direct-child destination has
  this installer's marker claiming the same skill. The destructive unknown-sibling regression
  passes, as does the guard that keeps a dropped add-on uninstallable from its own marker.
- **Round-12 MAJOR, invalid `addons` container types become core-only: CLOSED.** Only an absent
  field and explicit `false` are core-only. Strings, `true`, `null`, objects, and numbers now
  produce an unusable selection; the absent/`false` guards still pass.
- **Round-12 MAJOR, manifest keys escape the payload: CLOSED.** Keys are checked as canonical
  relative payload paths and each resolved entry must be a strict descendant before filesystem
  access. The external-sentinel regression passes.
- **Round-12 kimi-1 MAJOR, install leaves a partial fleet after commit failure: CLOSED for an
  isolated invocation.** All units stage before commit, and a later commit failure reverses
  earlier commits. The immovable-destination regression writes zero units and leaves no staging
  directories. The new first finding is a separate isolation/collision failure in overlapping
  transactions, not a failure of that single-process regression.
- **Round-12 kimi-1 NIT, discrimination count: CLOSED.** The implementation record now separates
  four fix-proving failures from two over-correction guards.
- **Recorded follow-ups: still open and non-blocking on their previously agreed terms.** These
  are B3.11 (only `parley-bidding` ships a manifest), the `dirExists` discovery-guard NIT,
  quarantine debris not discovered by `doctor`, and residual-disposal failures such as `uappnd`
  or delete-denying ACLs producing visible debris rather than a partial fleet.

## New findings

### [MAJOR] Install has neither a unique physical destination set nor transaction isolation

`installCommand` builds the fleet from logical target paths and immediately hands it to
`installFleetAtomically` (`lib/installer.js:656-658`). That function neither canonicalizes and
deduplicates physical destinations nor locks them (`lib/installer.js:1458-1552`). Its undo record
contains only ordinary destination/temp/backup paths, so rollback assumes no other writer has
changed them since this process committed (`lib/installer.js:1675-1759`). That assumption is
false both within one plan and across processes.

I reproduced both arms in isolated homes:

- I symlinked Agy's and Gemini's configured skill-container directories to the same physical
  directory, then ran `install --target all --include-undetected`. It returned `ok: true` for
  both targets even though their core units resolved to one physical destination. Gemini's later
  commit overwrote Agy's specialized core. Gemini's core was valid; Agy's core was malformed,
  missing `skills/SKILL.md`, and carried a Gemini marker. Symlinked configuration parents are
  supported by the resolver, so this is not a destination-payload symlink that the copy guard
  rejects.
- I synchronized two real installer processes after staging and interleaved their commits. A
  committed core, B committed the same core and bidding skill, then A's bidding commit failed.
  A's rollback moved B's new core aside and restored the stale pre-run core. B nevertheless
  resumed and returned `ok: true` for all six units. The restored core contained my old sentinel,
  `doctor` called it `valid` because core still has no manifest, and a hidden bidding backup was
  left behind. A reported failure, but B's successful result no longer described the fleet it
  had installed.

The rename transaction is atomic only while it has exclusive ownership of a collision-free
destination set. Before release, resolve each unit to its physical destination, reject aliases,
and acquire deterministic locks for every affected skills root/destination before preflight;
hold them through commit, rollback, and cleanup. Add regressions for aliased target roots and for
a second process committing between another process's commit and rollback.

### [MAJOR] A symlinked manifest is still trusted as payload authority

The payload walker deliberately ignores `parley-addon.json` before calling `lstat`
(`lib/addon-manifest.js:27-45`). `hasManifest` uses `stat`, which follows a link, and both
`manifestFileHash` and `readManifest` read the resolved target (`lib/addon-manifest.js:93-124`).
`verifyPayload` validates declared payload entries but never validates the manifest file itself
(`lib/addon-manifest.js:195-258`). Thus the cycle-16 key confinement does not confine the file
that supplies those keys, hashes, runtime policy, and ultimately the health verdict.

After a normal install I moved the installed bidding manifest to an external file and replaced
it with a symlink to that byte-identical file. `verifyPayload` returned `ok: true`; `doctor`
reported the skill `valid` and `managed`, with no problems. The install marker's manifest hash
also matched because hashing followed the link. Health now depends on mutable bytes outside the
installed destination even though add-on payload symlinks are otherwise forbidden.

Treat the manifest as part of the trust boundary: inspect it with `lstat` and require a regular,
non-symlink file before reading or hashing it. `hasManifest`, `readManifest`,
`manifestFileHash`, source validation, and installed health must agree on that rule. Regress both
direct verification and `doctor` with an external manifest sentinel.

### [MINOR] Install dry-run bypasses the fail-closed preflight used by a real install

Both fleet preflight and transaction planning are guarded by `!dryRun`; dry-run falls back to
`installTarget` (`lib/installer.js:622-679`). That fallback also skips `preflightSkillUnit` when
dry (`lib/installer.js:1352-1382`). Consequently it does not surface marker-derived selection
problems or other predictable preflight blockers, even though no write is needed to discover
them.

I installed a valid core-only Codex target, changed its marker's `addons` value to an object, and
ran install dry-run. It returned `ok: true` and said the core would be replaced. The identical
real install returned `ok: false` and blocked on the unusable add-on selection. The real command
is safe, but dry-run cannot be relied on to predict whether that command will run.

Run the same read-only fleet preflight for dry-run and report the same blocked/skipped plan;
only staging and commit should be omitted. Add a parity regression using the malformed marker
introduced by cycle 16.

## Release judgement

Not releasable as 2.1.0. The one release gate is to finish the filesystem trust boundary: every
input file must be regular and confined, and every mutation destination must be physically
unique and exclusively owned for the transaction's lifetime. That boundary must also be the one
reported by dry-run. Closing the three findings above satisfies that gate.

## What I verified

- Reviewed the complete implementation at `dd8d756`, the cycle-16 diff
  `5100f34..dd8d756`, the full `49fc3ec..dd8d756` round-7 delta, the ratified design, cycles
  10-16 of `IMPLEMENTATION.md`, the round-9 void record, and all round-12 reviews. I independently
  mapped marker, manifest, runtime, install, rollback, and health paths from source and a
  temporary code graph.
- Ran the six-item cycle-16 discrimination against both revisions from temporary copies. At
  `5100f34`, all four proof regressions failed and both guards passed. At `dd8d756`, all six
  passed. The add-on payload itself is unchanged from `714712f`.
- Ran the complete suite from a clean temporary copy: 344 Node tests passed with zero failures;
  the Python leg passed 54/54 on Python 3.14.6. The Node suite also passed 344/344 with system
  Python 3.9.6 first on `PATH`; the Python leg refuses that version by design.
- Verified the shipped bidding manifest: 47 files and aggregate
  `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`.
  The adapter validator checked four adapters with no errors, all 16 JSON files parsed, and no
  symlink or Python cache artifact was present in the shipped add-on.
- Ran `npm pack --dry-run --json`: 202 package entries, including the manifest plus all 47
  bidding payload files, with no cache artifacts. Both requested diff ranges pass
  `git diff --check`.
- Attacked install staging, late commit failure, reverse rollback, backup disposal, name
  collisions, logical-path aliases, dry-run, destination symlinks, and synchronized concurrent
  runs. Ordinary staging and single-process rollback remained clean; the three cases reported
  above did not.
- Every mutation probe used a temporary repository copy or isolated home. The reviewed repository
  remained at `dd8d756` with a clean working tree until this completed review file was written.
