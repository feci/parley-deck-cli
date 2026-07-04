package app

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/driver"
)

func TestScrubAndTruncate(t *testing.T) {
	// Every leaked shape MUST be gone from the scrubbed output.
	cases := map[string]string{
		"labeled api_key":      "api_key=sk-supersecretvalue1234567890",
		"authorization bearer": "Authorization: Bearer abcREALtoken1234567890",
		"standalone bearer":    "bearer abcREALtoken1234567890",
		"openai sk":            "using sk-abcdefghijklmnop1234 now",
		"github pat":           "token ghp_abcdefghijklmnopqrstuvwxyz0123",
		"aws access key":       "AKIAIOSFODNN7EXAMPLE",
		"password":             "password: hunter2secretvalue",
	}
	leaks := []string{"supersecretvalue", "REALtoken", "abcdefghijklmnop", "ghp_abcdefghij", "AKIAIOSFODNN7EXAMPLE", "hunter2secretvalue"}
	for name, in := range cases {
		got := scrubAndTruncate(in)
		for _, leak := range leaks {
			if strings.Contains(got, leak) {
				t.Errorf("%s: leaked %q in %q", name, leak, got)
			}
		}
	}
	long := strings.Repeat("line\n", 300)
	if n := strings.Count(scrubAndTruncate(long), "line"); n > evidenceMaxLines {
		t.Fatalf("not truncated to %d lines: %d", evidenceMaxLines, n)
	}
}

func TestReplaceSection(t *testing.T) {
	doc := "# Title\n\n## Summary\n\nx\n\n## Validation evidence\n\nold\n\n## Notes\n\nkeep\n"
	out := replaceSection(doc, "## Validation evidence", "## Validation evidence\n\nNEW\n")
	if strings.Contains(out, "old") {
		t.Fatalf("old content not replaced:\n%s", out)
	}
	if !strings.Contains(out, "NEW") || !strings.Contains(out, "## Notes") || !strings.Contains(out, "keep") {
		t.Fatalf("replaced too much or too little:\n%s", out)
	}
	// Absent heading → appended.
	out2 := replaceSection("# T\n\nbody\n", "## Validation evidence", "## Validation evidence\n\nADDED\n")
	if !strings.Contains(out2, "ADDED") {
		t.Fatal("absent section should append")
	}
}

func TestRunChecksContractWritesEvidenceAndVetoes(t *testing.T) {
	idea := t.TempDir()
	os.WriteFile(filepath.Join(idea, "IMPLEMENTATION.md"), []byte("---\nidea: x\n---\n\n## Summary of work\n\ndone\n\n## Validation evidence\n\n(pending)\n"), 0o644)
	o := driverImplOps{ideaDir: idea, root: idea, out: io.Discard}

	pass := []driver.CheckCriterion{{Name: "ok", Command: "true"}}
	if okPass, _ := o.runChecksContract(context.Background(), pass); !okPass {
		t.Fatal("passing contract should return true")
	}
	body, _ := os.ReadFile(filepath.Join(idea, "IMPLEMENTATION.md"))
	if !strings.Contains(string(body), "| ok | 0 |") {
		t.Fatalf("evidence table not written:\n%s", body)
	}

	fail := []driver.CheckCriterion{{Name: "boom", Command: "exit 3"}}
	okFail, detail := o.runChecksContract(context.Background(), fail)
	if okFail {
		t.Fatal("failing contract must veto (return false)")
	}
	if !strings.Contains(detail, "boom") {
		t.Fatalf("veto detail should name the failing criterion: %q", detail)
	}
}
