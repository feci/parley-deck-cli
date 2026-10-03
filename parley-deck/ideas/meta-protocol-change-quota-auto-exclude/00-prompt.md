---
idea: meta-protocol-change-quota-auto-exclude
author: user
facilitator: claude-1
facilitator_participates: true
created: 2026-10-03
track: deliberation
participants: [codex-1, claude-1]
excluded: kimi-1 — weekly Kimi quota exhausted until about 2026-10-03, hosted-PONG probe hung with no output (preflight class deadline-after-output) — confirmed 2026-10-03
excluded: zcode-1 — HTTP 429 Weekly/Monthly Limit Exhausted, reset 2026-10-05 06:14 local (preflight class provider-failure:rate-limit) — confirmed 2026-10-03
auto_implement: false
require_model_diversity: true
status: round-01
---

## Problem / idea

How should Parley drop an agent that has run out of tokens or credits automatically, without asking
the owner, while never going below a minimum of two participants? Design the rule (a protocol change
under §7) and the smallest CLI and skill changes it needs.

This run covers **design only: Phase 0 through Phase 4 (FINAL)**. Nothing is implemented, released or
merged in this run. The owner's standing rule is that every parley improvement is shown to them before
implementation and release.

Two source files are copied verbatim into this idea. **Every participant must read both in full:**

- `source-context/ORGANIZER-BRIEF.md`: the owner's controlling brief for this idea, as relayed. It holds
  the design question, the locators to verify and the constraints.
- `source-context/EVIDENCE-quota-incidents.md`: five quota-exhaustion incidents with provenance labels
  (PRIMARY, SECONDARY, RECALL).

A third file, `source-context/preflight-2026-10-03T1223.json`, is the raw readiness output summarized
below.

The brief's "Facts to verify at HEAD" are **locators, not conclusions**, and the evidence file marks
its own hypotheses as such. Verify each claim at HEAD and quote what you find (§15).

## The owner's words (verbatim, Slovak) and translation

> "nastartuj /parley-deck a navrhni upravu, resp. malu zmenu, kde vyradis agentov z rosteru automaticky,
> ak nemaju dost tokenov, resp. kreditov, min. roster je 2"

"Start /parley-deck and propose an adjustment, or rather a small change, where you drop agents from the
roster automatically if they do not have enough tokens, or rather credits. The minimum roster is 2."

<!-- Original language: Slovak. The translation is the relay session's, copied from source-context/ORGANIZER-BRIEF.md. -->

## The design question (verbatim from the brief)

Participants decide. FINAL must answer at least these:

1. **Scope.** Pre-idea (§9.0) only, or also mid-idea after the quorum locked? Four of the five incidents in
   the evidence file were mid-idea.
2. **What counts as "not enough tokens/credits".** Which provider signals qualify: usage or credit limit
   exhausted, 429 with a reset window, 403 weekly limit, 503 "reset after". How to tell a window exhaustion
   (hours or days) from a transient throttle (seconds), for example with a retry-after or reset threshold.
   Whether the classifier gap in the evidence (503 classified as overloaded) must be fixed. How to treat
   disguised exhaustion (a 400 or a hang from a disabled connection). **Fail closed:** anything not positively
   classified as exhaustion takes today's user-confirmed path, never the automatic one.
3. **What "roster" means here.** Per-idea quorum exclusion (temporary, re-probed at the next idea) versus
   changing the machine roster in `~/.parley/agents.toml`. Say what happens when the quota resets, and
   whether re-inclusion is automatic or confirmed, using the recorded `reset_at` where it exists.
4. **The minimum of 2.** Count only non-facilitator participants. Decide what happens when automatic
   exclusion would go below 2 (stop and escalate, with the §1 solo exception unchanged). Cover two agents
   exhausting at once, and whether the result depends on order.
5. **Interactions.** The designated implementer running out mid-Phase 5. Reviewer counts per track and
   `require_model_diversity` after an exclusion. Signoffs, rounds and open findings already filed by the
   excluded agent: an excluded reviewer's open CRITICAL or MAJOR finding must not vanish by exclusion, and
   §15.3 forbids resolution by count. The drafter. In-flight artifacts.
6. **Integrity.** Exclusion must be triggered only by a machine-classified provider error with its evidence
   recorded, never by content, disagreement or slowness alone. It must not become a way to remove a
   dissenter.
7. **Recording and notice.** Where the automatic exclusion is recorded (an `excluded:` entry with an
   automatic marker, the run manifest, a non-blocking owner inbox notice), and what `parley status`,
   `parley wait` and the organizer brief report.
8. **Configuration.** Default-on or opt-in. A knob name, and how it relates to
   `roster_change_policy = "confirm-breaking"`. Whether the floor of 2 is fixed or configurable, with 2 as the
   lower bound.
9. **Size.** The owner asked for a *small* change. Prefer reusing the existing classifiers, the preflight
   exclusion record and the §5/§9.0 text over any new subsystem. FINAL lists the exact protocol hunks, all
   three COOPERATION.md copies (two in the CLI under the drift guard, plus the skill's bundled
   `skills/parley-deck/references/COOPERATION.md`), the CLI files and the tests, and states what is
   deliberately left out.

## Quorum for this idea: owner decision (2026-10-03)

Readiness (below) passed only claude-1 and codex-1. The relay asked the owner, who answered (verbatim, as
relayed in the brief):

> Question: "Who should form the quorum for this idea while kimi-1 and zcode-1 have no credits?"
> Selected: **"codex-1 + claude-1 (Recommended)"**, described as "Same as in librade on 01.10: claude-1 joins
> as a separate participant process, and the organizer stays a separate claude-1 process. kimi-1 and zcode-1
> are dropped for this idea only and can rejoin when their quota returns. The role concentration is recorded
> under §15.5. Starts right away."

- **Quorum:** codex-1 and claude-1. kimi-1 and zcode-1 are excluded **for this idea only** under §9.0, with
  the owner's explicit confirmation above (frontmatter `excluded:`). Both stay in the §2 roster and are
  re-probed at the next idea. Preflight did not record these exclusions (it was not run with `--yes`); the
  owner decision above is the §9.0 confirmation, written here in the CLI's own `excluded:` line format.
- **Stable through FINAL.** kimi-1 and zcode-1 are not re-included during design, even if their quota
  returns. The owner pre-authorized their rejoining, so a later implementation or review phase may add them.
- **Role concentration (§15.5).** claude-1 is both the declared facilitator and a participant, in separate
  processes. The facilitator is the organizer session. The participant is a separately launched headless
  claude-1 process with its own fresh context at each launch. The organizer writes no participant content:
  no round file, signoff, consensus content, FINAL or review. Its procedural calls are provisional until the
  signoff gate passes. If the claude-1 participant drafts `consensus.md`, §15.5 requires the one-line
  role-concentration record and a `## Drafter position changes` section.
- **Model diversity.** The two participants run different model families: codex-1 runs OpenAI `gpt-6-astra`
  and claude-1 runs Anthropic `claude/claude-opus-5-5[1m]`, both at effort `max`.
- **Quota failures during this run.** The current protocol governs this run, not the rule under design. If
  codex-1 or the claude-1 participant hits a quota or credit failure, the organizer writes a blocking
  `inbox/claude-1-to-user_meta-protocol-change-quota-auto-exclude_quota.md` with the verbatim provider error
  and stops launching. Dropping either agent would leave one participant, below the owner's minimum of 2.

## Readiness (§9.0)

`parley preflight --dir . --json`, run 2026-10-03 12:23 CEST from this worktree on CLI 1.50.0. It took
94 s and exited 3 (pending gates). Raw output: `source-context/preflight-2026-10-03T1223.json`.

| rosterId | available | class | reason | CLI version |
|---|---|---|---|---|
| claude-1 | true | ready | | 2.1.288 (Claude Code) |
| codex-1 | true | ready | | codex-cli 0.159.3 |
| kimi-1 | false | deadline-after-output | deadline-after-output | 0.42.0 |
| zcode-1 | false | provider-failure | provider-failure:rate-limit | zcode-app-cli 3.7.7-13 |

Gates (verbatim):

- `kimi-1 readiness is unresolved (deadline-after-output) — investigate and re-check; the agent is not excluded`
- `zcode-1 reports a provider failure (provider-failure:rate-limit) — resolve the provider/auth state; the agent is not excluded`
- An unrelated historical gate: `ideas/meta-protocol-change-devx-speed: idea frontmatter declares
  facilitator: claude-1 but also lists it in participants: without facilitator_participates: true`. That
  idea is left untouched. The designated-implementer run recorded the same gate.

Protocol freshness is `source-advisory` (protocolRole `source`, advisory only, never writes
COOPERATION.md). The skill installer and both runtime skill markers are at 2.14.0.

## Workspace, transport and commits

- **CLI worktree** (all commands run here):
  `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/quota-auto-exclude`. It is on branch
  `quota-auto-exclude`, based on origin/main a8634cc (v1.50.0 plus the designated-implementer delivery
  records). This deck is its `parley-deck/`, with protocolRole `source`.
- **Skill worktree** (read-only in this run): `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/quota-auto-exclude-skill`.
  It is on branch `quota-auto-exclude`, based on origin/main a5664d8 (skill 2.14.0). FINAL must name the
  skill-side changes. Nobody edits the skill in this run.
- **Transport.** The live header names `github-pr`. Per the owner brief, this run uses local canonical files
  committed on branch `quota-auto-exclude`, with no PR, push or merge. That is an explicit per-run override
  of the §11.B mechanics, as in the designated-implementer run. The global transport header is unchanged.
- **Commits.** Participants write their own files and do not run `git commit`, `git push` or `git merge`.
  The organizer sweeps each participant's files into a commit under that participant's prefix,
  `[<agent-id>] meta-protocol-change-quota-auto-exclude: ...`, marked as swept by the organizer. Orchestration
  files use `[claude-1] meta-protocol-change-quota-auto-exclude: ...`.

## Constraints

- **Design only**, Phase 0 through Phase 4. No implementation, release, publication or merge. Nobody edits
  the CLI code, the skill or any COOPERATION.md copy in this run. FINAL describes those changes.
- **Track: deliberation.** §4.0 forces it for a protocol change (§7). The owner's "small change" constrains
  the size of the diff, not the track. Cross-review is capped at 3 rounds after round 1. If the participants
  still disagree, the question goes to the owner with both positions stated. Never resolve by count (§15.3).
- **Verification integrity (§15).** Every claim about current behavior carries a HEAD locator and a quote.
  Brief and evidence statements are testimony until verified. Round 1 needs a non-empty
  `## Existing alternatives` section with that exact heading (§15.6a).
- **Fail closed.** Anything not positively classified as exhaustion keeps today's user-confirmed path.
- **Prefer reuse** of the existing classifiers, the preflight exclusion record and the §5/§9.0 text over any
  new subsystem.
- Record your protocol attestation (`context_mode`, `source_sha256`, `packet_sha256`, `fallback_reason`) in
  your own artifact.
- English in every artifact. No secrets or tokens. Never open or drive Google Chrome; ego-browser is the only
  permitted browser. Do not touch the windows-portability worktree or its run.

## Non-goals (verbatim from the brief)

No model or effort changes, no provider, gateway or credential changes, no quota purchase, no
re-enabling of disabled gateway connections, no change to the 1.49.0 organizer rules or the 1.50.0
implementer chain beyond what this rule needs, and no unrelated code.

## Kickoff metadata (facilitator)

- Organizer protocol context: `parley protocol packet --phase 0 --track deliberation --audience facilitator
  --flag protocol_change`, giving context_mode=packet,
  source_sha256=b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388,
  packet_sha256=0744ea2200b578bdfbac56cb2e9d4c4ce4e736a4ef28f8e0f92519880a29b5ca, with fallback_reason absent.
  The omitted §3 and §15.5 were read from the full source on demand.
- Shared memory (OpenViking through MCPAnywhere) is not connected in the organizer session, so local sources
  are used.
- `parley run` mints its own timestamped slug and cannot adopt a pre-written kickoff. The launch path and
  every driver gap are recorded in `organizer-notes.md`.
