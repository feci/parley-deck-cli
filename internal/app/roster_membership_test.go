package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// deckWith writes a scratch deck and an isolated machine config, returning the deck root.
func deckWith(t *testing.T, deckTOML, machineTOML string) string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "parley-deck"), 0o755); err != nil {
		t.Fatal(err)
	}
	if deckTOML != "" {
		if err := os.WriteFile(filepath.Join(root, "parley-deck", "agents.toml"), []byte(deckTOML), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	home := t.TempDir()
	t.Setenv("PARLEY_HOME", home)
	if machineTOML != "" {
		if err := os.WriteFile(filepath.Join(home, "agents.toml"), []byte(machineTOML), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// A deck that declares two members must resolve to exactly two, no matter how many the
// machine configures. Before the fix, membership was the layered UNION, so this deck
// resolved to five and — because participant selection reads the same authority — a run
// would have deliberated with three agents nobody put on the deck.
func TestDeckMembershipIsTheDeckFileNotTheLayeredUnion(t *testing.T) {
	root := deckWith(t,
		"[roster.claude-1]\nadapter = \"claude\"\n[roster.kimi-1]\nadapter = \"kimi\"\n",
		"[roster.claude-1]\nadapter = \"claude\"\nmodel = \"machine-model\"\n"+
			"[roster.codex-1]\nadapter = \"codex\"\n[roster.hermes-1]\nadapter = \"hermes\"\n"+
			"[roster.opencode-1]\nadapter = \"opencode\"\n")

	active, inactive, ok := RosterMembership(root)
	if !ok {
		t.Fatal("membership not resolved")
	}
	if len(active)+len(inactive) != 2 {
		t.Fatalf("membership = %d agents, want 2 (deck declares claude-1, kimi-1); active=%v inactive=%v",
			len(active)+len(inactive), active, inactive)
	}
	for _, id := range []string{"codex-1", "hermes-1", "opencode-1"} {
		if active[id] || inactive[id] {
			t.Errorf("%s came from the machine config but is not on this deck", id)
		}
	}
	// Values still inherit: membership is the deck's, the model is the machine's.
	rows, err := resolveRoster(root, nil, rosterViewOpts{})
	if err != nil {
		t.Fatal(err)
	}
	var got string
	for _, r := range rows {
		if r.Agent == "claude-1" {
			got = r.Model
		}
	}
	if got != "machine-model" {
		t.Errorf("claude-1 model = %q, want the inherited machine value %q", got, "machine-model")
	}
}

// A deck with no roster of its own may DISPLAY the machine roster, but every row must say
// so. Silently adopting it is how an accident of one machine became a committed §2 table.
func TestRosterlessDeckMarksInheritedRows(t *testing.T) {
	root := deckWith(t, "", "[roster.claude-1]\nadapter = \"claude\"\n")
	rows, err := resolveRoster(root, nil, rosterViewOpts{})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("rows = %d, want 1", len(rows))
	}
	if !hasStatus(rows[0], "inherited-roster") {
		t.Errorf("inherited row status = %v, want it to contain inherited-roster", rows[0].Status)
	}
}

// render must refuse to write an inherited roster into the committed protocol file.
func TestRenderRefusesToCommitAnInheritedRoster(t *testing.T) {
	root := deckWith(t, "", "[roster.claude-1]\nadapter = \"claude\"\n")
	if _, _, err := renderRosterTable(root, nil, false); err == nil {
		t.Fatal("render accepted an inherited roster without --adopt-inherited")
	}
	if _, _, err := renderRosterTable(root, nil, true); err != nil {
		t.Fatalf("--adopt-inherited should permit it: %v", err)
	}
}

// A §2-only ID is reported `unmapped`, never dropped and never auto-added.
func TestSection2OnlyIDIsReportedNotDropped(t *testing.T) {
	root := deckWith(t, "[roster.claude-1]\nadapter = \"claude\"\n", "")
	coop := "## 2. Active roster\n\n| Agent ID | Workspace dir | Role |\n| --- | --- | --- |\n" +
		"| `claude-1` | . | participant |\n| `ghost-1` | ghosts/ | reviewer |\n"
	if err := os.WriteFile(filepath.Join(root, "parley-deck", "COOPERATION.md"), []byte(coop), 0o644); err != nil {
		t.Fatal(err)
	}
	rows := section2OnlyRows(root, map[string]bool{"claude-1": true})
	if len(rows) != 1 || rows[0].Agent != "ghost-1" {
		t.Fatalf("section2-only rows = %+v, want exactly ghost-1", rows)
	}
	if !hasStatus(rows[0], "unmapped") || !hasStatus(rows[0], "section2-only") {
		t.Errorf("ghost-1 status = %v, want unmapped and section2-only", rows[0].Status)
	}
}

// The membership gate keys on whether the BLOCK existed, not on which field is written.
// Keying on "+ adapter = " let `roster set new-9 --model X --yes` create a member with no
// second confirmation.
func TestMembershipGateCatchesNewBlockWrittenWithAnyField(t *testing.T) {
	for _, field := range []string{"model", "effort", "speed"} {
		if got := membershipChange([]string{"+ " + field + " = \"x\""}, false, true); got == "" {
			t.Errorf("a new block written with only --%s was not treated as a membership change", field)
		}
	}
	if got := membershipChange([]string{"+ model = \"x\""}, true, true); got != "" {
		t.Errorf("changing an existing member's model is a settings change, got %q", got)
	}
	// A revival must not report itself as a retirement.
	if got := membershipChange([]string{"- active = false", "+ active = true"}, true, false); !strings.Contains(got, "reactivat") {
		t.Errorf("reactivation reported as %q", got)
	}
}

func hasStatus(r rosterRow, code string) bool {
	for _, s := range r.Status {
		if s == code {
			return true
		}
	}
	return false
}
