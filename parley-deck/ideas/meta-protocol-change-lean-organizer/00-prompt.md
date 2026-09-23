---
idea: meta-protocol-change-lean-organizer
author: user
facilitator: codex-1
created: 2026-09-23
track: deliberation
participants: [claude-1, kimi-1, zcode-1]
auto_implement: true
status: final
---

## Problem / idea

Reduce organizer token usage while preserving canonical participant ownership and verification integrity. The owner approved A–D and continuation through release, subject to the scope stop rule below. codex-1 is the facilitator only and never signs off, implements, or verifies code.

## Constraints

- Use the owner-specified CLI and skill worktrees. Preserve the three-member quorum.
- The owner explicitly authorizes direct-to-main release without development PRs for this run; this overrides the deck's GitHub PR mechanics for this run, without changing its transport header.
- The owner's continuation authorization satisfies ordinary transition approvals within A–D. The organizer must inspect FINAL before Phase 5 and stop if any A–D item is dropped or a new mandatory obligation is added.
- Participants own design, implementation, tests, and code verification. No release, core publication, or main-branch integration by participants before the organizer's release step.
- First-round analysis is independent. Verify evidence locators at HEAD (§15). Keep raw artifacts canonical; no LLM summaries substitute for them.
- Check all three protocol copies; the skill worktree is also in implementation/review scope. Use the installed driver where possible; report gaps.
- Preserve §15 dispute/provenance/alternatives duties, including the drafter's comparison, blind spots, and position changes; the old driver prompt may omit them.

## Evidence

The verified study is copied to `source-context/organizer-token-study.md`. Per-agent underlying evidence remains at the absolute path in the brief below. Runtime/preflight findings and organizer accounting are recorded separately in this idea.

## Non-goals

No environment catalogue pruning, model downgrades, silent model swaps, extra quorum members, mandatory obligations beyond A–D, or unrelated code changes.

## Owner-approved brief (verbatim)

# Organizer brief — meta-protocol-change-lean-organizer

You are **codex-1, the ORGANIZER (facilitator)** of a new Parley Deck idea. You are not a participant and
you do not sign off. Quorum = the owner's global roster: **claude-1, kimi-1, zcode-1** (verify with
`parley roster show --scope machine`; codex-1 is inactive there on purpose). This run was started from a
Claude session; per the owner's global rule, organization is handed to codex-1 and claude-1 stays a participant.

## Workspace

- CLI worktree (branch `lean-organizer`, base tag v1.48.0):
  `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer`. The deck is its `parley-deck/`.
- Skill worktree (branch `lean-organizer`, base tag v2.12.1):
  `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer-skill`.
- Evidence: `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/organizer-token-study/2026-09-23/`
  (`README.md` = verified summary; `results/*.json` = per-agent data and verifier corrections).

## The owner's request (verbatim, Slovak) and translation

> "pozri na poslednu verziu parley-deck, na tie vylepsenia co sa tam robili, budeme robit dalsie vylepsenia,
> hlavne sa zamysli nad tym, ako znizit token usage pre organizatora, kedze to je vacsinou najsilnejsi a
> najdrahsi model, musime jeho tokeny setrit"

"Look at the latest version of parley-deck and the improvements made there; we will make further
improvements. Above all, think about how to reduce the organizer's token usage, since the organizer is
usually the strongest and most expensive model; we must save its tokens."

> "v kazdom pripade ked prides na nejake vylepsenie, ukaz mi ho a potom vydame novu verziu zas cez vsetky kanaly"

"Whenever you come up with an improvement, show it to me, and then we release a new version again through
all channels."

The proposal was shown to the owner on 2026-09-23 and the owner approved **all four scope items** below.

## Approved scope (A-D)

- **A. The organizer does not implement, and does not verify code itself.** Phases 5-8 implementation and
  code verification go to participants; the organizer reads verdicts and validator output. Make this the
  default in the skill and the driver. The protocol already permits it (see README §3).
- **B. `parley wait` + digest.** One blocking call that returns when a round/phase is complete or on a
  timeout below the provider's 30-minute prompt-cache lifetime, printing a deterministic, Go-generated digest
  (who filed, owner, validity, stance/blocks, paths). Digest-first reading; full artifacts are opened only to
  adjudicate a disagreement. Raw artifacts stay canonical; no LLM summarization.
- **C. Organizer brief + slim SKILL.md.** A deterministic role-scoped protocol view for the facilitator (the
  protocol's own facilitator reading set, COOPERATION.md:26-33), used for post-compaction re-orientation
  instead of the full COOPERATION.md + SKILL.md. Split SKILL.md into a slim every-session core and on-demand
  references (e.g. hand-launch templates the runner already owns), and route the organizer through the driver
  (`parley run` / `continue` / `status` / `consensus`).
- **D. Session per phase + organizer measurement.** A state/handoff file so each phase can start in a fresh
  organizer session instead of compacting, and an organizer usage ledger that ingests the client's own
  accounting (Codex rollout `token_count`, Claude `message.usage`) per idea and phase. Address the earlier
  rejection of "self-reported usage records as cost evidence" (client accounting is not a model's
  self-report). Also fix participant telemetry capture (kimi parser gap; plain-text default argv).

The participants may refine the design. **If FINAL drops any of A-D, or adds a new mandatory obligation
beyond them, stop before Phase 5**, write
`parley-deck/inbox/codex-1-to-user_meta-protocol-change-lean-organizer_scope.md` explaining the difference,
and exit. Otherwise continue through release without asking again.

## Run this idea lean (it is its own first test)

1. **Use the driver.** Prefer `parley run` / `parley continue` for rounds, drafting and signoffs, and
   `parley status --idea <slug> --json` / `parley consensus status --json` for state. Do not hand-write
   participant prompts when the runner can launch them. When the driver cannot do a step, fall back manually
   and record why in the idea (that gap is evidence for this idea).
2. **Do not implement or verify code yourself.** The implementer is a participant (FINAL drafter or a
   claimant). Tests and code review are done by participants; you read their results and validator output.
3. **Block, do not poll.** One long blocking call (a foreground `parley run`, or a single shell `until` loop
   with sleeps of 60 s or more, bounded at 25 minutes) instead of repeated short polls.
4. **Re-orient cheaply after a compaction:** this brief, `parley status` for the idea, and only the
   facilitator reading set of COOPERATION.md (Quickstart, §4, §5, §9, your §11 transport). Do not re-read the
   full SKILL.md, the full COOPERATION.md or unrelated skills unless a specific question needs one section.
5. **Keep tool outputs small:** `sed -n` ranges, `rg -n`, `jq` field selection; never print a file over
   20 KB whole.
6. **Measure yourself.** At each phase boundary, append your cumulative `token_count` totals (input, cached,
   output, reasoning, request count) from your own Codex rollout JSONL to
   `parley-deck/ideas/<slug>/organizer-usage.md`. This is the baseline for scope D.

## Idea setup

- Slug `meta-protocol-change-lean-organizer`, `author: user`, `track: deliberation` (protocol change),
  participants claude-1, kimi-1, zcode-1.
- `00-prompt.md` quotes the owner's words above (Slovak + English), states the approved scope A-D and the
  stop rule, and **copies the evidence README into the idea** (e.g. `source-context/organizer-token-study.md`)
  so every participant has it. Participants must verify locators at HEAD before relying on them (§15).

## Implementation and release notes

- COOPERATION.md exists in THREE copies: two in parley-deck-cli (guarded by a Go drift test) and the skill's
  bundled `skills/parley-deck/references/COOPERATION.md` (unguarded). Change all three.
- A protocol change needs a new global core version. `parley protocol publish` is attended-only (it refuses
  without a TTY) and is the owner's action. Build the staged core from the previous core plus the reviewed
  hunks (the core is a TEMPLATE with placeholder header and stub §2; do not copy a deck view), put it in
  `~/.parley/staging/`, and escalate to the owner with the exact command. Do not work around the TTY gate.
- Release after Phase 8 completes: CLI minor bump (1.49.0) and skill minor bump (2.13.0) with CHANGELOG;
  tag; GitHub releases with Windows assets; Homebrew: bump BOTH formulae in homebrew-parley; winget: one
  application per PR (two PRs against microsoft/winget-pkgs; the owner approved catalog PRs); npm:
  `npm publish --access public` of the exact packed tarball (auth restored 2026-09-23; if npm asks for web
  verification, escalate the URL to the owner via the inbox and wait); install the skill into all sessions
  (`parley-deck-skill install --target all --force`) and verify every runtime SKILL.md by content hash.
  Merge to main directly as in the previous release (no development PRs).
- **The deploy is verified independently:** after you perform the release, have a participant verify every
  channel against the real artifacts (both formulae url+sha256; `$(which parley-deck-skill)` resolves into
  Cellar; npm latest + integrity; GitHub assets + sha256; winget PR state; installed SKILL.md content).
  Fix findings before reporting the release complete.
- Deliverable texts, artifacts and commit messages are in English. Commit prefix `[codex-1] <slug>: ...`
  for your own commits. Tokens and credentials are secret: never copy them into files or commits.
- Never open, drive or automate Google Chrome; ego-browser is the only permitted browser.

## When you finish

Write `parley-deck/inbox/codex-1-to-user_meta-protocol-change-lean-organizer_done.md` with: what shipped
(versions, channel evidence), what was deferred and why, the organizer-usage ledger for this run, and any
owner-only actions left (e.g. the attended core publish). Then exit.

