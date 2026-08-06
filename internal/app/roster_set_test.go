package app

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeDeckTOML(t *testing.T, body string) string {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, "parley-deck")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "agents.toml"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// Preview is the default. Without --yes nothing may be written, so an unattended
// invocation cannot rewrite a roster by accident.
func TestRosterSetPreviewWritesNothing(t *testing.T) {
	body := "[roster.kimi-1]\nadapter = \"kimi\"\n"
	root := writeDeckTOML(t, body)
	var out, errb bytes.Buffer
	if code := runRoster([]string{"set", "kimi-1", "--dir", root, "--scope", "deck", "--effort", "max"}, &out, &errb); code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, errb.String())
	}
	if !strings.Contains(out.String(), "Nothing was written") {
		t.Fatalf("preview must say nothing was written:\n%s", out.String())
	}
	got, _ := os.ReadFile(filepath.Join(root, "parley-deck", "agents.toml"))
	if string(got) != body {
		t.Fatalf("file changed without --yes:\n%s", got)
	}
}

// The config files carry comments recording WHY a value is pinned. A marshal/unmarshal
// round-trip would delete them all, so the writer is line-based and this guards it.
func TestRosterSetPreservesCommentsAndOtherBlocks(t *testing.T) {
	body := `# deck defaults

[agents.kimi]
model = "kimi-code/k3"   # pinned after a probe

[roster.kimi-1]
adapter = "kimi"
# why: kimi reads its own config
`
	root := writeDeckTOML(t, body)
	var out, errb bytes.Buffer
	if code := runRoster([]string{"set", "kimi-1", "--dir", root, "--scope", "deck", "--effort", "max", "--yes"}, &out, &errb); code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, errb.String())
	}
	got, _ := os.ReadFile(filepath.Join(root, "parley-deck", "agents.toml"))
	for _, want := range []string{"# deck defaults", "# pinned after a probe", "# why: kimi reads its own config", "[agents.kimi]", "effort = \"max\""} {
		if !strings.Contains(string(got), want) {
			t.Errorf("result lost %q:\n%s", want, got)
		}
	}
}

// Re-applying the same value is a no-op, so `set` is safe to run from a checklist.
func TestRosterSetIsIdempotent(t *testing.T) {
	root := writeDeckTOML(t, "[roster.kimi-1]\nadapter = \"kimi\"\neffort = \"max\"\n")
	var out, errb bytes.Buffer
	if code := runRoster([]string{"set", "kimi-1", "--dir", root, "--scope", "deck", "--effort", "max", "--yes"}, &out, &errb); code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, errb.String())
	}
	if !strings.Contains(out.String(), "already matches") {
		t.Fatalf("expected a no-op, got:\n%s", out.String())
	}
}

// --state inactive must MARK, never delete: a past idea's participant list has to stay
// interpretable, so roster history is retained permanently.
func TestRosterSetInactiveMarksRatherThanDeletes(t *testing.T) {
	root := writeDeckTOML(t, "[roster.antigravity-1]\nadapter = \"agy\"\n")
	var out, errb bytes.Buffer
	if code := runRoster([]string{"set", "antigravity-1", "--dir", root, "--scope", "deck", "--state", "inactive", "--yes"}, &out, &errb); code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, errb.String())
	}
	got, _ := os.ReadFile(filepath.Join(root, "parley-deck", "agents.toml"))
	if !strings.Contains(string(got), "[roster.antigravity-1]") {
		t.Fatalf("the entry was removed instead of marked:\n%s", got)
	}
	if !strings.Contains(string(got), "active = false") {
		t.Fatalf("expected active = false:\n%s", got)
	}
}

// `--scope deck` must target the COMMITTED file. A roster change written to the
// gitignored agents.local.toml is invisible to the repository, which is precisely how a
// deck silently diverges from its own history.
func TestRosterSetDeckScopeTargetsCommittedFile(t *testing.T) {
	root := t.TempDir()
	target, err := rosterScopeFile(root, "deck")
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(target) != "agents.toml" {
		t.Fatalf("deck scope writes %q, want agents.toml", filepath.Base(target))
	}
	if strings.Contains(target, "agents.local.toml") {
		t.Fatalf("deck scope must never write the gitignored local file: %s", target)
	}
}
