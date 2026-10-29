package runner

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"

	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/quota"
	"parley-deck-cli/internal/telemetry"
)

func runExhaustionFixup(ctx context.Context, opts Options) Result {
	selected, _ := selectedAgents(opts.Idea.Participants, opts.Agents, resolveMapping(opts))
	if len(selected) != 1 {
		return Result{AgentID: "implementer", ExitError: "fix-up requires one selected implementer"}
	}
	agent := selected[0]
	contract, err := os.ReadFile(filepath.Join(opts.Idea.Path, "review", "consensus.md"))
	if err != nil {
		return Result{AgentID: agent.ID, ExitError: err.Error()}
	}
	step := fmt.Sprintf("fixup:%x", sha256.Sum256(contract))
	path, err := protocol.ExecutionPath(opts.Idea.Path)
	if err != nil {
		return Result{AgentID: agent.ID, ExitError: err.Error()}
	}
	var result Result
	settled, err := RunParticipantStep(ctx, ParticipantStepOptions{Root: opts.Root, Idea: opts.Idea.Slug, Agent: agent,
		Step: step, SingleAttempt: true, HardTimeout: timeoutForAgent(opts.Timeout, agent), Files: []string{path},
		Validate: func() StepValidation {
			if err := ValidateFixupArtifact(opts.Idea.Path, opts.Idea.Slug); err != nil {
				return StepValidation{Reason: err.Error()}
			}
			return freshExecutionArtifact(opts.Root, opts.Idea.Slug, agent.ID, step, path)
		},
	}, func(child context.Context, _ int, _ string) error {
		result = RunFixup(child, opts)
		if result.ExitError != "" {
			return fmt.Errorf("%s", result.ExitError)
		}
		return nil
	})
	result.AgentID, result.OutputPath = agent.ID, path
	result.InvocationID, result.QuotaEvidence, result.ArtifactOK = settled.Record.InvocationID, settled.Evidence, settled.Valid
	if err != nil {
		result.ExitError = err.Error()
	} else if settled.Valid {
		result.ExitError = ""
	}
	opts.Phase, opts.RoundLabel = "fixup", "fixup"
	return quotaSettleSingle(ctx, opts, result)
}

func participantArtifactStep(opts Options, agent string) string {
	if opts.ArtifactName != "" {
		if opts.exhaustionPolicy && opts.Phase == "review-consensus" {
			return fmt.Sprintf("review-consensus/%d:%s", opts.Round, filepath.ToSlash(opts.ArtifactName))
		}
		return opts.RoundLabel + ":" + filepath.ToSlash(opts.ArtifactName)
	}
	return filepath.ToSlash(filepath.Join(opts.RoundLabel, agent+".md"))
}

// A prior execution record is input, not output from the new code step.
func freshExecutionArtifact(root, idea, agent, step, path string) StepValidation {
	raw, err := os.ReadFile(path)
	if err != nil {
		return StepValidation{Reason: err.Error()}
	}
	before, err := ParticipantBaseline(root, idea, agent, step, path)
	if err != nil && !os.IsNotExist(err) {
		return StepValidation{Integrity: err}
	}
	if err == nil && bytes.Equal(before, raw) {
		return StepValidation{Reason: "unchanged prior execution artifact is not output from this code step"}
	}
	return StepValidation{Valid: true, SourceSHA256: fmt.Sprintf("%x", sha256.Sum256(raw))}
}

func invocationExhaustion(r telemetry.Record, v *participantValidationRecord) (*quota.Evidence, error) {
	return telemetry.InvocationExhaustion(r, v)
}

func exhaustionAdmission(root, idea, agent string) error {
	return telemetry.ExhaustionAdmission(root, idea, agent)
}
