package app

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"parley-deck-cli/internal/driver"
)

// driver_checks.go executes the list-form `checks:` completion contract and records the
// evidence table into IMPLEMENTATION.md's `## Validation evidence` section
// (completion-contracts-evidence-ledger). Output is secret-scrubbed and truncated.

const (
	evidenceMaxLines = 100
	evidenceMaxBytes = 4096
)

// secretPattern redacts obvious credential-shaped tokens from recorded output.
var secretPattern = regexp.MustCompile(`(?i)(token|secret|password|api[_-]?key|bearer|authorization)[=:\s]+\S+`)

type criterionResult struct {
	name     string
	exitCode int
	ok       bool
	dur      time.Duration
	tail     string
}

// runChecksContract runs every criterion (sh -c, cwd = repo root), writes the evidence
// table, and returns (allPass, detail). Any non-zero exit fails closed.
func (o driverImplOps) runChecksContract(ctx context.Context, criteria []driver.CheckCriterion) (bool, string) {
	results := make([]criterionResult, 0, len(criteria))
	allPass := true
	for _, c := range criteria {
		fmt.Fprintf(o.out, "driver: contract check %q ...\n", c.Name)
		start := time.Now()
		cmd := exec.CommandContext(ctx, "sh", "-c", c.Command)
		cmd.Dir = o.root
		var buf bytes.Buffer
		cmd.Stdout = &buf
		cmd.Stderr = &buf
		err := cmd.Run()
		res := criterionResult{
			name:     c.Name,
			ok:       err == nil,
			dur:      time.Since(start),
			exitCode: exitCodeOf(err),
			tail:     scrubAndTruncate(buf.String()),
		}
		if !res.ok {
			allPass = false
		}
		results = append(results, res)
	}
	if err := o.writeValidationEvidence(results); err != nil {
		fmt.Fprintf(o.out, "driver: warning — could not write validation evidence: %v\n", err)
	}
	if allPass {
		return true, fmt.Sprintf("contract: %d/%d criteria passed", len(results), len(results))
	}
	var failed []string
	for _, r := range results {
		if !r.ok {
			failed = append(failed, fmt.Sprintf("%s (exit %d)", r.name, r.exitCode))
		}
	}
	// Descriptive message so the author can fix the failing command (§14 stopping).
	return false, "contract failed: " + strings.Join(failed, ", ") + " — see IMPLEMENTATION.md ## Validation evidence"
}

func exitCodeOf(err error) int {
	if err == nil {
		return 0
	}
	if ee, ok := err.(*exec.ExitError); ok {
		return ee.ExitCode()
	}
	return -1
}

// scrubAndTruncate redacts credential-shaped tokens and caps the output to a bounded
// tail (last evidenceMaxLines lines, then evidenceMaxBytes bytes).
func scrubAndTruncate(s string) string {
	s = secretPattern.ReplaceAllString(s, "$1=«redacted»")
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	if len(lines) > evidenceMaxLines {
		lines = lines[len(lines)-evidenceMaxLines:]
	}
	out := strings.Join(lines, "\n")
	if len(out) > evidenceMaxBytes {
		out = "…" + out[len(out)-evidenceMaxBytes:]
	}
	return out
}

// writeValidationEvidence overwrites the `## Validation evidence` section of
// IMPLEMENTATION.md with the latest per-criterion table (git history keeps prior cycles).
func (o driverImplOps) writeValidationEvidence(results []criterionResult) error {
	path := filepath.Join(o.ideaDir, "IMPLEMENTATION.md")
	body, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var tbl strings.Builder
	tbl.WriteString("## Validation evidence\n\n")
	tbl.WriteString("<!-- driver-populated (completion-contracts): overwritten each cycle; git history keeps prior runs -->\n\n")
	tbl.WriteString("| criterion | exit | duration | result |\n")
	tbl.WriteString("|---|---|---|---|\n")
	for _, r := range results {
		verdict := "PASS"
		if !r.ok {
			verdict = "FAIL"
		}
		tbl.WriteString(fmt.Sprintf("| %s | %d | %s | %s |\n", r.name, r.exitCode, r.dur.Round(time.Millisecond), verdict))
	}
	for _, r := range results {
		if r.tail != "" {
			tbl.WriteString(fmt.Sprintf("\n<details><summary>%s output (scrubbed, truncated)</summary>\n\n```\n%s\n```\n</details>\n", r.name, r.tail))
		}
	}

	updated := replaceSection(string(body), "## Validation evidence", tbl.String())
	return os.WriteFile(path, []byte(updated), 0o644)
}

// replaceSection replaces the `heading` section (up to the next `## ` or EOF) with
// replacement, appending it if the heading is absent.
func replaceSection(doc, heading, replacement string) string {
	idx := strings.Index(doc, heading)
	if idx < 0 {
		if !strings.HasSuffix(doc, "\n") {
			doc += "\n"
		}
		return doc + "\n" + replacement
	}
	rest := doc[idx+len(heading):]
	next := strings.Index(rest, "\n## ")
	if next < 0 {
		return doc[:idx] + strings.TrimRight(replacement, "\n") + "\n"
	}
	return doc[:idx] + strings.TrimRight(replacement, "\n") + rest[next:]
}
