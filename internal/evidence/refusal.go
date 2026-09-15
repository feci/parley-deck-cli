package evidence

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"parley-deck-cli/internal/fsutil"
	"parley-deck-cli/internal/telemetry"
)

const VerificationRefusalDirectory = "verification-refusals"

var refusalPendingName = regexp.MustCompile(`^observation-[0-9]+\.json$`)
var refusalHash = regexp.MustCompile(`^[0-9a-f]{64}$`)

// VerificationRefusal records an observation, not a model invocation or proof
// that an unobserved process terminated. Raw diagnostics and output are excluded.
// Independent observers share VerificationID once a frozen request exists.
type VerificationRefusal struct {
	Version        int                `json:"version"`
	ObservationID  string             `json:"observation_id"`
	VerificationID string             `json:"verification_id"`
	ObservedAt     time.Time          `json:"observed_at"`
	Observer       string             `json:"observer"`
	Stage          string             `json:"stage"`
	Idea           *string            `json:"idea"`
	RunID          *string            `json:"run_id"`
	Verifier       *string            `json:"verifier"`
	InvocationID   *string            `json:"invocation_id"`
	RequestSHA256  *string            `json:"request_sha256"`
	OriginalSHA256 *string            `json:"original_report_sha256"`
	ReceiptSHA256  *string            `json:"receipt_sha256"`
	Executions     []RefusalExecution `json:"executions"`
}

type RefusalExecution struct {
	Name          *string `json:"name"`
	Status        Status  `json:"status"`
	CommandSHA256 *string `json:"command_sha256"`
}

type RefusalEntry struct {
	SHA256       string               `json:"sha256"`
	Record       *VerificationRefusal `json:"record"`
	Canonical    bool                 `json:"canonical"`
	PendingFiles []string             `json:"pending_files"`
	Problem      string               `json:"problem"`
}

func RefusalHash(value string) *string {
	if !refusalHash.MatchString(value) {
		return nil
	}
	return &value
}

func refusalDigest(data []byte) string { sum := sha256.Sum256(data); return hex.EncodeToString(sum[:]) }

func encodeRefusal(r VerificationRefusal) ([]byte, error) {
	if r.Version != 1 || !refusalHash.MatchString(r.ObservationID) || !refusalHash.MatchString(r.VerificationID) || r.ObservedAt.IsZero() {
		return nil, errors.New("invalid refusal identity or version")
	}
	if r.Observer != "driver" && r.Observer != "helper" {
		return nil, errors.New("invalid refusal observer")
	}
	switch r.Stage {
	case "preflight", "identity", "original-report", "bindings", "runtime", "request", "launch", "acceptance", "execution", "publication":
	default:
		return nil, errors.New("invalid refusal stage")
	}
	for _, v := range []*string{r.Idea, r.RunID, r.Verifier, r.InvocationID} {
		if v != nil && telemetry.SafeLabel(*v) == nil {
			return nil, errors.New("unsafe refusal label")
		}
	}
	for _, v := range []*string{r.RequestSHA256, r.OriginalSHA256, r.ReceiptSHA256} {
		if v != nil && RefusalHash(*v) == nil {
			return nil, errors.New("invalid refusal digest")
		}
	}
	if r.RequestSHA256 != nil && *r.RequestSHA256 != r.VerificationID {
		return nil, errors.New("refusal request identity differs")
	}
	if len(r.Executions) > 4096 {
		return nil, errors.New("too many refusal executions")
	}
	for _, e := range r.Executions {
		if e.Name != nil && telemetry.SafeLabel(*e.Name) == nil {
			return nil, errors.New("unsafe refusal criterion")
		}
		if e.CommandSHA256 != nil && RefusalHash(*e.CommandSHA256) == nil {
			return nil, errors.New("invalid refusal command digest")
		}
		switch e.Status {
		case StatusPass, StatusFail, StatusSkipped, StatusNotRun:
		default:
			return nil, errors.New("invalid refusal execution status")
		}
	}
	data, err := json.MarshalIndent(r, "", "  ")
	if len(data) > 1<<20 {
		return nil, errors.New("refusal exceeds size limit")
	}
	return data, err
}

func decodeRefusal(data []byte) (*VerificationRefusal, error) {
	var r VerificationRefusal
	if json.Unmarshal(data, &r) != nil {
		return nil, errors.New("invalid refusal JSON")
	}
	canonical, err := encodeRefusal(r)
	// Exact canonical bytes enforce mandatory nullable fields, duplicate/alias/
	// unknown-field refusal and a stable content-addressed recovery target.
	if err != nil || !bytes.Equal(data, canonical) {
		return nil, errors.New("invalid or noncanonical refusal record")
	}
	return &r, nil
}

func refusalDirs(ideaDir string) (string, string, error) {
	dir, err := filepath.Abs(ideaDir)
	if err != nil {
		return "", "", err
	}
	dir, err = filepath.EvalSymlinks(dir)
	if err != nil {
		return "", "", err
	}
	guard, err := reportGuardDirectory(dir)
	if err != nil {
		return "", "", err
	}
	return filepath.Join(guard, "refusals"), filepath.Join(dir, VerificationRefusalDirectory), nil
}

// Refuse existing aliases, including a symlink in any administration component.
// Already-existing system aliases are canonicalized before this helper is used.
func realRefusalDir(path string, create bool) error {
	info, err := os.Lstat(path)
	if err == nil {
		if !info.IsDir() {
			return errors.New("refusal storage must use real directories")
		}
		if filepath.Dir(path) != path {
			return realRefusalDir(filepath.Dir(path), false)
		}
		return nil
	}
	if !os.IsNotExist(err) || !create {
		return err
	}
	if err := realRefusalDir(filepath.Dir(path), true); err != nil {
		return err
	}
	if err := os.Mkdir(path, 0700); err != nil && !os.IsExist(err) {
		return err
	}
	info, err = os.Lstat(path)
	if err != nil || !info.IsDir() {
		return errors.New("refusal storage directory changed")
	}
	return syncRefusalDir(filepath.Dir(path))
}

func syncRefusalDir(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return fsutil.SyncFile(f)
}

// RetainVerificationRefusal does not acquire the report guard: a failed helper
// must retain its observation even if the parent stops or the lock origin is
// unavailable. Each observer gets its own exclusive file. Interrupted writes
// remain visible as invalid pending entries; they never become acceptance.
func RetainVerificationRefusal(ideaDir string, r VerificationRefusal) (string, error) {
	data, err := encodeRefusal(r)
	if err != nil {
		return "", err
	}
	sum := refusalDigest(data)
	pending, _, err := refusalDirs(ideaDir)
	if err != nil {
		return sum, err
	}
	if err := realRefusalDir(pending, true); err != nil {
		return sum, err
	}
	f, err := os.CreateTemp(pending, "observation-*.json")
	if err != nil {
		return sum, err
	}
	_, writeErr := f.Write(data)
	syncErr := fsutil.SyncFile(f)
	closeErr := f.Close()
	return sum, errors.Join(writeErr, syncErr, closeErr, syncRefusalDir(pending))
}

// InspectVerificationRefusals is read-only, including for a pristine idea. Both
// canonical history and the local pending spool are inventoried. Duplicate
// copies of exact bytes coalesce, while interrupted/malformed entries stay visible.
func InspectVerificationRefusals(ideaDir string) ([]RefusalEntry, error) {
	pending, canonical, err := refusalDirs(ideaDir)
	if err != nil {
		return nil, err
	}
	entries := map[string]*RefusalEntry{}
	var invalid []RefusalEntry
	for _, dir := range []string{pending, canonical} {
		if err := realRefusalDir(dir, false); os.IsNotExist(err) {
			continue
		} else if err != nil {
			return nil, err
		}
		files, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			data, readErr := readBoundedEvidenceFile(filepath.Join(dir, file.Name()), 1<<20)
			r, decodeErr := decodeRefusal(data)
			sum := refusalDigest(data)
			if readErr != nil || decodeErr != nil || (dir == pending && !refusalPendingName.MatchString(file.Name())) || (dir == canonical && file.Name() != sum+".json") {
				// Do not echo untrusted filenames or parser excerpts into portable output.
				invalid = append(invalid, RefusalEntry{Problem: "unreadable, incomplete or conflicting refusal storage entry"})
				continue
			}
			entry := entries[sum]
			if entry == nil {
				entry = &RefusalEntry{SHA256: sum, Record: r}
				entries[sum] = entry
			}
			if dir == canonical {
				entry.Canonical = true
			} else {
				entry.PendingFiles = append(entry.PendingFiles, file.Name())
			}
		}
	}
	identities := map[string]string{}
	for sum, entry := range entries {
		id := entry.Record.ObservationID
		if previous, found := identities[id]; found && previous != sum {
			entry.Problem = "conflicting observation identity"
			entries[previous].Problem = entry.Problem
		} else {
			identities[id] = sum
		}
	}
	result := make([]RefusalEntry, 0, len(entries)+len(invalid))
	for _, entry := range entries {
		result = append(result, *entry)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].SHA256 < result[j].SHA256 })
	return append(result, invalid...), nil
}

// PublishVerificationRefusal recovers exactly one observed record under the
// existing report guard. It never recreates a missing lock origin, overwrites
// conflicting bytes, removes pending history, or grants completion. commit must
// verify that the exact canonical bytes reached Git history before returning.
func PublishVerificationRefusal(ctx context.Context, ideaDir, sum string, commit func(string, []byte) error) error {
	if !refusalHash.MatchString(sum) || commit == nil {
		return errors.New("exact refusal digest and commit callback required")
	}
	return WithReportWriter(ctx, ideaDir, func(_ *ReportWriter) error {
		entries, err := InspectVerificationRefusals(ideaDir)
		if err != nil {
			return err
		}
		var data []byte
		for _, entry := range entries {
			if entry.SHA256 == sum && entry.Record != nil && entry.Problem == "" {
				data, err = encodeRefusal(*entry.Record)
				break
			}
		}
		if err != nil {
			return err
		}
		if data == nil {
			return errors.New("requested refusal record is unavailable")
		}
		_, dir, err := refusalDirs(ideaDir)
		if err != nil {
			return err
		}
		if err := realRefusalDir(dir, true); err != nil {
			return err
		}
		path := filepath.Join(dir, sum+".json")
		current, err := readBoundedEvidenceFile(path, 1<<20)
		if os.IsNotExist(err) {
			f, err := os.CreateTemp(dir, ".refusal-"+sum+"-*.tmp")
			if err != nil {
				return err
			}
			defer os.Remove(f.Name())
			if _, err := f.Write(data); err != nil {
				f.Close()
				return err
			}
			if err := fsutil.SyncFile(f); err != nil {
				f.Close()
				return err
			}
			if err := f.Close(); err != nil {
				return err
			}
			if err := fsutil.ReplaceSyncedFile(f.Name(), path); err != nil {
				return err
			}
		} else if err != nil || !bytes.Equal(current, data) {
			return errors.New("canonical refusal conflicts with the retained observation")
		}
		// A crash before rename can leave a partial sibling staging file. Its
		// name binds the requested digest; only a prefix of these exact retained
		// bytes is attributable to this publication. Under the same guard there
		// is no live cooperating publisher of such a file. Preserve divergent or
		// unrelated entries as unresolved evidence rather than deleting them.
		stagedFiles, err := os.ReadDir(dir)
		if err != nil {
			return err
		}
		for _, file := range stagedFiles {
			if !strings.HasPrefix(file.Name(), ".refusal-"+sum+"-") || !strings.HasSuffix(file.Name(), ".tmp") {
				continue
			}
			stagedPath := filepath.Join(dir, file.Name())
			staged, err := readBoundedEvidenceFile(stagedPath, 1<<20)
			if err != nil || !bytes.HasPrefix(data, staged) {
				return errors.New("interrupted refusal staging conflicts with the exact record")
			}
			if err := os.Remove(stagedPath); err != nil {
				return err
			}
		}
		if err := syncRefusalDir(dir); err != nil {
			return err
		}
		if err := commit(path, data); err != nil {
			return fmt.Errorf("canonical refusal is not verified committed: %w", err)
		}
		current, err = readBoundedEvidenceFile(path, 1<<20)
		if err != nil || !bytes.Equal(current, data) {
			return errors.New("canonical refusal changed during commit")
		}
		return nil
	})
}
