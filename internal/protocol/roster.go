package protocol

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// roster.go reads the §2 "Active agents (roster)" table from COOPERATION.md so that
// named roster presets (named-roster-presets) can validate their members against the
// canonical roster and fail closed on unknown/inactive IDs.

var (
	// A roster row: first cell is a backtick-wrapped Agent ID.
	rosterRowRe = regexp.MustCompile("^\\|\\s*`([a-z0-9][a-z0-9-]*)`\\s*\\|")
	// The §2 roster table header (distinguishes it from the host-handle table).
	rosterHeaderRe = regexp.MustCompile(`(?i)^\|\s*Agent ID\s*\|\s*Workspace`)
)

// ReadRosterIDs parses the §2 roster table in the deck's COOPERATION.md and returns
// the set of active Agent IDs plus the subset marked inactive. `ok` is false when the
// roster table cannot be found/parsed, so callers can fail closed rather than fall
// back to installed agents.
func ReadRosterIDs(root string) (active map[string]bool, inactive map[string]bool, ok bool) {
	path := filepath.Join(root, DeckDir, "COOPERATION.md")
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, false
	}
	defer f.Close()

	active = map[string]bool{}
	inactive = map[string]bool{}
	inTable := false
	sawHeader := false
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if rosterHeaderRe.MatchString(line) {
			inTable = true
			sawHeader = true
			continue
		}
		if !inTable {
			continue
		}
		if !strings.HasPrefix(line, "|") {
			// Table ended (blank line or prose) — stop at the first roster table.
			break
		}
		if strings.HasPrefix(line, "| -") || strings.HasPrefix(line, "|-") {
			continue // separator row
		}
		m := rosterRowRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		id := m[1]
		active[id] = true
		if strings.Contains(strings.ToLower(line), "inactive") {
			inactive[id] = true
		}
	}
	if err := sc.Err(); err != nil {
		return nil, nil, false
	}
	return active, inactive, sawHeader && len(active) > 0
}
