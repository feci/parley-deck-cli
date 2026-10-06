package quota

import (
	"fmt"
	"os"
	"path/filepath"
)

// ReadApplied distinguishes absence (checked replay may recover it) from a
// present but contradictory receipt. Both require full checked replay before a
// new receipt may be published; neither is proof of completed publication.
func ReadApplied(ideaDir, id string) (bool, error) {
	path := filepath.Join(ideaDir, "quota-applied", id)
	st, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if !st.Mode().IsRegular() {
		return false, fmt.Errorf("invalid quota applied receipt: %s", path)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	if string(raw) != id+"\n" {
		return false, nil
	}
	return true, nil
}
