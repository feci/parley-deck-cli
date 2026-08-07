---
idea: meta-protocol-change-global-core-protocol
agent: codex-1
date: 2026-08-07
verdict: BLOCK
---
# Signoff revision 2 — codex-1
## Verdict

BLOCK. Revision 2 contains the requested mechanisms, but retained text contradicts three of my
conditions, and G8 introduces a direct contradiction with the snapshot guarantee.

## Answers to 1-3

1. My conditions are not all met.

   - **Condition 1 — NOT MET.** `[PRIMARY — consensus.md:136-161,258-275]` D7 and G7 make the
     snapshot the sole later-phase input and require continuation after deleting the release, but
     D8 still says a missing pinned release blocks, while G8 requires continuation to byte-compare
     the on-disk release.
   - **Condition 2 — NOT MET.** `[PRIMARY — consensus.md:103-106,188-197]` D10 adds dependency
     declarations, change reports, and reconfirmation, but D4 still requires every extension to
     record the hash of the core block it replaces; an extension replaces no block.
   - **Condition 3 — NOT MET.** `[PRIMARY — consensus.md:69-70,169-178]` D9 is correctly narrowed,
     but §1 still adopts blanket “detection and attribution” for every unmanaged agent.
   - **Condition 4 — MET.** `[PRIMARY — consensus.md:246-265]` G2 has the attended-publisher
     exception and G7 covers the production entry points plus release-deletion continuation.

2. **Yes.** `[PRIMARY — consensus.md:272-275]` G8 is the new defect: as written, missing release
   bytes make its mandatory continuation comparison impossible. `[PRIMARY — consensus.md:165-168,
   250-251; signoffs/codex-1.md:102-104]` The earlier successful-probe location and trusted-cleanup
   requirement also remains unspecified. `[PRIMARY — full comparison of consensus.md with
   signoffs/codex-1.md]` I found no other revision-1 agreement lost.

3. **No; I BLOCK.** The contradictions are implementation-changing and must be removed before
   `FINAL.md`.

## Conditions (if any)

1. Scope D8 to adoption/rendering. In G8, verify release bytes on adoption and, on continuation,
   only when the release is present; always verify the snapshot against the effective hash.
2. Make D4 operation-specific: replace records an expected target hash; extend records D10's
   dependency set/hashes.
3. Replace §1's blanket attribution sentence with D9's attended-only attribution and
   `DETECTED-UNATTRIBUTED` rule.
4. Put the confinement sentinel in a dedicated non-release probe location, clean it through the
   trusted launcher, and never probe by modifying release bytes.
