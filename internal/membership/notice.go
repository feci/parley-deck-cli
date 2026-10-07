package membership

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
)

// An applied receipt ends ordinary notice publication. It records a completed
// attempt, not proof of delivery. Notice copies belong to the owner; their bytes
// and paths never decide membership. Only receipt failures are returned here.
func publishNotice(root, ideaDir string, b quota.Batch, applied bool) error {
	inbox := filepath.Join(root, protocol.DeckDir, "inbox")
	name := "parley-to-user_" + b.ID + ".md"
	if !applied {
		noticeDiagnostic(b.ID, publishMissingNotice(inbox, name, b.Notice()))
	}
	if b.Owner == nil || b.Owner.Authority != nil {
		return nil
	}

	// Historical manual revisions could carry an inaccurate authority label.
	// The clarification has its own receipt, independent of the old batch's
	// receipt. Neither display file becomes authority, including after owner edits.
	correctionID := b.ID + "-manual-authority"
	corrected, err := quota.ReadApplied(ideaDir, correctionID)
	if err != nil {
		return err
	}
	if corrected {
		return quota.SyncPath(filepath.Join(ideaDir, "quota-applied", correctionID))
	}
	_, historical, err := inspectNotice(inbox, name, b.LegacyManualNotice())
	noticeDiagnostic(b.ID, err)
	correctionName := "parley-to-user_" + correctionID + ".md"
	found, _, err := inspectNotice(inbox, correctionName, "")
	noticeDiagnostic(correctionID, err)
	if !historical && !found && err == nil {
		return nil
	}
	noticeDiagnostic(correctionID, publishMissingNotice(inbox, correctionName, b.Notice()))
	return quota.DurableWrite(filepath.Join(ideaDir, "quota-applied", correctionID), []byte(correctionID+"\n"), false)
}

func noticeDiagnostic(id string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "quota notice %s: publication diagnostic (non-blocking; delivery unconfirmed): %v\n", id, err)
	}
}

func publishMissingNotice(inbox, name, text string) error {
	if err := projectionFault("notice"); err != nil {
		return err
	}
	found, _, err := inspectNotice(inbox, name, "")
	if err != nil || found {
		return err
	}
	return quota.DurableWrite(filepath.Join(inbox, name), []byte(text), true)
}

// Inspect only for publication choices. Existing regular copies are owner-owned
// regardless of content. The historical label is solely a clarification hint.
func inspectNotice(inbox, name, legacy string) (found, historical bool, err error) {
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
