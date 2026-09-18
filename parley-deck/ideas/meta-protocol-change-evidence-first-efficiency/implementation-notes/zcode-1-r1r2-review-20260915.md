---
agent: zcode-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-15
kind: independent-source-review
reviews: claude-1 R1/R2 native candidate
candidate-hashes:
  source: 4f35b839bf8762c8a211c0efc7b60131fd28e16cdcdff0cf04a9c976b5b7ffbb
  new-test: c6fc513b89f501a8c278eaeae37119940bdfe9669890a5b1c0727f689a893
  app-test: 28791fbbb50ea69b7aec58859cd1c9d795d8b6c219634f68e9ab40ab798a2732
status: supporting-review
---

# zcode-1 source review of the claude-1 R1/R2 candidate (2026-09-15)

Read-only source review of `.parley-runtime/input/claude-r1r2-candidate/` (not this
checkout's baseline). I read the candidate `reservation_recovery.go`,
`reservation_recovery_validation_test.go`, `agents_exec_test.go`, the owned candidate
note, the stop-diagnosis results JSON, and the checkout's HEAD
`reservation_recovery.go`, `state.go`, `state_test.go`, `budget/ledger.go`,
`budget/cycle_binding.go` as dependencies. No execution: every test outcome below is
claude-1's attributed testimony (its note plus `results-*.json` with log SHA-256s),
not my run. Hash values are quoted from claude-1's artifacts; I could not recompute
them in this restricted launch. This is a supporting review of R1/R2 only — no
whole-audit signoff, and the R3 and N1 cross-review/step residuals remain open and
separate.

## Refutation attempts

### R1-a — original intent/policy/ledger/state identity preserved

I compared HEAD's single-pass body against the candidate's `withReservationRecovery`
plus `recoverReservationChecked` check by check; every HEAD check reappears, under the
guard, in both observations: binding load and guarded reload
(`reservation_recovery.go:303-323` candidate vs `:287-307` HEAD), ledger Inspect
(`:325` vs `:308`), canonical state read (`:328` vs `:312`), intent read and root
identity (`:332-343` vs `:316-327`), `validateIntentBefore` and `compareIntentPrefix`
(`:344-349` vs `:328-333`), the charged-row branch with the verbatim missing-row
exactness test (`:360-362` vs `:355`) and existing-row identity (`:365-369` vs
`:360-364`), and `validateState` on the post-append history (`:371` vs `:365`). The
identity/argument checks moved intact to `:386-392`. One delta: HEAD's uncharged
branch set `p.Status` only after its checks passed (`HEAD :342`); the candidate's
capture sets it before the content check (`candidate :351`), so an error return can
carry a populated preview where HEAD's did not — cosmetic, error-path only.

### R1-b — guarded publication

Publication (`persist` for a missing row, durability sync otherwise) runs inside the
second observation's callback (`:406-431`), i.e. under the cycle guard, after the
equality gate `reservationRecovery.same` (`:293-297`) requires unchanged store
dir/scope, policy JSON, ledger JSON, raw state bytes, missing flag, state JSON and
preview digest. The equality set is well-founded: `budget.Snapshot` is plain data
(`ledger.go:83-88`), as is `CycleBinding` (`cycle_binding.go:34-37`). The no-charge
refusal and expected-preview check precede `syncIntent` and persist, matching HEAD's
order (`HEAD :374-389`).

### R1-c — missing-row/replay and bounded retry

The restart is bounded by `changed && attempt == 0` (`:433-435`): at most two full
passes, and every non-drift failure (structural, content, identity, publication)
returns without retry — matching the note's claims. The replay path is real in code:
after a concurrent publish, state bytes differ, `same` fails, the restart re-observes
the now-existing row, the existing-row identity check passes, and only the durability
sync runs — `persist` is unreachable on that path (`:414-425`).

### R1-d — Begin-derived Finish witness

The decisive test uses a real `Begin` handle (`validation_test.go:48`), a real child
process changing source and exiting 7 (`:63-68`), the public
`PreviewReservationRecovery` (`:59`), and a real `run.Finish` under a 2 s context
(`:105`). `chargeFixture` charges through the unmodified `Observer`
(`state_test.go:48-58`), so no fixture-constructed handle is involved. The pause rides
the unexported `check` seam, but the wrapper runs the real `checkStateSnapshots`
first (`:78-81`); the public entry points bind that real check unconditionally
(`reservation_recovery.go:278-280`). The contention being fixed is genuine in source:
`Finish` publishes via `withStateControl` (guard + `validateState` only,
`state.go:607-663`, `:482-522`), while HEAD held the same guard across the full
archive/resolution content check (`checkStateSnapshots`, `state.go:665-698`).

### R2 — failure-output diagnostic

The diff to `agents_exec_test.go` is exactly a two-line comment plus `stdout=%s` in
the failure `Fatalf` (candidate `:44-46` vs HEAD `:44`); assertions and fixture are
unchanged. The printed `--json` record is content-free by test design (candidate
`:55-59` assert no content leak). R2 aids classification of a future failing run
only; it cannot classify the original historical failure, whose run printed no
record.

## Findings

### [MINOR] Unguarded content window is not closed for archive/resolution bytes
The recheck equality covers policy/ledger/state/intent/preview but never re-reads
archive or resolution content; `withReservationRecovery` "reads no archive or
resolution content" (candidate `:282-301`). A guard-respecting writer that mutated
archive bytes without a state write would land between the content check and the
recheck undetected. I found no such writer: production `CaptureSnapshot` callers
(`state.go:361`, `state.go:635`, `continuation.go:207`) all write state under the
guard too, which `same` catches; non-guard tampering exposed HEAD equally. Record
the invariant (no archive-only guarded writers) or extend the recheck.

### [MINOR] Uncharged-intent drift has no test witness
The `no-intent-preview-recheck` control passed (claude-1 testimony), so nothing in
the suite pins the intent/preview clauses of the recheck on the
`intent-without-published-charge` path. By reading, the production `same` does cover
it — the preview digest embeds `IntentSHA256` (`:351`, `:296`), and the second
observation re-reads and re-validates the intent — but that is exactly the
by-reading-only status the note discloses (Limitations 1, 3). Add the subtest.

### [MINOR] Registered-stop failure is environment-correlated and unexplained
Attributed testimony: the stop test failed in claude-1's restricted environment on
the candidate (22.1 s) and on verbatim HEAD (22.9 s), both exit 1
(`results-stop-diagnosis-20260915T231025.json`); the coordinator reports the same
candidate passing in the default environment at 6.1 s. Both failing together in the
restricted environment shows the failure there is not candidate-specific, but it does
not causally exclude anything about the original historical failure (a different,
unclassified run), and the cause remains unresolved. This launch still has no
registered-Stop contention evidence; keep that item open and diagnose separately.

### [NIT] Error-path preview payload divergence
See R1-a last sentence: on an uncharged content-check failure the candidate returns a
preview with `Status: "intent-without-published-charge"` where HEAD returned an empty
Status. No consumer of the error-path preview was found; harmless.

## Dispositions evaluated

- Note Limitation 4 (concurrent change during `checkStateSnapshots` surfaces as
  content error, not drift, and is not retried): concur — a disclosed fail-closed
  trade; HEAD would have waited up to 30 s under the guard.
- Note Limitation 5 (equality assumes no read-time-computed JSON fields): the struct
  definitions support it (`ledger.go:83-88`); the ledger `read` path remains unread
  by me, and a false drift only refuses (fail-closed). Low risk.
- Note's stop conclusion ("this candidate does not cause it"): valid only for the
  restricted environment where both versions fail identically; not a general causal
  exclusion — see the MINOR above.

## Verdict summary

No CRITICAL or MAJOR finding; the four refutation targets (identity preservation,
guarded publication, missing-row/replay exactness, bounded retry) survived source
attack at every locator I tried, and the Finish witness is structurally real.
Non-refutation is bounded by this review's scope and by attributed, unexecuted test
testimony. R1/R2 remain open pending the normal non-owner gate; this note is that
input, not a close.

## Self-correction (2026-09-16; appended at owner direction — every byte above preserved)

### 1. Header `new-test` hash mistranscribed — corrected value is attributed

The header's `candidate-hashes.new-test` line reads
`c6fc513b89f501a8c278eaeae37119940bdfe9669890a5b1c0727f689a893` — 61 hex characters;
the `e78` before the final `9a893` was dropped when I copied it. The corrected value is
`c6fc513b89f501a8c278eaeae37119940bdfe9669890a5b1c0727f68e789a893` for
`internal/trajectory/reservation_recovery_validation_test.go`. It is **attributed, not
computed**: quoted from claude-1's owned candidate note (§ Hashes, "Tested and frozen"
table) and byte-identical in the `sources_before`/`sources_after` entries of all three
`results-*.json` manifests (candidate, controls, stop-diagnosis — seven occurrences,
exact-string verified this session). I could not recompute the digest in this restricted
launch; provenance remains claude-1 testimony (SECONDARY), as the intro already
discloses for all quoted hashes. The header line above is superseded by this value; the
mistranscription was mine — the candidate's artifacts were internally consistent.

### 2. Stop-failure causal disposition narrowed

My MINOR finding said the two failures together 'shows the failure there is not
candidate-specific', and the note-disposition under `## Dispositions evaluated` above
was 'valid only for the restricted environment where both versions fail identically'.
Too strong. Narrowed: the stop-diagnosis runs establish **occurrence on both versions**
in that environment — the restricted environment alone reproduces the failure on the
note's verbatim-HEAD package, so the failure there cannot be attributed solely to the
candidate. They do not by themselves causally exclude a **candidate-specific
contribution** to the failure's occurrence or timing there: an identical symptom (same
test line, same message, both exit 1) does not prove an identical causal mechanism, and
nothing in those two runs discriminates an environment-only cause from an
environment-triggered, candidate-amplified one. The original historical failure remains
unclassified by this evidence. The finding itself — environment-correlated, unexplained,
kept open for separate diagnosis — is unchanged; no finding is altered, withdrawn, or
suppressed by this correction.
