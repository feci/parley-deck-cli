---
idea: meta-protocol-change-quota-auto-exclude
drafted-by: claude-1
date: 2026-10-03
---

Role concentration (§15.5): the claude-1 participant drafts this file and is also the idea's declared facilitator
(`facilitator_participates: true`), run as a separate process from the claude-1 organizer, which writes no
consensus content. The frontmatter says `drafted-by: claude-1`, replacing the driver scaffold's `user`.

Protocol context attestation: `context_mode=full`,
`source_sha256=b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388`,
`packet_sha256=b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388`, `fallback_reason` absent.
I read every round file: `round-01/`, `round-02/` and `round-03/`, each from both claude-1 and codex-1. The two round-03 files were
written in parallel, so neither answers the other. Below, each new round-03 proposal is marked as agreed
or open. HEAD `22cfddffed5c827e5aa28dfac9e7a82409ebd7cd`. I checked new locators myself; each is tagged
where it is used.

## User direction

1. The owner's instruction, about 18:55 CEST (original language Slovak), quoted verbatim in both round-03
   files:

   > "sakra tak pouzi len claude a codex"

   Translation (the relay's): "Damn, then just use claude and codex." The relay read it this way: the
   quorum stays codex-1 and claude-1 through FINAL, and implementation and review are planned with those
   two only.

2. The owner's answer to the organizer's note
   `inbox/claude-1-to-user_meta-protocol-change-quota-auto-exclude_codex-auth.md` (original language
   Slovak). That note reported that the gateway had rejected codex-1 with "401 Unauthorized … invalidated
   oauth token" after codex-1's round-03 file was complete. The answer, verbatim:

   > "pokracuj"

   Translation: "Continue." Source: `inbox/user-to-claude-1_meta-protocol-change-quota-auto-exclude_codex-auth-answer.md`.
   The relay re-checked codex-1 at 20:23 CEST and it replied `PONG`. The organizer's own probe at 19:57
   CEST had also returned `PONG`.

## Agreed decisions

One rule, structured by the nine design questions in `00-prompt.md`. Every item is agreed by both round-03
files unless it is marked **OPEN** (both positions stated) or **drafter proposal** (new in this draft;
codex-1's signoff accepts or blocks it).

**Staffing for this idea (owner direction 1).** FINAL names only codex-1 and claude-1:

- codex-1 implements; the separately launched claude-1 participant is the single non-implementer reviewer
  (deliberation: "all non-implementers"). The claude-1 organizer stays organizer only, and the §15.5 role
  concentration is disclosed in each phase.
- `require_model_diversity: true` is met: an OpenAI implementer and an Anthropic reviewer.
- Phases 5–8 run **attended**. Under `auto_implement` the driver refuses auto-completion with "fewer than
  two independent reviewers" (`COOPERATION.md:692–697`, codex-1 r3). That gate is not lowered and two
  claude-1 processes never count as two reviewers. The close takes both participants' review-consensus
  signoffs and current-tree criterion evidence; `internal/driver/impl.go:283–288` already prints the manual
  route.
- Under this rule, this idea has one non-facilitator (codex-1), so automatic exclusion could never fire here.

**Q1 Scope.**
- **Kickoff.** The decision is made inside `parley run`, over the exact proposed participant set, after all
  probes return. Standalone `parley preflight` only reports candidates and evidence.
- **Mid-idea.** The decision is made at a settled participant-dispatch batch, after every affected writer
  has stopped, and before any next dispatch, signoff evaluation or close.
- The organizer's own exhaustion is out of scope. No organizer failover and no replacement agent.
- **Delivery in two stages, specified in one FINAL.** Stage 1 is kickoff. Stage 2 is mid-idea.
  - Until stage 2 ships, the published text marks mid-idea application "not yet in force" (the §7 pattern).
  - `IMPLEMENTATION.md` cannot reach `status: complete` on stage 1 alone. If the owner chooses stage 1
    only, that narrower authorization is recorded together with a linked stage-2 follow-up idea.

**Q2 What counts as exhaustion (fail closed).** Every one of these must hold:
1. The invocation failed (exit ≠ 0 or a structured error result).
2. It produced no valid completed artifact, and no later attempt in the same batch succeeded. A valid
   artifact or a later success always wins.
3. Its terminal error was captured by that adapter's recognizer, and that recognizer's **native-error
   provenance is established**: located CLI behavior or source, plus recorded positive fixtures and
   adversarial negative ones. The negatives cover an assistant quoting `API Error: 429 … Limit Exhausted`,
   a tool emitting it, and each of these followed by a failure with no valid artifact.
   - An adapter without established provenance is listed as **diagnostic-only, automatic exclusion
     unsupported**, and stays on the existing human path.
   - This applies to every adapter, including claude/text and codex. An adapter name is not evidence.
   - `is_error` alone is not enough. The structured error must identify a terminal provider failure, not a
     local budget or tool failure.
4. The error's own text (for a gateway pass-through, the upstream message) states one of these:
   - explicit account-credit or account-quota exhaustion;
   - or an explicitly **reached** or exhausted named usage allowance. A message that only names a policy
     does not count.
5. A stated reset is at least **60 minutes** from observation.
   - A known reset under 60 minutes takes the existing path, even when the message has a weekly or monthly
     phrase.
   - A contradictory, past or unparseable reset gates. It is never relabeled `unknown` to gain eligibility.
6. **No stated reset: OPEN.**
   - *codex-1 (r3 item 2):* explicit account credits or quota, or a reached hourly-or-longer allowance,
     qualifies.
   - *claude-1 (r3 disagreement 2):* only credit or account-quota exhaustion, or a named allowance of
     24 hours or longer (daily, weekly, monthly), qualifies. An hourly or N-hour allowance needs a stated
     reset at least 60 minutes away, because a window's length bounds the time to reset only from above.
     An hourly allowance always resets within 60 minutes, so qualifying it without a reset defeats the
     threshold both participants hold.
   - *Drafter proposal:* adopt the narrowing. No incident outcome changes.
   - codex-1 has not seen it. If codex-1's signoff objects, round 4 opens on this item. If it is still
     unresolved after that, it goes to the owner with both positions (§15.3).

**Never qualifies:**
- a bare 429, a generic `quota exceeded`, a bare `credit balance`;
- a 5xx or a "(reset after N)" without allowance semantics, or a long `Retry-After` alone;
- 400s and auth errors;
- hangs and every watchdog class;
- artifact text, tool output and model prose.

No generic `reset after ⇒ quota` regex. Gateway-pool phrases ("Unavailable", "cooling down", "all accounts
exhausted") are not positive forms until the gateway's source establishes their cause.

**The 503 gap.** Align the preflight and runner classifiers so a bare 503 is a provider gate in both, not a
`--yes`-excludable process failure. Generic provider classification stays separate from the authorizing
predicate, and nothing unrecognized falls into automatic process-failure exclusion. **Disguised exhaustion**
(a 400 or a hang) stays on the human path.

**Q3 What "roster" means.**
- Per-idea quorum only. Neither `agents.toml` is written.
- `roster_change_policy` does not gate the rule, because no roster file is written. The text states both
  facts.
- The stated reset is a hint, never a timer. There is no same-idea automatic rejoin and no quota polling.
- Re-inclusion within the idea stays owner-confirmed, under the existing catch-up rules. The agent gets a
  fresh probe at the next idea.
- A notice may suggest one owner-authorized relaunch after the stated reset plus 5 minutes.
  - The time is labeled a provider estimate.
  - There is no suggestion when the reset is unknown, and no retry worker.

**Q4 The minimum of 2.**
- At least 2 distinct, usable **non-facilitator** participants must survive the whole batch.
  - The facilitator never counts, even under `facilitator_participates: true`.
  - Unresolved failures do not count as usable, and duplicate ids do not inflate the count.
- All or nothing, evaluated once per batch, independent of order.
- Below the floor, nothing is applied and one blocking escalation lists every candidate and the
  arithmetic. The §1 solo exception is unchanged.
- The floor is fixed at 2. An idea with two non-facilitators can never auto-exclude.

**Q5 Interactions.**
- **Designee or pin.** If the batch holds a per-idea designee or the pinned implementer, nothing in the batch
  is applied. The existing three-exit gate opens, with the evidence prefilled.
- **Global-default designee** before a pin exists: today's fall-through.
- **Drafter.** A candidate who has started a canonical consensus or FINAL draft blocks automatic application.
  Otherwise today's drafter fallback runs over the current set.
- **Gates are re-evaluated and never waived:** reviewer count per track, the LE-7/11 two-reviewer close,
  independent goal-checker eligibility, `require_model_diversity` and `strict_gate`. A shortfall escalates.
- **Filed artifacts survive.**
  - An excluded agent's ❌ stands. Engaging its counter-proposal does not erase it. It stands until that
    agent withdraws it after an owner-confirmed re-inclusion, or the owner rules and the ruling is quoted
    into the next artifact.
  - Its `DISPUTED` claims stand (§15.3: never resolved by count).
  - Its open CRITICAL and MAJOR findings, and every strict-gate finding, need an explicit disposition
    backed by independent evidence. "Author excluded" is never a disposition.
- **In-flight work.**
  - Partial or invalid files are preserved, identified as incomplete, and never promoted because a file
    exists at the path.
  - A failed round resumes only through a new terminal round evaluation bound to the transition, after
    every survivor artifact is validated. The earlier `round.incomplete` and the failed invocation are kept.
  - Closed FINAL and IMPLEMENTATION artifacts stay frozen.

**Q6 Integrity.**
- The only trigger is a recognized provider error from the failed invocation itself.
- The evidence is recorded before application: rule id, invocation id, scrubbed excerpt, raw reset string,
  UTC observation and reset times.
- Content, disagreement and slowness never trigger an exclusion.
- Exclusion removes no filed objection, so it cannot be used to remove a dissenter.

**Q7 Recording and notice.**
- **Membership.** The crossed round-2 concessions are reconciled. Both round-03 files hold this design.
  - `00-prompt.md` `participants:` is the **current** set.
  - Kickoff exclusions are filtered out before the first `participants:` line, the creation event, the
    manifest and the first dispatch are written.
  - A mid-idea exclusion rewrites `participants:` inside the transition.
  - The kickoff quorum and every authorized revision are kept in immutable history. `run.created` and
    participant artifacts are never rewritten.
  - `known` signers are the identities that were members under that recorded history. `required` signers,
    and every dispatch and await consumer, use the current `participants:`.
  - Kickoff-excluded ids never become known. Nothing derives membership by subtracting `excluded:` lines.
  - The signoff-append gate stays narrow (`consensus.go:236`).
  - The line-451 sentence is amended to say this.
- **Transition.**
  - One durable **batch-level** commit record, never one per agent. It carries:
    - the idea and run identity and the prior membership revision;
    - the before and after sets;
    - every candidate invocation id, the resolved policy and scope, and the rule ids;
    - the observation and reset times and the scrubbed evidence references.
  - The whole batch is validated under idea-scoped serialization before commit. Projections are reconciled
    before any further dispatch, signoff evaluation or close: the prompt, the manifest, captured consumers
    and the failed-round evaluation.
  - A committed but unreconciled transition is a recoverable pending state. Read-only `status`, `wait` and
    the organizer brief report it and never repair it.
  - A new run for the same idea finds the record before it acts. Missing or contradictory history blocks.
  - The commit needs checked durability and a fail-closed path for a truncated record.
  - One driving run per idea. A competing run under another run id is rejected or serialized.
- **Marker.** It carries an automatic marker, the rule id, the reset in UTC RFC3339 or `unknown`, the
  scrubbed raw provider reset string, the transition id and the recorded date. It never says "confirmed".
- **Notices.** One `blocking: no` owner notice per transition, deduplicated by transition id. It names who
  was excluded, the decisive reason, the reset hint, the survivors and any remaining gates. A floor, role
  or integrity failure gets `blocking: yes` instead.
- **Surfaces.**
  - `status`, `wait` and the organizer brief use the same current set and list each automatic exclusion
    with its reset hint, plus any pending transition.
  - `wait` keeps its exit contract: a successful reduction is not an exit-4 condition. It refreshes its
    participant list on every poll.

**Q8 Configuration.**
- `[defaults].quota_auto_exclude`, a presence-aware boolean, with a deck override and a per-idea
  `quota_auto_exclude: false` in `00-prompt.md`.
- On by default only for ideas created after ratification and delivery. Legacy ideas stay on confirmation.
- The resolved policy **and the authorized scope** (kickoff-only, or kickoff plus mid-idea) are recorded at
  kickoff and reused by every later run of the idea.
  - A binary upgrade or a resume never widens scope. Widening goes through the existing owner-confirmed
    path.
  - A malformed config or an ambiguous record fails closed.
- The floor (2) and the threshold (60 minutes) are fixed constants.
- The knob set to off means "no automatic quorum reduction", not byte-identical behavior. Two fixes change
  existing behavior for every deck and carry their own regression expectations:
  - `parley run --yes` no longer keeps or dispatches an excluded id (C1);
  - a bare preflight 503 can no longer be excluded with `--yes`.

**Q9 Size.** FINAL lists the following.
- **Protocol hunks.** Identical in `parley-deck/COOPERATION.md`, `internal/protocol/defaults/COOPERATION.md`
  (drift guard `TestEmbeddedDefaultMatchesLiveDeck`) and the skill's
  `skills/parley-deck/references/COOPERATION.md`:
  - **§9.0 "Excluding".** The binding quota auto-exclusion sub-bullet, with the signal, provenance,
    threshold, floor and batch, role guards, record, notice, return and policy.
  - **§5.** "does not silently shrink quorum" gains "except the recorded §9.0 quota auto-exclusion, which is
    never silent". One bullet says that filed signoffs, ❌s, `DISPUTED` claims and findings survive, with
    known versus required signers.
  - **§4 Phase 5.** A cross-reference saying a designee or pin is never auto-excluded, plus the line-451
    amendment.
  - **§0 `[defaults]` list** gains the key. The **Phase 0 template** gains the per-idea opt-out.
  - **One-line cross-references,** only where they are needed: Phase 3/6/7 for historical signers and
    review gates, and the §9 checklist for notices.
  - **`meta/protocol-changelog.md`** gets an entry.
  - **Packet map.** No new entry for text inside §5 and §9.0, which are `include: always` (D1, both
    participants). The Phase 0 and Phase 5 blocks are `include: when`, so the binding rule text belongs in
    §9.0, and packet rendering must be checked. Any new heading needs a fresh map review.
  - **Ratification.** Owner ratification and the attended `parley protocol publish` (§7).
- **CLI, stage 1:**
  - a new `internal/telemetry/quota.go` for the shared typed refinement and the per-adapter recognizers,
    with their support table;
  - callers `internal/telemetry/usage.go`, `internal/runner/failclass.go`, `internal/runner/telemetry.go`,
    the terminal-result paths, and `internal/app/preflight_liveness.go`;
  - `internal/app/preflight.go` and `internal/app/app.go` (C1 filter, resolved policy and scope, status);
  - `internal/config/runtime.go`;
  - `internal/protocol/workspace.go`, `internal/runcontrol/runcontrol.go`, and the manifest and runstate
    fields;
  - `internal/app/wait.go` and `internal/app/organizer.go`.
- **CLI, stage 2:**
  - a new `internal/driver/quota.go`, plus `driver.go`, `loop.go` and `phasedigest.go`;
  - `internal/runstate/runstate.go` and `internal/runmanifest/manifest.go` for replay;
  - `internal/consensus/consensus.go` for the history-backed `known` at `:122` and `:258`;
  - `internal/app/driver_consensus.go` and `driver_impl.go`;
  - store durability as needed;
  - runner `phase58.go` and `acp.go` evidence paths.
  - No change to `internal/agents/discover.go`.
- **Skill.** `SKILL.md` (the organizer reads the notice, never re-includes silently, and never infers
  exhaustion itself), `references/ROSTER_AND_PROTOCOL.md` and the bundled protocol copy.
- **Tests:**
  - native versus quoted or tool-emitted errors, and supported versus unsupported adapters;
  - short, unknown, contradictory and past resets, and "hourly" or "5-hour" with no reset (subject to the
    Q2 OPEN item);
  - a valid artifact or a later success wins;
  - every permutation of a simultaneous batch, 4→2 and 3→1, duplicates, and the facilitator floor;
  - protected roles and the started drafter;
  - kickoff filtering and replay, including that `parley run --yes` never dispatches an excluded id;
  - a bare preflight 503 that cannot be excluded with `--yes`;
  - a stage-1-to-stage-2 upgrade that does not widen scope;
  - interrupted and truncated commits, duplicate replay, and competing runs with different run ids;
  - a stub at the canonical round path that does not count;
  - after removal, an excluded agent's ❌ yields `TriageBlocked`, not `TriageMalformed`;
  - a kickoff-excluded id that never becomes a known signer;
  - retained findings and the reviewer, diversity and goal-check gates;
  - notices, `status`, `wait` and the brief agreeing.

  Run the affected Go packages, the drift and packet checks, and the skill's `test/installer.test.js` and
  `test/lean-organizer.test.js`. These are future checks. None was run in this design run.
- **Deliberately left out:**
  - detecting disguised exhaustion;
  - gateway introspection and any OmniRoute-specific rule;
  - timers, quota polling, retries, backoff and predictive budgeting;
  - billing APIs and purchases;
  - automatic same-idea re-inclusion;
  - any `agents.toml` write;
  - a configurable floor or threshold;
  - automatic replacement of a designee, a pin or a started drafter;
  - organizer self-exhaustion and failover;
  - the claude `--output-format json` switch, which is a scoped follow-up.

  The 300–450 LOC estimate stays UNVERIFIED. Full mid-idea support is larger than the owner's "small".

**Incident table for FINAL.** It separates "the message would qualify" from "a reduction would execute".
Row 8 is new in this draft; the reading is mine, PRIMARY.

| # | Incident | Outcome under this rule |
|---|---|---|
| 1/2 | zcode-1 429 "Weekly/Monthly Limit Exhausted", reset +49 h (kickoff) | Message qualifies only through a provenance-verified zcode recognizer. At the actual kickoff kimi-1 was unresolved, so codex-1 was the only usable non-facilitator survivor. The floor fails and the owner is asked. |
| 3 | kimi-1 probe hang | human path |
| 4 | codex-1 503 "Unavailable (reset after 5h 51m 11s)" | human path (gateway-local text; an organizer process) |
| 5 | kimi-1 400 from a disabled connection | human path |
| 6 | kimi-1 403 "weekly usage limit" during an implementation attempt | Three-exit gate if kimi-1 was the designee or pinned implementer. Otherwise it needs raw terminal evidence from a verified recognizer, plus the floor. The excerpt is a summary, so no automatic outcome is asserted. |
| 7 | claude-1 503 "Unavailable (reset after 55m 29s)", this idea, round 2 | human path (gateway-local text; under 60 min; claude/text unsupported; floor) |
| 8 | codex-1 401 "invalidated oauth token … (reset after 12s)", this idea, round 3 | never a candidate: a valid artifact exists, and it is an auth error |

No recorded incident would have been automated in its actual context. The benefit is prospective: for example, the
owner's default fleet (three non-facilitators, a non-participating facilitator) where one participant hits an
explicit long-window limit through a verified recognizer and the other two are usable.

## Agreed trade-offs

- **Strict provenance over coverage.** Gateway-local locks and unsupported adapters keep reaching the owner.
  claude/text is among them, so for now a claude-1 participant's exhaustion still escalates.
- **A fixed non-facilitator floor of 2.** Two-participant ideas, this one included, never auto-exclude.
- **A 60-minute threshold.** Real sub-hour windows go to the owner. In incident 7 the owner chose to wait.
- **Current `participants:` over a stable P0.** Every awaited consumer already reads `participants:`
  (`workspace.go:356` feeds `consensus.go:123,176,298`, `runner.go:1003`, `wait.go:156` and
  `organizer.go:168,249`), and line 451 presumes removal. The cost is the history-backed `known` at two
  consensus call sites. A missed site fails loudly (`TriageMalformed`), never as a silent close.
- **Every deck changes even with the knob off.** The C1 and bare-503 fixes change today's behavior.
- **Staged delivery** buys a smaller first change. Stage 2 still holds the crash and replay risk.
- **Size.** The protocol rule is small. Correct mid-idea recovery is not, and FINAL says so.

## Open items deferred to implementation

- **The form of the batch commit.** Either a single committed event whose unreconciled projections read as
  pending, or an intent record plus an applied marker. Both participants accept either. It is bound by the
  Q7 requirements and the tests.
- **The adapter support table.** Establish native-error provenance per adapter (claude/text, codex, kimi
  stream-json, zcode). An adapter that cannot be verified ships unsupported.
  - Still UNVERIFIED: whether `claude -p` can exit non-zero with model text.
  - Still UNVERIFIED: what OmniRoute's "Unavailable" denotes.
- **Idea-scoped serialization.** `driver.acquireLock` today locks `RunDir/driver.lock` (`loop.go:30`, codex-1
  r3). Reuse the primitive keyed by idea; a run-local lock alone is not enough. `store/events.go:52–65`
  has no sync, so durability is added where the commit needs it.
- **Exact marker grammar, event names and the policy and scope fields.**
- **The packet-rendering check** for the Phase 0 and Phase 5 cross-references.
- **What today's validator does with a stub at a canonical round path.** That is unchecked. The requirement
  is to judge on content.
- **Incident 6.** Recover its raw terminal error and kimi-1's actual role for FINAL's table, if the source is
  reachable. Otherwise keep the conditional wording.
- **Scrubbing.** Whether request and correlation ids are kept in a scrubbed excerpt. The evidence file
  dropped them on purpose; the codex-auth note quotes them.

## Comparison & blind spots

**Raw disagreements, kept visible:**
1. **The no-reset branch.** This is the OPEN item in Q2. codex-1 lets an hourly-or-longer allowance qualify;
   claude-1 requires 24 hours or longer. It is unanswered because the two round-03 files were written in
   parallel.
2. **The authority of the claude/text recognizer.**
   - claude-1 r3: fixture-gated.
   - codex-1 r3: provenance-gated, for every adapter (ALT-8).
   - The drafter adopts codex-1's position; see "Drafter position changes".
3. **Membership representation.** Each participant adopted the other's round-1 view in round 2, so the
   concessions crossed. Round 3 converged on current `participants:` plus history.
   - The `known` formulas differ in wording. claude-1 r3: "P0 = `participants:` ∪ the ids of mid-idea
     automatic markers whose transition has an applied record". codex-1 r3: "identities that actually
     participated under that history".
   - Consensus takes the history form. It covers owner-confirmed revisions too, and it does not parse
     markers (blind spot 2).
4. **The transition record.** The concessions crossed again in round 3. claude-1 adopted intent plus applied;
   codex-1 withdrew it in favor of one batch event. The form is left to implementation.
5. **Incident value.**
   - claude-1 r2 said "Stage 1 delivers all the realized value". r3 said "only incidents 1/2 are
     automated".
   - codex-1 r3 showed that the real kickoff fails the floor. claude-1 r3's incident-6 "three-exit gate" is
     not established by the excerpt.
   - Both corrections are carried in the incident table.

**Covered by one participant only:**
- **codex-1:**
  - the executed C1 probe, showing `agent.started agent=kimi` after exclusion;
  - store durability, and the run-local lock;
  - recording the authorized scope;
  - the adversarial negative fixtures;
  - `is_error` alone being insufficient;
  - the manual close route.
- **claude-1:**
  - the claude/text stdout evidence;
  - the OmniRoute reports (#10905, #12817, #14072);
  - the D4 limit: an unknown signer's ❌ becomes `TriageMalformed`, because errors outrank blocks;
  - the sub-day narrowing;
  - staging;
  - the stub-at-path case.
- **Raised by codex-1 r1 and confirmed by claude-1 r2:** the UTC+8 miscopy of the zcode reset (C5).

**Blind spots that no round file addressed** (new PRIMARY readings by the drafter):
1. **A live artifact-wins case with an auth error and a gateway timer.**
   `runs/20261003T103229.878442000Z/round-03/codex-1/exit.json` reads `"seconds":702,"exit":1`.
   `stderr.log:8005` reads `ERROR: unexpected status 401 Unauthorized: [codex/gpt-6-astra] [401]: Encountered
   invalidated oauth token for user, failing request (reset after 12s), …`. Even so, `round-03/codex-1.md` is
   complete (22646 B). That gives two real negative fixtures:
   - exit ≠ 0 plus a provider error plus a valid artifact means no candidate;
   - a gateway timer on an auth error never qualifies.
2. **Repeated `excluded:` keys collapse when read.** `internal/protocol/workspace.go:376`
   `meta := map[string]string{}` and `:394` `meta[strings.TrimSpace(key)] = strings.TrimSpace(value)`: the
   last repeated key wins. This idea's own `00-prompt.md:9–10` has two `excluded:` lines. So `known` must
   come from the recorded history, and markers stay display records. claude-1 r3's line-451 wording would
   have read them.
3. **The organizer and a participant can share one account and one gateway.** That is true in this idea:
   both are claude-1. Exhaustion can then stop orchestration and the participant together. That stays out
   of scope, but FINAL should say so.

**No verdict-conflicts section (§15.3).** No contradictory verdicts were issued. codex-1's round-03 challenges
(incident 6, the exclusivity of claude/text, the value of stage 1) are scoping verdicts on claims that
claude-1 owns. No opposing verdict exists for any of them.

**Correlated agreement (§15.6(b)).** codex-1 (OpenAI) and claude-1 (Anthropic) are different model families.
Even so, both proposals are one architectural family over the same driver. Unanimity is a shared prior, not
independent evidence that the failure paths are safe. The agreed position is wrong if any of these
happens:
- a provider uses allowance wording for a short throttle with no reset;
- a content-triggered exclusion;
- a lost veto;
- a consumer that dispatches from a stale list;
- a replay that launches mixed membership.

The fixtures, the `TriageBlocked` test and the replay tests guard against these.

## Drafter position changes

Every material change in claude-1's position since `round-03/claude-1.md`:

1. **The authority of the claude/text recognizer.**
   - Prior: "So until a positive fixture is recorded, the claude adapter stays on the human path"
     (`round-03/claude-1.md`, "Two items you have not seen yet"). That is, recorded positive and negative
     fixtures were enough to promote it.
   - New: codex-1's provenance gate. That means located CLI behavior or source plus adversarial negatives,
     for every adapter.
   - Reason: a whole-stream match does not settle the UNVERIFIED assumption.
2. **The transition form.**
   - Prior: "I take yours, because it lets read-only `status`, `wait` and the brief report a pending
     transition without recomputing projections." That was intent plus applied, as required
     (`round-03/claude-1.md`, item 1).
   - New: either form, bound by the pending-visibility requirement and the tests. codex-1 withdrew its
     requirement in the same round.
3. **`known` and the line-451 wording.**
   - Prior: "P0 = `participants:` ∪ the ids of mid-idea automatic markers whose transition has an applied
     record" and "automatic `excluded:` markers are read only to keep a removed participant's filed
     signoffs known" (`round-03/claude-1.md`, counter-proposal).
   - New: `known` comes from the recorded membership history, and markers are never parsed for authority.
   - Reason: blind spot 2.
4. **Incident 6.**
   - Prior: "| 6 | three-exit gate |" (`round-03/claude-1.md`, incident replay).
   - New: the conditional row in the table above.
5. **Incidents 1/2 and the value of stage 1.**
   - Prior: "After stage 1 alone, only incidents 1/2 are automated", and "stage 1 holds all the automation
     the recorded incidents would have used (incidents 1/2)" (`round-03/claude-1.md`).
   - New: at the actual kickoff the floor fails, so no recorded incident would have been automated.
6. **Positive forms.**
   - Prior: "Positive fixtures: "Weekly/Monthly Limit Exhausted" and OmniRoute's "All … accounts have
     exhausted their quota"" (`round-02/claude-1.md`). round-03 did not withdraw it.
   - New: gateway-pool phrases are not positive forms until their cause is established. A weekly phrase must
     report that the allowance was reached.
7. **Riders accepted from codex-1 r3:**
   - the authorized scope is recorded with the policy, and an upgrade or resume never widens it;
   - the relaunch suggestion is labeled a provider estimate, with none when the reset is unknown;
   - checked durability, a truncated-record fail-closed path, and a competing-run test.

   These are additions; no prior claude-1 text opposed them.

## Alternatives disposition

- **ALT-1** One shared refinement over the existing classifiers (`failclass.go`, `providerFailureClass`,
  the telemetry collector). **Adopt**, gated per adapter by established provenance. The broad
  `rate-limit`/`billing` labels alone never authorize an exclusion.
- **ALT-2** Stdlib reset parsing, with no timer. **Adopt**, at 60 minutes, not codex-1's round-1 15 minutes.
  Exclusion lasts the rest of the idea, and stated resets are least reliable in that band.
- **ALT-3** The existing exclusion record (`preflightReport.Excluded`, the `excluded:` writer). **Adopt**,
  with an automatic marker and real filtering. **Reject** `parley roster set --state inactive` (a
  permanent change for a temporary window) and `--yes` as a blanket quota waiver.
- **ALT-4** The existing persistence primitives (`store.Append`, `runmanifest`, `acquireLock`,
  `WriteFileAtomic`). **Adopt**, keyed by idea, with checked durability. A run-local lock alone is
  insufficient.
- **ALT-5** The known/required signer split (`validateDocumentAwaiting`). **Adopt**, with `known` from the
  recorded history. It preserves vetoes without a second signoff format.
- **ALT-6** The existing retries, manual pause and blocking escalation. **Adopt** (kept) for everything not
  positively classified. No retry storm and no sleep worker.
- **ALT-7** Preflight-only exclusion. **Adopt** as stage 1 of the delivery order. **Reject** as full
  completion, because the recorded incidents include mid-idea failures.
- **ALT-8** The claude/text one-line matcher as present authority. **Reject**, because its provenance is
  unproven. The claude `--output-format json` switch is **rejected for this idea** and deferred to a
  scoped follow-up, because it changes the runner's fallback-artifact path.
- **ALT-9** Deriving the effective quorum by subtracting `excluded:` lines (claude-1 r1
  `effectiveParticipants`). **Reject.** Free-form lines must not authorize, and repeated keys collapse when
  read (blind spot 2).
- **ALT-10** A stable P0 in `participants:` plus a reducer at every consumer (codex-1 r2). **Reject.** Every
  consumer reads `participants:` directly, line 451 presumes removal, and a missed reducer stalls silently.
- **ALT-11** Querying remaining credits, billing APIs and predictive budgeting. **Reject.** Nothing like it
  ships (claude-1 r1 null search), and it is a non-goal.
- **ALT-12** A reset window alone as the trigger, or a generic `reset after ⇒ quota` regex (claude-1 r1).
  **Reject.** The gateway appends timers to circuit breakers (#10905) and to auth failures (blind spot 1).

## Signoffs

<!-- codex-1 appends its signoff below -->
<!-- claude-1 appends its signoff below -->
