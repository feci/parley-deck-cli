package budget

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"parley-deck-cli/internal/fsutil"
)

// A durable witness is published before the first grant can return. It lets a
// later caller distinguish a fresh ledger from a lost previously charged file.
// The witness contains no charges and never substitutes for the missing ledger.
// Deleting all history, including this witness, is not same-UID tamper protection.
func (s Store) continuityPath() string { return filepath.Join(s.Dir, "ledger-established") }
func (s Store) continuityBytes() []byte {
	return []byte("parley-budget-established/v1\n" + key(s.Scope) + "\n")
}

func (s Store) checkContinuity(ledgerMissing bool) error {
	data, err := readLockOrigin(s.continuityPath()) // same bounded regular-file reader
	if os.IsNotExist(err) {
		return nil
	} // fresh, or existing pre-witness ledger
	if err != nil {
		return fmt.Errorf("budget continuity witness unreadable: %w", err)
	}
	if !bytes.Equal(data, s.continuityBytes()) {
		return errors.New("budget continuity witness is malformed or belongs to another scope")
	}
	if ledgerMissing {
		return errors.New("previously charged budget ledger is missing; restore its preserved history before continuing; refusing a zero-count reset")
	}
	return nil
}

// Called under the same store lock after the ledger persistence barrier and
// before returning permission. A crash before this publication cannot have
// granted work; a crash after it leaves the witness guarding the charged ledger.
func (s Store) establishContinuity() error {
	if _, err := os.Lstat(s.continuityPath()); err == nil {
		return s.checkContinuity(false)
	} else if !os.IsNotExist(err) {
		return err
	}
	f, err := os.CreateTemp(s.Dir, ".ledger-established-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	defer f.Close()
	if _, err := f.Write(s.continuityBytes()); err != nil {
		return err
	}
	if err := fsutil.SyncFile(f); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	if err := publishExclusive(f.Name(), s.continuityPath()); err != nil && !os.IsExist(err) {
		return err
	}
	return s.checkContinuity(false)
}
