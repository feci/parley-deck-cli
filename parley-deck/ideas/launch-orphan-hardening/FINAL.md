---
idea: launch-orphan-hardening
phase: final
date: 2026-06-07
drafter: claude
implementer: claude
status: final
participants: [claude, codex, agy, hermes]
---

## Problem (follow-up to launch-mkdir-resilience / 1.21.0)

After shipping 1.21.0 (`fsutil.MkdirAllResilient`, 75ms window), the owner STILL hit
`launch failed: mkdir .../runs/<id>` in `parley tui` → N. Verified diagnosis: the run is
created (idea + `events.jsonl` with `run.created`) but **`run.json` is missing** → the
deterministic failure is **`runmanifest.Write`** (the step after `Append`), whose
`MkdirAllResilient(runs/<id>)` fails on the *already-existing* run dir under virtio-fs
stale-cache. The 75ms retry window is shorter than virtio-fs's ~1s attr/entry cache
timeout. The failure could NOT be reproduced in 30k+ iterations (raw, multi-process,
600 real `runcontrol.Create` + concurrent scanner) — it is load-dependent. `runmanifest`
is written ONLY at create time (runcontrol.go), so a single failure orphans the run.

## Decision (owner-approved)

Two changes, both hardening the SAME class without depending on the exact cache window:

1. **Lengthen the retry window** in `fsutil.MkdirAllResilient`: `retryDelays` 75ms →
   `[0,15,35,100,250,500,1000]ms` = **8 attempts, ~1.9s worst case**, error-path only,
   healthy first-success path unchanged (no Stat, no sleep). Sized to outlast a virtio-fs
   ~1s cache timeout. Permission errors still fail fast.

2. **Make the manifest write non-fatal** in `runcontrol.Create`: a `runmanifest.Write`
   failure must NOT orphan an already-created run. On failure, append a
   `run.manifest_deferred` audit event and **continue**. Safe because the run is defined by
   `events.jsonl` and every `runmanifest.Load` caller already degrades gracefully when
   `run.json` is absent (verified: `runstate.loadManifestSnapshot`, app.go:880/940 all
   guard the error). Orphaning is now **deterministically impossible**.

Not written to stderr (would corrupt the Bubble Tea TUI) — recorded as a store event.
A seam `var writeManifest = runmanifest.Write` makes the best-effort path testable.

## Verification
`gofmt`/`go build`/`go vet` clean; `go test ./... -count=1` green. New
`TestCreateBestEffortManifest` (forces the manifest write to fail → Create still succeeds,
run.json absent, `run.created`+`run.manifest_deferred` present). `Test_GenuineFailure`
updated for the new 8-attempt schedule.

## Ship
Bump 1.21.0 → 1.22.0; local build for the owner to test in a FRESH `parley tui`; then PR,
tag v1.22.0, GitHub release + Windows assets, Homebrew bump.
