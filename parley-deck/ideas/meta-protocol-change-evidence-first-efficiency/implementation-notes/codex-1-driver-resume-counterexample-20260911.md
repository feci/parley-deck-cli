---
idea: meta-protocol-change-evidence-first-efficiency
author: codex-1
date: 2026-09-11
status: reproduced-open
source: 2194dca
---

# Driver resume still resets the configured step ceiling

After the new launch-policy checkpoint, Codex executed an additional negative
probe against the actual Driver.Run implementation. With MaxDriverSteps=1,
the first Run dispatches round 2 and halts at its cap. Constructing a new
Driver over the same idea/run and calling Run again dispatches round 3.
Observed runner calls: [2 3], expected [2]. Exit 1, package duration 0.367s.

This is a known incomplete D6/AC-B1 boundary, not a failure of the new launch
policy's six process surfaces. Driver.Run currently initializes steps=0 and
start=time.Now() on each entry in internal/driver/loop.go. A green ordinary
suite does not prove lifetime driver limits. Direct Advance and other manual
action paths also need the shared precharge boundary, not just a Run-local
counter fix.

## Executed reproduction

The test was added through a Go overlay; no failing test was silently placed
in the passing source suite. Existing package fixtures supplied valid protocol
artifacts; production Driver.Run, Advance, cursor and event persistence were
exercised with a fake process-dispatch seam. This is not a live-model trial.

Runtime files preserved on integration:

- .parley-runtime/driver_resume_budget_probe_test.go
- .parley-runtime/driver-resume-probe-overlay.json
- .parley-runtime/driver-resume-probe.log

Command: go test -count=1 -overlay .parley-runtime/driver-resume-probe-overlay.json
./internal/driver -run '^TestDriverResumeCannotResetStepCapRecoveryProbe$'

Reproducible package-driver test body:

```go
func TestDriverResumeCannotResetStepCapRecoveryProbe(t *testing.T) {
    parts := []string{"codex", "claude"}
    ideaDir, runDir := setupIdea(t, parts, "")
    writeAll(t, ideaDir, 1, parts)
    appendEvent(t, runDir, "round.completed", "round-01")
    fr := &fakeRunner{writeOnRun: func(round int) {
        writeAll(t, ideaDir, round, parts)
    }}
    first := newTestDriver(ideaDir, runDir, parts, 3, true, fr)
    first.cfg.MaxDriverSteps = 1
    if err := first.Run(context.Background()); err != nil { t.Fatal(err) }
    if len(fr.calls) != 1 { t.Fatalf("first run: %v", fr.calls) }
    resumed := newTestDriver(ideaDir, runDir, parts, 3, true, fr)
    resumed.cfg.MaxDriverSteps = 1
    if err := resumed.Run(context.Background()); err != nil { t.Fatal(err) }
    if len(fr.calls) != 1 { t.Fatalf("resume exceeded lifetime cap: %v", fr.calls) }
}
```

## Next implementation boundary

Reuse the durable shared reservation store and immutable policy activation,
charging applicable actions before execution across driver/manual/resume/BLOCK
paths. Keep failed attempts spent and wall time anchored across resumptions.
Preserve the existing inclusive fixup/cross-review limits and allow verification
of the final allowed attempt. Status queries, awaits and zero-fix signoffs must
not consume fixup cycles. Do not derive safety counts from participant prose
or silently treat historical accounting as zero. Operator migration/extensions
remain actual controls, not editable frontmatter. This correction is still open.

