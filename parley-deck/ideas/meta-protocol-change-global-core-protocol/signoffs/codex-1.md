---
idea: meta-protocol-change-global-core-protocol
agent: codex-1
date: 2026-08-07
verdict: BLOCK
---
# Signoff — codex-1
## Verdict

BLOCK. The core model is sound, but D7/D8 discard the exact-snapshot guarantee I retained in
round 2, D10 does not define a meaningful compatibility check for extensions, and the gates do
not yet protect the real phase-entry call sites. These are implementation-changing defects, not
wording reservations.

## Answers to 1-6

### 1. Treatment of my round positions

No. Stable registry IDs, sealed-by-default resolution, a generated current view, the deck lock,
no empty overlays, fail-closed stale replacements, and launch-path-scoped enforcement were
adopted.

- **Dropped/downgraded:** `[PRIMARY evidence: round-02/codex-1.md:52-55,93-96,211-217]` requires
  one immutable, deck-level snapshot per effective hash and says continuation reads it without
  requiring the global release. D7 makes the body conditional and places it per idea; D8 instead
  blocks when the global release is missing (`consensus.md:127-145`). I do not accept this. It
  defeats self-contained historical interpretation, and “deduped by hash” inside separate idea
  directories is not the deck-level deduplication I proposed.
- **Dropped:** `[PRIMARY evidence: round-02/codex-1.md:168-172,257-260]` requires extension
  dependency reporting and a per-core compatibility receipt. D10 checks only target existence,
  mode, and a base-block hash (`consensus.md:163-168`). That does not notice a changed sealed rule
  which makes an unchanged `ext-1` payload contradictory.
- **Misrepresented:** my position was detection with raw direct writes reported
  `DETECTED-UNATTRIBUTED` unless independent evidence identifies the actor
  (`round-02/codex-1.md:136-141,245-247`). D9 says “detection and attribution” without that bound
  (`consensus.md:153-155`).
- The two-operation overlay replaces my one-verb `provide` syntax. I accept that rejection: the
  registry still owns which targets are replaceable or extensible, so the safety property remains.

### 2. VC-1

I issue no verdict on the managed-prevention claim because I first asserted it and §15.1 forbids
an owner from verdicting its own claim. I accept the evidence-based resolution as the design
position.

`[PRIMARY evidence; executed check, not a verdict]` I attempted the direct, child, unlink, and
symlink-path probes with `/usr/bin/sandbox-exec -p '(version 1) (allow default) (deny file-write*
(subpath (param "GLOBAL")))' ...`. Every invocation stopped before applying the inner profile with
`sandbox-exec: sandbox_apply: Operation not permitted` (exit 71), so this outer-sandboxed session
could not independently reproduce VC-1.

`[SECONDARY — claude-1's PRIMARY check at round-02/claude-1.md:16-49] CONFIRMED:` the canonical
`/private/tmp` target denied direct write, child write, and unlink, while the `/tmp` symlink form
permitted the write. The unresolved-path limit is therefore correctly stated. I also accept the
facilitator exclusion, without self-verdicting the claim I own. The resolution still needs the
other limits from my round 2 (`round-02/codex-1.md:156-166`): inherited writable descriptors and
delegation to an unconfined helper/broker are outside what this filesystem probe proves.

### 3. D3 and user constraint 3

Yes. This is minimal, not tokenistic. `[PRIMARY evidence: consensus.md:85-116]` a deck can replace
the complete working-language rule and can add a deck-namespaced project-rule payload at `ext-1`;
both operations change effective protocol bytes. OOP-style inheritance permits the base to expose
only selected virtual members. The narrow surface is honest only if extension compatibility is
fixed as required below.

### 4. Ranked implementation scope

The order is almost right, but pinning cannot be a separately activatable rank 2. Core
installation/render/check may ship first, but **adoption or changing the current core must remain
disabled until the exact Phase-0 snapshot and all later-phase reads are implemented**. Core hash
verification likewise belongs before first resolution, not in rank 4. Overlay at rank 3 and the OS
sandbox deferral are safe; attended fleet migration may wait. Also, rank 1 tooling merely enables
conversion: it does not itself “convert 36” while migration remains DF-2
(`[PRIMARY evidence: consensus.md:180-199]`).

### 5. Gates

G5 states the right outcome but is too easy to satisfy with a resolver unit test while production
prompt builders still use the current file. `[PRIMARY evidence: internal/app/consensus_request_signoffs.go:718-752]`
the live signoff builder currently emits `filepath.Join(rootAbs, protocol.DeckDir,
"COOPERATION.md")` directly.

Add this binding gate:

> **G7 — production call-site pin test.** Start an idea under effective protocol A, adopt B, then
> capture the actual prompts/inputs produced by the production entry points for round 2, design
> consensus/signoff, implementation, review, review consensus/signoff, fix-up, resume/continue,
> steer, and inspect. Every entry point must resolve and expose snapshot A with the recorded hash;
> none may read the deck-current `COOPERATION.md`, the current core pointer, or core B. Instrument
> filesystem reads (or inject a resolver spy) so any forbidden read fails the test. Run the same
> test after deleting the global A release; continuation must still succeed from snapshot A.

### 6. Collective misses

- G2 says **no code path** may write the core, while D9 requires an attended publisher to create a
  new release (`consensus.md:153-155,209`). Scope G2 to autonomous/agent-accessible paths and make
  the attended, TTY-gated publisher the sole audited exception.
- D4 requires every operation to carry “the expected hash of the core block it replaces,” but an
  extension replaces no block (`consensus.md:98-103`). Define an extension dependency set (or
  conservatively depend on all sealed blocks) and hash it for D10.
- The failed-confinement branch of the sentinel probe writes successfully. Use a dedicated,
  non-release probe directory inside the resolved protected subtree, clean it through the trusted
  launcher, and never probe by mutating released core bytes.

## Conditions (if any)

1. Make the exact effective snapshot mandatory before Phase 0 closes, content-address it once per
   deck, and make it the sole protocol input for every later phase. Missing/tampered snapshot
   blocks continuation; missing global release blocks adoption/rendering only.
2. For every nonempty overlay/core pair, compare declared extension dependencies (all sealed
   blocks by default), issue a reviewable change report, and require reconfirmation when any target
   or dependency changes. Preserve the existing hard blocks for missing/tombstoned/sealed targets.
3. Replace D9's attribution overclaim with: supported attended publication is attributable; an
   unexplained hash mismatch is `DETECTED-UNATTRIBUTED`. Record helper/broker and inherited-FD
   limits on any prevention claim.
4. Add G7 above and repair G2's attended-publisher exception before implementation begins.
