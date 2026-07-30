---
idea: integrate-parley-bidding-addon
review-round: 2
agent: hermes-1
date: 2026-07-30
reviewed-commit: 89069b0
---

## Verdict
BLOCK

## Round-1 findings — closed or not

### Finding 1 (CRITICAL): Expected installed unit with missing or unreadable marker reported `valid`

**CLOSED.** The fix is in `skillUnitStatus` (lib/installer.js:1212-1255), centralized so every add-on
gets it, not just `parley-bidding`. `readMarkerState` (line 1402-1416) distinguishes absent
(`present: false`) from unreadable (`present: true, readable: false`) from foreign-name
(`state.marker.name !== PACKAGE_JSON.name`). Each case emits a distinct problem string and sets
`status: "malformed"`.

Re-measured at 89069b0 with the round-1 probe matrix (child process, python3 isolated so the
runtime check does not confound the marker result):

```
missingMarker:     {ok:false, status:"malformed",
                    problems:["no parley-deck-skill install marker: ..."]}
unreadableMarker:  {ok:false, status:"malformed",
                    problems:["the parley-deck-skill install marker is unreadable or is not valid JSON"]}
gut+deleteMarker:  {ok:false, status:"malformed",
                    problems:["no parley-deck-skill install marker: ..."]}
foreignMarker:     {ok:false, status:"malformed",
                    problems:["the install marker belongs to \"some-other-tool\", not parley-deck-skill"]}
untouched:         {ok:true,  status:"valid", problems:[]}
```

Additional edge cases I probed and confirmed correct:
- Marker as a directory → `present: false` (fileExists checks isFile) → "no marker" message.
- Dangling symlink → `present: false` → "no marker" message.
- JSON array (`[1,2,3]`) → `present: true, readable: false` → "unreadable" message.
- JSON `null` → `present: true, readable: false` → "unreadable" message.
- Legacy 2.0.0 marker (no `markerSchema` field) on the add-on → `status: "valid"`, no marker
  problem, `manifestProblems` returns `[]` at the legacy-skip path (line 1343-1347). Confirmed
  healthy.
- Legacy 2.0.0 core marker (no `addons` field) → `markerAddonNames` returns `[]`, so
  `expectedAddonNames` returns `[]`, so only the core skill is checked. `doctor` reports `ok:
  true` with only `parley-deck` in the skill list. No false "missing addon" problem.
- `--no-addons` install → `doctor` reports `ok: true`, only core skill, `runtime: null`.
- `--only parley-worktrees` install → `doctor` reports `ok: true`, core + parley-worktrees only,
  no `parley-bidding` unit. No false "missing" problem.

The marker check does not break any legitimate case I could find.

### Finding 2 (CRITICAL): Python interpreter minimum floor check in the npm test runner, not in `doctor`

**CLOSED.** `runtimeAvailability` (lib/installer.js:1263-1286) reads `runtime.python` from the
validated manifest via `addonManifest.readManifest`, probes `python3` via `probePython3`, and
returns `{ok, requirement, detail}`. `doctorCommand` (line 376-378) now gates health on both
`status === "valid"` AND `(!skill.runtime || skill.runtime.ok)`. The payload stays `valid` and
the unit gets `runtime: {ok:false, detail}` — the two answers are separate, per B6. `doctor`
exits non-zero.

Re-measured at 89069b0:

```
noPython (PATH has no python3):   {ok:false, status:"valid", problems:[],
                                   runtime:{ok:false, detail:"python3 is not available, but this skill requires >=3.10"}}
python 3.9 (below floor):         {ok:false, status:"valid", problems:[],
                                   runtime:{ok:false, detail:"python3 is 3.9, but this skill requires >=3.10"}}
python 3.14 (meets floor):        {ok:true,  status:"valid", problems:[],
                                   runtime:{ok:true,  detail:"python3 3.14"}}
```

The probe cache (`let pythonProbe = null`, line 1290) is correct: instrumented `spawnSync` and
confirmed exactly 1 `python3` call across 14 targets with `parley-bidding` in `doctor --target
all --include-undetected`. The cache is per-process, which is the right granularity for a single
CLI invocation.

The `doctor --json` shape is unambiguous: each skill entry has `status` (string), `problems`
(array), `missing` (array), `runtime` (null or `{ok, requirement, detail}`), and `marker` (the
parsed marker or null). The top-level `ok` is a single boolean. A consumer can distinguish
"payload broken" (`status !== "valid"`) from "payload valid but unavailable" (`status ===
"valid"` and `runtime && !runtime.ok`).

## New findings

### [MAJOR] Pre-existing test `doctor reports add-on skills per target` fails when system python3 < 3.10

**Where:** test/installer.test.js:537-547
**What:** The fix-up introduced the `probePython3` runtime check into `skillUnitStatus`, which
`doctorCommand` calls. The pre-existing test at line 537 installs all add-ons (including
`parley-bidding`) and asserts `result.ok === true` (line 542). But `probePython3` (line 1291-
1303) calls `spawnSync("python3", ...)` with no `env` option, so it inherits `process.env.PATH`
— not the test context's `env: { PATH: "" }`. On any machine where the system `python3` is below
the declared `>=3.10` floor, `doctor` reports `ok: false` and the assertion fails.

This test was passing at 714712f and fails at 89069b0. Confirmed by checking out 714712f's
`lib/installer.js` and `test/installer.test.js` and running the test in isolation: PASS at
714712f, FAIL at 89069b0.

**Why it matters:** The implementer claims "npm test → 286 node + 54 Python, 0 fail"
(IMPLEMENTATION.md line 162). That claim is false on any machine with python3 < 3.10. On this
review machine, `npm test` produces 285 pass / 1 fail, not 286 / 0. The default macOS python3 at
`/usr/bin/python3` is 3.9.6 — a developer on macOS running `npm test` gets a failure. The CI
workflow (`.github/workflows/test.yml`) uses `setup-python` with 3.10 and 3.13, so CI passes,
but the claim of "0 fail" is environment-dependent and the regression is real.

The root cause is that `probePython3` reads `process.env`, not `context.env`. The new tests in
`test/bidding-addon.test.js` correctly work around this by running `doctorCommand` in a child
process with a controlled `PATH` (`doctorInChildWithPath`, line 433-465). But the pre-existing
test in `installer.test.js` was not updated to account for the new runtime check.

**Evidence:**
```
$ python3 --version
Python 3.9.6

$ node --test --test-name-pattern="doctor reports add-on skills per target" test/installer.test.js
✖ doctor reports add-on skills per target
  expected: true, actual: false  (at line 542)

$ PATH="/opt/homebrew/bin:$PATH" node --test  # python3 3.14
ℹ pass 286  ℹ fail 0

$ node --test  # python3 3.9 (system default)
ℹ pass 285  ℹ fail 1
```

**Fix:** Update `test/installer.test.js:537` to either (a) run in a child process with a
controlled `PATH` that has a python3 >= 3.10 (matching the pattern in
`test/bidding-addon.test.js:doctorInChildWithPath`), or (b) change the assertion from
`result.ok === true` to check only `skill.status === "valid"` for each skill (which is what the
test actually cares about — add-on skills appearing per target). Option (b) is the smaller
change and matches the test's intent: it verifies the skill *list*, not the runtime health.

Alternatively, `probePython3` could accept the context and use `context.env` when available,
falling back to `process.env` for the real CLI. But that is a larger change and the current
behavior (reading the real environment) is correct for production usage — the issue is test
isolation, not production logic.

## What I verified and found correct

1. **F1 marker check is complete and centralized.** `skillUnitStatus` (line 1212-1255) handles
   all marker states: absent, unreadable, foreign-name, and valid. Each emits a distinct problem
   string. The `problems` array is assembled from both `validation.problems` (manifest integrity)
   and the marker checks, so a unit can report multiple problems. The `ok` boolean requires both
   `validation.missing.length === 0` and `problems.length === 0`.

2. **F2 runtime check is correct in production.** `runtimeAvailability` reads the manifest, parses
   the `>=X.Y` floor, probes `python3` once (cached per process), and returns a structured
   result. The floor comparison is correct: `probe.major < wanted[0] || (probe.major ===
   wanted[0] && probe.minor < wanted[1])`. An unsupported spec format (not `>=X.Y`) returns
   `{ok: false, detail: "unsupported python requirement ..."}` rather than silently passing.

3. **B6 separation is clean.** A byte-valid payload with an absent or below-floor interpreter
   reports `status: "valid"` and `runtime: {ok: false}`. `doctor` exits non-zero. The text
   output shows `unavailable: <detail>` as a separate line from `integrity:` problems. The
   stderr summary distinguishes "missing or malformed" from "installed but operationally
   unavailable".

4. **`doctor --json` shape is unambiguous.** Top-level `ok` is a boolean. Each skill has
   `status`, `problems`, `missing`, `runtime`, and `marker`. `runtime` is `null` for units that
   declare no runtime requirement or when the payload is malformed. It is `{ok, requirement,
   detail}` for units with a runtime declaration and a valid payload. A consumer can
   programmatically distinguish all cases.

5. **No regression in install, uninstall, status, or paths.** `installCommand` and
   `uninstallCommand` do not call `skillUnitStatus` (they use `installSkillUnit` /
   `uninstallSkillUnit`). `pathsCommand` calls `targetStatus` → `skillUnitStatus` but always
   returns `ok: true` (line 356), so the runtime check does not affect its exit code.
   `statusCommand` also always returns `ok: true` (line 392). The `runtime` field flows through
   to `status` output via `enrichRuntimeStatus` (line 458-465, which spreads `...status`), but
   this is informational and does not break the `status` command's contract.

6. **Legacy 2.0.0 markers are not broken.** A marker with no `markerSchema` field is treated as
   legacy (line 1343-1347): `manifestProblems` returns `[]`, and the marker check in
   `skillUnitStatus` passes (the marker is present, readable, and has `name ===
   "parley-deck-skill"`). The legacy core marker with no `addons` field is treated as core-only
   (`markerAddonNames` returns `[]`), so no add-ons are expected and `doctor` reports healthy.

7. **`--no-addons` and `--only` installs are healthy.** `doctor` derives the expected set from
   the installed core marker (line 866-868), so intentionally omitted add-ons are not flagged as
   missing or malformed. Verified with `--no-addons` (only core skill, `ok: true`) and `--only
   parley-worktrees` (core + worktrees, no bidding, `ok: true`).

8. **Probe cache is correct.** Instrumented `spawnSync` and confirmed 1 `python3` call across 14
   targets with `parley-bidding` in a single `doctor --target all` invocation. The module-level
   `let pythonProbe = null` is reset per process (per CLI invocation), which is the right
   granularity.

9. **The eight new tests in `test/bidding-addon.test.js` are well-designed.** The interpreter
   tests run in a child process with a controlled `PATH` (`doctorInChildWithPath`), correctly
   isolating the probe from the system environment. The marker tests use direct `doctorCommand`
   calls (no child process needed since they don't depend on the runtime probe — `runtime` is
   `null` when `status` is `malformed`). All 31 tests in the file pass.

10. **Python suite passes when python3 >= 3.10.** With python3 3.14 on PATH: 54 tests across 7
    files, 0 fail. The manifest check also passes (47 files, aggregate sha256 verified).

## Open questions for the implementer

1. The claim "npm test → 286 node + 54 Python, 0 fail" is true only when the system python3 is
   >= 3.10. On macOS with the default `/usr/bin/python3` (3.9.6), `npm test` produces 285 node
   pass / 1 fail, and the Python runner exits 1 before running any Python tests. Was this tested
   only on machines with python3 >= 3.10, or was the pre-existing test failure observed and
   not addressed?

2. `probePython3` (line 1291-1303) calls `spawnSync` with no `env` option, so it reads
   `process.env.PATH`, not `context.env.PATH`. In the real CLI this is correct (the context env
   IS `process.env` via `makeContext`, line 311). But it means any test that calls
   `doctorCommand` with a context that sets `env.PATH = ""` will still have the probe find the
   real `python3`. Is this intentional? If so, should the pre-existing test at
   `installer.test.js:537` be updated to account for it?
