---
idea: protocol-and-skill-audit
status: final
author: claude-1
consensus-date: 2026-08-21
participants: [claude-1, codex-1, hermes-1, kimi-1, opencode-1, zcode-1]
signoffs: [claude-1 ✅, codex-1 ✅, hermes-1 ✅, kimi-1 ✅, opencode-1 ✅, zcode-1 ✅]
---

# FINAL — a six-agent audit of the protocol and the skill, and the three correction passes it took to close

## Final plan / specification

Six agents audited `parley-deck-cli`, `COOPERATION.md` and `parley-deck-skill` for defects. **47
findings were filed; 37 are confirmed, 11 contested, and 33 of the confirmed set are fixed.** Four
are deferred with recorded reasons. The work shipped as fix-up cycles against `IMPLEMENTATION.md`,
was reviewed in `review/` to a unanimous consensus, and is described commit-by-commit in
`IMPLEMENTATION.md` and `CHANGELOG.md`.

The audit's single largest result is not any one finding. It is a **defect class** that recurred in
nine of the fixes and three times more in the audit's own artifacts:

> **A printed rule binds only where enforcement lives.**

Its sibling appeared almost as often: **a message that misstates its own effect.** Concretely:

- `consensus finalize` and the driver each owned a FINAL validator; only one checked the declared
  slug and status. `protocol.ValidateFinal` is now the single gate.
- The driver validated `reviewed-commit` on review artifacts; `consensus draft --review` did not.
  `protocol.ValidateReviewArtifact` now binds both.
- `parley run` probed the roster in preflight; standalone `parley preflight` probed every installed
  adapter family, including one removed from the roster, and reported "Ready".
- `COOPERATION.md` instructs the bootstrap to run `parley roster render`; doing so broke the repo's
  own drift guard.
- `masked-by-env` was in the closed STATUS vocabulary with nothing able to emit it.
- The unknown-freshness gate's own displayed `--yes` remedy had no branch for it, so a freshly
  initialized deck could never be reported ready.

## Purpose / user-visible outcome

Gates that used to accept emptiness now require content; gates that bound one entry point now bind
both. A user sees: `consensus draft` refusing a blank round, `finalize` refusing another idea's
FINAL, `draft --review` refusing a review that names no tree, `preflight` probing the actual roster,
`roster show` able to report `masked-by-env`, `init` writing real provenance instead of
placeholders, and `preflight --yes` able to clear the gate it raises on a new deck.

Nothing here is a feature. It is enforcement catching up with what the protocol already said.

## Context & orientation

- Findings: `round-01/*.md`. Adversarial verification: `round-02/*.md`.
- Ledger: `consensus.md` §1–§5, plus the correction sections appended to it.
- Work: `IMPLEMENTATION.md`. Review: `review/consensus.md` (triage `ready`, five of five ✅).
- The verification design was asymmetric **by construction and this is recorded, not hidden**:
  @codex-1 wrote 24 of the 47 findings, so it was the only verifier not permitted to assess them,
  and the two verifiers who did assess them confirmed 66 of the 74 findings they touched (~89%).
  All nine REFUTED verdicts in the corpus are @codex-1's.

## Observable acceptance criteria

- `go test ./...` — 27 packages, exit 0, on an unsandboxed host (see the known limitation below).
- `npm test` in `parley-deck-skill` — 391 pass / 0 fail, exit 0.
- `review/consensus.md` triage `ready`, zero outstanding agreed fixes.
- `consensus status protocol-and-skill-audit` triage `ready`, six of six ✅.
- Every fix has a test that fails when **that fix alone** is reverted, checked one at a time.

## Idempotence & recovery

Every fix is a validation gate or a resolver; none migrates data. Re-running any command is safe.
Two changes touch existing state and both were measured before enforcing: the implementation-status
vocabulary rejects **zero** of this deck's 72 live `IMPLEMENTATION.md` files beyond what the old
check already rejected, and `parley init` writes `protocolSha256` only for decks it creates. The
`reviewed-commit` gate binds new drafts only — rounds whose consensus already exists are never
revalidated.

## Known risks / de-risking

- **The stricter gates can reject historical artifacts.** Measured for each: the implementation gate
  adds zero new rejections; the review gate binds only new drafts; the FINAL gate rejects two closed
  ideas' `FINAL.md`, recorded as a deferred follow-up.
- **`preflight --yes` clears the freshness gate per confirmation, not permanently**, when the
  installed skill exposes no packaged protocol body. That is the deliberate price of refusing to
  invent a packaged hash — the alternative recreates the original defect with a fabricated value.
- **Known environment limitation:** `go test ./...` exits 1 inside a Codex Seatbelt sandbox, failing
  only `TestDurableKillEndToEndRealProcess`, because seatbelt denies `sysctl kern.boottime` so
  `procctl` cannot verify process identity and correctly refuses to kill. Not a code defect. Shim
  the sysctl or run with `-skip TestDurableKillEndToEndRealProcess`, and say which.
- **Deferred, unchanged:** @codex-1/F6, F8, F14 and @kimi-1/F1 — reasons in `IMPLEMENTATION.md`,
  independently rechecked and upheld by @codex-1.
- **Open protocol follow-up:** `COOPERATION.md:454` documents the IMPLEMENTATION.md status
  vocabulary as `implemented | fix-up-cycle-N | complete`, but five live files use
  `ready-for-review` or `in-progress`. The gate accepts all five so it rejects no existing work;
  reconciling the two needs its own idea.
- **Cross-repo follow-up:** `parley-deck-cli` and `parley-deck-skill` are separate git repositories.
  A disposition citing a commit must name the repo, and a fix spanning both produces two commits.
  Nothing enforces that yet — and this idea got it wrong once.

## What this idea proved about its own method

The design consensus needed **three signoff rounds** and every one found something real:

1. The ledger reported **verdicts nobody cast** — @kimi-1's row read `37/8/2/1`, summing to 48
   against 42 assessed, crediting it with two REFUTED and one UNREPRODUCIBLE it never issued. It
   also mis-recorded three of @zcode-1's CONFIRMED verdicts as PARTIAL-only, which silently dropped
   two confirmed findings out of the fix list entirely.
2. The **correction had its own off-by-one**, folding a supplemental measurement into a corpus count.
   Three reviewers caught it independently by re-deriving the numbers themselves.
3. A disposition **cited a commit that did not contain the fix** — the fix was in the other
   repository. The `npm test` result had been measured on uncommitted code.

The review consensus needed three fix-up cycles and had the same shape: it was **silent on four of
@codex-1's six MAJOR findings**, and Phase 7 had to enforce "silence is not a disposition" on
itself before it could enforce it on anything else.

**The method that worked was not agreement.** It was requiring every reviewer to re-derive the
numbers instead of reading them, and to mutation-check a fix instead of reading its diff — @codex-1
found a fix with no regression test that way, and @kimi-1 found a commit that did not contain what
it claimed. **Two reviewers agreeing proved nothing here; a reviewer who ran the command proved
everything.** Where a reviewer reported PRIMARY evidence it had not executed, that signoff's
conclusions were discarded and the fact recorded — twice, for the same participant.

## References

- `consensus.md` — the ledger and its three correction passes
- `IMPLEMENTATION.md` — every fix, its commit, and the mistakes made while fixing
- `review/consensus.md` — the Phase 7 record and the four MAJORs it initially dropped
- `round-01/*.md`, `round-02/*.md` — findings and adversarial verification
- `signoff-*.md` — six participants' design signoffs, three rounds
- `docs/agent-cli-mechanics.md:31-32` — the Codex Seatbelt `sysctl` artifact
- `CHANGELOG.md` 1.46.0
