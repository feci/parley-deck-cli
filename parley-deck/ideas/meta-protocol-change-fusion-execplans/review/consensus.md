---
idea: meta-protocol-change-fusion-execplans
review-cycle: 1
drafted-by: claude
date: 2026-06-19
reviewed-commit: 8477ea6
---

## Agreed fixes
- **[MINOR] Decision Log deviations cross-reference** (from hermes/review/round-01
  [MINOR]). Edit 2's `## Decision Log` description dropped the spec parenthetical that
  distinguishes it from the `## Deviations from FINAL.md` head above it. Fix: append
  "Deviations still go under `## Deviations from FINAL.md` above." to the Decision Log
  description, byte-identically in both COOPERATION.md copies.
- **[NIT→fix] Harmonize the general rigor trigger** (from hermes/review/round-01
  [NIT]). The ratified preamble lists four triggers (complex / `auto_implement` /
  driver-managed / **pipeline**); the Phase 4/5 prose listed only three. Fix: add
  "pipeline" to the general self-containedness trigger in Phase 4 and the two Phase 5
  spots, so the general trigger matches the preamble everywhere. (The idempotence
  clause keeps its intentionally tighter `auto_implement` / action / pipeline trigger.)
- **[NIT→fix] Split the dense Phase 4 paragraph** (from hermes/review/round-01 [NIT]).
  Move the `Idempotence & recovery` clause into its own paragraph for scannability.

## Deferred follow-ups
- None. (The `parley retro` confident-error *mining* tooling remains the pre-declared
  deferred follow-up from `FINAL.md`, governed by §13, out of scope for this idea.)

## Dismissed findings
- codex's verification note that `internal/runner/TestDurableKillEndToEndRealProcess`
  failed *in codex's sandbox* (`no recorded boot id; not killed`) — environment
  artifact (the durable-kill test needs a sysctl/boot-id the codex sandbox lacks). Not
  a finding: `go test ./...` is green in the implementer's environment and this PR
  changes no `internal/runner` files. 1-line rationale recorded.

## Signoffs
<!-- Each active participant (implementer included) APPENDS their signoff block. -->

### Signoff: claude — 2026-06-19
Status: ✅ ACCEPT
Notes: Implementer. All three agreed fixes applied in fix-up cycle 1 (Decision Log
deviations cross-reference; "pipeline" added to the general rigor trigger in Phase 4 +
both Phase 5 spots; Phase 4 idempotence clause split into its own paragraph). Drift
guard + `go test ./...` + `go build` re-verified green; the two copies' diff stayed
symmetric (13/13). codex's runner-test note dismissed as a sandbox artifact (no runner
files touched). Zero remaining agreed fixes.

### Signoff: codex — 2026-06-19
Status: ✅ ACCEPT
Notes: I verified the three agreed fixes are present in the live `COOPERATION.md` and accept the consensus disposition. No remaining agreed fixes from my side; the runner-test note remains dismissed as a sandbox artifact.

### Signoff: hermes — 2026-06-19
Status: ✅ ACCEPT
Notes: Verified all three of my findings are resolved: the Decision Log now references `## Deviations from FINAL.md` above (byte-identical in both COOPERATION.md copies), "pipeline" is present in the general self-containedness trigger in Phase 4 and both Phase 5 spots, and the idempotence clause is split into its own paragraph. Zero remaining agreed fixes from my side.
