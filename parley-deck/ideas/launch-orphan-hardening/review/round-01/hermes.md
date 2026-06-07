---
agent: hermes
idea: launch-orphan-hardening
round: 1
phase: review
date: 2026-06-07
---

## 1. CORRECTNESS of best-effort manifest write
- File: internal/runcontrol/runcontrol.go:82-92: On writeManifest failure, appends "run.manifest_deferred" event (best-effort, _ ignored) and continues; Create returns success.
- Verified callers (no hard requirement for run.json):
  - internal/runstate/runstate.go:181-183: loadManifestSnapshot does `runmanifest.Load`; on err returns `Manifest{}, false` then applyManifestDefaults only fills fields if present (graceful, no crash).
  - internal/app/app.go:880-882: `manifest, err := runmanifest.Load(...)`; `if err == nil { payload.Manifest = &manifest }` (optional field).
  - internal/app/app.go:940-942: `if manifest, err := runmanifest.Load(...); err == nil && manifest.Status != "" { return manifest.Status }` (falls back to other status).
- No evidence in resume/continue/driver paths of hard `run.json` requirement (events.jsonl + runStore define the run).
- Conclusion: SAFE. Non-fatal does not hide a launch-critical error; manifest is convenience metadata only. No CRITICAL/MAJOR.

## 2. Longer retry window (~1.9s, 8 attempts)
- File: internal/fsutil/fsutil.go:29-37: retryDelays = [0,15,35,100,250,500,1000]ms (7 sleeps after initial, total 8 mkdir attempts).
- Healthy path (line 49-52): single `mkdirAll` success returns immediately (no Stat, no loop, no sleep). Zero cost to success path.
- Permission fail-fast (line 58-60): returns on first ErrPermission without sleeps.
- Window sized to outlast ~1s virtio-fs cache; reasonable. No MAJOR issues.

## 3. run.manifest_deferred Append + MkdirAllResilient
- File: internal/runcontrol/runcontrol.go:87-91: `_ = runStore.Append(...)` (the Append path itself calls MkdirAllResilient on runDir).
- If that MkdirAllResilient also fails it is silently ignored; run still proceeds.
- Per IMPLEMENTATION.md note: acceptable (no orphan created; events already written earlier). No issue found. NIT only (could log but TUI constraint noted in FINAL.md).

## 4. TESTS
- internal/runcontrol/runcontrol_test.go:75-116 (TestCreateBestEffortManifest): correctly forces failure via seam (line 82), asserts Create err==nil (92), RunID/RunDir present (94), runmanifest.Load fails (97), both run.created + run.manifest_deferred events present (113). Seam restored via t.Cleanup (83). Passes criteria.
- internal/fsutil/fsutil_test.go: (Test_GenuineFailure): 8 mkdirCalls (105,120), distinct last err on attempt 8 (114-118), sleeps exactly [15,35,100,250,500,1000] (123-130; first 0ms retry has no sleep). Correct for new schedule.
- No findings.

## 5. BUILD/REGRESSION
- Changes isolated to error path + test seam. Normal Create (manifest success), --no-tui, status, resume, Home all unaffected.
- go build/vet/test would pass (per IMPLEMENTATION.md verification + code inspection). No regression.

No findings at CRITICAL, MAJOR, MINOR.

Overall verdict: ACCEPT
