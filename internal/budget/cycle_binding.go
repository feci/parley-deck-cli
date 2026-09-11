package budget

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
)

// CyclePolicy preserves the protocol cap and known pre-upgrade charged count.
// Unlike launch/step policy, a zero cycle cap forbids that protocol operation
// (the fast track has no cross-review rounds). Policy is runtime state, never
// a grant taken from participant-authored headings or extension frontmatter.
type CyclePolicy struct {
	Version  int    `json:"version"`
	Scope    string `json:"scope"`
	Idea     string `json:"idea"`
	IdeaPath string `json:"idea_path"`
	Kind     Kind   `json:"kind"`
	Maximum  int    `json:"maximum"`
	Carried  int    `json:"carried"`
}

type CycleBinding struct {
	Policy CyclePolicy
	Store  Store
}

func cycleScope(ctx context.Context, root, idea string, kind Kind, history bool) (string, string, []string, error) {
	if idea == "" || (kind != Fixup && kind != CrossReview) {
		return "", "", nil, errors.New("cycle budget requires an idea and a protocol cycle kind")
	}
	dir, scope, roots, err := launchScope(ctx, root, idea, history)
	scope = "parley-cycles/v1:" + string(kind) + ":" + key(scope)
	return filepath.Join(filepath.Dir(dir), "cycles-"+key(scope)), scope, roots, err
}

func readCyclePolicy(path string) (CyclePolicy, error) {
	var p CyclePolicy
	data, err := readLockOrigin(path)
	if err != nil {
		return p, err
	}
	if err := validateHistoryJSON(json.NewDecoder(bytes.NewReader(data)), 0); err != nil {
		return p, err
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&p); err != nil {
		return p, err
	}
	if dec.Decode(new(any)) != io.EOF {
		return p, errors.New("trailing cycle policy data")
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return p, err
	}
	for _, name := range []string{"version", "scope", "idea", "idea_path", "kind", "maximum", "carried"} {
		if raw, ok := fields[name]; !ok || bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
			return p, errors.New("incomplete cycle policy")
		}
	}
	if p.Version != 1 || p.Scope == "" || p.Idea == "" || !strings.HasPrefix(p.IdeaPath, "parley-deck/") || strings.Contains(p.IdeaPath, "\\") || filepath.IsAbs(p.IdeaPath) || filepath.ToSlash(filepath.Clean(p.IdeaPath)) != p.IdeaPath || (p.Kind != Fixup && p.Kind != CrossReview) || p.Maximum < 0 || p.Carried < 0 {
		return p, errors.New("invalid cycle policy")
	}
	return p, nil
}

func LoadCycleBinding(ctx context.Context, root, idea string, kind Kind) (*CycleBinding, error) {
	dir, scope, _, err := cycleScope(ctx, root, idea, kind, false)
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
		return nil, errors.New("cycle budget scope must be a real directory")
	}
	p, err := readCyclePolicy(filepath.Join(dir, "policy.json"))
	if err != nil {
		return nil, fmt.Errorf("saved cycle policy unavailable; preserve charges and recover the original policy: %w", err)
	}
	if p.Scope != scope || p.Idea != idea || p.Kind != kind {
		return nil, errors.New("cycle budget scope mismatch")
	}
	return &CycleBinding{Policy: p, Store: Store{Dir: filepath.Join(dir, "ledger"), Scope: scope}}, nil
}

// EnsureCycleBinding accepts only the known current driver cursor/marker floor.
// Unmigrated invocation history or other charged runs must be reconciled by an
// operator; the caller's current cursor is not permission to forget those runs.
func EnsureCycleBinding(ctx context.Context, root, idea string, kind Kind, maximum, carried int, currentRun, ideaDir string) (*CycleBinding, error) {
	if ideaDir == "" {
		ideaDir = filepath.Join(root, "parley-deck", "ideas", idea)
	}
	relative, err := filepath.Rel(root, ideaDir)
	if err != nil || filepath.IsAbs(relative) || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return nil, errors.New("cycle idea path escapes the runtime origin")
	}
	relative = filepath.ToSlash(relative)
	if !strings.HasPrefix(relative, "parley-deck/") {
		return nil, errors.New("cycle idea must be inside the canonical deck")
	}
	if maximum < 0 || carried < 0 {
		return nil, errors.New("invalid cycle ceilings or carried charges")
	}
	old, err := LoadCycleBinding(ctx, root, idea, kind)
	if err != nil {
		return nil, err
	}
	if old != nil {
		if old.Policy.IdeaPath != relative {
			return nil, errors.New("cycle idea path differs from frozen scope")
		}
		if old.Policy.Maximum != maximum {
			return nil, errors.New("cycle ceiling is frozen; a changed flag or track is not an operator extension")
		}
		state, err := old.Store.Inspect(ctx)
		if err != nil {
			return nil, err
		}
		if old.Count(state) < carried {
			return nil, errors.New("driver history exceeds the shared cycle count; explicit accounting reconciliation is required")
		}
		return old, nil
	}
	dir, scope, roots, err := cycleScope(ctx, root, idea, kind, true)
	if err != nil {
		return nil, err
	}
	if err := refuseUnmigratedCycles(roots, idea, relative, kind, carried, currentRun); err != nil {
		return nil, err
	}
	release, err := AcquireResourceGuard(ctx, dir)
	if err != nil {
		return nil, err
	}
	defer release()
	path := filepath.Join(dir, "policy.json")
	if p, err := readCyclePolicy(path); err == nil {
		if p.Scope != scope || p.Idea != idea || p.IdeaPath != relative || p.Kind != kind || p.Maximum != maximum || p.Carried != carried {
			return nil, errors.New("conflicting first cycle policy")
		}
		b := &CycleBinding{Policy: p, Store: Store{Dir: filepath.Join(dir, "ledger"), Scope: scope}}
		if _, err := b.Store.Inspect(ctx); err != nil {
			return nil, err
		}
		return b, nil
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	if err := refuseUnmigratedCycles(roots, idea, relative, kind, carried, currentRun); err != nil {
		return nil, err
	}
	b := &CycleBinding{Policy: CyclePolicy{Version: 1, Scope: scope, Idea: idea, IdeaPath: relative, Kind: kind, Maximum: maximum, Carried: carried}, Store: Store{Dir: filepath.Join(dir, "ledger"), Scope: scope}}
	if _, err := os.Lstat(filepath.Join(b.Store.Dir, "ledger.json")); err == nil || !os.IsNotExist(err) {
		return nil, errors.New("cycle ledger without policy requires recovery")
	}
	if _, err := b.Store.update(ctx, func(*Snapshot, time.Time) error { return nil }); err != nil {
		return nil, err
	}
	data, err := json.MarshalIndent(b.Policy, "", "  ")
	if err != nil {
		return nil, err
	}
	if err := writeSynced(path, data); err != nil {
		return nil, err
	}
	return b, nil
}

func (b *CycleBinding) Count(s Snapshot) int {
	n := b.Policy.Carried
	for _, entry := range s.Entries {
		if entry.Kind == b.Policy.Kind {
			if n == int(^uint(0)>>1) {
				return n
			}
			n++
		}
	}
	return n
}

func (b *CycleBinding) Reserve(ctx context.Context, id string) (int, error) {
	remaining := b.Policy.Maximum - b.Policy.Carried
	limits := Limits{Actions: map[Kind]int{}, Denied: map[Kind]bool{}}
	if remaining <= 0 {
		limits.Denied[b.Policy.Kind] = true
	} else {
		limits.Actions[b.Policy.Kind] = remaining
	}
	zero := int64(0)
	state, err := b.Store.Reserve(ctx, Request{ID: id, Kind: b.Policy.Kind, ReserveMicros: &zero}, limits)
	if err != nil {
		return 0, err
	}
	return b.Count(state), nil
}
