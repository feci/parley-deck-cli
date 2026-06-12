---
idea: runner-hardening-kindly
cycle: 3
drafted-by: claude
date: 2026-06-12
reviewed-commit: cec1d62
outstanding_agreed_fixes: 0
---

<!-- Cycle 1 (reviewed-commit 8a5d4c7) follows; cycle 2 is appended below the
     cycle-1 signoffs. -->

## Agreed fixes

From codex (review/round-01/codex.md):

1. **[MAJOR] ACP activity + attempt contract (D1/D3).** (a) Mark ACP protocol
   milestones (`agent.acp.initialized`, `agent.acp.session_opened`,
   `agent.acp.prompt_completed`) as activity, not only SessionUpdate — a live
   init/session sequence must never classify as no_first_output; (b) thread
   attempt_id through the ACP marker (runID:agentID:attemptID) and the ACP
   `agent.started` payload; (c) implement the same retry-once-for-
   no_first_output loop for ACP attempts (attempt loop lifts to runAgent's ACP
   branch with attemptID threaded into runACPAgent).
2. **[MAJOR] Snapshot retention on move-back failure (D9).** Cleanup becomes
   conditional: a failed MoveArtifactBack retains the snapshot directory
   (marker removed so the stale sweep never erases the recovery artifact);
   terminal events report the LIVE canonical artifact path with the snapshot
   path carried only as recovery metadata in
   review.snapshot_artifact_move_failed.
3. **[MAJOR] Move-aside safety (D3).** The attempt-1 invalid artifact moves to
   a destination that is guaranteed not to exist (unique suffix when
   `.attempt-1.invalid` is taken); a rename failure removes the invalid
   artifact (it was created by this attempt — the preexisted guard already
   excludes foreign files) instead of leaving it on the canonical path. Tests
   for the pre-existing-destination and rename-failure cases.
4. **[MAJOR] Phase 8 fix-up joins the hardened path.** RunFixup switches from
   CommandFor+cmd.Run to the supervised exec path (process group + procctl
   marker + cleanParticipantEnv + counting writers + waitSupervised + watchdog
   events). No retry for fix-ups (a code-mutating phase is not safely
   re-runnable on a watchdog kill); terminal events remain
   agent.fixup_finished/failed with the classified payload.
5. **[MINOR] failEarly classification (D5).** failEarly routes through the
   classified payload (failure_class + recovery_hint from classifyFailure over
   the setup error text; bounded tails when logs exist).
6. **[MINOR→merged] Consult provenance (codex m7 + agy MAJOR).** The consult
   frontmatter gains `session_id` (written even when empty — one-shot CLIs
   expose none today; agy's `consult:<agent>` marker value was judged not to
   be a session id) and `timeout_ms` records the EFFECTIVE timeout RunConsult
   used, not the raw flag.
7. **[TEST] Lock the hint table.** Extend TestClassifyFailure to assert exact
   class/hint pairs for the implemented table (see Dismissed #1 for why the
   implemented strings, not the round-01 strings, are the contract).

From agy: the session_id finding is merged into fix 6.

## Deferred follow-ups

- ACP-path claude marker shedding via acp.MergedEnv (recorded deviation;
  codex concurred for shedding specifically).
- waitSupervised 1s tick granularity: documented as accepted (hermes MINOR);
  a doc sentence lands with fix 4's comments.

## Dismissed findings

1. **codex [MINOR] "hints must match round-01 strings verbatim."** Dismissed
   per Phase 7 triage: agy — the taxonomy's author — reviewed the implemented
   paraphrases in this same round and explicitly accepted them as improvements
   ("no code change is requested"). The author's acceptance supersedes
   verbatim fidelity; the UX contract is locked going forward by the exact
   class/hint test (agreed fix 7). codex may contest in signoff.

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: claude — 2026-06-12
Status: ✅ ACCEPT
Notes: Seven agreed fixes for fix-up cycle 1; the strings dispute resolved by author-acceptance + test-locking.

### Signoff: codex — 2026-06-12
Status: ✅ ACCEPT
Notes: My round-01 findings are faithfully triaged; I accept the hint-string dismissal because author-accepted paraphrases plus exact-pair tests lock the contract.

### Signoff: hermes — 2026-06-12
Status: ✅ ACCEPT
Notes: Round-01 triage faithful; verbatim-hint dismissal accepted (author acceptance + test lock supersedes).

### Signoff: agy — 2026-06-12
Status: ✅ ACCEPT
Notes: Round-01 findings faithfully triaged; I concur with the hint-string dismissal as the taxonomy author.

## Cycle 2 (review/round-02 → fix-up cycle 2)

Reviewed commit: 6e20f1e. Verdicts: codex ACCEPT-WITH-FIXES, agy ACCEPT,
hermes ACCEPT. codex verified fixes 1, 4, 5, 6, 7 and flagged remainders of
fixes 2 and 3; agy independently surfaced the same finishACP gap as a NIT.
Both remainders are agreed:

### Agreed fixes (cycle 2)

1. **[MAJOR remainder] ACP live-path on publish failure (codex fix-2
   NOT-FIXED, merged with agy's NIT).** `finishACP` assigns
   `result.OutputPath` only on the publish success branch; an ACP review
   snapshot publish failure would emit a terminal event carrying the snapshot
   path. Mirror `finalizeExecResult`: whenever `publishArtifact` returns a
   non-empty live path, the terminal event reports it — even when the
   move-back failed.
2. **[TEST remainder] rename-failure case (codex fix-3 NOT-FIXED).**
   `TestMoveAsideInvalidArtifact` covers only the pre-existing-destination
   case; the agreed rename-failure case (rename fails → the invalid artifact
   is removed from the canonical path) is missing. Add it (forced
   deterministically via a destination basename beyond NAME_MAX, which fails
   the rename while the source stays removable).

### Dismissed findings (cycle 2)

- None.

### Dispositions (cycle 2)

- `TestDurableKillEndToEndRealProcess` under the codex seatbelt sandbox:
  re-confirmed by all three reviewers (codex reproduced the identical
  kern.boottime/boot-id attribution failure); remains dismissed as an
  environment artifact.

### Signoffs (cycle 2)

<!-- Each agent APPENDS their cycle-2 signoff block below. Do NOT edit others' blocks. -->

### Signoff: claude — 2026-06-12 (cycle 2)
Status: ✅ ACCEPT
Notes: Two remainders agreed for fix-up cycle 2; both are narrow (one ACP branch mirror, one test case).

### Signoff: codex — 2026-06-12 (cycle 2)
Status: ✅ ACCEPT
Notes: Cycle-2 triage matches my round-02 findings, including both runner remainders and the sandbox-test dismissal.

### Signoff: hermes — 2026-06-12 (cycle 2)
Status: ✅ ACCEPT
Notes: Cycle-2 triage faithfully reflects my round-02 ACCEPT verdict and sandbox dismissal.

### Signoff: agy — 2026-06-12 (cycle 2)
Status: ✅ ACCEPT
Notes: Verified all cycle 2 fixes and concur with the triage and sandbox dismissal.

## Cycle 3 (review/round-03 → complete)

Reviewed commit: cec1d62. Verdicts: codex ACCEPT, agy ACCEPT, hermes ACCEPT.
All three reviewers verified both cycle-2 remainders (finishACP live-path
mirror; rename-failure test case) with zero new findings, and re-confirmed
the TestDurableKillEndToEndRealProcess seatbelt-sandbox dismissal. Zero
agreed fixes remain — the implementation closes as complete on cycle-3
signoff.

### Agreed fixes (cycle 3)

- None.

### Dismissed findings (cycle 3)

- None.

### Signoffs (cycle 3)

<!-- Each agent APPENDS their cycle-3 signoff block below. Do NOT edit others' blocks. -->

### Signoff: claude — 2026-06-12 (cycle 3)
Status: ✅ ACCEPT
Notes: Zero agreed fixes after two fix-up cycles; idea completes and ships as 1.24.0.

### Signoff: codex — 2026-06-12 (cycle 3)
Status: ✅ ACCEPT
Notes: My round-03 review is reflected accurately: both remainders verified, no new findings, and the sandbox dismissal remains accepted.

### Signoff: hermes — 2026-06-12 (cycle 3)
Status: ✅ ACCEPT
Notes: Cycle-3 summary matches my round-03 ACCEPT verdict with zero remaining issues.

### Signoff: agy — 2026-06-12 (cycle 3)
Status: ✅ ACCEPT
Notes: Both cycle-2 remainders verified and sandbox dismissal re-confirmed; I accept the final complete status.
