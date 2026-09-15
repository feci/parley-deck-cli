package trajectory

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/fsutil"
)

const parentRecoveryName = "parent-recovery.json"

type ParentObservation struct {
	Kind   string `json:"kind"`
	SHA256 string `json:"sha256,omitempty"`
}

// The original parent stays untouched. Derived records describe independently
// retained executions, not a successful parent write or new execution permission.
type ParentRecoveryPreview struct {
	Version      int               `json:"version"`
	Root         string            `json:"root"`
	Idea         string            `json:"idea"`
	RunID        string            `json:"run_id"`
	StateSHA256  string            `json:"state_sha256"`
	Original     ParentObservation `json:"original"`
	Derived      ParentDerivation  `json:"derived"`
	ParentSHA256 string            `json:"parent_sha256"`
}

func (p ParentRecoveryPreview) SHA256() string { raw, _ := canonical(p); return digest(raw) }

type ParentRecovery struct {
	Version     int                   `json:"version"`
	Preview     ParentRecoveryPreview `json:"preview"`
	SHA256      string                `json:"sha256"`
	RecoveredAt time.Time             `json:"recovered_at"`
}

func compatibleUnpublishedParent(p, want ParentResult) bool {
	return (p.Version == 0 || p.Version == want.Version) &&
		(p.RunID == "" || p.RunID == want.RunID) && (p.RequestPath == "" || p.RequestPath == want.RequestPath) &&
		(p.RequestSHA256 == "" || p.RequestSHA256 == want.RequestSHA256) && (p.InvocationID == "" || p.InvocationID == want.InvocationID) &&
		(p.TerminalSHA256 == "" || p.TerminalSHA256 == want.TerminalSHA256) && (p.ReceiptSHA256 == "" || p.ReceiptSHA256 == want.ReceiptSHA256) &&
		(p.Assessment == nil || sameJSON(p.Assessment, want.Assessment)) && (p.FailureStage == "" || p.FailureStage == "parent-publication")
}

// Decode complete top-level fields before accepting a truncated publication.
// json.Unmarshal alone leaves a zero value on a syntax error, which would hide
// contrary fields that were already written before the interruption.
func decodeUnpublishedParent(raw []byte, want ParentResult) (ParentResult, bool, error) {
	var old ParentResult
	d := json.NewDecoder(bytes.NewReader(raw))
	token, err := d.Token()
	if errors.Is(err, io.EOF) {
		return old, false, nil
	}
	if err != nil || token != json.Delim('{') {
		return old, false, errors.New("parent is not an interrupted JSON object")
	}
	fields := map[string]json.RawMessage{}
	shape, _ := json.Marshal(want)
	var known map[string]json.RawMessage
	_ = json.Unmarshal(shape, &known)
	for d.More() {
		remaining := bytes.TrimSpace(raw[d.InputOffset():])
		if len(fields) > 0 {
			remaining = bytes.TrimSpace(bytes.TrimPrefix(remaining, []byte(",")))
		}
		if len(remaining) == 0 {
			return old, false, nil
		}
		for key := range known {
			encoded, _ := json.Marshal(key)
			if len(remaining) < len(encoded) && bytes.HasPrefix(encoded, remaining) {
				return old, false, nil
			}
		}
		token, err = d.Token()
		if err != nil {
			// A partial key cannot carry a value. Other syntax errors refuse.
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				return old, false, nil
			}
			return old, false, err
		}
		key, ok := token.(string)
		if !ok || known[key] == nil || fields[key] != nil {
			return old, false, errors.New("parent contains an unknown or repeated field")
		}
		start := d.InputOffset()
		var value json.RawMessage
		if err = d.Decode(&value); err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				// Accept an unwritten value or an exact scalar prefix. A partially
				// written structured assessment cannot conceal conflicting facts.
				partial := bytes.TrimSpace(raw[start:])
				partial = bytes.TrimSpace(bytes.TrimPrefix(partial, []byte(":")))
				if len(partial) == 0 || key != "assessment" && bytes.HasPrefix(known[key], partial) {
					return old, false, nil
				}
			}
			return old, false, errors.New("parent contains an ambiguous or conflicting partial value")
		}
		if bytes.Equal(value, []byte("null")) && key != "assessment" || key == "trajectory_pending" && !bytes.Equal(value, []byte("true")) || key == "version" && !bytes.Equal(value, known[key]) {
			return old, false, errors.New("retained parent contains an explicit conflicting value")
		}
		fields[key] = value
		complete, _ := json.Marshal(fields)
		decoder := json.NewDecoder(bytes.NewReader(complete))
		decoder.DisallowUnknownFields()
		if err = decoder.Decode(&old); err != nil || !compatibleUnpublishedParent(old, want) {
			return old, false, errors.New("retained parent contains conflicting identity, failure or assessment")
		}
	}
	if len(bytes.TrimSpace(raw[d.InputOffset():])) == 0 {
		return old, false, nil
	}
	token, err = d.Token()
	if errors.Is(err, io.EOF) {
		return old, false, nil
	}
	if err != nil || token != json.Delim('}') {
		return old, false, errors.New("parent contains invalid JSON")
	}
	if _, err = d.Token(); !errors.Is(err, io.EOF) {
		return old, false, errors.New("parent contains trailing JSON")
	}
	return old, true, nil
}

func observeUnpublishedParent(dir *os.Root, base string, derived ParentDerivation) (ParentObservation, error) {
	var observed ParentObservation
	path := filepath.Join(base, "parent-result.json")
	info, err := dir.Lstat(path)
	if os.IsNotExist(err) {
		return ParentObservation{Kind: "missing"}, nil
	}
	if err != nil {
		return observed, err
	}
	if info.IsDir() {
		f, err := dir.Open(path)
		if err != nil {
			return observed, err
		}
		defer f.Close()
		names, err := f.Readdirnames(1)
		current, nameErr := dir.Lstat(path)
		if !errors.Is(err, io.EOF) || len(names) != 0 || nameErr != nil || !os.SameFile(info, current) {
			return observed, errors.New("parent publication directory is nonempty or changed")
		}
		return ParentObservation{Kind: "empty-directory"}, nil
	}
	if !info.Mode().IsRegular() || info.Size() > 1<<20 {
		return observed, errors.New("unpublished parent must be a bounded regular file or an empty directory")
	}
	f, err := dir.Open(path)
	if err != nil {
		return observed, err
	}
	defer f.Close()
	opened, err := f.Stat()
	if err != nil || !os.SameFile(info, opened) {
		return observed, errors.New("unpublished parent changed while opening")
	}
	raw, err := io.ReadAll(io.LimitReader(f, (1<<20)+1))
	if err != nil || len(raw) > 1<<20 {
		return observed, errors.New("unpublished parent exceeds its read bound")
	}
	after, err := f.Stat()
	current, nameErr := dir.Lstat(path)
	if err != nil || nameErr != nil || !os.SameFile(opened, current) || int64(len(raw)) != opened.Size() || after.Size() != opened.Size() || after.ModTime() != opened.ModTime() {
		return observed, errors.New("unpublished parent changed while reading")
	}
	old, complete, err := decodeUnpublishedParent(raw, derived.Result)
	if err != nil {
		return observed, err
	}
	if complete && sameJSON(old, derived.Result) {
		return observed, errors.New("complete parent result is already available; use reconciliation")
	}
	if complete {
		want, _ := json.MarshalIndent(old, "", "  ")
		if !bytes.Equal(raw, want) && !bytes.Equal(raw, append(want, '\n')) {
			return observed, errors.New("parent has ambiguous complete JSON; preserve it for explicit inspection")
		}
	}
	observed = ParentObservation{Kind: "incomplete", SHA256: digest(raw)}
	if old.FailureStage == "parent-publication" {
		observed.Kind = "failed-publication"
	}
	return observed, nil
}

func originalParentDigest(s State, root, runID string, result ParentResult) string {
	raw, _ := json.MarshalIndent(result, "", "  ")
	sha := digest(raw)
	// Existing resolutions may have pinned the other already-supported encoding.
	alternate := digest(append(raw, '\n'))
	for _, r := range s.Resolutions {
		if r.Preview.Root == root && r.Preview.RunID == runID && r.Preview.ParentSHA256 == alternate {
			return alternate
		}
	}
	return sha
}

func newParentRecoveryPreview(ctx context.Context, b budget.CycleBinding, s State, root, runID string) (ParentRecoveryPreview, error) {
	var p ParentRecoveryPreview
	derived, err := deriveParentEvidence(ctx, b, s, root, runID)
	if err != nil {
		return p, err
	}
	dir, err := os.OpenRoot(root)
	if err != nil {
		return p, err
	}
	defer dir.Close()
	original, err := observeUnpublishedParent(dir, filepath.Join(".parley-runtime", "trajectory-verification", runID), derived)
	if err != nil {
		return p, err
	}
	raw, _ := canonical(s)
	return ParentRecoveryPreview{1, root, s.Policy.Idea, runID, digest(raw), original, derived, originalParentDigest(s, root, runID, derived.Result)}, nil
}

func readParentRecovery(dir *os.Root, base string) (ParentRecovery, string, error) {
	var r ParentRecovery
	sha, err := readReconciliationJSON(dir, filepath.Join(base, parentRecoveryName), &r)
	return r, sha, err
}
func validateParentRecovery(dir *os.Root, base, root string, s State, derived ParentDerivation, r ParentRecovery) error {
	p := r.Preview
	if r.Version != 1 || p.Version != 1 || !validHash(p.StateSHA256) || r.SHA256 != p.SHA256() || r.RecoveredAt.IsZero() || r.RecoveredAt.Before(derived.CompletedAt) {
		return errors.New("invalid parent recovery record")
	}

	if p.Root != root || p.Idea != s.Policy.Idea || p.RunID != derived.Result.RunID || !filepath.IsAbs(p.Root) || filepath.Join(p.Root, base, "request.json") != derived.Result.RequestPath || !sameJSON(p.Derived, derived) || p.ParentSHA256 != originalParentDigest(s, p.Root, p.RunID, derived.Result) {
		return errors.New("parent recovery differs from original independently executed evidence")
	}
	observed, err := observeUnpublishedParent(dir, base, derived)
	if err != nil {
		return err
	}
	if observed != p.Original {
		return errors.New("original parent publication state changed after recovery")
	}
	return nil
}

// Only the target's exact, previously bound parent facts may be reconstructed
// while obtaining the normal state guard. All other resolutions are mandatory.
func parentRecoveryResolutionCheck(root, runID string) func(context.Context, budget.CycleBinding, State) error {
	return func(ctx context.Context, b budget.CycleBinding, s State) error {
		for _, resolution := range s.Resolutions {
			old := resolution.Preview
			actual, err := readResolutionEvidence(ctx, b, s, old)
			if err != nil && old.Unchanged == nil && old.Root == root && old.RunID == runID && old.RecoverySHA256 == "" {
				dir, e := os.OpenRoot(root)
				if e != nil {
					return e
				}
				_, _, recoveryErr := readParentRecovery(dir, filepath.Join(".parley-runtime", "trajectory-verification", runID))
				dir.Close()
				if !os.IsNotExist(recoveryErr) {
					return err
				}
				p, e := newParentRecoveryPreview(ctx, b, s, root, runID)
				if e != nil {
					return e
				}
				actual = parentPreview(s, root, runID, p.ParentSHA256, "", p.Derived)
				err = nil
			}
			if err != nil {
				return err
			}
			if err = compareReconciledParent(actual, old); err != nil {
				return err
			}
		}
		return nil
	}
}

func syncParentRecovery(dir *os.Root, base string) error {
	f, err := dir.Open(filepath.Join(base, parentRecoveryName))
	if err != nil {
		return err
	}
	err = errors.Join(fsutil.SyncFile(f), f.Close())
	if err != nil {
		return err
	}
	parent, err := dir.Open(base)
	if err != nil {
		return err
	}
	return errors.Join(fsutil.SyncFile(parent), parent.Close())
}
func publishParentRecovery(dir *os.Root, base string, r ParentRecovery) error {
	return publishParentRecoveryWithSync(dir, base, r, syncParentRecovery)
}

func publishParentRecoveryWithSync(dir *os.Root, base string, r ParentRecovery, sync func(*os.Root, string) error) error {
	raw, err := canonical(r)
	if err != nil || len(raw) > 1<<20 {
		return errors.New("parent recovery exceeds its publication bound")
	}
	var token [16]byte
	if _, err = rand.Read(token[:]); err != nil {
		return err
	}
	stage := filepath.Join(base, ".parent-recovery-"+hex.EncodeToString(token[:])+".tmp")
	f, err := dir.OpenFile(stage, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer dir.Remove(stage)
	_, writeErr := f.Write(raw)
	if err = errors.Join(writeErr, fsutil.SyncFile(f), f.Close()); err != nil {
		return err
	}
	if err = dir.Rename(stage, filepath.Join(base, parentRecoveryName)); err != nil {
		return err
	}
	return sync(dir, base)
}

func parentRecovery(ctx context.Context, root, idea, runID, expected string, apply bool, publish func(*os.Root, string, ParentRecovery) error) (ParentRecoveryPreview, error) {
	var result ParentRecoveryPreview
	root, err := canonicalRoot(root)
	if err != nil {
		return result, err
	}
	if !runtimeID(idea) || !runtimeID(runID) || apply && !validHash(expected) {
		return result, errors.New("parent recovery requires exact idea/run and preview identity")
	}
	err = withStateResolutionCheck(ctx, root, idea, parentRecoveryResolutionCheck(root, runID), func(b budget.CycleBinding, _ budget.Snapshot, s State) error {
		dir, err := os.OpenRoot(root)
		if err != nil {
			return err
		}
		defer dir.Close()
		base := filepath.Join(".parley-runtime", "trajectory-verification", runID)
		retained, _, err := readParentRecovery(dir, base)
		if err == nil {
			derived, e := deriveParentEvidence(ctx, b, s, root, runID)
			if e != nil {
				return e
			}
			if e = validateParentRecovery(dir, base, root, s, derived, retained); e != nil {
				return e
			}
			if apply && expected != retained.SHA256 {
				return errors.New("parent recovery replay changed its original preview")
			}
			if apply {
				if e = syncParentRecovery(dir, base); e != nil {
					return e
				}
			}
			result = retained.Preview
			return nil
		}
		if !os.IsNotExist(err) {
			return err
		}
		p, err := newParentRecoveryPreview(ctx, b, s, root, runID)
		if err != nil {
			return err
		}
		if apply {
			if expected != p.SHA256() {
				return errors.New("parent recovery evidence or state changed since preview")
			}
			if err = publish(dir, base, ParentRecovery{1, p, expected, time.Now().UTC()}); err != nil {
				return err
			}
		}
		result = p
		return nil
	})
	if err == nil && result.Version == 0 {
		err = errors.New("trajectory is not active")
	}
	return result, err
}
func PreviewParentRecovery(ctx context.Context, root, idea, runID string) (ParentRecoveryPreview, error) {
	return parentRecovery(ctx, root, idea, runID, "", false, publishParentRecovery)
}
func RecoverParent(ctx context.Context, root, idea, runID, expected string) (ParentRecoveryPreview, error) {
	return parentRecovery(ctx, root, idea, runID, expected, true, publishParentRecovery)
}
