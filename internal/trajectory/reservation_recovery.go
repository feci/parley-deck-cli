package trajectory

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/fsutil"
)

// This is the original precharge state, not a reconstructed source observation.
// It is published exclusively before the ledger transaction and never rewritten.
type reservationIntent struct {
	Version      int                           `json:"version"`
	Root         string                        `json:"root"`
	PreparedAt   time.Time                     `json:"prepared_at"`
	Accounting   budget.CycleReservationIntent `json:"accounting"`
	Before       State                         `json:"before"`
	BeforeSHA256 string                        `json:"before_sha256"`
}

func intentName(entry string) string {
	return filepath.Join("reservation-intents", entry+".json")
}

func openIntentRoot(b budget.CycleBinding, create bool) (*os.Root, error) {
	dir, err := os.OpenRoot(filepath.Dir(b.Store.Dir))
	if err != nil {
		return nil, err
	}
	if create {
		if err = dir.Mkdir("reservation-intents", 0700); err != nil && !os.IsExist(err) {
			dir.Close()
			return nil, err
		}
	}
	info, err := dir.Lstat("reservation-intents")
	if err != nil || !info.IsDir() {
		dir.Close()
		return nil, errors.New("original reservation intent directory is missing or not a real directory")
	}
	if create {
		if err = syncVerificationDirectory(dir); err != nil {
			dir.Close()
			return nil, err
		}
	}
	return dir, nil
}

func syncIntent(dir *os.Root, entry string) error {
	f, err := dir.Open(intentName(entry))
	if err != nil {
		return err
	}
	if err = errors.Join(fsutil.SyncFile(f), f.Close()); err != nil {
		return err
	}
	f, err = dir.Open("reservation-intents")
	if err != nil {
		return err
	}
	return errors.Join(fsutil.SyncFile(f), f.Close(), syncVerificationDirectory(dir))
}

func publishReservationIntent(dir *os.Root, i reservationIntent) (string, error) {
	raw, err := canonical(i)
	if err != nil || len(raw) > 16<<20 {
		return "", errors.New("precharge intent exceeds its publication bound")
	}
	f, err := dir.OpenFile(intentName(i.Accounting.EntryKey), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return "", err
	}
	// A partial publication stays visible and cannot be overwritten by retry.
	_, writeErr := f.Write(raw)
	if err = errors.Join(writeErr, fsutil.SyncFile(f), f.Close()); err != nil {
		return "", err
	}
	if err = syncIntent(dir, i.Accounting.EntryKey); err != nil {
		return "", err
	}
	return digest(raw), nil
}

func readReservationIntent(dir *os.Root, entry string) (reservationIntent, string, error) {
	var i reservationIntent
	if !validHash(entry) {
		return i, "", errors.New("reservation recovery requires the original hashed entry key")
	}
	sha, err := readReconciliationJSON(dir, intentName(entry), &i)
	if err != nil {
		return i, "", err
	}
	raw, err := canonical(i.Before)
	if err != nil || i.Version != 1 || i.PreparedAt.IsZero() || i.BeforeSHA256 != digest(raw) || i.Accounting.EntryKey != entry {
		return i, "", errors.New("original precharge intent is incomplete or changed")
	}
	return i, sha, nil
}

func (o *Observer) PrepareCycleReservation(ctx context.Context, b budget.CycleBinding, ledger budget.Snapshot, accounting budget.CycleReservationIntent) error {
	o.intentSHA256 = ""
	s, raw, err := readState(statePath(b))
	if err != nil {
		return err
	}
	if o.prepared == nil || !bytes.Equal(raw, o.prepared) {
		return errors.New("trajectory changed before durable reservation preparation")
	}
	if err = accounting.CheckBefore(b.Policy, ledger); err != nil {
		return err
	}
	if err = validateState(s, b, ledger); err != nil {
		return err
	}
	root, err := canonicalRoot(o.Root)
	if err != nil {
		return err
	}
	i := reservationIntent{1, root, time.Now().UTC(), accounting, s, digest(raw)}
	dir, err := openIntentRoot(b, true)
	if err != nil {
		return err
	}
	defer dir.Close()
	o.intentSHA256, err = publishReservationIntent(dir, i)
	return err
}

// Prefix accounting uses actual retained entries. It never creates ledger facts
// or discards an unexpected entry from the real ledger.
func intentPrefixLedger(s budget.Snapshot, before State) (budget.Snapshot, error) {
	prefix := s
	prefix.Entries = make(map[string]budget.Reservation, len(before.Attempts))
	for _, a := range before.Attempts {
		e, ok := s.Entries[a.Charge.EntryKey]
		if !ok {
			return prefix, errors.New("an earlier original charge is missing")
		}
		prefix.Entries[a.Charge.EntryKey] = e
	}
	return prefix, nil
}

func validateIntentBefore(b budget.CycleBinding, ledger budget.Snapshot, i reservationIntent) error {
	if err := i.Accounting.CheckPolicy(b.Policy); err != nil {
		return err
	}
	prefix, err := intentPrefixLedger(ledger, i.Before)
	if err != nil {
		return err
	}
	if i.Accounting.Policy.TrajectorySHA256 != b.Policy.TrajectorySHA256 ||
		i.Accounting.Policy.Idea != b.Policy.Idea || i.Accounting.Policy.IdeaPath != b.Policy.IdeaPath ||
		i.PreparedAt.Before(ledger.StartedAt) {
		return errors.New("original intent differs from the frozen trajectory authority")
	}
	if err = i.Accounting.CheckBefore(i.Accounting.Policy, prefix); err != nil {
		return err
	}
	if err = validateState(i.Before, b, prefix); err != nil {
		return err
	}
	h := trajectoryHistory(i.Before)
	if h.Unreconciled > 0 || h.ReviewPending || len(h.InconclusivePending) > 0 {
		return errors.New("precharge intent bypassed a pending trajectory gate")
	}
	return nil
}

func intentAttempt(b budget.CycleBinding, ledger budget.Snapshot, i reservationIntent, sha string) (Attempt, error) {
	var a Attempt
	if err := validateIntentBefore(b, ledger, i); err != nil {
		return a, err
	}
	charge, err := i.Accounting.CheckCharge(ledger)
	if err != nil {
		return a, err
	}
	if charge.ReservedAt.Before(i.PreparedAt) {
		return a, errors.New("charge predates its retained preparation")
	}
	before, archive, err := currentTrajectorySource(i.Before)
	if err != nil || !before.Clean {
		return a, errors.New("precharge intent lacks its original clean source")
	}
	return Attempt{Sequence: len(i.Before.Attempts) + 1, Charge: charge, Before: before, BeforeArchive: archive, ReservationIntentSHA256: sha}, nil
}

func compareIntentPrefix(s State, i reservationIntent) error {
	n := len(i.Before.Attempts)
	if len(s.Attempts) < n || len(s.Resolutions) < len(i.Before.Resolutions) || len(s.Continuations) < len(i.Before.Continuations) {
		return errors.New("trajectory lost the original precharge history")
	}
	prefix := s
	prefix.Version = i.Before.Version // Version 3 may have been added by a later resolution.
	prefix.Attempts = prefix.Attempts[:n]
	prefix.Resolutions = prefix.Resolutions[:len(i.Before.Resolutions)]
	prefix.Continuations = prefix.Continuations[:len(i.Before.Continuations)]
	if !sameJSON(prefix, i.Before) {
		return errors.New("trajectory differs from its retained precharge history")
	}
	return nil
}

func checkReservationIntents(ctx context.Context, b budget.CycleBinding, s State) error {
	var dir *os.Root
	defer func() {
		if dir != nil {
			dir.Close()
		}
	}()
	var ledger budget.Snapshot
	for _, a := range s.Attempts {
		if a.ReservationIntentSHA256 == "" {
			if _, err := os.Lstat(filepath.Join(filepath.Dir(b.Store.Dir), intentName(a.Charge.EntryKey))); !os.IsNotExist(err) {
				return errors.New("charged history lost its existing precharge intent reference")
			}
			continue // Historical rows remain explicitly without precharge evidence.
		}
		var err error
		if dir == nil {
			dir, err = openIntentRoot(b, false)
			if err != nil {
				return err
			}
			ledger, err = b.Store.Inspect(ctx)
			if err != nil {
				return err
			}
		}
		i, sha, err := readReservationIntent(dir, a.Charge.EntryKey)
		if err != nil {
			return err
		}
		want, err := intentAttempt(b, ledger, i, sha)
		if err != nil {
			return err
		}
		if sha != a.ReservationIntentSHA256 || a.Sequence != want.Sequence || !sameJSON(a.Charge, want.Charge) || a.Before != want.Before || a.BeforeArchive != want.BeforeArchive {
			return errors.New("charged trajectory changed its original reservation intent")
		}
		if err = compareIntentPrefix(s, i); err != nil {
			return err
		}
	}
	return nil
}

type ReservationRecoveryPreview struct {
	Version         int      `json:"version"`
	Root            string   `json:"root"`
	Idea            string   `json:"idea"`
	EntryKey        string   `json:"entry_key"`
	IntentSHA256    string   `json:"intent_sha256"`
	BeforeSHA256    string   `json:"before_sha256"`
	Status          string   `json:"status"`
	Attempt         *Attempt `json:"attempt"`
	Permission      string   `json:"permission"`
	ExecutionStatus string   `json:"execution_status"`
}

func (p ReservationRecoveryPreview) SHA256() string { raw, _ := canonical(p); return digest(raw) }

func PreviewReservationRecovery(ctx context.Context, root, idea, entry string) (ReservationRecoveryPreview, error) {
	return recoverReservation(ctx, root, idea, entry, "", false, writeState)
}

func RecoverReservation(ctx context.Context, root, idea, entry, expected string) (ReservationRecoveryPreview, error) {
	return recoverReservation(ctx, root, idea, entry, expected, true, writeState)
}

func recoverReservation(ctx context.Context, root, idea, entry, expected string, apply bool, persist func(string, State) error) (ReservationRecoveryPreview, error) {
	var p ReservationRecoveryPreview
	root, err := canonicalRoot(root)
	if err != nil {
		return p, err
	}
	if !runtimeID(idea) || !validHash(entry) || apply && !validHash(expected) {
		return p, errors.New("reservation recovery requires exact original identity and preview")
	}
	b, err := budget.LoadCycleBinding(ctx, root, idea, budget.Fixup)
	if err != nil {
		return p, err
	}
	if b == nil {
		return p, errors.New("original cycle policy is missing")
	}
	wait, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	release, err := budget.AcquireResourceGuard(wait, filepath.Dir(b.Store.Dir))
	if err != nil {
		return p, err
	}
	defer release()
	b, err = budget.LoadCycleBinding(ctx, root, idea, budget.Fixup)
	if err != nil {
		return p, err
	}
	if b == nil {
		return p, errors.New("original cycle policy disappeared")
	}
	ledger, err := b.Store.Inspect(ctx)
	if err != nil {
		return p, err
	}
	s, raw, err := readState(statePath(*b))
	if err != nil {
		return p, err
	}
	dir, err := openIntentRoot(*b, false)
	if err != nil {
		return p, err
	}
	defer dir.Close()
	i, sha, err := readReservationIntent(dir, entry)
	if err != nil {
		return p, err
	}
	if i.Root != root {
		return p, errors.New("recovery worktree differs from the original intent")
	}
	if err = validateIntentBefore(*b, ledger, i); err != nil {
		return p, err
	}
	if err = compareIntentPrefix(s, i); err != nil {
		return p, err
	}
	p = ReservationRecoveryPreview{Version: 1, Root: root, Idea: idea, EntryKey: entry, IntentSHA256: sha, BeforeSHA256: i.BeforeSHA256, Permission: "none", ExecutionStatus: "not-established"}
	if _, charged := ledger.Entries[entry]; !charged {
		if err = validateState(s, *b, ledger); err != nil {
			return p, err
		}
		if err = checkStateSnapshots(ctx, *b, s); err != nil {
			return p, err
		}
		p.Status = "intent-without-published-charge"
		if apply {
			return p, errors.New("no matching spent charge exists; recovery grants no reservation or retry")
		}
		return p, nil
	}
	a, err := intentAttempt(*b, ledger, i, sha)
	if err != nil {
		return p, err
	}
	p.Status, p.Attempt = "charged-observation", &a
	missing := len(s.Attempts) == len(i.Before.Attempts)
	if missing {
		if digest(raw) != i.BeforeSHA256 || len(ledger.Entries) != len(s.Attempts)+1 {
			return p, errors.New("recovery may repair only this exact missing charged row")
		}
		s.Attempts = append(s.Attempts, a)
	} else {
		old := s.Attempts[a.Sequence-1]
		if old.ReservationIntentSHA256 != sha || !sameJSON(old.Charge, a.Charge) || old.Before != a.Before || old.BeforeArchive != a.BeforeArchive {
			return p, errors.New("existing charged row differs from the original recovery")
		}
	}
	if err = validateState(s, *b, ledger); err != nil {
		return p, err
	}
	if err = checkStateSnapshots(ctx, *b, s); err != nil {
		return p, err
	}
	if !apply {
		return p, nil
	}
	if p.SHA256() != expected {
		return p, errors.New("reservation recovery changed since its exact preview")
	}
	if err = syncIntent(dir, entry); err != nil {
		return p, err
	}
	if missing {
		err = persist(statePath(*b), s)
	} else {
		// Complete a failed durability barrier without rewriting the original bytes.
		f, openErr := dir.Open("trajectory.json")
		if openErr != nil {
			return p, openErr
		}
		err = errors.Join(fsutil.SyncFile(f), f.Close(), syncVerificationDirectory(dir))
	}
	return p, err
}
