package protocol

import (
	"fmt"
	"os"
)

// RoundOneRequiredSection is the §15.6(a) duty: a round-1 artifact must enumerate what the
// proposal builds by hand and, for each, what the toolchain already ships.
const RoundOneRequiredSection = "## Existing alternatives"

// ValidateRoundOneArtifact enforces §15.6(a) on a round-1 design artifact.
//
// It exists because design rounds had no existence-or-shape gate at all, unlike review rounds
// (ValidateReviewArtifact). That asymmetry let a participant CLI exit 0 having written nothing,
// about ten times in one idea, with only a human reading a directory listing between a phantom
// participant and the record.
//
// The section must be a real heading with at least one non-blank line under it. A substring
// mention or an empty heading is a rubber-stamp, not work shown -- the same rule the review gate
// applies, for the same reason: the half a scanner cannot check decays to prose compliance rates
// even when the prompt carries it.
func ValidateRoundOneArtifact(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("round-01 artifact unreadable: %w", err)
	}
	if !HasNonEmptySection(string(raw), RoundOneRequiredSection) {
		return fmt.Errorf("round-01 artifact %s is missing a non-empty %q section (§15.6a)", path, RoundOneRequiredSection)
	}
	return nil
}
