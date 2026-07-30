---
idea: skills-cli-install-path
review-round: 15
agent: codex-1
date: 2026-07-29
reviewed-commit: 070cf51
---

## Summary

Reviewed clean branch `readme-skill-catalogue` at `070cf51`. The nineteen supplied probes
reproduce exactly: fifteen unsupported or continued forms are refused, three legitimate
single-line forms stay green, and the deliberately broken path is actually run and fails.
The last result confirms that the guard did not become a universal refuser. `npm test` is
green at 253/253 on the unmodified tree.

Cycle 15 fixes the exact round-14 blockquote probe, its nested and token-split variants, and
the exact round-13 continuation probes. It does not make the logical-line reconstruction
sound. Removing every continuation line's leading whitespace before splicing can erase the
only separator between `node` and `--test`; the guard then sees `node--test`, reports 12/12,
and never runs the command that `/bin/sh` executes and rejects. This contradicts the stated
trade-off that the divergence can only make the guard refuse more.

## Prior finding dispositions

| Filed finding | Status | Current verification |
|---|---|---|
| Round 01 MAJOR — D-1 dropped a reconcilable Gemini channel | **FIXED** | Only `test/design-addons.test.js` changed since round 13. The repository manifest still names `skills/parley-deck/SKILL.md`; the staged-install rewrite test and the repository-path test both pass in the 253-test suite. |
| Round 01 MINOR — G1 was recorded as an unconditional pass | **FIXED** | The qualified fix-up record remains intact. No cycle-14/15 commit touched installer detection or the record. |
| Round 01 MINOR — installed README documented checkout-only paths | **FIXED** | The documented repository/native destination distinction is unchanged; the package still contains the five expected skill roots. |
| Round 01 MINOR — “whichever coding agents you have” exceeded the measured scope | **FIXED** | The panel still says “the coding agents it supports,” attributes the installer upstream, and states no agent count. |
| Round 02 MAJOR — shipped instructions retained `addons/` paths | **FIXED** | No `addons/` directory, root `SKILL.md`, or `addons/` reference exists. Adding one temporary shipped reference made the guard fail and name `skills/__round15_probe__.md:1`. |
| Round 03 MAJOR — tracker exemplar published a zero-pass directory command | **FIXED** | The two distinct shipped commands independently pass 159/159 and 35/35. Reintroducing `node --test skills/parley-tracker/bin` in a temporary shipped file turned the guard red. |
| Round 04 MAJOR — guard ignored fenced commands | **FIXED** | The extractor fixture still captures fenced and inline forms; the deliberately broken inline path is run and fails. |
| Round 05 MAJOR — extractor ignored tilde fences and its fixture could not prove coverage | **FIXED** | The executed extractor fixture still asserts the tilde form and the tab-separated form. |
| Round 06 MAJOR — only the first command on one line was checked | **FIXED** | The fixture still asserts both commands from its same-line pair; unsupported compound forms remain whole and refused. |
| Round 07 MAJOR — tab whitespace was missed and two targets became one argument | **FIXED** | The tab and multi-target fixture assertions pass; actual command execution still goes through `/bin/sh`. |
| Round 08 MAJOR — punctuation was silently repaired and single quotes were misread | **FIXED** | The strict whole-command grammar is unchanged, and the published-command test still runs the captured shell spelling rather than normalising argv. |
| Round 09 MAJOR — surrounding shell context was discarded | **FIXED** | The supplied environment-prefix and `cd ... &&` probes are both refused; the fenced substitution is also refused as a whole unit. |
| Round 09 MINOR/NIT — README contradiction and claimed-deleted lexer | **FIXED** | The README wording remains qualified and the handwritten argv lexer remains absent; cycles 14/15 changed only the guard test. |
| Round 10 MAJOR — fenced substitution was repaired and double-backtick spans were misread | **FIXED** | The fenced substitution is refused, while the legitimate double-backtick command is green and runs tests. |
| Round 11 MAJOR — handwritten fence state mishandled fake closes and blockquotes | **FIXED** | The single-line blockquoted command is green; the current extractor no longer carries fence state. |
| Round 12 MAJOR — span-or-line discriminator dropped executable command text | **FIXED** | The compound/context forms and continuation-after-target probe remain red by refusal; legitimate inline and blockquote forms remain green. |
| Round 13 MAJOR — a continuation before `--test` bypassed physical-line detection | **FIXED for the exact probe; class still open** | The exact split, the inside-`node` split, and the inside-`--test` split are refused. The finding below is a different boundary created by cycle 15's whitespace deletion. |
| Round 14 predicted finding — blockquote container markers interrupted splicing | **FIXED** | The exact blockquote probe, nested blockquote, list/indent variants, and token split inside a blockquote are all refused. Round 14 still has no signoff; this disposition is based on the credited probe and cycle-15 measurement. |

## Verification results

### Cycle-15 claims and guard behavior

- `git diff 82507b5..070cf51 --name-only` names only
  `test/design-addons.test.js`.
- The supplied harness baseline is 12 pass / 0 fail. P1–P3, P5–P6, P9–P18 are
  **REFUSED**; P4, P7 and P19 are **GREEN**; P8 is **RAN-AND-FAILED**. This supports all
  nineteen claimed classifications, including the important non-refusal control.
- `npm test`: 253 pass / 0 fail.
- The two distinct commands currently published under `skills/` run 159/159 and 35/35.
- `npm pack --dry-run --json`: 153 files, 145 under `skills/`, exactly five
  `skills/<name>/SKILL.md` roots, no `addons/`, and no root `SKILL.md`.

### Trade-offs

- Refusing a legitimate backslash-continued command is an acceptable fail-closed policy for
  this guard. The supplied valid-target continuation is refused, no shipped file uses the
  form, and flattening such a documentation command is a reasonable authoring requirement.
- Dropping a continuation line's leading whitespace is **not** merely a conservative version
  of that policy. The measurement below shows that it can erase a required token separator
  and make detection weaker than `/bin/sh`, producing a false green. I do not accept this
  trade-off as stated.

### Cleanliness

All temporary Markdown probes were removed. `git status --short` in `parley-deck-skill` was
empty after the probes and before this review was written. No tracked implementation file was
modified.

## Findings

### [MAJOR] Stripping continuation indentation can merge `node` and `--test` and skip the command

`logicalLines` strips `^[\s>]*` from every physical line at
`test/design-addons.test.js:245`, then removes the backslash and joins the parts without a
separator at lines 246–249. That is not always a conservative difference from the shell.

I added this fenced command to one temporary shipped Markdown file:

```bash
node\
  --test skills/parley-worktrees/round15-leading-space-missing.test.js
```

`node --test test/design-addons.test.js` reported **12 pass / 0 fail**. Running the copied
command through `/bin/sh` exited **1** with:

```text
Could not find 'skills/parley-worktrees/round15-leading-space-missing.test.js'
```

The shell removes only the backslash-newline and preserves the two leading spaces, yielding
`node  --test ...`. The guard removes those spaces before splicing, yielding
`node--test ...`; `/node\s+--test/` therefore never detects it. This directly contradicts
the cycle-15 claim that whitespace removal only causes extra refusals.

Preserve shell-content whitespace while removing Markdown container markers. Leading
indentation does not need to be deleted for the unanchored detector, and blockquote markers
should be removed separately from any content indentation. If a container boundary is
ambiguous, refuse the spliced unit rather than normalising it into an undetectable spelling.
Add the exact no-space-before-backslash probe above as a regression.

### Signoff: codex-1 — 2026-07-29
Status: ❌ BLOCK
