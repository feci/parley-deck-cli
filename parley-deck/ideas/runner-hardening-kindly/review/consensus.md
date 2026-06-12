---
idea: runner-hardening-kindly
cycle: 1
drafted-by: claude
date: 2026-06-12
reviewed-commit: 8a5d4c7
outstanding_agreed_fixes: 7
---

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
