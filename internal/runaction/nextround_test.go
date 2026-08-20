package runaction

import (
	"strings"
	"testing"
)

// codex-1/F11: the CLI labelled the action "open round-02 of <idea>" and then printed
// `parley run --auto "continue <idea>"`. `run` CREATES a new timestamped idea and starts round-01,
// so following the recommended command forked an unrelated idea and left the original stalled.
func TestOpenNextRoundSurfacesContinueNotRun(t *testing.T) {
	got := Command(NextAction{Kind: KindOpenNextRound, IdeaSlug: "my-idea", RunID: "run-7"}, "", "")
	if strings.Contains(got, "parley run") {
		t.Fatalf("advancing an existing idea must not surface `parley run`, which creates a new one: %q", got)
	}
	if !strings.HasPrefix(got, "parley continue ") {
		t.Fatalf("want a continue command, got %q", got)
	}

	// With no run id, the idea slug is the handle.
	got = Command(NextAction{Kind: KindOpenNextRound, IdeaSlug: "my-idea"}, "", "")
	if got != "parley continue my-idea" {
		t.Fatalf("got %q", got)
	}

	// No idea at all: nothing to suggest.
	if got := Command(NextAction{Kind: KindOpenNextRound}, "", ""); got != "" {
		t.Fatalf("want empty, got %q", got)
	}
}
