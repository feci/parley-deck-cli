---
agent: claude-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-10
worktree-head: d05fd55
worktree-head-provenance: launcher-supplied snapshot; SECONDARY, not independently confirmed (no shell in this session)
artifact-kind: implementation handoff — protocol-text slice only
not-a-signoff: true
not-a-phase-6-review: true
not-an-acceptance: true
files-changed:
  - parley-deck/COOPERATION.md (§4.0.1 LE-7/LE-11 line; §4 Phase-8 "Close-decision integrity")
  - internal/protocol/defaults/COOPERATION.md (byte-identical mirror of both edits)
  - parley-deck/meta/protocol-changelog.md (one new UNRELEASED entry, prepended)
files-not-touched:
  - internal/driver/**, internal/app/**, internal/runner/**, internal/fsutil/** (other owners)
  - FINAL.md, consensus.md, signoffs, earlier handoffs, meta/packet-applicability.yaml
  - parley-deck-skill reference protocol copy (separate worktree, outside this slice)
resolves: MAJOR "Fail-closed goal check contradicts the protocol text in force"
  (./claude-1-runtime-review-20260910.md) — the protocol-text half only
---

# Protocol close-integrity fix — §4 LE-7 goal-done check

This closes the protocol/code contradiction I filed as a MAJOR in
`claude-1-runtime-review-20260910.md`: the in-tree driver fails closed on a broken
goal check while the ratified §4 text said the goal check is "fail-open on its own
error". I own the protocol text; I do not own the code, and nothing here certifies
it. This is not a signoff, not a Phase-6 review, and not an acceptance of any
acceptance criterion.

## 1. Exact changed obligations

Two places in Section 4 changed, identically in both protocol copies.

**§4.0.1, the LE-7/LE-11 plain-English line.** Was:

> Before an auto-driven close, a goal-done check verifies FINAL's observable
> acceptance criteria; reservations or too-few reviewers escalate rather than close.

It now adds, between those two clauses: the check "can only withhold a close, never
establish one — a missing, self or unavailable checker, a failed run, or an
inconclusive or reserved verdict leaves completion unverified and escalates, and a
textual pass never replaces current-tree criterion evidence."

**§4 Phase 8, "Close-decision integrity (LE-7/LE-11)".** The removed sentence was:

> The goal-check is defense-in-depth on top of the review consensus and fail-open on
> its own error (a broken or inconclusive checker never blocks a review-clean idea).

Four obligations replace it. Completion is **unverified**, and the driver escalates
for a human decision instead of passing, when:

| # | Condition | Previously |
| --- | --- | --- |
| 1 | The checker is missing, is the implementer, or cannot be resolved and launched | fail-open: close proceeded |
| 2 | The check execution fails or exits non-zero | fail-open: close proceeded |
| 3 | The verdict is inconclusive (no parseable verdict) | fail-open: close proceeded |
| 4 | The verdict is only a pass-with-reservations, not an unqualified pass | not addressed by the old text |

And one non-substitution obligation, which is the part that binds to FINAL D3: a
textual goal-check verdict "is defense in depth on top of the review consensus and
never substitutes for the current-tree independent criterion evidence a close already
requires" — a self-issued verdict, a stale code tree, a skipped or no-execution
report, a missing criterion, or a partial original scope cannot close an
implementation. A textual PASS therefore cannot be offered as the criterion evidence;
it can only sit on top of it.

The paragraph also states the trade in one sentence, because it is a real operator-
visible consequence rather than a new mechanism: an unavailable checker now halts a
review-clean close until a human restores an independent checker or rules on it.

## 2. What deliberately did not change

- **The conditional trigger.** The goal-done check still runs only "under
  `auto_implement` or `strict_gate`". No other idea acquires a goal check.
- **Design-only routing stays lighter.** "A design-only idea keeps the lighter close
  (conditional rigor)" survives verbatim as the paragraph's last sentence, and the
  `ACCEPT-WITH-RESERVATIONS` and fewer-than-two-independent-reviewer refusals remain
  scoped to `auto_implement` exactly as before.
- No §4.0 track-table cell, quorum rule, signoff rule, reviewer count, fix-up budget
  or escalation mechanism was touched.
- No heading was added, removed or reworded, so every
  `meta/packet-applicability.yaml` locator (which addresses heading lines) still
  resolves. The map itself is unchanged — changing a classification would be its own
  §7 act.
- FINAL.md, consensus.md, every signoff and every earlier handoff are untouched.

## 3. Known compatibility impact on strict / auto close

This is a behavioural tightening for ideas that already close through the driver, and
it can strand a review-clean idea. Stated plainly so nobody discovers it at a close:

- **`strict_gate: true` alone is enough to reach the gate — no `auto_implement`
  needed.** `internal/driver/close_integrity_test.go:87-105`
  (`TestStrictDesignOnlyCompletesAndRunsGoalCheck`) sets `strict_gate: true` with a
  reserved triage and **one** reviewer, and asserts the goal check runs and the idea
  completes (PRIMARY, current worktree). Under the new obligations that same idea
  escalates instead of completing whenever conditions 1–4 hit.
- **The two most likely real triggers are environmental, not substantive:** a
  reviewer whose CLI cannot be resolved at close time, and a checker that exceeds the
  bounded one-shot deadline (`internal/app/driver_impl.go:400`, 2 minutes). Neither
  says anything about the implementation's quality, and both now escalate on every
  driver tick until a human intervenes.
- **The resolution path is human, by design:** restore an independent checker, or
  record an operator ruling. The protocol offers the driver no waiver, which is the
  point of the change — but it means an operator with no second agent available
  cannot auto-close a `strict_gate` idea at all.
- **Migration question still open** (carried over from my review, unanswered): does
  this reach existing decks mid-flight, or only ideas opened after the change? Both
  protocol copies are read live, so my expectation is "immediately, all decks" — that
  is an owner expectation, not a verdict, and it decides whether this is a migration
  hazard or only forward-looking. It is a fair question for the operator.
- The embedded default (`internal/protocol/defaults/COOPERATION.md`) is what
  `parley init` writes into a **new** deck, so new decks get the tightened text from
  the next build; no existing deck file is rewritten by this change.

## 4. Code correspondence (PRIMARY, non-owner reads)

I read the current worktree source; I ran nothing. The updated text matches what the
code already does:

- `internal/app/driver_impl.go:381-383` — `checker == "" || checker == o.implementer`
  → `false, "goal-check has no independent checker"` (obligation 1).
- `:384-387` — `agents.ResolveParticipant` error → `false, "goal-check checker
  unavailable"` (obligation 1).
- `:390-392` — evidence-directory `MkdirAll` failure → `false` (obligation 2).
- `:406-408` — `res.ExitError != "" || res.AgentExit != 0` → `false, "goal-check
  checker failed; completion is unverified"` (obligation 2).
- `:409-416` — `FAIL` → false with the answer; `PASS` → true; **default → `false,
  "goal-check inconclusive; completion is unverified"`** (obligation 3).
- `:441` — the verdict must be exactly `PASS`; a qualified pass falls through to the
  inconclusive branch (obligation 4).
- `internal/driver/impl.go:272-276` — the gate is `d.cfg.AutoImplement ||
  d.cfg.StrictGate`, and `!ok` returns `ActionEscalated` (the preserved conditional
  trigger, and escalate-not-pass).
- `internal/driver/impl.go:32-36` — the `ImplOps.GoalCheck` interface doc already
  states the rule the protocol now carries ("Missing, failed or ambiguous execution
  returns false … never replaces typed criterion evidence").

**Two code items remain open and are not mine to fix:**

1. `internal/driver/impl.go:271` still carries the stale comment "a checker error is
   advisory (fail-open inside GoalCheck)", six lines from the interface doc that says
   the opposite. It now contradicts the protocol too. Owner: the driver slice.
2. `parseGoalVerdict` last-wins aggregation (my separate MAJOR) can still turn a
   stated FAIL into PASS on a trailing template echo. The launch brief reports codex-1
   is correcting mixed-verdict aggregation; **I did not verify that correction and do
   not certify it.** Until it lands, obligation 3's protocol text is stricter than the
   parser's behaviour on that one input shape — the protocol is the safe side of that
   gap, not the unsafe one.
3. `IMPLEMENTATION.md:132` still records "Deviations from FINAL.md: None" for the
   fail-closed code change. That file is not mine to edit; with this §7 text change
   landed the code is no longer a deviation, so the honest fix there is a reference to
   this changelog entry rather than a deviation record. Flagging, not doing.

## 5. Checks I actually performed — and did not

Performed:

- PRIMARY reads of `FINAL.md` (D3, AC-E1/AC-E2), my own
  `claude-1-runtime-review-20260910.md`, both `COOPERATION.md` copies,
  `meta/protocol-changelog.md`, `meta/packet-applicability.yaml:87-95`,
  `internal/app/driver_impl.go:360-448`, `internal/driver/impl.go:20-48,259-286`,
  `internal/driver/close_integrity_test.go:87-105`.
- Confirmed by grep that this repo contains exactly two `COOPERATION.md` files, so
  the mirror set is complete **within this repo**.
- Compared the two edited regions line by line after editing: they render
  identically, at the same pre-existing 7-line offset. *This is an owner statement
  about my own edit, not a verdict (§15.1); the drift guard is the independent check.*
- Confirmed the applicability map keys on heading lines
  (`meta/packet-applicability.yaml:14`), and that I changed no heading.

**Not performed — do not read any of these as passed:**

- No `go build`, no `go test`, no `go test ./internal/protocol -run Drift`. **No shell
  was available in this session.** `internal/protocol/drift_test.go` is the machine
  check that both copies match outside the allowlisted zones, and it has not been run
  against this edit by me.
- No `parley protocol packet check` run, so the packet renderer has not been re-checked
  against the edited source.
- No git operation of any kind; `d05fd55` is the launcher's word, not my observation.
- Nothing about AC-E1/AC-E2 typed criterion evidence, AC-P1/P2, or any other acceptance
  criterion is assessed or claimed here. This slice changes protocol text; it delivers
  no criterion evidence and closes no gate.

## 6. What a reviewer should check

1. Run `go test ./internal/protocol/...` — the drift guard is the real proof the two
   copies match; my line comparison is not a substitute.
2. Read the new Phase-8 paragraph against `FINAL.md` D3 and decide whether the
   non-substitution sentence is faithful to "A self verdict, stale code-tree,
   skipped/no-execution report, missing criterion, or partial original scope cannot
   close the whole implementation".
3. Decide whether the one-sentence statement of the operator trade-off (§1, last
   paragraph) belongs in the protocol or should be dropped as commentary. I judged it
   descriptive of an existing mechanism rather than a new obligation; that judgement
   is challengeable and I will not defend it as settled.
4. Confirm the `#### Stopping judgment` block carrying this paragraph is
   `include: when, phases: [7, 8]` in the applicability map
   (`meta/packet-applicability.yaml:92-95`) — an optimized packet outside phases 7–8
   omits it under its existing trigger, "any close, cap or trajectory decision on an
   implementation". Unchanged by me, but it is the block whose obligations just got
   stricter, so it is worth a second reader's eye.
5. Propagate the same two edits to the parley-deck-skill reference copy. **Not done
   here** — it lives in a separate worktree outside this slice.
