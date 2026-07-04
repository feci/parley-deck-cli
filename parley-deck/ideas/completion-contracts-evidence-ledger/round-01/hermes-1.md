---
agent: hermes-1
idea: completion-contracts-evidence-ledger
round: 1
date: 2026-07-04
---

## Summary

The proposal targets a real gap: Phase 8 "complete" is a prose claim plus reviewer consensus, and the one overclaim incident on record (codex BLOCK on protocol-restructure-appendices) was claim-accuracy, not a code defect. But the brief bundles four changes — a `done_when:` contract, an evidence ledger artifact, driver enforcement, and §5/§8 protocol edits — when the minimal fix is mostly already present in the protocol.

Key observation: the protocol ALREADY has the verification primitive this proposal needs. `checks:` (LE-4, COOPERATION.md line 236) is a build/test gate the driver runs as `sh -c` at Phase 5/8, and `RunChecks` in `internal/app/driver_impl.go:210` already captures exit code + combined stdout/stderr. `IMPLEMENTATION.md` already has a `## Validation evidence` section (COOPERATION.md lines 446-448) asking the implementer to record "the commands run and what they proved." And the strict_gate close path (`internal/driver/impl.go:202-216`) already demonstrates the exact enforcement pattern this proposal needs: a deterministic fail-closed veto on a drafter's clean claim.

So the smallest protocol-text change that delivers evidence-grounded completion is NOT four amendments. It is: (a) generalize `checks:` from one command to an optional named list, (b) one sentence making the existing `## Validation evidence` section the canonical ledger when the list is present, and (c) one fail-closed veto line in the Phase 8 close rule. The driver change is ~20 lines mirroring the strict_gate veto.

## Proposed approach

Minimal delta, in order of size:

1. **Extend `checks:`, do not add `done_when:`.** Today `checks:` is a single string. Make it ALSO accept a YAML list of named criteria, each `{name, command, expect}` where `expect` defaults to `exit: 0`. A single string keeps today's semantics exactly (backward compatible). A list turns on the contract. One field, two shapes — no new frontmatter key, no overlap confusion. The driver's `RunChecks` already reads `checks:` from frontmatter (`driver_impl.go:221-223`); the list case loops and records per-criterion `(name, exit_code, duration, output_tail)`.

2. **The ledger IS `## Validation evidence` in IMPLEMENTATION.md.** No new `review/evidence.md` artifact. The section already exists (COOPERATION.md line 446) and already asks for commands + what they proved. The protocol edit is: when `checks:` is a list, the driver writes the per-criterion table into that section (append-per-fix-up-cycle) and reviewers receive IMPLEMENTATION.md as review input — which they already do. This avoids a second source of truth that would drift from IMPLEMENTATION.md.

3. **One fail-closed veto line in §8.** Mirror strict_gate's pattern (`impl.go:208-210`): when `checks:` is a list, the driver vetoes `status: complete` unless every criterion passed in the latest evidence entry. The veto can only fail a close claim, never auto-pass one — same LE-2 shape as strict_gate. This is ~20 lines of Go, reusing `RunChecks`'s existing capture.

4. **Protocol text scope.** The "evidence beats prose" rule MUST be scoped to "when `checks:` is a list." Written as a general principle it retroactively changes review semantics for every existing idea — a backward-compat break disguised as additive. Two sentences in §5 and §8, conditional on the list shape. Nothing in §0-§4, nothing in the appendices.

Total protocol-text delta: ~6-10 sentences across two sections, one frontmatter shape change, one new Go veto branch. No new artifact path. No new lifecycle rule beyond "append per cycle" which IMPLEMENTATION.md already supports.

## Concerns / open questions

- **Is the ledger even necessary, or is the fix a nudge?** The cited overclaim was a prose-accuracy problem. The implementer already had `## Validation evidence` available and did not fill it rigorously. A protocol amendment that says "when `checks:` is a list, the driver populates `## Validation evidence` automatically" may be the whole fix — the contract is just "make the existing section machine-filled instead of human-filled." If so, the `done_when:` naming and the separate ledger artifact are framing overhead. Question for the facilitator: is the goal a new contract concept, or is it just "stop trusting prose when checks ran"?

- **`output pattern` / `file existence` matchers.** The brief lists these as criterion types. Each is a new mini-spec: glob vs regex, encoding, error messages. The minimal version supports only `exit: 0` (what `checks:` does today) and treats anything beyond as a future extension. Adding matchers now is speculative surface area. Trade-off: less expressive contracts vs. no parser to maintain.

- **Append-only vs. rewrite-per-cycle.** The brief says append-only per fix-up cycle. IMPLEMENTATION.md `## Validation evidence` is currently free-form. Forcing append-only introduces a lifecycle rule the protocol doesn't currently enforce on any section. Minimal: let the driver overwrite the section each cycle with the latest run (the audit trail lives in git history of IMPLEMENTATION.md, which is already the case for every other section). Append-only is a nice property but is extra machinery for a benefit git already provides.

- **Interaction with `strict_gate`.** Both are opt-in close gates. An idea could set both `checks: [list]` and `strict_gate: true`. The order of enforcement matters: does a green evidence ledger satisfy strict_gate, or does strict_gate still require a fresh full-scope review round? The brief is silent. Minimal answer: they are independent — evidence ledger is a necessary condition (checks passed), strict_gate adds the sufficient condition (clean review). But this MUST be stated, or the two gates will collide in the driver.

- **Digest/truncation/secret-safety.** The brief raises this. `RunChecks` already captures full `buf.String()` today (`driver_impl.go:218`). The minimal safe choice: record exit code + last N lines of output (truncated at a fixed cap, e.g. 4KB), no digest scheme. Secrets are bounded by the same trust model as RunChecks — the contract author wrote the command, and the command is visible in review. A digest adds crypto surface area for no security gain over truncation.

## Risks

- **Overreach via generalization.** The biggest risk is writing "evidence beats prose claims" as a universal protocol principle. That silently changes every closed idea's review semantics and every future idea without a contract. The rule must be conditional on the contract's presence, or the proposal stops being backward compatible despite the claim.

- **Two verification fields.** If the proposal keeps `done_when:` as a new key alongside `checks:`, authors will confuse which one the driver honors. Two keys for the same concept is a protocol smell. Extending `checks:` to accept a list avoids this entirely.

- **Second source of truth.** A separate `review/evidence.md` drifts from `IMPLEMENTATION.md ## Validation evidence`. Reviewers already read IMPLEMENTATION.md; making them also read a second artifact for the same facts is friction and a drift hazard. One canonical location is more minimal.

- **Solving a social problem with machinery.** The root cause was an implementer's prose overclaim where checks had run. The protocol already required recording evidence; it wasn't followed. A new artifact doesn't fix non-compliance with the existing one — but a driver that AUTO-fills the section from `RunChecks` output does, because it removes the human step that failed. The risk is building a contract concept when the actual fix is "auto-populate the existing section."

- **Drift guard cost.** Any §5/§8 edit lands in both `parley-deck/COOPERATION.md` and `internal/protocol/defaults/COOPERATION.md`. A 10-sentence delta is 20 sentences to keep in sync. Keep the delta small or the maintenance cost dominates the benefit.
