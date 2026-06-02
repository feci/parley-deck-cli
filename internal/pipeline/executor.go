package pipeline

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// PipelineDir is <deck>/pipelines/<slug> — the home of the manifest, cursor,
// gates, and effects ledger for one pipeline. deckDir is the parley-deck/ dir.
func PipelineDir(deckDir, slug string) string {
	return filepath.Join(deckDir, "pipelines", slug)
}

// ManifestPath is the pipeline.yaml location.
func ManifestPath(deckDir, slug string) string {
	return filepath.Join(PipelineDir(deckDir, slug), "pipeline.yaml")
}

// RunPath is the pipeline-run.json cursor location.
func RunPath(deckDir, slug string) string {
	return filepath.Join(PipelineDir(deckDir, slug), RunFileName)
}

// BlockWorkspace is the per-block engine workspace, an ordinary idea dir so the
// unchanged Phase 0-8 engine runs in it without colliding with sibling blocks.
func BlockWorkspace(deckDir, slug, blockID string) string {
	return filepath.Join(deckDir, "ideas", slug+"__"+blockID)
}

// Action is what the driver determined should happen next for a pipeline.
type Action string

const (
	ActionRunBlock   Action = "run-block"    // current block has not finished its engine round
	ActionAwaitGate  Action = "await-gate"   // a boundary gate is open and needs resolution
	ActionSeededNext Action = "seeded-next"  // advanced: next block's prompt was seeded
	ActionDone       Action = "done"         // final block complete; pipeline finished
	ActionRejected   Action = "rejected"     // a boundary gate was rejected; pipeline stopped
)

// AdvanceResult is the outcome of one driver step.
type AdvanceResult struct {
	Action  Action
	Block   string // the block the action concerns (current or newly-current)
	Edge    string // set for await-gate
	Detail  string
}

// BlockComplete reports whether a block's engine round has finished and its
// output artifact is finalized. It is injected so the driver stays decoupled
// from the consensus/runner internals and is unit-testable.
type BlockComplete func(b Block) (bool, error)

// Driver walks a manifest, evaluating boundary gates and seeding each block's
// kickoff from the previous block's output. It mutates the cursor in place;
// the caller persists it with PipelineRun.Save.
type Driver struct {
	DeckDir  string
	Manifest Manifest
	Now      time.Time
}

// Advance performs one supervised-first step of the state machine and returns
// what happened. The cursor is updated in place.
func (d Driver) Advance(run *PipelineRun, complete BlockComplete) (AdvanceResult, error) {
	if run.Status == StatusCompleted {
		return AdvanceResult{Action: ActionDone, Block: lastBlockID(d.Manifest)}, nil
	}
	if run.Status == StatusFailed {
		return AdvanceResult{Action: ActionRejected, Detail: "pipeline previously stopped on a rejected gate"}, nil
	}

	if d.Manifest.IsDAG() {
		return d.advanceDAG(run, complete)
	}

	idx := d.blockIndex(run.CurrentBlock)
	if idx < 0 {
		return AdvanceResult{}, fmt.Errorf("current block %q not found in manifest", run.CurrentBlock)
	}
	cur := d.Manifest.Blocks[idx]

	done, err := complete(cur)
	if err != nil {
		return AdvanceResult{}, fmt.Errorf("check block %q complete: %w", cur.ID, err)
	}
	if !done {
		run.Status = StatusRunning
		return AdvanceResult{Action: ActionRunBlock, Block: cur.ID, Detail: "block has not finished its engine round + consensus"}, nil
	}

	// Current block is finalized. Is it the last one?
	if idx == len(d.Manifest.Blocks)-1 {
		run.recordCompleted(cur.ID)
		run.CurrentBlock = ""
		run.PendingGate = ""
		run.Status = StatusCompleted
		return AdvanceResult{Action: ActionDone, Block: cur.ID, Detail: "final block complete"}, nil
	}

	next := d.Manifest.Blocks[idx+1]
	edge := EdgeID(cur.ID, next.ID)

	gate, exists, err := LoadGate(d.DeckDir, run.PipelineSlug, edge)
	if err != nil {
		return AdvanceResult{}, err
	}
	if !exists {
		gate = NewGate(run.PipelineSlug, cur.ID, next.ID, next.Risk, next.GatePolicy, d.Now)
		if AutoApproveWithDecider(d.Manifest.Autonomy, next.Risk, d.Manifest.Decider != "") {
			by := "policy:auto-left"
			if d.Manifest.Decider != "" && d.Manifest.Autonomy != AutonomyAutoLeft {
				by = "decider:" + d.Manifest.Decider
			}
			gate.Resolve(true, by, d.Now)
		}
		if err := SaveGate(d.DeckDir, gate); err != nil {
			return AdvanceResult{}, err
		}
	}

	switch gate.Status {
	case GateRejected:
		run.Status = StatusFailed
		run.PendingGate = edge
		return AdvanceResult{Action: ActionRejected, Block: next.ID, Edge: edge, Detail: "boundary gate rejected"}, nil
	case GateOpen:
		run.Status = StatusBlockedOnGate
		run.PendingGate = edge
		return AdvanceResult{Action: ActionAwaitGate, Block: next.ID, Edge: edge, Detail: gate.Prompt}, nil
	case GateApproved:
		if _, err := d.seedNext(cur, next); err != nil {
			return AdvanceResult{}, err
		}
		run.recordCompleted(cur.ID)
		run.CurrentBlock = next.ID
		run.PendingGate = ""
		run.Status = StatusRunning
		return AdvanceResult{Action: ActionSeededNext, Block: next.ID, Edge: edge, Detail: "seeded next block kickoff"}, nil
	default:
		return AdvanceResult{}, fmt.Errorf("gate %q has unknown status %q", edge, gate.Status)
	}
}

// advanceDAG is the single-active DAG executor (Execution: dag). The next block
// is chosen by TOPOLOGICAL readiness (all inbound-edge sources complete), not by
// blocks[] order, so a dependent block can never run before its prerequisites
// even if listed first. One block is active at a time; full parallel multi-active
// execution is future work. The manifest's first listed block should be a root.
func (d Driver) advanceDAG(run *PipelineRun, complete BlockComplete) (AdvanceResult, error) {
	if cur := run.CurrentBlock; cur != "" {
		b, ok := d.findBlockByID(cur)
		if !ok {
			return AdvanceResult{}, fmt.Errorf("current block %q not found in manifest", cur)
		}
		done, err := complete(b)
		if err != nil {
			return AdvanceResult{}, fmt.Errorf("check block %q complete: %w", cur, err)
		}
		if !done {
			run.Status = StatusRunning
			return AdvanceResult{Action: ActionRunBlock, Block: cur, Detail: "block has not finished its engine round + consensus"}, nil
		}
		run.recordCompleted(cur)
	}
	if d.Manifest.AllBlocksComplete(run.CompletedBlocks) {
		run.CurrentBlock = ""
		run.PendingGate = ""
		run.Status = StatusCompleted
		return AdvanceResult{Action: ActionDone, Block: lastBlockID(d.Manifest), Detail: "all blocks complete"}, nil
	}

	// Topologically ready blocks (inbound sources complete); gates applied below.
	candidates := d.Manifest.ReadyBlocks(run.CompletedBlocks, func(string) bool { return true })
	if len(candidates) == 0 {
		return AdvanceResult{}, fmt.Errorf("dag deadlock: no ready block but pipeline incomplete (completed: %v)", run.CompletedBlocks)
	}
	nextID := candidates[0]
	next, _ := d.findBlockByID(nextID)
	from := run.CurrentBlock
	if from == "" {
		from = "start"
	}
	edge := EdgeID(from, nextID)

	gate, exists, err := LoadGate(d.DeckDir, run.PipelineSlug, edge)
	if err != nil {
		return AdvanceResult{}, err
	}
	if !exists {
		gate = NewGate(run.PipelineSlug, from, nextID, next.Risk, next.GatePolicy, d.Now)
		if AutoApproveWithDecider(d.Manifest.Autonomy, next.Risk, d.Manifest.Decider != "") {
			by := "policy:auto-left"
			if d.Manifest.Decider != "" && d.Manifest.Autonomy != AutonomyAutoLeft {
				by = "decider:" + d.Manifest.Decider
			}
			gate.Resolve(true, by, d.Now)
		}
		if err := SaveGate(d.DeckDir, gate); err != nil {
			return AdvanceResult{}, err
		}
	}
	switch gate.Status {
	case GateRejected:
		run.Status = StatusFailed
		run.PendingGate = edge
		return AdvanceResult{Action: ActionRejected, Block: nextID, Edge: edge, Detail: "boundary gate rejected"}, nil
	case GateOpen:
		run.Status = StatusBlockedOnGate
		run.PendingGate = edge
		return AdvanceResult{Action: ActionAwaitGate, Block: nextID, Edge: edge, Detail: gate.Prompt}, nil
	case GateApproved:
		var fromPtr *Block
		if run.CurrentBlock != "" {
			if fb, ok := d.findBlockByID(run.CurrentBlock); ok {
				fromPtr = &fb
			}
		}
		if _, err := SeedBlockPrompt(d.DeckDir, d.Manifest, fromPtr, next, d.Now); err != nil {
			return AdvanceResult{}, err
		}
		run.CurrentBlock = nextID
		run.PendingGate = ""
		run.Status = StatusRunning
		return AdvanceResult{Action: ActionSeededNext, Block: nextID, Edge: edge, Detail: "seeded next ready block"}, nil
	default:
		return AdvanceResult{}, fmt.Errorf("gate %q has unknown status %q", edge, gate.Status)
	}
}

func (d Driver) findBlockByID(id string) (Block, bool) {
	for _, b := range d.Manifest.Blocks {
		if b.ID == id {
			return b, true
		}
	}
	return Block{}, false
}

// seedNext writes the next block's driver-authored kickoff (§12.5).
func (d Driver) seedNext(from, to Block) (string, error) {
	return SeedBlockPrompt(d.DeckDir, d.Manifest, &from, to, d.Now)
}

// SeedBlockPrompt writes a block's driver-authored kickoff (00-prompt.md) into
// its engine workspace (§12.5). When prior is nil the block is the pipeline's
// first and is seeded from the manifest alone. The driver authors kickoff
// material only; it never writes a participant's round/consensus/signoff/final
// artifact.
func SeedBlockPrompt(deckDir string, m Manifest, prior *Block, to Block, now time.Time) (string, error) {
	ws := BlockWorkspace(deckDir, m.IdeaSlug, to.ID)
	if err := os.MkdirAll(filepath.Join(ws, "round-01"), 0o755); err != nil {
		return "", fmt.Errorf("create block workspace: %w", err)
	}
	derivedFrom := "-"
	var inputs strings.Builder
	for _, in := range to.InputArtifacts {
		fmt.Fprintf(&inputs, "  - %s\n", in)
	}
	if prior != nil {
		priorArtifact := filepath.Join(BlockWorkspace(deckDir, m.IdeaSlug, prior.ID), valueOr(prior.OutputArtifact, "FINAL.md"))
		derivedFrom = priorArtifact
		if inputs.Len() == 0 {
			fmt.Fprintf(&inputs, "  - %s\n", priorArtifact)
		}
	}
	if inputs.Len() == 0 {
		fmt.Fprintf(&inputs, "  - (none; first block — work from the pipeline idea)\n")
	}
	origin := "the previous block's finalized output"
	if prior == nil {
		origin = "the pipeline manifest (first block)"
	}
	prompt := fmt.Sprintf(`---
idea: %s__%s
author: driver
created: %s
participants: [%s]
status: round-01
pipeline_slug: %s
block_id: %s
artifact_kind: %s
derived_from: [%s]
---

## Problem / idea

This is pipeline block %q (kind: %s, stage: %s) of pipeline %q, seeded automatically by the driver from %s. Produce the block's output artifact: %s.

## Inputs (read these)

%s
## Block contract

- role lens: %s
- output artifact: %s
- risk: %s

Run the normal Phase 1-4 cooperation rounds and reach consensus. Each participant writes only its own artifact; the driver does not write participant content.
`,
		m.IdeaSlug, to.ID,
		now.Format("2006-01-02"),
		strings.Join(m.Participants, ", "),
		m.IdeaSlug, to.ID, valueOr(to.Stage, string(to.Kind)),
		derivedFrom,
		to.ID, to.Kind, valueOr(to.Stage, "-"), m.IdeaSlug, origin, valueOr(to.OutputArtifact, "FINAL.md"),
		inputs.String(),
		valueOr(to.RoleLens, "-"), valueOr(to.OutputArtifact, "FINAL.md"), riskOrNormal(to.Risk),
	)
	path := filepath.Join(ws, "00-prompt.md")
	if err := os.WriteFile(path, []byte(prompt), 0o644); err != nil {
		return "", fmt.Errorf("write seeded prompt: %w", err)
	}
	return path, nil
}

func (d Driver) blockIndex(id string) int {
	for i, b := range d.Manifest.Blocks {
		if b.ID == id {
			return i
		}
	}
	return -1
}

func (r *PipelineRun) recordCompleted(id string) {
	for _, c := range r.CompletedBlocks {
		if c == id {
			return
		}
	}
	r.CompletedBlocks = append(r.CompletedBlocks, id)
	r.UpdatedAt = nowOr(r.UpdatedAt)
}

func lastBlockID(m Manifest) string {
	if len(m.Blocks) == 0 {
		return ""
	}
	return m.Blocks[len(m.Blocks)-1].ID
}

func valueOr(v, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return v
}

func nowOr(t time.Time) time.Time {
	if t.IsZero() {
		return time.Now()
	}
	return t
}
