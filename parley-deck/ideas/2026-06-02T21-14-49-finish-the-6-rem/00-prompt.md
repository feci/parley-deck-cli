---
idea: 2026-06-02T21-14-49-finish-the-6-rem
author: user
created: 2026-06-02
participants: [codex, claude, hermes]
status: final
---

## Problem / idea

Finish the 6 remaining COOPERATION.md §12 pipeline items in parley-deck-cli. This is an IMPLEMENTATION idea (Phase 0-8): agree the implementation plan, then build it. The §12 protocol text and the internal/pipeline package (manifest, run cursor, gate + central policy evaluator, effects ledger, provider interface + VercelProvider, watcher types, Driver state machine) and the `parley pipeline validate|start|status|run-block|continue|auto|execute|record-effect|gate` command already exist and ship in cli 1.8.0. Build the remaining 6 items ON TOP without breaking existing behavior or the additive/backward-compatible guarantees.

Give independent Phase-1 analysis of the proposed approach for EACH item: what is sound, what is risky, what is missing, and the smallest correct implementation. Keep changes additive; preserve the agents-write-markdown / driver-executes boundary and non-bypassable production gates.

ITEM 1 — auto-loop drives action/implementation/watcher blocks.
Today `parley pipeline auto` drives only deliberation blocks and stops at others. Proposed: auto also handles (a) action blocks — run the Phase 1-4 plan rounds+consensus+finalize, then for the execute step, perform a dry-run automatically and STOP at the production gate (never auto-execute prod); (b) implementation blocks — drive via the new Phase 5-8 runner (item 2); (c) watcher blocks — finalize the MONITORING.md spec then hand off to the watch scheduler (item 4) without blocking. Question: exact stop/continue semantics per kind so auto never performs an unsafe action.

ITEM 2 — Phase 5-8 implementation+code-review runner.
Today only RunRoundOne/RunRound (Phase 1-4) exist. Proposed: a runner that (5) launches the implementer to produce IMPLEMENTATION.md + code on a branch, (6) launches each non-implementer reviewer to write review/round-NN/<agent>.md with CRITICAL/MAJOR/MINOR/NIT, (7) drafts review/consensus.md + signoffs, (8) fix-up loop until zero agreed fixes. Reuse the existing launch machinery (CommandFor, isolated home, validation). Question: how much to automate vs leave to the facilitator; artifact paths; how to detect "zero agreed fixes"; testability with fake agents.

ITEM 3 — machine-readable execute output + finalize MCP execution contract.
Today `pipeline execute` prints a human ProviderCall and writes the effect ledger; the driver/harness performs the actual MCP call. Proposed: add `--json` so execute emits {provider_call, effect_digest, status, idempotency_key} for an external harness/CI to consume and then call `record-effect`. Document the contract: agents write markdown, CLI plans + ledgers, harness executes via MCP, record-effect persists external_ref + reconcile. Question: exact JSON schema; how a harness signals success/failure back; reconcile-on-resume responsibilities.

ITEM 4 — live monitoring scheduler (`parley pipeline watch`).
Today watcher types + breach fingerprint/dedupe exist but nothing evaluates signals. Proposed: `parley pipeline watch SLUG [--signals FILE] [--once]` that loads MONITORING.md (or a watch spec), reads current signal values from a provided JSON file or a pluggable source, evaluates thresholds, dedupes breaches by fingerprint, and for each new breach either auto-opens a remediation idea (only predeclared low-risk classes) or writes a notify/gate record. Question: signal-source abstraction (avoid hard vendor coupling); dedupe persistence across runs; how remediation ideas link back to the pipeline; --once vs loop.

ITEM 5 — WinGet manifests for cli 1.6/1.7/1.8 + skill 1.3.0.
Proposed: generate the winget manifest files (installer/locale/version YAMLs) under the packaging dir for each new version, matching the existing manifest structure. Note that actual publication needs the GitHub release .exe assets (built by CI) + a PR to the winget-pkgs fork; we generate the manifests, not publish. Question: portable .exe vs installer manifest shape; per-version directory layout; what must be filled from real release assets vs templated.

ITEM 6 — per-block transport + decider-agent tie-break + DAG execution.
(6a) Per-block transport: allow each manifest block to override the pipeline transport (the schema reserves a single sticky transport in v1). (6b) Decider-agent: a policy hook that auto-resolves ONLY low-risk, non-production boundary gates via a configured decider agent, with block-and-wait as the default. (6c) DAG execution: replace linear-only with topological execution honoring edges[]; the driver advances any block whose inputs are all complete, still gating per edge. Question: backward compatibility (linear manifests unchanged); how DAG interacts with the cursor (multiple ready blocks); credential-scope risk of per-block transport; whether decider-agent belongs in §12 text (it is already reserved there).

CONSTRAINTS: additive + backward-compatible (existing manifests/decks unchanged); preserve agents-write-markdown / driver-executes boundary; production mutations always gate (non-bypassable); vendor-neutral (providers/transports behind interfaces); Go 1.26; full test coverage with fake agents (no reliance on launching real hosted agents in unit tests); English-only protocol files.

DELIVER in your round-01: per-item assessment + the minimal correct implementation + top risks; call out any item that should be split, deferred, or needs a §12 protocol-text change vs pure code; and the recommended build order.

## Constraints

- Local-directory transport for this initial run.
- Human-in-the-loop mode by default unless the run was started with auto mode.

## Non-goals

- Do not make unrelated repository changes.
