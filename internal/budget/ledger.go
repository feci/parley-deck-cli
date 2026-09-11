// Package budget owns durable, precharged action reservations. Caller identities
// and limits are trusted runtime input, not authentication of a human operator.
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
	"os"
	"path/filepath"
	"strings"
	"time"

	"parley-deck-cli/internal/fsutil"
)

var (
	ErrReserved        = errors.New("action already reserved: reconcile it before any new execution")
	ErrLimit           = errors.New("budget exhausted")
	ErrUnknownCost     = errors.New("monetary exposure is unknown")
	ErrCostOverflow    = errors.New("monetary exposure overflows microdollars")
	ErrClockSkew       = errors.New("budget clock moved backwards")
	ErrLockContention  = errors.New("budget lock contention exhausted the acquisition deadline")
	ErrSnapshotChanged = errors.New("budget ledger was replaced during inspection; retry the read")
)

type Kind string

const (
	Launch      Kind = "launch"
	DriverStep  Kind = "driver-step"
	Fixup       Kind = "fixup"
	CrossReview Kind = "cross-review"
)

type Limits struct {
	// Lifetime spent-action totals; zero/absent means unlimited. Denied can
	// forbid a kind without changing the historical meaning of zero.
	Actions    map[Kind]int
	Denied     map[Kind]bool
	CostMicros int64
	// Refuse unknown new reservations even without a monetary cap. A provider
	// may still report unknown terminal cost; this does not invent a price.
	RequireKnownCost bool
	// Elapsed wall time from the first reservation, including pauses/resumes.
	WallClock time.Duration
}

type Request struct {
	ID   string
	Kind Kind
	// A conservative maximum, not a guessed price. nil is unknown exposure.
	ReserveMicros *int64
}

type Reservation struct {
	Kind          Kind      `json:"kind"`
	ReservedAt    time.Time `json:"reserved_at"`
	ReserveMicros *int64    `json:"reserve_micros"`
	Settled       bool      `json:"settled"`
	ActualMicros  *int64    `json:"actual_micros"`
	// Operator ceilings preserve unknown observed cost. They do not become
	// provider usage, erase a charge, or authorize repeating the action.
	Reconciliations []Reconciliation `json:"reconciliations,omitempty"`
}

type Reconciliation struct {
	ID            string    `json:"id"`
	At            time.Time `json:"at"`
	CeilingMicros int64     `json:"ceiling_micros"`
	Reason        string    `json:"reason"`
}

type Snapshot struct {
	Schema    int                    `json:"schema"`
	Scope     string                 `json:"scope"`
	StartedAt time.Time              `json:"started_at"`
	Entries   map[string]Reservation `json:"entries"`
}

type Store struct {
	Dir, Scope string
	now        func() time.Time
	persist    func(string, []byte) error
}

// Inspect reads an existing atomically published ledger without creating locks,
// origins or directories. It returns a complete point-in-time snapshot, which
// may precede a concurrent replacement. The maps have no retained owner.
func (s Store) Inspect(ctx context.Context) (Snapshot, error) {
	if s.Dir == "" || s.Scope == "" {
		return Snapshot{}, errors.New("budget directory and scope are required")
	}
	path := filepath.Join(s.Dir, "ledger.json")
	state, err := readStableSnapshot(ctx, path, read)
	if err != nil {
		return Snapshot{}, err
	}
	if state.Scope != s.Scope {
		return Snapshot{}, errors.New("budget scope mismatch")
	}
	return state, nil
}

// Retry only an observed atomic replacement, never malformed data or a missing
// ledger. A busy writer cannot turn this diagnostic operation into an unbounded
// wait; the caller can retry ErrSnapshotChanged after five unsuccessful reads.
func readStableSnapshot(ctx context.Context, path string, readFile func(string) (Snapshot, error)) (Snapshot, error) {
	for attempt := 0; attempt < 5; attempt++ {
		if err := ctx.Err(); err != nil {
			return Snapshot{}, err
		}
		state, err := readFile(path)
		if !errors.Is(err, ErrSnapshotChanged) {
			return state, err
		}
	}
	return Snapshot{}, ErrSnapshotChanged
}

func validKind(k Kind) bool { return k == Launch || k == DriverStep || k == Fixup || k == CrossReview }
func key(id string) string  { sum := sha256.Sum256([]byte(id)); return hex.EncodeToString(sum[:]) }

// Reserve commits before returning permission for exactly one new action. A
// repeated ID is never permission to repeat work. A retry needs a distinct ID.
// Failed/crashed attempts stay charged; status and signoffs do not use Fixup.
func (s Store) Reserve(ctx context.Context, req Request, limits Limits) (Snapshot, error) {
	if req.ID == "" || !validKind(req.Kind) || limits.CostMicros < 0 || limits.WallClock < 0 || (req.ReserveMicros != nil && *req.ReserveMicros < 0) {
		return Snapshot{}, errors.New("invalid reservation or limits")
	}
	for kind, n := range limits.Actions {
		if !validKind(kind) || n < 0 {
			return Snapshot{}, errors.New("invalid action limit")
		}
	}
	for kind := range limits.Denied {
		if !validKind(kind) {
			return Snapshot{}, errors.New("invalid denied action kind")
		}
	}
	return s.update(ctx, func(state *Snapshot, now time.Time) error {
		id := key(req.ID)
		if _, exists := state.Entries[id]; exists {
			return ErrReserved
		}
		if limits.Denied[req.Kind] {
			return fmt.Errorf("%w: %s is forbidden", ErrLimit, req.Kind)
		}
		if limits.WallClock > 0 && now.Sub(state.StartedAt) >= limits.WallClock {
			return fmt.Errorf("%w: wall clock", ErrLimit)
		}
		count := 0
		for _, entry := range state.Entries {
			if entry.Kind == req.Kind {
				count++
			}
		}
		if ceiling := limits.Actions[req.Kind]; ceiling > 0 && count >= ceiling {
			return fmt.Errorf("%w: %s", ErrLimit, req.Kind)
		}
		if (limits.CostMicros > 0 || limits.RequireKnownCost) && req.ReserveMicros == nil {
			return ErrUnknownCost
		}
		if limits.CostMicros > 0 {
			exposure, err := state.ExposureError()
			if err != nil {
				return err
			}
			if exposure > limits.CostMicros || *req.ReserveMicros > limits.CostMicros-exposure {
				return fmt.Errorf("%w: monetary ceiling", ErrLimit)
			}
		}
		state.Entries[id] = Reservation{Kind: req.Kind, ReservedAt: now, ReserveMicros: copyInt(req.ReserveMicros)}
		return nil
	})
}

// ReconcileUnknown records an explicit operator-supplied conservative ceiling
// for an unknown observation. The calling operator control MUST obtain the
// actual decision; this API is not human authentication and must never consume
// participant-authored frontmatter grants. A known observation is immutable.
// The decision ID is idempotent; conflicting replay refuses. Every adjustment
// remains in order, and the original unknown cost and spent action are retained.
func (s Store) ReconcileUnknown(ctx context.Context, id, decisionID string, ceiling int64, reason string) (Snapshot, error) {
	if id == "" || strings.TrimSpace(decisionID) == "" || len(decisionID) > 128 || ceiling < 0 || strings.TrimSpace(reason) == "" || len(reason) > 1024 {
		return Snapshot{}, errors.New("invalid operator cost reconciliation")
	}
	return s.update(ctx, func(state *Snapshot, now time.Time) error {
		k := key(id)
		entry, exists := state.Entries[k]
		if !exists {
			return errors.New("cannot reconcile an unreserved action")
		}
		for _, old := range entry.Reconciliations {
			if old.ID == decisionID {
				if old.CeilingMicros != ceiling || old.Reason != reason {
					return errors.New("conflicting operator decision")
				}
				return nil
			}
		}
		observed := entry.ReserveMicros
		if entry.Settled {
			observed = entry.ActualMicros
		}
		if observed != nil {
			return errors.New("operator reconciliation cannot replace a known monetary observation")
		}
		if entry.ReserveMicros != nil && ceiling < *entry.ReserveMicros {
			return errors.New("operator reconciliation cannot lower the retained conservative reservation while actual cost is unknown")
		}
		entry.Reconciliations = append(entry.Reconciliations, Reconciliation{ID: decisionID, At: now, CeilingMicros: ceiling, Reason: reason})
		state.Entries[k] = entry
		return nil
	})
}

// Settle retains the spent action and reconciles its monetary exposure only.
// Unknown terminal cost remains unknown. Conflicting terminal prices fail.
func (s Store) Settle(ctx context.Context, id string, actualMicros *int64) (Snapshot, error) {
	if id == "" || (actualMicros != nil && *actualMicros < 0) {
		return Snapshot{}, errors.New("invalid settlement")
	}
	return s.update(ctx, func(state *Snapshot, _ time.Time) error {
		k := key(id)
		entry, exists := state.Entries[k]
		if !exists {
			return errors.New("cannot settle an unreserved action")
		}
		if entry.Settled {
			if (entry.ActualMicros == nil) != (actualMicros == nil) || (actualMicros != nil && *entry.ActualMicros != *actualMicros) {
				return errors.New("conflicting terminal cost")
			}
			return nil
		}
		entry.Settled, entry.ActualMicros = true, copyInt(actualMicros)
		state.Entries[k] = entry
		return nil
	})
}

func copyInt(n *int64) *int64 {
	if n == nil {
		return nil
	}
	v := *n
	return &v
}

// ExposureError sums observed cost and outstanding reservations, retaining
// distinct unknown-cost and overflow errors. Explicit operator ceilings are
// the fallback only while the original observation is unknown.
func (s Snapshot) ExposureError() (int64, error) {
	var total int64
	for _, entry := range s.Entries {
		n := entry.ReserveMicros
		if entry.Settled && entry.ActualMicros != nil {
			n = entry.ActualMicros
		} else if len(entry.Reconciliations) > 0 {
			ceiling := &entry.Reconciliations[len(entry.Reconciliations)-1].CeilingMicros
			if n == nil || *ceiling > *n {
				n = ceiling
			}
		}
		if n == nil || *n < 0 {
			return 0, ErrUnknownCost
		}
		if *n > math.MaxInt64-total {
			return 0, ErrCostOverflow
		}
		total += *n
	}
	return total, nil
}

func (s Store) update(ctx context.Context, change func(*Snapshot, time.Time) error) (Snapshot, error) {
	if s.Dir == "" || s.Scope == "" {
		return Snapshot{}, errors.New("budget directory and scope are required")
	}
	// Bound lock contention even when a caller supplies Background. A shorter
	// caller deadline still wins. This does not bound unrelated filesystem I/O.
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := fsutil.MkdirAllResilient(s.Dir, 0o700); err != nil {
		return Snapshot{}, err
	}
	release, err := lock(ctx, filepath.Join(s.Dir, "ledger.lock"))
	if err != nil {
		return Snapshot{}, err
	}
	defer release()
	now := time.Now().UTC()
	if s.now != nil {
		now = s.now().UTC()
	}
	path := filepath.Join(s.Dir, "ledger.json")
	state, err := read(path)
	if witnessErr := s.checkContinuity(os.IsNotExist(err)); witnessErr != nil {
		return Snapshot{}, witnessErr
	}
	if os.IsNotExist(err) {
		state = Snapshot{Schema: 1, Scope: s.Scope, StartedAt: now, Entries: map[string]Reservation{}}
	} else if err != nil {
		return Snapshot{}, err
	}
	if state.Scope != s.Scope {
		return state, errors.New("budget scope mismatch")
	}
	if now.Before(state.StartedAt) {
		return state, ErrClockSkew
	}
	for _, entry := range state.Entries {
		if now.Before(entry.ReservedAt) {
			return state, ErrClockSkew
		}
		for _, adjustment := range entry.Reconciliations {
			if now.Before(adjustment.At) {
				return state, ErrClockSkew
			}
		}
	}
	if err := change(&state, now); err != nil {
		return state, err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return state, err
	}
	if len(data) > 16<<20 {
		return state, errors.New("budget ledger exceeds limit")
	}
	write := writeSynced
	if s.persist != nil {
		write = s.persist
	}
	if err := write(path, append(data, '\n')); err != nil {
		return state, fmt.Errorf("persist budget before work: %w", err)
	}
	if err := s.establishContinuity(); err != nil {
		return state, fmt.Errorf("persist budget continuity before work: %w", err)
	}
	return state, nil
}

func read(path string) (Snapshot, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return Snapshot{}, err
	}
	if !info.Mode().IsRegular() || info.Size() > 16<<20 {
		return Snapshot{}, errors.New("budget ledger must be a bounded regular file")
	}
	f, err := os.Open(path)
	if err != nil {
		return Snapshot{}, err
	}
	defer f.Close()
	opened, err := f.Stat()
	if err != nil {
		return Snapshot{}, err
	}
	if !os.SameFile(info, opened) {
		return Snapshot{}, ErrSnapshotChanged
	}
	data, err := io.ReadAll(io.LimitReader(f, (16<<20)+1))
	if err != nil || len(data) > 16<<20 {
		return Snapshot{}, errors.New("cannot read budget ledger")
	}
	if err := checkJSON(json.NewDecoder(bytes.NewReader(data))); err != nil {
		return Snapshot{}, err
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var state Snapshot
	if err := dec.Decode(&state); err != nil {
		return Snapshot{}, err
	}
	if dec.Decode(new(any)) != io.EOF {
		return Snapshot{}, errors.New("trailing budget data")
	}
	if state.Schema != 1 || state.Scope == "" || state.StartedAt.IsZero() || state.Entries == nil {
		return Snapshot{}, errors.New("incomplete budget state")
	}
	for id, entry := range state.Entries {
		digest, err := hex.DecodeString(id)
		if err != nil || len(digest) != 32 || !validKind(entry.Kind) || entry.ReservedAt.Before(state.StartedAt) || (entry.ReserveMicros != nil && *entry.ReserveMicros < 0) || (entry.ActualMicros != nil && (*entry.ActualMicros < 0 || !entry.Settled)) {
			return Snapshot{}, errors.New("invalid budget entry")
		}
		seen := map[string]bool{}
		previous := entry.ReservedAt
		for _, adjustment := range entry.Reconciliations {
			if strings.TrimSpace(adjustment.ID) == "" || len(adjustment.ID) > 128 || seen[adjustment.ID] || adjustment.At.Before(previous) || adjustment.CeilingMicros < 0 || strings.TrimSpace(adjustment.Reason) == "" || len(adjustment.Reason) > 1024 {
				return Snapshot{}, errors.New("invalid operator reconciliation")
			}
			seen[adjustment.ID] = true
			previous = adjustment.At
		}
	}
	return state, nil
}

// Reject duplicate keys and case aliases before encoding/json can replace a
// charged field. Null is allowed only for explicitly unknown monetary values.
func checkJSON(dec *json.Decoder) error {
	return checkJSONValue(dec, 0, "")
}

func checkJSONValue(dec *json.Decoder, depth int, field string) error {
	if depth > 8 {
		return errors.New("budget JSON nesting exceeds schema bound")
	}
	tok, err := dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := tok.(json.Delim); ok {
		if delim == '[' && field == "reconciliations" {
			for dec.More() {
				if err := checkJSONValue(dec, depth+1, ""); err != nil {
					return err
				}
			}
			_, err := dec.Token()
			return err
		}
		if delim != '{' {
			return errors.New("unexpected budget JSON container")
		}
		seen := map[string]bool{}
		for dec.More() {
			t, err := dec.Token()
			if err != nil {
				return err
			}
			name, ok := t.(string)
			if !ok || name != strings.ToLower(name) || seen[name] {
				return errors.New("duplicate or aliased budget field")
			}
			seen[name] = true
			if err := checkJSONValue(dec, depth+1, name); err != nil {
				return err
			}
		}
		_, err = dec.Token()
		var required []string
		if seen["schema"] {
			required = []string{"schema", "scope", "started_at", "entries"}
		}
		if seen["kind"] {
			required = []string{"kind", "reserved_at", "reserve_micros", "settled", "actual_micros"}
		}
		if seen["id"] {
			required = []string{"id", "at", "ceiling_micros", "reason"}
		}
		for _, name := range required {
			if !seen[name] {
				return errors.New("incomplete budget object")
			}
		}
		return err
	}
	if depth == 0 {
		return errors.New("budget state must be an object")
	}
	if tok == nil && field != "reserve_micros" && field != "actual_micros" {
		return errors.New("null budget field")
	}
	return nil
}

func writeSynced(path string, data []byte) error {
	f, err := os.CreateTemp(filepath.Dir(path), ".budget-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	defer f.Close()
	if _, err = f.Write(data); err != nil {
		return err
	}
	if err = fsutil.SyncFile(f); err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return fsutil.ReplaceSyncedFile(f.Name(), path)
}
