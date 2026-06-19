---
idea: meta-protocol-change-preflight-readiness
status: final
author: claude-1
consensus-date: 2026-06-19
participants: [claude-1, codex-1, hermes-1, antigravity-1]
quorum-note: "signoffs claude-1 ✅, codex-1 ✅ (re-sign after block fix), hermes-1 ✅; antigravity-1 waived per operator ruling (agy hung on signoff append) — its round-01 on record."
---

## Purpose / user-visible outcome
At the start of every idea, `parley` automatically (1) checks protocol freshness and,
in a *consumer* project, brings the deck onto a newer installed protocol (auto for
additive, confirm for breaking) preserving project zones; and (2) hosted-PONG-pings the
roster, gating any per-idea exclusion / re-inclusion behind explicit user confirmation.
Delivered as protocol text (§9.0 + §5 + §7 carve-out, both COOPERATION.md copies) and a
`parley preflight` CLI wired into `parley run` before idea creation — never deadlocking
unattended runs.

## Context & orientation
Build on existing code (file refs are starting points — verify at implementation):
`internal/app/version_status.go` (`parleyDeckSkillStatus`, 10s-bounded);
`internal/agents/discover.go` (`Discover`, LookPath + `--version` 4s); the `run`
command + `runcontrol.Create` + `StartAutoAnswerer` in `internal/app/app.go`;
`internal/runner/supervision.go` + `failclass.go` (mid-idea watchdog);
`internal/protocol/` (drift guard + allowlist); `parley-deck/meta/version.json`
(`protocolSha256`, `packagedProtocolSha256`, `deckVersion`) — add `protocolRole`.

## Final specification

### Part A — Protocol text (both `COOPERATION.md` copies, drift-guard lockstep)
1. **New §9.0 "Pre-idea readiness check"** (runs as step 0 of §9, at idea start):
   - **Freshness:** compare live `protocolSha256` vs installed `packagedProtocolSha256`
     (via `parley-deck-skill status`). `protocolRole: source` → advisory only, **never
     writes**. `consumer` + newer installed protocol → **additive** (semver minor/patch
     of `deckVersion`) auto-syncs the protocol body **preserving project zones**
     (header, §0 transport, §2 roster — the drift-guard allowlist), records
     `Protocol synced:` + `meta/protocol-sync_<ISO-ts>.md`; **breaking** (major) →
     pause for user confirmation. Missing/unknown `protocolRole` → do NOT auto-write;
     one-time confirm + backfill.
   - **Roster ping:** hosted-PONG round-trip per rostered participant via its real
     configured invocation (operator ruling), `command -v` pre-check first, ~90s bound,
     `agy` gets `--add-dir`, concurrent behind a global deadline, process-group kill on
     timeout. Available = exits in time with the exact sentinel; else unavailable.
   - **Gates (both user-confirmed, per-idea):** exclude an unavailable agent → confirm
     + record in `00-prompt.md`; re-include a previously-excluded now-available agent →
     confirm. Quorum locks at Phase 0. A mid-idea hang falls to §5 + the supervision
     watchdog and the same per-idea waive.
2. **§5** add: "Quorum for an idea is set at Phase 0 from the §9.0 readiness check;
   agents excluded there (user-confirmed) do not count toward this idea's quorum; the
   quorum locks once Phase 0 completes. Excluding the last non-facilitator still
   requires the §1 user-authorized solo exception."
3. **§7 carve-out**: "Applying an upstream-ratified protocol version via the §9.0
   freshness sync — when it is additive/compatible and preserves project-specific
   zones — is a maintenance sync, not a protocol change, and does not require a
   meta-protocol-change idea."

### Part B — `parley preflight` CLI + `parley run` wiring
- `parley preflight [--dir DIR] [--json] [--yes] [--ping-timeout D] [--no-ping]`.
- ONE shared `preflight(root, opts) (Report, []Gate, error)` used by the standalone
  command and `parley run`.
- **Roster-ID ↔ runtime-ID map** from `meta/headless-agents.local.json` (the `agents`
  keys are roster IDs; each entry's `cli` is the runtime). Reports + `00-prompt.md` use
  roster IDs; an entry with no runnable `cli` → `not_configured` (never name-matched).
- **Exit codes:** `0` ready · `3` pending gate (names gate + prints confirm command) ·
  `1` hard failure (no workspace, or excluding would leave <2 participants → §1) ·
  `2` usage.
- **Wiring:** `parley run` calls `preflight` **before** `runcontrol.Create`, reusing the
  already-discovered agent set. Attended → route gates via existing HITL (`parley
  answer`)/`confirmLaunch` and stop. Unattended (`--auto`, no TTY) → **hard-stop
  non-zero; never read stdin, never auto-exclude, never auto-adopt a breaking bump.**
  The exclude / re-include / breaking-freshness gates are **excluded from
  `StartAutoAnswerer`'s auto-answer set**. Additive consumer freshness sync is the one
  write preflight may do without confirmation (gated by `protocolRole`/source). `resume`
  / `continue` → roster re-check only (no freshness, no gates).
- `--no-ping` → Tier-0 presence-only (the cheap fallback).

### Part C — `meta/version.json`
- Add `protocolRole: "source" | "consumer"`. **This repo = `"source"`** (so freshness
  is advisory here and never regresses the canonical deck). `parley init` writes
  `"consumer"`. Absent → confirm + backfill.

## Observable acceptance criteria
1. `parley preflight` in THIS (source) repo: prints freshness = advisory/source (writes
   nothing to `COOPERATION.md`), hosted-PONG roster table, exit `0` when all available.
2. With an agent unavailable: exit `3`, names the exclude gate, prints the confirm
   command; with `--yes`-style confirm path records `excluded:` in `00-prompt.md`.
3. Unattended `parley run --auto` with an unmet gate: hard-stops non-zero **without
   creating a half-open idea** and without hanging.
4. `go test ./...` green incl. drift guard (`TestEmbeddedDefaultMatchesLiveDeck`);
   `go build ./...` ok; new unit tests for role/additive-vs-breaking + gate exit codes.

## Idempotence & recovery
Protocol text edits are pure git-reversible. `parley preflight` is read-mostly; its only
write (additive consumer sync) is atomic + recorded + zone-preserving + git-reversible,
and never fires in a source repo. The roster ping spawns bounded child process groups
killed on timeout (no leaks). Re-running preflight is safe (idempotent report; sync is a
no-op once hashes match).

## Known risks / de-risking
Source-repo regression (→ `protocolRole` fail-closed); unattended deadlock (→ gates
excluded from auto-answerer, hard-stop); roster-ID/runtime-ID drift (→ explicit map);
false availability (→ sentinel PONG + supervision net); ping cost (→ concurrency +
per-process cache). De-risk by implementing + unit-testing the role/additive classifier
and the gate exit-code paths first, before the run-wiring.

## References
- Consensus: ./consensus.md (signoffs claude-1/codex-1/hermes-1 ✅; agy waived)
- Round-01: ./round-01/{claude-1,codex-1,hermes-1,antigravity-1}.md
- Design brief: ./reference/design-brief.md
