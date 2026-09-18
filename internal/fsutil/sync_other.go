//go:build !darwin

package fsutil

import "os"

// SyncFile requires the platform's normal file synchronization to succeed.
func SyncFile(file *os.File) error { return file.Sync() }
