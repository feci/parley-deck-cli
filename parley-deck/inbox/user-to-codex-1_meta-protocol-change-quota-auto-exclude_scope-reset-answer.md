---
from: user
to: codex-1
idea: meta-protocol-change-quota-auto-exclude
phase: implementation
blocking: no
date: 2026-10-04
---

## Owner answers to `codex-1-to-user_meta-protocol-change-quota-auto-exclude_scope-reset.md`

Relayed by the owner's Claude Code session on 2026-10-04 at about 02:35 CEST. The relay asked both questions
in Slovak and gave three options for question 1: A (build a reliable zcode channel, codex-1's
recommendation), B (defer, no supported adapter) and a pragmatic option. The selected options are below,
verbatim. A translation follows each one.

**Question 1, AC2 and zcode support.** Selected: **"Pragmaticky: stderr zcode stačí"** ("Pragmatic: zcode's
stderr is enough"). The option read (Slovak, verbatim): "Moje odporúčanie, ide o najrýchlejšiu funkčnú
cestu. Pri zcode stačí JSON `responseBody` v stderr, ak proces skončil chybou, nenechal výstup a každý
záznam o chybe je 429 „Limit Exhausted“ s `reset_at`. Je to výslovná výnimka z FINAL. Riziko: chyba
sub-agenta by mohla vyradiť agenta, ktorý zlyhal z iného dôvodu. Tlmia to minimum 2 a notifikácia."

Translation: "My recommendation, the fastest path that works. For zcode, the `responseBody` JSON in stderr is
enough when the process ended with an error, left no output, and every error record is a 429 'Limit
Exhausted' with `reset_at`. This is an explicit exception to FINAL. Risk: a sub-agent's error could exclude an
agent that failed for a different reason. The minimum of 2 and the notice soften that."

**Question 2, display time without a timezone.** Selected: **"Áno, podľa návrhu claude-1 (Recommended)"**.
This adopts the bounded interpretation exactly as your note states it. A display clock counts only when it
appears together with a machine reset value (`reset_at` or `retry_after`) in the same record and agrees with
it within one second. Missing, contradictory or display-only resets still gate.

## What this authorizes (an owner-directed deviation from FINAL; record it in IMPLEMENTATION.md)

- The zcode adapter becomes a **supported** recognizer, but only under this owner-defined evidence rule. All
  of the following must hold:
  - the zcode process exited non-zero;
  - it produced no valid artifact, and no later attempt in the batch succeeded;
  - stderr contains at least one provider error record whose `responseBody` JSON is a 429 with explicit
    exhaustion text ("Limit Exhausted", or allowance semantics as in FINAL §4.4) and a machine reset value;
  - **every** provider error record in that stderr agrees: each is that same exhaustion class, and their
    reset values agree within the tolerance. A mixed or contradictory record set gates;
  - the run ends with zcode's turn-failure line;
  - the reset clears FINAL's 60-minute threshold.
- Positive fixtures are this run's recorded zcode stderr (evidence incidents 1 and 2). Adversarial fixtures
  include a quoted 429 in assistant or tool text, a mixed 429 plus other-error stderr, a reset under 60
  minutes, a display-clock-only reset and a success after a 429.
- Everything else in FINAL is unchanged: the floor of 2, the role guards, fail-closed, the record and
  notice, both stages and claude-1's binding full review. Other adapters stay diagnostic-only unless they
  have native evidence.
- Record the deviation in `IMPLEMENTATION.md` with this note quoted, and carry the protocol wording into the
  §9.0 hunk under §7. The claude-1 review checks the deviation as implemented, not whether to make it.
