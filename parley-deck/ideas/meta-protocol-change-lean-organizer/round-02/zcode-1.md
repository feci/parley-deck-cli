---
agent: zcode-1
idea: meta-protocol-change-lean-organizer
round: 2
date: 2026-09-23
responding-to: [claude-1/round-01, kimi-1/round-01]
---

## Protocol context attestation

```json
{"context_mode": "full", "source_sha256": "12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18", "packet_sha256": "12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18", "fallback_reason": null, "body_path": "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer/.parley-runtime/protocol-packets/full-phase2-deliberation-12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18.md"}
```

I recomputed the SHA-256 of the emitted phase-2 packet body; it matches both fields above and is
byte-identical to `parley-deck/COOPERATION.md` at HEAD `4004a99` (1,380 lines, 108,400 B). The diff
`4ce8fa7..4004a99` touches only `parley-deck/ideas/**` files, so every code locator my round-1
verified at `4ce8fa7` is still valid; all new locators below were verified at `4004a99` in this
session. Organizer operational notes in `organizer-notes.md` are treated as observations, not code
verdicts; every code claim I rely on was re-checked by me at HEAD.

## Position changes since prior round

1. **C flag renamed: `--role facilitator` withdrawn, `--audience participant|facilitator` adopted**
   (claude-1's finding). Verified: `SourceInfo.Role` at `internal/protocolpacket/packet.go:113-122`
   is the deck's `protocolRole` (`source|consumer`, switched in `internal/protocolpacket/source.go`
   `ResolveSource`), so a `--role` flag would collide with a live, different concept. Semantics
   unchanged from my round-1: whole verbatim sections plus a complete omission index; the
   attestation gains an `audience` field; `context_mode` values stay untouched for everyone else.
2. **The 43–50 KB facilitator-body target is withdrawn; replaced by a measured ≤70,000 B ceiling.**
   My round-1 echoed the study's 43–50 KB figure (SECONDARY). Reproduced claude-1's PRIMARY byte
   table exactly — Quickstart 2,414 + §4 35,991 + §5 1,334 + §9 5,624 + §11.B 9,291 = 54,654 B —
   and measured §15 (lines 1246–1380) at 8,056 B; the facilitator set plus §15 is 62,710 B (57.8%).
   The ratified `neverCut` floor (`internal/protocolpacket/applicability.go:193-231`) additionally
   pins §0, §4.0, §6, §14 and the transport-matched §11 subsection in every packet. I now adopt
   claude-1's position that cutting inside §4 beyond per-phase scoping to chase the brief's
   43–50 KB is not worth the risk; the honest converged target is ~63–66 KB with a hard ceiling of
   70,000 B.
3. **Facilitator reading set = protocol table + §15 — converged, and no longer a choice.** I
   adopted kimi-1's "plus §15" lean and then found it is forced: `neverCut` pins §15.1–15.4 and
   §15.7 at kernel phases `[1,2,3,5,6,7,8]` (`applicability.go:190,207-212`), so an audience view
   cannot omit §15 at deliberation phases even if we wanted to. Kimi-1's concern 1 closes with no
   open question remaining.
4. **Digest foundation moves to the shipped digest.** My round-1 proposed building `parley wait`'s
   digest on the status/consensus-status parsers. Claude-1's finding is better and verified:
   `BuildRoundDigest` already exists (`internal/driver/digest.go:10-48`), is built at
   `internal/driver/driver.go:488`, and its only consumer is the TUI Home tab
   (`internal/tui/roundsummary.go:9-38`) — invisible to a headless organizer. New plan: extend it
   to a `PhaseDigest` (rounds, reviews, consensus, implementation), surface it on stdout via
   `parley wait`, and take validity from the three shipped validators. My fail-closed `unparsed`
   stance rule and the open-the-raw-artifact-to-adjudicate discipline carry over unchanged.
5. **Exit codes enriched from 0/2/1 to 0/3/4/1** (claude-1's scheme), with one refinement —
   see the response to claude-1 below: missing artifacts keep waiting; only a validator-failing
   artifact or a blocking escalation/driver-error event exits early with `4`; timeout exits `3`
   and still prints the partial digest.
6. **Timeout resolved as configuration with a concrete bounded default.** Adopting claude-1's
   "configuration, not a constant" principle together with a 25-minute default. To the standing
   question: **yes — the configuration still provides a default bounded wait below the requested
   30-minute cache lifetime (25 m), and we assert no provider-wide TTL evidence anywhere.** The
   30-minute figure is the owner's brief requirement (testimony); claude-1's claim 16 tags it
   UNVERIFIED and I could not verify it either. The acceptance criterion asserts only
   `default < 30m && configurable`, and code comments/docs may cite the brief, never provider
   documentation.
7. **D's re-orientation surface converges on one computed command.** My round-1
   `parley status --organizer` bootstrap and claude-1's `parley organizer brief --idea <slug>` are
   the same surface; I adopt claude-1's command with the ≤8,192 B bound, byte-identical-on-
   unchanged-tree determinism, and the read-only-deck (writes-no-file) test. Kimi-1's driver-owned
   handoff record survives as non-canonical run state under `runs/` — it satisfies the brief's
   "state/handoff file" letter without creating a canonical sidecar (details in the response to
   kimi-1).
8. **Usage command renamed `parley usage record` → `parley usage ingest`** (claude-1/kimi-1
   naming), with claude-1's mechanical rule adopted verbatim: it accepts a path, never a number;
   streaming ingest for the 228 MB Codex rollout; ≤1 KB stdout; idempotent re-ingest; kimi-1's
   round-trip (independent re-parse) criterion; my `ambiguous` attribution labeling retained.
9. **SKILL.md core budget fixed at ≤20,000 B.** Kimi-1 explicitly deferred the number to
   consensus ("≤ 25 KB, to be fixed in consensus"); claude-1 and I independently landed on 20 KB.
   With the relocation-script criterion (nothing dropped, every moved heading lands in exactly one
   reference) as the safety net, the tighter number costs nothing.
10. **Locator self-correction (§15.3 discipline).** My round-1 cited
    `parley-deck/COOPERATION.md:439` for the drafter's pre-publish verification duty; the duty is
    at `:437` ("Before publishing `FINAL.md`, the drafter MUST verify that every active
    non-facilitator participant has created the expected canonical artifacts…"). The
    default-implementer sentence is `:443-444`, and §15.5's provisional-calls text spans
    `:1340-1343` — my point locator `:1342` was inside it but imprecise. Peers' locators were
    correct; mine is corrected.

## Responses to others

### @claude-1

Your load-bearing claims are real; I verified each at HEAD `4004a99` before adopting them: the
digest exists and is TUI-only (`internal/driver/digest.go:10-48`, built at
`internal/driver/driver.go:488`, sole consumer `internal/tui/roundsummary.go:9-38`); the `Role`
name collision (`internal/protocolpacket/packet.go:113-122` + the `source|consumer` switch in
`source.go`); no `wait`/`usage`/`organizer` verbs in the dispatch table; the never-cut floor is
ratified and per-phase (`internal/protocolpacket/applicability.go:193-231`); and your byte table
reproduces to the byte in my hands (54,654 B literal facilitator set). Your framing — A–D is
mostly a routing and surfacing problem over machinery that already exists — is the correct one and
my round-2 design is rebuilt on it.

**Adopted:** the optional `facilitator:` declaration with preflight fail-closed on
facilitator-also-in-participants and the driver ineligibility predicate (A); `PhaseDigest` +
blocking `parley wait` over the append-only event store (`internal/store/events.go:42,68`) with
exit codes 0/3/4/1 (B); `--audience` with the safety property that the audience dimension may only
omit blocks already omittable at that phase/track/transport and never below `neverCut` (C1); the
SKILL.md split with relocation-only discipline and the six-command naming test (C2); the computed,
never-stored organizer brief (C3); and the path-only streaming usage ingest (D2). Your "what I
deliberately do not propose" list matches my round-1 guardrails; I co-sign all of it.

**Counter-proposals and answers:**

- **Byte ceiling 70,000 B, not 75,000 B.** Your ≤75,000 B target was derived from phase-1 scoping
  alone (73,181 B remains). With the audience dimension additionally omitting §1–§3, §8 and §10
  (~9.6 KB, none never-cut), ~63.5 KB is reachable without touching §4's interior. A 75 KB ceiling
  would pass a regression that quietly adds §2 back; 70,000 B is tight enough to catch that and
  loose enough for zone variance. My counter stands unless you have a case for the slack.
- **Your concern 6 (wait exits 4 on invalid artifacts): yes, refined.** Missing ≠ invalid. An
  artifact that has not landed yet keeps the wait blocking; exit `4` fires only when a present
  artifact fails the shared validator (the awaited condition cannot be met until it is repaired)
  or when a `driver.error`/blocking escalation event arrives. Without this split, one slow
  participant would be misread as blocked — the exact misread a facilitator should never make from
  a digest.
- **Your concern 2 (facilitator who is also a participant): your lean is right; recorded as
  decided.** Audience scoping applies only when the declared facilitator is not a participant;
  `facilitator_participates: true` → full context, no narrowing. A signing agent must never
  silently lose reading set.
- **Your claim 16 (30-minute TTL UNVERIFIED): agreed, and resolved** as position change 6 — the
  default is 25 m, bounded below the requested 30 m, configuration-first, and no artifact of this
  release asserts provider TTL evidence. Your instinct not to bake the figure into Go is exactly
  why the property survives.
- **Your concern 1 (22 KB skills catalogue outside A–D): agreed, unchanged scope.** It stays out;
  the D ledger is what will make the residual visible rather than arguable. Neither of us proposes
  a scope change, and I will not volunteer one.
- **Your concern 5 (unguarded skill copy): all three of us now converge** on a repo-side
  content-parity test for `skills/parley-deck/references/COOPERATION.md`, recorded in FINAL as
  tooling hygiene — not an obligation on any agent. My round-1 said "kept optional"; I have
  dropped that hedge (see also the response to kimi-1).

### @kimi-1

Your central A-repair claim is confirmed at HEAD and it is the most consequential finding of
round 1: `buildConsensusDraftPrompt` (`internal/app/driver_consensus.go:112-131`) hardcodes five
headings — `## Agreed decisions`, `## Trade-offs accepted`, `## Deferred follow-ups`,
`## Dismissed findings`, `## Signoffs` — carrying none of §15.3 `## Verdict conflicts`, §15.5
`## Drafter position changes`, §15.6 `## Alternatives disposition`, nor the advisory
`## Comparison & blind spots`, while `buildFinalDraftPrompt` (`:142-167`) already derives its list
from `protocol.RequiredFinalSections`. The protocol's own Phase-3 header
(`parley-deck/COOPERATION.md:361-365`) explicitly carries the §15.5/§15.6 drafter duties. This
round is itself live evidence: the cross-review instructions had to re-impose the §15 duties by
hand because the driver's prompts omit them. **Adopted in full:** one
`protocol.RequiredConsensusSections` constant, prompt and gate derived from it (parity test), and
your mitigation for the heading-rename risk (search for heading consumers during implementation)
goes into the acceptance set.

**Also adopted:** §15 in the facilitator set (now forced by `neverCut` — see position change 3);
stance parsing fail-visible as `unparsed`, never omitted; extending
`meta/packet-applicability.yaml` rather than a sibling file (one §7 surface the drift story
already knows); the kimi telemetry parser built from fixtures captured this run with an honest
`coverage: none` fallback; per-phase ledger granularity; and the skill-copy content-parity check
as repo tooling — your framing, which claude-1 and I have now joined, closing claude-1's concern 5.

**Counter-proposals and disagreements:**

- **The handoff record is run state, not the re-orientation surface.** I side with claude-1's
  invariant (no stored summary sidecar the organizer must trust) but keep your letter of the
  brief: the driver MAY write its per-phase handoff record under `runs/` (driver state per §3,
  non-canonical, advisory — `parley status`/brief recomputation stays the authority), and the
  brief's D wording "a state/handoff file" is satisfied by that driver-written record. What I
  drop from your proposal is its mandatory acceptance criterion (iii) ("handoff record exists
  with schema-valid fields"): the FINAL-mandatory criterion becomes the *computed brief's* output
  contract; the `runs/` record is an implementation detail a driver unit test may cover. FINAL
  should state this reading of D explicitly so the organizer can check it against the stop rule
  with open eyes — D's substance (fresh-session bootstrap + ledger) is preserved, not dropped.
- **SKILL.md budget 20,000 B, not 25,000 B.** You deferred the number to consensus; claude-1 and
  I each measured 20 KB workable, and the relocation-script criterion guarantees nothing is lost
  by the tighter cap. If implementation genuinely cannot fit 20 KB, the consensus round is where
  to show the bytes, not the FINAL.
- **Your AC (v) reworded.** "Default timeout constant < 30 min" becomes
  `default timeout < 30m && configurable`, with the no-provider-TTL-evidence rule attached
  (position change 6). The bound is ours to guarantee by choice of default; it must never be
  presented as a verified provider property.

## New concerns / questions

1. **Operational evidence from this run, recorded under B — not new mandatory scope.** Two
   observations from the organizer's fallback notes, both of which I treat as observations, not
   code verdicts: (a) `parley status` reported the round as complete-looking while two artifacts
   failed heading validation, and the *reason* was invisible until a disposable helper ran the
   existing validator — so the PhaseDigest validity column must carry the failing check name
   (e.g. `missing required section "X"`), not a bare boolean; this folds into B's digest design
   and adds no obligation. (b) The driver refused the round-02 launch because cross-review
   accounting referenced an unavailable historical worktree path, and labeled the escalation
   author `claude` although codex-1 organizes this run. That is evidence for the owner-sanctioned
   manual-fallback path and for attribution correctness in D's ledger; I propose no guard, ledger
   or code change around it in this release, and I verified no guards or ledgers were altered.
2. **Question for consensus (small):** should `parley preflight` soft-warn — not fail — when no
   `facilitator:` field is declared, nudging decks to declare one? I lean yes; a warning creates
   no obligation. No answer needed before FINAL if the quorum is indifferent; default is no
   warning.
3. **Retained disagreement: none blocking.** The three round-1 designs compose into one
   buildable A–D. The only open deltas are numeric bindings I have counter-proposed concretely
   (70,000 B ceiling; 20,000 B skill core; exit-code map with the missing/invalid split), plus
   the two small questions above.
4. **Mandatory-obligation audit against the owner's stop rule.** My converged proposal adds no
   mandatory obligation beyond A–D: the `facilitator:` field is optional and absent-decks are
   byte-identical; the audience view is opt-in and facilitator-only; `parley wait` is an
   alternative to polling, never a required step; the skill-copy parity check is repo tooling;
   the consensus-prompt repair is inside A's "default in the skill and the driver" mandate and is
   explicitly anticipated by this idea's constraints ("the old driver prompt may omit them"). If
   FINAL drifts into writing any of these as participant duties, that is new scope and the
   organizer should stop.

## Current proposal

The converged design, stated as deltas on my round-1 proposal (unchanged detail remains there);
peers' mechanisms credited above.

- **A — organizer boundary, enforced.** Optional `facilitator: <agent-id>` in `00-prompt.md`;
  `parley preflight` fails closed (exit non-zero, both fields named) when the declared facilitator
  also appears in `participants:` unless `facilitator_participates: true`; the driver refuses the
  declared facilitator as drafter, implementer, reviewer or goal-done checker; a deck without the
  field produces a byte-identical run plan (regression test). If no participant implementer can be
  launched, the driver escalates — it never silently falls back to facilitator implementation.
  Skill text states the boundary as the default practice. Consensus/FINAL drafting prompts derive
  their section lists from `protocol.RequiredConsensusSections` / `RequiredFinalSections`
  including the §15.3/§15.5/§15.6 drafter duties, with a prompt/gate parity test. This run is the
  witness: `IMPLEMENTATION.md` must record `implementer:` ≠ facilitator.
- **B — `parley wait` + PhaseDigest.** Extend the shipped `BuildRoundDigest` to a `PhaseDigest`
  over design rounds, review rounds, consensus signoff state and implementation status; per-agent
  columns `agent/path/bytes/owner/valid`(with failing-check reason)`/stance_flags/unparsed/
  fell_back`, all mechanically derived, sorted, byte-stable on an unchanged tree. Command:
  `parley wait --idea <slug> --for round-NN|consensus|review-round-NN|implementation|any
  [--timeout 25m] [--json]`, blocking on the append-only event store; exit `0` condition met,
  `3` timeout (partial digest printed, outstanding agents named), `4` invalid artifact or blocking
  escalation/driver error, `1` usage/IO error. Default timeout 25 m, configurable; no
  provider-TTL claims anywhere. Digest is printed, never canonical; any ❌, DISPUTED or
  `unparsed` sends the organizer to the raw artifact.
- **C — audience-scoped packet, slim skill, generated brief.** `--audience participant|facilitator`
  on `parley protocol packet`; `audiences:` key in `meta/packet-applicability.yaml` (§7-governed);
  audience may only omit what phase/track/transport already omits, never below `neverCut`;
  facilitator set = Quickstart, §4, §5, §9, active §11 subsection, §15; body ≤70,000 B with a
  complete omission index; `facilitator_participates: true` → full context; attestation gains
  `audience`; `full` remains the default for everyone else. SKILL.md split into a ≤20,000 B core
  (naming `run`, `continue`, `wait`, `status`, `consensus`, `preflight`) plus `references/`, with
  the relocation script proving nothing dropped. `parley organizer brief --idea <slug>` ≤8,192 B,
  computed from live state, never stored (read-only-deck test).
- **D — fresh sessions plus honest measurement.** Driver may write a per-phase handoff record
  under `runs/` (non-canonical, advisory; recomputation is authoritative); re-orientation =
  organizer brief + `parley status --idea --json` + facilitator packet. `parley usage ingest
  --agent <id> --source codex-rollout|claude-jsonl --path <file> --idea <slug> --phase <n>`:
  path-only (cannot accept a number), streaming (bounded memory on the 228 MB rollout), ≤1 KB
  stdout, idempotent, appends one ledger row plus a machine-readable sibling carrying source path
  and parser id; unresolved attribution labeled `ambiguous`; round-trip re-parse test. Kimi
  telemetry parser from fixtures captured this run; if structured usage is unobtainable, a visible
  `coverage: none` record — never a fabricated number.
- **Protocol text (three copies) + staging.** Minimal permissive deltas: §9.0 item 1 gains the
  facilitator audience view; §4 Phase 5 and the Quickstart Facilitator row gain the
  participants-by-default sentence; §11 advisory line preferring one blocking wait over polls.
  Skill copy gains the repo-side parity test (tooling, not obligation). Staged core built from
  the previous core template plus these hunks; `parley protocol publish` remains the owner's
  attended action. Release mechanics per the brief (CLI 1.49.0, skill 2.13.0, all channels,
  participant-verified deploy).
- **Acceptance criteria:** my round-1 AC set (AC-A1, AC-B1/B2, AC-C1/C2, AC-D1/D2/D3, AC-P1)
  re-bound to the round-2 numbers — AC-C1 ceiling 70,000 B and `--audience` naming; AC-B exit
  codes 0/3/4/1 with the missing/invalid split and validity-reason column; AC-C2 core ≤20,000 B;
  AC-D1 renamed `usage ingest` with path-only + streaming + idempotence; plus new: consensus-prompt
  parity test (A) and the brief's determinism/read-only test (C3).
- **Process.** Default FINAL drafter remains kimi-1 as first round-01 writer; I have no objection
  and do not claim it. I volunteer for Phase-5 implementer (and as alternate drafter should
  kimi-1 decline) — willing, not insisting. No signoff is given in this round.

## Provenance addendum (§15.2)

- **PRIMARY, verified this round at HEAD `4004a99`** (CLI worktree; skill worktree unchanged at
  `d1e57d5`): packet/protocol byte-identity and SHA-256 (recomputed); `git diff --stat
  4ce8fa7..4004a99` = idea files only; `internal/driver/digest.go:10-48` + `driver.go:488` +
  `internal/tui/roundsummary.go:9-38` (digest exists, TUI-only); `internal/protocolpacket/
  packet.go:113-122` + `source.go` role switch (`Role` = deck `protocolRole`); 
  `internal/app/driver_consensus.go:82-100` (participant-only drafter), `:112-131`
  (consensus prompt omits §15 duties and canonical headings), `:142-167` (FINAL prompt uses
  `protocol.RequiredFinalSections`); `internal/protocolpacket/applicability.go:190-191`
  (`kernel=[1,2,3,5,6,7,8]`, `fullFifteen=[1,2,3,6,7]`) and `:193-231` (`neverCut` incl. §15 at
  kernel phases); `internal/store/events.go:42,68` (append-only event store);
  `internal/sessionstore/sessionstore.go:13-30` (session registry); byte measurements via `awk`/
  `wc -c` (2414 / 35991 / 1334 / 5624 / 9291 / 8056; set+§15 = 62,710 B of 108,400 B);
  `parley-deck/COOPERATION.md:361-365` (Phase-3 §15.5/§15.6 drafter duties), `:437`
  (pre-publish verification duty — corrects my round-1 `:439`), `:443-444` (default implementer),
  `:1340-1343` (§15.5 provisional calls).
- **SECONDARY (unchanged from round-1):** all quantitative figures from
  `source-context/organizer-token-study.md`, relied on for ranking only; no acceptance criterion
  depends on them.
- **UNVERIFIED (standing):** provider-wide prompt-cache TTL (claude-1 claim 16 concurred; the
  design no longer depends on it — position change 6); kimi structured-usage shape
  (implementation-time discovery with the honest `coverage: none` path).
- **Verdict conflicts (§15.3):** none with peers — every peer claim I relied on was confirmed at
  HEAD; one self-correction of my own round-1 locator (`:439` → `:437`), recorded above.
- **Roles:** participant only; I hold no facilitator role in this idea. Organizer notes were read
  as operational observations and were not treated as verification verdicts on any code.
