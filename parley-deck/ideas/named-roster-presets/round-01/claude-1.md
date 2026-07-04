---
agent: claude-1
idea: named-roster-presets
round: 1
date: 2026-07-04
---

## Summary

Presets should be pure expansion sugar over the existing config layering: a
`[rosters.<name>]` table resolved through the same precedence chain as agent specs
(built-ins → `~/.parley/agents.toml` → deck config), expanded ONCE at idea creation
into the `participants:` list, which stays the sole source of truth. No new runtime
concept — after expansion the system cannot tell a preset was used (except for an
audit note).

## Proposed approach

1. **Config shape** — extend the TOML schema read by `internal/config` (the
   `LoadAgentSpecs` layer already merges central + deck files):
   ```toml
   [rosters.council]
   participants = ["claude", "codex", "hermes", "agy"]
   [rosters.pair]
   participants = ["claude", "codex"]
   [rosters.default_by_track]   # optional
   fast = "pair"
   standard = "trio"
   deliberation = "council"
   ```
   Deck file overrides central per-key (a deck redefining `rosters.pair` wins; a deck
   adding `rosters.release` extends).

2. **Resolution function** — pure Go, unit-testable:
   `ResolveRoster(name string, track track.Track, cfg) ([]string, string, error)`
   - explicit `--roster <name>` wins; else if idea has explicit participants, use them;
     else `default_by_track[track]` if configured; else no preset (today's behavior).
   - Returns the expanded participant list + a provenance string ("roster preset
     'council' from ~/.parley/agents.toml") for the audit trail.
   - Errors loudly on: unknown preset, empty preset, participant not in §2 roster /
     not discoverable. Never silently drops an agent.

3. **Expansion point** — idea scaffolding (`parley idea new` or the equivalent
   kickoff path) writes the EXPANDED list into `00-prompt.md` `participants:` and adds
   an HTML comment `<!-- roster: council (preset) -->` for provenance. Nothing
   downstream changes: driver, quorum, track policy all read `participants:` as today.

4. **Track-policy interaction** — expansion happens BEFORE track policy; §4.0 minimums
   then validate the result (e.g. `fast` + preset of 1 participant = the existing
   non-solo hard error). A preset larger than the track's MaxReviewers is fine — the
   driver already truncates reviewers by track.

5. **Protocol text** — none needed. `participants:` semantics are unchanged;
   presets are machine-local convenience (like headless-agents.local.json). At most a
   one-line mention in the Quickstart, and I would even skip that in v1.

## Concerns / open questions

1. Roster IDs vs agent keys: presets use central config keys (`claude`, `codex`);
   deck roster IDs are `claude-1`, `codex-1`. Need one mapping rule — propose: preset
   holds config keys; expansion maps key → deck roster ID via the §2 roster table and
   fails if ambiguous (two roster IDs for one key).
2. Where does `parley idea new` live today? If there is no scaffolding subcommand yet,
   scope creep risk — v1 could be `parley roster expand <preset>` printing the list,
   plus TUI/facilitator usage, with scaffolding integration when it exists.
3. `default_by_track` interacts with the §4.0 classifier: if the idea author declared
   no track and the classifier picks one later, the preset default must use the
   DECLARED track only — never a guessed one. Absent both ⇒ no preset.

## Risks

- Silent quorum shrink if expansion tolerated unknown agents — mitigated by hard error.
- Config drift between central and deck presets confusing users — mitigated by
  provenance string recorded at expansion.
- Over-generalizing into per-preset model overrides (explicit non-goal).
