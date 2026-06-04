package steer

import (
	"testing"
	"time"
)

func TestSubmitQueuesAgentSteerWithMonotonicIDs(t *testing.T) {
	dir := t.TempDir()
	now := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)

	first, err := Submit(dir, Request{Target: TargetAgent, Agent: "codex", Text: "focus on the parser", CreatedBy: "tui", SegmentID: "segment-0002"}, now)
	if err != nil {
		t.Fatal(err)
	}
	if first.ID != "steer-0001" || first.DeliveryMode != DeliveryQueuedNewAttempt || first.Status != StatusQueued {
		t.Fatalf("first result = %+v", first)
	}

	second, err := Submit(dir, Request{Target: TargetDeck, Text: "wrap up", CreatedBy: "cli"}, now.Add(time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if second.ID != "steer-0002" {
		t.Fatalf("second id = %q, want steer-0002 (monotonic)", second.ID)
	}

	queued, err := List(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(queued) != 2 {
		t.Fatalf("queued = %d, want 2", len(queued))
	}
	if queued[0].ID != "steer-0001" || queued[0].Agent != "codex" || queued[0].Target != "agent" {
		t.Fatalf("queued[0] = %+v", queued[0])
	}
	if queued[0].SegmentID != "segment-0002" || queued[0].Status != StatusQueued || queued[0].Mode != DeliveryQueuedNewAttempt {
		t.Fatalf("queued[0] = %+v", queued[0])
	}
	if queued[1].Target != "deck" || queued[1].Text != "wrap up" {
		t.Fatalf("queued[1] = %+v", queued[1])
	}
}

func TestSubmitValidates(t *testing.T) {
	dir := t.TempDir()
	now := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)

	if _, err := Submit(dir, Request{Target: TargetAgent, Agent: "codex", Text: "  "}, now); err == nil {
		t.Fatal("empty text must be rejected")
	}
	if _, err := Submit(dir, Request{Target: TargetAgent, Text: "do x"}, now); err == nil {
		t.Fatal("agent target without an agent must be rejected")
	}
	if _, err := Submit(dir, Request{Target: Target("bogus"), Text: "do x"}, now); err == nil {
		t.Fatal("unknown target must be rejected")
	}
}
