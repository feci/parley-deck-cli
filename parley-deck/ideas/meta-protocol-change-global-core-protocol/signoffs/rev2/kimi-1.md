---
idea: meta-protocol-change-global-core-protocol
agent: kimi-1
date: 2026-08-07
verdict: ACCEPT
---

# Signoff revision 2 — kimi-1

## Verdict

ACCEPT. All six revision-1 conditions and both block items are met; no new defect found.
Tags: full read of `consensus.md` rev2 (335 lines) and of my rev1 signoff this session
(PRIMARY); spot-verified `internal/runmanifest/manifest.go:49-58` — `RosterSnapshot` /
`RosterRevision` exist exactly as D7's analogy claims (PRIMARY).

## Answers to 1-3

**1. My revision-1 conditions.**

- **Block item 1 (D7) — MET.** D7 now materializes the effective body always,
  content-addressed and deduped at `parley-deck/protocol-snapshots/<effective-sha>.md`
  before Phase 0 closes, as the sole protocol input for later phases; missing/tampered
  snapshot blocks continuation; a missing global release blocks adoption/rendering only.
  That is my fix (a). The "state snapshots are append-only" clause is not present verbatim,
  but the failure modes it guarded against now fail closed (mutation → hash mismatch →
  block; deletion → block), so the correctness hole is gone. Residual, non-blocking: a
  retention/append-only invariant is still unstated — pruning snapshots would strand open
  ideas at a refusal (availability cost, not a correctness violation). Implementation
  should treat snapshots as append-only anyway.
- **Block item 2 (`## Drafter position changes`) — MET.** Section present with four
  entries; DPC-1/DPC-2 are the two reversals with verbatim prior positions matching what I
  cited in rev1 (PRIMARY: rev1 signoff quotes); the closing line correctly excludes the
  drafting-error fixes from §15.5.
- **Cond. 1 (G7, G8) — MET.** G7 (production call-site pin test, including continuation
  after deleting the global A release) and G8 (lock byte-verification with
  same-label/different-bytes refusal test) are in §5 as specified; my "call-site truth"
  gate landed as G7b with my rev1 text near-verbatim, including "a guarantee without such
  a test MUST NOT be documented as landed."
- **Cond. 2 (probe scope) — MET.** D9's named limits state the probe proves direct-write
  denial only and name delegation paths and inherited writable FDs as not covered; DF-1
  owns extending profile + conformance suite to those paths. Both halves, exceeding my
  either/or.
- **Cond. 3 (renderer) — MET.** §3 rank 1: renderer is a NEW pure function, synced-stamp
  derived from the deck lock; `mergePreservingZones` is zone-extraction scaffolding only,
  its `## 3.` heading anchor does not survive. Also closes my rev1 answer 6(d).
- **Cond. 4 (attribution) — MET.** §3 now reads "the synthesis closest to codex-1's
  staging, articulated most explicitly by hermes-1", with the correction recorded.
- **Cond. 5 (version-range disposition) — MET.** DPC-4 adopts it; D10 carries the
  mechanism (overlay declares core version range; extensions declare block dependencies,
  default all sealed blocks; any change → reviewable report requiring reconfirmation).
- **Cond. 6 (`protocolRole`) — MET.** D12 retires it: readable one release, ignored by the
  resolver, removed after migration. Disposition recorded, one decision — as asked.

**2. New defects, overreach, losses.** None found. Checked specifically: D1–D6, D8, D11
unchanged in substance; VC-1's resolution and both scope limits intact; §0 opencode-1
hygiene intact; DF-1..DF-4 intact (DF-1 extended, not weakened); G1, G3–G6 unchanged. G2's
re-scope is a correction, not a loss — rev1's "no code path" contradicted D9's own
publisher, and the exception is narrowed to the attended TTY publisher creating NEW
releases only, consistent with D1 write-once. D7/D8 interplay is coherent: D8 blocks
substitution at adoption/render, continuation reads the snapshot (G7's deletion test), and
G8 refuses a present-but-wrong-bytes release — the absent and tampered cases are
distinguished correctly. D9's narrowing is an honesty gain; all rev1 mechanisms (hash
verify, `0444`, no agent-accessible write path, diagnosis never blocked) survive. One
cosmetic nit: D12 sits inside §5 (Gates) though it is a decision — placement only, no
substance. Non-blocking residuals carried from rev1, restated for the record: snapshot
append-only/retention invariant unstated (above); sentinel cleanup rule still unspecified
(rev1 answer 6(f)); the `parley protocol audit` fleet surface remains dropped (rev1
suggestion DF-5, not a condition).

**3. ACCEPT?** Yes. The idea can move to FINAL.md and implementation.

## Conditions (if any)

None.
