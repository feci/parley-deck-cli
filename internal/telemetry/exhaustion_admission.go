package telemetry

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"parley-deck-cli/internal/quota"
	"path/filepath"
	"reflect"
)

type ParticipantValidationRecord struct {
	InvocationID   string `json:"invocation_id"`
	Step           string `json:"step"`
	TerminalSHA256 string `json:"terminal_sha256"`
	Valid          bool   `json:"valid"`
	Reason         string `json:"reason,omitempty"`
	Integrity      string `json:"integrity,omitempty"`
	SourceSHA256   string `json:"source_sha256,omitempty"`
}

func ParticipantTerminalHash(r Record) string {
	raw, _ := json.Marshal(r)
	return fmt.Sprintf("%x", sha256.Sum256(raw))
}

func ReadParticipantValidation(root, key string, r Record) (*ParticipantValidationRecord, error) {
	path := filepath.Join(root, ".parley-runtime", "invocations", r.InvocationID, "participant-result.json")
	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var v ParticipantValidationRecord
	if quota.DecodeRevisionRequest(raw, &v) != nil || v.InvocationID != r.InvocationID || v.Step != key || v.TerminalSHA256 != ParticipantTerminalHash(r) {
		return nil, fmt.Errorf("contradictory participant validation receipt")
	}
	return &v, nil
}

func InvocationExhaustion(r Record, v *ParticipantValidationRecord) (*quota.Evidence, error) {
	if !r.Metadata.ExhaustionPolicy || r.Outcome == nil || r.Outcome.ExhaustionEvidence == nil || !r.Outcome.ExhaustionEvidence.Eligible {
		return nil, nil
	}
	if v == nil || v.Integrity != "" {
		return nil, fmt.Errorf("provider exhaustion validation pending or corrupt: %s", r.InvocationID)
	}
	if v.Valid {
		return nil, nil
	}
	a, err := ParticipantAttempt(r, false, v.Reason)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, fmt.Errorf("exhaustion terminal lacks failed attempt")
	}
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	e := *r.Outcome.ExhaustionEvidence
	e.Exhaustion = &quota.ExhaustedStep{Idea: r.Metadata.Idea, Agent: r.Metadata.Agent,
		Step: r.Metadata.ParticipantLogicalStep, Attempt: *a, TerminalSHA256: ParticipantTerminalHash(r),
		ValidationSHA256: fmt.Sprintf("%x", sha256.Sum256(raw)), CLIVersion: r.Metadata.CLIVersion, ExplicitJSON: r.Metadata.ExplicitJSON}
	if err := quota.ValidateExhaustionEvidence(e, r.Metadata.Idea, r.Metadata.Agent); err != nil {
		return nil, err
	}
	return &e, nil
}

// ExhaustionAdmission derives suppression from the existing invocation ledger.
// It runs regardless of the current policy or membership transition outcome.
// Recovery validates old output through RunParticipantStep before this gate.
func ExhaustionAdmission(root, idea, agent string) error {
	if idea == "" || agent == "" {
		return nil
	}
	paths, err := filepath.Glob(filepath.Join(root, ".parley-runtime", "invocations", "*", "requested.json"))
	if err != nil {
		return err
	}
	var failures, successes []Record
	for _, path := range paths {
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var req Record
		if json.Unmarshal(raw, &req) != nil {
			return fmt.Errorf("unreadable invocation request during exhaustion admission")
		}
		if req.Metadata.Idea != idea || req.Metadata.Agent != agent || !req.Metadata.ExhaustionPolicy {
			continue
		}
		if req.InvocationID != filepath.Base(filepath.Dir(path)) || req.Type != "invocation.requested" || req.Metadata.ParticipantStep == "" {
			return fmt.Errorf("invalid exhaustion invocation binding")
		}
		raw, err = os.ReadFile(filepath.Join(filepath.Dir(path), "terminal.json"))
		if err != nil {
			return fmt.Errorf("exhaustion admission awaits settled invocation %s: %w", req.InvocationID, err)
		}
		var r Record
		if json.Unmarshal(raw, &r) != nil || r.Type != "invocation.terminal" || r.InvocationID != req.InvocationID || !reflect.DeepEqual(r.Metadata, req.Metadata) || r.Outcome == nil || r.CompletedAt == nil || r.CompletedAt.Before(req.RequestedAt) {
			return fmt.Errorf("contradictory exhaustion terminal evidence")
		}
		// A refused launch never started a child or created validation work.
		if !r.Outcome.DispatchAttempted && r.StartedAt == nil {
			continue
		}
		v, err := ReadParticipantValidation(root, r.Metadata.ParticipantStep, r)
		if err != nil {
			return err
		}
		if v == nil || v.Integrity != "" {
			return fmt.Errorf("exhaustion admission requires validation recovery for %s; no new child authorized", r.InvocationID)
		}
		if v.Valid {
			successes = append(successes, r)
		}
		if e, err := InvocationExhaustion(r, v); err != nil {
			return err
		} else if e != nil {
			failures = append(failures, r)
		}
	}
	for _, failure := range failures {
		won := false
		for _, success := range successes {
			if success.StartedAt != nil && success.StartedAt.Before(*failure.CompletedAt) && success.CompletedAt.After(*failure.CompletedAt) {
				won = true
				break
			}
		}
		if !won {
			return fmt.Errorf("provider exhaustion permanently suppresses %s in idea %s (invocation %s); no retry, probe, or reconnection", agent, idea, failure.InvocationID)
		}
	}
	return nil
}
