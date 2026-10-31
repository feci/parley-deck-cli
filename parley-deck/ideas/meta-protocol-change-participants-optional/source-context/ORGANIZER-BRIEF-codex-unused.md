# Organizer brief: meta-protocol-change-participants-optional

You are **codex-1, the organizer and the implementer** of a new Parley Deck idea. The owner wants Claude
tokens saved: codex-1 organizes and implements, the owner's Claude Code session only relays, and claude-1
takes no part.

## Workspace

- CLI worktree: `/Users/tomasfecko/AI_WORKSPACE/parley-deck/worktrees/participants-optional` (branch
  `participants-optional`, from origin/main `604ccfc`, CLI 1.56.0). The deck is its `parley-deck/`.
- Skill worktree: `/Users/tomasfecko/AI_WORKSPACE/parley-deck/worktrees/participants-optional-skill`
  (branch `participants-optional`, from origin/main `7cdfb24`, skill 2.20.0).
- `codex` is at `/Users/tomasfecko/.nvm/versions/node/v24.15.0/bin/codex`. `kimi` is at
  `~/.kimi-code/bin/kimi` (not on PATH). `zcode` is on PATH.
- Never prune or declare worktrees.

## The owner's words (verbatim, Slovak) and translation

The owner set this rule on 2026-10-31 in another session (project librade-chart-analytics, file
`librade-chart-analytics/docs/decisions/2026-10-31-parley-participants-optional.md`):

> "ak hocijaky agent nema tokeny hned ho vyluc, pokial existuje organizer alebo implementer, tak moze parley
> pokracovat, participanti su volitelny"

Translation: "If any agent has no tokens, exclude it at once. While an organizer or an implementer exists,
Parley can continue. Participants are optional."

On 2026-10-31 the relay asked the owner whether to make this rule a protocol change. The owner selected
"Áno, nová idea (Recommended)" ("Yes, a new idea"). That option said: the floor falls to the organizer or the
implementer, participants become optional, the idea decides who reviews when the participants drop out (for
example an organizer with a different model), and the release is CLI 1.57.0 and skill 2.21.0 on all channels.

Earlier binding owner rulings stay in force (shipped in 1.56.0, idea
`meta-protocol-change-credit-exhaustion-failover`):

- An agent whose provider is out of credits or tokens is excluded at once, with no retry and no reconnect.
- The organizer drops out: the implementer takes over the organizer role. The implementer drops out: the
  organizer takes over the implementer role. Both drop out: the first remaining participant takes over both
  roles. A participant drops out: no succession.

## Current state (verify at HEAD, §15)

Read core 2.20.0 §9.0 "Automatic exclusion" (the CLI worktree's `parley-deck/COOPERATION.md`) and the ideas
`meta-protocol-change-credit-exhaustion-failover`, `meta-protocol-change-participant-dropout-review-gate-timing`
and `meta-protocol-change-participant-dropout`. Do not redo their work.

The gap: 1.56.0 keeps the floor "two distinct usable non-facilitators". Below the floor, the run escalates and
waits for the owner. The owner wants the run to continue while an organizer or an implementer is usable.

## Design target (participants decide the details in FINAL)

1. **New floor.** The run continues while at least one usable agent holds the organizer or the implementer
   role (after the shipped succession). Participants are optional. Decide:
   - whether the rule also lets an idea **start** with no participant, or only covers participants that drop
     out (the owner's words "participanti su volitelny" suggest both; say what you choose and why);
   - what happens when nobody can hold the organizer or the implementer role: then escalate once.
2. **Review without participants.** Today a non-implementer reviews, goal-checks and signs. Decide what
   replaces the independent review when no independent participant remains. Candidates:
   - the organizer reviews when it is a different agent from the implementer (prefer a different model);
   - when one agent holds both roles, close with a recorded "no independent review" marker, the required
     checks (tests, `checks:`) green, and a clear statement of the reduced assurance in IMPLEMENTATION.md and
     the released note;
   - a later independent review when a participant returns at the next idea.

   Keep §15 honest: an agent never gives a verdict on its own work. A missing review is recorded as missing,
   never as a pass.
3. **Design rounds and consensus with fewer agents.** Say how rounds, consensus and FINAL work with one or two
   agents. The skill's non-solo requirement says a run with no other agent stops unless the owner authorizes a
   solo exception. State that this owner rule is a standing, recorded authorization for this case, and give its
   exact bounds.
4. **What stays.** Filed dissent, ❌ and findings of a dropped agent stay binding. A drop never removes dissent.
   The shipped exhaustion recognition, the no-retry rule, the succession order and the transient two-attempt
   path do not change.
5. **Size.** The core skill `SKILL.md` is 19,989 of 20,000 bytes (11 bytes free). The phase packets are 68,448
   and 69,880 of 70,000 bytes. Do not raise the limits. Make room by compaction without the loss of a rule.
   List the exact protocol hunks for all three COOPERATION.md copies, plus the CLI files, the skill files and the
   tests.
6. **Writing standard.** Write the new protocol prose and all artifacts in ASD-STE100 (Appendix B).

## Quorum and roles

- `00-prompt.md`: `author: user`, no `facilitator:` field (codex-1 facilitates as an ordinary participant),
  `participants: [codex-1, kimi-1, zcode-1]`, `implementer: codex-1`, `track: deliberation` (a protocol
  change), `auto_implement: true`, `require_model_diversity: true`. Quote the owner.
- Run `parley preflight`. If an agent fails, the shipped exclusion and dropout rules apply. Never touch
  credentials, keys or gateway settings.
- kimi-1 and zcode-1 are the independent reviewers and goal checkers, in separate processes with focused
  briefs. You never review your own code.
- **Close authority.** Normal gates apply. If one reviewer drops, the shipped single-reviewer exception applies.
  If it does not apply mechanically, the owner pre-authorizes the close when all of these hold:
  - the final review of the remaining reviewer has no open CRITICAL or MAJOR;
  - all current participants filed a review-consensus ACCEPT;
  - a fresh goal check by a reviewer process passed;
  - the current-tree AC evidence is recorded.

## Process, retries and release

- **No proposal stop.** The owner said: make the new version and release it on all channels. Go from design
  through implementation to release without stopping. Stop only if an owner decision is truly open. In that
  case, write one blocking `codex-1-to-user_<slug>_decision.md` with the options and a recommendation.
- **Driver and fallback.** Try the driver (`parley run`, `continue`, `wait`, `status`) once at each phase. The
  owner authorized native CLI fallback for recorded driver refusals in the previous idea; the same authority
  applies here. If a gate refuses, use native CLI processes, record the gap with its exact message in
  organizer-notes.md, and continue. Do not declare runs or worktrees. Do not apply legacy accounting.
- **Fix-up.** At most 5 cycles. If a CRITICAL or MAJOR is still unresolved after cycle 5, escalate.
- **Retries.**
  - A provider 429 or 503, or a timeout: wait 15 minutes and relaunch, at most 8 times per step.
  - A timeout: relaunch with 2400 s, then with 3600 s.
  - An auth or credit error on codex-1: stop with a blocking note.
- **Host reboot.** If the relay restarts you after a host reboot, a killed child is not a participant failure.
- **Waiting.** Never end your turn while you wait. Wait only with blocking tool calls. End only when the work
  is done, or deliberately after a blocking owner note.
- **Conventions.** English everywhere, in STE. Commit prefix
  `[codex-1] meta-protocol-change-participants-optional: ...`. Never copy secrets. Never drive Google Chrome.
- **Release.** Merge to `main`. Then:
  - Versions: CLI **1.57.0**, skill **2.21.0**, core **2.21.0**. Build core 2.21.0 on the staged
    `~/.parley/staging/COOPERATION-2.20.0.md` plus exactly this idea's reviewed hunks. Stage it as
    `~/.parley/staging/COOPERATION-2.21.0.md`.
  - GitHub releases.
  - Both Homebrew formulae.
  - A skill-only winget PR (one application per PR). The skill 2.20.0 PR #465809 may still be open: check it,
    and do not combine applications. The CLI winget stays held while Windows is experimental.
  - Install the skill into all managed runtimes and the four generic locations, with the same method and
    hash verification as the 1.56.0 release (see that idea's release receipts).
  - npm publication and `parley protocol publish` are owner-only (npm login with OTP, and a TTY gate). Do not
    try them, and do not work around the gates. The relay updates the owner script.
  - A reviewer process (zcode-1 or kimi-1) verifies all channels in a separate process.
- **Finish.** Write `inbox/codex-1-to-user_meta-protocol-change-participants-optional_released.md` with:
  - the channel evidence;
  - the still-pending owner commands, in version order: npm and core for each still-unpublished version from
    2.15.0 to 2.21.0;
  - the deferred items;
  - usage.

  Then exit.
