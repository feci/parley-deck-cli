# Organizer brief — meta-protocol-change-designated-implementer

You are **codex-1, the ORGANIZER (facilitator)** of a new Parley Deck idea. Declare it: put
`facilitator: codex-1` in `00-prompt.md` (the 1.49.0 declared-facilitator mode). You are not a participant,
you do not sign off, and you never implement or verify code. Quorum = the owner's global roster:
**claude-1, kimi-1, zcode-1** (verify with `parley roster show --scope machine`). This run was started from
a Claude session; per the owner's global rule, organization is handed to codex-1 and claude-1 stays a participant.

## Workspace

- CLI worktree (branch `designated-implementer`, base origin/main 9134c7a = v1.49.0 plus 9 unreleased
  lean-organizer Linux-CI remediation commits):
  `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/designated-implementer`. The deck is its `parley-deck/`.
- Skill worktree (branch `designated-implementer`, base origin/main 8161e5e = skill 2.13.0 metadata):
  `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/designated-implementer-skill`.

## The owner's words (verbatim, Slovak) and translation

> "cize organizator ma specialne miesto v parley-deck protokole, rozmyslam ze este by nebolo zle setnut aj
> implementatora, ktory vlastne napise kod, ak treba, alebo teda spravi tu pracu co treba a ostatny to
> ohodnotia, kto to robi teraz? je to v protokole predpisane?"

"So the organizer has a special place in the parley-deck protocol. I am thinking it would be good to also set
the implementer, who actually writes the code if needed, or does whatever work is needed, while the others
evaluate it. Who does it now? Is it prescribed in the protocol?"

> "ano otvor na to ideu a navrhnite to cez /parley-deck cim skor"

"Yes, open an idea for it and design it through /parley-deck as soon as possible."

> "a ked na mna nebudete mat otazky tak to rovno releasni a deployni cez vsetky kanaly"

"And if you have no questions for me, release it and deploy it through all channels directly."

## Facts to verify at HEAD before relying on them (§15; these are locators, not conclusions)

In `parley-deck/COOPERATION.md` at origin/main:
- Phase 4 (~L407-409): the FINAL drafter is the idea's `author:`; with `author: user` the default drafter is
  the first agent to submit a round-01 file; any participant may volunteer (`Drafter: yes`).
- Phase 5 (~L443): the default implementer is the FINAL drafter; another participant may claim via
  `inbox/<from>-to-all_<slug>_impl-claim.md` before work begins. A declared facilitator never implements or
  verifies code unless `facilitator_participates: true` (also ~L883).
- `roles:` (~L311) are advisory lenses only: they do not change quorum, signoff weight, artifact ownership,
  drafter eligibility or roster membership.
- Reviewers are all non-implementers (~L234); no merge without an invokable non-implementer reviewer (~L538);
  the goal-done check is by a fresh non-implementer (~L689).
- Code: two separate unexported `resolveImplementer` functions exist, `internal/app/driver_impl.go:104` and
  `internal/consensus/consensus.go:751`. Check whether they agree.
- Precedent: `facilitator:` / `facilitator_participates:` shipped in 1.49.0 (idea
  `meta-protocol-change-lean-organizer` in the lean-organizer worktree deck).
- Observed in that run: kimi-1 claimed drafting, zcode-1 claimed implementation, and claude-1's first code
  review found 1 CRITICAL and 8 MAJOR findings in the implementation, all fixed in later cycles.
- Hypothesis to test, not a conclusion: with `author: user`, the default drafter and therefore the default
  implementer is decided by who files round-01 first, i.e. by speed rather than by suitability.
- Cost context: implementation is the heaviest token work of a run; the organizer-token study measured it at
  ~30-40% of an organizer's input when the organizer implemented
  (`/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/organizer-token-study/2026-09-23/README.md`).

## The design question

Should the protocol let the owner designate the implementer, the agent that executes FINAL (code or other
work) while the others evaluate it, and if so how? Participants decide. Address at least: scope of the role
(code and non-code work); per-idea field versus a global default and their precedence; what happens when the
designated agent is unavailable; drafter versus implementer (FINAL must stay self-contained for a
non-drafter); reviewer model diversity when one agent implements everything; whether the claim mechanism
stays as the fallback; driver and preflight support and fail-closed validation; whether the implementer takes
part in design rounds and signoffs; consistency of the two `resolveImplementer` copies.

**Owner-decision boundary.** Ship the mechanism with the global default UNSET, so behaviour without it is
exactly today's. Do not choose which agent is the owner's default implementer; after release, report a
recommendation with evidence and let the owner set it. That is the only owner question anticipated.

Non-goals: no quorum or roster change, no model or effort downgrade, no change to the 1.49.0 organizer rules
beyond what this role needs, no unrelated code changes.

## How to run it (lean; use the 1.49.0 tooling)

- Classify the track per §4.0. The owner asked for speed: use the lightest track the protocol permits for a
  meta-protocol change, and say why in `00-prompt.md`.
- Use the driver and the new commands: `parley run`/`continue`, `parley wait` (blocking, under 25 minutes per
  call), `parley organizer brief` for re-orientation after compaction, `parley status --json`. Do not
  hand-write participant prompts when the runner can launch them; record every driver gap you hit.
- Do not implement or verify code yourself. Keep tool outputs small. Record your own usage at each phase
  boundary in `parley-deck/ideas/<slug>/organizer-usage.md` (use the new usage ingest if it works).
- `00-prompt.md`: `author: user`, `facilitator: codex-1`, participants claude-1, kimi-1, zcode-1; quote the
  owner's words above; copy this brief into `source-context/` so every participant has it.

## Release (owner-authorized when there are no owner questions)

- **Wait for the other release first.** Another codex-1 organizer (process 21386, idea
  `meta-protocol-change-lean-organizer`) is still finishing release 1.49.x on main: Linux CI remediation, the
  winget CLI 1.49.0 PR, npm 2.13.0 (owner login pending) and core 2.13.0 (owner publish pending). Do not start
  your release until that run has written `codex-1-to-user_meta-protocol-change-lean-organizer_done.md` in its
  deck inbox or its process has exited. Then integrate onto the latest origin/main and take the next minor
  versions above whatever is released.
- Channels: GitHub releases with Windows assets; Homebrew (bump BOTH formulae); winget one application per PR;
  npm `npm publish --access public` of the exact packed tarball (if auth fails, write an owner inbox note
  asking the owner to run `! npm login` in the Claude session, and continue the other channels); install the
  skill into all sessions and verify every runtime SKILL.md by content hash. Merge to main directly, no
  development PRs. A participant independently verifies every channel before you report completion.
- Protocol core: three COOPERATION.md copies must change (two in the CLI, guarded; the skill's bundled
  `skills/parley-deck/references/COOPERATION.md`). Stage the next core from the core TEMPLATE on top of the
  latest core: the published 2.13.0 if the owner has published it, otherwise the staged
  `~/.parley/staging/COOPERATION-2.13.0.md`. It must contain exactly this idea's reviewed hunks on top. Verify,
  then escalate the exact attended `parley protocol publish` command; never work around the TTY gate.
- English for every artifact and commit; commit prefix `[codex-1] <slug>: ...`. Never copy secrets. Never
  open or drive Google Chrome.

## When you finish

Write `parley-deck/inbox/codex-1-to-user_meta-protocol-change-designated-implementer_done.md`: what shipped,
channel evidence, deferred items, the organizer-usage ledger, the recommended default implementer with its
evidence, and any owner-only actions left. Then exit.

