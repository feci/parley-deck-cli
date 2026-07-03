---
idea: meta-protocol-change-devx-speed
status: final
author: claude-1
consensus-date: 2026-07-03
participants: [claude-1, codex-1, hermes-1, antigravity-1]
---

## Final plan / specification

Introduce **conditional rigor** into the Parley Deck protocol so ceremony scales to risk,
and **restructure `COOPERATION.md`** for developer usability. This is a `deliberation`-track
change to the protocol itself (§7). It is additive-with-restructure: no safety invariant is
removed; the full lifecycle survives intact as the `deliberation` track.

> **Owner-ratified defaults (2026-07-03):** track shape = **3 tracks + off-ramp**;
> standard-track speed posture = **balanced** (auto-advance, human gate only at
> FINAL→implementation). Open items 1–5 in consensus.md resolved accordingly.

### 1. Add a `track:` field (the core change)

`00-prompt.md` gains `track: fast | standard | deliberation` (default `standard`). Track is
chosen by an **objective, script-checkable classifier** at Phase 0 and recorded in frontmatter.

**Classifier (first matching row wins; deliberation triggers are checked first):**

| Force `deliberation` if ANY | Else `fast` if ALL | Else `standard` (default) |
|---|---|---|
| protocol change (§7); security/auth/secrets/payments/privacy/prod-infra; data migration / irreversible / destructive; `strict_gate: true`; `auto_implement`; pipeline/action (§12); public-API or schema break; > ~15 files / ~1000 LOC | reversible; ≤ ~3–5 files / ~300 LOC; no security/data surface; mechanically verifiable (lint/type/test) | everything else |

**Off-ramp (not a track):** trivial reversible work (typo, docs, one-file rename, dep bump
with green tests) should **not invoke Parley at all** — do it normally; do not claim Parley
verification. A new "When NOT to use Parley" note states this so §1's non-solo rule is never
stretched with solo exceptions for tiny work.

### 2. Per-track phase behavior

| Phase | `fast` | `standard` (default) | `deliberation` |
|---|---|---|---|
| 0 Kickoff | lightweight; **no §9.0 ping** (freshness only if `consumer`) | full §9.0; parallel round-1 | full §9.0 |
| 1 Round-1 | single parallel round | parallel round-1 | parallel round-1 |
| 2 Cross-review | **skipped** | **capped at 2 rounds** → escalate/upgrade if still ❌ | unbounded (current) |
| 3+4 Consensus/FINAL | **collapsed**: `FINAL.md` with embedded signoffs | separate but **simultaneous draft** | separate (current) |
| Reviewers (Phase 6) | **1 model-diverse** reviewer | **2** reviewers | **all** non-implementers |
| 7 Review consensus | reviewer ✅ = consensus | reviewers who reviewed sign off | all sign off (current) |
| 8 Fix-up | cap **1** cycle → escalate; **fix-only verification** ok | cap **2** cycles; fix-only verification for narrow fixes | unbounded; full re-review; `strict_gate` available |
| Timeout/agent | **5 min** | **15 min** | **30 min** (current) |
| Auto-advance | full (pauses only for the 1 required signoff) | auto-advance, **human gate at FINAL→impl** | human gate each transition |

**Invariants on every track (unchanged):** ≥1 independent non-facilitator artifact (non-solo,
§1); refutation-default review (LE-1) — reviewer count shrinks, refutation discipline does not;
round-1 independence; append-only ✅/🟡/❌ signoff mechanism; files-canonical audit trail;
§14 human brake; English-only; no-secrets.

### 3. Track binding, challenge & mid-idea upgrade

- Track is **binding** once Phase 0 closes but **challengeable** in round-1.
- **Safety valve:** any participant may force-upgrade to a stricter track via an inbox note
  before round-1 closes; a reviewer may force `fast→standard` with a MAJOR/CRITICAL finding
  citing a trigger. Down-tiering below the classifier floor needs a recorded user OK.
- **Mid-idea upgrade:** if implementation reveals a higher-risk surface (e.g. touches auth),
  any participant force-upgrades; the idea re-runs the **current phase** under the stricter
  track's rules. (Recorded in `00-prompt.md` + an inbox note.)
- **Roster-size floor:** with only 2 participants, `standard`'s "2 reviewers" degrades to 1
  (same as `fast`) — the trigger accounts for roster size, not just risk.

### 4. Document restructure (progressive disclosure)

- **~150-line core, in reading order:** header → expanded **Quickstart/TL;DR (~line 20)** →
  role "Who are you? → read this" table → §1 scope (compressed) → **§4 phases-by-track**
  (one paragraph per phase with the track table above) → §5 quorum → §6 conflict rules → §7.
- **Appendices (content preserved, moved behind pointers):** §9 session checklist, §11
  transport mechanics, §12 pipelines, §13 retro, §14 outer loop.
- **Consolidate LE-1…LE-11** inline jargon into one "Loop-engineering rules" block in plain
  English, cross-referenced from the phases that use them.
- Net-core-length discipline: the tiering rules must fit in ~40 lines (a trigger table + the
  per-track phase table), or the DevX gain is eaten by the amendment.

### 5. Tooling & config (skill + CLI, separate from protocol text)

- `parley init` / `parley run` infer or ask for the track, template `00-prompt.md` with
  `track:`, and seed **per-track timeouts** (5/15/30) into `agents.toml`.
- Track selection is a **deterministic routing decision** (no model judgment); phase
  sequencing stays deterministic; model judgment is reserved for round/review content.
- Both copies of the protocol stay in lockstep (`TestEmbeddedDefaultMatchesLiveDeck`); the
  skill's `references/COOPERATION.md` fallback is re-synced (body-verbatim rule).

## Purpose / user-visible outcome

A developer can start a normal task in <5 minutes reading ~150 lines, and a typical
`standard`/`fast` interaction finishes in a fraction of today's wall-clock (single/dual
reviewer, 5–15 min timeouts, capped rounds, collapsed consensus, auto-advance) — while
high-risk and protocol/security/production work keeps the full deliberation ceremony.

## Context & orientation

Target files (implementation phase, if approved): `parley-deck/COOPERATION.md` +
`internal/protocol/defaults/COOPERATION.md` (edit BOTH — drift guard), classifier + `track:`
handling in the driver/loop code, `parley init`/`run` templating, `agents.toml` seeding,
`parley-deck-skill/references/COOPERATION.md` re-sync, CHANGELOG + version bump + all release
channels. Prior art: fusion-execplans (conditional rigor), strict_gate, LE-1…LE-11, §14.

## Observable acceptance criteria

1. `00-prompt.md` supports `track:`; absent → defaults to `standard`; a script can classify a
   task from objective inputs and any deliberation trigger forces `deliberation` (no down-tier).
2. `fast`/`standard` ideas complete with the reduced reviewer count, capped rounds, collapsed
   consensus, and per-track timeouts described in §2 — verifiable on a sample idea end-to-end.
3. Every track still produces ≥1 independent non-facilitator artifact and refutation-recorded
   review; removing them is rejected by the driver's validation (non-solo + LE-1 preserved).
4. `COOPERATION.md` core is ≤ ~200 lines before the first appendix; §9/§11/§12/§13/§14 are
   appendices; a "When NOT to use Parley" note and a role table exist; LE jargon is consolidated.
5. Both protocol copies remain byte-identical modulo the allowlisted zones
   (`TestEmbeddedDefaultMatchesLiveDeck` green); skill fallback re-synced.
6. `meta/protocol-changelog.md` records the change; `protocolSha256` bump handled per §9.0.

## Idempotence & recovery

Restructure is one atomic edit with a cross-reference audit (section numbers shift). The
`track:` field is additive and back-compatible (absent = `standard`). Existing `ideas/` are
**prospective-only** — no retroactive `track:` backfill (blind-spot resolution). If the
restructure breaks a cross-reference, `parley retro` + the drift test catch it before release.

## Known risks / de-risking

- **Track-gaming / culture drift** → binding objective triggers + escalation valve.
- **Single-reviewer blind spot (fast)** → fast is low-risk/reversible only; irreversible forced to deliberation.
- **Restructuring breaks refs / consumer sync** → atomic edit, ref audit, §9.0 freshness sync.
- **Amendment accretion** → hard ~40-line budget for the tiering section.

## References
- Consensus: ./consensus.md
- Rounds: ./round-01/{claude-1,codex-1,hermes-1,antigravity-1}.md
- Prompt: ./00-prompt.md

<!-- NEXT STEPS (not yet done): collect Phase-3 signoffs from all four participants on the
     settled shape; owner confirms consensus.md Open items 1–5; then implement (deliberation
     track, full review) and release across all channels. -->
