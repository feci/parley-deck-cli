//go:build !windows

package fsutil

import "time"

// TEMPORARY replace-failure diagnostic, POSIX stub: the registry exists so
// cross-platform callers compile; nothing samples it off-Windows.
func DiagReaderEnter(path string, at time.Time) {}
func DiagReaderExit(path string)                {}
