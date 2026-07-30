---
idea: integrate-parley-bidding-addon
review-round: 1
agent: codex-1
date: 2026-07-30
reviewed-commit: 714712f
---

## Verdict

BLOCK

## Findings

No CRITICAL or NIT findings.

### [MAJOR] `doctor` does not enforce the add-on's declared Python minimum

**Where:** `lib/installer.js:1208-1221`; `scripts/run-python-tests.js:46-74`

**What:** Installed-unit health is derived only from payload validation. The Python probe and
the `>=3.10` check exist in the npm test runner, but `doctor` never performs an operational
runtime check.

**Why it matters:** FINAL.md B6 requires `doctor` to distinguish payload-valid from
operationally unavailable and to fail health when the declared interpreter minimum is absent.
Without Python, the installed payload can be byte-valid while its published commands are
unusable, yet `doctor` exits 0.

**Evidence:** With `python3` unavailable, the measured result is
`{"doctorOk":true,"status":"valid","problems":[]}`. The only runtime-floor enforcement is
`resolveInterpreter()` in `scripts/run-python-tests.js`; `skillUnitStatus()` in
`lib/installer.js` calls only `validateInstalledPayload()`.

**Fix:** During `doctor` and `status`, read the validated manifest's `runtime.python`
requirement, probe the declared interpreter, and report a separate operationally-unavailable
state/problem when it is missing or below the floor. `doctor` must exit non-zero. Add negative
tests for a missing interpreter and for an interpreter below 3.10.

### [MAJOR] Expected installer-managed units pass health checks without a readable marker

**Where:** `lib/installer.js:1212-1219`; `lib/installer.js:1253-1258`;
`lib/installer.js:1315-1320`

**What:** `readMarker()` collapses a missing marker and any read/JSON failure to `null`, and
`manifestProblems()` treats that result as having no problems. An expected add-on containing
only `SKILL.md` therefore passes the generic required-file check as `valid`.

**Why it matters:** This violates the ratified condition that an expected installed unit with
a missing or unreadable marker must be unhealthy. Deleting the marker together with the
manifest and inventoried payload defeats the integrity anchor and restores the gutted-tree
false green that B3 was intended to close.

**Evidence:** The measured results are
`{"missingMarker":{"doctorOk":true,"status":"valid","problems":[]},"unreadableMarker":{"doctorOk":true,"status":"valid","problems":[]}}`.
The control flow is explicit: `readMarker()` returns `null`, `manifestProblems()` returns `[]`,
and `validateInstalledPayload()` requires only `SKILL.md` for an add-on.

**Fix:** Preserve marker read state so missing and unreadable are distinguishable, and fail
health for either state whenever the unit is expected from the core install selection or an
explicit `--only` request. Retain legacy compatibility only for a readable, installer-owned
legacy marker. Add negative tests that delete and corrupt the expected add-on marker and
assert malformed/unhealthy status plus non-zero `doctor`.

### [MAJOR] Preflight misses predictable copy failures and leaves a partial selected set

**Where:** `lib/installer.js:897-922`; `lib/installer.js:957`;
`lib/installer.js:1064-1075`; `lib/installer.js:1171-1175`

**What:** `preflightSkillUnit()` traverses a source add-on only when it carries a manifest.
Manifest-free add-ons are permitted, but their source tree is first traversed by
`copyRecursive()` after earlier units have already been replaced. A symlink is therefore a
predictable source defect that passes preflight and is rejected only during the sequential
write loop.

**Why it matters:** FINAL.md B5 requires every unit and destination to be preflighted before
the first write and requires zero writes for a predictable failure. The implementation can
return `ok:false` after installing the core and every preceding add-on, which is precisely the
partial-fleet state B5 forbids. IMPLEMENTATION.md:74-75 overstates the implemented guarantee.

**Evidence:** I staged an otherwise valid package with a final manifest-free add-on
`zz-broken` containing `SKILL.md` plus a symlink. Installation failed with
`Refusing to copy symlink in skill payload`, but the destination already contained
`parley-deck`, `parley-bidding`, `parley-design`, `parley-design-check`, `parley-tracker`, and
`parley-worktrees`.

**Fix:** Add a generic, read-only copyability traversal for every source unit during
preflight (manifested or not), including symlink and readability checks. Prefer staging and
validating every selected unit before replacing the first destination. Add a regression test
whose late manifest-free unit contains a symlink and assert that the skills destination
remains nonexistent.

### [MINOR] D-2's evidence does not establish that the source was untouched

**Where:** `IMPLEMENTATION.md:127-132`

**What:** D-2 calls file-count and cache scans a "byte-level check" and concludes that the
source was untouched. Neither check observes the bytes of any ordinary source file, and the
later source-vs-integrated comparison is not a before/after source baseline.

**Why it matters:** The read-only source constraint is a provenance guarantee. A content edit
that preserves the 48 paths and creates no cache passes the cited checks, so the evidence does
not support the attached claim even though I found no indication that such an edit occurred.

**Evidence:** The documented replacement records only `48 files, 0 caches` before and after.
My current source comparison confirms the documented path/content-difference count, but it
cannot reconstruct the source's pre-copy bytes.

**Fix:** Compare a sorted path-plus-SHA-256 inventory captured before copying with one captured
afterwards, or narrow D-2 to the facts actually established and stop calling it byte-level or
proof that the source was untouched.

### [MINOR] The Python guard accepts the backslash composition D-3 says it refuses

**Where:** `test/design-addons.test.js:1109-1113`;
`skills/parley-bidding/SKILL.md:111-117`

**What:** The extractor marks a command assembled from backslash-continued physical lines by
appending `\`, but the Python arm removes that marker before applying `PUBLISHED_PYTHON`.
Consequently both shipped multiline commands are accepted. This contradicts FINAL.md F5
(`shell syntax refused`) and IMPLEMENTATION.md D-3's factual claim that `\` remains refused.

**Why it matters:** The static arm remains safe from command execution, and accepting `<` and
`>` as placeholder characters is justified for that reason. But the documented grammar and
the guard's actual acceptance language differ, so this release's evidence claims a stricter
publication contract than the test enforces.

**Evidence:** `npm test` passes the Python published-command test while `SKILL.md` publishes
two backslash-continued commands. The acceptance path is explicit at lines 1112-1113:
`replace(/\s*\\$/, "")` discards the extractor's splicing sentinel before matching.

**Fix:** Either preserve the sentinel so a continued command is refused and rewrite the two
published commands as one line each, or ratify a narrow continuation exception and update F5,
D-3, the comments, and grammar tests to state it honestly.

## What I verified and found correct

- `node --test test/bidding-addon.test.js`: 23/23 pass.
- `npm test`: 278 Node tests and 54 Python tests pass; manifest check passes.
- `node scripts/build-addon-manifest.js --check`: 47 inventoried payload files and aggregate
  `sha256:7854adf1…`.
- `npm pack --dry-run` with an isolated writable npm cache: 202 files, including 48 under
  `skills/parley-bidding/`; prepack ran and no cache artifact was listed.
- The shipped adapter validator reports all four adapters valid.
- A readable current-schema marker anchors both the manifest hash and aggregate payload digest.
- With that marker intact, manifest deletion, payload drift, and a self-consistent
  manifest-plus-payload replacement are rejected.
- A current-schema marker missing its `manifest` field is rejected rather than treated as
  legacy.
- The shipped manifest declares Python `>=3.10`, and the npm test runner fails rather than
  skips when that interpreter requirement is unavailable.
- The Node extractor arm retains the original `node`/`--test` patterns and existing assertion
  behavior after parameterization. The Python arm's non-executing treatment makes literal
  `<placeholder>` arguments safe to accept.
- Source comparison found 48 source files, no Python cache, one dropped `.gitignore`, one
  added manifest, and the documented nine content-different paths.
- D-4 honestly drops runtime-exposure claims as NOT TESTED, and D-5's six-skill wording is
  present in README.md. The schema-identity change in F1 is also explicitly recorded as
  compatibility NOT TESTED.

## Open questions for the implementer

- What backward-compatible JSON shape will expose payload validity separately from operational
  availability while keeping `doctor` exit behavior unambiguous?
- Will expected-unit marker enforcement be centralized in `skillUnitStatus()` so every add-on,
  not only `parley-bidding`, receives the same fail-closed behavior?
- Will B5 be implemented as complete preflight traversal or as stage-all-before-replace? The
  latter also catches future deterministic staging failures without growing a second copy of
  `copyRecursive()`'s refusal rules.
