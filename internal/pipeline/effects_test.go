package pipeline

import (
	"testing"
	"time"
)

func TestIdempotencyKeyStability(t *testing.T) {
	a := IdempotencyKey("p", "deploy", "vercel", "deploy.production", "app", HashRequest([]byte(`{"x":1}`)))
	b := IdempotencyKey("p", "deploy", "vercel", "deploy.production", "app", HashRequest([]byte(`{"x":1}`)))
	if a != b {
		t.Fatalf("same inputs must yield same key:\n%s\n%s", a, b)
	}
	if KeyDigest(a) != KeyDigest(b) {
		t.Fatal("same key must yield same digest")
	}
	c := IdempotencyKey("p", "deploy", "vercel", "deploy.production", "OTHER", HashRequest([]byte(`{"x":1}`)))
	if a == c {
		t.Fatal("different target must yield a different key")
	}
	d := IdempotencyKey("p", "deploy", "vercel", "deploy.production", "app", HashRequest([]byte(`{"x":2}`)))
	if a == d {
		t.Fatal("different request body must yield a different key")
	}
}

func TestEffectRoundTripAndReconcile(t *testing.T) {
	deck := t.TempDir()
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	e := NewEffect("demo", "deploy", "vercel", "deploy.production", "app", HashRequest([]byte("{}")), RiskProduction, now)
	if e.Status != EffectPlanned || len(e.Attempts) != 1 {
		t.Fatalf("new effect = %+v", e)
	}
	if err := SaveEffect(deck, e); err != nil {
		t.Fatalf("SaveEffect: %v", err)
	}

	e.Advance(EffectExecuting, "dpl_123", "calling provider", now.Add(time.Second))
	if !e.NeedsReconcile() {
		t.Fatal("executing effect must need reconcile")
	}
	if err := SaveEffect(deck, e); err != nil {
		t.Fatalf("SaveEffect(executing): %v", err)
	}

	got, ok, err := LoadEffect(deck, "demo", e.IdempotencyKey)
	if err != nil || !ok {
		t.Fatalf("LoadEffect ok=%v err=%v", ok, err)
	}
	if got.ExternalRef != "dpl_123" {
		t.Fatalf("external_ref = %q", got.ExternalRef)
	}
	if len(got.Attempts) != 2 {
		t.Fatalf("attempts = %d, want 2 (append-only)", len(got.Attempts))
	}
	if !got.NeedsReconcile() {
		t.Fatal("reloaded executing effect must need reconcile")
	}

	got.Advance(EffectSucceeded, "dpl_123", "confirmed via get_deployment", now.Add(2*time.Second))
	if got.NeedsReconcile() {
		t.Fatal("succeeded effect must not need reconcile")
	}
}

func TestBlockHasSucceededEffect(t *testing.T) {
	deck := t.TempDir()
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	// A planned (not succeeded) effect for "deploy" -> not complete.
	planned := NewEffect("p", "deploy", "vercel", "deploy.production", "app", HashRequest([]byte("{}")), RiskProduction, now)
	if err := SaveEffect(deck, planned); err != nil {
		t.Fatal(err)
	}
	if ok, _ := BlockHasSucceededEffect(deck, "p", "deploy"); ok {
		t.Fatal("planned effect must NOT count as a succeeded block effect")
	}
	// Advance it to succeeded.
	planned.Advance(EffectSucceeded, "dpl_1", "done", now)
	if err := SaveEffect(deck, planned); err != nil {
		t.Fatal(err)
	}
	if ok, _ := BlockHasSucceededEffect(deck, "p", "deploy"); !ok {
		t.Fatal("succeeded effect must count")
	}
	// A different block has no effect.
	if ok, _ := BlockHasSucceededEffect(deck, "p", "other"); ok {
		t.Fatal("unrelated block must not be complete")
	}
}

func TestSameKeyOverwritesSameFile(t *testing.T) {
	deck := t.TempDir()
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	e1 := NewEffect("demo", "b", "vercel", "deploy.preview", "app", HashRequest([]byte("{}")), RiskNormal, now)
	e2 := NewEffect("demo", "b", "vercel", "deploy.preview", "app", HashRequest([]byte("{}")), RiskNormal, now)
	if EffectPath(deck, "demo", e1.IdempotencyKey) != EffectPath(deck, "demo", e2.IdempotencyKey) {
		t.Fatal("identical actions must map to the same effect file (no duplicate)")
	}
}
