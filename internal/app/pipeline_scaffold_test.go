package app

import (
	"strings"
	"testing"

	"parley-deck-cli/internal/protocol"
)

// Review round 1, @codex-1 MAJOR: the scaffold `consensus finalize` writes carries
// `status: final` from its FIRST step, so a frontmatter-only completion test called an empty
// outline a completed block and the pipeline advanced past it.
func TestScaffoldWithStatusFinalIsNotACompletedBlock(t *testing.T) {
	var scaffold strings.Builder
	scaffold.WriteString("---\nidea: demo\nstatus: final\n---\n\n")
	for _, section := range protocol.RequiredFinalSections {
		scaffold.WriteString(section + "\n\n<fill in, or write N/A if this idea is trivial or design-only>\n\n")
	}
	if isFinalized(scaffold.String()) {
		t.Fatal("an unwritten scaffold declaring status: final was treated as a completed block")
	}

	var written strings.Builder
	written.WriteString("---\nidea: demo\nstatus: final\n---\n\n")
	for _, section := range protocol.RequiredFinalSections {
		written.WriteString(section + "\n\nReal content.\nSecond line.\nThird line.\n\n")
	}
	if !isFinalized(written.String()) {
		t.Fatal("a written FINAL was not recognised as complete")
	}

	// No status at all is still not final.
	if isFinalized("---\nidea: demo\n---\n\n## Final plan / specification\n\nx\ny\nz\n") {
		t.Fatal("an artifact without status: final was treated as final")
	}
}
