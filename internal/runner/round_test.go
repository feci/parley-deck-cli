package runner

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/store"
)

func TestBuildRoundPromptIncludesPriorAndRound(t *testing.T) {
	prompt := BuildRoundPrompt(
		agents.Discovery{Spec: agents.Spec{ID: "codex"}},
		protocol.IdeaStatus{Slug: "demo"},
		2, "/x/round-02/codex.md", "/x/questions",
		"PRIOR-ROUND-CONTENT-MARKER",
	)
	for _, want := range []string{"cross-review round 2", "round: 2", "PRIOR-ROUND-CONTENT-MARKER", "Responses to other participants"} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("prompt missing %q:\n%s", want, prompt)
		}
	}
}

func TestGatherPriorRounds(t *testing.T) {
	idea := t.TempDir()
	mustWrite(t, filepath.Join(idea, "round-01", "codex.md"), "CODEX-R1")
	mustWrite(t, filepath.Join(idea, "round-01", "claude.md"), "CLAUDE-R1")
	mustWrite(t, filepath.Join(idea, "round-01", "_index.md"), "INDEX-SHOULD-BE-SKIPPED")
	mustWrite(t, filepath.Join(idea, "round-02", "codex.md"), "CODEX-R2")

	got, err := gatherPriorRounds(idea, 3)
	if err != nil {
		t.Fatalf("gatherPriorRounds: %v", err)
	}
	for _, want := range []string{"CODEX-R1", "CLAUDE-R1", "CODEX-R2"} {
		if !strings.Contains(got, want) {
			t.Fatalf("prior digest missing %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "INDEX-SHOULD-BE-SKIPPED") {
		t.Fatal("_index.md must be skipped")
	}
}

func TestValidateRoundArtifactRound2(t *testing.T) {
	dir := t.TempDir()
	good := filepath.Join(dir, "good.md")
	mustWrite(t, good, "---\nagent: codex\nidea: demo\nround: 2\n---\n\n## Summary\nok\n")
	if err := ValidateRoundArtifact(good, "codex", "demo", 2); err != nil {
		t.Fatalf("valid round-2 artifact rejected: %v", err)
	}
	wrongRound := filepath.Join(dir, "wrong.md")
	mustWrite(t, wrongRound, "---\nagent: codex\nidea: demo\nround: 1\n---\n\n## Summary\nok\n")
	if err := ValidateRoundArtifact(wrongRound, "codex", "demo", 2); err == nil {
		t.Fatal("round:1 frontmatter must fail validation for round 2")
	}
	noSections := filepath.Join(dir, "nosec.md")
	mustWrite(t, noSections, "---\nagent: codex\nidea: demo\nround: 2\n---\n\njust prose, no headings\n")
	if err := ValidateRoundArtifact(noSections, "codex", "demo", 2); err == nil {
		t.Fatal("artifact without section headings must fail")
	}
}

func TestRunRoundCrossReviewWithHeadlessAgent(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	idea, err := protocol.CreateIdea(root, "Cross review task", []string{"fake"})
	if err != nil {
		t.Fatal(err)
	}
	// Seed a round-01 artifact so the cross-review prompt has prior content.
	mustWrite(t, filepath.Join(idea.Path, "round-01", "fake.md"), "---\nagent: fake\nidea: "+idea.Slug+"\nround: 1\n---\n\n## Summary\nr1\n")

	runStore := store.New(filepath.Join(root, protocol.DeckDir, "runs", "xr-run"))
	results := RunRound(context.Background(), Options{
		Root:  root,
		RunID: "xr-run",
		Idea:  idea,
		Round: 2,
		Agents: []agents.Discovery{
			{
				Spec: agents.Spec{
					ID:           "fake",
					HeadlessArgs: []string{"-test.run=TestFakeRound2AgentHelper", "--", "parley-fake-round2"},
					PromptMode:   agents.PromptStdin,
				},
				Path:  os.Args[0],
				Found: true,
			},
		},
		Timeout: 5 * time.Second,
		Store:   runStore,
	})
	if len(results) != 1 || !results[0].ArtifactOK || results[0].ExitError != "" {
		t.Fatalf("unexpected result: %+v", results)
	}
	if _, err := os.Stat(filepath.Join(idea.Path, "round-02", "fake.md")); err != nil {
		t.Fatalf("round-02 artifact missing: %v", err)
	}
}

// TestFakeRound2AgentHelper is a subprocess helper that emits a round-02
// cross-review artifact at the path named in the prompt.
func TestFakeRound2AgentHelper(t *testing.T) {
	if !hasArg("parley-fake-round2") {
		return
	}
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		t.Fatal(err)
	}
	out := regexp.MustCompile(`(?m)^- Create exactly this file and no other protocol artifact: (.+)$`).FindStringSubmatch(string(input))
	idea := regexp.MustCompile(`(?m)^idea: (.+)$`).FindStringSubmatch(string(input))
	if len(out) != 2 || len(idea) != 2 {
		t.Fatalf("prompt missing output path or idea:\n%s", string(input))
	}
	body := "---\nagent: fake\nidea: " + idea[1] + "\nround: 2\n---\n\n## Summary\nx\n\n## Responses to other participants\ny\n\n## Refined position\nz\n\n## Remaining disagreements\nnone\n"
	if err := os.WriteFile(out[1], []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	os.Exit(0)
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
