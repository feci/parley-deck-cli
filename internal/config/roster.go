package config

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// roster.go implements named roster presets (named-roster-presets): expand-at-creation
// sugar over the layered config. `participants:` in 00-prompt.md stays canonical; a
// preset just supplies that list. Members are §2 canonical roster IDs.

// RosterPreset is one named participant preset with its winning source layer.
type RosterPreset struct {
	Name         string
	Participants []string
	Source       string // human-readable winning layer (e.g. "parley-deck/agents.toml")
}

// RosterConfig is the merged view of all [rosters.*] presets and the optional
// [defaults.track_rosters] track→preset default map, across the layered config files.
type RosterConfig struct {
	Presets      map[string]RosterPreset
	TrackDefault map[string]string // track name → preset name
}

// LoadRosterPresets merges [rosters.*] and [defaults.track_rosters] across the same
// layered config chain as LoadAgentSpecs/LoadDefaults (central → deck → env). Merge is
// per-preset-name (a higher layer replaces a preset's participant list, never appends)
// and per-track-key for the track defaults.
func LoadRosterPresets(root string) (RosterConfig, error) {
	out := RosterConfig{Presets: map[string]RosterPreset{}, TrackDefault: map[string]string{}}
	for _, item := range configLayers(root) {
		data, err := os.ReadFile(item.path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return out, err
		}
		var cfg fileConfig
		if err := toml.Unmarshal(data, &cfg); err != nil {
			return out, fmt.Errorf("%s: %w", item.path, err)
		}
		for name, ro := range cfg.Rosters {
			// Per-name replace: the winning (higher) layer owns the whole list.
			out.Presets[name] = RosterPreset{
				Name:         name,
				Participants: append([]string(nil), ro.Participants...),
				Source:       item.source,
			}
		}
		if cfg.Defaults != nil {
			for track, preset := range cfg.Defaults.TrackRosters {
				out.TrackDefault[strings.ToLower(strings.TrimSpace(track))] = strings.TrimSpace(preset)
			}
		}
	}
	return out, nil
}

// RosterResolution is the outcome of expanding a preset selection.
type RosterResolution struct {
	Participants []string
	Preset       string // preset name used, or "" when none applied
	Source       string // winning layer for that preset, or ""
	Provenance   string // one-line HTML-comment-ready provenance, or ""
}

// ResolveRoster expands a preset selection into concrete §2 participant IDs.
//
//   - explicitPreset != "": that preset MUST exist (hard error otherwise).
//   - explicitPreset == "" && track != "": use track default when one is configured;
//     a configured-but-unknown track default is a hard error.
//   - otherwise: no preset applies (RosterResolution.Participants is nil) — the caller
//     keeps today's behavior.
//
// rosterIDs is the §2 active roster (from protocol.ReadRosterIDs); inactive maps IDs
// that appear in §2 but are marked inactive. Both drive fail-closed validation.
func ResolveRoster(cfg RosterConfig, explicitPreset, track string, rosterIDs, inactive map[string]bool) (RosterResolution, error) {
	name := strings.TrimSpace(explicitPreset)
	if name == "" {
		if t := strings.ToLower(strings.TrimSpace(track)); t != "" {
			if d, ok := cfg.TrackDefault[t]; ok && d != "" {
				name = d
			}
		}
	}
	if name == "" {
		return RosterResolution{}, nil // no preset applies
	}

	preset, ok := cfg.Presets[name]
	if !ok {
		return RosterResolution{}, fmt.Errorf("unknown roster preset %q (known: %s)", name, knownPresetNames(cfg))
	}
	if len(preset.Participants) == 0 {
		return RosterResolution{}, fmt.Errorf("roster preset %q is empty", name)
	}
	seen := map[string]bool{}
	for _, id := range preset.Participants {
		if seen[id] {
			return RosterResolution{}, fmt.Errorf("roster preset %q lists %q twice", name, id)
		}
		seen[id] = true
		if len(rosterIDs) > 0 && !rosterIDs[id] {
			return RosterResolution{}, fmt.Errorf("roster preset %q references %q, which is not in the §2 active roster", name, id)
		}
		if inactive[id] {
			return RosterResolution{}, fmt.Errorf("roster preset %q references %q, which is marked inactive in the §2 roster", name, id)
		}
	}
	return RosterResolution{
		Participants: append([]string(nil), preset.Participants...),
		Preset:       name,
		Source:       preset.Source,
		Provenance:   fmt.Sprintf("<!-- roster-preset: %s (source: %s) -->", name, preset.Source),
	}, nil
}

func knownPresetNames(cfg RosterConfig) string {
	if len(cfg.Presets) == 0 {
		return "none defined"
	}
	names := make([]string, 0, len(cfg.Presets))
	for n := range cfg.Presets {
		names = append(names, n)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}
