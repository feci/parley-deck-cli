---
idea: openviking-context-structure
drafted-by: claude-1
date: 2026-08-31
rounds: 2
participants: [claude-1, codex-1, hermes-1, kimi-1]
status: awaiting-signoffs
---

## Question asked

Is the way OpenViking structures data interesting enough to inspire Parley Deck, and if so, what
should we take?

## Verdict

**Yes as a hypothesis generator. No as a system, a dependency, or a file format.**

Exactly one OpenViking idea survives contact with this repository: **staged reading** — navigate on
a compact view, expand to full text on demand. Everything else is rejected. And the idea that
survives turns out to be **something Parley Deck already implements twice and never uses**, so the
work is not adoption but activation — and even that is gated behind a measurement we cannot yet run.

**The production decision today is: change nothing.**

## What we take

**The principle only:** L0/L1 navigation with L2 (full artifact) on demand.

## What we reject, and why

| Rejected | Decisive reason |
| --- | --- |
| `viking://` URI scheme | git paths + blob SHAs already address artifacts; a new scheme breaks `cat`/`grep`/`git` and needs 41-deck migration |
| md5 logical identity | kimi-1, unrefuted: the uri is an input to `md5(account_id:uri)`, so it is **not** rename-stable — it does not deliver the property it was wanted for |
| Vector store / semantic index | new runtime dependency; no measured query workload exists to justify it |
| OpenViking server or client | hermes-1 ships `OpenVikingMemoryProvider`, but availability is not a reason to adopt; it needs a server we will not run |
| Stored `.abstract.md` / `.overview.md` sidecars | withdrawn by their own proposer (hermes-1) in round 2: ~10–50 tokens saved per agent per round against a 220,115 B prompt, plus a staleness owner |
| Cross-deck recall index (41 decks) | speculative infrastructure; kimi-1's M3 (`grep -ri` vs index on 10 real questions) must run first and may show it is convenience, not capability |

## What we found in our own code (all PRIMARY, verified 2026-08-31)

This is the substantive result of the idea.

1. `internal/runner/round_index.go:96` `BuildRoundIndex` writes `round-NN/_index.md` — a per-round
   L1 view: per-agent status, approx-token counts, H2 section lists, first-sentence excerpt per
   section. Deterministic, `derived: true`, **zero model tokens**. 2,220 B for a 2-agent round.
2. `internal/driver/digest.go:48` `BuildRoundDigest` is a **second** deterministic extractor.
3. `internal/runner/runner.go:940` `gatherPriorRounds` loops `for r := 1; r < round` and inlines
   every earlier `.md` in full, **explicitly skipping `_index.md`**. Generated, then deliberately
   discarded.
4. Each extractor has exactly **one** call site — `runner.go:238` and `driver.go:458` — both on the
   `parley run` path. Under manual facilitation neither runs.
5. Consequently only **three** `_index.md` files exist across the whole deck, all June 2026, all
   `round-01`. This very idea produced none.
6. `internal/driver/loop.go:175`, verbatim: *"the runners do not yet emit `agent.usage`, so this is
   0 in practice until that telemetry lands."* **There is no token telemetry.**

## Measured cost of the status quo

On `protocol-generation-bias`: a round-03 participant prompt carries **220,115 B** (round-01 99,848
+ round-02 120,267), times 4 participants. Whole idea: 29 `.md` files, 599,718 B.

**These are bytes, not tokens, and no participant may restate them as tokens.** All three
byte-derived savings figures from round 1 were withdrawn by their own authors: claude-1's "41.9%
saving" (bytes, not tokens), codex-1's "~372k tokens" translation of 1,488,365 B, and kimi-1's
614,718 B / 30 files (wrong denominator — it counted a non-`.md` file; the `.md` figure is 29 /
599,718 B).

## Agreed decisions

**A1. Production behaviour does not change.** `gatherPriorRounds` keeps full-inline. No new file,
scheme, index, service or dependency ships.

**A2. Derive in memory; store nothing.** If a compact view is built, it is derived from the current
participant artifacts at prompt-build time and discarded. The stored `_index.md` is **not** read as
a source of prompt context. All four participants converged here after self-correction: hermes-1
withdrew the stored sidecar, kimi-1 demoted "consume the stored index", codex-1 withdrew its
hash-bound stored-index extension, claude-1 withdrew `## Summary` extraction as shape-dependent.

*Consequence to record:* because nothing is stored, codex-1's round-01 SHA-256 freshness binding is
**moot by construction** — a derive-at-use design has no staleness window. It is dropped, not
rejected: if any future design stores a derived artifact, the binding becomes mandatory again.

**A3. Full text is never optional for load-bearing content.** Round N-1 always inlines in full
(rebuttal needs exact claims), and any artifact containing `❌ BLOCK`, `Counter-proposal`,
`DISPUTED` or `ALT-` always expands to full text. Any read, parse or empty-outline error falls back
to full text and prints the reason.

**A4. A pre-registered gate, fixed now so it cannot be renegotiated once numbers exist.**
- **M1 — cost.** Replay over the top-quartile most expensive closed ideas. Report sanitized prompt
  bytes, the `bytes_div_4` heuristic **explicitly labelled blind**, provider-reported input tokens,
  wall time, and full-source expansions. Count the **whole interaction**, including on-demand reads.
  Pass at **≥50% median reduction in real provider input tokens**.
- **M2 — safety.** Blind gold set built from prior blocks, counter-proposals, verdict conflicts and
  adopted/rejected `ALT-` entries. A fresh agent given only the compact view plus on-demand reads
  must reach **≥90% answerability, zero lost blockers, zero false convergence**. Any lost blocker
  or adopted alternative is an automatic failure.
- **M3 — recall (optional, later).** `grep -ri` versus a deck index on 10 real "have we deliberated
  X?" questions. If grep answers acceptably, say so and drop the index.

**A5. M1 cannot run today — telemetry is a prerequisite, not a detail.** `agent.usage` is not
emitted by headless runners (finding 6). Prompt/token telemetry must land **before** M1 is
meaningful. Any plan that assumes the measurement is available is wrong.

**A6. A passing gate does not authorise a default flip.** The current round prompt instructs
participants to read every prior artifact, so changing what the runner inlines is a protocol change
— a new core version, published attended by the owner. A separate protocol-change idea is required.

**A7. Size-triggered, never global.** If anything is ever activated, it activates above a measured
size threshold. A two-round `fast` idea must not pay for this.

**A8. Our own unanimity is treated as a shared prior, not evidence (§15.6b).** Four language models
converged on "compress the context" within one round. The disconfirming test is M2, and it is
cheap. If M2 fails, the whole line dies regardless of how good M1's numbers look. The rejection of
OpenViking's URI/vector/runtime pieces is likewise **provisional** and reopens if counterevidence
appears.

## Open disagreement (recorded, not resolved)

**Surface of the measurement instrument.** kimi-1, hermes-1 and claude-1 would ship a read-only
`parley context round-pack … IDEA` (stdout only, writes nothing) as a normal CLI command. codex-1
would keep it behind a benchmark-only flag, on the grounds that a shipped command invites use
before the gate passes. Both agree on the mechanism, the gate and A1; they differ on discoverability.
Consensus does not resolve it — the drafter of `FINAL.md` decides, and the choice is reversible.

## Deferred follow-ups

- `parley-context-telemetry` — emit `agent.usage` (prompt/input tokens) from headless runners.
  Prerequisite for M1. **Owner: unassigned.**
- `index-generation-coverage` — `writeRoundIndex` runs only under `parley run`; manual facilitation
  produces no index. Worth fixing on its own merits (the TUI and `_index.md` navigation value), and
  separable from this idea. **Owner: unassigned.**
- `deck-recall-index` — kimi-1's C2/M3. Parked behind M3.

## Corrections logged during this idea

| # | Correction | Found by |
| --- | --- | --- |
| C1 | `00-prompt.md` said "22 artifacts, ~491 kB"; measured 29 `.md`, 599,718 B. Propagated into hermes-1/round-01 before it was caught. | claude-1 (own error) |
| C2 | "41.9% saving" is bytes, not tokens | codex-1 → claude-1 conceded |
| C3 | "~372k tokens" withdrawn — bytes are not tokenizer output | codex-1 (own) |
| C4 | "30 files / 614,718 B" counted a non-`.md` file | kimi-1 (own) |
| C5 | `gatherPriorRounds` cited at `phase58.go:295-301`; actually `runner.go:940`. Right claim, wrong locator. | kimi-1 (own) |
| C6 | "generation already happens per round / is FREE" — false; only 3 files exist, `parley run` path only | claude-1 → kimi-1 and hermes-1 both revised |
| C7 | `.abstract.md` sidecar withdrawn on its proposer's own cost figures | hermes-1 (own) |

## Signoffs

<!-- Each participant appends its own block below. Do not edit another agent's block. -->

### codex-1
- status: reservations
- date: 2026-08-31
- notes: I accept the null production decision, the derive-at-use correction, and the round-pack surface split; those represent my final position fairly. My source check found two factual overstatements in findings 1 and 4: `BuildRoundIndex` returns content at `round_index.go:96`, while `writeRoundIndex` writes it at lines 83–93 and the direct `BuildRoundIndex` call site is line 89, not `runner.go:238`. Finding 6 should also be narrowed to “headless runners do not emit `agent.usage` or provider input-token telemetry,” because ACP emits `agent.acp.usage` at `internal/runner/acp.go:387–395`; that event still does not satisfy M1.

### hermes-1
- status: accept
- date: 2026-08-31
- notes: My round-01 sidecar proposal (.abstract.md) and the 491 kB / 22-artifact figure are both correctly withdrawn/revised in A1–A8 and the corrections table (C1, C4, C7). The "Open disagreement" section fairly records the surface split: I moved to kimi-1's C1 (consume _index.md) with derived-at-use safeguards, which matches my round-02 self-correction. PRIMARY source check of finding 1 verified: BuildRoundIndex is at line 96 (returns string); writeRoundIndex (line 83) writes it; runner.go:238 is not the direct caller — the draft locator was imprecise, as codex-1 noted. Finding 6 verified from source (runner.go:940+): gatherPriorRounds skips _index.md; headless telemetry is 0; ACP's separate agent.acp.usage event does not satisfy M1. I accept the verdict (no production change, derive-at-use only, pre-registered M1/M2 gate, telemetry prerequisite A5, null option preserved, shared-prior caution A8/§15.6b fairly recorded). No substantive defect; nothing blocks.

### kimi-1
- status: reservations
- date: 2026-08-31
- notes: My round-02 position — C1 demoted from "consume the stored _index.md" to derive-at-pack-build, store nothing — is represented accurately in A2, as are my self-corrections (C4 byte count, C5 locator), the md5 rename-stability rejection, and M3 as the gate for the recall index. The "Open disagreement" section states the round-pack surface split fairly: I proposed a normal read-only `parley context round-pack` command (stdout-only, writes nothing), and I accept that the FINAL.md drafter decides the surface. I independently verified two findings against source (PRIMARY, 2026-08-31): finding 3 — `gatherPriorRounds` at `internal/runner/runner.go:940`, loops `for r := 1; r < round`, explicitly skips `_index.md` (line 953) and inlines every other `.md` in full; finding 6 — `internal/driver/loop.go:174-175`, the quoted comment is verbatim, and A5's "headless runners" scoping is correct. Reservation: finding 1's wording is the same defect class as my own C5 — right claim, wrong attribution. `BuildRoundIndex` (`round_index.go:96`) only returns the string; `writeRoundIndex` (`round_index.go:83-94`) performs the `os.WriteFile` at line 90, and the `runner.go:238` call site belongs to `writeRoundIndex`, as codex-1 and hermes-1 also recorded. No decision in A1–A8 depends on this, so it is a wording correction for the FINAL.md drafter, not a block.

### claude-1
- status: accept
- date: 2026-08-31
- notes: Accepting as drafter, and recording that both reservations are against MY drafting error,
  not against a decision. codex-1 and kimi-1 independently caught that findings 1 and 4 conflate
  `BuildRoundIndex` with `writeRoundIndex`; hermes-1 confirmed it in its own source check. I
  verified it myself (PRIMARY, 2026-08-31): `writeRoundIndex` (`round_index.go:83`) performs the
  `os.WriteFile` at line 90 and is what `runner.go:238` calls; `BuildRoundIndex` (line 96) only
  returns the string. This is exactly the defect class kimi-1 logged against itself as C5 — right
  claim, wrong attribution — and I reproduced it one section after recording it. Also verified
  codex-1's narrowing of finding 6: `internal/runner/acp.go:387-395` does emit `agent.acp.usage`
  with `used`/`size`, so "there is no token telemetry" is too broad; that event reports context-window
  occupancy, not provider input tokens, so A5 stands but its wording must be narrowed. Both
  corrections are carried into FINAL.md as C8 and C9; no decision in A1-A8 changes.
- counter-proposal: n/a
