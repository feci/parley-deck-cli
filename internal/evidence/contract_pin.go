package evidence

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"parley-deck-cli/internal/fsutil"
)

// The original contract is an obligation, separate from replaceable reports and
// per-run cursors. In Git it lives in administration, outside agent-authored idea
// files. Deleting all runtime/Git history remains outside this cooperative model.
func ReadContractPin(ideaDir string) (string, error) {
	dir, err := filepath.Abs(ideaDir)
	if err != nil {
		return "", err
	}
	dir, err = filepath.EvalSymlinks(dir)
	if err != nil {
		return "", err
	}
	guardDir, err := reportGuardDirectory(dir)
	if err != nil {
		return "", err
	}
	data, err := readBoundedEvidenceFile(filepath.Join(guardDir, "checks-contract"), 65)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	if len(data) != 65 || data[64] != '\n' {
		return "", errors.New("original checks witness is malformed")
	}
	if decoded, err := hex.DecodeString(string(data[:64])); err != nil || len(decoded) != 32 {
		return "", errors.New("original checks witness is malformed")
	}
	return string(data[:64]), nil
}

func PinChecksContract(ctx context.Context, ideaDir, digest string) error {
	return WithReportWriter(ctx, ideaDir, func(w *ReportWriter) error { return w.PinChecksContract(digest) })
}

// PinChecksContract is called before any named-contract agent/check work. It
// adds an obligation once; an existing different contract is never overwritten.
func (w *ReportWriter) PinChecksContract(digest string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.active {
		return errors.New("evidence writer transaction has ended")
	}
	if digest == "" {
		return nil
	}
	decoded, err := hex.DecodeString(digest)
	if err != nil || len(decoded) != 32 {
		return errors.New("invalid original checks digest")
	}
	prior, err := ReadContractPin(w.dir)
	if err != nil {
		return err
	}
	if prior != "" {
		if prior != digest {
			return errors.New("original named checks contract changed")
		}
		return nil
	}
	guardDir, err := reportGuardDirectory(w.dir)
	if err != nil {
		return err
	}
	f, err := os.CreateTemp(guardDir, ".checks-contract-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	defer f.Close()
	if _, err := fmt.Fprintln(f, digest); err != nil {
		return err
	}
	if err := fsutil.SyncFile(f); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return fsutil.ReplaceSyncedFile(f.Name(), filepath.Join(guardDir, "checks-contract"))
}
