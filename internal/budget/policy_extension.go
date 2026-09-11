package budget

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"path/filepath"
	"time"
)

// PolicyCeilings are absolute lifetime ceilings. An unchanged zero axis keeps
// its original unlimited meaning; an extension cannot turn a finite axis off.
type PolicyCeilings struct {
	Actions     int   `json:"actions"`
	CostMicros  int64 `json:"cost_micros"`
	WallClockNS int64 `json:"wall_clock_ns"`
}

type PolicyExtension struct {
	ID             string         `json:"id"`
	At             time.Time      `json:"at"`
	PreviousSHA256 string         `json:"previous_sha256"`
	Ceilings       PolicyCeilings `json:"ceilings"`
	Spent          int            `json:"spent"`
	ExposureMicros *int64         `json:"exposure_micros"`
	StartedAt      time.Time      `json:"started_at"`
	Reason         string         `json:"reason"`
}

type PolicyExtensionRequest struct {
	DecisionID, ExpectedPolicySHA256, Reason string
	Ceilings                                 PolicyCeilings
}

type PolicyStatus struct {
	Kind           Kind            `json:"kind"`
	Policy         json.RawMessage `json:"policy"`
	PolicySHA256   string          `json:"policy_sha256"`
	Spent          int             `json:"spent"`
	ExposureMicros *int64          `json:"exposure_micros"`
	StartedAt      time.Time       `json:"started_at"`
	LedgerDir      string          `json:"ledger_dir"`
}

type runtimePolicy interface {
	ceilings() PolicyCeilings
	originalPolicy() runtimePolicy
	withGrant(PolicyExtension) runtimePolicy
	extensionMetadata() (int, *PolicyCeilings, []PolicyExtension)
	policyKind() Kind
	validateBase() error
}

func (p LaunchPolicy) ceilings() PolicyCeilings {
	return PolicyCeilings{p.MaxLaunches, p.MaxCostMicros, p.WallClockMS * int64(time.Millisecond)}
}
func (p LaunchPolicy) originalPolicy() runtimePolicy {
	if p.Original != nil {
		p.MaxLaunches, p.MaxCostMicros, p.WallClockMS = p.Original.Actions, p.Original.CostMicros, p.Original.WallClockNS/int64(time.Millisecond)
	}
	p.Version, p.Original, p.Extensions = 1, nil, nil
	return p
}
func (p LaunchPolicy) withGrant(e PolicyExtension) runtimePolicy {
	base := p.originalPolicy().ceilings()
	p.Version, p.Original = 2, &base
	p.MaxLaunches, p.MaxCostMicros, p.WallClockMS = e.Ceilings.Actions, e.Ceilings.CostMicros, e.Ceilings.WallClockNS/int64(time.Millisecond)
	p.Extensions = append(append([]PolicyExtension(nil), p.Extensions...), e)
	return p
}
func (p LaunchPolicy) extensionMetadata() (int, *PolicyCeilings, []PolicyExtension) {
	return p.Version, p.Original, p.Extensions
}
func (p LaunchPolicy) policyKind() Kind { return Launch }

func (p StepPolicy) ceilings() PolicyCeilings {
	return PolicyCeilings{Actions: p.MaxSteps, WallClockNS: p.WallClockNS}
}
func (p StepPolicy) originalPolicy() runtimePolicy {
	if p.Original != nil {
		p.MaxSteps, p.WallClockNS = p.Original.Actions, p.Original.WallClockNS
	}
	p.Version, p.Original, p.Extensions = 1, nil, nil
	return p
}
func (p StepPolicy) withGrant(e PolicyExtension) runtimePolicy {
	base := p.originalPolicy().ceilings()
	p.Version, p.Original = 2, &base
	p.MaxSteps, p.WallClockNS = e.Ceilings.Actions, e.Ceilings.WallClockNS
	p.Extensions = append(append([]PolicyExtension(nil), p.Extensions...), e)
	return p
}
func (p StepPolicy) extensionMetadata() (int, *PolicyCeilings, []PolicyExtension) {
	return p.Version, p.Original, p.Extensions
}
func (p StepPolicy) policyKind() Kind { return DriverStep }

func policyDigest(p runtimePolicy) string {
	raw, _ := json.Marshal(p)
	h := sha256.Sum256(raw)
	return hex.EncodeToString(h[:])
}

func validPolicyCeilings(kind Kind, c PolicyCeilings) bool {
	return c.Actions >= 0 && c.CostMicros >= 0 && c.WallClockNS >= 0 &&
		(kind == Launch && c.WallClockNS%int64(time.Millisecond) == 0 || kind == DriverStep && c.CostMicros == 0)
}

func validateIncrease(old, next PolicyCeilings, spent int, exposure *int64, elapsed time.Duration) error {
	if old == next {
		return errors.New("an extension must increase at least one finite ceiling")
	}
	if next.Actions != old.Actions && (old.Actions == 0 || !finiteCycleMaximum(next.Actions) || next.Actions <= old.Actions || next.Actions <= spent) {
		return errors.New("new action ceiling must be finite and exceed both the prior ceiling and spent count")
	}
	if next.CostMicros != old.CostMicros {
		if exposure == nil {
			return fmt.Errorf("%w: reconcile exposure before increasing the monetary ceiling", ErrUnknownCost)
		}
		if old.CostMicros == 0 || next.CostMicros <= old.CostMicros || next.CostMicros <= *exposure || next.CostMicros == math.MaxInt64 {
			return errors.New("new monetary ceiling must be finite and exceed the prior ceiling and exposure")
		}
	}
	if next.WallClockNS != old.WallClockNS && (old.WallClockNS == 0 || next.WallClockNS <= old.WallClockNS || next.WallClockNS <= int64(elapsed) || next.WallClockNS == math.MaxInt64) {
		return errors.New("new wall-clock ceiling must be finite and exceed the prior ceiling and elapsed lifetime")
	}
	return nil
}

func validateRuntimePolicy(p runtimePolicy) error {
	if err := p.validateBase(); err != nil {
		return err
	}
	v, original, history := p.extensionMetadata()
	if v == 1 {
		if original != nil || history != nil {
			return errors.New("v1 policy cannot carry extension fields")
		}
		return nil
	}
	if v != 2 || original == nil || !validPolicyCeilings(p.policyKind(), *original) || len(history) == 0 || len(history) > maxCycleExtensions {
		return errors.New("invalid runtime policy extension history")
	}
	prior := p.originalPolicy()
	if err := prior.validateBase(); err != nil {
		return err
	}
	seen := map[string]bool{}
	var previousTime, started time.Time
	spent := 0
	for _, e := range history {
		if !validCycleDecision(e.ID, e.Reason, e.PreviousSHA256) || seen[e.ID] || e.At.IsZero() || e.StartedAt.IsZero() || e.At.Before(e.StartedAt) || e.At.Before(previousTime) || !started.IsZero() && e.StartedAt != started || e.Spent < spent || e.ExposureMicros != nil && *e.ExposureMicros < 0 || !validPolicyCeilings(p.policyKind(), e.Ceilings) || e.PreviousSHA256 != policyDigest(prior) {
			return errors.New("invalid runtime policy extension chain")
		}
		if err := validateIncrease(prior.ceilings(), e.Ceilings, e.Spent, e.ExposureMicros, e.At.Sub(e.StartedAt)); err != nil {
			return err
		}
		seen[e.ID], previousTime, started, spent = true, e.At, e.StartedAt, e.Spent
		prior = prior.withGrant(e)
	}
	if prior.ceilings() != p.ceilings() {
		return errors.New("effective policy ceilings differ from recorded decisions")
	}
	return nil
}

// Required fields, duplicate/case checks and bounded regular-file reads preserve
// v1 fail-closed behavior while permitting the explicitly versioned v2 history.
func readRuntimePolicy(path string, target any, required []string) error {
	data, err := readStepHistoryFile(path, 1<<20)
	if err != nil {
		return err
	}
	if err := validateHistoryJSON(json.NewDecoder(bytes.NewReader(data)), 0); err != nil {
		return err
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(target); err != nil {
		return err
	}
	if dec.Decode(new(any)) != io.EOF {
		return errors.New("trailing runtime policy data")
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	check := func(fields map[string]json.RawMessage, names []string) error {
		for _, name := range names {
			v, ok := fields[name]
			if !ok || bytes.Equal(bytes.TrimSpace(v), []byte("null")) && name != "reserve_micros" && name != "exposure_micros" {
				return errors.New("incomplete runtime policy object")
			}
		}
		return nil
	}
	if err := check(fields, required); err != nil {
		return err
	}
	var version int
	if err := json.Unmarshal(fields["version"], &version); err != nil {
		return err
	}
	_, hasOriginal := fields["original"]
	_, hasExtensions := fields["extensions"]
	if version == 1 {
		if hasOriginal || hasExtensions {
			return errors.New("v1 policy cannot carry extension fields")
		}
		return nil
	}
	if err := check(fields, []string{"original", "extensions"}); err != nil {
		return err
	}
	var original map[string]json.RawMessage
	if err := json.Unmarshal(fields["original"], &original); err != nil {
		return err
	}
	ceilings := []string{"actions", "cost_micros", "wall_clock_ns"}
	if err := check(original, ceilings); err != nil {
		return err
	}
	var history []map[string]json.RawMessage
	if err := json.Unmarshal(fields["extensions"], &history); err != nil {
		return err
	}
	for _, e := range history {
		if err := check(e, []string{"id", "at", "previous_sha256", "ceilings", "spent", "exposure_micros", "started_at", "reason"}); err != nil {
			return err
		}
		var fields map[string]json.RawMessage
		if err := json.Unmarshal(e["ceilings"], &fields); err != nil {
			return err
		}
		if err := check(fields, ceilings); err != nil {
			return err
		}
	}
	return nil
}

type runtimeBinding struct {
	policy runtimePolicy
	store  Store
}

func loadRuntimeBinding(ctx context.Context, root, idea string, kind Kind) (*runtimeBinding, error) {
	switch kind {
	case Launch:
		b, err := LoadLaunchBinding(ctx, root, idea)
		if err != nil {
			return nil, err
		}
		if b != nil {
			return &runtimeBinding{b.Policy, b.Store}, nil
		}
	case DriverStep:
		b, err := LoadStepBinding(ctx, root, idea)
		if err != nil {
			return nil, err
		}
		if b != nil {
			return &runtimeBinding{b.Policy, b.Store}, nil
		}
	default:
		return nil, errors.New("runtime policy kind must be launch or driver-step")
	}
	return nil, errors.New("initialize or migrate the original policy before inspecting or extending it")
}

func (b runtimeBinding) current() (runtimeBinding, error) {
	path := filepath.Join(filepath.Dir(b.store.Dir), "policy.json")
	var p runtimePolicy
	var err error
	if b.policy.policyKind() == Launch {
		p, err = readLaunchPolicy(path)
	} else {
		p, err = readStepPolicy(path)
	}
	if err != nil {
		return b, err
	}
	if policyDigest(b.policy.originalPolicy()) != policyDigest(p.originalPolicy()) {
		return b, errors.New("original runtime policy authority changed")
	}
	b.policy = p
	return b, nil
}

func (b runtimeBinding) inspect(ctx context.Context) (PolicyStatus, error) {
	s, err := b.store.Inspect(ctx)
	if err != nil {
		return PolicyStatus{}, err
	}
	spent := 0
	for _, e := range s.Entries {
		if e.Kind != b.policy.policyKind() {
			return PolicyStatus{}, errors.New("unexpected action kind in policy ledger")
		}
		spent++
	}
	_, _, history := b.policy.extensionMetadata()
	for _, e := range history {
		if e.Spent > spent || !e.StartedAt.Equal(s.StartedAt) {
			return PolicyStatus{}, errors.New("policy grant differs from retained charges or activation clock")
		}
	}
	var exposure *int64
	value, err := s.ExposureError()
	if err == nil {
		exposure = &value
	} else if !errors.Is(err, ErrUnknownCost) {
		return PolicyStatus{}, err
	}
	raw, err := json.Marshal(b.policy)
	if err != nil {
		return PolicyStatus{}, err
	}
	return PolicyStatus{b.policy.policyKind(), raw, policyDigest(b.policy), spent, exposure, s.StartedAt, b.store.Dir}, nil
}

func InspectRuntimeBudget(ctx context.Context, root, idea string, kind Kind) (PolicyStatus, error) {
	b, err := loadRuntimeBinding(ctx, root, idea, kind)
	if err != nil {
		return PolicyStatus{}, err
	}
	return b.inspect(ctx)
}

func ExtendRuntimeBudget(ctx context.Context, root, idea string, kind Kind, r PolicyExtensionRequest) (PolicyStatus, error) {
	if !validCycleDecision(r.DecisionID, r.Reason, r.ExpectedPolicySHA256) || !validPolicyCeilings(kind, r.Ceilings) {
		return PolicyStatus{}, errors.New("invalid finite runtime policy extension request")
	}
	b, err := loadRuntimeBinding(ctx, root, idea, kind)
	if err != nil {
		return PolicyStatus{}, err
	}
	return b.extend(ctx, r, writeSynced)
}

func (b runtimeBinding) extend(ctx context.Context, r PolicyExtensionRequest, persist func(string, []byte) error) (PolicyStatus, error) {
	release, err := AcquireResourceGuard(ctx, filepath.Dir(b.store.Dir))
	if err != nil {
		return PolicyStatus{}, err
	}
	defer release()
	b, err = b.current()
	if err != nil {
		return PolicyStatus{}, err
	}
	// Settlement and operator cost reconciliation use the ledger lock directly.
	// Hold that same lock through publication so the recorded exposure is current
	// at the grant, without rewriting even one byte of charged history.
	releaseLedger, err := lock(ctx, filepath.Join(b.store.Dir, "ledger.lock"))
	if err != nil {
		return PolicyStatus{}, err
	}
	defer releaseLedger()
	status, err := b.inspect(ctx)
	if err != nil {
		return PolicyStatus{}, err
	}
	_, _, history := b.policy.extensionMetadata()
	for _, e := range history {
		if e.ID == r.DecisionID {
			if e.PreviousSHA256 != r.ExpectedPolicySHA256 || e.Ceilings != r.Ceilings || e.Reason != r.Reason {
				return PolicyStatus{}, errors.New("conflicting runtime extension decision")
			}
			return status, nil
		}
	}
	if status.PolicySHA256 != r.ExpectedPolicySHA256 {
		return PolicyStatus{}, errors.New("runtime policy changed since inspection; inspect again")
	}
	if len(history) >= maxCycleExtensions {
		return PolicyStatus{}, errors.New("runtime extension history is full; preserve it for explicit migration")
	}
	now := time.Now().UTC()
	if b.store.now != nil {
		now = b.store.now().UTC()
	}
	if now.Before(status.StartedAt) || len(history) > 0 && now.Before(history[len(history)-1].At) {
		return PolicyStatus{}, ErrClockSkew
	}
	if err := validateIncrease(b.policy.ceilings(), r.Ceilings, status.Spent, status.ExposureMicros, now.Sub(status.StartedAt)); err != nil {
		return PolicyStatus{}, err
	}
	b.policy = b.policy.withGrant(PolicyExtension{r.DecisionID, now, r.ExpectedPolicySHA256, r.Ceilings, status.Spent, status.ExposureMicros, status.StartedAt, r.Reason})
	if err := validateRuntimePolicy(b.policy); err != nil {
		return PolicyStatus{}, err
	}
	raw, err := json.MarshalIndent(b.policy, "", "  ")
	if err != nil {
		return PolicyStatus{}, err
	}
	if len(raw) > 1<<20 {
		return PolicyStatus{}, errors.New("runtime policy exceeds size bound")
	}
	if err := persist(filepath.Join(filepath.Dir(b.store.Dir), "policy.json"), append(raw, '\n')); err != nil {
		return PolicyStatus{}, fmt.Errorf("persist runtime policy extension: %w", err)
	}
	return b.inspect(ctx)
}

func policyHistory(p runtimePolicy) []PolicyCeilings {
	_, _, history := p.extensionMetadata()
	out := []PolicyCeilings{p.originalPolicy().ceilings()}
	for _, e := range history {
		out = append(out, e.Ceilings)
	}
	return out
}

func (p LaunchPolicy) LimitHistory() []Limits {
	out := []Limits{}
	for _, c := range policyHistory(p) {
		out = append(out, Limits{Actions: map[Kind]int{Launch: c.Actions}, CostMicros: c.CostMicros, WallClock: time.Duration(c.WallClockNS)})
	}
	return out
}

func (p StepPolicy) accepts(steps int, wall time.Duration) bool {
	for _, c := range policyHistory(p) {
		if (steps == 0 || steps == c.Actions) && (wall == 0 || int64(wall) == c.WallClockNS) {
			return true
		}
	}
	return false
}

func (b *LaunchBinding) Reserve(ctx context.Context, id string) (Snapshot, error) {
	rb := runtimeBinding{b.Policy, b.Store}
	release, err := AcquireResourceGuard(ctx, filepath.Dir(b.Store.Dir))
	if err != nil {
		return Snapshot{}, err
	}
	defer release()
	rb, err = rb.current()
	if err != nil {
		return Snapshot{}, err
	}
	status, err := rb.inspect(ctx)
	if err != nil {
		return Snapshot{}, err
	}
	if err := rb.checkClock(status.StartedAt); err != nil {
		return Snapshot{}, err
	}
	p := rb.policy.(LaunchPolicy)
	return rb.store.Reserve(ctx, Request{ID: id, Kind: Launch, ReserveMicros: p.ReserveMicros}, p.Limits())
}

func (b *LaunchBinding) Inspect(ctx context.Context) (PolicyStatus, error) {
	rb, err := (runtimeBinding{b.Policy, b.Store}).current()
	if err != nil {
		return PolicyStatus{}, err
	}
	return rb.inspect(ctx)
}

func (b *StepBinding) Current() (*StepBinding, error) {
	rb, err := (runtimeBinding{b.Policy, b.Store}).current()
	if err != nil {
		return nil, err
	}
	return &StepBinding{rb.policy.(StepPolicy), rb.store}, nil
}

func (b *StepBinding) reserve(ctx context.Context, id string) (Snapshot, error) {
	release, err := AcquireResourceGuard(ctx, filepath.Dir(b.Store.Dir))
	if err != nil {
		return Snapshot{}, err
	}
	defer release()
	b, err = b.Current()
	if err != nil {
		return Snapshot{}, err
	}
	rb := runtimeBinding{b.Policy, b.Store}
	status, err := rb.inspect(ctx)
	if err != nil {
		return Snapshot{}, err
	}
	if err := rb.checkClock(status.StartedAt); err != nil {
		return Snapshot{}, err
	}
	zero := int64(0)
	return b.Store.Reserve(ctx, Request{ID: id, Kind: DriverStep, ReserveMicros: &zero}, b.Policy.limits())
}

func (b runtimeBinding) checkClock(started time.Time) error {
	now := time.Now().UTC()
	if b.store.now != nil {
		now = b.store.now().UTC()
	}
	_, _, history := b.policy.extensionMetadata()
	if now.Before(started) || len(history) > 0 && now.Before(history[len(history)-1].At) {
		return ErrClockSkew
	}
	return nil
}
