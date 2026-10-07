# Cycle-5 producer evidence — 2026-10-08

CLI 2705a1e850132f74aed2dbb87149df5491bfe94b / skill 99b3f3f9fee161e8e61e585ad6c8e5b1bbbd4d2b. Signed plan and all reservations are
in review/consensus.md (preserved under source-context/cycle5-plan-20261008/signed-plan.md).
No independent acceptance claim. Role concentration: codex-1 organizer+implementer.

G21–G25 implemented as the living IMPLEMENTATION cycle5 section describes. Shared
notice publication is moved, not reinterpreted as authority. Scoped alias checks use
the baseline root-first exact physical-scope predicate; no ancestor-name search.

Executed `go test ./internal/membership ./internal/runcontrol -run TestQuotaCycle5
-count=1 -v`: exit0 after correcting an unused import and a malformed legacy test
fixture. Both attempts retained in fixup-5/focused-{initial,corrected}.log.
The true legacy fixture is created through CreateIdeaFull; it does not delete
quota authority fields to simulate legacy. Read-only directory case executes here,
not skipped. New tests cover no writes on aliased refusal, normal/ancestor scope,
physical lease sharing and competing/nested runs, off-scope continuity, correction
receipt integrity/recovery and safe/unsafe/missing/read-only kickoff inboxes.

Executed cycle3/4/5, lifetime, cross-process, automatic kickoff and embedded drift
regressions on membership/runcontrol/protocol: all pass. Exact regex/log:
.parley-runtime/quota-implementation/fixup-5/regression-local.log.

Full host checks running at dispatch: fixup-5/host/{input-hashes,results,summary}.json
and named logs. Parent does not mutate product while they run. Required Go full,
build/vet/race/cross-build, shared/local and runner/driver checks are included.
Full skill test is running at dispatch; npm run manifest:addons completed before
freeze; source SKILL.md remains 19907 bytes. Final result will be appended here.

All normative protocol copies unchanged, deck/skill full hash 73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e.
No provider probes/captures, roster edits, new worktree or model change. AC2 remains
NOT MET / owner-waived. Windows remains knowingly broken as disclosed, no runtime
fix; explicit branch handoff and original CI log preserved. No close or release yet.

Full skill suite completed exit0: 399 Node, 54 Python and six manifests; SKILL.md 19907 bytes. Exact frozen source hashes/test-log hash are in fixup-5/skill-result.json. No assertion changed.

Candidate platform CI is recorded in source-context/codex-1-cycle5-platform-ci.md. Linux+skill PASS, Windows actual failure with full raw log/hash; macOS pending. Independent review initially failed provider502 before an artifact; same-step bounded retry is pending, no code verdict.
