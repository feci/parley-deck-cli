---
idea: meta-protocol-change-preflight-readiness
title: "Pre-idea readiness check: protocol auto-freshness + roster liveness ping with user-confirmed roster gates"
kind: meta-protocol-change + tooling
author: claude-1
status: round-01
transport: github-pr
participants: [claude-1, codex-1, hermes-1, antigravity-1]
roles:
  claude-1: architecture & correctness (lead synthesis)
  codex-1: risk & edge-cases / feasibility
  hermes-1: performance & simplicity
  antigravity-1: docs / consistency / UX
created: 2026-06-19
---

## Problem / idea
Make two things automatic at the start of every idea so the operator never has to
reason about session freshness or roster availability:
1. **Protocol auto-freshness** — facilitator checks for a newer installed protocol and
   brings the session onto it (auto for additive, confirm for breaking), preserving
   project-specific zones; off in the protocol *source* repo.
2. **Roster liveness ping** — facilitator pings all rostered agents; unavailable ones
   are excluded **per-idea** only with **user confirmation**, and re-including a
   previously-excluded agent also needs **user confirmation**.

Deliver as protocol text (§9.0 + §5 + §7 carve-out, both COOPERATION.md copies) **and**
a `parley preflight` CLI.

## Constraints
- Operator already LOCKED 3 decisions (see `reference/design-brief.md`): freshness =
  auto-additive/confirm-breaking; exclusion = per-idea temporary; scope = protocol +
  `parley preflight` CLI. Round-01 refines mechanics, does NOT relitigate these.
- Respect invariants: §1 non-solo, §5 quorum, §7 protocol-change gate, drift-guard
  lockstep + genericized embedded default, append-only signoffs.
- Must NOT regress the source repo (the source-vs-consumer inversion is real here).
- Unattended `parley run` / auto-drive must not deadlock on the new gates.

## Non-goals
- No confidence/majority gates, no auto-overwrite without zone preservation, no silent
  roster change, no change to phases/artifact shapes beyond the readiness check.

## Process
- Phase 1: independent `round-01/<agent-id>.md` (per the lenses above; read
  `reference/design-brief.md` first).
- Phase 2/3: cross-review → `consensus.md` + signoffs.
- Phase 4-8: FINAL.md → implement (protocol + `parley preflight`) → review → release.
- Pre-idea: roster pinged (dogfooding Feature 2); record any confirmed exclusion here.

## Readiness (this idea) — recorded by claude-1, 2026-06-19
- **Protocol freshness:** source repo → auto-adopt OFF (advisory); proceeding on the
  live canonical deck (it is ahead of the published skill — the source-vs-consumer case).
- **Roster ping (real PONG round-trip, bounded):** claude-1 ✅ (facilitator) · codex-1
  ✅ PONG · hermes-1 ✅ PONG (GLM 5.2) · antigravity-1 ✅ PONG (agy, `--add-dir "$(pwd)"`
  + 1m bound). Full 4-participant quorum available; no exclusion gate triggered.
