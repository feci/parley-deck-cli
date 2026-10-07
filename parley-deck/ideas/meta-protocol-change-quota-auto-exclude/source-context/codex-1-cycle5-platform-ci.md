# Cycle-5 platform CI — 2026-10-08

Frozen CLI 2705a1e850132f74aed2dbb87149df5491bfe94b:
https://github.com/feci/parley-deck-cli/actions/runs/37698607971.
Linux PASS, Windows Build PASS / Test FAIL; macOS still running at this checkpoint.
Windows job 113056369910 raw log retained at
.parley-runtime/quota-implementation/review-cycle-5/windows-ci.log,
SHA256 5933115e4c3722a545085464ee7cf3fa87aa58b00f5546536dd79cb4400f9356,
126208 bytes. 327 FAIL entries including subtests. They are not relabeled untested.
The log includes unchanged syscall.Mkfifo compilation failure, quota directory-sync
Access-is-denied failures and other platform failures. All raw findings remain
available to independent review. Windows remains explicitly broken for the documented
quota/durable operations, experimental, and CLI winget held; this cycle makes no
Windows runtime fix. The concrete unmerged-track handoff remains authoritative.

Frozen skill 99b3f3f9fee161e8e61e585ad6c8e5b1bbbd4d2b:
https://github.com/feci/parley-deck-skill/actions/runs/37698607419.
PASS on both Python 3.10 and 3.13 CI jobs, including npm tests and package dry run.
The local full skill suite also passes. Cross-platform tests are not native Windows
CLI validation. Final macOS status will be appended before close.

Final platform state: macOS PASS too. Linux and macOS Build/Test both succeed;
Windows Build succeeds and Test fails as documented. Overall workflow conclusion
is failure because the Windows job fails. Skill CI remains PASS in both jobs.
