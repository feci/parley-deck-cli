package app

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/hitl"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/store"
)

func TestAgentsListPrintsResolvedRuntime(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	bin := t.TempDir()
	writeFakeCLI(t, bin, "codex", "codex test 1.0")
	t.Setenv("PATH", bin)

	var stdout, stderr bytes.Buffer
	code := Run([]string{"agents", "list", "--dir", root}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	out := stdout.String()
	for _, want := range []string{"codex", "yes", "codex test 1.0", "workspace-write", "on-failure", "cli-default"} {
		if !strings.Contains(out, want) {
			t.Fatalf("list output missing %q:\n%s", want, out)
		}
	}
}

func TestAgentsVerifyCheapPath(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	bin := t.TempDir()
	writeFakeCLI(t, bin, "codex", "codex test 1.0")
	t.Setenv("PATH", bin)

	var stdout, stderr bytes.Buffer
	code := Run([]string{"agents", "verify", "--dir", root, "--agent", "codex"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "codex: installed version=codex test 1.0") {
		t.Fatalf("stdout=%q", stdout.String())
	}
}

func TestAgentsVerifyFullCodexReportsGitSmokeFailure(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	bin := t.TempDir()
	writeFakeCLI(t, bin, "codex", "codex test 1.0")
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))

	var stdout, stderr bytes.Buffer
	code := Run([]string{"agents", "verify", "--dir", root, "--agent", "codex", "--full", "--yes"}, &stdout, &stderr)
	if code == 0 {
		t.Fatalf("expected failure stdout=%s stderr=%s", stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "git status") {
		t.Fatalf("stderr=%q", stderr.String())
	}
}

func TestRunAnswerUpdatesQuestionAndEventLog(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	runID := "answer-run"
	runDir := filepath.Join(root, protocol.DeckDir, "runs", runID)
	question, err := hitl.New(runDir).Create(hitl.Question{
		Agent:  "codex",
		Prompt: "Which branch?",
		Risk:   hitl.RiskNormal,
	})
	if err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	code := Run([]string{"answer", "--dir", root, runID, question.ID, "main branch"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Answered "+question.ID) {
		t.Fatalf("stdout=%q", stdout.String())
	}

	questions, err := hitl.New(runDir).List()
	if err != nil {
		t.Fatal(err)
	}
	if len(questions) != 1 || questions[0].Status != hitl.StatusAnswered || questions[0].Answer != "main branch" {
		t.Fatalf("questions=%+v", questions)
	}
	events, err := store.New(runDir).Load()
	if err != nil {
		t.Fatal(err)
	}
	if got := events[len(events)-1].Type; got != "hitl.answered" {
		t.Fatalf("last event=%s, want hitl.answered", got)
	}
}

func writeFakeCLI(t *testing.T, dir, name, version string) {
	t.Helper()
	path := filepath.Join(dir, name)
	body := "#!/bin/sh\nif [ \"$1\" = \"--version\" ]; then echo '" + version + "'; exit 0; fi\ncat >/dev/null\nexit 0\n"
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
}
