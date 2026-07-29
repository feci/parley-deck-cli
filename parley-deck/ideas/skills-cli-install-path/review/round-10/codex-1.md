---
idea: skills-cli-install-path
review-round: 10
agent: codex-1
date: 2026-07-29
---

## Summary

Re-reviewed `parley-deck-skill` branch `readme-skill-catalogue` at
`4a75eb6a13ba630c210193a45fc817bd9ba585c4`. The implementation repository was clean before
and after review; all mutations below were made in isolated scratch exports.

The layout, installer destinations, Gemini reconciliation, package contents, G8 update path,
and all findings through round 08 still reproduce as fixed. Cycle 10 also turns all three
round-09 probes red, fixes the README detection contradiction, and removes the dead lexer.
It does not, however, close the extraction boundary. Backticks used for shell command
substitution inside a fenced command are still treated as Markdown inline delimiters. I
constructed a broken command whose substitution is deleted before validation; the guard runs
the valid remainder and stays green. I also reproduced the opposite error with a standard
double-backtick inline code span. This is one remaining MAJOR finding.

## Prior finding dispositions

### FIXED — round 01 [MAJOR] D-1 removed a reconcilable Gemini channel

D-1 remains correctly withdrawn. Repository `gemini-extension.json` has
`contextFileName: "skills/parley-deck/SKILL.md"`, and that path exists in the repository. A
native Gemini install stages `contextFileName: "SKILL.md"`, and that path exists in the flat
destination.

Both regression tests are genuine. In separate scratch copies:

- changing the repository value to `SKILL.md` made the focused repository-manifest test fail
  with actual `SKILL.md` versus expected `skills/parley-deck/SKILL.md`;
- changing the staged rewrite to `BROKEN.md` made the focused staged-install test fail with
  actual `BROKEN.md` versus expected `SKILL.md`.

The README restores `gemini extensions install <repo-url>` and explicitly says the Gemini CLI
has not been run end to end. That is honest. `gemini` is absent on this host, so G11 remains
NOT TESTED rather than being reported as a pass.

### FIXED — round 01 [MINOR] G1 was recorded as an unconditional pass

The later implementation record qualifies the precondition and records the pre-existing
Antigravity-to-Gemini detection ordering. In a fresh scratch `HOME`, the first exact sequence
returned install 0, status 0, doctor 1: Antigravity created the runtime evidence that made
Gemini detectable only after target selection, leaving five Gemini units missing. A second
identical install returned 0/0/0 with all 25 units valid. This matches the prior parent
measurement and is not a layout regression.

### FIXED — round 01 [MINOR] The installed README named checkout-only paths

The README distinguishes checkout paths under `skills/parley-deck/` from installed paths
`SKILL.md` and `references/COOPERATION.md`. All four paths exist in the stated contexts.

### FIXED — round 01 [MINOR] Universal-agent wording exceeded measured scope

“Whichever coding agents you have” is absent. The support surface is attributed upstream and
no numeric agent count is claimed. As a fresh measurement, cached `skills` 1.5.20 printed
`Installing to all 75 agents`; the native installer returned fourteen targets from
`paths --target all --include-undetected --json`. The panel's “longer list” comparison is
therefore measured.

### FIXED — round 02 [MAJOR] Shipped content referenced the removed `addons/` tree

`grep -rn "addons/" skills/` returned no matches. There is no `addons/` directory,
compatibility copy, or symlink, and the package contains no `addons/` entry. Reintroducing one
stale path in scratch made the focused guard fail and name
`skills/parley-tracker/SKILL.md:82`.

### FIXED — round 03 [MAJOR] The tracker exemplar ran zero passing tests

Both exemplar occurrences publish
`node --test "skills/parley-tracker/bin/*.test.js"`, which runs 35 tests with 35 passes.
Reverting them to `node --test skills/parley-tracker/bin` made the focused guard fail with
`MODULE_NOT_FOUND`, zero child passes, and the exact broken command.

### FIXED — round 04 [MAJOR] The guard ignored fenced commands

The real fenced design-check command runs 159 tests with 159 passes. The real inline tracker
command runs 35 tests with 35 passes. The current guard executes both.

### FIXED — round 05 [MAJOR] The guard ignored tilde-fenced commands

The structure-independent scan still covers the prior tilde-fence probe, and the extractor
fixture still contains that form.

### FIXED — round 06 [MAJOR] The guard checked only the first command on a line

The fixture and current unitizer retain both commands from the prior valid-first,
broken-second same-line probe.

### FIXED — round 07 [MAJOR] The guard missed shell whitespace and combined targets

The fixture retains a tab-separated command and a multi-target command. The legitimate
two-target form remains accepted by the current strict grammar and real-shell runner.

### FIXED — round 08 [MAJOR] The guard repaired punctuation and broke single quotes

Cycle 9 removed normalization and the hand-written argv lexer. The trailing-period probe now
turns the suite red, while the single-quoted glob runs 35 tests and stays green.

### PARTIALLY FIXED — round 09 [MAJOR] The guard executed a substring, not the published command

All three exact round-09 probes now turn the focused suite red:

- the argument-only backtick substitution is rejected;
- the environment-prefixed command is rejected with its prefix intact;
- `cd skills/parley-tracker/bin && node --test "*.test.js"` is rejected with the named
  bare-form refusal rather than a fabricated zero-test result.

But the extraction boundary is still fragmentary when a fenced shell command contains
backticks after an otherwise valid target. The new MAJOR finding below is a direct remaining
instance of this round-09 defect, so this disposition is only PARTIALLY FIXED.

### FIXED — round 09 [MINOR] The install panel contradicted itself about detection

Detection is no longer claimed as exclusive to the native installer. The remaining
distinguishers are `doctor`/`status` health checks and project-metadata sync, all implemented
and measured.

### FIXED — round 09 [NIT] The claimed-deleted lexer remained as dead code

The hand-written argv lexer and its stale comment are gone.

## Refutation attempts and gate evidence

### G7 — real install without `--full-depth`

I exported exact HEAD and ran the requested spelling with an isolated scratch `HOME` and no
`--full-depth`:

`HOME=<scratch> npx -y skills@latest add <path> --agent claude-code --yes --copy`

The wrapper waited for package resolution and then failed before the CLI ran because
`registry.npmjs.org` could not resolve (`ENOTFOUND`). I therefore cannot independently
reproduce the network-backed `@latest` resolution.

Running the cached `skills` 1.5.20 CLI with the same add arguments, source, scratch
project/HOME, and no `--full-depth` printed `Found 5 skills` and `Installed 5 skills`.
`skills list --json` listed exactly the five expected skills. The installed core top level
was exactly `SKILL.md`, `agents`, and `references`; `bin`, `lib`, and `package.json` were
absent. The layout behavior required by G7 passes.

### Native destinations and validation

The extracted `validateInstalledPayload` function has the same SHA-256 at pre-move parent
`94a4889` and current HEAD. Actual current core destinations were:

- Codex: marker, `SKILL.md`, `agents/`, `references/`, both manifests, README, and license;
- Antigravity: the Codex shape plus fabricated `skills/SKILL.md`;
- Gemini: the Codex shape, with staged `contextFileName: "SKILL.md"`.

No per-target installed destination shape changed. `agy plugin validate` against the actual
scratch Antigravity destination exited 0 with one skill and two agents processed.

### Tests, package, and the four previously NOT TESTED gates

- `npm test`: 253 pass, 0 fail.
- `npm pack --dry-run --json`: 153 files; 145 under `skills/`; exactly five
  `skills/<name>/SKILL.md` roots; no `addons/`; no root `SKILL.md`.
- **G3 Windows:** PARTIALLY TESTED. With cached `pkg` runtimes, both x64 and ARM64 executables
  built successfully; `file` identifies PE32+ x86-64 and AArch64 binaries. Windows execution
  and installation remain NOT TESTED because neither Windows nor Wine is present.
- **G8 `skills update`: PASS.** I installed all five from a scratch `file://` Git remote,
  advanced that remote with a new core reference file, and ran `skills update --yes`. It
  reported all five updated, and the new file appeared in the installed core.
- **G9 Homebrew:** NOT TESTED against this head. Homebrew and the formula are installed, but
  the stable formula downloads tag `v1.5.0` at commit `94127c2`, before this implementation;
  `brew test` or `brew upgrade` would test different material. The implementation repository
  contains no candidate formula.
- **G10 WinGet:** NOT TESTED. `winget` and a Windows host are absent, and the repository has
  only older published manifests, not a candidate manifest for this head.
- **G11 Gemini URL install:** NOT TESTED because the Gemini executable is absent. The two
  static resolutions and their independent negative tests pass, but they do not replace the
  end-to-end CLI gate.

The README panel satisfies A4/F1/F3/F4/F5: neither route is labelled recommended; no numeric
agent count is claimed; the longer-list comparison is measured; `skills list` is present and
lists all five after G7; the native installer remains in the same screenful; and F4 is not
triggered because the five-skill move passes. The Gemini URL caveat is explicit. I found no
other unmeasured panel claim.

The exact network-backed G7 wrapper is the only passing implementation claim I could not
rerun, for environmental DNS reasons. The initial unconditional G1 table is superseded by
the cycle-1 correction. All other relevant pass claims reproduce.

## Findings

### [MAJOR] Fenced command substitutions are still stripped before strict validation

`test/design-addons.test.js:228-246` decides that every backtick pair on a line is an inline
Markdown span, even when the line is a shell command inside a fenced block. It removes those
spans before applying `SUPPORTED_COMMAND`. That makes the strict grammar inspect a repaired
remainder, not the whole published command.

I inserted this command into a fenced `bash` block in a shipped skill:

```bash
node --test `printf %s --test-reporter=definitely-missing-reporter` "skills/parley-tracker/bin/*.test.js"
```

Copied verbatim through `/bin/sh -c`, it exited 7 with `ERR_MODULE_NOT_FOUND` because the
substitution selects a nonexistent test reporter. The focused published-command guard exited
0. It deleted the substitution as though it were an inline code span, executed only
`node --test "skills/parley-tracker/bin/*.test.js"`, and observed 35 passing tests. This is a
false green for the exact command-substitution class cycle 10 claims to refuse.

The same unitizer also fails in the opposite direction. I published the valid command in a
standard CommonMark double-backtick inline span:

```markdown
Run ``node --test "skills/parley-tracker/bin/*.test.js"``.
```

The rendered command runs 35 tests with 35 passes, but the guard reconstructed
`Run  node --test "skills/parley-tracker/bin/*.test.js" .` and rejected it. Thus the current
unitizer both omits meaningful shell text and manufactures surrounding text.

Fix the command-unit boundary rather than extending the shell grammar. A fenced command line
must remain whole before strict validation, and inline code must honor equal-length CommonMark
delimiter runs rather than the single-backtick regex. Add both probes above as regressions:
the substitution command must turn the suite red without executing a repaired fragment, and
the double-backtick span must execute 35 tests and stay green.

### Signoff: codex-1 — 2026-07-29
Status: ❌ BLOCK
