---
agent: hermes
idea: tui-protocol-visibility
round: 1
date: 2026-06-11
---

Phase derivation (cursor.go:87 Rebuild): cursor 7-state (PhaseRound/Consensus/Final/Impl/Review/Done/Blocked) maps to 00-prompt phases 0-8 via COOPERATION.md §4. Mapping holds for round-01→cross-review (cursor.CurrentRound + CrossReviewRounds from protocol.ReadFrontmatter) and blocked→reopen (driver.go:35 ActionReopened after BLOCK signoff). Scaffold/empty FINAL.md (cursor.go:106 finalScaffoldReason=="") correctly stays in consensus, not PhaseFinal. Fix-up cycles (review/consensus.md + IMPLEMENTATION.md status:fix-up-cycle-N) correctly land in PhaseReview (cursor.go:102 highestReviewRound>=1). Proof for step 7 vs 6 vs 8: review/consensus.md existence + signoff count vs review/round-NN/ dirs vs IMPLEMENTATION.md complete (cursor.go:100). Counter-proposal: add explicit cursor.ReviewRound field derived from highestReviewRound + signoff triage to avoid 6/7/8 ambiguity on partial review dirs. ACCEPT with this addition.

Delivery merge (events primary + disk fallback + ✓? on disagreement): on virtio-fs, cases are (1) half-written artifact (Stat succeeds but content truncated), (2) stale ENOENT on existing file (event present, Stat misses), (3) file present post-reattach but no agent.finished event. Rule is sound — events are the live source, disk only fallback, disagreement surfaces as ✓? for operator visibility. fsutil-style resilience (retry/keep-last) should wrap exactly the disk-fallback Stat/ReadFile paths in consensus missingRoundArtifacts (not the event path). ACCEPT.

Reconcile cadence (15s running/60s done + two-consecutive-agreement before regression): I/O budget (~1-2 ReadDir for highestRound/highestReviewRound + 3-4 Stat for FINAL/IMPLEMENTATION/review-consensus + <=2 ReadFile for frontmatter) is realistic and already matches cursor.Rebuild + consensus.Status. Design forgot the ReadFile of 00-prompt.md status: field (needed for Phase 0 kickoff vs round-01) and review/consensus.md frontmatter review-cycle: N. Add both; budget still fits. ACCEPT with those two extra reads gated behind the same 15s/60s timer.

consensus.Status (consensus.go:92) + frontmatter parsing (cursor.go:188 readFrontmatterField): edge cases handled — missing signoff heading returns TriageMalformed (consensus.go:27), RESERVATIONS treated as non-ACCEPT (TriageReserved), review=true reviewer set = participants minus implementer (consensus.go:313 rule). No import cycle with protocol/workspace.go. ACCEPT.

run.phase event payload ({phase, current_round, impl_status, reason}): sufficient for ribbon/status line. Driver crash between cursor.Save (cursor.go:52 atomic rename) and event append is covered by disk-reconcile safety net on next refresh (Rebuild authoritative). ACCEPT.

Participants chain (m.opts → run.created payload → 00-prompt.md frontmatter): m.opts wins for live run (authoritative at start), run.created payload next, 00-prompt.md frontmatter as final static fallback. Disagreement case (e.g. dynamic participant list) correctly resolved by opts precedence. ACCEPT.