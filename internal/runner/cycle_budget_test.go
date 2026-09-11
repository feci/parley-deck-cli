package runner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/store"
)

func failedCycleAgent(id string) agents.Discovery {
	return agents.Discovery{Spec: agents.Spec{ID: id, HeadlessArgs: []string{"-test.run=TestFakeFailingAgentHelper", "--", "parley-failing-agent"}, PromptMode: agents.PromptStdin}, Path: os.Args[0], Found: true}
}

func TestCycleManualFixupSixthAttemptRefusedAndReviewStaysFree(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	declareTestLaunchSource(t, root)
	idea, err := protocol.CreateIdea(root, "Manual fixup cap", []string{"builder", "reviewer"})
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(idea.Path, "00-prompt.md"), "---\nidea: "+idea.Slug+"\nparticipants: [builder, reviewer]\ntrack: deliberation\nstatus: implemented\n---\nFixture\n")
	idea.Participants = []string{"builder"}
	for attempt := 0; attempt < 6; attempt++ {
		runID := fmt.Sprintf("fixup-%d", attempt)
		result := RunFixup(context.Background(), Options{Root: root, Idea: idea, RunID: runID, Agents: []agents.Discovery{failedCycleAgent("builder")}, Timeout: 5 * time.Second, Store: store.New(filepath.Join(root, protocol.DeckDir, "runs", runID))})
		if result.ExitError == "" {
			t.Fatalf("failed operation passed: %+v", result)
		}
	}
	started, denied := 0, 0
	for _, r := range terminalRecords(t, root) {
		if r.StartedAt != nil {
			started++
		}
		if r.Outcome.FailureClass != nil && *r.Outcome.FailureClass == "budget_refused" {
			denied++
		}
	}
	if started != 5 || denied != 1 {
		t.Fatalf("sixth fixup boundary: started=%d refused=%d", started, denied)
	}
	b, err := budget.LoadCycleBinding(context.Background(), root, idea.Slug, budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	// A verification process after the fifth allowed fixup does not consume a
	// sixth fixup. Its actual process failure must remain a process failure.
	r, err := RunMeasured(context.Background(), ExecOptions{Root: root, Agent: failedCycleAgent("reviewer"), Prompt: "verification fixture", Timeout: time.Second, Info: LaunchInfo{Idea: idea.Slug, Phase: "review"}})
	if err == nil || r.StartedAt == nil {
		t.Fatalf("verification refused by fixup cap: %+v %v", r, err)
	}
	s, err := b.Store.Inspect(context.Background())
	if err != nil || b.Count(s) != 5 {
		t.Fatalf("verification charged a fixup: %+v %v", s, err)
	}
}

func TestCycleCrossReviewGroupThirdAllowedFourthRefused(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	declareTestLaunchSource(t, root)
	idea, err := protocol.CreateIdea(root, "Grouped cross review cap", []string{"alpha", "beta"})
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(idea.Path, "00-prompt.md"), "---\nidea: "+idea.Slug+"\nparticipants: [alpha, beta]\ntrack: deliberation\nstatus: round-01\n---\nFixture\n")
	for _, id := range idea.Participants {
		mustWrite(t, filepath.Join(idea.Path, "round-01", id+".md"), "---\nagent: "+id+"\nidea: "+idea.Slug+"\nround: 1\n---\n\n## Summary\nFixture\n")
	}
	for attempt := 0; attempt < 4; attempt++ {
		runID := fmt.Sprintf("cross-%d", attempt)
		results := RunRound(context.Background(), Options{Root: root, Idea: idea, RunID: runID, Round: 2, Agents: []agents.Discovery{failedCycleAgent("alpha"), failedCycleAgent("beta")}, Timeout: 5 * time.Second, Store: store.New(filepath.Join(root, protocol.DeckDir, "runs", runID))})
		if len(results) != 2 {
			t.Fatalf("group result count: %+v", results)
		}
		for _, r := range results {
			if r.ExitError == "" {
				t.Fatalf("failed child passed: %+v", r)
			}
		}
	}
	started, denied := 0, 0
	for _, r := range terminalRecords(t, root) {
		if r.StartedAt != nil {
			started++
		}
		if r.Outcome.FailureClass != nil && *r.Outcome.FailureClass == "budget_refused" {
			denied++
		}
	}
	if started != 6 || denied != 2 {
		t.Fatalf("group budget charged per child or reset: started=%d refused=%d", started, denied)
	}
	b, err := budget.LoadCycleBinding(context.Background(), root, idea.Slug, budget.CrossReview)
	if err != nil {
		t.Fatal(err)
	}
	s, err := b.Store.Inspect(context.Background())
	if err != nil || b.Count(s) != 3 {
		t.Fatalf("shared cross-review count: %+v %v", s, err)
	}
}

func TestCycleCeilingsUseStrictTrackInputs(t *testing.T) {
	for _, tc := range []struct {
		header     string
		fix, cross int
		fail       bool
	}{
		{"track: deliberation\n", 5, 3, false},
		{"track: standard\n", 2, 2, false},
		{"track: fast\n", 1, 0, false},
		{"", 3, 4, false},
		{"track: fast\nauto_implement: TRUE\n", 0, 0, true},
		{"track: standard\nstrict_gate: true\n", 0, 0, true},
		{"track: fast\ntrack: deliberation\n", 0, 0, true},
		{"track: null\n", 0, 0, true},
		{"track: standrd\n", 0, 0, true},
	} {
		t.Run(tc.header, func(t *testing.T) {
			dir := t.TempDir()
			mustWrite(t, filepath.Join(dir, "00-prompt.md"), "---\nparticipants: [a, b]\n"+tc.header+"---\nFixture\n")
			f, c, _, err := cycleCeilings(dir)
			if (err != nil) != tc.fail || (!tc.fail && (f != tc.fix || c != tc.cross)) {
				t.Fatalf("fix=%d cross=%d err=%v", f, c, err)
			}
		})
	}
}

func TestCycleMissingPromptCannotFreezeLegacyGrant(t *testing.T) {
	if _, _, _, err := cycleCeilings(t.TempDir()); err == nil {
		t.Fatal("missing track authority froze a legacy grant")
	}
}

func TestCycleHandoffDoesNotActivatePolicy(t *testing.T) {
	root := t.TempDir()
	writeLaunchProtocol(t, root)
	ideaDir := filepath.Join(root, "parley-deck", "ideas", "idea")
	mustWrite(t, filepath.Join(ideaDir, "00-prompt.md"), "---\nparticipants: [a, b]\ntrack: deliberation\n---\nFixture\n")
	packet, err := WriteHandoffPacket(HandoffOptions{Root: root, RunID: "handoff", Idea: "idea", Phase: "fixup", Agent: failedCycleAgent("a"), Prompt: "fixture"})
	if err != nil {
		t.Fatal(err)
	}
	if packet.InvocationID == "" {
		t.Fatal("handoff lacks evidence")
	}
	b, err := budget.LoadCycleBinding(context.Background(), root, "idea", budget.Fixup)
	if err != nil || b != nil {
		t.Fatalf("unobserved handoff activated policy: %+v %v", b, err)
	}
	records := terminalRecords(t, root)
	if len(records) != 1 || records[0].StartedAt != nil || records[0].Outcome.Status != "unobserved-handoff" {
		t.Fatalf("handoff claimed execution: %+v", records)
	}
}

func TestCycleCancelledLaunchKeepsTerminalAndNoCharge(t *testing.T) {
	root := t.TempDir()
	writeLaunchProtocol(t, root)
	ideaDir := filepath.Join(root, "parley-deck", "ideas", "idea")
	mustWrite(t, filepath.Join(ideaDir, "00-prompt.md"), "---\nparticipants: [a, b]\ntrack: deliberation\n---\nFixture\n")
	ctx := context.Background()
	b, err := budget.EnsureCycleBinding(ctx, root, "idea", budget.Fixup, 5, 0, "", ideaDir)
	if err != nil {
		t.Fatal(err)
	}
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	r, err := RunMeasured(cancelled, ExecOptions{Root: root, Agent: failedCycleAgent("a"), Prompt: "fixture", Timeout: time.Second, Info: LaunchInfo{Idea: "idea", Phase: "fixup"}})
	if !errors.Is(err, context.Canceled) || r.StartedAt != nil {
		t.Fatalf("cancelled launch: %+v %v", r, err)
	}
	records := terminalRecords(t, root)
	if len(records) != 1 || records[0].Outcome.FailureClass == nil || *records[0].Outcome.FailureClass != "cancelled" {
		t.Fatalf("cancellation lost terminal: %+v", records)
	}
	s, err := b.Store.Inspect(ctx)
	if err != nil || b.Count(s) != 0 {
		t.Fatalf("cancelled request charged: %+v %v", s, err)
	}
}
