---
agent: kimi-1
idea: addon-manifest-coverage
review-round: 1
date: 2026-08-01
reviewed-commit: 205416d
---

## Summary

Reviewed against FINAL.md as amended by Amendment 1. Everything below was executed, not read,
unless a finding says otherwise. Method: `git archive` of both commits extracted to /tmp (the
repo itself untouched), head suite run with the repo's own `node_modules` symlinked in, the new
and reworked test files copied into the base tree and run there, plus direct driver scripts
calling `lib/installer.js` for the shape probes.

Reproduced and confirmed:

- `npm test` at 205416d: 378/378 node, 54/54 python3, `--check` green on all six skills. The
  IMPLEMENTATION.md numbers are accurate.
- `test/manifest-coverage.test.js` at 23a9856: the four fix-proving tests fail, the three
  survival guards pass. "health does not confer ownership" also fails at base — deviation 2's
  re-classification of FINAL.md item 7 is correct and I reproduce its premise: the test cannot
  establish a healthy-unmanaged fleet on the base commit (only one unit in six could be
  healthy-unmanaged there), so it is fix-dependent. Filed correctly in the shipped file.
- Deviation 1's scoping rule (`state.present ? unit.packageRoot : null`) holds under attack.
  Marker present + root extras deleted afterwards → `malformed` naming `plugin.json` and
  `gemini-extension.json` (fail-safe, correct: our marker claims the copy-plan shape). Marker
  absent + native shape → `malformed` with "no marker", which is exactly deferred follow-up 3,
  and the record describes it accurately. I found no path where our own installer produces a
  marker-present core without the root extras (every target kind stages the same
  `PAYLOAD_ENTRIES`; the generic `--dest` target too).
- F4: a 2.1.0-shaped install (`markerSchema: 2, manifest: false`, built from a manifest-stripped
  package copy) goes `malformed` under the head build with the message naming the repair, and
  re-running install **without** `--force` repairs it to `valid` — the CHANGELOG's instruction
  ("re-run your usual install command") is accurate. F4 + unselected (a later `--only` run
  excluded the stale unit) produces both problem strings coherently. No unrepairable legitimate
  install found.
- FINAL.md verification item 5: one appended byte in a payload → `--check` exit 1
  ("STALE parley-addon.json"), and install preflight writes nothing (dest absent, ok:false).
- `skills/parley-bidding/` is byte-untouched between the commits; its aggregate still
  `sha256:7854adf1…`. F6 docs updated; the CHANGELOG keeps the 2.1.0 "Known residual" paragraph
  marked as resolved rather than rewriting history, which is the right call.
- The four reworked existing tests preserve their discrimination. I read each against its base
  version: the `packageWithoutManifest()` fixture constructs the same branch conditions the
  repo previously supplied by accident (source-manifest-absent), and the two schema-exemption
  tests still pin the exemption to "source ships no manifest". `design-addons.test.js` still
  budgets exactly four doctrine files with the manifest excluded and asserted separately.
  `installer.test.js`'s missing-list test uses a present-marker bare tree, so it genuinely
  exercises the derived list — strengthened, not relaxed.

One MAJOR: the new `corePayloadFiles` walk fails open, and I reproduce a gutted managed core
reporting `valid` with `doctor` exit 0 under a damaged CLI package — the exact false-green
class this idea ships to kill. One MINOR on unscoped foreign-install coverage (gemini /
antigravity). Four NITs. No CRITICAL.

## Findings

### [CRITICAL] None.

### [MAJOR] `corePayloadFiles` fails open: a damaged CLI package silently empties the core's required-file list and turns a gutted managed core green

**What is wrong.** `corePayloadFiles` (lib/installer.js:2251) derives the core's required-file
list by walking `context.packageRoot` at validation time. Every failure in that walk is
swallowed: `walk` catches `statSync` errors and contributes nothing (:2255-2259), and
`listVisibleEntries` catches `readdirSync` errors and returns `[]` (:2568-2574). For every
non-antigravity target kind the derived list can therefore degrade to **completely empty**, and
`validateInstalledPayload` then computes `missing = []`, `problems = []`. For a managed core
(marker present) there is no second net: `manifestProblems` is add-on-only (:2301), and
`unmanagedButVerified` is not consulted when a marker is present.

**Reproduced** (driver script against /tmp copy of 205416d, codex target):

- Native install, then the installed core gutted to *only* `.parley-deck-skill-install.json`;
  `doctor` run with a `packageRoot` reduced to `lib/` + `package.json` (skills tree and root
  extras gone — a partially corrupted global/npm install of the CLI itself):
  `status: "valid"`, `managed: true`, `missing: []`, `problems: []`, **`doctor` ok: true**.
  With the intact package the same gutted tree is `malformed` with 8 missing files. So the
  health answer for an identical installed tree depends on whether the CLI's own package
  happens to be readable, and the degraded answer is the green one.
- Partial damage is enough: package copy with only `skills/parley-deck/references/` removed,
  installed core with `references/COOPERATION.md` deleted → `valid`, `missing: []`. The
  unreadable subtree silently drops out of the required list and deletions inside it become
  invisible.
- Nothing ever notices the degradation: with the same broken package and an *intact* fleet,
  all six units report `valid` and `doctor` exits 0. `validatePayload(packageRoot)` — the only
  package-integrity check — gates `installCommand` alone (lib/installer.js:607); `doctor` never
  runs it. Install from such a package fails loudly (`firstCopyObstacle` refuses missing
  sources), so the two paths derive opposite answers from the same damage: install says the
  package is broken, `doctor` says every previously installed tree is healthy.

**Why it matters.** This reopens D2/D4 — the false-green class F5 exists to close — through a
single realistic precondition (npm corruption, partial extraction, AV quarantine, or a future
packaging mistake trimming `skills/` from the published tarball). It fails silently, and it
fails in the direction this whole change is built to prevent. The implementer flagged this
exact spot for review; the suspicion is confirmed by execution. F5's design note says "same
constant the staging loop copies from, so the two cannot disagree" — they *can* disagree,
because the staging loop's read (`firstCopyObstacle`/`copyRecursive`) fails loud while the
validation walk fails silent.

**Suggested fix.** Fail closed. Minimal form: when `packageRoot` is supplied, any
`PAYLOAD_ENTRIES` source that is missing, unreadable, or not the expected kind must contribute
a validation problem (e.g. `cannot inspect the packaged payload to derive the core's required
files: <entry>`) rather than contributing nothing — have `corePayloadFiles` return
`{ files, incomplete }` and `validateInstalledPayload` turn `incomplete` into a problem.
Stronger form: generate the derived list at build time into a packaged artifact (verified by
`--check`, like the manifests) and read that at validation time, so `doctor`'s answer never
depends on a runtime walk of the package. Either way, add a regression: gutted managed core +
unreadable package source subtree → not `valid`.

### [MINOR] The D1 fix does not cover the gemini or antigravity targets: a byte-perfect foreign copy of the core is still `malformed` there

**What is wrong.** With no marker present, `skillUnitStatus` passes `packageRoot: null`, so the
historical per-kind fallback lists apply — and those lists demand files a verbatim foreign copy
of `skills/parley-deck/` cannot contain (`gemini-extension.json` for gemini;
`skills/SKILL.md` and `plugin.json` for antigravity). `unmanagedButVerified` succeeding does
not help: `payloadOk` requires `missing.length === 0` independently (lib/installer.js:2092), so
the floor list convicts before the manifest proof can acquit.

**Reproduced** at 205416d: verbatim copy of all six `skills/*` directories into
`~/.gemini/extensions` → five add-ons `valid-unmanaged`, core `malformed`
(`missing: ["gemini-extension.json"]`), `doctor` exit 1. Same into
`~/.gemini/config/plugins` (agy) → core `malformed` (`missing: ["skills/SKILL.md",
"plugin.json"]`), exit 1. The codex-shaped copy from the shipped regression test passes
exactly because the codex fallback list happens to name only files inside
`skills/parley-deck/`.

**Why it matters.** FINAL.md verification item 1 ("foreign verbatim copy of all six → six
`valid-unmanaged`, `doctor` exit 0") names no target, and the CHANGELOG's claim is unqualified
("Such a tree now reports `valid-unmanaged`"). The shipped regression covers codex only. For
gemini the red is arguably *truthful* — an extension without `gemini-extension.json` genuinely
cannot load — but then the claim needs scoping, because as written the docs promise a result
two targets do not deliver. (Whether a real foreign installer ever writes to those directories
is an open question below; it changes whether this is a live false red or a doc bug.)

**Suggested fix.** Pick one and say so: (a) scope the claim — document in the CHANGELOG and in
the test header that the `valid-unmanaged` proof covers the codex-shaped skills layout, and
that gemini/antigravity extension layouts keep their runtime-manifest requirement on the
unmanaged path; or (b) make a successful `unmanagedButVerified` supersede the fallback floor
list for core kinds (safe for B3: a gutted unmarked tree never verifies, so it stays
`malformed`). Option (a) is the smaller change and consistent with "no per-target layout
manifests".

### [NIT] The F4 fix-proving test fails at 23a9856 only at fixture construction — fix-dependent by the implementer's own deviation-2 standard

**What is wrong.** `bidding-addon.test.js`'s "a manifest:false marker is unhealthy once the
skill ships a manifest" is the F4 fix-proving regression (FINAL.md item 3), and
IMPLEMENTATION.md states "Every fix-proving regression confirmed failing at `23a9856`".
Reproduced: copied into the base tree, the test fails at
`packageWithoutManifest`'s assert ("parley-worktrees was expected to ship a manifest to
remove", test/bidding-addon-head.test.js:67) — it never reaches the behavior it proves. At
`23a9856` the test cannot construct its own precondition (a packaged manifest for that skill
exists only after the fix), which is precisely the property deviation 2 used to re-classify
FINAL.md item 7: "fix-dependent, not invariant… cannot establish its own precondition on the
base commit". The same standard was not applied to the implementer's own F4 test, and the
"confirmed failing" claim is vacuous for it — the base run proves nothing about the
`declared === false && sourceHasManifest` branch, because it never executes.

**Why it matters.** The behavior itself is correct (probe E above: red with the right message,
repair without `--force`, verified at head), so this is about the honesty of the classification
ledger this deck runs on — the same overclaim family as Amendment 1.1, caught this time in the
implementer's verification summary rather than the drafter's.

**Suggested fix.** Add the deviation-2 treatment to the test header (state that the fixture is
fix-dependent and the base failure is at construction), or write the base-runnable variant:
install from the base package (`manifest: false` marker), then `doctor` with a packageRoot in
which that skill *has* been given a manifest — that asserts `malformed` at head and executes
the branch-less `return []` at base, so its base failure actually discriminates.

### [NIT] `listAddons()` enumerates directories without requiring `SKILL.md`; F2's binding text says "every `skills/*/SKILL.md` directory"

**What is wrong.** `scripts/build-addon-manifest.js`'s `listAddons()` accepts any directory
under `skills/` (`statSync(...).isDirectory()`), while the installer's `discoverAddons`
requires `skills/<name>/SKILL.md` (lib/installer.js:871) and FINAL.md F2 scopes the enumeration
to "every `skills/*/SKILL.md` directory". Reproduced: a stray `skills/scratch-notes/` directory
containing only `notes.md` makes `--check` exit 1 with `scratch-notes: MISSING
parley-addon.json` — a build failure demanding integrity metadata for a non-skill, with a
message that would send a contributor looking for the wrong fix. In generate mode the tool
would happily *create* a manifest there.

**Why it matters.** Small, but the two components now disagree on what a "packaged skill" is,
and the mandatory-coverage gate is exactly where a confusing failure costs the most time (a
stray backup/notes directory under `skills/` breaks `npm test` and `prepack`).

**Suggested fix.** Filter `listAddons()` on `fs.existsSync(path.join(dir, "SKILL.md"))`, the
same predicate `discoverAddons` uses, and note in the header comment that the core is included
deliberately (it has a SKILL.md, so the filter does not re-exclude it).

### [NIT] Amendment 1.2's "asserted per target" guard covers two of the staging shapes

**What is wrong.** The survival guard "the natively installed core does not carry
parley-addon.json" asserts `codex` and `gemini` only. Amendment 1.2 ratified the guard as
"asserted per target", and the one staging shape that *differs* — antigravity, which
additionally stages `skills/SKILL.md` (lib/installer.js:1857-1862) — is uncovered, as is the
generic `--dest` target. The guard's substance (keeping the manifest out of `PAYLOAD_ENTRIES`)
is indeed target-independent, so the codex assertion catches the actual mutation it exists
for; this is a gap against the amendment's letter, not against the risk.

**Suggested fix.** Extend the loop to `["codex", "gemini", "agy"]` (and a generic `--dest`
install if cheap), or record in the test header that per-target was narrowed to per-staging-
shape with the reasoning, so the amendment and the test agree.

### [NIT] The fix-proving gut test hard-codes `parley-worktrees`, for which "gutted to SKILL.md" removes no payload

**What is wrong.** In `manifest-coverage.test.js`, "a natively installed add-on gutted to
SKILL.md is malformed even with its marker intact" guts a fixed four-skill list including
`parley-worktrees` — a single-file payload. For that skill `gutTo(..., ["SKILL.md", MARKER])`
removes only `parley-addon.json`, so the `malformed` verdict comes from the
manifest-deleted-versus-marker branch (`declared` object, `present` false), not from payload
gutting; if payload verification broke but manifest-deletion detection survived, the
worktrees assertion would still pass. The same file's survival guard handles this correctly by
deriving subjects from `listPayloadFiles(...).length > 1` (and IMPLEMENTATION.md shows the
implementer knew), so the fix-proving test is the inconsistent one.

**Suggested fix.** Derive the gut list the same way the survival guard does (payload size > 1),
so the test name describes what every subject actually exercises and a future single-file skill
does not silently join with a vacuous assertion.

## Open questions

1. **Version.** `package.json` and `package-lock.json` still say 2.1.0 at 205416d while the
   CHANGELOG in the same commit declares `## 2.2.0 — 2026-08-01`. The 2.1.0 release was bumped
   in its own release commit (`23a9856`), so this presumably follows the same pattern — confirm
   the release commit is planned before publish, since `writeMarker` stamps
   `PACKAGE_JSON.version` into every new marker and a publish from this tree would write
   "2.1.0" into 2.2.0-era markers.
2. **Foreign installers and the gemini/antigravity layouts.** Does the universal `skills` CLI
   (or any foreign installer in scope for D1) actually write to `~/.gemini/extensions` or
   `~/.gemini/config/plugins`? If yes, the MINOR above is a live false red for those users and
   leans toward fix option (b); if no, it is documentation-only and option (a) closes it.
3. **Portable (`pkg`) builds.** `pkg.assets` includes `skills/**/*`, so the five new manifests
   and the walked payload trees should be present in the snapshot filesystem — but I did not
   build a portable binary to confirm `corePayloadFiles`'s runtime walk and
   `unmanagedButVerified`'s source reads work from inside `process.pkg`. Worth one
   `build:portable:current` smoke run before the release commit.
