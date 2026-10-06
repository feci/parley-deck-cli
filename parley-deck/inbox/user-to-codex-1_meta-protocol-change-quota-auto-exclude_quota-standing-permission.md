---
from: user
to: codex-1
idea: meta-protocol-change-quota-auto-exclude
phase: review-consensus
blocking: no
date: 2026-10-06
---

## Owner answer to `codex-1-to-user_meta-protocol-change-quota-auto-exclude_fix-consensus-quota-stop-20261005.md`

Relayed by the owner's Claude Code session. Received 2026-10-06 at about 12:18 CEST. The question and the
answer, verbatim (Slovak), with a translation:

> Question: "Zopakované review prešlo: claude-1 našiel 0 CRITICAL, 3 MAJOR, 2 MINOR a 1 NIT a codex na ne
> pripravil 5 opráv. Podpis plánu opráv od claude-1 však o 21:57 znova spadol na „429 All claude accounts
> have exhausted their quota (reset after 5m)“. Claude pool na gatewayi sa opakovane vyčerpá na pár minút a
> každá takáto zastávka teraz čaká hodiny na teba. Ako ďalej?"
> Selected: **"Trvalé povolenie: počkať a opakovať (Recommended)"**. The option read: "Pri chybe kvóty claude
> s resetom do 60 minút codex počká do resetu plus 2 minúty a krok zopakuje, najviac 3-krát na jeden krok.
> Až potom sa zastaví a ozve sa. Platí len pre túto ideu a nikoho to nevyradí."

Translation: "Standing permission: wait and retry. On a claude quota error with a reset within 60 minutes,
codex waits until the reset plus 2 minutes and repeats the step, at most 3 times per step. Only then does it
stop and report. This applies to this idea only and excludes nobody."

## Standing rule for the rest of this idea (it replaces the brief's "do not spin retries" for this case only)

- **When it applies.** The claude-1 participant (review, signoff or channel verification) fails with a
  quota or rate-limit error that carries a stated reset of **60 minutes or less**. Example: `429 … All claude
  accounts have exhausted their quota … (reset after 5m)`.
- **What to do.** Wait with blocking tool calls until the stated reset plus 2 minutes, then relaunch **the
  same step** with the same inputs.
- **Limit.** At most **3 relaunches per step**. After the third failure, or on any error without a reset of
  60 minutes or less, or on an auth or credit error, stop and write a blocking note as before.
- **Records.** Record every attempt in `organizer-notes.md` with the verbatim error.
- **Never** change the model or provider, substitute a reviewer or exclude an agent. The same rule applies if
  codex-1's own model hits a short quota window.

## Now

The claude-1 Phase-7 signoff on the cycle-2 fix plan in `review/consensus.md` is authorized under this rule.
Then continue with fix-up cycle 2, re-review, the attended-close request and the release, all as
`IMPL-ORGANIZER-BRIEF.md` says.
