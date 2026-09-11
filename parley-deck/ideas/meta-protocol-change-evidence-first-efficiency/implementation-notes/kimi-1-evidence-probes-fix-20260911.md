---
agent: kimi-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-11
artifact-kind: owner fix handoff
base-commit: 89a4305 (working-tree edits only; no git mutations)
not-a-signoff: true
---

# Kimi evidence slice — probes fix 2026-09-11 (OWNER-RUN, NOT INDEPENDENT)

Status: the three independently reproduced findings from codex-1's
2026-09-11 review (`.parley-runtime/codex-independent-probes-20260911.log`)
are fixed in my `internal/evidence/**` slice and pass focused owner-run tests,
including codex-1's own overlay probes used read-only. This is NOT an
independent verdict, NOT Phase-6 review, NOT whole-idea completion, and NOT
integration at the production close call site (codex-owned). All tests below
were executed by the owner; another participant must reproduce before
integration.

Changed files (working tree only, nothing committed):

- `internal/evidence/execute.go` — capture accounting + strict test2json parsing
- `internal/evidence/evidence.go` — count validity on retained verifier executions
- `internal/evidence/execute_test.go`, `internal/evidence/evidence_test.go` —
  new permanent regression tests

## Fix 1 — [MAJOR] verifier skipped counts: reconcileRerun count validity

`reconcileRerun` (internal/evidence/evidence.go) previously compared skipped
counts only when both were nonnegative, so a retained gotest-json re-run with
`SkippedCases` -1, -2 or -100 passed both `AttestExecution` and `Evaluate`.
Now, mirroring the record-level rule:

- Anything below -1 in any of the three counts on the retained execution is
  malformed in every format and is refused ("independent re-run has an invalid
  (negative) case count"). (Executed/failed negatives were already refused by
  the existing `<=0` / `!=0` checks; skipped was the gap.)
- A `gotest-json` re-run must carry a real skipped count: a true
  `go test -json` parse always yields one, so -1 there is refused
  ("carries no skipped-case count — not a real go test -json parse").
- `-1` remains legitimate ONLY for the envelope format, which really lacks the
  field (the `Envelope` struct has no `skipped_cases`). Verified both ways:
  envelope re-run with -1 attests and closes; with -2 it is refused.

Format compatibility cannot bypass reconciliation: executed counts are always
reconciled (a re-run must be >0, and equals the original's), skipped counts are
reconciled whenever the re-run's format carries the field, and the command
hash, output hash, exit status and before/after tree digests bind regardless
of format. I checked the cross-format case (envelope re-run -1 vs gotest-json
original 0): no deeper issue — nothing count-like goes unreconciled, and an
opaque shell re-run remains refused as an attestation basis.

## Fix 2 — [MAJOR] duplicate Go-test Action fields: strict event parsing

`ParseGoTestJSON` used `json.Unmarshal`, whose last-wins, case-insensitive
field match decoded `{"Action":"fail","Action":"pass","Test":"TestX"}` as a
pass. Replaced with the same discipline as the envelope parser:

- New `parseTestEventStrict` token-walks each `{`-prefixed line and rejects:
  duplicate `Action`/`Test` fields (either order, fail→pass and pass→fail),
  case-aliased semantic fields (`action`, `TEST`), non-string semantic values,
  trailing content after the object, and malformed event lines.
- New `parseGoTestJSONEvidence` returns `recognized=true` + a non-nil error
  for such streams. `RunCriterion` treats that as fail-closed: the record
  keeps `FormatGoTestJSON` with ALL counts unknown (-1), `StatusFail`, and a
  diagnostics header — there is NO fall-through to an opaque shell PASS and no
  fall-through between parsers (an envelope line still suppresses test2json
  parsing entirely).
- The exported `ParseGoTestJSON` keeps its 4-result signature (the
  facilitator's overlay probes compile against it unchanged) and reports a
  rejected stream as recognized with zeroed counts — never as a pass.
- Real-stream support preserved: the full real field set (Time, Package,
  Output, Elapsed), unknown/future fields including nested values, duplicate
  NON-semantic fields (no case effect), package-level and build-failure
  events, and plain noise lines (build text, `FAIL\tpkg [build failed]`) are
  all still handled. Verified against three REAL captured test2json streams
  (see below).

## Fix 3 — [MINOR] capture boundary accounting

`cappedWriter.Write` computed remaining capacity AFTER appending, so a 16-byte
write to a 16-byte cap, split 8+8 writes, and a single 9-byte write all set
`overflow` falsely. Overflow is now decided from the PRE-WRITE remaining
capacity: a chunk that fits exactly is retained in full; only a chunk larger
than the room left overflows. Covered by exact-bound, split-bound,
under-bound, one-byte-over (single and split), crossing-split, and
empty-write-at-bound cases.

## Exact commands and results (owner-run)

Environment for every command (the write sandbox permits only the worktree;
`/private/tmp` and the default `~/Library/Caches/go-build` are not writable —
probe: `touch /private/tmp/kimi-write-test` → "Operation not permitted"):

```
export TMPDIR=$PWD/.parley-runtime/tmp GOCACHE=$PWD/.parley-runtime/go-cache \
       GOTMPDIR=$PWD/.parley-runtime/tmp GIT_CEILING_DIRECTORIES=$PWD/.parley-runtime/tmp
```

1. Focused evidence suite (all prior focused tests + all new regression tests):

   ```
   go test -count=1 -p 1 -json ./internal/evidence \
     -run 'Evidence|Criterion|Envelope|TreeDigest|Typed|Verifier|Barrier|AttestExecution|HelperProcess|ParseGoTestJSON|CappedWriter|RerunCountValidity|InvalidRerunSkipped'
   ```

   exit=0. Terminal events: **44 pass, 0 fail, 1 skip**. Full stream retained
   at `.parley-runtime/kimi-evidence-probes-fix-20260911.jsonl`.
   New permanent regression tests among the passes:
   `TestCappedWriterBoundaryAccounting` (7 subcases),
   `TestParseGoTestJSONStrictFailClosed` (11 adversarial streams),
   `TestParseGoTestJSONRealStreamsPreserved`,
   `TestRunCriterionMalformedGoTestJSONFailsClosed`,
   `TestRunCriterionGoTestJSONNoiseTolerated`,
   `TestAttestExecutionRerunCountValidity`,
   `TestEvaluateInvalidRerunSkippedCountsRejected`.

2. Focused app suite (unchanged from codex-1's pattern):

   ```
   go test -count=1 -p 1 -json ./internal/app \
     -run 'Evidence|Criterion|Envelope|TreeDigest|Typed|Verifier|Barrier|AttestExecution|HelperProcess'
   ```

   exit=0. Terminal events: **13 pass, 0 fail, 0 skip**. Stream retained at
   `.parley-runtime/kimi-evidence-probes-fix-app-20260911.jsonl`.

3. Facilitator overlay probes, used READ-ONLY (source untouched):

   ```
   go test -count=1 -overlay .parley-runtime/codex-evidence-probes-overlay.json \
     ./internal/evidence -run '^TestCodexProbe' -v
   ```

   exit=0. All five PASS, including the three that failed at 89a4305:
   `TestCodexProbeNegativeVerifierSkipped`, `TestCodexProbeCaptureExactBound`,
   `TestCodexProbeDuplicateGoTestAction`; the two previously-passing probes
   still pass. Log: `.parley-runtime/kimi-codex-probes-after-fix-20260911.log`.

4. Real-stream compatibility check (throwaway module-internal main at
   `.parley-runtime/tmp/parsecheck/main.go`, feeds captured real
   `go test -json` output through `ParseGoTestJSON`):

   - my evidence stream → `executed=44 failed=0 skipped=1 recognized=true`
     (exact match to the run's own tallies)
   - my app stream → `executed=13 failed=0 skipped=0 recognized=true`
   - codex-1's retained stream
     (`.parley-runtime/codex-independent-evidence-20260911.jsonl`) →
     `executed=42 failed=0 skipped=0 recognized=true`, matching codex-1's
     reported "42 terminal test-pass events, zero skips/failures"

## Environment incident (reported, not suppressed)

The first two focused runs with `GOCACHE=.parley-runtime/gocache` failed at
LINK time with `unexpected fault address` / `SIGBUS` inside `cmd/link`
(`decodetypeName`, deadcode pass), for both packages. `go build`/`go vet` of
the package succeeded; only test-binary linking faulted. This is an
environment-level fault of the Go build cache on this network-mounted volume
(mmap page-in), not a code failure — switching to the other prepared cache dir
(`GOCACHE=.parley-runtime/go-cache`) produced a clean, repeatable pass (runs
1–4 above). The failing logs are overwritten by the passing streams; the
SYMPTOM is recorded here. I make no causal claim beyond this observation;
codex-1's local-disk runs are unaffected by it. Independent verifiers on local
storage should use their normal TMPDIR/GOCACHE.

## Skip recorded accurately (not a pass)

`TestTreeDigestUnsupportedEntryFails` **SKIPPED** on this mount, exactly as
before:

```
tree_report_test.go:102: no unsupported-entry probe available on this
filesystem: listen unix parley-evidence-probe.sock: bind: operation not supported
```

(mkfifo → not supported; the relative-path unix-socket fallback → bind not
supported). Codex-1 independently ran this test to PASS (0.07s) on local
TMPDIR at 89a4305; the owner sandbox still cannot exercise FIFO/socket
creation, so local-capable fixture verification remains an independent
responsibility. The skip is a coverage gap on MY runs, not a defect claim.

## Remaining gaps / non-claims

- Production close call-site integration (actual independent
  execution/attestation wiring) remains codex-owned and is untouched here.
- The trusted-caller boundary is unchanged: `AttestExecution`/`Evaluate`
  guarantee internal consistency and binding of persisted evidence; they do
  not prove a caller actually executed anything, and runtime attribution is
  not authentication of a human (FINAL D3).
- All results above are owner-run; per the protocol they earn nothing until
  independently reproduced. No signoff, no whole-idea acceptance, no FINAL or
  IMPLEMENTATION edits are made or implied by this handoff.
