---
idea: meta-protocol-change-devx-speed
drafted-by: claude-1
date: 2026-07-03
---

## Agreed decisions (unanimous across all four round-01 analyses)

All four participants (claude-1, codex-1, hermes-1, antigravity-1) independently reached
the same diagnosis and the same core remedy. The convergence was strong enough that no
cross-review round was needed to settle the substance; the only open points are a small set
of design choices flagged below for the owner.

### D1 — Conditional rigor via a `track:` field is the central change
Replace the single fixed 9-phase lifecycle with **risk-tiered tracks**, selected by an
**objective, script-checkable classifier** and recorded as `track:` in `00-prompt.md`
frontmatter (§4 Phase 0). This is the structural enabler for every other speed lever. All
four proposed this independently. (claude-1 §2, codex-1 §2, hermes-1 §2, antigravity-1 §2.)

### D2 — Track shape: 3 Parley tracks + 1 explicit "don't use Parley" off-ramp
- **`fast`** — small, reversible, no security/data surface (≲3–5 files). Collapse Phase 0–4
  into one `FINAL.md` with embedded signoffs; **single model-diverse reviewer**; no cross-
  review rounds; fix-up capped at 1 cycle; short timeouts.
- **`standard`** (default) — ordinary features/refactors, reversible. Parallel round-1;
  cross-review **capped at 2 rounds**; **2 reviewers** (not full roster); simultaneous
  FINAL draft; fix-up capped at 2 cycles.
- **`deliberation`** — today's full protocol **unchanged**. Forced by any high-risk trigger.
- **Off-ramp (codex-1's "Solo Scratch" / claude-1's "Direct"):** trivial reversible work
  (typo, docs, one-file rename) should **not invoke Parley at all** — do it normally and do
  not claim Parley verification. This is stated explicitly as "when NOT to use Parley,"
  which cleanly avoids stretching §1's non-solo rule with solo exceptions for tiny work.

### D3 — Objective triggers; `deliberation` is forced, never opt-down
`deliberation` is mandatory (no down-tiering) when ANY of: protocol change (§7);
security/auth/secrets/payments/privacy/production-infra; data migration / irreversible /
destructive; `strict_gate: true`; `auto_implement`; pipeline/action block (§12); public-API
or persisted-schema break; or change size above the standard ceiling. Default track =
`standard`. The author may always **upgrade**; downgrading below the classifier floor needs
a recorded user OK. (All four; codex-1 §2.D, hermes-1 §6.8.)

### D4 — Track is binding but challengeable; a safety valve prevents under-tiering
Any participant may force-upgrade to `deliberation` (inbox note before round-1 closes), and
a reviewer may force `fast → standard` by filing a MAJOR/CRITICAL finding that cites a
trigger. A **mid-idea upgrade path** (risk discovered during implementation → re-run the
current phase under the stricter track) is required. (hermes-1 §2 + OQ2, codex-1 §2.)

### D5 — Document restructure: progressive disclosure (biggest DevX win)
Restructure `COOPERATION.md` into a **~150-line core** (quickstart + track table + phases-
by-track + quorum + conflict rules) with **§9, §11, §12, §13, §14 moved to clearly-marked
appendices**. Move an expanded TL;DR/quickstart to the top (~line 20), add a role-based
"Who are you? → read this" table, and **consolidate the LE-1…LE-11 jargon** into plain
English in one place (today "LE-1" appears inline with no context). No content is deleted —
only reordered behind progressive-disclosure pointers. (hermes-1 §3, codex-1 §3, all four.)

### D6 — Speed levers (all preserve the safety core)
- **Reviewer scaling:** fast=1, standard=2, deliberation=all — non-solo + refutation-default
  kept on every track (only redundant coverage is dropped).
- **Collapse Phase 3+4** for fast (embedded signoffs) and simultaneous draft for standard.
- **Tiered timeouts** 5 / 15 / 30 min per track (protocol recommends; skill seeds
  `agents.toml`) — attacks the single largest wall-clock sink (uniform 30 min today).
- **Cap cross-review at 2 rounds** for standard (today unbounded, §4 Phase 2).
- **§9.0 readiness ping:** skip for fast; cached-TTL liveness for standard; full for delib.
- **Track-aware auto-advance driver** and **fix-only verification** for narrow fix-ups.
- **Parallel/streaming rounds** for standard (open round N+1 on available files).

### D7 — Modern agentic concepts, mapped (not bolted on)
Conditional rigor / right-altitude → the tiering model. Deterministic workflow vs model-
driven → track selection is deterministic routing; round content stays model-driven.
Progressive disclosure / context engineering → the doc restructure. Spec-driven dev,
plan mode, lead+subagent orchestration, parallel worktrees, refutation gates, right-sized
autonomy, "closing the loop," the bitter lesson → **mostly already present** (FINAL-as-spec,
§1 helpers, §11.B sub-branches, LE-1, §14, LE-7/11); the work is to make them **track-aware**
rather than to invent new machinery. (codex-1 §5, hermes-1 §5.)

## Agreed trade-offs
- A `fast` track trades redundant multi-reviewer coverage for speed; mitigated by mandatory
  refutation-default even at 1 reviewer, and by irreversible work being forced to `deliberation`.
- Restructuring shifts section numbers and changes `protocolSha256` — a one-time consumer
  freshness sync (§9.0) and a cross-reference audit are required.
- Net core length must stay disciplined: the new tiering rules must fit in ~40 lines or the
  restructure's DevX gain is eaten by the amendment (hermes-1 Risk 5).

## Open items for the owner (escalated — these are product-level calls)
1. **Track count/shape:** 3 tracks + off-ramp (this draft) vs codex-1's 4 explicit tracks
   (Solo Scratch as a named track). Substance identical; only whether the off-ramp is a
   named track or a "don't use Parley" note.
2. **How aggressive on speed vs safety** for the default `standard` track (e.g. is auto-
   advance-with-human-gate-only-at-FINAL acceptable, or keep human gates?).
3. **Scope of this idea:** deliver the ratified design (proposal) only, or proceed to
   implement the rewrite (COOPERATION.md + embedded default + skill fallback) now — a
   `deliberation`-track change to the protocol itself.
4. **Two-participant floor** (hermes-1 OQ5): standard with 2 participants → 1 reviewer.
5. **§9 checklist home:** short core block vs Appendix F (hermes-1 OQ3).

## Comparison & blind spots
- **Unanimous & high-confidence:** the tiering diagnosis, objective triggers, doc
  restructure, reviewer scaling, tiered timeouts, "what MUST stay."
- **Only hermes-1 covered deeply:** track-gaming culture drift, mid-idea upgrade mechanics,
  two-participant roster floor, restructuring/cross-reference risk, timeout enforcement
  location (protocol vs skill). These are the real implementation risks.
- **Only codex-1 covered:** an explicit machine-readable track schema + reducing §11 PR/MR
  mirror bookkeeping via tooling.
- **Blind spot (no participant addressed):** migration/back-compat for the ~60 existing
  ideas under `ideas/` — do they get a retroactive `track:`, or is the field prospective
  only? Should be settled in FINAL.

## Signoffs

<!-- Owner resolved Open items 1–3 on 2026-07-03: "choď, implementuj to podľa odporúčaných
     defaultov" → 3 tracks + off-ramp; balanced standard-track speed (human gate at
     FINAL→implementation). Each block below was authored by that participant in
     signoffs/<agent-id>.md and assembled here verbatim (concurrent-append avoidance). -->

### Signoff: codex-1 — 2026-07-03
Status: ✅ ACCEPT
Notes: None.

### Signoff: hermes-1 — 2026-07-03
Status: ✅ ACCEPT

### Signoff: antigravity-1 — 2026-07-03
Status: ✅ ACCEPT
Notes: The tiered rigor model and documentation restructure faithfully reflect my round-01 position. The balanced auto-advance posture for the standard track provides the necessary speed-up while maintaining the human safety brake where it matters most.

### Signoff: claude-1 (implementer/facilitator) — 2026-07-03
Status: ✅ ACCEPT — design is unanimous (✅ ×4, zero reservations). Owner approved the
recommended defaults. Proceeding to Phase 5 implementation (already drafted) → Phase 6 review.
