# Cycle-4 host verification — codex-1 producer evidence

Verified 2026-10-07T22:15:11.117510+00:00. CLI product e04852ef36846b8ec18f581c4e8ae4f28b0c2262,
skill b9596ddd5c37633026a26031f1bccaee76d2622f. This report is producer execution,
never independent acceptance or a close. Source and roster hashes stayed unchanged.

| Check | Exit | Seconds |
| --- | ---: | ---: |
| native-boot | 0 | 0.008 |
| native-lease | 0 | 4.808 |
| full | 0 | 947.455 |
| build | 0 | 3.026 |
| binary | 0 | 1.15 |
| vet | 0 | 1.321 |
| windows-build | 0 | 2.154 |
| race | 0 | 33.907 |
| shared-volume | 0 | 27.127 |
| local-control | 0 | 22.722 |
| shared-runner-driver | 0 | 1.481 |
| reviewer-zadv | 0 | 0.348 |
| reviewer-changing-countdowns | 0 | 3.275 |

All 13 host checks pass. Full Go suite uses `go test ./... -count=1 -timeout 45m`.
Full build and vet pass; Windows/amd64 cross-build passes (not runtime verification).
Race covers pidlease, membership, app, consensus, quota, protocol, telemetry and runner.
Affected tests, including the unchanged native crash writer, pass on both the actual
shared volume and /tmp. The explicit shared runner/driver repeat passes, resolving
the sandbox's killed invocation on the HOST; its original unknown-cause failure is
retained. Native boot and lease checks pass without changing any assertion.

Formatting: 97 changed Go files, none unformatted. Skill full suite passes
399 Node tests, 54 Python tests, six payload manifests; SKILL.md 19907 bytes.
Skill hashes stayed unchanged. Initial stale-manifest failure remains under
fixup-4/skill/before-manifest-*; regenerating the payload manifest was the correction.
No real provider command ran: local provider guards used, deny log absent.

Full packets phases0/5/8 match source/packet SHA256
73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e, no fallback.
Deck/skill protocol bytes are identical; embedded bootstrap passes the existing
normalization drift test in the full suite, with no assertion or guard change.
No normative protocol hunk changed in cycle4.

Candidate CI: Linux and macOS PASS; Windows BUILD passes but TEST fails, including
quota kickoff directory-sync Access is denied and an unchanged syscall.Mkfifo test.
This is executed failing Windows evidence, not a runtime pass or merely unexecuted.
Windows remains experimental and CLI winget held as the owner brief specifies. See
codex-1-cycle4-platform-ci.md and retained job metadata. Skill candidate CI passes.

Native AC2 is NOT MET and explicitly owner-waived for this release, never PASS.
R5-MAJOR-2 remains accepted/deferred. Container/PID namespace behavior is unverified.
Review/round-07/claude-1.md remains the independent verdict; it is still in progress.
Exact evidence copies/hashes: fixup-cycle-4-20261008/manifest.json. Private command
logs remain at .parley-runtime/quota-implementation/fixup-4/host/<check>.log.
