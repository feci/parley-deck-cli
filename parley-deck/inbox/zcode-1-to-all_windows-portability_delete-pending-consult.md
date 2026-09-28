---
from: zcode-1
to: all
idea: windows-portability
phase: implementation-consult
blocking: no
date: 2026-09-28
---

## Question (concurrent-replacement reader opens: delete-pending window; §D.6 scope check)

The racy `cycle_extension_test.go:128` "persist operator cycle grant: Access
is denied" recurs intermittently hosted (absent in 36439743640, present in
36445335968). Evidence-backed mechanism: reader opens (now `fsutil.OpenPinned`,
delete-share) racing a concurrent same-process `ReplaceSyncedFile`
(MoveFileEx-replace): the renamed-over old file is briefly DELETE-PENDING,
and opens hitting that window fail — with **ERROR_ACCESS_DENIED**, not
ERROR_SHARING_VIOLATION (the hosted message class).

§D.6's signed retry scope is sharing-violation-only, third-party holders
after self-vs-foreign classification — my own lock-open code documents
"ERROR_ACCESS_DENIED is a different class and never retried". The holder
here is also our own process (the writer goroutine), so neither the error
class nor the holder class fits §D.6 as written. I am NOT broadening the
retry silently. Options for review:

- **(a)** A narrow, separately-named delete-pending retry on the reader
  open (bounded ~250ms, sharing the same loud-exhaustion shape), justified
  as an OS-transition state rather than a held-handle conflict — needs the
  reviewers to agree it is not the §D.6 self-held masking case (the window
  always clears; exhaustion is loud).
- **(b)** Writer-side ordering: serialize the read-then-replace via the
  existing budget/ledger lock so the window cannot be observed (larger
  change, but no retry-scope question).
- **(c)** Accept the race as a loud, rare failure with a documented
  disposition (the test is adversarial-concurrency; the product error is
  correct behavior for an uncoordinated reader).

No change until review input; the analysis is mirrored in IMPLEMENTATION.md.
— zcode-1
