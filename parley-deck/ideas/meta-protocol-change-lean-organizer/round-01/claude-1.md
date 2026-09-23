---
agent: claude-1
idea: meta-protocol-change-lean-organizer
round: 1
date: 2026-09-23
---

## Protocol context attestation

```json
{"context_mode": "full", "source_sha256": "12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18", "packet_sha256": "12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18", "fallback_reason": null, "body_path": "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer/.parley-runtime/protocol-packets/full-phase1-deliberation-12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18.md"}
```

Read in full. `source_sha256 == packet_sha256` here because `context_mode: full` emits the live
authority verbatim; I confirmed the emitted body is byte-identical to `parley-deck/COOPERATION.md`
(`shasum -a 256` of both → `12e4b31cd3f6c106…`, 1,380 lines, 108,400 B).

## Summary

The organizer's bill is `requests × resident context`, and the study locates almost all of it in two
multipliers — the post-compaction re-orientation floor (~250–265 KB re-read after each of 50
compactions) and the request count inflated by hand-launching, polling and self-implementation. My
central finding from verifying the toolchain at HEAD is that **A–D is mostly a routing and surfacing
problem, not a building problem**: a deterministic LLM-free round digest, the three artifact
validators, the phase/track/transport-scoped packet renderer, and a blocking foreground driver all
already exist — they are emitted only into a TUI the headless organizer cannot read, or gated behind
a flag nobody turns on, while `SKILL.md` at v2.12.1 never names `parley run`, `continue`, `status`,
`consensus` or `preflight` even once and therefore actively teaches the hand-launching the study
measured. I propose implementing A–D as four thin, checkable mechanisms over that existing
machinery: an optional `facilitator:` declaration the driver enforces (A), `parley wait` returning
an extended `PhaseDigest` on stdout (B), an `--audience` dimension in the existing applicability map
plus a SKILL.md split and a *generated* organizer brief (C), and a streaming `parley usage ingest`
that may read only client-written accounting files and never a number typed by a model (D). I also
report, plainly, that the single largest measured lever — the 22 KB per-request skills catalogue and
the standing "read in full" instructions — lies outside this release and A–D cannot reach it.

## Proposed approach

### 0. The cost identity, and which term each scope item moves

Organizer cost ≈ `Σ_requests (resident context)`, price-weighted by cached vs uncached.

| Term | What inflates it (study §2) | Scope item that moves it |
|---|---|---|
| request count | polling (9.1–12.2%), 70 hand-written prompts, self-verification | B (one blocking call), A (delegation), C (driver routing) |
| resident-context floor | ~250–265 KB re-read in the first 10–15 calls after each compaction | C (brief ≪ SKILL.md + COOPERATION.md) |
| number of compactions | 50 compactions; 42% of all *uncached* input | D (fresh session per phase) |
| price weight | reasoning at xhigh; output ~0.67% of tokens but ~19–27% of cost | out of scope (non-goal: no model/effort downgrade) |

Two consequences I want on the record before the design:

1. The study's own simulation says cutting reads **outside** the re-orientation window saves ~0%
   because the session is compaction-saturated. So "read less in general" is not a design. Only the
   post-compaction floor and the compaction count are levers. C and D are therefore the load-bearing
   items; A and B are what make a fresh short session *sufficient* to run a phase.
2. Instruction-only fixes are disproven for this exact problem: after the owner's 2026-09-15
   "just manage the others", Go edits fell 452 → 2 but the hourly input rate *rose* (12.9M/h →
   15.0M/h). **Every item below must therefore land as a mechanism the tooling enforces or emits,
   not as a sentence in a document.** Where I propose protocol text, it exists to license a
   mechanism, never to substitute for one.

### A — The organizer neither implements nor verifies code

**What the protocol already permits.** Nothing forces the facilitator into these roles. The drafter
is the `author:`, or with `author: user` the first round-01 writer or a volunteer
(`parley-deck/COOPERATION.md:407,409`); the implementer is the FINAL drafter unless claimed
(`:443`); the facilitator's procedural calls are provisional (`:1340-1343`, §15.5); and the
facilitator's *verification* duty is existence, ownership and validity — not code correctness
(`:437`, `:890`). The driver already selects the drafter from `participants:` only
(`internal/app/driver_consensus.go:82-100`, `firstHeadlessAgent`).

So the hole is narrow and specific: **when the facilitator is also listed in `participants:`**, every
default quietly routes back to it. That is precisely the shape the measured session had.

**Mechanism.**

1. `00-prompt.md` gains an **optional** `facilitator: <agent-id>` field. Absent → today's behaviour,
   byte for byte. This matters for the stop rule: it adds no obligation to any existing or future
   deck that does not set it.
2. `parley preflight` fails closed when `facilitator:` names an agent that also appears in
   `participants:`, unless `facilitator_participates: true` is also set. The failure names both
   fields and exits non-zero; it is a readiness finding, not a silent downgrade.
3. The driver refuses to select the declared facilitator as drafter, implementer, reviewer or
   goal-done checker. This is one predicate reusing the existing participant filter.
4. Protocol text: one sentence in §4 Phase 5 and one in §4 Phase 6 recording that a declared
   non-participant facilitator is ineligible for implementer and reviewer, and one line in §9.0
   recording the preflight check. No new duty is created for anyone; an existing default is named.

**Why this is enough for "does not verify code itself".** Verification does not vanish — it moves to
participants, and the organizer's substitute for reading code becomes the digest's validator column
(B). Without B, A is an instruction, and instructions are disproven here.

**Observable criteria (A)**

- `parley preflight` on a deck whose `00-prompt.md` sets `facilitator: X` and lists `X` in
  `participants:` exits non-zero and prints both field names; adding `facilitator_participates: true`
  makes it exit 0.
- With `facilitator: X` set and `X` absent from `participants:`, a full auto-drive run never launches
  `X` for a drafter, implementer, reviewer or goal-done role — assertable from the run event log
  (`agent.started` payloads) in a driver test with a stub adapter.
- A deck with no `facilitator:` field produces a byte-identical run plan to v1.48.0 (regression test).

### B — `parley wait` + a deterministic phase digest

**What already exists, and why it does not help the organizer today.**

`internal/driver/digest.go:10-48` already builds exactly the artifact scope B asks for: a
"deterministic, LLM-free position map of a completed round" with a per-agent presence flag, a
capped verbatim `## Summary` slice, a degraded-extraction flag, and keyword stance counters that the
file's own comment is careful to call "FLAGS (hints), never verdicts". It is built at
`internal/driver/driver.go:488` and appended as a `round.digest` event — and the only consumer in
the tree is `internal/tui/roundsummary.go:9-38`, the **Home tab renderer**. A headless organizer
never sees it. `parley run` also opens the TUI by default (`--no-tui` to disable).

So B is not "build a digest". B is: **extend the digest past design rounds, and give a headless
caller a blocking command that returns it on stdout.**

**Mechanism.**

1. Generalize `RoundDigest` → `PhaseDigest`, covering design rounds (today), review rounds, consensus
   signoff state, and implementation status. Reuse the shipped validators for the validity column —
   `ValidateRoundOneArtifact` (`internal/protocol/roundartifact.go:23`), `ValidateReviewArtifact`
   (`reviewartifact.go:17`), `ValidateFinal` (`finalsections.go:92`). Nothing new is invented; the
   digest becomes the place their output is *shown*.
2. Per-agent columns, all mechanically derived: `agent`, `path`, `filed`, `bytes`, `owner` (the
   `agent:` frontmatter value, so a mis-owned file is visible), `valid` (validator verdict verbatim),
   `stance_flags` (existing keyword counters), `fell_back`.
3. New command:
   `parley wait --idea <slug> [--for round|consensus|review|implementation|any] [--timeout D] [--json]`.
   It blocks on the append-only event log (`internal/store/events.go:42,68`) and returns on the
   first of: the awaited boundary reached, the timeout, or a `driver.error` / escalation event.
4. Exit codes carry the outcome so a wrapper never has to parse prose: `0` boundary reached,
   `3` timeout (not a failure — re-invoke), `4` escalation/blocked, `1` usage/IO error. A timeout
   prints the partial digest, so even a timed-out call is informative rather than wasted.
5. `--timeout` defaults from `~/.parley [defaults.loop]`, seeded conservatively. **I deliberately do
   not hard-code the brief's "below the provider's 30-minute prompt-cache lifetime" figure**: I could
   not locate an authoritative source for that TTL from this worktree, it is provider- and
   plan-specific, and a number baked into Go becomes wrong silently. Making it configuration keeps
   the intent and drops the unverifiable constant. See §15 provenance.

**Why this is a real saving and not a shell trick.** The study counts 1,488 of 1,539 `write_stdin`
polls at a 1 s yield, with ~9M input tokens attributable to waiting on participants alone. Each poll
is a *request* carrying the full resident context. One blocking call replaces the whole sequence with
one request. The saving is at the cached-read rate rather than the uncached rate — cached reads are
~54% of cost, so it is smaller than the token share suggests, and I would rather state that here than
let a later measurement look like a shortfall.

**Observable criteria (B)**

- `parley wait --idea X --for round --timeout 1s --json` against a deck with a complete round-01
  exits `0` and emits one object per participant, each carrying `path`, `filed`, `owner`, `valid`,
  `stance_flags`.
- Determinism: two consecutive invocations over an unchanged tree emit byte-identical JSON.
- No-LLM property, asserted by test: every string field in the emitted JSON is a path, a fixed enum
  value, a verbatim ≤`digestPositionCap` slice of the artifact, or a number. No field is model-written.
- Against a deck with round-01 *incomplete* and no driver running, the command exits `3` at the
  timeout, prints the partial digest, and prints which agents are outstanding.
- Round-trip: driving a full round with `--no-tui` plus one `parley wait` makes exactly one organizer
  tool call between launch and digest (assertable from the run's own event log).
- Raw artifacts stay canonical: the digest carries the path of every file it summarizes, and a
  regression test asserts a digest never replaces or rewrites a round file.

### C — Audience-scoped protocol view, slim SKILL.md, generated organizer brief

This is three separable pieces and I want them judged separately, because they are not equally
valuable and the brief has already framed one of them more optimistically than the bytes support.

**C1 — `--audience` in the existing packet renderer.**

`parley protocol packet` already scopes by `--phase`, `--track`, `--transport`, `--idea` and `--flag`
(`internal/app/protocol_packet.go:48-60`). It has no audience input. It *does* have a `Role` field —
`internal/protocolpacket/packet.go:116`, `source.go:41,64-76` — but that is the deck's
`protocolRole` (`source|consumer`), a completely different concept. **A new flag must therefore not
be called `--role`**; I propose `--audience participant|facilitator`.

The classification lives in `parley-deck/meta/packet-applicability.yaml`, whose header states it *is
protocol* and that changing a classification is a §7 change — which is exactly what this idea is. The
mechanism is to add an `audiences:` key to the existing `when:` grammar and let `packet check` prove
it, with one hard safety property: **the audience dimension may only omit blocks already omittable at
that phase/track/transport; it can never cut below the ratified never-cut floor**
(`internal/protocolpacket/applicability.go:193-231`). That property is a one-line check, and it is
what keeps an audience view from quietly becoming a different protocol.

The protocol already names the facilitator's reading set — `parley-deck/COOPERATION.md:34`
("Facilitator | Quickstart, §4, §5, §9, your §11 transport"). C1 makes the renderer able to express
a rule the document has carried all along.

**What the bytes actually say, measured at HEAD (108,400 B total):**

| Region | Bytes | Note |
|---|---:|---|
| Quickstart (12–44) | 2,414 | facilitator set |
| §4 (203–732) | 35,991 | facilitator set |
| §5 (733–745) | 1,334 | facilitator set |
| §9 (843–892) | 5,624 | facilitator set |
| §11.B (934–1037) | 9,291 | facilitator set, github-pr |
| **literal facilitator reading set** | **54,654** | **50.4% of full** |
| Phases 5–8 (441–699) | 16,287 | omittable at phase 1 |
| §11.C (1038–1100) | 4,195 | omittable by transport today |
| Appendix A (1101–1126) | 2,021 | bootstrap only |
| §12 + §13 (1127–1205) | 12,716 | opt-in, not in never-cut |
| **omittable at phase 1** | **35,219** | leaves **73,181 B = 67.5%** |

**The honest reading:** 67.5% sits inside the band the study already measured for the *existing*
optimizer without flags (54.8–71.7%). Most of that saving is the **phase** dimension, which is
already built and simply not enabled. The audience dimension is a real but **modest increment** on
top of it. I would rather say so in round 1 than have it discovered in review. The brief's "a pure
organizer needs ~43–50 KB of the 108.4 KB" is reachable only by omitting inside §4 beyond what the
phase scoping already does, and I do not think that is worth the risk against the never-cut floor.

**C2 — split `SKILL.md`. This is the larger lever of the three.**

`skills/parley-deck/SKILL.md` at tag `v2.12.1` is 990 lines / **62,411 B** and, verified by
`grep -oE 'parley [a-z-]+'`, names only `parley roster` (9), `parley protocol` (4), `parley init` (4),
`parley agents` (2), `parley layer` (1). It **never once** names `run`, `continue`, `consensus`,
`status` or `preflight`. Meanwhile it devotes whole sections to `## Headless Agent Configuration`,
`## Generic CLI Invocation Contract` and `## Timeout Policy`. The skill is not neutral about the
measured behaviour — it *teaches* hand-launching. That single fact explains the study's "ran
`parley run` once, `parley status` never, wrote 70 prompts by hand" better than any discipline
failure does.

Proposed split — relocation only, no rule text deleted:

- **Slim core** (target ≤ 20,000 B): Core Rule, Non-Solo, Required Protocol Context, Startup Flow, a
  new **Driver-first** section (`run` / `continue` / `wait` / `status` / `consensus` / `preflight`,
  with hand-launching named as the fallback that must be recorded), File Ownership, Escalation,
  Quality Gates, Keep It Small.
- **`references/HEADLESS_LAUNCH.md`**: Headless Agent Configuration, Generic CLI Invocation Contract,
  Timeout Policy, Agent Capability Discovery.
- **`references/ARTIFACT_TEMPLATES.md`**: every fenced round / review / consensus / implementation
  template body.
- **`references/ROSTER_AND_PROTOCOL.md`**: roster verbs, global core, drift check, coverage
  checklist, transport/deck bootstrap.

**C3 — `parley organizer brief --idea <slug>`: generated, never stored.**

The post-compaction read becomes one command whose output is computed from live state: the packet
attestation + facilitator body path, `parley status --idea --json` (**measured 1,459 B** for this
idea; whole-deck `parley status` is 10,897 B), the current `PhaseDigest`, the driver's next action,
and the phase pointer from D. Target ≤ 8 KB, against the ~250–265 KB the study measured.

**This must be a view, not a file.** The study's §6 lists "stored summary sidecars" among the
already-rejected patterns, and I agree with that rejection: a stored brief goes stale, becomes a
second source of truth, and creates an obligation to maintain it. Computing it on demand from
`00-prompt.md`, the artifact tree and the event log keeps raw artifacts canonical and adds no duty to
any participant. If the design later drifts toward persisting it, that is the moment to stop.

**Observable criteria (C)**

- `parley protocol packet check` passes with the `audiences:` key present, and fails when any
  audience rule would omit a never-cut block at that phase/track/transport (negative test).
- `parley protocol packet --audience facilitator --phase 1 --json` emits a valid attestation; an
  unrecognized `--audience` value falls back to full context with a stated reason, matching today's
  unknown-value behaviour.
- Byte target: the phase-1 `deliberation` / `github-pr` facilitator body is ≤ 75,000 B of the
  108,400 B source, and its omission index lists every omitted block with a trigger.
- `wc -c skills/parley-deck/SKILL.md` ≤ 20,000; the slim core names every `references/*.md` file it
  relocated content into; a script proves every relocated heading appears in exactly one reference
  file and none was dropped.
- The slim core names `parley run`, `continue`, `wait`, `status`, `consensus` and `preflight` at
  least once each (the direct inverse of the verified v2.12.1 state).
- `parley organizer brief --idea <slug>` emits ≤ 8,192 B, is byte-identical on two consecutive runs
  over an unchanged tree, and writes no file (assertable: the command runs with the deck mounted
  read-only).

### D — Session per phase + organizer usage ledger

**D1 — session per phase.** Once C3 exists, "start a fresh session per phase" needs no new artifact:
the brief *is* the handoff, recomputed. I specifically propose **not** to introduce a hand-written
state/handoff file — it would be a stored summary sidecar (rejected), and writing it would be a new
mandatory obligation the stop rule forbids. The protocol change is a single permissive line in §9
recording that a facilitator may re-orient from the generated brief instead of re-reading SKILL.md
and the full COOPERATION.md. Permissive, not mandatory.

**D2 — organizer usage ledger.** I verified the client accounting is parseable exactly as the brief
assumes. In the measured rollout (228 MB, 3,851 `token_count` records), each record is:

```json
{"timestamp":"2026-09-09T22:01:18.014Z","type":"event_msg",
 "payload":{"type":"token_count","info":{
   "total_token_usage":{"input_tokens":17218,"cached_input_tokens":15488,
     "cache_write_input_tokens":0,"output_tokens":7,"reasoning_output_tokens":0,
     "total_tokens":17225},
   "last_token_usage":{…},"model_context_window":258400}}}
```

Two design facts fall straight out of that. First, `model_context_window: 258400` independently
matches the study's "window 258.4k", which is a small corroboration of its instrumentation. Second,
**the source file is 228 MB**: the ingest must stream, and the organizer must never read it. That is
the criterion that makes this feature safe rather than self-defeating.

Mechanism: `parley usage ingest --agent <id> --source codex-rollout|claude-jsonl --path <file>
--idea <slug> --phase <n>` streams the file, takes the last `total_token_usage`, and appends one row
to `parley-deck/ideas/<slug>/organizer-usage.md` plus a machine-readable sibling.

**On the prior rejection of "self-reported usage records as cost evidence."** I think the distinction
holds, and it needs to be *mechanical* rather than argued: the model does not author
`payload.info.total_token_usage` — the CLI writes it from the provider's response, in a file the
model does not control. The check that keeps the distinction real is:

> `parley usage ingest` MUST NOT accept a token count as an argument. It accepts a **path** and
> parses the file itself; the recorded row carries the source path and the parser id.

A number typed by an agent is `RECALL` and inadmissible; a number the tool parsed from the client's
own file is `PRIMARY` with a locator. That is the same provenance rule §15.2 already applies to
everything else, and it is testable.

**D3 — participant telemetry fix.** `internal/telemetry/usage.go` handles `claude` (`:230`), `codex`
(`:252`) and `opencode` (`:261`), plus a zcode envelope (`:190`, `:286`). **There is no `kimi` case**
— the study's "kimi-1 0/20 invocations with usage" is a parser gap confirmed at HEAD, not flaky
adapters. Fix: add the kimi case and default the runner argv to the structured-output form where the
adapter supports it.

**Observable criteria (D)**

- `parley usage ingest` rejects any attempt to pass a count directly (no such flag exists; a test
  asserts the flag set).
- Ingesting the 228 MB rollout emits ≤ 1 KB to stdout, appends exactly one ledger row, and completes
  with bounded memory (streaming decoder, asserted by a large-fixture test).
- The ledger row carries idea, phase, agent, source path, parser id, and the six `total_token_usage`
  fields verbatim.
- A kimi-1 runner launch with structured output records a non-nil `usage` in its telemetry record —
  the direct inverse of the measured 0/20.
- Re-ingesting the same file for the same idea+phase is idempotent (no duplicate row).

### What I deliberately do not propose

Named so the drafter does not have to guess, and so the stop rule is easy to check against my
position:

- No LLM or lossy summarization anywhere. Every digest, brief and ledger row is Go-generated from
  files or event logs, and carries the path of what it summarizes.
- No stored summary sidecar, and no hand-written handoff file.
- No tool that selects which protocol text an agent sees at retrieval time. `--audience` is a
  ratified, human-reviewed classification in a protocol file checked by `packet check`, not a
  retrieval heuristic.
- No new mandatory obligation. `facilitator:` is optional; the §9 brief line is permissive;
  `parley wait` is an alternative to polling, not a required step.
- No model or effort downgrade, no roster change, no environment-catalogue pruning (the owner's §7
  levers are outside this release, and I say below why that matters).

### Consolidated acceptance shape

Every criterion above is either a command with an exit code, a byte bound on a named file, a
determinism assertion, or a negative test. None requires a human to judge whether the organizer
"felt leaner". I would push back on any criterion that does.

## Existing alternatives

§15.6(a). For each mechanism this proposal would otherwise build by hand, what the toolchain already
ships, with a locator. All locators verified at HEAD `4ce8fa7` in this worktree; the skill worktree is
at tag `v2.12.1`, HEAD `d1e57d5`.

| Mechanism I would build | Already shipped | Locator | Why it does not currently serve the organizer |
|---|---|---|---|
| Deterministic LLM-free round digest | `BuildRoundDigest` / `RoundDigest` | `internal/driver/digest.go:10-48`; built at `internal/driver/driver.go:488` | Emitted only as a `round.digest` event, rendered only by the TUI (`internal/tui/roundsummary.go:9-38`); design rounds only; no stdout path |
| Per-artifact validity column | `ValidateRoundOneArtifact`, `ValidateReviewArtifact`, `ValidateFinal` | `internal/protocol/roundartifact.go:23`, `reviewartifact.go:17`, `finalsections.go:92` | Run inside the driver; verdicts are not collected into anything an organizer reads |
| Blocking wait | foreground `parley run` | `internal/app/app.go:80`; `parley run --help` | Blocks, but opens a TUI by default and prints progress lines, not a digest; a headless organizer gets no structured boundary signal |
| Event stream to wait on | append-only JSONL store | `internal/store/events.go:42,68`; types incl. `run.phase`, `round.completed`, `driver.error` | Exists and is sufficient; nothing exposes a blocking read of it |
| Compact state digest | `parley status --idea --json` | `internal/app/app.go:705-706`; measured 1,459 B | Process state only — no artifact paths, ownership, validity or stance |
| Signoff state digest | `parley consensus status --json` | `internal/app/app.go:68` | Exists; not composed into a single phase view |
| Role-scoped protocol view | phase/track/transport/flag-scoped packet + applicability map + never-cut list | `internal/app/protocol_packet.go:48-60`; `parley-deck/meta/packet-applicability.yaml`; `internal/protocolpacket/applicability.go:193-231` | No audience dimension; `Role` there is the deck's `protocolRole` (`source|consumer`), `internal/protocolpacket/packet.go:116`, `source.go:41,64-76` — a name collision to avoid |
| Facilitator reading set | already protocol text | `parley-deck/COOPERATION.md:34` | Stated for humans; no renderer input expresses it |
| Non-facilitator drafter/implementer | drafter = author / first round-1 writer / volunteer; implementer = drafter unless claimed; drafter must be a participant | `parley-deck/COOPERATION.md:407,409,443`; `internal/app/driver_consensus.go:82-100` | Already true when the facilitator is not a participant; nothing enforces that precondition |
| Participant usage capture | `agent.usage` events + adapter parsers | `internal/telemetry/usage.go:230,252,261` (+ zcode envelope `:190,:286`) | No kimi case; nothing captures the organizer at all |
| Client-side organizer accounting | Codex rollout `token_count` records | `~/.codex/sessions/2026/09/10/rollout-2026-09-10T00-01-07-01a08830-*.jsonl` (3,851 records; 228 MB) | Exists and is parseable; no command ingests it, and its size forbids reading it into context |
| Protocol copy drift guard | normalized byte-compare of two copies | `internal/protocol/drift_test.go:12` (`liveDeckPath`) | Guards deck ↔ `internal/protocol/defaults/` only; the skill's `references/COOPERATION.md` is unguarded |
| Session resume state | `internal/sessionstore` | `internal/sessionstore/sessionstore.go:18-30` | Tracks agent CLI sessions for resume; not an organizer phase handoff |

**Scoped null.** I found no shipped mechanism for: a blocking phase-boundary wait, an
audience-scoped protocol render, an organizer usage ledger, or a generated organizer brief. Sources
consulted: the `parley` dispatch table (`internal/app/app.go:51-114` — there is no `wait`, `usage` or
`organizer` verb), `parley run --help` and the printed usage block, `internal/driver/`,
`internal/protocolpacket/`, `internal/telemetry/`, `internal/store/`, `internal/sessionstore/`, the
v2.12.1 `SKILL.md` and its `references/`.

**Design alternatives I considered and set aside**

- *Ship the digest as a TUI improvement instead of stdout.* Rejected: the organizer is headless by
  construction; a TUI is unreadable to it, which is exactly why the existing digest has not helped.
- *A stored `handoff.md` per phase.* Rejected: stored summary sidecar (already-rejected pattern), and
  writing it would be a new mandatory obligation.
- *Cut inside §4 to reach the brief's 43–50 KB target.* Set aside: the marginal bytes come out of
  phase prose that never-cut protects per-phase, and the risk/benefit is poor once phase scoping is
  enabled.
- *Reduce polling by shortening the driver's internal poll interval.* Rejected: internal polling
  costs the organizer nothing; only organizer-visible tool calls carry resident context.
- *Enforce A by removing the facilitator from the roster.* Rejected: the roster is machine-global and
  shared across decks; a per-idea `facilitator:` field is the correctly-scoped mechanism.

## Concerns / open questions

1. **The largest measured lever is outside this release.** The study's §7 puts a 22 KB skills
   catalogue in *every* request (~5.6k tokens/request inferred) and names AGENTS.md "read in full"
   lines as a driver of post-compaction re-reads. Neither is in A–D, and neither is a parley-deck
   artifact. I am not proposing a scope change — the owner drew this line deliberately. I am
   flagging that if the release is later measured against total organizer spend rather than against
   the A–D criteria, it will look like a shortfall for reasons A–D never claimed to cover. The
   ledger from D is what will let that be seen rather than argued.

2. **Which audience view do we emit for a facilitator who is also a participant?** If
   `facilitator_participates: true`, the facilitator needs the participant reading set *and* §9/§11.
   Simplest safe answer: audience scoping applies only when the facilitator is not a participant;
   otherwise full context. I lean that way but want the other participants' view — getting this
   wrong silently narrows what a signing agent reads, which is the one failure mode here I consider
   serious.

3. **Timeout semantics vs the cache TTL.** I could not verify the "30-minute prompt-cache lifetime"
   from any source available in this worktree, and it is provider- and plan-specific. I propose
   configuration with a conservative default rather than a constant. If someone can supply a
   `PRIMARY` locator for the figure, seeding the default from it is fine — baking it into Go is not.

4. **`parley wait` and the §14 human brake.** A long blocking call is not an automated loop (it is
   driven by a human-started session), but the boundary deserves a sentence: `wait` must never
   *advance* a phase, only observe one. It is a read of the event log. I would treat any proposal to
   let `wait` take an action as out of scope.

5. **Is the skill's unguarded protocol copy in scope?** The three copies differ today — deck 108,400 B
   (`12e4b31c…`), `internal/protocol/defaults/` 108,154 B (`7dc3061c…`), skill `references/`
   108,244 B (`d45a7e86…`) — as expected for project-specific zones, but only the first two are
   guarded. Since C edits the skill anyway, adding a normalized drift check there is cheap. I raise
   it as a question rather than asserting it into scope, because a *required* check could be read as
   a new obligation; my own view is that a repo test is tooling, not an obligation on any agent.

6. **Digest fidelity when an artifact is malformed.** The existing digest deliberately never errors
   ("a display feature must never block protocol advancement" — `internal/driver/digest.go:44-47`).
   Once the organizer *relies* on the digest instead of reading files, silent degradation becomes a
   correctness question, not a display question. I propose the `fell_back` and `valid` columns be
   surfaced prominently and that `parley wait` exit `4` when any artifact is invalid — so degradation
   is loud. Worth a counter-proposal if someone thinks that over-blocks.

## Risks

| # | Risk | Severity | Mitigation |
|---|---|---|---|
| R1 | Audience scoping omits a rule the organizer then breaks | High | Never-cut floor unchanged; audience may only omit what phase/track/transport already omits; `packet check` negative test; omission index with triggers |
| R2 | The digest becomes a de-facto substitute for artifacts, eroding "raw artifacts canonical" | High | Digest carries every path; no model text; `wait` exits `4` on invalid artifacts; protocol text states the digest is a pointer, never evidence |
| R3 | `parley wait` hangs and the organizer silently stalls | Medium | Bounded `--timeout`; distinct exit `3`; partial digest printed on timeout; `driver.error` breaks the wait |
| R4 | D2 is read as re-litigating the rejected "self-reported usage" pattern | Medium | The ingest accepts a path, never a number; the ledger records source path + parser id; the distinction is enforced by the flag set, not by argument |
| R5 | SKILL.md split loses a rule into a reference nobody opens | Medium | Relocation only; a script proves every moved heading lands in exactly one reference; the slim core names every reference file |
| R6 | Measured saving underdelivers because the wins are cached-rate | Medium | Stated up front; the D ledger measures it per phase instead of asserting it |
| R7 | A ships as text and changes nothing, exactly as the 2026-09-15 instruction did | Medium | A is a preflight failure and a driver predicate, not a sentence; the regression test is that a deck without `facilitator:` is unchanged |
| R8 | Scope creep past A–D triggers the organizer's stop rule | Medium | Every item above maps to A, B, C or D; the "what I do not propose" list is written to make a drop or an addition easy to spot |
| R9 | Organizer brief drifts toward being stored | Low | Generated-only; test runs it against a read-only deck |
| R10 | `--audience` confused with `protocolRole` | Low | Distinct flag name; documented in the same sentence as `protocolRole` |
| R11 | Ingesting a 228 MB rollout blows up the organizer's context | Low, if designed | Streaming parser; ≤1 KB stdout; large-fixture test |

## §15 provenance

### Verdicts I issue

All `PRIMARY` verdicts below were produced in this worktree at HEAD `4ce8fa7` (skill worktree
`v2.12.1`, HEAD `d1e57d5`) with the command or file quoted. Tags per §15.2.

| # | Claim | Verdict | Tag | Evidence |
|---|---|---|---|---|
| 1 | The emitted full packet is byte-identical to the live deck protocol | CONFIRMED | `PRIMARY` | `shasum -a 256` of the packet body and `parley-deck/COOPERATION.md` both `12e4b31cd3f6c106…`; 1,380 lines / 108,400 B |
| 2 | The packet renderer has no audience/role input | CONFIRMED | `PRIMARY` | `internal/app/protocol_packet.go:48-60` — flags are `dir/phase/track/transport/idea/optimize/json/print/flag` only |
| 3 | `Role` in protocolpacket means the deck's `protocolRole`, not an agent role | CONFIRMED | `PRIMARY` | `internal/protocolpacket/packet.go:116`; `source.go:41,64-76` (`switch role` over `source`/`consumer`) |
| 4 | SKILL.md v2.12.1 never names `run`/`continue`/`consensus`/`status`/`preflight` | CONFIRMED | `PRIMARY` | `grep -oE 'parley [a-z-]+' skills/parley-deck/SKILL.md \| sort \| uniq -c` → `roster 9, protocol 4, init 4, agents 2, layer 1`; file is 990 lines / 62,411 B; `git describe` → `v2.12.1` |
| 5 | A deterministic LLM-free round digest already exists and is TUI-only | CONFIRMED | `PRIMARY` | `internal/driver/digest.go:10-48`; built `internal/driver/driver.go:488`; sole consumer `internal/tui/roundsummary.go:9-38` |
| 6 | Round/review/final artifact validators already exist | CONFIRMED | `PRIMARY` | `internal/protocol/roundartifact.go:23`, `reviewartifact.go:17`, `finalsections.go:92` |
| 7 | `parley wait` does not exist today | CONFIRMED | `PRIMARY` | dispatch table `internal/app/app.go:51-114`; `parley --version` → `parley 1.48.0` |
| 8 | `internal/telemetry/usage.go` has no kimi case | CONFIRMED | `PRIMARY` | cases at `:230` claude, `:252` codex, `:261` opencode; zcode envelope `:190,:286`; no kimi match |
| 9 | Compact state digests are small | CONFIRMED | `PRIMARY` | measured: `parley status --idea meta-protocol-change-lean-organizer --json` = 1,459 B; `parley status` = 10,897 B |
| 10 | The protocol already names a facilitator reading set | CONFIRMED | `PRIMARY` | `parley-deck/COOPERATION.md:34` |
| 11 | Nothing forces the facilitator to draft, implement or verify code | CONFIRMED | `PRIMARY` | `parley-deck/COOPERATION.md:407,409,443,437,890,1340-1343`; drafter restricted to participants at `internal/app/driver_consensus.go:82-100` |
| 12 | Three COOPERATION.md copies exist; the drift guard covers two | CONFIRMED | `PRIMARY` | deck 108,400 B `12e4b31c…`; `internal/protocol/defaults/` 108,154 B `7dc3061c…`; skill `references/` 108,244 B `d45a7e86…`; guard at `internal/protocol/drift_test.go:12` names only the first two |
| 13 | The never-cut list is ratified and per-phase | CONFIRMED | `PRIMARY` | `internal/protocolpacket/applicability.go:193-231`; map header in `parley-deck/meta/packet-applicability.yaml` states it is protocol and §7-governed |
| 14 | Section byte sizes used in the C table | CONFIRMED | `PRIMARY` | `awk 'NR>=a && NR<=b' parley-deck/COOPERATION.md \| wc -c` per row; boundaries from `grep -n '^## \|^### Phase'` |
| 15 | Codex rollout `token_count` is machine-written, streamable, and reports a 258,400-token window | CONFIRMED | `PRIMARY` | `~/.codex/sessions/2026/09/10/rollout-…01a08830-….jsonl`, 228 MB, 3,851 `token_count` records; first record quoted verbatim in D2 above |
| 16 | "Timeout below the provider's 30-minute prompt-cache lifetime" | UNVERIFIED | `RECALL` | No authoritative locator reachable from this worktree; provider- and plan-specific. Design consequence: configuration, not a constant |

Claim 16 is the only place where my design deviates from the brief's wording, and it deviates
*toward* the brief's intent rather than away from it.

### Testimony I transcribe but do not verdict

Per §15.1, I mark the following as **unverified testimony** and do not own it: every quantitative
figure from `source-context/organizer-token-study.md` §§1–5 (566.7M input tokens, 95.1% cached share,
~$1,008 and ~$3,870 estimates, 3,717/3,728 and 13,294 request counts, 50 compactions, the 30–40% /
33–39% / 42% / 9.1–12.2% attributions, the 54.8–71.7% optimized-packet band, the 37–77% fresh-context
simulation, and the per-agent telemetry ratios). The measured sessions are outside this worktree and
the study's `results/*.json` were not re-derived by me; the study states it was produced by six
readers and three independent verifiers, which is its own provenance, not mine. I rely on these
figures only to *rank* levers, never as an acceptance criterion — every criterion I propose is
measured at implementation time by a command in this repo. Two narrow corroborations I did obtain
first-hand are recorded as claims 15 and 8 above (the 258.4k window, and that kimi's missing usage is
a parser gap rather than adapter flakiness).

### Ownership and role notes

- I am an **owner** of every design proposition in `## Proposed approach` and issue no verification
  verdict on any of them (§15.1). The verdicts above are about the state of the tree, not about
  whether my design is a good one.
- I hold no facilitator role in this idea; `codex-1` is organizer only and is not a participant or
  signer, so §15.5's role-concentration disclosure does not apply to me here. If I am later the
  consensus or FINAL drafter, §15.5's `## Drafter position changes` applies to me then.
- Round-1 independence: written before reading `kimi-1` or `zcode-1`'s round-01 files. I read
  `00-prompt.md`, the study, the full protocol packet, and source in both worktrees — nothing else in
  `round-01/`.
- `EXEMPTION-CLAIM` check (§15.4): I make no claim to avoid a named known obstacle. Where I touch a
  previously-rejected pattern (D2 vs "self-reported usage as cost evidence"), I do not claim an
  exemption by adjective — I propose a mechanical discriminator (the ingest cannot accept a number)
  and invite the other participants to break it.
