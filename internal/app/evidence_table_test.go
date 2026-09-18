package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	gmhtml "github.com/yuin/goldmark/renderer/html"
	"os"
	"parley-deck-cli/internal/driver"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"

	"parley-deck-cli/internal/evidence"
)

func TestEvidenceTableTamperingCannotClose(t *testing.T) {
	for _, kind := range []string{"claims", "counts", "deleted"} {
		t.Run(kind, func(t *testing.T) {
			root, dir := gateScratchRepo(t, twoCriterionContract())
			op := gateOps(root, dir)
			report := attestContract(t, root, dir, op, "reviewer")
			if err := evidence.Save(dir, report); err != nil {
				t.Fatal(err)
			}
			if gate := op.EvidenceCloseGate("reviewer"); !gate.Allowed {
				t.Fatal(gate.Reasons)
			}
			path := filepath.Join(dir, "IMPLEMENTATION.md")
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			changed := string(raw)
			switch kind {
			case "claims":
				replacement, e := replaceValidationEvidence(raw, []byte("## Validation evidence\n\nEvery requirement independently passed.\n\n"))
				if e != nil {
					t.Fatal(e)
				}
				changed = string(replacement)
			case "counts":
				changed = strings.ReplaceAll(changed, "| 1/0/0 |", "| 999/0/0 |")
			case "deleted":
				rest, _, e := splitValidationEvidence(raw)
				if e != nil {
					t.Fatal(e)
				}
				changed = string(rest)
			}
			if changed == string(raw) {
				t.Fatal("mutation did not land")
			}
			if err := os.WriteFile(path, []byte(changed), 0644); err != nil {
				t.Fatal(err)
			}
			if gate := op.EvidenceCloseGate("reviewer"); gate.Allowed {
				t.Fatal("changed or missing human evidence still permits closure")
			}
		})
	}
}

func TestEvidenceTableRehashedForgeryAndLegacyReportRefuse(t *testing.T) {
	for _, kind := range []string{"rehashed-forgery", "legacy-binding", "exact-section-deletion"} {
		t.Run(kind, func(t *testing.T) {
			root, dir := gateScratchRepo(t, twoCriterionContract())
			op := gateOps(root, dir)
			report := attestContract(t, root, dir, op, "reviewer")
			path := filepath.Join(dir, "IMPLEMENTATION.md")
			doc, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			rel := "parley-deck/ideas/idea-x/IMPLEMENTATION.md"
			switch kind {
			case "rehashed-forgery":
				doc = bytes.ReplaceAll(doc, []byte("| 1/0/0 |"), []byte("| 999/0/0 |"))
				_, section, e := splitValidationEvidence(doc)
				if e != nil {
					t.Fatal(e)
				}
				report.ExtraDigests[rel+validationEvidenceBindingSuffix] = sha256Hex(string(section))
			case "legacy-binding":
				delete(report.ExtraDigests, rel+validationEvidenceBindingSuffix)
			case "exact-section-deletion":
				doc, _, err = splitValidationEvidence(doc)
				if err != nil {
					t.Fatal(err)
				}
				if sha256Hex(string(doc)) != report.ExtraDigests[rel] {
					t.Fatal("deletion changed supposedly bound remainder")
				}
			}
			if err := os.WriteFile(path, doc, 0644); err != nil {
				t.Fatal(err)
			}
			if err := evidence.Save(dir, report); err != nil {
				t.Fatal(err)
			}
			gate := op.EvidenceCloseGate("reviewer")
			if gate.Allowed || !strings.Contains(strings.Join(gate.Reasons, ";"), "validation evidence") {
				t.Fatalf("wrong projection refusal: %+v", gate)
			}
		})
	}
}

func TestEvidenceTableCommonMarkScope(t *testing.T) {
	prefix := "---\nstatus: implemented\nexample: |\n  ## Validation evidence\n---\n\n# Implementation\n\n"
	examples := []string{
		"```md\n## Validation evidence\nexample\n```\n\n",
		"~~~md\n## Validation evidence\nexample\n~~~\n\n",
		"    ## Validation evidence\n    example\n\n",
		"<!--\n## Validation evidence\nexample\n-->\n\n",
		"<pre>\n## Validation evidence\nexample\n</pre>\n\n",
		"> ## Validation evidence\n> example\n\n",
		"- Nested item\n\n  ## Validation evidence\n\n  example\n\n",
		"### Validation evidence\nexample\n\n",
	}
	for i, example := range examples {
		for _, heading := range []string{"## Validation evidence", "  ## Validation evidence ###", "Validation evidence\n-------------------", "## **Validation evidence**", "## Validation &#101;vidence"} {
			t.Run(fmt.Sprintf("%d-%s", i, heading), func(t *testing.T) {
				before := prefix + example + "End of example.\n\n"
				section := heading + "\n\nACTUAL\n\n"
				after := "# Bound follow-up\n\nPreserve exactly.\n"
				rest, actual, err := splitValidationEvidence([]byte(before + section + after))
				if err != nil || string(actual) != section || string(rest) != before+after {
					t.Fatalf("wrong CommonMark span: actual=%q rest=%q error=%v", actual, rest, err)
				}
			})
		}
	}
	for _, doc := range []string{
		"## Validation evidence\nfirst\n\n## Validation evidence\nsecond\n",
		"## Validation evidence\nfirst\n\nValidation evidence\n---\nsecond\n",
		"## Validation evidence\nfirst\n\n## **Validation evidence**\nsecond\n",
	} {
		if _, _, err := splitValidationEvidence([]byte(doc)); err == nil {
			t.Fatal("ambiguous managed sections accepted")
		}
	}
}

func TestEvidenceTableRenderingCannotInjectClaimsOrMarkup(t *testing.T) {
	name := "unit | [PASS](https://example.invalid)\n## Validation evidence <img src=x onerror=alert(1)>"
	diagnostic := "```\n## Forged scope\n<script>alert(1)</script>\n````\nAuthorization: Bearer tokenvalue"
	records := []evidence.CriterionRecord{{Name: name, Status: evidence.StatusSkipped, Command: evidence.CommandEvidence{ExitCode: 0, DurationMillis: 123, ExecutedCases: 0, FailedCases: -1, SkippedCases: 1, Diagnostics: diagnostic}}}
	rendered, err := renderValidationEvidence(records, "\n")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(rendered, []byte("| 0/unknown/1 | SKIPPED |")) || !bytes.Contains(rendered, []byte("123 ms")) || bytes.Contains(rendered, []byte("tokenvalue")) {
		t.Fatalf("misrepresented actual results: %s", rendered)
	}
	doc := append([]byte("---\nstatus: implemented\n---\n\n"), rendered...)
	doc = append(doc, []byte("## Bound scope\nunchanged\n")...)
	rest, section, err := splitValidationEvidence(doc)
	if err != nil || !bytes.Equal(section, rendered) || !bytes.Contains(rest, []byte("## Bound scope\nunchanged")) {
		t.Fatalf("rendered data escaped its managed section: %v", err)
	}
	var html bytes.Buffer
	md := goldmark.New(goldmark.WithExtensions(extension.GFM), goldmark.WithRendererOptions(gmhtml.WithUnsafe()))
	if err := md.Convert(rendered, &html); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(html.String(), "<script>") || strings.Contains(html.String(), "<img") || strings.Contains(html.String(), "<h2>Forged scope") {
		t.Fatalf("active injected markup: %s", html.String())
	}
	if !strings.Contains(html.String(), "&lt;script&gt;") {
		t.Fatal("diagnostic code content was lost")
	}
}

func TestEvidenceTableCRLFAndMissingSectionPreparation(t *testing.T) {
	for _, newline := range []string{"\n", "\r\n"} {
		for _, missing := range []bool{false, true} {
			t.Run(fmt.Sprintf("%q-%t", newline, missing), func(t *testing.T) {
				root, dir := gateScratchRepo(t, twoCriterionContract())
				path := filepath.Join(dir, "IMPLEMENTATION.md")
				doc, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				if missing {
					doc, _, err = splitValidationEvidence(doc)
					if err != nil {
						t.Fatal(err)
					}
				}
				doc = bytes.ReplaceAll(doc, []byte("\n"), []byte(newline))
				if err := os.WriteFile(path, doc, 0644); err != nil {
					t.Fatal(err)
				}
				op := gateOps(root, dir)
				report := attestContract(t, root, dir, op, "reviewer")
				authorizeCompletion(t, root, dir, report, "reviewer")
				flipStatusComplete(t, dir)
				if gate := op.EvidenceCloseGate("reviewer"); !gate.Allowed {
					t.Fatal(gate.Reasons)
				}
				completed, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				if err := verifyValidationEvidence(completed, "parley-deck/ideas/idea-x/IMPLEMENTATION.md", report); err != nil {
					t.Fatal(err)
				}
			})
		}
	}
}

func TestEvidenceTablePreservesPreExecutionScope(t *testing.T) {
	command := passJSONCmd + "; printf '\\n## Added during execution\\nUnverified scope\\n' >> parley-deck/ideas/idea-x/IMPLEMENTATION.md"
	criteria := []driver.CheckCriterion{{Name: "unit", Command: command}}
	root, dir := scratchContract(t, criteria)
	op := gateOps(root, dir)
	ok, detail := op.runChecksContract(context.Background(), criteria)
	if ok || !strings.Contains(detail, "implementation changed during check execution") {
		t.Fatalf("scope drift was accepted: %t %s", ok, detail)
	}
	report, err := evidence.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Records) != 1 || report.Records[0].Status != evidence.StatusPass {
		t.Fatal("failed attempt lost actual command result")
	}
	rest, rel, err := implementationRestDigest(root, dir)
	if err != nil {
		t.Fatal(err)
	}
	if report.ExtraDigests[rel] == rest {
		t.Fatal("scope drift was rebound to the post-execution document")
	}
}

func TestEvidenceTableProjectionSurvivesJSONUnicodeNormalization(t *testing.T) {
	records := []evidence.CriterionRecord{{Name: "unit", Status: evidence.StatusPass, Command: evidence.CommandEvidence{ExecutedCases: 1, FailedCases: 0, SkippedCases: 0, Diagnostics: evidence.ScrubAndTruncate(strings.Repeat("😀", 1600) + "x")}}}
	original, err := renderValidationEvidence(records, "\n")
	if err != nil {
		t.Fatal(err)
	}
	raw, err := json.Marshal(records)
	if err != nil {
		t.Fatal(err)
	}
	var persisted []evidence.CriterionRecord
	if err := json.Unmarshal(raw, &persisted); err != nil {
		t.Fatal(err)
	}
	reloaded, err := renderValidationEvidence(persisted, "\n")
	if err != nil {
		t.Fatal(err)
	}
	if !utf8.Valid(original) || !bytes.Equal(original, reloaded) {
		t.Fatal("Unicode truncation changes the projection after report persistence")
	}
}
