# Organizer brief: meta-protocol-change-participant-dropout-review-gate-timing

You are **codex-1, the organizer and the implementer** of a new Parley Deck idea. The owner wants Claude
tokens saved: codex-1 organizes and implements, the owner's Claude Code session only relays, and claude-1
takes no part.

## Workspace

- **CLI worktree:** `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/review-gate-timing`. Branch
  `review-gate-timing`, from origin/main `85c5a8f` (CLI 1.52.0, participant dropout shipped). The deck is its
  `parley-deck/`.
- **Skill worktree:** `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/review-gate-timing-skill`.
  Branch `review-gate-timing`, from `dbdb919` (skill 2.16.0).
- **Prior art.** Read it and do not redo it:
  - `parley-deck/ideas/meta-protocol-change-participant-dropout/`, especially its released note's "Limits and
    deferred work", which names this follow-up, and its organizer notes;
  - `parley-deck/ideas/meta-protocol-change-quota-auto-exclude/`.

## The owner's words (verbatim, Slovak) and translation

Earlier rule (2026-10-08):

> "organizator a implementator su vacsinou codex a claude, takze ucastnici su kimi a zcode, ak vypadnu z
> hocikakeho dovodu, tak v tej idei uz nepokracuju a mozu sa zapojit az do dalsej"

"The organizer and implementer are usually codex and claude, so the participants are kimi and zcode. If they
drop out for any reason, they do not continue in that idea and can join only the next one."

Also: "participanti nie su nevyhnutne potrebny obaja, staci jeden a to by nemalo zaseknut parley-deck"
("the participants are not necessarily both needed; one is enough, and that should not stall parley-deck").

Today (2026-10-09), answering whether to start this follow-up:

> "ano spusti ideu ako poriesit ze participant vypadne v strede idei, ved tam mame nejake timeouts, takze
> kludne rozved participant-dropout-review-gate-timing"

"Yes, start the idea on how to handle a participant dropping out in the middle of an idea. We do have some
timeouts there, so feel free to expand participant-dropout-review-gate-timing."

## The problem (verify at HEAD, §15)

1.52.0 drops a failed non-protected participant from the idea. The 1.52.0 released note itself says:

> "Existing reviewer, diversity, strict and goal-check gates still apply: **auto_implement 3→2 still blocks
> with only one independent reviewer**. This release supplies no automatic waiver."

In the owner's standard setup (organizer, implementer codex-1, participants kimi-1 and zcode-1), one dropout
mid-idea leaves zcode-1 as the only independent reviewer. The LE-7/LE-11 auto-close then refuses, and the run
waits for the owner. That is exactly the stall the owner wants gone.

## Design target. Participants decide the details in FINAL.

1. **Review gate after a recorded dropout.** When the independent-reviewer count falls to 1 **only because of
   a recorded dropout or exclusion**, one model-diverse, non-implementer reviewer is enough for the review
   round, the review consensus, the goal check and the auto-close. Decide:
   - whether the goal check may be a fresh process of the same remaining reviewer;
   - how `require_model_diversity` and `strict_gate` read in this case;
   - that the recorded reason is checked, not inferred.

   Without a recorded dropout, nothing changes.
2. **Mid-idea timing.** Use the timeouts that already exist (per-track ceilings, watchdog classes
   `no_first_output` / `stalled` / `timeout`) so that a hung participant is detected and dropped promptly,
   not after a long silent wait.
   - **Observed:** kimi hung with no output for 7+ minutes on a one-word probe.
   - **Decide:** how a no-output or stalled watchdog verdict feeds the 1.52.0 one-retry-then-drop path, and
     the bounds.
   - **Decide:** what happens to a participant whose failure lands in the middle of a round, a consensus or
     a signoff batch. For example, consensus proceeds with the remaining signers, and an earlier filed ❌
     stays binding.
3. **Floor and protected roles stay as shipped.** The floor is the implementer plus at least one other
   participant, and the organizer never counts. Below the floor, escalate once. Removing a dissenter must stay
   impossible.
4. **Compaction is mandatory.** Per the 1.52.0 note, the phase-1 packet is at 69,963 of 70,000 bytes and the
   core `SKILL.md` at 19,995 of 20,000 bytes. This idea must make room by compacting existing text, without
   losing a rule, and must not raise the limits.
5. **Size.** Keep it small. List the exact protocol hunks for all three COOPERATION.md copies, the CLI files
   and the tests.

## Quorum and roles for THIS idea

- **Write `00-prompt.md` with:**
  - `author: user`;
  - no `facilitator:` field (codex-1 facilitates as an ordinary participant);
  - `participants: [codex-1, kimi-1, zcode-1]` and `implementer: codex-1`;
  - `track: deliberation`, `auto_implement: true` and `require_model_diversity: true`.

  Quote the owner's words.
- **kimi-1 is currently broken.** The relay probed it at about 02:30 CEST:
  `403 [kimi-coding] … excluded by this API key's connection allowlist / quota scope`. That needs an
  owner-side gateway key change. Run `parley preflight`. If kimi-1 still fails, the shipped 1.52.0 rule and
  the owner's rule apply: kimi-1 is excluded for this idea, and codex-1 + zcode-1 stay at the floor. Record
  the exclusion. Do not touch any credential, key or gateway setting.
- **Implementer and reviewer:** codex-1 implements. zcode-1 is the independent reviewer, run as a separate
  process with focused review briefs. You never review your own code.
- **Dogfood the gap.** With kimi out, this run hits exactly the one-reviewer gate it is fixing. Until the fix
  ships, the close is **pre-authorized by the owner** when all of these hold:
  - the final zcode-1 review has no open CRITICAL or MAJOR;
  - both current participants' review-consensus ACCEPT blocks exist;
  - a fresh zcode-1 goal-check process passes;
  - the current-tree AC evidence is recorded.

  Record this as the owner-attended close authority, quoting the owner.

## Process

- **Design, then implementation, then release.** Go straight through without a proposal stop unless an owner
  decision remains open. If one does, write one blocking `codex-1-to-user_<slug>_decision.md` with the options
  and a recommendation.
- **Implementation.** Write the `IMPLEMENTATION.md` plan first, then implement exactly FINAL. Fix-up is capped
  at 5 cycles. If a CRITICAL or MAJOR is still unresolved after cycle 5, escalate.
- **Retries.**
  - A provider 429 or 503, or a timeout, in your own or zcode's step: wait 15 minutes and relaunch, at most
    8 times per step.
  - A timeout relaunches first with 2400 s, then with 3600 s.
  - An auth or credit error on codex-1 stops the run with a blocking note.
- **Driver.** Use the driver (`parley run`, `continue`, `wait`, `status`) and record any gap. D6 is known.
  Never prune or declare worktrees.
- **Never end your turn while you are waiting.** Wait only with blocking tool calls. End only when the work is
  done, or deliberately right after writing a blocking owner note.
- **Conventions.** English everywhere. Commit prefix:
  `[codex-1] meta-protocol-change-participant-dropout-review-gate-timing: ...`. Never copy secrets, and never
  drive Google Chrome.

## Release (after close)

- Merge to `main`. The targets are CLI **1.53.0**, skill **2.17.0** and core **2.17.0**. Build core 2.17.0 on
  the staged `~/.parley/staging/COOPERATION-2.16.0.md`, plus exactly this idea's reviewed hunks. Neither 2.15.0
  nor 2.16.0 is published yet.
- **Channels:**
  - GitHub releases;
  - both Homebrew formulae;
  - a skill-only winget PR (the CLI winget stays held while Windows is experimental);
  - the skill installed into all managed runtimes **and** the four generic locations
    `~/.hermes/profiles/{ldx,librade,testprofile}/skills/parley-deck` and
    `~/.config/opencode/skills/parley-deck`, using
    `parley-deck-skill install --target generic --dest <dir> --force`. Verify every SKILL.md by hash.
- zcode-1 verifies the channels, as a separate process with a short brief.
- **Finish** with `inbox/codex-1-to-user_meta-protocol-change-participant-dropout-review-gate-timing_released.md`.
  It must contain:
  - the channel evidence;
  - every owner-only command, in order: npm and core for 2.15.0, then 2.16.0, then 2.17.0, listing only the
    ones not yet published at that time;
  - the deferred items;
  - usage.

  Then exit.
