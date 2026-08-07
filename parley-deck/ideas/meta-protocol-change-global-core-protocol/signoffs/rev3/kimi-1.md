---
idea: meta-protocol-change-global-core-protocol
agent: kimi-1
date: 2026-08-07
verdict: ACCEPT
---
# Signoff revision 3 — kimi-1

## Verdict

ACCEPT.

## Answers to 1-3

1. **Yes — all four fixed, and nothing agreed in rev 1 or rev 2 is broken.**

   - **Fix 1 (G8 vs D7) — FIXED.** `[PRIMARY — consensus.md:283-291]` G8 now byte-compares the
     release against the lock on **adoption** always; on **continuation** only when the release is
     present, and always verifies the snapshot against the recorded effective hash — stating the
     reason itself (D7 requires continuation with the release deleted). `[PRIMARY —
     consensus.md:158-159,167-169]` D7 and D8's body both scope a missing release to blocking
     *adoption and rendering* only, never continuation. This is codex-1's condition 1 verbatim, and
     it stays consistent with G7's release-deletion continuation run `[PRIMARY —
     consensus.md:275-276]`.
   - **Fix 2 (D4 provenance) — FIXED.** `[PRIMARY — consensus.md:111-114]` A replace records the
     expected target-block hash; an extend "replaces nothing" and records D10's dependency set and
     hashes. Matches D10's split `[PRIMARY — consensus.md:199-203]` and leaves G4's fail-closed
     hash rule binding on the replace, where it has an answer.
   - **Fix 3 (§1 attribution) — FIXED.** `[PRIMARY — consensus.md:73-77]` §1's adopted wording is
     now detection for the facilitator/unmanaged agents with attribution only through the attended
     publisher, and `DETECTED-UNATTRIBUTED` for an unexplained mismatch — identical in substance to
     D9 `[PRIMARY — consensus.md:180-185]` and DPC-3 `[PRIMARY — consensus.md:332-336]`. The
     leftover codex-1 caught is gone.
   - **Fix 4 (confinement probe) — FIXED.** `[PRIMARY — consensus.md:175-179]` The sentinel lives
     in a dedicated non-release location (`~/.parley/protocol/.probe/`), probing never modifies
     release bytes, and the trusted launcher cleans it up. The location is outside
     `~/.parley/protocol/core/`, so it does not collide with G2's write ban `[PRIMARY —
     consensus.md:257-260]`.

2. **No new contradiction introduced.** `[PRIMARY — full read of consensus.md against the four
   edits]` I checked the seams the edits touch: the probe path sits outside `core/` (no G2
   conflict); D4's "changed base hash fails closed" now binds only to the replace, since an
   extend's dependency-hash changes route through D10's reviewable report/reconfirmation — the two
   mechanisms partition cleanly; G8's continuation release check is a drift check, not use of the
   release as protocol input, so D7's "sole protocol input" stands; and G8's "verify only when
   present" is exactly what G7's deleted-release test needs. Non-blocking nit `[SECONDARY —
   consensus.md:164]`: D8's heading still says "a missing pinned release BLOCKS" unqualified while
   the body scopes it — heading shorthand, the normative sentences are unambiguous, not a
   condition. `[RECALL]` Rev 1/rev 2 agreements I signed (always-materialized snapshot, G7 call-site
   coverage, G7b call-site truth, D12, ranking attribution) are all still present and unchanged.

3. **ACCEPT.** The four retained contradictions are resolved with no new ones; the idea can move
   to `FINAL.md` and implementation.

## Conditions (if any)

None.
