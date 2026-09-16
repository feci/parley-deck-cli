---
agent: zcode-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-16
kind: follow-up source review — publication-race reachability + NIT dispositions (read-only; no tests executed in this launch)
amends: zcode-1-recovery-source-review-20260916.md (preserved unchanged; this note withdraws its [MINOR] via SELF-CORRECTION and dispositions its two [NIT]s)
reviewed-source: /private/var/folders/yt/p2sr23f12_qcfx_w2z5c1p4r0000gn/T/parley-recovery-publication-pmwo7ia9/source
---

## Scope

One question: is the claimed interleaving — "two concurrent applies both observe NotExist
(parent_recovery.go:612-615), both pass digest validation, both rename" — reachable under
supported cooperating calls once the previously unread withState layer is read? Plus
dispositions of the two NITs. No whole-candidate re-review; no known gate re-opened; no
acceptance or signoff.

## Decisive chain (all PRIMARY: located and quoted this session, read-only)

1. `parent_recovery.go:585` — `recoveredParent` passes the whole read/derive/publish
   callback to `withState`: `err = withState(ctx, root, idea, func(b budget.CycleBinding, _ budget.Snapshot, s State) error {`.
2. `state.go:466-467` — `withState` delegates: `return withStateResolutionCheck(ctx, root, idea, checkResolutions, fn)`;
   `state.go:475` — `return withValidatedState(ctx, root, idea, check, fn)`.
3. `state_validation.go:49` — the caller's `fn` runs only inside the second
   `withStateControl` callback: `return fn(b, ledger, s)`; the file header states it:
   "Only an unchanged, fully validated authority may reach the caller callback, which
   still holds the guard" (`state_validation.go:11-12`).
4. `state.go:492/496/521` — `withStateControl` acquires the guard before reading state and
   calls `fn` before releasing it: `release, err := budget.AcquireResourceGuard(wait,
   filepath.Dir(b.Store.Dir))` … `defer release()` … `return fn(*b, ledger, s)`; the wait
   budget is `context.WithTimeout(ctx, 30*time.Second)` (`state.go:490`).
5. `resource_guard.go:21-22,32` — `AcquireResourceGuard` → `lockWithReady(ctx,
   filepath.Join(dir, "resource"), tryLock, unlock, …)`.
6. `lock_unix.go:14` — `err := unix.Flock(int(f.Fd()), unix.LOCK_EX|unix.LOCK_NB)`;
   EWOULDBLOCK → `(false, nil)`; `lock.go:174-184` retry loop (`held, err := take(f)` at
   :179) waits until acquired or the 30 s context expires, then fails closed.

## SELF-CORRECTION — withdrawing the [MINOR] concurrent-first-apply race

Verdict (PRIMARY, by the chain above): the claimed interleaving is **not reachable**
through supported cooperating calls. Same root+idea resolves the same persisted binding
(`state.go:483/497` — `LoadCycleBinding(ctx, root, idea, budget.Fixup)`), hence the same
`b.Store.Dir`, hence the same `resource` guard file. Exclusive flock applies per open
file description — so it excludes a second process AND a second goroutine of the same
process. The first apply completes stage→write→fsync→rename→dirsync
(`parent_recovery.go:561-573`) before its `defer release()`; a contender waits (≤30 s)
or fails closed, and once inside, its `readRecoveredParent` (:592) finds the record,
takes the retained replay-only branch (:593-610, refusal `if apply && expected !=
retained.SHA256` at :601) and never reaches `publishRecoveredParent`. That helper has
exactly one call site — :623, on the NotExist branch (grep over the tree: definition
:551, call :623, no others). "Both observe NotExist" requires both callbacks inside the
guard simultaneously; LOCK_EX excludes exactly that.

Narrow facts that still hold: `dir.Rename(stage, …)` at :570 does atomically replace, and
the record embeds `time.Now().UTC()` (:623) — but under supported calls only one apply
per run reaches them, so "validated and replayed … never rewritten" (:643-646) is
accurate for the supported surface.

Withdrawal scope (delimiting, not re-claiming):
- Writers bypassing `recoveredParent`/the guard — direct FS mutation, deleting the lock
  inode/origin (`lock.go:26` "Never delete a live lock inode or its origin file";
  `resource_guard.go:19-20`) — are outside supported cooperating calls; not re-framed
  as a finding (this is not the same-UID-hostile-mutation substitution).
- Proven exclusion scope is same-host. The lock's own contract says cross-host is
  refused, not raced (`lock.go:23-24`: "another cache environment or hostname refuses …
  This is same-origin coordination, not distributed cross-host locking") — quoted
  design intent, not my execution evidence.
- Whether two DISTINCT idea bindings could validly target one runID directory is
  unestablished by anything I read; no claim either way.

The old note stands otherwise; only this MINOR is withdrawn. No executable fixture is
supplied because no reachable interleaving exists; manufacturing one would be an
execution claim I cannot make (no tests run in this launch).

## NIT dispositions

Single-slot Observe closure (old note, `trajectory_verifier_recovery.go:262-264`):
`Observe` is invoked only via `(*launchEvidence).notify` (`telemetry.go:160-164`), called
at `telemetry.go:149` (constructor), `:170` (`started`), and `:237` (inside `finish`,
itself `once.Do`-guarded at `:175`); the relaunch reads `observed` at `:282-288` after
the synchronous `runner.RunConsult` return (`:277`). No concrete defect established by
this chain; a full absence-of-race proof would require the runner's whole goroutine
model, outside this follow-up. Per instruction and §15.1 (weakening my own item):
reclassified **unknown, not a finding**.

`--dir` not bound to a ticket root (old note): classification clarified as a **design
note only, not a defect** — a wrong dir/idea fails closed with no write
(`canonicalRoot` + `runtimeID` gate at `parent_recovery.go:578-583`; absent live state,
`withStateControl` errors or `:630-631` returns "trajectory is not active"), so no
incorrect behavior is reachable. The item compares relaunch's `Ticket.Root` refusal with
the publication route's state-rooted validation: hardening commentary, not a fault.

## Provenance & integrity

All locators above were located and quoted in this session from the fixture tree via
read-only content search; no execution, no network, no MCP, no git, no configuration
changes. Tags: chain = PRIMARY; Observe description = PRIMARY as to located call sites,
UNKNOWN as to full thread-safety; cross-host stance = PRIMARY-located quoted comment.
This is original audit supporting review — not an amendment round artifact, not an
acceptance, and not a signoff. No prior testimony was upgraded; no self-owned claim was
strengthened (both changes are weakenings/clarifications of my own items).
