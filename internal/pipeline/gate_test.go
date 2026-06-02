package pipeline

import (
	"testing"
	"time"
)

func TestAutoApproveProductionNeverAutoResolves(t *testing.T) {
	cases := []struct {
		autonomy Autonomy
		risk     Risk
		want     bool
	}{
		{AutonomySupervised, RiskLow, false},
		{AutonomySupervised, RiskProduction, false},
		{AutonomyAutoLeft, RiskLow, true},
		{AutonomyAutoLeft, RiskNormal, true},
		{AutonomyAutoLeft, RiskHigh, false},
		{AutonomyAutoLeft, RiskProduction, false},
	}
	for _, c := range cases {
		if got := AutoApprove(c.autonomy, c.risk); got != c.want {
			t.Errorf("AutoApprove(%s, %s) = %v, want %v", c.autonomy, c.risk, got, c.want)
		}
	}
}

func TestGateRoundTripAndResolve(t *testing.T) {
	deck := t.TempDir()
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	g := NewGate("demo", "a", "b", RiskProduction, "release", now)
	if g.Status != GateOpen {
		t.Fatalf("new gate status = %q, want open", g.Status)
	}
	if err := SaveGate(deck, g); err != nil {
		t.Fatalf("SaveGate: %v", err)
	}
	loaded, ok, err := LoadGate(deck, "demo", EdgeID("a", "b"))
	if err != nil || !ok {
		t.Fatalf("LoadGate ok=%v err=%v", ok, err)
	}
	if loaded.Risk != RiskProduction {
		t.Fatalf("risk = %q", loaded.Risk)
	}
	loaded.Resolve(true, "tomas", now)
	if err := SaveGate(deck, loaded); err != nil {
		t.Fatalf("SaveGate(resolved): %v", err)
	}
	reloaded, _, _ := LoadGate(deck, "demo", EdgeID("a", "b"))
	if reloaded.Status != GateApproved || reloaded.ApprovedBy != "tomas" {
		t.Fatalf("resolved gate = %+v", reloaded)
	}
}

func TestLoadGateMissing(t *testing.T) {
	_, ok, err := LoadGate(t.TempDir(), "demo", "x->y")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if ok {
		t.Fatal("expected missing gate ok=false")
	}
}
