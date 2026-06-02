// Package pipeline implements COOPERATION.md §12: composing the unchanged
// Phase 0-8 cooperation engine into ordered pipeline blocks declared by a
// pipeline.yaml manifest. This file covers the manifest schema, parsing, and
// validation. v1 is linear-only: a non-linear edge graph is a validation error.
package pipeline

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// SchemaVersion is the only manifest schema version a v1 driver accepts.
const SchemaVersion = 1

// BlockKind enumerates the four §12.1 block kinds. Each block is one invocation
// of the existing cooperation engine.
type BlockKind string

const (
	KindDeliberation  BlockKind = "deliberation"
	KindImplementation BlockKind = "implementation"
	KindAction        BlockKind = "action"
	KindWatcher       BlockKind = "watcher"
)

// Risk is the §12.2 risk level. "production" mutations are non-bypassable.
type Risk string

const (
	RiskLow        Risk = "low"
	RiskNormal     Risk = "normal"
	RiskHigh       Risk = "high"
	RiskProduction Risk = "production"
)

// Autonomy is the §12.8 per-pipeline autonomy setting.
type Autonomy string

const (
	AutonomySupervised Autonomy = "supervised"
	AutonomyAutoLeft   Autonomy = "auto-left"
)

// Block is one entry in the manifest's ordered block list (§12.1, §12.3).
type Block struct {
	ID                  string    `yaml:"id" json:"id"`
	Kind                BlockKind `yaml:"kind" json:"kind"`
	Stage               string    `yaml:"stage,omitempty" json:"stage,omitempty"`
	RoleLens            string    `yaml:"role_lens,omitempty" json:"role_lens,omitempty"`
	InputArtifacts      []string  `yaml:"input_artifacts,omitempty" json:"input_artifacts,omitempty"`
	OutputArtifact      string    `yaml:"output_artifact,omitempty" json:"output_artifact,omitempty"`
	Risk                Risk      `yaml:"risk,omitempty" json:"risk,omitempty"`
	ProviderCapabilities []string `yaml:"provider_capabilities,omitempty" json:"provider_capabilities,omitempty"`
	GatePolicy          string    `yaml:"gate_policy,omitempty" json:"gate_policy,omitempty"`
}

// Edge is a reserved-for-future DAG edge (§12.3). A v1 driver accepts only a
// single linear chain that matches block order; anything else is rejected.
type Edge struct {
	From   string `yaml:"from" json:"from"`
	To     string `yaml:"to" json:"to"`
	GateID string `yaml:"gate_id,omitempty" json:"gate_id,omitempty"`
}

// Manifest is the parsed pipeline.yaml (§12.3).
type Manifest struct {
	SchemaVersion int      `yaml:"schema_version" json:"schema_version"`
	IdeaSlug      string   `yaml:"idea_slug" json:"idea_slug"`
	Autonomy      Autonomy `yaml:"autonomy,omitempty" json:"autonomy,omitempty"`
	Transport     string   `yaml:"transport" json:"transport"`
	Participants  []string `yaml:"participants" json:"participants"`
	Blocks        []Block  `yaml:"blocks" json:"blocks"`
	Edges         []Edge   `yaml:"edges,omitempty" json:"edges,omitempty"`
}

var validKinds = map[BlockKind]bool{
	KindDeliberation: true, KindImplementation: true, KindAction: true, KindWatcher: true,
}

var validRisks = map[Risk]bool{
	RiskLow: true, RiskNormal: true, RiskHigh: true, RiskProduction: true,
}

var validTransports = map[string]bool{
	"local-dir": true, "github-pr": true, "gitlab-mr": true,
}

// Parse decodes manifest bytes and validates them. The returned manifest is
// safe for a v1 driver to execute.
func Parse(data []byte) (Manifest, error) {
	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return Manifest{}, fmt.Errorf("parse pipeline manifest: %w", err)
	}
	if m.Autonomy == "" {
		m.Autonomy = AutonomySupervised
	}
	if err := m.Validate(); err != nil {
		return Manifest{}, err
	}
	return m, nil
}

// ParseFile reads and parses a pipeline.yaml manifest from disk.
func ParseFile(path string) (Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, fmt.Errorf("read pipeline manifest: %w", err)
	}
	return Parse(data)
}

// Validate enforces the v1 manifest rules, including the linear-only edge graph.
func (m Manifest) Validate() error {
	if m.SchemaVersion != SchemaVersion {
		return fmt.Errorf("unsupported schema_version %d (want %d)", m.SchemaVersion, SchemaVersion)
	}
	if m.IdeaSlug == "" {
		return fmt.Errorf("idea_slug is required")
	}
	if !validTransports[m.Transport] {
		return fmt.Errorf("invalid transport %q (want local-dir|github-pr|gitlab-mr)", m.Transport)
	}
	if m.Autonomy != AutonomySupervised && m.Autonomy != AutonomyAutoLeft {
		return fmt.Errorf("invalid autonomy %q (want supervised|auto-left)", m.Autonomy)
	}
	if len(m.Blocks) == 0 {
		return fmt.Errorf("at least one block is required")
	}

	seen := make(map[string]bool, len(m.Blocks))
	for i, b := range m.Blocks {
		if b.ID == "" {
			return fmt.Errorf("block %d: id is required", i)
		}
		if seen[b.ID] {
			return fmt.Errorf("duplicate block id %q", b.ID)
		}
		seen[b.ID] = true
		if !validKinds[b.Kind] {
			return fmt.Errorf("block %q: invalid kind %q", b.ID, b.Kind)
		}
		if b.Risk != "" && !validRisks[b.Risk] {
			return fmt.Errorf("block %q: invalid risk %q", b.ID, b.Risk)
		}
	}

	return m.validateLinear()
}

// validateLinear enforces §12.3 / §12.9: edges, if present, must describe the
// single linear chain that matches block order. Any branch, join, cycle,
// skipped block, or unknown endpoint is rejected.
func (m Manifest) validateLinear() error {
	if len(m.Edges) == 0 {
		return nil // linear order is implied by blocks[] order
	}
	if len(m.Edges) != len(m.Blocks)-1 {
		return fmt.Errorf("non-linear manifest: expected %d edges for %d blocks, got %d", len(m.Blocks)-1, len(m.Blocks), len(m.Edges))
	}
	for i, e := range m.Edges {
		wantFrom := m.Blocks[i].ID
		wantTo := m.Blocks[i+1].ID
		if e.From != wantFrom || e.To != wantTo {
			return fmt.Errorf("non-linear manifest: edge %d must be %s->%s, got %s->%s (v1 is linear-only)", i, wantFrom, wantTo, e.From, e.To)
		}
	}
	return nil
}
