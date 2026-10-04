package app

import (
	"context"
	"path/filepath"
	"time"

	"parley-deck-cli/internal/membership"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runmanifest"
	"parley-deck-cli/internal/store"
)

// A pipeline block owns one run identity across implementation, review, signoff
// and fix-up, including the intervals between calls to the individual runners.
func quotaPipelineStart(ctx context.Context, root, dir string) (context.Context, func(), error) {
	v, err := protocol.InspectQuota(dir)
	if err != nil {
		return ctx, nil, err
	}
	if !v.History.MidIdea() {
		return ctx, func() {}, nil
	}
	run := "pipe-drive-" + store.NewRunID(time.Now())
	ctx, release, err := membership.Acquire(ctx, dir, run)
	if err != nil {
		return ctx, nil, err
	}
	h := v.History
	m := runmanifest.New(runmanifest.Options{Root: root, RunID: run, IdeaSlug: h.Kickoff.Idea, Participants: h.Current, QuotaKickoff: h.Kickoff})
	m.QuotaRevision = h.Revision
	if err = runmanifest.Write(root, run, m); err != nil {
		release()
		return ctx, nil, err
	}
	if err = store.New(filepath.Join(root, protocol.DeckDir, "runs", run)).AppendDurable(store.Event{Type: "run.created", Data: map[string]any{"idea": h.Kickoff.Idea, "participants": h.Current, "quota_kickoff": h.Kickoff, "mode": "pipeline"}}); err != nil {
		release()
		return ctx, nil, err
	}
	if _, err = membership.Before(ctx, root, dir, run); err != nil {
		release()
		return ctx, nil, err
	}
	return ctx, release, nil
}
func quotaPipelineRun(ctx context.Context, fallback string) string {
	if id := membership.DrivingRunID(ctx); id != "" {
		return id
	}
	return fallback
}
