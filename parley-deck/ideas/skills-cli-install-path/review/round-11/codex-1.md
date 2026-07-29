---
idea: skills-cli-install-path
review-round: 11
agent: codex-1
date: 2026-07-29
---

## Summary

Reviewed `readme-skill-catalogue` at `39919bf`, the head recorded for fix-up cycle 11.
The layout, installer, packaging, Gemini reconciliation, native destination shapes, README
panel, and prior regression fixes reproduce. I block on one new defect in the cycle-11
published-command guard: its handwritten fence state does not follow CommonMark fence
boundaries or container prefixes. I defeated it in both directions in scratch copies.

## Prior finding dispositions

| Prior finding | Status | Current verification |
|---|---|---|
| Round 01 MAJOR — D-1 dropped a reconcilable Gemini channel | **FIXED** | Repository `contextFileName` is `skills/parley-deck/SKILL.md` and resolves in the repository; a native Gemini install rewrites it to `SKILL.md` and that resolves in the flat destination. Mutating either value independently made its focused test fail 0/1. |
| Round 01 MINOR — G1 was recorded as an unconditional pass | **FIXED** | In a fresh HOME, first `install --target all --force` returned 0, then status exposed Gemini as missing and doctor returned 1; a second install made status and doctor return 0. Fix-up cycle 1 now records these preconditions and the pre-existing detection-order defect rather than calling the clean-HOME sequence unconditionally green. |
| Round 01 MINOR — installed README documented checkout-only paths | **FIXED** | `README.md:155-158` now distinguishes repository paths from flat installed-destination paths; both shapes exist where stated. |
| Round 01 MINOR — “whichever coding agents you have” exceeded the measured scope | **FIXED** | The panel now says “the coding agents it supports” and explicitly leaves the supported-agent list to the upstream project. |
| Round 02 MAJOR — shipped instructions retained `addons/` paths | **FIXED** | `grep -rn "addons/" skills/` returned nothing and no `addons/` tree exists. Reintroducing one stale path made the full suite 252/1 and named `skills/parley-tracker/SKILL.md:79`. |
| Round 03 MAJOR — tracker exemplar published a zero-pass directory command | **FIXED** | Both real published commands run non-zero suites: design-check 159/159 and tracker 35/35. Reverting the exemplar to `node --test skills/parley-tracker/bin` made the full suite 252/1 and named the command. |
| Round 04 MAJOR — guard ignored fenced commands | **FIXED for the filed probe** | A properly fenced command substitution is now captured whole and refused, making the suite 252/1. |
| Round 05 MAJOR — extractor ignored tilde fences | **FIXED for the filed probe** | The passing extractor fixture now exercises and asserts the tilde-fenced command. |
| Round 06 MAJOR — only the first command on a line was checked | **FIXED** | The fixture asserts both commands from its same-line pair. |
| Round 07 MAJOR — tab whitespace was missed and two targets were passed as one | **FIXED** | The fixture asserts both forms; the two-target command also ran 35/35 manually and remained green when published in a double-backtick span. |
| Round 08 MAJOR — punctuation was silently repaired and single quotes were misread | **FIXED** | The trailing-period probe now makes the suite 252/1 for zero tests; the single-quoted glob runs 35/35 manually and the full suite remains 253/0. |
| Round 09 MAJOR — surrounding shell context was discarded | **FIXED for the filed probes** | The fixture captures the environment-prefixed, `cd ... &&`, and substitution forms whole and asserts that the strict grammar refuses them. |
| Round 09 MINOR — README contradicted itself about detection | **FIXED** | The universal path still states detection; the native comparison now claims only measured health checks and metadata sync as additions. |
| Round 09 NIT — claimed-deleted lexer remained | **FIXED** | The dead lexer and its stale normalisation code are absent. |
| Round 10 MAJOR — fenced substitutions were stripped and double-backtick spans misread | **FIXED for both filed probes** | The fenced substitution makes the suite 252/1 with the named refusal. A legitimate double-backtick two-target command runs 35/35 manually and leaves the suite 253/0. Narrowing span extraction back to one backtick made the extractor fixture fail on `node --test double/span.test.js`, so that fixture is capable of failing. |

## Refutation attempts and gate results

### Universal and native installation

- Exact G7 with `HOME=<scratch> npx -y skills@latest add <scratch-copy> --agent
  claude-code --yes --copy`, with no `--full-depth`, could not resolve
  `registry.npmjs.org` (`ENOTFOUND`). The network-backed `@latest` resolution is the one
  passing claim I could not reproduce on this host.
- The cached `skills` 1.5.20 CLI with the same add arguments found and installed all five.
  The core contained `SKILL.md`, `agents/`, and `references/`; `bin/`, `lib/`, and
  `package.json` were absent.
- Native install/status/doctor had only the clean-HOME detection-order behavior recorded
  above. After the second install, Codex, Antigravity, Gemini, Claude, and Hermes plus all
  their add-ons were valid; doctor returned 0. `agy plugin validate` passed with one skill
  and two agents processed.
- Relative destination path sets at `39919bf` are identical to parent `94a4889` for Codex,
  Antigravity, and Gemini. Codex and Gemini have the marker, README, license, root
  `SKILL.md`, both manifests, `agents/`, and `references/`; Antigravity adds only its
  established fabricated `skills/SKILL.md`. This agrees with unchanged
  `validateInstalledPayload`.

### Tests, package, shipped commands, and deferred gates

- `npm test`: 253 pass, 0 fail.
- `npm pack --dry-run --json`: 153 files, 145 under `skills/`, exactly five
  `skills/<name>/SKILL.md` files, zero `addons/` entries, and no root `SKILL.md`.
- The two published test commands independently run 159/159 and 35/35.
- The two Gemini tests genuinely fail when their respective consumer is broken: changing
  the repository value to `SKILL.md` failed the repository test; changing the staged rewrite
  to `skills/parley-deck/SKILL.md` failed the destination test.
- G3 is partially testable here and passed: current HEAD cross-built PE32+ Windows x64 and
  ARM64 executables. Windows execution/install remains **NOT TESTED** because neither Windows
  nor Wine is available.
- G8 is testable and passed. I installed all five from a scratch `file://` Git remote,
  committed a new core reference to that remote, ran `skills update --yes`, saw all five
  reported updated, and found the new reference in the installed core.
- G9 remains **NOT TESTED** for this branch. Homebrew is installed, but the tap formula points
  at the already-published `v1.5.0` archive; `brew upgrade`/`brew test` would not exercise
  `39919bf`.
- G10 remains **NOT TESTED** because WinGet and Wine are unavailable and this repository has
  no candidate WinGet manifest to validate.
- G11 remains **NOT TESTED** because Gemini CLI is unavailable. The README says this plainly.
- The README panel satisfies F1/F3/F4/F5: universal first, neither path called recommended,
  no agent count asserted, five skills measured, and `skills list` present as verification.
  I found no additional unmeasured first-party claim.

## Findings

### [MAJOR] Fence state still repairs a broken command and rejects a legitimate one

`test/design-addons.test.js:232-236` treats any line beginning with a matching fence marker
as a close. In CommonMark, a closing fence may be followed only by spaces or tabs. A line
beginning with three backticks followed by `not-a-closing-fence` inside a fenced block is
content, not a closer. The
unitizer also does not recognize the `>` container prefix of a fenced block inside a
blockquote.

I measured both failure directions:

1. **False green.** I appended this valid fenced block to a shipped skill:

   ````markdown
   ```bash
   ```not-a-closing-fence
   node --test `printf %s --test-reporter=definitely-missing-reporter` "skills/parley-tracker/bin/*.test.js"
   ```
   ````

   The published command returned 7 manually with `ERR_MODULE_NOT_FOUND`. The unitizer
   incorrectly closed at the second line, treated the shell backticks as a Markdown span,
   dropped the substitution, ran the repaired 35-test remainder, and left the full suite
   green at 253/0.

2. **False failure.** I published the legitimate command in a standard blockquoted fence:

   ````markdown
   > ```bash
   > node --test "skills/parley-tracker/bin/*.test.js"
   > ```
   ````

   The copied command ran 35/35 manually. The full suite fell to 252/1 because the unitizer
   passed the source-line `>` into the strict grammar and refused
   `> node --test "skills/parley-tracker/bin/*.test.js"`.

This is the same material class as the prior guard findings: a broken published verification
command can still ship under a green suite, while valid CommonMark can be blocked. Replace
the handwritten fence/span unitizer with standards-compliant CommonMark parsing in test code
(a development-only dependency would not add a runtime dependency), preserve fenced code
content verbatim, and extract code-span content after container parsing. Add both exact
probes above as regression fixtures and prove the false-green mutation turns red while the
blockquote case stays green.

### Signoff: codex-1 — 2026-07-29
Status: ❌ BLOCK
