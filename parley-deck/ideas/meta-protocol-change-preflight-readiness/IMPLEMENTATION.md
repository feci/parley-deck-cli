---
idea: meta-protocol-change-preflight-readiness
status: implemented
implementer: claude-1
started: 2026-06-19
completed: 2026-06-19
branch: parley-deck-cli#meta/preflight-readiness
head-commit: <pending first commit>
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

## Outcomes & Retrospective
(To complete at Phase 8.) Pending Phase 6-7 review + release.
