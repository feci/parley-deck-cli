# Organizer brief — meta-protocol-change-quota-auto-exclude

You are **claude-1, the ORGANIZER (facilitator)** of a new Parley Deck idea, launched headlessly from the
owner's Claude Code session (the relay). That session is not a participant and writes no protocol content.
Follow the owner's global rule in `~/.claude/CLAUDE.md` ("Parley Deck — global roster and organizer"),
except where the owner's quorum decision below overrides it for this idea.

This session covers **design only: Phase 0 through Phase 4 (FINAL)**. Then you write the owner proposal note
and exit. The owner's standing rule is that every parley improvement is shown to them before
implementation and release, so do not implement, release or merge anything.

## Workspace

- CLI worktree (branch `quota-auto-exclude`, base origin/main a8634cc = v1.50.0 plus the designated-implementer
  delivery records): `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/quota-auto-exclude`.
  The deck is its `parley-deck/` (protocolRole `source`). Run every command from this worktree.
- Skill worktree (branch `quota-auto-exclude`, base origin/main a5664d8 = skill 2.14.0):
  `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/quota-auto-exclude-skill`. FINAL must name the
  skill-side changes; nobody edits the skill in this session.
- Evidence file: `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/runs-handoff/quota-auto-exclude-2026-10-03/EVIDENCE-quota-incidents.md`.
  Copy it and this brief into `parley-deck/ideas/<slug>/source-context/` so every participant reads them.

## The owner's words (verbatim, Slovak) and translation

> "nastartuj /parley-deck a navrhni upravu, resp. malu zmenu, kde vyradis agentov z rosteru automaticky,
> ak nemaju dost tokenov, resp. kreditov, min. roster je 2"

"Start /parley-deck and propose an adjustment, or rather a small change, where you drop agents from the
roster automatically if they do not have enough tokens, or rather credits. The minimum roster is 2."

## Quorum for THIS idea — owner decision (2026-10-03)

Today's preflight passed only claude-1 and codex-1. zcode-1 returns HTTP 429 "Weekly/Monthly Limit Exhausted"
until 2026-10-05 06:14 local. kimi-1's weekly quota is exhausted until about 2026-10-03, and its probe hung
with no output (see the evidence file). The relay asked the owner, and the owner chose this:

> Question: "Who should form the quorum for this idea while kimi-1 and zcode-1 have no credits?"
> Selected: **"codex-1 + claude-1 (Recommended)"**, described as "Same as in librade on 01.10: claude-1 joins
> as a separate participant process, and the organizer stays a separate claude-1 process. kimi-1 and zcode-1
> are dropped for this idea only and can rejoin when their quota returns. The role concentration is recorded
> under §15.5. Starts right away."

So, in `00-prompt.md`:

- `author: user`, `facilitator: claude-1`, `facilitator_participates: true`,
  `participants: [codex-1, claude-1]`, `track: deliberation` (§4.0: a protocol change forces deliberation;
  the owner's "small change" constrains the size of the diff, not the track), `auto_implement: false`,
  `require_model_diversity: true`.
- `excluded:` records both agents with their reason and `confirmed 2026-10-03`. kimi-1: weekly quota
  exhausted, probe hung. zcode-1: 429 Weekly/Monthly Limit Exhausted, reset 2026-10-05 06:14. Quote the owner
  decision above. Record the claude-1 role concentration under §15.5.
- **Process separation is binding.** You, the organizer, never write participant content: no round file, no
  signoff, no FINAL and no review. Every claude-1 participant artifact comes from a separately launched
  claude-1 process, through the driver or claude-1's configured invocation, each with its own fresh context.
  Do not pass your own opinions to it beyond the neutral kickoff material that codex-1 also receives.
- Keep the quorum stable through FINAL. Do not re-include kimi-1 or zcode-1 during design, even if they come
  back. The owner pre-authorized their rejoining, so a later implementation or review phase may add them, for
  example as extra model-diverse reviewers. Record that as a next-phase note in the proposal.
- If codex-1 or the claude-1 participant hits a quota or credit failure during this run, the current protocol
  still applies, because the new rule is not in force yet. Write a blocking
  `inbox/claude-1-to-user_<slug>_quota.md` with the verbatim provider error and stop launching. Do not spin
  retries. Dropping either agent would leave one participant, below the owner's minimum of 2.

## Facts to verify at HEAD before relying on them (§15; these are locators, not conclusions)

- `parley-deck/COOPERATION.md` §9.0 (~L876-893): excluding an unavailable agent requires **explicit user
  confirmation**, is per-idea and temporary, and the agent stays in the §2 roster. Re-including also needs
  confirmation. Excluding the last non-facilitator needs the §1 solo exception. Quorum locks once Phase 0
  completes, and a mid-idea unavailability falls to §5 and the watchdog.
- §5 (~L743-755): quorum locks at Phase 0. "a mid-idea unavailability does not silently shrink quorum". A
  valid idea normally has at least two active participants (~L751). The inactive-for-more-than-2-rounds drop
  needs a deadline and a ping first (~L753). Two-participant ideas: both must sign (~L754).
- §4.0 per-track table: reviewer counts by track. With two participants, standard's "2 reviewers" degrades
  to 1. Phase 5 designated-implementer availability rules (~L449): a per-idea designee failing the ping gates
  behind three recorded exits, while a global-default designee falls through with a notice.
- `internal/app/preflight_liveness.go` ~L33-36: `ClassProviderFailure` reads "Environmental, not the agent
  being dead. **Blocking but never an exclusion.**" The `providerFailureRules` at ~L627-635 cover rate-limit,
  auth, billing, overloaded and model-not-found.
- `internal/runner/failclass.go` ~L23-46: the run-level classifier, in which order matters. Check how the
  four observed exhaustion texts in the evidence file classify: 429 with `reset_at`/`retry_after`, 503
  "Unavailable (reset after 5h 51m 11s)", 403 "weekly usage limit", and a 400 or a silent hang from a disabled
  gateway connection. Also `internal/telemetry/usage.go` ~L252.
- `parley preflight --yes` "confirm excluding unavailable agents (records the exclusion)". Find where and
  how it records.
- `~/.parley/agents.toml` `[defaults] roster_change_policy = "confirm-breaking"` ("auto-add new agents;
  confirm drops/breaking changes"). The new behaviour must state how it relates to this knob.
- Frozen run roster: `parley continue` uses run.Participants plus the manifest RosterSnapshot
  (`internal/app/app.go` runContinue). This matters for any mid-run exclusion.

## The design question

How should Parley drop an agent that has run out of tokens or credits automatically, without asking the
owner, while never going below a minimum of two participants? Participants decide. FINAL must answer at
least these:

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

Non-goals: no model or effort changes, no provider, gateway or credential changes, no quota purchase, no
re-enabling of disabled gateway connections, no change to the 1.49.0 organizer rules or the 1.50.0
implementer chain beyond what this rule needs, and no unrelated code.

## How to run it (lean; use the 1.49.0+ tooling)

- Run `parley preflight` first and record its result. Use the driver: `parley run` / `parley continue`,
  `parley wait` (blocking, under 25 minutes per call), `parley status --json`, and `parley organizer brief`
  to re-orient after a compaction. Do not hand-write participant prompts when the runner can launch them,
  and record every driver gap you hit in `organizer-notes.md`. Deliberation timeout is about 30 minutes per
  agent.
- You do not implement or verify code. Participants verify every locator above and quote what they find.
  Keep your own tool outputs small. Record your usage at each phase boundary in
  `parley-deck/ideas/<slug>/organizer-usage.md`.
- Cross-review is capped at 3 rounds after round 1 (§4.0 deliberation). If the participants still disagree,
  escalate to the owner with a blocking inbox note that states both positions. Never resolve by count (§15.3).
- English in every artifact and commit. Commit on branch `quota-auto-exclude` with the prefix
  `[claude-1] meta-protocol-change-quota-auto-exclude: ...`; participants use their own prefix. Do not push
  or merge. Never copy secrets or tokens. Never open or drive Google Chrome; ego-browser is the only
  permitted browser. Do not touch the windows-portability worktree or run.

## When you finish (after FINAL is published and closed)

Write `parley-deck/inbox/claude-1-to-user_meta-protocol-change-quota-auto-exclude_proposal.md` with
frontmatter `blocking: yes` and `phase: final`. It holds the owner's approval request in plain English:

- the rule in about ten lines;
- what would have happened in each incident from the evidence file;
- the exact protocol changes as a hunk summary, and the CLI and skill changes;
- risks and their mitigations;
- every open owner decision, each with the participants' recommendation;
- the implementation size and plan, the implementer (the global default implementer is codex-1), the
  reviewers, and whether kimi-1 or zcode-1 should rejoin for review;
- the organizer-usage ledger.

Then exit.
