---
agent: hermes-1
idea: meta-protocol-change-preflight-readiness
review-round: 1
date: 2026-06-19
reviewed-commit: c9e872d
---

## Summary

The readiness machinery is sound where it matters most: the hosted-PONG probe is
bounded, concurrent, and process-group-killed-and-reaped on timeout (no leaked
children); freshness is fail-closed for source/unknown/missing; gates are never
written as HITL questions so the auto-answerer cannot touch them; and the
unattended path never reads stdin and stops before `runcontrol.Create` (no
half-open idea). Two correctness gaps block a clean bill: the `--yes` confirm
flow is a dead flag (the documented confirm command cannot clear a gate and the
exclusion is never recorded), and `classifyBump` fails *open* on unparseable
versions, which can auto-apply a breaking consumer bump.

## Findings

[MAJOR] `--yes` is a no-op; the confirm loop and exclusion recording are not
implemented (undelivered acceptance criterion #2, undeclared deviation).
`preflightOptions.Yes` is populated from the `--yes` flag at
`internal/app/preflight.go:128` but is **never read** anywhere in `preflight`,
`classifyAndSyncFreshness`, `checkRoster`, or `runTaskPreflight` (grep for
`opts\.Yes`/`\.Yes\b` in `internal/app` finds only the assignment at :128 and
the run-path opts at :186 which omits `Yes` entirely). Consequences:

- Every gate's `Confirm` string is `parley preflight --dir <root> --yes`
  (preflight.go:247, :302, :320). Re-running that command sets `Yes=true` and
  then re-enters `preflight`, which **still emits the same gate and returns exit
  3**. The documented confirm command cannot clear a gate — it is a dead end.
- `excluded:` is recorded nowhere (`excluded:` does not appear in `internal/`
  at all). FINAL.md §9.0 and acceptance criterion #2 require the exclusion to be
  recorded in `00-prompt.md` (`excluded: [<roster-id> — reason — confirmed
  <date>]`); no code writes it. There is also a chicken-and-egg here the
  implementation does not resolve: preflight runs *before* idea creation, so
  `00-prompt.md` does not exist yet — the exclusion would have to be plumbed
  forward into the idea-creation step, and it is not.
- In the run path, `runTaskPreflight` builds `preflightOptions{Root, NoPing}`
  with no `Yes` (:186), so `parley run --yes` does not confirm preflight gates
  either. The only working escape from a gate is `--no-preflight`, which is
  documented as a "CI escape" and bypasses the safety check entirely (including
  real availability gates) rather than confirming them.

This fails safe (the gate still stops the run), so not CRITICAL, but a
documented user-facing flow is non-functional and a stated acceptance criterion
is unmet. It is also not listed among the two declared deviations in
IMPLEMENTATION.md. Fix: either make `--yes` suppress the exclude-agent gate
(treat as confirmed-excluded), return exit 0 when ≥2 remain, and plumb the
exclusion list into `00-prompt.md` at idea creation; or remove the `--yes` flag
and the `Confirm` strings and document `--no-preflight`/manual edit as the
escape. As shipped the flag misleads.

[MAJOR] `classifyBump` fails open on unparseable versions — can auto-sync a
breaking consumer bump. `classifyBump` (preflight.go:347) returns `bumpPatchMinor`
(additive → auto-write) when *either* version fails to parse. The comment
(:345) calls this "conservative," but for a freshness sync the conservative
choice is to **gate** when the major-vs-minor classification is uncertain, not
to auto-write. A genuinely breaking (major) packaged bump whose `deckVersion`
is unparseable would be silently auto-synced via `syncConsumerProtocol`,
replacing the §3-onward body without the user confirmation the spec mandates
for breaking changes — the exact failure mode the breaking gate exists to
prevent. Mitigations: consumer-only (this source repo never reaches the path);
the write is atomic, zone-preserving, recorded, and git-reversible. But it
contradicts the spec's fail-closed intent and the comment mislabels the
direction. Recommend: return `bumpMajor` (gate) on parse failure, or at least
gate when the packaged version is unparseable. `parseMajor` also strips only a
leading `v` and splits on the first `.` — pre-release suffixes
(`1.4.0-rc1` → head `1`) happen to work, but `1` alone parses to major 1, so a
bare `1` vs `2` is classified correctly; the real hole is the unparseable
fallback direction, not the parser.

[MINOR] Source repo pays the `parley-deck-skill status` shell-out (≤10s, plus a
≤5s legacy fallback) even on `--no-ping`, though the classification is
predetermined. `classifyAndSyncFreshness` always calls `parleyDeckSkillStatus`
(:272) before reading `version.json`. In a source repo the result is always
`source-advisory` / no-write regardless of what the skill reports, so the probe
only contributes `SkillVersion` to the report. For the `--no-ping` "fast/CI"
path in a source repo this is avoidable latency. Consider reading `version.json`
first and skipping the skill-status probe when `role == source` (make
`SkillVersion` optional in the report), or gating the probe on `!NoPing`.

[MINOR] Per-probe timeout == global deadline, so the "global deadline" adds no
tighter bound. `checkRoster` (preflight.go:527) wraps the parent ctx with
`context.WithTimeout(ctx, timeout)` where `timeout == PingTimeout` (90s), and
each `hostedPONG` independently creates `context.WithTimeout(ctx, timeout)`
with the same 90s (:579). All probes start in the same loop, so the per-probe
and global deadlines coincide; the global deadline never fires before a
per-probe one. The comment at :534 ("Global deadline bounds all concurrent
probes together") oversells it. The global ctx does still wire parent
cancellation, so it is not useless — but if the intent is a wall-clock bound
tighter than per-probe, the two values should differ. Harmless as-is.

[MINOR] Roster report uses runtime IDs while the §2 roster (edited in this same
PR) uses `-1` IDs — an active mismatch. The preflight table prints
`agent.ID` (`codex`/`claude`/`agy`/`hermes`) via `rosterEntry.RosterID`
(preflight.go:546), but the §2 roster table in `parley-deck/COOPERATION.md`
(this PR's diff) now lists `claude-1`/`codex-1`/`hermes-1`/`antigravity-1`.
Cross-referencing the report against §2 will mismatch every row. This is the
declared deviation #2 / Phase-6 follow-up; flagging here so it is not lost —
the PR that introduces the §2 `-1` names is the same PR whose tool output cannot
use them.

[NIT] `bumpUnknown` (preflight.go:338) is declared but never returned by
`classifyBump` (only `bumpPatchMinor` / `bumpMajor`). Dead constant; remove or
wire it in as the parse-failure result if the [MAJOR] above is fixed by gating.

## What is correct (for the record)

- Probe lifecycle (lens Q1): `SetNewProcessGroup` uses `Setsid` so `pgid == pid`
  (`procctl_unix.go:14`); `Capture` records PID/PGID; on `probeCtx.Done()`,
  `KillGroup` signals `-pgid` (SIGTERM → 1.5s grace → SIGKILL) with explicit
  self-protection against killing parley's own group
  (`procctl_unix.go:36-38`), then `<-waitErr` reaps. ESRCH on an already-dead
  group is handled. No leaked children. Concurrency is one goroutine per rostered
  agent behind a `sync.WaitGroup` (roster is small). Bounded.
- `--no-ping` fast path (Q2): presence-only, no goroutine, no shell-out;
  `TestPreflightNoPingPresenceOnly` asserts the probe is never called. The run
  path reuses the already-discovered slice (no re-discovery).
- Freshness fail-closed (Q3): `source` returns immediately with no write
  (test `TestClassifyFreshnessSourceIsAdvisoryNoWrite` asserts `COOPERATION.md`
  is byte-identical); `unknown-role` → gate, no write; missing `version.json` →
  `no-version-json`, no write. The only write is consumer-additive, atomic via
  `fsutil.WriteFileAtomic` (temp + rename), and unreachable in this source repo.
- Gates never auto-answered (Q4): preflight gates live only in the in-memory
  `preflightReport`; they are never written as HITL questions, and
  `runTaskPreflight` stops before `runcontrol.Create`, so `StartAutoAnswerer`
  (which only auto-answers `RiskLow`+`DefaultAnswer`+open per
  `internal/hitl/hitl.go:151`) never starts on a gate.
- Unattended hard-stop never blocks (Q4): `attendedRun` → `isTerminal(os.Stdin)`;
  the unattended branch writes to stderr and returns `(code, true)`; no function
  in the preflight path reads stdin.
- Run-path hosted-PONG default (Q6): acceptable given the operator ruling and
  the `--no-ping` (presence-only) / `--no-preflight` (skip) escapes; the two
  run-tests pass `--no-ping` to stay hermetic. The ~90s/run cost is the
  documented tradeoff.
- Minimal/clean (Q5): 718 lines is reasonable for the scope; the zone-merge is a
  pure function with focused unit tests; the multi-key fallbacks in
  `skillStatusPackagedDeckVersion`/`skillStatusPackagedProtocol` are defensive
  against an external payload shape and acceptable.
