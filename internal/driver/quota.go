package driver

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"parley-deck-cli/internal/membership"
)

// Rebind interfaces cover closures captured at startup; required participants
// must never come from run.created or a pre-transition constructor after replay.
type consensusMembership interface{ WithParticipants([]string) ConsensusOps }
type implMembership interface{ WithParticipants([]string) ImplOps }
type roundMembership interface{ WithParticipants([]string) RoundRunner }

func (d *Driver) quotaBefore(ctx context.Context) (context.Context, func(), error) {
	ctx, release, err := membership.Acquire(ctx, d.cfg.IdeaDir, fmtRunID(d.cfg.RunDir))
	if err != nil {
		return ctx, nil, err
	}
	h, err := membership.Before(ctx, d.cfg.Root, d.cfg.IdeaDir, fmtRunID(d.cfg.RunDir))
	if err != nil {
		release()
		return ctx, nil, err
	}
	if h != nil && !reflect.DeepEqual(d.cfg.Participants, h.Current) {
		ids := append([]string(nil), h.Current...)
		if o, ok := d.cfg.Consensus.(consensusMembership); ok {
			d.cfg.Consensus = o.WithParticipants(ids)
		} else if d.cfg.Consensus != nil {
			release()
			return ctx, nil, fmt.Errorf("consensus consumer cannot rebind membership")
		}
		if o, ok := d.cfg.Impl.(implMembership); ok {
			d.cfg.Impl = o.WithParticipants(ids)
		} else if d.cfg.Impl != nil {
			release()
			return ctx, nil, fmt.Errorf("implementation consumer cannot rebind membership")
		}
		if r, ok := d.runner.(roundMembership); ok {
			d.runner = r.WithParticipants(ids)
		} else if d.runner != nil {
			release()
			return ctx, nil, fmt.Errorf("round consumer cannot rebind membership")
		}
		d.cfg.Participants = ids
		updated := New(d.cfg, d.runner)
		d.cfg = updated.cfg
		d.trackErr = updated.trackErr
	}
	return ctx, release, nil
}

func fmtRunID(path string) string {
	i := len(path) - 1
	for i >= 0 && !os.IsPathSeparator(path[i]) {
		i--
	}
	return path[i+1:]
}
