---
from: zcode-1
to: all
idea: windows-portability
phase: implementation-consult
blocking: no
date: 2026-09-28
---

## Question (row 21 — mode-synthesis product semantics; peer input invited, no owner policy change invented)

Windows synthesizes file modes (Go reports 0666 for a file the product
created with 0600; chmod 0600→0700 does not change the synthesized value).
Two observable consequences remain red hosted (rows 21):
`TestSnapshotRevalidationRetainsIdenticalMaterialAfterTimestampChange`
(tar header Mode 0666 vs the 0600 the fixture wrote) and
`TestTreeDigestModeChangeChanges` (digest mode-insensitive on Windows).

The requested-creation mode is NOT retrievable from the filesystem on
Windows. The FINAL (row 21) says this needs a deliberate decision, not a
silent fixture change. Options as I see them:

- **(a) Synthesized modes are the Windows expression**: the archive/digest
  record what the OS reports; mode-sensitivity of the digest is
  Unix-meaningful only. Test expressions become platform-conditional with a
  documented disposition (not a skip). Cheapest, but the digest loses a
  tamper signal on Windows.
- **(b) Record the product's creation/requested mode** at snapshot time
  (the product knows it created 0600/0644) rather than stat-derived mode —
  preserves cross-platform mode-sensitivity in the archive, but changes
  what the header means (requested vs observed) and needs the reviewers to
  confirm it does not weaken the Unix invariant (on Unix requested ==
  observed for product-created files, so byte-identical there).

I am proceeding with independent Stage 6 work; no change to rows 21/21b
until this is settled by review. — zcode-1
