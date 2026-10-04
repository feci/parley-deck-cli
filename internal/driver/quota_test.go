package driver

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/membership"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/runmanifest"
	"parley-deck-cli/internal/store"
)

type quotaConsensusConsumer struct {
	fakeConsensus
	ids []string
}

func (o *quotaConsensusConsumer) WithParticipants(ids []string) ConsensusOps {
	o.ids = append([]string(nil), ids...)
	return o
}

type quotaImplConsumer struct {
	fakeImpl
	ids []string
}

func (o *quotaImplConsumer) WithParticipants(ids []string) ImplOps {
	o.ids = append([]string(nil), ids...)
	return o
}

type quotaRoundConsumer struct {
	ids        []string
	dispatched []string
}

func (o *quotaRoundConsumer) WithParticipants(ids []string) RoundRunner {
	o.ids = append([]string(nil), ids...)
	return o
}
func (o *quotaRoundConsumer) RunRound(context.Context, int) error {
	o.dispatched = append([]string(nil), o.ids...)
	return nil
}

func TestQuotaDriverReconcilesPendingBeforeRebindingEveryConsumer(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PARLEY_HOME", t.TempDir())
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	p := quota.NewPolicy(nil, nil)
	ids := []string{"a", "b", "c"}
	idea, k, err := protocol.CreateIdeaWithQuota(root, "quota driver", ids, nil, "deliberation", "", "old-run", &p, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, run := range []string{"old-run", "new-run"} {
		m := runmanifest.New(runmanifest.Options{Root: root, RunID: run, IdeaSlug: idea.Slug, Participants: ids, QuotaKickoff: k})
		if err = runmanifest.Write(root, run, m); err != nil {
			t.Fatal(err)
		}
		if err = store.New(filepath.Join(root, protocol.DeckDir, "runs", run)).AppendDurable(store.Event{Type: "run.created", Data: map[string]any{"idea": idea.Slug, "participants": ids, "quota_kickoff": k}}); err != nil {
			t.Fatal(err)
		}
	}
	now := time.Now()
	reset := now.Add(2 * time.Hour)
	e := quota.Evidence{Eligible: true, InvocationID: "c", Provenance: "synthetic-test", RuleID: "test", Excerpt: "weekly limit exhausted", ObservedAt: now, ResetAt: &reset}
	h, _ := quota.ReadHistory(idea.Path)
	decision := quota.Evaluate(p, ids, []quota.Member{{ID: "a", Usable: true}, {ID: "b", Usable: true}, {ID: "c", Evidence: &e}}, quota.Roles{})
	batch := quota.NewBatch(h, "old-run", "consensus", ids, decision, now)
	if err = quota.CommitBatch(idea.Path, batch); err != nil {
		t.Fatal(err)
	}
	c := &quotaConsensusConsumer{ids: ids}
	r := &quotaRoundConsumer{ids: ids}
	impl := &quotaImplConsumer{ids: ids}
	d := New(Config{Root: root, IdeaDir: idea.Path, IdeaSlug: idea.Slug, RunDir: filepath.Join(root, protocol.DeckDir, "runs", "new-run"), Participants: ids, Consensus: c, Impl: impl}, r)
	ctx, release, err := d.quotaBefore(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer release()
	for _, got := range [][]string{d.cfg.Participants, c.ids, r.ids, impl.ids} {
		if strings.Join(got, ",") != "a,b" {
			t.Fatal("stale consumer", got)
		}
	}
	if err = r.RunRound(ctx, 2); err != nil {
		t.Fatal(err)
	}
	if strings.Join(r.dispatched, ",") != "a,b" {
		t.Fatal(r.dispatched)
	}
	if _, _, err = membership.Acquire(context.Background(), idea.Path, "competitor"); err == nil {
		t.Fatal("lifetime lease released before dispatch")
	}
	for _, run := range []string{"old-run", "new-run"} {
		m, err := runmanifest.Load(root, run)
		if err != nil || m.QuotaRevision != 1 || strings.Join(m.Participants, ",") != "a,b" {
			t.Fatal(m, err)
		}
	}
	// Removing authority must gate a later tick rather than recreating a kickoff.
	if err = os.Remove(filepath.Join(idea.Path, quota.KickoffFile)); err != nil {
		t.Fatal(err)
	}
	if _, _, err = d.quotaBefore(ctx); err == nil {
		t.Fatal("missing history repaired")
	}
}
