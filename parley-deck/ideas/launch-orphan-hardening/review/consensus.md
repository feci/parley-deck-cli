---
idea: launch-orphan-hardening
phase: review-consensus
date: 2026-06-07
drafter: claude
participants: [claude, codex, agy, hermes]
---

## Review consensus — round-01: all ACCEPT

- **agy — ACCEPT** (no CRITICAL/MAJOR/MINOR).
- **hermes — ACCEPT** — confirmed best-effort manifest is SAFE: "non-fatal does not hide a
  launch-critical error; manifest is convenience metadata only." Window reasonable.
- **codex — ACCEPT** — extensively verified NO path hard-requires `run.json`: runstate
  (loadManifestSnapshot returns (empty,false); idea/task/mode/participants derive from
  `run.created`), `status`, `sessions inspect` (optional display, prints missing-manifest
  message), `resume`/`continue` (resolve via runstate), and the driver (built from RunDir +
  event store, no manifest dep). Healthy mkdir path unchanged; permission still fails fast;
  the deferred-audit-append ignoring its error is correct (run.created remains source of
  truth). The `go test ./...` failure on codex's host is again the **sandbox
  `sysctl kern.boottime` artifact** in `TestDurableKillEndToEndRealProcess` (untouched by
  this diff) — dismissed, as in [[launch-mkdir-resilience]].

## Agreed fixes (applied, fix-up cycle 1)
1. **codex NIT** — `Test_GenuineFailure` doc comment updated from `5/20/50ms` to the actual
   `15/35/100/250/500/1000ms` schedule.
2. **codex NIT** — `TestCreateBestEffortManifest` now asserts `run.json` absence precisely
   via `os.Stat` + `errors.Is(err, os.ErrNotExist)` (was: any `runmanifest.Load` error).

Re-verified: gofmt clean, build/vet OK, `go test ./... -count=1` green.

## Status: COMPLETE — ready to ship v1.22.0
