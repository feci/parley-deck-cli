---
from: claude-1
to: user
idea: meta-protocol-change-quota-auto-exclude
phase: final
blocking: yes
date: 2026-10-03
---

## Question

The design for "drop agents from the roster automatically when they run out of tokens or credits, minimum
roster 2" is finished. Both participants signed it (codex-1 ✅, claude-1 ✅) and FINAL is published in
`parley-deck/ideas/meta-protocol-change-quota-auto-exclude/FINAL.md`.

Do you approve this design, with your answers to decisions D1 to D7 below, so that codex-1 can implement
it in a later, attended session?

Nothing is implemented, published, pushed or merged. Everything is committed on branch
`quota-auto-exclude` in the CLI worktree. FINAL is the authoritative text, and this note summarizes it.

## The rule in brief

1. When a participant's invocation **fails**, it leaves **no valid artifact**, and the provider's error
   **explicitly says** the account's credits or quota are exhausted, or a named usage allowance was
   reached, Parley drops that agent from **this idea's quorum only**, without asking you.
2. The error must come from a recognizer whose provenance is verified for that CLI adapter. Without one,
   the adapter stays on today's path.
3. A stated reset must be **at least 60 minutes** away. With no stated reset, only credit or quota
   exhaustion, or a daily, weekly or monthly allowance, qualifies.
4. At least **2 non-facilitator participants** must remain after the whole batch of failures. The rule is
   all or nothing. Otherwise nothing changes and you are asked, as today.
5. A designated or pinned implementer, or an agent that has started a consensus or FINAL draft, is never
   dropped automatically.
6. Anything ambiguous stays on today's path, where you confirm: a bare 429, a 5xx or a "(reset after N)"
   timer, a 400, an authentication error, a hang, or any text from an artifact, a tool or the model. In
   short, it fails closed.
7. The dropped agent's filed ❌s, `DISPUTED` claims and findings remain in force. Dropping an agent can
   never remove a dissenter.
8. No `agents.toml` file is ever written. The agent is probed again at the next idea and rejoins the same
   idea only with your confirmation.
9. You get one non-blocking inbox notice for each drop. `parley status`, `parley wait` and the organizer
   brief show the current quorum.
10. A new setting, `[defaults].quota_auto_exclude`, is on for new ideas and has a per-idea opt-out. It is
    delivered in two stages: kickoff first, then mid-idea.

## What would have happened in each incident

| # | Incident | Under this rule |
|---|---|---|
| 1, 2 | zcode-1 429 "Weekly/Monthly Limit Exhausted", reset in 49 h, at this idea's kickoff | The message qualifies, given a verified zcode recognizer. But only codex-1 was a usable non-facilitator, so the floor fails and you are asked, as happened. |
| 3 | kimi-1 readiness probe hung | Asked, as today: a hang is not evidence. |
| 4 | codex-1 503 "Unavailable (reset after 5h 51m 11s)" (windows-portability organizer) | Asked, as today: gateway text that states no cause, and an organizer is out of scope. |
| 5 | kimi-1 400 from a disabled gateway connection | Asked, as today. |
| 6 | kimi-1 403 "weekly usage limit" during an implementation | If kimi-1 was the designated or pinned implementer, the existing three-exit handoff gate opens. Otherwise its raw error and the floor decide. The excerpt cannot settle which, so FINAL keeps this row conditional. |
| 7 | claude-1 503 "Unavailable (reset after 55m 29s)" (this run, round 2) | Asked, as today. You chose to wait and relaunch once. |
| 8 | codex-1 401 "invalidated oauth token … (reset after 12s)" (this run, round 3) | Never a candidate: its artifact was complete, and it is an authentication error. |

None of the recorded incidents would have been automated in its real context. The benefit is
prospective. In your default fleet (three non-facilitators and a non-participating organizer), a
participant that hits a verified long-window limit is dropped with a notice and the run continues.

## Changes (FINAL §13)

**Protocol.** Identical hunks go into `parley-deck/COOPERATION.md`,
`internal/protocol/defaults/COOPERATION.md` (drift guard) and the skill's
`references/COOPERATION.md`.

- **§9.0 "Excluding".** The binding rule: signal, provenance, threshold, floor, role guards, record, notice,
  return, the policy key, and "mid-idea: not yet in force" until stage 2 ships.
- **§5.** "does not silently shrink quorum" gains "except the recorded §9.0 quota auto-exclusion". One
  bullet says that filed signoffs, ❌s and findings survive.
- **§4 Phase 5.** A designee or pin is never auto-excluded. The line-451 sentence is amended to the new
  membership model: `participants:` is the current set, and history is kept.
- **§0 `[defaults]`** gains the key, and the **Phase 0 template** gains the opt-out.
- **Cross-references** in Phases 3, 6 and 7 and in the §9 checklist.
- A **changelog** entry.
- **Packet map.** No new entry: §5 and §9.0 are always included. Rendering is checked for phases 0, 5
  and 8.

**CLI, stage 1 (kickoff):**

- a new `internal/telemetry/quota.go`, holding the recognizers and their support table;
- wiring in `usage.go`, `failclass.go`, `telemetry.go`, the runner and `preflight_liveness.go`;
- `preflight.go` and `app.go` for the kickoff decision and the C1 fix (below);
- `config/runtime.go` for the setting;
- `workspace.go`, `runcontrol`, `runmanifest` and `runstate` for the records;
- `wait.go` and `organizer.go` for the views.

**CLI, stage 2 (mid-idea):**

- a new `internal/driver/quota.go`, holding the batch decision, a durable transition record and a
  per-idea lock;
- the driver loop, replay, history-based signers in `consensus.go`, the driver's consensus and
  implementation paths, the store and the runner evidence paths.

**Skill.** `SKILL.md` and `references/ROSTER_AND_PROTOCOL.md`, plus the bundled protocol copy.

**Tests.** FINAL lists about 15 required case groups and 21 acceptance criteria (AC1 to AC21). Negative
fixtures come from this run's own failures.

**Two behavior changes for every deck, even with the setting off:**

- **C1.** `parley run --yes` no longer keeps and dispatches an agent it excluded. codex-1 showed today's
  behavior with an executed probe.
- **Bare 503.** A preflight 503 can no longer be excluded with `--yes`.

## Risks and mitigations

- **A provider uses quota wording for a short throttle.** Mitigated by the per-adapter provenance gate,
  adversarial fixtures, the 60-minute threshold and the 24-hour rule.
- **A dissent or veto is lost.** Prevented because filed ❌s, claims and findings survive. A missed code
  path fails loudly as "malformed", never as a silent close (AC9).
- **A crash in the middle of a transition.** Mitigated by one durable batch record, a fail-closed path for
  truncated records, reconciliation before any next dispatch, and fault-injection tests.
- **Size.** The protocol change is small, but correct mid-idea recovery is not. The 300 to 450 line
  estimate is unverified, and both participants say full support is larger than "small".
- **Limited near-term benefit.** Two-participant ideas never auto-exclude, and claude's text-mode errors
  and gateway-pool locks still reach you until their recognizers are verified.
- **Shared account.** The organizer and a participant can share one account and one gateway, as claude-1
  does here. Exhaustion can then stop both. That is out of scope and stated in FINAL.

## Your decisions, with the participants' recommendations

- **D1 (R1). What "roster" means.** It means this idea's quorum only, not the machine roster.
  *Recommended.*
- **D2 (R2). Default.** On by default for ideas created after delivery. Older ideas keep asking.
  *Recommended.*
- **D3 (R3). Fixed values.** The floor (2) and the threshold (60 minutes) are fixed, not configurable.
  *Recommended.*
- **D4 (R4). Delivery.** Two stages under this one FINAL. Both are needed for "complete". The alternative
  is to authorize stage 1 only, with a recorded follow-up idea for stage 2. *Both stages recommended.*
- **D5 (R5). One driving run per idea.** This applies only to ideas whose recorded scope includes mid-idea
  exclusion, so it is not a third change with the setting off. The alternative is a deck-wide singleton,
  which would be a third behavior change. *Scoped version recommended* (the FINAL drafter's reading).
- **D6. Driver accounting blocker (gap 11, owner action).** A legacy run record without an idea field,
  `parley-deck/runs/20260510T194003Z`, makes the driver halt before every new idea's first cross-review
  round. Details are in the non-blocking note `claude-1-to-user_…_driver-gap-11.md`. Two options:
  - run `parley budget migrate apply --kind cross-review` per idea, from an attended terminal;
  - authorize a separate lasting fix.

  *Recommendation (organizer):* a lasting fix, because the per-idea declaration is not stored and would
  repeat for every idea. Until then, recorded fallbacks keep runs moving.
- **D7. Implementation authorization.** Approve Phases 5 to 8 with the staffing below, in a later
  attended session.

## Implementation size, plan and staffing

- **Size.** Stage 1 touches 14 existing CLI files plus one new file. Stage 2 touches 11, two of them shared
  with stage 1, plus one new file. Then come the protocol hunks in three copies and two skill files. The
  line count is unverified.
- **Implementer: codex-1.** It is recorded in FINAL's frontmatter (`implementer: codex-1`) and is also the
  machine default implementer. codex-1 writes the `IMPLEMENTATION.md` plan and checklist before any code.
- **Reviewer: the separately launched claude-1 participant.** It is the **single non-implementer
  reviewer**, because the deliberation track has all non-implementers review.
- **Role concentration (§15.5).** claude-1 is both the organizer and the reviewer, in separate processes.
- **Attended close.** The driver refuses unattended auto-completion with fewer than two independent
  reviewers, and two claude-1 processes never count as two. So Phases 5 to 8 run attended, and the close
  takes both participants' review signoffs and current-tree evidence for AC1 to AC21.
- **kimi-1 and zcode-1 are not planned back** for implementation or review, per your instruction "sakra
  tak pouzi len claude a codex". For your information only: kimi-1 passed readiness again at 19:28 CEST,
  and zcode-1 was still rate-limited.
- **Protocol publication.** After your ratification, it needs the attended `parley protocol publish`
  (§7).

## How this run went

- **Rounds.** Round 1 (independent), cross-review rounds 2 and 3 (the cap allows 3), then consensus, two
  ✅ signoffs and FINAL.
- **The participants converged** on everything except one item. That item, the no-reset branch, closed at
  signoff when codex-1 adopted the stricter 24-hour rule.
- **Interruptions, all recorded in `organizer-notes.md`:**
  - two organizer deaths: the launching session exited, and later the tmux server died;
  - the claude-1 participant's gateway quota error at 14:04, after which you chose a single relaunch;
  - codex-1's gateway 401 at 19:26. It was back by 19:57, and you answered "pokracuj";
  - two 30-minute timeouts, claude-1's round 3 and the driver's consensus drafter. Each was relaunched once
    with a 45-minute limit, under the skill's timeout policy.
- **Driver gaps.** Twelve are recorded. The driver itself scaffolded consensus, ran both signoffs and
  authored FINAL. The recorded fallbacks covered rounds 1 to 3 and the consensus draft. The gaps would
  make a good follow-up idea for the driver.

## Organizer usage

Cumulative client token totals from each organizer transcript, recorded in `organizer-usage.md` and
`usage-ledger.jsonl`. Participants' usage is not included.

| Organizer session | Last event (UTC) | Input (uncached) | Cache read | Cache write | Output |
|---|---|---:|---:|---:|---:|
| 1 (Phases 0 to 2, killed at about 13:02 CEST) | 11:01:10 | 622 | 54,766,551 | 1,090,018 | 292,155 |
| 2 (round-2 relaunch, quota stop) | 12:07:39 | 264 | 14,718,418 | 458,124 | 257,691 |
| 3 (claude-1 round-2 relaunch) | 16:23:54 | 216 | 9,166,220 | 860,913 | 122,670 |
| 4 (round 3 to FINAL, snapshot before this note) | 19:20:50 | 882 | 99,529,951 | 5,884,383 | 617,637 |
| **Total** | | **1,984** | **178,181,140** | **8,293,438** | **1,290,153** |

## What I need from you

1. Answer D1 to D5: approve the recommendations or change them.
2. Choose an option for D6, or tell me to leave it for now.
3. Answer D7: approve or decline implementation by codex-1 with claude-1 reviewing, attended, in a later
   session.

Your answer is quoted under `## User direction` in the next canonical artifact, which is
`IMPLEMENTATION.md` or the ratification record.
