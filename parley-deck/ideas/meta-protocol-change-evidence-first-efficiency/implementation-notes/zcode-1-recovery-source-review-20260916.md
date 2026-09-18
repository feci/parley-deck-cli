---
agent: zcode-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-16
reviewed-source: /private/var/folders/yt/p2sr23f12_qcfx_w2z5c1p4r0000gn/T/parley-recovery-publication-pmwo7ia9/source
kind: supporting source review of the cumulative verifier-refusal recovery candidate (read-only; no tests executed in this launch)
supersedes-layer-scope: kimi-1 attended app/relaunch layer + codex-1 recover-verifier-parent CLI layer in the newer publication tree
---

## Ownership boundary (§15.1)

I coauthored earlier layers (bound-verifier runner binding, verifier-parent lineage,
verifier-refusal recovery candidate/correction). I issue **no verdicts on my own prior
claims**; where those layers matter I quote code only as context or mark reliance on them
as outside this review. Verdicts below target the Kimi-owned app/library layer and the
Codex-owned CLI route in the reviewed tree. Coordinator runtime measurements are recorded
as **testimony**, never as my execution. All PRIMARY verdicts carry locator + quotation.

## Refutation attempts

R1 — "The publication route launches nothing." Tried to find any launch/charge path in
`runTrajectoryVerifierParentRecovery`. `internal/app/trajectory_verifier_recovery.go:331-335`:
apply calls only `trajectory.PublishRecoveredParent(...)`, preview only
`trajectory.PreviewRecoveredParent(...)`; no `runner.`/discovery calls in the function.
Verdict **CONFIRMED (PRIMARY)** that the route itself executes no process; anything deeper
sits in the library callback (R2) and the runner (not restated here).

R2 — "Publication writes only the immutable observation." In
`internal/trajectory/parent_recovery.go:576-633` the `withState` callback reads state,
reads/validates/writes only `parent-recovered.json` (`recoveredParentName`, line 19/405).
No write target names `parent-result.json` in lines 405-651. The doc comment at 643-648
states "Publication executes no model and spends, refunds or resets nothing." — consistent
with the code I read. Verdict **CONFIRMED (PRIMARY)** for the recovered-parent section;
`withState` lock/side-effect internals are not re-verified (residual unknown).

R3 — "Original failed parent result is preserved, not overwritten or conflated."
`parent_recovery.go:513-523` refuses unless the retained original is exactly the pre-start
budget refusal: Version 1, RunID match, RequestPath/RequestSHA256 equality,
`original.InvocationID != effective.Recovery.RefusedInvocationID` / TerminalSHA256 equality
(else error at 485/523: "original retained result is not the exact pre-start budget
refusal"), `ReceiptSHA256 != "" || original.Assessment != nil || original.FailureStage != "launch" || !original.TrajectoryPending`.
A missing parent-result refuses too (515: "requires the retained original refused parent
result"). Mirrored at the app layer, `trajectory_verifier_recovery.go:156-161`. Verdict
**CONFIRMED (PRIMARY)**.

R4 — "Apply is exact-digest gated; replay never rewrites." First apply:
`parent_recovery.go:620-623` — `if expected != p.SHA256() { return errors.New("recovered
parent evidence or state changed since preview") }` before
`publishRecoveredParent(dir, base, RecoveredParentRecord{1, p, expected, time.Now().UTC()})`.
Replay over a retained record: 598 validates it against freshly derived live evidence, and
601-602 — `if apply && expected != retained.SHA256 { return errors.New("recovered parent
observation replay changed its original preview") }` — after which only `syncRecoveredParent`
runs (604-608), never a rewrite. Verdict **CONFIRMED (PRIMARY)**.

R5 — "Preview is read-only." `PreviewRecoveredParent` (639-641) passes `apply=false`;
`recoveredParent` then never enters a publish branch (619-626 gated on `apply`). CLI side
`trajectory_verifier_recovery.go:326` rejects malformed authority combinations
(`(!*yes && *expected != "") || (*yes && *expected == "")` → exit 2). The new test asserts
`os.Stat` NotExist after preview (`trajectory_verifier_publication_test.go:73-74`).
Verdict **CONFIRMED (PRIMARY)**.

R6 — "Relaunch admission refuses duplicates and stale plans." Plan build
(`trajectory_verifier_recovery.go:144-186`) requires the captured recovery identity, the
retained refused parent (153-161), helper scope (162), roster-resolved headless verifier
(173-179). Apply rechecks the plan digest before launching (231-233: "plan changed since
preview; inspect again"). Protocol precheck runs before the pinned invocation is spent
(265-268). Post-launch battery: pinned-invocation escape check 279-281, terminal identity
282-288, request immutability 290-294, receipt containment 298-304, assessment 305-307.
App-layer guards **CONFIRMED (PRIMARY)**. The actual duplicate-relaunch/consumed-invocation
refusal lives in the runner binding I coauthored — **no verdict issued**; the end-to-end
assertion (`trajectory_verifier_publication_test.go:48` "duplicate relaunch reused its
consumed invocation") is test code, not my execution.

R7 — "Route reachable from the real CLI." `cmd/parley/main.go:10` —
`os.Exit(app.Run(os.Args[1:], os.Stdout, os.Stderr))`; `internal/app/app.go:64` —
`case "trajectory":`; `internal/app/trajectory.go:23-24` —
`if len(args) > 0 && args[0] == "recover-verifier-parent" { return runTrajectoryVerifierParentRecovery(...)`.
Usage surfaced at `trajectory.go:53`. Verdict **CONFIRMED (PRIMARY)**.

R8 — "Stranded-evidence mechanism." Without the route: relaunch cannot re-consume the bound
invocation (R6 context; runner-owned), the missing-parent path refuses a retained *failed*
parent (R3 code), so a completed replacement whose publication failed stays unresolved.
The mechanism is code-confirmed; whether that was "ergonomic" or a "recovery gap" is a
position difference — kimi-1 note line 225: "possible later ergonomic, not a correctness
gap (every read revalidates)." vs codex-1's recovery-gap reading under FINAL.md:252-254
("An external provider or spend stop leaves work pending... Do not overwrite the historical
assessment."). Not a verdict conflict: the facts agree, the classification differs.

## Findings

### [MINOR] Concurrent first-apply can replace the published observation's bytes
`parent_recovery.go:570` — `dir.Rename(stage, filepath.Join(base, recoveredParentName))`
atomically **replaces** an existing destination, and the record embeds a fresh timestamp
(line 623 `time.Now().UTC()`). Two concurrent applies that both observe NotExist (612-615
check-then-act) both pass digest validation and both rename; the loser's bytes are
overwritten with a same-lineage, same-digest record differing only in PublishedAt. The
comment (548-550: "A retained record is validated and replayed by the caller; it is never
rewritten here") holds under the attended single-operator model, not under concurrent
applies. This discipline is inherited verbatim from the missing-parent route ("same
stage-fsync-rename discipline", 548-549), so it is a pre-existing pattern, not a regression
introduced by this layer. Minimal correction: refuse when the destination exists before
rename, or drop the wall-clock timestamp from the canonical record.

### [NIT] Single-slot `observed` closure in relaunch
`trajectory_verifier_recovery.go:262-264` — `Observe: func(record telemetry.Record) {
observed = record }`. Any record arriving after the terminal record makes the 282-288
battery fail closed (safe), but whether Observe is synchronous/single-threaded is a runner
internals question I do not verdict on (coauthored layer).

### [NIT] Publication route does not bind `--dir` to a ticket root
Unlike relaunch (`readRetainedVerifierRun`, line 50: `req.Ticket.Root != abs` refusal),
`recover-verifier-parent` relies on `canonicalRoot` + `runtimeID` validation + live state
lookup (`parent_recovery.go:578-583`); a wrong dir/idea fails closed with no state. Design
note, not a defect.

## Testimony register (not my execution)

- codex-1 note lines 21-27: before-control failed in 15.658s (route absent), logs claimed
  at `managed-continuation-20260915/app-recovery-publication-control/before.*`; lines
  44-50: 87 passing Kimi-run test/subtest events, app package failing only on the
  rewritten-terminal expected message, now corrected. Launch supplement: after-control +
  exact rewritten-terminal guard test pass in 27.234s. **UNVERIFIED testimony** — I could
  not locate the `app-recovery-publication-control/{before,after}-result.json` files from
  this launch (globs under the worktree idea dir and the reachable temp roots returned
  nothing; the temp parent dir is permission-denied to my tools).
- The corrected guard expectation I did verify in-tree: `trajectory_verifier_recovery_test.go:526`
  ("original-parent-rewritten"), :542 ("replacement-terminal-rewritten"), :556 message
  "original requested lifecycle differs from the verifier terminal" — refusal remains
  mandatory; only the expected string changed.

## Unknowns / criteria needing executed tests

- Whole-scope suite has not been run; the focused app test was still running at launch.
  No pass claim is made here.
- `withState` lock/rollback semantics; runner-side consumed-invocation refusal and
  fresh-ID protocol-refusal recording; budget admission/settlement inside `RunConsult` —
  earlier/coauthored layers, deliberately not verdicted.
- Replay byte-equality and "no extra process/charge" are code-confirmed (R4/R1) and
  test-asserted (`trajectory_verifier_publication_test.go:44-58, 63, 88-101, 105-110`)
  but not executed by me.
- 1 MiB publication bound (`parent_recovery.go:553`) could refuse an unusually large but
  legitimate lineage — theoretical only.

## Position summary

The Kimi app/relaunch layer and the Codex `recover-verifier-parent` route, as source, hold
against every refutation I could mount read-only: admission, exact-digest authority,
original-failure preservation, replay immutability, and read-only preview all check out at
the quoted locators. The one concrete weakness found is the concurrent-apply rename race
(MINOR, inherited pattern). Completion evidence still requires the running tests to finish
and a whole-scope pass — that remains open, and nothing here accepts the candidate.

