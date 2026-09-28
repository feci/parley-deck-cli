package fsutil

import (
	"errors"
	"fmt"
	"os"
)

// ErrDirEntryDurabilityUnsupported is the FINAL §B/AC-DUR-1 named,
// fail-closed refusal: on Windows there is no directory-entry durability
// mechanism that preserves os.Root containment (FlushFileBuffers is
// foreclosed on read-only directory handles — H2 — and path-based
// alternatives were rejected by the owner), so the entry-durability barrier
// refuses rather than continuing without the guarantee.
var ErrDirEntryDurabilityUnsupported = errors.New("directory-entry durability is not available on Windows while os.Root containment is preserved; refusing rather than continuing without the guarantee (run the operation on macOS/Linux)")

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
