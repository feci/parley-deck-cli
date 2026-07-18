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

## Notes for reviewers

- The identity/adapter split is the crux: `runAgent` keys artifacts on `agent.ID` (now the
  participant/roster ID after resolution); the three vendor branches key on `agent.Adapter()`
  (family). Confirm no other site switches on `agent.ID` for vendor behavior (grep clean).
- `ResolveParticipant` never guesses (no prefix heuristic); `proposeFamily` (init-time only)
  may propose via `-N` strip / alias, but the runtime resolver requires an explicit mapping.
- Speed never touches model/effort in Go (guard test); the legacy `profiles.fast` downgrade
  lived only in deck JSON/skill convention.
