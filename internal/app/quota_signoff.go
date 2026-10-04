package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"parley-deck-cli/internal/membership"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/runmanifest"
	"parley-deck-cli/internal/store"
)

func quotaSignoffStart(ctx context.Context, root, slug string, dry bool) (context.Context, string, func(), *quota.History, error) {
	dir := filepath.Join(root, protocol.DeckDir, "ideas", slug)
	v, err := protocol.InspectQuota(dir)
	if err != nil {
		return ctx, "", nil, nil, err
	}
	run := membership.DrivingRunID(ctx)
	if run == "" {
		run = store.NewRunID(time.Now())
	}
	if dry || !v.History.MidIdea() {
		return ctx, run, func() {}, v.History, nil
	}
	ctx, release, err := membership.Acquire(ctx, dir, run)
	if err != nil {
		return ctx, run, nil, nil, err
	}
	if membership.DrivingRunID(ctx) != run {
		release()
		return ctx, run, nil, nil, fmt.Errorf("quota signoff run identity mismatch")
	}
	if _, e := runmanifest.Load(root, run); e != nil {
		if !os.IsNotExist(e) {
			release()
			return ctx, run, nil, nil, e
		}
		h := v.History
		m := runmanifest.New(runmanifest.Options{Root: root, RunID: run, IdeaSlug: slug, Mode: "consensus-signoff", Participants: h.Current, QuotaKickoff: h.Kickoff})
		m.QuotaRevision = h.Revision
		if err = runmanifest.Write(root, run, m); err != nil {
			release()
			return ctx, run, nil, nil, err
		}
		if err = store.New(filepath.Join(root, protocol.DeckDir, "runs", run)).AppendDurable(store.Event{Type: "run.created", Data: map[string]any{"idea": slug, "participants": h.Current, "quota_kickoff": h.Kickoff, "mode": "consensus-signoff"}}); err != nil {
			release()
			return ctx, run, nil, nil, err
		}
	}
	h, err := membership.Before(ctx, root, dir, run)
	if err != nil {
		release()
		return ctx, run, nil, nil, err
	}
	return ctx, run, release, h, nil
}
