---
idea: 2026-06-02T21-14-49-finish-the-6-rem
status: final
author: claude
consensus-date: 2026-06-02
participants: [codex, claude, hermes]
---

## Final plan / specification

### Goal

Finish the six remaining §12 pipeline items as additive code on the shipped 1.8.0 surface, preserving all backward-compat and safety invariants. The detailed, signed-off decisions are in `./consensus.md` (unanimous ✅ codex/claude/hermes); this FINAL records the build contract.

### Scope

Build order **3 → 5 → 4 → 2 → 1 → 6**:
1. `pipeline execute --json` (versioned schema; ledger written before print) + `references/EXECUTION_CONTRACT.md`.
2. WinGet manifests for cli 1.6.0/1.7.0/1.8.0 + skill 1.3.0 (templated; sha/URL placeholders).
3. `pipeline watch SLUG --signals FILE [--once]` with a `SignalSource` interface, persisted breach dedupe, low-risk-only auto-remediation.
4. Phase 5-8 runner: `runner.RunImplementation` + `runner.RunReviewRound`; machine-readable `agreed_fixes` detection.
5. `pipeline auto` per-kind dispatch (action=dry-run→stop at prod gate; implementation=runner; watcher=finalize+handoff) with explicit stop statuses; auto NEVER executes prod.
6. 6a per-block transport (`Block.Transport`), 6b decider-agent (`Manifest.Decider`, low-risk non-prod only), 6c DAG execution (cursor `completed_blocks[]`/`ready_blocks[]`, additive schema bump) + a §12 text note.

### Implementation details

Per `consensus.md` §Agreed decisions 1-6. Key invariants: additive only; linear manifests + existing decks/cursors byte-for-byte unchanged (round-trip tested); agents write markdown, driver/harness executes via MCP; production capabilities (`*.production`,`*.rollback`) non-bypassable regardless of flags/manifest/decider; vendor-neutral via interfaces (`SignalSource`, `Provider`, transport allow-list).

### Tests

Full unit coverage with fake agents (helper-process pattern). Required: execute --json golden/shape; watch same/changed/resolved-then-reopened breach; Phase 5-8 runner artifact creation + failed review + fix-up loop + zero-agreed-fixes parse; auto stops at action prod gate (no real execute) + drives impl; per-block transport override; decider rejects production; DAG ready-set + old-cursor round-trip. `go build`/`vet` + full suite green before review.

### Non-goals

No real production mutation fired by the CLI. No WinGet publication (assets/PR are external). No removal/rewrite of existing §12 behavior. No new vendor dependency.

### Verification

Phase 6-8: codex + hermes review `review/round-NN/`, `review/consensus.md` with machine-readable `agreed_fixes`, fix-up until zero. Ship as a CLI minor release (tag + Homebrew tap); §12 text note applied to canonical + dogfood COOPERATION.md.

## References

- Consensus: ./consensus.md
- Rounds: ./round-01/
