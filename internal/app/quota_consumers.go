package app

import (
	"context"
	"parley-deck-cli/internal/membership"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/store"
	"path/filepath"
	"reflect"
)

// Every independently callable action refreshes its captured identities. This
// also covers a reduction during the same driver tick (e.g. signoff requests).
func (o driverImplOps) quotaCurrent() (driverImplOps, error) {
	ids, _, err := protocol.QuotaMembers(o.ideaDir, o.base.Idea.Participants)
	if err != nil {
		return o, err
	}
	h, err := quota.ReadHistory(o.ideaDir)
	if err != nil {
		return o, err
	}
	if h.LastRoleTransition() != nil || !reflect.DeepEqual(ids, o.base.Idea.Participants) {
		o = o.WithParticipants(ids).(driverImplOps)
	}
	return o, nil
}

func (o driverImplOps) Handover(ctx context.Context) error {
	return runner.RunRoleHandover(ctx, o.base)
}
func (o driverConsensusOps) Handover(ctx context.Context) error {
	run := membership.DrivingRunID(ctx)
	return runner.RunRoleHandover(ctx, runner.Options{Root: o.root, RunID: run, Idea: protocol.IdeaStatus{Slug: o.ideaSlug, Path: o.ideaDir, Participants: o.participants}, Agents: o.discovered, Store: store.New(filepath.Join(o.root, protocol.DeckDir, "runs", run))})
}
