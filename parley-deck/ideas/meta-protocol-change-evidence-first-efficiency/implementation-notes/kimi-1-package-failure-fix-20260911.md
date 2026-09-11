---
agent: kimi-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-11
artifact-kind: owner fix handoff (post-timeout recovery)
not-a-signoff: true
---

# Kimi evidence slice — package-failure fix 2026-09-11 (OWNER-RUN, NOT INDEPENDENT)

Correction invocation f750e9ad-b977-44c4-95f3-cc9ce51a77ee timed out at 1800.913s; its source and test edits survived in the worktree and are retained unchanged. This is the missing handoff. Not Phase-6 review, not final acceptance, not production closure integration (verifier close wiring remains absent, facilitator-side).

Fix (codex-1's MAJOR probe TestCodexProbePackageFailureCannotPass): parseGoTestJSONEvidence previously skipped package-level events, so a masked exit 0 turned package fail/build-fail into PASS. Now:

- execute.go — package-level fail/build-fail events increment a new failedPackages signal (sticky failure; never folded into test-case counts, no invented executed cases for a failed build); RunCriterion fails closed when FailedPackages > 0 despite exit 0, with a diagnostics header; the ParseGoTestJSON wrapper reports such streams recognized with zeroed counts — never a pass.
- evidence.go — CommandEvidence.FailedPackages (-1 = unknown/unsupported for shell/envelope); Evaluate, AttestExecution and reconcileRerun reject >0, require the count on gotest-json records/re-runs, and reconcile it cross-execution.
- execute_test.go / evidence_test.go — masked-exit fail/build-fail, mixed-package, real-stream, and evaluation/attestation cases.

Owner-run results (preserved logs, owner-executed): evidence suite 69 pass / 0 fail / 1 skip (known unsupported-filesystem skip preserved); app suite 13 / 0 / 0; six read-only codex overlay probes pass; real-stream check surfaced failedPackages=2. Usage/cost: unknown. SECONDARY (facilitator/Codex, not mine): selected evidence/app suite 53 events clean; independent Codex full evidence-package run passed 85 events, zero skips/failures.
