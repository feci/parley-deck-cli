---
agent: kimi-1
idea: meta-protocol-change-lean-organizer
round: 1
date: 2026-09-23
---

## Context attestation

```json
{"context_mode": "full", "source_sha256": "12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18", "packet_sha256": "12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18", "fallback_reason": null, "body_path": "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer/.parley-runtime/protocol-packets/full-phase1-deliberation-12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18.md"}
```

I read the rendered body in full (1,380 lines, all sections §0–§15 including appendices) before
writing this file. Round-1 independence kept: at the time of writing, `round-01/` contained no
other participant file (verified by directory listing), and I have not read any other
participant's artifact. Sources I did read: the packet body, `00-prompt.md`,
`source-context/organizer-token-study.md`, and the HEAD verifications listed under
`## Evidence & provenance`.

## Summary

The study's cost model is credible and its levers are mostly already built but unused: the
organizer pays for re-read context, and the biggest slices (self-implementation/verification,
post-compaction re-orientation, polling) map directly onto the four approved items. My proposal:
(A) codify the organizer boundary as a *default* in skill + driver and repair the driver's
consensus/FINAL drafting prompts so the delegated path emits the full §15 duties; (B) add one
new command, `parley wait`, that blocks on round/phase completion and returns a deterministic,
Go-generated digest — never stored as canonical state; (C) add a `--role facilitator` dimension
to the existing packet renderer (the protocol's own facilitator reading set, verbatim blocks +
omission index) and split SKILL.md into a slim driver-first core plus on-demand references;
(D) write a driver-owned per-phase handoff record and an organizer usage ledger ingested from
the *client's own* accounting files, and fix the kimi telemetry parser gap. No scope item
requires new mandatory participant obligations; the only protocol-text change is the §7-sanctioned
facilitator role view plus minimal clarifying lines.

## Scope reading and design position

### A. Organizer does not implement, and does not verify code itself

**Position: adopt, mostly as default + tooling repair, minimal protocol text.** The protocol
already permits this: with `author: user` the FINAL drafter is the first round-1 writer or a
volunteer (Phase 4), the implementer is the FINAL drafter unless claimed (Phase 5), and the
facilitator's verification duty is existence/ownership/validity of artifacts, not code
correctness (Phase 4 pre-publish check; §9 item 6). What is missing at HEAD is that the
*default path* does not enforce or even narrate this boundary:

1. **Skill text** (skill worktree `skills/parley-deck/SKILL.md`): state the boundary explicitly —
   the facilitator launches participants through the driver and reads their artifacts and
   validator output; it does not edit code and does not re-run code verification itself. This is
   a default and a discipline, not a new gate: no new mandatory obligation falls on participants.
2. **Driver drafting prompts are under-specified for §15.** At HEAD,
   `internal/app/driver_consensus.go:112-131` (`buildConsensusDraftPrompt`) instructs the drafter
   to produce `## Agreed decisions / ## Trade-offs accepted / ## Deferred follow-ups /
   ## Dismissed findings / ## Signoffs`. That omits the canonical Phase-3 sections
   (`## Agreed trade-offs`, `## Open items deferred to implementation`, `## Comparison & blind
   spots`) and every §15 drafter duty (§15.3 `## Verdict conflicts` when conflicts exist, §15.5
   `## Drafter position changes`, §15.6 `## Alternatives disposition`) — and two heading names
   drift from the protocol template. The 00-prompt constraint ("Preserve §15
   dispute/provenance/alternatives duties… the old driver prompt may omit them") makes this
   repair in-scope: if the organizer routes drafting through the driver, the driver must not
   produce artifacts that fail the protocol. Fix by generating the consensus section list from a
   single `protocol.RequiredConsensusSections` constant (the FINAL prompt already does exactly
   this with `protocol.RequiredFinalSections`, `driver_consensus.go:142-167`), so prompt and gate
   cannot drift apart again.
3. **Observable criteria (draft, for FINAL):**
   - Go test: the consensus drafting prompt contains the canonical Phase-3 section names and the
     strings "Comparison & blind spots", "Drafter position changes", "Alternatives disposition",
     "Verdict conflicts"; a second test asserts the section list is derived from the protocol
     constant (parity test with the gate).
   - Skill test: SKILL.md core contains the organizer boundary statement and routes
     implementation/verification to participants via driver commands.
   - Manual witness (this run): Phases 5–8 of this very idea are executed by a participant, with
     codex-1 acting only as facilitator — recorded in IMPLEMENTATION.md frontmatter
     (`implementer:` ≠ facilitator).

### B. `parley wait` + digest

**Position: adopt as one new read-only command; digest is derived, printed, never canonical.**
At HEAD the command surface (`internal/app/app.go:52-110`) has no `wait` or `digest`; the driver
waits internally during `parley run` but a manually-driving organizer has no blocking primitive,
which is exactly the polling the study measures (9.1–12.2% of input; 1,488 of 1,539 polls at a
1 s yield — SECONDARY, study §2).

Design:

- **Command:** `parley wait --idea <slug> [--round NN | --phase <phase>] [--timeout 25m]
  [--interval 15s] [--json]`. Blocks until the awaited condition holds or the timeout fires.
  Default timeout 25 min — below the 30-minute provider prompt-cache lifetime named in the
  approved brief; `--interval` default ≥ 10 s so the wait itself is not a poll storm.
- **Condition:** round complete = every quorum participant's expected artifact exists and passes
  the *existing* driver artifact validation (frontmatter, required sections); phase complete per
  the §4.0 per-track table. The predicate must reuse the driver's validation functions — one
  source of truth, no parallel reimplementation that can drift.
- **Digest (deterministic, Go-generated, no LLM text):** idea, phase/round, completeness,
  per-participant row (expected path, present/missing, bytes, mtime UTC, validity verdict from
  the shared validator, owner from frontmatter `agent:`), stance/blocks parsed mechanically —
  ✅/🟡/❌ from consensus signoff blocks, finding counts by severity tag from review headings —
  plus open `to: user` escalations. Stable ordering (00-prompt participant order). Same deck
  state → byte-identical digest (golden-file test).
- **Exit codes:** 0 = condition met; a distinct non-zero = timeout; another = error. This makes
  `parley wait … && parley status …` scriptable by the organizer without parsing.
- **Anti-summary guardrails (the prior rejection is explicit: no LLM/lossy summaries, no stored
  summary sidecars):** the digest is recomputed on each call from raw artifacts and printed to
  stdout; it is never written into `ideas/` as canonical state; every row carries the raw path so
  adjudication opens the artifact itself. Run logs under `runs/` may capture it as telemetry —
  that is already non-canonical.
- **Observable criteria:** (i) on a fixture deck, `parley wait --idea X --round 1` returns within
  one interval of the last participant file landing, exit 0, digest lists all three participants;
  (ii) with a missing file it blocks until timeout, exits non-zero, digest marks the file
  `missing`; (iii) digest determinism golden test; (iv) digest contains no LLM-generated prose
  (structural test: all fields parse from typed sources); (v) default timeout constant < 30 min.

### C. Organizer brief + slim SKILL.md

**Position: adopt; implement the role view inside the existing renderer, not as a new artifact
class.** At HEAD the renderer has no agent-role input — `SourceInfo.Role`
(`internal/protocolpacket/packet.go:113-122`) is the *deck* role (source/consumer), not a
facilitator/participant dimension. The packet machinery already does phase/track scoping with
verbatim blocks plus an omission index, governed by the live `meta/packet-applicability.yaml`.
This idea is the anticipated §7 change for the organizer role (study §6 invariant: full context
stays default "until … a §7 protocol change says otherwise (this idea is such a §7 change for
the organizer role)").

Design:

1. **`parley protocol packet --role facilitator`** (protocol vocabulary; "organizer" is this
   run's synonym): an additional deterministic filter over the already-parsed blocks of the live
   resolved authority. Emission rules identical to the existing packet: every included block
   verbatim, a complete omission index whose triggers say when to read the full source,
   `source_sha256` still computed over the full authority. The role→sections map lives in
   protocol-governed config (extend `packet-applicability.yaml` or a sibling role map), so
   changing it is itself a §7 change. `full` remains the default; the role view is opt-in per
   launch and recorded in the attestation (additive `role` field; `context_mode` enum unchanged).
   - Proposed facilitator set: the protocol's own reading set (Quickstart, §4, §5, §9, active
     §11 subsection — COOPERATION.md facilitator row) **plus §15**, because the organizer reads
     verdicts, provenance tags and DISPUTED states when digesting signoffs. Open for
     cross-review; the map makes the answer explicit and auditable either way.
2. **Organizer brief + re-orientation procedure:** the post-compaction floor becomes: the idea
   brief (00-prompt) + `parley status --idea <slug> --json` + the facilitator packet — replacing
   the full COOPERATION.md + full SKILL.md re-read the study measures at ~250–265 KB per
   compaction (SECONDARY, study §2.2). SKILL.md core documents this procedure.
3. **Slim SKILL.md:** split the 62,411-byte SKILL.md into a slim every-session core (role
   boundary from A, driver-first routing, protocol-packet/attestation rule, lean-output rules)
   and move the rest into `references/` loaded on demand — hand-launch templates in particular
   should defer to the runner, which already owns launch templates (the study §3 documents 70
   hand-written prompts while `parley run` was used once). The core must name `parley run`,
   `continue`, `status`, `consensus`, `preflight`, and the new `wait`; at HEAD SKILL.md mentions
   none of them (verified: no matches for those command names in
   `skills/parley-deck/SKILL.md`).
4. **Observable criteria:** (i) `--role facilitator` output contains §9.0 text verbatim and
   excludes e.g. Appendix A, with `source_sha256` equal to the full-authority hash and a complete
   omission index; renderer tests cover all four context modes unchanged plus the new role;
   (ii) SKILL.md core below an agreed byte budget (proposal: ≤ 25 KB, to be fixed in consensus)
   with every moved section present under `references/` — enforced by extending the existing
   manifest-coverage test pattern in the skill repo; (iii) SKILL.md core names the six driver
   commands above (test).

### D. Session per phase + organizer measurement

**Position: adopt; reuse existing run/session state, mechanize the ledger the organizer is
already keeping by hand.**

1. **Per-phase handoff record (driver-owned, non-canonical).** The driver already writes run
   state under `parley-deck/runs/` and a session registry exists (`internal/sessionstore`,
   sessions keyed by workspace/run/idea with terminal flags). Add a handoff record written at
   each phase transition — idea, phase, round, per-participant artifact status (the same
   computation as the B digest, shared code), open escalations, next expected action, gate
   status — so a *fresh* organizer session starts from brief + status + handoff + facilitator
   packet instead of compacting a growing conversation. Deliberately placed under `runs/`
   (driver state, §3) rather than `ideas/`, so it introduces no new canonical artifact type and
   no participant obligation. Where `parley run` is not used, `parley status --json` plus the
   digest already supplies the same content; the handoff file is the driver's persisted form.
2. **Organizer usage ledger from client accounting.** Mechanize what the brief has codex-1 doing
   manually (`organizer-usage.md` exists in this idea): a CLI ingest (e.g. `parley usage ingest
   --idea <slug> --phase <p> --client codex|claude --file <path>`) that parses the *client's own*
   accounting events — Codex rollout `token_count`, Claude `message.usage` — and appends
   cumulative totals (input, cached, output, reasoning, request count) with provenance (source
   path, event type, ingest time) to `ideas/<slug>/organizer-usage.md`. This answers the earlier
   rejection of "self-reported usage records as cost evidence": the numbers are the client
   runtime's API-side accounting read from its own logs, not a model's estimate; the tool never
   accepts free-text numbers.
3. **Fix participant telemetry capture.** At HEAD `internal/telemetry/usage.go` parses claude,
   codex, opencode and a zcode envelope — no kimi path, matching the measured 0/20 kimi capture
   (SECONDARY, study §4). Add a kimi usage parser with fixtures captured from real kimi output
   during this run, and fix the kimi adapter's default argv to request structured output (the
   "plain-text default argv" gap). Exact kimi event shape is an implementation-time discovery
   item — to be taken from observed output, not assumed.
4. **Observable criteria:** (i) kimi parser unit test: structured kimi usage event → populated
   Usage; adapter argv includes the structured flag (test); (ii) ingest appends a row whose
   fields equal the values computed by an independent re-parse of the same client file
   (round-trip test), and refuses free-text input; (iii) after a driver phase transition on a
   fixture idea, the handoff record exists with schema-valid fields; (iv) ledger rows carry
   client, phase, source path and event-type provenance.

## Existing alternatives

Mechanisms this proposal touches, and what the toolchain already ships (locators verified at
HEAD unless tagged otherwise):

- **ALT-1 — State digests.** `parley status --idea X --json` (~2.2 KB) and `parley consensus
  status --json` (~1.8 KB) already exist (study §3, SECONDARY; command surface verified at
  `internal/app/app.go:68,74`). They are snapshot printers: no blocking wait, no per-artifact
  validity/stance digest. B extends rather than duplicates them; the digest renderer should be
  shared code.
- **ALT-2 — Driver-internal waiting.** The driver already waits, validates artifacts, and
  escalates on deadlines inside `parley run` (study §3, SECONDARY). No standalone blocking wait
  exists for a manually-driving organizer — B fills that gap and the lean brief falls back to a
  bounded shell `until` loop only when the driver cannot run.
- **ALT-3 — Packet renderer.** `internal/protocolpacket` already renders phase/track-scoped
  packets with verbatim blocks and an omission index from the live authority; its only `Role`
  concept is the deck's source/consumer role (`packet.go:113-122`). C adds a role dimension to
  this renderer — no new renderer, no bundled organizer snapshot.
- **ALT-4 — Session/run state.** `internal/sessionstore` (session registry: run_id, idea,
  participants, terminal flag) and `runs/` run state already persist coordination state. D's
  handoff builds on these rather than introducing a parallel store.
- **ALT-5 — Telemetry collectors.** `internal/telemetry/usage.go` already ingests claude / codex
  / opencode / zcode usage events; D adds kimi and organizer-side ingest instead of a new
  telemetry subsystem.
- **ALT-6 — Skill references.** The skill already ships `references/` (COOPERATION.md,
  WORKED_EXAMPLES.md, compatibility.json) loaded on demand; C's split extends an existing
  pattern rather than inventing one.
- **ALT-7 — Prior rejected directions (scoped null).** Per study §6 (SECONDARY): LLM or lossy
  summarization of artifacts, stored summary sidecars, tools that select which protocol text an
  agent sees via retrieval (cognee/graphify/vector stores), asymmetric reviewer context, silent
  model swaps, self-reported usage as cost evidence. This proposal deliberately routes around
  all of them: digests are deterministic and non-canonical; the role packet is a §7-governed
  renderer input with verbatim blocks, not retrieval; usage comes from client accounting files.
  Sources consulted: the study's prior-decisions list and the invariants in the packet body §9.1
  and §15.

## Concerns / open questions

1. **Facilitator packet composition.** Should the facilitator role set be exactly the protocol
   table (Quickstart, §4, §5, §9, §11) or that plus §15? I lean to "plus §15" (the organizer
   reads verdicts and DISPUTED states), but this is the kind of boundary the quorum should fix
   explicitly since the map is protocol-governed.
2. **Digest stance parsing robustness.** Signoff blocks are semi-structured markdown; mechanical
   ✅/🟡/❌ extraction is easy for well-formed blocks and must fail visibly (an `unparsed` row),
   never silently, for malformed ones. Failure mode preference: show `unknown`, not omit.
3. **`parley wait` vs. driver `parley run` overlap.** If the organizer uses a foreground
   `parley run`, waiting is internal already; `wait` serves the manual/fallback path and
   post-compaction checks. Worth one consensus sentence to prevent two competing "wait"
   semantics.
4. **Where the role map lives.** Extending `meta/packet-applicability.yaml` keeps one protocol
   file; a sibling `packet-roles.yaml` keeps concerns separate. Slight preference: extend the
   existing map (one §7 surface, one file the drift story already knows).
5. **Ledger granularity.** Per phase boundary (as the brief does manually) vs. per driver step.
   Per-phase is the right default: cheap, comparable across ideas; per-step is telemetry the
   runner already emits for launched agents.
6. **kimi structured-output shape** is unverified at this writing; the parser must be built from
   captured fixtures this run, and if kimi CLI cannot emit machine-readable usage, the fallback
   is a plainly-marked `coverage: none` record — never a fabricated number.

## Risks

- **Digest-first reading degrades into summary-substitution** (a previously rejected direction).
  Mitigation: digest is deterministic, printed, path-carrying, never persisted as canonical; the
  discipline "open the raw artifact to adjudicate any disagreement" is written into SKILL.md
  core and the protocol hunk. Residual risk: behavioral (organizers over-trusting the digest) —
  partially mitigated by the `unparsed`-visible failure mode.
- **Role-scoped packet hides a section the organizer turns out to need.** Mitigation: complete
  omission index with triggers, `full` stays the default and one flag away, the role map is
  protocol-governed so enrichment is a normal §7 change. The study's own reading-set estimate
  (43–50 KB of 108.4 KB needed — SECONDARY) bounds the intended saving.
- **Prompt-repair scope creep (A).** Fixing `buildConsensusDraftPrompt` changes driver behavior
  for every deck using auto-drive; heading renames (`Trade-offs accepted` → `Agreed trade-offs`)
  could break downstream parsers that match the old headings. Mitigation: search for heading
  consumers during implementation; the driver gate itself must validate the new canonical list
  (single constant).
- **kimi telemetry may be unobtainable** if the CLI emits no structured usage; risk accepted
  with a fail-visible `coverage: none`, which is still better than the current silent 0/20.
- **Three-copy protocol drift.** The skill's bundled `references/COOPERATION.md` is unguarded at
  HEAD (verified: the CLI's drift test `internal/protocol/drift_test.go` guards the two CLI
  copies; no test covers the skill copy). This idea edits all three; I propose adding a
  content-parity check for the skill copy as implementation hygiene (tooling, not a new protocol
  obligation) — flag for consensus.
- **Handoff record becomes stale truth.** It is driver-written at transitions and advisory;
  `parley status` recomputation remains the authority. State that explicitly in its schema doc.

## Evidence & provenance (§15)

Verified at CLI-worktree HEAD `4ce8fa72a60f0d3e7818c81056e6b023aa7f36d4` (2026-09-23) and
skill-worktree HEAD `d1e57d5` (2026-09-18). Tags per §15.2; all `PRIMARY` entries were checked
by me with the quoted locator in this session.

| Claim relied on | Locator | Tag / result |
|---|---|---|
| Driver drafts consensus/FINAL via the first available headless participant | `internal/app/driver_consensus.go:82-95` (`runDrafter`) | PRIMARY — CONFIRMED |
| Consensus drafting prompt omits §15.3/§15.5/§15.6 and `Comparison & blind spots`; heading names drift from the Phase-3 template | `internal/app/driver_consensus.go:112-131` (read in full) | PRIMARY — CONFIRMED; study claim at §3 upheld |
| FINAL prompt generates its section list from `protocol.RequiredFinalSections` | `internal/app/driver_consensus.go:142-167` | PRIMARY — CONFIRMED (pattern to copy) |
| No `wait`/`digest` command exists | command dispatch, `internal/app/app.go:52-110` | PRIMARY — CONFIRMED |
| Packet renderer has no agent-role input (`Role` = deck source/consumer role) | `internal/protocolpacket/packet.go:113-122` | PRIMARY — CONFIRMED |
| Telemetry parses claude/codex/opencode/zcode; no kimi path | `internal/telemetry/usage.go` (`consume` switch) | PRIMARY — CONFIRMED; matches study §4's 0/20 kimi capture |
| Session registry exists (run_id/idea/participants/terminal) | `internal/sessionstore/sessionstore.go:13-30` | PRIMARY — CONFIRMED |
| Facilitator reading set is defined by the protocol | deck `COOPERATION.md` reading table | PRIMARY — CONFIRMED with a locator correction: the study's `COOPERATION.md:26-33` resolves to the table spanning lines 29–34 (facilitator row on line 34) at this HEAD |
| Drafter = author, else first round-1 writer/volunteer when `author: user` | `COOPERATION.md` Phase 4 (lines 405–409) | PRIMARY — CONFIRMED |
| Facilitator verification duty = existence/ownership/validity | Phase 4 pre-publish (line 437); §9 item 6 (line 890) | PRIMARY — CONFIRMED |
| Facilitator procedural calls are provisional | §15.5 (lines 1340–1343) | PRIMARY — CONFIRMED |
| Two CLI copies of COOPERATION.md + Go drift guard; skill copy unguarded | `parley-deck/COOPERATION.md` (108,400 B), `internal/protocol/defaults/COOPERATION.md`, `internal/protocol/drift_test.go`; skill `skills/parley-deck/references/COOPERATION.md` (108,244 B) | PRIMARY — CONFIRMED |
| SKILL.md is 62,411 B; its Required Protocol Context section directs full-body reads; it nowhere names `parley run`/`continue`/`status`/`consensus`/`preflight` | skill `skills/parley-deck/SKILL.md` (section read; negative search for command names) | PRIMARY — CONFIRMED (negative claim verified by search, not by absence of reading) |
| Per-agent study data and verifier corrections exist | `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/organizer-token-study/2026-09-23/` (`README.md`, `results/*.json`, `verify-*.json` — listing verified) | PRIMARY (existence) — contents SECONDARY via the copied study |
| All quantitative cost figures (566.7M input tokens, ~$1,008, 33–39% re-orientation, 9.1–12.2% polling, 43–50 KB facilitator set, etc.) | `source-context/organizer-token-study.md` §§1–5 | SECONDARY — relied on via the copied, triply-verified study; I did not recompute the underlying session logs. No acceptance criterion in this proposal depends on any single figure; they size the prize, they do not gate the design. |
| Kimi CLI structured-usage event shape | — | RECALL-adjacent: deliberately UNVERIFIED; flagged as an implementation-time discovery (Concerns §6), not relied on. |

No verdict conflicts known to me at writing time (no other round-01 file existed when I wrote
this). If a later round contradicts any PRIMARY claim above with its own locator, §15.3 applies
and I will engage the evidence, not the count.
