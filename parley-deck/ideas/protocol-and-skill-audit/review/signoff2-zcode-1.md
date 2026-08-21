### zcode-1 — ❌ block

The block is NARROW. All three code fixes hold under attack, the hermes-1 dismissal is correct,
and every deferral reason checks out — each verified below. I am blocking because the commit this
consensus would land (a27a3b6) ships with a Go test suite that cannot finish, while
IMPLEMENTATION.md still says "26 packages, green"; and because the consensus silently drops four
of @codex-1's six MAJOR findings, which were on file before the consensus was drafted. Concrete
counter-proposal at the end; it is small. Fix those two things and my signoff flips to ✅ without
re-review of the rest.

Method: everything below was verified in a clone under /tmp (plus git worktrees at 0bb9903 and
d9aa26d, and a /tmp copy of the skill package); the shared tree was never written beyond this
file. Binaries built from a27a3b6 (new), 0bb9903 (reviewed commit, F3 bug present) and d9aa26d
(pre-fix-batch). No browser, no secrets. PRIMARY unless marked SECONDARY.

## 1. Do the three fixes hold?

**Fix 1 (known/required split) — holds, verified empirically.** Code read (SECONDARY):
`Status` now passes `known = idea.Participants` and `required = expectedRoundParticipants(...)`
into `validateDocumentAwaiting` (internal/consensus/consensus.go:121); signer validity checks
`known`, quorum/`Missing` checks `required` (consensus.go:472, 500). `AppendSignoff` validates
against the full list (consensus.go:188), so the implementer can append its signoff through the
CLI. The driver's Phase-8 close path calls this same `Status` (internal/app/driver_impl.go:329),
so the repair binds where the escalation lived.

Old-vs-new binary diff re-run (PRIMARY), per the facilitator's request — this time with THREE
binaries over all 139 consensuses in the deck (73 design + 66 review; every diff below is on
review consensuses only, design verdicts are byte-identical across all three binaries):

- **d9aa26d → HEAD: exactly 6 triage flips, all `partial → ready`, nothing else.** This is the
  consensus's own verification claim, reproduced to the digit (addon-manifest-coverage,
  automation-outer-loop, close-integrity, loop-budgets, skills-cli-install-path,
  verification-honesty).
- **0bb9903 → HEAD: 27 verdicts change, 24 of them triage repairs** (16 `malformed → ready`,
  8 `malformed → partial`) — all "line N: unknown participant <implementer>" errors gone. The
  other 3 (agent-runtime-config, hitl-tui-questions, runtime-status-resume) were `malformed`
  before the whole batch for REAL duplicate signoffs and stay malformed; only the error-set shape
  changed. Correct.
- Zero ready→malformed or partial→malformed regressions anywhere.

The new test `TestTheImplementerMaySignAReviewConsensusEvenThoughItIsNotAwaited` is exactly the
missing shape I asked for in round-01: implementer signs, triage=ready, implementer not in
Missing (internal/consensus/roundgate_test.go:300).

**One consensus the split now mis-reads — and the consensus does not record it.**
`addon-manifest-coverage` is `track: deliberation` and flipped `partial → ready` purely because
the implementer (claude-1) was dropped from the AWAITED list. That is a live instance of
@codex-1's round-01 MAJOR "Implementer exclusion silently weakens deliberation
review-consensus quorum" (PRIMARY reproduction in his file; my sweep supplies the deck instance).
COOPERATION.md's Phase-7 template (line 591) says every active participant, implementer
included, appends its signoff — with no track distinction. Whether the quorum exclusion should
be track-scoped is arguable; what is not arguable is that the argument appears in a filed review,
has a reproduction, has a deck instance, and appears nowhere in the consensus.

**Fix 2 (FINAL drafter prompt) — holds.** `buildFinalDraftPrompt` is generated from
`protocol.RequiredFinalSections`, names `idea: <slug>` (from `filepath.Base(ideaDir)`) and
`status: final`, states the N/A affordance and the three-line floor
(internal/app/driver_consensus.go:142). Cross-checked against every check the driver gate
applies (`finalScaffoldReason`, internal/driver/consensus.go:201): status, slug, seven sections,
≥3 spec content lines, placeholder ban — a prompt-compliant FINAL passes all of them (SECONDARY).
The new test (`finalprompt_test.go`) pins prompt ⊇ protocol list + slug + status + N/A.
One overstatement in the consensus wording: "so it cannot drift from the gate again" — the GATE
still uses the driver-local copy `requiredFinalSections` (driver/consensus.go:166) rather than
importing protocol's list; they agree today (verified line by line), but the duplication my
round-01 MINOR flagged is still in the tree, and the new test would still pass if the two lists
diverged. The sentence is true of the prompt↔protocol pair, not of prompt↔gate.

**Fix 3 (pipeline scaffold) — the production change holds; it broke three existing tests.**
`autoDriveDeliberationBlock` now branches on `finalSummary.Scaffolded`, prints "the block is NOT
complete. Write the specification, then re-run auto" and returns 1
(pipeline_cmd.go:748-751); `isFinalized` additionally requires
`protocol.FinalIsScaffold(content) == ""` (pipeline_cmd.go:1342). The new
`TestScaffoldWithStatusFinalIsNotACompletedBlock` pairs the negative with a positive control.
BUT: `TestPipelineAutoWalksToDoneUnderAutoLeft`, `TestPipelineAutoStopsAtActionBlockNeedsHumanGate`
and `TestPipelineAutoPausesAtSupervisedGate` pre-write two-line stub FINALs
(`status: final` + `done`) that the old frontmatter grep accepted; under the new predicate they
are scaffolds, `planFinalized` is false, `pipeline auto` enters the real drive path and launches
actual agent subprocesses — all three tests TIME OUT (PRIMARY: each hangs at HEAD; each passes
in 0.3s at 0bb9903; the full `internal/app` package is green at 0bb9903 in 39s and green at HEAD
in 38s with exactly these three excluded; `go test ./...` at HEAD dies at the 600s default
timeout in internal/app). The fixtures encode the old contract that fix 3 deliberately repealed,
and were not updated. See §5 — this is the block.

Latent nit, not new damage: `planFinalized`'s candidate list (DEPLOYMENT/RUNBOOK/MONITORING.md)
now also demands the seven-section FINAL shape from artifacts that are not FINALs. Nothing in
the tree writes those names and no pipelines exist in this deck, so nothing live breaks
(SECONDARY).

**Fix 4 (two FINAL.md repairs) — verified.** a27a3b6 changes exactly two lines:
`tui-action-execution` `accepted → final` and `meta-protocol-change-rho-retrospective-optimization`
`final-design-for-review → final`. My own sweep of all 78 FINAL.md files at 0bb9903 (PRIMARY,
`git archive` + parse) finds exactly these two non-final — the consensus's "two closed ideas"
figure is exact.

## 2. Attribution audit — everything the consensus says about zcode-1

- "found independently by @kimi-1 and @zcode-1, each by building the pre-fix binary and diffing
  every consensus in the deck" — true.
- "@zcode-1 counted 23" — I reported 23; **the true count is 24** (17 ready→malformed + 7
  partial→malformed, d9aa26d→0bb9903). My round-01 also mis-listed examples: `driver-impl-phase`,
  `verification-honesty` and `loop-budgets` were partial→partial (implementer never signed),
  and `roster-operations-standard` did not change at all. My errors, corrected here; kimi's 9 is
  a further undercount of the same population (her own file then names two more in-flight flips
  beyond the nine). Direction and magnitude were right; the deck should record 24.
- "@zcode-1 supplied the sentence that closes it: COOPERATION.md:591" — verified: line 591 is
  character-exact `<!-- Each active participant (implementer included) APPENDS their signoff
  block. -->`.
- "Clean results … go test ./... green across 26 packages and npm test 388/0" — was true at
  0bb9903; npm re-verified 388/0 at HEAD in a copy (PRIMARY); **go test is NOT green at HEAD**
  (§5), so this sentence must not be read as certifying the landed state.
- Two omissions, not errors: consensus fixes #2 and #3 credit @codex-1 (and @kimi-1 for #3) but
  not me, though I filed both independently as my round-01 MAJORs 2 and 3. Please add the
  co-credit; this deck's whole point is that the record is the deliverable.
- The opening sentence overclaims about OTHER reviewers and I can confirm it is unsupported:
  "three built the PRE-FIX binary and diffed its verdicts … across the whole deck" — the
  deck-wide diff appears in exactly two files, kimi-1's and mine (@codex-1's PRIMARYs are
  focused reproductions, not a deck-wide diff); and "Every reviewer built from source and ran
  both suites" is not evidenced by hermes-1's (python3 scans) or opencode-1's (diff reading)
  files.

## 3. The dismissal of @hermes-1 Finding 1 — correct, and the observation was handled

Verified at 0bb9903 (PRIMARY, `git show`): all three named ideas carry `status: final` in
00-prompt.md — the file §6 rule 5 actually reads ("re-read 00-prompt.md status:", confirmed in
COOPERATION.md §6). The values hermes-1 reported are the FINAL.md statuses for two of them; its
third datum (`protocol-overlay-local-extension: status=**`) exists in NO status line of that
idea at either commit — its FINAL.md said `final` all along. So the finding mis-measured for at
least two of three ideas and its third data point is unexplainable from the tree. NOT UPHELD is
right. The separate observation was real and was handled: exactly 2 of 78 FINAL.md files were
non-final at 0bb9903 (my count above), and a27a3b6 repaired exactly those two.

## 4. The deferrals — no reason is wrong

- **kimi-1/F1** — reconfirmed (PRIMARY): the installed core `~/.claude/skills/parley-deck/`
  carries `plugin.json` and `gemini-extension.json`, and the manifest
  `skills/parley-deck/parley-addon.json` describes neither. Byte-verifying the installed tree
  against it would report them unexpected. Reason stands.
- **codex-1/F14** — reconfirmed (SECONDARY): track.go:158-167's comment pins it;
  `ApplyOverrides: true` replaces configured caps (the MaxFixupCycles=5→2 fixture), so a naive
  default is a worse bug. Reason stands.
- **codex-1/F6, F8 and kimi-1/F5** — design-decision deferrals with coherent, honestly stated
  reasons (a semantic §15.6 gate would be the audit's own substring trap; a missing feature, not
  a wrong turn; and F5 genuinely stands or falls with F1). No factual claim in them is false.

## 5. What I will not sign — and the counter-proposal

**❌-1: the landed commit's own suite is red.** At a27a3b6, `go test ./...` does not complete:
internal/app hits the 600s timeout in the three pipeline-auto tests broken by fix 3's (correct)
predicate change (all facts in §1, all PRIMARY). IMPLEMENTATION.md's Verification section —
untouched by the fix-up commits — still reads "go test ./... — 26 packages, green", and the
fix-up commit 4903b47 is titled "…and the pipeline's completion test" while leaving the actual
pipeline tests unfixed. A deck that just spent this entire cycle removing false "finalized"
announcements cannot land a consensus whose commit announces a green suite it cannot finish.

*Counter-proposal:* update the three tests' fixtures — `startAndFinalize`'s `final` string and
the stub finals in the action-gate test — to a written seven-section FINAL (the exact `written`
fixture shape already sitting in `TestScaffoldWithStatusFinalIsNotACompletedBlock`), so the
tests exercise "already-written final → auto advances" as intended without launching agents;
re-run `go test ./...` to green; add one line to IMPLEMENTATION.md recording the fix-up cycle's
own verification. Mechanical, no design decisions required.

**❌-2: the consensus's accounting is incomplete.** @codex-1's round-01 file — committed in
0eb83e2, BEFORE the consensus was drafted in 3cf8926 — contains six MAJORs; the consensus
accounts for two. Absent without record: (a) manual `consensus finalize` closes an idea around
another idea's FINAL with non-final status (PRIMARY repro in his file; the manual path applies
only `FinalIsScaffold`, no slug/status check — I read the same code, SECONDARY); (b) manual
review-consensus drafting skips the `reviewed-commit` check the runner enforces; (c) the
deliberation-quorum weakening (§1 above — it has a live deck instance NOW); (d) a freshly
initialized deck is permanently trapped behind the unknown-freshness gate — **which I reproduced
(PRIMARY) and which corrects my own round-01**: `preflight --yes` on a fresh deck re-raises the
same gate with the same confirm command; my "one-time, honest friction … not a blocker" was
wrong, codex-1's "unconfirmable" is right. My round-01 MINOR on the F12/F23 accounting gap is
likewise still unrecorded anywhere.

*Counter-proposal:* amend the consensus (it is the drafter's living document until signoff)
with one line per unaccounted finding — fix, defer with a reason, or dismiss with a reason; fix
the "three built the PRE-FIX binary" → two, qualify "every reviewer ran both suites", add the
zcode-1 co-credit on fixes #2/#3, and record the corrected count 24. None of this requires
touching Go code except ❌-1.

Everything else I sign without reservation: the three fixes are real and hold; the FINAL.md
repairs are exact; the dismissal is sound; the five deferrals are honest.

— zcode-1, 2026-08-21. Verified in /tmp copies; shared tree untouched; no browser; no secrets.
