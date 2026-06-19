---
agent: claude-1
idea: meta-protocol-change-preflight-readiness
round: 1
lens: architecture & correctness (lead synthesis)
date: 2026-06-19
---

## Summary
The 3 decisions are locked; the architecture risk is concentrated in **two correctness
hazards**: (1) the freshness sync must never regress the *source* repo or clobber
project zones, and (2) the readiness gates must never deadlock an unattended
`parley run`. I propose grounding both in mechanisms that already exist
(`meta/version.json`, the drift-guard allowlist, §1/§5) and making `parley preflight`
a pure *reporter+gater* that performs at most one reversible, confirmed action.

## Proposed mechanics

### Freshness (decision 1)
- **Version line:** track the *protocol* via `deckVersion` + `packagedProtocolSha256`
  in `meta/version.json` — NOT the `parley-deck-cli` version. They are different lines
  (skill `1.3.x` vs CLI `1.29.x`); conflating them is a bug.
- **Direction is decided by role + version order, drift only by hash:**
  `protocolSha256 != packagedProtocolSha256` means "differs", not "behind". Adopt only
  when `protocolRole == "consumer"` AND `semver(packagedDeckVersion) > semver(projectDeckVersion)`.
  Add `protocolRole: "source" | "consumer"` to `meta/version.json`; `parley init` writes
  `"consumer"`; the source repo carries `"source"`. **Fail safe:** if `protocolRole`
  is missing/unknown → treat as `source` (advise-only), never auto-write. This makes
  the dangerous direction the one that requires an explicit opt-in.
- **additive vs breaking:** two signals, both must say "additive" to auto-apply —
  (a) semver: minor/patch = additive, major = breaking; (b) structural: the new body is
  a **superset** (only adds sections/lines; modifies/removes nothing in shared zones).
  Any modify/remove, or a major bump → confirm. This catches a minor-versioned but
  rule-changing edit.
- **Zone preservation:** reuse the drift-guard allowlist (header, §2 roster, transport)
  verbatim — the sync rewrites only the shared body, transplanting project zones
  unchanged. Reversible (git); record `Protocol synced:` + `meta/protocol-sync_<date>.md`.

### Roster ping (decision 2)
- **Probe = the PONG we just ran:** `command -v <cli>` (missing → `unavailable:missing`)
  then a bounded trivial PONG (responds in ≤ timeout → `available`; timeout/empty →
  `unavailable:unresponsive`). agy MUST get `--add-dir "$(pwd)"`. Default ping timeout
  ~60-90s (separate, much shorter than round/review timeouts) so a hang is caught fast.
- **Per-idea exclusion (temporary):** record in the idea's `00-prompt.md` readiness
  block (`excluded: [<id> — reason — confirmed <date>]`); the agent stays in §2 roster
  and is re-pinged next idea. **Both gates user-confirmed:** exclude needs confirm;
  re-include (a previously-excluded agent now answering PONG, not in current quorum)
  needs confirm. Excluding the last non-facilitator still hits the §1 non-solo stop.
- **agy nuance (from today):** a short PONG can pass while a heavy round/signoff still
  hangs. So "available" from the ping is necessary, not sufficient — keep the existing
  watchdog/first-output timeout on the actual round, and a mid-idea hang downgrades to
  the same user-confirmed per-idea waive. Document this explicitly.

### `parley preflight` (decision 3)
- `parley preflight [--idea <slug>] [--json] [--yes]`: runs the skill `status` +
  roster ping, prints a readiness report, and returns a typed result:
  `ready` | `needs-confirmation(freshness|roster)` | `blocked(non-solo)`.
- **Side-effect discipline:** it performs at most ONE action — the confirmed freshness
  sync — and only with `--yes` or interactive confirm; the roster gates only *report*
  (the facilitator/user acts). No silent writes.
- **No deadlock:** wire into `parley run` as a pre-idea step that is **advisory by
  default**; under unattended auto-drive a `needs-confirmation`/`blocked` result
  **stops with a clear status** (not a hang, not a silent proceed) — mirrors how
  auto_implement already gates. `--yes` pre-authorizes additive freshness only, never a
  roster change.

## Concerns / open questions
- `protocolRole` bootstrap for existing consumer projects that lack the field (the
  fail-safe `=source` is conservative but means existing consumers won't auto-adopt
  until re-synced once — acceptable; flag in release notes).
- Ping cost: 4 live round-trips per idea adds latency + tokens. Mitigate with a short
  TTL cache (skip re-ping within N minutes in the same session) — but never cache an
  *exclusion* decision.
- Keep `parley preflight` read-mostly so it is safe to run anytime / in CI.

## Risks
- **Source-repo regression** if role detection is wrong → fail-safe to source/advise-only.
- **Unattended deadlock** → typed stop, never block-and-wait in auto mode.
- **Zone clobber** → reuse drift-guard allowlist; superset-only auto-apply.
- Scope creep into a full "update manager" — keep it to detect + (confirmed) sync +
  ping; nothing more.
