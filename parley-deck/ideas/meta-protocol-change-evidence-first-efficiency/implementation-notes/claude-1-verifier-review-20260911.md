---
agent: claude-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-11
artifact-kind: partial independent source review
reviewed-checkpoint: 0a2022b
not-a-signoff: true
phase: implementation-notes
---

# Independent partial source review — Codex production verifier checkpoint 0a2022b

## Scope and standing

This is a NEW partial source review of the production independent-verifier path at
checkpoint `0a2022b`. It is **not** a Phase-6 review artifact, not a full acceptance,
and not a signoff. It is not a republication of the earlier budget review; that review
covered `chargedFixupAttempts`/cap boundaries and is not restated here except where the
verifier path touches it (it does not).

**Files I read directly (PRIMARY, this worktree, at this checkpoint):**

- `parley-deck/ideas/meta-protocol-change-evidence-first-efficiency/FINAL.md` (D3, AC-E1, AC-E2, Idempotence & recovery)
- `parley-deck/ideas/meta-protocol-change-evidence-first-efficiency/implementation-notes/codex-1-verifier-checkpoint-20260911.md`
- `internal/app/evidence_verify.go`
- `internal/app/evidence_verify_test.go`
- `internal/app/driver_evidence.go`
- `internal/app/driver_checks.go`
- `internal/app/driver_impl.go`
- `internal/driver/impl.go`
- `internal/runner/consult.go`
- `internal/evidence/evidence.go`, `internal/evidence/report.go`, `internal/evidence/tree.go`

**I executed nothing.** No shell is available in this launch and none was used. Every
verdict below is from reading the source at the locators given. Where I could not reach
a definition, I say so and mark the item unverified rather than guessing.

**Provenance separation (§15.2).** Everything under "Findings" is `PRIMARY` — located
source, quoted, with the mechanism traced. Everything under "Facilitator-reported
executions (testimony, not verified by me)" is transcribed testimony I explicitly do not
own and do not endorse as established.

## Facilitator-reported executions (testimony, not verified by me)

Transcribed from `codex-1-verifier-checkpoint-20260911.md` lines 32–45 and from the launch
brief. I did not observe any of these runs, did not see the logs, and assign them **no**
verdict. Recorded so a later reader can locate them, not as support for any claim here:

- Full `go test -count=1 ./...` reported PASS on this source before the stronger overlay.
- Scoped vet over app/driver/runner/evidence reported PASS.
- Windows app test-binary cross-build reported PASS; **Windows runtime explicitly untested**.
- A stronger positive-case overlay reported FAIL in 3.224s after a successful completion,
  retained privately at `.parley-runtime/evidence-post-completion-counterexample-20260911.log`.

`TEST-EXECUTION UNVERIFIED` for all of the above from my position. Independently of that:
a passing in-process fixture suite is **not** a live real-model verifier trial. The suite at
`internal/app/evidence_verify_test.go:34` drives a shell script fixture agent
(`evidence_verify_test.go:103-108`), not a configured model CLI, so AC-E2's "an independent
verifier executes the concurrency counterexample and corrected case" is **not** evidenced by
it. `NOVELTY/COVERAGE UNVERIFIED` for any claim that the path is proven end-to-end with a
real model.

---

## Findings

Severity is assigned against FINAL D3 and AC-E1/AC-E2. "Fail-closed in the safe direction"
lowers severity but does not erase a finding: D3's `Idempotence & recovery` paragraph makes
resumability a stated obligation, so a permanent false-deny is a defect, not a nicety.

### [CRITICAL] Deleting the named `checks:` list between driver ticks removes the evidence gate entirely

**Mechanism.** Both enforcement points are conditioned on the *current* shape of the
contract, and nothing durable records that a named contract ever existed.

- `internal/driver/impl.go:257` reads `hadNamedChecks` and `impl.go:287` re-reads
  `hasNamedChecks`; `impl.go:291` requires the verifier only when
  `hadNamedChecks || hasNamedChecks`.
- `internal/app/driver_impl.go:479-486` (`Complete`) calls `EvidenceCloseGate` only in the
  `else if isList` branch.

Both reads happen inside the *same* `Advance` call, so the in-tick deletion the comment at
`impl.go:285-286` defends against ("Removing the original list cannot bypass its independent
evidence gate") is genuinely caught. The gap is **across ticks**.

**Counterexample.** An idea runs several fix-up cycles under a list-form `checks:`, so
`parley-deck/ideas/<slug>/EVIDENCE.json` exists with real records. Before the closing tick,
the `checks:` list is removed from `00-prompt.md` (left absent, not malformed — a malformed
list is caught by `impl.go:258` and `driver_impl.go:229`). On the next tick:
`ReadChecksContract` returns `isList=false, err=nil` at `impl.go:257` and `impl.go:287`, so
`impl.go:291` is false and **`VerifyCompletionEvidence` is never called**; `Complete` at
`driver_impl.go:479` takes the non-list branch and **`EvidenceCloseGate` is never called**;
`driver_impl.go:487-521` writes `status: complete`. The idea closes with zero independent
criterion evidence, while a stale `EVIDENCE.json` sits in the idea directory making it look
evidenced.

**Why this is CRITICAL and not MAJOR.** Under `auto_implement` the implementer/fix-up agent
has repository write access (`driver_impl.go:466-473` re-invokes it on the working tree), and
`00-prompt.md` is an ordinary file in that tree. The invariant the whole slice rests on —
stated verbatim at `driver_impl.go:475-477`, "an implementer cannot short-circuit review" — is
defeated by a one-line edit to a file the implementer can write. This is the removal of the
*entire* original scope, which D3 names explicitly: "a … partial original scope … cannot close
the whole implementation."

**Correction (concrete).** Make the *presence of prior evidence* the trigger, not the current
contract shape:

1. In `driverImplOps.Complete` (`driver_impl.go:479`) and at `impl.go:291`, require the gate
   when `ReadChecksContract` reports a list **OR** `evidence.ReportPath(o.ideaDir)` exists.
2. Have `EvidenceCloseGate` compare the report's recorded criterion names against the current
   contract and deny on shrinkage: a report naming criteria the contract no longer lists is
   exactly "partial original scope", and `driver_evidence.go:74-76`'s current message ("no
   named checks contract — the criterion scope is unknown") is the right denial for it.
3. Additionally persist the original named scope into the driver-owned run cursor on the first
   contract run, so removal is detectable even if `EVIDENCE.json` is also deleted. The cursor is
   driver-authored and lives in the run directory, matching the trust argument already made for
   `chargedFixupAttempts` at `impl.go:533-556`: "A number that is a safety boundary must not be
   authored by the party it constrains." The same reasoning applies to the criterion scope.

Item 1 alone closes the counterexample above and is a small change; items 2–3 harden it.

### [MAJOR] Post-completion gate self-invalidation (reproduced; recorded here with mechanism and a correction that does not silently widen the exclusion)

Already reproduced by Codex and listed as open blocker 1 in the checkpoint note; retained here
per the no-suppression rule, with my own mechanism trace and a specific correction. Kimi is
separately working an exact verifier-bound status transition — **this review assumes no part of
that is done.**

**Mechanism.** `definedEvidenceArtifacts` (`driver_evidence.go:56-59`) excludes
`IMPLEMENTATION.md` *wholesale* from the tested-tree digest. The compensating binding is
`implementationRestDigest` (`driver_checks.go:170-183`), which strips **only** the
`## Validation evidence` section (`driver_checks.go:180`) and hashes everything else —
including the frontmatter. `Complete` (`driver_impl.go:487-521`) rewrites the frontmatter
`status:` line. That rewrite therefore changes `restDigest` while
`report.ExtraDigests[implRel]` still holds the pre-close value, so the very next
`EvidenceCloseGate` denies at `driver_evidence.go:106-108`: "IMPLEMENTATION.md non-evidence
content changed after the evidence was recorded."

**Counterexample.** `Driver.Advance` completes successfully; a fresh `EvidenceCloseGate(drafter)`
immediately afterwards on an unchanged tree returns `Allowed=false`. This matches the overlay
failure described at checkpoint note lines 41-45. The tree digest itself is unaffected
(`IMPLEMENTATION.md` is excluded from it), so the denial is exclusively the `ExtraDigests`
binding.

**Consequences beyond the overlay.** This is not only a test artifact:

- A resumed or re-entered driver run cannot re-verify a completed idea. `FINAL.md`
  "Idempotence & recovery" requires resume "from IMPLEMENTATION.md, current branches and
  durable run state"; after close, that state is self-denying.
- Any later audit that re-runs the gate on a legitimately closed idea gets a denial that is
  indistinguishable from a real scope tamper. The gate's most important message becomes noise.
- It creates pressure toward the one fix that must not be taken — widening the exclusion set —
  which the source itself warns against at `evidence_verify.go:351-353` ("Never expand
  exclusions to hide it") and which the checkpoint note correctly rules out.

**Correction (concrete), in preference order.**

1. **Preferred — bind the transition, do not hide it.** Extend the report with a second,
   explicitly-named binding for the post-close state: before writing `status: complete`,
   `Complete` computes the restDigest of the file *as it will be after* the status rewrite and
   records it as a distinct field (e.g. `ExtraDigests[implRel+"#closed"]`, or a typed
   `ClosingDigests map[string]string`), then performs the write. `EvidenceCloseGate` accepts a
   match against **either** the pre-close or the recorded post-close binding, and against
   nothing else. Both values are driver-computed, both are persisted, and any edit that is not
   exactly the authorized status transition still fails. This is the "exact independently
   authorized before/after transition … designed and reviewed explicitly" that checkpoint
   blocker 1 asks for, and it adds no silent exclusion.
   - Ordering matters: the report re-save must happen *before* the `IMPLEMENTATION.md` write,
     so a crash between them leaves a report whose post-close binding simply never matches —
     fail-closed — rather than a closed file with no binding.
2. **Acceptable fallback — normalized status line.** Have `implementationRestDigest` canonicalize
   the single frontmatter `status:` value to a fixed placeholder before hashing. This is narrow
   (one key, first frontmatter block only) and the value is independently validated by
   `ImplementationStatus` (`driver_impl.go:198-204`) and the state machine at `impl.go:102-114`,
   so it is not a usable hiding place. **But** it is still a widening of what the digest ignores,
   it must be documented in `driver_checks.go` and in the report schema rather than applied
   quietly, and it does not generalize if any other driver-owned frontmatter field is later
   rewritten. I recommend it only if (1) is judged too invasive for this slice.
3. **Reject** — broad frontmatter exclusion, or excluding `IMPLEMENTATION.md` from
   `ExtraDigests` altogether. Either re-opens the hidden-scope-edit hole that `ExtraDigests`
   exists to close (`evidence.go:112-118`, `driver_checks.go:126-132`).

### [MAJOR] The retained `VerifierExecution` tree digests are copied from the request, never measured around the re-run

**Mechanism.** `internal/evidence/evidence.go:60-64` documents the contract:
"TreeBeforeSHA256/TreeAfterSHA256 are the tested-tree digests the verifier computed
**immediately before and after its re-run**". The only production producer is
`evidence_verify.go:265`:

```go
rerun := evidence.VerifierExecution{Command: rec.Command, TreeBeforeSHA256: req.TreeSHA256, TreeAfterSHA256: req.TreeSHA256}
```

Both fields are the same constant, copied out of the frozen request. They are not measured.

**Counterexample.** `reconcileRerun`'s stability check at `evidence.go:404-406` —
"tested tree changed during the independent re-run" — is structurally incapable of firing
against this producer: `vr.TreeBeforeSHA256 != vr.TreeAfterSHA256` is `req.TreeSHA256 !=
req.TreeSHA256`, always false. Likewise `evidence.go:407-409` degenerates to comparing the
request digest against the report digest, which `verificationBindings` already established at
`evidence_verify.go:137`. Two of the three retained tree assertions are tautologies.

**What actually provides the property today.** The real stability guard is the
`verificationBindings` call bracketing each criterion at `evidence_verify.go:254` and
`evidence_verify.go:259`, each of which recomputes `evidence.TreeDigest` (`evidence_verify.go:167`)
and compares to `req.TreeSHA256`. So the *behaviour* is presently correct — this is not a live
false-accept. The defect is that the persisted artifact records an assertion it did not make,
and the defence-in-depth layer that is supposed to catch a regression in the surrounding code
cannot catch one.

**Why it matters concretely.** `verificationBindings` is O(whole tree) and is invoked `2N+2`
times per verification (`evidence_verify.go:254`, `:259`, `:273`, plus `:241`). That is exactly
the shape a future "this is redundant, hoist it out of the loop" refactor removes. If the two
in-loop calls are dropped, tree stability during execution silently stops being checked, while
`EvidenceCloseGate` keeps reporting — from the persisted record — that before == after == the
tested tree. A later auditor reading `EVIDENCE.json` cannot tell the difference.

**Correction.** In `evidence_verify.go`, measure the digests where the doc comment says they are
measured:

```go
before, err := evidence.TreeDigest(root, exclusions...)   // immediately before RunCriterion
rec := evidence.RunCriterion(ctx, root, criterion.Name, criterion.Command, req.Verifier)
after, err := evidence.TreeDigest(root, exclusions...)    // immediately after
rerun := evidence.VerifierExecution{Command: rec.Command, TreeBeforeSHA256: before, TreeAfterSHA256: after}
```

Keep both `verificationBindings` calls as the outer binding check. The exclusion set is already
computed inside `verificationBindings` via `definedEvidenceArtifacts` (`evidence_verify.go:163`);
hoist it into `executeIndependentVerification` once and reuse it so the added digests cost
nothing beyond the two measurements. After this change `evidence.go:404-409` becomes a real
check, and the doc comment at `evidence.go:60-64` becomes true.

### [MAJOR] A receipt-write failure after a successful `Save` leaves the idea permanently unclosable

**Mechanism.** The helper attests into the in-memory report and persists it at
`evidence_verify.go:282` (`evidence.Save`), then reads it back and sets `receipt.UpdatedSHA256`
at `:285-289`. The receipt itself is written only by the deferred closure at
`evidence_verify.go:222-227`. If that deferred write fails — read-only runtime directory, full
disk, or the receipt path removed — `errors.Join` makes the helper return an error and exit 1,
while `EVIDENCE.json` on disk **already carries the verifier provenance**.

On the parent side, `evidence_verify.go:376-378` then reports "no valid helper execution
receipt" and denies. On any retry, the helper reaches
`evidence_verify.go:245-249` and refuses outright:

```go
if rec.Provenance.Verifier != "" || rec.Provenance.VerifierRerun != nil {
    return errors.New("fresh execution report required; an earlier verifier cannot be replayed")
}
```

**Counterexample.** Make `.parley-runtime/evidence-verification/attempt-*/` unwritable between
the helper's `Save` and its deferred receipt write (or exhaust the disk at that instant). The
canonical report is attested; every subsequent `VerifyCompletionEvidence` fails at the replay
guard; there is no in-band recovery. Hand-editing `EVIDENCE.json` to strip the provenance is
then rejected by `verificationBindings` (`evidence_verify.go:237-240`, "canonical evidence
differs from the frozen original report") on the next attempt whose frozen original was taken
before the edit, and by the parent's `reflect.DeepEqual` at `:405-407` after it. The idea is
terminal. This is the failure class Codex's blocker 3 names ("report and receipt persistence
failures") — I am recording the exact mechanism and a correction, not re-raising it as new.

**Correction.** Recovery already has everything it needs and is not being used:

1. The driver retains byte-exact pre-attestation bytes at
   `<attempt>/original-report.json` (`evidence_verify.go:340`) and their digest in
   `req.ReportSHA256`. On detecting an orphaned attestation — persisted verifier provenance with
   no corresponding valid receipt — the **driver** (never the verifier) restores the report from
   those retained bytes after confirming their digest, records the restoration, and re-runs the
   verification from a clean state. This keeps "fresh execution report required" fully intact;
   it does not weaken the replay guard at all.
2. Independently, relax the receipt's `UpdatedSHA256` from *required* to *if-present-must-match*
   at `evidence_verify.go:390`. The per-record reconciliation at `:398-404` (name, pass status,
   `reflect.DeepEqual` on the retained command, executor identity) plus the whole-report
   `DeepEqual` at `:405-407` already establish everything `UpdatedSHA256` establishes. Removing
   the hard dependency on the *last* write in the helper shrinks the terminal window.
3. Do **not** "fix" this by scoping the replay guard to a different verifier identity. That
   would make an already-attested report re-attestable by the same identity, which is the
   property the guard exists to deny.

### [MAJOR] A concurrent writer landing between the pre-save byte comparison and `Save` is silently clobbered

Acknowledged in-source at `evidence_verify.go:276-277` and as checkpoint blocker 2; Codex makes
no concurrency-safety claim for this checkpoint. Recorded with the specific window and a
correction, since the launch brief asks for the concurrent-writer and race boundaries.

**Mechanism.** `evidence_verify.go:278-281` reads the canonical report and compares bytes to the
frozen original; `evidence_verify.go:282` then calls `evidence.Save`, which is internally atomic
(temp + `SyncFile` + `Rename`, `report.go:41-68`) but is a *separate* operation from the
comparison. Nothing holds a lock across the two.

**Counterexample (the one the parent does not catch).** A second process — a concurrent
`runChecksContract` cycle (`driver_checks.go:101`, `:155`) or a second driver run on the same
worktree — completes a check cycle on the *same unchanged tree* and writes a **new**
`EVIDENCE.json` in which one criterion is now `fail` (a genuinely flaky test, or a criterion that
regressed). That write lands strictly between the helper's read at `:278` and its rename at
`:282`. The helper's rename then replaces the fresh failing report with the frozen original plus
its attestations. The parent's detection at `:386-392` compares against
`receipt.UpdatedSHA256`, which is the helper's *own* bytes, so the comparison succeeds; the
per-record checks at `:398-404` succeed (the report is the original, correctly attested); the
tree digest is unchanged, so `EvidenceCloseGate` passes. **A concurrently observed failing
result is destroyed and the close proceeds.**

Note the asymmetry: a concurrent write landing *after* the helper's `Save` **is** caught, at
`:390` (`receipt.UpdatedSHA256 != sha256Hex(updatedRaw)`) and again at `:398-404` (a foreign
report carries no `VerifierRerun`). Only the compare→rename window is unguarded.

**Correction.**

1. Serialize report writers with an advisory lock file created `O_CREATE|O_EXCL` in the idea
   directory (e.g. `.EVIDENCE.lock`, carrying pid + run id + timestamp for stale-lock
   diagnosis), acquired by **both** `writeTypedEvidence` (`driver_checks.go:133-163`) and the
   helper across compare+`Save` (`evidence_verify.go:278-284`). Failure to acquire is a denial,
   never a wait-forever and never a silent override.
2. Give `evidence.Save` an optional compare-and-swap form — `SaveIfUnchanged(ideaDir, r,
   expectedSHA)` — that re-reads and re-compares immediately before the rename. This narrows the
   window to the rename itself without fully closing it; it is a useful belt to the lock's
   braces, not a substitute for it.
3. Have the helper verify its own write: `evidence_verify.go:285-289` currently reads the file
   back only to hash it for the receipt. It should also compare the read-back bytes to the bytes
   it marshalled and fail if they differ, which detects a racing writer that lands immediately
   after the rename.

### [MAJOR] The verifier command is POSIX-quoted and handed to the agent as a shell string; Windows is unsupported and untested

**Mechanism.** `evidence_verify.go:363-364` builds the command with a POSIX single-quote quoter:

```go
quote := func(s string) string { return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'" }
command := quote(executable) + " evidence verify --request " + quote(requestPath) + " --request-sha256 " + quote(requestSHA)
```

and the prompt at `:365-368` instructs the agent to "Execute the exact verifier command below".
On Windows `cmd.exe` does not treat `'` as a quoting character at all, and the paths involved
contain backslashes and may contain spaces — the test itself deliberately uses a spaced path
(`filepath.Join(t.TempDir(), "parley fixture")`, `evidence_verify_test.go:42`), so the quoting is
load-bearing, not cosmetic.

**Counterexample.** On Windows, `executable` is `C:\Users\...\parley.exe`; the emitted string is
`'C:\Users\...\parley.exe' evidence verify --request '...'`. A `cmd.exe`-hosted agent attempts to
run a program literally named `'C:\Users\...\parley.exe'` and fails. The failure surfaces as
`evidence_verify.go:372-374` "independent verifier process failed or its launch was not
observed" — a correct fail-closed outcome, but one that is indistinguishable from a real verifier
failure and that makes the named-contract close path unusable on that platform.

The test acknowledges the gap by skipping (`evidence_verify_test.go:35-37`: "POSIX fixture agent;
Windows runtime requires its own execution host"), and the checkpoint note states the cross-build
passed but "runtime not tested". So the platform is knowingly unvalidated; the finding is that the
*code* has a POSIX assumption baked into a string that is handed to a model, with no guard.

**Correction.** Either (a) emit the command through a platform-aware quoter selected on
`runtime.GOOS`, or better (b) stop shipping a shell string: pass the request path and digest as
plain prompt fields and have the agent invoke `parley evidence verify` with them as arguments,
so no quoting layer exists. Until one of those lands, add an explicit early denial in
`VerifyCompletionEvidence` on `runtime.GOOS == "windows"` with a precise reason, so the failure
is legible rather than masquerading as a verifier failure. A legible refusal is honest; a
mis-quoted command that fails for an unrelated-looking reason is not.

### [MAJOR] A refused independent verification leaves no durable, committed record

**Mechanism.** On refusal, the helper's receipt is written to
`<root>/.parley-runtime/evidence-verification/attempt-*/result.json`
(`evidence_verify.go:226`) and the agent logs to the same directory
(`evidence_verify.go:371`). That whole tree is runtime-local and expected to be
git-ignored — the fixture writes exactly that ignore rule at
`evidence_verify_test.go:73`. `commitEvidence` (`driver_checks.go:223-245`) commits only
`IMPLEMENTATION.md` and `EVIDENCE.json`, and is not called on this path at all. The only
operator-visible trace is the escalation string at `impl.go:296-298`, which lives in process
output.

**Why it is a finding.** `FINAL.md` "Idempotence & recovery" requires "Preserve failed
invocation records", and D2 requires terminal records for real outcomes. A refused independent
verification is precisely the event that most needs to survive: it is the difference between
"nobody tried" and "an independent verifier tried and the evidence did not hold". Today, after
a restart or a log rotation, those two states are indistinguishable in the canonical record.

**Correction.** Persist a bounded, scrubbed attempt summary into a committed artifact —
verifier id, invocation id, request digest, per-criterion statuses from `receipt.Executions`,
and the refusal reason — reusing the existing scrubber `scrubAndTruncate`
(`driver_checks.go:197-217`) so no raw output leaks. Two viable shapes: an append-only
`## Independent verification attempts` section in `IMPLEMENTATION.md`, or a sibling
`EVIDENCE-ATTEMPTS.json`.

**Trade to state explicitly, not to hide:** either shape adds content to a digest-bound path.
`IMPLEMENTATION.md` is already wholesale-excluded from the tree digest but bound by
`implementationRestDigest`, so an append-only section there **must** be stripped by
`implementationRestDigest` exactly as `## Validation evidence` is (`driver_checks.go:180`), or
every append will invalidate the binding — the same class of problem as the MAJOR above. A new
`EVIDENCE-ATTEMPTS.json` must be added to `definedEvidenceArtifacts` (`driver_evidence.go:56-59`)
**and** given its own binding, because an excluded-but-unbound path is a hiding place. This is
the exact hazard `evidence_verify.go:351-353` warns about; any new evidence artifact inherits the
obligation. I recommend `EVIDENCE-ATTEMPTS.json` with an `ExtraDigests` entry, since it keeps the
attempt log out of the human-authored document entirely.

### [MAJOR] Failed runtime-directory creation poisons the tree permanently when `.parley-runtime/` is not ignored

**Mechanism and ordering.** `VerifyCompletionEvidence` validates bindings once at
`evidence_verify.go:328` (tree still clean), then creates its runtime directory and writes two
files into the tree:

- `evidence_verify.go:331` `os.MkdirAll(base)` — an empty directory, invisible to both
  `git ls-files` and the walk in `listTreeFiles` (`tree.go:162-203`), so harmless;
- `evidence_verify.go:340` writes `original-report.json`;
- `evidence_verify.go:344` writes `request.json`.

It then re-checks bindings at `evidence_verify.go:353`. If `.parley-runtime/` is not ignored,
those two files are untracked-and-not-ignored, so `git ls-files -c -o --exclude-standard`
(`tree.go:164`) includes them, `TreeDigest` changes, and the check correctly denies.

**Counterexample.** Run a named-contract close in a deck whose `.gitignore` lacks
`.parley-runtime/`. The attempt denies with "tested tree changed before or during independent
verification" (`evidence_verify.go:177`). **Nothing removes the directory** — there is no
cleanup on any failure path in `VerifyCompletionEvidence`. The two files stay in the tree, so
the tree digest now permanently differs from `report.TreeSHA256`. Every subsequent attempt
fails, and each one adds another `attempt-*` directory, compounding the divergence. Adding the
ignore rule afterwards does **not** recover: the already-created files remain untracked and the
newly-added `.gitignore` line is itself a tracked-tree change. Only manual deletion of
`.parley-runtime/` plus a fresh evidence cycle recovers, and the error message never says so.

**Correction.**

1. Probe before writing: check that the runtime base path is out of digest scope *before*
   creating anything — write a single probe file, recompute `TreeDigest`, delete the probe, and
   if the digest moved, refuse with an actionable message naming the exact remedy ("add
   `.parley-runtime/` to .gitignore"). A `git check-ignore -q .parley-runtime` probe is cheaper
   and sufficient inside a git work tree; keep the digest probe as the non-git fallback.
2. Clean up on failure: `defer` removal of the `attempt-*` directory when the verification did
   not reach a state worth retaining, or — better, since the agent logs are genuinely useful —
   retain it **only** once the ignore precondition is known to hold, which item 1 establishes.
3. Distinguish the two causes in the message. "The driver's own runtime files changed the tested
   tree because `.parley-runtime/` is not ignored" and "the code under test changed during
   verification" are different operator actions; today both print the same string from
   `evidence_verify.go:177`.

### [MINOR] Receipt-to-record matching is positional, and depends on an unenforced ordering

**Mechanism.** The helper appends executions in `criteria` order
(`evidence_verify.go:250-258`), while the parent compares them **by index** against
`updated.Records` (`evidence_verify.go:398-404`). Ordering coincidence is what makes this work:
`runChecksContract` iterates `criteria` when building records (`driver_checks.go:79-97`), so
`report.Records` happens to share the contract's order. Nothing enforces it —
`verificationBindings` (`evidence_verify.go:147-162`) establishes only a name-keyed bijection.

**Counterexample.** Reorder the two entries in the `checks:` list between the recording cycle
and the closing tick, leaving both names and commands intact. `report.Records` is `[A, B]`;
`receipt.Executions` is `[B, A]`; `AttestExecution` mutates in place by name so `updated.Records`
stays `[A, B]`; the parent's `receipt.Executions[0].Name != updated.Records[0].Name` fires and
denies with "retained attestation differs from actual helper execution". Fail-closed, so this is
MINOR — but the message blames the attestation for what is a contract reordering, which will cost
an operator real time.

**Correction.** Key the comparison by criterion name — build
`map[string]evidence.CriterionRecord` from `receipt.Executions`, look up each
`updated.Records[i].Name`, and deny if any name is absent or duplicated. That is the same shape
`verificationBindings` already uses at `evidence_verify.go:147-156` and it removes the hidden
ordering dependency entirely.

### [MINOR] `reflect.DeepEqual` over the whole `Report` compares `time.Time` structurally

**Mechanism.** `evidence_verify.go:405-407` guards the central "verifier changed the original
evidence instead of only attesting it" property with
`reflect.DeepEqual(updated, original)`. `Report.GeneratedAt` is a `time.Time`
(`evidence.go:110`), whose struct comparison (`wall`/`ext`/`loc`) is documented-fragile.

Both values here are decoded from JSON produced by the same `json.MarshalIndent` path
(`report.go:36`), so monotonic readings are stripped and the location is UTC in both — the
comparison is correct today, and I am not claiming a live defect. The finding is that the
strongest integrity check in the parent rests on a comparison Go's own documentation warns
against for this type, and it would break silently under an unrelated change (a `Location`
difference, or a future `Report` field with a `time.Time` or a map with differing nil-vs-empty
representation).

**Correction.** Compare canonical serializations rather than structs: after the provenance reset
at `evidence_verify.go:403`, `json.Marshal` both and compare bytes, or hash both and compare
digests. That is exact, stable under refactoring, and produces a diffable artifact for the error
message. Alternatively, keep `DeepEqual` but compare `GeneratedAt` separately with `.Equal()`.

### [MINOR] `verificationBindings` re-hashes the entire tree `2N+2` times per verification

**Mechanism.** Each call runs `driver.ReadChecksContract`, `definedEvidenceArtifacts`,
`evidence.TreeDigest` (a full read-and-hash of every in-scope file, `tree.go:60-124`) and
`implementationRestDigest`. It is invoked at `evidence_verify.go:241`, then at `:254` and `:259`
inside the per-criterion loop, then at `:273` — `2N+2` full tree hashes for `N` criteria, on top
of the criteria themselves. On a repository of this size that is a real cost, and it is the
motivation a future contributor will have for the loop-hoisting refactor that Finding
"[MAJOR] retained tree digests are copied from the request" shows would silently remove the
stability property.

**Correction.** Hoist the invariant parts — the exclusion list and the contract read — out of the
loop, and keep only the tree and rest digests inside it. Combined with the measured before/after
digests recommended above, the per-criterion cost becomes two tree hashes instead of two full
binding evaluations, the behaviour strengthens rather than weakens, and the incentive to delete
the in-loop checks goes away.

### [NIT] `driver_evidence.go`'s wiring comment no longer describes the real call sites

`driver_evidence.go:24-38` presents the gate as un-wired ("Wiring (Codex-owned driver state
machine; this file is the API, the call site is theirs)") and gives an illustrative snippet. The
gate now has two real call sites — `driver_impl.go:482` inside `Complete` and
`evidence_verify.go:408` at the end of `VerifyCompletionEvidence`. A reader auditing where the
gate is enforced is pointed at a hypothetical instead of at them.

**Correction.** Replace the snippet with the two actual locations and state that `Complete`'s
call is the last check before the status write. This matters more than typical comment drift,
because the CRITICAL above turns on exactly which call sites are conditioned on the contract
shape.

### [NIT] Inconsistent path-escape predicates between the two modules

`driver_evidence.go:53` checks `rel == ".." || strings.HasPrefix(rel, ".."+sep) || filepath.IsAbs(rel)`.
`evidence_verify.go:210` omits the exact `rel == ".."` case. In `evidence_verify.go` the omission
is covered by the adjacent `filepath.Base(rel) != "request.json"` test, so there is no hole — but
two spellings of the same security predicate in one slice is how one of them eventually gets
copied without its companion guard. Extract a single `isContained(root, path) bool` helper and
use it in both.

---

## Refutation attempts that did NOT break the implementation

Recorded per LE-1: a "no findings here" statement is only credible with the attempts shown.
These are traced attacks against the source, not executions.

- **Textual PASS substitution.** The fixture case `text-pass-without-helper`
  (`evidence_verify_test.go:53`) mirrors what I traced: printing `GOAL-CHECK: PASS` without
  invoking the helper leaves no `result.json`, so `evidence_verify.go:376-378` denies. The
  goal-check verdict path (`driver_impl.go:410-417`) is separately gated and runs *before*
  `VerifyCompletionEvidence` (`impl.go:280-284` then `:296`), so a textual pass cannot reach the
  close. Consistent with D3 and LE-7. **Could not break.**
- **Self-verification.** Blocked at four independent layers: `evidence_verify.go:295`
  (`o.drafter == o.implementer`), `evidence_verify.go:153` (executor equals requesting verifier),
  `evidence.go:332-334` (`AttestExecution` refuses a self-executed criterion), and
  `evidence.go:228-244` (`Evaluate`'s two self-verdict reasons). **Could not break.**
- **Hand-populated provenance.** `Provenance.Verifier` without a `VerifierRerun` is denied at
  `evidence.go:246-248`; a fabricated `VerifierRerun` must match the criterion's exact command
  hash, carry a non-empty output hash, exit 0, use a case-proving format, and reproduce the
  executed/skipped/failed-package counts (`evidence.go:362-423`). **Could not break within the
  stated trust boundary.** Outside it — a same-UID process able to fabricate mutually consistent
  hashes — the package documents the limit honestly at `evidence.go:163-169` and I agree the
  framing is correct: this is attribution, **explicitly not adversarial authentication**, and the
  code says so rather than marketing it as proof.
- **Masked package failure.** The fixture pipes through `cat` to mask `go test`'s exit code
  (`evidence_verify_test.go:68`) and `TestMain` exits 1 (`:87-91`). `FailedPackages > 0` denies at
  `evidence.go:274-276` and `evidence.go:377-379` even though a real test case passed and the exit
  code is 0. This is a genuinely good check and it is the kind of thing a regex over output would
  have missed — it matches D3's "Unknown command output is not semantically certified by a text
  regex." **Could not break.**
- **Skipped-only execution.** `evidence.go:211-212` and `evidence.go:259-260` both deny.
  **Could not break.**
- **Stale tree.** Recomputed at close (`driver_evidence.go:89`) and at four points during
  verification; `evidence.go:189-191` denies on mismatch. **Could not break.**
- **Partial scope via a shrunken report.** `evidence_verify.go:144` requires
  `len(criteria) == len(report.Records)`, and `evidence.go:201-205` requires a record per required
  criterion. **Could not break** — note this is the mirror image of the CRITICAL above, which
  shrinks the *contract* rather than the *report*, and is not covered.
- **Request tampering.** The request path must resolve inside the driver-owned runtime directory,
  be named `request.json`, sit one level below the base, and match the digest the driver passed
  on the command line (`evidence_verify.go:190`, `:201-211`). The root must be canonical and
  equal to `req.Root` (`:193-200`). **Could not break.**
- **Unattributed helper invocation.** `evidence_verify.go:216-218` requires `PARLEY_RUN_ID`,
  `PARLEY_AGENT_ID` and a non-empty `PARLEY_PROC_MARKER` matching the request, and
  `TestEvidenceHelperRefusesMissingRuntimeIdentity` (`evidence_verify_test.go:174-203`) asserts no
  report and no receipt are written on refusal. The refusal ordering is correct: the identity check
  at `:216` precedes the deferred receipt installation at `:222`, so an unattributed call writes
  nothing at all. **Could not break** — subject to the open verification item below.

## Open verification items I could not close from source

Stated as unverified rather than asserted either way (§15.2 — `RECALL` support does not carry a
verdict):

1. **The actual invocation marker binding.** `evidence_verify.go:379` requires
   `receipt.ProcessMarker == res.InvocationID`, where the receipt value is the child's
   `PARLEY_PROC_MARKER` (`evidence_verify.go:221`) and `res.InvocationID` is captured through the
   telemetry observer at `consult.go:112-117`. Whether `execAgentProcess` exports
   `PARLEY_PROC_MARKER` with exactly that invocation id is **the load-bearing fact** for this
   whole binding, and I could not locate `execAgentProcess` — it is not at
   `internal/runner/exec.go`. `MARKER-BINDING UNVERIFIED`. A reviewer with a shell should confirm
   (a) the env var is set from the same value the observer reports, and (b) that
   `consult.go:112-117` overwrites `result.InvocationID` on every observed record, so a
   requested-then-terminal pair with differing ids cannot leave the wrong one in place. If (b)
   does not hold, the equality check at `:379` compares the *last* record's id against the
   child's, which may be a different value.
2. **`ReadChecksContract` semantics.** I rely on `isList=false, err=nil` for an absent `checks:`
   key, which is what the CRITICAL counterexample turns on. This is strongly implied by
   `driver_impl.go:229-240` (the non-list path falls through to the scalar/`go test` resolution)
   but I did not read `ReadChecksContract` itself. `CONTRACT-ABSENCE-SEMANTICS UNVERIFIED`. If an
   absent key instead returns an error, the CRITICAL's severity drops — the deletion would
   escalate at `impl.go:258` rather than close. **This is worth checking first**, before acting on
   the CRITICAL, because it is a five-minute read that determines the severity.
3. **`runner.RunCriterion`.** I did not locate its definition (not at `internal/evidence/run.go`,
   `exec.go`, or `criterion.go`). Its format detection, case counting, and output hashing are what
   `reconcileRerun` reconciles against, so its correctness is assumed here, not verified.
   `RUNCRITERION UNVERIFIED`.

## Recommendations, in the order I would act

1. **Read `ReadChecksContract`** to settle open item 2. It determines whether item 2 below is a
   CRITICAL or a MAJOR, and it is the cheapest thing on this list.
2. **Gate on evidence presence, not contract shape** (CRITICAL). `driver_impl.go:479` and
   `impl.go:291` should require the gate when `EVIDENCE.json` exists, regardless of the current
   `checks:` shape, and `EvidenceCloseGate` should deny on scope shrinkage. This is the only
   finding here that permits a false close.
3. **Bind the status transition explicitly** (MAJOR, reproduced). Prefer the recorded post-close
   digest over any exclusion widening. Coordinate with Kimi's verifier-bound status transition
   work rather than landing a second mechanism — but **do not treat that work as done**; nothing
   in this checkpoint corrects it, and the overlay failure stands at `0a2022b`.
4. **Measure the re-run tree digests** (MAJOR). Small, local, turns two tautologies into real
   checks, and removes the refactor hazard that the `2N+2` hashing cost creates.
5. **Serialize report writers** (MAJOR) and **make an orphaned attestation recoverable from the
   retained original bytes** (MAJOR). Both are prerequisites for any concurrency claim; the
   checkpoint correctly makes none.
6. **Persist refused verifications durably** (MAJOR), with the new artifact's digest binding
   designed in from the start rather than added later.
7. **Probe-and-clean the runtime directory** (MAJOR) and **refuse legibly on Windows** (MAJOR)
   before anyone runs this outside a POSIX host with `.parley-runtime/` already ignored.
8. Then the MINORs and NITs.

## Assessment

The architecture of this path is sound and the fail-closed discipline is real, not decorative:
the layering of receipt binding, `AttestExecution` reconciliation, and `Evaluate`'s independent
re-check means that most of the attacks I traced are stopped at two or three places rather than
one. The masked-package-failure and skipped-test denials in particular defend properties that a
text-matching gate would have missed, which is exactly what D3 demands.

The one finding that permits a **false close** is the CRITICAL — removing the named contract
across ticks bypasses both enforcement points, because both ask what the contract *is* rather
than what the evidence *says it was*. Every other MAJOR here fails in the safe direction; they
are defects because D3 also obliges the system to be resumable, auditable, and honest about what
it recorded, and a permanently self-denying gate, a tautological retained assertion, a terminal
unclosable state, and a destroyed concurrent failure each violate one of those.

Nothing here should be read as acceptance of the slice. AC-E2 requires an independent verifier to
execute the concurrency counterexample and the corrected case; the process fixtures do not supply
that, no live real-model verifier trial has been evidenced to me, and I ran nothing. My position
is unchanged from the design signoff: independent verification remains an obligation, not a
passed gate.
