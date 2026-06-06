package runner

import (
	"time"

	"parley-deck-cli/internal/procctl"
	"parley-deck-cli/internal/store"
)

// KillAgentDurable kills (or clears) an agent using only the run's event log —
// the cross-restart / observational path. It reconstructs the agent's durable
// process identity from the latest agent.started event and:
//   - if a terminal event already follows → no-op ("already terminated");
//   - if the recorded process is the live, attributed process → group-kills the
//     whole tree and records agent.killed;
//   - if the recorded process is gone (or the live PID is NOT attributable, i.e.
//     reused by something else) → clears the stale "running" badge with a
//     synthetic agent.failed, and NEVER signals the unattributable PID.
func KillAgentDurable(st store.Store, agentID string) KillResult {
	sp, seg, found, terminal := latestAgentProc(st, agentID)
	if !found {
		// No process record (old run / never started) — clear the badge so it
		// doesn't hang forever.
		return clearStale(st, agentID, seg, "no recorded process; cleared stale status")
	}
	if terminal {
		return KillResult{AgentID: agentID, SegmentID: seg, Message: agentID + " already finished"}
	}
	if !procctl.Alive(sp) {
		// Provably dead (the common 2-day-stale case): clear the hanging badge.
		return clearStale(st, agentID, seg, "process already exited; cleared stale status")
	}
	if ok, reason := procctl.Attributed(sp); !ok {
		// Alive but unattributable: do NOT signal it (it may be a reused PID) and
		// do NOT auto-clear (it might be our own live process the probe couldn't
		// verify). Surface the refusal so the user can decide.
		return KillResult{AgentID: agentID, Failed: true, SegmentID: seg, Message: "process verification failed (" + reason + "); not killed"}
	}
	if err := procctl.KillGroup(sp); err != nil {
		return KillResult{AgentID: agentID, Failed: true, SegmentID: seg, Message: "kill failed: " + err.Error()}
	}
	_ = st.Append(store.Event{
		Time: time.Now().UTC(),
		Type: "agent.killed",
		Data: map[string]any{"agent": agentID, "segment_id": seg, "source": "reattach", "pid": sp.PID, "pgid": sp.PGID},
	})
	return KillResult{AgentID: agentID, Killed: true, SegmentID: seg, Message: "killed " + agentID + " and its process tree"}
}

// DurableKillAt is KillAgentDurable against a run dir (the reattach seam source).
func DurableKillAt(runDir, agentID string) KillResult {
	return KillAgentDurable(store.New(runDir), agentID)
}

// AgentLiveness reports a projected-running agent's true process state for the
// TUI badge: "live" (attributed process alive), "stale" (no attributable live
// process), or "" (not running per the event log).
func AgentLiveness(st store.Store, agentID string) string {
	sp, _, found, terminal := latestAgentProc(st, agentID)
	if !found || terminal {
		return ""
	}
	if !procctl.Alive(sp) {
		return "stale"
	}
	if ok, _ := procctl.Attributed(sp); ok {
		return "live"
	}
	return "stale" // alive but not attributable → our process is gone
}

// AgentLivenessAt is AgentLiveness against a run dir.
func AgentLivenessAt(runDir, agentID string) string {
	return AgentLiveness(store.New(runDir), agentID)
}

// clearStale records a synthetic terminal event so a stale "running" badge that
// has no terminal event stops hanging the run; no OS process is signaled.
func clearStale(st store.Store, agentID, seg, msg string) KillResult {
	_ = st.Append(store.Event{
		Time: time.Now().UTC(),
		Type: "agent.failed",
		Data: map[string]any{"agent": agentID, "segment_id": seg, "error": "stale process cleared by user", "stale_cleared": true},
	})
	return KillResult{AgentID: agentID, Cleared: true, SegmentID: seg, Message: msg}
}

// latestAgentProc returns the durable identity from the agent's most recent
// agent.started, its segment, whether one was found, and whether a terminal
// event (finished/failed/killed) followed it (idempotency / two-parley guard).
func latestAgentProc(st store.Store, agentID string) (sp procctl.Spawned, seg string, found, terminal bool) {
	events, err := st.Load()
	if err != nil {
		return procctl.Spawned{}, "", false, false
	}
	startedIdx := -1
	for i, e := range events {
		if dataString(e.Data, "agent") != agentID {
			continue
		}
		if e.Type == "agent.started" {
			startedIdx = i
			seg = dataString(e.Data, "segment_id")
			sp = procctl.Spawned{
				PID:       dataInt(e.Data, "pid"),
				PGID:      dataInt(e.Data, "pgid"),
				BootID:    dataString(e.Data, "boot_id"),
				ProcStart: dataString(e.Data, "proc_start"),
				Command:   dataString(e.Data, "command"),
				Marker:    dataString(e.Data, "proc_marker"),
			}
			found = true
			terminal = false
		}
	}
	if startedIdx >= 0 {
		for _, e := range events[startedIdx+1:] {
			if dataString(e.Data, "agent") != agentID {
				continue
			}
			switch e.Type {
			case "agent.finished", "agent.failed", "agent.killed", "agent.skipped":
				terminal = true
			}
		}
	}
	return sp, seg, found, terminal
}

func dataString(data map[string]any, key string) string {
	if data == nil {
		return ""
	}
	s, _ := data[key].(string)
	return s
}

func dataInt(data map[string]any, key string) int {
	if data == nil {
		return 0
	}
	switch v := data[key].(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	}
	return 0
}
