package membership

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"parley-deck-cli/internal/pidlease"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/runmanifest"
)

// RecordManual imports an already recorded knob-off confirmation at a mutation
// boundary. Read surfaces can show that confirmation, but never write history.
func RecordManual(root, dir string) error {
	v, err := protocol.InspectQuota(dir)
	if err != nil {
		return err
	}
	if v.Manual == nil {
		return nil
	}
	if err = revisionOpen(dir); err != nil {
		return err
	}
	path, err := lockPath(dir, "revision")
	if err != nil {
		return err
	}
	l, err := pidlease.TryAcquire(path, filepath.Base(dir)+"/manual-revision")
	if err != nil {
		return err
	}
	defer l.Release()
	v, err = protocol.InspectQuota(dir)
	if err != nil {
		return err
	}
	if v.Manual == nil {
		return nil
	}
	if v.Pending != "" {
		return fmt.Errorf("%s", v.Pending)
	}
	if err = RequireStopped(root, v.History.Kickoff.Idea); err != nil {
		return err
	}
	// Retain historical vetoes/findings for a manual exclusion too.
	b := *v.Manual
	removed := []string{}
	for _, id := range b.Decision.Before {
		if !Has(b.Decision.After, id) {
			removed = append(removed, id)
		}
	}
	obs, err := protocol.CaptureQuotaObligations(dir, removed)
	if err != nil {
		return err
	}
	b = b.BindRetained(obs)
	if _, err = validateManifests(root, b.RunID, v.History); err != nil {
		return err
	}
	if err = quota.CommitBatch(dir, b); err != nil {
		return err
	}
	h, err := quota.ReadHistory(dir)
	if err != nil {
		return err
	}
	return reconcileLocked(context.Background(), root, dir, b.RunID, h)
}

// Revise requires a committed owner answer that spells out the exact requested
// set and policy. This is the explicit same-idea return/scope path; never called
// by binary upgrade, reset timers or automatic batch exclusion.
func Revise(ctx context.Context, root, dir, run string, ids []string, p quota.Policy, a quota.Authority, catchup map[string]string) (*quota.Batch, error) {
	if err := revisionOpen(dir); err != nil {
		return nil, err
	}
	v, err := protocol.InspectQuota(dir)
	if err != nil {
		return nil, err
	}
	h := v.History
	if h == nil {
		return nil, fmt.Errorf("owner revision requires recorded kickoff history")
	}
	if v.Pending != "" || v.Manual != nil {
		return nil, fmt.Errorf("reconcile recorded membership before owner revision")
	}
	m, e := runmanifest.Load(root, run)
	if e != nil {
		return nil, e
	}
	if m.RunID != run || m.IdeaSlug != h.Kickoff.Idea || m.QuotaRevision != h.Revision || !reflect.DeepEqual(m.QuotaKickoff, h.Kickoff) || !reflect.DeepEqual(m.Participants, h.Current) {
		return nil, fmt.Errorf("owner revision requires current idea-bound manifest")
	}
	if err = quota.ValidateAuthority(root, h.Kickoff.Idea, a); err != nil {
		return nil, err
	}
	// Serialize the explicit revision even when switching from a policy-off or
	// kickoff-only idea. Normal driving runs in those modes remain unrestricted.
	path, err := lockPath(dir, "revision")
	if err != nil {
		return nil, err
	}
	l, err := pidlease.TryAcquire(path, filepath.Base(dir)+"/owner-revision")
	if err != nil {
		return nil, err
	}
	defer l.Release()
	ctx, release, err := Acquire(ctx, dir, run)
	if err != nil {
		return nil, err
	}
	defer release()
	if err = RequireStopped(root, h.Kickoff.Idea); err != nil {
		return nil, err
	}
	// Acquire the future lifetime lease before publishing an enabling revision.
	if !h.MidIdea() && p.Enabled && p.Scope == quota.KickoffAndMidIdea {
		ctx, release, err = acquireScoped(ctx, dir, run)
		if err != nil {
			return nil, err
		}
		defer release()
	}
	unlock, err := projectionLock(dir)
	if err != nil {
		return nil, err
	}
	defer unlock()
	h, err = quota.ReadHistory(dir)
	if err != nil {
		return nil, err
	}
	proof := map[string]string{}
	for _, id := range ids {
		if Has(h.Known, id) {
			continue
		}
		rel := catchup[id]
		if rel != filepath.ToSlash(filepath.Join("round-01", id+".md")) {
			return nil, fmt.Errorf("new identity %s needs canonical catch-up artifact", id)
		}
		path := filepath.Join(dir, filepath.FromSlash(rel))
		if err = protocol.ValidateQuotaCatchup(path, h.Kickoff.Idea, id); err != nil {
			return nil, err
		}
		raw, e := os.ReadFile(path)
		if e != nil {
			return nil, e
		}
		proof[id] = string(raw)
	}
	b := quota.NewRevision(h, run, ids, p, quota.Revision{Authority: &a, Catchup: proof}, time.Now())
	removed := []string{}
	for _, id := range h.Current {
		if !Has(ids, id) {
			removed = append(removed, id)
		}
	}
	obs, err := protocol.CaptureQuotaObligations(dir, removed)
	if err != nil {
		return nil, err
	}
	b = b.BindRetained(obs)
	if _, err = validateManifests(root, run, h); err != nil {
		return nil, err
	}
	if err = quota.CommitBatch(dir, b); err != nil {
		return nil, err
	}
	h, err = quota.ReadHistory(dir)
	if err != nil {
		return &b, err
	}
	if err = reconcileLocked(ctx, root, dir, run, h); err != nil {
		return &b, err
	}
	return &b, nil
}

func revisionOpen(dir string) error {
	m, e := protocol.ReadFrontmatter(filepath.Join(dir, "00-prompt.md"))
	if e != nil {
		return e
	}
	if m["status"] == "final" || m["status"] == "closed" {
		return fmt.Errorf("closed idea membership is frozen; use the existing reopen protocol first")
	}
	return nil
}
