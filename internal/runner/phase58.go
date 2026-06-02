package runner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/protocol"
)

// RunImplementation runs Phase 5: it launches a single implementer to produce
// IMPLEMENTATION.md (and code on a branch) per FINAL.md. opts.Idea.Participants
// must contain exactly the implementer. Reuses the shared launch machinery.
func RunImplementation(ctx context.Context, opts Options) Result {
	opts.Phase = "implementation"
	opts.ArtifactName = "IMPLEMENTATION.md"
	if opts.RoundLabel == "" {
		opts.RoundLabel = "implementation"
	}
	selected := selectedAgents(opts.Idea.Participants, opts.Agents)
	if len(selected) == 0 {
		return Result{AgentID: "implementer", ExitError: "no implementer available in participants"}
	}
	return runAgent(ctx, opts, selected[0])
}

// RunReviewRound runs Phase 6 review round N: each reviewer (opts.Idea.Participants
// should be the non-implementer reviewers) writes review/round-NN/<agent>.md.
func RunReviewRound(ctx context.Context, opts Options) []Result {
	if opts.Round < 1 {
		opts.Round = 1
	}
	opts.Phase = "review"
	opts.RoundLabel = filepath.Join("review", roundLabel(opts.Round))
	return RunRoundOne(ctx, opts)
}

// BuildImplementationPrompt is the Phase 5 implementer prompt.
func BuildImplementationPrompt(agent agents.Discovery, idea protocol.IdeaStatus, outputPath string) string {
	finalPath := filepath.Join(idea.Path, "FINAL.md")
	return fmt.Sprintf(`You are %s, the implementer for Parley Deck idea %s (Phase 5).

Rules:
- Implement strictly according to %s. Read it first.
- Before multi-file changes or changes outside parley-deck/, write your plan/checklist into the file below.
- Create exactly this file: %s
- Record unavoidable deviations from FINAL.md in that file.
- Do not edit other agents' protocol files.
- The first line of the file must be exactly "---".
- Return a short confirmation with branch, files changed, checks run.

Required file shape:
---
idea: %s
status: implemented
implementer: %s
started: %s
completed: %s
branch: <repo-path>#<branch-name>
head-commit: <sha-or-short-sha>
---

## Summary of work
## Implementation plan / checklist
## Deviations from FINAL.md
## Notes for reviewers
`, agent.ID, idea.Slug, finalPath, outputPath, idea.Slug, agent.ID, time.Now().Format("2006-01-02"), time.Now().Format("2006-01-02"))
}

// BuildReviewPrompt is the Phase 6 reviewer prompt for review round N.
func BuildReviewPrompt(agent agents.Discovery, idea protocol.IdeaStatus, round int, outputPath, context string) string {
	return fmt.Sprintf(`You are %s, a reviewer for Parley Deck idea %s (Phase 6, review round %d).

Rules:
- Review the implementation against FINAL.md and IMPLEMENTATION.md (below).
- Create exactly this review file: %s
- Do not edit implementation files or other agents' files.
- Use ONLY these severity tags: CRITICAL, MAJOR, MINOR, NIT.
- Each finding states what is wrong, why it matters, and the concrete fix.
- The first line of the file must be exactly "---".
- Return only a short confirmation with the path written.

Required file shape:
---
agent: %s
idea: %s
review-round: %d
date: %s
---

## Summary
## Findings
### [CRITICAL] <title>
### [MAJOR] <title>
### [MINOR] <title>
### [NIT] <title>
## Open questions

Context (FINAL.md, IMPLEMENTATION.md, prior review rounds):
%s
`, agent.ID, idea.Slug, round, outputPath, agent.ID, idea.Slug, round, time.Now().Format("2006-01-02"), context)
}

// gatherReviewContext concatenates FINAL.md, IMPLEMENTATION.md, and any prior
// review-round artifacts to seed a code-review round.
func gatherReviewContext(ideaPath string, round int) (string, error) {
	var b strings.Builder
	for _, name := range []string{"FINAL.md", "IMPLEMENTATION.md"} {
		data, err := os.ReadFile(filepath.Join(ideaPath, name))
		if err == nil {
			fmt.Fprintf(&b, "\n===== %s =====\n%s\n", name, string(data))
		}
	}
	for r := 1; r < round; r++ {
		dir := filepath.Join(ideaPath, "review", roundLabel(r))
		entries, err := os.ReadDir(dir)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return "", err
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") || e.Name() == "_index.md" {
				continue
			}
			data, err := os.ReadFile(filepath.Join(dir, e.Name()))
			if err != nil {
				return "", err
			}
			fmt.Fprintf(&b, "\n===== review/%s/%s =====\n%s\n", roundLabel(r), e.Name(), string(data))
		}
	}
	return b.String(), nil
}

// validateArtifactForPhase dispatches artifact validation by phase.
func validateArtifactForPhase(opts Options, outputPath, agentID string) error {
	switch opts.Phase {
	case "implementation":
		return ValidateImplementationArtifact(outputPath, opts.Idea.Slug)
	case "review":
		return ValidateReviewArtifact(outputPath, agentID, opts.Idea.Slug, roundNumber(opts))
	}
	return ValidateRoundArtifact(outputPath, agentID, opts.Idea.Slug, roundNumber(opts))
}

// ValidateImplementationArtifact checks an IMPLEMENTATION.md.
func ValidateImplementationArtifact(path, ideaSlug string) error {
	meta, err := protocol.ReadFrontmatter(path)
	if err != nil {
		return err
	}
	if got := strings.Trim(strings.TrimSpace(meta["idea"]), `"'`); got != ideaSlug {
		return fmt.Errorf("%s frontmatter idea=%q, want %q", path, got, ideaSlug)
	}
	if strings.TrimSpace(meta["status"]) == "" {
		return fmt.Errorf("%s frontmatter missing status", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if !strings.Contains(string(data), "## Summary of work") {
		return fmt.Errorf("%s missing '## Summary of work'", path)
	}
	return nil
}

// ValidateReviewArtifact checks a review/round-NN/<agent>.md.
func ValidateReviewArtifact(path, agentID, ideaSlug string, round int) error {
	meta, err := protocol.ReadFrontmatter(path)
	if err != nil {
		return err
	}
	for key, want := range map[string]string{
		"agent":        agentID,
		"idea":         ideaSlug,
		"review-round": strconv.Itoa(round),
	} {
		if got := strings.Trim(strings.TrimSpace(meta[key]), `"'`); got != want {
			return fmt.Errorf("%s frontmatter %s=%q, want %q", path, key, got, want)
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if !strings.Contains(string(data), "## Findings") {
		return fmt.Errorf("%s missing '## Findings'", path)
	}
	return nil
}
