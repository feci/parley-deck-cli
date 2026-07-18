---
idea: composite-agent-naming-and-roster-reinit
status: implemented
implementer: claude-1
started: 2026-07-18
branch: parley-deck-cli#composite-agent-naming-and-roster-reinit
---

## Summary of work

Implemented the ratified design in staged, independently-green commits on branch
`composite-agent-naming-and-roster-reinit`. Roster ID stays the identity; the composite
`family_model_effort` name is derived for display; the two-namespace schism is fixed with a
fail-closed resolver; autonomous write is first-class with corrected built-ins; `fast` is the
standard speed on a separate axis. `go build ./... && go vet ./... && go test ./...` green
(25 packages, 0 fail) after every stage.

## Implementation plan / checklist

- [x] **S1 — naming core** `internal/agents/naming.go` (+test): `SanitizeSection`, `Compose`,
      `Parse` (fail-closed, right-to-left), `NormalizeEffortToken`/`EffortDisplayForm`
      (camelCase `xHigh`/`cliDefault`), `StripParenTier` (agy), `RenderDisplayName`.
      Grammar `family_model_effort` (`_` meaning, `-` word, `.` version) per the user's
      post-signoff decision (FINAL amendment 1).
- [x] **S2 — Spec split + resolver** `Spec.AdapterID`/`Adapter()`, `internal/agents/resolve.go`
      `ResolveParticipant` (exact spec-ID → explicit `[roster.*]` mapping → hard error, no
      prefix heuristic); `runner.selectedAgents` resolves, `runAgent` artifact path + vendor
      dispatch (env clean / hermes env / isolated home) split identity vs `Adapter()`.
      Non-breaking (nil mapping = exact-ID; legacy decks unchanged).
- [x] **S3 — config + speed** `Spec.ModelLabel` + config `model_label`; central template
      `speed = "fast"`; speed-invariance guard test (LoadAgentSpecs keeps model+reasoning
      under `fast`). codex model corrected to `gpt-5.6-sol` in `~/.parley/agents.toml`
      (FINAL amendment 2).
- [x] **S4 — autonomous built-ins** `AutonomousWrite{Mode,Args,Scope}` + `Declared()`; AUTO
      column in the runtime matrix; posture fixes: claude `acceptEdits→bypassPermissions`,
      codex `on-failure→never`, hermes `+--yolo` (agy already autonomous); claude
      `ModelLabel="Opus 4.8 1m"`. CHANGELOG note.
- [x] **S5 — `parley roster` command** `internal/app/roster.go`: `roster show` (resolved
      roster + composite names), `roster init` (proposes+writes `[roster.*]` mapping,
      fail-closed on unresolved, idempotent, `--dry-run/--yes/--json`);
      `config.LoadRosterAdapters` + a runner `RosterMappingLoader` inject the mapping into
      every run path so the driver can run a `claude-1, …` roster.
- [x] **S7 — skill** `parley-deck-skill/SKILL.md`: an **Autonomous Execution (required)**
      section (the yolo requirement + per-CLI mapping) and an **Agent display names & roster
      init** section (composite names + `parley roster show|init` + fast-axis). Protocol
      wording unchanged, so both COOPERATION.md copies + the drift guard are untouched.
- [x] Checks: `go build ./...`, `go vet ./...`, `gofmt -l` (my files clean), `go test ./...`
      green.

## Deviations from FINAL.md

- **S6 (§2 `Display name` column + digest/TUI display rendering) DEFERRED** as a follow-up.
  Rationale: the composite names are already user-visible via `parley roster show` (S5); the
  §2-table column rewrite carries real drift-guard risk (`TestEmbeddedDefaultMatchesLiveDeck`)
  for low incremental benefit, and digest/TUI rendering needs invasive driver→config plumbing.
  Roster IDs remain the on-screen identity during runs (correct); the composite is a lookup
  via `roster show`. Tracked as deferred; not required for the feature to work.
- **`roster init` is proposal-based, not fully interactive.** It records the roster-ID→family
  mapping and is idempotent; per-agent interactive model/effort re-selection (vs the current
  config values) is a lighter follow-up — the config write path and `model_label`/`speed`
  fields already support it.
- FINAL amendments 1 (underscore grammar) and 2 (codex `gpt-5.6-sol`) are post-signoff user
  decisions, recorded in FINAL.md.

## Fix-up cycle 1 (review round-01)

codex-1 ran a refutation-default review (verdict BLOCK) and found real defects. Fixed
(commit 54b5282, full suite green):

- **[CRITICAL] uncommitted compile blocker** — `naming.go`/`naming_test.go` (with
  `RenderDisplayName`) were left out of the S5 commit, so a clean checkout failed
  `undefined: agents.RenderDisplayName`. Committed.
- **[CRITICAL] path traversal** — `ResolveParticipant` now validates the participant against
  `^[a-z0-9][a-z0-9-]*$` before it reaches `filepath.Join`, so a malicious
  `[roster."../../x"]` can't make a CLI write outside the deck.
- **[CRITICAL] fail-open resolution** — `RunRoundOne` emits a failed result for every
  unresolved participant, so a round is `round.incomplete`, never silently completed with a
  partial/empty quorum.
- **[MAJOR] run-selection wiring** — `selectedParticipantIDs` resolves a roster id via the
  `[roster.*]` mapping, so `parley run --participants claude-1` works (was rejected as "not
  installed").
- **[MAJOR] roster init** — idempotency judged against the TARGET file (not the layered
  stack); `--scope` validated; `--json` actually writes and reports the real outcome
  (`written`/`unchanged`/`dry-run`/`needs-confirmation`); an existing-but-invalid mapping
  fails closed; the write is atomic (temp+rename) and skips already-present blocks.
- **[MINOR] Parse canonical round-trip** (rejects `x-high`, lowercase `xhigh`, `_02`);
  exact-ID resolution preserves an already-explicit adapter.

### Deferred from review (documented, not fixed this cycle)

- **[MAJOR] deeper app-level roster-ID wiring** — beyond run selection, a few paths still
  compare participant strings to raw family discovery ids: `preflight` readiness ping,
  `driver_consensus` drafter attribution, `consensus request-signoffs`, and TUI `steer`.
  They work for family-id rosters and for the driver's round runner (which resolves), but
  full roster-id support there is a scoped follow-up (one shared resolution boundary). The
  primary `parley run` + runner path is wired.
- **[CRITICAL, contested] autonomous-write confinement honesty** — `AutonomousWrite.Scope="workspace"`
  is asserted from the vendor's own scoping flags (`--add-dir {root}`, `--sandbox
  workspace-write`, cwd), NOT an OS sandbox. Codex's `workspace-write` is a real sandbox;
  claude/agy/hermes rely on flag+cwd confinement. This is the intended, documented mechanism
  (the CLIs offer no stronger primitive); a true OS-enforced sandbox + live sentinel probe is
  a follow-up. `AUTO=yes` means "declared + flag-scoped", not "OS-jailed".

## Notes for reviewers

- The identity/adapter split is the crux: `runAgent` keys artifacts on `agent.ID` (now the
  participant/roster ID after resolution); the three vendor branches key on `agent.Adapter()`
  (family). Confirm no other site switches on `agent.ID` for vendor behavior (grep clean).
- `ResolveParticipant` never guesses (no prefix heuristic); `proposeFamily` (init-time only)
  may propose via `-N` strip / alias, but the runtime resolver requires an explicit mapping.
- Speed never touches model/effort in Go (guard test); the legacy `profiles.fast` downgrade
  lived only in deck JSON/skill convention.
