package fsutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ErrDirEntryDurabilityUnsupported is the FINAL §B/AC-DUR-1 named,
// fail-closed refusal: on Windows there is no directory-entry durability
// mechanism that preserves os.Root containment (FlushFileBuffers is
// foreclosed on read-only directory handles — H2 — and path-based
// alternatives were rejected by the owner), so the entry-durability barrier
// refuses rather than continuing without the guarantee.
var ErrDirEntryDurabilityUnsupported = errors.New("directory-entry durability is not available on Windows while os.Root containment is preserved; refusing rather than continuing without the guarantee (run the operation on macOS/Linux)")

// StageSuffix is the deterministic-stage derivation of FINAL §A: the stage
// name is always final+StageSuffix, created with O_EXCL and NOT removed on
// the error path, so an interrupted publication leaves a visible blocker
// whose existence prevents replay. Requirement (a): the suffix sits outside
// every valid-final grammar of the §B plain rows — A1's fixed artifact names
// (request.json, parent-result.json), B1's fixed runtime directories
// (.parley-runtime, trajectory-verification) and B2's telemetry IDs (UUID v4:
// lowercase hex and dashes only, no dots). The dotted suffix can therefore
// never be part of a valid final name; the pin is durable_grammar_test.go.
const StageSuffix = ".parley-staging"

// OpenSharedDelete opens name for reading with delete-share semantics:
// routed through an os.Root of its directory, which on Windows passes
// FILE_SHARE_DELETE so a concurrent atomic replacement (ReplaceSyncedFile /
// WriteFileAtomic) is not blocked by this reader (§D.6/AC-LOCK-1 share-flag-
// correct opens; §D.6 prefers exactly this rooting over a raw CreateFile
// helper). On POSIX root.Open is the same read-only open. Callers must have
// verified the path is a regular file (the rooted open does not follow
// symlinks that escape the directory).
func OpenSharedDelete(name string) (*os.File, error) {
	root, err := os.OpenRoot(filepath.Dir(name))
	if err != nil {
		return nil, err
	}
	f, err := root.Open(filepath.Base(name))
	root.Close()
	if err != nil {
		return nil, err
	}
	return f, nil
}

// DirHandle is the named directory-handle type of the entry-durability
// contract: directory handles reach their barrier ONLY as a DirHandle through
// SyncDir — the named type keeps them out of SyncFile at compile time, making
// the contract the audit of record rather than a grep (FINAL §B/AC-DUR-1).
type DirHandle struct {
	file *os.File
}

// DirHandleOf wraps an already-open handle, failing closed unless it is a
// real directory (callers hold os.Root containment and open "." or a rooted
// subpath themselves; the constructor only re-verifies what they assert).
func DirHandleOf(f *os.File) (*DirHandle, error) {
	if f == nil {
		return nil, errors.New("fsutil.DirHandleOf: nil handle")
	}
	info, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("fsutil.DirHandleOf: stat: %w", err)
	}
	if !info.IsDir() {
		return nil, errors.New("fsutil.DirHandleOf: handle is not a directory")
	}
	return &DirHandle{file: f}, nil
}

// File exposes the underlying handle for Close by the owning caller.
func (d *DirHandle) File() *os.File { return d.file }
