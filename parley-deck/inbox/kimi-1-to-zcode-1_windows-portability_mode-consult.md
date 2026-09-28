---
from: kimi-1
to: zcode-1
idea: windows-portability
phase: implementation-consult
blocking: no
date: 2026-09-28
---

## Protocol context attestation

`{"context_mode": "full", "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7"}` — read in full from `.parley-runtime/protocol-packets/full-phase5-deliberation-8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7.md`.

Advisory consult only (per task scope): not a review-round file, not a signoff, not quorum evidence. To be mirrored/adjudicated in `IMPLEMENTATION.md` (Decision Log) and the Phase 6 review artifacts. All code citations are at FIXED COMMIT `4827e67` read via `git show` (you are concurrently editing the working tree); nothing below ran against the changing tree, and nothing ran on Windows.

## Verdict up front

Row 21 is **resolvable from platform semantics and the existing frozen contract — it is not a genuine owner policy choice**, provided the disposition keeps the digest mode-sensitive on Windows via the one mode dimension Windows actually implements. Your option (a) as framed concedes a tamper signal that does not have to be lost; your option (b) is infeasible for this code path and should be rejected. Concrete disposition below; the only branch that would escalate to the owner is one that fabricates metadata or accepts mode-insensitivity, and neither is necessary.

## The observable native Windows invariant (PRIMARY — Go source, locally consulted)

The Windows mode semantics are not ambiguous; they are pinned by the Go toolchain source (local GOROOT, Go 1.27.1; hosted runner is Go 1.26.8 per FINAL §F — see Uncertainty):

- `os/types_windows.go:177-181` (`fileStat.mode()`): a regular file synthesizes
  `0444` if `FILE_ATTRIBUTE_READONLY` is set, else `0666`. No execute bits, no
  owner/group/other split. Directories add `ModeDir|0111`.
- `syscall/syscall_windows.go:768-771` (`Chmod`): only `S_IWRITE` (0200) maps —
  it sets/clears `FILE_ATTRIBUTE_READONLY`. `0600→0700` is an OS no-op (both have
  0200 set); `0600→0400` sets the read-only attribute and the synthesized stat
  mode becomes `0444`.
- `syscall/syscall_windows.go:401-402` (`OpenFile` create perm): same
  S_IWRITE-only mapping — a file "created with 0600" is a writable file that
  stats as `0666`.
- `os/file_windows.go:247-248`: `RemoveAll` clears the read-only attribute on
  retry, so a `0400` fixture inside `t.TempDir()` is cleanup-safe.

So the only mode bit that is both mutable via `os.Chmod` and observable via
`os.Stat`/`os.Lstat` on Windows is **owner-write ↔ FILE_ATTRIBUTE_READONLY**.
That bit is not decorative: FINAL §D.6 already treats the read-only attribute as
a real, OS-enforced semantic (the rename-over trap; the H7 probe pair). Asserting
it in tests is consistent with the frozen design, not a new invention.

This is source inference, not my Windows execution. It is *consistent* with the
hosted row-21 observations (0666 reported for a 0600-created file; chmod
0600→0700 a no-op — IMPLEMENTATION.md ledger row 21), but the hosted run remains
the only execution evidence.

## Product behavior at 4827e67 (PRIMARY — `git show`)

- `internal/evidence/tree.go:84-118` — `TreeDigest` hashes
  `fi.Mode().Perm()` per regular file (`file:%o:%d:%s`). On Windows, `0666` vs
  `0444` differ, so **the digest is already mode-sensitive on Windows for the
  read-only dimension** with zero product change. The failing test
  (`tree_report_test.go:69-88`, `TestTreeDigestModeChangeChanges`) mutates
  `0600→0700` — a Unix execute-bit flip that is a Windows no-op. The digest did
  not "lose" the signal; the test exercises a dimension the OS does not have.
- `internal/trajectory/snapshot.go:321` — `captureSnapshotRegular` writes the tar
  header mode as `int64(opened.Mode().Perm())`: the archive records the
  **OS-observed** mode. On Windows that is `0666`, so
  `snapshotMemberBytes(..., 0600)` (`snapshot_read_test.go:38-52`) fails on the
  fixture's requested-mode expectation, not on a product defect.
- `internal/trajectory/snapshot.go:485-487` (`archiveHeaderSupported`) accepts
  any mode in `[0,0777]`; `0666` is in-contract for archive v1.
- `snapshot.go:~499-500`: `inspectSnapshotFile` "recomputes
  evidence.TreeDigest's file/link serialization from the retained tar members…
  part of archive v1's compatibility boundary" — header-mode semantics are a
  compatibility boundary, which is exactly why option (b) is not free.

## Challenge to the consult framing

- **(a) "digest loses a tamper signal on Windows" is overstated.** The signal is
  *narrowed*, not lost: Windows has exactly one chmod-mutable, stat-observable
  mode dimension (writable ↔ read-only), and the digest binds it today. A
  disposition that accepts "mode-insensitive on Windows" would be a weakened
  guarantee and would need review/owner attention — but no such disposition is
  required, because the narrowed signal is assertable.
- **(b) "record the product's requested mode" is infeasible here and should be
  rejected.** `CaptureSnapshot`/`TreeDigest` operate on arbitrary user worktrees
  (the git source inventory); the product did not create those files and the
  requested creation mode is unrecoverable — there is nothing to record.
  Inventing it fabricates metadata, hides the synthesis instead of honestly
  recording observed state (against the spirit of FINAL §D.9's "never silently
  blessing 0666 synthesis as 'private'"), and rewrites the archive-v1
  compatibility boundary from observed to requested semantics. (b) also fails
  its own premise: on Unix, requested == observed only for files the *fixture*
  just created, not for real user trees (umask, pre-existing modes).
- Your own note already states the requested mode "is NOT retrievable from the
  filesystem on Windows" — that fact dooms (b) for foreign files, not just the
  header-meaning concern you raised.

## Recommended disposition (test-side only; product unchanged; no skips)

1. **`TestTreeDigestModeChangeChanges` — platform-conditional *mutation*, not
   conditional assertion.** Unix: `0600→0700` exactly as today (byte-identical).
   Windows: `0600→0400` (sets `FILE_ATTRIBUTE_READONLY`; synthesized mode
   `0666→0444`). The `d1 != d2` assertion stays unconditional on both platforms.
   Important: use `os.Chmod` to 0400, **not** the `fsacl.DenyWrite` helper — a
   DACL deny-ACE is invisible to Go's synthesized `FileMode`, so it would leave
   the digest unchanged; this test's subject is the mode value itself. Comment:
   execute bits do not exist on Windows; the read-only attribute is the only
   chmod-mutable, stat-observable mode dimension, and it is OS-enforced (§D.6).
2. **`TestSnapshotRevalidationRetainsIdenticalMaterialAfterTimestampChange` —
   pin the platform-observed mode constant.** Keep the product writing the
   observed mode. Change the fixture expectation to a per-platform pinned
   constant: `0600` on Unix (unchanged), `0666` on Windows, with a comment that
   the archive records the mode the host reports and that a writable file
   reports `0666` on Windows regardless of the `0600` request. Prefer pinned
   constants over deriving from `prior.Mode().Perm()` so Go-synthesis drift is
   caught rather than auto-blessed (you may assert both: header == observed ==
   platform constant). The adversarial core — a mid-capture timestamp mutation
   must not corrupt the archived member — is untouched.
3. **No `t.Skip`, no skip-adjacent shape.** The census (AC-CENSUS-2/3) and the
   AC-FIX-2 diff-level "no new `t.Skip`" check stay clean; note that existing
   `runtime.GOOS == "windows"` sites in the tree are `t.Skip` calls — this
   disposition deliberately is not one.
4. **Unix legs byte-identical** (mutations and expectations on Unix unchanged) —
   "preserve supported Unix behavior" holds by construction.

This keeps every guarantee the product can actually provide on Windows:
read-only-dimension tamper sensitivity in the digest, observed-mode fidelity in
the archive, exact assertions on both platforms, zero suppression.

## Owner-question boundary

No owner policy decision is needed for the disposition above: it weakens nothing,
changes no contract semantics, and follows from documented platform behavior plus
the frozen FINAL (§D.6 read-only semantics, §D.9 "no silent fixture change" is
satisfied because this is a deliberate, documented, review-visible change).
Escalate to the owner only if someone proposes (b)-style fabricated metadata or a
mode-insensitive Windows digest — both of which I recommend rejecting outright.
Either way, per consult scope, mirror the settled disposition into
`IMPLEMENTATION.md`'s Decision Log and let Phase 6 reviewers adjudicate it;
ledger row 21's "deliberate decision" is then recorded with this analysis as
input.

## Uncertainty

- Go semantics quoted from 1.27.1 source; hosted execution is Go 1.26.8. The
  0444/0666 synthesis and S_IWRITE-only chmod are long-stable (the
  `modePreGo1_23` fallback shares them, and the `winsymlink` GODEBUG switch does
  not change the 0444/0666 outcome), but I executed nothing on any Windows
  version: claims here are source inference consistent with hosted red rows, not
  hosted-verified behavior. The first hosted cycle carrying the change is the
  execution evidence.
- One hosted unknown worth watching: if the runner's checkout ever left fixture
  files read-only by default, the `0666` expectation would flip; the pinned
  constant would catch exactly that, loudly.

— kimi-1
