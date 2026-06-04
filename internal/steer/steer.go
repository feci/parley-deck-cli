// Package steer records mid-run steering / follow-up prompts as durable events
// on a run's event log. A steering request is the user giving an agent (or the
// whole deck) its next instruction after a task has started.
//
// This slice (Slice 4 of tui-interactivity-overhaul) records a QUEUED new
// attempt; it does NOT execute delivery. Recording intent performs no
// side-effect, so it bypasses no HITL/action gate — actual delivery (live
// acp_prompt, native resume, or launching the queued attempt) is gated and is
// deferred to a later slice. The driver, not the TUI, will own that execution.
package steer

import (
	"fmt"
	"strings"
	"time"

	"parley-deck-cli/internal/store"
)

// Target is who a steering request is aimed at.
type Target string

const (
	TargetAgent Target = "agent"
	TargetDeck  Target = "deck"
)

const (
	// DeliveryQueuedNewAttempt is the only delivery mode in this slice: the steer
	// is queued for the agent's next attempt. Live delivery modes (acp_prompt,
	// native_resume) are deferred.
	DeliveryQueuedNewAttempt = "new_attempt"

	StatusQueued = "queued"
)

// Request is a steering instruction submitted from the TUI composer or the
// `parley steer` CLI.
type Request struct {
	Target    Target
	Agent     string // required when Target == TargetAgent
	Text      string
	CreatedBy string // "tui" | "cli"
	SegmentID string // current segment of the targeted agent, if known
	Risk      string // hitl risk vocabulary; defaults to "normal"
}

// Result reports the assigned steer id and how it was queued.
type Result struct {
	ID           string `json:"id"`
	DeliveryMode string `json:"delivery_mode"`
	Status       string `json:"status"`
}

// Queued is a steering request projected from the event log for display.
type Queued struct {
	ID        string `json:"id"`
	Target    string `json:"target"`
	Agent     string `json:"agent,omitempty"`
	Text      string `json:"text"`
	SegmentID string `json:"segment_id,omitempty"`
	Status    string `json:"status"`
	Mode      string `json:"mode,omitempty"`
}

func (r Request) validate() error {
	if strings.TrimSpace(r.Text) == "" {
		return fmt.Errorf("steering text is required")
	}
	switch r.Target {
	case TargetAgent:
		if strings.TrimSpace(r.Agent) == "" {
			return fmt.Errorf("an agent is required when steering a single agent")
		}
	case TargetDeck:
	default:
		return fmt.Errorf("unknown steer target %q", r.Target)
	}
	return nil
}

// Submit records a steering request as durable steer.requested + steer.delivered
// events on the run's event log and returns the assigned id. It queues the steer
// for the agent's next attempt; it does not execute it.
func Submit(runDir string, req Request, now time.Time) (Result, error) {
	if err := req.validate(); err != nil {
		return Result{}, err
	}
	risk := strings.TrimSpace(req.Risk)
	if risk == "" {
		risk = "normal"
	}
	s := store.New(runDir)
	id := nextSteerID(s)
	if err := s.Append(store.Event{
		Time: now,
		Type: "steer.requested",
		Data: map[string]any{
			"id":         id,
			"target":     string(req.Target),
			"agent":      req.Agent,
			"text":       req.Text,
			"created_by": req.CreatedBy,
			"segment_id": req.SegmentID,
			"risk":       risk,
		},
	}); err != nil {
		return Result{}, err
	}
	if err := s.Append(store.Event{
		Time: now,
		Type: "steer.delivered",
		Data: map[string]any{
			"id":     id,
			"mode":   DeliveryQueuedNewAttempt,
			"status": StatusQueued,
		},
	}); err != nil {
		return Result{ID: id}, err
	}
	return Result{ID: id, DeliveryMode: DeliveryQueuedNewAttempt, Status: StatusQueued}, nil
}

// nextSteerID returns the next monotonic steer-NNNN id by counting prior
// steer.requested events.
func nextSteerID(s store.Store) string {
	events, err := s.Load()
	if err != nil {
		return "steer-0001"
	}
	n := 0
	for _, e := range events {
		if e.Type == "steer.requested" {
			n++
		}
	}
	return fmt.Sprintf("steer-%04d", n+1)
}

// List projects the queued steering requests from a run's event log, in
// submission order.
func List(runDir string) ([]Queued, error) {
	events, err := store.New(runDir).Load()
	if err != nil {
		return nil, err
	}
	return project(events), nil
}

func project(events []store.Event) []Queued {
	byID := map[string]*Queued{}
	order := []string{}
	for _, e := range events {
		id := dataString(e.Data, "id")
		if id == "" {
			continue
		}
		switch e.Type {
		case "steer.requested":
			if _, ok := byID[id]; !ok {
				order = append(order, id)
			}
			byID[id] = &Queued{
				ID:        id,
				Target:    dataString(e.Data, "target"),
				Agent:     dataString(e.Data, "agent"),
				Text:      dataString(e.Data, "text"),
				SegmentID: dataString(e.Data, "segment_id"),
				Status:    StatusQueued,
			}
		case "steer.delivered":
			if q, ok := byID[id]; ok {
				if mode := dataString(e.Data, "mode"); mode != "" {
					q.Mode = mode
				}
				if status := dataString(e.Data, "status"); status != "" {
					q.Status = status
				}
			}
		}
	}
	out := make([]Queued, 0, len(order))
	for _, id := range order {
		out = append(out, *byID[id])
	}
	return out
}

func dataString(data map[string]any, key string) string {
	if data == nil {
		return ""
	}
	if v, ok := data[key].(string); ok {
		return v
	}
	return ""
}
