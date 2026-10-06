---
from: user
to: codex-1
idea: meta-protocol-change-quota-auto-exclude
phase: review-consensus
blocking: no
date: 2026-10-06
---

## Owner answer to `codex-1-to-user_meta-protocol-change-quota-auto-exclude_fix-consensus-timeout-20261006.md`

Relayed by the owner's Claude Code session, 2026-10-06, about 13:26 CEST. Verbatim (Slovak), with a
translation:

> Question: "Podpis claude-1 tentoraz nespadol na kvótu, ale do 20-minútového limitu bez jedinej chyby
> nevrátil nič. Codex sa preto znova zastavil a pýta sa. Rozšíriť trvalé povolenie aj na takéto timeouty,
> aby sa na teba pri každom nemuselo čakať?"
> Selected: **"Áno, aj timeouty (Recommended)"**. The option read: "Pri timeoute codex krok zopakuje s dlhším
> limitom (40, potom 60 minút), najviac 2-krát na jeden krok. Rovnaký postup pri timeoute už predpisuje
> skill. Kvóty platia ako doteraz. Zastaví sa iba pri chybe autentifikácie alebo kreditov, alebo keď vyčerpá
> pokusy."

Translation: "Yes, timeouts too. On a timeout, codex repeats the step with a longer limit (40, then 60
minutes), at most twice per step. The skill already prescribes this for timeouts. Quota handling stays as
before. It stops only on an auth or credit error, or when it runs out of attempts."

## Standing rule, extending `…_quota-standing-permission.md` for the rest of this idea

- **Silent timeouts.** If any participant step (claude-1 or codex-1) times out with no provider error,
  relaunch the same step with the same inputs and a longer ceiling: first **2400 s**, then **3600 s**. That
  is at most **2 timeout relaunches per step**. Record each attempt in `organizer-notes.md`.
- **Short quota windows.** The rule in `…_quota-standing-permission.md` still applies unchanged: a reset of
  60 minutes or less, wait for the reset plus 2 minutes, at most 3 relaunches per step.
- **When to stop and write a blocking note.**
  - an auth or credit error;
  - a quota error with no reset, or a reset longer than 60 minutes;
  - the attempts for that step are used up;
  - any decision that changes scope, FINAL, the quorum or the owner's prior rulings.
- **Never** change the model or provider, substitute a reviewer or exclude an agent.

## Now

Relaunch claude-1's Phase-7 signoff on the unchanged cycle-2 plan with a 2400 s ceiling. Then continue with
fix-up cycle 2, re-review, the attended-close request and the release, all as `IMPL-ORGANIZER-BRIEF.md`
says. Routine procedural steps that these rules or the skill already cover need no owner question.
