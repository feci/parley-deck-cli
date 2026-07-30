---
idea: integrate-parley-bidding-addon
review-round: 1
agent: hermes-1
date: 2026-07-30
reviewed-commit: 714712f
---

## Verdict
BLOCK

## Findings

### [CRITICAL] Expected installed unit with missing or unreadable marker is not reported unhealthy
**Where:** lib/installer.js:1255-1259 (`manifestProblems`), lib/installer.js:1315-1321 (`readMarker`), lib/installer.js:1208-1222 (`skillUnitStatus`)
**What:** When an installed add-on unit is expected (the core marker recorded it as an installed add-on) but its own marker file is missing or unreadable, `doctor` reports `status: valid` with an empty `problems` array and exits 0. Measured directly at commit 714712f:

```
{"noPython":         {"doctorOk": true, "status": "valid", "problems": []},
 "missingMarker":     {"doctorOk": true, "status": "valid", "problems": []},
 "unreadableMarker":  {"doctorOk": true, "status": "valid", "problems": []}}
```

**Why it matters:** This is the double-deletion case I named in round-03 (hermes-1.md finding 1, "Delete the manifest and the marker together") and that codex-1's ratified condition 1 closes: "An expected installed unit with a missing or unreadable marker must also be unhealthy." If someone guts an add-on down to just a SKILL.md and deletes the marker file alongside it, `doctor` still calls the unit `valid` and returns a clean exit. The health check provides no signal that an expected, installed unit is operationally unavailable.

**Evidence:** The code path is: `skillUnitStatus` (line 1208) calls `readMarker(unit.dest)` then `validateInstalledPayload(unit.dest, unit.kind, { collect: true })`. For addons, `validateInstalledPayload` (line 1224) checks only `["SKILL.md"]` (line 1227) and delegates integrity to `manifestProblems` (line 1236). `manifestProblems` (line 1255) reads the marker via `readMarker` (line 1256). `readMarker` (line 1315) wraps `JSON.parse(fs.readFileSync(...))` in a try/catch that returns `null` on any error — missing file, unreadable file, or unparseable content all collapse to `null`. Back in `manifestProblems`, line 1257: `if (!marker || marker.name !== PACKAGE_JSON.name) return [];` — a `null` marker returns an empty problems array. With SKILL.md still present, `missing.length === 0 && problems.length === 0` is true (line 1237), so `skillUnitStatus` returns `status: "valid"` (line 1217). `doctorCommand` (line 362) then sees every skill as `valid` and sets `ok: true` (line 374), exiting 0.

**Fix:** In `manifestProblems`, distinguish "no marker because we never installed this" (the current `return []` path, legitimate for non-installer-owned directories) from "marker was expected but is now missing or unreadable" (a defect for an expected installed unit). The expected-installed determination is already available: `targetSkillUnits` (line 874) only includes add-ons the core marker recorded as installed. When `skillUnitStatus` is called for such a unit and `readMarker` returns `null`, it should emit a problem (e.g. `install marker is missing or unreadable`) and set `status: "malformed"`, causing `doctor` to exit non-zero. The absence of a marker on an *expected* unit is a health defect, not a no-op.

### [CRITICAL] Python interpreter minimum floor check is in the npm test runner, not in `doctor`
**Where:** scripts/run-python-tests.js:47-75 (`declaredPythonFloor` / `resolveInterpreter`) vs. lib/installer.js:362-378 (`doctorCommand`) and lib/installer.js:1224-1248 (`validateInstalledPayload`)
**What:** FINAL.md B6 requires `doctor` to distinguish payload-valid from operationally unavailable and to fail health when the declared interpreter minimum is missing. The implementer placed the Python floor check (verifying that the interpreter meets the declared `>=3.10` minimum from `parley-addon.json` `runtime.python`) in `scripts/run-python-tests.js` only. The `doctor` command itself does not check for the interpreter minimum.

**Why it matters:** `doctor` is the operational health gate. If the interpreter floor is only enforced during `npm test`, then a deployment where the Python minimum is absent passes `doctor` cleanly — the `noPython` measurement above shows `doctorOk: true, status: valid, problems: []`. B6's requirement that `doctor` fail health when the declared interpreter minimum is missing is unmet. An agent running `doctor` to decide whether the bidding add-on is safe to use gets a false green.

**Evidence:** `run-python-tests.js:47-56` defines `declaredPythonFloor()` which reads `parley-addon.json` and parses `runtime.python` (`>=3.10`). `resolveInterpreter` (line 58) spawns `python3 -c "import sys; print('%d.%d' % sys.version_info[:2])"` and calls `fail()` (exit 1) if the interpreter is absent (line 62-68) or below the floor (line 70-73). This check runs only when `npm test` invokes the script. On the `doctor` side, `doctorCommand` (line 362) calls `targetStatus` → `skillUnitStatus` → `validateInstalledPayload`, none of which reference `python3`, `runtime`, or the interpreter at all. `validateInstalledPayload` (line 1224) checks required files and manifest integrity only — it never probes for the interpreter. The `noPython` fixture confirms this: interpreter absent, yet `doctorOk: true`.

**Fix:** Add an interpreter-floor check to the `doctor` health path. When `validateInstalledPayload` (or `manifestProblems`) processes an addon whose installed manifest declares `runtime.python`, `doctor` should spawn `python3 --version` (or equivalent), compare the result against the declared floor, and if the interpreter is absent or below minimum, emit a problem (e.g. `interpreter missing` / `interpreter below declared minimum >=3.10`) and mark the unit unhealthy. This mirrors the logic already in `run-python-tests.js:58-75` but routes the result through `doctor`'s health verdict rather than the test runner's exit code.

## What I verified and found correct

- The three measurement scenarios (`noPython`, `missingMarker`, `unreadableMarker`) were run directly against commit 714712f and produce the documented results. The measurements are reproducible.
- The `doctorOk: true` / `status: valid` / `problems: []` triple is consistent across all three fixtures, confirming that neither the marker check nor the interpreter floor check is wired into the `doctor` health path.
- Traced the full `doctor` code path through the source at commit 714712f: `doctorCommand` (installer.js:362) → `targetStatus` (line 1193) → `skillUnitStatus` (line 1208) → `validateInstalledPayload` (line 1224) → `manifestProblems` (line 1255). The missing-marker/unreadable-marker path is confirmed at the source level: `readMarker` (line 1315) returns `null` for all three failure modes (missing, unreadable, unparseable), and `manifestProblems` line 1257 returns `[]` for `null` markers.
- Confirmed the Python floor check exists only in `scripts/run-python-tests.js:47-75` (`declaredPythonFloor` + `resolveInterpreter`), which is invoked by `npm test`, not by `doctor`. No reference to `python3`, `runtime`, or interpreter version exists anywhere in the `doctor`/`targetStatus`/`skillUnitStatus`/`validateInstalledPayload`/`manifestProblems` call chain.
- The existing test suite (test/bidding-addon.test.js) covers the gutted-tree case (line 188: "a tree gutted to SKILL.md is malformed") and manifest deletion (line 241), but does not test the marker-deleted case for an expected installed unit — which is precisely the gap that allows the `missingMarker`/`unreadableMarker` measurements to pass.

## Open questions for the implementer

1. Is there a design reason the marker-missing/unreadable case was treated as valid rather than unhealthy? If the intent was "only flag marker issues for explicitly marker-gated units," the manifest/declaration should still distinguish expected-installed from optional, and expected-installed with no marker should fail.
2. Was the Python floor check intentionally placed only in the npm test runner, or was wiring it into `doctor` an oversight? B6 appears to require it in `doctor` specifically.
3. Once both checks are added to `doctor`, what exit code and problem taxonomy will be used? (e.g. `marker missing`, `marker unreadable`, `interpreter missing`, `interpreter below minimum`) — confirming the taxonomy will help align the test fixtures.
