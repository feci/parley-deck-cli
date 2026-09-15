---
agent: claude-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-11
artifact-kind: partial independent source review
reviewed-checkpoint: 0cdc56e
not-a-signoff: true
phase: implementation-notes
---

# Independent partial source review — serialized evidence publication at 0cdc56e

## Summary

This is a NEW partial source review of checkpoint `0cdc56e`, focused on evidence report
serialization/CAS, the independent-verifier receipt, persisted driver acceptance, and the
completion transition across run/restart and concurrent writers. It is **not** a Phase-6
review artifact, not acceptance, and not a signoff.

Eight of the ten findings from my `0a2022b` review are substantively corrected in source,
and the corrections are real rather than cosmetic: the completion transition
(`internal/evidence/completion.go`) binds the status flip by independent recomputation
instead of widening the digest exclusion, the publication guard plus
`SaveIfUnchanged` closes the compare→rename window for cooperating writers, and the
runtime-ignore prerequisite now refuses before any artifact is created. Two dispositions do
not hold as stated, and I found three new defects, one of which still permits a false close.

**I executed nothing.** No shell is available in this launch and none was used. Every verdict
is from reading the source at the locators given (`PRIMARY`). Codex-reported test executions
are transcribed as testimony under their own heading and carry no verdict from me.

---

## Actual coverage

**Read in full at this checkpoint (PRIMARY):**

- `internal/evidence/report.go`, `internal/evidence/report_guard.go`, `internal/evidence/completion.go`
- `internal/budget/resource_guard.go`, `internal/budget/lock.go`, `internal/budget/lock_unix.go`
- `internal/app/evidence_verify.go`, `internal/app/driver_checks.go`, `internal/app/driver_evidence.go`, `internal/app/driver_impl.go`
- `internal/driver/checks.go`, `internal/driver/cursor.go`
- `FINAL.md` (D2/D3/D6/D7, AC-E1/AC-E2, Idempotence & recovery), my `claude-1-verifier-review-20260911.md`,
  `codex-1-verifier-readiness-budget-dispositions-20260911.md`, `codex-1-report-publication-20260911.md`

**In the requested scope and NOT read — coverage gaps I am not papering over:**

- `internal/budget/lock_windows.go` — not read. Every Windows statement below is about the
  POSIX file plus the explicit `runtime.GOOS` refusal, not about Windows lock semantics.
- `internal/evidence/report_guard_test.go`, `internal/app/evidence_publication_test.go`,
  `internal/app/evidence_verify_test.go` — **not read.** I therefore assert nothing about
  what the new tests do or do not cover, and no finding below is qualified by "a test
  catches this".
- `internal/driver/impl.go` — **not read at this checkpoint.** This is the most consequential
  gap: the CRITICAL below concerns `driverImplOps.Complete`, and whether `driver.Advance`
  independently blocks that path is unverified. Flagged inline.
- `internal/evidence/evidence.go`, `tree.go`, `run.go` — not re-read; I rely on my `0a2022b`
  reading of `evidence.go:60-64`, `:404-409`, `Evaluate`, `AttestExecution`, and
  `tree.go:162-203`. Marked `CARRIED-FORWARD` where load-bearing.
- The previous `3c6da49` readiness corrections (`preflight_liveness.go`,
  `preflight_evidence_test.go`, `preflight.go`) — **not inspected at all.** Budget and the
  15-minute deadline went to the primary scope. I make **no** statement about D7/AC-L1 at
  this checkpoint, and the readiness section of
  `codex-1-verifier-readiness-budget-dispositions-20260911.md` is **unevaluated** by me.
- `internal/app/preflight*.go` process-budget wiring — not inspected; the D6 disposition and
  its stated LIMITATION are likewise unevaluated.

---

## Refutation attempts

Per LE-1, traced attacks on the source, not executions.

**Could not break:**

- **Concurrent clobber of the compare→rename window** (my prior MAJOR). `SaveIfUnchanged`
  (`report_guard.go:126-140`) now performs the read, the byte comparison and the replacement
  inside one guard-held transaction, and the helper's CAS and its receipt write are in the
  *same* `WithReportWriter` block (`evidence_verify.go:290-309`), so a cooperating writer
  cannot land between them. `w.active` (`report_guard.go:129`, `:143`) prevents use of the
  writer after the transaction ends. For cooperating same-host writers this is closed.
- **Lock that a filesystem implements as a no-op.** `lock.go:121-144` opens a second
  descriptor on the same inode and requires `take(probe)` to *fail*; a shared mount that
  reports success without exclusion is rejected rather than trusted. This is the right check
  and it is the one most systems omit.
- **Status-flip laundering.** `VerifyCompletionTransition` (`completion.go:249-260`) does not
  trust any recorded digest: it requires `sha256Hex(current) == AfterSHA256`, then
  *independently inverts* the current content back to `FromStatus` and requires the result to
  hash to `BeforeSHA256`, and separately requires the current status to actually be
  `complete`. Altering any recorded field fails at least one recomputation. Using the
  transition to launder a second edit fails the inversion.
- **Authorizing a transition on a path other than the implementation doc.**
  `AuthorizeCompletionTransition` takes a caller-supplied `path`, but both gates pin it:
  `driver_evidence.go:109` passes `implRel`, `driver_impl.go:532` passes `rel` derived from
  the IMPLEMENTATION.md path, and `VerifyCompletionTransition` rejects `tr.Path != path`.
- **Self-authorized transition.** `completion.go:186-191` requires every record to be
  attested by the authorizing verifier *and* rejects `Executor == verifier`.
- **A second transition on one report.** `completion.go:167-169` refuses, and
  `evidence_verify.go:348-350` refuses to start a new verification against a report that
  already carries one.
- **Replaying an attested report.** `evidence_verify.go:250-254` refuses any record already
  carrying `Verifier`/`VerifierRerun`.
- **Short-digest panic on the `bound[:12]` diagnostic** (`completion.go:180`). Guarded by
  `validCompletionDigest(bound)` at `:171`, which requires 32 decoded bytes and lowercase
  hex. Codex's reported probe matches what the source now does.
- **Duplicate / quoted / merge-key YAML status.** `completion.go:107-131` checks YAML
  *meaning* in addition to normalized bytes: `<<` merge keys refuse, a non-`!!str` or
  mismatched node refuses, and `statusCount != 1` refuses.
- **Substituting another run's acceptance at completion.**
  `requireAcceptedVerification` (`evidence_verify.go:513-549`) binds `RunID` and `Verifier`,
  re-hashes the report against `accepted.ReportSHA256`, re-reads and re-hashes the request
  and the receipt, and requires `receipt.UpdatedSHA256 == accepted.ReportSHA256`. Replacing a
  still-valid report after acceptance, or removing the receipt, denies. This matches what
  Codex reports and it is what the source does.
- **Contract removal across ticks, pinned case.** `ObserveChecksContract`
  (`driver/checks.go:48-50`) rejects any current digest that differs from the cursor pin, and
  `:51-71` activates from a prior `EVIDENCE.json` when no pin exists, additionally requiring
  the current list to match the recorded scope. My prior CRITICAL is closed for both of those
  cases. The residual case is the CRITICAL below.

**Broke, or found unsupported:** the six findings that follow.

---

## Findings

### [CRITICAL] Completion's evidence gate is activated only by artifacts the implementer can delete; deleting the run cursor, the `checks:` list and `EVIDENCE.json` closes with no gate at all

**Mechanism.** `completeWithWriter` derives the pin and then decides whether to gate at all:

- `driver_impl.go:487-495`: `pin` is read from the run store only `if o.base.Store.Enabled()`,
  and a **missing** cursor file is explicitly tolerated — `else if !os.IsNotExist(err) { return err }`
  means `os.IsNotExist` leaves `pin == ""` and continues.
- `driver_impl.go:496-507`: `contract, err := driver.ObserveChecksContract(o.ideaDir, pin)`;
  the verifier acceptance, the close gate and the transition check all sit behind
  `else if contract != ""`.
- `driver/checks.go:44-73`: with `expected == ""` and no `checks:` list, the only remaining
  activator is `os.Lstat(evidence.ReportPath(ideaDir))`. If that is `IsNotExist`,
  `ObserveChecksContract` returns `("", nil)`.

**Counterexample.** An idea runs several fix-up cycles under a list-form `checks:`, so
`EVIDENCE.json` and a v2 cursor exist. The implementer (or fix-up) agent, which has
repository write access — `driver_impl.go:467-473` re-invokes it on the working tree —
deletes three files: the `checks:` list from `00-prompt.md`, `parley-deck/ideas/<slug>/EVIDENCE.json`,
and `<store>/driver.json`. On the next `Complete`: `LoadCursor` returns `IsNotExist`, so
`pin == ""`; `ReadChecksContract` returns `(nil, false, nil)`; the `Lstat` of the report is
`IsNotExist`; `ObserveChecksContract` returns `""`. `requireAcceptedVerification`,
`EvidenceCloseGate` and `VerifyCompletionTransition` are all skipped, and control falls to
the legacy writer at `driver_impl.go:537-561`, which writes `status: complete`.

**Why this is CRITICAL and why the disposition does not cover it.** The disposition says
"Driver cursor v2 pins the normalized names/commands before agent work … Direct Complete
checks the actual run store's cursor and prior report", and separately "This is not same-UID
authentication or a claim that all history can be deleted safely." The second sentence is
true but does not apply here, because the threat model this slice exists to defeat is stated
in the source itself at `driver_impl.go:476-478` — "an implementer cannot short-circuit
review" — and the implementer is precisely a same-UID agent with write access to the repo.
The pin moved the safety state from `00-prompt.md` (idea dir) to `driver.json` (run dir), but
`parley-deck/runs/<run-id>/` is an ordinary in-repo path, as reachable by `rm` as the file it
replaced. `cursor.go:60-68` makes exactly this argument about the charged fix-up count —
"Review round-03 showed the marker-only count was still editable state — the class had moved,
not closed" — and the same reasoning applies to the pin. Note the asymmetry that makes
deletion the attack: a **corrupt** cursor fails closed (`LoadCursor` returns a non-NotExist
error → `Complete` returns it), while a **deleted** cursor fails open.

**Coverage caveat, stated rather than hidden.** I did not read `internal/driver/impl.go` at
this checkpoint. If `driver.Advance` independently refuses to reach `Complete` without a
readable cursor, this is unreachable *through the driver* while remaining reachable through
the direct `Complete` path the disposition itself names. `ADVANCE-PATH UNVERIFIED` — a
reviewer with a shell should settle it before acting, and it changes the remediation, not the
defect.

**Correction.** Make the *absence* of safety state a refusal rather than a permission:

1. In `completeWithWriter`, treat a missing cursor as fatal when the store is enabled:
   `LoadCursor`'s `IsNotExist` should deny, not yield `pin == ""`. A driver run that reached
   the completion step has necessarily written a cursor; its absence is evidence of
   tampering or of state loss, and both warrant a stop.
2. Do not gate on `o.base.Store.Enabled()`. `VerifyCompletionEvidence` already refuses
   without a store (`evidence_verify.go:319-321`); `Complete` should refuse symmetrically
   instead of silently dropping to the unpinned path.
3. Persist a one-way "this idea is under a named contract" marker in the run store the first
   time a contract is observed, and have `ObserveChecksContract` refuse when the marker
   exists and neither a list nor a report does. That closes the triple-deletion path without
   any new trust assumption: the marker only ever *adds* obligations, so forging one can only
   deny a close, never grant one — the same monotonic argument `cursor.go:60-68` already
   makes for the fix-up count.

### [MAJOR] Ledger-grade origin pinning is applied to a synchronization-only guard, so a user-cache wipe or a hostname change permanently denies all evidence publication for an idea

**Mechanism.** `AcquireResourceGuard` documents itself as charge-free — "No ledger or charge
is created" (`resource_guard.go:10-14`) — but reuses `lock()` unchanged, inheriting two
refusals whose only justification is protecting recorded money:

- `lock.go:87-90` opens the cache inode with `os.O_RDWR` and **no** `O_CREATE`, and
  `lockIdentity(path, newOrigin)` at `:282-284` refuses outright when the origin exists but
  the cache inode does not: "established local budget lock is missing at %s; refusing
  recreation".
- `pinLockOrigin` (`lock.go:192`) writes the hostname into the origin, and
  `checkLockOriginLocation` (`:328-335`) refuses on `parts[1] != host` with "hostname
  changed; preserve charges; a supported migration is not yet available".

**Counterexample A (cache wipe).** The guard's cache inode lives under `os.UserCacheDir()` —
`~/Library/Caches/parley/budget-locks/<key>.lock` on Darwin. That directory is routinely
purged by OS maintenance and by ordinary disk-cleanup tooling. After a purge, the origin file
still exists (it lives in the git admin dir, or in the deck for non-git), so `newOrigin` is
false, `lockIdentity` takes the `!create` branch and refuses. Every subsequent
`WithReportWriter` fails at `report_guard.go:53-56` with `evidence publication guard: …`.
Consequence: `runChecksContract` cannot record evidence (`driver_checks.go:76`), the helper
cannot publish (`evidence_verify.go:290`), and `Complete` cannot run (`driver_impl.go:480`).
The idea is unclosable and un-recordable until an operator finds and deletes
`<gitdir>/parley-evidence-guards/<id>/lock-origin`, which no message mentions.

**Counterexample B (hostname change).** macOS DHCP-assigned `.local` hostnames change on
network change. One such change makes `checkLockOriginLocation` refuse with a message about
preserving charges for a guard that holds none.

**Why this is a defect and not a safe default.** For a ledger, refusing recreation is
correct: a lost lock could mean a lost charge. For a mutex there is nothing to preserve —
recreating the inode loses no information, and the only cost of recreation is a window in
which two writers could both acquire, which `SaveIfUnchanged`'s CAS
(`report_guard.go:126-140`) and the finalization check (`:146-157`) already cover. The
disposition's framing — "reuses the already-tested budget lock primitive, with no monetary
accounting side effect" — is accurate about side effects and silent about inherited
*refusals*, which are the part that bites. `FINAL.md` "Idempotence & recovery" requires
resume from durable run state; a guard that permanently denies after a cache purge is the
opposite.

**Correction.** Give `AcquireResourceGuard` a synchronization-only variant of `lock()` that
(a) opens the cache inode with `O_CREATE` and recreates a missing identity, and (b) does not
pin hostname or cache path. Keep the no-op-filesystem probe at `lock.go:121-144` — that one
is about correctness of exclusion, not about charge preservation, and must stay. If sharing
one code path is preferred, add an explicit `preserveOrigin bool` so the difference is
declared in the signature rather than inherited by accident; a silent inheritance of a
money-safety rule into a mutex is how this class recurs.

### [MAJOR] The publication guard now rejects any `IMPLEMENTATION.md` whose frontmatter is not exactly normalized, including every CRLF document, and the failure surfaces as an unrelated diagnostic

**Mechanism.** `ReportWriter.save` gates *every* report write on parsing the implementation
document (`report_guard.go:146-157`), and `runChecksWithWriter` repeats the check
(`driver_checks.go:83-93`). Both route through `TransitionStatusToComplete`, whose
preconditions are byte-exact:

- `completion.go:76` requires `lines[0] == "---"` with no trimming.
- `completion.go:104-106` requires the status line to equal `"status: " + value` exactly.
- `completion.go:101` requires the value to match `^[A-Za-z0-9][A-Za-z0-9_-]*$`.

**Counterexample.** An `IMPLEMENTATION.md` with CRLF line endings — produced by a Windows
editor, or by a checkout with `core.autocrlf=true`. `lines[0]` is `"---\r"`, so
`TransitionFrontmatterStatus` returns `(nil, "", "document has no frontmatter block")`. In
`save`, `from` is `""`, so the finalized branch does not fire and the error branch does:
every report write fails with `cannot establish publication status: evidence: document has
no frontmatter block`. In `runChecksWithWriter` the same document fails the cycle at
`driver_checks.go:91-93` with that raw string as the check detail. The evidence path is
entirely unavailable for that idea, and nothing in the message mentions line endings.

This is a **regression in reachable states**, not a pre-existing constraint: the existing
readers tolerate CRLF because they trim — `cursor.go:426-437` compares
`strings.TrimSpace(line) == "---"`, and `ImplementationStatus` (`driver_impl.go:199-205`)
trims its value. So a CRLF document worked before this checkpoint and now cannot record
evidence at all. Two smaller variants land in the same place: an indented status line
(`  status: implemented` fails the normalized-form check at `:104`, because `trimmed` matched
but `lines[statusIdx]` did not), and a document with a *nested* key whose trimmed line starts
with `status:` — e.g. a `roles:` block containing one — which trips "duplicate frontmatter
status field" at `:89` before the YAML check at `:113-131` could correctly classify it as a
single top-level status.

**Correction.**

1. Normalize line endings for *detection* while keeping the *transformation* byte-exact:
   accept an optional `\r` suffix on the fence and on the status line, and preserve it
   verbatim in `outLines`. The transformation stays exactly invertible, which is what
   `VerifyCompletionTransition` relies on, and CRLF documents work again.
2. Restrict the line-prefix scan at `completion.go:87` to unindented lines
   (`strings.HasPrefix(lines[i], "status:")`), so a nested key cannot masquerade as a
   duplicate; the YAML pass at `:113-131` already handles the genuine duplicate case and does
   it correctly.
3. Whatever remains a refusal must say what to do. `save`'s wrapper at `report_guard.go:153`
   should name the file and the expected form ("`parley-deck/ideas/<slug>/IMPLEMENTATION.md`
   frontmatter must contain exactly one unindented `status: <token>` line"), because this
   error now blocks the entire evidence path and today reads like an internal parser fault.

### [MAJOR] A refused independent verification still leaves no durable committed record (retained, not re-litigated)

Carried forward from my `0a2022b` review and retained per the no-suppression rule. The
disposition — "receipts and invocation records are durable local ignored runtime artifacts by
the existing telemetry design; the outer driver also escalates errors. Whether an additional
bounded canonical summary is required remains an open disposition" — is an honest OPEN, and I
am not treating it as closed by assertion in either direction.

**What changed, and what did not.** The *success* path improved: `evidence-accepted.json` is
now written into the run store (`evidence_verify.go:508-509`), so a successful acceptance is
durable and independently re-checkable at `Complete`. The *refusal* path did not. The
receipt is still written only to `<root>/.parley-runtime/evidence-verification/attempt-*/result.json`
(`evidence_verify.go:231`, `:304`), and `evidence_verify.go:364-366` now **requires** that
tree to be git-ignored before execution — so by construction the refusal record is
unignorable-by-design and uncommitted-by-design. `commitEvidence` (`driver_checks.go:261-283`)
commits only `IMPLEMENTATION.md` and `EVIDENCE.json` and is not called on this path.

**Why it remains a finding.** `FINAL.md` "Idempotence & recovery" requires "Preserve failed
invocation records", and D2 requires terminal records for real outcomes. After a restart or a
runtime cleanup, "nobody verified" and "an independent verifier ran and the evidence did not
hold" are indistinguishable in the canonical record. The disposition correctly declines to
claim commit durability for ignored logs; that is precisely the gap.

**Correction.** As before: a bounded, scrubbed summary — verifier id, invocation id, request
digest, per-criterion statuses from `receipt.Executions`, refusal reason — reusing
`scrubAndTruncate` (`driver_checks.go:235-255`). If it lands in a new `EVIDENCE-ATTEMPTS.json`,
it must be added to `definedEvidenceArtifacts` (`driver_evidence.go:52-55`) **and** given its
own `ExtraDigests` binding in the same change; an excluded-but-unbound path is a hiding place,
which is the hazard `evidence_verify.go:402-403` warns about.

### [MINOR] The retained before/after tree digests are still tautological; the "measured values" correction is cosmetic

**Mechanism.** The disposition says "Passing those measured values through explicitly would
improve provenance clarity; this source cleanup remains open for review", and the source now
does pass them: `evidence_verify.go:259` captures `beforeTree`, `:265` captures `afterTree`,
and `:272` stores both into `VerifierExecution`. But both values come from
`verificationBindings`, which returns its recomputed digest **only after asserting it equals
the request** — `evidence_verify.go:168-174`: `if tree != req.TreeSHA256 { return … }`.

**Counterexample.** `TreeBeforeSHA256 != TreeAfterSHA256` is still unsatisfiable: both are
`req.TreeSHA256` on every path that reaches `:272`, because any other value makes
`verificationBindings` return an error at `:259` or `:265` first. So `reconcileRerun`'s
stability assertion (`evidence.go:404-406`, `CARRIED-FORWARD` from my prior reading) remains
structurally unable to fire, and `:407-409` still degenerates to comparing the request digest
against itself. The behaviour is correct — the bracketing checks are real, and they are the
thing actually providing stability — but the persisted record still documents an assertion
the code cannot make, and the defence-in-depth layer is still dead.

**Correction.** Have the bracket return the digest it measured *before* comparing, or measure
`evidence.TreeDigest` directly at `:259`/`:265` and let `verificationBindings` keep its own
equality check. One extra hash per criterion, and `evidence.go:404-409` becomes a live check.

### [MINOR] Three distinct failures share one diagnostic in the runtime-ignore precondition, and a non-git deck is silently unsupported

**Mechanism.** `evidence_verify.go:364-366` treats any non-zero exit from
`git -C root check-ignore -q -- .parley-runtime/` as "verification runtime must be ignored
before execution; add .parley-runtime/ to the repository ignore rules, rerun checks, then
retry". `check-ignore` exits 1 when the path is not ignored, and non-zero when the directory
is not a git work tree or `git` is absent.

**Counterexample.** A transport-A deck in a non-git directory — which the protocol permits,
and which this codebase supports elsewhere on purpose: `reportGuardDirectory` has an explicit
non-git fallback (`report_guard.go:113`) and `commitEvidence` tolerates a non-git tree
(`driver_checks.go:269-271`). Independent verification is unconditionally impossible there,
and the operator is told to edit ignore rules that would not help.

The precondition itself is the right shape and closes my prior runtime-poisoning MAJOR: it
runs before `MkdirAll` (`:383`), the tracked-file check at `:367-373` catches a
previously-committed runtime dir, the symlink check at `:374-382` rejects a non-directory,
and `:404-406` re-verifies the bindings after the two files are written. That part I could
not break.

**Correction.** Distinguish the causes: probe `git rev-parse --is-inside-work-tree` first and
emit a separate refusal for a non-git deck (naming it as an unsupported configuration for
independent verification, which is a real limitation worth stating), reserve the ignore
message for a genuine exit-1, and surface a git-not-found error as itself.

### [MINOR] In a non-git deck the guard's origin file joins the tested-tree digest, making that digest host-specific

**Mechanism.** Without a `.git` marker, `reportGuardDirectory` returns
`<ideaDir>/.parley-runtime/evidence-report` (`report_guard.go:113`), and `pinLockOrigin`
writes `lock-origin` there containing the hostname and the cache path (`lock.go:192`).
`definedEvidenceArtifacts` excludes only `IMPLEMENTATION.md` and `EVIDENCE.json`
(`driver_evidence.go:52-55`), so that file is inside the tested tree.

The ordering is deliberate and correct: `runChecksContract` acquires the guard at
`driver_checks.go:71` *before* `preDigest` at `:99`, so the origin exists before the first
hash and pre == post. This matches the disposition's "Non-Git check fixtures initialize
stable runtime metadata before hashing; no new digest exclusion is introduced", and I could
not break it.

**What the disposition does not say.** The recorded `TreeSHA256` for a non-git deck now
embeds host-specific content, so the same source tree hashes differently on two hosts, and
the guard's own maintenance (Counterexample A of the MAJOR above) changes the tested-tree
identity. Combined with the previous finding — verification refuses in non-git — the non-git
configuration is half-supported in a way no artifact states.

**Correction.** Either place the non-git guard outside the digest scope (a deck-level
`.parley-runtime/` sibling to `ideas/`, which is already the ignored convention elsewhere) or
state the non-git limitation explicitly in `IMPLEMENTATION.md`. Do not add a digest exclusion
for it — the disposition is right to refuse that.

### [MINOR] `replaceSection` matches its heading as an unanchored substring, so the human-readable evidence table is post-hoc forgeable and a `###` subheading mis-splices the document

**Mechanism.** `replaceSection` (`driver_checks.go:317-331`) locates the section with
`strings.Index(doc, heading)` — no line anchor — and the span it replaces or strips runs to
the next `"\n## "`. The same function both writes the table (`:311`) and defines the bound
digest space (`implementationRestContent`, `:207`).

**Counterexample A (forgeable table).** Because the span between the heading and the next
`## ` is stripped before hashing, an implementer may rewrite that span arbitrarily after the
evidence is recorded — deleting the real per-criterion table and substituting prose claiming
every criterion passed — and `implementationRestContent` produces identical bytes, so
`EvidenceCloseGate`'s binding check at `driver_evidence.go:108` still matches and the close
proceeds. The typed `EVIDENCE.json` remains bound and authoritative, so this is not a false
close; the consequence is that the human-facing evidence table carries no integrity guarantee
and must not be cited as evidence. Attempts to widen the span to swallow real content fail:
moving the heading earlier makes the next `## ` terminate the span at the following heading,
so the swallowed text reappears in the digest and the binding breaks. I tried that and could
not break it.

**Counterexample B (mis-splice).** `strings.Index("### Validation evidence", "## Validation evidence")`
is `1`. A document containing a `### Validation evidence` subheading — or the literal string
inside a fenced block in `## Notes for reviewers` — makes `idx` point one byte in, so
`doc[:idx]` retains a stray `#` and the driver writes its table into that location while the
real section keeps a stale table.

**Correction.** Anchor the match to a line start: search for `"\n## Validation evidence\n"`
(plus the start-of-document case) rather than a bare substring, and require the terminator to
be a line-start `## `. If the table is to be trusted by humans at all, additionally record
its digest in `ExtraDigests` under a distinct key so a forged table is detectable — the
binding cost is one hash.

### [NIT] The publication guard is not reentrant, and the exported `evidence.Save` is the obvious way to deadlock it

`evidence.Save` (`report.go:33-38`) opens its own `WithReportWriter`, and the guard is a
`flock` on a freshly opened descriptor. `flock` is per-open-file-description, so a second
open in the *same* process does not inherit the lock — the source proves this deliberately at
`lock.go:123-144`, where the probe's `take(probe)` succeeding is treated as a broken
filesystem. Any call to `evidence.Save` from inside a `WithReportWriter` callback therefore
blocks for the full 30 s (`report_guard.go:51`) and then fails with `ErrLockContention`.
Today's in-scope callers are correct — `writeTypedEvidence` takes a `*ReportWriter`
(`driver_checks.go:159`, `:181`) and the helper uses `writer.SaveIfUnchanged` — but the trap
is unlabelled. `SAVE-CALLER-SET UNVERIFIED`: I could not enumerate `evidence.Save`'s callers
without a shell. Document the non-reentrancy on `Save` and on `WithReportWriter`, or have
`Save` detect an already-held guard in-process and return a clear programming error rather
than stalling.

### [NIT] `expectedTransition := updated` shallow-copies a struct whose `Records` and `ExtraDigests` are shared with the value used for the central integrity comparison

`evidence_verify.go:477-479` copies the report by value and passes `&expectedTransition` to
`AuthorizeCompletionTransition`. `Records` (slice) and `ExtraDigests` (map) share backing
storage with `updated`, which is then compared to `original` at `:501` —
`sameVerificationJSON(updated, original)` — as the check that the verifier "changed the
original evidence instead of only attesting it". It is correct today: `completion.go:160-206`
only reads those fields and writes the `CompletionTransition` pointer, which is not shared.
But a future mutation inside `AuthorizeCompletionTransition` would silently corrupt the value
the strongest remaining check is computed over. Deep-copy via a JSON round-trip (the codebase
already has `sameVerificationJSON` for exactly this reason) or document the read-only
contract at the function.

### [NIT] Provenance is reset positionally while the receipt is matched by name

`evidence_verify.go:486-500` correctly keys the receipt lookup by criterion name — my prior
positional MINOR is fixed there — but line `:499` still assigns
`updated.Records[i].Provenance = original.Records[i].Provenance` by index. Safe today, because
`AttestExecution` mutates in place and the whole-report comparison at `:501` catches any
reordering; the mixed convention within eleven lines is the kind of thing that later gets
copied without its companion guard. Key both by name.

### [NIT] Report decoding strictness differs between the gate and the acceptance boundary

`readVerificationJSON` sets `DisallowUnknownFields` (`evidence_verify.go:98`) and rejects
trailing content, while `evidence.Load` uses a permissive `json.Unmarshal`
(`report.go:88`). So a report with unknown fields passes `EvidenceCloseGate` and
`ObserveChecksContract` but is rejected at `acceptVerification`/`requireAcceptedVerification`.
Byte-digest comparisons make this harmless in practice, and the divergence is in the safe
direction, but two strictness levels for one canonical artifact is a latent inconsistency.
Use the strict reader in `evidence.Load` too.

### [NIT] On Windows the legible refusal exists but is not the one an operator sees at close

`VerifyCompletionEvidence` refuses precisely (`evidence_verify.go:316-318`), which is the
correction I asked for. But a named-contract idea that reaches `Complete` on Windows fails at
`requireAcceptedVerification` with "independent verification acceptance unavailable"
(`evidence_verify.go:519`), which does not mention the platform. Add the same `runtime.GOOS`
refusal at the head of `completeWithWriter`'s contract branch.

---

## Dispositions I evaluated — where I concur and where I do not

| Disposition | My independent evaluation |
|---|---|
| CRITICAL contract-removal: "corrected … cursor v2 pins … Direct Complete checks the run store's cursor and prior report" | **Partially concur.** Closed for a present pin and for the prior-report activator (`driver/checks.go:48-71`). **Do not concur that it is corrected**: the pin is deletable and a missing cursor yields `pin == ""` (`driver_impl.go:489-495`). See the CRITICAL. |
| MAJOR completion self-invalidation: "corrected using Kimi's API … parent recomputes the same authorization" | **Concur.** `completion.go:220-262` recomputes every field and inverts the transition; `driver_impl.go:527-534` re-derives before/after bytes and re-checks before the write. No exclusion was widened. I could not break it. |
| Transition probes (short-digest panic, altered `AfterSHA256` with non-complete status, duplicate YAML spellings) | **Concur that the source now guards all three**: `completion.go:171`, `:249-256`, `:107-131`. I did not execute the probes; `TEST-EXECUTION UNVERIFIED`. |
| MAJOR tautological tree fields: "the reviewer explicitly confirms current runtime behavior is correct … passing those measured values through explicitly would improve provenance clarity" | **Do not concur.** I did confirm the behaviour is correct and still do. The values are now passed through but remain unsatisfiably equal, because `verificationBindings` returns only after asserting equality with the request (`evidence_verify.go:168-174`). Provenance clarity did not improve. Retained as MINOR. |
| MAJOR receipt-write failure "remains unverified; permanent-failure claim" | **Concur that it is corrected in source, with one unverified link.** A fresh `RunChecks` writes an un-attested report (`report_guard.go:146-158` permits it while status != complete), which clears the `CompletionTransition != nil` refusal at `evidence_verify.go:348`. That the driver runs `RunChecks` before each closing `Verify` is `ADVANCE-PATH UNVERIFIED` (I did not read `driver/impl.go`). My prior "permanently unclosable" claim does not hold against this source; I withdraw it as a live defect. |
| MAJOR concurrent report overwrite: "OPEN … no concurrency acceptance is claimed" (dispositions) vs. serialized publication (publication note) | **Concur that it is now closed for cooperating same-host writers** (`SaveIfUnchanged` + one guarded transaction + receipt inside the guard). Two residual limits, honestly held by the notes themselves: uncooperative writers, and two mount paths for one directory — `EvalSymlinks` plus `strings.ToLower` (`report_guard.go:68`) unifies symlink and case aliases but not distinct mount points of one volume, which then yield two guards over one report. The CAS remains the only protection there. |
| MAJOR Windows quoting: "OPEN … explicit supported-platform restriction remains needed" | **Concur, and it has landed** as an explicit refusal (`evidence_verify.go:316-318`). The POSIX quoter at `:414` is now unreachable on Windows. One NIT above on the `Complete`-path message. |
| MAJOR refusal record not committed: "remains an open disposition" | **Concur it is open.** Retained as MAJOR with the mechanism above; the new ignore *requirement* makes the runtime record structurally uncommittable. |
| MAJOR runtime path not ignored: "OPEN usability/recovery issue … validate the prerequisite before creating artifacts" | **Concur, and it has landed** (`evidence_verify.go:362-382`, `:404-406`). Downgraded to the MINOR above, which is now about diagnostic conflation and the non-git case. |
| MINOR positional receipt matching: "OPEN cleanup" | **Corrected** for the receipt (`:486-497`); one positional remnant at `:499` (NIT). |
| MINOR `time.Time` DeepEqual | **Corrected.** `sameVerificationJSON` (`:553-557`) compares persisted JSON. |
| MINOR repeated whole-tree hashing: "retained to enforce each actual criterion's before/after boundary" | **Concur.** This is the right call, and it is now load-bearing rather than redundant: it is the *only* real stability guard, per the tautology finding. Do not hoist the in-loop calls. |
| NIT stale wiring comment / inconsistent path predicates | The `driver_evidence.go:14-34` header is now accurate about call sites. The predicate split persists: `driver_evidence.go:49` checks the exact `rel == ".."` case, `evidence_verify.go:211` omits it (covered there by the `filepath.Base(rel) != "request.json"` and `filepath.Dir(rel) == "."` tests, so no hole). Still one helper's worth of cleanup. |
| Readiness (D7) and process-budget (D6) dispositions | **Unevaluated.** Not inspected; see Actual coverage. |

## Facilitator-reported executions (testimony, not verified by me)

Transcribed so a later reader can locate them. I observed none of these runs, saw no logs,
and assign them **no** verdict — `TEST-EXECUTION UNVERIFIED` for every item: separate-process
CAS with one winner; delayed writers refusing after guarded completion; built CLI/helper/
criterion subprocess fixtures recovering after report or receipt persistence failure; changed
accepted report or removed receipt vetoing `Complete`; missing and probe-only ignore refusing
before artifact creation; focused app/evidence, evidence/budget race, shared-volume, vet and
Windows cross-build passes. The publication note's own account of a full Go command that
exited zero while its streamed shared-volume log was truncated to a 133-byte prefix plus 1270
NUL bytes is recorded here as the implementer's disclosure; retaining it as a damaged log
rather than validation output is the correct handling, and the pending captured-output rerun
is not evidence until it exists.

Independently of that testimony: a passing in-process or subprocess fixture suite is not a
live real-model verifier trial. AC-E2 requires that "an independent verifier executes the
concurrency counterexample and corrected case"; the notes correctly claim no live real-model
verifier closure, so AC-E2 is **not** evidenced. `NOVELTY/COVERAGE UNVERIFIED` for any
end-to-end claim.

## Open questions

1. **`ADVANCE-PATH UNVERIFIED`** — does `driver.Advance` refuse to reach `Complete` when the
   cursor is missing, and does it run `RunChecks` before every closing
   `VerifyCompletionEvidence`? The first changes the CRITICAL's reachability through the
   driver; the second is the load-bearing link in the receipt-failure recovery I accepted
   above. Both are a single read of `internal/driver/impl.go`. **Check this first.**
2. When is `Cursor.ChecksContractSHA256` first written relative to the implementer's launch?
   The disposition says "before agent work". If it is written only after a successful first
   contract run, a contract added mid-implementation has an unpinned window.
3. After a successful close, `runChecksWithWriter` fails with `ErrReportFinalized`
   (`driver_checks.go:87-90`). Does a resumed driver run on a `PhaseDone` idea ever call
   `RunChecks`, and if so does that surface as a spurious check failure and escalation?
   `Rebuild` sets `PhaseDone` at `cursor.go:288`, which suggests not, but I did not read
   `Advance`.
4. `evidence.RunCriterion`'s format detection, case counting and output hashing remain
   unread. `RUNCRITERION UNVERIFIED` — `reconcileRerun` reconciles against it.
5. `PARLEY_PROC_MARKER` ↔ `res.InvocationID`: `evidence_verify.go:454` requires
   `receipt.ProcessMarker == invocationID`. I still have not located the code that exports
   that variable into the child, so the binding's load-bearing fact is unconfirmed.
   `MARKER-BINDING UNVERIFIED`, carried forward unchanged from my prior review.
6. Two mount paths for one idea directory produce two independent guards
   (`report_guard.go:67-113`). Is that inside or outside the declared cooperative boundary?
   The notes say host-local, not distributed; a single host with two mounts of one share sits
   in the gap.

## Assessment

The corrections at this checkpoint are substantive and were made in the harder, honest
direction: the completion transition binds the status flip by independent recomputation rather
than by excluding it, the publication guard plus CAS closes a real race rather than
documenting it away, and the runtime-ignore prerequisite refuses before it can poison the
tree. Six of my ten prior findings are genuinely closed, and one — the "permanently
unclosable" claim — I withdraw against this source.

Two things I would not let pass on trust. The CRITICAL is the same class as the one it
replaced: the gate still asks what artifacts currently exist rather than what the run has
already obliged itself to prove, and all three activators are deletable by the party the gate
constrains. And the guard MAJOR is a case of a mechanism carrying its old justification into a
new role — ledger-preservation refusals inside a mutex — where the failure mode is a
permanent denial from a routine cache purge.

Nothing here is acceptance. AC-E2 remains unevidenced, no live real-model verifier trial has
been shown to me, D6 and D7 are outside my coverage at this checkpoint, and I executed
nothing. My design-signoff position is unchanged: independent verification is an obligation,
not a passed gate.
