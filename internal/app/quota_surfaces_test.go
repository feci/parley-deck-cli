package app

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/driver"
	"parley-deck-cli/internal/membership"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/runmanifest"
	"parley-deck-cli/internal/runstate"
	"parley-deck-cli/internal/store"
)

func TestQuotaStatusWaitBriefAgreePendingAndAppliedReadOnly(t *testing.T) {
	root, dir := seedWaitIdea(t, []string{"a", "b", "c"})
	t.Setenv("PARLEY_HOME", t.TempDir())
	p := quota.NewPolicy(nil, nil)
	k := quota.NewKickoff("wait-idea", "quota-surfaces", p, []string{"a", "b", "c"}, nil, time.Now())
	if err := quota.WriteKickoff(dir, k); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "00-prompt.md")
	raw, _ := os.ReadFile(path)
	raw = bytes.Replace(raw, []byte("status:"), []byte("quota_auto_exclude: true\nquota_auto_exclude_scope: "+quota.KickoffAndMidIdea+"\nstatus:"), 1)
	os.WriteFile(path, raw, 0600)
	for _, id := range []string{"a", "b"} {
		os.WriteFile(filepath.Join(dir, "round-01", id+".md"), []byte(validRoundOne(id)), 0600)
	}
	os.WriteFile(filepath.Join(dir, "round-01", "c.md"), []byte("incomplete"), 0600)
	m := runmanifest.New(runmanifest.Options{Root: root, RunID: k.RunID, IdeaSlug: k.Idea, Participants: k.Participants, QuotaKickoff: &k})
	if err := runmanifest.Write(root, k.RunID, m); err != nil {
		t.Fatal(err)
	}
	s := store.New(filepath.Join(root, protocol.DeckDir, "runs", k.RunID))
	s.AppendDurable(store.Event{Type: "run.created", Data: map[string]any{"idea": k.Idea, "participants": k.Participants, "quota_kickoff": k}})
	s.AppendDurable(store.Event{Type: "agent.failed", Data: map[string]any{"agent": "c", "error": "quota exhaustion"}})
	now := time.Now().UTC()
	reset := now.Add(2 * time.Hour)
	ev := quota.Evidence{Eligible: true, InvocationID: "fixture-c", RuleID: "test", Provenance: "synthetic-test", ObservedAt: now, ResetAt: &reset, RawReset: "2h", Excerpt: "weekly limit exhausted"}
	members := []quota.Member{{ID: "a", Usable: true}, {ID: "b", Usable: true}, {ID: "c", Evidence: &ev}}
	h, _ := quota.ReadHistory(dir)
	d := quota.Evaluate(p, h.Current, members, quota.Roles{})
	batch := quota.NewBatch(h, k.RunID, "round-01", k.Participants, d, now)
	if err := quota.CommitBatch(dir, batch); err != nil {
		t.Fatal(err)
	}
	withWaitPoll(t, 5*time.Millisecond)
	for _, pending := range []bool{true, false} {
		if !pending {
			ctx, release, err := membership.Acquire(context.Background(), dir, k.RunID)
			if err != nil {
				t.Fatal(err)
			}
			_, err = membership.Before(ctx, root, dir, k.RunID)
			release()
			if err != nil {
				t.Fatal(err)
			}
		}
		before, _ := os.ReadFile(path)
		digest := driver.BuildPhaseDigest(root, k.Idea, dir, k.Participants)
		if strings.Join(digest.Participants, ",") != "a,b" || (digest.QuotaPending != "") != pending {
			t.Fatal(digest)
		}
		state, err := runstate.LoadRun(root, k.RunID)
		if err != nil || strings.Join(state.Participants, ",") != "a,b" || (state.QuotaPending != "") != pending {
			t.Fatal(state, err)
		}
		var out, errs bytes.Buffer
		if code := runOrganizer([]string{"brief", "--dir", root, "--idea", k.Idea}, &out, &errs); code != 0 {
			t.Fatal(code, errs.String())
		}
		for _, want := range []string{"automatic exclusion: c", ev.ResetHint()} {
			if !strings.Contains(out.String(), want) {
				t.Fatal("brief missing", want, out.String())
			}
		}
		code, waitOut, waitErr := runWaitFor(t, "--dir", root, "--idea", k.Idea, "--for", "round", "--timeout", "20ms")
		want := 0
		if pending {
			want = 3
		}
		if code != want {
			t.Fatal(code, waitOut, waitErr)
		}
		if !strings.Contains(waitOut, ev.ResetHint()) {
			t.Fatal("wait missing reset", waitOut)
		}
		for _, action := range state.NextActions {
			if action.AgentID == "c" {
				t.Fatal("planner kept excluded dispatch", action)
			}
		}
		surface := quotaSurfaceText(dir)
		if !strings.Contains(surface, ev.ResetHint()) || !strings.Contains(surface, "automatic exclusion: c") {
			t.Fatal(surface)
		}
		after, _ := os.ReadFile(path)
		if !bytes.Equal(before, after) {
			t.Fatal("read surface repaired prompt")
		}
	}
}
