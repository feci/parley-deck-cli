package driver

import (
	"os"
	"path/filepath"
	"testing"
)

func writeIdea(t *testing.T, fm string) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "00-prompt.md"), []byte(fm), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestReadChecksContractAbsent(t *testing.T) {
	dir := writeIdea(t, "---\nidea: x\nstatus: round-01\n---\n\n## Problem\n")
	c, isList, err := ReadChecksContract(dir)
	if err != nil || isList || c != nil {
		t.Fatalf("absent checks: got c=%v isList=%v err=%v", c, isList, err)
	}
}

func TestReadChecksContractScalarIsLegacy(t *testing.T) {
	dir := writeIdea(t, "---\nidea: x\nchecks: go test ./...\nstatus: round-01\n---\n")
	c, isList, err := ReadChecksContract(dir)
	if err != nil || isList || c != nil {
		t.Fatalf("scalar checks should be legacy: c=%v isList=%v err=%v", c, isList, err)
	}
}

func TestReadChecksContractList(t *testing.T) {
	dir := writeIdea(t, "---\nidea: x\nchecks:\n  - name: unit\n    command: go test ./...\n  - name: build\n    command: go build ./...\nstatus: round-01\n---\n")
	c, isList, err := ReadChecksContract(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !isList || len(c) != 2 {
		t.Fatalf("list: isList=%v len=%d", isList, len(c))
	}
	if c[0].Name != "unit" || c[1].Command != "go build ./..." {
		t.Fatalf("parsed = %+v", c)
	}
}

func TestReadChecksContractMalformedFailsClosed(t *testing.T) {
	cases := []string{
		"---\nchecks:\n  - name: a\n    command: x\n  - name: a\n    command: y\n---\n", // dup name
		"---\nchecks:\n  - name: \"\"\n    command: x\n---\n",                              // empty name
		"---\nchecks:\n  - name: a\n    command: \"\"\n---\n",                              // empty command
		"---\nchecks: []\n---\n",                                                          // empty list
	}
	for i, fm := range cases {
		dir := writeIdea(t, fm)
		_, isList, err := ReadChecksContract(dir)
		if !isList || err == nil {
			t.Fatalf("case %d should fail closed: isList=%v err=%v", i, isList, err)
		}
	}
}
