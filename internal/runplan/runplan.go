package runplan

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/hitl"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runaction"
)

const (
	KindAnswerQuestion  = runaction.KindAnswerQuestion
	KindRetryAgent      = runaction.KindRetryAgent
	KindOpenNextRound   = runaction.KindOpenNextRound
	KindDraftConsensus  = runaction.KindDraftConsensus
	KindRequestSignoffs = runaction.KindRequestSignoffs
	KindFinalize        = runaction.KindFinalize
	KindInspect         = runaction.KindInspect

	RiskLow    = runaction.RiskLow
	RiskNormal = runaction.RiskNormal
	RiskHigh   = runaction.RiskHigh
)

type NextAction = runaction.NextAction

type Input struct {
	RunID        string
	IdeaSlug     string
	Participants []string
	Terminal     bool
	Outcome      string
	Questions    []hitl.Question
	Agents       []AgentState
	RoundStatus  string
	CurrentRound string
}

type AgentState struct {
	ID           string
	State        string
	ArtifactPath string
	Error        string
	Reason       string
}

func Plan(root string, input Input) []NextAction {
	var actions []NextAction
	for _, question := range input.Questions {
		if question.Status != hitl.StatusOpen {
			continue
		}
		actions = append(actions, NextAction{
			ID:          "answer-question." + question.ID,
			Kind:        KindAnswerQuestion,
			RunID:       input.RunID,
			IdeaSlug:    input.IdeaSlug,
			Phase:       "hitl",
			AgentID:     question.Agent,
			Risk:        riskFromQuestion(question.Risk),
			RequiresYes: false,
			Summary:     fmt.Sprintf("Answer HITL question %s for %s", question.ID, valueOr(question.Agent, "run")),
		})
	}

	if input.IdeaSlug == "" || input.IdeaSlug == "unknown" {
		return appendInspectIfEmpty(actions, input, "Inspect run state; idea slug is unknown")
	}

	ideaDir := filepath.Join(root, protocol.DeckDir, "ideas", input.IdeaSlug)
	participants := participantOrder(input)
	if len(participants) == 0 {
		return appendInspectIfEmpty(actions, input, "Inspect run state; no participants are known")
	}

	round := currentRound(input)
	missingRoundArtifact := false
	for _, participant := range participants {
		artifactRel := filepath.ToSlash(filepath.Join(round, participant+".md"))
		if artifactExists(ideaDir, artifactRel) {
			continue
		}
		missingRoundArtifact = true
		agent := agentByID(input.Agents, participant)
		if shouldRetryAgent(input, agent) {
			actions = append(actions, NextAction{
				ID:           "retry-agent." + round + "." + participant,
				Kind:         KindRetryAgent,
				RunID:        input.RunID,
				IdeaSlug:     input.IdeaSlug,
				Phase:        "round",
				Round:        round,
				AgentID:      participant,
				ArtifactPath: artifactRel,
				Risk:         RiskNormal,
				RequiresYes:  true,
				Summary:      fmt.Sprintf("Retry missing or failed %s artifact for %s", round, participant),
			})
		}
	}

	if missingRoundArtifact {
		return appendInspectIfEmpty(actions, input, "Inspect incomplete round state")
	}

	// Round complete. If the cross-review policy (00-prompt cross_review_rounds,
	// default 1) wants another round and no consensus exists yet, the next protocol
	// step is to open that round — executed automatically by internal/driver under
	// `parley run --auto` (local-dir), surfaced here so `parley continue` no longer
	// jumps a completed independent round straight to consensus.
	if next := nextCrossReviewRound(ideaDir, round); next != "" {
		if _, statusErr := consensus.Status(root, input.IdeaSlug, false); errors.Is(statusErr, os.ErrNotExist) {
			actions = append(actions, NextAction{
				ID:          "open-next-round." + input.IdeaSlug,
				Kind:        KindOpenNextRound,
				RunID:       input.RunID,
				IdeaSlug:    input.IdeaSlug,
				Phase:       "round",
				Round:       next,
				Risk:        RiskNormal,
				RequiresYes: true,
				Summary:     "Open " + next + " (cross-review) before drafting consensus",
			})
			return actions
		}
	}

	summary, err := consensus.Status(root, input.IdeaSlug, false)
	if errors.Is(err, os.ErrNotExist) {
		actions = append(actions, NextAction{
			ID:          "draft-consensus." + input.IdeaSlug,
			Kind:        KindDraftConsensus,
			RunID:       input.RunID,
			IdeaSlug:    input.IdeaSlug,
			Phase:       "consensus",
			Round:       round,
			Risk:        RiskNormal,
			RequiresYes: true,
			Summary:     "Draft consensus from completed round artifacts",
		})
		return actions
	}
	if err != nil {
		return appendInspectIfEmpty(actions, input, "Inspect consensus state; consensus status could not be read")
	}

	switch summary.Triage {
	case consensus.TriageMalformed:
		actions = append(actions, NextAction{
			ID:          "inspect.consensus-malformed",
			Kind:        KindInspect,
			RunID:       input.RunID,
			IdeaSlug:    input.IdeaSlug,
			Phase:       "consensus",
			Risk:        RiskHigh,
			RequiresYes: false,
			Summary:     "Inspect malformed consensus before continuing",
		})
	case consensus.TriagePartial:
		actions = append(actions, NextAction{
			ID:          "request-signoffs." + input.IdeaSlug,
			Kind:        KindRequestSignoffs,
			RunID:       input.RunID,
			IdeaSlug:    input.IdeaSlug,
			Phase:       "consensus",
			Risk:        RiskNormal,
			RequiresYes: true,
			Summary:     "Request missing consensus signoffs: " + strings.Join(summary.Missing, ","),
		})
	case consensus.TriageBlocked:
		actions = append(actions, NextAction{
			ID:          "inspect.consensus-blocked",
			Kind:        KindInspect,
			RunID:       input.RunID,
			IdeaSlug:    input.IdeaSlug,
			Phase:       "consensus",
			Risk:        RiskHigh,
			RequiresYes: false,
			Summary:     "Inspect blocked consensus before opening another round",
		})
	case consensus.TriageReady, consensus.TriageReserved:
		if !fileExists(filepath.Join(ideaDir, "FINAL.md")) {
			actions = append(actions, NextAction{
				ID:          "finalize." + input.IdeaSlug,
				Kind:        KindFinalize,
				RunID:       input.RunID,
				IdeaSlug:    input.IdeaSlug,
				Phase:       "finalization",
				Risk:        RiskNormal,
				RequiresYes: true,
				Summary:     "Finalize accepted consensus",
			})
		}
	}

	return appendInspectIfEmpty(actions, input, "No recoverable action; inspect artifacts and logs")
}

func shouldRetryAgent(input Input, agent AgentState) bool {
	if agent.ID == "" {
		return input.Terminal || input.RoundStatus == "incomplete"
	}
	switch agent.State {
	case "failed", "skipped":
		return true
	case "pending":
		return input.Terminal || input.RoundStatus == "incomplete"
	default:
		return false
	}
}

func currentRound(input Input) string {
	round := strings.TrimSpace(input.CurrentRound)
	if round == "" {
		return "round-01"
	}
	return round
}

// nextCrossReviewRound returns the next round label ("round-NN") when the
// cross-review policy in 00-prompt.md (cross_review_rounds, default 1) wants
// another cross-review round after the given completed round, or "" when the
// budget is exhausted (→ proceed to consensus). round-01 is the independent
// round; the policy adds N cross-review rounds, so the last deliberation round is
// 1+cross_review_rounds.
func nextCrossReviewRound(ideaDir, completedRound string) string {
	n := roundOrdinal(completedRound)
	if n < 1 || n >= 1+readCrossReviewRounds(ideaDir) {
		return ""
	}
	return fmt.Sprintf("round-%02d", n+1)
}

func roundOrdinal(label string) int {
	num := runaction.RoundNumber(label)
	if num == "" {
		return 0
	}
	n, err := strconv.Atoi(num)
	if err != nil {
		return 0
	}
	return n
}

// readCrossReviewRounds reads cross_review_rounds from the idea 00-prompt.md
// frontmatter; defaults to 1. N=0 is an explicit straight-to-consensus bypass.
func readCrossReviewRounds(ideaDir string) int {
	const def = 1
	data, err := os.ReadFile(filepath.Join(ideaDir, "00-prompt.md"))
	if err != nil {
		return def
	}
	inFrontmatter := false
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			}
			break
		}
		if inFrontmatter && strings.HasPrefix(trimmed, "cross_review_rounds:") {
			value := strings.TrimSpace(strings.TrimPrefix(trimmed, "cross_review_rounds:"))
			if n, err := strconv.Atoi(value); err == nil && n >= 0 {
				return n
			}
			return def
		}
	}
	return def
}

func participantOrder(input Input) []string {
	seen := map[string]bool{}
	var participants []string
	for _, participant := range input.Participants {
		participant = strings.TrimSpace(participant)
		if participant == "" || seen[participant] {
			continue
		}
		seen[participant] = true
		participants = append(participants, participant)
	}
	for _, agent := range input.Agents {
		id := strings.TrimSpace(agent.ID)
		if id == "" || seen[id] {
			continue
		}
		seen[id] = true
		participants = append(participants, id)
	}
	return participants
}

func agentByID(agents []AgentState, id string) AgentState {
	for _, agent := range agents {
		if agent.ID == id {
			return agent
		}
	}
	return AgentState{ID: id, State: "pending"}
}

func artifactExists(ideaDir, artifactRel string) bool {
	return fileExists(filepath.Join(ideaDir, filepath.FromSlash(artifactRel)))
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir() && info.Size() > 0
}

func riskFromQuestion(risk string) string {
	switch strings.TrimSpace(strings.ToLower(risk)) {
	case hitl.RiskLow:
		return RiskLow
	case hitl.RiskHigh:
		return RiskHigh
	default:
		return RiskNormal
	}
}

func appendInspectIfEmpty(actions []NextAction, input Input, summary string) []NextAction {
	if len(actions) > 0 {
		return actions
	}
	return append(actions, NextAction{
		ID:          "inspect." + valueOr(input.RunID, input.IdeaSlug),
		Kind:        KindInspect,
		RunID:       input.RunID,
		IdeaSlug:    input.IdeaSlug,
		Phase:       "inspection",
		Risk:        RiskLow,
		RequiresYes: false,
		Summary:     summary,
	})
}

func valueOr(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
