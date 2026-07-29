---
idea: skills-cli-install-path
implementer: claude-1
date: 2026-07-29
status: fix-up-cycle-6
target: parley-deck-skill
head-commit: 46b5730
prior-commits: [951d7a5 move+installer+panel, f8e3a1c gemini path, 085799e cycle-1, a05bac7 cycle-2, bddbf1a cycle-3, fa1fdb1 cycle-4, 4f7fd32 cycle-5]
---

## What was done

The layout move from `FINAL.md`, the installer repointing, and the README panel.

```
skills/parley-deck/{SKILL.md, references/, agents/}   ← was the repository root
skills/parley-design/  parley-design-check/  parley-tracker/  parley-worktrees/
                                                      ← was addons/*
```

`addons/` is gone. A **rename**, not a copy — no second skill tree, so no drift guard is
needed (S4b). `plugin.json` and `gemini-extension.json` stayed at the repository root per A3,
with the installer's *source* path for them set explicitly.

### Installer (`lib/installer.js`)

`CORE_SKILL_DIR` and `CORE_SKILL_NAME` introduced. `REQUIRED_PAYLOAD_FILES` and
`COMPATIBILITY_FILE` resolve under it. `PAYLOAD_ENTRIES` / `OPTIONAL_PAYLOAD_ENTRIES` became
`{from, to}` pairs — **the source moved, the destination shape did not**, so an installed skill
still has `SKILL.md` at its root and every per-target validator is unchanged. `ADDONS_DIR`
`"addons"` → `"skills"`, with `discoverAddons` skipping `parley-deck` so the core is never
offered as an add-on. `packagedProtocolPath` and the core `skillSha256` repointed. The
Antigravity staging that fabricates `skills/SKILL.md` in the destination now reads from the
new source.

### Packaging

`package.json` `files` ships `skills/` and no longer `addons/` or a root `SKILL.md`.
`pkg.assets` (which the standalone binaries embed) rewritten the same way.

### Tests

`test/design-addons.test.js` repointed. `test/installer.test.js`'s synthetic package now
builds the new layout. **No test assertion was weakened** — only the paths they build and read.

## Gates

| Gate | Result |
|---|---|
| **G1** `install --target all`, `status --target all` | **PASS** — rc 0; 30 `valid` lines; the only non-valid line is `project: missing protocol`, which is expected (this repo is not a Parley deck) |
| **G2** `npm pack` ships `skills/`, not `addons/` or a root `SKILL.md` | **PASS** — 145 `skills/` entries, `skills/parley-deck/SKILL.md` present, 0 `addons/` entries, 0 root `SKILL.md` |
| **G3** Windows binary build + install | **NOT TESTED** — no Windows host. `pkg.assets` is updated but unproven |
| **G4** legacy Gemini `contextFileName` | **PASS for the supported path** — `install --target gemini` rc 0 and `~/.gemini/extensions/parley-deck/SKILL.md` exists, so `contextFileName: "SKILL.md"` still resolves. See the deviation below for the unsupported path |
| **G5** `agy plugin validate` | **PASS** — `✔ skills: 1 processed` |
| **G6** `npm test` | **PASS** — 247 pass, 0 fail |
| **G7** real install, no `--full-depth`, isolated `HOME` | **PASS, and it is the point of the idea.** `Found 5 skills` / `Installed 5 skills` with **no flag**. Installed core contains exactly `SKILL.md`, `agents`, `references`; `bin`, `lib`, `package.json`, `test`, `dist` are all **absent**. Both the discovery defect (S1) and the whole-repo-as-skill defect (S3) are fixed and *measured*, not asserted |
| **G8** `skills update` | **NOT TESTED** — `update` reported "No project skills to update" after a `--copy` install, so the run proved nothing about our layout. kimi-1 root-caused it: `skills` skips `sourceType: "local"` by design, so a filesystem-path install has nothing to update from |
| **G9** Homebrew `brew upgrade` / `brew test` | **NOT TESTED** — the formula points at a published tarball; untestable before the release |
| **G10** WinGet `winget validate` / `winget install` | **NOT TESTED** — `winget` is not installed on this host (macOS) |
| **G11** `gemini extensions install <url>` | **NOT TESTED and now unsupported** — see below |

G7 verbatim:

```text
$ HOME=/tmp/g7/home npx -y skills@latest add <repo> --agent claude-code --yes --copy
◇  Found 5 skills
◇  Installed 5 skills
     ✓ parley-deck  ✓ parley-design  ✓ parley-design-check  ✓ parley-tracker  ✓ parley-worktrees
$ ls .claude/skills/parley-deck/     → SKILL.md  agents  references
$ test -e .claude/skills/parley-deck/{bin,lib,package.json,test,dist}  → all absent
```

## Deviation D-1 — the legacy Gemini extension-URL path is dropped, not fixed

kimi-1 called `gemini-extension.json`'s `contextFileName` "the sharpest single edge", and it
is. The value is relative to the extension directory, and the two consumers now disagree:

- `install --target gemini` copies `SKILL.md` to the destination root → `"SKILL.md"` resolves.
- `gemini extensions install <repo-url>` treats the **repository** as the extension → there is
  no root `SKILL.md` any more → it would resolve to nothing.

**One value cannot satisfy both.** I could not test the second path (`gemini` is not installed
here), and shipping a documented command that my own analysis says is broken is exactly what
`parley-design`'s honesty rule forbids. So the README line for it is **removed**, and the
README now states plainly that installing this repository through
`gemini extensions install <url>` is not supported and why. The supported `--target gemini`
path is verified working.

This is a **capability we lost**, not a defect we fixed, and it is the real price of the
restructure. It is written here rather than left for someone to discover.

## What reviewers should attack

1. The four `NOT TESTED` gates. None may be reported as a pass.
2. D-1 — is dropping the extension-URL path the right call, or should the restructure be
   reshaped to keep it?
3. Whether any per-target *destination* shape changed. I claim none did, because only
   `PAYLOAD_ENTRIES.from` moved. `validateInstalledPayload` is untouched — verify that.
4. The README panel's wording against F1/F3/F4/F5.

---

## Fix-up cycle 1

Review round 01: **codex-1 ❌ BLOCK** (1 MAJOR, 3 MINOR), **hermes-1 ✅ ACCEPT**.

### The MAJOR overturns my deviation D-1, and codex-1 is right

I claimed one `contextFileName` could not serve both consumers, and dropped the legacy Gemini
extension-URL channel on that basis. codex-1 refuted it with the upstream source: Gemini's
extension manager joins `contextFileName` to the extension root and **accepts nested paths**.
So the reconciliation it proposed works, and it is strictly better than what I did:

- **repository** `gemini-extension.json` → `"skills/parley-deck/SKILL.md"`, which is where the
  file actually is in a checkout, so an extension installed from the repo URL resolves;
- the native installer **rewrites the staged copy** to `"SKILL.md"`, matching the flat
  destination shape it has always produced.

One canonical skill tree, two consumers, no second copy of `SKILL.md`. **D-1 is withdrawn and
the channel is restored.** I had reached for "remove the capability" when "reconcile the two
consumers" was available — the removal was disclosed honestly, but honest disclosure of an
unnecessary loss is still an unnecessary loss.

Two regression tests were added, exactly as codex-1 asked: one asserting the repository
manifest points at a file that exists **in the repository**, one asserting a staged install's
manifest points at a file that exists **in the destination**. The second caught a real bug
immediately: my first helper used `readJsonFile`, which returns a `{status, value}` wrapper,
so it wrote the wrapper back and destroyed the manifest. Fixed to parse directly.

| Finding | Action |
|---|---|
| MAJOR — D-1 removes a reconcilable channel | **Withdrawn.** Manifest reconciled both ways, channel restored, 2 tests added |
| MINOR — G1 recorded as unconditional pass but is environment-dependent | Gate restated below with its preconditions. codex-1 found that in a clean `HOME` the Antigravity install is what creates the evidence that then detects Gemini — an **ordering defect that predates this change** and is recorded as a follow-up, not fixed here |
| MINOR — README ships into destinations where its nested paths do not exist | Now states both shapes explicitly: checkout paths *and* installed-destination paths |
| MINOR — "whichever coding agents you have" exceeds measured scope | → "the coding agents it supports … theirs to state, not ours" |

### Verification

```text
$ npm test                              → pass 249, fail 0   (2 new tests)
$ install --target gemini --force       → rc 0; staged contextFileName "SKILL.md", file present
$ cat gemini-extension.json             → "skills/parley-deck/SKILL.md", file present in repo
$ G7 re-run, no --full-depth            → Found 5 / Installed 5; core = SKILL.md + agents +
                                          references; bin, lib, package.json all absent
```

### Gate status, restated honestly

**G1 — PASS with preconditions.** On this host, with runtimes already present,
`install --target all --force` → rc 0 and 30 `valid` lines. codex-1 is right that in a clean
`HOME` the sequence is order-dependent; the deterministic form is
`--target all --include-undetected`. Not a layout regression — reproducible on the parent
commit.

**Still NOT TESTED, unchanged:** G3 Windows binary · G8 `skills update` · G9 Homebrew ·
G10 WinGet · **G11 `gemini extensions install <url>`** — the manifest now supports it and two
tests guard both values, but no Gemini CLI is available here to run it end to end. That is
stated in the README itself, not only here.

### New follow-up

Antigravity-before-Gemini detection ordering in a clean `HOME` (pre-existing, found by codex-1).

---

## Fix-up cycle 2

Review round 01 completed with all three reviewers: **codex-1 ❌ BLOCK** (1 MAJOR, 3 MINOR),
**kimi-1 🟡 ACCEPT-WITH-RESERVATIONS** (1 MAJOR, 1 MINOR, 2 NIT), **hermes-1 ✅ ACCEPT**.

### kimi-1's MAJOR is a miss of mine, and the kind that would have shipped

**I renamed the directory and did not sweep the rename through the shipped skill content.**
Six files still told the reader to run `addons/…` paths that no longer exist — including
`parley-design-check`'s own documented run command,
`node addons/parley-design-check/bin/check.js`, which kimi-1 proved broken. Every filled
exemplar in `parley-tracker/templates/` had the same class of defect in its `files:`
frontmatter and its `Verify:` commands.

A layout move is not done when the tests pass. It is done when the **instructions the skills
give their readers** still resolve. The installer had no reason to fail — none of these paths
is code the installer runs — so nothing in my gate list could have caught it. kimi-1 found it
by reading what the skills tell people to type.

| Finding | Action |
|---|---|
| **MAJOR-1** — `addons/…` references stale in shipped skill content | Swept all 6 files. `grep -rn "addons/" skills/` → **0**. Verified the two commands that were provably broken now resolve |
| **MINOR-1** — same class in the tracker's filled exemplars | Same sweep. **Re-ran the exemplars' own validator**: `validate.js --strict --dir templates` → *All 3 ticket(s) passed*. The exemplars still self-pass, which is the property they exist to have |
| **NIT-1** — panel comparative unmeasured | Already reworded in cycle 1 to "a longer list than this package's own installer covers, and theirs to state, not ours". Verified absent |
| **NIT-2** — G8's reason imprecise | Corrected above using kimi-1's root cause |

### G8, root-caused by kimi-1 rather than hand-waved

My note said "upstream bookkeeping". The actual mechanism: `skills` skips
`sourceType === "local"` by design, and `skills-lock.json` records `"local"` for all five
because the source was a filesystem path — a local copy has nothing to update *from*. **G8 is
therefore only meaningful post-merge against the published remote.** That is a different
re-test from the one my note implied, which is why the imprecision mattered.

kimi-1 also confirmed G9 and G10 are genuinely untestable here for reasons independent of this
change: `packaging/homebrew/Formula/` in this repo is empty (pre-existing), and the WinGet
manifests reference GitHub-release binaries rather than repo paths, so the move does not
invalidate them.

### Verification

```text
$ grep -rn "addons/" skills/                                   → 0
$ node skills/parley-tracker/bin/validate.js --strict --dir skills/parley-tracker/templates
                                                               → All 3 ticket(s) passed, rc 0
$ test -f skills/parley-design-check/bin/check.js               → present, matches its SKILL.md
$ npm test                                                      → pass 249, fail 0
```

### Remaining NOT TESTED

G3 Windows · G8 `skills update` (post-merge, against the remote) · G9 Homebrew (pre-release) ·
G10 WinGet (no host) · G11 `gemini extensions install <url>` (no Gemini CLI). None is reported
as a pass.

---

## Fix-up cycle 3

Review round 02: **codex-1 ❌ BLOCK** (1 MAJOR). codex-1's review ran against the cycle-1 tree,
before the cycle-2 sweep, so the 34 stale paths it counted were already gone by the time it
filed. **But it asked for something cycle 2 did not do, and it was right to.**

### The false-green codex-1 found

`skills/parley-design-check/SKILL.md:372` documented its own test command. With the stale
path it **exited 0 while running zero tests** — a verification path that reports success by
matching nothing. That is worse than a broken command, because a broken command is loud.

Verified after the sweep:

```text
$ node --test "skills/parley-design-check/test/*.test.js"
ℹ tests 159   ℹ pass 159   ℹ fail 0
```

### The guard, which is the actual deliverable of this cycle

Cycle 2 fixed the paths. It added **nothing that would notice them coming back** — and
`npm test` structurally cannot notice, because the moved files are valid; only their
*instructions* were wrong. Two assertions added to `test/design-addons.test.js`:

- every shipped `.md` / `.js` / `.json` / `.yaml` under `skills/` is scanned for a live
  `addons/` path, reporting **file and line** for each offender;
- the package ships no `addons/` directory at all.

**The guard was proved to fail before it was trusted.** Reverting one line of
`skills/parley-tracker/SKILL.md:79` back to `addons/parley-tracker` drops the suite to
`pass 250, fail 1` with the offending file and line named; restoring it returns
`pass 251, fail 0`.

codex-1 also warned against restoring an `addons/` compatibility tree to make the old paths
work. Not done — that would reintroduce the second tree S4b exists to prevent.

### Verification

```text
$ grep -rn "addons/" skills/                                    → 0
$ npm test                                                       → pass 251, fail 0
$ npm test, with one addons/ path reintroduced                   → pass 250, fail 1 (guard fires)
$ node --test "skills/parley-design-check/test/*.test.js"        → 159 tests, 159 pass
```

### Lesson worth keeping

Three reviewers looked at this change. The installer was correct from the first commit; both
real defects were in **what the shipped files tell a human to type**. Neither the test suite
nor any of my seven gates was capable of seeing them. kimi-1 found them by reading the
instructions; codex-1 found the false-green by *running* one.

---

## Fix-up cycle 4

Review round 03: **codex-1 ❌ BLOCK** (1 MAJOR), **kimi-1 ✅ ACCEPT** (all four of its round-01
findings verified FIXED by running them, not by reading my claims).

### codex-1 found a second instance of the class — and one my cycle-3 guard could not see

`skills/parley-tracker/templates/subtask.md:68` and `:74` publish
`node --test skills/parley-tracker/bin`. On Node v26.5.0 that directory form fails with
`MODULE_NOT_FOUND`: **one failed harness entry, zero passing tests.** It is a filled exemplar,
so its acceptance criterion and Definition of Done tell an implementer that this command
proves the tests are green. It proves nothing, and the box cannot be ticked honestly.

**The cycle-3 guard was blind to it by construction.** That guard greps for `addons/`. This
path is perfectly correct — the *form* is wrong. Fixing symptoms one at a time was never going
to converge; the guard had to be for the class.

Fixed to `node --test "skills/parley-tracker/bin/*.test.js"` → **35 tests, 35 pass**.

### The real deliverable: a guard for published verification commands

New test: extract **every** `` `node --test …` `` command from every shipped `.md` under
`skills/`, run each one, and assert both `fail == 0` **and `pass > 0`**. A command that exits
0 while running nothing now fails the suite.

Two bugs of mine surfaced while building it, and both are worth recording because each would
have produced a guard that silently passed:

1. **The summary parser looked for `# pass N`.** Node's default reporter prints `ℹ pass N`.
   My first version read zero and would have been "fixed" by relaxing the assertion — the
   wrong direction entirely. It now accepts either form and **treats an unparseable summary as
   a failure**, not as zero.
2. **A test runner spawned inside a test runner emits nothing to stdout.** It inherits
   `NODE_TEST_CONTEXT` and reports through the parent, so the child's output was empty. The
   child env is now stripped of `NODE_TEST*` so it behaves exactly as it does for a person
   typing the published command.

**Proved to fail before being trusted**, like the cycle-3 guard: reverting the exemplar to the
directory form drops the suite to `pass 251, fail 1`; restoring it gives `pass 252, fail 0`.

### Verification

```text
$ node --test "skills/parley-tracker/bin/*.test.js"        → 35 tests, 35 pass, 0 fail
$ node --test "skills/parley-design-check/test/*.test.js"  → 159 tests, 159 pass, 0 fail
$ npm test                                                  → pass 252, fail 0
$ npm test, exemplar reverted to the directory form         → pass 251, fail 1 (guard fires)
```

### Findings per round

`r01: 8 (3 reviewers) · r02: 1 · r03: 1`. Every one of them was in shipped *instructions*,
never in the installer.

---

## Fix-up cycle 5

Review round 04: **codex-1 ❌ BLOCK** (1 MAJOR).

### My cycle-4 guard did not do what cycle 4 claimed it did

I wrote that it runs "every published `node --test` command". It ran the **inline** ones. The
extractor was `` /`(node --test [^`]+)`/g `` — a single-backtick code span — and
`skills/parley-design-check/SKILL.md:371-373` publishes its command inside a **fenced block**.
So the 159-test command, the one whose false green started this whole thread, was never
checked by the guard built to protect it.

codex-1 did not argue this; it measured it in scratch copies:

- a broken **inline** command → `npm test` 251/1 (caught)
- the same broken command **fenced** → `npm test` 252/0 (**missed**)

That is the third time in this idea that a guard or claim of mine was weaker than its
description. The pattern is consistent and worth stating plainly: **I keep verifying the
instance I just fixed rather than the class I claimed to cover.**

### What changed

- Extraction pulls commands from **both** inline code spans and fenced blocks, strips a
  leading `$ ` prompt, normalises whitespace and deduplicates.
- The extractor is now a **named function tested against a fixture** containing one valid and
  one broken command in *each* form, plus a non-test line that must be ignored. An extractor
  that silently misses a form yields a guard that silently passes — so the extractor itself
  needed a test, not just the thing it feeds.
- The runner assertion now requires `published.size >= 2`, so the suite fails if discovery
  ever regresses to finding only one command again.

### Proved, both forms

```text
$ npm test                                                   → pass 253, fail 0
$ break the FENCED design-check command  → pass 252, fail 1   (was 252/0 before this cycle)
$ break the INLINE tracker command       → pass 252, fail 1
$ both restored                                              → pass 253, fail 0
```

### Findings per round

`r01: 8 (3 reviewers) · r02: 1 · r03: 1 · r04: 1`. Every finding in this idea has been about
shipped instructions or about a guard that was narrower than its claim — none about the
installer, which has been correct since the first commit.

---

## Fix-up cycle 6 — the enumeration is abandoned

Review round 05: **codex-1 ❌ BLOCK** (1 MAJOR). Tilde fences (`~~~`) are ordinary markdown
fences and the cycle-5 extractor read only backtick fences. codex-1 measured it: a broken
command in a `~~~` block left the suite at **253/0**, while the same command in a ``` block
gave 252/1.

**Adding `~~~` would have been the same mistake in a new costume.** Four rounds, four
enumerations, each narrower than its claim:

| cycle | what the guard enumerated | what it missed |
|---|---|---|
| 3 | the string `addons/` | a correct path in a broken *form* |
| 4 | inline code spans | fenced blocks |
| 5 | inline spans + backtick fences | tilde fences |

Fence syntax is an **open set** — backtick, tilde, longer runs, indented blocks, HTML. An
enumeration of an open set is incomplete by definition, so the next round would have found the
next member. This is the same lesson the `parley-design-check` work reached the hard way:
patching a family member-by-member never converges; you have to close the class.

### What cycle 6 does

The extractor **no longer parses markdown structure at all**. It scans every line for
`node --test `, takes the command, and terminates it at the closing backtick if it is inside
an inline span. There is no container to miss because containers are never consulted.

A false positive is harmless and arguably correct here: if a shipped file prints a
`node --test` command *anywhere*, in any context, that command should work.

Two defects in my own first line-scan were caught by the fixture before commit: it ran past
the closing backtick of a mid-sentence inline span, and it swallowed `` `) — COMMIT-SHA `` from
a real Definition-of-Done line. Both fixed by terminating at the backtick.

### Proved against five container forms, not argued

Each form inserted into a shipped skill file in turn, with a deliberately missing test file:

```text
~~~ tilde fence            → pass 252, fail 1     (was 253/0 in cycle 5 — codex-1's case)
``` backtick fence         → pass 252, fail 1
    indented block         → pass 252, fail 1     (never tested before)
inline span, mid-sentence  → pass 252, fail 1
```` four-backtick fence   → pass 252, fail 1     (never tested before)
restored                   → pass 253, fail 0
```

The fixture test now asserts the extractor finds commands in all of those forms plus prose,
and ignores a non-`node --test` line.

### Findings per round

`r01: 8 (3 reviewers) · r02: 1 · r03: 1 · r04: 1 · r05: 1`. Not one was in the installer.
