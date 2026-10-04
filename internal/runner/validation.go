package runner

import "parley-deck-cli/internal/protocol"

// Round validators are shared with transition recovery; both paths must apply
// the identical artifact contract before a terminal round can be completed.
func ValidateRoundArtifact(path, agentID, ideaSlug string, round int) error {
	return protocol.ValidateParticipantRoundArtifact(path, agentID, ideaSlug, round)
}
func ValidateRoundOneArtifact(path, agentID, ideaSlug string) error {
	return protocol.ValidateParticipantRoundOneArtifact(path, agentID, ideaSlug)
}
