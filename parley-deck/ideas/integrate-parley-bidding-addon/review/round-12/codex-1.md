---
idea: integrate-parley-bidding-addon
review-round: 12
agent: codex-1
date: 2026-07-31
reviewed-commit: 5100f34
---

## Verdict

BLOCK

Cycle 15 fixes the separator-based escape and the quarantine transaction survived the
adversarial cases I tried. The release is still blocked by one CRITICAL continuation of the
round-11 marker finding and two MAJOR stored-data validation failures.

## Outstanding findings — closed or not

- **Round-11 CRITICAL, marker-controlled removal: NOT CLOSED.** The out-of-root traversal arm is
  closed: invalid entries fail closed and the second guard confines destinations to direct
  children of the skills directory. However, the round-11 requirement also said *unknown*
  entries must fail closed. A syntactically valid unknown name still lets mutable marker data
  select an unrelated sibling skill for forced deletion. The marker container itself is also
  not schema-validated. See the first two findings below.
- **Round-11 MAJOR, inaccurate removability prediction: CLOSED.** `firstRemovalObstacle` is gone.
  Phase A now quarantines the entire plan before Phase B deletes anything. Rename collision,
  Phase-A failure, rollback failure, a symlinked destination, dry-run, and concurrent
  uninstall probes all failed safely or produced an honest, recoverable result.
- **Round-11 MAJOR, over-broad legacy exemption: CLOSED.** The downgrade is now limited to a
  packaged source that ships no manifest, and the regression uses the genuinely legacy
  `parley-worktrees` case.
- **Round-11 MINOR, raw EACCES from `hashFile`: CLOSED.** Manifest hashing is inside the
  verifier's error boundary and the regression reports `unreadable (EACCES)`.
- **Recorded follow-up, quarantine observability: still open and non-blocking.** A process or
  machine crash during Phase A introduces a recoverable state in which a prefix of the fleet is
  intact under `.removing` names while its normal destinations are absent. `doctor` reports the
  destinations missing but does not discover the quarantines. This is the crash-stop variant of
  the already-recorded Phase-B-debris visibility limit, not a reason to restore a removability
  predicate or a separate 2.1.0 blocker.
- The previously deferred B3.11 follow-up (only `parley-bidding` currently ships a manifest) and
  the round-9 `dirExists` discovery-guard NIT are unchanged.

## New findings

### [CRITICAL] A plain unknown marker name can still delete an unrelated sibling skill

`markerAddonNames` now rejects dangerous syntax, but it never establishes that an accepted name
is a discovered Parley add-on or an existing installer-owned legacy add-on
(`lib/installer.js:993-1026`). `targetSkillUnits` consequently constructs a unit whose
`addon` is `null` for an unknown name (`lib/installer.js:1077-1107`). Finally,
`removeFleetAtomically` deliberately waives the destination-ownership check under `--force`
(`lib/installer.js:1404-1417`). Confinement therefore limits the blast radius to the skills
directory, but does not limit it to this package's skills.

I installed the Codex target, changed the core marker to
`"addons": ["unrelated-sentinel"]`, and created
`.codex/skills/unrelated-sentinel/KEEP`. `uninstall --target codex --force --json` returned
`ok: true`, reported both `parley-deck` and `unrelated-sentinel` as `removed`, and deleted the
sentinel. The real Parley add-on directories remained because the corrupt marker had excluded
them. This is still marker-controlled deletion of unrelated user data; only its location has
been narrowed since round 11.

`--force` may waive ownership for a destination the caller explicitly selected. It must not
waive authority for an unknown destination selected only by mutable stored data. Before release,
a recorded name should be accepted only if it is in the packaged add-on registry, or—so an add-on
removed from a newer package can still be uninstalled—if the existing direct-child destination
has a readable installer marker whose skill identity matches that name. Unknown and reserved
names must make the core malformed and block both mutation plans. Add an in-skills unrelated
sentinel regression for forced uninstall.

### [MAJOR] Invalid `addons` container types silently become a healthy core-only selection

At `lib/installer.js:1013-1015`, every non-array value is treated as `names: []` with no problem.
That is correct for an absent legacy field and the intentional current value `false`, but it also
accepts strings, `true`, `null`, objects, and numbers. The promised “unusable selection” branch
only validates entries after `Array.isArray` has already passed.

Starting from a valid `--no-addons` install, I replaced `false` with the string
`"parley-bidding"`. `doctor --json` returned `ok: true` and a `valid` core. Default install
returned success and repaired the marker instead of blocking; forced uninstall returned success
and removed the core. I reproduced the same core-only interpretation with `true`, `null`, `{}`,
and `42`. Thus health is falsely green for malformed current-schema state and neither mutation
preflight implements the cycle-15 fail-closed claim.

Distinguish an absent legacy property and the explicit `false` value from every other non-array
value. Invalid container types must set `markerProblem`, just like invalid array entries, with
doctor/install/uninstall regressions.

### [MAJOR] Manifest file keys can escape the add-on payload root

`readManifest` validates each hash value but does not validate its key as a canonical relative
payload path (`lib/addon-manifest.js:120-161`). `verifyPayload` then feeds each stored key to
`path.join(root, rel)` without a descendant check (`lib/addon-manifest.js:168-189`). The aggregate
digest only proves that the manifest agrees with its own key/value map; it does not add path
confinement.

In a temporary copy, I added `"../outside-sentinel"` with the external file's correct digest,
recomputed the aggregate, and `verifyPayload` returned `ok: true`. I then tested the real install
path with a manifest entry for `../parley-deck/SKILL.md`. Source verification passed, install
returned `ok: true`, and the installed bidding skill remained `valid` in doctor because the
relative path resolved to the installed sibling core. The manifest therefore certifies a file
that is not in the add-on payload, and both source preflight and installed health can be made
dependent on arbitrary sibling content.

Reject absolute paths, separators native to the wrong manifest syntax, empty/`.`/`..` segments,
and non-canonical aliases when reading the manifest. Resolve every accepted key and require it to
be a strict descendant of the payload root before any filesystem access. Regress both direct
verification and install source preflight with an external sentinel.

## Release judgement

Not releasable as 2.1.0. The single release gate is to complete the stored-data path boundary:
no marker selection or manifest key may become a filesystem target until it is schema-valid,
confined, and authorized for this package. That requires closing all three findings above and
adding their destructive/false-green regressions.

## What I verified

- Reviewed the complete implementation at `5100f34`, the cycle-15 diff
  `12f9071..5100f34`, the full `49fc3ec..5100f34` round-7 delta, the ratified design, cycles
  10–15 of `IMPLEMENTATION.md`, round 9's void record, and all round-11 reviews.
- Mapped the implementation call graph independently. The marker-to-path bridge is
  `markerAddonNames` → `expectedAddonNames`/`targetSkillUnits`; install uses the separate staged
  `copyPayloadAtomically` commit path.
- Ran the complete suite from a clean archive of the reviewed commit: 338 Node tests passed with
  zero failures; the Python leg passed 54/54 on Python 3.14.6. The Node suite also passed 338/338
  with `/usr/bin/python3` 3.9.6 first on `PATH`; the Python leg refuses that version by design.
- Verified the shipped bidding manifest: 47 files and aggregate
  `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`.
  The adapter validator checked four adapters with no errors, and all JSON/schema parsing checks
  passed.
- Ran `npm pack --dry-run --json`: 202 package entries, including 48 bidding entries, with no
  cache files. Built the current macOS arm64 portable executable, installed all six skills into
  an isolated home, and obtained an all-valid `doctor` result.
- Attacked the quarantine transaction directly. A pre-existing quarantine-name collision blocked
  the fleet with every destination intact. A later Phase-A rename failure rolled earlier renames
  back; an injected rollback failure produced an error naming the exact intact quarantine rather
  than deleting it. Dry-run changed nothing. A destination symlink was unlinked without touching
  its external target. Two synchronized uninstalls yielded one complete success and one blocked
  run, with no residue. Install's `.bak` namespace is distinct from uninstall's `.removing`
  namespace, and the interleavings I inspected do not give either command a path outside its
  resolved unit destinations.
- Reproduced all three findings against temporary copies/homes at exactly `5100f34`. No probe
  mutated the reviewed repository, and the source working tree was clean immediately before this
  review file was written.
