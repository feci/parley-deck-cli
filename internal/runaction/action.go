package runaction

const (
	KindAnswerQuestion  = "answer-question"
	KindRetryAgent      = "retry-agent"
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
