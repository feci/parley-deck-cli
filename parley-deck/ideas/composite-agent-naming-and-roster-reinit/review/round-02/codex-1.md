---
agent: codex-1
idea: composite-agent-naming-and-roster-reinit
review-round: 2
date: 2026-07-19
reviewed-commit: 31fd5c733968f8ae102953081188f58f94ac7703
responding-to: [codex-1/review/round-01, kimi-1/review/round-01]
---

## Summary

The fix-up resolves the clean-build blocker, the concrete round-one quorum omission,
participant traversal, canonical display-name parsing, explicit roster-ID selection,
preflight, consensus drafting, signoff collection, steering, and the valid-mapping forms of
the implementation safety gates. The focused tests, build, vet, formatting, and diff checks
are green.

The branch is still not safe to accept. The ratified fail-closed autonomous-write honesty
rule remains knowingly violated: claude, agy, and hermes still advertise
`Scope="workspace"` although their bypass/yolo processes are not workspace-confined. The
fix-up also leaves the default `parley run` path on raw family identities and introduces
machine-scope and quoted-TOML failure modes in `roster init`. These are observable behavior,
not merely missing convenience commands.

## Refutation attempts

- Read the supplied 962-line `fixup.diff` from `cdef991` through HEAD and inspected every
  changed resolver, app, config, runner, and test hunk against the current source.
- Traced all callers of `ResolveParticipant`, `selectedParticipantIDs`, `selectedAgents`,
  `rosterMappingFor`, `modelOf`, and `discoveryFor`, including the graphified call paths for
  app selection, preflight, Phase 5-8, consensus/signoff, and steering.
- Replayed the original explicit roster-ID scenarios: selection, preflight, drafter choice,
  signoff request, steering, LE-3 model comparison, and LE-7 discovery now reach the mapped
  family while retaining the roster ID as identity.
- Checked traversal inputs against `pathSafeParticipant`, canonical/non-canonical display
  spellings against the re-compose check, exact-ID specs with an existing `AdapterID`,
  target-file idempotency, invalid scope, JSON mutation, and repeat writes.
- Challenged `roster init --scope machine` with conflicting central/session layers and the
  empty-adapter case with a valid quoted TOML table key. Both expose gaps described below.
- Ran `go build ./...`, `go vet ./...`, `git diff --check cdef991..HEAD`, focused
  `go test ./internal/agents ./internal/app ./internal/config`, and `gofmt -l` over every
  changed Go file: all passed/clean.
- Ran `go test ./...`: every package passed except the known environment-dependent
  `internal/runner.TestDurableKillEndToEndRealProcess`, which failed because this host has no
  recorded boot ID. This is the same unrelated failure recorded in round 1 and is not
  attributed to this fix-up.

## Responses to prior reviewers

### @codex-1 — review/round-01

The compile, round-one fail-open, traversal, parse canonicality, and explicit-adapter
findings are fixed as claimed. The roster writer is materially safer, and the explicit
roster-ID app paths now work. I do not withdraw the autonomous-confinement finding: neither
the supplied diff nor current launch code adds confinement or stops making the workspace
claim. The original app-wide finding is also only partially closed because the no-flag
default selection still emits family IDs.

### @kimi-1 — review/round-01

I concur that mapped preflight, drafter, signoff, steer, legacy path-safe IDs, and the valid
mapping forms of `modelOf`/`discoveryFor` are fixed. I also concur with Kimi's autonomous
honesty analysis. The empty-adapter fix covers the exact unquoted spelling in its new test,
but not equivalent valid TOML spellings; the residual can atomically install an invalid
duplicate table. Kimi's two NITs remain unchanged.

## Prior finding dispositions

### [CRITICAL] RESOLVED — Clean-checkout compile blocker

The supplied diff's `internal/agents/naming.go` hunk adds `RenderDisplayName`, and its
`naming_test.go` hunk adds the five principal rendering cases. The symbol is committed at
`internal/agents/naming.go:188`, and `go build ./...` now succeeds at HEAD.

### [CRITICAL] RESOLVED — Round one no longer completes with unresolved quorum members

The supplied `internal/runner/runner.go` hunk changes `selectedAgents` to return an
`unresolved` list and appends one failed `Result` per unresolved participant. Current
`runner.go:204-237` therefore makes mixed or zero-resolved rounds `round.incomplete` with
`FailureClass="roster-unresolved"`. Phase 6 inherits this through `RunReviewRound`.

### [CRITICAL] RESOLVED — Participant path traversal

The supplied `internal/agents/resolve.go` hunk adds `pathSafeParticipant` and invokes it
before both exact-ID and mapped resolution. Current `resolve.go:15-25,43-46` rejects slash,
backslash, `..`, and edge-dot traversal forms while preserving safe legacy IDs. The added
`TestResolveRejectsTraversal` covers exact/mapped unsafe inputs and safe underscore,
uppercase, and single-dot legacy IDs.

### [CRITICAL] OPEN — Autonomous-write workspace confinement is still falsely asserted

There is no confinement fix in `fixup.diff`. Current `internal/agents/discover.go:193-194,
218-219,267-268` still sets `Scope: "workspace"` for claude, agy, and hermes. Their effective
controls are `--add-dir`, cwd, and bypass/yolo flags; those grant/select context but do not
prevent an absolute or parent-directory write. `Declared()` at `discover.go:97-99` trusts
the string and makes `roster show` report `AUTO=yes`.

This directly contradicts FINAL/consensus §C: unverified confinement must leave the bit
unset. Calling the limitation "flag-scoped" in `IMPLEMENTATION.md` documents the
disagreement but does not satisfy or amend the ratified safety contract. Set `Scope` empty
for these adapters until an enforced sandbox exists, or add real OS confinement and a
negative outside-workspace test. Codex may retain the declaration for its effective
`workspace-write` sandbox.

### [MAJOR] OPEN — Roster resolution is still not app-wide on the default run path

The supplied diff correctly wires explicit IDs through `selectedParticipantIDs`,
`participantDiscoveries`, `firstHeadlessAgent`, `requestSignoffAgents`, and `discoverAgent`.
However, current `internal/app/app.go:2391-2394` returns `installedAgentIDs(discovered)`
without consulting the mapping whenever `--participants`/a preset is absent. `runTask` then
persists those raw families in `runcontrol.Create` at `app.go:1839-1844`.

Thus a normal `parley run --yes TASK` in a deck containing
`[roster.codex-1] adapter="codex"` creates a new idea/run with participant `codex`, not the
ratified identity `codex-1`; artifact paths, signoffs, runstate, and quorum all remain in the
old namespace. Resolve the ordered active §2 roster for default selection, falling back to
installed family IDs only for a legacy deck with no roster mapping. Add a no-flag end-to-end
test; `TestAppLevelRosterIDResolution` currently exercises only the explicit raw argument.

### [MAJOR] RESOLVED — Explicit roster-ID preflight is no longer vacuous

The supplied `internal/app/preflight.go` hunk routes each selected participant through
`ResolveParticipant` before `checkRoster`. For a valid `codex-1 -> codex` mapping, the report
now contains `codex-1`, the hosted probe runs with the codex adapter, and the `<2` non-solo
hard stop applies. This closes Kimi's concrete one- and two-participant scenarios.

### [MAJOR] RESOLVED — Direct roster-init scope/JSON/idempotency/atomic-write cases

The supplied `internal/app/roster.go` hunk validates `session|machine`, compares against
`RosterAdaptersInFile(target)`, mutates before JSON rendering, reports the real outcome, and
replaces the live file through `fsutil.WriteFileAtomic`. The added tests cover invalid scope,
`--yes --json`, an effective unknown adapter, an unquoted empty adapter, and repeat
idempotency. These close the exact round-1 misspelling, JSON no-op, inherited-idempotency,
and partial-append scenarios, subject to the scope/representation regressions below.

### [MINOR] RESOLVED — Parse accepts only canonical display spellings

The supplied `internal/agents/naming.go` hunk re-composes the parsed value and requires byte
equality at current `naming.go:170-177`. The added tests reject `x-high`, lowercase `xhigh`,
and `_02`; canonical round trips remain green.

### [MINOR] RESOLVED — Exact-ID resolution preserves an explicit adapter

The supplied resolver hunk captures `d.Spec.Adapter()` before replacing `Spec.ID`
(`internal/agents/resolve.go:47-52`), and `TestResolvePreservesExplicitAdapter` proves
`claude-1` retains adapter `claude`. Config/spec construction de-duplicates IDs before
discovery, so the original first-match ambiguity is not reachable through the app loader.

### [MINOR] RESOLVED — Legacy safe participant IDs are no longer rejected

The strict canonical-roster regex was replaced by the containment predicate above. The
tests demonstrate `my_cli`, `Claude-1`, and `claude.v1` still resolve by exact configured ID
without reopening separator traversal.

### [MINOR] RESOLVED — Valid roster mappings activate LE-3/LE-7 lookups

The supplied `internal/app/driver_impl.go` hunk maps roster ID to family in both `modelOf`
and `discoveryFor` (`driver_impl.go:124-131,361-368`). With a valid mapping, model diversity
is evaluated and the goal checker is discovered rather than silently skipped. This fixes
Kimi's stated valid-roster scenario, although the implementation should eventually call the
shared resolver instead of duplicating a permissive lookup.

### [MINOR] OPEN — `roster diff` and TOML `autonomous_write` override remain deferred

Current `internal/app/roster.go:45-77` still exposes only `show|init`, and
`internal/config/runtime.go:97-132` has no `autonomous_write` override. The new
`IMPLEMENTATION.md` deferred section makes the deviation visible, but the two ratified
surfaces remain unimplemented. They may be accepted only as explicit follow-up scope, not
described as delivered.

### [NIT] OPEN — Unresolved pseudo-results still emit no per-agent failure event

Current `internal/runner/runner.go:220-227` appends only a `Result`; segment targets at line
205 omit unresolved IDs, and there is no matching `agent.failed` event before the aggregate
`round.incomplete`. Runstate/TUI therefore cannot show the participant-level
`roster-unresolved` failure even though CLI/index output can. Kimi's suggested projection
parity fix remains unapplied.

### [NIT] OPEN — `selectedAgents` comment still describes fail-open skipping

Current `internal/runner/runner.go:358-362` still says unresolved participants are skipped,
even though `RunRoundOne` now materializes failed results. Update the contract and explain
that single-agent Phase 5/8/consensus callers fail when selection is empty.

## New fix-up regressions

### [MAJOR] OPEN — Machine-scope init reads project overrides and can copy them upward or mask a broken machine mapping

`rosterInit` calls layered `resolveRoster(root)` first (`internal/app/roster.go:207`), and
that helper loads central + deck + local + environment specs/mappings at lines 96 and 106.
Only afterward does init select the machine target and parse that target at lines 212-215.
It then skips by target key presence alone at lines 227-230.

Two fail-closed violations follow:

1. If the machine target is missing `claude-1` but the deck maps it to a project-only
   adapter, `--scope machine --yes` writes that deck-derived adapter into
   `~/.parley/agents.toml`, contrary to FINAL §B's "never copies deck values up" rule.
2. If the machine target contains `adapter="claud"` but the deck overrides the same ID with
   valid `adapter="claude"`, layered resolution validates the deck value, target parsing
   sees the key, and init reports `unchanged` while the requested machine scope remains
   broken outside this deck.

Build machine proposals and validation without higher project layers; validate the target's
actual adapter value against the machine family catalog rather than only checking key
presence. Add conflicting-layer tests in both directions.

### [MAJOR] OPEN — Empty-adapter protection is textual and can atomically install duplicate invalid TOML

`RosterAdaptersInFile` intentionally omits empty adapters. `writeRosterMappings` then tries
to detect that omitted table with the literal substring
`"[roster."+id+"]"` (`internal/app/roster.go:311-329`). A valid equivalent table such as
`[roster."claude-1"]` with `adapter=""` does not contain that spelling. Init appends
`[roster.claude-1]`, atomically replacing the target with two declarations of the same TOML
table and reporting success; the next config load fails.

Determine table presence structurally, including empty values and quoted keys, and validate
the complete candidate TOML before rename. Add quoted-key and dotted/inline-equivalent cases;
the current test covers only the writer's exact unquoted output spelling.

### [MINOR] OPEN — `RenderDisplayName` applies the agy parenthesized-tier rule to every family

The newly committed function calls `StripParenTier(label)` unconditionally and substitutes
the tier whenever reasoning is invalid/`cli-default` (`internal/agents/naming.go:188-198`).
The comments, FINAL, and consensus make this an agy-only rule. For example, a non-agy spec
with `ModelLabel="Model (High)"` and `Reasoning="cli-default"` is rendered as
`family_model_high`, silently moving a model qualifier into effort, rather than preserving
the qualifier in the model and rendering `cliDefault`.

Gate tier extraction/substitution on `family == "agy"`; sanitize the full label for other
families. Add the non-agy parenthesized-label and empty/invalid-label cases requested by the
round-1 compile finding—the supplied rendering test includes neither.

### [MINOR] OPEN — Phase-8 implementation metadata was appended but not advanced

The supplied `IMPLEMENTATION.md` diff appends cycles 1-3, but its frontmatter still says
`status: implemented` and has no `completed` or `head-commit`. Phase 8 requires
`status: fix-up-cycle-3` plus the reviewed HEAD until a zero-fix consensus permits
`status: complete`. The stale resume/audit metadata does not describe the code this review
actually inspected.

## Verdict

**BLOCK**
