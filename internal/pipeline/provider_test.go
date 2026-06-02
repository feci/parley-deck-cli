package pipeline

import (
	"errors"
	"testing"
)

func TestVercelSupports(t *testing.T) {
	p := VercelProvider{}
	for _, c := range []Capability{CapDeployPreview, CapDeployProduction, CapRuntimeRollback, CapMonitorAlert} {
		if !p.Supports(c) {
			t.Errorf("vercel should support %s", c)
		}
	}
	if p.Supports(CapIssueCreate) {
		t.Error("vercel should not support issue.create")
	}
}

func TestProductionGateGuard(t *testing.T) {
	p := VercelProvider{}
	// Real production execution without a gate must be refused.
	if _, err := p.Plan(ActionRequest{Capability: CapDeployProduction, Target: "app", DryRun: false, GateApproved: false}); !errors.Is(err, ErrProductionGateRequired) {
		t.Fatalf("expected ErrProductionGateRequired, got %v", err)
	}
	// Dry-run production is allowed (no mutation).
	if _, err := p.Plan(ActionRequest{Capability: CapDeployProduction, Target: "app", DryRun: true}); err != nil {
		t.Fatalf("dry-run production should be allowed: %v", err)
	}
	// Production with an approved gate is allowed and marked production.
	call, err := p.Plan(ActionRequest{Capability: CapDeployProduction, Target: "app", GateApproved: true})
	if err != nil {
		t.Fatalf("gated production should be allowed: %v", err)
	}
	if call.Args["production"] != "true" {
		t.Fatalf("expected production arg, got %v", call.Args)
	}
}

func TestRollbackIsProductionGuarded(t *testing.T) {
	p := VercelProvider{}
	if _, err := p.Plan(ActionRequest{Capability: CapRuntimeRollback, Target: "app"}); !errors.Is(err, ErrProductionGateRequired) {
		t.Fatalf("rollback must be production-gated, got %v", err)
	}
}

func TestPlanPreviewAndUnsupported(t *testing.T) {
	p := VercelProvider{}
	call, err := p.Plan(ActionRequest{Capability: CapDeployPreview, Target: "app", DryRun: true})
	if err != nil {
		t.Fatalf("preview plan: %v", err)
	}
	if call.Tool == "" || !call.DryRun || call.Provider != "vercel" {
		t.Fatalf("unexpected call %+v", call)
	}
	if _, err := p.Plan(ActionRequest{Capability: CapIssueCreate, Target: "x"}); !errors.Is(err, ErrUnsupportedCapability) {
		t.Fatalf("expected ErrUnsupportedCapability, got %v", err)
	}
}

func TestNoopProviderStillGuardsProduction(t *testing.T) {
	p := NoopProvider{}
	if !p.Supports(CapNotifySend) {
		t.Fatal("noop supports everything")
	}
	if _, err := p.Plan(ActionRequest{Capability: CapDeployProduction, Target: "x"}); !errors.Is(err, ErrProductionGateRequired) {
		t.Fatalf("noop must still guard production, got %v", err)
	}
}
