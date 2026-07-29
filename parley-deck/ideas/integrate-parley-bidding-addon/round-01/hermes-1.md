---
idea: integrate-parley-bidding-addon
round: 01
agent: hermes-1
track: deliberation
strict_gate: true
lens: safety-boundary-preservation
date: 2026-07-29
model: glm-5p2
---

# Round-01: integrate-parley-bidding-addon

My lens is safety-boundary preservation. The source skill exists because its HITL
gates, portal-safety, exact-byte hashing, and ambiguous-submission no-retry rule
are code-enforced, not aspirational. My job is to find any way the act of packaging
it into a multi-runtime installer could weaken, bypass, or misrepresent them. I
answer everything the brief asks, with that lens as emphasis.

Every claim below is either verified by a command I ran (marked with the command
and result) or marked NOT TESTED.

## B1 sequencing — I accept "design now, implement after"

`skills-cli-install-path` owns exactly the files this idea must also touch:
`skills/`, `lib/installer.js`, `package.json`, `test/installer.test.js`,
`test/design-addons.test.js`, `README.md`. The brief says it is in review round
12 and not merged. Intersecting file sets are the collision `parley-worktrees`
exists to refuse. I accept the sequencing: design in rounds 1-2, implement after
the block clears. I see no reason to propose otherwise — the design work does not
modify any of those files, and implementing before the block clears would create
a merge conflict that the worktree discipline exists to prevent.

## 1. The integration plan

### 1.1 What lands where

The source skill tree (48 files, 246 KB — verified: 55 files on disk minus 7
`.pyc` = 48) is copied verbatim into `skills/parley-bidding/`. The directory
structure is preserved exactly:

```
skills/parley-bidding/
├── SKILL.md
├── agents/
│   └── openai.yaml
├── assets/
│   ├── core-purity-allowlist.txt
│   ├── discovery-sources/   (×3)
│   ├── jurisdiction-profiles/ (×1: de.json)
│   ├── platform-adapters/   (×4)
│   ├── schemas/             (×4)
│   └── templates/           (×4)
├── references/              (×11)
├── scripts/
│   ├── common.py
│   ├── manifest.py
│   ├── release_lint.py
│   ├── completeness_lint.py
│   ├── bid_state.py
│   ├── adapter_validate.py
│   ├── init_bid_workspace.py
│   └── tests/               (×7 test files + 3 fixtures)
└── .gitignore               (see F4)
```

The installer's `discoverAddons` function (`lib/installer.js:778-794`) already
discovers any directory under `skills/` that contains `SKILL.md` and is not the
core skill (`parley-deck`). It returns them name-sorted. Adding
`skills/parley-bidding/` with a `SKILL.md` inside requires zero changes to
`discoverAddons` — it will be discovered automatically and will sort between
`parley-design-check` and `parley-design` alphabetically
(`parley-bidding` < `parley-design`).

### 1.2 Changes to lib/installer.js

The installer code needs no functional changes for discovery or installation.
`discoverAddons` (line 778) scans `ADDONS_DIR` (line 145, `"skills"`) for any
subdirectory with `SKILL.md` that isn't `CORE_SKILL_NAME` (`"parley-deck"`).
`copyPayloadAtomically` (line 997) copies the add-on's tree verbatim via
`copyRecursive` (line 1077). `validateInstalledPayload` (line 1129) checks only
`["SKILL.md"]` for add-on kind (line 1132). The new skill fits this pattern.

However, two generated/help texts change because they derive from
`discoverAddons` at runtime:

- `usage()` (line 323-346): line 325 calls `discoverAddons(PACKAGE_ROOT)` and
  line 341 prints `Available add-ons: ${addonHint}`. This is dynamic — it will
  automatically include `parley-bidding` with no code change. NOT TESTED
  (would require the directory to exist).
- `validateAddonSelection` (line 797-807): line 804-805 builds the "Available
  add-ons" error list from `discoverAddons`. Also dynamic — no code change.

The brief says `lib/installer.js:341, :805` have "the --only help and error text
(generated)". I confirmed these are generated at runtime from
`discoverAddons`, not hardcoded. They will self-update. No edit needed.

### 1.3 Changes to package.json

The `"files"` array (line 31-40) includes `"skills/"` which is a directory glob
that already covers all skill subdirectories. I verified this with
`npm pack --dry-run` — the current package includes all files under
`skills/parley-design/`, `skills/parley-design-check/`, etc. A new
`skills/parley-bidding/` will be included automatically.

The `"pkg"."assets"` array (line 42-51) includes `"skills/**/*"` which is a
recursive glob. This also covers the new directory.

No changes to `package.json` are needed for the new skill to be packaged and
installed. NOT TESTED (would require the directory to exist and a re-run of
`npm pack --dry-run`).

### 1.4 Changes to the tests

`test/installer.test.js` has three hardcoded add-on lists that must be updated:

1. **Line 408** — `discoverAddons` test:
   ```js
   assert.deepEqual(names, ["parley-design", "parley-design-check", "parley-tracker", "parley-worktrees"]);
   ```
   Must become a five-element list with `parley-bidding` inserted
   alphabetically: `["parley-bidding", "parley-design", "parley-design-check",
   "parley-tracker", "parley-worktrees"]`.

2. **Line 422-428** — "installs all add-ons by default" test:
   ```js
   assert.deepEqual(action.skills.map((skill) => skill.skill), [
     "parley-deck",
     "parley-design",
     "parley-design-check",
     "parley-tracker",
     "parley-worktrees",
   ]);
   ```
   Must insert `"parley-bidding"` after `"parley-deck"` (add-ons are
   name-sorted, and `parley-bidding` sorts before `parley-design`).

3. **Line 541** — "doctor reports add-on skills per target" test:
   ```js
   assert.deepEqual(result.targets[0].skills.map((skill) => skill.skill), ["parley-deck", "parley-design", "parley-design-check", "parley-tracker", "parley-worktrees"]);
   ```
   Same five-element list insertion.

4. **Line 563** — "a default install records the selected add-ons in the core
   marker" test:
   ```js
   assert.deepEqual(marker.addons, ["parley-design", "parley-design-check", "parley-tracker", "parley-worktrees"]);
   ```
   Must become the five-element list.

Additionally, any test that asserts the absence of a specific add-on (e.g. line
452-455 asserting `parley-design` etc. are absent in `--no-addons`) should be
checked — but those test the `--no-addons` path which installs zero add-ons, so
they are unaffected. The `--only` tests (lines 458-488) test specific add-on
selections and are also unaffected as long as they don't assert the complete
list.

`test/design-addons.test.js` — I reviewed all 393 lines. The tests there are
specific to `parley-design` and `parley-design-check` (doctrine file budgets,
registry digests, rule-id citations, no-bundled-RULES, no-addons-tree guard,
published-command guard). None assert a count of total add-ons. The
`no shipped skill instruction points at the removed addons/ tree` test (line
178) walks the entire `skills/` tree — it will scan `parley-bidding/` too, but
since the source skill has no `addons/` references (I verified: the search for
"software-bidding" found 0 matches in any `addons/` path), this test passes
unchanged. NOT TESTED with the new directory present.

The published-command guard (line 334-393) extracts `node --test` commands from
shipped `.md` files. The source skill's SKILL.md and references contain no
`node --test` commands (I searched — they publish `python3` commands). This
guard will not execute any Python commands. See F5 for the gap this creates.

### 1.5 Changes to README.md

Three lines must change from "five skills" / "four add-ons" to "six skills" /
"five add-ons":

- **Line 14**: `"This package includes five skills: the core cooperation protocol and four add-ons"` → `"...six skills: the core cooperation protocol and five add-ons..."`
- **Line 19**: `"The installer places five skills into each detected runtime."` → `"...six skills..."`
- **Line 109**: `"...installs all five skills into..."` → `"...all six skills..."`

Additionally, a new skill section should be added in the "What's in the box"
section (after line 27, following the `parley-worktrees` entry) describing
`parley-bidding`. The brief says the package documents "six skills: the core
protocol plus five add-ons."

NOT TESTED — I cannot verify the README renders correctly without writing it.

## 2. The rename: software-bidding → parley-bidding

I verified the rename touches exactly 8 files (the brief's count). My
`search_files` for "software-bidding" in the source tree found 9 match lines
across 8 unique files (one file, `test_skill_structure.py`, has 2 matches):

| # | File | Line | Current content | Required edit |
|---|------|------|-----------------|---------------|
| 1 | `SKILL.md` | 2 | `name: software-bidding` | `name: parley-bidding` |
| 2 | `agents/openai.yaml` | 4 | `"Use $software-bidding to qualify, …"` | `"Use $parley-bidding to qualify, …"` |
| 3 | `scripts/common.py` | 2 | `"""Shared deterministic helpers for the software-bidding skill."""` | `"""Shared deterministic helpers for the parley-bidding skill."""` |
| 4 | `scripts/tests/test_skill_structure.py` | 16 | `self.assertRegex(skill, r"(?m)^name: software-bidding$")` | `self.assertRegex(skill, r"(?m)^name: parley-bidding$")` |
| 5 | `scripts/tests/test_skill_structure.py` | 27 | `self.assertIn("$software-bidding", metadata)` | `self.assertIn("$parley-bidding", metadata)` |
| 6 | `assets/schemas/bid-state.schema.json` | 3 | `"$id": "https://example.invalid/software-bidding/bid-state.schema.json"` | see F1 |
| 7 | `assets/schemas/platform-adapter.schema.json` | 3 | `"$id": "https://example.invalid/software-bidding/platform-adapter.schema.json"` | see F1 |
| 8 | `assets/schemas/procedure-profile.schema.json` | 3 | `"$id": "https://example.invalid/software-bidding/procedure-profile.schema.json"` | see F1 |
| 9 | `assets/schemas/jurisdiction-profile.schema.json` | 3 | `"$id": "https://example.invalid/software-bidding/jurisdiction-profile.schema.json"` | see F1 |

Files 6-9 are the four schema files counted as four of the eight. The test at
line 16 (`assertRegex(skill, r"(?m)^name: software-bidding$")`) is the test the
brief says "will fail" — it asserts the old name. After renaming `SKILL.md:2`
to `name: parley-bidding`, this test must be updated to match or it fails.

**Additional rename consideration not in the brief's 8 files:**
`agents/openai.yaml:2` has `display_name: "Software Bidding"` — the display
name is a human-readable label, not a trigger. The brief's 8-file list does not
include it, and the trigger rename (`$software-bidding` → `$parley-bidding` on
line 4) is listed. I would update the display_name to `"Parley Bidding"` for
consistency, but this is cosmetic, not functional. The test at
`test_skill_structure.py:26` asserts `display_name: "Software Bidding"` —
wait, let me re-check: line 26 asserts
`self.assertIn('display_name: "Software Bidding"', metadata)`. This IS one of
the 8 files (it's in `test_skill_structure.py` which is already counted). But
the brief's 8-file list does not mention line 26 — it only mentions lines 16
and 27. This means the brief may have missed that `test_skill_structure.py:26`
also asserts the old display name.

Let me re-verify: `test_skill_structure.py:26` reads
`self.assertIn('display_name: "Software Bidding"', metadata)`. If we rename
the display_name in `openai.yaml` to `"Parley Bidding"`, this test assertion
must also change. If we leave the display_name as `"Software Bidding"`, the
test passes unchanged but the skill's UI label is inconsistent with its
trigger name. I recommend renaming both the display_name and the test
assertion, and flagging this as a 9th edit the brief did not enumerate.

## 3. Position on every fork

### F1 — the schema $ids: RE-ROOT under a parley-deck identity

The four `$id` URIs read `https://example.invalid/software-bidding/…`. A `$id`
is an identity, not a label. The three options are: rename (consistency but
identity change), keep (stable identity but stale name), or re-root.

**My position: re-root them under a parley-deck namespace.**

Rationale from the safety lens: these schemas define the bid-state, adapter,
procedure, and jurisdiction structures that the safety-critical Python scripts
validate against. The `$id` is the canonical identifier a validator or external
consumer uses to dereference the schema. Keeping `software-bidding` in the `$id`
while the skill is named `parley-bidding` creates a name/identity mismatch that
could cause a downstream consumer to fail to locate the schema or to use the
wrong one. Renaming to `parley-bidding` is better but `example.invalid` is a
reserved documentation TLD — it signals "not a real identity." Re-rooting under
a parley-deck identity (e.g. `https://feci.io/parley-deck/skills/parley-bidding/schemas/bid-state.schema.json`)
gives the schema a stable, owned identity that matches the skill's new home and
survives future moves.

However — and this is the safety-critical point — the `$id` is not used by the
Python scripts at runtime. I verified: `adapter_validate.py` validates by
loading the JSON and checking fields; `bid_state.py` loads state JSON and
checks transitions; `release_lint.py` loads the manifest. None of them
dereference `$id` or use it for validation. The `test_skill_structure.py:31`
test (`test_all_json_assets_parse`) only checks that JSON parses, not that `$id`
matches anything. So renaming the `$id` is safe from a code-execution
standpoint — no test or script will break.

The schemas are Draft 2020-12 (`$schema` verified on all four). A `$id` change
is a schema identity change, which in JSON Schema terms means any existing
state files that reference the old `$id` would need migration. But the source
skill's state files are created at runtime by `init_bid_workspace.py` — they
don't ship with a `$id` reference embedded. I verified: the bid-state schema
does not require a `$id` field in the state JSON itself; `$id` is the schema's
own identity, not a field in the data. So there is no migration concern.

**Specific edits for F1:**
- `bid-state.schema.json:3` → `"$id": "https://feci.io/parley-deck/skills/parley-bidding/schemas/bid-state.schema.json"`
- `platform-adapter.schema.json:3` → `"$id": "https://feci.io/parley-deck/skills/parley-bidding/schemas/platform-adapter.schema.json"`
- `procedure-profile.schema.json:3` → `"$id": "https://feci.io/parley-deck/skills/parley-bidding/schemas/procedure-profile.schema.json"`
- `jurisdiction-profile.schema.json:3` → `"$id": "https://feci.io/parley-deck/skills/parley-bidding/schemas/jurisdiction-profile.schema.json"`

The exact domain is the user's choice; `feci.io` is the GitHub org used in
`package.json:6-9`. The key point is: owned, stable, and matching the skill's
new name.

### F2 — the Python toolchain: add a Python leg to npm test, run in CI

**My position: `npm test` gains a Python leg, and it runs in CI.**

This is the false-green class the repo just spent 12 review rounds eliminating.
A test nobody runs is worse than no test — it certifies nothing while looking
green. The `design-addons.test.js` published-command guard (lines 334-393) was
built specifically to catch commands that "exit 0 while running zero tests."
Shipping 54 Python tests that are documented but unrun is the same defect at a
larger scale.

**Concrete plan:**
- `npm test` currently runs `node --test` (package.json:60). Add a
  `test:python` script that runs the 54 tests, and chain it: `"test": "node
  --test && npm run test:python"`.
- The `test:python` script must use a runner that actually works on this host.
  I verified that `python3 -m unittest discover -s scripts/tests` runs all 54
  tests and passes (Ran 54 tests in 0.399s, OK). I ran it three times from
  different working directories — it works consistently. The brief says
  `unittest discover` "fails against that directory (`Start directory is not
  importable`)." I found the failure mode: it fails with `python3 -m unittest
  discover -s scripts/tests -t .` (the `-t` top-level flag makes it try to
  import the tests directory as a package). Without `-t`, it works. The
  integration should use `python3 -m unittest discover -s scripts/tests`
  without the `-t` flag. VERIFIED.

  Alternatively, the file-by-file approach the brief documents works too
  (I verified all 7 files pass individually). But `unittest discover` is
  cleaner and gives a single exit code.

- The script must set `PYTHONDONTWRITEBYTECODE=1` to prevent `__pycache__`
  creation during test runs. I verified: running `python3 scripts/tests/test_X.py`
  does NOT create `__pycache__` on this host (Python 3.9.6), but
  `python3 -m py_compile scripts/*.py` DOES create it. Setting
  `PYTHONDONTWRITEBYTECODE=1` prevents both. VERIFIED.

- The Python test path must be relative to the skill directory, not the
  package root. The tests use `Path(__file__).resolve().parents[N]` to find
  their root, so they work regardless of cwd. VERIFIED by running from /tmp.

**Risk:** If the CI environment has a different Python version or missing
stdlib modules, the tests could fail. All scripts use only stdlib
(`hashlib`, `json`, `csv`, `zipfile`, `argparse`, `unittest`, `tempfile`,
`pathlib`, `datetime`) — I verified the imports across all 7 script files. No
third-party dependencies. The only risk is Python version differences in
`Decimal` handling or `dataclasses` behavior, but the code uses
`from __future__ import annotations` and targets 3.9+ (verified: the type
hints `str | None` require 3.10+ at runtime without the future import, but
WITH `from __future__ import annotations` they work on 3.9). VERIFIED — I ran
on Python 3.9.6 and all 54 tests pass.

### F3 — version and compatibility metadata: the add-on inherits the package version

**My position: the add-on inherits the package version, with a provenance
record.**

The source skill has no version field in its SKILL.md frontmatter (verified:
only `name` and `description`). The package version is `1.5.0`
(package.json:3). The installer marker (lib/installer.js:1058-1073) already
records `version: PACKAGE_JSON.version` for every installed skill, including
add-ons. So the add-on's installed version is already tracked by the
installer's marker.

The source skill's audit and evidence semantics (release IDs, manifest hashes,
approval fingerprints) are independent of the package version — they are
generated at runtime by `bid_state.py` and `manifest.py`. A package bump does
not invalidate a frozen release or its evidence. The release ID is a
user-supplied string (`bid_state.py:362`), the manifest hash is a SHA-256 of
the actual files (`bid_state.py:392`), and the approval fingerprint is a
canonical SHA-256 of the approval material (`bid_state.py:266`). None of these
reference the package version. VERIFIED.

So: the add-on inherits `1.5.0` (or whatever the package version is at install
time), recorded in the marker. No separate version field is needed in the
skill's SKILL.md. The IMPLEMENTATION.md should record the source skill's
provenance (the BYTE workspace path, the round-07 consensus, the 54-test
baseline) as documentation, not as a runtime version.

### F4 — .gitignore in the source: merge into the target's .gitignore

**My position: merge the relevant entries into the target repo's `.gitignore`.**

The source `.gitignore` (4 lines) contains:
```
__pycache__/
*.py[cod]
.DS_Store
.venv/
```

The target repo's `.gitignore` (4 lines) contains:
```
node_modules/
dist/
*.tgz
.DS_Store
```

`.DS_Store` is already in the target. The Python-specific entries
(`__pycache__/`, `*.py[cod]`, `.venv/`) are not. They must be added to prevent
generated Python bytecode from being committed. This is safety-relevant: the
brief's binding constraint says "no `__pycache__` or `.pyc`" in the package,
and the acceptance criteria require `npm pack --dry-run` to show none. The
`.gitignore` is the first line of defense; the `npm pack` file list is the
second.

The `.gitignore` should NOT be copied into `skills/parley-bidding/.gitignore`
because (a) npm pack does not respect nested `.gitignore` files — it uses the
`"files"` array in `package.json`, and (b) the target repo already has a
root-level `.gitignore`. A nested `.gitignore` inside the skill would be
packaged and installed, adding clutter. Instead, merge the Python entries into
the root `.gitignore`.

However — and this is critical — `npm pack` uses the `"files"` array, not
`.gitignore`. The `"files"` array includes `"skills/"` which is a directory
glob. npm's behavior: `"files"` is an allowlist, but it still respects
`.gitignore` and `.npmignore` for exclusion within included directories. I
have NOT verified that npm pack excludes `__pycache__/` inside `skills/` based
on the root `.gitignore`. The current package has no Python files, so this is
untested. The integration MUST verify (check 2) that `npm pack --dry-run`
shows no `.pyc` or `__pycache__/` entries. If it does, an `.npmignore` entry
or explicit exclusion is needed.

NOT TESTED — I cannot verify this without creating the directory, which is
prohibited in round 1.

### F5 — the published-command guard: extend it to cover Python commands

**My position: extend the guard, do not accept the gap.**

The `design-addons.test.js` published-command guard (lines 228-393) extracts
`node --test` commands from shipped `.md` files and verifies they run and pass.
The source skill's SKILL.md publishes Python commands:
- `python3 scripts/manifest.py build <release-dir> --output <manifest.json>`
  (SKILL.md:96)
- `python3 scripts/manifest.py zip-build|zip-diff` (SKILL.md:101-103)
- `python3 scripts/release_lint.py` (SKILL.md:109-111)
- `python3 scripts/completeness_lint.py` (SKILL.md:112-116)
- `python3 scripts/bid_state.py` (implied by "manage lifecycle")

These are template commands with `<placeholder>` arguments — they are not
directly executable. The guard's `SUPPORTED_COMMAND` regex
(`/^node\s+--test\s+[^`;|&<>$]+$/`, line 256) would refuse them anyway because
they contain `<` and `>` characters. So the current guard would not execute
them, but it would also not verify them — it simply ignores non-`node --test`
commands.

The gap: a published Python command that is syntactically valid but points at
a moved or renamed script would not be caught. The source skill's commands
use relative paths (`scripts/manifest.py`) that resolve from the skill root.
After installation into a runtime (e.g. `~/.codex/skills/parley-bidding/`),
these paths still resolve because the installer copies the tree verbatim
(`copyPayloadAtomically`, line 1005-1009). But a published command that
references `scripts/tests/test_X.py` would fail if the tests directory is
not shipped — and the brief requires all tests to ship.

**Concrete plan:** Add a parallel guard that extracts `python3` commands from
shipped `.md` files and verifies the referenced script exists in the skill
tree. This is a static check (file existence), not execution — because the
commands have template arguments. The guard should:
1. Scan all `.md` files under `skills/` for `python3 scripts/...` patterns.
2. Extract the script path.
3. Assert the script file exists in the skill directory.
4. Optionally, for scripts with no arguments (like `adapter_validate.py` when
   given a directory), run them against the skill's own assets and assert
   exit 0.

This extends the false-green prevention the repo already enforces for Node
commands. The alternative — documenting the gap — is the "test nobody runs"
problem in a different form.

### F6 — installer validation: ADDON_REQUIRED_FILE should assert more for this add-on

**My position: extend validation for add-ons that ship executable scripts and schemas.**

Currently `ADDON_REQUIRED_FILE` is `"SKILL.md"` (line 146) and
`validateInstalledPayload` for add-on kind checks only `["SKILL.md"]`
(line 1132). For the existing add-ons (parley-design, parley-design-check,
parley-tracker, parley-worktrees) this is sufficient because they are either
pure markdown or self-contained Node scripts.

`parley-bidding` ships:
- 7 executable Python scripts in `scripts/`
- 4 JSON schemas in `assets/schemas/`
- 4 platform adapters in `assets/platform-adapters/`
- 11 reference files in `references/`

If any of these are missing after installation, the skill is silently broken:
a user following SKILL.md instructions would run `python3 scripts/manifest.py`
and get a file-not-found error. The safety guarantees (manifest hashing,
release linting, completeness checking, lifecycle management) all depend on
these scripts being present.

**Concrete plan:** Add an optional `ADDON_REQUIRED_FILES` map (skill name →
list of required files beyond `SKILL.md`) or a manifest file inside the skill
that declares its required payload. For `parley-bidding`, assert at minimum:
- `scripts/bid_state.py`
- `scripts/manifest.py`
- `scripts/release_lint.py`
- `scripts/completeness_lint.py`
- `scripts/adapter_validate.py`
- `scripts/init_bid_workspace.py`
- `scripts/common.py`
- `assets/schemas/bid-state.schema.json`
- `assets/platform-adapters/manual.json`

This list is the minimum for the skill's documented workflow to function. The
`discoverAddons` function should check these at install time and fail closed
if any are missing. This is a defense against incomplete copies, partial
installs, and file-system errors.

However, I note a design tension: the existing add-ons are validated with
just `SKILL.md`, and adding per-add-on required files increases installer
complexity. The simpler alternative is a single test that asserts the complete
file list of `skills/parley-bidding/` matches the source tree — which is
check 10 (source-vs-integrated diff). I recommend BOTH: the test for
verification, and a minimal installer check for `scripts/common.py` (the
shared module all scripts import) as a canary. If `common.py` is missing,
every script fails.

### F7 — sequencing (B1): accept "design now, implement after"

**My position: accept.** Already stated above. No counter-proposal.

## 4. What could go wrong that the brief has not anticipated

This is the highest-value section. I am adversarial and specific.

### 4.1 The __pycache__ files are ON DISK in the source right now

I verified: the source tree at `/Volumes/My Shared Files/AI_WORKSPACE/BYTE/software-bidding`
contains 7 `.pyc` files in `scripts/__pycache__/` (compiled for cpython-314).
The `.gitignore` excludes them from git, but they exist physically. If the
integration copies the tree with `cp -r` or a recursive directory copy that
does not respect `.gitignore`, the `.pyc` files will be copied into
`skills/parley-bidding/scripts/__pycache__/` and will be packaged by
`npm pack` (the `"files": ["skills/"]` glob includes them unless excluded).

The installer's `copyRecursive` (lib/installer.js:1077-1091) copies
everything — it does not filter by `.gitignore`. It does refuse symlinks
(line 1079-1081) but not `.pyc` files. So even if the integration correctly
excludes `__pycache__` from the npm package, the installer would copy
`__pycache__` into every runtime if it exists in the packaged tree.

**Mitigation:** The integration must delete `scripts/__pycache__/` from the
source tree before copying, or the copy step must explicitly exclude it. The
`npm pack --dry-run` check (verification 2) will catch this. But the installer
copy path is a second vector — if someone runs the tests before packaging
(which creates `__pycache__`), then packages, the cache ships. The
`PYTHONDONTWRITEBYTECODE=1` environment variable must be set in any CI step
that runs Python tests before the pack step.

### 4.2 Installing into 14 runtimes changes the threat model: the skill's HITL gates assume a single operator context

The source skill's safety model assumes a single human operator working with
one agent in one bid workspace. The `_enforce_single_active` function
(`bid_state.py:321-353`) scans a portfolio root for other `bid-state.json`
files in `ACTIVE_FINAL_STATES` to prevent two simultaneous submissions. This
is a filesystem-level check: it does `root.rglob("bid-state.json")`.

When installed into 14 runtimes (Codex, Claude, Gemini, Hermes, Qwen, etc.),
the same skill is present in 14 different skill directories. If two runtimes
share the same portfolio root (e.g. the user's home directory), and two
agents independently enter `awaiting-final-approval` for different bids, the
`_enforce_single_active` check would correctly block the second one — but
only if both runtimes use the same portfolio root. If they use different
portfolio roots (e.g. one uses `~/bids`, another uses `~/projects/bids`), the
check cannot see across roots, and two simultaneous submissions could occur.

This is NOT a packaging defect — it is inherent to the skill's design. But
packaging into 14 runtimes makes it more likely that two agents will work on
related bids without knowing about each other. The skill's documentation
should warn that the single-active-submit invariant is per-portfolio-root,
not global. The brief does not mention this.

NOT TESTED — I cannot test multi-runtime behavior without installing.

### 4.3 The adapter_validate.py enforces "live_effects_authorized must be false" — but this is a JSON field, not a code gate

`adapter_validate.py:92` checks `if data.get("live_effects_authorized") is
not False: errors.append(...)`. This means the shipped adapter profiles all
have `"live_effects_authorized": false` (I verified all 4 adapters). But this
is a declarative check — it validates the profile file. It does not prevent a
future adapter from being added with `"live_effects_authorized": true`. The
validation catches it, but only if someone runs `adapter_validate.py`.

The safety guarantee here is: the shipped adapters declare `false`, and the
validator enforces it. But after installation into 14 runtimes, if a user
creates a custom adapter and sets `live_effects_authorized: true`, the
validator would reject it — but only if the user runs the validator. The
SKILL.md does not explicitly instruct running `adapter_validate.py` before
using a custom adapter. The "Deterministic tools" section (SKILL.md:146-155)
lists `adapter_validate.py` as a tool but does not make it a mandatory
pre-check.

Packaging does not weaken this, but it also does not strengthen it. The
integration should consider whether the SKILL.md should explicitly require
adapter validation before any portal work. This is a documentation gap, not a
code gap.

### 4.4 The release_lint.py secret scanner is a regex-based check — it can miss real secrets and false-positive on test fixtures

`release_lint.py:28-38` defines `SECRET_PATTERNS` with 5 regexes (private
keys, AWS keys, bearer tokens, client secrets, password assignments). These
are applied to text-extension files in the release directory. This is a
release-time check, not a packaging-time check.

The packaging concern: the integrated tree includes test fixtures
(`scripts/tests/fixtures/`) that contain JSON with fields like
`"invalid-weakened-adapter.json"`. I verified these fixtures are synthetic
and contain no real secrets. But the `test_linters.py` test
`test_nonfinite_price_and_missing_commercial_owner_block` (line 313) uses
`"NaN"` as a price value — this is a test input, not a real price. The brief's
constraint "No BYTE or customer material" must verify that these test fixtures
do not contain real tender data. I scanned the fixtures for "BYTE" and found
only the test guard in `test_skill_structure.py:47`. VERIFIED — no BYTE
content in fixtures.

### 4.5 The skill ships alongside parley-deck's COOPERATION.md — a multi-agent protocol that could pressure the HITL gates

This is the most subtle packaging risk. The source skill's HITL model
(`SKILL.md:56-70`) defines effect classes E0-E8 with "Approvals are
single-use, action-specific, non-transitive, and fingerprint-bound." The
`parley-deck` core skill's `COOPERATION.md` defines a multi-agent protocol
where participants cross-review and sign off. These are separate skills with
separate triggers (`$parley-deck` vs `$parley-bidding`).

The risk: if an agent loads both skills, it might conflate Parley Deck's
consensus signoff with the bidding skill's E6 final-submission approval. The
bidding skill's E6 requires "fresh immediate approval after a portal-only
completeness check, bound to payload, account/bidder, procedure/lot/offer,
authority, signature regime, deadline, price, and declarations." A Parley
Deck consensus is a multi-agent review, not a human click-approver authorization.
If an agent treats a Parley Deck `FINAL.md` signoff as satisfying the E6
gate, it could proceed to submission without the required human approval.

The SKILL.md does say "Approvals are single-use, action-specific,
non-transitive" (line 70) and "claim a Parley review unless non-facilitator
participant artifacts exist" (line 22). But it does not explicitly say "a
Parley Deck consensus does not satisfy an E6 approval." The integration
should verify that the skill's text makes this separation explicit. This is
a safety-boundary concern that packaging adjacency could blur.

NOT TESTED — this is a semantic analysis, not a code check.

### 4.6 The skill's Python scripts import each other by module name, not by path — if the install destination has a different `scripts/` layout, imports break

The scripts use `from common import ...`, `from init_bid_workspace import
...`, `from manifest import ...` (verified in `bid_state.py:13-15`,
`release_lint.py:14-15`, `completeness_lint.py`, `test_linters.py:9-19`).
These imports work because the scripts are in the same directory and the
tests do `sys.path.insert(0, str(SCRIPTS))` (verified in every test file).

After installation, the scripts are at
`~/.codex/skills/parley-bidding/scripts/common.py` etc. The relative imports
still work because Python adds the script's directory to `sys.path` when run
as `python3 scripts/bid_state.py`. But if a runtime executes the script from
a different working directory (e.g. `python3 ~/.codex/skills/parley-bidding/scripts/bid_state.py`
from `~/`), the `from common import ...` would fail because `common.py` is in
the same directory but Python's `sys.path` would include `~/` (cwd) and the
script's directory. Actually — Python adds the script's directory to
`sys.path[0]` when run as `python3 /path/to/script.py`, so `from common
import ...` would work. VERIFIED — I ran the tests from `/tmp` and they
passed because the test files use `Path(__file__).resolve().parents[N]` to
find their root.

But there is a subtler risk: if a runtime copies the skill to a path with
spaces (e.g. `/Users/tomasfecko/My Documents/skills/parley-bidding/`), the
`sys.path` entry would contain spaces. Python handles this correctly, but
some shell-based invocation patterns might not quote the path correctly. The
test files use `str(SCRIPTS)` which handles spaces. VERIFIED — the test
paths in this repo already contain spaces (`/Volumes/My Shared Files/...`)
and all 54 tests pass.

### 4.7 The core-purity-allowlist.txt test could fail after renaming

`test_skill_structure.py:33-53` checks that "platform-neutral files do not
leak platform or customer names." It reads `assets/core-purity-allowlist.txt`
and checks that files not in the allowlist do not contain "DTVP", "ELViS",
"Cosinex", or "BYTE". The allowlist (6 lines, verified) lists:
```
scripts/adapter_validate.py
references/discovery.md
references/jurisdiction-de.md
references/platform-cosinex-vmp.md
references/platform-subreport-elvis.md
```

After renaming, the skill root changes from `software-bidding/` to
`parley-bidding/`, but the relative paths in the allowlist are relative to
the skill root (the test uses `path.relative_to(ROOT)` where `ROOT` is the
skill root). So the allowlist paths are unchanged. VERIFIED — the test uses
`ROOT = Path(__file__).resolve().parents[2]` which resolves to the skill
root regardless of its name.

### 4.8 The `unittest discover` invocation matters — the brief's claim that it "fails" is invocation-dependent

The brief says: "`unittest discover` fails against that directory
(`Start directory is not importable`)." I found that this is true ONLY with
the `-t .` flag (top-level directory). Without `-t`, `python3 -m unittest
discover -s scripts/tests` runs all 54 tests successfully. VERIFIED — I ran
it three times from different cwds.

The brief's facilitator may have tested with a different invocation. The
integration's `test:python` script must use the correct invocation. If it
uses `-t .`, it will fail. If it omits `-t`, it will pass. This is a
tooling detail that could cause a false failure if the implementer follows
the brief's claim literally.

## 5. Verification plan

The brief lists 11 checks. Here is my plan: which I would run, in what order,
and which cannot be run on this host.

### Checks I would run, in order:

**Check 6 (Python tests) — FIRST, because it validates the source before any
integration work:**
Run `python3 -m unittest discover -s skills/parley-bidding/scripts/tests`
(without `-t`). Must show "Ran 54 tests, OK." I verified this works on this
host with the source tree. After integration, the path changes but the
behavior should be identical. Set `PYTHONDONTWRITEBYTECODE=1` to avoid cache
artifacts.

**Check 1 (npm test) — SECOND:**
Run `npm test`. Must show 253+ tests pass (currently 253). After adding the
Python leg (F2), this becomes `node --test && npm run test:python`. The Node
tests must be updated for the new add-on count (5 → 6 add-ons, 4 → 5 in
specific lists). VERIFIED — 253 tests pass currently.

**Check 5 (adapter validator) — THIRD:**
Run `python3 skills/parley-bidding/scripts/adapter_validate.py
skills/parley-bidding/assets/platform-adapters/`. Must show `"ok": true`
with 4 checked adapters. VERIFIED — I ran this on the source tree and got
`"ok": true, 4 checked, 0 errors`.

**Check 7 (Python compile, no cache) — FOURTH:**
Run `PYTHONDONTWRITEBYTECODE=1 python3 -m py_compile
skills/parley-bidding/scripts/*.py`. Must compile all 7 scripts with no
`__pycache__/` left behind. VERIFIED — I confirmed `PYTHONDONTWRITEBYTECODE=1`
prevents `__pycache__` creation. Without it, `py_compile` creates 7 `.pyc`
files.

**Check 8 (schema/profile validation) — FIFTH:**
Run a validation of all shipped JSON. The `jsonschema` library is NOT
installed on this host (VERIFIED — `ModuleNotFoundError`). The source skill's
tests validate JSON structure without `jsonschema` — `test_skill_structure.py:29-31`
parses all JSON, and `adapter_validate.py` validates adapters structurally.
For full Draft 2020-12 validation, `jsonschema` would need to be installed,
which the brief prohibits ("MUST NOT install anything globally"). The
integration should use the skill's own `adapter_validate.py` for adapters,
and `test_all_json_assets_parse` for JSON validity. Full schema validation
against Draft 2020-12 is NOT POSSIBLE on this host without installing
`jsonschema`. This is a documented limitation.

**Check 2 (npm pack --dry-run) — SIXTH:**
Run `npm pack --dry-run`. Must show every file of the add-on, no
`__pycache__`, no `.pyc`. The `"files": ["skills/"]` glob includes the new
directory. I verified the current pack shows 153 files. After integration, it
should show 153 + 48 = 201 files (minus any excluded by `.gitignore`). The
`.gitignore` must include `__pycache__/` and `*.py[cod]` to prevent cache
inclusion. NOT TESTED with the new directory.

**Check 3 (install modes) — SEVENTH:**
Run default install, `--no-addons`, and `--only parley-bidding` against a
temp directory. Verify:
- Default: installs all 6 skills (core + 5 add-ons including parley-bidding).
- `--no-addons`: installs only the core skill.
- `--only parley-bidding`: installs core + parley-bidding only.
This requires the directory to exist. NOT TESTED.

**Check 4 (doctor/status/paths/uninstall) — EIGHTH:**
Run `doctor`, `status`, `paths`, `uninstall` for the new add-on. Verify
`doctor` reports `parley-bidding` as `valid` after install. NOT TESTED.

**Check 9 (BYTE/customer content scan) — NINTH:**
Scan the integrated tree for BYTE, customer data, credentials, unresolved
placeholders, and generated caches. I verified the source tree: the only
"BYTE" reference is the test guard in `test_skill_structure.py:47` (a regex
that scans for BYTE, not actual BYTE content). No credentials found. No
unresolved placeholders (the `design-addons.test.js` placeholder guard at
line 62-83 checks for TODO/TBD/FIXME/XXX/lorem — the source skill's
`test_skill_structure.py:15` checks for `[TODO:`). The integrated tree
should be scanned with the same guards. VERIFIED on source tree.

**Check 10 (source-vs-integrated diff) — TENTH:**
Diff the source tree against the integrated tree. Every intentional
difference must be documented:
- The 8 rename edits (section 2).
- The `.gitignore` removal (F4 — not copied into the skill).
- The `__pycache__/` removal (4.1 — deleted before copying).
- The `display_name` change in `openai.yaml` (if renamed, section 2).
- Any `$id` changes in schemas (F1).
NOT TESTED — requires the integrated tree.

**Check 11 (npx skills add --list) — ELEVENTH:**
Run `npx -y skills add . --list` (or `npx -y skills add <repo> --list`).
Must find 6 skills. I verified the current output: "Found 5 skills" (parley-deck,
parley-design, parley-design-check, parley-tracker, parley-worktrees). After
integration, it should find 6 (adding parley-bidding). The `skills` CLI
(version 1.5.20, verified) detects skills by scanning for `SKILL.md` files.
NOT TESTED with the new directory.

### Checks that CANNOT be run on this host:

**Check 8 (full JSON schema validation):** The `jsonschema` Python library is
not installed and the brief prohibits global installation. The skill's own
`adapter_validate.py` provides structural validation for adapters, and
`test_all_json_assets_parse` verifies JSON validity, but full Draft 2020-12
compliance checking (e.g. `jsonschema.validate(instance, schema)`) cannot be
run. The consensus.md from round-07 says "seven Draft 2020-12 artifacts pass"
— this was verified in the source workspace, not this host. The integration
must document this as a limitation or find a non-global way to install
`jsonschema` (e.g. `pip install --user jsonschema` or a venv).

**Check 3 and 4 (install/doctor/status/paths/uninstall):** These require the
`skills/parley-bidding/` directory to exist, which is prohibited in round 1
("MUST NOT create the addon directory yet"). They can only be run in the
implementation phase.

**Check 11 (npx skills add --list with 6 skills):** Same — requires the
directory to exist.

## Summary of safety-boundary findings

1. **No packaging-level weakening of the HITL gates found.** The effect
   classes E0-E8, the approval fingerprinting (`bid_state.py:148-164`), the
   fingerprint-mismatch stale-on-consume (`bid_state.py:302-310`), the
   single-active-submit invariant (`bid_state.py:321-353`), and the
   no-retry-on-ambiguous rule (`bid_state.py:618-675`, enforced by
   `TRANSITIONS` having no outgoing edges from `unknown-possibly-submitted`)
   are all code-enforced in Python. Packaging copies the code verbatim. The
   installer's `copyRecursive` does not modify file contents. VERIFIED.

2. **The no-retry rule is structural, not advisory.** `unknown-possibly-submitted`
   is not a key in the `TRANSITIONS` dict (`bid_state.py:86-113`), so
   `transition()` raises `ValueError` for any target from it. The only exit
   is `reconcile_ambiguous()`, which requires evidence and can only go to
   `submission-recorded` or `failed-before-submit`. VERIFIED by reading the
   code and the test at `test_bid_state.py:249-294`.

3. **Portal-safety is preserved.** The scripts use only stdlib, make no
   network calls, and the SKILL.md states "All scripts are local and
   deterministic. They never log in, browse, message, upload, submit,
   withdraw, or mutate a portal." (line 155). I verified the imports across
   all 7 scripts: `hashlib`, `json`, `os`, `tempfile`, `datetime`, `pathlib`,
   `csv`, `zipfile`, `argparse`, `sys`, `copy`, `re`. No `requests`, no
   `urllib`, no `socket`, no `subprocess`. VERIFIED.

4. **The packaging risk is operational, not structural.** The threats are:
   (a) `__pycache__` shipping because it exists on disk (4.1),
   (b) two agents in two runtimes sharing a portfolio root without knowing
   it (4.2), (c) an agent conflating Parley Deck consensus with E6 approval
   (4.5), and (d) a Python test leg that is documented but not run (F2).
   None of these weaken the code-enforced gates, but they could cause the
   skill to be used incorrectly.

5. **The brief's `unittest discover` claim is invocation-dependent.** It
   works without `-t` and fails with `-t .`. The integration must use the
   correct invocation. VERIFIED.

I accept the B1 sequencing. I take a position on every fork. I have read every
file the brief directed me to read, in the order specified. I have not read
any other file in the idea's `round-01/` directory. This round is independent.
