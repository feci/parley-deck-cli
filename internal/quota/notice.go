package quota

import "fmt"

// LegacyManualNotice is retained display history, never owner authority.
func (b Batch) LegacyManualNotice() string {
	if b.Owner == nil || b.Owner.Authority != nil {
		return ""
	}
	return fmt.Sprintf("---\nfrom: parley\nto: user\nidea: %s\nblocking: no\ntransition: %s\n---\n\nOwner-confirmed membership/policy revision. Current participants: %v. Policy: %+v.\n", b.Idea, b.ID, b.Decision.After, b.Policy)
}
