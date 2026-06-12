package runner

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"parley-deck-cli/internal/fsutil"
	"parley-deck-cli/internal/procctl"
)

// Review snapshots (runner-hardening-kindly D9, adapted from kindly's
// create_snapshot): Phase 6 reviewers read a disposable shared-clone checkout
// on LOCAL tmp instead of the live virtio-fs tree. A clean tree pins HEAD; a
// dirty tree is captured as a temp-index snapshot commit parented on HEAD so
// uncommitted implementations are reviewed exactly as they stand. The reviewer
// writes its artifact inside the snapshot; the runner validates it and moves
// it back to the canonical deck path.

// ReviewSnapshot is one reviewer's disposable checkout.
type ReviewSnapshot struct {
	Dir  string // snapshot worktree root (local tmp)
	SHA  string // reviewed commit (HEAD or the snapshot commit)
	Mode string // "head" | "snapshot-commit"

	liveRoot string
	marker   string
}

// snapshotUnavailable distinguishes "fall back to the live tree" causes from
// hard errors; callers emit review.snapshot_fallback with the reason.
type snapshotUnavailable struct{ reason string }

func (e snapshotUnavailable) Error() string { return "snapshot unavailable: " + e.reason }

// gitProbe builds a read-only git command with optional locks disabled
// (consensus D8: probes must never write .git on the weakly-coherent mount).
func gitProbe(dir string, args ...string) *exec.Cmd {
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(), "GIT_OPTIONAL_LOCKS=0")
	return cmd
}

// CreateReviewSnapshot builds the checkout for one reviewer. A nil snapshot
// with a snapshotUnavailable error means "review the live tree" (not a git
// repo, staged/worktree divergence, or creation failure).
func CreateReviewSnapshot(liveRoot, ideaSlug, roundLabel, agentID, runID string) (*ReviewSnapshot, error) {
	if err := gitProbe(liveRoot, "rev-parse", "--is-inside-work-tree").Run(); err != nil {
		return nil, snapshotUnavailable{reason: "not a git work tree"}
	}
	// kindly's pre-check: a snapshot commit holds ONE version per path. When
	// staged content diverges from the worktree, only a live review keeps both
	// versions in scope.
	staged, _ := gitProbe(liveRoot, "diff", "--cached", "--name-only").Output()
	unstaged, _ := gitProbe(liveRoot, "diff-files", "--name-only").Output()
	if hasIntersection(splitLines(string(staged)), splitLines(string(unstaged))) {
		return nil, snapshotUnavailable{reason: "staged content diverges from the working tree"}
	}

	repoHash := sha1.Sum([]byte(liveRoot))
	base := filepath.Join(os.TempDir(), "parley-review-snapshots", hex.EncodeToString(repoHash[:6]))
	dir := filepath.Join(base, ideaSlug, strings.ReplaceAll(roundLabel, string(filepath.Separator), "-"), agentID)
	marker := dir + ".pid"
	if err := fsutil.MkdirAllResilient(filepath.Dir(dir), 0o755); err != nil {
		return nil, snapshotUnavailable{reason: "snapshot base: " + err.Error()}
	}
	sweepStaleSnapshots(filepath.Dir(dir))
	if ownerAlive(marker) {
		// A live concurrent run owns the stable path; step aside (kindly rule).
		dir = fmt.Sprintf("%s-%d", dir, os.Getpid())
		marker = dir + ".pid"
	} else {
		_ = os.RemoveAll(dir)
		_ = os.Remove(marker)
	}
	if err := writeMarker(marker, runID, agentID); err != nil {
		return nil, snapshotUnavailable{reason: "marker: " + err.Error()}
	}

	snap := &ReviewSnapshot{Dir: dir, liveRoot: liveRoot, marker: marker}
	if err := snap.materialize(); err != nil {
		snap.Cleanup()
		return nil, snapshotUnavailable{reason: err.Error()}
	}
	return snap, nil
}

func (s *ReviewSnapshot) materialize() error {
	clone := exec.Command("git", "clone", "--quiet", "--shared", "--no-checkout", "--", s.liveRoot, s.Dir)
	clone.Env = append(os.Environ(), "GIT_OPTIONAL_LOCKS=0")
	if out, err := clone.CombinedOutput(); err != nil {
		return fmt.Errorf("clone: %v: %s", err, strings.TrimSpace(string(out)))
	}
	// Carry local-only ignore rules into the clone (kindly behavior).
	if exclude, err := gitProbe(s.liveRoot, "rev-parse", "--path-format=absolute", "--git-path", "info/exclude").Output(); err == nil {
		src := strings.TrimSpace(string(exclude))
		if data, err := os.ReadFile(src); err == nil {
			_ = os.MkdirAll(filepath.Join(s.Dir, ".git", "info"), 0o755)
			_ = os.WriteFile(filepath.Join(s.Dir, ".git", "info", "exclude"), data, 0o644)
		}
	}

	dirty, _ := gitProbe(s.liveRoot, "status", "--porcelain", "--untracked-files=all").Output()
	if strings.TrimSpace(string(dirty)) == "" {
		head, err := gitProbe(s.liveRoot, "rev-parse", "HEAD").Output()
		if err != nil {
			return fmt.Errorf("rev-parse HEAD: %v", err)
		}
		s.SHA = strings.TrimSpace(string(head))
		s.Mode = "head"
	} else {
		// Temp-index snapshot commit of the live working tree, into the CLONE's
		// object store, parented on HEAD (kindly create_snapshot; verified by
		// codex with a local-tmp clone and a virtio-fs origin).
		tempIndex := s.Dir + ".index"
		defer os.Remove(tempIndex)
		gitDir := filepath.Join(s.Dir, ".git")
		snapGit := func(args ...string) *exec.Cmd {
			cmd := exec.Command("git", append([]string{"-C", s.liveRoot, "--git-dir", gitDir, "--work-tree", "."}, args...)...)
			cmd.Env = append(os.Environ(), "GIT_OPTIONAL_LOCKS=0", "GIT_INDEX_FILE="+tempIndex)
			return cmd
		}
		if out, err := snapGit("read-tree", "HEAD").CombinedOutput(); err != nil {
			return fmt.Errorf("read-tree: %v: %s", err, strings.TrimSpace(string(out)))
		}
		if out, err := snapGit("add", "-A").CombinedOutput(); err != nil {
			return fmt.Errorf("add -A: %v: %s", err, strings.TrimSpace(string(out)))
		}
		treeOut, err := snapGit("write-tree").Output()
		if err != nil {
			return fmt.Errorf("write-tree: %v", err)
		}
		commitCmd := exec.Command("git", "--git-dir", gitDir,
			"-c", "user.name=parley", "-c", "user.email=parley@localhost",
			"commit-tree", strings.TrimSpace(string(treeOut)), "-p", "HEAD", "-m", "parley review snapshot")
		commitCmd.Env = append(os.Environ(), "GIT_OPTIONAL_LOCKS=0")
		commitOut, err := commitCmd.Output()
		if err != nil {
			return fmt.Errorf("commit-tree: %v", err)
		}
		s.SHA = strings.TrimSpace(string(commitOut))
		s.Mode = "snapshot-commit"
	}

	checkout := gitProbe(s.Dir, "checkout", "--quiet", "--detach", s.SHA)
	if out, err := checkout.CombinedOutput(); err != nil {
		return fmt.Errorf("checkout: %v: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// MoveArtifactBack copies the reviewer's artifact from the snapshot to the
// canonical deck path: copy + fsync + rename WITHIN the target directory —
// rename(2) across devices would fail (consensus D9/hermes).
func (s *ReviewSnapshot) MoveArtifactBack(relPath, canonicalPath string) error {
	src := filepath.Join(s.Dir, relPath)
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("snapshot artifact: %w", err)
	}
	if err := fsutil.MkdirAllResilient(filepath.Dir(canonicalPath), 0o755); err != nil {
		return err
	}
	tmp := canonicalPath + ".snapshot-tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		f.Close()
		_ = os.Remove(tmp)
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		_ = os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, canonicalPath)
}

// Abandon keeps the snapshot directory for manual recovery (a failed artifact
// move-back, consensus D9 review fix 2) and removes only the marker — the
// stale sweep skips directories without a .pid marker, so the retained
// artifact survives later runs.
func (s *ReviewSnapshot) Abandon() {
	if s == nil {
		return
	}
	if s.marker != "" {
		_ = os.Remove(s.marker)
	}
}

// Cleanup deletes the snapshot worktree and its marker (plain delete — the
// shared clone never wrote the origin's .git).
func (s *ReviewSnapshot) Cleanup() {
	if s == nil {
		return
	}
	if s.Dir != "" {
		_ = os.RemoveAll(s.Dir)
	}
	if s.marker != "" {
		_ = os.Remove(s.marker)
	}
}

type snapshotMarker struct {
	PID    int    `json:"pid"`
	BootID string `json:"boot_id"`
	RunID  string `json:"run_id"`
	Agent  string `json:"agent"`
}

func writeMarker(path, runID, agentID string) error {
	data, _ := json.Marshal(snapshotMarker{PID: os.Getpid(), BootID: procctl.CurrentBootID(), RunID: runID, Agent: agentID})
	return os.WriteFile(path, data, 0o644)
}

func ownerAlive(markerPath string) bool {
	data, err := os.ReadFile(markerPath)
	if err != nil {
		return false
	}
	var m snapshotMarker
	if json.Unmarshal(data, &m) != nil || m.PID <= 0 {
		return false
	}
	if m.BootID != "" && m.BootID != procctl.CurrentBootID() {
		return false // marker from a previous boot — its pid space is meaningless
	}
	return procctl.Alive(procctl.Spawned{PID: m.PID})
}

// sweepStaleSnapshots heals crashed runs: any sibling snapshot whose marker
// names a dead process is removed (kindly's stray sweep).
func sweepStaleSnapshots(parent string) {
	entries, err := os.ReadDir(parent)
	if err != nil {
		return
	}
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".pid") {
			continue
		}
		marker := filepath.Join(parent, name)
		if ownerAlive(marker) {
			continue
		}
		_ = os.RemoveAll(strings.TrimSuffix(marker, ".pid"))
		_ = os.Remove(marker)
	}
}

func splitLines(s string) []string {
	var out []string
	for _, line := range strings.Split(s, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			out = append(out, line)
		}
	}
	return out
}

func hasIntersection(a, b []string) bool {
	set := make(map[string]bool, len(a))
	for _, item := range a {
		set[item] = true
	}
	for _, item := range b {
		if set[item] {
			return true
		}
	}
	return false
}
