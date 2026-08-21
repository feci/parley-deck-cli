
### Signoff: kimi-1 — 2026-08-21
Status: ❌ BLOCK
Notes: Fixes 1, 2, 4 and the production half of fix 3 verified PRIMARY and hold — but fix-up cycle 1 (4903b47) changed isFinalized without updating the pipeline test fixtures that encoded the old frontmatter-only finality, so `go test ./...` is RED at HEAD (internal/app hangs to the 10-min timeout, launching a real `codex` process from tests). This deck does not sign a tree its own suite rejects.
Counter-proposal (required if ❌): In internal/app/pipeline_cmd_test.go make the four stub artifacts satisfy protocol.FinalIsScaffold (one shared writtenFinal() helper — all seven RequiredFinalSections headings + ≥3 spec content lines) at line 51, lines 104-105 and line 141, then re-run `go test ./...` to green and record the run in the fix-up note. Proven in an isolated copy: the four affected tests pass in 0.33s, full internal/app green in 36s.

## Evidence base

PRIMARY (I ran it): all builds and sweeps in an isolated copy at /tmp/pdv (rsync of HEAD a27a3b6
plus `git archive` extractions of 0bb99031 and a1926ae); the shared tree was not modified.
Binaries: parley-base (a1926ae, pre-F2/F3), parley-old (0bb99031, the reviewed commit),
parley-new (a27a3b6, HEAD). Sweep: `consensus status --review --json` over all 87 idea dirs with
each binary. `go test` on both extractions and the HEAD copy. Fixture patch applied and tested in
the copy only.
SECONDARY (I read it): the three fix commits' code, the two repaired FINAL.md files, `git show`
at 0bb99031 for the dismissal, the skill installer source for the F1/F5 reason.

## 1. Do the three fixes hold?

**Fix 1 (known/required split) — holds.** `Status(..., review=true)` passes the full participant
list as KNOWN and the implementer-reduced list as AWAITED (`internal/consensus/consensus.go:121-122`);
identity is validated against KNOWN, missing/triage against AWAITED
(`validateDocumentAwaiting`, consensus.go:452-518). Fail-closed paths intact: unresolvable
implementer → full list awaited (consensus.go:654-656); reduction to zero participants → full
list (663-665); non-review path unchanged (446-448). An implementer-only signoff still reads
partial — ready cannot be gamed. An implementer 🟡/❌ still requires notes/counter-proposal.

Re-ran the old-vs-new diff, wider than before (three binaries, whole deck, current tree):

- old(0bb99031) → new(HEAD): 27 verdict diffs, zero regressions. 24 are `malformed →
  ready/partial` — exactly the F2/F3 damage undone, including the two in-flight ideas from my
  round-01 finding: `rho-retro-tooling` malformed→**ready**, `build-companion-skills`
  malformed→**partial** (awaiting only antigravity-1, matching my filed note). 3 are
  `malformed → malformed` where the error text now names the REAL problem (e.g.
  agent-runtime-config: the masked `duplicate signoff for codex` replaces the old
  `unknown participant codex`) — strictly better diagnostics, triage unchanged.
- base(a1926ae) → new(HEAD): exactly **6 triage flips, all `partial → ready`**, plus 23
  same-triage corrections that only drop the implementer from the awaited list (the agreed F3
  semantics). The consensus claim "6 flips, all partial → ready, zero regressions" is **exact**.

No consensus in this deck mis-reads under the split. The 22 ideas whose status errors (no review
consensus) error identically under all three binaries.

**Fix 2 (drafter prompt) — holds.** `buildFinalDraftPrompt` generates its section list from
`protocol.RequiredFinalSections` and names the idea slug for the frontmatter
(`internal/app/driver_consensus.go:142-167`); `TestFinalDraftPromptDescribesWhatTheGateRequires`
passes at HEAD (PRIMARY).

**Fix 3 (scaffold finality) — production code holds; the commit broke the suite.** `isFinalized`
now requires `status: final` AND `protocol.FinalIsScaffold(content) == ""`
(`internal/app/pipeline_cmd.go:1342-1354`); `autoDriveDeliberationBlock` reports "FINAL.md
scaffold written; the block is NOT complete" and exits non-zero (pipeline_cmd.go:743-751); the
new `TestScaffoldWithStatusFinalIsNotACompletedBlock` passes. That is precisely the repair my
round-01 F5 finding asked for.

**But** the commit did not update the fixtures that encoded the old semantics. PRIMARY, HEAD copy:

- `startAndFinalize` (pipeline_cmd_test.go:51) writes `---\nstatus: final\n---\n\ndone\n` — its
  own comment says it exists "so the auto-loop control flow can be exercised without launching
  any agent". Under the new isFinalized that stub is a scaffold, so the pipeline enters
  `autoDriveDeliberationBlock` and launches a real agent: `TestPipelineAutoWalksToDoneUnderAutoLeft`
  and `TestPipelineAutoStopsAtActionBlockNeedsHumanGate` (stub at line 141) HANG until the
  10-minute `go test` timeout, with goroutine dumps inside `runner.execAgentProcess` — on this
  machine the tests exec the real `/opt/homebrew/bin/codex`. `TestPipelineAutoPausesAtSupervisedGate`
  shares the fixture. Full-suite run at HEAD: `FAIL parley-deck-cli/internal/app 602.236s`
  (25 other packages ok).
- With agent CLIs removed from PATH the same test fails fast instead (exit 1, "participant
  selection failed"). There is no environment in which these tests pass at HEAD.
- Control: the identical suite on the 0bb99031 extraction, same machine, same PATH: green in
  37.5s. The breakage is attributable to the fix-up commits, not to my environment.
- `TestActionBlockCompleteNeedsSucceededEffectNotJustPlan` still passes, but its first assertion
  no longer constructs a finalized plan (the stub DEPLOYMENT.md/FINAL.md at lines 104-105 are
  scaffolds now), so "a finalized plan alone must not complete an action block" is vacuously
  satisfied — a weakened test, the class opencode-1 certified absent at the reviewed commit. The
  weakening arrived with the fix-up, after that certification.

This contradicts "go test ./... green after every commit" for the fix-up cycle (the reviewed
commit 0bb99031 is genuinely green — verified). Block stands until the fixtures move.

## 2. Attribution audit (sentences naming kimi-1)

- "found independently by @kimi-1 and @zcode-1 … @kimi-1 counted 9 flips to malformed" — that IS
  what I filed; the attribution is faithful. **Correction against myself, recorded for the
  record:** re-measuring now (both on the current tree and on a pristine 0bb99031 extraction)
  the F2/F3 change flipped **24** review consensuses into `malformed`, not 9; zcode-1's 23 was
  nearly exact, my 9 was an undercount (all 24 restore correctly at HEAD, so the fix verification
  is unaffected).
- Fix 3 "MAJOR, @codex-1 and @kimi-1" — accurate; matches my round-01 F5 finding.
- "three built the PRE-FIX binary and diffed its verdicts across the whole deck" — accurate for me.
- "Every reviewer built from source and ran both suites" — **not evidenced by my round-01 file**:
  it records building both binaries and the full-deck status sweep, not go/npm suite runs. Minor
  over-generalization; suggest "all five filed; three diffed binaries", which is fully accurate.
- "@kimi-1/F1, F5 remain unfixed with stated reasons … no reviewer challenged any of those
  reasons" — accurate.
- The zcode-1 quote of COOPERATION.md:591 ("Each active participant (implementer included)
  APPENDS their signoff block") — verified verbatim (PRIMARY, grep).

## 3. The dismissal of hermes-1 Finding 1

Correct (PRIMARY, `git show 0bb99031:…/00-prompt.md`): all three ideas carried `status: final` in
`00-prompt.md` at the reviewed commit; the values hermes-1 reported (`accepted`, `**`,
`final-design-for-review`) are the FINAL.md `status:` — a different file from both the repair
target and the one §6 rule 5 reads. The finding as filed measured the wrong file; dismissed
rightly.

The observation it produced was handled: commit a27a3b6 repairs exactly two FINAL.md files with
one-line `status:` changes to `final` (tui-action-execution: `accepted`→`final`;
meta-protocol-change-rho-retrospective-optimization: `final-design-for-review`→`final`), which is
what `finalScaffoldReason` requires (driver/consensus.go:209). The third name,
protocol-overlay-local-extension, already had `status: final` in FINAL.md at 0bb99031 — the `**`
was an artifact of hermes-1's own parser; nothing to repair. Housekeeping note: the consensus's
"Deferred follow-ups" still lists these two FINAL.md files as open because the consensus
(3cf8926) predates the repair (a27a3b6); a reader should be pointed at the repair, not sent to
redo it.

## 4. The deferrals — reasons checked

- **codex-1/F14** — sound. Recorded in `PolicyFor`'s own comment (internal/track/track.go:156-165);
  the override-stomps-config mechanism (MaxFixupCycles 5→2) is the real semantics of
  `ApplyOverrides`, and a per-knob default policy is a design decision, not a fix.
- **codex-1/F6** — sound. §15.6 (COOPERATION.md:1346-1354) demands a steelmanned adversarial
  alternative for judgment calls — a semantic property of prose; a mechanical gate would be the
  substring-matching trap this audit keeps finding.
- **codex-1/F8** — sound. No collapsed-final branch exists in the planner or auto-driver
  (grep: fast-track handling is only the round short-circuit and model-diversity gate); it is a
  missing feature of the fast route, not a wrong turn.
- **kimi-1/F1** — verified against the installer source (SECONDARY, lib/installer.js:147-160 of
  parley-deck-skill): the installed core is staged from the core subdir (SKILL.md, agents/,
  references/) PLUS package-root `plugin.json`, `gemini-extension.json` (required) and
  `README.md`, `LICENSE` (optional), while `parley-addon.json` (6 file entries) describes only
  the source subdir — so byte-verifying the installed tree against that manifest reports the four
  package-root files as unexpected. "Needs a payload manifest describing the installed shape" is
  the correct reason; reverting the half-working attempt was right.
- **kimi-1/F5** — the fold into F1 is coherent: the core marker cannot honestly carry a manifest
  anchor until that manifest exists.

No reason is wrong.

## 5. What I will not sign

- Any record implying HEAD is suite-green. It is not; the counter-proposal above is the whole
  gap, it is proven in a copy, and it should clear in one small commit plus one `go test ./...`
  run.
- The sentence "Every reviewer built from source and ran both suites" as a claim about me (see
  §2). Also noted: "npm test 388/0" refers to the sibling skill repo's suite
  (parley-deck-skill/package.json); no fix-up commit touches that repo, so its state is invariant
  here — I did not re-run it (SECONDARY to this signoff).
