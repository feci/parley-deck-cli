---
agent: gemini
idea: repo-map-mvp
round: 1
date: 2026-05-17
---

## Summary

The `repo-map-mvp` should provide a high-density, low-token representation of the repository structure and its public API. From a context-efficiency perspective, the map must prioritize "signal" (exported symbols, core project structure) over "noise" (implementation details, dependency artifacts). The output shape must be dual-purpose: a human-readable Markdown tree for quick inspection and a stable JSON schema for programmatic consumption by future agentic slices.

## Proposed approach

### 1. Command Shape
Implement `parley context map` with the following interface:
- `parley context map [--dir DIR] [--format markdown|json] [--limit N] [--exclude PATTERN]`
- **Default format**: `markdown`.
- **Default limit**: 1000 files (safety guardrail).
- **Default excludes**: `.git/`, `vendor/`, `node_modules/`, `go.sum`, and common binary/build artifacts.

### 2. Data Model
A lean, hierarchical model to represent the repository:
```go
type RepoMap struct {
    Root        string    `json:"root"`
    Files       []File    `json:"files"`
    GeneratedAt time.Time `json:"generated_at"`
    Truncated   bool      `json:"truncated,omitempty"`
}

type File struct {
    Path    string   `json:"path"` // Relative to root
    Symbols []Symbol `json:"symbols,omitempty"`
}

type Symbol struct {
    Name     string `json:"name"`
    Kind     string `json:"kind"` // struct, interface, func, method, type, const
    Line     int    `json:"line"`
    Exported bool   `json:"exported"`
}
```

### 3. Output Shapes
- **Markdown**: An indented tree view. Files containing Go symbols will list them immediately below the filename.
  ```markdown
  - internal/app/
    - app.go
      - struct App
      - func Run(args []string, stdout, stderr io.Writer) int
  ```
- **JSON**: The `RepoMap` struct serialized. This allows agents to "search" the map without parsing text.

### 4. Logic & Guardrails
- **Go Parsing**: Use `go/parser` and `go/ast`. Extract only top-level declarations and methods.
- **Context Efficiency**: By default, only include **Exported** symbols in Go. This significantly reduces token count while preserving the "Public API" map.
- **Tree Traversal**: Use `filepath.WalkDir`. Respect `.gitignore` if present.
- **Limit Behavior**: If the file limit is reached, stop walking and set `Truncated: true`. The CLI should emit a warning to `stderr`.

### 5. Tests
- **Symbol Extraction**: Unit tests in `internal/repomap/golang_test.go` using various Go source snippets (interfaces, nested structs, methods).
- **Ignore Logic**: Tests ensuring `.git` and `vendor` are skipped.
- **Formatters**: Snapshot tests for both Markdown and JSON output.

## Concerns / open questions
- **Language Detection**: For MVP, should we only attempt symbol extraction for `.go` files? (Recommendation: Yes, keep it simple).
- **Relative vs Absolute Paths**: Always use relative paths from the workspace root to ensure determinism across different environments.
- **Signature Detail**: Should we include full function signatures in the map? (Proposal: Include names and parameters for functions/methods, but skip body-related info to save tokens).

## Risks
- **Large Repositories**: Even with a 1000-file limit, a repo with many small files and many symbols might exceed prompt windows. We may need a "depth" limit or "package-only" mode in the future.
- **Binary/Large Files**: The walker must be careful not to "read" large non-text files. It should filter by extension or use a header-sniffing check.
- **Standard Library**: We should decide if we want to map `internal/` packages. (Proposal: Yes, as they are part of the repo's internal map for agents).
