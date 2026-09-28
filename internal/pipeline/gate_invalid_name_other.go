//go:build !windows

package pipeline

// errInvalidName is always false on POSIX: every byte but NUL and '/' is
// representable in a filename, so a legacy raw gate name is always readable.
func errInvalidName(err error) bool { return false }
