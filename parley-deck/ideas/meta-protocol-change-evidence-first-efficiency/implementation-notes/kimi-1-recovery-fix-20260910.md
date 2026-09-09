# kimi-1 recovery fix handoff — 2026-09-10

Run: continuation of f7248aa0 (timed out). Deadline 30m from ~00:48 CEST (hard stop ~01:18, code freeze ~01:10 for final notes).
Snapshot under fix: evidence source SHA256 a71f397b8419ba715e3ffa1eea99bb1c6fafd0fe2603f206038252fa5f3c84a5.
Scope: internal/evidence/**, internal/app/driver_checks.go(+test), driver_evidence.go(+test) only. No git mutations.

## Facilitator-reported failures being addressed (5)

1. Evaluate: empty Provenance.Verifier + foreign ClosureOptions.Verifier closed a pass. Fix: require persisted per-record verifier attestation bound to tree+criterion (new attestation API recording independently EXECUTED checks); reject empty scope, dup/missing criteria, invalid/negative counts, contradictory exit/status, absent CurrentTreeSHA256 (was: zero reasons).
2. ParseEnvelope: earlier valid pass survived a later malformed PARLEY-EVIDENCE line. Fix: last envelope line wins; malformed/null/partial/duplicate/negative/conflicting counts fail closed; no fallthrough to unrelated parser.
3. TreeDigest: symlink retarget (a→b) and mode change (0600→0700) invisible. Fix: hash entry type+mode+link target; fail on lstat errors/unsupported entries; canonicalize root (/var vs /private/var) before escape check.
4. writeTypedEvidence hashed only AFTER command run. Fix: digest before execute, re-verify after, fail attempt on change; preserve failed-report evidence; narrow IMPLEMENTATION.md exclusion to generated evidence content only.
5. Save suppressed ENOTTY/EINVAL/EOPNOTSUPP from tmp.Sync. Fix: use internal/fsutil.SyncFile (ENOTTY-only fallback to syscall.Fsync); other errors remain failures.

## Additional in-slice concerns

- RunCriterion: unbounded output buffer + only kills sh (child-held pipes defeat cancellation) → bounded capture/hash + procctl process-group cleanup (procctl read-only).
- Do not persist raw Command strings (secret-bearing): digest + safe representation only.
- Add serial-vs-barrier concurrency fixture test.

## Status log

- 00:48 start; read evidence.go; handoff created early per instructions.
- (updated as work proceeds below)

## Unresolved at deadline

- (to be filled in final update — nothing claimed complete without tests)
