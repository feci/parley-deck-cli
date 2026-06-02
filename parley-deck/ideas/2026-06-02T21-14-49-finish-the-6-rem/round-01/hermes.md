---
agent: hermes
idea: 2026-06-02T21-14-49-finish-the-6-rem
round: 1
date: 2026-06-02
---

## Summary
Independent Phase-1 analysis for the 6 items. All changes additive; preserve agents-write-markdown/driver-executes boundary and non-bypassable prod gates. Recommended order: 3,1,2,6,4,5 (core safety first, then automation, packaging last). Items 3 and 6 require minor §12 text clarifications; others pure code.

## Proposed approach
ITEM 1: Extend auto to drive action (dry-run only, stop at gate), impl (reuse item2 runner), watcher (finalize spec + handoff). Stop semantics: action=stop-before-prod, impl=stop-after-review, watcher=non-blocking. Minimal: add block-type switch in auto loop + dry-run flag.

ITEM 2: New Phase5-8 runner: (5) implementer branch+IMPLEMENTATION.md, (6) reviewers write CRITICAL/MAJOR etc under review/round-NN/, (7) consensus+signoffs, (8) fix-up until zero CRITICAL/MAJOR agreed. Reuse CommandFor/isolated homes. Detect zero fixes via consensus count. Test with fake agents.

ITEM 3: Add --json to execute emitting {provider_call, effect_digest, status, idempotency_key}. Document contract: agents MD, CLI plans/ledgers, harness MCP, record-effect persists external_ref. Schema: versioned, minimal fields.

ITEM 4: `pipeline watch SLUG [--signals FILE] [--once]`: load MONITORING.md, eval thresholds from JSON/pluggable source, dedupe by fingerprint (persist in ledger), low-risk auto-remediation idea or notify/gate. Abstraction via SignalSource interface.

ITEM 5: Generate winget YAMLs (installer/locale/version) under packaging/ per version; template most fields, fill release .exe from CI assets. Dir layout: per-version subdirs. No publish.

ITEM 6: (6a) per-block transport override (additive field), (6b) decider-agent hook for low-risk gates only (default block), (6c) DAG topo exec with cursor tracking ready blocks. Backward compat for linear manifests. Credential scope limited to block.

## Concerns / open questions
- Exact JSON schema + reconcile semantics for item 3 (needs §12 update).
- "Zero agreed fixes" detection heuristic for item 2; fake-agent test harness coverage.
- SignalSource pluggability and dedupe persistence for item 4.
- Decider-agent policy location (code vs §12) for item 6.

## Risks
- High: unsafe auto-execution if stop semantics wrong (item1) — mitigate with explicit dry-run + gate.
- Med: DAG cursor complexity + multiple ready blocks (item6) — keep simple topo + single-active default.
- Low: WinGet manifest drift from real assets (item5) — template + CI validation only.
- Trade-off: more automation (1,2,4) vs facilitator control; chose minimal runner + explicit gates.