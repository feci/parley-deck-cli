---
idea: integrate-parley-bidding-addon
review-round: 8
agent: hermes-1
date: 2026-07-30
reviewed-commit: ebe269e
---

## Verdict

ACCEPT

## Outstanding findings — closed or not

### codex-1 round-1 MAJOR — preflight missed predictable copy failures, leaving a partial fleet — CLOSED

This was the most important finding never dispositioned before cycle 8. `preflightSkillUnit`
walked the source tree only for add-ons carrying a manifest, so a manifest-free add-on
containing a symlink passed preflight and failed inside the sequential write loop — after the
core and every preceding add-on had already been replaced. That is exactly the partial-fleet
state B5 forbids.

**Fixed and verified.** `firstCopyObstacle` (lib/installer.js:1020) is a read-only recursive
walk that mirrors what `copyRecursive` (lib/installer.js:1364) refuses or fails on: symlinks,
non-regular files, unreadable entries. It is called for every unit during preflight
(lib/installer.js:990), not only manifested add-ons. `copySourcesFor` (lib/installer.js:1003)
enumerates an add-on's single directory or the core's `PAYLOAD_ENTRIES` plus
`OPTIONAL_PAYLOAD_ENTRIES`, so the core source is walked too — the gap codex-1 caught in
cycle 9, when the first version of this check ran only when `unit.addon` was truthy.

I independently reproduced both regression scenarios:

- **Manifest-free add-on with a symlink, sorting last.** Staged a package with `zz-broken`
  containing `SKILL.md` plus a symlink, sorting after all six real skills.
  `installCommand` returned `ok: false`; the skills destination did not exist (zero writes);
  the failure message was "Refusing to copy symlink in skill payload: link.md".
- **Symlink in the core source.** Placed a symlink in
  `skills/parley-deck/references/linked.md`. Same result: `ok: false`, zero writes,
  "Refusing to copy symlink in skill payload: linked.md".
- **Clean package (no false positive).** A staging with no defects installed all six units
  normally: `ok: true`, six directories present.

I also verified the preflight is genuinely read-only: a SHA-256 inventory of the staged
source tree before and after a preflight-triggered failure was identical.

The two committed regressions ("a symlink in a manifest-free add-on is caught before the first
write" and "a symlink in the CORE source is caught before the first write too") assert exactly
the zero-writes property I measured. Both pass.

### codex-1 round-1 MINOR (D-2) — evidence did not establish the source was untouched — CLOSED

D-2 called a file count plus a cache scan a "byte-level check" proving the source "untouched".
Neither observes any file's bytes, so a content edit preserving all 48 paths would pass it.

**Fixed.** D-2 in IMPLEMENTATION.md (line 128) is retitled "PRE-7's proving artefact replaced,
and the replacement claimed more than it showed" and narrowed to what the checks establish:
48 files, zero caches, every path and content difference vs the integrated copy accounted for.
The before/after SHA inventory that would have proven byte-level integrity was never captured
and cannot be reconstructed; the record now says so plainly. The verification table in
IMPLEMENTATION.md (line 195) and the consensus verification table (line 159) both state the
narrowed claim — the consensus row explicitly says "Not proof it is byte-for-byte untouched —
the same evidence D-2 was narrowed for."

### codex-1 round-1 MINOR (D-3) — backslash grammar contradiction — CLOSED

D-3 said backslash "remains refused" while the Python arm strips the extractor's splice
sentinel before matching, so both shipped multi-line commands are accepted.

**Fixed.** D-3 in IMPLEMENTATION.md (line 146) is retitled to include "and backslash
continuations" and states the exception: the extractor marks a continued unit by re-appending
`\`, the Python arm strips it, and both shipped multi-line commands are accepted. Safe because
this arm never executes what it finds — the node arm, which does execute, still refuses
splices. The guard comment in test/design-addons.test.js (line 1109) states the exception
explicitly. The grammar test now asserts both directions: a continued command is accepted
after sentinel stripping, and a non-trailing backslash is still refused.

I independently reconstructed `PUBLISHED_PYTHON` from the test file and verified: the raw
grammar refuses a trailing backslash; after sentinel stripping, the shipped multi-line
continuation is accepted; an internal backslash (`build a\ b`) is still refused; `&&`,
backtick, and `$var` are all refused.

### Round-7 closures — intact

The cycles 8-9 diff is purely additive: zero lines removed from `lib/installer.js`,
`test/bidding-addon.test.js`, and `test/design-addons.test.js` (verified by counting `^-[^-]`
lines in each file's diff). Therefore the two round-7 closures could not have been disturbed:

- **Omitted-`skill` identity (codex-1 round-6 MAJOR).** The guard at lib/installer.js:1446 is
  `state.marker.skill !== unit.skill` with no `undefined` exemption. An absent identity gets
  its own message (line 1457-1458). `installerOwnsDestination` (lib/installer.js:1743)
  requires exact equality. Intact, unchanged.
- **`selected` from the core marker (kimi-1 round-6 MINOR).** `targetSkillUnits`
  (lib/installer.js:904) reads the recorded selection via `markerAddonNames` for read commands.
  Intact, unchanged.

### Round-6 hermes-1 NIT — `valid-unselected` masks `valid-unmanaged` — not closed, not blocking

Unchanged at `ebe269e`, as expected — this is a documented design choice carried as deferred
follow-up #2 in the consensus. The status precedence remains
`malformed > valid-unselected > valid-unmanaged > valid`, with `managed: false`
disambiguating provenance. Not a blocker.

## New findings

None.

## Release judgement

Releasable as 2.1.0. The three findings from codex-1's round-1 review that were never
dispositioned are now fixed and verified: the B5 preflight covers every unit (add-ons and
core), with zero writes on a predictable failure confirmed by independent probe; D-2's
overstated "byte-level" claim is narrowed in both documents; D-3's backslash contradiction is
corrected and asserted in both directions. The shipped payload `skills/parley-bidding/` has
not changed since `714712f` — the entire cycle 8-9 diff touches only installer code and tests.
The full suite passes on both interpreters. No further change is required for release.

## What I verified

1. **316/316 node tests pass, 0 fail, on python3 3.9.6** (system Python).
   `PYTHONDONTWRITEBYTECODE=1 npm test`. The Python leg refuses by design
   ("python3 is 3.9, but the add-on declares >=3.10") and runs zero Python tests.

2. **316/316 node tests pass, 0 fail, on python3 3.14.6** (Homebrew), and **54/54 Python
   tests** across 7 files. Manifest check: `parley-bidding: ok (47 files,
   sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d)`.

3. **B5 preflight — symlink in manifest-free add-on (codex-1 round-1 MAJOR).** Independent
   probe: staged a package with `zz-broken` (manifest-free, sorting last) containing a
   symlink. `installCommand` returned `ok: false`; `~/.codex/skills` did not exist (zero
   writes); message "Refusing to copy symlink in skill payload: link.md".

4. **B5 preflight — symlink in core source (codex-1 round-8 catch).** Independent probe:
   symlink in `skills/parley-deck/references/linked.md`. `ok: false`; zero writes; message
   "Refusing to copy symlink in skill payload: linked.md".

5. **B5 preflight — no false positive.** Clean staging installed all six units normally:
   `ok: true`, six directories present.

6. **B5 preflight is read-only.** SHA-256 inventory of the staged source tree before and
   after a preflight-triggered failure was identical — no writes to the source.

7. **D-3 backslash grammar.** Reconstructed `PUBLISHED_PYTHON` from the test file. Raw
   grammar refuses trailing `\`; after sentinel stripping, the shipped continuation is
   accepted; internal `\` still refused; `&&`, backtick, `$var` all refused.

8. **B6 runtime gate.** With python3 absent (PATH=""): `bidding.runtime.ok: false`,
   `doctor.ok: false` (exit 1). With python3 3.14 on PATH: `bidding.runtime.ok: true`,
   `doctor.ok: true`. Payload stays `valid` in both cases — the gate distinguishes
   payload-valid from operationally-unavailable, as B6 requires.

9. **Payload frozen.** `git diff 714712f..ebe269e -- skills/parley-bidding/` is empty. 48
   tracked files, zero `__pycache__`/`.pyc`/`.pytest_cache`/`.DS_Store`/`node_modules`.

10. **Diff scope.** `git diff 49fc3ec..ebe269e` touches only `lib/installer.js` (+91),
    `test/bidding-addon.test.js` (+59), `test/design-addons.test.js` (+22) — all insertions,
    zero deletions. No drive-by changes to the payload or unrelated code.

11. **Round-7 closures intact.** The identity check (lib/installer.js:1446) and
    `installerOwnsDestination` (lib/installer.js:1743) are unchanged — the additive diff
    could not have disturbed them.

12. **D-2 narrowing in both documents.** IMPLEMENTATION.md D-2 (line 128) and the
    verification table (line 195) no longer claim "byte-level" or "untouched". The consensus
    verification table (line 159) explicitly says "Not proof it is byte-for-byte untouched".
    The "several hermes rounds" claim is gone from the participation section (line 645 now
    says "round 5" only); the sole remaining mention at line 736 is the cycle-9 changelog
    describing the fix.

13. **`npm pack --dry-run`.** 202 total files, 48 under `skills/parley-bidding/`, zero
    caches.

14. **`npx skills add --list`.** "Found 6 skills", `parley-bidding` among them.

15. **Working tree clean.** `git status` empty before and after all probes; no tree mutation.
