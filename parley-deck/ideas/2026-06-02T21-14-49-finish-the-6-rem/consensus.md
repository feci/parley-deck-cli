---
idea: 2026-06-02T21-14-49-finish-the-6-rem
drafted-by: claude
date: 2026-06-02
---

## Agreed decisions

Round-01 converged (codex, claude, hermes) with no substantive blocker. All six items are built as **additive** code on the shipped 1.8.0 §12 surface; linear manifests and existing decks stay byte-for-byte unchanged; the agents-write-markdown / driver-executes boundary and non-bypassable production gates are preserved.

**Build order:** 3 → 5 → 4 → 2 → 1 → 6 (smallest/independent first; the keystone Phase 5-8 runner before auto-for-impl; DAG last).

1. **Item 3 — `execute --json`.** Add `--json` to `pipeline execute` emitting a **versioned** schema: `{schema_version, status (dry_run|pending_gate|ready_for_harness|recorded|error), provider_call{provider,tool,args,dry_run}, effect_digest, idempotency_key, gate{required,state}}`. The ledger entry is written **before** the JSON is printed. Harness contract documented in `references/EXECUTION_CONTRACT.md`: agents write markdown; CLI plans + ledgers; harness performs the MCP call only when the gate permits; `record-effect` persists `external_ref`/status; resume reconciles non-terminal effects.
2. **Item 5 — WinGet manifests.** Generate versioned manifest dirs for cli 1.6.0/1.7.0/1.8.0 and skill 1.3.0 under the existing packaging layout (version + installer + locale YAMLs), templated from the current manifests. `InstallerSha256`/URL fields carry an explicit `PLACEHOLDER-FILL-FROM-RELEASE-ASSET` marker — never invented. We generate; publication waits for immutable GitHub release assets + the winget-pkgs PR.
3. **Item 4 — `pipeline watch`.** `parley pipeline watch SLUG --signals FILE [--once]`. A `SignalSource` interface (JSON-file source first) keeps it vendor-neutral. Evaluate thresholds from a **structured** watch spec (embedded YAML/sidecar, not prose). Dedupe by `Breach.Fingerprint` **persisted** under `pipelines/<slug>/breaches/` with lifecycle (open/resolved) so the same ongoing breach does not respawn. Per new breach: open a remediation idea ONLY for predeclared low-risk classes (linked via `derived_from` + parent slug + fingerprint), else write a notify/gate record. `--once` is the tested core; loop is a thin sleep wrapper.
4. **Item 2 — Phase 5-8 runner.** Add `runner.RunImplementation`, `runner.RunReviewRound`, and fix-up support reusing CommandFor/isolated-home/validation. Artifact paths: `IMPLEMENTATION.md`, `review/round-NN/<agent>.md`, `review/consensus.md`. **"Zero agreed fixes" is detected from a machine-readable field** (a fenced `agreed_fixes:` list / structured section in `review/consensus.md`), never inferred from prose. Loop orchestration lives in the app/pipeline layer; launch primitives in runner. Tested with fake agents (helper-process pattern).
5. **Item 1 — auto for all block kinds.** Extend `runPipelineAuto` with a per-kind dispatch: deliberation = today; **action** = plan rounds+consensus+finalize, then `execute --dry-run` automatically, then STOP at the production gate with an explicit status (`needs_human_gate`); **implementation** = drive the item-2 runner, stop after review with agreed fixes or on completion; **watcher** = finalize the spec, print the `watch` command, continue. Explicit auto-stop statuses: `needs_human_gate`, `needs_external_harness`, `needs_artifact`, `failed_validation`. **auto NEVER calls real production execution** (no non-dry-run execute).
6. **Item 6 — split into 6a/6b/6c.**
   - **6a per-block transport:** optional `Block.Transport` overriding the manifest transport (same allow-list); effective transport recorded per block in state. Small, pure code.
   - **6b decider-agent:** optional `Manifest.Decider`; the central policy evaluator MAY consult it to auto-resolve ONLY low-risk, non-production boundary gates; block-and-wait stays default; production gates remain ineligible/non-bypassable. Small.
   - **6c DAG execution:** generalize the cursor to `completed_blocks[]` + `ready_blocks[]` (additive, `schema_version` bump, zero-value defaulting; `current_block` still populated for linear back-compat). The driver advances any block whose inbound edges are all complete and whose edge gates are resolved. Linear is the degenerate case. Largest/riskiest; ships last with the most tests.

## Agreed trade-offs

- More automation (auto, watch, Phase 5-8) concentrates trust in the driver and adds protocol strictness (structured `agreed_fixes`, versioned JSON) — accepted for determinism and testability.
- Per-block transport widens credential scope — accepted only with effective-transport logging per block.
- DAG complicates status/display and likely needs a **§12 text clarification** (cursor semantics + ready-set + multi-active artifacts) — that clarification is in-scope as a small protocol note, not a new meta-change idea, since §12 already reserves `edges[]`.
- Item 5 is packaging-only and must not block runtime work.

## Open items deferred to implementation

- Exact `references/EXECUTION_CONTRACT.md` wording and `execute --json` golden test.
- The structured watch-spec format inside `MONITORING.md` (embedded YAML block).
- The structured `agreed_fixes` format inside `review/consensus.md`.
- §12 text note for DAG cursor semantics + decider-agent eligibility (applied to both canonical + dogfood COOPERATION.md at finalize time).
- WinGet `InstallerSha256`/URL values (filled from real release assets later).

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: codex — 2026-06-02
Status: ✅ ACCEPT
Notes: Accept. The consensus captures codex's safety requirements: structured automation contracts, dry-run-only auto at production boundaries, non-bypassable gates, additive cursor/schema evolution, and DAG last with §12 clarification.
Counter-proposal (required if ❌): N/A

### Signoff: hermes — 2026-06-02
Status: ✅ ACCEPT
Notes: The consensus fully reflects the independent analysis and recommended build order from round-01. All safety invariants (additive changes, non-bypassable gates, dry-run-only auto, agents-write-markdown boundary) are preserved.
Counter-proposal (required if ❌): N/A

### Signoff: claude — 2026-06-02
Status: ✅ ACCEPT
Notes: Drafter. Plan + build order 3-5-4-2-1-6 reflect round-01 convergence; codex's structured-contract refinements (versioned execute --json, machine-readable agreed_fixes, sidecar watch spec, explicit auto-stop statuses, DAG needs §12 note) adopted. Safety invariants preserved.
