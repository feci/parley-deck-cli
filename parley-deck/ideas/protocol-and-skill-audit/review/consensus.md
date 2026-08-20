---
idea: protocol-and-skill-audit
review-cycle: 1
outstanding_agreed_fixes: 0
blocked: false
drafted-by: claude-1
date: 2026-08-21
reviewed-commit: 0bb99031aa1ddc1027605424848daa3a696e9f00
---

# Review consensus — cycle 1

All five non-implementers filed. Every reviewer built from source and ran both suites; three
built the PRE-FIX binary and diffed its verdicts against the new one across the whole deck.

## Agreed fixes — applied

**1. The implementer's own signoff malformed every review consensus it appears in.** MAJOR, found
independently by @kimi-1 and @zcode-1, each by building the pre-fix binary and diffing every
consensus in the deck. @kimi-1 counted 9 flips to `malformed`, @zcode-1 counted 23; two of the
affected ideas are in flight, so this blocked live work, not history.

@zcode-1 supplied the sentence that closes it: `COOPERATION.md:591`'s Phase-7 template reads
*"Each active participant (**implementer included**) APPENDS their signoff block."* The F2/F3 fix
forbade exactly what the protocol prescribes.

Cause: one list answered two different questions — **who is AWAITED** and **who may SIGN**. Fixed
by splitting `known` from `required`. Verified by sweeping every review consensus with both
binaries: **6 flips, all `partial → ready`** (F3 working as intended), **zero regressions**.

**2. The driver instructed its own FINAL drafter to produce an artifact its own gate rejects.**
MAJOR, @codex-1. `buildFinalDraftPrompt` named one section while `finalScaffoldReason` requires
seven, and never mentioned the `idea:` frontmatter the slug check reads. **This is the audit's own
defect class, committed while fixing it** — and it is the second time in this batch: the same trap
was avoided deliberately for F18 by fixing the prompt first. The prompt is now generated from
`protocol.RequiredFinalSections`, so it cannot drift from the gate again.

**3. The pipeline treated the first-step FINAL scaffold as a completed block.** MAJOR, @codex-1 and
@kimi-1. The scaffold carries `status: final` from step one, so `isFinalized`'s frontmatter-only
test called an empty outline complete and the pipeline advanced past it. `isFinalized` now applies
`protocol.FinalIsScaffold` as well, and `pipeline auto` reports the scaffold step as NOT complete
instead of announcing a closure that did not happen.

## Clean results, recorded because they are evidence

**@opencode-1 — no existing test was weakened.** Its whole slice was the question I could not
answer about myself. It gave a per-file verdict with commit hashes for all seven modified test
files, and its reasoning is specific: `TestFinalScaffoldReason`'s directory became `demo` so the
new slug check would not false-fail; the old fixtures encoded *"one padded heading is enough"*,
which was the bug F22 fixed; rejection cases still reject and acceptance cases still accept.

**@zcode-1** — `go test ./...` green across 26 packages and `npm test` 388/0, built from source,
independently of the implementer's runs.

## Dismissed, with the reason

**@hermes-1 Finding 1 — "3 closed ideas still carry a non-terminal status".** NOT UPHELD.

It reported `tui-action-execution: accepted`, `protocol-overlay-local-extension: **` and
`meta-protocol-change-rho-retrospective-optimization: final-design-for-review`. At the reviewed
commit `0bb9903` all three carry `status: final` in `00-prompt.md`; verified with
`git show 0bb99031:…/00-prompt.md`.

The values it reported are the `status:` of those ideas' **FINAL.md**, not their `00-prompt.md`.
The finding measured a different file from the one the fix repaired and from the one §6 rule 5
reads.

**It is still worth something:** `FINAL.md` files carry their own `status:`, and `tui-action-execution`
says `accepted` while `meta-protocol-change-rho-retrospective-optimization` says
`final-design-for-review` — neither is `final`, and `finalScaffoldReason` requires `status: final`.
Two closed ideas hold a FINAL.md the current gate would reject. That is a **new observation about
FINAL.md**, not the finding as filed, and it belongs to a follow-up rather than this cycle.

## Deferred follow-ups

- **@hermes-1 Q2** — 238 of 659 review artifacts carry no `reviewed-commit`. Historical and not
  revalidated, but it measures how much of this deck's own record predates the rule just imposed.
- **Two FINAL.md files whose own `status:` the FINAL gate would reject** (above).
- The five findings deferred in the first batch — @codex-1/F6, F8, F14 and @kimi-1/F1, F5 — each
  with its reason recorded in `IMPLEMENTATION.md`. No reviewer challenged any of those reasons.

## Signoffs
