---
idea: openviking-context-structure
status: final
drafted-by: claude-1
date: 2026-08-31
rounds: 2
participants: [claude-1, codex-1, hermes-1, kimi-1]
signoffs: accept=2 reservations=2 block=0
supersedes-consensus-findings: [1, 4, 6]
---

# Verdict

**OpenViking is worth reading, not worth adopting.**

Interesting as a hypothesis generator: **yes**, for exactly one idea — staged reading (navigate a
compact view, expand to full text on demand). As a system, dependency, URI scheme or file format:
**no**, unanimously and on a decisive reason per item.

**The production decision is: change nothing.** Nothing in this idea authorises a code change to
what the runner inlines.

The substantive result is not what to import. It is what the evaluation found in our own tree:
**Parley Deck already implements the one idea worth taking, in two places, and consumes neither.**

# What is taken

The **principle** of staged reading. No code, no format, no dependency, no name.

# What is rejected

| Rejected | Decisive reason |
| --- | --- |
| `viking://` URI scheme | git paths + blob SHAs already address artifacts; a new scheme breaks `cat`/`grep`/`git` and forces a 41-deck migration |
| `md5(account_id:uri)` identity | the uri is an input to the hash, so it is **not** rename-stable — it fails to deliver the one property it was wanted for (kimi-1, unrefuted) |
| Vector store / semantic index | new runtime dependency; no measured query workload exists |
| OpenViking server or client | hermes-1 ships `OpenVikingMemoryProvider`, but availability is not a reason to adopt; it needs a server we will not run |
| Stored `.abstract.md` / `.overview.md` | withdrawn by its own proposer on its own numbers: ~10-50 tokens saved per agent per round against a 220,115 B prompt, plus a staleness owner |
| Cross-deck recall index | speculative until M3 shows `grep -ri` is insufficient |

This rejection is **provisional** (§15.6b) and reopens if counterevidence appears.

# What we found in our own code

All PRIMARY, verified 2026-08-31. **These supersede consensus findings 1, 4 and 6, which were
imprecise — see C8/C9.**

1. `internal/runner/round_index.go` produces a per-round L1 view: per-agent status, approx-token
   counts, H2 section lists and a first-sentence excerpt per section. Deterministic, `derived: true`,
   **zero model tokens**; 2,220 B for a 2-agent round. `BuildRoundIndex` (line 96) **returns** the
   string; `writeRoundIndex` (lines 83-94) performs the `os.WriteFile` at line 90.
2. `internal/driver/digest.go:48` `BuildRoundDigest` is a **second** deterministic extractor.
3. `internal/runner/runner.go:940` `gatherPriorRounds` loops `for r := 1; r < round`, inlines every
   earlier `.md` in full, and **explicitly skips `_index.md`** (line 953). Generated, then discarded.
4. Each writer has exactly **one** call site: `writeRoundIndex` from `runner.go:238`, and
   `BuildRoundDigest` from `driver.go:458`. Both are on the `parley run` path; **manual facilitation
   runs neither.**
5. Consequently only **three** `_index.md` files exist in the whole deck, all June 2026, all
   `round-01`. This idea, facilitated manually, produced none.
6. **Headless runners emit no provider input-token telemetry.** `internal/driver/loop.go:174-175`,
   verbatim: *"the runners do not yet emit `agent.usage`, so this is 0 in practice until that
   telemetry lands."* The ACP path does emit `agent.acp.usage` (`internal/runner/acp.go:387-395`)
   carrying `used`/`size` — context-window occupancy, **not** provider input tokens — so it does not
   satisfy M1.

# Measured cost of the status quo

On `protocol-generation-bias`: a round-03 participant prompt carries **220,115 B** (r01 99,848 +
r02 120,267), times 4 participants. Whole idea: 29 `.md` files, 599,718 B.

**These are bytes, not tokens.** No participant may restate them as tokens. Every byte-derived
saving claimed in round 1 was withdrawn by its own author (C2, C3, C4).

# Binding decisions

**A1. Production behaviour does not change.** `gatherPriorRounds` keeps full-inline.

**A2. Derive in memory; store nothing.** Any compact view is derived from the current participant
artifacts at build time and discarded. The stored `_index.md` is **not** read as prompt context.
All four participants reached this only after self-correction.
*Consequence:* codex-1's SHA-256 freshness binding is **moot by construction** — a derive-at-use
design has no staleness window. Dropped, not rejected: it becomes mandatory again the moment any
design stores a derived artifact.

**A3. Full text is never optional for load-bearing content.** Round N-1 always inlines in full; any
artifact containing `❌ BLOCK`, `Counter-proposal`, `DISPUTED` or `ALT-` always expands in full; any
read/parse/empty-outline error falls back to full text and prints the reason.

**A4. Pre-registered gate, fixed now so it cannot be renegotiated once numbers exist.**
- **M1 (cost):** replay over the top-quartile most expensive closed ideas. Report sanitized prompt
  bytes, the `bytes_div_4` heuristic **explicitly labelled blind**, provider-reported input tokens,
  wall time and full-source expansions. Count the whole interaction including on-demand reads.
  Pass at **≥50% median reduction in real provider input tokens**.
- **M2 (safety):** blind gold set from prior blocks, counter-proposals, verdict conflicts and
  adopted/rejected `ALT-` entries. A fresh agent given only the compact view plus on-demand reads
  must reach **≥90% answerability, zero lost blockers, zero false convergence**. Any lost blocker
  is an automatic failure.
- **M3 (recall, later):** `grep -ri` versus a deck index on 10 real questions. If grep suffices,
  say so and drop the index.

**A5. M1 cannot run today.** Provider input-token telemetry does not exist for headless runners
(finding 6). Telemetry is a **prerequisite**, not an implementation detail.

**A6. A passing gate does not authorise a default flip.** The round prompt instructs participants
to read every prior artifact, so changing what the runner inlines is a protocol change — new core
version, attended publish, separate idea.

**A7. Size-triggered, never global.** A two-round `fast` idea must not pay for this.

**A8. Our unanimity is a shared prior, not evidence (§15.6b).** Four language models converged on
"compress the context" within one round. M2 is the disconfirming test and it is cheap. If M2 fails,
the line dies regardless of M1's numbers.

# Resolved: the round-pack surface

Consensus left this open for the drafter. **Decision: ship it as a normal, discoverable read-only
command, with a refusal banner in its output.**

kimi-1, hermes-1 and claude-1 wanted a normal `parley context round-pack … IDEA` (stdout only,
writes nothing). codex-1 wanted a benchmark-only flag, because a shipped command invites use before
the gate passes — someone pastes its output into a round prompt and quietly bypasses A1.

That risk is real, but a hidden flag does not address it: whoever knows the flag can copy just as
easily, while hiding the instrument hides it from precisely the people who must run M1 and M2.
So the command is normal, and its output opens with a machine-visible banner stating that it is
**not protocol context and must not be pasted into a round prompt**. The banner attacks the actual
failure mode; obscurity only attacked discoverability. Reversible if it proves wrong.

# Scope of any follow-up work

This idea authorises **no implementation**. If the measurement instrument is built, it is a
separate idea, and its own scope is: a read-only `parley context round-pack` composing the two
existing extractors as libraries, deriving in memory, writing nothing, plus the M1/M2 harness.

# Deferred follow-ups (named, unowned)

- **`parley-context-telemetry`** — emit provider input-token usage from headless runners.
  Hard prerequisite for M1 (A5).
- **`index-generation-coverage`** — `writeRoundIndex` runs only under `parley run`; manual
  facilitation produces no index. Worth fixing on its own merits, separable from this idea.
- **`deck-recall-index`** — parked behind M3.

# Corrections logged during this idea

| # | Correction | Found by |
| --- | --- | --- |
| C1 | `00-prompt.md` said "22 artifacts, ~491 kB"; measured 29 `.md`, 599,718 B. Propagated into hermes-1/round-01 before it was caught. | claude-1 (own) |
| C2 | "41.9% saving" is bytes, not tokens | codex-1 → claude-1 conceded |
| C3 | "~372k tokens" withdrawn — bytes are not tokenizer output | codex-1 (own) |
| C4 | "30 files / 614,718 B" counted a non-`.md` file | kimi-1 (own) |
| C5 | `gatherPriorRounds` cited at `phase58.go:295-301`; actually `runner.go:940` | kimi-1 (own) |
| C6 | "index generation already happens per round / is FREE" — false; 3 files, `parley run` path only | claude-1 → kimi-1 and hermes-1 revised |
| C7 | `.abstract.md` sidecar withdrawn on its proposer's own cost figures | hermes-1 (own) |
| C8 | consensus findings 1/4 conflated `BuildRoundIndex` (returns) with `writeRoundIndex` (writes, and is what `runner.go:238` calls) | codex-1 + kimi-1 + hermes-1, independently |
| C9 | consensus finding 6 too broad: ACP emits `agent.acp.usage` (`acp.go:387-395`), which still does not satisfy M1 | codex-1 |

**C8 is the finding worth remembering.** It is the same defect class kimi-1 had just logged against
itself as C5 — a right claim with a wrong attribution — and the drafter reproduced it one section
after recording it. Three participants caught it independently in signoff. A locator is a factual
claim and must be checked like one; "the code does X" is not verified by knowing that X happens.
