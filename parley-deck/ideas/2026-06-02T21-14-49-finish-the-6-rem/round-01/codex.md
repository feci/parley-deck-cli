---
agent: codex
idea: 2026-06-02T21-14-49-finish-the-6-rem
round: 1
date: 2026-06-02
---

## Summary

The six items are compatible with the existing §12 pipeline if implemented as additive extensions to the current manifest, cursor, gate, ledger, provider, watcher, and Driver layers. The safest plan is to keep agents responsible for markdown artifacts and code proposals, while the CLI/driver owns state transitions, JSON contracts, gates, ledgers, and external execution handoff.

Recommended build order:

1. Item 3: machine-readable execute output and MCP execution contract.
2. Item 2: Phase 5-8 implementation/review runner.
3. Item 1: auto-loop support for action, implementation, and watcher blocks.
4. Item 4: watch scheduler.
5. Item 6a/6c: per-block transport and DAG cursor.
6. Item 6b: decider-agent hook, only after gate policy tests are strong.
7. Item 5: WinGet manifests, in parallel once release metadata is known.

Item 6 should be split. Per-block transport, decider-agent, and DAG execution each touch different risk surfaces. DAG execution also likely needs §12 text updates because linear cursor semantics become insufficient. Item 5 is packaging work and should not block runtime pipeline work.

## Proposed approach

Item 1: auto-loop drives action/implementation/watcher blocks.

What is sound: extending `pipeline auto` is reasonable if it remains a driver over existing block semantics, not a bypass around gates. The proposed action-block behavior is safe if auto may run planning/finalization and dry-run execution, then must stop at production mutation gates.

Smallest correct implementation: add a dispatch table keyed by block kind. Deliberation keeps current behavior. Action blocks run Phase 1-4, finalize, invoke `execute --dry-run --json`, persist any planned effect metadata already allowed by current ledgers, then return a stopped status such as `blocked_by_gate` when production execution is required. Implementation blocks call the new Phase 5-8 runner. Watcher blocks finalize or validate `MONITORING.md`, enqueue/register watch metadata, and return without entering a long-running loop.

Stop semantics: auto may continue past non-mutating completed blocks, may stop on missing artifacts, failed validation, unresolved consensus, review fixes, external harness requirement, or any production gate. Auto must never call real provider execution for production effects.

Trade-off: this keeps auto useful but conservative; it may require users or CI to call `continue` more often at boundaries.

Item 2: Phase 5-8 implementation and code-review runner.

What is sound: reusing `CommandFor`, isolated homes, validation, and fake-agent launch machinery keeps the runner testable and avoids a second orchestration model.

Smallest correct implementation: add runner methods for `RunImplementation`, `RunReviewRound`, `RunReviewConsensus`, and `RunFixupCycle`. Use fixed artifact paths: `IMPLEMENTATION.md`, `review/round-NN/<agent>.md`, `review/consensus.md`, and updated `IMPLEMENTATION.md` entries per fix-up cycle. The implementer writes code and `IMPLEMENTATION.md`; non-implementers write review files with exact severities `CRITICAL`, `MAJOR`, `MINOR`, `NIT`. The facilitator/driver may draft review consensus, but signoffs remain participant-authored.

Detecting zero agreed fixes: parse `review/consensus.md` for a machine-readable agreed-fixes section, preferably YAML frontmatter or a fenced JSON/YAML block with `agreed_fixes: []`. Do not infer zero fixes from prose alone. If missing, block and ask for consensus repair.

Trade-off: requiring a structured consensus field adds protocol strictness, but makes Phase 8 deterministic and testable.

Item 3: machine-readable execute output and MCP execution contract.

What is sound: `execute --json` is the right boundary between CLI planning/ledgering and external harness MCP execution. It preserves the agents-write-markdown / driver-executes split and keeps vendor-specific execution outside core policy.

Smallest correct implementation: add `--json` to `pipeline execute` with a stable schema:

```json
{
  "status": "dry_run|pending_gate|ready_for_harness|recorded|error",
  "provider_call": {
    "provider": "string",
    "operation": "string",
    "target": "string",
    "arguments": {},
    "dry_run": true
  },
  "effect_digest": "sha256:...",
  "idempotency_key": "string",
  "gate": {
    "required": true,
    "state": "open|approved|denied|not_required"
  }
}
```

Harness contract: harness reads JSON, performs the MCP call only when gate state permits it, then calls `record-effect` with `idempotency_key`, `effect_digest`, `status`, `external_ref`, provider result summary, and error details when failed. Resume must reconcile ledger entries whose planned effect exists but lacks terminal success/failure.

Trade-off: the schema should be small and versioned; adding too much provider-specific detail will weaken vendor neutrality.

Item 4: live monitoring scheduler.

What is sound: watcher types and breach fingerprints already provide the core dedupe model, so `pipeline watch` can be additive.

Smallest correct implementation: add `parley pipeline watch SLUG [--signals FILE] [--once] [--interval DURATION]`. Define a `SignalSource` interface returning named values plus timestamps. Implement JSON-file source first. Parse `MONITORING.md` only if it has a structured watch spec; otherwise require a separate structured spec generated during watcher finalization. Persist dedupe state in the existing run/pipeline state area, keyed by watcher id plus breach fingerprint. For each new breach, either open a remediation idea only for allowlisted low-risk classes, or write a notify/gate event in the ledger.

Linkage: remediation ideas should carry parent pipeline slug, watcher id, breach fingerprint, source signal ids, and the triggering effect digest if applicable.

Trade-off: `--once` is easy to test and CI-friendly; loop mode should be thin polling over the same evaluator and not embed vendor monitoring clients.

Item 5: WinGet manifests for CLI 1.6/1.7/1.8 and skill 1.3.0.

What is sound: generating manifests without publishing is low risk and separable from runtime features.

Smallest correct implementation: inspect the existing packaging layout and add versioned manifest directories matching that shape. Use installer, locale, and version YAML files as the repo already does. Leave fields that require real release assets explicit and verifiable: download URL, SHA256, installer type, architecture, product code if applicable, release notes URL, and version.

Portable vs installer: if release assets are standalone `.exe` binaries, use the portable manifest shape. If CI produces an installer, use installer metadata. Do not guess asset hashes.

Trade-off: templated manifests are useful, but submitting to winget-pkgs must wait for immutable GitHub release assets and verified hashes.

Item 6: per-block transport, decider-agent tie-break, and DAG execution.

6a per-block transport: sound if credentials and execution scopes remain attached to the block and provider. Smallest implementation is an optional `transport` override on block config, defaulting to manifest-level transport. Validate that overrides are declared up front and never mutate cursor history. Risk is credential leakage or cross-transport confusion, so logs and ledgers must record effective transport per block.

6b decider-agent: sound only for low-risk, non-production gates where policy explicitly permits automated resolution. Smallest implementation is a gate policy hook that may ask a configured decider agent for a recommendation, then records the decision artifact and rationale. Default remains block-and-wait. Production mutation gates must be non-bypassable and ineligible.

6c DAG execution: sound but larger than the others. Smallest implementation is to preserve linear manifests as a DAG where each block depends on the previous block, then add `edges[]` validation, topological ready-set calculation, and a cursor state that tracks per-block status rather than a single current index. A block is ready when all inbound dependencies are complete and any edge gate is satisfied.

Trade-off: DAG support improves parallelism but complicates status, auto-loop ordering, retries, and user display. It needs protocol text updates for cursor semantics, ready-block selection, and artifact expectations when multiple blocks are active.

## Concerns / open questions

Item 1 needs explicit status names for auto stops so CI and humans can distinguish `needs_human_gate`, `needs_external_harness`, `needs_artifact`, and `failed_validation`.

Item 2 needs a structured review consensus format. Without it, "zero agreed fixes" becomes unreliable prose interpretation.

Item 3 should include a schema version, even if the initial version is `1`, and should define whether `execute --json` writes the planned ledger entry before or after printing JSON.

Item 4 should avoid making markdown parsing the source of truth unless `MONITORING.md` already has structured sections. A sidecar or embedded YAML block is safer.

Item 5 depends on actual release artifact type and SHA256 values. Missing real assets should block publication but not block manifest template generation.

Item 6 needs a §12 protocol update for DAG cursor semantics and probably for decider-agent eligibility. Per-block transport may be pure code if §12 already reserves it, but credential-scope rules should still be documented.

## Risks

Highest risk: accidentally allowing `auto` or decider-agent behavior to perform or approve production mutations. Mitigation: central policy evaluator must reject those paths regardless of command flags or manifest settings.

High risk: DAG cursor migration breaking existing linear manifests. Mitigation: treat old manifests as a generated linear DAG in tests and keep old status output stable unless `edges[]` is present.

Medium risk: implementation runner launches real agents in tests or depends on hosted credentials. Mitigation: fake agents must cover artifact creation, failed reviews, fix-up loops, consensus parsing, and launch failures.

Medium risk: `execute --json` schema churn breaking harnesses. Mitigation: add schema versioning, golden tests, and documented compatibility rules.

Medium risk: watcher dedupe either suppresses real repeated incidents or spams remediation ideas. Mitigation: persist fingerprints with explicit lifecycle states and include tests for same breach, changed breach, and resolved-then-reopened breach.

Low risk: WinGet manifests become stale before assets publish. Mitigation: keep them generated from release metadata and verify hashes in packaging tests or a release checklist.
