package trajectory

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"time"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/evidence"
)

type History struct {
	Decision               Decision `json:"decision"`
	Unreconciled           int      `json:"unreconciled"`
	RequiredReviewSequence int      `json:"required_review_sequence"`
	ReviewPending          bool     `json:"review_pending"`
	InconclusivePending    []int    `json:"inconclusive_pending"`
}

type ContinuationPreview struct {
	Version     int     `json:"version"`
	Root        string  `json:"root"`
	StateSHA256 string  `json:"state_sha256"`
	Sequence    int     `json:"sequence"`
	Source      Source  `json:"source"`
	History     History `json:"history"`
}

func (p ContinuationPreview) SHA256() string { data, _ := canonical(p); return digest(data) }

type Continuation struct {
	Preview                  ContinuationPreview `json:"preview"`
	SHA256                   string              `json:"sha256"`
	DecisionID               string              `json:"decision_id"`
	Reason                   string              `json:"reason"`
	Archive                  SnapshotRef         `json:"archive"`
	ReviewThrough            int                 `json:"review_through"`
	AcknowledgedInconclusive bool                `json:"acknowledged_inconclusive"`
	At                       time.Time           `json:"at"`
}

func continuationAt(s State, sequence int) *Continuation {
	for i := range s.Continuations {
		if s.Continuations[i].Preview.Sequence == sequence {
			return &s.Continuations[i]
		}
	}
	return nil
}

func trajectoryHistory(s State) History {
	h := History{Unreconciled: len(s.Attempts) - len(s.Resolutions)}
	ack := 0
	for _, c := range s.Continuations {
		if c.ReviewThrough > ack {
			ack = c.ReviewThrough
		}
	}
	for i, r := range s.Resolutions {
		a := r.Preview.Assessment
		h.Decision.Assessments = append(h.Decision.Assessments, a)
		switch a.Outcome {
		case Regression:
			h.Decision.Consecutive++
		case NoRegression:
			h.Decision.Consecutive = 0
		case Inconclusive:
			h.Decision.Consecutive = 0
			h.Decision.Pending = true
			c := continuationAt(s, i+1)
			if c == nil || !c.AcknowledgedInconclusive {
				h.InconclusivePending = append(h.InconclusivePending, i+1)
			}
		}
		if h.Decision.Consecutive >= 2 {
			h.RequiredReviewSequence = i + 1
			if !h.Decision.ReviewRequired {
				h.Decision.ReviewRequired = true
				h.Decision.TriggerSequence = i + 1
			}
		}
	}
	h.Decision.Pending = h.Decision.Pending || h.Unreconciled > 0
	h.ReviewPending = h.RequiredReviewSequence > ack
	return h
}

func InspectHistory(ctx context.Context, root, idea string) (History, error) {
	var h History
	found := false
	err := withState(ctx, root, idea, func(_ budget.CycleBinding, _ budget.Snapshot, s State) error {
		found = true
		h = trajectoryHistory(s)
		return nil
	})
	if err == nil && !found {
		err = errors.New("trajectory is not active")
	}
	return h, err
}

func currentTrajectorySource(s State) (Source, SnapshotRef, error) {
	if len(s.Attempts) == 0 {
		return s.Policy.Baseline, s.BaselineArchive, nil
	}
	a := s.Attempts[len(s.Attempts)-1]
	if a.After == nil || a.AfterArchive == nil {
		return Source{}, SnapshotRef{}, errors.New("charged attempt has no retained terminal source")
	}
	if c := continuationAt(s, a.Sequence); c != nil {
		return c.Preview.Source, c.Archive, nil
	}
	return *a.After, *a.AfterArchive, nil
}

func continuationPreview(ctx context.Context, root string, s State) (ContinuationPreview, error) {
	var p ContinuationPreview
	h := trajectoryHistory(s)
	if len(s.Attempts) == 0 || h.Unreconciled != 0 {
		return p, errors.New("continuation requires every charged attempt to be reconciled")
	}
	actual, err := Observe(ctx, root)
	if err != nil {
		return p, err
	}
	a := s.Attempts[len(s.Attempts)-1]
	if a.After == nil || !actual.Clean || actual.Tree.SHA256 != a.After.Tree.SHA256 {
		return p, errors.New("continuation requires a clean commit with the exact retained after-source bytes")
	}
	if _, err = gitOutput(ctx, root, "merge-base", "--is-ancestor", a.After.Tree.Commit, actual.Tree.Commit); err != nil {
		return p, errors.New("promoted source does not descend from the retained patch HEAD")
	}
	data, _ := canonical(s)
	return ContinuationPreview{1, root, digest(data), a.Sequence, actual, h}, nil
}

func PreviewContinuation(ctx context.Context, root, idea string) (ContinuationPreview, error) {
	var p ContinuationPreview
	root, err := canonicalRoot(root)
	if err != nil {
		return p, err
	}
	err = withState(ctx, root, idea, func(_ budget.CycleBinding, _ budget.Snapshot, s State) error {
		if c := continuationAt(s, len(s.Attempts)); c != nil {
			p = c.Preview
			return nil
		}
		var err error
		p, err = continuationPreview(ctx, root, s)
		return err
	})
	if err == nil && p.Version == 0 {
		err = errors.New("trajectory is not active")
	}
	return p, err
}

// Continue is reached only through the platform-attended CLI control. These
// parameters are explicit operator assertions, not model/human authentication.
// The transition records promotion and acknowledgment; it never removes a
// regression, grants budget, changes a verdict or closes an implementation.
func Continue(ctx context.Context, root, idea, expected, id, reason string, review, inconclusive bool) (ContinuationPreview, error) {
	var applied ContinuationPreview
	if !runtimeID(id) || strings.TrimSpace(reason) == "" || len(reason) > 2048 || evidence.ScrubAndTruncate(reason) != reason {
		return applied, errors.New("continuation requires a safe decision id and reason")
	}
	root, err := canonicalRoot(root)
	if err != nil {
		return applied, err
	}
	found := false
	err = withState(ctx, root, idea, func(b budget.CycleBinding, ledger budget.Snapshot, s State) error {
		found = true
		for _, c := range s.Continuations {
			if c.DecisionID == id {
				if c.SHA256 == expected && c.Reason == reason && (c.ReviewThrough > 0) == review && c.AcknowledgedInconclusive == inconclusive {
					applied = c.Preview
					return nil
				}
				return errors.New("continuation decision id was already used differently")
			}
		}
		if continuationAt(s, len(s.Attempts)) != nil {
			return errors.New("this attempt already has a retained continuation decision")
		}
		p, err := continuationPreview(ctx, root, s)
		if err != nil {
			return err
		}
		if expected != p.SHA256() {
			return errors.New("source or history changed since continuation preview")
		}
		if p.History.ReviewPending != review || (len(p.History.InconclusivePending) > 0) != inconclusive {
			return errors.New("continuation must explicitly acknowledge each pending review or inconclusive outcome")
		}
		// Only the latest unresolved result can be acknowledged here. Earlier
		// unresolved evidence could not legitimately be followed by a new charge.
		for _, seq := range p.History.InconclusivePending {
			if seq != p.Sequence {
				return errors.New("earlier unresolved outcome needs explicit recovery")
			}
		}
		archive, err := CaptureSnapshot(ctx, root, snapshotDirectory(b), p.Source)
		if err != nil {
			return err
		}
		reviewThrough := 0
		if review {
			reviewThrough = p.History.RequiredReviewSequence
		}
		s.Version = 3
		s.Continuations = append(s.Continuations, Continuation{p, expected, id, reason, archive, reviewThrough, inconclusive, time.Now().UTC()})
		if err = validateState(s, b, ledger); err != nil {
			return err
		}
		if err = writeState(statePath(b), s); err != nil {
			return err
		}
		applied = p
		return nil
	})
	if err == nil && !found {
		return applied, errors.New("trajectory is not active")
	}
	return applied, err
}

func validateTransitions(s State) error {
	if len(s.Resolutions) > len(s.Attempts) || len(s.Continuations) > len(s.Resolutions) || (s.Version == 2 && (len(s.Resolutions) != 0 || len(s.Continuations) != 0)) {
		return errors.New("invalid trajectory transition coverage or version")
	}
	for i, r := range s.Resolutions {
		p := r.Preview
		if (p.RecoverySHA256 != "" && !validHash(p.RecoverySHA256)) || p.Version != 1 || p.Sequence != i+1 || p.ChargeKey != s.Attempts[i].Charge.EntryKey || !filepath.IsAbs(p.Root) || !runtimeID(p.RunID) || !validHash(p.StateSHA256) || !validHash(p.ParentSHA256) || r.SHA256 != p.SHA256() || r.At.IsZero() || s.Attempts[i].Terminal == nil || r.At.Before(s.Attempts[i].Terminal.At) {
			return errors.New("resolution changed its charge, preview or chronology")
		}
		if p.Assessment.Outcome != Regression && p.Assessment.Outcome != NoRegression && p.Assessment.Outcome != Inconclusive {
			return errors.New("invalid reconciled outcome")
		}
	}
	last := 0
	ids := map[string]bool{}
	for _, c := range s.Continuations {
		p := c.Preview
		if p.Version != 1 || p.Sequence <= last || p.Sequence > len(s.Resolutions) || !filepath.IsAbs(p.Root) || !validHash(p.StateSHA256) || c.SHA256 != p.SHA256() || !runtimeID(c.DecisionID) || ids[c.DecisionID] || c.At.Before(s.Resolutions[p.Sequence-1].At) || !sourceValid(p.Source) || !p.Source.Clean || !c.Archive.valid() || strings.TrimSpace(c.Reason) == "" || len(c.Reason) > 2048 || evidence.ScrubAndTruncate(c.Reason) != c.Reason {
			return errors.New("continuation changed its retained decision or source")
		}
		a := s.Attempts[p.Sequence-1]
		if a.After == nil || p.Source.Tree.SHA256 != a.After.Tree.SHA256 {
			return errors.New("continuation rewrote material patch history")
		}
		prefix := s
		prefix.Attempts = prefix.Attempts[:p.Sequence]
		prefix.Resolutions = prefix.Resolutions[:p.Sequence]
		prefix.Continuations = nil
		for _, earlier := range s.Continuations {
			if earlier.Preview.Sequence < p.Sequence {
				prefix.Continuations = append(prefix.Continuations, earlier)
			}
		}
		h := trajectoryHistory(prefix)
		if !sameJSON(h, p.History) || (h.ReviewPending && (c.ReviewThrough != h.RequiredReviewSequence)) || (!h.ReviewPending && c.ReviewThrough != 0) || c.AcknowledgedInconclusive != (len(h.InconclusivePending) > 0) {
			return errors.New("continuation bypassed its original review or inconclusive gate")
		}
		ids[c.DecisionID] = true
		last = p.Sequence
	}
	return nil
}
