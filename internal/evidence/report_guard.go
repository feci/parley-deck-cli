package evidence

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"parley-deck-cli/internal/budget"
)

var (
	ErrReportChanged   = errors.New("evidence report changed; rerun checks before another verification")
	ErrReportFinalized = errors.New("completed implementation evidence cannot be replaced")
)

// ReportWriter owns a cooperative publication transaction. Methods are valid
// only within WithReportWriter's callback. It coordinates report replacement,
// the generated evidence table, verifier receipts and the final status write.
// Uncooperative filesystem edits and same-UID tampering remain outside this
// boundary. No lock is held while an independent model process is running.
type ReportWriter struct {
	dir    string
	mu     sync.Mutex
	active bool
}

func WithReportWriter(ctx context.Context, ideaDir string, fn func(*ReportWriter) error) error {
	dir, err := filepath.Abs(ideaDir)
	if err != nil {
		return err
	}
	dir, err = filepath.EvalSymlinks(dir)
	if err != nil {
		return err
	}
	guardDir, err := reportGuardDirectory(dir)
	if err != nil {
		return err
	}
	// A missing caller deadline must not leave a CLI waiting forever on another
	// live process. Explicit earlier caller cancellation still wins.
	wait, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	release, err := budget.AcquireResourceGuard(wait, guardDir)
	if err != nil {
		return fmt.Errorf("evidence publication guard: %w", err)
	}
	defer release()
	w := &ReportWriter{dir: dir, active: true}
	defer func() { w.mu.Lock(); w.active = false; w.mu.Unlock() }()
	return fn(w)
}

// Metadata stays outside the tested Git worktree. Linked worktrees get their
// own report identities; aliases of one actual idea directory share a guard.
// Without Git, a stable local runtime origin is created before check hashing.
// It is not silently excluded from a non-Git digest.
func reportGuardDirectory(dir string) (string, error) {
	sum := sha256.Sum256([]byte(strings.ToLower(dir)))
	identity := hex.EncodeToString(sum[:])
	for ancestor := dir; ; ancestor = filepath.Dir(ancestor) {
		marker := filepath.Join(ancestor, ".git")
		info, err := os.Lstat(marker)
		if err == nil {
			gitDir := marker
			if !info.IsDir() {
				if !info.Mode().IsRegular() || info.Size() > 16<<10 {
					return "", errors.New("invalid Git administration marker for evidence guard")
				}
				f, err := os.Open(marker)
				if err != nil {
					return "", err
				}
				data, readErr := io.ReadAll(io.LimitReader(f, (16<<10)+1))
				f.Close()
				if readErr != nil {
					return "", readErr
				}
				value := strings.TrimSpace(string(data))
				if len(data) > 16<<10 || !strings.HasPrefix(value, "gitdir: ") {
					return "", errors.New("invalid Git directory pointer for evidence guard")
				}
				gitDir = strings.TrimSpace(strings.TrimPrefix(value, "gitdir: "))
				if gitDir == "" {
					return "", errors.New("empty Git directory pointer")
				}
				if !filepath.IsAbs(gitDir) {
					gitDir = filepath.Join(ancestor, gitDir)
				}
			}
			gitDir, err = filepath.EvalSymlinks(gitDir)
			if err != nil {
				return "", err
			}
			return filepath.Join(gitDir, "parley-evidence-guards", identity), nil
		}
		if !os.IsNotExist(err) {
			return "", err
		}
		if filepath.Dir(ancestor) == ancestor {
			break
		}
	}
	return filepath.Join(dir, ".parley-runtime", "evidence-report"), nil
}

// Save replaces the current report under the shared guard. Once Complete has
// published its authorized status transition, delayed writers must refuse.
func (w *ReportWriter) Save(report *Report) ([]byte, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.save(report)
}

// SaveIfUnchanged performs the byte comparison and replacement within the same
// transaction. A changed, missing, or unreadable original is never overwritten.
func (w *ReportWriter) SaveIfUnchanged(original []byte, report *Report) ([]byte, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.active {
		return nil, errors.New("evidence writer transaction has ended")
	}
	current, err := os.ReadFile(ReportPath(w.dir))
	if err != nil {
		return nil, errors.Join(ErrReportChanged, err)
	}
	if !bytes.Equal(original, current) {
		return nil, ErrReportChanged
	}
	return w.save(report)
}

func (w *ReportWriter) save(report *Report) ([]byte, error) {
	if !w.active {
		return nil, errors.New("evidence writer transaction has ended")
	}
	doc, err := os.ReadFile(filepath.Join(w.dir, "IMPLEMENTATION.md"))
	if err == nil {
		_, from, err := TransitionStatusToComplete(doc)
		if from == "complete" {
			return nil, ErrReportFinalized
		}
		if err != nil {
			return nil, fmt.Errorf("cannot establish publication status: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	return saveReport(w.dir, report)
}
