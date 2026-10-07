# Cycle-4 platform CI evidence

Candidate CLI commit e04852ef36846b8ec18f581c4e8ae4f28b0c2262.
Run https://github.com/feci/parley-deck-cli/actions/runs/37692987015
Windows job 113037539987: failure. Build success; test fails.
Leading compiler error: internal/evidence/tree_report_test.go:97:20 undefined syscall.Mkfifo.
That file is unchanged against FINAL 27e42b8. Other Windows runtime/fixture failures exist;
this is not a claim that all failures are pre-existing or that Windows passed.
Cycle4-specific FAIL entries printed in the log: 33 (absence is not execution proof).
Raw log SHA256 b36a9e47732dd37de0d1c637bfe9d28004dce4fcc396545674b4000da609de2d, bytes 112905, at
.parley-runtime/quota-implementation/review-cycle-4/windows-ci.log. Full job metadata beside it.
Logs obtained from the job API before overall run completion; terminal escapes stripped
for safe plain-text storage. CLI Windows remains experimental and CLI winget is held
under the controlling owner brief. Independent reviewer should assess this evidence.
Linux candidate CI passed; macOS is still running at this observation. Skill CI passed
https://github.com/feci/parley-deck-skill/actions/runs/37692990193 at b9596dd.

Windows cycle4 membership fixtures also fail before the notice checks because directory
sync at quota kickoff returns `Access is denied` (for example cycle4_test.go:56 in
TestQuotaCycle4OwnerNoticeCopies/applied/live/annotation). There are 33 cycle4 FAIL
entries including subtests. Thus Windows has executed failures, not merely absent
evidence. Build succeeds; runtime support remains experimental, with CLI winget held.
No claim is made that these failures are all pre-existing.

Final candidate CI: Linux SUCCESS, macOS SUCCESS, Windows FAILURE; overall run FAILURE because of Windows. Windows Build SUCCESS/Test FAILURE. This is the recorded experimental platform limit, not an all-platform green claim.
