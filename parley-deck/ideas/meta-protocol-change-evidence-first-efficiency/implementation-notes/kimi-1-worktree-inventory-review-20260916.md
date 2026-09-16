---
agent: kimi-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-16
reviewed-source: immutable snapshot /private/var/folders/yt/p2sr23f12_qcfx_w2z5c1p4r0000gn/T/parley-complete-recovery-1h8k1qp7/source — claude-1 MRW-1 stage-one slice only (internal/budget/worktree_inventory.go + test, internal/app/budget_worktree.go + test, internal/app/budget.go dispatch, candidate note)
review-kind: supporting-source-only
---

# Supporting source review: claude-1 read-only worktree inventory (MRW-1 stage one)

Scope: the slice above only, against the launch-task invariants and MRW-1 item 1. No test, Git, or
CLI run; no round-02 content, phase change, consensus, signoff, or acceptance — round-02 stays
blocked. Verdicts are PRIMARY (locator + quote) or explicitly weaker; I own no code in this slice.

## Refutation attempts and PRIMARY verifications

1. **Retain-all + unknown-never-zero — holds.** `worktree_inventory.go:189-201` retains every
   registration and appends an uncertainty line per unavailable row ("retained unpruned, not
   counted as absent"); `:211-225` classifies only via `os.Lstat` on the raw path into available/
   missing/not-directory/indeterminate, `history_coverage` limited to local/unknown (`:29-32`).
   Closed key sets are test-asserted (`worktree_inventory_test.go:120-121`): no count/cap/floor
   or attestation field exists to read as zero.
2. **launchScope refusal unweakened — holds.** Snapshot `binding.go:101-104` still raises
   `"historical worktree is unavailable: %s"`, identical to integration-root `binding.go:102`; root
   `budget.go` has no `worktree` branch — the addition is the only dispatch change (`:56-60`, usage
   `:74`); the witness test asserts `ConfigureLaunchBudget` still fails with that string (`:143-145`).
3. **Refusal to mutate — holds.** The new files issue no write/remove/mkdir: only two `rev-parse`
   probes, one `git worktree list --porcelain -z`, `Lstat`/`EvalSymlinks`, JSON to stdout
   (`worktree_inventory.go:124-207`; `budget_worktree.go:18-46`); dispatch is ungated
   (`budget.go:56-60`). `budget_worktree_test.go:100-142` refuses apply/attest/prune/repair/
   recover/remove, stray flags, and `--expected-observation-sha256` even attended+supported, with
   no state created and no unimplemented verb in the usage text.
4. **Malformed authority/porcelain — holds.** The strict parser (`worktree_inventory.go:230-325`)
   refuses unknown/duplicate attributes, duplicate/relative/empty paths, branch+detached,
   bare+HEAD/branch, neither-bare-nor-HEAD, bad HEAD (non-40/64-hex), non-`refs/` branch, valued
   flags, nested records, empty/unterminated output, >16 MiB, >4096 rows; 17 fabricated cases plus
   a counter-test accepting every shape real Git emits (`:315-375`).
5. **Bare and non-Git roots — holds.** Bare: `--show-prefix` failure is non-fatal (uncertainty
   line, registrations still read; `:173-180`) and the parser accepts a `bare` record. Non-Git:
   `git:false`, `covers_repository:false`, explicit `single-directory/non-git` scope, digest still
   computed (`:151-161`); broken authority (`.git` present, `rev-parse` fails) refuses, mirroring
   `binding.go:122-130` (`:136-150`).
6. **Symlink/alias — holds, conservative.** Aliased registration → `not-directory`/unknown
   (Lstat does not follow), stricter than `launchScope`'s `os.Stat`; candidate Limitation 7
   discloses the divergence. Caller-root aliases are canonicalized (`:126-134`): alias roots of
   one repository share a digest.
7. **Cross-root digest identity + deterministic restoration — holds.** `OperatorRecord`
   (`:101-109`) marshals the closed six-field `observedWorktrees`, registrations sorted by raw
   path, no clock, root/prefix excluded, `git_common_dir` carrying identity; digest reuses `key()`
   (`ledger.go:131`). Test 2 asserts identical digests from main worktree, linked worktree, and
   subdirectory; test 3's rename round-trip restores the exact prior digest (recreate would stay
   prunable; rename is the right witness).

## Findings (none blocking the read-only stage)

- **[MINOR-1] `history_coverage: local` over-claims for available-but-unreadable Git metadata.**
  Classification rests solely on Lstat+IsDir (`:215-219`), yet the comment says coverage records
  "whether a registration's local history is readable here" (`:27-28`) — a gutted or mismatched
  `.git` pointer is not checked. Mitigated in-band: `prunable` is carried verbatim and each such
  row adds an uncertainty line (`:198-200`); Limitation 10 names the case. Correction now: narrow
  the comment to the check performed. Attest stage: `local` must not read as "history verified";
  whether Git's prunable catches every mismatched pointer is UNVERIFIED (RECALL; no Git run).
- **[MINOR-2] Non-Git digest is location-independent.** Root is excluded by design (`:83-88`), so
  every non-Git directory hashes identically. Honest for a stage-one preview (uncertainty
  discloses scope), but the documented pinning seam (`--expected-observation-sha256`) would be
  vacuous across non-Git directories; decide before any attest stage consumes the digest.
- **[NIT-1] Output shape deviates from MRW-1 item 1's "partitioned into available/unavailable".**
  Flat `registrations[]` + per-row availability is information-equivalent and the digest is
  stronger; record the deviation in the attest-stage spec.
- **[NIT-2]** No-mutation witness compares admin-dir entry counts, not name sets
  (`worktree_inventory_test.go:131-134`); same-count delete+create would pass. Test-strength only.
- **[NIT-3]** `assertWorktreeFields(observedWorktrees{}, …)` would miss a future field carrying
  `omitempty`; today's six fields have none.

## Existing alternatives and residual uncertainty

Alternatives in-tree: `launchScope`'s own enumeration (`binding.go:95-113`) is welded to the
refusal the task forbids weakening; duplicating the identical wire read is justified (a shared
helper is a later refactor, drift risk noted). Skipping `readStepHistoryFile`/`validateHistoryJSON`
is right — no file or JSON input exists here. Residual: I executed nothing; every real-Git behavior
claim (removed registration kept with HEAD/branch; restore clears prunable; byte-identical rerun)
comes from reading code and tests — the coordinator's independent run, not this note, is the verdict.

## Conclusion

Scoped support for stage one as a read-only inventory: every launch-task invariant (retain all,
unknown-never-zero, no prune/recreate/mutation, launch admission unweakened, stable material
digest, no attest/apply) verified by PRIMARY reads. Required: MINOR-1's comment precision now;
MINOR-2 and the NIT-1 deviation recorded before the attestation/apply stage. Not whole-audit
acceptance, not a test claim, not a signoff.
