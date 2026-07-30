---
idea: integrate-parley-bidding-addon
status: fix-up-cycle-5
implementer: claude-1
started: 2026-07-30
completed: n/a
branch: parley-deck-skill#integrate-parley-bidding-addon
head-commit: 3634cc8
design-pr: n/a
implementation-pr: n/a
---

## Summary of work

`parley-bidding` is added as the sixth skill in `parley-deck-skill`, together with the
payload-integrity mechanism that makes shipping a 47-file security-relevant payload
defensible. Both in one change, per `FINAL.md` — the user was offered the option of deferring
the integrity mechanism and declined it.

The defect the integrity work exists against, stated once: `validateInstalledPayload` asked a
single question of an add-on directory — is `SKILL.md` there? — so a tree gutted to that one
file reported `valid`. `lib/installer.js:1132` was that line.

## Preconditions (F7) — measured, not assumed

| precondition | evidence |
|---|---|
| zero agreed fixes | `review/consensus.md` for `skills-cli-install-path`: zero agreed fixes, four deferred follow-ups, three ✅ ACCEPT signoffs |
| merged | `parley-deck-skill` `main` at `a544dcd` (`release: parley-deck-skill 2.0.0`), pushed |
| released | npm `2.0.0`, GitHub `v2.0.0`, Homebrew `2.0.0`, two WinGet PRs open |
| rebase | this branch is cut from `a544dcd`, so F7's rebase is a no-op — recorded rather than performed |
| six overlapping files re-read | `lib/installer.js`, `package.json`, `test/installer.test.js`, `test/design-addons.test.js`, `README.md`, `skills/` layout — surveyed on the merged tree before any copy |

### Baselines captured before the first write

- `npm test` → **253 pass / 0 fail**
- source inventory → **48 files**, 246 KB, **zero** symlinks, **zero** `__pycache__`/`.pyc`
- Python baseline → **54 tests**, split **4+20+2+3+15+3+7**, all OK, run against a **copy**
  with `PYTHONDONTWRITEBYTECODE=1 python3 -B`, leaving zero cache artefacts in copy or source

## What was built

### The payload (T1–T4)

- `.gitignore` merged **before** any copy (F4). Three of the source's four rules were new;
  `.DS_Store` was already present.
- 47 files copied to `skills/parley-bidding/`. The nested `.gitignore` was dropped per F4.1 —
  `copyRecursive` filters nothing, so it would otherwise have landed in every runtime.
- Rename `software-bidding` → `parley-bidding`: **13 occurrences across 8 files**. The
  contract said twelve; the thirteenth is real and is recorded here rather than silently
  absorbed. Schema `$id` hosts stay `example.invalid` per F1.2.
- The ratified B1/B2 consent paragraph inserted byte-exact into **both** `SKILL.md` (after the
  E0–E8 ladder) and `references/parley-integration.md`.

**Placement note, not a deviation.** `FINAL.md` says the paragraph goes "beside E3b". It sits
immediately after the full effect ladder rather than interrupting it between the `E3b` and
`E4` bullets. The sentence names `E3b, E5, E6, E7 or E8`; placing it beside `E3b` would
forward-reference four bullets the reader has not yet reached. Text unchanged, verified
byte-exact against the `FINAL.md` blockquote by unwrapping it programmatically rather than by
re-typing and comparing against my own typing.

### The integrity mechanism (T5–T6)

- `lib/addon-manifest.js` — a new shared module. One definition of the digest, used by both
  the generator and the installer, so the two cannot drift.
- `parley-addon.json` — **generated, never hand-written**: 47 entries, excludes itself,
  aggregate digest, `runtime.python: ">=3.10"`, and no second version number (F3).
- `scripts/build-addon-manifest.js` — generates and, with `--check`, verifies. Generic: with
  no arguments it refreshes only add-ons that already carry a manifest, so it never decides on
  its own that an add-on ought to have one.
- `lib/installer.js`:
  - `MARKER_SCHEMA = 2`, and the marker now records `manifest: {aggregate, sha256}` — or
    `false` when the source genuinely shipped none.
  - `preflightSkillUnit` validates **every** unit and destination before the first write;
    `installTarget` aborts the whole target on any blocker (B5).
  - staged bytes are re-verified before the marker is written, then the complete staged unit
    is validated, then the atomic replace happens — codex-1's ordering, unchanged.
  - `manifestProblems` is the marker-anchored check; `doctor`/`status` print `integrity:` lines.

### Verification, packaging, documentation (T7–T15)

- `scripts/run-python-tests.js` — seven files individually, `python3 -B` with
  `PYTHONDONTWRITEBYTECODE=1`, per-file counts asserted, a stray-bytecode scan, and a hard
  **failure** when no interpreter is present. Measured: with `python3` off `PATH` it exits 1.
- `package.json` — `npm test` now chains node → Python → manifest check. `prepack` carries the
  package preflight and was confirmed to fire under `npm pack --dry-run`.
- `.github/workflows/release-portable.yml` — `setup-python` pinned to the declared floor.
- `.github/workflows/test.yml` — **new**. The only workflow before this ran on `release:
  published`, so the Python leg would never have run until after a release existed. Runs the
  suite on every push and PR across Python 3.10 and 3.13.
- `test/design-addons.test.js` — the published-command guard **generalized**, not duplicated:
  the extractor is now parameterized on a `{binary, flag}` shape, with `NODE_COMMAND` and
  `PYTHON_COMMAND` arms. All 253 pre-existing assertions stayed green through the refactor,
  which is the evidence that the node arm is unchanged.
- `test/bidding-addon.test.js` — **new**, 23 tests in cycle 0 and **31** after fix-up cycle 1,
  covering every negative case codex-1 attached to its amendment vote.
- `README.md` — six skills, the add-on's own section, and the default-install availability
  duty stated as a callout.

## Deviations from FINAL.md

### D-1 — B3.11 contradicts B3.13; resolved by a third option, and ratified separately

`FINAL.md` requires the manifest to be both "generic and optional" (B3.11) and "required for
`parley-bidding`" (B3.13). Those cannot both hold: name-keying re-introduces the hardcoded
add-on registry `discoverAddons` exists to avoid, and presence-keying makes B3.13's own
acceptance test unimplementable, since deleting the manifest would silently downgrade the
add-on to `SKILL.md`-only validity.

Escalated to the user
(`inbox/claude-to-user_integrate-parley-bidding-addon_manifest-requirement-fork.md`). The user
chose **marker-anchored** — a third option not present in `FINAL.md`.

Because that design was never reviewed, it was **not** treated as satisfying B3 on the user
ruling alone. A targeted amendment round (`round-03/`) put it to the three participants. All
three returned `Amendment: ACCEPT WITH CONDITIONS`. codex-1's four conditions are implemented
as listed above: versioned marker schema, two stored values rather than a boolean, the
validate → stage → revalidate → mark → validate → replace ordering, and the negative tests.

**The boundary codex-1 drew, documented rather than implied.** The marker cannot detect a
manifest omitted *before the first install ever observed the source*. B3 is therefore only
satisfied in combination with the release-time inventory gate (B7) below. And the whole
mechanism is **defect detection after a validated install, not tamper resistance**: anyone who
can rewrite the payload can rewrite the marker beside it. `lib/addon-manifest.js` says so at
the top of the file.

### D-2 — PRE-7's proving artefact replaced

`git status --porcelain software-bidding` returns `?? software-bidding/` in the parent repo —
the source is *untracked* there, so the check can never be empty and proves nothing. Replaced
with a byte-level check: file count plus a cache scan before and after. Deviation from the
stated method, not from its intent. **Result: source still 48 files, 0 caches, untouched.**

### D-3 — the Python command grammar accepts `<` and `>`

Open question F5.3/F5.5, now closed. `SUPPORTED_COMMAND` refuses `<`/`>` because it hands the
string to `/bin/sh`, where they redirect. The Python arm is **static** — F5 says so — and
hands nothing to a shell, so those characters carry no hazard, while all five published Python
commands legitimately carry `<placeholder>` arguments. `;` `|` `&` backtick `$` and `\` are
still refused, because they mean the published line is not one self-contained command.

### D-4 — B4.3 resolved as NOT TESTED, and the claim dropped

Only five of the fourteen runtime CLIs exist on this machine (`codex`, `claude`, `agy`,
`hermes`, `kimi`). An isolated-`HOME` probe of Claude Code could not proceed either — the
isolated config carries no credentials. Per B4 the choice is prove per-target recognition or
stop claiming the target, so **the claim is dropped**: the README states that the payload
installs into and validates in all fourteen destinations (measured), and that whether a
runtime then *exposes* it as an invocable skill is **NOT TESTED**.

### D-5 — the D-10 index target is `README.md`

The contract's wording is that *documentation* says "six skills: the core protocol plus five
add-ons". `README.md` now says exactly that, in three places. No add-on index was added to
`skills/parley-deck/SKILL.md`: the core skill is distributed standalone, and a list of its
siblings inside it would couple the protocol to this package's contents.

## Verification — every item run, none asserted

| check | result |
|---|---|
| `npm test` (node leg) | **309 pass / 0 fail**, measured on python3 **3.9.6 and 3.14** (253 before this idea; 278 at cycle 0) |
| `npm test` (Python leg) | **54 pass / 0 fail on supported interpreters only.** On 3.9.6 the leg **refuses to run** — `python3 is 3.9, but the add-on declares >=3.10` — and zero Python tests execute through it. That is the F2 contract working, not a failure. Run file-by-file, 3.9.6 does pass all 54; that is a *compatibility* observation, not the package gate. |
| Python leg | **54 tests**, 4+20+2+3+15+3+7, 0 fail |
| Python floor `>=3.10` | **measured** on 3.10, 3.11 and 3.14 — 54/54 on each |
| Python absent | runner exits **1** with a failure message, does not skip |
| `npm pack --dry-run` | 202 files; **48** under `skills/parley-bidding/`; **no** `__pycache__`/`.pyc`; no nested `.gitignore`; `prepack` fires |
| default install / `--no-addons` / `--only parley-bidding` | all three behave; covered by tests |
| `doctor` / `status` / `paths` / `uninstall` | all four exercised against the new add-on |
| adapter validator | all **4** shipped adapters OK |
| shipped JSON | **16** files parse; **4** schemas carry `$schema` 2020-12 and `example.invalid` `$id`s |
| B3 negative test | gutted tree, four single-file deletions, one flipped byte, stray `.pyc`, deleted manifest, corrupt manifest, stripped marker field, newer marker schema, and a **self-consistent manifest+payload swap** — every one reports `malformed` |
| B5 zero-writes | a corrupt source payload leaves the destination **non-existent**, not half-written |
| B7 cross-channel | identical aggregate `sha256:7854adf1…` across **repository, npm tarball, portable binary install, native install** |
| source vs integrated | 1 path dropped (`.gitignore`), 1 added (`parley-addon.json`), **9** files differ in content — 8 rename, 2 consent paragraph — each listed above |
| secret / identity scan | no credentials, no `BYTE`, no customer data; the only email-shaped strings are `operator@example.invalid` in test fixtures |
| cache / symlink scan | zero in repository, tarball, and every install |
| `npx skills add … --list` | **"Found 6 skills"**, `parley-bidding` among them |
| portable binary | builds; installs 49 files (47 payload + manifest + marker); `doctor` valid |

## Release implications

- **Version.** This adds a skill and new installer validation. That is a feature, so the
  correct number is **2.1.0**, not a patch. Existing install modes are unchanged and
  backward compatible; a marker written by 2.0.0 has no `markerSchema` and is treated as
  legacy rather than broken, so upgrading in place does not report a healthy install as
  malformed.
- **Availability, not permission.** A routine `install --force` upgrade places a
  procurement-portal skill into every covered runtime, including for users who never asked for
  one. Every gate still binds. This is stated in `README.md` as a callout, per the
  documentation duty.
- **Channels.** npm, GitHub release, Homebrew formula, WinGet manifests — the same four as
  2.0.0. WinGet PRs #409827 and #409828 from the 2.0.0 release are still open upstream; a
  2.1.0 submission stacks on top of them rather than waiting.
- **`FINAL.md` carries a hard stop before publishing.** Phase 5 ends here, with the diff and
  this evidence. The user has since asked for the release explicitly; it therefore proceeds
  **after** Phases 6–8, not in place of them.

## Fix-up cycle 1 — review round 1

Three reviewers, three **BLOCK** verdicts, on the same two defects. `agy-1` could not
participate (account quota exhausted, ~6 days to reset) and is recorded as an outage, not an
accept. Both defects were requirements the amendment round had already ratified and that I did
not implement — this is a compliance failure on my side, not a disagreement about design.

### F1 — an expected unit with no readable marker reported `valid`

`codex-1` MAJOR, `hermes-1` CRITICAL, `kimi-1` CRITICAL. All three measured the same matrix:

```json
{"noPython":        {"doctorOk": true, "status": "valid", "problems": []},
 "missingMarker":   {"doctorOk": true, "status": "valid", "problems": []},
 "unreadableMarker":{"doctorOk": true, "status": "valid", "problems": []}}
```

`codex-1`'s ratified condition 1 ends: *"An expected installed unit with a missing or
unreadable marker must also be unhealthy."* `hermes-1`'s round-03 finding 1 named the same
double-deletion path. I implemented neither half.

Why it mattered more than it looks: every negative test in cycle 0 preserved the marker, so the
suite was green while the cheapest real-world gutting — copy `SKILL.md` alone, or an extraction
that stops early — produced no marker at all and reported `valid` with exit 0.

**Fixed** in `skillUnitStatus`, not in the add-on-specific path, so every add-on gets it.
`readMarkerState` now distinguishes absent from unreadable, which `kimi-1` asked for
explicitly. Measured after the fix: gut+delete-marker, missing marker, and unreadable marker
all report `malformed` with `doctorOk: false`; an untouched install is unchanged.

### F2 — B6 was enforced in the test runner, not in `doctor`

`codex-1` MAJOR, `hermes-1` CRITICAL, `kimi-1` MAJOR. B6 assigns the interpreter check to
`doctor`; I had put it only in `scripts/run-python-tests.js`, which never runs on a real
install. `kimi-1` called it precisely: a placement error, not a missing feature.

**Fixed.** `doctor` reads `runtime.python` from the validated manifest, probes `python3` once
per process, and fails health when it is absent or below the floor. Per B6 the two answers stay
separate — the payload is still `valid`, the unit is reported `unavailable`, and `doctor` exits
1. Add-ons that declare no runtime are unaffected (`runtime: null`). Measured with a stubbed
interpreter at 3.9 (fails), at 3.10 (passes), and with none at all (fails).

### F3 — `head-commit` front matter was stale

`kimi-1` MINOR. It recorded `a544dcd`, which is `main`. Corrected to the fix-up head.

### kimi-1's open question 4

`CHANGELOG.md` and `RELEASING.md` belong to **this** idea — they are T14 of the plan, written
while the reviewers were running and uncommitted when `kimi-1` looked. They are in the fix-up
commit.

### A consequence of F1's fix that I am NOT resolving unilaterally — round-2 reviewers, please rule

`codex-1`'s condition is unconditional: *"An expected installed unit with a missing or
unreadable marker must also be unhealthy."* I implemented it as written. It has a cost I found
by measuring, and it is a **false alarm on a documented install path**:

```text
$ npx -y skills add feci/parley-deck-skill --skill parley-bidding --agent claude-code
$ parley-deck-skill doctor --target claude
claude/parley-bidding: malformed …/.claude/skills/parley-bidding
  integrity: no parley-deck-skill install marker: this directory was not installed by this
             tool, or the marker was removed
exit=1
```

Measured: the universal `skills` CLI — the install path this repository's own README
recommends first — copies the tree faithfully **including `parley-addon.json`**, and writes no
marker of ours. `verifyPayload` on that directory returns `ok: true`: the payload is *provably*
byte-intact, and `doctor` calls it malformed anyway.

So the strict rule trades a real false positive for its true positive. A refinement exists:

> When the marker is absent but the tree carries a manifest that fully verifies, the integrity
> concern the condition protects is satisfied by **stronger** evidence than the marker — report
> the unit as valid-but-not-installed-by-this-tool rather than malformed. When the source
> add-on ships a manifest and the installed tree has neither manifest nor marker, that is still
> the gutting signal and stays malformed. This stays generic: whether the source ships a
> manifest is read from `unit.addon.root`, so no add-on is named in code.

I am not applying that on my own judgement, because it narrows a condition three participants
ratified, and the last time I substituted my own reading for a ratified condition it produced
the two BLOCKs above. Round-2 verdicts should say which they want:

- **(a)** keep it unconditional as shipped — a foreign-installer tree is unhealthy to this tool,
  and the false alarm is the honest answer;
- **(b)** take the refinement above;
- **(c)** something else.

Whatever is chosen lands in fix-up cycle 2 with tests either way.

### Test delta

278 → **286** node tests. Eight new: missing marker, unreadable marker, a marker belonging to
another tool, the double deletion, interpreter absent, interpreter below floor, interpreter at
floor, and an add-on with no runtime declaration. The interpreter tests run in a child process
with a controlled `PATH`, because the probe is cached per process and reads the real
environment.

## Fix-up cycle 2 — review round 2

`codex-1` **BLOCK** (1 MAJOR, 1 MINOR), `hermes-1` **BLOCK** (1 MAJOR), `kimi-1` **ACCEPT**
(1 NIT). Both round-1 findings were confirmed closed by all three, each from its own
re-measurement rather than from reading the diff.

### F4 — the interpreter probe answered for the wrong environment

`codex-1` MAJOR. `probePython3` called `spawnSync` with no `env`, so it searched the parent
process's `PATH`. The installer is also a library: `run(argv, io)` takes an effective
environment, and every other part of the installer honours `context.env`. codex measured a
direct false green — parent had Python 3.14, caller declared `PATH: ""`, `doctor` exited 0
with `runtime.ok: true`. Checking the wrong environment does not answer an
operational-availability question at all, so B6 was not actually closed.

**Fixed.** The probe resolves against `context.env` and is memoized **per PATH** rather than
once per process, so two contexts in one process get two answers. Added a 5s timeout: a health
check must not hang on a `PATH` entry sitting on a stalled mount.

### F5 — my "286 pass / 0 fail" was environment-dependent

`hermes-1` MAJOR, and it corrects a claim in this file rather than only the code. The
pre-existing test `doctor reports add-on skills per target` asserts `result.ok`, and `ok` now
includes availability. On this machine `python3` is Homebrew 3.14, so it passed. On macOS's
default `/usr/bin/python3` (3.9.6) it fails — hermes measured 285 pass / 1 fail there. The
evidence I published was true of my machine, not of the change.

**Fixed** the way hermes proposed: the assertion now checks the skill *list*, which is what
that test is about; health is covered by the tests written for it. One of my own new tests had
the same defect (it assumed the ambient interpreter met the floor) and now stubs both arms.

**Re-measured on both interpreters: the node leg passed on `python3` 3.9.6 and on 3.14.** (The Python leg refuses 3.9.6 by design; see the evidence table.)

### F6 — runtime probing leaked into `paths`, and `status` discarded it

`codex-1` MINOR. `skillUnitStatus` probed unconditionally, so `paths` — a path-discovery
command — executed a `PATH`-resolved program, with no timeout. Meanwhile `status` paid for the
probe and printed only `valid`, so the two commands disagreed about the same directory.

**Fixed.** Probing is now an explicit option of the traversal: `doctor` and `status` probe,
`paths` does not. `status` renders the `unavailable:` line.

### F7 — `kimi-1` NIT

"One or more installs are **installed but** operationally unavailable" — subject and reason
collided. Reworded.

### Answers to the open questions

- **`kimi-1` Q1 (Windows).** `python3`-only is the intended contract, and it is now stated in
  `CHANGELOG.md` rather than left implicit. It matches how the skill's own published commands
  invoke the interpreter, and it fails safe.
- **`kimi-1` Q2 (the legacy carve-out is invisible).** It is already distinguishable without a
  new field: `doctor --json` carries the parsed `marker` per unit, so a grandfathered unit is
  the one whose `marker.markerSchema` is absent. No code change; recorded here so the property
  is deliberate rather than accidental.
- **`codex-1`'s transient 4-test failure.** Not reproduced here either, in repeated runs on two
  interpreters. One plausible mechanism now removed: the process-global probe cache made
  results depend on which test ran first in a file. It is keyed per PATH as of this cycle.
  Recorded as unexplained rather than closed.

### Process note — a reviewer reset the working tree

`hermes-1` checked out `714712f`'s `lib/installer.js` to compare behaviour, and a `git reset`
landed in the repository under review (visible in `git reflog`). It discarded one uncommitted
edit of mine — the `status` printer — which I re-applied. No committed work was affected and
no reviewer finding depends on it. The lesson is mine: **commit before launching reviewers**,
because a review that runs commands in the tree can mutate it.

## Fix-up cycle 3 — review round 3

`codex-1` **BLOCK** (1 MAJOR, 1 MINOR), `hermes-1` **ACCEPT WITH CONDITIONS** (1 MINOR),
`kimi-1` **ACCEPT WITH CONDITIONS** (2 NITs). All four round-2 findings confirmed closed by
independent re-measurement.

### The ruling I refused to make alone — unanimous (b)

All three chose **(b)**, and `kimi-1` measured a worse consequence than the one I disclosed:
following the README's first recommendation fully (`--skill '*'`, all six skills) produced
**six `malformed` verdicts and exit 1 on six byte-perfect payloads**. Its argument is the one
that settles it — when `doctor`'s verdict contradicts its own strongest evidence, users learn
that `malformed` can mean "perfectly fine", and the word stops meaning anything. That is the
untested-and-green failure this whole mechanism exists to kill, inverted into tested-and-red.

**Implemented exactly as ratified**, with `kimi-1`'s precisions:

| case | verdict |
|---|---|
| marker absent, source ships a manifest, tree's manifest fully verifies | **`valid-unmanaged`**, `managed: false`, `marker: null`, health unaffected, runtime still probed |
| marker present but unreadable | `malformed` — corrupted management metadata, not "never installed here" |
| marker naming another installer | `malformed` |
| marker absent, manifest absent where the source ships one | `malformed` — the gutting signal, preserved |
| legacy `markerSchema`-less marker | unchanged |

A distinct status string rather than `valid` + a flag, so automation can require tool-managed
installs without parsing prose; `managed` carries the same fact as a boolean, which is what
`codex-1` asked for. The rule is generic — it asks `unit.addon.root` what the packaged source
ships, so no add-on is named in installer code.

**Measured after the change**, on the README-first full install: `parley-bidding`
`valid-unmanaged`; the other five still `malformed`.

### The residual, stated rather than hidden

Only `parley-bidding` ships a manifest today, so the other five units have nothing to verify
against and an unmarked copy of them stays `malformed` — `doctor` exits 1 on that path.
`kimi-1` named this as follow-up scope and gave the two ways to close it: ship manifests for
every unit (the honest end state; the mechanism already generalizes), or state the limit in
the docs. **Both were done in part**: the limit is now stated plainly in `CHANGELOG.md` and
referenced from `README.md`; shipping the remaining manifests is **deferred**, because
`FINAL.md` B3.11 holds the other four add-ons unaffected by this change and I will not widen a
ratified boundary in a fix-up cycle. It is the first follow-up in the consensus.

### The other findings

- **`codex-1` MINOR — a broken shim satisfied the floor.** `4.not-a-version` parsed to
  `{major: 4, minor: 0}` and passed `>=3.10`: fail-open on the one check whose job is to fail
  closed. The parse is now anchored `^(\d+)\.(\d+)$` and requires both integers.
- **`kimi-1` NIT — the cache key was narrower than the spawn.** Keyed on `PATH` while
  `spawnSync` received the whole environment, so `PYTHONHOME` and friends could change the
  answer without changing the key. Now keyed on all of them.
- **`kimi-1` NIT — `status` explained add-ons but not the core.** A foreign-installed core is
  exactly the unit whose verdict a user most needs explained. The printers are unified.
- **`hermes-1` MINOR — `status` always exits 0.** Kept: it is the command's original contract
  and hermes did not block. Now documented, with `doctor` named as the health gate.

### `codex-1`'s unexplained transient, revisited

Not reproduced: **6 sequential and 4 concurrent full runs, all 296/0**, on two interpreters.
One plausible mechanism was removed in cycle 2 (the process-global probe cache made a result
depend on which test ran first). Still recorded as unexplained rather than closed.

## Fix-up cycle 4 — review round 4

`codex-1` **BLOCK** (3 MAJOR, 2 MINOR), `hermes-1` **ACCEPT** (1 MINOR), `kimi-1` incomplete —
its process ended with a `PENDING MEASUREMENT` draft after finishing its code read. Its two
named suspects were taken here regardless, and it is invited to round 5 rather than counted as
an accept.

The round-3 ruling opened a door, and `codex-1` walked through it four times. Three of the five
findings are about behaviour my change made **inconsistent** rather than newly broken — health
became strict about ownership while the mutation paths stayed loose — which is exactly the kind
of seam a fix-up cycle should be looking for.

### F8 — `valid-unmanaged` trusted the installed manifest as its own authority

**MAJOR.** `unmanagedButVerified` used the packaged source manifest only as a yes/no capability
check, then verified the installed payload against whichever manifest sat beside it. That
recognizes *any* self-consistent tree, not the packaged one. And because `runtime` lives
outside the payload aggregate, **deleting that single field from the installed manifest
disabled the B6 interpreter check without rehashing a single file** — codex measured
`ok: true`, `valid-unmanaged`, `runtime: null`, no problems, with an empty `PATH`.

**Fixed** by anchoring the proof to the packaged source: the source manifest must itself
verify, the installed manifest's bytes must equal the source's, and the installed payload must
match. A laundered tree is `malformed` again — which also makes the round-3 security argument
stronger than it was when it was accepted.

### F9 — mutations authorized on a pathname, not on ownership

**MAJOR.** Install preflight and `installSkillUnit` treated *any* file at the marker path as
authorization. codex put a `FOREIGN-SENTINEL` in a tree marked `name:"other-installer"`, ran a
plain `install --only parley-bidding` with no `--force`, and it was replaced and the sentinel
deleted. Health had called that same tree `malformed`; the two disagreed.

**Fixed** with one parsed predicate — present, readable, our name, matching skill — used by
health, install and uninstall alike.

### F10 — uninstall could remove the core and then refuse an add-on

**MAJOR.** `uninstall` had no preflight. With a managed core and an unmanaged add-on it removed
the core, then refused the add-on and reported failure. **Fixed**: removal is atomic across the
selected set, the same rule B5 gave installation.

### F11 — the documented opt-out could leave the skill installed and invisible

**MAJOR, and the one that matters most for this particular add-on.** After a universal
install, `install --force --no-addons` wrote only the core and recorded a core-only selection.
The bidding directory stayed on disk, still reachable by the runtime, and **disappeared from
`doctor`'s traversal** — so a green `doctor` was not evidence that the opt-out had taken
effect. The README tells users to reach for exactly that flag to avoid a procurement skill they
did not ask for.

**Fixed**: read-only commands now surface an add-on directory present on disk but absent from
the recorded selection, mark it `selected: false`, and fail health with a message naming the
remedy. Install and uninstall keep the selection-only view — this is a visibility rule, not a
licence to delete directories the user did not select.

### F12/F13 — the two MINORs

- A **directory or dangling symlink** at the marker path read as "entirely absent" and so took
  the `valid-unmanaged` branch, against round 3's precision that only an absent marker
  qualifies. `fileExists` → `lstatSync`; only `ENOENT` is absent.
- The probe cache key enumerated five variables and **joined them with a separator that can
  occur inside a value**. codex showed a `PYENV_VERSION`-selected shim reusing a stale verdict;
  `kimi-1` independently flagged the collision. Now collision-free JSON over the whole
  effective environment.

### `kimi-1`'s second suspect, and `hermes-1`'s MINOR

- `managed` is now present on every unit shape, including `missing`, so the JSON is uniform.
- `status` names missing files, as `doctor` does.

### Measured after the cycle

**309 node tests, 0 fail, on python3 3.9.6 and 3.14; 54 Python tests, 0 fail, on supported interpreters (3.10, 3.11, 3.14) — the leg intentionally refuses 3.9.6.** Six new tests: the
runtime-field-removal probe, a laundered tree, a directory/dangling-symlink marker, the
`--no-addons` and excluding-`--only` trails, atomic uninstall, and a foreign-marked
destination refusing a plain install.

The README-first residual is unchanged and still deferred: `parley-bidding` reads
`valid-unmanaged`, the five skills that ship no manifest read `malformed`.

## Fix-up cycle 5 — review round 5

`codex-1` **BLOCK** (2 MAJOR, 2 MINOR). `hermes-1` and `kimi-1` did not complete this round;
`kimi-1` left a `PENDING MEASUREMENT` draft, `hermes-1`'s process produced no output.

codex opened this round with a provisional `ACCEPT` written before measuring, then flipped it
to `BLOCK` after measuring. That is the right order and it caught two real defects.

### F14 — health used half of the shared ownership predicate

**MAJOR.** Round 4 required one ownership answer across health, install and uninstall.
`installerOwnsDestination` compares package name **and** skill identity; `skillUnitStatus`
compared only the name. So a marker naming a *different* skill reported `valid`,
`managed: true` — while both mutation commands refused the same directory. codex measured it
by changing one field: an upgrade or uninstall refused a tree the health gate had just
approved.

**Fixed.** Health now reports a distinct identity problem and `malformed`. The three commands
give the same answer, which is what round 4 asked for and what I only half delivered.

### F15 — a read command's filter was treated as the recorded selection

**MAJOR, and self-inflicted by cycle 4.** `expectedAddonNames` lets an explicit `--only` /
`--no-addons` override the core marker. Cycle 4's residual detection then used that
command-local set as its baseline, so on a **healthy full install** `doctor --only
parley-bidding` labelled the four other *recorded* add-ons "not part of the recorded
selection", failed health, and advised deleting them.

A narrowing flag must narrow. **Fixed**: residual detection now runs only for unflagged read
commands. With an explicit selector the command inspects exactly what was requested.

### F16/F17 — the two MINORs

- **The probe cache ignored the working directory.** `spawnSync` got no `cwd`, so a relative
  `PATH` entry resolved against the process directory; two calls sharing an environment but
  not a directory reused each other's verdict. The probe now takes the effective `cwd`, passes
  it to `spawnSync`, and includes its resolved value in the key.
- **My validation record overstated its own evidence — again.** It said `npm test` produced
  305 node **plus 54 Python** passes on 3.9.6 and 3.14. It does not: on 3.9.6 the Python leg
  **refuses to run** (`python3 is 3.9, but the add-on declares >=3.10`) and zero Python tests
  execute through it. Run file-by-file, 3.9.6 does pass all 54 — a compatibility fact, not the
  package gate. The evidence table now separates the two legs and names the interpreters each
  was measured on. This is the second time a reviewer has caught me stating a result more
  broadly than I measured it; both times the code was right and the claim was not.

### Measured after the cycle

- node leg: **309 / 309**, on python3 3.9.6 and 3.14
- Python leg: **54 / 54** on 3.10, 3.11 and 3.14; refuses 3.9.6 by design
- three new tests: ownership agreement across doctor/install/uninstall, a read filter that
  narrows instead of accusing, and unrelated sibling directories being ignored

## Notes for reviewers

Three places worth adversarial attention:

1. **`test/design-addons.test.js` was refactored, not extended.** The extractor that took 21
   review rounds and 29 fix-up cycles is now parameterized. The node arm's 253 assertions all
   pass, but a reviewer should check whether parameterization weakened any of them rather than
   trusting the count.
2. **The self-consistent-swap boundary.** `verifyPayload` alone *passes* that case by
   construction; only the marker catches it. If the marker is writable, so is the check.
3. **`manifestProblems` treats a marker with no `markerSchema` as legacy.** That is the one
   path where the manifest check is skipped entirely. It exists so 2.0.0 installs do not turn
   malformed on upgrade — but it is also the mechanism's softest edge.
