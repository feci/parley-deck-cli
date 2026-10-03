---
agent: claude-1
idea: meta-protocol-change-quota-auto-exclude
round: 1
date: 2026-10-03
---

## Summary

I propose a narrow **quota auto-exclusion** carve-out from the §9.0 rule that exclusion needs explicit
user confirmation. It drops an agent from **this idea's quorum only**, never from `agents.toml`. It
fires only on a positively classified, evidence-recorded provider exhaustion, meaning a reset window
of at least 60 minutes or an allowlisted credit or quota-exhausted phrase. That evidence must come
from stderr, the exit error or a structured provider error event, never from free-form agent output.
The same rule applies pre-idea (preflight) and mid-idea (at a launch that failed with no artifact).
Exclusions at one decision point are evaluated **as a batch, all or nothing**, against a fixed floor
of 2 effective participants. Below the floor nothing is auto-excluded and today's blocking escalation
applies. Filed dissent survives an exclusion. The excluded agent's ❌ signoffs and open
CRITICAL/MAJOR findings keep their full weight, so the rule removes a duty to sign, never an
objection. Everything not positively classified (400s, hangs, timeouts, short 429s, unknown text)
keeps today's user-confirmed path.

Protocol context attestation: `context_mode=full`,
`source_sha256=b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388`,
`packet_sha256=b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388`, `fallback_reason`
absent. HEAD: `04b22e265cc8bb96975bd59cddaf11b737d4cde9` on branch `quota-auto-exclude`.

## Proposed approach

### A. HEAD facts this design rests on (§15 verdicts; the brief's locators were treated as testimony)

| # | Claim | Verdict | Provenance |
|---|---|---|---|
| V1 | §9.0 requires explicit user confirmation for exclusion and for re-inclusion. Exclusion is per-idea. | CONFIRMED | PRIMARY: `parley-deck/COOPERATION.md:880-885`: "**Excluding** an unavailable agent from this idea's quorum requires **explicit user confirmation** … Exclusion is **per-idea and temporary**" and "**Re-including** … **also** requires explicit user confirmation (no silent quorum expansion)." |
| V2 | §5 forbids a silent mid-idea shrink and sets a normal floor of two. | CONFIRMED | PRIMARY: `COOPERATION.md:748-751`: "a mid-idea unavailability does not silently shrink quorum — it falls to the async rules below and the runtime watchdog" and "A valid Parley Deck idea normally has at least two active participants." |
| V3 | A preflight provider failure is never an exclusion, even with `--yes`. | CONFIRMED | PRIMARY: `internal/app/preflight_liveness.go:33-35`: "Environmental, not the agent being dead. Blocking but never an exclusion." Also `internal/app/preflight.go:406-411`: provider failures "are BLOCKING gates that are NOT auto-excluded and NOT waivable into an exclusion by --yes". |
| V4 | The run classifier labels the verbatim codex 503 "reset after 5h 51m 11s" as `overloaded`. | CONFIRMED | PRIMARY, executed: I ran `go test -overlay … -run TestClaude1QuotaProbe ./internal/runner/` with a throwaway overlay test that calls HEAD `classifyFailure("", "", text)`. Output: `B-codex-503-verbatim -> overloaded`, `A-zcode-429-verbatim -> rate-limit`, synthetic `403 … weekly usage limit -> rate-limit`, synthetic `HTTP 403 Forbidden -> auth`, synthetic `400 Bad Request: ambiguous model k3 -> invalid-request`, synthetic `429 …; retry-after: 2 -> rate-limit`, synthetic OpenAI `insufficient_quota` text without "429" `-> unknown`, synthetic `credit balance is too low -> billing`. The overlay was outside the worktree and has been deleted. |
| V5 | At preflight, the same 503 text is **not** a provider failure. It becomes `process-failure`, which `--yes` **does** exclude as "definite". | CONFIRMED | PRIMARY, executed: HEAD `providerFailureClass(503 text)` returned `""` and `providerFailureClass(429 text)` returned `"rate-limit"` (same overlay method, `./internal/app/`). Code: `preflight_liveness.go:130-135` (exit≠0 → `ClassProcessFailure` unless a provider class matches); `:181` `isDefiniteUnavailable` includes `ClassProcessFailure`; `preflight.go:417-421` records `excluded` under `--yes`. So **today `--yes` excludes a 503 quota window but cannot exclude a 429 quota window**. That inconsistency exists independently of this idea. |
| V6 | The run classifier's corpus includes the agent's **stdout** tail, i.e. content. | CONFIRMED | PRIMARY: `internal/runner/failclass.go:58-65`: `b.WriteString(exitError)` … `b.Write(tailOfFile(stderrPath, failTailBytes))` … `b.Write(tailOfFile(stdoutPath, failTailBytes))` and then `corpus := b.String()`. A participant that crashes while writing about "429"/"usage limit", as every participant of *this* idea does, is classified `rate-limit`. |
| V7 | Telemetry classifies 401/403 as `auth-error` and 429 as `rate-limit`, so a 403 weekly limit can get two different labels. | CONFIRMED | PRIMARY: `internal/telemetry/usage.go:248-253`: `case 401, 403: c.failure = "auth-error"` / `case 429: c.failure = "rate-limit"`. |
| V8 | `--yes` exclusions are written as `excluded:` frontmatter lines and nothing reads them back. | CONFIRMED | PRIMARY: writer `internal/protocol/workspace.go:186-191` (`excludedBlock += "excluded: " + line`); format `preflight.go:419` `"%s — %s — confirmed %s"`; `COOPERATION.md:451`: "nothing reads `excluded:`". A non-test grep for `Excluded`/`excluded` finds only writers (app.go, preflight*.go, runcontrol.go, workspace.go). |
| V9 | The preflight floor counts *ready probe rows*, not non-facilitator effective participants. | CONFIRMED | PRIMARY: `preflight.go:435`: `if len(report.Roster) > 0 && available < 2 {`, where `available` counts `entry.Available` rows. |
| V10 | `parley continue` freezes the participant set at the run roster. | CONFIRMED | PRIMARY: `internal/app/app.go:1230`: `applyRosterSnapshotToParticipants(run.Participants, …, m.RosterSnapshot, stderr)`; the driver is built with `Participants: run.Participants` (`app.go:~1245`). |
| V11 | `roster_change_policy` is parsed but not enforced anywhere yet. | CONFIRMED | PRIMARY: `internal/config/runtime.go:287,550-551` parse it; template `:653`: `roster_change_policy = "confirm-breaking"  # auto-add new agents; user confirms drops/breaking changes`; `internal/app/roster_migrate.go:61-63`: "per-deck confirmation honoring roster_change_policy … Those are deferred". |
| V12 | `parley run --yes` keeps the excluded ids inside `participants:`. | UNVERIFIED | Code reading only, not executed: `app.go:1920-1943` passes `Participants: participants` (unfiltered) next to `Excluded: preflightExcluded`, while `workspace_test.go:80-81` calls the writer with an already-filtered list. Needs an executed check before FINAL relies on it. |

### B. The rule (proposed protocol text, in substance)

1. **Qualifying signal (fail closed).** An agent is *quota-exhausted* only when all of these hold:
   - (a) the launch or probe **failed** (non-zero exit or a structured error result) **and produced no
     valid artifact**. An artifact-wins success is never an exclusion.
   - (b) the evidence is read only from **stderr, the exit-error text, or a structured provider error
     event** (the telemetry collector's `is_error`/`turn.failed`/`error` events). The free-form stdout
     tail never counts (V6).
   - (c) one of these **positive** forms is present:
     - a parsed reset window **≥ 60 min** from now: `reset_at` (RFC3339), a `retry_after`/`retry-after`
       number of seconds, or a "reset after 5h 51m 11s" duration;
     - or an allowlisted exhaustion phrase with no window: `limit exhausted`,
       `(weekly|monthly|daily|N-hour) (usage )?limit`, `quota (exceeded|exhausted)`,
       `insufficient_quota`, `credit (balance|exhausted)`. In that case `reset_at` is recorded as
       `unknown`.
   - A 429 with a short or absent window, any 5xx without "reset after", 400s, auth errors without a
     limit phrase, hangs and every watchdog class (`no_first_output`, `stalled`, `timeout`) are
     **not** exhaustion. They keep today's path.
2. **Scope.** The rule applies at two decision points. **Pre-idea**: the §9.0 probe, where the agent
   is not written into `participants:` and gets an `excluded:` line. **Mid-idea**: a participant launch
   in any phase. Four of the five incidents were mid-idea, so pre-idea alone would not answer the
   owner. The organizer and facilitator process is out of scope because the driver does not launch it
   (incident 4's dying codex-1 *organizer* is not covered).
3. **Meaning of "roster".** Per-idea quorum only. The machine roster `~/.parley/agents.toml` and the
   deck `agents.toml` are never written. Exhaustion is temporary, and `active = false` is a membership
   change. The agent is re-probed at the next idea under the unchanged §9.0 rule, so it rejoins
   automatically at the next idea if it then answers PONG. **Within the same idea, re-inclusion stays
   user-confirmed** (V1), with no flapping and no silent expansion. The recorded `reset_at` appears in
   the notice so the owner can re-include it in one step.
4. **Floor = 2, batch, all or nothing.** Let *E* = `participants` − existing exclusions. A declared
   facilitator counts only when `facilitator_participates: true`, which matches the owner's own count
   for this idea ("Dropping either agent would leave one participant"). At a decision point, collect
   every agent that qualified there. Preflight decides after all probes return. Mid-idea, it decides
   when the runner's `wg.Wait()` closes the step (`internal/runner/runner.go:~220-230`).
   **If |E| − |qualified| ≥ 2, exclude them all. Otherwise exclude none** and take today's blocking
   escalation with the classified evidence attached. This is order-independent within a batch.
   Exhaustions at *later* decision points are judged against the then-current *E*. The floor is still
   never crossed, and the outcome is deterministic for a given event order. The §1 solo exception is
   unchanged. Consequence: **a two-participant idea can never auto-exclude**, so the rule mainly helps
   ideas with three or more participants.
5. **Dissent and integrity.** The exclusion trigger is a provider error and never content,
   disagreement or slowness. Exclusion removes the agent's *duty to act*, never the *weight* of what it
   already filed:
   - a ❌ BLOCK it filed stays a block until its counter-proposal is engaged by a new round (§4
     Phase 3);
   - its open CRITICAL/MAJOR findings stay open and must get an explicit disposition with reasoning in
     `review/consensus.md`. "Author excluded" is never a disposition, §15.3 forbids resolution by count,
     and under `strict_gate` they keep the gate open;
   - its ✅ signoffs and round files stay canonical and untouched (append-only).

   A dissenter therefore gains nothing by being excluded, so exclusion cannot be used to remove one.
6. **Implementer and other interactions.**
   - **Per-idea designee or pinned implementer** (`implementer:` in `00-prompt.md` or in
     `IMPLEMENTATION.md`): never auto-excluded. The existing three-exit gate applies (`COOPERATION.md:449`),
     with the evidence and `reset_at` prefilled. An in-flight branch must not change hands silently,
     and the 1.50.0 chain stays untouched (non-goal).
   - **Global-default designee**: existing fall-through with a notice.
   - **Drafter**: if excluded, the existing fallback ("first agent to have submitted a round-01 file")
     runs over *E*.
   - **Reviewers**: the set is recomputed over *E* minus the implementer. Too few reviewers for the
     track, the LE-7/11 "fewer than two independent reviewers" close rule and `require_model_diversity`
     are **re-evaluated, never waived**, so a shortfall escalates.
   - **Cross-review**: "address every other active agent" is validated against *E*, while the excluded
     agent's existing files stay readable context.
   - **Partial artifacts** from the failed launch stay in run logs and are never committed as round
     files.
7. **Recording and notice.**
   - **Idea record.** The existing `excluded:` line format with an `auto` marker and *recorded*, not
     *confirmed*: `excluded: zcode-1 — auto: quota-exhausted (rate-limit; reset_at 2026-10-04T22:14:57Z)
     — recorded 2026-10-03`. Mid-idea, the driver appends this line with the same atomic frontmatter
     rewrite it already uses for `status:` (`internal/driver/driver.go:595` `setIdeaStatus`).
     `participants:` stays as the record of the locked quorum.
   - **Run record.** One `agent.excluded` event in `events.jsonl` with class, parsed `reset_at`, a
     secret-scrubbed evidence tail (`sanitizeTail`/`scrubSecrets`) and the floor arithmetic. Later steps
     reuse the shipped `agent.skipped` event with `reason: auto-excluded:quota`.
   - **Owner notice.** A non-blocking `inbox/<driver-agent>-to-user_<slug>_quota-auto-exclude.md`
     (`blocking: no`) with the evidence, `reset_at` and the re-include instruction. Below the floor it
     becomes today's `blocking: yes` escalation.
   - **Status surfaces.** `parley status` and `parley organizer brief` list *E* and each auto-exclusion
     with its `reset_at`. `parley wait` returns at the exclusion event as a "degradation", consistent
     with its contract to return "loudly on degradation" (§11 advisory). I have not read the
     `parley wait` code; this is a proposal.
8. **Configuration.** One knob, `[defaults] quota_auto_exclude = true|false` in `~/.parley/agents.toml`,
   overridable per deck, with a per-idea `quota_auto_exclude: false` opt-out in `00-prompt.md`.
   **I recommend default-on, because the owner asked for "automatically".** The floor is **fixed at 2,
   with no knob**: the owner stated it, §5 already says "at least two", and an upward-only knob is
   speculative. The 60-minute threshold is a constant. `roster_change_policy` is **orthogonal**: it
   governs membership edits to roster files (V11), and this rule never writes one, so it can never be a
   "drop" under that policy. The protocol text says so explicitly.
9. **Pre-existing defects fixed by this change (needed for it to work).**
   - (i) **V5:** the probe must use the same exhaustion predicate, so the 503 window and the 429 window
     classify alike. A provider failure that does *not* qualify stays a non-excludable gate.
   - (ii) **V12, if confirmed:** `parley run` must not keep excluded ids in `participants:`.
   - (iii) **V8:** the driver needs one `effectiveParticipants(ideaDir, runParticipants)` helper that
     subtracts `excluded:` ids. It is used by round completion, cross-review validation, reviewer
     selection, drafter fallback and `parley continue` (V10).

### C. Replay of the evidence incidents

| Incident | Today | Under this rule |
|---|---|---|
| 1/2 zcode-1 429, `reset_at` +49 h (pre-idea) | blocking gate, owner asked | auto-excluded if the floor holds (4 rostered → 3); notice only |
| 3 kimi-1 probe hang (pre-idea) | blocking gate | unchanged: a hang is not positive evidence. **So the 2026-10-03 kickoff would still have needed the owner once.** |
| 4 codex-1 503 "reset after 5h 51m" | organizer died | out of scope (organizer). As a *participant* launch it would qualify (window ≥ 60 min); as pinned implementer it would go to the three-exit gate |
| 5 kimi-1 400 from a disabled connection | stalled | unchanged (disguised exhaustion; the CLI never saw the quota text) |
| 6 kimi-1 403 "weekly usage limit" at implementation | waiver asked, quorum unchanged | phrase qualifies. As a reviewer it would be auto-excluded; as implementer it goes to the three-exit gate, with the evidence prefilled |

### D. Change inventory (for FINAL; nothing is implemented in this run)

- **Protocol hunks**, identical in all three copies: `parley-deck/COOPERATION.md`, the drift-guarded
  `internal/protocol/defaults/COOPERATION.md`, and the skill's
  `skills/parley-deck/references/COOPERATION.md`.
  - §9.0: one new sub-bullet "Quota auto-exclusion" under "Excluding", with the signal, floor, batch,
    marker and notice.
  - §5: the "does not silently shrink quorum" sentence gains "except the recorded §9.0 quota
    auto-exclusion, which is never silent", plus one bullet on filed signoffs and findings surviving
    exclusion.
  - §4 Phase 5: one clause saying a per-idea or pinned implementer is never auto-excluded.
  - `meta/packet-applicability.yaml` must be checked for the new §9.0/§5 text (it exists at HEAD; I
    did not read its block map). A core change also needs owner ratification and the attended
    `parley protocol publish` (§7).
- **CLI changes, est. 300–450 LOC with tests:**
  - `internal/runner/failclass.go`: a `quotaExhaustion(stderr, exitErr, events) (resetAt, ok)`
    predicate and a `quota-exhausted` class ahead of `overloaded`;
  - `internal/app/preflight_liveness.go` and `preflight.go`: shared predicate, batch floor over *E*,
    auto `excluded:` line;
  - `internal/app/app.go`: filter participants on `parley run` and apply *E* in `runContinue`;
  - `internal/driver/driver.go`: mid-idea batch decision, frontmatter append, *E* in round and
    cross-review checks;
  - `internal/app/driver_impl.go`: reviewer set over *E*, designee exclusion → gate;
  - `internal/runstate`, status and organizer-brief rendering;
  - `internal/config/runtime.go`: the knob and its template line.
- **Tests:**
  - classifier table with the verbatim 429 and 503 texts plus negatives (`retry-after: 2`, bare 503,
    400, `HTTP 403 Forbidden`, a stdout-only "429" that must **not** qualify);
  - floor cases: 3→2 excludes, 2→1 escalates, 4 with two simultaneous → 2 excludes, 3 with two
    simultaneous → none and escalates, permuted input order gives the same result;
  - an excluded agent's ❌ still blocks consensus, and its open MAJOR blocks review close;
  - a per-idea designee's exhaustion → gate;
  - knob off → today's behavior byte-for-byte.
- **Skill:** a SKILL.md paragraph on the organizer's handling (read the notice, never re-include
  silently) and the bundled COOPERATION.md copy.
- **Deliberately left out:**
  - detecting disguised exhaustion (400/hang) and gateway dashboard introspection;
  - automatic in-idea re-inclusion;
  - any `agents.toml` write;
  - a configurable floor or threshold;
  - auto-replacing a designated or pinned implementer;
  - retries, backoff and predictive token budgeting;
  - organizer self-exhaustion.

## Existing alternatives

| Hand-built element in this proposal | Closest thing already shipped (locator) | Reuse decision | Constraint-forced or inherited |
|---|---|---|---|
| Exhaustion class | `internal/runner/failclass.go` `failureRules` (`rate-limit`, `billing`); `internal/app/preflight_liveness.go:627-636` `providerFailureRules`; `internal/telemetry/usage.go:243-253` `api_error_status` mapping | Extend these and add no new classifier subsystem. The existing classes alone are insufficient: they include stdout (V6), conflate a 2 s throttle with a 49 h window (V4, case E), and disagree on 503/403 (V5, V7). | Constraint-forced (fail-closed integrity) |
| Reset-window parser | Go stdlib `time.Parse(time.RFC3339Nano, …)` for `reset_at`; `strconv.Atoi` for `retry_after`; `time.ParseDuration` for "5h51m11s" once spaces are stripped; `net/http.ParseTime` for HTTP-date `Retry-After` | Use the stdlib. Only the regex that locates these fields is hand-written. No in-tree parser extracts them: a grep of `failclass.go`, `preflight_liveness.go` and `usage.go` found none. | Constraint-forced |
| Exclusion record | `excluded:` writer `internal/protocol/workspace.go:186-191`; line format `internal/app/preflight.go:419`; frontmatter rewriter `internal/driver/driver.go:595` `setIdeaStatus` | Reuse the format, adding `auto`/`recorded`, and the atomic rewriter | Inherited |
| Reader of `excluded:` (effective quorum *E*) | none: `COOPERATION.md:451` "nothing reads `excluded:`" | New helper, about 20 LOC | Constraint-forced (mid-idea needs it) |
| Floor check | `internal/app/preflight.go:435` `available < 2` hard stop | Reuse, recounting over *E* with facilitator handling | Inherited, recount constraint-forced |
| Batch boundary | `internal/runner/runner.go:~220-230` `wg.Wait()` before `round.completed`/`round.incomplete` | Reuse as the decision point | Inherited |
| Run event and status display | `agent.skipped` with `reason` in `internal/runstate/runstate.go` (`applyAgentEvent`, event summary) | Reuse for later steps. Add one `agent.excluded` event for the evidence payload. | Inherited, one event constraint-forced (evidence must be recorded) |
| Evidence scrubbing | `internal/app/preflight_liveness.go` `sanitizeTail`/`scrubSecrets` | Reuse | Inherited |
| Membership change | `parley roster set <id> --state inactive` (`internal/app/roster.go`, `roster_set.go`) | **Rejected**: a permanent membership change gated by `--confirm-breaking`, wrong for a temporary quota window | n/a |
| Re-inclusion | §9.0 per-idea re-probe at the next idea (`COOPERATION.md:882-883`), with no code | Reuse unchanged | Inherited |
| Config knob | `[defaults]` parsing pattern for `RosterChangePolicy` (`internal/config/runtime.go:287,550-551,653`) | Reuse the pattern | Inherited |
| Querying remaining credits ahead of time | **Null.** I found nothing shipped. LE-5 cost telemetry (`agent.usage`) tracks spend, not provider balance. Sources consulted: `failclass.go`, `preflight_liveness.go`, `telemetry/usage.go`, `COOPERATION.md` §8 LE-5. | Not built (left out) | n/a |

## Concerns / open questions

1. **"Roster" interpretation (owner decision).** The owner said "vyradis agentov z rosteru"
   ("drop agents from the roster"). I read that as the per-idea quorum, not `agents.toml`, because
   quota resets and a roster write is a confirm-breaking membership change. If the owner meant the
   machine roster, the design changes materially.
2. **Floor counting.** The brief says "count only non-facilitator participants". Read literally, this
   idea (codex-1 plus the participating claude-1) has one, yet the owner accepted it as two. I propose
   counting a declared facilitator only under `facilitator_participates: true`. The owner should
   confirm.
3. **Default-on vs opt-in.** I recommend default-on per the owner's words. The counter-argument is
   that it flips §9.0's confirmation default for every deck at once. The core-change ratification is
   the mitigation.
4. **Threshold.** A fixed 60 minutes covers every observed window (5 h, 49 h, weekly) with wide margin
   over seconds-scale throttles. Should it instead track the track's per-agent timeout (about 5, 15 or
   30 min)? I prefer a constant for size.
5. **Implementer carve-out.** Escalating for a pinned or designated implementer means incidents 4 and
   6 would still reach the owner. That is deliberate. Is that acceptable, or does the owner want
   automatic re-designation?
6. **V12** needs an executed check (codex-1 or round 2) before FINAL cites it as a defect.
7. **Packet applicability.** New §9.0/§5 text may need block entries in `meta/packet-applicability.yaml`.
   Changing that file is itself a §7 change; it belongs in the same idea.

## Risks

- **Content-triggered false exclusion.** Mitigated by stderr/exit/structured-only evidence, the
  failed-and-no-artifact requirement, and a test with a stdout-only "429" that must not qualify.
- **Provider wording drift.** An unrecognized text simply misses and falls back to the user-confirmed
  path. That is safe but loses automation; the allowlist is extended by normal ideas.
- **Erosion of the confirmation principle.** Kept narrow by the single class, a recorded marker, a
  non-blocking notice, in-idea re-inclusion that stays confirmed, and a knob plus per-idea opt-out.
- **Lost model diversity or reviewer count after an exclusion.** Never waived: the LE-3/LE-7/11 gates
  re-evaluate and escalate.
- **Secrets in recorded evidence.** Reuse `scrubSecrets`. The zcode body carries only
  `retry_after`/`reset_at`, but other providers may echo headers.
- **Limited payoff on two-participant decks**, which is the owner's current active fleet. The floor
  makes auto-exclusion impossible there. The honest benefit today is pre-idea pruning on decks with
  three or more participants, plus clearer and faster escalations.
- **Scope creep from fixes (i)–(iii).** They are necessary for correctness, not extras. If round 2
  finds them too large for "small", (iii) is the first candidate to split into a prerequisite idea.
