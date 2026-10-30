package protocol

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"parley-deck-cli/internal/quota"
)

var draftSignoff = regexp.MustCompile(`(?m)^\s*(?:#{1,6}\s+Signoff\b|Status:\s*(?:✅|🟡|❌|ACCEPT|BLOCK))`)

// PlanDraftRecovery never edits a draft. The batch binds any later relocation.
func PlanDraftRecovery(ideaDir string, candidates []quota.Candidate) ([]quota.DraftArchive, error) {
	var archives []quota.DraftArchive
	for _, name := range []string{"consensus.md", "FINAL.md", "review/consensus.md"} {
		path := filepath.Join(ideaDir, filepath.FromSlash(name))
		raw, err := os.ReadFile(path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		meta, err := strictArtifactMeta(raw)
		if err != nil {
			return nil, fmt.Errorf("shared draft %s integrity: %w", name, err)
		}
		author := strings.Trim(meta["drafted-by"], "\"'")
		if author == "" {
			author = strings.Trim(meta["author"], "\"'")
		}
		lost := false
		for _, c := range candidates {
			if c.Agent == author && c.Evidence.RuleID == quota.ProviderExhaustionRule {
				lost = true
			}
		}
		if !lost {
			continue
		}
		if strings.Trim(meta["idea"], "\"'") != filepath.Base(ideaDir) || author == "" || (meta["author"] != "" && meta["drafted-by"] != "" && strings.Trim(meta["author"], "\"'") != author) {
			return nil, fmt.Errorf("shared draft %s has contradictory identity", name)
		}
		body := string(raw)
		fences := 0
		for _, line := range strings.Split(body, "\n") {
			if strings.HasPrefix(line, "```") || strings.HasPrefix(line, "~~~") {
				fences++
			}
		}
		if fences%2 != 0 {
			return nil, fmt.Errorf("shared draft %s has a truncated construct", name)
		}
		complete := false
		if name == "FINAL.md" {
			complete = ValidateFinal(body, filepath.Base(ideaDir)) == ""
		} else {
			headings := []string{"## Agreed decisions", "## Agreed trade-offs", "## Open items deferred to implementation"}
			if name == "review/consensus.md" {
				headings = []string{"## Agreed fixes", "## Deferred follow-ups", "## Dismissed findings"}
			}
			complete = meta["date"] != "" && HasHeadingLine(body, "## Signoffs")
			for _, heading := range headings {
				complete = complete && HasNonEmptySection(body, heading)
			}
		}
		if complete {
			continue
		}
		if draftSignoff.MatchString(body) || strings.ContainsAny(body, "✅🟡❌") || meta["status"] == "final" || meta["status"] == "complete" {
			return nil, fmt.Errorf("signed or frozen incomplete draft %s requires integrity repair", name)
		}
		hash := fmt.Sprintf("%x", sha256.Sum256(raw))
		archives = append(archives, quota.DraftArchive{Path: name, Author: author, SHA256: hash, Archive: "draft-archive/" + strings.ReplaceAll(name, "/", "-") + "-" + hash + ".md"})
	}
	return archives, nil
}
