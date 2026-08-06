---
agent: hermes-1
idea: skill-sync-cli-1-39
review-round: 1
date: 2026-08-06
reviewed-commit: 661af98
---

## Summary

The implementation ships D1–D5 faithfully. D1's opencode row and D2's replacement paragraph are
verbatim matches against FINAL.md's adopted text. D3's manual/CLI branch split is structurally
correct. D5's `writeModeArgs` removal is complete — enumerated with `find` + `/usr/bin/grep`, not
the facilitator's aliased `grep`. D4's guard is real: I desynced `compatibility.json` to prove
`prepack` exits 1 and the test assertion fails, then restored it exactly. The one declared
deviation (a single `versionSyncProblem()` function with two callers instead of two independent
copies) is within D4's decision. No CRITICAL or MAJOR findings. Two MINORs and one NIT.

## Scope

What I read (PRIMARY):

- `COOPERATION.md` §15 (lines 1176–1316) — verification integrity rules.
- `FINAL.md` (full) — decisions D1–D6, dissent, verification record.
- `IMPLEMENTATION.md` (full) — claimed work and declared deviation.
- `git show 661af98` — full diff of all 9 changed files.
- `scripts/build-addon-manifest.js` — full file at HEAD (190 lines), and the pre-commit tail
  (`git show 661af98~1:scripts/build-addon-manifest.js | tail -5`).
- `skills/parley-deck/SKILL.md` — lines 240–289 (Autonomous Execution + table), 360–389 (Headless
  Agent Configuration), 800–849 (Generic CLI Invocation Contract). Pre-commit equivalent of lines
  249–255.
- `skills/parley-deck/references/WORKED_EXAMPLES.md` — full (100 lines).
- `package.json` — full (72 lines).

What I ran (PRIMARY):

- `npm test` → 386 node tests pass, 0 fail. Python tests skipped (python3 is 3.9, add-on requires
  >=3.10; exit 0, not a failure). The new test
  `compatibility.json skillVersion tracks package.json version` is in the pass list.
- `npm run prepack` → exit 0, all 6 addon manifests `ok`.
- Guard-fire test: desynced `compatibility.json` `skillVersion` to `"2.3.0"`, ran `npm run
  prepack` → exit 1 with
  `build-addon-manifest: compatibility.json skillVersion is "2.3.0" but package.json version is
  "2.4.0" — update skills/parley-deck/references/compatibility.json`. Also ran
  `node -e "...versionSyncProblem()"` → returned the same problem string. Restored the file
  from a `/tmp` backup; verified sha256 `f37af38b...` matches the committed version and
  `git status` is clean. Post-restore `npm run prepack` → exit 0.
- `writeModeArgs` enumeration: `find . -path ./.git -prune -o -path ./node_modules -prune -o -type
  f ... -print | xargs /usr/bin/grep -l "writeModeArgs"` and
  `/usr/bin/grep -rn "writeModeArgs" ...`. Tool named: `/usr/bin/grep` and `find` (not the
  facilitator's aliased `grep`).
- Verbatim text comparison: Python script comparing FINAL.md's adopted D1/D2 text against SKILL.md
  byte-for-byte.

What I did NOT do:

- I did not run the Python test suite (skipped by the runner on this machine due to python 3.9 <
  3.10). The IMPLEMENTATION.md claim of "54 python tests across 7 files" is not verified by me.
- I did not inspect `CHANGELOG.md` for correctness beyond confirming it exists and mentions
  `writeModeArgs` in narrative/migration context only.
- I did not verify the `parley-addon.json` aggregate hash by recomputing it independently; I
  relied on `npm run prepack` (which runs `build-addon-manifest.js --check` and verifies the
  payload) passing as evidence the manifest is consistent.

## Refutation attempts

**D4 acceptance criterion: "set `skillVersion` to the version actually shipped, and add exactly one
equality assertion `compatibility.skillVersion === package.version` … No new script, job, or
checker." + "close the `prepack` half."**

- Attempted to break: desynced `compatibility.json` to `2.3.0` while `package.json` stays `2.4.0`.
  Result: `npm run prepack` exits 1 (guard fires). `node -e "...versionSyncProblem()"` returns the
  problem string. The test `compatibility.json skillVersion tracks package.json version` would fail
  (the function returns non-null, `assert.equal(null)` fails). Both gates closed. Restored
  afterwards. (PRIMARY)
- Attempted to break: checked whether a second copy of the comparison exists (FINAL says "exactly
  one equality assertion"). There is one `versionSyncProblem()` function in
  `scripts/build-addon-manifest.js:89-97`, called from `main()` (line 103) and required by
  `test/manifest-coverage.test.js:498`. No duplicate comparison logic. (PRIMARY)

**D5 acceptance criterion: "The field is removed from the documented config shape … carries the
migration rule."**

- Attempted to break: searched for any remaining `writeModeArgs` field in a config shape or code.
  `find` + `/usr/bin/grep` found references only in: (a) `CHANGELOG.md` narrative (lines 27, 31),
  (b) `WORKED_EXAMPLES.md:44` migration note, (c) `SKILL.md:381,383,820,833` — all in the "there is
  no separate write-mode list" / migration context. The JSON shape at `SKILL.md:362-374` no longer
  has the field. `WORKED_EXAMPLES.md:31` has `--workspace-write` inside `headlessArgs`. No field
  survives as a live config key. (PRIMARY)

**D1 verbatim match (kimi-1's wording).**

- Attempted to break: compared FINAL.md's blockquoted adopted text against SKILL.md:250
  byte-for-byte. Exact match (the blockquote wraps in FINAL; when joined, the text is identical).
  (PRIMARY)

**D2 verbatim match (kimi-1's converged text).**

- Attempted to break: compared FINAL.md's blockquoted replacement paragraph against SKILL.md:252
  byte-for-byte. 1124 characters, exact match. (PRIMARY)

**D6: `references/COOPERATION.md` untouched.**

- Attempted to break: `git diff 661af98~1..661af98 --stat -- skills/parley-deck/references/
  COOPERATION.md` → empty. `shasum -a 256` of the file at both commits → identical
  (`a72bc52e...`). (PRIMARY)

**`require.main === module` guard does not change existing behaviour.**

- Pre-commit: `main()` was called unconditionally at the bottom of the file. Post-commit: `main()`
  is called inside `if (require.main === module)`. When the script is run directly
  (`node scripts/build-addon-manifest.js …`), `require.main === module` is true and `main()`
  executes — identical to before. When the script is `require()`d by the test, `main()` does not
  execute — this is the new behaviour that enables the test to call `versionSyncProblem()` without
  triggering the full script. The `npm test`, `npm run prepack`, `npm run manifest:addons`, and
  `npm run manifest:check` scripts all invoke the file as a direct process, so the guard is true
  and behaviour is unchanged. (PRIMARY)

## Findings

### [CRITICAL]

None. The guard fires, the tests pass, the decisions land, and the declared deviation is within
D4's scope.

### [MAJOR]

None.

### [MINOR]

### [MINOR] D2's separated "A vendor flag change is a config edit, not a skill revision." line is an
undeclared structural addition

**File:** `skills/parley-deck/SKILL.md:254`
**What is wrong:** FINAL.md D2 says the adopted replacement "replaces **both** the inverted
sentence and the confinement-only fail-closed sentence, rather than adding a third paragraph."
The pre-commit text was a single sentence containing three clauses: (1) the inverted
source-of-truth claim, (2) "a vendor flag change is a config edit, not a skill revision", (3) the
confinement fail-closed rule. The implementation replaced clauses (1) and (3) with the new D2
paragraph (verbatim match, confirmed), and extracted clause (2) as a standalone line
(`SKILL.md:254`). This is not a paraphrase or content change — the words are identical — but it is
a structural change (one sentence → two paragraphs) that FINAL does not explicitly decide. The
IMPLEMENTATION.md deviation section does not mention it.

This is minor because: the clause's meaning is unchanged, the separation arguably improves
readability, and FINAL's "rather than adding a third paragraph" refers to not adding a *new*
paragraph of content, not to prohibiting the existing clause from being line-broken out. But a
reviewer checking "does D2 land as decided" will notice the extra line and wonder whether it was
intended.

**Fix:** Either (a) fold "A vendor flag change is a config edit, not a skill revision." into the
end of the D2 paragraph as a final sentence (closest to the pre-commit structure), or (b) add a
one-line note to IMPLEMENTATION.md's deviation section acknowledging the structural split. (a) is
preferred since it makes the replacement truly drop-in.

### [MINOR] The `npm test` script already runs `build-addon-manifest.js --check`, making the test
file's assertion the third gate, not the second

**File:** `test/manifest-coverage.test.js:498`, `package.json:60`
**What is wrong:** IMPLEMENTATION.md describes two callers of `versionSyncProblem()`: `main()` (for
`prepack`) and the test file (for `npm test`). But `package.json:60` defines `test` as
`node --test && node scripts/run-python-tests.js && node scripts/build-addon-manifest.js --check`.
So during `npm test`, the guard fires twice: once via the test-file assertion (`node --test`) and
once via the `--check` call at the end of the script. This is not wrong — redundancy is harmless
and the double-fire is exactly the kind of belt-and-suspenders D4 wants — but the IMPLEMENTATION.md
narrative of "one comparison with two callers" understates it: `npm test` exercises the function
through both callers in a single run.

**Fix:** No code change needed. Optionally, IMPLEMENTATION.md could note that `npm test` exercises
both the test assertion and the `--check` script call, so the single `npm test` invocation hits
the guard twice. This is documentation accuracy, not a correctness issue.

### [NIT]

### [NIT] D2 paragraph is a 1124-character single line

**File:** `skills/parley-deck/SKILL.md:252`
**What is wrong:** The D2 replacement paragraph is one unbroken line of 1124 characters. This
matches FINAL's adopted text verbatim (so it is not a deviation), but it is hard to read in a
diff or editor, and a future edit to one clause will produce a noisy diff. This is a readability
nit only — the content is correct.

**Fix:** Soft-wrap the paragraph at sentence boundaries in a future revision. Do not change it now;
any rewrap would break the verbatim match this review confirmed.

## Open questions

1. **The Python tests were not run on this machine** (python 3.9 < the add-on's >=3.10 floor; the
   runner skips them with exit 0). IMPLEMENTATION.md claims "54 python tests across 7 files" pass.
   I cannot confirm or refute this. If a reviewer with python >=3.10 is available, running
   `node scripts/run-python-tests.js` would close the gap. This is not a finding against the
   implementation — the skip is an environment constraint, not a code defect.

2. **Is the `versionSyncProblem()` placement in `main()` before the `--check` vs. write
   decision correct?** The function is called at `build-addon-manifest.js:103`, before the
   `if (check)` branch at line 155. This means the version sync fails even on a
   `manifest:addons` (write) run, not just `--check`. That seems intentional — a desynced version
   should block manifest regeneration too — but FINAL D4 only mentions `prepack` and the test
   harness. The implementer may want to confirm this is desired, though I consider it correct
   behaviour (a write run that succeeds while the version is desynced would produce a manifest
   certifying a broken state).
