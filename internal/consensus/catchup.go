package consensus

import (
	"strings"

	"parley-deck-cli/internal/protocol"
)

// This is the pre-existing CLI representation of a §5 NON-PARTICIPANT note,
// not a new signoff status or authority to change membership. Exact notes keep
// examples, quotations and ordinary BLOCKs out of the catch-up exception.
func catchupDecline(s Signoff) bool {
	status, err := CanonicalStatus(s.Status)
	return err == nil && status == StatusBlock &&
		strings.TrimSpace(s.Notes) == "❌ NON-PARTICIPANT" &&
		strings.TrimSpace(s.CounterProposal) != ""
}

func quotaCatchupDecliners(dir string, review bool) ([]string, string, error) {
	v, err := protocol.InspectQuota(dir)
	if err != nil {
		return nil, "", err
	}
	if review || v.History == nil || v.History.Policy().Enabled {
		return nil, v.Pending, nil
	}
	return v.Catchup, v.Pending, nil
}
