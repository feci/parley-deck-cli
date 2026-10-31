package driver

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"parley-deck-cli/internal/membership"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/runmanifest"
	"parley-deck-cli/internal/store"
	"parley-deck-cli/internal/telemetry"
)

// optionalIdea creates a v3 idea with participants a..d and, when reduce is
// set, settles typed failures of b, c and d so only the implementer remains.
func optionalIdea(t *testing.T, p quota.Policy, reduce bool) (string, string, string) {
	t.Helper()
	root := t.TempDir()
	t.Setenv("PARLEY_HOME", t.TempDir())
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	ids := []string{"a", "b", "c", "d"}
	idea, k, err := protocol.CreateIdeaWithQuota(root, "optional participants", ids, nil, "deliberation", "", "run-1", &p, nil)
	if err != nil {
		t.Fatal(err)
	}
	run := "run-1"
	if err = runmanifest.Write(root, run, runmanifest.New(runmanifest.Options{Root: root, RunID: run, IdeaSlug: idea.Slug, Participants: idea.Participants, QuotaKickoff: k})); err != nil {
		t.Fatal(err)
	}
	s := store.New(filepath.Join(root, protocol.DeckDir, "runs", run))
	if err = s.AppendDurable(store.Event{Type: "run.created", Data: map[string]any{"idea": idea.Slug, "participants": idea.Participants, "quota_kickoff": k}}); err != nil {
		t.Fatal(err)
	}
	body := fmt.Sprintf("---\nagent: a\nidea: %s\nround: 1\n---\n", idea.Slug)
	for _, h := range []string{"Summary", "Proposed approach", "Concerns / open questions", "Risks", "Existing alternatives"} {
		body += "\n## " + h + "\nConcrete fixture evidence.\n"
	}
	if err = os.WriteFile(filepath.Join(idea.Path, "round-01", "a.md"), []byte(body), 0600); err != nil {
		t.Fatal(err)
	}
	if reduce {
		ms := []quota.Member{{ID: "a", Usable: true}}
		for _, id := range []string{"b", "c", "d"} {
			now := time.Now().UTC()
			code := 1
			first := quota.FailedAttempt{InvocationID: id + "-1", Status: "failed", FailureClass: "process_failure", ExitCode: &code, ObservedAt: now, Excerpt: "exit_code: 1"}
			second := first
			second.InvocationID, second.RetryOf, second.ObservedAt = id+"-2", first.InvocationID, now.Add(5*time.Second)
			e, err := telemetry.ParticipantEvidence(idea.Slug, id, "round-01/"+id+".md", "fixture", []quota.FailedAttempt{first, second})
			if err != nil {
				t.Fatal(err)
			}
			ms = append(ms, quota.Member{ID: id, Evidence: &e})
		}
		ctx, release, err := membership.Acquire(context.Background(), idea.Path, run)
		if err != nil {
			t.Fatal(err)
		}
		defer release()
		if b, err := membership.Settle(ctx, root, idea.Path, run, "round-01", ids, ms); err != nil || b == nil {
			t.Fatalf("settle: %+v %v", b, err)
		}
	}
	return root, idea.Path, filepath.Join(root, protocol.DeckDir, "runs", run)
}

func TestOptionalParticipantsKeepTheDriverAlive(t *testing.T) {
	root, dir, runDir := optionalIdea(t, quota.NewOptionalParticipantPolicy(nil, nil), true)
	d := New(Config{Root: root, IdeaDir: dir, RunDir: runDir, Participants: []string{"a"}}, nil)
	if d.trackErr != nil {
		t.Fatalf("v3 sole role holder escalated: %v", d.trackErr)
	}
	root, dir, runDir = optionalIdea(t, quota.NewOptionalParticipantPolicy(nil, nil), false)
	if d = New(Config{Root: root, IdeaDir: dir, RunDir: runDir, Participants: []string{"a"}}, nil); d.trackErr == nil {
		t.Fatal("an unproved solo membership passed the non-solo track check")
	}
}
