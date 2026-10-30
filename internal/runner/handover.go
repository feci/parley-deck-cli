package runner

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"parley-deck-cli/internal/membership"
	"parley-deck-cli/internal/protocol"
)

// RunRoleHandover lets only the selected successor author the takeover record.
// It does not authorize code changes or another agent's artifact edits.
func RunRoleHandover(ctx context.Context, opts Options) error {
	ctx, release, err := membership.Acquire(ctx, opts.Idea.Path, opts.RunID)
	if err != nil {
		return err
	}
	defer release()
	if _, err = membership.Before(ctx, opts.Root, opts.Idea.Path, opts.RunID); err != nil {
		return err
	}
	a, err := protocol.ResolveActiveExecution(opts.Idea.Path)
	var pending *protocol.PendingHandover
	if !errors.As(err, &pending) {
		return err
	}
	r := pending.Transition
	selected, _ := selectedAgents([]string{r.Successor}, opts.Agents, resolveMapping(opts))
	if len(selected) != 1 {
		return fmt.Errorf("selected handover successor %s is unavailable", r.Successor)
	}
	agent := selected[0]
	if err = os.MkdirAll(filepath.Dir(a.Path), 0700); err != nil {
		return err
	}
	logs := filepath.Join(opts.Root, ".parley-runtime", "handovers", r.ID)
	if err = os.MkdirAll(logs, 0700); err != nil {
		return err
	}
	prompt := fmt.Sprintf(`You are %s, the selected successor for Parley Deck idea %s.
Create only %s. Do not change code or any other artifact.
Read FINAL.md, the canonical phase artifacts, and the complete predecessor execution chain.
Preserve every predecessor byte, dissent, finding, and validation result.
Record the takeover and the next action in the new file.
Use this exact frontmatter authority:
---
idea: %s
transition: %s
author: %s
organizer: %s
implementer: %s
phase: %s
base: %s
base-sha256: %s
takeover-branch: %s
takeover-head: %s
date: %s
status: in-progress
---
Include non-empty sections '## Takeover' and '## Resume state'. State the affected roles: %v.
Include the next action and exact canonical resume paths.
If implementation has started, include all ordinary Phase 5 fields and sections.
Retain the prior progress, cycle count, evidence, and findings through the base chain.
If base is none and implementation has not started, do not claim review readiness.
Never review code that you implement. Return a short path confirmation.
`, r.Successor, opts.Idea.Slug, a.Path, opts.Idea.Slug, r.ID, r.Successor, r.After.Organizer, r.After.Implementer, r.Phase, r.Base, r.BaseSHA256, r.Branch, r.Head, time.Now().UTC().Format("2006-01-02"), r.Roles) + protocol.OutputTextInstruction
	ceiling := timeoutForAgent(opts.Timeout, agent)
	step := "handover:" + r.ID
	_, err = RunParticipantStep(ctx, ParticipantStepOptions{Root: opts.Root, Idea: opts.Idea.Slug, Agent: agent, Step: step, SingleAttempt: true, HardTimeout: ceiling, Files: []string{a.Path}, Validate: func() StepValidation {
		if e := protocol.ValidateRoleHandover(a.Path, opts.Idea.Slug, *r); e != nil {
			return StepValidation{Reason: e.Error()}
		}
		if _, e := protocol.ResolveActiveExecution(opts.Idea.Path); e != nil {
			return StepValidation{Integrity: e}
		}
		raw, e := os.ReadFile(a.Path)
		if e != nil {
			return StepValidation{Integrity: e}
		}
		return StepValidation{Valid: true, SourceSHA256: fmt.Sprintf("%x", sha256.Sum256(raw))}
	}}, func(child context.Context, _ int, _ string) error {
		child = WithLaunchInfo(child, LaunchInfo{RunID: opts.RunID, Idea: opts.Idea.Slug, Phase: "implementation", ArtifactPath: a.Path, Store: opts.Store})
		res := RunConsult(child, ConsultOptions{Root: opts.Root, Agent: agent, Prompt: prompt, Timeout: ceiling,
			StdoutPath: filepath.Join(logs, "stdout.log"), StderrPath: filepath.Join(logs, "stderr.log")})
		if res.ExitError != "" {
			return fmt.Errorf("%s", res.ExitError)
		}
		return nil
	})
	return err
}
