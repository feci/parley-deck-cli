package budget

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"
)

var ErrActionConflict = errors.New("action identity was reused with different input or accounting")

// ActionInput is supplied by a typed runtime boundary. Only input hashes are
// retained: no prompt, command, environment, transcript or artifact body.
// Runtime-input binds the caller's actual request; policy-only explicitly lacks
// that stronger binding and must not be presented as whole-operation identity.
type ActionInput struct {
	Operation   string `json:"operation"`
	Basis       string `json:"basis"`
	InputSHA256 string `json:"input_sha256"`
	RunID       string `json:"run_id,omitempty"`
}
type actionInputKey struct{}

func WithActionInput(ctx context.Context, input ActionInput) context.Context {
	return context.WithValue(ctx, actionInputKey{}, input)
}

// InheritActionInput retains the outer grouped operation for nested runners.
func InheritActionInput(ctx context.Context, input ActionInput) context.Context {
	if _, ok := ctx.Value(actionInputKey{}).(ActionInput); ok {
		return ctx
	}
	return WithActionInput(ctx, input)
}

// ActionInputDigest is a deterministic digest of typed caller input. Callers
// must provide the fields their operation actually reads, never a self verdict.
func ActionInputDigest(input any) (string, error) {
	raw, err := json.Marshal(input)
	if err != nil {
		return "", errors.New("typed action input cannot be encoded")
	}
	if len(raw) > 4<<20 {
		return "", errors.New("typed action input exceeds digest bound")
	}
	return key(string(raw)), nil
}

// ActionIdentity separates a logical operation from this one spent attempt.
// EntrySHA256 is the ledger key of the caller's original unique attempt ID.
// LimitsSHA256 freezes the effective reservation ceilings, not a later grant.
// ReservationSHA256 pins the original charge, including its ledger epoch. It is
// populated inside the reservation transaction, once its timestamp is known.
type ActionIdentity struct {
	Version           int         `json:"version"`
	Input             ActionInput `json:"input"`
	LogicalSHA256     string      `json:"logical_sha256"`
	EntrySHA256       string      `json:"entry_sha256"`
	LimitsSHA256      string      `json:"limits_sha256"`
	ReservationSHA256 string      `json:"reservation_sha256"`
}

func (a ActionIdentity) Digest() string { return migrationDigest(a) }
func actionIdentityDigest(a *ActionIdentity) string {
	if a == nil {
		return ""
	}
	return a.Digest()
}
func actionLogicalDigest(scope string, input ActionInput) string {
	return migrationDigest(struct {
		Scope string      `json:"scope"`
		Input ActionInput `json:"input"`
	}{scope, input})
}
func validActionLabel(s string, limit int, empty bool) bool {
	if len(s) > limit || s == "" && !empty {
		return false
	}
	for _, r := range s {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || strings.ContainsRune("-_.:/", r)) {
			return false
		}
	}
	return true
}
func validateActionInput(in ActionInput) error {
	if !validActionLabel(in.Operation, 80, false) || !validActionLabel(in.RunID, 200, true) || !originHash(in.InputSHA256) || in.Basis != "runtime-input" && in.Basis != "policy-only" && in.Basis != "launch-metadata" {
		return errors.New("action input must be typed, bounded and hash-bound with an explicit basis")
	}
	return nil
}
func originalChargeDigest(scope string, startedAt time.Time, entryKey string, entry Reservation) string {
	return migrationDigest(struct {
		Scope         string    `json:"scope"`
		StartedAt     time.Time `json:"started_at"`
		EntryKey      string    `json:"entry_key"`
		Kind          Kind      `json:"kind"`
		ReservedAt    time.Time `json:"reserved_at"`
		ReserveMicros *int64    `json:"reserve_micros"`
	}{scope, startedAt, entryKey, entry.Kind, entry.ReservedAt, entry.ReserveMicros})
}
func (a ActionIdentity) validate(scope, entryKey string, startedAt time.Time, entry Reservation) error {
	if e := validateActionInput(a.Input); e != nil {
		return e
	}
	if a.Version != 1 || !originHash(a.EntrySHA256) || a.EntrySHA256 != entryKey || !originHash(a.LimitsSHA256) || a.LogicalSHA256 != actionLogicalDigest(scope, a.Input) {
		return errors.New("action identity differs from ledger scope, entry or typed input")
	}
	if !originHash(a.ReservationSHA256) || a.ReservationSHA256 != originalChargeDigest(scope, startedAt, entryKey, entry) {
		return errors.New("action identity differs from its original reservation")
	}
	return nil
}
func actionIdentityFor(ctx context.Context, scope string, req Request, limits Limits) (*ActionIdentity, error) {
	in, ok := ctx.Value(actionInputKey{}).(ActionInput)
	if !ok {
		return nil, nil
	} // retained legacy API without invented semantics
	if e := validateActionInput(in); e != nil {
		return nil, e
	}
	return &ActionIdentity{Version: 1, Input: in, LogicalSHA256: actionLogicalDigest(scope, in), EntrySHA256: key(req.ID), LimitsSHA256: migrationDigest(limits)}, nil
}
func cloneActionIdentity(a *ActionIdentity) *ActionIdentity {
	if a == nil {
		return nil
	}
	copy := *a
	return &copy
}
func sameOriginalRequest(entry Reservation, req Request, action *ActionIdentity) bool {
	// The already-validated original charge retains its epoch and timestamp; a
	// duplicate lookup must not invent the timestamp of a hypothetical retry.
	original := cloneActionIdentity(entry.Action)
	if original != nil {
		original.ReservationSHA256 = ""
	}
	return entry.Kind == req.Kind && reflect.DeepEqual(entry.ReserveMicros, req.ReserveMicros) && reflect.DeepEqual(original, action)
}

// ActionReceipt reports the original charge. It deliberately contains no
// execution permission or operation-success assertion. A terminal cost is not
// an operation terminal, and reading a receipt cannot relaunch its child.
type ActionReceipt struct {
	Scope           string          `json:"scope"`
	EntryKey        string          `json:"entry_key"`
	LedgerStartedAt time.Time       `json:"ledger_started_at"`
	ReservedAt      time.Time       `json:"reserved_at"`
	Kind            Kind            `json:"kind"`
	ReserveMicros   *int64          `json:"reserve_micros"`
	Identity        *ActionIdentity `json:"identity,omitempty"`
	IdentitySHA256  string          `json:"identity_sha256,omitempty"`
	Binding         string          `json:"binding"`
	Permission      string          `json:"permission"`
	ExecutionStatus string          `json:"execution_status"`
}

func actionReceipt(state Snapshot, entryKey string) (ActionReceipt, error) {
	entry, ok := state.Entries[entryKey]
	if !ok {
		return ActionReceipt{}, errors.New("action entry is missing; absence is not an unused action")
	}
	binding := "legacy-unbound"
	if entry.Action != nil {
		if e := entry.Action.validate(state.Scope, entryKey, state.StartedAt, entry); e != nil {
			return ActionReceipt{}, e
		}
		binding = entry.Action.Input.Basis
	}
	return ActionReceipt{Scope: state.Scope, EntryKey: entryKey, LedgerStartedAt: state.StartedAt, ReservedAt: entry.ReservedAt, Kind: entry.Kind, ReserveMicros: copyInt(entry.ReserveMicros), Identity: cloneActionIdentity(entry.Action), IdentitySHA256: actionIdentityDigest(entry.Action), Binding: binding, Permission: "none", ExecutionStatus: "not-established"}, nil
}
func (s Store) InspectActions(ctx context.Context) ([]ActionReceipt, error) {
	state, e := s.Inspect(ctx)
	if e != nil {
		return nil, e
	}
	keys := make([]string, 0, len(state.Entries))
	for k := range state.Entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]ActionReceipt, 0, len(keys))
	for _, k := range keys {
		r, e := actionReceipt(state, k)
		if e != nil {
			return nil, e
		}
		out = append(out, r)
	}
	return out, nil
}

// ReplayAction is read-only exact accounting replay. The original identity must
// be present; a historical unbound entry cannot be retroactively certified.
// Later settlement and operator extensions do not rewrite this original receipt.
func (s Store) ReplayAction(ctx context.Context, entryKey, expectedIdentitySHA256 string) (ActionReceipt, error) {
	if !originHash(entryKey) || !originHash(expectedIdentitySHA256) {
		return ActionReceipt{}, errors.New("action replay requires exact entry and identity hashes")
	}
	state, e := s.Inspect(ctx)
	if e != nil {
		return ActionReceipt{}, e
	}
	r, e := actionReceipt(state, entryKey)
	if e != nil {
		return r, e
	}
	if r.Identity == nil || r.IdentitySHA256 != expectedIdentitySHA256 {
		return ActionReceipt{}, fmt.Errorf("%w: original action identity is missing or differs from the requested replay", ErrActionConflict)
	}
	return r, nil
}

func sessionActionInput(ctx context.Context) ActionInput {
	input, _ := ctx.Value(actionInputKey{}).(ActionInput)
	return input
}
func checkSessionActionInput(ctx context.Context, input ActionInput) error {
	if sessionActionInput(ctx) != input {
		return errors.New("nested accounting session changed its original logical action input")
	}
	return validateActionInput(input)
}
