---
agent: kimi-1
idea: meta-protocol-change-phase-packet-and-fixup-budget
round: 1
date: 2026-08-11
---

## Summary

The phase packet is sound only if omission is *visible*, so the design below makes every packet
carry a machine-generated, always-complete index of what it did NOT load. Packets must be
rendered on demand from the single `COOPERATION.md` source (never committed — a committed packet
is a fourth stale copy), keyed on phase × track × idea flags. The instruction change lands in
three text paths (skill standing line, §9 checklist, CLI prompt templates) and is backstopped by
the wall-clock metric the regression idea already established, because hand-written prompts
cannot be governed by text. For the fix-up budget: cap `deliberation` at 5 cycles with a
recorded-escalation extension mechanism, reusing the existing "escalation threshold, not close
criterion" machinery — and explicitly reject any severity floor, because the tail data shows
fresh MAJORs at rounds 19–24.

One honesty note up front, per the brief: a rule missing from both the packet AND the omission
index produces no in-loop signal — under Phase 2 rule 1 that silence reads as agreement.
Detection of that case exists only at generation time (deterministic regeneration + diff) and
ex-post (retro audit). I say so plainly rather than pretend the packet can catch it. The design
therefore forbids hand-curation of the index: inclusion of full text is curated, the index is
not.

## Proposed approach

### Q1 — Section-by-phase mapping (current file: 1,372 lines; PRIMARY, `wc -l parley-deck/COOPERATION.md` → `1372`)

Design rule: every packet = **[always-loaded core]** + **[phase sections]** + **[complete
generated omission index]**. "On-demand" below means: not in the packet, present in the index
with a one-line digest, loadable mid-phase when the digest matches the situation.

**Always-loaded core (every phase, every track — each item is cheap, and each one's omission
converts silence into consent or corrupts the audit trail):**

- Packet header: active `Transport:` value, roster list, idea track, `sourceSha256` of the
  COOPERATION.md rendered from.
- §4.0 track table + invariants (lines 200–257) — it OVERRIDES phase defaults (line 233–237);
  an agent without it applies the wrong ceremony and cannot know.
- §6 rules 1–6 entire (lines 732–744, ~13 lines) — includes rule 3 (never edit another's file)
  and rule 4 (copy snippets; scope-material asymmetry). See per-section verdict below.
- Escalation to user (lines 686–718) + §8 inbox addressing (788–812) — the safety valve, and
  the only path by which a packet-caused problem can be reported mid-phase.
- §15.1 + §15.2 (1239–1290) — verdict ownership and provenance tags bind on every track
  (§15.7 table, 1362–1372) in every phase that writes or challenges a claim, which is all of
  them.
- The phase's own §4 subsection (artifact template + rules — non-negotiable).
- §3 directory layout (164–195) — where files go; needed in every writing phase, ~30 lines.

**Phase-specific inclusion:**

| Phase | Load-bearing (full text in packet) | Reason | On-demand (index only) |
|---|---|---|---|
| 1 | Phase 1 (306–325); §15.4 (1311–1322) | Round 1 is where recommendations first appear; exemption-claim admissibility gates them. The independence rule (324) is in the phase text. | §15.3/15.5/15.6, §5, §7, §14, §12, §13, §11 |
| 2 | Phase 2 (326–353); §15.3 (1291–1309); §15.4; §15.6(a) (1339–1361); §5 whole (719–731) | Verdict challenges are issued and resolved in rounds (15.3). **Non-obvious:** on `standard`, §15.6(a) lives as a section *inside an existing round-02 file* — omit it from the P2 packet and a `standard` idea closes without the adversarial-alternative section and nobody is prompted to notice. §5 is 13 lines and governs mid-round drops/pings. | §7, §12–§14, §11 |
| 3 | Phase 3 (354–397); §15.3; §15.5 (1324–1337); §15.6; §5 | Consensus is where verification integrity bites: DISPUTED claims, `## Verdict conflicts`, `## Drafter position changes`, correlated-agreement close conditions, silent-past-deadline rule (393). | §7, §12–§14, §11 |
| 5 | Phase 5 (434–500); §4.0.1 LE-4/LE-5 text (258–270); §15.2 | IMPLEMENTATION.md is the living contract; validation-evidence claims carry provenance. **Conditional:** §14 full text when `auto_implement: true` or §12 pipeline blocks are present (generator keys off 00-prompt frontmatter) — the human brake binds exactly that context. | §15.3–15.6, §7, §12 (unless flagged), §13, §11 |
| 6 | Phase 6 (501–557, includes LE-1 refutation-default and the review-brief no-suppression rule, 537–556); §15.4; §15.2 | Severity taxonomy, refutation attempts, no-suppression were each bought with a failure and gate what the reviewer writes. | §15.5/15.6, §7, §12–§14, §11 |
| 7 | Phase 7 (558–587); **strict-gate driver block (638–644)**; §15.5; §15.3 | **Non-obvious:** "the Phase 7 review-consensus drafter sets the machine-readable `closing_review_round` and `strict_gate_clean` fields" (`COOPERATION.md:638-641`) — that duty lives in the Phase 8 text but binds the P7 drafter. §15.5: review consensus is the same drafter/signoff shape as consensus.md; include rather than argue scope. | §7, §12–§14, §11 |
| 8 | Phase 8 whole (588–685); LE-4 completion contract | Fix-up format, strict gate, stopping judgment, budget semantics, close integrity all live here. | §15.3–15.6, §7, §12–§14, §11 |

**Reference-only in every packet (index digest, load at need):** §0/§11 transport mechanics
(879–1086, ~208 lines — load the active subsection at publish time; the transport *value* is
already in the packet header); §10 TL;DR (subsumed by the index); §13 retro (facilitator
tooling); Appendix A (adoption); §12 (conditional, see P5); §2 roster annotations (the roster
list is in the header); §9 (superseded by packet loading itself; §9.0 readiness is
facilitator-side at Phase 0).

**Verdicts on the four named sections:**

- **§15** — load-bearing, split by subsection as above: 15.1/15.2 in the always-core; 15.3 in
  P2/P3/P7; 15.4 in P1/P2/P6; 15.5 in P3/P7; 15.6 in P2 (standard) and P3. In P5/P8 only
  15.1/15.2 carry; the rest is on-demand. It was bought with a real failure
  (`meta-protocol-change-verification-integrity`, ratified 2026-08-04, line 1234) — which is
  why the split keeps every subsection reachable and the binding subsections loaded.
- **§7** — on-demand everywhere except Phase 0 of a protocol-change idea (where it is the
  phase's operative text). Reachable-on-demand is sufficient *only because* the index carries
  its one-line prohibition everywhere: "the protocol changes only via an idea (+ user
  ratification for core); never hand-edit COOPERATION.md." The digest carries the norm; the
  full text is needed only to *perform* a change.
- **§6 rule 3** — load-bearing in every phase. It is one sentence inside a 13-line section;
  include §6 entire in the core. Its failure mode (a silent cross-edit) is exactly the kind
  that produces no signal until the audit trail is already corrupt.
- **§14** — conditional load-bearing: full text in P0/P5 packets of ideas flagged
  `auto_implement`, pipeline, or loop tooling; index one-liner everywhere else ("an automated
  loop may only draft candidates; promotion, implementation, merge, roster, and consensus
  changes require a recorded human or full-quorum gate"). The failure it was bought for
  (`automation-outer-loop`, ratified 2026-06-24, line 1230) is an automation-context failure,
  so conditioning on automation flags matches the blast radius instead of taxing every reader.

### Q2 — Packet production and honesty

- **Render, don't store.** `parley protocol packet --phase N` (skill-invokable) renders the
  packet from the live `COOPERATION.md` at invocation time. Nothing is committed. This deck has
  three copies already — deck view, CLI embedded default (`internal/protocol/defaults/COOPERATION.md`,
  go:embed at `internal/protocol/workspace.go:21`), skill packaged copy
  (`.../skills/parley-deck/references/COOPERATION.md`) — plus the drift guard in
  `meta/version.json` (`protocolSha256` f67a4c1e… vs `packagedProtocolSha256` 00c033cf…,
  currently differing; PRIMARY, `cat parley-deck/meta/version.json`). A committed packet would be
  copy number four. A rendered view cannot go stale; the stale-second-copy failure mode is
  eliminated by construction, not guarded against.
- **One curated artifact:** the inclusion map (phase × track × flags → section-id list) lives in
  a single machine-readable table next to the protocol source. It is code-reviewed like code.
  Everything else in the packet is derived.
- **What reports what a packet omits:** every packet ends with the generated **omission index** —
  for every omitted section: heading, one-line digest of its MUST/NEVER sentences, line count.
  The index is built from structural headings (`##`/`###`) plus normative-keyword extraction,
  never hand-curated. Headings are structural — a rule phrased without "MUST" still surfaces
  under its heading, so the index cannot silently lose a section the way a curated one can.
- **Honesty check:** the packet header carries `sourceSha256` of the COOPERATION.md it was
  rendered from plus the generator version. `parley preflight` (already the §9.0 vehicle)
  verifies the hash against the deck's live file; mismatch → fall back to full read and warn.
  CI extends the existing drift-guard pattern: regenerate packets deterministically and diff.
- **Plain statement on detectability:** omission of full *text* is detectable at runtime (index
  digest + on-demand load). Omission from the *index* is detectable only at generation time
  (deterministic regeneration diff) and ex-post (§13 retro sampling a full read). Within a live
  idea there is no signal — accept this and put the enforcement where detection exists.

### Q3 — Where the instruction changes

Confirmed three text paths, plus a fourth in code (PRIMARY, grep):

1. Skill standing line — `.../skills/parley-deck/SKILL.md:12`: "Always read
   `parley-deck/COOPERATION.md` first." → "Load the phase packet for the phase you owe work in;
   its omission index lists everything else that exists — load omitted sections on demand."
   Same change at `SKILL.md:116` (workflow step 1).
2. §9 session-start checklist item 1 (`COOPERATION.md:869`) → load packets per owed phase;
   transport and roster come from the packet header; protocol-changelog check stays.
3. CLI prompt templates — the runtime composes prompts in code (e.g. `runner.go` `gatherPriorRounds`
   and the round instruction at `:989`; PRIMARY, `ideas/protocol-read-cost-regression/FINAL.md:40`).
   These change in the same patch and are code-reviewed — this is the enforceable path, and the
   one the established constraint says rank 1 must live in (instruction layer, since the runner
   never reads COOPERATION.md: PRIMARY, `grep -c 'COOPERATION' internal/runner/runner.go
   internal/runner/phase58.go internal/app/driver_consensus.go` → `0`, `0`, `0`).
4. Hand-written facilitator prompts (like the one driving this very round, which says "Read
   parley-deck/COOPERATION.md") — **cannot be governed by text.** Any prompt can always say
   "read everything." What bounds them is measurement, not norm: the regression idea established
   the metric (full protocol read = 3.3× median per-call wall clock, n=3/arm; PRIMARY,
   `ideas/protocol-read-cost-regression/FINAL.md:19-21`). Add per-call wall-clock/byte tracking
   to the §13 retro surface; a reversion to full-read prompting shows up in the numbers within
   one idea. Normative text makes packet-loading the compliance state; telemetry makes reversion
   visible. Claiming text can stop a hand-written prompt would be a false comfort — I don't.

### Q4 — Fix-up budget

Current state, and a divergence worth recording (all PRIMARY):

- §4.0 table: `fast` cap 1, `standard` cap 2, `deliberation` **unbounded** (`COOPERATION.md:229`).
- Driver code: default `MaxFixupCycles = 3` (`internal/driver/driver.go:103-105`: `if
  cfg.MaxFixupCycles <= 0 { cfg.MaxFixupCycles = 3 }`); fast → 1, standard → 2
  (`internal/track/track.go:149,159`); deliberation → `ApplyOverrides: false`, i.e. the driver
  default (3) stands (`track.go:150-158`, `case Deliberation`). So "unbounded" is literally true only for
  human-driven deliberation runs; driver-managed deliberation is already capped at 3 today.
  The protocol text and the tool already disagree — part of this change is making them agree.

Proposal:

- **`deliberation` cap: 5 fix-up cycles.** Anchored on the measured mean of 5.1 review rounds
  (PRIMARY, `ideas/protocol-read-cost-regression/FINAL.md:20`: "review rounds 1.6 -> 5.1 (max
  24)"), noting the mean is itself inflated by the 19–24 tail, so 5 covers the bulk of ideas.
  `fast`/`standard` stay at 1/2 — they are already finite and were not the measured problem.
- **At the cap: ESCALATE, never close.** Reuse the existing machinery verbatim — Phase 8
  already states "`MaxFixupCycles` … are escalation thresholds, not close criteria. Hitting the
  budget never marks an implementation complete; it requires human review of the trajectory and
  either a new fix-up plan, a recorded operator ruling, or a decision to abandon/defer"
  (`COOPERATION.md:661-664`), and the LE-5 budget hit already escalates as "a durable blocking
  inbox note" (`:666-673`). The cap adds no new mechanism, only a number where "unbounded" stood.
- **Extension path:** the user (or full quorum) grants +N cycles, recorded in IMPLEMENTATION.md's
  Decision Log, repeatable. Each grant requires the trajectory summary stopping judgment already
  defines (findings per cycle, severity, fresh-vs-relitigated, `:646-659`). This keeps the
  19–24-round tail *possible* but makes each extension a recorded decision with evidence
  attached, instead of an accident of nobody stopping.
- **No severity floor, and say why in the normative text:** the two worst ideas kept producing
  **fresh MAJORs at rounds 19–24**, so "late findings are trivial" is empirically false here
  (PRIMARY, `ideas/protocol-read-cost-regression/consensus.md:73-74`: "the two worst ideas kept
  finding fresh MAJORs at rounds 19–24**, so a severity floor would not have stopped them" —
  quoted verbatim; same substance at `FINAL.md:30-31`). The cap is a checkpoint with a human,
  not a termination rule.

### Q5 — Never cut, and detection

Never-cut list (the always-loaded core of Q1, restated as the ratifiable invariant): Phase 2's
"Silence = implicit agreement" semantics; §6 rule 3 (never edit another's file); English-only /
no-secrets; the escalation-to-user path; the non-solo floor (§1 one-liner); §15.2 provenance;
§4.0 track table + invariants; active transport + roster; the phase's own artifact template;
and the omission index itself.

Detection if a packet cuts one anyway:

1. **Generation-time (the only layer that catches index omission):** the index is derived, not
   curated; deterministic regeneration + diff fails on any divergence. A hand-edited packet or
   generator regression is caught here, before any agent runs.
2. **Runtime (catches staleness, makes text-omission visible):** `sourceSha256` in the header;
   any participant verifies one hash; the omission index makes "a rule exists about X"
   discoverable even when the text isn't loaded.
3. **In-loop: nothing, and I won't pretend otherwise.** If a binding rule is missing from both
   packet and index, no participant is prompted to notice, and Phase 2 rule 1 converts that
   silence into recorded agreement at the signoff gate. The signoff gate ratifies the record; it
   does not audit the packet. The honest mitigation is structural (rule 1: the index can never
   be hand-curated) plus ex-post sampling (§13 retro periodically runs one idea under full-read
   and diffs outcomes), not an in-loop check.

## Concerns / open questions

- **Packet = f(phase, track, flags) adds a frontmatter dependency.** Conditional inclusion
  (§12, §14 keyed off `auto_implement` / pipeline flags) means a mis-flagged 00-prompt gets a
  thinner packet. Mitigation: the classifier already fail-closes to the stricter track
  (`COOPERATION.md:213-218`); the omission index digest still names the conditional sections.
  Open: should the generator fail-closed too (include conditionals when flags are ambiguous)?
  I lean yes — cost of over-inclusion is bytes, cost of under-inclusion is invisible.
- **Index digest quality.** A one-line MUST digest can mislead by compression. Mitigation:
  digest = verbatim extracted sentences, not summaries; the generator MUST NOT paraphrase.
  Open: who owns digest fidelity review when the extractor changes?
- **Round-2+ context is separate from protocol context.** The re-read term (prior round files
  embedded by the runtime, `FINAL.md:40`) is out of scope for this packet change but is the
  other half of the cost; the packet change must not claim to fix it.
- **Cap = 5 is anchored on a mean of 5.1 that includes the tail.** If the median is well below
  5, half of deliberation ideas would escalate once — acceptable if the escalation is cheap
  (one inbox note + recorded grant), annoying if it becomes ceremony. Open: someone should
  compute the median from the regression dataset before Phase 3 fixes the number.
- **Skill vs deck vs core.** §7 says the deck's COOPERATION.md becomes a generated view of a
  global core (`:754-757`); the packet generator should render from whatever the deck's live
  file is, so it survives that migration unchanged. Worth one normative line.

## Risks

- **False-confidence risk (the headline one):** a packet system makes omission *less* visible
  than full-read does, because full-read at least puts every rule in front of the agent. The
  omission index is the only counterweight; if reviewers sign off on this idea without the
  "index is never hand-curated" invariant in FINAL.md, the change is net-negative.
- **Reversion risk:** hand-written facilitator prompts quietly return to "read everything"
  (cost regression, silent) or, worse, a stale committed packet gets copy-pasted into prompts
  (correctness regression, silent). Mitigated only by the render-don't-store rule and the
  wall-clock telemetry; both must ship together.
- **Conditional-inclusion gap:** an `auto_implement` idea whose flag is dropped mid-idea (or
  never set) runs P5 without §14 loaded. The index one-liner is the thin remaining defense.
- **Cap anchoring risk:** 5 may be wrong in either direction; the escalation mechanism makes a
  wrong-low cap annoying rather than fatal, and a wrong-high cap just preserves today's tail.
  Directionally safe, numerically uncertain.
- **Divergence debt:** the protocol already says "unbounded" while the driver defaults to 3
  (see Q4). If this idea fixes only the text and not the seeded defaults (`parley init`,
  `~/.parley [defaults.loop]`, `:666-673`), the disagreement persists under a new number.
