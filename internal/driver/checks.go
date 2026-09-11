package driver

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"parley-deck-cli/internal/evidence"

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

// ObserveChecksContract preserves the original named scope across driver ticks.
// A prior report also activates the gate on legacy runs without a cursor pin;
// deleting the current checks list cannot turn that history into a scalar task.
// This is runtime safety state, not authentication against same-UID tampering.
func ObserveChecksContract(ideaDir, expected string) (string, error) {
	obligation, err := evidence.ReadContractPin(ideaDir)
	if err != nil {
		return "", fmt.Errorf("cannot read original contract witness: %w", err)
	}
	if obligation != "" {
		if expected != "" && expected != obligation {
			return "", fmt.Errorf("cursor differs from original contract witness")
		}
		expected = obligation
	}
	criteria, named, err := ReadChecksContract(ideaDir)
	if err != nil {
		return "", err
	}
	digest := ChecksContractDigest
	current := ""
	if named {
		current = digest(criteria)
	}
	if expected != "" && current != expected {
		return "", fmt.Errorf("original named checks contract was removed or changed")
	}
	if expected == "" {
		_, err := os.Lstat(evidence.ReportPath(ideaDir))
		if err != nil && !os.IsNotExist(err) {
			return "", fmt.Errorf("cannot inspect prior named evidence: %w", err)
		}
		if err == nil {
			if !named {
				return "", fmt.Errorf("prior named evidence exists but its checks contract was removed")
			}
			report, err := evidence.Load(ideaDir)
			if err != nil {
				return "", err
			}
			prior := make([]CheckCriterion, 0, len(report.Records))
			for _, rec := range report.Records {
				prior = append(prior, CheckCriterion{Name: rec.Name, Command: rec.Command.Command})
			}
			if digest(prior) != current {
				return "", fmt.Errorf("current checks differ from the original recorded scope")
			}
		}
	}
	return current, nil
}

// ChecksContractDigest is the canonical normalized names/commands binding.
func ChecksContractDigest(criteria []CheckCriterion) string {
	data, _ := json.Marshal(criteria)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
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
		// A YAML parse error must NOT silently fall back to legacy when `checks:` is
		// written in block/list form (review fix): that would fail open. Detect the
		// list shape syntactically and fail closed; only a scalar/absent checks: falls
		// through to the legacy scalar path.
		if looksLikeChecksList(fm) {
			return nil, true, fmt.Errorf("checks: list is not valid YAML: %w", err)
		}
		return nil, false, nil
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

// checksListRe matches a `checks:` key that starts a block/list (nothing but an
// optional comment after it on the same line, then an indented `-` item).
var checksListRe = regexp.MustCompile(`(?m)^checks:[ \t]*(#.*)?\n[ \t]*-`)

// looksLikeChecksList reports whether the frontmatter writes `checks:` in block/list
// form (as opposed to an inline scalar), used to fail closed on a malformed list.
func looksLikeChecksList(fm string) bool {
	return checksListRe.MatchString(fm)
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
