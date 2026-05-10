package store

import (
	"testing"
	"time"
)

func TestAppendAndLoad(t *testing.T) {
	s := New(t.TempDir())
	now := time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)
	if err := s.Append(Event{Time: now, Type: "run.created", Data: map[string]any{"mode": "hitl"}}); err != nil {
		t.Fatal(err)
	}

	events, err := s.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Fatalf("len(events) = %d", len(events))
	}
	if events[0].Type != "run.created" {
		t.Fatalf("type = %q", events[0].Type)
	}
}
