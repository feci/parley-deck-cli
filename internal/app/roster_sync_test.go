package app

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func syncFixture(t *testing.T, machine, deck string) (root string) {
	t.Helper()
	// PARLEY_HOME names the central config DIRECTORY itself. Seeding
	// $PARLEY_HOME/.parley/agents.toml is what let a machine-scope write land in a file
	// no resolver reads while still reporting success.
	home := t.TempDir()
	if err := os.WriteFile(filepath.Join(home, "agents.toml"), []byte(machine), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PARLEY_HOME", home)
	return writeDeckTOML(t, deck)
}

const machineRoster = "[roster.claude-1]\nadapter = \"claude\"\nmodel = \"opus-5\"\n[roster.kimi-1]\nadapter = \"kimi\"\nmodel = \"k3\"\n"

// Rebase removes a deck override that merely restates the machine value, so the deck
// goes back to INHERITING. A copy-down would freeze the deck the day after it was
// written and every later central change would stop at its door.
func TestRosterSyncRemovesRedundantOverridesSoTheDeckInherits(t *testing.T) {
	root := syncFixture(t, machineRoster, "[roster.claude-1]\nadapter = \"claude\"\nmodel = \"opus-5\"\n")
	var out, errb bytes.Buffer
	if code := runRoster([]string{"sync", "--dir", root, "--yes"}, &out, &errb); code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, errb.String())
	}
	got, _ := os.ReadFile(filepath.Join(root, "parley-deck", "agents.toml"))
	if strings.Contains(string(got), "opus-5") {
		t.Fatalf("redundant override survived:\n%s", got)
	}
	if !strings.Contains(string(got), "[roster.claude-1]") {
		t.Fatalf("membership must survive a rebase:\n%s", got)
	}
}

// A deliberate pin is the entire point of the layering. Sync may remove it, but never
// silently: the preview names it and prints the exact --keep needed to retain it.
func TestRosterSyncNamesDeliberatePinsAndTheKeepFlag(t *testing.T) {
	root := syncFixture(t, machineRoster, "[roster.kimi-1]\nadapter = \"kimi\"\nmodel = \"k2-legacy\"\n")
	var out, errb bytes.Buffer
	if code := runRoster([]string{"sync", "--dir", root}, &out, &errb); code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, errb.String())
	}
	s := out.String()
	if !strings.Contains(s, "DELIBERATE PINS") || !strings.Contains(s, "k2-legacy") {
		t.Fatalf("pin not enumerated:\n%s", s)
	}
	if !strings.Contains(s, "--keep kimi-1.model") {
		t.Fatalf("preview must print the exact --keep to retain it:\n%s", s)
	}
	if !strings.Contains(s, "Nothing was written") {
		t.Fatalf("preview is the default:\n%s", s)
	}
}

func TestRosterSyncKeepExemptsAPin(t *testing.T) {
	root := syncFixture(t, machineRoster, "[roster.kimi-1]\nadapter = \"kimi\"\nmodel = \"k2-legacy\"\n")
	var out, errb bytes.Buffer
	if code := runRoster([]string{"sync", "--dir", root, "--keep", "kimi-1.model", "--yes"}, &out, &errb); code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, errb.String())
	}
	got, _ := os.ReadFile(filepath.Join(root, "parley-deck", "agents.toml"))
	if !strings.Contains(string(got), "k2-legacy") {
		t.Fatalf("--keep must exempt the pin:\n%s", got)
	}
	if strings.Contains(string(got), "adapter = \"kimi\"") {
		t.Fatalf("the redundant adapter override should still have been rebased away:\n%s", got)
	}
}

// One direction only. A deck-only value must never travel up into the machine file.
func TestRosterSyncNeverWritesTheMachineFile(t *testing.T) {
	root := syncFixture(t, machineRoster, "[roster.kimi-1]\nadapter = \"kimi\"\nmodel = \"k2-legacy\"\n")
	machinePath := filepath.Join(os.Getenv("PARLEY_HOME"), "agents.toml")
	before, _ := os.ReadFile(machinePath)
	var out, errb bytes.Buffer
	runRoster([]string{"sync", "--dir", root, "--yes"}, &out, &errb)
	after, _ := os.ReadFile(machinePath)
	if string(before) != string(after) {
		t.Fatalf("sync modified the machine file:\n%s", after)
	}
}

// Round-trip guard for the path bug: a machine-scope write must land in the file the
// resolver reads, and a deck must then inherit it.
func TestRosterSetMachineScopeWritesTheFileTheResolverReads(t *testing.T) {
	home := t.TempDir()
	if err := os.WriteFile(filepath.Join(home, "agents.toml"), []byte("[roster.claude-1]\nadapter = \"claude\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PARLEY_HOME", home)
	root := writeDeckTOML(t, "")

	var out, errb bytes.Buffer
	if code := runRoster([]string{"set", "claude-1", "--dir", root, "--scope", "machine", "--model", "opus-9", "--yes"}, &out, &errb); code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, errb.String())
	}
	got, err := os.ReadFile(filepath.Join(home, "agents.toml"))
	if err != nil {
		t.Fatalf("machine write did not land in the central file: %v", err)
	}
	if !strings.Contains(string(got), "opus-9") {
		t.Fatalf("central file does not carry the write:\n%s", got)
	}
	if _, err := os.Stat(filepath.Join(home, ".parley", "agents.toml")); err == nil {
		t.Fatal("a stray $PARLEY_HOME/.parley/agents.toml was written")
	}
}
