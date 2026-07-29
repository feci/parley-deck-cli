---
idea: integrate-parley-bidding-addon
round: 01
agent: kimi-1
date: 2026-07-29
---

# kimi-1 — round 01 — packaging & distribution mechanics

Host facts used below: macOS, node v26.5.0, npm 11.17.0, python3 3.14.6
(`/opt/homebrew/bin/python3`), no `pytest` module (TESTED: import fails), no
`skills` CLI on PATH (TESTED: `which` finds nothing). Target repo on branch
`readme-skill-catalogue`, clean, HEAD `c3aa392`.

## 1. The integration plan

Everything below was read on disk today; line numbers are from
`readme-skill-catalogue` HEAD.

### 1.1 What lands where

One new tree: `skills/parley-bidding/` — a **filtered** copy of the 48-file,
246 KB source (`du -sk` TESTED: 246). The filter excludes `__pycache__/`,
`*.py[cod]`, `.DS_Store`, `.venv/` — the source's own `.gitignore` rules —
plus the `.gitignore` file itself (F4, §3). The copy must be made with an
exclude-capable tool (`rsync --exclude`, or `tar` with `--exclude`), **never
bare `cp -R`**: at my first listing today the source contained
`scripts/__pycache__/` with seven `*.cpython-314.pyc` files — contradicting the
brief's "no `__pycache__`" — and minutes later the directory was gone, removed
by an actor that was not me (I ran only `find`, `grep`, `ls`, and Python with
`PYTHONDONTWRITEBYTECODE=1`; nothing destructive). Both observations are true;
the source is a shared, mutable volume. The integration must snapshot source
file hashes at copy time and record them in `IMPLEMENTATION.md`, so the
source-vs-integrated diff (check #10) compares against the state actually
copied, not against a later state of a moving tree.

No other repository location changes. The add-on is discovered dynamically:
`discoverAddons` (`lib/installer.js:778-794`) scans `ADDONS_DIR = "skills"`
(`:145`), skips only `CORE_SKILL_NAME = "parley-deck"` (`:787`), and accepts
any directory with `SKILL.md` (`:146`, `:790`). `usage()` builds its
`Available add-ons:` line from discovery (`:325-326`, `:341`), and the
`--only` error text likewise (`:797-807`, throw at `:805`) — the brief's
"hardcoded" sites at 341/805 are in fact generated and self-updating (TESTED
by reading; also TESTED end-to-end: the v1.5.0 portable binary's dry-run
install listed all four current add-ons from the pkg snapshot).

### 1.2 `lib/installer.js` — no functional change, one comment

- `:785-786` — comment "…a generic skill installer sees all five / as
  siblings" → "all six". The only edit in this file.
- Nothing else: discovery, `--only`, marker recording (`writeMarker`
  `:1056-1075`), and install/uninstall/status/doctor are name-agnostic.

One design change I *do* propose lives here, but it is generic, not
bidding-specific — see F6 (§3): `validateInstalledPayload` (`:1129-1148`)
gains an optional declarative required-files manifest for `kind === "addon"`
(currently hardcoded to `["SKILL.md"]` at `:1131-1132`).

### 1.3 `package.json` — no change required to ship

- `files` (`:31-40`) already allowlists `skills/` whole — the new tree ships
  in the npm tarball with zero edits (mechanism TESTED today via
  `npm pack --dry-run` baseline: 153 files, 952.6 kB unpacked).
- `pkg.assets` (`:41-52`) already contains `"skills/**/*"` — the portable
  binaries embed the new tree with zero edits (mechanism TESTED: today's
  `dist/parley-deck-skill-v1.5.0-macos-arm64` installed all 154 files of the
  five current skills into a temp `$HOME` and `doctor` reported all five
  `valid`).
- The only candidate edit is under F2 (§3): none to `"scripts"` — the Python
  leg is a `node --test` file, so `"test": "node --test"` (`:60`) and the
  RELEASING.md preflight stay untouched.

### 1.4 Tests

- `test/installer.test.js` — four assertion sites hardcode the four-add-on
  set and must gain `parley-bidding` (alphabetical first under
  `localeCompare`): `:408` (`discoverAddons` deepEqual), `:422-428` (default
  install skills list), `:432-436` (per-add-on `existsSync`), `:541` (doctor
  skills list), `:563` (core marker `addons` list). The `--no-addons` /
  `--only` assertions (`:450`, `:452-455`, `:464-467`, `:478-482`, `:590`,
  `:602`, `:611`, `:620`, `:638`) stay as-is; consider adding
  `parley-bidding` to the absence checks at `:452-455` for symmetry.
- `test/design-addons.test.js` — no count edits needed (`published.size >= 2`
  at `:352` is unaffected: the source tree contains no `node --test` string —
  TESTED by grep — so the published-command guard `:334-393` sees the same
  two commands). Two *additions*, both generic:
  1. a cache scan over `skills/` that fails on `__pycache__`, `*.pyc`,
     `.DS_Store` (prevents the exact dirty-tree state I observed this
     morning from ever shipping);
  2. the F5 python branch of the published-command guard (§3).
- New `test/bidding-addon.test.js` — the F2 Python leg (§3). `node --test`
  discovers `*.test.js` repo-wide, so it joins `npm test` automatically
  (TESTED baseline: 253 tests pass today, including the shipped add-ons'
  own `*.test.js` files, which is precedent that shipping tests inside an
  add-on tree is normal here: `skills/parley-tracker/bin/claim.test.js`,
  `skills/parley-design-check/test/`).

### 1.5 README.md

- `:3` hook comment, `:14-15` ("five skills … four add-ons"), `:19-21`
  ("places five skills … the other four … All five install by default"),
  `:109` ("installs all five skills") → six / five add-ons.
- `:23-27` bullet list gains a `parley-bidding` bullet; a new `###` section
  after `parley-worktrees` (`:90-103`) describes it — and that section is
  where two load-bearing sentences go (§4, items B and F): the E3b-vs-default-
  consent precedence, and the new python3 prerequisite.
- `:248-264` repository-layout tree gains `parley-bidding/`. The brief did
  not list this block; it rots otherwise.
- `:121-122` ("fourteen named runtimes") unchanged.

### 1.6 IMPLEMENTATION.md (new, at release time)

Provenance (source path + content hashes at copy time), F1–F7 decisions, the
run verification table, the intentional-difference list (§2 and §4 item E),
and the release implications (§4 item A).

## 2. The rename: `software-bidding` → `parley-bidding`

The brief's "8 files" is correct; its per-line list is **incomplete by three
lines**. Verified grep for `software-bidding` over the source returns exactly
nine hits in eight files, and two more lines are name-derived human labels
that must move with the rename or they become stale references of exactly the
class this repo spent 12 review rounds eliminating. All edits, by file:

1. `SKILL.md:2` — `name: software-bidding` → `name: parley-bidding`.
2. `SKILL.md:6` — H1 `# Software Bidding` → `# Parley Bidding`
   (**missed by the brief**; the H1 is the display twin of `display_name`).
3. `agents/openai.yaml:2` — `display_name: "Software Bidding"` →
   `"Parley Bidding"` (**missed by the brief**; otherwise the Codex UI shows
   a name no installed skill has).
4. `agents/openai.yaml:4` — `Use $software-bidding to …` → `$parley-bidding`.
5. `scripts/common.py:2` — module docstring `for the software-bidding skill`
   → `parley-bidding`.
6. `scripts/tests/test_skill_structure.py:16` — regex
   `(?m)^name: software-bidding$` → `parley-bidding`.
7. `scripts/tests/test_skill_structure.py:26` — `assertIn('display_name:
   "Software Bidding"', metadata)` → `"Parley Bidding"` (**missed by the
   brief**; required by edit 3, and this is the test the brief warned will
   fail — it fails on two assertions, not one).
8. `scripts/tests/test_skill_structure.py:27` — `assertIn("$software-bidding",
   metadata)` → `"$parley-bidding"`.
9. `assets/schemas/bid-state.schema.json:3` — `$id` path segment (see F1).
10. `assets/schemas/jurisdiction-profile.schema.json:3` — same.
11. `assets/schemas/platform-adapter.schema.json:3` — same.
12. `assets/schemas/procedure-profile.schema.json:3` — same.

Not touched: `SKILL.md:3` description and `openai.yaml:3` short_description
contain no skill name; `references/` contain no occurrence (TESTED: the grep
above covers the whole tree).

### F1 position — the schema `$id`s: **rename the path segment, keep the
`example.invalid` host**

New form: `https://example.invalid/parley-bidding/<file>.schema.json`.

Evidence that this is safe, all TESTED today by grep:

- No `$ref` exists in any of the four schemas — nothing resolves one schema
  against another.
- No script reads `$id` or `example.invalid` — the only `example.invalid`
  hits in `scripts/` are synthetic fixture accounts/URLs in the Python tests,
  unrelated to schema identity.
- No script writes `$schema`/`$id` into generated bid workspaces — so even
  workspaces created by the standalone skill carry no reference to the old
  identity. The rename breaks no existing consumer, inside or outside the
  package.

Arguments against the alternatives. **Keep `…/software-bidding/…`:** a `$id`
is an identity, but it is also the only place the retired name would survive;
a schema whose identity names a skill that does not exist is precisely the
stale-path rot the `addons/` guard (design-addons.test.js `:178-204`) exists
to catch. Stability of identity has value only when something resolves
against it — here, nothing does. **Re-root under a real parley-deck URL:**
worse than both. `example.invalid` is RFC-2606-reserved and honestly
non-resolving; a real domain asserts a canonical, fetchable governance
location that does not exist and invites resolution attempts. Rename the
segment, keep the reserved host.

## 3. Positions on every fork

- **F1 — schema `$id`s:** rename the path segment to `parley-bidding`, keep
  `example.invalid`. §2 above.
- **F2 — Python toolchain:** `npm test` gains a Python leg, implemented as
  `test/bidding-addon.test.js` (a `node --test` file, so `"test"` stays
  `node --test`). It locates `python3` (PATH, then common homes), and —
  per test file, matching the facilitator's measured table — runs
  `PYTHONDONTWRITEBYTECODE=1 python3 scripts/tests/test_X.py`, asserting exit
  0, `OK`, and the per-file counts 4/20/2/3/15/3/7, total **54 reproduced,
  not asserted**. The exact-count assertion is deliberate: a future skill
  edit that adds or loses a test fails red until someone consciously updates
  the expectation — the same philosophy as the doctrine budgets in
  `test/design-addons.test.js:22-28`. If `python3` is absent, the test
  **fails with a named reason**; a skip would be the false-green this repo
  just spent 12 rounds killing. Consequences accepted: contributors need
  python3 to run the suite; CI (`release-portable.yml` runs `npm ci &&
  npm test` on `ubuntu-latest`, which ships python3 preinstalled — NOT TESTED
  from this host) exercises the leg at release. Never `pytest` (not installed
  here, TESTED) and no new dependency. Note: the brief says
  `unittest discover` "fails against that directory" — that is **stale**:
  `python3 -m unittest discover -s scripts/tests` ran 54/OK/exit 0 on Python
  3.14.6 today (TESTED). I still prescribe per-file execution: it isolates
  failures and yields the per-file counts the brief's table records.
- **F3 — version metadata:** the add-on inherits the package version, full
  stop. `writeMarker` already stamps `PACKAGE_JSON.version` into every
  installed skill's marker (`installer.js:1061`), and the skill's audit
  semantics are anchored by content hashes (`manifest.py`, raw SHA-256), not
  by version numbers — a second version field creates a drift axis with no
  consumer. Provenance lives in `IMPLEMENTATION.md` (source path + hashes +
  date). If the skill's semantics later change, that is a package
  minor/major bump like any other content change.
- **F4 — source `.gitignore`:** **drop it from the packaged tree; hoist its
  rules into the root `.gitignore`; enforce with the new cache-scan test.**
  TESTED basis: with `"files": ["skills/"]`, npm 11.17.0 both (a) refuses to
  ship a nested `skills/x/.gitignore` itself and (b) honors its rules —
  a simulation with a nested `.gitignore` containing `__pycache__/`/`*.pyc`
  packed 3 files, silently excluding a planted `.pyc`; the identical tree
  without it packed 4, `.pyc` included. So keeping the nested file makes
  `npm pack` look clean **even when the repo is dirty**, while the portable
  build (pkg embeds from disk, no ignore semantics) and repo-checkout
  installs (`copyRecursive` has no filters — `installer.js:1077-1091`)
  would still carry the filth: an install-mode divergence and a false-green
  in one. Dropping it makes all three channels see the same tree, the root
  `.gitignore` (currently only `node_modules/ dist/ *.tgz .DS_Store` — does
  not cover Python caches) keeps them out of git, and the test makes a dirty
  tree fail red instead of ship.
- **F5 — published-command guard:** extend it, but **never execute** the
  published python commands. The commands in `SKILL.md:95-116` carry
  `<placeholder>` arguments — they are templates; executing templates either
  fails spuriously or, worse, partially runs tooling. Add a python branch to
  the guard that extracts `python3\s+scripts/\S+\.py` references from shipped
  markdown and asserts each referenced file **exists and compiles** (the
  MAJOR-1 dangling-reference class, closed), plus a `>0` count assertion so
  silent disappearance fails. Real execution coverage comes from the F2 leg,
  which is stronger than running doc snippets. Accepting the gap in writing
  is not acceptable — the gap is closable for ~40 lines of test.
- **F6 — installer validation:** yes, assert more than `SKILL.md` — but
  generically, not bidding-specifically. `validateInstalledPayload`
  (`installer.js:1129-1148`) gains support for an optional declarative
  manifest inside the add-on (e.g. `.parley-deck-addon-files.json`: a list of
  required relative paths); when present, install-time validation and
  `doctor`/`status` require those files; when absent, today's `["SKILL.md"]`
  behavior stands. No add-on-specific strings in the installer. The manifest
  itself is pinned by a package test that recomputes it from the tree, so it
  cannot rot into a false-green. Justification: `copyPayloadAtomically`
  cannot lose files silently (a missing source throws in `copyRecursive`),
  so the real enemy is **post-install gutting** — a runtime cleaner or user
  deletes `scripts/`, and `doctor` today reports `valid`. For this add-on the
  scripts *are* the safety guarantees (freeze, manifest, no-retry
  enforcement); an install that lost them must read as `malformed`, not
  `valid`.
- **F7 — sequencing (B1):** accept "design now, implement after" without
  amendment. The intersecting file set is exact (`skills/`,
  `lib/installer.js`, `package.json`, both test files, `README.md`), and the
  collision rule is this project's own ratified discipline. Two additions
  from today's evidence: (a) because the source tree proved mutable mid-round
  (the `__pycache__` episode, §1.1), the implementation must copy from a
  hash-snapshotted source state; (b) rounds 1–2 should pre-agree the exact
  test baselines (253 Node tests + 54 Python tests) so post-merge
  verification is a comparison, not a rediscovery.

## 4. What could go wrong that the brief has not anticipated

Ordered by how much they worry me. The brief preserves the skill's guarantees
as *text*; these are the ways the **act of packaging** weakens them in fact.

**A. The threat model changes because the distribution unit changes.** Today
`software-bidding` is a personal skill its owner deliberately installed.
After integration it rides the **default install**: every `npx -y
parley-deck-skill@latest install` puts a procurement-portal workflow into
up to 14 runtimes, including ones the user rarely audits. Worse, the upgrade
path: `expectedAddonNames` (`installer.js:846-863`) uses the package default
for a flag-less `install`, so the next routine `install --force` update of
every *existing* user silently adds `parley-bidding` to every runtime they
have (their old marker pins the read/uninstall view, not the next install).
The HITL gates survive — they bind the agent reading the skill — but
*availability* of portal-adjacent workflows expands by default to users who
never asked for a bidding tool. Not a blocker; a hard documentation duty:
the README section and release notes must say plainly that the sixth skill
operates procurement portals under HITL gates, and `IMPLEMENTATION.md` must
record this expansion as a release implication.

**B. Two consent defaults collide, and packaging creates the collision.** The
bidding skill's E3b gate (`SKILL.md:63`, `references/parley-integration.md:5-16`)
requires showing roster/providers/data classes/exact packet and obtaining
**tender-scoped** approval before any disclosure to model backends. The deck's
local agent contract (README `:226-236`) defaults to **YES** "for sending the
task brief plus necessary repository context to external CLI backends."
Standalone, these never meet. Packaged together, `parley-integration.md:3`
explicitly routes bid challenges through "the active project's
`parley-deck/COOPERATION.md`" — so a facilitator running a Parley challenge
under the deck's default consent would disclose tender content to external
backends *without* the E3b gate as the bidding skill defines it, and both
documents would claim to have been followed. The skill's own
stricter-gate rule ("A platform adapter may impose stricter gates but never
weaker ones", `SKILL.md:70`) covers adapters, not host protocols. The brief
forbids touching `COOPERATION.md` (correctly — that is a separate
meta-protocol idea), so the fix lands in the integration's own text: one
sentence in the new README section and in `IMPLEMENTATION.md` stating that
the deck's default transport consent never satisfies E3b, and that tender
disclosure requires the skill's own tender-scoped gate. Without that
sentence, installing the two side by side manufactures a compliance gap.

**C. The deterministic guarantees are load-bearing Python; the package
neither ships nor requires an interpreter.** `engines` declares only
`node >= 18` (`package.json:53-55`). On a python-less machine the skill
still *reads* as if freeze, exact-byte manifests, no-retry enforcement and
adapter validation exist — but none of them can run, and an agent mid-
workflow will improvise around the missing tools. That is precisely the
failure the scripts exist to prevent, reintroduced by omission. The README
must name python3 as this add-on's prerequisite (no other add-on has one)
and the skill should be read as guidance-only without it. (The scripts
themselves are stdlib-only — TESTED by import scan: argparse, csv, json, sys,
re, zipfile, hashlib, os, tempfile, shutil, copy, decimal, datetime,
pathlib, typing, plus `subprocess`/`unittest` in tests. No network module
anywhere, so portal-safety of the deterministic tooling holds at the import
level.)

**D. The copy channel has no hygiene filter, and one channel's false-green
hides another's dirt.** §1.1 and F4: `copyRecursive` copies whatever is in
the repo tree, including caches; pkg embeds whatever is on disk; only npm
filters (via ignore files). Three channels, three answers, unless the repo
tree itself is kept clean and a test proves it. Related observed hazard: the
source volume is shared and was modified by a third party while I worked —
any integration step that assumes a stable source (including a bare `cp -R`
or a diff run hours after the copy) is unsound.

**E. File modes do not survive any install mode.** Five source scripts carry
`+x` (`adapter_validate`, `common`, `init_bid_workspace`, `manifest`,
`release_lint`; `bid_state.py` and `completeness_lint.py` do not — the source
is itself inconsistent), and all carry `#!/usr/bin/env python3` shebangs.
`copyRecursive` (`installer.js:1089-1090`) writes default modes; npm
normalizes shipped modes; pkg's snapshot likewise. So every installed copy
loses the exec bit. Functionally benign — every documented invocation goes
through `python3 scripts/…` — but an agent trusting the shebang and trying
`./scripts/manifest.py` gets permission denied. Record as an intentional
difference in check #10; do not "fix" by adding mode-preservation to the
installer (scope creep for zero functional gain).

**F. Smaller, confirmed-clean items (TESTED, listed so reviewers don't
re-chase them):** no symlinks in the source (0 found — so the symlink
refusal at `installer.js:1079-1081` never trips today); no `addons/` or
`node --test` strings anywhere in the source (the two repo guards are
unaffected); no `$ref`/runtime `$id` coupling (§2); the core skill is
structurally always installed with any add-on (`targetSkillUnits`,
`installer.js:869-886` puts core at `units[0]`), so `parley-integration.md`'s
dependency on the core protocol can never be installed away; existing
installs' markers keep `doctor` green after the package upgrade
(`markerAddonNames` `:829-838`, `:856-861`) — no false breakage, but also no
signal that a new add-on exists until reinstall (release-note duty); no
`.gitattributes` in the target repo, so CRLF-mangling on Windows checkouts
is possible in theory — Python/JSON/CSV/markdown all tolerate CRLF, so I
rate this noise and propose no change; legacy markers without an `addons`
field stay healthy (`:626-639` test).

**G. What I checked and found clean:** credential/customer scan (only the
linters' own secret-detection regexes and prohibition prose match); the DTVP
adapter ships sanitized maturity metadata only, `"live_effects_authorized":
false`, no procedure IDs, accounts, or prices; placeholder scan clean (only
the linters' marker patterns themselves); `scripts/tests/fixtures/` are
synthetic (`example.invalid` accounts/URLs).

## 5. Verification plan

Order chosen so cheap, source-side gates fail before expensive install-mode
runs. "TESTED" below means I already ran it today against the pre-integration
state and report the observed result; in Phase 5 every one is re-run against
the integrated tree.

1. **#9 hygiene scan of the integrated tree** (BYTE/customer strings,
   credentials, placeholders, caches) — cheapest falsifier, run first.
   Pre-run on the source today: clean (TESTED, §4.G).
2. **#10 source-vs-integrated diff** against the hash snapshot taken at copy
   time — establishes that the tree is what we meant to copy before any
   installer touches it. Intentional differences to document: the 12 rename
   lines (§2), the dropped `.gitignore` + hoisted root rules (F4), stripped
   exec bits (§4.E), excluded caches (§1.1).
3. **#6 the 54 Python tests** in the integrated tree, per-file,
   `PYTHONDONTWRITEBYTECODE=1`. Baseline TESTED today in the source:
   4+20+2+3+15+3+7 = 54, all `OK`, Python 3.14.6.
4. **#7 compile checks with no cache artefacts** — `compile()` built-in over
   all 14 `*.py` (not `py_compile`, which *writes* `.pyc` even under
   `PYTHONDONTWRITEBYTECODE`; not `compileall`, same reason). TESTED today:
   14 files compile, no artefacts written (verified by post-run `find`).
5. **#8 every shipped JSON validates** — the 4 schemas parse + are Draft
   2020-12 artifacts (round-07 consensus already verified seven artifacts;
   re-verified here only for parse via the skill's own test), 4 adapters, 3
   discovery sources, 1 jurisdiction profile. The skill's
   `test_all_json_assets_parse` covers parse; adapter semantics are #5 below.
6. **#5 the adapter validator** — `python3 scripts/adapter_validate.py` over
   all four adapters; expect the locked maturities (dtvp
   live-submission-validated scoped to the 2026-07-22 flow, nrw/elvis
   research-only, manual P0-ceiling). NOT RUN by me yet (runs in Phase 5
   against the integrated tree; the 4 validator tests passed inside #6).
7. **#1 `npm test`** — full Node suite incl. the new F2 leg, F5 branch, and
   cache scan. Baseline TESTED pre-integration: 253/253 pass.
8. **#2 `npm pack --dry-run`** — assert every add-on file present and zero
   `__pycache__`/`.pyc`. Baseline mechanism TESTED (153 files today; nested-
   `.gitignore` behavior characterized, §3 F4). Note `npm pack --dry-run`
   on npm 11 writes no tarball (TESTED: no `.tgz` appeared, `git status`
   clean).
9. **#3 install modes** — default, `--no-addons`, `--only parley-bidding`
   into a temp `$HOME`; assert the six-skill set, the core-only set, and
   core+bidding respectively; assert no `__pycache__` lands. Mechanism
   TESTED for the current five skills via the portable binary (154 files,
   deep trees intact).
10. **#4 `doctor`, `status`, `paths`, `uninstall`** for the new add-on —
    including the F6 manifest negative case (delete `scripts/` →
    `malformed`). Mechanism TESTED for current skills (`doctor` all-valid
    after portable install).
11. **Portable build leg (not in the brief's 11 — I am adding it):**
    `npm run build:portable:current`, then a temp-`$HOME` install from the
    freshly built binary asserting `parley-bidding` arrives complete. The
    brief's list never exercises the pkg channel for the new tree; my lens
    requires it. Feasible on this host (`@yao-pkg/pkg` is in
    `node_modules`, and today's `dist/` v1.5.0 binary proves the toolchain
    runs here — TESTED).
12. **#11 `npx skills add <repo> --list` shows six** — **cannot run on this
    host now, for two independent reasons:** (a) the `skills` CLI is not
    installed and `npx -y skills …` would fetch from the network (not
    permitted this round); (b) the command resolves the *published GitHub
    repo*, which cannot contain `parley-bidding` until merge — so even in
    Phase 5 it is a **post-merge** check unless the CLI accepts a local path
    (NOT TESTED). Pre-merge substitute: assert locally that the integrated
    repo contains exactly six `SKILL.md` files, one per `skills/*` directory
    — that is the discovery invariant the universal installer relies on
    (TESTED for the current five by `find`; trivially re-run for six).

Explicitly not runnable here and not substituted: `pytest` (absent, and not
needed — the suite is unittest); any live Gemini CLI validation (README
`:217` already carries that caveat for the core skill); any real-portal
check (forbidden by the brief's own constraints).

## Sign-off position

Design sound; sequencing (F7) accepted; the plan above is implementable
mechanically once `skills-cli-install-path` clears, with the three
brief-missed rename lines (§2), the E3b precedence sentence (§4.B), the
python3 prerequisite documentation (§4.C), and the portable-build
verification leg (§5.11) added to the shared design.
