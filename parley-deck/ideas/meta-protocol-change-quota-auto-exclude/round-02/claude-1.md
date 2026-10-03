---
agent: claude-1
idea: meta-protocol-change-quota-auto-exclude
round: 2
date: 2026-10-03
responding-to: [claude-1/round-01, codex-1/round-01]
---

## Summary

Round 2 moves me most of the way to codex-1's round-1 design. Two pieces of new PRIMARY evidence
decide our main disagreements against my own round-1 text:

1. **The claude adapter reports provider errors only on stdout.** My failed relaunch exited 1 with an
   empty stderr and a one-line stdout `API Error: 503 …`, because the shipped claude headless argv is
   `--output-format text`. My round-1 rule ("the free-form stdout tail never counts") would never see a
   claude quota error.
2. **The owner's gateway (OmniRoute) appends "(reset after N)" to circuit-breaker cooldowns as well as
   to quota locks.** Its stated resets have also been wrong in both directions. A reset window alone
   therefore does not establish quota exhaustion.

**Withdrawn (SELF-CORRECTION, effective now):**
- the standalone "parsed reset window ≥ 60 min" trigger;
- the stdout ban.

**Adopted from codex-1:**
- an adapter-recognized native terminal error that carries explicit allowance or credit semantics;
- the strictly non-facilitator floor count;
- removing the excluded id from `participants:`;
- an explicit supersession of `round.incomplete` before a failed round resumes.

**Still open, all narrow:**
- the remaining-reset threshold (60 vs 15 min);
- how heavy the mid-idea transition record must be;
- one adapter-recognizer choice.

**New finding:** under the converged evidence rule, no recorded mid-idea incident (4–7) would have
been auto-excluded. I therefore propose one FINAL with staged delivery, kickoff path first.

Protocol context attestation: `context_mode=full`,
`source_sha256=b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388`,
`packet_sha256=b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388`, `fallback_reason`
absent. HEAD `7595ee1f01b17835d52d6fe6c4bd69438eefe7bd`. `git diff --stat 04b22e2 HEAD -- internal/ cmd/`
is empty, so every round-1 code locator still holds. OpenViking tools were not loaded in this launch, so
this round uses local sources and the located web sources cited below. Per the launch rules, I have not
read `round-02/codex-1.md`.

## User direction

Source: `parley-deck/inbox/user-to-claude-1_meta-protocol-change-quota-auto-exclude_quota-answer.md`. The
owner's answer, verbatim:

> Question: "The claude-1 participant's round 2 failed with a gateway quota error: "503
> [claude/claude-opus-5-5] Unavailable (reset after 55m 29s)", which should reset at about 15:00. How should
> the run continue?"
> Selected: **"Relaunch once after 15:05 (Recommended)"**. The option read: "The organizer relaunches the
> claude-1 participant's round 2 once, after the reset. codex-1 keeps working on its round 2. If it fails on
> quota again, the organizer stops and reports back."

This file is the output of that single relaunch. I use the answer below as design input, not as a ruling
on the rule. It is the only owner decision we have on a sub-hour gateway window, and it chose "wait for the
stated reset and relaunch once". It did not choose a solo exception, and exclusion was impossible here
because of the floor.

## Responses to other participants

### @codex-1

**Verdicts on your round-01 claims** (non-owner verdicts, §15.1; each is written here in my own file):

| # | Your claim | Verdict | Provenance |
|---|---|---|---|
| C1 | A `--yes` exclusion is recorded, but the id stays in `participants:`. | CONFIRMED | PRIMARY (static reading, not executed). `internal/app/app.go:1922` `code, excluded, stop := runTaskPreflight(…, participants, …)`. `:1939-1943` `runcontrol.Create(runcontrol.CreateOptions{ … Participants:   participants, Excluded:       preflightExcluded,` and the slice is not reassigned in between. `internal/app/preflight.go:287` `return 0, report.Excluded, false` returns only the exclusions. `internal/protocol/workspace.go:224` writes `strings.Join(participants, ", ")`. This also settles my round-01 V12, which I raised as UNVERIFIED and did not own. |
| C2 | An already-failed round stays blocked. | CONFIRMED | PRIMARY. `internal/driver/driver.go:467-468`: `case "round.incomplete":` / `return false, nil // authoritative block`. |
| C3 | Consensus already separates who may sign from who is awaited. | CONFIRMED | PRIMARY. `internal/consensus/consensus.go:523-524`: "validateDocumentAwaiting separates who may SIGN (known) from who is AWAITED (required)." |
| C4 | A reset duration is not a contract that a 503 means quota. | CONFIRMED, scoped to OmniRoute | PRIMARY, located reports (quoted from the fetched issue pages). **OmniRoute 3.8.49, [#10905](https://github.com/diegosouzapw/OmniRoute/issues/10905)**, a connection failure: "✕ [antigravity/gemini-3.6-flash-high] [502]: fetch failed (cause: ECONNREFUSED) (reset after 30s)" and "The circuit breaker triggers for all models (reset after 30s)." **3.8.50, [#12817](https://github.com/diegosouzapw/OmniRoute/issues/12817)**, a weekly-quota lock whose 503 body has no quota word: "codex \| all 1 active accounts cooling down for model gpt-5.6-luna-medium (reset after 38h 42m 37s)" … "lastErrorCode=null, lastError=undefined". After a voucher restored the quota, requests "continued to return 503". **3.8.50, [#14072](https://github.com/diegosouzapw/OmniRoute/issues/14072)**: "All antigravity accounts have exhausted their quota (reset after 5m)", yet the account stayed exhausted "(observed: 11+ hours)". Scope: these are reporters' logs, not the gateway source. None shows the word "Unavailable" from incidents 4 and 7, so what "Unavailable" denotes remains **UNVERIFIED**. |
| C5 | The zcode reset must not be rendered as Berlin local time. | CONFIRMED, with a live instance | PRIMARY. The evidence body (`source-context/EVIDENCE-quota-incidents.md:32`) reads "Your limit will reset at 2026-10-05 06:14:57 (reset after 49h 8m 50s)" with `"reset_at":"2026-10-04T22:14:57.003Z"`. This idea's own `00-prompt.md:10` records "reset 2026-10-05 06:14 local". 22:14:57Z is 00:14:57 CEST; 06:14:57 is the provider's UTC+8 wall clock. The miscopy already sits in a canonical record. |

**New PRIMARY evidence about the claude adapter (bears on your recognizer rule):**
- `internal/agents/discover.go:234`: `HeadlessArgs: []string{"-p", "--model", "{model}", "--effort", "{effort}", "--output-format", "text", …}`.
- `runs/20261003T103229.878442000Z/round-02/claude-1/relaunch-1/exit.json` reads `"seconds":127,"exit":1`.
  - `stderr.log` is 0 bytes.
  - `stdout.log` is 216 bytes, exactly one line: `API Error: 503 [claude/claude-opus-5-5] Unavailable (reset after 55m 29s). This is a server-side issue, usually temporary — try again in a moment. If it persists, check your inference gateway (omniroute.marao.sk).`
  - That launcher's claude argv matches `discover.go:234`.
- `internal/telemetry/usage.go:424-441` `StructuredArgs` returns true only for `--json`, or for `--output-format`/`--format` set to `json`/`stream-json`. With this argv there is no structured envelope for the collector's claude `result`/`is_error` branch (`usage.go:244-262`).

**Where I adopt your position:**

1. **Provenance.** I adopt "adapter-recognized native terminal error". One refinement: the per-adapter
   recognizer table must name claude's text mode explicitly. Otherwise claude, the adapter of a
   participant in this very idea, falls into your "ambiguous mixed text" bucket with no automation at all.
   - **Proposed claude/text recognizer.** All of these must hold:
     - exit ≠ 0;
     - no valid artifact;
     - the **entire** stdout is one line matching `^API Error: (\d{3}) `. This is a whole-stream match,
       never a tail substring.
   - **Assumption.** `claude -p` cannot exit non-zero with model-authored final text. This is
     UNVERIFIED (one observation) and must be pinned by recorded fixtures before the recognizer is
     trusted.
   - **Alternative.** Switch claude to `--output-format json` and reuse the collector's claude branch.
     That is sturdier, but it changes stdout for the runner's fallback-artifact path. The kimi stream-json
     unwrap in `internal/runner/kimistream.go` is the precedent. I would make it a follow-up, not part of
     this idea.
2. **503 and reset windows.** Withdrawn as above. I adopt your rule: no generic `reset after ⇒ quota`.
   - Observed pattern (an inference from five samples, UNVERIFIED as a general rule):
     - **Pass-through form.** `[provider/model] [NNN]: <upstream message> (reset after …)` carries the
       upstream's own semantics. Examples: the zcode 429 at `EVIDENCE-quota-incidents.md:32`, and #10905.
     - **Gateway-local form.** `[provider/model] Unavailable (reset after …)` and "all N active accounts
       cooling down" state no cause. Example: incident 4 at `EVIDENCE-quota-incidents.md:53`, "503 Service
       Unavailable: [codex/gpt-6-astra] Unavailable (reset after 5h 51m 11s)".
   - So the recognizer should judge the **upstream message**.
   - Positive fixtures: "Weekly/Monthly Limit Exhausted" and OmniRoute's "All … accounts have exhausted
     their quota".
   - Negative fixtures: "Unavailable" and "cooling down", until the gateway source establishes their
     cause.
3. **Floor count.** I concede. `00-prompt.md:64`: "Count only non-facilitator participants." My round-1
   reading mixed up the owner's authorization of this two-participant run with the floor arithmetic. The
   owner's default fleet is unaffected because its facilitator does not participate.
4. **Representation.** I concede: remove the id from `participants:` and add the `excluded:` marker.
   - Existing protocol text already requires this: "a §9.0 exclusion of the designee must also remove
     the id from `participants:`, or the designation stays live" (§4 Phase 5).
   - Designation validation reads only `participants:`.
   - This replaces my round-1 `effectiveParticipants` helper in most consumers. `known` signers become
     `participants:` ∪ the `excluded:` ids, which matches C3.
5. **Roles and the rest.** I adopt the following:
   - **Drafter gate.** A candidate who has started a consensus or FINAL draft gates. Otherwise today's
     fallback applies.
   - **Designee or pin.** If any candidate is a per-idea designee or the pinned implementer, nothing in
     that batch is applied automatically.
   - **Defaults.** Default-on is resolved and recorded at kickoff, and runs already in flight keep the
     confirmation behavior.
   - **`parley wait`.** It does not return for a successful recorded reduction. I concede that this is
     not "degradation"; the notice and status/brief suffice.
   - **Kickoff.** The exclusion decision is made inside `parley run`, over the exact participant set.
     Standalone `parley preflight` only reports candidates.
   - **`roster_change_policy`.** We agree on behavior: no roster file is written and the policy does not
     gate the rule. The protocol text should state both facts rather than calling it either an
     "exception" or "orthogonal".

**Where I refine or still differ:**

- **Mid-idea transition record (refinement).** I accept C2, and I accept that replay must consume
  exclusions (your `runstate.go:130-138` locator, which I did not re-read; SECONDARY on your reading).
  - **Simpler commit point.** Make the run's append-only `events.jsonl` `agent.excluded` event the
    single commit point. It is keyed by the settled batch's invocation ids.
  - **Derived mirrors.** The `participants:`/`excluded:` rewrite, the manifest mirror and the
    superseding round evaluation are derived. Resume re-applies them idempotently and skips any that
    are already present.
  - **What that removes.** No separate pending/applied state is needed.
  - **Determinism.** A crash before the event leaves the batch undecided, and the same terminal evidence
    gives the same decision again. A relaunched agent's later success supersedes old evidence because
    the decision reads only the settled batch's own invocations.
  - **Idea-scoped lock.** I accept one (reuse `driver.acquireLock`) if two live runs on one idea are
    actually possible. FINAL should state whether they are.
- **Threshold (disagreement): I keep 60 min. Your proposal is 15.**
  - Exclusion is the action that holds for the rest of the idea. Escalation is the conservative one.
  - C4 shows stated resets are least trustworthy in exactly the band where the number matters: 5 m stated
    against 11 h actual, and a 38 h countdown that outlived a restore.
  - The one owner decision in that band (55 min, this idea) chose waiting over losing the participant.
    Under 15 min, with an explicit-quota text and three or more non-facilitators, the rule would have
    dropped that participant for the whole idea.
  - **Counter-proposal:** 60 min. In addition, the blocking notice for an explicit-quota error with a
    stated reset under 60 min, or one blocked by the floor, prefills "relaunch once after the stated reset
    + 5 min" as the recommended option. This is text only, with no sleep worker, consistent with your
    ALT-6.
- **Staging (new; consistent with your ALT-7 labelling demand).** The table under "Refined position"
  shows that only incidents 1/2 (kickoff) would have been automated. Mid-idea automation changes no
  recorded outcome; it serves future cases such as a reviewer who hits an explicit weekly limit. Stage 1
  delivers all the realized value with no transition machinery. Stage 2 carries the mid-idea work you
  specified.

## Refined position

1. **Trigger.** All of the following must hold:
   - The invocation failed (exit ≠ 0 or a structured error result) and produced no valid artifact.
   - Its terminal error was captured by that adapter's recognizer: a structured envelope, the claude/text
     whole-stdout line, or codex's native error line. Each recognizer is pinned by recorded fixtures.
   - That error's own text (for a gateway pass-through, the upstream message) carries **explicit
     allowance or credit semantics**:
     - exhausted credits, for example `insufficient_quota` or `credit balance is too low`;
     - or an explicitly reached or exhausted hourly-or-longer allowance, for example
       `Weekly/Monthly Limit Exhausted`, "weekly usage limit" or "accounts have exhausted their quota".
   - If a reset is stated, it is at least the threshold from observation (60 min; still disputed).
     Otherwise `reset: unknown`.
   - **Never** counts: a bare 429; a generic `quota exceeded`, which may be a per-minute rate quota; a
     5xx or "(reset after N)" without allowance semantics; a long `Retry-After` alone; 400s; auth errors;
     hangs and watchdog classes; any artifact text, tool output or model prose.
2. **Decision points.**
   - Kickoff, inside `parley run`.
   - Stage 2: the settled dispatch batch, before any next dispatch, signoff evaluation or close.
   - A failure of the organizer itself stays out of scope.
3. **Floor.** At least 2 **non-facilitator** participants must remain after the whole batch. The
   facilitator never counts, and the decision is all or nothing. Otherwise nothing is applied and a
   blocking escalation lists every candidate. The §1 solo exception is unchanged.
4. **Role guards** (they also make the batch all or nothing):
   - a per-idea designee or pin: the existing three-exit gate, prefilled with the evidence;
   - a started drafter: gate;
   - a global-default designee before Phase 5: today's fall-through.
5. **Record.**
   - Remove the id from `participants:`.
   - Add an `excluded:` line marked automatic, never "confirmed". It carries the UTC RFC3339 stated reset
     and the raw provider string (C5) plus the invocation id.
   - Write one `agent.excluded` event, which is the commit point. It holds the original and effective
     sets, the rule id and a scrubbed excerpt.
   - Send one deduplicated, non-blocking inbox notice.
   - `status` and the organizer brief show the effective quorum.
6. **Dissent and gates.** Exclusion resolves nothing and waives nothing:
   - A filed ❌ still stands.
   - A `DISPUTED` claim still stands.
   - Open CRITICAL/MAJOR findings, and every strict-gate finding, still need an explicit disposition.
   - Reviewer count, the LE-7/11 two-reviewer close, the goal checker and `require_model_diversity` are
     re-evaluated, and a shortfall escalates.
7. **Return.**
   - The stated reset is a hint, never a timer.
   - There is no automatic rejoin within the idea; re-inclusion stays owner-confirmed.
   - The agent gets a fresh probe at the next idea.
8. **Configuration.**
   - `[defaults].quota_auto_exclude` is a presence-aware boolean.
   - It defaults to true for ideas created after ratification and delivery.
   - The resolved value is recorded at kickoff.
   - An idea can opt out with `quota_auto_exclude: false`.
9. **Staging.** One FINAL with one ratified protocol text.
   - **Stage 1:**
     - the shared recognizer and refinement for all three classifiers;
     - the kickoff path;
     - the C1 filtering fix;
     - notices.
   - **Stage 2:** mid-idea application:
     - `round.incomplete` supersession;
     - replay;
     - `known`/`required` signers;
     - the lock.
   - Until stage 2 ships, the text marks mid-idea application "not yet in force" (the §7 pattern), so it
     never describes intended behavior as present fact.

Incident replay under this rule:

| Incident | Outcome | Decisive reason |
|---|---|---|
| 1/2 zcode 429 "Weekly/Monthly Limit Exhausted", reset +49 h (kickoff) | auto-exclude if at least 2 non-facilitators remain | explicit allowance in the pass-through message |
| 3 kimi probe hang | human path | silence is not evidence |
| 4 codex 503 "Unavailable (reset after 5h 51m 11s)" | human path | gateway-local text with no stated cause (C4); the organizer is out of scope anyway |
| 5 kimi 400 from a disabled connection | human path | no quota semantics |
| 6 kimi 403 "weekly usage limit" at implementation | three-exit gate | allowance qualifies, but the implementer is pinned |
| 7 claude-1 503 "Unavailable (reset after 55m 29s)" (this round) | human path | gateway-local text; also the two-participant floor; also under 60 min |

The honest consequence: for the owner's OmniRoute-fronted fleet, gateway-local locks keep reaching the
owner. Closing that gap belongs at the gateway (surface the cause) or in a later gateway-aware recognizer
backed by the gateway's source, not in this rule.

## Remaining disagreements

1. **Threshold: 60 min (mine) vs 15 min (codex-1).** This is a policy choice, not a fact. My
   counter-proposal is above. If round 3 does not settle it, it should go to the owner as one line.
2. **Weight of the transition record.** I propose the event as the commit point; codex-1 proposed an
   explicit pending/applied transition. Both must pass the same tests: a crash at each write, a double
   replay, competing continuations and stale evidence. FINAL can fix the tests and leave the mechanism
   to the implementer.
3. **claude/text recognizer vs a switch to `--output-format json`.** codex-1 has not seen this evidence
   yet.
4. **Staging.** This is new, and codex-1 has not responded to it.
5. **Still UNVERIFIED:**
   - what OmniRoute's "Unavailable" denotes;
   - whether `claude -p` can exit non-zero with model text.
6. **Owner items, not disagreements between us:**
   - "roster" means the per-idea quorum;
   - default-on;
   - the threshold.
