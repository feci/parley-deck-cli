package quota

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LegacyManualNotice is retained display history, never owner authority.
func (b Batch) LegacyManualNotice() string {
	if b.Owner == nil || b.Owner.Authority != nil {
		return ""
	}
	return fmt.Sprintf("---\nfrom: parley\nto: user\nidea: %s\nblocking: no\ntransition: %s\n---\n\nOwner-confirmed membership/policy revision. Current participants: %v. Policy: %+v.\n", b.Idea, b.ID, b.Decision.After, b.Policy)
}

// PublishNotice preserves owner-owned copies and uses checked exclusive publication
// on safe destinations, creating a missing inbox. Callers report errors as diagnostics;
// these display files never authorize membership or block a completed transition.
func PublishNotice(inbox, name, text string) error {
	found, _, err := InspectNoticeForPublication(inbox, name, "")
	if err != nil || found {
		return err
	}
	return DurableWrite(filepath.Join(inbox, name), []byte(text), true)
}

func NoticeDiagnostic(id string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "quota notice %s: publication diagnostic (non-blocking; delivery unconfirmed): %v\n", id, err)
	}
}

// Inspect only for publication choices. Existing regular copies are owner-owned
// regardless of content. The historical label is solely a clarification hint.
func InspectNoticeForPublication(inbox, name, legacy string) (found, historical bool, err error) {
	// Check the deck and both inbox directories before looking at destinations,
	// so an archived/ or inbox/ symlink is not followed even when a file is absent.
	for _, dir := range []string{filepath.Dir(inbox), inbox, filepath.Join(inbox, "archived")} {
		st, e := os.Lstat(dir)
		if os.IsNotExist(e) {
			continue
		}
		if e != nil {
			return found, historical, e
		}
		if !st.IsDir() {
			return found, historical, fmt.Errorf("unsafe quota notice directory: %s", dir)
		}
	}
	for _, dir := range []string{inbox, filepath.Join(inbox, "archived")} {
		path := filepath.Join(dir, name)
		st, e := os.Lstat(path)
		if os.IsNotExist(e) {
			continue
		}
		if e != nil {
			return found, historical, e
		}
		if !st.Mode().IsRegular() {
			return found, historical, fmt.Errorf("unsafe quota notice destination: %s", path)
		}
		found = true
		if legacy != "" {
			raw, e := os.ReadFile(path)
			if e != nil {
				return found, historical, e
			}
			historical = historical || strings.HasPrefix(string(raw), legacy)
		}
	}
	return found, historical, nil
}
