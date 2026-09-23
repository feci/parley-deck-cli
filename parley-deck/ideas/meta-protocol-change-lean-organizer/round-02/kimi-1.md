---
agent: kimi-1
idea: meta-protocol-change-lean-organizer
round: 2
date: 2026-09-23
responding-to: [claude-1/round-01, zcode-1/round-01]
---

## Protocol context attestation

```json
{"context_mode": "full", "source_sha256": "12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18", "packet_sha256": "12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18", "fallback_reason": null, "body_path": "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer/.parley-runtime/protocol-packets/full-phase2-deliberation-12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18.md"}
```

Read in full (1,380 lines) before writing; `source_sha256 == packet_sha256` as expected for
`context_mode: full`. I then read `00-prompt.md`, `organizer-notes.md`, and both peer round-01
files. Round-2 discipline: I write only this file; I do not edit peers' artifacts.

**Organizer operational notes are testimony, not verdicts (§15.1).** Everything in
`organizer-notes.md` — the kimi stream-json readiness envelope, the `continue --auto` fresh-run
startup gap, and the historical-worktree registration blocker that stopped the driver's
cross-review accounting before any round-02 launch — is organizer-recorded operational
observation I did not independently re-verify. I treat it as SECONDARY testimony consistent with
the owner's standing authorization ("When the driver cannot do a step, fall back manually and
record why in the idea"), and as design-relevant evidence for B (state surfaces must expose
*why* a wait is blocked), never as a code verdict. No guard, ledger, or worktree history was
altered for this fallback; that constraint carries into my proposal unchanged. This is manual
cross-review 1 of the deliberation track's maximum 3.

## Position changes since prior round

1. **Flag name: `--role` → `--audience`.** claude-1 is right that `Role` inside
   `internal/protocolpacket` already denotes the deck's `protocolRole` (`source|consumer`) — a
   collision I myself cited at `internal/protocolpacket/packet.go:113-122` in round 1 (PRIMARY).
   Reusing the word would invite exactly the confusion his R10 names. I adopt
   `--audience participant|facilitator`, with the attestation gaining an additive `audience`
   field and `context_mode` values unchanged (zcode-1's invariant, which I share: participant
   launches keep the ratified full default).
2. **SKILL.md core budget: ≤ 25 KB → ≤ 20 KB.** Both peers independently converged on ≤ 20 KB;
   the tighter bound is measurable, and my 25 KB was a placeholder marked "to be fixed in
   consensus". Fixed now: ≤ 20,000 B.
3. **Timeout default: self-correction (§15.1, weakening, effective immediately).** My round-01
   text said the 25-minute default is "below the 30-minute provider prompt-cache lifetime named
   in the approved brief". The *owner-requested* 30-minute figure is brief guidance; no
   provider-wide TTL evidence is available from this worktree, and claude-1's claim-16
   (UNVERIFIED/RECALL) is the correct provenance state — I adopt it. Corrected position:
   **timeout configuration still provides a default bounded wait below the requested 30-minute
   cache lifetime** — a 25-minute default seeded in `~/.parley [defaults.loop]` configuration,
   overridable per call, documented as a conservative heuristic under the owner's requested
   bound, never a baked constant asserting a provider TTL. A provider change then needs a config
   edit, not a code change. This keeps the brief's intent (one bounded blocking call whose
   return lands inside the cache window the owner budgeted for) while dropping the unverifiable
   claim.
4. **D1 handoff: from "driver-owned record under `runs/`" to "computed bootstrap + existing
   driver records".** zcode-1's `--organizer` bootstrap and claude-1's anti-staleness argument
   together improve my round-01 shape. New position: the organizer-facing per-phase handoff is
   **computed on demand** (`parley status --idea <slug> --organizer`, or the composed
   `parley organizer brief` — see responses) from live state; the driver additionally persists
   per-phase records under `runs/` (it already writes handoff packets,
   `internal/runner/handoff.go:38` — PRIMARY from my round 1) as non-canonical driver state. No
   new canonical artifact type, no hand-written file, no participant obligation.
5. **Adopted from peers without residue:** claude-1's audience safety property (an audience view
   may only omit what phase/track/transport scoping already omits; never below the ratified
   never-cut floor), his `PhaseDigest` generalization (design rounds → review rounds → signoff
   state → implementation status) reusing the shipped validators, and his loud-degradation exit
   semantics (`wait` exits distinctly on invalid artifacts and on escalation events, not just on
   timeout).

## Responses to others

### @claude-1 — round-01

Agreements I adopt outright: mechanism-over-text (your 2026-09-15 instruction-failure evidence
matches my ALT-7 reasoning); the digest-already-exists finding (TUI-only `BuildRoundDigest`) as
B's core; `--audience` naming; the C1 byte table; "no cut inside §4"; the generated-only
organizer brief with a read-only-deck test; ingest-accepts-a-path-never-a-number as the
mechanical answer to the self-report rejection (I proposed the same discriminator; your
flag-set test is the right shape).

**Disagreement 1 — the facilitator reading set needs §15, not only the literal Quickstart
table.** Your C1 scopes to Quickstart, §4, §5, §9, §11.B (54,654 B measured, PRIMARY).
Counter-proposal: the facilitator audience set is the protocol's own table **plus §15**.
Reasons: (a) §15.5 binds the facilitator directly — its procedural calls are provisional, and
its dispute-adjudication limits are exactly what an organizer exercises when digesting
signoffs; (b) scope A makes the organizer a *consumer* of verdicts, provenance tags and
`DISPUTED` states — §15 is the grammar of everything it will read in the digest; (c) the
Quickstart itself names §15 as binding on every track. Cost is modest: §15 spans lines
1246–1380 (~11 KB by line-average arithmetic on your table), putting the set near ~66 KB — under
your own 75,000 B criterion and well under the full 108,400 B. This is the one substantive
boundary I want consensus to fix explicitly; the `audiences:` map makes either answer auditable,
but my position is +§15.

**Disagreement 2 (partial) — D1 should not drop the file form entirely.** You propose no
state/handoff file at all; the generated brief *is* the handoff. The owner-approved D text
explicitly names "a state/handoff file so each phase can start in a fresh organizer session", so
a FINAL with no file form at all risks reading as a dropped sub-item under your own stop rule.
Counter-proposal (my position change 4): computed bootstrap as the organizer surface **plus**
the driver's existing, non-canonical per-phase records under `runs/`. Nothing hand-written,
nothing canonical, nothing a participant must produce. Your sidecar rejection targeted stored
summaries becoming a second source of truth; a driver state record whose authority is always
the recomputed view does not have that failure mode. Please confirm this middle path addresses
the objection, or block it with reasons — this is the second point consensus must close.

Answers to your open questions: (2) facilitator-who-is-also-a-participant → full context;
audience scoping applies only to a non-participant facilitator — I concur, and it pairs with
your own `facilitator:` field. (3) Timeout — resolved by configuration, above; no PRIMARY TTL
locator is claimed by any of us. (4) `wait` must never advance a phase — agreed; it is a read of
the event log and the artifact tree. (5) Skill-copy drift check — in scope as implementation
hygiene (a skill-repo test along the existing manifest-coverage pattern), because 00-prompt puts
all three copies in implementation scope; it is tooling, not a protocol obligation on any agent.
(6) Loud degradation — adopted; `valid`/`fell_back` columns prominent, and an invalid artifact
or `unparsed` stance forces a distinct non-zero exit.

Your A mechanism (optional `facilitator:` field, fail-closed preflight, driver predicate
excluding the declared facilitator from drafter/implementer/reviewer/goal-done roles) is
stronger than my round-01 "skill text + prompt repair" and I adopt it, with one boundary check
against the stop rule: the field is optional, absent → byte-identical behavior, so it adds no
mandatory obligation — consistent. I keep my round-01 A.2 as a *separate, bounded* repair: the
driver's consensus drafting prompt must emit the full §15 duties (`## Verdict conflicts` when
conflicts exist, `## Drafter position changes`, `## Alternatives disposition`,
`## Comparison & blind spots`) generated from one `protocol.RequiredConsensusSections` constant
so prompt and gate cannot drift. That is squarely inside the 00-prompt constraint ("the old
driver prompt may omit them") and inside A ("make this the default in the driver"); it is not
unrelated scope.

### @zcode-1 — round-01

Agreements I adopt: participant-implementer resolution (FINAL drafter → inbox claimant → first
available participant; escalate rather than silently fall back to the facilitator); the
`--organizer` bootstrap built on existing status machinery and `WriteHandoffPacket`; keeping
`status` non-blocking and wrapping its parsing in a new command rather than changing its
behavior for existing callers; ambiguity-labeled usage attribution with the method stated in the
ledger header; fail-closed `unparsed` stance; the Windows CI leg for `wait`/`usage`; portable
polling where no driver event log exists. Your protocol-text delta list (§9.0 facilitator-view
line; Phase 5 + Quickstart facilitator row; §11 advisory prefer-one-blocking-wait) converges
with claude-1's and mine; the merged minimal set is in `## Current proposal`, renamed to
`--audience`.

Answers to your concerns: (1) kimi structured-usage shape UNVERIFIED — agreed, and the honest
`coverage: none` path stands; note the organizer's readiness log (testimony, per my attestation
note) reports a real stream-json envelope, so a fixture-based parser is plausible but must be
built from captured output at implementation time, not assumed. (6) Should `wait` return early
on a new blocking escalation? Yes — adopt claude-1's semantics: `wait` breaks on
`driver.error`/escalation events with a distinct exit code, so an organizer blocked in `wait`
wakes when a `to-user` escalation lands; unanswered escalations also surface in the bootstrap.
(7) Stance-parser reuse unverified — agreed; FINAL will carry "extract shared signoff/stance
parsing, do not fork it" as an implementation-time verification item (a code-review checklist
line, not a new protocol obligation).

**Disagreement — the ≤ 60 KB facilitator-packet cap (AC-C1).** If the facilitator set includes
§15 (my position above), claude-1's own byte table puts the set near ~66 KB, over your cap.
Counter-proposal: the byte budget is **derived from the ratified map, not fixed in advance** —
the acceptance criterion is "emitted body equals exactly the map's verbatim blocks, byte-counted
and recorded with the attestation", with an absolute guardrail of ≤ 75,000 B (claude-1's
criterion) for the phase-1 / deliberation / github-pr case. We fix the map first, then the
number falls out. If consensus picks the literal table without §15, your 60 KB stands
automatically.

Your sequencing sketch (protocol deltas → renderer → skill split → wait → driver resolution →
bootstrap/usage/telemetry → release prep) is the right Phase 5 order and I adopt it, with one
addition: the consensus-prompt repair rides with the driver-paths step, and the skill-copy
parity test rides with the skill-split step.

## New concerns / questions

1. **Retained disagreement register** (for the consensus drafter): (i) facilitator set = literal
   Quickstart table vs table + §15 — I hold +§15; (ii) D1 file form: driver-persisted
   non-canonical records + computed bootstrap (me, and I read zcode-1 as compatible) vs no file
   form at all (claude-1) — I hold the middle path; (iii) packet byte budget: map-derived with a
   ≤ 75 KB guardrail (me) vs fixed ≤ 60 KB (zcode-1) — resolves itself once (i) is fixed. No
   other substantive disagreement remains that I can see; everything else converged this round.
2. **Driver startup gaps are evidence, not scope.** The preflight kimi stream-json parser
   rejection, the `continue --auto` fresh-run launch gap, and the historical-worktree
   registration blocker are recorded organizer observations (testimony) and good B-side evidence
   that state surfaces must say *why* they block. Their repair is optional hardening inside the
   release's driver work, to be re-verified at HEAD during implementation — FINAL must not turn
   them into mandatory obligations, and the artifact-heading/formatting repairs already done
   this run are not design scope at all.
3. **§15 state.** No conflicting verdicts exist between the three round-01 files that I can
   find — where figures differ (e.g. per-idea `status --json` 1,459 B PRIMARY at this deck vs
   ~2.2 KB SECONDARY from the study) they are compatible, not contradictory, so §15.3's
   `## Verdict conflicts` section is not triggered yet. The study's quantitative corpus remains
   SECONDARY testimony for all three of us; no acceptance criterion rests on it. §15.6's
   consensus duties (correlated-agreement statement; `## Alternatives disposition` covering each
   round-01 ALT id) stand ready for the drafter.
4. **Drafter/implementer.** With `author: user`, the default FINAL drafter is the first round-01
   writer — recorded for this run as kimi-1 — and the default implementer is the FINAL drafter.
   **I am willing and claim both defaults**, unless claude-1 or zcode-1 volunteers by the normal
   Phase 4/5 route (an `inbox/` note before consensus signoff completes, signoff block stating
   `Drafter: yes`). If I draft, §15.5's `## Drafter position changes` and §15.6's
   `## Alternatives disposition` bind me at consensus, and this round's position changes are
   written to make that section easy to check.
5. **Cross-review budget.** This is cross-review 1 of the deliberation cap of 3. With the
   register in item 1, I believe consensus is reachable in round-03 without escalation — but the
   cap is an escalation trigger, not a target; if either retained disagreement survives
   round-03, it goes to the owner rather than being smoothed over.

## Current proposal

Converged A–D design, restated as the buildable shape I would carry into FINAL:

- **A — facilitator boundary, enforced.** Optional `facilitator:` in `00-prompt.md`; preflight
  fails closed when it names an agent also in `participants:` unless
  `facilitator_participates: true`; the driver never selects the declared facilitator as
  drafter, implementer, reviewer or goal-done checker; implementer resolution is
  participant-only with escalation, never silent facilitator fallback (zcode-1's rule). Driver
  consensus/FINAL prompts emit the full §15 drafter duties from one protocol constant (my
  repair). Skill text states the boundary as the default. Protocol text: one permissive sentence
  each in Phase 5, Phase 6, §9.0; absent field → byte-identical behavior. No new mandatory
  obligation.
- **B — `parley wait` + `PhaseDigest`.** New read-only command blocking on the event log (or
  bounded ≥ 10 s filesystem polling where no driver runs) until the round / consensus / review /
  implementation boundary, timeout, or an escalation/error event. Exit codes: `0` boundary
  reached; distinct codes for timeout (partial digest still printed), escalation or invalid
  artifact (loud), and usage error. Digest: deterministic, Go-generated, byte-identical on an
  unchanged tree; per-agent path / filed / bytes / owner / valid (verbatim verdicts of the
  existing validators) / stance flags / fell_back; extends the shipped TUI-only digest to all
  phases; carries every raw path; never persisted as canonical; no model-written field. Default
  bounded wait 25 minutes, seeded in configuration below the owner's requested 30-minute
  cache-lifetime guidance — a heuristic default, not provider-TTL evidence.
- **C — audience view, slim skill, generated brief.** `--audience facilitator` on the existing
  packet renderer via an `audiences:` key in `meta/packet-applicability.yaml` (a §7-governed
  protocol file; this idea is that §7 change), constrained to omit only what phase/track/
  transport scoping already omits and never the never-cut floor; verbatim blocks + complete
  omission index; additive `audience` attestation field, `context_mode` unchanged; facilitator
  set = Quickstart, §4, §5, §9, active §11 transport **plus §15** (my retained position);
  non-participant facilitators only, else full. SKILL.md split: ≤ 20,000 B driver-first core
  naming `run` / `continue` / `wait` / `status` / `consensus` / `preflight`, relocation-only into
  named `references/` files, proven by the skill repo's coverage test pattern. Organizer brief:
  one generated, read-only, ≤ 8 KB composition (attestation pointer + idea status + current
  digest + next action), byte-identical on an unchanged tree, never stored.
- **D — fresh session per phase + client-accounting ledger.** The organizer re-orients from the
  generated brief + `--organizer` bootstrap (≤ 4 KB: phase, owed artifacts, open escalations,
  resume command); the driver keeps its non-canonical per-phase records under `runs/` (extended
  handoff packets), satisfying D's state/handoff file without a canonical artifact or
  participant duty. `parley usage` ingests only client-written accounting files (Codex rollout
  `token_count`, Claude `message.usage`) by path — never a number argument — appending one
  provenance-carrying row (client, phase, source path, event type, ingest time; tokens only, no
  dollars; ambiguous attribution labeled) per phase boundary; idempotent re-ingest. Add the kimi
  telemetry parser from captured fixtures and structured-output default argv where the adapter
  supports it; honest `coverage: none` if kimi emits no usage.
- **Protocol copies & release shape.** All three COOPERATION.md copies carry the same permissive
  hunks; the CLI drift test keeps guarding the two CLI copies; a skill-repo parity check is
  added as implementation tooling (not a protocol obligation). Staged global core built from the
  previous core template + reviewed hunks; `parley protocol publish` remains the owner's
  attended action. CLI 1.49.0 / skill 2.13.0 after Phase 8, direct-to-main per the owner's
  run-scoped override.

**Stop-rule compliance (stated for the organizer's FINAL check):** nothing above drops A, B, C
or D — including D's handoff file (driver records + computed bootstrap) and D's ledger; nothing
adds a mandatory obligation beyond A–D — every new surface is optional (a field), permissive (a
protocol line), generated (brief / digest / bootstrap), or repo tooling (tests). The only
retained disagreements are the three in `## New concerns / questions` item 1, all inside A–D's
envelope.

**Provenance for this round (§15.2):** verdicts about peers' claims are SECONDARY resting on
their named round-01 PRIMARY verdicts (acyclic, terminating in their HEAD locators); where I had
my own PRIMARY from round 1 (`packet.go:113-122`, `handoff.go:38`,
`driver_consensus.go:112-167`, `usage.go`, the `app.go` dispatch table), I reuse it. The ~66 KB
§15-inclusive figure is arithmetic on claude-1's PRIMARY byte table — an estimate, marked as
such, to be replaced by the renderer's own byte count at implementation. Organizer-notes
content is testimony per the attestation note. I issue no verdict on any design proposition I
own (§15.1).
