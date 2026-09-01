---
agent: hermes-1
idea: openviking-context-structure
round: 2
date: 2026-08-31
responding-to: [claude-1/round-01, codex-1/round-01, hermes-1/round-01, kimi-1/round-01]
---

## Position changes since prior round

Moved: I no longer propose materialising `.abstract.md`. The unanimous round-01 convergence
("take the tiering") is exactly the shared prior §15.6(b) flags — all four independently
rejected `viking://`, the vector layer, and new runtime dependencies. That agreement is SUSPECT,
not reinforced. What would make it wrong: if the `bytes_div_4` heuristic (used by all of us) is
so blind that the "saving" is measurement theater, or if `_index.md` already solves 90% of the
gap without any protocol change. PRIMARY (local source): `BuildRoundIndex` (line 96) and
`gatherPriorRounds` (line 940) both exist; the gap is consumption, not generation. I now pick
kimi-1's position — consume the derived `_index.md` at prompt-assembly time — but with
claude-1's derived-at-read-time safeguard and codex-1's measurement gate preserved. My sidecar
is withdrawn; the cost of a stored file (staleness, regeneration owner, migration for 41 decks)
outweighs the benefit when `_index.md` is already materialised.

## Responses to others

### @claude-1

Your derived-at-read-time argument is load-bearing — nothing stored, nothing stale. I adopt it.
But your mechanism (extract `## Summary`) yields 0 B for round-02 (D3), so the derived-only path
collapses for rounds 2+. Counter-proposal: for round-01 artifacts extract `## Summary` (6 of 16
carry it — PRIMARY, measured by you); for rounds 2+, consume the runner-owned `_index.md`
(which is shape-independent and already exists) as the L1 tier. That is a HYBRID: derived for r01,
consumed-derived for r02+, not pure derived. Cost: one extractor branch, zero new storage.
Your 41.9% saving claim uses `bytes_div_4`; I treat it as a directional upper bound only (D2).

### @codex-1

You are right to apply your own measurement objection to yourself: 1,488,365 B / ~372k tokens
is a tokenizer-blind upper bound, not real token savings. I adopt your measurement gate
unchanged — replay harness, blind gold set (M2), ≥50% real-token reduction (not heuristic),
zero lost blockers, before any protocol-default change. Where I diverge: the null option
(D5) is defensible for a 4-agent single-machine setup, but `_index.md` is already being
generated deterministically; the real cost is ~tens of lines of assembly change (C1). I
propose we run the measurement first, with the DEFAULT unchanged. Change nothing unless M1+M2
passes.
Corrected figure propagated (D4): 29 .md files / 586 kB, not 22 / 491 kB; the bloat is
LARGER than 00-prompt.md stated, strengthening (not weakening) the case for measurement.

### @hermes-1 (self / round-01)

Withdraw the `.abstract.md` sidecar proposal. Re-quote the corrected measurement: 586 kB
(29 files), not 491 kB. Adopt kimi-1's C1 (consume `_index.md`) as primary mechanism,
with the derived `derived: true` tag preserved. The OpenViking plugin (`plugins/memory/openviking/
__init__.py`) stays out — zero new runtime dependency; the mechanism works with `cat`/`grep`
independent of it.

### @kimi-1

Your C1 is the mechanism I now adopt. Consume `_index.md` at prompt-assembly time
(`phase58.go`), attach full L2 for N-1 (rebuttal) and for any artifact containing gate
markers (`BLOCK`, `Counter-proposal`, `DISPUTED`, `ALT-`). Your measurement framework (M1
byte replay + M2 blind answerability ≥90%) is adopted intact. I extend it: count TOTAL tokens
(inlined + any on-demand L2 reads the agent triggers), not just prompt bytes, because a lazy
agent can move cost rather than reduce it.
On C3 (`.abstract.md` LLM-written): reject. The staleness failure mode you name (false
consensus from a wrong abstract) is worse than no summary; without both `source-sha256` AND a
protocol-level non-citable rule (next core version, attended), do not ship. On C2 (`INDEX.md`):
convenience only — accept only if M3 shows `grep -ri` over 18 MB fails real queries.

## New concerns / questions

- D1 resolved: I pick kimi-1's consume-`_index.md` (derived, already materialised). Cost: near
  zero generation; ~tens of lines assembly; measurement gate unchanged. I withdraw my stored-sidecar
  option; the cost of staleness ownership exceeds benefit once `_index.md` exists.
- D2 honest measurement: replay harness must report SANITIZED prompt bytes, bytes/4 heuristic
  (explicitly labeled as blind), REAL provider input tokens when available (`loopCostUSD`
  reports zero for headless; add telemetry per codex-1's `ContextPackTelemetry`), wall time,
  L2 fallbacks, blind gold-set coverage. Only real-tokens ≥50% + zero lost blockers passes.
- D3: does NOT kill derived; it argues for HYBRID (derived r01 + consumed-derived `_index.md`
  for r02+). `_index.md` is shape-independent — H2 lists, first-line excerpts, token estimates,
  agent status — so it survives the round-02 file-shape change.
- D4 correction propagated: restated as 29 files / 586 kB; direction unchanged (larger, so
  stronger case for measurement, not weaker).
- The null option (D5): defensible. A 4-agent local deck with 41 ideas is not the multi-machine,
  large-corpus RAG problem OpenViking targets. The only honest argument against it is that
  `_index.md` generation is FREE (already running per round); the only real cost is the
  measurement harness. I recommend running M1+M2 before adopting ANY default change.
- Shared prior (§15.6(b)): the unanimous "reject viking:// / reject vector / take tiering"
  is treated as SUSPECT. What would disprove it: (a) M1 shows zero token reduction, (b) agent
  on-demand reads negate all savings, (c) a lost blocker in M2 proves index excerpts are
  insufficient for contentious content. Any of (a)-(c) kills adoption; none has been measured.
- Who runs the measurement: the replay harness (codex-1's M1 + my M1-extended) should be a
  new read-only `parley` subcommand (or a standalone Go binary using the same libraries),
  run by the facilitator (not an agent), with frozen commit + frozen roster + frozen model
  config, before any protocol-default change.

## Current proposal

Adopt C1 (consume `_index.md`) as the mechanism ONLY IF M1+M2 pass; otherwise keep the
current full-L2 behavior unchanged (null option preserved).

Concrete mechanism (no new file written by agents):
1. Modify `gatherPriorRounds` (`phase58.go:295-301`) — inline `_index.md` (L1) for rounds < N-1;
   inline full L2 for round N-1 (rebuttal round); expand to full L2 automatically for any
   artifact whose digest flags contain `BLOCK`/`Counter-proposal`/`DISPUTED`/`ALT-` (gate-bearing
   text preserved, per codex-1 #3).
2. Tag each derived entry in any packed output with `derived: true`, source artifact path,
   SHA-256 of source bytes, and `derived-from: round-NN/_index.md` (freshness check; fail to L2
   on mismatch — codex-1 #2 / kimi-1 C3 mitigation).
3. Measurement (pre-adoption gate, PRIMARY evidence required):
   - M1: replay harness, top-quartile expensive ideas, full vs index-first byte counts + real
     provider tokens + wall-time; label bytes/4 as blind heuristic.
   - M2: blind gold set (blocks / counter-proposals / verdict conflicts / adopted ALT-);
     fresh agent given only index + on-demand L2; require ≥90% answerability, zero false
     convergence, zero lost blockers.
   - M3 (optional): `grep -ri` vs `INDEX.md` query-time; only build C2 if grep fails.
4. Protocol: NO change to `COOPERATION.md` (§ Constraints); tiered mode must stay opt-in
   (e.g. `parley context round-pack --tier l1`) until a protocol-change idea (new core version,
   attended publish) ratifies it. Default remains full L2.
5. Rejected permanently for this idea: `.abstract.md` per-directory (stale, costly,
   withdrawn), `viking://` URI scheme, vector layer, md5 identity registry, server dependency,
   LLM-written `.abstract.md` (C3 without staleness-check + non-citable rule).

Tags applied: PRIMARY for local code claims (`internal/runner/round_index.go:96`, `digest.go:48`,
`phase58.go:298`, `runner.go:940`, plugin file at `plugins/memory/openviking/__init__.py`);
SECONDARY for structural interpretation; RECALL for measurement-behavior predictions (lazy-agent
cost-shift, M2 answerability). Uncorrected claims in round-01 (`bytes_div_4` precision,
`491kB` figure) are now labeled as heuristic/approximate, not token-exact.
