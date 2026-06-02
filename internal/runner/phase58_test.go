package runner

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/store"
)

func TestRunImplementationWritesArtifact(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	idea, err := protocol.CreateIdea(root, "Impl task", []string{"builder"})
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(idea.Path, "FINAL.md"), "---\nstatus: final\n---\n\nplan\n")

	res := RunImplementation(context.Background(), Options{
		Root:  root,
		RunID: "impl-run",
		Idea:  idea,
		Agents: []agents.Discovery{{
			Spec:  agents.Spec{ID: "builder", HeadlessArgs: []string{"-test.run=TestFakeImplHelper", "--", "parley-fake-impl"}, PromptMode: agents.PromptStdin},
			Path:  os.Args[0],
			Found: true,
		}},
		Timeout: 5 * time.Second,
		Store:   store.New(filepath.Join(root, protocol.DeckDir, "runs", "impl-run")),
	})
	if !res.ArtifactOK || res.ExitError != "" {
		t.Fatalf("unexpected result: %+v", res)
	}
	if _, err := os.Stat(filepath.Join(idea.Path, "IMPLEMENTATION.md")); err != nil {
		t.Fatalf("IMPLEMENTATION.md missing: %v", err)
	}
}

func TestRunReviewRoundWritesReviews(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	idea, err := protocol.CreateIdea(root, "Review task", []string{"rev1"})
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(idea.Path, "IMPLEMENTATION.md"), "---\nidea: "+idea.Slug+"\nstatus: implemented\n---\n\n## Summary of work\nx\n")

	results := RunReviewRound(context.Background(), Options{
		Root:  root,
		RunID: "rev-run",
		Idea:  idea,
		Round: 1,
		Agents: []agents.Discovery{{
			Spec:  agents.Spec{ID: "rev1", HeadlessArgs: []string{"-test.run=TestFakeReviewHelper", "--", "parley-fake-review"}, PromptMode: agents.PromptStdin},
			Path:  os.Args[0],
			Found: true,
		}},
		Timeout: 5 * time.Second,
		Store:   store.New(filepath.Join(root, protocol.DeckDir, "runs", "rev-run")),
	})
	if len(results) != 1 || !results[0].ArtifactOK || results[0].ExitError != "" {
		t.Fatalf("unexpected results: %+v", results)
	}
	if _, err := os.Stat(filepath.Join(idea.Path, "review", "round-01", "rev1.md")); err != nil {
		t.Fatalf("review artifact missing: %v", err)
	}
}

func TestValidateImplementationAndReviewArtifacts(t *testing.T) {
	dir := t.TempDir()
	impl := filepath.Join(dir, "IMPLEMENTATION.md")
	mustWrite(t, impl, "---\nidea: demo\nstatus: implemented\n---\n\n## Summary of work\nok\n")
	if err := ValidateImplementationArtifact(impl, "demo"); err != nil {
		t.Fatalf("valid impl rejected: %v", err)
	}
	if err := ValidateImplementationArtifact(impl, "other"); err == nil {
		t.Fatal("wrong idea slug must fail")
	}
	rev := filepath.Join(dir, "rev.md")
	mustWrite(t, rev, "---\nagent: rev1\nidea: demo\nreview-round: 2\n---\n\n## Findings\nnone\n")
	if err := ValidateReviewArtifact(rev, "rev1", "demo", 2); err != nil {
		t.Fatalf("valid review rejected: %v", err)
	}
	if err := ValidateReviewArtifact(rev, "rev1", "demo", 1); err == nil {
		t.Fatal("wrong review-round must fail")
	}
}

func TestFakeImplHelper(t *testing.T) {
	if !hasArg("parley-fake-impl") {
		return
	}
	input, _ := io.ReadAll(os.Stdin)
	out := regexp.MustCompile(`(?m)^- Create exactly this file: (.+)$`).FindStringSubmatch(string(input))
	idea := regexp.MustCompile(`(?m)^idea: (\S+)$`).FindStringSubmatch(string(input))
	if len(out) != 2 || len(idea) != 2 {
		t.Fatalf("impl prompt missing path/idea:\n%s", string(input))
	}
	body := "---\nidea: " + idea[1] + "\nstatus: implemented\nimplementer: builder\n---\n\n## Summary of work\nbuilt it\n\n## Implementation plan / checklist\n- [x] done\n\n## Deviations from FINAL.md\nnone\n\n## Notes for reviewers\nnone\n"
	if err := os.WriteFile(out[1], []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	os.Exit(0)
}

func TestFakeReviewHelper(t *testing.T) {
	if !hasArg("parley-fake-review") {
		return
	}
	input, _ := io.ReadAll(os.Stdin)
	out := regexp.MustCompile(`(?m)^- Create exactly this review file: (.+)$`).FindStringSubmatch(string(input))
	idea := regexp.MustCompile(`(?m)^idea: (\S+)$`).FindStringSubmatch(string(input))
	round := regexp.MustCompile(`(?m)^review-round: (\d+)$`).FindStringSubmatch(string(input))
	if len(out) != 2 || len(idea) != 2 || len(round) != 2 {
		t.Fatalf("review prompt missing path/idea/round:\n%s", string(input))
	}
	body := "---\nagent: rev1\nidea: " + idea[1] + "\nreview-round: " + round[1] + "\n---\n\n## Summary\nlgtm\n\n## Findings\n### [NIT] tiny\nnit.\n\n## Open questions\nnone\n"
	if err := os.WriteFile(out[1], []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	os.Exit(0)
}
