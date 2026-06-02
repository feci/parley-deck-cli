package pipeline

import (
	"path/filepath"
	"testing"
	"time"
)

func sampleManifest(t *testing.T) Manifest {
	t.Helper()
	m, err := Parse([]byte(validLinearYAML))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	return m
}

func TestNewPipelineRunPositionsAtFirstBlock(t *testing.T) {
	m := sampleManifest(t)
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	r := NewPipelineRun(m, now)

	if r.CurrentBlock != "business-spec" {
		t.Fatalf("current_block = %q, want business-spec", r.CurrentBlock)
	}
	if r.Status != StatusPending {
		t.Fatalf("status = %q, want pending", r.Status)
	}
	if r.PipelineSlug != "demo-pipeline" {
		t.Fatalf("pipeline_slug = %q", r.PipelineSlug)
	}
	if r.SchemaVersion != SchemaVersion {
		t.Fatalf("schema_version = %d", r.SchemaVersion)
	}
}

func TestPipelineRunRoundTrip(t *testing.T) {
	m := sampleManifest(t)
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	r := NewPipelineRun(m, now)
	r.Status = StatusBlockedOnGate
	r.CompletedBlocks = []string{"business-spec"}
	r.CurrentBlock = "technical-spec"
	r.PendingGate = "business-spec->technical-spec"
	r.EffectKeys = []string{"abc123"}

	path := filepath.Join(t.TempDir(), RunFileName)
	if err := r.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := LoadRun(path)
	if err != nil {
		t.Fatalf("LoadRun: %v", err)
	}

	if got.Status != StatusBlockedOnGate {
		t.Fatalf("status = %q", got.Status)
	}
	if got.CurrentBlock != "technical-spec" {
		t.Fatalf("current_block = %q", got.CurrentBlock)
	}
	if got.PendingGate != "business-spec->technical-spec" {
		t.Fatalf("pending_gate = %q", got.PendingGate)
	}
	if len(got.CompletedBlocks) != 1 || got.CompletedBlocks[0] != "business-spec" {
		t.Fatalf("completed_blocks = %v", got.CompletedBlocks)
	}
	if len(got.EffectKeys) != 1 || got.EffectKeys[0] != "abc123" {
		t.Fatalf("effect_keys = %v", got.EffectKeys)
	}
	if !got.UpdatedAt.Equal(now) {
		t.Fatalf("updated_at = %v, want %v", got.UpdatedAt, now)
	}
}
