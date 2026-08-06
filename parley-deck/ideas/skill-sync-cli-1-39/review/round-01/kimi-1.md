---
agent: kimi-1
idea: skill-sync-cli-1-39
review-round: 1
date: 2026-08-06
reviewed-commit: 661af98ef7caa68c0883ce51ba4ae42194967e59
---

# Review — skill-sync-cli-1-39, round 1 (kimi-1)

## Summary

The implementation lands D1–D5 as decided: both kimi-1-worded texts (D1 row, D2 paragraph) are
byte-identical to the adopted wording after whitespace normalization, the D3 branch split matches
the CLI's actual launch path, the D5 field is gone from every documented shape with the migration
rule in place, and the D4 guard is real — I reproduced it firing at both `prepack` and `node
--test` on a desynced copy. One MINOR: the CHANGELOG ships a `hermes`-incident allusion that D6
excluded, undeclared in `IMPLEMENTATION.md`.

## Scope

Read: `COOPERATION.md` §15 (and the Phase 6 format at §4); `FINAL.md`; `IMPLEMENTATION.md`; the
full diff of `661af98` in `parley-deck-skill`; the changed regions of
`skills/parley-deck/SKILL.md` (lines 205–234, 238–265, 355–395, 795–845), all of
`references/WORKED_EXAMPLES.md`, `scripts/build-addon-manifest.js` (whole file),
`test/manifest-coverage.test.js` (tail), `package.json`, `RELEASING.md`; and, as CLI ground truth,
`internal/runner/runner.go` (lines 1020–1139), `internal/app/app.go` (`agents` command), and
`internal/agents/discover.go:539` in `parley-deck-cli`.

Ran: `npm test` and `npm run prepack` in the real skill repo; a desync experiment on an `rsync`
copy under `/tmp/kimi-1-review/` (both repositories untouched; see D4 below); verbatim text
comparisons via a Python script in `/tmp`; absence checks via `find … -print0 | xargs -0
/usr/bin/grep` (named per FINAL follow-up 2 — no `grep -r`/bare `rg` was relied on for any
negative claim).

Not done: I did not read the other reviewers' round-01 review files (independence); I did not
audit pre-existing SKILL.md content outside the changed regions; I did not run the Python test
files standalone (they ran inside `npm test`); I did not execute any git write command.

§15.1 note: I authored the D1/D2 wording adopted in FINAL, but the verdicts below are about the
*implementation's* content — claims first asserted by the implementer in commit `661af98`, which I
checked `PRIMARY`. I issue no verdict on any FINAL-era claim I own; where those matter here
(e.g. what `package.json:60`/`:66` run) I re-derived them from the current files.

## Refutation attempts

### D1 — opencode row lands, verbatim

Compared FINAL.md §D1's adopted blockquote against `skills/parley-deck/SKILL.md:250` with a
whitespace-normalizing Python script. Result: **identical** — `| opencode | `run --auto` — the
prompt is an argv positional, not stdin. … what may change between versions. |`. `PRIMARY`
(script: extract `> ` lines between `### D1`/`### D2`, collapse whitespace, string-compare; output
`D1 MATCH : True`).

### D2 — replacement paragraph lands, verbatim, replacing both sentences

Same method, FINAL.md §D2 blockquote vs `SKILL.md:252`: **identical** (`D2 MATCH : True`).
`PRIMARY`. The old inverted sentence and the old confinement-only sentence are both gone from the
paragraph; the surviving vendor-flag clause was kept as its own line at `SKILL.md:254` (*"A vendor
flag change is a config edit, not a skill revision."*), which D2 neither required nor forbade.
`PRIMARY` (diff hunk `@@ -247,8 +247,11`).

### D3 — manual/CLI branch split, checked against the CLI itself

`SKILL.md:815–835` carries branch A (manual, 5 steps, write-enabling flag folded into
`headlessArgs` at step 2) and branch B (CLI). Branch B's three claims verified against the CLI,
`PRIMARY`:

- *"resolved `headless_args` is the complete argv template … launched as-is"* and *"nothing is
  appended afterwards"* — `internal/runner/runner.go:1097–1122` (`buildAgentInvocation`) builds
  argv **only** from `agent.HeadlessArgs`, substituting whole-token `{root}`/`{prompt}`; no other
  argument is added. `execAgentProcess` (`runner.go:1031`) launches exactly that vector.
- *"`prompt_mode` only decides whether the prompt is wired to stdin"* — `runner.go:1056–1058`:
  `if agent.PromptMode == agents.PromptStdin { cmd.Stdin = strings.NewReader(prompt) }`. No other
  use of `PromptMode` affects argv.
- D2's references exist: `parley agents list` is a real command (`internal/app/app.go:121`,
  `:420`), and `AUTO=no` is really surfaced when a declared mode's enabling args are absent
  (`internal/agents/discover.go:539`).

The "same boundary" for the config shape is at `SKILL.md:381` (*"`headlessArgs` is the whole
invocation … Nothing is appended to it afterwards"*), and `WORKED_EXAMPLES.md` carries the updated
example plus the note at line 44. `PRIMARY` (file reads cited above).

### D4 — bump + guard, attacked on a copy

- Real repo, unmodified: `npm test` → **386 node tests pass, 0 fail** (including the new
  `✔ compatibility.json skillVersion tracks package.json version`), 54 python tests across 7
  files OK, all six add-on manifests `ok` (parley-deck aggregate
  `sha256:e6de5305…`, matching the regenerated `parley-addon.json` in the diff). `npm run prepack`
  → exit 0. `PRIMARY` (commands and output quoted; IMPLEMENTATION.md's numbers reproduced exactly).
- Guard attack: copied the repo (`rsync -a --exclude .git --exclude dist`) to
  `/tmp/kimi-1-review/repo-copy`, edited **only the copy's**
  `skills/parley-deck/references/compatibility.json` to `"skillVersion": "2.3.0"`, then:
  - `npm run prepack` → exit **1**, output `build-addon-manifest: compatibility.json skillVersion
    is "2.3.0" but package.json version is "2.4.0" — update
    skills/parley-deck/references/compatibility.json`. No per-skill `ok` lines — the check fires
    before the manifest loop. `PRIMARY`.
  - `node --test test/manifest-coverage.test.js` → the D4 test **fails** (`✖ compatibility.json
    skillVersion tracks package.json version`), plus 9 payload-integrity tests fail because
    `compatibility.json` is inside the hashed payload — exactly the side effect IMPLEMENTATION.md
    reported. `PRIMARY`.
  - Restoration: copied the backup back, `npm run prepack` → exit 0 again, then deleted the copy.
    **The real repositories were never desynced**; `skillVersion` in the real repo reads `2.4.0`
    (`PRIMARY`, grep of the file after cleanup).
- `require.main === module` guard: `node -e "require('./scripts/build-addon-manifest.js')"` (in
  the copy) exits 0, runs no manifest work, and exposes `versionSyncProblem` as a function.
  `PRIMARY`. No pre-existing consumer could observe the change: `find . -type f -not -path
  './.git/*' -print0 | xargs -0 /usr/bin/grep -l build-addon-manifest` over the whole tree matches
  only `test/manifest-coverage.test.js` (the require is *added by this commit*), `CHANGELOG.md`,
  `package.json`, and the script itself — `bin/` and `lib/` never reference it. `PRIMARY`
  (find-based enumeration, tool named). Direct execution still enters `main()` — proven by both
  prepack runs above. Behavior preserved.

**Is the declared deviation within the decision?** Yes, in my reading. D4 sanctioned "putting the
equality check inside `scripts/build-addon-manifest.js`" for the `prepack` half, and asked for
"exactly one equality assertion" in the test harness. The implementation ships exactly one new
test with one assertion (`test/manifest-coverage.test.js:497–500`) and exactly one comparison
(`scripts/build-addon-manifest.js:93`), shared by both gates — no new script, job, or checker
file. The deviation (single shared function instead of two literal copies) is declared in
`IMPLEMENTATION.md` and is a shape difference: both gates demonstrably close (attacks above). This
is my reviewer assessment, not a verdict on a claim I own.

### D5 — writeModeArgs removal, enumerated the ignore-proof way

`find . -type f -not -path './.git/*' -print0 | xargs -0 /usr/bin/grep -l <pattern>` (BSD
`/usr/bin/grep`, no ignore-file honoring; `node_modules` included) over the skill repo:

- `writeModeArgs`: only `SKILL.md:383` and `SKILL.md:820` (the migration rule, as decided),
  `references/WORKED_EXAMPLES.md:44` (the migration note, as decided), and `CHANGELOG.md:27,31`
  (the release note describing the removal). The documented JSON shapes (`SKILL.md:361–376`,
  `WORKED_EXAMPLES.md:30–42`) are clean. `PRIMARY`.
- `write_mode_args`: zero matches. `PRIMARY`.
- `writeMode`: the three files above, plus `SKILL.md:220` — a capability-*discovery* item
  ("`writeMode`: how to allow narrow workspace writes for one protocol artifact"), not the removed
  config field; the concept legitimately survives and its answer now lands in `headlessArgs` per
  `SKILL.md:381`/`:820`. Not a finding. One unrelated vendored hit (`node_modules/postject`
  stream API). `PRIMARY`.
- CLI side: the same find-based enumeration in `parley-deck-cli` finds **no `.go` file** declaring
  `writeModeArgs`/`write_mode_args` (matches are deck artifacts, logs, and the live deck's own
  `meta/headless-agents.local.json` — follow-up 1 territory, out of scope here). The CHANGELOG's
  "No such field exists in the CLI" holds for code. `PRIMARY`.

### D6 — what stayed out

`git show 661af98 --stat` lists 9 files; `skills/parley-deck/references/COOPERATION.md` is not
among them — untouched as decided. `PRIMARY`. No promotion history, probe versions/outputs, ACP
availability, kimi exit-1 message, or CLI struct names appear in the shipped docs (`PRIMARY`, full
diff read). One exclusion is violated — see MINOR-1.

### Internal consistency — fail-closed paragraph vs branch B

No contradiction. `SKILL.md:252` (D2) says a config override can replace launch arguments
wholesale and silently drop the enabling flag, and mandates fail-closed (`AUTO=no`) when the
effective argv cannot be inspected or lacks a required argument. Branch B (`SKILL.md:835`) states
the same wholesale-override hazard and explicitly cites it as the reason the Autonomous Execution
check reads the effective argv. The two texts reinforce each other; branch B's CLI claims are
accurate per the `runner.go` citations above. `PRIMARY`.

## Findings

### CRITICAL

None.

### MAJOR

None.

### [MINOR] CHANGELOG ships the D6-excluded `hermes` incident allusion, undeclared

- **File/line:** `CHANGELOG.md:29` — *"A separate list is precisely the mental model that made the
  `hermes` regression invisible."*
- **What is wrong:** D6 decides that "the `hermes` incident narrative" stays out of this change.
  The 2.4.0 release notes nonetheless name the incident. `IMPLEMENTATION.md` declares exactly one
  deviation (the D4 shape) and does not declare this one. Severity is MINOR, not higher: the text
  sits in release notes rather than skill guidance, it is one clause, and it changes no next
  action — which was D6's own rationale for exclusion. If the deck reads D6 as scoped only to
  `SKILL.md`/`WORKED_EXAMPLES.md` content, downgrade to NIT; I read "what stays out" as covering
  the shipped change. `PRIMARY` (grep output: `29:  model that made the `hermes` regression
  invisible.`; D6 text at `FINAL.md:88`).
- **Fix:** drop or neutralize the clause, e.g. *"A separate list is precisely the two-list mental
  model that 1.39.0 removed."* — or record a deviation in `IMPLEMENTATION.md` if the allusion is
  deliberately kept.

### [NIT] D5 migration rule adapted `headless_args` → `headlessArgs` without a note

- **File/line:** `skills/parley-deck/SKILL.md:383` and `:820`,
  `references/WORKED_EXAMPLES.md:44`.
- **What is wrong:** FINAL.md's adopted migration rule (quoted at `FINAL.md:83–84`) says merge
  into "that agent's `headless_args`"; the implementation writes `headlessArgs`. The adaptation is
  *correct in context* — the documented JSON shape being edited uses camelCase (`SKILL.md:364`),
  so the snake_case spelling would have sent readers to a nonexistent key — but the difference
  from the adopted quote is nowhere noted. `PRIMARY` (normalized comparison script output quoted
  in my transcript; D5 FINAL vs SKILL texts differ exactly in `a writeModeArgs`/`headlessArgs`
  articles-plus-spelling and the added rationale sentence).
- **Fix:** none required for correctness. Optionally, one parenthetical in `IMPLEMENTATION.md`
  ("migration rule spelled `headlessArgs` to match the JSON shape") so the adaptation is on the
  record rather than silent.

## Observations (no severity)

- D3's adopted text lists "no permission, model, thinking, **profile or prompt** arguments …
  synthesized afterwards"; branch B's nothing-appended list (`SKILL.md:833`) omits "prompt". Not a
  gap: branch B's first two bullets already fix the prompt's only two channels (a `{prompt}` token
  inside the template, or stdin wiring), so no prompt argument can be synthesized. `PRIMARY`.
- Branch B's "substituted **inside** `headless_args`" traces to D3's adopted wording; the CLI
  performs *whole-token* substitution (`runner.go:1100–1107` matches the entire arg), so a config
  embedding `{prompt}` inside a larger token (e.g. `--prompt={prompt}`) would not be substituted.
  The wording is the ratified decision's own, so this is not an implementation finding — noted in
  case a future idea wants to tighten it. `PRIMARY`.
- `IMPLEMENTATION.md`'s "only two files in the repository ever mentioned it" was accurate for the
  pre-change scope it enumerated; post-commit the string also lives in `CHANGELOG.md` (the removal
  note) — the enumeration method (find-based) and conclusion both hold. `PRIMARY`.

## Open questions

1. MINOR-1's severity hinges on D6's scope: does "what stays out" cover the release notes, or only
   the skill guidance? I flagged it as I read it; the fix-up cycle or consensus should rule.
2. D4's accepted cost makes release ordering two-file (bump `compatibility.json` together with
   `package.json`, and only then regenerate manifests — while drifted, `manifest:addons` now fails
   too). `RELEASING.md` never mentions `compatibility.json`/`skillVersion` (`PRIMARY`: grep exit
   1). FINAL neither required nor forbade documenting this — worth one line in `RELEASING.md` so
   the next release learns the coupling from the runbook rather than from a red `prepack`?
