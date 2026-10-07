---
from: user
to: codex-1
idea: meta-protocol-change-quota-auto-exclude
phase: review-consensus
blocking: no
date: 2026-10-07
---

## Owner direction: finish now (supersedes the wait in `…_long-quota-answer.md` and `…_provider-stop-answer.md`)

Relayed by the owner's Claude Code session on 2026-10-07 at about 23:20 CEST. Verbatim (Slovak):

> "sakra tak to fixni a dokonci a deployni cez vsetky kanaly, taha sa to dlho"

Translation: "Damn, then fix it, finish it and deploy it through all channels, this is dragging on."

## Relay facts (PRIMARY, 23:18 CEST)

- `ANTHROPIC_BASE_URL` points every Claude CLI call at the OmniRoute gateway, so there is no direct route.
- Two probes ran, each with a 103,660-byte prompt (FINAL + consensus + review consensus): one with
  `--model 'claude-opus-5-5[1m]'` (plain id) and one with `--model 'claude/claude-opus-5-5[1m]'`
  (prefixed id). **Both returned `PONG`.** Large requests pass right now. The 22:01 failure came after
  857 s of an agentic session, so the gateway pool is intermittent, not hard down.
- The relay's own long sessions use the plain id and kept working through the pool errors that hit the
  prefixed id.

## What the owner authorizes now

1. **Start immediately.** Do not wait for 2026-10-09. The relay killed the auto-resume daemon.
2. **claude-1 stays the reviewer and the model stays Opus 5.5.** For the rest of this idea, launch claude-1
   with the plain model id `claude-opus-5-5[1m]` instead of `claude/claude-opus-5-5[1m]`. This is the
   same model at the same max effort; only the gateway route id changes. Record it in
   `organizer-notes.md`. Do not edit `agents.toml`.
3. **Provider errors no longer stop you.**
   - If a claude-1 or codex-1 step fails with a 429 or 503 quota/unavailable error, whatever reset it
     states, wait 15 minutes and relaunch the same step. Do this at most 8 times per step.
   - Silent timeouts keep the earlier rule: relaunch with 2400 s, then 3600 s.
   - Stop only on an auth or credit error, or when the attempts for a step run out.
4. **Fix-up.** Do cycle 4 now: the plan signoff, the implementation and the full re-review. If the cycle-4
   re-review finds new findings, fix them narrowly in cycle 5 without asking. That is the last cycle the
   protocol's deliberation cap allows. Escalate only if a CRITICAL needs a change of scope or FINAL, or
   if findings remain after cycle 5.
5. **Close is pre-confirmed.** The owner's "dokonci a deployni" is the attended-close confirmation,
   provided all of the following hold:
   - the final re-review has no open CRITICAL or MAJOR;
   - both review-consensus signoffs exist;
   - current-tree evidence for AC1 to AC21 is recorded, with AC2 owner-waived as already decided.

   Record this verbatim as the close authority. Do not write a separate close-request note.
6. **Release immediately after the close,** on all channels, as `IMPL-ORGANIZER-BRIEF.md` says:
   - merge to `main`;
   - CLI 1.51.0, skill 2.15.0 and core 2.15.0 staged;
   - GitHub releases;
   - both Homebrew formulae;
   - a winget PR for the skill only;
   - install the skill into all local runtimes and verify each by hash;
   - claude-1 verifies every channel independently, using a short brief.

   The two steps only the owner can run (`! npm publish …` and `! parley protocol publish …`) go in ONE
   final note, `codex-1-to-user_meta-protocol-change-quota-auto-exclude_released.md`, with the exact
   commands.
7. **Unchanged:** the owner's earlier answers (the zcode stderr rule, the AC2 waiver and the follow-up),
   no other reviewer, no quorum change, and English artifacts.
