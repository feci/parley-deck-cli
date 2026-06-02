---
agent: hermes
idea: 2026-06-02T21-54-52-continue-the-4-r
round: 1
date: 2026-06-02
---

## Summary
Four additive pipeline items analyzed. All preserve linear back-compat, agents-write-markdown, non-bypassable gates. Item 3 is infra-blocked; safe step is gh release create on tags. Recommended order: 4,1,2,3.

## Proposed approach
ITEM 1: Extend cursor schema with ready_blocks[]/active_blocks[] (additive, zero-value defaults, current_block kept). auto launches all ready blocks concurrently; track per-block round state in active_blocks map. Deterministic sort by block id for tests. Concurrency bound via flag (default 4). Trade-off: more cursor state vs simpler single-active.

ITEM 2: agreed_fixes as YAML frontmatter count in consensus.md (written by drafter from review aggregation). Add ReviewAgreedFixes helper. Loop in auto: review->consensus->signoffs->if fixes>0 and <max (default 3) run fixup->re-review (round-02+). Stop at 0/max/BLOCK. Fake agents for tests. Trade-off: one extra machine field vs full fix list.

ITEM 3: CONFIRMED BLOCKED - no releases for 1.6-1.9 (only tags); no real .exe/SHA256. Safe actionable: gh release create v1.6.0..v1.9.0 from tags (triggers CI assets). Do not invent hashes or open winget PR.

ITEM 4: Add agy-specific "write-first" preamble in runner launch (discover.go note + prompt branch). Generic preamble safe if short. Stdout fallback (validate starts with ---) only on --print failure path. Test with fake stdout capture. Trade-off: minor vendor note vs deeper abstraction.

## Concerns / open questions
ITEM1: test determinism with concurrent fakes; cursor migration for old runs. ITEM2: exact frontmatter key name; "applied" definition for re-review. ITEM4: preamble wording to avoid over-directive.

## Risks
ITEM1: race on shared state if not locked (mitigate with per-run mutex). ITEM2: infinite loop if count never hits 0 (cap prevents). ITEM3: release creation is low-risk but changes tag visibility. ITEM4: preamble may affect other agents (keep generic + short). All items additive; full fake-agent coverage required.