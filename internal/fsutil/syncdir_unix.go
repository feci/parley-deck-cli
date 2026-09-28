//go:build !windows

package fsutil

// SyncDir makes a directory's entries durable on POSIX: same semantics the
// previous inline SyncFile(dir) calls had (darwin's F_FULLFSYNC fallback
// included), byte-identical at every call site (AC-DUR-5).
func SyncDir(d *DirHandle) error {
	return SyncFile(d.file)
}
