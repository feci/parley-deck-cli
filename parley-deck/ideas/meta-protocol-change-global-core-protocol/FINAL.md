---
idea: meta-protocol-change-global-core-protocol
status: final
drafted-by: claude-1
date: 2026-08-07
participants: [claude-1, codex-1, hermes-1, kimi-1]
absent: [opencode-1]
consensus: revision 3, accepted by codex-1, hermes-1, kimi-1
track: deliberation
---

# FINAL — global core protocol with local override/extension

Supersedes nothing yet in force; it defines the target state and the slice implemented now. The
full decision text, the resolved verdict conflict, and the drafter's position reversals are in
`consensus.md` (revision 3) and are not repeated here.

## The change in one paragraph

`COOPERATION.md` stops being a hand-edited per-deck store and becomes a **generated view**, exactly
as §2's roster table did. Authority moves to an **immutable, versioned, content-addressed core**
under `~/.parley/protocol/core/<version>/`. A deck may **replace** a block the core marks
replaceable and **extend** at one declared point, addressed by permanent block ID; everything else
is sealed. Each idea **pins** the effective protocol and keeps a materialized snapshot, so an open
idea completes under the version it started with while the next idea picks up the current one.

## Binding decisions

D1 immutable versioned core store in `~/.parley`; write-once releases.
D2 registry of permanent block IDs, sealed by default, tombstones on deletion, ID addressing.
D3 v1 open surface: replaceable `s6.6` (working language); one extension point `ext-1`; six
   identity slots; everything else sealed, including §7 and §15.
D4 one committed overlay per deck, two operations, operation-specific provenance, empty overlays
   forbidden.
D5 core-owned resolution order; a contradictory extension is incompatible, not higher-precedence.
D6 the deck keeps a committed, generated `COOPERATION.md`; render and check report what they change.
D7 the effective snapshot is ALWAYS materialized and is the sole protocol input for later phases.
D8 deck lock; a missing pinned release blocks adoption and rendering, never continuation.
D9 enforcement per launch path; no claim without a runtime proof; `DETECTED-UNATTRIBUTED`.
D10 registry-driven compatibility check, including extension dependency sets; never auto-migrate.
D11 §7 gains a blast-radius clause: a core change needs the meta idea AND user ratification.
D12 `protocolRole` retired, superseded by the deck lock.

## Implementation slice for this cycle

Ranked in `consensus.md` §3. This cycle ships **rank 1 and the protocol-text change**:

1. **`parley protocol` verbs** — `render`, `check`, and the core store they read.
2. **The §7 blast-radius clause (D11)** in all three `COOPERATION.md` copies, plus the changelog
   entry required by G6.

Ranks 2–4 (per-idea pinning, the overlay, the detection-layer enforcement) and DF-1/DF-2 remain
ratified and scheduled, not built here. **Nothing in the shipped slice may claim a guarantee it does
not implement** — G7b makes that binding.

## Gates the implementation must satisfy

G1 `protocol render` idempotent and reports every block replaced or removed.
G2 no autonomous/agent-accessible write path into the core; attended TTY publisher is the sole
   audited exception, creating a new release rather than modifying one.
G3 no confinement claim without the preflight write-probe passing, using the RESOLVED path and a
   dedicated non-release probe location.
G4 overlay targeting a sealed/missing/tombstoned ID, or a stale base hash, fails closed; empty
   overlay rejected.
G5 an idea's pinned protocol is what its later phases read; proven by test.
G6 this protocol change recorded in `meta/protocol-changelog.md` in the §7 template format.
G7 production call-site pin test, including continuation after deleting the global release.
G7b call-site truth — no guarantee documented as landed without an end-to-end test of the real
   entry point.
G8 lock byte-verification, scoped: adoption always; continuation only when the release is present;
   the snapshot always.

Gates G5, G7 and G8 bind the ranks not shipped this cycle; they apply when those ranks land.
