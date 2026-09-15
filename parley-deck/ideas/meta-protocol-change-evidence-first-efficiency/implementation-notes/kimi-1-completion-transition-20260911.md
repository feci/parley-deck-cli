---
agent: kimi-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-11
artifact-kind: owner fix handoff
base-commit: 0a2022b (working-tree edits only; no git mutations)
not-a-signoff: true
---

# Kimi evidence slice — completion status-transition integrity 2026-09-11 (OWNER-RUN, NOT INDEPENDENT)

Scope: ONLY the exact status-transition integrity fix. Not Phase-6 review, not
final acceptance, not Codex-side integration (evidence_verify.go /
driver_impl.go are untouched; the wiring contract is below). Shared
writer/compare-before-save locking remains separately open and is NOT claimed.

## Problem (facilitator counterexample)

Independent verification succeeds through Driver.Advance; `Complete`
(driver_impl.go:478) then flips frontmatter to `status: complete`; any later
`EvidenceCloseGate` denied with "IMPLEMENTATION.md non-evidence content
changed after the evidence was recorded" — the flip changes the bound
non-evidence rest digest (`Report.ExtraDigests`). The suite never exercised
that last invariant. Status/frontmatter must stay bound; the fix is an exact
authorized BEFORE/AFTER transition, never an exclusion.

## Implementation (owned files only)

- `internal/evidence/completion.go` (new):
  - `TransitionFrontmatterStatus(doc, to) (out, from, err)` — pure,
    deterministic: rewrites exactly one NORMALIZED `status: <value>`
    frontmatter line, byte-exact elsewhere. Fail-closed on no/unterminated
    frontmatter, missing/empty/duplicate status, non-normalized line,
    multi-token/comment value, no-op target, or an invalid `to`.
  - `TransitionStatusToComplete(doc)` — the only transition the gate honors.
  - `CompletionTransition{Path, FromStatus, ToStatus, BeforeSHA256,
    AfterSHA256, AuthorizedBy}` on `Report` (`completion_transition`,
    omitempty). Old reports carry no field and stay valid only in their
    ORIGINAL state — never silently upgraded.
  - `AuthorizeCompletionTransition(r, path, currentContent, verifier)` — the
    adapter the independently invoked verifier calls AFTER attesting all
    criteria and BEFORE Save. Requires: non-nil report, non-empty verifier, no
    existing transition, path bound in ExtraDigests, `sha256(currentContent)
    == ExtraDigests[path]` (the authorized-from state IS the attested state;
    drift means re-record), ≥1 record all attested (Verifier + VerifierRerun)
    BY this verifier, verifier != executor, and a transitionable document.
    AFTER digest is recomputed by applying the transformation — caller
    digests are never accepted.
  - `VerifyCompletionTransition(r, path, boundDigest, currentContent,
    verifier) []string` — gate-side recomputation: right path, target
    `complete`, well-formed source, authorizer == selected verifier,
    BeforeSHA256 == bound digest, current digest == AfterSHA256, AND an
    independent inverse transformation of the current content reproduces
    exactly the recorded state (proves the flip is the ONLY change). Any
    tampered field fails a recomputation.
- `internal/evidence/evidence.go`: Report field + doc (above).
- `internal/app/driver_checks.go`: `implementationRestContent` exposes the
  exact bound bytes (IMPLEMENTATION.md minus the generated `## Validation
  evidence` section); `implementationRestDigest` now wraps it (unchanged).
- `internal/app/driver_evidence.go`: on a bound-digest mismatch the gate
  reconciles ONLY via `VerifyCompletionTransition`; otherwise the original
  deny stands (message preserved, reasons appended).

Trusted-caller boundary unchanged: internal consistency/binding is
guaranteed; the orchestrator is trusted to pass the real current document and
to be the attester it names. Runtime attribution is not same-UID
authentication (FINAL D3).

## Codex wiring contract (exact names + ordering)

1. Helper (evidence_verify.go, after the AttestExecution loop, BEFORE
   `evidence.Save`): `restContent, implRel, err := implementationRestContent(root, ideaDir)`
   then `evidence.AuthorizeCompletionTransition(&report, implRel, restContent, req.Verifier)`;
   on error fail the verification (nothing is saved).
2. Parent reconciliation (evidence_verify.go ~L390-407): `updated` now
   carries `completion_transition` — the current `reflect.DeepEqual(updated,
   original)` would fail "verifier changed the original evidence". Validate
   the transition, then exclude it from the comparison (like Provenance):
   non-nil, Path == the ExtraDigests key, AuthorizedBy == o.drafter, ToStatus
   == "complete", FromStatus well-formed, BeforeSHA256 ==
   `original.ExtraDigests[path]` == sha256 of the CURRENT rest content (doc
   not yet flipped at this point), AfterSHA256 == sha256 of
   `TransitionStatusToComplete(rest)`; then nil the field on `updated` before
   DeepEqual. Receipt flow (`UpdatedSHA256`) is unaffected.
3. Complete (driver_impl.go): replace the inline frontmatter loop with
   `evidence.TransitionStatusToComplete(data)` and write the returned bytes
   (keep tmp+rename) so the applied bytes equal the authorized AFTER bytes.
   Stricter than the old loop: a missing status field now errors instead of
   inserting; duplicates error instead of replacing all. The gate call inside
   Complete still runs pre-flip (original state — passes).

Report schema: optional `completion_transition{path, from_status, to_status,
before_sha256, after_sha256, authorized_by}`. No migration; old reports load
with nil and require an exact digest match as before.

## Failure/rollback obligations

Authorized-but-unapplied (Complete vetoed/failed after authorization): the
doc is in the original state, bound == current, gate allows — Complete may be
retried. Partial writes are impossible (tmp+rename). A new record cycle
writes a FRESH report with no transition; stale transitions can never
reconcile (BeforeSHA256 must equal the current binding). One transition per
recorded report: changing an authorization requires re-recording.

## Tests (owner-run; NOT independent)

Env: TMPDIR/GOTMPDIR=$PWD/.parley-runtime/tmp,
GOCACHE=$PWD/.parley-runtime/go-cache, GIT_CEILING_DIRECTORIES=$PWD/.parley-runtime/tmp.

- `go vet ./internal/evidence ./internal/app` → clean.
- `go test -count=1 -p 1 -json ./internal/evidence ./internal/app -run
  'Evidence|Criterion|Envelope|TreeDigest|Typed|Verifier|Barrier|AttestExecution|HelperProcess|ParseGoTestJSON|CappedWriter|RerunCountValidity|InvalidRerunSkipped|Transition|ChecksContract'`
  → exit=0; **90 test-level pass, 0 fail, 1 skip** (stream:
  `.parley-runtime/kimi-completion-transition-20260911.jsonl`). Skip is the
  known host skip `TestTreeDigestUnsupportedEntryFails` (no FIFO/socket on
  this mount) — unchanged, recorded as skip, not pass.
- New tests (all pass): TransitionFrontmatterStatus (15 fail-closed cases +
  byte-exact/invertible positives), AuthorizeCompletionTransition (positive +
  9 refusals + nil), VerifyCompletionTransition (round-trip + 9 negatives + 6
  tampered records), JSON round-trip/old-report compat; gate tests: the fixed
  counterexample (authorized flip closes), unauthorized flip denied,
  wrong-verifier denied, extra-edit denied, unapplied-transition allowed.
- NOT run: full `go test ./...`, any process-level/Codex integration
  (unimplemented). Usage/cost: unknown. Owner-run earns nothing until
  independently reproduced; no signoff, no AC closure.
