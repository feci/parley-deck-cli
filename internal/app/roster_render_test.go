package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func renderFixture(t *testing.T, tomlBody, coopBody string) string {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, "parley-deck")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "agents.toml"), []byte(tomlBody), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "COOPERATION.md"), []byte(coopBody), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PARLEY_HOME", t.TempDir())
	return root
}

const coopSkeleton = `# Protocol

## 2. Active agents (roster)

Prose explaining that this table is generated.

| Agent ID | Workspace dir | Role | State |
| -------- | ------------- | ---- | ----- |
| ` + "`claude-1`" + ` | ../claude/ | facilitator | active |

## 3. Next section

Body.
`

// THE GATE (FINAL.md G4). A non-idempotent generator recreates the drift it exists to
// end: every render would be a diff, and nobody could tell a real change from noise.
func TestRosterRenderIsIdempotent(t *testing.T) {
	root := renderFixture(t, "[roster.claude-1]\nadapter = \"claude\"\n[roster.kimi-1]\nadapter = \"kimi\"\n", coopSkeleton)
	path := filepath.Join(root, "parley-deck", "COOPERATION.md")

	var out, errb strings.Builder
	if code := rosterRender(root, false, true, false, &out, &errb); code != 0 {
		t.Fatalf("first render exit=%d stderr=%s", code, errb.String())
	}
	first, _ := os.ReadFile(path)

	out.Reset()
	errb.Reset()
	if code := rosterRender(root, false, true, false, &out, &errb); code != 0 {
		t.Fatalf("second render exit=%d stderr=%s", code, errb.String())
	}
	second, _ := os.ReadFile(path)

	if string(first) != string(second) {
		t.Fatalf("render is not idempotent:\n--- first ---\n%s\n--- second ---\n%s", first, second)
	}
	if !strings.Contains(out.String(), "already matches") {
		t.Fatalf("a second render must report no change, got:\n%s", out.String())
	}
}

// Generation must never lose project data that config does not carry yet. The workspace
// dir and role only ever existed as prose in the old table.
func TestRosterRenderPreservesLegacyCells(t *testing.T) {
	root := renderFixture(t, "[roster.claude-1]\nadapter = \"claude\"\n", coopSkeleton)
	var out, errb strings.Builder
	if code := rosterRender(root, false, true, false, &out, &errb); code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, errb.String())
	}
	got, _ := os.ReadFile(filepath.Join(root, "parley-deck", "COOPERATION.md"))
	if !strings.Contains(string(got), "../claude/") || !strings.Contains(string(got), "facilitator") {
		t.Fatalf("render dropped values only the legacy table carried:\n%s", got)
	}
}

// Only the table is regenerated: the prose around it explains WHY the section is
// generated and must survive, as must every other section.
func TestRosterRenderTouchesOnlyTheTable(t *testing.T) {
	root := renderFixture(t, "[roster.claude-1]\nadapter = \"claude\"\n[roster.kimi-1]\nadapter = \"kimi\"\n", coopSkeleton)
	var out, errb strings.Builder
	rosterRender(root, false, true, false, &out, &errb)
	got, _ := os.ReadFile(filepath.Join(root, "parley-deck", "COOPERATION.md"))
	for _, want := range []string{"# Protocol", "Prose explaining that this table is generated.", "## 3. Next section", "Body."} {
		if !strings.Contains(string(got), want) {
			t.Errorf("render removed %q:\n%s", want, got)
		}
	}
	if !strings.Contains(string(got), "`kimi-1`") {
		t.Errorf("render did not add the new member:\n%s", got)
	}
}

// Inactive rows sort after active ones, then by ID, so ordering cannot depend on map
// iteration order.
func TestRosterRenderOrderingIsDeterministic(t *testing.T) {
	toml := "[roster.zulu-1]\nadapter = \"claude\"\n[roster.alpha-1]\nadapter = \"kimi\"\n[roster.retired-1]\nadapter = \"claude\"\nactive = false\n"
	root := renderFixture(t, toml, coopSkeleton)
	table, _, err := renderRosterTable(root, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	iAlpha := strings.Index(table, "`alpha-1`")
	iZulu := strings.Index(table, "`zulu-1`")
	iRetired := strings.Index(table, "`retired-1`")
	if !(iAlpha < iZulu && iZulu < iRetired) {
		t.Fatalf("want active alpha-1 < zulu-1 < inactive retired-1, got:\n%s", table)
	}
}
