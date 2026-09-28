package trajectory

import (
	"errors"
	"os"
	"parley-deck-cli/internal/fsutil"
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
	// Row-30 AFTER rule: a Root-derived stat — the lazy os.Lstat's implicit
	// dwShareMode=0 identity open would refuse concurrent atomic replacements
	// of this file (the claude-1 delete-pending consult finding).
	named, err := fsutil.PinLstat(file)
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
