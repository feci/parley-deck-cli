//go:build windows

package fsutil

// SyncDir is the fail-closed named refusal (FINAL §B/AC-DUR-1): the rooted
// sites' refusal emitter. It never attempts a mechanism — a green API call is
// never cited as a guarantee — and every Unix-only directory sync stays in
// syncdir_unix.go behind //go:build !windows.
func SyncDir(d *DirHandle) error {
	return ErrDirEntryDurabilityUnsupported
}
