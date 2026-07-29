---
idea: skills-cli-install-path
implementer: claude-1
date: 2026-07-29
status: fix-up-cycle-13
target: parley-deck-skill
head-commit: 82507b5
prior-commits: [951d7a5 move+installer+panel, f8e3a1c gemini path, 085799e cycle-1, a05bac7 cycle-2, bddbf1a cycle-3, fa1fdb1 cycle-4, 4f7fd32 cycle-5, 46b5730 cycle-6, c642636 cycle-7, cycle-8]
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

---

## Fix-up cycle 7

Review round 06: **codex-1 ❌ BLOCK** (1 MAJOR). I asked it to try to defeat the cycle-6 guard.
It did, on the first attempt, with one line:

```markdown
first run `node --test "…/bin/*.test.js"`; then run `node --test definitely-missing-second.test.js`
```

The guard called `line.indexOf("node --test ")` **once per line** and never resumed scanning.
First command valid, second broken, suite green at 253/0. Published alone on its own line, the
same broken command turns it red — so the bypass was purely positional. codex-1 also noted the
fixture could not have exposed it, because every fixture case held at most one command per line.

### The fix is also a simplification

Line-by-line `indexOf` is replaced by **one global match over the whole document**:

```js
markdown.matchAll(/node --test ([^`\n]*)/g)
```

A command runs to the first backtick — an inline span's close — or end of line, whichever comes
first. This drops the line loop, the manual index arithmetic and the separate backtick-split,
and it is strictly more complete: every occurrence, every line, every container.

That is the fifth guard revision in this idea, and the first one that got *smaller*. Each
earlier revision added a case; this one removed the machinery that made cases possible.

### Fixture

Extended with two commands on one line — the shape that hid the bypass — asserting both are
found.

### Proved

```text
$ codex-1's exact probe (valid first, broken second, one line)  → pass 252, fail 1
   (before this cycle: 253, 0 — the bypass)
$ ~~~ tilde fence            → pass 252, fail 1
$     indented block         → pass 252, fail 1
$ ```` four-backtick fence   → pass 252, fail 1
$ restored                                                      → pass 253, fail 0
```

### Findings per round

`r01: 8 (3 reviewers) · r02: 1 · r03: 1 · r04: 1 · r05: 1 · r06: 1`. Six consecutive rounds
where codex-1 found exactly one real defect, five of them in the guard rather than in the
product. The installer has not been touched since the first commit.

---

## Fix-up cycle 8 — the guard finally models what it claims to run

Review round 07: **codex-1 ❌ BLOCK** (1 MAJOR), and it is the sharpest framing of the whole
idea: *"a mismatch between the claimed shell command surface and a handwritten
extractor/runner that supports only one exact whitespace form and one target argument."*

Two measured defects, in opposite directions:

1. **A missed break.** A literal **tab** between `--test` and the target. A shell treats it as
   ordinary argument whitespace, so the published command is valid syntax and fails when typed.
   My regex required one literal ASCII space, so the guard never saw it: **253/0**.
2. **A manufactured failure.** `node --test a.test.js b.test.js` — a legitimate two-target
   command that runs 35 tests when typed. The runner passed both paths as **one** argument, so
   Node reported `Could not find "claim.test.js validate.test.js"` and the suite failed
   **252/1**. The guard was blaming the product for its own bug.

The second is the more dangerous class: a guard that produces false failures gets relaxed, and
a relaxed guard is how the original false-green comes back.

### What changed

- **Whitespace:** `/node\s+--test\s+([^`\n]*)/g`. Any whitespace, any amount, tabs included.
- **Argument boundaries:** the captured text is tokenised into real `argv` — whitespace
  separates, double quotes group — and passed to `execFileSync` as separate arguments. Quoted
  globs still arrive as one argument; multiple targets arrive as several.
- **Fixture:** now carries a tab-separated command and a multi-target command, and asserts both
  are discovered.

### Proved, both directions

```text
$ tab between --test and a missing target  → pass 252, fail 1   (before: 253/0, missed)
$ legitimate two-target command            → pass 253, fail 0   (before: 252/1, false failure)
$ restored                                 → pass 253, fail 0
```

A transcription slip while editing the fixture — a dropped comma between two array elements —
broke the file into a `SyntaxError` and took the whole suite from 253 to 241. Caught by
running it, fixed, recorded here rather than tidied away.

### Findings per round

`r01: 8 (3 reviewers) · r02: 1 · r03: 1 · r04: 1 · r05: 1 · r06: 1 · r07: 1`.

---

## Fix-up cycle 9 — the hand-written shell parser is deleted

Review round 08: **codex-1 ❌ BLOCK** (1 MAJOR). Both directions again, both measured:

1. **False green.** `node --test "…/bin/*.test.js".` — a trailing period the shell concatenates
   onto the target. Typed, Node exits 0 having run **zero tests**. My normaliser stripped the
   period as "prose punctuation", repaired the target back into the valid glob, and the suite
   stayed **253/0**. The guard was *manufacturing* the exact false green it exists to catch.
2. **Manufactured failure.** `node --test 'skills/parley-tracker/bin/*.test.js'` — legitimate,
   runs 35 tests. My lexer understood only double quotes, so the single quotes survived into
   the argv value, Node matched nothing, and the suite failed **252/1**.

### The fix codex-1 prescribed, and why it is the last one

*"Fix it without another case-by-case parser extension: capture command text without silently
rewriting it, and use a fail-closed exact-command-to-argv registry (or an actual, explicitly
scoped shell) so an unrecognised form fails as unsupported instead of being guessed."*

So:

- **Capture verbatim.** No stripping, no whitespace normalisation, no re-quoting. The command
  stored is the command published.
- **Execute through `/bin/sh -c`.** Quoting, whitespace and trailing characters are interpreted
  by a real shell — exactly as they are for a person who copies the line and presses enter.
  There is no approximation left to be wrong.
- **Fail closed.** Anything containing `;`, `|`, `&`, `<`, `>` or `$(` is **refused before
  execution** with a named reason, rather than guessed at or quietly run.

Six revisions of this guard hand-wrote progressively better approximations of shell syntax.
Every one of them was wrong in both directions at once. The correct move was to stop
approximating and delegate to the thing that defines the semantics.

### Proved

```text
$ trailing-period command    → pass 252, fail 1   (before: 253/0 — the false green)
$ single-quoted glob         → pass 253, fail 0   (before: 252/1 — the false failure)
$ restored                   → pass 253, fail 0

fail-closed, refused before execution and named:
$ node --test a.test.js; rm -rf /tmp/nope   → pass 252, fail 1
$ node --test $(echo evil).test.js          → pass 252, fail 1
$ node --test a.test.js | grep x            → pass 252, fail 1
  message: "published command uses shell syntax this guard refuses to interpret: …"
```

### Findings per round

`r01: 8 (3 reviewers) · r02: 1 · r03: 1 · r04: 1 · r05: 1 · r06: 1 · r07: 1 · r08: 1`.
Eight consecutive single-finding rounds, seven of them in the guard. The product — the layout
move, the installer, the packaging — has needed no change since round 01.

---

## Fix-up cycle 10 — the extraction boundary, which was the real bug all along

Review round 09: **codex-1 ❌ BLOCK** (1 MAJOR, 1 MINOR, 1 NIT).

Cycle 9 delegated *parsing* to a real shell but kept capturing from the word `node` onward,
discarding everything around it. **The guard was executing a substring of the published
command.** Three measured cases, two directions:

| probe | typed by a person | the guard |
|---|---|---|
| ``node --test `printf missing.test.js` `` | exits 1, file not found | stopped at the substitution backtick, dropped the fragment, **253/0 green** |
| `NODE_OPTIONS='--require ./missing.cjs' node --test "…"` | exits 1, preload missing | discarded the assignment, ran the valid suffix, **253/0 green** |
| `cd skills/parley-tracker/bin && node --test "*.test.js"` | 35 tests, 35 pass | discarded `cd … &&`, ran the suffix from the repo root, **252/1 false failure** |

The `&` refusal added in cycle 9 never fired, because the ampersands were **outside the
captured substring**. A fail-closed check that inspects only what the extractor kept is not
fail-closed at all.

### The fix

Extraction now yields **whole command units**, never fragments: the content of each inline
code span, or the whole line where there are none. Then a single strict grammar decides:

```js
const SUPPORTED_COMMAND = /^node\s+--test\s+[^`;|&<>$]+$/;
```

The command, its targets, and nothing else. Any environment prefix, `cd … &&`, command
substitution, pipe or trailing operator is **refused by name**, not guessed at and not
silently skipped:

> published command is not a bare `node --test <targets>` form, so this guard refuses to
> interpret it rather than execute a fragment of it: `cd skills/parley-tracker/bin && node --test "*.test.js"`

Narrow and fail-closed beats broad and approximate. That is the whole lesson of rounds 03–09.

### Proved

```text
$ backtick substitution   → pass 252, fail 1   (was 253/0 — false green)
$ env-prefixed command    → pass 252, fail 1   (was 253/0 — false green)
$ cd … && command         → pass 252, fail 1   (was 252/1 for the WRONG reason; now a
                                                named refusal, not a fabricated zero-test)
$ restored                → pass 253, fail 0
```

Also fixed this cycle:

- **MINOR** — the install panel said the universal installer "detects your agents" and three
  lines later claimed detection as something our installer adds. The G7 output itself prints
  `Agent detected`. Detection is dropped from the exclusive list; `doctor`/`status` health
  checks and project-metadata sync remain, and those are measured.
- **NIT** — cycle 9 claimed the hand-written argv lexer was deleted. It was still there as dead
  code with a stale comment. Now actually deleted, so the record and the code agree.

### Findings per round

`r01: 8 (3 reviewers) · r02: 1 · r03: 1 · r04: 1 · r05: 1 · r06: 1 · r07: 1 · r08: 1 · r09: 3`.

---

## Fix-up cycle 11 — a backtick means two different things, and the guard knew only one

Review round 10: **codex-1 ❌ BLOCK** (1 MAJOR). Both directions, both measured:

1. **False green.** Inside a fenced `bash` block:
   ``node --test `printf %s --test-reporter=definitely-missing-reporter` "…/*.test.js"``.
   Typed verbatim it exits 7 with `ERR_MODULE_NOT_FOUND`. The unitizer treated the backtick
   pair as a Markdown span, **deleted the substitution**, and executed the repaired remainder —
   35 passing tests, suite green at **253/0**. It repaired a broken command into a working one,
   which is the false-green class cycle 10 claimed to have closed.
2. **Manufactured failure.** `` Run ``node --test "…/*.test.js"``. `` — a standard CommonMark
   **double-backtick** span. The unitizer read the doubled backticks as two separate
   delimiters and reconstructed `Run  node --test "…" .`, then rejected it as not a bare
   command.

### The missing fact

Whether a backtick delimits a Markdown span or is shell syntax depends on exactly one thing:
**whether the line is inside a fenced block.** The unitizer never tracked fence state, so it
applied span rules to shell text and shell text to spans.

Now it does:

- it walks the document tracking the open fence marker (``` or `~~~`, any length ≥ 3);
- **inside** a fence the line *is* the command, verbatim — no span parsing, so a command
  substitution survives into the unit and is then refused by the strict grammar;
- **outside** a fence, a run of N backticks opens a span that a matching run of N closes, so
  ``` ``…`` ``` is one span rather than two delimiters.

### Proved — the two new probes, and all nine earlier ones re-run

```text
$ fenced command substitution        → pass 252, fail 1, refused by name  (was 253/0)
$ double-backtick inline span        → pass 253, fail 0                   (was a false failure)

regression sweep, unchanged behaviour:
  backtick fence / tilde fence / indented / two-per-line / tab / trailing period
                                     → pass 252, fail 1   (each)
  single-quoted glob / two targets   → pass 253, fail 0   (each, legitimately)
  cd … &&                            → pass 252, fail 1   (named refusal)
  restored                           → pass 253, fail 0
```

The fixture now carries both new forms and asserts the grammar's verdict on each.

### Findings per round

`r01: 8 (3 reviewers) · r02–r08: 1 each · r09: 3 · r10: 1`.

---

## Fix-up cycle 12 — the fence tracking is deleted, because it was never the right question

Review round 11: **codex-1 ❌ BLOCK** (1 MAJOR). Two more CommonMark rules my cycle-11 fence
tracker got wrong: a closing fence may be followed only by spaces or tabs (so
` ```not-a-closing-fence ` is *content*), and a fenced block inside a blockquote carries a `>`
prefix. False green in the first case, false failure in the second.

**Adding those two rules would have bought one more round.** I was reimplementing CommonMark
inside a test, one reviewer-found rule at a time.

### The question was wrong

I had been asking *"is this line inside a fence?"* — which requires a Markdown parser. The
question that actually separates the two meanings of a backtick is much smaller:

> **Does a backtick span contain the *whole* command?**

- Inline publication wraps the whole command: `` `node --test "x"` `` → the span *is* the
  command.
- A fenced shell line wraps only a substitution: ``node --test `printf …` "x"`` → **no span
  contains `node --test`**, so the unit is the whole line, backticks included, and the strict
  grammar refuses it.

One rule. No fence state, no closing-fence rule, no blockquote rule, no CommonMark. The
unitizer got smaller for the second time, and this time it lost an entire category of bug
rather than a case.

### Proved — all four probes from rounds 10 and 11

```text
$ fake closing fence + substitution   → pass 252, fail 1   (refused)
$ blockquoted fence, valid command    → pass 253, fail 0
$ fenced substitution                 → pass 252, fail 1   (refused)
$ double-backtick inline span, valid  → pass 253, fail 0
$ restored                            → pass 253, fail 0
```

**A correction to my own testing.** My first run of these four reported two regressions. They
were not regressions — my probe harness passed the fixture through `printf "%b"`, which mangled
the `%s` inside the command, so the file under test was not the file I meant to write. Running
the extractor directly against an exact heredoc showed the correct verdicts, and re-running the
probes with `cat` instead of `printf` confirmed all four. **The tool was right and my
measurement was wrong** — recorded here because a fix-up log that hides a bad measurement is
worth less than no log.

### Findings per round

`r01: 8 (3 reviewers) · r02–r08: 1 each · r09: 3 · r10: 1 · r11: 1`.

---

## Fix-up cycle 13

Review round 12: **codex-1 ❌ BLOCK** (1 MAJOR). Two more false greens, both from deciding
**once per line** instead of once per command:

1. **Compound.** ``echo `node --test "…"`; node --test missing.test.js`` — the extractor
   emitted only the span's valid command and discarded the surrounding shell *and* the broken
   second command. Suite stayed **253/0**.
2. **Backslash continuation.** `node --test "…claim.test.js" \` + a second line. The guard ran
   only the first physical line; `/bin/sh` stripped the trailing backslash and the five passing
   claim tests reported green. The continuation was never seen. **253/0**.

### What changed

The discriminator no longer tries to work out *which part* of a mixed line to run:

- a line ending in `\` is a command that does not fit on this line → **the whole line becomes
  the unit**, and the grammar refuses it;
- a line where `node --test` appears **both** inside a span and outside it is compound or
  substituted → **the whole line becomes the unit**, refused;
- only when *nothing outside the spans* runs the command are the spans executed individually.

`SUPPORTED_COMMAND` also now excludes `\`. Without that the continuation line still matched —
the grammar permitted the trailing backslash, `sh` stripped it, and the truncated half ran
green. **P2 passed on the first attempt of this cycle and I caught it only by running the
probe**, which is why the probe set is re-run in full every cycle rather than trusted.

### Proved — all six probes from rounds 10, 11 and 12

```text
compound span+shell           → 252/1     backslash continuation      → 252/1
fake closing fence            → 252/1     blockquoted fence (valid)   → 253/0
fenced substitution           → 252/1     double-backtick (valid)     → 253/0
restored                      → 253/0
```

### Findings per round

`r01: 8 (3 reviewers) · r02–r08: 1 each · r09: 3 · r10: 1 · r11: 1 · r12: 1`.

## Fix-up cycles 14 and 15 — the continuation, and the container the continuation hid in

Round 13 was `❌ BLOCK` with one MAJOR: a continuation *before* `--test` bypassed the guard.
Round 14 was `❌ BLOCK` too — but nobody signed it, because the reviewer never got to finish.

### Cycle 14 — logical lines before detection

`publishedTestCommands` filtered each **physical** line with `/node\s+--test/` before it looked
for a continuation, so it saw a continuation only when one line already held both tokens.
codex-1's probe put the break between them:

```text
node \
  --test skills/parley-worktrees/round13-definitely-missing.test.js
```

Copied out and run: exit 1. Suite: **253/0**. Neither line reached the refusal or the grammar.

Physical lines are now spliced into logical ones — exactly as a shell removes
backslash-newline — **before** anything looks for a command, and a spliced unit is emitted with
its backslash restored so the strict grammar refuses it by name. Cycle 14b added codex-1's
exact probe and both token-boundary splits it named (inside `node`, inside `--test`), because
cycle 14's own fixture had covered neither.

### Round 14 — no signoff exists, and why

`codex exec` was terminated mid-run by an upstream OpenAI content filter
(*"This content was flagged for possible cybersecurity risk"*) while writing its probe harness.
It had spent ~119k tokens, wrote no review file, and left two untracked files behind
(`.round14-probes.js`, `skills/round14-codex-probe.md`) which were removed. **Round 14 has no
reviewer signoff.** This is recorded as a tool outage, not as an accept.

Two things survived it, and both are usable:

1. Its last recorded step: *"testing whether Markdown container syntax can still interrupt the
   reconstructed logical shell line."*
2. Its unrun harness, in which `P13-blockquote` carried `expected: "GREEN"` — codex-1 had
   **predicted the false green before measuring it**.

### Cycle 15 — the prediction was correct

Its probe, run verbatim on a clean tree:

```text
> ```bash
> node \
>   --test skills/parley-worktrees/round14-blockquote-missing.test.js
> ```
```

Guard **12/0 — green**. The same command copied out: **exit 1**.

A blockquote marker sits at the start of *every* line it contains, but the container-noise
strip ran **once, after splicing**, so it only ever cleaned the first line. The interior
markers stayed embedded — `node > --test x` — which matches no command pattern, so the command
was never tested at all.

This is not a new rule. It is the existing rule applied to every physical line instead of one
of them. Because the strip now precedes reconstruction, it also closes a shape raw splicing
could never reach: the tokens themselves split across lines inside a blockquote.

**Named divergence:** `/bin/sh` preserves a continuation line's leading whitespace; this does
not. The guard reconstructs more aggressively than a shell — it refuses more and executes
less. Fail-closed, and stated here rather than left to be discovered.

### Proved — nineteen probes, clean tree, cycle 15 in place

```text
split after --test            → 252/1     round-13 exact probe        → 252/1
split inside `node`           → 252/1     split inside `--test`       → 252/1
valid half + poisoned tail    → 252/1     continuation, valid target  → 252/1
fenced substitution           → 252/1     blockquote continuation     → 252/1
cd … && …                     → 252/1     indented continuation       → 252/1
NODE_OPTIONS=… prefix         → 252/1     list-item continuation      → 252/1
                                          nested blockquote           → 252/1
plain inline (valid)          → 253/0     tokens split in blockquote  → 252/1
double-backtick (valid)       → 253/0     codex-1's exact probe       → 252/1
blockquote single (valid)     → 253/0
genuinely broken path         → 252/1  ← RUNS and fails, not refused
```

The last line is the one that matters most: the guard did not become a refuser that verifies
nothing.

### Findings per round

`r01: 8 (3 reviewers) · r02–r08: 1 each · r09: 3 · r10–r13: 1 each · r14: 1 (predicted by the
reviewer, confirmed by me after the reviewer was cut off)`.

## Fix-up cycle 16 — the marker goes, the whitespace stays

Round 15 ran (codex-1, `❌ BLOCK`) and refuted a claim **I** had made in cycle 15.

I wrote that dropping a continuation line's leading whitespace could only cause extra
refusals — fail-closed. codex-1 measured the opposite:

```text
node\
  --test skills/parley-worktrees/round15-leading-space-missing.test.js
```

A shell removes only the backslash-newline, so a reader gets `node  --test x` — it runs, and
it exits 1. Deleting the indentation produced `node--test x`, which the detector cannot see.
Guard: **12 pass / 0 fail**. The same false-green class as round 14, from the opposite side.

Cycle 15 had conflated two things that must be separated:

- the container **marker** must go — it sits at the head of every line it contains;
- the content's **whitespace** must stay — deleting it merges tokens the shell keeps apart.

Cycle 16 removes exactly the marker: up to three spaces of lead-in, the `>`, and at most one
space of padding, per nesting level — all CommonMark consumes. Nothing else is touched.

**The divergence cycle 15 named as acceptable is gone rather than justified.** The
reconstruction is no longer "more aggressive than a shell"; it is what the reader copies out
of the rendered page.

### Proved — twenty-two probes, clean tree

Every continuation shape refused, with and without a space before the backslash, inside
blockquotes, nested blockquotes, indented blocks and list items, and with the tokens
themselves split across lines. All three valid forms still green. A genuinely broken path
still **runs and fails** rather than being refused.

### Process defect, and the correction

Rounds 03–15 were reviewed by **codex-1 alone**. Thirteen consecutive single-reviewer rounds,
each finding exactly one real problem from one perspective. Round 16 runs the **full roster** —
`codex-1`, `agy-1`, `hermes-1`, `kimi-1` — with `antigravity-1` reactivated for it. Each
reviewer works in its own git worktree with its own probe harness, because probing writes
temporary files into the tree and concurrent reviewers on one checkout corrupt each other's
measurements. I hit exactly that contamination myself earlier in this session and read three
probe results wrong before the error message named a file I had not written.

### Findings per round

`r01: 8 (3 reviewers) · r02–r08: 1 each · r09: 3 · r10–r13: 1 each · r14: 1 (predicted by the
reviewer, confirmed by me after it was cut off) · r15: 1 (refuted a claim of mine)`.

## Round 16 and cycles 17–18 — the roster change paid for itself in one round

Round 16 was the first full-roster review since round 02. **All four reviewers signed
`❌ BLOCK`, and three of the four findings came from the agents that had been absent for
thirteen rounds.**

| reviewer | finding | disposition |
|---|---|---|
| `agy-1` | a node flag between the tokens (`node --no-warnings --test x`, `node -r ./setup.js --test x`) matched no detection pattern, so the command was **skipped entirely** | cycle 17 |
| `codex-1` | markdown **rendering** synthesizes commands no source scanner sees: an escaped backslash, emphasis inside the flag, a numeric entity | cycle 18 |
| `kimi-1` **and** `hermes-1`, independently | cycle 16 stripped a `>` it could not prove was a marker and **executed the mutated text** — a fenced `> node --test x` certified green while the reader's copy is a redirection that creates a file named `node` and exits 127 | cycle 18 |
| `kimi-1` | the zero-width continuation boundary — nothing on either side of the break | cycle 17 |
| `hermes-1` | the same repair-instead-of-refuse defect for the `$ ` prompt | cycle 18 |

`agy-1` found its MAJOR in its **first review after reactivation**, in a class thirteen
single-reviewer rounds had not touched. That is the measurable answer to whether the roster
gap mattered.

### Cycle 17 — detection must be broader than acceptance

`agy-1`'s finding exposed the mistake every false green in seventeen rounds shares: **detection
used the same pattern as acceptance.** Whatever the grammar would refuse, detection also failed
to see — so the command was not refused, it was skipped, and skipping reads as success.

A unit is now a candidate if it mentions `node` and `--test` at all, in any order, with
anything between them. What runs is decided afterwards and only by the grammar.

### Cycle 18 — stop approximating markdown; parse it

All four reviewers independently concluded the class cannot be closed by a scanner over source
text. The user ratified the publication-contract-plus-parser direction.

**The contract:** a verification command MUST be the whole text of a single code node — one
inline span, or one line of one code block — in canonical form. Anything else that renders as
such a command is refused by name.

- A real CommonMark parser (`commonmark`, a new devDependency; CI already runs `npm ci`)
  produces the AST.
- `publishedTestCommands` returns command → **provenance** (`code` | `prose`), and provenance
  is checked **before** form: a canonical command reaching the reader out of prose is refused
  rather than run, so it cannot pass by happening to work.
- Inside a code node a continuation is spliced raw and emitted **with** its backslash, so the
  grammar refuses it. One node, one line.

**What the parser buys that no line heuristic could:** a `>` is a marker or content depending
only on its container. A fence *inside* a blockquote yields the bare command and runs; a `>`
*inside* a fence stays in the literal and is refused. `kimi-1` named that asymmetry precisely
and it is why cycle 16's per-line stripper could not be repaired.

### The fixture was structurally wrong, and the old design could not tell

Rewritten as a **well-formed** document. The previous fixture was authored for a line scanner,
so its structure was accidental: an unclosed fence had been swallowing the blockquote cases
below it, and its assertions still passed **because the scanner did not model containers
either**. The parser exposed it on the first run. A fixture that agrees with a defect is not
evidence.

### Proved — thirty-one probes, clean tree

Every finding from rounds 13–16 turns the guard red; all three valid forms stay green; a
genuinely broken path still **runs and fails** rather than being refused. Suite 253/253.

### Findings per round

`r01: 8 (3 reviewers) · r02–r08: 1 each · r09: 3 · r10–r13: 1 each · r14: 1 (predicted by the
reviewer, confirmed by me after it was cut off) · r15: 1 (refuted a claim of mine) · r16: 5
across 4 reviewers (two more refuted claims of mine)`.
