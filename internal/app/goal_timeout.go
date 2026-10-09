package app

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/track"
)

// Unlike trackTimeoutCeiling's UI fallback, a declared-but-malformed track is
// not permission for the longest goal check. Only an absent field defaults.
func goalCheckTimeout(ideaDir string, configuredMS int) (time.Duration, error) {
	meta, err := protocol.ReadFrontmatter(filepath.Join(ideaDir, "00-prompt.md"))
	if err != nil {
		return 0, err
	}
	t := track.Standard
	if raw, present := meta["track"]; present {
		value := strings.ToLower(strings.Trim(strings.TrimSpace(raw), "\"'"))
		switch track.Track(value) {
		case track.Fast, track.Standard, track.Deliberation:
			t = track.Track(value)
		default:
			return 0, fmt.Errorf("invalid goal-check track %q", raw)
		}
	}
	ceiling := trackTimeoutCeiling(string(t))
	// Compare milliseconds before conversion so a huge configured integer cannot
	// overflow time.Duration and turn a configured upper bound into a default.
	if configuredMS > 0 && int64(configuredMS) < ceiling.Milliseconds() {
		ceiling = time.Duration(configuredMS) * time.Millisecond
	}
	return ceiling, nil
}
