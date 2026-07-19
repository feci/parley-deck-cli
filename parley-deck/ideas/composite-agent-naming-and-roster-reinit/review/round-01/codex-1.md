---
agent: codex-1
idea: composite-agent-naming-and-roster-reinit
review-round: 1
date: 2026-07-18
reviewed-commit: cdef99120072f2a57933a57569bb19fe65217e80
---

## Summary

The committed implementation is not safe to merge. The identity/adapter split is correct at the three requested runner vendor-dispatch sites (`cleanParticipantEnv`, the Hermes environment, and `isolatedAgentHome`), and the speed field does not map to model or effort. However, the reviewed commit does not compile, unresolved roster identities are silently omitted from rounds, participant IDs can escape artifact directories, and several app-level paths still bypass roster resolution. The roster writer and autonomous-write declaration also have fail-closed and safety defects.

S6 is intentionally deferred per `IMPLEMENTATION.md` and the review brief; it is not reported as a finding.

## Refutation attempts

- Inspected the committed files with `git show HEAD:<path>` because the shared working tree contains uncommitted edits to `internal/agents/naming.go` and `naming_test.go`. Those edits are not part of reviewed commit `cdef991` and were not credited to it.
- Traced artifact identity through `selectedAgents` and `runAgent`: resolved roster IDs do key run logs, artifacts, frontmatter validation, and events. Traced vendor dispatch separately: `cleanParticipantEnv(agent.Adapter(), ...)`, the Hermes branch, and `isolatedAgentHome` all use `Adapter()`.
- Exercised resolver cases mentally and against its tests: exact legacy IDs and explicit mappings work, and the source slice is not mutated. Error propagation and pre-existing `AdapterID` preservation do not.
- Checked naming with the five amended examples, edge dots, repeated dots, all-digit models, instance suffixes, and non-canonical effort spellings. Generated names are path-safe, but `Parse` accepts multiple spellings for the same parsed value.
- Grepped every Go use of `Speed`: it is loaded and displayed independently; no current path maps it to `Model` or `Reasoning`.
- Audited roster initialization for a missing map, an invalid existing adapter, JSON writes, machine/session scope isolation, repeat execution, concurrent execution, and interrupted append.
- Ran `go test ./...` on the shared working tree only as supplementary evidence. The uncommitted naming patch makes the new packages compile there; relevant agent/app/config tests passed, while an unrelated pre-existing runner durable-kill test failed because the environment had no recorded boot ID. This does not change the clean-commit build blocker below.

## Findings

### CRITICAL

#### [CRITICAL] The reviewed commit does not compile

**File:** `internal/app/roster.go:129`

`resolveRoster` calls `agents.RenderDisplayName`, but commit `cdef991` contains no definition of that symbol: committed `internal/agents/naming.go` ends with `isAllDigits`, and the committed naming tests do not define or exercise rendering.

**Failing scenario:** On a clean checkout of `cdef991`, `go build ./...` or `go test ./...` fails with `undefined: agents.RenderDisplayName`. The function currently visible in the shared working tree is an uncommitted change and therefore cannot make the reviewed revision releasable.

**Fix:** Commit a `RenderDisplayName` implementation and focused tests for all five amended display names, the agy parenthesized-tier rule, empty/invalid model labels, and non-agy parenthesized labels. Re-run build, vet, formatting, and all tests from a clean tree.

#### [CRITICAL] Resolution errors are discarded, so a round can complete with part or all of quorum missing

**File:** `internal/runner/runner.go:352`

`selectedAgents` ignores every `ResolveParticipant` error. `RunRoundOne` then sizes its result set from only the successfully resolved agents and emits `round.completed` if those agents succeed. If no participant resolves, it emits a completed round with `total: 0`. The injected loader also discards `LoadRosterAdapters` errors at `internal/app/roster.go:24`.

**Failing scenario:** An idea has `participants: [claude-1, codex-1]`, both CLIs are installed, but the mapping contains only `claude-1 -> claude`. The runner launches Claude, silently omits Codex, writes a one-agent round index, and can emit `round.completed`. With no mappings it can complete a zero-agent round. This violates the fail-closed resolver contract and the locked quorum.

**Fix:** Make mapping loading and participant selection return errors. Resolve every participant before launching any process; if any participant is unresolved, abort the phase or emit an explicit failed result for that identity and ensure the phase is `round.incomplete`. Add mixed-resolved, zero-resolved, malformed-config, and Phase 5/6/7 propagation tests.

#### [CRITICAL] Unvalidated participant identities permit path traversal in run logs and protocol artifacts

**Files:** `internal/agents/resolve.go:20`; `internal/runner/runner.go:364`; `internal/runner/runner.go:367`

`ResolveParticipant` copies an arbitrary participant string into `Spec.ID`, and `runAgent` uses it directly in `filepath.Join` for both the run directory and `<agent>.md`. The §2 reader validates roster rows, but neither the resolver nor the run path proves that an idea participant/mapping key came through that reader. Manually edited frontmatter and quoted TOML keys can therefore contain `..` path segments.

**Failing scenario:** A deck contains `[roster."../../../../tmp/reviewer"] adapter = "claude"` and an idea lists that same participant. Resolution succeeds, after which the runner creates logs below the escaped directory and prompts Claude to write an artifact at a cleaned path outside `review/round-NN`. A malicious repository can make the host CLI write outside the deck.

**Fix:** Validate participant identities before resolution and again before filesystem use with the canonical roster-ID grammar (`[a-z0-9][a-z0-9-]*`), rejecting separators, edge dots, and `..`. Add a containment check after joining paths as defense in depth. Test both exact-ID and mapped traversal inputs.

#### [CRITICAL] `AutonomousWrite.Scope="workspace"` is asserted without workspace confinement

**Files:** `internal/agents/discover.go:97`; `internal/agents/discover.go:193`; `internal/agents/discover.go:219`; `internal/agents/discover.go:268`; `internal/runner/runner.go:1023`

`Declared` trusts only a non-empty mode plus the literal scope string. Claude and agy are marked workspace-confined because they carry `--add-dir {root}`, which grants/adds access rather than enforcing an OS sandbox; Hermes is marked workspace-confined with only `--yolo`. Setting `cmd.Dir = root` changes the working directory but does not prevent writes elsewhere. The static test checks declarations and flags, not confinement.

**Failing scenario:** A Hermes or Claude participant running in the declared autonomous mode executes a tool action against `../outside`, an absolute path, or a user configuration file. Nothing in the runner constrains the process to `root`, yet `roster show` reports `AUTO=yes`, contradicting the ratified fail-closed honesty rule.

**Fix:** Leave `Scope` empty/unverified for adapters without demonstrable confinement, or launch them through an enforced workspace sandbox. Make `Declared` depend on verified enforcement rather than a self-asserted string, and add a negative outside-workspace sentinel test. Codex's `workspace-write` sandbox can remain declared when its effective invocation is verified.

### MAJOR

#### [MAJOR] Roster resolution is not wired through the rest of the workflow

**Files:** `internal/app/app.go:2391`; `internal/app/preflight.go:206`; `internal/app/driver_consensus.go:97`; `internal/app/consensus_request_signoffs.go:251`; `internal/runner/steer.go:270`

Only runner phase selection uses `ResolveParticipant`. New-run selection, preflight, consensus drafting, signoff requests, implementation metadata/model lookups, goal checks, and live steering still compare roster strings directly with raw discovery IDs. Those raw IDs are families such as `claude`, not stable identities such as `claude-1`.

**Failing scenario:** With a valid `claude-1 -> claude` mapping, `parley run --participants claude-1 ...` rejects `claude-1` as not installed. If an idea is seeded manually, its initial runner round can launch, but auto-consensus reports no headless participant, `consensus request-signoffs` reports no configured runner entry, and a TUI steer to `claude-1` is rejected as unknown.

**Fix:** Introduce one shared, error-returning roster-resolution boundary and use resolved `Discovery` values everywhere participant identities enter app/driver/steer code. Preserve roster ID as identity and adapter as launch family. Add an end-to-end roster-ID test covering preflight, round, consensus/signoff, review, and steering; retain a nil-mapping legacy spec-ID case.

#### [MAJOR] Existing but invalid mappings are reported as initialized

**File:** `internal/app/roster.go:111`

When a mapping exists, `resolveRoster` does not verify that its adapter names a configured family. It indexes `byFamily[family]`, receives a zero-value spec on a typo, falls back only for display, and leaves `Family` non-empty. `rosterInit` then skips the ID solely because it exists in the effective mapping.

**Failing scenario:** `[roster.claude-1] adapter = "claud"` makes `roster init --yes` print “already initialized” and exit 0, although `ResolveParticipant` will later fail because `claud` is not discovered.

**Fix:** Treat a mapping whose adapter is absent from the configured/discovered family catalog as unresolved and exit nonzero. Validate every existing mapping, not only proposed additions, and test typoed and unavailable adapters.

#### [MAJOR] Scope idempotency is computed from all layers instead of the target file

**Files:** `internal/app/roster.go:190`; `internal/app/roster.go:248`

`rosterInit` calls layered `LoadRosterAdapters` to decide whether the selected target needs writing. A mapping inherited from another layer therefore suppresses a write to the requested scope. Scope values are also not validated; anything other than exactly `machine` silently selects the session file.

**Failing scenario:** The deck `agents.toml` already contains all mappings and the central file contains none. `parley roster init --scope machine --yes` reports the machine roster already initialized and leaves `~/.parley/agents.toml` unchanged. Separately, `--scope machien --yes` silently appends to the deck file.

**Fix:** Validate `scope` against `session|machine`; parse and compare the chosen target file for idempotency, using layered data only as proposal defaults. After writing, reload that target and verify the requested entries are present.

#### [MAJOR] `--json` turns roster initialization into an undocumented no-op

**File:** `internal/app/roster.go:211`

The JSON branch returns before confirmation and write handling. It does so even when `dry_run` is false and `--yes` is present, so structured-output callers receive exit 0 and a payload naming proposed additions while no state changes.

**Failing scenario:** `parley roster init --scope session --yes --json` outputs `"dry_run": false` and the expected `adds`, exits successfully, but does not create or update `parley-deck/agents.toml`. A script can proceed under the false assumption that roster resolution is ready.

**Fix:** Separate mutation from rendering. Perform the same validation/confirmation/write path in text and JSON modes, then serialize the actual outcome (`written`, `unchanged`, or `dry-run`). Add a test that reloads the mapping after `--yes --json`.

#### [MAJOR] Direct TOML append is neither atomic nor concurrency-safe

**File:** `internal/app/roster.go:266`

`appendToFile` opens the live configuration with `O_APPEND` and writes the whole block directly. Two concurrent initializations can both observe a missing entry and append duplicate `[roster.<id>]` tables; interruption or a short write can leave invalid TOML. The next configuration load then fails for every run using that layer.

**Failing scenario:** Two `roster init --yes` processes start against an unmapped deck. Both compute the same `toWrite`, both append, and the resulting file contains duplicate table declarations that `go-toml` rejects.

**Fix:** Lock the target, re-read and parse it after acquiring the lock, apply missing entries structurally, and write a complete validated document through the repository's atomic-write helper. Verify the newly written TOML before replacing the original. Add concurrent and injected-write-failure tests.

### MINOR

#### [MINOR] `Parse` accepts non-canonical aliases for one display name

**File:** `internal/agents/naming.go:146`

The effort section is normalized by deleting punctuation, and the instance suffix accepts any digit string whose numeric value is at least two. Consequently distinct strings parse to the same `ParsedName` and do not round-trip exactly.

**Failing scenario:** `Parse("codex_gpt-5.6-sol_x-high")` succeeds as effort `xhigh`, and `Parse("codex_530_xHigh_02")` succeeds as instance `2`; composing either result produces a different canonical name. That violates the documented fail-closed grammar and makes collision/name comparisons ambiguous.

**Fix:** After parsing, re-compose and require byte-for-byte equality with the input, or explicitly validate the effort display token and canonical decimal instance spelling. Add these two cases to `TestParseFailsClosed`.

#### [MINOR] Exact-ID resolution destroys an already explicit adapter

**File:** `internal/agents/resolve.go:23`

On an exact match, the resolver unconditionally assigns `AdapterID = Spec.ID` instead of preserving `Spec.Adapter()`. This violates the split for callers that already hold a stable-ID spec with a distinct adapter.

**Failing scenario:** Resolving participant `claude-1` from `Discovery{Spec: Spec{ID: "claude-1", AdapterID: "claude"}, Found: true}` returns `Adapter() == "claude-1"`. Claude environment cleaning and any other family dispatch are then skipped.

**Fix:** Capture and preserve the existing effective adapter before assigning the participant identity. Also fail on multiple found discoveries with the same match key rather than returning the first one. Add exact-ID-with-adapter and duplicate-match tests.

### NIT

No NIT findings. This review intentionally excludes stylistic preferences.

## Open questions

None. The failures above have concrete fixes and do not depend on S6.

## Overall verdict

**BLOCK** — the committed revision is uncompilable and has multiple fail-open identity and workspace-safety defects.
