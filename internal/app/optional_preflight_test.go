package app

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/quota"
)

// V3 (§9.0): one usable role holder passes the kickoff floor; v2 still stops.
func TestOptionalParticipantsPreflightFloor(t *testing.T) {
	for _, tc := range []struct {
		name     string
		policy   quota.Policy
		hardStop bool
	}{
		{"v2", quota.NewExhaustionPolicy(nil, nil), true},
		{"v3", quota.NewOptionalParticipantPolicy(nil, nil), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := tc.policy
			root, _, _ := fixupAppFixture(t, p)
			ids := []string{"alpha", "beta", "gamma"}
			snapshot := &quota.RoleSnapshot{Version: 1, Order: ids, Organizer: "alpha"}
			roles := quota.Roles{Snapshot: snapshot, Assignment: quota.RoleAssignment{Organizer: "alpha", Implementer: "alpha", ImplementerSource: "legacy"}}
			var discovered []agents.Discovery
			for _, id := range ids {
				mode := "transient"
				if id == "alpha" {
					mode = "success"
				}
				discovered = append(discovered, agents.Discovery{Spec: agents.Spec{ID: id, HeadlessArgs: []string{"-test.run=^TestExhaustionReadinessHelper$", "--", "--exhaustion-readiness-helper", mode, ""}, PromptMode: agents.PromptStdin}, Found: true, Path: os.Args[0]})
			}
			report, code, err := preflight(context.Background(), preflightOptions{Root: root, QuotaPolicy: &p, QuotaRoles: roles, PingTimeout: time.Second}, discovered, io.Discard, io.Discard)
			if err != nil {
				t.Fatal(err)
			}
			if (code == 1) != tc.hardStop {
				t.Fatalf("code=%d decision=%+v", code, report.QuotaDecision)
			}
			if !tc.hardStop && (report.QuotaDecision == nil || !report.QuotaDecision.Applied || len(report.QuotaDecision.After) != 1 || report.QuotaDecision.After[0] != "alpha") {
				t.Fatalf("v3 decision=%+v", report.QuotaDecision)
			}
		})
	}
}
