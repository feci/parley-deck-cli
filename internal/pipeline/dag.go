package pipeline

// ReadyBlocks returns the IDs of blocks that may run next: not yet completed,
// with every inbound edge originating from a completed block AND that edge's
// gate approved. A block with no inbound edges (a root) is ready immediately.
// gateApproved is consulted per edge id (EdgeID(from,to)); pass a function that
// returns true when no gate is required.
//
// For a linear manifest this yields the single next incomplete block, so the
// same engine serves both linear and DAG execution.
func (m Manifest) ReadyBlocks(completed []string, gateApproved func(edgeID string) bool) []string {
	done := make(map[string]bool, len(completed))
	for _, c := range completed {
		done[c] = true
	}
	incoming := make(map[string][]Edge, len(m.Blocks))
	for _, e := range m.Edges {
		incoming[e.To] = append(incoming[e.To], e)
	}
	var ready []string
	for _, b := range m.Blocks {
		if done[b.ID] {
			continue
		}
		blocked := false
		for _, e := range incoming[b.ID] {
			if !done[e.From] {
				blocked = true
				break
			}
			if gateApproved != nil && !gateApproved(EdgeID(e.From, e.To)) {
				blocked = true
				break
			}
		}
		if !blocked {
			ready = append(ready, b.ID)
		}
	}
	return ready
}

// AllBlocksComplete reports whether every block id is in completed.
func (m Manifest) AllBlocksComplete(completed []string) bool {
	done := make(map[string]bool, len(completed))
	for _, c := range completed {
		done[c] = true
	}
	for _, b := range m.Blocks {
		if !done[b.ID] {
			return false
		}
	}
	return true
}
