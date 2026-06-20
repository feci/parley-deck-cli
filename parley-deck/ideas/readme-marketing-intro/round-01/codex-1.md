---
agent: codex-1
idea: readme-marketing-intro
round: 1
date: 2026-06-20
---

## Summary

The intro should sell Parley Deck as a protocol-backed workflow, not just a CLI wrapper. The technically accurate framing is: file-backed multi-agent cooperation with explicit phases, per-agent artifacts, consensus gates, supervised automation, and tooling that helps run and resume the protocol. Avoid claims that imply ungated autonomy, guaranteed quality, or provider-specific execution beyond what the protocol and shipped commands actually support.

## Proposed approach

Put this immediately after `# parley-deck-cli`, before `## Install`:

```markdown
Parley Deck turns multi-agent AI work into a durable, reviewable process. Instead
of losing decisions in a chat transcript, each agent writes its own canonical
artifact, the group moves through explicit rounds, and the final design plus
implementation record live in your repository.

Use it when you want agents to disagree productively, converge before code is
written, and leave enough context for a fresh human or model to resume the work
later. The `parley` CLI provides the local TUI, agent discovery, readiness checks,
round orchestration, consensus/signoff helpers, advisory consults, retrospective
scans, and supervised pipeline machinery around the protocol.

What it includes:

- Transport choices for local files, GitHub PRs, or GitLab MRs.
- An eight-phase lifecycle: kickoff, independent analysis, cross-review,
  consensus, finalization, implementation, code review, and fix-up.
- Non-solo cooperation by design: stable agent IDs, one artifact per agent per
  round, and no editing another agent's file.
- Self-contained design output in `FINAL.md` plus a living `IMPLEMENTATION.md`
  for progress, decisions, surprises, validation, and outcomes.
- Pre-idea readiness checks with protocol freshness and roster liveness before a
  new idea starts.
- A live TUI and auto-drive support for advancing protocol phases, with code
  implementation and side effects gated by the protocol instead of assumed.
- Optional pipeline blocks for staged work, gates, idempotency, reconciliation,
  and durable effect records.
- Advisory retrospectives over the deck's own history, so process improvements
  are proposed through the same quorum-gated workflow rather than applied
  automatically.

Parley Deck is adopted and adapted from good prior art: OpenRouter Fusion's
compare-not-merge synthesis, OpenAI ExecPlans/PLANS.md-style resumable plans,
RHO's retrospective optimization loop, kindly's review/supervision discipline,
and preflight-style readiness checks. The twist is the repository-backed quorum:
files are canonical, automation is supervised first, and consensus beats vibes.
```

Feature claim map:

- "Transport choices for local files, GitHub PRs, or GitLab MRs" maps directly to COOPERATION.md Section 0 and Section 11.
- "Eight-phase lifecycle" maps to COOPERATION.md Section 4. Use the eight named phases exactly; do not invent additional phases in the marketing bullets.
- "Non-solo cooperation", "stable agent IDs", and "one artifact per agent per round" map to Sections 1, 2, and 6. This is protocol-backed, but still depends on the active runner/participants actually producing files.
- "`FINAL.md` plus `IMPLEMENTATION.md`" maps to Sections 4 and the prompt's ground truth. This is safe to market as resumability support, not as a guarantee that every old run is fully recoverable.
- "Pre-idea readiness checks" maps to Section 9.0 and shipped `parley preflight`. Include liveness/freshness, but do not imply unavailable agents are silently dropped; exclusion/re-inclusion is user-confirmed.
- "Live TUI and auto-drive" maps to shipped `parley tui`, `parley run`, `parley resume`, and driver code. Important wording: auto-drive advances protocol phases; implementation/fix-up automation is gated by `auto_implement`, `--no-implement`, review state, and protocol gates.
- "Pipeline blocks" maps to Section 12 and shipped `parley pipeline`. Keep the wording provider-agnostic; local CLI agents author markdown, and the driver is the actor for approved side effects.
- "Advisory retrospectives" maps to Section 13 and shipped `parley retro`. The safe claim is "scan/select/diagnose/propose"; `retro propose` may scaffold one new `00-prompt.md`, but it does not apply improvements.
- "Advisory consults" can be mentioned because shipped `parley consult` exists, but it must be clearly non-canonical and not quorum evidence.

For the adapted `parley-deck-skill` README, keep the same protocol story but swap the CLI paragraph:

```markdown
The `parley-deck-skill` installer puts the Parley Deck cooperation instructions
and fallback protocol snapshot into supported agent runtimes, then helps check
and sync project metadata. It complements the `parley` CLI: the skill teaches
agents how to participate; the CLI orchestrates runs, readiness checks, TUI
views, consensus helpers, pipelines, and retrospectives.
```

Only claim the installed skill command's own shipped surfaces: `install`, `doctor`, `status`, `sync-project`, `uninstall`, `paths`, and `--version`. Do not say `parley-deck-skill` runs rounds, opens the TUI, runs `retro`, or launches pipelines unless those commands move into that installer.

## Concerns / open questions

- The current `docs/cli-reference.md` is behind `internal/app/app.go` for `preflight`, `retro`, `consult`, `pipeline`, and default auto-drive. If README points users to the CLI reference as exhaustive, either update that doc in the implementation or phrase the intro as "see `parley help`" for the full current command list.
- "Autonomy" needs strict wording. The protocol is supervised-first; block boundaries gate by default, production mutations are non-bypassable, and code-writing implementation only auto-runs for opted-in ideas.
- "Agent supervision" is real: first-output watchdog, stall guard, heartbeats, failure classification, stdout fallback, and artifact-wins behavior are implemented. But market it as run supervision and recovery signals, not as reliability guarantees.
- "Comparison & blind spots" is a consensus drafting lens, not an automatic truth engine. Avoid "finds every blind spot" or "proves the best answer".
- The inspirations paragraph should say "adopted and adapted" and name the deck's change. Do not imply Parley Deck invented compare-not-merge synthesis, resumable plans, retrospective optimization, review dispositions, or preflight readiness.

## Risks

- Over-claiming automation is the largest risk. Phrases like "fully autonomous delivery", "self-healing pipelines", or "agents deploy for you" would contradict Sections 12.4, 12.8, and the `auto_implement` gate.
- Over-claiming enforcement is also risky. Some protocol rules, especially round-01 independence across PR surfaces, are social/procedural unless the chosen transport and runner setup enforce isolation.
- Provider-specific wording can age badly. Keep provider examples out of the top intro unless a specific shipped provider surface is documented and intended as public.
- The CLI and skill are separate release artifacts. The README pair should not blur `parley` commands with `parley-deck-skill` installer commands.
- If the README intro says "any CLI/model", qualify it as vendor/model-agnostic roster support, subject to installed commands, auth, sandbox permissions, and `parley agents verify` results.
