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
	ErrReserved    = errors.New("action already reserved: reconcile it before any new execution")
	ErrLimit       = errors.New("budget exhausted")
	ErrUnknownCost = errors.New("monetary exposure is unknown")
)

type Kind string

const (
	Launch      Kind = "launch"
	DriverStep  Kind = "driver-step"
	Fixup       Kind = "fixup"
	CrossReview Kind = "cross-review"
)

type Limits struct {
	Actions    map[Kind]int
	CostMicros int64
	WallClock  time.Duration
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
	return s.update(ctx, func(state *Snapshot, now time.Time) error {
		id := key(req.ID)
		if _, exists := state.Entries[id]; exists {
			return ErrReserved
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
		if cap := limits.Actions[req.Kind]; cap > 0 && count >= cap {
			return fmt.Errorf("%w: %s", ErrLimit, req.Kind)
		}
		if limits.CostMicros > 0 {
			if req.ReserveMicros == nil {
				return ErrUnknownCost
			}
			exposure, known := state.Exposure()
			if !known {
				return ErrUnknownCost
			}
			if exposure > limits.CostMicros || *req.ReserveMicros > limits.CostMicros-exposure {
				return fmt.Errorf("%w: monetary ceiling", ErrLimit)
			}
		}
		state.Entries[id] = Reservation{Kind: req.Kind, ReservedAt: now, ReserveMicros: copyInt(req.ReserveMicros)}
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

// Exposure is the sum of observed terminal cost and outstanding conservative
// reservations. It never represents unknown cost as zero or permits overflow.
func (s Snapshot) Exposure() (int64, bool) {
	var total int64
	for _, entry := range s.Entries {
		n := entry.ReserveMicros
		if entry.Settled {
			n = entry.ActualMicros
		}
		if n == nil || *n < 0 || *n > math.MaxInt64-total {
			return 0, false
		}
		total += *n
	}
	return total, true
}

func (s Store) update(ctx context.Context, change func(*Snapshot, time.Time) error) (Snapshot, error) {
	if s.Dir == "" || s.Scope == "" {
		return Snapshot{}, errors.New("budget directory and scope are required")
	}
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
	if os.IsNotExist(err) {
		state = Snapshot{Schema: 1, Scope: s.Scope, StartedAt: now, Entries: map[string]Reservation{}}
	} else if err != nil {
		return Snapshot{}, err
	}
	if state.Scope != s.Scope || now.Before(state.StartedAt) {
		return state, errors.New("budget scope or clock mismatch")
	}
	for _, entry := range state.Entries {
		if now.Before(entry.ReservedAt) {
			return state, errors.New("budget reservation is in the future")
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
	if err != nil || !os.SameFile(info, opened) {
		return Snapshot{}, errors.New("budget ledger changed during open")
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
	}
	return state, nil
}

// Reject duplicate keys and case aliases before encoding/json can replace a
// charged field. Null is allowed only for explicitly unknown monetary values.
func checkJSON(dec *json.Decoder) error {
	tok, err := dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := tok.(json.Delim); ok {
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
			var value json.RawMessage
			if err := dec.Decode(&value); err != nil {
				return err
			}
			if bytes.Equal(value, []byte("null")) && name != "reserve_micros" && name != "actual_micros" {
				return errors.New("null budget field")
			}
			if len(value) > 0 && value[0] == '{' {
				if err := checkJSON(json.NewDecoder(bytes.NewReader(value))); err != nil {
					return err
				}
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
		for _, name := range required {
			if !seen[name] {
				return errors.New("incomplete budget object")
			}
		}
		return err
	}
	return errors.New("budget state must be an object")
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
