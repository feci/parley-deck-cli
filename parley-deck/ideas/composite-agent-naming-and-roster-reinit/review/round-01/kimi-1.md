---
agent: kimi-1
idea: composite-agent-naming-and-roster-reinit
review-round: 1
date: 2026-07-18
reviewed-commit: 8130afcc6b9cdc7921071221bbced0efa72bf869   # code state = 54b5282 (fix-up cycle 1); 8130afc touches docs only
---

## Summary

I reviewed the branch at HEAD `8130afc` (the code is unchanged since `54b5282`, the
fix-up commit for codex-1's round-01 review). Note: the `impl.diff` handed to reviewers
predates that fix-up — I reviewed the working tree / HEAD, not the stale diff.

The core of the feature is sound: the identity/adapter split is correct at every
vendor-dispatch site I could find, the resolver is fail-closed and traversal-safe, the
naming package survives every edge input I threw at it, speed never maps to
model/effort, and legacy family-ID decks are untouched. codex-1's four CRITICALs are
either verifiably fixed (compile blocker, path traversal, fail-open rounds) or
documented as contested/deferred (autonomous-write confinement).

What remains: the roster-ID wiring stops at the runner. Three app/driver paths still
compare roster IDs against raw family discovery IDs, so a roster-ID idea cannot
complete the auto flow (consensus drafting and signoff collection hard-fail; TUI steer
is rejected). Worse, preflight's §1 non-solo hard-stop and the §9.0 availability ping
go silently vacuous for roster-ID participant sets. These were acknowledged as a
deferred follow-up in IMPLEMENTATION.md, but the deferral note understates the safety
angle. S6 deferred is intentional per the brief and is not reported.

Verification performed: `go build ./...`, `go vet ./...` clean; `go test ./...` green
(25 pkgs, 0 fail) at HEAD; `gofmt -l` on the changed packages flags only
`internal/app/preflight{,_test}.go`, which are equally unformatted on `main`
(pre-existing, not this branch). I also ran `parley roster show` and
`parley roster init --dry-run` in this real deck: the four §2 rows render with the
amended grammar (`claude_opus-4.8-1m_max`, `agy_gemini-3.5-flash_high`,
`codex_cli-default_cliDefault` — the honest bootstrap form, `hermes_glm-5p2_cliDefault`),
the dry-run proposes the correct four `[roster.*]` mappings and writes nothing.

## Refutation attempts (what I tried to break, and could not)

- **Vendor dispatch keyed on `agent.ID`:** grepped every family-name comparison in
  `internal/runner` and `internal/app`. `cleanParticipantEnv` (runner.go:1022), the
  hermes env branch (runner.go:1104), and `isolatedAgentHome` (runner.go:1154) all key
  on `agent.Adapter()`. The only remaining family-name switches are on raw,
  pre-resolution discoveries where ID == family by construction (`probePrompt`'s
  codex case, app.go:2110; preflight's gemini exclusion). `runAgent` keys artifacts,
  frontmatter validation, review snapshots, events, and the tracker on the roster ID.
  No missed vendor-dispatch site found.
- **Resolver fail-closedness:** exact-ID, mapping, unknown, uninstalled-mapped-family,
  and non-canonical-ID inputs all behave; the input slice is not mutated (code +
  `TestResolveDoesNotMutateInput`); the mapping family value is never path-joined.
- **Naming:** all-digit model + instance disambiguation (`codex_530_xHigh` vs
  `codex_530_xHigh_2`), `..`/edge-dot rejection, leading-zero instance, non-canonical
  effort spellings (`x-high`, lowercase `xhigh`), empty sections, unicode labels,
  effort-in-vocabulary enforcement (`fast`/`deep` as "effort" fails closed to the
  roster-ID fallback). `SanitizeSection` output always satisfies `sectionRe`
  (path-safe: no `..`, no edge dots, ASCII only).
- **Speed invariance:** every Go use of `Speed` (config load, defaults fan-out in
  `applyFile`, the prompt's "speed:" line at runner.go:854) touches only `Spec.Speed`;
  no speed→model/effort path exists. Guard test passes.
- **Legacy nil-mapping decks:** exact-ID resolution is the pre-change behavior;
  runner/app/config suites (which exercise family-ID fake agents) all pass.
- **roster init:** idempotent against the target file, atomic write, fail-closed on
  unresolved ids, `--json` performs and reports the real outcome, typoed adapter in an
  existing mapping fails closed (roster.go:125-131).

## Findings

### CRITICAL

None remaining. codex-1's compile blocker, path traversal, and fail-open round
completion are verifiably fixed at HEAD; the fourth (autonomous-write confinement) I
re-raise as MAJOR below with my own reasoning.

### MAJOR

#### [MAJOR] Roster-ID resolution stops at the runner — consensus drafting, signoff collection, and TUI steer still match raw family IDs

**Files:** `internal/app/driver_consensus.go:97-108` (`firstHeadlessAgent`, the match at :103); `internal/app/consensus_request_signoffs.go:251-268` (`requestSignoffAgents`, the lookup at :258); `internal/runner/steer.go:270-277` (`discoverAgent`, the match at :272)

All three compare participant strings (roster IDs like `claude-1`) against
`Discovery.ID`, which after the split is the *family* (`claude`). This is the
un-fixed remainder of codex-1's round-01 MAJOR ("not wired through the rest of the
workflow"); 54b5282 wired only `selectedParticipantIDs`.

**Failing scenario:** Deck with `[roster.claude-1] adapter = "claude"` (written by
`roster init`) and an idea with `participants: [claude-1, codex-1]`:

1. `parley run --auto --participants claude-1,codex-1` — selection, preflight (see
   next finding), and all rounds work; then the driver's consensus phase calls
   `runDrafter` → `firstHeadlessAgent(discovered, [claude-1, codex-1])` → no
   `agent.ID` equals a roster ID → `Draft` errors with "no headless idea participant
   available to draft consensus". **A roster-ID idea cannot complete the auto flow.**
2. `parley consensus request-signoffs <idea>` → targets are the idea's roster IDs
   (validated against the 00-prompt participant set) → `byID["claude-1"]` misses →
   hard error "participant claude-1 has no configured runner entry". Signoff
   collection is impossible for roster-ID ideas.
3. TUI steer to `claude-1` during a live run → the participants check
   (steer.go:131) passes but `discoverAgent("claude-1")` misses → rejected as
   "unknown agent claude-1". (Even if it matched, `BuildSteerPrompt` at steer.go:298
   would tail `runs/<id>/agents/claude/stdout.log` while the round wrote under
   `agents/claude-1/` — the fix must use the resolved Discovery, whose ID is the
   roster ID.)

**Fix:** Resolve once at each boundary with `agents.ResolveParticipant(id, discovered,
config.LoadRosterAdapters(root))` and use the returned Discovery (roster ID as
identity, `Adapter()` for launch): in `firstHeadlessAgent` (match participants
against resolved identity), in `requestSignoffAgents` (resolve each target, error if
unresolvable — it already errors on unknown, so fail-closed semantics are preserved),
and in `discoverAgent`. Add an end-to-end roster-ID test that drives a run through
round-01 → consensus draft → signoff request → steer.

#### [MAJOR] Preflight's §1 non-solo hard-stop and §9.0 availability ping are silently vacuous for roster-ID participant sets

**Files:** `internal/app/preflight.go:206-218` (`participantDiscoveries`), hard-stop guard at `internal/app/preflight.go:337`

`runTaskPreflight` (preflight.go:233) passes the selected participant IDs — which may
now be roster IDs after the 54b5282 fix to `selectedParticipantIDs` — into
`participantDiscoveries`, which matches them against raw discovery (family) IDs. For
roster-ID participants the result is an **empty** slice. `preflight` then computes
`report.Roster` from that empty set, and the non-solo hard-stop at :337 is guarded by
`len(report.Roster) > 0`, so it never fires; no ping probes run either.

**Failing scenario:** `parley run --participants claude-1 --yes` (a SOLO roster-ID
run) with a valid mapping: selection succeeds, preflight evaluates zero agents,
returns exit 0 "ready", and the run launches one participant — the exact §1 non-solo
violation preflight exists to hard-stop. With two roster-ID participants the run
proceeds but the ratified §9.0 full hosted-PONG ping of the actual participant set is
skipped entirely, so an unreachable/broken family is discovered only at round time.
Both failures are *silent* — no warning that the readiness check evaluated nothing.

**Fix:** Resolve participant IDs through the `[roster.*]` mapping (same resolver)
before building the preflight set, so `report.Roster` is non-empty and both the
`< 2` hard-stop and the ping apply. As defense in depth, treat an empty resolved set
with non-empty input as a preflight error, not as vacuous readiness.

#### [MAJOR] `AutonomousWrite.Scope = "workspace"` is asserted for claude/agy/hermes without demonstrated confinement (ratified honesty rule violated)

**Files:** `internal/agents/discover.go:194` (claude), `internal/agents/discover.go:219` (agy), `internal/agents/discover.go:268` (hermes); `Declared()` at `internal/agents/discover.go:97-99`

Consensus §C ratified the fail-closed honesty rule: where workspace confinement
cannot be demonstrated, the bit stays unset. For codex, `--sandbox workspace-write
--cd {root}` is a real OS-enforced sandbox — demonstrated. For claude
(`--permission-mode bypassPermissions --add-dir {root}`), agy
(`--dangerously-skip-permissions --add-dir {root}`), and hermes (`--yolo`), nothing
confines writes to the workspace: `--add-dir` *grants* access, `cmd.Dir = root`
(runner.go:1059) sets cwd but no jail, and the bypass flags then auto-approve every
tool call. A participant in these modes can write anywhere the OS user can, yet
`Declared()` returns true and `roster show` / the runtime matrix print `AUTO=yes`
with Scope "workspace". IMPLEMENTATION.md contests this ("flag+cwd is the intended,
documented mechanism; AUTO=yes means flag-scoped, not OS-jailed"), but the ratified
contract says unverified ⇒ bit unset — if the primitive can't demonstrate
confinement, the honest declaration is `Scope: ""`, not a redefinition of what the
bit means. Redefining the semantics in IMPLEMENTATION.md after ratification is a
design change that never went through the deck.

**Failing scenario:** A user reads `AUTO=yes` (workspace) for hermes-1 and runs it
unattended on a deck; the agent's tools write to `~/.config/...` or `../other-repo`
with full auto-approval. The displayed guarantee ("workspace") does not hold.

**Fix:** Per the ratified rule: set `Scope: "workspace"` only for codex (verified
sandbox); leave `Scope` empty for claude/agy/hermes (so `Declared()` is false and
AUTO renders "no"/"unverified") until a demonstrable confinement lands, OR take an
explicit amendment through the deck that redefines the `AUTO` column semantics and
rename the surface so it no longer implies confinement. Do not keep the current
silent redefinition.

### MINOR

#### [MINOR] `rosterIDRe` rejects legacy non-canonical custom family IDs (narrow legacy regression)

**File:** `internal/agents/resolve.go:14` (regex), enforced at :31

`ResolveParticipant` now hard-errors on any participant outside
`^[a-z0-9][a-z0-9-]*$` — including in rule 1 (exact spec-ID match). Spec IDs can come
from `[agents.<key>]` config sections, and TOML bare keys allow `_` (and quoted keys
allow `.`, uppercase). A deck that pre-change ran `--participants my_cli` against
`[agents.my_cli]` now fails at selection and at the runner ("invalid participant
id"). The design promised "legacy spec-ID participants keep working via resolver
rule 1". The traversal fix is correct and must stay — the participant becomes a path
segment — but the grammar check conflates "safe to join into a path" with "known to
the local config".

**Fix:** Keep the grammar check for the mapping path (rule 2) and for anything
written into `00-prompt.md`/§2, but for rule-1 exact matches against a discovered
spec ID, accept any ID that is itself a configured spec ID and additionally passes a
containment check (no `/`, `\`, or `..`) at the `filepath.Join` sites. Add a legacy
underscore-ID test.

#### [MINOR] `roster init` can print "already initialized" while the effective mapping is absent

**Files:** `internal/app/roster.go:315` (textual guard), with `internal/config/runtime.go:241` (empty adapters skipped) and `internal/app/roster.go:227` (idempotency skip)

`RosterAdaptersInFile` skips entries whose adapter is empty/whitespace, so
`rosterInit` plans to write them; but `writeRosterMappings` then skips any id whose
literal string `[roster.<id>]` appears *anywhere* in the file — including an
empty-adapter table or a comment.

**Failing scenario:** target file contains
```toml
[roster.claude-1]
adapter = ""
```
(or a comment `# see [roster.claude-1]`). `parley roster init --yes` → the planned
entry is textually "present" → `wrote == 0` → outcome `unchanged` → prints "Roster
already initialized (…): every §2 roster id already maps to a family." The mapping is
effectively absent, every subsequent run fails to resolve `claude-1`, and the command
whose job is to fix exactly this reports success.

**Fix:** Decide presence from the parsed TOML (a table with a non-empty `adapter`),
not substring matching. An empty-adapter table should be repaired in place or
reported as a hard error telling the user to fix the block — never silently skipped.

#### [MINOR] `driver_impl` safety checks go silently vacuous for roster-ID decks

**Files:** `internal/app/driver_impl.go:120-127` (`modelOf`), `internal/app/driver_impl.go:354-361` (`discoveryFor`)

Both look up roster IDs (`o.implementer`, `o.drafter`, `o.reviewers`) in raw
family-ID discoveries.

**Failing scenario:** A roster-ID deck under auto-implement: (a) `GoalCheck` always
takes the `!ok` branch — "goal-check skipped — checker claude-1 not discovered
(advisory)" — so the LE-7 goal-done gate never executes; (b)
`reviewersShareImplementerModel` always returns `("", false)` because
`modelOf("claude-1")` misses, so the LE-3 model-diversity check never fires — **even
with `require_model_diversity: true`, a configured hard gate is silently disabled**
(driver_impl.go:177 unreachable). Both are fail-open with no warning that the lookup
missed.

**Fix:** Resolve ids through the mapping (same resolver) before lookup; if a
configured hard gate's inputs can't be resolved, fail loud rather than vacuous pass.

#### [MINOR] Command/config surface gaps vs the ratified design: no `roster diff`, no `autonomous_write` TOML field

**Files:** `internal/app/roster.go:70-77` (only `show|init` dispatched); `internal/config/runtime.go:97-132` (`agentOverride` has no `autonomous_write`)

FINAL decision 5 and §B specify `parley roster init|show|diff`; `diff` is absent
(FINAL's verification section even requires a `roster init --dry-run` — present —
but `diff` was the cross-layer disagreement surfacing tool; with layered mapping,
`init --scope machine` can silently keep a central `claude-1→X` while the deck maps
`claude-1→Y`, and nothing shows the skew). FINAL S3 lists `autonomous_write` as a
config field; only `model_label` was added, so the consensus §C promise "a vendor
flag change is a config edit, not a skill revision" is not deliverable — the bit is
settable only in Go code.

**Fix:** Add `roster diff` (layered mapping vs target file, plus §2 roster skew) and
an `autonomous_write` block in `agentOverride` mapped in `applyOverride`, or record
both as explicitly deferred deviations in IMPLEMENTATION.md (which currently lists
neither).

### NIT

#### [NIT] Unresolved participants emit no per-agent failure event

**File:** `internal/runner/runner.go:220-227`

The fail-closed pseudo-result flips the round to `round.incomplete` (correct), but no
`agent.failed` store event is appended, so the runstate projection/TUI shows an
incomplete round with no failed-agent badge (CLI `printRunResults` does show the
result). Emit the event (agent id, `segment_id`, class `roster-unresolved`) for
projection parity.

#### [NIT] Stale comment on `selectedAgents`

**File:** `internal/runner/runner.go:356-360`

The doc comment still says "Unresolvable participants are skipped (matching the
pre-split behavior of an absent participant)" — they are now returned in the
`unresolved` out-parameter and turned into failed results by `RunRoundOne`. Update
the comment to describe the new contract (and note that the phase58 single-agent
callers rely on the empty-`selected` check instead).

## Open questions

- The implementer's fix-up notes defer the MAJOR wiring gaps as "a scoped follow-up."
  Given that a roster-ID idea cannot complete `Draft`/`RequestSignoffs` today, should
  the user-facing docs (SKILL.md, `roster init` success output) state that roster-ID
  participants currently work for rounds/review/implementation but not for the
  driver's consensus/signoff automation, until the follow-up lands?

## Overall verdict

**ACCEPT-WITH-FIXES** — the identity/adapter split, resolver, naming, speed axis, and
roster-init writer are correct and well-tested, and legacy decks are unaffected. The
three MAJORs are merge-gating for the feature's stated purpose ("the driver can run a
`claude-1, …` roster"): the roster-ID wiring must reach consensus drafting, signoff
collection, and steer; preflight must not evaluate an empty set as ready; and the
`AUTO` workspace claim must be brought back in line with the ratified fail-closed
honesty rule (or amended openly). None of these requires redesign — one shared
resolution boundary plus the scope/declaration fixes covers all of them.
