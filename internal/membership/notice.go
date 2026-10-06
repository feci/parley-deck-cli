package membership

import (
	"path/filepath"

	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
)

// An applied receipt proves the publication step completed. The recipient may
// subsequently archive or delete the notice. Before a receipt exists, a checked
// archived copy proves delivery; if both copies disappeared, replay publishes
// once again rather than permanently gating an ordinary inbox deletion.
func publishNotice(root, ideaDir string, b quota.Batch, applied bool) error {
	inbox := filepath.Join(root, protocol.DeckDir, "inbox")
	name := "parley-to-user_" + b.ID + ".md"
	want := b.Notice()
	legacy := b.LegacyManualNotice()
	manual := b.Owner != nil && b.Owner.Authority == nil
	found, historical, err := checkedNotice(inbox, name, want, legacy)
	if err != nil {
		return err
	}
	if !found && !applied {
		if err := quota.DurableWrite(filepath.Join(inbox, name), []byte(want), true); err != nil {
			return err
		}
	}
	if !manual {
		return nil
	}

	// Old manual revisions had an inaccurate display label, never owner
	// authority. Preserve those bytes and deliver a separate clarification once.
	// Its own receipt also survives archival/deletion independently of the old
	// transition receipt, which predates this clarification.
	correctionID := b.ID + "-manual-authority"
	corrected, err := quota.ReadApplied(ideaDir, correctionID)
	if err != nil {
		return err
	}
	correctionFound, _, err := checkedNotice(inbox, "parley-to-user_"+correctionID+".md", want, "")
	if err != nil {
		return err
	}
	if corrected {
		return quota.SyncPath(filepath.Join(ideaDir, "quota-applied", correctionID))
	}
	if !historical && !correctionFound {
		return nil
	}
	if !correctionFound {
		if err := quota.DurableWrite(filepath.Join(inbox, "parley-to-user_"+correctionID+".md"), []byte(want), true); err != nil {
			return err
		}
	}
	return quota.DurableWrite(filepath.Join(ideaDir, "quota-applied", correctionID), []byte(correctionID+"\n"), false)
}

// Validate every existing copy, even with an applied receipt. A corrupt notice
// cannot serve as publication proof or be silently accepted as a completed step.
func checkedNotice(inbox, name, want, legacy string) (found, historical bool, err error) {
	paths, historical, err := quota.InspectNotice(inbox, name, want, legacy)
	if err != nil {
		return false, false, err
	}
	for _, path := range paths {
		if err := quota.SyncPath(path); err != nil {
			return false, false, err
		}
	}
	return len(paths) > 0, historical, nil
}
