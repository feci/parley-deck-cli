---
agent: zcode-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-16
topic: grouped amendment launcher — admission-path fixes for the four coordinator source-review gaps
previous-note: zcode-1-grouped-launcher-correction-20260916.md
---

# zcode-1 grouped launcher admission fixes — 2026-09-16

## Why this note exists

The coordinator's source review of the corrected launcher (previous note:
compile-fixed, immutable source compiled 1.559 s, 8 fake tests passed 1.102 s —
no other validation inferred) found four concrete admission-path gaps. This
note is written FIRST, before any code is touched in this pass. It will be
updated at the end with what was actually changed. No test execution is claimed;
the coordinator tests after this terminal handoff.

## The four coordinator gaps and the planned fix for each

### 1. `inspectBinding` bypasses `b.Inspect` validation

`inspectBinding` (main.go) used `b.Store.Inspect` + `b.Count(state)`. That path
skips everything `CycleBinding.Inspect` (internal/budget/cycle_extension.go:145)
validates beyond the raw ledger read:

- `checkProtocolMigrationCharges` (internal/budget/protocol_migration.go:185) —
  a migrated binding must still hold every imported charge verbatim and its
  accounting origin; and
- extension-history consistency (`e.Spent > spent || e.At.Before(StartedAt)` →
  "cycle grant differs from retained charge history").

**Fix:** call the full read-only `b.Inspect(ctx)` and derive `Count` from the
returned `status.Spent`. A binding whose retained imported charge was removed or
rewritten now refuses at admission — before discovery — instead of surviving
until a later runner-side refusal.

**Semantic test:** build a real migrated cycle binding in the temp fixture via
the public `budget.InspectProtocolMigration` + `budget.MigrateProtocolBudget`
APIs (fixture-local only; never real amendment state), with one imported
cross-review charge, maximum frozen at 3. Preview is admissible (count 1 of 3).
Then rewrite the retained ledger so the imported charge is gone. Both preview
and `-execute -expect-digest` must refuse naming the imported charge, with the
discovery seam provably never invoked, no writes, and no cycle spent.

### 2. `configLayerPaths` hashed the explicit config only, not the effective env layer

When `-config` is omitted, `config.LoadAgentSpecs` reads the ambient
`PARLEY_HEADLESS_AGENT_CONFIG` layer (internal/config/runtime.go:392-397:
trimmed, relative paths joined to root, non-optional) — a layer the preview
digest never bound. Changing that file between preview and execute was
invisible to `-expect-digest`.

**Fix:** hash the **effective** env layer with the same path resolution
`config.configLayers` uses: the explicit `-config` when given (the launcher
overrides the env var with it before `LoadAgentSpecs`, so only that layer is
required — the ambient value it displaces is not), else the ambient
`PARLEY_HEADLESS_AGENT_CONFIG` value. The layer is required when set: the env
layer is non-optional in `LoadAgentSpecs`, so a missing file must refuse at
preview, not fail later at execute. Its hash line uses the
`PARLEY_HEADLESS_AGENT_CONFIG:<abs>` label so the report shows what was bound.

**Env hygiene:** execute applies the explicit layer via `os.Setenv`; the run
function now restores the previous ambient value (or unsets it) after the
launch, so one in-process run never leaks its override. No global file is ever
written.

**Test:** with the ambient env layer set and `-config` omitted, the preview
report's `config_hashes` contains the `PARLEY_HEADLESS_AGENT_CONFIG:<abs>`
entry; changing only that file after preview makes
`-execute -expect-digest` refuse with "digest changed" before discovery (spy on
the discovery seam), with no writes and no cycle spent.

### 3. The 900 s shared context started AFTER discovery

`executeRound` created `context.WithTimeout(ctx, groupedTimeout)` after spec
loading, discovery and resolution — only the grouped `RunRound` shared the
clock; all peer preparation ran on the unbounded caller context.

**Fix:** create the launch context at the top of `executeRound`, before the
pre-launch binding re-check and any execution-stage discovery, and pass it to
the binding re-inspection, discovery and the single `RunRound`. All peer
preparation and launch share one 900 s clock.

**Test:** a context-aware fake on the discovery seam records
`ctx.Deadline()` when discovery runs and asserts the deadline is already set
and ≈ `groupedTimeout` away (no sleep; the old ordering fails this because
discovery ran on a deadline-less context).

### 4. `organizer_preserved` reported the overall `kept` bool

`rep.OrganizerKept = &kept` bound the report field to the shared overall-success
variable (and, being a pointer to it, to its later mutations): any peer failure
rendered `organizer_preserved: false` even when the organizer bytes were
untouched.

**Fix:** compute organizer artifact equality (`preview hash == post-round
hash`) independently; `organizer_preserved` reports only that equality.
Overall success (`ok`, exit 2 path) still requires every participant to succeed
AND the organizer artifact to be unchanged.

**Failed-peer canary test:** one fake peer fails (helper exits non-zero,
writes no artifact); the round is unsuccessful (exit 2, `ok: false`) while
`organizer_preserved` is `true` and the organizer bytes on disk are identical.

## Gates retained exactly (no new CLI/product scope)

- exact `-v2` slug target; frozen four-member set, never replaced;
- skip semantics: organizer pre-existing artifact retained via `Overwrite=false`,
  spawn canary unchanged;
- exactly ONE `runner.RunRound` per execute ("oneRunRound");
- "no N" choice retained — no migration logic added to the launcher; the new
  test's migration calls are fixture-local public-API calls inside the test
  binary only;
- 900 s grouped timeout value unchanged; exit codes 0/1/2 unchanged.

## Boundaries

- Writes only the three allocated files (launcher `main.go`, `main_test.go`,
  this note). No Bash, no tests, no Git, no subagents, no providers, no real
  execute, no migration.
- Concurrent Kimi driver/session and Claude migration edits are not touched;
  targeted reads of published APIs only (`budget`, `config`, `runner`, `agents`).
- Old notes and their recorded failures are retained untouched.

## Status checklist (updated live)

- [x] Note written first, before touching the Go files.
- [ ] Fix 1: `inspectBinding` → full `b.Inspect`, `Count` from `status.Spent`.
- [ ] Fix 2: effective env layer hashed; env restored after execute.
- [ ] Fix 3: shared deadline before execution-stage discovery.
- [ ] Fix 4: `organizer_preserved` independent of overall success.
- [ ] Four new tests added (migrated-ledger corruption, env drift, deadline,
      failed-peer organizer canary); existing gates untouched.
- [ ] End-of-pass update appended below; no test claims until the coordinator
      executes.
