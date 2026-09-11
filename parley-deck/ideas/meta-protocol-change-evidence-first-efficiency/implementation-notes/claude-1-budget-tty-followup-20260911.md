---
agent: claude-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-11
type: partial-source-review
reviewed-commit: 9cd1e42
also-anchored: internal/budget/** introduced at 4bd11b7, read at current HEAD
scope: internal/budget/{ledger.go,lock.go,lock_unix.go,lock_windows.go,ledger_test.go},
  internal/runner/launch.go (RunInteractive + ValidateInteractiveDelivery),
  internal/runner/launch_test.go, internal/runner/handoff.go, internal/telemetry/record.go,
  internal/telemetry/usage.go, internal/procctl/terminal_unix.go,
  internal/app/consensus_request_signoffs.go, internal/fsutil/{replace_unix.go,replace_windows.go,fsutil.go},
  scripts/test_terminal_launch.py, docs/agent-runtime-configuration.md
signoff: none
supersedes: nothing — claude-1-terminal-review-20260911.md stands unchanged
---

# Partial source review — terminal corrections re-evaluation + budget foundation

## 0. What this is and is not

A **partial source review of two slices**. It is **not** Phase-6 full-scope acceptance,
not a review-consensus signoff, not an AC closure, and not a signature of any kind. My
prior note `claude-1-terminal-review-20260911.md` and every earlier artifact of mine are
**unchanged**; nothing here retracts or amends them.

Per the launch brief: no source edit, no git mutation, no browser, no global config, no
other agent invoked. **Shell was unavailable to me — I executed nothing.** I claim no
test execution, no build, and no runtime observation.

Also per the brief and explicitly restated here: **a signoff-process failure remains a
failure even when an artifact like this one exists.** This file is evidence for a future
review round; it is not that round.

## 1. Provenance (§15.2)

- **PRIMARY** — every claim below about *what the code says* comes from reading the named
  files at the stated line locators in `worktrees/evidence-first-integration`. The
  negative claim "`internal/budget` has zero non-test callers" comes from a repo-wide
  regex over all `*.go` for `internal/budget|budget\.(Store|Reserve|Settle|Limits|Request)`
  → **no matches**; a caller reached by reflection or codegen would escape it.
- **Environment/brief-reported, not verified by me** — that this tree *is* commit `9cd1e42`,
  and that `internal/budget/**` was introduced at `4bd11b7`. I ran no git command. Re-anchor
  the locators if the tree has moved.
- **Facilitator-reported execution, explicitly not mine** — full `go test ./...` PASS after
  the terminal corrections; targeted race + vet PASS; Windows runner+budget cross-compile
  PASS with Windows runtime untested; five real-PTY scenarios PASS (including a nonzero
  parent `/dev/tty` fd, an actual timeout, and start/terminal evidence-write failures with
  descendant cleanup); budget tests PASS on ordinary local temp, on a shared mount, under
  race, and with a killed subprocess lock owner. **These are reported results.** Where a
  verdict of mine depends on one, I say so and tag it SECONDARY on the facilitator.
- **SECONDARY (on codex-1)** — the `Ctty` parent-fd semantics, cited by the owner to Go's
  local `syscall/exec_libc.go:30-34`. I did not read the Go sources (outside permitted
  scope). See §2.4.
- **RECALL (UNVERIFIED)** — three background facts remain memory-only and are tagged where
  used: (a) POSIX `flock` locks attach to the *open file description*, so two separate
  `open()`s in one process do conflict (this is what makes the new probe meaningful);
  (b) Windows `LockFileEx` byte-range locks also conflict across handles within one
  process; (c) Linux `MAX_ARG_STRLEN` caps a single argv element at 128 KiB. (a) is
  corroborated — not proven — by `TestConcurrentReservationsCannotOverspend` passing
  in-process, as facilitator-reported. (b) is **uncorroborated**; Windows runtime is
  untested by anyone.

---

## 2. Part 1 — my terminal findings re-evaluated against the corrections

I checked each correction against source rather than against the disposition text. Verdicts
below are mine.

### 2.1 Corrections I concur are closed

| Prior finding | Verdict | PRIMARY evidence |
| --- | --- | --- |
| **MAJOR-1** unobserved streams serialized as measured zeros | **Closed. Concur.** | `record.go:63-64` both counters are now `*int64`; `telemetry.go:112` nils them on the direct-terminal path; `usage.go:285-286` and `record.go:296-297` clone through; `launch_test.go:68` asserts both nil for the terminal record and `handoff_test.go:28` asserts both nil for the handoff record. The handoff case — the one I flagged as worse because it carries no `stream_coverage` — is now explicitly covered. |
| **MAJOR-2** no fixture exercises cleanup/restore on the real-terminal branch | **Closed. Concur.** | `launch_test.go:172-219` `TestInteractiveRealTerminalFailures` runs four scenarios (`timeout`, `start-write`, `failed-start`, `terminal-write`) under the PTY harness, each asserting **both** halves I asked for: the parent can still read the terminal (`:207-212`) and the descendant did not survive (`:213-216`). The `timeout` case is exactly the transition I called riskiest — killing a child that owns the foreground, then re-taking it. `scripts/test_terminal_launch.py:51` cross-checks that all 5 interactions actually ran (`sent_child != 1 or sent_parent != 5`), so a silently-skipped scenario fails the harness rather than passing quietly. That last detail is what makes this a real fixture and not a smoke test. |
| **MAJOR-3** whole protocol in one argv element | **Closed as scoped. Concur.** | `launch.go:309-313` bounds each expanded element at `120<<10` and returns `interactive-prompt-argument-too-large-use-file` — a `*protocolContextError`, so it is classified and recorded, not an opaque `E2BIG`. `launch_test.go:221-250` asserts both the bound and `StartedAt == nil` (refused before spawn). `docs/agent-runtime-configuration.md:123-125` now says "Prefer file mode: argument mode exposes the full protocol and task through the process command line". The owner's framing — a conservative per-element bound, not a universal argv-capacity claim, with total argv/env overflow still landing as a recorded failed start — is accurate and matches `launch.go:309-310`. See NEW-2 for a residual operational note. |
| **MINOR-1** no launch-mode guard | **Closed as a mislabel risk. Concur** (but see NEW-1). | `launch.go:262-264`. |
| **MINOR-2** restoration error discarded | **Closed. Concur.** | `launch.go:327` now `fmt.Errorf("cannot restore terminal foreground process group: %w", err)`. |
| **MINOR-3** delivery contract not validated at selection | **Closed, and better than I proposed. Concur.** | `ValidateInteractiveDelivery` is extracted and exported (`launch.go:239-251`), called at selection (`consensus_request_signoffs.go:359-360`) *and* at the spawn gate (`launch.go:266`) — one predicate, two call sites, so they cannot drift. The orphaned-trail half is also fixed: `consensus_request_signoffs.go:490` appends `agent.handoff.failed`. |
| **MINOR-4** timeout budget applied twice | **Closed. Concur.** | The poll loop now selects on the **same** `agentCtx` that bounded the spawn (`consensus_request_signoffs.go:511`), not a fresh deadline. One session budget. |
| **MINOR-5** handoff prompts in a tracked path | **Closed. Concur.** | `handoff.go:74` writes `PromptPath` into `evidence.invocation.Dir` (under the gitignored `.parley-runtime/`), and `handoff_test.go:51` asserts the old `parley-deck/runs/.../handoff-prompt.md` path is **absent** — a negative assertion, which is the right shape for a "must not be written here" rule. |
| NIT: `SysProcAttr` replaced wholesale | **Closed. Concur.** | `terminal_unix.go:28-32` now allocates only when nil and sets exactly the four fields it owns, explicitly zeroing `Setsid`/`Setctty` rather than assuming them. |
| NIT: no `WaitDelay` | **Closed. Concur.** | `launch.go:317` `cmd.WaitDelay = 2 * time.Second`. |
| NIT: docs silent on arg-mode exposure | **Closed. Concur.** | `docs/agent-runtime-configuration.md:123-125`. |

### 2.2 Dispositions I concur with that are *retentions*, not fixes

- **Constant `/bin/sh` restoration dependency, explicitly retained.** I concur with
  retaining it. `terminal_unix.go:40` runs a compile-time-constant argv (`-c "exit 0"`)
  with no task text and no captured output, so it is correctly not agent telemetry; the
  alternative the comment rejects (mutating the parent's process-wide signal disposition,
  or stopping its background group) is strictly worse. `/bin/sh` is POSIX-mandated. With
  MINOR-2 fixed, a missing or failing `/bin/sh` is now *diagnosable* rather than a bare
  "cannot restore", which was most of my objection. **No residual finding.**
- **Windows is a no-op with runtime untested.** Recorded, not resolved. Unchanged from my
  prior note; still correct to carry as a disclosed gap rather than a claim.

### 2.3 Findings that remain OPEN (I do not concur that these are closed; the owner does not claim they are)

- **MINOR-6 — `spawn-tty` is honoured on exactly one launch surface, silently.** Not in the
  MINOR-1..5 correction set, and `IMPLEMENTATION.md:22-23` independently records "complete
  launch coverage ... remain open". I concur the gap is **disclosed**; I do not agree it is
  resolved. It still bounds what this slice can contribute to AC-T1, and the specific
  sharp edge stands: the handoff text advertises `Invoke: spawn-tty` on packets that no
  surface will spawn. Disclosure is not a fix; this should stay on the open list until
  either the other surfaces are wired or validation warns.
- **Open question 3 — a non-zero exit after a *valid* signoff append still aborts the whole
  batch** (`consensus_request_signoffs.go:489-494`: a `runInteractiveTTY` error returns
  before the poll loop can observe the appended signoff). Not addressed in the
  dispositions. With a human terminal session, exiting non-zero is an ordinary way to end,
  so inherited headless behaviour is now reachable in a much more likely way. Still wants
  an explicit decision.

### 2.4 Open question 1 (`Ctty` parent-fd semantics) — resolved for the success path

The owner cites Go's `syscall/exec_libc.go:30-34` (owner-PRIMARY; **SECONDARY for me**, I
did not read it) *and* supplies independent empirical corroboration I can verify by
inspection: `launch_test.go:149` opens `/dev/tty` directly, which is necessarily a fd
above 2, and passes it as `input` — so the passing PTY fixture now *does* discriminate
between the parent-fd and child-fd readings, which my prior note said it could not. That
is the right kind of evidence and I accept it **for the success path**.

**Residual:** `TestInteractiveRealTerminalFailures:202` passes `os.Stdin` (fd 0). So the
*restore-after-failure* path — the one that also passes `Ctty: int(input.Fd())` into the
`/bin/sh` shim at `terminal_unix.go:42` — is still only exercised at fd 0. Small, and I do
not think it is wrong; naming it so the coverage claim stays exact. (NIT-A below.)

### 2.5 New findings in the corrected terminal code

- **[MINOR] NEW-1 — the new launch-mode guard is the one refusal that produces no record.**
  `launch.go:262-264` returns a plain error *before* `beginLaunch`, so a wrong-`launch_mode`
  call emits **zero** telemetry records. Every other exit in this function is accounted:
  the delivery refusal explicitly begins and finishes evidence (`:267-277`), the
  nil-descriptor check (`:294`) sits after the deferred `finish` at `:285`, and so does the
  oversized-argument refusal. The function's stated invariant — which my prior refutation
  #9 verified and which `launch_test.go:73-121` fixtures — is one terminal record per
  attempt. This guard is a latent path today (the only caller is mode-gated upstream), so
  this is a uniformity defect, not a live accounting hole. Fix: route it through the same
  `beginLaunch`/`finish` shape as the delivery refusal, with its own reason string.
- **[NIT] NEW-2 — arg mode now has ~14 KB of headroom, not "plenty".** The bound is 120 KiB
  = 122,880 bytes; this session's own shadow audit reports the live protocol source at
  108,400 bytes, and `prepared` is attestation + protocol + task. So a correctly configured
  arg-mode agent is already within roughly 13% of the ceiling and will start refusing after
  modest protocol growth. The refusal is clean, named and pre-spawn — this is working as
  designed, and the fail direction is right. It is an *operational* note: arg mode should be
  understood as near-term unusable for full-protocol launches, which strengthens the docs'
  new "prefer file mode" guidance rather than contradicting it.
- **[NIT] NIT-A** — restore-path `Ctty` only exercised at fd 0 (§2.4).
- **[NIT] NIT-B** — the refusal paths still omit `stream_coverage` rather than labelling it
  `not-observed-terminal` (`directTerminal` is set at `launch.go:283`, after both refusal
  returns). This was a NIT in my prior note and is now *materially* harmless, because
  MAJOR-1's fix means those records carry `null` byte counters rather than zeros. Noting
  the labelling asymmetry survives; I would not hold anything for it.
- **Scope note, not a finding** — `runner.go:596-597`, `acp.go:162-163`, `phase58.go:90`,
  `consult.go:129` still emit plain `int64` byte counts from `supervision.Snapshot`
  (`supervision.go:106-107`). Those are captured paths where a count is genuinely observed,
  so nothing is misstated. But a consumer that aggregates **run-store event payloads**
  rather than telemetry records does not get MAJOR-1's nullability. If the ≥20-attempt
  reconciliation reads events rather than records, check which surface it uses.

---

## 3. Part 2 — `internal/budget/**` (new, read at HEAD; introduced at 4bd11b7)

### 3.1 What I verified the code actually does

The owner's framing is **accurate and I confirm it by inspection**:

- **Precharge before permission.** `Reserve` → `update` (`ledger.go:87`, `:170-219`) writes
  the ledger via `writeSynced` → `fsutil.ReplaceSyncedFile` **before** returning. On unix
  that is temp-write → `SyncFile` → `rename` → **parent-directory fsync**
  (`fsutil/replace_unix.go:19-31`). So the charge is durable before the caller is
  authorized. `TestPersistFailureNeverAuthorizesWork:103-128` pins both halves, including
  the nastier one: a *replaced-then-errored* write must keep the conservative charge.
- **Replay is not a second permission.** `ledger.go:88-90` → `ErrReserved` on a repeated ID.
- **Failed/crashed reservations stay charged.** No refund path exists; `count`
  (`:95-100`) includes settled entries, so caps are lifetime totals.
- **Known vs unknown exposure is reconciled, never coerced to zero.** `Exposure`
  (`:155-168`) returns `(_, false)` if any effective amount is nil, and `Reserve`
  (`:104-115`) refuses with `ErrUnknownCost` rather than assuming.
- **Bounded/validated persisted state.** `read` (`:221-264`) rejects non-regular files,
  >16 MiB, TOCTOU inode swaps (`os.SameFile`), unknown fields, trailing data, duplicate and
  case-aliased keys, nulls outside the two monetary fields, bad hex ids, invalid kinds,
  `reserved_at` before `started_at`, negatives, and `actual_micros` set without `settled`.
  `TestCorruptPersistedChargeFailsClosed:143-162` covers eight of those.
- **Synchronized replacement + host-local cache kernel locks keyed by canonical ledger dir.**
  `lock.go:17-88`. Kernel ownership means a crashed holder is released by the OS with no
  PID file and no stale-lock race — pinned by `TestKernelLockSurvivesProcessDeath:164-199`.
- **The two-handle probe.** `lock.go:61-76` opens a *second* descriptor and re-locks; if the
  second attempt succeeds, exclusion is not real and `lock` fails closed. Given RECALL(a),
  this is a sound same-host detector for the exact no-op-flock failure that motivated it.

**It is genuinely unwired.** Repo-wide search for `internal/budget` / `budget.Store` etc.
returns **no matches** outside the package. The owner's "foundation rather than a completed
budget gate" (`IMPLEMENTATION.md:629-631`) is exactly right, and I will not claim otherwise.

### 3.2 Refutation attempts (LE-1) — inspection only, no execution

1. *"A crash cannot produce a free retry."* Tried to find a path where a charge is written
   but not durable, or durable but not returned as `ErrReserved`. The fsync-rename-fsync-dir
   chain and the `replaced=true` fixture close it. **Could not break it.**
2. *"Cost cannot be silently assumed zero."* Tried to get a nil amount to contribute 0 to the
   total. `Exposure` returns `known=false` on the first nil (`:162`). **Could not break it.**
3. *"Overflow cannot admit a reservation."* `*n > math.MaxInt64-total` (`:162`), and the
   admission test `*req.ReserveMicros > limits.CostMicros-exposure` (`:112`) runs only after
   `exposure <= limits.CostMicros`, so the subtraction cannot underflow. **Could not break it.**
4. *"Two writers on one host cannot both win."* **Broken — see B-MAJOR-1 and B-MAJOR-2.**
   The lock key is derived from `os.UserCacheDir()` and from a canonicalization whose step
   order leaves a symlinked cwd unresolved.
5. *"Corrupt state fails closed."* Tried `null`, `{}`, trailing data, dup key, case alias,
   negative, null-`settled`, unknown field — all fixtured. Tried to smuggle a value past
   `DisallowUnknownFields` via duplicate keys: `checkJSON` (`:268-317`) runs first and
   rejects them. **Could not break correctness** — but see B-MINOR-6 on unbounded recursion.
6. *"A settled action cannot be re-priced."* `:133-137` rejects a conflicting terminal cost
   and accepts an identical repeat idempotently. **Could not break it.** (Its cost is
   B-MAJOR-4: there is then no legitimate correction path either.)
7. *"The store enforces a budget."* **Broken by construction** — `ReserveMicros` is
   caller-asserted and `Settle` is caller-reported; nothing here observes real spend. This
   is an accounting ledger with a lock, which is the correct scope for a foundation, but it
   must not be described as a production budget enforcement.

### 3.3 Findings

#### [MAJOR] B-MAJOR-1 — the lock is **user- and environment-scoped**, not host-scoped as documented

`lock.go:29` keys the lock directory off `os.UserCacheDir()`. Two processes on the **same
host** reserving against the **same ledger directory** under different users — or merely
with a different `HOME` / `XDG_CACHE_HOME` — resolve **different lock files** and do not
exclude each other. Both proceed; the last `ReplaceSyncedFile` wins and the other's charges
are lost. There is no detection and no error: the failure is **silent and indistinguishable
from correct operation**.

This is not hypothetical for this project's own shape: a `sudo`-elevated run, a
launchd/systemd unit, a container that mounts the shared volume but not the user cache, or
a CI shell with an overridden cache dir all produce it. The doc comment at `lock.go:13-16`
says "This coordinates runtimes on one host, not distributed writers" — the delivered
guarantee is *one host, one user, one cache environment*.

Because "a repeated ID is never permission to repeat work" is the package's entire value,
a silent loss of exclusion defeats the primary invariant rather than degrading it.

Fix directions (any one closes it): derive the lock directory from a fixed host-wide path
instead of user env; or record the resolved lock path inside the ledger and refuse when a
second writer resolves a different one — turning a silent divergence into a loud one; or,
at minimum, narrow the documented claim to match the code and state the constraint in
caller-facing docs.

#### [MAJOR] B-MAJOR-2 — canonicalization order can yield two lock keys for one ledger

`lock.go:18-25` does `EvalSymlinks(filepath.Dir(path))` **then** `filepath.Abs(...)`. For a
**relative** `Store.Dir`, `EvalSymlinks` resolves only the components it is handed; the cwd
prefix is prepended afterwards by `Abs` and is never symlink-resolved. Two processes that
reach the same directory via differently-symlinked cwds compute different `canonical`
strings → different `key()` → different lock files → no exclusion. Same silent shape as
B-MAJOR-1.

This is live in this very workspace layout: `/Volumes/My Shared Files/...` worktrees and, on
darwin, `/tmp` → `/private/tmp` and `/var` → `/private/var`. The correct order is **`Abs`
first, then `EvalSymlinks`** — a one-line change.

The committed tests cannot see it: `testStore` (`ledger_test.go:18`) always passes
`t.TempDir()`, an absolute path, so both the parent and the subprocess fixture agree by
construction.

#### [MAJOR] B-MAJOR-3 — the shared-volume regression that motivated the whole design has no committed test

Per the launch brief and `IMPLEMENTATION.md:216`, the host-local cache **and** the two-handle
probe exist because an actual shared-volume test showed `flock` succeeding without
exclusion — **23 admissions against a five-call cap**. That is the single most important
behavioural fact in this package.

`internal/budget/ledger_test.go` contains **no** test for it. Nothing exercises a
no-op-`flock` filesystem, and nothing asserts that the probe at `lock.go:61-76` rejects one.
The facilitator's shared-mount pass was a manual `TMPDIR` run, not a committed regression.

Consequence: a future refactor that moves the lock back onto the ledger directory, or that
deletes the probe as redundant, **reproduces the original 23-admissions bug with a fully
green suite**. For a durable-charge invariant, that is the wrong safety property.

Suggested fix: make `tryLock` injectable in tests and assert that a
always-succeeds implementation drives `lock()` to the "does not provide verified exclusion"
error. That is a pure unit test, needs no shared mount, and pins the design decision where
a reader will find it. Optionally add an env-gated test against an operator-supplied path.

#### [MAJOR] B-MAJOR-4 — one unknown-cost entry permanently wedges cost-limited reservations, and there is no reconcile API

`Reserve` refuses with `ErrUnknownCost` whenever `Exposure()` is unknown (`ledger.go:106-111`),
and `Exposure()` is unknown if **any single entry** has a nil effective amount (`:162`).
An entry acquires a nil amount in two ordinary ways:

1. It was reserved while `limits.CostMicros == 0` — i.e. **no cost limit configured**, the
   natural default — because the `ReserveMicros != nil` requirement is gated behind
   `CostMicros > 0` (`:104-106`).
2. It was settled with a nil actual — the documented "unknown terminal cost remains
   unknown" (`:139`). A **crashed** action never settles at all and keeps whatever it
   reserved, so a crash under (1) is permanent.

`TestUnknownExposureAndSettlement:95-100` asserts this as intended, and as a *fail-closed*
rule it is right. The problem is that it is **unrecoverable**: the package exposes no
`Reconcile`, no `Abandon`, no operator-priced override. Once poisoned, the only remedy is
deleting `ledger.json` — which also discards every replay-refusal charge, i.e. the recovery
for "I can no longer compute a budget" is "surrender the duplicate-execution protection".

The trap for the wiring work is specific and easy to hit: **a caller that runs with no cost
limit silently poisons cost accounting for any later cost-limited caller in the same scope.**
Fail-closed and unrecoverable are different properties, and §12.7's reconcile duty is
currently unmet by this package.

Suggested fix: an explicit `Reconcile(ctx, id, micros, reason)` (and/or `Abandon`) that
records an operator decision *in the entry* so the audit trail survives; plus a distinct
sentinel so a caller can distinguish "over budget" from "budget no longer computable".

#### [MINOR] Cost/count semantics that constrain callers

- **B-MINOR-1 — a cap of `0` means unlimited, not forbidden.** `ledger.go:101` (`cap > 0`),
  `:92`, `:104`. With `n < 0` rejected at `:83-85`, `Limits` **cannot express "deny this
  kind"**. A BLOCK path that must refuse an action outright has no way to say so.
- **B-MINOR-2 — action caps are lifetime totals, not concurrency.** `count` at `:95-100`
  includes settled entries. Correct for precharge semantics and for the 23-admissions bug,
  but a caller wanting "N in flight" will get something materially different. Document on
  `Limits.Actions`.
- **B-MINOR-3 — wall clock runs from *ledger creation*, in wall time.** `StartedAt` is set
  on first use (`:189`) and compared at `:92`. Reusing a `Dir` across runs silently inherits
  the older origin. Worse, `now.Before(state.StartedAt)` (`:193`) and
  `now.Before(entry.ReservedAt)` (`:197`) make any **backward clock step** (NTP correction,
  VM snapshot restore) a hard total failure of `Reserve` *and* `Settle` — so the store
  cannot even record a settlement to unwedge itself. `time.Time` is persisted, so no
  monotonic reading is available and the check is defensible; it should at least produce an
  error that names clock skew, since "budget scope or clock mismatch" (`:194`) will send an
  operator to the wrong place.
- **B-MINOR-4 — overflow is reported as "unknown".** `:162` returns `(0,false)` for both a
  nil amount and an int64 overflow, so an arithmetic overflow surfaces as `ErrUnknownCost`.
  Both fail closed; the diagnosis is wrong and, given B-MAJOR-4, sends the operator down a
  wedge-recovery path for an arithmetic problem.
- **B-MINOR-5 — no forward compatibility, including for draining.** `read` rejects
  `Schema != 1` (`:254`) and any kind outside `validKind` (`:259`). A newer binary that adds
  a `Kind` writes a ledger an older binary cannot read **at all** — it cannot even settle
  outstanding charges. Fine for a single-version deployment; a real constraint for mixed
  versions on one shared volume, which is precisely this project's topology.
- **B-MINOR-6 — `checkJSON` recurses with no depth bound.** `:295-299` recurses on every
  nested object across a file of up to 16 MiB, allocating a fresh `json.Decoder` and a
  `RawMessage` copy per level, *before* any schema check runs. Deeply nested input is
  O(depth) stack and O(depth × size) copies. The input is a local, trusted file so this is
  availability-only — but the real schema needs depth 3, and a bound is one line.
- **B-MINOR-7 — unbounded entry growth ends in a two-way wedge.** The 16 MiB ceiling
  (`:208-210`, `:226`) is the only bound on entry count. When it trips, `Reserve` **and**
  `Settle` both fail, so the ledger cannot be drained below the limit. Every operation also
  rewrites the whole file and scans every entry — fine at hundreds of actions, not at 10^5.
- **B-MINOR-8 — lock acquisition has no fairness and no intrinsic ceiling.**
  `lock.go:48-87` polls `LOCK_NB` every 20 ms with no backoff; the **only** bound is the
  caller's context, and there is no queueing, so a contended caller can starve. Callers
  MUST pass a bounded context — a `context.Background()` caller spins indefinitely.
- **B-MINOR-9 — a cache dir on a no-op-flock filesystem is a total outage with a
  misdirecting error.** If `os.UserCacheDir()` lands on NFS/SMB (network home directories
  are common), the probe may fail and **every** budget operation returns "budget lock
  filesystem does not provide verified exclusion" (`lock.go:75`). Failing closed is right;
  the message does not name the resolved lock path, so an operator will inspect the ledger
  directory — which is not the problem. Include the path.
- **B-MINOR-10 — case folding is fail-safe on darwin and fail-open on Linux.**
  `lock.go:26-28` lowercases on windows/darwin only. On a case-**sensitive** APFS volume
  this over-locks (two distinct dirs share one lock: safe, less concurrency). On Linux with
  a case-**insensitive** mount (ciopfs, NTFS/exFAT) it under-locks: two spellings of one
  directory take two locks and exclude nothing. Same silent shape as B-MAJOR-1/2.

#### [NIT]

- `ledger.go:101` shadows the builtin `cap`.
- `update` returns the live `state` alongside errors (`:194`, `:198`, `:202`), so the
  returned `Snapshot.Entries` map is shared rather than defensively copied — unlike the
  `clone` discipline in `internal/telemetry`. Harmless today; nothing retains it.
- `release()` from `lock()` swallows the unlock error and is not idempotent. Only one call
  site (`ledger.go:181`, deferred), so fine — worth a comment.
- The probe costs an extra `open` + 2 `flock` syscalls on **every** ledger operation, not
  once per process. Cheap and worth keeping; just not free.
- `lock.go:13-16`'s "one host" should say one host **and one user cache** once B-MAJOR-1 is
  resolved either way.

### 3.4 Test-coverage observations (inspection of the committed suite only)

The suite is well-aimed — `TestPersistFailureNeverAuthorizesWork`'s replaced-then-errored
case and `TestKernelLockSurvivesProcessDeath`'s kill-the-owner case are both exactly the
right experiments. Gaps I can see by reading:

- No shared/no-op-flock regression (B-MAJOR-3) — the most important gap.
- `TestConcurrentReservationsCannotOverspend:45-69` is **in-process goroutines**, which
  exercises the cross-FD flock path but not cross-**process** counting. The observed bug was
  cross-process. `TestKernelLockSurvivesProcessDeath` proves cross-process *exclusion* but
  never checks a cap. Nothing asserts "N processes, cap M ⇒ exactly M admissions".
- That same test is bounded by the **cost** cap (5 × 3 = 15) before the action cap of 8, so
  the action cap under concurrency is untested.
- Unexercised `read` validations: symlinked ledger (the `IsRegular` guard), oversize,
  bad-hex entry id, `actual_micros` with `settled: false`, `reserved_at` before
  `started_at`, scope mismatch, backward clock.
- `TestBudgetLockChild:214` sleeps 30 s; if the parent's `Kill` path is ever skipped this
  is a 30 s tax. Cleanup at `:171-176` covers the normal case.

### 3.5 Does the stated boundary suffice?

**The boundary is the right call.** Refusing to claim distributed cross-host locking, and
adding a probe rather than trusting `flock`, is the correct response to the observed
evidence — I want to say that plainly, because it is the strongest decision in this slice.

**But as implemented the guarantee is narrower than the boundary states**, in the two ways
above (user/env scope, canonicalization order), and **both narrowings are silent**. For a
package whose entire purpose is "a repeated ID is never permission to repeat work", a
silent loss of exclusion is worse than a loud absence of it. So: the boundary suffices *as
a design*; the *documented* boundary does not yet match the *delivered* one. Close
B-MAJOR-1/2, or narrow the claim in `lock.go:13-16` and in caller docs to exactly what is
guaranteed. Either resolution is acceptable; leaving the gap between them is not.

### 3.6 Integration constraints for future callers

Stated as constraints, **not** as a claim that the store enforces production budgets — it
does not, and nothing calls it yet.

1. **Each run/scope needs its own `Dir`** — `Scope` is pinned per directory (`:193`), the
   wall-clock origin is the directory's creation (`:189`), and entries accumulate forever.
2. **Always pass a bounded context** to `Reserve`/`Settle` (B-MINOR-8).
3. **Generate a fresh ID per execution attempt.** The store deliberately cannot distinguish
   "crashed before doing the work" from "already did the work". All retry policy lives in
   the caller, and every attempt is charged.
4. **Always supply a conservative `ReserveMicros`, even when no cost limit is set**, or you
   permanently poison cost accounting for the scope (B-MAJOR-4).
5. **`Settle` every completed action**, with a known figure whenever one exists; a nil
   actual is a permanent unknown, and a wrong figure cannot be corrected (`:133-137`).
6. **Build crash reconciliation yourself, or add it here.** There is no
   reconcile/abandon entry point; §12.7's reconcile duty is unmet by this package.
7. **A cap of 0 is unlimited**; there is no way to deny a kind (B-MINOR-1).
8. **This is an accounting ledger with a lock, not a spend gate.** Nothing observes real
   cost; `ReserveMicros` is asserted and `Settle` is reported.
9. **Treat Windows as unvalidated.** `lock_windows.go` is compile-checked only
   (facilitator-reported cross-compile PASS, runtime untested), and the probe — which is
   load-bearing for correctness — depends on RECALL(b), cross-handle `LockFileEx` conflict
   within one process, which no one has run.
10. **Single-version deployments only**, until B-MINOR-5 is addressed.

---

## 4. Summary of my verdicts

- **Concur closed:** MAJOR-1, MAJOR-2, MAJOR-3 (as scoped), MINOR-1..MINOR-5, and all four
  prior NITs that were addressed. The corrections are real and, in two cases
  (MINOR-3's single shared predicate, MAJOR-2's harness interaction-count assertion),
  better than what I proposed.
- **Concur with retention:** the constant `/bin/sh` dependency; Windows-untested status.
- **Still open, not disputed by the owner:** MINOR-6 (disclosed ≠ fixed), open question 3.
- **Resolved:** open question 1, for the success path, on owner-PRIMARY + an independent
  nonzero-fd fixture I can verify by inspection.
- **New this pass (terminal):** NEW-1 [MINOR], NEW-2/NIT-A/NIT-B [NIT].
- **New this pass (budget):** B-MAJOR-1, B-MAJOR-2, B-MAJOR-3, B-MAJOR-4; B-MINOR-1..10;
  five NITs.
- **No CRITICAL** found in either slice.

## 5. Standing obligations this note does not touch

Unchanged and still open: all driver/manual/resume/BLOCK/launch budget integration; complete
launch coverage (AC-T1); the exact ratified packet experiment (AC-P2); the ≥20-attempt
reconciliation (AC-T2); the twelve-task pilot (AC-X1/X2); the Kimi typed-evidence
integration; readiness/liveness validation; any quorum amendment; and the full Phase-6/7/8
review and signoff cycle. **Nothing here waives any of them, and nothing here is a
signature.**
