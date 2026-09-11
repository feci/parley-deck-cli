---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-11
artifact-kind: partial implementation checkpoint
not-a-signoff: true
---

# Independent execution path and outstanding closure defects

The production driver now calls a non-implementer configured CLI through
runner.RunConsult. The selected agent must invoke the runtime-only command
`parley evidence verify --request PATH --request-sha256 SHA`. Textual PASS,
missing execution, failed process, self-verification, failed/skipped/package-
failed tests, changed code or changed non-evidence scope cannot close the named
contract in the current process fixtures. An adapter lacking the independent
execution hook cannot close a named contract.

Each attempt retains exact original report bytes, a hashed request and actual
criterion executions in a helper receipt. Before and after each criterion,
the helper verifies current scope/commands, tested tree and non-evidence
implementation digest. Parent validation binds the receipt to the actual unique
invocation marker, verifier, request and persisted report, and only accepts
verifier provenance added to original evidence. Complete rechecks the evidence
immediately before writing its generated status.

Runtime process attribution is not human authentication or protection against
malicious same-UID writers. Tests use a fixture CLI agent plus an actual built
Parley helper and actual Go test subprocesses; they are not live model trials.
Only prior review consensus is stubbed in the Driver.Advance process fixture.

## Executed checks

- Full `go test -count=1 ./...`: PASS. App 84.811s, driver 4.522s,
  evidence 2.976s, runner 24.872s. Log:
  `.parley-runtime/evidence-verifier-full-go-20260911.log`.
- Scoped vet app/driver/runner/evidence: PASS.
- Windows app test-binary cross-build: PASS; runtime not tested.
- Corrected focused real-process suite: PASS. The initial marker expectation
  and implementation-scope fixture errors remain in historical logs.
- New, stronger positive-case overlay: FAIL after successful completion, in
  3.224s. Log `.parley-runtime/evidence-post-completion-counterexample-20260911.log`.
  It adds a fresh `EvidenceCloseGate` call after Driver.Advance completes;
  the gate refuses because the generated status transition changed the
  non-evidence IMPLEMENTATION.md digest. The overlay is retained privately.

## Open acceptance blockers

1. Preserve auditable evidence after the exact generated status transition.
   The separately hashed scope must remain binding; excluding status or broad
   frontmatter silently is not a correction. An exact independently authorized
   before/after transition can be designed and reviewed explicitly.
2. Serialize cooperative report writers and use a real snapshot/CAS boundary.
   Current pre-save byte comparison is not atomic with report Save. No concurrency
   safety claim is made for this checkpoint.
3. Complete's direct-call contract deletion/race behavior, report and receipt
   persistence failures, and current-source independent review remain open.
4. All-launch budget wiring, actual launch-surface trials, exact packet and
   comparative experiments and final owner-authored reviews/signatures remain
   unfinished. No final merge, release, deployment or global install is claimed.

Graphify vocabulary expansion used evidence/digest/complete/verification to
locate completion glue. The graph did not contain the new uncommitted helper;
all implementation and defect claims above are from direct source inspection
and actual executions, not inferred graph edges.
