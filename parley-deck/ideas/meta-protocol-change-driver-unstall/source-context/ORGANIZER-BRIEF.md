# Organizer brief: driver-unstall (D6 legacy accounting + goal-check ceiling)

You are **codex-1, the organizer and the implementer** of a new Parley Deck idea, slug `driver-unstall`. You
can prefix the slug with `meta-protocol-change-` if FINAL changes protocol text. The owner wants Claude tokens
saved, so codex-1 organizes and implements, the owner's Claude Code session only relays, and claude-1 takes
no part.

## Workspace

- CLI worktree: `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/driver-unstall`, branch
  `driver-unstall`, from origin/main `128e30b` (CLI 1.53.0). The deck is its `parley-deck/`.
- Skill worktree: `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/driver-unstall-skill`, branch
  `driver-unstall`, from `8ce4dec` (skill 2.17.0).

## Owner authority (verbatim, Slovak) and translation

- **2026-10-04.** The owner answered the D6 question "D6: driver sa pri každej novej idei zastaví na starom
  zázname z mája. Čo s tým?" ("D6: the driver stops at every new idea on an old record from May. What to do
  about it?") by selecting **"Samostatná trvalá oprava (Recommended)"**, a separate lasting fix in its own
  parley idea.
- **2026-10-08.** "participanti nie su nevyhnutne potrebny obaja, staci jeden a to by nemalo zaseknut
  parley-deck" ("…one is enough and that should not stall parley-deck").
- **2026-10-09.** "ano spusti ideu ako poriesit ze participant vypadne v strede idei, ved tam mame nejake
  timeouts…" ("yes, start the idea on how to handle a participant dropping out mid-idea, we have some
  timeouts…").

The owner's goal across all three is that **the driver must not stall**.

## The two stalls to fix (verify at HEAD, §15)

1. **D6, legacy accounting.** The driver halts before every new idea's first driver-launched cross-review
   round with `run round-02: cross-review accounting: historical cycle event lacks idea identity`. The cause
   is line 1 of `parley-deck/runs/20260510T194003Z/events.jsonl`, the May smoke run's `run.created` event
   with no `idea` field, committed in `3ec10ac`. The only remedy today is an attended, per-idea, unstored
   `parley budget migrate apply --declare-unscoped-run …`. It has hit **every** idea since 2026-10-03, and
   every organizer has had to fall back to hand-launched rounds. See the organizer notes of
   `meta-protocol-change-quota-auto-exclude` (driver gap 11), `meta-protocol-change-participant-dropout`
   and `meta-protocol-change-participant-dropout-review-gate-timing`.
   - Design a lasting fix that keeps the budget system's fail-closed integrity: the gate must never silently
     ignore real unscoped history.
   - Candidate directions for the participants to judge: a durable, recorded one-time declaration, a gate
     scoped to the idea's own history, or a provable migration of the legacy record.
2. **Goal-check ceiling.** The product goal-check ceiling is 120 s. Real independent goal checks took 518 s
   (review-gate-timing) and similar in earlier runs, so a driver auto-goal check times out and the close
   escalates.
   - Make the ceiling fit real runs. For example, derive it from the track's per-agent timeout and route it
     through the shipped one-retry and watchdog path, without weakening LE-7/LE-11.

Also consider the other recorded driver gaps from those three runs' organizer notes. Include only the ones
that are small and clearly stall the driver, such as the phase-pointer or ready-consensus reopen gaps. List
the rest as deferred. **Keep it small.**

## Quorum and roles

- Write `00-prompt.md` as follows:
  - `author: user`;
  - no `facilitator:` field, because codex-1 facilitates as an ordinary participant;
  - `participants: [codex-1, kimi-1, zcode-1]`;
  - `implementer: codex-1`;
  - the track per §4.0. If FINAL changes protocol text, `deliberation` is forced;
  - `auto_implement: true` and `require_model_diversity: true`.

  Quote the owner's words.
- Run `parley preflight`. kimi-1 has been failing with `403 … API key's connection allowlist`, which is an
  owner-side gateway fix. If it still fails, the shipped 1.52.0/1.53.0 dropout rules apply: exclude it for this
  idea and continue with codex-1 and zcode-1. Never touch credentials, keys or gateway settings.
- zcode-1 is the independent reviewer and goal-checker, in separate processes with focused briefs. You never
  review your own code.
- **Dogfood the 1.53.0 rule.** A recorded dropout permits one model-diverse reviewer. If the product exception
  does not apply mechanically, the close is **pre-authorized by the owner** under these conditions:
  - the final zcode-1 review has no open CRITICAL or MAJOR finding;
  - review-consensus ACCEPT blocks exist from both current participants;
  - a fresh zcode-1 goal check passes;
  - the current-tree AC evidence is recorded.

## Process, retries and release

- Go from design through implementation to release with no proposal stop, unless an owner decision stays
  open. In that case, write one blocking `codex-1-to-user_<slug>_decision.md` with the options and a
  recommendation.
- Fix-up cycles are capped at 5. If a CRITICAL or MAJOR is unresolved after cycle 5, escalate.
- Retries:
  - On a provider 429 or 503, or a timeout: wait 15 minutes and relaunch, at most 8 times per step.
  - Timeouts relaunch with 2400 s, then 3600 s.
  - An auth or credit error on codex-1 stops the run with a blocking note.
- Never end your turn while waiting. Wait only with blocking tool calls. End only when done, or deliberately
  after a blocking owner note.
- English everywhere, commit prefix `[codex-1] <slug>: ...`. Never copy secrets, never drive Google Chrome, and
  never prune or declare worktrees.
- **Release** (a CLI release always ships the skill too):
  - Merge to `main` and pick versions by impact. A protocol change means CLI 1.54.0, skill 2.18.0 and core
    2.18.0, staged on the staged 2.17.0. A CLI-only change means CLI 1.53.1 and skill 2.17.1, with no core.
  - Publish GitHub releases.
  - Bump both Homebrew formulae.
  - Open a skill-only winget PR. The CLI winget stays held.
  - Install the skill into all managed runtimes **and** the four generic locations
    `~/.hermes/profiles/{ldx,librade,testprofile}/skills/parley-deck` and
    `~/.config/opencode/skills/parley-deck`, using
    `parley-deck-skill install --target generic --dest <dir> --force`. Verify every SKILL.md by hash.
  - zcode-1 verifies the channels in a separate process.
- Finish with `inbox/codex-1-to-user_<slug>_released.md`. It carries:
  - the channel evidence;
  - every still-unpublished owner command, in version order (npm and core for 2.15.0, 2.16.0, 2.17.0 and
    any new one), each listed only if it is still unpublished at that time;
  - the deferred items;
  - usage.

  Then exit.
