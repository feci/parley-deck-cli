package app

import (
	"bytes"
	"errors"
	"fmt"
	"html"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"parley-deck-cli/internal/evidence"
)

const validationEvidenceHeading = "Validation evidence"
const validationEvidenceBindingSuffix = "#validation-evidence/v1"

// Keep exact document bytes. A normal report read is bounded too; an excluded
// carrier must not turn a symlink or unlimited stream into hidden input.
func readImplementationEvidence(ideaDir string) ([]byte, error) {
	path := filepath.Join(ideaDir, "IMPLEMENTATION.md")
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || info.Size() > 16<<20 {
		return nil, errors.New("implementation evidence must be a bounded regular file")
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	opened, err := f.Stat()
	if err != nil || !os.SameFile(info, opened) {
		return nil, errors.New("implementation evidence changed during open")
	}
	data, err := io.ReadAll(io.LimitReader(f, (16<<20)+1))
	if err != nil {
		return nil, err
	}
	if len(data) > 16<<20 {
		return nil, errors.New("implementation evidence exceeds limit")
	}
	return data, nil
}

type validationSection struct {
	start, end int
	found      bool
}

// Parse actual top-level CommonMark headings. Frontmatter, examples in fenced
// or indented code, comments, HTML blocks and nested headings are not sections.
// ATX variants and setext headings are recognized; duplicates are ambiguous.
func locateValidationEvidence(doc []byte) (validationSection, error) {
	bodyOffset := 0
	lines := bytes.SplitAfter(doc, []byte("\n"))
	if len(lines) > 0 && strings.TrimSuffix(strings.TrimSuffix(string(lines[0]), "\n"), "\r") == "---" {
		bodyOffset = len(lines[0])
		closed := false
		for _, line := range lines[1:] {
			bodyOffset += len(line)
			if strings.TrimSuffix(strings.TrimSuffix(string(line), "\n"), "\r") == "---" {
				closed = true
				break
			}
		}
		if !closed {
			return validationSection{}, errors.New("unterminated implementation frontmatter")
		}
	}
	source := doc[bodyOffset:]
	tree := goldmark.DefaultParser().Parse(text.NewReader(source))
	section := validationSection{end: len(doc)}
	for node := tree.FirstChild(); node != nil; node = node.NextSibling() {
		heading, ok := node.(*ast.Heading)
		if !ok || heading.Level > 2 {
			continue
		}
		if heading.Lines().Len() == 0 {
			return section, errors.New("unlocatable implementation heading")
		}
		offset := heading.Lines().At(0).Start
		start := bodyOffset + bytes.LastIndexByte(source[:offset], '\n') + 1
		title := strings.Join(strings.Fields(html.UnescapeString(string(heading.Text(source)))), " ")
		if section.found && section.end == len(doc) {
			section.end = start
		}
		if heading.Level == 2 && title == validationEvidenceHeading {
			if section.found {
				return section, errors.New("multiple top-level validation evidence sections")
			}
			section.start, section.end, section.found = start, len(doc), true
		}
	}
	return section, nil
}

func splitValidationEvidence(doc []byte) (rest, section []byte, err error) {
	span, err := locateValidationEvidence(doc)
	if err != nil {
		return nil, nil, err
	}
	if !span.found {
		return append([]byte(nil), doc...), nil, nil
	}
	rest = make([]byte, 0, len(doc)-(span.end-span.start))
	rest = append(rest, doc[:span.start]...)
	rest = append(rest, doc[span.end:]...)
	return rest, doc[span.start:span.end], nil
}

func validationNewline(doc []byte) string {
	if i := bytes.IndexByte(doc, '\n'); i > 0 && doc[i-1] == '\r' {
		return "\r\n"
	}
	return "\n"
}

func replaceValidationEvidence(doc, rendered []byte) ([]byte, error) {
	span, err := locateValidationEvidence(doc)
	if err != nil {
		return nil, err
	}
	var result []byte
	if span.found {
		result = append(result, doc[:span.start]...)
		result = append(result, rendered...)
		result = append(result, doc[span.end:]...)
	} else {
		result = append(result, doc...)
		result = append(result, []byte(validationNewline(doc)+validationNewline(doc))...)
		result = append(result, rendered...)
	}
	_, actual, err := splitValidationEvidence(result)
	if err != nil || !bytes.Equal(actual, rendered) {
		return nil, errors.New("generated validation evidence is not a distinct top-level section")
	}
	return result, nil
}

// JSON report persistence replaces each invalid UTF-8 byte with RuneError.
// Normalize before scrubbing to match the reloaded record, and afterward because
// byte-bounded diagnostic truncation can itself split a multibyte character.
func portableEvidenceText(value string) string {
	return string([]rune(evidence.ScrubAndTruncate(string([]rune(value)))))
}

func evidenceLabel(value string) string {
	value = portableEvidenceText(value)
	return strings.NewReplacer("|", "&#124;", "`", "&#96;", "*", "&#42;", "_", "&#95;", "[", "&#91;", "]", "&#93;", "\\", "&#92;", "\n", "&#10;", "\r", "&#13;", "\t", "&#9;").Replace(html.EscapeString(value))
}

// The table is a deterministic projection of typed ORIGINAL executions. It
// does not imply independent execution, acceptance or deployment. Timing and
// unknown counts remain explicit; diagnostics cannot close their own fence.
func renderValidationEvidence(records []evidence.CriterionRecord, newline string) ([]byte, error) {
	var out strings.Builder
	out.WriteString("## Validation evidence\n\n")
	out.WriteString("<!-- driver-populated (completion-contracts): original executions; independent acceptance is recorded separately -->\n\n")
	out.WriteString("Original command executions; independent acceptance is recorded separately.\n\n")
	out.WriteString("| criterion | exit | duration | cases (exec/fail/skip) | result |\n|---|---|---|---|---|\n")
	count := func(n int) string {
		if n < 0 {
			return "unknown"
		}
		return strconv.Itoa(n)
	}
	for _, rec := range records {
		switch rec.Status {
		case evidence.StatusPass, evidence.StatusFail, evidence.StatusSkipped, evidence.StatusNotRun:
		default:
			return nil, errors.New("unknown typed criterion status")
		}
		cases := "unknown"
		if rec.Command.ExecutedCases >= 0 {
			cases = count(rec.Command.ExecutedCases) + "/" + count(rec.Command.FailedCases) + "/" + count(rec.Command.SkippedCases)
		}
		fmt.Fprintf(&out, "| %s | %d | %d ms | %s | %s |\n", evidenceLabel(rec.Name), rec.Command.ExitCode, rec.Command.DurationMillis, cases, strings.ToUpper(string(rec.Status)))
	}
	for _, rec := range records {
		diagnostic := strings.ReplaceAll(portableEvidenceText(rec.Command.Diagnostics), "\r\n", "\n")
		if diagnostic == "" {
			continue
		}
		longest, run := 0, 0
		for _, r := range diagnostic {
			if r == '`' {
				run++
				if run > longest {
					longest = run
				}
			} else {
				run = 0
			}
		}
		fence := strings.Repeat("`", max(3, longest+1))
		fmt.Fprintf(&out, "\n<details><summary>%s output (scrubbed, truncated)</summary>\n\n%s\n%s\n%s\n\n</details>\n", evidenceLabel(rec.Name), fence, diagnostic, fence)
	}
	out.WriteString("\n")
	result := []byte(out.String())
	if newline == "\r\n" {
		result = bytes.ReplaceAll(result, []byte("\n"), []byte("\r\n"))
	}
	return result, nil
}

func verifyValidationEvidence(doc []byte, implRel string, report *evidence.Report) error {
	if report == nil {
		return errors.New("validation evidence has no typed report")
	}
	_, actual, err := splitValidationEvidence(doc)
	if err != nil {
		return err
	}
	if len(actual) == 0 {
		return errors.New("validation evidence section is missing")
	}
	bound := report.ExtraDigests[implRel+validationEvidenceBindingSuffix]
	if bound == "" {
		return errors.New("report has no validation evidence binding; rerun checks")
	}
	if bound != sha256Hex(string(actual)) {
		return errors.New("validation evidence changed after execution")
	}
	expected, err := renderValidationEvidence(report.Records, validationNewline(doc))
	if err != nil {
		return err
	}
	if !bytes.Equal(expected, actual) {
		return errors.New("validation evidence does not match the typed execution records")
	}
	return nil
}

// Initialize a missing managed section before any execution, so subsequently
// inserting a separator is not mistaken for a command's scope edit. The exact
// remaining document bytes are pinned before commands and compared afterward.
func (o driverImplOps) prepareValidationEvidence() (restDigest, implRel string, err error) {
	doc, err := readImplementationEvidence(o.ideaDir)
	if err != nil {
		return "", "", err
	}
	_, section, err := splitValidationEvidence(doc)
	if err != nil {
		return "", "", err
	}
	if section == nil {
		pending := []byte(strings.ReplaceAll("## Validation evidence\n\nPending execution.\n\n", "\n", validationNewline(doc)))
		updated, err := replaceValidationEvidence(doc, pending)
		if err != nil {
			return "", "", err
		}
		if err := writeVerificationBytes(filepath.Join(o.ideaDir, "IMPLEMENTATION.md"), updated); err != nil {
			return "", "", err
		}
	}
	return implementationRestDigest(o.root, o.ideaDir)
}
