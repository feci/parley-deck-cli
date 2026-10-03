---
idea: meta-protocol-change-quota-auto-exclude
status: final
author: claude-1
consensus-date: 2026-10-03
participants: [codex-1, claude-1]
implementer: codex-1
---

# FINAL: quota auto-exclusion with a floor of two non-facilitators

Role concentration (§15.5): the separately launched claude-1 participant process drafts this FINAL. claude-1
is also the declared facilitator (`facilitator_participates: true`); the claude-1 organizer process writes no
FINAL content. The driver launched this drafter ("driver: drafting FINAL via claude-1",
`runs/20261003T103229.878442000Z/driver-continue-4/stdout.log`).

Protocol context attestation: `context_mode=full`,
`source_sha256=b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388`,
`packet_sha256=b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388`, `fallback_reason` absent.
The shadow packet audit was not applied. HEAD at drafting: `7285b280400e324de67721adde8398fc66fbd317`.
Shared memory (OpenViking) was not consulted in this time-limited launch. This FINAL is built from the
idea's local canonical files only.

This is the design output of Phases 0 to 4 only. Nothing is implemented, published or merged by it.
Phase 5 starts only after the owner has reviewed this FINAL (the owner's standing rule, recorded in
`00-prompt.md`). The protocol text takes effect only after owner ratification and the attended
`parley protocol publish` (§7).

## Final plan / specification

### 1. The rule in brief

The rule is called **quota auto-exclusion**. When a participant's invocation fails, a provenance-verified
adapter recognizer captures the provider error, and that error states one of the following, the CLI removes
the participant from **this idea's quorum only**, without asking the owner:

- exhausted account credit or account quota;
- a reached named usage allowance whose stated reset is at least 60 minutes away;
- with no stated reset, a reached allowance of 24 hours or longer.

Two conditions must also hold. At least two distinct, usable non-facilitator participants must survive the
whole batch, and no protected role may be affected. The removal is recorded as an automatic `excluded:`
marker plus one durable batch transition record, and the owner gets one non-blocking inbox notice. Anything
not positively classified keeps today's user-confirmed path (fail closed). No `agents.toml` file is ever
written.

### 2. Items for owner ratification (§7)

Both participants recommend each of these. Each one is the owner's decision at ratification:

- **R1.** "Roster" in the owner's request means the per-idea quorum, not the machine or deck roster.
- **R2.** The rule is on by default for ideas created after ratification and delivery. Legacy ideas stay on
  confirmation.
- **R3.** The floor (2 non-facilitators) and the threshold (60 minutes) are fixed constants. Neither is
  configurable.
- **R4.** Delivery is in two stages under this one FINAL: stage 1 is kickoff, stage 2 is mid-idea. Both
  stages are required for `IMPLEMENTATION.md: complete`.
- **R5.** "One driving run per idea" is scoped as in section 10. That scope is the drafter's reading,
  answering a question raised in the claude-1 signoff.

### 3. Scope (Q1)

- **Kickoff.** The decision is made inside `parley run`, over the exact proposed participant set, after
  every readiness probe has returned. Standalone `parley preflight` only reports candidates and their
  evidence. It never applies an exclusion.
- **Mid-idea.** The decision is made at a settled participant-dispatch batch, after every affected writer
  has stopped, and before any next dispatch, signoff evaluation or close.
- **Out of scope.** The organizer's own exhaustion, organizer failover and replacement agents.
- **Staged delivery.** Stage 1 is kickoff and stage 2 is mid-idea.
  - Until stage 2 ships, the published protocol text marks mid-idea application "not yet in force" (the §7
    pattern).
  - `IMPLEMENTATION.md` cannot reach `status: complete` on stage 1 alone. If the owner authorizes stage 1
    only, that narrower authorization is recorded together with a linked stage-2 follow-up idea.

### 4. What counts as exhaustion (Q2): the authorizing predicate

An invocation is a candidate only if **all** of these hold:

1. **It failed:** exit ≠ 0, or a structured error result.
2. **No artifact, no later success.** It produced no valid completed artifact, and no later attempt in the
   same batch succeeded. A valid artifact or a later success always wins.
3. **Native-error provenance.** That adapter's recognizer captured the terminal error, and the recognizer's
   native-error provenance is established: located CLI behavior or source, plus recorded positive fixtures
   and adversarial negative fixtures.
   - The negatives cover an assistant quoting `API Error: 429 … Limit Exhausted`, a tool emitting the same
     text, and each of these followed by a failure with no valid artifact.
   - An adapter without established provenance is listed as **diagnostic-only, automatic exclusion
     unsupported**, and stays on the human path.
   - This applies to every adapter (claude/text, codex, kimi stream-json, zcode). An adapter name is not
     evidence.
   - `is_error` alone is not enough. The structured error must identify a terminal provider failure, not
     a local budget or tool failure.
4. **Exhaustion semantics.** The error's own text (for a gateway pass-through, the upstream message)
   states explicit account-credit or account-quota exhaustion, or an explicitly reached or exhausted named
   usage allowance. A message that only names a policy does not count.
5. **A stated reset at least 60 minutes from observation.** A known reset under 60 minutes takes the
   existing path, even when the message has a weekly or monthly phrase. A contradictory, past or
   unparseable reset gates. It is never relabeled `unknown` to gain eligibility.
6. **No stated reset.** This was OPEN in the consensus draft. Both participants closed it at signoff.
   - With no stated reset, only explicit account-credit or account-quota exhaustion qualifies, or an
     explicitly reached or exhausted named allowance of 24 hours or longer (daily, weekly, monthly).
   - An hourly or N-hour allowance under 24 hours qualifies only with a stated reset at least 60 minutes
     away.
   - This branch is a policy choice made with the remaining time unknown. It is not proof of a minimum
     remaining wait (codex-1 signoff).

**Never qualifies:**

- a bare 429, a generic `quota exceeded`, a bare `credit balance`;
- a 5xx or a "(reset after N)" without allowance semantics, or a long `Retry-After` alone;
- 400s and authentication errors;
- hangs and every watchdog class;
- artifact text, tool output and model prose.

There is no generic `reset after ⇒ quota` regex. Gateway-pool phrases ("Unavailable", "cooling down", "all
accounts exhausted") are not positive forms until the gateway's source establishes their cause.

**The 503 gap.** Align the preflight and runner classifiers so that a bare 503 is a provider gate in both,
not a `--yes`-excludable process failure. Generic provider classification stays separate from the
authorizing predicate. Nothing unrecognized falls into automatic process-failure exclusion.

**Disguised exhaustion** (a 400, or a hang from a disabled connection) stays on the human path.

### 5. What "roster" means, and return (Q3)

- **Per-idea quorum only.** Neither `~/.parley/agents.toml` nor `parley-deck/agents.toml` is written.
- **`roster_change_policy` does not gate the rule,** because no roster file is written. The protocol text
  states both facts.
- **The stated reset is a hint, never a timer.** There is no same-idea automatic rejoin and no quota
  polling.
- **Re-inclusion within the idea stays owner-confirmed,** under the existing catch-up rules. The agent gets
  a fresh readiness probe at the next idea.
- **Relaunch suggestion.** A notice may suggest one owner-authorized relaunch at the stated reset plus 5
  minutes, labeled a provider estimate. There is no suggestion when the reset is unknown, and no retry
  worker.

### 6. The minimum of 2 (Q4)

- **Who counts.** At least 2 distinct, usable **non-facilitator** participants must survive the whole
  batch. The facilitator never counts, even under `facilitator_participates: true`. Unresolved failures do
  not count as usable, and duplicate ids do not inflate the count.
- **All or nothing.** The batch is evaluated once, all or nothing, independent of the order in which
  failures arrived. Two agents exhausting at once are one batch.
- **Below the floor.** Nothing is applied, and one blocking escalation lists every candidate and the
  arithmetic. The §1 solo exception is unchanged.
- **The floor is fixed at 2.** An idea with two non-facilitators can never auto-exclude, and that includes
  this idea.

### 7. Interactions (Q5)

- **Designee or pin.** If the batch holds a per-idea designee or the pinned implementer, nothing in the
  batch is applied. The existing three-exit gate (§4 Phase 5) opens with the evidence prefilled.
- **Global-default designee.** Before a pin exists, today's fall-through applies unchanged.
- **Drafter.** A candidate who has started a canonical `consensus.md` or `FINAL.md` draft blocks automatic
  application. Otherwise today's drafter fallback runs over the current set.
- **Gates are re-evaluated, never waived.** The reviewer count per track, the LE-7/LE-11 two-reviewer
  close, independent goal-checker eligibility, `require_model_diversity` and `strict_gate`. A shortfall
  escalates.
- **Filed artifacts survive.**
  - An excluded agent's ❌ stands. Engaging its counter-proposal does not erase it. It stands until that
    agent withdraws it after an owner-confirmed re-inclusion, or until the owner rules and the ruling is
    quoted into the next artifact.
  - Its `DISPUTED` claims stand (§15.3: never resolved by count).
  - Its open CRITICAL and MAJOR findings, and every strict-gate finding, need an explicit disposition
    backed by independent evidence. "Author excluded" is never a disposition.
- **In-flight work.**
  - Partial or invalid files are preserved, identified as incomplete, and never promoted because a file
    exists at the canonical path.
  - A failed round resumes only through a new terminal round evaluation bound to the transition, after
    every survivor artifact has been validated. The earlier `round.incomplete` event and the failed
    invocation are kept.
  - Closed `FINAL.md` and `IMPLEMENTATION.md` artifacts stay frozen.

### 8. Integrity (Q6)

- **The only trigger** is a recognized provider error from the failed invocation itself.
- **Evidence first.** The evidence is recorded before application: the rule id, the invocation id, a
  scrubbed excerpt, the raw reset string, and the UTC observation and reset times.
- **Content, disagreement and slowness never trigger an exclusion.** A hang never triggers one either.
- **No dissenter removal.** Exclusion removes no filed objection, so it cannot be used to remove a
  dissenter.

### 9. Recording and notice (Q7)

**Membership.**

- `00-prompt.md` `participants:` is the **current** set.
- Kickoff exclusions are filtered out before the first `participants:` line, the creation event, the run
  manifest and the first dispatch are written.
- A mid-idea exclusion rewrites `participants:` inside the transition.
- The kickoff quorum and every authorized revision are kept in immutable membership history. `run.created`
  and participant artifacts are never rewritten.
- **Known and required signers.** `known` signers are the identities that were members under that recorded
  history. `required` signers, and every dispatch and await consumer, use the current `participants:`.
  Kickoff-excluded ids never become known.
- **No derived membership.** Nothing derives membership by subtracting `excluded:` lines. Repeated
  `excluded:` keys collapse when the frontmatter is read (see Context), so markers are display records and
  never authority.
- **Signoff gate.** The signoff-append gate stays narrow (`internal/consensus/consensus.go:236`). An
  excluded agent cannot append a signoff until an owner-confirmed re-inclusion puts it back.

**Transition.**

- **One record per batch.** There is one durable **batch-level** commit record, never one per agent. It
  carries:
  - the idea and run identity and the prior membership revision;
  - the before and after sets;
  - every candidate invocation id;
  - the resolved policy and scope, and the rule ids;
  - the observation and reset times and the scrubbed evidence references.
- **Validate, then reconcile.** The whole batch is validated under idea-scoped serialization before commit.
  Projections are reconciled before any further dispatch, signoff evaluation or close: the prompt, the
  manifest, captured consumers and the failed-round evaluation.
- **Pending state.** A committed but unreconciled transition is a recoverable pending state. Read-only
  `parley status`, `parley wait` and `parley organizer brief` report it and never repair it.
- **Later runs.** A new run for the same idea finds the record before it acts. Missing or contradictory
  history blocks.
- **Durability.** The commit needs checked durability and a fail-closed path for a truncated record.

**Marker.** The `excluded:` marker carries:

- an automatic marker;
- the rule id;
- the reset in UTC RFC3339, or `unknown`;
- the scrubbed raw provider reset string;
- the transition id;
- the recorded date.

It never says "confirmed". The raw string is kept because this idea's own kickoff miscopied a UTC+8
provider wall clock as local time: `reset_at` `2026-10-04T22:14:57.003Z` is 2026-10-05 00:14:57 CEST, not
06:14 (C5, raised by codex-1 in round 1 and confirmed by claude-1 in round 2).

**Notices.**

- One `blocking: no` owner notice per transition, deduplicated by transition id. It names who was excluded,
  the decisive reason, the reset hint, the survivors and any remaining gates.
- A floor, role or integrity failure gets one `blocking: yes` escalation instead.

**Surfaces.**

- `parley status`, `parley wait` and the organizer brief use the same current set. They list each automatic
  exclusion with its reset hint, plus any pending transition.
- `wait` keeps its exit contract: a successful reduction is not an exit-4 condition. It refreshes its
  participant list on every poll.

### 10. One driving run per idea: scope

**Decision (drafter's reading, ratification item R5).** "One driving run per idea" binds every driving run
of an idea whose **recorded scope includes mid-idea application** (stage 2), for the whole life of the run,
not only around the commit.

- A competing run of the same idea under a different run id is rejected, or waits (serialized), before it
  dispatches anything.
- It does **not** bind ideas whose recorded policy is off, legacy ideas, or kickoff-only ideas. So it is
  **not** a third behavior change that applies with the knob off.

**Basis.**

- The consensus Q8 lists exactly two changes that apply with the knob off (C1 and the bare 503). A
  deck-wide singleton would contradict that list.
- The hazard it prevents exists only where an automatic mid-idea transition can happen: a competing run
  dispatching from stale membership, or committing a conflicting batch.
- The kickoff decision is made inside the `parley run` that creates the idea. That last point is the
  drafter's reasoning, not an executed check.

**Owner alternative.** If the owner prefers a deck-wide singleton, it becomes a third knob-off behavior
change, with its own regression expectations. It must then be listed beside C1 and the bare 503 in the
release notes and tests.

**Mechanism.** `driver.acquireLock` today locks `RunDir/driver.lock` (`internal/driver/loop.go:30`), which
is a run-local lock. Reuse that primitive keyed by idea. A run-local lock alone is not enough.

### 11. Configuration (Q8)

- **The knob.** `[defaults].quota_auto_exclude` is a presence-aware boolean, with a deck override and a
  per-idea `quota_auto_exclude: false` in `00-prompt.md`.
- **Default.** On by default only for ideas created after ratification and delivery (R2). Legacy ideas
  stay on confirmation.
- **Recorded at kickoff.** The resolved policy **and the authorized scope** (kickoff-only, or kickoff plus
  mid-idea) are recorded at kickoff and reused by every later run of the idea.
  - A binary upgrade or a resume never widens scope. Widening goes through the existing owner-confirmed
    path.
  - A malformed config or an ambiguous record fails closed.
- **Fixed constants.** The floor (2) and the threshold (60 minutes) cannot be configured (R3).
- **Knob off** means "no automatic quorum reduction", not byte-identical behavior. Two fixes change
  existing behavior for every deck and carry their own regression expectations:
  - **C1:** `parley run --yes` no longer keeps or dispatches an excluded id;
  - **bare 503:** a bare preflight 503 can no longer be excluded with `--yes`.

### 12. Staffing for this idea's Phases 5 to 8 (owner direction)

- **Implementer: codex-1.** It is recorded in this file's frontmatter (`implementer:`, chain item 4) so the
  staffing does not depend on who drafted FINAL. The machine default resolves to the same agent:
  `~/.parley/agents.toml:15` reads `default_implementer = "codex-1"` at drafting (drafter's reading).
- **Reviewer: the separately launched claude-1 participant,** the single non-implementer reviewer
  (deliberation: "all non-implementers").
- **Organizer.** The claude-1 organizer stays organizer only. The §15.5 role concentration is disclosed in
  each phase.
- **Model diversity.** `require_model_diversity: true` is met: codex-1 runs OpenAI `gpt-6-astra` and
  claude-1 runs Anthropic `claude/claude-opus-5-5[1m]`.
- **No auto-completion.** Phases 5 to 8 run **attended**. Under `auto_implement` the driver refuses
  auto-completion with "fewer than two independent reviewers" (`parley-deck/COOPERATION.md:692–697`). That
  gate is not lowered, and two claude-1 processes never count as two reviewers.
- **The close.** The close takes both participants' review-consensus signoffs and current-tree criterion
  evidence. `internal/driver/impl.go:283–288` already prints the manual route ("add a reviewer or sign off
  manually").
- **kimi-1 and zcode-1** are not planned back for implementation or review (owner direction, quoted under
  Context).
- **Failures.** If codex-1 or claude-1 fails on an authentication, quota or credit error, the organizer
  stops and writes a blocking note with the verbatim error.
- **The rule cannot fire here.** Under this rule the idea has one non-facilitator (codex-1), so automatic
  exclusion could never fire for it.

### 13. Change inventory (Q9)

#### 13.1 Protocol hunks

The hunks are identical in all three copies:

- `parley-deck/COOPERATION.md` (this deck, protocolRole `source`);
- `internal/protocol/defaults/COOPERATION.md`, under the drift guard `TestEmbeddedDefaultMatchesLiveDeck`
  (`internal/protocol/drift_test.go`);
- the skill's bundled `skills/parley-deck/references/COOPERATION.md`, in the skill worktree
  `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/quota-auto-exclude-skill`.

1. **§9.0, under "Excluding"** (`COOPERATION.md:880`). Add the binding quota auto-exclusion sub-bullet. Its
   normative content is sections 3 to 9 and 11 above:
   - the signal and its provenance gate;
   - the 60-minute threshold and the 24-hour no-reset rule;
   - the non-facilitator floor of 2, evaluated per batch, all or nothing;
   - the role guards;
   - the automatic marker and the batch transition record;
   - the non-blocking notice;
   - per-idea return, with owner-confirmed re-inclusion and a next-idea probe;
   - the statement that `roster_change_policy` does not gate the rule because no roster file is written;
   - the policy key;
   - "mid-idea application: not yet in force" until stage 2 ships.
2. **§5, the sentence at `COOPERATION.md:748`.** "a mid-idea unavailability does not silently shrink
   quorum" gains "except the recorded §9.0 quota auto-exclusion, which is never silent". Add one bullet: an
   excluded participant's filed signoffs, ❌s, `DISPUTED` claims and findings survive its exclusion. It
   stays a known signer through the recorded membership history, while required signers follow the
   current `participants:`. Exclusion is never a withdrawal or a disposition.
3. **§4 Phase 5.** Add a cross-reference: a per-idea designee or the pinned implementer is never
   auto-excluded, and a qualifying failure opens the three-exit gate with the evidence prefilled. Amend
   the `COOPERATION.md:451` sentence ("Designations are validated against `participants:` only — nothing
   reads `excluded:` — so a §9.0 exclusion of the designee must also remove the id from `participants:`,
   or the designation stays live") to state the section-9 membership model:
   - `participants:` is the current set;
   - membership history is immutable;
   - known signers come from that history;
   - nothing derives membership from `excluded:` lines.
4. **§0 `[defaults]` list.** Add `quota_auto_exclude`. **Phase 0 template:** add the optional per-idea
   `quota_auto_exclude: false` line.
5. **One-line cross-references,** only where needed:
   - Phase 3, 6 and 7: historical signers keep their force, and review gates are re-evaluated and never
     waived;
   - the §9 checklist: read automatic-exclusion notices.
6. **`meta/protocol-changelog.md`:** one entry naming this idea.
7. **Packet map.** No new entry is needed for text inside §5 and §9.0, which are `include: always`
   (`parley-deck/meta/packet-applicability.yaml:98–99` and `:114–115`: claim D1 from codex-1's round 2,
   CONFIRMED by claude-1 in round 3).
   - Phase 0 (`:45–47`, `include: when`, phases [0]) and Phase 5 (`:65–67`, `include: when`, phases
     [5, 8]) are conditional. So the binding rule text lives in §9.0, and the Phase 0 and Phase 5 hunks
     are cross-references only.
   - Check packet rendering for phases 0, 5 and 8.
   - Any new heading needs a fresh map review.
8. **Ratification.** Owner ratification, then the attended `parley protocol publish` (§7).

#### 13.2 CLI, stage 1 (kickoff)

- **New `internal/telemetry/quota.go`:** the shared typed refinement, the per-adapter recognizers and
  their support table, with `internal/telemetry/quota_test.go`.
- **Callers:** `internal/telemetry/usage.go`, `internal/runner/failclass.go`, `internal/runner/telemetry.go`,
  the terminal-result paths in `internal/runner/runner.go`, and `internal/app/preflight_liveness.go`.
- **`internal/app/preflight.go` and `internal/app/app.go`:** the C1 filter, the resolved policy and scope,
  the batch floor at kickoff, and status.
- **`internal/config/runtime.go`:** the presence-aware knob.
- **Creation and records:** `internal/protocol/workspace.go` and `internal/runcontrol/runcontrol.go`
  (filtered first `participants:` line, creation event and marker), plus the policy, scope and evidence
  fields in `internal/runmanifest/manifest.go` and `internal/runstate/runstate.go`.
- **Surfaces:** `internal/app/wait.go` and `internal/app/organizer.go`.

#### 13.3 CLI, stage 2 (mid-idea)

- **New `internal/driver/quota.go`,** with `internal/driver/quota_test.go`: the batch decision, the
  transition and the idea-scoped lock. Plus `internal/driver/driver.go`, `internal/driver/loop.go` and
  `internal/driver/phasedigest.go`.
- **Replay:** `internal/runstate/runstate.go` and `internal/runmanifest/manifest.go`.
- **`internal/consensus/consensus.go`:** the history-backed `known` at `:122` and `:258`.
- **`internal/app/driver_consensus.go` and `internal/app/driver_impl.go`:** the reviewer set, the designee
  gate and the drafter guard.
- **Store durability,** as the commit needs it (`internal/store/events.go`).
- **Runner evidence paths:** `internal/runner/phase58.go` and `internal/runner/acp.go`.
- **No change** to `internal/agents/discover.go`.

#### 13.4 Skill

These live in the skill worktree and nobody edits them in this design run:

- `skills/parley-deck/SKILL.md`: the organizer reads the notice, never re-includes silently, and never
  infers exhaustion itself;
- `skills/parley-deck/references/ROSTER_AND_PROTOCOL.md`;
- the bundled `skills/parley-deck/references/COOPERATION.md`.

#### 13.5 Tests

Test files to extend:

- `internal/runner/hardening_test.go`;
- `internal/telemetry/usage_test.go`;
- `internal/app/preflight_liveness_test.go` and `internal/app/preflight_test.go`;
- `internal/config/runtime_test.go`;
- `internal/runcontrol/runcontrol_test.go`;
- `internal/runstate/runstate_test.go` and `internal/runmanifest/manifest_test.go`;
- `internal/consensus/consensus_test.go`;
- `internal/app/wait_test.go` and `internal/app/organizer_test.go`;
- the two new `quota_test.go` files.

Required cases:

- native versus assistant-quoted or tool-emitted errors, and supported versus unsupported adapters;
- short, unknown, contradictory and past resets;
- "hourly" and "5-hour" limits with no stated reset, which must **not** qualify;
- a valid artifact or a later success wins;
- every permutation of a simultaneous batch, 4→2 and 3→1, duplicate ids, and the facilitator floor;
- protected roles (per-idea designee, pin, started drafter) and the global-default fall-through;
- kickoff filtering and replay, including that `parley run --yes` never dispatches an excluded id;
- a bare preflight 503 that cannot be excluded with `--yes`;
- a stage-1-to-stage-2 upgrade that does not widen an existing idea's scope;
- interrupted and truncated commits, duplicate replay, and competing runs with different run ids for the
  same idea;
- a stub at the canonical round path that does not count;
- after removal, an excluded agent's ❌ yields `TriageBlocked`, not `TriageMalformed`;
- a kickoff-excluded id that never becomes a known signer;
- retained findings, and the reviewer, diversity and goal-check gates;
- notices, `status`, `wait` and the brief agreeing.

Recorded negative fixtures from this idea:

- the claude/text 503 "Unavailable (reset after 55m 29s)" one-line stdout (incident 7);
- the codex 503 "Unavailable (reset after 5h 51m 11s)" (incident 4);
- the codex 401 "invalidated oauth token … (reset after 12s)" that has a valid artifact beside it
  (incident 8).

#### 13.6 Checks to run (future, none run in this design run)

- `go build ./...`, `go vet ./...` and `go test ./...`, or at least every affected package;
- `go test ./internal/protocol -run TestEmbeddedDefaultMatchesLiveDeck`;
- a byte comparison of the skill's bundled copy against `parley-deck/COOPERATION.md`;
- the packet-rendering check for phases 0, 5 and 8;
- in the skill repository, `test/installer.test.js` and `test/lean-organizer.test.js`.

#### 13.7 Deliberately left out

- detecting disguised exhaustion (a 400, or a hang);
- gateway introspection and any OmniRoute-specific rule;
- timers, quota polling, retries, backoff and predictive budgeting;
- billing APIs and purchases;
- automatic same-idea re-inclusion;
- any `agents.toml` write;
- a configurable floor or threshold;
- automatic replacement of a designee, a pin or a started drafter;
- organizer self-exhaustion and failover;
- the claude `--output-format json` switch, which is a scoped follow-up because it changes the runner's
  fallback-artifact path;
- model, effort, provider, gateway or credential changes (non-goals of the brief).

The 300 to 450 LOC estimate stays UNVERIFIED. Full mid-idea support is larger than the owner's "small".

#### 13.8 Deferred to implementation

These are bound by the requirements and tests above.

- **The form of the batch commit.** Either a single committed event whose unreconciled projections read as
  pending, or an intent record plus an applied marker. Both participants accept either.
- **The adapter support table.** Establish native-error provenance per adapter. An adapter that cannot be
  verified ships unsupported. Two points remain UNVERIFIED:
  - whether `claude -p` can exit non-zero with model text;
  - what OmniRoute's "Unavailable" denotes.
- **Idea-scoped serialization** (section 10), plus durability where the commit needs it.
  `internal/store/events.go:52–65` opens with `O_APPEND`, writes the line and returns, with no sync call
  (codex-1 round 3).
- **Names and fields.** The exact marker grammar, the event names, and the policy and scope fields.
- **Packet rendering.** The packet-rendering check for the Phase 0 and Phase 5 cross-references.
- **Stubs at canonical round paths.** What today's validator does with a stub at a canonical round path is
  unchecked. The requirement is to judge on content.
- **Incident 6.** Recover its raw terminal error and kimi-1's actual role, if the source is reachable.
  Otherwise the conditional row stays.
- **Scrubbing.** Whether request and correlation ids are kept in a scrubbed excerpt.
- **Implementation plan.** codex-1 writes the `IMPLEMENTATION.md` plan and checklist before any code
  (§4 Phase 5).

### 14. Incident table

The table separates "the message would qualify" from "a reduction would execute".

| # | Incident | Outcome under this rule |
|---|---|---|
| 1/2 | zcode-1 429 "Weekly/Monthly Limit Exhausted", reset +49 h (kickoff) | The message qualifies only through a provenance-verified zcode recognizer. At the actual kickoff kimi-1 was unresolved, so codex-1 was the only usable non-facilitator survivor. The floor fails, and the owner is asked. |
| 3 | kimi-1 probe hang | Human path. |
| 4 | codex-1 503 "Unavailable (reset after 5h 51m 11s)" | Human path: gateway-local text, and the process was an organizer. |
| 5 | kimi-1 400 from a disabled connection | Human path. |
| 6 | kimi-1 403 "weekly usage limit" during an implementation attempt | If kimi-1 was the designee or the pinned implementer, the three-exit gate opens. Otherwise it needs raw terminal evidence from a verified recognizer, plus the floor. The excerpt is a summary, so no automatic outcome is asserted, and the evidence check is deferred (13.8). |
| 7 | claude-1 503 "Unavailable (reset after 55m 29s)", this idea, round 2 | Human path: gateway-local text, under 60 minutes, claude/text unsupported, and the floor. |
| 8 | codex-1 401 "invalidated oauth token … (reset after 12s)", this idea, round 3 | Never a candidate: a valid artifact exists, and it is an authentication error. Row 8 is a claude-1 reading (see Context). |

No recorded incident would have been automated in its actual context. The benefit is prospective. For
example, take the owner's default fleet: three non-facilitators and a non-participating facilitator. If one
participant hits an explicit long-window limit through a verified recognizer and the other two are usable,
the run continues with a notice instead of a question.

## Purpose / user-visible outcome

The owner asked, in Slovak (verbatim):

> "nastartuj /parley-deck a navrhni upravu, resp. malu zmenu, kde vyradis agentov z rosteru automaticky,
> ak nemaju dost tokenov, resp. kreditov, min. roster je 2"

The relay's translation: "Start /parley-deck and propose an adjustment, or rather a small change, where you
drop agents from the roster automatically if they do not have enough tokens, or rather credits. The minimum
roster is 2." (Original language: Slovak.)

What the owner sees after delivery:

- **A run continues without a question.** In an idea with three or more non-facilitator participants, a
  participant whose provider reports verified long-window exhaustion is dropped from that idea without a
  question. The owner gets one non-blocking inbox notice that names who was dropped, why, the provider's
  reset estimate, the survivors and any remaining gates.
- **The same view everywhere.** `parley status`, `parley wait` and the organizer brief show the current
  quorum and each automatic exclusion.
- **Fail closed everywhere else.** The owner is still asked whenever the floor of two non-facilitators would
  be broken, or the dropped agent is the designated or pinned implementer or a started drafter. The owner is
  also asked whenever the error is ambiguous: a bare 429, a 5xx or a timer without allowance semantics, a
  400, an authentication error, a hang, or an adapter whose error provenance is not verified.
- **No permanent change.** No machine or deck roster entry changes. The agent is re-probed at the next idea
  and rejoins the same idea only with the owner's confirmation.
- **Two fixes for every deck,** even with the knob off. `parley run --yes` stops keeping and dispatching an
  excluded agent (C1), and a bare preflight 503 can no longer be excluded with `--yes`.
- **Honest limits.** Two-participant ideas, this one included, never auto-exclude. claude/text and gateway
  locks keep reaching the owner until their recognizers are verified. None of the recorded incidents would
  have been automated in its actual context.

## Context & orientation

### Owner directions for this idea (verbatim, original language Slovak)

- **Quorum decision at kickoff.** Selected "codex-1 + claude-1 (Recommended)". kimi-1 and zcode-1 are
  excluded for this idea only under §9.0 with owner confirmation (`00-prompt.md` frontmatter `excluded:`).
- **About 18:55 CEST.** "sakra tak pouzi len claude a codex". Translation (the relay's): "Damn, then just
  use claude and codex." The quorum stays codex-1 and claude-1 through FINAL, and implementation and review
  are planned with those two only.
- **About 20:35 CEST.** "pokracuj". Translation: "Continue." This answered the codex-1 gateway 401 note:
  continue with consensus, signoffs and FINAL, and stop with a new blocking note if codex-1 fails again on an
  authentication, quota or credit error.

### Current behavior at HEAD (provenance as recorded in the round files)

This FINAL summarizes statuses and originates no verdict (§15.1).

- **§9.0 needs explicit confirmation.** Excluding an unavailable agent "requires **explicit user
  confirmation**" (`parley-deck/COOPERATION.md:880`).
- **§5.** "a mid-idea unavailability does not silently shrink quorum" (`:748`).
- **Phase 5.** "Designations are validated against `participants:` only — nothing reads `excluded:` — so a
  §9.0 exclusion of the designee must also remove the id from `participants:`, or the designation stays
  live" (`:451`).
- **LE-7/LE-11.** The close refuses auto-completion with "fewer than two independent reviewers"
  (`:692–697`, PRIMARY, codex-1 round 3).
- **C1: `parley run --yes` keeps excluded ids in `participants:` and dispatches them.**
  - claude-1 round 2 C1: a static reading of `internal/app/app.go:1922–1943`.
  - codex-1's executed probe showed `agent.started agent=kimi` after a confirmed exclusion. claude-1
    disclosed this as corroboration (round 3, D5) and issued no PRIMARY verdict on the dispatch claim.
- **The 503 gap.**
  - claude-1 round 1 V4 (PRIMARY, executed overlay test): the runner classifier labels the codex 503 "reset
    after 5h 51m 11s" as `overloaded`.
  - claude-1 round 1 V5: preflight treats a bare 503 as a `--yes`-excludable process failure.
  - codex-1 round 3 item 7 accepts the resulting fix.
- **Known and required signers.**
  - `internal/consensus/consensus.go:122–123` passes `idea.Participants` as known and
    `reviewConsensusVoters(...)` as required.
  - `:236` is the narrow signoff-append gate.
  - `:553–554` and `:587–590`: a signer missing from `known` is an error, and errors outrank blocks, so its
    ❌ yields `TriageMalformed`. This is claim D4, CONFIRMED with that scope limit by claude-1 round 3.
- **`wait`.**
  - `internal/app/wait.go:156` builds the digest from the startup participant list (D3).
  - `:10–19`: a `blocking: no` note is not an exit-4 condition (D2).
  - Both claims are codex-1's from round 2, CONFIRMED PRIMARY by claude-1 round 3.
- **Membership consumers.** `internal/protocol/workspace.go:356` fills `IdeaStatus.Participants` from
  `participants:`. That set feeds `consensus.go:123,176,298`, `runner.go:1003`, `wait.go:156` and
  `organizer.go:168,249` (claude-1 round 3, PRIMARY).
- **Repeated keys collapse** (`workspace.go:376` and `:394`): the last repeated frontmatter key wins, and
  this idea's own `00-prompt.md` has two `excluded:` lines.
  - This is a claude-1 reading (consensus blind spot 2), not independently confirmed.
  - codex-1's round 3 reached the history-based `known` independently, so no decision rests on this reading
    alone.
- **The lock is run-local.** `internal/driver/loop.go:30` locks `filepath.Join(d.cfg.RunDir,
  "driver.lock")`, a run-local lock (codex-1 round 3).
- **Event appends are not synced.** `internal/store/events.go:52–65` has no sync, and `:84–85` rejects
  malformed JSON on load (codex-1 round 3).
- **Manual close route.** `internal/driver/impl.go:283–288` prints "add a reviewer or sign off manually"
  (codex-1 round 3).
- **claude/text output.**
  - The claude headless argv uses `--output-format text` (`internal/agents/discover.go:234`).
  - `internal/telemetry/usage.go:423–440` recognizes only JSON formats.
  - The round-2 claude-1 failure (`runs/20261003T103229.878442000Z/round-02/claude-1/relaunch-1/`) exited 1
    after 127 s, with an empty stderr and the provider error as its whole one-line stdout.
  - codex-1 round 3 CONFIRMED all three (PRIMARY). The stdout observation holds for that invocation only;
    it does not establish that claude reports provider errors only on stdout.
- **Row 8 evidence.** `runs/20261003T103229.878442000Z/round-03/codex-1/exit.json` reads
  `"seconds":702,"exit":1`, and `stderr.log:8005` holds the 401 "invalidated oauth token … (reset after
  12s)" line, yet `round-03/codex-1.md` is complete (22646 B).
  - This is a claude-1 reading (consensus blind spot 1), not independently confirmed.
  - organizer-notes.md records the same 401 as an orchestration fact.

### Verification status

- **No verdict conflicts and no `DISPUTED` claims.** consensus.md records none (§15.3), so no claim enters
  this FINAL under a DISPUTED heading.
- **Drafter readings.** Blind spots 1 to 3 and incident row 8 are claude-1 readings with locators. The
  claude-1 signoff re-read them at `f1054d8` and, sharing the identity, issued no verdict. They are cited
  here as claude-1 readings, not as independently confirmed.
- **UNVERIFIED.**
  - the 300 to 450 LOC estimate;
  - whether `claude -p` can exit non-zero with model text;
  - what OmniRoute's "Unavailable" denotes;
  - incident 6's raw terminal error and kimi-1's role in it;
  - today's validator behavior on a stub at a canonical round path.

### Alternatives disposition (consistent with consensus.md)

| ID | Alternative | Disposition |
|---|---|---|
| ALT-1 | One shared refinement over the existing classifiers | Adopted, gated per adapter by established provenance. |
| ALT-2 | Stdlib reset parsing, with no timer | Adopted, at 60 minutes. |
| ALT-3 | The existing exclusion record, with an automatic marker and real filtering | Adopted. Rejected: `parley roster set --state inactive`, and `--yes` as a blanket quota waiver. |
| ALT-4 | The existing persistence primitives, keyed by idea with checked durability | Adopted. |
| ALT-5 | The known/required signer split, with `known` from history | Adopted. |
| ALT-6 | The existing retries, manual pause and blocking escalation | Adopted (kept) for everything not positively classified. |
| ALT-7 | Preflight-only exclusion | Adopted as stage 1. Rejected as full completion. |
| ALT-8 | The claude/text one-line matcher as present authority | Rejected. The `--output-format json` switch is deferred. |
| ALT-9 | Deriving membership by subtracting `excluded:` lines | Rejected. |
| ALT-10 | A stable P0 in `participants:` plus a reducer at every consumer | Rejected. |
| ALT-11 | Credit querying, billing APIs and predictive budgeting | Rejected. |
| ALT-12 | A reset window alone as the trigger, or a generic `reset after ⇒ quota` regex | Rejected. |

### Correlated agreement (§15.6(b))

- **One architectural family.** codex-1 (OpenAI) and claude-1 (Anthropic) are different model families.
  Even so, their two proposals are **one architectural family** over the same driver, and they converged
  through three cross-review rounds. Their unanimity is a shared prior, not independent evidence that the
  failure paths are safe.
- **What would prove the position wrong:**
  - a provider that uses allowance wording for a short throttle with no reset;
  - a content-triggered exclusion;
  - a lost veto;
  - a consumer that dispatches from a stale list;
  - a replay that launches mixed membership.
- **Guards.** The adversarial fixtures, the `TriageBlocked` test and the replay tests guard against these.

### Run facts for whoever implements

- **Workspace and transport.** The CLI worktree is
  `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/quota-auto-exclude`, on branch
  `quota-auto-exclude`. The transport is local canonical files, as a per-run override of §11.B, with no
  PR, push or merge.
- **Driver gaps.** organizer-notes.md records driver gaps 1 to 12 met in this run. Gap 11 is a legacy run
  record without an idea identity, and it blocks driver-launched cross-review rounds until the owner's
  accounting decision (`parley budget migrate apply`). Expect to use the recorded fallbacks, or to have that
  decision first, when running Phases 5 to 8.

## Observable acceptance criteria

Each criterion is checked on the current tree by the non-implementer reviewer (claude-1) and recorded in
`IMPLEMENTATION.md` `## Validation evidence`. Stage-2 criteria are AC10 to AC14 and the mid-idea parts of
AC9 and AC15.

1. **AC1, protocol text.**
   - The hunks of 13.1 are present in `parley-deck/COOPERATION.md`, and
     `go test ./internal/protocol -run TestEmbeddedDefaultMatchesLiveDeck` passes.
   - The skill's `skills/parley-deck/references/COOPERATION.md` is byte-identical to the deck copy.
   - `meta/protocol-changelog.md` has the entry.
   - Packets rendered for phases 0, 5 and 8 contain the §9.0 rule text and the conditional
     cross-references.
   - Until stage 2 ships, the mid-idea paragraph says "not yet in force".
2. **AC2, positive classification.** A provenance-verified zcode recognizer classifies the recorded 429
   "Weekly/Monthly Limit Exhausted" with `reset_at` 49 h away as a candidate, and records the rule id, the
   invocation id, a scrubbed excerpt, the raw reset string and the UTC times.
3. **AC3, negative classification.** Each of these is **not** a candidate:
   - a bare 429, a generic `quota exceeded`, a bare `credit balance`;
   - the 503 "Unavailable (reset after 5h 51m 11s)" and the 503 "Unavailable (reset after 55m 29s)";
   - the 401 "invalidated oauth token … (reset after 12s)";
   - a 400 from a disabled connection;
   - a hang, and every watchdog class;
   - an assistant-quoted, and a tool-emitted, `API Error: 429 … Limit Exhausted` followed by a failure with
     no valid artifact;
   - "hourly" and "5-hour" limits with no stated reset;
   - a weekly phrase with a stated reset 30 minutes away;
   - a past, contradictory or unparseable reset;
   - a long `Retry-After` alone;
   - any failure with a valid artifact, or with a later success in the same batch.
4. **AC4, unsupported adapters.** An adapter without established provenance never produces an automatic
   exclusion. The support table lists it as diagnostic-only, and claude/text stays there until it is
   verified.
5. **AC5, the floor.**
   - Four non-facilitators with two candidates → both are excluded (4→2).
   - Three non-facilitators with two candidates → nothing is applied, and exactly one blocking escalation
     lists both candidates and the arithmetic (3→1).
   - Every permutation of the same simultaneous candidates gives the same result.
   - Duplicate ids do not inflate the count, and unresolved failures do not count as usable.
   - The facilitator never counts, even under `facilitator_participates: true`.
   - An idea with two non-facilitators never auto-excludes.
6. **AC6, kickoff (C1).**
   - In `parley run`, an excluded id never appears in the first `participants:` line, in `run.created`,
     in the manifest, or in any `agent.started` event.
   - `parley run --yes` with a confirmed exclusion never dispatches the excluded id.
   - Standalone `parley preflight` writes no exclusion.
7. **AC7, the bare 503.** A bare preflight 503 is a provider gate that `--yes` cannot exclude, and the
   runner and preflight classify it the same way.
8. **AC8, protected roles.**
   - A batch holding a per-idea designee or the pinned implementer applies nothing, and opens the
     three-exit gate with the evidence prefilled.
   - A global-default designee before a pin falls through as today.
   - A candidate who has started a canonical consensus or FINAL draft blocks application.
9. **AC9, filed artifacts survive.**
   - After a mid-idea exclusion, the excluded agent's ❌ yields `TriageBlocked`, not `TriageMalformed`.
   - Its `DISPUTED` claims and its open CRITICAL, MAJOR and strict-gate findings remain and need an
     explicit disposition.
   - It cannot append a new signoff.
   - A kickoff-excluded id never becomes a known signer.
10. **AC10, gates.** After an exclusion, the reviewer count per track, the LE-7/LE-11 close, goal-checker
    eligibility, `require_model_diversity` and `strict_gate` are re-evaluated. A shortfall escalates, and
    none is waived.
11. **AC11, transition durability.**
    - With a fault injected at each write of the batch commit, the result is either no transition or a
      recoverable pending state. A pending state blocks dispatch, signoff evaluation and close until it is
      reconciled.
    - A truncated record fails closed.
    - Replaying a committed transition twice is idempotent and sends one notice.
    - No dispatch ever uses mixed membership.
12. **AC12, serialization.** Two driving runs with different run ids on the same stage-2-scoped idea: the
    second is rejected or serialized before it dispatches anything. A policy-off idea's runs are
    unaffected.
13. **AC13, partial artifacts.**
    - A stub at a canonical round path left by a failed invocation is preserved, identified as
      incomplete, and never counted.
    - The round resumes only through a new terminal evaluation after the survivor artifacts validate.
    - The earlier `round.incomplete` event and the failed invocation are kept.
14. **AC14, history.**
    - `run.created` and participant artifacts are byte-unchanged by a transition.
    - The kickoff quorum and each revision are readable from the membership history.
    - A new run for the same idea finds the record before acting. Missing or contradictory history blocks.
15. **AC15, records and surfaces.**
    - The marker carries every field of section 9 and never says "confirmed".
    - Exactly one `blocking: no` notice is sent per transition. Floor, role and integrity failures produce
      `blocking: yes`.
    - `status`, `wait` and the organizer brief show the same current set, each automatic exclusion with
      its reset hint, and any pending transition, without repairing anything.
    - `wait` does not exit 4 on a successful reduction, and it refreshes its participant list on every
      poll.
16. **AC16, configuration.**
    - Per-idea `false` overrides a `true` default, and the deck override works.
    - A malformed config fails closed.
    - Ideas created before delivery stay on confirmation.
    - The policy and scope recorded at kickoff are reused by later runs.
    - A stage-1-to-stage-2 binary upgrade does not widen an existing idea's scope.
17. **AC17, return.**
    - Both `agents.toml` files are byte-unchanged across an exclusion.
    - There is no same-idea timer or rejoin.
    - The next idea re-probes the agent.
    - The relaunch suggestion is the reset plus 5 minutes, labeled a provider estimate, and is absent when
      the reset is unknown.
18. **AC18, integrity.** No code path triggers an exclusion from artifact content, tool output, a
    disagreement or elapsed time.
19. **AC19, completion.** `IMPLEMENTATION.md` reaches `status: complete` only with both stages delivered,
    or with a recorded owner authorization for stage 1 only and a linked stage-2 follow-up idea.
20. **AC20, checks.** Every check in 13.6 passes on the current tree.
21. **AC21, close.** The close is attended. It needs both participants' review-consensus signoffs and
    current-tree criterion evidence, with no unattended auto-completion.

## Idempotence & recovery

**This design artifact.**

- `FINAL.md` is static. If it is later invalidated, open `meta-protocol-change-quota-auto-exclude-v2`; do
  not edit this file.
- Nothing in the design run changed CLI code, the skill or any COOPERATION.md copy, so there is nothing to
  roll back.

**State that matters in the implementation:**

- each idea's recorded policy and authorized scope;
- the immutable membership history;
- the batch transition records and their evidence records;
- the notice deduplication keys (transition ids);
- the preserved failed invocations and partial artifacts.

**Safe to rerun:**

- **Classification** is a pure function of the recorded terminal evidence, so re-classifying gives the
  same result.
- **The kickoff decision** is recomputed inside `parley run` before the idea is created. A crash before
  creation persists nothing, and after creation the first `participants:` line, `run.created` and the
  manifest agree by construction.
- **Replay.** Replaying a committed transition reconciles the projections idempotently.
- **Notices** are deduplicated by transition id.
- **`status`, `wait` and the organizer brief** are read-only and never repair anything.

**Recovery:**

- **Pending transition.** A committed but unreconciled transition is a pending state. Resume reconciles it
  before any dispatch, signoff evaluation or close.
- **Bad history.** A truncated record, or missing or contradictory history, fails closed with a blocking
  escalation.
- **Competing runs.** A competing run of a stage-2-scoped idea is rejected or serialized (section 10).

**Needs a human gate:**

- a batch below the floor;
- a designee, a pin or a started drafter in the batch;
- every error not positively classified;
- re-inclusion within the idea;
- widening an idea's recorded scope;
- stage-1-only completion;
- missing or contradictory history;
- owner ratification of R1 to R5;
- the attended `parley protocol publish`.

**Rollback:**

- Set `[defaults].quota_auto_exclude = false` at deck or machine level. Ideas created after that use
  confirmation only. An existing idea keeps its recorded policy.
- Revert the CLI commits on the implementation branch.
- A published protocol release is write-once, so reverting the protocol text takes a new
  meta-protocol-change idea and a new attended release (§7).

## Known risks / de-risking

- **Provider wording drift.** A provider could use allowance wording for a short throttle. De-risked by the
  per-adapter provenance gate, the adversarial fixtures, the 60-minute threshold and the 24-hour rule when
  no reset is stated. An unrecognized text misses and falls back to the owner.
- **Content-triggered or dissent-driven exclusion.** De-risked because the only trigger is the failed
  invocation's own recognized terminal error, with AC18, and because filed ❌s, `DISPUTED` claims and
  findings survive (AC9).
- **A lost veto through a missed `known` site.** De-risked by history-backed `known` and the
  `TriageBlocked` test. A missed site fails loudly as `TriageMalformed`, never as a silent close.
- **Stale consumers or mixed-membership replay.** De-risked by current `participants:`, reconciliation
  before any dispatch, the `wait` refresh, the idea-scoped lock and the replay and fault-injection tests
  (AC11 to AC14).
- **Crash or truncation mid-transition.** Event appends are not synced today
  (`internal/store/events.go:52–65`). De-risked by the checked-durability requirement and the fail-closed
  truncated-record path.
- **Size.** The protocol rule is small, but correct mid-idea recovery is not. The 300 to 450 LOC estimate
  is UNVERIFIED. Staging gives the owner a smaller first change, but stage 1 alone cannot complete the
  idea.
- **Behavior changes with the knob off.** C1 and the bare-503 fix change behavior for every deck. Each has
  its own regression expectation (AC6, AC7) and goes in the release notes. A deck-wide singleton (section
  10, owner alternative) would add a third.
- **Limited near-term benefit.**
  - Two-participant ideas never auto-exclude.
  - claude/text and gateway-local locks keep reaching the owner.
  - No recorded incident would have been automated.
  - This is stated so the owner can weigh the size against the benefit at ratification.
- **Shared account and gateway** (consensus blind spot 3, a claude-1 reading). The organizer and a
  participant can share one account and one gateway, as claude-1 does in this idea. Exhaustion can then
  stop orchestration and the participant together. This is out of scope, and it is stated here so nobody
  assumes the rule covers it.
- **Review depth for this idea.** There is a single non-implementer reviewer (claude-1), with the §15.5
  role concentration. De-risked by the attended close, refutation-default review against AC1 to AC21, and
  the existing refusal to auto-complete with fewer than two independent reviewers.
- **Correlated agreement.** The two proposals are one architectural family (Context, §15.6(b)). The guards
  are the fixtures and the `TriageBlocked`, replay and serialization tests, not the agreement itself.
- **Secrets in evidence.** Excerpts are scrubbed. Whether request and correlation ids stay in them is
  deferred (13.8). No token or credential is ever recorded.

## References

- **Kickoff and owner brief:** `./00-prompt.md`, including the owner directions, the quorum decision, the
  readiness table and the C5 correction note.
- **Owner brief, copied verbatim:** `./source-context/ORGANIZER-BRIEF.md`.
- **Incident evidence, with provenance labels:** `./source-context/EVIDENCE-quota-incidents.md`.
- **Raw readiness output:** `./source-context/preflight-2026-10-03T1223.json`.
- **Consensus:** `./consensus.md`, with the signoffs of codex-1 ✅ ACCEPT and claude-1 ✅ ACCEPT.
- **Rounds:** `./round-01/`, `./round-02/` and `./round-03/` (codex-1.md and claude-1.md in each).
- **Orchestration record:** `./organizer-notes.md`, covering driver gaps 1 to 12, the launches, the
  failures and the probes.
- **Run evidence:** `parley-deck/runs/20261003T103229.878442000Z/`, including:
  - `round-02/claude-1/relaunch-1/` (incident 7);
  - `round-03/codex-1/` (incident 8);
  - `probe-codex-1929/` and `probe-codex-1957/`;
  - `driver-continue-4/` (signoffs and this FINAL's launch).
- **Archived owner escalations and answers:** `parley-deck/inbox/archived/`, for the quota, codex-auth and
  worktree-prune notes.
- **Protocol:**
  - `parley-deck/COOPERATION.md`: §0 `[defaults]`, §4 Phase 0 and Phase 5 (`:451`), LE-7/LE-11
    (`:692–697`), §5 (`:748`), §7, §9.0 (`:880`) and §15;
  - `parley-deck/meta/packet-applicability.yaml` (`:45–47`, `:65–67`, `:98–99`, `:114–115`);
  - `parley-deck/meta/protocol-changelog.md`.
- **Protocol copies:** `internal/protocol/defaults/COOPERATION.md` and `internal/protocol/drift_test.go`;
  the skill worktree `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/quota-auto-exclude-skill`.
- **Code locators** are listed under Context and section 13. They were read at HEAD `7285b28`.
  `git diff --stat 22cfddf HEAD -- internal/ cmd/ parley-deck/COOPERATION.md
  parley-deck/meta/packet-applicability.yaml` is empty at drafting, which extends the claude-1 signoff's
  "nothing ... has changed since `22cfddf`" to this HEAD.
