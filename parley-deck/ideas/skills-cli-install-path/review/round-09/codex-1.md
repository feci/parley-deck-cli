---
idea: skills-cli-install-path
review-round: 09
agent: codex-1
date: 2026-07-29
---

## Summary

Re-reviewed `parley-deck-skill` branch `readme-skill-catalogue` at
`279a72c15cc4c76f384d80db60fe5f24d6ead75d`. The implementation repository was clean before
and after review; every mutation below was made in an isolated scratch export.

The layout, installer, destination shapes, Gemini reconciliation, package, G8 update path,
and every concrete finding I filed in rounds 01–08 now reproduce as fixed. Cycle 9 also fixes
both round-08 probes. I nevertheless defeated the replacement guard in both requested
directions: a broken published shell command is silently omitted, while a legitimate
published shell command is stripped of its required prefix and rejected. The implementation
therefore still does not enforce “every published `node --test` command,” and its fail-closed
claim is false.

## Prior finding dispositions

### FIXED — round 01 [MAJOR] D-1 removed a reconcilable Gemini channel

D-1 remains withdrawn, which is the right decision. Repository
`gemini-extension.json.contextFileName` is `skills/parley-deck/SKILL.md` and resolves in the
repository. A native Gemini install stages `contextFileName: "SKILL.md"`, which resolves in
the flat destination. The two focused tests are genuine: changing the repository value to
`SKILL.md` failed its test with the wrong actual value, and changing the staged rewrite to
`BROKEN.md` independently failed its test. The README restores
`gemini extensions install <repo-url>` and says plainly that the Gemini CLI has not been run
end to end. `gemini` is absent here, so G11 remains honestly NOT TESTED.

### FIXED — round 01 [MINOR] G1 was recorded as an unconditional pass

The later fix-up text records the precondition and the pre-existing detection-order defect.
In a fresh scratch `HOME`, the first exact sequence returned install 0, status 0, doctor 1:
the Antigravity install made Gemini detectable only after target selection, leaving five
Gemini units missing. A second exact install returned 0/0/0 with all 25 detected units valid.
The same first-pass behavior occurs at parent `94a4889`; this is not a layout regression.

### FIXED — round 01 [MINOR] The installed README named checkout-only paths

The README now distinguishes checkout paths under `skills/parley-deck/` from installed paths
`SKILL.md` and `references/COOPERATION.md`. All four paths resolve in their stated contexts.

### FIXED — round 01 [MINOR] Universal-agent wording exceeded measured scope

“Whichever coding agents you have” is gone. The support surface is attributed upstream and no
numeric agent count is claimed. As a fresh measurement, cached `skills` 1.5.20 printed
`Installing to all 75 agents`; the native installer returned fourteen targets from
`paths --target all --include-undetected --json`. The “longer list” comparison is measured.

### FIXED — round 02 [MAJOR] Shipped content referenced the removed `addons/` tree

`grep -rn "addons/" skills/` returned no matches. There is no `addons/` directory,
compatibility copy, or symlink. Reintroducing one stale path in scratch made the focused
guard and the full suite fail at 252/1, naming `skills/parley-tracker/SKILL.md:79`.

### FIXED — round 03 [MAJOR] The tracker exemplar ran zero passing tests

Both exemplar occurrences now publish
`node --test "skills/parley-tracker/bin/*.test.js"`, which runs 35 tests, all passing.
Reverting them to `node --test skills/parley-tracker/bin` made the guard fail with
`MODULE_NOT_FOUND`, zero child passes, and the exact broken command.

### FIXED — round 04 [MAJOR] The guard ignored fenced commands

A missing fenced command and a missing inline command independently made the current guard
fail. The two real published commands also pass directly: design-check runs 159 tests and
tracker runs 35.

### FIXED — round 05 [MAJOR] The guard ignored tilde-fenced commands

A missing command inserted in a `~~~bash` block made the current guard fail and named the
missing target. Narrowing the extractor back to literal post-`--test` space also made its
fixture test fail, confirming the fixture remains capable of exposing a narrower extractor.

### FIXED — round 06 [MAJOR] The guard checked only the first command on a line

The exact valid-first/broken-second same-line shape made the current guard fail on the second
command.

### FIXED — round 07 [MAJOR] The guard missed shell whitespace and combined targets

A tab-separated missing target made the guard fail, while a legitimate two-target command
passed.

### FIXED — round 08 [MAJOR] The guard repaired punctuation and broke single quotes

The trailing-period command now fails because it runs zero tests. The single-quoted glob
passes and runs 35 tests. A semicolon-bearing suffix is refused before execution with the
named unsupported-syntax reason.

## Gate and compatibility evidence

### G7 — real install without `--full-depth`

I exported the exact review commit to scratch and ran the requested
`HOME=<scratch> npx -y skills@latest add <path> --agent claude-code --yes --copy` spelling
without `--full-depth`. The wrapper failed before the CLI ran because
`registry.npmjs.org` could not resolve (`ENOTFOUND`), so I cannot independently reproduce the
network-backed `@latest` resolution.

Running the locally cached `skills` 1.5.20 executable with the same arguments, same scratch
source, and still no `--full-depth` printed `Found 5 skills` and `Installed 5 skills`.
Exactly five installed `SKILL.md` roots exist. The core top level is exactly `SKILL.md`,
`agents`, and `references`; `bin/`, `lib/`, and `package.json` are absent. `skills list`
lists all five. The layout behavior required by G7 passes.

### Native destinations, status, doctor, and manifests

The extracted `validateInstalledPayload` function has the same SHA-256 at `94a4889` and
`279a72c`. I installed both commits into separate scratch homes and diffed relative path
sets:

- Codex: 13 entries at each commit, empty diff;
- Antigravity: 15 entries at each commit, empty diff;
- Gemini: 13 entries at each commit, empty diff.

Current Codex has the flat core payload, both manifests, README, license, and marker.
Antigravity adds the fabricated `skills/SKILL.md`. Gemini has the flat shape and staged
`contextFileName: "SKILL.md"`. No per-target installed destination shape changed.
`agy plugin validate` passed with one skill and two agents processed.

### Tests, package, and platform/channel gates

- `npm test`: 253 pass, 0 fail.
- `npm pack --dry-run --json`: 153 files; 145 under `skills/`; exactly five
  `skills/<name>/SKILL.md` roots; no `addons/`; no root `SKILL.md`.
- G3 Windows: PARTIALLY TESTED. Using the cached x64 and ARM64 pkg runtimes against a scratch
  export built both executables; `file` identifies PE32+ x86-64 and AArch64 binaries.
  Windows execution/install remains NOT TESTED because neither Windows nor Wine is present.
- G8 `skills update`: PASS. I installed all five from a scratch `file://` Git remote,
  advanced the remote with a new core reference, and ran `skills update --yes`. It reported
  all five updated and installed the new reference.
- G9 Homebrew: NOT TESTED against this head. Homebrew is present, but its stable formula
  downloads the fixed published `v1.5.0` tag tarball; this repository contains no candidate
  formula. `brew test`/`brew upgrade` would test different source.
- G10 WinGet: NOT TESTED. `winget` and Windows are absent, and this repository contains only
  packaging instructions, not a candidate manifest to validate.

The implementation's passing layout, manifest, native-install, package, test, G5, and
cycle-9 round-08-probe claims reproduce. The exact network-backed G7 wrapper is the one pass
I cannot rerun for environmental reasons. The initial unconditional G1 table is true only
under the later recorded precondition and is superseded by fix-up cycle 1.

## Findings

### [MAJOR] The cycle-9 guard still executes a substring, not the published shell command

`test/design-addons.test.js:234-236` starts capture at `node`, discards everything before it,
and stops at the first backtick or newline. The unsupported-shell check at line 244 therefore
sees only that truncated substring. Passing the substring to `/bin/sh -c` at line 343 does not
restore text the extractor already threw away.

I measured three failures:

1. In a fenced block I published
   ``node --test `printf definitely-missing-backtick.test.js` ``. Typed exactly through
   `/bin/sh -c`, it exited 1 with `Could not find 'definitely-missing-backtick.test.js'`.
   The extractor stopped at the first substitution backtick, discarded the now-argumentless
   fragment, and the full suite stayed green at 253/0.
2. I published
   `NODE_OPTIONS='--require ./definitely-missing-preload.cjs' node --test "skills/parley-tracker/bin/*.test.js"`.
   Typed exactly, it exited 1 because the preload is missing. The guard discarded the
   assignment, executed the already-valid suffix, and the full suite stayed green at 253/0.
3. In the other direction,
   `cd skills/parley-tracker/bin && node --test "*.test.js"` legitimately ran 35 tests with
   35 passes. The guard discarded `cd ... &&`, ran the suffix from the repository root, and
   failed the full suite at 252/1 as a zero-test command. It did not give the promised
   refusal for `&`, because the ampersands were outside the captured substring.

This is both a false green and a manufactured failure, and the backtick case is an
unrecognised form silently skipped rather than failed closed. Fix the extraction boundary,
not the shell lexer: validate the entire published shell form before extracting/executing
anything, or adopt an explicitly narrow command grammar/registry and fail with a named reason
whenever surrounding prefix, substitution, continuation, or suffix text is outside it. Add
all three probes as regressions.

### [MINOR] The install panel contradicts itself about agent detection

`README.md:118` says the universal installer detects agents. Three lines later,
`README.md:121-122` says this package's installer adds “detection … that the universal one
does not.” The G7 output itself prints `Agent detected`, so detection is not an exclusive
differentiator. This weakens the otherwise honest F1 comparison. Drop `detection` from the
exclusive list or name the specific behavior that differs; health checks and project-metadata
sync remain valid measured differentiators.

### [NIT] The claimed-deleted hand-written lexer remains as dead code

Fix-up cycle 9 says the argv lexer was deleted, but its unused `argv` function and stale
round-07 comment remain at `test/design-addons.test.js:314-324`. Delete them so the code and
implementation record agree.

### Signoff: codex-1 — 2026-07-29
Status: ❌ BLOCK
