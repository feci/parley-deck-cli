package membership

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/quotatest"
	"parley-deck-cli/internal/runmanifest"
	"parley-deck-cli/internal/store"
	"parley-deck-cli/internal/telemetry"
)

func TestDropoutCrashDerivedProofNeverInventsTerminal(t *testing.T) {
	oldDead, oldGroup := crashDeadLocal, crashGroupStopped
	t.Cleanup(func() { crashDeadLocal, crashGroupStopped = oldDead, oldGroup })
	dead := false
	crashDeadLocal = func(string, string, int) bool { return dead }
	crashGroupStopped = func(int) bool { return dead }
	root := t.TempDir()
	now := time.Now().Add(-time.Minute)
	pid := 12345
	requested := telemetry.Record{Type: "invocation.requested", InvocationID: "orphan", RequestedAt: now, Metadata: telemetry.Metadata{Idea: "idea", Agent: "c", AttemptOrdinal: 1}}
	started := requested
	started.Type = "invocation.started"
	started.StartedAt = &now
	started.PID = &pid
	started.Writer = &telemetry.WriterIdentity{Host: "fixture", Boot: "boot", SupervisorPID: 12346, ProcessGroup: pid}
	dir := filepath.Join(root, ".parley-runtime", "invocations", "orphan")
	os.MkdirAll(dir, 0700)
	raw, _ := json.Marshal(started)
	os.WriteFile(filepath.Join(dir, "started.json"), raw, 0600)
	if _, err := ParticipantCrash(root, requested); err == nil {
		t.Fatal("unresolved writer admitted")
	}
	dead = true
	derived, err := ParticipantCrash(root, requested)
	if err != nil {
		t.Fatal(err)
	}
	if derived.Type != "invocation.crash-settled" || derived.Outcome == nil || *derived.Outcome.FailureClass != "crash" {
		t.Fatalf("%+v", derived)
	}
	if _, err = os.Stat(filepath.Join(dir, "terminal.json")); !os.IsNotExist(err) {
		t.Fatal("fake terminal created", err)
	}
	again, err := ParticipantCrash(root, requested)
	if err != nil || !again.CompletedAt.Equal(*derived.CompletedAt) {
		t.Fatal(again, err)
	}
	requested.Metadata.Agent = "wrong"
	if _, err := ParticipantCrash(root, requested); err == nil {
		t.Fatal("mismatched request admitted")
	}
}

func TestDropoutBlockingDecisionDeduplicatesAcrossRuns(t *testing.T) {
	root, dir, _ := fixture(t, quota.NewParticipantPolicy(nil, nil))
	d := quota.Evaluate(quota.NewParticipantPolicy(nil, nil), []string{"a", "b", "c", "d"}, dropoutMembers(t, filepath.Base(dir)), quota.Roles{})
	d.Applied = false
	d.Block = "reviewer gate"
	for _, run := range []string{"first", "resumed"} {
		if err := Block(root, dir, run, "round-01", d); !IsBlocked(err) {
			t.Fatal(err)
		}
	}
	paths, _ := filepath.Glob(filepath.Join(root, protocol.DeckDir, "inbox", "parley-to-user_quota-block-*.md"))
	if len(paths) != 1 {
		t.Fatal(paths)
	}
}

func dropoutMembers(t *testing.T, idea string) []quota.Member {
	t.Helper()
	ms := []quota.Member{{ID: "a", Usable: true}, {ID: "b", Usable: true}}
	for _, id := range []string{"c", "d"} {
		now := time.Now().UTC()
		code := 1
		first := quota.FailedAttempt{InvocationID: id + "-1", Status: "failed", FailureClass: "process_failure", ExitCode: &code, ObservedAt: now, Excerpt: "exit_code: 1"}
		second := first
		second.InvocationID = id + "-2"
		second.RetryOf = first.InvocationID
		second.ObservedAt = now.Add(5 * time.Second)
		e, err := telemetry.ParticipantEvidence(idea, id, "round-01/"+id+".md", "fixture", []quota.FailedAttempt{first, second})
		if err != nil {
			t.Fatal(err)
		}
		ms = append(ms, quota.Member{ID: id, Evidence: &e})
	}
	return ms
}

func TestDropoutEveryReturnPathAfterDowngrade(t *testing.T) {
	root, dir, run := fixture(t, quota.NewParticipantPolicy(nil, nil))
	ctx := leaseFixture(t, dir, run)
	if _, err := Settle(ctx, root, dir, run, "round-01", []string{"a", "b", "c", "d"}, dropoutMembers(t, filepath.Base(dir))); err != nil {
		t.Fatal(err)
	}
	off := quota.Policy{Enabled: false, Scope: quota.KickoffAndMidIdea}
	authority := quotatest.Authority(t, root, filepath.Base(dir), "disable", (quota.RevisionDirective{Participants: []string{"a", "b"}, Policy: off}).Text())
	if _, err := Revise(ctx, root, dir, run, []string{"a", "b"}, off, authority, nil); err != nil {
		t.Fatal(err)
	}
	rejoin := quotatest.Authority(t, root, filepath.Base(dir), "return", (quota.RevisionDirective{Participants: []string{"a", "b", "c"}, Policy: off}).Text())
	if _, err := Revise(ctx, root, dir, run, []string{"a", "b", "c"}, off, rejoin, nil); err == nil || !strings.Contains(err.Error(), "permanently dropped") {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "00-prompt.md")
	raw, _ := os.ReadFile(path)
	changed := strings.Replace(string(raw), "participants: [a, b]", "participants: [a, b, c]", 1)
	if changed == string(raw) {
		t.Fatal("fixture edit missed")
	}
	os.WriteFile(path, []byte(changed), 0600)
	if _, err := protocol.InspectQuota(dir); err == nil || !strings.Contains(err.Error(), "permanently dropped") {
		t.Fatal(err)
	}
	if _, err := Before(ctx, root, dir, run); err == nil {
		t.Fatal("recovery erased permanent return ban")
	}
	os.WriteFile(path, raw, 0600)
	v, err := protocol.InspectQuota(dir)
	if err != nil {
		t.Fatal(err)
	}
	if protocol.ManualCatchupTarget(dir, v, "c", filepath.Join(dir, "round-01", "c.md")) {
		t.Fatal("catch-up admitted permanent drop")
	}
}

func TestDropoutKickoffNoticeReceiptReplay(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PARLEY_HOME", t.TempDir())
	protocol.InitWorkspace(root)
	p := quota.NewParticipantPolicy(nil, nil)
	ms := dropoutMembers(t, "readiness-batch")
	for i := range ms {
		if ms[i].Evidence != nil {
			ms[i].Evidence.Failure.Step = "readiness"
		}
	}
	d := quota.Evaluate(p, []string{"a", "b", "c", "d"}, ms, quota.Roles{})
	idea, k, err := protocol.CreateIdeaWithQuota(root, "kickoff replay", d.After, nil, "deliberation", "", "kickoff-run", &p, &d)
	if err != nil {
		t.Fatal(err)
	}
	runmanifest.Write(root, k.RunID, runmanifest.New(runmanifest.Options{Root: root, RunID: k.RunID, IdeaSlug: idea.Slug, Participants: k.Participants, QuotaKickoff: k}))
	store.New(filepath.Join(root, protocol.DeckDir, "runs", k.RunID)).AppendDurable(store.Event{Type: "run.created", Data: map[string]any{"idea": idea.Slug, "participants": k.Participants, "quota_kickoff": k}})
	// Crash point: kickoff+manifest exist, publication/receipt have not run.
	ctx, release, err := Acquire(context.Background(), idea.Path, k.RunID)
	if err != nil {
		t.Fatal(err)
	}
	defer release()
	if _, err = Before(ctx, root, idea.Path, k.RunID); err != nil {
		t.Fatal(err)
	}
	id := k.Transition.ID
	notice := filepath.Join(root, protocol.DeckDir, "inbox", "parley-to-user_"+id+".md")
	if _, err = os.Stat(notice); err != nil {
		t.Fatal(err)
	}
	if ok, err := quota.ReadApplied(idea.Path, id); err != nil || !ok {
		t.Fatal(ok, err)
	}
	archive := filepath.Join(filepath.Dir(notice), "archived")
	os.MkdirAll(archive, 0700)
	os.Rename(notice, filepath.Join(archive, filepath.Base(notice)))
	os.WriteFile(filepath.Join(archive, filepath.Base(notice)), []byte("owner edited"), 0600)
	// An interruption after publication but before its receipt preserves archived edits.
	os.Remove(filepath.Join(idea.Path, "quota-applied", id))
	if _, err = Before(ctx, root, idea.Path, k.RunID); err != nil {
		t.Fatal(err)
	}
	if _, err = os.Stat(notice); !os.IsNotExist(err) {
		t.Fatal("notice re-created after archival", err)
	}
	h, err := quota.ReadHistory(idea.Path)
	if err != nil || h.CheckReturn([]string{"c"}) == nil {
		t.Fatal(h, err)
	}
}
