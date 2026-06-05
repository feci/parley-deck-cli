package runaction

import (
	"fmt"
	"strings"
)

const (
	KindAnswerQuestion  = "answer-question"
	KindRetryAgent      = "retry-agent"
	KindOpenNextRound   = "open-next-round"
	KindDraftConsensus  = "draft-consensus"
	KindRequestSignoffs = "request-signoffs"
	KindFinalize        = "finalize"
	KindInspect         = "inspect"

	RiskLow    = "low"
	RiskNormal = "normal"
	RiskHigh   = "high"
)

type NextAction struct {
	ID           string `json:"id"`
	Kind         string `json:"kind"`
	RunID        string `json:"run_id,omitempty"`
	IdeaSlug     string `json:"idea_slug,omitempty"`
	Phase        string `json:"phase,omitempty"`
	Round        string `json:"round,omitempty"`
	AgentID      string `json:"agent_id,omitempty"`
	ArtifactPath string `json:"artifact_path,omitempty"`
	Risk         string `json:"risk,omitempty"`
	RequiresYes  bool   `json:"requires_yes,omitempty"`
	Summary      string `json:"summary,omitempty"`
}

func Command(action NextAction, fallbackRunID, fallbackIdeaSlug string) string {
	runID := valueOr(action.RunID, fallbackRunID)
	idea := valueOr(action.IdeaSlug, fallbackIdeaSlug)
	switch action.Kind {
	case KindAnswerQuestion:
		questionID := strings.TrimPrefix(action.ID, "answer-question.")
		if runID == "" || questionID == "" || questionID == action.ID {
			return ""
		}
		return fmt.Sprintf("parley answer %s %s <answer>", runID, questionID)
	case KindOpenNextRound:
		if idea == "" {
			return ""
		}
		// Visibility only: the next cross-review round is opened automatically by
		// internal/driver under `parley run --auto` (local-dir). The surfaced
		// command re-runs the task to advance the idea.
		return fmt.Sprintf("parley run --auto --dir . \"continue %s\"", idea)
	case KindDraftConsensus:
		if idea == "" {
			return ""
		}
		if round := RoundNumber(action.Round); round != "" {
			return fmt.Sprintf("parley consensus draft --round %s %s", round, idea)
		}
		return fmt.Sprintf("parley consensus draft %s", idea)
	case KindRequestSignoffs:
		if idea == "" {
			return ""
		}
		return fmt.Sprintf("parley consensus request-signoffs --yes %s", idea)
	case KindFinalize:
		if idea == "" {
			return ""
		}
		return fmt.Sprintf("parley consensus finalize %s", idea)
	case KindInspect:
		if runID != "" {
			return fmt.Sprintf("parley status --run %s", runID)
		}
		if idea != "" {
			return fmt.Sprintf("parley status --idea %s", idea)
		}
	}
	return ""
}

func RoundNumber(round string) string {
	round = strings.TrimSpace(round)
	if !strings.HasPrefix(round, "round-") {
		return ""
	}
	n := strings.TrimLeft(strings.TrimPrefix(round, "round-"), "0")
	if n == "" {
		return ""
	}
	return n
}

func valueOr(value, fallback string) string {
	if strings.TrimSpace(value) != "" {
		return value
	}
	return fallback
}
