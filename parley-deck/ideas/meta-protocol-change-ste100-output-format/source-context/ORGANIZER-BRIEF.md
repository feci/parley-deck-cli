# Organizer brief: meta-protocol-change-ste100-output-format

You are **codex-1, the organizer and the implementer** of a new Parley Deck idea. The owner wants Claude
tokens saved: codex-1 organizes and implements, the owner's Claude Code session only relays, and claude-1
takes no part.

## Workspace

- CLI worktree: `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/ste100` (branch `ste100`, from
  origin/main `24fe38f`, CLI 1.54.0). The deck is its `parley-deck/`.
- Skill worktree: `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/ste100-skill` (branch `ste100`,
  from `0549544`, skill 2.18.0).

## The owner's words (verbatim, Slovak) and translation

> "pridaj do parley-deck protokolo language format ASDSTE-100, vsetko co komunikujes, ako formulujes text a
> vytvaras obsah, bude to default output text format for the protocol, sprav ideu na to ako to zapracovat do
> protokolu a potom rovno vydaj novu verziu cez vsetky kanaly"

Translation: "Add the ASD-STE100 language format to the parley-deck protocol. Everything you communicate,
how you formulate text and create content: it will be the default output text format for the protocol. Make
an idea on how to work it into the protocol, and then release the new version directly through all
channels."

## Design target (participants decide the details in FINAL)

ASD-STE100 (Simplified Technical English) becomes the **default writing standard for all English text that
Parley Deck agents produce under the protocol**. That covers round files, consensus, FINAL, IMPLEMENTATION,
reviews, signoffs, inbox notes, PR/MR mirrors, commit messages and generated briefs.

FINAL must settle these points:

1. **Which rules to adopt, and how to cite them.** ASD-STE100 is a controlled language: an approved
   vocabulary with one meaning per word, writing rules and limits on sentence length. Use the current issue
   of the specification and verify which issue that is (§15, with a provenance tag).
   - **The ASD dictionary is copyrighted.** Do not copy the dictionary or long passages of the specification
     into the repository. Summarize the writing rules in your own words, cite the source, and say how agents
     get the full specification.
   - Decide the concrete core rules agents must follow. Typical candidates:
     - procedural sentences of at most 20 words and descriptive sentences of at most 25;
     - paragraphs of at most 6 sentences;
     - one instruction per sentence;
     - imperative mood for procedures;
     - active voice;
     - simple verb tenses;
     - approved words with one meaning each, plus permitted technical names and technical verbs;
     - no `-ing` forms used as verbs.
2. **Where it lives, inside the size limits.** The phase packet is at about 68,846 of 70,000 bytes, and the
   core `SKILL.md` is at about 18,990 of 20,000 bytes. Do not raise the limits.
   - Put the rule summary in a new on-demand reference, for example `references/STE100.md` in the skill
     and a matching protocol appendix or reference, with a short binding pointer in COOPERATION.md and
     SKILL.md.
   - Update the packet applicability map if it is needed.
3. **Exemptions.** Some text must never be rewritten:
   - verbatim quotations of the owner or of other sources, including the Slovak quotes;
   - code, commands, paths, identifiers, log and error excerpts;
   - schema keys;
   - frontmatter;
   - generated or validator text;
   - other agents' text that is quoted.

   Also state how STE relates to the existing English-only rule. Chat with the owner in another language
   (for example Slovak) is outside the protocol's file rule. Say so explicitly. STE applies to English
   output only.
4. **Default and opt-out.** On by default for new ideas and for every new artifact. Decide whether there is a
   per-idea or per-deck opt-out and what it is called. Say how 1.54.0 decks and open ideas adopt it. Existing
   artifacts and the existing protocol text are **not** rewritten. You may recommend a gradual
   re-authoring as deferred work.
5. **Checking.** Decide whether a small mechanical check is worth it, and keep it advisory. Candidates:
   sentence and paragraph length, passive voice heuristics, and a short list of forbidden ambiguous words.
   Reviewers can then cite real violations. Any check must not block on exempt text. Keep it small. A full
   dictionary checker is out of scope, because of the copyright and the size limits.
6. **Brief and prompt templates.** The runner and driver templates that brief participants must tell them
   to write in STE. They must stay within the size limits.
7. **Size.** Keep the change small. List the exact protocol hunks for all three COOPERATION.md copies, plus
   the CLI files, the skill files and the tests.

**Dogfood it.** Write this idea's own new artifacts in STE once the rules are drafted (from round 2 on), and
let the reviewers check them.

## Quorum and roles

- `00-prompt.md`: `author: user`, no `facilitator:` field (codex-1 facilitates as an ordinary participant),
  `participants: [codex-1, kimi-1, zcode-1]`, `implementer: codex-1`, `track: deliberation` (a protocol
  change), `auto_implement: true`, `require_model_diversity: true`. Quote the owner.
- Run `parley preflight`. kimi-1 still fails with `403 … API key's connection allowlist` (relay probe
  2026-10-13 14:10 CEST). That is an owner-side gateway fix. The shipped dropout rules apply: exclude kimi-1
  for this idea and continue with codex-1 and zcode-1. Never touch credentials, keys or gateway settings.
- zcode-1 reviews and goal-checks independently, in separate processes with focused briefs. You never review
  your own code.
- **Close authority.** If the 1.53.0 single-reviewer exception does not apply mechanically, the owner
  pre-authorizes the close when all of these hold:
  - the final zcode-1 review has no open CRITICAL or MAJOR;
  - both current participants filed a review-consensus ACCEPT;
  - a fresh zcode-1 goal check passed;
  - the current-tree AC evidence is recorded.

## Process, retries and release

- **No proposal stop.** The owner said release directly. Go from design through implementation to release
  without stopping. Stop only if an owner decision is truly open; in that case write one blocking
  `codex-1-to-user_<slug>_decision.md` with a recommendation.
- **Fix-up.** At most 5 cycles. If a CRITICAL or MAJOR is still unresolved after cycle 5, escalate.
- **Retries.**
  - A provider 429 or 503, or a timeout: wait 15 minutes and relaunch, at most 8 times per step.
  - A timeout: relaunch with 2400 s, then with 3600 s.
  - An auth or credit error on codex-1: stop with a blocking note.
- **Driver.** D6 is fixed in 1.54.0, but the owner has not yet run the one-time
  `parley budget legacy apply`, so the round-02 driver refusal can still occur. Use the recorded fallback.
  Never apply the legacy declaration yourself. Never prune or declare worktrees.
- **Waiting.** Never end your turn while waiting. Wait only with blocking tool calls. End only when done, or
  deliberately after a blocking owner note.
- **Conventions.** English everywhere. Commit prefix `[codex-1] meta-protocol-change-ste100-output-format: ...`.
  Never copy secrets. Never drive Google Chrome.
- **Release.** Merge to `main`. Then:
  - Versions: CLI **1.55.0**, skill **2.19.0**, core **2.19.0**. Build core 2.19.0 on the staged
    `~/.parley/staging/COOPERATION-2.18.0.md` plus exactly this idea's reviewed hunks.
  - GitHub releases.
  - Both Homebrew formulae.
  - A skill-only winget PR. The CLI winget stays held.
  - Install the skill into all managed runtimes **and** the four generic locations
    (`~/.hermes/profiles/{ldx,librade,testprofile}/skills/parley-deck` and
    `~/.config/opencode/skills/parley-deck`) with
    `parley-deck-skill install --target generic --dest <dir> --force`. Verify every SKILL.md by hash.
  - zcode-1 verifies all channels in a separate process.
- **Finish.** Write `inbox/codex-1-to-user_meta-protocol-change-ste100-output-format_released.md` with:
  - the channel evidence;
  - the still-pending owner commands, in version order: the D6 legacy apply from the driver-unstall note,
    then npm and core for each still-unpublished version from 2.15.0 to 2.19.0;
  - the deferred items;
  - usage.

  Then exit.
