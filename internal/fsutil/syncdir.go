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

// PinLstat stats name through an os.Root of its directory and returns an
// EAGERLY pinned FileInfo. Row 30 (hosted-confirmed 36443342532): on Windows
// an os.Stat/os.Lstat FileInfo defers its identity to the first
// os.SameFile call, re-opening BY PATH at comparison time — a deferred pin
// re-resolves to whatever now lives at the path and is structurally
// incapable of detecting a same-path replacement. A Root-derived stat
// captures the identity (volume + file index) at stat time. The same-shaped
// invariant therefore binds every os.SameFile guard: the PINNED ("before")
// operand must be handle- or Root-derived, never os.Stat/os.Lstat.
//
// Containment honesty (claude-1 consult §7, accepted): only the FINAL
// component is contained — filepath.Dir(name) is resolved with ambient
// authority (symlinks/junctions/.. in the parent chain are followed), and no
// origin identity is pinned. Callers needing parent-resolution containment
// must hold their own root and pin against it. On POSIX this is the same
// lstat.
func PinLstat(name string) (os.FileInfo, error) {
	root, err := os.OpenRoot(filepath.Dir(name))
	if err != nil {
		return nil, err
	}
	defer root.Close()
	return root.Lstat(filepath.Base(name))
}

// OpenPinned opens name for reading and returns BOTH the open file and its
// eagerly pinned identity, taken through ONE directory handle so pin and
// open cannot observe different directories: the pin is a Root.Lstat (row 30
// eager identity; see PinLstat) and the open is Root.Open, which on Windows
// carries FILE_SHARE_DELETE so a concurrent atomic replacement
// (ReplaceSyncedFile / WriteFileAtomic) is not blocked by this reader
// (§D.6/AC-LOCK-1). Callers compare with os.SameFile(pinned, opened) — that
// comparison now works on every platform. Containment scope as documented on
// PinLstat: final component only; the parent resolution is the caller's.
func OpenPinned(name string) (*os.File, os.FileInfo, error) {
	root, err := os.OpenRoot(filepath.Dir(name))
	if err != nil {
		return nil, nil, err
	}
	defer root.Close()
	base := filepath.Base(name)
	info, err := root.Lstat(base)
	if err != nil {
		return nil, nil, err
	}
	f, err := root.Open(base)
	if err != nil {
		return nil, nil, err
	}
	return f, info, nil
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
