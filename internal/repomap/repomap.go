package repomap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const SchemaVersion = 1

type Options struct {
	Root     string
	MaxFiles int
}

type Map struct {
	SchemaVersion int    `json:"schema_version"`
	Root          string `json:"root"`
	MaxFiles      int    `json:"max_files"`
	Truncated     bool   `json:"truncated"`
	Counts        Counts `json:"counts"`
	Files         []File `json:"files"`
}

type Counts struct {
	Files   int `json:"files"`
	Symbols int `json:"symbols"`
}

type File struct {
	Path       string   `json:"path"`
	Kind       string   `json:"kind"`
	SizeBytes  int64    `json:"size_bytes"`
	Package    string   `json:"package,omitempty"`
	Imports    []string `json:"imports,omitempty"`
	Symbols    []Symbol `json:"symbols,omitempty"`
	ParseError string   `json:"parse_error,omitempty"`
}

type Symbol struct {
	Kind     string `json:"kind"`
	Name     string `json:"name"`
	Receiver string `json:"receiver,omitempty"`
	Exported bool   `json:"exported"`
	Line     int    `json:"line"`
}

func Build(opts Options) (Map, error) {
	root := opts.Root
	if root == "" {
		root = "."
	}
	maxFiles := opts.MaxFiles
	if maxFiles <= 0 {
		maxFiles = 1000
	}
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return Map{}, err
	}

	candidates, err := collectFiles(rootAbs)
	if err != nil {
		return Map{}, err
	}
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Path < candidates[j].Path
	})

	truncated := false
	if len(candidates) > maxFiles {
		candidates = candidates[:maxFiles]
		truncated = true
	}

	files := make([]File, 0, len(candidates))
	symbols := 0
	for _, candidate := range candidates {
		file := File{
			Path:      candidate.Path,
			Kind:      kindForPath(candidate.Path),
			SizeBytes: candidate.SizeBytes,
		}
		if file.Kind == "go" {
			enrichGoFile(rootAbs, &file)
			symbols += len(file.Symbols)
		}
		files = append(files, file)
	}

	return Map{
		SchemaVersion: SchemaVersion,
		Root:          ".",
		MaxFiles:      maxFiles,
		Truncated:     truncated,
		Counts: Counts{
			Files:   len(files),
			Symbols: symbols,
		},
		Files: files,
	}, nil
}

type fileCandidate struct {
	Path      string
	SizeBytes int64
}

var ignoredDirNames = map[string]bool{
	".cache":       true,
	".git":         true,
	".gocache":     true,
	".gomodcache":  true,
	".idea":        true,
	".next":        true,
	".venv":        true,
	".vscode":      true,
	"__pycache__":  true,
	"build":        true,
	"dist":         true,
	"node_modules": true,
	"target":       true,
	"vendor":       true,
}

func collectFiles(rootAbs string) ([]fileCandidate, error) {
	var files []fileCandidate
	err := filepath.WalkDir(rootAbs, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == rootAbs {
			return nil
		}
		rel, err := filepath.Rel(rootAbs, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if entry.IsDir() {
			if shouldSkipDir(rel, entry.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		files = append(files, fileCandidate{Path: rel, SizeBytes: info.Size()})
		return nil
	})
	return files, err
}

func shouldSkipDir(rel, name string) bool {
	if rel == "parley-deck/runs" || strings.HasPrefix(rel, "parley-deck/runs/") {
		return true
	}
	return ignoredDirNames[name]
}

func kindForPath(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return "go"
	case ".md", ".markdown":
		return "markdown"
	case ".json":
		return "json"
	case ".txt", ".toml", ".yaml", ".yml", ".sh", ".bash", ".zsh", ".mod", ".sum":
		return "text"
	default:
		return "other"
	}
}

func enrichGoFile(rootAbs string, file *File) {
	data, err := os.ReadFile(filepath.Join(rootAbs, filepath.FromSlash(file.Path)))
	if err != nil {
		file.ParseError = err.Error()
		return
	}
	fset := token.NewFileSet()
	parsed, err := parser.ParseFile(fset, file.Path, data, parser.ParseComments)
	if err != nil {
		file.ParseError = err.Error()
	}
	if parsed == nil {
		return
	}
	file.Package = parsed.Name.Name
	file.Imports = extractImports(parsed)
	file.Symbols = extractSymbols(fset, parsed)
}

func extractImports(file *ast.File) []string {
	imports := make([]string, 0, len(file.Imports))
	for _, spec := range file.Imports {
		path, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			path = strings.Trim(spec.Path.Value, `"`)
		}
		imports = append(imports, path)
	}
	sort.Strings(imports)
	return imports
}

func extractSymbols(fset *token.FileSet, file *ast.File) []Symbol {
	var symbols []Symbol
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			kind := "func"
			receiver := ""
			if d.Recv != nil && len(d.Recv.List) > 0 {
				kind = "method"
				receiver = receiverName(d.Recv.List[0].Type)
			}
			symbols = append(symbols, Symbol{
				Kind:     kind,
				Name:     d.Name.Name,
				Receiver: receiver,
				Exported: ast.IsExported(d.Name.Name),
				Line:     fset.Position(d.Pos()).Line,
			})
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					symbols = append(symbols, Symbol{
						Kind:     "type",
						Name:     s.Name.Name,
						Exported: ast.IsExported(s.Name.Name),
						Line:     fset.Position(s.Pos()).Line,
					})
				case *ast.ValueSpec:
					kind := "var"
					if d.Tok == token.CONST {
						kind = "const"
					}
					for _, name := range s.Names {
						symbols = append(symbols, Symbol{
							Kind:     kind,
							Name:     name.Name,
							Exported: ast.IsExported(name.Name),
							Line:     fset.Position(name.Pos()).Line,
						})
					}
				}
			}
		}
	}
	sort.SliceStable(symbols, func(i, j int) bool {
		if symbols[i].Line != symbols[j].Line {
			return symbols[i].Line < symbols[j].Line
		}
		if symbols[i].Kind != symbols[j].Kind {
			return symbols[i].Kind < symbols[j].Kind
		}
		if symbols[i].Receiver != symbols[j].Receiver {
			return symbols[i].Receiver < symbols[j].Receiver
		}
		return symbols[i].Name < symbols[j].Name
	})
	return symbols
}

func receiverName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		return receiverName(e.X)
	case *ast.SelectorExpr:
		return e.Sel.Name
	case *ast.IndexExpr:
		return receiverName(e.X)
	case *ast.IndexListExpr:
		return receiverName(e.X)
	default:
		return fmt.Sprintf("%T", expr)
	}
}

func RenderJSON(m Map, w io.Writer) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.Write(append(data, '\n'))
	return err
}

func RenderMarkdown(m Map, w io.Writer) error {
	var b bytes.Buffer
	fmt.Fprintln(&b, "# Repository Map")
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "- schema_version: %d\n", m.SchemaVersion)
	fmt.Fprintf(&b, "- root: %s\n", m.Root)
	fmt.Fprintf(&b, "- files: %d\n", m.Counts.Files)
	fmt.Fprintf(&b, "- symbols: %d\n", m.Counts.Symbols)
	fmt.Fprintf(&b, "- max_files: %d\n", m.MaxFiles)
	fmt.Fprintf(&b, "- truncated: %t\n", m.Truncated)
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "## Files")
	for _, file := range m.Files {
		fmt.Fprintf(&b, "\n- `%s` (%s, %d bytes)\n", file.Path, file.Kind, file.SizeBytes)
		if file.Package != "" {
			fmt.Fprintf(&b, "  - package: `%s`\n", file.Package)
		}
		if len(file.Imports) > 0 {
			fmt.Fprintf(&b, "  - imports: `%s`\n", strings.Join(file.Imports, "`, `"))
		}
		if file.ParseError != "" {
			fmt.Fprintf(&b, "  - parse_error: `%s`\n", file.ParseError)
		}
		for _, symbol := range file.Symbols {
			name := symbol.Name
			if symbol.Receiver != "" {
				name = symbol.Receiver + "." + name
			}
			fmt.Fprintf(&b, "  - %s `%s` line %d exported=%t\n", symbol.Kind, name, symbol.Line, symbol.Exported)
		}
	}
	_, err := w.Write(b.Bytes())
	return err
}
