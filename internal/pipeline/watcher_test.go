package pipeline

import (
	"testing"
	"time"
)

func TestBreachDedupeByFingerprint(t *testing.T) {
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	b1 := Breach{Signal: "p99", Target: "api", Threshold: "<500ms", Observed: "900ms", Class: "latency", At: now}
	b2 := Breach{Signal: "p99", Target: "api", Threshold: "<500ms", Observed: "950ms", Class: "latency", At: now.Add(time.Minute)}
	b3 := Breach{Signal: "5xx", Target: "api", Threshold: "<1%", Observed: "4%", Class: "errors", At: now}
	if b1.Fingerprint() != b2.Fingerprint() {
		t.Fatal("same signal/target/threshold/class must share a fingerprint")
	}
	got := DedupeBreaches([]Breach{b1, b2, b3})
	if len(got) != 2 {
		t.Fatalf("dedupe = %d, want 2", len(got))
	}
}

func TestCanAutoOpenOnlyLowRiskClasses(t *testing.T) {
	m := Monitoring{BreachClasses: []BreachClass{
		{Name: "latency", Risk: RiskLow, AutoOpen: true},
		{Name: "errors", Risk: RiskHigh, AutoOpen: true},
		{Name: "data-loss", Risk: RiskProduction, AutoOpen: true},
		{Name: "noise", Risk: RiskLow, AutoOpen: false},
	}}
	cases := map[string]bool{
		"latency":   true,  // predeclared low-risk + auto-open
		"errors":    false, // high risk never auto-opens
		"data-loss": false, // production never auto-opens
		"noise":     false, // auto-open disabled
		"unknown":   false, // unknown class gates
	}
	for class, want := range cases {
		got := m.CanAutoOpen(Breach{Class: class})
		if got != want {
			t.Errorf("CanAutoOpen(%s) = %v, want %v", class, got, want)
		}
	}
}
