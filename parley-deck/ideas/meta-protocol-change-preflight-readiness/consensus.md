---
idea: meta-protocol-change-preflight-readiness
phase: 3-consensus (facilitator draft)
drafted-by: claude-1
date: 2026-06-19
participants: [claude-1, codex-1, hermes-1, antigravity-1]
status: consensus-signoff
---

# Consensus — Pre-idea readiness check

All four participants wrote independent round-01 (full 4-quorum; all PONGed at ping).
Strong convergence on architecture. The 3 operator decisions stay locked; round-01
**refined the mechanics** and surfaced two improvements worth the operator's eye
(marked ★). Grounded heavily in existing code by hermes-1.

## Agreed design

### Freshness (decision 1: auto additive / confirm breaking)
- **Reuse `parleyDeckSkillStatus`** (`internal/app/version_status.go`, already 10s-bounded);
  drift = string compare `meta/version.json.protocolSha256` vs the status payload's
  packaged hash. Computed **once per process** (the auto-driver is one process) — no
  file TTL cache in v1; `--no-preflight` escape for CI.
- **Explicit `meta/version.json.protocolRole: "source" | "consumer"`** decides
  direction (not a semver heuristic). **Source → advisory only, writes nothing** → the
  source-vs-consumer inversion is structurally impossible (this repo is `source`).
- ★ **Missing/unknown role → do NOT auto-write; one-time confirm + backfill the field**
  (codex-1's fail-safe, chosen over a silent default in either direction). After the
  one confirm, the field is set and never prompts again.
- **Additive vs breaking = semver of `deckVersion`** for v1 (major = breaking → confirm;
  minor/patch = additive → auto). Safety net = zone-preservation + the sync record +
  the §7 carve-out (human-auditable). A structural section-diff is a **v2** hardening,
  not needed for v1.
- **Zone preservation** via the existing drift-guard allowlist (header, §0 transport,
  §2 roster). Consumer auto-sync writes the single project `COOPERATION.md`; the
  two-copy lockstep is a *source-repo* concern (and source = no auto-write).
- Record each sync to `meta/protocol-sync_<ISO-timestamp>.md` (timestamp, not bare
  date, to avoid same-day collisions) with old/new hash, version, type, preserved
  zones, diff summary; update the `Protocol synced:` header line.

### Roster ping (decision 2: per-idea temporary, both gates user-confirmed)
- **Default ping = hosted PONG round-trip per agent** (★ operator ruling 2026-06-19,
  chosen over hermes-1's cheaper Tier-0 proposal — strongest liveness signal). Each
  rostered participant gets a bounded sentinel PONG via its **real configured
  invocation** (not an ad-hoc call). `command -v` is a pre-check (missing CLI →
  `unavailable:missing`, no PONG attempted). Bound each probe **~90s** (hermes inits
  ~40s, per codex-1); **agy MUST get `--add-dir "$(pwd)"`**; run probes **concurrently
  behind one global deadline**; kill the process group on timeout. Available = exits in
  time with the exact sentinel content; timeout/empty/exit-error = unavailable (reason
  recorded). Accepted cost: ~4× latency+tokens per idea start; mitigated by concurrency
  + a per-process cache (no re-PONG within the same idea).
- **Mid-idea safety net unchanged:** a PONG can pass while a heavy round later hangs
  (e.g. agy). The supervised round-01 watchdog (`internal/runner/supervision.go`
  first-event 120s / stall 30m) + `failclass.go` still classify a mid-idea hosted hang;
  it downgrades to the same per-idea user-confirmed waive.
- Probe via the same launch mechanics as a real round (prompt mode, root, configured
  args), sentinel written to a **non-canonical probe dir** (never an idea round path).
- **Roster-ID ↔ runtime-ID map** (codex-1's catch): roster uses `claude-1`/`codex-1`/
  `hermes-1`/`antigravity-1`; runtimes are `claude`/`codex`/`hermes`/`agy`. Add a
  machine-readable map (in `meta/headless-agents.local.json` or a roster map); reports +
  `00-prompt.md` use **roster IDs**; a roster ID with no runtime → `not_configured`
  (never silently skipped or name-matched). Distinguish `--participants` selection
  (`not_selected`) from availability exclusion.
- **Two gates, user-confirmed, per-idea, recorded in `00-prompt.md`:**
  `excluded: [<roster-id> — reason — confirmed_by=user confirmed_at=<date>]`;
  `reincluded: [...]`. Re-include detected by scanning recent ideas' `excluded:` lines
  for an id now Tier-0 green with no later `reincluded` (closes the chain → no
  forever-prompt). Quorum is **locked at Phase 0** (no mid-idea change; a mid-idea hang
  falls to §5 async rules).

### `parley preflight` CLI (decision 3)
- One command: `parley preflight [--dir DIR] [--json] [--yes] [--ping-timeout D]
  [--no-ping]`; default = freshness check + **hosted-PONG roster ping** + readiness
  report, **no irreversible action** except the locked additive consumer auto-sync.
  `--no-ping` falls back to a Tier-0 presence-only check for a quick look.
- **Exit codes reuse existing semantics** (`app.go`): `0` ready · `3` pending gate
  (breaking-freshness / exclude / re-include — names the gate + prints the confirm
  command; reuses the existing exit-3 "pending manual handoff") · `1` hard failure (no
  workspace, or excluding would leave <2 participants → the §1 non-solo hard-stop) ·
  `2` usage.
- **ONE shared `preflight(root, opts) (report, gates)` function** called by both the
  standalone command and `parley run`.
- **Wiring (the deadlock-safety invariant):** `parley run` calls `preflight` **BEFORE
  `runcontrol.Create`** (reusing the `discovered` slice already computed), so a gate
  returns before any half-open idea/round exists.
  - Attended → route gates through the existing HITL (`parley answer`) and stop for
    confirmation (mirrors `confirmLaunch`/`--yes`).
  - Unattended (`--auto`, no TTY) → **hard-stop non-zero; never read stdin, never
    auto-exclude, never auto-adopt a breaking bump.** The two new gates are **excluded
    from `StartAutoAnswerer`'s auto-answer set** — this one rule is the whole
    deadlock-safety argument and keeps "user-confirmed" intact in every mode.
- `resume`/`continue` → Tier-0 roster re-check only (no freshness, no gates).

### Protocol text (both COOPERATION.md copies, drift-guard lockstep)
- New **§9.0 "Pre-idea readiness check"**: freshness (auto-additive/confirm-breaking,
  source=advisory) + **hosted-PONG roster ping (operator ruling)** + both gates —
  wired as step 0 of §9.
- **§5** sentence: quorum is set at Phase 0 from the readiness check; excluded agents
  don't count for this idea; quorum locks once Phase 0 completes.
- **§7 carve-out**: applying an upstream-ratified, additive, zone-preserving protocol
  sync is NOT a protocol-change idea.

## REJECT / out of scope
- A structural section-diff engine (v2), a `meta/exclusions.json` index (v2), a
  cross-process TTL cache (v2). (hermes-1's Tier-0-default ping was considered and
  overridden by the operator ruling — hosted PONG is the default; Tier-0 survives as
  the `command -v` pre-check + `--no-ping`.)
- Silent role default in either direction; silent roster change; auto-answering the new
  gates; auto-adopting a breaking bump.

## Operator rulings (2026-06-19)
1. **Roster ping = hosted PONG every idea** — operator chose the strongest liveness
   signal over hermes-1's cheaper Tier-0 proposal. Cost accepted; mitigated by
   concurrency + per-process cache + process-group kill on timeout. (hermes-1's Tier-0
   reuse-of-`agents.Discover` survives as the `command -v` pre-check and `--no-ping`.)
2. **Missing `protocolRole` → one-time confirm + backfill** (not a silent default), so
   we never guess the dangerous direction.

## Signoffs
<!-- Each participant APPENDS their own block. -->

### Signoff: claude-1 — 2026-06-19
Status: ✅ ACCEPT
Notes: Facilitator/architecture. Consensus reflects the operator ruling (hosted PONG
every idea) and the fail-closed `protocolRole` handling. My two correctness anchors are
in: (1) the deadlock-safety invariant — preflight before `runcontrol.Create`, the two
new gates excluded from `StartAutoAnswerer`, unattended = hard-stop never block; (2) the
roster-ID↔runtime-ID map so quorum/records use roster IDs and a hosted PONG never
mislabels an identity. Source-repo stays advisory (no auto-write). Ship.

### Signoff: codex-1 — 2026-06-19
Status: ❌ BLOCK
Notes: I accept the operator-ruling direction, but the draft is internally inconsistent: it still says §9.0 should use "Tier-0 default ping" and rejects hosted PONG by default, contradicting the agreed design and operator rulings. This would let FINAL/implementation choose the wrong default.
Counter-proposal: Amend only those stale consensus lines so hosted PONG every idea is the default; keep Tier-0 as the `command -v` pre-check, `--no-ping` fallback, and resume/continue re-check.

### Signoff: codex-1 — 2026-06-19 (re-sign after fix)
Status: ✅ ACCEPT
Notes: The stale live consensus lines are fixed: hosted PONG every idea is now the stated default and the text no longer rejects hosted PONG by default. Remaining Tier-0 wording is non-default context, with the old failure quoted only inside my previous historical BLOCK.

### Signoff: hermes-1 — 2026-06-19
Status: ✅ ACCEPT
Notes: Consensus is internally consistent and faithfully records the operator ruling (hosted PONG every idea as default) plus my Tier-0 design surviving as the `command -v` pre-check, `--no-ping` fallback, and resume/continue re-check. Reservation (recorded, not blocking): I still hold that Tier-0-as-default would be cheaper and sufficient given the existing mid-idea watchdog already catches hosted hangs — but I accept the operator's call for the strongest liveness signal, and my fallback paths are preserved intact.

### Facilitator note: antigravity-1 signoff waived — 2026-06-19 (claude-1)
antigravity-1 (agy) **PONGed green at the readiness ping** (short bounded probe) but
then **hung on the heavier consensus-signoff append** — the exact mid-idea-hang case
this very consensus documents (§ "Mid-idea safety net"). The **operator confirmed a
per-idea waive** of agy's signoff (the Feature-2 exclude gate, dogfooded live). This is
a documented per-participant tooling waive, **not** a solo run: signoff quorum =
claude-1 ✅ + codex-1 ✅ + hermes-1 ✅ (three participants), and agy's own
`round-01/antigravity-1.md` (docs/UX lens) already endorsed the design and is canonical
evidence. Consensus is reached; proceeding to Phase 4.
