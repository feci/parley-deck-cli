package app

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/quota"
)

func TestQuotaPreflightReportsOnlyAndBare503(t *testing.T) {
	root := sourceWorkspace(t)
	old := pingProbe
	t.Cleanup(func() { pingProbe = old })
	pingProbe = func(_ context.Context, _ string, a agents.Discovery, _ time.Duration) readinessObservation {
		if a.ID == "c" {
			return classifyReadiness("", "503", 1, false, false, "")
		}
		return readinessObservation{Class: ClassReady, Ready: true}
	}
	ds := []agents.Discovery{{Spec: agents.Spec{ID: "a"}, Found: true}, {Spec: agents.Spec{ID: "b"}, Found: true}, {Spec: agents.Spec{ID: "c"}, Found: true}}
	var out, errs bytes.Buffer
	report, code, err := preflight(context.Background(), preflightOptions{Root: root, Yes: true}, ds, &out, &errs)
	if err != nil || code != 3 || len(report.Excluded) > 0 || len(report.Gates) != 1 || report.Gates[0].Kind != gateProviderFailure {
		t.Fatal(report, code, err)
	}
	if providerFailureClass("503") != "overloaded" {
		t.Fatal("bare503 bypass")
	}
	files, _ := os.ReadDir(filepath.Join(root, "parley-deck", "ideas"))
	if len(files) > 0 {
		t.Fatal("standalone preflight created an idea")
	}
}
func TestQuotaPreflightWholeBatchAndOneBlock(t *testing.T) {
	root := sourceWorkspace(t)
	os.MkdirAll(filepath.Join(root, "parley-deck", "inbox"), 0755)
	old := pingProbe
	t.Cleanup(func() { pingProbe = old })
	now := time.Now().UTC()
	pingProbe = func(_ context.Context, _ string, a agents.Discovery, _ time.Duration) readinessObservation {
		if a.ID == "c" || a.ID == "d" {
			return readinessObservation{Class: ClassProviderFailure, QuotaEvidence: &quota.Evidence{Eligible: true, InvocationID: a.ID, RuleID: "synthetic-policy-test", Provenance: "synthetic-test", ObservedAt: now}}
		}
		return readinessObservation{Class: ClassReady, Ready: true}
	}
	ds := []agents.Discovery{{Spec: agents.Spec{ID: "a"}, Found: true}, {Spec: agents.Spec{ID: "b"}, Found: true}, {Spec: agents.Spec{ID: "c"}, Found: true}, {Spec: agents.Spec{ID: "d"}, Found: true}}
	var out, errs bytes.Buffer
	p := quota.NewPolicy(nil, nil)
	report, code, err := preflight(context.Background(), preflightOptions{Root: root, QuotaPolicy: &p}, ds, &out, &errs)
	if err != nil || code != 0 || !report.QuotaDecision.Applied || len(report.QuotaDecision.After) != 2 {
		t.Fatal(report, code, err)
	}
	// Standalone command sees identical candidates but never applies them.
	report, code, err = preflight(context.Background(), preflightOptions{Root: root, Yes: true}, ds, &out, &errs)
	if err != nil || code != 3 || report.QuotaDecision != nil || len(report.Excluded) != 0 {
		t.Fatal(report, code, err)
	}
	for i := 0; i < 2; i++ {
		report, _, err = preflight(context.Background(), preflightOptions{Root: root, QuotaPolicy: &p}, ds[1:], &out, &errs)
		if err != nil || report.QuotaDecision.Applied || !strings.Contains(report.QuotaDecision.Block, "1 usable") {
			t.Fatal(report, err)
		}
	}
	notes, _ := filepath.Glob(filepath.Join(root, "parley-deck", "inbox", "parley-to-user_quota-kickoff-*.md"))
	if len(notes) != 1 {
		t.Fatal(notes)
	}
}
