package trajectory

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/evidence"
)

const MaxCapturedGitBytes int64 = 512 << 20

func checkCapturedGitStorage(ctx context.Context, root string) error {
	data, err := gitOutput(ctx, root, "rev-parse", "--git-path", "objects")
	if err != nil {
		return err
	}
	objects := strings.TrimSpace(string(data))
	if !filepath.IsAbs(objects) {
		objects = filepath.Join(root, objects)
	}
	if info, err := os.Lstat(filepath.Join(objects, "info", "alternates")); err == nil && info.Size() > 0 {
		return errors.New("captured verification requires independent Git objects; alternates are unsupported")
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	var total int64
	entries := 0
	return filepath.WalkDir(objects, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if err = ctx.Err(); err != nil {
			return err
		}
		entries++
		if entries > MaxSnapshotEntries {
			return errors.New("Git object inventory exceeds verification bound")
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() || info.Size() > MaxCapturedGitBytes-total {
			return errors.New("Git object history is unsupported or exceeds the captured verification bound")
		}
		total += info.Size()
		return nil
	})
}

// CapturedRequest binds a complete charged attempt to its original material
// scope and actual archived before/after worktrees. Unlike Request, it allows
// dirty output with the same original HEAD. The two exact source archives bind
// the complete change; it makes no clean-commit or synthetic ancestry claim.
type CapturedRequest struct {
	Version       int                `json:"version"`
	Kind          string             `json:"kind"`
	Idea          string             `json:"idea"`
	PolicySHA256  string             `json:"policy_sha256"`
	AttemptSHA256 string             `json:"attempt_sha256"`
	Sequence      int                `json:"sequence"`
	Charge        budget.CycleCharge `json:"charge"`
	InvocationID  string             `json:"invocation_id"`
	Implementer   string             `json:"implementer"`
	Verifier      string             `json:"verifier"`
	Before        Source             `json:"before"`
	After         Source             `json:"after"`
	BeforeArchive SnapshotRef        `json:"before_archive"`
	AfterArchive  SnapshotRef        `json:"after_archive"`
	Criteria      []CriterionBinding `json:"criteria"`
}

func (r CapturedRequest) SHA256() (string, error) {
	if r.Version != 1 || r.Kind != "captured-worktree" || !safeLabel(r.Idea) ||
		!validHash(r.PolicySHA256) || !validHash(r.AttemptSHA256) ||
		r.Sequence < 1 || r.Sequence > MaxPatches || !safeLabel(r.InvocationID) ||
		!safeLabel(r.Implementer) || !safeLabel(r.Verifier) || r.Implementer == r.Verifier ||
		!sourceValid(r.Before) || !sourceValid(r.After) || !r.Before.Clean ||
		r.Before.Tree.SHA256 == r.After.Tree.SHA256 ||
		!r.BeforeArchive.valid() || !r.AfterArchive.valid() || r.BeforeArchive == r.AfterArchive ||
		r.Charge.Kind != budget.Fixup || r.Charge.Scope == "" ||
		!validHash(r.Charge.EntryKey) || r.Charge.StartedAt.IsZero() ||
		r.Charge.ReservedAt.Before(r.Charge.StartedAt) ||
		(r.Charge.ReserveMicros != nil && *r.Charge.ReserveMicros < 0) ||
		len(r.Criteria) == 0 || len(r.Criteria) > MaxCriteria {
		return "", errors.New("invalid captured source request or non-independent verifier")
	}
	seen := map[string]bool{}
	for _, c := range r.Criteria {
		if !safeLabel(c.Name) || !validHash(c.CommandSHA256) || seen[c.Name] {
			return "", errors.New("invalid captured material scope")
		}
		seen[c.Name] = true
	}
	data, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return digest(data), nil
}

func capturedRequest(s State, verifier string) (CapturedRequest, error) {
	return capturedRequestAt(s, verifier, len(s.Attempts))
}

func capturedRequestAt(s State, verifier string, sequence int) (CapturedRequest, error) {
	sha, err := s.Policy.SHA256()
	if err != nil {
		return CapturedRequest{}, err
	}
	if sequence < 1 || sequence > len(s.Attempts) {
		return CapturedRequest{}, errors.New("no charged patch attempt to verify")
	}
	a := s.Attempts[sequence-1]
	if a.Launch == nil || a.Terminal == nil || a.After == nil || a.AfterArchive == nil || a.Terminal.SnapshotError != "" {
		return CapturedRequest{}, errors.New("charged attempt lacks its actual terminal or retained source")
	}
	if a.Before.Tree.SHA256 == a.After.Tree.SHA256 {
		return CapturedRequest{}, errors.New("unchanged source is not a new patch regression")
	}
	data, err := canonical(a)
	if err != nil {
		return CapturedRequest{}, err
	}
	r := CapturedRequest{
		Version: 1, Kind: "captured-worktree", Idea: s.Policy.Idea,
		PolicySHA256: sha, AttemptSHA256: digest(data), Sequence: a.Sequence,
		Charge: a.Charge, InvocationID: a.Launch.InvocationID,
		Implementer: s.Policy.Implementer, Verifier: verifier,
		Before: a.Before, After: *a.After, BeforeArchive: a.BeforeArchive,
		AfterArchive: *a.AfterArchive, Criteria: s.Policy.Criteria,
	}
	_, err = r.SHA256()
	return r, err
}

// FreezeCaptured reads the complete shared charge/state/archive authority.
// Caller-supplied patch inventories, criterion reductions or expected sources
// cannot replace the actual attempt. This freezes a request, not acceptance.
func FreezeCaptured(ctx context.Context, root, idea, verifier string) (CapturedRequest, error) {
	var request CapturedRequest
	err := withState(ctx, root, idea, func(b budget.CycleBinding, _ budget.Snapshot, s State) error {
		var err error
		request, err = capturedRequest(s, verifier)
		if err != nil {
			return err
		}
		_, err = checkCapturedActivationQuorum(ctx, b, s, request)
		return err
	})
	if err != nil {
		return CapturedRequest{}, err
	}
	if _, err = request.SHA256(); err != nil {
		return CapturedRequest{}, errors.New("captured trajectory is not active")
	}
	return request, nil
}

func checkCapturedAuthority(ctx context.Context, root string, r CapturedRequest) (string, error) {
	want, err := r.SHA256()
	if err != nil {
		return "", err
	}
	var archiveDir string
	err = withState(ctx, root, r.Idea, func(b budget.CycleBinding, _ budget.Snapshot, s State) error {
		current, err := capturedRequestAt(s, r.Verifier, r.Sequence)
		if err != nil {
			return err
		}
		actual, err := current.SHA256()
		if err != nil || actual != want {
			return errors.New("original charge, source, scope or invocation changed after comparison freeze")
		}
		if _, err := checkCapturedActivationQuorum(ctx, b, s, current); err != nil {
			return err
		}
		archiveDir = snapshotDirectory(b)
		return nil
	})
	if err != nil {
		return "", err
	}
	if archiveDir == "" {
		return "", errors.New("required captured trajectory authority disappeared")
	}
	return archiveDir, nil
}

// CapturedWorkspace owns only newly allocated private verification directories.
// Original sources, Git index and branch references are never mutation targets.
// A caller may inspect the roots but cannot replace the frozen request.
type CapturedWorkspace struct {
	mu      sync.Mutex
	owner   string
	origin  string
	before  string
	after   string
	request CapturedRequest
	closed  bool
	used    bool
}

func (w *CapturedWorkspace) BeforeRoot() string { return w.before }
func (w *CapturedWorkspace) AfterRoot() string  { return w.after }
func (w *CapturedWorkspace) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed {
		return nil
	}
	w.closed = true
	return os.RemoveAll(w.owner)
}

// OpenCaptured restores archives and attaches separate local copies of the
// original Git object history, detached at the recorded real commits. It does
// not create synthetic commits. Exact original Source observations must match.
//
// The archive currently retains worktree bytes, not arbitrary historical index
// blobs or Git object backups. Unavailable commits or staged-state differences
// therefore refuse visibly; current files or new commits cannot fill that gap.
func OpenCaptured(ctx context.Context, root string, r CapturedRequest, parent string) (_ *CapturedWorkspace, resultErr error) {
	archiveDir, err := checkCapturedAuthority(ctx, root, r)
	if err != nil {
		return nil, err
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return nil, err
	}
	root, err = filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	if parent == "" {
		parent = os.TempDir()
	}
	if _, err = gitOutput(ctx, parent, "rev-parse", "--git-dir"); err == nil {
		return nil, errors.New("captured verification parent must be outside every Git repository")
	}
	if err = ctx.Err(); err != nil {
		return nil, err
	}
	owner, err := os.MkdirTemp(parent, "trajectory-verifier-")
	if err != nil {
		return nil, err
	}
	defer func() {
		if resultErr != nil {
			_ = os.RemoveAll(owner)
		}
	}()
	before, err := RestoreSnapshot(ctx, archiveDir, r.BeforeArchive, r.Before, owner)
	if err != nil {
		return nil, err
	}
	after, err := RestoreSnapshot(ctx, archiveDir, r.AfterArchive, r.After, owner)
	if err != nil {
		return nil, err
	}
	template := filepath.Join(owner, "empty-template")
	if err = os.Mkdir(template, 0700); err != nil {
		return nil, err
	}
	if err = checkCapturedGitStorage(ctx, root); err != nil {
		return nil, err
	}
	for i, source := range []struct {
		path string
		want Source
	}{{before, r.Before}, {after, r.After}} {
		clone := filepath.Join(owner, fmt.Sprintf("git-copy-%d", i))
		// Local clone copies object files and avoids an upload-pack process.
		// --no-hardlinks prevents shared writable object storage; --no-checkout
		// avoids filters, and the empty template supplies no execution hooks.
		if _, err = gitOutput(ctx, root, "clone", "--quiet", "--local", "--no-hardlinks", "--no-checkout", "--template="+template, "--", root, clone); err != nil {
			return nil, errors.New("cannot independently copy the original Git object history")
		}
		if err = checkCapturedGitStorage(ctx, clone); err != nil {
			return nil, err
		}
		if err = os.Rename(filepath.Join(clone, ".git"), filepath.Join(source.path, ".git")); err != nil {
			return nil, errors.New("cannot attach isolated Git metadata to retained source")
		}
		if _, err = gitOutput(ctx, source.path, "update-ref", "--no-deref", "HEAD", source.want.Tree.Commit); err != nil {
			return nil, errors.New("recorded source commit is unavailable in isolated history")
		}
		if _, err = gitOutput(ctx, source.path, "read-tree", "--reset", source.want.Tree.Commit); err != nil {
			return nil, errors.New("recorded source index cannot be reconstructed from its commit")
		}
		actual, err := Observe(ctx, source.path)
		if err != nil || actual != source.want {
			return nil, errors.New("retained worktree and original Git status cannot be reproduced exactly; preserve staged state for explicit recovery")
		}
	}
	if _, err = checkCapturedAuthority(ctx, root, r); err != nil {
		return nil, err
	}
	// Clone pointer-bearing fields before retaining a caller's request.
	data, err := canonical(r)
	if err != nil {
		return nil, err
	}
	var frozen CapturedRequest
	if err = json.Unmarshal(data, &frozen); err != nil {
		return nil, err
	}
	return &CapturedWorkspace{owner: owner, origin: root, before: before, after: after, request: frozen}, nil
}

func (w *CapturedWorkspace) check(ctx context.Context) error {
	if _, err := checkCapturedAuthority(ctx, w.origin, w.request); err != nil {
		return err
	}
	for _, source := range []struct {
		path string
		want Source
	}{{w.before, w.request.Before}, {w.after, w.request.After}} {
		actual, err := Observe(ctx, source.path)
		if err != nil || actual != source.want {
			return errors.New("isolated captured source or original Git identity changed")
		}
	}
	return nil
}

// VerifyCaptured executes original criteria in AB/BA order on the actual
// captured worktrees, including dirty failure output. A workspace is single-use
// so a failed or interrupted helper cannot quietly overwrite its observations
// by replaying execution. The parent must retain every returned partial result.
// This function itself neither invokes nor authenticates the selected model.
func VerifyCaptured(ctx context.Context, w *CapturedWorkspace, criteria []Criterion) (Observation, error) {
	return verifyCaptured(ctx, w, criteria, nil, nil, nil)
}

// The journal hooks are private: callers cannot replace source validation or
// turn an observation-write failure into an accepted execution.
func verifyCaptured(ctx context.Context, w *CapturedWorkspace, criteria []Criterion, guard func() error, retain func(int, Execution) error, control func(int) evidence.CriterionStartControl) (Observation, error) {
	if w == nil {
		return Observation{}, errors.New("captured verification workspace is required")
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed || w.used {
		return Observation{}, errors.New("captured verification workspace is closed or already used")
	}
	w.used = true
	r := w.request
	sha, err := r.SHA256()
	if err != nil {
		return Observation{}, err
	}
	o := Observation{Version: 1, RequestSHA256: sha, Verifier: r.Verifier}
	if len(criteria) != len(r.Criteria) {
		return o, errors.New("captured verification omitted original material criteria")
	}
	for i, c := range criteria {
		if c.Name != r.Criteria[i].Name || digest([]byte(c.Command)) != r.Criteria[i].CommandSHA256 {
			return o, errors.New("captured verification changed frozen criterion command or order")
		}
	}
	check := func() error {
		if guard != nil {
			if err := guard(); err != nil {
				return err
			}
		}
		return w.check(ctx)
	}
	if err = check(); err != nil {
		return o, err
	}
	ordinal := 0
	for _, c := range criteria {
		o.Pairs = append(o.Pairs, Pair{Name: c.Name})
		p := &o.Pairs[len(o.Pairs)-1]
		for _, step := range []struct {
			after bool
			index int
		}{{false, 0}, {true, 0}, {true, 1}, {false, 1}} {
			if err = check(); err != nil {
				return o, err
			}
			root, source, target := w.before, r.Before, &p.Before[step.index]
			if step.after {
				root, source, target = w.after, r.After, &p.After[step.index]
			}
			var startControl evidence.CriterionStartControl
			if control != nil {
				startControl = control(ordinal + 1)
			}
			execution := evidence.RunCriterionControlled(ctx, root, c.Name, c.Command, r.Verifier, startControl)
			*target = Execution{Complete: execution.Complete, Record: execution.Record, TreeBeforeSHA256: source.Tree.SHA256}
			actual, observeErr := Observe(ctx, root)
			if observeErr == nil {
				target.TreeAfterSHA256 = actual.Tree.SHA256
			}
			ordinal++
			if retain != nil {
				if err = retain(ordinal, *target); err != nil {
					return o, fmt.Errorf("captured execution retention failed: %w", err)
				}
			}
			if observeErr != nil {
				return o, fmt.Errorf("captured criterion post-source unavailable: %w", observeErr)
			}
			if err = check(); err != nil {
				return o, err
			}
			if !execution.Complete {
				return o, errors.New("captured criterion execution is incomplete; retain partial observations")
			}
		}
	}
	if _, err = AssessCaptured(r, o); err != nil {
		return o, err
	}
	return o, nil
}

func AssessCaptured(r CapturedRequest, o Observation) (Assessment, error) {
	sha, err := r.SHA256()
	if err != nil {
		return Assessment{}, err
	}
	if o.Version != 1 || o.RequestSHA256 != sha || o.Verifier != r.Verifier {
		return Assessment{}, errors.New("missing or mismatched captured patch observation")
	}
	return assessPairs(r.Criteria, r.Before.Tree, r.After.Tree, r.Verifier, o.Pairs)
}
