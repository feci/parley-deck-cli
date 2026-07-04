package tui

import (
	"encoding/json"
	"fmt"
	"strings"
)

// roundsummary.go renders the consolidated round digest in the Home tab
// (tui-round-summary). The digest is produced by the driver and carried on
// `round.digest` events; this is a pure renderer over already-consumed events —
// no disk scan, no markdown parsing on the render path.

// digestView is a local mirror of driver.RoundDigest for decoding the event blob
// (kept local to avoid a tui→driver import edge).
type digestView struct {
	Idea      string `json:"idea"`
	Round     int    `json:"round"`
	Total     int    `json:"total"`
	Completed int    `json:"completed"`
	Lines     []struct {
		Agent    string `json:"agent"`
		Position string `json:"position"`
		Fell     bool   `json:"fell_back"`
		Present  bool   `json:"present"`
	} `json:"lines"`
	FlagBlock   int    `json:"flag_block"`
	FlagCounter int    `json:"flag_counter"`
	FlagAccept  int    `json:"flag_accept"`
	FlagEscal   int    `json:"flag_escalate"`
	Next        string `json:"next"`
}

// latestRoundDigest returns the most recent round.digest event decoded, or ok=false.
func (m liveModel) latestRoundDigest() (digestView, bool) {
	for i := len(m.events) - 1; i >= 0; i-- {
		e := m.events[i]
		if e.Type != "round.digest" {
			continue
		}
		blob, _ := e.Data["digest"].(string)
		if blob == "" {
			continue
		}
		var dv digestView
		if err := json.Unmarshal([]byte(blob), &dv); err != nil {
			continue
		}
		return dv, true
	}
	return digestView{}, false
}

// renderRoundDigest formats the latest digest as a bounded block for the Home tab.
// It caps itself to maxRows lines so it can never push the chips/roster/runs list
// off-screen (the explicit regression guard); overflow is truncated with a note.
func renderRoundDigest(dv digestView, width, maxRows int) string {
	if maxRows < 3 {
		return ""
	}
	header := sectionTitle(fmt.Sprintf("Round %02d digest — complete (%d/%d)", dv.Round, dv.Completed, dv.Total))

	// Trailer: flags (HINTS, never verdicts) + the next-action line.
	trailer := []string{
		mutedStyle.Render(truncateText(fmt.Sprintf("  mentions: %d block · %d counter-proposal · %d accept · %d escalate",
			dv.FlagBlock, dv.FlagCounter, dv.FlagAccept, dv.FlagEscal), width-1)),
	}
	if dv.Next != "" {
		trailer = append(trailer, mutedStyle.Render(truncateText("  next: "+dv.Next, width-1)))
	}

	// Agent rows fill whatever remains after header + trailer (hard cap on total).
	agentBudget := maxRows - 1 - len(trailer)
	var agentRows []string
	for _, ln := range dv.Lines {
		if len(agentRows) >= agentBudget {
			if agentBudget > 0 { // replace the last row with an overflow note
				agentRows[agentBudget-1] = mutedStyle.Render(fmt.Sprintf("  … %d more (open the agent tabs)", len(dv.Lines)-agentBudget+1))
			}
			break
		}
		switch {
		case !ln.Present:
			agentRows = append(agentRows, warnStyle.Render(fmt.Sprintf("  @%-13s [no artifact]", ln.Agent)))
		case ln.Fell:
			agentRows = append(agentRows, fmt.Sprintf("  @%-13s %s", ln.Agent, truncateText(ln.Position+"  [no Summary — fell back]", width-18)))
		default:
			agentRows = append(agentRows, fmt.Sprintf("  @%-13s %s", ln.Agent, truncateText(ln.Position, width-18)))
		}
	}

	out := append([]string{header}, agentRows...)
	out = append(out, trailer...)
	return strings.Join(out, "\n")
}
