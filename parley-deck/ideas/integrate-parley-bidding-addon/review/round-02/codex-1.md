---
idea: integrate-parley-bidding-addon
review-round: 2
agent: codex-1
date: 2026-07-30
reviewed-commit: 89069b0
---

## Verdict

BLOCK

## Round-1 findings — closed or not

1. **`doctor` does not enforce the add-on's declared Python minimum — NOT FULLY
   CLOSED.** The normal one-shot CLI path is fixed: with `python3` absent I measured
   `doctorOk:false`, payload `status:"valid"`, and
   `runtime:{ok:false,requirement:">=3.10",detail:"python3 is not available..."}`. The real
   CLI exits 1 in both text and JSON modes. However, the probe ignores the effective
   environment accepted by `installer.run`; the resulting false-green path is the MAJOR
   finding below.
2. **Expected installer-managed units pass health checks without a readable marker —
   CLOSED.** My fresh matrix produced `doctorOk:false` and `status:"malformed"` for a
   deleted marker, a mode-000 unreadable marker, invalid marker JSON, a foreign marker, and
   the `SKILL.md`-only gutted tree with its marker deleted. Missing and unreadable markers
   have distinct problem messages.

The other three observations in my round-1 file (manifest-free preflight atomicity, D-2
evidence wording, and the backslash-continuation grammar) were not part of fix-up cycle 1
and this diff does not touch their code or evidence. Per the round-2 instruction, I did not
re-litigate their accepted disposition.

## New findings

### [MAJOR] The Python probe ignores the runner's effective environment and can still return healthy

**Where:** `lib/installer.js:310-318`, `lib/installer.js:1263-1275`,
`lib/installer.js:1288-1302`

**What:** `makeContext` makes `io.env` the command's effective environment, and the rest of
the installer uses `context.env` for command discovery and subprocesses. `probePython3`
instead calls `spawnSync` without `env`, so it searches the parent process environment. Its
process-global cache also cannot distinguish two command contexts with different effective
environments.

**Why it matters:** This leaves a direct false-green through the exported `run(argv, io)`
entry point: the target environment can have no `python3`, while `doctor` probes a different
environment and exits 0. B6 is an operational-availability guarantee, so checking the wrong
environment does not close the round-1 defect.

**Evidence:** In a fresh Node process I installed and ran:

`installer.run(["doctor","--target","codex","--json"], {env:
{...process.env, HOME: isolatedHome, PATH: ""}, ...})`

The parent process had Python 3.14. The measured result was:

```json
{
  "exitCode": 0,
  "doctorOk": true,
  "status": "valid",
  "runtime": {
    "ok": true,
    "requirement": ">=3.10",
    "detail": "python3 3.14"
  }
}
```

Running the real CLI in a child whose actual `PATH` was empty correctly produced exit 1 and
`runtime.ok:false`, proving that the discrepancy is environment selection rather than
payload state.

**Fix:** Thread the command context into runtime probing and pass `env: context.env` to
`spawnSync`. Memoize once per command invocation, or key the cache by all environment fields
that affect executable resolution, rather than using one process-global answer. Add a test
where the parent has Python but `installer.run` receives `PATH:""`; it must exit 1. Give the
probe a finite timeout while changing this code.

### [MINOR] Runtime probing leaks into `paths`, while text `status` hides the result

**Where:** `lib/installer.js:353-359`, `lib/installer.js:384-400`,
`lib/installer.js:1212-1254`, `lib/installer.js:1481-1525`

**What:** Runtime availability is computed unconditionally inside `skillUnitStatus`, which
is shared by `doctor`, `status`, and `paths`. As a result, `paths` now launches `python3` and
prints an availability diagnostic. Conversely, `status` launches the same probe but its
text renderer prints only `addon parley-bidding: valid`, discarding `runtime.ok:false`.

**Why it matters:** A path-discovery command now executes a PATH-resolved program and can
block indefinitely because the probe has no timeout. The command intended to summarize
status performs that work but presents a false-looking human summary. This is a collateral
CLI regression and makes equivalent JSON/text invocations disagree about actionable state.

**Evidence:** I installed into an isolated home with a `python3` stub that writes a sentinel.
`paths --target codex --json` exited 0 and created the sentinel. With no `python3`, the
measured text outputs were:

- `status`: `addon parley-bidding: valid` with no unavailable line.
- `paths`: `parley-bidding: valid` followed by `unavailable: python3 is not available...`.
- `doctor`: the same separate availability line and exit 1, which is correct.

**Fix:** Make runtime probing an explicit option of the status traversal. Enable it for
`doctor` and, if intended, `status`; disable it for `paths`. If `status` probes runtime,
render the runtime failure in text and incorporate an appropriate action/summary instead of
printing only `valid`. Add negative tests for both commands.

## What I verified and found correct

- The reviewed worktree was exactly `89069b0913cc`; the fix diff changes only
  `lib/installer.js` and `test/bidding-addon.test.js`.
- My round-1 matrix now gives:
  - no Python: payload `valid`, runtime unavailable, `doctorOk:false`;
  - missing marker: `malformed`, `doctorOk:false`;
  - unreadable marker: `malformed`, `doctorOk:false`;
  - gutted to `SKILL.md` plus marker deletion: `malformed`, `doctorOk:false`.
- A foreign marker is `malformed`. A readable installer-owned legacy marker with no
  `markerSchema` remains `valid`.
- An intact `--no-addons` install reports only valid core. An install using
  `--only parley-tracker` reports core plus tracker and does not invent a missing
  `parley-bidding` unit.
- Default install, `--only parley-bidding`, core-only install, and add-on uninstall tests
  pass; I found no install or uninstall regression in this fix diff.
- Interpreter stubs below 3.10 fail and at 3.10 pass in the new tests. Add-ons without a
  runtime declaration retain `runtime:null`.
- The actual `doctor --json` shape is unambiguous on the normal CLI path: top-level `ok` is
  false, payload `status` stays `valid`, `runtime.ok` is false, and the process exits 1.
- `node --test test/bidding-addon.test.js`: 31/31 pass.
- Final `npm test`: 286/286 Node tests and 54/54 Python tests pass; the manifest check reports
  47 files and aggregate `sha256:7854adf1...`.
- The reviewed repository remained clean, and no `__pycache__`, `.pyc`, or `.pyo` appeared
  under `skills/parley-bidding/`.

## Open questions for the implementer

- One earlier `npm test` run produced four simultaneous failures in the new marker tests
  (282/286), although six immediate full Node-suite reruns and the final full npm run all
  passed. I could not reproduce it again or tie it to a code path, so I have not elevated it
  to a finding. Is there any known concurrent test or filesystem mutation that could explain
  that transient result?
