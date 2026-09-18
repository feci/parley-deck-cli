---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-10
artifact-kind: partial-slice independent check
source-sha256: 00a428ba7f10be58809f8982e9745a7d772271f89da386d7ee2171b1318a57e4
not-a-signoff: true
---

# Kimi recovery checkpoint — independent rejection remains

The second owner attempt c16ced7b-14d0-4e82-80d6-0f88d035fb36 ended at its
30-minute hard deadline: 1800.905 seconds, exit 143, failure_class=timeout.
Its early handoff exists but never acquired a final status. All source and
owner-authored notes were preserved; this slice is not merged into integration.
Token counts, reported model and monetary cost are unknown, not zero. The
requested model was kimi-code/k3. Raw diagnostics contain two ChatProviderError
retry notifications, but no quota or sandbox-denial marker was observed; that
is not a complete diagnosis of elapsed time.

## PRIMARY checks on the stopped snapshot

A private isolated copy under .parley-runtime/evidence-probe-2 preserves the
production files and the executable probe. It fixes several previously
reproduced cases: empty verifier attestation now rejects; a malformed final
envelope rejects; symlink retarget and executable-mode changes alter the digest.
The actual probe output also exposes two remaining false acceptance paths:

```
duplicate_field_rejected=false
unrelated_unbound_rerun_attestation_accepted=true close_reasons=[]
```

### [MAJOR] Independent rerun is not bound or retained

internal/evidence/evidence.go AttestExecution accepts an unrelated command hash
and unrelated output hash, does not receive the rerun's tested-tree digest, and
stores only verifier identity plus the ORIGINAL report tree hash. The supplied
rerun evidence is discarded. A caller can therefore attest a different command
or tree and Evaluate returns no reasons. This is a concrete API consistency
failure even within the stated same-runtime asserted-identity trust boundary;
it is not a demand for cryptographic authentication of the human.

Retain the independent execution record, bind criterion, exact command and
before/after tested tree, validate hashes/counts/status, and reconcile it at
close. Integrating a mere name-stamping path would violate FINAL D3.

### [MAJOR] Duplicate JSON field can replace the executed-case count

ParseEnvelope accepts a single JSON object containing executed_cases=0 followed
by executed_cases=1. Standard json.Unmarshal keeps the latter. Reject duplicate,
aliased, unknown and otherwise contradictory envelope fields before semantics;
do not treat permissive decoding as evidence of execution.

## Owner-package/integration fixtures that still fail

The focused command was go test ./internal/evidence ./internal/app with the
Evidence/Criterion/Envelope/TreeDigest/Typed/Verifier/Barrier name filter:

- TestTreeDigestUnsupportedEntryFails: the Git file inventory omits an untracked
  FIFO, so unsupported entries are not actually rejected as documented.
- TestRunChecksContractEvidenceWriteFailureVetoes: its fixture fails earlier at
  an empty-tree digest and never exercises the intended persistence failure.
- TestEvidenceCloseGatePositive: the old positive fixture has no persisted
  verifier attestation and is now correctly rejected; a real positive needs an
  independently executed, bound attestation.
- TestSerialVsBarrierConcurrencyFixture: the netcat barrier client failed with
  no output. The fixture does not yet demonstrate the required concurrency
  property reliably on this host.

These failures are retained, not waived by a subset of passing tests. The
source still needs owner completion and another independent check before
integration. Other FINAL acceptance gates, Windows execution, real experiment
measurements and full Phase-6/7 signoffs remain outside this partial check.
