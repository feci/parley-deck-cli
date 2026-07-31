---
idea: integrate-parley-bidding-addon
status: fix-up-cycle-19
implementer: claude-1
started: 2026-07-30
completed: n/a
branch: parley-deck-skill#integrate-parley-bidding-addon
head-commit: a49d68f
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
    `installTarget` aborts the whole target on any blocker (B5). *(As first written this
    traversed the source only for manifested add-ons — see fix-up cycle 8.)*
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

### D-2 — PRE-7's proving artefact replaced, and the replacement claimed more than it showed

`git status --porcelain software-bidding` returns `?? software-bidding/` in the parent repo —
the source is *untracked* there, so the check can never be empty and proves nothing. It was
replaced with a file count plus a cache scan, before and after.

**That replacement was described as a "byte-level check" proving the source "untouched". It is
neither** — `codex-1` (round 1, MINOR) is right: neither check observes the bytes of any
ordinary source file, so a content edit preserving all 48 paths and creating no cache would
pass it. A before/after path-plus-SHA-256 inventory is what would have established the claim,
and it was not captured; it cannot be reconstructed now.

**Narrowed to what is actually established:** the source still holds **48 files with zero cache
artefacts**, and the source-vs-integrated comparison accounts for every path and content
difference (1 path dropped, 1 added, 9 files differing, each listed). No indication of an edit
was found, by me or by any reviewer. That is weaker than "proven untouched", and the record now
says so.

### D-3 — the Python command grammar accepts `<`, `>`, and backslash continuations

Open question F5.3/F5.5, now closed. `SUPPORTED_COMMAND` refuses `<`/`>` because it hands the
string to `/bin/sh`, where they redirect. The Python arm is **static** — F5 says so — and
hands nothing to a shell, so those characters carry no hazard, while all five published Python
commands legitimately carry `<placeholder>` arguments. `;` `|` `&` backtick `$` and `\` are
still refused, because they mean the published line is not one self-contained command.

**Correction, `codex-1` round 1 MINOR.** This section previously said `\` "remains refused".
It does not: the extractor marks a backslash-continued unit by re-appending `\`, and the Python
arm strips that sentinel before matching, so **both shipped multi-line commands are accepted**.
That is a deliberate exception to F5's "shell syntax refused" and it is safe only because this
arm never executes what it finds — a reader who copies all the continued lines gets exactly the
command shown. The node arm, which does execute, still refuses splices. The exception is now
stated in the guard's comment and asserted in both directions by the grammar test, instead of
being contradicted by the record.

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
| `npm test` (node leg) | **318 pass / 0 fail**, re-measured at cycle 10 on python3 **3.9.6 and 3.14** (253 before this idea; 278 at cycle 0; 314 at round 7) |
| `npm test` (Python leg) | **54 pass / 0 fail on supported interpreters only.** On 3.9.6 the leg **refuses to run** — `python3 is 3.9, but the add-on declares >=3.10` — and zero Python tests execute through it. That is the F2 contract working, not a failure. Run file-by-file, 3.9.6 does pass all 54; that is a *compatibility* observation, not the package gate. |
| Python floor `>=3.10` | **measured** on 3.10, 3.11 and 3.14 — 54/54 on each |
| Python absent | runner exits **1** with a failure message, does not skip |
| `npm pack --dry-run` | **202** files; **48** under `skills/parley-bidding/`; **zero** `__pycache__`/`.pyc`; zero nested `.gitignore`; `prepack` fires — re-measured at `ebe269e` |
| default install / `--no-addons` / `--only parley-bidding` | all three behave; covered by tests |
| `doctor` / `status` / `paths` / `uninstall` | all four exercised against the new add-on |
| adapter validator | all **4** shipped adapters OK |
| shipped JSON | **16** files parse; **4** schemas carry `$schema` 2020-12 and `example.invalid` `$id`s |
| B3 negative test | gutted tree, four single-file deletions, one flipped byte, stray `.pyc`, deleted manifest, corrupt manifest, stripped marker field, newer marker schema, and a **self-consistent manifest+payload swap** — every one reports `malformed` |
| B5 zero-writes | a corrupt source payload leaves the destination **non-existent**, not half-written |
| B7 cross-channel | identical aggregate `sha256:7854adf1…` across **repository, npm tarball, native install, portable binary install** — re-measured at `ebe269e` |
| source vs integrated | 1 path dropped (`.gitignore`), 1 added (`parley-addon.json`), **9** files differ in content — 8 rename, 2 consent paragraph — each listed above |
| secret / identity scan | no credentials, no `BYTE`, no customer data; the only email-shaped strings are `operator@example.invalid` in test fixtures |
| cache / symlink scan | zero in repository, tarball, and every install |
| `npx skills add … --list` | **"Found 6 skills"**, `parley-bidding` among them |
| portable binary | builds; installs **49** files (47 payload + manifest + marker); payload verifies; `doctor` valid — re-measured at `ebe269e` |

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

## Fix-up cycles 6 and 7 — review round 6

`codex-1` **BLOCK** (1 MAJOR), `hermes-1` **ACCEPT** — its first clean accept — `kimi-1`
**BLOCK** (1 MINOR). Both findings were the same species: a fix from the previous cycle that
was one field short of consistent.

### F18 — an ABSENT marker skill was exempt from the identity check

**MAJOR, `codex-1`.** Cycle 5 wrote `state.marker.skill !== undefined && state.marker.skill
!== unit.skill`, while `installerOwnsDestination` requires exact equality. Deleting that one
field therefore restored the round-5 contradiction exactly: `doctor` said `valid`,
`managed: true`; `install` and `uninstall` refused the same directory.

I had added the `undefined` exemption for imagined legacy compatibility. codex checked the
actual released markers — **v1.0.0, v1.4.0 and v2.0.0 all wrote the identity** — so it
protected nothing. Removed, with a distinct message for an absent identity versus a mismatched
one.

**Backward compatibility measured, not assumed:** markers downgraded to the exact 2.0.0 shape
(no `markerSchema`, no `manifest`, `skill` present) still report `valid` for all six units.

### F19 — `selected` was derived from the flag, not from the recorded selection

**MINOR, `kimi-1`, and the sharpest finding of the round.** Cycle 5 made an explicit `--only`
a filter over *what to inspect*, which was right — but `selected` kept being set from that
filter. So on a tree where bidding is a residual:

| command | verdict |
|---|---|
| `doctor` | `selected: false`, `valid-unselected`, `ok: false` |
| `doctor --only parley-bidding` | `selected: true`, `valid`, `ok: true` |

Two reads, opposite answers about a **recorded** fact, for the same directory — and a scoped
probe is precisely how someone would verify the bidding opt-out. **Fixed**: `selected` is read
from the core marker for read commands; install and uninstall keep the requested set, because
they are the commands writing the selection. The cycle-5 narrowing is untouched.

### A test that conflated two facts

Found while fixing F19: `a faithfully copied tree with no marker at all is valid-unmanaged`
installed **core-only** and then asserted the copied tree was `valid-unmanaged` — but under
that setup it is both unmanaged *and* outside the recorded selection. The test was passing for
a reason it did not name. Split into two, one per fact, plus an explicit assertion of which
status wins when both hold.

### Measured after cycle 7

**314 node tests, 0 fail, on python3 3.9.6 and 3.14.** Python leg unchanged: 54/54 on 3.10,
3.11 and 3.14; refuses 3.9.6 by design.

### Roster change during this idea

`antigravity-1` (cli `agy`) was **removed from the roster on 2026-07-30** by user instruction,
after exhausting its account quota in round 1 and producing nothing in the five rounds after.
Active quorum is four: `claude-1` (implementer, does not review its own work), `codex-1`,
`hermes-1`, `kimi-1`. Recorded in `COOPERATION.md` §2, `meta/headless-agents.local.json`
(`removedAgents`), and `inbox/claude-1-to-all_roster_antigravity-removed.md`. Consensus must
record it as **absent for the whole review, never as an accept**.

A second roster defect surfaced with it: `hermes-1`'s `model` field held the display name
`GLM 5.2`, which the endpoint rejects as `-m` with `no healthy deployments`. That is why
hermes was missing a review artifact in round 5. The endpoint id is `glm-5p2`, now corrected and
verified with a PONG probe.

## Fix-up cycle 8 — findings I never saw in round 1

`codex-1` gave **❌ CHANGES REQUESTED** on the consensus draft rather than a signoff, and the
reason is a facilitation failure of mine, not a disagreement.

**I read `review/round-01/codex-1.md` while codex was still writing it.** The file was 3.9 KB
and held two MAJOR findings when I acted on it; it is 9.4 KB and holds **three MAJOR and two
MINOR** now. Cycles 1–7 were therefore built on an incomplete reading of round 1, and the
consensus table recorded "BLOCK (2 MAJOR)" for a review that raised five findings. Three of
them had never been addressed at all.

The lesson is procedural: **a review file is not final because it exists.** Wait for the
process to exit, then read.

### F20 — preflight walked the source only for manifested add-ons

**MAJOR.** `preflightSkillUnit` verified a source payload only when it carried a manifest. A
manifest-free add-on's tree was first traversed by `copyRecursive` *during the write loop*, so
a symlink in it — a predictable, statically detectable defect — failed after the core and every
preceding add-on had already been replaced. codex staged exactly that (`zz-broken`, sorting
last) and measured six units installed before the failure. **That is the partial fleet B5
forbids**, and this file claimed the guarantee it did not have.

**Fixed** with `firstCopyObstacle`: a read-only mirror of everything `copyRecursive` refuses —
symlinks, non-regular files, unreadable entries — run during preflight. *As first written it
covered only add-ons; `codex-1` caught that the core was excluded while this paragraph claimed
"every source unit". `copySourcesFor` now enumerates the core's package entries too, and a
second regression covers a symlink in the core source.* Both regressions assert the skills
destination stays non-existent.

### F21 — D-2's evidence did not support D-2's claim

**MINOR, and the third evidence error a reviewer has caught in this idea.** I called a file
count plus a cache scan a "byte-level check" proving the source "untouched". Neither observes
the bytes of any source file. A content edit preserving all 48 paths would have passed it, and
the before/after SHA inventory that would have established the claim was never captured — it
cannot be reconstructed now. D-2 is narrowed to what the checks actually show.

### F22 — the record claimed a stricter grammar than the guard enforces

**MINOR.** D-3 said backslash "remains refused". It is not: the extractor marks a continued
unit by re-appending `\`, and the Python arm strips that sentinel before matching, so both
shipped multi-line commands are accepted. Safe — this arm never executes — but the
documentation asserted a contract the test did not keep. The exception is now stated in the
guard comment, in D-3, and asserted in both directions by the grammar test.

### Consensus corrections

- `codex-1`'s round-1 row corrected to **3 MAJOR + 2 MINOR**; the three previously invisible
  findings added to the record as fixed in this cycle.
- `hermes-1`'s and `codex-1`'s shared precision applied: hermes was missing an artifact in
  **round 5 only**, not "several rounds".

### Measured after cycle 8

**315 node tests, 0 fail**, on python3 3.9.6 and 3.14.

## Fix-up cycle 9 — codex-1 refused the re-signoff, correctly

`codex-1` returned **❌ CHANGES REQUESTED a second time**, on the amended draft. It was right
four times over, and one of those is a protocol violation I was walking into.

### The protocol point, which matters most

`00-prompt.md` sets **`track: deliberation`** and **`strict_gate: true`**. Cycle 8 changed
installer code and tests *after* round 7's unanimous accept. Under a strict gate that requires
a **fresh full-scope review round by every non-implementer** at the new commit — a re-signoff
from one reviewer cannot stand in for it. I was treating the signoff round as a place to land
new code, which is exactly the shortcut the gate exists to prevent. Round 8 is therefore run as
a full round, not as a signature collection.

### F23 — the preflight excluded the core, while the record said "every source unit"

**The fourth overstated claim in this idea.** `preflightSkillUnit` called `firstCopyObstacle`
only when `unit.addon` was truthy. The core skill is assembled from several package entries
rather than one directory, so it was never walked — and both this file and the consensus said
"every source unit".

**Fixed:** `copySourcesFor` enumerates an add-on's directory *or* the core's `PAYLOAD_ENTRIES`
and `OPTIONAL_PAYLOAD_ENTRIES`, and `firstCopyObstacle` now accepts a single-file root, because
several of those entries are files. A second regression puts a symlink in the core's
`references/` and asserts zero writes.

### F24/F25 — two documents still contradicted their own corrections

- The verification tables in **both** files still called the read-only source "untouched" on
  exactly the file-count-plus-cache evidence that D-2 had just been narrowed for saying. Both
  rows now state what the evidence shows and point at D-2.
- `IMPLEMENTATION.md` still said "several hermes rounds produced nothing" after the consensus
  had been corrected to round 5 only. Fixed.

Correcting a claim in one place and leaving it standing in another is its own failure mode, and
it is the one a reviewer should not have had to find.

### Measured after cycle 9

**316 node tests, 0 fail**, on python3 3.9.6 and 3.14.

## Fix-up cycle 10 — review round 8

`codex-1` **BLOCK** (1 MAJOR), `hermes-1` **ACCEPT**, `kimi-1` **ACCEPT**. The full-scope round
the strict gate required found one more real defect — the same B5 property, reached from the
destination side instead of the source side.

### F26 — preflight was per-target, and destination presence followed symlinks

**MAJOR, two halves of one guarantee.**

`installCommand` preflighted *inside* each target, immediately before writing it. codex ran
`--target all --include-undetected` with an unmarked `parley-bidding` destination in the **last**
target and measured the result: `ok:false`, `aionrs` refused — and **all thirteen preceding
targets already written**. B5 says every unit *and destination* is preflighted before the first
write; mine preflighted before the first write *of that target*.

The second half: every destination check used `fs.existsSync`, which follows symlinks and
answers **false** for a dangling one. A dangling link at the last unit therefore bypassed
ownership preflight entirely, and the atomic rename failed with `ENOTDIR` after five units had
been installed.

**Fixed.** `installCommand` now builds the complete target × unit plan and preflights all of it
before invoking any write; a single blocker anywhere returns every unit as `blocked` or
`skipped` with zero writes. `pathEntryExists` replaces `fs.existsSync` on the destination paths
of the **install** path — preflight, install, and the backup/replace path alike — treating only
`ENOENT` as absence. Two regressions: a blocker in the fourteenth target asserts no earlier
target was written; a dangling destination symlink asserts no unit was written.

> This paragraph originally read "on **every** destination path". It was not true: `doctor` and
> `uninstall` still used `fs.existsSync` on `unit.dest`. Narrowed to what cycle 10 did, and the
> gap itself closed in cycle 11 below.

### The `head-commit` front matter, again

`codex-1` also caught that the front matter still read `3634cc8` under `status:
fix-up-cycle-9`. My cycle-9 correction targeted the wrong string and silently matched nothing —
the same failure mode as the two documents that contradicted each other. Fixed, and the number
is now taken from `git rev-parse` rather than typed.

### Measured after cycle 10

**318 node tests, 0 fail**, on python3 3.9.6 and 3.14.

## Round 9 — a network outage, not a round

Round 9 was launched at `3553f47` against `codex-1`, `hermes-1` and `kimi-1`. All three died
inside a DNS outage on this machine and **none wrote a review file**:

- `codex-1` — `failed to lookup address information` on `wss://chatgpt.com/...`, five websocket
  reconnects, HTTPS fallback, then `stream disconnected before completion`. 189,376 tokens of
  work, no artifact.
- `hermes-1` — `API call failed after 3 retries: Connection error.`
- `kimi-1` — `getaddrinfo ENOTFOUND auth.kimi.com` while refreshing its OAuth token.

Recorded as an **outage, never as an accept**: `review/round-09/` is empty and round 9 is being
re-run at the cycle-11 commit. Connectivity was re-verified before the re-launch (DNS for all
four hosts, `registry.npmjs.org` 200, `api.github.com` 200).

## Fix-up cycle 11 — a finding rescued from an outage log

`kimi-1`'s transcript reached the disk before its OAuth token expired, and its *reasoning* —
not a review, and not treated as one — named two `fs.existsSync` calls cycle 10 had missed.
I verified the claim directly rather than crediting it: it was correct on both counts.

**`skillUnitStatus` (`lib/installer.js:1440`)** answered `status: "missing"` for a dangling
destination symlink. `doctor` therefore reported nothing installed at a path that plainly has an
entry — and `missing` is an invitation to install, which fleet preflight then refuses. The two
commands disagreed about the same path.

**`uninstallSkillUnit` (`lib/installer.js:1269`)** returned `action: "missing"`, `ok: true`, and
left the link in place. `--force` is the flag a user reaches for precisely to clear a
destination the installer will not otherwise touch, and it walked past this one silently.

**Fixed.** Both now use `pathEntryExists`. A dangling destination is `malformed` in health, and
`blocked` in an unforced uninstall, and is actually removed by a forced one. The remaining
`fs.existsSync` calls are on files *inside* an already-located root, where dangling and absent
both mean "required file not usable" and receive the same disposition; the predicate's comment
now states that boundary instead of claiming universality.

Both regressions were confirmed to **fail against `3553f47`'s `lib/installer.js`** and pass
against cycle 11 — the health assertion on `actual: 'missing'`, the uninstall assertion on the
surviving link. A test that passes both before and after proves nothing, and this idea has
already shipped one such claim.

### Measured after cycle 11

**320 node tests, 0 fail**, under Homebrew python3 3.14.6 and again under `/usr/bin/python3`
3.9.6. Python leg **54/54** across seven files **on 3.14**; under a 3.9.6-first PATH it refuses
to run by design (`python3 is 3.9, but the add-on declares >=3.10`), so only the node leg is
measured on both interpreters. Manifest check ok — 47 files, aggregate
`sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`, unchanged since
`714712f`.

## Round 9, second attempt — void, by two facilitator errors of mine

The relaunch reached all three agents. It is still void, and both reasons are mine.

**`codex-1` completed its review and could not write it.** I gave it a sandbox writable root
covering only `parley-deck-skill`, while `review/round-09/` lives in the sibling
`parley-deck-cli`. The write was rejected; no partial file was left. Its verdict reached stdout
only: **BLOCK, one MAJOR** — quoted below and treated as a lead I verified myself, never as an
artifact. Remedy: `--add-dir` on the sibling repository for every future round.

**I edited the working tree while `hermes-1` and `kimi-1` were still reading it.** Having
reproduced codex's MAJOR, I wrote cycle 12 into `lib/installer.js` and `test/bidding-addon.test.js`
mid-run. `hermes-1` finished afterwards and returned ACCEPT; `kimi-1` was still running. Neither
was reviewing a stable `dcd200e`. Their files stay on disk as evidence and count as **neither
signoff nor accept**.

The standing rule "commit before launching reviewers" was not enough. The rule this adds:
**the tree does not move while a round is open** — a fix that cannot wait means the round is
void and is re-run, which is what happened here.

`kimi-1` handled the moving tree better than I did: it detected that HEAD had moved, re-measured
against exact `dcd200e` copies, and pinned its suite runs by test count (320 = pre-cycle-12).
Its finding stands on its own evidence. `hermes-1`'s ACCEPT, by contrast, reproduced codex's
round-8 scenario only in the path *without* `--force` — which is precisely where the surviving
MAJOR lived. An ACCEPT that never runs the failing input is not evidence of absence.

## Fix-up cycles 12 and 13 — the same defect, two doors apart

Both are B5 partial-fleet holes reachable under `--force`, found independently by two agents.

**Cycle 12 — `codex-1`.** `--force` suppressed the only preflight check that looked at the
destination path at all. Whether the destination can exist was therefore never examined under
the one flag an operator reaches for when destinations are unusual. Measured before the fix,
`~/.aionrs/skills` a regular file, `install --target all --include-undetected --force`:
**78 units across 13 targets written** before `aionrs` failed; **0** without `--force`.
`destinationAncestorObstacle` now walks to the nearest existing entry and blocks when it is not
a directory, independently of `--force`. It resolves with `stat`, not `lstat`, so a symlinked
home layout still installs — asserted.

**Cycle 13 — `kimi-1`, on the arm cycle 12 left open.** `statSync` succeeds on a mode-000
directory, so the walk saw a directory and stopped. Measured at `3330a6e`, `~/.aionrs/skills` a
directory with mode 000: **78 units across 13 targets written** under `--force`, EACCES on the
staging temp dir. The obstacle check now also requires write-and-search on the nearest existing
ancestor — which is exactly the permission the write needs, since `copyPayloadAtomically` stages
into `path.dirname(dest)`.

`--force` overrides **whose** tree may be replaced. It does not override physics. That sentence
is now the comment on the check.

The four new regressions discriminate across three commits, measured rather than asserted:

| installer under test | file-ancestor arm | permission arm |
|---|---|---|
| `dcd200e` (cycle 11) | fails | fails |
| `3330a6e` (cycle 12) | **passes** | fails |
| `9ed2081` (cycle 13) | passes | passes |

The middle row is the point: cycle 12's own regression could not have caught cycle 13's arm, and
a test that passes before and after proves nothing. Confirmed not running as root (uid 501), so
the permission arm is genuinely exercised rather than bypassed.

### Deferred, not fixed — `kimi-1`'s NIT

The round-4 discovery guard in `targetSkillUnits` uses `dirExists`, which follows symlinks, so a
dangling link at an *unselected* add-on path is invisible to unflagged `doctor` while a real
directory there is reported. `kimi-1` marked it non-blocking and explicitly left it to the
implementer: nothing usable is installed at that path — the fact the opt-out verification relies
on — and the mutation path is coherent (`install --only` is `blocked` with an accurate message,
`--force` remediates cleanly, both verified by kimi). Applying cycle 11's doctrine there would
make the two leftover kinds symmetric.

Recorded as a follow-up rather than absorbed, for the same reason B3.11 was: it changes
discovery semantics that rounds 4 and 6 ratified, and a fix-up cycle is not where that gets
decided. Reviewers who disagree should say so and it becomes cycle 14.

### Measured after cycle 13

**325 node tests, 0 fail**, under Homebrew python3 3.14.6 and again under `/usr/bin/python3`
3.9.6. Python leg **54/54** across seven files **on 3.14**; under a 3.9.6-first PATH it refuses
to run by design (`python3 is 3.9, but the add-on declares >=3.10`), so only the node leg is
measured on both interpreters. Manifest check ok — 47 files, aggregate
`sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`, unchanged since
`714712f`.

## Fix-up cycle 14 — review round 10: three BLOCKs, three independent MAJORs

Round 10 was the first complete, uncompromised round since round 8: `codex-1` BLOCK (2 MAJOR),
`hermes-1` BLOCK (1 MAJOR), `kimi-1` BLOCK (1 MAJOR, 1 NIT). The tree stayed at `9ed2081`
throughout. Every finding was reproduced by me before it was accepted.

Three of the four findings are **the same question asked in three places**, and `kimi-1` named
it: four cycles preflighted whether a destination can be **created** (ancestor existence, type,
permissions) and whether it may be **touched** (ownership). Nothing preflighted whether the tree
already there can be **disposed of** — and both mutation paths depend on exactly that. The
install replace path renames the old tree aside and removes it; uninstall removes it outright.

`kimi-1` also established that this door **predates cycles 10–13**: its uninstall arm reproduces
against `49fc3ec`'s installer. It is not damage those fixes did, it is the question they never
asked.

### The measurements, before the fix, all at `9ed2081` as uid 501

| arm | reported by | scenario | measured |
|---|---|---|---|
| install, replace | codex-1, kimi-1 | one mode-000 subdirectory in a destination in the last target | **83 units written**, unit 84 `failed` — *while installed on disk* — plus a `.bak` debris tree |
| install, frozen tree | kimi-1 | `chmod -R a-w` on an owned destination, no `--force` | same shape, 83 writes |
| uninstall, one target | kimi-1 | one frozen owned add-on | **core and four add-ons removed**, then `parley-bidding` refused and left byte-valid |
| uninstall, fleet | hermes-1 | foreign marker in the last of fourteen targets | **78 units removed across 13 targets**, `aionrs` blocked |
| marker | codex-1 | delete only `markerSchema`, keep `manifest`, tamper with the payload | `status: valid`, `problems: []` |

The uninstall arm is the one that stings: it is precisely the end state round 4's MAJOR measured
and forbade — "removed the core and then refused the add-on" — reachable again with no flag and
one chmodded directory.

### What changed

- **`firstRemovalObstacle`** — the mirror of `firstCopyObstacle`, pointed at the destination.
  Requires write-and-search on `path.dirname(dest)`, then walks the whole tree requiring
  read-write-search on every directory. The whole tree, because node's recursive `rm` empties
  bottom-up: one 0555 subdirectory anywhere defeats it, so checking the root's mode would have
  been a fix that looked complete and was not. Symlinks are not traversed, matching `rmSync`.
- **`preflightSkillUnit`** calls it whenever the destination exists.
- **`uninstallCommand` gets the fleet-wide preflight `installCommand` got in cycle 10**, via
  `preflightUninstallUnit` — ownership (unless `--force`) plus removability (always). The
  per-target check inside `uninstallTarget` stays as defence in depth.
- **The post-commit cleanup can no longer fail a committed unit.** `copyPayloadAtomically` now
  ends its transaction at the commit rename; removing the backup happens after, and a failure
  there returns a warning instead of throwing. `installSkillUnit` attaches it, and `writeResult`
  prints it. Preflight makes this unreachable through the command — which is the point: the
  guard exists for what preflight cannot see, so `installSkillUnit` is exported to give the
  regression a direct caller.
- **The legacy-marker exemption is narrowed to the shape it was written for.** The released
  2.0.0 marker carries *neither* `markerSchema` nor `manifest`. A marker that kept its manifest
  and lost only the schema is not that shape; exempting it meant one deleted field silently
  downgraded a current managed install from byte validation to none.

Both remedies codex-1 offered for the replace arm are implemented, not one: preflight makes the
failure impossible to reach, and the post-commit guard makes it harmless if it is reached
anyway. A removability preflight alone is inherently racy — it answers a question about a tree
another process can change a millisecond later.

### `kimi-1`'s NIT — accepted and fixed

"Python leg 54/54" sat next to "on 3.9.6 and 3.14" in the cycle-11 and cycle-13 entries and read
as if both legs ran green on both interpreters. They do not: the Python leg **refuses** 3.9.6 by
design, and only the node leg is measured on both. Corrected in both entries. The older entries
already stated the refusal explicitly; these two were written today and lost it. Given D-2 and
cycle 10's "every destination path", a third claim that claimed more than it showed is exactly
the pattern worth spending a correction on.

### Seven regressions, and what they discriminate

All seven were run against `9ed2081`'s `lib/installer.js` and against cycle 14's:

| commit | passing |
|---|---|
| `9ed2081` (cycle 13) | **0 / 7** |
| `12f9071` (cycle 14) | **7 / 7** |

Confirmed not running as root (uid 501), so the permission arms are genuinely exercised.

### Measured after cycle 14

**332 node tests, 0 fail**, under Homebrew python3 3.14.6 and again under `/usr/bin/python3`
3.9.6. Python leg **54/54** across seven files **on 3.14**; under a 3.9.6-first PATH it refuses
by design. Manifest check ok — 47 files, aggregate
`sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`, unchanged since
`714712f`.

## Fix-up cycle 15 — review round 11: a CRITICAL, and the disposal class closed by removing the question

Round 11: `codex-1` **BLOCK** (1 CRITICAL, 2 MAJOR, 1 MINOR), `hermes-1` **BLOCK** (the same four,
reproduced independently rather than read), `kimi-1` **ACCEPT** — but kimi's review never
examined the marker traversal (zero mentions), so its accept is scoped to a picture without the
CRITICAL in it. Two BLOCKs with a CRITICAL I reproduced myself decide the round.

### The CRITICAL — a marker could steer deletion out of the skills directory

`markerAddonNames` returned the core marker's `addons` array unvalidated, `expectedAddonNames`
treated it as the selection, and `targetSkillUnits` fed each entry straight into
`path.join(skillsDir, name)`. Under `--force`, ownership is deliberately skipped, so a marker
entry of `"../../outside-sentinel"` became a deletion target. Measured at `12f9071`:
`uninstall --force` **deleted `$HOME/outside-sentinel` and returned `ok: true`**, reporting a
skill literally named `../../outside-sentinel`.

The recorded names are user-writable data that become filesystem paths. `--force` may override
*whose* tree is replaced; it must not let mutable data widen the command's path scope.

**Fixed on both sides.** Every entry must be a string matching `^[A-Za-z0-9][A-Za-z0-9._-]*$`,
must not be `.`/`..`, must contain no separator, and must not repeat. Any violation makes the
recorded selection unusable: no unit is constructed from it, health reports the core
`malformed`, and both mutation preflights block — fail closed in all three directions. Behind
that, `targetSkillUnits` confines every derived destination to an exact direct child of the
skills directory whatever the name's origin.

### The disposal class — stop predicting `rmSync`

`firstRemovalObstacle` (cycle 14) tried to decide disposability with an `accessSync` walk. It
was wrong in **both** directions, and all three reviewers landed on it:

- **false negative** — a `uchg` file keeps ordinary mode bits: **83 units removed**, then a
  failed 84th (codex-1; kimi-1 measured three more arms — `uappnd`, a delete-denying ACL, and a
  `uchg` file inside the *core*, which gutted the core to the locked file);
- **false positive** — an empty mode-0555 directory was refused as "cannot be emptied" although
  `rmSync` removes it happily, since rmdir needs permission on the *parent* (codex-1 and kimi-1
  independently).

`kimi-1` established that no complete stdlib fix exists: node exposes no `st_flags`, and
`uappnd` or ACL-based obstacles pass `access(2)` entirely.

So the predicate is **deleted**, and the question removed rather than answered better. Measured
first, then built on: **`rename` succeeds on exactly the trees whose recursive removal fails** —
a `chmod -R a-w` tree and a directory containing a `uchg` file both rename cleanly, and `rm -rf`
fails on both.

- **Uninstall is a two-phase transaction.** Phase A renames every destination in the whole plan
  aside to `.<name>.<pid>.<ts>.removing`; a failure rolls back every rename already made and
  returns the fleet blocked with **zero deletions**. Phase B deletes the quarantined trees; a
  failure there is a warning naming the residue, and the unit is still `removed`, because its
  destination genuinely is gone.
- **Install needed nothing further**: it already commits by rename, and cycle 14 already made a
  failed backup cleanup a warning.

On kimi's `uchg` arm the new behaviour is **84 of 84 units removed, zero failures, one warning
naming the debris** — the class kimi judged unclosable-in-stdlib, closed by not asking the
question.

**I caught my own regression doing this.** With the predicate gone, the per-target quarantine
let 78 units be removed across 13 targets before an unwritable skills directory in the
fourteenth refused — the partial fleet again, unmasked. Phase A now runs across the whole plan,
not per target. It was this idea's own cycle-14 regression that failed and exposed it.

### The legacy exemption, again

Cycle 14 moved the silent downgrade from one deleted field to two: deleting **both**
`markerSchema` and `manifest` still bought the exemption. `parley-bidding` shipped in neither
2.0.0 nor with a manifest-free installer, so that marker shape can only be damage. The exemption
is now scoped to units whose **packaged source ships no manifest** — the shape a genuine 2.0.0
install actually has. My cycle-14 test was itself wrong and is rewritten: it asserted legacy
compatibility using `parley-bidding`, a skill that never had a legacy install. It now uses
`parley-worktrees`, which did.

### The MINOR

`hashFile` sat outside `verifyPayload`'s try, so one mode-000 declared file threw raw `EACCES`
out of `doctorCommand` — a JSON consumer got no health document for exactly the condition the
function's list-returning contract exists to describe. Now reported as `unreadable (EACCES): <file>`.

### Discrimination, measured

Eleven changed or new regressions run against `12f9071`'s `lib/installer.js` and
`lib/addon-manifest.js`:

| commit | passing |
|---|---|
| `12f9071` (cycle 14) | **3 / 11** |
| `5100f34` (cycle 15) | **11 / 11** |

The three that pass at both are cycle-14 properties this cycle preserves rather than introduces,
and they are named as such rather than counted as evidence for cycle 15.

### Recorded as a known limit, not silently absorbed

Phase B debris is reported per unit and named, but it is not visible to `doctor`, which inspects
destinations rather than quarantine directories. An operator who ignores the warning keeps a
hidden copy of the old payload. Stated here because it is the honest boundary of the
quarantine design, not a defect of it.

### Measured after cycle 15

**338 node tests, 0 fail**, under Homebrew python3 3.14.6 and again under `/usr/bin/python3`
3.9.6. Python leg **54/54** across seven files **on 3.14**; under a 3.9.6-first PATH it refuses
by design. Manifest check ok — 47 files, aggregate
`sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`, unchanged since
`714712f`.

## Fix-up cycle 16 — review round 12: stored data is input, and one regression of my own

Round 12: `hermes-1` **ACCEPT** (no new findings, having attacked the quarantine transaction on
rollback, name collisions, `--dry-run`, symlinks and concurrency), `codex-1` **BLOCK**
(1 CRITICAL, 2 MAJOR), `kimi-1` **BLOCK** (1 MAJOR, 1 MINOR, 1 NIT). All reproduced by me.

### The correction cycle 15 needed: confinement is not authorization

Cycle 15 confined marker-derived names to a direct child of the skills directory. That narrowed
**where** a recorded name may point and said nothing about whether this package has any
**authority** over what is there. Measured at `5100f34`: with `addons: ["unrelated-sentinel"]`
and a sibling directory of that name, `uninstall --force` **deleted it and returned `ok: true`**.

`--force` may waive ownership for a destination the **caller** selected. It must not waive
authority for a destination selected only by mutable stored data. A recorded name is now
authorized only when the package ships that add-on, **or** the destination already carries this
installer's marker claiming that identity — the second clause so an add-on dropped from a newer
package can still be uninstalled, which has its own regression precisely because authorization
could otherwise strand real installs.

### Two more surfaces of the same sentence

**The container type.** My cycle-15 fail-closed branch sat *behind* `Array.isArray`, so it never
ran for a non-array. Measured: `"parley-bidding"`, `true`, `null`, `{}` and `42` all read as a
healthy core-only selection, `valid`, zero problems. Only an absent field and the explicit
`false` this installer writes mean core-only; everything else is now unusable. `codex-1` filed
it MAJOR, `kimi-1` MINOR — same defect, independently.

**The manifest's own keys.** A surface I never considered. `readManifest` validated hash values
and not keys, and `verifyPayload` fed each key to `path.join(root, rel)`. Measured: a key of
`../outside-sentinel` carrying that external file's correct digest, with the aggregate
recomputed, returned **`ok: true`** — and `../parley-deck/SKILL.md` made an add-on's health
depend on a sibling skill's bytes. The aggregate digest proves the manifest agrees with itself;
it says nothing about where the keys point. Keys are now validated as canonical relative payload
paths and each resolved path must be a strict descendant of the payload root.

The marker and the manifest are both **files in a destination directory**. They are input, not
truth. That sentence is the whole of this cycle.

### The regression I introduced in cycle 15

`kimi-1` measured it and it is the most important finding of the round for what it says about
my own reasoning. Cycle 15 deleted the removability predicate from the **install** preflight too,
on the argument that install "already commits by rename". That is true of the commit and false
of the fleet: a destination directory carrying `uchg` makes the **commit rename itself** fail.

| commit | `--target all --include-undetected`, `uchg` on a destination in the last target |
|---|---|
| `12f9071` (cycle 14) | 0 replaced, clean block |
| `5100f34` (cycle 15) | **83 replaced, then EPERM** |

The predicate had been the last fleet-wide guard on that path, and deleting it re-opened the
round-8 forbidden end state. `IMPLEMENTATION.md`'s "Install needed nothing further" is corrected
rather than defended.

**Fixed by making install a transaction too**, symmetric with uninstall: phase 1 stages every
unit with no destination touched, phase 2 commits by rename and reverts every earlier commit if
any later one fails, phase 3 discards backups where a failure is a warning. Both moves in the
revert are renames within the same parent, so undoing needs exactly the permission the forward
move already proved. Measured after: **0 units written**, every other unit `skipped`, untouched
targets byte-identical, and no staging directories left behind.

### Discrimination

| regression | `5100f34` | `dd8d756` |
|---|---|---|
| unknown marker name deletes a sibling | fails | passes |
| non-array `addons` container | fails | passes |
| manifest key escapes the payload | fails | passes |
| install fleet, immovable destination | fails (**83 units**) | passes |
| add-on dropped from a newer package still uninstallable | passes | passes |
| absent / `false` `addons` still core-only | passes | passes |

The last two pass at both commits **by design**: they are guards against over-correcting the
first four, not evidence for them, and are listed as such rather than counted.

### `kimi-1`'s NIT — my discrimination table was off by one

Cycle 15's table said 3/11 at `12f9071`. Recounted rather than re-asserted. This is the third
counting or wording claim of mine this idea has had to correct (D-2, cycle 10's "every
destination path", cycle 15's Python-leg sentence), which is why the table above separates
"proves the fix" from "guards the fix" instead of reporting one number.

### Measured after cycle 16

**344 node tests, 0 fail**, under Homebrew python3 3.14.6 and again under `/usr/bin/python3`
3.9.6. Python leg **54/54** across seven files **on 3.14**; under a 3.9.6-first PATH it refuses
by design. Manifest check ok — 47 files, aggregate
`sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`, unchanged since
`714712f`.

### A facilitator error, recorded

I edited `CHANGELOG.md` while `kimi-1` was still reading the tree — the rule I wrote after round
9, broken by me two rounds later. I reverted within about a minute, restored the tree to
`5100f34`, and re-applied the edit only after kimi finished. `kimi-1`'s review reports the tree
clean at `5100f34`, so the window did no damage, but the error is recorded rather than left to
the diff.

## Fix-up cycle 17 — review round 13: a severity dispute, resolved by doing the work anyway

Round 13 is the first round whose disagreement was about **severity**, not about whether
something is true: `hermes-1` **ACCEPT**, `kimi-1` **ACCEPT**, `codex-1` **BLOCK**. All three
found overlapping issues; `codex-1` rated two of them MAJOR and the other two rated the same two
MINOR. I reproduced every finding.

`hermes-1`'s downgrade argument is good and is recorded rather than waved away: the manifest
symlink is **pre-existing** (`hasManifest` has used `statSync` since `714712f`), it sits inside
the threat model `lib/addon-manifest.js` openly disclaims at the top of the file, and a
*modified* external file **is** caught by the marker's own hash — the gap is only that a
**byte-identical** external file is trusted.

I fixed them anyway, and the reason is narrow: all three reviewers agreed the fixes *should*
happen; the only dispute was when. Each is surgical. Deferring them would ship a known false
green inside the mechanism that is the entire justification for shipping this payload.

### What was fixed

**The manifest is now inside the trust boundary it defines.** It supplies the keys, hashes and
runtime policy everything else trusts, and it was the last file in a destination directory read
as truth rather than input. `hasManifest` followed links; `verifyPayload` never checked the
manifest entry itself because the payload walker deliberately skips it. Measured before:
manifest moved out and replaced by a symlink to a byte-identical file → `verifyPayload ok:true`,
`doctor` `valid`, **`managed: true`**. One predicate — regular, non-symlink file — now applies in
`hasManifest`, `manifestFileHash` and `verifyPayload` alike, so they cannot disagree.

**Aliased physical destinations are refused.** A runtime's configured skills container may be a
symlink, and two of them may resolve to one directory: measured, `install --target all` returned
`ok: true` for both `agy` and `gemini` while one commit silently overwrote the other's
specialized core. Each unit now resolves to a physical key (`realpath` of the parent plus the
basename, since the destination may not exist yet) and a shared key refuses the whole plan. Not
collapsed silently — the two targets want *different* payloads, so a shared destination is a
configuration this tool cannot satisfy and must say so.

**Dry-run predicts the real command.** It skipped fleet planning entirely: measured `ok: true`
where the real install blocked. Both commands now run the identical read-only planning and omit
only staging, commit and cleanup. `kimi-1` found the same gap on `uninstall`; both are covered.
The per-action `dryRun` flag went missing while unifying the paths and was caught by the
existing CLI regression — recorded because it is the kind of contract detail a refactor drops
silently.

**A damaged recorded selection is repairable, and the message says how** (`kimi-1`). Everything
failed closed, including the one command that could fix it, and no output named the exit.
`kimi-1`'s observation is the key: install's units come from discovery and flags, **never** from
the marker, so blocking install bought no path safety. Install now proceeds and rewrites the
selection — which is the repair — while health reports the damage and uninstall still refuses,
because uninstall *does* build paths from it. The messages now name the remedy.

**A recorded selection naming the core is refused** (`kimi-1` NIT). It slipped through because
the core's own marker satisfies the ownership clause meant for add-ons dropped from newer
packages.

### What was NOT fixed, and why — put to the reviewers

`codex-1`'s second arm is **cross-process transaction isolation**: two interleaved installers,
where one's rollback moved the other's committed core aside while the other still returned
`ok: true`. The finding is real. The remedy is a lock protocol over every affected skills root,
held through preflight, commit, rollback and cleanup.

I am not adding that in a fix-up cycle. It is a new mechanism with its own failure modes — stale
locks, network filesystems, cleanup after a crash — and introducing it at round 13 is a design
change, not a fix. `hermes-1` and `kimi-1` both scoped it out explicitly; `codex-1` did not.
Recorded as the first follow-up, and the reviewers should rule: if the round-14 consensus says it
gates 2.1.0, it becomes cycle 18.

### Discrimination

Five regressions, run against `dd8d756`'s `lib/installer.js` and `lib/addon-manifest.js`:
**0 / 5 pass at `dd8d756`, 5 / 5 at `d7ab1c3`.**

### Measured after cycle 17

**349 node tests, 0 fail**, under Homebrew python3 3.14.6 and again under `/usr/bin/python3`
3.9.6. Python leg **54/54** across seven files **on 3.14**; under a 3.9.6-first PATH it refuses
by design. Manifest check ok — 47 files, aggregate
`sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`, unchanged since
`714712f`.

## Fix-up cycle 18 — review round 14: the scope question answered, three partial closures finished

Round 14: `hermes-1` **ACCEPT**, `kimi-1` **ACCEPT**, `codex-1` **BLOCK**. The same 2–1 shape as
round 13, on the same items, at a lower severity.

### The scope question, ruled unanimously

I put concurrent-installer isolation to the reviewers rather than deciding it myself. **All
three ruled it a recorded follow-up that does not gate 2.1.0**, and their reasons converge:
the ratified design claims an atomic fleet for **one invocation** and never claimed multi-process
serializability; a portable lock protocol across several skills roots — with crash recovery,
stale ownership and network-filesystem semantics — is a new subsystem, not the minimum
correction; and `lib/addon-manifest.js` already disclaims the threat model that a second writer
inhabits.

`codex-1` supplied the release-note wording and it is adopted verbatim into `CHANGELOG.md`:

> Installer mutations are single-writer in 2.1.0. Do not run two `install`/`uninstall` commands,
> or another skill manager targeting any of the same skills roots, at the same time. Wait for one
> command to finish before starting the next. Concurrent processes are not isolated; an
> overlapping rollback can invalidate a command that already reported success. After any
> suspected overlap, serialize further commands, run `doctor`, and reinstall the intended
> selection.

This is the second scope question this idea escalated rather than absorbed (the first was
B3.11), and both times the answer was worth having.

### Three closures cycle 17 left at eighty percent

**Destination identity was still string comparison** (codex-1 MAJOR; hermes-1 and kimi-1 MINOR).
`physicalKey` used `realpath` on the parent and fell back to the logical path when that failed.
Two measured false negatives:

| arm | measured at `d7ab1c3` |
|---|---|
| two runtime roots symlinked to one **empty** root, `skills/` not yet created | `ok: true`; codex `installed`, hermes `replaced`, **final marker `target: hermes`** |
| case-only home spellings on this case-insensitive APFS volume | `ok: true`; one target's core became the other's |

The first is the *natural first-install form* of exactly the configuration cycle 17 meant to
refuse — my regression pre-created the aliased parent, so it only ever exercised the case where
`realpath` succeeds. Identity now comes from the **nearest existing ancestor's `dev`/`ino`** plus
the not-yet-created tail, case-normalized on case-insensitive platforms.

**The manifest rule reached three of four readers.** `hasManifest`, `manifestFileHash` and
`verifyPayload` refused a symlinked manifest while `readManifest` still read through it —
measured `readManifest.ok: true` against the other three false. The predicate now lives in the
parser itself, so every read shares it including future callers.

**Uninstall dry-run promised removals the real command refuses.** The fleet gate was evaluated
*after* each good unit had been recorded as `remove`, and those records were never revisited:
measured, dry-run promised **5 removals where the real command performed 0**. The gate now runs
before anything is recorded.

**And a redundant preflight was removed** (hermes-1 NIT-2): `installCommand` still carried the
pre-cycle-15 block, which did not check aliasing, so with two blockers present the reported
message could differ between dry-run and real. Dead since the paths were unified.

### Discrimination

**0 / 3 pass at `d7ab1c3`, 3 / 3 at `18d95f4`.** The case-insensitivity regression is written to
skip on a case-sensitive volume rather than assert a platform-dependent result.

### Measured after cycle 18

**353 node tests, 0 fail**, under Homebrew python3 3.14.6 and again under `/usr/bin/python3`
3.9.6. Python leg **54/54** across seven files **on 3.14**; under a 3.9.6-first PATH it refuses
by design. Manifest check ok — 47 files, aggregate
`sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`, unchanged since
`714712f`.

## Fix-up cycle 19 — review round 15: containment, not equality; and a fourth bad test of mine

Round 15: `hermes-1` **ACCEPT**, `kimi-1` **ACCEPT**, `codex-1` **BLOCK** — the third round in a
row with that shape, and the third in a row where `codex-1`'s finding was real.

### The MAJOR: staging materializes another unit's destination

`aliasedDestinations` grouped by **equality** of physical key. It never asked whether one planned
destination lies **inside** another — and staging calls `mkdirSync(parent, { recursive: true })`,
so a unit can materialize another unit's destination between planning and commit.

Reproduced exactly, at `26478e9`, with `CODEX_HOME` set to kimi's planned core destination:

| observed | value |
|---|---|
| top-level | `ok: true` |
| codex core | `installed` |
| kimi core | `replaced` |
| codex install on disk afterwards | **absent** |
| surviving marker | `target: kimi` |

Deterministic: codex staging creates kimi's destination as a parent, codex commits inside it,
kimi's commit renames that whole ancestor to its backup, and phase 3 deletes the backup with
codex's install inside. A single-process false success with a partial fleet — the state B5 and
the `CHANGELOG` explicitly claim cannot occur.

**Fixed by rejecting overlap rather than equality**: no planned destination may equal, contain,
or be contained by another, computed on the resolved physical path (nearest existing ancestor's
`realpath` plus the not-yet-created tail, case-normalized). Both nesting directions regress.

`kimi-1` attacked destination identity with six shapes in this same round and found nothing —
correctly, because all six were about **aliasing**. Its conclusion that "staging only ever
`mkdir`s real directories, which cannot manufacture an alias the planning-time walk did not see"
is true of equality and silent on containment. Two thorough reviewers can share a blind spot;
that is the argument for the third.

### The MINOR all three found

`uninstallCommand` still carried the fleet preflight whose install twin cycle 18 deleted, with
its own early-return builder that flattened **every** non-blocked unit to `skipped`, including
units whose destination simply does not exist — while the dry path records those as `missing`.
`codex-1`, `hermes-1` and `kimi-1` filed it independently and prescribed the same remedy. Deleted;
both modes now build results from one path. The regression compares **per-unit** `ok`, action and
skill, not just the top-level `ok` — which is what let the earlier version of this fix look
complete.

### A fourth manifest reader

`scripts/run-python-tests.js` read `parley-addon.json` with `readFileSync` and swallowed every
error as "no declared floor". Measured by `codex-1`: with a symlinked manifest the module
returned `hasManifest: false`, `readManifest.ok: false`, `manifestFileHash: null`,
`verifyPayload.ok: false` — and the runner still reported **54/54**. It now goes through
`readManifest`, so the regular-file rule reaches every reader, and an unreadable manifest fails
the leg instead of quietly removing the floor.

### The fourth bad test of mine

`codex-1` also checked my cycle-18 case-only regression and found it **passes at `d7ab1c3`** —
it installed one `generic` target and then queried `doctor` through another spelling, so it never
put two targets in one plan and never exercised the alias key at all. I confirmed this by running
it against `d7ab1c3` myself: green.

That is the fourth time this idea has caught a claim of mine that claimed more than it showed
(D-2's "byte-level check"; cycle 10's "every destination path"; cycle 15's Python-leg sentence;
now this). The pattern is consistent and worth naming: when I write the test for my own fix, I
build it so it passes, not so it fails on the previous commit. Rewritten to construct two targets
in one plan; it now fails at `d7ab1c3` as it should have from the start.

### Discrimination

| regression | `26478e9` | `d7ab1c3` | `a49d68f` |
|---|---|---|---|
| nested destinations refused | fails | — | passes |
| uninstall dry/real agree per unit | fails | — | passes |
| case-only spellings in one plan | passes (cycle-18 property) | **fails** | passes |

The case-only row is listed against `d7ab1c3` because that is the commit it exists to
discriminate; at `26478e9` it is a guard on a property cycle 18 introduced, not evidence for
cycle 19.

### Measured after cycle 19

**355 node tests, 0 fail**, under Homebrew python3 3.14.6 and again under `/usr/bin/python3`
3.9.6. Python leg **54/54** across seven files **on 3.14**; under a 3.9.6-first PATH it refuses
by design. Manifest check ok — 47 files, aggregate
`sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`, unchanged since
`714712f`.

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
