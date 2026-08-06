---
agent: codex-1
idea: skill-sync-cli-1-39
review-round: 1
date: 2026-08-06
reviewed-commit: 661af98ef7caa68c0883ce51ba4ae42194967e59
---

## Summary

[PRIMARY — `FINAL.md:23-89`; commit `661af98`; locators and executed checks below] D1, D2, D4,
and D5 land; D1's row and D2's paragraph match the adopted text after normalizing only Markdown
line wrapping. D3's two branches exist, but the unqualified `Headless Agent Configuration`
paragraph contradicts manual branch A, so D3 does not land completely as decided. D6 is also
breached by one incident-specific changelog clause. I found one MAJOR and one MINOR finding.

## Scope and provenance

[PRIMARY] I read the complete live
`/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/parley-deck-cli/parley-deck/COOPERATION.md`
(`1-1316`), including §15 at `1176-1316`; the complete `00-prompt.md`, `FINAL.md`, and
`IMPLEMENTATION.md`; kimi-1's round-02 wording; and the complete commit diff from `661af98^` to
`661af98` in the implementation repository. I inspected all nine changed files, with focused
line-numbered reads of every changed executable and instructional passage.

[PRIMARY] I ran `npm test`, `npm run prepack`, a side-effect-free import check, an exact-text
comparison for D1/D2, `git diff --check`, and a full-commit-tree search for all spellings of
`writeModeArgs`. I also exercised both version gates against a temporary archive of `661af98`
outside both repositories.

[PRIMARY] I did not publish or pack a release, invoke live `opencode`/`kimi` agents, inspect or
migrate external deck configs, or modify any implementation file. I did not re-verdict codex-1's
design-time claims recorded in `FINAL.md:95-123`; this review is limited to the content and
behavior of commit `661af98`.

## Refutation attempts

### D1 and D2 — adopted text

[PRIMARY — `FINAL.md:23-49`; `skills/parley-deck/SKILL.md:243-252`] I extracted the D1 blockquote
and D2 blockquote from `FINAL.md`, stripped only blockquote prefixes, collapsed whitespace caused
by line wrapping, and compared them with the shipped table row and paragraph. The check returned:

```text
d1ExactAfterWhitespaceNormalization: true
d2ExactAfterWhitespaceNormalization: true
d1Length/rowLength: 228/228
d2Length/paragraphLength: 1124/1124
```

[PRIMARY] D1's `opencode` row is therefore not paraphrased, and D2's replacement is not
paraphrased. The original inverted source-of-truth sentence is absent from the changed passage.

### D3 — manual/CLI boundary

[PRIMARY — `FINAL.md:51-60`; `skills/parley-deck/SKILL.md:349-383,813-835`;
`skills/parley-deck/references/WORKED_EXAMPLES.md:16-46`] I tried to follow the new local JSON
shape and then branch A as one manual-facilitator workflow. The configuration paragraph says
`headlessArgs` is the whole invocation and that nothing is appended, while branch A subsequently
instructs the facilitator to append model/thinking/profile flags and deliver the prompt. The
worked example also keeps model metadata outside `headlessArgs` and puts no prompt in that array.
This produced the MAJOR finding below. Branch B itself states the CLI behavior required by D3:
resolved `headless_args` is launched as-is, placeholders are substituted inside it,
`prompt_mode` controls only stdin, and no later arguments are synthesized.

### D4 — equality guard and declared implementation-shape deviation

[PRIMARY — `package.json:3,60-66`; `skills/parley-deck/references/compatibility.json:4`;
`scripts/build-addon-manifest.js:85-105,186-190`; `test/manifest-coverage.test.js:489-500`] The
shipped values are both `2.4.0`. The commit contains one decisive version comparison,
`compat.skillVersion === pkg.version`, in `versionSyncProblem()`, and one named Node-test
assertion that requires the helper and expects `null`.

[PRIMARY] On the unmodified implementation tree, `npm test` exited 0 with this relevant output:

```text
tests 386; pass 386; fail 0
python 3.14: 54 tests OK across 7 files
parley-bidding: ok
parley-deck: ok
parley-design: ok
parley-design-check: ok
parley-tracker: ok
parley-worktrees: ok
```

[PRIMARY] On the same unmodified tree, `npm run prepack` exited 0 and printed `ok` for all six
packaged skills after directly running `node scripts/build-addon-manifest.js --check`.

[PRIMARY] I created a temporary `git archive 661af98` at
`/tmp/skill-sync-cli-1-39-codex-1.vlYz6A`, linked the existing dependencies, and changed only its
`package.json` version from `2.4.0` to `2.4.0-review-desync`. `npm run prepack` exited 1 with:

```text
build-addon-manifest: compatibility.json skillVersion is "2.4.0" but package.json version is "2.4.0-review-desync" — update skills/parley-deck/references/compatibility.json
```

[PRIMARY] In that same temporary mismatch, `npm test` exited 1; its named guard was the sole
failure (`385` pass, `1` fail), and the assertion reported the same mismatch string as its actual
value. I restored the temporary `package.json` to `2.4.0`; its SHA-256 and
`git show 661af98:package.json | shasum -a 256` then both equaled
`aed485da9360f48bed0cb94aed9cfd13be6afd0952e0e0736a212d6c99f8cc61`.

[PRIMARY] The import check
`node -e 'require("./scripts/build-addon-manifest.js"); process.stdout.write("require returned without main side effects\\n")'`
printed only `require returned without main side effects`. Together with the successful direct
`prepack` invocation, this shows the `require.main === module` guard suppresses `main()` only on
import and preserves direct-script execution.

[PRIMARY] I raise no finding on the declared D4 deviation. One shared comparison is exercised by
the existing Node harness and by the existing manifest script; both gates demonstrably fail on a
mismatch, and the commit adds no new script, job, or checker. This is within D4's required behavior
and its explicit permission to put the check inside `build-addon-manifest.js`.

### D5 — exhaustive removal check

[PRIMARY] I did not rely on normal `grep`/`rg`. With `661af98` at clean HEAD, I enumerated every
tracked file in the reviewed commit and fed that NUL-delimited list to the system grep, which does
not consult ignore files:

```text
git ls-tree -r --name-only -z 661af98 |
  xargs -0 /usr/bin/grep -nHE 'writeModeArgs|write_mode_args|WriteModeArgs'
```

[PRIMARY] The only results were `CHANGELOG.md:27,31`, `SKILL.md:383,820`, and
`WORKED_EXAMPLES.md:44`. Each is release-note or migration/removal wording; no documented JSON
field, recipe step, or snake/capitalized variant remains. The pre-change tree's substantive hits
were the JSON field and recipe step in `SKILL.md` and the JSON field in `WORKED_EXAMPLES.md`, and
all three are removed in `661af98`.

[PRIMARY — `skills/parley-deck/SKILL.md:364-383,820`;
`skills/parley-deck/references/WORKED_EXAMPLES.md:29-44`] The shipped migration rule tells users to
merge the legacy arguments into `headlessArgs` and remove `writeModeArgs`, and the write-enabling
argument appears inside `headlessArgs` in both documented JSON shapes. I issue no verdict here on
the design choice codex-1 previously owned; this is raw implementation-coverage evidence.

### D6 and changed-file scope

[PRIMARY] `git diff --name-status 661af98^ 661af98` lists exactly nine modified files:
`CHANGELOG.md`, the two package files, the manifest builder, `SKILL.md`, the generated core
manifest, `WORKED_EXAMPLES.md`, `compatibility.json`, and the existing Node test file.
`git diff --exit-code 661af98^ 661af98 -- skills/parley-deck/references/COOPERATION.md` exited 0
with no output, so the bundled protocol snapshot is untouched as D6 requires.

[PRIMARY — `FINAL.md:86-89`; `CHANGELOG.md:27-31`] One incident-specific clause nevertheless
escaped D6's exclusion; see the MINOR finding.

## Findings

### [CRITICAL] None

[PRIMARY] No CRITICAL finding was found in the reviewed commit and exercised paths.

### [MAJOR] The unscoped JSON guidance contradicts manual branch A

[PRIMARY — `skills/parley-deck/SKILL.md:349-383,815-835`;
`skills/parley-deck/references/WORKED_EXAMPLES.md:16-46`; `FINAL.md:51-60`] The generic JSON
section says, **"`headlessArgs` is the whole invocation"** and **"Nothing is appended to it
afterwards."** Manual branch A then says to add `headlessArgs`, append model/thinking/profile
flags, and deliver the prompt as a later step. The example reinforces the ambiguity: its
`headlessArgs` contains only `--non-interactive` and `--workspace-write`, while model fields remain
separate and no prompt appears in the array. D3 required the manual/CLI boundary to apply to both
`Headless Agent Configuration` and `WORKED_EXAMPLES.md`; the absolute paragraph instead applies
branch B's CLI rule to the manual JSON shape too. A facilitator cannot follow both instructions,
and may duplicate model/prompt arguments or omit them.

Fix: label the camel-case JSON shape and worked example as manual-facilitator inputs. State there
that the write-enabling arguments belong inside `headlessArgs`, while branch A still appends any
separately configured model/thinking/profile arguments and delivers the prompt according to
`promptMode`. Reserve **complete argv template / nothing appended** for branch B's snake-case
Parley CLI `headless_args`, or show separate manual and CLI shapes in the configuration section.

### [MINOR] D6's `hermes` incident exclusion leaked into the changelog

[PRIMARY — `CHANGELOG.md:27-31`; `FINAL.md:86-89`] D6 says the `hermes` incident narrative stays
out because it changes no next action, but the release note says the removed list is the model
that made **"the `hermes` regression invisible."** This is incident-specific history that FINAL
explicitly excluded. It does not affect runtime behavior, but it exceeds the ratified release
scope.

Fix: delete the incident clause or make it mechanism-only, for example: "A separate list teaches
the two-list launch model the CLI does not implement."

### [NIT] None

[PRIMARY] No additional NIT finding was found.

## Open questions

[PRIMARY] None. The two findings have local fixes that do not require a new design decision.
