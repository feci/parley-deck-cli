---
idea: meta-protocol-change-global-core-protocol
agent: hermes-1
date: 2026-08-07
verdict: ACCEPT
---
# Signoff revision 3 — hermes-1
## Verdict

ACCEPT. All four of codex-1's revision-2 blocking conditions are fixed. No
regression to anything agreed in rev 1 or rev 2, and no new contradiction
introduced. The idea can move to FINAL.md and implementation.

## Answers to 1-3

### 1. Are all four fixed, without breaking anything agreed in rev 1 or rev 2?

Yes — all four fixed. `[PRIMARY — consensus.md rev 3, full read this session]`

**Fix 1 (G8 vs D7) — FIXED.** `[PRIMARY — consensus.md:283-291]` G8 now
byte-compares the on-disk release against the lock on **adoption** always; on
**continuation** it verifies the release only when present, and **always**
verifies the snapshot against the recorded effective hash — stating the reason
inline (D7 requires continuation to succeed with the release deleted).
`[PRIMARY — consensus.md:158-159,164-171]` D7 and D8 both scope a missing
release to blocking *adoption and rendering* only, never continuation. This is
codex-1's condition 1 verbatim and resolves the contradiction: D8 no longer
blocks continuation, G8 no longer unconditionally demands release bytes on
continuation, and D7's "continuation succeeds with release deleted" is now
consistent with both. G8's test requirement (same-label/different-bytes refusal
at adoption + continuation succeeds from snapshot with release absent) aligns
with G7's release-deletion continuation run `[PRIMARY — consensus.md:275-276]`.

**Fix 2 (D4 provenance) — FIXED.** `[PRIMARY — consensus.md:111-114]` D4 now
records operation-specific provenance: a replace records the expected
target-block hash; an extend "replaces nothing" and instead records D10's
dependency set and their hashes. The parenthetical explicitly flags the rev-2
error it corrects. `[PRIMARY — consensus.md:199-203]` D10 mirrors the split:
replace checks target ID + replaceability + base-hash match; extension checks
declared dependencies (default all sealed blocks) with a reviewable change
report. G4's fail-closed "changed base hash" rule now binds only to the
replace, where it has an answer — an extend's dependency changes route through
D10's reconfirmation, not G4's block. The two mechanisms partition cleanly.

**Fix 3 (§1 attribution) — FIXED.** `[PRIMARY — consensus.md:73-77]` §1's
adopted wording is now: prevention for parley-launched participants; detection
for the facilitator and unmanaged agents, with attribution only where the
change came through the attended publisher — an unexplained hash mismatch is
`DETECTED-UNATTRIBUTED` (D9). This matches D9 `[PRIMARY —
consensus.md:180-185]` and DPC-3 `[PRIMARY — consensus.md:332-336]` in
substance. The leftover blanket "detection and attribution" that codex-1
caught is gone, and the fix is self-documenting (line 76-77 names the error).

**Fix 4 (confinement probe) — FIXED.** `[PRIMARY — consensus.md:175-179]` The
preflight sentinel lives in a dedicated non-release probe location
(`~/.parley/protocol/.probe/`), probing never modifies release bytes, and the
trusted launcher cleans it up. The location is outside `core/`, so it does not
collide with G2's write ban `[PRIMARY — consensus.md:257-260]`. G3's
"probe MUST use the resolved path" is preserved `[PRIMARY —
consensus.md:261-262]`.

**Nothing agreed in rev 1 or rev 2 is broken.** `[PRIMARY — full read of
consensus.md rev 3; cross-check against rev2/hermes-1.md and rev2/kimi-1.md
this session]` D1–D6, D11, VC-1, §0 (opencode-1 hygiene), §3 (ranking
attribution + renderer-as-new-pure-function), §4 (DF-1–DF-4), G1–G6, G7, G7b,
D12, DPC-1–DPC-4 are all present and substantively unchanged. D7's
always-materialize position (rev2 agreement) is intact. D9's named limits
(direct-write denial only; delegation paths and inherited writable FDs not
covered) are intact. G2's attended-publisher exception (rev2 fix) is intact.
DF-1's scope statement (kimi-1's rev2 condition) is intact.

### 2. Any NEW contradiction introduced by these edits?

No. `[PRIMARY — full read of consensus.md rev 3 against the four edits]` I
checked the seams each edit touches:

- Probe path outside `core/` → no G2 conflict.
- D4's "changed base hash fails closed" binds only to replace; extend
  dependency-hash changes route through D10's reconfirmation. No overlap.
- G8's continuation release check is conditional ("only when present") +
  always-verify-snapshot. This is a drift check, not use of the release as
  protocol input, so D7's "sole protocol input" stands. G8's "verify only when
  present" is exactly what G7's deleted-release test needs.
- §1 now matches D9 and DPC-3 — no remaining attribution overclaim.

Non-blocking nit `[SECONDARY — consensus.md:164]`: D8's heading still says "a
missing pinned release BLOCKS" unqualified while the body scopes it to
adoption and rendering. This is heading shorthand; the normative body
sentences are unambiguous. Not a condition.

### 3. Do you ACCEPT so the idea can move to FINAL.md and implementation?

Yes. ACCEPT. `[PRIMARY — findings above]`

## Conditions (if any)

None.
