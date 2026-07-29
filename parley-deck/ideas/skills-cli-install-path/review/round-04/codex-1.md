---
idea: skills-cli-install-path
review-round: 04
agent: codex-1
date: 2026-07-29
---

## Summary

Re-reviewed `parley-deck-skill` branch `readme-skill-catalogue` at
`fa1fdb13e7e66085c7a0a72f1f1b7550443e2b4b`. All six findings I filed in rounds 01–03 are
fixed in the shipped behavior, including the tracker exemplar that caused the round-03
block. The Gemini reconciliation, native destination shapes, packaging, universal install,
stale-path guard, and current published commands all reproduce.

I found one remaining MAJOR in fix-up cycle 4. Its new guard does not find every published
`node --test` command as claimed: it extracts inline code spans only and misses fenced code
blocks. One of the two existing unique published commands is already fenced. Adding a new
broken fenced command in a scratch copy leaves all 252 tests green, so the regression guard
does not cover the class it says it covers.

## Prior finding dispositions

### FIXED — round 01 [MAJOR] D-1 removed a reconcilable Gemini channel

D-1 remains withdrawn, and reconciling the two consumers was the right decision:

- Repository `gemini-extension.json` uses
  `contextFileName: "skills/parley-deck/SKILL.md"`; that path resolves in the repository.
- A staged Gemini install uses `contextFileName: "SKILL.md"`; that path resolves in the
  destination.
- Removing `contextFileName` from both parsed manifests leaves equal objects, so the actual
  rewrite changed no other field.
- Changing the repository value to `SKILL.md` in a scratch copy made the focused test fail
  with exit 1 and the actual/expected values.
- Changing the staged rewrite to `BROKEN.md` in another scratch copy made its focused test
  fail with exit 1 and the actual/expected values.

Both reconciliation tests can fail when their protected values are broken. The README
restores `gemini extensions install <repo-url>` and explicitly says the Gemini CLI has not
been run end to end. That is honest; `gemini` is absent on this host, so G11 remains NOT
TESTED rather than being presented as a pass.

### FIXED — round 01 [MINOR] G1 was recorded as an unconditional pass

Fix-up cycle 1 still qualifies the clean-HOME detection ordering. I reran the exact sequence
in a fresh scratch HOME:

- first `install --target all --force`: exit 0;
- first `status --target all`: exit 0, with five Gemini units missing;
- first `doctor --target all`: exit 1;
- second identical install: exit 0;
- second status and doctor: every one of the 25 detected core/add-on units valid, both exit 0.

The same first-run behavior exists at pre-move parent `94a4889`; it is not a layout
regression. The current implementation record discloses the precondition and the
Antigravity-to-Gemini ordering follow-up.

### FIXED — round 01 [MINOR] The installed README named checkout-only paths

The README distinguishes repository paths under `skills/parley-deck/` from flat installed
paths `SKILL.md` and `references/COOPERATION.md`. Both stated shapes exist.

### FIXED — round 01 [MINOR] Universal-agent wording exceeded measured scope

The panel now attributes the supported-agent surface to the upstream installer, states no
agent count, and does not say “whichever coding agents you have.” The cached upstream CLI
contains more named target configurations than this package's fourteen, while the all-five
claim is backed by the real G7 install below.

### FIXED — round 02 [MAJOR] Shipped instructions referenced the removed `addons/` tree

The exact `grep -rn "addons/" skills/` command returns no matches. There is no `addons/`
directory, compatibility tree, or symlink, and the package contains no `addons/` entry.

The cycle-3 guard is genuine. Reintroducing one stale path in
`skills/parley-tracker/SKILL.md` in a scratch copy made `npm test` report 251 pass and 1 fail;
the assertion named `skills/parley-tracker/SKILL.md:82`. Restoring the path returns the
unmodified suite to 252/252.

The shipped runnable instructions also resolve: the design-check help command exits 0; all
three tracker exemplars pass strict validation; the epic's strict-directory and single-file
commands pass in a constructed `tickets/` tree; and the story claim command writes
`status: in-progress` and `assignee: me` in a scratch copy.

### FIXED — round 03 [MAJOR] The tracker exemplar's test command failed with zero passes

Both exemplar occurrences now publish
`node --test "skills/parley-tracker/bin/*.test.js"`. Running it executes 35 tests, all
passing. Reverting the exemplar to `node --test skills/parley-tracker/bin` in a scratch copy
makes `npm test` report 251 pass and 1 fail, with `MODULE_NOT_FOUND` and the child suite's
zero-pass result. The new guard therefore catches the original regression.

The other published test command,
`node --test "skills/parley-design-check/test/*.test.js"`, independently executes 159 tests,
all passing. The current files are correct; the new finding below concerns the guard's
incomplete discovery of published commands.

## Refutation attempts and gate results

### G7 — install without `--full-depth`

I copied `fa1fdb1` to scratch and ran exactly:

`HOME=<scratch> npx -y skills@latest add <scratch-repo> --agent claude-code --yes --copy`

with no `--full-depth`. The npm wrapper failed before the CLI ran because
`registry.npmjs.org` could not resolve (`ENOTFOUND`), so I cannot reproduce the network-backed
`@latest` resolution in this sandbox.

Using the locally cached `skills` 1.5.20 executable—the same version used by the ratified
design evidence—with otherwise identical arguments and still no `--full-depth` produced
`Found 5 skills` and `Installed 5 skills`. `skills list --json` returned five entries.
The installed core contained `SKILL.md`, `agents/`, and `references/`; `bin/`, `lib/`, and
`package.json` were absent.

### Native destinations, status, doctor, and Antigravity

`validateInstalledPayload` is unchanged from `94a4889`. After the converged scratch install,
I listed the actual core destinations and compared their relative path sets with installs
made from that parent:

- Codex: 13 entries, identical to the parent.
- Antigravity: 15 entries, identical to the parent, including fabricated
  `skills/SKILL.md`.
- Gemini: 13 entries, identical to the parent, with staged
  `contextFileName: "SKILL.md"`.

`agy plugin validate` against the installed Antigravity destination exits 0 with one skill
and two agents processed. I found no per-target installed-destination regression.

### Tests and package

- `npm test`: 252 passed, 0 failed.
- `npm pack --dry-run --json`: 153 files; 145 under `skills/`; exactly five
  `skills/<name>/SKILL.md` roots; no `addons/`; no root `SKILL.md`.

### Previously NOT TESTED gates

- **G3 Windows:** PARTIALLY TESTED. Using the two cached target runtimes explicitly, both
  Windows x64 and ARM64 binaries build from this head; `file` identifies them as PE32+
  x86-64 and Aarch64 executables. Windows execution/install remains NOT TESTED because this
  host has neither Windows nor Wine.
- **G8 `skills update`: PASS.** I installed all five from an isolated `file://` Git remote,
  committed a new core reference file to that remote, and ran `skills update --yes`. It
  reported all five updated, and the new file appeared in the installed core.
- **G9 Homebrew:** NOT TESTED against this head. Homebrew is installed, but the tapped
  formula downloads the published `v1.5.0` tag and this repository has no candidate formula;
  `brew test` or `brew upgrade` would test different source.
- **G10 WinGet:** NOT TESTED. `winget` and a Windows execution environment are unavailable
  on this macOS host; this repository contains only `packaging/winget/README.md`.

### README panel and implementation claims

The panel satisfies A4/F1/F3/F4/F5: it is first under Install, neither path is labelled
recommended, it states no numeric universal-agent count, it includes `skills list`, the
native installer and doctor remain in the same screenful, and the all-five statement is
backed by G7. F4's core-only fallback is not triggered because the move and five-skill
install pass.

The current test, package, destination, Gemini, stale-path, and exemplar claims reproduce.
The exact network-backed `npx ...@latest` wrapper remains unreproducible here because of DNS,
and the initial clean-HOME G1 row requires the later recorded qualification. The cycle-4
claim that the guard runs every published `node --test` command does not reproduce and is
the finding below.

## Findings

### [MAJOR] The published-command guard ignores fenced `node --test` commands

`test/design-addons.test.js:228` extracts only text matching a single-backtick inline span
with ``/`(node --test [^`]+)`/g``.

The tracker command is inline, so the guard runs it. The design-check command at
`skills/parley-design-check/SKILL.md:371-373` is in a fenced code block, so the extractor
does not see it. Enumerating with the test's own regex finds only the tracker command even
though shipped Markdown contains two unique published commands.

I tested both sides in scratch copies:

- Adding a new broken inline command made `npm test` fail 251/1.
- Adding the same broken command in a fenced block made `npm test` remain 252/252.

This is not a theoretical alternate style: one of the two commands already uses the missed
fenced form. A future regression of the 159-test design-check command would therefore return
the same false green that fix-up cycle 4 claims to prevent.

Extend discovery to executable lines in fenced code blocks as well as inline code spans,
normalize and deduplicate the resulting commands, and add a focused extractor fixture that
contains one valid inline command, one valid fenced command, and a broken command in each
form. The suite should prove both forms are discovered and both broken forms turn it red.

### Signoff: codex-1 — 2026-07-29
Status: ❌ BLOCK
