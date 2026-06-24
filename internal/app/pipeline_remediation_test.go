package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/pipeline"
)

// LE-10: a watcher-opened remediation idea is a non-active candidate, not a round-01
// idea claiming an empty quorum (the non-solo Phase-0 invariant).
func TestOpenRemediationIdeaIsCandidate(t *testing.T) {
	deck := t.TempDir()
	b := pipeline.Breach{Signal: "latency_p99", Target: "api", Threshold: "<200ms", Observed: "350ms", Class: "slo"}
	slug, err := openRemediationIdea(deck, "mypipe", b, "abc123", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(deck, "ideas", slug, "00-prompt.md"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(data)
	if !strings.Contains(s, "status: candidate") {
		t.Fatal("LE-10: a watcher-opened idea must be status: candidate")
	}
	if strings.Contains(s, "status: round-01") {
		t.Fatal("LE-10: must not auto-open at round-01")
	}
	if strings.Contains(s, "participants: []") {
		t.Fatal("LE-10: must not claim an empty participant quorum")
	}
	if !strings.Contains(s, "## Promotion") {
		t.Fatal("expected a Promotion note explaining the candidate→round-01 gate")
	}
}
