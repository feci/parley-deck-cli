---
idea: meta-protocol-change-lean-organizer
status: final
author: kimi-1
consensus-date: 2026-09-23
participants: [claude-1, kimi-1, zcode-1]
---

# FINAL — meta-protocol-change-lean-organizer

## Protocol context attestation (Phase-4 drafting)

```json
{"context_mode": "full", "source_sha256": "12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18", "packet_sha256": "12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18", "fallback_reason": null, "body_path": "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer/.parley-runtime/protocol-packets/full-phase4-deliberation-12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18.md"}
```

This artifact is the single source of truth for the idea. It is static and frozen at
publication; the living companion is `IMPLEMENTATION.md` (Phase 5). It is written to be
self-contained: a fresh implementer needs only this file, the two worktrees named below, and the
referenced on-disk evidence — no session transcripts.

**Owner stop rule (binding pre-Phase-5 gate, from `00-prompt.md`).** The organizer (codex-1)
inspects this FINAL before Phase 5. **If FINAL drops any of A–D, or adds a new mandatory
obligation beyond them, the organizer stops before Phase 5**, writes
`parley-deck/inbox/codex-1-to-user_meta-protocol-change-lean-organizer_scope.md` explaining the
difference, and exits. Otherwise the run continues through release without asking again. The
A–D scope & acceptance table in `## Observable acceptance criteria` is the lookup for that check:
every row's "New mandatory obligation?" column is **No**. ("Obligation" = a duty falling on any
agent or deck; repo tests and tooling are not obligations.)

**Roles (frozen at consensus).** Organizer/facilitator: codex-1 — never a participant, never
signs off, implements, or verifies code. FINAL drafter: **kimi-1** (claimed via
`inbox/kimi-1-to-all_meta-protocol-change-lean-organizer_drafter-claim.md` before signoff
completion; signoff states `Drafter: yes`). Phase-5 implementer: **the FINAL drafter (kimi-1) by
default**. zcode-1 stated willingness — "willing, not insisting" (round-02 Process) — and intends
to perfect an implementer claim via
`inbox/zcode-1-to-all_meta-protocol-change-lean-organizer_impl-claim.md` before Phase-5 work
begins; as of FINAL publication that claim file does not exist, so the default stands. If the
claim is perfected before implementation starts, zcode-1 implements and kimi-1 and claude-1 are
independent reviewers; a claim posted after implementation starts is void. claude-1 claims neither
role and remains willing if the quorum prefers. Exactly one drafter and exactly one implementer at
any time; no parallel conflicting writers.

---

## Final plan / specification

Four owner-approved scope items (A–D), each shipped as enforced tooling plus minimal protocol
text. Design rule carried from consensus: **mechanism over text** — instructions alone did not
change the organizer's token rate in the measured study, so every item lands as deterministic Go
tooling or generated artifacts, not as prose discipline.

### A — Pure organizer is the DEFAULT for declared facilitator runs

1. `00-prompt.md` gains an **optional** `facilitator: <agent-id>` field. Absent → byte-identical
   v1.48.0 behavior (regression test). The field is what makes a run a "declared facilitator run".
2. For a declared facilitator run, the default is the pure organizer: the declared facilitator
   does not implement and does not verify code; participants own drafting, implementation, tests,
   and code verification; the facilitator reads verdicts and validator output. The driver refuses
   to select the declared facilitator as drafter, implementer, reviewer, or goal-done checker, and
   **escalates rather than silently falling back** to facilitator implementation when no
   participant implementer can be launched.
3. Compatibility/role exceptions, explicit and enumerated: (i) `facilitator_participates: true`
   opts the facilitator back into participation — `parley preflight` fails closed (non-zero exit,
   naming both fields) when `facilitator:` names an agent also in `participants:` without it; with
   it, the facilitator keeps full protocol context (no audience narrowing — a signing agent never
   silently loses reading set) and role eligibility; (ii) a deck with no `facilitator:` field is
   untouched, byte-identical; (iii) the §1 user-authorized solo-exception path is unchanged;
   (iv) the owner's per-run direction is preserved — A changes the default, not the owner's
   authority.
4. **Driver drafting-prompt repair.** `buildConsensusDraftPrompt` / `buildFinalDraftPrompt`
   (`internal/app/driver_consensus.go:112-131`, `:142-167` at HEAD `ffa4587`) must emit the
   canonical Phase-3 section names and every §15 drafter duty (`## Verdict conflicts` when
   conflicts exist, `## Drafter position changes`, `## Alternatives disposition`,
   `## Comparison & blind spots`), with both section lists derived from one
   `protocol.RequiredConsensusSections` / `RequiredFinalSections` constant and a **parity test**
   proving prompt and gate read the same value. Verified basis: at HEAD `ffa4587` the file emits
   `## Trade-offs accepted` / `## Deferred follow-ups` / `## Dismissed findings` and a `grep -c`
   for the three §15 duty headings returns 0 (claude-1 signoff, PRIMARY re-verification). Scope
   basis: `00-prompt.md` constraint ("Preserve §15 dispute/provenance/alternatives duties… the old
   driver prompt may omit them"). Adds no duty to any agent — §15.5/§15.6 already bind on every
   track (§15.7).
5. Skill text states the boundary as the default and routes the organizer through the driver
   (`run` / `continue` / `wait` / `status` / `consensus` / `preflight`); hand-launching is the
   recorded fallback. Protocol text: one permissive sentence each in §4 Phase 5, §4 Phase 6, §9.0,
   and the Quickstart facilitator row — in **all three COOPERATION.md copies** (deck view
   `parley-deck/COOPERATION.md`; `internal/protocol/defaults/`; skill
   `skills/parley-deck/references/COOPERATION.md`).

### B — `parley wait` + `PhaseDigest` (blocking read, never a phase actor)

1. Extend the shipped TUI-only `BuildRoundDigest` (`internal/driver/digest.go:10-48`, built at
   `internal/driver/driver.go:488`, persisted as `round.digest` events) into a `PhaseDigest`
   covering design rounds, review rounds, consensus signoff state, and implementation status.
   Per-agent columns, all mechanically derived: `agent`, `path`, `filed`, `bytes`, `owner`
   (frontmatter `agent:`), `valid` (the shipped validators' verdict verbatim, **with the failing
   check named**, e.g. `missing required section "X"`), `stance_flags`, `unparsed`, `fell_back`.
   Validity reuses `ValidateRoundOneArtifact` / `ValidateReviewArtifact` / `ValidateFinal` — one
   source of truth.
2. Digest guardrails (all adopted): deterministic and byte-identical over an unchanged tree
   (golden test); a structural test asserts no field is model-written; fail-closed `unparsed`,
   never a guess; the next-action line is a **fixed enumeration**, never generated prose; every
   row carries its raw path; the digest is printed, never persisted as canonical (`runs/`
   telemetry capture excepted as non-canonical); any ❌ / `DISPUTED` / `unparsed` / adverse
   validity sends the organizer to the raw artifact. Raw artifacts stay canonical.
3. Command: `parley wait --idea <slug> --for round|consensus|review|implementation|any
   [--timeout D] [--json]`, blocking on the append-only event log (`internal/store/events.go`),
   with bounded portable polling (interval ≥ 10 s) where no driver event log exists —
   Windows-safe, no fsnotify dependency. **Exit codes:** `0` boundary reached; `3` timeout —
   partial digest printed, outstanding agents named; `4` a **present** artifact fails the shared
   validator (immediate, validator error verbatim) **or** a blocking escalation / `driver.error`
   event arrives (immediate); `1` usage/IO error. Missing ≠ invalid: a not-yet-filed artifact
   keeps waiting; an invalid one exits loudly. Early return on a new unanswered `to-user`
   escalation is included.
4. **Timeout:** default 25 minutes — bounded below the owner-requested 30 minutes, with **no
   universal provider TTL asserted anywhere** (the 30-minute figure is owner-brief testimony,
   `UNVERIFIED`/`RECALL` as a provider claim; no shipped code, help text, doc, or comment may
   state it as a provider fact). Configuration-first: seeded in the per-user defaults, overridable
   per call and per deck; acceptance criterion is exactly `default < 30m && configurable`; a hard
   ceiling rejects `--timeout` values above the active track's §4.0 agent timeout. The exact
   config key (reuse shipped `[defaults.timeouts]` keys vs a new seeded key) is deferred to
   implementation — both satisfy the criterion.
5. `wait` only **observes**; it never advances a phase (§14 human-brake boundary, stated in its
   docs and skill section).

### C — Audience-scoped protocol packet, slim SKILL.md, generated organizer brief

1. `parley protocol packet` gains `--audience participant|facilitator`. (Name note: `--role` /
   top-level `role` were rejected — `Role` / `source.role` already denote the deck's
   `protocolRole`, a live collision verified PRIMARY by all three participants.) The
   classification lives in an `audiences:` key in `meta/packet-applicability.yaml` — the existing
   §7-governed protocol file, one §7 surface. Hard safety property: **the audience dimension may
   only omit blocks already omittable at that phase/track/transport; it can never cut below the
   ratified never-cut floor** (`internal/protocolpacket/applicability.go:193-231`), and
   `packet check` proves it with a negative test. Attestation gains an additive `audience` field
   (never a top-level `role` key); `context_mode` values are unchanged; `full` remains the default
   for everyone; an unrecognized audience falls back to full context with a stated reason;
   `facilitator_participates: true` ⇒ full context.
2. **Facilitator verbatim reading set** (whole blocks, verbatim, with a complete omission index
   whose triggers say when to read the full source): **retain** Quickstart, §4, §5, §9, the
   active-transport §11 subsection (the protocol's own facilitator row, `COOPERATION.md:34`, table
   spanning `:29-34`), **§2** (the same row assigns the facilitator "keep the roster"), and **§15**
   — not a preference: the ratified never-cut floor pins `## 15.` and §15.1–15.4, §15.7 at kernel
   phases {1,2,3,5,6,7,8} and §15.5/§15.6 at {1,2,3,6,7}
   (`internal/protocolpacket/applicability.go:190-191,206-215`, re-verified PRIMARY at HEAD
   `ffa4587` by claude-1 at signoff), so the map cannot express its omission at deliberation
   phases. **Omit (the named omission set):** §1, §3, §8, §10, §12, §13, Appendix A, and
   non-active §11 subsections.
3. **Acceptance shape (Z-1 resolved; R-1/R-2 corrections applied — see open items 13/14).** The
   binding, gating criterion is the **named retention set + named omission set + complete omission
   index**: the facilitator body MUST contain each retained block verbatim and MUST verifiably
   omit each named-omission-set block — the named-omission-set assertion (each named block
   verified absent) is the gating check, byte-counted and recorded with the attestation. Secondary
   guardrail: **≤ 70,000 B** for the phase-1 / deliberation / github-pr facilitator body. Corrected
   framing per R-1: 70,000 B is **conservative headroom above the measured whole-section floor**,
   not an arithmetic limit — claude-1's PRIMARY measurements at HEAD `ffa4587`: named omission set
   27,420 B; floor ≈ 38.1 KB at phase 2 / ≈ 38.0 KB at phase 1 without §2; ≈ 42.1 KB with §2
   retained; `--optimize` output 65,415 B at phase 1 and 65,516 B at phase 2 (101 B apart).
   zcode-1's 43–50 KB facilitator-body target is the figure that was withdrawn (zcode-1 round-02
   position change 2, `round-02/zcode-1.md:30`); claude-1's ≤ 45,000 B secondary ceiling was never
   withdrawn and remains on record (`round-02/claude-1.md:188`, `:315`) — it is not insisted on,
   because the named-set criterion does the real work. Per R-2: the measured facilitator body is
   recorded against **both** the floor and the guardrail in the same test run, and a body under
   70,000 B that retains any named-omission-set block fails C.3 regardless of size. The
   `--optimize` baseline is measured in the same test run — what `--optimize` cuts is unknown and
   is measured, not assumed. If the §2-inclusive map-derived body exceeds the guardrail, the bytes
   are shown to the quorum before changing either the map or the ceiling; the never-cut floor is
   untouchable either way.
4. **Slim SKILL.md: core ≤ 20,000 B.** Relocation only, no rule text deleted: core keeps
   frontmatter (description verbatim — trigger reliability), Core Rule, Non-Solo, Required
   Protocol Context, Automation Mode, Startup Flow, a new driver-first section, File Ownership,
   Escalation, Quality Gates; `references/` absorbs headless-launch configuration, generic CLI
   invocation contract, timeout policy, artifact templates, roster / global-core / drift detail.
   The core names `parley run`, `continue`, `wait`, `status`, `consensus`, `preflight` at least
   once each (the direct inverse of the verified v2.12.1 state, which names none) and links every
   reference file; a skill-repo script proves every moved heading lands in exactly one reference
   and none was dropped. If 20,000 B genuinely cannot fit without dropping a rule, the bytes come
   back to the quorum — the cap is not silently raised.
5. **`parley organizer brief --idea <slug>`: computed, never stored.** Content: packet attestation
   + facilitator body path, `parley status --idea --json` state, current `PhaseDigest`, driver's
   next action, phase pointer from D. Contract: ≤ 8,192 B, byte-identical across two runs over an
   unchanged tree, writes no file (asserted against a read-only deck). A stored brief is a rejected
   sidecar pattern; any drift toward persisting it is a stop moment. Post-compaction
   re-orientation = generated brief + facilitator packet + `parley status` — documented as the
   lean loop in the slim skill core.

### D — Fresh session per phase + client-accounting measurement

1. **State/handoff FILE form is retained and supported** (owner requirement for this consensus):
   the driver writes a per-phase handoff record under `runs/<run-id>/` at each phase transition,
   built on the shipped handoff machinery (`internal/runner/handoff.go` `WriteHandoffPacket`) and
   sharing the PhaseDigest computation. It is non-canonical driver state (§3 `runs/`), advisory,
   and its schema doc states the recomputed view (`organizer brief` / `parley status`) is
   authoritative on any disagreement. **No agent is ever obliged to author or refresh a handoff
   file by hand** — no new manual-refresh obligation. The user-facing mandatory acceptance
   criterion is the computed brief's output contract (C.5); the `runs/` record is covered by a
   driver unit test as repo tooling.
2. **Organizer usage ledger:** `parley usage ingest --agent <id> --source
   codex-rollout|claude-jsonl --path <file> --idea <slug> --phase <n>`. Binding properties: **no
   flag accepts a token count** — the tool accepts a path and parses the client's own accounting
   file itself (the mechanical answer to the earlier rejection of "self-reported usage records as
   cost evidence": a number typed by a model is `RECALL` and inadmissible; a number the tool parses
   from the client-written file is `PRIMARY` with a locator); streaming with bounded memory over
   the verified 228 MB Codex rollout; ≤ 1 KB stdout; one ledger row per ingest carrying idea,
   phase, agent, source path, parser id, ingest time, and the six `total_token_usage` fields
   verbatim — tokens only, no dollar claim; idempotent re-ingest; round-trip re-parse equality
   test.
3. **Attribution is explicit and bounded — never slug scanning:** records attribute by the
   explicit `--idea` / `--phase` arguments plus run-record timestamp windows; residue the windows
   cannot resolve is labeled `ambiguous`, with the method stated in the ledger header. A slug
   appearing in a session file is never attribution evidence — it appears for many reasons
   (discussion, grep, pasted briefs), and this run's own `organizer-usage.md` accounting correction
   (ambiguous source selector picked the wrong rollout; superseded snapshots) is same-class
   operational testimony.
4. **Participant telemetry:** add the `kimi` case to `internal/telemetry/usage.go` (three `case`
   labels exist at HEAD `ffa4587` — `claude:230`, `codex:252`, `opencode:261`; no kimi path; zcode
   is envelope-handled near `:190`/`:286` — the resolved verdict conflict, see consensus.md
   `## Verdict conflicts`), built from fixtures captured live this run
   (`source-context/kimi-readiness.jsonl` is testimony that a structured envelope exists);
   structured-output argv defaults only where the adapter supports them without breaking `-p`
   semantics. If kimi emits no machine-readable usage, record a visible `coverage: none` — never a
   fabricated number.
5. Protocol text: one permissive §9 line recording that a facilitator may re-orient from the
   generated brief instead of re-reading SKILL.md and the full COOPERATION.md, and one §11
   advisory line preferring one blocking `parley wait` over repeated short polls — all three
   copies.

### Cross-cutting decisions (binding)

- **Repos in scope:** the CLI worktree
  (`/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer`, branch
  `lean-organizer`, base tag v1.48.0, HEAD `ffa4587` at consensus) **and** the sibling skill
  worktree (`/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer-skill`,
  branch `lean-organizer`, base tag v2.12.1, HEAD `d1e57d5` at FINAL drafting). Protocol text
  lands in **all three COOPERATION.md copies** (deck view; `internal/protocol/defaults/`; skill
  `skills/parley-deck/references/COOPERATION.md`). The CLI Go drift test keeps guarding the two
  CLI copies; a skill-repo content-parity check for the third copy is added as **repo tooling, not
  a protocol obligation**.
- **Global core:** staged from the previous core TEMPLATE plus the reviewed hunks (placeholder
  header, stub §2 — never a deck-view copy), placed in `~/.parley/staging/`, with an escalation to
  the owner carrying the exact command. **`parley protocol publish` is the owner's attended-only
  action** (TTY-gated); no participant or organizer works around the gate.
- **Release (after Phase 8 completes, organizer's step):** CLI 1.49.0, skill 2.13.0, CHANGELOGs,
  tags, GitHub releases with Windows assets, both Homebrew formulae in homebrew-parley, two winget
  PRs (one application per PR against microsoft/winget-pkgs), `npm publish --access public` of the
  exact packed tarball (if npm asks for web verification, escalate the URL to the owner via the
  inbox and wait), skill install `parley-deck-skill install --target all --force` verified by
  content hash of every runtime SKILL.md — per the owner brief; direct-to-main under the owner's
  run-scoped override (no development PRs; the deck's transport header is unchanged). The deploy
  is independently **participant-verified** against real channel artifacts (both formulae
  url+sha256; `$(which parley-deck-skill)` resolves into Cellar; npm latest + integrity; GitHub
  assets + sha256; winget PR state; installed SKILL.md content) before the release is reported
  complete; findings are fixed first.
- **Driver/runtime gaps observed this run are evidence, not automatic new scope:** the kimi
  stream-json preflight parser rejection, the round-01 heading strictness + `parley status`
  validity concealment, the `continue --auto` fresh-run startup gap, and the historical-worktree
  registration blocker (with its mislabeled escalation author) are recorded organizer operational
  testimony (`organizer-notes.md`). B's validity-with-reason column covers the concealment class
  inside scope B. Any further repair is optional hardening inside the release's driver work,
  re-verified at HEAD during implementation — nothing here adds a mandatory obligation or an
  unrelated process gate.

### Alternatives disposition (§15.6(c) — FINAL may not contradict a recorded adoption)

Adopted (design constraints on the implementer): **ALT-1** reuse the `parley status --json` /
`parley consensus status --json` parsing inside the digest (wrap in a blocking verb; do not change
`status`); **ALT-3** build the audience dimension on the existing packet renderer + applicability
map (one §7-governed surface); **ALT-5** build D on `internal/sessionstore` +
`WriteHandoffPacket`; **ALT-6** add the kimi case to the shipped telemetry switch (no new
subsystem); **ALT-7** extend the skill's existing `references/` on-demand pattern; **ALT-2**
foreground `parley run` adopted only as a complement, never as the wait primitive.

Rejected (do not reintroduce): **ALT-2** (as wait primitive — a driver cycle is not a passive
read), **ALT-4** (`--optimize` packet as the organizer view — experimental, participant-oriented,
measured 65,516 B still carrying non-facilitator sections), **ALT-8** (all previously rejected
directions upheld: LLM/lossy summarization, stored summary sidecars, retrieval-selected protocol
text, asymmetric reviewer context, silent model swaps, self-reported usage as cost evidence),
**ALT-9** (hand-written stored `handoff.md` — stale sidecar + manual-refresh obligation),
**ALT-10** (TUI digest — the organizer is headless), **ALT-11** (cutting inside §4 to chase
43–50 KB), **ALT-12** (shortening the driver's internal poll interval — invisible to the
organizer), **ALT-13** (enforcing A by roster removal — wrong scope), **ALT-14** (making `status`
blocking — breaks scripted callers), **ALT-15** (`--role` / top-level `role` — live collision),
**ALT-16** (usage attribution by slug mentions), **ALT-17** (hand-written static organizer brief),
**ALT-18** (25-minute timeout as a baked Go constant — configuration-first instead). Full
reasoning: `consensus.md` `## Alternatives disposition`.

### Open items deferred to implementation

Carried verbatim from `consensus.md` (items 1–12), plus the two accepted signoff reservations
logged there (items 13–14):

1. Exact configuration key for the `wait` default (reuse shipped `[defaults.timeouts]` keys vs a
   new seeded key) — must satisfy `default < 30m && configurable` plus the per-track ceiling.
2. Whether the §2-inclusive facilitator map fits the ≤ 70,000 B guardrail at each phase; if not,
   show the bytes to the quorum before changing the map or the ceiling. The never-cut floor is not
   negotiable. (Extended by R-2 — item 14.)
3. What `--optimize` actually cuts: the audience × optimizer interaction is measured during
   implementation; the acceptance-criterion baseline is captured in the same test run, not
   transcribed from any round file.
4. Kimi structured-usage event shape (fixtures captured this run); honest `coverage: none` path
   if unobtainable.
5. Signoff/stance parsing: verify reuse of the consensus-status parser or extract shared parsing —
   do not fork it (a code-review checklist line, not a protocol obligation).
6. Heading-consumer search for the consensus-prompt repair (downstream parsers matching old
   headings).
7. Skill-copy parity check shape (extend the skill repo's manifest-coverage test pattern vs a
   `compatibility.json` sha — implementation choice).
8. Exact `references/` file names and section allocation for the SKILL.md split (claude-1's
   HEADLESS_LAUNCH / ARTIFACT_TEMPLATES / ROSTER_AND_PROTOCOL sketch is the starting point).
9. Handoff-record schema and exact `runs/` location; schema doc wording stating recomputation is
   authoritative.
10. Machine-readable ledger sibling format alongside `organizer-usage.md`.
11. Whether `parley preflight` soft-warns when no `facilitator:` field is declared. Default: **no
    warning** (zcode-1 confirmed its lean-yes was advisory).
12. `## Alternatives disposition` consistency duty: this FINAL must not contradict any adoption
    recorded in `consensus.md` (§15.6(c)); a contradiction blocks signoff and escalates to the
    owner.
13. **R-1 (claude-1 reservation, accepted by kimi-1, concurred by zcode-1 — record correction,
    applied in this FINAL):** the corrected attribution/arithmetic framing of C.3 — zcode-1's
    43–50 KB target is the withdrawn figure; claude-1's ≤ 45,000 B secondary ceiling was never
    withdrawn; the ≤ 70,000 B guardrail is conservative headroom above the measured ≈ 42 KB
    whole-section floor (claude-1 PRIMARY at HEAD `ffa4587`), not an arithmetic limit. kimi-1's
    §15.1 `SELF-CORRECTION` in `consensus.md` is the canonical correction record. Disposition:
    **carried into C.3 above.**
14. **R-2 (claude-1 reservation, accepted by kimi-1, concurred by zcode-1 — extends open item
    2):** the named-omission-set assertion (each named block verified absent) is the **gating**
    check in C.3; the facilitator body is recorded against both the floor and the guardrail in the
    same test run; a body under 70,000 B that retains any named-omission-set block fails C.3
    regardless of size. The never-cut floor stays untouchable. Disposition: **carried into C.3
    above and into the acceptance table below.**

---

## Purpose / user-visible outcome

Reduce organizer (facilitator) token usage while preserving canonical participant ownership and
verification integrity. The organizer is usually the strongest and most expensive model on the
roster; the measured study (`source-context/organizer-token-study.md`) shows its cost is driven by
re-reading resident context — self-implementation/self-verification (~30-40% of input),
post-compaction re-orientation (~33-39%), compaction itself (42% of uncached input), and polling
(~9-12%) — and that instructions alone did not change the rate.

User-visible outcomes after this ships:

- A deck that declares `facilitator:` gets a pure organizer by default: the driver never assigns
  it drafter/implementer/reviewer/goal-done roles and escalates instead of silently falling back.
- `parley wait` — one blocking call replacing poll loops, returning a deterministic, Go-generated
  per-agent digest with validator verdicts and exit codes an organizer can branch on.
- `parley protocol packet --audience facilitator` — a verbatim, omission-indexed,
  §7-governed role-scoped protocol view for cheap post-compaction re-orientation.
- A slim SKILL.md core (≤ 20,000 B) that routes organizers through the driver by name.
- `parley organizer brief` — a computed, never-stored re-orientation surface.
- `parley usage ingest` — a per-idea, per-phase organizer usage ledger built from the client's own
  accounting files (no self-reported numbers), plus repaired kimi participant telemetry.
- Protocol text recording the pure-organizer default and the lean re-orientation loop, identical
  in all three COOPERATION.md copies, staged for the owner's attended core publish, and released
  as CLI 1.49.0 / skill 2.13.0 on all channels.

---

## Context & orientation

**Idea state.** Track `deliberation`; quorum locked at Phase 0: claude-1, kimi-1, zcode-1.
codex-1 is organizer only (never participant, signoff, implementer, or code verifier). Consensus
reached 2026-09-23: claude-1 🟡 ACCEPT-WITH-RESERVATIONS (R-1, R-2 — both accepted and carried
into this FINAL), kimi-1 ✅ ACCEPT (`Drafter: yes`, with the §15.1 SELF-CORRECTION), zcode-1 ✅
ACCEPT (concurring with R-1/R-2). No ❌. One verdict conflict (telemetry adapter case list) was
resolved by grep evidence — claude-1's PRIMARY verdict (`grep -n 'case "'
internal/telemetry/usage.go` → `claude:230`, `codex:252`, `opencode:261`); it is RESOLVED, not
DISPUTED, and no decision or acceptance criterion depends on the phrasing. Full record:
`consensus.md` `## Verdict conflicts`.

**Worktrees and evidence:**

- CLI worktree: `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer`
  (branch `lean-organizer`, base v1.48.0, HEAD `ffa4587` at consensus). The deck is its
  `parley-deck/`. Test suite: `go build ./... && go test ./...`.
- Skill worktree:
  `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer-skill` (branch
  `lean-organizer`, base v2.12.1, HEAD `d1e57d5` at FINAL drafting). Test suite: `npm test`
  (`node --test && node scripts/run-python-tests.js && node scripts/build-addon-manifest.js
  --check`).
- Verified study (six readers, three independent verifiers): `source-context/organizer-token-study.md`;
  per-agent raw data at `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/organizer-token-study/2026-09-23/`.
- This run's operational evidence (driver gaps, manual fallbacks, usage accounting):
  `organizer-notes.md`, `organizer-usage.md`, `source-context/kimi-readiness.jsonl`,
  `source-context/preflight.json`, `source-context/machine-roster.json`, and
  `inbox/claude-to-user_meta-protocol-change-lean-organizer_driver-error.md`.

**Key code locators at HEAD `ffa4587` (re-verify at implementation HEAD):** driver consensus/FINAL
prompts `internal/app/driver_consensus.go:112-131`, `:142-167`; TUI digest
`internal/driver/digest.go:10-48` (built at `internal/driver/driver.go:488`); event log
`internal/store/events.go`; packet applicability + never-cut floor
`internal/protocolpacket/applicability.go:190-231`; packet role collision
`internal/protocolpacket/packet.go:113-122` (+ `source.go`); telemetry switch
`internal/telemetry/usage.go:229-298` (structured-flag heuristic `:380-395`,
`internal/runner/telemetry.go:147`); kimi plain `-p` argv `internal/agents/discover.go:331-358`;
handoff machinery `internal/runner/handoff.go` `WriteHandoffPacket`; applicability map (protocol,
§7-governed) `parley-deck/meta/packet-applicability.yaml`.

**Correlated-agreement statement (§15.6(b)).** The three participants are three model families,
but their unanimity is a **shared prior, not independent evidence**: all three read the same
study, the same two worktrees, and share a visible common prior that deterministic Go tooling is
trustworthy and LLM summarization is suspect (study §6, unchallenged by anyone). The three designs
were never independent mechanisms — they are one family (route the organizer through existing
driver / renderer / digest machinery) with per-item variations, and round-02 explicitly rebuilt
two of them onto claude-1's framing. What would make the agreed position wrong: the digest
becoming a de-facto substitute for artifacts; the audience map omitting a block an organizer
operationally needs despite `packet check`; measured savings failing because the out-of-scope
skills catalogue dominates; kimi CLI proving unable to emit structured usage (D3 degrades to
`coverage: none`); the lean loop not being adopted behaviorally; the 30-minute cache guidance
being wrong in the short direction.

**Known blind spots flagged for implementation (measure, do not silently assume):** (i) §15.5/§15.6
are never-cut at phases {1,2,3,6,7} but **not pinned at phases 5 and 8** — check what a
facilitator view loses there and whether it matters; (ii) no participant verified whether the
skill repo has a Windows CI leg today; (iii) the real provider cache window was never measured —
the 25-minute default is owner-guidance heuristic; (iv) the `organizer-usage.md` same-day
accounting correction is testimony motivating explicit attribution (D.3), not a verified defect in
a shipped tool.

---

## Observable acceptance criteria

**Owner stop-rule lookup table** (per claude-1's accepted drafting ask N-4: the organizer's
pre-Phase-5 check is a lookup on this table, not a judgment call). Every acceptance line is a
command with an exit code, a byte bound on a named file, a determinism assertion, or a negative
test. Product acceptance tests implementing A–D are normal refinements of this table; nothing here
adds a mandatory obligation beyond owner-approved A–D, and no unrelated process gate is created.

| Item | Ships (mechanisms) | Acceptance (observable, code-reviewable) | New mandatory obligation? |
|---|---|---|---|
| **A** — organizer does not implement / verify code, default for declared facilitator runs | optional `facilitator:` field; preflight fail-closed + `facilitator_participates:` exception; driver role-ineligibility predicate with escalate-not-fallback; `RequiredConsensusSections` prompt repair + parity test; skill boundary text; permissive protocol lines in all 3 copies | `parley preflight` exits non-zero naming both fields when `facilitator:` ∈ `participants:` without the exception flag; exit 0 with it; fixture auto-drive run launches the declared facilitator for no drafter/implementer/reviewer/goal-done role (event-log assertion); absent-field deck → byte-identical run plan (regression test); emitted prompt contains canonical Phase-3 sections + §15.3/15.5/15.6 strings, derived from the shared constant (parity test); this run's `IMPLEMENTATION.md` records `implementer:` ≠ facilitator | **No** — optional field, permissive text, repo tests |
| **B** — `parley wait` + PhaseDigest | extend shipped digest; new read-only blocking verb; shipped validators reused; event-log wait with portable polling fallback | exit 0 with per-agent `path/filed/bytes/owner/valid(with failing check)/stance/unparsed/fell_back`; missing file → blocks to timeout exit 3, partial digest names outstanding agents; present-but-invalid → immediate exit 4 with verbatim validator error; unanswered `to-user` escalation / `driver.error` → exit 4; byte-identical digest on unchanged tree (golden test); structural no-model-written-field test; `default < 30m && configurable`, ceiling rejects `--timeout` above the active track's §4.0 timeout; digest never rewrites or replaces a round file (regression test); round-trip: one organizer tool call between launch and digest under `--no-tui` | **No** — an alternative to polling, never a required step |
| **C** — audience packet + slim skill + generated brief | `--audience` on the existing renderer; `audiences:` key in the §7-governed applicability map; SKILL.md relocation split; computed `parley organizer brief` | `packet check` passes with `audiences:` and fails any audience rule omitting a never-cut block (negative test); **gating check (R-2): each named-omission-set block (§1, §3, §8, §10, §12, §13, Appendix A, non-active §11) verified absent — a body under 70,000 B retaining any named-omission-set block fails regardless of size**; facilitator body contains the named retention set verbatim (Quickstart, §4, §5, §9, active §11, §2, §15), complete omission index with triggers, `source_sha256` over the full authority, additive `audience` attestation field; map-derived byte count recorded against **both** the measured floor (≈ 42 KB with §2) and the ≤ 70,000 B guardrail in the same test run, `--optimize` baseline captured in that run; unknown audience → full fallback with reason; `facilitator_participates: true` → full context; SKILL.md core ≤ 20,000 B, names the six driver commands, relocation script proves nothing dropped, frontmatter/Core Rule verbatim; brief ≤ 8,192 B, byte-identical ×2, writes no file (read-only-deck test) | **No** — opt-in audience, generated views |
| **D** — fresh session per phase + client-accounting ledger + telemetry | driver-written per-phase handoff records under `runs/` (non-canonical, advisory); `parley usage ingest` path-only streaming ingest; kimi telemetry parser | handoff record written at each phase transition with schema-valid fields (driver unit test); schema doc states recomputation authoritative; **no count-accepting flag exists** (flag-set test); 228 MB fixture streams with bounded memory, ≤ 1 KB stdout, exactly one ledger row with idea/phase/agent/source-path/parser-id/six-fields-verbatim; idempotent re-ingest; round-trip re-parse equality; attribution by explicit args + timestamp windows, residue labeled `ambiguous`, method stated in ledger header; kimi fixture → non-nil usage or visible `coverage: none`; structured argv only where adapter-supported (test) | **No** — driver-generated records; path-only tool; no manual refresh duty |
| **Cross-cutting** | all three protocol copies; CLI drift test; skill-copy parity check (tooling); staged core from previous template + reviewed hunks in `~/.parley/staging/`; owner-only attended `parley protocol publish` with exact-command escalation; release CLI 1.49.0 / skill 2.13.0 all channels, direct-to-main per owner override; participant-verified deploy | both worktrees' test suites green incl. the new tests; the three copies carry identical hunks; staged core exists with project zones preserved (placeholder header, stub §2); escalation note contains the exact attended publish command; post-release channel verification findings fixed before completion is reported | **No** — publish is owner-only by existing mechanism; checks are repo tooling |
| **Evidence, not scope** | startup parser gap, heading-strictness/status concealment, stale-worktree registration blocker recorded as testimony; B's validity-reason column covers the concealment class | mentioned in FINAL's risks/context (done — see `## Known risks / de-risking`); any repair is optional hardening re-verified at HEAD in implementation | **No** — explicitly not mandatory new scope |

### Validation plan (per worktree, per criterion)

The implementer validates each criterion with the named command class, and records results in
`IMPLEMENTATION.md` `## Validation evidence` with the commands run and what they proved; reviewers
re-run or inspect them in Phase 6 (refutation-default: try to break each criterion).

**CLI worktree** (`cd <cli-worktree> && go build ./... && go test ./...` — all green, including
new tests):

- A: preflight fail-closed/exception exit-code tests; driver role-ineligibility test asserting the
  declared facilitator is never selected as drafter/implementer/reviewer/goal-done checker and
  that the no-participant-implementer case escalates (event-log assertion); absent-field
  byte-identical run-plan regression test; consensus/FINAL prompt parity test against
  `protocol.RequiredConsensusSections` / `RequiredFinalSections` (prompt and gate read one
  constant); heading-consumer search documented in `IMPLEMENTATION.md`.
- B: digest golden test (byte-identical over an unchanged tree); structural test that no digest
  field is model-written; exit-code tests for 0 / 3 (partial digest names outstanding) / 4
  (present-but-invalid with verbatim validator string; escalation/`driver.error`) / 1; missing ≠
  invalid test; timeout default + ceiling rejection test; polling fallback test without fsnotify
  (Windows-safe); digest-never-rewrites regression test.
- C: `packet check` suite with `audiences:` including the negative never-cut test; named retention
  set verbatim-in-body test; **named-omission-set absence test (the gating assertion)**; byte
  counts (floor, guardrail, `--optimize` baseline) emitted by the same test run and recorded;
  unknown-audience full-fallback test; attestation additive-`audience` test; organizer-brief
  contract tests (≤ 8,192 B, byte-identical ×2, read-only-deck no-write).
- D: handoff-record driver unit test (schema-valid fields at each transition); usage-ingest
  flag-set test (no count-accepting flag); streaming test over the 228 MB Codex rollout fixture
  with bounded memory and ≤ 1 KB stdout; one-row ledger shape test (six `total_token_usage`
  fields verbatim); idempotent re-ingest and round-trip re-parse equality tests; attribution
  window + `ambiguous` residue test; kimi telemetry parser test from this run's captured fixtures
  (`coverage: none` path included); adapter-argv support test.
- Cross-cutting: the existing Go drift test stays green over the two CLI protocol copies;
  Windows CI leg covers `wait`/`usage`.

**Skill worktree** (`cd <skill-worktree> && npm test` — all green, including new tests):

- SKILL.md relocation script: every moved heading lands in exactly one `references/` file, none
  dropped; core ≤ 20,000 B byte test; core names `parley run`, `continue`, `wait`, `status`,
  `consensus`, `preflight`; frontmatter description and Core Rule verbatim; every reference file
  linked from core.
- Content-parity check for the third COOPERATION.md copy (`skills/parley-deck/references/`) —
  repo tooling, shape per open item 7.

**Cross-repo:** the three protocol copies carry identical hunks (diff-visible); staged core exists
in `~/.parley/staging/` built from the previous core template + reviewed hunks; the owner
escalation carries the exact attended `parley protocol publish` command. If a byte bound genuinely
cannot be met (facilitator guardrail, skill core cap), the bytes come back to the quorum before
any map/ceiling/cap change — the bound is not silently raised.

---

## Idempotence & recovery

**State that matters:** the two worktree git branches (`lean-organizer` in each worktree); this
idea's directory under `parley-deck/ideas/`; the staged core under `~/.parley/staging/`; release
channel state (tags, GitHub releases, Homebrew formulae, winget PRs, npm package, installed
skills). `runs/` content (handoff records, telemetry) is non-canonical driver state — advisory
only; the recomputed view (`organizer brief` / `parley status`) is authoritative on any
disagreement.

**Safe to rerun at any time:** both test suites; `parley wait` (read-only observe); `parley
organizer brief` (computed, writes nothing); `parley protocol packet` (read-only render);
`parley usage ingest` (idempotent re-ingest by design — same input file, same ledger row);
digest/packet/brief generation (deterministic, byte-identical over an unchanged tree).

**Human / owner gates (never worked around):** `parley protocol publish` is the owner's
attended-only TTY-gated action — the run ends with an escalation carrying the exact command and
the staged core path; the release step (merges to main, tags, channel publishes) is the
organizer's step after Phase 8 completes, under the owner's standing authorization for this run;
participants never release, publish core, or integrate to main before that step. If npm publish
requests web verification, escalate the URL to the owner via inbox and wait.

**Recovery:** implementation deviations from this FINAL are logged in `IMPLEMENTATION.md`
`## Deviations from FINAL.md`, never silently absorbed. Fix-up cycles (Phase 8, cap 5 on
`deliberation`) apply agreed fixes on the same branches. This FINAL is frozen: if it is later
invalidated, open `meta-protocol-change-lean-organizer-v2` — do not edit this file. Pre-merge
mistakes are ordinary git work on the feature branches; post-release problems follow the normal
revert/patch-release path of each repo. A failed or partial channel publish is reconciled against
the real channel artifact (the participant-verified deploy checklist), not by re-running blindly.

---

## Known risks / de-risking

1. **Savings land mostly at cached-read rates** — polling tokens are ~54%-of-cost cached reads;
   the token share overstates the cost share. De-risking: D's ledger measures the real effect per
   phase instead of asserting it.
2. **The largest measured lever is out of scope** — the ~22 KB skills catalogue in every request
   and the AGENTS.md read-in-full lines (study §7) are owner-environment levers A–D cannot reach.
   If the release is measured against total organizer spend it may look like a shortfall for
   reasons never claimed; the ledger makes that visible rather than arguable.
3. **Digest-first over-trust** — an organizer may stop opening raw artifacts. Mitigations:
   fail-closed `unparsed`, loud exit-4 degradation, path-carrying rows, fixed-enumeration
   next-action line, the written adjudication discipline (any ❌/`DISPUTED`/`unparsed`/adverse
   validity → raw artifact), and the residual-risk statement in the slim skill core.
4. **Byte guardrail slack masks a retained block** — 70,000 B sits ~28 KB above the measured
   ≈ 42 KB floor, so size alone cannot detect a body carrying an omission-set block. Mitigation
   (R-2, gating): the named-omission-set absence assertion gates C.3; body recorded against floor
   and guardrail in the same run.
5. **25-minute default is a heuristic** — if a real provider cache window is shorter, a wait may
   straddle expiry; cost is one re-invocation, mitigated by exit-3 partial digests naming
   outstanding agents. No artifact asserts a provider TTL.
6. **Consensus-prompt repair changes driver-emitted headings for every auto-drive deck** —
   mitigated by the heading-consumer search during implementation and the prompt/gate parity test.
7. **Kimi structured usage may be unobtainable** — accepted: a visible `coverage: none` beats a
   fabricated number; the D criterion passes either way.
8. **Windows support honesty** — a Windows CI leg is added for `wait`/`usage` so the release ships
   Windows assets honestly; whether the skill repo has a Windows CI leg today is unverified (check
   at implementation).
9. **Behavioral adoption risk** — tooling can exist and still not be used; the slim skill core
   routes the organizer through the driver by name, and the generated brief is the documented
   re-orientation surface.
10. **Observed driver/runtime gaps this run** (evidence, not scope): kimi stream-json preflight
    parser rejection; round-01 heading strictness with `parley status` concealing the validity
    reason; `continue --auto` fresh-run startup gap; historical-worktree registration blocker with
    a mislabeled escalation author. B's validity-with-reason column covers the concealment class;
    any further repair is optional hardening, re-verified at HEAD during implementation.
11. **Role-concentration / process risk** — codex-1 (organizer) is a non-participant; kimi-1
    drafted consensus and FINAL as a participant holding no facilitator role, so no
    facilitator-drafter concentration exists; §15.5/§15.6 drafter duties were discharged in
    `consensus.md` (`## Drafter position changes`, `## Alternatives disposition`, `## Comparison
    & blind spots`, correlated-agreement statement) and are carried into this FINAL.

---

## References

- Consensus (decisions, trade-offs, open items, verdict conflicts, drafter position changes, all
  signoffs incl. R-1/R-2 and the §15.1 SELF-CORRECTION): `./consensus.md`
- Rounds: `./round-01/` (claude-1.md, kimi-1.md, zcode-1.md), `./round-02/` (claude-1.md,
  kimi-1.md, zcode-1.md)
- Owner kickoff + approved brief (verbatim, incl. A–D and the stop rule): `./00-prompt.md`
- Verified evidence study: `./source-context/organizer-token-study.md`; raw per-agent data:
  `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/organizer-token-study/2026-09-23/`
- Run evidence: `./organizer-notes.md`, `./organizer-usage.md`,
  `./source-context/kimi-readiness.jsonl`, `./source-context/preflight.json`,
  `./source-context/machine-roster.json`
- Claims/escalations: `../../inbox/kimi-1-to-all_meta-protocol-change-lean-organizer_drafter-claim.md`;
  `../../inbox/claude-to-user_meta-protocol-change-lean-organizer_driver-error.md`
- Live protocol (authority, attested above): `../../COOPERATION.md` (deck view);
  applicability map `../../meta/packet-applicability.yaml`
- CLI worktree: `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer`
  (branch `lean-organizer`, base v1.48.0)
- Skill worktree:
  `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer-skill` (branch
  `lean-organizer`, base v2.12.1)
