---
agent: hermes-1
idea: skill-sync-cli-1-39
review-round: 2
date: 2026-08-06
reviewed-commit: b769ced
responding-to: [review/consensus.md]
---

## Summary

Verdict: **accept**. Zero agreed fixes remain open.

The fix-up commit b769ced (on top of the reviewed 661af98) applied AF-1, AF-2, and AF-3 — the
three fixes that touched the skill repository — correctly and without introducing new
contradictions. AF-4 and AF-5 were documentation corrections to `IMPLEMENTATION.md` in the CLI
repo and are also resolved. `npm test` passes (386/386, 0 fail). `npm run prepack` passes (all 6
addon manifests `ok`). No code or test files were touched by the fix-up; only `CHANGELOG.md`,
`SKILL.md`, and the regenerated `parley-addon.json` changed.

## Scope

Declared scope: Phase 8 re-review of fix-up cycle 1, commit b769ced, against the review consensus
AF-1 through AF-5 and the dismissed NIT.

What I read (PRIMARY):

- `review/consensus.md` — the agreed fixes AF-1..AF-5, the D3 verdict conflict resolution, and the
  dismissed NIT.
- `review/round-01/hermes-1.md` — my own round-1 findings.
- `FINAL.md` — decisions D1-D6, dissent, verification record.
- `IMPLEMENTATION.md` — updated with AF-4 and AF-5 corrections.
- `git show b769ced` — the full fix-up diff (3 files: CHANGELOG.md, SKILL.md, parley-addon.json).
- `git diff 661af98..b769ced --stat` — confirms only 3 files changed, 8 insertions, 8 deletions.
- `git diff 661af98..b769ced -- scripts/ test/ package.json` — empty; no code/test/package changes.
- `skills/parley-deck/SKILL.md` — lines 240-420 (Autonomous Execution, Headless Agent
  Configuration) and lines 800-870 (Generic CLI Invocation Contract, branches A and B).
- `skills/parley-deck/references/WORKED_EXAMPLES.md` — full (100 lines).
- `skills/parley-deck/references/compatibility.json` — `skillVersion: "2.4.0"`.
- `package.json` — `version: "2.4.0"`, `test` and `prepack` scripts.
- `CHANGELOG.md` — lines 20-35.
- `skills/parley-deck/parley-addon.json` — current vs. pre-fix-up hashes.

What I ran (PRIMARY):

- `npm test` → 386 node tests pass, 0 fail. Python tests skipped (python3 is 3.9, add-on requires
  >=3.10; exit 0, not a failure). The test `compatibility.json skillVersion tracks package.json
  version` is in the pass list.
- `npm run prepack` → exit 0, all 6 addon manifests `ok`.
- Python script comparing FINAL.md's adopted D2 text against SKILL.md line 252 byte-for-byte.

What I did NOT do:

- I did not re-run the D4 guard-fire desync test. The fix-up diff proves `scripts/build-addon-
  manifest.js`, `test/manifest-coverage.test.js`, `package.json`, and `compatibility.json` are
  unchanged from 661af98 (empty diff). The guard I verified in round 1 is the same code. (PRIMARY)
- I did not run the Python test suite (same environment constraint as round 1: python 3.9 < 3.10).

## My round-1 findings — resolved status

### [MINOR] D2's separated "A vendor flag change is a config edit, not a skill revision." line
(AF-3)

**Resolved.** (PRIMARY) The fix-up folded the sentence back into the D2 paragraph as its final
sentence. I verified this with a byte-for-byte comparison: FINAL.md's adopted D2 text (1124
characters) + " A vendor flag change is a config edit, not a skill revision." (61 characters
including leading space) = SKILL.md line 252 (1185 characters). Exact match. The standalone line
at the old SKILL.md:254 is gone; the paragraph is now a single line again, matching the pre-commit
structure where the clause was part of one sentence.

### [MINOR] The `npm test` narrative understated what a single run exercises (AF-4)

**Resolved.** (PRIMARY) `IMPLEMENTATION.md` now contains a "Correction (review AF-4, hermes-1)"
section stating that `npm test` exercises the guard through both callers — the test assertion and
the `--check` invocation — in a single run, and that `prepack` exercises only the script caller.
This matches what I found in round 1: `package.json:60` is `node --test && node scripts/run-
python-tests.js && node scripts/build-addon-manifest.js --check`.

### [NIT] D2 paragraph is a 1124-character single line (dismissed)

**Not fixed, correctly.** (PRIMARY) The NIT was dismissed by consensus: rewrapping would break the
verbatim match. The paragraph is now 1185 characters (with the appended AF-3 sentence) and remains
a single unbroken line. This is the right call — the verbatim property is worth more than diff
ergonomics. No action needed.

## AF-1 — the MAJOR (codex-1): scope the absolute "nothing is appended" paragraph

**Resolved. No new contradiction introduced.** (PRIMARY)

Pre-fix-up (at 661af98), the manual JSON section read:

> **`headlessArgs` is the whole invocation.** There is no separate write-mode argument list: the
> flag that lets an agent write its own artifact must be **inside** `headlessArgs`, alongside
> everything else the launch needs. Nothing is appended to it afterwards.

The absolute "Nothing is appended to it afterwards" applied to the manual JSON shape, but branch A
steps 3-4 explicitly append model/thinking/profile flags (step 3) and the prompt (step 4) after
`headlessArgs`. That was the contradiction codex-1 found.

Post-fix-up (at b769ced), the section now reads:

> This shape is **manual-facilitator input**: it is what you read when you assemble and run the
> command yourself (branch A of "Generic CLI Invocation Contract"). The Parley CLI reads its own
> snake-case configuration instead, where `headless_args` is the complete argv template and
> nothing is appended to it — see branch B.
>
> **There is no separate write-mode argument list.** The flag that lets an agent write its own
> artifact belongs **inside** `headlessArgs`. Model, thinking and profile flags remain separate
> fields and are appended by branch A at launch; the write-enabling flag is not one of them.

I checked the four-way consistency the task asked for:

1. **New paragraph vs. branch A.** The new paragraph says "Model, thinking and profile flags
   remain separate fields and are appended by branch A at launch." Branch A step 3 says "Add
   model/thinking/profile flags only when discovered or configured." These agree: the manual path
   appends model/thinking/profile flags after `headlessArgs`. The write-enabling flag goes inside
   `headlessArgs` (new paragraph) — branch A step 2 says "Add `headlessArgs` — including the flag
   that lets the agent write its own artifact." Consistent. (PRIMARY)

2. **New paragraph vs. branch B.** The new paragraph says "The Parley CLI reads its own snake-case
   configuration instead, where `headless_args` is the complete argv template and nothing is
   appended to it — see branch B." Branch B says "Nothing is appended afterwards — no permission
   flag, no model flag, no thinking flag, no profile flag, no separate write-mode list." The
   "nothing appended" rule is now correctly attributed to branch B (the CLI path), not to the
   manual section. Consistent. (PRIMARY)

3. **New paragraph vs. WORKED_EXAMPLES.md.** WORKED_EXAMPLES.md says "the write-enabling flag
   (`--workspace-write` above) sits **inside** `headlessArgs`. There is no separate write-mode
   argument list." The new paragraph says the same: "The flag that lets an agent write its own
   artifact belongs **inside** `headlessArgs`." Consistent. (PRIMARY)

4. **Branch A vs. branch B.** Branch A assembles a command in 5 steps (cli + headlessArgs +
   model/thinking/profile + prompt + timeout). Branch B launches `headless_args` as-is with
   nothing appended. These describe two different activities and no longer contradict each other
   within one section. The old contradiction — the manual section claiming "nothing is appended"
   while branch A appends — is gone. (PRIMARY)

The `headlessArgs` field itself remains in the JSON shape (SKILL.md:363), which is correct: AF-1
scoped the absolute paragraph, not the field.

## AF-2 — MINOR (codex-1 and kimi-1): remove the `hermes` allusion from CHANGELOG

**Resolved.** (PRIMARY) The CHANGELOG previously read "A separate list is precisely the mental
model that made the `hermes` regression invisible." It now reads "A separate list teaches a
two-list launch model the CLI does not implement." The `hermes` incident narrative is gone.
`/usr/bin/grep -c "hermes" CHANGELOG.md` returns 0.

## AF-3 — MINOR (hermes-1): fold the vendor-flag sentence back into the D2 paragraph

**Resolved.** (PRIMARY) See "My round-1 findings" above. The D2 paragraph at SKILL.md:252 is now a
verbatim match against FINAL.md's adopted D2 text with the vendor-flag sentence appended as the
final sentence. 1185 characters, exact match.

## AF-4 — MINOR (hermes-1): correct the IMPLEMENTATION.md narrative

**Resolved.** (PRIMARY) `IMPLEMENTATION.md` now contains a "Correction (review AF-4, hermes-1)"
paragraph stating that `npm test` exercises the guard through both callers in a single run, and
`prepack` exercises only the script caller. This is documentation accuracy only, no code change.

## AF-5 — NIT (kimi-1): note the `headless_args` → `headlessArgs` adaptation

**Resolved.** (PRIMARY) `IMPLEMENTATION.md` now contains an "Adaptation (review AF-5, kimi-1)"
paragraph recording that FINAL's snake_case `headless_args` was applied to the manual JSON shape as
camelCase `headlessArgs`, and that this is a deliberate adaptation of ratified text to its target
context.

## Did the fix-up break anything?

No. (PRIMARY)

- `npm test` → 386 pass, 0 fail. Exit 0.
- `npm run prepack` → all 6 addon manifests `ok`. Exit 0.
- The fix-up diff touched only `CHANGELOG.md`, `SKILL.md`, and `parley-addon.json` (the
  regenerated hash). No code, test, or package files changed.
- `compatibility.json` `skillVersion` is `"2.4.0"`, matching `package.json` `version: "2.4.0"`.
- The `parley-addon.json` aggregate hash was correctly regenerated (`e6de5305...` →
  `73e5d925...`) and `npm run prepack` verifies the payload integrity.

## Anything new shipped in the fix-up that the consensus did not agree?

No. (PRIMARY) The fix-up diff is exactly 3 files with 8 insertions and 8 deletions:

1. `CHANGELOG.md` — AF-2 (removed the `hermes` allusion, replaced with the mechanism description).
2. `skills/parley-deck/SKILL.md` — AF-1 (scoped the absolute paragraph, labelled the shape as
   manual-facilitator input) and AF-3 (folded the vendor-flag sentence back into the D2 paragraph).
3. `skills/parley-deck/parley-addon.json` — regenerated hash, a mechanical consequence of the
   SKILL.md edit.

AF-4 and AF-5 were applied to `IMPLEMENTATION.md` in the CLI repo (commit 3bda823), not to the
skill repo. Both are documentation corrections the consensus agreed.

No new content, no new files, no scope creep. The `parley-addon.json` regeneration is a mechanical
artifact of the SKILL.md edit, not an independent change.

## Position changes since prior review round

No position changes. All three of my round-1 findings are resolved (two by fix-up, one dismissed
by consensus and correctly left alone). The MAJOR filed by codex-1 (AF-1) is resolved without
introducing a new contradiction. I have no new findings.

## Updated findings

### [CRITICAL]

None.

### [MAJOR]

None.

### [MINOR]

None.

### [NIT]

None.

## Open questions

None. The Python test suite remains unverified on this machine (python 3.9 < 3.10), same as round
1. This is an environment constraint, not a code defect, and the fix-up did not touch any Python
code or tests.
