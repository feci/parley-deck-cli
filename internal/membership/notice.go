package membership

import (
	"path/filepath"

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
		quota.NoticeDiagnostic(b.ID, publishMissingNotice(inbox, name, b.Notice()))
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
	_, historical, err := quota.InspectNoticeForPublication(inbox, name, b.LegacyManualNotice())
	quota.NoticeDiagnostic(b.ID, err)
	if err != nil && !historical {
		return nil
	}
	correctionName := "parley-to-user_" + correctionID + ".md"
	found, _, err := quota.InspectNoticeForPublication(inbox, correctionName, "")
	quota.NoticeDiagnostic(correctionID, err)
	if !historical && (err != nil || !found) {
		return nil
	}
	quota.NoticeDiagnostic(correctionID, publishMissingNotice(inbox, correctionName, b.Notice()))
	return quota.DurableWrite(filepath.Join(ideaDir, "quota-applied", correctionID), []byte(correctionID+"\n"), false)
}

func publishMissingNotice(inbox, name, text string) error {
	if err := projectionFault("notice"); err != nil {
		return err
	}
	return quota.PublishNotice(inbox, name, text)
}

// PublishKickoffNotice uses the same completed-attempt receipt as later batches.
// A crash before the receipt retries publication, preserving any owner-owned copy.
func PublishKickoffNotice(root, ideaDir string, k *quota.Kickoff) error {
	if k == nil || k.Transition == nil {
		return nil
	}
	id := k.Transition.ID
	applied, err := quota.ReadApplied(ideaDir, id)
	if err != nil {
		return err
	}
	if applied {
		return nil
	}
	quota.NoticeDiagnostic(id, publishMissingNotice(filepath.Join(root, protocol.DeckDir, "inbox"), "parley-to-user_"+id+".md", k.Notice()))
	return quota.DurableWrite(filepath.Join(ideaDir, "quota-applied", id), []byte(id+"\n"), false)
}
