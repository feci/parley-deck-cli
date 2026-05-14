package runner

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"parley-deck-cli/internal/protocol"
)

func ValidateArtifact(path, agentID, ideaSlug string, round int) error {
	meta, err := protocol.ReadFrontmatter(path)
	if err != nil {
		return err
	}
	for key, want := range map[string]string{
		"agent": agentID,
		"idea":  ideaSlug,
		"round": strconv.Itoa(round),
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
