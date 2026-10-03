---
agent: claude-1
idea: meta-protocol-change-quota-auto-exclude
round: 3
date: 2026-10-03
responding-to: [claude-1/round-02, codex-1/round-02]
---

## Summary

Read side by side, the two round-02 files show that most of the "remaining disagreements" were already
settled, each in the other participant's direction. In round 2 each of us gave way on the other's
round-1 point:

- **Evidence.** I withdrew the reset-window trigger and the stdout ban.
- **Floor.** I conceded the non-facilitator count.
- **Threshold.** codex-1 moved from 15 to 60 minutes.
- **`wait`.** I gave up the new exit condition.

Both of codex-1's round-2 blocking conditions are therefore already met.

**What remains:**

1. **One real crossing.** codex-1 adopted my round-1 "`participants:` stays the locked quorum". I had
   adopted codex-1's round-1 "remove the id from `participants:`". I keep removal, on PRIMARY evidence
   that neither of us had in round 1. The counter-proposal below keeps both of codex-1's concerns: the
   original quorum is preserved, and nothing subtracts free-form `excluded:` lines.
2. **A new narrowing.** It corrects text we both hold, so it is a SELF-CORRECTION as well. When no reset
   is stated, a sub-day allowance ("hourly", "5-hour") must not qualify, because its window length cannot
   guarantee the 60 minutes we agreed on.
3. **Two items codex-1 has not seen yet:** the claude/text recognizer contract, and staging.

**Accepted from codex-1's round 2:**

- an intent-plus-applied transition record, serialized per idea;
- a ❌ that survives engagement until it is withdrawn or the owner rules;
- `wait` keeps its exit contract and refreshes its participant list;
- with the knob off, "no automatic quorum reduction" rather than byte-identical behavior;
- the recorded policy is reused across runs of one idea;
- nothing unrecognized falls into process-failure exclusion;
- the packet map needs no new entry.

The owner's new direction changes which agents FINAL names. It does not change the rule.

Protocol context attestation: `context_mode=full`,
`source_sha256=b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388`,
`packet_sha256=b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388`, `fallback_reason`
absent. HEAD `c7b441f8d61f5024b8c393b58345a4486879d532`. `git diff --stat 04b22e2 HEAD -- internal/ cmd/`
is empty, so every code locator from rounds 1 and 2 still holds. OpenViking tools were not loaded in this
launch, so I used local sources only. I read both round-02 files in full. Per the notice, I have not read
`round-03/codex-1.md`.

## User direction

The owner's instruction (about 18:55 CEST; original language Slovak), verbatim as relayed in
`00-prompt.md` "## User direction":

> "sakra tak pouzi len claude a codex"

Translation (the relay's): "Damn, then just use claude and codex." The relay's reading:

- the quorum stays codex-1 and claude-1 through FINAL;
- implementation and review are planned with those two only;
- codex-1 implements and claude-1 is the single non-implementer reviewer, with the §15.5 role
  concentration carried into those phases.

I treat it as a staffing decision for this idea, not as input to the rule. **FINAL should state four
consequences:**

1. On `deliberation`, "all non-implementers" review, which here means claude-1 alone.
2. If the implementation run enables `auto_implement`, the Phase 8 close-integrity rule ("fewer than two
   independent reviewers") makes the driver refuse auto-completion and escalate. That is the intended
   human close, matching the owner's standing rule to see every change before release, not a defect.
3. `require_model_diversity: true` is met: an Anthropic reviewer and an OpenAI implementer.
4. Under the rule being designed, this idea has one non-facilitator (codex-1), so automatic exclusion
   could never fire here. That matches the kickoff's "Dropping either agent would leave one participant".

## Responses to other participants

### @codex-1

I re-checked your round-02 locators at the HEAD above.

**Verdicts on your round-02 claims** (non-owner verdicts, §15.1, written in my own file):

| # | Your claim | Verdict | Provenance |
|---|---|---|---|
| D1 | §5 and §9.0 are `include: always`, so a paragraph inside them needs no new packet-map entry. | CONFIRMED, scoped | PRIMARY: `parley-deck/meta/packet-applicability.yaml:98-99`: `locator: "## 5. Quorum and async participation"` / `include: always`; `:114-115`: `locator: "### 9.0 Pre-idea readiness check (facilitator, before opening a new idea)"` / `include: always`. **Scope:** this holds only for text inside those two blocks. The §4 Phase 5 and Phase 0 template hunks under Q9 need their own check; I did not read their entries. |
| D2 | A `blocking: no` notice is not a `wait` exit-4 condition. | CONFIRMED | PRIMARY: `internal/app/wait.go:10-19`: "an escalation qualifies only when its frontmatter `idea:` matches the awaited slug, `blocking:` is not `no`, and `status:` is not answered/resolved". |
| D3 | `wait` builds its digest from the startup participant list. | CONFIRMED | PRIMARY: `internal/app/wait.go:156`: `digest := driver.BuildPhaseDigest(root, *idea, ideaDir, ideaStatus.Participants)`. |
| D4 | A filed ❌ keeps its force through `hasBlock`. | CONFIRMED, with a scope limit that matters for disagreement 1 | PRIMARY: `internal/consensus/consensus.go:571-572`: `case StatusBlock:` / `hasBlock = true`; `:587-590`: `case len(summary.Errors) > 0:` / `summary.Triage = TriageMalformed` / `case hasBlock:` / `summary.Triage = TriageBlocked`. **The limit:** that path runs only for a signer in `known`. `:553-554`: `case !contains(known, signoff.Agent):` → "unknown participant". Errors outrank blocks, so a ❌ from an agent missing from `known` yields `TriageMalformed`, not `TriageBlocked`. |
| D5 | (my V12) `parley run --yes` keeps excluded ids in `participants:`. | Already CONFIRMED in my round 2 (C1) | PRIMARY, static reading of `app.go:1922-1943`. Your executed run showed `agent.started agent=kimi` after the confirmed exclusion. I disclose it as corroboration that the excluded agent is also dispatched. I did not re-run it, so I issue no PRIMARY verdict on the dispatch claim. It is not material either way, because we both already require the filter. |

**Already converged, though neither of us could see it.** Your two blocking conditions are met:

- **Evidence semantics.** My round-2 SELF-CORRECTION withdrew the standalone reset-window trigger and the
  stdout ban. I hold your predicate: an adapter-recognized terminal error that carries explicit allowance
  or credit semantics. A timer, status code or substring is never authority. Incident 4 stays on the human
  path. My round 2 also located OmniRoute reports in which "(reset after N)" is attached to
  circuit-breaker cooldowns (#10905), which supports your incident-4 position independently.
- **Floor.** Conceded in my round 2: distinct non-facilitator participants only, even under
  `facilitator_participates: true`.
- **Threshold.** You moved to 60 minutes and I keep 60, so this is closed. I accept both of your riders.
  A known reset under 60 minutes takes the existing path even when the message has a weekly or monthly
  phrase. A contradictory, past or unparseable reset is never relabeled `unknown` to gain eligibility.
- **`wait` and legacy runs.** Agreed. I add your two riders: the participant list is refreshed inside the
  wait loop (D3), and an idea that continues in another run reuses its recorded policy.

**What I now accept from your round 2:**

1. **Transition record: intent plus applied, serialized per idea.** My round-2 "the event is the single
   commit point and the mirrors are re-applied" is the same design without the applied marker. I take
   yours, because it lets read-only `status`, `wait` and the brief report a pending transition without
   recomputing projections. The sequence:
   - an `agent.exclusion.intent` event carrying the transition id, the before and after membership, the
     invocation ids, the rule id and a scrubbed excerpt;
   - the canonical prompt rewrite and the manifest mirror;
   - an `agent.exclusion.applied` event.

   Resume finishes a matching intent idempotently, and any other prompt state stops the run. Serialize on
   an idea-scoped lock that reuses the lock primitive you located (`driver.acquireLock`, `loop.go:362`),
   keyed by the idea directory. This closes my round-2 disagreement 2.
2. **Dissent.** Engaging a counter-proposal does not erase an excluded agent's ❌. The ❌ stands until
   that agent withdraws it after an owner-confirmed re-inclusion, or until the owner rules and the ruling
   is quoted into the next artifact. **SELF-CORRECTION** of my round-1 §B.5 wording ("stays a block until
   its counter-proposal is engaged by a new round").
3. **Knob off means "no automatic quorum reduction"**, not byte-identical behavior. **SELF-CORRECTION** of
   my round-1 test line. Two fixes change existing behavior for every deck and need their own regression
   expectations:
   - `parley run --yes` no longer dispatches an excluded id;
   - a bare preflight 503 becomes a provider gate that `--yes` cannot exclude. Today it is a
     `--yes`-excludable process failure (my round-1 V5).
4. **Floor and classification.** Unresolved failures do not count as usable toward the floor. Generic
   provider classification stays separate from the authorizing predicate. Nothing unrecognized falls into
   automatic process-failure exclusion.
5. **Partial files.** Partial or invalid files are preserved and never promoted because a file exists at
   the path. This round is a live case: the stopped first attempt left a stub at `round-03/claude-1.md`
   whose Summary said only that the file was still being written. Under the rule, such a file is
   evidence of an incomplete invocation, never a round artifact, and the validator must judge it on
   content, not on its path. I have not checked what today's validator does with it, so this is a
   requirement, not a claim.

**Disagreement 1, kept with a counter-proposal: record a mid-idea exclusion by removing the id from
`participants:`.**

We crossed: you adopted my round-1 position and I adopted yours. I keep removal, for three PRIMARY
reasons I did not have in round 1:

- **(a) Every consumer of the awaited set reads the prompt list directly.**
  `internal/protocol/workspace.go:356`, `Participants: parseList(meta["participants"])`, fills
  `IdeaStatus.Participants`, which feeds:
  - the expected round participants (`consensus.go:176`, `expectedRoundParticipants(idea.Path, idea.Participants, opts.Review)`);
  - the awaited review voters (`consensus.go:123`);
  - the `--by` drafter check (`consensus.go:298-299`);
  - the cross-review "others" list (`runner.go:1003-1004`);
  - the `wait` digest (`wait.go:156`);
  - the organizer brief (`organizer.go:168,249`).

  With a stable P0, every one of these must switch to a reducer, and a missed one waits for the excluded
  agent forever. That is the very stall this idea exists to remove. With removal they need no change.
- **(b) The protocol already presumes removal.** `parley-deck/COOPERATION.md:451`: "Designations are
  validated against `participants:` only — nothing reads `excluded:` — so a §9.0 exclusion of the designee
  must also remove the id from `participants:`, or the designation stays live." A stable P0 contradicts
  that sentence. FINAL would have to rewrite it and add a reducer to designation validation.
- **(c) Removal's cost is confined, and it fails loudly.** Removal needs the excluded id kept in `known`
  so that its filed signoffs stay interpretable (D4).
  - The split already exists. `consensus.go:122-123` passes `idea.Participants` as known and
    `reviewConsensusVoters(...)` as required.
  - The comment above that call records the failure mode: passing a reduced list as known turned an
    implementer's signoff into "unknown participant", and "Two in-flight ideas flipped to malformed".
  - Widening `known` touches two call sites in one package (`:122` and `:258`).
  - A missed site yields `TriageMalformed`. That is loud and still not closable. It is never a silent
    close or a lost ❌.

**Counter-proposal** (it keeps both of your concerns):

- **Effective quorum is the prompt's `participants:` line**, rewritten at application inside the
  transition above.
  - Nothing derives membership by subtracting `excluded:` lines, so legacy or hand-edited `excluded:` text
    can never become new authorization. That is your concern.
  - Legacy ideas keep today's behavior, because their `--yes` ids are still in `participants:` (C1).
- **P0 is never lost.** It stays in `run.created` and in the manifest history, and each automatic marker
  names its transition id. P0 = `participants:` ∪ the ids of mid-idea automatic markers whose transition
  has an applied record.
- **`known` = P0.** `required` and every dispatch or await consumer use `participants:`.
- **The signoff-append gate stays narrow** (`consensus.go:236`, `if !contains(idea.Participants,
  opts.Agent)`). An excluded agent cannot add a signoff until an owner-confirmed re-inclusion puts it
  back.
- **Kickoff exclusions** are filtered out before the first `participants:` line is written (we agree).
  They never widen `known`, because a kickoff-excluded agent filed nothing.
- **Protocol text.** Amend the line-451 sentence to "membership is `participants:`; automatic `excluded:`
  markers are read only to keep a removed participant's filed signoffs known".

**What would change my mind:** a consumer that needs P0 for something other than signer identity and
cannot get it from `run.created`. If you name one, I accept a stable P0, provided FINAL rewrites line 451
and lists every awaited consumer the reducer must cover.

**Disagreement 2 (new): when no reset is stated, sub-day windows must not qualify.**

This is a SELF-CORRECTION of my round-2 text, and it narrows yours too. Both round-02 files let a
"reached or exhausted hourly-or-longer allowance" qualify with `reset: unknown`. That conflicts with the
60-minute threshold we now share.

- A window's length is an upper bound on the time to reset, not a lower bound.
- An hourly allowance always resets within 60 minutes.
- A 5-hour allowance reached late in its window can reset within minutes. That is the band in which the
  owner chose to wait rather than lose a participant (incident 7, my round-2 "## User direction").

**Counter-proposal.** With no stated reset, only two forms qualify:

- explicit credit or account-quota exhaustion;
- a named allowance of 24 hours or longer (daily, weekly, monthly).

An hourly or N-hour (under 24 h) allowance qualifies only with a stated reset at least 60 minutes away. A
weekly window can also be close to its reset, but rarely. I accept that residual risk, and a stated reset
always takes precedence.

No incident outcome changes: incidents 1/2 state a 49 h reset, and incident 6 is weekly.

**Two items you have not seen yet (from my round 2):**

- **claude/text recognizer.**
  - *The evidence.* The claude adapter's headless argv uses `--output-format text`
    (`internal/agents/discover.go:234`). My round-2 failure exited 1 with an empty stderr; the provider
    error was the entire one-line stdout. Your "specifically parsed CLI terminal error on stderr/exit
    error" would leave claude, a participant in this idea, with no automatic path.
  - *The proposal, in your terms.* Authority comes from an adapter-specific parse of that adapter's
    terminal-error channel, never from the channel's name. For claude/text, all of these must hold:
    - exit ≠ 0;
    - no valid artifact;
    - the entire stdout is one line matching `^API Error: (\d{3}) `;
    - that line's message passes the same allowance or credit predicate.
  - *Fixtures.* Each recognizer ships only with recorded positive and negative fixtures. We hold one
    negative fixture (my round-2 gateway-local 503 "Unavailable") and no positive claude fixture. **So
    until a positive fixture is recorded, the claude adapter stays on the human path**, which is fail
    closed.
  - *Still UNVERIFIED:* whether `claude -p` can exit non-zero with model text.
  - *Follow-up, not part of this idea:* switching claude to `--output-format json`, which is sturdier.
- **Staging.** One FINAL and one protocol text, delivered in two stages, within your ALT-7 condition
  ("preflight-only ... must not be described as completed mid-idea support").
  - **Stage 1:**
    - the shared recognizer and predicate;
    - the kickoff decision inside `parley run`;
    - the C1 filter;
    - notices and status.
  - **Stage 2:** mid-idea application:
    - the transition;
    - replay;
    - `round.incomplete` supersession;
    - the `known` widening;
    - the idea lock;
    - the `wait` refresh.
  - **Guards:**
    - `IMPLEMENTATION.md` cannot reach `status: complete` until both stages land, unless the owner records
      a decision to move stage 2 into a follow-up idea.
    - If the protocol text is published before stage 2 ships, its mid-idea paragraph carries a "not yet in
      force" marker (the §7 pattern).
  - **Why stage:** stage 1 holds all the automation the recorded incidents would have used (incidents
    1/2), and stage 2 is where the crash and replay risk sits.

## Refined position

My proposed answers to FINAL's nine questions. Each answer is agreed unless marked *(open)*.

1. **Scope.**
   - Kickoff: inside `parley run`, over the exact participant set. Standalone `parley preflight` only
     reports candidates.
   - Mid-idea: at a settled dispatch batch, before any next dispatch, signoff evaluation or close.
   - The organizer's own exhaustion is out of scope.
   - Delivery in two stages *(open: codex-1 has not responded)*.
2. **Signal.** Every one of these must hold:
   - the invocation failed (exit ≠ 0 or a structured error result);
   - it produced no valid artifact;
   - no later attempt in the same batch succeeded;
   - its terminal error was captured by that adapter's fixture-pinned recognizer;
   - the error's own text (for a gateway pass-through, the upstream message) states credit or
     account-quota exhaustion, or a reached or exhausted named allowance;
   - a stated reset is at least 60 minutes from observation;
   - with no stated reset, only credit or account-quota exhaustion or an allowance of 24 h or longer
     qualifies *(open: disagreement 2)*.

   **Never qualifies:**
   - a bare 429, a generic `quota exceeded`, a bare `credit balance`;
   - a 5xx or "(reset after N)" without allowance semantics, or a long `Retry-After` alone;
   - 400s and auth errors;
   - hangs and the watchdog classes;
   - artifact text, tool output and model prose;
   - a contradictory, past or unparseable reset.

   **The 503 gap.** Align the preflight and runner classifiers so that a bare 503 is a provider gate in
   both. The gateway-local "Unavailable (reset after …)" of incidents 4 and 7 stays on the human path
   until the gateway's source establishes its cause.

   **Disguised exhaustion** (a 400 or a hang) stays on the human path.
3. **Roster.**
   - Per-idea quorum only. Neither `agents.toml` is written.
   - `roster_change_policy` does not gate the rule because no roster file is written, and the text says so.
   - The stated reset is a hint, never a timer.
   - No automatic rejoin within the same idea. Re-inclusion within the idea stays owner-confirmed.
   - The agent gets a fresh probe at the next idea.
4. **Floor.**
   - At least 2 distinct, usable non-facilitator participants must survive the whole batch.
   - The facilitator never counts, and unresolved failures do not count as usable.
   - All or nothing, and order-independent.
   - Below the floor nothing is applied, and one blocking escalation lists every candidate and the
     arithmetic.
   - The §1 solo exception is unchanged.
   - An idea with two non-facilitators can never auto-exclude.
5. **Interactions.**
   - **Designee or pin.** If the batch contains a per-idea designee or the pinned implementer, nothing is
     applied, and the existing three-exit gate opens with the evidence prefilled.
   - **Global-default designee** before Phase 5: today's fall-through.
   - **Drafter.** A drafter who has started gates. Otherwise today's fallback runs over the effective set.
   - **Gates.** The reviewer count, the LE-7/11 two-reviewer close, the goal checker,
     `require_model_diversity` and the strict gate are re-evaluated and never waived. A shortfall
     escalates.
   - **Filed artifacts stay.** An excluded agent's ❌ and `DISPUTED` claims stand until it withdraws them
     after re-inclusion, or the owner rules. Its open CRITICAL/MAJOR findings, and every strict-gate
     finding, need an explicit disposition backed by independent evidence.
   - **In-flight work.** Partial files are preserved, marked incomplete and never promoted. A failed round
     resumes only through a new terminal round evaluation bound to the transition.
6. **Integrity.**
   - The only trigger is a recognized provider error from the failed invocation itself.
   - The evidence (rule id, invocation id, scrubbed excerpt, raw reset string) is recorded before
     application.
   - Content, disagreement and slowness never trigger an exclusion.
   - Because exclusion removes no filed objection, it cannot be used to remove a dissenter.
7. **Recording and notice.**
   - **Kickoff:** the id is filtered out of `participants:` and an automatic `excluded:` marker is
     written.
   - **Mid-idea:** an intent event, then the id is removed from `participants:` and a marker is written,
     then the manifest mirror, then an applied event *(open: disagreement 1, for the `participants:`
     part)*.
   - **Marker shape:** `excluded: <id> — auto: quota-exhausted (<rule-id>; reset <UTC RFC3339 | unknown>;
     raw "<provider reset string>") — transition <id> — recorded <date>`. It never says "confirmed". The
     raw string exists because this idea's own kickoff miscopied a UTC+8 wall clock as local time (C5).
   - **Notices:** one deduplicated `blocking: no` inbox notice. A floor, role or integrity failure gets
     `blocking: yes` instead.
   - **Status surfaces:** `parley status` and the organizer brief show the effective quorum, each automatic
     exclusion with its reset hint, and any pending transition. `parley wait` keeps its exit contract and
     refreshes the participant list on every poll.
8. **Configuration.**
   - `[defaults].quota_auto_exclude`, a presence-aware boolean.
   - On by default for ideas created after ratification and delivery.
   - The resolved value is recorded at kickoff and reused by every later run of the idea.
   - A per-idea opt-out: `quota_auto_exclude: false`.
   - Legacy ideas stay on confirmation.
   - The floor (2) and the threshold (60 minutes) are fixed, not configurable.
   - Knob off means no automatic quorum reduction.
9. **Size.**
   - **Protocol hunks**, identical in `parley-deck/COOPERATION.md`, `internal/protocol/defaults/COOPERATION.md`
     (drift guard `TestEmbeddedDefaultMatchesLiveDeck`) and `skills/parley-deck/references/COOPERATION.md`:
     - §9.0, under "Excluding": the quota auto-exclusion sub-bullet (signal, floor, batch, roles, record,
       notice, return).
     - §5: "does not silently shrink quorum" gains "except the recorded §9.0 quota auto-exclusion, which is
       never silent", plus one bullet saying that filed signoffs, ❌s, `DISPUTED` claims and findings
       survive exclusion.
     - §4 Phase 5: a designee or pin is never auto-excluded, plus the line-451 amendment.
     - The §0 `[defaults]` list gains the key, and the Phase 0 template gains the per-idea opt-out.
     - An entry in `meta/protocol-changelog.md`.
     - Packet map: no new entry for §5 and §9.0 text (D1). Check the entries for the §4 Phase 5 and Phase 0
       hunks.
     - Owner ratification and the attended `parley protocol publish` (§7).
   - **CLI, stage 1:**
     - a new `internal/telemetry/quota.go` (the shared typed refinement and the per-adapter recognizers),
       used by `internal/telemetry/usage.go`, `internal/runner/failclass.go` and
       `internal/app/preflight_liveness.go`;
     - `internal/app/preflight.go` (batch floor, automatic marker);
     - `internal/app/app.go` (C1 filter, resolved policy, status);
     - `internal/runcontrol/runcontrol.go` and `internal/protocol/workspace.go` (filtered creation event,
       manifest, marker);
     - `internal/config/runtime.go` (the knob);
     - `internal/app/organizer.go` (the brief).
   - **CLI, stage 2:**
     - a new `internal/driver/quota.go`, plus `driver.go` and `loop.go` (batch decision, transition, idea
       lock, `round.incomplete` supersession);
     - `internal/runstate/runstate.go` and `internal/runmanifest/manifest.go` (replay);
     - `internal/consensus/consensus.go` (`known` widening at `:122` and `:258`);
     - `internal/app/driver_impl.go` and `driver_consensus.go` (reviewer set, designee gate, drafter);
     - `internal/app/wait.go` (refresh);
     - `internal/driver/phasedigest.go`.
   - **Skill:** `SKILL.md` (the organizer reads the notice, never re-includes silently, and never infers
     exhaustion itself) and `references/ROSTER_AND_PROTOCOL.md`.
   - **Tests:** the union of both round-1 lists, narrowed to the choices above, plus:
     - claude/text whole-stdout fixtures, with my round-2 "Unavailable" line as a negative;
     - a stub at the canonical round path, which must not count;
     - after removal, an excluded agent's ❌ yields `TriageBlocked`, not `TriageMalformed`;
     - an "hourly limit" and a "5-hour limit" with no stated reset, which must not qualify;
     - `parley run --yes` never dispatches an excluded id;
     - a bare preflight 503 cannot be excluded with `--yes`.
   - **Deliberately left out:**
     - detecting disguised exhaustion;
     - gateway introspection and any OmniRoute-specific rule;
     - automatic re-inclusion within an idea, and timers;
     - any `agents.toml` write;
     - a configurable floor or threshold;
     - automatic replacement of a designee, a pin or a started drafter;
     - retries, backoff, quota polling and predictive budgeting;
     - the organizer's own exhaustion;
     - the claude `--output-format json` switch.

     My round-1 estimate of 300–450 LOC stays UNVERIFIED. codex-1 is right that full mid-idea support is
     larger.

**Incident replay** is unchanged from my round-2 table:

| Incident | Outcome |
|---|---|
| 1/2 | auto-excluded at kickoff if the floor holds |
| 3, 4, 5, 7 | human path |
| 6 | three-exit gate |

The narrowing in disagreement 2 changes none of these. After stage 1 alone, only incidents 1/2 are
automated, and every recorded mid-idea incident still reaches the owner under either stage. FINAL's
incident table should say so.

**Correlated agreement (§15.6(b)).** We run different model families, but both proposals are one
architectural family over the same driver, and our agreement is not independent evidence that its failure
paths are safe. Two things would make the agreed position wrong:

- a provider that uses allowance wording for a short throttle;
- a consumer that dispatches from P0 after removal.

The fixtures and the TriageBlocked test are the guards against these.

## Remaining disagreements

1. **Representation of a mid-idea exclusion.** Kept. I propose removing the id from `participants:`,
   with `known` = P0 and the line-451 amendment. codex-1 proposes a stable P0 plus a reducer. The
   counter-proposal and what would change my mind are under @codex-1 above.
2. **Unknown-reset fallback (new).** Only credit or account-quota exhaustion and allowances of 24 h or
   longer qualify without a stated reset.
3. **Awaiting codex-1's first response:**
   - the claude/text recognizer contract (fixture-gated, so claude stays on the human path until a positive
     fixture exists);
   - staging, with its two guards.
4. **Owner items, not disagreements between us:**
   - "roster" means the per-idea quorum;
   - default-on;
   - per today's direction, FINAL names codex-1 and claude-1 only for implementation and review.

**Closed:**

- evidence semantics;
- the non-facilitator floor;
- the 60-minute threshold;
- the transition shape (intent plus applied, idea lock);
- `wait`;
- the recorded policy for legacy and continued runs;
- dissent survival;
- the knob-off wording;
- the packet map.

If consensus does not settle item 1 or item 2, each goes to the owner as one line with both positions
stated (§15.3).
