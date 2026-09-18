---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-11
artifact-kind: implementer dispositions and executed evidence
responding-to: claude-1-budget-tty-followup-20260911.md
not-a-signoff: true
---

# Budget foundation and residual TTY review dispositions

Claude's own partial source review at 9cd1e42 is preserved unchanged. These
are Codex's corrections, positions and executed tests, not Claude's agreement
or full Phase-6/7/8 acceptance. Full production launch/action integration is
still missing; no budget acceptance criterion is claimed complete.

## Budget findings

- B-MAJOR-1: an immutable `lock-origin` is published atomically before taking
  a local lock, pinning hostname and resolved local-cache lock path. A second
  environment refuses if it differs; first-writer bootstrap is tested under
  competing origins. Unix uses a synchronized staged file plus hardlink and
  directory barrier; Windows uses non-replacing, write-through MoveFileEx.
  There is no non-atomic publication fallback. This remains one pinned origin,
  not cross-host distributed locking or a promise about hostnames as identity.
- B-MAJOR-2: absolute-path resolution precedes symlink resolution. A relative
  path under a symlinked working directory and the absolute physical path
  share one lock in the executed fixture.
- B-MAJOR-3: `TestNoOpKernelLockIsRefused` injects an always-successful lock
  primitive and requires refusal. Separate-process tests admit exactly five
  of twelve attempts under the action cap and separately under the cost cap.
  These tests also pass with TMPDIR on the actual shared volume.
- B-MAJOR-4: `ReconcileUnknown` and the attended `parley budget reconcile`
  control append a distinct, idempotent operator decision with a conservative
  monetary ceiling and reason. They retain the original unknown cost and
  spent action. A conflicting decision replay or attempt to replace a known
  observation refuses. The exact ledger/scope/action are explicit CLI inputs;
  participant frontmatter is never a grant. Tests require unattended refusal,
  unchanged state on refusal, scope validation, exact replay, retained unknown
  observation and retained replay refusal after recovery. This is recovery
  infrastructure; the launch policy/callers remain open.
- B-MINOR-1/2: `Limits.Denied` expresses refusal independently of historical
  zero/unlimited semantics; `Actions` is documented as lifetime charged totals.
- B-MINOR-3: wall time explicitly includes pause/resume; backward movement has
  a distinct `ErrClockSkew`. It still fails closed rather than inventing a clock
  correction. `Inspect` remains available without that mutation-time check.
- B-MINOR-4: arithmetic overflow has `ErrCostOverflow`, separate from unknown
  observation. Both stop a monetary ceiling.
- B-MINOR-5: unsupported/mixed schema versions still fail closed. The current
  implementation is explicitly not certified for mixed-version writers.
- B-MINOR-6: validation now walks decoder tokens with an eight-level nesting
  bound instead of copying a RawMessage at each nested level.
- B-MINOR-7: remains an explicit scalability limitation: 16 MiB ceiling, full
  file rewrites, no compaction/draining protocol that forgets replay identities.
  No claim that arbitrarily large ledgers remain writable.
- B-MINOR-8: store operations impose a 30-second lock-contention deadline when
  the caller has no shorter deadline. This does not bound filesystem syscalls
  and does not claim fairness.
- B-MINOR-9: failed exclusion names the resolved local lock path.
- B-MINOR-10: conservative case folding applies on all OSes, including Linux
  with case-insensitive mounts. It can overlock distinct case-sensitive paths.
- NITs: no builtin `cap` shadow; lock release is idempotent. Returned snapshots
  have no retained internal mutable owner. The per-operation exclusion probe
  is intentionally retained. Cache removal or migration during use is not
  supported; neither the live lock inode nor its origin may be deleted.

## Terminal findings

- NEW-1: wrong launch mode now follows the recorded refusal path and has an
  explicit reason. A new fixture requires a failed terminal record with no
  actual spawn. Its initial test run failed because the expected-reason table
  omitted the newly added scenario; the corrected table then passed.
- NIT-A: all four real-PTY failure cases now pass a `/dev/tty` descriptor above
  fd 2, while verifying restored parent input and no surviving descendant.
- NIT-B: pre-render delivery/mode refusals explicitly mark terminal streams
  unobserved, with null counters.
- MINOR-6: handoff instructions distinguish configured invocation from actual
  execution and explicitly state that automatic spawn-tty currently exists
  only on consensus request-signoffs. Other surfaces are print-only. This
  fixes misleading capability text; it does not add other spawn surfaces.
- NEW-2: the conservative argument bound and file-mode preference remain.
- Open question 3: retain nonzero-exit failure even when an artifact exists.
  Process success and artifact existence are different observations; an
  abnormal process outcome cannot silently become success. The artifact is
  preserved for reconciliation, and this position does not forge acceptance.

## Executed validation

- `go test -count=1 ./internal/budget`: PASS, including cross-process caps,
  no-op-lock rejection, competing origins, alias exclusion, operator recovery,
  durable write failure, clock/overflow and bounded parser cases.
- Same budget suite with TMPDIR under integration's
  `.parley-runtime/budget-review-tmp`: PASS (1.137s). An initial command used
  a nonexistent temporary directory and failed before Go created a workdir;
  the explicitly created directory above is the actual shared-volume run.
- `go test -race -count=1 ./internal/budget`: PASS (4.729s).
- `go vet ./internal/budget ./internal/runner`: PASS.
- Windows budget test binary cross-compiles; runtime remains untested.
- `go test -count=1 ./internal/runner -run 'Interactive|Handoff'`: PASS (3.222s).
- `python3 scripts/test_terminal_launch.py`: PASS, one child and five parent
  input handshakes; all four failure cases use nonzero parent terminal fds.
- Full `go test ./...` passed after foundation/TTY corrections, logged in
  `.parley-runtime/budget-review-fixes-full-suite-20260911.log`. The later
  operator CLI/Inspect additions have a separate targeted app/budget PASS;
  they were not part of that already-running full-suite snapshot.

## Still required

Independent re-evaluation of these changes; caller policy shared across
launch/manual/driver/resume/BLOCK; legacy-charge migration without resets;
enforced action/cost/time ceilings and actual operator extensions; complete
readiness and criterion execution; exact experiments and full review consensus.
No merge, release, global installation, quorum change or pilot amendment.
