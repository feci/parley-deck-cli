---
idea: integrate-parley-bidding-addon
review-round: 16
agent: kimi-1
date: 2026-07-31
reviewed-commit: a49d68f
---

## Verdict
ACCEPT

## Outstanding findings — closed or not

All three round-15 findings are closed, and I re-measured each rather than trusting the table:

- **Containment MAJOR — closed.** `resolvedDestination` + `overlaps` in `aliasedDestinations`
  (lib/installer.js:1374-1459) refuse equality, containment, and contained-by on the resolved
  physical path, in both nesting directions, for install and uninstall alike. I ran the
  a49d68f regression against a `git archive` copy of `26478e9`: it **fails** there with
  `ok: true` (the false success codex-1 measured) and **passes** at a49d68f, both directions.
- **Uninstall single result path MINOR — closed.** The real-only fleet preflight is deleted;
  both modes build results from `removeFleetAtomically`. The new regression fails at `26478e9`
  with exactly the old flattening in the diff (`parley-tracker:true:missing` dry vs
  `parley-tracker:false:skipped` real) and passes at a49d68f.
- **Fourth manifest reader MINOR — closed.** `scripts/run-python-tests.js` now reads the floor
  through `readManifest` and fails on any refusal. Verified both directions of the floor
  itself: under python3 3.14 the leg runs 54/54; under `/usr/bin/python3` 3.9.6 (PATH shim) it
  exits 1 with `python3 is 3.9, but the add-on declares >=3.10` — the message proving the
  floor came through the shared parser.
- **My cycle-18 non-discriminating test — closed.** The rewritten case-only regression puts
  two spellings in one plan and **fails at `d7ab1c3`** (`ok: true`, both written) as it should
  have originally; it passes at a49d68f.

Full suite at a49d68f, measured here: **355 node tests, 0 fail**; python leg **54/54** on
3.14; manifest check ok, 47 files, aggregate `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`.

## New findings

Two, both follow-ups. Neither blocks.

**MINOR — the manifest's `runtime.python` spec has two consumers that disagree, and the
release gate is the permissive one.** The grammar `/^>=\s*(\d+)\.(\d+)$/` is enforced
fail-closed by `runtimeAvailability` (lib/installer.js:2046-2049: an unparsable spec is
`ok:false`, "unsupported python requirement"), but `declaredPythonFloor` in
scripts/run-python-tests.js:53-61 returns `null` on the same spec — floor silently
unenforced, suite green. The malformed value can only enter through the sanctioned
`--runtime-python` flag, after which `build-addon-manifest.js --check` byte-pins it (the
runtime block sits outside the payload aggregate), so `npm test` ships it green while every
installed tree fails `doctor`. Measured end-to-end on a /tmp copy at a49d68f:
`--runtime-python "3.10"` → `--check` ok (aggregate unchanged) → python leg 54/54 on 3.14 →
install `ok:true` → `doctor` `ok:false` with
`runtime: {"ok":false,"requirement":"3.10","detail":"unsupported python requirement \"3.10\""}`.
This is the same species as the round-15 MINOR cycle 19 closed — a read failure swallowed as
"no declared floor" — one arm further along: the unreadable-manifest arm now fails, the
readable-but-unparsable-spec arm still swallows. It does not block 2.1.0: the shipped manifest
declares a well-formed `>=3.10`, `--check` pins its bytes, and triggering this needs a
maintainer typo through the build flag. Recorded follow-up; the cheapest fix is to fail in
`declaredPythonFloor` when `runtime.python` is present but unparsable (mirroring
`runtimeAvailability`), or to validate the runtime block in `readManifest` itself.

**NIT — dead code and stale comments from the cycle-19 deletion.** `preflightUninstallUnit`
(lib/installer.js:687-712) has no callers: its only call site was the deleted preflight, and
`removeFleetAtomically` inlines the logic. Its doc comment and the `pathEntryExists` comment
(lib/installer.js:2261-2262) still describe `uninstallSkillUnit`, a function that no longer
exists. No behavioral impact.

Also probed and dismissed, stated so the next round does not have to redo it:

- *Planning-to-commit relation change (the round-16 prompt's first question).* I enumerated
  every mutation site in lib/installer.js: `mkdirSync`, `writeFileSync`, `renameSync`,
  `rmSync` only — no symlink creation anywhere, and `copyRecursive`/`firstCopyObstacle`
  refuse symlinks in sources. Path resolution changes only via symlink create/remove or
  ancestor replacement; staging/commit do neither (renames stay within the same parent and
  never cross units, since containment is refused at planning). A tail that "becomes existing"
  is a real directory created along the already-computed tail — the resolved string is
  unchanged. A destination whose resolution depended on a file another unit writes is exactly
  the containment case now refused. What remains is cross-process change, which is the
  recorded single-writer limit in CHANGELOG "Known limits" (round-14 ruling — not
  re-litigated).
- *Case-twin add-on* (`skills/Parley-Deck` beside `parley-deck`). Not constructible on a
  case-insensitive volume: the two spellings are one directory — I accidentally proved this by
  overwriting and then deleting `skills/parley-deck` through its twin spelling in a /tmp
  copy. The destination-level equivalent (two homes differing only in case, one plan) is the
  round-14/15 regression, verified discriminating above; `physicalKey` lowercases tails.
- *`overlaps` root-prefix gap.* At string level `overlaps("/", "/parley-bidding")` is false —
  root containment is not seen. Reachable only via `--target generic --dest /` (generic never
  shares a plan with other targets). Measured as non-root: all six units blocked at preflight
  (EROFS), zero writes, `ok:false`; as root the commit would fail attempting to rename `/`.
  Fails closed both ways; not a finding.
- *Uninstall dry/real beyond the sampled dispositions (the prompt's second question).* I
  measured a four-scenario matrix, not only the regression's sample: a mixed 14-target fleet
  (foreign-marker blocked + deleted missing + clean units) — per-unit shapes **identical**
  between dry and real; a clean 84-unit fleet — per-unit `ok` identical, action differing only
  by the deliberate `remove`→`removed` tense, mirroring install's `install`→`installed`
  (scenario 3); a corrupt recorded selection (`addons: ["../../escape"]`) — identical, core
  blocked in both modes. Every disposition is assigned before the mode branch, so agreement is
  structural; the only dry-can't-predict residue (quarantine rename failure) is the documented
  "removability is deliberately not predicted" design. Uninstall's per-action JSON lacks the
  `dryRun: true` flag install carries on the action object — pre-existing, and the flag is
  present per skill and at top level; noted, not filed.
- *Fifth manifest reader / validator bypass (the prompt's third question).* None. Every
  `readFileSync` in lib/bin/scripts accounted for: the parser itself, the staged-copy gemini
  rewrite, payload copying, `readMarkerState` (all marker reads route here), `sha256File`/
  `readJsonFile` (informational only), and `build-addon-manifest.js`'s byte-compare guarded by
  `hasManifest` + `verifyPayload`. No Python code in the add-on reads `parley-addon.json`.
  Marker-field consumers each validate what they trust (name/skill compared, `addons`
  name-validated + confined + authorized, `manifest` structurally checked, `markerSchema`
  range-checked). The one stored value with a permissive consumer is the MINOR above.

## Release judgement

Releasable as 2.1.0. The round-15 MAJOR is genuinely closed — containment is computed where
equality was, the proof now covers the only in-process materialization primitive, and both
regressions discriminate against the commits they name. The two new findings are a
release-gate robustness gap that requires an authoring error to reach (MINOR) and dead code
(NIT); both are recorded follow-ups, not blockers.

## What I verified

- Working tree clean at `a49d68f`; no mutations to the repo — all cross-commit comparison done
  on `git archive` copies under /tmp, all experiments against /tmp package copies or tmp HOMEs.
- `npm test` at a49d68f: 355/355 node, 54/54 python (3.14), manifest check ok (47 files,
  aggregate unchanged since `714712f`, matching the record).
- Discrimination (a49d68f test file against archived old trees): nested-destination regression
  fails at `26478e9`; uninstall dry/real regression fails at `26478e9`; case-only regression
  fails at `d7ab1c3`. All three pass at a49d68f.
- Q1: full mutation-site enumeration (no symlink creation); case-twin addon non-constructible
  (volume-level identity of the two spellings); `--dest /` refuses cleanly as non-root.
- Q2: four-scenario dry/real matrix — dispositions identical except the deliberate tense pair.
- Q3: all `readFileSync` sites audited; the `runtime.python` consumer disagreement reproduced
  end-to-end (`--check` ok → suite green → install ok → `doctor` red).
- Python leg under `/usr/bin/python3` 3.9.6 refuses with exit 1 naming the declared floor.
- CHANGELOG "Known limits" carries the single-writer ruling; left untouched per the round-14
  settlement. `IMPLEMENTATION.md` and `review/round-09/VOID.md` present in the idea records.
