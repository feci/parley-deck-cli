package budget

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

// Availability of one registered worktree path. Only "available" means the
// registration has a local working directory to read. Every other state leaves
// that worktree's history unknown; none of them is an absence of history and
// none is ever reported as a zero count.
const (
	WorktreeAvailable     = "available"
	WorktreeMissing       = "missing"
	WorktreeNotDirectory  = "not-directory"
	WorktreeIndeterminate = "indeterminate"
)

// History coverage is a location classification, never a count or a readability
// guarantee. Local means only that the registered path exists as a directory;
// Git metadata and historical contents have not been inspected. This stage never
// counts, caps or attests history, so an unavailable worktree cannot read as zero.
const (
	WorktreeHistoryLocal   = "local"
	WorktreeHistoryUnknown = "unknown"
)

const (
	worktreeInventorySchema = 1
	maxWorktreeListBytes    = 16 << 20
	maxWorktreeRegistration = 4096
)

// WorktreeRegistration is one raw registration exactly as Git reports it. The
// registered identity is preserved verbatim: an unavailable path is retained,
// never pruned, normalized away or recreated.
type WorktreeRegistration struct {
	RegisteredPath  string `json:"registered_path"`
	ResolvedPath    string `json:"resolved_path"`
	Head            string `json:"head"`
	Branch          string `json:"branch"`
	Detached        bool   `json:"detached"`
	Bare            bool   `json:"bare"`
	Locked          bool   `json:"locked"`
	LockReason      string `json:"lock_reason"`
	Prunable        bool   `json:"prunable"`
	PrunableReason  string `json:"prunable_reason"`
	Availability    string `json:"availability"`
	HistoryCoverage string `json:"history_coverage"`
}

// WorktreeInventory is a read-only observation. CoversRepository states whether
// the registrations below came from Git worktree authority at all; a non-Git
// root reports false and must not be read as repository-wide coverage.
// Uncertainty carries every fact this observation could not settle.
// ObservationSHA256 covers only the load-bearing registration and availability
// facts (see observedWorktrees): Root and RepoPrefix report where the caller
// stood but are excluded, as is any clock, so an unchanged repository
// re-inspects to the same digest from anywhere inside it.
type WorktreeInventory struct {
	Schema            int                    `json:"schema"`
	Root              string                 `json:"root"`
	Git               bool                   `json:"git"`
	GitCommonDir      string                 `json:"git_common_dir"`
	RepoPrefix        string                 `json:"repo_prefix"`
	SourceScope       string                 `json:"source_scope"`
	CoversRepository  bool                   `json:"covers_repository"`
	Registrations     []WorktreeRegistration `json:"registrations"`
	Uncertainty       []string               `json:"uncertainty"`
	ObservationSHA256 string                 `json:"observation_sha256"`
}

// observedWorktrees is the canonical digest input: the load-bearing
// availability and registration facts, with registrations ordered by their raw
// registered path so Git's enumeration order cannot change the digest.
//
// Three observed values are deliberately excluded. No timestamp takes part, so
// a clock cannot invalidate identical previews. Neither Root nor RepoPrefix
// takes part either: both describe where the caller stood, not what the
// repository registers, so the same repository previewed from any of its own
// worktrees or subdirectories hashes identically. GitCommonDir carries the
// repository identity.
type observedWorktrees struct {
	Schema           int                    `json:"schema"`
	Git              bool                   `json:"git"`
	CoversRepository bool                   `json:"covers_repository"`
	SourceScope      string                 `json:"source_scope"`
	GitCommonDir     string                 `json:"git_common_dir"`
	Registrations    []WorktreeRegistration `json:"registrations"`
}

// OperatorRecord renders the exact bytes a later attested operator record would
// have to quote. It is a pure function of this observation and grants nothing:
// no apply, attest or confirm path in this package consumes it.
func (i WorktreeInventory) OperatorRecord() ([]byte, error) {
	// An absent and an empty registration list are the same observation, so the
	// digest input normalizes both to one encoding.
	sorted := append([]WorktreeRegistration{}, i.Registrations...)
	slices.SortFunc(sorted, func(a, b WorktreeRegistration) int {
		return strings.Compare(a.RegisteredPath, b.RegisteredPath)
	})
	return json.Marshal(observedWorktrees{worktreeInventorySchema, i.Git, i.CoversRepository, i.SourceScope, i.GitCommonDir, sorted})
}

func (i WorktreeInventory) digest() (string, error) {
	raw, err := i.OperatorRecord()
	if err != nil {
		return "", err
	}
	return key(string(raw)), nil
}

// InspectWorktreeInventory enumerates every worktree registration, including
// registrations whose working directory is gone. It is read-only: it runs
// standard Git probes and stats each registered path. It never prunes, repairs,
// recreates a path, takes a lock or writes a byte, and it leaves
// launchScope(..., true)'s refusal on unavailable history untouched.
func InspectWorktreeInventory(ctx context.Context, root string) (WorktreeInventory, error) {
	i := WorktreeInventory{Schema: worktreeInventorySchema, Registrations: []WorktreeRegistration{}, Uncertainty: []string{}}
	abs, err := filepath.Abs(root)
	if err != nil {
		return WorktreeInventory{}, err
	}
	abs, err = filepath.EvalSymlinks(abs)
	if err != nil {
		return WorktreeInventory{}, err
	}
	i.Root = abs
	common, gitErr := exec.CommandContext(ctx, "git", "-C", abs, "rev-parse", "--path-format=absolute", "--git-common-dir").Output()
	if gitErr != nil {
		if ctx.Err() != nil {
			return WorktreeInventory{}, ctx.Err()
		}
		// A broken Git repository is refused, exactly as shared-scope resolution
		// refuses it. It must not degrade into a single-directory answer that
		// silently drops every other registration in the repository.
		for ancestor := abs; ; ancestor = filepath.Dir(ancestor) {
			if _, e := os.Lstat(filepath.Join(ancestor, ".git")); e == nil || !os.IsNotExist(e) {
				return WorktreeInventory{}, fmt.Errorf("cannot resolve Git worktree authority: %w", gitErr)
			}
			if filepath.Dir(ancestor) == ancestor {
				break
			}
		}
		// A genuine non-Git root is explicit and claims nothing: no Git
		// authority was found, so this observation describes only the one
		// directory inspected and establishes no repository-wide coverage.
		i.SourceScope = "single-directory/non-git"
		i.Uncertainty = append(i.Uncertainty,
			"no Git worktree authority was found at or above this root; no registration could be enumerated",
			"this observation covers only the inspected directory and does not establish repository-wide history coverage")
		if i.ObservationSHA256, err = i.digest(); err != nil {
			return WorktreeInventory{}, err
		}
		return i, nil
	}
	gitDir, err := filepath.EvalSymlinks(strings.TrimSpace(string(common)))
	if err != nil {
		return WorktreeInventory{}, err
	}
	i.Git, i.CoversRepository, i.GitCommonDir = true, true, gitDir
	i.SourceScope = "git-worktree-list/" + gitDir
	// The prefix is orientation only and is not part of the digest, so a root
	// with no work tree (a bare repository) still has inspectable
	// registrations. Authority is already established by the common dir above,
	// and the registration read below must still succeed.
	if prefix, e := exec.CommandContext(ctx, "git", "-C", abs, "rev-parse", "--show-prefix").Output(); e == nil {
		i.RepoPrefix = strings.TrimSpace(string(prefix))
	} else {
		if ctx.Err() != nil {
			return WorktreeInventory{}, ctx.Err()
		}
		i.Uncertainty = append(i.Uncertainty, "this root has no reportable work-tree prefix, so the inspected location is not described; the registrations below are unaffected")
	}
	out, err := exec.CommandContext(ctx, "git", "-C", abs, "worktree", "list", "--porcelain", "-z").Output()
	if err != nil {
		return WorktreeInventory{}, fmt.Errorf("cannot read Git worktree registrations: %w", err)
	}
	registrations, err := parseWorktreeRegistrations(out)
	if err != nil {
		return WorktreeInventory{}, err
	}
	for n := range registrations {
		classifyWorktreeAvailability(&registrations[n])
		r := registrations[n]
		if r.Availability != WorktreeAvailable {
			i.Uncertainty = append(i.Uncertainty, fmt.Sprintf("registration %q is %s; its history is unknown here and is retained unpruned, not counted as absent", r.RegisteredPath, r.Availability))
		}
		if r.Locked {
			i.Uncertainty = append(i.Uncertainty, fmt.Sprintf("registration %q is locked; the lock is reported, not cleared", r.RegisteredPath))
		}
		if r.Prunable {
			i.Uncertainty = append(i.Uncertainty, fmt.Sprintf("registration %q is reported prunable by Git; it is retained and not pruned here", r.RegisteredPath))
		}
	}
	i.Registrations = registrations
	if i.ObservationSHA256, err = i.digest(); err != nil {
		return WorktreeInventory{}, err
	}
	return i, nil
}

// Availability is decided by Lstat on the raw registered path, so an aliased
// registration is never reported as an available directory.
func classifyWorktreeAvailability(r *WorktreeRegistration) {
	r.Availability, r.HistoryCoverage = WorktreeIndeterminate, WorktreeHistoryUnknown
	info, err := os.Lstat(r.RegisteredPath)
	switch {
	case err == nil && info.IsDir():
		r.Availability, r.HistoryCoverage = WorktreeAvailable, WorktreeHistoryLocal
		if resolved, e := filepath.EvalSymlinks(r.RegisteredPath); e == nil {
			r.ResolvedPath = resolved
		}
	case err == nil:
		r.Availability = WorktreeNotDirectory
	case os.IsNotExist(err):
		r.Availability = WorktreeMissing
	}
}

// Registration input is parsed strictly. Malformed, ambiguous, duplicated or
// out-of-scope attributes are refused rather than guessed at, because these
// facts are load-bearing for the observation digest.
func parseWorktreeRegistrations(out []byte) ([]WorktreeRegistration, error) {
	if len(out) > maxWorktreeListBytes {
		return nil, fmt.Errorf("worktree registrations exceed %d bytes", maxWorktreeListBytes)
	}
	var registrations []WorktreeRegistration
	var current *WorktreeRegistration
	seen := map[string]bool{}
	attributes := map[string]bool{}
	closeRecord := func() error {
		if current == nil {
			return nil
		}
		if !current.Bare && current.Head == "" {
			return fmt.Errorf("worktree registration %q reports neither a HEAD nor bare", current.RegisteredPath)
		}
		if current.Bare && (current.Head != "" || current.Branch != "" || current.Detached) {
			return fmt.Errorf("worktree registration %q is both bare and checked out", current.RegisteredPath)
		}
		if current.Branch != "" && current.Detached {
			return fmt.Errorf("worktree registration %q is both detached and on a branch", current.RegisteredPath)
		}
		registrations = append(registrations, *current)
		current, attributes = nil, map[string]bool{}
		return nil
	}
	for _, field := range strings.Split(string(out), "\x00") {
		if field == "" {
			if err := closeRecord(); err != nil {
				return nil, err
			}
			continue
		}
		key, value, hasValue := strings.Cut(field, " ")
		if key == "worktree" {
			if current != nil {
				return nil, errors.New("worktree registration is not terminated before the next one")
			}
			if !hasValue || value == "" || !filepath.IsAbs(value) {
				return nil, fmt.Errorf("worktree registration path is empty or not absolute: %q", value)
			}
			if seen[value] {
				return nil, fmt.Errorf("duplicate worktree registration path: %q", value)
			}
			seen[value] = true
			current = &WorktreeRegistration{RegisteredPath: value}
			attributes = map[string]bool{"worktree": true}
			continue
		}
		if current == nil {
			return nil, fmt.Errorf("worktree attribute %q appears outside a registration", key)
		}
		if attributes[key] {
			return nil, fmt.Errorf("duplicate worktree attribute %q for %q", key, current.RegisteredPath)
		}
		attributes[key] = true
		switch key {
		case "HEAD":
			if !validWorktreeObjectName(value) {
				return nil, fmt.Errorf("worktree registration %q reports an invalid HEAD: %q", current.RegisteredPath, value)
			}
			current.Head = value
		case "branch":
			if !strings.HasPrefix(value, "refs/") {
				return nil, fmt.Errorf("worktree registration %q reports an unqualified branch: %q", current.RegisteredPath, value)
			}
			current.Branch = value
		case "detached", "bare":
			if hasValue {
				return nil, fmt.Errorf("worktree attribute %q for %q carries an unexpected value", key, current.RegisteredPath)
			}
			if key == "detached" {
				current.Detached = true
			} else {
				current.Bare = true
			}
		case "locked":
			current.Locked, current.LockReason = true, value
		case "prunable":
			current.Prunable, current.PrunableReason = true, value
		default:
			// Fail closed: an unrecognized attribute is a fact this stage cannot
			// place in a load-bearing digest, so it is refused, not dropped.
			return nil, fmt.Errorf("unsupported worktree attribute %q for %q", key, current.RegisteredPath)
		}
	}
	if current != nil {
		return nil, errors.New("trailing unterminated worktree registration")
	}
	if len(registrations) == 0 {
		return nil, errors.New("no worktree registration was reported for an existing repository")
	}
	if len(registrations) > maxWorktreeRegistration {
		return nil, fmt.Errorf("worktree registrations exceed %d entries", maxWorktreeRegistration)
	}
	return registrations, nil
}

func validWorktreeObjectName(name string) bool {
	if len(name) != 40 && len(name) != 64 {
		return false
	}
	for _, c := range name {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			return false
		}
	}
	return true
}
