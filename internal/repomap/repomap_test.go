package repomap

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildExtractsGoSymbolsAndParseErrors(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "internal/sample/sample.go", `package sample

import (
	"fmt"
	"strings"
)

type Thing struct{}

const ExportedConst = "value"

var localVar = 1

func Do() {
	fmt.Println(strings.TrimSpace(" ok "))
}

func hidden() {}

func (t *Thing) Run() {}
`)
	writeFile(t, root, "internal/sample/broken.go", "package sample\nfunc broken(")
	writeFile(t, root, "README.md", "# sample\n")

	m, err := Build(Options{Root: root, MaxFiles: 100})
	if err != nil {
		t.Fatal(err)
	}
	if m.SchemaVersion != 1 || m.Root != "." || m.Truncated {
		t.Fatalf("unexpected map metadata: %+v", m)
	}
	sample := findFile(t, m, "internal/sample/sample.go")
	if sample.Kind != "go" || sample.Package != "sample" {
		t.Fatalf("sample file=%+v", sample)
	}
	if got := strings.Join(sample.Imports, ","); got != "fmt,strings" {
		t.Fatalf("imports=%q", got)
	}
	wantSymbols := map[string]struct {
		kind     string
		receiver string
		exported bool
	}{
		"Thing":         {kind: "type", exported: true},
		"ExportedConst": {kind: "const", exported: true},
		"localVar":      {kind: "var", exported: false},
		"Do":            {kind: "func", exported: true},
		"hidden":        {kind: "func", exported: false},
		"Run":           {kind: "method", receiver: "Thing", exported: true},
	}
	for name, want := range wantSymbols {
		got := findSymbol(t, sample, name)
		if got.Kind != want.kind || got.Receiver != want.receiver || got.Exported != want.exported || got.Line == 0 {
			t.Fatalf("symbol %s=%+v want %+v with line", name, got, want)
		}
	}

	broken := findFile(t, m, "internal/sample/broken.go")
	if broken.ParseError == "" {
		t.Fatalf("broken file missing parse error: %+v", broken)
	}
	if findFile(t, m, "README.md").Kind != "markdown" {
		t.Fatalf("README kind mismatch")
	}
}

func TestBuildIgnoresTransientDirsAndTruncatesDeterministically(t *testing.T) {
	root := t.TempDir()
	for _, path := range []string{
		".git/config",
		"node_modules/pkg/index.js",
		"parley-deck/runs/run/events.jsonl",
		"b.txt",
		"a.txt",
		"c.txt",
	} {
		writeFile(t, root, path, "x\n")
	}

	m, err := Build(Options{Root: root, MaxFiles: 2})
	if err != nil {
		t.Fatal(err)
	}
	if !m.Truncated || m.Counts.Files != 2 || len(m.Files) != 2 {
		t.Fatalf("unexpected truncation metadata: %+v", m)
	}
	got := []string{m.Files[0].Path, m.Files[1].Path}
	if strings.Join(got, ",") != "a.txt,b.txt" {
		t.Fatalf("paths=%v", got)
	}
	for _, file := range m.Files {
		if strings.Contains(file.Path, ".git") || strings.Contains(file.Path, "node_modules") || strings.Contains(file.Path, "parley-deck/runs") {
			t.Fatalf("ignored path included: %+v", file)
		}
	}
}

func TestRenderersAreDeterministic(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "main.go", "package main\n\nfunc main() {}\n")
	m, err := Build(Options{Root: root, MaxFiles: 100})
	if err != nil {
		t.Fatal(err)
	}

	var md1, md2 bytes.Buffer
	if err := RenderMarkdown(m, &md1); err != nil {
		t.Fatal(err)
	}
	if err := RenderMarkdown(m, &md2); err != nil {
		t.Fatal(err)
	}
	if md1.String() != md2.String() || !strings.Contains(md1.String(), "func `main`") {
		t.Fatalf("markdown mismatch or missing symbol:\n%s\n---\n%s", md1.String(), md2.String())
	}

	var json1, json2 bytes.Buffer
	if err := RenderJSON(m, &json1); err != nil {
		t.Fatal(err)
	}
	if err := RenderJSON(m, &json2); err != nil {
		t.Fatal(err)
	}
	if json1.String() != json2.String() {
		t.Fatalf("json mismatch:\n%s\n---\n%s", json1.String(), json2.String())
	}
	var decoded Map
	if err := json.Unmarshal(json1.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, json1.String())
	}
	if decoded.SchemaVersion != SchemaVersion {
		t.Fatalf("decoded schema=%d", decoded.SchemaVersion)
	}
}

func writeFile(t *testing.T, root, rel, body string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func findFile(t *testing.T, m Map, path string) File {
	t.Helper()
	for _, file := range m.Files {
		if file.Path == path {
			return file
		}
	}
	t.Fatalf("file %s not found in %+v", path, m.Files)
	return File{}
}

func findSymbol(t *testing.T, file File, name string) Symbol {
	t.Helper()
	for _, symbol := range file.Symbols {
		if symbol.Name == name {
			return symbol
		}
	}
	t.Fatalf("symbol %s not found in %+v", name, file.Symbols)
	return Symbol{}
}
