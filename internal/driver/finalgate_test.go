package driver

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fullFinal is a FINAL.md that satisfies the protocol template.
func fullFinal(slug string) string {
	var b strings.Builder
	b.WriteString("---\nidea: " + slug + "\nstatus: final\nauthor: a\n---\n\n")
	for _, s := range requiredFinalSections {
		b.WriteString(s + "\n\nSubstantive content for this section, stated plainly.\nA second line of real detail.\nA third line of real detail.\n\n")
	}
	return b.String()
}

func writeFinal(t *testing.T, slug, body string) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), slug)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "FINAL.md")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// codex-1/F22: one generic heading plus three padded lines passed as a complete specification.
func TestPaddedScaffoldIsRejectedForMissingSections(t *testing.T) {
	body := "---\nidea: demo\nstatus: final\n---\n\n## Final plan / specification\n\nyes\nok\n" +
		strings.Repeat("padding to cross the byte floor. ", 12) + "\n"
	reason := finalScaffoldReason(writeFinal(t, "demo", body))
	if reason == "" {
		t.Fatal("a one-heading padded scaffold was accepted as a final specification")
	}
	if !strings.Contains(reason, "missing required section") {
		t.Fatalf("reason should name the missing sections, got %q", reason)
	}
	for _, want := range []string{"Purpose", "Observable acceptance criteria", "Known risks"} {
		if !strings.Contains(reason, want) {
			t.Errorf("reason does not name %q: %s", want, reason)
		}
	}
}

// codex-1/F22 also declared the wrong idea slug and the gate never looked.
func TestFinalDeclaringAnotherIdeaIsRejected(t *testing.T) {
	reason := finalScaffoldReason(writeFinal(t, "real-idea", fullFinal("some-other-idea")))
	if !strings.Contains(reason, "closes idea") {
		t.Fatalf("slug mismatch not reported, got %q", reason)
	}
}

func TestFinalWithNoSlugIsRejected(t *testing.T) {
	body := strings.Replace(fullFinal("demo"), "idea: demo\n", "", 1)
	if reason := finalScaffoldReason(writeFinal(t, "demo", body)); !strings.Contains(reason, "no idea slug") {
		t.Fatalf("missing slug not reported, got %q", reason)
	}
}

// A complete FINAL must still pass; the gate must not reject real work.
func TestCompleteFinalIsAccepted(t *testing.T) {
	if reason := finalScaffoldReason(writeFinal(t, "demo", fullFinal("demo"))); reason != "" {
		t.Fatalf("a complete FINAL was rejected: %s", reason)
	}
}

// The protocol allows N/A content for trivial ideas; only the heading is mandatory.
func TestSectionsMayBeNA(t *testing.T) {
	var b strings.Builder
	b.WriteString("---\nidea: demo\nstatus: final\n---\n\n")
	b.WriteString(requiredFinalSections[0] + "\n\nThe real specification lives here.\nSecond line.\nThird line.\n\n")
	for _, s := range requiredFinalSections[1:] {
		b.WriteString(s + "\n\nN/A\n\n")
	}
	if reason := finalScaffoldReason(writeFinal(t, "demo", b.String())); reason != "" {
		t.Fatalf("N/A sections were rejected: %s", reason)
	}
}
