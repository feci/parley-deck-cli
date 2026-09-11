package budget

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

const maxCycleExtensions = 128

// CycleExtension is an operator decision, never a provider observation or an
// action reservation. Maximum is absolute, finite and strictly increasing.
type CycleExtension struct {
	ID             string    `json:"id"`
	At             time.Time `json:"at"`
	PreviousSHA256 string    `json:"previous_sha256"`
	Maximum        int       `json:"maximum"`
	Spent          int       `json:"spent"`
	Reason         string    `json:"reason"`
}

type CycleExtensionRequest struct {
	DecisionID, ExpectedPolicySHA256, Reason string
	Maximum                                  int
}

type CycleStatus struct {
	Policy       CyclePolicy `json:"policy"`
	PolicySHA256 string      `json:"policy_sha256"`
	Spent        int         `json:"spent"`
	StartedAt    time.Time   `json:"started_at"`
	LedgerDir    string      `json:"ledger_dir"`
}

func (p CyclePolicy) InitialMaximum() int {
	if p.OriginalMaximum != nil {
		return *p.OriginalMaximum
	}
	return p.Maximum
}

func cloneCyclePolicy(p CyclePolicy) CyclePolicy {
	if p.OriginalMaximum != nil {
		n := *p.OriginalMaximum
		p.OriginalMaximum = &n
	}
	p.Extensions = append([]CycleExtension(nil), p.Extensions...)
	return p
}

func sameCycleAuthority(a, b CyclePolicy) bool {
	return a.Scope == b.Scope && a.Idea == b.Idea && a.IdeaPath == b.IdeaPath && a.Kind == b.Kind && a.Carried == b.Carried && a.InitialMaximum() == b.InitialMaximum()
}

// CyclePolicyDigest addresses canonical JSON, not whitespace in a policy file.
func CyclePolicyDigest(p CyclePolicy) string {
	data, _ := json.Marshal(p)
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

func validCycleDecision(id, reason, digest string) bool {
	decoded, err := hex.DecodeString(digest)
	return strings.TrimSpace(id) != "" && len(id) <= 128 && strings.TrimSpace(reason) != "" && len(reason) <= 1024 && err == nil && len(decoded) == 32 && digest == strings.ToLower(digest)
}

func finiteCycleMaximum(n int) bool {
	// A cross-review maximum needs a representable following round ordinal.
	return n > 0 && n < int(^uint(0)>>1)
}

func validateCycleExtensions(p CyclePolicy, fields map[string]json.RawMessage) error {
	_, originalPresent := fields["original_maximum"]
	_, extensionsPresent := fields["extensions"]
	if p.Version == 1 {
		if originalPresent || extensionsPresent {
			return errors.New("v1 cycle policy cannot carry extension fields")
		}
		return nil
	}
	if p.Version != 2 || p.OriginalMaximum == nil || *p.OriginalMaximum <= 0 || len(p.Extensions) == 0 || len(p.Extensions) > maxCycleExtensions {
		return errors.New("invalid extended cycle policy")
	}
	var raw []map[string]json.RawMessage
	if err := json.Unmarshal(fields["extensions"], &raw); err != nil {
		return err
	}
	for _, fields := range raw {
		for _, name := range []string{"id", "at", "previous_sha256", "maximum", "spent", "reason"} {
			if value, ok := fields[name]; !ok || strings.TrimSpace(string(value)) == "null" {
				return errors.New("incomplete cycle extension")
			}
		}
	}
	prior := p
	prior.Version = 1
	prior.Maximum = *p.OriginalMaximum
	prior.OriginalMaximum = nil
	prior.Extensions = nil
	seen := map[string]bool{}
	previousTime := time.Time{}
	spent := p.Carried
	for _, e := range p.Extensions {
		if !validCycleDecision(e.ID, e.Reason, e.PreviousSHA256) || !finiteCycleMaximum(e.Maximum) || seen[e.ID] || e.At.IsZero() || e.At.Before(previousTime) || e.Spent < spent || e.Maximum <= prior.Maximum || e.Maximum <= e.Spent || e.PreviousSHA256 != CyclePolicyDigest(prior) {
			return errors.New("invalid cycle extension chain")
		}
		seen[e.ID] = true
		previousTime = e.At
		spent = e.Spent
		base := p.InitialMaximum()
		prior.Version = 2
		prior.OriginalMaximum = &base
		prior.Maximum = e.Maximum
		prior.Extensions = append(prior.Extensions, e)
	}
	if prior.Maximum != p.Maximum {
		return errors.New("cycle maximum differs from the recorded operator grant")
	}
	return nil
}

func (b *CycleBinding) current() (*CycleBinding, error) {
	p, err := readCyclePolicy(filepath.Join(filepath.Dir(b.Store.Dir), "policy.json"))
	if err != nil {
		return nil, err
	}
	if !sameCycleAuthority(b.Policy, p) {
		return nil, errors.New("frozen cycle authority changed")
	}
	return &CycleBinding{Policy: p, Store: b.Store}, nil
}

func (b *CycleBinding) Inspect(ctx context.Context) (CycleStatus, error) {
	state, err := b.Store.Inspect(ctx)
	if err != nil {
		return CycleStatus{}, err
	}
	spent := b.Count(state)
	for _, e := range b.Policy.Extensions {
		if e.Spent > spent || e.At.Before(state.StartedAt) {
			return CycleStatus{}, errors.New("cycle grant differs from retained charge history")
		}
	}
	return CycleStatus{Policy: cloneCyclePolicy(b.Policy), PolicySHA256: CyclePolicyDigest(b.Policy), Spent: spent, StartedAt: state.StartedAt, LedgerDir: b.Store.Dir}, nil
}

// InspectCycleBudget does not initialize policy, a ledger or its locks.
func InspectCycleBudget(ctx context.Context, root, idea string, kind Kind) (CycleStatus, error) {
	b, err := LoadCycleBinding(ctx, root, idea, kind)
	if err != nil {
		return CycleStatus{}, err
	}
	if b == nil {
		return CycleStatus{}, errors.New("cycle policy is not initialized")
	}
	return b.Inspect(ctx)
}

// ExtendCycleBudget is for the attended CLI/operator control. The caller must
// obtain the exact decision; this API is not authentication of a human identity.
func ExtendCycleBudget(ctx context.Context, root, idea string, kind Kind, r CycleExtensionRequest) (CycleStatus, error) {
	if !validCycleDecision(r.DecisionID, r.Reason, r.ExpectedPolicySHA256) || !finiteCycleMaximum(r.Maximum) {
		return CycleStatus{}, errors.New("invalid finite cycle extension request")
	}
	b, err := LoadCycleBinding(ctx, root, idea, kind)
	if err != nil {
		return CycleStatus{}, err
	}
	if b == nil {
		return CycleStatus{}, errors.New("initialize or migrate the original cycle policy before extending it")
	}
	return b.extend(ctx, r, writeSynced)
}

func (b *CycleBinding) extend(ctx context.Context, r CycleExtensionRequest, persist func(string, []byte) error) (CycleStatus, error) {
	release, err := AcquireResourceGuard(ctx, filepath.Dir(b.Store.Dir))
	if err != nil {
		return CycleStatus{}, err
	}
	defer release()
	b, err = b.current()
	if err != nil {
		return CycleStatus{}, err
	}
	status, err := b.Inspect(ctx)
	if err != nil {
		return CycleStatus{}, err
	}
	for _, old := range b.Policy.Extensions {
		if old.ID == r.DecisionID {
			if old.PreviousSHA256 != r.ExpectedPolicySHA256 || old.Maximum != r.Maximum || old.Reason != r.Reason {
				return CycleStatus{}, errors.New("conflicting cycle extension decision")
			}
			return status, nil
		}
	}
	if status.PolicySHA256 != r.ExpectedPolicySHA256 {
		return CycleStatus{}, errors.New("cycle policy changed since inspection; inspect the current policy before extending")
	}
	if b.Policy.InitialMaximum() == 0 {
		return CycleStatus{}, errors.New("a cycle extension cannot enable a skipped protocol phase")
	}
	if r.Maximum <= b.Policy.Maximum || r.Maximum <= status.Spent {
		return CycleStatus{}, errors.New("new finite maximum must exceed both the current maximum and spent count")
	}
	if len(b.Policy.Extensions) >= maxCycleExtensions {
		return CycleStatus{}, errors.New("cycle extension history is full; preserve it for explicit migration")
	}
	now := time.Now().UTC()
	if b.Store.now != nil {
		now = b.Store.now().UTC()
	}
	if now.Before(status.StartedAt) || len(b.Policy.Extensions) > 0 && now.Before(b.Policy.Extensions[len(b.Policy.Extensions)-1].At) {
		return CycleStatus{}, ErrClockSkew
	}
	base := b.Policy.InitialMaximum()
	b.Policy.Version = 2
	b.Policy.OriginalMaximum = &base
	b.Policy.Maximum = r.Maximum
	b.Policy.Extensions = append(b.Policy.Extensions, CycleExtension{ID: r.DecisionID, At: now, PreviousSHA256: r.ExpectedPolicySHA256, Maximum: r.Maximum, Spent: status.Spent, Reason: r.Reason})
	data, err := json.MarshalIndent(b.Policy, "", "  ")
	if err != nil {
		return CycleStatus{}, err
	}
	if len(data) > 1<<20 {
		return CycleStatus{}, errors.New("cycle policy exceeds size bound")
	}
	if err := persist(filepath.Join(filepath.Dir(b.Store.Dir), "policy.json"), append(data, '\n')); err != nil {
		return CycleStatus{}, fmt.Errorf("persist operator cycle grant: %w", err)
	}
	return b.Inspect(ctx)
}
