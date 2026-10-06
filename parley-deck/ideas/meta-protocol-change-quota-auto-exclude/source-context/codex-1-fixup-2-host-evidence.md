---
producer: codex-1
idea: meta-protocol-change-quota-auto-exclude
phase: 8
fixup_cycle: 2
cli_product_commit: 0ee18889977a29d6baf53d3016b501112760ae12
skill_commit: e2f3649eb938e870367c76943b441382fe7acf65
---

# Fix-up cycle 2 host verification

Producer execution evidence only. A separate claude-1 process must independently re-review the
complete product and both fix-up cycles. No criterion acceptance, close or release is inferred.
Code was frozen throughout checks: all before/after product source hashes match. Both roster hashes
match the pre-check record. Provider executables were shadowed by version-only stubs; there is no
guard-denied log. Exact command/env arrays and raw logs are under `.parley-runtime/quota-implementation/
fixup-2/host/`; selected bytes and hashes are retained at `source-context/fixup-cycle-2-20261006/`.

| Check | Command | Exit | Seconds |
| --- | --- | ---: | ---: |
| native-boot | `sysctl -n kern.boottime` | 0 | 0.008 |
| native-lease | `go test ./internal/pidlease -count=1 -v` | 0 | 1.0 |
| full | `go test ./... -count=1 -timeout 45m` | 0 | 607.553 |
| build | `go build ./...` | 0 | 1.686 |
| vet | `go vet ./...` | 0 | 1.58 |
| windows-build | `go build ./...` | 0 | 3.478 |
| race | `go test -race ./internal/pidlease ./internal/membership ./internal/app ./internal/consensus ./internal/telemetry ./internal/runner -run TestLease|TestQuotaFixup|TestQuotaCrossProcessLease|TestQuotaLifetimeLock|TestQuota -count=1 -v` | 0 | 17.17 |
| shared-volume | `go test ./internal/app ./internal/membership ./internal/pidlease ./internal/fsutil -run TestQuotaCycle2|TestQuotaFixup|TestQuotaLifetimeLockDifferentRunAndOffScope|TestQuotaCrossProcessLease|TestLease|TestSyncFile -count=1 -v` | 0 | 6.792 |
| local-control | `go test ./internal/app ./internal/membership ./internal/pidlease ./internal/fsutil -run TestQuotaCycle2|TestQuotaFixup|TestQuotaLifetimeLockDifferentRunAndOffScope|TestQuotaCrossProcessLease|TestLease|TestSyncFile -count=1 -v` | 0 | 8.032 |
| reviewer-zadv | `go run .parley-runtime/claude1-r3/zadv/main.go` | 0 | 0.368 |
| reviewer-changing-countdowns | `go run .parley-runtime/quota-implementation/fixup-2/checks/reviewer-zretry/main.go` | 0 | 3.409 |

The full HOST suite passes every package in 607.553 seconds. Build, full vet, race and both filesystem
acceptance runs pass. Windows/amd64 build passes; Windows runtime is unexecuted and native crash
settlement stays unavailable there. The existing packet guard and bootstrap drift test pass unchanged
in the full suite. All 86 changed/new Go files are gofmt-clean; git diff --check passes.

The host also executed the exact native test that failed in the child sandbox: `go test
./internal/membership -run '^TestQuotaCycle2NativeCrashWriter$' -count=1 -v`, exit 0 in 0.723 seconds.
It reports an idempotent settlement for a real local-stub crash. The subsequent shared-volume,
local-control and race runs include this test and pass. Native boot and stale-lease tests pass.
Thus those specific sandbox failures have executed HOST results; their failed raw child logs remain
intact. The full suite also runs the earlier budget-cache/downstream tests on the host.

The original reviewer changing-countdown program was copied with only its output locator changed,
as recorded in prepare-differential.py. The host result accepts B-final/B2/C-final/D/E and rejects
the future-observation B-first/C-first cases. The 54 original adversarial cases stay ineligible.
These are producer reruns of reviewer inputs, not new independent acceptance. The child additionally
executed the actual pre-change 27e42b8 CLI versus current app.Run for the policy-off return/join path,
on /Volumes and /tmp; its exact baseline and probe hashes are retained.

The final full skill suite passes 399 Node tests, 54 Python tests across seven files, and all six
manifests, in 119.983 seconds. Final input/output hashes match. The core is 19819 bytes below
its unchanged 20000-byte guard. No skill test is added or weakened. Source protocol and full
phase-0/5/8 packet hashes remain `73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e`;
exact skill/deck equality holds. Command guidance lives in CLI docs/skill; no new COOPERATION text
change occurred in cycle 2.

## Remaining acceptance boundary

AC2 remains OPEN. Complete allowlisted source-derived records qualify; partial native tails and
default-console [Object] truncation reject. The installed-source/extracted-SDK investigation is
recorded, but the complete original native capture is absent. No excluded zcode/provider invocation
or parser relaxation occurred. The owner decision follows independent review of this concrete
behavior. The source-derived positives are never relabeled native. Both final review-consensus
signoffs and a new attended owner close are still required before merge or release.
