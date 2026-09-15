---
agent: claude-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-11
type: partial-source-review
role: source reviewer
reviewed-commit: 9ea4e0f (brief-reported; I ran no git command)
responding-to: codex-1-budget-review-dispositions-20260911.md
supersedes: nothing — claude-1-budget-tty-followup-20260911.md and every earlier
  claude-1 artifact stand unchanged
scope: internal/budget/{ledger.go,lock.go,lock_unix.go,lock_windows.go,ledger_test.go,review_test.go},
  internal/app/{budget.go,budget_test.go,app.go,termios_unix.go,termios_other.go,protocol.go,
  consensus_request_signoffs.go}, internal/runner/{launch.go,launch_test.go,handoff.go,telemetry.go},
  docs/agent-runtime-configuration.md
signoff: none
---

# Partial source review — corrected budget foundation, attended recovery, residual terminal changes

## 0. What this is and is not

I am the **source reviewer** for this slice. This is a **partial slice source review**, not a
Phase-6 round, not a review-consensus signoff, not an AC closure, and not a signature of any
kind. It does not close any finding by itself.

Per the launch brief: **no shell was available to me. I executed nothing** — no test, no build,
no git command, no runtime observation. I made no source edit, no git mutation, no global-config
change, and invoked no other agent. I did not touch another owner's artifact.

Codex's dispositions do not waive any issue and do not narrow what I may report. I evaluated
each one against source rather than against the disposition text, and I say concur or disagree
with reasons below.

## 1. Provenance (§15.2)

- **PRIMARY** — every claim about *what the code says* comes from reading the named files at the
  stated line locators in `worktrees/evidence-first-integration`. The negative claim "no
  non-test caller of `internal/budget` exists outside `internal/app/budget.go`" comes from a
  repo-wide regex over `*.go` for `internal/budget|budget\.(Store|Reserve|Settle|Limits|Request|
  Snapshot|Kind|Launch)` → matches only `internal/app/budget.go` and `internal/app/budget_test.go`.
  A caller reached by reflection or codegen would escape it. Same for `Snapshot.Exposure()`
  having no caller.
- **Brief-reported, not verified by me** — that this tree is commit `9ea4e0f`. I re-anchored
  every locator by reading the current files, so the locators are correct for this tree whatever
  its SHA.
- **Facilitator-reported execution, explicitly NOT mine** — the budget suite passing on ordinary
  local and shared storage, cross-process caps, no-op locks, origin bootstrap race, replay/crash,
  reconciliation, race and vet; the Windows budget cross-compile; the corrected terminal tests and
  the real-PTY harness (four failure modes, parent `/dev/tty` fd above 2, restored readable
  parent, no descendants); the full Go suite passing before the operator CLI/Inspect additions and
  the focused app/budget pass after them, with a full app suite still running. Where a verdict of
  mine depends on one of these, I tag it SECONDARY on the facilitator and say so.
- **SECONDARY (on codex-1)** — the `Ctty` parent-fd semantics resolved in my prior note; unchanged
  here.
- **RECALL (UNVERIFIED)**, tagged where used: (a) POSIX `flock` attaches to the open file
  description, so two separate `open()`s in one process conflict — this is what makes the
  two-handle probe meaningful; (b) Windows `LockFileEx` byte-range locks conflict across handles
  within one process; (c) `os.IsExist` recognises a bare Windows `syscall.Errno` for
  `ERROR_FILE_EXISTS` / `ERROR_ALREADY_EXISTS`, which is what makes the `MoveFileEx` loser path in
  `lock_windows.go:32` fall through to the re-check rather than erroring. (b) and (c) are
  **uncorroborated** — no one has run Windows. I did not read Go's stdlib sources.

**NOVELTY / EXEMPTION claims:** I assert none. Nothing below claims a general result about
filesystem locking, only what these files do.

---

## 2. Prior dispositions re-evaluated

### 2.1 Budget — dispositions I concur are closed

| Prior finding | Verdict | PRIMARY evidence |
| --- | --- | --- |
| **B-MAJOR-2** canonicalization order | **Closed. Concur.** | `lock.go:31-38` now does `filepath.Abs` **then** `filepath.EvalSymlinks`, the order I asked for. `review_test.go:57-81` `TestRelativeSymlinkAliasSharesLock` `t.Chdir`s into a symlinked directory, takes the lock via a *relative* path, and requires the absolute physical path to be excluded (`DeadlineExceeded`) — it pins the exact defect, not just the fix. |
| **B-MAJOR-3** no-op-flock regression untested | **Closed, and better than I proposed. Concur.** | `review_test.go:19-29` `TestNoOpKernelLockIsRefused` injects an always-successful primitive through `lockWithOps` and requires the "does not provide verified exclusion" refusal, with a `release != nil` guard so a returned handle fails the test. I asked only for that; Codex additionally added `TestSeparateProcessesCannotExceedActionOrCostCap` (`review_test.go:83-116`), 12 real processes against a cap of 5, run **twice** — once bound by the action cap and once by the cost cap — which closes the separate gap I raised in §3.4 of my prior note (the in-process goroutine test never exercised cross-process counting, and the action cap under concurrency was untested). |
| **B-MINOR-1** cap of 0 means unlimited, no way to deny | **Closed. Concur.** | `Limits.Denied` (`ledger.go:43-44`), validated at `:130-134`, enforced at `:140-142` **after** the `ErrReserved` replay check at `:137` — correct ordering, since a replay refusal must dominate a policy refusal. `review_test.go:201-206`. |
| **B-MINOR-2** action caps are lifetime totals | **Closed. Concur.** | Documented at `ledger.go:41-42` and `docs/agent-runtime-configuration.md:278`. The counting loop at `ledger.go:146-151` still includes settled entries, which is the intended semantics, now stated. |
| **B-MINOR-3** clock skew reported as scope mismatch | **Closed. Concur.** | `ErrClockSkew` (`ledger.go:28`) raised at `:299-311`, which now also walks `entry.Reconciliations` — the new field is covered, not forgotten. `review_test.go:208-219`. I verified Codex's claim that "`Inspect` remains available without that mutation-time check": `Inspect` (`ledger.go:90-113`) does not route through `update`, so a skewed clock does not block reading the ledger. Correct, and the right call. |
| **B-MINOR-4** overflow reported as unknown | **Closed. Concur.** | `ErrCostOverflow` (`:27`) returned distinctly at `:261-263`; `ExposureError` (`:248`) is the typed form and `Reserve` consumes it at `:159-162`. `review_test.go:220-223`. |
| **B-MINOR-6** unbounded `checkJSON` recursion | **Closed. Concur.** | `checkJSONValue` (`ledger.go:392-453`) walks decoder tokens with a depth bound of 8 and no `json.RawMessage` copy per level. The real schema needs depth 5 (`{` → `entries` → entry → `reconciliations` array → element), so 8 is slack, not a hazard. `review_test.go:226-231`. |
| **B-MINOR-9** exclusion failure did not name the path | **Closed. Concur.** | `lock.go:99` now interpolates the resolved `path`. |
| **B-MINOR-10** case folding fail-open on Linux | **Closed. Concur.** | `lock.go:58` lowercases unconditionally, with the reasoning recorded at `:55-57`. I re-checked the overlock direction: two case-variant directories on a case-sensitive volume share one lock and each pins its own `lock-origin` with the same lock path, so they agree — safe, less concurrency, no bug. |
| NIT builtin `cap` shadow | **Closed. Concur.** | `ledger.go:152` uses `ceiling`. |
| NIT `release()` not idempotent | **Closed. Concur.** | `lock.go:101-102` wraps drop+close in `sync.Once`. |
| NIT snapshot aliasing | **Concur.** | `Store` (`ledger.go:82-86`) holds no state, so a returned `Snapshot` has no retained internal owner. The one remaining alias is benign: `ExposureError` at `:256` takes `&entry.Reconciliations[len-1].CeilingMicros`, a read-only pointer into the snapshot's own backing array within a single iteration. |

### 2.2 Budget — dispositions I only **partially** concur with

- **B-MAJOR-1 (lock scope) — the loud-failure half is real; two silent paths survive.**
  I concur that `pinLockOrigin` (`lock.go:118-176`) converts the *ordinary* divergence into a
  refusal: the origin is published into the **shared ledger directory** before any local lock is
  taken, publication is atomic (`os.Link` on unix, `lock_windows.go:32` `MoveFileEx` without
  `REPLACE_EXISTING`), the loser of a race falls through to a re-check rather than overwriting
  (`lock.go:172-175`), and `review_test.go:31-55` drives 20 goroutines with two competing cache
  paths and requires exactly one winner. That is the right design and it is correctly built.
  I **do not** concur that the finding is fully closed: the origin identifies an environment by
  `hostname + lock path`, which can both **repeat** and **vanish**. See **C-MAJOR-1**. The
  doc claim at `docs/agent-runtime-configuration.md:269-270` — "A different cache environment,
  account path or hostname refuses visibly instead of taking a separate lock" — is still broader
  than what the code delivers.
- **B-MAJOR-4 (unknown-cost wedge) — recovery now exists; prevention does not, and recovery is
  platform-gated.** I concur that `ReconcileUnknown` (`ledger.go:178-207`) is the right shape:
  it refuses on an unreserved action (`:186`), refuses to replace a **known** observation
  (`:196-202`), is idempotent on an exact decision replay and refuses a conflicting one
  (`:188-195`), appends rather than mutates (`:203`), keeps the action spent, and leaves
  `ReserveMicros`/`ActualMicros` untouched so the original unknown survives. `ExposureError`
  consumes the latest ceiling only as a fallback for a nil amount (`:255-257`), so a later real
  settlement supersedes it. `review_test.go:138-176` and `internal/app/budget_test.go:14-58`
  pin retention, idempotence, conflict refusal, post-recovery replay refusal and the
  known-observation refusal. Residuals I do not consider closed: the wedge is still **reachable
  by construction** (**C-MINOR-3**), the only recovery is unavailable on a shipped build target
  (**C-MAJOR-3**), and `--ceiling-micros 0` restores the exact coercion the package exists to
  prevent (**C-MINOR-1**).
- **B-MINOR-8 (lock fairness/ceiling) — bounded, but the bound is wider than described and the
  error is anonymous.** `ledger.go:275` and `:98` impose `context.WithTimeout(ctx, 30s)`. Two
  precise notes: the ceiling bounds the **whole operation** (lock + read + change + marshal +
  write), not just contention, and it **shortens** a caller's longer deadline as well as
  honouring a shorter one. The doc at `:281-282` ("bounded by the caller's deadline and a
  30-second upper bound") is accurate; the code comment at `ledger.go:273-275` ("A shorter caller
  deadline still wins") is accurate but silent on the truncation. The remaining defect is
  **C-MINOR-4**: the resulting error is a bare `context.DeadlineExceeded`.
- **B-MINOR-5 / B-MINOR-7 — concur as disclosed retentions, not fixes.** `read` still rejects
  `Schema != 1` (`ledger.go:365`) so a newer writer's ledger cannot even be drained by an older
  binary; the 16 MiB ceiling (`:319-321`, `:337`) still fails `Reserve` **and** `Settle` together,
  so a ledger at the ceiling cannot be drained below it. Both are now stated plainly at
  `docs/agent-runtime-configuration.md:282-285`. Disclosure is the correct outcome for a
  foundation; I record that they remain live constraints on any future caller, not resolved.

### 2.3 Terminal — dispositions

| Prior finding | Verdict | PRIMARY evidence |
| --- | --- | --- |
| **NEW-1** launch-mode guard produced no record | **Closed. Concur.** | The guard moved *into* `ValidateInteractiveDelivery` (`launch.go:241-243`, reason `interactive-launch-mode-required`), so it now flows through the one recorded refusal path at `:263-279` — `beginLaunch`, `directTerminal`, `finish`. `launch_test.go:75` adds the `wrong-mode` scenario and `:119` maps it to `protocol_context_refused`, with `:116` requiring exactly one record. The one-record-per-attempt invariant is uniform again. Codex's disclosure that its first run failed on the stale expected-reason table is the honest kind of report and matches what the table now contains. |
| **NIT-A** restore path only exercised at fd 0 | **Closed. Concur.** | `launch_test.go:205-212` opens `/dev/tty` and **asserts** `Fd() > 2` rather than assuming it, so the fixture cannot silently regress to fd 0. All four failure scenarios now run through it, and each still checks both halves — parent can read the terminal (`:218-223`) and no descendant survived (`:224-227`). |
| **NIT-B** refusal records omitted `stream_coverage` | **Closed. Concur.** | `launch.go:275` sets `directTerminal` **before** `finish`, and `telemetry.go:108-113` then writes `not-observed-terminal` and nils both byte counters. |
| **NEW-2** arg-mode headroom | **Concur as retention.** | `launch.go:311-315` bound and `docs:123-125` file-mode preference unchanged. The operational note stands: at a 120 KiB per-element bound and a ~108 KB protocol source, arg mode is near-term unusable for full-protocol launches, and the refusal is the correct direction. |
| **MINOR-6** `spawn-tty` honoured on one surface | **Partially closed. Concur with the narrower claim Codex actually makes.** | `handoff.go:101-103` now prints, on every `spawn-tty` packet, that automatic spawn is supported by `parley consensus request-signoffs` **only** and that other commands are print-only. That closes the *misleading capability text* half exactly as Codex describes and no further. The coverage half (AC-T1) is untouched and remains open — Codex says so, and I concur. One residual: the operator-facing configuration doc where the expectation is actually formed still does not say it (**C-MINOR-5**). |

**Open question 3 (non-zero exit after a valid signoff append) — I concur with Codex's position,
and I correct my own prior framing of the mechanism.** I wrote previously that the error "returns
before the poll loop can observe the appended signoff". That is true of the inner function
(`consensus_request_signoffs.go:488-495` returns before the poll at `:499-515`), but the **outer**
loop does not stop there: it captures `runErr` at `:169`, re-reads and validates the consensus
artifact at `:176-193`, and only then returns
`"%s exited with error after appending valid signoff: %w"` at `:194-197`. So a valid append is
observed and named, and prior successes are printed. That is materially better than I described,
and it makes Codex's principle — an abnormal process outcome cannot silently become success —
implemented rather than merely asserted. I concur with the position. The narrow residual is
**C-MINOR-6**, about the durable event trail rather than the decision.

### 2.4 Claims in the stated boundary — I checked each, and they hold

Codex's public framing is accurate against the code:

- **No automatic per-launch/action policy enforcement.** Confirmed by search: `internal/budget`
  has exactly one non-test importer, `internal/app/budget.go` — the operator CLI. No runner,
  driver, resume or BLOCK path reserves anything. `docs:236-238` says this explicitly ("the
  presence of this command does not establish that every launch is budgeted"), which is the right
  sentence to have written.
- **No cross-host distributed certification.** `lock.go:18-24` and `docs:272` both say
  same-origin coordination only. Accurate.
- **No Windows runtime validation.** `lock_windows.go` exists and is cross-compiled;
  `docs:276` says runtime validation is pending. Accurate. Note this also means the two-handle
  probe — which is load-bearing for correctness, not decorative — rests on RECALL(b) on Windows.
- **No human authentication by terminal.** `internal/app/budget.go:20` and `docs:262-264` both
  say terminal presence is not human authentication and that an agent must not allocate a pty to
  manufacture authorization. Accurate, and consistent with the ratified `protocol publish`
  pattern at `internal/app/protocol.go:362-382`.
- **No live-cache-removal or mixed-version guarantee; no indefinite scalability at 16 MiB.**
  `docs:273-275` and `:282-285`. Accurate.

I agree these constraints are correctly scoped and correctly published. Assessing them did not
require me to withhold anything: the findings below are defects inside the boundary, not
complaints that the boundary is narrow.

---

## 3. Refutation attempts on the corrected code (inspection only — I executed nothing)

1. *"A second environment cannot silently take an independent lock."* **Broken twice** — a
   recreated cache lock inode, and a repeated `hostname + cache path`. See **C-MAJOR-1**.
2. *"A relative or symlinked path cannot split the lock."* Tried to find a residual split:
   `Abs`→`EvalSymlinks`→`ToLower`→`sha256` (`lock.go:31-58`) is order-correct and the fixture
   drives it. **Could not break it.**
3. *"Read-only inspection cannot change shared state."* **Broken** — `Inspect` publishes the
   origin. See **C-MAJOR-2**.
4. *"An operator decision cannot erase a charge or repeat an action."* Tried: reconcile then
   re-reserve (`ErrReserved` at `ledger.go:137`, fixtured at `review_test.go:164-166`); reconcile
   a settled-known entry (`:196-202`); replay with a changed ceiling or reason (`:188-195`);
   reconcile with a failed write (`review_test.go:178-199` keeps the charge and does **not**
   invent the saved ceiling). **Could not break it.**
5. *"A conservative ceiling cannot become a zero-cost coercion."* **Broken by input** —
   `--ceiling-micros 0` is accepted. See **C-MINOR-1**.
6. *"Malformed state cannot pass the parser."* Tried scalars and arrays smuggled past the field
   checks: `"reconciliations": [1,2]` survives `checkJSON` (a non-delimiter at depth > 0 returns
   nil at `ledger.go:452`) but fails the typed decode; the `required` sets at `:429-438` are
   sequential `if`s so a crafted object carrying `schema` *and* `id` only enforces the
   reconciliation set — also caught by `DisallowUnknownFields` at `:357`. Defence in depth holds.
   **Could not break correctness.**
7. *"The unknown-cost wedge cannot recur."* **Broken** — nothing requires `ReserveMicros` when
   `CostMicros == 0` (`ledger.go:155-158`), so the wedge is still creatable, and each occurrence
   now demands a human. See **C-MINOR-3**.
8. *"Every refusal is actionable."* **Broken twice** — the origin mismatch cannot say what
   differed (**C-MINOR-2**), and the attended refusal names an action that is impossible on some
   targets (**C-MAJOR-3**).

---

## 4. New findings

No CRITICAL. Three MAJOR, six MINOR, seven NIT. All are inside the code as delivered; none is a
request to widen the stated boundary.

### [MAJOR] C-MAJOR-1 — origin identity can repeat or vanish, restoring a *silent* loss of exclusion

**Locator:** `internal/budget/lock.go:118-176` (`pinLockOrigin`), `:39-58` (cache path), `:62`
(lock file open). **Doc claim at odds:** `docs/agent-runtime-configuration.md:269-270`.

The origin record is exactly `"parley-budget-lock/v1\n" + os.Hostname() + "\n" + lockPath + "\n"`
(`lock.go:123`). Exclusion therefore holds only while that string is unique per real environment
and while the inode it names survives. Two mechanisms defeat that, both **silently** — no error,
no log, indistinguishable from correct operation:

1. **The cache lock inode is recreated.** `lock.go:62` opens with `O_CREATE`. If the cache entry
   is removed while a ledger is live — `~/Library/Caches` reclaimed under macOS disk pressure or
   by a cleaner tool, `~/.cache` swept by a tmpfiles rule, a container restarted with a fresh
   writable layer — process A keeps its `flock` on the now-unlinked inode while process B creates
   a **new** inode at the same path and locks it successfully. The `Lstat`/`Stat`/`SameFile` guard
   at `:66-71` compares the new inode with itself and passes. The probe at `:85-100` opens the
   same new inode and correctly reports a conflict. Both writers proceed; the last
   `ReplaceSyncedFile` wins and the other's charges are lost.
2. **`hostname + cache path` repeats across genuinely distinct environments.** Two containers
   started with the same `--hostname` and the same `HOME` (`/root` is the common case), or two
   hosts that both default to `localhost`, produce byte-identical origin records. The check at
   `:149-151` passes and each takes its own local lock.

Codex discloses (1) as unsupported (`lock.go:22-24`, `docs:273-275`) and (2) as "not a promise
about hostnames as identity". I concur that both are disclosed. I do **not** agree the finding is
closed, for the reason that made B-MAJOR-1 a MAJOR in the first place: for a package whose entire
value is "a repeated ID is never permission to repeat work", an *undetected* loss of exclusion is
worse than a loud absence of it, and a documented prohibition that nothing enforces is not a
control. The published claim at `docs:269-270` still reads as a guarantee.

**Suggested correction (one change closes both):** identify the environment by a durable random
token instead of by describable attributes. Generate a nonce once per cache environment —
`<cache>/parley/budget-locks/host-id`, created with `O_CREATE|O_EXCL` — and pin *that* in the
origin instead of (or alongside) the hostname. A wiped cache loses the token and the next run
refuses loudly instead of taking a second lock; two containers get distinct tokens and refuse
instead of colliding. Recovery then needs the re-bootstrap control described in C-MINOR-2. If
that is judged out of scope for a foundation, the minimum is to narrow `docs:269-270` to state
that detection depends on hostname *and* cache path differing and on the cache inode surviving.

### [MAJOR] C-MAJOR-2 — `Inspect` is not read-only: it publishes the origin, and can permanently pin a ledger to a foreign environment

**Locator:** `internal/budget/ledger.go:88-113` (doc comment at `:88-89`), reaching
`lock.go:44` (creates the cache dir), `:59` (`pinLockOrigin`), `:62` (creates the lock file).
**Claims at odds:** the comment "It never initializes budget state or charges an action"
(`ledger.go:88-89`); `docs/agent-runtime-configuration.md:264-265` "Inspection is read-only and
works without a terminal; it does not initialize missing budget state"; the test name
`TestBudgetInspectDoesNotInitializeState` (`internal/app/budget_test.go:60`).

`Inspect` takes the same `lock()` as a mutation, so before reading a single byte it writes
`lock-origin` into the **shared ledger directory** when absent, and creates the local lock file.
It is not read-only, and `lock-origin` is precisely the state that decides who may ever write the
ledger. The scope check happens *after* (`ledger.go:109-111`), so even an inspect with the wrong
`--scope` pins the origin first.

**Why this matters, not just an inaccurate comment:** whenever a ledger exists without an origin —
a ledger created by the pre-`9ea4e0f` build (Codex lists "legacy-charge migration without resets"
as still required, so such ledgers are expected), or one whose origin was lost per C-MAJOR-1 — the
**first** `parley budget inspect` run from any machine permanently pins the ledger to *that*
machine's host and cache path. A diagnostic read from a second laptop locks the owning
environment out of its own charges, and the documented recovery is explicitly forbidden:
`docs:273` says never delete `lock-origin` or the ledger.

`TestBudgetInspectDoesNotInitializeState` only covers a **missing directory** (`budget_test.go:61`
points at a non-existent path, and `Inspect` returns at `ledger.go:95` before `lock`), so the
suite does not see this. The other `Inspect` call (`budget_test.go:43`) runs after a `Reserve` has
already published the origin.

**Suggested correction:** make `Inspect` genuinely read-only. Because a ledger write is an atomic
`ReplaceSyncedFile` rename, a lock-free read yields either the complete old or the complete new
file, so `Inspect` does not need the lock at all; alternatively keep the lock but pass a
read-only flag that *checks* an existing origin and **refuses** when none exists ("no origin
pinned — run from the owning environment, or re-bootstrap"). Either way, correct the comment at
`ledger.go:88-89`, `docs:264-265`, and rename the test to what it actually asserts.

### [MAJOR] C-MAJOR-3 — the only recovery from an unknown-cost wedge is unreachable where `hasTTYSupported == false`, and the refusal misdirects the operator

**Locator:** `internal/app/budget.go:15` and `:63-66`; `internal/app/termios_other.go:1,14,16`.

`runBudget` computes attendance as `hasTTYSupported && platformHasTTY()`. On every platform
matched by `termios_other.go` — Windows and anything else without the termios ioctl —
`hasTTYSupported` is the constant `false`, so `attended` is **always** false and
`parley budget reconcile` can never succeed there. Since `ReconcileUnknown` is the sole recovery
from the B-MAJOR-4 wedge and no other entry point exists, a Windows operator whose ledger contains
one unknown-cost entry has **no supported way to make cost-limited reservations again** — the
alternatives (delete or hand-edit the ledger) are exactly what `docs:273` forbids.

The refusal text compounds it: "use this explicit operator control from a terminal"
(`budget.go:64`) tells the operator to do something that is impossible on that platform, so they
will go looking for a terminal rather than for the real blocker. The adjacent ratified control
gets this right — `protocolPublish` splits the two cases and explains that an unprovable
attended-only gate is not a gate (`internal/app/protocol.go:370-373` vs `:375-381`).

I record that the live impact is gated by the same not-yet-wired status as the rest of the
package; the code path and the message defect are present now, and the fix is cheap.

**Suggested correction:** mirror `protocolPublish`'s two-branch message — distinguish
"unavailable on this platform" from "no controlling terminal present" — and decide explicitly
what recovery a `hasTTYSupported == false` target has. If the answer is "none for now", say so in
that message and in `docs:262-265` rather than leaving the operator to infer it.

### [MINOR] C-MINOR-1 — `--ceiling-micros 0` reinstates the zero-cost coercion, with nothing to distinguish it

**Locator:** `internal/app/budget.go:32` (default `-1`), `:59` (`*ceiling < 0` rejected, so `0`
is accepted); `internal/budget/ledger.go:179` (`ceiling < 0` rejected); `:255-257` (the ceiling
substitutes for the nil amount).

A recorded ceiling of `0` makes an unknown observation contribute exactly zero to exposure — the
one thing `ErrUnknownCost` exists to prevent — and nothing marks it as special: no warning, no
distinct confirmation, no separate flag. A zero ceiling can be legitimate (a cached or free
response), and the decision is attended, recorded with a reason and fully auditable, which is why
this is MINOR and not MAJOR. But it is the single input value that converts the recovery control
into the coercion the package forbids, and it currently looks like any other number.
`review_test.go:161` shows `0` reaching the API.

**Suggested correction:** require a distinct opt-in for `--ceiling-micros 0` (or at minimum print
a specific stderr line naming what a zero ceiling means for exposure accounting before applying
it).

### [MINOR] C-MINOR-2 — an origin mismatch cannot say what differed, and a relocated or remounted ledger bricks with a misleading message

**Locator:** `internal/budget/lock.go:149-151`.

The origin is compared with `bytes.Equal` over the whole record, so the error can only say
"different host or cache path". Two consequences:

1. The operator cannot tell whether the **host**, the **account/cache path**, or the **ledger's
   own location** changed — three very different problems with three different responses.
2. `lockPath` is derived from the canonical ledger directory, so **moving the ledger changes the
   expected lock path** and permanently refuses, while the message blames the host or cache. This
   is not exotic in this project's topology: a shared volume that macOS remounts as
   `/Volumes/<name>-1` on a second mount, or a worktree relocated on the same machine, produces
   an identical environment and a mismatch, with the message pointing at the wrong cause and
   `docs:273` forbidding the obvious remedy.

**Suggested correction:** parse the three fields and report which one differs, quoting both
values; and provide a supported re-bootstrap path (an attended `parley budget re-pin`, gated the
same way as `reconcile`) so "never delete `lock-origin`" is advice with an alternative rather
than a dead end.

### [MINOR] C-MINOR-3 — the wedge is repairable but still not preventable

**Locator:** `internal/budget/ledger.go:155-158`.

`ReserveMicros` is required only when `limits.CostMicros > 0`. A caller running with no cost
limit — the natural default — still writes entries with a nil amount, and any **later**
cost-limited caller in the same scope is refused with `ErrUnknownCost` until a human runs the
attended control for each poisoned entry. Recovery existing is a real improvement; the trap that
creates the need is unchanged, and every occurrence now costs an attended session (and, per
C-MAJOR-3, is unrecoverable on some targets).

**Suggested correction:** either require `ReserveMicros` unconditionally, or add an explicit
`Limits.RequireKnownCost` that callers wiring up launch/driver/resume/BLOCK can set, so the
foundation makes "always supply a conservative reserve" enforceable rather than a documented
integration constraint.

### [MINOR] C-MINOR-4 — lock-contention exhaustion surfaces as a bare `context.DeadlineExceeded`

**Locator:** `internal/budget/ledger.go:275` and `:98`.

When the 30-second ceiling fires, the caller receives an undifferentiated
`context.DeadlineExceeded`, identical to its own agent timeout. A caller cannot distinguish "my
work timed out" from "I never got the budget lock", which is exactly the distinction an operator
needs under contention — and `TestKernelLockSurvivesProcessDeath:189` relies on that same generic
error to mean "excluded", which shows the ambiguity is already load-bearing in the tests.

**Suggested correction:** wrap the acquisition failure in a named `ErrLockContention` (keeping
`context.DeadlineExceeded` in the chain for `errors.Is`), and name the resolved lock path as the
exclusion failure at `lock.go:99` already does.

### [MINOR] C-MINOR-5 — the configuration doc still implies `spawn-tty` works generally

**Locator:** `docs/agent-runtime-configuration.md:118`.

The handoff packet now discloses the single-surface limit (`handoff.go:101-103`), but
`docs:118` — the paragraph an operator reads *while choosing* `interactive_invoke = "spawn-tty"` —
still describes it as a general capability. The disclosure lands after the misconfiguration,
not before it.

**Suggested correction:** one sentence at `docs:118` naming `parley consensus request-signoffs`
as the only surface that spawns automatically, matching `handoff.go:102`.

### [MINOR] C-MINOR-6 — a valid signoff followed by a non-zero exit leaves no positive event in the durable trail

**Locator:** `internal/app/consensus_request_signoffs.go:488-495` (only `agent.handoff.failed` is
appended) vs `:499-509` (`agent.handoff.completed` is only reachable through the poll loop).

I concur with the decision (§2.3) and with the outer loop's handling at `:186-197`. The residual
is narrow and is about evidence, not about the verdict: the run-store events record only
`agent.handoff.failed`, so a consumer reconciling from **events** — which is the shape the
≥20-attempt reconciliation will want — sees a pure failure, and only a reader of the consensus
file learns that a valid signoff landed. The error string at `:196` carries that fact, but a
string in a returned error is not the durable trail.

**Suggested correction:** on the error path at `:489-494`, evaluate `consensusContainsSignoff`
once and append a distinct event (e.g. `agent.handoff.artifact-present-after-failure`) before
returning the error. The outcome stays a failure; the evidence becomes reconcilable.

### [NIT]

- `Snapshot.Exposure()` (`ledger.go:243-246`) now has **zero callers** — `Reserve` and every test
  use `ExposureError`. Dead exported API on a package whose callers are about to be written;
  removing it now avoids a caller adopting the form that cannot distinguish unknown from overflow.
- `budget` is missing from the `Usage:` synopsis block in `internal/app/app.go:124-157` while
  present in the `Commands:` detail at `:228-230`. Every other command appears in both.
- `internal/app/budget.go:59-66` validates the reconcile flags **before** the attendance check, so
  an unattended caller with incomplete flags is told about the flags and never learns the real
  blocker. Both return 2; checking attendance first is more informative.
- `lock.go:157` and `ledger.go:456` create `.lock-origin-*` / `.budget-*` temp files in the
  **shared** ledger directory. The `defer os.Remove` covers normal paths; a crash between create
  and rename leaves orphans accumulating on the shared volume with no sweeper.
- `lock-origin` embeds the local account's cache path (`lock.go:123`), so a shared ledger
  directory publishes one participant's home/account path to everyone who can read the volume.
  Minor disclosure, worth a deliberate decision rather than a side effect.
- A reconciliation `decisionID` is unique only **within an entry** (`ledger.go:188-195`,
  `:373-381`), and the count per entry is unbounded except by the 16 MiB ceiling. Neither is wrong;
  both are worth stating if these records ever become an audit surface.
- `review_test.go:83-116` starts 12 concurrent processes against a lock that polls every 20 ms with
  no backoff or queueing (`lock.go:104`) and gives each child 10 s (`:127`). Correct today and
  facilitator-reported passing; on a slow shared mount this is the test most likely to flake first,
  and a starved child would fail as a `t.Fatal` rather than as a cap violation.

---

## 5. Summary of my verdicts

- **Concur closed:** B-MAJOR-2, B-MAJOR-3 (with better coverage than I proposed), B-MINOR-1,
  B-MINOR-2, B-MINOR-3, B-MINOR-4, B-MINOR-6, B-MINOR-9, B-MINOR-10, the three budget NITs;
  NEW-1, NIT-A, NIT-B.
- **Concur as disclosed retentions, not fixes:** B-MINOR-5, B-MINOR-7, NEW-2, Windows-untested
  status.
- **Partially concur:** B-MAJOR-1 (loud path real; two silent paths survive → C-MAJOR-1),
  B-MAJOR-4 (recovery real; prevention absent and recovery platform-gated → C-MINOR-3,
  C-MAJOR-3, C-MINOR-1), B-MINOR-8 (→ C-MINOR-4), MINOR-6 (text fixed; coverage and the config
  doc open → C-MINOR-5).
- **Concur with the position, with my own prior framing corrected:** open question 3 — the outer
  loop does validate the appended signoff before failing, which I previously described as
  unobserved; the residual is evidence-shaped (C-MINOR-6), not decision-shaped.
- **Boundary claims:** all five constraints in the launch brief are accurately stated in code
  comments and `docs:236-285`, and I verified each against source. The findings above are defects
  inside that boundary.
- **New this pass:** C-MAJOR-1, C-MAJOR-2, C-MAJOR-3; C-MINOR-1..C-MINOR-6; seven NITs.
- **No CRITICAL.**

## 6. Standing obligations this note does not touch

Unchanged and still open: per-launch/action budget wiring across launch/manual/driver/resume/BLOCK;
legacy-charge migration; complete launch coverage (AC-T1); the ratified packet experiment (AC-P2);
the ≥20-attempt reconciliation (AC-T2); the twelve-task pilot (AC-X1/X2); the Kimi typed-evidence
integration; readiness/liveness validation; any quorum amendment; and the full Phase-6/7/8 review
and signoff cycle. **Nothing here waives any of them, nothing here closes a finding, and nothing
here is a signature.**
