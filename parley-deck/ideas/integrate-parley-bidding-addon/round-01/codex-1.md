---
idea: integrate-parley-bidding-addon
round: 01
agent: codex-1
date: 2026-07-29
---

# Position

Integrate the standalone snapshot as `skills/parley-bidding/`, keep it separate from the
core protocol, and treat it as a safety-sensitive add-on rather than as an ordinary
`SKILL.md`-only directory. I accept B1/F7 exactly: finish design now, but do not start Phase
5 until `skills-cli-install-path` reaches zero agreed fixes and its result is merged into the
target integration branch. The two ideas overlap in `skills/`, `lib/installer.js`,
`package.json`, both named test files, and `README.md` (`00-prompt.md:30-39`).

The inspected target is clean on branch `readme-skill-catalogue`. Baseline observations:

- **RUN:** the source contains 48 regular files, no symlinks, and no `__pycache__`, `.pyc`, or
  `.pyo`. Five Python files have an executable bit; the other Python files are regular
  `0644` files.
- **RUN:** all seven source test files passed individually with
  `PYTHONDONTWRITEBYTECODE=1 python3 -B`, reproducing exactly 54 tests
  (`4+20+2+3+15+3+7`).
- **RUN:** `python3 -B scripts/adapter_validate.py assets/platform-adapters` validated all
  four adapter profiles with zero errors.
- **RUN, target baseline only:** `npm test` passed 253/253 tests. The integrated add-on is
  absent, so the integration is **NOT TESTED**.
- **RUN, target baseline only:** `npm pack --dry-run --json` reported 153 package entries.
  The first attempt failed because the user's default npm cache is root-owned; retrying with
  `npm_config_cache=/private/tmp/parley-npm-cache-codex-1` passed. The integrated add-on is
  absent, so its npm distribution is **NOT TESTED**.

Relevant safety evidence was read in the source at `SKILL.md:10-22`, `SKILL.md:55-70`,
`SKILL.md:88-120`, and `SKILL.md:122-155`; `references/hitl-and-recovery.md:21-46`;
`references/parley-integration.md:3-16,31-41`; and
`references/platform-adapter-contract.md:21-57`. The implementation enforcers are visible at
`scripts/adapter_validate.py:30-57,73-149` and `scripts/common.py:28-61`. The ratified design
requires a neutral core and non-weakenable platform profiles (`FINAL.md:72-85`), and the
round-07 review consensus says the standalone tree closed strict review with the authority,
proof-ceiling, immutable-release, no-retry, and recovery guarantees intact
(`review/round-07/consensus.md:15-34`).

# Integration plan

Everything below is **PROPOSED — NOT TESTED** until Phase 5 runs it on the post-B1 target.

## 1. The skill tree

Create `skills/parley-bidding/` only after B1 clears.

1. Copy the source tree without its nested `.gitignore`.
2. Apply the eight-file rename listed below.
3. Add `skills/parley-bidding/parley-addon.json`, a package-owned compatibility and integrity
   manifest. It should declare:
   - manifest schema version `1`;
   - skill name `parley-bidding`;
   - independent skill version `1.0.0`;
   - runtime requirement `Python >=3.10`;
   - a lexicographically sorted inventory of every payload file except the manifest itself,
     with raw SHA-256 for each;
   - one aggregate SHA-256 over the canonical path/hash list; and
   - the five POSIX-executable Python paths as metadata, without making POSIX mode a Windows
     validity condition.
4. Preserve every non-renamed source byte. With the source `.gitignore` merged at repository
   level and one new manifest added, the integrated skill remains 48 files: 39 byte-identical
   source files, 8 intentionally edited files, and 1 new manifest in place of the source-only
   `.gitignore`.
5. Add two packaging-context sentences to `SKILL.md` near the Parley integration boundary:
   Parley Deck's generic external-context default never satisfies E3b for tender material,
   and no agent consensus/signoff can substitute for a named human E5/E6/E7 approval. These
   are stricter clarifications, not weaker semantics. They are an intentional source diff
   separate from the rename.

`SKILL.md:145-155` says all seven scripts are local and deterministic and never operate a
portal. The integrated tree must retain exactly those seven modules and seven test files.
Do not add a browser driver, network client, credential helper, or portal mutation script.

## 2. `lib/installer.js`

The existing path model is already correct:

- keep `CORE_SKILL_DIR = skills/parley-deck` (`lib/installer.js:115-125`);
- keep `ADDONS_DIR = skills` (`lib/installer.js:145-146`);
- do not add bidding files to core `PAYLOAD_ENTRIES` (`lib/installer.js:132-144`);
- let `discoverAddons()` find `skills/parley-bidding/SKILL.md` dynamically
  (`lib/installer.js:775-793`); and
- keep help and unknown-name output generated from discovery
  (`lib/installer.js:325-341,796-805`).

Make these concrete changes:

1. Replace the stale count-bearing comment at `lib/installer.js:785-786` ("all five") with
   "all packaged skills". No generated list should be hardcoded to six.
2. Add an optional `ADDON_MANIFEST_FILE = "parley-addon.json"` contract. Legacy add-ons
   without it remain `SKILL.md`-validated for backward compatibility; when the manifest is
   present, discovery and package preflight must reject a missing file, hash mismatch,
   duplicate/escaping path, symlink, generated cache, malformed version, or wrong skill name.
3. Pass the complete skill unit into `validateInstalledPayload`, rather than only
   `unit.kind`. Update both call sites at `lib/installer.js:1035-1036` and
   `lib/installer.js:1114-1125`. For `parley-bidding`, `validateInstalledPayload` must verify
   the full manifest and raw-byte hashes, not only `SKILL.md`; the current add-on branch at
   `lib/installer.js:1129-1147` is too weak.
4. Record `addonVersion` and `payloadSha256` in the add-on marker. Preserve the core marker's
   existing `addons: false | string[]` shape at `lib/installer.js:1056-1074`, so old
   core-only and selective installs remain healthy (`lib/installer.js:823-862`).
5. Make `doctor` check the declared Python requirement and report a distinct
   payload-valid/runtime-unavailable result. Installation may still copy documentation on a
   Python-less host, but `doctor` must not call the add-on operationally healthy.
6. Preserve executable bits on POSIX when `copyRecursive` writes regular files. The function
   is described as copying the tree verbatim at `lib/installer.js:1005-1009`, but
   `fs.writeFileSync` at `lib/installer.js:1077-1091` preserves bytes, not the source mode.
   The documentation invokes scripts through Python, so Windows validity must not depend on
   an executable bit.
7. Preflight every selected unit before the first destination mutation. Today
   `installTarget` maps directly over units (`lib/installer.js:888-891`), so one unmarked
   `parley-bidding` collision can leave a partially updated six-skill set and a core marker
   that claims the full selection. A predictable collision or invalid package must cause
   zero writes for that target.
8. Add fail-closed cache checks for `__pycache__`, `*.pyc`, and `*.pyo` in packaged and
   installed manifest validation. `listVisibleEntries` only filters `.DS_Store`
   (`lib/installer.js:1203-1208`).

The output shapes for install, status, paths, doctor, and uninstall already operate on
`targetSkillUnits` (`lib/installer.js:865-905,1099-1126,1211-1283`). They should not gain a
bidding-specific branch; the stronger manifest is the extension point.

## 3. `package.json`, lockfile, CI, npm and portable assets

On the inspected branch:

- bump the distribution by one minor version, `1.5.0` to `1.6.0`, because the public skill
  catalogue grows; update `package-lock.json:2-16` in the same change;
- add `bidding` and `procurement` keywords;
- keep production dependencies empty;
- add `ajv` as a dev dependency so Draft 2020-12 schemas and their shipped instances can be
  validated in the Node suite without adding a runtime dependency; and
- retain `"test": "node --test"` because the Node suite will invoke the Python leg.

No new npm/package glob is required: `"files": ["skills/", ...]` already ships the complete
tree (`package.json:31-40`), and `"pkg.assets": ["skills/**/*", ...]` already embeds it in
portable binaries (`package.json:41-51`). Tests must prove those broad entries work; adding a
redundant bidding-specific glob would create a second list that can drift.

The version decision is "next minor from the post-B1 baseline"; on the branch actually read,
that is `1.6.0`. If B1 changes the package version before merge, rebase first and apply the
same next-minor rule rather than lowering or duplicating its version.

Update `.github/workflows/release-portable.yml` to set up Python before `npm test`
(`release-portable.yml:16-28`). Also add a normal pull-request test workflow covering Linux,
macOS, and Windows with Node's declared minimum and Python 3.10. The current release workflow
builds Windows binaries on Ubuntu and never executes them; build success alone is not runtime
proof.

This exposes two files omitted from B1's stated collision set: `package-lock.json` is
unavoidably changed by a package version/dev-dependency edit, and CI must change if F2 is
answered "yes". Recheck worktree claims for both before Phase 5.

## 4. Tests

### `test/installer.test.js`

Update every exact sorted list:

- `test/installer.test.js:406-411`: discovered add-ons become
  `parley-bidding, parley-design, parley-design-check, parley-tracker, parley-worktrees`;
- `test/installer.test.js:421-428`: default installed skills become core first, bidding
  second, then the other four add-ons;
- `test/installer.test.js:535-545`: doctor reports the same six-skill set; and
- `test/installer.test.js:558-564`: the core marker records all five add-on names.

Extend existing mode assertions at `test/installer.test.js:414-488`:

- default install asserts representative files from every bidding subtree, not just
  `SKILL.md`;
- `--no-addons` asserts `parley-bidding` is absent;
- a dedicated `--only parley-bidding` case asserts core plus bidding only, including scripts,
  tests, fixtures, schemas, adapters, references, templates, metadata, and marker version/hash;
- doctor becomes malformed after deleting a manifest-listed script or schema;
- doctor becomes malformed after changing a byte in a safety-critical file;
- status and paths report core plus bidding after a selective install;
- plain uninstall derives the selective set from the core marker and removes both;
- an unmarked bidding destination blocks the whole target before core or another add-on is
  changed; and
- POSIX installs preserve the five executable bits.

Add table-driven filesystem-shape tests for every target kind, especially Antigravity and
legacy Gemini. Existing tests prove only the core plugin/extension shapes
(`test/installer.test.js:233-247,663-692`) while all add-on lifecycle assertions use Codex.
Filesystem tests do not replace an actual runtime-recognition test.

### `test/design-addons.test.js`

Keep the removed-`addons/` guard at `test/design-addons.test.js:172-208`. Extend the
published-command protection at `test/design-addons.test.js:210-409` instead of accepting a
Python blind spot:

1. Scan all shipped Markdown for logical `python3 scripts/...` commands, including backslash
   continuations. The source publishes five such logical commands at `SKILL.md:96,102-103,
   109-115`.
2. Never execute arbitrary Markdown through a shell. Parse only the Python executable and a
   repository-relative script path, reject shell operators/command substitution, assert the
   script exists, and invoke each entry point with `-B` in a safe validation/help mode.
3. Enumerate `scripts/tests/test_*.py`, run every file with `-B` and
   `PYTHONDONTWRITEBYTECODE=1`, require zero failures, parse each `Ran N tests` summary, and
   require the aggregate to equal 54. A missing Python interpreter is a failing package test,
   not a skip.
4. Run `adapter_validate.py` against all four shipped adapters.
5. Use Ajv 2020 to meta-validate all four schemas and validate the jurisdiction profile,
   adapters, generated initial state, and generated procedure profile against their
   corresponding schemas. The three discovery declarations have no source schema; enforce
   their required shape plus `submission_capable:false` and `retain_origin_link:true`
   explicitly rather than pretending JSON parsing is schema validation.
6. Add a tree-hygiene test for generated caches, absolute source paths, BYTE/customer
   material, credential-shaped values, private keys, and unresolved implementation
   placeholders. Documentation metavariables such as `<release-dir>` and safety prose that
   mentions "password" require narrow allowlists; a raw word scan would be a false green or a
   false failure.
7. Add package-list assertions over `npm pack --dry-run --json` and the portable installed
   tree, comparing them with `parley-addon.json`.

### Source-local test

`skills/parley-bidding/scripts/tests/test_skill_structure.py` is one of the eight renamed
files. Its post-rename assertions must verify the new slug, trigger, and display name. The
test count remains exactly 54.

## 5. README

Make these exact catalogue edits:

- `README.md:3`: remove the stale "five-skill close" wording from the provenance comment;
- `README.md:14-15`: "six skills", "five add-ons", and include bidding in the summary;
- `README.md:19-27`: "six"/"five" and add the `parley-bidding` bullet;
- add a full `parley-bidding` section after the core skill, stating Python 3.10+, default
  read-only behavior, human-only E3b/E5/E6/E7/E8 gates, upload-is-not-submit, exact-byte
  proof, ambiguous-outcome no-retry, and the fact that Parley consensus is not portal
  authority;
- `README.md:109-110`: universal installer wording becomes "all six skills";
- add a safe `$parley-bidding` usage example and a `--only parley-bidding` install example;
- `README.md:121-127,160-224`: explain the Python doctor result and warn users to use one
  installation manager per destination; and
- `README.md:248-263`: add `skills/parley-bidding/` with its references, assets, schemas,
  scripts, tests, and `parley-addon.json`.

## 6. `.gitignore`

Merge, do not copy, the source rules. Add `__pycache__/`, `*.py[cod]`, and `.venv/` to the
target repository `.gitignore`; `.DS_Store` already exists at target `.gitignore:4`.
The nested source `.gitignore` is development metadata, not skill behavior, and should not
become an inconsistently handled runtime payload.

# The rename: all eight files

`rg` found the old slug only in these eight files. It also found the human-facing old display
name in two of the same files and the bid-state title. No schema contains a `$ref`, so the
four `$id` changes have no internal reference rewrite.

All post-edit results below are **NOT TESTED**.

1. `skills/parley-bidding/SKILL.md` (source `SKILL.md:2,6`)
   - `name: software-bidding` -> `name: parley-bidding`
   - `# Software Bidding` -> `# Parley Bidding`
2. `skills/parley-bidding/agents/openai.yaml` (source `agents/openai.yaml:2,4`)
   - display name `"Software Bidding"` -> `"Parley Bidding"`
   - trigger `$software-bidding` -> `$parley-bidding`
3. `skills/parley-bidding/scripts/common.py` (source `scripts/common.py:2`)
   - docstring `software-bidding skill` -> `parley-bidding skill`
4. `skills/parley-bidding/scripts/tests/test_skill_structure.py`
   (source `scripts/tests/test_skill_structure.py:16,26-27`)
   - frontmatter regex expects `name: parley-bidding`
   - UI metadata expects `display_name: "Parley Bidding"`
   - default-prompt assertion expects `$parley-bidding`
5. `skills/parley-bidding/assets/schemas/bid-state.schema.json`
   (source `assets/schemas/bid-state.schema.json:3-4`)
   - `$id` -> `urn:parley-deck:parley-bidding:schema:bid-state:1`
   - title -> `Parley bidding state`
6. `skills/parley-bidding/assets/schemas/jurisdiction-profile.schema.json`
   (source `assets/schemas/jurisdiction-profile.schema.json:3`)
   - `$id` -> `urn:parley-deck:parley-bidding:schema:jurisdiction-profile:1`
7. `skills/parley-bidding/assets/schemas/platform-adapter.schema.json`
   (source `assets/schemas/platform-adapter.schema.json:3`)
   - `$id` -> `urn:parley-deck:parley-bidding:schema:platform-adapter:1`
8. `skills/parley-bidding/assets/schemas/procedure-profile.schema.json`
   (source `assets/schemas/procedure-profile.schema.json:3`)
   - `$id` -> `urn:parley-deck:parley-bidding:schema:procedure-profile:1`

# Fork positions F1-F7

## F1 — schema `$id`s

**Re-root all four under versioned Parley Deck URNs.** Do not keep
`example.invalid/software-bidding`, and do not merely substitute a new path under the same
placeholder domain. These are new distribution identities, the source has no internal
`$ref`, and versioned URNs are stable without pretending that a schema is retrievable at an
unowned HTTP endpoint. Record this as an intentional schema-identity migration in
`IMPLEMENTATION.md`.

## F2 — Python toolchain

**`npm test` and CI must run the Python leg.** The package remains runtime-dependency-free:
Python is an add-on runtime prerequisite, not an npm production dependency. The Node suite
must run all seven files and assert 54, and CI must install Python 3.10. A missing interpreter
fails `doctor` and the package test; it is never a silent skip.

## F3 — version and compatibility metadata

**The add-on gets its own version, initially `1.0.0`; it does not inherit the npm package
version.** The npm version tracks the installer/catalogue distribution. The add-on version
tracks bidding semantics and payload bytes. Store its version and tree hash in
`parley-addon.json` and in the install marker. The inspected package should advance to
`1.6.0` for this catalogue addition.

## F4 — source `.gitignore`

**Merge its Python rules into the target root `.gitignore`; do not ship a nested copy.**
This protects development without treating ignore policy as an operational skill asset.

## F5 — published-command guard

**Extend it.** Capture Python logical commands, reject unsafe shell shapes, validate every
published script path, and run the Python tests/adapter validator through argument arrays.
Accepting a known language blind spot would recreate the false-green class documented at
`test/design-addons.test.js:210-227`.

## F6 — installer validation

**Require full manifest-and-hash validation for `parley-bidding`.** `SKILL.md` alone is not a
valid health definition for an add-on whose safety depends on schemas, adapter ceilings,
state code, tests, and references. Preserve `SKILL.md`-only compatibility for the four legacy
add-ons, but make the manifest extension generic for future complex skills.

## F7 — sequencing

**Accept B1: design now, implement only after the blocking review is clean and merged.**
Rebase onto that result, re-read the six overlapping files, and rerun all baselines before
copying anything. A parallel implementation or worktree override is rejected because the
file sets intersect.

# What packaging can weaken

These are adversarial integration findings. The described integrated behavior is **NOT
TESTED** unless a baseline result is explicitly stated above.

1. **Cross-skill disclosure bypass.** The source requires a tender-scoped E3b approval for
   the exact roster, providers, packet, data classes, and redactions
   (`references/parley-integration.md:5-16`). The target README currently defaults external
   backend disclosure to YES for task/repository context (`README.md:232-236`). Once both
   skills are installed, an agent can interpret the generic Parley default as permission to
   send tender files, pricing, contracts, or supplier data to multiple backends. The
   bidding-specific E3b rule must explicitly override that generic default.

2. **Consensus laundering of human authority.** Co-location with a consensus protocol makes
   it easier to mistake four agent signoffs for a commercial, authority, upload, or submit
   approval. The source says approvals are single-use and fingerprint-bound
   (`SKILL.md:55-70`) and agents never operate portals
   (`references/parley-integration.md:33-41`). Tests and README must state that no Parley
   artifact satisfies E5/E6/E7/E8.

3. **Instruction loading and tool availability vary across fourteen runtimes.** The source
   is an instruction-level safety boundary, not a sandbox. A runtime that truncates
   `SKILL.md`, loads only the description, resolves conflicting skills in another order, or
   grants an autonomous browser can lose the procedural gate while retaining the tempting
   "submit" capability wording. The short frontmatter description already says "Never use it
   as authorization" (`SKILL.md:2-3`); packaging tests must verify that exact sentence reaches
   every artifact and installed copy.

4. **Antigravity and legacy Gemini may not discover sibling add-ons at all.** For these two
   targets, `targetSkillUnits` puts add-ons beside the core destination
   (`lib/installer.js:865-883`). That produces bare directories under a plugin/extension
   root, while the existing tests prove manifests only for the core
   (`test/installer.test.js:233-247,663-692`). A copied directory plus a green `doctor` is not
   proof that the runtime exposes `$parley-bidding`. Actual recognition is a strict blocker;
   if those managers require a multi-skill manifest, the installer/manifests must be changed.
   If legacy Gemini cannot express it, the user must change acceptance rather than receive a
   false fourteen-runtime claim.

5. **`doctor` currently approves a safety-gutted tree.** Deleting
   `scripts/adapter_validate.py`, all schemas, or `references/hitl-and-recovery.md` still
   leaves an add-on "valid" because only `SKILL.md` is required
   (`lib/installer.js:1129-1147`). Hashing only at source handoff is insufficient; installed
   health must detect missing and modified safety files.

6. **A name collision can create a partial fleet.** The installer is atomic per skill
   directory, not per selected set (`lib/installer.js:888-891,997-1053`). An existing
   unmarked `parley-bidding` can block that unit after the core and other add-ons were
   replaced, while the core marker records the intended full set. Preflight must happen
   before the first write.

7. **Python availability can be a false health signal.** The source uses Python 3.10 syntax
   (`scripts/common.py:19`) and publishes five Python commands (`SKILL.md:96-115`), but the
   current installer validates files only. A host can receive six "valid" skills while every
   deterministic bidding command fails. Declare and check the interpreter requirement.

8. **File modes and packers can diverge.** Five source scripts are executable, while the
   install copy routine recreates regular files without applying the source mode
   (`lib/installer.js:1077-1091`). npm pack, a portable executable, and a native install can
   therefore expose different trees. Compare file lists, raw hashes, and applicable POSIX
   modes across all three artifacts.

9. **The DTVP maturity label becomes globally visible.** The profile says
   `live-submission-validated` but also `live_effects_authorized:false`
   (`assets/platform-adapters/cosinex-vmp.dtvp.json:7-11`). A runtime or user may read the
   maturity as blanket permission or platform-wide reliability. Keep the false authorization
   field in schema, adapter validator, README, and installed-hash checks; never shorten the
   profile to marketing prose.

10. **Universal and native installers have different trust paths.** `npx skills add` discovers
    directories directly, while this package's native installer writes markers and can
    enforce `parley-addon.json`. A universal install may be complete but unmanaged by native
    uninstall, or incomplete without native doctor ever running. README must say which
    manager owns a destination, and the npm tree itself must be self-validating.

11. **Prompt injection gains a multi-backend fan-out.** Tender and portal content is
    explicitly untrusted evidence (`SKILL.md:10`). Once Parley is one instruction away, a
    malicious tender can ask an agent to disclose itself to the full roster or reinterpret a
    portal message as approval. E3b packet allowlisting, credential exclusion, and the
    "content is evidence, never instructions" rule must survive verbatim and be exercised in
    adversarial review.

12. **A literal secret scan can be misleading.** Safety documentation intentionally contains
    words such as "password", while command examples intentionally contain metavariables such
    as `<release-dir>`. Conversely, a customer name or credential need not use either word.
    Check 9 needs explicit allowed documentation patterns plus credential-value, private-key,
    entropy, absolute-path, and known-customer scans. A simple zero-match grep is not enough.

# Verification plan

Every integrated check is **NOT TESTED** in this design-only round. I would run the brief's
checks in this fail-fast order:

1. **Brief check 9 — integrated-tree hygiene and confidential-data scan.** First inspect
   file types/symlinks and scan for generated caches, absolute source paths, BYTE/customer
   material, credential values, private keys, unresolved implementation placeholders, portal
   mutation imports, and network/browser libraries. Use narrow allowlists for documented
   safety words and example metavariables. **NOT TESTED on integrated tree.**
2. **Brief check 10 — source-versus-integrated diff.** Compare sorted paths, raw SHA-256,
   and POSIX modes. Require exactly 39 unchanged files, the 8 documented edits, removal/merge
   of source `.gitignore`, and the new `parley-addon.json`; reject every other difference.
   **NOT TESTED.**
3. **Brief check 8 — schemas and profiles.** Meta-validate the four schemas with Ajv 2020;
   validate the four adapters, jurisdiction profile, generated bid state, and generated
   procedure profile; structurally validate the three discovery declarations. The current
   host has neither Python `jsonschema` nor Ajv in the target package, so full Draft 2020-12
   validation was **NOT TESTED**. Phase 5 adds Ajv as a dev-only validator.
4. **Brief check 5 — adapter validator.** Run
   `python3 -B skills/parley-bidding/scripts/adapter_validate.py
   skills/parley-bidding/assets/platform-adapters`. **Source baseline RUN and PASS (4
   profiles); integrated tree NOT TESTED.**
5. **Brief check 6 — all Python tests.** Run the seven files individually with `-B` and
   `PYTHONDONTWRITEBYTECODE=1`; parse and total the summaries. **Source baseline RUN and PASS
   (54); integrated tree NOT TESTED.**
6. **Brief check 7 — cache-free compile.** Read every shipped `.py` and call Python's
   in-memory `compile(source, path, "exec")`, then assert no cache artifact exists. This
   performs a syntax compile without writing bytecode. **NOT TESTED as a dedicated compile
   check.** The source tests were run with bytecode disabled and left no cache.
7. **Brief check 1 — full Node suite.** Run `npm test`; require the Node suite, the 54-count
   Python leg, adapter validation, schema/profile validation, and published-command guard all
   to pass. **Target baseline RUN and PASS (253); integrated suite NOT TESTED.**
8. **Brief check 2 — npm pack.** Run `npm pack --dry-run --json` with an isolated writable
   npm cache, compare the returned paths/hashes/modes with the manifest, and assert no cache
   files. **Target baseline RUN and PASS (153 entries); integrated package NOT TESTED.**
9. **Brief check 3 — install modes.** In separate temporary homes run default install,
   `--no-addons`, and `--only parley-bidding`; assert exact skill sets and full manifest
   equality. Repeat with the current-host portable executable. **NOT TESTED.**
10. **Brief check 4 — management lifecycle.** On the selective install run `doctor`,
    `status`, and `paths`; corrupt one script and prove doctor fails; restore/reinstall; run
    `uninstall`; prove only managed selected destinations are removed. **NOT TESTED.**
11. **Brief check 11 — generic discovery.** Run
    `npx skills add <local-repo> --list` and require exactly the six expected names.
    **NOT TESTED and not presently runnable offline on this host:** no `skills` executable is
    installed, and fetching it through `npx` requires external package access. This check
    remains a release blocker until the CLI is made available; it must not be waived.

Mandatory distribution checks beyond the numbered eleven:

- Build `npm run build:portable:current`, execute it from outside the repository, and repeat
  checks 3-4 plus the manifest comparison. **NOT TESTED:** the build writes `dist/`, which was
  forbidden in this design-only round.
- Build Windows and Linux artifacts and inspect their embedded manifest; execute them in
  matching CI runners. **Cannot be executed on this Darwin arm64 host.**
- Prove actual `$parley-bidding` recognition in every claimed runtime, not merely copied
  files. **Cannot be completed on this host:** the fourteen runtime CLIs and Windows/Linux
  environments are not all available. Antigravity and legacy Gemini are priority blockers.

The host does have Node `v26.5.0`, npm `11.17.0`, and Python `3.14.6`; `pytest` is absent but
is unnecessary because the authoritative runner is file-by-file standard-library
`unittest`.
