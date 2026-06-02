---
agent: claude
idea: 2026-06-02T21-14-49-finish-the-6-rem
round: 1
date: 2026-06-02
---

## Summary

All six are tractable as additive code on top of the shipped 1.8.0 surface, but they are not equal in size or risk. Items 3, 4, 5 are small/medium and self-contained. Item 2 (Phase 5-8 runner) is the keystone — item 1 (auto for impl blocks) depends on it. Item 6 splits into three sub-features of very different cost: per-block transport (small), decider-agent (small), DAG execution (the single largest and riskiest change). I recommend building in dependency order 3 → 5 → 4 → 2 → 1 → 6, and treating DAG as the one item that may need explicit §12 text plus the most test scaffolding. Everything must keep linear manifests and existing decks byte-for-byte unchanged, keep the agents-write-markdown / driver-executes boundary, and keep production gates non-bypassable.

## Proposed approach

- **Item 3 (execute --json):** smallest. Add `--json` to `execute`; emit `{idempotency_key, effect_digest, status, provider_call, dry_run}`. record-effect already takes the digest. Document the harness contract in a short `references` note. No new concepts. Do first to lock the execution contract the other items lean on.
- **Item 5 (WinGet manifests):** mechanical. Generate `packaging/winget/manifests/f/Feci/ParleyDeckCli/<ver>/` (and skill) version+installer+locale YAMLs templated from the existing 1.x manifests, with InstallerSha256 left as a documented placeholder to fill from the real CI `.exe` assets. We generate, not publish. No Go code.
- **Item 4 (watch):** medium. `parley pipeline watch SLUG --signals FILE [--once]`. Define a `SignalSource` interface (default: JSON file reader) so it stays vendor-neutral. Load `MONITORING.md` watch spec (parse frontmatter/YAML block), evaluate thresholds, dedupe by `Breach.Fingerprint` persisted in `pipelines/<slug>/breaches/`, and per new breach either record a notify/gate file or (predeclared low-risk class only) open a remediation idea linked via `derived_from`. `--once` evaluates a single pass; loop mode is a thin wrapper with a sleep, but keep `--once` the tested core.
- **Item 2 (Phase 5-8 runner):** the keystone. Add `runner.RunImplementation` (launch implementer → IMPLEMENTATION.md + branch), `runner.RunReviewRound` (each non-implementer writes `review/round-NN/<agent>.md`), and detection of "zero agreed fixes" by reading `review/consensus.md`'s Agreed-fixes section. Reuse CommandFor/isolated-home/validation. Keep the loop orchestration in the app/pipeline layer; keep the launch primitives in runner. Test with fake agents (helper-process pattern, like round_test.go).
- **Item 1 (auto for all kinds):** once item 2 exists, extend `runPipelineAuto`: action → plan rounds+consensus+finalize, then `execute --dry-run` automatically, then STOP at the production gate (never auto-exec prod); implementation → drive item 2's loop; watcher → finalize spec then print the `watch` command and continue/stop. Pure orchestration glue over items 2/3/4.
- **Item 6 (transport/decider/DAG):**
  - per-block transport: add optional `Block.Transport`; default to manifest transport; validate against the same allow-list. Small.
  - decider-agent: add `Manifest.Decider` (agent id) + policy so `AutoApprove` may consult a decider ONLY for low-risk non-prod gates; block-and-wait stays default. Small; the gate/policy seam already exists.
  - DAG execution: generalize the cursor from a single `current_block` to a set of ready/completed blocks; the driver advances every block whose incoming edges are all from completed blocks and whose edge gates are resolved. Keep linear as the degenerate case. This is the big one and the only piece I'd gate behind extra tests + possibly a §12 clarification that a v1.1 driver MAY execute declared DAGs.

## Concerns / open questions

- **DAG + cursor model:** `pipeline-run.json` today has a single `current_block`. DAG needs `ready_blocks[]` / `completed_blocks[]`. Risk: breaking the existing cursor schema. Mitigation: additive fields + keep `current_block` populated to the first ready block for linear back-compat; bump a cursor schema_version with zero-value defaulting.
- **Phase 5-8 testability:** real implementer/reviewer launches are slow/flaky. We must unit-test the loop with fake agents and a file-based "zero agreed fixes" predicate, not live agents.
- **watch dedupe persistence:** breaches must dedupe across invocations, so fingerprints need to be persisted (a `breaches/` dir), not just in-memory.
- **auto + action blocks:** auto must do dry-run then HARD STOP at the prod gate. It must never call execute without `--dry-run` unless a human-approved gate already exists — and even then I'd keep real execution out of `auto` entirely (auto reaches the gate; a human runs `execute`).
- **Scope:** this is ~6 features. I suggest one FINAL.md plan but multiple implementation commits/sub-reviews, or splitting DAG into its own follow-up idea if review load is too high.

## Risks

- **Over-scope in one idea** → shallow review. Mitigation: build order with independent, separately-testable commits; DAG last and optionally its own idea.
- **Cursor schema break** (DAG) → existing pipelines fail to load. Mitigation: additive + schema_version + defaulting; round-trip test old cursor.
- **auto performing an unsafe action** → the cardinal risk. Mitigation: auto never runs `execute` without `--dry-run`; production gate non-bypassable; explicit test that auto stops at an action block's prod gate.
- **Vendor coupling in watch** → mitigate with a `SignalSource` interface (file source default).
- **WinGet hashes** → never invent InstallerSha256; leave a clear placeholder + doc that it is filled from real release assets.
- Smallest first increment I would ship even if time runs short: items 3 + 5 + 4 + the per-block-transport and decider-agent halves of item 6 (all small/medium, fully testable), deferring Phase 5-8 runner (item 2), auto-for-impl (item 1), and DAG (item 6c) if review bandwidth is the bottleneck — but the goal is all six.
