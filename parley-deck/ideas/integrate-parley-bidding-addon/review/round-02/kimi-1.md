---
idea: integrate-parley-bidding-addon
review-round: 2
agent: kimi-1
date: 2026-07-30
reviewed-commit: 89069b0
---

## Verdict
ACCEPT

## Round-1 findings — closed or not

### 1. [CRITICAL] `doctor` reported an expected unit with a missing/unreadable marker as healthy — CLOSED
Re-measured my exact round-1 probe matrix against the real bin at `89069b0`, driving
`install --target codex --yes` then `doctor --target codex --json` in child processes with a
controlled `HOME`:

| probe | round 1 | round 2 (measured) |
|---|---|---|
| `missingMarker` (delete `.parley-deck-skill-install.json`) | `valid`, exit 0 | `malformed`, exit 1, problem: "no parley-deck-skill install marker: …" |
| `unreadableMarker` (`{ not json`) | `valid`, exit 0 | `malformed`, exit 1, problem: "the parley-deck-skill install marker is unreadable or is not valid JSON" |
| gut to `SKILL.md` + delete marker (double deletion) | `valid`, exit 0 | `malformed`, exit 1, same missing-marker problem |

Missing and unreadable are distinct, separately diagnosable messages — what I asked for in
open question 2. The check lives in `skillUnitStatus` (`lib/installer.js:1212`), so it applies
to every add-on and to the core unit, not only to `parley-bidding`: deleting the **core**
marker also fails health (measured: core `malformed`, exit 1). A foreign marker is a third
distinct problem ("the install marker belongs to \"some-other-installer\", not
parley-deck-skill").

### 2. [MAJOR] Python floor enforced in the test runner, not in `doctor` — CLOSED
Measured with a controlled `PATH` per child process:

| probe | round 1 | round 2 (measured) |
|---|---|---|
| `noPython` (empty `PATH`) | `valid`, exit 0 | payload stays `valid`, `runtime: {ok:false, requirement:">=3.10", detail:"python3 is not available…"}`, `doctor.ok:false`, exit 1 |
| python3 stubbed at 3.9 | n/a | `runtime.ok:false`, "python3 is 3.9, but this skill requires >=3.10", exit 1 |
| python3 stubbed at 3.10 | n/a | `runtime.ok:true`, exit 0 |
| real host python3 (3.14) | n/a | `runtime.ok:true`, exit 0 |

Payload validity and operational availability are reported as separate answers per B6:
`status` stays `valid` while health fails, and `doctor --json` carries both fields
unambiguously (`status`, `problems[]`, `runtime:{ok,requirement,detail}|null`, top-level
`ok`). Human output prints the per-unit `unavailable:` line plus a distinct stderr reason.
The probe is cached per process (`lib/installer.js:1290`) and reads the doctor process's own
environment — correct for a single-shot CLI run; my stubbed-PATH measurements prove it honors
the `PATH` doctor itself runs with.

### 3. [MINOR] IMPLEMENTATION.md front matter recorded the wrong head commit — CLOSED
Front matter now reads `head-commit: 89069b0`, matching the reviewed tree.

Open question 4 from round 1 is also answered: `CHANGELOG.md` / `RELEASING.md` belong to this
idea and are committed in `89069b0`; the worktree is clean (`git status --short` empty).

## New findings

### [NIT] Human-mode doctor stderr reads "One or more installs are installed but operationally unavailable."
**Where:** `lib/installer.js:1502-1509` (`writeResult`, doctor branch)
**What:** When the only failure is an unavailable runtime, the reasons list produces
"One or more installs are **installed but** operationally unavailable." — the subject
("installs") and the first word of the reason ("installed") collide.
**Why it matters:** Cosmetic only; the message is still unambiguous and the exit code is
correct. Not worth a cycle of its own.
**Evidence:** measured human-mode output of the `noPython` probe (exit 1).
**Fix:** Rephrase the reason to "operationally unavailable" so the sentence reads "One or
more installs are missing or malformed, or operationally unavailable." Can ride along with
any future touch of this file.

## What I verified and found correct

- **The reviewed commit is as stated:** HEAD is `89069b0`, one commit on `714712f`, worktree
  clean. All measurements below ran against the real bin in child processes with controlled
  `HOME`/`PATH`, not against a reading of the source.
- **No remaining path reports an unrunnable payload as healthy** among the ones I could
  construct: missing interpreter, below-floor interpreter, missing marker, unreadable marker,
  foreign marker, gutted tree, gutted core — every one fails health with exit 1.
- **The new marker rule breaks no legitimate case** (all measured):
  - a legacy 2.0.0-shape marker (no `markerSchema`, no `manifest` anchor) on an intact
    add-on → `valid`, exit 0 — upgrades in place do not turn malformed;
  - `install --no-addons` → doctor checks only the core unit, exit 0;
  - `install --only parley-tracker` → `parley-bidding` absent and *not* reported missing,
    exit 0;
  - the foreign-marker rule cannot false-positive on any historical marker: `name:
    "parley-deck-skill"` has been written since the installer's introducing commit
    (`344f797`), verified with `git log -S`.
- **Measured boundary (disclosed by the implementer, correct tradeoff):** a legacy-marker
  unit gutted to `SKILL.md` *with the marker kept* still reports `valid` (exit 0) — the
  `markerSchema`-less carve-out skips the manifest anchor check by design, and closing it
  would declare every 2.0.0 install malformed. Stated from measurement, not inference. Any
  reinstall rewrites the marker at schema 2 and closes the gap; the round-1 double-deletion
  path is closed regardless because deleting the marker is now itself fatal.
- **Fail-safe directions everywhere I pushed:** an unparseable `runtime.python` spec
  (anything not `>=X.Y`) reports `ok:false` ("unsupported python requirement"), never a
  silent pass; a `python3` that prints garbage is treated as not found.
- **Probe cache correctness:** module-level, one probe per process, consulted per unit; the
  child-process tests control `PATH` precisely because of this, and my independent
  empty-PATH / stub-PATH / real-PATH measurements agree with theirs.
- **No regression in `install` / `uninstall` / `status` / `paths`** (measured on fresh
  installs): all exit 0; uninstall removes both trees; reinstall + doctor green. The
  `readMarker` semantic narrowing (non-object JSON now reads as unreadable instead of being
  returned) is safe for its other callers (`uninstallSkillUnit`, `isInstallerOwnedSkill`,
  `markerAddonNames`) — they only ever dereferenced `.name`, which yields the same
  not-ours verdict as before.
- **Claimed test counts reproduce exactly:** `npm test` at `89069b0` → 286 node tests, 0
  fail; Python suite 54 tests OK across 7 files (python 3.14); the add-on manifest check
  passes (`47 files, sha256:7854adf1…`). The eight new tests cover the round-1 matrix
  including the double deletion and both interpreter directions. I did not run the Python
  suite outside the runner (it already invokes `python3 -B`); no payload file was touched
  and the worktree stayed clean throughout.
- **`doctor --json` shape is unambiguous:** `ok`, per-unit `status`, `problems`, and
  `runtime` (`null` for units with no declared runtime, for malformed units, and for
  missing units) give automation everything it needs without parsing prose.

## Open questions for the implementer

1. The interpreter probe is `python3`-only (`spawnSync("python3", …)`). On Windows hosts
   where only `python` exists — or where `python3` resolves to the Store app-execution alias
   — doctor will report the add-on unavailable. That is the fail-safe direction and matches
   `scripts/run-python-tests.js`, so I am not blocking on it, but is `python3`-only the
   intended contract on Windows, or should the probe (and the manifest vocabulary) say so
   explicitly in the docs?
2. The legacy carve-out above is silent: a unit validated under a `markerSchema`-less marker
   is indistinguishable in doctor's output from a fully anchored one. Do you want
   `doctor --json` to surface that distinction (e.g. a `markerSchema` field on the unit, or a
   `legacy: true` note) so automation can tell "validated against the anchor" from
   "grandfathered"? Visibility only — no health change. Fine to defer to a follow-up idea.
