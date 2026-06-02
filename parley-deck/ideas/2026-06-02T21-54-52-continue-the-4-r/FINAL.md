---
idea: 2026-06-02T21-54-52-continue-the-4-r
status: final
author: claude
consensus-date: 2026-06-02
participants: [codex, claude, hermes]
---

## Final plan / specification

### Goal
Implement follow-ups 1, 2, 4 (additive) and surface item 3's infra blocker. Signed-off decisions in ./consensus.md (unanimous codex/claude/hermes).

### Scope (build order 4 -> 2 -> 1)
- Item 4: generic write-first preamble in round/review prompts + stdout-capture fallback in runAgent (strict `---` validation; record agent.stdout_fallback). agy spec note.
- Item 2: review/consensus.md `outstanding_agreed_fixes`(+`blocked`) contract + ReviewAgreedFixes helper; bounded (maxCycles 3) auto Phase 7/8 loop for implementation blocks; fail-closed in auto if field absent.
- Item 1: additive cursor `ready_blocks[]`/`active_blocks[]` (schema bump + default from current_block); advanceDAG returns the full ready set; auto launches up to `--max-active` (default 4); status lists ready/active/blocked/complete; atomic writes.
- Item 3: BLOCKED — no GitHub releases/.exe assets for 1.6-1.9. Escalate release-creation to user (inbox); no faked hashes, no placeholder upstream PR.

### Tests
Fake-agent tests: agy stdout fallback (writes/valid-md/narration); Phase 8 loop (stubbed outstanding -> fixup -> zero); diamond DAG parallel readiness + old-cursor round-trip. go build/vet + full suite green.

### Non-goals
No invented winget hashes; no upstream PR without real assets; no vendor lock; no removal of existing behavior.

### Verification
Phase 6-8 review by codex+hermes; ship as CLI minor release + skill note if protocol text changes.

## References
- Consensus: ./consensus.md
- Rounds: ./round-01/
