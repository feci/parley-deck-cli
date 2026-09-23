---
agent: zcode-1
idea: meta-protocol-change-lean-organizer
round: 1
date: 2026-09-23
---

## Protocol context attestation (§9.0 item 1 / SKILL.md "Required Protocol Context")

```json
{"context_mode": "full", "source_sha256": "12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18", "packet_sha256": "12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18", "fallback_reason": null, "body_path": "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer/.parley-runtime/protocol-packets/full-phase1-deliberation-12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18.md"}
```

I read the emitted body in full (1,380 lines) and independently recomputed its SHA-256; it matches
both `source_sha256` and `packet_sha256` above. Active transport noted: `github-pr`, with the owner's
run-scoped override to direct-main integration recorded in `00-prompt.md`. Provenance tags below use
`@CLI` = CLI worktree HEAD `4ce8fa7` and `@SK` = skill worktree HEAD `d1e57d5`.

## Summary

The study's central finding — organizer cost is dominated by resident-context re-reads (post-compaction
re-orientation ~33–39% of input, compaction 42% of uncached input), not by output — means the effective
levers are: keep the organizer out of implementation/verification (A), make waiting and state-reading
one bounded blocking call with a deterministic digest instead of polling loops (B), cut what a fresh or
compacted organizer session must read to a role-scoped, attested protocol view plus a slim skill core
(C), and stop compounding context across phases while measuring the result from client accounting (D).
I propose concrete CLI/skill/protocol changes for all four: a new `parley wait` command with a
byte-stable Go-generated digest, a `--role facilitator` input on the existing packet renderer, a
SKILL.md split into a ≤20 KB core plus `references/`, a phase-bootstrap status shape, a `parley usage
record` tool for the organizer ledger, and minimal, permissive protocol text deltas in all three
COOPERATION.md copies. No new mandatory obligations are introduced; raw artifacts stay canonical.

## Proposed approach

### Design principles

1. **The organizer's tokens are saved by changing what it must hold in context, not what it may read.**
   Every mechanism below is an index or a routing change; the canonical artifacts are untouched and
   remain the audit trail. Nothing is an LLM summary (the prior rejection stands; a deterministic
   structured extraction is not summarization, and it never substitutes for adjudication reading).
2. **Defaults change only for the facilitator role.** Participant context pipeline (full protocol
   prepended per launch) is out of scope and stays the ratified default. This idea is the §7 change
   the study itself scopes to the organizer role.
3. **Fail closed everywhere.** A digest that cannot parse a stance prints `unparsed`, never a guess;
   a usage attribution that cannot be resolved is labeled ambiguous, never split arbitrarily; a kimi
   usage shape that does not exist is recorded as a gap, never faked.
4. **No new mandatory obligations beyond A–D** (owner stop rule). Protocol text is permissive-default:
   it changes what the facilitator does when participants are available, and never binds a deck that
   cannot staff them.

### A. Organizer does not implement or verify code

Today the protocol already permits full delegation: the drafter is the initiator or the first round-1
writer/volunteer [PRIMARY: `parley-deck/COOPERATION.md:407-409` @CLI], the default implementer is the
FINAL drafter unless claimed [PRIMARY: `parley-deck/COOPERATION.md:443-446` @CLI], facilitator calls
are provisional [PRIMARY: `parley-deck/COOPERATION.md:1342` @CLI], and the facilitator's verification
duty is existence/ownership/validity [PRIMARY: `parley-deck/COOPERATION.md:439` and `:890` @CLI]. The
driver already drafts consensus/FINAL through the first available headless participant
[PRIMARY: `internal/app/driver_consensus.go:82-95` @CLI]. What is missing is that none of this is the
stated default in the skill or enforced by the driver.

Changes:

- **Driver (`internal/app/driver_impl.go` and Phase 5–8 paths):** when `auto_implement` is set (or a
  run reaches Phase 5), the implementer agent is resolved as a *participant* (FINAL drafter, else the
  first `inbox/` claimant, else first available participant) — never the facilitator session running
  the driver. Code verification remains participant reviewers plus the machine gates that already
  exist (LE-4 `checks:` validation, LE-7/11 goal-done check via a fresh non-implementer). The
  organizer consumes: review files, `## Validation evidence` tables, and driver validator output. If
  no participant implementer can be launched, the driver escalates (existing §1 non-solo behavior) —
  it does not silently fall back to facilitator implementation.
- **Skill:** the facilitator sections state the boundary explicitly — the facilitator does not
  implement or verify code when a participant implementer exists; it reads verdicts and validator
  output. The organizer-facing loop routes through `parley run` / `continue` / `wait` / `status` /
  `consensus` (these commands exist [PRIMARY: `internal/app/app.go:60-108` dispatch @CLI] but the
  current SKILL.md never mentions them — verified: `grep -nE "parley (run|continue|consensus|status|preflight)"`
  over `skills/parley-deck/SKILL.md` returns nothing @SK).
- **Protocol text (all three copies):** one sentence in §4 Phase 5 and the Quickstart role table's
  Facilitator row: when the facilitator is not a participating implementer, implementation and code
  verification are participant work; the facilitator reads verdicts and validator output. Framed as
  the default when participants are available — an owner may still direct otherwise; no new
  obligation is created.

### B. `parley wait` + digest

No `wait` subcommand exists today [PRIMARY: `internal/app/app.go:60-108` @CLI has no `wait` case].

- **Command:** `parley wait --idea <slug> --for <condition> [--timeout 25m] [--interval 20s] [--json]`.
  Conditions: `round-NN` (every `participants:` entry has an artifact that passes the existing driver
  round validation), `consensus` (all signoff blocks appended), `review-round-NN`, `complete`
  (`IMPLEMENTATION.md` `status: complete`).
- **Blocking:** a single foreground process that polls the filesystem internally at `--interval`
  (portable polling, not fsnotify — the release ships Windows assets; CLI-side polling cost is
  negligible; the token saving is the organizer's). Default timeout 25 m, deliberately below the
  ~30 m provider prompt-cache lifetime cited in the brief; configurable.
- **Exit codes:** `0` condition met, `2` timeout (digest still printed, showing who is outstanding),
  `1` error (bad idea, unreadable state).
- **Digest:** deterministic, Go-generated, sorted by participant ID. Per participant: `filed` (bool),
  artifact `path`, `owner` (frontmatter `agent:`), `validity` (driver validator verdict), `stance`
  (`✅`/`🟡`/`❌`/absent, parsed by the same rules `parley consensus status` uses; unparseable →
  `unparsed`, fail closed). Plus: condition, met/timeout, blocks list, and a fixed-enumeration
  next-action line (e.g. `awaiting: kimi-1 round-01`) — never free prose. Two consecutive runs on an
  unchanged tree must produce a byte-identical digest body so the organizer can diff phases cheaply.
- **Digest is an index, not evidence:** the organizer opens full artifacts only to adjudicate — any
  `❌`, `unparsed`, `DISPUTED`, or surprising validity verdict triggers reading the raw file. Raw
  artifacts stay canonical.

### C. Organizer brief + slim SKILL.md

- **Role-scoped packet:** add `--role facilitator|participant` to `parley protocol packet`. The
  renderer's `Request` struct currently has no role input [PRIMARY:
  `internal/protocolpacket/packet.go:53-62` @CLI]. `--role facilitator` emits the protocol's own
  facilitator reading set — Quickstart, §4, §5, §9, and the active §11 transport — as whole verbatim
  sections plus an omission index (which sections are absent and when to read them). The attestation
  gains a `role` field; `context_mode` remains `full` for the default and for participants, so the
  ratified "full context is the default" invariant is untouched for everyone but the facilitator.
  The runner's per-launch protocol prepend (full body verbatim [PRIMARY:
  `internal/runner/protocol_context.go:111` @CLI]) is participant plumbing and is **not** changed.
- **The organizer brief is the facilitator packet** — one generated artifact, no parallel static brief
  to drift. Post-compaction re-orientation becomes: facilitator packet (~43–50 KB of the 108.4 KB
  deck view; deck COOPERATION.md is 108,400 B [PRIMARY: `wc -c` @CLI]) + `parley status --idea <slug>`
  + `parley wait`. The skill documents this three-step loop; no full SKILL.md or full COOPERATION.md
  re-read.
- **Slim SKILL.md:** split the current 62,411 B file [PRIMARY: `wc -c` @SK] into a ≤20 KB every-session
  core and on-demand `references/` (the directory already exists [PRIMARY:
  `skills/parley-deck/references/` @SK]). Core keeps verbatim: frontmatter (trigger reliability — the
  description must not change), Core Rule, Non-Solo Requirement, Required Protocol Context, Automation
  Mode, Startup Flow, kickoff + round/consensus templates, commit conventions, and a new facilitator
  driver-routing section. Moved to references/: hand-launch templates, headless agent configuration,
  timeout policy, transport-selection detail, roster deep-dive, drift-check internals, coverage
  checklist detail. A pointer index in the core links each moved section.

### D. Session per phase + organizer measurement

- **Phase bootstrap for fresh organizer sessions:** `parley status --idea <slug> --organizer` emits a
  bounded (≤4 KB) bootstrap: current phase, per-participant owed/artifact state, open escalations,
  and the exact resume command (`parley continue …`). Built on the existing status machinery and the
  handoff-packet concepts already in `internal/runner/handoff.go` [PRIMARY: `WriteHandoffPacket`,
  `handoff.go:38` @CLI]. Each phase can then start in a fresh organizer session with: facilitator
  packet + bootstrap + one blocking wait — no compaction carry-over, which attacks the 42%-of-uncached
  compaction cost and the ~33–39% re-orientation share identified by the study (SECONDARY:
  `source-context/organizer-token-study.md` §2, figures not recomputed by me).
- **Organizer usage ledger tooling:** `parley usage record --idea <slug> [--source codex|claude]
  [--session <path>]`. It parses the *client's own* accounting — Codex rollout JSONL `token_count`
  events, Claude session JSONL `message.usage` — attributes records to the idea by slug mentions in
  the session file (phase boundaries from `.parley-runtime` run records when available; unresolved
  splits labeled ambiguous), and appends a row to `parley-deck/ideas/<slug>/organizer-usage.md` in
  the format this idea already uses [PRIMARY: file exists with Phase 0→1 row @CLI]. Tokens only, no
  dollar claim. This is deliberately *client accounting*, addressing the earlier rejection of
  self-reported usage records as cost evidence: the rejected thing was a model asserting its own
  usage; this is measured client-runtime data with a named source path per row.
- **Participant telemetry fixes:** the collector switch has cases for `claude`, `codex`, `zcode` and
  no `kimi` case [PRIMARY: `internal/telemetry/usage.go:229-298` @CLI], and `structured` is only true
  when the argv already carries `--json`/`--output-format json` [PRIMARY:
  `internal/telemetry/usage.go:380-395` and `internal/runner/telemetry.go:147` @CLI], while kimi's
  headless mode is plain `kimi -p` [PRIMARY: `internal/agents/discover.go:331-358` @CLI]. Fixes: add
  a `kimi` case driven by a captured structured-output fixture; add adapter-appropriate structured
  flags to default headless argv only where the adapter supports them without breaking `-p` semantics.
  If kimi emits no usage in structured mode, `parley usage` records that gap explicitly — honest
  failure, no fabricated numbers.

### Protocol text deltas (all three COOPERATION.md copies + staged core)

COOPERATION.md exists in three copies: deck view and embedded default in the CLI worktree, plus the
skill's bundled reference [PRIMARY: `parley-deck/COOPERATION.md` 108,400 B and
`internal/protocol/defaults/COOPERATION.md` 108,154 B @CLI; `skills/parley-deck/references/COOPERATION.md`
108,244 B @SK]. Deltas, kept minimal and permissive:

1. §9.0 item 1: an official facilitator launch may consume the attested role-scoped facilitator view;
   participant launches keep the full default.
2. §4 Phase 5 + Quickstart Facilitator row: implementation and code verification default to
   participants when available; the facilitator reads verdicts and validator output.
3. §11 facilitator practice (advisory): prefer one blocking `parley wait` + digest over repeated
   short polls.

The two CLI copies are guarded by a Go drift test; the skill copy is unguarded — I propose a skill-repo
test that checks `references/COOPERATION.md` against a sha recorded in `references/compatibility.json`
at release time, closing that gap cheaply. The new global core version is staged from the previous
core template plus these hunks; `parley protocol publish` stays the owner's attended action.

### Implementation sketch (sequencing for Phase 5)

1. Protocol deltas (three copies) + drift tests + staged core rendering.
2. `--role facilitator` on the packet renderer + tests (verbatim-section, omission-index, attestation).
3. SKILL.md split + skill test suite + packaging/link checks.
4. `parley wait` + digest + golden determinism/timeout tests.
5. Driver participant-implementer resolution + A-related tests.
6. `--organizer` bootstrap, `parley usage record` (Codex + Claude fixtures), kimi telemetry fixture.
7. CHANGELOGs, version bumps (CLI 1.49.0, skill 2.13.0), release prep — release execution remains the
   organizer's owner-gated step; participants verify channels afterward per the brief.

### Observable acceptance criteria (draft for FINAL)

- **AC-B1:** on a fixture deck with one missing round-01 artifact, `parley wait --for round-01
  --timeout 2s --json` exits 2 and the digest names the missing participant; after the artifact is
  written, the same command exits 0 with per-participant filed/owner/validity/stance.
- **AC-B2:** two consecutive `--json` runs on an unchanged tree produce a byte-identical digest body.
- **AC-C1:** `parley protocol packet --phase 1 --track deliberation --role facilitator --json` emits
  an attested body containing Quickstart, §4, §5, §9 and the active §11 transport verbatim, ≤60 KB,
  with an omission index; omitting `--role` yields the unchanged full default.
- **AC-C2:** slim SKILL.md ≤20 KB; every moved section exists under `references/` and is linked from
  the core; skill tests green; `parley-deck-skill install --target all --force` verifies by content hash.
- **AC-A1:** a driver fixture run with `auto_implement` records a participant agent as implementer in
  `IMPLEMENTATION.md`, and the driver event log shows no facilitator code-verification launch.
- **AC-D1:** `parley usage record` against checked-in Codex and Claude fixtures reproduces expected
  token totals (golden test) and appends exactly one ledger row per run.
- **AC-D2:** a kimi structured fixture yields recorded usage tokens, or the recorded-gap path is
  exercised and visible in output.
- **AC-D3:** `parley status --idea <slug> --organizer` output ≤4 KB and contains phase, owed
  artifacts, resume command, open escalations.
- **AC-P1:** all three COOPERATION.md copies carry deltas 1–3; the CLI drift test is green; a staged
  core exists under `~/.parley/staging/` with project zones preserved; the escalation inbox note
  contains the exact attended publish command.

## Existing alternatives

- **Blocking wait:** the toolchain ships no wait command [PRIMARY: `internal/app/app.go:60-108` @CLI].
  The shipped alternative is a foreground `parley run` (blocks, but runs a whole driver cycle, not a
  passive wait) or the brief's manual `until`-loop fallback — the latter spends organizer tool calls.
  Rejected because both either do too much (run) or cost polls (loop); adopt the new `wait` for the
  narrow wait-then-digest need.
- **State digest:** `parley status` and `parley consensus status` already exist [PRIMARY:
  `internal/app/app.go:68,74,536` @CLI]; measured at 10.7 KB deck-wide / 2.2 KB per-idea / 1.8 KB
  consensus (SECONDARY: study §3, sizes not recomputed). Adopt: reuse their parsing for the digest
  rather than a new parser — `wait` wraps status-shaped data in a blocking condition, it does not
  duplicate it. Reject extending `status` with blocking instead of a new command: `status` is
  read-only fast everywhere and is scripted by other tools; a blocking default would be a behavior
  change for existing callers.
- **Role-scoped context:** the shipped `--optimize` packet exists but is experimental, participant-
  oriented, and the `Request` struct has no role field [PRIMARY:
  `internal/protocolpacket/packet.go:53-62` @CLI; mode defaults documented at
  `internal/app/protocol_packet.go:23,55` @CLI]. Reject reusing `--optimize` for the organizer
  (omission-based, unratified as a default); adopt a separate `--role` whole-section view, attested,
  facilitator-only.
- **Organizer brief:** a hand-written static brief file would duplicate protocol content and drift;
  the protocol already names the facilitator reading set [PRIMARY: the Quickstart "Who are you?"
  table, Facilitator row at `parley-deck/COOPERATION.md:34` @CLI — note the study cites `:26-33`,
  which is off by one; the table spans `:29-34`]. Adopt: generate the brief from the renderer so the
  reading set has one source.
- **Session handoff:** `WriteHandoffPacket` already models per-launch handoff [PRIMARY:
  `internal/runner/handoff.go:38` @CLI]. Adopt/extend for the per-phase organizer bootstrap rather
  than a new state file format.
- **Usage measurement:** telemetry today covers only runner-launched participant processes
  [PRIMARY: `internal/telemetry/usage.go` collector wired at `internal/runner/telemetry.go:148` @CLI];
  nothing ingests organizer client sessions. No shipped alternative; the manual ledger in this idea is
  the practice the new command automates.

## Concerns / open questions

1. **kimi's structured usage shape is unknown to me.** Whether `kimi -p` can emit structured usage at
   all must be established by a live capture during implementation; AC-D2 is written to pass honestly
   either way (recorded gap is a valid outcome). This is `UNVERIFIED` — I make no claim that the shape
   exists.
2. **Role view wording vs the ratified packet trial.** Adding a role-scoped view beside `full`/`packet`
   needs precise §9.0 wording so nobody reads it as licensing omission-based context for participants.
   I propose the attestation carry an explicit `role` field and that `context_mode` values remain
   unchanged for non-facilitator launches.
3. **Usage attribution across long sessions.** A Claude session spanning months and many ideas cannot
   be split exactly; I propose whole-session attribution when only one idea matches, run-record
   timestamp windows otherwise, and an explicit `ambiguous` label rather than a guess. The ledger's
   headers must state the method.
4. **Digest next-action line** must stay a fixed enumeration of state, not generated prose, or it
   drifts toward advisory summarization.
5. **Skill split and trigger reliability:** the frontmatter description and Core Rule must remain
   verbatim; if reviewers know of other load-bearing anchors (external docs linking into SKILL.md
   sections), the moved sections need stable redirect notes.
6. **`parley wait` and escalations:** v1 watches protocol phase conditions only; an unanswered
   `to-user` escalation surfaces via the bootstrap/status, not the wait condition. Worth confirming
   with other participants whether a wait should also return early on a new blocking escalation.
7. **Stance parsing reuse:** I assert the consensus-status parsing is reusable for the digest, but I
   have not read that parser [PRIMARY only for its existence]. Implementation must verify or extract
   shared parsing rather than fork it.

## Risks

- **False "valid" in a digest** would let an organizer skip needed adjudication. Mitigation: fail-closed
  `unparsed`, golden determinism tests, and the rule that any ❌/DISPUTED/unparsed forces reading the
  raw artifact.
- **Scope creep into participant context** would violate the ratified full-default invariant. Mitigation:
  the role flag is facilitator-only; C1 acceptance explicitly checks the unchanged default.
- **Three-copy drift** — two copies are test-guarded, the skill copy is not; the compatibility.json
  hash check is cheap insurance but is new surface (kept optional in my proposal).
- **Behavioral risk:** token savings assume organizers actually adopt the lean loop; the skill routing
  and facilitator packet make the cheap path the default path, and the D ledger makes regression
  measurable per phase.
- **Windows/portability** for `wait` polling and path handling; the release ships Windows assets, so
  the wait/usage commands need a Windows CI leg.
- **Attribution privacy/size:** scanning client session files must be bounded (date-filtered, size-
  capped) and never copy credentials into the ledger; the ledger row records paths and counts only.
- **The 25-minute default timeout** encodes an assumption about provider cache lifetime; it stays a
  flag so a changed provider reality needs no code change.

## Provenance (§15.2)

Claims above are tagged inline. Summary of bases:

- **PRIMARY (verified by me at HEAD):** CLI worktree `4ce8fa7` — `internal/app/app.go:60-108`
  (subcommand dispatch; no `wait`); `internal/app/driver_consensus.go:82-95` (participant drafter);
  `internal/protocolpacket/packet.go:53-62` (Request has no role); `internal/app/protocol_packet.go:23,55`
  (full default, `--optimize` experimental); `internal/runner/protocol_context.go:111` (full body
  prepended verbatim to launches); `internal/runner/telemetry.go:147` + `internal/telemetry/usage.go:229-298,380-395`
  (adapter coverage, structured-args heuristic; no kimi case); `internal/agents/discover.go:331-358`
  (kimi headless `-p`, plain argv); `internal/runner/handoff.go:38` (handoff packet);
  `parley-deck/COOPERATION.md` line counts `29-34` (role table; Facilitator row `:34`), `:407-409`,
  `:439`, `:443-446`, `:890`, `:1342` (role rules); file sizes 108,400 B / 108,154 B; idea files
  (`organizer-usage.md` with Phase 0→1 row). Skill worktree `d1e57d5` — SKILL.md 62,411 B, no driver-
  command mentions (grep), `references/` layout, references/COOPERATION.md 108,244 B. Packet SHA-256
  recomputed and matched.
- **SECONDARY (named dependency, not recomputed by me):** all quantitative figures from
  `source-context/organizer-token-study.md` (cost shares, re-orientation share, compaction share,
  polling share, status output sizes, telemetry coverage counts, Claude session figures). The study
  states three independent verifiers recomputed headline numbers; I rely on it as a named source, and
  its per-agent data at `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/organizer-token-study/2026-09-23/`
  (directory and `README.md`/`results/` existence verified by me).
- **Study locator correction (PRIMARY):** the facilitator reading-set row is at
  `parley-deck/COOPERATION.md:34`, not `:26-33` as the study cites; the table spans `:29-34`.
- **UNVERIFIED:** existence of a kimi structured usage shape (Concern 1); reusability of the
  consensus-status stance parser (Concern 7). Neither is load-bearing for the design's shape; both
  are implementation-time checks with honest-failure paths.
- **Novelty scope:** I claim only "no `wait` subcommand and no role input exist at the HEADs I
  inspected" — I do not claim nothing similar exists elsewhere.
