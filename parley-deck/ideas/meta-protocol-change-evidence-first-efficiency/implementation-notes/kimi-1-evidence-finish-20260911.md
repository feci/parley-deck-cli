---
agent: kimi-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-11
artifact-kind: owner recovery handoff
supersedes: checkpoint 7f2676f partial work (preserved, not deleted)
not-a-signoff: true
---

# Kimi evidence slice — finish attempt 2026-09-11 (OWNER-RUN, NOT INDEPENDENT)

Status: the six audit priorities from
`implementation-notes/codex-1-kimi-recovery-check-20260910.md` are addressed in
code and pass focused owner-run tests. This is NOT an independent verdict, NOT
Phase-6 review, NOT AC-E1/E2 closure, and NOT a whole-idea completion claim.
All tests below were executed by the owner (kimi-1); another participant must
reproduce them before integration.

## Priority 1 — AttestExecution binding (internal/evidence/evidence.go)

- `Provenance.VerifierTreeSHA256` (a bare name + tree stamp that discarded the
  rerun) is replaced by `Provenance.VerifierRerun *VerifierExecution`, which
  RETAINS the verifier's actual independent execution: the rerun's
  CommandEvidence plus `TreeBeforeSHA256`/`TreeAfterSHA256`.
- `AttestExecution(r, criterion, verifier, rerun VerifierExecution)` now
  refuses unless: the criterion record is pass-typed; the reran command hash
  equals the record's exact command hash (an unrelated command is refused);
  the output hash is non-empty; exit 0; structured format; >0 executed and 0
  failed cases; executed/skipped counts equal the original record's; and
  before == after == the report's tested tree digest.
- `Evaluate` re-reconciles every one of these bindings from the PERSISTED
  report at close time (`reconcileRerun`): a supplied verifier name without a
  consistent retained execution, a different command, diverging counts, a
  missing output hash, or a foreign/unstable verification tree all fail
  closed. A supplied name is never a verdict.
- Trusted-caller boundary is documented on `AttestExecution` and `Evaluate`:
  the API guarantees internal consistency and binding of persisted evidence
  and trusts the calling orchestrator to have ACTUALLY executed the rerun it
  attests; a same-UID caller fabricating mutually consistent hashes is outside
  the boundary. Runtime attribution is not cryptographic authentication
  (FINAL D3).

## Priority 2 — ParseEnvelope strictness (internal/evidence/execute.go)

`json.Unmarshal` (last-wins, case-insensitive key match) is replaced by a
token-walking `parseEnvelopeStrict` that rejects: duplicate fields (the
reproduced `executed_cases:0` then `:1` case), case-aliased fields
(`Executed_Cases`), unknown fields, null/non-integer values (strings,
booleans, fractions, nested values), non-object payloads, and any trailing
content after the closing brace. Existing negative/conflicting/partial checks
and the duplicate-envelope-LINES guard remain.

## Priority 3 — fixtures

- `TestTreeDigestUnsupportedEntryFails`: root cause was `git ls-files -c -o`
  omitting untracked non-regular entries entirely. `listTreeFiles` now
  cross-checks the working tree with `rejectUnsupportedEntries` (a type-only
  WalkDir, pruning only `.git`) so fifos/sockets/devices fail the digest
  instead of vanishing. HOST LIMITATION: this host's mandated
  TMPDIR volume supports neither FIFOs (`mkfifo` → EOPNOTSUPP) nor unix
  sockets (`bind` → EOPNOTSUPP; absolute paths also exceed the AF_UNIX
  sun_path limit at 113 chars), so the fixture SKIPs here with an explicit
  reason after trying fifo then a relative-path socket. On a capable host it
  exercises the rejection. The walk's no-error path is exercised here by
  every other TreeDigest test.
- `TestRunChecksContractEvidenceWriteFailureVetoes`: fixture now contains an
  in-scope `code.go` so the pre-execution tree digest succeeds and the veto is
  exercised at the intended persistence step (it previously died earlier at
  the empty-tree digest).
- `TestEvidenceCloseGatePositive`: now performs a REAL independent rerun — the
  verifier re-executes every criterion via `RunCriterion`, computes
  before/after tree digests, attests via the new bound API, persists the
  report, and then the close gate allows.
- `TestSerialVsBarrierConcurrencyFixture`: netcat removed. Both halves are the
  test binary re-executing itself (`TestHelperProcess`, env-gated): the server
  listens on loopback, publishes the port atomically, blocks in Accept; the
  client polls the port file and exits 1 when no listener appears. Serial
  execution fails both halves (explicit failure); concurrent execution passes
  both with structured executed-case proof (corrected success). No nc, no
  FIFOs, no port pre-allocation race.

## Changed files (all inside the owned slice)

- internal/evidence/evidence.go — VerifierExecution, retained/bound
  AttestExecution, reconcileRerun in Evaluate, trusted-caller docs.
- internal/evidence/execute.go — strict envelope parser.
- internal/evidence/tree.go — rejectUnsupportedEntries cross-check.
- internal/evidence/evidence_test.go — new-API positives + adversarial
  attestation/reconciliation cases.
- internal/evidence/execute_test.go — duplicate/aliased/unknown/trailing/
  string/fractional/array envelope cases + RunCriterion duplicate-field case.
- internal/evidence/tree_report_test.go — fifo→socket→explicit-skip fixture
  (also gofmt-normalized; it carried pre-existing drift).
- internal/app/driver_checks_test.go — write-failure fixture reaches the
  intended step.
- internal/app/driver_evidence_test.go — positive gate with real persisted
  rerun; portable helper-process barrier fixture.
- This handoff. No other source files touched. No git mutations performed.

EVIDENCE.json schema note: `verifier_tree_sha256` is gone;
`verifier_rerun{command,tree_before_sha256,tree_after_sha256}` replaces it.
Grep confirmed no references to the old field/API outside the owned slice, so
no integration glue is required for compilation. The runtime wiring of WHO
calls AttestExecution with real digests at close remains facilitator-side
call-site work (unchanged, already documented in driver_evidence.go).

## Executed commands and observed output (owner-run)

Env for all: `TMPDIR=$PWD/.parley-runtime/tmp`,
`GIT_CEILING_DIRECTORIES=$PWD/.parley-runtime/tmp`,
`GOCACHE=$PWD/.parley-runtime/go-cache`.

- `go vet ./internal/evidence ./internal/app` → clean.
- `go test -count=1 -p 1 ./internal/evidence ./internal/app -run
  'Evidence|Criterion|Envelope|TreeDigest|Typed|Verifier|Barrier|AttestExecution|HelperProcess'`
  → `ok internal/evidence 6.033s`, `ok internal/app 11.491s` (final run; an
  earlier identical run: 3.636s / 23.147s). Every previously failing fixture
  group now passes; TestTreeDigestUnsupportedEntryFails SKIP (host, above).
- ENVIRONMENTAL: the same two packages in one default (parallel) invocation
  twice gave `signal: killed` for internal/evidence at 0.002–0.003s while app
  passed; `go test ./internal/evidence` alone → ok 3.901s;
  `go test ./internal/fsutil ./internal/procctl` (parallel) → ok,ok;
  `go test ./internal/evidence ./internal/fsutil` (parallel) → ok,ok. The
  kill happens only when the large app test binary starts concurrently with
  another binary on this shared-volume host; it is not a code failure.
  Reviewers should use `-p 1` or per-package invocations here.

## Remaining failures / honest limits

- TestTreeDigestUnsupportedEntryFails is SKIP (not PASS) on this host —
  unsupported-entry rejection is code-fixed but demonstrated only on
  filesystems that can host fifos/sockets.
- Full `go test ./...` was deliberately NOT run (launch constraint: no broad
  loops); only the focused filter above plus the small package probes.
- `internal/app/protocol_test.go` has pre-existing gofmt drift; not mine,
  untouched.
- No independent verification, no review-consensus, no AC closure, no claim
  that the slice is integration-ready. Checkpoint 7f2676f and all prior
  artifacts preserved.
