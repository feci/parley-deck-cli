---
agent: zcode-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-16
topic: grouped amendment launcher — correction pass after invocation timeout and coordinator compile failure
previous-invocation: 1d612975 (timed out 1800s, no own note)
---

# zcode-1 grouped launcher correction — 2026-09-16

## Why this note exists

The previous invocation (`1d612975`) timed out at 1800 s and closed with **no own
note**. This note is written FIRST, before any file is touched in this pass, to
restore the audit trail honestly. It will be updated at the end of the pass with
what was actually corrected. No test execution is claimed anywhere in this note;
the coordinator executes tests after this terminal handoff.

## Record of the failure being corrected

- **Timeout:** invocation `1d612975` hit its 1800 s ceiling and produced no own
  note. Both owned Go files under `.parley-runtime/grouped-amendment-launcher-20260916/`
  were retained exactly as written — no partial edits were lost, and none were
  silently altered.
- **Coordinator compile failure:** compiling the immutable published-recovery444
  tree plus my two Go control files failed in 0.95 s with two mismatch classes,
  reported by the coordinator:
  1. `main.go:251,256,259,260` reference `rep.PeerArtifactsPresent`, but the
     struct field is `PeerArtifacts`.
  2. `main_test.go:98` calls `discoverAgents` without the leading `context.Context`
     parameter the current signature requires.

## Correction scope for this pass (boundaries)

- Fix exactly those two compile mismatch classes in my two owned files:
  `main.go` and `main_test.go` under `.parley-runtime/grouped-amendment-launcher-20260916/`.
- Inspect ONLY these two files for other obvious faults. No broad repository
  rereading. Targeted lookup of the two mismatched symbols' declarations is the
  only external consultation needed to make the fix compile-correct.
- No Bash, no test execution, no Git operations, no subagents, no providers,
  no migration, no apply, no signoff in this pass.

## Invariants this launcher must hold (restated from the correction contract)

- **Exact target:** production currently permits `l.ideaSlug` other than the
  target idea whenever a matching member/status exists. That permissiveness is
  removed: the launcher runs ONLY `meta-protocol-change-evidence-first-efficiency-v2`.
  The isolated test fixture slug is adapted accordingly, and a wrong-slug refusal
  is covered.
- **Read-only preview**; a **missing binding never triggers discovery or writes**.
- Exactly **one** `runner.RunRound`, with a **900 s shared context**, using the
  **existing cycle only**.
- **Codex preserved**; **3 fake peers / one charge** in the fixture helper.
- **No real execute, no migration; "no N" choice retained.**
- Fake test helper must produce complete valid output.

## Attribution note (not mine)

Tracked-drift entries under `internal/app/budget*` / `internal/runner/telemetry*`
and related app entries were concurrent authorized Claude-owned changes, not
attributable to zcode-1. They are not edited in this pass.

## Status checklist (updated live)

- [x] Note written first, before touching the Go files.
- [x] Compile mismatches corrected in `main.go` / `main_test.go`.
- [x] Two-file inspection against the invariants above completed.
- [x] Update appended below with results; no test claims until the coordinator
      executes.

---

## Update — end of correction pass (2026-09-16)

### Compile mismatches fixed (exactly the two coordinator-reported classes)

1. `main.go` — the `report` struct field is `PeerArtifacts` (JSON tag
   `peer_artifacts_present`, unchanged for output compatibility). All four
   references at the former lines 251/256/259/260 now use `rep.PeerArtifacts`.
   No behavior change: append / length checks / `strings.Join` over the same
   slice of present peer artifact paths.
2. `main_test.go` — the injected fake in `withFakeDiscovery` now matches the
   real seam type `func(context.Context, []agents.Spec) []agents.Discovery`
   (verified against `agents.Discover` in `internal/agents/discover.go:476`).
   The production call site in `main.go` already passed `ctx` first and needed
   no change.

### Exact-target invariant enforced

- `collectPlan` now refuses any `-idea` slug other than
  `meta-protocol-change-evidence-first-efficiency-v2` as its **first** check,
  before any workspace read — so a different slug with a merely matching member
  set and `round-02` status can no longer launch, in preview or execute.
- `-idea` flag help and the file header now state the exact-target contract.
- `fixtureSlug` in the test now equals `targetIdea`: enforcement is slug-exact,
  so the isolated fixture carries the real target slug inside its temp
  workspace.
- New `TestRefusesWrongIdeaSlug`: the pre-v2 slug (closest possible wrong
  target, differing only by the missing `-v2` suffix) must be refused in both
  preview and execute with the exact target named, with no mutations and no
  cycle-ledger change.

### Two-file inspection results (other obvious faults)

- **Read-only preview** — intact: `collectPlan` only reads/hashes; no writes
  outside the single grouped `RunRound` in execute mode.
- **Missing binding never discovery/write** — intact: `inspectBinding` only
  loads/inspects and refuses on `nil`; execute re-checks admissibility
  immediately before any spec load, discovery, or spawn.
- **One `runner.RunRound`, 900 s shared context, existing cycle only** —
  intact and untouched (`groupedTimeout = 900 * time.Second`,
  `context.WithTimeout`, single call, no retries).
- **Codex preserved / 3 fake peers / one charge** — intact: `Overwrite=false`,
  organizer pre-validation and pre/post hash comparison, spawn canary, cycle
  before=1 → after=2 assertion.
- **Fake helper output verified complete and valid** against the real
  validator (`internal/runner/validation.go`): round ≥2 requires matching
  `agent`/`idea`/`round` frontmatter plus a section with content — the peer
  helper body satisfies it; the fixture round-01 files carry the four required
  sections plus a non-empty `## Existing alternatives` (§15.6a). The helper's
  prompt-parsing regexes are grounded in the actual runner prompt
  (`internal/runner/runner.go:853`, `:1011`, required-file-shape block at
  `:1021-1025`).
- No migration logic exists or was added; "no N" choice unchanged; no other
  faults found in the two files.

### Boundaries respected

Only the three allocated files were written. No Bash, no test execution, no
Git, no subagents, no providers, no migration, no apply, no signoff. Symbol-level
lookups (`agents.Discover`, runner prompt lines, `ValidateRoundArtifact`) were
the only reads outside the two owned files — the minimum needed to make the
fixes compile-exact and to verify the helper output, not a broad reread.

### Honest test status

**No tests were run and none are claimed.** Compilation and test outcomes are
unknown until the coordinator executes them after this terminal handoff.
