# Organizer brief: meta-protocol-change-participant-dropout

You are **codex-1, the organizer and the implementer** of a new Parley Deck idea. The owner wants Claude tokens
saved, so codex-1 does the organizing and the implementation. The owner's Claude Code session only relays and
checks minimally, and claude-1 takes no part.

## Workspace

- CLI worktree: `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/participant-dropout`. It is on
  branch `participant-dropout`, based on origin/main `d16ee9c` (CLI 1.51.0, quota auto-exclusion shipped). The
  deck is its `parley-deck/`.
- Skill worktree: `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/participant-dropout-skill`, on
  branch `participant-dropout`, based on `352a475` (skill 2.15.0).
- Prior art (read it, do not redo it): `parley-deck/ideas/meta-protocol-change-quota-auto-exclude/`, namely its
  FINAL, IMPLEMENTATION and the review rounds. Also `COOPERATION.md` §9.0 "Quota auto-exclusion" and §5.

## The owner's words (verbatim, Slovak) and translation

> "ok, cize je v parley deck automaticke pravidlo, ze ak niektory z participant modelov nie je dostupny tak sa
> preskoci? ak mame organizatora, implementatora a potom participantov, tak participanti nie su nevyhnutne
> potrebny obaja, staci jeden a to by nemalo zaseknut parley-deck"

"OK, so is there an automatic rule in parley deck that if one of the participant models is unavailable, it is
skipped? If we have an organizer, an implementer and then participants, the participants are not necessarily
both needed. One is enough, and that should not stall parley-deck."

The relay explained that 1.51.0 excludes automatically only on an explicit credit or quota exhaustion message
with a reset of at least 60 minutes, and only through a verified recognizer (today only zcode, possibly inert).
It also explained that a hang, a 503, a timeout or an auth error still escalates to the owner. It then asked
whether to skip an unavailable participant on other failures too, keeping the floor of 2 and never applying
it to the implementer. The owner answered:

> "organizator a implementator su vacsinou codex a claude, takze ucastnici su kimi a zcode, ak vypadnu z
> hocikakeho dovodu, tak v tej idei uz nepokracuju a mozu sa zapojit az do dalsej"

"The organizer and the implementer are usually codex and claude, so the participants are kimi and zcode. If
they drop out **for any reason**, they do not continue in that idea, and they can join only from the next one."

## The rule the owner wants (the design target)

A **non-protected participant**, typically kimi-1 or zcode-1 and never the organizer or the implementer, whose
dispatched step fails **for any reason** is dropped from **this idea only**. That covers a provider error of
any kind (429, 503, 401), a hang, a timeout, a crash or an unparseable or invalid output. The dropped agent:

- does not rejoin this idea;
- is probed again at the next idea;
- keeps every filed ❌, finding and DISPUTED claim in force. The quota rule already guarantees this.

The run continues without asking the owner, as long as the floor holds. The owner gets one non-blocking notice.

The participants must settle these points in FINAL:

1. **Reuse 1.51.0.** Broaden the trigger of the shipped mechanism rather than building a second one. It already
   has the floor, the protected roles, the batch transition record, the notice, the preserved dissent and the
   stage-1 and stage-2 paths. State what the quota-specific predicate becomes, or whether it stays as the
   stricter sub-case, for example for the stated reset.
2. **"Any reason" against a transient blip.** The owner said "for any reason". Decide the minimal retry before a
   drop, for example one relaunch or the existing timeout relaunch, so that a single network glitch does not cost
   an agent the whole idea. Keep it small and fixed, and make the owner's intent the default.
3. **The floor.** The owner's minimum is 2: the implementer plus at least one other participant, and the
   organizer never counts. If a drop would go below it (both kimi and zcode out), apply nothing and escalate
   once. FINAL should state the recommended owner options in that escalation, for example authorizing a
   separate claude or codex reviewer process. Do not choose one silently.
4. **Protected roles stay protected.** That means the implementer or designee, the organizer and the author of a
   started consensus or FINAL draft. Re-evaluate the reviewer count and model diversity after a drop, and let
   shortfalls escalate.
5. **Integrity.** A drop must never become a way to remove a dissenter. Record the verbatim failure evidence
   for every drop.
6. **Default and opt-out.** Use the existing `[defaults].quota_auto_exclude` knob, or a renamed or additional
   one, with a per-idea opt-out, on by default for new ideas. Say how 1.51.0 decks migrate.
7. **Size.** Keep it as small as possible. List the exact protocol hunks for all three COOPERATION.md copies,
   the CLI files and the tests.

## Quorum and roles for THIS idea

- `00-prompt.md`: `author: user`, **no `facilitator:` field** (codex-1 facilitates as an ordinary participant,
  as the owner's model has it), `participants: [codex-1, kimi-1, zcode-1]`, `implementer: codex-1`,
  `track: deliberation` (a protocol change), `auto_implement: true`, `require_model_diversity: true`. Quote
  both owner messages.
- Run `parley preflight` first. claude-1 is not a participant here.
- **Implementer:** codex-1. **Reviewers:** kimi-1 and zcode-1, the non-implementers. Give each a focused
  review brief: the diff, FINAL, the ACs and the validator output.
- **During this run the rule is not in force yet.** If kimi-1 or zcode-1 fails, the shipped rules apply. As
  standing owner permission for this idea, also apply the owner's intent now: a participant that fails twice on
  the same step (the original attempt plus one relaunch) is dropped from this idea, provided codex-1 plus one
  other participant remain. Record it in `00-prompt.md` as an owner-authorized `excluded:` entry and give a
  non-blocking notice. If both drop, write a blocking `codex-1-to-user` note.

## Process

- **Design (Phases 0-4).** Rounds, consensus with all three signoffs, then FINAL. If FINAL implements the
  owner's rule above and leaves **no open owner decision**, continue directly to implementation and release
  without a proposal stop. The owner has asked twice to finish without delays. If an owner decision remains,
  write one blocking `codex-1-to-user_<slug>_decision.md` with the options and a recommendation.
- **Implementation (Phases 5-8).** Write the `IMPLEMENTATION.md` plan before any code, then implement exactly
  FINAL. Reviews come from kimi-1 and zcode-1, separate processes with their own verdicts; you never review
  your own code. Fix-up is capped at 5 cycles. Unresolved CRITICAL or MAJOR findings after cycle 5 escalate.
- **Close.** It is pre-authorized when both final review-consensus signoffs exist, there is no open CRITICAL or
  MAJOR, and the current-tree AC evidence is recorded.
- **Retries.** If a provider (429 or 503) or a timeout affects your own or a participant's step: wait 15
  minutes and relaunch, at most 8 times per step. Timeouts relaunch with 2400 s and then 3600 s. An auth or
  credit error on codex-1 stops the run with a blocking note.
- **Use the driver** (`parley run`, `parley continue`, `parley wait` and `parley status`). Record any gaps; the
  D6 legacy-record gap is known. Never prune or declare worktrees.
- **Never end your turn while waiting.** Wait only with blocking tool calls. End only when you are done or
  deliberately after a blocking owner note.
- English everywhere. Use the commit prefix `[codex-1] meta-protocol-change-participant-dropout: ...`. Never
  copy secrets, and never drive Google Chrome.

## Release (after close)

- Merge to `main`. Target versions are CLI **1.52.0**, skill **2.16.0** and core **2.16.0**.
- The core needs care: core 2.15.0 is staged but not yet published by the owner. Build 2.16.0 on the staged
  2.15.0 plus exactly this idea's reviewed hunks, and say so.
- Channels:
  - GitHub releases;
  - both Homebrew formulae;
  - the skill-only winget PR (CLI winget stays held while Windows is experimental);
  - the skill installed into all managed runtimes **and** the four unmanaged generic locations
    `~/.hermes/profiles/{ldx,librade,testprofile}/skills/parley-deck` and
    `~/.config/opencode/skills/parley-deck`, each via
    `parley-deck-skill install --target generic --dest <dir> --force`, every SKILL.md verified by hash.
- Channel verification by kimi-1 or zcode-1, whichever is available, as a separate process with a short brief.
- Finish with `inbox/codex-1-to-user_meta-protocol-change-participant-dropout_released.md`. It carries the
  channel evidence, the owner-only commands (`! npm publish …` for 2.16.0 and `! parley protocol publish …`;
  if 2.15.0 is still unpublished, list it first), the deferred items and the usage. Then exit.
