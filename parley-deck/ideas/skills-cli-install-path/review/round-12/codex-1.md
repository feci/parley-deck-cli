---
idea: skills-cli-install-path
review-round: 12
agent: codex-1
date: 2026-07-29
---

## Summary

Reviewed clean branch `readme-skill-catalogue` at `c3aa392`. The layout, packaging, native
destination shapes, Gemini reconciliation, README corrections, and every previously filed
probe reproduce. I block on one new MAJOR defect in the cycle-12 published-command guard:
the span-or-line discriminator still discards executable parts of a published command. I
produced two independent broken commands that return non-zero when copied verbatim while the
full suite remains green at 253/253.

## Prior finding dispositions

| Filed finding | Status | Current verification |
|---|---|---|
| Round 01 MAJOR — D-1 dropped a reconcilable Gemini channel | **FIXED** | Repository `contextFileName` is `skills/parley-deck/SKILL.md` and resolves; native Gemini staging rewrites it to `SKILL.md`, which resolves in the flat destination. Independently mutating either consumer made its focused test fail 0/1. |
| Round 01 MINOR — G1 was recorded as an unconditional pass | **FIXED** | Fix-up cycle 1 now states the clean-HOME precondition and ordering defect. The first install/status/doctor sequence returned 0/0/1 at both current HEAD and parent `94a4889`; explicitly installing newly detected Gemini made status and doctor return 0. This is not a layout regression. |
| Round 01 MINOR — installed README documented checkout-only paths | **FIXED** | README now distinguishes `skills/parley-deck/...` in a checkout from root `SKILL.md` and `references/...` in a native destination; both measured shapes exist. |
| Round 01 MINOR — “whichever coding agents you have” exceeded measured scope | **FIXED** | The panel now says “the coding agents it supports” and leaves the supported-agent list to the upstream project. |
| Round 02 MAJOR — shipped instructions retained `addons/` paths | **FIXED** | `grep -rn "addons/" skills/` returned nothing, and no `addons/` tree exists. Reintroducing one stale path made the full suite 252/1 and named `skills/parley-tracker/SKILL.md:79`. |
| Round 03 MAJOR — tracker exemplar published a zero-pass directory command | **FIXED** | The two real commands run 159/159 and 35/35. Reverting the exemplar to `node --test skills/parley-tracker/bin` made the full suite 252/1 and named the failing command. |
| Round 04 MAJOR — guard ignored fenced commands | **FIXED for the filed probe** | A newly added broken fenced command makes the guard fail 0/1. |
| Round 05 MAJOR — extractor ignored tilde fences | **FIXED for the filed probe** | The extractor fixture covers the tilde form; narrowing detection back to literal post-`--test` space made the fixture fail on its tab case. |
| Round 06 MAJOR — only the first command on a line was checked | **FIXED for the filed two-span probe** | The fixture still asserts both commands in its same-line pair. |
| Round 07 MAJOR — tab whitespace was missed and two targets became one argument | **FIXED** | The fixture asserts both forms; the real multi-target command remains executable. |
| Round 08 MAJOR — punctuation was repaired and single quotes were misread | **FIXED** | The trailing-period and single-quoted-glob behavior remains guarded as recorded. |
| Round 09 MAJOR — surrounding shell context was discarded | **FIXED for the filed probes** | Environment-prefix, `cd ... &&`, and substitution units are captured whole and refused. |
| Round 09 MINOR — README contradicted itself about detection | **FIXED** | The native comparison now claims only health checks and metadata sync as additions. |
| Round 09 NIT — claimed-deleted lexer remained | **FIXED** | The lexer and stale normalization code remain absent. |
| Round 10 MAJOR — fenced substitution was stripped and double-backtick span misread | **FIXED for both filed probes** | Fenced substitution is refused; a valid double-backtick command passes. |
| Round 11 MAJOR — handwritten fence state mishandled fake closes and blockquotes | **FIXED for both filed probes** | Fake-close substitution and ordinary fenced substitution fail by named refusal; the blockquoted fence and valid double-backtick span pass. |

## Refutation attempts and gate results

### Universal and native installation

- The literal G7 wrapper
  `HOME=<scratch> npx -y skills@latest add <scratch-copy> --agent claude-code --yes --copy`
  could not resolve `registry.npmjs.org` (`ENOTFOUND`). This is the one passing wrapper claim
  I could not reproduce in this sandbox.
- The locally cached `skills` 1.5.20 CLI—the version used by the design evidence—with the
  same add arguments found and installed five skills without `--full-depth`. The installed
  core contained `SKILL.md`, `agents/`, and `references/`; it excluded `bin/`, `lib/`, and
  `package.json`. `skills list` then listed all five.
- The required native install/status/doctor sequence has only the documented clean-HOME
  ordering behavior: first-pass return codes were 0/0/1 because Antigravity creates the
  evidence that makes Gemini detectable after target resolution. Parent `94a4889` produced
  the same return codes. After explicit Gemini installation, every detected core/add-on unit
  was valid and doctor returned 0. `agy plugin validate` passed with one skill and two agents.
- `validateInstalledPayload` is textually unchanged from `94a4889`, and relative path-set
  diffs for actual Codex, Antigravity, and Gemini core destinations were empty. Codex and
  Gemini have the marker, README, license, root `SKILL.md`, both manifests, `agents/`, and
  `references/`; Antigravity adds its established fabricated `skills/SKILL.md`.
- Both Gemini values resolve in their own consumer. The repository manifest retained its
  other keys, and the staged manifest retained the same four-key shape while changing only
  `contextFileName`. The two mutation tests prove the regression checks can fail.

### Tests, package, shipped instructions, and deferred gates

- `npm test`: 253 pass, 0 fail.
- `npm pack --dry-run --json`: 153 files, 145 under `skills/`, exactly five
  `skills/<name>/SKILL.md` files, zero `addons/` entries, and no root `SKILL.md`.
- The tracker strict validator reports all three shipped templates valid. The documented
  design-check CLI resolves and prints its help. The only published commands claiming to
  run tests execute non-zero suites: 159/159 and 35/35.
- The stale-path guard, old exemplar guard, new inline/fenced-command checks, extractor
  fixture, and both Gemini reconciliation tests all turned red under their respective
  scratch mutations.
- G3 is partly testable here and passes that part: current HEAD cross-built PE32+ Windows x64
  and ARM64 executables. Windows execution/install remains **NOT TESTED** because neither
  Windows nor Wine is available.
- G8 is testable and **PASS**. I installed all five from a scratch `file://` Git remote,
  committed a new core reference to that remote, ran `skills update --yes`, saw all five
  reported updated, and found the new reference in the installed core.
- G9 remains **NOT TESTED** for this branch. Homebrew is installed, but this repository's
  formula directory is empty; testing the external formula would exercise an already
  published release, not `c3aa392`.
- G10 remains **NOT TESTED**: WinGet and Wine are unavailable, and the repository has no
  candidate WinGet manifest beyond its packaging README.
- G11 remains **NOT TESTED** because Gemini CLI is unavailable. README states that limitation
  honestly while keeping both supported manifest shapes.
- The README panel satisfies A4/F1/F3/F4/F5: universal first, neither path labelled
  recommended, no agent count asserted for the universal tool, all five skills measured,
  and `skills list` present and verified. I found no additional unmeasured first-party claim.

## Findings

### [MAJOR] The span-or-line discriminator still drops executable command text

`test/design-addons.test.js:230-245` processes physical lines independently. If any backtick
span on a line contains `node --test`, line 244 emits only those spans and discards the rest
of the line. It also never joins a shell continuation whose following line does not repeat
`node --test`. Consequently, the strict grammar at line 256 never sees the discarded text.

I measured two false greens in separate scratch copies:

1. A fenced compound command:

   ```text
   echo `node --test "skills/parley-tracker/bin/*.test.js"`; node --test definitely-missing-outside-span.test.js
   ```

   The copied command returned 1. The extractor emitted only the valid command inside the
   backticks, ignored both the surrounding shell and the broken second command, and left the
   full suite at 253 pass / 0 fail.

2. A backslash-continued command:

   ```text
   node --test "skills/parley-tracker/bin/claim.test.js" \
     skills/parley-worktrees/round12-failing-probe.js
   ```

   The second target was a deliberate failing test, so the copied two-line command returned
   1. The guard executed only the first physical line; `/bin/sh` removed its terminal
   backslash and ran the five passing claim tests. The continuation line was ignored, and the
   full suite again stayed 253/0.

This disproves the cycle-12 claim that “does a span contain the whole command?” is a complete
discriminator. The decision is still made once per line, not once per command occurrence,
and physical lines are not shell command units.

Do not add two more Markdown/shell cases to this heuristic. Replace it with a fail-closed
registry of the exact published test commands and their source occurrences: every raw
`node --test` occurrence must map bijectively to one registered, single-line whole command
with no executable prefix/suffix or continuation, every unregistered or ambiguous occurrence
must fail, and the registered commands must still be executed with non-zero pass counts.
That keeps the current two-command surface explicit and makes adding any future command a
conscious test update. Add both probes above as regression fixtures and prove they turn the
suite red.

### Signoff: codex-1 — 2026-07-29
Status: ❌ BLOCK
