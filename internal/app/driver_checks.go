package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"parley-deck-cli/internal/driver"
	"parley-deck-cli/internal/evidence"
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
	record   evidence.CriterionRecord
}

// runChecksContract runs every criterion (sh -c, cwd = repo root) through the
// typed evidence executor, writes the markdown evidence table AND the typed
// EVIDENCE.json report, and returns (allPass, detail). Any non-zero exit fails
// closed. A structured proof of zero executed cases (empty test2json stream)
// or an all-skip run also fails closed — an exit-0 that ran nothing is not a
// pass. An evidence-write failure (markdown table or typed report) is itself a
// failure of the cycle, never a warning plus PASS.
//
// Tested-tree identity is taken BEFORE any criterion runs and re-verified
// AFTER: a command (or a concurrent writer) that changes the code tree during
// execution invalidates the whole attempt — the report still records the
// pre-execution digest (preserving the failed attempt's evidence) and the
// cycle fails closed.
func (o driverImplOps) runChecksContract(ctx context.Context, criteria []driver.CheckCriterion) (bool, string) {
	var ok bool
	var detail string
	err := evidence.WithReportWriter(ctx, o.ideaDir, func(w *evidence.ReportWriter) error {
		ok, detail = o.runChecksWithWriter(ctx, criteria, w)
		return nil
	})
	if err != nil {
		return false, "contract: evidence-write failure (publication guard): " + err.Error()
	}
	return ok, detail
}

func (o driverImplOps) runChecksWithWriter(ctx context.Context, criteria []driver.CheckCriterion, writer *evidence.ReportWriter) (bool, string) {
	// Do not let a delayed check cycle alter the table after completion.
	data, err := os.ReadFile(filepath.Join(o.ideaDir, "IMPLEMENTATION.md"))
	if err != nil {
		return false, err.Error()
	}
	_, status, err := evidence.TransitionStatusToComplete(data)
	if status == "complete" {
		return false, evidence.ErrReportFinalized.Error()
	}
	if err != nil {
		return false, err.Error()
	}

	if _, err := driver.ObserveChecksContract(o.ideaDir, driver.ChecksContractDigest(criteria)); err != nil {
		return false, err.Error()
	}
	if err := writer.PinChecksContract(driver.ChecksContractDigest(criteria)); err != nil {
		return false, err.Error()
	}
	excl, err := definedEvidenceArtifacts(o.root, o.ideaDir)
	if err != nil {
		return false, fmt.Sprintf("contract: evidence artifact scoping: %v — no completion", err)
	}
	preDigest, err := evidence.TreeDigest(o.root, excl...)
	if err != nil {
		return false, fmt.Sprintf("contract: pre-execution tree digest: %v — no completion", err)
	}
	results := make([]criterionResult, 0, len(criteria))
	allPass := true
	for _, c := range criteria {
		fmt.Fprintf(o.out, "driver: contract check %q ...\n", c.Name)
		// Executor provenance is the asserted runtime identity of the
		// implementer the driver acts for — attribution, not authentication
		// (see the evidence package trust boundary).
		rec := evidence.RunCriterion(ctx, o.root, c.Name, c.Command, o.implementer)
		res := criterionResult{
			name:     c.Name,
			ok:       rec.Status == evidence.StatusPass,
			dur:      time.Duration(rec.Command.DurationMillis) * time.Millisecond,
			exitCode: rec.Command.ExitCode,
			tail:     rec.Command.Diagnostics,
			record:   rec,
		}
		if !res.ok {
			allPass = false
		}
		results = append(results, res)
	}
	if err := o.writeValidationEvidence(results); err != nil {
		return false, fmt.Sprintf("contract: evidence-write failure (validation table): %v — no completion", err)
	}
	if err := o.writeTypedEvidence(results, preDigest, excl, writer); err != nil {
		return false, fmt.Sprintf("contract: evidence-write failure (typed report): %v — no completion", err)
	}
	// Commit the driver-authored evidence immediately so it does not leave the tree
	// dirty and trip the next fix-up cycle's gitTreeClean guard (review fix): mirrors
	// the driver committing other artifacts. Best-effort: a commit failure only warns.
	o.commitEvidence()
	if allPass {
		return true, fmt.Sprintf("contract: %d/%d criteria passed", len(results), len(results))
	}
	var failed []string
	for _, r := range results {
		if !r.ok {
			failed = append(failed, fmt.Sprintf("%s (%s, exit %d)", r.name, r.record.Status, r.exitCode))
		}
	}
	// Descriptive message so the author can fix the failing command (§14 stopping).
	return false, "contract failed: " + strings.Join(failed, ", ") + " — see IMPLEMENTATION.md ## Validation evidence"
}

// writeTypedEvidence builds and atomically persists the typed EVIDENCE.json
// report for this cycle. preDigest is the tested-tree identity taken BEFORE
// the criteria ran; the tree is digested again here and a mismatch fails the
// whole attempt (the code under test changed during execution) — but the
// report is still persisted first, so the failed attempt leaves evidence.
//
// The tested-tree digest excludes only the defined evidence artifacts
// (EVIDENCE.json and IMPLEMENTATION.md). Because excluding the whole
// IMPLEMENTATION.md would hide edits to its NON-evidence sections, the report
// additionally binds the digest of IMPLEMENTATION.md with ONLY the generated
// ## Validation evidence section removed (Report.ExtraDigests); the close
// gate recomputes and compares it.
func (o driverImplOps) writeTypedEvidence(results []criterionResult, preDigest string, excl []string, writer *evidence.ReportWriter) error {
	postDigest, err := evidence.TreeDigest(o.root, excl...)
	if err != nil {
		return fmt.Errorf("post-execution tree digest: %w", err)
	}
	restDigest, implRel, err := implementationRestDigest(o.root, o.ideaDir)
	if err != nil {
		return fmt.Errorf("non-evidence implementation digest: %w", err)
	}
	report := &evidence.Report{
		Idea:           o.ideaSlug,
		ReviewedCommit: evidence.ReviewedCommit(o.root),
		TreeSHA256:     preDigest,
		TreeDirty:      evidence.TreeDirty(o.root),
		GeneratedAt:    time.Now().UTC(),
		ExtraDigests:   map[string]string{implRel: restDigest},
	}
	for _, r := range results {
		report.Records = append(report.Records, r.record)
	}
	// Persist first: even a failed attempt (e.g. tree changed mid-run) leaves
	// its evidence artifact for the audit trail.
	if _, err := writer.Save(report); err != nil {
		return err
	}
	if postDigest != preDigest {
		return fmt.Errorf("tested tree changed during check execution (pre %s…, post %s…) — the run did not test a stable tree",
			preDigest[:12], postDigest[:12])
	}
	return nil
}

// implementationRestContent returns the exact bound bytes of the idea's
// IMPLEMENTATION.md with ONLY the driver-generated `## Validation evidence`
// section removed, plus the file's slash-separated path relative to root.
// This is the digest space of Report.ExtraDigests for this path — the bytes a
// verifier authorizes a completion transition from, and the bytes the close
// gate recomputes against.
func implementationRestContent(root, ideaDir string) (content []byte, relSlash string, err error) {
	rel, err := filepath.Rel(root, ideaDir)
	if err != nil {
		return nil, "", err
	}
	relSlash = filepath.ToSlash(filepath.Join(rel, "IMPLEMENTATION.md"))
	body, err := os.ReadFile(filepath.Join(ideaDir, "IMPLEMENTATION.md"))
	if err != nil {
		return nil, "", err
	}
	stripped := replaceSection(string(body), "## Validation evidence", "")
	return []byte(stripped), relSlash, nil
}

// implementationRestDigest returns the digest of implementationRestContent.
// Any edit to the non-evidence content changes this digest, so the broad
// tree-digest exclusion of IMPLEMENTATION.md cannot hide scope edits.
func implementationRestDigest(root, ideaDir string) (digest, relSlash string, err error) {
	content, relSlash, err := implementationRestContent(root, ideaDir)
	if err != nil {
		return "", "", err
	}
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:]), relSlash, nil
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

// commitEvidence commits the driver-authored evidence artifacts (the
// IMPLEMENTATION.md evidence table and the typed EVIDENCE.json report) so the
// tree stays clean between fix-up cycles. Best-effort and non-fatal: a non-git
// tree or a no-op commit is silently fine.
func (o driverImplOps) commitEvidence() {
	implRel := filepath.Join(o.ideaDir, "IMPLEMENTATION.md")
	evidenceRel := filepath.Join(o.ideaDir, evidence.ReportFileName)
	git := func(args ...string) error {
		cmd := exec.Command("git", append([]string{"-C", o.root}, args...)...)
		cmd.Env = append(os.Environ(), "GIT_OPTIONAL_LOCKS=0")
		return cmd.Run()
	}
	if git("rev-parse", "--is-inside-work-tree") != nil {
		return // not a git tree → nothing to commit
	}
	if err := git("add", implRel, evidenceRel); err != nil {
		fmt.Fprintf(o.out, "driver: warning — could not stage validation evidence: %v\n", err)
		return
	}
	// `git commit` is a no-op error when nothing changed; ignore that case.
	if err := git("commit", "-m", "[driver] "+o.ideaSlug+": validation evidence", "--", implRel, evidenceRel); err != nil {
		// Only warn if the files actually have staged changes (a real failure).
		if diff := exec.Command("git", "-C", o.root, "diff", "--cached", "--quiet", "--", implRel, evidenceRel).Run(); diff != nil {
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
	tbl.WriteString("| criterion | exit | duration | cases (exec/fail/skip) | result |\n")
	tbl.WriteString("|---|---|---|---|---|\n")
	for _, r := range results {
		cases := "unknown" // opaque shell output: no executed-case semantics
		if r.record.Command.ExecutedCases >= 0 {
			cases = fmt.Sprintf("%d/%d/%d", r.record.Command.ExecutedCases, max(r.record.Command.FailedCases, 0), max(r.record.Command.SkippedCases, 0))
		}
		tbl.WriteString(fmt.Sprintf("| %s | %d | %s | %s | %s |\n", r.name, r.exitCode, r.dur.Round(time.Millisecond), cases, strings.ToUpper(string(r.record.Status))))
	}
	for _, r := range results {
		if r.tail != "" {
			tbl.WriteString(fmt.Sprintf("\n<details><summary>%s output (scrubbed, truncated)</summary>\n\n```\n%s\n```\n</details>\n", r.name, r.tail))
		}
	}

	updated := replaceSection(string(body), "## Validation evidence", tbl.String())
	return writeVerificationBytes(path, []byte(updated))
}

// replaceSection replaces the `heading` section (up to the next `## ` or EOF) with
// replacement, appending it if the heading is absent.
func replaceSection(doc, heading, replacement string) string {
	start, end := -1, len(doc)
	offset := 0
	fence := byte(0)
	fenceLen := 0
	for _, raw := range strings.SplitAfter(doc, "\n") {
		line := strings.TrimSuffix(strings.TrimSuffix(raw, "\n"), "\r")
		trimmed := strings.TrimLeft(line, " ")
		indent := len(line) - len(trimmed)
		if indent <= 3 && len(trimmed) >= 3 && (trimmed[0] == '`' || trimmed[0] == '~') {
			n := 0
			for n < len(trimmed) && trimmed[n] == trimmed[0] {
				n++
			}
			if fence == 0 && n >= 3 {
				fence = trimmed[0]
				fenceLen = n
				offset += len(raw)
				continue
			}
			if fence == trimmed[0] && n >= fenceLen && strings.TrimSpace(trimmed[n:]) == "" {
				fence = 0
				offset += len(raw)
				continue
			}
		}
		if fence == 0 {
			if start < 0 && strings.TrimRight(line, " \t") == heading {
				start = offset
			} else if start >= 0 && strings.HasPrefix(line, "## ") {
				end = offset
				break
			}
		}
		offset += len(raw)
	}
	if start < 0 {
		if replacement == "" {
			return doc
		}
		if !strings.HasSuffix(doc, "\n") {
			doc += "\n"
		}
		return doc + "\n" + replacement
	}
	return doc[:start] + strings.TrimRight(replacement, "\n") + "\n" + doc[end:]
}
