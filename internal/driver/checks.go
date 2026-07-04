package driver

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// checks.go implements the list-form `checks:` completion contract
// (completion-contracts-evidence-ledger). A scalar `checks:` keeps today's behavior;
// a YAML list of {name, command} activates the contract: the driver runs each criterion
// and writes a per-criterion evidence table into IMPLEMENTATION.md's `## Validation
// evidence` section. This reader is the single place that understands the list shape.

// CheckCriterion is one named completion criterion (expects exit 0).
type CheckCriterion struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
}

// ReadChecksContract inspects the `checks:` frontmatter of 00-prompt.md.
//   - absent or scalar → (nil, false, nil): the caller uses today's scalar path.
//   - a YAML list → (criteria, true, nil) after validation.
//   - a malformed list → (nil, true, err): fail closed (present but invalid).
func ReadChecksContract(ideaDir string) ([]CheckCriterion, bool, error) {
	raw, err := os.ReadFile(filepath.Join(ideaDir, "00-prompt.md"))
	if err != nil {
		return nil, false, nil // no idea file → no contract (legacy)
	}
	fm := extractFrontmatter(string(raw))
	if fm == "" {
		return nil, false, nil
	}
	// Decode only the `checks` node so scalar vs list is detectable.
	var probe struct {
		Checks yaml.Node `yaml:"checks"`
	}
	if err := yaml.Unmarshal([]byte(fm), &probe); err != nil {
		return nil, false, nil // unparseable frontmatter → treat as legacy (scalar path handles it)
	}
	if probe.Checks.Kind == 0 || probe.Checks.Kind == yaml.ScalarNode {
		return nil, false, nil // absent or scalar → legacy
	}
	if probe.Checks.Kind != yaml.SequenceNode {
		return nil, true, fmt.Errorf("checks: must be a scalar command or a list of {name, command}")
	}
	var criteria []CheckCriterion
	if err := probe.Checks.Decode(&criteria); err != nil {
		return nil, true, fmt.Errorf("checks: list is malformed: %w", err)
	}
	seen := map[string]bool{}
	for i, c := range criteria {
		c.Name = strings.TrimSpace(c.Name)
		c.Command = strings.TrimSpace(c.Command)
		if c.Name == "" {
			return nil, true, fmt.Errorf("checks[%d]: empty name", i)
		}
		if seen[c.Name] {
			return nil, true, fmt.Errorf("checks: duplicate criterion name %q", c.Name)
		}
		seen[c.Name] = true
		if c.Command == "" {
			return nil, true, fmt.Errorf("checks[%q]: empty command", c.Name)
		}
		criteria[i] = c
	}
	if len(criteria) == 0 {
		return nil, true, fmt.Errorf("checks: list is empty")
	}
	return criteria, true, nil
}

// extractFrontmatter returns the text between the leading `---` fences, or "".
func extractFrontmatter(doc string) string {
	s := strings.TrimLeft(doc, " \t\n")
	if !strings.HasPrefix(s, "---") {
		return ""
	}
	s = s[len("---"):]
	end := strings.Index(s, "\n---")
	if end < 0 {
		return ""
	}
	return s[:end]
}
