package trajectory

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/evidence"
	"parley-deck-cli/internal/fsutil"
)

// Source retains actual worktree observations without changing Git history.
// A dirty or unavailable post-state is evidence of unfinished work, never a
// substituted clean snapshot or permission to discard the patch.
type Source struct {
	Tree         Tree   `json:"tree"`
	StatusSHA256 string `json:"status_sha256"`
	Clean        bool   `json:"clean"`
}
type Policy struct {
	Version     int                `json:"version"`
	Idea        string             `json:"idea"`
	Scope       string             `json:"scope"`
	StartedAt   time.Time          `json:"started_at"`
	Implementer string             `json:"implementer"`
	Baseline    Source             `json:"baseline"`
	Criteria    []CriterionBinding `json:"criteria"`
}
type Launch struct {
	InvocationID string `json:"invocation_id"`
	Agent        string `json:"agent"`
}
type Terminal struct {
	At            time.Time `json:"at"`
	Status        string    `json:"status"`
	ExitCode      *int      `json:"exit_code"`
	SnapshotError string    `json:"snapshot_error"`
}
type Attempt struct {
	ReservationIntentSHA256 string             `json:"reservation_intent_sha256,omitempty"`
	Sequence                int                `json:"sequence"`
	Charge                  budget.CycleCharge `json:"charge"`
	Before                  Source             `json:"before"`
	BeforeArchive           SnapshotRef        `json:"before_archive"`
	Launch                  *Launch            `json:"launch"`
	After                   *Source            `json:"after"`
	AfterArchive            *SnapshotRef       `json:"after_archive"`
	Terminal                *Terminal          `json:"terminal"`
}
type State struct {
	Version         int            `json:"version"`
	Policy          Policy         `json:"policy"`
	BaselineArchive SnapshotRef    `json:"baseline_archive"`
	Attempts        []Attempt      `json:"attempts"`
	Resolutions     []Resolution   `json:"resolutions,omitempty"`
	Continuations   []Continuation `json:"continuations,omitempty"`
}

func statePath(b budget.CycleBinding) string {
	return filepath.Join(filepath.Dir(b.Store.Dir), "trajectory.json")
}
func sourceValid(s Source) bool {
	return validCommit(s.Tree.Commit) && validHash(s.Tree.SHA256) && validHash(s.StatusSHA256) && s.Clean == (s.StatusSHA256 == digest(nil))
}
func (p Policy) SHA256() (string, error) {
	if err := trajectoryVersion(p.Version); err != nil {
		return "", err
	}
	if p.Version != 2 {
		return "", errors.New("unsupported frozen trajectory policy version")
	}
	if !safeLabel(p.Idea) || p.Scope == "" || p.StartedAt.IsZero() || !safeLabel(p.Implementer) || !sourceValid(p.Baseline) || !p.Baseline.Clean || len(p.Criteria) == 0 || len(p.Criteria) > MaxCriteria {
		return "", errors.New("invalid frozen trajectory policy")
	}
	seen := map[string]bool{}
	for _, c := range p.Criteria {
		if !safeLabel(c.Name) || !validHash(c.CommandSHA256) || seen[c.Name] {
			return "", errors.New("invalid trajectory material scope")
		}
		seen[c.Name] = true
	}
	data, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return digest(data), nil
}

func trajectoryVersion(version int) error {
	if version == 1 {
		return errors.New("trajectory v1 retained digests only; preserve its history for explicit recovery; current files cannot reconstruct past attempts")
	}
	if version != 2 && version != 3 {
		return errors.New("unsupported trajectory version")
	}
	return nil
}

// Observe brackets the source digest with Git observations. It accepts dirty
// results so failed work remains visible. It persists no diff/status bodies.
func Observe(ctx context.Context, root string) (Source, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return Source{}, err
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return Source{}, err
	}
	top, err := gitOutput(ctx, root, "rev-parse", "--show-toplevel")
	if err != nil {
		return Source{}, err
	}
	actual, err := filepath.EvalSymlinks(strings.TrimSpace(string(top)))
	if err != nil || actual != root {
		return Source{}, errors.New("trajectory source must be an exact Git worktree root")
	}
	read := func() (string, []byte, error) {
		commit, err := gitOutput(ctx, root, "rev-parse", "--verify", "HEAD^{commit}")
		if err != nil {
			return "", nil, err
		}
		status, err := gitOutput(ctx, root, "status", "--porcelain=v1", "--untracked-files=all")
		return strings.TrimSpace(string(commit)), status, err
	}
	head, status, err := read()
	if err != nil {
		return Source{}, err
	}
	tree, err := trajectoryTreeDigest(ctx, root)
	if err != nil {
		if errors.Is(err, evidence.ErrLocalSourceExcludes) {
			return Source{}, evidence.ErrLocalSourceExcludes
		}
		return Source{}, errors.New("trajectory cannot digest the actual source")
	}
	head2, status2, err := read()
	if err != nil {
		return Source{}, err
	}
	tree2, err := trajectoryTreeDigest(ctx, root)
	if err != nil {
		if errors.Is(err, evidence.ErrLocalSourceExcludes) {
			return Source{}, evidence.ErrLocalSourceExcludes
		}
		return Source{}, errors.New("trajectory cannot recheck the actual source")
	}
	if head != head2 || !bytes.Equal(status, status2) || tree != tree2 {
		return Source{}, errors.New("trajectory source changed during observation")
	}
	s := Source{Tree{head, tree}, digest(status), len(status) == 0}
	if !sourceValid(s) {
		return Source{}, errors.New("invalid actual source identity")
	}
	return s, nil
}
func NewPolicy(ctx context.Context, root, idea, implementer string, criteria []Criterion) (Policy, string, error) {
	b, err := budget.LoadCycleBinding(ctx, root, idea, budget.Fixup)
	if err != nil {
		return Policy{}, "", err
	}
	if b == nil {
		return Policy{}, "", errors.New("initialize the original fixup budget before opting in")
	}
	ledger, err := b.Store.Inspect(ctx)
	if err != nil {
		return Policy{}, "", err
	}
	source, err := Observe(ctx, root)
	if err != nil {
		return Policy{}, "", err
	}
	p := Policy{Version: 2, Idea: idea, Scope: b.Policy.Scope, StartedAt: ledger.StartedAt, Implementer: implementer, Baseline: source}
	for _, c := range criteria {
		if strings.TrimSpace(c.Command) == "" || len(c.Command) > 16<<10 || strings.ContainsRune(c.Command, 0) {
			return Policy{}, "", errors.New("invalid material command")
		}
		p.Criteria = append(p.Criteria, CriterionBinding{c.Name, digest([]byte(c.Command))})
	}
	_, err = p.SHA256()
	original := b.Policy
	original.TrajectorySHA256 = ""
	return p, budget.CyclePolicyDigest(original), err
}
func canonical(v any) ([]byte, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	return append(data, '\n'), err
}
func writeState(path string, s State) error {
	data, err := canonical(s)
	if err != nil {
		return err
	}
	if len(data) > 16<<20 {
		return errors.New("trajectory state exceeds its bound")
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".trajectory-*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err = f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err = fsutil.SyncFile(f); err != nil {
		f.Close()
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return fsutil.ReplaceSyncedFile(f.Name(), path)
}
func readState(path string) (State, []byte, error) {
	var s State
	info, err := os.Lstat(path)
	if err != nil {
		return s, nil, err
	}
	if !info.Mode().IsRegular() || info.Size() > 16<<20 {
		return s, nil, errors.New("trajectory state must be a bounded regular file")
	}
	f, err := os.Open(path)
	if err != nil {
		return s, nil, err
	}
	defer f.Close()
	opened, err := f.Stat()
	if err != nil || !os.SameFile(info, opened) {
		return s, nil, errors.New("trajectory state changed during open")
	}
	data, err := io.ReadAll(io.LimitReader(f, (16<<20)+1))
	if err != nil {
		return s, nil, err
	}
	if len(data) > 16<<20 {
		return s, nil, errors.New("trajectory state exceeds its bound")
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err = dec.Decode(&s); err != nil {
		return s, nil, err
	}
	if err = trajectoryVersion(s.Version); err != nil {
		return s, nil, err
	}
	expected, err := canonical(s)
	if err != nil {
		return s, nil, err
	}
	if !bytes.Equal(data, expected) {
		return s, nil, errors.New("trajectory state is incomplete, ambiguous or noncanonical")
	}
	return s, data, nil
}
func validateState(s State, b budget.CycleBinding, ledger budget.Snapshot) error {
	sha, err := s.Policy.SHA256()
	if err != nil {
		return err
	}
	if (s.Version != 2 && s.Version != 3) || !s.BaselineArchive.valid() || s.Attempts == nil || len(s.Attempts) > MaxPatches || sha != b.Policy.TrajectorySHA256 || s.Policy.Scope != b.Policy.Scope || s.Policy.Idea != b.Policy.Idea || !s.Policy.StartedAt.Equal(ledger.StartedAt) || b.Policy.Carried != 0 || len(s.Attempts) != len(ledger.Entries) {
		return errors.New("trajectory policy or complete charged history is missing or changed")
	}
	if err := validateTransitions(s); err != nil {
		return err
	}
	seen := map[string]bool{}
	before := s.Policy.Baseline
	beforeArchive := s.BaselineArchive
	for i, a := range s.Attempts {
		if a.ReservationIntentSHA256 != "" && !validHash(a.ReservationIntentSHA256) {
			return errors.New("invalid original reservation intent reference")
		}
		if i > 0 {
			if len(s.Resolutions) < i {
				return errors.New("an unverified trajectory attempt was followed by another charge")
			}
			prefix := s
			prefix.Attempts = prefix.Attempts[:i]
			prefix.Resolutions = prefix.Resolutions[:i]
			prefix.Continuations = nil
			for _, c := range s.Continuations {
				if c.Preview.Sequence <= i {
					prefix.Continuations = append(prefix.Continuations, c)
				}
			}
			h := trajectoryHistory(prefix)
			if h.ReviewPending || len(h.InconclusivePending) > 0 {
				return errors.New("charged history bypassed a trajectory review gate")
			}
			if c := continuationAt(s, i); c != nil {
				before, beforeArchive = c.Preview.Source, c.Archive
			}
		}
		if a.Sequence != i+1 || seen[a.Charge.EntryKey] || a.Charge.Kind != budget.Fixup || a.Before != before || a.BeforeArchive != beforeArchive || !a.BeforeArchive.valid() || !sourceValid(a.Before) || !a.Before.Clean {
			return errors.New("trajectory has a gap, replay or changed before-state")
		}
		seen[a.Charge.EntryKey] = true
		if err := a.Charge.Check(ledger); err != nil {
			return err
		}
		if a.Launch != nil && (!safeLabel(a.Launch.InvocationID) || a.Launch.Agent != s.Policy.Implementer) {
			return errors.New("trajectory launch identity changed")
		}
		if a.Terminal == nil {
			if a.After != nil || a.AfterArchive != nil {
				return errors.New("trajectory has a post-state without a terminal observation")
			}
		} else {
			if a.Launch == nil || a.Terminal.At.Before(a.Charge.ReservedAt) || (a.Terminal.Status != "process-exited" && a.Terminal.Status != "failed") || (a.Terminal.ExitCode != nil && *a.Terminal.ExitCode < -1) {
				return errors.New("invalid trajectory terminal observation")
			}
			if a.After == nil {
				if a.Terminal.SnapshotError != "source-unavailable" || a.AfterArchive != nil {
					return errors.New("missing actual post-state")
				}
			} else if !sourceValid(*a.After) {
				return errors.New("invalid actual post-state")
			} else if a.AfterArchive == nil {
				if a.Terminal.SnapshotError != "archive-unavailable" {
					return errors.New("post-state lacks its required archive")
				}
			} else if !a.AfterArchive.valid() || a.Terminal.SnapshotError != "" {
				return errors.New("invalid post-state archive")
			}
		}
		if a.After != nil {
			before = *a.After
			if a.AfterArchive != nil {
				beforeArchive = *a.AfterArchive
			}
		}
	}
	return nil
}
func Activate(ctx context.Context, root, expected string, p Policy) error {
	sha, err := p.SHA256()
	if err != nil {
		return err
	}
	return budget.ActivateCycleTrajectory(ctx, root, p.Idea, expected, sha, func(b budget.CycleBinding, ledger budget.Snapshot) error {
		if p.Scope != b.Policy.Scope || !p.StartedAt.Equal(ledger.StartedAt) {
			return errors.New("trajectory policy differs from the accounting origin")
		}
		actual, err := Observe(ctx, root)
		if err != nil {
			return err
		}
		if actual != p.Baseline {
			return errors.New("source changed since trajectory preview")
		}
		archive, err := CaptureSnapshot(ctx, root, snapshotDirectory(b), p.Baseline)
		if err != nil {
			return err
		}
		s := State{Version: 2, Policy: p, BaselineArchive: archive, Attempts: []Attempt{}}
		expectedBytes, err := canonical(s)
		if err != nil {
			return err
		}
		_, old, err := readState(statePath(b))
		if err == nil {
			if !bytes.Equal(old, expectedBytes) {
				return errors.New("conflicting trajectory initialization; preserve existing history")
			}
			return nil
		}
		if !os.IsNotExist(err) {
			return err
		}
		return writeState(statePath(b), s)
	})
}

// Observer is one live reservation's snapshot capture. Its methods are called
// by budget only under the existing common cycle resource guard.
type Observer struct {
	intentSHA256 string
	Root         string
	prepared     []byte
	before       Source
	archive      SnapshotRef
}

func (o *Observer) BeforeCycle(ctx context.Context, b budget.CycleBinding, ledger budget.Snapshot) error {
	o.prepared = nil
	s, data, err := readState(statePath(b))
	if err != nil {
		return err
	}
	if err = validateState(s, b, ledger); err != nil {
		return err
	}
	if err = checkStateSnapshots(ctx, b, s); err != nil {
		return err
	}
	h := trajectoryHistory(s)
	if h.Unreconciled > 0 || h.ReviewPending || len(h.InconclusivePending) > 0 {
		return errors.New("trajectory awaits independent reconciliation or an attended review decision; further fixup is refused")
	}
	before, archive, err := currentTrajectorySource(s)
	if err != nil {
		return err
	}
	if !before.Clean {
		return errors.New("retained dirty patch requires explicit clean source promotion before another fixup")
	}
	actual, err := Observe(ctx, o.Root)
	if err != nil {
		return err
	}
	if actual != before {
		return errors.New("source differs from the frozen trajectory baseline")
	}
	o.before = actual
	o.archive = archive
	o.prepared = data
	return nil
}
func (o *Observer) AfterCycle(ctx context.Context, b budget.CycleBinding, ledger budget.Snapshot, key string) error {
	s, data, err := readState(statePath(b))
	if err != nil {
		return err
	}
	if o.prepared == nil || !bytes.Equal(data, o.prepared) {
		return errors.New("trajectory changed across reservation; charge remains spent")
	}
	o.prepared = nil
	dir, err := openIntentRoot(b, false)
	if err != nil {
		return err
	}
	defer dir.Close()
	i, sha, err := readReservationIntent(dir, key)
	if err != nil {
		return err
	}
	if o.intentSHA256 == "" || sha != o.intentSHA256 || !sameJSON(s, i.Before) {
		return errors.New("original precharge intent is missing or changed; charge remains spent")
	}
	a, err := intentAttempt(b, ledger, i, sha)
	if err != nil {
		return err
	}
	if a.Before != o.before || a.BeforeArchive != o.archive {
		return errors.New("precharge intent changed the original source")
	}
	s.Attempts = append(s.Attempts, a)
	if err = validateState(s, b, ledger); err != nil {
		return err
	}
	if err = checkStateSnapshots(ctx, b, s); err != nil {
		return err
	}
	return writeState(statePath(b), s)
}
func withState(ctx context.Context, root, idea string, fn func(budget.CycleBinding, budget.Snapshot, State) error) error {
	return withStateResolutionCheck(ctx, root, idea, checkResolutions, fn)
}

// Source/charge validation is unconditional. Recovery may reconstruct one exact
// parent observation while all other retained resolutions are checked normally.
func withStateResolutionCheck(ctx context.Context, root, idea string, check func(context.Context, budget.CycleBinding, State) error, fn func(budget.CycleBinding, budget.Snapshot, State) error) error {
	b, err := budget.LoadCycleBinding(ctx, root, idea, budget.Fixup)
	if err != nil {
		return err
	}
	if b == nil {
		return nil
	}
	wait, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	release, err := budget.AcquireResourceGuard(wait, filepath.Dir(b.Store.Dir))
	if err != nil {
		return err
	}
	defer release()
	b, err = budget.LoadCycleBinding(ctx, root, idea, budget.Fixup)
	if err != nil {
		return err
	}
	if b == nil {
		return errors.New("trajectory accounting disappeared")
	}
	if b.Policy.TrajectorySHA256 == "" {
		if _, err := os.Lstat(statePath(*b)); os.IsNotExist(err) {
			return nil
		}
		return errors.New("trajectory initialization is incomplete; exact activation replay is required")
	}
	ledger, err := b.Store.Inspect(ctx)
	if err != nil {
		return err
	}
	s, _, err := readState(statePath(*b))
	if err != nil {
		return err
	}
	if err = validateState(s, *b, ledger); err != nil {
		return err
	}
	if err = checkSourceSnapshots(ctx, *b, s); err != nil {
		return err
	}
	if err = check(ctx, *b, s); err != nil {
		return err
	}
	return fn(*b, ledger, s)
}
func Inspect(ctx context.Context, root, idea string) (*State, error) {
	var result *State
	err := withState(ctx, root, idea, func(_ budget.CycleBinding, _ budget.Snapshot, s State) error { result = &s; return nil })
	return result, err
}
func RequireResolved(ctx context.Context, root, idea string) error {
	return withState(ctx, root, idea, func(_ budget.CycleBinding, _ budget.Snapshot, s State) error {
		h := trajectoryHistory(s)
		if h.Unreconciled > 0 || h.ReviewPending || len(h.InconclusivePending) > 0 {
			return errors.New("trajectory has pending independent evidence or review decisions")
		}
		if len(s.Resolutions) > 0 && s.Resolutions[len(s.Resolutions)-1].Preview.Assessment.Outcome != NoRegression {
			return errors.New("latest charged patch has no confirmed clean material outcome")
		}
		expected, _, err := currentTrajectorySource(s)
		if err != nil {
			return err
		}
		actual, err := Observe(ctx, root)
		if err != nil {
			return err
		}
		if actual != expected {
			return errors.New("trajectory source changed without a retained charged patch")
		}
		return nil
	})
}

// Run is a runtime-only handle; it contains no authority to replay a model call.
type Run struct {
	root, idea   string
	charge       budget.CycleCharge
	invocation   string
	policySHA256 string
}

func Begin(ctx context.Context, root, idea, agent, invocation string) (*Run, error) {
	expected, err := budget.ActiveCycleTrajectory(ctx)
	if err != nil {
		return nil, err
	}
	var run *Run
	err = withState(ctx, root, idea, func(b budget.CycleBinding, _ budget.Snapshot, s State) error {
		if b.Policy.TrajectorySHA256 != expected {
			return errors.New("active trajectory policy changed before launch")
		}
		charge, err := budget.ActiveCycleCharge(ctx)
		if err != nil {
			return err
		}
		if len(s.Attempts) == 0 {
			return errors.New("trajectory launch has no charged attempt")
		}
		a := &s.Attempts[len(s.Attempts)-1]
		if a.Charge.EntryKey != charge.EntryKey || a.Launch != nil || agent != s.Policy.Implementer || !safeLabel(invocation) {
			return errors.New("trajectory launch changed implementer, charge or invocation")
		}
		actual, err := Observe(ctx, root)
		if err != nil {
			return err
		}
		if actual != a.Before {
			return errors.New("source changed between reservation and model launch")
		}
		a.Launch = &Launch{invocation, agent}
		if err = writeState(statePath(b), s); err != nil {
			return err
		}
		run = &Run{root, idea, charge, invocation, expected}
		return nil
	})
	if err == nil && expected != "" && run == nil {
		err = errors.New("required trajectory disappeared before launch")
	}
	return run, err
}
func (r *Run) Finish(ctx context.Context, status string, exit *int) error {
	if r == nil {
		return nil
	}
	observed := false
	err := withState(ctx, r.root, r.idea, func(b budget.CycleBinding, ledger budget.Snapshot, s State) error {
		observed = true
		if b.Policy.TrajectorySHA256 != r.policySHA256 {
			return errors.New("active trajectory policy changed before terminal publication")
		}
		if err := r.charge.Check(ledger); err != nil {
			return err
		}
		a := &s.Attempts[len(s.Attempts)-1]
		if a.Charge.EntryKey != r.charge.EntryKey || a.Launch == nil || a.Launch.InvocationID != r.invocation || a.Terminal != nil {
			return errors.New("trajectory terminal changed or was already published")
		}
		actual, snapshotErr := Observe(ctx, r.root)
		var archiveErr error
		a.Terminal = &Terminal{At: time.Now().UTC(), Status: status, ExitCode: exit}
		if snapshotErr == nil {
			a.After = &actual
			archive, err := CaptureSnapshot(ctx, r.root, snapshotDirectory(b), actual)
			archiveErr = err
			if err == nil {
				a.AfterArchive = &archive
			} else {
				a.Terminal.SnapshotError = "archive-unavailable"
			}
		} else {
			a.Terminal.SnapshotError = "source-unavailable"
		}
		if err := validateState(s, b, ledger); err != nil {
			return err
		}
		if err := writeState(statePath(b), s); err != nil {
			return err
		}
		if archiveErr != nil {
			return fmt.Errorf("trajectory retained unavailable archive: %w", archiveErr)
		}
		if snapshotErr != nil {
			return fmt.Errorf("trajectory retained an unavailable post-state: %w", snapshotErr)
		}
		return nil
	})
	if err == nil && !observed {
		return errors.New("required trajectory disappeared before terminal publication")
	}
	return err
}

func checkStateSnapshots(ctx context.Context, b budget.CycleBinding, s State) error {
	if err := checkSourceSnapshots(ctx, b, s); err != nil {
		return err
	}
	return checkResolutions(ctx, b, s)
}
func checkSourceSnapshots(ctx context.Context, b budget.CycleBinding, s State) error {
	if err := checkReservationIntents(ctx, b, s); err != nil {
		return err
	}
	dir := snapshotDirectory(b)
	if err := CheckSnapshot(ctx, dir, s.BaselineArchive, s.Policy.Baseline); err != nil {
		return fmt.Errorf("baseline archive unavailable: %w", err)
	}
	for _, a := range s.Attempts {
		if err := CheckSnapshot(ctx, dir, a.BeforeArchive, a.Before); err != nil {
			return fmt.Errorf("before archive unavailable: %w", err)
		}
		if a.AfterArchive != nil {
			if a.After == nil {
				return errors.New("archive without actual post-state")
			}
			if err := CheckSnapshot(ctx, dir, *a.AfterArchive, *a.After); err != nil {
				return fmt.Errorf("after archive unavailable: %w", err)
			}
		}
	}
	for _, c := range s.Continuations {
		if err := CheckSnapshot(ctx, dir, c.Archive, c.Preview.Source); err != nil {
			return fmt.Errorf("promoted source archive unavailable: %w", err)
		}
	}
	return nil
}
