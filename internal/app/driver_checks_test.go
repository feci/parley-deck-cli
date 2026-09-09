package app

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/driver"
	"parley-deck-cli/internal/evidence"
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
	root, idea := gateScratchRepo(t, twoCriterionContract())
	o := driverImplOps{ideaDir: idea, root: root, ideaSlug: "x", implementer: "kimi-1", out: io.Discard}

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

const passJSONLine = `printf '%s\n' '{"Action":"run","Test":"TestA"}' '{"Action":"pass","Test":"TestA"}' '{"Action":"pass","Package":"x"}'`

// Positive: structured test2json output yields typed evidence with real
// executed-case counts, persisted to EVIDENCE.json and the markdown table.
func TestRunChecksContractCountsExecutedCases(t *testing.T) {
	root, idea := gateScratchRepo(t, twoCriterionContract())
	o := driverImplOps{ideaDir: idea, root: root, ideaSlug: "x", implementer: "kimi-1", out: io.Discard}

	ok, detail := o.runChecksContract(context.Background(), []driver.CheckCriterion{{Name: "unit", Command: passJSONLine}})
	if !ok {
		t.Fatalf("structured pass should return true: %s", detail)
	}
	report, err := evidence.Load(idea)
	if err != nil {
		t.Fatalf("typed report not persisted: %v", err)
	}
	if len(report.Records) != 1 || report.Records[0].Command.ExecutedCases != 1 || report.Records[0].Command.Format != evidence.FormatGoTestJSON {
		t.Fatalf("expected 1 executed case via gotest-json, got %+v", report.Records)
	}
	if report.Records[0].Provenance.Executor != "kimi-1" {
		t.Fatalf("executor provenance missing: %+v", report.Records[0].Provenance)
	}
	body, _ := os.ReadFile(filepath.Join(idea, "IMPLEMENTATION.md"))
	if !strings.Contains(string(body), "| 1/0/0 |") {
		t.Fatalf("markdown table lacks case counts:\n%s", body)
	}
}

// Adversarial: exit 0 with structured proof of ZERO executed cases vetoes the
// cycle — an empty test run is not a pass.
func TestRunChecksContractZeroExecutionVetoes(t *testing.T) {
	root, idea := gateScratchRepo(t, twoCriterionContract())
	o := driverImplOps{ideaDir: idea, root: root, ideaSlug: "x", implementer: "kimi-1", out: io.Discard}

	zero := `printf '%s\n' '{"Action":"start","Package":"x"}' '{"Action":"pass","Package":"x"}'`
	ok, detail := o.runChecksContract(context.Background(), []driver.CheckCriterion{{Name: "unit", Command: zero}})
	if ok {
		t.Fatal("zero-execution structured output must veto")
	}
	if !strings.Contains(detail, "not-run") {
		t.Fatalf("veto detail should type the zero-execution verdict: %q", detail)
	}
	report, err := evidence.Load(idea)
	if err != nil || report.Records[0].Command.ExecutedCases != 0 {
		t.Fatalf("report must record 0 executed cases: %v %+v", err, report)
	}
}

// Adversarial: opaque shell output containing the word PASS is exit-code
// evidence only — never an executed-case count.
func TestRunChecksContractUnknownOutputNotCertified(t *testing.T) {
	root, idea := gateScratchRepo(t, twoCriterionContract())
	o := driverImplOps{ideaDir: idea, root: root, ideaSlug: "x", implementer: "kimi-1", out: io.Discard}

	ok, _ := o.runChecksContract(context.Background(), []driver.CheckCriterion{{Name: "unit", Command: "echo 'PASS all green'"}})
	if !ok {
		t.Fatal("shell exit 0 still passes the per-cycle gate")
	}
	report, err := evidence.Load(idea)
	if err != nil {
		t.Fatal(err)
	}
	c := report.Records[0].Command
	if c.Format != evidence.FormatShell || c.ExecutedCases != -1 {
		t.Fatalf("opaque output must stay shell/unknown: %+v", c)
	}
	body, _ := os.ReadFile(filepath.Join(idea, "IMPLEMENTATION.md"))
	if !strings.Contains(string(body), "| unknown |") {
		t.Fatalf("markdown table must mark unknown case counts:\n%s", body)
	}
}

// Adversarial: an evidence-write failure vetoes the cycle outright — never a
// warning plus PASS.
func TestRunChecksContractEvidenceWriteFailureVetoes(t *testing.T) {
	idea := t.TempDir()
	os.WriteFile(filepath.Join(idea, "IMPLEMENTATION.md"), []byte("---\nidea: x\n---\n\n## Validation evidence\n\n(pending)\n"), 0o644)
	os.Chmod(idea, 0o555)
	t.Cleanup(func() { os.Chmod(idea, 0o755) })
	o := driverImplOps{ideaDir: idea, root: idea, ideaSlug: "x", implementer: "kimi-1", out: io.Discard}

	ok, detail := o.runChecksContract(context.Background(), []driver.CheckCriterion{{Name: "unit", Command: "true"}})
	if ok {
		t.Fatal("evidence-write failure must veto, not warn-and-pass")
	}
	if !strings.Contains(detail, "evidence-write failure") {
		t.Fatalf("veto detail should name the evidence-write failure: %q", detail)
	}
}
