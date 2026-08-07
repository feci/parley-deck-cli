---
idea: meta-protocol-change-global-core-protocol
agent: hermes-1
date: 2026-08-07
verdict: ACCEPT
---
# Signoff revision 2 — hermes-1

## Verdict

ACCEPT. All three of my revision-1 conditions are met. Revision 2 introduces no
new defect, overreach, or misrepresentation, and nothing agreed in revision 1 was
lost. The idea can move to FINAL.md and implementation.

## Answers to 1-3

### 1. Are each of my revision-1 conditions met?

All three are met [PRIMARY: consensus.md revision 2, full read].

**Condition 1 — `## Drafter position changes` per §15.5.** MET.
`consensus.md:286-330` adds the section with four entries: DPC-1 (prevention
reversal, round 1 → round 2), DPC-2 (snapshot materialization reversal, round 2
→ consensus), DPC-3 (attribution narrowed), DPC-4 (overlay version range and
extension dependencies). Both of the drafter's own reversals are recorded with
prior position quoted verbatim. The closing paragraph
(`consensus.md:328-330`) correctly distinguishes drafting-error corrections from
design-position changes. This is the section my rev-1 condition 1 required, and
it goes further by recording DPC-2 through DPC-4 which I did not specifically ask
for but which §15.5 requires.

**Condition 2 — Add G7 (end-to-end call-site pin test).** MET.
`consensus.md:258-265` adds G7, attributed to codex-1. It is a superset of my
proposed G7: my version tested "the actual §9 read path" against the pinned
hash; the consensus G7 enumerates the production entry points (round 2, design
consensus/signoff, implementation, review, review consensus/signoff, fix-up,
resume/continue, steer, inspect), forbids reads of deck-current
`COOPERATION.md` / current core pointer / core B, and adds the
delete-global-A-release continuation test. G7b (`consensus.md:266-271`, kimi-1)
and G8 (`consensus.md:272-275`, kimi-1) were also added, covering call-site
truth for all named guarantees and lock byte-verification respectively. All
three address the "documented as landed, wrong at the call site" failure class I
flagged in rev-1 Q5.

**Condition 3 — Record `protocolRole: source` disposition.** MET.
`consensus.md:277-284` (D12) retires `protocolRole`, superseded by the deck
lock. It is kept readable for one release for backwards compatibility, ignored
by the resolver, and removed after migration. The disposition is "retire" —
going further than my condition asked (record the disposition) by actually
deciding it, which is appropriate since kimi-1 also raised it as condition 6.

### 2. Does revision 2 introduce any NEW defect, overreach, or misrepresentation?

No [PRIMARY: full read of consensus.md revision 2; cross-check against rev-1
signoffs/hermes-1.md, signoffs/codex-1.md, signoffs/kimi-1.md].

**D7 reversal** (lines 130-153): The new always-materialize position is sound.
The argument — inputs can stop being present (prune, migrate, fresh clone),
making pinned bytes unrecoverable, breaking user constraint 1 — is correct on
its own terms. Dedup-by-effective-hash at a deck-level content-addressed path
means storage is per-distinct-protocol, not per-idea, so the storage objection
does not survive dedup. The "blocks continuation" / "blocks adoption and
rendering only" split is precise. No overreach.

**D9 narrowing** (lines 163-184): Attribution only via attended publisher;
DETECTED-UNATTRIBUTED for unexplained mismatch; named limits (probe proves
direct-write denial only; delegation paths and inherited writable FDs not
covered). This is more honest than rev-1's flat "detection and attribution."
No overreach.

**D10 extension** (lines 186-197): Extensions declare core-block dependencies
(default all sealed blocks), change report requires reconfirmation, overlay
declares core version range. This closes the gap codex-1 identified (rev-1 D10
checked only target existence/mode/base hash, not changed sealed rules). No
overreach.

**G2 scoping** (lines 246-249): Changed from "no code path" to "no autonomous
or agent-accessible code path" with attended TTY-gated publisher as sole
audited exception. The old G2 contradicted D9's own publisher; the fix is
necessary and correct. No overreach.

**§3 attribution correction** (lines 207-209): Now reads "the synthesis closest
to codex-1's staging, articulated most explicitly by hermes-1 and accepted by
all four." This is accurate — codex-1 had pinning 2nd / overlay 4th, which
matches the consensus order; my round-2 order was overlay 2nd / pinning 3rd.
The "articulated most explicitly by hermes-1" is a fair characterization. No
misrepresentation. [RECALL: my rev-1 signoff Q1 flagged the old attribution as
MISATTRIBUTED but not blocking; this fixes it.]

**§3 rank 1 correction** (lines 213-217): Renderer is a NEW pure function;
mergePreservingZones is zone-extraction scaffolding only. This matches my own
round-2 self-correction, which I noted as ADOPTED in rev-1 Q1. No
misrepresentation.

**D12** (lines 277-284): protocolRole retired. Sound — under D1/D8 the deck
lock and release-publishing-repo status carry the source/consumer distinction.
No overreach.

**Nothing agreed in revision 1 was lost.** I verified D1–D6, D8, D11, VC-1, §4
(DF-1 through DF-4), and G1–G6 are all present and substantively unchanged.
D9 and D10 were narrowed/extended, not weakened. DF-1 gained the direct-write
scope statement (kimi-1's condition) — an addition, not a loss. The §15.6(b)
deliberation record I noted in rev-1 Q6 but explicitly did not make a condition
remains absent; that is not a regression since I did not require it.

### 3. Do you ACCEPT so the idea can move to FINAL.md and implementation?

Yes. ACCEPT.

## Conditions (if any)

None. All revision-1 conditions are met, and the revision-2 changes are
surgical — each traces to a specific reviewer condition (codex-1 conds. 1-4,
kimi-1 conds. 1-6) with no collateral damage to agreed positions.
