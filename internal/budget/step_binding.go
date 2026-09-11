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
	"time"
)

// StepPolicy freezes lifetime driver-transition ceilings. It is activated by
// actual runtime configuration, never a participant's grant or prose. Zero
// means unlimited; a saved policy remains authoritative when flags are absent.
type StepPolicy struct {
	Version     int    `json:"version"`
	Scope       string `json:"scope"`
	Idea        string `json:"idea"`
	MaxSteps    int    `json:"max_steps"`
	WallClockNS int64  `json:"wall_clock_ns"`
}
type StepBinding struct {
	Policy StepPolicy
	Store  Store
}

func (p StepPolicy) limits() Limits {
	return Limits{Actions: map[Kind]int{DriverStep: p.MaxSteps}, WallClock: time.Duration(p.WallClockNS)}
}
func stepScope(ctx context.Context, root, idea string, history bool) (string, string, []string, error) {
	dir, scope, roots, err := launchScope(ctx, root, idea, history)
	return filepath.Join(filepath.Dir(dir), "steps-"+key(scope)), "parley-steps/v1:" + key(scope), roots, err
}
func readStepPolicy(path string) (StepPolicy, error) {
	var p StepPolicy
	data, err := readLockOrigin(path)
	if err != nil {
		return p, err
	}
	if err = checkJSON(json.NewDecoder(bytes.NewReader(data))); err != nil {
		return p, err
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err = dec.Decode(&p); err != nil {
		return p, err
	}
	if dec.Decode(new(any)) != io.EOF {
		return p, errors.New("trailing step policy data")
	}
	var fields map[string]json.RawMessage
	if err = json.Unmarshal(data, &fields); err != nil {
		return p, err
	}
	for _, name := range []string{"version", "scope", "idea", "max_steps", "wall_clock_ns"} {
		if _, ok := fields[name]; !ok {
			return p, errors.New("incomplete step policy")
		}
	}
	if p.Version != 1 || p.Scope == "" || p.Idea == "" || p.MaxSteps < 0 || p.WallClockNS < 0 {
		return p, errors.New("invalid step policy")
	}
	return p, nil
}
func LoadStepBinding(ctx context.Context, root, idea string) (*StepBinding, error) {
	if idea == "" {
		return nil, nil
	}
	dir, scope, _, err := stepScope(ctx, root, idea, false)
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
		return nil, errors.New("step scope must be a real directory")
	}
	p, err := readStepPolicy(filepath.Join(dir, "policy.json"))
	if err != nil {
		return nil, fmt.Errorf("saved step policy unavailable; preserve history and recover the original policy: %w", err)
	}
	if p.Scope != scope || p.Idea != idea {
		return nil, errors.New("step scope mismatch")
	}
	return &StepBinding{Policy: p, Store: Store{Dir: filepath.Join(dir, "ledger"), Scope: scope}}, nil
}

// EnsureStepBinding applies saved limits even when the next process omits flags.
// A differing nonzero runtime value is not permission for an extension.
func EnsureStepBinding(ctx context.Context, root, idea string, steps int, wall time.Duration) (*StepBinding, error) {
	if steps < 0 || wall < 0 {
		return nil, errors.New("invalid driver step limits")
	}
	old, err := LoadStepBinding(ctx, root, idea)
	if err != nil {
		return nil, err
	}
	if old != nil {
		if (steps != 0 && steps != old.Policy.MaxSteps) || (wall != 0 && int64(wall) != old.Policy.WallClockNS) {
			return nil, errors.New("driver step policy is frozen; an explicit operator extension is required")
		}
		return old, nil
	}
	if steps == 0 && wall == 0 {
		return nil, nil
	}
	if idea == "" {
		return nil, errors.New("driver step budget requires an idea identity")
	}
	dir, scope, roots, err := stepScope(ctx, root, idea, true)
	if err != nil {
		return nil, err
	}
	if err = refuseUnmigratedSteps(roots, idea); err != nil {
		return nil, err
	}
	release, err := AcquireResourceGuard(ctx, dir)
	if err != nil {
		return nil, err
	}
	defer release()
	if old, err = LoadStepBinding(ctx, root, idea); err == nil && old != nil {
		if old.Policy.MaxSteps != steps || old.Policy.WallClockNS != int64(wall) {
			return nil, errors.New("conflicting first step policy")
		}
		return old, nil
	}
	// Missing policy in our newly pinned scope is expected only before the first
	// ledger exists. Never use configuration to reset retained/partial charges.
	policyPath := filepath.Join(dir, "policy.json")
	if _, e := os.Lstat(policyPath); e == nil || !os.IsNotExist(e) {
		return nil, errors.New("step policy exists but cannot be read")
	}
	if err = refuseUnmigratedSteps(roots, idea); err != nil {
		return nil, err
	}
	b := &StepBinding{Policy: StepPolicy{Version: 1, Scope: scope, Idea: idea, MaxSteps: steps, WallClockNS: int64(wall)}, Store: Store{Dir: filepath.Join(dir, "ledger"), Scope: scope}}
	if _, e := os.Lstat(filepath.Join(b.Store.Dir, "ledger.json")); e == nil || !os.IsNotExist(e) {
		return nil, errors.New("step ledger without policy requires recovery, not fresh activation")
	}
	if _, err = b.Store.update(ctx, func(*Snapshot, time.Time) error { return nil }); err != nil {
		return nil, err
	}
	data, err := json.MarshalIndent(b.Policy, "", "  ")
	if err != nil {
		return nil, err
	}
	if err = writeSynced(policyPath, data); err != nil {
		return nil, err
	}
	return b, nil
}

func (b *StepBinding) Check(ctx context.Context) (Snapshot, error) {
	state, err := b.Store.Inspect(ctx)
	if err != nil {
		return state, err
	}
	if err := b.checkTime(state); err != nil {
		return state, err
	}
	if b.Policy.MaxSteps > 0 && StepCount(state) >= b.Policy.MaxSteps {
		return state, fmt.Errorf("%w: lifetime driver steps", ErrLimit)
	}
	return state, nil
}

func (b *StepBinding) checkTime(state Snapshot) error {
	now := time.Now().UTC()
	if now.Before(state.StartedAt) {
		return ErrClockSkew
	}
	if b.Policy.WallClockNS > 0 && now.Sub(state.StartedAt) >= time.Duration(b.Policy.WallClockNS) {
		return fmt.Errorf("%w: lifetime wall clock", ErrLimit)
	}
	return nil
}

// StepCount includes failed and unfinished attempts, not only successful work.
func StepCount(s Snapshot) int {
	n := 0
	for _, e := range s.Entries {
		if e.Kind == DriverStep {
			n++
		}
	}
	return n
}
