package app

import (
	"fmt"
	"io"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/runmanifest"
)

// applyRosterSnapshot pins freshly discovered agents back to what the run FROZE at
// creation.
//
// Without this, `continue` re-discovers configuration on every resumption, so changing a
// machine default — or running `roster sync`, which is designed to change deck values —
// could silently move a running idea onto a different model mid-deliberation. The run
// record would then describe a deliberation that never happened that way.
//
// The snapshot wins for the fields it froze. An agent absent from the snapshot is left
// as discovered (it joined later, or the run predates snapshots) and is reported, because
// silently launching an unfrozen agent inside a frozen run is the same class of surprise.
func applyRosterSnapshot(discovered []agents.Discovery, snapshot []runmanifest.RosterSnapshotEntry, warn io.Writer) []agents.Discovery {
	if len(snapshot) == 0 {
		return discovered
	}
	frozen := make(map[string]runmanifest.RosterSnapshotEntry, len(snapshot))
	for _, e := range snapshot {
		frozen[e.Adapter] = e
	}
	out := make([]agents.Discovery, 0, len(discovered))
	for _, d := range discovered {
		e, ok := frozen[d.Spec.Adapter()]
		if !ok {
			if warn != nil {
				fmt.Fprintf(warn, "warning: %s is not in this run's roster snapshot — launching it as currently configured\n", d.Spec.ID)
			}
			out = append(out, d)
			continue
		}
		if e.Model != "" && e.Model != agents.Unknown {
			d.Spec.Model = e.Model
		}
		if e.Effort != "" && e.Effort != agents.Unknown {
			d.Spec.Reasoning = e.Effort
		}
		if e.Speed != "" {
			d.Spec.Speed = e.Speed
		}
		out = append(out, d)
	}
	return out
}
