package trajectory

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"parley-deck-cli/internal/evidence"
)

// Freeze binds independent clean Git snapshots and exact command hashes before
// execution. Commands are supplied in memory; no raw commands, diff bodies,
// environment or absolute workspace paths are persisted in the request.
func Freeze(ctx context.Context, r Request, beforeRoot, afterRoot string, criteria []Criterion) (Request, error) {
	if len(criteria) == 0 || len(criteria) > MaxCriteria {
		return Request{}, errors.New("invalid material criterion scope")
	}
	beforeRoot, afterRoot, err := distinctRoots(beforeRoot, afterRoot)
	if err != nil {
		return Request{}, err
	}
	r.Version = 1
	r.Before, err = snapshot(ctx, beforeRoot)
	if err != nil {
		return Request{}, err
	}
	r.After, err = snapshot(ctx, afterRoot)
	if err != nil {
		return Request{}, err
	}
	r.PatchSHA256, err = patchDigest(ctx, afterRoot, r.Before.Commit, r.After.Commit)
	if err != nil {
		return Request{}, err
	}
	r.Criteria = nil
	for _, c := range criteria {
		if strings.TrimSpace(c.Command) == "" || len(c.Command) > 16<<10 || strings.ContainsRune(c.Command, 0) {
			return Request{}, errors.New("empty or excessive material criterion command")
		}
		r.Criteria = append(r.Criteria, CriterionBinding{Name: c.Name, CommandSHA256: digest([]byte(c.Command))})
	}
	if err := r.validate(); err != nil {
		return Request{}, err
	}
	return r, nil
}

// Verify executes each frozen criterion twice on both snapshots, in AB/BA
// order. It records only actual executions; an interrupted run returns the
// partial observation plus an error. The orchestrator must retain that failure
// and must never turn the partial result into a confirmation.
//
// The runtime caller must independently select and observe r.Verifier. This
// function alone does not launch or authenticate an independent model agent.
func Verify(ctx context.Context, r Request, beforeRoot, afterRoot string, criteria []Criterion) (Observation, error) {
	sha, err := r.SHA256()
	if err != nil {
		return Observation{}, err
	}
	o := Observation{Version: 1, RequestSHA256: sha, Verifier: r.Verifier}
	beforeRoot, afterRoot, err = distinctRoots(beforeRoot, afterRoot)
	if err != nil {
		return o, err
	}
	if len(criteria) != len(r.Criteria) {
		return o, errors.New("verification omitted material criteria")
	}
	for i, c := range criteria {
		if c.Name != r.Criteria[i].Name || digest([]byte(c.Command)) != r.Criteria[i].CommandSHA256 {
			return o, errors.New("verification changed the frozen command or criterion order")
		}
	}
	check := func() error {
		if err := ctx.Err(); err != nil {
			return err
		}
		before, err := snapshot(ctx, beforeRoot)
		if err != nil {
			return err
		}
		after, err := snapshot(ctx, afterRoot)
		if err != nil {
			return err
		}
		patch, err := patchDigest(ctx, afterRoot, r.Before.Commit, r.After.Commit)
		if err != nil {
			return err
		}
		if before != r.Before || after != r.After || patch != r.PatchSHA256 {
			return errors.New("frozen before/after snapshot or patch changed")
		}
		return nil
	}
	if err := check(); err != nil {
		return o, err
	}
	for _, c := range criteria {
		o.Pairs = append(o.Pairs, Pair{Name: c.Name})
		p := &o.Pairs[len(o.Pairs)-1]
		for _, run := range []struct {
			after bool
			index int
		}{{false, 0}, {true, 0}, {true, 1}, {false, 1}} {
			if err := check(); err != nil {
				return o, err
			}
			root, tree, target := beforeRoot, r.Before, &p.Before[run.index]
			if run.after {
				root, tree, target = afterRoot, r.After, &p.After[run.index]
			}
			actual := evidence.RunCriterionDetailed(ctx, root, c.Name, c.Command, r.Verifier)
			*target = Execution{Complete: actual.Complete, Record: actual.Record, TreeBeforeSHA256: tree.SHA256}
			// Record the actual post-execution digest even on drift. Never fill
			// the expected digest into an observation of a different tree.
			observed, snapErr := snapshot(ctx, root)
			if snapErr == nil {
				target.TreeAfterSHA256 = observed.SHA256
			}
			if snapErr != nil {
				return o, fmt.Errorf("criterion snapshot after execution: %w", snapErr)
			}
			if err := check(); err != nil {
				return o, err
			}
			if !actual.Complete {
				return o, errors.New("criterion execution is incomplete; retain the observation without confirmation")
			}
		}
	}
	if _, err := Assess(r, o); err != nil {
		return o, err
	}
	return o, nil
}

func distinctRoots(a, b string) (string, string, error) {
	canonical := func(s string) (string, error) {
		if s == "" {
			return "", errors.New("snapshot root is required")
		}
		x, err := filepath.EvalSymlinks(s)
		if err != nil {
			return "", err
		}
		return filepath.Abs(x)
	}
	x, err := canonical(a)
	if err != nil {
		return "", "", err
	}
	y, err := canonical(b)
	if err != nil {
		return "", "", err
	}
	for _, pair := range [][2]string{{x, y}, {y, x}} {
		rel, err := filepath.Rel(pair[0], pair[1])
		if err != nil {
			return "", "", err
		}
		if rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))) {
			return "", "", errors.New("baseline and patched snapshots must be separate non-nested roots")
		}
	}
	return x, y, nil
}

type boundedOutput struct {
	buf      bytes.Buffer
	overflow bool
}

func (b *boundedOutput) Write(p []byte) (int, error) {
	room := (16 << 20) - b.buf.Len()
	if len(p) > room {
		b.overflow = true
	}
	if room > 0 {
		n := len(p)
		if n > room {
			n = room
		}
		_, _ = b.buf.Write(p[:n])
	}
	return len(p), nil
}
func gitOutput(ctx context.Context, root string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "git", append([]string{"-C", root, "--no-pager"}, args...)...)
	var out boundedOutput
	cmd.Stdout = &out
	// Git diagnostics can contain local paths or source snippets; error
	// messages here identify the failed operation without copying stderr.
	if err := cmd.Run(); err != nil {
		return nil, errors.New("Git snapshot operation failed")
	}
	if out.overflow {
		return nil, errors.New("Git snapshot output exceeds the comparison bound")
	}
	return out.buf.Bytes(), nil
}
func snapshot(ctx context.Context, root string) (Tree, error) {
	top, err := gitOutput(ctx, root, "rev-parse", "--show-toplevel")
	if err != nil {
		return Tree{}, err
	}
	actual, err := filepath.EvalSymlinks(strings.TrimSpace(string(top)))
	if err != nil || actual != root {
		return Tree{}, errors.New("snapshot must name its exact Git worktree root")
	}
	status, err := gitOutput(ctx, root, "status", "--porcelain=v1", "--untracked-files=all")
	if err != nil {
		return Tree{}, err
	}
	if len(status) != 0 {
		return Tree{}, errors.New("snapshot is dirty or contains untracked material")
	}
	commit, err := gitOutput(ctx, root, "rev-parse", "--verify", "HEAD^{commit}")
	if err != nil {
		return Tree{}, err
	}
	head := strings.TrimSpace(string(commit))
	if !validCommit(head) {
		return Tree{}, errors.New("snapshot has invalid commit identity")
	}
	tree, err := evidence.TreeDigest(root)
	if err != nil {
		return Tree{}, errors.New("cannot digest the complete snapshot tree")
	}
	return Tree{Commit: head, SHA256: tree}, nil
}
func patchDigest(ctx context.Context, root, before, after string) (string, error) {
	if !validCommit(before) || !validCommit(after) || before == after {
		return "", errors.New("patch must name distinct full commit identities")
	}
	if _, err := gitOutput(ctx, root, "merge-base", "--is-ancestor", before, after); err != nil {
		return "", errors.New("baseline is not an ancestor of the patched snapshot")
	}
	diff, err := gitOutput(ctx, root, "diff", "--no-ext-diff", "--no-textconv", "--binary", "--full-index", before, after, "--")
	if err != nil {
		return "", err
	}
	if len(diff) == 0 {
		return "", errors.New("unchanged patch cannot establish a new regression")
	}
	return digest(diff), nil
}
