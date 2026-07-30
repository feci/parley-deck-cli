---
idea: integrate-parley-bidding-addon
review-round: 3
agent: hermes-1
date: 2026-07-30
reviewed-commit: 5c324ef
---

## Verdict

ACCEPT WITH CONDITIONS

## Round-2 findings — closed or not

### F4 (codex-1 MAJOR) — probe resolved against the wrong environment — CLOSED

`probePython3` (lib/installer.js:1296-1316) now accepts `env`, resolves against it
(`const effective = env && typeof env === "object" ? env : process.env`), passes it to
`spawnSync` as `env: effective`, and is memoized per-PATH via a `Map` keyed on
`effective.PATH` rather than a single process-global scalar. The 5s `timeout: 5000` is
present. `doctorCommand` (line 374) and `statusCommand` (line 388) both thread
`{ probeRuntime: true, env: context.env }` through `targetStatus` → `skillUnitStatus` →
`runtimeAvailability` → `probePython3`.

Re-measured the direct regression: installed into an isolated HOME, then called
`installer.run(["doctor","--target","codex","--json"], { env: { ...process.env, HOME,
PATH: "" } })` in-process. The parent has python3 3.9.6; the declared environment has no
python3. Result: `bidding.status: "valid"`, `bidding.runtime.ok: false`, `result.ok:
false`, exit 1. Before the fix this path returned `runtime.ok: true` and exit 0 — codex-1's
exact measurement. The per-PATH memoization test ("the interpreter probe is memoized per
PATH") stubs two different PATHs in one process and confirms two different answers, which a
process-global cache could not produce. Confirmed passing.

### F5 (hermes-1 MAJOR) — pre-existing test asserted result.ok — CLOSED

test/installer.test.js:539-547: the assertion `assert.equal(result.ok, true)` is replaced
with `assert.deepEqual(result.targets[0].skills.map(...), ["parley-deck","parley-bidding",
...])` plus a per-skill `status === "valid"` check. The comment explains why: `ok` now
includes availability, and this context declares `PATH: ""`, so the probe correctly finds
no interpreter. The test now checks the skill LIST, which is what it was always about.

Re-measured on this machine (python3 3.9.6, the default `/usr/bin/python3`):
`node --test` → 290 pass / 0 fail. The one test that was my round-2 finding now passes
without any PATH tricks. The implementer's claim of "290/290 with 3.9.6 AND with 3.14" is
reproduced — I am on 3.9.6 and see 290/0.

One of the implementer's own new tests ("run() honors the caller's environment…") had the
same ambient-interpreter assumption and was fixed in the same cycle: it now stubs both arms
rather than relying on the host meeting the floor. The memoization test does the same.

### F6 (codex-1 MINOR) — probing leaked into paths, status discarded the result — CLOSED

`targetStatus` now takes an `options` argument; `skillUnitStatus` checks
`ok && options && options.probeRuntime` before calling `runtimeAvailability` (line 1255).
`pathsCommand` (line 358) passes `{ probeRuntime: false }`; `doctorCommand` and
`statusCommand` pass `{ probeRuntime: true, env: context.env }`. The `status` text renderer
(line 1541-1543) now prints `unavailable: <detail>` when `skill.runtime && !skill.runtime.ok`,
and the `integrity:` lines (line 1544-1546) for problems. Re-measured: `paths` with a
sentinel `python3` stub on PATH does not create the sentinel file; `doctor` in the same
environment does. The "paths does not launch an interpreter" test confirms this directly.

### F7 (kimi-1 NIT) — stderr wording — CLOSED

Line 1521: `reasons.push("operationally unavailable")` — the "installed but" collision is
gone. The full sentence now reads "One or more installs are missing or malformed, or
operationally unavailable." Measured in the `noPython` path.

### codex-1's transient 4-test failure — addressed, not closed

The process-global probe cache was the most plausible mechanism: a test that ran first with
one PATH could poison the cache for a later test with a different PATH. Keying the cache per
PATH (F4's fix) removes that mechanism. The implementer records it as "unexplained rather
than closed" — I agree with that disposition. I ran `node --test` five times in succession
and saw 290/0 every time; I cannot reproduce the transient either.

## Ruling on the skills-CLI question

**(b)** — report valid-but-not-installed-by-this-tool when the marker is absent but the tree
carries a manifest that fully verifies, while a tree with neither manifest nor marker (where
the source ships a manifest) stays malformed.

I reproduced the implementer's measurement end-to-end. I installed `parley-bidding` from the
local repository via the universal `skills` CLI (v1.5.21) with `--agent claude-code
--copy`, then ran `doctor --target claude`:

```
claude/parley-bidding: malformed …/.claude/skills/parley-bidding
  integrity: no parley-deck-skill install marker: this directory was not installed by this tool, or the marker was removed
exit=1
```

The skills CLI copies the tree faithfully — `parley-addon.json` is present — and writes no
`.parley-deck-skill-install.json`. I confirmed `verifyPayload` on that directory returns
`{ ok: true, problems: [] }`: the payload is provably byte-intact, all 47 files match their
hashes, the aggregate digest matches, and there are no undeclared files. `doctor` calls it
malformed anyway, purely because the marker is absent.

Why (b), not (a):

1. The README's Install section leads with the `skills` CLI in a `[!TIP]` callout —
   "One command, most agents" — and presents this package's own installer as the second
   option. A tree installed by the path we recommend first is not unhealthy to this tool's
   users; calling it `malformed` with exit 1 is a false alarm on our own documented
   primary install path.

2. The marker exists to anchor the manifest-to-payload binding: it records which manifest
   was installed so a self-consistent swap is detectable. But when `verifyPayload` returns
   `ok: true`, the manifest beside the payload already proves byte-integrity — stronger
   evidence than the marker provides. The marker's own binding check
   (`verified.manifest.aggregate !== declared.aggregate`) is moot when there is no
   `declared` to compare against; the only question left is "is the payload intact?", and
   `verifyPayload` answers it directly. Reporting `malformed` discards stronger evidence
   in favour of weaker.

3. The gutting signal the marker check exists to close (F1: copy `SKILL.md` alone, delete
   the marker alongside) is preserved. In that case the installed tree has neither a
   verifying manifest (files are missing or modified) nor a marker — so it stays
   `malformed` under (b). The refinement only changes the verdict when the manifest
   verifies, which a gutted tree cannot.

4. The implementer notes the refinement "stays generic: whether the source ships a manifest
   is read from `unit.addon.root`". I confirmed `unit.addon` is populated by
   `targetSkillUnits` (line 886-894) from `discoverAddons`, and `addonManifest.hasManifest`
   on the source root correctly returns `true` for `parley-bidding` and `false` for the
   other four add-ons. The data needed to implement (b) is already in the unit object;
   no add-on is named in code.

One implementation note for whoever lands this: the new status should not be `valid` (which
implies installer-verified) nor `malformed`. A third value — e.g. `verified` or
`intact` — with a distinct problem/notice string ("payload verifies but was not installed
by this tool") keeps `doctor`'s exit code honest (the payload IS intact) without collapsing
the foreign-installer case into either the healthy or broken bucket. Whether `doctor` exits
0 or 1 for this state is a design call I leave to the implementer; I lean toward exit 0
(the payload is usable) with the notice printed, but would accept exit 1 if the implementer
argues that "not installed by this tool" is itself a health defect the user should act on.
Either way, it must not say `malformed` for a byte-perfect tree.

## New findings

### [MINOR] status always exits 0 even when an add-on is unavailable or malformed

**Where:** lib/installer.js:393 (`statusCommand` returns `ok: true` unconditionally)

**What:** `statusCommand` hardcodes `ok: true` regardless of what `targetStatus` finds. A
`status` run against a target with a `malformed` add-on or an `unavailable` runtime exits 0.
`doctor` exits 1 for the same directory. The text output now correctly prints
`unavailable:` and `integrity:` lines (F6's fix), but the exit code still says "everything
is fine."

**Why it matters:** This is not a regression from the fix-up cycle — `status` has always
returned `ok: true` (it is an informational command, not a health gate). But the fix-up
cycle made `status` probe the runtime and print the result, which creates a new
expectation: if `status` reports a problem, a script consuming its exit code will miss it.
A CI step that runs `status` instead of `doctor` (they are sibling commands) gets a false
green on an unavailable or malformed add-on. This is a pre-existing design choice that the
new probing makes more visible, not a defect introduced by this diff.

**Evidence:** Installed into an isolated HOME, deleted the marker from `parley-bidding`,
ran `status --target claude`: text output shows `addon parley-bidding: malformed` with the
`integrity:` line, but `echo $?` is 0. Same directory via `doctor`: exit 1.

**Fix:** Either (a) leave as-is and document that `status` is informational-only (exit 0
always) while `doctor` is the health gate (exit 1 on any problem), or (b) make `status`'s
`ok` reflect `skill.status === "valid" && (!skill.runtime || skill.runtime.ok)` the same way
`doctor` does. (a) is the smaller change and is consistent with the command's original
contract. I am not blocking on this — flagging it because the new probing makes the
discrepancy more surprising.

## What I verified and found correct

1. **290/290 node tests pass on python3 3.9.6.** This was my round-2 finding (F5): the
   pre-existing test failed on 3.9 and passed on 3.14. The fix is correct — the assertion
   now checks the skill list, not `result.ok`. I am on 3.9.6 and see 290/0, matching the
   implementer's claim exactly. The four new tests (run() honors caller env, memoization
   per PATH, paths does not launch interpreter, status reports unavailable in text) all
   pass.

2. **The probe resolves against `context.env`, not `process.env`.** Traced the full path:
   `makeContext` (line 311) sets `env: io.env || process.env`; `doctorCommand` (line 374)
   passes `{ probeRuntime: true, env: context.env }`; `skillUnitStatus` (line 1255) passes
   `options.env` to `runtimeAvailability` (line 1265); `runtimeAvailability` (line 1275)
   passes `env` to `probePython3` (line 1296); `probePython3` (line 1297) uses `effective`
   as the `spawnSync` `env`. The chain is unbroken. The direct regression test
   confirms it behaviourally.

3. **Per-PATH memoization is correct.** `pythonProbes` is a `Map` keyed on
   `String(effective.PATH)` with a sentinel `"\u0000inherit"` for `undefined`. Two contexts
   with different PATHs in one process get two different probes. The test stubs a 3.12
   python3 and an empty PATH and confirms the first returns `ok: true` and the second
   `ok: false`, with distinct `detail` strings. A process-global cache would have returned
   the first answer for both.

4. **`paths` does not probe.** `pathsCommand` (line 358) passes `{ probeRuntime: false }`.
   `skillUnitStatus` (line 1255) gates on `options && options.probeRuntime`, so `runtime`
   is `null` in the `paths` output. The sentinel test confirms no `python3` process is
   spawned. `doctor` in the same environment does spawn it.

5. **`status` now reports unavailable and integrity in text.** Lines 1541-1546 print
   `unavailable: <detail>` and `integrity: <problem>` for each add-on skill. The text
   output matches the JSON output for the same directory — the round-2 disagreement is
   gone.

6. **F1 (round-1 CRITICAL) stays closed.** The marker check in `skillUnitStatus`
   (lines 1237-1243) is unchanged from cycle 1: absent marker → "no parley-deck-skill
   install marker", unreadable → "the … marker is unreadable or is not valid JSON",
   foreign-name → "the install marker belongs to …". Each sets `status: "malformed"`.
   Re-measured: missing marker, unreadable marker, foreign marker, gutted+deleted-marker
   all report `malformed` with the correct problem string. An untouched install reports
   `valid`.

7. **F2 (round-1 CRITICAL) stays closed.** `runtimeAvailability` (line 1265-1288) reads
   `runtime.python` from the manifest, parses the `>=X.Y` floor, probes `python3` via the
   now-env-aware `probePython3`, and returns `{ok, requirement, detail}`. Below floor →
   `ok: false` with "python3 is X.Y, but this skill requires >=3.10". Absent → `ok: false`
   with "python3 is not available, but this skill requires >=3.10". At/above floor →
   `ok: true`. An unparseable spec → `ok: false` with "unsupported python requirement".
   `doctor`'s top-level `ok` (line 378-380) requires both `status === "valid"` AND
   `(!skill.runtime || skill.runtime.ok)`.

8. **Legacy 2.0.0 markers stay healthy.** A marker with no `markerSchema` field:
   `manifestProblems` (line 1356-1359) returns `[]` at the legacy-skip path. The marker
   check in `skillUnitStatus` passes (present, readable, `name === "parley-deck-skill"`).
   `status: "valid"`. No regression on upgrade.

9. **`--no-addons` and `--only` installs stay healthy.** `expectedAddonNames` derives the
   expected set from the core marker (line 868-871), so intentionally omitted add-ons are
   not flagged. Verified: `--no-addons` → only core skill, `ok: true`; `--only
   parley-worktrees` → core + worktrees, no `parley-bidding` unit, no false "missing".

10. **The design-addons.test.js refactor is sound.** The published-command extractor is
    parameterized on `{ binary, flag }` (line 277-278): `NODE_COMMAND` and `PYTHON_COMMAND`
    are fresh objects per use (line 281, `rx()` creates a new `RegExp` each call, avoiding
    the `/g` `lastIndex` trap). The node arm's assertions are unchanged in substance — the
    extractor logic, the `mentionsCommand` predicate, the `publishedCommands` function, and
    the refusal grammar are all shared. The Python arm reuses every line of the extraction
    machinery. The 253 pre-existing node-arm assertions pass, which is the evidence the
    refactor did not weaken them. The Python arm's `PUBLISHED_PYTHON` regex (line 1069)
    accepts `<`/`>` (placeholders, no shell) and refuses `;|&` `` ` `` `$` `\\` — the
    documented D-3 deviation, correct.

11. **The 5s timeout is present and bounded.** `spawnSync` at line 1302-1306 includes
    `timeout: 5000`. A health check that hangs on a stalled NFS mount is bounded. The
    timeout kills the child and produces `run.error` (SIGTERM), which maps to
    `probe.found: false` → `runtime.ok: false`. Fail-safe.

12. **The `pythonProbes` Map is module-level and never cleared.** This is correct for a
    single CLI invocation (one process, one `run()` call). The per-PATH key means two
    contexts in one process (e.g. a test) get independent answers. The Map grows by one
    entry per distinct PATH per process — bounded by the number of distinct environments,
    which is 1 for the real CLI and small for tests. No leak concern.

13. **No regression in install, uninstall, or paths.** `installCommand` and
    `uninstallCommand` do not call `skillUnitStatus`. `pathsCommand` returns `ok: true`
    (line 357) regardless. The `options` parameter is additive — `skillUnitStatus(unit)`
    without options still works (the `options && options.probeRuntime` guard at line 1255
    defaults to no probing, which is the pre-cycle-2 behaviour).

14. **CHANGELOG entries are accurate.** The "doctor and status report operational
    availability separately from payload validity" entry (line 34-42) describes the
    behaviour I measured. The "installed skill directory with no readable marker is now
    reported malformed" entry (line 56-59) describes the F1 fix — and honestly discloses
    the skills-CLI false positive: "This includes trees copied by a third-party skill
    installer, which do not write our marker: doctor reports them as not installed by this
    tool." That disclosure is what made the (a)/(b)/(c) question possible.

## Open questions for the implementer

1. The `status` exit-code question above (MINOR finding): is `status` intentionally
   informational-only (always exit 0), or should its `ok` reflect health the way `doctor`'s
   does now that it probes and prints the result? If intentionally informational, consider
   documenting that `status` is not a health gate and `doctor` is.

2. For option (b): when `verifyPayload` passes but the marker is absent, should `doctor`
   exit 0 (payload usable) or 1 (not installed by this tool)? I lean toward 0 with a notice,
   but the implementer should decide and document it. The `status` field needs a third
   value — what should it be called, and should `doctor --json` carry a `foreignInstalled:
   true` flag so automation can distinguish it from a native install?

3. The `pythonProbes` Map is never cleared between tests in the same process. The
   per-PATH key makes this safe in practice (each test uses a distinct PATH or the same
   one). But if a future test reuses the same PATH string with a different stub binary
   (e.g. overwriting the stub file), the cached result would be stale. Should the cache be
   cleared in a test helper, or is the per-PATH key sufficient forever? Not blocking —
   flagging for awareness.
