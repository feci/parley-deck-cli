package trajectory

import (
	"errors"
	"os"
)

func sameSnapshotFile(a, b os.FileInfo) bool {
	return a != nil && b != nil && os.SameFile(a, b) && a.Mode() == b.Mode() && a.Size() == b.Size()
}

// Called only after every full archive/canonical/tree hash and link check.
// A timestamp transition alone permits one complete strict reread. Any other
// metadata change, or a second transition, refuses. No sleep/retry-until-green.
func checkSnapshotFileStability(f *os.File, file string, initial, opened os.FileInfo, allowTimestampTransition bool, recheck func() error) error {
	ended, err := f.Stat()
	if err != nil {
		return err
	}
	named, err := os.Lstat(file)
	if err != nil {
		return err
	}
	if !sameSnapshotFile(initial, opened) || !sameSnapshotFile(initial, ended) || !sameSnapshotFile(initial, named) {
		return errors.New("snapshot identity, size or mode changed during verification")
	}
	stable := initial.ModTime().Equal(opened.ModTime()) && opened.ModTime().Equal(ended.ModTime()) && ended.ModTime().Equal(named.ModTime())
	if stable {
		return nil
	}
	if !allowTimestampTransition {
		return errors.New("snapshot timestamp changed during stable verification")
	}
	return recheck()
}
