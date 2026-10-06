package quota

import (
	"fmt"
	"os"
	"path/filepath"
)

// LegacyManualNotice is retained display history, never owner authority.
func (b Batch) LegacyManualNotice() string {
	if b.Owner == nil || b.Owner.Authority != nil {
		return ""
	}
	return fmt.Sprintf("---\nfrom: parley\nto: user\nidea: %s\nblocking: no\ntransition: %s\n---\n\nOwner-confirmed membership/policy revision. Current participants: %v. Policy: %+v.\n", b.Idea, b.ID, b.Decision.After, b.Policy)
}

// InspectNotice checks every extant live/archived copy without repairing it.
// Absence is normal inbox deletion, not corruption. Callers decide whether a
// durable receipt proves prior delivery or checked publication is still needed.
func InspectNotice(inbox, name, want, legacy string) (paths []string, historical bool, err error) {
	for _, dir := range []string{inbox, filepath.Join(inbox, "archived")} {
		path := filepath.Join(dir, name)
		st, e := os.Lstat(path)
		if os.IsNotExist(e) {
			continue
		}
		if e != nil {
			return nil, false, e
		}
		if !st.Mode().IsRegular() {
			return nil, false, fmt.Errorf("invalid quota notice: %s", path)
		}
		raw, e := os.ReadFile(path)
		if e != nil {
			return nil, false, e
		}
		if string(raw) != want && (legacy == "" || string(raw) != legacy) {
			return nil, false, fmt.Errorf("corrupt or mismatched quota notice: %s", path)
		}
		paths = append(paths, path)
		historical = historical || (legacy != "" && string(raw) == legacy)
	}
	return paths, historical, nil
}
