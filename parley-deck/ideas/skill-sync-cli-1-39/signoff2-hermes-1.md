### hermes-1 — revision 2

Verdict: accept

I signed off on revision 1 with no reservations. Revision 2 rewrites VC-1, corrects
decision 4's `prepack` clause, and discloses a measurement-tool defect. Below I assess
each change, re-verify what I can reach, and record one thing my own revision-1 signoff
got wrong that revision 2 fixes.

----------------------------------------------------------------

## codex-1's three required changes — all met

codex-1's revision-1 block (consensus.md:378-391) set three conditions. I assess each
against the revised text.

### (i) claim (b) marked UNVERIFIED and removed from the closure — MET

`PRIMARY` The consensus now states, at lines 226-230: "The 23-deck aggregate therefore
stands `UNVERIFIED` by a non-owner and is used below only as illustration, never as
load-bearing evidence — the condition codex-1 set." The closure (lines 293-299) is
explicit: "VC-1 closes on the design argument, not on the workspace counts … The
aggregate in (b) illustrates the teaching but is `UNVERIFIED` and carries nothing."

The claim-(b) text (lines 190-224) now carries the exact commands (`find … -print0 |
xargs -0 grep -l`), the inputs (the workspace root path), the output (the 23/12/10/9
table at lines 210-216), and locators (the named 10-deck list at lines 218-221, plus
the `igm-app` worked example at lines 222-224). This cures the §15.2 malformation
codex-1 identified: revision 1's `PRIMARY` label had no command/inputs/output/locators
and read as `RECALL`.

### (ii) claim (a) narrowed, three inferences deleted — MET

`PRIMARY` Claim (a) is now scoped (lines 153-156) to "no Go code in `parley-deck-cli`
declares a `writeModeArgs`/`write_mode_args` field or reads
`meta/headless-agents.local.json`." Lines 182-186 explicitly withdraw the three
overreaching inferences: "it does not show the field has 'no consumer', that deleting
it is a 'zero behaviour change', or that 'the protected cost does not exist'. Those
three inferences are withdrawn from revision 1. The file does have a consumer — a
facilitator hand-assembling a command, exactly as codex-1's position held."

I independently re-verified the scoped claim (a) for this signoff:
- `PRIMARY` I ran a repository-wide search (`search_files`, ripgrep-backed) for
  `headless-agents|headless_agents|writeModeArgs|write_mode_args|WriteModeArgs` across
  all `*.go` files under `parley-deck-cli` — zero hits.
- `PRIMARY` I read `internal/config/runtime.go:97-132` (`agentOverride`): it has
  `HeadlessArgs` and no write-mode field. I read `internal/config/runtime.go:134-154`
  (`configLayers`): it enumerates exactly `~/.parley/agents.toml` (line 141),
  `parley-deck/agents.toml` (144), `parley-deck/agents.local.toml` (145), and
  `$PARLEY_HEADLESS_AGENT_CONFIG` (147-151) — no JSON path.

This reproduces my revision-1 check and the drafter's `find`-based re-derivation. I am
a non-owner of claim (a) (claude-1 owns it, `round-01/claude-1.md` item 4), so my
`CONFIRMED` is admissible under §15.1.

### (iii) the legacy migration rule added — MET

`PRIMARY` Lines 306-314 adopt codex-1's condition: "when an existing
`headless-agents.local.json` contains `writeModeArgs`, merge its arguments into that
agent's `headlessArgs` and remove the field." The rule is stated as a skill instruction
(no deck is edited by this idea), and the follow-up is deferred to `FINAL.md`. This is
the rule codex-1's block required ("add an explicit legacy migration rule: when a
manual JSON config contains it, merge its arguments into `headlessArgs` before launch
and update that config").

----------------------------------------------------------------

## kimi-1's reservation R1 — correctly recorded

`PRIMARY` R1 (kimi-1's block, consensus.md:450-462) said the draft's claim that the
assertion runs "under both `npm test` and the existing `prepack` lifecycle" was
imprecise, because `prepack` runs only `build-addon-manifest.js --check` and does not
invoke `node --test`, so a test-file assertion leaves `npm publish` ungated.

`PRIMARY` Revision 2 corrects this at lines 83-96: "`package.json:60` defines `test` as
`node --test && …`; `package.json:66` defines `prepack` as
`node scripts/build-addon-manifest.js --check` alone. `prepack` does not invoke
`node --test` … Implementation MUST close the `prepack` half explicitly, by one of
(a) extend the `prepack` command to run the version check, or (b) put the equality
check inside `scripts/build-addon-manifest.js`, which `prepack` already runs."

I independently verified the prepack facts against the skill source repo
(`/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/parley-deck-skill`), at HEAD
(commit `76b31bc`):
- `PRIMARY` `package.json:60`: `"test": "node --test && node scripts/run-python-tests.js
  && node scripts/build-addon-manifest.js --check"`.
- `PRIMARY` `package.json:66`: `"prepack": "node scripts/build-addon-manifest.js --check"`.
  `prepack` does not invoke `node --test`.

Ownership note (§15.1): I first raised the prepack gap in `round-02/hermes-1.md @codex-1`,
so I co-own the underlying claim and issue no verdict on it here. The admissible
non-owner verdict is the drafter's (`PRIMARY`, consensus.md:88-89, verified against
`package.json:60,66`). My reads above are evidence, not a self-verification.

----------------------------------------------------------------

## The defective measurement tool — assessment

The drafter disclosed that revision-1 numbers came through a shell shim: `grep` in the
facilitator's environment is aliased to `ugrep --ignore-files`, which honours
`.gitignore` and silently omits ignored paths with exit status 1 (consensus.md:232-284,
drafter position change 10 at line 333). Revision 1's first causal explanation —
"`grep -r` silently under-reports on this filesystem" — was the drafter's second false
causal explanation in this idea and was corrected mid-signoff to the shell-alias cause.

### Is the corrected measurement adequately evidenced?

Yes, where it is load-bearing.

- Claim (a) — the load-bearing factual claim for the VC-1 closure — does NOT rest on
  the defective tool. `PRIMARY` I re-verified it with my own search (zero hits). The
  consensus records independent non-owner confirmations from codex-1 (`rg`, exit 1, no
  output) and kimi-1 (`rg` over all `*.go`, zero hits) at lines 170-173. The drafter's
  own re-derivation uses `find … -print0 | xargs -0`, which does not consult ignore
  files (line 161-163). So (a) has three independent non-owner `PRIMARY` confirmations
  plus a `find`-based owner run — well-evidenced, not instrument-correlated.

- Claim (b) aggregate (12/10/9) — the drafter re-derived it with `find … -print0 |
  xargs -0` and a `while read` loop (lines 192-208). It is marked `UNVERIFIED` by a
  non-owner and demoted to illustration (lines 226-230, 293-299). I cannot reach the 22
  decks outside this repository and do not verify the aggregate. `PRIMARY` I did verify
  the one in-repo deck: `parley-deck/meta/headless-agents.local.json` sets
  `writeModeArgs` on all four agents — claude-1 (lines 28-31, non-empty), codex-1
  (78-81, non-empty), hermes-1 (109-112, non-empty), kimi-1 (154, empty array) — and
  for claude-1/codex-1/hermes-1 the flags are duplicated in `headlessArgs` (lines
  18-25, 68-75, 103-106). This matches the consensus description at lines 218-221
  ("this repo's own deck, the only non-exclusive one") and kimi-1's verification. Since
  the aggregate carries nothing load-bearing, its UNVERIFIED status is adequate.

### Does anything else in the draft still rest on the defective tool?

I checked every measurement claim in the consensus:

- `references/COOPERATION.md` currency (line 119-120): verified by `git diff
  v1.38.0..v1.39.0`, not by `grep`. Not affected. (Facilitator-owned; I did not
  re-verify, as in my revision-1 signoff — no position of mine depends on it.)
- "nothing under `test/` or `scripts/` reads `compatibility.json`'s `skillVersion`
  today" (lines 104-105): `PRIMARY` I verified this at HEAD with `/usr/bin/grep -rn
  "skillVersion" . --include='*.js' --include='*.json'` in the skill repo — zero hits
  in `test/` or `scripts/` at HEAD. This could have been a false negative if the
  drafter used the ugrep shim, but my non-defective-tool check confirms it. Not
  affected.
- The `runner.go:1094-1108` reference (item 3) and `package.json:60,66` reads
  (decision 4): direct file reads, not `grep` over a tree. Not affected.
- The `parley-addon.json` sha256-pin claim (lines 75-77): verified by reading the
  manifest-check code, not by tree grep. Not affected.

Nothing else in the draft rests on the defective tool.

### A note on the mid-signoff correction

`PRIMARY` The consensus discloses at lines 601-609 that revision 2 was corrected while
signoff requests were already out (the filesystem-claim text was replaced with the
ugrep-shell-alias explanation). The corrected text reaches the same numbers (12/10/9)
by the same `find`-based commands; no decision moves. I am signing off on the corrected
text. My revision-1 signoff predates revision 2 and never quoted the filesystem claim,
so I have nothing to revise on this basis.

----------------------------------------------------------------

## What my own revision-1 signoff got wrong

My revision-1 signoff (consensus.md:525-527) argued: "Claim (a) shows there is no
schema — no validator, no loader, no consumer in the CLI. Deleting the field from the
documentation changes the runtime behaviour of zero decks. The cost codex-1 was
protecting against does not exist; the argument's factual premise is removed, not
sidestepped. That is a direct defeat."

That is the same overreach revision 2 withdraws. Claim (a) establishes "no Go loader
and no Go field" — it does not establish "no consumer" (the manual facilitator is a
consumer) or "zero behaviour change" or "the protected cost does not exist." I made the
inference the drafter made, and codex-1's block correctly identified it as overreaching.
`RECALL` I did not recognise this in my revision-1 signoff because I shared the
majority's delete position and read (a) as stronger than its scope warranted.

Revision 2 closes VC-1 on the design argument alone: a two-field shape teaches a
two-list launch model the CLI does not implement, and (a) establishes no Go loader
exists whose compatibility deletion could break. That is the argument I was actually
making in rounds 1 and 2 (my round-1 position rested on `agentOverride` having no
`writeModeArgs` field, `PRIMARY`), without the overreach. I accept the revision-2
closure as the correct and cleaner framing.

----------------------------------------------------------------

## The six edits

Edits 1, 2, 3, and 6 are unchanged from revision 1. I assessed them in my revision-1
signoff (consensus.md:533-540) and have nothing to add.

Edit 4 (`compatibility.json` bump plus assertion): the `prepack` correction (lines
83-96) is the change. I verified the underlying facts at HEAD (`PRIMARY`, above). The
consensus correctly requires implementation to close the `prepack` half by option (a)
or (b). `PRIMARY` I confirmed option (b) is viable: `prepack` (package.json:66) already
runs `build-addon-manifest.js --check` at HEAD, so placing the equality check inside
that script gates `npm publish` without a new script.

Observation for Phase 5 (non-blocking): `PRIMARY` the skill source repo's working tree
already contains uncommitted changes implementing D4 option (b) — a `versionSyncProblem`
function at `scripts/build-addon-manifest.js:85-96` and a test at
`test/manifest-coverage.test.js:489-500` (which references "Idea `skill-sync-cli-1-39`,
decision D4"). At HEAD these do not exist. This is an in-progress implementation
consistent with the consensus direction; it does not affect this signoff, but the
implementer should be aware the working tree is ahead of HEAD.

Edit 5 (`writeModeArgs` deleted): per VC-1 closure above. Agree.

----------------------------------------------------------------

## §15.5 drafter position changes

`PRIMARY` The table at lines 320-333 now records ten changes. Changes 1-5 match the
round files I read (I ratified 1-4 in my revision-1 signoff at consensus.md:542-544;
change 5 I ratified at lines 523-524). Changes 6-9 record the revision-1-to-revision-2
corrections codex-1's block and kimi-1's reservation forced. Change 10 records the
mid-signoff correction of the second false causal explanation. `PRIMARY` I verify
change 10 against the consensus body: lines 232-284 give the corrected ugrep-shell-alias
explanation, and line 333 records the withdrawal. The table preserves change 8's
original (now-superseded) causal claim ("under-reports on this filesystem") alongside
change 10's correction — this is §15.5 working as intended (file the change, do not
smooth it), not an inconsistency to fix.

## §15.6 correlated agreement

Unchanged from revision 1. `PRIMARY` The consensus at lines 340-360 records the
correlated-agreement caveat, the what-would-have-to-be-true condition, and the
one-family note. Round 1 was not unanimous (I opposed the guard), so no steelman round
was required. I have nothing to add to my revision-1 assessment
(consensus.md:546-548).

----------------------------------------------------------------

## Scope declared

Read in full for this signoff: `COOPERATION.md` §15 (lines 1176-1316);
`consensus.md` (all 610 lines, including the embedded revision-1 signoff blocks at
lines 362-598 and the revision-2 disclosure at 601-609).

Ran for this signoff (`PRIMARY`):
- Repository-wide search for
  `headless-agents|headless_agents|writeModeArgs|write_mode_args|WriteModeArgs` across
  all `*.go` files in `parley-deck-cli` — zero hits.
- Read `internal/config/runtime.go:97-154` (`agentOverride` struct and `configLayers`).
- Read `parley-deck/meta/headless-agents.local.json` (all 179 lines) to verify the
  in-repo deck against consensus lines 218-221.
- Located the skill source repo at
  `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/parley-deck-skill`; read
  `package.json:55-69` and verified `prepack`/`test` at HEAD via `git show HEAD:package.json`.
- Ran `/usr/bin/grep -rn "skillVersion"` across the skill repo to confirm "nothing
  reads `skillVersion` today" (consensus lines 104-105) at HEAD — zero hits in
  `test/`/`scripts/`.
- Read `scripts/build-addon-manifest.js:80-114` and `test/manifest-coverage.test.js:485-500`
  in the skill repo working tree (uncommitted changes implementing D4 option (b)).

Did not run: the skill test suite (`npm test`); the `parley` CLI; live vendor probes;
the `v1.38.0..v1.39.0` COOPERATION.md currency diff (facilitator-owned; no position of
mine depends on it).

Could not reach: the 22 `headless-agents.local.json` decks outside this repository
(claim (b) aggregate). I do not verify the 12/10/9 aggregate and do not guess at it.

I ran no git write commands and edited no file other than this signoff.
