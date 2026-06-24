---
agent: codex-1
idea: loop-engineering-research
round: 1
date: 2026-06-22
---

## Summary

Parley Deck is already a strong inner cooperation loop for one idea: durable files, non-solo quorum, maker/checker phases, review/fix-up repetition, run event logs, driver cursors, and pipeline effects. The missing loop-engineering layer is the outer driver above an idea: recurring discovery, explicit loop budgets, a goal-done verifier, and durable cross-run goal state.

CLI finding: `parley run` now auto-drives round-01 through consensus/finalization and optionally Phase 5-8, `internal/driver` has re-entrant disk-derived state, and `parley pipeline auto/watch` covers block sequencing and one-shot monitoring. But the caps are fragmented: `internal/driver/loop.go` has a hard-coded 30 minute await deadline, `driver.Config` has hidden `MaxRounds` and `MaxFixupCycles`, `pipeline auto` has `--max-blocks`, `--rounds`, a hard-coded implementation `maxCycles = 3`, and a DAG wave cap. There is no top-level cost budget, no user-facing driver-step budget, no cross-run goal ledger, and no independent goal-done check comparable to the `/goal` pattern.

## Six-block mapping

- Automations / Trigger: partial. `parley run` is human-triggered; `pipeline watch` can be scheduled externally and can auto-open low-risk remediation ideas. Missing: a first-class discovery loop that reads commits/issues/CI/logs and opens candidate ideas with dedupe and provenance.
- Worktrees: covered by `parley-worktrees`, but the CLI driver should treat worktree allocation as a bounded execution resource for future outer loops, not as an informal implementation detail.
- Skills: covered. Parley Deck itself is a skill, and repo docs/config externalize launch intent. No major change beyond making loop policy explicit.
- Plugins / Connectors: partial. COOPERATION.md §12 has provider capabilities, gates, effects, and idempotency. The Go CLI emits provider calls rather than performing side effects, which is the right boundary for now.
- Maker / Checker: partial. Phase 5 vs Phase 6 is structurally correct, and `internal/app/driver_impl.go` chooses a non-implementer review-consensus drafter when possible. Missing: review prompts do not explicitly default to refutation, and the CLI does not enforce or prefer a different verification model/profile.
- Durable State: partial. Per-idea artifacts, `parley-deck/runs/<run-id>/run.json`, `events.jsonl`, `driver.json`, and `pipeline-run.json` are good. Missing: durable project-level state for standing goals and "what should the loop consider next" across runs.

## Placement: Protocol, CLI, Defaults

| Concern | Protocol invariant | CLI flag / feature | `~/.parley` `[defaults]` |
|---|---|---|---|
| Outer-loop triggers | Auto-opened work must record source, dedupe key, risk, participants, and gate policy; it must not create participant artifacts or bypass quorum. | New `parley loop tick` one-shot command, intended for cron/GitHub Actions/MCP events; adapters under a new `internal/loop`. | Default trigger policy: disabled by default, allowed sources, max ideas per tick, default risk policy. |
| Max iterations | Any automated driver loop must declare iteration ceilings; hitting a ceiling escalates, never marks success. | `parley run --max-driver-steps`, `--max-rounds`, `--max-fixup-cycles`; align `pipeline auto --max-blocks` and DAG wave caps with the same budget object. | `[defaults.loop] max_driver_steps`, `max_rounds`, `max_fixup_cycles`, `max_blocks_per_tick`. |
| Max wall-clock | Automated loops must have wall-clock ceilings and durable stop reasons. | `--max-wall-clock`, `--round-deadline`; replace `internal/driver/loop.go`'s hard-coded `roundDeadline`. | `[defaults.loop] max_wall_clock_ms`, `round_deadline_ms`. |
| Max cost | External-backend automation must have a cost policy; unknown cost under a strict budget must halt or require an explicit override. | `--max-cost-usd`; runner emits `agent.usage` / `loop.budget` events when telemetry is available. Do not blindly pass provider spend caps through if they abort before artifacts land. | `[defaults.loop] max_cost_usd`, `strict_cost_accounting`, optional per-backend cost policy. |
| Goal-done gate | Automated completion must be checked against a durable done condition and observable acceptance criteria; budget exhaustion is not done. | `--goal <condition>`, `parley goal check`, and a driver gate before `ActionComplete`; checker should be a non-maker agent/profile. | `[defaults.loop] goal_checker`, `goal_check_model/profile`, `goal_check_required_for_auto`. |
| Durable goal state | Project goals are canonical project state, not local preference. They must live under `parley-deck/` and survive runs. | `parley goals list|add|pause|close`; `parley loop tick` reads open goals before discovering new work. | Only local behavior defaults, never canonical goal content. |

## Proposed approach

Treat Parley Deck as two nested loops:

1. Existing per-idea cooperation loop: Phase 0-8, canonical files, consensus, implementation review, fix-up.
2. New outer discovery/goal loop: trigger -> discover candidate -> dedupe -> open idea or update goal state -> optionally run the existing per-idea loop within explicit budgets -> stop on done, budget, risk gate, or human question.

The protocol should define only invariants and durable artifact contracts. The CLI should implement flags, scheduler-friendly commands, event accounting, and provider adapters. `~/.parley/agents.toml` should provide operator-local defaults for budgets and trigger policy, with project overrides in `parley-deck/agents.toml`; neither should hold project goal truth.

## Concerns / open questions

- Cost accounting will be incomplete for some CLIs. That should not block the design, but strict budget mode must treat unknown cost as a halt condition instead of pretending the budget is enforceable.
- A true daemon is not necessary. A one-shot `parley loop tick` is easier to audit, test, and schedule from cron/GitHub Actions than a long-running resident process.
- `pipeline watch` currently auto-opens remediation ideas with `participants: []` in `openRemediationIdea`. That is acceptable only as a candidate scaffold. If it is an active idea, it conflicts with the non-solo Phase 0 invariant.
- HITL fatigue gets worse if the outer loop opens too many low-value ideas. The CLI needs per-tick caps and dedupe before it needs more connectors.

## Risks

- Runaway cost: highest risk once triggers can start work without a human typing `parley run`.
- False completion: current clean review consensus is strong, but it is not the same as a separate goal-done verifier.
- Comprehension debt: auto-opened ideas can flood the deck with plausible work that nobody understands. Discovery must prefer fewer, higher-confidence candidates.
- Protocol bloat: budgets and goal-state should be generic invariants, not a large provider-specific automation framework inside COOPERATION.md.

## Prioritized recommendations

1. ADOPT - Add a unified loop budget contract. Change `COOPERATION.md` §4 "Stopping judgment", §9.0, and §12 to require explicit max iterations, max wall-clock, and max cost for automated loops, with "budget hit = escalate, not complete." Change `internal/driver/driver.go`, `internal/driver/loop.go`, `internal/runmanifest/manifest.go`, and `internal/app/app.go` to carry a `LoopBudget` into `parley run`. Add CLI flags `--max-driver-steps`, `--max-wall-clock`, `--max-cost-usd`, `--max-rounds`, `--max-fixup-cycles`, and defaults under `~/.parley/agents.toml` `[defaults.loop]`. Effort: medium. Risk: low.

2. ADOPT - Add a goal-done gate for auto-driven completion. Change `COOPERATION.md` Phase 4/5/8 language so driver-managed work must name a durable done condition derived from `FINAL.md` observable acceptance criteria. Change `internal/driver/impl.go` so `ActionComplete` is gated by a goal-check result when auto mode is active. Add `parley run --goal`, `parley goal check`, and `run.json` fields for `goal_condition`, `goal_check_status`, and `goal_checked_by`. Effort: medium. Risk: medium, because bad checkers can false-green.

3. ADOPT - Fix auto-opened watcher ideas before expanding triggers. Change `internal/app/pipeline_cmd.go` `openRemediationIdea` to populate participants from the pipeline manifest or create a clearly non-active candidate file outside `ideas/`. Amend `COOPERATION.md` §12.11 to say watcher-created remediation ideas must satisfy Phase 0 quorum before `status: round-01`. Effort: small. Risk: low.

4. ADAPT - Add a scheduler-friendly outer-loop command, not a resident daemon. Add feature `parley loop tick --source <git|ci|issues|signals> --dry-run --max-open N --yes`, implemented in a new `internal/loop` package and documented in `docs/cli-reference.md`. It should discover candidates, dedupe, write trigger provenance into `00-prompt.md`, then optionally call the existing `parley run` path within the loop budget. Protocol change: new COOPERATION.md subsection "Automation triggers" or an expansion of §12. Effort: large. Risk: medium.

5. ADOPT - Add durable cross-run goal state. Change COOPERATION.md §3 directory layout to include `parley-deck/STATE.md` as the canonical standing-goal ledger, with constrained entries: `goal_id`, status, done condition, source, current idea/run, next allowed action, budget counters, and last checked timestamp. Add CLI feature `parley goals list|add|pause|close|check`. Do not put goal truth in `~/.parley`; only defaults belong there. Effort: medium. Risk: medium due to merge conflicts unless entries are append/update disciplined.

6. ADAPT - Add cost telemetry as a best-effort budget input. Change `internal/runner` and `internal/store/events.go` usage conventions to emit `agent.usage` events where CLIs expose tokens/cost, and `loop.budget` events after each driver tick. Keep `docs/agent-cli-mechanics.md` guidance against unsafe provider spend-cap flags that can abort before an artifact is written. Effort: medium-large. Risk: medium because backend telemetry is uneven.

7. ADAPT - Make verification explicitly refutational. Change COOPERATION.md §4 Phase 6 and `internal/runner/phase58.go` `BuildReviewPrompt` to say reviewers should try to disprove the implementation against FINAL.md, tests, and acceptance criteria. Extend `parley-deck/agents.toml` / runtime config with an optional `verification_profile` or `review_model` preference; prefer a different model/profile than the implementer when available, but do not make that a hard protocol requirement unless the roster can satisfy it. Effort: small-medium. Risk: low.

8. REJECT - Do not adopt a fully autonomous discover -> implement -> merge/deploy daemon. Preserve COOPERATION.md §12.4 and §12.8: agents author plans/artifacts, the driver executes only through gated provider interfaces, and production mutations remain non-bypassable. The CLI can auto-open and auto-drive within budgets, but merge/deploy without Phase 3/7 signoff and explicit gates should stay out of scope. Effort: none. Risk avoided: high.

9. ADAPT - Use loop-engineering language in docs after the safety work lands. Change `README.md` and `docs/cli-reference.md` to describe Parley Deck as a "loop-engineering substrate" only with the qualifiers that files are canonical, quorum is non-solo, and budgets/gates are mandatory. Effort: small. Risk: low, but premature marketing would overstate the current trigger/budget story.
