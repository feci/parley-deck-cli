package pipeline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// GateStatus is the lifecycle of a §12.8 boundary gate.
type GateStatus string

const (
	GateOpen     GateStatus = "open"
	GateApproved GateStatus = "approved"
	GateRejected GateStatus = "rejected"
)

// Gate is a typed promote-to-next-block approval at
// pipelines/<slug>/gates/<edge-id>.gate.json. It reuses the HITL question/risk
// model: a human (or, for low-risk non-production gates under auto-left
// autonomy, the policy evaluator) resolves it before the driver may advance.
type Gate struct {
	ID            string     `json:"id"`
	PipelineSlug  string     `json:"pipeline_slug"`
	Edge          string     `json:"edge"`
	FromBlock     string     `json:"from_block"`
	ToBlock       string     `json:"to_block"`
	Risk          Risk       `json:"risk"`
	Status        GateStatus `json:"status"`
	Prompt        string     `json:"prompt"`
	Policy        string     `json:"policy,omitempty"`
	ApprovedBy    string     `json:"approved_by,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	AnsweredAt    time.Time  `json:"answered_at,omitempty"`
}

// EdgeID is the stable identifier for the boundary between two blocks.
func EdgeID(from, to string) string { return from + "->" + to }

// AutoApprove reports whether a boundary gate may be resolved without a human
// under the given autonomy, with no decider configured.
func AutoApprove(autonomy Autonomy, risk Risk) bool {
	return AutoApproveWithDecider(autonomy, risk, false)
}

// AutoApproveWithDecider is the single central policy evaluator (§12.8/§12.11):
// production-risk mutations are NEVER auto-approvable.
//   - auto-left autonomy auto-resolves low/normal, non-production boundaries.
//   - a configured decider agent auto-resolves ONLY low-risk, non-production
//     boundaries (never normal/high/production) — strictly narrower than
//     auto-left, per FINAL.md item 6b.
// Supervised + no decider = block-and-wait (default).
func AutoApproveWithDecider(autonomy Autonomy, risk Risk, hasDecider bool) bool {
	if risk == RiskProduction {
		return false
	}
	if autonomy == AutonomyAutoLeft && (risk == RiskLow || risk == "" || risk == RiskNormal) {
		return true
	}
	if hasDecider && risk == RiskLow {
		return true
	}
	return false
}

// NewGate builds an open boundary gate for the edge from->to.
func NewGate(slug, from, to string, risk Risk, policy string, now time.Time) Gate {
	return Gate{
		ID:           EdgeID(from, to),
		PipelineSlug: slug,
		Edge:         EdgeID(from, to),
		FromBlock:    from,
		ToBlock:      to,
		Risk:         risk,
		Status:       GateOpen,
		Prompt:       fmt.Sprintf("Approve advancing pipeline %q from block %q to %q (risk=%s)?", slug, from, to, riskOrNormal(risk)),
		Policy:       policy,
		CreatedAt:    now,
	}
}

func riskOrNormal(r Risk) Risk {
	if r == "" {
		return RiskNormal
	}
	return r
}

// Resolve marks the gate approved or rejected.
func (g *Gate) Resolve(approve bool, by string, now time.Time) {
	if approve {
		g.Status = GateApproved
	} else {
		g.Status = GateRejected
	}
	g.ApprovedBy = by
	g.AnsweredAt = now
}

// GatePath returns the on-disk location of a gate file.
func GatePath(deckDir, slug, edgeID string) string {
	return filepath.Join(PipelineDir(deckDir, slug), "gates", edgeID+".gate.json")
}

// SaveGate writes a gate atomically.
func SaveGate(deckDir string, g Gate) error {
	path := GatePath(deckDir, g.PipelineSlug, g.Edge)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create gates dir: %w", err)
	}
	data, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return fmt.Errorf("encode gate: %w", err)
	}
	data = append(data, '\n')
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("write gate: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("commit gate: %w", err)
	}
	return nil
}

// LoadGate reads a gate file; the bool reports whether it exists.
func LoadGate(deckDir, slug, edgeID string) (Gate, bool, error) {
	path := GatePath(deckDir, slug, edgeID)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Gate{}, false, nil
	}
	if err != nil {
		return Gate{}, false, fmt.Errorf("read gate: %w", err)
	}
	var g Gate
	if err := json.Unmarshal(data, &g); err != nil {
		return Gate{}, false, fmt.Errorf("parse gate: %w", err)
	}
	return g, true, nil
}
