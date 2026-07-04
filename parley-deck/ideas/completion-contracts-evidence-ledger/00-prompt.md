---
idea: completion-contracts-evidence-ledger
author: user
created: 2026-07-04
track: deliberation
participants: [claude-1, codex-1, hermes-1, antigravity-1]
roles:
  claude-1: facilitation + protocol-text coherence
  codex-1: driver/Go internals + failure modes
  hermes-1: protocol minimalism + backward compatibility
  antigravity-1: consumer-deck / DevX perspective
status: round-01
---

## Problem / idea

Inspired by Hermes Agent v0.18.0 ("The Judgment Release"), which shipped completion
contracts for its `/goal` loop: the user states what "done" looks like, and the loop
judges completion **against recorded evidence** (actually running the project's checks),
not against the model's own claim of success.

Parley Deck has the same gap. Today "complete" in Phase 8 is a declaration by the
implementer plus reviewer consensus. The driver runs checks (RunChecks), but:

1. There is no per-idea, machine-checkable statement of what "done" means.
2. Check results are not recorded as a canonical artifact — reviewers cannot see
   proof; they see the implementer's prose claim ("full suite green"), which has
   already caused a real overclaim incident in this deck (codex BLOCK on
   protocol-restructure-appendices round-01 was claim-accuracy, not a code defect).

Proposal to deliberate:

- **`done_when:` contract** — an optional block in `00-prompt.md` frontmatter (or a
  dedicated section) listing mechanically verifiable completion criteria: commands to
  run and expected outcomes (exit code 0, output pattern, file existence).
- **Evidence ledger** — a canonical artifact (e.g. `review/evidence.md`) written by the
  driver (or implementer tooling), recording for each criterion: the command, exit
  code, duration, and an output digest. Append-only per fix-up cycle.
- **Driver enforcement** — when a `done_when:` contract exists, the driver refuses to
  mark Phase 8 complete until every criterion passed in the latest evidence entry;
  reviewers receive the ledger as review input so "checks pass" claims are grounded.
- **Protocol text** — amend §5 (implementation) and §8 (fix-up/stopping) to define the
  contract, the ledger artifact, and the rule that evidence beats prose claims
  (extends the existing artifact-wins principle).

## Constraints

- Additive and backward compatible: ideas without `done_when:` behave exactly as today.
- All-track invariants preserved (non-solo, refutation-default review, §14 human brake,
  audit trail, English-only, no-secrets).
- Evidence must never contain secrets; digests/truncation must be safe by construction.
- The ledger is evidence, not judgment: reviewers still decide severity and acceptance;
  a green ledger does not auto-close review.
- Driver changes live in parley-deck-cli (Go); protocol changes go to BOTH
  COOPERATION.md copies (drift guard enforces).
- Commands in the contract execute locally — the design must state the trust model
  (same trust level as RunChecks today; no network side effects expected; contract is
  authored in the idea and visible in review).

## Non-goals

- No auto-close of review on green evidence (human/reviewer judgment stays).
- No remote execution service; local checks only.
- No change to strict_gate semantics (it remains a separate, stricter opt-in).
- No retroactive contracts for already-closed ideas.
