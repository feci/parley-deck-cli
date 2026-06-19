---
idea: meta-protocol-change-preflight-readiness
kind: reference
drafted-by: claude-1 (facilitator)
date: 2026-06-19
---

# Design brief — Pre-idea readiness check (locked decisions)

Shared evidence for the panel. The **goals + 3 design decisions are already locked by
the operator**; round-01 should pressure-test the design, surface risks/edge cases,
and refine the *mechanics* — not relitigate the decisions.

## Goal (operator, verbatim intent)
At the start of *every* idea, the facilitator should automatically (a) check whether a
newer protocol version is installed and bring the session onto it, so the operator
never has to reason about whether the session has the updated skill; and (b) ping all
rostered agents and, if one is unavailable, adjust the roster — but **excluding an
agent OR re-including a previously-excluded agent both require user confirmation**.

## What already exists (build on it, don't reinvent)
- `parley-deck-skill` CLI is installed; `parley-deck-skill status --target all
  --project . --json` reports installer version, runtime skill versions, project
  metadata, compatibility warnings. `sync-project` updates only `meta/version.json`
  today (it does NOT overwrite `COOPERATION.md`).
- `parley-deck/meta/version.json` already carries `protocolSha256` (live deck),
  `packagedProtocolSha256` (protocol bundled in the installed skill), `deckVersion`,
  `source`. → drift is one hash comparison.
- §9 "Session-start checklist" exists but only says *read* COOPERATION.md; it does NOT
  mandate a freshness sync or a roster ping.
- §1 "Non-solo execution requirement" + §5 "Quorum" already forbid silently collapsing
  to solo and require a recorded, user-authorized solo exception.
- Two-copy protocol + drift guard (allowlisted project-specific zones: header, §2
  roster, transport) — see [[embedded-default-drift-guard]].

## LOCKED decision 1 — protocol auto-freshness: **auto additive, confirm breaking**
At idea start the facilitator runs a freshness check (`parley-deck-skill status`, else
hash-compare `protocolSha256` vs `packagedProtocolSha256`). If the project is a
**consumer** and the installed skill's protocol is newer:
- **Additive / compatible** bump → **auto-sync** the protocol *body* into
  `COOPERATION.md`, **preserving project-specific zones** (roster §2, transport §0,
  header) via the same allowlist discipline as the drift guard; update the
  `Protocol synced:` line + drop a `meta/protocol-sync_<date>.md` record; report the
  diff summary; continue on the new version.
- **Breaking / semantic** bump (major version, or a changed/removed rule) → detect +
  **pause for user confirmation** before writing.
- Never a blind copy.

⚠️ **CRITICAL PITFALL (live, in THIS repo): source vs consumer inversion.** Here
`protocolSha256 ≠ packagedProtocolSha256` because this deck is the protocol **SOURCE**
and is *ahead* of the published skill (§13, Fusion/ExecPlans aren't in npm `1.3.1`). A
naive "hashes differ → adopt the packaged one" would **regress** the canonical
protocol. The rule MUST detect role — proposed: `meta/version.json.protocolRole:
"source" | "consumer"` (or derive "which side is newer" + a source flag). In a *source*
repo, auto-adopt is **off** (advise only).

§7 interaction: a freshness sync is *adopting upstream-ratified text*, NOT inventing a
change → carve out an explicit exception so it does not require a meta-protocol-change
idea.

## LOCKED decision 2 — roster ping + **per-idea temporary** exclusion, both gates user-confirmed
At idea start, for each rostered participant: probe availability (`command -v` + a
trivial PONG with a short timeout). Build an available/unavailable table.
- **Exclude** an unavailable agent from THIS idea's quorum → STOP + **user
  confirmation**; record in `00-prompt.md` (`excluded: [<id> — reason — confirmed
  <date>]`). **Per-idea / temporary**: the agent stays in the §2 roster and is
  re-pinged at the next idea (no roster edit).
- **Re-include** a previously-excluded (now-available) agent into quorum → **also user
  confirmation** (no silent quorum expansion).
- Excluding the last non-facilitator still hits the §1 non-solo block.
- Hung agent (e.g. agy past its ping timeout) = **unavailable** → user-confirmed
  exclusion. This formalizes the ad-hoc agy waive done on
  [[fusion-execplans-shipped]] and the agy hang in [[agy-headless-regression]].

## LOCKED decision 3 — scope: protocol text **+** `parley preflight` CLI
- **Protocol:** new §9.0 "Pre-idea readiness check" (freshness + roster ping + both
  gates) + a §7 carve-out for upstream sync + a §5 sentence on per-idea quorum from the
  ping. Both `COOPERATION.md` copies, drift-guard lockstep.
- **Tooling:** a `parley preflight` command that runs `status` + roster ping, prints
  the readiness report, and surfaces the two gates (it performs no irreversible action
  without confirmation). Wireable into `parley run` / auto-drive as a pre-idea step.

## Questions for the panel (refine, don't relitigate)
1. Best mechanism + exact source/consumer detection (the `protocolRole` field vs
   deriving from version ordering). How to be safe in the source repo.
2. Exactly what "additive/compatible" vs "breaking" means operationally (semver of
   `deckVersion`? a structural protocol-section diff? both?).
3. The roster-ping probe per agent (what's a reliable, cheap PONG for claude/codex/
   hermes/agy; timeout; how to classify "hung" vs "slow").
4. Where the per-idea exclusion is recorded and how re-include is detected next idea.
5. `parley preflight` surface: flags, output shape, exit codes, and how/whether it
   gates `parley run` (auto-drive) without breaking unattended runs.
6. Anything that conflicts with existing invariants (§1/§5/§7/§9, drift guard) or that
   we should explicitly NOT do.
