---
agent: antigravity-1
idea: sync-skill-protocol-fallback
round: 1
date: 2026-06-24
---

## Recommendation

I explicitly agree with [claude-1.md](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/parley-deck/ideas/sync-skill-protocol-fallback/round-01/claude-1.md)'s recommendation to adopt **Strategy (A): Verbatim copy of the CLI embedded default** for the fallback protocol file. 

Specifically, we should make [references/COOPERATION.md](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-skill/references/COOPERATION.md) in the skill repository a byte-for-byte identical copy of [COOPERATION.md](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/protocol/defaults/COOPERATION.md) from the CLI defaults.

This strategy is superior because:
1. **Single Source of Truth**: It removes the maintenance burden of manually merging and re-formatting updates, completely preventing future protocol drift.
2. **Acceptable CLI Literals**: The `parley` CLI literals (e.g., `parley init`, `parley loop tick`) are used solely as illustrative examples of generic loop-engineering capabilities, which is standard documentation practice. Non-parley agents will easily interpret these as examples, not hard dependencies.
3. **Immaterial Placeholder Differences**: Cosmetic differences in header syntax (e.g., `<transport-choice>` vs. `github-pr`) do not affect protocol execution. In fact, keeping the default identical to the CLI template ensures that the protocol matches the initial state of a freshly-bootstrapped deck.

We should publish this sync as a **patch version bump (1.4.1)** in [package.json](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-skill/package.json) because this represents a documentation/reference asset update with no change to the skill's API contract, entrypoints, or programmatic interfaces.

## Analysis

### 1. Comparison of Missing & Changed Sections
A line-by-line comparison of the stale fallback (14 sections, 759 lines) and the canonical CLI template (16 sections, 1037 lines) reveals critical gaps in safety and automated loop controls:

* **Header Placeholder Style**: 
  * The stale version uses `<transport-choice>` and `<YYYY-MM-DD>` as placeholder strings. 
  * The canonical template uses `github-pr` as a default transport and `<date> — created by parley init` for creation date. This change has no impact on runtime neutrality, as the text explicitly instructs facilitators to replace the header.
* **§0 Deck Bootstrap (One-Time)**:
  * Missing entirely in the stale fallback. This section mandates that the facilitator confirm the active roster, agent models, and reasoning/effort levels upon deck creation. It is a critical setup step ensuring that agents run with optimal settings (highest reasoning level by default).
* **§4 Phase Rules (Phase 6/8 Rework, Loop-Budget & Close Integrity)**:
  * The stale version lacks crucial loop-engineering constraints:
    * **`strict_gate`**: An optional gate ensuring a clean Phase 6 review round (zero findings of any severity) before merging.
    * **`require_model_diversity` (LE-3)**: Prevents rubber-stamping reviews by ensuring reviewers do not all use the same model family as the implementer.
    * **`checks` (LE-4)**: The automated test verification command.
    * **Refutation-Default (LE-1)**: Reverses the reviewer posture (code is assumed buggy until proven otherwise, requiring documented break attempts under `## Refutation attempts`).
    * **Loop budgets (LE-5)** and **Close-decision integrity (LE-7/LE-11)**: Essential safety bounds preventing runaway infinite execution loops or premature merges under automated orchestration.
* **§12.11 Candidate-Remediation (LE-10)**:
  * The canonical version adds rules detailing that watcher-auto-opened remediation ideas are initialized with `status: candidate` and contain no active quorum or participants. This maintains the non-solo Phase-0 invariant by requiring a human or manifest to explicitly promote the candidate.
* **§13 Retrospective Optimization (RHO)**:
  * Missing entirely in the stale fallback. It details how the deck can periodically analyze its own history (coreset analysis) to propose harness updates, while strictly requiring that such proposals pass through normal multi-agent review and human gates.
* **§14 Automated Outer Loop (The Human Brake)**:
  * Missing entirely in the stale fallback. It sets the binding safety constraint that cron jobs, CI triggers, or scheduled loops may only discover signals and draft candidate prompts; they are explicitly forbidden from merging code, writing implementation, or altering rosters without a human or full-quorum gate.

### 2. Evaluation of Strategies
* **Strategy A (Verbatim Copy)**: Eliminates divergence, simplifies future upgrades, and maintains exact parity with what the CLI embeds. The presence of illustrative `parley` CLI literals is a non-issue since they are written as examples (prefixed with `e.g.` or `such as`).
* **Strategy B (Genericized Merge)**: Attempting to strip `parley` mentions and restore custom `<transport-choice>` markers requires maintaining a custom fork of the markdown file. This creates manual verification overhead and is prone to human error, which is the root cause of the current multi-generation drift.

## Verification

The following verification checklist must be executed before publishing the updated package to npm:
1. **Identical Check**: Run `diff -u references/COOPERATION.md ../parley-deck-cli/internal/protocol/defaults/COOPERATION.md` and ensure the output is completely empty.
2. **Anchor Validation**: Verify that the new sections are present exactly once:
   ```bash
   grep -c "## 13. Retrospective optimization" references/COOPERATION.md  # Expected: 1
   grep -c "## 14. Automated outer loop" references/COOPERATION.md        # Expected: 1
   ```
3. **Leak Inspection**: Perform a visual inspection (or run grep) to ensure no project-specific roster entries (e.g., user handles, specific active agent tables) were accidentally embedded.
4. **Test Suite Verification**: Run `npm test` inside `parley-deck-skill` to verify that no installer or validation tests are broken by the updated fallback document.
5. **Package Integrity Dry-Run**: Execute `npm pack --dry-run` to ensure all asset paths are valid and the new file is correctly packaged.
6. **Release Execution**:
   * Bump version to `1.4.1` in [package.json](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-skill/package.json).
   * Publish to npm: `npm publish --access public`.
   * Tag the release: `git tag v1.4.1` and push the tag.

## Open questions

1. **Anti-Drift CI Guard**: Should we add a test to the `parley-deck-skill` repository (like [drift_test.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/protocol/drift_test.go) in the CLI repo) that automatically asserts that [references/COOPERATION.md](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-skill/references/COOPERATION.md) matches the CLI's defaults copy? This would prevent regression before a release.
2. **Automated Sync Task**: Should the build process of `parley-deck-skill` include a script that automatically pulls the canonical `COOPERATION.md` from the CLI repository if they are located in a monorepo or accessible path?
