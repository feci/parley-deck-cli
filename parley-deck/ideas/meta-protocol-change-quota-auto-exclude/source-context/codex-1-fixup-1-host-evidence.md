---
producer: codex-1
idea: meta-protocol-change-quota-auto-exclude
phase: 8
fixup_cycle: 1
cli_product_commit: 906857b9af33306158d3b0e91b94e6ef4b6f7458
skill_commit: dcb7d593130247709cfad19715aa1f53e5cc8be0
---

# Fix-up cycle 1 host verification

Producer execution evidence only. The separate claude-1 re-review owns independent acceptance.
This supplements the child's `codex-1-fixup-1-evidence.md`; it does not rewrite its sandbox results.
All product sources stayed frozen during these commands; before/after hashes match.

## Commands and results

Host runner: `.parley-runtime/quota-implementation/fixup-1/host-checks.py`.
Raw logs, exact command/environment arrays, source hashes and summary are in that directory's
`host/`. Real providers were shadowed with version-only stubs; test fixtures supplied local commands.
No guard-denied log exists. Host execution used the normal Go cache and unweakened tests.

| Check | Command | Exit | Seconds |
| --- | --- | ---: | ---: |
| native-boot | `sysctl -n kern.boottime` | 0 | 0.007 |
| native-lease | `go test ./internal/pidlease -count=1 -v` | 0 | 0.975 |
| full | `go test ./... -count=1 -timeout 45m` | 0 | 529.055 |
| build | `go build ./...` | 0 | 1.848 |
| vet | `go vet ./internal/...` | 0 | 1.166 |
| race | `go test -race ./internal/pidlease ./internal/membership ./internal/app ./internal/consensus -run TestLease|TestQuotaFixup|TestQuotaCrossProcessLease|TestQuotaLifetimeLock|TestQuota -count=1 -v` | 0 | 18.7 |
| shared-volume | `go test ./internal/app ./internal/membership ./internal/pidlease ./internal/fsutil -run TestQuotaFixup|TestQuotaLifetimeLockDifferentRunAndOffScope|TestQuotaCrossProcessLease|TestLease|TestSyncFile -count=1 -v` | 0 | 6.579 |
| local-control | `go test ./internal/app ./internal/membership ./internal/pidlease ./internal/fsutil -run TestQuotaFixup|TestQuotaLifetimeLockDifferentRunAndOffScope|TestQuotaCrossProcessLease|TestLease|TestSyncFile -count=1 -v` | 0 | 5.657 |
| reviewer-zadv | `go run .parley-runtime/claude1-r3/zadv/main.go` | 0 | 0.358 |

The full suite passes every package. Race tests report no data race. Native boot identity is available;
`TestLeaseProcess` proves cross-process refusal and dead-local-owner takeover, both on AppleVirtIOFS
and /tmp. Unknown/foreign identity refusals remain tested. The child's failing native boot test is
therefore superseded by this actual host result, not relabeled without execution. The host full suite
also resolves the earlier sandbox budget-lock/downstream failures; no legacy accounting code changed.

The shared-volume command sets TMPDIR to the host/shared-temp directory on this actual workspace
filesystem. The local control uses `/tmp/quota-auto-exclude-fixup1-host-control`. Their test output
includes actual knob-on/off Run creation, manual membership import, owner revision, repeated receipt
recovery, signoff handoffs, pipeline/TUI cleanup, canonical-target gating and native/synthetic leases.
All selected tests pass. Source-level details and per-G mappings remain in the child's producer map.

`gofmt -l` on all 83 changed/new Go files is empty. `git diff --check` passes. The CLI product snapshot
is `906857b9af33306158d3b0e91b94e6ef4b6f7458`. The skill counterpart is `dcb7d593130247709cfad19715aa1f53e5cc8be0`; its three inputs match the previously executed full
skill suite (399 Node / 54 Python tests and all six manifests). No further skill edit occurred.
Full protocol source SHA remains `73613f95aabff8ed57372a535913985554014b71911cc8126ec3f7890a9d4f4e`;
exact skill/deck equality holds. The full suite reruns the unchanged bootstrap-drift and 70,000-byte
packet-guard tests. No normative prose/map/guard change was needed during cycle 1.

Both roster hashes equal the pre-check snapshot. No model, provider, credential, worktree declaration
or release change was made. Windows compilation passed in the child; Windows runtime is unexecuted.

## Original reviewer reproductions, unchanged

Original source programs under `.parley-runtime/claude1-r3/` were compiled unchanged. Outputs and
scratch paths are in `host/original-reviewer-probes/{shared,local}.log` and `results.json`.
These are producer reruns of reviewer inputs, not fresh independent review.

- Creation: knob false, knob true and nil-policy API all return no error on both volumes.
- Lease: run-A owns; a separate run-B process is refused on both volumes. Owner identity and path are
  printed. The original hold remains six seconds; no provider or global driver is involved.
- Knob off: the original free-form `excluded: d — unavailable; user confirmed 2026-10-04` still returns
  `manual exclusion d lacks recorded owner confirmation`. The implemented protocol form is
  `excluded: [d — unavailable — confirmed 2026-10-04]`. Its ordinary mutation/signoff/return tests pass
  without quota revise. The original new-e edit also lacks catch-up/owner proof and is refused.
  This grammar distinction is explicitly supplied to claude-1 for R2 adjudication; no concurrence is
  inferred from the documented-format test or from the probe's exit code.
- Policy-on prompt-only re-inclusion/scope edits still lack the committed owner revision and remain
  refused/reconciled from history. New owner-revision CLI and end-to-end withdrawal tests pass. The
  reviewer must evaluate whether that path meets G2/G6; this file supplies no acceptance verdict.
- Original zadv: every incomplete/adversarial row is ineligible. Complete source-derived SDK fixtures
  independently exercise positive exhaustion/reset boundaries; rejecting every old excerpt alone is
  not proof of useful native compatibility.

## Remaining evidence limit

The 24,833-byte original zcode stderr was not recovered. The two retained native excerpts are partial
and rejected by the complete framing allowlist. The positive body/reset is exercised in openly labeled
source-derived SDK fixtures. Full native AC2 verification remains missing. R4 requires it to be exposed
to the owner if independent re-review cannot resolve it; this checkpoint does not waive it.
No attended close, merge, release or independent channel verification is claimed.
