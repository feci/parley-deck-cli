---
idea: integrate-parley-bidding-addon
review-round: 1
agent: kimi-1
date: 2026-07-30
reviewed-commit: 714712f
---

## Verdict

BLOCK

Two independent kimi-1 passes reached this verdict (see the note at the top of "What I
verified"). The branch advanced to `89069b0` ("fix-up cycle 1") while the review ran; at
diff level that commit addresses findings 1 and 2. Findings 3 (B5) and 5 are untouched at
`89069b0` and remain open at branch head.

## Findings

### [CRITICAL] An expected add-on gutted to `SKILL.md` reports `valid` when the install marker is gone too

**Where:** `lib/installer.js:1255-1259` (`manifestProblems` early return), `lib/installer.js:1315-1320` (`readMarker` collapses every failure mode to `null`), `lib/installer.js:1227` (add-on required files are `["SKILL.md"]` only) — all at 714712f.

**What:** The install marker (`.parley-deck-skill-install.json`) lives inside the add-on directory it anchors. Gut an installed `parley-bidding` down to `SKILL.md` and the marker is deleted along with the manifest and payload; `readMarker` then returns `null`, `manifestProblems` returns `[]`, and the unit passes the `SKILL.md`-only check: `doctor` reports `valid` and exits 0. The same false green occurs with a corrupt (unparseable) marker and with a marker naming another tool. The only marker states that are actually checked are readable, installer-owned, and current-schema.

**Why it matters:** This is the B3 defect the entire integrity mechanism exists to close, reproduced by deleting one more file — and "gut the tree and confirm `doctor` says `malformed`" is the contract's own negative test. The marker is a dotfile: `cp src/* dst` and most cleanup scripts drop it as a matter of course, so marker loss is not an exotic companion to payload damage. The shipped gutted-tree test (`test/bidding-addon.test.js:188-198`) deliberately preserves the marker, so it cannot see this. codex-1's ratified amendment condition 1 (`round-03/codex-1.md`) — "an expected installed unit with a missing or unreadable marker must also be unhealthy" — is unimplemented, and condition 4's required marker-deletion/corruption negative tests do not exist in the 23-test suite. The implementer's "Notes for reviewers" #3 also claims the legacy no-`markerSchema` path is "the one path where the manifest check is skipped entirely"; that is factually wrong — absent, corrupt, and foreign markers skip it too.

**Evidence:** Measured twice, independently, against 714712f (once against a pristine `git archive` export, once against a clean `git worktree` at the commit):

```
A) gutted to SKILL.md, marker deleted  → doctor.ok=true  status=valid  problems=[]
B) marker corrupt ("{ not json")       → doctor.ok=true  status=valid  problems=[]
C) marker.name = "other-tool"          → doctor.ok=true  status=valid  problems=[]
D) control: delete adapter_validate.py → doctor.ok=false status=malformed, names the file
```

Probe A was reproduced in both passes with identical output.

**Fix:** Preserve marker read state (absent vs unreadable vs foreign) instead of collapsing to `null`. For an expected unit — one the core marker's recorded selection or an explicit `--only` makes expected — a missing or unreadable marker is itself `malformed`, with distinct problem strings. Legacy compatibility stays only for a readable, installer-owned, genuinely-older marker. Add the marker-deletion and marker-corruption negative tests. (`89069b0` implements exactly this shape in `skillUnitStatus` via `readMarkerState`, with both tests — diff-verified, not yet re-run by me.)

### [MAJOR] B6: `doctor` never checks the declared Python floor

**Where:** at 714712f the floor probe lives only in `scripts/run-python-tests.js`; the `doctor` path (`doctorCommand` `lib/installer.js:362` → `targetStatus:1193` → `skillUnitStatus:1208` → `validateInstalledPayload:1224`) contains no reference to `python3`, `runtime`, or the interpreter — grep for `runtimeAvailability|probePython3|skill.runtime` in the 714712f installer is empty.

**What:** On a host with no `python3` (or one below the manifest's declared `>=3.10`), the add-on can be byte-perfect and is dead, yet `doctor` reports `valid` and exits 0. The floor is enforced where it gates tests, not where it gates operational health.

**Why it matters:** FINAL.md B6 assigns this to `doctor` explicitly: it must distinguish *payload-valid* from *operationally unavailable* and fail health when the declared interpreter minimum is missing. A test runner does not run on production installs; `doctor` does. No deviation records the gap.

**Evidence:** Probe E on the export: intact install with `PATH` empty (no `python3` anywhere) → `install.ok=true`, `doctor.ok=true`, `parley-bidding status=valid`, `problems=[]`. The manifest declares `runtime.python: ">=3.10"`; nothing in the health path reads it.

**Fix:** During `doctor`/`status`, read `runtime.python` from the validated manifest, probe the interpreter (once per process), and report operational unavailability as a distinct state from payload validity — payload still `valid`, unit unhealthy, `doctor` exit non-zero. Negative tests: interpreter absent, interpreter below floor. (`89069b0` adds `runtimeAvailability`/`probePython3` and wires `skill.runtime.ok` into `doctorCommand`'s exit — diff-verified.)

### [MAJOR] B5: preflight is per-target and manifest-blind — predictable failures still produce partial writes

**Where:** `preflightSkillUnit` (`lib/installer.js:897`), `installTarget` (`:927`, preflights only its own target's units), `installCommand` (`:593`, maps targets sequentially, writing as it goes).

**What:** Two vectors, one root cause — preflight does not cover "every unit and destination before the first write":

- **(a)** Manifest-free add-on sources are never traversed in preflight (`:911` gates source validation on `hasManifest`). A symlink in a late manifest-free add-on passes preflight and is refused by `copyRecursive` only after the core and every earlier add-on have been replaced. (Independently reproduces codex-1's finding; crediting him with first report.)
- **(b)** Targets are preflighted one at a time, interleaved with writes. With the contract's own example defect — an existing unmarked `parley-bidding` at a later target — every earlier target is fully installed before the later one is blocked. The run ends `ok:false` with writes already on disk.

**Why it matters:** B5: "Preflight every unit and destination before the first write; a predictable failure — an existing unmarked `parley-bidding`, for instance — must produce **zero** writes." IMPLEMENTATION.md's claim ("a corrupt source payload leaves the destination non-existent") is true only per destination. If per-target isolation is the intended reading, that deviates from the contract text and should have been recorded as a deviation; it was not. Side observation: the blocked target's top-level message is "Not attempted: another skill in this install failed preflight" because the result shape takes `skills[0]` (the core unit) rather than the actual blocker — the real reason only appears inside `skills[]`.

**Evidence:** Measured on the export (first pass) and reproduced on a clean worktree (second pass):

```
F) staged package, manifest-free zz-broken add-on containing a symlink
   → install.ok=false, destination written anyway:
     parley-deck, parley-bidding, parley-design, parley-design-check,
     parley-tracker, parley-worktrees all present; zz-broken absent
G) --target all, codex + qwen detected, hand-copied unmarked parley-bidding
   in qwen's skills dir
   → target codex: "installed"; target qwen: blocked at preflight;
     overall ok=false; codex parley-bidding written = true
G') --target all --include-undetected, unmarked parley-bidding planted at
    aionrs (last of 14 targets)
    → 13 targets report "installed", aionrs "skipped", overall ok=false;
      codex core + parley-bidding + manifest on disk = true
```

**Fix:** Preflight every selected unit of every resolved target before the first write (or stage all units and replace only after every stage validates). Add a read-only copyability traversal (symlink, readability) for manifest-free sources. Regression tests for both vectors. Surface the blocking unit's own message at target level.

### [MINOR] IMPLEMENTATION.md front matter recorded the wrong head commit — since amended

**Where:** `IMPLEMENTATION.md` front matter.

**What:** At the reviewed state it recorded `head-commit: a544dcd` — that is `main` ("release: parley-deck-skill 2.0.0"), not the implementation head under review (`714712f`). The record self-certified the wrong tree.

**Evidence:** `git log --oneline` on the skill repo; the file's front matter as read before the fix-up.

**Status:** Amended in fix-up cycle 1 — the file now records `head-commit: 89069b0` and `status: fix-up-cycle-1`. No further action; recorded here because the reviewed state was wrong.

### [MINOR] D-3's claim that `\` is still refused is false in the only path that exercises it

**Where:** `test/design-addons.test.js:1112` (`joined = command.replace(/\s*\\$/, "")` strips the extractor's splicing sentinel before `PUBLISHED_PYTHON` is applied; sentinel appended at `:431`); `skills/parley-bidding/SKILL.md:111-116`.

**What:** The extractor marks a backslash-continued command by re-appending ` \`; the Python arm removes that marker before matching, so continued commands are *accepted* — and the shipped `SKILL.md` relies on it, publishing two backslash-continued commands (`release_lint.py`, `completeness_lint.py`) that pass the suite. The grammar unit test asserting `PUBLISHED_PYTHON` refuses a trailing backslash is vacuous for extracted commands, because the sentinel is stripped first. So D-3's "`;` `|` `&` backtick `$` and `\` are still refused" and F5's "shell syntax refused" do not hold for the shipped docs. (Concurrent with codex-1's MINOR finding; verified independently. Line citations re-verified at 714712f in the second pass.)

**Why it matters:** The arm is static, so there is no execution hazard and accepting `<`/`>` placeholders remains justified. The harm is a false factual claim in the deviation record and a publication contract the test claims to enforce but doesn't.

**Evidence:** 278/278 node tests pass at 714712f including the Python arm, while `SKILL.md:111-116` ships two continued commands.

**Fix:** Either refuse the sentinel and publish those two commands on one line each, or ratify a narrow continuation exception and correct D-3, the F5 wording, and the arm's comment.

No NIT findings.

## What I verified and found correct

Two kimi-1 passes over the same commit, the second started before the first's output was
visible, so the headline findings and the counts below were derived twice independently
(once via `git archive` export, once via a clean `git worktree` at 714712f with
`npm ci`). Where a probe was run in both, it is marked "×2".

- **The worktree hazard itself.** The branch advanced from `714712f` to `89069b0` mid-review; one early probe in each pass accidentally measured the fixed code (it reported the gutted case `malformed`). Every number above and below is from 714712f, measured in isolation from the live tree.
- `npm test` at 714712f: **278 node tests pass / 0 fail**, then the Python leg **54/54** (4+20+2+3+15+3+7), then `build-addon-manifest.js --check` ok — matches IMPLEMENTATION.md. ×2
- Baseline reproduced: `node --test` at `a544dcd` (main) is **253 pass / 0 fail**, so "was 253" is exact, and 253 + 23 new + 2 guard-arm tests = 278 is consistent.
- Python floor claim "measured on 3.10, 3.11 and 3.14 — 54/54 on each" **reproduced**: 54/54 on 3.10.20 and 3.11.15 (run directly, `python3 -B`, `PYTHONDONTWRITEBYTECODE=1`, zero cache artefacts left) and on 3.14 via `npm test`.
- `npm pack --dry-run` at 714712f: **202 files total, 48 under `skills/parley-bidding/`, zero `__pycache__`/`.pyc`, no nested `.gitignore`, and `prepack` fires** (its `--check` output precedes the tarball listing). All as claimed.
- **The rename count is exactly right, and the contract was wrong.** In the source (`/Volumes/My Shared Files/AI_WORKSPACE/BYTE/software-bidding`): 9 `software-bidding` occurrences + 4 `Software Bidding` display-name occurrences = **13 across the same 8 files** — the "thirteenth is real" claim verifies, and FINAL.md's "twelve lines" was an undercount. The integrated tree contains zero occurrences of either variant; the only repo reference to the old name is the negative assertion in `test/bidding-addon.test.js`.
- **Source-vs-integrated diff as claimed:** source is 48 files, 0 caches, and genuinely untracked in the BYTE repo (`?? software-bidding/` — D-2's premise holds, so the porcelain check would indeed have proven nothing). 9 files differ from the integrated tree: the 8 rename files plus `references/parley-integration.md` (SKILL.md is both rename and consent), matching "9 files differ — 8 rename, 2 consent".
- **The consent paragraph is byte-exact in both files**, verified by unwrapping the FINAL.md blockquote programmatically (not by re-typing): present verbatim in `skills/parley-bidding/SKILL.md` and `references/parley-integration.md`.
- Manifest: 47 entries, excludes itself, aggregate `sha256:7854adf1…` as recorded, `runtime.python: ">=3.10"`, no `version`/`semver` field. Repository leg matches; I did not re-derive the tarball/binary install channels.
- With the marker intact, the mechanism works: gutted tree, four single-file deletions, single-byte mutation, stray `.pyc`, manifest deletion, manifest corruption, stripped marker field, newer marker schema, and a self-consistent manifest+payload swap all report `malformed` — via the 23 `bidding-addon` tests passing at 714712f plus control probe D. ×2
- **The guard refactor is faithful — verified behaviorally, not just by reading.** I loaded the `main` and 714712f extractors side by side and compared them on 13 adversarial inputs (`n$(printf '')ode --test x`, case variants, brace expansion, split spans, trailing operators, the python3 shape as a non-target): `mentionsATestCommand` produced **zero** divergent answers, and `publishedTestCommands` returned identical occurrence lists on a mixed code/prose markdown sample. Same patterns, same flags, fresh `RegExp` per call (no `/g` `lastIndex` carry-over). Parameterization weakened no assertion.
- **The Python arm's grammar is sound** where it matters: the script path must be a single plain filename under `scripts/` (no `/`, no leading dot — traversal refused), composition characters `; | & \` $` backtick are refused, and accepting `<`/`>` is justified because the arm never executes anything. SKILL.md publishes exactly five `python3 scripts/*.py` commands. The one false note is the backslash claim — finding 5.
- README states the runtime-exposure claim honestly as **NOT TESTED** (D-4: "installs into and validates in every runtime this installer covers — measured across all fourteen destinations" — consistent with probe G', where 13 targets installed cleanly and the 14th was my own planted blocker) and carries the six-skills wording (D-5). D-1's marker-anchored design matches what round-03 ratified — the defect is in its implementation (finding 1), not the design.
- Not verified by me: `npx skills add <repo> --list` finding six; the portable-binary build/install claims (49 files, `doctor` valid); the B7 aggregate on the tarball/binary/native channels.

## Open questions for the implementer

1. The branch advanced to `89069b0` ("fix-up cycle 1 — ready for re-review") while this review ran. At diff level it closes findings 1 and 2 (marker-state preservation in `skillUnitStatus`, both negative tests, `runtimeAvailability` wired into `doctor`'s exit). Should re-review target that head? It does not touch `installCommand`/`installTarget` or the guard arm — is a second fix-up cycle planned for findings 3 and 5?
2. Is per-target preflight the intended reading of B5? If yes, record it as a deviation (D-6) and narrow the "zero writes" claim to per-destination scope; if no, finding 3 stands as written, and it stands at `89069b0` too.
3. The fix at `89069b0` makes marker-missing and marker-unreadable distinct problem strings — confirmed in the diff. Will the re-review baseline (`npm test` at the fix-up head) be recorded in IMPLEMENTATION.md alongside them?
4. Finding 4 is already moot (the record now says `head-commit: 89069b0`). For cycle 2, please keep the head commit and the reviewed state in lockstep as further fix-ups land.
