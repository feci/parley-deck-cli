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

// secretPatterns redact credential-shaped tokens from recorded output before it is
// written into the (committable) evidence section. Layered: labeled key/value pairs
// AND standalone provider token shapes, so an unlabeled `sk-…`/`ghp_…`/bearer value is
// caught even when its label was on a different line (review hardening).
var secretPatterns = []*regexp.Regexp{
	// Authorization / Bearer headers: redact the token that follows.
	regexp.MustCompile(`(?i)(authorization\s*[:=]\s*)(bearer\s+)?\S+`),
	regexp.MustCompile(`(?i)\bbearer\s+[A-Za-z0-9._\-]+`),
	// Labeled secrets: token/secret/password/api_key = value.
	regexp.MustCompile(`(?i)(token|secret|password|passwd|api[_-]?key|access[_-]?key|private[_-]?key)(\s*[:=]\s*)\S+`),
	// Standalone provider token shapes.
	regexp.MustCompile(`\bsk-[A-Za-z0-9_\-]{16,}`),
	regexp.MustCompile(`\bgh[pousr]_[A-Za-z0-9]{20,}`),
	regexp.MustCompile(`\bxox[baprs]-[A-Za-z0-9\-]{10,}`),
	regexp.MustCompile(`\bAKIA[0-9A-Z]{16}\b`),
	regexp.MustCompile(`\beyJ[A-Za-z0-9_\-]{10,}\.[A-Za-z0-9_\-]{10,}\.[A-Za-z0-9_\-]{5,}`), // JWT
}

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
	} else {
		// Commit the driver-authored evidence immediately so it does not leave the tree
		// dirty and trip the next fix-up cycle's gitTreeClean guard (review fix): mirrors
		// the driver committing other artifacts. Best-effort: a commit failure only warns.
		o.commitEvidence()
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
	// Labeled patterns keep their key (group 1/2) and redact the value; standalone
	// token shapes are replaced whole.
	s = secretPatterns[0].ReplaceAllString(s, "$1«redacted»")
	s = secretPatterns[2].ReplaceAllString(s, "$1$2«redacted»")
	for i, re := range secretPatterns {
		if i == 0 || i == 2 {
			continue
		}
		s = re.ReplaceAllString(s, "«redacted»")
	}
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

// commitEvidence commits the driver-authored IMPLEMENTATION.md evidence write so the
// tree stays clean between fix-up cycles. Best-effort and non-fatal: a non-git tree or
// a no-op commit is silently fine.
func (o driverImplOps) commitEvidence() {
	rel := filepath.Join(o.ideaDir, "IMPLEMENTATION.md")
	git := func(args ...string) error {
		cmd := exec.Command("git", append([]string{"-C", o.root}, args...)...)
		cmd.Env = append(os.Environ(), "GIT_OPTIONAL_LOCKS=0")
		return cmd.Run()
	}
	if git("rev-parse", "--is-inside-work-tree") != nil {
		return // not a git tree → nothing to commit
	}
	if err := git("add", rel); err != nil {
		fmt.Fprintf(o.out, "driver: warning — could not stage validation evidence: %v\n", err)
		return
	}
	// `git commit` is a no-op error when nothing changed; ignore that case.
	if err := git("commit", "-m", "[driver] "+o.ideaSlug+": validation evidence", "--", rel); err != nil {
		// Only warn if the file actually has staged changes (a real failure).
		if diff := exec.Command("git", "-C", o.root, "diff", "--cached", "--quiet", "--", rel).Run(); diff != nil {
			fmt.Fprintf(o.out, "driver: warning — could not commit validation evidence: %v\n", err)
		}
	}
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
