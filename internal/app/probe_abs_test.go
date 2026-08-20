package app

import (
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/agents"
)

// The call site, not the helper: `agents verify --full` used to build its probe directory by
// joining a possibly-relative root, and the resulting relative path is what reached the agent.
// hermes resolves relative paths against $HOME whatever the process cwd is (measured; `--in` and
// `--no-restore-cwd` do not change it), so the file landed in ~/parley-deck/... while the verifier
// looked in the repository and reported "did not create <path>".
func TestProbeDirIsAbsoluteEvenWhenRootIsRelative(t *testing.T) {
	for _, root := range []string{".", "", "some/relative/deck"} {
		got := probeDirFor(root, "RUN123")
		if !filepath.IsAbs(got) {
			t.Errorf("probeDirFor(%q) = %q, which is relative; an agent that does not share parley's cwd writes it elsewhere", root, got)
		}
		if !strings.HasSuffix(got, filepath.Join("meta", "runtime-probes", "RUN123")) {
			t.Errorf("probeDirFor(%q) = %q, lost its runtime-probes suffix", root, got)
		}
	}
}

// An absolute root must be preserved exactly, not re-rooted at the process cwd.
func TestProbeDirPreservesAnAbsoluteRoot(t *testing.T) {
	root := t.TempDir()
	got := probeDirFor(root, "RUN123")
	if !strings.HasPrefix(got, root) {
		t.Fatalf("probeDirFor(%q) = %q; an absolute root must be preserved", root, got)
	}
}

// Whatever the caller computed, the prompt must hand the agent that exact path.
func TestProbePromptNamesTheProbePath(t *testing.T) {
	out := filepath.Join(probeDirFor(t.TempDir(), "RUN123"), "hermes.md")
	prompt := probePrompt(agents.Discovery{Spec: agents.Spec{ID: "hermes"}}, out, "# sentinel")
	if !strings.Contains(prompt, out) {
		t.Fatalf("probe prompt does not name %q", out)
	}
}
