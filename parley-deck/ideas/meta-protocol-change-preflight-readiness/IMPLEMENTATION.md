---
idea: meta-protocol-change-preflight-readiness
status: complete
implementer: claude-1
started: 2026-06-19
completed: 2026-06-19
branch: parley-deck-cli#meta/preflight-readiness
head-commit: c9e872d (impl) → fix-up cycle 1 (see below)
design-pr: n/a (single-repo; design + impl in one PR)
implementation-pr: <pending>
---

## Summary of work
Implemented the FINAL spec: protocol text (§9.0 pre-idea readiness check + §5 quorum
sentence + §7 carve-out, both COOPERATION.md copies, drift-guard lockstep), the
`meta/version.json` `protocolRole: "source"` field, and a new `parley preflight` Go
command (freshness check + hosted-PONG roster ping + gates + exit codes) wired into
`parley run` before idea creation. Build + full `go test ./...` + drift guard green.

## Implementation plan / checklist
- [x] Files/areas: COOPERATION.md ×2 (§9.0/§5/§7), meta/version.json (protocolRole),
  internal/app/preflight.go (+ _test.go), internal/app/app.go (dispatch + run wiring +
  usage), internal/fsutil/fsutil.go (WriteFileAtomic).
- [x] Checks: `go build ./...`, `go test ./...` (incl. drift guard), real `parley
  preflight` run.
- [x] Risk notes: source-repo fail-closed; unattended hard-stop; gates not
  auto-answered; run-path honors the hosted-PONG operator ruling.

## Deviations from FINAL.md
1. **Run-path ping default — corrected to honor the operator ruling.** The
   implementation agent initially wired the `parley run` pre-check as presence-only
   (Tier-0) to avoid ~90s/run latency + run-test shell-outs. The implementer (claude-1)
   changed it to **hosted-PONG by default** (the operator ruling "hosted PONG every
   idea"), with `--no-ping` (presence-only) and `--no-preflight` (skip) as opt-outs; the
   two run-path unit tests pass `--no-ping` so they stay hermetic (no live shell-out).
2. **Roster IDs in reports = runtime IDs (`codex`/`claude`/`agy`/`hermes`), not the
   §2 `-1` IDs.** Go discovery (`agents.Discover` via `config.LoadAgentSpecs`) keys on
   the agents.toml/built-in spec IDs (bare names); the `claude-1`/`codex-1`/`hermes-1`/
   `antigravity-1` names exist only in COOPERATION.md §2 and the gitignored
   `headless-agents.local.json`, which the Go runtime does not read. So the consensus's
   "reports use roster IDs + roster-ID↔runtime-ID map" is **only partially realized** —
   the tool uses the IDs the runtime actually knows. **Phase-6 finding / follow-up**
   (resolve the §2-vs-runtime naming mismatch separately; not a blocker — the tool is
   functional and correct with runtime IDs).

## Notes for reviewers
- Verify: source-repo advisory (no COOPERATION.md write); exit codes 0/3/1/2; gates not
  auto-answered (no RiskLow+DefaultAnswer); unattended hard-stop never reads stdin;
  `mergePreservingZones` keeps header→§2 and replaces §3-onward; the run-wiring honors
  the hosted-PONG ruling without breaking the run/TUI path.
- The two deviations above (run-path corrected to hosted-PONG = honoring the ruling;
  roster-ID = runtime IDs = Phase-6 follow-up).

## Progress
- [x] (2026-06-19) Protocol §9.0/§5/§7 both copies; drift guard green.
- [x] (2026-06-19) meta/version.json protocolRole=source.
- [x] (2026-06-19) parley preflight implemented (delegated build + implementer review).
- [x] (2026-06-19) Run-path corrected to hosted-PONG default (operator ruling); tests
  updated to `--no-ping`; full `go test ./...` + build green.
- [ ] Phase 6 review (codex-1, hermes-1; agy likely waived) → Phase 7 → release 1.30.0.

## Decision Log
- Decision: run-path ping defaults to hosted PONG (not presence-only). Rationale:
  operator ruled "hosted PONG every idea"; `--no-ping`/`--no-preflight` give the speed
  escape; run-tests pass `--no-ping` to stay hermetic. Date/Author: 2026-06-19 / claude-1.
- Decision: ship reports with runtime IDs; defer §2-vs-runtime `-1` naming
  reconciliation to a follow-up. Rationale: the `-1` names aren't in any source the Go
  runtime reads; mapping is out of this idea's scope. Date/Author: 2026-06-19 / claude-1.

## Surprises & Discoveries
- Feature 2 was validated LIVE during this very idea's own deliberation: agy PONGed
  green at the readiness ping but hung on the heavier consensus-signoff append → the
  operator confirmed a per-idea waive (the exact gate being built). Evidence: consensus
  signoffs (claude-1/codex-1/hermes-1 ✅, agy waived).
- The `-1` roster IDs (from the prior roster-update task) live only in §2 + a gitignored
  JSON the Go discovery never reads — surfaced the §2-vs-runtime naming gap (above).

## Validation evidence
- `go build ./...` OK; `go test ./...` all packages ok incl.
  `internal/protocol/TestEmbeddedDefaultMatchesLiveDeck`.
- `parley preflight --no-ping` in this source repo: freshness=source/advisory (no
  write), roster presence table, exit 0. `--ping-timeout 100s` (real hosted PONG): all
  4 installed agents PONGed green, exit 0.

## Fix-up cycle 1
status: complete
completed: 2026-06-19

### Fixes applied (all 6 agreed; verified by claude-1 + re-signed by codex-1/hermes-1)
- [CRITICAL] §1 hard-stop now evaluated on the selected `--participants` set written to
  `00-prompt.md` (not the full discovered set) — `runTaskPreflight` rewritten +
  `participantDiscoveries`; `TestRunParticipantsSubsetHardStopsSolo`.
- [MAJOR] `--yes` confirms gates + records `excluded:` into `00-prompt.md`
  (`CreateOptions.Excluded` → `CreateIdeaWithExclusions`) + backfills `protocolRole`.
- [MAJOR] `--json` returns the real exit code (no longer masked to 0).
- [MAJOR] `parley init` writes `protocolRole:consumer`; non-workspace → exit 1;
  absent metadata in an existing deck → role/backfill gate.
- [MAJOR] hosted-PONG exact sentinel (`isExactPONG`) + echoed-prompt/commentary tests.
- [MAJOR] `classifyBump` fail-closed to `bumpMajor` on parse failure; dead `bumpUnknown` removed.

### Deviations from agreed fixes
None. The `excluded:` recording used full `CreateOptions` threading (the cleaner option).

## Outcomes & Retrospective
Shipped §9.0 readiness check (protocol, both copies) + `parley preflight` (freshness +
hosted-PONG ping + gates) wired into `parley run`, honoring the operator's hosted-PONG
ruling. Phase 6 review caught a real CRITICAL (§1 bypass via `--participants`) + 5 MAJORs
— all fixed in one cycle and re-signed. `go build` + `go test ./...` + drift guard green.
Two MINORs deferred (roster-ID `-1` reconciliation; preflight perf). Lessons: (1) the §1
non-solo guard must always evaluate the *exact* set written to `00-prompt.md`, never a
superset; (2) Feature 2 was validated live — agy PONGed at the ping but hung on the
heavier signoff append → operator-confirmed per-idea waive, the exact gate this idea
builds.
