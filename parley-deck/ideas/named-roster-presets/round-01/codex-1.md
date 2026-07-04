---
agent: codex-1
idea: named-roster-presets
round: 1
date: 2026-07-04
---

## Summary

Named roster presets fit the current runtime config model if they are treated as another layered config value, not as protocol state. The preset should be resolved before `runcontrol.Create` writes `00-prompt.md`; after that point only the expanded `participants:` list matters for quorum.

The implementation should mirror `internal/config/runtime.go`: start from built-in defaults, apply `~/.parley/agents.toml`, then apply deck config so the deck can override a central preset. A preset that cannot expand into valid, available, canonical participant IDs should fail before any idea directory is created.

## Proposed approach

Add rosters to the same TOML parse path as agent specs:

```go
type fileConfig struct {
    Defaults *globalDefaults          `toml:"defaults"`
    Agents   map[string]agentOverride `toml:"agents"`
    Rosters  map[string]rosterOverride `toml:"rosters"`
}

type rosterOverride struct {
    Participants []string `toml:"participants"`
}

type globalDefaults struct {
    // existing fields...
    TrackRosters *trackRosterDefaults `toml:"track_rosters"`
}

type trackRosterDefaults struct {
    Fast         string `toml:"fast"`
    Standard     string `toml:"standard"`
    Deliberation string `toml:"deliberation"`
}
```

The TOML shape should be:

```toml
[rosters.council]
participants = ["claude-1", "codex-1", "hermes-1", "antigravity-1"]

[rosters.pair]
participants = ["claude-1", "codex-1"]

[defaults.track_rosters]
fast = "pair"
standard = "pair"
deliberation = "council"
```

Preset participants should be canonical Parley participant IDs: the same strings used in `COOPERATION.md` §2 and in `00-prompt.md participants:`. They should not be CLI family aliases unless the deck has explicitly made those IDs canonical. In this repo, that means `codex-1` rather than `codex` unless the runtime/roster ID mapping is deliberately changed.

Add a config loader beside `LoadAgentSpecs`, for example `LoadRosterPresets(root) (RosterConfig, error)`, that uses the existing `configLayers(root)` order. The merge rule should be atomic by preset name: if a higher layer defines `[rosters.council]`, its `participants` list replaces the lower-layer list instead of appending. Track defaults merge per field, so a deck can override only `fast` while inheriting central `standard` and `deliberation`.

Built-ins can be an empty roster map in v1. If built-in presets are added later, they should be just the lowest layer in the same chain:

1. built-in roster presets
2. `~/.parley/agents.toml`
3. `parley-deck/agents.toml`
4. any existing higher local/env config layer, if the runtime continues to apply those above deck config

The expansion point should be in the idea-creation CLI path, before preflight and before `runcontrol.Create`. Today `runTask` discovers runtime specs through `discoverConfigured`, then calls `selectedParticipantIDs(discovered, *participantsFlag)`, then preflight, then `runcontrol.Create`, which writes `participants: [...]`. Add `--roster NAME` there and to any future `parley idea new` entry point:

1. reject `--participants` plus `--roster` as ambiguous;
2. if `--roster NAME` is set, load layered roster presets and expand `NAME`;
3. if neither is set and an idea track is already known, look up `[defaults.track_rosters].<track>` and expand it when present;
4. otherwise keep the current default of installed configured agents;
5. run the existing availability/preflight logic on the expanded concrete list;
6. pass only the concrete participants to `runcontrol.Create`.

This keeps `protocol.CreateIdeaWithExclusions` simple. It should not know which preset was used, and it should not write `roster:` into frontmatter. `00-prompt.md participants:` remains the canonical quorum.

Validation should happen at expansion time and fail closed:

- unknown preset: return a hard error such as `unknown roster preset "pair"`;
- empty preset or duplicate participant IDs: hard error, because the author likely made a quorum mistake;
- preset references an ID not present in `COOPERATION.md` §2 active roster: hard error before discovery/preflight;
- preset references a rostered ID that is not discovered/installed/available: hard error or preflight hard stop, with no silent shrinking;
- track default references an unknown preset: hard error only when that track default is selected;
- explicit `participants:` in an existing `00-prompt.md` is never re-expanded from config, even if the central or deck preset changes later.

To validate against §2, add a small protocol helper such as `protocol.ReadRosterIDs(root)` that extracts the active roster IDs from `COOPERATION.md`. If that parser cannot confidently read the roster table, roster-preset expansion should fail instead of falling back to installed agent IDs.

Tests should cover the layering and the CLI boundary:

- central preset inherited when deck has none;
- deck preset with the same name replaces central participants;
- track defaults merge per track;
- unknown preset fails before idea creation;
- preset member not in §2 fails;
- preset member in §2 but not discovered fails without shrinking;
- `--participants` and `--roster` together fail;
- the generated `00-prompt.md` contains only expanded `participants: [...]`.

## Concerns / open questions

The current runtime config keys and the protocol roster IDs may not always be the same. `LoadAgentSpecs` discovers `agents.Spec.ID` values, while this idea's active participant IDs include suffixes such as `codex-1`. If Parley wants presets to reference §2 IDs, the discovery layer must be able to select those IDs exactly. If Parley wants presets to reference runtime family IDs such as `codex`, then the §2 roster validation requirement becomes weaker. I recommend exact canonical IDs for v1 and a separate explicit alias/mapping design if needed.

Track-linked defaults need a known track before expansion. For existing `parley run`, the current creation path does not appear to accept a `--track` flag or write `track:` in `CreateIdeaWithExclusions`. Either add track selection to the creation command first, or limit automatic track-default expansion to entry points that already know the track. Do not infer a track after the idea has been written and then mutate participants.

The environment override layer exists in `configLayers`. The idea brief only names central and deck config, but the implementation should decide whether `PARLEY_HEADLESS_AGENT_CONFIG` and `agents.local.toml` are allowed to override rosters too. My preference is yes, because using the same layer iterator avoids surprising precedence splits.

## Risks

The main risk is silent quorum drift. If a preset is expanded lazily after `00-prompt.md` exists, or if unavailable agents are dropped automatically, later consensus and review gates no longer match the author's intended panel. Expansion must be one-time, before idea creation, and failures must stop creation.

The second risk is ID ambiguity. Accepting both `codex` and `codex-1` as casual synonyms would make artifacts, signoffs, and runtime launches disagree. Either require one canonical ID everywhere or introduce an explicit, test-covered mapping.

The third risk is overloading presets with behavior. V1 should not include per-preset model, reasoning, transport, or track-policy overrides. Presets should choose participants only; all agent runtime fields continue to come from `[agents.*]`, and track policy continues to apply after expansion.
