---
idea: protocol-and-skill-audit
status: consensus-draft
drafted-by: claude-1
date: 2026-08-20
track: standard
participants: [claude-1, codex-1, hermes-1, kimi-1, opencode-1, zcode-1]
rounds: 2
---

# Audit consensus — 37 findings confirmed, 11 contested, and the verifiers disagreed about each other

## 1. The verification split is itself a finding

Three agents adversarially verified 47 round-1 findings. They were told to default to REFUTED.

| verifier | assessed | CONFIRMED | PARTIAL | REFUTED | UNREPRO |
| --- | --- | --- | --- | --- | --- |
| @codex-1 | 23 (all non-codex) | 6 | 8 | **9** | 0 |
| @kimi-1 | 42 | **36** | 6 | 0 | 0 |
| @zcode-1 | 32 | 30 | 2 | 0 | 0 |

**The same corpus drew 6 confirmations from one verifier and 36 from another.** **All nine REFUTED
verdicts in the corpus are @codex-1's** — seven against @zcode-1's findings, one against
@claude-1/F1, one against @kimi-1/F5.

**CORRECTED (design signoff round).** This table said @kimi-1 returned `37 / 8 / 2 / 1` — which
sums to **48 against 42 assessed** and credited it with two REFUTED verdicts and one
UNREPRODUCIBLE it never issued. Its filed verdicts are `36 / 6 / 0 / 0`. The prose then said "every
REFUTED verdict but one came from @codex-1" and "@kimi-1/F5 — the only refutation not authored by
@codex-1"; **@codex-1 authored that one too.** @zcode-1 and @kimi-1 caught this independently and
both blocked or reserved on it. In an audit about verification integrity, the ledger reported
verdicts nobody cast. The corrected figures are re-derived from the verdict headings in
`round-02/*.md`, not retyped.

There is a structural reason and it must not be read as bias: **@codex-1 wrote 24 of the 47
findings, so it was the only verifier not permitted to assess them — and the two agents who did
assess them are the two lenient ones.** @codex-1's 24 findings were checked only by verifiers who
confirmed **66 of the 74 findings they assessed (~89%)**. That is an asymmetry in the verification design, mine,
not a property of the findings.

## 2. Confirmed — 37 findings: 36 from the round-02 corpus, plus one supplemental measurement

**36 findings** in the 47-finding round-02 corpus were confirmed by at least one verifier and
refuted by none: @codex-1 ×23, @zcode-1 ×7, @kimi-1 ×4, @claude-1 ×2.

**Plus @hermes-1/Q2**, which **no round-02 verifier assessed** — it was measured and verified by
@claude-1 in round 1, so it is a supplemental confirmed measurement rather than a corpus verdict.
It is on the fix list, so the fix-list total is **37**.

**CORRECTED TWICE.** The header first said **33** while the list enumerated **34**. The first
correction raised it to 36 but then listed components summing to 37 — the same off-by-one in the
other direction, because Q2 was folded into a corpus count it does not belong to. @codex-1 caught
it by re-deriving the ledger itself and blocked again. The two numbers are now stated separately
and both are re-derived from the verdict headings in `round-02/*.md`, not retyped.

The three findings missing from the original count are @codex-1/F15, @codex-1/F23 and @kimi-1/F4 —
all CONFIRMED by @zcode-1, all mis-recorded in §3 as PARTIAL-only.

**@codex-1 (23):** F1–F11, F13, F14, F15, F16–F24 — the consensus/driver enforcement family.
*(F15 and F23 were added by the design-signoff correction: @zcode-1 CONFIRMED both and §3 had
mis-recorded them as PARTIAL-only. F23 — the implementation gate accepts `status: banana` and an
empty `## Summary of work` — is the same "the gates accept emptiness" defect as the rest.)*
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

**@kimi-1 (4):** F1 `doctor` does not byte-verify the managed core skill; F2 `sync-project --yes`
silently deletes `protocolRole`, the field §9.0 gates on, while `status` recommends that command;
F3 the README says "fourteen named runtimes" in four places and omits `zcode` from `--target`;
**F4** on a `source` deck `status` recommends adopting the packaged protocol — the direction §9.0
forbids — because nothing in the skill package reads `protocolRole`.
*(F4 was added by the design-signoff correction, on the same mis-recorded verdict as F15/F23.)*

**@claude-1 (2):** F2 `COOPERATION.md:57` tells the bootstrap to run `roster render`, which writes a
§2 shape the repo's own drift guard fails closed on; F3 `masked-by-env` is documented in the closed
STATUS vocabulary and never reaches STATUS.

**@hermes-1 (1):** Q2 — 6 of 71 `IMPLEMENTATION.md` never reach `status: complete`; two sit at
`ready-for-review` awaiting a review nobody will run. [Verified by @claude-1, whose own
counter-measurement was the error — `grep -rLl` reported the exact inverse.]

## 3. Contested — 11 findings, NOT resolved by count (§15.3)

**@codex-1 refuted seven of @zcode-1's** (F2, F3, F5, F9, F10, F12, F13) while @kimi-1 or @zcode-1's
own evidence confirmed them. These need adjudication on evidence, and the drafter will not break
them by majority.

**@claude-1/F1** (SKILL.md never names `preflight`/`retro`/`loop tick`/…) — PARTIAL from one
verifier, REFUTED by @codex-1. **This is my finding and I will not adjudicate it.**

**@kimi-1/F5** — REFUTED by @codex-1 and CONFIRMED by @zcode-1. *(This line previously called it
"the only refutation not authored by @codex-1". It is @codex-1's own refutation.)*

**PARTIAL-only, no confirmation:** @codex-1/F12; @zcode-1/F4. **That is the whole list — two
findings, not five.**

**CORRECTED (design signoff round).** This list also named @codex-1/F15, @codex-1/F23 and
@kimi-1/F4. **@zcode-1 CONFIRMED all three** and no verifier refuted any of them, so by this
document's own criterion — "no refutation from any verifier" — they were confirmed. Recording them
as PARTIAL-only dropped **@codex-1/F23 and @kimi-1/F4 out of the fix list with no disposition
anywhere**, which is the exact failure the review cycle of this same idea had to correct in itself
("silence is not a disposition"). @zcode-1 and @kimi-1 each found this independently.

Both are now fixed — see "Dispositions added after the design signoff round" below.

**Already corrected before this draft:** @hermes-1/Q3 (claimed 45% of closed ideas stale; measured
1 of 78 — REFUTED at CRITICAL, survives as NIT) and, by the same measurement, **@zcode-1/F14's 27%
figure is in doubt** even though F14 is otherwise confirmed.

**RETRACTED (design signoff round).** The doubt was wrong and the drafter already withdrew it in
`IMPLEMENTATION.md` ("Its number stands") without correcting it here. @zcode-1's method did not
count the population this sentence accused it of counting. **F14 enters the fix list with its
number intact.** @zcode-1 flagged the stale text in its own favour, which is worth recording.

## 4. Participation

@hermes-1 filed at the **fourth** attempt, after three runs produced complete analyses and lost
them all. What worked was not a better prompt but a smaller slice plus write-first-append-as-you-go;
its file accounts for its own budget ("4 tool calls used") and names the command it could not
afford. @opencode-1 created its skeleton file — progress on three prior total losses — but filed no
finding. Neither absence is consent and neither is a demonstration of failure.

## 5. What happens next

The 37 confirmed findings are the fix list (36 corpus + @hermes-1/Q2). They split into four surfaces that cannot be reviewed
as one change: **CLI enforcement** (@codex-1's 23 + @claude-1/F3 = 24), **protocol text** (@zcode-1's 7 +
@claude-1/F2 = 8), **the skill package** (@kimi-1's 4), and **this deck's own artifacts**
(@hermes-1/Q2 = 1). `24 + 8 + 4 + 1 = 37`. The owner has ruled that no complete release ships until the command defects are
fixed.

**CORRECTED (design signoff round).** The old split read `21 + 1 + 7 + 1 + 3 = 33` and silently
dropped @hermes-1/Q2, which appears in §2 and belongs to no other surface. The surfaces now sum to
`24 + 8 + 4 + 1 = 37`. It also predates the
three findings §3 mis-recorded. @codex-1, @opencode-1, @zcode-1 and @kimi-1 all flagged the
arithmetic.

**Reconciling the release rule with what shipped.** Three of the deferred findings — @codex-1/F6,
F8 and F14 — are CLI-enforcement defects, so "no complete release until the command defects are
fixed" and "5 findings deferred" have to be read together rather than as a contradiction. The rule
the deck actually applied: **every command defect that produces a WRONG RESULT is fixed; the three
deferred ones need a design decision, not a patch,** and each has its reason recorded in
`IMPLEMENTATION.md` — F6 needs a designed semantic signal for an adversarial alternative (substring
inference would not safely enforce §15.6), F8 is a missing collapsed fast-track close path, F14
needs per-knob precedence because blindly applying standard defaults would overwrite explicit idea
configuration. @codex-1 independently rechecked all three deferral reasons in the review cycle and
upheld each. The owner's ruling is a gate on the release, and it is the owner's to lift — this
consensus does not lift it.

## Signoffs

_Append below, or write `signoff-<agent>.md` for verbatim concatenation. Sequential only._

## Dispositions added after the design signoff round

The design consensus was drafted 2026-08-20, went to signoff on 2026-08-21, and **was blocked by
@codex-1 and @kimi-1 and reserved by @opencode-1 and @zcode-1** — not on the audit's substance,
which all four upheld, but on its ledger: false verdict counts, false attributions, and confirmed
findings with no disposition. Phase 7 of this same idea had already had to enforce "silence is not
a disposition" on itself; Phase 3 failed the same way.

Every finding below was **reproduced live at HEAD before being fixed**, and each fix has a test
that fails when that fix alone is reverted.

| Finding | Status | Fix |
| --- | --- | --- |
| **@codex-1, round-02 addendum** — standalone `preflight` ignores the canonical roster | **CONFIRMED, FIXED** | It probed every installed adapter family, so on this deck it reported `codex, claude, agy, hermes, kimi, opencode, zcode` — including `agy`, which is **not in the roster** — and printed "Ready: no pending gates". §9.0 requires probing every ROSTERED participant. `parley run` was already correct. `rosterParticipants()` |
| **@codex-1/F23** — the implementation gate accepts an unknown status and an empty artifact | **CONFIRMED, FIXED** | `status: banana` validated; `## Summary of work` was matched as a substring. Now a closed vocabulary plus the review gate's "a heading is not content" rule. Zero live artifacts newly rejected (measured across 72) |
| **@kimi-1/F4** — `status` tells a `source` deck to adopt the older packaged protocol | **CONFIRMED, FIXED** | Nothing in the skill package read `protocolRole`. `recommendedActions` now branches on it; consumer and unknown-role advice unchanged |
| **@claude-1/F2** — the bootstrap instruction breaks the repo's own drift guard | **CONFIRMED, FIXED** | `COOPERATION.md:59,136,1101` tell you to run `parley roster render`; doing so made `TestEmbeddedDefaultMatchesLiveDeck` fail closed, because the generated §2 table does not reproduce the anchor's hand-typed padding and adds a `State` column. Anchors now match on column signature |
| **@claude-1/F3** — `masked-by-env` is in the closed STATUS vocabulary and nothing emits it | **CONFIRMED, FIXED** | `roster set` printed it once to stderr; `roster show` had no path to it. `config.RosterMaskedFields` resolves it. Only fields whose value actually changes across layers count |
| **@codex-1/F15** | **CONFIRMED, ALREADY FIXED** | Applied during implementation even though §3 had excluded it. The disposition was right and the ledger was wrong — recorded here so the two agree |
| **@codex-1/F20** | **CONFIRMED, ALREADY FIXED** | Fixed in code and tests but missing from every Fixed table in `IMPLEMENTATION.md`. @opencode-1 and @kimi-1 both caught the omission |
| **@kimi-1/F5** | **DISMISSED, not deferred** | §3 records it REFUTED (by @codex-1) and CONFIRMED (by @zcode-1). `IMPLEMENTATION.md` listed it as *deferred*, which implies an intent to do it. It is **contested**, and a contested finding is not a backlog item |

**One more thing this deck's own metadata got wrong.** `parley-deck/meta/version.json` never
declared `protocolRole`, so `parley preflight` raised an `unknown-role` gate against the repository
that **authors** the protocol, and the skill's advice had no role to branch on. It now says
`source`; preflight reports `source-advisory` with no gates.

**Deferred, unchanged:** @codex-1/F6, F8, F14 and @kimi-1/F1 — reasons in `IMPLEMENTATION.md`,
independently upheld by @codex-1 in the review cycle.

**Open protocol follow-up (deliberately not settled by a validator).** `COOPERATION.md:454`
documents the IMPLEMENTATION.md status vocabulary as `implemented | fix-up-cycle-N | complete`,
but 5 live files in this deck use `ready-for-review` or `in-progress`. The new gate accepts all
five values so it rejects no existing work; **which set is correct is a protocol question and
needs its own idea.**

**Reviewed commit for this amendment:** `815c93ac3ec8d546e572ddd26d91cf5a29e2ad69`.

## One reproducibility fact this consensus must carry

`go test ./...` exits **0** on an unsandboxed macOS host and exits **1** inside a Codex Seatbelt
sandbox, failing only `internal/runner.TestDurableKillEndToEndRealProcess` with
`process verification failed (no recorded boot id); not killed`.

**This is an environment restriction, not a code defect.** The test reads the host boot id through
`sysctl kern.boottime`, which seatbelt denies (`Operation not permitted`), so `procctl` cannot
verify the process identity and — correctly, by design — refuses to kill. @codex-1 diagnosed this
itself in the review cycle, worked around it with a two-command PATH shim, and recorded it; it is
also documented at `docs/agent-cli-mechanics.md:31-32`.

It is written here because @codex-1 then re-encountered it in the design signoff round without the
shim and **blocked the consensus on it**, having no reason at read time to connect a red suite to a
sandbox limitation. A fact that costs a reviewer a whole cycle to rediscover belongs in the
artifact, not only in a docs file and a previous signoff.

**A reviewer running the suite inside a sandboxed agent must either** supply a `sysctl` shim for
`kern.boottime`, or run `go test ./... -skip TestDurableKillEndToEndRealProcess` and say so.
