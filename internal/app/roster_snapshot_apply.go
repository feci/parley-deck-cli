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
	// Key by ROSTER ID, not by adapter family. Two roster IDs may share an adapter and
	// still run different models — the implementation's own contract says so — and an
	// adapter-keyed map made the second entry overwrite the first, so both continuations
	// launched the last entry's model. The freeze must be per member.
	frozen := make(map[string]runmanifest.RosterSnapshotEntry, len(snapshot))
	byAdapter := make(map[string]runmanifest.RosterSnapshotEntry, len(snapshot))
	for _, e := range snapshot {
		frozen[e.Agent] = e
		// Adapter is only a FALLBACK, for runs frozen before roster IDs were recorded.
		// First writer wins so the fallback cannot silently reorder.
		if _, seen := byAdapter[e.Adapter]; !seen {
			byAdapter[e.Adapter] = e
		}
	}
	out := make([]agents.Discovery, 0, len(discovered))
	for _, d := range discovered {
		e, ok := frozen[d.Spec.ID]
		if !ok {
			e, ok = byAdapter[d.Spec.Adapter()]
		}
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
		// Pin the AUTO posture too. Model/effort/speed alone leave the launch shape free
		// to change under a running idea: dropping an auto-approve flag from the machine
		// config would alter what the continuation is permitted to do while the frozen
		// row still reported AUTO=yes.
		if len(e.LaunchArgs) > 0 {
			d.Spec.HeadlessArgs = append([]string(nil), e.LaunchArgs...)
		}
		out = append(out, d)
	}
	return out
}
