---
idea: integrate-parley-bidding-addon
review-round: 4
agent: codex-1
date: 2026-07-30
reviewed-commit: b180127
---

## Verdict

BLOCK

## Round-3 findings — closed or not

1. **The faithfully copied, markerless manifested tree was mislabeled `malformed` —
   CLOSED for that intended case.** In an isolated archive of `b180127`, I installed a
   managed core, copied `skills/parley-bidding` byte-for-byte without its marker, and put a
   sentinel `python3` shim on the caller's `PATH`. `doctor` returned exit 0 and the unit
   returned `status:"valid-unmanaged"`, `managed:false`, `marker:null`, no problems, and
   `runtime:{ok:true,detail:"python3 3.12"}`; the sentinel proved the runtime was still
   probed. A normal install without `--force` blocked before changing the markerless tree,
   while `install --force` replaced it and wrote a managed marker.
2. **Malformed high-major interpreter output passed the floor — CLOSED.** A `python3` shim
   that printed `4.not-a-version` now produced `runtime.ok:false`, the “python3 is not
   available” detail, and `doctor` exit 1. The parse at `lib/installer.js:1361` is anchored
   and requires both numeric components.

Those closures do not close the review: the implementation admits states beyond the faithful
copy that the ruling was intended to recognize.

## New findings

### [MAJOR] `valid-unmanaged` trusts the installed manifest as its own authority
**Where:** `lib/installer.js:1285-1300`, `lib/installer.js:1308-1311`;
`lib/addon-manifest.js:65-90`, `lib/addon-manifest.js:116-210`

**What:** `unmanagedButVerified` uses the packaged source manifest only as a boolean
capability check. It then verifies the installed payload solely against the installed
manifest, without comparing that manifest, its aggregate, or its raw hash with the packaged
source. The installed manifest can therefore redefine both the accepted payload and the
runtime requirement while still earning `valid-unmanaged`.

**Why it matters:** This is not a request for tamper resistance, which the design explicitly
disclaims. It is an ordinary package-identity and health false green. The new verdict is meant
to recognize the packaged tree copied by another installer, but the current predicate also
recognizes an arbitrary self-consistent tree. Worse, `runtime` is outside the payload
aggregate, so a single manifest-field edit can bypass B6 without rehashing any payload file.

**Evidence:** I copied the packaged tree without a marker, deleted only
`runtime:{"python":">=3.10"}` from its installed `parley-addon.json`, and left every file hash
and the aggregate untouched. `verifyPayload(dest)` returned `ok:true`; the installed manifest
hash no longer matched the packaged source; with an empty `PATH`, `doctor` nevertheless
returned `ok:true`, `status:"valid-unmanaged"`, `runtime:null`, and no problems. In a second
probe I changed `SKILL.md`, regenerated a manifest from that replacement tree, omitted the
runtime field, and obtained the same healthy verdict.

**Fix:** Anchor the unmanaged proof to the packaged source: require the source manifest
itself to verify, require the installed payload to verify, and require the installed manifest
bytes (or an equivalently complete canonical digest covering file inventory, hashes, aggregate,
and runtime metadata) to match the packaged source manifest. If older/different unmanaged
versions must be recognized, give them a non-healthy or explicitly version-unknown state
rather than calling them the currently packaged payload. Add negative tests for runtime-field
removal and a self-consistent replacement tree with no marker.

### [MAJOR] Lifecycle mutations use marker-path existence instead of actual ownership
**Where:** `lib/installer.js:906-920`, `lib/installer.js:987-1001`,
`lib/installer.js:1032-1074`

**What:** Install preflight and `installSkillUnit` consider any existing marker pathname
sufficient authorization to replace a destination. They do not require a readable marker
whose `name` is `parley-deck-skill` and whose `skill` matches the selected unit. This
contradicts the new health classification, which correctly says a foreign or unreadable
marker is malformed and unmanaged. Uninstall uses a stricter ownership predicate, but it
does not preflight the selected set before deleting earlier managed units.

**Why it matters:** A foreign manager's tree can be destroyed without the explicit
`--force` override. The inverse path can partially mutate before reporting that it refused
the unmanaged unit: the ownership boundary is neither consistent nor fail-before-write.

**Evidence:** For both (a) a readable marker with `name:"other-installer"` and (b) invalid
marker JSON, `doctor` first reported `malformed` with the expected ownership/readability
problem. I added a `FOREIGN-SENTINEL` file and ran a normal
`install --only parley-bidding` with no `--force`; it returned `ok:true`, reported the add-on
`replaced`, deleted the sentinel, and wrote a new `parley-deck-skill` marker. Conversely,
with a managed core and markerless `valid-unmanaged` bidding tree,
`uninstall --only parley-bidding` returned false and preserved bidding, but had already
removed the core.

**Fix:** Centralize one parsed ownership predicate (marker present, regular and readable;
correct installer name; correct skill identity) and use it in install preflight, install,
and uninstall. Without `--force`, preflight every selected unit before the first mutation.
Only an explicit `--force` may replace or remove a destination whose ownership is absent,
unreadable, foreign, or mismatched. Add foreign-marker and unreadable-marker install tests
with a sentinel, plus a no-partial-uninstall test.

### [MAJOR] `--no-addons` and excluding `--only` selections hide an existing unmanaged add-on
**Where:** `lib/installer.js:827-880`, `lib/installer.js:883-903`;
`README.md:20-37`, `README.md:260-263`

**What:** A selective install writes only the selected units and records that selection in
the core marker, but it neither reconciles nor reports known add-on directories already
present beside the core. Subsequent unflagged `doctor` and `status` derive their traversal
only from the core marker, so the leftover directory disappears from health output while
remaining on disk and available to the runtime.

**Why it matters:** This is the precise README-first transition the new verdict makes
legitimate. The documented opt-out says to use `--no-addons`, or an excluding `--only`, to
leave bidding out. After a universal install, those commands can instead leave the
markerless bidding skill installed and make the native health command stop showing it.
`doctor` returning green is then not evidence that the availability opt-out took effect.

**Evidence:** I simulated the README-first full copy: before the native command, bidding was
`valid-unmanaged`. After successful `install --force --no-addons`, the bidding directory
still existed, still had no marker, was absent from the default doctor's unit list, and
`doctor.ok` was true. Repeating with `install --force --only parley-design` produced the
same result: only core and design were reported while bidding remained installed.

**Fix:** Define transition semantics rather than testing only fresh homes. A safe minimum is
to detect known sibling add-on directories omitted from the recorded selection and keep them
visible as unselected-present/unmanaged, with a clear removal action; do not silently delete
foreign-owned trees. If `--force` plus a selection is intended as desired-state
reconciliation, document that destructive meaning explicitly and preflight all removals.
Add universal-full-copy to `--no-addons` and excluding-`--only` transition tests.

### [MINOR] A non-regular marker entry qualifies as “entirely absent”
**Where:** `lib/installer.js:1254-1265`, `lib/installer.js:1465-1492`

**What:** `readMarkerState` calls `fileExists`, which returns true only when `statSync`
reports a file. A directory or dangling symlink at the marker pathname is therefore returned
as `present:false`, not present-but-unreadable. Because the marker pathname is excluded from
manifest inventory, an otherwise verifying tree takes the `valid-unmanaged` branch.

**Why it matters:** This violates the round-3 precision that only an entirely absent marker
qualifies and every present-but-unreadable marker is malformed. JSON then says `marker:null`
even though a filesystem entry is present.

**Evidence:** I created a directory named `.parley-deck-skill-install.json` inside an intact
markerless bidding tree. `fs.existsSync(markerPath)` was true, but the unit returned
`status:"valid-unmanaged"`, `managed:false`, `marker:null`, and no integrity problems.

**Fix:** Use `lstatSync` for marker state. Only `ENOENT` is absent; a directory, symlink,
device, unreadable regular file, or invalid JSON is present-but-unreadable/malformed. Add
directory and dangling-symlink regression cases.

### [MINOR] The probe cache still omits variables that select the interpreter
**Where:** `lib/installer.js:1333-1346`

**What:** The new cache key enumerates five variables, but a stable `PATH` can resolve a
version-manager shim whose answer is selected by another variable such as
`PYENV_VERSION`. The second environment then reuses the first environment's result.

**Why it matters:** This remains a false-green route for library callers that invoke
`run(argv, io)` with multiple effective environments in one process. It is narrower than the
original PATH bug because the ordinary one-shot CLI uses one environment.

**Evidence:** With one fixed `PATH`, my `python3` shim printed 3.12 unless
`PYENV_VERSION=old`, when it printed 3.9. The first doctor call with `new` returned healthy
3.12. The second call in the same process with `old` should have failed the `>=3.10` floor,
but returned the cached 3.12 result and stayed healthy.

**Fix:** Key on a stable serialization of the complete effective environment plus the
effective working directory, or scope/remove the process-global cache. Add a same-PATH,
different-`PYENV_VERSION` regression in both call orders.

## What I verified and found correct

- The reviewed repository was clean at
  `b180127166d94b95122f6ce94a32cb0e32c7e35a`; `git diff --check
  5c324ef..b180127` was clean. The fix-up changes only `CHANGELOG.md`, `README.md`,
  `lib/installer.js`, and `test/bidding-addon.test.js`.
- From a `git archive` of the reviewed commit, with the checkout's dependency directory
  supplied read-only through `NODE_PATH`, `PYTHONDONTWRITEBYTECODE=1 npm test` passed:
  **296/296 Node tests**, **54/54 Python tests** on Python 3.14, and the 47-file manifest
  check at aggregate
  `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`.
- The intended `valid-unmanaged` case has the required JSON shape, does not itself fail
  health, and still fails health when its unchanged runtime probe is unavailable.
- An absent marker does not grant mutation ownership: install without force blocked and
  preserved the tree; install with force replaced it and became managed; uninstall without
  force preserved the unmanaged add-on.
- File-form unreadable markers and markers naming another installer are `malformed`.
  The defect is their separate treatment by install, not the health verdict.
- The legacy `markerSchema`-less carve-out is unchanged. I reproduced its known behavior:
  after removing the schema/manifest anchor, the manifest, and the entire scripts directory,
  it still returned `valid`, `managed:true`, exit 0. I do not reopen that explicitly accepted
  compatibility boundary.
- `status` explains core integrity failures and remains exit 0. A corrupt core marker printed
  its `integrity:` detail; the returned exit code was 0. The changelog correctly names
  `doctor` as the health gate.
- I agree with deferring manifests for the other five units. B3.11 deliberately preserves
  their manifest-free behavior, and the changelog now states the exact one-unmanaged/five-
  malformed result. This review does not reopen that ratified scope decision.
- No `__pycache__`, `.pyc`, or `.pyo` appeared in the reviewed repository, and the reviewed
  worktree remained clean.

## Open questions for the implementer

1. For an intact unmanaged tree from an older package version, should source-manifest
   mismatch be `malformed`, or should it receive a separate version-unknown/drift verdict?
   It must not receive the current package's unqualified `valid-unmanaged` health proof.
2. Are `--no-addons` and `--only` intended as desired-state reconciliation or only as
   copy-selection flags? The current README reads as desired-state for bidding availability,
   while the implementation is selection-only and makes leftovers invisible.
