---
agent: kimi-1
idea: skill-sync-cli-1-39
review-round: 2
date: 2026-08-06
reviewed-commit: b769ced807268fae538fe934be96021784eb0025
---

# Re-review — skill-sync-cli-1-39, round 2 (kimi-1)

## Verdict

**accept.** All five agreed fixes (AF-1 … AF-5) are landed as the consensus specified, the AF-1
contradiction is gone with no new one introduced, the D2 paragraph remains a verbatim match with
the AF-3 sentence appended, `npm test` and `npm run prepack` are green at the fix-up commit, and
the fix-up ships nothing the consensus did not agree.

## Scope

Read: `review/consensus.md`; my own `review/round-01/kimi-1.md`; `FINAL.md`; the updated
`IMPLEMENTATION.md`; the full diff of `b769ced` in `parley-deck-skill`; the changed regions of
`skills/parley-deck/SKILL.md` (lines 240–270, 350–399, 795–860); all of
`references/WORKED_EXAMPLES.md`; `CHANGELOG.md` (2.4.0 notes); `package.json:50–72`; and, as CLI
ground truth for the new paragraph's snake-case claim, `internal/config/runtime.go:104` and
`internal/config/runtime_test.go:31` in `parley-deck-cli`.

Ran: `npm test` and `npm run prepack` in the real skill repo at `b769ced`; a
whitespace-normalizing Python comparison script in `/tmp/kimi1-r2-verify.py` (D1 row, D2 paragraph
plus the AF-3 appended sentence, sentence-uniqueness and old-absolute-absence checks); `git log` /
`git status --porcelain` / `git show` (read-only) in the skill repo; content greps of `SKILL.md`
and `CHANGELOG.md`.

Not done: I did **not** desync anything this round — the fix-up touches no guard file (diff is
`CHANGELOG.md`, `SKILL.md`, `parley-addon.json` only), so the D4 guard proof from round 1 stands
unrevisited and nothing needed restoring. I did not read the other reviewers' round-02 files
(independence). I executed no git write command; both repositories were left byte-identical to how
I found them (`git status --porcelain` empty after my runs, `PRIMARY`).

§15.1 note: I authored the D1/D2 wording adopted in FINAL. The checks below are mechanical
comparisons of the *implementation's* current text against FINAL plus the consensus-agreed
modifications — claims first asserted by the implementer in `b769ced`, checked `PRIMARY`. I issue
no verdict on any FINAL-era claim I own. AF-5 was my own NIT; what I assess below is the
implementer's new `IMPLEMENTATION.md` paragraph, not my own finding.

## My round-1 findings

1. **[MINOR] CHANGELOG `hermes` incident allusion** (my round-1 MINOR-1, consensus AF-2):
   **resolved.** `CHANGELOG.md` now reads *"A separate list teaches a two-list launch model the CLI
   does not implement."* — mechanism without incident, exactly the agreed fix. `grep -i hermes`
   over the whole `CHANGELOG.md` returns zero matches. `PRIMARY` (diff hunk `@@ -25,8 +25,8`;
   grep). The replacement wording echoes FINAL's recorded-dissent closure ("a two-field shape
   teaches a two-list launch model the CLI does not implement", `FINAL.md:101`), so it is grounded
   in ratified text, not new narrative. `PRIMARY`.
2. **[NIT] `headless_args` → `headlessArgs` adaptation undeclared** (my round-1 NIT, consensus
   AF-5): **resolved.** `IMPLEMENTATION.md` gains an *"Adaptation (review AF-5, kimi-1)"* paragraph
   stating the migration rule was written snake_case in FINAL and applied as `headlessArgs` to
   match the JSON shape, "recorded here rather than left silent" — precisely the one-parenthetical
   fix I asked for. `PRIMARY` (`IMPLEMENTATION.md:53–56`).

## AF-1 (the MAJOR) — contradiction check

**The old contradiction is gone.** The unscoped absolute *"Nothing is appended to it afterwards."*
appears nowhere in `SKILL.md` (script output: `old absolute present anywhere : False`). `PRIMARY`.
The manual JSON section now reads (`SKILL.md:379`, `:381`):

- *"This shape is **manual-facilitator input** … (branch A of "Generic CLI Invocation Contract").
  The Parley CLI reads its own snake-case configuration instead, where `headless_args` is the
  complete argv template and nothing is appended to it — see branch B."*
- *"**There is no separate write-mode argument list.** The flag that lets an agent write its own
  artifact belongs **inside** `headlessArgs`. Model, thinking and profile flags remain separate
  fields and are appended by branch A at launch; the write-enabling flag is not one of them."*

So the "complete argv template / nothing appended" claim now survives only in explicitly
CLI-scoped sentences (`:379` naming branch B, and branch B itself at `:833`), while the manual
section keeps only what is true of both branches — the consensus fix verbatim. `PRIMARY`.

**No new contradiction.** Cross-checked the new paragraph against all three neighbours, `PRIMARY`
(file reads cited in Scope):

- *vs branch A* (`SKILL.md:815–825`): step 2 states the write-enabling flag belongs in
  `headlessArgs` with no separate write-mode list — near-identical to `:381`'s first two
  sentences. Step 3 appends model/thinking/profile flags "only when discovered or configured";
  `:381` says those flags "are appended by branch A at launch". I considered whether the dropped
  condition is a new contradiction and rule it is not: the sentence's job is to attribute the
  appending mechanism to branch A (contrasting with the write-enabling flag, which "is not one of
  them"), and it explicitly defers to branch A, where the condition lives. Step 4 (prompt
  delivery) is untouched and uncontradicted — the manual section no longer says anything is or is
  not appended after `headlessArgs`.
- *vs branch B* (`SKILL.md:827–835`): `:379`'s scoped restatement ("complete argv template …
  nothing is appended … see branch B") matches branch B's "complete argv template … launched
  as-is" and "Nothing is appended afterwards". The claim "The Parley CLI reads its own snake-case
  configuration instead" is accurate: the CLI declares `HeadlessArgs []string
  \`toml:"headless_args"\`` at `internal/config/runtime.go:104` (used in
  `internal/config/runtime_test.go:31`), and per FINAL's verification record the CLI does not read
  `meta/headless-agents.local.json` — "its own … instead" is the correct boundary. `PRIMARY`
  (runtime.go read this round; the FINAL record cited as `SECONDARY` — ratified deck text I did
  not re-derive).
- *vs `WORKED_EXAMPLES.md`*: line 44 keeps the flag-inside-`headlessArgs` note and the migration
  rule and carries no "nothing appended" absolute; the example JSON (`:30–41`) keeps model and
  thinking as separate fields, matching `:381`'s "remain separate fields". Consistent. `PRIMARY`.

## AF-3 — D2 paragraph still verbatim, with the sentence appended

`SKILL.md:252` now ends *"… treat its autonomous bit as unset (fail-closed) rather than escalating
to a full-filesystem bypass. A vendor flag change is a config edit, not a skill revision."* and the
standalone line formerly at `:254` is gone. `PRIMARY` (diff hunk `@@ -249,9 +249,7`; grep counts
the sentence exactly once in the file). The whitespace-normalizing comparison of FINAL §D2's
adopted blockquote **plus the appended sentence** against `:252` returns `D2+appended MATCH :
True`; the appended words are byte-identical to the removed standalone line, so this is hermes-1's
option (a), drop-in, as agreed. `PRIMARY` (script `/tmp/kimi1-r2-verify.py`). The dismissed NIT's
premise is untouched: the paragraph remains one unrewrapped line, and the paragraph minus the
appended final sentence is still the verbatim FINAL text. `PRIMARY`. D1's `opencode` row also
still matches verbatim (`D1 MATCH : True`). `PRIMARY`.

## AF-4 — IMPLEMENTATION.md narrative corrected

The new *"Correction (review AF-4, hermes-1)"* paragraph states `package.json:60` is `node --test
&& node scripts/run-python-tests.js && node scripts/build-addon-manifest.js --check`, so one
`npm test` run exercises the guard through both callers while `prepack` exercises only the script
caller. Verified against the file: `package.json:60` and `:66` contain exactly those commands.
`PRIMARY`. This also matches what I observed empirically this round — the `npm test` output ends
with the six `ok` manifest lines from the `--check` invocation, including the D4 test by name
earlier in the run. `PRIMARY`.

## Did the fix-up break anything? — no

At `b769ced` (confirmed directly on top of `661af98`: `git log --oneline` shows `b769ced` then
`661af98`; `PRIMARY`):

- `npm test` → **386 node tests pass, 0 fail**, including `✔ compatibility.json skillVersion
  tracks package.json version`; 54 python tests across 7 files OK; all six add-on manifests `ok`,
  parley-deck aggregate `sha256:73e5d9258ed8d88ec95c9e63c25041af117e63de71e233bf0a9eaa2cc00ad37e`.
  `PRIMARY` (ran it; output quoted).
- `npm run prepack` → exit 0, all six manifests `ok`, same parley-deck aggregate. `PRIMARY`.
- The regenerated aggregate in the `b769ced` diff of `parley-addon.json` is
  `sha256:73e5d925…ad37e` — identical to what both gates recompute from the edited payload. The
  manifest regen is therefore exactly the mechanical consequence of the SKILL.md edit, and the D4
  guard (untouched by this commit) holds the tree in sync. `PRIMARY` (diff vs command output).

## Anything new shipped that consensus did not agree? — no

The skill-repo diff touches exactly three files: `CHANGELOG.md` (AF-2), `SKILL.md` (AF-1 and AF-3
hunks only), `parley-addon.json` (hash regen, required by D4's guard and already declared in
`IMPLEMENTATION.md`'s file table). The deck-side `IMPLEMENTATION.md` gains exactly the AF-4 and
AF-5 paragraphs. `PRIMARY` (full diff read; `IMPLEMENTATION.md` read). New wording inside the AF-1
rewrite — the branch-B pointer sentence and the "model/thinking/profile flags … appended by branch
A" sentence — is the scoped replacement the AF-1 fix describes ("label the shape as
manual-facilitator input … 'complete argv template / nothing appended' belongs to branch B's
snake-case `headless_args` alone"), and both sentences are verified consistent above. No
D6-excluded content (promotion history, probe outputs, ACP availability, kimi exit-1 message, CLI
struct names) appears anywhere in the fix-up. `PRIMARY`.

## Open questions

None from this round. My round-1 open question 1 (D6's scope over release notes) was resolved by
the consensus adopting AF-2. Question 2 (one line in `RELEASING.md` about the
`compatibility.json`/`package.json` coupling) was not part of the agreed fixes and remains a
non-blocking suggestion for a future change; I do not hold the verdict on it.
