package runner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/protocol"
)

func TestSanitizeForContextRemovesSupportedReasoningFences(t *testing.T) {
	input := strings.Join([]string{
		"keep before",
		"<think>remove this</think>",
		"keep middle",
		"<thought>remove this too</thought>",
		"<thinking>remove nested words</thinking>",
		"keep after",
	}, "\n")

	got := SanitizeForContext(input)
	for _, unwanted := range []string{"remove this", "remove this too", "remove nested words", "<think>", "<thought>", "<thinking>"} {
		if strings.Contains(got, unwanted) {
			t.Fatalf("sanitized output still contains %q:\n%s", unwanted, got)
		}
	}
	for _, want := range []string{"keep before", "keep middle", "keep after"} {
		if !strings.Contains(got, want) {
			t.Fatalf("sanitized output missing %q:\n%s", want, got)
		}
	}
}

func TestSanitizeForContextDropsMalformedOpenFence(t *testing.T) {
	got := SanitizeForContext("visible\n<think>unfinished hidden text")
	if strings.Contains(got, "unfinished hidden text") {
		t.Fatalf("malformed open fence was not removed: %q", got)
	}
	if !strings.Contains(got, "visible") {
		t.Fatalf("visible text was removed: %q", got)
	}
}

func TestBuildRoundIndexIsDeterministicAndExtractsH2Only(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	idea, err := protocol.CreateIdea(root, "Index test task", []string{"beta", "alpha"})
	if err != nil {
		t.Fatal(err)
	}
	alpha := filepath.Join(idea.Path, "round-01", "alpha.md")
	alphaBody := `---
agent: alpha
idea: index-test-task
round: 1
---

<think>private reasoning</think>

## Summary
Alpha summary paragraph.

### Detail
This heading must not become an index section.

## Risks
Alpha risk paragraph.
`
	if err := os.WriteFile(alpha, []byte(alphaBody), 0o644); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(alpha)
	if err != nil {
		t.Fatal(err)
	}

	results := []Result{
		{AgentID: "beta", OutputPath: filepath.Join(idea.Path, "round-01", "beta.md"), ExitError: "exit status 7"},
		{AgentID: "alpha", OutputPath: alpha, ArtifactOK: true},
	}
	first := BuildRoundIndex(idea, "round-01", results)
	second := BuildRoundIndex(idea, "round-01", results)
	if first != second {
		t.Fatalf("round index output is not deterministic\nfirst:\n%s\nsecond:\n%s", first, second)
	}
	for _, want := range []string{
		"artifact: round-index",
		"| alpha | ok |",
		"| beta | failed |",
		"## alpha",
		"Summary: Alpha summary paragraph.",
		"Risks: Alpha risk paragraph.",
		"Approx tokens heuristic",
	} {
		if !strings.Contains(first, want) {
			t.Fatalf("index missing %q:\n%s", want, first)
		}
	}
	for _, unwanted := range []string{"private reasoning", "Detail:"} {
		if strings.Contains(first, unwanted) {
			t.Fatalf("index contains unwanted %q:\n%s", unwanted, first)
		}
	}
	after, err := os.ReadFile(alpha)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatalf("source artifact was modified")
	}
}

func TestBuildRoundIndexIncludesSkippedWithoutRecognizedSections(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	idea, err := protocol.CreateIdea(root, "Skipped index task", []string{"fake"})
	if err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(idea.Path, "round-01", "fake.md")
	if err := os.WriteFile(output, []byte("already done\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	index := BuildRoundIndex(idea, "round-01", []Result{{
		AgentID:    "fake",
		OutputPath: output,
		Skipped:    true,
		SkipReason: "artifact already exists",
	}})
	for _, want := range []string{"| fake | skipped |", "Sections: no recognized H2 sections", "artifact already exists"} {
		if !strings.Contains(index, want) {
			t.Fatalf("index missing %q:\n%s", want, index)
		}
	}
}
