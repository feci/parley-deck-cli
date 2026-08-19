---
idea: protocol-and-skill-audit
status: consensus-draft
drafted-by: claude-1
date: 2026-08-20
track: standard
participants: [claude-1, codex-1, hermes-1, kimi-1, opencode-1, zcode-1]
rounds: 2
---

# Audit consensus — 33 findings confirmed, 14 contested, and the verifiers disagreed about each other

## 1. The verification split is itself a finding

Three agents adversarially verified 47 round-1 findings. They were told to default to REFUTED.

| verifier | assessed | CONFIRMED | PARTIAL | REFUTED | UNREPRO |
| --- | --- | --- | --- | --- | --- |
| @codex-1 | 23 (all non-codex) | 6 | 8 | **9** | 0 |
| @kimi-1 | 42 | **37** | 8 | 2 | 1 |
| @zcode-1 | 32 | 30 | 2 | 0 | 0 |

**The same corpus drew 6 confirmations from one verifier and 37 from another.** Every REFUTED
verdict but one came from @codex-1, and every finding @codex-1 refuted belongs to @zcode-1 or
@claude-1.

There is a structural reason and it must not be read as bias: **@codex-1 wrote 24 of the 47
findings, so it was the only verifier not permitted to assess them — and the two agents who did
assess them are the two lenient ones.** @codex-1's 24 findings were checked only by verifiers who
confirmed ~92% of everything they touched. That is an asymmetry in the verification design, mine,
not a property of the findings.

## 2. Confirmed — 33 findings, no refutation from any verifier

**@codex-1 (21):** F1–F11, F13, F14, F16–F22, F24 — the consensus/driver enforcement family.
Twenty variants of one defect: **the gates accept emptiness.** A blank round passes as complete
(F1, F17); an empty `FINAL.md` closes an idea (F5); three padded lines pass as a full specification
(F22); signoff-shaped headings outside `## Signoffs` satisfy the gate (F20); an empty
`responding-to:` passes cross-review (F16); a review with no `reviewed-commit` validates (F18).
Plus track enforcement: an explicit track the classifier rejects is accepted (F13) and a missing
`track:` does not apply the documented default (F14).

**@zcode-1 (7):** F1, F6, F7, F8, F11, F14, F15 — §2's roster-authority text is false for this deck;
§12.12 cites a slug that exists nowhere; Quickstart omits §15 and §10; §3's layout omits
`agents.toml` and `runs/`; §11.B contradicts its own branch-protection advice; `learn` and
`preset list` are invisible in `--help`.

**@kimi-1 (3):** F1 `doctor` does not byte-verify the managed core skill; F2 `sync-project --yes`
silently deletes `protocolRole`, the field §9.0 gates on, while `status` recommends that command;
F3 the README says "fourteen named runtimes" in four places and omits `zcode` from `--target`.

**@claude-1 (2):** F2 `COOPERATION.md:57` tells the bootstrap to run `roster render`, which writes a
§2 shape the repo's own drift guard fails closed on; F3 `masked-by-env` is documented in the closed
STATUS vocabulary and never reaches STATUS.

**@hermes-1 (1):** Q2 — 6 of 71 `IMPLEMENTATION.md` never reach `status: complete`; two sit at
`ready-for-review` awaiting a review nobody will run. [Verified by @claude-1, whose own
counter-measurement was the error — `grep -rLl` reported the exact inverse.]

## 3. Contested — 14 findings, NOT resolved by count (§15.3)

**@codex-1 refuted seven of @zcode-1's** (F2, F3, F5, F9, F10, F12, F13) while @kimi-1 or @zcode-1's
own evidence confirmed them. These need adjudication on evidence, and the drafter will not break
them by majority.

**@claude-1/F1** (SKILL.md never names `preflight`/`retro`/`loop tick`/…) — PARTIAL from one
verifier, REFUTED by @codex-1. **This is my finding and I will not adjudicate it.**

**@kimi-1/F5** — REFUTED outright, the only refutation not authored by @codex-1.

**PARTIAL-only, no confirmation:** @codex-1/F12, F15, F23; @kimi-1/F4; @zcode-1/F4.

**Already corrected before this draft:** @hermes-1/Q3 (claimed 45% of closed ideas stale; measured
1 of 78 — REFUTED at CRITICAL, survives as NIT) and, by the same measurement, **@zcode-1/F14's 27%
figure is in doubt** even though F14 is otherwise confirmed. Both counted a population that
includes genuinely open ideas. **F14 enters the fix list with its number struck, not its claim.**

## 4. Participation

@hermes-1 filed at the **fourth** attempt, after three runs produced complete analyses and lost
them all. What worked was not a better prompt but a smaller slice plus write-first-append-as-you-go;
its file accounts for its own budget ("4 tool calls used") and names the command it could not
afford. @opencode-1 created its skeleton file — progress on three prior total losses — but filed no
finding. Neither absence is consent and neither is a demonstration of failure.

## 5. What happens next

The 33 confirmed findings are the fix list. They split into three surfaces that cannot be reviewed
as one change: **CLI enforcement** (@codex-1's 21 + @claude-1/F3), **protocol text** (@zcode-1's 7 +
@claude-1/F2), and **the skill package** (@kimi-1's 3). The owner has ruled that no complete release
ships until the command defects are fixed.

## Signoffs

_Append below, or write `signoff-<agent>.md` for verbatim concatenation. Sequential only._
