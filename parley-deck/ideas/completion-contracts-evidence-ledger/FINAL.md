---
idea: completion-contracts-evidence-ledger
status: final
drafter: claude-1
track: deliberation
date: 2026-07-04
participants: [claude-1, codex-1, hermes-1, antigravity-1]
---

## Decision

Deliver evidence-grounded completion by EXTENDING the existing `checks:` primitive, not
by adding a new `done_when:` key or a separate ledger artifact. Unanimous round-02
convergence (✅ ×4). Absent or scalar `checks:` is byte-for-byte today's behavior.

## Design (ratified in consensus.md)

1. **`checks:` scalar → optional named list.** `checks:` in `00-prompt.md` accepts a
   scalar command (today) OR a list of `{name, command}` criteria (each expects exit 0).
   The list form activates the completion contract. Never both shapes.

2. **`ReadChecksContract(ideaDir)`** (new, `internal/driver/transport.go`, yaml.v3 over
   the frontmatter block): returns absent/scalar → legacy signal; list → a validated
   `[]{name, command}`. Fail closed on malformed list (empty/dup names, missing command).

3. **Driver-populated ledger = the existing `## Validation evidence` section** of
   IMPLEMENTATION.md. For list-form `checks:`, `driverImplOps.RunChecks` loops each
   criterion (`sh -c`, cwd = repo root), captures exit code + duration + a
   **secret-scrubbed, ~100-line/4KB-truncated** output tail, and OVERWRITES the section
   with the latest per-criterion table each cycle (git history preserves prior cycles).

4. **One Phase-8 fail-closed veto** (`internal/driver/impl.go`): with list-form
   `checks:`, `status: complete` requires the latest run ALL-PASS at current HEAD. The
   veto can only fail a close claim, never auto-pass one (mirrors strict_gate). A failing
   criterion escalates via §14 stopping-judgment — no auto-retry loop; the failure
   message names the failing criterion + command so the author can fix it.

5. **strict_gate independent:** list-`checks:` is necessary (checks passed); strict_gate
   adds sufficient (clean full-scope review). Green evidence does NOT satisfy strict_gate.

6. **Protocol text (BOTH COOPERATION.md copies + skill fallback), scoped to the list
   shape:** LE-4 (~2 sentences: `checks:` accepts scalar or named list; list activates
   the contract and the driver writes `## Validation evidence` each cycle); Phase-5
   IMPLEMENTATION.md template (~1 sentence: driver-populated when `checks:` is a list);
   Phase-8 (~1 sentence: list-`checks:` completion requires a latest all-pass run).

## Verification (done criteria)

- Tests: scalar/absent `checks:` unchanged (byte-for-byte); list pass/fail; malformed
  list fails closed; evidence section written with bounded scrubbed output (no raw
  unbounded dump, no secrets); list-checks run before completion; strict_gate still
  requires its own clean review.
- `go build ./...`, `go vet`, `gofmt -l` clean; drift guard green (both COOPERATION.md
  copies identical outside the allowlist).

## Non-goals

No new `done_when:` key; no separate `review/evidence.md`; no append-only rule; no
matcher grammar (v1 exit-0 only); no auto-close on green evidence.

## Signoffs

<!-- each participant appends its own block -->
