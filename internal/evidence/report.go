package evidence

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"parley-deck-cli/internal/fsutil"
)

// report.go persists and reloads the typed evidence Report.
//
// The report filename is fixed (EVIDENCE.json, inside the idea directory) so
// the set of "defined evidence artifacts" excluded from the tested-tree digest
// is enumerable and small. Save is atomic (temp file + rename): a partial or
// failed write must never leave a half-written report that a later close could
// read as valid. Any write error is returned to the caller, and the caller
// MUST treat it as a failure of the whole completion attempt — never warn and
// continue as if evidence existed.

// ReportFileName is the canonical typed-evidence artifact name.
const ReportFileName = "EVIDENCE.json"

// ReportPath returns the canonical report path for an idea directory.
func ReportPath(ideaDir string) string {
	return filepath.Join(ideaDir, ReportFileName)
}

// Save atomically persists the report. A failure here invalidates the whole
// completion attempt (evidence-write failure is a failure, not a warning).
func Save(ideaDir string, r *Report) error {
	return WithReportWriter(context.Background(), ideaDir, func(w *ReportWriter) error {
		_, err := w.Save(r)
		return err
	})
}

func saveReport(ideaDir string, r *Report) ([]byte, error) {
	if r == nil {
		return nil, fmt.Errorf("evidence: cannot save a nil report")
	}
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("evidence: marshal report: %w", err)
	}
	final := ReportPath(ideaDir)
	tmp, err := os.CreateTemp(ideaDir, ".evidence-*.json.tmp")
	if err != nil {
		return nil, fmt.Errorf("evidence: create temp report: %w", err)
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return nil, fmt.Errorf("evidence: write report: %w", err)
	}
	// Durability barrier before the atomic rename. fsutil.SyncFile requests a
	// full sync and falls back to ordinary fsync ONLY when the volume rejects
	// the device-specific operation with ENOTTY (Darwin F_FULLFSYNC on shared
	// volumes); every other error — EINVAL, EOPNOTSUPP, real I/O failures —
	// aborts the save. A sync failure is a write failure, never a warning.
	if err := fsutil.SyncFile(tmp); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return nil, fmt.Errorf("evidence: sync report: %w", err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return nil, fmt.Errorf("evidence: close report: %w", err)
	}
	if err := fsutil.ReplaceSyncedFile(tmpName, final); err != nil {
		_ = os.Remove(tmpName)
		return nil, fmt.Errorf("evidence: commit report: %w", err)
	}
	return data, nil
}

// Load reads a previously saved report. A missing or corrupt report is an
// error: closure callers treat that as fail-closed (no evidence).
func Load(ideaDir string) (*Report, error) {
	data, err := os.ReadFile(ReportPath(ideaDir))
	if err != nil {
		return nil, fmt.Errorf("evidence: read report: %w", err)
	}
	var r Report
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("evidence: corrupt report: %w", err)
	}
	return &r, nil
}
