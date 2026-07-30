---
idea: skills-cli-install-path
review-round: 16
agent: agy-1
date: 2026-07-30
reviewed-commit: c45601f
---

## Summary

Reviewed clean branch `readme-skill-catalogue` at commit `c45601f` in isolated worktree `wt-agy`. 

The author's claims regarding cycles 14, 15, and 16 were measured using the provided probe harness (`probe-agy.sh`). All 22 probes in the harness reproduce as claimed (18 refusals, 3 green controls, 1 ran-and-failed control). The full test suite on the clean tree passes at 253/253.

However, the guard remains vulnerable to false-green certifications. Node CLI options positioned between `node` and `--test` (such as `--no-warnings`, `--trace-warnings`, `--enable-source-maps`, `--experimental-strip-types`, or `-r ./module.js`) bypass the `/node\s+--test/` pre-filter in `publishedTestCommands`. When a document publishes such a command and points to a missing or broken test file, a reader copying the command receives `exit 1`, while the guard completely misses the command and reports `12 pass / 0 fail` (GREEN). This is one remaining **MAJOR** finding.

Sixteen consecutive fix-up cycles of hand-written regexes, line splicers, and Markdown container strippers demonstrate that attempting to parse open-ended Markdown and shell syntax is an endless task. A structural rule requiring published verification commands to stand alone as single-line commands of the form `node --test <targets>` would close this class by construction.

---

## What was verified and how

### Harness & Baseline Verification
- Workspace: Isolated git worktree `/private/tmp/claude-501/-Volumes-My-Shared-Files-AI-WORKSPACE-parley-deck/5dc331bd-5ddf-45e0-b6c2-d519d8c05128/scratchpad/wt-agy` at commit `c45601f`.
- Full test suite execution: `npm test` passed with `253 pass / 0 fail`.
- Harness execution: `zsh /private/tmp/claude-501/-Volumes-My-Shared-Files-AI-WORKSPACE-parley-deck/5dc331bd-5ddf-45e0-b6c2-d519d8c05128/scratchpad/probe-agy.sh`.
  - Baseline (no probe file): `12 pass / 0 fail`.
  - Probes P1–P3, P5–P6, P9–P18, P20–P22: `REFUSED` (11 pass / 1 fail).
  - Probes P4, P7, P19: `GREEN` (12 pass / 0 fail).
  - Probe P8 (genuinely broken path): `RAN-AND-FAILED` (11 pass / 1 fail).

### Named Trade-off Evaluation
Cycle 16 explicitly refuses legitimate `node --test` commands that are split across lines using backslash continuations. This is an acceptable fail-closed policy: no shipped file currently uses backslash continuations for test commands, and requiring documentation commands to be written on a single line is a clean authoring boundary.

### Cleanliness
All temporary probe files were removed after testing. `git status` in `wt-agy` was verified clean. No tracked files were modified.

---

## Findings

### [MAJOR] Node CLI options before `--test` bypass detection and produce false greens

`publishedTestCommands` filters lines in `test/design-addons.test.js` at line 273 using:
```js
if (!/node\s+--test/.test(text)) continue;
```
This regex requires the literal word `node`, followed immediately by whitespace, followed immediately by `--test`.

When a published shell command includes Node.js CLI flags or options prior to `--test` (for example, `node --no-warnings --test <targets>`, `node --trace-warnings --test <targets>`, `node --enable-source-maps --test <targets>`, `node --experimental-strip-types --test <targets>`, or `node -r ./setup.js --test <targets>`), the regex `/node\s+--test/` returns `false`. `publishedTestCommands` silently drops the command without adding it to `units` or passing it to `SUPPORTED_COMMAND`.

#### Reproduction & Empirical Measurement

1. In `wt-agy`, create a temporary Markdown file `skills/__probe_node_options__.md` containing:
   ```markdown
   Run: `node --no-warnings --test skills/parley-tracker/bin/missing.test.js`
   ```
2. Run the guard test:
   ```bash
   node --test test/design-addons.test.js
   ```
   **Guard Output:** `ℹ pass 12 / ℹ fail 0` (GREEN). The guard does not see or execute the command.
3. Run the exact published command directly in `/bin/sh`:
   ```bash
   node --no-warnings --test skills/parley-tracker/bin/missing.test.js
   ```
   **Shell Output:** Exits with code **1** (`Could not find 'skills/parley-tracker/bin/missing.test.js'`).

The guard certifies success while the reader's command fails.

#### Recommendation
Update the line pre-filter in `publishedTestCommands` to allow flags between `node` and `--test` (e.g. `/node(?:\s+\S+)*\s+--test/`), or adopt the structural rule outlined below.

---

### Architectural Assessment: Close the class by construction

Sixteen review cycles have each uncovered a single blind spot in a hand-written Markdown/shell extractor:
- Round 03: Ignored broken command forms.
- Round 04: Ignored fenced code blocks.
- Round 05: Ignored tilde fences.
- Round 06: Missed second command on same line.
- Round 07: Missed tab whitespace.
- Round 08: Rewrote punctuation & misread single quotes.
- Round 09: Discarded surrounding shell context.
- Round 10: Repaired fenced substitutions.
- Round 11: Mishandled fake fence closes.
- Round 12: Dropped command text across line splits.
- Round 13: Bypassed continuations before `--test`.
- Round 14: Failed on Markdown blockquote markers.
- Round 15: Merged tokens by stripping indentation.
- Round 16: Missed Node CLI flags before `--test`.

This 16-cycle history demonstrates that attempting to parse open-ended Markdown and reconstruct shell command units using custom regexes is an uncloseable problem space. 

Instead of adding a 17th regex patch, the project should enforce a structural constraint:
> **Rule:** Every published verification command in documentation MUST be written as a standalone, single-line fenced code block of the explicit form `node --test <targets>`.

Enforcing this simple authoring rule simplifies `publishedTestCommands` into a predictable, fail-closed check and eliminates the entire class of Markdown-reconstruction edge cases by construction.

---

### Signoff: agy-1 — 2026-07-30
Status: ❌ BLOCK
