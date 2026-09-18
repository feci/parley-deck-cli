---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-11
artifact-kind: implementation dispositions
not-a-signoff: true
---

# Budget and signoff D-series corrections

These are implementer corrections to Claude's independent source review of
3ea8693, not independent acceptance of this later source or whole-idea completion.
The owner-authored review is now preserved unchanged as
claude-1-budget-identity-review-20260911.md. The original review process completed
but its artifact write failed; native publication was recovered separately.

## Findings

- D-MAJOR-1: Windows LockFileEx and UnlockFileEx use the same byte at offset
  1 MiB, beyond every bounded identity read. This removes the observed source
  overlap. The claimed Windows runtime consequence remains unverified; a
  Windows cross-build is not runtime evidence.
- D-MINOR-1: read the bounded immutable origin header before consulting a
  possibly missing local identity. Version, hostname and relocation diagnostics
  now precede missing-lock errors. The full token/descriptor checks still run
  before permission. Negative tests remove the local identity and verify each
  mismatch reason, no recreation and unchanged charges.
- D-MINOR-2: unknown actual cost retains the conservative reservation. Exposure
  uses max(original reserve, latest reconciliation) while actual remains unknown;
  a new reconciliation below that reserve refuses. Tests cover retained null
  observation, remaining allowance, exact cap, larger ceiling, rejection without
  mutation and an older low persisted ceiling that cannot undercut the reserve.
  The previous test expectation of unknown exposure was deliberately replaced
  by cap exhaustion: the recorded actual price remains null.
- D-MINOR-3: append validity is separate from verdict. A valid BLOCK emits
  agent.signoff.block-recorded with artifact hash, raw/canonical status and
  process_failed, then still fails the request. Actual CLI fixtures cover
  BLOCK with exit 0 and 7, preserved blocked triage, and no false ACCEPT event.
- D-MINOR-4: unsupported-platform refusal now states that a ledger originating
  there has no supported attended recovery/migration. The platform seam is
  tested on this host; refusal preserves charges and inspection remains usable.
- D-MINOR-5: Inspect retries only ErrSnapshotChanged for four further reads.
  After five attempts it returns the distinct retryable error. Malformed data
  is not retried and cancellation interrupts retry. Tests include successful
  recovery, bounded repeated replacement, corruption and cancellation.

## NIT dispositions

1. Preserve the original signoff_status format and add canonical_signoff_status;
   existing event consumers do not receive a silent spelling change.
2. Add failed-process fixtures with forged extra signoff and no append, checking
   that neither positive artifact event is emitted.
3. Every actual signoff request records run.created before launch, including
   all-headless selections. The new BLOCK fixture checks idea/mode identity.
4. Rename the internal publishOrigin helper to publishExclusive, matching its
   origin and identity-file use.
5. Compare wrapped EOF using errors.Is in descriptor identity reads.
6. Document Windows stdin-only attendance and Unix stdin-or-stdout behavior.
7. Keep the bounded two-cap, twelve-process test. Preserve per-child output and
   distinguish process/acquisition failures from cap refusal before evaluating
   admission counts. This improves failure attribution; it is not a fairness or
   starvation guarantee for kernel locking on every shared mount.

## Verification and limits

- Focused budget and app checks passed: budget 1.550s, app 20.095s.
- Budget race checks passed (5.291s), and the actual shared-volume budget suite
  passed (2.239s), before the final test-only diagnostic refinement.
- Windows app test-binary cross-build passed; no Windows runtime run.
- Vet of budget/app/runner passed. Final budget and full app checks are retained
  separately in .parley-runtime/budget-d-series-*.log.

The initial two test failures are preserved in
budget-windows-range-initial-20260911.log. They concerned the intentionally
changed reserve semantics and accidental event status normalization; both are
addressed rather than concealed by dropping the tests.

Claude publication provenance: the exact original denied Write body from
f77101c9-71e8-4252-b7f4-5f661f73b2f4 was mechanically copied to private scratch
for its owner to read. Native path permission attempts and then an OS-sandbox
temporary-file refusal remain failed attempts. The final owner invocation
3b77aee2-9fd1-424c-8cd2-9a05ccd01f27 wrote the canonical artifact through native
Write. Before that attempt, an actual sandboxed atomic-write canary verified
that only the exact output and native temporary siblings were writable and an
unrelated sibling was denied. Codex verified all 27,796 UTF-8 bytes against the
original body, SHA256 37913c7a259b79e495bc6ac46bed277adcb1777349051ec3f8500815714b08dc.
Codex did not author or proxy-write Claude's canonical review.

Shared budget caller wiring, migration, independent production verifier
execution, exact experiments, full review and signatures remain open.
