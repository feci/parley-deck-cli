package protocol

import (
	"crypto/sha256"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"

	"parley-deck-cli/internal/quota"
)

type ActiveExecution struct {
	Path       string
	Transition *quota.RoleTransition
	Chain      []string
}

// PendingHandover permits only the selected successor's handover action.
type PendingHandover struct{ Transition *quota.RoleTransition }

func (e *PendingHandover) Error() string {
	return "selected role handover pending: " + e.Transition.Path()
}

// ResolveActiveExecution selects only validated history. It never uses filenames
// to select a successor or falls back after an invalid handover.
func ResolveActiveExecution(ideaDir string) (ActiveExecution, error) {
	h, err := quota.ReadHistory(ideaDir)
	if err != nil {
		return ActiveExecution{}, err
	}
	return resolveActiveExecution(ideaDir, h)
}

func resolveActiveExecution(ideaDir string, h *quota.History) (ActiveExecution, error) {
	a := ActiveExecution{Path: filepath.Join(ideaDir, "IMPLEMENTATION.md")}
	var transitions []*quota.RoleTransition
	if h != nil {
		if h.Kickoff != nil && h.Kickoff.Transition != nil && h.Kickoff.Transition.Decision.RoleTransition != nil {
			transitions = append(transitions, h.Kickoff.Transition.Decision.RoleTransition)
		}
		for _, b := range h.Batches {
			if b.Decision.RoleTransition != nil {
				transitions = append(transitions, b.Decision.RoleTransition)
			}
		}
	}
	for i, r := range transitions {
		if r.Base != "none" {
			base := filepath.Join(ideaDir, filepath.FromSlash(r.Base))
			if err := regularExecutionFile(base); err != nil {
				return a, fmt.Errorf("handover base: %w", err)
			}
			raw, err := os.ReadFile(base)
			if err != nil || fmt.Sprintf("%x", sha256.Sum256(raw)) != r.BaseSHA256 {
				return a, fmt.Errorf("missing or changed handover base %s", r.Base)
			}
			if i == 0 {
				a.Chain = append(a.Chain, base)
			}
		}
		a.Path, a.Transition = filepath.Join(ideaDir, filepath.FromSlash(r.Path())), r
		if err := regularExecutionFile(a.Path); os.IsNotExist(err) && i == len(transitions)-1 {
			return a, &PendingHandover{Transition: r}
		} else if err != nil {
			return a, err
		}
		if err := ValidateRoleHandover(a.Path, filepath.Base(ideaDir), *r); err != nil {
			return a, err
		}
		a.Chain = append(a.Chain, a.Path)
	}
	if len(a.Chain) == 0 {
		a.Chain = []string{a.Path}
	}
	return a, nil
}

func regularExecutionFile(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("execution authority is not a regular file: %s", path)
	}
	parent, err := os.Lstat(filepath.Dir(path))
	if err != nil || !parent.IsDir() || parent.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("unsafe execution authority directory")
	}
	return nil
}

func ValidateRoleHandover(path, idea string, r quota.RoleTransition) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	m, err := strictArtifactMeta(raw)
	if err != nil {
		return fmt.Errorf("invalid role handover: %w", err)
	}
	want := map[string]string{"idea": idea, "transition": r.ID, "author": r.Successor,
		"implementer": r.After.Implementer, "organizer": r.After.Organizer, "base": r.Base,
		"base-sha256": r.BaseSHA256, "phase": r.Phase, "takeover-branch": r.Branch, "takeover-head": r.Head}
	for key, value := range want {
		if strings.Trim(m[key], "\"'") != value {
			return fmt.Errorf("handover %s contradicts immutable transition", key)
		}
	}
	for _, heading := range []string{"## Takeover", "## Resume state"} {
		if !HasNonEmptySection(string(raw), heading) {
			return fmt.Errorf("handover lacks %s", heading)
		}
	}
	status := strings.Trim(m["status"], "\"'")
	if status == "" || !ValidImplementationStatus(status) || m["date"] == "" {
		return fmt.Errorf("handover status/date is invalid")
	}
	// Before implementation, takeover is in-progress and supplies no review proof.
	implementation := r.Base == "IMPLEMENTATION.md" || m["started"] != "" || status != "in-progress"
	if r.Base != "none" && r.Base != "IMPLEMENTATION.md" {
		baseMeta, err := ReadFrontmatter(filepath.Join(filepath.Dir(filepath.Dir(path)), filepath.FromSlash(r.Base)))
		if err != nil {
			return err
		}
		implementation = implementation || baseMeta["started"] != ""
	}
	if implementation {
		for _, key := range []string{"started", "completed", "branch", "head-commit", "design-pr", "implementation-pr"} {
			if m[key] == "" {
				return fmt.Errorf("active execution record lacks %s", key)
			}
		}
		for _, heading := range []string{"## Summary of work", "## Implementation plan / checklist", "## Deviations from FINAL.md", "## Notes for reviewers", "## Progress", "## Decision Log", "## Surprises & Discoveries", "## Validation evidence", "## Outcomes & Retrospective"} {
			if !HasNonEmptySection(string(raw), heading) {
				return fmt.Errorf("active execution record lacks %s", heading)
			}
		}
	}
	return nil
}

// An early handover is role authority, not evidence that Phase 5 has started.
func ExecutionStarted(ideaDir string) (bool, error) {
	a, err := ResolveActiveExecution(ideaDir)
	if err != nil {
		return false, err
	}
	m, err := ReadFrontmatter(a.Path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return a.Transition == nil || m["started"] != "", nil
}

// This grammar rejects partial or duplicate metadata instead of repairing it.
func strictArtifactMeta(raw []byte) (map[string]string, error) {
	text := string(raw)
	if !strings.HasPrefix(text, "---\n") || !strings.HasSuffix(text, "\n") {
		return nil, fmt.Errorf("incomplete artifact framing")
	}
	front, _, ok := strings.Cut(text[4:], "\n---\n")
	if !ok {
		return nil, fmt.Errorf("incomplete frontmatter")
	}
	var document yaml.Node
	if err := yaml.Unmarshal([]byte(front), &document); err != nil {
		return nil, err
	}
	if len(document.Content) != 1 || document.Content[0].Kind != yaml.MappingNode {
		return nil, fmt.Errorf("metadata must be a mapping")
	}
	m := map[string]string{}
	nodes := document.Content[0].Content
	for i := 0; i < len(nodes); i += 2 {
		key, value := nodes[i], nodes[i+1]
		if key.Kind != yaml.ScalarNode || key.Tag != "!!str" || key.Value == "" {
			return nil, fmt.Errorf("invalid metadata key")
		}
		if _, exists := m[key.Value]; exists {
			return nil, fmt.Errorf("duplicate artifact metadata")
		}
		if value.Kind != yaml.ScalarNode && value.Kind != yaml.SequenceNode {
			return nil, fmt.Errorf("unsupported metadata construct")
		}
		for _, child := range value.Content {
			if child.Kind != yaml.ScalarNode {
				return nil, fmt.Errorf("unsupported metadata sequence")
			}
		}
		m[key.Value] = value.Value
	}
	return m, nil
}

func ExecutionPath(ideaDir string) (string, error) {
	a, err := ResolveActiveExecution(ideaDir)
	return a.Path, err
}

func ReadExecution(ideaDir string) ([]byte, error) {
	p, err := ExecutionPath(ideaDir)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(p)
}

func ReadExecutionFrontmatter(ideaDir string) (map[string]string, error) {
	p, err := ExecutionPath(ideaDir)
	if err != nil {
		return nil, err
	}
	return ReadFrontmatter(p)
}

func StatExecution(ideaDir string) (os.FileInfo, error) {
	p, err := ExecutionPath(ideaDir)
	if err != nil {
		return nil, err
	}
	return os.Stat(p)
}

func ExecutionContext(ideaDir string) (string, error) {
	a, err := ResolveActiveExecution(ideaDir)
	if err != nil {
		return "", err
	}
	var out strings.Builder
	for _, path := range a.Chain {
		raw, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(&out, "\nExecution record: %s\n%s\n", path, raw)
	}
	return out.String(), nil
}
