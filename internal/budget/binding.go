package budget

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// LaunchPolicy is configured through an attended operator control. Values are
// frozen for this scope; zero count/cost/time ceilings mean unlimited. Unknown
// observed costs stay unknown even when ReserveMicros supplies a safe ceiling.
type LaunchPolicy struct {
	Version         int               `json:"version"`
	Scope           string            `json:"scope"`
	Idea            string            `json:"idea"`
	MaxLaunches     int               `json:"max_launches"`
	MaxCostMicros   int64             `json:"max_cost_micros"`
	WallClockMS     int64             `json:"wall_clock_ms"`
	ReserveMicros   *int64            `json:"reserve_micros"`
	Original        *PolicyCeilings   `json:"original,omitempty"`
	Extensions      []PolicyExtension `json:"extensions,omitempty"`
	MigrationSHA256 string            `json:"migration_sha256,omitempty"`
}

func (p LaunchPolicy) Limits() Limits {
	return Limits{Actions: map[Kind]int{Launch: p.MaxLaunches}, CostMicros: p.MaxCostMicros, WallClock: time.Duration(p.WallClockMS) * time.Millisecond}
}
func (p LaunchPolicy) validate() error { return validateRuntimePolicy(p) }

func (p LaunchPolicy) validateBase() error {
	if p.MigrationSHA256 != "" && !validCycleDecision("migration", "migration", p.MigrationSHA256) {
		return errors.New("invalid launch migration reference")
	}
	if (p.Version != 1 && p.Version != 2) || p.Scope == "" || p.MaxLaunches < 0 || p.MaxCostMicros < 0 || p.WallClockMS < 0 || p.WallClockMS > int64((1<<63-1)/time.Millisecond) || (p.ReserveMicros != nil && *p.ReserveMicros < 0) {
		return errors.New("invalid frozen launch budget policy")
	}
	if p.MaxCostMicros > 0 && (p.ReserveMicros == nil || *p.ReserveMicros > p.MaxCostMicros) {
		return errors.New("monetary launch policy requires a conservative reservation no larger than its cap")
	}
	return nil
}

type LaunchBinding struct {
	Policy LaunchPolicy
	Store  Store
}

// Declarable stat classes for a registered worktree whose working directory is
// gone. Any other stat error stays a hard refusal: an unreadable path is not an
// established absence, so it can never be declared away.
const (
	UnavailableMissing      = "stat-enoent"
	UnavailableNotDirectory = "stat-not-directory"
	// UnknownHistory is the only admissible history value for a declared row.
	// It is never a count and never zero.
	UnknownHistory = "unknown"
)

// UnavailableRoot records one registered worktree the operator declared
// unavailable for a single request. It is request-scoped: it is an argument to
// one decision, never repository state, and no ordinary bootstrap reads one.
type UnavailableRoot struct {
	Path        string `json:"path"`
	Observation string `json:"observation"`
	History     string `json:"history"`
}

// Canonical declaration form: absolute verbatim porcelain paths, sorted and
// deduplicated, so an exact replay cannot split on ordering, repetition or the
// difference between a nil and an empty list.
func normalizeDeclaredUnavailable(declared []string) ([]string, error) {
	if len(declared) == 0 {
		return nil, nil
	}
	if len(declared) > maxWorktreeRegistration {
		return nil, errors.New("too many declared-unavailable worktrees")
	}
	seen := map[string]bool{}
	canonical := make([]string, 0, len(declared))
	for _, path := range declared {
		if path == "" || !filepath.IsAbs(path) || strings.ContainsAny(path, "\x00\n\r") {
			return nil, fmt.Errorf("declared-unavailable worktree must be a verbatim absolute registered path: %q", path)
		}
		if seen[path] {
			continue
		}
		seen[path] = true
		canonical = append(canonical, path)
	}
	sort.Strings(canonical)
	return canonical, nil
}

// Declared paths are compared verbatim against the registered porcelain value.
// No normalization, symlink resolution or case folding widens the match.
func isDeclaredUnavailable(declared []string, worktree string) bool {
	for _, path := range declared {
		if path == worktree {
			return true
		}
	}
	return false
}

// Scope location is independent of run IDs and linked-worktree paths. A deck
// nested in a repository retains its prefix, so unrelated decks do not collide.
// Non-Git roots use local runtime storage and do not claim cross-root sharing.
func launchScope(ctx context.Context, root, idea string, inspectHistory bool) (dir, scope string, roots []string, err error) {
	dir, scope, roots, _, err = launchScopeDeclared(ctx, root, idea, inspectHistory, nil)
	return
}

// launchScopeDeclared is launchScope plus one request-scoped operator
// declaration. With a nil declaration it is exactly launchScope, which is what
// every ordinary Ensure/Configure caller passes: unavailable history stays a
// refusal for them. A declared path is retained as an explicitly unknown row
// and is excluded from roots, so it contributes no evidence, floor or count.
func launchScopeDeclared(ctx context.Context, root, idea string, inspectHistory bool, declared []string) (dir, scope string, roots []string, missing []UnavailableRoot, err error) {
	if idea == "." || idea == ".." || strings.ContainsAny(idea, "/\\\x00\n\r") {
		return "", "", nil, nil, errors.New("invalid launch budget idea")
	}
	declared, err = normalizeDeclaredUnavailable(declared)
	if err != nil {
		return
	}
	root, err = filepath.Abs(root)
	if err != nil {
		return
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return
	}
	common, gitErr := exec.CommandContext(ctx, "git", "-C", root, "rev-parse", "--path-format=absolute", "--git-common-dir").Output()
	base := filepath.Join(root, ".parley-runtime", "launch-budgets")
	prefix := ""
	roots = []string{root}
	// A declaration is only meaningful against enumerated Git registrations. It
	// is refused, not silently ignored, anywhere it could not be checked.
	if len(declared) > 0 && (!inspectHistory || gitErr != nil) {
		err = errors.New("declared-unavailable worktrees require enumerated Git worktree history")
		return
	}
	if gitErr == nil {
		gitDir := strings.TrimSpace(string(common))
		gitDir, err = filepath.EvalSymlinks(gitDir)
		if err != nil {
			return
		}
		base = filepath.Join(gitDir, "parley-launch-budgets")
		out, e := exec.CommandContext(ctx, "git", "-C", root, "rev-parse", "--show-prefix").Output()
		if e != nil {
			err = e
			return
		}
		prefix = strings.TrimSpace(string(out))
		if inspectHistory {
			out, e = exec.CommandContext(ctx, "git", "-C", root, "worktree", "list", "--porcelain", "-z").Output()
			if e != nil {
				err = e
				return
			}
			roots = nil
			registered := map[string]bool{}
			for _, field := range strings.Split(string(out), "\x00") {
				if strings.HasPrefix(field, "worktree ") {
					worktree := strings.TrimPrefix(field, "worktree ")
					registered[worktree] = true
					// Missing/unavailable worktrees are unknown history. A nested deck
					// that does not exist in an accessible worktree has no local files
					// to migrate; it still inherits the common policy if later created.
					info, e := os.Stat(worktree)
					observation := ""
					if e == nil && !info.IsDir() {
						observation = UnavailableNotDirectory
					} else if os.IsNotExist(e) {
						observation = UnavailableMissing
					}
					if isDeclaredUnavailable(declared, worktree) {
						// Bind the actual stat class, not the operator's word for it.
						// Non-ENOENT errors are never declarable, and a path that is
						// readable again refuses rather than migrating over evidence.
						if observation == "" {
							if e != nil {
								err = fmt.Errorf("historical worktree is unavailable: %s", worktree)
							} else {
								err = fmt.Errorf("declared-unavailable worktree is available again: %s", worktree)
							}
							return
						}
						missing = append(missing, UnavailableRoot{worktree, observation, UnknownHistory})
						continue
					}
					if e != nil || !info.IsDir() {
						err = fmt.Errorf("historical worktree is unavailable: %s", worktree)
						return
					}
					candidate := filepath.Join(worktree, filepath.FromSlash(prefix))
					if _, e := os.Stat(candidate); os.IsNotExist(e) {
						continue
					} else if e != nil {
						err = e
						return
					}
					roots = append(roots, candidate)
				}
			}
			// A declaration binds a retained registration. A path Git no longer
			// registers is a stale registry, not an unknown history row.
			for _, path := range declared {
				if !registered[path] {
					err = fmt.Errorf("declared-unavailable worktree is not a retained registration: %s", path)
					return
				}
			}
			sort.Slice(missing, func(a, b int) bool { return missing[a].Path < missing[b].Path })
		}
	} else {
		if ctx.Err() != nil {
			err = ctx.Err()
			return
		}
		// A broken Git repository must not fall back to a new local ledger.
		for ancestor := root; ; ancestor = filepath.Dir(ancestor) {
			if _, e := os.Lstat(filepath.Join(ancestor, ".git")); e == nil || !os.IsNotExist(e) {
				err = fmt.Errorf("cannot resolve shared Git budget scope: %w", gitErr)
				return
			}
			if filepath.Dir(ancestor) == ancestor {
				break
			}
		}
	}
	scope = "parley-launch/v1:" + key(prefix+"\x00"+idea)
	dir = filepath.Join(base, key(scope))
	return
}

// LoadLaunchBinding is read-only. An absent scope has no configured policy;
// an existing scope whose policy disappeared is an error, never unbudgeted.
func LoadLaunchBinding(ctx context.Context, root, idea string) (*LaunchBinding, error) {
	dir, scope, _, err := launchScope(ctx, root, idea, false)
	if err != nil {
		return nil, err
	}
	info, err := os.Lstat(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, errors.New("budget scope must be a real directory")
	}
	p, err := readLaunchPolicy(filepath.Join(dir, "policy.json"))
	if err != nil {
		return nil, fmt.Errorf("configured launch policy unavailable; preserve history and complete operator recovery: %w", err)
	}
	if p.Scope != scope || p.Idea != idea {
		return nil, errors.New("launch policy scope mismatch")
	}
	return &LaunchBinding{Policy: p, Store: Store{Dir: filepath.Join(dir, "ledger"), Scope: scope}}, nil
}

func readLaunchPolicy(path string) (LaunchPolicy, error) {
	p, err := readLaunchPolicyRaw(path)
	if err != nil {
		return p, err
	}
	if p.MigrationSHA256 == "" {
		for _, name := range []string{"migration.json", "migration-active"} {
			if _, err := os.Lstat(filepath.Join(filepath.Dir(path), name)); err == nil || !os.IsNotExist(err) {
				return p, errors.New("launch policy lost its migration reference")
			}
		}
	}
	return p, checkLaunchMigrationPolicy(filepath.Dir(path), p)
}

func readLaunchPolicyRaw(path string) (LaunchPolicy, error) {
	var p LaunchPolicy
	if err := readRuntimePolicy(path, &p, []string{"version", "scope", "idea", "max_launches", "max_cost_micros", "wall_clock_ms", "reserve_micros"}); err != nil {
		return p, err
	}
	return p, p.validate()
}

// ConfigureLaunchBudget freezes a first policy only. The CLI must establish
// attendance and obtain explicit values. No participant artifact supplies them.
// An exact replay preserves every charge and the original wall-clock origin.
// Policy extensions and legacy migration are separate controls, not resets here.
func ConfigureLaunchBudget(ctx context.Context, root, idea string, policy LaunchPolicy) (*LaunchBinding, error) {
	if policy.MigrationSHA256 != "" {
		return nil, errors.New("configure cannot manufacture a migration reference")
	}
	dir, scope, roots, err := launchScope(ctx, root, idea, true)
	if err != nil {
		return nil, err
	}
	policy.Version, policy.Scope, policy.Idea = 1, scope, idea
	if err := policy.validate(); err != nil {
		return nil, err
	}
	// Reject ordinary legacy state without leaving a partially enabled scope.
	if _, err := os.Lstat(dir); os.IsNotExist(err) {
		if err := refuseUnmigratedLaunches(roots, idea); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	// Publish the scope directory before examining history. Concurrent launches
	// now refuse its incomplete policy. Earlier launches already left a request.
	release, err := AcquireResourceGuard(ctx, dir)
	if err != nil {
		return nil, err
	}
	defer release()
	path := filepath.Join(dir, "policy.json")
	if prior, err := readLaunchPolicy(path); err == nil {
		policy.MigrationSHA256 = prior.MigrationSHA256
		a, _ := json.Marshal(prior.originalPolicy())
		b, _ := json.Marshal(policy)
		if !bytes.Equal(a, b) {
			return nil, errors.New("launch policy is frozen; changing it requires a separate explicit operator extension, not configure")
		}
		binding := &LaunchBinding{Policy: prior, Store: Store{Dir: filepath.Join(dir, "ledger"), Scope: scope}}
		if _, err := binding.Inspect(ctx); err != nil {
			return nil, err
		}
		return binding, nil
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	if _, err := os.Lstat(filepath.Join(dir, "migration.json")); err == nil || !os.IsNotExist(err) {
		return nil, errors.New("pending migration requires exact operator replay, not fresh configuration")
	}
	if err := refuseUnmigratedLaunches(roots, idea); err != nil {
		return nil, err
	}
	ledger := Store{Dir: filepath.Join(dir, "ledger"), Scope: scope}
	// A ledger with a missing policy is prior/partial state, never a new grant.
	if _, err := os.Lstat(filepath.Join(ledger.Dir, "ledger.json")); err == nil {
		return nil, errors.New("existing ledger has no frozen policy; preserve charges and restore the original policy before continuing")
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	if _, err := ledger.update(ctx, func(*Snapshot, time.Time) error { return nil }); err != nil {
		return nil, err
	}

	data, err := json.MarshalIndent(policy, "", "  ")
	if err != nil {
		return nil, err
	}
	if err := writeSynced(path, data); err != nil {
		return nil, err
	}
	return &LaunchBinding{Policy: policy, Store: ledger}, nil
}

func refuseUnmigratedLaunches(roots []string, idea string) error {
	for _, root := range roots {
		if _, err := os.Stat(root); err != nil {
			return fmt.Errorf("cannot inspect historical worktree for budget configuration: %w", err)
		}
		history := filepath.Join(root, ".parley-runtime", "invocations")
		for _, path := range []string{filepath.Dir(history), history} {
			info, err := os.Lstat(path)
			if os.IsNotExist(err) {
				continue
			}
			if err != nil {
				return err
			}
			if !info.IsDir() {
				return errors.New("historical invocation storage must use real directories; preserve history before budget configuration")
			}
		}
		entries, err := os.ReadDir(history)
		if os.IsNotExist(err) {
			entries, err = nil, nil
		}
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				return errors.New("historical invocation entry is not a real directory; accounting migration is required")
			}
			path := filepath.Join(history, entry.Name(), "requested.json")
			var record struct {
				Metadata *struct {
					Idea *string `json:"idea"`
				} `json:"metadata"`
			}
			data, err := readLockOrigin(path)
			if err != nil {
				return fmt.Errorf("historical invocation request unavailable; accounting migration is required: %w", err)
			}
			if validateHistoryJSON(json.NewDecoder(bytes.NewReader(data)), 0) != nil || json.Unmarshal(data, &record) != nil || record.Metadata == nil || record.Metadata.Idea == nil {
				return errors.New("unknown historical invocation accounting; migration is required before configuring a fresh policy")
			}
			if *record.Metadata.Idea == idea {
				return errors.New("existing invocations need charge migration before configuring a policy; refusing a zero-count start")
			}
		}
		if idea != "" {
			ideaDir := filepath.Join(root, "parley-deck", "ideas", idea)
			err := filepath.WalkDir(ideaDir, func(path string, entry os.DirEntry, err error) error {
				if os.IsNotExist(err) && path == ideaDir {
					return nil
				}
				if err != nil {
					return err
				}
				if !entry.IsDir() && filepath.Base(path) != "00-prompt.md" {
					return errors.New("existing idea artifacts need legacy accounting before configuring a fresh launch budget")
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Telemetry may contain explicit null observations and arrays, but duplicate or
// case-aliased fields cannot reassign a historical invocation to another idea.
func validateHistoryJSON(dec *json.Decoder, depth int) error {
	if depth > 24 {
		return errors.New("historical invocation JSON nesting exceeds limit")
	}
	token, err := dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := token.(json.Delim); ok {
		switch delim {
		case '{':
			seen := map[string]bool{}
			for dec.More() {
				keyToken, err := dec.Token()
				if err != nil {
					return err
				}
				name, ok := keyToken.(string)
				if !ok || name != strings.ToLower(name) || seen[name] {
					return errors.New("duplicate or aliased historical invocation field")
				}
				seen[name] = true
				if err := validateHistoryJSON(dec, depth+1); err != nil {
					return err
				}
			}
		case '[':
			for dec.More() {
				if err := validateHistoryJSON(dec, depth+1); err != nil {
					return err
				}
			}
		default:
			return errors.New("invalid historical invocation JSON")
		}
		if _, err := dec.Token(); err != nil {
			return err
		}
	} else if depth == 0 {
		return errors.New("historical invocation must be an object")
	}
	if depth == 0 {
		if dec.Decode(new(any)) != io.EOF {
			return errors.New("trailing historical invocation JSON")
		}
	}
	return nil
}
