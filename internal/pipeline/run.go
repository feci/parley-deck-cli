package pipeline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// RunStatus is the durable pipeline cursor status (§12.6).
type RunStatus string

const (
	StatusPending       RunStatus = "pending"
	StatusRunning       RunStatus = "running"
	StatusBlockedOnGate RunStatus = "blocked_on_gate"
	StatusCompleted     RunStatus = "completed"
	StatusFailed        RunStatus = "failed"
)

// PipelineRun is the durable cursor/index at pipelines/<slug>/pipeline-run.json
// (§12.6). It is a cursor only: per-effect records live in their own files
// under effects/ and are referenced here by key. The driver is the only writer.
type PipelineRun struct {
	SchemaVersion   int       `json:"schema_version"`
	PipelineSlug    string    `json:"pipeline_slug"`
	Status          RunStatus `json:"status"`
	CurrentBlock    string    `json:"current_block,omitempty"`
	CompletedBlocks []string  `json:"completed_blocks"`
	// ReadyBlocks / ActiveBlocks support multi-active DAG execution. They are
	// additive: a linear/single-active run leaves them empty and uses
	// current_block, so older cursors load unchanged.
	ReadyBlocks  []string `json:"ready_blocks,omitempty"`
	ActiveBlocks []string `json:"active_blocks,omitempty"`
	PendingGate  string   `json:"pending_gate,omitempty"`
	EffectKeys      []string  `json:"effect_keys,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// RunFileName is the cursor filename within pipelines/<slug>/.
const RunFileName = "pipeline-run.json"

// NewPipelineRun creates a fresh cursor positioned at the manifest's first
// block. The caller supplies the timestamp so callers (and tests) stay
// deterministic.
func NewPipelineRun(m Manifest, now time.Time) PipelineRun {
	current := ""
	if len(m.Blocks) > 0 {
		current = m.Blocks[0].ID
	}
	return PipelineRun{
		SchemaVersion:   SchemaVersion,
		PipelineSlug:    m.IdeaSlug,
		Status:          StatusPending,
		CurrentBlock:    current,
		CompletedBlocks: []string{},
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// LoadRun reads a pipeline-run.json cursor from disk.
func LoadRun(path string) (PipelineRun, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return PipelineRun{}, fmt.Errorf("read pipeline run: %w", err)
	}
	var r PipelineRun
	if err := json.Unmarshal(data, &r); err != nil {
		return PipelineRun{}, fmt.Errorf("parse pipeline run: %w", err)
	}
	if r.CompletedBlocks == nil {
		r.CompletedBlocks = []string{}
	}
	return r, nil
}

// Save writes the cursor atomically (temp file + rename) so a crash mid-write
// cannot corrupt the durable state the driver resumes from.
func (r PipelineRun) Save(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create pipeline dir: %w", err)
	}
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("encode pipeline run: %w", err)
	}
	data = append(data, '\n')
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("write pipeline run: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("commit pipeline run: %w", err)
	}
	return nil
}
