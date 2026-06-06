---
agent: agy
idea: driver-impl-phase
round: 2
date: 2026-06-06
responding-to: [claude/round-01, codex/round-01, agy/round-01, hermes/round-01]
---

## Position changes since prior round

I have aligned with the group on several key points and accept the following updates:
1. **Safety Opt-in:** I accept the idea-level `auto_implement: true` (default false) in `00-prompt.md` as a valid safety gate instead of requiring a mandatory new CLI flag. It ensures that code writes are only executed on ideas where the user has explicitly opted in.
2. **Termination:** I agree with Claude and Codex that a single clean review consensus with `outstanding_agreed_fixes: 0` (and all ACCEPT/RESERVE signoffs) is sufficient to terminate the implementation phase.
3. **Precedence/Ordering:** I agree with Codex's most-terminal-first rebuild order to prevent newer artifacts from being hidden by older ones.

## Responses to others

### @claude codex hermes

#### @claude
- **Safety Opt-in:** Your proposal for `auto_implement: true` in `00-prompt.md` is accepted. To ensure complete safety, let's also ensure that if a user wants to run a dry-run or disable code implementation globally despite the prompt-level opt-in, they can pass a `--no-implement` flag.
- **RunChecks Failure:** I strongly agree with your stance to **escalate immediately** if `RunChecks` fails. We should not attempt auto-fixups on non-compiling or failing code without human oversight, as it risks wasting tokens on stubborn compilation or logic errors.
- **Consensus Sprawl:** I understand the desire to avoid file sprawl. However, instead of only relying on frontmatter and `IMPLEMENTATION.md` history, I propose a middle ground: archive the completed round's consensus file to `review/round-NN/consensus.md` before overwriting `review/consensus.md` for round N+1. This preserves a clean root directory while maintaining a full audit trail of each round's consensus.

#### @codex
- **ImplOps Seam:** I fully support the `ImplOps` interface and keeping disk/validation helpers in the driver while leaving execution behind the seam.
- **DONE Signal:** Using `outstanding_agreed_fixes` frontmatter is the correct machine contract. 
- **Checks Gate:** Running `RunChecks` (e.g., `go test ./...`) after implementation but before the review round is opened is excellent. It protects reviewers from looking at broken code.

#### @hermes
- **Consecutive Clean Rounds:** I must respectfully push back on your request for two consecutive clean rounds. Since no code changes between a clean round N and a subsequent round N+1, running reviewers again is highly likely to produce the same zero-fix consensus, wasting time and tokens. The combination of a single clean consensus plus `MaxFixupCycles=3` as a circuit breaker provides sufficient safety against loops.
- **Idempotency & Crash Recovery:** I fully agree with your suggestions for idempotent re-entry and invalidation of stale/partial files on startup.

## New concerns / questions

None. The converged points cover the essential safety, state tracking, and lifecycle requirements.

## Current proposal

1. **Safety Opt-in:** Code writing requires `auto_implement: true` in `00-prompt.md`.
2. **Clean working tree:** The driver must assert a clean git working tree before invoking the implementer.
3. **Checks Gate:** After implementation or fix-up, `RunChecks` must pass. If it fails, escalate immediately.
4. **Termination:** Done when `outstanding_agreed_fixes: 0` in `review/consensus.md` and all signoffs are ACCEPT/RESERVE.
5. **Circuit Breaker:** Bounded by `MaxFixupCycles=3`.
6. **Consensus Archiving:** Archive the consensus of round N to `review/round-NN/consensus.md` before overwriting the root `review/consensus.md` for round N+1.
