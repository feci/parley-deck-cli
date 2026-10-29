package runner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/quotatest"
	"parley-deck-cli/internal/runmanifest"
	"parley-deck-cli/internal/store"
)

func TestSelectedHandoverChildAndReplay(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PARLEY_HOME", t.TempDir())
	writeLaunchProtocol(t, root)
	p := quota.NewExhaustionPolicy(nil, nil)
	snapshot := &quota.RoleSnapshot{Version: 1, Order: []string{"old-1", "test-1", "peer-1"}, Organizer: "old-1"}
	idea, k, err := protocol.CreateIdeaWithQuota(root, "handover fixture", snapshot.Order, nil, "deliberation", "", "run", &p, nil, snapshot)
	if err != nil {
		t.Fatal(err)
	}
	m := runmanifest.New(runmanifest.Options{Root: root, RunID: "run", IdeaSlug: idea.Slug, Participants: idea.Participants, QuotaKickoff: k})
	if err := runmanifest.Write(root, "run", m); err != nil {
		t.Fatal(err)
	}
	h, _ := quota.ReadHistory(idea.Path)
	e := quotatest.Exhaustion(idea.Slug, "old-1", "round-02")
	d := quota.Evaluate(p, h.Current, []quota.Member{{ID: "old-1", Evidence: &e}, {ID: "test-1", Usable: true}, {ID: "peer-1", Usable: true}}, quota.Roles{Snapshot: snapshot, Assignment: quota.RoleAssignment{Organizer: "old-1", Implementer: "old-1", ImplementerSource: "designation"}})
	d = d.BindHandover("consensus", "none", "", "test", "abc123")
	if err := quota.CommitBatch(idea.Path, quota.NewBatch(h, "run", "consensus", nil, d, time.Now())); err != nil {
		t.Fatal(err)
	}
	r := d.RoleTransition
	raw := fmt.Sprintf("---\nidea: %s\ntransition: %s\nauthor: test-1\norganizer: test-1\nimplementer: test-1\nphase: consensus\nbase: none\nbase-sha256:\ntakeover-branch: test\ntakeover-head: abc123\ndate: 2026-10-29\nstatus: in-progress\n---\n\n## Takeover\nBoth roles transfer.\n\n## Resume state\nContinue consensus at consensus.md.\n", idea.Slug, r.ID)
	source := filepath.Join(root, "fixture.md")
	target := filepath.Join(idea.Path, r.Path())
	os.WriteFile(source, []byte(raw), 0600)
	quote := func(s string) string { return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'" }
	agent := telemetryShell("cat >/dev/null; cp "+quote(source)+" "+quote(target), false)
	agent.Found = true
	opts := Options{Root: root, RunID: "run", Idea: idea, Agents: []agents.Discovery{agent}, Store: store.New(filepath.Join(root, protocol.DeckDir, "runs", "run"))}
	if err := RunRoleHandover(context.Background(), opts); err != nil {
		t.Fatal(err)
	}
	if err := RunRoleHandover(context.Background(), opts); err != nil {
		t.Fatal(err)
	}
	records := terminalRecords(t, root)
	if len(records) != 1 || records[0].Metadata.Agent != "test-1" || records[0].Metadata.ParticipantLogicalStep != "handover:"+r.ID {
		t.Fatalf("records: %+v", records)
	}
	if started, err := protocol.ExecutionStarted(idea.Path); started || err != nil {
		t.Fatalf("handover implied Phase 5: %v %v", started, err)
	}
}
