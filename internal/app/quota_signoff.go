package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/fsutil"
	"parley-deck-cli/internal/runner"
	"path/filepath"
	"strconv"
	"strings"
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

// The protocol round, not mutable draft bytes or run IDs, names a signoff step.
func participantSignoffStep(path string) string {
	base := filepath.Dir(path)
	entries, _ := os.ReadDir(base)
	latest := 0
	for _, e := range entries {
		if e.IsDir() && strings.HasPrefix(e.Name(), "round-") {
			n, err := strconv.Atoi(strings.TrimPrefix(e.Name(), "round-"))
			if err == nil && n > latest {
				latest = n
			}
		}
	}
	return fmt.Sprintf("%s/round-%02d/signoff", filepath.Base(base), latest)
}

func participantSignoffOptions(root, idea string, agent agents.Discovery, path, beforeRaw string, before consensus.Summary, review bool) runner.ParticipantStepOptions {
	step := participantSignoffStep(path)
	integrity := func(raw string) error {
		if !strings.HasPrefix(raw, beforeRaw) {
			return fmt.Errorf("%s changed existing consensus prefix; preserved for repair", agent.ID)
		}
		for _, id := range signoffHeaderAgents(strings.TrimPrefix(raw, beforeRaw)) {
			if id != agent.ID {
				return fmt.Errorf("%s altered shared signoff ownership; preserved for repair", agent.ID)
			}
		}
		return nil
	}
	return runner.ParticipantStepOptions{Root: root, Idea: idea, Agent: agent, Step: step, Files: []string{path},
		Validate: func() runner.StepValidation {
			raw, err := os.ReadFile(path)
			if err != nil {
				return runner.StepValidation{Integrity: err}
			}
			if err = integrity(string(raw)); err != nil {
				return runner.StepValidation{Integrity: err}
			}
			after, err := consensus.Status(root, idea, review)
			if err != nil {
				return runner.StepValidation{Integrity: err}
			}
			if _, err = validateRequestedSignoff(before, after, agent.ID, beforeRaw, string(raw)); err != nil {
				return runner.StepValidation{Reason: err.Error()}
			}
			return runner.StepValidation{Valid: true}
		},
		AfterFailure: func() error {
			raw, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if err = integrity(string(raw)); err != nil {
				return err
			}
			if string(raw) == beforeRaw {
				return nil
			}
			// RunParticipantStep has durably copied these exact invalid bytes first.
			return fsutil.WriteFileAtomic(path, []byte(beforeRaw), 0600)
		},
	}
}

// Recover only an interrupted invalid OWN append after the common stopped-writer
// check. A live writer or shared-prefix change remains a blocking integrity error.
func recoverParticipantSignoffs(ctx context.Context, root, idea string, review bool, h *quota.History) error {
	if !h.MidIdea() || !h.Policy().Dropout() {
		return nil
	}
	dir := filepath.Join(root, protocol.DeckDir, "ideas", idea)
	path := filepath.Join(dir, "consensus.md")
	if review {
		path = filepath.Join(dir, "review", "consensus.md")
	}
	summary, err := consensus.Status(root, idea, review)
	if err != nil {
		return err
	}

	roles, err := membership.Roles(dir)
	if err != nil {
		return err
	}
	for _, id := range h.Current {
		if roles.Protected(id) {
			continue
		}
		pending, err := runner.ParticipantStepPending(root, idea, id, participantSignoffStep(path))
		if err != nil {
			return err
		}
		if !pending {
			continue
		}
		baseline, err := runner.ParticipantBaseline(root, idea, id, participantSignoffStep(path), path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		before := summary
		before.Signoffs = nil
		// Parsing the baseline does not require rewriting the shared document. Its
		// prior signoffs are the unchanged prefix subset of the current summary.
		for _, sig := range summary.Signoffs {
			if sig.Agent != id {
				before.Signoffs = append(before.Signoffs, sig)
			}
		}
		opts := participantSignoffOptions(root, idea, agents.Discovery{Spec: agents.Spec{ID: id}}, path, string(baseline), before, review)
		opts.RecoverOnly = true
		result, err := runner.RunParticipantStep(ctx, opts, nil)
		if err != nil && result.Evidence == nil {
			return err
		}
		after, e := consensus.Status(root, idea, review)
		if e != nil {
			return e
		}
		summary = after
	}
	if len(summary.Errors) == 0 {
		return nil
	}
	return errors.New("shared signoff remains malformed; no attributable own-suffix recovery")
}
