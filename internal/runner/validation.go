package runner

import (
	"fmt"
	"os"
	"strings"

	"parley-deck-cli/internal/protocol"
)

func ValidateRoundOneArtifact(path, agentID, ideaSlug string) error {
	meta, err := readRoundOneFrontmatter(path)
	if err != nil {
		return err
	}
	for key, want := range map[string]string{
		"agent": agentID,
		"idea":  ideaSlug,
		"round": "1",
	} {
		if got := strings.Trim(strings.TrimSpace(meta[key]), `"'`); got != want {
			return fmt.Errorf("%s frontmatter %s=%q, want %q", path, key, got, want)
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	body := string(data)
	for _, section := range []string{
		"## Summary",
		"## Proposed approach",
		"## Concerns / open questions",
		"## Risks",
	} {
		if !strings.Contains(body, section) {
			return fmt.Errorf("%s missing required section %q", path, section)
		}
	}
	return nil
}

func readRoundOneFrontmatter(path string) (map[string]string, error) {
	meta, err := protocol.ReadFrontmatter(path)
	if err != nil {
		return nil, err
	}
	if hasAnyFrontmatter(meta) {
		return meta, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parseLeadingFrontmatterWithoutOpeningFence(string(data)), nil
}

func hasAnyFrontmatter(meta map[string]string) bool {
	for _, key := range []string{"agent", "idea", "round"} {
		if strings.TrimSpace(meta[key]) != "" {
			return true
		}
	}
	return false
}

func parseLeadingFrontmatterWithoutOpeningFence(data string) map[string]string {
	lines := strings.Split(data, "\n")
	meta := map[string]string{}
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			if i == 0 || len(meta) == 0 {
				return map[string]string{}
			}
			return meta
		}
		if trimmed == "" {
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			return map[string]string{}
		}
		key = strings.TrimSpace(key)
		if key == "" || strings.ContainsAny(key, " \t") {
			return map[string]string{}
		}
		meta[key] = strings.TrimSpace(value)
	}
	return map[string]string{}
}
