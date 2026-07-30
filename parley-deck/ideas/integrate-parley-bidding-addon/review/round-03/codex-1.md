---
idea: integrate-parley-bidding-addon
review-round: 3
agent: codex-1
date: 2026-07-30
reviewed-commit: 5c324ef
---

## Verdict

BLOCK

## Round-2 findings — closed or not

1. **F4 — the Python probe ignores the runner's effective environment — CLOSED.**
   `doctor` and `status` pass `context.env` into the probe; the child is spawned with that
   environment; results are cached by its `PATH`; and the call has a 5-second timeout. I
   reproduced the original parent-has-Python/caller-has-empty-`PATH` case and it now returns
   `runtime.ok:false` and exit 1. A deliberately stalled `python3` was terminated after
   5.008 seconds and reported unavailable.
2. **F5 — the add-on-list test depends on the host Python floor — CLOSED.** The test now
   asserts the skill list and each payload status rather than the availability-inclusive
   top-level `result.ok`. The full 290-test Node suite passes with `/usr/bin/python3` 3.9.6
   first on `PATH`, and it also passes with Python 3.14.6.
3. **F6 — `paths` probes while text `status` hides availability — CLOSED.** Runtime probing
   is now an explicit traversal option. The new sentinel test proves `paths` does not launch
   `python3`, while `doctor` does; `status` text includes the indented `unavailable:` detail.
4. **F7 — duplicated “installed” wording — CLOSED.** The stderr sentence now reads
   “One or more installs are operationally unavailable.”

## Ruling on the skills-CLI question

**(b).** A byte-valid foreign-installer tree is not malformed. Payload integrity and
installer ownership are separate facts, just as payload validity and runtime availability
are now separate facts.

For an absent marker only, when the packaged source add-on carries `parley-addon.json` and
the installed tree's manifest fully verifies, report the unit as valid-but-unmanaged (for
example, `status:"valid"` plus `managed:false` or an equivalent explicit field/message).
Do not synthesize a marker and do not let `install`, update, or `uninstall` treat the
directory as owned. The missing marker still means provenance is not anchored to an install
event, so the output must not imply the stronger managed-install guarantee.

An unreadable marker remains malformed. An unmarked installed tree with neither the manifest
that the packaged source ships nor a marker also remains malformed; that preserves the
gutted-tree detection ratified in round 3. This is a required fix before acceptance because
the README recommends the universal installer first, and the current output assigns the
wrong integrity state to its faithfully copied payload.

## New findings

### [MAJOR] The recommended universal-installer path is mislabeled `malformed`

**Where:** `README.md:140-161`; `CHANGELOG.md:56-59`; `lib/installer.js:1214-1255`;
`lib/installer.js:1344-1353`

**What:** `skillUnitStatus` unconditionally adds an integrity problem when an expected
directory has no Parley Deck install marker. For an unmarked add-on, `manifestProblems`
performs no manifest validation at all. A tree copied faithfully by the README's first
installer therefore becomes `malformed` solely because a different installer owns it, even
when its complete manifest verifies.

**Why it matters:** `malformed` says the payload is structurally or byte-wise defective. In
this supported installation path the known fact is instead “payload internally valid,
management provenance external.” Conflating those states makes the health output factually
wrong and prevents the manifest from providing useful cross-installer health evidence. This
was disclosed before round 2; this round's requested ruling promotes it to a blocking
compatibility finding.

**Evidence:** I copied the repository's `skills/parley-bidding` directory unchanged into an
isolated Claude skills directory, without a marker. `addonManifest.verifyPayload(dest)`
returned `ok:true`, no problems, 47 files, aggregate
`sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`.
`doctor` nevertheless returned exit 1 and
`parley-bidding:{status:"malformed", problems:["no parley-deck-skill install marker: ..."]}`.
The README recommends `npx -y skills add feci/parley-deck-skill` before presenting the
package-native installer and its health checks.

**Fix:** Implement ruling (b) generically. If the marker is absent, the packaged source
declares a manifest, and `verifyPayload` succeeds, return a distinct valid-but-unmanaged
state and continue any non-mutating runtime probe. Keep ownership-sensitive mutations
fail-closed. Add regression cases for (1) a complete externally copied manifested tree,
(2) the same tree with a payload mismatch, and (3) an unmarked tree gutted to `SKILL.md`.
Update the changelog's unconditional-malformed statement.

### [MINOR] Malformed high-major interpreter output passes the Python floor

**Where:** `lib/installer.js:1307-1313`

**What:** The probe considers the interpreter found when only `major` is an integer. A
non-numeric `minor` is coerced to zero. Output such as `4.not-a-version` therefore becomes
`{found:true, major:4, minor:0}` and satisfies the declared `>=3.10` floor.

**Why it matters:** A PATH-resolved shim or broken executable that is not a working Python
interpreter can make `doctor` return a false healthy result. This is an edge case, but it is
the fail-open direction for the operational-availability check.

**Evidence:** In an isolated install I supplied an executable named `python3` that exits 0
after printing `4.not-a-version`. At `5c324ef`, `doctor` returned
`{"doctorOk":true,"runtime":{"ok":true,"requirement":">=3.10","detail":"python3 4.0"}}`.

**Fix:** Parse trimmed stdout with an anchored two-integer pattern such as
`^(\d+)\.(\d+)$`; reject the probe unless both captures are valid integers and no extra
stdout remains. Add a negative regression test using a malformed version whose numeric major
is above the required floor.

## What I verified and found correct

- The reviewed repository was clean at exactly
  `5c324ef4905c7d5199bfa22a22ecc10fe3c31a0e`; the fix-up diff changes only
  `CHANGELOG.md`, `lib/installer.js`, `test/bidding-addon.test.js`, and
  `test/installer.test.js`. `git diff --check 89069b0..5c324ef` is clean.
- `PYTHONDONTWRITEBYTECODE=1 npm test` passes: **290/290 Node tests**, **54/54 Python
  tests** across seven files, and the 47-file add-on manifest check.
- With `/usr/bin/python3` 3.9.6 first on `PATH`, the complete Node suite is independently
  **290/290**. With Homebrew Python 3.14.6, the complete npm suite is green.
- The focused bidding and installer suites pass **82/82**.
- The F4 regression now returns payload `valid`, runtime unavailable, top-level
  `doctor.ok:false`, and exit 1 when the caller's `PATH` is empty even though the parent
  process has Python.
- Two `run()` calls in one process with different `PATH` values receive different probe
  results. A stalled interpreter is bounded by the documented 5-second timeout.
- `paths` does not launch the interpreter and returns `runtime:null`; `doctor` and `status`
  do probe. Text `status` exposes the unavailable detail.
- No `__pycache__`, `.pyc`, or `.pyo` appeared under `skills/parley-bidding/`, and no file
  under that directory was modified.

## Open questions for the implementer

None. The required foreign-installer behavior is specified by ruling (b), and the malformed
version-output case has a deterministic fix.
